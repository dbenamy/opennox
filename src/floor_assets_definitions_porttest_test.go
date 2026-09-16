//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"strconv"
	"strings"
	"testing"
	"unsafe"
)

func floorAssetsFacades() []string {
	var names []string
	for i := 0; i < 12; i++ {
		p := *memmap.PtrPtr(0x587000, 26488+4*uintptr(i))
		if p == nil {
			return names
		}
		names = append(names, alloc.GoString((*byte)(p)))
	}
	panic("unterminated facade table")
}
func floorAssetsScratch(mode, count int) []byte {
	b := bytes.Repeat([]byte{0xa5}, 64)
	for i := 0; i < count; i++ {
		b[0] = 0
		if mode == 1 || (mode == 2 && i%2 != 0) {
			name := []byte("outside-" + strconv.Itoa(i))
			copy(b, name)
			b[len(name)] = 0
		}
	}
	return b
}
func TestFloorAssetsDefinitions(t *testing.T) {
	o := newFloorAssetsOwner(t)
	facades := floorAssetsFacades()
	if len(facades) == 0 {
		t.Fatal("actual facade table empty")
	}
	names := []string{"", facades[0], strings.ToLower(facades[0]), strings.Repeat("N", 31), "Unknown", "\x80raw", "\x00tail", "ab\x00tail"}
	type record struct {
		Edge                bool
		Name                string
		Count, Mode, Format int
		Dims                [3]byte
		Color               [3]byte
		Return, Consumed    int
		Counts              [2]uint32
		Metadata            [32]byte
		Scratch             []byte
	}
	var rows []record
	for _, edge := range []bool{false, true} {
		cap := 176
		dimensions := [][3]byte{{0, 0, 0}, {1, 1, 1}, {2, 3, 2}, {255, 1, 1}, {255, 255, 0}}
		if edge {
			cap = 64
			dimensions = [][3]byte{{0, 0, 0}, {1, 1, 1}, {2, 3, 2}, {255, 1, 0}, {0, 255, 255}}
		}
		for _, count := range []int{0, 1, cap - 1, cap} {
			for _, name := range names {
				effectiveName := strings.SplitN(name, "\x00", 2)[0]
				for _, dims := range dimensions {
					for mode := 0; mode < 3; mode++ {
						colors := [][3]byte{{0, 0, 0}, {255, 0, 255}, {17, 129, 254}, {255, 255, 255}}
						if edge {
							colors = colors[:1]
						}
						formats := []int{0}
						if edge {
							formats = []int{0, 1, 255}
						}
						for _, color := range colors {
							for _, format := range formats {
								o.reset()
								index := count
								op := 0
								frames := int(dims[0]) * int(dims[1]) * int(dims[2])
								payload := floorAssetsFloor(name, color, dims, mode)
								if edge {
									*o.edgeCount = uint32(count)
									op = 1
									frames = 2 * int(dims[0]) * (int(dims[1]) + int(dims[2]))
									payload = floorAssetsEdge(name, dims, mode, byte(format), 0x454e4420)
								} else {
									*o.count = uint32(count)
								}
								if count < cap {
									row := unsafe.Slice((*byte)(unsafe.Pointer(&o.defs[count])), 60)
									if edge {
										row = o.edges[count*60:][:60]
									}
									for i := range row {
										row[i] = 0x5a
									}
									clear(row[32:36])
									row[31] = 0
								}
								wantDefs, wantEdges := o.metadata()
								wantRet, wantConsumed := 0, 0
								wantScratch := bytes.Repeat([]byte{0xa5}, 64)
								wantCounts := [2]uint32{*o.count, *o.edgeCount}
								if count < cap {
									target := wantDefs[index*60:][:60]
									if edge {
										target = wantEdges[index*60:][:60]
									}
									if edge {
										copy(target[:32], effectiveName)
										target[len(effectiveName)] = 0
									} else {
										clear(target[:31])
										copy(target[:31], effectiveName)
									}
									binary.LittleEndian.PutUint16(target[46:], 0)
									wantRet = 1
									wantConsumed = len(payload)
									wantScratch = floorAssetsScratch(mode, frames)
									if edge {
										binary.LittleEndian.PutUint32(target[36:], 0xfffffff9)
										binary.LittleEndian.PutUint32(target[40:], 0x10203040)
										target[57] = 9
										target[54] = dims[0]
										target[55] = 7
										target[56] = 7
										if format == 1 {
											wantRet = 0
											wantConsumed = 17 + len(name)
											wantScratch = bytes.Repeat([]byte{0xa5}, 64)
										} else {
											target[53] = dims[1]
											target[52] = dims[2]
											binary.LittleEndian.PutUint16(target[44:], uint16(2*(int(dims[1])+int(dims[2]))))
											wantCounts[1]++
										}
									} else {
										binary.LittleEndian.PutUint32(target[36:], 0xfffffff1)
										binary.LittleEndian.PutUint32(target[40:], 0x12345678)
										binary.LittleEndian.PutUint16(target[44:], uint16(dims[0])*uint16(dims[1]))
										v := nox_color_rgb_4344A0(int(color[0]), int(color[1]), int(color[2]))
										if color == [3]byte{255, 0, 255} {
											v = 0x80000000
										}
										binary.LittleEndian.PutUint32(target[48:], v)
										target[52] = dims[1]
										target[53] = dims[0]
										target[54] = dims[2]
										target[55] = 5
										target[56] = 5
										target[57] = 7
										target[58] = 0
										for _, f := range facades {
											if effectiveName == f {
												target[58] = 1
											}
										}
										wantCounts[0]++
										wantConsumed -= 4
									}
								}
								ret, consumed := o.invoke(t, op, payload)
								for i := range o.defs {
									if o.defs[i].Data32 != nil {
										t.Fatal("definition allocated image data")
									}
								}
								for i := 0; i < 64; i++ {
									if *(*unsafe.Pointer)(unsafe.Pointer(&o.edges[60*i+32])) != nil {
										t.Fatal("edge definition allocated image data")
									}
								}

								defs, edges := o.metadata()
								if ret != wantRet || consumed != wantConsumed || [2]uint32{*o.count, *o.edgeCount} != wantCounts || !bytes.Equal(defs, wantDefs) || !bytes.Equal(edges, wantEdges) || !bytes.Equal(o.scratch[:64], wantScratch) {
									t.Fatalf("definition edge%v count%d name%q dims%v mode%d format%d ret%d/%d consumed%d/%d", edge, count, name, dims, mode, format, ret, wantRet, consumed, wantConsumed)
								}
								rows = append(rows, record{edge, name, count, mode, format, dims, color, ret, consumed, wantCounts, sha256.Sum256(append(defs, edges...)), append([]byte(nil), o.scratch[:64]...)})
							}
						}
					}
				}
			}
		}
	}
	floorAssetsCapture(t, "definitions", rows, "e72c5717e6ed60608fd938d346df873b456a998dcbabd1a39fa892e34500690a")
}
