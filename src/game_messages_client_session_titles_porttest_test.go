//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"reflect"
	"testing"
	"unsafe"
)

// nox_swprintf renders a null wide-string pointer as "(null)"; an allocated
// empty string remains empty. Keep that distinction when using Go strings.
func TestGameMessageClientSessionMissingQuestTitles(t *testing.T) {
	c, _, _ := newEffectsFullOwner(t)
	units, configure, _, free := c.srv.PortTestEscortPlayers()
	t.Cleanup(free)
	configure(2)
	pl, other := units[0].UpdateDataPlayer().Player, units[1].UpdateDataPlayer().Player
	oldServer, oldPlayer := noxServer, legacy.Get_dword_8531A0_2576()
	noxServer = c.srv
	t.Cleanup(func() { noxServer = oldServer; legacy.Set_dword_8531A0_2576(oldPlayer) })
	printer, index, region := teamRuntimeJoinTextOwner(t, c.srv.Server)
	var entries []strman.Entry
	for _, category := range []string{"Spell", "Guide"} {
		prefix := "plyrspel.c:"
		if category == "Guide" {
			prefix = "PlyrGide.c:"
		}
		entries = append(entries, strman.Entry{ID: strman.ID(prefix + "Award" + category), Vals: []strman.Variant{{Str: "Self [%s]"}}}, strman.Entry{ID: strman.ID(prefix + "Award" + category + "ToOther"), Vals: []strman.Variant{{Str: "%s [%s]"}}})
	}
	entries = append(entries, strman.Entry{ID: "guimsg.c:systemmsg", Vals: []strman.Variant{{Str: "%s"}}})
	set, restore := c.srv.PortTestMeterStrings(entries...)
	t.Cleanup(restore)
	set(0)
	defs, restoreDefs := c.srv.PortTestAISpellDefs()
	t.Cleanup(restoreDefs)
	defs([]server.PortTestSpellClassDef{{Index: 1, Valid: true}})
	c.srv.Spells.DefByInd(1).Title = ""
	guides := serverConfigOwnBytes(t, 0x5D4594, 740076, 41*28)
	clear(guides)
	binary.LittleEndian.PutUint32(guides[28:], uint32(uintptr(unsafe.Pointer(alloc.InternCString16("")))))
	binary.LittleEndian.PutUint32(guides[32:], 1)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	pl.Active = 1
	pl.NetCodeVal = 17
	pl.SetName("Player Ω")
	for on := 0; on < 2; on++ {
		for _, kind := range []byte{30, 31} {
			for _, id := range []byte{0, 1, 255} {
				for local := 0; local < 2; local++ {
					binary.LittleEndian.PutUint32(connected, uint32(on))
					legacy.Set_dword_8531A0_2576(other)
					if local != 0 {
						legacy.Set_dword_8531A0_2576(pl)
					}
					clear(region)
					*index = 0
					printer.lines = nil
					data := []byte{240, kind, id, 17, 0}
					before := bytes.Clone(data)
					n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, 240, data)
					var want []string
					if on != 0 {
						title := "(null)"
						if id == 1 {
							title = ""
						}
						who := "Player Ω"
						if local != 0 {
							who = "Self"
						}
						want = []string{fmt.Sprintf("%s [%s]", who, title)}
					}
					if n != 5 || !bytes.Equal(data, before) || !reflect.DeepEqual(printer.lines, want) {
						t.Fatalf("on=%d kind=%d id=%d local=%d: %q want %q", on, kind, id, local, printer.lines, want)
					}
				}
			}
		}
	}
}
