//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

type prefabGroupRecord struct {
	ID    uint32
	Kind  byte
	Name  string
	Items [][2]uint32
}

func prefabGroupSnapshot(g *server.MapGroup) prefabGroupRecord {
	r := prefabGroupRecord{ID: g.Index(), Kind: byte(g.GroupType()), Name: g.ID()}
	var prev *server.MapGroupItem
	for p := g.List; p != nil; p = p.Next8 {
		if p.Prev12 != prev {
			panic("group item reverse link")
		}
		r.Items = append(r.Items, [2]uint32{p.Raw0, p.Raw4})
		prev = p
	}
	return r
}
func prefabGroupOwner(t *testing.T) *server.Server {
	t.Helper()
	old := noxServer
	s := &server.Server{}
	s.MapGroups.Init()
	noxServer = &Server{Server: s}
	t.Cleanup(func() { prefabGroupReset(s); s.MapGroups.Free(); noxServer = old })
	return s
}
func prefabGroupReset(s *server.Server) {
	for p := s.MapGroups.Refs; p != nil; {
		next := p.Next4
		alloc.Free(p)
		p = next
	}
	s.MapGroups.Refs = nil
	s.MapGroups.Reset()
}
func TestPrefabRuntimeGroupRead(t *testing.T) {
	prefabRuntimeTables(t)
	s := prefabGroupOwner(t)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	current := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 2598188)), 80)
	oldName := append([]byte(nil), current...)
	clear(current)
	copy(current, "arena")
	t.Cleanup(func() { copy(current, oldName) })
	blocked := memmap.PtrUint32(0x5D4594, 739992)
	oldBlocked := *blocked
	t.Cleanup(func() { *blocked = oldBlocked })
	oldFile := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	t.Cleanup(func() { cryptfile.Close(); cryptfile.SetGlobal(oldFile) })
	path := filepath.Join(t.TempDir(), "groups.bin")
	type row struct {
		Name         string
		Return       uint64
		Position     int64
		Groups, Refs []prefabGroupRecord
	}
	var rows []row
	for _, version := range []uint16{0, 1, 2, 3, 4, 0x8000} {
		for kind := byte(0); kind < 5; kind++ {
			for _, flags := range []uint32{0, 1, 0x200000, 0x400000, 0x400001, 0x600000} {
				for _, wall := range []uint32{0, 4} {
					for _, length := range []int{0, 1, 46, 47, 63, 75} {
						if version == 2 && (length == 0 || length > 73) {
							continue
						}
						for _, items := range []int{0, 2} {
							prefabGroupReset(s)
							noxflags.ResetGame()
							noxflags.SetGame(noxflags.GameFlag(flags))
							*blocked = wall
							name := strings.Repeat("g", length)
							wireName := name
							if version == 2 {
								wireName = "o:" + name
							}
							var b bytes.Buffer
							binary.Write(&b, binary.LittleEndian, version)
							binary.Write(&b, binary.LittleEndian, uint32(2))
							for group := 0; group < 2; group++ {
								b.WriteByte(byte(len(wireName)))
								b.WriteString(wireName)
								b.WriteByte(kind)
								binary.Write(&b, binary.LittleEndian, uint32(101+group))
								binary.Write(&b, binary.LittleEndian, uint32(items))
								for i := 0; i < items; i++ {
									binary.Write(&b, binary.LittleEndian, uint32(7+10*group+i))
									if kind == 2 {
										binary.Write(&b, binary.LittleEndian, uint32(70+10*group+i))
									}
								}
							}
							if err := os.WriteFile(path, b.Bytes(), 0600); err != nil {
								t.Fatal(err)
							}
							if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
								t.Fatal(err)
							}
							ret := legacy.PortTestPrefabCall(33, [6]uint32{})
							// Position uses the underlying ordinary file, whose offsets match key -1.
							position, err := cryptfile.Global().File.File.Seek(0, io.SeekCurrent)
							if err != nil {
								t.Fatal(err)
							}
							r := row{Name: fmt.Sprintf("v%d/k%d/f%x/w%d/n%d/i%d", version, kind, flags, wall, length, items), Return: ret, Position: position}
							for g := s.MapGroups.GetFirstMapGroup(); g != nil; g = g.Next() {
								r.Groups = append(r.Groups, prefabGroupSnapshot(g))
							}
							var prev *server.MapGroupRef
							for ref := s.MapGroups.Refs; ref != nil; ref = ref.Next4 {
								if ref.Prev8 != prev {
									t.Fatal("group reference reverse link")
								}
								r.Refs = append(r.Refs, prefabGroupSnapshot(ref.Field0))
								prev = ref
							}
							wantRet := uint64(1)
							rejectedName := int16(version) < 3 && 5+1+length >= 53
							if int16(version) > 3 || rejectedName || kind > 3 && items > 0 {
								wantRet = 0
							}
							if ret != wantRet {
								t.Fatalf("%s return %d want %d", r.Name, ret, wantRet)
							}
							var expected []prefabGroupRecord
							if int16(version) <= 3 && !rejectedName && wall&4 == 0 && noxflags.GameFlag(flags)&(noxflags.GameHost|noxflags.GameFlag22) != 0 {
								resolved := name
								if int16(version) < 2 {
									resolved = "arena.map:" + name
								} else if version == 2 {
									resolved = "arena:" + name
								}
								groups := 2
								if kind > 3 && items > 0 {
									groups = 1
								}
								for group := groups - 1; group >= 0; group-- {
									item := prefabGroupRecord{ID: uint32(101 + group), Kind: kind, Name: resolved}
									if kind <= 3 {
										for i := items - 1; i >= 0; i-- {
											v := [2]uint32{uint32(7 + 10*group + i), 0}
											if kind == 2 {
												v[1] = uint32(70 + 10*group + i)
											}
											item.Items = append(item.Items, v)
										}
									}
									expected = append(expected, item)
								}
							}
							actual := r.Groups
							if flags&0x400000 != 0 {
								actual = r.Refs
								if len(r.Groups) != 0 {
									t.Fatal("prefab groups leaked into world list")
								}
							} else if len(r.Refs) != 0 {
								t.Fatal("world groups leaked into prefab list")
							}
							if !reflect.DeepEqual(actual, expected) {
								t.Fatalf("%s group state got %+v want %+v", r.Name, actual, expected)
							}
							rows = append(rows, r)
							cryptfile.Close()
						}
					}
				}
			}
		}
	}
	spellbookCapture(t, "prefab-runtime-group-read", rows, "2d32f5e357d1c8225a40d321b87ba711684ebc563c276be5908d4ea9cfb49856")
}
