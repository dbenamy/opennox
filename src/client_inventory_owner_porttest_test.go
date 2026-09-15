//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
)

type uiInventoryOwner struct {
	*meterOwner
	words       []*uint32
	regions     [][]byte
	items       []*client.Drawable
	shortReturn uint32
}

func newUIInventoryOwner(t *testing.T) *uiInventoryOwner {
	o := &uiInventoryOwner{meterOwner: newMeterOwner(t)}
	var restore func()
	o.words, restore = legacy.PortTestUIInventoryWords()
	t.Cleanup(restore)
	for _, r := range [][2]uintptr{{1049788, 8}, {1049848, 4}, {1062536, 52}, {1063100, 1024}} {
		b := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, r[0])), r[1])
		old := append([]byte(nil), b...)
		o.regions = append(o.regions, b)
		t.Cleanup(func() { copy(b, old) })
	}
	o.items = o.meterOwner.items(t)
	bow, quiver := o.weapons(t)
	o.items = append(o.items, bow, quiver)
	return o
}
func (o *uiInventoryOwner) reset(t *testing.T) {
	o.plain(t)
	for _, p := range o.words {
		*p = 0
	}
	for _, b := range o.regions {
		clear(b)
	}
	for _, dr := range o.items {
		dr.NetCode32 = 0
		for _, off := range []uintptr{292, 294, 448, 450} {
			*(*uint16)(unsafe.Add(dr.C(), off)) = 0
		}
		for _, off := range []uintptr{368, 372} {
			*(*uint32)(unsafe.Add(dr.C(), off)) = 0
		}
	}
	*(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 4)) = 0
	*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 3684)) = 0
	o.shortReturn = 0
	*o.words[4] = uint32(uintptr(o.parent.C()))
	*o.words[5] = uint32(uintptr(o.parent.C()))
}
func (o *uiInventoryOwner) cell(index int) uint32 {
	return uint32(uintptr(unsafe.Pointer(&o.grid[index*148])))
}
func (o *uiInventoryOwner) norm(v uint32) uint32 {
	start := o.cell(0)
	if v >= start && v < start+uint32(len(o.grid)) {
		return 0xee000000 + v - start
	}
	for _, off := range []uintptr{1049788, 1049792} {
		if v == uint32(uintptr(memmap.PtrOff(0x5D4594, off))) {
			return 0xef000000 + uint32(off)
		}
	}
	return o.normalize(v)
}

type uiInventoryResult struct {
	Case, Op                  int
	Args                      [3]uint32
	Return                    uint32
	Named, Regions, Equipment []uint32
	Grid                      string
	Items                     [][7]uint32
	Messages                  [][]byte
	Hidden                    bool
	PlayerByte                byte
	Meters                    [2][2]uint32
}

func (o *uiInventoryOwner) call(t *testing.T, id, op int, a, b, c uint32) uiInventoryResult {
	t.Helper()
	ret := legacy.PortTestUIInventoryCall(op, a, b, c)
	if op == 18 && o.shortReturn != 0 {
		if ret != uint32(int32(int16(o.shortReturn))) {
			t.Fatalf("durability pointer return %#x expected low16 of %#x", ret, o.shortReturn)
		}
		ret = 0xedffffff
	} else {
		ret = o.norm(ret)
	}
	r := uiInventoryResult{Case: id, Op: op, Args: [3]uint32{a, b, c}, Return: ret, Hidden: o.parent.GetFlags().Has(0x10)}
	for _, p := range o.words {
		r.Named = append(r.Named, o.norm(*p))
	}
	for _, reg := range o.regions {
		for i := 0; i < len(reg); i += 4 {
			r.Regions = append(r.Regions, o.norm(binary.LittleEndian.Uint32(reg[i:])))
		}
	}
	grid := append([]byte(nil), o.grid...)
	for i := 0; i < 84; i++ {
		binary.LittleEndian.PutUint32(grid[i*148:], o.norm(binary.LittleEndian.Uint32(grid[i*148:])))
	}
	r.Grid = fmt.Sprintf("%x", sha256.Sum256(grid))
	for _, v := range o.equipment {
		r.Equipment = append(r.Equipment, o.norm(v))
	}
	for _, dr := range o.items {
		r.Items = append(r.Items, [7]uint32{dr.NetCode32, uint32(*(*uint16)(unsafe.Add(dr.C(), 292))), uint32(*(*uint16)(unsafe.Add(dr.C(), 294))), uint32(*(*uint16)(unsafe.Add(dr.C(), 448))), uint32(*(*uint16)(unsafe.Add(dr.C(), 450))), o.norm(*(*uint32)(unsafe.Add(dr.C(), 368))), o.norm(*(*uint32)(unsafe.Add(dr.C(), 372)))})
	}
	o.c.srv.NetList.ByInd(31, netlist.Kind0).Each(func(b []byte) bool { r.Messages = append(r.Messages, append([]byte(nil), b...)); return false })
	r.PlayerByte = *(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 3684))
	for i := 0; i < 2; i++ {
		m := o.meters.Records[5+i]
		r.Meters[i] = [2]uint32{m.Current, m.Maximum}
	}
	return r
}
func uiInventoryCapture(t *testing.T, label string, rows []uiInventoryResult, want string) {
	t.Helper()
	data, err := json.Marshal(rows)
	if err != nil {
		t.Fatal(err)
	}
	if p := os.Getenv("OPENNOX_INVENTORY_CAPTURE"); p != "" {
		if err := os.WriteFile(p+"-"+label+".json", data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	got := fmt.Sprintf("%x", sha256.Sum256(data))
	t.Logf("%s: %d results %s", label, len(rows), got)
	if got != want {
		t.Fatalf("%s hash %s want frozen C %s", label, got, want)
	}
}
