//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"testing"
	"unsafe"

	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/noxrender"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"golang.org/x/image/font/basicfont"
)

var bookResourceNames = []string{"ArrowNW", "ArrowN", "ArrowNE", "ArrowW", "ArrowE", "ArrowSW", "ArrowS", "ArrowSE", "BookOfKnowledge", "GuideTabLit", "SpellTabLit", "BookPageForward", "BookPageBackward"}

type spellbookOwner struct {
	*inventoryWindowOwner
	words   map[string]*uint32
	region  []uint32
	font    unsafe.Pointer
	refs    []*legacy.ImageRef
	anims   []*legacy.ImageRefAnim
	missing string
}

func newSpellbookOwner(t *testing.T) *spellbookOwner {
	t.Helper()
	o := &spellbookOwner{inventoryWindowOwner: newInventoryWindowOwner(t)}
	t.Cleanup(o.c.srv.Server.PortTestBookSpellOwner())
	var entries []strman.Entry
	for _, pair := range [][2]string{{"Size", "Size"}, {"Small", "Small"}, {"Medium", "Medium"}, {"Large", "Large"}, {"ManaCost", "Mana cost"}, {"PowerLevel", "Power level %d"}, {"SpellInstant", "Instant"}, {"SpellTargeted", "Targeted"}, {"SpellAtLocation", "At location"}, {"SpellHostile", "Hostile"}, {"EmptyBook", "No known entries"}, {"ToolTipAbilityTab", "Abilities"}, {"ToolTipSpellTab", "Spells"}, {"ToolTipGuideTab", "Creatures"}} {
		entries = append(entries, strman.Entry{ID: strman.ID("guibook.c:" + pair[0]), Vals: []strman.Variant{{Str: pair[1]}}})
	}
	configureStrings, restoreStrings := o.c.srv.Server.PortTestMeterStrings(entries...)
	t.Cleanup(restoreStrings)
	configureStrings(0)
	var restore func()
	o.words, restore = legacy.PortTestBookWords()
	t.Cleanup(restore)
	raw := unsafe.Slice(memmap.PtrUint8(0x5D4594, 1046612), 932)
	saved := bytes.Clone(raw)
	t.Cleanup(func() { copy(raw, saved) })
	o.region = unsafe.Slice((*uint32)(unsafe.Pointer(&raw[0])), len(raw)/4)
	vec := *legacy.PortTestBookVector()
	t.Cleanup(func() { *legacy.PortTestBookVector() = vec })
	o.font, restore = o.c.Render().GetFonts().PortTestWindowFont(basicfont.Face7x13)
	t.Cleanup(restore)
	o.c.dataRefs[uint32(uintptr(o.font))] = 0xee100001
	callbacks := legacy.PortTestBookCallbacks()
	var names []string
	for name := range callbacks {
		names = append(names, name)
	}
	sort.Strings(names)
	for i, name := range names {
		o.c.callbackRefs[callbacks[name]] = 0xee200001 + uint32(i)
	}
	for i := 0; i < 2; i++ {
		ref, freeRef := alloc.New(legacy.ImageRef{})
		t.Cleanup(freeRef)
		anim, freeAnim := alloc.New(legacy.ImageRefAnim{})
		t.Cleanup(freeAnim)
		*ref = *o.shiny
		*anim = *o.animation
		ref.Field_24 = unsafe.Pointer(anim)
		ref.SetName(bookResourceNames[11+i])
		o.refs = append(o.refs, ref)
		o.anims = append(o.anims, anim)
		o.c.dataRefs[uint32(uintptr(ref.C()))] = 0xee300001 + uint32(i)
		o.c.dataRefs[uint32(uintptr(anim.C()))] = 0xee400001 + uint32(i)
		o.c.dataRefs[uint32(uintptr(unsafe.Pointer(anim.ImagesPtr)))] = 0xee500001
	}
	oldLoad := legacy.Nox_xxx_gLoadImg
	t.Cleanup(func() { legacy.Nox_xxx_gLoadImg = oldLoad })
	legacy.Nox_xxx_gLoadImg = func(name string) *noxrender.Image {
		for i, n := range bookResourceNames[:11] {
			if name == n {
				o.loads = append(o.loads, name)
				if o.missing == name {
					return nil
				}
				return o.images[i]
			}
		}
		return oldLoad(name)
	}
	oldAnim := legacy.Nox_xxx_gLoadAnim
	t.Cleanup(func() { legacy.Nox_xxx_gLoadAnim = oldAnim })
	legacy.Nox_xxx_gLoadAnim = func(name string) *legacy.ImageRef {
		for i, n := range bookResourceNames[11:] {
			if name == n {
				o.loads = append(o.loads, name)
				if o.missing == name {
					return nil
				}
				return o.refs[i]
			}
		}
		return oldAnim(name)
	}
	// Windows can retain the owned animation callbacks; destroy them before the
	// animation records are freed and before previous globals are restored.
	t.Cleanup(func() { o.c.GUI.DestroyAll(); o.c.GUI.FreeDestroyed(); o.c.GUI.FreeDestroyed() })
	return o
}
func (o *spellbookOwner) resetBook(t *testing.T) {
	o.inventoryWindowOwner.reset(t)
	for _, p := range o.words {
		*p = 0
	}
	clear(o.region)
	*legacy.PortTestBookVector() = [2]float32{}
	o.missing = ""
	for _, anim := range o.anims {
		anim.OnEnd = nil
		anim.Field_3 = 0
	}
	*o.words["nox_win_width"], *o.words["nox_win_height"] = 640, 480
	*o.words["nox_xxx_aNox_cfg_0_587000_132136"] = ^uint32(0)
	*o.words["dword_8531A0_2576"] = uint32(uintptr(unsafe.Pointer(&o.players[0])))
	*o.words["nox_player_netCode_85319C"] = 73
	o.players[0].Active = 1
	o.players[0].NetCodeVal = 73
	clear(unsafe.Slice((*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 3696)), 178))
	*memmap.PtrPtr(0x852978, 8) = nil
	noxflags.UnsetGame(noxflags.GetGame())
	noxflags.SetGame(noxflags.GameHost)
	o.loads = nil
	o.sounds = nil
	clear(o.pix.Pix)
	o.collect()
}
func (o *spellbookOwner) bookCall(name string, a ...uint32) uint32 {
	var args [4]uint32
	copy(args[:], a)
	return legacy.PortTestBookInvoke(name, args)
}
func (o *spellbookOwner) bookWindow() *gui.Window {
	return (*gui.Window)(unsafe.Pointer(uintptr(*o.words["nox_win_unk1"])))
}

