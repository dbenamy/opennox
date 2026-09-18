//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestMapMetadataInfo(t *testing.T) {
	owned := serverConfigOwnBytes(t, 0x973F18, 2400, 1478)
	defaults := serverConfigOwnBytes(t, 0x5D4594, 741376, 8)
	copy(defaults, []byte{65, 11, 12, 13, 90, 21, 22, 23})
	dir := t.TempDir()
	type record struct {
		Name  string
		IO    legacy.PortTestMapSectionWire
		State []byte
	}
	var records []record
	for _, read := range []bool{true, false} {
		for _, ver := range []uint16{0, 1, 2, 3, 4, 32768, 65535} {
			if !read && ver != 3 {
				continue
			}
			for _, lens := range [][2]int{{0, 31}, {1, 15}, {15, 1}, {31, 0}} {
				for salt := 0; salt < 3; salt++ {
					name := fmt.Sprintf("read=%t/v=%d/length=%d,%d/salt=%d", read, ver, lens[0], lens[1], salt)
					t.Run(name, func(t *testing.T) {
						for i := range owned {
							owned[i] = byte(1 + (i*17+salt*31)%254)
						}
						header := owned[8:1470]
						// Existing C strlen must see valid, bounded destination strings too.
						for j, off := range []int{1398, 1430} {
							header[off+lens[j]] = 0
						}
						want := bytes.Clone(owned)
						wire := binary.LittleEndian.AppendUint16(nil, ver)
						if ver >= 1 && ver <= 3 {
							fixed := bytes.Clone(header[:1396])
							if read {
								for i := range fixed {
									fixed[i] = byte(i*29 + salt*67)
								}
								copy(want[8:1404], fixed)
							}
							wire = append(wire, fixed...)
							if ver == 2 {
								wire = append(wire, 37, 219)
								want[1404], want[1405] = 37, 219
							} else if read {
								want[1404], want[1405] = 2, 16
							}
						}
						if ver == 3 {
							for j, off := range []int{1398, 1430} {
								s := bytes.Clone(header[off : off+lens[j]])
								if read {
									s = bytes.Repeat([]byte{byte('a' + salt + j)}, lens[j])
									copy(want[8+off:], s)
									want[8+off+len(s)] = 0
								}
								wire = append(wire, byte(len(s)))
								wire = append(wire, s...)
							}
						} else if ver < 3 {
							want[1406], want[1438] = 65, 90
						}
						ret := uint32(1)
						pos := int64(len(wire))
						if ver > 3 {
							ret = 0
							wire = append(wire, 0x91, 0x92, 0x93)
							pos = 2
						}
						result := legacy.PortTestMapMetadata(legacy.PortTestMapSectionIO{Function: "info", Read: read, Data: wire}, dir)
						if result.Return != ret || result.Position != pos || !bytes.Equal(result.Data, wire) {
							t.Fatalf("IO got %+v, expected return %d position %d bytes %x", result, ret, pos, wire)
						}
						if !bytes.Equal(owned, want) {
							for i := range want {
								if owned[i] != want[i] {
									t.Fatalf("state byte %d got %d want %d", i, owned[i], want[i])
								}
							}
						}
						if !bytes.Equal(defaults, []byte{65, 11, 12, 13, 90, 21, 22, 23}) {
							t.Fatal("default backing mutated")
						}
						records = append(records, record{name, result, bytes.Clone(owned)})
					})
				}
			}
		}
	}
	spellbookCapture(t, "map-metadata-info", records, "64ae983faeb2eb437f161268ef2a4c1ad21d2de6fd25569c863256169a8214d3")
}
