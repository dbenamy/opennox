//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
)

func TestTileRasterNilImageAndOrigins(t *testing.T) {
	o := newTileRasterOwner(t)
	type result struct {
		Setup     int
		Fill, Nil bool
		Flag      byte
		Flat      uint32
		Buffer    [32]byte
	}
	var rows []result
	for setup, values := range [][4]int32{{0, 0, 0, 0}, {31, 47, 0, 0}, {100, -13, 7, 11}} {
		for _, fill := range []bool{false, true} {
			for _, nilImage := range []bool{false, true} {
				for _, flag := range []byte{0, 1} {
					o.resetBuffer()
					o.defs[0].Color48 = 0x12341234
					o.defs[0].Field58 = flag
					defBefore := o.defs[0]
					*o.words["flatFlag"] = 7
					for i, key := range []string{"originX", "originY", "scrollX", "scrollY"} {
						*o.words[key] = uint32(values[i])
					}
					nox_client_texturedFloors_154956 = !fill
					nox_xxx_tileSetDrawFn_481420()
					beforeWords := map[string]uint32{}
					for key, p := range o.words {
						beforeWords[key] = *p
					}
					pos := image.Pt(40+int(values[0]-values[2]), 188+int(values[1]-values[3]))
					b := o.bytes()
					want := append([]byte(nil), b...)
					stride := int(*o.words["stride"])
					if fill || !nilImage {
						for y := 0; y < 46; y++ {
							x, width, src := tileRasterRow(y)
							for i := 0; i < width; i++ {
								index := (188*stride + 80 + y*stride + 2*(x+i)) % len(b)
								value := uint16(0x1234)
								if !fill {
									value = binary.LittleEndian.Uint16(o.raw[src+2*i:])
								}
								binary.LittleEndian.PutUint16(want[index:], value)
							}
						}
					}
					var img noxrender.ImageHandle
					if !nilImage {
						img = o.images[0].C()
					}
					legacy.PortTestTileRasterDispatch(pos, img, 0)
					if !bytes.Equal(b, want) {
						t.Fatalf("setup%d fill%v nil%v flag%d changed wrong pixels", setup, fill, nilImage, flag)
					}
					expected := uint32(7)
					if fill && flag != 0 {
						expected = 1
					}
					for key, p := range o.words {
						v := beforeWords[key]
						if key == "flatFlag" {
							v = expected
						}
						if *p != v {
							t.Fatalf("unexpected word change %s", key)
						}
					}
					if o.defs[0] != defBefore {
						t.Fatal("raster callback changed tile definition")
					}
					rows = append(rows, result{setup, fill, nilImage, flag, *o.words["flatFlag"], sha256.Sum256(b)})
				}
			}
		}
	}
	tileRasterCapture(t, "nil-origins", rows, "5dbfe27adaddadff79aa3874cc1f2a571bb86dc293be746068ac98d97ba6b943")
}
