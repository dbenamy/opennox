//go:build porttest

package opennox

import (
	"bytes"
	"testing"
)

func thingSkipAligned(b []byte, offset int, v byte) []byte {
	for (offset+len(b))%8 != 0 {
		b = append(b, 0xb1)
	}
	b = append(b, v)
	return thingSkipFill(b, 7, 0xc1)
}
func thingSkipWallRecord(offset, materials, n, pattern, mode int, end uint32) (b, wantScratch []byte, ret, consumed int) {
	wantScratch = bytes.Repeat([]byte{0xa5}, 256)
	name := func(seed byte) {
		b = thingSkipName(b, n, seed)
		for i := 0; i < n; i++ {
			wantScratch[i] = seed + byte(i%17)
		}
		wantScratch[n] = 0
	}
	b = thingSkipWord(b, 0x12345678)
	name(0x41)
	b = thingSkipFill(b, 14, 0x61)
	b = thingSkipAligned(b, offset, byte(materials))
	for i := 0; i < materials; i++ {
		name(byte(0x51 + i))
		if i == 7 && materials > 8 {
			consumed = len(b)
			ret = 0
			return
		}
	}
	for j := 0; j < 3; j++ {
		name(byte(0x71 + j))
	}
	b = append(b, 0xd7)
	for side := 0; side < 15; side++ {
		count := 0
		switch pattern {
		case 1:
			count = 1
		case 2:
			count = side % 4
		case 3:
			if side == 14 {
				count = 255
			}
		}
		b = thingSkipAligned(b, offset, byte(count))
		for i := 0; i < count; i++ {
			for face := 0; face < 4; face++ {
				b = thingSkipFill(b, 8, 0x91)
				b = thingSkipRef(b, mode == 1 || (mode == 2 && (i+face)%2 == 0), n)
			}
		}
	}
	b = thingSkipWord(b, end)
	if end == 0x454e4420 {
		ret = 1
	}
	consumed = len(b)
	return
}
func TestThingSkipsWalls(t *testing.T) {
	var rows []thingSkipResult
	for offset := 0; offset < 8; offset++ {
		for _, materials := range []int{0, 1, 8, 9} {
			for _, n := range []int{0, 1, 31, 255} {
				for pattern := 0; pattern < 4; pattern++ {
					for mode := 0; mode < 3; mode++ {
						for _, end := range []uint32{0x454e4420, 0x12345678} {
							b, scratch, ret, consumed := thingSkipWallRecord(offset, materials, n, pattern, mode, end)
							// Retain trailing input after the rejected ninth-material count. The original
							// stops after reading the first eight names, without touching this suffix.
							if materials > 8 {
								b = thingSkipName(b, n, 0x79)
								b = append(b, 0x31, 0x42, 0x53)
							}
							rows = append(rows, thingSkipInvoke(t, 6, offset, b, ret, consumed, scratch))
						}
					}
				}
			}
		}
	}
	thingSkipCapture(t, "walls", rows, "657885a2ecd7241a10a7564ea541b30ed3b8034b1286ffe080cfc0427553faac")
}
