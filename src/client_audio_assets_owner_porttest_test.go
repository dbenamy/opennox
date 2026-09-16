//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/timer"
	"os"
	"sort"
	"strings"
	"testing"
	"unsafe"
)

type clientAudioAssetsOwner struct {
	rows, base, scratch []byte
	enabled             *uint32
	catalog             *unsafe.Pointer
	owned               *audioStructXxx
	ticks               int
}

func newClientAudioAssetsOwner(t *testing.T) *clientAudioAssetsOwner {
	t.Helper()
	o := new(clientAudioAssetsOwner)
	oldClock := timer.PlatformTicks
	timer.PlatformTicks = func() uint64 { o.ticks++; return 123456789 }
	var restore func()
	o.rows, o.enabled, o.catalog, restore = legacy.PortTestClientAudioAssetsOwner()
	// Real table/default-timer/name initialization, without playback contexts.
	legacy.Sub_451850(nil, nil)
	o.base = append([]byte(nil), o.rows...)
	var freeScratch func()
	o.scratch, freeScratch = alloc.Make([]byte{}, 256*1024)
	t.Cleanup(func() {
		restore()
		if o.owned != nil {
			o.owned.Free()
		}
		freeScratch()
		timer.PlatformTicks = oldClock
	})
	o.setCatalog(t, nil)
	o.reset()
	return o
}
func (o *clientAudioAssetsOwner) reset() {
	copy(o.rows, o.base)
	*o.enabled = 1
	o.ticks = 0
	for i := range o.scratch {
		o.scratch[i] = 0xa5
	}
}
func (o *clientAudioAssetsOwner) setCatalog(t *testing.T, names []string) {
	t.Helper()
	*o.catalog = nil
	if o.owned != nil {
		o.owned.Free()
	}
	p, _ := alloc.New(audioStructXxx{})
	o.owned = p
	p.size4 = uint(len(names))
	if len(names) > 0 {
		rows, _ := alloc.Make([]audioStructYyy{}, len(names))
		p.arr0 = &rows[0]
		for i, name := range names {
			if len(name) >= 24 {
				t.Fatal("fixture sample name exceeds metadata boundary")
			}
			copy(rows[i].field0[:24], name)
			binary.LittleEndian.PutUint32(rows[i].field0[24:], uint32(i)^0x76543210)
			binary.LittleEndian.PutUint32(rows[i].field0[28:], uint32(i)*17+3)
			rows[i].field32 = uint32(i) ^ 0xabcdef01
		}
		sort.SliceStable(rows, func(i, j int) bool {
			return strings.ToLower(alloc.GoStringS(rows[i].field0[:])) < strings.ToLower(alloc.GoStringS(rows[j].field0[:]))
		})
	}
	*o.catalog = p.C()
}
func (o *clientAudioAssetsOwner) metadata(t *testing.T) []byte {
	t.Helper()
	out := append([]byte(nil), o.rows...)
	for i := 0; i < 1023; i++ {
		off := 200*i + 84
		if !bytes.Equal(o.rows[off:off+4], o.base[off:off+4]) {
			t.Fatalf("sound%d name pointer changed", i)
		}
		clear(out[off : off+4])
	}
	return out
}
func (o *clientAudioAssetsOwner) invoke(t *testing.T, op int, payload []byte) (ret, consumed int) {
	t.Helper()
	input := append(append([]byte(nil), payload...), 0x73, 0xa9)
	raw, _ := alloc.CloneSlice(input)
	f := binfile.NewMemFile(unsafe.Pointer(&raw[0]), len(raw))
	defer f.Free()
	switch op {
	case 0:
		ret = legacy.Nox_thing_read_audio_415660(f, o.scratch)
	case 1:
		ret = legacy.Nox_thing_read_AVNT_452890(f, o.scratch)
	case 2:
		ret = legacy.PortTestClientAudioRecord(f, o.scratch)
	default:
		t.Fatal("reader op")
	}
	if !bytes.Equal(raw, input) {
		t.Fatal("reader mutated input")
	}
	return ret, len(raw) - len(f.Data())
}
func clientAudioAssetsCapture(t *testing.T, label string, rows any, want string) {
	t.Helper()
	b, e := json.Marshal(rows)
	if e != nil {
		t.Fatal(e)
	}
	if p := os.Getenv("OPENNOX_CLIENT_AUDIO_ASSETS_CAPTURE"); p != "" {
		if e := os.WriteFile(p+"-"+label+".json", b, 0600); e != nil {
			t.Fatal(e)
		}
	}
	got := fmt.Sprintf("%x", sha256.Sum256(b))
	t.Logf("%s %s", label, got)
	if want != "" && got != want {
		t.Fatalf("%s got%s want%s", label, got, want)
	}
}
func clientAudioAssetsName(b []byte, name string) []byte {
	b = append(b, byte(len(name)))
	return append(b, []byte(name)...)
}
func clientAudioAssetsWord(b []byte, v uint32) []byte { return binary.LittleEndian.AppendUint32(b, v) }
func clientAudioAssetsShort(b []byte, v int16) []byte {
	return binary.LittleEndian.AppendUint16(b, uint16(v))
}
