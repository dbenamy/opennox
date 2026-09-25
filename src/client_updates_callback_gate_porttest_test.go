//go:build porttest

package opennox

import (
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
)

// Preserve the original Linux/386 callback gate, including the zero machine
// result captured from the old magic-trail void export.
func TestClientUpdatesCallbackGate(t *testing.T) {
	c, pix, effects, env := newUpdateTestOwner(t)
	observer, words, result, restoreObserver := legacy.PortTestXferSoundRawObserver()
	defer restoreObserver()
	*result = 0

	reset := func(seed uint32) {
		c.resetCase(effects, pix, seed, 127)
		if c.Objs.FirstList5() != nil {
			t.Fatal("update fixture reset left list-5 drawables behind")
		}
		env.Reset(0)
		noxflags.ResetGame()
		*words = [3]uint32{}
		*result = 0
	}
	newDrawable := func() *client.Drawable {
		dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(320, 448))
		if dr == nil {
			t.Fatal("could not create update probe drawable")
		}
		return dr
	}
	installObserver := func(dr *client.Drawable) {
		dr.Field_115 = observer
		c.Objs.List5Add(dr)
	}
	checkObserver := func(want uint32, dr *client.Drawable) {
		t.Helper()
		if words[0] != want {
			t.Fatalf("observer callback count = %d, want %d", words[0], want)
		}
		if want != 0 && (words[1] != uint32(uintptr(c.Viewport().C())) || words[2] != uint32(uintptr(dr.C()))) {
			t.Fatalf("observer callback args = %#v, want viewport/drawable %#x/%#x", words, uintptr(c.Viewport().C()), uintptr(dr.C()))
		}
	}

	// A nil primary update is the root loop's unconditional secondary-call path.
	reset(101)
	dr := newDrawable()
	dr.ClientUpdateFuncPtr = nil
	installObserver(dr)
	c.sub_49BD70(c.Viewport())
	checkObserver(1, dr)

	// The same raw observer can be used as the primary int-return callback and
	// the secondary void callback. This checks the exact nonzero test with
	// positive, negative, and high-bit int32 results, without inventing an
	// adapter that normalizes the callback return.
	for i, returnWord := range []int32{0, 1, -1, -1 << 31, 1<<31 - 1} {
		reset(uint32(111 + i))
		dr = newDrawable()
		dr.ClientUpdateFuncPtr = observer
		installObserver(dr)
		*result = returnWord
		c.sub_49BD70(c.Viewport())
		wantCount := uint32(2)
		if returnWord == 0 {
			wantCount = 1
		}
		checkObserver(wantCount, dr)
	}

	// A nil secondary function suppresses only the second call; the primary
	// update still runs and preserves its mutation/return behavior.
	reset(121)
	dr = newDrawable()
	dr.ClientUpdateFuncPtr = legacy.PortTestClientUpdateCloudCallback()
	dr.ZVal = 10
	*(*byte)(unsafe.Add(dr.C(), 432)) = 4
	c.Objs.List5Add(dr)
	c.sub_49BD70(c.Viewport())
	if dr.ZVal != 14 {
		t.Fatalf("cloud update without secondary changed ZVal to %d", dr.ZVal)
	}
	checkObserver(0, dr)

	// The fixture exposes the same identity used by the update table. Check its
	// int-dispatch result directly on one drawable; the owner mutates/spawns, so
	// use a fresh reset before checking the real root-loop gate.
	magic := legacy.PortTestClientUpdateMagicCallback()
	if magic == nil {
		t.Fatal("magic update callback is absent")
	}
	reset(103)
	dr = newDrawable()
	dr.ClientUpdateFuncPtr = magic
	*(*uint32)(unsafe.Add(dr.C(), 432)) = uint32(dr.PosVec.X)
	*(*uint32)(unsafe.Add(dr.C(), 436)) = uint32(dr.PosVec.Y)
	installObserver(dr)
	got := int(client.CallDrawableUpdateResult(magic, c.Viewport(), dr))
	if got != 0 {
		t.Fatalf("magic callback int result = %d, want 0 for qualified binary", got)
	}
	if len(c.Calls) != 5 {
		t.Fatalf("direct magic callback attempted %d trail sparks, want 4", len(c.Calls)-1)
	}

	reset(103)
	dr = newDrawable()
	dr.ClientUpdateFuncPtr = magic
	*(*uint32)(unsafe.Add(dr.C(), 432)) = uint32(dr.PosVec.X)
	*(*uint32)(unsafe.Add(dr.C(), 436)) = uint32(dr.PosVec.Y)
	installObserver(dr)
	c.sub_49BD70(c.Viewport())
	checkObserver(0, dr)
	if len(c.Calls) != 5 {
		t.Fatalf("magic callback attempted %d trail sparks, want 4", len(c.Calls)-1)
	}

	for _, call := range c.Calls[1:] {
		if call.Ref == 0 {
			t.Fatal("magic callback spark creation failed")
		}
	}

	// A retained int-return cloud updater is the nonzero positive control. It
	// wraps ZVal at 16 bits, returns 1, and must enter the same Field_115 call.
	reset(107)
	dr = newDrawable()
	dr.ClientUpdateFuncPtr = legacy.PortTestClientUpdateCloudCallback()
	dr.ZVal = 0xfffe
	*(*byte)(unsafe.Add(dr.C(), 432)) = 3
	installObserver(dr)
	c.sub_49BD70(c.Viewport())
	if dr.ZVal != 1 {
		t.Fatalf("cloud positive control ZVal = %d, want 1 after 16-bit wrap", dr.ZVal)
	}
	checkObserver(1, dr)

	// Pause exits before either primary or secondary callback.
	reset(109)
	dr = newDrawable()
	dr.ClientUpdateFuncPtr = legacy.PortTestClientUpdateCloudCallback()
	dr.ZVal = 10
	*(*byte)(unsafe.Add(dr.C(), 432)) = 4
	installObserver(dr)
	noxflags.SetGame(noxflags.GamePause)
	c.sub_49BD70(c.Viewport())
	if dr.ZVal != 10 {
		t.Fatalf("paused cloud updater changed ZVal to %d", dr.ZVal)
	}
	checkObserver(0, dr)
}
