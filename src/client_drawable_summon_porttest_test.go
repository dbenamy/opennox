//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

type drawableSummonOwner struct {
	*effectsTestClient
	FailChild bool
}

func (c *drawableSummonOwner) Nox_new_drawable_for_thing(typ int) *client.Drawable {
	if c.FailChild {
		return nil
	}
	dr := c.Client.Nox_new_drawable_for_thing(typ)
	if dr != nil {
		c.next++
		c.refs[dr] = c.next
	}
	return dr
}
func TestClientDrawableSummons(t *testing.T) {
	base, pix, env := newEffectsFullOwner(t, "SummonEffect")
	c := &drawableSummonOwner{effectsTestClient: base}
	old := legacy.GetClient
	legacy.GetClient = func() legacy.Client { return c }
	t.Cleanup(func() { legacy.GetClient = old })
	words, restore := legacy.PortTestDrawableEffectGlobals()
	t.Cleanup(restore)
	type row struct {
		Dir, Failure       int
		Frame              uint32
		Created, After     int
		Before, AfterState [][]uint32
		Calls              []effectsSpawnCall
		Logic, Other       int
	}
	var rows []row
	for _, frame := range []uint32{0, 100, 0xffffffff} {
		for dir := 0; dir < 256; dir++ {
			for failure := 0; failure < 3; failure++ {
				c.resetCase(env, pix, 17, frame)
				for _, p := range words {
					*p = 0
				}
				c.FailChild = failure == 2
				if failure == 1 {
					c.FailEvery = 1
				}
				owner := uint32(0x89ab1234)
				legacy.PortTestDrawableEffect(0, nil, [5]uint32{300, 400, owner, 4, uint32(dir)})
				parent := c.Objs.List1
				var child *client.Drawable
				want := 2
				if failure == 1 {
					want = 0
				} else if failure == 2 {
					want = 1
				}
				if c.Objs.Count != want {
					t.Fatalf("summon failure%d count%d want%d", failure, c.Objs.Count, want)
				}
				var before [][]uint32
				if parent != nil {
					child = *(**client.Drawable)(unsafe.Add(parent.C(), 432))
					if failure == 0 {
						if child == nil || c.refs[child] == 0 || child.NextPtr != nil || child.PosVec != image.Pt(300, 400) || child.AnimInd != 8 || parent.AnimStart != frame || *(*uint32)(unsafe.Add(parent.C(), 436)) != (owner<<16|owner>>16) {
							t.Fatal("summon ownership/metadata")
						}
						// Independent octant ranges from the audited shipped direction table.
						wantDir := byte(5)
						for _, r := range [][3]int{{18, 46, 8}, {47, 83, 7}, {84, 110, 6}, {111, 147, 3}, {148, 172, 0}, {173, 209, 1}, {210, 236, 2}} {
							if dir >= r[0] && dir <= r[1] {
								wantDir = byte(r[2])
							}
						}
						if child.AnimDir != wantDir {
							t.Fatalf("summon direction %d got%d want%d", dir, child.AnimDir, wantDir)
						}
						before = c.snapshotDrawables(t, child)
						before[0][109] = c.refs[child]
					} else {
						if child != nil {
							t.Fatal("failed child linked")
						}
						before = c.snapshotDrawables(t)
					}
				}
				created := c.Objs.Count
				if failure == 0 {
					legacy.PortTestDrawableEffect(1, nil, [5]uint32{0x4321})
					if c.Objs.Count != created || c.Objs.List1 != parent {
						t.Fatal("unmatched summon removed")
					}
					legacy.PortTestDrawableEffect(1, nil, [5]uint32{owner & 0xffff})
					if c.Objs.Count != 50 {
						t.Fatalf("summon finish count%d want50 point effects", c.Objs.Count)
					}
					for dr := c.Objs.List1; dr != nil; dr = dr.NextPtr {
						if dr == parent || dr == child {
							t.Fatal("summon allocation retained")
						}
					}
				}
				rows = append(rows, row{dir, failure, frame, created, c.Objs.Count, before, c.snapshotDrawables(t), append([]effectsSpawnCall(nil), c.Calls...), c.srv.Rand.Logic.Index(), c.srv.Rand.Other.Index()})
			}
		}
	}
	drawableStateCapture(t, "summons", rows, "3f62476ae211c3b83079c10265b0785555b336312529fbeb4944207c2e0bfa28")
}
