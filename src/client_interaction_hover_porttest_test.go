//go:build porttest && !server

package opennox

import (
	"encoding/binary"
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

// Cursor hover uses the actual client visibility renderer, which the server build
// deliberately replaces with an unreachable stub. Geometry and selection run here
// on both production client variants, without substituting a visibility result.
func TestClientInteractionHoverCircleAndGates(t *testing.T) {
	o := newMeterOwner(t, "Polyp", "Glyph", "HoverItem")
	words, restore := legacy.PortTestClientInteractionWords()
	defer restore()
	local := serverConfigOwnBytes(t, 0x852978, 8, 4)
	depth := serverConfigOwnBytes(t, 0x5D4594, 1096628, 4)
	glyph := serverConfigOwnBytes(t, 0x5D4594, 1096632, 4)
	polyp := serverConfigOwnBytes(t, 0x5D4594, 1096648, 4)
	glyphID := uint32(o.c.Things.TypeByID("Glyph").Index())
	binary.LittleEndian.PutUint32(glyph, glyphID)
	restoreGlyph := legacy.PortTestClientInteractionGlyphCache(glyphID)
	defer restoreGlyph()
	dr, free := alloc.New(client.Drawable{})
	defer free()
	self, freeSelf := alloc.New(client.Drawable{})
	defer freeSelf()
	pl, freePlayer := alloc.New(server.Player{})
	defer freePlayer()
	point, freePoint := alloc.New([2]int32{})
	defer freePoint()
	ptr := func(d *client.Drawable) uint32 { return uint32(uintptr(d.C())) }
	reset := func() {
		*dr = client.Drawable{PosVec: image.Pt(100, 100), ObjClass: object.Class(0x200), ZSizeMax: 20, NetCode32: 0x1234, TypeIDVal: uint32(o.c.Things.TypeByID("HoverItem").Index()), DrawFuncPtr: legacy.PortTestSpriteAnimationCallback(2)}
		dr.Shape.Kind = server.ShapeKindCircle
		dr.Shape.Circle.R = 10
		*self = client.Drawable{}
		*pl = server.Player{}
		*words["dword_8531A0_2576"] = uint32(uintptr(unsafe.Pointer(pl)))
		binary.LittleEndian.PutUint32(local, ptr(self))
		clear(depth)
		clear(polyp)
		*words["dword_5d4594_1096640"] = 0
		*words["dword_5d4594_1096636"] = 0
		*words["nox_client_spriteUnderCursorXxx_1096644"] = 0
	}
	type row struct {
		Label              string
		X, Y               int32
		Hover, Usable      bool
		Depth, UsableDepth uint32
	}
	var captured []row
	check := func(label string, x, y int32, wantHover, wantUsable bool, wantDepth, wantUsableDepth uint32) {
		t.Helper()
		*point = [2]int32{x, y}
		interactionCall("nox_xxx_clientOnCursorHover_477050", uintptr(dr.C()), uintptr(unsafe.Pointer(point)))
		h, u := *words["dword_5d4594_1096640"], *words["nox_client_spriteUnderCursorXxx_1096644"]
		d, ud := binary.LittleEndian.Uint32(depth), *words["dword_5d4594_1096636"]
		if (h == ptr(dr)) != wantHover || (u == ptr(dr)) != wantUsable || d != wantDepth || ud != wantUsableDepth {
			t.Fatalf("hover %s (%d,%d): selected%v/%v depths%d/%d want%v/%v %d/%d", label, x, y, h == ptr(dr), u == ptr(dr), d, ud, wantHover, wantUsable, wantDepth, wantUsableDepth)
		}
		if binary.LittleEndian.Uint32(polyp) != uint32(o.c.Things.TypeByID("Polyp").Index()) {
			t.Fatal("Polyp cache before gates")
		}
		captured = append(captured, row{label, x, y, h == ptr(dr), u == ptr(dr), d, ud})
	}
	for y := int32(68); y <= 112; y++ {
		for x := int32(88); x <= 112; x++ {
			reset()
			// Independent capsule contract: distance to the vertical segment.
			nearestY := max(int32(80), min(int32(100), y))
			dx, dy := int(x-100), int(y-nearestY)
			hit := dx*dx+dy*dy < 100
			d := uint32(0)
			if hit {
				d = 100
			}
			check("circle", x, y, hit, hit, d, d)
		}
	}
	// A 20-pixel-tall diamond prism: two strict diamond caps plus the side band.
	abs := func(v int32) int32 {
		if v < 0 {
			return -v
		}
		return v
	}
	for y := int32(68); y <= 112; y++ {
		for x := int32(88); x <= 112; x++ {
			reset()
			dr.Shape.Kind = server.ShapeKindBox
			dr.Shape.Box = server.ShapeBox{W: 14.142272, H: 14.142272, LeftTop: 0, LeftBottom: -10, LeftBottom2: -10, LeftTop2: 0, RightTop: 10, RightBottom: 0, RightBottom2: 0, RightTop2: 10}
			dx := abs(x - 100)
			hit := dx+abs(y-80) < 10 || dx+abs(y-100) < 10 || dx < 10 && y > 80 && y < 100
			d := uint32(0)
			if hit {
				d = 100
			}
			check("box", x, y, hit, hit, d, d)
		}
	}

	for _, g := range []struct {
		name          string
		change        func()
		hover, usable bool
	}{
		{"ordinary", func() {}, true, true},
		{"self", func() { binary.LittleEndian.PutUint32(local, ptr(dr)) }, false, false},
		{"hidden", func() { dr.ObjFlags = 0x8000 }, false, false},
		{"invisible", func() { dr.Buffs = 1 }, false, false},
		{"see-invisible", func() { dr.Buffs = 1; self.Buffs = 1 << 21 }, true, true},
		{"monster-excluded", func() { dr.ObjClass = 2; dr.ObjSubClass = 0x4000 }, false, false},
		{"uninteresting-class", func() { dr.ObjClass = 0 }, false, false},
		{"Polyp-exception", func() { dr.ObjClass = 0; dr.TypeIDVal = uint32(o.c.Things.TypeByID("Polyp").Index()) }, true, true},
		{"no-draw-function", func() { dr.DrawFuncPtr = nil }, false, false},
		{"player-info-missing", func() { dr.ObjClass = 4 }, false, false},
		{"static-ineligible", func() { dr.ObjClass = 0x400000 }, false, false},
		{"static-eligible", func() { dr.ObjClass = 0x400000; dr.ObjSubClass = 0x80 }, true, true},
		{"monster-anim10", func() { dr.ObjClass = 2; dr.AnimInd = 10 }, false, false},
		{"monster-anim11", func() { dr.ObjClass = 2; dr.AnimInd = 11 }, true, true},
		{"no-shape", func() { dr.Shape.Kind = 0 }, false, false},
		{"no-player", func() { *words["dword_8531A0_2576"] = 0 }, true, false},
		{"glyph-warrior", func() { dr.TypeIDVal = glyphID }, true, false},
		{"glyph-wizard", func() { dr.TypeIDVal = glyphID; *(*byte)(unsafe.Add(unsafe.Pointer(pl), 2251)) = 1 }, true, true},
	} {
		reset()
		g.change()
		d, ud := uint32(0), uint32(0)
		if g.hover {
			d = 100
		}
		if g.usable && g.name != "glyph-wizard" {
			ud = 100
		}
		check(g.name, 100, 90, g.hover, g.usable, d, ud)
	}
	// Equal depth must keep the first selection; greater depth replaces both.
	reset()
	*point = [2]int32{100, 90}
	interactionCall("nox_xxx_clientOnCursorHover_477050", uintptr(dr.C()), uintptr(unsafe.Pointer(point)))
	*self = *dr
	self.ZVal = 0
	binary.LittleEndian.PutUint32(local, 0)
	interactionCall("nox_xxx_clientOnCursorHover_477050", uintptr(self.C()), uintptr(unsafe.Pointer(point)))
	if *words["dword_5d4594_1096640"] != ptr(dr) || *words["nox_client_spriteUnderCursorXxx_1096644"] != ptr(dr) {
		t.Fatal("equal hover depth changed first selection")
	}
	self.ZVal = 1
	interactionCall("nox_xxx_clientOnCursorHover_477050", uintptr(self.C()), uintptr(unsafe.Pointer(point)))
	if *words["dword_5d4594_1096640"] != ptr(self) || *words["nox_client_spriteUnderCursorXxx_1096644"] != ptr(self) || binary.LittleEndian.Uint32(depth) != 101 || *words["dword_5d4594_1096636"] != 101 {
		t.Fatal("greater hover depth did not replace selection")
	}
	interactionCapture(t, "hover-circle-gates", captured)
}

// The effects fixture supplies a synthetic mouse getter; this integration check
// delegates that boundary to the actual owned input device.
type interactionHoverMouseClient struct {
	legacy.Client
	input interface{ GetMousePos() image.Point }
}

func (c *interactionHoverMouseClient) GetMousePos() image.Point { return c.input.GetMousePos() }

func TestClientInteractionHoverEnumeration(t *testing.T) {
	o := newInventoryTransactionOwner(t, "Polyp", "Glyph", "HoverItem")
	o.reset(t)
	oldClient := legacy.GetClient
	proxy := &interactionHoverMouseClient{oldClient(), o.c.Inp}
	legacy.GetClient = func() legacy.Client { return proxy }
	defer func() { legacy.GetClient = oldClient }()
	words, restore := legacy.PortTestClientInteractionWords()
	defer restore()
	serverConfigOwnBytes(t, 0x852978, 8, 4)
	depth := serverConfigOwnBytes(t, 0x5D4594, 1096628, 4)
	glyph := serverConfigOwnBytes(t, 0x5D4594, 1096632, 4)
	serverConfigOwnBytes(t, 0x5D4594, 1096648, 4)
	restoreTypes := o.c.srv.Server.PortTestRewardTypes([]string{"Glyph"}, nil, true, 0, 0)
	defer restoreTypes()
	restoreGlyph := legacy.PortTestClientInteractionGlyphCache(uint32(o.c.Things.TypeByID("Glyph").Index()))
	defer restoreGlyph()
	if o.c.srv.Types.IndByID("Glyph") == 0 {
		t.Fatal("nonempty server Glyph lookup")
	}
	oldOffset := partViewportOff
	defer func() { partViewportOff = oldOffset }()
	o.c.Inp.SetMouseBounds(image.Rect(0, 0, 256, 256))
	var objects []*client.Drawable
	for i := 0; i < 4; i++ {
		dr := o.item(t, "HoverItem", uint32(i+1))
		dr.ObjClass = 0x200
		dr.ObjFlags = 0
		dr.Buffs = 0
		dr.DrawFuncPtr = legacy.PortTestSpriteAnimationCallback(2)
		dr.PosVec = image.Pt(100, 100)
		dr.ZSizeMin = 0
		dr.ZSizeMax = 20
		dr.Shape.Kind = server.ShapeKindCircle
		dr.Shape.Circle.R = 10
		if i == 1 || i == 2 {
			dr.ZVal = 1
		}
		if i == 3 {
			dr.PosVec = image.Pt(500, 500)
			dr.Shape.Circle.R = 1000
		}
		objects = append(objects, dr)
	}
	var indexed []*client.Drawable
	remove := func() {
		for _, dr := range indexed {
			o.c.Objs.Index2DRemove(dr, dr.GetExt())
		}
		indexed = nil
	}
	defer remove()
	type row struct {
		Reverse       bool
		Offset        image.Point
		Hover, Usable int
		Depth         uint32
		Glyph         int
	}
	var captured []row
	for _, reverse := range []bool{false, true} {
		for _, offset := range []image.Point{{0, 0}, {20, 20}} {
			remove()
			order := []int{0, 1, 2, 3}
			want := 2
			if reverse {
				order = []int{0, 2, 1, 3}
				want = 1
			}
			for _, i := range order {
				o.c.Objs.AddIndex2D(objects[i])
				indexed = append(indexed, objects[i])
			}
			partViewportOff = offset
			o.c.Inp.ChangeMousePos(image.Pt(100, 90).Sub(offset), true)
			binary.LittleEndian.PutUint32(depth, 0x7fffffff)
			*words["dword_5d4594_1096636"] = 0x7fffffff
			*words["dword_5d4594_1096640"] = uint32(uintptr(objects[3].C()))
			*words["nox_client_spriteUnderCursorXxx_1096644"] = uint32(uintptr(objects[3].C()))
			clear(glyph)
			interactionCall("nox_xxx_clientEnumHover_476FA0")
			selected := uint32(uintptr(objects[want].C()))
			if *words["dword_5d4594_1096640"] != selected || *words["nox_client_spriteUnderCursorXxx_1096644"] != selected || binary.LittleEndian.Uint32(depth) != 101 || *words["dword_5d4594_1096636"] != 101 || binary.LittleEndian.Uint32(glyph) != uint32(o.c.srv.Types.IndByID("Glyph")) {
				index := func(v uint32) int {
					for i, dr := range objects {
						if v == uint32(uintptr(dr.C())) {
							return i
						}
					}
					return -1
				}
				var visits []int
				o.c.Objs.EachInRect(image.Rect(4, -6, 196, 186), func(dr *client.Drawable) { visits = append(visits, index(uint32(uintptr(dr.C())))) })
				t.Fatal("hover enumeration spatial bound/order/reset/translation", reverse, offset, "selected", index(*words["dword_5d4594_1096640"]), index(*words["nox_client_spriteUnderCursorXxx_1096644"]), "depth", binary.LittleEndian.Uint32(depth), *words["dword_5d4594_1096636"], "glyph", binary.LittleEndian.Uint32(glyph), o.c.srv.Types.IndByID("Glyph"), "mouse", o.c.Inp.GetMousePos(), "visits", visits)
			}
			captured = append(captured, row{reverse, offset, want, want, 101, o.c.srv.Types.IndByID("Glyph")})
		}
	}
	remove()
	interactionCall("nox_xxx_clientEnumHover_476FA0")
	if *words["dword_5d4594_1096640"] != 0 || *words["nox_client_spriteUnderCursorXxx_1096644"] != 0 || binary.LittleEndian.Uint32(depth) != 0 || *words["dword_5d4594_1096636"] != 0 {
		t.Fatal("empty hover enumeration clears previous selection")
	}
	interactionCapture(t, "hover-enumeration", captured)
}
