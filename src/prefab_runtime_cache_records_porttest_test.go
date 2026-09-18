//go:build porttest

package opennox

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestPrefabRuntimeCacheRecords(t *testing.T) {
	words, restore := legacy.PortTestPrefabRuntimeGlobals()
	defer restore()
	type record struct{ Node, Payload []uint32 }
	type row struct {
		Name    string
		Records []record
		Lookups []uint32
		Count   uint32
	}
	var rows []row
	coords := []int32{0, 1, -1, 127, 128, 255, 256, -257, 32767, 10737416, 16777217, math.MinInt32, math.MaxInt32}
	for _, kind := range []string{"tiles", "walls", "waypoints"} {
		orientations := []uint32{0}
		if kind == "tiles" {
			orientations = []uint32{0, 1, 2, 255, 0x3f800000, 0x3f800001, 0x80000001, 0x7fc00001}
		}
		for _, x := range coords {
			for _, orientation := range orientations {
				*words[kind] = 0
				count := memmap.PtrUint32(0x5D4594, 1599560)
				*count = 0
				var nodes, payloads []unsafe.Pointer
				ids := map[uint32]uint32{0: 0}
				norm := func(v uint32) uint32 {
					n, ok := ids[v]
					if !ok {
						t.Fatalf("unowned cache pointer %x", v)
					}
					return n
				}
				for i := int32(0); i < 3; i++ {
					xx, yy := x+i, -x-i
					var raw uint64
					switch kind {
					case "tiles":
						raw = legacy.PortTestPrefabCall(17, [6]uint32{uint32(xx), uint32(yy), orientation})
					case "walls":
						raw = legacy.PortTestPrefabCall(19, [6]uint32{uint32(xx), uint32(yy)})
					case "waypoints":
						raw = legacy.PortTestPrefabCall(22, [6]uint32{uint32(700 + i), math.Float32bits(float32(xx)), math.Float32bits(float32(yy))})
					}
					if raw == 0 {
						t.Fatal("cache constructor failed")
					}
					node := unsafe.Pointer(uintptr(uint32(raw)))
					payload := *(*unsafe.Pointer)(node)
					nodes = append(nodes, node)
					payloads = append(payloads, payload)
					ids[uint32(raw)] = uint32(100 + i)
					ids[uint32(uintptr(payload))] = uint32(200 + i)
				}
				r := row{Name: fmt.Sprintf("%s/x%d/o%08x", kind, x, orientation), Count: *count}
				if norm(*words[kind]) != 102 {
					t.Fatal("cache list head")
				}
				for i := 0; i < 3; i++ {
					nWords, pWords := 3, 9
					nextOff, prevOff := 1, 2
					if kind == "tiles" {
						nWords, pWords, nextOff, prevOff = 6, 5, 4, 5
					} else if kind == "waypoints" {
						pWords = 129
					}
					n := append([]uint32(nil), unsafe.Slice((*uint32)(nodes[i]), nWords)...)
					p := append([]uint32(nil), unsafe.Slice((*uint32)(payloads[i]), pWords)...)
					n[0], n[nextOff], n[prevOff] = norm(n[0]), norm(n[nextOff]), norm(n[prevOff])
					wantN := make([]uint32, nWords)
					wantP := make([]uint32, pWords)
					wantN[0] = uint32(200 + i)
					if i > 0 {
						wantN[nextOff] = uint32(100 + i - 1)
					}
					if i < 2 {
						wantN[prevOff] = uint32(100 + i + 1)
					}
					xx, yy := x+int32(i), -x-int32(i)
					switch kind {
					case "tiles":
						px := float64(xx) * 46
						py := float64(float32(float64(yy) * 46))
						if byte(orientation) == 1 {
							px += 23
						} else {
							py += 23
						}
						wantN[1], wantN[2], wantN[3] = math.Float32bits(float32(px)), math.Float32bits(float32(py)), uint32(byte(orientation))
					case "walls":
						wantP[1] = uint32(byte(xx))<<8 | uint32(byte(yy))<<16
					case "waypoints":
						p[121], p[122] = norm(p[121]), norm(p[122])
						wantP[0], wantP[2], wantP[3], wantP[120] = uint32(700+i), math.Float32bits(float32(xx)), math.Float32bits(float32(yy)), 0x1000000
						if i > 0 {
							wantP[121] = uint32(200 + i - 1)
						}
						if i < 2 {
							wantP[122] = uint32(200 + i + 1)
						}
					}
					if !reflect.DeepEqual(n, wantN) || !reflect.DeepEqual(p, wantP) {
						t.Fatalf("%s/node%d record mismatch node=%x want=%x", r.Name, i, n, wantN)
					}
					r.Records = append(r.Records, record{n, p})
					if kind == "walls" {
						got := uint32(legacy.PortTestPrefabCall(20, [6]uint32{uint32(byte(xx)), uint32(byte(yy))}))
						if norm(got) != uint32(100+i) {
							t.Fatal("wall lookup byte coordinates")
						}
						exact := uint32(legacy.PortTestPrefabCall(20, [6]uint32{uint32(xx), uint32(yy)}))
						want := uint32(0)
						if xx >= 0 && xx < 256 && yy >= 0 && yy < 256 {
							want = uint32(100 + i)
						}
						if norm(exact) != want {
							t.Fatal("wall lookup full-width coordinates")
						}
						r.Lookups = append(r.Lookups, norm(got), norm(exact))
					}
				}
				if kind == "tiles" && r.Count != 3 {
					t.Fatal("tile count")
				}
				rows = append(rows, r)
				*words[kind] = 0
				for i, p := range payloads {
					if kind == "waypoints" {
						alloc.FreePtr(p)
					} else {
						legacy.PortTestPrefabReleasePayload(p)
					}
					legacy.MapPrefabFreeNode(nodes[i])
				}
			}
		}
	}
	spellbookCapture(t, "prefab-runtime-cache-records", rows, "0770aaf6b42c27a54a3a820cf16dcc3fd9d12e39dac81d01bba808e7e11338e6")
}
