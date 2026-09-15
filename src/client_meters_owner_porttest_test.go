//go:build porttest

package opennox

import (
	"fmt"
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noximage"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/input"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"golang.org/x/image/font/basicfont"
)

type meterOwner struct {
	grid      []byte
	equipment []uint32
	text      []meterTextCall
	sounds    [][2]int
	*entryOwner
	meters  *legacy.PortTestMeterEnvironment
	font    unsafe.Pointer
	windows []*gui.Window
}

func newMeterOwner(t *testing.T, extraNames ...string) *meterOwner {
	o := &meterOwner{entryOwner: newEntryOwner(t, append([]string{"RedPotion", "BluePotion", "CurePoisonPotion", "RedApple", "Meat", "Quiver", "Bow"}, extraNames...)...)}
	t.Cleanup(legacy.PortTestClientSoundObserver(func(id, volume int) { o.sounds = append(o.sounds, [2]int{id, volume}) }))
	legacy.GetClient = func() legacy.Client { return &meterClient{o.objectDrawingOwner.proxy, o} }
	var restoreInventory func()
	o.grid, o.equipment, restoreInventory = legacy.PortTestMeterInventory()
	t.Cleanup(restoreInventory)
	oldNet := o.c.srv.NetList
	o.c.srv.NetList = netlist.New()
	o.c.srv.NetList.Init()
	t.Cleanup(func() { o.c.srv.NetList.Free(); o.c.srv.NetList = oldNet })
	for _, reg := range blobdata.PortTestMeterTables() {
		b := unsafe.Slice((*byte)(memmap.PtrOff(reg.Base, reg.Offset)), len(reg.Data))
		old := append([]byte(nil), b...)
		copy(b, reg.Data)
		t.Cleanup(func() { copy(b, old) })
	}
	for _, reg := range [][2]uintptr{{1049732, 52}, {1050000, 4}, {1063640, 4}, {1317000, 2048}, {1309840, 1280}} {
		b := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, reg[0])), int(reg[1]))
		old := append([]byte(nil), b...)
		clear(b)
		t.Cleanup(func() { copy(b, old) })
	}
	legacy.Sub_4AEE30()
	o.meters = legacy.PortTestNewMeterEnvironment()
	t.Cleanup(o.meters.Restore)
	language, restoreStrings := o.c.srv.Server.PortTestMeterStrings()
	t.Cleanup(restoreStrings)
	language(0)
	oldStrings := strMan
	strMan = o.c.srv.Strings()
	t.Cleanup(func() { strMan = oldStrings })
	nox_win_width, nox_win_height = 256, 256
	for i, name := range []string{"RedPotion", "BluePotion", "CurePoisonPotion", "RedApple", "Meat", "Quiver", "Bow"} {
		typ := o.c.Things.TypeByID(name)
		data, freeData := alloc.Make([]uint32{}, 2)
		t.Cleanup(freeData)
		data[0], data[1] = 8, uint32(uintptr(o.images[i].C()))
		typ.DrawFunc = legacy.PortTestSpriteAnimationCallback(2)
		typ.DrawData = unsafe.Pointer(&data[0])
		typ.ObjClass, typ.ObjFlags = 0x10, 0x40000000
		o.c.dataRefs[uint32(uintptr(typ.DrawData))] = 0xe7000010 + uint32(i)
		o.c.callbackRefs[typ.DrawFunc] = 0xe5000002
	}
	for i, off := range []uintptr{1090316, 1090852, 1091388} {
		o.c.refs[(*client.Drawable)(memmap.PtrOff(0x5D4594, off))] = 0xea000001 + uint32(i)
	}
	var free func()
	o.font, free = o.c.Render().GetFonts().PortTestWindowFont(basicfont.Face7x13, "small")
	t.Cleanup(free)
	oldTimes := inputKeyTimeoutsOld
	inputKeyTimeoutsOld = make(map[byte]uint32)
	t.Cleanup(func() { inputKeyTimeoutsOld = oldTimes })
	o.c.ctrl = new(CtrlEventHandler)
	o.seat = &entrySeat{}
	o.c.Inp = input.New(o.c.Log, o.seat, false, 0)
	o.pix = noximage.NewImage16(image.Rect(0, 0, 256, 256))
	o.c.r.SetPixBuffer(o.pix)
	o.c.r.circleSeg.Init(o.c.r)
	t.Cleanup(o.c.r.circleSeg.Free)
	o.c.r.Data().SetClipRect(o.pix.Rect)
	o.c.r.Data().SetClipRect2(image.Rect(0, 0, 255, 255))
	o.c.r.Data().SetRect3(o.pix.Rect)
	*o.c.Viewport() = noxrender.Viewport{Screen: o.pix.Rect, World: o.pix.Rect, Size: o.pix.Rect.Size()}
	oldLoad := legacy.Nox_xxx_gLoadImg
	t.Cleanup(func() { legacy.Nox_xxx_gLoadImg = oldLoad })
	names := map[string]int{"PoisonTube": 10, "HealthManaTubes": 11, "WarriorPoisonTube": 12, "WarriorHealthTube": 13}
	for i := 1; i <= 10; i++ {
		names[fmt.Sprintf("HealthMana%d", i)] = i - 1
	}
	legacy.Nox_xxx_gLoadImg = func(name string) *noxrender.Image {
		i, ok := names[name]
		if !ok {
			panic("unowned meter image: " + name)
		}
		o.named = append(o.named, name)
		return o.images[i]
	}
	return o
}
func (o *meterOwner) collect() {
	o.windows = nil
	var add func(*gui.Window)
	add = func(w *gui.Window) {
		o.windows = append(o.windows, w)
		for c := w.Field100Ptr; c != nil; c = c.Prev() {
			add(c)
		}
	}
	for w := o.c.GUI.Last(); w != nil; w = w.Prev() {
		add(w)
	}
}
func (o *meterOwner) normalize(v uint32) uint32 {
	if v == 0 {
		return 0
	}
	if v == uint32(uintptr(o.font)) {
		return 0xe3000001
	}
	for i, w := range o.windows {
		if v == uint32(uintptr(w.C())) {
			return 0xe1000001 + uint32(i)
		}
	}
	if n, ok := o.c.imageRefs[v]; ok {
		return n
	}
	if n, ok := o.c.callbackRefs[unsafe.Pointer(uintptr(v))]; ok {
		return n
	}
	for dr, n := range o.c.refs {
		if v == uint32(uintptr(dr.C())) {
			return n
		}
	}
	if n, ok := o.c.dataRefs[v]; ok {
		return n
	}
	for i := 1; o.c.Things.TypeByInd(i) != nil; i++ {
		if typ := o.c.Things.TypeByInd(i); typ != nil && v == uint32(uintptr(typ.C())) {
			return 0xe6000000 + uint32(i)
		}
	}

	for op := 0; op < 37; op++ {
		p := legacy.PortTestMeterCallback(op)
		if p != nil && v == uint32(uintptr(p)) {
			return 0xed000000 + uint32(op)
		}
	}
	for i := range o.players {
		if v == uint32(uintptr(unsafe.Pointer(&o.players[i]))) {
			return 0xe4000001 + uint32(i)
		}
	}
	for _, r := range [][3]uintptr{{0x5D4594, 1090276, 6032}, {0x587000, 147904, 488}, {0x5D4594, 1096672, 516}} {
		start := uint32(uintptr(memmap.PtrOff(r[0], r[1])))
		if v >= start && v <= start+uint32(r[2]) {
			return 0xf0000000 + uint32(r[1]) + v - start
		}
	}
	return v
}

