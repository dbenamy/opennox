//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/client/audio/ail"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/dialog"
	"github.com/opennox/opennox/v1/legacy/timer"
)

type shopUIOwner struct {
	*inventoryWindowOwner
	shopWords    map[string]*uint32
	shopRegions  [][]byte
	shopReady    bool
	missingShop  bool
	messageOwner bool
}

func newShopUIOwner(t *testing.T) *shopUIOwner {
	o := &shopUIOwner{inventoryWindowOwner: newInventoryWindowOwner(t, "Shopkeeper", "ShopkeeperWarriorsRealm", "ShopkeeperConjurerRealm", "ShopkeeperWizardRealm", "ShopkeeperLandOfTheDead", "ShopkeeperMagicShop", "ShopkeeperPurple", "ShopkeeperYellow")}
	// Use the real dialogue owner; the minimal GUI environment does not initialize
	// the game's global audio modules. Shop actions queue/clear dialogue filenames;
	// this fixture does not pump audio streams.
	oldDialog := legacy.Dialogs
	var state [6]uint32
	state[0] = 1
	var driver ail.Driver
	var timers [4]timer.TimerGroup
	legacy.Dialogs = dialog.NewDialog("dialog", &state[0], &state[1], &state[2], &state[3], &state[4], &driver, &state[5], o.c.srv.Strings, &timers[0], &timers[1], &timers[2], &timers[3], nil, nil, nil, nil, nil, nil, nil)
	t.Cleanup(func() { legacy.Dialogs = oldDialog })
	oldRender, oldRemembered := nox_client_renderGUI_80828, nox_xxx_xxxRenderGUI_587000_80832
	t.Cleanup(func() { nox_client_renderGUI_80828, nox_xxx_xxxRenderGUI_587000_80832 = oldRender, oldRemembered })
	var restore func()
	o.shopWords, restore = legacy.PortTestShopUIWords()
	t.Cleanup(restore)
	for _, r := range [][2]uintptr{{1097300, 52}, {1097352, 1028}, {1098396, 4}, {1098400, 84}, {1098484, 4}, {1098492, 52}, {1098556, 20}, {1098584, 8}, {1098608, 8}, {1097292, 4}, {1107040, 4}, {823804, 1932}} {
		b := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, r[0])), r[1])
		old := append([]byte(nil), b...)
		o.shopRegions = append(o.shopRegions, b)
		t.Cleanup(func() { copy(b, old) })
	}
	for i := range legacy.PortTestShopUICells() {
		o.c.dataRefs[uint32(txptr(unsafe.Pointer(&legacy.PortTestShopUICells()[i])))] = 0xed600001 + uint32(i)
	}
	for i, p := range legacy.PortTestShopUICallbacks() {
		if p != nil {
			o.c.callbackRefs[p] = 0xed700000 + uint32(i)
		}
	}
	oldLoad := legacy.Nox_xxx_gLoadImg
	t.Cleanup(func() { legacy.Nox_xxx_gLoadImg = oldLoad })
	images := map[string]*noxrender.Image{}
	for i, name := range []string{"ShopBase", "ShopTradeMode", "ShopIdentifyMode", "ShopRepairMode", "ShopExitMode", "ShopInventoryBar1", "ShopInventoryBar2", "ShopInventorySlider", "ShopInventorySliderSelected", "ShopInventoryUp", "ShopInventoryUpSelected", "ShopInventorydown", "ShopInventorydownSelected", "ShopTextBorder", "ShopkeeperPic", "ShopkeeperWarriorPic", "ShopkeeperConjurerPic", "ShopkeeperWizardPic", "ShopkeeperLandOfTheDeadPic", "ShopkeeperMagicShopPic", "ShopKeeperPic", "ShopKeeperWarriorPic", "ShopKeeperConjurerPic", "ShopKeeperWizardPic", "ShopKeeperLandOfTheDeadPic", "ShopKeeperMagicShopPic", "ShopKeeperPurplePic", "ShopKeeperBrownPic"} {
		images[name] = o.images[i%len(o.images)]
	}
	images["ShopInventoryDown"] = images["ShopInventorydown"]
	images["ShopInventoryDownSelected"] = images["ShopInventorydownSelected"]
	legacy.Nox_xxx_gLoadImg = func(name string) *noxrender.Image {
		if im := images[name]; im != nil {
			o.loads = append(o.loads, name)
			return im
		}
		return oldLoad(name)
	}
	oldResource := legacy.Nox_new_window_from_file
	t.Cleanup(func() { legacy.Nox_new_window_from_file = oldResource })
	legacy.Nox_new_window_from_file = func(name string, fn gui.WindowFunc) *gui.Window {
		if name == "Shop.wnd" {
			o.loads = append(o.loads, name)
			if o.missingShop {
				return newWindowFromString(o.c.GUI, "", fn)
			}
			return newWindowFromString(o.c.GUI, shopUIResource(), fn)
		}
		return oldResource(name, fn)
	}
	// Preserve all existing fixture strings when adding shop-specific entries.
	path := filepath.Join(t.TempDir(), "shop-strings.json")
	if err := o.c.srv.Strings().WriteJSON(path, false); err != nil {
		t.Fatal(err)
	}
	var data struct {
		Lang    int            `json:"lang"`
		Entries []strman.Entry `json:"entries"`
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(raw, &data); err != nil {
		t.Fatal(err)
	}
	ids := make(map[strman.ID]bool)
	for _, e := range data.Entries {
		ids[e.ID] = true
	}
	for i, pair := range [][2]string{{"SellInstructions", "Choose items to sell."}, {"RepairInstructions", "Choose equipment to repair."}, {"NotEnoughGold", "Need %d more gold."}, {"ShopInformationTitle", "Shop information"}, {"BuyLabel", "Buy quantity"}, {"SellLabel", "Sell quantity"}, {"RepairLabel", "Repair item"}, {"FixtureGreeting", "Welcome to the shop."}, {"FixtureSilentGreeting", "A quiet welcome."}} {
		id := strman.ID("GUIShop.c:" + pair[0])
		if !ids[id] {
			v := strman.Variant{Str: pair[1]}
			if pair[0] == "FixtureGreeting" {
				v.Str2 = "shop-greeting.wav"
			}
			data.Entries = append(data.Entries, strman.Entry{ID: id, Vals: []strman.Variant{v}})
			o.c.dataRefs[uint32(txptr(unsafe.Pointer(alloc.InternCString16(v.Str))))] = 0xed800001 + uint32(i)
			if v.Str2 != "" {
				o.c.dataRefs[uint32(txptr(unsafe.Pointer(alloc.InternCString(v.Str2))))] = 0xed900002
			}
		}
	}

	o.c.dataRefs[uint32(txptr(unsafe.Pointer(alloc.InternCString(""))))] = 0xed900001

	raw, err = json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	if err = o.c.srv.Strings().ReadJSON(path); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(o.releaseShop)
	return o
}
func shopUIResource() string {
	s := "FONT = small; WINDOW 3800 0 0 340 300 USER; STATUS = ENABLED+ABOVE; CHILD "
	for _, v := range [][5]int{{3801, 0, 0, 50, 20}, {3802, 50, 0, 50, 20}, {3803, 100, 0, 50, 20}, {3804, 150, 0, 50, 20}, {3805, 200, 0, 100, 40}, {3806, 0, 50, 300, 200}, {3807, 305, 50, 20, 200}, {3808, 305, 25, 20, 20}, {3809, 305, 255, 20, 20}, {3810, 0, 275, 300, 20}} {
		typ, data := "PUSHBUTTON", "STYLE = MOUSETRACK;"
		if v[0] == 3805 || v[0] == 3806 {
			typ, data = "USER", ""
		}
		if v[0] == 3807 {
			typ, data = "VERTSLIDER", "DATA = 0 300;"
		}
		if v[0] == 3810 {
			typ, data = "STATICTEXT", "DATA = 1 0 WindowDir:Blank;"
		}
		s += fmt.Sprintf("WINDOW %d %d %d %d %d %s; STATUS = ENABLED+NOFOCUS; %s END ", v[0], v[1], v[2], v[3], v[4], typ, data)
	}
	return s + "END END"
}
func (o *shopUIOwner) releaseShop() {
	if o.shopReady {
		legacy.PortTestShopUI(13)
		o.shopReady = false
	}
}
func (o *shopUIOwner) reset(t *testing.T) {
	o.releaseShop()
	legacy.Dialogs.Sub_44D8F0()
	o.inventoryWindowOwner.reset(t)
	nox_client_renderGUI_80828, nox_xxx_xxxRenderGUI_587000_80832 = true, false
	for _, p := range o.shopWords {
		*p = 0
	}
	for _, b := range o.shopRegions {
		clear(b)
	}
	o.missingShop = false
}
func (o *shopUIOwner) shopWindow() *gui.Window {
	return (*gui.Window)(unsafe.Pointer(uintptr(*o.shopWords["dword_5d4594_1098576"])))
}
func (o *shopUIOwner) constructShop(t *testing.T) {
	t.Helper()
	if ret := legacy.PortTestShopUI(4); ret != 1 {
		t.Fatalf("shop constructor %d", ret)
	}
	o.shopReady = true
	o.collect()
}

