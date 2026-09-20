//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"slices"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/strman"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientLessons(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	o.reset(t)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	printer, index, region := teamRuntimeJoinTextOwner(t, o.c.srv.Server)
	set, restore := o.c.srv.PortTestMeterStrings(
		strman.Entry{ID: "guimsg.c:systemmsg", Vals: []strman.Variant{{Str: "System: %s"}}},
		strman.Entry{ID: "cdecode.c:Eliminated", Vals: []strman.Variant{{Str: "Out %s"}}},
	)
	t.Cleanup(restore)
	set(0)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	config := serverConfigOwnBytes(t, 0x5D4594, 371380, 58)
	for i := range o.players {
		o.players[i].Active = 0
	}
	pl := &o.players[0]
	pl.NetCodeVal = 7
	pl.SetName("Ada Ω")
	player := unsafe.Slice((*byte)(unsafe.Pointer(pl)), int(unsafe.Sizeof(*pl)))
	old := bytes.Clone(player)
	t.Cleanup(func() { copy(player, old) })
	type row struct {
		On, Host, Elimination, Present, Code, Limit, Return int
		Frame, Wins, Losses                                 uint32
		Stored                                              [3]uint32
		Text                                                []string
		Sounds                                              [][2]int
		Head                                                uint32
		Ring                                                []byte
	}
	var rows []row
	scores := [][2]uint32{{0, 0}, {1, 1}, {255, 255}, {0x7fffffff, 65535}, {0x80000000, 0x80000000}, {0xffffffff, 0xffffffff}}
	for on := 0; on < 2; on++ {
		for host := 0; host < 2; host++ {
			for elimination := 0; elimination < 2; elimination++ {
				for present := 0; present < 2; present++ {
					for _, code := range []uint16{7, 0x8007, 8, 0xffff} {
						for _, score := range scores {
							for _, limit := range []uint16{0, 1, 65535} {
								for _, frame := range []uint32{0, 100, 0xffffffff} {
									noxflags.ResetGame()
									noxflags.SetGame(noxflags.GameFlag(host | elimination*1024))
									binary.LittleEndian.PutUint32(connected, uint32(on))
									binary.LittleEndian.PutUint16(config[54:], limit)
									o.c.srv.SetFrame(frame)
									pl.Active = byte(present)
									for i, v := range []uint32{17, 29, 31} {
										binary.LittleEndian.PutUint32(player[2136+4*i:], v)
									}
									want := bytes.Clone(player)
									data := []byte{78, byte(code), byte(code >> 8), 0, 0, 0, 0, 0, 0, 0, 0}
									binary.LittleEndian.PutUint32(data[3:], score[0])
									binary.LittleEndian.PutUint32(data[7:], score[1])
									input := bytes.Clone(data)
									found := on != 0 && present != 0 && code&0x7fff == 7
									if found && host == 0 {
										binary.LittleEndian.PutUint32(want[2136:], score[0])
										binary.LittleEndian.PutUint32(want[2140:], score[1])
										binary.LittleEndian.PutUint32(want[2144:], frame)
									}
									var text []string
									var sounds [][2]int
									head := uint32(2)
									// The old dispatcher compares its signed int temporary against the
									// promoted unsigned-short limit, independently of the stored score.
									if found && elimination != 0 && int32(score[1]) >= int32(limit) {
										text = []string{"System: Out Ada Ω"}
										sounds = [][2]int{{312, 100}}
										head = 0
									}
									clear(region)
									*index = 2
									printer.lines = nil
									o.sounds = nil
									n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(78), data)
									if n != 11 || !bytes.Equal(data, input) || !bytes.Equal(player, want) || !slices.Equal(printer.lines, text) || !slices.Equal(o.sounds, sounds) || *index != head {
										t.Fatalf("lessons on%d host%d elimination%d present%d code%x score%x limit%d frame%x: return%d text%v sounds%v head%d playerEqual%v", on, host, elimination, present, code, score, limit, frame, n, printer.lines, o.sounds, *index, bytes.Equal(player, want))
									}
									rows = append(rows, row{On: on, Host: host, Elimination: elimination, Present: present, Code: int(code), Limit: int(limit), Return: n, Frame: frame, Wins: score[0], Losses: score[1], Stored: [3]uint32{binary.LittleEndian.Uint32(player[2136:]), binary.LittleEndian.Uint32(player[2140:]), binary.LittleEndian.Uint32(player[2144:])}, Text: slices.Clone(printer.lines), Sounds: slices.Clone(o.sounds), Head: *index, Ring: bytes.Clone(region)})
								}
							}
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-lessons", rows)
}

func TestGameMessageClientExperience(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	o.reset(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	storage := serverConfigOwnBytes(t, 0x5D4594, 1062540, 12)
	type row struct {
		On      int
		Value   uint32
		Return  int
		Storage []byte
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, value := range []uint32{0, 1, 1024, 0x7fffffff, 0x80000000, 0xffffffff} {
			for i := range storage {
				storage[i] = byte(0xa0 + i)
			}
			binary.LittleEndian.PutUint32(connected, uint32(on))
			want := bytes.Clone(storage)
			if on != 0 {
				binary.LittleEndian.PutUint32(want[4:], value)
			}
			data := []byte{110, 0, 0, 0, 0}
			binary.LittleEndian.PutUint32(data[1:], value)
			input := bytes.Clone(data)
			n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(110), data)
			if n != 5 || !bytes.Equal(data, input) || !bytes.Equal(storage, want) {
				t.Fatal("experience", on, value, n, storage, want)
			}
			rows = append(rows, row{on, value, n, bytes.Clone(storage)})
		}
	}
	interactionCapture(t, "game-progress-experience", rows)
}

func TestGameMessageClientTreasure(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	o.reset(t)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	printer, index, region := teamRuntimeJoinTextOwner(t, o.c.srv.Server)
	set, restore := o.c.srv.PortTestMeterStrings(
		strman.Entry{ID: "guimsg.c:systemmsg", Vals: []strman.Variant{{Str: "System: %s"}}},
		strman.Entry{ID: "cdecode.c:SH_NearVictory", Vals: []strman.Variant{{Str: "Near %s"}}},
	)
	t.Cleanup(restore)
	set(0)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	for i := range o.players {
		o.players[i].Active = 0
	}
	pl := &o.players[0]
	pl.NetCodeVal = 7
	pl.SetName("Ada Ω")
	player := unsafe.Slice((*byte)(unsafe.Pointer(pl)), int(unsafe.Sizeof(*pl)))
	old := bytes.Clone(player)
	t.Cleanup(func() { copy(player, old) })
	type row struct {
		On, Host, Present, Code, Return int
		Frame, PreviousFrame            uint32
		Before                          [2]uint32
		Incoming                        [2]uint16
		Stored                          [3]uint32
		Text                            []string
		Head                            uint32
		Ring                            []byte
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for host := 0; host < 2; host++ {
			for present := 0; present < 2; present++ {
				for _, code := range []uint16{7, 0x8007, 8, 0xffff} {
					for _, previous := range [][2]uint32{{0, 1}, {5, 6}, {0xffffffff, 0}, {9, 9}} {
						for _, incoming := range [][2]uint16{{0, 1}, {1, 1}, {65534, 65535}, {65535, 0}} {
							for _, frames := range [][2]uint32{{0, 0}, {1, 0}, {100, 100}, {100, 101}, {0xffffffff, 0}, {0, 0xffffffff}, {0xffffffff, 0xffffffff}} {
								noxflags.ResetGame()
								noxflags.SetGame(noxflags.GameFlag(host))
								binary.LittleEndian.PutUint32(connected, uint32(on))
								o.c.srv.SetFrame(frames[0])
								pl.Active = byte(present)
								for i, v := range []uint32{previous[0], previous[1], frames[1]} {
									binary.LittleEndian.PutUint32(player[2152+4*i:], v)
								}
								want := bytes.Clone(player)
								data := []byte{85, byte(code), byte(code >> 8), 0, 0, 0, 0}
								binary.LittleEndian.PutUint16(data[3:], incoming[0])
								binary.LittleEndian.PutUint16(data[5:], incoming[1])
								input := bytes.Clone(data)
								found := on != 0 && present != 0 && code&0x7fff == 7
								if found && host == 0 && frames[0] > frames[1] {
									binary.LittleEndian.PutUint32(want[2152:], uint32(incoming[0]))
									binary.LittleEndian.PutUint32(want[2156:], uint32(incoming[1]))
									binary.LittleEndian.PutUint32(want[2160:], frames[0])
								}
								var text []string
								head := uint32(2)
								// Announcements use the resulting stored values even when the host
								// or an equal/older frame prevented a record update. Subtraction wraps.
								if found && binary.LittleEndian.Uint32(want[2152:]) == binary.LittleEndian.Uint32(want[2156:])-1 {
									text = []string{"System: Near Ada Ω"}
									head = 0
								}
								clear(region)
								*index = 2
								printer.lines = nil
								o.sounds = nil
								n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(85), data)
								if n != 7 || !bytes.Equal(data, input) || !bytes.Equal(player, want) || !slices.Equal(printer.lines, text) || len(o.sounds) != 0 || *index != head {
									t.Fatalf("treasure on%d host%d present%d code%x old%v incoming%v frames%x: return%d text%v head%d playerEqual%v", on, host, present, code, previous, incoming, frames, n, printer.lines, *index, bytes.Equal(player, want))
								}
								rows = append(rows, row{On: on, Host: host, Present: present, Code: int(code), Return: n, Frame: frames[0], PreviousFrame: frames[1], Before: previous, Incoming: incoming, Stored: [3]uint32{binary.LittleEndian.Uint32(player[2152:]), binary.LittleEndian.Uint32(player[2156:]), binary.LittleEndian.Uint32(player[2160:])}, Text: slices.Clone(printer.lines), Head: *index, Ring: bytes.Clone(region)})
							}
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-treasure", rows)
}
