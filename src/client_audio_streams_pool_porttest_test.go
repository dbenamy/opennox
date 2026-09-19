//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func audioStreamPointer(p unsafe.Pointer) uint32 { return uint32(uintptr(p)) }
func audioStreamMemory(p uint32, n int) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(uintptr(p))), n)
}

func TestClientAudioStreamsPool(t *testing.T) {
	type record struct {
		Count, Size              int
		Popped, Returned, Reused []uint32
	}
	var rows []record
	for _, count := range []int{1, 2, 3, 32, 200} {
		for _, size := range []int{1, 2, 3, 4, 24, 84, 576} {
			t.Run(fmt.Sprintf("count%d-size%d", count, size), func(t *testing.T) {
				pool := legacy.PortTestAudioStreamCall("sub_4BD280", uint32(count), uint32(size))
				if pool == 0 {
					t.Fatal("pool allocation")
				}
				defer legacy.PortTestAudioStreamCall("sub_4BD2D0", pool)
				var pointers []uint32
				row := record{Count: count, Size: size}
				for i := 0; i < count; i++ {
					p := legacy.PortTestAudioStreamCall("sub_4BD2E0", pool)
					want := pool + 8 + uint32(i*(size+4))
					if p != want || !bytes.Equal(audioStreamMemory(p, size), make([]byte, size)) {
						t.Fatal("pool initial allocation order/layout/zeroing")
					}
					pointers = append(pointers, p)
					row.Popped = append(row.Popped, p-pool)
					for j := range audioStreamMemory(p, size) {
						audioStreamMemory(p, size)[j] = byte(i*37 + j*13 + 7)
					}
				}
				if legacy.PortTestAudioStreamCall("sub_4BD2E0", pool) != 0 {
					t.Fatal("exhausted pool allocated another block")
				}
				for cycle := 0; cycle < 3; cycle++ {
					var order []int
					for i := 0; i < count; i++ {
						index := i
						if cycle == 1 {
							index = count - 1 - i
						} else if cycle == 2 {
							index = (i + 3) % count
						}
						order = append(order, index)
						p := pointers[index]
						got := legacy.PortTestAudioStreamCall("sub_4BD300", pool, p)
						if got != p-4 {
							t.Fatal("pool return header")
						}
						row.Returned = append(row.Returned, got-pool)
					}
					for i := count - 1; i >= 0; i-- {
						index := order[i]
						p := legacy.PortTestAudioStreamCall("sub_4BD2E0", pool)
						if p != pointers[index] {
							t.Fatal("returned pool blocks are not LIFO")
						}
						for j, value := range audioStreamMemory(p, size) {
							if value != byte(index*37+j*13+7) {
								t.Fatal("pool return changed payload")
							}
						}
						row.Reused = append(row.Reused, p-pool)
					}
					if legacy.PortTestAudioStreamCall("sub_4BD2E0", pool) != 0 {
						t.Fatal("pool over-allocation after reuse")
					}
				}
				rows = append(rows, row)
			})
		}
	}
	spellbookCapture(t, "client-audio-streams-pool", rows, "e8dbc94935fff466d9b26e66c7386e3cc57a254bd2e31bafd8d14e57267e55e4")
}

func TestClientAudioStreamsReferences(t *testing.T) {
	words, free := alloc.Make([]uint32{}, 24)
	defer free()
	address := audioStreamPointer(unsafe.Pointer(&words[0]))
	type record struct{ Before, Increment, AfterIncrement, Decrement, AfterDecrement, Read uint32 }
	var rows []record
	for _, before := range []uint32{0, 1, 2, 99, 0x7fffffff, 0x80000000, 0xfffffffe, 0xffffffff} {
		for i := range words {
			words[i] = uint32(i)*7919 + 3
		}
		original := append([]uint32(nil), words...)
		words[3] = before
		ret := legacy.PortTestAudioStreamCall("sub_4BD650", address)
		after := words[3]
		if ret != address || after != before+1 {
			t.Fatal("reference increment/width")
		}
		words[3] = before
		dec := legacy.PortTestAudioStreamCall("sub_4BD660", address)
		want := before - 1
		stored := want
		if int32(want) < 0 {
			stored = 0
		}
		if dec != want || words[3] != stored {
			t.Fatal("reference decrement/clamp/return")
		}
		read := legacy.PortTestAudioStreamCall("sub_4BD680", address)
		if read != stored {
			t.Fatal("reference getter")
		}
		for i := range words {
			if i != 3 && words[i] != original[i] {
				t.Fatal("reference operation changed neighbor")
			}
		}
		rows = append(rows, record{before, ret - address, after, dec, words[3], read})
	}
	spellbookCapture(t, "client-audio-streams-references", rows, "01f43c735b5b453cd1f0efd19864e529783ef2f5e4b188795384d4a7ec3ccb47")
}
