//go:build porttest

package opennox

import (
	"bytes"
	"errors"
	"image"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/noxcrypt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/cnxz"
	"github.com/opennox/opennox/v1/server"
)

func TestWorldGridMapSave(t *testing.T) {
	c, _, _ := newEffectsFullOwner(t)
	oldServer := noxServer
	noxServer = c.srv
	t.Cleanup(func() { noxServer = oldServer })
	t.Cleanup(c.srv.PortTestMinimapWalls())
	oldFile := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	t.Cleanup(func() { cryptfile.Close(); cryptfile.SetGlobal(oldFile) })
	oldSections := noxMapSections
	t.Cleanup(func() { noxMapSections = oldSections })
	for _, reg := range [][3]uintptr{{0x587000, 253112, 5}, {0x5D4594, 2487252, 8}, {0x5D4594, 739980, 8}} {
		b := memmap.Slice(reg[0], reg[1])[:reg[2]]
		old := append([]byte(nil), b...)
		t.Cleanup(func() { copy(b, old) })
	}
	copy(memmap.Slice(0x587000, 253112)[:5], []byte(".nxz\x00"))
	c.srv.Map.Debug.Add("fixture", "map-save")
	c.srv.Map.Debug.Add("fixture", "second-value")
	type result struct {
		Mode, Compression, Walls int
		Return                   bool
		Min                      [2]uint32
		Size                     [2]uint32
		Wire                     []byte
		Dir                      byte
	}
	var rows []result
	for walls := 0; walls < 2; walls++ {
		var magic *server.MagicWall
		if walls != 0 {
			w := c.srv.Walls.CreateAtGrid(image.Pt(17, 23))
			w.Dir0 = 7
			c.srv.Walls.CreateAtGrid(image.Pt(31, 9))
			magic = &server.MagicWall{Field0: 1, Dir4: 7, PrevDir13: 2, Wall8: w}
			c.srv.spells.walls.list = magic
		}
		for mode := 0; mode < 3; mode++ {
			for compression := 0; compression < 2; compression++ {
				if magic != nil {
					magic.Wall8.Dir0 = 7
				}
				noxMapSections = []mapSection{{Name: "DebugData", Fnc: nox_server_mapRWDebugData_5060D0}}
				if mode == 1 {
					noxMapSections = append(noxMapSections, mapSection{Name: "FixtureFailure", Fnc: func(cf *cryptfile.CryptFile, _ unsafe.Pointer) error {
						cf.WriteU32(0x12345678)
						return errors.New("intentional section failure")
					}})
				}
				path := filepath.Join(t.TempDir(), "world.map")
				if mode == 2 {
					path = filepath.Join(path, "missing.map")
				}
				ok := legacy.Nox_xxx_mapSaveMap_51E010(path, compression)
				if ok != (mode == 0) || cryptfile.Global() != nil {
					t.Fatalf("mode%d comp%d return%v file%v", mode, compression, ok, cryptfile.Global())
				}
				r := result{Mode: mode, Compression: compression, Walls: walls, Return: ok, Min: [2]uint32{memmap.Uint32(0x5D4594, 2487252), memmap.Uint32(0x5D4594, 2487256)}, Size: [2]uint32{memmap.Uint32(0x5D4594, 739980), memmap.Uint32(0x5D4594, 739984)}}
				if magic != nil {
					r.Dir = magic.Wall8.Dir0
					want := byte(7)
					if mode == 1 {
						want = 2
					}
					if r.Dir != want {
						t.Fatalf("wall direction mode%d got%d want%d", mode, r.Dir, want)
					}
				}
				if mode != 2 {
					var err error
					r.Wire, err = os.ReadFile(path)
					if err != nil {
						t.Fatal(err)
					}
					want := [2]uint32{256, 256}
					if walls != 0 {
						want = [2]uint32{17, 9}
					}
					if r.Min != want || r.Size != want {
						t.Fatalf("wall bounds %+v want%v", r, want)
					}
					if mode == 0 {
						if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, crypt.MapKey); err != nil {
							t.Fatal(err)
						}
						cf := cryptfile.Global()
						if _, err := mapReadCryptHeader(cf); err != nil {
							t.Fatal(err)
						}
						x, err := cf.ReadU32()
						if err != nil {
							t.Fatal(err)
						}
						y, err := cf.ReadU32()
						if err != nil {
							t.Fatal(err)
						}
						if [2]uint32{x, y} != want {
							t.Fatalf("serialized bounds %d,%d", x, y)
						}
						name, err := cf.ReadString8()
						if err != nil || name != "DebugData" {
							t.Fatalf("section %q %v", name, err)
						}
						n, err := cf.ReadAlignedU32()
						if err != nil || n == 0 {
							t.Fatalf("section size %d %v", n, err)
						}
						version, err := cf.ReadU16()
						if err != nil || version != 1 {
							t.Fatalf("debug version %d %v", version, err)
						}
						count, err := cf.ReadU32()
						if err != nil || count != 2 {
							t.Fatalf("debug count %d %v", count, err)
						}
						for _, wantValue := range []string{"map-save", "second-value"} {
							key, err := cf.ReadString32()
							if err != nil || key != "fixture" {
								t.Fatalf("key %q %v", key, err)
							}
							v, err := cf.ReadString32()
							if err != nil || v != wantValue {
								t.Fatalf("value %q %v", v, err)
							}
						}
						end, err := cf.ReadU8()
						if err != nil || end != 0 {
							t.Fatalf("terminator %d %v", end, err)
						}
						cryptfile.Close()
						if compression != 0 {
							dst := filepath.Join(filepath.Dir(path), "roundtrip.map")
							if err := cnxz.DecompressFile(path[:len(path)-4]+".nxz", dst); err != nil {
								t.Fatal(err)
							}
							b, err := os.ReadFile(dst)
							if err != nil || !bytes.Equal(b, r.Wire) {
								t.Fatalf("compression round trip %v", err)
							}
						}
					}
				}
				rows = append(rows, r)
			}
		}
		c.srv.spells.walls.list = nil
	}
	drawableStateCapture(t, "world-save", rows, "5d3a18efb63c740b213dab4601409eb8800faa17fcf836199d498d1bb2652548")
}
