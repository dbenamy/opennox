//go:build porttest

package opennox

import (
	"encoding/binary"
	"testing"
)

func TestThingSkipsAudio(t *testing.T) {
	var rows []thingSkipResult
	for _, count := range []int32{-2147483648, -1, 0, 1, 3} {
		for _, n := range []int{0, 1, 31, 255} {
			for _, list := range []int{0, 1, 3} {
				for offset := 0; offset < 8; offset++ {
					b := thingSkipWord(nil, uint32(count))
					for i := int32(0); i < count; i++ {
						b = thingSkipName(b, n, 0x41)
						b = thingSkipFill(b, 9, 0x81)
						for j := 0; j < list; j++ {
							b = thingSkipName(b, max(1, n), 0x51)
						}
						b = append(b, 0)
					}
					rows = append(rows, thingSkipInvoke(t, 0, offset, b, 1, len(b), nil))
				}
			}
		}
	}
	thingSkipCapture(t, "audio", rows, "c9ea00dc6ae76204249978d9fbe9e416a87786a5ac3d6487dcaf222d82360e24")
}
func thingSkipSpellRecord(b []byte, ability bool, n, long, mode int) []byte {
	b = thingSkipName(b, n, 0x41)
	extra := 3
	refs := 2
	if ability {
		extra = 1
		refs = 3
	}
	b = thingSkipFill(b, extra, 0x21)
	if !ability {
		b = thingSkipName(b, n, 0x61)
	}
	for j := 0; j < refs; j++ {
		b = thingSkipRef(b, mode == 1 || (mode == 2 && j%2 == 0), n)
	}
	if !ability {
		b = thingSkipFill(b, 4, 0x71)
	}
	b = thingSkipName(b, n, 0x31)
	b = binary.LittleEndian.AppendUint16(b, uint16(long))
	b = thingSkipFill(b, long, 0x51)
	for j := 0; j < 3; j++ {
		b = thingSkipName(b, n, byte(0x41+j))
	}
	return b
}
func TestThingSkipsSpellsAbilities(t *testing.T) {
	var rows []thingSkipResult
	for _, ability := range []bool{false, true} {
		op := 1
		if ability {
			op = 2
		}
		for _, count := range []int32{-1, 0, 1, 3} {
			for _, n := range []int{0, 1, 31, 255} {
				for _, long := range []int{0, 1, 255, 32767} {
					for mode := 0; mode < 3; mode++ {
						for _, offset := range []int{0, 3, 7} {
							b := thingSkipWord(nil, uint32(count))
							for i := int32(0); i < count; i++ {
								b = thingSkipSpellRecord(b, ability, n, long, mode)
							}
							rows = append(rows, thingSkipInvoke(t, op, offset, b, 1, len(b), nil))
						}
					}
				}
			}
		}
	}
	thingSkipCapture(t, "spells-abilities", rows, "62b79a08d4088d3101357d93d270e575bb1e0b877306b60a30e6675507daec1d")
}
func TestThingSkipsImages(t *testing.T) {
	var rows []thingSkipResult
	for _, count := range []int32{-1, 0, 1, 3} {
		for _, kind := range []byte{0, 1, 2, 3, 255} {
			for _, frames := range []int{0, 1, 3, 255} {
				for _, n := range []int{0, 1, 255} {
					for mode := 0; mode < 3; mode++ {
						b := thingSkipWord(nil, uint32(count))
						for i := int32(0); i < count; i++ {
							b = thingSkipName(b, n, 0x41)
							b = append(b, kind)
							num := 1
							if kind == 2 {
								num = frames
								b = append(b, byte(frames), 0x73)
								b = thingSkipName(b, n, 0x51)
							}
							for j := 0; j < num; j++ {
								b = thingSkipRef(b, mode == 1 || (mode == 2 && j%2 == 0), n)
							}
						}
						rows = append(rows, thingSkipInvoke(t, 3, int(kind%8), b, 1, len(b), nil))
					}
				}
			}
		}
	}
	thingSkipCapture(t, "images", rows, "9221f54959b72482dab6df3101b01a4c889fab7bd5736bf577440c8fb8df30d6")
}
