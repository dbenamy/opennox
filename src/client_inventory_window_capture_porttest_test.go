//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

type inventoryWindowWidget struct {
	Window, Data []uint32
	Text         string
	Rows         []gui.ScrollListBoxItem
}
type inventoryWindowResult struct {
	Case, Step, Op   int
	Return, Captured uint32
	MapRefresh       int
	Named            map[string]uint32
	Regions          [][]byte
	Windows          []inventoryWindowWidget
	Loads            []string
	Inventory        inventoryTransactionResult
	Display          inventoryDisplayResult
}

func (o *inventoryWindowOwner) capture(t *testing.T, id, step, op int, a, b, c, d uintptr) inventoryWindowResult {
	t.Helper()
	ret := o.invoke(op, a, b, c, d)
	o.collect()
	r := inventoryWindowResult{Case: id, Step: step, Op: op, Return: o.norm(ret), Named: make(map[string]uint32), Loads: append([]string(nil), o.loads...)}
	r.MapRefresh = o.c.GUI.ValYYY
	r.Captured = o.norm(uint32(uintptr(o.c.GUI.Captured().C())))
	pointers := map[string]bool{"dword_5d4594_1049976": true, "dword_5d4594_1049992": true, "dword_5d4594_1049996": true, "dword_5d4594_1050008": true, "dword_5d4594_1062452": true, "dword_5d4594_1062456": true, "dword_5d4594_1062468": true, "dword_5d4594_1062476": true, "dword_5d4594_1062480": true, "dword_5d4594_1062492": true, "dword_5d4594_1062508": true, "dword_5d4594_1062524": true, "dword_5d4594_1062528": true, "dword_5d4594_1063116": true, "dword_5d4594_1063120": true, "dword_5d4594_1063636": true, "nox_win_unk5": true, "nox_gui_itemAmount_dialog_1319228": true, "nox_gui_itemAmount_item_1319256": true, "dword_5d4594_1319232": true, "dword_5d4594_1319236": true, "dword_5d4594_1309820": true, "dword_5d4594_1321228": true, "nox_wnd_quitMenu_825760": true}
	for n, p := range o.windowWords {
		v := *p
		if pointers[n] {
			v = o.norm(v)
		}
		r.Named[n] = v
	}
	for i, b := range o.windowRegions {
		cp := append([]byte(nil), b...)
		norm := func(off int) { binary.LittleEndian.PutUint32(cp[off:], o.norm(binary.LittleEndian.Uint32(cp[off:]))) }
		switch i {
		case 8:
			norm(32)
			norm(92)
			for off := 128; off <= 156; off += 4 {
				norm(off)
			}
		case 9:
			for off := 0; off < len(cp); off += 20 {
				norm(off + 4)
				norm(off + 8)
				norm(off + 12)
			}
		case 0:
			norm(8)
		case 2:
			for off := 0; off < len(cp); off += 140 {
				norm(off)
			}
		case 3, 4, 5, 12:
			for off := 0; off < len(cp); off += 4 {
				norm(off)
			}
		}
		r.Regions = append(r.Regions, cp)
	}
	for i, w := range o.windows {
		v := inventoryWindowWidget{Window: append([]uint32(nil), unsafe.Slice((*uint32)(w.C()), 101)...)}
		for _, n := range []int{13, 15, 17, 19, 21, 23, 59, 93, 94, 95, 96, 97, 98, 99, 100} {
			v.Window[n] = o.norm(v.Window[n])
		}
		if w.WidgetData != nil && uintptr(w.WidgetData) > 4096 {
			v.Window[8] = 0xef100000 + uint32(i)
			words := func(size uintptr) []uint32 {
				return append([]uint32(nil), unsafe.Slice((*uint32)(w.WidgetData), size/4)...)
			}
			switch style := w.DrawData().Style; {
			case style&(gui.StyleVertSlider|gui.StyleHorizSlider) != 0:
				v.Data = words(unsafe.Sizeof(gui.SliderData{}))
			case style&gui.StyleStaticText != 0:
				data := (*gui.StaticTextData)(w.WidgetData)
				v.Data = words(unsafe.Sizeof(*data))
				if data.Text != nil {
					v.Data[0] = 0xef200000 + uint32(i)
					v.Text = alloc.GoString16(data.Text)
				}
			case style&gui.StyleScrollListBox != 0:
				data := (*gui.ScrollListBoxData)(w.WidgetData)
				v.Data = words(unsafe.Sizeof(*data))
				if data.Items != nil {
					v.Data[6] = 0xef300000 + uint32(i)
					v.Rows = append([]gui.ScrollListBoxItem(nil), unsafe.Slice(data.Items, int(data.Count))...)
				}
				for _, j := range []int{7, 8, 9} {
					v.Data[j] = o.norm(v.Data[j])
				}
				if data.Field_3 != 0 {
					t.Fatal("multi-select list needs an owned selection-array capture")
				}
			default:
				t.Fatalf("unowned window%d widget data style%x", w.ID(), style)
			}
		}
		r.Windows = append(r.Windows, v)
	}
	r.Inventory = o.inventoryTransactionOwner.snapshot(t, id, op, ret)
	r.Inventory.Regions = nil
	r.Display = o.displaySnapshot(t, id, -1, uint64(ret))
	return r
}
func inventoryWindowCapture(t *testing.T, label string, rows []inventoryWindowResult, want string) {
	t.Helper()
	b, err := json.Marshal(rows)
	if err != nil {
		t.Fatal(err)
	}
	if prefix := os.Getenv("OPENNOX_INVENTORY_WINDOW_CAPTURE"); prefix != "" {
		if err := os.WriteFile(prefix+"-"+label+".json", b, 0600); err != nil {
			t.Fatal(err)
		}
	}
	hash := fmt.Sprintf("%x", sha256.Sum256(b))
	t.Logf("%s: %d results %s", label, len(rows), hash)
	// Empty expectations are only for the in-progress C baseline. Lock before qualification.
	if want != "" && hash != want {
		t.Fatalf("%s hash%s want frozen C%s", label, hash, want)
	}
}
