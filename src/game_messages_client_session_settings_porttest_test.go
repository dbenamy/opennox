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

func TestGameMessageClientSessionSettings(t *testing.T) {
	o := newServerOptionsOwner(t)
	words, restore := legacy.PortTestClientSessionWords()
	t.Cleanup(restore)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	serverConfigOwnBytes(t, 0x5D4594, 3512, 4)
	currentMap := serverConfigOwnBytes(t, 0x85B3FC, 36, 80)
	clear(currentMap)
	copy(currentMap, "arena")
	oldVersion := versionCode
	t.Cleanup(func() { versionCode = oldVersion })
	oldPlayState := legacy.GameGetPlayState
	legacy.GameGetPlayState = func() int { return 0 }
	t.Cleanup(func() { legacy.GameGetPlayState = oldPlayState })
	t.Cleanup(noxflags.PortTestGameFlags(0))
	oldExtra := legacy.PortTestServerConfigScalar("flags-get", 0, 0)
	t.Cleanup(func() { legacy.PortTestServerConfigScalar("flags-set", int(oldExtra), 0) })
	oldUpdated := legacy.PortTestServerConfigScalar("updated-get", 0, 0)
	t.Cleanup(func() {
		legacy.PortTestServerConfigScalar("updated-clear", 0, 0)
		if oldUpdated != 0 {
			legacy.PortTestServerConfigScalar("updated-set", 0, 0)
		}
	})
	if teamUICall("ctf-construct", 0, 0) != 1 || teamUICall("ball-construct", 0, 0) != 1 {
		t.Fatal("settings HUD owners")
	}
	type row struct {
		On, Variant                                             int
		OldFlags, NewFlags, OldGeneration, Generation           uint32
		Flags, StoredGeneration, Changed, Extra, Version, Limit uint32
		Settings                                                []byte
		CTFHidden, BallHidden                                   bool
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, oldFlags := range []uint32{0, 1, 32, 33, 64} {
			for _, newFlags := range []uint32{0, 32, 64, 128, 4096, 0x20000} {
				for _, oldGen := range []uint32{0, 100, 0xffffffff} {
					for _, generation := range []uint32{0, 100, 0x80000000, 0xffffffff} {
						for variant := 0; variant < 3; variant++ {
							noxflags.ResetGame()
							noxflags.SetGame(noxflags.GameFlag(oldFlags))
							binary.LittleEndian.PutUint32(connected, uint32(on))
							*memmap.PtrUint32(0x5D4594, 3512) = oldGen
							*words["settingsChanged"] = 99
							for i := range o.settings {
								o.settings[i] = 0xa5
							}
							binary.LittleEndian.PutUint16(o.settings[54:], 7)
							o.settings[56] = 8
							versionCode = 123
							legacy.PortTestServerConfigScalar("flags-set", 456, 0)
							legacy.PortTestServerConfigScalar("limit-set", 9, 0)
							legacy.PortTestServerConfigScalar("updated-clear", 0, 0)
							o.window("ctf").Hide()
							o.window("ball").Hide()
							data := make([]byte, 20)
							data[0] = 175
							binary.LittleEndian.PutUint32(data[1:], generation)
							binary.LittleEndian.PutUint32(data[5:], 0xabcdef01)
							binary.LittleEndian.PutUint32(data[9:], newFlags)
							binary.LittleEndian.PutUint32(data[13:], 0xaabbccdd)
							data[17] = 255
							data[18] = 7
							data[19] = 8
							if variant == 1 {
								data[18] = 255
							}
							if variant == 2 {
								data[19] = 255
							}
							input := bytes.Clone(data)
							want := bytes.Clone(o.settings)
							wantFlags := oldFlags
							wantGen := oldGen
							wantChanged := uint32(0)
							wantExtra, wantVersion, wantLimit := uint32(456), uint32(123), uint32(9)
							if generation > oldGen {
								wantGen = generation
								if oldFlags&1 == 0 {
									wantFlags = oldFlags&^524272 | newFlags
									wantExtra = 0xaabbccdd
									wantVersion = 0xabcdef01
									wantLimit = 255
									if data[18] != 7 || data[19] != 8 {
										wantChanged = 1
									}
									binary.LittleEndian.PutUint16(want[54:], uint16(data[18]))
									want[56] = data[19]
									if newFlags&128 == 0 {
										binary.LittleEndian.PutUint16(want[52:], uint16(newFlags))
										copy(want, "arena\x00")
									}
								}
							}
							n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(175), data)
							gotGen := memmap.Uint32(0x5D4594, 3512)
							extra := uint32(legacy.PortTestServerConfigScalar("flags-get", 0, 0))
							limit := uint32(legacy.PortTestServerConfigScalar("limit-get", 0, 0))
							if n != 20 || !bytes.Equal(input, data) || !bytes.Equal(o.settings, want) || uint32(noxflags.GetGame()) != wantFlags || gotGen != wantGen || *words["settingsChanged"] != wantChanged || extra != wantExtra || uint32(versionCode) != wantVersion || limit != wantLimit {
								t.Fatal("settings freshness/host/copy", on, oldFlags, newFlags, oldGen, generation, variant)
							}
							rows = append(rows, row{on, variant, oldFlags, newFlags, oldGen, generation, uint32(noxflags.GetGame()), gotGen, *words["settingsChanged"], extra, uint32(versionCode), limit, bytes.Clone(o.settings), o.window("ctf").GetFlags().IsHidden(), o.window("ball").GetFlags().IsHidden()})
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-settings", rows)
}
