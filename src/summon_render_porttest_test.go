//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"testing"
	"unsafe"
)

func TestSummonRendering(t *testing.T) {
	o := newSummonOwner(t)
	var rows []summonResult
	for _, size := range []uint32{1, 2, 4} {
		for _, icon := range []bool{false, true} {
			for _, health := range [][3]uint16{{0, 0, 0}, {0, 100, 0}, {50, 100, 0}, {100, 100, 1}, {65535, 65535, 1}} {
				for _, flash := range []uint32{0, 1} {
					o.reset(t)
					o.init(t)
					*o.summonWords["dword_5d4594_1320992"] = 0
					*o.word(1321200) = 2
					r := o.record(0)
					r[0] = 71
					r[1] = uint32(o.c.Things.TypeByID("PortSmallCreature").Index())
					r[2] = 1
					r[5] = size
					r[6] = 0
					r[7] = flash
					o.call("sub_4C2F70")
					if icon {
						*memmap.PtrUint32(0x5D4594, 740116) = uint32(uintptr(o.images[1].C()))
					}
					h := o.auxiliary[0]
					*(*uint32)(unsafe.Pointer(&h[0])) = 71
					h[4] = byte(health[2])
					*(*uint16)(unsafe.Pointer(&h[6])) = health[0]
					*(*uint16)(unsafe.Pointer(&h[8])) = health[1]
					*(*uint32)(unsafe.Pointer(&h[12])) = 1
					ret := o.call("nox_xxx_guiDrawSummonBox_4C1FE0", *o.summonWords["dword_5d4594_1321036"])
					o.check(t, r[7] == 0, "draw consumes flash once")
					rows = append(rows, o.snapshot(fmt.Sprintf("size%d-icon%t-health%v-flash%d", size, icon, health, flash), ret))
				}
			}
		}
	}
	spellbookCapture(t, "summon-rendering", rows, "8e9e47604f6f06526301a15e8f338526993ce3f3912e0f3274f1cb0db3845e6a")
}
func TestSummonHighlights(t *testing.T) {
	o := newSummonOwner(t)
	var rows []summonResult
	for _, mouse := range []image.Point{{0, 0}, {552, 40}, {577, 14}, {639, 479}} {
		for _, prior := range []uint32{0, 1} {
			o.reset(t)
			o.init(t)
			*o.summonWords["dword_5d4594_1320992"] = 0
			*o.word(1321200) = 2
			*o.word(1321212) = prior
			o.c.Mouse = mouse
			dr := o.item(t, "PortSmallCreature", 71)
			old := o.c.Objs.List1
			dr.NextPtr = old
			o.c.Objs.List1 = dr
			*txword(dr, 120) = 0x40000042
			r := o.record(0)
			r[0] = 71
			r[1] = dr.TypeIDVal
			r[2] = 1
			r[5] = 1
			r[6] = 0
			o.call("sub_4C2F70")
			ret := o.call("nox_xxx_guiDrawSummonBox_4C1FE0", *o.summonWords["dword_5d4594_1321036"])
			got := o.snapshot(fmt.Sprintf("mouse%v-prior%d", mouse, prior), ret)
			got.SpriteFlags = []uint32{*txword(dr, 120)}
			rows = append(rows, got)
			o.check(t, *txword(dr, 120)&^uint32(0x40000000) == 0x42, "hover preserves other sprite flags")
			o.check(t, o.call("nox_xxx_sprite_4C3220", uint32(uintptr(dr.C()))) == 1, "drawable membership uses code")
			o.call("nox_xxx_cliSummonOnDieOrBanish_4C3140", 71, 1)
			o.check(t, *txword(dr, 120) == 0x42, "removal clears highlight")
			got = o.snapshot(fmt.Sprintf("mouse%v-prior%d-remove", mouse, prior), 0)
			got.SpriteFlags = []uint32{*txword(dr, 120)}
			rows = append(rows, got)
			o.c.Objs.List1 = old
			dr.NextPtr = nil
		}
	}
	spellbookCapture(t, "summon-highlights", rows, "c042bd0b66ecd82c655cb554bcf1bcdbd2bc1fa5447998a2d902eac777848bc2")
}
func TestSummonControlEvents(t *testing.T) {
	o := newSummonOwner(t)
	var rows []summonResult
	for _, op := range []string{"nox_xxx_wndSummonBigButtonProc_4C24B0", "nox_xxx_wndSummonProc_4C2B10", "nox_xxx_clientOrderCreature_4C2A60"} {
		for _, event := range []uint32{0, 5, 6, 7, 17, 18, 0xffffffff} {
			for _, selected := range []bool{false, true} {
				o.reset(t)
				o.init(t)
				r := o.record(0)
				r[0] = 71
				r[1] = uint32(o.c.Things.TypeByID("PortSmallCreature").Index())
				r[2] = 1
				r[5] = 1
				r[6] = 0
				o.call("sub_4C2F70")
				if selected {
					*o.summonWords["dword_5d4594_1321204"] = o.ptr(0)
				}
				w := *o.summonWords["dword_5d4594_1321040"]
				if op == "nox_xxx_wndSummonProc_4C2B10" {
					w = *o.summonWords["dword_5d4594_1321036"]
				}
				// Window word 8 (widget data) selects the command in the third handler.
				(*[101]uint32)(unsafe.Pointer(uintptr(w)))[8] = 4
				ret := o.call(op, w, event, uint32(560)|(uint32(0)<<16))
				want := uint32(0)
				if event >= 5 && event <= 7 {
					want = 1
				}
				o.check(t, ret == want, "event dispatch return")
				rows = append(rows, o.snapshot(fmt.Sprintf("%s-event%d-selected%t", op, event, selected), ret))
			}
		}
	}
	for _, op := range []string{"sub_4C24A0", "sub_4C2BD0", "sub_4C2BE0"} {
		o.reset(t)
		ret := o.call(op)
		want := uint32(1)
		if op == "sub_4C2BD0" {
			want = 0
		}
		o.check(t, ret == want, "constant callback return")
		rows = append(rows, o.snapshot(op, ret))
	}
	spellbookCapture(t, "summon-control-events", rows, "f2feb5303724b631fb677d971be272b55572d576b33e646a5d3222494ee2a9cc")
}
func TestSummonTooltipsAndMenuDraw(t *testing.T) {
	o := newSummonOwner(t)
	var rows []summonResult
	pos, free := alloc.New([2]int32{})
	defer free()
	p := uint32(uintptr(unsafe.Pointer(pos)))
	for _, command := range []uint32{0, 1, 2, 3, 4, 5, 6, 0xffffffff} {
		o.reset(t)
		o.init(t)
		ret := o.call("sub_4C1CA0", command)
		rows = append(rows, o.snapshot(fmt.Sprintf("command%d-icon", command), ret))
		ret = o.call("sub_4C2CE0")
		rows = append(rows, o.snapshot(fmt.Sprintf("command%d-tip", command), ret))
	}
	for _, selected := range []int{-1, 0, 1} {
		for _, hover := range []bool{false, true} {
			for _, command := range []uint32{0, 1, 3, 4, 5} {
				o.reset(t)
				o.init(t)
				r := o.record(0)
				r[0] = 71
				r[1] = uint32(o.c.Things.TypeByID("PortSmallCreature").Index())
				if selected == 1 {
					r[1] = uint32(o.c.Things.TypeByID("CarnivorousPlant").Index())
				}
				r[2] = 1
				r[5] = 1
				r[6] = 0
				o.call("sub_4C2F70")
				if selected >= 0 {
					*o.summonWords["dword_5d4594_1321204"] = o.ptr(0)
				}
				*pos = [2]int32{100, 100}
				o.call("nox_xxx_wndSummonCreateList_4C2560", p)
				o.collect()
				var win *gui.Window
				for _, w := range o.windows {
					if w.Parent() != nil && uint32(uintptr(w.Parent().C())) == *o.summonWords["dword_5d4594_1321044"] && (*[101]uint32)(w.C())[8] == command {
						win = w
						break
					}
				}
				if win == nil {
					t.Fatal("menu command child missing")
				}
				o.c.Mouse = image.Pt(0, 0)
				if hover {
					o.c.Mouse = win.GlobalPos().Add(image.Pt(1, 1))
				}
				ret := o.call("sub_4C27F0", uint32(uintptr(win.C())))
				rows = append(rows, o.snapshot(fmt.Sprintf("selection%d-hover%t-command%d", selected, hover, command), ret))
			}
		}
	}
	spellbookCapture(t, "summon-tooltips-menu-draw", rows, "98464c50265bab0771efb8b9e62db94bbb6bbb959df486b64442a8f91c25e959")
}

