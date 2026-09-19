//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestChatBubbleHUDPlacement(t *testing.T) {
	type fixture struct {
		name string
		rect [4]int32
		guiY int32
	}
	cases := []fixture{
		{"clear", [4]int32{300, 300, 330, 320}, 300},
		{"top", [4]int32{20, 20, 40, 40}, 55},
		{"left", [4]int32{20, 370, 60, 400}, 323},
		{"middle", [4]int32{300, 420, 340, 440}, 386},
		{"right", [4]int32{570, 300, 600, 330}, 249},
		{"summons", [4]int32{580, 80, 600, 100}, 80},
		{"outside-left", [4]int32{-1, 300, 30, 320}, 300},
		{"outside-bottom", [4]int32{300, 481, 330, 500}, 481},
		{"left-inclusive-edge", [4]int32{111, 353, 130, 375}, 331},
	}
	type record struct {
		Name  string
		Rect  [4]int32
		Arrow uint32
	}
	var rows []record
	for _, gui := range []bool{false, true} {
		for inventory := uint8(0); inventory < 4; inventory++ {
			for _, summon := range []uint32{0, 1} {
				for _, tc := range cases {
					name := fmt.Sprintf("gui=%t/inventory=%d/summon=%d/%s", gui, inventory, summon, tc.name)
					t.Run(name, func(t *testing.T) {
						_ = newChatBubbleRenderOwner(t)
						nox_client_renderGUI_80828 = gui
						*memmap.PtrUint8(0x5D4594, 1049868) = inventory
						*memmap.PtrUint32(0x5D4594, 1321060) = summon
						got, arrow := tc.rect, uint32(1)
						want, wantArrow := tc.rect, uint32(1)
						if gui {
							want[1] = tc.guiY
						}
						switch tc.name {
						case "top":
							want[1] = 55
							if inventory == 1 || inventory == 2 {
								want[1] = 279
							}
							wantArrow = 0
						case "outside-left", "outside-bottom":
							wantArrow = 0
						case "left", "middle", "right", "left-inclusive-edge":
							if gui {
								wantArrow = 0
							}
						case "summons":
							if gui && summon != 0 {
								want[1], wantArrow = 145, 0
							}
						}
						legacy.PortTestChatBubblePlace(&got, &arrow)
						if got != want || arrow != wantArrow {
							t.Fatalf("rect/arrow %v/%d want %v/%d", got, arrow, want, wantArrow)
						}
						rows = append(rows, record{name, got, arrow})
					})
				}
			}
		}
	}
	spellbookCapture(t, "chat-bubble-hud-placement", rows, "5da8c1f3b8929f9910f672c358d2fab277c2c8cb1bb40ff109caa658fc0995c4")
}
