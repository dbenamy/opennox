//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"image"
	"os"
	"sort"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

type inventoryDisplayText struct {
	Kind, Text string
	Rect       image.Rectangle
	Return     int
}
type inventoryDisplayRender struct {
	*meterRender
	owner *inventoryDisplayOwner
}

func (r *inventoryDisplayRender) DrawString(face font.Face, s string, pos image.Point) int {
	v := r.meterRender.DrawString(face, s, pos)
	r.owner.displayText = append(r.owner.displayText, inventoryDisplayText{"plain", s, image.Rectangle{Min: pos}, v})
	return v
}
func (r *inventoryDisplayRender) DrawStringWrapped(face font.Face, s string, rect image.Rectangle) int {
	v := r.meterRender.DrawStringWrapped(face, s, rect)
	r.owner.displayText = append(r.owner.displayText, inventoryDisplayText{"wrapped", s, rect, v})
	return v
}

type inventoryDisplayClient struct {
	*inventoryTransactionClient
	owner *inventoryDisplayOwner
}

func (c *inventoryDisplayClient) R2() legacy.Render2 {
	return &inventoryDisplayRender{&meterRender{&objectDrawingRender{c.owner.c.r, c.owner.objectDrawingOwner}, c.owner.meterOwner}, c.owner}
}

type inventoryDisplayOwner struct {
	*inventoryTransactionOwner
	displayWords   map[string]*uint32
	displayRegions [][]byte
	displayText    []inventoryDisplayText
	smallFont      unsafe.Pointer
	cursorText     []uint16
	language       func(int)
}