func TestSummonMenuBounds(t *testing.T) {
	o := newSummonOwner(t)
	var rows []summonResult
	pos, free := alloc.New([2]int32{})
	defer free()
	p := uint32(uintptr(unsafe.Pointer(pos)))
	for _, screen := range [][2]uint32{{640, 480}, {1024, 768}, {480, 640}} {
		for _, corner := range [][2]int32{{0, 0}, {1, 0}, {0, 1}, {1, 1}} {
			o.reset(t)
			*o.words["nox_win_width"] = screen[0]
			*o.words["nox_win_height"] = screen[1]
			o.init(t)
			*pos = [2]int32{corner[0] * int32(screen[0]-1), corner[1] * int32(screen[1]-1)}
			o.call("nox_xxx_wndSummonCreateList_4C2560", p)
			w := (*gui.Window)(unsafe.Pointer(uintptr(*o.summonWords["dword_5d4594_1321044"])))
			r := image.Rectangle{Min: w.GlobalPos(), Max: w.GlobalPos().Add(w.Size())}
			if r.Min.X < 0 || r.Min.Y < 0 || r.Max.X > int(screen[0]) || r.Max.Y > int(screen[1]) {
				t.Fatalf("summon menu %v outside screen %v at corner %v", r, screen, corner)
			}
			rows = append(rows, o.snapshot(fmt.Sprintf("screen%v-corner%v", screen, corner), 0))
		}
	}
	spellbookCapture(t, "summon-menu-bounds", rows, "a0b1e77a473097af5763706bdab66cd195d1a187ae467012c0a7e19dc9ac063c")
}

