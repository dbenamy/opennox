//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"reflect"
	"testing"
	"unsafe"
)

func TestGameMessageClientSessionGauntletScoreboard(t *testing.T) {
	o := newScoreboardOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		On, Present, Hidden, Quest int
		Mode                       uint32
		State                      scoreboardResult
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for present := 0; present < 2; present++ {
			for hidden := 0; hidden < 2; hidden++ {
				for quest := 0; quest < 2; quest++ {
					for _, mode := range []uint32{0, 1, 2, 4} {
						setup := func() {
							o.resetRank(t)
							if present != 0 {
								o.constructRank(t)
								w := (*gui.Window)(unsafe.Pointer(uintptr(*o.rankWords["dword_5d4594_1090048"])))
								w.Show()
								if hidden != 0 {
									w.Hide()
								}
							}
							if quest != 0 {
								noxflags.SetGame(4096)
							}
							*o.rankWords["dword_5d4594_1090120"] = mode
							binary.LittleEndian.PutUint32(connected, uint32(on))
						}
						setup()
						if on != 0 && legacy.PortTestScoreboard(25, 0, 0, 0) == 0 {
							legacy.PortTestScoreboard(26, 0, 0, 0)
						}
						want := o.rankCapture(t, 0, 0)
						setup()
						data := []byte{240, 20}
						input := bytes.Clone(data)
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(240), data)
						got := o.rankCapture(t, 0, 0)
						if n != 2 || !bytes.Equal(input, data) || !reflect.DeepEqual(got, want) {
							t.Fatal("quest scoreboard", on, present, hidden, quest, mode, n)
						}
						rows = append(rows, row{on, present, hidden, quest, mode, got})
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-gauntlet-scoreboard", rows)
}
func TestGameMessageClientSessionGauntletLevel(t *testing.T) {
	o := newInventoryWindowOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		On    int
		Level uint16
		State inventoryWindowResult
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, level := range []uint16{0, 1, 255, 32767, 32768, 65535} {
			setup := func() {
				o.reset(t)
				o.construct(t)
				*o.windowWords["dword_5d4594_1049844"] = 99
				binary.LittleEndian.PutUint32(connected, uint32(on))
			}
			setup()
			if on != 0 {
				legacy.PortTestInventoryWindow(14, uintptr(level), 0, 0, 0)
			}
			want := o.capture(t, 0, 0, 7, 0, 0, 0, 0)
			setup()
			data := []byte{240, 29, byte(level), byte(level >> 8)}
			input := bytes.Clone(data)
			n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(240), data)
			got := o.capture(t, 0, 0, 7, 0, 0, 0, 0)
			if n != 4 || !bytes.Equal(input, data) || !reflect.DeepEqual(got, want) {
				t.Fatal("quest inventory level", on, level, n)
			}
			rows = append(rows, row{on, level, got})
		}
	}
	interactionCapture(t, "game-client-session-gauntlet-level", rows)
}
