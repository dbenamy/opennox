//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"strings"
	"testing"
	"unsafe"
)

type serverOptionsControl struct {
	id   int
	kind string
}

// The controls match the IDs and widget types that the remaining C subpanels use.
// Geometry, strings, images and widget contents are deliberately controlled.
func serverOptionsSubpanelResource(name string, faithful ...bool) string {
	root := 0
	var controls []serverOptionsControl
	switch name {
	case "rulelist.wnd":
		root = 10169
		controls = []serverOptionsControl{{10170, "SCROLLLISTBOX"}, {10179, "VERTSLIDER"}, {10177, "PUSHBUTTON"}, {10178, "PUSHBUTTON"}, {10171, "ENTRYFIELD"}, {10172, "PUSHBUTTON"}, {10173, "PUSHBUTTON"}, {10174, "PUSHBUTTON"}, {10175, "PUSHBUTTON"}, {10176, "USER"}}

	case "access.wnd":
		root = 10151
		controls = []serverOptionsControl{{10123, "SCROLLLISTBOX"}, {10124, "CHECKBOX"}, {10125, "CHECKBOX"}, {10126, "ENTRYFIELD"}, {10127, "CHECKBOX"}, {10128, "ENTRYFIELD"}, {10129, "CHECKBOX"}, {10130, "ENTRYFIELD"}, {10131, "CHECKBOX"}, {10132, "ENTRYFIELD"}, {10133, "ENTRYFIELD"}, {10102, "CHECKBOX"}, {10103, "CHECKBOX"}, {10104, "ENTRYFIELD"}, {10105, "SCROLLLISTBOX"}, {10187, "VERTSLIDER"}, {10185, "PUSHBUTTON"}, {10186, "PUSHBUTTON"}, {10109, "SCROLLLISTBOX"}, {10190, "VERTSLIDER"}, {10188, "PUSHBUTTON"}, {10189, "PUSHBUTTON"}, {10200, "SCROLLLISTBOX"}, {10203, "VERTSLIDER"}, {10201, "PUSHBUTTON"}, {10202, "PUSHBUTTON"}, {10111, "ENTRYFIELD"}, {10112, "PUSHBUTTON"}, {10113, "PUSHBUTTON"}, {10191, "PUSHBUTTON"}, {10192, "PUSHBUTTON"}, {10136, "ENTRYFIELD"}, {10205, "STATICTEXT"}, {10206, "RADIOBUTTON"}, {10207, "RADIOBUTTON"}}
	case "general.wnd":
		root = 10300
		controls = []serverOptionsControl{{10303, "STATICTEXT"}, {10301, "CHECKBOX"}, {10302, "CHECKBOX"}, {10304, "CHECKBOX"}, {10305, "CHECKBOX"}, {10306, "CHECKBOX"}, {10319, "PUSHBUTTON"}}
	case "advanced.wnd":
		root = 10168
		controls = []serverOptionsControl{{10167, "RADIOBUTTON"}, {10164, "RADIOBUTTON"}, {10165, "RADIOBUTTON"}, {10166, "RADIOBUTTON"}, {10148, "PUSHBUTTON"}}
	case "advserv.wnd":
		root = 2100
		controls = []serverOptionsControl{{2101, "STATICTEXT"}, {2102, "CHECKBOX"}, {2103, "CHECKBOX"}, {2130, "PUSHBUTTON"}, {2104, "SCROLLLISTBOX"}, {2105, "STATICTEXT"}, {2106, "RADIOBUTTON"}, {2107, "RADIOBUTTON"}, {2108, "RADIOBUTTON"}, {2109, "RADIOBUTTON"}, {2110, "ENTRYFIELD"}}
	case "objlst.wnd":
		root = 1500
		controls = []serverOptionsControl{{1510, "SCROLLLISTBOX"}, {1513, "PUSHBUTTON"}, {1514, "PUSHBUTTON"}, {1515, "PUSHBUTTON"}, {1516, "PUSHBUTTON"}}
		for id := 1520; id <= 1533; id++ {
			controls = append(controls, serverOptionsControl{id, "CHECKBOX"})
		}
	case "spelllst.wnd":
		root = 1100
		controls = []serverOptionsControl{{1120, "CHECKBOX"}, {1121, "CHECKBOX"}, {1122, "CHECKBOX"}, {1123, "CHECKBOX"}, {1124, "CHECKBOX"}, {1125, "CHECKBOX"}, {1126, "CHECKBOX"}, {1127, "CHECKBOX"}, {1128, "CHECKBOX"}, {1129, "CHECKBOX"}, {1130, "CHECKBOX"}, {1131, "CHECKBOX"}, {1132, "CHECKBOX"}, {1133, "CHECKBOX"}, {1110, "SCROLLLISTBOX"}, {1112, "SCROLLLISTBOX"}, {1113, "PUSHBUTTON"}, {1114, "PUSHBUTTON"}, {1115, "PUSHBUTTON"}, {1116, "PUSHBUTTON"}}

	default:
		panic(name)
	}
	if name == "advserv.wnd" && len(faithful) > 0 && faithful[0] {
		controls = append(controls, serverOptionsControl{2119, "HORZSLIDER"}, serverOptionsControl{2120, "STATICTEXT"})
	}
	var b strings.Builder
	fmt.Fprintf(&b, "FONT = small; WINDOW %d 0 0 300 300 USER; STATUS = ENABLED; CHILD\n", root)
	for _, c := range controls {
		data := ""
		switch c.kind {
		case "SCROLLLISTBOX":
			data = "DATA = 256 1 0 0 1 0 0;"
			if len(faithful) > 0 && faithful[0] && (c.id == 10123 || c.id == 10200) {
				data = "DATA = 256 10 0 1 0 1 0;"
			}
		case "ENTRYFIELD":
			data = "DATA = 64 -1 0 0;"
		case "STATICTEXT":
			data = "DATA = 0 0 WindowDir:Blank;"
		case "VERTSLIDER", "HORZSLIDER":
			data = "DATA = 0 100;"
		case "RADIOBUTTON":
			data = "DATA = 0;"
		}
		fmt.Fprintf(&b, "WINDOW %d 0 0 100 40 %s; STATUS = ENABLED; %s END\n", c.id, c.kind, data)
	}
	b.WriteString("END END")
	return b.String()
}
func (o *serverOptionsOwner) installSubpanels(t *testing.T, faithful ...bool) *[]string {
	for off, name := range map[uintptr]string{127824: "access.wnd", 173556: "general.wnd", 180048: "advserv.wnd"} {
		table := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, off)), 36)
		saved := append([]byte(nil), table...)
		t.Cleanup(func() { copy(table, saved) })
		for i := 0; i < 9; i++ {
			*memmap.PtrPtr(0x587000, off+4*uintptr(i)) = unsafe.Pointer(alloc.InternCString(name))
		}
	}
	old := legacy.Nox_new_window_from_file
	var loads []string
	legacy.Nox_new_window_from_file = func(name string, fn gui.WindowFunc) *gui.Window {
		kind := strings.ToLower(name)
		switch kind {
		case "access.wnd", "laccess.wnd", "waccess.wnd":
			kind = "access.wnd"
		case "general.wnd", "lgeneral.wnd", "wgeneral.wnd":
			kind = "general.wnd"
		case "advserv.wnd", "ladvserv.wnd", "wadvserv.wnd":
			kind = "advserv.wnd"
		case "advanced.wnd", "spelllst.wnd", "rulelist.wnd", "objlst.wnd":
		default:
			return old(name, fn)
		}
		loads = append(loads, kind)
		return newWindowFromString(o.c.GUI, serverOptionsSubpanelResource(kind, faithful...), fn)
	}
	t.Cleanup(func() { legacy.Nox_new_window_from_file = old })
	return &loads
}
