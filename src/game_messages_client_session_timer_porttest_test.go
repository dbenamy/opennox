//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noxnet/netmsg"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientSessionTimer(t *testing.T) {
	o := newReliableReportsOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	enabled := serverConfigOwnBytes(t, 0x587000, 4660, 4)
	deadline := serverConfigOwnBytes(t, 0x5D4594, 3468, 8)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	oldTicks := legacy.PlatformTicks
	t.Cleanup(func() { legacy.PlatformTicks = oldTicks })
	type row struct {
		On, Host, Return, Calls       int
		Frame, Age, Enabled, Duration uint32
		Ticks, Deadline               uint64
		State                         legacy.PortTestReliableReportState
		Local                         []byte
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for host := 0; host < 2; host++ {
			for _, frame := range []uint32{0, 100, 0xffffffff} {
				for _, age := range []uint32{0, 29, 30, 31, 0xffffffff} {
					for _, flag := range []uint32{0, 1, 0xffffffff} {
						for _, duration := range []uint32{0, 1, 0x80000000, 0xffffffff} {
							for _, tick := range []uint64{0, 0xffffffff, 0x100000000} {
								o.reset()
								o.s.SetFrame(frame)
								noxflags.ResetGame()
								noxflags.SetGame(noxflags.GameFlag(host))
								binary.LittleEndian.PutUint32(connected, uint32(on))
								binary.LittleEndian.PutUint32(enabled, 0x1234)
								binary.LittleEndian.PutUint64(deadline, 0x100000014)
								calls := 0
								legacy.PlatformTicks = func() uint64 { calls++; return tick }
								data := make([]byte, 13)
								data[0] = 211
								binary.LittleEndian.PutUint32(data[1:], flag)
								binary.LittleEndian.PutUint32(data[5:], duration)
								binary.LittleEndian.PutUint32(data[9:], frame-age)
								before := bytes.Clone(data)
								n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(211), data)
								state := o.state()
								local := o.s.NetList.CopyPacketsA(31, netlist.Kind0)
								wantEnabled, wantDeadline, wantCalls := uint32(0x1234), uint64(0x100000014), 0
								var packet []byte
								to, ordered := byte(0), byte(0)
								if on != 0 {
									if age >= 30 {
										if host == 0 {
											packet = []byte{212}
											to = 31
										} else if !bytes.Equal(local, []byte{212}) {
											t.Fatal("timer local retry", local)
										}
									} else {
										wantEnabled = uint32(bool2int(flag != 0))
										if host != 0 {
											wantCalls++
											packet = make([]byte, 13)
											packet[0] = 211
											to = 159
											ordered = 1
											binary.LittleEndian.PutUint32(packet[1:], wantEnabled)
											binary.LittleEndian.PutUint32(packet[5:], uint32(20-uint32(tick)))
											binary.LittleEndian.PutUint32(packet[9:], frame)
										}
										if flag != 0 {
											wantCalls++
											wantDeadline = uint64(int64(int32(duration)) + int64(uint32(tick)))
										}
									}
								}
								if n != 13 || !bytes.Equal(before, data) || binary.LittleEndian.Uint32(enabled) != wantEnabled || binary.LittleEndian.Uint64(deadline) != wantDeadline || calls != wantCalls {
									t.Fatalf("timer on%d host%d frame%d age%d flag%x duration%x tick%x return%d calls%d/%d enabled%x/%x deadline%x/%x", on, host, frame, age, flag, duration, tick, n, calls, wantCalls, binary.LittleEndian.Uint32(enabled), wantEnabled, binary.LittleEndian.Uint64(deadline), wantDeadline)
								}
								wantNodes := bool2int(packet != nil)
								if len(state.Nodes) != wantNodes {
									t.Fatal("timer queue count", len(state.Nodes), wantNodes)
								}
								if packet != nil {
									p := state.Nodes[0]
									if !bytes.Equal(p.Data, packet) || p.To != to || p.Ordered != ordered || p.Priority != 1 {
										t.Fatal("timer notification ordering/route", p, packet)
									}
								}
								if !(on != 0 && age >= 30 && host != 0) && len(local) != 0 {
									t.Fatal("unexpected local timer data", local)
								}
								rows = append(rows, row{on, host, n, calls, frame, age, flag, duration, tick, binary.LittleEndian.Uint64(deadline), state, local})
							}
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-timer", rows)
}

func TestGameMessageClientSessionLatency(t *testing.T) {
	o := newReliableReportsOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	for i := range o.units {
		o.units[i].UpdateDataPlayer().Player.Active = 0
	}
	pl := o.units[0].UpdateDataPlayer().Player
	raw := unsafe.Slice((*byte)(pl.C()), int(unsafe.Sizeof(*pl)))
	type row struct {
		On, Present, Return       int
		PlayerCode, Code, Latency uint16
		Stored                    uint16
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for present := 0; present < 2; present++ {
			for _, id := range []uint16{7, 0x8007, 65535} {
				for _, code := range []uint16{0, 7, 8, 0x8007, 65535} {
					for _, latency := range []uint16{0, 1, 255, 65535} {
						binary.LittleEndian.PutUint32(connected, uint32(on))
						pl.Active = byte(present)
						pl.NetCodeVal = uint32(id)
						binary.LittleEndian.PutUint16(raw[2148:], 0x1234)
						want := bytes.Clone(raw)
						if present != 0 && code == id {
							binary.LittleEndian.PutUint16(want[2148:], latency)
						}
						data := []byte{215, byte(code), byte(code >> 8), byte(latency), byte(latency >> 8)}
						input := bytes.Clone(data)
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(215), data)
						if n != 5 || !bytes.Equal(data, input) || !bytes.Equal(raw, want) {
							t.Fatal("latency owner/full code/connection independence", on, present, id, code, latency, n)
						}
						rows = append(rows, row{on, present, n, id, code, latency, binary.LittleEndian.Uint16(raw[2148:])})
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-latency", rows)
}
