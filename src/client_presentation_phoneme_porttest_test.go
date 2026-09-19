//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"testing"
)

type clientPresentationPhonemes struct {
	*chatBubbleRenderOwner
	images []*noxrender.Image
	calls  []string
	names  []string
}

func newClientPresentationPhonemes(t *testing.T, fail int) *clientPresentationPhonemes {
	t.Helper()
	o := &clientPresentationPhonemes{chatBubbleRenderOwner: newChatBubbleRenderOwner(t)}
	oldGUI := o.c.GUI
	o.c.GUI = gui.New(o.c.Render())
	t.Cleanup(func() {
		o.c.GUI.DestroyAll()
		o.c.GUI.FreeDestroyed()
		o.c.GUI = oldGUI
	})
	for _, r := range blobdata.PortTestClientPresentationTables() {
		copy(serverConfigOwnBytes(t, r.Base, r.Offset, len(r.Data)), r.Data)
	}
	clear(serverConfigOwnBytes(t, 0x5D4594, 1096564, 64))
	clear(serverConfigOwnBytes(t, 0x587000, 151272, 32))
	raw := make([][]byte, 8)
	for i := range raw {
		raw[i] = spriteAnimationTestImage(i)
		*memmap.PtrPtr(0x587000, 151272+4*uintptr(i)) = memmap.PtrOff(0x587000, 151304+16*uintptr(i))
		o.names = append(o.names, alloc.GoString(memmap.PtrUint8(0x587000, 151304+16*uintptr(i))))
	}
	var free func()
	o.images, free = o.c.r.GetBag().PortTestSpriteImages(raw)
	t.Cleanup(free)
	old := legacy.Nox_xxx_gLoadImg
	t.Cleanup(func() { legacy.Nox_xxx_gLoadImg = old })
	legacy.Nox_xxx_gLoadImg = func(name string) *noxrender.Image {
		o.calls = append(o.calls, name)
		for i, n := range o.names {
			if n == name {
				if i == fail {
					return nil
				}
				return o.images[i]
			}
		}
		t.Fatalf("unowned phoneme image %q", name)
		return nil
	}
	return o
}
func TestClientPresentationPhonemeInit(t *testing.T) {
	type record struct {
		Fail           int
		Calls          []string
		Loaded         [8]bool
		Window         bool
		Position, Size image.Point
	}
	var rows []record
	for fail := -1; fail < 8; fail++ {
		t.Run(fmt.Sprint(fail), func(t *testing.T) {
			o := newClientPresentationPhonemes(t, fail)
			w := (*gui.Window)(legacy.PortTestPresentationPhonemeInit())
			if w != nil {
				t.Cleanup(w.Destroy)
			}
			wantCalls := 8
			if fail >= 0 {
				wantCalls = fail + 1
			}
			if len(o.calls) != wantCalls || (w != nil) != (fail < 0) {
				t.Fatal("phoneme initialization/failure boundary")
			}
			var loaded [8]bool
			for i := range loaded {
				p := *memmap.PtrPtr(0x5D4594, 1096564+4*uintptr(i))
				loaded[i] = p != nil
				want := fail < 0 || i < fail
				if loaded[i] != want {
					t.Fatal("loaded prefix")
				}
				if want && noxrender.ImageHandle(p) != o.images[i].C() {
					t.Fatal("wrong image")
				}
			}
			var pos, size image.Point
			if w != nil {
				pos, size = w.GlobalPos(), w.Size()
				if pos != image.Pt(270, 190) || size != image.Pt(1, 1) {
					t.Fatal("phoneme window placement")
				}
			}
			rows = append(rows, record{fail, append([]string(nil), o.calls...), loaded, w != nil, pos, size})
		})
	}
	spellbookCapture(t, "client-presentation-phoneme-init", rows, "eb549fa8fb07fcd7ededc62cf78655230cfbf73c8206ef01106356c8bbe66521")
}
func TestClientPresentationPhonemeFrames(t *testing.T) {
	o := newClientPresentationPhonemes(t, -1)
	w := (*gui.Window)(legacy.PortTestPresentationPhonemeInit())
	if w == nil {
		t.Fatal("phoneme window")
	}
	t.Cleanup(w.Destroy)
	type record struct {
		Index             int
		Frame, Age, Stamp uint32
		Pixels, Next      string
	}
	var rows []record
	for index := 0; index < 8; index++ {
		for _, frame := range []uint32{0, 1, 100, 0xffffffff} {
			for _, age := range []uint32{0, 1, 3, 4, 7} {
				for i := 0; i < 8; i++ {
					*memmap.PtrUint32(0x5D4594, 1096596+4*uintptr(i)) = 0
				}
				stamp := frame - age
				o.c.srv.SetFrame(stamp)
				legacy.PortTestPresentationPhonemeMark(index)
				p := memmap.PtrUint32(0x5D4594, 1096596+4*uintptr(index))
				if *p != stamp {
					t.Fatal("mark frame")
				}
				o.c.srv.SetFrame(frame)
				clear(o.pix.Pix)
				blank := effectsPixelHash(o.pix)
				w.Draw()
				actual := effectsPixelHash(o.pix)
				after := stamp
				if age > 3 {
					after = 0
				}
				if *p != after {
					t.Fatal("expiry boundary")
				}
				clear(o.pix.Pix)
				if stamp != 0 {
					pos := image.Pt(320+int(memmap.Int32(0x587000, 151208+8*uintptr(index)))-16, 240+int(memmap.Int32(0x587000, 151212+8*uintptr(index)))-41)
					o.c.r.DrawImageAt(o.images[index], pos)
				}
				if actual != effectsPixelHash(o.pix) || (stamp != 0 && actual == blank) {
					t.Fatalf("icon placement/pixels index=%d frame=%d age=%d", index, frame, age)
				}
				clear(o.pix.Pix)
				w.Draw()
				next := effectsPixelHash(o.pix)
				if after == 0 && next != blank {
					t.Fatal("expired icon still draws")
				}
				rows = append(rows, record{index, frame, age, *p, actual, next})
			}
		}
	}
	spellbookCapture(t, "client-presentation-phoneme-frames", rows, "aa5b01bccabf35ac0a89c2ccd1e59e756788f895ac77902bac04d35a00ce65c5")
}
