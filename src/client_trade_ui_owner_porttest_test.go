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
)

type tradeUIOwner struct {
	*inventoryWindowOwner
	tradeWords                  map[string]*uint32
	tradeRegions                [][]byte
	tradeReady                  bool
	missingTrade, missingAmount bool
	callbacks                   [][7]uint32
	accept, cancel              unsafe.Pointer
}

func newTradeUIOwner(t *testing.T) *tradeUIOwner {
	o := &tradeUIOwner{inventoryWindowOwner: newInventoryWindowOwner(t)}
	var restore func()
	o.tradeWords, restore = legacy.PortTestTradeUIWords()
	t.Cleanup(restore)
	for _, r := range [][2]uintptr{{1319284, 560}, {1320308, 560}, {1319844, 128}, {1319972, 128}, {1320100, 64}, {1320164, 24}, {1320188, 52}, {1320240, 64}, {1320304, 4}, {1320868, 64}, {1320960, 4}, {823804, 1932}} {
		b := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, r[0])), r[1])
		old := append([]byte(nil), b...)
		o.tradeRegions = append(o.tradeRegions, b)
		t.Cleanup(func() { copy(b, old) })
	}
	for side, cells := range legacy.PortTestTradeUICells() {
		for i := range cells {
			o.c.dataRefs[uint32(txptr(unsafe.Pointer(&cells[i])))] = 0xed500001 + uint32(side*4+i)
		}
	}
	for i, p := range legacy.PortTestTradeUICallbacks() {
		o.c.callbackRefs[p] = 0xed400000 + uint32(i)
	}
	o.accept, o.cancel, restore = legacy.PortTestTradeUIObserve(func(v [7]uint32) { o.callbacks = append(o.callbacks, v) })
	t.Cleanup(restore)
	o.c.callbackRefs[o.accept], o.c.callbackRefs[o.cancel] = 0xed410001, 0xed410002
	oldLoad := legacy.Nox_xxx_gLoadImg
	t.Cleanup(func() { legacy.Nox_xxx_gLoadImg = oldLoad })
	images := map[string]*noxrender.Image{}
	for i, name := range []string{"TradeBase", "TradeLeftAcceptPushed", "TradeLeftAcceptLit", "TradeRightAcceptLit", "TradeCancelLit", "TradeGold"} {
		images[name] = o.images[i]
	}
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
		if name == "Trade.wnd" {
			o.loads = append(o.loads, name)
			if o.missingTrade {
				return newWindowFromString(o.c.GUI, "", fn)
			}
			return newWindowFromString(o.c.GUI, tradeUIResource(), fn)
		}
		if name == "MultMove.wnd" && o.missingAmount {
			o.loads = append(o.loads, name)
			return newWindowFromString(o.c.GUI, "", fn)
		}
		return oldResource(name, fn)
	}
	// Preserve all existing fixture strings when adding trade-specific entries.
	path := filepath.Join(t.TempDir(), "trade-strings.json")
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
	for _, pair := range [][2]string{{"TradeMain", "Trade"}, {"TradePlayerName", "Player"}, {"TradeVendorName", "Partner"}, {"TradePlayerAccept", "Accept"}, {"TradeVendorAccept", "Partner accepted"}, {"TradeCancel", "Cancel trade"}, {"TotalValueLabel", "Value %d"}, {"TradeGUIItemNotFound", "Trade item missing"}, {"guiinv.c:DrawablesExhausted", "Drawables exhausted"}} {
		id := strman.ID("guitrade.c:" + pair[0])
		if pair[0] == "guiinv.c:DrawablesExhausted" {
			id = strman.ID(pair[0])
		}
		if !ids[id] {
			data.Entries = append(data.Entries, strman.Entry{ID: id, Vals: []strman.Variant{{Str: pair[1]}}})
		}
	}
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
	t.Cleanup(o.releaseTrade)
	return o
}
func tradeUIResource() string {
	s := "FONT = small; WINDOW 3700 0 0 320 220 USER; STATUS = ENABLED+ABOVE; CHILD "
	for _, v := range [][5]int{{3702, 5, 5, 140, 20}, {3703, 170, 5, 140, 20}, {3704, 10, 40, 100, 100}, {3705, 180, 40, 100, 100}, {3708, 10, 170, 60, 20}, {3709, 180, 170, 60, 20}, {3710, 110, 190, 70, 20}, {3711, 10, 145, 90, 20}, {3712, 110, 145, 65, 20}, {3713, 180, 145, 90, 20}} {
		typ, data := "STATICTEXT", "DATA = 1 0 WindowDir:Blank;"
		if v[0] == 3704 || v[0] == 3705 {
			typ, data = "USER", ""
		}
		if v[0] >= 3708 && v[0] <= 3710 {
			typ, data = "PUSHBUTTON", "STYLE = MOUSETRACK;"
		}
		s += fmt.Sprintf("WINDOW %d %d %d %d %d %s; STATUS = ENABLED+NOFOCUS; %s END ", v[0], v[1], v[2], v[3], v[4], typ, data)
	}
	return s + "END END"
}
func (o *tradeUIOwner) releaseTrade() {
	if o.tradeReady {
		legacy.PortTestTradeUI(23)
		legacy.PortTestTradeUI(19)
		o.tradeReady = false
	}
}
func (o *tradeUIOwner) reset(t *testing.T) {
	o.releaseTrade()
	o.inventoryWindowOwner.reset(t)
	for _, p := range o.tradeWords {
		*p = 0
	}
	for _, b := range o.tradeRegions {
		clear(b)
	}
	o.callbacks = nil
	o.missingTrade = false
	o.missingAmount = false
}
func (o *tradeUIOwner) tradeWindow() *gui.Window {
	return (*gui.Window)(unsafe.Pointer(uintptr(*o.tradeWords["dword_5d4594_1320940"])))
}
func (o *tradeUIOwner) constructTrade(t *testing.T) {
	t.Helper()
	if ret := legacy.PortTestTradeUI(34); ret != 1 {
		t.Fatalf("trade constructor %d", ret)
	}
	o.tradeReady = true
	o.collect()
}

