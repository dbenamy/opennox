//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestBriefingChapterTable(t *testing.T) {
	o := newBriefingOwner(t)
	ret := legacy.PortTestBriefing(0, 0, 0, 0)
	if alloc.GoString16((*uint16)(unsafe.Pointer(uintptr(ret)))) != "Fixture credits." {
		t.Fatal("credits return")
	}
	if len(o.loads) != 67 {
		t.Fatalf("chapter/credits images requested %d want67", len(o.loads))
	}
	for ci, cl := range []string{"Warrior", "Wizard", "Conjurer"} {
		for ch := 0; ch < 11; ch++ {
			row := unsafe.Slice((*uint32)(memmap.PtrOff(0x5D4594, 831300+uintptr((ci*11+ch)*32))), 8)
			duration := uint32(ch + 3)
			if ch == 0 {
				duration = 1
			}
			if ch == 1 {
				duration = []uint32{4, 2, 3}[ci]
			}
			for ki, kind := range []string{"Begin", "Loss"} {
				off := ki * 4
				want := fmt.Sprintf("%s %s chapter %d.", cl, kind, ch+1)
				voice := fmt.Sprintf("voice_%s_%s_%d", cl, kind, ch+1)
				if row[off] == 0 || alloc.GoString16((*uint16)(unsafe.Pointer(uintptr(row[off+1])))) != want || alloc.GoString((*byte)(unsafe.Pointer(uintptr(row[off+2])))) != voice || row[off+3] != duration {
					t.Fatalf("chapter table class%s chapter%d %s", cl, ch+1, kind)
				}
				if o.loads[ch*6+ci*2+ki] != fmt.Sprintf("%sChapter%s%d", cl, kind, ch+1) {
					t.Fatal("chapter image request order")
				}
			}
		}
	}
	briefingCapture(t, "chapter-table", []briefingResult{o.snapshotBriefing(t, 0, ret)}, "5d224a483890f3cea8ffdf54b224c1c0d784e3edc31ca85975b10de8203d3a68")
}
func TestBriefingStageCaptionAndImage(t *testing.T) {
	o := newBriefingOwner(t)
	var rows []briefingResult
	for _, v := range []uint32{0, 1, 65535, 0x7fffffff, 0x80000000, 0xffffffff} {
		o.resetBriefing(t)
		legacy.PortTestBriefing(11, uintptr(v), 0, 0)
		ret := legacy.PortTestBriefing(12, 0, 0, 0)
		if ret != v || memmap.Uint32(0x5D4594, 832468) != v {
			t.Fatal("stage round trip")
		}
		rows = append(rows, o.snapshotBriefing(t, 12, ret))
	}
	for _, name := range []string{"", "A custom briefing."} {
		o.resetBriefing(t)
		p := alloc.InternCString16(name)
		ret := legacy.PortTestBriefing(10, txptr(unsafe.Pointer(p)), 0, 0)
		if ret != uint32(txptr(unsafe.Pointer(p))) || memmap.Uint32(0x5D4594, 832464) != ret {
			t.Fatal("caption owner")
		}
		rows = append(rows, o.snapshotBriefing(t, 10, ret))
	}
	for _, present := range []bool{false, true} {
		o.resetBriefing(t)
		var p uintptr
		if present {
			p = uintptr(o.images[0].C())
		}
		ret := legacy.PortTestBriefing(9, p, 0, 0)
		if ret == 0 || memmap.Uint32(0x5D4594, 832460) != ret {
			t.Fatal("briefing image")
		}
		if present && ret != uint32(p) {
			t.Fatal("replaced supplied image")
		}
		if !present && (len(o.loads) != 1 || o.loads[0] != "WarriorChapterBegin8") {
			t.Fatal("missing image fallback")
		}
		rows = append(rows, o.snapshotBriefing(t, 9, ret))
	}
	briefingCapture(t, "stage-caption-image", rows, "3fdbe48d8581876862202e9a81fd41d9f3a31e3618afdc7fbf6bda94beb506ce")
}
func TestBriefingScoreComparison(t *testing.T) {
	o := newBriefingOwner(t)
	pair, free := alloc.New([8]uint32{})
	defer free()
	var rows []briefingResult
	for _, a := range []uint32{0, 1, 65535, 0x7fffffff, 0x80000000, 0xffffffff} {
		for _, b := range []uint32{0, 1, 65535, 0x7fffffff, 0x80000000, 0xffffffff} {
			pair[3], pair[7] = a, b
			ret := legacy.PortTestBriefing(8, txptr(unsafe.Pointer(&pair[0])), txptr(unsafe.Pointer(&pair[4])), 0)
			want := int32(0)
			if a > b {
				want = -1
			} else if a < b {
				want = 1
			}
			if int32(ret) != want {
				t.Fatalf("unsigned score order %x %x", a, b)
			}
			rows = append(rows, briefingResult{Op: 8, Return: ret})
		}
	}
	_ = o
	briefingCapture(t, "score-comparison", rows, "311e6d15c549a800678406c7cca805312e7b8e668917891ab178eb8813ec2b72")
}
func TestBriefingSpriteCacheLifecycle(t *testing.T) {
	o := newBriefingOwner(t)
	var rows []briefingResult
	for cycle := 0; cycle < 12; cycle++ {
		o.resetBriefing(t)
		ret := legacy.PortTestBriefing(7, 0, 0, 0)
		if ret != *o.briefWords["dword_5d4594_832536"] {
			t.Fatal("sprite initializer return")
		}
		for i, n := range briefingSpriteWords {
			p := *o.briefWords[n]
			if p == 0 {
				t.Fatalf("missing sprite %s", n)
			}
			dr := o.objects[o.identities[(*client.Drawable)(unsafe.Pointer(uintptr(p)))]].Drawable
			if int(dr.TypeIDVal) != o.c.Things.TypeByID(briefingSpriteNames[i]).Index() || *txword(dr, 120)&0x1000000 == 0 {
				t.Fatalf("sprite type/flag %s", n)
			}
		}
		events := len(o.events)
		again := legacy.PortTestBriefing(7, 0, 0, 0)
		if again != ret || len(o.events) != events {
			t.Fatal("cached sprites reallocated")
		}
		rows = append(rows, o.snapshotBriefing(t, 7, ret))
		o.releaseBriefing()
		for _, n := range briefingSpriteWords {
			if *o.briefWords[n] != 0 {
				t.Fatal("sprite cleanup retained pointer")
			}
		}
		for _, obj := range o.objects {
			if obj.Live {
				t.Fatal("sprite not freed by actual cleanup")
			}
		}
	}
	briefingCapture(t, "sprite-lifecycle", rows, "7db12f469ed96b02a77965c93413e7de683c906d77a368fe578336c15752a67a")
}