func newInventoryDisplayOwner(t *testing.T, extraNames ...string) *inventoryDisplayOwner {
	o := &inventoryDisplayOwner{inventoryTransactionOwner: newInventoryTransactionOwner(t, append([]string{"ArcherArrow", "ArcherBolt"}, extraNames...)...)}
	var restore func()
	o.displayWords, restore = legacy.PortTestInventoryDisplayWords()
	t.Cleanup(restore)
	o.cursorText = tooltipMapped(t, 1096676, 256)
	for i, name := range []string{"Warrior", "Wizard", "Conjurer"} {
		p := memmap.PtrPtr(0x587000, 29456+uintptr(i*4))
		old := *p
		*p = unsafe.Pointer(alloc.InternCString(name))
		t.Cleanup(func() { *p = old })
	}
	for i := 1; o.c.Things.TypeByInd(i) != nil; i++ {
		typ := o.c.Things.TypeByInd(i)
		typ.PrettyName = alloc.InternCString16(fmt.Sprintf("Item%d", i))
	}
	o.weapon.Desc8 = alloc.InternCString16("Bow")
	o.armor.Desc8 = alloc.InternCString16("Armor")
	oldClient := legacy.GetClient
	proxy := &inventoryDisplayClient{o.tx, o}
	legacy.GetClient = func() legacy.Client { return proxy }
	t.Cleanup(func() { legacy.GetClient = oldClient })
	for _, r := range [][3]uintptr{{0x5D4594, 1049868, 4}, {0x5D4594, 1049908, 92}, {0x5D4594, 1050012, 4}, {0x5D4594, 1062588, 512}, {0x5D4594, 1063124, 512}, {0x5D4594, 1063644, 128}, {0x85B3FC, 940, 20}, {0x5D4594, 2597996, 4}, {0x587000, 136192, 224}, {0x5D4594, 1568264, 4}} {
		b := unsafe.Slice((*byte)(memmap.PtrOff(r[0], r[1])), r[2])
		old := append([]byte(nil), b...)
		o.displayRegions = append(o.displayRegions, b)
		t.Cleanup(func() { copy(b, old) })
	}
	small := *basicfont.Face7x13
	small.Advance = 6
	small.Width = 6
	small.Height = 11
	small.Ascent = 9
	o.smallFont, restore = o.c.Render().GetFonts().PortTestWindowFont(&small, "display-small")
	t.Cleanup(restore)
	o.c.dataRefs[uint32(uintptr(o.smallFont))] = 0xe3000002
	t.Cleanup(func() { delete(o.c.dataRefs, uint32(uintptr(o.smallFont))) })
	oldStats := o.c.srv.Players.Stats
	t.Cleanup(func() { o.c.srv.Players.Stats = oldStats })
	entries := []strman.Entry{{ID: "guimsg.c:systemmsg", Vals: []strman.Variant{{Str: "System: %s"}}}}
	add := func(id, text string) {
		if !strings.Contains(id, ":") {
			id = "guiinv.c:" + id
		}
		entries = append(entries, strman.Entry{ID: strman.ID(id), Vals: []strman.Variant{{Str: text}}})
	}
	for _, pair := range [][2]string{{"WindowDir:Blank", ""}, {"thing.db:AnkhGUI", "Extra lives"}, {"GeneralPrint:TooltipKeyIcon", "Shared keys"}, {"OpenInventoryTT", "Open inventory"}, {"ToolTipWeapon2Area", "Secondary weapon"}, {"ObjectTooFar", "Object too far"}, {"NoObject", "No object"}, {"DropLabel", "Drop amount"}, {"Journal:QuestLabel", "Quest:"}, {"Journal:CompletedLabel", "Complete:"}, {"Journal:HintLabel", "Hint:"}, {"Journal:PortInventoryEntry", "A fixture journal entry with enough text to exercise the real wrapped text renderer across several lines."}, {"StatsLevel", "Level %d"}, {"StatsEXP", "Experience %d / %d"}, {"StatsHealth", "Health"}, {"StatsMana", "Mana"}, {"StatsStrength", "Strength"}, {"StatsSpeed", "Speed"}, {"StatsArmor", "Armor"}, {"MinMaxFormat", "%d / %d"}, {"IdentifyItem", "Identify"}, {"IdentifyWeight", "Weight %d"}, {"IdentifyDurability", "Durability %d / %d"}, {"IdentifyDurabilityIndestructable", "Indestructible"}, {"IdentifyDurabilityNoDamage", "No damage"}, {"IdentifyDurabilitySlight", "Slight damage"}, {"IdentifyDurabilityModerate", "Moderate damage"}, {"IdentifyDurabilitySevere", "Severe damage"}, {"JournalModeTT", "Journal"}, {"InventoryModeTT", "Inventory"}, {"StatsModeTT", "Statistics"}, {"PaperDollModeTT", "Paper doll"}, {"CloseInventoryTT", "Close"}, {"Weapon2CantUse", "Cannot use secondary weapon"}, {"ElaborateNameFormat", "%s the %s"}} {
		add(pair[0], pair[1])
	}
	for _, pair := range [][2]string{{"StatsArmorLabel", "Armor"}, {"DollWeight", "Weight"}, {"DollRegionError", "No equipment region"}, {"ToolTipDrag", "Drag equipment here"}, {"WeaponDamageLabelNA", "Damage unavailable"}, {"WeaponDamageLabel", "Damage %.2f"}, {"WeaponDamageLabelUnknownPlus", "Damage %.2f + unknown"}, {"BaseDamageLabel", "Base %.2f"}, {"StrengthDamageLabel", "Strength %.2f"}, {"FireDamageLabel", "Fire %.2f"}, {"ElectricalDamageLabel", "Lightning %.2f"}, {"ArmorValueLabelNA", "Armor unavailable"}, {"ArmorValueLabel", "Armor %d"}, {"IdentifySpecialAttributes", "Special attributes"}, {"IdentifyUnknown", "Unknown"}, {"thing.db:IdentifyDescription", "Choose an item"}} {
		add(pair[0], pair[1])
	}
	for class, name := range []string{"Warrior", "Wizard", "Conjurer"} {
		for _, level := range []int{-128, -1, 0, 1, 5, 10, 11, 127} {
			add(fmt.Sprintf("experience:%s%d", name, level), fmt.Sprintf("Rank%d-%d", class, level))
		}
	}
	for i := 0; i < 6; i++ {
		add(fmt.Sprintf("Modifier.c:WindowEffect%d", i), fmt.Sprintf("Status effect %d", i))
	}
	t.Cleanup(o.c.srv.Server.PortTestInventoryDisplayBalance())
	configure, restore := o.c.srv.Server.PortTestMeterStrings(entries...)
	t.Cleanup(restore)
	o.language = configure
	configure(0)
	oldStrings := strMan
	strMan = o.c.srv.Strings()
	t.Cleanup(func() { strMan = oldStrings })
	return o
}
func (o *inventoryDisplayOwner) reset(t *testing.T) {
	o.renderEnv.Reset()
	o.inventoryTransactionOwner.reset(t)
	for _, p := range o.displayWords {
		*p = 0
	}
	for _, b := range o.displayRegions {
		clear(b)
	}
	palette := map[string]uint32{"nox_color_white_2523948": 0x7fff7fff, "nox_color_black_2650656": 0, "nox_color_red_2589776": 0x7c007c00, "nox_color_blue_2650684": 0x001f001f, "nox_color_cyan_2649820": 0x03ff03ff, "nox_color_orange_2614256": 0x7e007e00, "nox_color_yellow_2589772": 0x7fe07fe0, "nox_color_violet_2598268": 0x7c1f7c1f}
	for n, v := range palette {
		*o.displayWords[n] = v
	}
	*o.displayWords["dword_8531A0_2576"] = uint32(uintptr(unsafe.Pointer(&o.players[0])))
	*o.displayWords["dword_5d4594_1063636"] = uint32(uintptr(o.smallFont))
	*o.displayWords["dword_5d4594_1062456"] = uint32(uintptr(o.parent.C()))
	*o.displayWords["dword_5d4594_1062476"] = uint32(uintptr(o.parent.C()))
	o.c.srv.Players.Stats.Base = server.ClassStats{Health: 100, Mana: 80, Speed: 100, Strength: 60}
	o.c.srv.Players.Stats.Warrior = server.ClassStats{Health: 150, Mana: 0, Speed: 100, Strength: 80}
	o.c.srv.Players.Stats.Wizard = server.ClassStats{Health: 75, Mana: 150, Speed: 110, Strength: 40}
	o.c.srv.Players.Stats.Conjurer = server.ClassStats{Health: 100, Mana: 100, Speed: 105, Strength: 60}
	for _, off := range []uintptr{2235, 2239, 2243, 2247} {
		*(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), off)) = 50
	}
	*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 2251)) = 0
	*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 3684)) = 1
	*memmap.PtrUint32(0x85B3FC, 940) = 0x03e003e0
	*memmap.PtrUint32(0x85B3FC, 944) = 0x42104210
	*memmap.PtrUint32(0x85B3FC, 956) = 0x42104210
	*memmap.PtrUint32(0x5D4594, 2597996) = 0x03e003e0
	o.c.Mouse = image.Pt(10, 20)
	o.c.MouseReads = 0
	o.displayText = nil
	o.drawTrace = nil
	o.rawDeleted = nil
	o.namedCalls = nil
	clear(o.cursorText)
	for _, m := range o.mods {
		m.Attack40 = server.ModifierEffFnc{}
		m.AttackPreHit52 = server.ModifierEffFnc{}
		m.Defend76 = server.ModifierEffFnc{}
	}
}

