//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"reflect"
	"testing"
	"unsafe"
)

func TestGameMessageClientSessionRespawn(t *testing.T) {
	o := newPlayerStateEquipmentOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	pl := o.units[0].UpdateDataPlayer().Player
	type row struct {
		On, Present, Class, Mask, Flag int
		Mode                           uint32
		Code                           uint16
		State                          playerStateEquipmentSnapshot
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for present := 0; present < 2; present++ {
			for class := 0; class < 3; class++ {
				for _, mode := range []uint32{0, 1, 2048, 4096} {
					for _, code := range []uint16{17, 0x8011, 65535} {
						for _, mask := range []byte{0, 1, 31, 128, 255} {
							for _, flag := range []byte{0, 1, 2, 255} {
								noxflags.ResetGame()
								noxflags.SetGame(noxflags.GameFlag(mode))
								binary.LittleEndian.PutUint32(connected, uint32(on))
								pl.Active = byte(present)
								pl.NetCodeVal = 17
								pl.WeaponEquip, pl.ArmorEquip = 0x40000000, 0x20000000
								*(*byte)(unsafe.Add(pl.C(), 2251)) = byte(class)
								for i := 0; i < 5; i++ {
									*(*byte)(unsafe.Add(pl.C(), 2268+i)) = byte(i + 1)
								}
								for i := range pl.Weapon {
									pl.Weapon[i] = server.EquipmentData{Field0: 0x123, Field4: [4]unsafe.Pointer{unsafe.Pointer(o.s.Modif.Nox_xxx_modifGetDescById413330(1)), nil, nil, nil}, Field20: uint32(100 + i)}
								}
								for i := range pl.Armor {
									pl.Armor[i] = server.EquipmentData{Field0: 0x456, Field4: [4]unsafe.Pointer{nil, unsafe.Pointer(o.s.Modif.Nox_xxx_modifGetDescById413330(16)), nil, nil}, Field20: uint32(200 + i)}
								}
								before := *pl
								if on != 0 && present != 0 && code == 17 && flag&1 != 0 {
									legacy.PortTestPlayerStateEquipment("respawn", pl, mask, 0, 0, nil)
								}
								want := *pl
								wantSnapshot := o.snapshot(t, pl)
								*pl = before
								data := binary.LittleEndian.AppendUint16([]byte{233}, code)
								data = binary.LittleEndian.AppendUint32(data, 0xffffffff)
								data = append(data, mask, flag)
								input := bytes.Clone(data)
								n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(233), data)
								got := o.snapshot(t, pl)
								if n != 9 || !bytes.Equal(data, input) || *pl != want || !reflect.DeepEqual(got, wantSnapshot) {
									t.Fatal("respawn connection/player/mask/flag", on, present, class, mode, code, mask, flag)
								}
								rows = append(rows, row{on, present, class, int(mask), int(flag), mode, code, got})
							}
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-respawn", rows)
}
