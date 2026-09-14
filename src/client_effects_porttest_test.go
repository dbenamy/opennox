//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"os"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noximage"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

// Same raw-word oracle as noxrender's pixelHash16; dimensions participate, PNG
// encoding and RGB expansion do not. Keep this fixture out of production builds.
func effectsPixelHash(p *noximage.Image16) string {
	h := sha256.New()
	var size [8]byte
	binary.LittleEndian.PutUint32(size[:4], uint32(p.Rect.Dx()))
	binary.LittleEndian.PutUint32(size[4:], uint32(p.Rect.Dy()))
	h.Write(size[:])
	var word [2]byte
	for y := p.Rect.Min.Y; y < p.Rect.Max.Y; y++ {
		start := p.PixOffset(p.Rect.Min.X, y)
		for _, v := range p.Pix[start : start+p.Rect.Dx()] {
			binary.LittleEndian.PutUint16(word[:], v)
			h.Write(word[:])
		}
	}
	return hex.EncodeToString(h.Sum(nil))
}

func newEffectsTestOwner(t *testing.T) (*Client, *noximage.Image16) {
	t.Helper()
	core := new(server.Server)
	cli, free := client.PortTestEffectsClient(core, []string{"DrainManaOrb", "HealOrb", "CharmOrb", "WhiteOrb", "ManaBombOrb", "WhiteMoveOrb", "BlueMoveOrb", "CyanSpark", "DynamicChainLightning", "DynamicEnergyBolt", "DynamicLightning", "GreenZap", "OrbRay", "PlasmaRay", "RainOrbBlue", "RainOrbWhite", "WhiteSpark", "BlueSpark", "YellowSpark", "GreenSpark"})
	t.Cleanup(free)
	c := &Client{Client: cli, srv: &Server{Server: core}, r: NewNoxRender(cli.Render())}
	cli.ExtClient = unsafe.Pointer(c)
	old := legacy.GetClient
	legacy.GetClient = func() legacy.Client { return c }
	t.Cleanup(func() { legacy.GetClient = old })
	restore := legacy.PortTestEffectsOrbEnvironment()
	t.Cleanup(restore)
	d, freeData := noxrender.NewRenderData()
	t.Cleanup(freeData)
	rect := image.Rect(0, 0, 96, 96)
	d.SetClip(true)
	d.SetClipRect(rect)
	d.SetClipRect2(image.Rect(0, 0, 95, 95))
	d.SetRect3(rect)
	pix := noximage.NewImage16(rect)
	c.r.SetData(d)
	c.r.SetPixBuffer(pix)
	*cli.Viewport() = noxrender.Viewport{Screen: rect, World: rect, Size: rect.Size()}
	return c, pix
}

func TestClientEffectsOrbProbe(t *testing.T) {
	c, pix := newEffectsTestOwner(t)
	cli := c.Client
	dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 70))
	if dr == nil || cli.Objs.List1 != dr || cli.Objs.Count != 1 {
		t.Fatal("production drawable factory did not populate its owner")
	}
	state := unsafe.Slice((*byte)(unsafe.Pointer(dr)), int(unsafe.Sizeof(*dr)))
	state[444] = 12
	state[445] = 3
	state[446] = 2
	blank := effectsPixelHash(pix)
	if got := legacy.PortTestEffectsOrb(cli.Viewport(), dr, false); got != 1 {
		t.Fatalf("stationary orb return = %d", got)
	}
	if state[444] != 12 || state[446] != 1 || dr.PosVec != image.Pt(48, 70) {
		t.Fatal("orb changed unexpected drawable fields")
	}
	hash := effectsPixelHash(pix)
	if hash == blank {
		t.Fatal("original C orb did not render any pixels")
	}
	effectsDiagnostic(t, "orb", pix)
	const want = "2bc7703f6e4a987ab887c7a2c3cc13bc33a59715eebb9c5a5566c4daabff1b6f"
	if hash != want {
		t.Fatalf("orb framebuffer = %s, want repeated original-C capture %s", hash, want)
	}
}

