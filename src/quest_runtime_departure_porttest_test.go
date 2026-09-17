//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
)

func TestQuestRuntimeDepartureTracking(t *testing.T) {
	o := newQuestRuntimeOwner(t)
	stamps := unsafe.Slice(memmap.PtrUint32(0x5D4594, 1556172), 32)
	type state struct {
		Words [32]uint32
		Bytes [96]byte
	}
	snapshot := func() [3]state {
		var out [3]state
		for i := range o.units {
			d := o.units[i].UpdateData
			for j := 0; j < 32; j++ {
				out[i].Words[j] = objectXferGetWord(d, 324+4*j)
			}
			copy(out[i].Bytes[:], unsafe.Slice((*byte)(unsafe.Add(d, 452)), 96))
		}
		return out
	}
	type row struct {
		Name   string
		Return uint64
		Stamps [32]uint32
		State  [3]state
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "quest-runtime-departure", rows, "e184d7fc205b77f1c6402dc8303bb2ed691faa8af102a543d4d0ec0d0c15aac5")
	}()
	for slot := 0; slot < 32; slot++ {
		for _, fps := range []int{0, 1, 30, 60} {
			for _, age := range []uint32{0, uint32(30 * fps), uint32(30*fps) + 1, 0xffffffff} {
				for _, participation := range []uint32{0, 1, 2} {
					name := fmt.Sprintf("slot%d/fps%d/age%d/participation%d", slot, fps, age, participation)
					t.Run(name, func(t *testing.T) {
						clear(stamps)
						stamps[slot] = 100
						o.s.SetTickRate(uint32(fps))
						o.s.SetFrame(100 + age)
						present := false
						for i := range o.units {
							d := o.units[i].UpdateData
							pl := o.units[i].UpdateDataPlayer().Player
							objectXferSetWord(pl.C(), 4792, participation)
							if int(pl.PlayerInd) == slot && participation == 1 {
								present = true
							}
							for j := 0; j < 32; j++ {
								objectXferSetWord(d, 324+4*j, 0xaabb0000+uint32(j))
							}
							for j := 0; j < 96; j++ {
								*(*byte)(unsafe.Add(d, 452+j)) = byte(1 + j)
							}
						}
						want := snapshot()
						stampWant := uint32(100)
						if present {
							stampWant = 0
						} else if age > uint32(30*fps) {
							stampWant = 0
							for i := range want {
								want[i].Words[slot] = 0
								for k := 0; k < 3; k++ {
									want[i].Bytes[slot+32*k] = 0
								}
							}
						}
						rv := questRuntimeCall("sub_4D7A80", nil)
						got := snapshot()
						if rv != 32 || stamps[slot] != stampWant || got != want {
							t.Fatalf("tracking return%d stamp%d want%d state match%t", rv, stamps[slot], stampWant, got == want)
						}
						var st [32]uint32
						copy(st[:], stamps)
						for i, v := range st {
							if i != slot && v != 0 {
								t.Fatal("unrelated slot changed")
							}
						}
						rows = append(rows, row{name, rv, st, got})
					})
				}
			}
		}
	}
	// A zero timestamp is disabled, even with maximum elapsed frame age.
	clear(stamps)
	o.s.SetFrame(0xffffffff)
	before := snapshot()
	questRuntimeCall("sub_4D7A80", nil)
	if snapshot() != before {
		t.Fatal("disabled timestamps changed relationships")
	}
	for _, frame := range []uint32{0, 1, 0x80000000, 0xffffffff} {
		o.s.SetFrame(frame)
		for slot := 0; slot < 32; slot++ {
			if rv := questRuntimeCall("sub_4D7A60", nil, uint32(slot)); rv != uint64(slot) || stamps[slot] != frame {
				t.Fatal("departure stamp", slot, rv, stamps[slot])
			}
		}
		if rv := questRuntimeCall("sub_4D7B40", nil); rv != 0 {
			t.Fatal("reset result")
		}
		for _, v := range stamps {
			if v != 0 {
				t.Fatal("reset retained timestamp")
			}
		}
	}
}
