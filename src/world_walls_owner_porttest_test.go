//go:build porttest

package opennox

import (
	"encoding/binary"
	noxcolor "github.com/opennox/libs/color"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"runtime"
	"testing"
	"unsafe"
)

type worldWallOwner struct {
	*objectRenderOwner
	words      map[string]*uint32
	regions    [][]byte
	wallImages []*noxrender.Image
}

func worldWallImage(frame int) []byte {
	const w, h = 12, 8
	b := make([]byte, 17)
	binary.LittleEndian.PutUint32(b, w)
	binary.LittleEndian.PutUint32(b[4:], h)
	for y := 0; y < h; y++ {
		b = append(b, 2, w)
		for x := 0; x < w; x++ {
			color := noxcolor.RGB5551Color(byte(40+x*13+frame*20), byte(50+y*19), byte(200-frame*35)).Color16()
			b = binary.LittleEndian.AppendUint16(b, color)
		}
	}
	return b
}
func newWorldWallOwner(t *testing.T) *worldWallOwner {
	base := newObjectRenderOwner(t)
	o := &worldWallOwner{objectRenderOwner: base}
	var restore func()
	o.words, restore = legacy.PortTestWorldWallWords()
	t.Cleanup(restore)
	t.Cleanup(o.c.srv.PortTestMinimapWalls())
	for _, reg := range [][3]uintptr{{0x587000, 80808, 4}, {0x587000, 80816, 4}, {0x5D4594, 805848, 4}, {0x973F18, 12, 4}, {0x973F18, 52, 4}} {
		buf := unsafe.Slice((*byte)(memmap.PtrOff(reg[0], reg[1])), reg[2])
		old := append([]byte(nil), buf...)
		o.regions = append(o.regions, buf)
		t.Cleanup(func() { copy(buf, old) })
	}
	for _, reg := range blobdata.PortTestWorldWallTables() {
		buf := unsafe.Slice((*byte)(memmap.PtrOff(reg.Base, reg.Offset)), len(reg.Data))
		old := append([]byte(nil), buf...)
		copy(buf, reg.Data)
		t.Cleanup(func() { copy(buf, old) })
	}
	// C retains pixel addresses in its actual row table until cleanup.
	var pixels runtime.Pinner
	pixels.Pin(&o.pix.Pix[0])
	t.Cleanup(pixels.Unpin)
	oldBuffer := noxPixBuffer
	noxPixBuffer.img = o.pix
	noxPixBuffer.rows = nil
	noxPixBuffer.freeRows = nil
	nox_video_initPixbufferRows_486230()
	ownedRows := noxPixBuffer.freeRows
	t.Cleanup(func() { legacy.Set_nox_pixbuffer_rows_3798784(nil); ownedRows(); noxPixBuffer = oldBuffer })
	raw := [][]byte{worldWallImage(0), worldWallImage(1), worldWallImage(2), worldWallImage(3)}
	o.wallImages, restore = o.c.r.GetBag().PortTestWorldWallImages(raw)
	t.Cleanup(restore)
	for i, img := range o.wallImages {
		o.c.imageRefs[uint32(uintptr(img.C()))] = 0xeb000000 + uint32(i)
	}
	def := o.c.srv.Walls.DefByInd(0)
	for layer := 0; layer < 4; layer++ {
		for dir := 0; dir < 15; dir++ {
			for variant := 0; variant < 16; variant++ {
				def.Sprite8432[layer][dir][variant] = unsafe.Pointer(o.wallImages[(dir+layer+variant)%len(o.wallImages)].C())
				def.DrawOffs752[layer][dir][variant] = image.Pt(50+variant%3, -72+layer)
			}
		}
	}
	o.resetWalls(t)
	return o
}
func (o *worldWallOwner) resetWalls(t *testing.T) {
	o.resetRender(1, 120)
	o.c.srv.Walls.Reset()
	for _, b := range o.regions {
		clear(b)
	}
	for _, p := range o.words {
		*p = 0
	}
	*o.words["edgeMaxX"] = 96
	*o.words["edgeMaxY"] = 96
	o.c.srv.Walls.DefByInd(0).Flags32 = 0
	o.c.Viewport().World = image.Rectangle{Min: image.Pt(190, 190), Max: image.Pt(286, 286)}
}
