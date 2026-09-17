//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestCreatureXferPostloadReferences(t *testing.T) {
	s := newCreatureXferOwner(t)
	var rows []struct {
		Case    string
		Entries [][6]uint32
	}
	defer func() { spellbookCapture(t, "creature-xfer-postload", rows, "") }()
	run := func(label string, ids []int, stack byte, present, destroyed bool) {
		t.Run(label, func(t *testing.T) {
			u := newCreatureXferObject(t, s, "Monster")
			ref := newItemXferObject(t, s, "Gold")
			ref.ScriptIDVal = 12345
			if destroyed {
				ref.ObjFlags |= object.FlagDestroyed
			}
			oldList := s.Objs.List
			s.Objs.List = nil
			if present {
				s.Objs.List = ref
			}
			t.Cleanup(func() { s.Objs.List = oldList })
			wp, free := alloc.New(server.Waypoint{})
			t.Cleanup(free)
			wp.Index = 9876
			oldWP := s.WPs.List
			s.WPs.List = nil
			if present {
				s.WPs.List = wp
			}
			t.Cleanup(func() { s.WPs.List = oldWP })
			entries := unsafe.Slice((*[6]uint32)(unsafe.Add(u.UpdateData, 552)), len(ids))
			expected := make([][6]uint32, len(ids))
			normalized := make([][6]uint32, len(ids))
			for i, id := range ids {
				entries[i] = [6]uint32{uint32(id), 0xdead0001, 0xdead0002, 0xdead0003, 0xdead0004, 0xdead0005}
				n := int(*memmap.PtrUint32(0x587000, uintptr(255604+16*id)))
				for j := 0; j < n; j++ {
					kind := *memmap.PtrUint32(0x587000, uintptr(255608+16*id+4*j))
					if kind == 1 {
						entries[i][1+2*j] = uint32(ref.ScriptIDVal)
					} else if kind == 2 {
						entries[i][1+2*j] = wp.Index
					}
				}
				expected[i] = entries[i]
				normalized[i] = entries[i]
				if stack&0x80 != 0 {
					continue
				}
				for j := 0; j < n; j++ {
					kind := *memmap.PtrUint32(0x587000, uintptr(255608+16*id+4*j))
					word := 1 + 2*j
					if kind == 1 || kind == 2 {
						expected[i][word] = 0
						normalized[i][word] = 0
						if present && (kind == 2 || !destroyed) {
							if kind == 1 {
								expected[i][word] = uint32(uintptr(ref.CObj()))
							} else {
								expected[i][word] = uint32(uintptr(unsafe.Pointer(wp)))
							}
							normalized[i][word] = 1
						}
					}
				}
			}
			*(*byte)(unsafe.Add(u.UpdateData, 544)) = stack
			ret := legacy.PortTestCreatureXferHelper(6, u, nil, 0)
			if ret != uint32(uintptr(u.UpdateData)) {
				t.Fatal("postload return owner")
			}
			for i := range entries {
				if entries[i] != expected[i] {
					t.Fatalf("entry%d reference identity/untouched fields", i)
				}
			}
			rows = append(rows, struct {
				Case    string
				Entries [][6]uint32
			}{t.Name(), normalized})
		})
	}
	for id := 0; id < 72; id++ {
		for _, mode := range []int{0, 1, 2} {
			run(fmt.Sprintf("id%d-mode%d", id, mode), []int{id}, 0, mode != 0, mode == 2)
		}
	}
	for _, count := range []int{2, 24} {
		ids := make([]int, count)
		for i := range ids {
			ids[i] = (i*7 + 3) % 72
		}
		run(fmt.Sprintf("stack%d", count), ids, byte(count-1), true, false)
	}
	for _, stack := range []byte{128, 255} {
		run(fmt.Sprintf("disabled%d", stack), []int{3, 10}, stack, true, false)
	}
}
