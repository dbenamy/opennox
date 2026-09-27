//go:build porttest

package opennox

import (
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestFinalMenuCallbackArguments(t *testing.T) {
	old := legacy.Sub_4A18E0
	t.Cleanup(func() { legacy.Sub_4A18E0 = old })
	win, free := alloc.New(gui.Window{})
	t.Cleanup(free)
	callback := legacy.PortTestFinalMenuCallback()
	for _, w := range []*gui.Window{nil, win} {
		for _, result := range []int{-1 << 31, -1, 0, 1, 1<<31 - 1} {
			calls := 0
			legacy.Sub_4A18E0 = func(got *gui.Window, code, a, b int) int {
				calls++
				if got != w || uint32(code) != 0x80000007 || uint32(a) != 0xfedcba98 || uint32(b) != 0x87654321 {
					t.Fatalf("callback arguments: %p %08x %08x %08x", got, uint32(code), uint32(a), uint32(b))
				}
				return result
			}
			resp := callback(w, &gui.RawEvent{Event: -2147483641, Arg1: 0xfedcba98, Arg2: 0x87654321})
			if calls != 1 || gui.EventRespInt(resp) != result || (resp == nil) != (result == 0) {
				t.Fatalf("result %d: calls=%d response=%v", result, calls, resp)
			}
		}
	}
}

type finalMenuEvent struct {
	order *[]string
	after func()
}

func (e finalMenuEvent) EventCode() int { *e.order = append(*e.order, "code"); return 123 }
func (e finalMenuEvent) EventArgsC() (uintptr, uintptr) {
	*e.order = append(*e.order, "args")
	e.after()
	return 7, 9
}

func TestFinalMenuCallbackHookOrder(t *testing.T) {
	old := legacy.Sub_4A18E0
	t.Cleanup(func() { legacy.Sub_4A18E0 = old })
	legacy.Sub_4A18E0 = func(*gui.Window, int, int, int) int { t.Fatal("stale hook invoked"); return 0 }
	callback := legacy.PortTestFinalMenuCallback()
	var order []string
	event := finalMenuEvent{&order, func() {
		legacy.Sub_4A18E0 = func(w *gui.Window, code, a, b int) int {
			order = append(order, "hook")
			if w != nil || code != 123 || a != 7 || b != 9 {
				t.Fatal("event values changed")
			}
			return -1
		}
	}}
	if got := gui.EventRespInt(callback(nil, event)); got != -1 {
		t.Fatal(got)
	}
	if len(order) != 3 || order[0] != "code" || order[1] != "args" || order[2] != "hook" {
		t.Fatal(order)
	}
}

func TestFinalOptionsTooltipCallbacks(t *testing.T) {
	catalog := legacy.PortTestMapCatalogOpen(3)
	t.Cleanup(catalog.Close)
	o := newServerOptionsOwner(t)
	defer noxflags.PortTestGameFlags(0)()
	o.installConstructor(t)
	if got := o.call("construct", 0, ""); got != 1 {
		t.Fatal("constructor", got)
	}
	cursor := unsafe.Slice(memmap.PtrUint16(0x5D4594, 1096676), 256)
	old := append([]uint16(nil), cursor...)
	t.Cleanup(func() { copy(cursor, old) })
	for _, tc := range []struct {
		id   uint
		name string
	}{{10331, "assign"}, {10333, "damage"}} {
		w := o.options.ChildByID(tc.id)
		if w.TooltipFuncPtr == nil {
			t.Fatal("missing tooltip callback")
		}
		for flags := 0; flags < 256; flags++ {
			w.DrawData().Field0 = uint32(0xabcdef00) | uint32(flags)
			w.TooltipFunc(0xfedcba98)
			want := tc.name + " off"
			if flags&4 != 0 {
				want = tc.name + " on"
			}
			if got := alloc.GoString16(&cursor[0]); got != want {
				t.Fatalf("%s/%x: %q want %q", tc.name, flags, got, want)
			}
		}
	}
}

func TestFinalConversationTooltipNoop(t *testing.T) {
	win, free := alloc.New(gui.Window{})
	t.Cleanup(free)
	cursor := unsafe.Slice(memmap.PtrUint16(0x5D4594, 1096676), 256)
	old := append([]uint16(nil), cursor...)
	t.Cleanup(func() { copy(cursor, old) })
	clear(cursor)
	copy(cursor, []uint16{'k', 'e', 'e', 'p', 0})
	win.SetTooltipFunc(legacy.PortTestFinalConversationTooltip())
	for _, arg := range []uintptr{0, 1, 0xffffffff} {
		win.TooltipFunc(arg)
		if got := alloc.GoString16(&cursor[0]); got != "keep" {
			t.Fatal(got)
		}
	}
}

func TestFinalParticleCallbackTraversal(t *testing.T) {
	o := newObjectDrawingOwner(t)
	o.c.GUI = gui.New(o.c.Render())
	env := legacy.PortTestNewScreenEnvironment()
	t.Cleanup(env.Restore)
	snapshot, free := legacy.PortTestEffectsScreenParticles(4)
	t.Cleanup(free)
	args := [10]int32{0, 48, 48, 0, 0, 0, 2, 0, 0, 0}
	first, second, third := legacy.PortTestScreenParticleCreate(args), legacy.PortTestScreenParticleCreate(args), legacy.PortTestScreenParticleCreate(args)
	if first == nil || second == nil || third == nil {
		t.Fatal("particle creation")
	}
	original := first.Draw_fnc
	var calls []*legacy.Nox_screenParticle
	var created *legacy.Nox_screenParticle
	key, restore := legacy.PortTestObserveFinalPtr2(func(vp, raw unsafe.Pointer) int32 {
		if vp != o.c.Viewport().C() {
			t.Fatal("viewport argument")
		}
		p := (*legacy.Nox_screenParticle)(raw)
		calls = append(calls, p)
		if p == third {
			p.Field_44 = first // Traversal must use the next pointer saved before the call.
			args[3] = 1
			created = legacy.PortTestScreenParticleCreate(args)
		}
		return -1 << 31 // Draw traversal ignores the callback result.
	})
	t.Cleanup(restore)
	for _, p := range []*legacy.Nox_screenParticle{first, second, third} {
		p.Draw_fnc = key
	}
	legacy.PortTestScreenParticlesDraw(o.c.Viewport())
	third.Field_44 = second
	for _, p := range []*legacy.Nox_screenParticle{first, second, third} {
		p.Draw_fnc = original
	}
	if len(calls) != 3 || calls[0] != third || calls[1] != second || calls[2] != first {
		t.Fatal("callback traversal", calls)
	}
	if created == nil || created.Field_24 != 48<<16 || len(snapshot()) != 4 {
		t.Fatal("new head was visited or lost")
	}
}

func TestFinalPlayerSectionForeignResults(t *testing.T) {
	for _, result := range []int32{-1 << 31, -1, 0, 1, 1<<31 - 1} {
		calls := 0
		key, restore := legacy.PortTestObserveFinalPtr(func(arg unsafe.Pointer) int32 {
			calls++
			if arg != nil {
				t.Fatal("section argument must be nil")
			}
			return result
		})
		got := legacy.PortTestFinalPlayerSectionCall(key)
		restore()
		if int32(got) != result || calls != 1 {
			t.Fatalf("%d: %d calls=%d", result, got, calls)
		}
	}
}

func TestFinalMapCallbackGuardsAndResults(t *testing.T) {
	arg, free := alloc.New(uint32(123))
	t.Cleanup(free)
	calls := 0
	result := int32(1)
	key, restore := legacy.PortTestObserveFinalPtr(func(got unsafe.Pointer) int32 {
		calls++
		if got != unsafe.Pointer(arg) {
			t.Fatal("map callback argument")
		}
		return result
	})
	t.Cleanup(restore)
	if err := portTestFinalMapSection(nil, "absent-final-callback-section", nil); err == nil {
		t.Fatal("nil callback guard lost")
	}
	if err := portTestFinalMapSection(nil, "absent-final-callback-section", key); err != nil || calls != 0 {
		t.Fatal("unknown section must skip callback", err, calls)
	}
	for _, value := range []int32{-1 << 31, -1, 0, 1, 1<<31 - 1} {
		result = value
		before := calls
		err := portTestFinalMapSection(unsafe.Pointer(arg), "ObjectData", key)
		if (err != nil) != (value == 0) || calls != before+1 {
			t.Fatalf("%d: err=%v calls=%d", value, err, calls)
		}
	}
}

type finalFlameServer struct {
	legacy.Server
	deleted []*server.Object
}

func (s *finalFlameServer) DelayedDelete(u *server.Object) { s.deleted = append(s.deleted, u) }

func TestFinalFlameRegisteredCallback(t *testing.T) {
	newEffectsFullOwner(t)
	old := legacy.GetServer
	proxy := &finalFlameServer{Server: old()}
	legacy.GetServer = func() legacy.Server { return proxy }
	t.Cleanup(func() { legacy.GetServer = old })
	u, free := alloc.New(server.Object{})
	t.Cleanup(free)
	u.Field34 = proxy.S().Frame() // Expiry bypasses geometry; existing owner tests cover it.
	u.CallUpdate()
	if len(proxy.deleted) != 0 {
		t.Fatal("nil callback invoked")
	}
	u.Update = legacy.PortTestFinalFlameCallback()
	if u.Update == nil {
		t.Fatal("missing flame callback")
	}
	u.CallUpdate()
	u.CallUpdate()
	if len(proxy.deleted) != 2 || proxy.deleted[0] != u || proxy.deleted[1] != u {
		t.Fatal("flame owner dispatch", proxy.deleted)
	}
}
