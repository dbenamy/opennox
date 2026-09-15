//go:build porttest

package opennox

import (
	"testing"
	"unsafe"

	"github.com/opennox/libs/spell"
	"github.com/opennox/opennox/v1/client"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

type tooltipOwner struct {
	*objectRenderOwner
	language        func(int)
	dr              *client.Drawable
	scratch, cursor []uint16
	text            [][]uint16
}
type tooltipResult struct {
	Case, Step            int
	Return                uint32
	Text, Scratch, Cursor []uint16
	Drawables             [][]uint32
	Messages              [][]byte
}

func tooltipWide(t *testing.T, text []uint16) []uint16 {
	p, free := alloc.Make([]uint16{}, len(text)+1)
	copy(p, text)
	t.Cleanup(free)
	return p
}
func tooltipMapped(t *testing.T, offset uintptr, count int) []uint16 {
	p := unsafe.Slice((*uint16)(memmap.PtrOff(0x5D4594, offset)), count)
	old := append([]uint16(nil), p...)
	t.Cleanup(func() { copy(p, old) })
	return p
}
func newTooltipOwner(t *testing.T) *tooltipOwner {
	o := &tooltipOwner{objectRenderOwner: newObjectRenderOwner(t)}
	configure, restore := o.c.srv.Server.PortTestTooltipStrings()
	o.language = configure
	t.Cleanup(restore)
	configure(0)
	o.scratch = tooltipMapped(t, 1317000, 1024)
	o.cursor = tooltipMapped(t, 1096676, 256)
	fallback := tooltipMapped(t, 1319048, 1)
	// The only observed fallback is the static empty UTF16 cell. Do not assume
	// additional space beyond it or populate adjacent mapped globals.
	if fallback[0] != 0 {
		t.Fatal("unexpected nonempty tooltip fallback")
	}
	o.dr = o.newUnlinked(t)
	oldNet := o.c.srv.NetList
	o.c.srv.NetList = netlist.New()
	o.c.srv.NetList.Init()
	t.Cleanup(func() { o.c.srv.NetList.Free(); o.c.srv.NetList = oldNet })
	noxflags.SetGame(noxflags.GameHost)
	for _, s := range [][]uint16{{'P', 'r', 'e', 't', 't', 'y'}, {'W', 'e', 'a', 'p', 'o', 'n'}, {'A', 'r', 'm', 'o', 'r'}, {'M', '0'}, {'M', '1'}, {'M', '2'}, {'S', '3'}, {}, {0xd800, 'X', 0xdc01}} {
		o.text = append(o.text, tooltipWide(t, s))
	}
	typ := o.c.Things.TypeByInd(4)
	oldPretty := typ.PrettyName
	typ.PrettyName = &o.text[0][0]
	t.Cleanup(func() { typ.PrettyName = oldPretty })
	o.weapon.Desc8 = &o.text[1][0]
	o.armor.Desc8 = &o.text[2][0]
	for i, m := range o.mods {
		*(*unsafe.Pointer)(unsafe.Add(m.C(), 8)) = unsafe.Pointer(&o.text[3+i][0])
		*(*unsafe.Pointer)(unsafe.Add(m.C(), 12)) = unsafe.Pointer(&o.text[6][0])
	}
	spells, restoreSpells := o.c.srv.Server.PortTestAISpellDefs()
	t.Cleanup(restoreSpells)
	spells([]server.PortTestSpellClassDef{{Index: 1, Valid: true}, {Index: 2, Valid: true}})
	o.c.srv.Spells.DefByInd(spell.ID(1)).Title = "Spark"
	o.c.srv.Spells.DefByInd(spell.ID(2)).Title = "Flame Ω"
	oldAbilities := o.c.srv.abilities.defs
	t.Cleanup(func() { o.c.srv.abilities.defs = oldAbilities })
	for i, name := range []string{"", "Berserk", "Cry Ω", "Harpoon", "Eye", "Sneak"} {
		if i > 0 {
			o.c.srv.abilities.defs[i] = AbilityDef{name: name, field24: 1}
		}
	}
	for i := 1; i <= 2; i++ {
		p := unsafe.Slice((*uint32)(memmap.PtrOff(0x5D4594, uintptr(740076+28*i))), 2)
		old := [2]uint32{p[0], p[1]}
		t.Cleanup(func() { copy(p, old[:]) })
		title := tooltipWide(t, []uint16{'G', 'u', 'i', 'd', 'e', uint16('0' + i)})
		p[0], p[1] = uint32(uintptr(unsafe.Pointer(&title[0]))), 1
	}
	o.prepare(0, 0, 0, 0)
	return o
}
func (o *tooltipOwner) prepare(class, sub, metadata, mask uint32) {
	w := unsafe.Slice((*uint32)(unsafe.Pointer(o.dr)), 128)
	w[28], w[29] = class, sub
	clear(w[108:112])
	if class&0x13001000 != 0 {
		for i, m := range o.mods {
			if mask&(1<<i) != 0 {
				w[108+i] = uint32(uintptr(m.C()))
			}
		}
	} else {
		w[108] = metadata
	}
	for i := range o.scratch {
		o.scratch[i] = uint16(0x8100 + i%256)
	}
	for i := range o.cursor {
		o.cursor[i] = uint16(0x9100 + i%256)
	}
	o.c.srv.NetList.ResetByInd(31, netlist.Kind0)
}
func (o *tooltipOwner) invoke(t *testing.T, id, step int, dr *client.Drawable) tooltipResult {
	p := legacy.PortTestTooltip(dr)
	r := tooltipResult{Case: id, Step: step, Scratch: append([]uint16(nil), o.scratch...), Cursor: append([]uint16(nil), o.cursor...)}
	switch p {
	case &o.scratch[0]:
		r.Return = 1
	case o.c.Things.TypeByInd(4).PrettyName:
		r.Return = 2
	default:
		t.Fatal("tooltip returned an unowned string")
	}
	for i := 0; i < 1024; i++ {
		v := *(*uint16)(unsafe.Add(unsafe.Pointer(p), i*2))
		if v == 0 {
			break
		}
		r.Text = append(r.Text, v)
		if i == 1023 {
			t.Fatal("unterminated tooltip")
		}
	}
	r.Drawables = o.c.snapshotDrawables(t, o.unlinked...)
	for _, row := range r.Drawables {
		w := row[1:]
		if w[28]&0x13001000 != 0 {
			for i := 108; i < 112; i++ {
				if w[i] != 0 {
					ref, ok := o.modRefs[w[i]]
					if !ok {
						t.Fatal("unknown modifier")
					}
					w[i] = ref
				}
			}
		}
	}
	o.c.srv.NetList.ByInd(31, netlist.Kind0).Each(func(b []byte) bool { r.Messages = append(r.Messages, append([]byte(nil), b...)); return false })
	return r
}
