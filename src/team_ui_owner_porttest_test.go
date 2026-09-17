//go:build porttest

package opennox

import (
	"fmt"
	"image"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"golang.org/x/image/font/basicfont"
)

type teamUIOwner struct {
	*entryOwner
	words             map[string]*uint32
	loads             []string
	missing           bool
	configureLanguage func(int)
}

func teamUIHUDResource(ctf bool) string {
	var b strings.Builder
	b.WriteString("FONT = small;\nWINDOW\n 8800 0 0 160 320 USER;\n STATUS = ENABLED;\n CHILD\n")
	n := 1
	if ctf {
		n = 16
	}
	for i := 0; i < n; i++ {
		fmt.Fprintf(&b, " WINDOW\n %d %d %d 14 10 USER;\n STATUS = ENABLED;\n END\n", 8811+i, 16+18*(i%4), 16+14*(i/4))
	}
	b.WriteString(" END\nEND\n")
	return b.String()
}
func newTeamUIOwner(t *testing.T) *teamUIOwner {
	t.Helper()
	o := &teamUIOwner{entryOwner: newEntryOwner(t)}
	oldNet := o.c.srv.NetList
	o.c.srv.NetList = netlist.New()
	o.c.srv.NetList.Init()
	t.Cleanup(func() { o.c.srv.NetList.Free(); o.c.srv.NetList = oldNet })
	configure, restore := o.c.srv.PortTestMeterStrings(
		strman.Entry{ID: "playrlst.c:Rename", Vals: []strman.Variant{{Str: "Rename team"}}},
		strman.Entry{ID: "playrlst.c:NewName", Vals: []strman.Variant{{Str: "New team name"}}},
		strman.Entry{ID: "playrlst.c:RedFlag", Vals: []strman.Variant{{Str: "\tred flag"}}},
		strman.Entry{ID: "playrlst.c:BlueFlag", Vals: []strman.Variant{{Str: "\tblue flag"}}},
		strman.Entry{ID: "playrlst.c:Title1", Vals: []strman.Variant{{Str: "Chat teams"}}},
		strman.Entry{ID: "GUI_CTF.c:FlagHomeTT", Vals: []strman.Variant{{Str: "home"}}},
		strman.Entry{ID: "GUI_CTF.c:YourFlagCarriedTT", Vals: []strman.Variant{{Str: "your flag carried"}}},
		strman.Entry{ID: "GUI_CTF.c:TheirFlagCarriedTT", Vals: []strman.Variant{{Str: "their flag carried"}}},
		strman.Entry{ID: "GUI_CTF.c:FlagAwayTT", Vals: []strman.Variant{{Str: "away"}}},
	)
	t.Cleanup(restore)
	configure(0)
	o.configureLanguage = configure
	oldStrings := strMan
	strMan = o.c.srv.Strings()
	t.Cleanup(func() { strMan = oldStrings })
	_, restoreFont := o.c.Render().GetFonts().PortTestWindowFont(basicfont.Face7x13, "small", "large")
	t.Cleanup(restoreFont)
	oldLoad := legacy.Nox_xxx_gLoadImg
	legacy.Nox_xxx_gLoadImg = func(name string) *noxrender.Image {
		if name == "FlagTeamBorder" {
			return o.images[4]
		}
		if name == "UISlider" {
			return o.images[5]
		}
		if name == "UISliderLit" {
			return o.images[6]
		}
		return oldLoad(name)
	}
	t.Cleanup(func() { legacy.Nox_xxx_gLoadImg = oldLoad })
	oldParser := legacy.Nox_new_window_from_file
	legacy.Nox_new_window_from_file = func(name string, fn gui.WindowFunc) *gui.Window {
		o.loads = append(o.loads, name)
		if o.missing {
			return newWindowFromString(o.c.GUI, "", fn)
		}
		var resource string
		switch name {
		case "GUI_CTF.wnd":
			resource = teamUIHUDResource(true)
		case "team-ui-player-0.wnd", "team-ui-player-1.wnd", "team-ui-player-2.wnd", "team-ui-player-3.wnd", "team-ui-player-4.wnd", "team-ui-player-5.wnd", "team-ui-player-6.wnd", "team-ui-player-7.wnd", "team-ui-player-8.wnd":
			resource = teamUIPlayersResource()
		case "gui_fb.wnd":
			resource = teamUIHUDResource(false)
		default:
			t.Fatalf("unowned team UI resource %q", name)
		}
		w := newWindowFromString(o.c.GUI, resource, fn)
		if w != nil {
			for i := 0; i < 16; i++ {
				if c := w.ChildByID(uint(8811 + i)); c != nil {
					d := c.DrawData()
					d.SetBackgroundImage(o.images[0])
					d.SetEnabledImage(o.images[1])
					d.SetDisabledImage(o.images[2])
				}
			}
		}
		return w
	}
	t.Cleanup(func() { legacy.Nox_new_window_from_file = oldParser })
	oldW, oldH, oldMax := nox_win_width_game, nox_win_height_game, noxVideoMax
	nox_win_width_game, nox_win_height_game = 640, 480
	noxVideoMax = image.Pt(800, 600)
	t.Cleanup(func() { nox_win_width_game, nox_win_height_game, noxVideoMax = oldW, oldH, oldMax })
	installListboxPalette(t)
	for _, region := range [][3]int{{0x587000, 129048, 36}, {0x587000, 128968, 80}, {0x5D4594, 371380, 58}, {0x5D4594, 371688, 4}} {
		b := memmap.BlobByAddr(uintptr(region[0])).Data[region[1] : region[1]+region[2]]
		old := append([]byte(nil), b...)
		t.Cleanup(func() { copy(b, old) })
		clear(b)
	}
	for i := 0; i < 9; i++ {
		*memmap.PtrPtr(0x587000, 129048+uintptr(4*i)) = unsafe.Pointer(alloc.InternCString(fmt.Sprintf("team-ui-player-%d.wnd", i)))
	}
	for i := 0; i < 10; i++ {
		*memmap.PtrUint8(0x587000, 128968+uintptr(8*i)) = uint8(i + 1)
	}
	ballType := legacy.PortTestMapDrawableTeamWord()
	oldBallType := *ballType
	*ballType = 0
	t.Cleanup(func() { *ballType = oldBallType })
	var done func()
	o.words, done = legacy.PortTestTeamUIWords()
	t.Cleanup(done)
	return o
}
func (o *teamUIOwner) window(name string) *gui.Window {
	return (*gui.Window)(unsafe.Pointer(uintptr(*o.words[name])))
}
func teamUICall(op string, code, value int) int {
	return legacy.PortTestTeamUI(op, nil, nil, code, value, "")
}

func teamUIPlayersResource() string {
	var b strings.Builder
	b.WriteString("FONT = small;\nWINDOW\n 10500 0 0 500 400 USER;\n STATUS = ENABLED;\n CHILD\n")
	for i := 0; i < 2; i++ {
		fmt.Fprintf(&b, " WINDOW\n %d %d 20 180 240 SCROLLLISTBOX;\n STATUS = ENABLED;\n DATA = 128 1 0 0 1 %d 0;\n END\n", 10501+i, 10+210*i, 1-i)
	}
	for _, id := range []int{10503, 10507, 10509, 10515, 10516, 10518, 10519} {
		fmt.Fprintf(&b, " WINDOW\n %d 0 280 40 20 PUSHBUTTON;\n STATUS = ENABLED;\n END\n", id)
	}
	for _, id := range []int{10517, 10520} {
		fmt.Fprintf(&b, " WINDOW\n %d 190 20 16 240 VERTSLIDER;\n STATUS = ENABLED;\n DATA = 0 100;\n END\n", id)
	}
	b.WriteString(" WINDOW\n 10504 0 0 180 20 STATICTEXT;\n STATUS = ENABLED;\n DATA = 0 0 playrlst.c:Title1;\n END\n END\nEND\n")
	return b.String()
}