type shopUIResult struct {
	RenderGUI    bool
	MessageState []uint32
	Dialogue     [2]string
	Case, Op     int
	Return       uint32
	Named        map[string]uint32
	Regions      [][]byte
	Window       inventoryWindowResult
}

func (o *shopUIOwner) shopCapture(t *testing.T, id, op int, args ...uintptr) shopUIResult {
	// Keep identities for known widgets whose cached C addresses outlive Destroy.
	// Record them while still live, without reading a destroyed widget later.
	o.collect()
	for i, w := range o.windows {
		o.c.dataRefs[uint32(txptr(w.C()))] = 0xe1000001 + uint32(i)
	}
	ret := legacy.PortTestShopUI(op, args...)
	r := shopUIResult{Case: id, Op: op, Return: o.norm(ret), Named: make(map[string]uint32)}
	r.Dialogue = [2]string{legacy.Dialogs.FileToRead(), legacy.Dialogs.CurrentPlayingFile()}
	r.Window = o.capture(t, id, 0, 7, 0, 0, 0, 0)
	r.RenderGUI = nox_client_renderGUI_80828
	if o.messageOwner {
		r.MessageState = []uint32{memmap.Uint32(0x5D4594, 830240), o.norm(uint32(txptr(nox_gui_curDialog_830224.C()))), o.norm(uint32(txptr(dword_5d4594_830228.C()))), o.norm(uint32(txptr(dword_5d4594_830232.C()))), o.norm(uint32(txptr(dword_5d4594_830236.C())))}
	}

	for n, p := range o.shopWords {
		v := *p
		switch n {
		case "dword_5d4594_1098456", "dword_5d4594_1098576", "dword_5d4594_1098580", "dword_5d4594_1098596", "dword_5d4594_1098600", "dword_5d4594_1098604":
			v = o.norm(v)
		}
		r.Named[n] = v
	}
	for i, b := range o.shopRegions {
		cp := append([]byte(nil), b...)
		norm := func(off int) { binary.LittleEndian.PutUint32(cp[off:], o.norm(binary.LittleEndian.Uint32(cp[off:]))) }
		if i == 8 {
			norm(0)
		}
		if i == 3 || i == 7 {
			for off := 0; off < len(cp); off += 4 {
				norm(off)
			}
		}
		r.Regions = append(r.Regions, cp)
	}

	return r
}
func shopUICapture(t *testing.T, label string, rows []shopUIResult, want string) {
	t.Helper()
	b, err := json.Marshal(rows)
	if err != nil {
		t.Fatal(err)
	}
	if p := os.Getenv("OPENNOX_SHOP_UI_CAPTURE"); p != "" {
		if err = os.WriteFile(p+"-"+label+".json", b, 0600); err != nil {
			t.Fatal(err)
		}
	}
	got := fmt.Sprintf("%x", sha256.Sum256(b))
	t.Logf("%s: %d results %s", label, len(rows), got)
	if want != "" && got != want {
		t.Fatalf("%s hash %s want frozen expectation %s", label, got, want)
	}
}
