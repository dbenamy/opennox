//go:build porttest

package opennox

import (
	noxcolor "github.com/opennox/libs/color"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"testing"
	"unsafe"
)

func TestClientScreenEffectsMaidenOwner(t *testing.T) {
	oldGameplay := noxflags.GetGamePlay()
	noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0))
	t.Cleanup(func() {
		noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0))
		noxflags.SetGamePlay(oldGameplay)
	})
	o := newObjectDrawingOwner(t)
	c := o.c
	oldClient, oldServer := noxClient, noxServer
	noxClient, noxServer = c.Client, c.srv
	t.Cleanup(func() { noxClient, noxServer = oldClient, oldServer })
	t.Cleanup(c.srv.PortTestScreenNPCs(2))
	// Own the auxiliary chat/ally lookup tables consulted by the real renderer.
	for _, reg := range [][2]uintptr{{1197368, 4}, {1200916, 512}} {
		b := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, reg[0])), reg[1])
		old := append([]byte(nil), b...)
		clear(b)
		t.Cleanup(func() { copy(b, old) })
	}
	npc := c.srv.NPCs.New(7)
	for i := range npc.Color8 {
		npc.Color8[i] = uint32(noxcolor.RGB5551Color(byte(20+i*35), byte(220-i*25), byte(40+i*15)).Color32())
	}
	data, free := alloc.New(client.MonsterDrawData{})
	t.Cleanup(free)
	data.Size = 772
	ani := &data.Anim[0]
	ani.Size = 48
	ani.Cnt40 = 1
	ani.Kind = client.AnimSlave
	for i := range ani.Frames {
		ani.Frames[i] = (*noxrender.ImageHandle)(unsafe.Pointer(&o.frames[0]))
	}
	dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48))
	dr.NetCode32 = 7
	dr.DrawData = unsafe.Pointer(data)
	dr.DrawFuncPtr = legacy.PortTestScreenMaidenCallback()
	c.callbackRefs[dr.DrawFuncPtr] = 0xec000004
	c.dataRefs[uint32(uintptr(dr.DrawData))] = 0xed000001
	blank := effectsPixelHash(o.pix)
	if dr.CallDraw(c.Viewport()) != 1 || len(o.drawTrace) != 1 || effectsPixelHash(o.pix) == blank {
		t.Fatal("real maiden/monster renderer did not draw")
	}
	prior := effectsPixelHash(o.pix)
	clear(o.pix.Pix)
	npc.Color8[0] = uint32(noxcolor.RGB5551Color(255, 0, 0).Color32())
	dr.CallDraw(c.Viewport())
	if effectsPixelHash(o.pix) == prior {
		t.Fatal("NPC material did not affect maiden pixels")
	}
	c.snapshotDrawables(t)
}
