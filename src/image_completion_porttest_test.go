//go:build porttest

package opennox

import (
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
)

func TestImageCompletionTimingAndReference(t *testing.T) {
	o := newObjectRenderOwner(t)
	a := o.animation
	a.ImagesSz, a.Field_2_1, a.Field_3 = 3, 1, 10
	var refs []unsafe.Pointer
	key, restore := legacy.PortTestObserveImageEnd(func(p unsafe.Pointer) { refs = append(refs, p) })
	t.Cleanup(restore)
	a.OnEnd = key
	for _, tc := range []struct{ ts, frame, calls int }{
		{10, 0, 0}, {11, 0, 0}, {12, 1, 0}, {13, 1, 0},
		{14, 2, 1}, {15, 2, 2}, {16, 2, 3}, {100, 2, 4},
	} {
		got := o.c.getCursorAnimFrame(o.shiny, tc.ts-int(o.c.GetInputSeq()))
		if got != o.images[tc.frame] || len(refs) != tc.calls {
			t.Fatalf("ts=%d: frame=%p want=%p, callbacks=%d want=%d", tc.ts, got, o.images[tc.frame], len(refs), tc.calls)
		}
	}
	for i, ref := range refs {
		if ref != o.shiny.C() {
			t.Fatalf("callback %d received %p, want reference %p", i, ref, o.shiny.C())
		}
	}
}

func TestImageCompletionLookupAfterCallback(t *testing.T) {
	o := newObjectRenderOwner(t)
	a := o.animation
	a.ImagesSz = 3
	calls := 0
	key, restore := legacy.PortTestObserveImageEnd(func(p unsafe.Pointer) {
		calls++
		// The caller has already selected the slice, but has not read its cell.
		o.frames[2] = o.frames[7]
		a.ImagesPtr = (*noxrender.ImageHandle)(unsafe.Pointer(&o.frames[20]))
	})
	t.Cleanup(restore)
	a.OnEnd = key
	got := o.c.getCursorAnimFrame(o.shiny, 2-int(o.c.GetInputSeq()))
	if calls != 1 || got != o.images[7] {
		t.Fatalf("calls=%d, image=%p want=%p", calls, got, o.images[7])
	}
}

func TestImageCompletionNonCompletionPaths(t *testing.T) {
	o := newObjectRenderOwner(t)
	a := o.animation
	a.ImagesSz, a.Field_2_1 = 3, 1
	calls := 0
	key, restore := legacy.PortTestObserveImageEnd(func(unsafe.Pointer) { calls++ })
	t.Cleanup(restore)
	a.OnEnd = key
	if got := o.c.getCursorAnimFrame(nil, 100); got != nil {
		t.Fatal("nil reference produced image")
	}
	a.AnimType = 2
	a.Field_3 = 100 // Loop timing ignores the one-shot start sequence.
	for ts := 0; ts < 14; ts++ {
		if got := o.c.getCursorAnimFrame(o.shiny, ts-int(o.c.GetInputSeq())); got != o.images[(ts/2)%3] {
			t.Fatalf("loop frame at %d", ts)
		}
	}
	a.AnimType = 1
	if got := o.c.getCursorAnimFrame(o.shiny, 0); got != nil {
		t.Fatal("unsupported animation produced image")
	}
	if calls != 0 {
		t.Fatalf("unexpected completion calls: %d", calls)
	}
	a.AnimType, a.Field_3, a.OnEnd = 0, 0, nil
	if got := o.c.getCursorAnimFrame(o.shiny, 100-int(o.c.GetInputSeq())); got != o.images[2] {
		t.Fatal("nil callback changed last frame")
	}
}

func TestImageCompletionBookOwners(t *testing.T) {
	o := newObjectRenderOwner(t)
	words, restore := legacy.PortTestBookWords()
	t.Cleanup(restore)
	p, q := words["dword_5d4594_1046868"], words["dword_5d4594_1046872"]
	keys := legacy.PortTestBookCallbacks()
	spell, creature := keys["nox_xxx_bookClickSpell_45B1F0"], keys["nox_xxx_bookClickCreature_45B200"]
	if spell == nil || creature == nil || spell == creature {
		t.Fatal("completion identities missing or shared")
	}
	for _, tc := range []struct {
		key  unsafe.Pointer
		want uint32
	}{{spell, 0}, {creature, 1}} {
		o.animation.OnEnd = tc.key
		for i := 0; i < 2; i++ {
			*p, *q = 7, 9
			o.c.getCursorAnimFrame(o.shiny, 3-int(o.c.GetInputSeq()))
			if *p != tc.want || *q != tc.want {
				t.Fatalf("book completion: %d/%d want %d", *p, *q, tc.want)
			}
		}
	}
	if unsafe.Sizeof(*o.animation) != 16 || unsafe.Offsetof(o.animation.ImagesPtr) != 4 || unsafe.Offsetof(o.animation.Field_3) != 12 {
		t.Fatal("animation ABI changed")
	}
}
