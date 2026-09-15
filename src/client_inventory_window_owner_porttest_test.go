//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/noximage"
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

type inventoryWindowOwner struct {
	missingIdentify bool
	*inventoryDisplayOwner
	windowWords   map[string]*uint32
	windowRegions [][]byte
	loads         []string
	animationRef  *legacy.ImageRef
	amountReady   bool
}

// A newly authored minimal resource exercises the actual GUI parser and widgets.
// Original asset geometry is covered by the separate headless scenario.
const inventoryWindowIdentifyResource = `FONT = small;
WINDOW
 9150 0 0 200 200 USER;
 STATUS = ENABLED+HIDDEN+NOFOCUS;
 CHILD
  WINDOW
   9151 16 2 148 40 STATICTEXT;
   STATUS = ENABLED;
   DATA = 0 0 thing.db:IdentifyDescription;
  END
  WINDOW
   9155 0 40 50 50 PUSHBUTTON;
   STATUS = ENABLED+IMAGE+NOFOCUS;
   STYLE = MOUSETRACK;
  END
  WINDOW
   9156 0 90 180 100 SCROLLLISTBOX;
   STATUS = ENABLED+NOFOCUS;
   DATA = 64 1 0 0 0 0 0;
  END
 END
END
`

func newInventoryWindowOwner(t *testing.T) *inventoryWindowOwner {
	o := &inventoryWindowOwner{inventoryDisplayOwner: newInventoryDisplayOwner(t, "Gold", "Identify", "AutoMap", "QuestGoldPile", "QuestGoldChest")}
	oldMapRefresh := o.c.GUI.ValYYY
	t.Cleanup(func() { o.c.GUI.ValYYY = oldMapRefresh })
	var restore func()
	o.windowWords, restore = legacy.PortTestInventoryWindowWords()
	t.Cleanup(restore)
	for _, r := range [][3]uintptr{{0x852978, 0, 12}, {0x5D4594, 1098380, 16}, {0x5D4594, 1098636, 8400}, {0x5D4594, 1062472, 4}, {0x5D4594, 1062500, 8}, {0x5D4594, 1062532, 4}, {0x5D4594, 1049724, 8}, {0x5D4594, 1064848, 4}, {0x5D4594, 1319068, 208}, {0x587000, 27332, 120}, {0x5D4594, 251624, 4}, {0x5D4594, 1096672, 4}, {0x973A20, 16, 16}, {0x5D4594, 1096424, 4}} {
		b := unsafe.Slice((*byte)(memmap.PtrOff(r[0], r[1])), r[2])
		old := append([]byte(nil), b...)
		o.windowRegions = append(o.windowRegions, b)
		t.Cleanup(func() { copy(b, old) })
	}
	for i, p := range legacy.PortTestInventoryWindowCallbacks() {
		o.c.callbackRefs[p] = 0xed200000 + uint32(i)
	}
	names := []string{"InventoryBase", "InventoryIdentifyBase", "InventoryTray1", "InventoryTray2", "InventoryTray3", "InventoryTraySpecial", "InventoryTrayIdentifyLit", "InventoryTrayMapLit", "InventoryUpButton", "InventoryUpButtonLit", "InventoryDownButton", "InventoryDownButtonLit", "InventorySliderButton", "InventorySliderButtonLit", "InventoryEquipRing", "InventoryQuickItemRing", "InventoryCloseButtonLit", "InventoryJournalButtonLit", "InventoryInventoryButton", "InventoryInventoryButtonLit", "InventoryDollButtonLit", "InventoryStatsButton", "InventoryStatsButtonLit", "GUIFist", "SharedKeyMode", "CurrentWeapon", "CurrentWeaponLit", "DefaultLBUpButton", "DefaultLBUpButtonLit", "DefaultLBUpButtonDis", "DefaultLBDownButton", "DefaultLBDownButtonLit", "DefaultLBDownButtonDis", "DefaultSliderThumb", "DefaultSliderThumbLit", "DefaultSliderThumbDis", "MultiMoveBase", "MultiMoveUpLit", "MultiMoveDownLit", "MultiMoveYesPressed", "MultiMoveNoPressed", "MultiMoveBaseNoTag", "MultiMoveYesPressedNoTag", "MultiMoveNoPressedNoTag"}
	images := make(map[string]*noxrender.Image, len(names))
	for i, n := range names {
		images[n] = o.images[i%len(o.images)]
	}
	oldLoad := legacy.Nox_xxx_gLoadImg
	t.Cleanup(func() { legacy.Nox_xxx_gLoadImg = oldLoad })
	legacy.Nox_xxx_gLoadImg = func(name string) *noxrender.Image {
		o.loads = append(o.loads, name)
		if im := images[name]; im != nil {
			return im
		}
		return oldLoad(name)
	}
	ref, free := alloc.New(legacy.ImageRef{})
	t.Cleanup(free)
	*ref = *o.shiny
	ref.SetName("ExtraLives")
	o.animationRef = ref
	o.c.dataRefs[uint32(uintptr(ref.C()))] = 0xe8000001
	oldAnim := legacy.Nox_xxx_gLoadAnim
	t.Cleanup(func() { legacy.Nox_xxx_gLoadAnim = oldAnim })
	legacy.Nox_xxx_gLoadAnim = func(name string) *legacy.ImageRef {
		if name != "ExtraLives" {
			panic("unowned inventory animation: " + name)
		}
		o.loads = append(o.loads, name)
		return ref
	}
	oldParser := legacy.Nox_new_window_from_file
	t.Cleanup(func() { legacy.Nox_new_window_from_file = oldParser })
	legacy.Nox_new_window_from_file = func(name string, fn gui.WindowFunc) *gui.Window {
		o.loads = append(o.loads, name)
		switch name {
		case "identify.wnd":
			if o.missingIdentify {
				return newWindowFromString(o.c.GUI, "", fn)
			}
			return newWindowFromString(o.c.GUI, inventoryWindowIdentifyResource, fn)
		case "MultMove.wnd":
			return newWindowFromString(o.c.GUI, inventoryWindowAmountResource(), fn)
		default:
			panic("unowned inventory resource: " + name)
		}
	}
	icons := make([]unsafe.Pointer, len(o.images))
	for i, im := range o.images {
		icons[i] = unsafe.Pointer(im)
	}
	t.Cleanup(o.c.srv.PortTestInventoryWindowSpells(icons))
	t.Cleanup(o.releaseAmount)
	return o
}
func (o *inventoryWindowOwner) reset(t *testing.T) {
	o.releaseAmount()
	o.missingIdentify = false
	o.c.GUI.ValYYY = 0
	o.inventoryDisplayOwner.reset(t)
	for _, p := range o.windowWords {
		*p = 0
	}
	for _, b := range o.windowRegions {
		clear(b)
	}
	// Reset only the newly owned words, then reinstall the shared display context.
	*o.displayWords["dword_5d4594_1063636"] = uint32(uintptr(o.smallFont))
	*o.displayWords["dword_5d4594_1062456"] = uint32(uintptr(o.parent.C()))
	*o.displayWords["dword_5d4594_1062476"] = uint32(uintptr(o.parent.C()))
	geometry := blobdata.PortTestInventoryWindowGeometry()
	copy(unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 136192)), len(geometry)), geometry)
	*memmap.PtrPtr(0x852978, 8) = o.items[0].C()
	*txword(o.items[0], 276) = 0
	o.players[0].Journal = nil
	*o.windowWords["dword_587000_136184"] = ^uint32(224)
	nox_win_width, nox_win_height = 640, 480
	*o.meters.NamedWord("nox_win_width"), *o.meters.NamedWord("nox_win_height") = 640, 480
	o.pix = noximage.NewImage16(image.Rect(0, 0, 640, 480))
	o.c.r.SetPixBuffer(o.pix)
	o.c.r.Data().SetClipRect(o.pix.Rect)
	o.c.r.Data().SetClipRect2(image.Rect(0, 0, 639, 479))
	o.c.r.Data().SetRect3(o.pix.Rect)
	*o.c.Viewport() = noxrender.Viewport{Screen: o.pix.Rect, World: o.pix.Rect, Size: o.pix.Rect.Size()}
	*o.windowWords["dword_587000_183456"], *o.windowWords["dword_587000_183460"] = 103, 38
	// The unchanged modifier helpers consume this normally preloaded six-row table.
	for i := 0; i < 6; i++ {
		row := unsafe.Slice((*uint32)(memmap.PtrOff(0x587000, 27332+uintptr(i*20))), 5)
		row[0] = 1 << i
		row[2] = uint32(uintptr(o.images[i].C()))
		key := alloc.InternCString(fmt.Sprintf("WindowEffect%d", i))
		row[3] = uint32(uintptr(unsafe.Pointer(key)))
		o.c.dataRefs[row[3]] = 0xe9000000 + uint32(i)
	}
	*memmap.PtrUint32(0x5D4594, 251624) = 1
	o.loads = nil
}
func (o *inventoryWindowOwner) mainWindow() *gui.Window {
	return (*gui.Window)(unsafe.Pointer(uintptr(*o.windowWords["dword_5d4594_1062456"])))
}
func (o *inventoryWindowOwner) invoke(op int, a, b, c, d uintptr) uint32 {
	return legacy.PortTestInventoryWindow(op, a, b, c, d)
}
func (o *inventoryWindowOwner) construct(t *testing.T) {
	t.Helper()
	if v := o.invoke(15, 0, 0, 0, 0); v == 0 || v != *o.windowWords["dword_5d4594_1062456"] {
		t.Fatalf("inventory constructor return %#x", v)
	}
	o.collect()
}
func inventoryWindowPoint(x, y int) uintptr {
	return uintptr(uint32(uint16(x)) | uint32(uint16(y))<<16)
}