type spellbookResult struct {
	Case              string
	Return            uint32
	Words             map[string]uint32
	Region            []uint32
	Windows           [][]uint32
	Animations        [][4]uint32
	Loads             []string
	Sounds            [][2]int
	Vector            [2]float32
	Pixels            string
	Known             []uint32
	Captured, Focused uint32
}

func (o *spellbookOwner) bookSnapshot(label string, ret uint32) spellbookResult {
	o.collect()
	r := spellbookResult{Case: label, Return: o.normalize(ret), Words: make(map[string]uint32), Region: append([]uint32(nil), o.region...), Loads: append([]string(nil), o.loads...), Sounds: append([][2]int(nil), o.sounds...), Vector: *legacy.PortTestBookVector()}
	r.Pixels = effectsPixelHash(o.pix)
	r.Known = append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 3696)), 178)...)
	r.Captured = o.normalize(uint32(uintptr(o.c.GUI.Captured().C())))
	r.Focused = o.normalize(uint32(uintptr(o.c.GUI.Focused().C())))
	pointers := map[string]bool{"nox_win_unk1": true, "dword_8531A0_2576": true, "dword_5d4594_1047516": true, "dword_5d4594_1046924": true, "dword_5d4594_1046928": true, "dword_5d4594_1046944": true, "dword_5d4594_1046948": true, "dword_5d4594_1046952": true, "dword_5d4594_1046956": true}
	for name, p := range o.words {
		v := *p
		if pointers[name] {
			v = o.normalize(v)
		}
		r.Words[name] = v
	}
	for _, off := range []int{1046644, 1046660, 1046856, 1046888, 1046892, 1046896, 1046900, 1046908, 1046912, 1046916, 1046920} {
		i := (off - 1046612) / 4
		r.Region[i] = o.normalize(r.Region[i])
	}
	for _, w := range o.windows {
		words := append([]uint32(nil), unsafe.Slice((*uint32)(w.C()), 101)...)
		for _, i := range []int{13, 15, 17, 19, 21, 23, 59, 93, 94, 95, 96, 97, 98, 99, 100} {
			words[i] = o.normalize(words[i])
		}
		r.Windows = append(r.Windows, words)
	}
	for _, anim := range o.anims {
		words := *(*[4]uint32)(anim.C())
		words[0] = o.normalize(words[0])
		words[1] = o.normalize(words[1])
		r.Animations = append(r.Animations, words)
	}
	return r
}
func spellbookCapture(t *testing.T, label string, rows any, want string) {
	t.Helper()
	raw, err := json.Marshal(rows)
	if err != nil {
		t.Fatal(err)
	}
	if path := os.Getenv("OPENNOX_SPELLBOOK_CAPTURE"); path != "" {
		if err := os.WriteFile(path+"-"+label+".json", raw, 0600); err != nil {
			t.Fatal(err)
		}
	}
	got := fmt.Sprintf("%x", sha256.Sum256(raw))
	t.Log(label, got)
	if want != "" && got != want {
		t.Fatalf("%s changed: %s", label, got)
	}
}