func TestSummonSlotTooltips(t *testing.T) {
	o := newSummonOwner(t)
	var rows []summonResult
	pos, free := alloc.New([2]int32{})
	defer free()
	p := uint32(uintptr(unsafe.Pointer(pos)))
	for _, xy := range [][2]int32{{-39, 0}, {-38, 0}, {-37, 0}, {-1, 0}, {0, 0}, {37, 37}, {38, 0}, {75, 75}, {76, 76}} {
		for _, active := range []bool{false, true} {
			o.reset(t)
			o.init(t)
			w := (*gui.Window)(unsafe.Pointer(uintptr(*o.summonWords["dword_5d4594_1321036"])))
			w.SetPos(image.Pt(100, 100))
			if active {
				r := o.record(0)
				r[0] = 71
				r[1] = uint32(o.c.Things.TypeByID("PortSmallCreature").Index())
				r[2] = 1
				r[5] = 1
				r[6] = 0
				o.call("sub_4C2F70")
			}
			base := w.GlobalPos()
			*pos = [2]int32{int32(base.X) + xy[0], int32(base.Y) + xy[1]}
			ret := o.call("sub_4C2C60", uint32(uintptr(w.C())), p)
			inside := xy[0]/38 == 0 && xy[1]/38 == 0
			o.check(t, (ret != 0) == (active && inside), "tooltip uses C truncating division and selected cell")
			// Capture the returned text via the real tooltip callback instead of a raw string pointer.
			packed := uint32(uint16(pos[0])) | uint32(uint16(pos[1]))<<16
			ret = o.call("sub_4C2C20", uint32(uintptr(w.C())), 0, packed)
			rows = append(rows, o.snapshot(fmt.Sprintf("xy%v-active%t", xy, active), ret))
		}
	}
	spellbookCapture(t, "summon-slot-tooltips", rows, "3ec72d5e87eed2f6d367a0030270f028d122ddb75bee60d9fd16781d4175ce43")
}
