//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func TestClientInventoryScalars(t *testing.T) {
	o := newUIInventoryOwner(t)
	var rows []uiInventoryResult
	for _, v := range []uint32{0, 1, 2, 5, 6, 7, 127, 128, 255, 256, 65535, 0x7fffffff, 0x80000000, 0xffffffff} {
		o.reset(t)
		for _, op := range []int{4, 5, 6, 7, 8, 9, 11, 12, 13, 16, 17, 19, 20, 23, 27, 28} {
			*o.words[1] = v
			*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 3684)) = byte(v)
			rows = append(rows, o.call(t, len(rows), op, v, ^v, 0))
		}
	}
	for _, index := range []uint32{0, 1, 127, 128, 255, 256, 257, 0xffffffff} {
		for _, bits := range []uint32{0, 0x80000000, 1, 0x007fffff, 0x00800000, 0x3f800000, 0x7f7fffff, 0x7f800000, 0xff800000, 0x7fc12345} {
			o.reset(t)
			rows = append(rows, o.call(t, len(rows), 10, index, bits, 0))
		}
	}
	o.reset(t)
	*o.meters.NamedWord("dword_8531A0_2576") = 0
	rows = append(rows, o.call(t, len(rows), 16, 0, 0, 0))
	uiInventoryCapture(t, "inventory-scalars", rows, "fb9d8cdf78d3d73ca4aaf6655518fce96d9505438c6265c0ff39c26a93b2a055")
}
func TestClientInventoryScalarContracts(t *testing.T) {
	o := newUIInventoryOwner(t)
	o.reset(t)
	if legacy.PortTestUIInventoryCall(6, 255, 0, 0) != 0xffffffff || legacy.PortTestUIInventoryCall(7, 0, 0, 0) != 255 {
		t.Fatal("signed setter/unsigned getter")
	}
	if legacy.PortTestUIInventoryCall(10, 257, 0x80000000, 0) != 1 || memmap.Uint32(0x5D4594, 1063104) != 0x80000000 {
		t.Fatal("byte index or negative zero")
	}
	for _, v := range []uint32{0, 5, 6, 7, 0xffffffff} {
		*o.words[1] = v
		legacy.PortTestUIInventoryCall(19, 0, 0, 0)
		want := v
		if v == 6 {
			want = 0
		}
		if *o.words[1] != want {
			t.Fatal("cancel changed unrelated mode")
		}
	}
	*o.meters.NamedWord("dword_8531A0_2576") = 0
	if legacy.PortTestUIInventoryCall(16, 0, 0, 0) != 1 {
		t.Fatal("missing player default")
	}
}
func (o *uiInventoryOwner) stack(index, count int) {
	o.slot(index, o.items[0], byte(count), 0)
	for i := 0; i < count; i++ {
		binary.LittleEndian.PutUint32(o.grid[index*148+4+i*4:], uint32(0x1000+i))
	}
}
func TestClientInventoryGridQueries(t *testing.T) {
	o := newUIInventoryOwner(t)
	var rows []uiInventoryResult
	for index := 0; index < 84; index++ {
		for _, count := range []int{0, 1, 2, 31, 32} {
			o.reset(t)
			o.stack(index, count)
			for _, code := range []uint32{0x1000, uint32(0x1000 + count - 1), 0xffffffff} {
				for _, op := range []int{3, 21, 22} {
					rows = append(rows, o.call(t, len(rows), op, code, 0, 0))
				}
			}
			for _, op := range []int{15, 25} {
				rows = append(rows, o.call(t, len(rows), op, o.items[0].TypeIDVal, 0, 0))
			}
			for _, op := range []int{24, 26} {
				rows = append(rows, o.call(t, len(rows), op, uint32(index/21), uint32(index%21), 0))
			}
		}
	}
	for _, x := range []uint32{0, 3, 4, 0xffffffff} {
		for _, y := range []uint32{0, 19, 20, 21, 0xffffffff} {
			o.reset(t)
			for _, op := range []int{24, 26} {
				rows = append(rows, o.call(t, len(rows), op, x, y, 0))
			}
		}
	}
	uiInventoryCapture(t, "inventory-grid", rows, "5538bf038ef38d4f95d62a4bd49b210992bea12c08568cd38656ddf9e77a9ae9")
}
func TestClientInventorySearchOrderContract(t *testing.T) {
	o := newUIInventoryOwner(t)
	o.reset(t)
	o.stack(20, 2)
	o.stack(21, 2)
	// Row0/column1 precedes row20/column0 despite column-major storage.
	if ret := legacy.PortTestUIInventoryCall(15, o.items[0].TypeIDVal, 0, 0); ret != o.cell(21) {
		t.Fatalf("type search order %#x", ret)
	}
	legacy.PortTestUIInventoryCall(3, 0x1001, 0, 0)
	if memmap.Uint32(0x5D4594, 1049788) != o.cell(21) || memmap.Uint32(0x5D4594, 1049792) != 1 {
		t.Fatal("stack search order/index")
	}
	clear(o.grid[21*148:][:148])
	if ret := legacy.PortTestUIInventoryCall(15, o.items[0].TypeIDVal, 0, 0); ret != o.cell(20) {
		t.Fatal("extra valid row excluded")
	}
	if legacy.PortTestUIInventoryCall(24, 0, 20, 0) != 0 || legacy.PortTestUIInventoryCall(26, 0, 20, 0) != 0 {
		t.Fatal("extra search row exposed as visible cell")
	}
	o.reset(t)
	o.stack(83, 32)
	if legacy.PortTestUIInventoryCall(3, 0x101f, 0, 0) == 0 || memmap.Uint32(0x5D4594, 1049788) != o.cell(83) || memmap.Uint32(0x5D4594, 1049792) != 31 {
		t.Fatal("last valid cell/code")
	}
	*memmap.PtrUint32(0x5D4594, 1049792) = 0xabcdef01
	if legacy.PortTestUIInventoryCall(3, 0xffffffff, 0, 0) != 0 || memmap.Uint32(0x5D4594, 1049792) != 0xabcdef01 {
		t.Fatal("missing code must preserve previous result")
	}
}
func (o *uiInventoryOwner) link(a, b int) {
	*(*uint32)(unsafe.Add(o.items[a].C(), 368)) = uint32(uintptr(o.items[b].C()))
}
func TestClientInventoryEquipmentQueries(t *testing.T) {
	o := newUIInventoryOwner(t)
	var rows []uiInventoryResult
	for slot := 0; slot < 9; slot++ {
		for _, length := range []int{0, 1, 2, 3} {
			for _, flag := range []uint32{0, 0x1000, 0x1000000, 0x1001000} {
				o.reset(t)
				if length > 0 {
					o.equipment[slot] = uint32(uintptr(o.items[0].C()))
				}
				for i := 0; i < 3; i++ {
					o.items[i].ObjClass = 0
					if i+1 < length {
						o.link(i, i+1)
					}
				}
				o.items[1].ObjClass = 0x10 | 0x1000
				*(*uint32)(unsafe.Add(o.items[1].C(), 112)) = flag
				for _, typ := range []uint32{o.items[0].TypeIDVal, o.items[1].TypeIDVal, o.items[2].TypeIDVal, 0xffffffff} {
					rows = append(rows, o.call(t, len(rows), 1, typ, 0, 0))
				}
				rows = append(rows, o.call(t, len(rows), 2, 0, 0, 0))
			}
		}
	}
	for _, length := range []int{0, 1, 2, 3} {
		for _, position := range []int{-1, 0, 1, 2} {
			o.reset(t)
			o.equipment[7] = uint32(uintptr(o.items[6].C()))
			if length > 0 {
				o.equipment[8] = uint32(uintptr(o.items[0].C()))
			}
			for i := 0; i < length; i++ {
				o.items[i].TypeIDVal = uint32(o.c.Things.TypeByID("RedPotion").Index())
				if i == position {
					o.items[i].TypeIDVal = o.items[5].TypeIDVal
				}
				if i+1 < length {
					o.link(i, i+1)
				}
			}
			rows = append(rows, o.call(t, len(rows), 0, 0, 0, 0))
		}
	}
	uiInventoryCapture(t, "inventory-equipment", rows, "4ddb3f204531b6c5f9ab646f9af2eb212a28d2ca0463b272f8c6496fd5bfaec5")
}
func TestClientInventoryBowHeadContract(t *testing.T) {
	o := newUIInventoryOwner(t)
	o.reset(t)
	o.equipment[7] = uint32(uintptr(o.items[6].C()))
	o.equipment[8] = uint32(uintptr(o.items[0].C()))
	o.link(0, 5)
	if legacy.PortTestUIInventoryCall(0, 0, 0, 0) != uint32(uintptr(o.items[0].C())) {
		t.Fatal("Bow within slot8 selects slot8 head")
	}
	*(*uint32)(unsafe.Add(o.items[0].C(), 368)) = 0
	if legacy.PortTestUIInventoryCall(0, 0, 0, 0) != o.equipment[7] {
		t.Fatal("non-Bow slot8 falls back to slot7")
	}
}
func TestClientInventoryItemUpdates(t *testing.T) {
	o := newUIInventoryOwner(t)
	var rows []uiInventoryResult
	for _, mode := range []int{0, 1, 2, 3} {
		for _, code := range []uint32{0, 0x1000, 0x1001, 0xffff} {
			for _, value := range []uint32{0, 1, 255, 65535, 65536, 0xffffffff} {
				o.reset(t)
				if mode == 0 || mode == 3 {
					o.stack(83, 2)
				}
				if mode == 1 || mode == 3 {
					*memmap.PtrUint32(0x5D4594, 1049848) = uint32(uintptr(o.items[1].C()))
					o.items[1].NetCode32 = 0x1000
				}
				for _, op := range []int{21, 22} {
					rows = append(rows, o.call(t, len(rows), op, code, 0, 0))
				}
				gridMatch := (mode == 0 || mode == 3) && (code == 0x1000 || code == 0x1001)
				if gridMatch {
					o.shortReturn = uint32(uintptr(o.items[0].C()))
				} else if (mode == 1 || mode == 3) && code != 0x1000 {
					o.shortReturn = uint32(uintptr(o.items[1].C()))
				}
				rows = append(rows, o.call(t, len(rows), 18, code, value, ^value))
				o.shortReturn = 0
				for _, linked := range []uint32{0, 1} {
					binary.LittleEndian.PutUint32(o.grid[83*148+132:], linked)
					rows = append(rows, o.call(t, len(rows), 30, code, value, ^value))
				}
			}
		}
	}
	uiInventoryCapture(t, "inventory-item-updates", rows, "c651316e97bd90ab39882450554b02b252bae199dcc5240b8a60b71224687fb8")
}
func TestClientInventoryUpdateContracts(t *testing.T) {
	o := newUIInventoryOwner(t)
	o.reset(t)
	o.stack(83, 2)
	legacy.PortTestUIInventoryCall(18, 0x1001, 0x12345678, 0xffffabcd)
	if *(*uint16)(unsafe.Add(o.items[0].C(), 292)) != 0x5678 || *(*uint16)(unsafe.Add(o.items[0].C(), 294)) != 0xabcd {
		t.Fatal("durability stores")
	}
	binary.LittleEndian.PutUint32(o.grid[83*148+132:], 1)
	legacy.PortTestUIInventoryCall(30, 0x1001, 12, 34)
	if o.meters.Records[5].Current != 12 || o.meters.Records[6].Maximum != 34 {
		t.Fatal("linked charge meter update")
	}
	if *(*uint16)(unsafe.Add(o.items[0].C(), 448)) != 12 || *(*uint16)(unsafe.Add(o.items[0].C(), 450)) != 34 {
		t.Fatal("item charge stores")
	}
}
func TestClientInventoryPotionUse(t *testing.T) {
	o := newUIInventoryOwner(t)
	var rows []uiInventoryResult
	for _, cursor := range []uint32{0, 1, 2, 8, 0xffffffff} {
		for _, pause := range []bool{false, true} {
			for _, present := range []bool{false, true} {
				o.reset(t)
				noxflags.ResetGame()
				if pause {
					noxflags.SetGame(noxflags.GamePause)
				}
				*memmap.PtrUint32(0x5D4594, 1096672) = cursor
				if present {
					o.stack(83, 2)
				}
				r := o.call(t, len(rows), 14, o.items[0].TypeIDVal, 0, 0)
				rows = append(rows, r)
				if cursor == 1 || pause || !present {
					if len(r.Messages) != 0 {
						t.Fatal("potion use ignored gate")
					}
				} else if len(r.Messages) != 1 || !bytes.Equal(r.Messages[0], []byte{116, 0, 16}) {
					t.Fatalf("potion message %x", r.Messages)
				}
			}
		}
	}
	uiInventoryCapture(t, "inventory-potion-use", rows, "febbfe95315831602eec7156ec170728bcb5c14244b0cb12fa84c3a6a49196ce")
}
func TestClientInventoryWeaponSelection(t *testing.T) {
	o := newUIInventoryOwner(t)
	var rows []uiInventoryResult
	for _, player := range []bool{false, true} {
		for _, mask := range []uint32{0, 2, 4, 6, 8, 0x7ffffff} {
			for _, equipped := range []bool{false, true} {
				for _, found := range []bool{false, true} {
					o.reset(t)
					*(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 4)) = mask
					if !player {
						*o.meters.NamedWord("dword_8531A0_2576") = 0
					}
					bow := o.items[5]
					bow.NetCode32 = 0x1000
					if equipped {
						o.equipment[7] = uint32(uintptr(bow.C()))
					}
					if found {
						o.slot(83, bow, 1, 0x1000)
					}
					*o.words[2] = o.cell(83)
					rows = append(rows, o.call(t, len(rows), 27, 0, 0, 0))
					r := o.call(t, len(rows), 29, 0, 0, 0)
					rows = append(rows, r)
					if player && mask&4 != 0 && equipped && found {
						if r.Return != o.norm(uint32(uintptr(bow.C()))) {
							t.Fatal("equipped Bow inventory selection")
						}
					} else if r.Return != 0 {
						t.Fatal("unexpected selected weapon")
					}
				}
			}
		}
	}
	uiInventoryCapture(t, "inventory-weapon-selection", rows, "442b4b21c61bfe5b3ade36b3eb2d818678da76c68d1488cdc2f874897c0b9d34")
}
