//go:build porttest

package opennox

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestGameMessageServerGauntletRespawn(t *testing.T) {
	var direct, messages []legacy.PortTestRoamSpec
	for class := byte(0); class < 3; class++ {
		for _, present := range []bool{false, true} {
			for _, dead := range []bool{false, true} {
				for _, delay := range []uint32{0, 123, 0xffffffff} {
					makeCase := func(dispatch bool) legacy.PortTestRoamSpec {
						op := 67
						if present && dead {
							op = 79
						}
						s := controlsBase(op)
						controlsStats(&s)
						p := s.Callbacks.Shop
						o := p.TemporaryUpdates.World.Objectives
						a := o.Attack
						p.Resources.PlayerClass = class
						a.Controls.Equipment, a.Controls.Corpse = true, true
						a.ActorWords = map[int]uint32{16: 0}
						if dead {
							a.ActorWords[16] = 0x8000
						}
						a.UpdateWords = map[int]uint32{548: delay}
						o.PlayerDataWords[0] = map[int]uint32{3684: 5, 4700: 1}
						ref := 0
						if present {
							ref = 1
						}
						o.PlayerDataRefs = []map[int]int{{2056: ref}}
						if dispatch {
							a.Controls.GameMessage, a.Controls.GameMessageLength = []byte{240, 3}, 2
						}
						return s
					}
					direct = append(direct, makeCase(false))
					messages = append(messages, makeCase(true))
				}
			}
		}
	}
	want, got := controlsRun(t, direct), controlsRun(t, messages)
	gameMessageStateGuards(t, want)
	gameMessageStateGuards(t, got)
	for i := range got {
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Fatalf("gauntlet respawn differs from qualified operation: case %d", i)
		}
	}
	interactionCapture(t, "game-server-gauntlet-respawn", got)
}

func TestGameMessageServerGauntletLeave(t *testing.T) {
	o := newReliableReportsOwner(t)
	old := legacy.Sub_4DD0B0
	t.Cleanup(func() { legacy.Sub_4DD0B0 = old })
	for i := range o.units {
		u := &o.units[i]
		calls := 0
		legacy.Sub_4DD0B0 = func(got *server.Object) {
			if got != u {
				t.Fatal("gauntlet leave routed to wrong player")
			}
			calls++
		}
		data := []byte{240, 27}
		pl := u.UpdateDataPlayer().Player
		if n := legacy.Nox_xxx_netOnPacketRecvServ_51BAD0_net_sdecode_switch(pl.PlayerIndex(), data, pl, u, u.UpdateData); n != 2 || calls != 1 || !bytes.Equal(data, []byte{240, 27}) {
			t.Fatal("gauntlet leave length, input or action count")
		}
	}
}
