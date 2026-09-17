//go:build porttest

package opennox

import (
	"encoding/binary"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

type mapPolygonsOwner struct {
	words                                            map[string]*uint32
	vertices, polygons, selection, control, defaults []byte
}

func newMapPolygonsOwner(t *testing.T) *mapPolygonsOwner {
	t.Helper()
	o := new(mapPolygonsOwner)
	var restore func()
	o.words, restore = legacy.PortTestMapPolygonWords()
	t.Cleanup(restore)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	o.vertices = serverConfigOwnBytes(t, 0x5D4594, 535844, 16*1024)
	o.polygons = serverConfigOwnBytes(t, 0x5D4594, 552228, 140*256)
	o.selection = serverConfigOwnBytes(t, 0x5D4594, 534820, 4*256)
	o.control = serverConfigOwnBytes(t, 0x5D4594, 588072, 12)
	o.defaults = serverConfigOwnBytes(t, 0x587000, 60364, 128)
	for _, p := range [][]byte{o.vertices, o.polygons, o.selection, o.control, o.defaults} {
		clear(p)
	}
	for _, p := range o.words {
		*p = 0
	}
	copy(o.defaults, "Default Polygon")
	o.defaults[100] = 117
	o.defaults[101] = 38
	o.defaults[102] = 219
	o.defaults[126] = 9
	// Reset frees only allocations created inside this owner; old pointer-bearing
	// records and extracted globals are restored by the later cleanups.
	t.Cleanup(func() { *memmap.PtrUint32(0x5D4594, 588076) = 1; legacy.PortTestMapPolygonPointer("reset", nil, 0) })
	legacy.PortTestMapPolygonPointer("reset", nil, 0)
	return o
}
func (o *mapPolygonsOwner) vertex(index int) []byte  { return o.vertices[16*index : 16*(index+1)] }
func (o *mapPolygonsOwner) polygon(index int) []byte { return o.polygons[140*index : 140*(index+1)] }
func (o *mapPolygonsOwner) pointer(index int) unsafe.Pointer {
	return unsafe.Pointer(&o.polygon(index)[0])
}
func (o *mapPolygonsOwner) polygonID(p unsafe.Pointer) int {
	if p == nil {
		return 0
	}
	return int(binary.LittleEndian.Uint32(unsafe.Slice((*byte)(p), 140)[80:]))
}
func (o *mapPolygonsOwner) vertexID(p unsafe.Pointer) int {
	if p == nil {
		return 0
	}
	return int(*(*uint32)(p))
}

func (o *mapPolygonsOwner) construct(t *testing.T, points [][2]float32) unsafe.Pointer {
	t.Helper()
	if len(points) < 3 || len(points) > 256 {
		t.Fatal("fixture polygon vertex count")
	}
	index := int(*o.words["polygons"])
	for i, pt := range points {
		id := legacy.PortTestMapPolygonScalar("free-vertex", nil)
		legacy.PortTestMapPolygonVertex(pt[0], pt[1], id, 0)
		binary.LittleEndian.PutUint32(o.selection[4*i:], uint32(id))
	}
	binary.LittleEndian.PutUint16(o.control, uint16(len(points)))
	legacy.PortTestMapPolygonScalar("construct", nil)
	if *o.words["polygons"] != uint32(index+1) || binary.LittleEndian.Uint16(o.control) != 0 {
		t.Fatal("fixture construct")
	}
	return o.pointer(index)
}
