//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/strman"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"reflect"
	"strings"
	"testing"
)

func TestGameMessageClientSessionSettings2(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	words, restore := legacy.PortTestClientSessionWords()
	t.Cleanup(restore)
	messageWords, restoreMessages := legacy.PortTestClientInteractionWords()
	t.Cleanup(restoreMessages)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	settings := serverConfigOwnBytes(t, 0x5D4594, 371380, 58)
	saved := serverConfigOwnBytes(t, 0x5D4594, 1200708, 58)
	enabled := serverConfigOwnBytes(t, 0x587000, 4660, 4)
	deadline := serverConfigOwnBytes(t, 0x5D4594, 3468, 8)
	oldTicks := legacy.PlatformTicks
	legacy.PlatformTicks = func() uint64 { return 0x100000007 }
	t.Cleanup(func() { legacy.PlatformTicks = oldTicks })
	configure, restoreStrings := o.c.srv.PortTestMeterStrings(strman.Entry{ID: "guimsg.c:systemmsg", Vals: []strman.Variant{{Str: "System: %s"}}}, strman.Entry{ID: "cdecode.c:OptionsChanged", Vals: []strman.Variant{{Str: "Options changed"}}})
	t.Cleanup(restoreStrings)
	configure(0)
	s, desc := server.PortTestRuleServerSetup()
	oldSpells := o.c.srv.Spells
	o.c.srv.Spells = s.Spells
	t.Cleanup(func() { o.c.srv.Spells = oldSpells })
	var sounds [][2]int
	t.Cleanup(legacy.PortTestClientSoundObserver(func(id, volume int) { sounds = append(sounds, [2]int{id, volume}) }))
	type row struct {
		On, Host, Dirty, Change int
		Name                    string
		Duration                uint32
		Settings, Saved         []byte
		Enabled, Changed        uint32
		Deadline                uint64
		Spells                  []bool
		Console                 []string
		Sounds                  [][2]int
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for host := 0; host < 2; host++ {
			for dirty := 0; dirty < 2; dirty++ {
				for change := 0; change < 4; change++ {
					for _, name := range []string{"", "short", strings.Repeat("N", 16)} {
						for _, duration := range []uint32{0, 1, 0x80000000, 0xffffffff} {
							o.reset(t)
							noxflags.ResetGame()
							noxflags.SetGame(noxflags.GameFlag(host))
							binary.LittleEndian.PutUint32(connected, uint32(on))
							*words["settingsChanged"] = uint32(dirty)
							*messageWords["dword_5d4594_825736"] = 0
							sounds = nil
							for i := range settings {
								settings[i] = 0x5a
							}
							for i := range saved {
								saved[i] = 0xa5
							}
							binary.LittleEndian.PutUint32(enabled, 9)
							binary.LittleEndian.PutUint64(deadline, 0x123456789)
							for _, def := range desc.Spells {
								o.c.srv.Spells.DefByInd(spell.ID(def.Index)).Enabled = true
							}
							data := make([]byte, 49)
							data[0] = 176
							copy(data[1:17], name)
							for i := 17; i < 37; i++ {
								data[i] = byte(13*i + change)
							}
							binary.LittleEndian.PutUint32(data[37:], 0x80000001)
							binary.LittleEndian.PutUint32(data[41:], 0xdeadbeef)
							binary.LittleEndian.PutUint32(data[45:], duration)
							input := bytes.Clone(data)
							copy(saved[24:52], data[17:45])
							if change != 0 {
								saved[24+[]int{0, 0, 20, 24}[change]] ^= 1
							}
							wantSettings, wantSaved := bytes.Clone(settings), bytes.Clone(saved)
							wantEnabled := uint32(9)
							wantDeadline := uint64(0x123456789)
							wantDirty := uint32(dirty)
							if host == 0 {
								clear(wantSettings[9:24])
								copy(wantSettings[9:24], name)
								copy(wantSettings[24:52], data[17:45])
								copy(wantSaved, wantSettings)
								if change != 0 {
									wantDirty = 1
								}
								wantEnabled = 0
								if duration != 0 {
									wantEnabled = 1
									wantDeadline = uint64(int64(int32(duration)) + 7)
								}
							}
							n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(176), data)
							var wantConsole []string
							var wantSounds [][2]int
							if host == 0 && on != 0 && wantDirty != 0 {
								wantConsole = []string{"6:System: Options changed"}
								wantSounds = [][2]int{{310, 100}}
							}
							if n != 49 || !bytes.Equal(input, data) || !bytes.Equal(settings, wantSettings) || !bytes.Equal(saved, wantSaved) || *words["settingsChanged"] != wantDirty || binary.LittleEndian.Uint32(enabled) != wantEnabled || binary.LittleEndian.Uint64(deadline) != wantDeadline || !reflect.DeepEqual(o.console, wantConsole) || !reflect.DeepEqual(sounds, wantSounds) {
								t.Fatal("settings2 state/notification", on, host, dirty, change, name, duration, o.console, sounds)
							}
							var states []bool
							for _, def := range desc.Spells {
								got := o.c.srv.Spells.DefByInd(spell.ID(def.Index)).Enabled
								want := true
								if host == 0 {
									want = data[17+def.Index/8]&(1<<uint(def.Index%8)) != 0
								}
								if got != want {
									t.Fatal("settings spell mask", def.Index)
								}
								states = append(states, got)
							}
							rows = append(rows, row{on, host, dirty, change, name, duration, bytes.Clone(settings), bytes.Clone(saved), binary.LittleEndian.Uint32(enabled), *words["settingsChanged"], binary.LittleEndian.Uint64(deadline), states, append([]string(nil), o.console...), append([][2]int(nil), sounds...)})
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-settings2", rows)
}
