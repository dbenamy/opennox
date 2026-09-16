//go:build porttest

package opennox

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/libs/things"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestSpellbookAssetPages(t *testing.T) {
	path := os.Getenv("OPENNOX_SPELLBOOK_ASSETS")
	if path == "" {
		t.Skip("set OPENNOX_SPELLBOOK_ASSETS to the original Nox data directory")
	}
	o := newSpellbookOwner(t)
	restore, err := o.c.r.GetBag().PortTestAssetBag(filepath.Join(path, "video.bag"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(restore)
	tf, err := things.Open(filepath.Join(path, "thing.bin"))
	if err != nil {
		t.Fatal(err)
	}
	defer tf.Close()
	images, err := tf.ReadImages()
	if err != nil {
		t.Fatal(err)
	}
	byName := make(map[string]things.Image)
	for _, im := range images {
		byName[im.Name] = im
	}
	oldImages := nox_images_arr1_787156
	t.Cleanup(func() { nox_images_arr1_787156 = oldImages })
	nox_images_arr1_787156 = nil
	load := func(r things.ImageRef) *noxrender.Image {
		if r.Ind < 0 {
			t.Fatalf("book requires external image %q", r.Name)
		}
		im := o.c.r.GetBag().ImageByIndex(r.Ind)
		if im == nil || len(im.Pixdata()) == 0 {
			t.Fatalf("missing original image %d", r.Ind)
		}
		o.c.imageRefs[uint32(uintptr(im.C()))] = 0xee800000 + uint32(r.Ind)
		return im
	}
	for i, name := range bookResourceNames {
		im, ok := byName[name]
		if !ok {
			t.Fatalf("missing original book resource %s", name)
		}
		var ref *legacy.ImageRef
		if i < 11 {
			if im.Img == nil {
				t.Fatal("original still image kind")
			}
			ref, restore = alloc.New(legacy.ImageRef{})
			t.Cleanup(restore)
			ref.SetName(name)
			ref.RefKind = 1
			ref.Field_24 = unsafe.Pointer(uintptr(im.Img.Ind))
			load(*im.Img)
		} else {
			if im.Ani == nil || len(im.Ani.Frames) == 0 {
				t.Fatal("original book animation kind")
			}
			ref = o.refs[i-11]
			anim := o.anims[i-11]
			frames, freeFrames := alloc.Make([]noxrender.ImageHandle{}, len(im.Ani.Frames))
			t.Cleanup(freeFrames)
			for j, r := range im.Ani.Frames {
				frames[j] = load(r).C()
			}
			anim.ImagesPtr = &frames[0]
			anim.ImagesSz = uint8(len(frames))
			anim.Field_2_1 = im.Ani.Field
			anim.AnimType = uint8(im.Ani.Kind)
			o.c.dataRefs[uint32(uintptr(unsafe.Pointer(anim.ImagesPtr)))] = 0xee900001 + uint32(i-11)
		}
		nox_images_arr1_787156 = append(nox_images_arr1_787156, ref)
	}
	oldLoad, oldAnim := legacy.Nox_xxx_gLoadImg, legacy.Nox_xxx_gLoadAnim
	t.Cleanup(func() { legacy.Nox_xxx_gLoadImg, legacy.Nox_xxx_gLoadAnim = oldLoad, oldAnim })
	legacy.Nox_xxx_gLoadImg = func(name string) *noxrender.Image { o.loads = append(o.loads, name); return nox_xxx_gLoadImg(name) }
	legacy.Nox_xxx_gLoadAnim = func(name string) *legacy.ImageRef { o.loads = append(o.loads, name); return nox_xxx_gLoadAnim(name) }
	t.Cleanup(func() { o.c.GUI.DestroyAll(); o.c.GUI.FreeDestroyed(); o.c.GUI.FreeDestroyed() })

	oldAbilities := o.c.srv.abilities.defs
	t.Cleanup(func() { o.c.srv.abilities.defs = oldAbilities })
	o.c.srv.abilities.defs[1] = AbilityDef{name: "Berserker Charge", field24: 1}
	o.c.srv.abilities.defs[2] = AbilityDef{name: "War Cry", field24: 1}
	guides := unsafe.Slice(memmap.PtrUint32(0x5D4594, 740076), 41*7)
	oldGuides := append([]uint32(nil), guides...)
	t.Cleanup(func() { copy(guides, oldGuides) })
	clear(guides)
	for i, name := range []string{"Wolf", "Urchin"} {
		title := tooltipWide(t, titlesToWords(name))
		guides[7*(i+1)] = uint32(uintptr(unsafe.Pointer(&title[0])))
		guides[7*(i+1)+1] = 1
	}
	var rows []spellbookResult
	for pass := 0; pass < 2; pass++ {
		o.resetBook(t)
		if o.bookCall("nox_xxx_bookInit_45B9D0") != 1 {
			t.Fatal("original-asset book initialization")
		}
		o.bookCall("nox_xxx_bookSetColor_45AC40")
		*memmap.PtrUint32(0x5D4594, 1047508) = 2
		*memmap.PtrUint32(0x5D4594, 1046940) = 1
		*memmap.PtrUint32(0x5D4594, 1046960) = 1
		*memmap.PtrUint32(0x5D4594, 1046964) = 2
		for _, view := range []uint32{0, 1, 2, 3} {
			if view < 2 {
				*o.words["dword_5d4594_1046868"] = view
				*o.words["dword_5d4594_1046872"] = view
			} else if view == 2 {
				o.bookCall("nox_xxx_bookMoveToPage_45B930", 0)
			} else {
				o.bookCall("nox_xxx_book_45B210", 0, 5)
			}
			for frame := 0; frame < 256; frame++ {
				clear(o.pix.Pix)
				o.displayText = nil
				blank := effectsPixelHash(o.pix)
				ret := o.bookCall("nox_xxx_bookDrawList_45BD40", *o.words["nox_win_unk1"])
				if ret != 1 || effectsPixelHash(o.pix) == blank {
					t.Fatal("original book draws pixels")
				}
				if view < 2 && len(o.displayText) != 2 {
					t.Fatal("original-asset book must render both entry names")
				}
				rows = append(rows, o.bookSnapshot(fmt.Sprintf("pass%d-view%d-frame%d", pass, view, frame), ret))
				if view < 2 || *o.words["dword_5d4594_1046868"] != view {
					break
				}
				if frame == 255 {
					t.Fatal("original page animation failed to complete")
				}
				o.c.Inp.Tick()
			}
		}
	}
	spellbookCapture(t, "asset-pages", rows, "aab7dc659cb46b961797c2661d6c6279761120cdc35582f4415584b3a5a29cdc")
}
