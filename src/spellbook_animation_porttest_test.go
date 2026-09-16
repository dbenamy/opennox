//go:build porttest

package opennox

import (
	"fmt"
	"math"
	"testing"
	"unsafe"

	"github.com/opennox/libs/things"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestSpellbookImmediateAddition(t *testing.T) {
	o := newSpellbookOwner(t)
	configure, restore := o.c.srv.Server.PortTestAISpellDefs()
	t.Cleanup(restore)
	bar := unsafe.Slice(memmap.PtrUint32(0x5D4594, 1049220), 64)
	oldBar := append([]uint32(nil), bar...)
	t.Cleanup(func() { copy(bar, oldBar) })
	t.Cleanup(legacy.PortTestBookQuickbar(unsafe.Pointer(&bar[0])))
	for _, part := range []struct {
		base, off uintptr
		data      []byte
	}{{0x587000, 133488, blobdata.PortTestBookQuickbarSpacing()}, {0x581450, 9872, blobdata.PortTestClientEffectsCurveTable()}} {
		raw := unsafe.Slice(memmap.PtrUint8(part.base, part.off), len(part.data))
		old := append([]byte(nil), raw...)
		t.Cleanup(func() { copy(raw, old) })
		copy(raw, part.data)
	}
	expanded := memmap.PtrUint32(0x5D4594, 1049476)
	oldExpanded := *expanded
	t.Cleanup(func() { *expanded = oldExpanded })
	oldGUI := nox_client_renderGUI_80828
	t.Cleanup(func() { nox_client_renderGUI_80828 = oldGUI })
	for p, id := range legacy.PortTestBookQuickbarCallbacks() {
		o.c.callbackRefs[p] = id
	}
	oldLoad := legacy.Nox_xxx_gLoadImg
	t.Cleanup(func() { legacy.Nox_xxx_gLoadImg = oldLoad })
	legacy.Nox_xxx_gLoadImg = func(name string) *noxrender.Image {
		if name == "QuickBarBase" {
			return o.images[0]
		}
		return oldLoad(name)
	}
	var rows []struct {
		Book      spellbookResult
		Bar       []uint32
		Particles [][]uint32
		RNG       [2]int
		RenderGUI bool
	}
	for _, width := range []uint32{640, 749, 750, 999, 1000, 1280} {
		for _, kind := range []uint32{2, 3, 4} {
			for mode := 0; mode < 4; mode++ {
				func() {
					o.resetBook(t)
					clear(bar)
					nox_client_renderGUI_80828 = false
					*expanded = 0
					*o.words["nox_win_width"] = width
					if o.bookCall("nox_xxx_bookInit_45B9D0") != 1 {
						t.Fatal("book setup")
					}
					for row := 0; row < 5; row++ {
						o.c.dataRefs[uint32(uintptr(unsafe.Pointer(&bar[row*10])))] = 0xee600001 + uint32(row)
					}
					legacy.PortTestBookQuickbarInit(unsafe.Pointer(&bar[0]), 229, 438)
					for i := 0; i < 5; i++ {
						parent := (*gui.Window)(unsafe.Pointer(uintptr(bar[52])))
						w := o.c.GUI.NewWindowRaw(parent, gui.StatusFlags(8), 10+36*i, 1, 10, 10, nil)
						bar[58+i] = uint32(uintptr(w.C()))
					}
					flags := uint32(things.SpellClassAny)
					if mode == 3 {
						flags |= 0x1000
					}
					configure([]server.PortTestSpellClassDef{{Index: 1, Flags: flags, Valid: true}})
					if mode == 1 {
						bar[0] = 1
					}
					if mode == 2 {
						for i := 0; i < 25; i++ {
							bar[2*i] = 2
						}
					}
					particles, freeParticles := legacy.PortTestEffectsScreenParticles(256)
					defer freeParticles()
					ret := o.bookCall("nox_xxx_bookFillAll_45D570", kind, 1)
					blocked := mode == 1 || mode == 2 || (mode == 3 && kind == 2)
					if *o.words["dword_5d4594_1047520"] != 0 {
						t.Fatal("non-coop addition must complete immediately")
					}
					if !blocked && (bar[0] != 1 || len(particles()) != 50 || memmap.Uint32(0x5D4594, 1046680) != 19) {
						t.Fatalf("addition kind%d mode%d slot%d particles%d path%d", kind, mode, bar[0], len(particles()), memmap.Uint32(0x5D4594, 1046680))
					}
					if !blocked {
						speed := 6.0
						if width >= 1000 {
							speed = 10
						} else if width >= 750 {
							speed = 8
						}
						v := legacy.PortTestBookVector()
						if math.Abs(math.Hypot(float64(v[0]), float64(v[1]))-speed) > 0.00001 {
							t.Fatal("screen-width animation speed")
						}
					}
					if blocked && len(particles()) != 0 {
						t.Fatal("blocked addition emitted particles")
					}
					b := append([]uint32(nil), bar...)
					o.collect()
					for i := 51; i <= 62; i++ {
						b[i] = o.normalize(b[i])
					}
					rows = append(rows, struct {
						Book      spellbookResult
						Bar       []uint32
						Particles [][]uint32
						RNG       [2]int
						RenderGUI bool
					}{o.bookSnapshot(fmt.Sprintf("width%d-kind%d-mode%d", width, kind, mode), ret), b, particles(), [2]int{o.c.srv.Rand.Logic.Index(), o.c.srv.Rand.Other.Index()}, nox_client_renderGUI_80828})
				}()
			}
		}
	}
	spellbookCapture(t, "immediate-addition", rows, "")
}
