//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
)

func TestPrefabRuntimeGroupWrite(t *testing.T) {
	s := prefabGroupOwner(t)
	t.Cleanup(noxflags.PortTestGameFlags(noxflags.GameHost))
	oldFile := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	t.Cleanup(func() { cryptfile.Close(); cryptfile.SetGlobal(oldFile) })
	path := filepath.Join(t.TempDir(), "group-out.bin")
	type row struct {
		Name   string
		Return uint64
		Wire   []byte
	}
	var rows []row
	for kind := byte(0); kind < 4; kind++ {
		for _, groups := range []int{0, 1, 3} {
			for _, items := range []int{0, 1, 3} {
				for _, length := range []int{0, 1, 75} {
					prefabGroupReset(s)
					for g := 0; g < groups; g++ {
						name := strings.Repeat(string(rune('a'+g)), length)
						if s.MapGroups.MapLoadAddGroup57C0C0(name, uint32(100+g), kind) != 1 {
							t.Fatal("group allocation")
						}
						for i := 0; i < items; i++ {
							if s.MapGroups.Sub57C130([]uint32{uint32(10*g + i), uint32(1000 + 10*g + i)}, uint32(100+g)) != 1 {
								t.Fatal("group item allocation")
							}
						}
					}
					if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
						t.Fatal(err)
					}
					ret := legacy.PortTestPrefabCall(33, [6]uint32{})
					if err := cryptfile.Close(); err != nil {
						t.Fatal(err)
					}
					got, err := os.ReadFile(path)
					if err != nil {
						t.Fatal(err)
					}
					var want bytes.Buffer
					binary.Write(&want, binary.LittleEndian, uint16(3))
					binary.Write(&want, binary.LittleEndian, uint32(groups))
					for g := groups - 1; g >= 0; g-- {
						name := strings.Repeat(string(rune('a'+g)), length)
						want.WriteByte(byte(len(name) + 1))
						want.WriteString(name)
						want.WriteByte(0)
						want.WriteByte(kind)
						binary.Write(&want, binary.LittleEndian, uint32(100+g))
						binary.Write(&want, binary.LittleEndian, uint32(items))
						for i := items - 1; i >= 0; i-- {
							binary.Write(&want, binary.LittleEndian, uint32(10*g+i))
							if kind == 2 {
								binary.Write(&want, binary.LittleEndian, uint32(1000+10*g+i))
							}
						}
					}
					if ret != 1 || !bytes.Equal(got, want.Bytes()) {
						t.Fatalf("k%d/g%d/i%d/n%d write differs", kind, groups, items, length)
					}
					rows = append(rows, row{fmt.Sprintf("k%d/g%d/i%d/n%d", kind, groups, items, length), ret, got})
				}
			}
		}
	}
	spellbookCapture(t, "prefab-runtime-group-write", rows, "a5bcbf34b3562ff77dddacc304e4c99aa08f3dbc17618da65059f61a46d18eb3")
}
