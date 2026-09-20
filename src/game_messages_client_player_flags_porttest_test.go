//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientPlayerFlags(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	local := serverConfigOwnBytes(t, 0x852978, 8, 4)
	type row struct {
		On, Present, Return int
		Value, Flags        uint32
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for present := 0; present < 2; present++ {
			for _, value := range []uint32{0, 1, 4, 0x10000, 0x400000, 0x1000000, 0x20000000, 0x80000000, 0xffffffff} {
				clear(local)
				c.resetCase(env, pix, 1, 100)
				dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(300, 400))
				if dr == nil {
					t.Fatal("allocation")
				}
				originalFlags := dr.ObjFlags
				if present != 0 {
					binary.LittleEndian.PutUint32(local, uint32(uintptr(dr.C())))
				}
				binary.LittleEndian.PutUint32(connected, uint32(on))
				raw := unsafe.Slice((*byte)(dr.C()), int(unsafe.Sizeof(*dr)))
				want := bytes.Clone(raw)
				if on != 0 && present != 0 {
					binary.LittleEndian.PutUint32(want[120:], value)
				}
				data := binary.LittleEndian.AppendUint32([]byte{102}, value)
				before := bytes.Clone(data)
				n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(102), data)
				if n != 5 || !bytes.Equal(data, before) || !bytes.Equal(raw, want) {
					t.Fatalf("player flags on%d present%d value%x return%d", on, present, value, n)
				}
				rows = append(rows, row{on, present, n, value, uint32(dr.ObjFlags)})
				// Arbitrary wire flags are data under test, not fixture list ownership.
				dr.ObjFlags = originalFlags
			}
		}
	}
	clear(local)
	gameMessageCapture(t, "game-client-player-flags", rows)
}
