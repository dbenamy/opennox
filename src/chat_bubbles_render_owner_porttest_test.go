//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/libs/noximage"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
)

type chatBubbleRenderOwner struct {
	*objectRenderOwner
	*chatBubbleStorage
	words map[string]*uint32
}

func newChatBubbleRenderOwner(t *testing.T, extraNames ...string) *chatBubbleRenderOwner {
	t.Helper()
	o := &chatBubbleRenderOwner{objectRenderOwner: newObjectRenderOwner(t, extraNames...), chatBubbleStorage: newChatBubbleStorage(t)}
	t.Cleanup(o.c.srv.PortTestMinimapTeamColors())
	var restore func()
	o.words, restore = legacy.PortTestChatBubbleRenderGlobals()
	t.Cleanup(restore)
	oldW, oldH, oldGUI := nox_win_width, nox_win_height, nox_client_renderGUI_80828
	t.Cleanup(func() { nox_win_width, nox_win_height, nox_client_renderGUI_80828 = oldW, oldH, oldGUI })
	nox_win_width, nox_win_height, nox_client_renderGUI_80828 = 640, 480, false
	*o.words["width"], *o.words["height"], *o.words["white"] = 640, 480, 0xffff
	clear(serverConfigOwnBytes(t, 0x5D4594, 1049868, 4))
	clear(serverConfigOwnBytes(t, 0x5D4594, 1321052, 128))
	binary.LittleEndian.PutUint32(serverConfigOwnBytes(t, 0x852978, 4, 4), 0x8000)
	binary.LittleEndian.PutUint32(serverConfigOwnBytes(t, 0x85B3FC, 956, 4), 0xffff)
	oldPix := o.c.r.PixBuffer()
	t.Cleanup(func() { o.c.r.SetPixBuffer(oldPix) })
	o.pix = noximage.NewImage16(image.Rect(0, 0, 640, 480))
	o.c.r.SetPixBuffer(o.pix)
	o.c.r.Data().SetClipRect(o.pix.Rect)
	o.c.r.Data().SetClipRect2(image.Rect(0, 0, 639, 479))
	o.c.r.Data().SetRect3(o.pix.Rect)
	*o.c.Viewport() = noxrender.Viewport{Screen: image.Rect(0, 0, 640, 480), World: image.Rect(0, 0, 640, 480), Size: image.Pt(640, 480)}
	o.c.srv.SetFrame(123)
	return o
}
