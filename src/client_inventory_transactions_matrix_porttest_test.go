//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/client"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"reflect"
	"sort"
	"testing"
	"unsafe"
)

func TestClientInventoryTransactionsStackMatrix(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	var rows []inventoryTransactionResult
	coord, free := alloc.Make([]int32{}, 2)
	defer free()
	mods, freeMods := alloc.Make([]uint32{}, 6)
	defer freeMods()
	lookup, freeLookup := alloc.New(legacy.PortTestInventoryLookup{})
	defer freeLookup()
	for _, name := range []string{"RedPotion", "Bow", "RedApple"} {
		for _, n := range []int{0, 1, 2, 3, 9, 16, 30, 31, 32} {
			for _, index := range []int{0, 19, 20, 21, 42, 63, 82, 83} {
				for _, op := range []int{2, 3, 4, 8, 15, 20, 25} {
					if op == 8 && n == 0 {
						continue
					}
					o.reset(t)
					dr := o.item(t, name, 567)
					cell := &legacy.PortTestInventoryCells()[index]
					cell.Drawable, cell.Count = dr, byte(n)
					for i := 0; i < n; i++ {
						cell.Codes[i] = uint32(100 + i)
					}
					*txword(dr, 292) = 0x12340023
					coord[0], coord[1] = -33, -44
					a, b, c, d := uintptr(0), uintptr(0), uintptr(0), uintptr(0)
					switch op {
					case 2, 3:
						a, b, c = 777, uintptr(dr.TypeIDVal), txptr(unsafe.Pointer(&mods[0]))
						if op == 3 {
							d = txptr(unsafe.Pointer(&coord[0]))
						}
					case 4:
						a, b = 777, uintptr(dr.TypeIDVal)
					case 8:
						lookup.Cell, lookup.Index = cell, uint32(n/2)
						a = txptr(unsafe.Pointer(lookup))
					case 15:
						a, b, c = txptr(dr.C()), uintptr(index/21), uintptr(index%21)
					case 20:
						*o.txwords[8] = uint32(index / 21)
						*o.txwords[9] = uint32(index % 21)
					case 25:
						a, b = uintptr(dr.TypeIDVal), 1
					}
					ret := o.call(op, a, b, c, d)
					r := o.snapshot(t, len(rows), op, ret)
					// Explicit output coordinates are values, never pointer-normalized.
					if op == 3 {
						r.Named = append(r.Named, uint32(coord[0]), uint32(coord[1]))
					}
					rows = append(rows, r)
				}
			}
		}
	}
	inventoryTransactionCapture(t, "stacks", rows, "56549f251b5941bc9fc432b6462d4cb13618b79754de82c9f068a24c4d53b391")
}

func TestClientInventoryTransactionsEquipmentMatrix(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	var rows []inventoryTransactionResult
	classes := []uint32{0, 0x10, 0x1000, 0x1000000, 0x2000000, 0x3001000}
	subclasses := []uint32{0, 1, 2, 4, 8, 16, 32, 64, 128, 256, 0x140, 0x144, 0x1ff, 0xffffffff}
	for _, class := range classes {
		for _, sub := range subclasses {
			o.reset(t)
			dr := o.item(t, "Bow", 123)
			dr.ObjClass, dr.ObjSubClass = object.Class(class), object.SubClass(sub)
			ret := o.call(11, txptr(dr.C()), 0, 0, 0)
			rows = append(rows, o.snapshot(t, len(rows), 11, ret))
		}
	}
	for _, slot := range []int{0, 1, 2, 3, 4, 5, 6, 7, 8} {
		for _, order := range [][]uint32{{0, 4, 64, 256}, {256, 64, 4, 0}, {16, 128, 16, 128}, {64, 64, 256, 256}} {
			o.reset(t)
			for i, sub := range order {
				dr := o.item(t, "RedApple", uint32(100+i))
				dr.ObjClass, dr.ObjSubClass = 0x2000000, object.SubClass(sub)
				ret := o.call(13, txptr(dr.C()), uintptr(slot), 0, 0)
				rows = append(rows, o.snapshot(t, len(rows), 13, ret))
				o.checkEquipment(t)
			}
			for _, code := range []int{102, 100, 103, 101, 999} {
				ret := o.call(9, uintptr(code), 0, 0, 0)
				rows = append(rows, o.snapshot(t, len(rows), 9, ret))
				o.checkEquipment(t)
			}
		}
	}
	inventoryTransactionCapture(t, "equipment", rows, "6d9608cd26b06f3626ffa92b0c4143770f210ae3bb752bb9316177856c7626e6")
}
func (o *inventoryTransactionOwner) checkEquipment(t *testing.T) {
	t.Helper()
	seen := make(map[uint32]bool)
	for _, head := range o.equipment {
		prev := uint32(0)
		for p := head; p != 0; {
			if seen[p] {
				t.Fatal("equipment cycle or duplicate slot membership")
			}
			seen[p] = true
			i, ok := o.identities[(*client.Drawable)(unsafe.Pointer(uintptr(p)))]
			if !ok || !o.objects[i].Live {
				t.Fatal("equipment contains non-live drawable")
			}
			obj := o.objects[i].Drawable
			if *txword(obj, 372) != prev {
				t.Fatal("equipment previous link")
			}
			prev = p
			p = *txword(obj, 368)
		}
	}
}