type meterResult struct {
	ImageNames           []string
	Messages             [][]byte
	Inventory, Equipment []uint32
	Text                 []meterTextCall
	Sounds               [][2]int
	Case, Step, Op       int
	Return               uint32
	Named, Records       []uint32
	Regions, Windows     [][]uint32
	PlayerBytes          []byte
	Cooldowns            [2]uint32
	Render               objectRenderResult
}

func (o *meterOwner) invoke(t *testing.T, id, step, op int, w *gui.Window, a, b, c, d int) meterResult {
	ret := legacy.PortTestMeterCall(op, w, a, b, c, d)
	o.collect()
	r := meterResult{Case: id, Step: step, Op: op, Return: o.normalize(ret)}
	o.c.srv.NetList.ByInd(31, netlist.Kind0).Each(func(b []byte) bool { r.Messages = append(r.Messages, append([]byte(nil), b...)); return false })
	r.Inventory = append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(&o.grid[0])), len(o.grid)/4)...)
	for i := 0; i < 84; i++ {
		r.Inventory[i*37] = o.normalize(r.Inventory[i*37])
	}
	r.Equipment = append([]uint32(nil), o.equipment...)
	for i := range r.Equipment {
		r.Equipment[i] = o.normalize(r.Equipment[i])
	}
	r.ImageNames = append([]string(nil), o.named...)
	r.Sounds = append([][2]int(nil), o.sounds...)
	r.Text = append([]meterTextCall(nil), o.text...)
	r.Named = o.meters.Named()
	for i := range r.Named {
		r.Named[i] = o.normalize(r.Named[i])
	}
	r.Records = append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(&o.meters.Records[0])), 35)...)
	for i := 0; i < 7; i++ {
		r.Records[5*i] = o.normalize(r.Records[5*i])
	}
	for ri, region := range o.meters.Regions {
		words := append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(&region[0])), len(region)/4)...)
		if ri == 0 {
			for _, offset := range []int{1090296, 1090832, 1091368, 1091900, 1092996, 1093000, 1093004, 1093008, 1093012, 1093016, 1093020, 1093024, 1093028, 1093032} {
				i := (offset - 1090276) / 4
				words[i] = o.normalize(words[i])
			}
			// The embedded potion drawables use the actual type/link initializer.
			for _, offset := range []int{1090316, 1090852, 1091388} {
				start := (offset - 1090276) / 4
				for _, slot := range []int{2, 75, 76, 83, 84, 87, 88, 90, 91, 92, 93, 94, 95, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 114, 115, 116, 124} {
					words[start+slot] = o.normalize(words[start+slot])
				}
			}
		}
		r.Regions = append(r.Regions, words)
	}
	for _, w := range o.windows {
		words := append([]uint32(nil), unsafe.Slice((*uint32)(w.C()), 101)...)
		for _, i := range []int{13, 15, 17, 19, 21, 23, 59, 93, 94, 95, 96, 97, 98, 99, 100} {
			words[i] = o.normalize(words[i])
		}
		r.Windows = append(r.Windows, words)
	}
	r.PlayerBytes = append([]byte(nil), unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 2240)), 16)...)
	r.Cooldowns = [2]uint32{inputKeyTimeoutsOld[4], inputKeyTimeoutsOld[17]}
	r.Render = o.renderResult(t, id, step, int(ret))
	// A returned raw reference must not reappear unnormalized in nested metadata.
	r.Render.Draw.Return = int(r.Return)
	return r
}

