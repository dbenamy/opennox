//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"image"
	"slices"
	"testing"
	"unsafe"

	noxcolor "github.com/opennox/libs/color"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/client/audio/ail"
	"github.com/opennox/opennox/v1/legacy/dialog"
	"github.com/opennox/opennox/v1/legacy/timer"
)

func TestClientInteractionConversationBlink(t *testing.T) {
	o := newMeterOwner(t)
	tint := serverConfigOwnBytes(t, 0x5D4594, 1107052, 4)
	binary.LittleEndian.PutUint32(tint, 0x03e0)
	parent := o.c.GUI.NewWindowRaw(nil, 8, 3, 4, 100, 100, nil)
	defer parent.Destroy()
	w := o.c.GUI.NewWindowRaw(parent, 8, 20, 30, 30, 14, nil)
	draw := w.DrawData()
	draw.BgColorVal = 0x80000000
	draw.EnColorVal = 0x80000000
	draw.DisColorVal = 0x80000000
	draw.HlColorVal = 0x80000000
	draw.TextColorVal = 0x80000000
	var state [6]uint32
	var driver ail.Driver
	var timers [4]timer.TimerGroup
	oldDialog := legacy.Dialogs
	defer func() { legacy.Dialogs = oldDialog }()
	legacy.Dialogs = dialog.NewDialog("dialog", &state[0], &state[1], &state[2], &state[3], &state[4], &driver, &state[5], o.c.srv.Strings, &timers[0], &timers[1], &timers[2], &timers[3], nil, nil, nil, nil, nil, nil, nil)
	baseline := *o.c.r.Data()
	type row struct {
		Sequence               uint
		Enabled, Queued, Blink bool
		Pixels                 string
	}
	var captured []row
	for _, seq := range []uint{0, 1, 7, 8, 15, 16, 23, 24, 29, 30, 31, 127, 128, 135, 136, 151, 152, 157, 158, 159, 255, 256, 264, 285, 286, 511, 0x10008} {
		for o.c.Inp.CurrentSeq() < seq {
			o.c.Inp.Tick()
		}
		for _, enabled := range []bool{false, true} {
			for _, queued := range []bool{false, true} {
				legacy.Dialogs.Sub_44D8F0()
				state[0] = 1
				if queued {
					legacy.Dialogs.PlayFile("BlinkVoice.wav", 100)
				}
				if !enabled {
					state[0] = 0
				}
				*o.c.r.Data() = baseline
				clear(o.pix.Pix)
				o.c.r.ClearPoints()
				if interactionCall("sub_479C40", uintptr(w.C()), uintptr(unsafe.Pointer(draw))) != 1 {
					t.Fatal("conversation blink draw return")
				}
				got := append([]uint16(nil), o.pix.Pix...)
				blink := !(enabled && queued) && byte(seq)&0x7f < 30 && byte(seq)&8 != 0
				*o.c.r.Data() = baseline
				clear(o.pix.Pix)
				o.c.r.ClearPoints()
				if blink {
					c := noxrender.SplitColor(noxcolor.RGBA5551(0x03e0))
					o.c.r.Data().SetColorInt54(noxrender.RGB{R: int(c.R), G: int(c.G), B: int(c.B)})
					o.c.r.Data().SetField262(4)
					for _, line := range [][2]image.Point{{{23, 34}, {52, 34}}, {{53, 34}, {53, 45}}, {{53, 46}, {23, 46}}, {{23, 45}, {23, 35}}} {
						o.c.r.AddPoint(line[0])
						o.c.r.AddPoint(line[1])
						o.c.r.DrawParticles49ED80(64)
					}
				}
				if !slices.Equal(got, o.pix.Pix) {
					t.Fatal("conversation blink gate/border geometry", seq, enabled, queued, blink)
				}
				nonzero := false
				raw := make([]byte, 2*len(got))
				for i, v := range got {
					nonzero = nonzero || v != 0
					binary.LittleEndian.PutUint16(raw[2*i:], v)
				}
				if nonzero != blink {
					t.Fatal("empty conversation blink fixture", seq, enabled, queued, blink)
				}
				captured = append(captured, row{seq, enabled, queued, blink, fmt.Sprintf("%x", sha256.Sum256(raw))})
			}
		}
	}
	interactionCapture(t, "conversation-blink", captured)
}
