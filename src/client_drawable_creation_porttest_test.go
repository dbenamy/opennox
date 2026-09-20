//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
)

func TestClientDrawableStreamCreation(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	c.srv.NetList = netlist.New()
	c.srv.NetList.Init()
	t.Cleanup(c.srv.NetList.Free)
	raw := memmap.Slice(0x5D4594, 1198020)[:255*8+8]
	saved := append([]byte(nil), raw...)
	t.Cleanup(func() { copy(raw, saved) })
	old := legacy.Nox_xxx_cliUpdateCameraPos_435600
	legacy.Nox_xxx_cliUpdateCameraPos_435600 = func(x, y int) {}
	t.Cleanup(func() { legacy.Nox_xxx_cliUpdateCameraPos_435600 = old })
	typ := c.Things.TypeByInd(4)
	savedClass := typ.ObjClass
	t.Cleanup(func() { typ.ObjClass = savedClass })
	type row struct {
		Stream                    int
		Class                     uint32
		Pos                       [2]int32
		Return                    int32
		Frame, Anim, Slave, Flags uint32
		Count                     int
	}
	var rows []row
	for _, stream := range []int{1, 2} {
		for _, cl := range []uint32{0, 0x200000, 0x200004, 0x200002, 0x20400000} {
			for _, point := range [][2]int32{{0, 0}, {1, 1}, {5999, 5999}, {6000, 6000}} {
				c.resetCase(env, pix, 1, 0xffffffff)
				clear(raw)
				typ.ObjClass = object.Class(cl)
				code := uint16(17)
				if cl&0x20400000 != 0 {
					code |= 0x8000
				}
				binary.LittleEndian.PutUint16(raw[8:], code)
				binary.LittleEndian.PutUint16(raw[10:], 4)
				packet := []byte{1}
				if stream == 2 {
					packet = []byte{0, 1}
				}
				packet = binary.LittleEndian.AppendUint16(packet, uint16(point[0]))
				packet = binary.LittleEndian.AppendUint16(packet, uint16(point[1]))
				ncoord := len(packet)
				packet = append(packet, 0x80, 253)
				if stream == 1 || cl&4 != 0 {
					packet = append(packet, 9)
				}
				ret, pos := legacy.PortTestDrawableStream(stream, packet, [2]int32{})
				dr := c.Objs.ByNetCode(code)
				wantRet := len(packet)
				if stream == 2 && cl&0x200000 == 0 {
					wantRet = ncoord
				}
				if dr == nil || c.Objs.Count != 1 || ret != int32(wantRet) || pos != point || dr.PosVec != image.Pt(int(point[0]), int(point[1])) {
					t.Fatalf("creation stream%d class%x point%v ret%d count%d", stream, cl, point, ret, c.Objs.Count)
				}
				rows = append(rows, row{stream, cl, pos, ret, dr.Field_72, dr.AnimInd, dr.AnimFrameSlave, uint32(dr.ObjFlags), c.Objs.Count})
			}
		}
	}
	for _, stream := range []int{1, 2} {
		for _, code := range []uint16{0, 17} {
			c.resetCase(env, pix, 1, 100)
			clear(raw)
			c.Objs.LoadError = false
			packet := []byte{0xff, byte(code), 0, 0, 0}
			if stream == 2 {
				packet = append([]byte{0}, packet...)
			}
			packet = binary.LittleEndian.AppendUint16(packet, 321)
			packet = binary.LittleEndian.AppendUint16(packet, 654)
			coordBytes := len(packet)
			packet = append(packet, 0, 7)
			ret, pos := legacy.PortTestDrawableStream(stream, packet, [2]int32{})
			want := int32(len(packet))
			if stream == 2 {
				want = -int32(coordBytes)
			}
			if ret != want || pos != [2]int32{321, 654} || c.Objs.Count != 0 || c.Objs.LoadError != (stream == 2 || code != 0) {
				t.Fatalf("missing type stream%d code%d ret%d loadError%v", stream, code, ret, c.Objs.LoadError)
			}
			rows = append(rows, row{Stream: stream, Pos: pos, Return: ret})
		}
	}
	drawableStateCapture(t, "creation", rows, "c0ec09e0e7d08f59d7e0d0e833c8ae5bf4432ce89a15d3d79245751fa9e38d44")
}
