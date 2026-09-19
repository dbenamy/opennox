//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestPlayerFilesAbilityRestore(t *testing.T) {
	o := newMatchRosterOwner(t)
	t.Cleanup(o.s.PortTestPlayerFileAbilities())
	oldAbilities := noxServer.abilities
	noxServer.abilities.Init(noxServer)
	t.Cleanup(func() { noxServer.abilities = oldAbilities })
	noxServer.abilities.defs[server.AbilityTreadLightly] = AbilityDef{duration: 77}
	oldDur := noxServer.spells.duration
	noxServer.spells.duration.Init(noxServer)
	t.Cleanup(func() { noxServer.spells.duration = oldDur })
	o.s.Spells.Dur.Init()
	t.Cleanup(o.s.Spells.Dur.Free)
	for _, off := range []uintptr{1569740, 1569744} {
		p := memmap.PtrUint32(0x5D4594, off)
		old := *p
		*p = 0
		t.Cleanup(func() { *p = old })
	}
	u := &o.units[0]
	p := u.UpdateDataPlayer().Player
	p.NetCodeVal = u.NetCode
	*(*byte)(unsafe.Add(unsafe.Pointer(p), 2251)) = 0
	p.SpellLvl[server.AbilityTreadLightly] = 5
	var rows []map[string]any
	for _, version := range []uint16{4, 5} {
		for _, present := range []byte{0, 1} {
			for _, first := range []byte{0, 1, 2} {
				for _, fourth := range []byte{0, 1, 2} {
					o.reset()
					flags.ResetGame()
					if present != 0 {
						flags.SetGame(2048)
					}
					u.ObjFlags = 0
					u.Buffs = 0
					u.BuffsDur = [32]uint16{}
					u.BuffsPower = [32]uint8{}
					ad := o.s.Abils.GetFor(u)
					ad.Cooldowns = [server.AbilityMax]int{}
					ad.Cooldowns[1] = 777
					ad.ExecList = nil
					noxServer.abilities.curxxx = 0
					input := binary.LittleEndian.AppendUint16(nil, version)
					input = append(input, present)
					if present != 0 {
						input = append(input, 0)
					}
					hasTail := version == 4 || present != 0
					if hasTail {
						input = append(input, first, fourth)
						input = binary.LittleEndian.AppendUint32(input, 0xffffffef)
						start := 1
						if first == 1 {
							start = 2
						}
						for i := start; i < 6; i++ {
							input = binary.LittleEndian.AppendUint32(input, 0)
						}
					}
					size := len(input)
					input = append(input, 0xde, 0xad, 0xbe, 0xef)
					ret, got, pos := playerFileSection(t, "nox_xxx_guiEnchantment_41B9C0", input, uint32(uintptr(unsafe.Pointer(u))), 0)
					wantCD := 0
					if !hasTail || first == 1 {
						wantCD = 777
					}
					wantSelected := server.Ability(0)
					if hasTail && first == 1 {
						wantSelected = 1
					}
					active := hasTail && fourth == 1
					if ret != 1 || pos != int64(size) || !bytes.Equal(got, input) || ad.Cooldowns[1] != wantCD || noxServer.abilities.curxxx != wantSelected || (ad.ExecList != nil) != active {
						t.Fatal("ability restoration", version, present, first, fourth, ret, pos, size, ad.Cooldowns, noxServer.abilities.curxxx)
					}
					state := o.state()
					wantNodes := 0
					if active {
						wantNodes = 1
					}
					if len(state.Nodes) != wantNodes {
						t.Fatal("ability messages", state)
					}
					if active {
						entry := ad.ExecList
						if entry.Abil != 4 || entry.Frame != 106 || entry.Active != 1 || entry.Next != nil || entry.Prev != nil || u.Buffs != 1<<server.ENCHANT_SNEAK || u.BuffsDur[server.ENCHANT_SNEAK] != 77 || u.BuffsPower[server.ENCHANT_SNEAK] != 5 || !bytes.Equal(state.Nodes[0].Data, []byte{byte(netmsg.MSG_REPORT_ACTIVE_ABILITIES), 4, 1}) {
							t.Fatal("active tread lightly", entry, u.Buffs, u.BuffsDur, u.BuffsPower, state)
						}
					} else if u.Buffs != 0 {
						t.Fatal("unexpected ability effect")
					}
					rows = append(rows, map[string]any{"version": version, "present": present, "first": first, "fourth": fourth, "return": ret, "position": pos, "cooldowns": ad.Cooldowns, "selected": noxServer.abilities.curxxx, "active": active, "buffs": u.Buffs, "queue": state})
				}
			}
		}
	}
	spellbookCapture(t, "player-files-ability-restore", rows, "5864836439dfbc2bc0cd5328ad9d8fac8f9ad01191908a7a1045299962097273")
}
