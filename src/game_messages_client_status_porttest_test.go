//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientExtendedStatus(t *testing.T) {
	o := newObjectDrawingOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		On, Present, Static, Generator, Count, Delay, Return int
		Previous, Value, Class, Flags, Countdown             uint32
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for present := 0; present < 2; present++ {
			for static := 0; static < 2; static++ {
				for gen := 0; gen < 2; gen++ {
					for _, old := range []uint32{0, 0x400, 0x800, 0xffffffff} {
						for _, value := range []uint32{0, 0x400, 0x800, 0xc00, 0x80000000, 0xffffffff} {
							for _, cd := range [][2]int{{0, 0}, {1, 0}, {7, 5}, {255, 255}} {
								o.reset(1, 100)
								o.generatorData(cd[0], cd[1], 1)
								binary.LittleEndian.PutUint32(connected, uint32(on))
								dr := o.c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(300, 400))
								if dr == nil {
									t.Fatal("drawable allocation")
								}
								originalFlags := dr.ObjFlags
								dr.ObjClass = 0x80000
								if static != 0 {
									dr.ObjClass |= object.Class(0x20400000)
								}
								if gen != 0 {
									dr.ObjClass |= 0x20000
								}
								dr.ObjFlags |= 0x20000000
								dr.NetCode32 = 17
								dr.Flags70Val = old
								dr.DrawData = unsafe.Pointer(&o.data[0])
								*txword(dr, 432) = 0x12345678
								raw := unsafe.Slice((*byte)(dr.C()), int(unsafe.Sizeof(*dr)))
								want := bytes.Clone(raw)
								code := uint16(17)
								if present == 0 {
									code = 18
								}
								if static != 0 {
									code |= 0x8000
								}
								data := binary.LittleEndian.AppendUint32([]byte{101, byte(code), byte(code >> 8)}, value)
								input := bytes.Clone(data)
								if on != 0 && present != 0 {
									binary.LittleEndian.PutUint32(want[280:], value)
									if gen != 0 {
										if old&0x400 == 0 && value&0x400 != 0 {
											binary.LittleEndian.PutUint32(want[432:], uint32(cd[0]*(cd[1]+1)))
										}
										if value&0x800 != 0 {
											binary.LittleEndian.PutUint32(want[112:], uint32(dr.ObjClass)&^0x80000)
											binary.LittleEndian.PutUint32(want[120:], uint32(dr.ObjFlags)&^0x20000000)
										}
									}
								}
								n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(101), data)
								if n != 7 || !bytes.Equal(data, input) || !bytes.Equal(raw, want) {
									t.Fatalf("status on%d present%d static%d gen%d old%x value%x count/delay%v", on, present, static, gen, old, value, cd)
								}
								rows = append(rows, row{on, present, static, gen, cd[0], cd[1], n, old, value, uint32(dr.ObjClass), uint32(dr.ObjFlags), *txword(dr, 432)})
								dr.ObjFlags = originalFlags
							}
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-extended-status", rows)
}
