//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestChatBubbleAttachmentAndVisibility(t *testing.T) {
	type record struct {
		Name           string
		Position       uint32
		X, Y           int32
		Visible, Arrow uint32
		Attached       bool
	}
	var rows []record
	for _, static := range []bool{false, true} {
		for _, matching := range []bool{false, true} {
			for _, restricted := range []bool{false, true} {
				for _, pos := range []image.Point{image.Pt(320, 300), image.Pt(0, 0), image.Pt(70000, 65536+300)} {
					name := fmt.Sprintf("static=%t/matching=%t/restricted=%t/pos=%v", static, matching, restricted, pos)
					t.Run(name, func(t *testing.T) {
						o := newChatBubbleRenderOwner(t)
						flags := noxflags.GameFlag(0)
						if restricted {
							flags = 2048
						}
						restore := noxflags.PortTestGameFlags(flags)
						defer restore()
						dr := o.drawable(1234, pos)
						dr.ObjClass = 0
						if static == matching {
							dr.ObjClass = 0x20000000
						}
						code := uint16(1234)
						if static {
							code |= 0x8000
						}
						fallback := uint32(320) | uint32(300)<<16
						p := o.create(t, code, "Attached", 8, fallback, 20)
						legacy.PortTestChatBubbleLayout(o.c.Viewport())
						attached := *(*unsafe.Pointer)(unsafe.Add(p, 668))
						if matching && attached != unsafe.Pointer(dr) || !matching && attached != nil {
							t.Fatal("wrong drawable lookup class")
						}
						wantPos := fallback
						if matching {
							wantPos = uint32(uint16(pos.X)) | uint32(uint16(pos.Y))<<16
						}
						if objectXferGetWord(p, 644) != wantPos {
							t.Fatal("source coordinates must truncate to 16 bits")
						}
						visible := uint32(1)
						if restricted && matching && pos != image.Pt(320, 300) {
							visible = 0
						}
						if objectXferGetWord(p, 660) != visible {
							t.Fatalf("visibility %d want %d", objectXferGetWord(p, 660), visible)
						}
						rows = append(rows, record{name, wantPos, int32(objectXferGetWord(p, 648)), int32(objectXferGetWord(p, 652)), visible, objectXferGetWord(p, 664), attached != nil})
					})
				}
			}
		}
	}
	spellbookCapture(t, "chat-bubble-attachment-visibility", rows, "56f4ed39c8e782bd37551fdd619060394e746fc7ffb4ed944bd305b81fca4076")
}

func TestChatBubbleCandidateAndArrange(t *testing.T) {
	type record struct {
		Name     string
		Accepted bool
		X, Y     int32
	}
	var rows []record
	for _, visible := range []uint32{0, 1} {
		for _, priorVisible := range []uint32{0, 1} {
			for _, ordinal := range []uint32{0, 1, 2} {
				for _, candidate := range []image.Point{image.Pt(300, 220), image.Pt(400, 300), image.Pt(20, 20), image.Pt(-1, 300)} {
					name := fmt.Sprintf("visible=%d/prior=%d/ordinal=%d/at=%v", visible, priorVisible, ordinal, candidate)
					t.Run(name, func(t *testing.T) {
						o := newChatBubbleRenderOwner(t)
						a := o.create(t, 1, "First", 5, 0, 20)
						b := o.create(t, 2, "Second", 6, 0, 20)
						for _, p := range []unsafe.Pointer{a, b} {
							objectXferSetWord(p, 648, 300)
							objectXferSetWord(p, 652, 220)
							objectXferSetWord(p, 672, 40)
							objectXferSetWord(p, 676, 13)
						}
						objectXferSetWord(a, 660, priorVisible)
						objectXferSetWord(a, 680, 1)
						objectXferSetWord(b, 660, visible)
						objectXferSetWord(b, 680, ordinal)
						got := legacy.PortTestChatBubbleCandidate(b, int32(candidate.X), int32(candidate.Y))
						// Existing candidate scan gates on the moving bubble's visibility. Prior
						// visibility is considered by arrangement, but not by this helper.
						want := candidate == image.Pt(400, 300) || candidate == image.Pt(300, 220) && (visible == 0 || ordinal <= 1)
						if got != want {
							t.Fatalf("candidate %t want %t", got, want)
						}
						if objectXferGetWord(b, 648) != 300 || objectXferGetWord(b, 652) != 220 {
							t.Fatal("candidate mutated position")
						}
						legacy.PortTestChatBubbleArrange(o.c.Viewport(), b)
						moved := objectXferGetWord(b, 648) != 300 || objectXferGetWord(b, 652) != 220
						if moved != (priorVisible != 0 && ordinal > 1) {
							t.Fatal("arrangement visibility/ordinal gate")
						}
						if moved && legacy.PortTestChatBubbleOverlap(a, b) != 0 {
							t.Fatal("arrangement left overlap")
						}
						rows = append(rows, record{name, got, int32(objectXferGetWord(b, 648)), int32(objectXferGetWord(b, 652))})
					})
				}
			}
		}
	}
	spellbookCapture(t, "chat-bubble-candidate-arrange", rows, "42f57ce354f507073c34ed38f4d413e27a63fac7706b94cbdd2e1db418fdde18")
}
