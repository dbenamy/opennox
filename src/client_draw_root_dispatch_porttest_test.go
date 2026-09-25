//go:build porttest

package opennox

import (
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// Exercise the registered addresses, including the original C forwarding path,
// independently of the real renderer contracts. Hooks must be read at call time.
func TestClientDrawRootCallbackDispatch(t *testing.T) {
	cases := []struct {
		name string
		key  unsafe.Pointer
		hook *func(*noxrender.Viewport, *client.Drawable) int
	}{
		{"debug", legacy.Get_nox_thing_debug_draw(), &legacy.Nox_thing_debug_draw},
		{"monster", legacy.Get_nox_thing_monster_draw(), &legacy.Nox_thing_monster_draw},
		{"vector", legacy.Get_nox_thing_vector_animate_draw(), &legacy.Nox_thing_vector_animate_draw},
		{"released-soul", legacy.Get_nox_thing_released_soul_draw(), &legacy.Nox_thing_vector_animate_draw},
		{"state", legacy.PortTestSpriteAnimationCallback(6), &legacy.Nox_thing_animate_state_draw},
		{"player", legacy.Get_nox_thing_player_draw(), &legacy.Nox_thing_player_draw},
		{"npc", legacy.Get_nox_thing_npc_draw(), &legacy.Nox_thing_npc_draw},
	}
	seen := make(map[unsafe.Pointer]string)
	for _, tc := range cases {
		if tc.key == nil || seen[tc.key] != "" {
			t.Fatalf("%s key is nil or aliases %s", tc.name, seen[tc.key])
		}
		seen[tc.key] = tc.name
	}
	if client.ThingDrawDefault != cases[0].key {
		t.Fatal("default object drawing does not use the debug callback")
	}
	vp, freeVP := alloc.New(noxrender.Viewport{})
	defer freeVP()
	dr, freeDR := alloc.New(client.Drawable{})
	defer freeDR()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			old := *tc.hook
			defer func() { *tc.hook = old }()
			dr.DrawFuncPtr = tc.key
			for _, want := range []int{-2147483648, -1, 0, 1, 2147483647} {
				calls := 0
				*tc.hook = func(gotVP *noxrender.Viewport, gotDR *client.Drawable) int {
					if gotVP != vp || gotDR != dr {
						t.Error("callback arguments changed")
					}
					calls++
					return want
				}
				if got := dr.CallDraw(vp); got != want || calls != 1 {
					t.Fatalf("drawable result=%d calls=%d want=%d", got, calls, want)
				}
				if got := legacy.CallDrawFunc(dr, vp); got != want || calls != 2 {
					t.Fatalf("legacy result=%d calls=%d want=%d", got, calls, want)
				}
				client.CallDrawableDrawDiscard(tc.key, vp, dr)
				if calls != 3 {
					t.Fatal("discard path did not invoke callback exactly once")
				}
				// A replacement closure must be observed even after prior calls.
				*tc.hook = func(gotVP *noxrender.Viewport, gotDR *client.Drawable) int {
					if gotVP != vp || gotDR != dr {
						t.Error("replacement callback arguments changed")
					}
					calls += 10
					return ^want
				}
				if got := dr.CallDraw(vp); got != ^want || calls != 13 {
					t.Fatalf("replacement result=%d calls=%d want=%d", got, calls, ^want)
				}
			}
			for _, args := range []struct {
				vp *noxrender.Viewport
				dr *client.Drawable
			}{{nil, dr}, {vp, nil}, {nil, nil}} {
				calls := 0
				*tc.hook = func(gotVP *noxrender.Viewport, gotDR *client.Drawable) int {
					if gotVP != args.vp || gotDR != args.dr {
						t.Error("nil argument changed")
					}
					calls++
					return -17
				}
				if got := client.CallDrawableDrawResult(tc.key, args.vp, args.dr); got != -17 || calls != 1 {
					t.Fatalf("nil argument result=%d calls=%d", got, calls)
				}
				client.CallDrawableDrawDiscard(tc.key, args.vp, args.dr)
				if calls != 2 {
					t.Fatal("nil argument discard call")
				}
			}
		})
	}
}