type tradeUIResult struct {
	Case, Op  int
	Return    uint32
	Named     map[string]uint32
	Regions   [][]byte
	Callbacks [][7]uint32
	Window    inventoryWindowResult
}

func (o *tradeUIOwner) tradeCapture(t *testing.T, id, op int, args ...uintptr) tradeUIResult {
	ret := legacy.PortTestTradeUI(op, args...)
	r := tradeUIResult{Case: id, Op: op, Return: o.norm(ret), Named: make(map[string]uint32), Callbacks: append([][7]uint32(nil), o.callbacks...)}
	r.Window = o.capture(t, id, 0, 7, 0, 0, 0, 0)
	for n, p := range o.tradeWords {
		v := *p
		if n != "dword_5d4594_1320944" && n != "dword_5d4594_1320948" {
			v = o.norm(v)
		}
		r.Named[n] = v
	}
	for i, b := range o.tradeRegions {
		cp := append([]byte(nil), b...)
		norm := func(off int) { binary.LittleEndian.PutUint32(cp[off:], o.norm(binary.LittleEndian.Uint32(cp[off:]))) }
		if i == 0 || i == 1 {
			for off := 0; off < len(cp); off += 140 {
				norm(off)
			}
		}
		if i == 5 {
			for off := 0; off < len(cp); off += 4 {
				norm(off)
			}
		}
		r.Regions = append(r.Regions, cp)
	}
	return r
}
func tradeUICapture(t *testing.T, label string, rows []tradeUIResult, want string) {
	t.Helper()
	b, err := json.Marshal(rows)
	if err != nil {
		t.Fatal(err)
	}
	if p := os.Getenv("OPENNOX_TRADE_UI_CAPTURE"); p != "" {
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
