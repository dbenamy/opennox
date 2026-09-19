//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestResourceDefinitionsCatalogLinks(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(o.s.PortTestRewardTypes([]string{"PortA", "PortB"}, nil, true, 0, 0))
	var mods []*server.ModifierEff
	var names []*byte
	for _, name := range []string{"ModA", "ModB"} {
		m, free := alloc.New(server.ModifierEff{})
		t.Cleanup(free)
		p, freeS := alloc.CString(name)
		t.Cleanup(freeS)
		mods = append(mods, m)
		names = append(names, p)
	}
	t.Cleanup(o.s.PortTestControlsModifiers(mods, names))
	offsets := []uintptr{208180, 210712, 210856, 211000, 209344}
	tables := make([][]uint32, 5)
	for i, off := range offsets {
		stride := 6
		start := off - 4
		if i == 0 {
			stride = 5
			start = off
		}
		tables[i] = unsafe.Slice(memmap.PtrUint32(0x587000, start), 4*stride)
		saved := append([]uint32(nil), tables[i]...)
		ii := i
		t.Cleanup(func() { copy(tables[ii], saved) })
	}
	var rows []map[string]any
	for mask := 0; mask < 32; mask++ {
		var frees []func()
		for i, table := range tables {
			for j := range table {
				table[j] = 0xa5a5a5a5
			}
			stride, nameword := 6, 1
			if i == 0 {
				stride, nameword = 5, 0
			}
			inputs := []string{"ModA", "missing", "ModB"}
			if i == 0 {
				inputs = []string{"#PortA", "PortB", "#missing"}
			}
			for j, name := range inputs {
				p, free := alloc.CString(name)
				frees = append(frees, free)
				table[j*stride+nameword] = uint32(uintptr(unsafe.Pointer(p)))
			}
			table[3*stride+nameword] = 0
			if mask&(1<<i) == 0 {
				table[nameword] = 0
			}
		}
		if legacy.PortTestResourceLinkCatalogs() != 0 {
			t.Fatal("catalog return")
		}
		var state [][]uint32
		for i, table := range tables {
			stride, nameword, valueword := 6, 1, 0
			if i == 0 {
				stride, nameword, valueword = 5, 0, 1
			}
			snap := append([]uint32(nil), table...)
			for j := 0; j < 3; j++ {
				want := uint32(0xa5a5a5a5)
				if mask&(1<<i) != 0 {
					if i == 0 {
						if j == 0 {
							want = uint32(o.s.Types.IndByID("PortA"))
						} else if j == 1 {
							want = uint32(o.s.Types.IndByID("PortB"))
						} else {
							want = 0
						}
					} else {
						if j == 1 {
							want = 0
						} else {
							m := mods[0]
							if j == 2 {
								m = mods[1]
							}
							want = uint32(uintptr(unsafe.Pointer(m)))
						}
					}
				}
				if table[j*stride+valueword] != want {
					t.Fatalf("catalog mask%d table%d row%d", mask, i, j)
				}
				if i > 0 && mask&(1<<i) != 0 && j != 1 {
					snap[j*stride+valueword] = uint32(j + 1)
				}
				snap[j*stride+nameword] = uint32(j + 1)
				if j == 0 && mask&(1<<i) == 0 {
					snap[nameword] = 0
				}
				for k := 0; k < stride; k++ {
					if k != nameword && k != valueword && table[j*stride+k] != 0xa5a5a5a5 {
						t.Fatal("adjacent catalog field")
					}
				}
			}
			state = append(state, snap)
		}
		rows = append(rows, map[string]any{"mask": mask, "tables": state})
		for _, f := range frees {
			f()
		}
	}
	spellbookCapture(t, "resource-definitions-catalog-links", rows, "204c57215cc922f7babe66caed6ae3d2287788a7a23fe03487f7f4da2eea5a9e")
}
func TestResourceDefinitionsMonsterSound(t *testing.T) {
	u, free := alloc.New(server.Object{})
	defer free()
	data, freeData := alloc.New(server.MonsterUpdateData{})
	defer freeData()
	set, freeSet := alloc.New(uint32(0x12345678))
	defer freeSet()
	u.UpdateData = unsafe.Pointer(data)
	if legacy.Nox_xxx_monsterGetSoundSet_424300(nil) != nil {
		t.Fatal("nil monster")
	}
	var rows []map[string]any
	for _, flags := range []uint32{0, 1, 2, 3, 4, 0x102, 0xffffffff} {
		for _, present := range []bool{false, true} {
			u.ObjClass = object.Class(flags)
			data.SoundSet122 = nil
			if present {
				data.SoundSet122 = unsafe.Pointer(set)
			}
			got := legacy.Nox_xxx_monsterGetSoundSet_424300(u)
			want := unsafe.Pointer(nil)
			if flags&2 != 0 && present {
				want = unsafe.Pointer(set)
			}
			if got != want {
				t.Fatal(fmt.Sprintf("monster sound %x %v", flags, present))
			}
			rows = append(rows, map[string]any{"flags": flags, "present": present, "found": got != nil})
		}
	}
	spellbookCapture(t, "resource-definitions-monster-sound", rows, "4f65333521fd007a469cfde08293e2d59910cfd8975c3dc9de6bd86d0311e724")
}
