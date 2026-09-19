//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestCharacterCreationInitialAppearance(t *testing.T) {
	o := newEntryOwner(t)
	colors := characterPalette(t)
	o.create(t, 0, 8, func(_ *gui.WindowData, d *gui.EntryFieldData) { d.Field_1040 = 25 })
	o.win.SetID(751)
	var controls []*gui.Window
	var ptrs []unsafe.Pointer
	for i := 0; i < 14; i++ {
		w := o.c.GUI.NewWindowRaw(o.parent, 8, 0, 0, 1, 1, nil)
		id := 720 + i
		if i >= 10 {
			id = 711 + i - 10
		}
		w.SetID(uint(id))
		controls = append(controls, w)
		ptrs = append(ptrs, w.C())
	}
	for i := 0; i < 4; i++ {
		w := o.c.GUI.NewWindowRaw(o.parent, 8, 0, 0, 1, 1, nil)
		w.SetID(uint(731 + i))
	}
	ptrs = append(ptrs, o.win.C())
	player, free := alloc.Calloc(1, 128)
	defer free()
	raw := unsafe.Slice((*byte)(player), 128)
	copy(raw, playerFileName("SavedHero"))
	copy(raw[71:74], colors[3*(64+3):][:3])
	for i, off := range []int{68, 74, 77, 80} {
		copy(raw[off:off+3], colors[3*(32+4+i):][:3])
	}
	for i := 0; i < 5; i++ {
		raw[83+i] = byte(5 + i)
	}
	defer legacy.PortTestCharacterAppearanceOwner(player, ptrs)()
	defer legacy.PortTestCharacterPaletteOwner(o.parent.C(), nil)()
	initial := unsafe.Slice(memmap.PtrUint16(0x587000, 171372), 5)
	oldInitial := append([]uint16(nil), initial...)
	defer copy(initial, oldInitial)
	copy(initial, []uint16{3, 7, 11, 17, 31})
	onceBytes := unsafe.Slice((*byte)(unsafe.Pointer(&modifierColorsOnce)), int(unsafe.Sizeof(modifierColorsOnce)))
	oldOnce := append([]byte(nil), onceBytes...)
	clear(onceBytes)
	defer copy(onceBytes, oldOnce)
	oldColors := o.c.srv.Modif.Colors
	defer func() { o.c.srv.Modif.Colors = oldColors }()
	for bank := 0; bank < 3; bank++ {
		for index := 0; index < 32; index++ {
			rgb := &o.c.srv.Modif.Colors[bank][index]
			off := 3 * (32*bank + index)
			rgb.R = colors[off]
			rgb.G = colors[off+1]
			rgb.B = colors[off+2]
		}
	}
	type row struct {
		Default, Return uint32
		Words           []uint32
		Name            string
	}
	var rows []row
	for _, value := range []uint32{0, 1, 2, 0xffffffff} {
		restore := legacy.PortTestCharacterDefaultOwner(value)
		o.text([]uint16{'B', 'e', 'f', 'o', 'r', 'e'})
		ret := legacy.PortTestCharacterSetup()
		restore()
		words := make([]uint32, 10)
		for i := range words {
			words[i] = *characterWord(controls[i], 32)
		}
		want := []uint32{2 | 2<<16, 1 | 9<<16, 1 | 9<<16, 1 | 9<<16, 1 | 9<<16, 3 << 16, 7 << 16, 11 << 16, 17 << 16, 31 << 16}
		name := "Before"
		if value == 0 {
			want = []uint32{2 | 3<<16, 1 | 4<<16, 1 | 5<<16, 1 | 6<<16, 1 | 7<<16, 5 << 16, 6 << 16, 7 << 16, 8 << 16, 9 << 16}
			name = "SavedHero"
			if ret != 9<<16 {
				t.Fatal("saved return", ret)
			}
		} else if ret != value {
			t.Fatal("default return", ret)
		}
		for i := range want {
			if words[i] != want[i] {
				t.Fatal("initial appearance", value, i, words, want)
			}
		}
		response := o.win.Func94(&gui.RawEvent{Event: 16413})
		p := unsafe.Pointer(uintptr(gui.EventRespInt(response)))
		got := alloc.GoString16((*uint16)(p))
		if got != name {
			t.Fatal("initial name", got, name)
		}
		rows = append(rows, row{value, ret, words, got})
	}
	spellbookCapture(t, "character-creation-initial-appearance", rows, "709996138e19f5cb8fb971fb3b87cccd26cf26217bdcaa199364ef50a413d075")
}
