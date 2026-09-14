//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

// This is a synthetic empty prefab record, not redistributed game data.
// The real loader reads its header, bounds and section terminator from disk.
func TestMapPopulationPrefabFile(t *testing.T) {
	handles.Init()
	defer handles.Release()
	t.Chdir(t.TempDir())
	var b bytes.Buffer
	word := func(v uint32) {
		if err := binary.Write(&b, binary.LittleEndian, v); err != nil {
			t.Fatal(err)
		}
	}
	word(0) // record byte count is not consulted when reading a selected entry
	b.WriteByte(byte(len("fixture")))
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
	b.WriteByte(126) // zero section-name length in the file's XOR section mode
	if err := os.WriteFile("population-fixture.lib", b.Bytes(), 0600); err != nil {
		t.Fatal(err)
	}
	s := populationCacheBase()
	s.Globals["decodedCache"] = roomValue(0)
	s.Records = append(s.Records, roomRecord(128))
	paintString(&s.Records[8], 0, "population-fixture.lib")
	s.Globals["dword_5d4594_1599588"] = roomArg(9)
	s.Globals["dword_5d4594_1599592"] = roomArg(9)
	s.Globals["dword_5d4594_1599480"] = roomValue(-1)
	s.Actions = []legacy.PortTestPaintAction{paintAction(29, roomArg(1), roomArg(6))}
	out := populationCapture(t, "prefab-file", []legacy.PortTestPaintSpec{s})
	step := out[0].Steps[0]
	if step.Return != 1 || step.Globals["dword_5d4594_1599480"] != 0 {
		t.Fatal("real prefab file did not load")
	}
	if step.Globals["prefabLoadCount"] != 0 {
		t.Fatal("real file test used seeded loader")
	}
	for i, want := range []uint32{46, 0, 0, 46, 92, 46, 46, 92} {
		if step.Globals[fmt.Sprintf("cacheBlob%d", 1599500+4*i)] != want {
			t.Fatalf("decoded bounds word %d differs", i)
		}
	}
}
