//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func TestOptionsViewportRefresh(t *testing.T) {
	type record struct {
		Menu                         bool
		Cut, Return, Selected, After int
		Hidden, OverlayHidden        bool
	}
	var rows []record
	for _, menu := range []bool{false, true} {
		t.Run(fmt.Sprint(menu), func(t *testing.T) {
			o := newOptionsAudioOwner(t, menu)
			prepareOptionsConstruction(t, o, false)
			noxflags.ResetGame()
			cached := memmap.BlobByAddr(0x587000).Data[172884:172888]
			oldCached := append([]byte(nil), cached...)
			t.Cleanup(func() { copy(cached, oldCached) })
			if legacy.PortTestOptionsConstruct(menu) != 1 {
				t.Fatal("construct")
			}
			off := 1309820
			if menu {
				off = 1309720
			}
			root := (*gui.Window)(unsafe.Pointer(uintptr(*o.words[off])))
			for _, cut := range []int{-2147483648, -1, 0, 65, 69, 70, 75, 79, 80, 85, 89, 90, 100, 2147483647} {
				for id := 311; id <= 314; id++ {
					root.ChildByID(uint(id)).DrawData().Field0 &^= 4
				}
				root.SetHidden(true)
				nox_video_cutSize = cut
				op := 1
				if menu {
					op = 0
				}
				ret := legacy.PortTestOptionsAction(op, 0)
				selected := 0
				for id := 311; id <= 314; id++ {
					if root.ChildByID(uint(id)).DrawData().Field0&4 != 0 {
						if selected != 0 {
							t.Fatal("multiple selected viewport controls")
						}
						selected = id
					}
				}
				want := 311
				switch {
				case cut > 89:
					want = 314
				case cut > 79:
					want = 313
				case cut > 69:
					want = 312
				}
				r := record{Menu: menu, Cut: cut, Return: ret, Selected: selected, After: nox_video_cutSize, Hidden: root.Flags.IsHidden()}
				if selected != want {
					t.Fatalf("viewport selection %+v want=%d", r, want)
				}
				if menu {
					expected := []int{65, 75, 85, 100}[want-311]
					if r.After != expected || !r.Hidden || *(*int32)(unsafe.Pointer(&cached[0])) != int32(cut) {
						t.Fatalf("menu refresh %+v", r)
					}
				} else {
					overlay := (*gui.Window)(unsafe.Pointer(uintptr(*o.words[1309824])))
					r.OverlayHidden = overlay.Flags.IsHidden()
					if r.After != cut || r.Hidden || r.OverlayHidden {
						t.Fatalf("in-game show %+v", r)
					}
				}
				rows = append(rows, r)
			}
			if !menu {
				ret := legacy.PortTestOptionsAction(4, 0)
				if *o.words[1309820] != 0 || root.Flags&gui.StatusDestroyed == 0 {
					t.Fatalf("destroy return=%d root=%x flags=%x", ret, *o.words[1309820], root.Flags)
				}
			}
		})
	}
	spellbookCapture(t, "options-viewport", rows, "eb2016a7036e29206b99323c2cfa5e8353606308b0d340b33561611ea6a54024")
}
