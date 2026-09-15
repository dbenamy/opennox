//go:build porttest

package opennox

import (
	noxcolor "github.com/opennox/libs/color"
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func TestClientMetersDrawGatesMatrix(t *testing.T) {
	o := newMeterOwner(t)
	var out []meterResult
	for mode := 0; mode < 5; mode++ {
		for record := 0; record < 2; record++ {
			o.plain(t)
			noxflags.ResetGame()
			*(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 2092)) = 0
			switch mode {
			case 1:
				*o.meters.NamedWord("nox_gameDisableMapDraw_5d4594_2650672") = 1
			case 2:
				noxflags.SetGame(0x100000)
			case 3:
				noxflags.SetGame(0x800000)
			case 4:
				*(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 2092)) = 1
				*(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 3680)) = 1
			}
			blank := effectsPixelHash(o.pix)
			r := o.invoke(t, mode*2+record, 0, 22, o.meters.Records[record].Window, 0, 0, 0, 0)
			if (r.Render.Draw.Pixels == blank) != (mode != 0) {
				t.Fatalf("mini-bar visibility mode%d record%d", mode, record)
			}
			out = append(out, r)
		}
	}
	for i, p := range [][2]int{{-2, -2}, {0, 0}, {20, 30}, {254, 254}, {256, 256}} {
		o.plain(t)
		o.c.r.Data().SetColor2(noxcolor.RGB5551Color(255, 255, 255))
		out = append(out, o.invoke(t, 10+i, 0, 24, nil, p[0], p[1], 0, 0))
	}
	meterCapture(t, "meters-draw-gates", out, len(out), "ebfbff18366c3bcaf8fe3ccc084f7604b29df86b1b57239f33745b1920ee2a4d")
}

func TestClientMetersManaToggleContract(t *testing.T) {
	o := newMeterOwner(t)
	p := memmap.PtrUint8(0x85B3FC, 12254)
	old := *p
	t.Cleanup(func() { *p = old })
	for _, mana := range []byte{0, 1} {
		o.plain(t)
		*p = mana
		*o.meters.NamedWord("dword_5d4594_1096252") = 1
		legacy.PortTestMeterCall(25, o.meters.Records[0].Window, 7, 0, 0, 0)
		if !o.meters.Records[2].Window.GetFlags().Has(16) {
			t.Fatal("health mini-bar not hidden")
		}
		if o.meters.Records[3].Window.GetFlags().Has(16) != (mana != 0) {
			t.Fatal("mana mini-bar toggle gate")
		}
	}
}
