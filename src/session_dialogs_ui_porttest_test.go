//go:build porttest

package opennox

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"golang.org/x/image/font/basicfont"
)

func sessionDialogResources(t *testing.T) *listboxOwner {
	t.Helper()
	o := newListboxOwner(t)
	serverConfigOwnBytes(t, 0x5D4594, 1309752, 4) // cached disconnect image handle
	_, restoreFont := o.c.Render().GetFonts().PortTestWindowFont(basicfont.Face7x13, "default", "large")
	t.Cleanup(restoreFont)
	oldStrings := strMan
	strMan = o.c.srv.Strings()
	t.Cleanup(func() { strMan = oldStrings })
	load := legacy.Nox_new_window_from_file
	images := legacy.Nox_xxx_gLoadImg
	t.Cleanup(func() { legacy.Nox_new_window_from_file = load; legacy.Nox_xxx_gLoadImg = images })
	legacy.Nox_xxx_gLoadImg = func(string) *noxrender.Image { return o.images[0] }
	legacy.Nox_new_window_from_file = func(name string, fn gui.WindowFunc) *gui.Window {
		b, err := os.ReadFile(filepath.Join(os.Getenv("OPENNOX_SPELLBOOK_ASSETS"), "window", name))
		if err != nil {
			t.Fatal(err)
		}
		w := newWindowFromString(o.c.GUI, string(b), fn)
		if w == nil {
			t.Fatal("resource constructor", name)
		}
		return w
	}
	return o
}
func sessionCall(op string, w *gui.Window, code int, a, b uint32) uintptr {
	var p unsafe.Pointer
	if w != nil {
		p = w.C()
	}
	return legacy.PortTestSessionDialogCall(op, p, code, a, b)
}
func sessionCapture(t *testing.T, name string, v any) {
	t.Helper()
	prefix := os.Getenv("OPENNOX_SESSION_DIALOGS_CAPTURE")
	if prefix == "" {
		return
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(prefix+"-"+name+".json", append(data, '\n'), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestSessionFilterResourceLifecycle(t *testing.T) {
	o := sessionDialogResources(t)
	words, restore := legacy.PortTestSessionDialogWords()
	defer restore()
	modes := serverConfigOwnBytes(t, 0x5D4594, 1193372, 8)
	settings := serverConfigOwnBytes(t, 0x5D4594, 1193388, 88)
	type row struct {
		Mode, Mask uint32
		Flags      map[uint]uint32
		Text       string
		Saved      [11]uint32
		Closed     bool
	}
	var rows []row
	for mode := uint32(0); mode < 3; mode++ {
		for mask := uint32(0); mask < 64; mask++ {
			clear(modes)
			clear(settings)
			binary.LittleEndian.PutUint32(modes, mode)
			fields := []int{0, 1, 2, 3, 5, 10}
			for bit, field := range fields {
				if mask&(1<<uint(bit)) != 0 {
					binary.LittleEndian.PutUint32(settings[4*field:], 2)
				}
			}
			ping := uint32(250 + mask)
			binary.LittleEndian.PutUint32(settings[16:], ping)
			p := sessionCall("filterOpen", o.parent, 0, 0, 0)
			if p == 0 {
				t.Fatal("filter constructor")
			}
			w := (*gui.Window)(unsafe.Pointer(p))
			controls := w.ChildByID(10012)
			if w.Parent() != o.parent || controls.Parent() != w || (*words["filter"] != uint32(p)) {
				t.Fatal("filter owners")
			}
			if controls.Flags.Has(gui.StatusHidden) != (mode != 2) {
				t.Fatal("control visibility", mode, mask)
			}
			if w.ChildByID(uint(10024+mode)).DrawData().Field0&4 == 0 {
				t.Fatal("filter mode selection", mode, mask)
			}
			for bit, id := range []uint{10028, 10029, 10030, 10015, 10014, 10018} {
				if (w.ChildByID(id).DrawData().Field0&4 != 0) != (mask&(1<<uint(bit)) != 0) {
					t.Fatal("checkbox state", mode, mask, id)
				}
			}
			if mode == 2 {
				if w.ChildByID(10031).Flags.Has(gui.StatusEnabled) != (mask&1 != 0) {
					t.Fatal("ping enabled", mask)
				}
				for _, id := range []uint{10016, 10017} {
					if w.ChildByID(id).Flags.Has(gui.StatusEnabled) != (mask&8 != 0) {
						t.Fatal("resolution enabled", mask)
					}
				}
			}
			data := (*gui.EntryFieldData)(w.ChildByID(10031).WidgetData)
			text := alloc.GoString16(&data.Text[0])
			if text != strconv.Itoa(int(ping)) {
				t.Fatal("ping text", text, ping)
			}
			r := row{Mode: mode, Mask: mask, Flags: map[uint]uint32{}, Text: text}
			for _, id := range []uint{10012, 10014, 10015, 10016, 10017, 10018, 10024, 10025, 10026, 10028, 10029, 10030, 10031} {
				r.Flags[id] = uint32(w.ChildByID(id).Flags)
			}
			// Resource-parser child notifications use numeric IDs, not window pointers.
			for _, id := range []uint32{0, 1, 575, 10012, 0xffffffff} {
				if sessionCall("filterEvent", w, 22, id, 0) != 0 {
					t.Fatal("numeric child notification")
				}
			}
			if sessionCall("filterEvent", w, 23, 0, 0) != 1 {
				t.Fatal("dialog query")
			}
			sessionCall("filterSave", nil, 0, 0, 0)
			for i := range r.Saved {
				r.Saved[i] = binary.LittleEndian.Uint32(settings[4*i:])
			}
			if mode == 2 {
				want := uint32(0)
				if mask&8 != 0 {
					want = 130
				}
				if r.Saved[3] != want {
					t.Fatal("resolution save", mask, r.Saved)
				}
			}
			sessionCall("filterClose", nil, 0, 0, 0)
			r.Closed = *words["filter"] == 0
			if !r.Closed {
				t.Fatal("filter root retained")
			}
			sessionCall("filterClose", nil, 0, 0, 0)
			o.c.GUI.FreeDestroyed()
			rows = append(rows, r)
		}
	}
	sessionCapture(t, "filter-resource", rows)
}

func TestSessionMOTDResourceLifecycle(t *testing.T) {
	o := sessionDialogResources(t)
	words, restore := legacy.PortTestSessionDialogWords()
	defer restore()
	type row struct {
		Round         int
		Thumb         image.Point
		Text          []string
		First, Second uintptr
		Flags         [2]uint32
	}
	var rows []row
	for round := 0; round < 3; round++ {
		p := sessionCall("motdOpen", nil, 0, 0, 0)
		if p == 0 {
			t.Fatal("MOTD constructor")
		}
		w := (*gui.Window)(unsafe.Pointer(p))
		list := w.ChildByID(4203)
		slider := w.ChildByID(4204)
		if list == nil || slider == nil || slider.Field100Ptr == nil {
			t.Fatal("MOTD child tree")
		}
		d := (*gui.ScrollListBoxData)(list.WidgetData)
		if d.Field_9 != slider.C() || d.Field_7 != w.ChildByID(4205).C() || d.Field_8 != w.ChildByID(4206).C() {
			t.Fatal("MOTD scrollbar links")
		}
		if slider.Field100Ptr.SizeVal != image.Pt(16, 10) || slider.DrawData().Window != list {
			t.Fatal("MOTD thumb/owner")
		}
		if sessionCall("motdShown", nil, 0, 0, 0) != 0 || sessionCall("motdClose", nil, 0, 0, 0) != 0 {
			t.Fatal("initial hidden state")
		}
		for _, text := range []string{"First", "", "Second"} {
			s, free := alloc.CString(text)
			legacy.PortTestSessionDialogCall("motdAdd", unsafe.Pointer(s), 0, 0, 0)
			free()
		}
		names := serverPanelsListNames(list)
		if fmt.Sprint(names) != "[First Second]" {
			t.Fatal("MOTD rows", names)
		}
		w.SetHidden(false)
		w.Flags |= gui.StatusEnabled
		list.Flags |= gui.StatusEnabled
		list.Focus()
		if sessionCall("motdShown", nil, 0, 0, 0) != 1 {
			t.Fatal("MOTD shown predicate")
		}
		first := sessionCall("motdClose", nil, 0, 0, 0)
		second := sessionCall("motdClose", nil, 0, 0, 0)
		if first != 1 || second != 0 || !w.Flags.Has(gui.StatusHidden) || w.Flags.Has(gui.StatusEnabled) || list.Flags.Has(gui.StatusEnabled) || o.c.GUI.Focused() != nil || len(serverPanelsListNames(list)) != 0 {
			t.Fatal("MOTD close state", first, second, w.Flags, list.Flags)
		}
		rows = append(rows, row{Round: round, Thumb: slider.Field100Ptr.SizeVal, Text: names, First: first, Second: second, Flags: [2]uint32{uint32(w.Flags), uint32(list.Flags)}})
		w.Destroy()
		o.c.GUI.FreeDestroyed()
		*words["motd"] = 0
		*words["motdList"] = 0
	}
	sessionCapture(t, "motd-resource", rows)
}

func TestSessionDisconnectResourceLifecycle(t *testing.T) {
	o := sessionDialogResources(t)
	words, restore := legacy.PortTestSessionDialogWords()
	defer restore()
	dim := legacy.PortTestBindingDimensions()
	old := [2]int32{*dim[0], *dim[1]}
	defer func() { *dim[0], *dim[1] = old[0], old[1] }()
	type row struct {
		Size, Dialog, Icon image.Point
		Input              []uintptr
		Hidden, Enabled    bool
	}
	var rows []row
	for _, size := range []image.Point{{640, 480}, {641, 481}, {1024, 768}, {1280, 960}} {
		*dim[0], *dim[1] = int32(size.X), int32(size.Y)
		if sessionCall("disconnectOpen", nil, 0, 0, 0) != 1 {
			t.Fatal("disconnect constructor")
		}
		w := (*gui.Window)(unsafe.Pointer(uintptr(*words["disconnect"])))
		icon := (*gui.Window)(unsafe.Pointer(uintptr(*words["disconnectIcon"])))
		if w == nil || icon == nil {
			t.Fatal("disconnect owners")
		}
		if w.Off != image.Pt(size.X/2-w.SizeVal.X/2, size.Y/2-w.SizeVal.Y/2) || icon.Off != image.Pt(size.X-50, size.Y/2+3) || icon.SizeVal != image.Pt(50, 50) {
			t.Fatal("disconnect geometry", size, w.Off, icon.Off)
		}
		for _, enabled := range []uint32{0, 1, 7, 0} {
			sessionCall("disconnectIcon", nil, 0, enabled, 0)
			if icon.Flags.Has(gui.StatusHidden) != (enabled == 0) {
				t.Fatal("icon visibility")
			}
		}
		var events []int
		w.SetFunc93(func(_ *gui.Window, e gui.WindowEvent) gui.WindowEventResp {
			events = append(events, e.EventCode())
			return nil
		})
		r := row{Size: size, Dialog: w.Off, Icon: icon.Off}
		for _, a := range []uint32{1, 57, 0, 0xffffffff} {
			r.Input = append(r.Input, sessionCall("disconnectInput", w, 21, a, 0))
		}
		if fmt.Sprint(r.Input) != "[1 0 0 0]" || fmt.Sprint(events) != "[5]" {
			t.Fatal("disconnect keys", r.Input, events)
		}
		for _, id := range []uint32{1, 575, 576, 577, 0xffffffff} {
			if sessionCall("disconnectEvent", w, 22, id, 0) != 0 {
				t.Fatal("numeric construction event")
			}
		}
		if sessionCall("disconnectEvent", w, 23, 0, 0) != 1 {
			t.Fatal("disconnect query")
		}
		sessionCall("disconnectShow", nil, 0, 1, 0)
		if w.Flags.Has(gui.StatusHidden) || !w.Flags.Has(gui.StatusEnabled) || o.c.GUI.Focused() != w {
			t.Fatal("disconnect shown/focus")
		}
		sessionCall("disconnectShow", nil, 0, 0, 0)
		r.Hidden = w.Flags.Has(gui.StatusHidden)
		r.Enabled = w.Flags.Has(gui.StatusEnabled)
		if !r.Hidden || r.Enabled || o.c.GUI.Focused() != nil {
			t.Fatal("disconnect hidden/focus")
		}
		if sessionCall("disconnectClose", nil, 0, 0, 0) != 0 || *words["disconnect"] != 0 || *words["disconnectIcon"] != 0 {
			t.Fatal("disconnect destroy")
		}
		o.c.GUI.FreeDestroyed()
		rows = append(rows, r)
	}
	sessionCapture(t, "disconnect-resource", rows)
}
