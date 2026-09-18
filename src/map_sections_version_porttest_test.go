//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestMapSectionsUnsupportedVersions(t *testing.T) {
	var cases []legacy.PortTestMapSectionSpec
	for _, kind := range []struct {
		name    string
		maximum uint16
	}{{"floor", 4}, {"walls", 7}, {"windows", 2}, {"breakable", 1}, {"secret", 2}} {
		versions := []uint16{kind.maximum + 1, 127, 32767}
		if kind.name == "floor" {
			versions = append(versions, 0, 1, 2, 32768, 65535)
		}
		for _, version := range versions {
			cases = append(cases, legacy.PortTestMapSectionSpec{Name: fmt.Sprintf("%s/v%d", kind.name, version), Paint: legacy.PortTestPaintSpec{Seed: 47, Globals: map[string]legacy.PortTestMapRoomArg{"section-magic-wall": {Value: 2}}, Actions: []legacy.PortTestPaintAction{{Op: 0}}}, IO: []legacy.PortTestMapSectionIO{{Function: kind.name, Read: true, Data: binary.LittleEndian.AppendUint16(nil, version)}}})
		}
	}
	out := mapSectionsRun(t, cases)
	for _, r := range out {
		t.Run(r.Name, func(t *testing.T) {
			if r.IO[0].Return != 0 || r.IO[0].Position != 2 {
				t.Fatalf("unsupported version return/position %d/%d", r.IO[0].Return, r.IO[0].Position)
			}
		})
	}
	spellbookCapture(t, "map-sections-version-rejection", out, "64e9e80a9c3a9efde15e8a3dda7ef0453e09ddaab6b0f5b90c79c0936ebcf6ef")
}
