//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestGameMessageClientSessionSettingsUI(t *testing.T) {
	catalog := legacy.PortTestMapCatalogOpen(9)
	t.Cleanup(catalog.Close)
	o := newServerOptionsOwner(t)
	root := *o.optionWords["root"]
	words, restore := legacy.PortTestClientSessionWords()
	t.Cleanup(restore)
	records := serverConfigOwnBytes(t, 0x5D4594, 371380, 116)
	acquired := serverConfigOwnBytes(t, 0x5D4594, 371704, 4)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	spellMask := serverConfigOwnBytes(t, 0x5D4594, 1045488, 20)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	type row struct {
		On, Open, Acquired, Slot, Toggle      int
		Mode                                  uint16
		Records, SpellMask                    []byte
		Weapons, Armor, AfterAcquired, Notice uint32
		Title, Score, Time                    string
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for open := 0; open < 2; open++ {
			for prior := 0; prior < 2; prior++ {
				for slot := 0; slot < 2; slot++ {
					for toggle := 0; toggle < 2; toggle++ {
						for _, mode := range []uint16{0x100, 0x20, 0x400, 0x1000} {
							noxflags.ResetGame()
							binary.LittleEndian.PutUint32(connected, uint32(on))
							binary.LittleEndian.PutUint32(acquired, uint32(prior))
							*words["settingsNotice"] = 99
							*o.optionWords["root"] = root
							if open == 0 {
								*o.optionWords["root"] = 0
							}
							clear(records)
							for i := 0; i < 2; i++ {
								st := records[58*i:][:58]
								copy(st, "Arena")
								copy(st[9:], "Initial")
								binary.LittleEndian.PutUint16(st[52:], 0x100)
								binary.LittleEndian.PutUint16(st[54:], 19)
								st[56] = 3
								for j := 24; j < 52; j++ {
									st[j] = byte(j)
								}
							}
							data := make([]byte, 60)
							data[0], data[1] = 177, byte(slot)
							st := data[2:]
							copy(st, "Arena")
							copy(st[9:], "Updated")
							for i := 24; i < 52; i++ {
								st[i] = byte(5*i + 3)
							}
							binary.LittleEndian.PutUint16(st[52:], mode)
							binary.LittleEndian.PutUint16(st[54:], 91)
							st[56] = 17
							records[slot*58+52] = st[52] ^ byte(toggle*32)
							want := bytes.Clone(records)
							copy(want[slot*58:][:58], st)
							if prior == 0 && open == 0 && slot == 1 {
								copy(want[:58], st)
							}
							input := bytes.Clone(data)
							n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(177), data)
							if n != 60 || !bytes.Equal(input, data) || !bytes.Equal(records, want) || binary.LittleEndian.Uint32(acquired) != 1 || *words["settingsNotice"] != 99 {
								t.Fatal("UI settings copy/acquired", on, open, prior, slot, toggle, mode)
							}
							second := records[58:]
							weapons, armor := memmap.Uint32(0x5D4594, 1045452), memmap.Uint32(0x5D4594, 1045456)
							if !bytes.Equal(spellMask, second[24:44]) || weapons != binary.LittleEndian.Uint32(second[44:]) || armor != binary.LittleEndian.Uint32(second[48:]) {
								t.Fatal("UI selection masks")
							}
							title, score, minutes := o.options.ChildByID(10119).DrawData().Text(), o.entry(10134), o.entry(10135)
							if open != 0 && slot == 1 {
								wantTitle := map[uint16]string{0x100: "Arena battle", 0x20: "Capture flag", 0x400: "Last survivor", 0x1000: "Quest adventure"}[mode]
								if title != wantTitle || score != "91" || minutes != "17" {
									t.Fatal("UI refreshed labels", title, score, minutes)
								}
							}
							rows = append(rows, row{on, open, prior, slot, toggle, mode, bytes.Clone(records), bytes.Clone(spellMask), weapons, armor, binary.LittleEndian.Uint32(acquired), *words["settingsNotice"], title, score, minutes})
						}
					}
				}
			}
		}
	}
	*o.optionWords["root"] = root
	interactionCapture(t, "game-client-session-settings-ui", rows)
}
