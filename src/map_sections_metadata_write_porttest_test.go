//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestMapSectionsMetadataWrite(t *testing.T) {
	var cases []legacy.PortTestMapSectionSpec
	var wants [][]byte
	for _, kind := range []string{"windows", "breakable", "secret"} {
		bit := byte(0x40)
		version := uint16(2)
		if kind == "breakable" {
			bit = 8
			version = 1
		}
		if kind == "secret" {
			bit = 4
		}
		for _, flag := range []byte{0, bit, 0xff, 0xff ^ bit} {
			for _, region := range []bool{false, true} {
				for _, point := range [][2]int{{12, 10}, {30, 30}} {
					quad := [8]uint32{12 * 23, 8 * 23, 10 * 23, 10 * 23, 14 * 23, 10 * 23, 12 * 23, 12 * 23}
					qw := map[int]uint32{}
					for i, v := range quad {
						qw[i*4] = v
					}
					sp := legacy.PortTestMapSectionSpec{Name: fmt.Sprintf("%s/flag%x/region%v/point%v", kind, flag, region, point), Paint: legacy.PortTestPaintSpec{Seed: 59, Records: []legacy.PortTestMapRoomRecord{{Size: 32, Words: map[int]uint32{16: 0x12345678, 20: 0xa5c30308, 24: 0x87654321, 28: 0xfedcba98}}, {Size: 32, Words: qw}}, Walls: []legacy.PortTestPaintWall{{X: point[0], Y: point[1], Words: map[int]uint32{4: uint32(flag) | uint32(point[0])<<8 | uint32(point[1])<<16}, Data: legacy.PortTestMapRoomArg{Slot: 1}}}, Actions: []legacy.PortTestPaintAction{{Op: 0}}}, IO: []legacy.PortTestMapSectionIO{{Function: kind}}}
					if region {
						sp.Paint.Actions[0].Args[0] = legacy.PortTestMapRoomArg{Slot: 2}
					}
					wire := binary.LittleEndian.AppendUint16(nil, version)
					count := uint16(0)
					// The actual wall iterator omits doors and broken walls before callbacks.
					if flag&bit != 0 && flag&0x30 == 0 && (!region || point[0] == 12) {
						count = 1
					}
					wire = binary.LittleEndian.AppendUint16(wire, count)
					if count != 0 {
						wire = binary.LittleEndian.AppendUint32(wire, uint32(point[0]))
						wire = binary.LittleEndian.AppendUint32(wire, uint32(point[1]))
						if kind == "secret" {
							wire = binary.LittleEndian.AppendUint32(wire, 0x12345678)
							wire = append(wire, 8, 3, 0xc3)
							wire = binary.LittleEndian.AppendUint32(wire, 0x87654321)
							wire = binary.LittleEndian.AppendUint32(wire, 0xfedcba98)
						}
					}
					cases = append(cases, sp)
					wants = append(wants, wire)
				}
			}
		}
	}
	out := mapSectionsRun(t, cases)
	for i, r := range out {
		t.Run(r.Name, func(t *testing.T) {
			got := r.IO[0]
			if got.Return != 1 || got.Position != int64(len(wants[i])) || !bytes.Equal(got.Data, wants[i]) {
				t.Fatalf("metadata serialization result%d position%d bytes%x want%x", got.Return, got.Position, got.Data, wants[i])
			}
		})
	}
	spellbookCapture(t, "map-sections-metadata-write", out, "bb31322318332a0809d1777c596875e47d28911ce16bd3f867b87f145f108488")
}
