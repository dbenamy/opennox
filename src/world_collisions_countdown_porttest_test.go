//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestWorldCollisionsCountdown(t *testing.T) {
	o := newWorldCollisionOwner(t)
	defer noxflags.PortTestGameFlags(0)()
	enabled := memmap.PtrUint32(0x587000, 4660)
	deadline := memmap.PtrUint64(0x5D4594, 3468)
	oldEnabled, oldDeadline, oldFlag := *enabled, *deadline, noxServer.flag3592
	t.Cleanup(func() { *enabled = oldEnabled; *deadline = oldDeadline; noxServer.flag3592 = oldFlag })
	var rows []struct {
		Name    string
		Enabled uint32
		Started bool
		Seconds int
	}
	defer func() {
		spellbookCapture(t, "world-collisions-countdown", rows, "6b13fe19017f1577e192fb9ff72071356f919e575621ebea8af9162ac1e857cb")
	}()
	for _, participants := range []int{0, 1, 2, 3} {
		for ready := 0; ready <= participants; ready++ {
			for _, seconds := range []int{0, 1, 29, 30, 61} {
				for _, running := range []bool{false, true} {
					name := fmt.Sprintf("participants%d/ready%d/seconds%d/running%t", participants, ready, seconds, running)
					t.Run(name, func(t *testing.T) {
						o.reset()
						o.balance(map[string]float64{"QuestExitTimerStart": float64(seconds)})
						*enabled = 0
						*deadline = 777
						noxServer.flag3592 = running
						if running {
							*enabled = 1
							*deadline = o.ticks + uint64(seconds)*1000
						}
						initialDeadline := *deadline
						for i := range o.units {
							u := &o.units[i]
							p := u.UpdateDataPlayer().Player
							objectXferSetWord(p.C(), 4792, 0)
							objectXferSetWord(u.UpdateData, 312, 0)
							if i < participants {
								objectXferSetWord(p.C(), 4792, 1)
							}
							if i < ready {
								objectXferSetWord(u.UpdateData, 312, uint32(uintptr(u.CObj())))
							}
						}
						before := platformTicks()
						legacy.PortTestWorldCollision(4, nil, nil, nil)
						after := platformTicks()
						wantSeconds := seconds
						if participants > 0 {
							wantSeconds -= int(float32(float64(ready) / float64(participants) * float64(seconds)))
						}
						started := participants > 0 && (!running || wantSeconds < seconds)
						wantEnabled := uint32(0)
						if participants > 0 {
							wantEnabled = 1
						}
						if *enabled != wantEnabled {
							t.Fatal("countdown enabled")
						}
						if started {
							min, max := before+uint64(wantSeconds)*1000, after+uint64(wantSeconds)*1000
							if *deadline < min || *deadline > max || !noxServer.flag3592 {
								t.Fatal("countdown deadline")
							}
						} else if *deadline != initialDeadline {
							t.Fatal("unchanged deadline")
						}
						rows = append(rows, struct {
							Name    string
							Enabled uint32
							Started bool
							Seconds int
						}{name, *enabled, started, wantSeconds})
					})
				}
			}
		}
	}
}
