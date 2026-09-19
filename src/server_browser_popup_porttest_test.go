//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"golang.org/x/image/font/basicfont"
	"image"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"unsafe"
)

func TestServerBrowserProximityPopup(t *testing.T) {
	o := newListboxOwner(t)
	_, restoreFont := o.c.Render().GetFonts().PortTestWindowFont(basicfont.Face7x13, "default")
	defer restoreFont()
	oldStrings := strMan
	strMan = o.c.srv.Strings()
	defer func() { strMan = oldStrings }()
	words, restore := legacy.PortTestServerBrowserWords()
	defer restore()
	table := serverConfigOwnBytes(t, 0x5D4594, 1307316, 404)
	clear(table)
	binary.LittleEndian.PutUint32(table[400:], 0xbeef1234)
	radius := serverConfigOwnBytes(t, 0x581450, 9720, 8)
	binary.LittleEndian.PutUint64(radius, math.Float64bits(10))
	resource, err := os.ReadFile(filepath.Join(os.Getenv("OPENNOX_SPELLBOOK_ASSETS"), "window", "proxlist.wnd"))
	if err != nil {
		t.Fatal(err)
	}
	loader := legacy.Nox_new_window_from_file
	defer func() { legacy.Nox_new_window_from_file = loader }()
	legacy.Nox_new_window_from_file = func(name string, fn gui.WindowFunc) *gui.Window {
		if name != "proxlist.wnd" {
			t.Fatal("unexpected popup resource", name)
		}
		w := newWindowFromString(o.c.GUI, string(resource), fn)
		if w == nil {
			t.Fatal("popup resource construction")
		}
		for _, id := range []uint{10061, 10062, 10063, 10064} {
			if w.ChildByID(id) == nil {
				t.Fatal("popup resource missing child", id)
			}
		}
		return w
	}
	imageLoader := legacy.Nox_xxx_gLoadImg
	defer func() { legacy.Nox_xxx_gLoadImg = imageLoader }()
	legacy.Nox_xxx_gLoadImg = func(name string) *noxrender.Image { return o.images[0] }
	point, freePoint := alloc.Make([]uint32{}, 2)
	defer freePoint()
	type row struct {
		Count   int
		Miss    bool
		Shown   bool
		Pos     image.Point
		Names   []string
		Indices []int
		Closed  bool
	}
	var rows []row
	for _, count := range []int{0, 1, 3, 100} {
		for _, miss := range []bool{false, true} {
			t.Run(fmt.Sprintf("%d/%t", count, miss), func(t *testing.T) {
				raw, free := alloc.Make([]byte{}, 172*(count+1))
				defer free()
				var records []unsafe.Pointer
				var want []string
				for i := 0; i < count; i++ {
					b := raw[i*172 : (i+1)*172]
					records = append(records, unsafe.Pointer(&b[0]))
					copy(b[12:28], "127.0.0.1")
					binary.LittleEndian.PutUint16(b[109:], 18590)
					binary.LittleEndian.PutUint16(b[44:], 100)
					binary.LittleEndian.PutUint16(b[46:], 80)
					binary.LittleEndian.PutUint32(b[96:], uint32(i*100))
					name := "127.0.0.1:18590"
					if i%2 == 1 {
						name = fmt.Sprintf("Server %d", i)
						copy(b[120:135], name)
					}
					want = append(want, fmt.Sprintf("%s   %dms", name, i*100))
				}
				head, freeHead := legacy.PortTestServerBrowserList(records)
				defer freeHead()
				point[0], point[1] = 100, 80
				if miss {
					point[0] = 101
				}
				p := legacy.PortTestServerBrowserPopup(o.parent.C(), unsafe.Pointer(&point[0]), head)
				shown := count > 0 && !miss
				r := row{Count: count, Miss: miss, Shown: legacy.PortTestServerBrowserPopupShown()}
				if (p != nil) != shown || r.Shown != shown {
					t.Fatal("popup visibility", count, miss, p, r.Shown)
				}
				if shown {
					w := (*gui.Window)(p)
					r.Pos = w.Off
					r.Names = serverPanelsListNames(w.ChildByID(10061))
					if r.Pos != image.Pt(216, 87) || !reflect.DeepEqual(r.Names, want) {
						t.Fatal("popup layout/rows", r.Pos, r.Names, want)
					}
					for i, record := range records {
						if legacy.PortTestServerBrowserPopupAt(i) != record {
							t.Fatal("popup record identity", i)
						}
						r.Indices = append(r.Indices, i)
					}
				}
				if legacy.PortTestServerBrowserPopupAt(int(*words["dword_5d4594_1307720"])) != nil {
					t.Fatal("popup outside count")
				}
				legacy.PortTestServerBrowserPopupClose()
				legacy.PortTestServerBrowserPopupClose()
				r.Closed = !legacy.PortTestServerBrowserPopupShown()
				if !r.Closed || binary.LittleEndian.Uint32(table[400:]) != 0xbeef1234 {
					t.Fatal("popup cleanup/table guard")
				}
				o.c.GUI.FreeDestroyed()
				rows = append(rows, r)
			})
		}
	}
	spellbookCapture(t, "server-browser-proximity-popup", rows, "b2879bab13cb61613d89de00fd656c54683507e4d8a7c799e5addf0de13364f8")
}
