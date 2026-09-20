//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"image"
	"testing"
	"unsafe"
)

func TestGameMessageClientSessionDeath(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	units, configure, _, free := c.srv.PortTestEscortPlayers()
	t.Cleanup(free)
	configure(2)
	pl, other := units[0].UpdateDataPlayer().Player, units[1].UpdateDataPlayer().Player
	oldPlayer, oldClient := legacy.Get_dword_8531A0_2576(), noxClient
	noxClient = c.Client
	t.Cleanup(func() { legacy.Set_dword_8531A0_2576(oldPlayer); noxClient = oldClient })
	windowWords, restore := legacy.PortTestInventoryWindowWords()
	t.Cleanup(restore)
	*windowWords["dword_5d4594_1319268"] = 0
	clear(serverConfigOwnBytes(t, 0x5D4594, 1049848, 4))
	c.GUI = gui.New(c.Render())
	parent := c.GUI.NewWindowRaw(nil, 8, 0, 0, 40, 40, nil)
	child := c.GUI.NewWindowRaw(parent, 8, 0, 0, 10, 10, nil)
	*windowWords["dword_5d4594_1062452"] = uint32(uintptr(parent.C()))
	t.Cleanup(func() { c.GUI.DestroyAll(); c.GUI.FreeDestroyed() })
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	oldDrag := legacy.Sub_478000
	calls := 0
	legacy.Sub_478000 = func() int { calls++; return oldDrag() }
	t.Cleanup(func() { legacy.Sub_478000 = oldDrag })
	normalize := func(v []server.EquipmentData) [][6]uint32 {
		out := make([][6]uint32, len(v))
		for i, e := range v {
			out[i][0], out[i][5] = e.Field0, e.Field20
			for j, p := range e.Field4 {
				if p != nil {
					if p != other.C() {
						t.Fatal("equipment pointer owner")
					}
					out[i][j+1] = 1
				}
			}
		}
		return out
	}
	type row struct {
		On, Quest, Present, Local, Code, Calls int
		Weapon, Armor                          uint32
		Weapons, Armors                        [][6]uint32
		Minimap                                byte
		Listed                                 bool
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for quest := 0; quest < 2; quest++ {
			for present := 0; present < 2; present++ {
				for local := 0; local < 2; local++ {
					for _, code := range []uint16{0, 17, 0x8011, 65535} {
						c.resetCase(env, pix, 1, 100)
						noxflags.ResetGame()
						if quest != 0 {
							noxflags.SetGame(4096)
						}
						binary.LittleEndian.PutUint32(connected, uint32(on))
						calls = 0
						child.Capture(true)
						pl.Active = byte(present)
						pl.NetCodeVal = 17
						pl.WeaponEquip, pl.ArmorEquip = 0xffffffff, 0xffffffff
						other.NetCodeVal = 999
						legacy.Set_dword_8531A0_2576(other)
						if local != 0 {
							legacy.Set_dword_8531A0_2576(pl)
						}
						for _, slots := range [][]server.EquipmentData{pl.Weapon[:], pl.Armor[:]} {
							for i := range slots {
								slots[i] = server.EquipmentData{Field0: 1 << uint(i), Field4: [4]unsafe.Pointer{other.C(), other.C(), other.C(), other.C()}, Field20: 0xcafe0000 + uint32(i)}
							}
						}
						want := *pl
						otherBefore := *other
						dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(300, 400))
						dr.ObjClass = 0
						dr.NetCode32 = 17
						c.Objs.MinimapAdd(dr, 3)
						active := on != 0 && present != 0 && code == 17
						if active && quest == 0 {
							want.WeaponEquip = 0
							// C starts the weapon sweep at byte2328 (four bytes into Weapon[0]).
							// It clears modifiers and the following record's mask, retaining tails.
							for i := range want.Weapon {
								want.Weapon[i].Field4 = [4]unsafe.Pointer{}
								if i != 0 {
									want.Weapon[i].Field0 = 0
								}
							}
							want.Armor[0].Field0 = 0
							for i := range want.Armor {
								if want.Armor[i].Field0&0xc0d == 0 {
									want.ArmorEquip &^= want.Armor[i].Field0
									want.Armor[i].Field0 = 0
									want.Armor[i].Field4 = [4]unsafe.Pointer{}
								}
							}
						}
						data := binary.LittleEndian.AppendUint16([]byte{232}, code)
						input := bytes.Clone(data)
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(232), data)
						wantCalls := 0
						if active && local != 0 {
							wantCalls = 1
						}
						wantMask := byte(3)
						if active {
							wantMask = 0
						}
						if n != 3 || !bytes.Equal(input, data) || *pl != want || *other != otherBefore || calls != wantCalls || (c.GUI.Captured() == child) != (wantCalls == 0) || dr.Field_71_0 != wantMask || (c.Objs.FirstMinimapList() == dr) != (wantMask != 0) {
							t.Fatal("death fields/gate/local UI", on, quest, present, local, code, calls)
						}
						rows = append(rows, row{on, quest, present, local, int(code), calls, pl.WeaponEquip, pl.ArmorEquip, normalize(pl.Weapon[:]), normalize(pl.Armor[:]), dr.Field_71_0, c.Objs.FirstMinimapList() == dr})
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-death", rows)
}
