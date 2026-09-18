//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"time"
	"unsafe"
)

func TestGameStatisticsMatchArrays(t *testing.T) {
	words, restore := legacy.PortTestStatisticsGlobals()
	t.Cleanup(restore)
	raw, free := alloc.Make([]byte{}, 640)
	defer free()
	record := unsafe.Pointer(&raw[0])
	ptr := func(off int) unsafe.Pointer { return *(*unsafe.Pointer)(unsafe.Add(record, off)) }
	defer func() {
		names := ptr(608)
		if names != nil {
			for i := 0; i < int(*words["array-count"]); i++ {
				legacy.PortTestStatisticsFree(*(*unsafe.Pointer)(unsafe.Add(names, 4*i)))
			}
		}
		for off := 608; off <= 632; off += 4 {
			legacy.PortTestStatisticsFree(ptr(off))
		}
	}()
	type row struct {
		Count  int
		Names  []string
		Words  [][3]uint32
		Flags  [][3]byte
		Source []byte
	}
	var rows []row
	for _, count := range []int{0, 1, 2, 31, 254, 2, 0, 31, 1, 0} {
		source, release := alloc.Make([]byte{}, 32*(count+1))
		for i := 0; i < count; i++ {
			r := source[32*i : 32*(i+1)]
			copy(r, fmt.Sprintf("N%07d", i))
			binary.LittleEndian.PutUint32(r[12:], 0x12340000+uint32(i))
			binary.LittleEndian.PutUint32(r[16:], 0xfedc0000+uint32(i))
			r[20], r[21], r[28] = byte(i*7), byte(i*17), byte(i*43)
			binary.LittleEndian.PutUint32(r[24:], 0x12345678+uint32(i))
		}
		before := uint32(time.Now().Unix())
		legacy.PortTestStatisticsCall("match-arrays", record, unsafe.Pointer(&source[0]), int32(count))
		after := uint32(time.Now().Unix())
		if binary.LittleEndian.Uint16(raw[6:]) != uint16(count) || *words["array-count"] != uint32(count) {
			t.Fatal("array count")
		}
		result := row{Count: count}
		for i := 0; i < count; i++ {
			namePtr := *(*unsafe.Pointer)(unsafe.Add(ptr(608), 4*i))
			name := unsafe.Slice((*byte)(namePtr), 10)
			end := bytes.IndexByte(name, 0)
			if end < 0 || string(name[:end]) != fmt.Sprintf("N%07d", i) {
				t.Fatal("owned player name")
			}
			get := func(off int) uint32 { return *(*uint32)(unsafe.Add(ptr(off), 4*i)) }
			flag := func(off int) byte { return *(*byte)(unsafe.Add(ptr(off), i)) }
			stamp := uint32(0x12345678) + uint32(i)
			duration := get(628)
			if duration < before-stamp || duration > after-stamp || binary.LittleEndian.Uint32(source[32*i+24:]) != duration {
				t.Fatal("elapsed time/source update")
			}
			if get(612) != 0x12340000+uint32(i) || get(616) != 0xfedc0000+uint32(i) || flag(620) != byte(i*7) || flag(624) != byte(i*17) || flag(632) != byte(i*43) {
				t.Fatal("array field copy")
			}
			result.Names = append(result.Names, string(name[:end]))
			result.Words = append(result.Words, [3]uint32{get(612), get(616), 0})
			result.Flags = append(result.Flags, [3]byte{flag(620), flag(624), flag(632)})
			binary.LittleEndian.PutUint32(source[32*i+24:], 0) // Checked wall-clock value only.
		}
		result.Source = bytes.Clone(source)
		release()
		rows = append(rows, result)
	}
	spellbookCapture(t, "game-statistics-match-arrays", rows, "ca15b39a60dfd614175fd2a77b62954ab5182472dd990afaa925d90ea3f7b1b1")
}
func TestGameStatisticsEventArray(t *testing.T) {
	raw, free := alloc.Make([]byte{}, 640)
	defer free()
	record := unsafe.Pointer(&raw[0])
	defer func() { legacy.PortTestStatisticsFree(*(*unsafe.Pointer)(unsafe.Add(record, 636))) }()
	serverConfigOwnBytes(t, 0x5D4594, 741308, 4)
	type row struct {
		Count  int
		Result uint32
		Data   []byte
	}
	var rows []row
	for _, count := range []int{0, 1, 2, 127, 128, 255, 256, 32767, 0, 2, 0} {
		data, release := alloc.Make([]byte{}, 2*count+1)
		for i := range data {
			data[i] = byte(i*71 + 13)
		}
		got := legacy.PortTestStatisticsCall("event-array", record, unsafe.Pointer(&data[0]), int32(count))
		dst := *(*unsafe.Pointer)(unsafe.Add(record, 636))
		output := bytes.Clone(unsafe.Slice((*byte)(dst), 2*count))
		if got != uint32(2*count) || *memmap.PtrUint32(0x5D4594, 741308) != uint32(count) || !bytes.Equal(output, data[:2*count]) {
			t.Fatal("event array copy/count")
		}
		release()
		rows = append(rows, row{count, got, output})
	}
	spellbookCapture(t, "game-statistics-event-array", rows, "03db0e7c93fa8a801ded52b58e080a8e09861b117edcf4051b85b88f0b1841e4")
}
