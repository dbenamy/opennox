//go:build porttest

package opennox

import (
	"encoding/binary"
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noximage"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"github.com/opennox/opennox/v1/server"
	"golang.org/x/image/font/basicfont"
)

type objectDrawingOwner struct {
	c             *effectsTestClient
	pix           *noximage.Image16
	effects       *legacy.PortTestEffectsEnvironment
	env           *legacy.PortTestObjectDrawEnvironment
	images        []*noxrender.Image
	data, frames  []uint32
	weapon, armor *server.Modifier
	mods          []*server.ModifierEff
	proxy         *objectDrawingClient
	unlinked      []*client.Drawable
	rawDeleted    []uint32
	namedCalls    []string
	drawTrace     []objectImageDraw
	modRefs       map[uint32]uint32
}

func objectMaterialImage(frame int) []byte {
	const width, height = 12, 6
	b := make([]byte, 17)
	binary.LittleEndian.PutUint32(b, width)
	binary.LittleEndian.PutUint32(b[4:], height)
	binary.LittleEndian.PutUint32(b[8:], uint32(frame%8))
	binary.LittleEndian.PutUint32(b[12:], uint32(frame/8))
	for y := 0; y < height; y++ {
		// Each row uses a different actual renderer material, with varied intensity.
		b = append(b, byte((y+1)<<4|4), width)
		for x := 0; x < width; x++ {
			b = append(b, byte(80+(x*13+frame*7)%176))
		}
	}
	return b
}
func newObjectDrawingOwner(t *testing.T) *objectDrawingOwner {
	t.Helper()
	t.Cleanup(handles.PortTestInit())
	c, pix, effects := newEffectsFullOwner(t, "ArrowTailLink", "WeakArrowTailLink")
	t.Cleanup(legacy.PortTestSpriteAnimationEnvironment())
	env := legacy.PortTestNewObjectDrawEnvironment()
	t.Cleanup(env.Restore)
	flags := noxflags.GetGame()
	engine := noxflags.GetEngine()
	t.Cleanup(func() {
		noxflags.ResetGame()
		noxflags.SetGame(flags)
		noxflags.ResetEngine()
		noxflags.SetEngine(engine)
	})
	light, freeLight := alloc.Make([]uint32{}, 6)
	t.Cleanup(freeLight)
	c.tiles.lightsOutBuf = light
	for i := range light {
		light[i] = 255
	}
	for x := range c.tiles.nox_arr2_853BC0 {
		for y := range c.tiles.nox_arr2_853BC0[x] {
			c.tiles.nox_arr2_853BC0[x][y] = noxrender.RGB{R: 255 << 16, G: 255 << 16, B: 255 << 16}
		}
	}
	raw := make([][]byte, 32)
	for i := range raw {
		raw[i] = objectMaterialImage(i)
	}
	images, freeImages := c.r.GetBag().PortTestSpriteImages(raw)
	t.Cleanup(freeImages)
	t.Cleanup(c.r.GetFonts().PortTestDefaultFont(basicfont.Face7x13))
	c.imageRefs = make(map[uint32]uint32)
	for i, img := range images {
		c.imageRefs[uint32(uintptr(img.C()))] = 0xe8000000 + uint32(i)
	}
	data, freeData := alloc.Make([]uint32{}, 14)
	t.Cleanup(freeData)
	frames, freeFrames := alloc.Make([]uint32{}, 160)
	t.Cleanup(freeFrames)
	for i := range frames {
		frames[i] = uint32(uintptr(images[(i+i/32)%32].C()))
	}
	c.dataRefs = map[uint32]uint32{uint32(uintptr(unsafe.Pointer(&data[0]))): 0xe9000001}
	c.callbackRefs = make(map[unsafe.Pointer]uint32)
	for op := 0; op < 18; op++ {
		c.callbackRefs[legacy.PortTestObjectDrawCallback(op)] = 0xea000000 + uint32(op)
	}
	for op := 0; op < 7; op++ {
		c.callbackRefs[legacy.PortTestSpriteAnimationCallback(op)] = 0xe5000000 + uint32(op)
	}
	weapon, freeWeapon := alloc.New(server.Modifier{})
	armor, freeArmor := alloc.New(server.Modifier{})
	c.srv.Modif.Dword_5d4594_251600 = weapon
	c.srv.Modif.Dword_5d4594_251608 = armor
	t.Cleanup(func() {
		c.srv.Modif.Dword_5d4594_251600 = nil
		c.srv.Modif.Dword_5d4594_251608 = nil
		freeArmor()
		freeWeapon()
	})
	weapon.TypeInd, armor.TypeInd = 4, 4
	for i := 0; i < 8; i++ {
		weapon.Colors12[i] = types.RGB{R: byte(30 + i*25), G: byte(220 - i*19), B: byte(80 + i*13)}
		armor.Colors12[i] = types.RGB{R: byte(210 - i*19), G: byte(20 + i*23), B: byte(50 + i*27)}
	}
	*weapon.ColorIndexes() = [4]int32{1, 3, 5, 6}
	*armor.ColorIndexes() = [4]int32{2, 4, 6, 1}
	mods := make([]*server.ModifierEff, 4)
	for i := range mods {
		mod, freeMod := alloc.New(server.ModifierEff{})
		t.Cleanup(freeMod)
		mod.Color24 = types.RGB{R: byte(240 - i*47), G: byte(30 + i*60), B: byte(150 - i*19)}
		mods[i] = mod
	}
	teams, freeTeams := alloc.Make([]server.Team{}, 9)
	t.Cleanup(freeTeams)
	c.srv.Teams.Arr = teams
	noxflags.ResetGame()
	noxflags.SetGame(noxflags.GameModeCoopTeam)
	for i, name := range []string{"Red team", "Blue team"} {
		c.srv.Teams.Create(server.TeamID(i+1)).SetNameAnd68(name, 0)
	}
	o := &objectDrawingOwner{c: c, pix: pix, effects: effects, env: env, images: images, data: data, frames: frames, weapon: weapon, armor: armor, mods: mods}
	o.installOwnership(t)
	o.reset(1, 120)
	return o
}
func (o *objectDrawingOwner) reset(seed, frame uint32) {
	for len(o.unlinked) > 0 {
		o.proxy.Nox_xxx_spriteDelete_45A4B0(o.unlinked[0])
	}
	o.rawDeleted = nil
	o.namedCalls = nil
	o.drawTrace = nil
	o.c.resetCase(o.effects, o.pix, seed, frame)
	o.env.Reset()
	noxflags.ResetGame()
	noxflags.UnsetEngine(noxflags.EngineNoSoftLights)
	o.c.srv.SetTickRate(60)
	for i := range o.c.tiles.lightsOutBuf {
		o.c.tiles.lightsOutBuf[i] = 255
	}
	for x := range o.c.tiles.nox_arr2_853BC0 {
		for y := range o.c.tiles.nox_arr2_853BC0[x] {
			o.c.tiles.nox_arr2_853BC0[x][y] = noxrender.RGB{R: 255 << 16, G: 255 << 16, B: 255 << 16}
		}
	}
	*o.c.Viewport() = noxrender.Viewport{Screen: o.pix.Rect, World: o.pix.Rect, Size: o.pix.Rect.Size()}

	for i := 0; i < 16; i++ {
		o.c.r.Data().SetMaterialRGB(i, 255, 255, 255)
	}
	clear(o.data)
	o.data[0] = 16
	o.data[1] = uint32(uintptr(unsafe.Pointer(&o.frames[0])))
	b := unsafe.Slice((*byte)(unsafe.Pointer(&o.data[0])), 56)
	b[8], b[9] = 32, 1
	o.data[3] = 2
}