type effectsOrbResult struct {
	Radius, Period, Counter byte
	Point                   image.Point
	Step, Return, Count     int
	Pixels                  string
	State                   []uint32
	Render                  []uint32
}

func TestClientEffectsOrbLifetime(t *testing.T) {
	c, pix := newEffectsTestOwner(t)
	var out []effectsOrbResult
	for _, radius := range []byte{1, 2, 5, 12, 31, 47} {
		for _, period := range []byte{0, 1, 3, 255} {
			for _, counter := range []byte{0, 1, 2, 255} {
				r := int(radius)
				for _, pos := range []image.Point{
					{48, 48}, {r - 1, 48}, {r, 48}, {r + 1, 48}, {95 - r, 48}, {96 - r, 48},
					{48, r - 1}, {48, r}, {48, r + 1}, {48, 95 - r}, {48, 96 - r},
				} {
					clear(pix.Pix)
					d := c.r.Data()
					d.Reset()
					d.SetClip(true)
					d.SetClipRect(pix.Rect)
					d.SetClipRect2(image.Rect(0, 0, 95, 95))
					d.SetRect3(pix.Rect)
					c.r.ClearPoints()
					dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, pos.Add(image.Pt(0, 22)))
					if dr == nil {
						t.Fatal("drawable allocation")
					}
					state := unsafe.Slice((*byte)(unsafe.Pointer(dr)), int(unsafe.Sizeof(*dr)))
					state[444], state[445], state[446] = radius, period, counter
					wantR, wantC := radius, counter
					alive := true
					for step := 0; step < 6 && alive; step++ {
						before := append([]byte(nil), state...)
						oldPixels := effectsPixelHash(pix)
						clipped := pos.X-int(wantR) < 0 || pos.Y-int(wantR) < 0 || pos.X+int(wantR) >= 96 || pos.Y+int(wantR) >= 96
						if !clipped && period != 0 {
							wantC--
							if wantC == 0 {
								wantC = period
								wantR--
								alive = wantR != 0
							}
						}
						wantRet := 1
						if !alive {
							wantRet = 0
						}
						got := legacy.PortTestEffectsOrb(c.Viewport(), dr, false)
						if got != wantRet || c.Objs.Count != wantRet {
							t.Fatalf("orb %d/%d/%d at%v step%d return/count %d/%d want%d", radius, period, counter, pos, step, got, c.Objs.Count, wantRet)
						}
						row := effectsOrbResult{Radius: radius, Period: period, Counter: counter, Point: pos, Step: step, Return: got, Count: c.Objs.Count, Pixels: effectsPixelHash(pix)}
						if clipped && row.Pixels != oldPixels {
							t.Fatal("clipped orb changed pixels")
						}
						if alive {
							before[444], before[446] = wantR, wantC
							if !bytes.Equal(before, state) {
								t.Fatal("orb mutated unexpected drawable bytes")
							}
							if c.Objs.List1 != dr {
								t.Fatal("live orb lost list membership")
							}
							// This fixture has one drawable with no pointer-valued fields in
							// its C prefix. The final Go client handle is intentionally omitted.
							row.State = append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(dr)), 128)...)
						} else if c.Objs.List1 != nil {
							t.Fatal("deleted orb retained list membership")
						}
						row.Render = append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(d)), int(unsafe.Sizeof(*d))/4)...)
						out = append(out, row)
					}
					if alive {
						c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
					}
				}
			}
		}
	}
	data, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	if prefix := os.Getenv("OPENNOX_CLIENT_EFFECTS_CAPTURE"); prefix != "" {
		if err := os.WriteFile(prefix+"-orb-lifetime.json", data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	t.Logf("orb-lifetime: %d snapshots %s", len(out), hash)
	const want = "02d3f918f3890220a797685bff2d41d0b9edf94e9573ee1c7e1c0e120cbf7edb"
	if hash != want {
		t.Fatalf("orb lifetime hash %s want repeated original-C capture %s", hash, want)
	}
}
