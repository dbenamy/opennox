//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy"
	"math/bits"
	"testing"
)

func TestGameplayReportsEliminationRules(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	type check struct {
		notified  uint32
		status    [3]uint32
		countdown bool
		team      bool
	}
	var checks []check
	for _, frame := range []uint32{599, 600, 601} {
		for _, mode := range []uint32{0, 1024, 1024 | 0x4000000} {
			for _, notified := range []uint32{0, 1} {
				for _, active := range []bool{false, true} {
					for _, teams := range []int{-1, 0, 1, 2} {
						for _, observers := range []uint32{0, 1, 3} {
							for _, limit := range []uint32{0, 2, 8} {
								s := gameplayReportsBase(36, reportObject(1))
								s.Owner.Frame = frame
								s.Owner.FPS = 30
								s.Lifecycle.GameFlags = mode
								p := s.Callbacks.Shop
								o := p.TemporaryUpdates.World.Objectives
								a := o.Attack
								o.Players = 3
								rules := &legacy.PortTestGameplayReportRules{Notified: notified, Limit: limit, CountdownActive: active, TeamMode: teams >= 0}
								if teams >= 0 {
									rules.Teams = teams
								}
								a.Controls.Reports.Rules = rules
								p.EffectsUse.Balance["SuddenDeathCountdown"] = []float64{7}
								c := check{notified: notified, team: teams >= 0}
								for j := 0; j < 3; j++ {
									flag := (observers >> j) & 1
									if j == 2 {
										flag = 1
									}
									c.status[j] = flag
									o.PlayerDataWords[j] = map[int]uint32{3680: flag, 2140: 0xffffffff}
									team := uint32(1)
									if j == 2 {
										team = 2
									}
									o.PlayerWords[j][52] = team
								}
								if mode&1024 != 0 && frame > 600 {
									c.notified = 1
									if notified == 0 {
										for j := range c.status {
											if c.status[j]&1 != 0 {
												c.status[j] |= 256
											}
										}
									}
								}
								eligible := mode&1024 != 0 && mode&0x4000000 == 0 && frame > 600 && !active
								survivors := 2 - bits.OnesCount32(observers)
								if teams < 0 {
									c.countdown = eligible && uint32(survivors+1) < limit
								} else {
									c.countdown = eligible && uint32(teams) < limit && teams > 0 && survivors == 1
								}
								a.Controls.Reports.Calls = map[uint32][5]legacy.PortTestGameplayReportArg{0: {reportObject(1)}, 1: {reportObject(1)}}
								p.Sequence = []legacy.PortTestShopAction{{Op: 1836}, {Op: 1839, Value: 1}}
								cases = append(cases, s)
								checks = append(checks, c)
							}
						}
					}
				}
			}
		}
	}
	out := controlsRun(t, cases)
	textByMode := map[bool]string{}
	for i, r := range out {
		step := r.Callbacks.Shop.Sequence[1]
		c := checks[i]
		rules := step.GameplayReports.Rules
		if rules == nil || rules.Notified != c.notified || rules.PlayerStatus != c.status {
			t.Fatalf("elimination case%d timestamp state", i)
		}
		n := 0
		if c.countdown {
			n = 1
		}
		if len(rules.Countdown) != n {
			t.Fatalf("elimination case%d countdown count%d want%d", i, len(rules.Countdown), n)
		}
		if n != 0 {
			call := rules.Countdown[0]
			if call.Seconds != 7 || call.Text != "Settings.c:SuddenDeathImminent" {
				t.Fatalf("elimination case%d countdown seconds%d text%q", i, call.Seconds, call.Text)
			}
			if text, ok := textByMode[c.team]; ok && text != call.Text {
				t.Fatalf("elimination case%d countdown text selection", i)
			}
			textByMode[c.team] = call.Text
		}
		if len(step.Packets) != 1 || len(step.Packets[0].Data) != 11 || step.Packets[0].Data[0] != 78 || binary.LittleEndian.Uint32(step.Packets[0].Data[7:]) != 0 {
			t.Fatalf("elimination case%d death counter wrap", i)
		}
	}
	if len(textByMode) != 2 {
		t.Fatal("both countdown rule paths must execute")
	}
	gameplayReportsCapture(t, "elimination-rules", out)
}
