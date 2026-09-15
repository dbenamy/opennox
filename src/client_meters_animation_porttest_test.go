//go:build porttest

package opennox

import (
	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func meterBubbles(record int) []int32 {
	return unsafe.Slice((*int32)(memmap.PtrOff(0x5D4594, 1093180+uintptr(record)*1536)), 384)
}

func TestClientMetersBubbleContract(t *testing.T) {
	o := newMeterOwner(t)
	o.plain(t)
	r := &o.meters.Records[0]
	r.Current, r.Maximum = 100, 100
	*o.meters.NamedWord("nox_client_renderBubbles_80844") = 1
	b := meterBubbles(0)
	for i := 0; i < 64; i++ {
		copy(b[i*6:], []int32{3, 120 * 16, 2, 16, 1, 0x7fff7fff})
	}
	before := o.c.srv.Rand.Other.Index()
	legacy.PortTestMeterCall(36, r.Window, 0, 0, 0, 0)
	for i := 0; i < 64; i++ {
		if b[i*6+1] != 119*16 || b[i*6+4] != 1 {
			t.Fatal("active bubble movement/lifetime")
		}
	}
	if o.c.srv.Rand.Other.Index() != before {
		t.Fatal("full active bubble pool consumed RNG")
	}
	legacy.PortTestMeterCall(2, nil, 0, 0, 0, 0)
	for i := 0; i < 64; i++ {
		if b[i*6+4] != 0 || b[i*6+1] != 119*16 {
			t.Fatal("color initialization must clear only bubble active flags")
		}
	}
}

func TestClientMetersBubbleMatrix(t *testing.T) {
	o := newMeterOwner(t)
	var out []meterResult
	id := 0
	for _, seed := range []int{1, 17, 999} {
		for record := 0; record < 2; record++ {
			for _, level := range []uint32{0, 1, 2, 25, 50, 99, 100} {
				for mode := 0; mode < 3; mode++ {
					o.plain(t)
					o.c.srv.Rand.Logic, o.c.srv.Rand.Other = prand.New(seed), prand.New(seed+1)
					legacy.PortTestMeterCall(2, nil, 0, 0, 0, 0)
					*o.meters.NamedWord("nox_client_renderBubbles_80844") = 1
					if mode == 2 {
						// This case exercises the loaded poison overlay, as in
						// gameplay. A nil image leaves the shared image-output
						// coordinates untouched and is not a valid draw trace.
						*memmap.PtrUint32(0x5D4594, 1091900) = uint32(uintptr(o.images[10].C()))
						legacy.PortTestMeterCall(3, nil, 1, 0, 0, 0)
					}
					b := meterBubbles(record)
					clear(b)
					if mode != 0 {
						for i := 0; i < 64; i++ {
							copy(b[i*6:], []int32{int32(i % 13), int32((i*2)%126) * 16, int32(i%3 + 1), int32(4 + i%45), int32(i % 2), 0x7fff7fff})
						}
					}
					r := &o.meters.Records[record]
					r.Current, r.Maximum = level, 100
					for step := 0; step < 4; step++ {
						out = append(out, o.invoke(t, id, step, 36, r.Window, 0, 0, 0, 0))
					}
					id++
				}
			}
		}
	}
	meterCapture(t, "meters-bubbles", out, len(out), "56ed92b179637625063f77a1c3c27a97efdb0241ac58a0c33c2f1232a8d41e78")
}

func TestClientMetersChargeRasterMatrix(t *testing.T) {
	o := newMeterOwner(t)
	var out []meterResult
	id := 0
	for _, maximum := range []uint32{0, 1, 2, 3, 7, 29, 30, 31, 32, 60, 61, 62, 100, 255} {
		for _, current := range []uint32{0, 1, maximum / 2, maximum, maximum + 1} {
			o.plain(t)
			o.c.Inp.Tick()
			r := &o.meters.Records[5]
			r.Current, r.Maximum = current, maximum
			for step := 0; step < 3; step++ {
				if step == 2 {
					o.c.Inp.Tick()
				}
				out = append(out, o.invoke(t, id, step, 18, r.Window, 0, 0, 0, 0))
			}
			id++
		}
	}
	meterCapture(t, "meters-charge-raster", out, len(out), "8957381b6dc1597772b65331d34de1e3a25c184358d9b9534f25c7d6c4ab7387")
}

func TestClientMetersChargeRasterContract(t *testing.T) {
	o := newMeterOwner(t)
	o.plain(t)
	o.c.Inp.Tick()
	r := &o.meters.Records[5]
	r.Current, r.Maximum = 3, 7
	legacy.PortTestMeterCall(18, r.Window, 0, 0, 0, 0)
	first := effectsPixelHash(o.pix)
	legacy.PortTestMeterCall(18, r.Window, 0, 0, 0, 0)
	if effectsPixelHash(o.pix) != first {
		t.Fatal("same input sequence redrew alpha charge rows")
	}
	clear(o.pix.Pix)
	blank := effectsPixelHash(o.pix)
	legacy.PortTestMeterCall(18, r.Window, 0, 0, 0, 0)
	if effectsPixelHash(o.pix) != blank {
		t.Fatal("same sequence drew cleared charge rows")
	}
	o.c.Inp.Tick()
	legacy.PortTestMeterCall(18, r.Window, 0, 0, 0, 0)
	if effectsPixelHash(o.pix) == blank {
		t.Fatal("next input sequence did not redraw charge rows")
	}
}
