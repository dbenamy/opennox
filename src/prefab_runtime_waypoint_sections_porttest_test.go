//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestPrefabRuntimeWaypointSections(t *testing.T) {
	s := newObjectXferOwner(t)
	words, restore := legacy.PortTestPrefabRuntimeGlobals()
	defer restore()
	path := filepath.Join(t.TempDir(), "waypoints.bin")
	type record struct {
		ID, X, Y, Flags uint32
		Name            string
		Links           [][2]uint32
	}
	type row struct {
		Name     string
		Return   uint64
		Position int64
		Records  []record
		Wire     []byte
	}
	var rows []row
	for _, version := range []uint16{0, 1, 2, 3, 4, 5, 0x8000, 0xffff} {
		for _, cache := range []bool{false, true} {
			noxflags.ResetGame()
			if cache {
				noxflags.SetGame(noxflags.GameFlag(0x400000))
			}
			for _, links := range []byte{0, 1, 2, 31, 32} {
				for _, nameLen := range []int{0, 1, 75} {
					for _, coord := range []uint32{0, 16777217, math.MaxUint32} {
						var b bytes.Buffer
						put := func(v any) {
							if err := binary.Write(&b, binary.LittleEndian, v); err != nil {
								t.Fatal(err)
							}
						}
						put(version)
						put(uint32(2))
						name := strings.Repeat("w", nameLen)
						for i := uint32(0); i < 2; i++ {
							put(uint32(100 + i))
							put(coord)
							put(coord + 1)
							if int16(version) >= 3 {
								b.WriteByte(byte(nameLen))
								b.WriteString(name)
							}
							put(uint32(0x80000001 + i))
							if int16(version) < 4 {
								put(uint32(links))
							} else {
								b.WriteByte(links)
							}
							for j := byte(0); j < links; j++ {
								put(uint32(100 + (j % 2)))
								if int16(version) >= 2 {
									b.WriteByte(j * 9)
								}
							}
						}
						if err := os.WriteFile(path, b.Bytes(), 0600); err != nil {
							t.Fatal(err)
						}
						if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
							t.Fatal(err)
						}
						ret := legacy.PortTestPrefabCall(34, [6]uint32{})
						pos, err := cryptfile.Global().File.File.Seek(0, io.SeekCurrent)
						if err != nil {
							t.Fatal(err)
						}
						cryptfile.Close()
						r := row{Name: fmt.Sprintf("read/v%d/cache%t/links%d/name%d/coord%x", version, cache, links, nameLen, coord), Return: ret, Position: pos}
						if int16(version) > 4 {
							if ret != 0 || pos != 2 || s.WPs.Pending != nil || *words["waypoints"] != 0 {
								t.Fatal("version rejection")
							}
						} else {
							if ret != 1 || pos != int64(b.Len()) {
								t.Fatal("waypoint read/cursor")
							}
							head := s.WPs.Pending
							if cache {
								head = *(**server.Waypoint)(unsafe.Pointer(uintptr(*words["waypoints"])))
							}
							count := 0
							for w := head; w != nil; w = w.WpNext {
								i := uint32(1 - count)
								wantX, wantY := coord, coord+1
								if int16(version) < 4 {
									wantX = math.Float32bits(float32(coord))
									wantY = math.Float32bits(float32(coord + 1))
								}
								wantName := ""
								if int16(version) >= 3 {
									wantName = name
								}
								gotName := alloc.GoString(&w.NameBuf[0])
								if w.Index != 100+i || math.Float32bits(w.PosVec.X) != wantX || math.Float32bits(w.PosVec.Y) != wantY || w.Flags != 0x80000001+i || gotName != wantName || w.PointsCnt != links {
									t.Fatalf("%s fields differ", r.Name)
								}
								rec := record{ID: w.Index, X: wantX, Y: wantY, Flags: w.Flags, Name: gotName}
								for j := byte(0); j < links; j++ {
									kind := byte(2)
									if int16(version) >= 2 {
										kind = j * 9
									}
									if w.Field348[j] != uint32(100+j%2) || w.Points[j].Ind != kind || w.Points[j].Waypoint != nil {
										t.Fatal("waypoint link layout")
									}
									rec.Links = append(rec.Links, [2]uint32{w.Field348[j], uint32(w.Points[j].Ind)})
								}
								r.Records = append(r.Records, rec)
								count++
							}
							if count != 2 {
								t.Fatal("waypoint count")
							}
							for w := head; w != nil; {
								next := w.WpNext
								alloc.Free(w)
								w = next
							}
							s.WPs.Pending = nil
							for n := *words["waypoints"]; n != 0; {
								p := unsafe.Pointer(uintptr(n))
								next := *(*uint32)(unsafe.Add(p, 4))
								legacy.MapPrefabFreeNode(p)
								n = next
							}
							*words["waypoints"] = 0
						}
						rows = append(rows, r)
					}
				}
			}
		}
	}
	noxflags.ResetGame()
	for _, count := range []int{0, 1, 3} {
		for _, links := range []byte{0, 1, 31, 32} {
			for _, length := range []int{0, 1, 75} {
				for i := 0; i < count; i++ {
					w := s.NewWaypoint(types.Ptf(float32(i)+0.25, -float32(i)-0.75))
					w.Index = uint32(100 + i)
					w.Flags = 0x80000000 | uint32(i&1)
					copy(w.NameBuf[:], strings.Repeat("n", length))
					w.PointsCnt = links
					for j := byte(0); j < links; j++ {
						w.Points[j].Waypoint = w
						w.Points[j].Ind = j * 9
					}
				}
				var want bytes.Buffer
				put := func(v any) { binary.Write(&want, binary.LittleEndian, v) }
				put(uint16(4))
				put(uint32(count))
				for w := s.WPs.List; w != nil; w = w.WpNext {
					put(w.Index)
					put(math.Float32bits(w.PosVec.X))
					put(math.Float32bits(w.PosVec.Y))
					want.WriteByte(byte(length))
					want.WriteString(strings.Repeat("n", length))
					put(w.Flags & 1)
					want.WriteByte(links)
					for j := byte(0); j < links; j++ {
						put(w.Index)
						want.WriteByte(j * 9)
					}
				}
				if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
					t.Fatal(err)
				}
				ret := legacy.PortTestPrefabCall(34, [6]uint32{})
				cryptfile.Close()
				got, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				if ret != 1 || !bytes.Equal(got, want.Bytes()) {
					t.Fatal("waypoint write wire")
				}
				rows = append(rows, row{Name: fmt.Sprintf("write/count%d/links%d/name%d", count, links, length), Return: ret, Wire: got})
				s.Nox_xxx_waypointDeleteAll_579DD0()
			}
		}
	}
	spellbookCapture(t, "prefab-runtime-waypoint-sections", rows, "0820d8ae08b9447ea29965fce70b9b5caeea997d3bdd0fb2edb12d7f5e72d010")
}
