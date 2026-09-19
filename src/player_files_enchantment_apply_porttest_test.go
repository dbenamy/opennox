//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/spell"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestPlayerFilesEnchantmentApply(t *testing.T) {
	o := newReliableReportsOwner(t)
	defer flags.PortTestGameFlags(2048)()
	oldDur := noxServer.spells.duration
	noxServer.spells.duration.Init(noxServer)
	t.Cleanup(func() { noxServer.spells.duration = oldDur })
	o.s.Spells.Dur.Init()
	t.Cleanup(o.s.Spells.Dur.Free)
	configure, restore := o.s.PortTestAISpellDefs()
	t.Cleanup(restore)
	id := server.ENCHANT_INFRAVISION.Spell()
	configure([]server.PortTestSpellClassDef{{Index: uint32(id), Valid: true}})
	o.s.Spells.DefByInd(id).Effect = spell.SPELL_INFRAVISION
	_, restoreBalance := o.s.PortTestEffectsUseEnvironment(map[string][]float64{"InfravisionEnchantDuration": {90}}, false, false, 0)
	t.Cleanup(restoreBalance)
	for _, off := range []uintptr{1569740, 1569744} {
		p := memmap.PtrUint32(0x5D4594, off)
		old := *p
		*p = 0
		t.Cleanup(func() { *p = old })
	}
	u := &o.units[0]
	p := u.UpdateDataPlayer().Player
	*(*byte)(unsafe.Add(unsafe.Pointer(p), 2251)) = 1
	var rows []map[string]any
	for _, version := range []uint16{0, 1, 2, 3, 4, 5, 0xffff} {
		for _, power := range []byte{0, 1, 5, 6, 255} {
			for _, timer := range []uint16{0, 1, 32768, 65535} {
				o.reset()
				u.ObjFlags = 0
				u.Buffs = 0
				u.BuffsDur = [32]uint16{}
				u.BuffsPower = [32]uint8{}
				input := binary.LittleEndian.AppendUint16(nil, version)
				input = append(input, 1, 1)
				name := server.ENCHANT_INFRAVISION.String()
				input = append(input, byte(len(name)))
				input = append(input, name...)
				input = binary.LittleEndian.AppendUint16(input, timer)
				if int16(version) >= 2 {
					input = append(input, power)
				}
				size := len(input)
				input = append(input, 0xde, 0xad, 0xbe, 0xef)
				wantPower := power
				if int16(version) < 2 {
					wantPower = 2
				}
				if wantPower > 5 {
					wantPower = 5
				}
				wantTimer := timer
				if timer == 0 {
					wantTimer = 30
				}
				var timers [32]uint16
				timers[server.ENCHANT_INFRAVISION] = wantTimer
				var powers [32]byte
				powers[server.ENCHANT_INFRAVISION] = wantPower
				ret, got, pos := playerFileSection(t, "nox_xxx_guiEnchantment_41B9C0", input, uint32(uintptr(unsafe.Pointer(u))), 0)
				if ret != 1 || pos != int64(size) || !bytes.Equal(got, input) || u.Buffs != 1<<server.ENCHANT_INFRAVISION || u.BuffsDur != timers || u.BuffsPower != powers {
					t.Fatal("enchant apply", version, power, timer, ret, pos, size, u.Buffs, u.BuffsDur, u.BuffsPower)
				}
				rows = append(rows, map[string]any{"version": version, "power": power, "timer": timer, "return": ret, "position": pos, "buffs": u.Buffs, "timers": u.BuffsDur, "powers": u.BuffsPower, "queue": o.state()})
			}
		}
	}
	spellbookCapture(t, "player-files-enchantment-apply", rows, "")
}
