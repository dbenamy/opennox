//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"testing"
	"unsafe"
)

func TestFloorAssetsImageBinding(t *testing.T) {
	o := newFloorAssetsOwner(t)
	type record struct {
		Edge, Found, Allocated bool
		Count, Mode, Format    int
		Dims                   [3]byte
		End                    uint32
		Return, Consumed       int
		Images                 []uint32
		Scratch                [32]byte
	}
	var rows []record
	for _, edge := range []bool{false, true} {
		cap := 176
		dimsList := [][3]byte{{0, 0, 0}, {1, 1, 1}, {2, 3, 2}, {255, 1, 1}, {255, 255, 0}}
		if edge {
			cap = 64
			dimsList = [][3]byte{{0, 0, 0}, {1, 1, 1}, {2, 3, 2}, {255, 1, 0}, {0, 255, 255}}
		}
		for _, count := range []int{1, 3, cap} {
			for _, found := range []bool{false, true} {
				for _, dims := range dimsList {
					for mode := 0; mode < 3; mode++ {
						formats := []int{0}
						ends := []uint32{0x454e4420}
						if edge {
							formats = []int{0, 1, 255}
							ends = append(ends, 0x12345678)
						}
						for _, format := range formats {
							for _, end := range ends {
								o.reset()
								index := count - 1
								for i := 0; i < count; i++ {
									name := fmt.Sprintf("Definition-%d", i)
									if edge {
										copy(o.edges[60*i:][:32], name)
									} else {
										copy(o.defs[i].NameBuf[:], name)
									}
								}
								name := fmt.Sprintf("Definition-%d", index)
								if !found {
									name = "Missing"
								}
								op := 4
								frames := int(dims[0]) * int(dims[1]) * int(dims[2])
								payload := floorAssetsFloor(name, [3]byte{255, 0, 255}, dims, mode)
								*o.count = uint32(count)
								if edge {
									*o.count = 0
									*o.edgeCount = uint32(count)
									op = 5
									frames = 2 * int(dims[0]) * (int(dims[1]) + int(dims[2]))
									payload = floorAssetsEdge(name, dims, mode, byte(format), end)
								}
								beforeDefs, beforeEdges := o.metadata()
								wantRet := 1
								wantConsumed := len(payload) - 4
								wantScratch := bytes.Repeat([]byte{0xa5}, 64)
								allocates := found && (!edge || format != 1)
								if edge {
									wantConsumed = len(payload)
									if end != 0x454e4420 {
										wantRet = 0
									}
								}
								if !found {
									wantRet = 0
									wantConsumed = 5 + len(name)
								} else if edge && format == 1 {
									wantRet = 0
									wantConsumed = 17 + len(name)
								} else {
									wantScratch = floorAssetsScratch(mode, frames)
								}
								ret, consumed := o.invoke(t, op, payload)
								afterDefs, afterEdges := o.metadata()
								expectIndex := -1
								if allocates {
									expectIndex = index
								}
								o.assertOnlyData(t, edge, expectIndex)
								wantCounts := [2]uint32{uint32(count), 0}
								if edge {
									wantCounts = [2]uint32{0, uint32(count)}
								}
								if [2]uint32{*o.count, *o.edgeCount} != wantCounts {
									t.Fatal("binding changed definition count")
								}

								if ret != wantRet || consumed != wantConsumed || !bytes.Equal(beforeDefs, afterDefs) || !bytes.Equal(beforeEdges, afterEdges) || !bytes.Equal(o.scratch[:64], wantScratch) {
									t.Fatalf("binding edge%v found%v count%d dims%v mode%d format%d end%x ret%d/%d consumed%d/%d", edge, found, count, dims, mode, format, end, ret, wantRet, consumed, wantConsumed)
								}
								ptr := o.defs[index].Data32
								if edge {
									ptr = *(*unsafe.Pointer)(unsafe.Pointer(&o.edges[60*index+32]))
								}
								if !allocates && ptr != nil {
									t.Fatal("rejected binding allocated data")
								}
								if allocates && frames > 0 && ptr == nil {
									t.Fatal("binding did not allocate image array")
								}
								var refs []uint32
								if allocates {
									for i, p := range unsafe.Slice((*uint32)(ptr), frames) {
										ref, ok := o.refs[p]
										if !ok {
											t.Fatal("unowned image handle")
										}
										want := uint32(i%8 + 1)
										if mode == 1 || (mode == 2 && i%2 != 0) {
											want = 0
										}
										if ref != want {
											t.Fatalf("image%d got%d want%d", i, ref, want)
										}
										refs = append(refs, ref)
									}
									if edge {
										tail := unsafe.Slice((*byte)(unsafe.Add(ptr, 4*frames)), frames)
										for _, b := range tail {
											if b != 0 {
												t.Fatal("edge allocation tail not zero")
											}
										}
									}
								}
								rows = append(rows, record{edge, found, ptr != nil, count, mode, format, dims, end, ret, consumed, refs, sha256.Sum256(o.scratch)})
							}
						}
					}
				}
			}
		}
	}
	floorAssetsCapture(t, "binding", rows, "b1e4d0288b1fb06bcf9bb28aeb20ba9ee02e4f34b48c277b160c3874c11a76be")
}

func TestFloorAssetsFirstNameMatch(t *testing.T) {
	o := newFloorAssetsOwner(t)
	type record struct {
		Edge             bool
		Return, Consumed int
		Images           []uint32
	}
	var rows []record
	for _, edge := range []bool{false, true} {
		o.reset()
		for i, name := range []string{"Duplicate", "Other", "Duplicate"} {
			if edge {
				copy(o.edges[60*i:][:32], name)
			} else {
				copy(o.defs[i].NameBuf[:], name)
			}
		}
		*o.count = 3
		op, n := 4, 1
		input := floorAssetsFloor("Duplicate", [3]byte{0, 0, 0}, [3]byte{1, 1, 1}, 0)
		wantConsumed := len(input) - 4
		if edge {
			*o.count = 0
			*o.edgeCount = 3
			op, n = 5, 2
			input = floorAssetsEdge("Duplicate", [3]byte{1, 1, 0}, 0, 0, 0x454e4420)
			wantConsumed = len(input)
		}
		ret, consumed := o.invoke(t, op, input)
		if ret != 1 || consumed != wantConsumed {
			t.Fatal("duplicate-name binding failed")
		}
		o.assertOnlyData(t, edge, 0)
		ptr := o.defs[0].Data32
		if edge {
			ptr = *(*unsafe.Pointer)(unsafe.Pointer(&o.edges[32]))
		}
		if ptr == nil {
			t.Fatal("first matching definition was not bound")
		}
		var refs []uint32
		for i, p := range unsafe.Slice((*uint32)(ptr), n) {
			ref, ok := o.refs[p]
			if !ok || ref != uint32(i+1) {
				t.Fatal("duplicate binding image mismatch")
			}
			refs = append(refs, ref)
		}
		rows = append(rows, record{edge, ret, consumed, refs})
	}
	floorAssetsCapture(t, "first-match", rows, "714b6b32de6da860eb869b5088f1f24f0b10875c6bc8d54e30286c813f617a72")
}
