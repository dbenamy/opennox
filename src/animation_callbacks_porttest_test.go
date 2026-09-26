//go:build porttest

package opennox

import (
	"image"
	"reflect"
	"testing"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
)

func TestAnimationCallbackReturnAndHook(t *testing.T) {
	old := legacy.WinMainMenuAnimOutStartFnc
	t.Cleanup(func() { legacy.WinMainMenuAnimOutStartFnc = old })
	key := legacy.Get_winMainMenuAnimOutStartFnc()
	if key == nil || key != legacy.Get_winMainMenuAnimOutStartFnc() {
		t.Fatal("unstable animation callback identity")
	}
	a := gui.Anim{Func12Ptr: key, Func13Ptr: key}
	for _, value := range []int{0, 1, -1, 0x7fffffff, -2147483648} {
		calls := 0
		legacy.WinMainMenuAnimOutStartFnc = func() int { calls++; return value }
		if got := a.Func12(); got != value || calls != 1 {
			t.Fatalf("Func12: got %d calls %d, want %d/1", got, calls, value)
		}
		if got := a.Func13(); got != value || calls != 2 {
			t.Fatalf("Func13: got %d calls %d, want %d/2", got, calls, value)
		}
	}
	// An already stored identity must observe a subsequently replaced hook.
	legacy.WinMainMenuAnimOutStartFnc = func() int { return 73 }
	if got := a.Func13(); got != 73 {
		t.Fatalf("callback retained an earlier hook: %d", got)
	}
}

func TestAnimationForeignCallbackBoundary(t *testing.T) {
	key, count, restore := legacy.PortTestOptionsDone()
	t.Cleanup(restore)
	a := gui.Anim{Func12Ptr: key, Func13Ptr: key}
	if got := a.Func12(); got != 1 || count() != 1 {
		t.Fatalf("foreign Func12: result %d calls %d", got, count())
	}
	if got := a.Func13(); got != 1 || count() != 2 {
		t.Fatalf("foreign Func13: result %d calls %d", got, count())
	}
}

func TestAnimationCompletionOrdering(t *testing.T) {
	o := newEntryOwner(t)
	t.Cleanup(gui.PortTestOwnAnimations())
	old := legacy.WinMainMenuAnimOutDoneFnc
	t.Cleanup(func() { legacy.WinMainMenuAnimOutDoneFnc = old })
	w := o.c.GUI.NewWindowRaw(o.parent, 8, 0, 0, 12, 12, nil)
	start, end := image.Pt(5, 9), image.Pt(-2, -3)
	a := gui.NewAnim(w, start, end, image.Pt(3, 4), image.Pt(-4, -5))
	a.StateID = 0x6f001
	w.SetPos(start)
	a.SetState(gui.AnimOut)
	gui.SetAnimGlobalState(gui.AnimOut)
	a.FncDoneOutPtr = legacy.Get_winMainMenuAnimOutDoneFnc()
	calls := 0
	completionReady := false
	legacy.WinMainMenuAnimOutDoneFnc = func() int {
		calls++
		completionReady = a.State() == gui.AnimOutDone && gui.AnimGlobalState() == gui.AnimOutDone && w.Offs() == end
		a.Free()
		return -1 // Completion discards the result without changing its call convention.
	}
	for i, want := range []image.Point{image.Pt(1, 4), image.Pt(-2, -1), end} {
		gui.AnimTick()
		if got := w.Offs(); got != want {
			t.Fatalf("tick %d position %v, want %v", i, got, want)
		}
		wantCalls := 0
		if i == 2 {
			wantCalls = 1
		}
		if calls != wantCalls {
			t.Fatalf("tick %d callbacks %d, want %d", i, calls, wantCalls)
		}
	}
	if !completionReady {
		t.Fatal("completion callback preceded final position/state")
	}
	if gui.FindAnimForStateID(0x6f001) != nil {
		t.Fatal("completion did not unlink its animation")
	}
	gui.AnimTick()
	if calls != 1 {
		t.Fatal("completed callback ran again")
	}
}

