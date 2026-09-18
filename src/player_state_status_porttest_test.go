//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestPlayerStateStatus(t *testing.T) {
	o := newMatchRosterOwner(t)
	type row struct {
		Name   string
		Status uint32
		Queue  legacy.PortTestReliableReportState
	}
	var rows []row
	for _, host := range []bool{false, true} {
		for _, op := range []string{"add", "remove", "report"} {
			for _, before := range []uint32{0, 1, 0x423, 0x100, 0xffffffff} {
				for _, mask := range []uint32{0, 1, 2, 0x20, 0x400, 0x423, 0x100, 0x80000000, 0xffffffff} {
					for _, code := range []uint32{0, 0x12345678, 0xffffffff} {
						name := fmt.Sprintf("host%t/%s/before%x/mask%x/code%x", host, op, before, mask, code)
						t.Run(name, func(t *testing.T) {
							o.reset()
							// Chat mode suppresses timer startup, tested separately below.
							flags := uint32(128)
							if host {
								flags |= 1
							}
							defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
							pl := o.units[0].UpdateDataPlayer().Player
							pl.Field3680 = before
							pl.NetCodeVal = code
							legacy.PortTestPlayerStateStatus(op, pl, mask)
							want := before
							if op == "add" {
								want |= mask
							}
							if op == "remove" {
								want &^= mask
							}
							state := o.state()
							reports := bool2int(op == "report" || host && mask&0x423 != 0)
							if pl.Field3680 != want || len(state.Nodes) != reports {
								t.Fatal("status/messaging", pl.Field3680, want, len(state.Nodes), reports)
							}
							if reports != 0 {
								data := make([]byte, 7)
								data[0] = 106
								binary.LittleEndian.PutUint16(data[1:], uint16(code))
								binary.LittleEndian.PutUint32(data[3:], want&0x423)
								n := state.Nodes[0]
								if !bytes.Equal(n.Data, data) || n.To != 255 || n.Recipients != 0x80000082 {
									t.Fatal("status report", n, data)
								}
							}
							rows = append(rows, row{name, pl.Field3680, state})
						})
					}
				}
			}
		}
	}
	spellbookCapture(t, "player-state-status", rows, "e50bce99eee95e155f56fa161a4d83c5b3affe4d1644e9cdbc0461d7e6ecb589")
}

func TestPlayerStateStatusTimer(t *testing.T) {
	o := newMatchRosterOwner(t)
	type row struct {
		Name     string
		Enabled  uint32
		Deadline uint64
		Queue    legacy.PortTestReliableReportState
	}
	var rows []row
	for _, host := range []bool{false, true} {
		for _, chat := range []bool{false, true} {
			for _, others := range []bool{false, true} {
				for _, enabled := range []uint32{0, 1, 0xffffffff} {
					for _, minutes := range []byte{0, 1, 255} {
						name := fmt.Sprintf("host%t/chat%t/others%t/enabled%x/minutes%d", host, chat, others, enabled, minutes)
						t.Run(name, func(t *testing.T) {
							o.reset()
							oldTicks := legacy.PlatformTicks
							tick := uint64(10000)
							legacy.PlatformTicks = func() uint64 { return tick }
							defer func() { legacy.PlatformTicks = oldTicks }()
							flags := uint32(256)
							if host {
								flags |= 1
							}
							if chat {
								flags |= 128
							}
							defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
							for i := range o.units {
								p := o.units[i].UpdateDataPlayer().Player
								p.Field3680 = 1
								if i != 0 && others {
									p.Field3680 = 0
								}
							}
							*memmap.PtrUint32(0x587000, 4660) = enabled
							*memmap.PtrUint64(0x5D4594, 3468) = 123456
							*memmap.PtrUint8(0x5D4594, 3500) = minutes
							pl := o.units[0].UpdateDataPlayer().Player
							legacy.PortTestPlayerStateStatus("remove", pl, 1)
							start := host && !chat && others && enabled == 0 && minutes != 0
							wantEnabled := enabled
							wantDeadline := uint64(123456)
							if start {
								wantEnabled = 1
								wantDeadline = uint64(uint32(legacy.PlatformTicks())) + uint64(minutes)*60000
							}
							if memmap.Uint32(0x587000, 4660) != wantEnabled || memmap.Uint64(0x5D4594, 3468) != wantDeadline {
								t.Fatal("timer transition", start, memmap.Uint32(0x587000, 4660), memmap.Uint64(0x5D4594, 3468), wantDeadline)
							}
							// Repeating the status clear after time advances must not restart the timer.
							tick += 1000
							legacy.PortTestPlayerStateStatus("remove", pl, 1)
							if memmap.Uint64(0x5D4594, 3468) != wantDeadline {
								t.Fatal("timer restarted")
							}
							rows = append(rows, row{name, memmap.Uint32(0x587000, 4660), memmap.Uint64(0x5D4594, 3468), o.state()})
						})
					}
				}
			}
		}
	}
	spellbookCapture(t, "player-state-status-timer", rows, "b67ad580389467d5e4f63b3617430fd58fa1dab0419f38134f90de86b92a79ca")
}
