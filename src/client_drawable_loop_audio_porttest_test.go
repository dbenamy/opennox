//go:build porttest

package opennox

import (
	"image"
	"math"
	"testing"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestClientDrawableLoopAudio(t *testing.T) {
	drawing := newObjectDrawingOwner(t)
	o := newAudioEventsOwner(t, 1)
	*o.words["dword_5d4594_1045432"] = 1
	dimensions := legacy.PortTestBindingDimensions()
	oldWidth := *dimensions[0]
	t.Cleanup(func() { *dimensions[0] = oldWidth })
	*drawing.c.Viewport() = noxrender.Viewport{Screen: image.Rect(10, 20, 330, 220), World: image.Rect(100, 200, 420, 400), Size: image.Pt(320, 200)}
	dr, free := alloc.New(client.Drawable{})
	t.Cleanup(free)
	listener, free := alloc.New(client.Drawable{})
	t.Cleanup(free)
	meta := o.metadata(1, 0, 1)
	mw := audioStreamWords(meta, 50)
	type row struct {
		Width, Range, DX, DY, Phase                 int32
		Flags                                       uint32
		Present                                     bool
		State, Volume, VolumeTarget, Pan, PanTarget uint32
	}
	var rows []row
	deltas := []image.Point{{0, 0}, {1, 1}, {-1, -1}, {99, 0}, {100, 0}, {-100, 0}, {0, 100}, {0, -100}, {70, 70}, {-71, -71}, {46340, 46340}, {-46340, -46340}, {0x7fffffff, 0}}
	for _, width := range []int32{96, 640, -640} {
		*dimensions[0] = width
		for _, limit := range []int32{-1, 0, 1, 2, 100, 0x7fffffff} {
			mw[16] = uint32(limit)
			for _, delta := range deltas {
				for _, flags := range []uint32{0, 4, 8, 12, 0x100} {
					o.eventCall("sub_4521F0")
					*dr = client.Drawable{}
					*listener = client.Drawable{}
					dr.PosVec = image.Pt(32, 64)
					listener.PosVec = dr.PosVec.Add(delta)
					dr.AudioLoop = 1
					dr.ObjFlags = 0x1000000
					dr.Flags70Val = flags
					dx, dy := int32(delta.X), int32(delta.Y)
					sum := dx*dx + dy*dy + 1
					percent := int32(0)
					if flags&12 == 0 && dx < limit && dy < limit && limit > 0 && sum >= 0 {
						distance := int32(math.Sqrt(float64(sum)))
						if distance < limit {
							percent = 100 * (limit - distance) / limit
							if percent > 100 {
								percent = 100
							}
							if percent < 0 {
								percent = 0
							}
						}
					}
					for phase := int32(0); phase < 2; phase++ {
						legacy.Sub_45A9B0(dr, listener)
						p := dr.Field_124
						// Negative squared-distance values expose C's target conversion behavior;
						// capture those separately while independently checking defined cases.
						if sum >= 0 && ((p != nil) != (percent != 0)) {
							t.Fatal("sound range/gate", width, limit, delta, flags, p != nil, percent)
						}
						r := row{Width: width, Range: limit, DX: dx, DY: dy, Phase: phase, Flags: flags, Present: p != nil}
						if p != nil {
							w := audioStreamWords(audioStreamPointer(p), 144)
							r.State, r.Volume, r.VolumeTarget, r.Pan, r.PanTarget = w[7], w[47], w[48], w[63], w[64]
							if sum >= 0 && r.VolumeTarget != uint32(percent*163)<<16 {
								t.Fatal("loop volume", width, limit, delta, phase, percent, r)
							}
							// C uses viewport word6 (World.Max.X), not World.Min.X.
							pan := int32(50*(32-420-10)) / (width / 2)
							if pan > 50 {
								pan = 50
							}
							if pan < -50 {
								pan = -50
							}
							if r.PanTarget != uint32(pan*8192/50+8192)<<16 {
								t.Fatal("loop pan", width, phase, r)
							}
						}
						rows = append(rows, r)
					}
					// Turning off the drawable's audible flag stops an existing event.
					p := dr.Field_124
					dr.ObjFlags = 0
					legacy.Sub_45A9B0(dr, listener)
					if p != nil && audioStreamWords(audioStreamPointer(p), 144)[7] != 4 {
						t.Fatal("loop stop")
					}
				}
			}
		}
	}
	o.eventCall("sub_4521F0")
	dr.Field_124 = nil
	interactionCapture(t, "client-drawable-loop-audio", rows)
}