func (o *meterOwner) plain(t *testing.T) {
	o.c.GUI.DestroyAll()
	o.c.GUI.FreeDestroyed()
	o.c.GUI.FreeDestroyed()
	o.parent = o.c.GUI.NewWindowRaw(nil, 8, 4, 6, 220, 220, nil)
	o.win = nil
	o.meters.Reset()
	// Override the reused effects fixture's 0xffff white with opaque RGB555.
	// The former means transparent in the text renderer's RGBA5551 format.
	*o.meters.NamedWord("nox_color_white_2523948") = 0x7fff7fff
	*o.meters.NamedWord("nox_color_yellow_2589772") = 0x7fe07fe0
	clear(o.grid)
	clear(o.equipment)
	o.c.srv.NetList.ResetAll()
	*o.meters.NamedWord("nox_win_width") = 256
	*o.meters.NamedWord("nox_win_height") = 256
	*o.meters.NamedWord("dword_5d4594_1090276") = uint32(uintptr(o.parent.C()))
	*o.meters.NamedWord("dword_5d4594_1096288") = uint32(uintptr(o.font))
	for i := range o.meters.Records {
		w := o.c.GUI.NewWindowRaw(o.parent, 8, 12+i*26, 8, 25, 125, nil)
		w.WidgetData = unsafe.Pointer(uintptr(i))
		o.meters.Records[i] = legacy.PortTestMeterRecord{Window: w, Current: 50, Maximum: 100, Color: 0x7c007c00, Alternate: 0x42104210}
	}
	for i := 2240; i < 2256; i++ {
		*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), i)) = byte(i*17 + 3)
	}
	*o.meters.NamedWord("dword_8531A0_2576") = uint32(uintptr(unsafe.Pointer(&o.players[0])))
	o.c.srv.SetFrame(120)
	clear(inputKeyTimeoutsOld)
	o.sounds = nil
	o.text = nil
	clear(o.pix.Pix)
	o.collect()
}
