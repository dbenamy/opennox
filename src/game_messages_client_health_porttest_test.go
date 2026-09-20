//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientHealthChanges(t *testing.T) {
	o := newCombatOverlayOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		On, Count, Code, Amount, Return int
		Records                         []combatHealthRecord
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, count := range []int{0, 1, 31, 32} {
			for _, code := range []uint16{0, 7, 32767, 32768, 65535} {
				for _, amount := range []uint16{0, 1, 32767, 32768, 65535} {
					legacy.PortTestCombatHealthClear()
					binary.LittleEndian.PutUint32(connected, uint32(on))
					stamp := uint32(o.c.GetInputSeq())
					var want []combatHealthRecord
					for i := 0; i < count; i++ {
						legacy.PortTestCombatHealthAdd(uint32(i), int16(i))
						want = append([]combatHealthRecord{{uint32(i), int16(i), stamp}}, want...)
					}
					data := []byte{66, byte(code), byte(code >> 8), byte(amount), byte(amount >> 8)}
					before := bytes.Clone(data)
					if count < 32 {
						want = append([]combatHealthRecord{{uint32(code), int16(amount), stamp}}, want...)
					}
					n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(66), data)
					got := o.healthRecords(t)
					if n != 5 || !bytes.Equal(data, before) || !reflect.DeepEqual(got, want) {
						t.Fatalf("health delta on%d count%d code%x amount%x ret%d got%v want%v", on, count, code, amount, n, got, want)
					}
					rows = append(rows, row{on, count, int(code), int(amount), n, got})
					o.c.Inp.Tick()
				}
			}
		}
	}
	gameMessageCapture(t, "game-client-health-changes", rows)
}

func TestGameMessageClientHealthMeters(t *testing.T) {
	o := newMeterOwner(t)
	o.plain(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	staminaMax := serverConfigOwnBytes(t, 0x587000, 157092, 4)
	allies := serverConfigOwnBytes(t, 0x5D4594, 1200916, 512)
	records := append([]legacy.PortTestMeterRecord(nil), o.meters.Records...)
	player := unsafe.Slice((*byte)(unsafe.Pointer(&o.players[0])), int(unsafe.Sizeof(o.players[0])))
	oldPlayer := bytes.Clone(player)
	t.Cleanup(func() { copy(player, oldPlayer) })
	type row struct {
		On, Kind, Code, Ally, Player, Value, Maximum, Return int
		Meters                                               [7][2]uint32
		Duration, Timeout                                    uint32
		Allies                                               []byte
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, kind := range []int{65, 67, 69, 71, 221, 222} {
			codes := []uint16{7, 8, 0x8007, 0xffff}
			if kind == 67 || kind == 71 {
				codes = []uint16{0}
			}
			values := []uint16{0, 1, 255, 32767, 32768, 65535}
			if kind == 65 || kind == 71 {
				values = nil
				for v := 0; v < 256; v++ {
					values = append(values, uint16(v))
				}
			}
			maxes := []uint16{100}
			if kind == 221 || kind == 222 {
				maxes = []uint16{0, 1, 255, 32767, 32768, 65535}
			}
			for _, code := range codes {
				for _, value := range values {
					for _, maximum := range maxes {
						allyModes := []int{0}
						if kind == 65 || kind == 221 {
							allyModes = []int{0, 1}
						}
						playerModes := []int{1}
						if kind == 221 || kind == 222 {
							playerModes = []int{0, 1}
						}
						for _, ally := range allyModes {
							for _, hasPlayer := range playerModes {
								copy(o.meters.Records, records)
								copy(player, oldPlayer)
								clear(allies)
								clear(inputKeyTimeoutsOld)
								binary.LittleEndian.PutUint32(connected, uint32(on))
								binary.LittleEndian.PutUint32(staminaMax, uint32(maximum))
								*o.meters.NamedWord("nox_player_netCode_85319C") = 7
								*o.meters.NamedWord("dword_5d4594_1096260") = 99
								*o.meters.NamedWord("dword_8531A0_2576") = 0
								if hasPlayer != 0 {
									*o.meters.NamedWord("dword_8531A0_2576") = uint32(uintptr(unsafe.Pointer(&o.players[0])))
								}
								if ally != 0 {
									binary.LittleEndian.PutUint32(allies, uint32(code))
									binary.LittleEndian.PutUint16(allies[6:], 17)
									binary.LittleEndian.PutUint16(allies[8:], 101)
									binary.LittleEndian.PutUint32(allies[12:], 1)
								}
								wantAllies := bytes.Clone(allies)
								wantPlayer := bytes.Clone(player)
								wantRecords := append([]legacy.PortTestMeterRecord(nil), records...)
								data := []byte{byte(kind)}
								switch kind {
								case 65:
									data = append(data, byte(code), byte(code>>8), byte(value))
								case 67:
									data = binary.LittleEndian.AppendUint16(data, value)
								case 69:
									data = binary.LittleEndian.AppendUint16(data, code)
									data = binary.LittleEndian.AppendUint16(data, value)
								case 71:
									data = append(data, byte(value))
								default:
									data = binary.LittleEndian.AppendUint16(data, code)
									data = binary.LittleEndian.AppendUint16(data, value)
									data = binary.LittleEndian.AppendUint16(data, maximum)
								}
								before := bytes.Clone(data)
								wantDuration, wantTimeout := uint32(99), uint32(0)
								if on != 0 {
									switch kind {
									case 65:
										if ally != 0 {
											binary.LittleEndian.PutUint16(wantAllies[6:], 2*value)
										}
									case 67:
										wantRecords[0].Current = uint32(int32(int16(value)))
									case 69:
										if code&0x7fff == 7 {
											wantRecords[1].Current = uint32(value)
										}
									case 71:
										wantRecords[4].Current = uint32(value)
										wantRecords[4].Maximum = uint32(maximum)
										if value != maximum {
											wantTimeout = 120
										}
									case 221, 222:
										if code&0x7fff == 7 {
											i, off := 0, 2247
											if kind == 222 {
												i, off = 1, 2243
											}
											wantRecords[i].Current = uint32(value)
											wantRecords[i].Maximum = uint32(maximum)
											wantDuration = 32
											if hasPlayer != 0 {
												binary.LittleEndian.PutUint32(wantPlayer[off:], uint32(maximum))
											}
										} else if kind == 221 && ally != 0 {
											binary.LittleEndian.PutUint16(wantAllies[6:], value)
											binary.LittleEndian.PutUint16(wantAllies[8:], maximum)
										}
									}
								}
								n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(kind), data)
								if n != len(data) || !bytes.Equal(data, before) || !reflect.DeepEqual(o.meters.Records, wantRecords) || !bytes.Equal(player, wantPlayer) || !bytes.Equal(allies, wantAllies) || *o.meters.NamedWord("dword_5d4594_1096260") != wantDuration || inputKeyTimeoutsOld[17] != wantTimeout {
									t.Fatalf("meter on%d kind%d code%x ally%d player%d value%d max%d ret%d", on, kind, code, ally, hasPlayer, value, maximum, n)
								}
								r := row{On: on, Kind: kind, Code: int(code), Ally: ally, Player: hasPlayer, Value: int(value), Maximum: int(maximum), Return: n, Duration: wantDuration, Timeout: inputKeyTimeoutsOld[17], Allies: bytes.Clone(allies[:16])}
								for i, m := range o.meters.Records {
									r.Meters[i] = [2]uint32{m.Current, m.Maximum}
								}
								rows = append(rows, r)
							}
						}
					}
				}
			}
		}
	}
	gameMessageCapture(t, "game-client-health-meters", rows)
}
