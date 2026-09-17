//go:build porttest

package opennox

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestItemXferEquipmentHealth(t *testing.T) {
	s := newItemXferOwner(t)
	path := filepath.Join(t.TempDir(), "health.bin")
	var rows []itemXferCaptureRow
	defer func() {
		spellbookCapture(t, "item-xfer-equipment-health", rows, "af7f8e014742cec75751206eb64e9cfa484c9184927f8a436db1b7831116ce56")
	}()
	modes := []struct{ solo, switching, quest, players bool }{
		{solo: true}, {switching: true}, {}, {quest: true}, {quest: true, players: true}, {players: true},
	}
	for _, name := range []string{"Weapon", "Armor"} {
		for _, v := range []int16{40, 41, 42, 61, 62, 64} {
			if name == "Armor" && v > 62 {
				continue
			}
			for mi, mode := range modes {
				for _, definition := range []bool{false, true} {
					for _, hp := range []uint16{0, 37, 100, 101, 32768, 65535} {
						t.Run(fmt.Sprintf("%s-v%d-mode%d-definition%v-hp%d", name, v, mi, definition, hp), func(t *testing.T) {
							oldSolo, oldSwitch := gameIsNotMultiplayer, gameIsSwitchToSolo
							gameIsNotMultiplayer, gameIsSwitchToSolo = mode.solo, mode.switching
							t.Cleanup(func() { gameIsNotMultiplayer, gameIsSwitchToSolo = oldSolo, oldSwitch })
							flags := noxflags.GameFlag(0x200000)
							if mode.quest {
								flags |= 0x1000
							}
							t.Cleanup(noxflags.PortTestGameFlags(flags))
							mask := uint32(0)
							if mode.players {
								mask = 1
							}
							t.Cleanup(s.PortTestItemXferQuestPlayers(mask))
							u := newItemXferObject(t, s, name)
							def, free := alloc.New(server.Modifier{})
							t.Cleanup(free)
							def.TypeInd = uint32(u.TypeInd)
							def.Durability52 = 0x12340051
							oldWeapon, oldArmor := s.Modif.Dword_5d4594_251600, s.Modif.Dword_5d4594_251608
							t.Cleanup(func() { s.Modif.Dword_5d4594_251600, s.Modif.Dword_5d4594_251608 = oldWeapon, oldArmor })
							s.Modif.Dword_5d4594_251600, s.Modif.Dword_5d4594_251608 = nil, nil
							if definition {
								if name == "Weapon" {
									s.Modif.Dword_5d4594_251600 = def
								} else {
									s.Modif.Dword_5d4594_251608 = def
								}
							}
							var p mapDrawableStream
							p.u16(uint16(v))
							objectXferEmptyBase(&p, v)
							p.Write(make([]byte, 4))
							hasHealth := name == "Armor" && v >= 41 || name == "Weapon" && v >= 42
							if hasHealth {
								p.u16(hp)
							}
							if name == "Armor" && v == 61 || name == "Weapon" && v == 63 {
								p.u8(7)
							}
							if name == "Armor" && v >= 62 || name == "Weapon" && v >= 64 {
								p.u32(123)
							}
							if err := os.WriteFile(path, p.Bytes(), 0600); err != nil {
								t.Fatal(err)
							}
							if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
								t.Fatal(err)
							}
							defer cryptfile.Close()
							if err := u.CallXfer(nil); err != nil {
								t.Fatal(err)
							}
							cur, max, initial := uint16(80), uint16(100), uint16(100)
							if hasHealth {
								if mode.solo || mode.switching || mode.quest && mode.players {
									cur = hp
									if cur > 100 {
										cur = 100
									}
								} else if definition {
									cur, max, initial = 81, 81, 81
								}
							}
							if u.HealthData.Cur != cur || u.HealthData.Max != max || u.HealthData.Field2 != initial {
								t.Fatalf("health=%d/%d/%d want=%d/%d/%d", u.HealthData.Cur, u.HealthData.Max, u.HealthData.Field2, cur, max, initial)
							}
							pos, err := cryptfile.Global().File.Seek(0, 1)
							if err != nil || pos != int64(p.Len()) {
								t.Fatalf("position=%d want=%d err=%v", pos, p.Len(), err)
							}
							itemXferCaptureCase(t, &rows, u)
						})
					}
				}
			}
		}
	}
}
