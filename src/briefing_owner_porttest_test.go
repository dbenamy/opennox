//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noximage"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/music"
	"golang.org/x/image/font/basicfont"
)

var briefingSpriteNames = []string{"GauntletExitB", "BeholderGenerator", "Ankh", "SoulGate", "SilverKey", "GoldKey", "QuestGoldChest", "QuestGoldPile", "DunMirChest4", "WarHammer", "HastePotion", "ConjurerSpellBook"}
var briefingSpriteWords = []string{"dword_5d4594_832496", "dword_5d4594_832492", "dword_5d4594_832500", "dword_5d4594_832504", "dword_5d4594_832508", "dword_5d4594_832512", "dword_5d4594_832516", "dword_5d4594_832520", "dword_5d4594_832524", "dword_5d4594_832528", "dword_5d4594_832532", "dword_5d4594_832536"}

const briefingWindowResource = `FONT = small; WINDOW 1000 0 0 640 480 USER; STATUS = ENABLED+HIDDEN; CHILD WINDOW 1010 0 0 640 480 STATICTEXT; STATUS = ENABLED; DATA = 0 0 WindowDir:Blank; END END END`

type briefingOwner struct {
	*inventoryWindowOwner
	briefWords  map[string]*uint32
	regions     [][]byte
	strings     map[string]strman.Variant
	defaultFont unsafe.Pointer
}
type briefingResult struct {
	Op       int
	Return   uint32
	Named    map[string]uint32
	Regions  [][]byte
	Text     []inventoryDisplayText
	Loads    []string
	Render   objectRenderResult
	Music    music.MusicState
	Fade     [5]bool
	Windows  []inventoryWindowWidget
	Captured uint32
}

