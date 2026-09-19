//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestChatBubbleRemoveAppend(t *testing.T) {
	type record struct {
		Name                     string
		AfterRemove, AfterAppend []chatBubbleRecord
	}
	var rows []record
	for _, remove := range []uint32{0, 1, 2, 3} {
		t.Run(fmt.Sprint(remove), func(t *testing.T) {
			o := newChatBubbleOwner(t)
			for code := uint16(1); code <= 3; code++ {
				o.create(t, code, fmt.Sprint("Bubble", code), 7, uint32(code)*100, 20)
			}
			legacy.PortTestChatBubbleRemove(remove)
			r := record{Name: fmt.Sprint(remove), AfterRemove: o.snapshot(t)}
			expected := 3
			if remove != 0 {
				expected--
			}
			if len(r.AfterRemove) != expected {
				t.Fatal("removed bubble count")
			}
			if remove != 0 && legacy.PortTestChatBubbleLookup(remove) != nil {
				t.Fatal("removed bubble remains")
			}
			o.create(t, 4, "Appended", 8, 400, 30)
			r.AfterAppend = o.snapshot(t)
			if len(r.AfterAppend) != expected+1 || r.AfterAppend[len(r.AfterAppend)-1].Code != 4 {
				t.Fatal("append after removal")
			}
			rows = append(rows, r)
		})
	}
	spellbookCapture(t, "chat-bubble-remove-append", rows, "07be6a642f156cd485e4e6240940dd0d549d8953b7cb64aada53a4db5c397fa5")
}
func TestChatBubbleCreateReplace(t *testing.T) {
	type record struct {
		Name              string
		Initial, Replaced []chatBubbleRecord
	}
	var rows []record
	for _, frame := range []uint32{0, 123, 0xfffffff0} {
		for _, length := range []byte{0, 7, 8, 63, 64, 255} {
			for _, duration := range []uint16{0, 1, 65535} {
				name := fmt.Sprintf("frame=%x/length=%d/duration=%d", frame, length, duration)
				t.Run(name, func(t *testing.T) {
					o := newChatBubbleOwner(t)
					o.s.SetFrame(frame)
					first := o.create(t, 65535, "Hello, 世界", length, 0xfffe8001, duration)
					o.create(t, 1, "Second", 1, 0, 1)
					r := record{Name: name, Initial: o.snapshot(t)}
					ttl := uint32(duration)
					if ttl == 0 {
						ticks := uint32(length) / 8
						if ticks > 8 {
							ticks = 8
						}
						ttl = o.s.TickRate() * (ticks + 2)
					}
					got := r.Initial[0]
					if got.Code != 65535 || got.Text != "Hello, 世界" || got.Length != uint32(length) || got.Pos != 0xfffe8001 || got.Expiry != frame+ttl {
						t.Fatalf("create fields %+v", got)
					}
					again := o.create(t, 65535, "短", 2, 0x00020003, 9)
					if again != first {
						t.Fatal("replacement allocated a second bubble")
					}
					r.Replaced = o.snapshot(t)
					if len(r.Replaced) != 2 || r.Replaced[0].Text != "短" || r.Replaced[0].Expiry != frame+9 || r.Replaced[0].Pos != 0x00020003 || r.Replaced[1] != r.Initial[1] {
						t.Fatal("replacement order/state")
					}
					legacy.PortTestChatBubbleClear()
					if *o.head != nil {
						t.Fatal("clear head")
					}
					o.create(t, 5, "After clear", 11, 0, 1)
					if len(o.snapshot(t)) != 1 {
						t.Fatal("create after clear")
					}
					legacy.PortTestChatBubbleDestroy()
					if *o.head != nil || legacy.Get_nox_alloc_chat_1197364() != nil {
						t.Fatal("destroy pool/head")
					}
					rows = append(rows, r)
				})
			}
		}
	}
	spellbookCapture(t, "chat-bubble-create-replace", rows, "acd8e05d2eb585e9f2076388a920cd4b60d23b06130f3302274c721288f8efcc")
}