type inventoryDisplayResult struct {
	Case, Op     int
	Return       uint64
	Text         []inventoryDisplayText
	Pixels       string
	Named        [][2]uint32
	Regions      [][]byte
	Render       objectRenderResult
	CursorText   string
	ReturnedText string
	WidgetText   []string
	Inventory    *inventoryTransactionResult `json:",omitempty"`
}

func (o *inventoryDisplayOwner) call(t *testing.T, id, op int, a, b, c uintptr) inventoryDisplayResult {
	ret := legacy.PortTestInventoryDisplay(op, a, b, c)
	return o.displaySnapshot(t, id, op, ret)
}

func (o *inventoryDisplayOwner) displaySnapshot(t *testing.T, id, op int, ret uint64) inventoryDisplayResult {
	var returnedText string
	if op == 11 && ret != 0 {
		returnedText = alloc.GoString16((*uint16)(unsafe.Pointer(uintptr(ret))))
		// The returned text may live in the string manager or tooltip scratch.
		// Record its content and non-null status, not its allocation address.
		ret = 1
	}
	// snapshot applies its own normalization; keep its input unnormalized.
	snapshotRet := uint32(ret)
	if op != 1 && op != 2 {
		ret = uint64(o.norm(uint32(ret)))
	}
	r := inventoryDisplayResult{Case: id, Op: op, Return: ret, Text: append([]inventoryDisplayText(nil), o.displayText...), Pixels: effectsPixelHash(o.pix)}
	r.CursorText = alloc.GoString16(&o.cursorText[0])
	r.ReturnedText = returnedText
	if op == 3 || op == 11 || op == 15 {
		state := o.snapshot(t, id, op, snapshotRet)
		// These broad transaction scratch regions overlap display pointers.
		// Display state is captured by its declared fields below instead.
		state.Regions = nil
		r.Inventory = &state
	}
	names := make([]string, 0, len(o.displayWords))
	for n := range o.displayWords {
		names = append(names, n)
	}
	sort.Strings(names)
	for i, n := range names {
		r.Named = append(r.Named, [2]uint32{uint32(i), o.norm(*o.displayWords[n])})
	}
	for index, b := range o.displayRegions {
		cp := append([]byte(nil), b...)
		if index == 1 { // Inventory image handles, all 23 declared words.
			for off := 0; off < len(cp); off += 4 {
				binary.LittleEndian.PutUint32(cp[off:], o.norm(binary.LittleEndian.Uint32(cp[off:])))
			}
		}
		r.Regions = append(r.Regions, cp)
	}
	r.Render = o.renderResult(t, id, 0, 0)
	// Display items use the actual static-image drawable callback, whose render
	// fixture does not interpret item modifier slots. These four declared item
	// pointers still need identities when recording display-specific metadata.
	for _, row := range r.Render.Draw.Drawables {
		owned := false
		for _, obj := range o.objects {
			owned = owned || obj.Live && row[0] == obj.Ref
		}
		if !owned {
			continue
		}
		for slot := 109; slot <= 112; slot++ {
			v := row[slot]
			if n, ok := o.modRefs[v]; ok {
				row[slot] = n
			} else if v != 0 && (v < 0xeb000000 || v >= 0xeb000004) {
				t.Fatalf("unowned display modifier at word%d: %x", slot-1, v)
			}
		}
	}
	return r
}
func inventoryDisplayCapture(t *testing.T, label string, rows []inventoryDisplayResult, want string) {
	t.Helper()
	b, err := json.Marshal(rows)
	if err != nil {
		t.Fatal(err)
	}
	if prefix := os.Getenv("OPENNOX_INVENTORY_DISPLAY_CAPTURE"); prefix != "" {
		if err := os.WriteFile(prefix+"-"+label+".json", b, 0600); err != nil {
			t.Fatal(err)
		}
	}
	h := fmt.Sprintf("%x", sha256.Sum256(b))
	t.Logf("%s: %d results %s", label, len(rows), h)
	if h != want {
		// Preserve intermittent differences even when an explicit capture was not requested.
		if os.Getenv("OPENNOX_INVENTORY_DISPLAY_CAPTURE") == "" {
			const dir = "../build/port-failures"
			if err := os.MkdirAll(dir, 0700); err != nil {
				t.Logf("cannot create mismatch artifact directory: %v", err)
			} else if f, err := os.CreateTemp(dir, "inventory-display-"+label+"-*.json"); err != nil {
				t.Logf("cannot create mismatch artifact: %v", err)
			} else {
				_, writeErr := f.Write(b)
				closeErr := f.Close()
				t.Logf("full mismatch capture: %s (write error: %v; close error: %v)", f.Name(), writeErr, closeErr)
			}
		}
		t.Fatalf("%s hash %s want frozen C %s", label, h, want)
	}
}

