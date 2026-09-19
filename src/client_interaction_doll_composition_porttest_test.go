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

	noxcolor "github.com/opennox/libs/color"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestClientInteractionDollComposition(t *testing.T) {
	var names []string
	for i := 0; i < 53; i++ {
		names = append(names, fmt.Sprintf("DollPiece%d", i))
	}
	o := newInventoryDisplayOwner(t, names...)
	o.reset(t)
	words, restore := legacy.PortTestClientInteractionWords()
	defer restore()
	active := serverConfigOwnBytes(t, 0x852978, 8, 4)
	table := serverConfigOwnBytes(t, 0x973A20, 16, 456)
	cloak := serverConfigOwnBytes(t, 0x5D4594, 1319052, 4)
	imageID := func(img *noxrender.Image) uint32 { return uint32(uintptr(img.C())) }
	for gender := 0; gender < 2; gender++ {
		binary.LittleEndian.PutUint32(table[4*gender:], imageID(o.images[gender]))
		binary.LittleEndian.PutUint32(table[8+4*gender:], imageID(o.images[2+gender]))
		for i := 0; i < 26; i++ {
			binary.LittleEndian.PutUint32(table[16+104*gender+4*i:], imageID(o.images[(4+3*gender+i)%32]))
		}
		for i := 0; i < 27; i++ {
			binary.LittleEndian.PutUint32(table[240+108*gender+4*i:], imageID(o.images[(7+5*gender+i)%32]))
		}
	}
	binary.LittleEndian.PutUint32(cloak, imageID(o.images[31]))
	armorTypes, weaponTypes := map[uint16]uint32{}, map[uint16]uint32{}
	var equipped []*client.Drawable
	for i, name := range names {
		dr := o.item(t, name, uint32(i+1))
		equipped = append(equipped, dr)
		if i < 26 {
			dr.ObjClass = 0x2000000
			armorTypes[uint16(dr.TypeIDVal)] = 1 << i
		} else {
			dr.ObjClass = 0x1000000
			weaponTypes[uint16(dr.TypeIDVal)] = 1 << (i - 26)
		}
		if i > 0 {
			*txword(equipped[i-1], 368) = uint32(uintptr(dr.C()))
		}
	}
	o.equipment[0] = uint32(uintptr(equipped[0].C()))
	restoreTypes := o.c.srv.Server.PortTestInventoryEnvironment(false, weaponTypes, armorTypes)
	defer restoreTypes()
	// Layer-specific tests cover modifier definitions. Here all pieces intentionally
	// retain the player's six body materials so composition/order is isolated.
	o.weapon.TypeInd, o.armor.TypeInd = 0, 0
	pl := &o.players[0]
	palette := []uint32{0x7c00, 0x03e0, 0x001f, 0x4210, 0x7fe0, 0x03ff}
	point, free := alloc.New([2]int32{})
	*point = [2]int32{3, 4}
	defer free()
	pos := image.Pt(14, 19)
	type row struct {
		Gender               int
		Same, Active, Player bool
		Armor, Weapon        uint32
		ReturnKind           string
		Return               int32
		Images               []uint32
		Pixels               string
	}
	var captured []row
	for gender := 0; gender < 2; gender++ {
		for _, same := range []bool{false, true} {
			for _, hasActive := range []bool{false, true} {
				for _, hasPlayer := range []bool{false, true} {
					for _, armor := range []uint32{0, 1, 2, 1 << 23, 1 << 24, 1 << 25, 0x3000002, 0x3ffffff} {
						for _, weapon := range []uint32{0, 1, 1 << 26, 0x7ffffff, 0x80000000} {
							colors := append([]uint32(nil), palette...)
							if same {
								colors[0] = colors[1]
							}
							for i, c := range colors {
								*(*uint32)(unsafe.Add(unsafe.Pointer(pl), 2292+4*i)) = c
							}
							*(*byte)(unsafe.Add(unsafe.Pointer(pl), 2252)) = byte(gender)
							*(*uint32)(unsafe.Pointer(pl)) = armor
							*(*uint32)(unsafe.Add(unsafe.Pointer(pl), 4)) = weapon
							binary.LittleEndian.PutUint32(active, 0)
							if hasActive {
								binary.LittleEndian.PutUint32(active, uint32(uintptr(equipped[0].C())))
							}
							*words["dword_8531A0_2576"] = 0
							if hasPlayer {
								*words["dword_8531A0_2576"] = uint32(uintptr(unsafe.Pointer(pl)))
							}
							for i := range o.pix.Pix {
								o.pix.Pix[i] = 0x4210
							}
							o.drawTrace = nil
							gotReturn := uint32(interactionCall("sub_4BF7E0", uintptr(unsafe.Pointer(point))))
							got := append([]uint16(nil), o.pix.Pix...)
							var gotImages []uint32
							for _, d := range o.drawTrace {
								gotImages = append(gotImages, d.Image)
							}
							var indices []int
							expectedReturn := uint32(0)
							kind := "zero"
							normalized := int32(0)
							if hasActive {
								expectedReturn = uint32(int32(int16(uintptr(equipped[0].C()))))
								kind = "active-pointer"
								if hasPlayer {
									body := gender
									if same {
										body += 2
									}
									indices = append(indices, body)
									for i := 0; i < 24; i++ {
										if armor&(1<<i) != 0 {
											indices = append(indices, (4+3*gender+i)%32)
										}
									}
									if armor&2 != 0 {
										indices = append(indices, 31)
									}
									for i := 24; i < 26; i++ {
										if armor&(1<<i) != 0 {
											indices = append(indices, (4+3*gender+i)%32)
										}
									}
									for i := 0; i < 27; i++ {
										if weapon&(1<<i) != 0 {
											indices = append(indices, (7+5*gender+i)%32)
										}
									}
									expectedReturn = uint32(int32(int16(weapon)))
									kind = "weapon-mask"
									normalized = int32(int16(weapon))
									if weapon&(1<<26) != 0 {
										expectedReturn = uint32(int32(int16(uintptr(equipped[52].C()))))
										kind = "weapon26-pointer"
										normalized = 0
									}
								}
							}
							if gotReturn != expectedReturn {
								t.Fatal("doll composition return", hasActive, hasPlayer, weapon, gotReturn, expectedReturn)
							}
							var wantImages []uint32
							for _, i := range indices {
								wantImages = append(wantImages, o.c.imageRefs[imageID(o.images[i])])
							}
							if !slices.Equal(gotImages, wantImages) {
								t.Fatal("doll layer order", gender, armor, weapon, gotImages, wantImages)
							}
							for i := range o.pix.Pix {
								o.pix.Pix[i] = 0x4210
							}
							o.c.r.DrawRectFilledOpaque(pos.X, pos.Y, 200, 200, noxcolor.RGBA5551(0))
							if hasActive && hasPlayer {
								for i, c := range []uint32{colors[1], colors[3], colors[5], colors[4], colors[2], colors[0]} {
									o.c.r.Data().SetMaterial(i+1, noxcolor.RGBA5551(c))
								}
								for _, i := range indices {
									o.c.r.DrawImageAt(o.images[i], pos)
								}
							}
							if !slices.Equal(got, o.pix.Pix) {
								t.Fatal("doll body palette/composition pixels", gender, same, hasActive, hasPlayer, armor, weapon)
							}
							raw := make([]byte, 2*len(got))
							for i, v := range got {
								binary.LittleEndian.PutUint16(raw[2*i:], v)
							}
							captured = append(captured, row{gender, same, hasActive, hasPlayer, armor, weapon, kind, normalized, gotImages, fmt.Sprintf("%x", sha256.Sum256(raw))})
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "doll-composition", captured)
}
