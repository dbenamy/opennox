//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestScriptInventoryStartup(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	orders := [][]int{nil, {0}, {0, 1, 2}, {2, 1, 0}, {1, 2, 0}}
	for mask := 0; mask < 8; mask++ {
		for _, order := range orders {
			for _, chapter := range []uint32{0, 1, 0xffffffff} {
				s := controlsBase(58)
				p := s.Callbacks.Shop
				p.TemporaryUpdates.World.Objectives.Players = 3
				for i := range p.Items {
					p.Items[i].Class = 8
					if mask&(1<<i) != 0 {
						p.Items[i].Class |= 0x40
					}
					p.Items[i].Flags = 0
				}
				p.TemporaryUpdates.World.Objectives.Attack.Controls.ScriptStartup = &legacy.PortTestScriptStartupSpec{Order: order, Chapter: chapter, Journal: []uint16{0, 1, 2, 4, 8, 14, 16, 0xffff}}
				cases = append(cases, s)
			}
		}
	}
	callbackHash(t, "script-inventory-startup", controlsRun(t, cases), "85dea812ebfb2d36ed1b2557432eb58c7b4aca6ea5a5dafd9a665342b688f369")
}

func TestScriptInventoryCarrySelection(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for variant := 0; variant < 12; variant++ {
		for _, free := range []int{0, 1, 80} {
			for _, reserve := range []int32{0, 1, 80, 81} {
				for _, notice := range []uint32{0, 1, 9} {
					for _, success := range []bool{false, true} {
						s := controlsBase(59)
						p := s.Callbacks.Shop
						p.Inventory.Linked = []int{0, 1, 2}
						p.Inventory.DropResult = success
						for i := range p.Items {
							p.Items[i] = legacy.PortTestShopItem{Type: uint16(23 + i), Class: 8, Worth: uint32(10 + 10*i)}
						}
						want := 0
						switch variant {
						case 1:
							p.Inventory.Linked = nil
							want = -1
						case 2:
							p.Items[0].Class |= 0x10
							want = 1
						case 3:
							p.Items[0].Flags |= 0x100
							want = 1
						case 4:
							p.Inventory.ItemTypes = []string{"Glyph"}
							want = 1
						case 5:
							p.Inventory.DropTable = [][3]uint32{{1, 23, 0}}
							want = 1
						case 6:
							p.Items[0].Worth = 999999
							p.Items[1].Worth = 1000000
							p.Items[2].Worth = 999998
							want = 2
						case 7:
							for i := range p.Items {
								p.Items[i].Worth = 999999
							}
							want = -1
						case 8:
							for i := range p.Items {
								p.Items[i].Worth = 10
							}
							p.Inventory.Linked = []int{2, 1, 0}
							want = 2
						case 9:
							p.Items[1].Worth = 1
							want = 1
						case 10:
							p.Items[0].Class |= 0x10
							p.Items[1].Flags |= 0x100
							p.Inventory.ItemTypes = []string{"", "", "Glyph"}
							want = -1
						case 11:
							p.Inventory.Linked = []int{1, 0}
							p.Items[0].Worth = 20
							want = 1
						}
						if free > int(reserve) {
							want = -1
						}
						p.TemporaryUpdates.World.Objectives.Attack.Controls.ScriptCarry = &legacy.PortTestScriptCarrySpec{FreeCells: free, Count: 31, SameType: true, Reserved: reserve, Notice: notice, Capacity: free, WantItem: want, Repeats: 2}
						cases = append(cases, s)
					}
				}
			}
		}
	}
	callbackHash(t, "script-inventory-carry-selection", controlsRun(t, cases), "ee653fca8c51b397b2729f32e25af94bed91b24b0b8d400f94c35af9311e6bf2")
}

