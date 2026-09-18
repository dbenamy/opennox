//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

type statisticsDecodedField struct {
	Tag  string
	Kind uint16
	Data []byte
}

func statisticsDecode(t *testing.T, data []byte) []statisticsDecodedField {
	t.Helper()
	if len(data) < 4 || int(binary.BigEndian.Uint16(data)) != len(data) {
		t.Fatal("report size")
	}
	var out []statisticsDecodedField
	for pos := 4; pos < len(data); {
		if pos+8 > len(data) {
			t.Fatal("short field header")
		}
		size := int(binary.BigEndian.Uint16(data[pos+6:]))
		end := pos + 8 + size
		if end > len(data) {
			t.Fatal("short payload")
		}
		out = append(out, statisticsDecodedField{string(data[pos : pos+4]), binary.BigEndian.Uint16(data[pos+4:]), append([]byte{}, data[pos+8:end]...)})
		pos = end
		for pos%4 != 0 {
			if pos >= len(data) || data[pos] != 0 {
				t.Fatal("nonzero/absent field padding")
			}
			pos++
		}
	}
	return out
}
func statisticsOwnReportTags(t *testing.T) {
	// The report functions mutate the suffixes of shared field-name strings.
	// Preserve the initialized data and observe it separately from the report.
	copy(serverConfigOwnBytes(t, 0x587000, 71480, 536), blobdata.PortTestStatisticsTags())
	serverConfigOwnBytes(t, 0x5D4594, 741660, 32)
}
func TestGameStatisticsReports(t *testing.T) {
	statisticsOwnReportTags(t)
	type row struct {
		Quest                  bool
		Mode, Count            int
		Value, Sequence, After uint32
		Bytes                  []byte
		Tags                   []byte
	}
	var rows []row
	for _, quest := range []bool{false, true} {
		for _, mode := range []int{0, 1, 2, 7} {
			if quest && mode != 0 {
				continue
			}
			for _, count := range []int{0, 1, 2, 9, 16, 32} {
				for _, value := range []uint32{0, 1, 127, 128, 255, 256, 0x80000001, 0xffffffff} {
					for _, sequence := range []uint32{0, 127, 255, 0xffffffff} {
						t.Run(fmt.Sprintf("quest%v/mode%d/count%d/v%x/seq%x", quest, mode, count, value, sequence), func(t *testing.T) {
							raw := make([]byte, 640)
							for off := 0; off < 536; off += 4 {
								binary.LittleEndian.PutUint32(raw[off:], value)
							}
							scene, name := 96, 352
							if quest {
								scene, name = 24, 280
							}
							clear(raw[scene : scene+256])
							copy(raw[scene:], "arena.map")
							clear(raw[name : name+128])
							copy(raw[name:], "Test match")
							*memmap.PtrUint32(0x5D4594, 741668) = sequence
							*memmap.PtrUint32(0x5D4594, 741672) = sequence
							events := []byte{0, 255, 1, 0, 255, 255}
							data := legacy.PortTestStatisticsReport(quest, mode, raw, count, events)
							fields := statisticsDecode(t, data)
							found := map[string][]statisticsDecodedField{}
							for _, f := range fields {
								found[f.Tag] = append(found[f.Tag], f)
								if f.Tag == "\x00\x00\x00\x00" {
									t.Fatal("uninitialized report field name")
								}
							}
							for _, tag := range []string{"IDNO", "GSKU", "GSTY"} {
								f := found[tag]
								if len(f) != 1 || f[0].Kind != 6 || len(f[0].Data) != 4 || binary.BigEndian.Uint32(f[0].Data) != uint32(int32(int8(value))) {
									t.Fatalf("%s scalar %v", tag, f)
								}
							}
							for tag, want := range map[string]string{"SCEN": "arena.map\x00", "GNAM": "Test match\x00"} {
								f := found[tag]
								if len(f) != 1 || f[0].Kind != 7 || string(f[0].Data) != want {
									t.Fatalf("%s string %v", tag, f)
								}
							}
							if quest || mode == 0 || mode == 1 || mode == 2 {
								field := found["PLRS"]
								wantCount := int16(count)
								if !quest && mode == 0 {
									wantCount = -1
								}
								if len(field) != 1 || field[0].Kind != 3 || len(field[0].Data) != 2 || binary.BigEndian.Uint16(field[0].Data) != uint16(int8(wantCount)) {
									t.Fatal("player count field", field, wantCount)
								}
							}
							if quest || mode == 1 || mode == 2 {
								tag, off := "IPL?", 612
								if quest {
									tag, off = "IPLS", 540
								}
								ips := found[tag]
								if len(ips) != count {
									t.Fatal("report player rows", len(ips), count)
								}
								for i, f := range ips {
									want := uint32(int32(int8(uint32(0x12340000 + off + (count-1-i)*131))))
									if f.Kind != 6 || len(f.Data) != 4 || binary.BigEndian.Uint32(f.Data) != want {
										t.Fatal("player IP array order/width", i, f, want)
									}
								}
							}
							after := *memmap.PtrUint32(0x5D4594, 741668)
							wantSequence := sequence
							if mode == 0 {
								wantSequence = 0
							}
							if mode == 1 || mode == 2 {
								wantSequence++
							}
							if quest {
								after = *memmap.PtrUint32(0x5D4594, 741672)
								wantSequence = sequence + 1
							}
							if after != wantSequence {
								t.Fatalf("sequence %x != %x", after, wantSequence)
							}
							rows = append(rows, row{quest, mode, count, value, sequence, after, data, bytes.Clone(unsafe.Slice(memmap.PtrUint8(0x587000, 71480), 536))})
						})
					}
				}
			}
		}
	}
	spellbookCapture(t, "game-statistics-reports", rows, "414e5f59c9f79de28cf4069304c4f246c7b158b9c2a529399e443118af2932f0")
}
