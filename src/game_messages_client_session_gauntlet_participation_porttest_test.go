//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"reflect"
	"testing"
)

func TestGameMessageClientSessionGauntletParticipation(t *testing.T) {
	c, _, _ := newEffectsFullOwner(t)
	units, configure, _, free := c.srv.PortTestEscortPlayers()
	t.Cleanup(free)
	configure(2)
	pl, other := units[0].UpdateDataPlayer().Player, units[1].UpdateDataPlayer().Player
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	var sounds [][2]int
	t.Cleanup(legacy.PortTestClientSoundObserver(func(id, v int) { sounds = append(sounds, [2]int{id, v}) }))
	type row struct {
		On, Host, Present int
		Code              uint16
		Initial, After    uint32
		Return            int
		Sounds            [][2]int
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for host := 0; host < 2; host++ {
			for present := 0; present < 2; present++ {
				for _, code := range []uint16{17, 0x8011, 65535} {
					for _, initial := range []uint32{0, 1, 0xffffffff} {
						noxflags.ResetGame()
						noxflags.SetGame(noxflags.GameFlag(host))
						binary.LittleEndian.PutUint32(connected, uint32(on))
						pl.Active = byte(present)
						pl.NetCodeVal = 17
						other.NetCodeVal = 999
						pl.Field4792 = initial
						want := *pl
						otherBefore := *other
						var wantSounds [][2]int
						if on != 0 {
							wantSounds = [][2]int{{1008, 100}}
							if present != 0 && host == 0 && code == 17 {
								want.Field4792 = 1
							}
						}
						sounds = nil
						data := []byte{240, 1, byte(code), byte(code >> 8)}
						input := bytes.Clone(data)
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(240), data)
						if n != 4 || !bytes.Equal(input, data) || *pl != want || *other != otherBefore || !reflect.DeepEqual(sounds, wantSounds) {
							t.Fatal("quest participation", on, host, present, code, initial, n, sounds)
						}
						rows = append(rows, row{on, host, present, code, initial, pl.Field4792, n, append([][2]int(nil), sounds...)})
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-gauntlet-participation", rows)
}
