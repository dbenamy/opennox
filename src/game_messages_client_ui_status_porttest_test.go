//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestGameMessageClientStatus(t *testing.T) {
	o := newInventoryDisplayOwner(t)
	meters := legacy.PortTestNewMeterEnvironment()
	t.Cleanup(meters.Restore)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	oldGUI, oldSetting := nox_client_renderGUI_80828, nox_xxx_xxxRenderGUI_587000_80832
	t.Cleanup(func() { nox_client_renderGUI_80828, nox_xxx_xxxRenderGUI_587000_80832 = oldGUI, oldSetting })
	oldEngine := noxflags.GetEngine()
	t.Cleanup(func() { noxflags.ResetEngine(); noxflags.SetEngine(oldEngine) })
	oldCode := legacy.ClientPlayerNetCode()
	t.Cleanup(func() { legacy.ClientSetPlayerNetCode(oldCode) })
	for i := range o.players {
		o.players[i].Active = 0
	}
	p := &o.players[0]
	p.NetCodeVal = 7
	type row struct {
		On, Host, Present, Graphics, GUI, Setting, Code int
		Initial, Incoming, Stored                       uint32
		Render                                          bool
		Meter                                           legacy.PortTestMeterRecord
		Mode                                            uint32
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for host := 0; host < 2; host++ {
			for present := 0; present < 2; present++ {
				for graphics := 0; graphics < 2; graphics++ {
					for prior := 0; prior < 2; prior++ {
						for setting := 0; setting < 2; setting++ {
							for _, code := range []uint16{7, 0x8007, 8, 65535} {
								for _, initial := range []uint32{0, 0x423, 0xffffffff} {
									for _, incoming := range []uint32{0, 1, 2, 0x20, 0x400, 0x423, 0xfffffbdd, 0xffffffff} {
										meters.Reset()
										legacy.ClientSetPlayerNetCode(7)
										meters.SetNamed(1, 0x1111)
										meters.SetNamed(2, 0x2222)
										meters.SetNamed(9, 0x5555)
										*memmap.PtrUint32(0x5D4594, 1091964) = 0x3333
										*memmap.PtrUint32(0x5D4594, 1092992) = 0x4444
										meters.Records[0].Color, meters.Records[0].Alternate = 0xaaaa, 0xbbbb
										nox_client_renderGUI_80828 = prior != 0
										nox_xxx_xxxRenderGUI_587000_80832 = setting != 0
										noxflags.ResetEngine()
										if graphics == 0 {
											noxflags.SetEngine(noxflags.EngineNoRendering)
										}
										reset := noxflags.PortTestGameFlags(noxflags.GameFlag(host))
										binary.LittleEndian.PutUint32(connected, uint32(on))
										p.Active = byte(present)
										p.Field3680 = initial
										data := []byte{106, byte(code), byte(code >> 8), 0, 0, 0, 0}
										binary.LittleEndian.PutUint32(data[3:], incoming)
										input := bytes.Clone(data)
										n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(106), data)
										reset()
										want := initial
										if present != 0 && code == 7 && host == 0 {
											want = (initial &^ 0x423) | (incoming & 0x423)
										}
										gui := prior != 0
										color, alternate, mode := uint32(0xaaaa), uint32(0xbbbb), uint32(0x5555)
										if present != 0 && code == 7 && graphics != 0 {
											if want&1 != 0 {
												gui = false
											} else if setting != 0 {
												gui = true
											}
											mode = (want >> 10) & 1
											if mode != 0 {
												color, alternate = 0x3333, 0x4444
											} else {
												color, alternate = 0x2222, 0x1111
											}
										}
										if n != 7 || !bytes.Equal(data, input) || p.Field3680 != want || nox_client_renderGUI_80828 != gui || meters.Records[0].Color != color || meters.Records[0].Alternate != alternate || meters.Named()[9] != mode {
											t.Fatal("client status mask/host/lookup/UI gates")
										}
										rows = append(rows, row{on, host, present, graphics, prior, setting, int(code), initial, incoming, p.Field3680, gui, meters.Records[0], mode})
									}
								}
							}
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-client-status", rows)
}
