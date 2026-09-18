//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
)

func questProgressWire(version uint16, names []string, kinds, values []uint32) []byte {
	var b bytes.Buffer
	binary.Write(&b, binary.LittleEndian, version)
	binary.Write(&b, binary.LittleEndian, uint32(len(names)))
	for i, name := range names {
		b.WriteByte(byte(len(name)))
		b.WriteString(name)
		binary.Write(&b, binary.LittleEndian, kinds[i])
		if kinds[i] <= 1 {
			binary.Write(&b, binary.LittleEndian, values[i])
		}
	}
	return b.Bytes()
}
func TestQuestProgressPersistence(t *testing.T) {
	t.Cleanup(legacy.PortTestQuestProgressOwner())
	defer cryptfile.Close()
	path := filepath.Join(t.TempDir(), "quest.bin")
	type row struct {
		Name     string
		Return   uint64
		Data     []byte
		Records  [][37]uint32
		Checksum uint32
		Position int64
	}
	var rows []row
	for _, flags := range []uint32{0, 2048, 4096, 6144} {
		restore := noxflags.PortTestGameFlags(noxflags.GameFlag(flags))
		for _, count := range []int{0, 1, 4, 32} {
			legacy.PortTestQuestProgress("reset", "*:*", 0)
			legacy.PortTestQuestProgress("namespace", "map", 0)
			var names []string
			var kinds, values []uint32
			for i := 0; i < count; i++ {
				name := fmt.Sprintf("map:key%d", i)
				v := uint32(i) * 0x1234567
				op := "set-int"
				if i%2 == 1 {
					op = "set-float"
				}
				legacy.PortTestQuestProgress(op, name, v)
				names = append([]string{name}, names...)
				kinds = append([]uint32{uint32(i % 2)}, kinds...)
				values = append([]uint32{v}, values...)
			}
			if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
				t.Fatal(err)
			}
			ret := legacy.PortTestQuestProgress("write", "", 0)
			crc := cryptfile.Global().PortTestChecksum()
			cryptfile.Close()
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			want := questProgressWire(1, names, kinds, values)
			if flags&2048 == 0 {
				want = questProgressWire(1, nil, nil, nil)
			}
			if ret != 1 || !bytes.Equal(data, want) {
				t.Fatalf("writer flags=%d count=%d ret=%d bytes=%x expected=%x", flags, count, ret, data, want)
			}
			rows = append(rows, row{fmt.Sprintf("write-%d-%d", flags, count), ret, data, legacy.PortTestQuestProgressSnapshot(), crc, int64(len(data))})
			if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
				t.Fatal(err)
			}
			ret = legacy.PortTestQuestProgress("read", "", 0)
			pos, _ := cryptfile.Global().File.Seek(0, 1)
			crc = cryptfile.Global().PortTestChecksum()
			cryptfile.Close()
			got := legacy.PortTestQuestProgressSnapshot()
			expected := count
			if flags&2048 == 0 {
				expected = 0
			}
			if ret != 1 || len(got) != expected {
				t.Fatal("roundtrip", ret, len(got), expected)
			}
			rows = append(rows, row{fmt.Sprintf("read-%d-%d", flags, count), ret, data, got, crc, pos})
		}
		restore()
	}
	for _, version := range []uint16{0, 1, 2, 0x7fff, 0x8000, 0xffff} {
		data := questProgressWire(version, []string{"map:a", "map:b", "map:a", "map:ignored", "map:float"}, []uint32{0, 1, 1, 7, 1}, []uint32{5, 0x3fc00000, 0xffffffff, 0, 0x80000000})
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
		legacy.PortTestQuestProgress("set-int", "stale:entry", 99)
		if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
			t.Fatal(err)
		}
		ret := legacy.PortTestQuestProgress("read", "", 0)
		pos, _ := cryptfile.Global().File.Seek(0, 1)
		crc := cryptfile.Global().PortTestChecksum()
		cryptfile.Close()
		snap := legacy.PortTestQuestProgressSnapshot()
		if int16(version) > 1 {
			if ret != 0 || len(snap) != 0 || pos != 2 {
				t.Fatal("unsupported version clear/position")
			}
		} else {
			if ret != 1 || len(snap) != 3 || legacy.PortTestQuestProgress("int", "map:a", 0) != 0xffffffff {
				t.Fatal("duplicate/unknown-kind contract")
			}
		}
		rows = append(rows, row{fmt.Sprintf("version-%d", version), ret, data, snap, crc, pos})
	}
	spellbookCapture(t, "quest-progress-persistence", rows, "874350d314bbeb0804cd169998c252a0f5b2f8760f078679e74590b634044180")
}
