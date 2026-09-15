//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"image"
	"testing"
	"unsafe"
)

type objectRenderOwner struct {
	*objectDrawingOwner
	renderEnv *legacy.PortTestObjectRenderEnvironment
	players   []server.Player
	shiny     *legacy.ImageRef
	animation *legacy.ImageRefAnim
}

func newObjectRenderOwner(t *testing.T) *objectRenderOwner {
	t.Helper()
	o := newObjectDrawingOwner(t, "Ghost")
	e := legacy.PortTestNewObjectRenderEnvironment(o.c.r.Data())
	t.Cleanup(e.Restore)
	players, freePlayers := o.c.srv.PortTestObjectRenderPlayers()
	t.Cleanup(freePlayers)
	oldClient, oldServer := noxClient, noxServer
	noxClient, noxServer = o.c.Client, o.c.srv
	t.Cleanup(func() { noxClient, noxServer = oldClient, oldServer })
	flags := noxflags.GetGamePlay()
	noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0))
	t.Cleanup(func() { noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0)); noxflags.SetGamePlay(flags) })
	width, height := nox_win_width, nox_win_height
	nox_win_width, nox_win_height = 96, 96
	t.Cleanup(func() { nox_win_width, nox_win_height = width, height })
	ref, freeRef := alloc.New(legacy.ImageRef{})
	t.Cleanup(freeRef)
	anim, freeAnim := alloc.New(legacy.ImageRefAnim{})
	t.Cleanup(freeAnim)
	ref.SetName("ShinySpot")
	ref.RefKind = 2
	ref.Field_24 = unsafe.Pointer(anim)
	anim.ImagesPtr = (*noxrender.ImageHandle)(unsafe.Pointer(&o.frames[0]))
	anim.ImagesSz = 4
	oldImages := nox_images_arr1_787156
	nox_images_arr1_787156 = []*legacy.ImageRef{ref}
	t.Cleanup(func() { nox_images_arr1_787156 = oldImages })
	r := &objectRenderOwner{objectDrawingOwner: o, renderEnv: e, players: players, shiny: ref, animation: anim}
	t.Cleanup(func() {
		for dr := o.c.Objs.List1; dr != nil; dr = dr.NextPtr {
			dr.TeamVal.ID = 0
		}
	})
	r.resetRender(1, 120)
	return r
}
func (o *objectRenderOwner) resetRender(seed, frame uint32) {
	// Team IDs are read-only comparison inputs, not membership registrations.
	for dr := o.c.Objs.List1; dr != nil; dr = dr.NextPtr {
		dr.TeamVal.ID = 0
	}
	o.reset(seed, frame)
	o.renderEnv.Reset()
	o.c.r.Set_dword_5d4594_3799484(0)
	o.c.r.Reset_dword_5d4594_3799476()
	clear(o.players)
	for i := 0; i < 2; i++ {
		o.players[i].Active = 1
		o.players[i].PlayerInd = byte(i)
		o.players[i].NetCodeVal = uint32(7 + i)
	}
	legacy.Set_dword_8531A0_2576(&o.players[0])
	for i := range o.c.tiles.nox_arr_956A00 {
		o.c.tiles.nox_arr_956A00[i] = 0
		clear(o.c.tiles.nox_arr_957820[i].arr[:])
	}
	*memmap.PtrUint32(0x85B3FC, 956) = 0x8c5f
	o.animation.ImagesSz = 4
	o.animation.ImagesPtr = (*noxrender.ImageHandle)(unsafe.Pointer(&o.frames[0]))
}
func (o *objectRenderOwner) drawable(id int, pos image.Point) *client.Drawable {
	dr := o.c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, pos)
	dr.NetCode32 = uint32(id)
	return dr
}

type objectRenderResult struct {
	Draw        objectDrawingResult
	Globals     []uint32
	Last        uint32
	Crop        uint32
	Bottom      int
	ShinyCached bool
}

func (o *objectRenderOwner) renderResult(t *testing.T, id, step, ret int) objectRenderResult {
	t.Helper()
	last := o.renderEnv.LastDrawable()
	var lastRef uint32
	if last != nil {
		var ok bool
		lastRef, ok = o.c.refs[last]
		if !ok {
			t.Fatal("last drawable cache points outside owner")
		}
	}
	cache := *memmap.PtrPtr(0x5D4594, 1321524)
	if cache != nil && cache != unsafe.Pointer(o.shiny) {
		t.Fatal("shiny cache points outside owned reference")
	}
	return objectRenderResult{o.result(t, id, step, ret), o.renderEnv.State(), lastRef, o.c.r.PortTestObjectRenderCrop(), o.c.r.Get_dword_5d4594_3799476(), cache != nil}
}
