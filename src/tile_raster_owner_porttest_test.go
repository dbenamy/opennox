//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"os"
	"testing"
	"unsafe"
)

type tileRasterOwner struct {
	*objectRenderOwner
	words  map[string]*uint32
	defs   []server.TileDef
	images []*noxrender.Image
	raw    []byte
}

func tileRasterCapture(t *testing.T, label string, rows any, want string) {
	t.Helper()
	raw, err := json.Marshal(rows)
	if err != nil {
		t.Fatal(err)
	}
	if prefix := os.Getenv("OPENNOX_TILE_RASTER_CAPTURE"); prefix != "" {
		if err := os.WriteFile(prefix+"-"+label+".json", raw, 0600); err != nil {
			t.Fatal(err)
		}
	}
	got := fmt.Sprintf("%x", sha256.Sum256(raw))
	t.Logf("%s: %s", label, got)
	// Expectations captured from original production C.
	if got != want {
		t.Fatalf("%s got%s want%s", label, got, want)
	}
}
func newTileRasterOwner(t *testing.T) *tileRasterOwner {
	o := &tileRasterOwner{objectRenderOwner: newObjectRenderOwner(t)}
	var restore func()
	o.words, o.defs, restore = legacy.PortTestTileRasterOwner()
	if o.c.tiles.noxTileBuf != nil {
		t.Fatal("fresh owner already has tile buffer")
	}
	oldTextured := nox_client_texturedFloors_154956
	t.Cleanup(func() { nox_client_texturedFloors_154956 = oldTextured })
	for _, r := range [][3]uintptr{{0x973CE0, 0, 576}, {0x973F18, 7696, 4}} {
		b := unsafe.Slice((*byte)(memmap.PtrOff(r[0], r[1])), r[2])
		old := append([]byte(nil), b...)
		t.Cleanup(func() { copy(b, old) })
	}
	*memmap.PtrUint32(0x973F18, 7696) = 1
	o.c.nox_xxx_initSight_485F80()
	o.c.nox_xxx_tileInitBuf_430DB0(96, 96)
	free := o.c.tiles.noxTileBufFree
	t.Cleanup(func() { restore(); free(); o.c.tiles.noxTileBuf = nil; o.c.tiles.noxTileBufFree = nil })
	legacy.Nox_xxx_tile_486060()
	o.raw = make([]byte, 2116)
	for i := range o.raw {
		o.raw[i] = byte(i*43 + i/7 + 19)
	}
	var undo func()
	o.images, undo = o.c.r.GetBag().PortTestWallEdgeImages([][]byte{o.raw}, []int{0})
	t.Cleanup(undo)
	return o
}
func (o *tileRasterOwner) bytes() []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(&o.c.tiles.noxTileBuf[0])), 2*len(o.c.tiles.noxTileBuf))
}
func (o *tileRasterOwner) resetBuffer() {
	for i := range o.c.tiles.noxTileBuf {
		o.c.tiles.noxTileBuf[i] = uint16(i*37 + 0xa351)
	}
	*o.words["flatFlag"] = 0
	*o.words["dirty"] = 0
}
func tileRasterRow(y int) (x, width, source int) {
	radius := min(y, 45-y)
	x, width = 23-radius, 2*radius+1
	for r := 0; r < y; r++ {
		source += 2 * (2*min(r, 45-r) + 1)
	}
	return
}

// Independent per-pixel expectation for the unrolled word-store phase.
func tileRasterFastFill(dst []byte, x, width int, color uint32) {
	for i := 0; i < width; i++ {
		value := uint16(color)
		if (x+i)&1 != 0 && i != 0 {
			value = uint16(color >> 16)
		}
		binary.LittleEndian.PutUint16(dst[2*i:], value)
	}
}
