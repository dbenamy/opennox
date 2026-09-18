//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"reflect"
	"testing"
)

func statisticsRecordExpected(f legacy.PortTestStatisticsField) []byte {
	switch f.Kind {
	case 2:
		return []byte{byte(f.Value)}
	case 3:
		b := make([]byte, 2)
		binary.BigEndian.PutUint16(b, uint16(int8(f.Value)))
		return b
	case 6:
		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, uint32(int32(int8(f.Value))))
		return b
	case 7:
		b := f.Data
		if i := bytes.IndexByte(b, 0); i >= 0 {
			b = b[:i]
		}
		out := append(append([]byte{}, b...), 0)
		return out[:uint16(len(out))]
	case 20:
		return f.Data[:uint16(len(f.Data))]
	default:
		panic(f.Kind)
	}
}
func statisticsRecordsCheck(t *testing.T, p *legacy.PortTestStatisticsRecords, kind uint16, fields []legacy.PortTestStatisticsField) []byte {
	t.Helper()
	before := p.Snapshot()
	got := p.Serialize()
	if !reflect.DeepEqual(before, p.Snapshot()) {
		t.Fatal("serialization changed host records")
	}
	if !bytes.Equal(got, p.Serialize()) {
		t.Fatal("serialization is not repeatable")
	}
	want := []byte{0, 0, byte(kind >> 8), byte(kind)}
	for i := len(fields) - 1; i >= 0; i-- {
		f := fields[i]
		data := statisticsRecordExpected(f)
		var tag [4]byte
		name := f.Tag
		if end := bytes.IndexByte([]byte(name), 0); end >= 0 {
			name = name[:end]
		}
		copy(tag[:], name)
		want = append(want, tag[:]...)
		want = append(want, byte(f.Kind>>8), byte(f.Kind), byte(len(data)>>8), byte(len(data)))
		want = append(want, data...)
		for len(want)%4 != 0 {
			want = append(want, 0)
		}
	}
	binary.BigEndian.PutUint16(want, uint16(len(want)))
	if !bytes.Equal(got, want) {
		t.Fatalf("record encoding mismatch: got %x want %x", got, want)
	}
	for i, n := range before {
		if n.Next != i-1 {
			t.Fatalf("node %d next %d", i, n.Next)
		}
	}
	return got
}
func TestGameStatisticsRecordWidths(t *testing.T) {
	// The legacy scalar API accepts signed char even for two/four-byte fields.
	type row struct {
		Path, Kind int
		Value      int32
		Bytes      []byte
	}
	var rows []row
	for path := 0; path < 3; path++ {
		for _, kind := range []int{2, 3, 6} {
			for _, v := range []int32{-2147483648, -65537, -257, -256, -129, -128, -1, 0, 1, 127, 128, 255, 256, 257, 65535, 2147483647} {
				t.Run(fmt.Sprintf("path%d/kind%d/value%d", path, kind, v), func(t *testing.T) {
					p := legacy.PortTestStatisticsRecordsOpen(0x8765)
					defer p.Close()
					f := legacy.PortTestStatisticsField{Kind: kind, Tag: "VAL?", Value: v}
					p.Add(f, path)
					rows = append(rows, row{path, kind, v, statisticsRecordsCheck(t, p, 0x8765, []legacy.PortTestStatisticsField{f})})
				})
			}
		}
	}
	spellbookCapture(t, "game-statistics-record-widths", rows, "2cb82be714e41e7b4bd372903bfeb9b9ca9f646144f848c37e1f2286c28aece7")
}
func TestGameStatisticsRecordOrdering(t *testing.T) {
	type row struct {
		Path, Count int
		Bytes       []byte
	}
	var rows []row
	for path := 0; path < 3; path++ {
		p := legacy.PortTestStatisticsRecordsOpen(0xabcd)
		func() {
			defer p.Close()
			var fields []legacy.PortTestStatisticsField
			rows = append(rows, row{path, 0, statisticsRecordsCheck(t, p, 0xabcd, fields)})
			for i, tag := range []string{"", "A", "AB", "ABC", "ABCD", "ABCDE", "A\x00BC", "same", "same"} {
				for _, kind := range []int{2, 3, 6, 7, 20} {
					data := make([]byte, i)
					for j := range data {
						data[j] = byte(0x80 + j)
					}
					f := legacy.PortTestStatisticsField{Kind: kind, Tag: tag, Value: int32(i * 57), Data: data}
					fields = append(fields, f)
					p.Add(f, path)
					rows = append(rows, row{path, len(fields), statisticsRecordsCheck(t, p, 0xabcd, fields)})
				}
			}
		}()
	}
	spellbookCapture(t, "game-statistics-record-ordering", rows, "13ce0401c6862e229ef397d1be7620143e95b7f9a1f0faa8a0b4aef20f36a978")
}
func TestGameStatisticsRecordReplacement(t *testing.T) {
	p := legacy.PortTestStatisticsRecordsOpen(7)
	defer p.Close()
	p.Add(legacy.PortTestStatisticsField{Kind: 7, Tag: "OLD", Data: []byte("previous storage")}, 0)
	var rows [][]byte
	for round := 0; round < 4; round++ {
		for _, kind := range []int{20, 7, 6, 3, 2} {
			for _, size := range []int{0, 1, 2, 3, 4, 7, 256, 1024} {
				f := legacy.PortTestStatisticsField{Kind: kind, Tag: "NEW", Value: int32(size), Data: bytes.Repeat([]byte{byte(1 + round)}, size)}
				p.Replace(0, f)
				rows = append(rows, statisticsRecordsCheck(t, p, 7, []legacy.PortTestStatisticsField{f}))
			}
		}
	}
	spellbookCapture(t, "game-statistics-record-replacement", rows, "720e6933c9377cc7f6844131eabdc48a1b23825a2dd667d82aaa75cb4352eac2")
}

