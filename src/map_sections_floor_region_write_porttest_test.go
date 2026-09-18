//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestMapSectionsFloorRegionWrite(t *testing.T) {
	var cases []legacy.PortTestMapSectionSpec
	var wants [][]byte
	node := [4]uint32{37, 0xffff8001, 19, 23}
	for _, point := range [][2]int{{12, 10}, {11, 11}, {14, 14}, {10, 10}} {
		for _, present := range []bool{false, true} {
			quad := [8]uint32{12 * 23, 10 * 23, 10 * 23, 12 * 23, 14 * 23, 12 * 23, 12 * 23, 14 * 23}
			qw := map[int]uint32{}
			for i, v := range quad {
				qw[i*4] = v
			}
			x, y := point[0]/2, point[1]/2
			half, base := uint32(2), 6
			if point[1]&1 != 0 {
				x = (point[0] - 1) / 2
				y = (point[1] + 1) / 2
				half = 1
				base = 1
			}
			words := map[int]uint32{0: 0}
			if present {
				words[0] = half
			}
			for i, v := range node {
				words[4*(base+i)] = v
			}
			sp := legacy.PortTestMapSectionSpec{Name: fmt.Sprintf("point%v/present%v", point, present), Paint: legacy.PortTestPaintSpec{Seed: 67, Records: []legacy.PortTestMapRoomRecord{{Size: 32, Words: qw}}, Cells: []legacy.PortTestPaintCell{{X: x, Y: y, Words: words}}, Actions: []legacy.PortTestPaintAction{{Op: 0, Args: [6]legacy.PortTestMapRoomArg{{Slot: 1}}}}}, IO: []legacy.PortTestMapSectionIO{{Function: "floor"}}}
			wire := []byte{4, 0}
			for _, v := range []uint32{10, 10, 5, 5} {
				wire = binary.LittleEndian.AppendUint32(wire, v)
			}
			if present && (point[0] == 12 || point[0] == 11) {
				wire = binary.LittleEndian.AppendUint16(wire, uint16(point[0]<<8|point[1]))
				wire = append(wire, mapSectionTileBytes([][4]uint32{node})...)
			}
			wire = append(wire, 255, 255)
			cases = append(cases, sp)
			wants = append(wants, wire)
		}
	}
	out := mapSectionsRun(t, cases)
	for i, r := range out {
		t.Run(r.Name, func(t *testing.T) {
			got := r.IO[0]
			if got.Return != 1 || got.Position != int64(len(wants[i])) || !bytes.Equal(got.Data, wants[i]) {
				t.Fatalf("floor region bytes %x want%x", got.Data, wants[i])
			}
		})
	}
	spellbookCapture(t, "map-sections-floor-region-write", out, "39e55c1ff2e304b58ca09a3899229a4fbfe43ab5fa5a7b20e5ae83efb6e4df1a")
}
