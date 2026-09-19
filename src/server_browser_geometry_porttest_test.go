//go:build porttest

package opennox

import (
	"encoding/binary"
	"math"
	"testing"
	"unsafe"

	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestServerBrowserModeNames(t *testing.T) {
	o := newEntryOwner(t)
	names := []string{"Quest", "CTF", "Highlander", "KotR", "Flagball", "Chat", "Arena"}
	var entries []strman.Entry
	for _, n := range names {
		entries = append(entries, strman.Entry{ID: strman.ID("noxworld.c:" + n), Vals: []strman.Variant{{Str: "mode:" + n}}})
	}
	set, restore := o.c.srv.Server.PortTestMeterStrings(entries...)
	defer restore()
	set(0)
	type row struct {
		Mode uint16
		Name string
	}
	var rows []row
	for mode := uint32(0); mode < 65536; mode++ {
		name := legacy.PortTestServerBrowserMode(uint16(mode))
		want := "Arena"
		// Highest-priority active mode wins, independent of unrelated status bits.
		for _, v := range []struct {
			mask uint32
			name string
		}{{0x1000, "Quest"}, {0x20, "CTF"}, {0x400, "Highlander"}, {0x10, "KotR"}, {0x40, "Flagball"}, {0x80, "Chat"}} {
			if mode&v.mask != 0 {
				want = v.name
				break
			}
		}
		if name != "mode:"+want {
			t.Fatalf("mode %#x: %q want %q", mode, name, want)
		}
		rows = append(rows, row{uint16(mode), name})
	}
	spellbookCapture(t, "server-browser-mode-names", rows, "8d9eaaff7082ea232196f4ebfe730c643bf75d3070c6d875e0a37ca486f14295")
}

func TestServerBrowserMapHit(t *testing.T) {
	radiusBytes := serverConfigOwnBytes(t, 0x581450, 9720, 8)
	record, free := alloc.Calloc(1, 172)
	defer free()
	coords, freeCoords := alloc.Make([]uint32{}, 2)
	defer freeCoords()
	type row struct {
		Radius string
		Server [2]int16
		Point  [2]uint32
		Hit    int
	}
	var rows []row
	centers := [][2]int16{{0, 0}, {100, 80}, {-1, -1}, {-32768, 32767}}
	for _, rad := range []float64{-1, 0, 1, 5, 10, math.Inf(1), math.NaN()} {
		binary.LittleEndian.PutUint64(radiusBytes, math.Float64bits(rad))
		for _, center := range centers {
			*(*int16)(unsafe.Add(record, 44)) = center[0]
			*(*int16)(unsafe.Add(record, 46)) = center[1]
			for _, dx := range []int32{-11, -10, -9, -4, -3, -1, 0, 1, 3, 4, 9, 10, 11} {
				for _, dy := range []int32{-11, -10, -9, -4, -3, -1, 0, 1, 3, 4, 9, 10, 11} {
					coords[0] = uint32(int32(center[0]) + dx)
					coords[1] = uint32(int32(center[1]) + dy)
					got := legacy.PortTestServerBrowserHit(unsafe.Pointer(&coords[0]), record)
					// The qualified C expression subtracts unsigned mouse coordinates before
					// converting to double. Preserve that observable wrap, including quadrants.
					x := float64(uint32(int32(center[0])) - coords[0])
					y := float64(uint32(int32(center[1])) - coords[1])
					want := 0
					if math.Sqrt(x*x+y*y) <= rad {
						want = 1
					}
					if got != want {
						t.Fatalf("radius %v server %v point %v: %d want %d", rad, center, coords, got, want)
					}
					rows = append(rows, row{Radius: serverBrowserFloatText(rad), Server: center, Point: [2]uint32{coords[0], coords[1]}, Hit: got})
				}
			}
		}
	}
	spellbookCapture(t, "server-browser-map-hit", rows, "1edf2217bb2f7a85446a91b97e4b85af2ea43e10861ef5a5ab1055b402adf04d")
}
func serverBrowserFloatText(v float64) string {
	switch {
	case math.IsNaN(v):
		return "nan"
	case math.IsInf(v, 1):
		return "inf"
	}
	return map[float64]string{-1: "-1", 0: "0", 1: "1", 5: "5", 10: "10"}[v]
}

func TestServerBrowserPopupClamp(t *testing.T) {
	out, free := alloc.Make([]uint32{}, 4)
	defer free()
	type row struct {
		X, Y     int32
		Out      [4]uint32
		Identity bool
	}
	var rows []row
	for _, x := range []int32{-2147483648, -101, -1, 0, 215, 216, 315, 316, 499, 500, 501, 600, 2147483547, 2147483647} {
		for _, y := range []int32{-2147483648, -21, -1, 0, 26, 27, 46, 47, 250, 251, 270, 271, 272, 451, 2147483647} {
			copy(out, []uint32{0xaabbccdd, 0x11223344, 0xabcdef01, 0x98765432})
			identity := legacy.PortTestServerBrowserClamp(x, y, unsafe.Pointer(&out[0])) == uintptr(unsafe.Pointer(&out[0]))
			wantX, wantY := uint32(x-100), uint32(y-20)
			if x+100 > 600 {
				wantX = 400
			}
			if y+180 > 451 {
				wantY = 251
			}
			if wantY < 27 {
				wantY = 27
			}
			if int32(wantX) < 216 {
				wantX = 216
			}
			if !identity || out[0] != wantX || out[1] != wantY || out[2] != 0xabcdef01 || out[3] != 0x98765432 {
				t.Fatalf("clamp %d,%d: %v want %d,%d identity=%v", x, y, out, wantX, wantY, identity)
			}
			rows = append(rows, row{x, y, [4]uint32{out[0], out[1], out[2], out[3]}, identity})
		}
	}
	spellbookCapture(t, "server-browser-popup-clamp", rows, "1c2ff85279dd3a9c8a979e8ad2d5e82a9acdd2df7ab50eb217ec6b280e178260")
}

func TestServerBrowserMapListCount(t *testing.T) {
	radius := serverConfigOwnBytes(t, 0x581450, 9720, 8)
	binary.LittleEndian.PutUint64(radius, math.Float64bits(10))
	coords, free := alloc.Make([]uint32{}, 2)
	defer free()
	var nodes []unsafe.Pointer
	for _, v := range [][2]int16{{0, 0}, {3, 4}, {6, 8}, {7, 8}, {-3, -4}, {10, 0}} {
		p, done := alloc.Calloc(1, 172)
		defer done()
		*(*int16)(unsafe.Add(p, 44)) = v[0]
		*(*int16)(unsafe.Add(p, 46)) = v[1]
		nodes = append(nodes, p)
	}
	type row struct {
		Count int
		Point [2]uint32
		Hits  int
	}
	var rows []row
	for n := 0; n <= len(nodes); n++ {
		head, restore := legacy.PortTestServerBrowserList(nodes[:n])
		for _, pt := range [][2]uint32{{0, 0}, {1, 1}, {3, 4}, {10, 10}, {0xffffffff, 0xffffffff}} {
			coords[0], coords[1] = pt[0], pt[1]
			got := legacy.PortTestServerBrowserCount(unsafe.Pointer(&coords[0]), head)
			want := 0
			for _, p := range nodes[:n] {
				want += legacy.PortTestServerBrowserHit(unsafe.Pointer(&coords[0]), p)
			}
			if got != want {
				t.Fatalf("list count %d/%v: %d want %d", n, pt, got, want)
			}
			rows = append(rows, row{n, pt, got})
		}
		restore()
	}
	if got := legacy.PortTestServerBrowserCount(unsafe.Pointer(&coords[0]), nil); got != 0 {
		t.Fatal("nil list", got)
	}
	spellbookCapture(t, "server-browser-map-list", rows, "9b1dc57391bc90f8d356f84d8f397a9eae3f0230f90d000a535fc071252fed68")
}
