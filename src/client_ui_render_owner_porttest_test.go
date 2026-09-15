//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
)

// Remaining C UI callers may select a different RenderData owner from R2.Data().
func TestClientUIRenderSeparateClipOwner(t *testing.T) {
	o := newObjectRenderOwner(t)
	active := o.c.r.Data()
	before := *active
	alternate, free := noxrender.NewRenderData()
	defer free()
	alternate.SetRect3(image.Rect(0, 0, 80, 80))
	*(*uint32)(alternate.C()) = 0x13572468
	env := legacy.PortTestNewObjectRenderEnvironment(alternate)
	defer env.Restore()
	if ret := legacy.PortTestUIRenderOp(0, [6]int32{-1}); ret != 0x13572468 {
		t.Fatal("alternate clip flag prior value")
	}
	if ret := legacy.PortTestUIRenderOp(1, [6]int32{10, 20, 30, 40}); ret != 1 {
		t.Fatal("alternate copy return")
	}
	if alternate.ClipRect() != image.Rect(10, 20, 40, 60) || alternate.ClipRect2() != image.Rect(10, 20, 39, 59) {
		t.Fatal("copy did not update selected C owner")
	}
	if ret := legacy.PortTestUIRenderOp(2, [6]int32{15, 35}); ret != 1 {
		t.Fatal("alternate narrowing return")
	}
	if alternate.ClipRect() != image.Rect(15, 20, 35, 60) || alternate.ClipRect2() != image.Rect(15, 20, 34, 59) {
		t.Fatal("narrowing used the wrong current owner")
	}
	legacy.PortTestObjectRenderClip(true)
	legacy.PortTestUIRenderOp(0, [6]int32{0})
	legacy.PortTestUIRenderOp(1, [6]int32{1, 2, 3, 4})
	if ret := legacy.PortTestObjectRenderClip(false); ret != 34 {
		t.Fatal("alternate restore return")
	}
	if *(*uint32)(alternate.C()) != 0xffffffff || alternate.ClipRect() != image.Rect(15, 20, 35, 60) || alternate.ClipRect2() != image.Rect(15, 20, 34, 59) {
		t.Fatal("save/restore lost selected C state")
	}
	if *active != before {
		t.Fatal("C clipping modified the distinct Go renderer state")
	}
	// Prove these are separate allocated owners, rather than two names for one buffer.
	if active == alternate {
		t.Fatal("fixture did not separate render owners")
	}
}
