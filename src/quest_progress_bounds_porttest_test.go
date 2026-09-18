//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
)

func TestQuestProgressBoundedInput(t *testing.T) {
	t.Cleanup(legacy.PortTestQuestProgressOwner())
	defer cryptfile.Close()
	legacy.PortTestQuestProgress("namespace", "map", 0)
	legacy.PortTestQuestProgress("set-int", "keep", 17)
	before := legacy.PortTestQuestProgressSnapshot()
	for _, name := range []string{":" + strings.Repeat("x", 131), strings.Repeat("x", 128), strings.Repeat("x", 255), strings.Repeat("x", 4096)} {
		for _, op := range []string{"set-int", "set-float", "find", "int", "float", "reset"} {
			legacy.PortTestQuestProgress(op, name, 0xffffffff)
			if !reflect.DeepEqual(before, legacy.PortTestQuestProgressSnapshot()) {
				t.Fatal("oversized name changed accepted records", op, len(name))
			}
		}
	}
	legacy.PortTestQuestProgress("namespace", strings.Repeat("n", 132), 0)
	if legacy.PortTestQuestProgress("int", "keep", 0) != 17 {
		t.Fatal("oversized namespace replaced prior one")
	}
	path := filepath.Join(t.TempDir(), "quest.bin")
	run := func(data []byte) uint64 {
		t.Helper()
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
		if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
			t.Fatal(err)
		}
		ret := legacy.PortTestQuestProgress("read", "", 0)
		cryptfile.Close()
		return ret
	}
	valid := questProgressWire(1, []string{"map:first", "map:second"}, []uint32{0, 1}, []uint32{17, 0x3fc00000})
	for end := 0; end < len(valid); end++ {
		if run(valid[:end]) != 0 {
			t.Fatal("incomplete save accepted", end)
		}
		snap := legacy.PortTestQuestProgressSnapshot()
		if len(snap) > 1 {
			t.Fatal("partial save fabricated records", end)
		}
	}
	if run(valid) != 1 || len(legacy.PortTestQuestProgressSnapshot()) != 2 {
		t.Fatal("valid save")
	}
	for _, size := range []int{132, 255} {
		data := questProgressWire(1, []string{"map:first", strings.Repeat("x", size)}, []uint32{0, 0}, []uint32{17, 99})
		if run(data) != 0 || len(legacy.PortTestQuestProgressSnapshot()) != 1 || legacy.PortTestQuestProgress("int", "map:first", 0) != 17 {
			t.Fatal("oversized save name lost valid prefix", size)
		}
	}
	var huge bytes.Buffer
	binary.Write(&huge, binary.LittleEndian, uint16(1))
	binary.Write(&huge, binary.LittleEndian, uint32(0xffffffff))
	if run(huge.Bytes()) != 0 || len(legacy.PortTestQuestProgressSnapshot()) != 0 {
		t.Fatal("unbounded record count did not stop at EOF")
	}
	data := questProgressWire(1, []string{"map:first\x00ignored"}, []uint32{0}, []uint32{73})
	if run(data) != 1 || legacy.PortTestQuestProgress("int", "map:first", 0) != 73 {
		t.Fatal("C string termination changed")
	}
}
