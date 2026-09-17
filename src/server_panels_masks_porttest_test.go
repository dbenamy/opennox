//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestServerPanelsMaskSetters(t *testing.T) {
	type row struct {
		Before, Mask uint32
		On           int
		After        uint32
		Pointer      bool
	}
	var byteRows, wordRows []row
	b, free := alloc.Make([]byte{}, 3)
	defer free()
	for before := 0; before < 256; before++ {
		for mask := 0; mask < 256; mask++ {
			for _, on := range []int{0, 1, -1} {
				b[0], b[1], b[2] = 0x5a, byte(before), 0xa5
				ptr := legacy.PortTestServerPanelsByte(&b[1], byte(mask), on)
				want := byte(before) &^ byte(mask)
				if on != 0 {
					want = byte(before) | byte(mask)
				}
				if b[1] != want || b[0] != 0x5a || b[2] != 0xa5 || !ptr {
					t.Fatalf("byte %x/%x/%d", before, mask, on)
				}
				// Exhaust independently without writing a 196k-record capture.
				if before%17 == 0 && mask%17 == 0 {
					byteRows = append(byteRows, row{uint32(before), uint32(mask), on, uint32(b[1]), ptr})
				}
			}
		}
	}
	words, freeWords := alloc.Make([]uint32{}, 3)
	defer freeWords()
	for _, before := range []uint32{0, 1, 0x55555555, 0xaaaaaaaa, 0x7fffffff, 0x80000000, 0xffffffff} {
		for _, mask := range []uint32{0, 1, 0x80, 0xff, 0x100, 0x8000, 0x10000, 0x1000000, 0x80000000, 0xffffffff} {
			for _, on := range []int{0, 1, -1} {
				words[0], words[1], words[2] = 0x12345678, before, 0x87654321
				ptr := legacy.PortTestServerPanelsWord(&words[1], mask, on)
				want := before &^ mask
				if on != 0 {
					want = before | mask
				}
				if words[1] != want || words[0] != 0x12345678 || words[2] != 0x87654321 || !ptr {
					t.Fatal("word mask or ownership")
				}
				wordRows = append(wordRows, row{before, mask, on, words[1], ptr})
			}
		}
	}
	spellbookCapture(t, "server-panels-byte-masks", byteRows, "cbd87dd949459e3da456a0863724622dcb9677eec7565491fc55fd07c534537b")
	spellbookCapture(t, "server-panels-word-masks", wordRows, "771b47269b73c7066b3337d8cf45a96f962bd899e4da7bdb9402ebf93403b5b5")
}
func TestServerPanelsSpellBits(t *testing.T) {
	type row struct {
		Index, On                   int
		Before, Mask, Result, After uint32
		Set                         bool
	}
	var rows []row
	words, free := alloc.Make([]uint32{}, 258)
	defer free()
	indices := []int{-2147483648, -8193, -8192, -8191, -33, -32, -31, -1, 0, 1, 31, 32, 33, 63, 64, 127, 128, 136, 159, 160, 255, 256, 8191, 8192, 8193, 2147483647}
	for _, index := range indices {
		for _, before := range []uint32{0, 0x55555555, 0xaaaaaaaa, 0xffffffff} {
			for _, on := range []int{0, 1, -1} {
				t.Run(fmt.Sprintf("%d-%x-%d", index, before, on), func(t *testing.T) {
					for i := range words {
						words[i] = before
					}
					words[0] = 0x12345678
					words[257] = 0x87654321
					group := int(byte(index / 32))
					mask := uint32(1) << uint(uint32(index)&31)
					result, set := legacy.PortTestServerPanelsBit(&words[1], index, on)
					want := before &^ mask
					wantResult := ^mask
					if on != 0 {
						want = before | mask
						wantResult = mask
					}
					if result != wantResult || set != (on != 0) || words[group+1] != want || words[0] != 0x12345678 || words[257] != 0x87654321 {
						t.Fatal("bit value/index/result/guards")
					}
					for i := 1; i < 257; i++ {
						if i != group+1 && words[i] != before {
							t.Fatal("unrelated word changed")
						}
					}
					rows = append(rows, row{index, on, before, mask, result, words[group+1], set})
				})
			}
		}
	}
	spellbookCapture(t, "server-panels-spell-bits", rows, "227e65da51585f1538ef1a3b6e8741f2286c6dfb3530c773e7e4e64689e54364")
}
func TestServerPanelsMaskQueries(t *testing.T) {
	b := unsafe.Slice(memmap.PtrUint8(0x5D4594, 1045448), 16)
	old := append([]byte(nil), b...)
	defer copy(b, old)
	type row struct {
		Bits   uint32
		Before byte
		Armor  bool
		Index  int
		Result bool
	}
	var rows []row
	for _, bits := range []uint32{0, 1, 0x55555555, 0xaaaaaaaa, 0x7fffffff, 0x80000000, 0xffffffff} {
		for _, preceding := range []byte{0, 0x5a, 255} {
			for _, armor := range []bool{false, true} {
				b[3] = preceding
				binary.LittleEndian.PutUint32(b[4:], bits)
				binary.LittleEndian.PutUint32(b[8:], bits)
				for index := 0; index < 256; index++ {
					mask := uint32(1) << uint(index&31)
					want := bits&mask != 0
					if !armor && index&31 == 31 {
						want = false
					}
					got := legacy.PortTestServerPanelsClass(byte(index), armor)
					if got != want {
						t.Fatalf("class %d bits %x armor %t got %t", index, bits, armor, got)
					}
					rows = append(rows, row{bits, preceding, armor, index, got})
				}
				for _, mask := range []uint32{0, 1, 0xff, 0x100, 0x01000101, 0x7fffffff, 0x80000000, 0x80000001, 0xffffffff} {
					want := bits&mask != 0
					if !armor {
						if int32(mask) <= 0 {
							want = preceding&byte(mask) != 0
						} else {
							shift := uint(0)
							for mask>>shift > 255 {
								shift += 8
							}
							want = byte(bits>>shift)&byte(mask>>shift) != 0
						}
					}
					if got := legacy.PortTestServerPanelsMaskQuery(mask, armor); got != want {
						t.Fatalf("query %x bits %x armor %t got %t", mask, bits, armor, got)
					}
				}
			}
		}
	}
	spellbookCapture(t, "server-panels-mask-queries", rows, "60c23547384230ba5b326e969354b31c1d6d087ecd2d58f7670c1b1892dc72ab")
}
