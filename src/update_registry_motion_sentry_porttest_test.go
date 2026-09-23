//go:build porttest

package opennox

import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func TestUpdateRegistryMotionSentry(t *testing.T) {
	o := newCollisionCoreOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	configure, unchanged, free := o.s.PortTestPathWalls()
	t.Cleanup(free)
	sentry := newObjectXferSimple(t, o.s)
	saved := *sentry
	data := collisionCoreGuarded(t, o, 12)
	words, restore := legacy.PortTestWorldMotionListGlobals()
	t.Cleanup(restore)
	t.Cleanup(func() { *sentry = saved })
	type row struct {
		Wall                           int
		Powered, Destroyed             bool
		Angle, Speed                   uint32
		Return, Flags, Head, NextAngle uint32
		End                            [2]uint32
	}
	var rows []row
	for wall := 0; wall < 3; wall++ {
		for _, powered := range []bool{false, true} {
			for _, destroyed := range []bool{false, true} {
				for _, angle := range []float32{0, 0.125, math.Pi / 4, math.Pi / 2, math.Pi, -math.Pi / 4, 20 * math.Pi} {
					for _, speed := range []float32{-0.1, 0, 0.1} {
						o.s.PortTestAIEmptyMap()
						configure(wall)
						*words["sentry"] = 0
						*sentry = saved
						worldGeometryResetObject(sentry, 1001, 100, 100, false)
						sentry.ObjFlags = 4
						sentry.UpdateData = data
						if powered {
							sentry.ObjFlags |= 0x1000000
						}
						if destroyed {
							sentry.ObjFlags |= 0x20
						}
						*(*[3]float32)(data) = [3]float32{angle, 0.375, speed}
						sentry.Pos39 = types.Pointf{1, 2}
						legacy.PortTestRegisteredUpdate(sentry, "SentryGlobeUpdate")
						rv := uint32(0)
						head := uint32(0)
						if *words["sentry"] != 0 {
							if *words["sentry"] != uint32(uintptr(sentry.CObj())) {
								t.Fatal("sentry list head")
							}
							head = 1001
						}
						if (head != 0) == destroyed {
							t.Fatal("sentry membership after update")
						}
						next := (*[3]float32)(data)[0]
						if powered && next != angle+speed || !powered && next != 0.375 {
							t.Fatal("sentry angular advance/reset")
						}
						if powered && angle == 0 && wall == 0 && sentry.Pos39 != (types.Pointf{700, 100}) {
							t.Fatal("sentry beam range", sentry.Pos39)
						}
						if powered && angle == 0 && wall == 1 && (sentry.Pos39.X <= 100 || sentry.Pos39.X >= 200) {
							t.Fatal("sentry beam must stop at real wall", sentry.Pos39)
						}
						if !unchanged() {
							t.Fatal("sentry changed wall")
						}
						rows = append(rows, row{wall, powered, destroyed, math.Float32bits(angle), math.Float32bits(speed), rv, uint32(sentry.ObjFlags), head, math.Float32bits(next), motionBits(sentry.Pos39)})
					}
				}
			}
		}
	}
	spellbookCapture(t, "update-registry-motion-sentry", rows, updateRegistryHashes["update-registry-motion-sentry"])
	updateRegistryReportCounts(t)
}
