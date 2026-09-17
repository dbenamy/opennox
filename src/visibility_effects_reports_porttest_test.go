//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestVisibilityEffectsSpecialUpdates(t *testing.T) {
	s, units, configure := visibilityEffectsPlayers(t)
	configure(3)
	u := newCreatureXferObject(t, s, "Monster")
	u.NetCode = 0x1234
	u.Field33 = 0x12345678
	u.Field5 = 0xABCDEF00
	u.Buffs = 0x11223344
	u.ZVal = 12.75
	u.ObjClass = object.ClassMonster | object.ClassImmobile | object.ClassVisibleEnable
	u.Extent = 0x2345
	u.TeamPtr().ID = 1
	for i := range units {
		units[i].TeamPtr().ID = 1
		units[i].TypeInd = u.TypeInd
	}
	colors := unsafe.Slice((*byte)(unsafe.Add(u.UpdateData, 2076)), 18)
	for i := range colors {
		colors[i] = byte(i + 1)
	}
	cache := memmap.PtrUint32(0x5D4594, 1556320)
	old := *cache
	*cache = uint32(u.TypeInd)
	t.Cleanup(func() { *cache = old })
	reset, snapshot, free := legacy.PortTestVisibilityReliableOwner()
	t.Cleanup(free)
	type row struct {
		Name          string
		Return, Flags uint32
		Packets       []legacy.PortTestShopPacketResult
		Direct        [][]byte
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "visibility-effects-special-updates", rows, "0f9647778519617c8a159e4584ef7e32bbbea71e49e0e1fa0ab055b7fe6f2cc6")
	}()
	flags := []uint32{0, 0xffffffff, 0xFF0F0000, 0x04CF0000}
	for i := 0; i < 32; i++ {
		flags = append(flags, 1<<i)
	}
	for _, slot := range []int{0, 1, 2} {
		for _, flag := range flags {
			name := fmt.Sprintf("player%d/flags%x", slot, flag)
			t.Run(name, func(t *testing.T) {
				reset()
				s.NetList.ResetAll()
				u.ObjClass = object.ClassMonster | object.ClassImmobile
				if flag&0x40000 != 0 {
					u.ObjClass |= object.ClassVisibleEnable
				}
				pl := units[slot].UpdateDataPlayer().Player
				ind := int(pl.PlayerInd)
				dst := (*uint32)(unsafe.Add(u.CObj(), 560+4*ind))
				*dst = flag
				rv := legacy.PortTestVisibilityEffects(12, &units[slot], u, nil, nil, [5]int32{}, nil, "")
				wantRV := uint32(1)
				if flag&0x0FFF0000 == 0 {
					wantRV = 0
				}
				if rv != wantRV || *dst != flag&^uint32(0x04CF0000|0x02000000) {
					t.Fatalf("return%d flags%x", rv, *dst)
				}
				code := uint16(0xA345)
				packet := func(op byte, tail ...byte) []byte { return append([]byte{op, byte(code), byte(code >> 8)}, tail...) }
				var expected [][]byte
				if flag&0x10000 != 0 {
					b := packet(107)
					b = binary.LittleEndian.AppendUint32(b, u.Field33)
					expected = append(expected, b)
				}
				if flag&0x20000 != 0 {
					expected = append(expected, packet(65, 40))
				}
				if flag&0x40000 != 0 {
					expected = append(expected, packet(56))
				}
				if flag&0x80000 != 0 {
					b := packet(101)
					b = binary.LittleEndian.AppendUint32(b, u.Field5)
					expected = append(expected, b)
				}
				if flag&0x800000 != 0 {
					b := packet(90)
					b = binary.LittleEndian.AppendUint32(b, u.Buffs)
					expected = append(expected, b)
				}
				if flag&0x2000000 != 0 {
					expected = append(expected, packet(103, 255, 255, 255, 255))
				}
				if flag&0x4000000 != 0 {
					expected = append(expected, append([]byte{105, 0x34, 0x12}, colors...))
				}
				got := snapshot()
				if len(got) != len(expected) {
					t.Fatalf("packets%d want%d", len(got), len(expected))
				}
				for i, p := range got {
					want := expected[len(expected)-1-i]
					if int(p.Recipient) != ind || !bytes.Equal(p.Data, want) {
						t.Fatalf("packet%d got%+v want%x", i, p, want)
					}
				}
				direct := visibilityEffectsPackets(s)
				for i, b := range direct {
					want := []byte{}
					if i == ind && flag&0x400000 != 0 {
						want = packet(94, 12)
					}
					if !bytes.Equal(b, want) {
						t.Fatalf("direct%d got%x want%x", i, b, want)
					}
				}
				rows = append(rows, row{name, rv, *dst, got, direct})
			})
		}
	}
}
func TestVisibilityEffectsDestroyReports(t *testing.T) {
	s, _, configure := visibilityEffectsPlayers(t)
	u := newCreatureXferObject(t, s, "Monster")
	u.NetCode = 0x1234
	reset, snapshot, free := legacy.PortTestVisibilityReliableOwner()
	t.Cleanup(free)
	type row struct {
		Name    string
		Packets []legacy.PortTestShopPacketResult
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "visibility-effects-destroy-reports", rows, "04de6f0dc161424687a0ea32c73951740ff4e63ea4b687459de41673b360e38b")
	}()
	for _, n := range []int{0, 1, 3} {
		configure(n)
		for _, class := range []object.Class{object.ClassSimple, object.ClassMonster, object.ClassPlayer} {
			for _, mask := range []uint32{0, 2, 8, 0x80000000, 0xffffffff} {
				for _, field := range []uint32{0, 64, 128, 255, 0xFFFFFF00} {
					name := fmt.Sprintf("n%d/class%x/mask%x/field%x", n, class, mask, field)
					t.Run(name, func(t *testing.T) {
						reset()
						u.ObjClass = class
						u.Field37 = mask
						u.Field5 = field
						legacy.PortTestVisibilityEffects(25, u, nil, nil, nil, [5]int32{}, nil, "")
						type pkt struct {
							to   byte
							data []byte
						}
						var want []pkt
						for _, slot := range []int{1, 3, 7, 31} {
							active := slot == 3 || n >= 1 && slot == 1 || n >= 2 && slot == 7 || n >= 3 && slot == 31
							if !active {
								continue
							}
							if mask&(1<<slot) != 0 {
								want = append(want, pkt{byte(slot), []byte{byte(field)>>6 | 49, 0x34, 0x12}})
							}
							if uint32(class)&6 != 0 {
								want = append(want, pkt{byte(slot), []byte{53, 0x34, 0x12}})
							}
						}
						got := snapshot()
						if len(got) != len(want) {
							t.Fatalf("packets%d want%d", len(got), len(want))
						}
						for i, p := range got {
							w := want[len(want)-1-i]
							if p.Recipient != w.to || !bytes.Equal(p.Data, w.data) {
								t.Fatalf("packet%+v want%+v", p, w)
							}
						}
						rows = append(rows, row{name, got})
					})
				}
			}
		}
	}
}

