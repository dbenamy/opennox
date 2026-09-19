//go:build porttest

package opennox

import (
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestCharacterCreationClassEvents(t *testing.T) {
	o := newEntryOwner(t)
	root := o.c.GUI.NewWindowRaw(nil, 8, 0, 0, 90, 90, nil)
	next := o.c.GUI.NewWindowRaw(root, 0, 0, 0, 1, 1, nil)
	next.SetID(610)
	var notices []string
	label := o.c.GUI.NewWindowRaw(root, 8, 0, 0, 1, 1, func(_ *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
		if ev.EventCode() == 16385 {
			a, _ := ev.EventArgsC()
			notices = append(notices, alloc.GoString16((*uint16)(unsafe.Pointer(a))))
		}
		return nil
	})
	label.SetID(605)
	child := o.c.GUI.NewWindowRaw(root, 8, 0, 0, 1, 1, nil)
	p, free := alloc.Calloc(1, 128)
	defer free()
	defer legacy.PortTestCharacterClassOwner(p)()
	defer legacy.PortTestCharacterClassWindow(root.C())()
	nextPtr := memmap.PtrPtr(0x5D4594, 1307728)
	oldNext := *nextPtr
	*nextPtr = next.C()
	defer func() { *nextPtr = oldNext }()
	selected := memmap.PtrUint32(0x5D4594, 1307740)
	oldSelected := *selected
	defer func() { *selected = oldSelected }()
	table := unsafe.Slice(memmap.PtrPtr(0x587000, 170208), 3)
	old := append([]unsafe.Pointer(nil), table...)
	defer copy(table, old)
	var entries []strman.Entry
	for i, name := range []string{"WarriorInfo", "WizardInfo", "ConjurerInfo"} {
		s, done := alloc.CString(name)
		defer done()
		table[i] = unsafe.Pointer(s)
		entries = append(entries, strman.Entry{ID: strman.ID("SelClass.c:" + name), Vals: []strman.Variant{{Str: name}}})
	}
	lang, done := o.c.srv.PortTestMeterStrings(entries...)
	defer done()
	lang(0)
	type row struct {
		Event, ID, Return int
		Class             byte
		Selected, Enabled uint32
		Notices           []string
	}
	var rows []row
	for _, event := range []int{0, 5, 16389, 16391} {
		for _, id := range []int{0, 600, 601, 602, 603, 604, 609, 611} {
			child.SetID(uint(id))
			notices = nil
			*selected = 0x11223344
			*(*byte)(unsafe.Add(p, 66)) = 0x7f
			next.Flags = 0
			ret := legacy.PortTestCharacterClassEvent(root.C(), child.C(), event)
			wantRet := 0
			if event == 16389 || event == 16391 {
				wantRet = 1
			}
			wantClass, wantSelected, wantEnabled := byte(0x7f), uint32(0x11223344), uint32(0)
			if event == 16389 && id >= 601 && id <= 603 {
				wantClass = byte(id - 601)
				wantSelected = uint32(id)
				wantEnabled = 8
				names := []string{"WarriorInfo", "WizardInfo", "ConjurerInfo"}
				if len(notices) != 1 || notices[0] != names[id-601] {
					t.Fatal("class text", notices)
				}
			} else if len(notices) != 0 {
				t.Fatal("unexpected class text")
			}
			gotClass := *(*byte)(unsafe.Add(p, 66))
			if ret != wantRet || gotClass != wantClass || *selected != wantSelected || uint32(next.Flags) != wantEnabled {
				t.Fatal("class selection", event, id, ret, gotClass, *selected, next.Flags)
			}
			rows = append(rows, row{event, id, ret, gotClass, *selected, uint32(next.Flags), append([]string(nil), notices...)})
		}
	}
	spellbookCapture(t, "character-creation-class-events", rows, "0f973e847e3df35b60b4d148bfbaa0c55ab28039788a98a2f4da098596cefdee")
}
func TestCharacterCreationClassShade(t *testing.T) {
	o := newEntryOwner(t)
	w := o.c.GUI.NewWindowRaw(nil, 8, 4, 5, 7, 9, nil)
	w.SetID(601)
	selected := memmap.PtrUint32(0x5D4594, 1307740)
	old := *selected
	defer func() { *selected = old }()
	type row struct {
		Selected uint32
		Pixels   string
	}
	var rows []row
	for _, id := range []uint32{0, 600, 601, 602, 603, 0xffffffff} {
		*selected = id
		for i := range o.pix.Pix {
			o.pix.Pix[i] = 0x7fff
		}
		before := effectsPixelHash(o.pix)
		if legacy.PortTestCharacterClassDraw(w.C()) != 1 {
			t.Fatal("shade return")
		}
		after := effectsPixelHash(o.pix)
		if (after == before) != (id == 601) {
			t.Fatal("selected class shading", id)
		}
		if o.pix.Pix[0] != 0x7fff {
			t.Fatal("shade escaped rectangle")
		}
		rows = append(rows, row{id, after})
	}
	spellbookCapture(t, "character-creation-class-shade", rows, "090f1637c55b06703bd95696345b66eaf91acce739a16b27f288ae20baff0c8b")
}
