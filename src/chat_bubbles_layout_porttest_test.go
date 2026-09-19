//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
)

func TestChatBubbleLayoutSingle(t *testing.T) {
	type record struct {
		Name                string
		X, Y, Width, Height int32
		Visible, Arrow      uint32
	}
	var rows []record
	for _, text := range []string{"A", "Speech bubble", "A long message that wraps over several lines of speech"} {
		for _, shifted := range []bool{false, true} {
			for _, pos := range []image.Point{image.Pt(320, 300), image.Pt(0, 300), image.Pt(639, 300), image.Pt(320, 0), image.Pt(320, 65535)} {
				name := fmt.Sprintf("%s/shift=%t/pos=%v", text, shifted, pos)
				t.Run(name, func(t *testing.T) {
					o := newChatBubbleRenderOwner(t)
					restore := noxflags.PortTestGameFlags(0)
					defer restore()
					v := o.c.Viewport()
					if shifted {
						v.World = v.World.Add(image.Pt(10, 20))
					}
					p := o.create(t, 1234, text, byte(len(text)), uint32(pos.X)|uint32(pos.Y)<<16, 20)
					size := o.c.r.GetStringSizeWrapped(nil, text, 128)
					size.X = min(size.X, 128)
					if size.X <= 0 || size.Y <= 0 {
						t.Fatal("font fixture must measure text")
					}
					cap := o.c.r.FontHeight(nil)
					wantX, wantY := pos.X-v.World.Min.X-size.X/2, pos.Y-v.World.Min.Y-64-size.Y
					arrow := uint32(1)
					if wantX < cap {
						wantX, arrow = cap, 0
					}
					if wantX > 640-size.X-cap {
						wantX, arrow = 640-size.X-cap, 0
					}
					if wantY < 2*cap+2 {
						wantY, arrow = 2*cap+2, 0
					}
					if wantY > 480-size.Y-cap {
						wantY, arrow = 480-size.Y-cap, 0
					}
					if wantY < 55 {
						wantY, arrow = 55, 0
					}
					legacy.PortTestChatBubbleLayout(v)
					got := record{name, int32(objectXferGetWord(p, 648)), int32(objectXferGetWord(p, 652)), int32(objectXferGetWord(p, 672)), int32(objectXferGetWord(p, 676)), objectXferGetWord(p, 660), objectXferGetWord(p, 664)}
					want := record{name, int32(wantX), int32(wantY), int32(size.X), int32(size.Y), 1, arrow}
					if got != want {
						t.Fatalf("layout %+v want %+v", got, want)
					}
					if objectXferGetWord(p, 668) != 0 || objectXferGetWord(p, 680) != 0 {
						t.Fatal("unexpected attachment/ordinal")
					}
					rows = append(rows, got)
				})
			}
		}
	}
	spellbookCapture(t, "chat-bubble-layout-single", rows, "a71c483ce24b72cfdee13fc4e299a7da334a5c25b51b049fc36dcaa7cd733eb9")
}

func TestChatBubbleExpiry(t *testing.T) {
	type record struct {
		Name           string
		Live, Appended []chatBubbleRecord
	}
	var rows []record
	for _, frame := range []uint32{0, 123, 0xfffffff0} {
		for _, delta := range []uint32{0, 1, 2, 16} {
			for mask := 0; mask < 8; mask++ {
				name := fmt.Sprintf("frame=%x/delta=%d/mask=%d", frame, delta, mask)
				t.Run(name, func(t *testing.T) {
					o := newChatBubbleRenderOwner(t)
					o.c.srv.SetFrame(frame)
					var expected []uint32
					for i := 0; i < 3; i++ {
						ttl := uint16(40)
						if mask&(1<<i) != 0 {
							ttl = 1
						}
						o.create(t, uint16(i+1), "Expiry", 6, uint32(300)|uint32(300)<<16, ttl)
						if frame+delta <= frame+uint32(ttl) {
							expected = append(expected, uint32(i+1))
						}
					}
					o.c.srv.SetFrame(frame + delta)
					legacy.PortTestChatBubbleLayout(o.c.Viewport())
					r := record{Name: name, Live: o.snapshot(t)}
					if len(r.Live) != len(expected) {
						t.Fatalf("live count %d want %d", len(r.Live), len(expected))
					}
					for i, v := range r.Live {
						if v.Code != expected[i] {
							t.Fatal("expiry order")
						}
					}
					o.create(t, 4, "Append", 6, 0, 40)
					r.Appended = o.snapshot(t)
					if len(r.Appended) != len(expected)+1 || r.Appended[len(expected)].Code != 4 {
						t.Fatal("append after expiry")
					}
					rows = append(rows, r)
				})
			}
		}
	}
	spellbookCapture(t, "chat-bubble-expiry", rows, "e820394c3b135ab1ee3a6dcc0c24daea88aec8cc79390e13f6f8bdf046b00d01")
}

func TestChatBubblePixels(t *testing.T) {
	type record struct {
		Text   string
		X      int
		Pixels string
	}
	var rows []record
	seen := make(map[string]bool)
	for _, text := range []string{"Hello", "Different text", "This longer message wraps onto several lines"} {
		for _, x := range []int{200, 400} {
			t.Run(fmt.Sprintf("%s/%d", text, x), func(t *testing.T) {
				o := newChatBubbleRenderOwner(t)
				restore := noxflags.PortTestGameFlags(0)
				defer restore()
				blank := effectsPixelHash(o.pix)
				o.create(t, 1234, text, byte(len(text)), uint32(x)|uint32(300)<<16, 20)
				legacy.PortTestChatBubbleDraw(o.c.Viewport())
				hash := effectsPixelHash(o.pix)
				if hash == blank {
					t.Fatal("bubble did not draw")
				}
				if seen[hash] {
					t.Fatal("text/position change did not alter pixels")
				}
				seen[hash] = true
				rows = append(rows, record{text, x, hash})
			})
		}
	}
	spellbookCapture(t, "chat-bubble-pixels", rows, "4c9ad61cf27bfc2003153eb47c7267749172f559a769a2ab0924a05440d9b3f6")
}