func TestScriptInventoryCarryStacks(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, flags := range []uint32{1, 2049, 4097} {
		for _, class := range []uint32{0, 0x10, 0x4000000, 0x4000010} {
			for _, count := range []byte{2, 3, 8, 9, 30, 31} {
				for _, same := range []bool{false, true} {
					for _, reserve := range []int32{0, 79, 80} {
						s := controlsBase(59)
						s.Lifecycle.GameFlags = flags
						p := s.Callbacks.Shop
						p.Inventory.Linked = []int{0}
						p.Items[0] = legacy.PortTestShopItem{Type: 23, Class: 8, Worth: 20}
						limit := 31
						if class&0x10 != 0 {
							limit = 3
							if flags&6144 != 0 {
								limit = 9
							}
						}
						capacity := 0
						if same && class&0x4000000 == 0 && int(count)+1 <= limit {
							capacity = 80
						}
						want := 0
						if capacity > int(reserve) {
							want = -1
						}
						p.TemporaryUpdates.World.Objectives.Attack.Controls.ScriptCarry = &legacy.PortTestScriptCarrySpec{Count: count, Class: class, SameType: same, Reserved: reserve, Capacity: capacity, WantItem: want, Repeats: 1}
						cases = append(cases, s)
					}
				}
			}
		}
	}
	callbackHash(t, "script-inventory-carry-stacks", controlsRun(t, cases), "a59c3196c203e0adb5d5ce5db568c55734bf322b2bce74b0e0cdcae87e64c62b")
}

func TestScriptInventoryHalberd(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for index := uint32(0); index < 4; index++ {
		for variant := 0; variant < 14; variant++ {
			for _, equipped := range []bool{false, true} {
				for _, allowed := range []bool{false, true} {
					s := controlsBase(60)
					s.Lifecycle.Place = true
					s.Lifecycle.Eligible = allowed
					p := s.Callbacks.Shop
					p.TemporaryUpdates.World.Objectives.Players = 3
					for i := range p.Items {
						p.Items[i] = legacy.PortTestShopItem{Type: uint16(23 + i), Class: 8}
					}
					order := []int{0, 1, 2}
					want := 0
					p.Items[0].Class = 0x1000000
					p.Items[0].Subclass = 0x800000
					switch variant {
					case 0:
						order = nil
						want = -1
					case 1:
						p.Items[0].Subclass = 0
						want = -1
					case 2:
						p.Items[0].Class = 8
						p.Items[1].Class = 0x1000000
						p.Items[1].Subclass = 0x1000000
						want = 1
					case 3:
						p.Items[0].Class = 8
						p.Items[2].Class = 0x1000000
						p.Items[2].Subclass = 0x2000000
						want = 2
					case 4:
						p.Items[1].Class = 0x1000000
						p.Items[1].Subclass = 0x4000000
					case 5:
						p.Items[1].Class = 0x1000000
						p.Items[1].Subclass = 0x4000000
						order = []int{2, 1, 0}
						want = 1
					case 6:
						p.Items[0].Class = 0x2000000
						want = -1
					case 7:
						p.Items[0].Subclass = 0x400000
						want = -1
					case 8:
						p.Items[0].Subclass = 0x8000000
						want = -1
					case 9:
						p.Items[0].Subclass = 0x1000000
					case 10:
						p.Items[0].Subclass = 0x2000000
					case 11:
						p.Items[0].Subclass = 0x4000000
					case 12:
						p.Items[0].Subclass = 0x7800000
					case 13:
						p.Items[0].Subclass = 0x87800000
					}
					if equipped && want >= 0 {
						p.Items[want].Flags = 0x100
					}
					p.TemporaryUpdates.World.Objectives.Attack.Controls.ScriptHalberd = &legacy.PortTestScriptHalberdSpec{Index: index, Order: order, WantDelete: want, Equipped: equipped && want >= 0 && allowed}
					cases = append(cases, s)
				}
			}
		}
	}
	callbackHash(t, "script-inventory-halberd", controlsRun(t, cases), "c33f3ac928be5f759e425de97f1a4a024882befdc26c83bceed4f299047b7a3c")
}
