//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
)

func TestObjectXferInventoryLoading(t *testing.T) {
	core := newObjectXferOwner(t)
	path := filepath.Join(t.TempDir(), "inventory.bin")
	typ := core.Types.ByID("PortElevatorShaft")
	if typ == nil {
		t.Fatal("inventory type")
	}
	for _, version := range []int{0, 59, 60, 65535} {
		for _, count := range []int{0, 1, 2, 7} {
			for _, missing := range []bool{false, true} {
				if count == 0 && missing {
					continue
				}
				t.Run(fmt.Sprintf("v%d-count%d-missing%v", version, count, missing), func(t *testing.T) {
					parent := newObjectXferSimple(t, core)
					before := core.Objs.Alive
					var stream mapDrawableStream
					expectedPos := 0
					for i := 0; i < count; i++ {
						bad := missing && i == count-1
						if version < 60 {
							name := "PortElevatorShaft"
							if bad {
								name = "NoSuchFixtureObject"
							}
							stream.u8(byte(len(name)))
							stream.WriteString(name)
						} else {
							code := uint16(1000 + typ.Ind())
							if bad {
								code = 999
							}
							stream.u16(code)
						}
						if bad {
							expectedPos = stream.Len()
							stream.u32(0xdeadbeef)
							break
						}
						var record mapDrawableStream
						record.u16(60)
						objectXferEmptyBase(&record, 60)
						record.u32(uint32(101 + i))
						stream.u32(uint32(record.Len()))
						stream.Write(record.Bytes())
						expectedPos = stream.Len()
					}
					if err := os.WriteFile(path, stream.Bytes(), 0600); err != nil {
						t.Fatal(err)
					}
					if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
						t.Fatal(err)
					}
					defer cryptfile.Close()
					got := legacy.Nox_xxx_xfer_4F3E30(version, parent, uint32(count))
					pos, err := cryptfile.Global().File.Seek(0, io.SeekCurrent)
					wantReturn, loaded := 1, count
					if missing {
						wantReturn = 0
						loaded--
					}
					if got != wantReturn || err != nil || pos != int64(expectedPos) {
						t.Fatalf("ret=%d want=%d pos=%d want=%d err=%v", got, wantReturn, pos, expectedPos, err)
					}
					if core.Objs.Alive != before+loaded {
						t.Fatalf("live objects=%d want=%d", core.Objs.Alive, before+loaded)
					}
					it := parent.InvFirstItem
					for i := loaded - 1; i >= 0; i-- {
						if it == nil {
							t.Fatal("short inventory chain")
						}
						trackObjectXferTyped(t, core, it)
						value := binary.LittleEndian.Uint32(unsafe.Slice((*byte)(unsafe.Add(it.UpdateData, 8)), 4))
						if it.InvHolder != parent || value != uint32(101+i) {
							t.Fatal("inventory holder/order/data")
						}
						if it.InvNextItem != nil && it.InvNextItem.Field125 != it {
							t.Fatal("inventory back link")
						}
						if i == loaded-1 && it.Field125 != nil {
							t.Fatal("head back link")
						}
						it = it.InvNextItem
					}
					if it != nil {
						t.Fatal("long inventory chain")
					}
					parent.InvFirstItem = nil
				})
			}
		}
	}
}

func TestObjectXferInventoryFactoryRejection(t *testing.T) {
	core := newObjectXferOwner(t)
	path := filepath.Join(t.TempDir(), "factory.bin")
	for _, count := range []int32{-2147483648, -1, 0, 1} {
		for _, code := range []uint16{900, 999} {
			t.Run(fmt.Sprintf("count%d-code%d", count, code), func(t *testing.T) {
				u := newObjectXferSimple(t, core)
				before := core.Objs.Alive
				var stream mapDrawableStream
				stream.u16(code)
				stream.u32(0x12345678)
				if err := os.WriteFile(path, stream.Bytes(), 0600); err != nil {
					t.Fatal(err)
				}
				if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
					t.Fatal(err)
				}
				defer cryptfile.Close()
				got := legacy.Nox_xxx_xfer_4F3E30(60, u, uint32(count))
				want, posWant := 1, int64(0)
				if count > 0 {
					want = 0
					posWant = 2
					if code == 900 {
						posWant = 6
					}
				}
				pos, err := cryptfile.Global().File.Seek(0, io.SeekCurrent)
				if got != want || pos != posWant || err != nil || core.Objs.Alive != before || u.InvFirstItem != nil {
					t.Fatalf("factory result=%d/%d pos=%d/%d alive=%d/%d err=%v", got, want, pos, posWant, core.Objs.Alive, before, err)
				}
			})
		}
	}
}