func TestClientInventoryTransactionsCompactionConservation(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	var rows []inventoryTransactionResult
	for seed := 0; seed < 64; seed++ {
		o.reset(t)
		code := uint32(100)
		var want []uint32
		for row := 0; row < 20; row++ {
			for col := 0; col < 4; col++ {
				if (row*13+col*7+seed)%5 > 1 {
					continue
				}
				name := []string{"RedApple", "Meat", "Bow"}[(row+col+seed)%3]
				dr := o.item(t, name, 0)
				cell := &legacy.PortTestInventoryCells()[col*21+row]
				cell.Drawable, cell.Count = dr, byte(1+(row*3+col+seed)%32)
				cell.Equipped = uint32((row + col) % 2)
				for j := 0; j < int(cell.Count); j++ {
					cell.Codes[j] = code
					want = append(want, code)
					code++
				}
			}
		}
		ret := o.call(7, 0, 0, 0, 0)
		// C exposes a derived grid pointer as a return even past the array.
		delta := ret - o.cell(0)
		if delta > 16000 {
			t.Fatalf("compaction return not derived from grid: %x", ret)
		}
		var got []uint32
		for _, cell := range legacy.PortTestInventoryCells() {
			if cell.Count > 32 {
				t.Fatal("stack overflow")
			}
			got = append(got, cell.Codes[:cell.Count]...)
		}
		sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("compaction seed%d lost or duplicated codes", seed)
		}
		rows = append(rows, o.snapshot(t, len(rows), 7, 0xea800000+delta))
	}
	inventoryTransactionCapture(t, "compaction", rows, "90e01e41e7810dd0a0f2fe2e462a6d191d012d4dd1a8eec4ac0706955ad8da7c")
}

func TestClientInventoryTransactionsCapacityMatrix(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	var rows []inventoryTransactionResult
	for _, flags := range []uint32{1, 2049, 4097, 6145} {
		for _, class := range []uint32{0, 0x10, 0x4000000, 0x4000010} {
			for _, count := range []byte{1, 2, 3, 8, 9, 10, 30, 31, 32} {
				for _, quantity := range []int{0, 1, 2, 3, 9, 31, 32, 100, -1} {
					o.reset(t)
					noxflags.ResetGame()
					noxflags.SetGame(noxflags.GameFlag(flags))
					dr := o.item(t, "RedApple", 123)
					dr.ObjClass = object.Class(class)
					for i := range legacy.PortTestInventoryCells() {
						cell := &legacy.PortTestInventoryCells()[i]
						cell.Drawable, cell.Count = dr, count
					}
					ret := o.call(25, uintptr(dr.TypeIDVal), uintptr(quantity), 0, 0)
					limit := 31
					if class&0x10 != 0 {
						limit = 3
						if flags&6144 != 0 {
							limit = 9
						}
					}
					want := uint32(0)
					if class&0x4000000 == 0 && quantity+int(count) <= limit {
						want = 80
					}
					if ret != want {
						t.Fatalf("flags%x class%x count%d quantity%d capacity%d want%d", flags, class, count, quantity, ret, want)
					}
					rows = append(rows, o.snapshot(t, len(rows), 25, ret))
				}
			}
		}
	}
	inventoryTransactionCapture(t, "capacity", rows, "66e46a8c080b3dd52028c046b6f7118eb809e0b04565ad3545996fefa26d304d")
}
