//go:build porttest

package opennox

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/opennox/libs/spell"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestCreatureXferBuffApplication(t *testing.T) {
	s := newCreatureXferOwner(t)
	configure, restore := s.PortTestAISpellDefs()
	t.Cleanup(restore)
	id := server.ENCHANT_INFRAVISION.Spell()
	configure([]server.PortTestSpellClassDef{{Index: uint32(id), Valid: true}})
	s.Spells.DefByInd(id).Effect = spell.SPELL_INFRAVISION
	_, restoreBalance := s.PortTestEffectsUseEnvironment(map[string][]float64{"InfravisionEnchantDuration": {90}}, false, false, 0)
	t.Cleanup(restoreBalance)
	for _, off := range []uintptr{1569740, 1569744} {
		p := memmap.PtrUint32(0x5d4594, off)
		old := *p
		*p = 0
		t.Cleanup(func() { *p = old })
	}
	path := filepath.Join(t.TempDir(), "buff-apply.bin")
	type row struct {
		Case   string
		Buffs  uint32
		Timers [32]uint16
		Powers [32]uint8
		CRC    uint32
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "creature-xfer-buff-application", rows, "eeb10d0e71494c8986fdb7dc42724f964bafe728770ac7ab6f81fa242275e32c")
	}()
	for _, version := range []uint16{1, 2} {
		for _, power := range []byte{0, 1, 5, 6, 255} {
			for _, timer := range []uint32{0, 1, 0x7fff, 0x8000, 0xffff, 0x10000, 0xffffffff} {
				t.Run(fmt.Sprintf("v%d-power%d-timer%d", version, power, timer), func(t *testing.T) {
					u := newCreatureXferObject(t, s, "Monster")
					p := new(mapDrawableStream)
					p.u16(version)
					p.u8(1)
					itemXferRewardName(p, server.ENCHANT_INFRAVISION.String())
					p.u8(power)
					p.u32(timer)
					if err := os.WriteFile(path, p.Bytes(), 0600); err != nil {
						t.Fatal(err)
					}
					if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
						t.Fatal(err)
					}
					defer cryptfile.Close()
					ret := legacy.PortTestCreatureXferHelper(2, u, nil, 0)
					pos, err := cryptfile.Global().File.Seek(0, 1)
					wantPower := power
					if wantPower > 5 {
						wantPower = 5
					}
					var timers [32]uint16
					timers[server.ENCHANT_INFRAVISION] = uint16(timer)
					var powers [32]uint8
					powers[server.ENCHANT_INFRAVISION] = wantPower
					if err != nil || pos != int64(p.Len()) || ret != 1 || u.Buffs != 1<<server.ENCHANT_INFRAVISION || u.BuffsDur != timers || u.BuffsPower != powers {
						t.Fatalf("ret=%d pos=%d buffs=%x timers=%v powers=%v", ret, pos, u.Buffs, u.BuffsDur, u.BuffsPower)
					}
					rows = append(rows, row{t.Name(), u.Buffs, u.BuffsDur, u.BuffsPower, cryptfile.Global().PortTestChecksum()})
				})
			}
		}
	}
}