func TestAnimationCompletionListIteration(t *testing.T) {
	o := newEntryOwner(t)
	t.Cleanup(gui.PortTestOwnAnimations())
	oldStart, oldDone := legacy.WinMainMenuAnimOutStartFnc, legacy.WinMainMenuAnimOutDoneFnc
	t.Cleanup(func() { legacy.WinMainMenuAnimOutStartFnc, legacy.WinMainMenuAnimOutDoneFnc = oldStart, oldDone })
	makeAnim := func(id gui.StateID) *gui.Anim {
		w := o.c.GUI.NewWindowRaw(o.parent, 8, 0, 0, 12, 12, nil)
		a := gui.NewAnim(w, image.Point{}, image.Pt(2, 3), image.Pt(-4, -4), image.Pt(4, 4))
		w.SetPos(image.Point{})
		a.StateID = id
		a.SetState(gui.AnimOut)
		return a
	}
	a, b := makeAnim(0x6f002), makeAnim(0x6f003)
	a.FncDoneOutPtr = legacy.Get_winMainMenuAnimOutStartFnc()
	b.FncDoneOutPtr = legacy.Get_winMainMenuAnimOutDoneFnc()
	var order []string
	legacy.WinMainMenuAnimOutStartFnc = func() int { order = append(order, "a"); a.Free(); return 1 }
	legacy.WinMainMenuAnimOutDoneFnc = func() int { order = append(order, "b"); b.Free(); return 1 }
	gui.AnimTick()
	if !reflect.DeepEqual(order, []string{"b", "a"}) {
		t.Fatalf("callback/free iteration order %v", order)
	}
	if gui.FindAnimForStateID(0x6f002) != nil || gui.FindAnimForStateID(0x6f003) != nil {
		t.Fatal("freed animations remain linked")
	}
}

func TestMainMenuAnimationCompletionAfterFree(t *testing.T) {
	o := newEntryOwner(t)
	t.Cleanup(gui.PortTestOwnAnimations())
	oldWindow, oldTop, oldBottom := winMainMenu, winMainMenuAnimTop, winMainMenuAnimBottom
	oldHook := legacy.WinMainMenuAnimOutStartFnc
	t.Cleanup(func() {
		winMainMenu, winMainMenuAnimTop, winMainMenuAnimBottom = oldWindow, oldTop, oldBottom
		legacy.WinMainMenuAnimOutStartFnc = oldHook
	})
	window := o.c.GUI.NewWindowRaw(nil, 8, 0, 0, 100, 100, nil)
	top := o.c.GUI.NewWindowRaw(window, 8, 0, 0, 50, 20, nil)
	bottom := o.c.GUI.NewWindowRaw(window, 8, 0, 40, 50, 20, nil)
	winMainMenu = window
	winMainMenuAnimTop = gui.NewAnim(top, image.Point{}, image.Pt(0, -20), image.Pt(0, 2), image.Pt(0, -2))
	winMainMenuAnimBottom = gui.NewAnim(bottom, image.Pt(0, 40), image.Pt(0, 60), image.Pt(0, -2), image.Pt(0, 2))
	winMainMenuAnimTop.StateID, winMainMenuAnimBottom.StateID = 0x6f004, 0x6f005
	winMainMenuAnimTop.Func13Ptr = legacy.Get_winMainMenuAnimOutStartFnc()
	calls, cleanedBeforeCallback := 0, false
	legacy.WinMainMenuAnimOutStartFnc = func() int {
		calls++
		cleanedBeforeCallback = winMainMenu == nil && winMainMenuAnimTop == nil && winMainMenuAnimBottom == nil &&
			window.Flags&gui.StatusDestroyed != 0 && top.Flags&gui.StatusDestroyed != 0 && bottom.Flags&gui.StatusDestroyed != 0 &&
			gui.FindAnimForStateID(0x6f004) == nil && gui.FindAnimForStateID(0x6f005) == nil
		return -7
	}
	if got := winMainMenuAnimOutDoneFnc(); got != 1 || calls != 1 || !cleanedBeforeCallback {
		t.Fatalf("main-menu completion result=%d calls=%d cleanup-before-callback=%v", got, calls, cleanedBeforeCallback)
	}
}
