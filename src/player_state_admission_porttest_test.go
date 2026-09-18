//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestPlayerStateAdmission(t *testing.T) {
	o := newMatchRosterOwner(t)
	serverConfigOwnBytes(t, 0x5D4594, 3520, 4)
	type row struct {
		Name   string
		Result int
	}
	var rows []row
	for _, flags := range []uint32{0, 128, 1024, 1152, 4096, 5120} {
		for _, count := range []int{0, 1, 5, 6, 7} {
			for _, status := range []uint32{0, 256} {
				for _, reentry := range []uint32{0, 1, 0xffffffff} {
					for _, score := range []uint32{0, 1, 0x7fffffff, 0x80000000, 0xffffffff} {
						for _, delta := range []uint32{600, 601, 0xffffffff} {
							for _, nilPlayer := range []bool{false, true} {
								name := fmt.Sprintf("flags%x/count%d/status%x/reentry%x/score%x/delta%x/nil%t", flags, count, status, reentry, score, delta, nilPlayer)
								t.Run(name, func(t *testing.T) {
									defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
									for i := 0; i < 32; i++ {
										p := o.s.Players.ByIndRaw(ntype.PlayerInd(i))
										p.Active = 0
										p.Field2140 = 0
										p.Field3680 = 0
									}
									for i := 0; i < count; i++ {
										p := o.s.Players.ByIndRaw(ntype.PlayerInd(i))
										p.Active = 1
										p.PlayerInd = byte(i)
									}
									candidate := o.s.Players.ByIndRaw(30)
									candidate.Field3680 = status
									if count > 0 {
										o.s.Players.ByIndRaw(0).Field2140 = score
									}
									*memmap.PtrUint32(0x5D4594, 3508) = reentry
									*memmap.PtrUint32(0x5D4594, 3520) = 0xfffffff0
									o.s.SetTickRate(30)
									o.s.SetFrame(0xfffffff0 + delta)
									var pl *server.Player = candidate
									if nilPlayer {
										pl = nil
									}
									want := true
									switch {
									case !nilPlayer && flags&4096 != 0:
										want = count < 6
									case !nilPlayer && status&256 != 0 && reentry == 0:
										want = false
									case flags&128 != 0:
										want = true
									case flags&1024 != 0:
										want = count == 0 || int32(score) <= 0 || delta <= 600
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
	spellbookCapture(t, "player-state-admission", rows, "150b26ab05b80357613901711a2f326214e3947e3f00a65519969af168295678")
}
