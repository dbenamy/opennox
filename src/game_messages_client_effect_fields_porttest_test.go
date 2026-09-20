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

func TestGameMessageClientDeltaZ(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		On, Present, Static, Value, Return int
		Frame                              uint32
		Words                              [6]uint32
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for present := 0; present < 2; present++ {
			for static := 0; static < 2; static++ {
				for _, frame := range []uint32{0, 100, 0xffffffff} {
					for value := 0; value < 256; value++ {
						c.resetCase(env, pix, 1, frame)
						binary.LittleEndian.PutUint32(connected, uint32(on))
						dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(300, 400))
						if dr == nil {
							t.Fatal("drawable allocation")
						}
						dr.NetCode32 = 17
						dr.ObjClass = 0
						if static != 0 {
							dr.ObjClass = object.Class(0x20400000)
						}
						for i := 0; i < 6; i++ {
							*txword(dr, 432+uintptr(4*i)) = 0x12345678 + uint32(i)
						}
						raw := unsafe.Slice((*byte)(dr.C()), int(unsafe.Sizeof(*dr)))
						want := bytes.Clone(raw)
						code := uint16(17)
						if present == 0 {
							code = 18
						}
						if static != 0 {
							code |= 0x8000
						}
						data := []byte{159, byte(code), byte(code >> 8), byte(value), byte(255 - value), byte(value ^ 0x5a)}
						input := bytes.Clone(data)
						if on != 0 && present != 0 {
							binary.LittleEndian.PutUint32(want[432:], frame)
							binary.LittleEndian.PutUint32(want[436:], math.Float32bits(float32(data[3])))
							binary.LittleEndian.PutUint32(want[440:], math.Float32bits(float32(int8(data[4]))))
							binary.LittleEndian.PutUint32(want[444:], math.Float32bits(float32(data[5])))
						}
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(159), data)
						if n != 6 || !bytes.Equal(input, data) || !bytes.Equal(raw, want) {
							t.Fatalf("deltaZ on%d present%d static%d value%d frame%x", on, present, static, value, frame)
						}
						r := row{On: on, Present: present, Static: static, Value: value, Return: n, Frame: frame}
						for i := range r.Words {
							r.Words[i] = *txword(dr, 432+uintptr(i*4))
						}
						rows = append(rows, r)
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-delta-z", rows)
}

func TestGameMessageClientWhiteFlash(t *testing.T) {
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	state := serverConfigOwnBytes(t, 0x5D4594, 1096516, 12)
	type row struct {
		On, Return int
		Previous   uint32
		State      []byte
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, old := range []uint32{0, 1, 2, 0xffffffff} {
			binary.LittleEndian.PutUint32(connected, uint32(on))
			for i := range state {
				state[i] = 0x5a
			}
			binary.LittleEndian.PutUint32(state[4:], old)
			want := bytes.Clone(state)
			if on != 0 {
				binary.LittleEndian.PutUint32(want[4:], 1)
			}
			data := []byte{154, 1, 0xff, 0x80, 0x7f}
			input := bytes.Clone(data)
			n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(154), data)
			if n != 5 || !bytes.Equal(data, input) || !bytes.Equal(state, want) {
				t.Fatal("white flash connection gate/state")
			}
			rows = append(rows, row{on, n, old, bytes.Clone(state)})
		}
	}
	interactionCapture(t, "game-progress-white-flash", rows)
}
