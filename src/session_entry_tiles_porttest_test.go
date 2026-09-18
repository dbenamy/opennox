//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestSessionEntryTileClear(t *testing.T) {
	grid, _, restore := legacy.PortTestTileCompositionGrid()
	defer restore()
	head := legacy.PortTestSessionEntryTileFreeHead()
	old := *head
	*head = 0
	defer func() { *head = old }()
	nodes, free := alloc.Make([][5]uint32{}, 4)
	defer free()
	ids := map[uint32]int{0: 0}
	for i := range nodes {
		nodes[i] = [5]uint32{uint32(10 + i), 20, 30, 40, 0}
		ids[uint32(uintptr(unsafe.Pointer(&nodes[i])))] = i + 1
	}
	ptr := func(i int) uint32 { return uint32(uintptr(unsafe.Pointer(&nodes[i]))) }
	nodes[0][4] = ptr(1)
	for x := range grid {
		for y := range grid[x] {
			grid[x][y] = [11]uint32{0xaabbccdd, uint32(x), 11, 12, 13, 0, uint32(y), 21, 22, 23, 0}
		}
	}
	grid[0][0][5] = ptr(0)
	grid[63][64][10] = ptr(2)
	grid[127][127][5] = ptr(3)
	type row struct {
		Pass  int
		Cells [][11]uint32
		Free  []int
	}
	var rows []row
	for pass := 0; pass < 2; pass++ {
		legacy.PortTestSessionEntryScalar("tile-clear", 0)
		r := row{Pass: pass}
		for x := range grid {
			for y := range grid[x] {
				got := grid[x][y]
				want := [11]uint32{0xaabbcc00, 255, 11, 12, 13, 0, 255, 21, 22, 23, 0}
				if got != want {
					t.Fatalf("tile%d,%d got%x want%x", x, y, got, want)
				}
				if (x == 0 && y == 0) || (x == 63 && y == 64) || (x == 127 && y == 127) {
					r.Cells = append(r.Cells, got)
				}
			}
		}
		seen := map[uint32]bool{}
		for p := *head; p != 0; p = (*[5]uint32)(unsafe.Pointer(uintptr(p)))[4] {
			id, ok := ids[p]
			if !ok || seen[p] {
				t.Fatal("invalid returned subtile list")
			}
			seen[p] = true
			r.Free = append(r.Free, id)
		}
		if len(seen) != 4 {
			t.Fatalf("returned subtiles%d want4", len(seen))
		}
		rows = append(rows, r)
	}
	spellbookCapture(t, "session-entry-tile-clear", rows, "8140baaa250cbdb2a279e59c93719717d729ede0bd929c5a28b61b16dac7ffc6")
}