func newBriefingOwner(t *testing.T) *briefingOwner {
	// InventoryWindowOwner already registers the two quest gold sprites.
	var extra []string
	for _, name := range briefingSpriteNames {
		if name != "QuestGoldChest" && name != "QuestGoldPile" {
			extra = append(extra, name)
		}
	}
	o := &briefingOwner{inventoryWindowOwner: newInventoryWindowOwner(t, extra...)}
	for i, name := range briefingSpriteNames {
		typ := o.c.Things.TypeByID(name)
		data, free := alloc.Make([]uint32{}, 2)
		t.Cleanup(free)
		data[0], data[1] = 8, uint32(uintptr(o.images[i%len(o.images)].C()))
		typ.DrawFunc = legacy.PortTestSpriteAnimationCallback(2)
		typ.DrawData = unsafe.Pointer(&data[0])
		typ.ObjClass, typ.ObjFlags = 0x10, 0x40000000
		o.c.dataRefs[uint32(txptr(typ.DrawData))] = 0xf0300000 + uint32(i)
	}
	var restore func()
	o.briefWords, restore = legacy.PortTestBriefingWords()
	t.Cleanup(restore)
	for _, r := range [][3]uintptr{{0x5D4594, 831228, 4}, {0x5D4594, 831248, 8}, {0x5D4594, 831264, 12}, {0x5D4594, 831280, 20}, {0x5D4594, 831300, 1056}, {0x5D4594, 832356, 4}, {0x5D4594, 832364, 112}, {0x5D4594, 832488, 4}, {0x5D4594, 832540, 12}, {0x587000, 122944, 12}, {0x587000, 122960, 48}, {0x587000, 123012, 4}, {0x5D4594, 806004, 40}} {
		b := unsafe.Slice((*byte)(memmap.PtrOff(r[0], r[1])), r[2])
		old := append([]byte(nil), b...)
		o.regions = append(o.regions, b)
		t.Cleanup(func() { copy(b, old) })
	}
	for i, p := range legacy.PortTestBriefingCallbacks() {
		o.c.callbackRefs[p] = 0xf0100000 + uint32(i)
	}
	o.defaultFont, restore = o.c.Render().GetFonts().PortTestWindowFont(basicfont.Face7x13, "default")
	t.Cleanup(restore)
	o.c.dataRefs[uint32(txptr(o.defaultFont))] = 0xf0200000
	oldTicks := legacy.PlatformTicks
	legacy.PlatformTicks = func() uint64 { return 123456 }
	t.Cleanup(func() { legacy.PlatformTicks = oldTicks })
	oldMusic := legacy.MusicModule
	legacy.MusicModule = &music.Module{}
	t.Cleanup(func() { legacy.MusicModule = oldMusic })
	o.strings = map[string]strman.Variant{}
	add := func(id, text, voice string) { o.strings[id] = strman.Variant{Str: text, Str2: voice} }
	for _, cl := range []string{"Warrior", "Wizard", "Conjurer"} {
		for chapter := 1; chapter <= 11; chapter++ {
			for _, kind := range []string{"Begin", "Loss"} {
				id := fmt.Sprintf("Briefing:%sChapter%s%d", cl, kind, chapter)
				add(id, fmt.Sprintf("%s %s chapter %d.", cl, kind, chapter), fmt.Sprintf("voice_%s_%s_%d", cl, kind, chapter))
			}
		}
	}
	add("Nox:Credits", "Fixture credits.", "credits_voice")
	for _, pair := range [][2]string{{"Noxworld.c:Stage", "Stage"}, {"GUIBrief.c:GauntletStatTitle", "Quest results"}, {"GUIBrief.c:GeneratorsDestroyed", "Generators"}, {"GUIBrief.c:numSecretsFound", "Secrets"}, {"GUIBrief.c:Kills", "Kills"}, {"GUIBrief.c:TotalScore", "Score"}, {"GeneralPrint:SecretsTotal", "Total secrets: %d"}, {"GeneralPrint:SecretsFound", "You found %d."}, {"GeneralPrint:SecretsNoneFound", "You found none."}, {"GeneralPrint:SecretsFoundByFriends", "Others found %d."}, {"GeneralPrint:SecretsNoneFoundByFriends", "Others found none."}, {"GeneralPrint:QuestSplash1", "Quest instructions"}, {"GeneralPrint:QuestSplash12", "Press a key to continue"}, {"Briefing:Custom", "A custom briefing."}} {
		add(pair[0], pair[1], "")
	}
	for i := 2; i <= 11; i++ {
		for _, suffix := range []string{"a", "b"} {
			add(fmt.Sprintf("GeneralPrint:QuestSplash%d%s", i, suffix), fmt.Sprintf("Instruction %d%s", i, suffix), "")
		}
	}
	o.installStrings(t)
	o.c.dataRefs[uint32(txptr(unsafe.Pointer(alloc.InternCString16(""))))] = 0xf0400000
	for i, off := range []uintptr{832540, 832544, 832548} {
		o.c.dataRefs[uint32(txptr(memmap.PtrOff(0x5D4594, off)))] = 0xf0400010 + uint32(i)
	}
	oldLoad := legacy.Nox_xxx_gLoadImg
	legacy.Nox_xxx_gLoadImg = func(name string) *noxrender.Image {
		if name == "MissingBriefingImage" {
			o.loads = append(o.loads, name)
			return nil
		}
		if strings.Contains(name, "ChapterBegin") || strings.Contains(name, "ChapterLoss") || name == "CreditsImage" || name == "CopyrightScreen" || name == "GauntletStartMines" || name == "MenuSystemBG" || name == "GauntletInstructionBackground" || name == "CustomBriefingImage" {
			o.loads = append(o.loads, name)
			h := sha256.Sum256([]byte(name))
			return o.images[int(h[0])%len(o.images)]
		}
		return oldLoad(name)
	}
	t.Cleanup(func() { legacy.Nox_xxx_gLoadImg = oldLoad })
	oldResource := legacy.Nox_new_window_from_file
	legacy.Nox_new_window_from_file = func(name string, fn gui.WindowFunc) *gui.Window {
		if name == "Briefing.wnd" {
			return newWindowFromString(o.c.GUI, briefingWindowResource, fn)
		}
		return oldResource(name, fn)
	}
	t.Cleanup(func() { legacy.Nox_new_window_from_file = oldResource })
	t.Cleanup(o.releaseBriefing)
	o.resetBriefing(t)
	return o
}
func (o *briefingOwner) installStrings(t *testing.T) {
	path := filepath.Join(t.TempDir(), "briefing-strings.json")
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
	for id, v := range o.strings {
		found := false
		for i := range data.Entries {
			if string(data.Entries[i].ID) == id {
				data.Entries[i].Vals = []strman.Variant{v}
				found = true
				break
			}
		}
		if !found {
			data.Entries = append(data.Entries, strman.Entry{ID: strman.ID(id), Vals: []strman.Variant{v}})
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
	// Normalize owned strings by stable content, not iteration/allocation order.
	for id, v := range o.strings {
		for kind, s := range []string{v.Str, v.Str2} {
			h := sha256.Sum256([]byte(fmt.Sprintf("%d:%s", kind, s)))
			ref := binary.LittleEndian.Uint32(h[:]) | 0x80000000
			var p unsafe.Pointer
			if kind == 0 {
				p = unsafe.Pointer(alloc.InternCString16(s))
			} else {
				p = unsafe.Pointer(alloc.InternCString(s))
			}
			o.c.dataRefs[uint32(txptr(p))] = ref
		}
		_ = id
	}
}
func (o *briefingOwner) releaseBriefing() {
	if o.briefWords != nil {
		legacy.PortTestBriefing(14, 0, 0, 0)
		o.c.GUI.FreeDestroyed()
	}
}
func (o *briefingOwner) resetBriefing(t *testing.T) {
	o.releaseBriefing()
	o.inventoryWindowOwner.reset(t)
	for _, p := range o.briefWords {
		*p = 0
	}
	for _, b := range o.regions {
		clear(b)
	}
	legacy.MusicModule = &music.Module{}
	o.c.r.FadeReset()
	for i, cl := range []string{"Warrior", "Wizard", "Conjurer"} {
		*memmap.PtrPtr(0x587000, 122944+uintptr(4*i)) = unsafe.Pointer(alloc.InternCString(cl))
	}
	for i, pos := range [][2]int{{30, 80}, {330, 80}, {30, 210}, {330, 210}, {30, 340}, {330, 340}} {
		*memmap.PtrUint32(0x587000, 122960+uintptr(i*8)) = uint32(pos[0])
		*memmap.PtrUint32(0x587000, 122964+uintptr(i*8)) = uint32(pos[1])
	}
	for i := 0; i < 6; i++ {
		p := &o.players[i]
		p.Active = 1
		p.NetCodeVal = uint32(100 + i)
		p.SetName(fmt.Sprintf("Player%d", i))
	}
	o.resize(640, 480)
	o.loads = nil
	o.displayText = nil
}
func (o *briefingOwner) resize(w, h int) {
	nox_win_width, nox_win_height = w, h
	o.pix = noximage.NewImage16(image.Rect(0, 0, w, h))
	o.c.r.SetPixBuffer(o.pix)
	o.c.r.Data().SetClipRect(o.pix.Rect)
	o.c.r.Data().SetClipRect2(image.Rect(0, 0, w-1, h-1))
	o.c.r.Data().SetRect3(o.pix.Rect)
	*o.c.Viewport() = noxrender.Viewport{Screen: o.pix.Rect, World: o.pix.Rect, Size: o.pix.Rect.Size()}
}
func (o *briefingOwner) constructBriefing(t *testing.T) {
	if legacy.PortTestBriefing(13, 0, 0, 0) == 0 {
		t.Fatal("briefing owner construction")
	}
	o.collect()
	if o.briefWindow().ChildByID(1010) == nil {
		t.Fatal("briefing text child")
	}
}
func (o *briefingOwner) briefWindow() *gui.Window {
	return (*gui.Window)(unsafe.Pointer(uintptr(*o.briefWords["nox_wnd_briefing_831232"])))
}
func (o *briefingOwner) snapshotBriefing(t *testing.T, op int, ret uint32) briefingResult {
	o.collect()
	ws := o.capture(t, 0, 0, 7, 0, 0, 0, 0)
	r := briefingResult{Windows: ws.Windows, Captured: ws.Captured, Op: op, Return: o.norm(ret), Named: map[string]uint32{}, Text: append([]inventoryDisplayText(nil), o.displayText...), Loads: append([]string(nil), o.loads...), Render: ws.Display.Render, Music: legacy.MusicModule.GetCurrentBlock()}
	for n, p := range o.briefWords {
		r.Named[n] = o.norm(*p)
	}
	for i, b := range o.regions {
		if i == 9 {
			continue
		} // Static class-name input pointers are independently checked via all generated entries.
		cp := append([]byte(nil), b...)
		switch i {
		case 2, 4:
			for off := 0; off < len(cp); off += 4 {
				binary.LittleEndian.PutUint32(cp[off:], o.norm(binary.LittleEndian.Uint32(cp[off:])))
			}
		case 6:
			for off := 0; off < 96; off += 16 {
				binary.LittleEndian.PutUint32(cp[off:], o.norm(binary.LittleEndian.Uint32(cp[off:])))
			}
			for _, off := range []int{96, 100} {
				binary.LittleEndian.PutUint32(cp[off:], o.norm(binary.LittleEndian.Uint32(cp[off:])))
			}
		}
		r.Regions = append(r.Regions, cp)
	}
	for i := range r.Fade {
		r.Fade[i] = o.c.r.CheckFade(noxrender.FadeKey(i))
	}
	return r
}
func briefingCapture(t *testing.T, label string, rows []briefingResult, want string) {
	t.Helper()
	raw, err := json.Marshal(rows)
	if err != nil {
		t.Fatal(err)
	}
	if prefix := os.Getenv("OPENNOX_BRIEFING_CAPTURE"); prefix != "" {
		if err = os.WriteFile(prefix+"-"+label+".json", raw, 0600); err != nil {
			t.Fatal(err)
		}
	}
	got := fmt.Sprintf("%x", sha256.Sum256(raw))
	t.Logf("%s: %d results %s", label, len(rows), got)
	if got != want {
		t.Fatalf("%s hash %s want %s", label, got, want)
	}
}
