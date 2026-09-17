//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"testing"

	"github.com/opennox/opennox/v1/client"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
)

func TestItemXferObeliskMinimap(t *testing.T) {
	c, _, _ := newEffectsFullOwner(t)
	s := newItemXferOwner(t)
	path := filepath.Join(t.TempDir(), "obelisk.bin")
	var rows []itemXferCaptureRow
	defer func() {
		spellbookCapture(t, "item-xfer-obelisk-minimap", rows, "682faf05e991e354997ef2f75194e6e9b62fe0ce0c3c9e05cb71d86d61a41d1c")
	}()
	for _, clientMode := range []bool{false, true} {
		for scene := 0; scene < 4; scene++ {
			for mini := 0; mini < 3; mini++ {
				for _, mana := range []uint32{0, 1, 50, 0xffffffff} {
					t.Run(fmt.Sprintf("client%v-scene%d-mini%d-mana%x", clientMode, scene, mini, mana), func(t *testing.T) {
						flags := noxflags.GameFlag(0x200000)
						if clientMode {
							flags |= 0x800
						}
						t.Cleanup(noxflags.PortTestGameFlags(flags))
						u := newItemXferObject(t, s, "Obelisk")
						objectXferSetCommon(u)
						objectXferSetWord(u.UpdateData, 0, mana)
						var dr *client.Drawable
						if scene != 0 {
							dr = c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(64, 128))
							if dr == nil {
								t.Fatal("drawable allocation")
							}
							defer c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
							dr.ObjClass = 0x400000
							dr.NetCode32 = u.Extent
							if scene == 1 {
								dr.NetCode32++
							}
							if scene == 2 {
								dr.ObjClass = 1
							}
							if mini > 0 {
								c.Objs.MinimapAdd(dr, 1)
							}
						}
						if mini == 2 {
							other := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(80, 144))
							if other == nil {
								t.Fatal("other drawable allocation")
							}
							defer c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(other)
							other.NetCode32 = u.Extent + 2
							c.Objs.MinimapAdd(other, 1)
						}
						want := objectXferCurrentRecord(61)
						want.u32(mana)
						visible := byte(0)
						if clientMode && scene == 3 && mini > 0 {
							visible = 1
						}
						want.u8(visible)
						if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
							t.Fatal(err)
						}
						if err := u.CallXfer(nil); err != nil {
							t.Fatal(err)
						}
						checksum := cryptfile.Global().PortTestChecksum()
						if err := cryptfile.Close(); err != nil {
							t.Fatal(err)
						}
						got, err := os.ReadFile(path)
						if err != nil {
							t.Fatal(err)
						}
						if !bytes.Equal(got, want.Bytes()) {
							t.Fatalf("obelisk bytes=%x want=%x", got, want.Bytes())
						}
						if objectXferGetWord(u.UpdateData, 0) != mana {
							t.Fatal("writer mutated mana")
						}
						itemXferCaptureCase(t, &rows, u, 0, checksum)
					})
				}
			}
		}
	}
}
