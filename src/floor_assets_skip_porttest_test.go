//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"strings"
	"testing"
)

func TestFloorAssetsSkipping(t *testing.T) {
	o := newFloorAssetsOwner(t)
	type record struct {
		Edge             bool
		Name             string
		Dims             [3]byte
		Mode, Format     int
		End              uint32
		Return, Consumed int
		Scratch          [32]byte
	}
	var rows []record
	for _, edge := range []bool{false, true} {
		dimsList := [][3]byte{{0, 0, 0}, {1, 1, 1}, {2, 3, 2}, {255, 1, 1}, {255, 255, 0}}
		if edge {
			dimsList = [][3]byte{{0, 0, 0}, {1, 1, 1}, {2, 3, 2}, {255, 1, 0}, {0, 255, 255}}
		}
		for _, name := range []string{"", "Example", strings.Repeat("Z", 31)} {
			for _, dims := range dimsList {
				for mode := 0; mode < 3; mode++ {
					for _, end := range []uint32{0x454e4420, 0x12345678} {
						formats := []int{0}
						if edge {
							formats = []int{0, 1, 255}
						}
						for _, format := range formats {
							o.reset()
							*o.count = 7
							*o.edgeCount = 9
							op := 2
							payload := floorAssetsFloor(name, [3]byte{17, 129, 254}, dims, mode)
							payload[len(payload)-4] = byte(end)
							payload[len(payload)-3] = byte(end >> 8)
							payload[len(payload)-2] = byte(end >> 16)
							payload[len(payload)-1] = byte(end >> 24)
							if edge {
								op = 3
								payload = floorAssetsEdge(name, dims, mode, byte(format), end)
							}
							defsBefore, edgesBefore := o.metadata()
							wantRet := 0
							if end == 0x454e4420 {
								wantRet = 1
							}
							wantConsumed := len(payload)
							wantScratch := bytes.Repeat([]byte{0xa5}, 64)
							if edge {
								copy(wantScratch, name)
								if format == 1 {
									wantRet = 0
									wantConsumed = 17 + len(name)
								}
							}
							ret, consumed := o.invoke(t, op, payload)
							defsAfter, edgesAfter := o.metadata()
							o.assertOnlyData(t, false, -1)
							if ret != wantRet || consumed != wantConsumed || !bytes.Equal(o.scratch[:64], wantScratch) || *o.count != 7 || *o.edgeCount != 9 || !bytes.Equal(defsBefore, defsAfter) || !bytes.Equal(edgesBefore, edgesAfter) {
								t.Fatalf("skip edge%v name%q dims%v mode%d format%d end%x ret%d/%d consumed%d/%d", edge, name, dims, mode, format, end, ret, wantRet, consumed, wantConsumed)
							}
							rows = append(rows, record{edge, name, dims, mode, format, end, ret, consumed, sha256.Sum256(o.scratch)})
						}
					}
				}
			}
		}
	}
	floorAssetsCapture(t, "skipping", rows, "0ed871e6198beaeedef4ba84a056545cb216308a3c3af0d15dd7d3bf43ddebc7")
}
