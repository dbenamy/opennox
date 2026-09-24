//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestPlayerDeathRecentAssist(t *testing.T) {
	o := newMatchRosterOwner(t)
	oldAbilities := noxServer.abilities
	noxServer.abilities.Init(noxServer)
	t.Cleanup(func() { noxServer.abilities = oldAbilities })
	t.Cleanup(legacy.PortTestPlayerDeathLookupOwner())
	callback := server.PortTestPlayerDeathCallback()
	if callback == nil {
		t.Fatal("PlayerDie registration")
	}
	type record struct {
		Name           string
		Scores, Deaths [3]uint32
		Recent         uint32
		Reports        legacy.PortTestReliableReportState
	}
	var rows []record
	for _, frame := range []uint32{0, 123, 0xfffffff0} {
		for _, age := range []uint32{0, 299, 300, 301, 0xffffffff} {
			for candidate := 0; candidate < 3; candidate++ {
				for presence := 0; presence < 4; presence++ {
					for enabled := 0; enabled < 2; enabled++ {
						name := fmt.Sprintf("frame=%x/age=%x/candidate=%d/presence=%d/enabled=%d", frame, age, candidate, presence, enabled)
						t.Run(name, func(t *testing.T) {
							o.reset()
							o.s.SetFrame(frame)
							o.s.PortTestCombatAudioReset()
							noxflags.ResetGame()
							noxflags.SetGame(noxflags.GameHost | noxflags.GameModeArena)
							noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0))
							for i := range o.units {
								u := &o.units[i]
								u.ObjFlags = 0
								u.TeamVal = server.ObjectTeam{}
								ud := u.UpdateDataPlayer()
								ud.State = 0
								pl := ud.Player
								pl.PlayerUnit = u
								pl.NetCodeVal = u.NetCode
								pl.Active = 1
								pl.Field3680 = 0
								pl.Lessons = 7
								pl.Field2140 = 11
								objectXferSetWord(unsafe.Pointer(ud), 280, 0)
							}
							u := &o.units[0]
							pl := u.UpdateDataPlayer().Player
							objectXferSetWord(u.CObj(), 520, uint32(uintptr(o.units[1].CObj())))
							objectXferSetWord(u.CObj(), 524, 0)
							info := o.units[candidate].UpdateDataPlayer().Player
							objectXferSetWord(pl.C(), 3600, uint32(enabled))
							objectXferSetWord(pl.C(), 3604, uint32(info.PlayerInd))
							objectXferSetWord(pl.C(), 3608, frame-age)
							if presence == 1 {
								info.Active = 0
							}
							if presence == 2 {
								info.PlayerUnit = nil
							}
							if presence == 3 {
								info.NetCodeVal = 0x76543210
							}
							scores := [3]uint32{7, 8, 7}
							deaths := [3]uint32{12, 11, 11}
							if enabled != 0 && age < 300 && presence == 0 && candidate == 2 {
								scores[2]++
							}
							callback(u)
							r := record{Name: name, Recent: objectXferGetWord(pl.C(), 3600), Reports: o.state()}
							for i := range o.units {
								p := o.units[i].UpdateDataPlayer().Player
								r.Scores[i] = uint32(p.Lessons)
								r.Deaths[i] = p.Field2140
							}
							if r.Scores != scores || r.Deaths != deaths || r.Recent != 0 {
								t.Fatalf("score/death/recent %v/%v/%d want %v/%v/0", r.Scores, r.Deaths, r.Recent, scores, deaths)
							}
							rows = append(rows, r)
						})
					}
				}
			}
		}
	}
	spellbookCapture(t, "player-death-recent-assist", rows, "8f05687cde6b2c9ed738c998d2584aca47e10c0a8577082276ca380fc778eb60")
}