func TestGameStatisticsRecordEndian(t *testing.T) {
	type row struct {
		Kind           int
		Value          int32
		Encode, Decode uint16
		Bytes          []byte
	}
	var rows []row
	for _, kind := range []int{2, 3, 6, 7, 20} {
		for _, value := range []int32{0, 1, 127, 128, 255, -1} {
			p := legacy.PortTestStatisticsRecordsOpen(0)
			func() {
				defer p.Close()
				f := legacy.PortTestStatisticsField{Kind: kind, Tag: "ENDN", Value: value, Data: []byte{1, 2, 3, 4, 5}}
				p.Add(f, 0)
				before := p.Snapshot()
				enc := p.Endian(0, true)
				dec := p.Endian(0, false)
				length := before[0].Length
				if enc != length>>8|length<<8 {
					t.Fatalf("encode return %x, length %x", enc, length)
				}
				want := uint16(kind - 3)
				if kind == 3 || kind == 6 {
					want = uint16(int8(value))
				}
				if dec != want || !reflect.DeepEqual(before, p.Snapshot()) {
					t.Fatalf("endian roundtrip kind%d value%d: %x want%x", kind, value, dec, want)
				}
				rows = append(rows, row{kind, value, enc, dec, statisticsRecordsCheck(t, p, 0, []legacy.PortTestStatisticsField{f})})
			}()
		}
	}
	spellbookCapture(t, "game-statistics-record-endian", rows, "9228a19d4751fc8121f03e14082befe164c1df481c41a36249f8caf004e65db2")
}
func TestGameStatisticsRecordLengthBoundaries(t *testing.T) {
	type row struct {
		Kind, Size int
		Hash       string
		Length     int
	}
	var rows []row
	for _, kind := range []int{7, 20} {
		for _, size := range []int{0, 1, 2, 3, 4, 255, 256, 65530, 65531, 65532, 65534, 65535, 65536, 65537} {
			p := legacy.PortTestStatisticsRecordsOpen(0xffff)
			func() {
				defer p.Close()
				f := legacy.PortTestStatisticsField{Kind: kind, Tag: "SIZE", Data: bytes.Repeat([]byte{0xab}, size)}
				p.Add(f, 0)
				got := statisticsRecordsCheck(t, p, 0xffff, []legacy.PortTestStatisticsField{f})
				rows = append(rows, row{kind, size, fmt.Sprintf("%x", sha256.Sum256(got)), len(got)})
			}()
		}
	}
	spellbookCapture(t, "game-statistics-record-length-boundaries", rows, "3e0a1cee6f7d0e903894ce0eb43993d9d4d187ddc02b8ed11794b60790b403a7")
}

func TestGameStatisticsRecordEmbeddedZero(t *testing.T) {
	type row struct {
		Kind, Path   int
		Input, Bytes []byte
	}
	var rows []row
	for _, kind := range []int{7, 20} {
		for path := 0; path < 3; path++ {
			for _, data := range [][]byte{nil, {}, {0}, {0, 0}, {0, 128, 255}, {128, 0, 255}, {128, 255, 0}, {1, 2, 0, 3, 0, 4}} {
				p := legacy.PortTestStatisticsRecordsOpen(0)
				func() {
					defer p.Close()
					f := legacy.PortTestStatisticsField{Kind: kind, Tag: "NULS", Data: data}
					p.Add(f, path)
					rows = append(rows, row{kind, path, data, statisticsRecordsCheck(t, p, 0, []legacy.PortTestStatisticsField{f})})
				}()
			}
		}
	}
	spellbookCapture(t, "game-statistics-record-embedded-zero", rows, "42a0d88e6d7b7059ba3931da5a9c725682d84b4d2809f1c7469c7d1185a67b64")
}
