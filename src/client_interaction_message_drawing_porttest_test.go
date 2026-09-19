//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"image"
	"os"
	"slices"
	"testing"
	"unsafe"

	noxcolor "github.com/opennox/libs/color"
	"github.com/opennox/opennox/v1/legacy"
	"golang.org/x/image/font/basicfont"
)

func TestClientInteractionMessageDrawing(t *testing.T) {
	o := newMeterOwner(t)
	words, restore := legacy.PortTestClientInteractionWords()
	defer restore()
	colors, restoreColors := legacy.PortTestInventoryDisplayWords()
	defer restoreColors()
	*colors["nox_color_black_2650656"] = 0x0000
	*colors["nox_color_white_2523948"] = 0x7fff
	older := serverConfigOwnBytes(t, 0x5D4594, 2597996, 4)
	binary.LittleEndian.PutUint32(older, 0x03e0)
	rows := serverConfigOwnBytes(t, 0x5D4594, 823804, 1932)
	serverConfigOwnBytes(t, 0x5D4594, 825740, 2)
	shadows := serverConfigOwnBytes(t, 0x587000, 107848, 32)
	blob, err := os.ReadFile("common/memmap/nox/blobdata/blob_587000.dat")
	if err != nil {
		t.Fatal(err)
	}
	copy(shadows, blob[107848:107880])
	wantShadows := []image.Point{{-1, -1}, {1, -1}, {-1, 1}, {1, 1}}
	for i, p := range wantShadows {
		if int32(binary.LittleEndian.Uint32(shadows[8*i:])) != int32(p.X) || int32(binary.LittleEndian.Uint32(shadows[8*i+4:])) != int32(p.Y) {
			t.Fatal("shipped message shadow table")
		}
	}
	resetFont := o.c.Render().GetFonts().PortTestDefaultFont(basicfont.Face7x13)
	defer resetFont()
	dim := legacy.PortTestBindingDimensions()
	oldw := *dim[0]
	defer func() { *dim[0] = oldw }()
	*dim[0] = 200
	vp := o.c.Viewport()
	oldVP := *vp
	defer func() { *vp = oldVP }()
	vp.Size = image.Pt(200, 120)
	vp.Screen = image.Rect(0, 10, 200, 130)
	oldFrame := o.c.srv.Frame()
	defer o.c.srv.SetFrame(oldFrame)
	for i := range rows {
		rows[i] = 0x5a
	}
	wantClear := append([]byte(nil), rows...)
	for slot := 0; slot < 3; slot++ {
		binary.LittleEndian.PutUint16(wantClear[644*slot:], 0)
		binary.LittleEndian.PutUint32(wantClear[644*slot+636:], 0)
		wantClear[644*slot+640] = 0
	}
	ret := interactionCall("sub_445450")
	if uintptr(ret) != uintptr(unsafe.Pointer(&rows[1288])) || *words["dword_5d4594_825736"] != 0 || !bytes.Equal(rows, wantClear) {
		t.Fatal("message clear ownership/row layout")
	}
	texts := []string{"One", "Two two", "Three"}
	type row struct {
		Frame              uint32
		Start, Mask, Drawn int
		Return             uint32
		Pixels             string
	}
	var captured []row
	for _, frame := range []uint32{100, 0xffffffff} {
		for start := 0; start < 3; start++ {
			for mask := 0; mask < 8; mask++ {
				clear(rows)
				*words["dword_5d4594_825736"] = uint32(start)
				o.c.srv.SetFrame(frame)
				for slot, text := range texts {
					for i, b := range []byte(text) {
						binary.LittleEndian.PutUint16(rows[644*slot+2*i:], uint16(b))
					}
					expiry := frame - 1
					if mask&(1<<slot) != 0 {
						expiry = frame
					}
					binary.LittleEndian.PutUint32(rows[644*slot+636:], expiry)
				}
				before := append([]byte(nil), rows...)
				clear(o.pix.Pix)
				gotReturn := uint32(interactionCall("nox_xxx_drawMessageLines_445530"))
				got := append([]uint16(nil), o.pix.Pix...)
				if !bytes.Equal(rows, before) {
					t.Fatal("message draw changed rows")
				}
				clear(o.pix.Pix)
				drawn := 0
				wantReturn := frame
				for order := 0; order < 3; order++ {
					slot := (start - order + 3) % 3
					if mask&(1<<slot) == 0 {
						break
					}
					text := texts[slot]
					// Face7x13 has 11px cap height; C adds a 4px gap.
					pos := image.Pt((200-7*len(text))/2, 85-15*order)
					o.c.r.Data().SetTextColor(noxcolor.RGBA5551(0x0000))
					for _, off := range wantShadows {
						o.c.r.DrawString(nil, text, pos.Add(off))
					}
					color := uint32(0x7fff)
					if order != 0 {
						color = 0x03e0
					}
					o.c.r.Data().SetTextColor(noxcolor.RGBA5551(color))
					o.c.r.DrawString(nil, text, pos)
					drawn++
					if order == 2 {
						wantReturn = uint32(slot)
					}
				}
				if gotReturn != wantReturn || !slices.Equal(got, o.pix.Pix) {
					t.Fatal("message drawing/order/expiry", frame, start, mask, gotReturn, wantReturn)
				}
				nonzero := false
				for _, v := range got {
					nonzero = nonzero || v != 0
				}
				if (drawn != 0) != nonzero {
					t.Fatal("empty message renderer fixture", drawn)
				}
				raw := make([]byte, 2*len(got))
				for i, v := range got {
					binary.LittleEndian.PutUint16(raw[2*i:], v)
				}
				captured = append(captured, row{frame, start, mask, drawn, gotReturn, fmt.Sprintf("%x", sha256.Sum256(raw))})
			}
		}
	}
	interactionCapture(t, "message-drawing", captured)
}
