//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestThingSkipsWallScratchCapacity(t *testing.T) {
	type record struct {
		Length, Capacity, Return, Consumed int
		Panicked                           bool
		Scratch                            [32]byte
	}
	var rows []record
	for _, capacity := range []int{256*1024 - 1, 256 * 1024} {
		for _, length := range []int{0, 1, 31, 64, 256 * 1024} {
			if length > capacity {
				continue
			}
			func() {
				input, wantPrefix, _, consumed := thingSkipWallRecord(0, 8, 31, 2, 2, 0x454e4420)
				raw, _ := alloc.CloneSlice(input)
				f := binfile.NewMemFile(unsafe.Pointer(&raw[0]), len(raw))
				defer f.Free()
				scratch, free := alloc.Make([]byte{}, 256*1024)
				defer free()
				for i := range scratch {
					scratch[i] = 0xa5
				}
				want := append([]byte(nil), scratch...)
				panicked := false
				ret := 0
				func() {
					defer func() {
						if recover() != nil {
							panicked = true
						}
					}()
					ret = legacy.Nox_thing_read_WALL_414F60(f, scratch[:length:capacity])
				}()
				wantPanic := capacity < 256*1024 || length == 0
				pos := len(raw) - len(f.Data())
				wantPos, wantRet := 0, 0
				if !wantPanic {
					copy(want, wantPrefix)
					wantPos, wantRet = consumed, 1
				}
				if panicked != wantPanic || ret != wantRet || pos != wantPos || !bytes.Equal(scratch, want) || !bytes.Equal(raw, input) {
					t.Fatalf("length%d capacity%d panic%v/%v return%d/%d cursor%d/%d", length, capacity, panicked, wantPanic, ret, wantRet, pos, wantPos)
				}
				rows = append(rows, record{length, capacity, ret, pos, panicked, sha256.Sum256(scratch)})
			}()
		}
	}
	thingSkipCapture(t, "scratch-capacity", rows, "6bfdd59d5861b06669856b724decc01a413bf782ba2419f143bb8d6217eafe38")
}
