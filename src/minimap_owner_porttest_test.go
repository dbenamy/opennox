//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"math"
	"testing"
	"unsafe"
)

type minimapOwner struct {
	*objectRenderOwner
	words           map[string]*uint32
	regions         [][]byte
	polygonVertices [][]uint32
}

func newMinimapOwner(t *testing.T) *minimapOwner {
	base := newObjectRenderOwner(t, "Crown", "GameBall")
	o := &minimapOwner{objectRenderOwner: base}
	var restore func()
	o.words, restore = legacy.PortTestMinimapWords()
	t.Cleanup(restore)
	t.Cleanup(o.c.srv.PortTestMinimapWalls())
	t.Cleanup(o.c.srv.PortTestMinimapTeamColors())
	for _, r := range [][3]uintptr{{0x5D4594, 1096300, 20}, {0x5D4594, 1096424, 4}, {0x5D4594, 535844, 16 * 12}, {0x5D4594, 552228, 140 * 4}, {0x5D4594, 588080, 4}, {0x852978, 8, 4}, {0x5D4594, 823804, 1932}} {
		b := unsafe.Slice((*byte)(memmap.PtrOff(r[0], r[1])), r[2])
		old := append([]byte(nil), b...)
		o.regions = append(o.regions, b)
		t.Cleanup(func() { copy(b, old) })
	}
	for j := 0; j < 2; j++ {
		b, free := alloc.Make([]uint32{}, 4)
		o.polygonVertices = append(o.polygonVertices, b)
		t.Cleanup(free)
	}
	for _, r := range blobdata.PortTestMinimapTables() {
		b := unsafe.Slice((*byte)(memmap.PtrOff(r.Base, r.Offset)), len(r.Data))
		old := append([]byte(nil), b...)
		copy(b, r.Data)
		t.Cleanup(func() { copy(b, old) })
	}
	o.resetMinimap(t)
	return o
}
func (o *minimapOwner) resetMinimap(t *testing.T) {
	o.resetRender(1, 120)
	o.c.srv.Walls.Reset()
	for _, b := range o.regions {
		clear(b)
	}
	*o.words["messageHead"] = 0
	*o.words["debugIterator"] = 0
	*o.words["zoom"] = 500
	*o.words["polygons"] = 1
	*o.words["gui"] = 0
	*o.words["width"] = 96
	*o.words["height"] = 96
	// Use real named type lookup/cache and real wall definitions on first draw.
	o.c.Viewport().World = image.Rect(182, 182, 230, 230)
}
func (o *minimapOwner) polygon(t *testing.T, index int, bounds image.Rectangle, level byte) {
	t.Helper()
	p := legacy.PortTestMinimapNewPolygon()
	if p == nil {
		t.Fatal("production polygon allocation")
	}
	*(*int32)(unsafe.Add(p, 88)) = int32(bounds.Min.X)
	*(*int32)(unsafe.Add(p, 92)) = int32(bounds.Min.Y)
	*(*int32)(unsafe.Add(p, 96)) = int32(bounds.Max.X)
	*(*int32)(unsafe.Add(p, 100)) = int32(bounds.Max.Y)
	vertices := o.polygonVertices[index]
	points := []image.Point{bounds.Min, image.Pt(bounds.Max.X, bounds.Min.Y), bounds.Max, image.Pt(bounds.Min.X, bounds.Max.Y)}
	for i, pt := range points {
		vertex := 1 + index*4 + i
		vertices[i] = uint32(vertex)
		off := uintptr(535844 + 16*vertex)
		*memmap.PtrUint32(0x5D4594, off) = uint32(vertex)
		*memmap.PtrUint32(0x5D4594, off+4) = math.Float32bits(float32(pt.X))
		*memmap.PtrUint32(0x5D4594, off+8) = math.Float32bits(float32(pt.Y))
		*memmap.PtrUint32(0x5D4594, off+12) = 1
	}
	*(*unsafe.Pointer)(unsafe.Add(p, 108)) = unsafe.Pointer(&vertices[0])
	*(*uint16)(unsafe.Add(p, 128)) = 4
	*(*byte)(unsafe.Add(p, 130)) = level
}
