//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestClientDrawableEffectTypes(t *testing.T) {
	// Names are the shipped 0x587000+161216 pointer table, audited against blob data.
	fire := []string{"PitifulFireball", "WeakFireball", "Fireball", "StrongFireball", "TitanFireball"}
	names := append([]string{"Spark", "BlueSpark", "YellowSpark", "CyanSpark", "GreenSpark", "Puff"}, fire...)
	names = append(names, "VioletSpark")
	c, _, _ := newEffectsFullOwner(t, append([]string{"Spark", "Puff", "VioletSpark"}, fire...)...)
	words, restore := legacy.PortTestDrawableEffectGlobals()
	t.Cleanup(restore)
	table := unsafe.Slice(memmap.PtrPtr(0x587000, 161216), 5)
	saved := append([]unsafe.Pointer(nil), table...)
	t.Cleanup(func() { copy(table, saved) })
	for i, name := range fire {
		p, free := alloc.CString(name)
		t.Cleanup(free)
		table[i] = unsafe.Pointer(p)
	}
	slots := []int{15, 2, 17, 18, 19, 20, 25, 26, 27, 28, 29, 3}
	type row struct {
		Missing, Return int
		Globals         []uint32
	}
	var rows []row
	for missing := -1; missing < len(names); missing++ {
		for i, p := range words {
			*p = 0xface0000 + uint32(i)
		}
		want := make([]uint32, len(words))
		for i, p := range words {
			want[i] = *p
		}
		restoreType := func() {}
		if missing >= 0 {
			restoreType = c.Cli().PortTestDrawableHideType(names[missing])
		}
		for i, name := range names {
			want[slots[i]] = uint32(c.Things.IndByID(name))
			if i == missing {
				break
			}
		}
		ret := int(legacy.PortTestDrawableEffect(5, nil, [5]uint32{}))
		restoreType()
		wantRet := 0
		if missing == -1 {
			wantRet = 1
		}
		if ret != wantRet {
			t.Fatalf("init missing%d return%d", missing, ret)
		}
		got := make([]uint32, len(words))
		for i, p := range words {
			got[i] = *p
			if got[i] != want[i] {
				t.Fatalf("missing%d word%d got%x want%x", missing, i, got[i], want[i])
			}
		}
		rows = append(rows, row{missing, ret, got})
	}
	drawableStateCapture(t, "types", rows, "7844c843adee3ec01b4c248eca24b12b677c2ab35bf3df0484ff6d42891b148b")
}
