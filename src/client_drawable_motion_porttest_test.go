//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"math"
	"testing"
	"unsafe"
)

func TestClientDrawableMotion(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	type result struct {
		Op, Age, Variant, Visible int
		Frame                     uint32
		Return, Count             int
		State                     [][]uint32
		Deleted                   []uint32
	}
	var out []result
	for op := 0; op < 2; op++ {
		for _, frame := range []uint32{0, 100, 0xffffffff} {
			for _, age := range []int{0, 1, 3, 10, 30} {
				for variant := 0; variant < 36; variant++ {
					for visible := 0; visible < 2; visible++ {
						c.resetCase(env, pix, 1, frame)
						points := []image.Point{{-10000, -10000}, {10000, -10000}, {10000, 10000}, {-10000, 10000}}
						if visible == 0 {
							points = nil
						}
						_, restore := c.Cli().PortTestObjectRenderSight(points)
						origin := []image.Point{{300, 400}, {1, 1}, {5887, 5887}, {0, 0}, {5888, 5888}, {30, 50}}[variant%6]
						dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, origin)
						if dr == nil {
							t.Fatal("factory")
						}
						raw := unsafe.Slice((*byte)(dr.C()), 512)
						dr.AnimStart = frame - uint32(age)
						dr.Field_81 = uint32(origin.X)
						dr.Field_82 = uint32(origin.Y)
						vx := []float32{-9.75, -1, 0, 1, 9.75, 35.125}[variant/6]
						vy := -vx / 2
						drag := []float32{0, 0.1, 0.25, 1}[variant%4]
						binary.LittleEndian.PutUint32(raw[468:], math.Float32bits(vx))
						binary.LittleEndian.PutUint32(raw[472:], math.Float32bits(vy))
						binary.LittleEndian.PutUint32(raw[476:], math.Float32bits(drag))
						target := []image.Point{{300, 400}, {305, 405}, {290, 400}, {0, 0}, {5887, 5887}, {330, 450}}[variant%6]
						if op == 1 {
							binary.LittleEndian.PutUint16(raw[432:], uint16(target.X))
							binary.LittleEndian.PutUint16(raw[434:], uint16(target.Y))
							raw[443] = []byte{0, 1, 9, 10, 50, 255}[variant/6]
						}
						ret := legacy.PortTestDrawableMotion(op, c.Viewport(), dr)
						if ret < 0 || ret > 1 || c.Objs.Count != ret || len(c.Deleted) != 1-ret {
							t.Fatalf("motion op%d return%d count%d deleted%v", op, ret, c.Objs.Count, c.Deleted)
						}
						if op == 0 && (drag == 0 || drag == 1) {
							want := origin
							if drag == 0 {
								want.X = int(float32(float64(origin.X) + float64(vx)*float64(age+1)))
								want.Y = int(float32(float64(origin.Y) + float64(float32(float64(vy)*float64(age+1)))))
							}
							alive := visible == 1 && want.X > 0 && want.Y > 0 && want.X < 5888 && want.Y < 5888
							if (ret == 1) != alive {
								t.Fatalf("exact motion age%d v%v drag%v pos%v ret%d", age, vx, drag, want, ret)
							}
							if alive && dr.PosVec != want {
								t.Fatalf("exact position %v want%v", dr.PosVec, want)
							}
						}
						if op == 1 && origin == target && ret != 0 {
							t.Fatal("arrived target retained drawable")
						}
						if op == 0 && visible == 0 && ret != 0 {
							t.Fatal("invisible moving sprite retained")
						}
						out = append(out, result{op, age, variant, visible, frame, ret, c.Objs.Count, c.snapshotDrawables(t), append([]uint32(nil), c.Deleted...)})
						restore()
					}
				}
			}
		}
	}
	drawableStateCapture(t, "motion", out, "58190061966a66b15cd8fc4eabd62a05885f4c199f89119233c6604b21fb910f")
}
