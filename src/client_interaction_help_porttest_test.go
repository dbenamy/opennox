//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"os"
	"testing"
	"unsafe"

	"github.com/opennox/libs/client/keybind"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type interactionFont struct {
	font.Face
	height int
}

func (f interactionFont) Metrics() font.Metrics {
	m := f.Face.Metrics()
	m.CapHeight = fixed.I(f.height)
	return m
}

func TestClientInteractionHelpDialog(t *testing.T) {
	o := sessionDialogResources(t)
	words, restore := legacy.PortTestClientInteractionWords()
	defer restore()
	options, restoreOptions := legacy.PortTestServerOptionsWords()
	defer restoreOptions()
	other := o.c.GUI.NewWindowRaw(nil, 8, 0, 0, 10, 10, nil)
	defer other.Destroy()
	*options["root"] = uint32(uintptr(other.C()))
	dim := legacy.PortTestBindingDimensions()
	oldw, oldh := *dim[0], *dim[1]
	defer func() { *dim[0], *dim[1] = oldw, oldh }()
	*dim[0], *dim[1] = 801, 601
	oldHelp := legacy.Get_nox_server_sanctuaryHelp_54276()
	defer legacy.Set_nox_server_sanctuaryHelp_54276(oldHelp)
	table := serverConfigOwnBytes(t, 0x587000, 164512, 32)
	serverConfigOwnBytes(t, 0x5D4594, 1304400, 1280)
	quest := serverConfigOwnBytes(t, 0x5D4594, 1556160, 4)
	blob, err := os.ReadFile("common/memmap/nox/blobdata/blob_587000.dat")
	if err != nil {
		t.Fatal(err)
	}
	names := make([]string, 8)
	for i := range names {
		off := int(binary.LittleEndian.Uint32(blob[164512+4*i:]) - 0x587000)
		end := bytes.IndexByte(blob[off:], 0)
		if end < 0 {
			t.Fatal("help filename terminator")
		}
		names[i] = string(blob[off : off+end])
		ptr, free := alloc.CString(names[i])
		defer free()
		binary.LittleEndian.PutUint32(table[4*i:], uint32(uintptr(unsafe.Pointer(ptr))))
	}
	configure, resetStrings := o.c.srv.Server.PortTestMeterStrings(
		strman.Entry{ID: "Sanchlp.wnd:Help", Vals: []strman.Variant{{Str: "Help %s"}}},
		strman.Entry{ID: "Sanchlp.wnd:ClientHelp", Vals: []strman.Variant{{Str: "Client %s"}}},
		strman.Entry{ID: "cdecode.c:KeyToChat", Vals: []strman.Variant{{Str: "Chat %s"}}},
	)
	defer resetStrings()
	oldStrings := strMan
	defer func() { strMan = oldStrings }()
	o.c.ctrl = new(CtrlEventHandler)
	o.c.ctrl.addBinding(&CtrlEventBinding{keys: []keybind.Key{keybind.KeyQ}, events: []keybind.Event{45}})
	o.c.ctrl.addBinding(&CtrlEventBinding{keys: []keybind.Key{keybind.KeyF12}, events: []keybind.Event{8}})
	var sounds [][2]int
	defer legacy.PortTestClientSoundObserver(func(id, v int) { sounds = append(sounds, [2]int{id, v}) })()
	load := legacy.Nox_new_window_from_file
	defer func() { legacy.Nox_new_window_from_file = load }()
	loaded := ""
	legacy.Nox_new_window_from_file = func(name string, fn gui.WindowFunc) *gui.Window { loaded = name; return load(name, fn) }
	type row struct {
		Language, Height int
		Flags            uint32
		Quest            bool
		Resource, Text   string
		Closed           bool
	}
	var rows []row
	for _, height := range []int{10, 11} {
		resetDefault := o.c.Render().GetFonts().PortTestDefaultFont(interactionFont{basicfont.Face7x13, height})
		_, resetFont := o.c.Render().GetFonts().PortTestWindowFont(interactionFont{basicfont.Face7x13, height}, "default", "large")
		for language := 0; language < 8; language++ {
			configure(language)
			strMan = o.c.srv.Strings()
			for _, flags := range []uint32{0, 1, 4096, 4097} {
				for _, isQuest := range []bool{false, true} {
					resetFlags := noxflags.PortTestGameFlags(noxflags.GameFlag(flags))
					clear(quest)
					if isQuest {
						binary.LittleEndian.PutUint32(quest, 1)
					}
					other.Show()
					interactionCall("nox_xxx_cliShowHelpGui_49C560")
					idx := language
					if height > 10 {
						idx = 2
					}
					if loaded != names[idx] {
						t.Fatal("help resource", language, height, loaded, names[idx])
					}
					closed := flags&4096 != 0 || isQuest
					if (interactionCall("sub_49C810") == 0) != closed {
						t.Fatal("help auto close", flags, isQuest)
					}
					text := ""
					if !closed {
						w := (*gui.Window)(unsafe.Pointer(uintptr(*words["dword_5d4594_1305680"])))
						text = interactionStaticText(w.ChildByID(4102))
						want := "Client Q Chat F12"
						if flags&1 != 0 {
							want = "Help Q Chat F12"
						}
						if text != want || w.Off != image.Pt((801-w.SizeVal.X)/2, (601-w.SizeVal.Y)/2) || other.GetFlags().IsHidden() != (flags&1 != 0) {
							t.Fatal("help text/layout", text, want, w.Off)
						}
						for _, code := range []uintptr{0, 22, 23} {
							if interactionCall("nox_xxx_wnd_49C760", uintptr(w.C()), code, 0xffffffff, 0) != 1 {
								t.Fatal("help numeric notification")
							}
						}
						checked := (language+height)%2 == 0
						w.ChildByID(4104).DrawData().Field0 &^= 4
						if checked {
							w.ChildByID(4104).DrawData().Field0 |= 4
						}
						interactionCall("nox_xxx_wnd_49C760", uintptr(w.C()), 16391, uintptr(w.ChildByID(4103).C()), 0)
						wantHelp := 1
						if checked {
							wantHelp = 0
						}
						if legacy.Get_nox_server_sanctuaryHelp_54276() != wantHelp {
							t.Fatal("help checkbox persistence")
						}
					}
					if *words["dword_5d4594_1305680"] != 0 || other.GetFlags().IsHidden() {
						t.Fatal("help owners/options restoration")
					}
					rows = append(rows, row{language, height, flags, isQuest, loaded, text, closed})
					resetFlags()
					o.c.GUI.FreeDestroyed()
				}
			}
		}
		resetFont()
		resetDefault()
	}
	interactionCapture(t, "help", rows)
}