func TestClientObjectDrawingOwner(t *testing.T) {
	o := newObjectDrawingOwner(t)
	c := o.c
	dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48))
	words := unsafe.Slice((*uint32)(dr.C()), 128)
	words[76] = uint32(uintptr(unsafe.Pointer(&o.data[0])))
	o.data[0], o.data[1] = 8, o.frames[0]
	dr.DrawFuncPtr = legacy.PortTestObjectDrawCallback(7)
	// No fixed-light flag: use the real world-light interpolation path.
	blank := effectsPixelHash(o.pix)
	if got := dr.CallDraw(c.Viewport()); got != 1 {
		t.Fatalf("weapon draw return %d", got)
	}
	first := effectsPixelHash(o.pix)
	if first == blank {
		t.Fatal("material image did not render")
	}
	if words[2] != o.frames[0] {
		t.Fatal("weapon did not retain owned image handle")
	}
	o.weapon.Colors12[1] = types.RGB{R: 255}
	clear(o.pix.Pix)
	dr.CallDraw(c.Viewport())
	if effectsPixelHash(o.pix) == first {
		t.Fatal("actual weapon material change did not affect pixels")
	}
	if got := c.Sub469920(image.Pt(48, 48)); len(got) != 3 || got[0] != 255 || got[1] != 255 || got[2] != 255 {
		t.Fatal("uniform world-light owner did not interpolate white")
	}
}
