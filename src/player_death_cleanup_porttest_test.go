//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestPlayerDeathAbilityCleanup(t *testing.T) {
	type row struct {
		Name      string
		Buffs     uint32
		Dur       [32]uint16
		Power     [32]uint8
		Cooldowns [server.AbilityMax]int
		Reports   legacy.PortTestReliableReportState
	}
	var rows []row
	for _, mask := range []uint32{0, 1, 0x80000000, 0xffffffff} {
		for _, count := range []int{0, 1, 5} {
			name := fmt.Sprintf("buffs=%x/abilities=%d", mask, count)
			t.Run(name, func(t *testing.T) {
				o := newMatchRosterOwner(t)
				oldAbilities := noxServer.abilities
				noxServer.abilities.Init(noxServer)
				t.Cleanup(func() { noxServer.abilities = oldAbilities })
				oldCore := o.s.Abils.ByUnit
				o.s.Abils.Reset()
				t.Cleanup(func() { o.s.Abils.ByUnit = oldCore })
				noxflags.ResetGame()
				noxflags.SetGame(noxflags.GameHost)
				u := &o.units[0]
				u.ObjFlags = 0
				u.Buffs = mask
				for i := range u.BuffsDur {
					u.BuffsDur[i] = uint16(100 + i)
					u.BuffsPower[i] = uint8(20 + i)
				}
				objectXferSetWord(u.CObj(), 520, 0)
				objectXferSetWord(u.CObj(), 524, 0)
				objectXferSetWord(u.UpdateDataPlayer().Player.C(), 3600, 0)
				objectXferSetWord(unsafe.Pointer(u.UpdateDataPlayer()), 280, 0)
				ad := o.s.Abils.GetFor(u)
				for i := range ad.Cooldowns {
					ad.Cooldowns[i] = 100 + i
				}
				entries := make([]server.ExecAbilityClass, count)
				for i := range entries {
					entries[i] = server.ExecAbilityClass{Abil: server.Ability(i + 1), Frame: 999, Active: 1}
					if i > 0 {
						entries[i].Prev = &entries[i-1]
						entries[i-1].Next = &entries[i]
					}
				}
				if count > 0 {
					ad.ExecList = &entries[0]
				}
				o.reset()
				ccall.CallVoidPtr(server.PortTestPlayerDeathCallback(), u.CObj())
				r := row{Name: name, Buffs: u.Buffs, Dur: u.BuffsDur, Power: u.BuffsPower, Cooldowns: ad.Cooldowns, Reports: o.state()}
				if r.Buffs != 0 || r.Dur != [32]uint16{} || r.Power != [32]uint8{} || r.Cooldowns != [server.AbilityMax]int{} || ad.ExecList != nil {
					t.Fatal("death cleanup retained buffs or abilities")
				}
				for _, entry := range entries {
					if entry != (server.ExecAbilityClass{}) {
						t.Fatal("active ability entry retained")
					}
				}
				rows = append(rows, r)
			})
		}
	}
	spellbookCapture(t, "player-death-ability-cleanup", rows, "1eeebe561b4108072aabe293a989d27651625e41899098da0478b06f5f8d8c40")
}
