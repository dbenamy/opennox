//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/opennox/libs/noxnet/netmsg"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientPickup(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	resetReliable, packets, freeReliable := legacy.PortTestVisibilityReliableOwner()
	t.Cleanup(freeReliable)
	names := make([]*byte, len(o.mods))
	for i, m := range o.mods {
		names[i] = *(**byte)(m.C())
	}
	t.Cleanup(o.c.srv.Server.PortTestControlsModifiers(o.mods, names))
	type row struct {
		On, Host, Code, Mode, Mods int
		Name                       string
		State                      inventoryTransactionResult
		Reliable                   []legacy.PortTestShopPacketResult
	}
	var rows []row
	for _, op := range []byte{75, 76} {
		for on := 0; on < 2; on++ {
			for host := 0; host < 2; host++ {
				for _, code := range []uint16{0, 7, 0x8007, 0xffff} {
					for _, name := range []string{"RedPotion", "Bow", "RedApple"} {
						for mode := 0; mode < 3; mode++ {
							for mi, ids := range [][4]byte{{255, 255, 255, 255}, {1, 2, 3, 4}, {4, 1, 0, 254}} {
								if op == 75 && mi != 0 {
									continue
								}
								o.reset(t)
								resetReliable()
								noxflags.ResetGame()
								noxflags.SetGame(noxflags.GameFlag(host))
								binary.LittleEndian.PutUint32(connected, uint32(on))
								typ := o.c.Things.TypeByID(name).Index()
								if mode == 1 {
									dr := o.item(t, name, 123)
									cell := &legacy.PortTestInventoryCells()[0]
									cell.Drawable, cell.Count = dr, 1
									cell.Codes[0] = 123
								}
								size := 5
								if op == 76 {
									size = 9
								}
								data := make([]byte, size)
								data[0] = op
								binary.LittleEndian.PutUint16(data[1:], code)
								binary.LittleEndian.PutUint16(data[3:], uint16(typ))
								if op == 76 {
									copy(data[5:], ids[:])
								}
								input := bytes.Clone(data)
								before := o.snapshot(t, 0, int(op), uint32(size))
								var n int
								run := func() { n = legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(op), data) }
								if mode == 2 {
									o.exhausted(t, run)
								} else {
									run()
								}
								if n != size || !bytes.Equal(data, input) {
									t.Fatalf("pickup %d return/input: %d", op, n)
								}
								cell := &legacy.PortTestInventoryCells()[0]
								wantCount := byte(0)
								if mode == 1 {
									wantCount = 1
								}
								if on != 0 && mode != 2 {
									wantCount++
								}
								if cell.Count != wantCount {
									t.Fatalf("pickup op%d on%d mode%d %s count%d want%d", op, on, mode, name, cell.Count, wantCount)
								}
								if on != 0 && mode != 2 && cell.Codes[cell.Count-1] != uint32(code&0x7fff) {
									t.Fatal("pickup masked inventory code")
								}
								if on != 0 && mode == 0 && uint32(cell.Drawable.ObjClass)&0x13001000 != 0 {
									for i := 0; i < 4; i++ {
										var want uint32
										if op == 76 && ids[i] >= 1 && ids[i] <= 4 {
											want = uint32(uintptr(o.mods[ids[i]-1].C()))
										}
										if *txword(cell.Drawable, 432+uintptr(4*i)) != want {
											t.Fatalf("pickup modifier %d", i)
										}
									}
									if *txword(cell.Drawable, 448) != 0xffffffff {
										t.Fatal("pickup modifier tail")
									}
								}
								state := o.snapshot(t, 0, int(op), uint32(n))
								ps := packets()
								if on != 0 && mode == 2 && host == 0 {
									want := []byte{241, byte(code & 0x7fff), byte((code & 0x7fff) >> 8)}
									if len(ps) != 1 || !bytes.Equal(ps[0].Data, want) || ps[0].Recipient != 31 {
										t.Fatalf("pickup failure report: %+v", ps)
									}
								} else if len(ps) != 0 {
									t.Fatalf("unexpected failure report: %+v", ps)
								}
								if on != 0 && mode == 2 && host != 0 {
									want := []byte{241, byte(code & 0x7fff), byte((code & 0x7fff) >> 8)}
									if len(state.Messages) != 1 || !bytes.Equal(state.Messages[0], want) {
										t.Fatalf("host pickup failure report: %v", state.Messages)
									}
								} else if len(state.Messages) != 0 {
									t.Fatalf("unexpected host report: %v", state.Messages)
								}
								if on == 0 {
									// Whole normalized state must remain untouched when disconnected.
									if !reflect.DeepEqual(before, state) {
										t.Fatal("disconnected pickup changed inventory")
									}
								}
								rows = append(rows, row{on, host, int(code), mode, mi, name, state, ps})
							}
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-pickup", rows)
}

func TestGameMessageClientInventoryNotifications(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		On, Code, Present, Count int
		Name                     string
		State                    inventoryTransactionResult
	}
	var rows []row
	for _, op := range []byte{77, 96, 97} {
		for on := 0; on < 2; on++ {
			for _, code := range []uint16{7, 0x8007, 8, 0xffff} {
				for present := 0; present < 2; present++ {
					for _, count := range []byte{1, 2, 32} {
						for _, name := range []string{"Bow", "Quiver", "RedApple"} {
							o.reset(t)
							binary.LittleEndian.PutUint32(connected, uint32(on))
							if present != 0 {
								dr := o.item(t, name, 7)
								cell := &legacy.PortTestInventoryCells()[0]
								cell.Drawable, cell.Count = dr, count
								cell.Codes[0] = 7
								for i := 1; i < int(count); i++ {
									cell.Codes[i] = uint32(100 + i)
								}
								*txword(dr, 448) = 0xffffffff
								if op == 97 {
									o.call(10, 7, 0, 0, 0)
								}
							}
							size := 3
							if op == 77 {
								size = 5
							}
							data := make([]byte, size)
							data[0] = op
							binary.LittleEndian.PutUint16(data[1:], code)
							if op == 77 {
								data[3], data[4] = 0xa5, 0x5a
							}
							before := o.snapshot(t, 0, int(op), uint32(size))
							input := bytes.Clone(data)
							n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(op), data)
							if n != size || !bytes.Equal(data, input) {
								t.Fatal("inventory notification return/input")
							}
							o.checkEquipment(t)
							after := o.snapshot(t, 0, int(op), uint32(n))
							if on == 0 && !reflect.DeepEqual(before, after) {
								t.Fatal("disconnected notification changed inventory")
							}
							admitted := on != 0 && present != 0 && code&0x7fff == 7
							total := 0
							for _, cell := range legacy.PortTestInventoryCells() {
								total += int(cell.Count)
							}
							want := 0
							if present != 0 {
								want = int(count)
							}
							if admitted && op == 77 {
								want--
							}
							if total != want {
								t.Fatalf("notification%d count%d want%d", op, total, want)
							}
							if admitted && op == 96 {
								wantEquipped := uint32(1)
								if name == "RedApple" {
									wantEquipped = 0
								}
								if legacy.PortTestInventoryCells()[0].Equipped != wantEquipped {
									t.Fatalf("equip flag %s", name)
								}
							}
							if admitted && op == 97 && legacy.PortTestInventoryCells()[0].Equipped != 0 {
								t.Fatal("unequip flag")
							}
							rows = append(rows, row{on, int(code), present, int(count), name, after})
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-inventory-notifications", rows)
}
