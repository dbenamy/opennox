//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func serverTextExtensionBase() legacy.PortTestRoamSpec {
	s := controlsBase(97)
	p := s.Callbacks.Shop
	p.Equipment.GameEx = 0
	p.Equipment.ActiveAbilities = 0
	p.Equipment.Definitions = []legacy.PortTestEquipmentDef{{Type: 23, Strength: 20}, {Type: 24, Strength: 20}, {Type: 25, Strength: 20}}
	p.Inventory.WeaponBits = map[uint16]uint32{23: 16, 24: 32, 25: 64}
	for i := range p.Items {
		p.Items[i].Type = uint16(23 + i)
		p.Items[i].Class = 0x1000000
		p.Items[i].Subclass = 16 << i
		p.Items[i].Flags = 0
		p.Items[i].Mods = [4]bool{}
	}
	return s
}

func TestServerTextWeaponRoll(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, status := range []byte{0, 1, 2, 3, 4, 128, 255} {
		for _, state := range []byte{0, 1, 2, 255} {
			for _, direction := range []int8{0, 1, -1, -128, 127} {
				for current := -1; current < 3; current++ {
					for mode := 0; mode < 7; mode++ {
						s := serverTextExtensionBase()
						p := s.Callbacks.Shop
						sp := &legacy.PortTestExtensionSpec{Status: status, State: state, Direction: direction, Current: current, Mode: mode, Order: []int{0, 1, 2}, Want: -1}
						p.TemporaryUpdates.World.Objectives.Attack.Controls.Extension = sp
						switch mode {
						case 1:
							p.Inventory.WeaponBits[24] = 2
						case 3:
							p.Equipment.Definitions[1].Strength = 31
						case 6:
							p.Inventory.WeaponBits[24] = 0
						}
						// Pick from the explicit fixture order, independently of object links.
						candidates := []int{0, 1, 2}
						if current >= 0 {
							if direction == 0 {
								candidates = nil
								for i := current - 1; i >= 0; i-- {
									candidates = append(candidates, i)
								}
							} else {
								candidates = candidates[current+1:]
							}
						}
						if status&3 == 0 && state != 1 {
							for _, i := range candidates {
								if i == 1 && (mode == 1 || mode == 2 || mode == 3 || mode == 6) {
									continue
								}
								sp.Want = i
								break
							}
						}
						specs = append(specs, s)
					}
				}
			}
		}
	}
	results := controlsRun(t, specs)
	serverTextCaptureStates(t, "server-text-weapon-roll", results)
	t.Logf("%d weapon-roll cases", len(specs))
}

func TestServerTextTrapDrop(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, status := range []byte{0, 1, 2, 3, 4, 128, 255} {
		for _, state := range []byte{0, 1, 2, 255} {
			for _, nilUnit := range []bool{false, true} {
				for _, success := range []bool{false, true} {
					for mask := 0; mask < 8; mask++ {
						for _, order := range [][]int{nil, {0, 1, 2}, {2, 1, 0}} {
							s := serverTextExtensionBase()
							p := s.Callbacks.Shop
							p.Inventory.DropResult = success
							sp := &legacy.PortTestExtensionSpec{Drop: true, Mode: mask % 2, Nil: nilUnit, Status: status, State: state, Current: -1, Order: order, Want: -1}
							p.TemporaryUpdates.World.Objectives.Attack.Controls.Extension = sp
							for i := range p.Items {
								p.Items[i].Class = 0x1000000
								if mask&(1<<i) != 0 {
									p.Items[i].Class = 0x1110000
								}
							}
							if !nilUnit && status&3 == 0 && state != 1 {
								for _, i := range order {
									if mask&(1<<i) != 0 {
										sp.Want = i
										break
									}
								}
							}
							specs = append(specs, s)
						}
					}
				}
			}
		}
	}
	results := controlsRun(t, specs)
	serverTextCaptureStates(t, "server-text-trap-drop", results)
	t.Logf("%d trap-drop cases", len(specs))
}

// Capture every normalized state field without duplicating hundreds of MB of
// mostly unchanged fixture memory. Independent contracts run before hashing.
func serverTextCaptureStates(t *testing.T, name string, results []legacy.PortTestRoamResult) {
	t.Helper()
	hashes := make([]string, len(results))
	for i, r := range results {
		b, err := json.Marshal(r)
		if err != nil {
			t.Fatal(err)
		}
		hashes[i] = fmt.Sprintf("%x", sha256.Sum256(b))
	}
	interactionCapture(t, name, hashes)
}
