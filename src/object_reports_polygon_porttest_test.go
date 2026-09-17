//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"image"
	"testing"
	"unsafe"
)

func objectReportsPolygonOwner(t *testing.T) *minimapOwner {
	t.Helper()
	words, restore := legacy.PortTestMinimapWords()
	t.Cleanup(restore)
	*words["polygons"] = 1
	for _, region := range [][2]uintptr{{535844, 16 * 12}, {552228, 140 * 4}, {588080, 4}} {
		p := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, region[0])), region[1])
		old := append([]byte(nil), p...)
		clear(p)
		t.Cleanup(func() { copy(p, old) })
	}
	o := &minimapOwner{}
	for i := 0; i < 2; i++ {
		v, free := alloc.Make([]uint32{}, 4)
		t.Cleanup(free)
		o.polygonVertices = append(o.polygonVertices, v)
	}
	return o
}
func TestObjectReportsPolygonLevels(t *testing.T) {
	_, units, _ := objectReportsPlayers(t)
	o := objectReportsPolygonOwner(t)
	o.polygon(t, 0, image.Rect(10, 10, 20, 20), 7)
	o.polygon(t, 1, image.Rect(30, 30, 40, 40), 200)
	monster, free := alloc.New(server.Object{})
	t.Cleanup(free)
	data, fd := alloc.Make([]uint32{}, 1)
	t.Cleanup(fd)
	monster.ObjClass = 2
	monster.UpdateData = unsafe.Pointer(&data[0])
	var rows []struct {
		Name   string
		Return int
	}
	defer func() {
		spellbookCapture(t, "object-reports-polygons", rows, "16288be028cb8173fa7b6d252cdcb50cb42b2e9e3483046554461997a07fd46e")
	}()
	for _, kind := range []string{"nil", "other", "player", "monster", "monster-unassigned"} {
		for _, cached := range []byte{0, 1, 7, 127, 128, 200, 255} {
			for _, point := range []types.Pointf{{0, 0}, {15, 16}, {35, 36}, {9.9, 15}, {10.1, 15}, {19.9, 15}, {20.1, 15}, {-15, -15}} {
				name := fmt.Sprintf("%s/cache%d/point%v", kind, cached, point)
				t.Run(name, func(t *testing.T) {
					var u *server.Object
					switch kind {
					case "other":
						u = monster
						u.ObjClass = 0
					case "player":
						u = &units[0]
						*(*byte)(unsafe.Add(unsafe.Pointer(u.UpdateDataPlayer().Player), 3668)) = cached
					case "monster":
						u = monster
						u.ObjClass = 2
						data[0] = 1
						*(*byte)(memmap.PtrOff(0x5D4594, 552228+140+130)) = cached
					case "monster-unassigned":
						u = monster
						u.ObjClass = 2
						data[0] = 0xdeadface
					}
					// Reset polygon one unless this case intentionally exercises its cached level.
					if kind != "monster" {
						*(*byte)(memmap.PtrOff(0x5D4594, 552228+140+130)) = 7
					}
					rv := legacy.PortTestObjectReports(8, nil, u, 0, 0, 0, &point)
					if (kind == "player" || kind == "monster") && cached != 0 {
						if rv != int(int8(cached)) {
							t.Fatalf("cached level %d", rv)
						}
					} else if point.X == 15 {
						want := 7
						if kind == "monster" {
							want = int(int8(cached))
						}
						if rv != want {
							t.Fatalf("first polygon %d want%d", rv, want)
						}
					} else if point.X == 35 && rv != -56 {
						t.Fatalf("second polygon signed level %d", rv)
					} else if (point.X == 0 || point.X < 0 || point.X == 9.9) && rv != 0 {
						t.Fatalf("outside level %d", rv)
					}
					rows = append(rows, struct {
						Name   string
						Return int
					}{name, rv})
				})
			}
		}
	}
}

func TestObjectReportsPolygonVertexRays(t *testing.T) {
	objectReportsPlayers(t)
	o := objectReportsPolygonOwner(t)
	o.polygon(t, 0, image.Rect(10, 10, 20, 20), 7)
	var rows []struct {
		Counter uint32
		Point   types.Pointf
		Return  int
	}
	defer func() {
		spellbookCapture(t, "object-reports-polygon-vertices", rows, "ef2551055a27807ed42cfb9d15f005527de5747b5f0cfd443d05c53c6c51d4ba")
	}()
	for counter := uint32(0); counter < 4; counter++ {
		for _, p := range []types.Pointf{{15, 15}, {15, 16}, {10, 10}, {20, 20}, {10, 15}, {20, 15}, {9.99, 15}, {20.01, 15}} {
			t.Run(fmt.Sprintf("counter%d/point%v", counter, p), func(t *testing.T) {
				*memmap.PtrUint32(0x5D4594, 588080) = counter
				rv := legacy.PortTestObjectReports(8, nil, nil, 0, 0, 0, &p)
				if rv != 0 && rv != 7 {
					t.Fatalf("invalid level%d", rv)
				}
				rows = append(rows, struct {
					Counter uint32
					Point   types.Pointf
					Return  int
				}{counter, p, rv})
			})
		}
	}
}
