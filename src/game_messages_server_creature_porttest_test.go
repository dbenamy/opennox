//go:build porttest

package opennox

import (
	"reflect"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageServerCreatureCommand(t *testing.T) {
	var direct, messages []legacy.PortTestRoamSpec
	for _, flags := range []uint32{0, 1, 2, 3, 0x100} {
		for _, code := range []uint16{0, 7, 0x8017, 0xffff} {
			for _, order := range []int{0, 1, 2, 3, 4, 5, 6, 255} {
				for _, second := range []bool{false, true} {
					makeCase := func(dispatch bool) legacy.PortTestRoamSpec {
						s := stateBase(23)
						v := s.MonsterState
						v.Source, v.Order, v.NetCode, v.Status = 2, order, 7, 0x80
						v.Own, v.Second, v.SecondEnabled, v.Broadcast = true, true, second, code == 0
						v.ServerCommand = &legacy.PortTestServerCreatureCommandSpec{Dispatch: dispatch, Invoke: flags&1 == 0 && code != 0xffff, Flags: flags, Code: code}
						return s
					}
					direct = append(direct, makeCase(false))
					messages = append(messages, makeCase(true))
				}
			}
		}
	}
	want, got := legacy.PortTestRoam(direct), legacy.PortTestRoam(messages)
	for i := range got {
		want[i].Nanos, got[i].Nanos = 0, 0
		if !got[i].Intact || !want[i].Intact || !got[i].MonsterState.Intact || !want[i].MonsterState.Intact || !got[i].Combat.Intact || !want[i].Combat.Intact {
			t.Fatalf("creature command guarded state changed: case %d", i)
		}
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Fatalf("creature command differs from qualified operation: case %d flags=%d code=%d order=%d", i, messages[i].MonsterState.ServerCommand.Flags, messages[i].MonsterState.ServerCommand.Code, messages[i].MonsterState.Order)
		}
	}
	interactionCapture(t, "game-server-creature-command", got)
}
