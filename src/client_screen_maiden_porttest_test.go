//go:build porttest

package opennox

import (
	noxcolor "github.com/opennox/libs/color"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"image"
	"testing"
	"unsafe"
)

func TestClientScreenEffectsMaiden(t *testing.T) {
	oldGameplay := noxflags.GetGamePlay()
	noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0))
	t.Cleanup(func() {
		noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0))
		noxflags.SetGamePlay(oldGameplay)
	})
	o := newObjectDrawingOwner(t)
	c := o.c
	env := legacy.PortTestNewScreenEnvironment()
	defer env.Restore()
	oldClient, oldServer := noxClient, noxServer
	noxClient, noxServer = c.Client, c.srv
	defer func() { noxClient, noxServer = oldClient, oldServer }()
	t.Cleanup(c.srv.PortTestScreenNPCs(2))
	npc := c.srv.NPCs.New(7)
	data, freeData := alloc.New(client.MonsterDrawData{})
	defer freeData()
	data.Size = 772
	for i := range data.Anim {
		ani := &data.Anim[i]
		ani.Size = 48
		ani.Cnt40 = 32
		ani.Val42 = 1
		for dir := range ani.Frames {
			ani.Frames[dir] = (*noxrender.ImageHandle)(unsafe.Pointer(&o.frames[(dir%5)*32]))
		}
	}
	c.dataRefs[uint32(uintptr(unsafe.Pointer(data)))] = 0xed000001
	callback := legacy.PortTestScreenMaidenCallback()
	c.callbackRefs[callback] = 0xec000004
	first, freeFirst := alloc.New(server.Object{})
	defer freeFirst()
	second, freeSecond := alloc.New(server.Object{})
	defer freeSecond()
	update, freeUpdate := alloc.Make([]byte{}, 2200)
	defer freeUpdate()
	first.UpdateData, second.UpdateData = unsafe.Pointer(&update[0]), unsafe.Pointer(&update[0])
	defer func() { c.srv.Objs.List = nil }()
	type result struct {
		Mode, Palette, Anim, Direction, Kind, Step int
		Draw                                       objectDrawingResult
		Globals                                    []uint32
	}
	var out []result
	id := 0
	for mode := 0; mode < 6; mode++ {
		for palette := 0; palette < 2; palette++ {
			for _, anim := range []int{0, 1, 8, 15} {
				for _, dir := range []int{0, 1, 3, 5, 7, 8} {
					for _, kind := range []client.AnimKind{client.AnimLoop, client.AnimSlave} {
						id++
						o.reset(uint32(id), 120)
						env.Reset()
						c.srv.Objs.List = nil
						for i := 0; i < 6; i++ {
							r, g, b := byte(20+i*35+palette*17), byte(220-i*25-palette*21), byte(40+i*15+palette*29)
							npc.Color8[i] = uint32(noxcolor.RGB5551Color(r, g, b).Color32())
							copy(update[2076+3*i:], []byte{b, r, g})
						}
						npc.IDVal = 7
						if mode == 1 {
							npc.IDVal = 99
						}
						first.Extent, second.Extent = 99, 99
						first.ObjNext = nil
						second.ObjPrev = nil
						if mode >= 2 {
							noxflags.SetGame(noxflags.GameFlag22)
						}
						if mode >= 3 {
							c.srv.Objs.List = first
						}
						if mode == 4 {
							first.Extent = 7
						}
						if mode == 5 {
							first.ObjNext = second
							second.ObjPrev = first
							second.Extent = 7
						}
						data.Anim[anim].Kind = kind
						dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48))
						dr.DrawData = unsafe.Pointer(data)
						dr.DrawFuncPtr = callback
						dr.NetCode32 = 7
						dr.AnimInd = uint32(anim)
						dr.AnimDir = byte(dir)
						for step := 0; step < 3; step++ {
							c.srv.SetFrame(uint32(120 + step*17))
							dr.AnimFrameSlave = uint32(step * 15)
							clear(o.pix.Pix)
							o.drawTrace = nil
							ret := dr.CallDraw(c.Viewport())
							if (mode == 1 && len(o.drawTrace) != 0) || (mode != 1 && len(o.drawTrace) != 1) {
								t.Fatal("maiden dispatch/ownership contract")
							}
							out = append(out, result{mode, palette, anim, dir, int(kind), step, o.result(t, id, step, ret), env.State()})
						}
					}
				}
			}
		}
	}
	effectsCapture(t, "screen-effects-maiden", out, len(out), "ad6a3815de864ed82b56a1b1c1f868705babe7f483a7bbade2fb14e77158701f")
}
