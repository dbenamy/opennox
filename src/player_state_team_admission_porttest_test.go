//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestPlayerStateTeamAdmission(t *testing.T) {
	o := newMatchRosterOwner(t)
	type row struct {
		Name   string
		Result int
	}
	var rows []row
	for _, flags := range []uint32{0, 128, 0x8000, 0x8020, 0x8040, 0x8010} {
		for _, special := range []byte{0, 127, 128, 255} {
			for _, gameplay := range []bool{false, true} {
				for _, groups := range []int{0, 1, 2, 3} {
					for _, capacity := range []byte{0, 1, 2, 3, 255} {
						for _, flagCapacity := range []uint32{0, 2, 4, 0xffffffff} {
							for _, remembered := range []uint32{0, 100, 999} {
								name := fmt.Sprintf("flags%x/special%d/teams%t/groups%d/cap%d/flags%d/remember%d", flags, special, gameplay, groups, capacity, flagCapacity, remembered)
								t.Run(name, func(t *testing.T) {
									defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
									o.s.Teams.Reset()
									o.s.Teams.ActiveCnt = 0
									for i := 0; i < groups; i++ {
										tm := o.s.Teams.Create(server.TeamID(i + 1))
										legacy.PortTestTeamRuntimeOther("group", tm, nil, 100+i, 0)
									}
									noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0))
									if gameplay {
										noxflags.SetGamePlay(4)
									}
									*memmap.PtrUint8(0x5D4594, 371380+53) = special
									*memmap.PtrUint8(0x5D4594, 371516+52) = capacity
									*o.roster["team-cap"] = flagCapacity
									pl := o.units[0].UpdateDataPlayer().Player
									pl.Field3680 = 0
									pl.Field2068 = remembered
									want := true
									if flags&0x8000 != 0 || flags&128 != 0 && special >= 128 {
										switch {
										case remembered == 0:
											want = false
										case remembered == 100 && groups > 0:
											want = true
										default:
											limit := int(capacity)
											if flags&96 != 0 || flags&16 != 0 && gameplay {
												if limit > 2 {
													limit = 2
												}
											}
											want = groups < limit && (flags&96 == 0 || int32(groups) < int32(flagCapacity))
										}
									}
									got := legacy.PortTestPlayerStateQuery("admission", pl, nil)
									if got != bool2int(want) {
										t.Errorf("admission=%d want%t", got, want)
									}
									rows = append(rows, row{name, got})
								})
							}
						}
					}
				}
			}
		}
	}
	spellbookCapture(t, "player-state-team-admission", rows, "4975d06634b697d4cd7f43c1abeb1efcd14677eecfd826fa523419024062adcf")
}
