//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"strings"
	"testing"
	"unsafe"
)

func TestClientAudioAssetsGetters(t *testing.T) {
	o := newClientAudioAssetsOwner(t)
	type record struct {
		Enabled     uint32
		ID          int32
		Word        uint32
		Slot, Delay int32
	}
	var rows []record
	for _, enabled := range []uint32{0, 1, 0xffffffff} {
		for _, id := range []int32{-2147483648, -2, -1, 0, 1, 2, 511, 1022, 1023, 2147483647} {
			for _, word := range []uint32{0, 1, 0x7fffffff, 0x80000000, 0xffffffff} {
				o.reset()
				*o.enabled = enabled
				binary.LittleEndian.PutUint32(o.rows[200+64:], word)
				before := sha256.Sum256(o.rows)
				p := legacy.PortTestClientAudioSlot(id)
				slot := int32(-1)
				if enabled != 0 && id >= 0 && id < 1023 {
					if p != unsafe.Pointer(&o.rows[200*int(id)]) {
						t.Fatal("wrong actual sound slot")
					}
					slot = id
				} else if p != nil {
					t.Fatal("invalid sound slot was accepted")
				}
				arg := unsafe.Pointer(&o.rows[200])
				if enabled == 0 {
					arg = nil
				}
				delay := legacy.PortTestClientAudioDelay(arg)
				want := int32(0)
				if enabled != 0 {
					want = int32(word)
				}
				if delay != want || before != sha256.Sum256(o.rows) || o.ticks != 0 {
					t.Fatal("getter changed owner or delay")
				}
				rows = append(rows, record{enabled, id, word, slot, delay})
			}
		}
	}
	clientAudioAssetsCapture(t, "getters", rows, "38278328bf4a85afca9e48e552d6960c60da3c68d6fc7a6dc64a997bd0bfaa06")
}
func TestClientAudioAssetsSampleLookup(t *testing.T) {
	o := newClientAudioAssetsOwner(t)
	type record struct {
		Group, Key string
		Count      int
		Index      int32
	}
	var rows []record
	sets := [][]string{{}, {"Alpha"}, {"", "Alpha"}, {"alpha", "ALPHA", "Alpha", "beta", "Beta", "zulu"}}
	for _, count := range []int{3, 16, 33, 255, 65537} {
		names := make([]string, count)
		for i := range names {
			names[i] = fmt.Sprintf("S%05d", i)
		}
		sets = append(sets, names)
	}
	for group, names := range sets {
		o.reset()
		o.setCatalog(t, names)
		catalog := unsafe.Slice(o.owned.arr0, int(o.owned.size4))
		var raw []byte
		if len(catalog) > 0 {
			raw = unsafe.Slice((*byte)(unsafe.Pointer(&catalog[0])), len(catalog)*36)
		}
		before := sha256.Sum256(raw)
		definitions := sha256.Sum256(o.rows)
		keys := []string{"", "missing", "ALPHA", "aLpHa", "beta", "ZULU", "S32767", "S32768", "S65534", "S65535", "S65536"}
		for _, i := range []int{0, 1, len(catalog) / 2, len(catalog) - 1} {
			if i >= 0 && i < len(catalog) {
				keys = append(keys, alloc.GoStringS(catalog[i].field0[:]))
			}
		}
		for _, key := range keys {
			data, free := alloc.CloneSlice(append([]byte(key), 0))
			original := sha256.Sum256(data)
			index := legacy.PortTestClientAudioSample(o.owned.C(), &data[0])
			if original != sha256.Sum256(data) {
				t.Fatal("lookup mutated key")
			}
			free()
			found := false
			for _, r := range catalog {
				if strings.EqualFold(alloc.GoStringS(r.field0[:]), key) {
					found = true
					break
				}
			}
			if found {
				if index < 0 || int(index) >= len(catalog) || !strings.EqualFold(alloc.GoStringS(catalog[index].field0[:]), key) {
					t.Fatal("lookup did not return a matching row")
				}
			} else if index != -1 {
				t.Fatal("missing sample lookup succeeded")
			}
			if before != sha256.Sum256(raw) || definitions != sha256.Sum256(o.rows) || o.ticks != 0 {
				t.Fatal("lookup changed actual owners")
			}
			rows = append(rows, record{fmt.Sprint(group), key, len(catalog), index})
		}
	}
	clientAudioAssetsCapture(t, "sample-lookup", rows, "4147b1f864904f88e969faa847245d15e844b50d5ff6f2039d3c96a65e8b28e6")
}