// Identity tokens are comparison values, not raw addresses. An image handle can
// happen to have the same word value in a 32-bit handle arena.
func TestClientInventoryDisplaySnapshotIdentityCollision(t *testing.T) {
	o := newInventoryDisplayOwner(t)
	o.reset(t)
	dr := o.item(t, "RedApple", 123)
	raw := uint32(uintptr(dr.C()))
	want := o.norm(raw)
	if want != 0xec000001 {
		t.Fatalf("first drawable identity %#x", want)
	}
	// Reproduce the observed collision without depending on mmap placement.
	o.c.imageRefs[want] = 0xe8000000
	if got := o.norm(want); got != 0xe8000000 {
		t.Fatalf("identity collision not established: %#x", got)
	}
	for _, op := range []int{3, 15} {
		r := o.displaySnapshot(t, op, op, uint64(raw))
		if r.Inventory == nil {
			t.Fatal("missing inventory snapshot")
		}
		if r.Return != uint64(want) || r.Inventory.Return != want {
			t.Fatalf("op %d: display return %#x, inventory return %#x; want identity %#x in both", op, r.Return, r.Inventory.Return, want)
		}
	}
	text, free := alloc.CString16("description")
	defer free()
	for _, v := range []uintptr{0, uintptr(unsafe.Pointer(text))} {
		r := o.displaySnapshot(t, 11, 11, uint64(v))
		if r.Inventory == nil {
			t.Fatal("missing text inventory snapshot")
		}
		wantReturn, wantText := uint32(0), ""
		if v != 0 {
			wantReturn, wantText = 1, "description"
		}
		if r.Return != uint64(wantReturn) || r.Inventory.Return != wantReturn || r.ReturnedText != wantText {
			t.Fatalf("text return normalization: display %#x inventory return %#x text %q", r.Return, r.Inventory.Return, r.ReturnedText)
		}
	}
}
