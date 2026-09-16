//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/sound"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

type clientAudioAssetsReadResult struct {
	Label                 string
	Return, Cursor, Ticks int
	Metadata, Scratch     [32]byte
}

func (o *clientAudioAssetsOwner) observe(t *testing.T, label string, op int, payload []byte) clientAudioAssetsReadResult {
	t.Helper()
	ret, pos := o.invoke(t, op, payload)
	if pos > len(payload) {
		t.Fatalf("%s consumed guard bytes: %d/%d", label, pos, len(payload))
	}
	return clientAudioAssetsReadResult{label, ret, pos, o.ticks, sha256.Sum256(o.metadata(t)), sha256.Sum256(o.scratch[:cap(o.scratch)])}
}
func clientAudioAssetsRecord(name string, priority int16, volume byte, delay int16, fields [4]byte, samples []string) []byte {
	b := clientAudioAssetsName(nil, name)
	b = clientAudioAssetsShort(b, priority)
	b = append(b, volume)
	b = clientAudioAssetsShort(b, delay)
	b = append(b, fields[:]...)
	for _, s := range samples {
		b = clientAudioAssetsName(b, s)
	}
	return append(b, 0)
}
func TestClientAudioAssetsRecords(t *testing.T) {
	o := newClientAudioAssetsOwner(t)
	o.setCatalog(t, []string{"Alpha", "beta", "multi.dot", ""})
	known := sound.ID(1).String()
	var rows []clientAudioAssetsReadResult
	for _, name := range []string{known, "missing", "", known + "\x00tail"} {
		for _, enabled := range []uint32{0, 1} {
			for _, v := range []int{0, 1, 2, 3, 127, 128, 255} {
				o.reset()
				*o.enabled = enabled
				p := clientAudioAssetsRecord(name, int16(v*257), byte(v), int16(v*257), [4]byte{byte(v), byte(255 - v), byte(v), byte(v)}, []string{"alpha.wav", "missing.wav", "BETA", "multi.dot.wav"})
				r := o.observe(t, fmt.Sprintf("%q/%d/%d", name, enabled, v), 2, p)
				accepted := sound.ByName(strings.Split(name, "\x00")[0]) != 0 && enabled != 0
				if !accepted {
					if r.Return != 1 || r.Cursor != len(p) || r.Ticks != 0 || !bytes.Equal(o.rows, o.base) {
						t.Fatal("unknown/disabled record mutated definitions or failed to skip")
					}
				} else {
					if r.Ticks != 1 {
						t.Fatal("record must initialize its real timer once")
					}
					if int8(v) >= 3 {
						if r.Return != 0 || r.Cursor != len(name)+10 {
							t.Fatal("mode rejection cursor/return")
						}
					} else if r.Return != 1 || r.Cursor != len(p) {
						t.Fatal("accepted record did not finish")
					}
					off := 200 * int(sound.ByName(strings.Split(name, "\x00")[0]))
					if int32(binary.LittleEndian.Uint32(o.rows[off+56:])) != int32(int8(v)) {
						t.Fatal("signed field widened incorrectly")
					}
				}
				rows = append(rows, r)
			}
		}
	}
	// Independent delay/priority signed-short boundaries, not tied to mode selection.
	for _, v := range []int16{-32768, -1, 0, 1, 32767} {
		o.reset()
		off := 200
		binary.LittleEndian.PutUint32(o.rows[off+64:], 321)
		p := clientAudioAssetsRecord(known, v, 255, v, [4]byte{}, nil)
		rows = append(rows, o.observe(t, fmt.Sprint("short/", v), 2, p))
		want := uint32(321)
		if v > 0 {
			want = uint32(v) * 15
		}
		if binary.LittleEndian.Uint32(o.rows[off+64:]) != want || int32(binary.LittleEndian.Uint32(o.rows[off+8:])) != int32(v) {
			t.Fatal("signed priority/delay contract")
		}
	}
	clientAudioAssetsCapture(t, "records", rows, "d629e1dcef6d3d8c855c9724c4619843f16d053247598186c36289b49b62c8f2")
}
func clientAudioAssetsTag(tag byte, v byte) []byte {
	b := []byte{tag}
	switch tag {
	case 1, 2, 3, 4, 5:
		b = append(b, v)
	case 6:
		b = append(b, v, 255-v)
	case 7:
		for _, s := range []string{"missing", "ALPHA", "beta"} {
			b = clientAudioAssetsName(b, s)
		}
		b = append(b, 0)
	case 8:
		b = clientAudioAssetsWord(b, 0x80000000|uint32(v))
		b = clientAudioAssetsWord(b, 0xffffffff-uint32(v))
	case 9, 10:
		b = clientAudioAssetsShort(b, int16(uint16(v)*257))
	}
	return b
}
func TestClientAudioAssetsEvents(t *testing.T) {
	o := newClientAudioAssetsOwner(t)
	o.setCatalog(t, []string{"Alpha", "beta"})
	var rows []clientAudioAssetsReadResult
	for _, known := range []bool{false, true} {
		name := "missing"
		if known {
			name = sound.ID(1).String()
		}
		for _, enabled := range []uint32{0, 1} {
			for tag := 0; tag < 256; tag++ {
				o.reset()
				*o.enabled = enabled
				p := append(clientAudioAssetsName(nil, name), clientAudioAssetsTag(byte(tag), 255)...)
				p = append(p, 0)
				r := o.observe(t, fmt.Sprintf("%t/%d/%d", known, enabled, tag), 1, p)
				if !known || enabled == 0 {
					if r.Return != 1 || r.Ticks != 0 || !bytes.Equal(o.rows, o.base) {
						t.Fatal("skipped AVNT changed state or propagated unknown-tag failure")
					}
				} else if (tag <= 10) != (r.Return == 1) {
					t.Fatal("event tag success contract")
				}
				rows = append(rows, r)
			}
		}
	}
	// Repeated and reordered tags preserve earlier fields and timer call count.
	for _, order := range [][]byte{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, {10, 9, 8, 7, 6, 5, 4, 3, 2, 1}, {3, 9, 3, 7, 7, 9, 3}} {
		for _, v := range []byte{0, 1, 127, 128, 255} {
			o.reset()
			p := clientAudioAssetsName(nil, sound.ID(1).String())
			timers := 0
			for _, tag := range order {
				p = append(p, clientAudioAssetsTag(tag, v)...)
				if tag == 3 {
					timers++
				}
			}
			p = append(p, 0)
			r := o.observe(t, fmt.Sprintf("%v/%d", order, v), 1, p)
			if r.Return != 1 || r.Cursor != len(p) || r.Ticks != timers {
				t.Fatal("event sequence cursor/timer contract")
			}
			rows = append(rows, r)
		}
	}
	clientAudioAssetsCapture(t, "events", rows, "3cecd7b4489458691ff5183af6a56079a61e29d7311ffb4327c4b5a3909f3e9e")
}
func TestClientAudioAssetsSampleLists(t *testing.T) {
	o := newClientAudioAssetsOwner(t)
	names := make([]string, 65537)
	for i := range names {
		names[i] = fmt.Sprintf("S%05d", i)
	}
	o.setCatalog(t, names)
	var rows []clientAudioAssetsReadResult
	for _, op := range []int{1, 2} {
		for _, count := range []int{0, 1, 31, 32, 33, 40} {
			if op == 2 && count > 32 {
				continue
			} // Original AUD valid-input capacity.
			o.reset()
			samples := []string{"missing"}
			for i := 0; i < count; i++ {
				samples = append(samples, names[[]int{0, 32767, 32768, 65534, 65535, 65536}[i%6]])
			}
			var p []byte
			if op == 1 {
				p = append(clientAudioAssetsName(nil, sound.ID(1).String()), 7)
				for _, s := range samples {
					p = clientAudioAssetsName(p, s)
				}
				p = append(p, 0, 0)
			} else {
				for i := range samples {
					samples[i] += ".wav"
				}
				p = clientAudioAssetsRecord(sound.ID(1).String(), 0, 0, 0, [4]byte{}, samples)
			}
			r := o.observe(t, fmt.Sprintf("%d/%d", op, count), op, p)
			expected := count - count/6 // 65535 narrows to sentinel -1; index65536 becomes0.
			if count%6 >= 5 {
				expected--
			}
			if op == 1 && expected > 32 {
				expected = 32
			}
			if r.Return != 1 || r.Cursor != len(p) || int(binary.LittleEndian.Uint32(o.rows[200+192:])) != expected {
				t.Fatalf("narrowing/count contract %d/%d", op, count)
			}
			rows = append(rows, r)
		}
	}
	// Long unsigned AVNT names and short signed AUD names, embedded NUL and extension tails.
	o.setCatalog(t, []string{"Alpha", "multi.dot"})
	for _, sample := range []string{"Alpha\x00ignored.wav", "multi.dot.wav", strings.Repeat("x", 127), strings.Repeat("x", 255)} {
		for _, op := range []int{1, 2} {
			if op == 2 && len(sample) > 127 {
				continue
			}
			o.reset()
			var p []byte
			if op == 1 {
				p = append(clientAudioAssetsName(nil, sound.ID(1).String()), 7)
				p = clientAudioAssetsName(p, sample)
				p = append(p, 0, 0)
			} else {
				p = clientAudioAssetsRecord(sound.ID(1).String(), 0, 0, 0, [4]byte{}, []string{sample})
			}
			rows = append(rows, o.observe(t, fmt.Sprintf("tail/%d/%q", op, sample), op, p))
		}
	}
	clientAudioAssetsCapture(t, "sample-lists", rows, "d40055048815871ca71380244ded275e86fea7b927f90bcdd6eb36e00613eb57")
}
func TestClientAudioAssetsBulk(t *testing.T) {
	o := newClientAudioAssetsOwner(t)
	o.setCatalog(t, []string{"Alpha"})
	var rows []clientAudioAssetsReadResult
	for _, count := range []int32{-2147483648, -1, 0, 1, 2, 3} {
		for _, failure := range []int{-1, 0, 1, 2} {
			o.reset()
			p := clientAudioAssetsWord(nil, uint32(count))
			stop := 4
			ret := 1
			calls := 0
			for i := 0; i < 3; i++ {
				mode := byte(0)
				if i == failure {
					mode = 3
				}
				rec := clientAudioAssetsRecord(sound.ID(i+1).String(), int16(i), 1, 1, [4]byte{0, 0, 0, mode}, []string{"Alpha.wav"})
				if int32(i) < count && ret == 1 {
					calls++
					if i == failure {
						stop = len(p) + len(sound.ID(i+1).String()) + 10
						ret = 0
					} else {
						stop = len(p) + len(rec)
					}
				}
				p = append(p, rec...)
			}
			r := o.observe(t, fmt.Sprintf("%d/%d", count, failure), 0, p)
			if r.Return != ret || r.Cursor != stop || r.Ticks != calls {
				t.Fatalf("bulk stop contract: %+v want %d/%d/%d", r, ret, stop, calls)
			}
			rows = append(rows, r)
		}
	}
	clientAudioAssetsCapture(t, "bulk", rows, "f8e256c6cc1d6c5f025c8646eb68e72745aead4e71aa318a3170832eec676556")
}
func TestClientAudioAssetsScratchCapacity(t *testing.T) {
	o := newClientAudioAssetsOwner(t)
	var rows []clientAudioAssetsReadResult
	for _, op := range []int{0, 1} {
		for _, n := range []int{1, 17, 256 * 1024} {
			o.reset()
			o.scratch = o.scratch[:n]
			p := clientAudioAssetsWord(nil, 0)
			if op == 1 {
				p = append(clientAudioAssetsName(nil, sound.ID(1).String()), 0)
			}
			rows = append(rows, o.observe(t, fmt.Sprintf("%d/%d", op, n), op, p))
			o.scratch = o.scratch[:cap(o.scratch)]
		}
		for _, kind := range []int{0, 1} {
			o.reset()
			full := o.scratch
			if kind == 0 {
				o.scratch = full[:0]
			} else {
				o.scratch = full[:1:1]
			}
			before := sha256.Sum256(o.rows)
			func() {
				defer func() {
					if recover() == nil {
						t.Error("invalid scratch accepted")
					}
				}()
				o.invoke(t, op, []byte{0, 0, 0, 0})
			}()
			o.scratch = full
			if before != sha256.Sum256(o.rows) || o.ticks != 0 {
				t.Fatal("rejected scratch changed state")
			}
		}
	}
	clientAudioAssetsCapture(t, "scratch-capacity", rows, "8839d0be9c250f287840dd9c149253425d9c79c2f94363f374672996575d8c64")
}
func TestClientAudioAssetsRealCatalog(t *testing.T) {
	t.Cleanup(handles.PortTestInit())
	assets := os.Getenv("OPENNOX_CLIENT_AUDIO_ASSETS")
	if assets == "" {
		t.Skip("asset path not configured")
	}
	o := newClientAudioAssetsOwner(t)
	catalog := sub_4866F0(filepath.Join(assets, "audio"), "")
	if catalog == nil {
		t.Fatal("real catalog could not load")
	}
	*o.catalog = nil
	o.owned.Free()
	o.owned = catalog
	*o.catalog = catalog.C()
	actual := unsafe.Slice(catalog.arr0, int(catalog.size4))
	if len(actual) == 0 {
		t.Fatal("empty real catalog")
	}
	var rows []clientAudioAssetsReadResult
	for i, r := range actual {
		key := alloc.GoStringS(r.field0[:])
		data, free := alloc.CloneSlice(append([]byte(key), 0))
		got := legacy.PortTestClientAudioSample(catalog.C(), &data[0])
		free()
		if got < 0 || int(got) >= len(actual) || !strings.EqualFold(alloc.GoStringS(actual[got].field0[:]), key) {
			t.Fatalf("real catalog lookup %d failed", i)
		}
	}
	for _, i := range []int{0, len(actual) / 2, len(actual) - 1} {
		key := alloc.GoStringS(actual[i].field0[:])
		for _, op := range []int{1, 2} {
			o.reset()
			var p []byte
			if op == 1 {
				p = append(clientAudioAssetsName(nil, sound.ID(1).String()), 7)
				p = clientAudioAssetsName(p, key)
				p = append(p, 0, 0)
			} else {
				p = clientAudioAssetsRecord(sound.ID(1).String(), 0, 1, 1, [4]byte{}, []string{key + ".wav"})
			}
			r := o.observe(t, fmt.Sprintf("%d/%d", op, i), op, p)
			if r.Return != 1 || r.Cursor != len(p) || binary.LittleEndian.Uint32(o.rows[200+192:]) != 1 {
				t.Fatal("real sample not bound")
			}
			rows = append(rows, r)
		}
	}
	clientAudioAssetsCapture(t, "real-catalog", rows, "bdd84b61caa8affe8e9ed30b173a323356c402f6160905090b1cabdd1e64c953")
}