func TestClientInventoryWindowConstruction(t *testing.T) {
	o := newInventoryWindowOwner(t)
	o.reset(t)
	o.construct(t)
	for _, id := range []uint{9102, 9103, 9105, 9107, 9111, 9150, 9151, 9155, 9156} {
		if o.mainWindow().ChildByID(id) == nil {
			t.Errorf("missing inventory child %d", id)
		}
	}
	for i, n := range []string{"Gold", "Identify", "AutoMap"} {
		c := &legacy.PortTestInventoryCells()[(i+1)*21-1]
		if c.Drawable == nil || c.Count != 1 || int(c.Drawable.TypeIDVal) != o.c.Things.TypeByID(n).Index() {
			t.Errorf("hidden-row %s item missing", n)
		}
	}
	if len(o.loads) < 29 {
		t.Fatalf("missing constructor asset requests: %v", o.loads)
	}
	t.Logf("constructor: %d windows, %d asset requests", len(o.windows), len(o.loads))
}
func TestClientInventoryWindowEventAdmission(t *testing.T) {
	o := newInventoryWindowOwner(t)
	o.reset(t)
	for _, ev := range []int{0, 4, 5, 6, 7, 8, 9, 12, 16, 17, 16391} {
		want := uint32(1)
		if ev == 8 || ev == 12 || ev == 16 {
			want = 0
		}
		if got := o.invoke(8, 0, uintptr(ev), 0, 0); got != want {
			t.Errorf("event%d got%d want%d", ev, got, want)
		}
	}
	for state := 0; state < 256; state++ {
		*memmap.PtrUint8(0x5D4594, 1049868) = byte(state)
		want := uint32(0)
		if state == 1 || state == 2 {
			want = 1
		}
		if got := o.invoke(31, 0, 0, 0, 0); got != want {
			t.Fatal(fmt.Sprintf("state%d got%d want%d", state, got, want))
		}
	}
}

func (o *inventoryWindowOwner) releaseAmount() {
	if o.amountReady {
		legacy.Nox_gui_itemAmount_free_4C03E0()
		o.amountReady = false
	}
}
func (o *inventoryWindowOwner) initAmount(t *testing.T) {
	t.Helper()
	if legacy.Nox_gui_itemAmount_init_4BFEF0() == 0 {
		t.Fatal("amount dialog constructor")
	}
	o.amountReady = true
	o.collect()
}
func inventoryWindowAmountResource() string {
	text := "FONT = small; WINDOW 3600 0 0 180 100 USER; STATUS = ENABLED+ABOVE; CHILD "
	for id := 3601; id <= 3607; id++ {
		typ := "PUSHBUTTON"
		data := "STYLE = MOUSETRACK;"
		if id == 3601 || id == 3607 {
			typ = "STATICTEXT"
			data = "DATA = 1 0 WindowDir:Blank;"
		}
		text += fmt.Sprintf("WINDOW %d %d %d 30 15 %s; STATUS = ENABLED+NOFOCUS; %s END ", id, 5+(id-3601)*22, 20+(id%2)*25, typ, data)
	}
	return text + "END END"
}
