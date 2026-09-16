//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"github.com/opennox/opennox/v1/server"
	"os"
	"testing"
	"unsafe"
)

type floorAssetsOwner struct {
	defs             []server.TileDef
	edges            []byte
	count, edgeCount *uint32
	refs             map[uint32]uint32
	scratch          []byte
}

func newFloorAssetsOwner(t *testing.T) *floorAssetsOwner {
	t.Helper()
	t.Cleanup(handles.PortTestInit())
	facade := blobdata.PortTestFloorFacadeData()
	table := unsafe.Slice(memmap.PtrUint8(0x587000, 26488), len(facade))
	saved := append([]byte(nil), table...)
	copy(table, facade)
	for i, off := range []uintptr{26536, 26556, 26576, 26596, 26612, 26636} {
		*memmap.PtrPtr(0x587000, 26488+4*uintptr(i)) = memmap.PtrOff(0x587000, off)
	}
	t.Cleanup(func() { copy(table, saved) })

	c, _, _ := newEffectsFullOwner(t)
	o := &floorAssetsOwner{refs: map[uint32]uint32{0: 0}}
	var undo func()
	o.defs, o.edges, o.count, o.edgeCount, undo = legacy.PortTestFloorAssetsOwner()
	t.Cleanup(undo)
	raw := make([][]byte, 8)
	for i := range raw {
		raw[i] = spriteAnimationTestImage(i)
	}
	imgs, freeImages := c.r.GetBag().PortTestSpriteImages(raw)
	t.Cleanup(freeImages)
	for i, img := range imgs {
		o.refs[uint32(uintptr(img.C()))] = uint32(i + 1)
	}
	for _, off := range []uintptr{251576, 251580, 1193192, 1193196} {
		p := memmap.PtrUint8(0x5D4594, off)
		old := *p
		*p = 0
		t.Cleanup(func() { *p = old })
	}
	var freeScratch func()
	o.scratch, freeScratch = alloc.Make([]byte{}, 256*1024)
	t.Cleanup(freeScratch)
	return o
}
func (o *floorAssetsOwner) reset() {
	legacy.PortTestFloorAssetsRelease(o.defs, o.edges)
	clear(o.defs)
	clear(o.edges)
	*o.count = 0
	*o.edgeCount = 0
	for i := range o.scratch {
		o.scratch[i] = 0xa5
	}
}
func floorAssetsCapture(t *testing.T, label string, rows any, want string) {
	t.Helper()
	b, e := json.Marshal(rows)
	if e != nil {
		t.Fatal(e)
	}
	if p := os.Getenv("OPENNOX_FLOOR_ASSETS_CAPTURE"); p != "" {
		if e := os.WriteFile(p+"-"+label+".json", b, 0600); e != nil {
			t.Fatal(e)
		}
	}
	got := fmt.Sprintf("%x", sha256.Sum256(b))
	t.Logf("%s: %s", label, got)
	if want != got {
		t.Fatalf("%s got%s want%s", label, got, want)
	}
}
func floorAssetsWord(b []byte, v uint32) []byte { return binary.LittleEndian.AppendUint32(b, v) }
func floorAssetsFrames(b []byte, count, mode int) []byte {
	for i := 0; i < count; i++ {
		if mode == 1 || (mode == 2 && i%2 != 0) {
			b = floorAssetsWord(b, 0xffffffff)
			name := []byte(fmt.Sprintf("outside-%d", i))
			b = append(b, byte(i%7), byte(len(name)))
			b = append(b, name...)
		} else {
			b = floorAssetsWord(b, uint32(i%8))
		}
	}
	return b
}
func floorAssetsFloor(name string, color [3]byte, dims [3]byte, mode int) []byte {
	b := floorAssetsWord(nil, 0x10203040)
	b = append(b, byte(len(name)))
	b = append(b, name...)
	b = append(b, color[:]...)
	b = floorAssetsWord(b, 0xfffffff1)
	b = floorAssetsWord(b, 0x12345678)
	b = append(b, 7, dims[0], dims[1], dims[2], 5)
	b = floorAssetsFrames(b, int(dims[0])*int(dims[1])*int(dims[2]), mode)
	return floorAssetsWord(b, 0x454e4420)
}
func floorAssetsEdge(name string, dims [3]byte, mode int, badFormat byte, end uint32) []byte {
	b := floorAssetsWord(nil, 0x11223344)
	b = append(b, byte(len(name)))
	b = append(b, name...)
	b = floorAssetsWord(b, 0xfffffff9)
	b = floorAssetsWord(b, 0x10203040)
	b = append(b, 9, dims[0], 7, badFormat, dims[1], dims[2])
	b = floorAssetsFrames(b, 2*int(dims[0])*(int(dims[1])+int(dims[2])), mode)
	return floorAssetsWord(b, end)
}
func (o *floorAssetsOwner) invoke(t *testing.T, op int, input []byte) (ret, consumed int) {
	t.Helper()
	raw, _ := alloc.CloneSlice(append(append([]byte(nil), input...), 0xc7, 0x5e))
	mf := binfile.NewMemFile(unsafe.Pointer(&raw[0]), len(raw))
	defer mf.Free()
	switch op {
	case 0:
		ret = legacy.Nox_thing_read_FLOR_411540(mf, o.scratch)
	case 1:
		ret = legacy.Nox_thing_read_EDGE_411850(mf, o.scratch)
	case 2:
		ret = legacy.Nox_thing_read_FLOR_414DB0(mf)
	case 3:
		ret = legacy.Nox_thing_read_EDGE_414E70(mf, o.scratch)
	case 4:
		ret = legacy.Nox_thing_read_floor_485B30(mf, o.scratch)
	case 5:
		ret = legacy.Nox_thing_read_edge_485D40(mf, o.scratch)
	default:
		t.Fatal("reader op")
	}
	if sha256.Sum256(raw) != sha256.Sum256(append(append([]byte(nil), input...), 0xc7, 0x5e)) {
		t.Fatal("reader mutated input")
	}
	return ret, len(raw) - len(mf.Data())
}
func (o *floorAssetsOwner) metadata() ([]byte, []byte) {
	defs := append([]byte(nil), unsafe.Slice((*byte)(unsafe.Pointer(&o.defs[0])), len(o.defs)*60)...)
	edges := append([]byte(nil), o.edges...)
	for i := range o.defs {
		clear(defs[i*60+32 : i*60+36])
	}
	for i := 0; i < 64; i++ {
		clear(edges[i*60+32 : i*60+36])
	}
	return defs, edges
}
func (o *floorAssetsOwner) assertOnlyData(t *testing.T, edge bool, index int) {
	t.Helper()
	for i := range o.defs {
		if (edge || i != index) && o.defs[i].Data32 != nil {
			t.Fatalf("unexpected floor image array %d", i)
		}
	}
	for i := 0; i < 64; i++ {
		if (!edge || i != index) && *(*unsafe.Pointer)(unsafe.Pointer(&o.edges[60*i+32])) != nil {
			t.Fatalf("unexpected edge image array %d", i)
		}
	}
}
