//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"math"
	"os"
	"testing"

	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

func TestPrefabRuntimeNonemptyFile(t *testing.T) {
	handles.Init()
	defer handles.Release()
	oldFile := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	defer func() { cryptfile.Close(); cryptfile.SetGlobal(oldFile) }()
	t.Chdir(t.TempDir())
	for _, tc := range []struct {
		name    string
		version uint16
		want    uint32
	}{
		{"decoded", 1, 0}, {"rejected", 0, math.MaxUint32},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var b bytes.Buffer
			word := func(v uint32) {
				if err := binary.Write(&b, binary.LittleEndian, v); err != nil {
					t.Fatal(err)
				}
			}
			word(0)
			b.WriteByte(7)
			b.WriteString("fixture")
			b.WriteByte(1)
			b.WriteByte(1)
			word(math.Float32bits(130.10765))
			word(math.Float32bits(260.2153))
			word(0xcafedead)
			word(32)
			word(32)
			for _, v := range []uint32{46, 0, 46, 92, 0, 46, 92, 46} {
				word(v)
			}
			var sections bytes.Buffer
			name := "DebugData\x00"
			sections.WriteByte(byte(len(name)))
			sections.WriteString(name)
			data := prefabDebugWire(tc.version, [][2]string{{"fixture-room", "north"}, {"fixture-room", "south"}})
			binary.Write(&sections, binary.LittleEndian, uint32(len(data)))
			sections.Write(data)
			sections.WriteByte(0)
			for _, v := range sections.Bytes() {
				b.WriteByte(v ^ 126)
			}
			if err := os.WriteFile("prefab-runtime.lib", b.Bytes(), 0600); err != nil {
				t.Fatal(err)
			}
			s := populationCacheBase()
			s.Globals["decodedCache"] = roomValue(0)
			s.Records = append(s.Records, roomRecord(128))
			paintString(&s.Records[8], 0, "prefab-runtime.lib")
			s.Globals["dword_5d4594_1599588"] = roomArg(9)
			s.Globals["dword_5d4594_1599592"] = roomArg(9)
			s.Globals["dword_5d4594_1599480"] = roomValue(-1)
			s.Actions = []legacy.PortTestPaintAction{paintAction(29, roomArg(1), roomArg(6))}
			out := populationRun([]legacy.PortTestPaintSpec{s})
			if !out[0].Intact || !out[0].ControlOK {
				t.Fatal("fixture state or guard changed")
			}
			step := out[0].Steps[0]
			if step.Globals["dword_5d4594_1599480"] != tc.want {
				t.Fatalf("loaded cache=%d want %d", step.Globals["dword_5d4594_1599480"], tc.want)
			}
			if step.Globals["prefabLoadCount"] != 0 {
				t.Fatal("used seeded loader")
			}
			if cryptfile.Global() != nil {
				t.Fatal("prefab file remained open")
			}
		})
	}
}
