//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"math"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientEnchantments(t *testing.T) {
	c, pix, effects := newEffectsFullOwner(t)
	env := legacy.PortTestNewClientParticleEnvironment()
	t.Cleanup(env.Restore)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	local := serverConfigOwnBytes(t, 0x852978, 8, 4)
	ui := serverConfigOwnBytes(t, 0x5D4594, 1062536, 8)
	typ := c.Things.TypeByInd(4)
	oldIntensity := typ.LightIntensity
	typ.LightIntensity = 3.25
	t.Cleanup(func() { typ.LightIntensity = oldIntensity })
	type row struct {
		On, Kind, Present, Static, Local int
		Old, Value, Item                 uint32
		Return                           int
		Buffs, Intensity, Radius, Fixed  uint32
		UI                               []byte
	}
	var rows []row
	run := func(on, kind, present, static, isLocal int, old, value, item uint32) {
		clear(local)
		c.resetCase(effects, pix, 1, 100)
		env.Reset()
		binary.LittleEndian.PutUint32(connected, uint32(on))
		binary.LittleEndian.PutUint32(ui, item|0xa5b6c700)
		binary.LittleEndian.PutUint32(ui[4:], 0x12345678)
		dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(300, 400))
		if dr == nil {
			t.Fatal("allocation")
		}
		dr.NetCode32 = 17
		dr.ObjClass = 0
		if static != 0 {
			dr.ObjClass = object.Class(0x20400000)
		}
		dr.Buffs = old
		dr.LightIntensity = 17
		dr.LightIntensityRad = 0x12345678
		dr.LightIntensityU16 = 0x87654321
		if isLocal != 0 {
			binary.LittleEndian.PutUint32(local, uint32(uintptr(dr.C())))
		}
		raw := unsafe.Slice((*byte)(dr.C()), int(unsafe.Sizeof(*dr)))
		before := bytes.Clone(raw)
		want := bytes.Clone(raw)
		wantUI := bytes.Clone(ui)
		code := uint16(17)
		if present == 0 {
			code = 18
		}
		if static != 0 {
			code |= 0x8000
		}
		data := []byte{byte(kind), byte(value)}
		restore := false
		if kind == 90 {
			data = binary.LittleEndian.AppendUint32([]byte{90, byte(code), byte(code >> 8)}, value)
			if on != 0 && present != 0 {
				binary.LittleEndian.PutUint32(want[124:], value)
				if isLocal != 0 {
					binary.LittleEndian.PutUint32(wantUI[4:], value)
				}
				restore = old&(1<<15) != 0 && value&(1<<15) == 0 && (isLocal == 0 || item&8 == 0)
			}
		} else if on != 0 {
			wantUI[0] = byte(value)
			restore = item&8 != 0 && value&8 == 0 && isLocal != 0 && old&(1<<15) == 0
		}
		if restore {
			legacy.PortTestClientDrawParticle(16, c.Viewport(), dr, [4]int32{int32(math.Float32bits(3.25))})
			copy(want[140:152], raw[140:152])
			copy(raw, before)
			if binary.LittleEndian.Uint32(want[140:]) != math.Float32bits(3.25) || binary.LittleEndian.Uint32(want[144:]) == 0 || binary.LittleEndian.Uint32(want[148:]) != 212992 {
				t.Fatal("nontrivial normal-light restoration")
			}
		}
		input := bytes.Clone(data)
		n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(kind), data)
		if n != len(data) || !bytes.Equal(data, input) || !bytes.Equal(raw, want) || !bytes.Equal(ui, wantUI) {
			t.Fatalf("enchantment on%d kind%d present%d static%d local%d old%x value%x item%x return%d", on, kind, present, static, isLocal, old, value, item, n)
		}
		rows = append(rows, row{on, kind, present, static, isLocal, old, value, item, n, dr.Buffs, math.Float32bits(dr.LightIntensity), dr.LightIntensityRad, dr.LightIntensityU16, bytes.Clone(ui)})
	}
	for on := 0; on < 2; on++ {
		for present := 0; present < 2; present++ {
			for static := 0; static < 2; static++ {
				for isLocal := 0; isLocal < 2; isLocal++ {
					for _, old := range []uint32{0, 1 << 15, 0xffffffff} {
						for _, value := range []uint32{0, 1, 1 << 15, 0xffffffff} {
							for _, item := range []uint32{0, 8} {
								run(on, 90, present, static, isLocal, old, value, item)
							}
						}
					}
				}
			}
		}
	}
	for on := 0; on < 2; on++ {
		for isLocal := 0; isLocal < 2; isLocal++ {
			for _, old := range []uint32{0, 1 << 15} {
				for _, item := range []uint32{0, 8, 255} {
					for value := uint32(0); value < 256; value++ {
						run(on, 91, 1, 0, isLocal, old, value, item)
					}
				}
			}
		}
	}
	clear(local)
	interactionCapture(t, "game-client-enchantments", rows)
}
