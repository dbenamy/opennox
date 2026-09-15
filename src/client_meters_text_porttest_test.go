//go:build porttest

package opennox

import (
	"image"
	"strconv"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
	"golang.org/x/image/font"
)

type meterTextCall struct {
	Text     string
	Position image.Point
	Return   int
}
type meterRender struct {
	*objectDrawingRender
	owner *meterOwner
}

func (r *meterRender) DrawString(face font.Face, text string, pos image.Point) int {
	ret := r.objectDrawingRender.DrawString(face, text, pos)
	r.owner.text = append(r.owner.text, meterTextCall{text, pos, ret})
	return ret
}

type meterClient struct {
	*objectDrawingClient
	owner *meterOwner
}

func (c *meterClient) R2() legacy.Render2 {
	return &meterRender{&objectDrawingRender{c.owner.c.r, c.owner.objectDrawingOwner}, c.owner}
}
func TestClientMetersLabelContract(t *testing.T) {
	o := newMeterOwner(t)
	for _, value := range []uint32{0, 999, 1000, 65535, 0x7fffffff, 0x80000000, 0xffffffff} {
		o.plain(t)
		o.meters.Records[6].Current = value
		blank := effectsPixelHash(o.pix)
		legacy.PortTestMeterCall(19, o.meters.Records[6].Window, 0, 0, 0, 0)
		if effectsPixelHash(o.pix) == blank {
			t.Fatal("charge label produced no visible glyph pixels")
		}
		want := strconv.FormatInt(int64(int32(value)), 10)
		if len(o.text) != 1 || o.text[0].Text != want {
			t.Fatalf("count%x text%v want%q", value, o.text, want)
		}
	}
}
