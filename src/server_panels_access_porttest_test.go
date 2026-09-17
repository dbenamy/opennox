//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestServerPanelsAccessNumbers(t *testing.T) {
	type row struct {
		ID, Event, Focus int
		Input, Shown     string
		Admission        []byte
		Limit, Updated   uint32
	}
	var rows []row
	cases := []struct {
		text  string
		value int
	}{{"", 0}, {"-1", 0}, {"0", 0}, {"1", 1}, {"14", 14}, {"15", 15}, {"32", 32}, {"33", 33}, {"65535", 65535}, {"65536", 65536}, {" 12", 12}, {"7tail", 7}, {"bad", 0}}
	for _, id := range []int{10126, 10128, 10130, 10132, 10133} {
		for _, event := range []int{16415, 16387} {
			for _, focus := range []int{0, 1} {
				if event == 16415 && focus == 1 {
					continue
				}
				for _, tc := range cases {
					t.Run(fmt.Sprintf("%d-%d-%d-%q", id, event, focus, tc.text), func(t *testing.T) {
						o := newServerOptionsOwner(t)
						o.installSubpanels(t, true)
						raw := legacy.PortTestServerPanelsConstruct("access", o.options, unsafe.Pointer(&o.settings[0]))
						if raw == 0 {
							t.Fatal("access constructor")
						}
						w := (*gui.Window)(unsafe.Pointer(uintptr(uint32(raw))))
						admission := unsafe.Slice(memmap.PtrUint8(0x5D4594, 371616), 10)
						seed := []byte{0x5a, 0xa5, 0x16, 0x27, 0x38, 0x49, 0x5a, 0x6b, 0x7c, 0x8d}
						copy(admission, seed)
						*memmap.PtrUint32(0x5D4594, 3464) = 19
						*o.optionWords["settings-updated"] = 0
						child := w.ChildByID(uint(id))
						p := alloc.InternCString16(tc.text)
						child.Func94(gui.AsWindowEvent(16414, uintptr(unsafe.Pointer(p)), 0))
						arg1, arg2 := uintptr(unsafe.Pointer(child)), uintptr(0)
						if event == 16387 {
							arg1, arg2 = uintptr(focus), uintptr(id)
						}
						if gui.EventRespInt(w.Func94(gui.AsWindowEvent(event, arg1, arg2))) != 0 {
							t.Fatal("edit return")
						}
						want := append([]byte(nil), seed...)
						shown := tc.text
						limit, updated := uint32(19), uint32(0)
						if focus == 0 {
							value := tc.value
							switch id {
							case 10126, 10128:
								if value > 14 {
									value = 14
									shown = "14"
								}
								if id == 10126 {
									want[1] = seed[1]&0xf0 | byte(value)
								} else {
									want[1] = seed[1]&15 | byte(value<<4)
								}
							case 10130:
								binary.LittleEndian.PutUint16(want[5:], uint16(value))
							case 10132:
								binary.LittleEndian.PutUint16(want[7:], uint16(value))
							case 10133:
								if value < 1 {
									value = 1
									shown = "1"
								}
								if value > 32 {
									value = 32
									shown = "32"
								}
								want[4] = byte(value)
								limit = uint32(value)
								if limit != 19 {
									updated = 1
								}
							}
						}
						ptr := gui.EventRespInt(child.Func94(gui.AsWindowEvent(16413, 0, 0)))
						gotText := alloc.GoString16((*uint16)(unsafe.Pointer(uintptr(uint32(ptr)))))
						if !bytes.Equal(admission, want) || gotText != shown || *memmap.PtrUint32(0x5D4594, 3464) != limit || *o.optionWords["settings-updated"] != updated {
							t.Fatalf("edit admission %x want %x shown %q want %q", admission, want, gotText, shown)
						}
						rows = append(rows, row{id, event, focus, tc.text, gotText, append([]byte(nil), admission...), limit, updated})
					})
				}
			}
		}
	}
	spellbookCapture(t, "server-panels-access-numbers", rows, "b1171d23fd4405da2e4be0e56378c1be0ec49c602db5cfaed639832a243c87f1")
}
