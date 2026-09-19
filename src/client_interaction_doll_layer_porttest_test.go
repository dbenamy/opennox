//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"image"
	"slices"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestClientInteractionDollLayer(t *testing.T) {
	o := newInventoryDisplayOwner(t)
	table, free := alloc.New([27]uint32{})
	defer free()
	cloak := serverConfigOwnBytes(t, 0x5D4594, 1319052, 4)
	for i := range table {
		table[i] = uint32(uintptr(o.images[i].C()))
	}
	binary.LittleEndian.PutUint32(cloak, uint32(uintptr(o.images[31].C())))
	type row struct {
		Cloth, Definition, Present, Overlay bool
		Mods, Index                         int
		Draws                               int
		Pixels                              string
	}
	var captured []row
	for _, cloth := range []bool{false, true} {
		for _, definition := range []bool{false, true} {
			for _, present := range []bool{false, true} {
				for _, overlay := range []bool{false, true} {
					for mask := 0; mask < 16; mask++ {
						o.reset(t)
						dr := o.item(t, "Bow", 77)
						dr.ObjClass = 0x1000000
						def := o.weapon
						if cloth {
							dr.ObjClass = 0x2000000
							def = o.armor
						}
						def.TypeInd = 0
						if definition {
							def.TypeInd = dr.TypeIDVal
						}
						if present {
							o.equipment[8] = uint32(uintptr(dr.C()))
						}
						for i := 0; i < 4; i++ {
							*txword(dr, 432+uintptr(4*i)) = 0
							if mask&(1<<i) != 0 {
								*txword(dr, 432+uintptr(4*i)) = uint32(uintptr(unsafe.Pointer(o.mods[i])))
							}
						}
						resetMaterials := func() {
							for i := 1; i <= 6; i++ {
								o.c.r.Data().SetMaterialRGB(i, 100+i, 150+i, 200+i)
							}
						}
						resetMaterials()
						clear(o.pix.Pix)
						o.drawTrace = nil
						index := []int{0, 1, 25, 26}[mask%4]
						flag := uintptr(0)
						if overlay {
							flag = 1
						}
						ret := uint32(interactionCall("sub_4BF9F0", 0xffffffff, uintptr(dr.TypeIDVal), 21, 18, uintptr(unsafe.Pointer(table)), uintptr(index), flag))
						wantReturn := uint32(0)
						if present {
							wantReturn = uint32(int32(int16(uintptr(dr.C()))))
						}
						if ret != wantReturn {
							t.Fatal("paper-doll layer signed pointer return", cloth, definition, present, ret, wantReturn)
						}
						got := append([]uint16(nil), o.pix.Pix...)
						expectedDraws := 0
						if present {
							expectedDraws = 1
						}
						if len(o.drawTrace) != expectedDraws {
							t.Fatal("doll missing equipment gate", len(o.drawTrace), expectedDraws)
						}
						clear(o.pix.Pix)
						resetMaterials()
						if present {
							if definition {
								for i := 1; i <= 6; i++ {
									c := def.Colors12[i]
									o.c.r.Data().SetMaterialRGB(i, int(c.R), int(c.G), int(c.B))
								}
								for i := 0; i < 4; i++ {
									if mask&(1<<i) != 0 {
										c := o.mods[i].Color24
										o.c.r.Data().SetMaterialRGB(int(def.ColorIndexes()[i]), int(c.R), int(c.G), int(c.B))
									}
								}
							}
							img := o.images[index]
							if overlay {
								img = o.images[31]
							}
							o.c.r.DrawImageAt(img, image.Pt(21, 18))
						}
						if !slices.Equal(got, o.pix.Pix) {
							t.Fatal("doll layer material/pixel contract", cloth, definition, present, overlay, mask, index)
						}
						nonzero := false
						raw := make([]byte, 2*len(got))
						for i, v := range got {
							nonzero = nonzero || v != 0
							binary.LittleEndian.PutUint16(raw[2*i:], v)
						}
						if nonzero != present {
							t.Fatal("doll layer nonempty pixels", present)
						}
						captured = append(captured, row{cloth, definition, present, overlay, mask, index, expectedDraws, fmt.Sprintf("%x", sha256.Sum256(raw))})
					}
				}
			}
		}
	}
	interactionCapture(t, "doll-layer", captured)
}
