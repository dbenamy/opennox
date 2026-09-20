//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"reflect"
	"testing"
	"unsafe"
)

func TestGameMessageClientSessionGauntletNotices(t *testing.T) {
	c, _, _ := newEffectsFullOwner(t)
	units, configurePlayers, _, free := c.srv.PortTestEscortPlayers()
	t.Cleanup(free)
	configurePlayers(2)
	pl, other := units[0].UpdateDataPlayer().Player, units[1].UpdateDataPlayer().Player
	oldServer, oldPlayer := noxServer, legacy.Get_dword_8531A0_2576()
	noxServer = c.srv
	t.Cleanup(func() { noxServer = oldServer; legacy.Set_dword_8531A0_2576(oldPlayer) })
	printer, index, region := teamRuntimeJoinTextOwner(t, c.srv.Server)
	var entries []strman.Entry
	for _, p := range [][2]string{{"plyrspel.c:AwardSpell", "Spell %s"}, {"plyrspel.c:AwardSpellToOther", "%s gets spell %s"}, {"PlyrGide.c:AwardGuide", "Guide %s"}, {"PlyrGide.c:AwardGuideToOther", "%s gets guide %s"}, {"ComAblty.c:AwardAbility", "Ability %s"}, {"ComAblty.c:AwardAbilityToOther", "%s gets ability %s"}, {"guimsg.c:systemmsg", "%s"}, {"cdecode.c:QuestNotice", "Notice %s"}} {
		entries = append(entries, strman.Entry{ID: strman.ID(p[0]), Vals: []strman.Variant{{Str: p[1]}}})
	}
	for i := 0; i < 5; i++ {
		entries = append(entries, strman.Entry{ID: strman.ID(fmt.Sprintf("cdecode.c:QuestClass%d", i)), Vals: []strman.Variant{{Str: fmt.Sprintf("Class %d", i)}}})
	}
	set, restoreStrings := c.srv.PortTestMeterStrings(entries...)
	t.Cleanup(restoreStrings)
	set(0)
	setDefs, restoreDefs := c.srv.PortTestAISpellDefs()
	t.Cleanup(restoreDefs)
	setDefs([]server.PortTestSpellClassDef{{Index: 1, Valid: true}, {Index: 5, Valid: true}, {Index: 136, Valid: true}})
	for _, id := range []spell.ID{1, 5, 136} {
		c.srv.Spells.DefByInd(id).Title = fmt.Sprintf("Spell title %d Ω", id)
	}
	oldAbilities := c.srv.abilities.defs
	t.Cleanup(func() { c.srv.abilities.defs = oldAbilities })
	for id := 1; id <= 5; id++ {
		c.srv.abilities.defs[id] = AbilityDef{name: fmt.Sprintf("Ability title %d Ω", id), field24: 1}
	}
	guides := serverConfigOwnBytes(t, 0x5D4594, 740076, 41*28)
	clear(guides)
	for _, id := range []int{1, 24, 40} {
		binary.LittleEndian.PutUint32(guides[28*id:], uint32(uintptr(unsafe.Pointer(alloc.InternCString16(fmt.Sprintf("Guide title %d Ω", id))))))
		binary.LittleEndian.PutUint32(guides[28*id+4:], 1)
	}
	classes := serverConfigOwnBytes(t, 0x587000, 160948, 204)
	clear(classes)
	for i, off := range []int{0, 40, 80, 120, 164} {
		copy(classes[off:], fmt.Sprintf("QuestClass%d", i))
	}
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	c.srv.SetFrame(100)
	type row struct {
		On, Present, Local, Kind, ID int
		Code                         uint16
		Messages                     []string
		Centered                     []byte
		Return                       int
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for present := 0; present < 2; present++ {
			for local := 0; local < 2; local++ {
				for _, kind := range []int{30, 31, 32, 33} {
					ids := []int{1, 5, 136}
					if kind == 31 {
						ids = []int{1, 24, 40}
					}
					if kind == 32 {
						ids = []int{1, 3, 5}
					}
					if kind == 33 {
						ids = []int{0, 1, 2, 3, 4}
					}
					for _, id := range ids {
						for _, code := range []uint16{17, 0x8011, 65535} {
							binary.LittleEndian.PutUint32(connected, uint32(on))
							pl.Active = byte(present)
							pl.NetCodeVal = 17
							other.NetCodeVal = 999
							pl.SetName("Player Ω")
							legacy.Set_dword_8531A0_2576(other)
							if local != 0 {
								legacy.Set_dword_8531A0_2576(pl)
							}
							clear(region)
							*index = 0
							printer.lines = nil
							data := []byte{240, byte(kind), byte(id), byte(code), byte(code >> 8)}
							length := 5
							if kind == 33 {
								data = make([]byte, 52)
								data[0], data[1], data[51] = 240, 33, byte(id)
								copy(data[2:], "QuestNotice")
								length = 52
							}
							input := bytes.Clone(data)
							n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(240), data)
							var want []string
							if on != 0 && (kind == 33 || present != 0 && code == 17) {
								message := ""
								switch kind {
								case 33:
									message = fmt.Sprintf("Notice Class %d", id)
								default:
									category := map[int]string{30: "Spell", 31: "Guide", 32: "Ability"}[kind]
									title := fmt.Sprintf("%s title %d Ω", category, id)
									message = fmt.Sprintf("%s %s", category, title)
									if local == 0 {
										lower := map[int]string{30: "spell", 31: "guide", 32: "ability"}[kind]
										message = fmt.Sprintf("Player Ω gets %s %s", lower, title)
									}
								}
								want = []string{message}
							}
							if n != length || !bytes.Equal(input, data) || !reflect.DeepEqual(printer.lines, want) {
								t.Fatal("quest localized notice", on, present, local, kind, id, code, n, printer.lines, want)
							}
							rows = append(rows, row{on, present, local, kind, id, code, append([]string(nil), printer.lines...), bytes.Clone(region), n})
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-gauntlet-notices", rows)
}
