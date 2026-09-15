//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/noxrender"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"testing"
	"unsafe"
)

func TestClientObjectRenderOwner(t *testing.T) {
	o := newObjectDrawingOwner(t)
	c := o.c
	env := legacy.PortTestNewObjectRenderEnvironment(c.r.Data())
	defer env.Restore()
	oldClient, oldServer := noxClient, noxServer
	noxClient, noxServer = c.Client, c.srv
	defer func() { noxClient, noxServer = oldClient, oldServer }()
	oldGameplay := noxflags.GetGamePlay()
	noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0))
	defer func() { noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0)); noxflags.SetGamePlay(oldGameplay) }()
	players, freePlayers := c.srv.PortTestObjectRenderPlayers()
	defer freePlayers()
	local := &players[0]
	remote := &players[1]
	local.Active, local.NetCodeVal = 1, 7
	remote.Active, remote.NetCodeVal = 1, 8
	legacy.Set_dword_8531A0_2576(local)
	dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48))
	dr.ObjClass = 4
	dr.NetCode32 = 8
	remote.Field3680 = 1
	blank := effectsPixelHash(o.pix)
	legacy.Nox_xxx_drawObject_4C4770_draw(c.Viewport(), dr, noxrender.ImageHandle(o.images[0].C()))
	if len(o.drawTrace) != 0 || effectsPixelHash(o.pix) != blank {
		t.Fatal("remote observer was rendered")
	}
	remote.Field3680 = 0
	legacy.Nox_xxx_drawObject_4C4770_draw(c.Viewport(), dr, noxrender.ImageHandle(o.images[0].C()))
	if len(o.drawTrace) != 1 || effectsPixelHash(o.pix) == blank {
		t.Fatal("visible remote player did not render")
	}
	// Repeated below-floor drawing must use the real C save/restore path.
	dr.ObjClass = 0
	dr.ZVal = 0xffff
	legacy.Nox_xxx_drawObject_4C4770_draw(c.Viewport(), dr, noxrender.ImageHandle(o.images[1].C()))
	clip, clip2 := c.r.Data().ClipRect(), c.r.Data().ClipRect2()
	legacy.Nox_xxx_drawObject_4C4770_draw(c.Viewport(), dr, noxrender.ImageHandle(o.images[2].C()))
	if c.r.Data().ClipRect() != clip || c.r.Data().ClipRect2() != clip2 {
		t.Fatal("below-floor draw failed to restore clipping")
	}
	if env.LastDrawable() != dr {
		t.Fatal("last-drawable cache does not reference actual owner")
	}
	c.snapshotDrawables(t)
}
func TestClientObjectRenderShinyOwner(t *testing.T) {
	o := newObjectDrawingOwner(t)
	c := o.c
	env := legacy.PortTestNewObjectRenderEnvironment(c.r.Data())
	defer env.Restore()
	ref, freeRef := alloc.New(legacy.ImageRef{})
	defer freeRef()
	anim, freeAnim := alloc.New(legacy.ImageRefAnim{})
	defer freeAnim()
	ref.SetName("ShinySpot")
	ref.RefKind = 2
	ref.Field_24 = unsafe.Pointer(anim)
	anim.ImagesPtr = (*noxrender.ImageHandle)(unsafe.Pointer(&o.frames[0]))
	anim.ImagesSz = 4
	old := nox_images_arr1_787156
	nox_images_arr1_787156 = []*legacy.ImageRef{ref}
	defer func() { nox_images_arr1_787156 = old }()
	dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(96, 96))
	dr.NetCode32 = 0
	c.srv.SetFrame(0)
	blank := effectsPixelHash(o.pix)
	if legacy.PortTestObjectRenderShiny(c.Viewport(), dr) != 0 || len(o.drawTrace) != 1 || effectsPixelHash(o.pix) == blank {
		t.Fatal("real named shiny animation did not render")
	}
	nox_images_arr1_787156 = nil // Subsequent calls must use the actual cached reference.
	c.srv.SetFrame(8)
	if legacy.PortTestObjectRenderShiny(c.Viewport(), dr) != 0 || len(o.drawTrace) != 1 {
		t.Fatal("shiny off-period rendered a frame")
	}
	c.srv.SetFrame(32)
	if legacy.PortTestObjectRenderShiny(c.Viewport(), dr) != 1 || len(o.drawTrace) != 2 {
		t.Fatal("cached shiny reference was not reused")
	}
	c.snapshotDrawables(t)
}

func TestClientObjectRenderGhostOwner(t *testing.T) {
	o := newObjectRenderOwner(t)
	typ := o.c.Things.IndByID("Ghost")
	if typ == 0 {
		t.Fatal("real Ghost type missing")
	}
	o.renderEnv.GhostType(0)
	dr := o.c.Nox_xxx_spriteLoadAdd_45A360_drawable(typ, image.Pt(48, 48))
	blank := effectsPixelHash(o.pix)
	legacy.Nox_xxx_drawObject_4C4770_draw(o.c.Viewport(), dr, noxrender.ImageHandle(o.images[0].C()))
	if o.renderEnv.State()[0] != uint32(typ) || len(o.drawTrace) != 1 || effectsPixelHash(o.pix) == blank {
		t.Fatal("lazy ghost lookup or rendering failed")
	}
	if legacy.PortTestObjectRenderGhost(o.c.Viewport(), dr) != 128 {
		t.Fatal("center ghost opacity")
	}
	o.c.snapshotDrawables(t)
}
