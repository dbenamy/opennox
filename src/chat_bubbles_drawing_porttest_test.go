//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"image"
	"testing"
)

func TestChatBubblePlayerDrawing(t *testing.T) {
	type record struct {
		Name, Pixels string
		Arrow        uint32
	}
	var rows []record
	hashes := make(map[string]string)
	for _, name := range []string{"", "Alice"} {
		for _, team := range []server.TeamID{0, 1, 2} {
			for _, clipped := range []bool{false, true} {
				for _, multiple := range []bool{false, true} {
					key := fmt.Sprintf("name=%s/team=%d/clipped=%t/multiple=%t", name, team, clipped, multiple)
					t.Run(key, func(t *testing.T) {
						o := newChatBubbleRenderOwner(t)
						restore := noxflags.PortTestGameFlags(0)
						defer restore()
						o.players[0].SetName(name)
						pos := image.Pt(300, 300)
						if clipped {
							pos.X = 0
						}
						dr := o.drawable(7, pos)
						dr.ObjClass = 4
						dr.TeamVal.ID = team
						p := o.create(t, 7, "Player speech", 13, 0, 20)
						if multiple {
							o.create(t, 1234, "Nearby speech", 13, uint32(300)|uint32(300)<<16, 20)
						}
						blank := effectsPixelHash(o.pix)
						legacy.PortTestChatBubbleDraw(o.c.Viewport())
						hash := effectsPixelHash(o.pix)
						if hash == blank {
							t.Fatal("player speech did not draw")
						}
						wantArrow := uint32(1)
						if clipped {
							wantArrow = 0
						}
						if objectXferGetWord(p, 664) != wantArrow {
							t.Fatal("clipped arrow")
						}
						hashes[key] = hash
						rows = append(rows, record{key, hash, wantArrow})
					})
				}
			}
		}
	}
	for _, clipped := range []bool{false, true} {
		for _, multiple := range []bool{false, true} {
			key := fmt.Sprintf("team=0/clipped=%t/multiple=%t", clipped, multiple)
			if hashes["name=/"+key] == hashes["name=Alice/"+key] {
				t.Fatal("player name did not draw")
			}
		}
	}
	spellbookCapture(t, "chat-bubble-player-drawing", rows, "695e0ecafe1411dcef50aa3ad388d11ff6164f8097cc452b0c21bca1db0098c0")
}

func TestChatBubbleEmptyDrawing(t *testing.T) {
	o := newChatBubbleRenderOwner(t)
	before := effectsPixelHash(o.pix)
	legacy.PortTestChatBubbleDraw(o.c.Viewport())
	if effectsPixelHash(o.pix) != before {
		t.Fatal("empty list drew pixels")
	}
	o.create(t, 1234, "Hidden", 6, 0, 20)
	restore := noxflags.PortTestGameFlags(2048)
	defer restore()
	legacy.PortTestChatBubbleDraw(o.c.Viewport())
	if effectsPixelHash(o.pix) != before {
		t.Fatal("hidden bubble drew pixels")
	}
	legacy.PortTestChatBubbleClear()
	legacy.PortTestChatBubbleDraw(o.c.Viewport())
	if effectsPixelHash(o.pix) != before {
		t.Fatal("cleared list drew pixels")
	}
	spellbookCapture(t, "chat-bubble-empty-drawing", []string{before}, "931975f4bcde125d552e7274894ff359085770a9442ef583bc949abc91f7fa05")
}
