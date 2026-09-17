//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/spell"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestServerConfigSettingsRefresh(t *testing.T) {
	type row struct {
		Flags          uint32
		Headless       bool
		Pattern        int
		Map            string
		Video          byte
		Record, Config []byte
		First, Second  uint64
	}
	var rows []row
	for _, flags := range []uint32{0x100, 0x101, 0x20, 0x21} {
		for _, headless := range []bool{false, true} {
			for pattern := 0; pattern < 2; pattern++ {
				for _, mapName := range []string{"arena", "map12345", "abcdefghijkl"} {
					for _, video := range []byte{0x13, 0xa5} {
						t.Run(fmt.Sprintf("%x/%t/%d/%s/%x", flags, headless, pattern, mapName, video), func(t *testing.T) {
							o := newServerOptionsOwner(t)
							defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
							engine := noxflags.GetEngine()
							noxflags.ResetEngine()
							if headless {
								noxflags.SetEngine(noxflags.EngineNoRendering)
							}
							t.Cleanup(func() { noxflags.ResetEngine(); noxflags.SetEngine(engine) })
							oldVideo := legacy.Sub_43BE50_get_video_mode_id
							legacy.Sub_43BE50_get_video_mode_id = func() int { return 3 }
							t.Cleanup(func() { legacy.Sub_43BE50_get_video_mode_id = oldVideo })
							s, desc := server.PortTestRuleServerSetup()
							oldTypes, oldWeapons, oldArmor, oldSpells := o.c.srv.Types, o.c.srv.Weapons, o.c.srv.Armor, o.c.srv.Spells
							o.c.srv.Types, o.c.srv.Weapons, o.c.srv.Armor, o.c.srv.Spells = s.Types, s.Weapons, s.Armor, s.Spells
							t.Cleanup(func() {
								o.c.srv.Types, o.c.srv.Weapons, o.c.srv.Armor, o.c.srv.Spells = oldTypes, oldWeapons, oldArmor, oldSpells
							})
							for _, it := range desc.Weapons {
								s.Types.ByID(it.Name).SetAllowed(true)
							}
							for _, it := range desc.Armors {
								s.Types.ByID(it.Name).SetAllowed(true)
							}
							for _, it := range desc.Spells {
								s.Spells.DefByInd(spell.ID(it.Index)).Enabled = true
							}
							mask := bytes.Repeat([]byte{255}, 28)
							if pattern != 0 {
								id := int(spell.SPELL_FIREBALL)
								s.Spells.DefByInd(spell.ID(id)).Enabled = false
								mask[id/8] &^= 1 << uint(id&7)
								weapon, armor := desc.Weapons[0], desc.Armors[0]
								s.Types.ByID(weapon.Name).SetAllowed(false)
								s.Types.ByID(armor.Name).SetAllowed(false)
								binary.LittleEndian.PutUint32(mask[20:], 0xffffffff&^weapon.Bit)
								binary.LittleEndian.PutUint32(mask[24:], 0xffffffff&^armor.Bit)
							}
							for i := range o.players {
								o.players[i].Active = 0
								o.players[i].PlayerInd = byte(i)
							}
							count, limit := 0, 0
							if pattern != 0 {
								count, limit = 3, 32
							}
							for i := 0; i < count; i++ {
								o.players[i].Active = 1
							}
							legacy.PortTestServerConfigScalar("limit-set", limit, 0)
							serverName := unsafe.Slice(memmap.PtrUint8(0x5D4594, 1324), 16)
							clear(serverName)
							copy(serverName, "HostABCDEFGHIJK")
							currentMap := serverConfigOwnBytes(t, 0x85B3FC, 36, 80)
							clear(currentMap)
							copy(currentMap, mapName)
							for i := range o.settings {
								o.settings[i] = 0xa5
							}
							oldFlags := uint16(flags) ^ 0x1000
							if pattern != 0 {
								oldFlags = uint16(flags) ^ 1
							}
							binary.LittleEndian.PutUint16(o.settings[52:], oldFlags)
							config := unsafe.Slice(memmap.PtrUint8(0x5D4594, 371516), 184)
							for i := range config {
								config[i] = 0x5a
							}
							config[102] = video
							expected := append([]byte(nil), o.settings...)
							clear(expected[:24])
							copy(expected[:8], mapName)
							copy(expected[9:24], "HostABCDEFGHIJK")
							copy(expected[24:44], mask[:20])
							if flags&1 != 0 {
								copy(expected[44:52], mask[20:])
							}
							selected := oldFlags
							if (oldFlags^uint16(flags))&0xfff0 != 0 {
								selected = uint16(flags)
							}
							binary.LittleEndian.PutUint16(expected[52:], selected)
							index := 0
							if flags&0x20 != 0 {
								index = 2
							}
							binary.LittleEndian.PutUint16(expected[54:], uint16(101+11*index))
							expected[56] = byte(13 + 7*index)
							expectedConfig := append([]byte(nil), config...)
							if video&0xef != 3 {
								expectedConfig[102] = video&0x80 | 3
							}
							if headless {
								count--
								limit--
							}
							expectedConfig[103], expectedConfig[104] = byte(count), byte(limit)
							*o.optionWords["settings-record-dirty"] = 0
							first := legacy.PortTestServerConfigScalar("refresh", 0, 0)
							if !bytes.Equal(o.settings, expected) || !bytes.Equal(config, expectedConfig) || first != 1 {
								t.Fatalf("refresh fields got %x / %x return%d want %x / %x", o.settings, config, first, expected, expectedConfig)
							}
							*o.optionWords["settings-record-dirty"] = 0
							second := legacy.PortTestServerConfigScalar("refresh", 0, 0)
							// C compares unsigned stored counts against signed decremented values,
							// and compares weapon bytes against a signed-char temporary. Hosted
							// snapshots retain high unused bits, so that comparison reports dirty again.
							wantDirty := uint64(bool2int(expectedConfig[102]&0xef != 3 || count < 0 || limit < 0 || flags&1 != 0))
							if !bytes.Equal(o.settings, expected) || !bytes.Equal(config, expectedConfig) || second != wantDirty {
								t.Fatal("repeat refresh", second, wantDirty)
							}
							rows = append(rows, row{flags, headless, pattern, mapName, video, append([]byte(nil), o.settings...), append([]byte(nil), config...), first, second})
						})
					}
				}
			}
		}
	}
	spellbookCapture(t, "server-config-settings-refresh", rows, "c24d4fc00f6e411b7a14072871a0924e19d4e8ff5cadcacf4a379c6f3fa92d77")
}
