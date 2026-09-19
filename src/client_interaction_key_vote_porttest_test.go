//go:build porttest

package opennox

import (
	"encoding/binary"
	"image"
	"slices"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
)

func TestClientInteractionKeyAndVoteWindows(t *testing.T) {
	o := sessionDialogResources(t)
	words, restore := legacy.PortTestClientInteractionWords()
	defer restore()
	dim := legacy.PortTestBindingDimensions()
	oldw, oldh := *dim[0], *dim[1]
	defer func() { *dim[0], *dim[1] = oldw, oldh }()
	*dim[0], *dim[1] = 801, 601
	voteMode := serverConfigOwnBytes(t, 0x5D4594, 1197304, 4)
	serverConfigOwnBytes(t, 0x5D4594, 1321220, 4)
	var sounds [][2]int
	defer legacy.PortTestClientSoundObserver(func(id, volume int) { sounds = append(sounds, [2]int{id, volume}) })()
	if interactionCall("sub_4BFC90") != 1 || interactionCall("sub_4C3390") != 1 {
		t.Fatal("key/vote constructors")
	}
	defer interactionCall("sub_4BFD10")
	defer interactionCall("sub_4C34A0")
	key := (*gui.Window)(unsafe.Pointer(uintptr(*words["dword_5d4594_1319060"])))
	vote := (*gui.Window)(unsafe.Pointer(uintptr(*words["dword_5d4594_1321216"])))
	if !key.GetFlags().IsHidden() || !vote.GetFlags().IsHidden() || vote.Off != image.Pt(751, 200) {
		t.Fatal("key/vote initial state")
	}
	type row struct {
		Old, Input, Mode uint32
		Hidden           bool
		Sounds           [][2]int
	}
	var rows []row
	for _, old := range []uint32{0, 1, 2, 0xffffffff} {
		for _, input := range []uint32{0, 1, 2, 0xffffffff} {
			*words["dword_5d4594_1319056"] = old
			sounds = nil
			key.Show()
			interactionCall("sub_4BFBB0", uintptr(input))
			want := old
			hidden := false
			if old == 0 && input == 1 {
				want = 1
			} else if old == 1 && input == 0 {
				want = 0
				hidden = true
			}
			if uint32(interactionCall("sub_4BFD30")) != want || key.GetFlags().IsHidden() != hidden {
				t.Fatal("key mode", old, input, want)
			}
			var expected [][2]int
			if old == 0 && input == 1 {
				expected = [][2]int{{1022, 100}}
				if key.Off != image.Pt(801/2-key.SizeVal.X/2, 601/2-key.SizeVal.Y/2) {
					t.Fatal("key centering")
				}
			}
			if !slices.Equal(sounds, expected) {
				t.Fatal("key sounds", old, input, sounds)
			}
			rows = append(rows, row{old, input, want, hidden, append([][2]int(nil), sounds...)})
		}
	}
	for _, input := range []uint32{0, 1, 2, 0xffffffff} {
		interactionCall("sub_48D4B0", uintptr(input))
		if binary.LittleEndian.Uint32(voteMode) != input || vote.GetFlags().IsHidden() != (input != 1) {
			t.Fatal("vote mode", input)
		}
	}
	vote.SetPos(image.Pt(25, 30))
	for _, off := range []image.Point{{0, 0}, {7, 9}, {-5, 3}} {
		vote.DrawData().ImgPtVal = off
		clear(o.pix.Pix)
		if interactionCall("sub_4C3410", uintptr(vote.C())) != 1 {
			t.Fatal("vote draw return")
		}
		got := append([]uint16(nil), o.pix.Pix...)
		clear(o.pix.Pix)
		o.c.R2().DrawImageAt(o.images[0], image.Pt(25, 30).Add(off))
		if !slices.Equal(got, o.pix.Pix) {
			t.Fatal("vote image offset", off)
		}
	}
	for _, code := range []uintptr{0, 22, 23, 16390} {
		if interactionCall("sub_4BFCD0", uintptr(key.C()), code, 0xffffffff, 0) != 0 {
			t.Fatal("key unrelated event")
		}
	}
	key.Show()
	sounds = nil
	interactionCall("sub_4BFCD0", uintptr(key.C()), 16391, uintptr(key.ChildByID(10803).C()), 0)
	if !key.GetFlags().IsHidden() || !slices.Equal(sounds, [][2]int{{766, 100}}) {
		t.Fatal("key button close", sounds)
	}
	interactionCall("sub_4BFD10")
	interactionCall("sub_4C34A0")
	if *words["dword_5d4594_1319060"] != 0 || *words["dword_5d4594_1319056"] != 0 || *words["dword_5d4594_1321216"] != 0 {
		t.Fatal("key/vote owners after close")
	}
	interactionCapture(t, "key-vote", rows)
}