func TestVisibilityEffectsSpecialGates(t *testing.T) {
	s, units, _ := visibilityEffectsPlayers(t)
	u := newCreatureXferObject(t, s, "Monster")
	pl := units[0].UpdateDataPlayer().Player
	for _, ind := range []byte{0, 1, 31, 32, 128, 255} {
		for _, present := range []bool{false, true} {
			t.Run(fmt.Sprintf("index%d/present%t", ind, present), func(t *testing.T) {
				old := pl.PlayerInd
				pl.PlayerInd = ind
				defer func() { pl.PlayerInd = old }()
				target := u
				if !present {
					target = nil
				}
				before := bytes.Clone(unsafe.Slice((*byte)(u.CObj()), int(unsafe.Sizeof(*u))))
				rv := legacy.PortTestVisibilityEffects(12, &units[0], target, nil, nil, [5]int32{}, nil, "")
				want := uint32(1)
				if present && ind < 32 {
					want = 0
				}
				if rv != want {
					t.Fatalf("return%d want%d", rv, want)
				}
				if !bytes.Equal(before, unsafe.Slice((*byte)(u.CObj()), len(before))) {
					t.Fatal("object mutated")
				}
			})
		}
	}
}

func TestVisibilityEffectsEnemyHealth(t *testing.T) {
	s, units, _ := visibilityEffectsPlayers(t)
	u := newCreatureXferObject(t, s, "Monster")
	units[0].TypeInd = u.TypeInd
	units[0].TeamPtr().ID = 1
	u.TeamPtr().ID = 1
	u.ObjClass = object.ClassMonsterGenerator
	defer func() { u.ObjClass = object.ClassMonster }()
	reset, snapshot, free := legacy.PortTestVisibilityReliableOwner()
	t.Cleanup(free)
	reset()
	if !s.IsEnemyTo(&units[0], u) {
		t.Fatal("expected real generator enemy classification")
	}
	dst := (*uint32)(unsafe.Add(u.CObj(), 564))
	*dst = 0x20000
	rv := legacy.PortTestVisibilityEffects(12, &units[0], u, nil, nil, [5]int32{}, nil, "")
	if rv != 1 || *dst != 0 || len(snapshot()) != 0 {
		t.Fatalf("return%d flags%x packets%v", rv, *dst, snapshot())
	}
}
