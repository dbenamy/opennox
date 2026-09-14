//go:build porttest

package opennox

import (
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
)

func TestClientParticleBubbleLifecycle(t *testing.T) {
	c, pix, effects := newEffectsFullOwner(t)
	env := legacy.PortTestNewClientParticleEnvironment()
	t.Cleanup(env.Restore)
	c.resetCase(effects, pix, 17, 120)
	env.Reset()
	c.srv.SetTickRate(60)
	dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48))
	bytes := unsafe.Slice((*byte)(dr.C()), 512)
	words := unsafe.Slice((*uint32)(dr.C()), 128)
	words[108], words[109] = 0xf81f, 0xffff
	dr.ZVal = 10
	bytes[440], bytes[441], bytes[442], bytes[443] = 4, 1, 1, 2
	bytes[444], bytes[445], bytes[446] = 3, 1, 1
	c.Objs.TransparentDecay(dr, 4)
	// Explicit early lifecycle: grow to five, reverse, and shrink to four;
	// vertical direction reverses on its independent timer.
	early := [][7]int{{5, 2, 2, 3, -1, 10, 124}, {5, 2, 1, 2, -1, 9, 124}, {4, 2, 2, 1, -1, 8, 124}, {4, 2, 1, 3, 1, 7, 124}, {4, 3, 3, 2, 1, 7, 184}}
	for frame := uint32(120); frame <= 139; frame++ {
		c.srv.SetFrame(frame)
		got := legacy.PortTestClientDrawParticle(5, c.Viewport(), dr, [4]int32{})
		if frame < 139 {
			if got != 1 || c.Objs.Count != 1 || c.Objs.DeadlineList != dr || len(c.Deleted) != 0 {
				t.Fatal("bubble deleted before fade completed")
			}
			if frame < 125 {
				want := early[frame-120]
				state := [7]int{int(bytes[440]), int(bytes[441]), int(bytes[442]), int(bytes[445]), int(int8(bytes[446])), int(dr.ZVal), int(dr.Deadline)}
				if state != want {
					t.Fatalf("bubble frame%d state%v want%v", frame, state, want)
				}
			}
		} else if got != 0 || c.Objs.Count != 0 || c.Objs.DeadlineList != nil || len(c.Deleted) != 1 || c.Deleted[0] != 1 {
			t.Fatal("bubble final fade ownership")
		}
	}
	if len(c.Calls) != 1 || c.srv.Rand.Logic.Index() != 17 || c.srv.Rand.Other.Index() != 18 {
		t.Fatal("bubble lifecycle changed RNG or allocated")
	}
}

func TestClientParticleVortexLifecycle(t *testing.T) {
	c, pix, effects := newEffectsFullOwner(t)
	env := legacy.PortTestNewClientParticleEnvironment()
	t.Cleanup(env.Restore)
	c.resetCase(effects, pix, 17, 120)
	env.Reset()
	dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48))
	words := unsafe.Slice((*uint32)(dr.C()), 128)
	bytes := unsafe.Slice((*byte)(dr.C()), 512)
	words[108], words[109] = 0xbdef, 0xffff
	words[110], words[111] = 48, 48
	bytes[448], bytes[449], bytes[450], bytes[451] = 0, 2, 20, 1
	if got := legacy.PortTestClientDrawParticle(11, c.Viewport(), dr, [4]int32{}); got != 1 || bytes[448] != 2 || bytes[450] != 50 || dr.ZVal != 1 || c.Objs.Count != 1 {
		t.Fatal("vortex angle/height/radius transition")
	}
	c.snapshotDrawables(t)
	if got := legacy.PortTestClientDrawParticle(11, c.Viewport(), dr, [4]int32{}); got != 0 || c.Objs.Count != 0 || len(c.Deleted) != 1 || c.Deleted[0] != 1 {
		t.Fatal("vortex viewport deletion")
	}
	if len(c.Calls) != 1 || c.srv.Rand.Logic.Index() != 17 || c.srv.Rand.Other.Index() != 18 {
		t.Fatal("vortex lifecycle changed RNG or allocated")
	}
}
