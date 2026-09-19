//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"image"
	"reflect"
	"testing"
	"unsafe"
)

func TestClientPresentationChant(t *testing.T) {
	o := newObjectRenderOwner(t)
	player := o.drawable(7, image.Pt(48, 48))
	for _, r := range blobdata.PortTestClientPresentationTables() {
		copy(serverConfigOwnBytes(t, r.Base, r.Offset, len(r.Data)), r.Data)
	}
	clear(serverConfigOwnBytes(t, 0x5D4594, 1096596, 32))
	clear(serverConfigOwnBytes(t, 0x5D4594, 1303504, 16))
	treeWord, restore := legacy.PortTestPresentationChantOwner()
	t.Cleanup(restore)
	tree, free := alloc.Make([]server.PhonemeLeaf{}, 10)
	t.Cleanup(free)
	phon := []int{0, 1, 2, 3, 5, 6, 7, 8, 4}
	for i, p := range phon {
		tree[i].Pho[p] = &tree[i+1]
	}
	tree[9].Ind = 1
	t.Cleanup(o.c.srv.Server.PortTestSpellLifecycle([]server.PortTestSpellLifecycleDef{{Index: 1, Valid: true, Enabled: true, Phonemes: phon}}, &tree[0]))
	var sounds [][2]int
	t.Cleanup(legacy.PortTestClientSoundObserver(func(id, volume int) { sounds = append(sounds, [2]int{id, volume}) }))
	type record struct {
		Voice                  int
		Start, Frame, Deadline uint32
		Cancel                 bool
		Index, Active          byte
		Stamps                 [8]uint32
		Sounds                 [][2]int
	}
	var rows []record
	icon := [9]int{0, 1, 2, 3, 1, 4, 5, 6, 7}
	sound := [9]int{193, 186, 187, 192, 0, 188, 191, 190, 189}
	for voice := 0; voice < 3; voice++ {
		player.ObjClass = 0
		if voice != 0 {
			player.ObjClass = 4
		}
		*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 2252)) = byte(voice / 2)
		for _, start := range []uint32{0, 120, 0xfffffffe} {
			for _, cancel := range []bool{false, true} {
				clear(unsafe.Slice(memmap.PtrUint32(0x5D4594, 1096596), 8))
				sounds = nil
				o.c.srv.SetFrame(start)
				legacy.PortTestPresentationChantStart(1)
				active, index, deadline := byte(1), byte(0), start
				var stamps [8]uint32
				var expected [][2]int
				if *treeWord != uint32(uintptr(unsafe.Pointer(&tree[0]))) || memmap.Uint32(0x5D4594, 1303516) != start || memmap.Uint8(0x5D4594, 1303504) != 1 || memmap.Uint8(0x5D4594, 1303512) != 0 {
					t.Fatal("chant start state")
				}
				for step := 0; step < 40; step++ {
					frame := start + uint32(step)
					o.c.srv.SetFrame(frame)
					if cancel && step == 4 {
						legacy.PortTestPresentationChantClear()
						active = 0
					}
					if active != 0 && frame >= deadline {
						p := phon[index]
						id := sound[p]
						if voice == 2 && id != 0 {
							id += 8
						}
						expected = append(expected, [2]int{id, 100})
						stamps[icon[p]] = frame
						index++
						deadline = frame + 3
						if index == 9 {
							active = 0
						}
					}
					legacy.PortTestPresentationChantTick()
					if memmap.Uint8(0x5D4594, 1303504) != active || memmap.Uint8(0x5D4594, 1303512) != index || memmap.Uint32(0x5D4594, 1303516) != deadline || *treeWord != uint32(uintptr(unsafe.Pointer(&tree[index]))) {
						t.Fatalf("chant schedule/tree voice%d frame%d", voice, frame)
					}
					if !reflect.DeepEqual(sounds, expected) || *(*[8]uint32)(memmap.PtrOff(0x5D4594, 1096596)) != stamps {
						t.Fatal("chant sound/icon sequence")
					}
					rows = append(rows, record{voice, start, frame, deadline, cancel, index, active, stamps, append([][2]int(nil), sounds...)})
				}
			}
		}
	}
	spellbookCapture(t, "client-presentation-chant", rows, "b014fbb6c8afb9f4db3ed018029da4ddb19b7c20fb515e69ea2941aac5a1bfcd")
}
