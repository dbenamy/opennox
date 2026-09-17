//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
)

func TestSpatialTargetingPrediction(t *testing.T) {
	o := newCollisionCoreOwner(t)
	a, b := newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s)
	sa, sb := *a, *b
	t.Cleanup(func() { *a = sa; *b = sb })
	type row struct {
		Source, Target, Velocity [2]uint32
		Speed                    uint32
		Output                   [2]uint32
		Return                   uint32
	}
	var rows []row
	positions := []types.Pointf{{0, 0}, {3, 4}, {-3, -4}, {.001, 100000}, {16777216, -16777216}, {100.00001, 99.99999}}
	for _, ap := range positions {
		for _, bp := range positions {
			for _, vel := range []types.Pointf{{0, 0}, {1, -1}, {.00001, 12345}, {-3, 4}} {
				for _, speed := range []float32{0, -1, .001, 1, 5, 12345} {
					*a = sa
					*b = sb
					a.PosVec = ap
					b.PosVec = bp
					b.VelVec = vel
					beforeA, beforeB := *a, *b
					out := types.Pointf{17, -19}
					rv := legacy.PortTestSpatialPredict(a, b, speed, &out)
					if rv != uint32(uintptr(b.CObj())) {
						t.Fatal("prediction must return target")
					}
					if *a != beforeA || *b != beforeB {
						t.Fatal("prediction changed objects")
					}
					if speed != 0 && vel == (types.Pointf{}) && out != bp {
						t.Fatal("stationary target prediction moved")
					}
					if ap == (types.Pointf{}) && bp == (types.Pointf{3, 4}) && vel == (types.Pointf{1, -1}) && speed == 5 && out != (types.Pointf{4, 3}) {
						t.Fatal("3-4-5 prediction", out)
					}
					rows = append(rows, row{motionBits(ap), motionBits(bp), motionBits(vel), math.Float32bits(speed), motionBits(out), 1002})
				}
			}
		}
	}
	spellbookCapture(t, "spatial-targeting-prediction", rows, "f1238eaccfcf8e62c00674314e46fafb2c5e6b038208799af591f312cabe8c8f")
}

func TestSpatialTargetingEligibility(t *testing.T) {
	o := newCollisionCoreOwner(t)
	a, b := newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s)
	sa, sb := *a, *b
	t.Cleanup(func() { *a = sa; *b = sb })
	old := noxflags.GetGamePlay()
	noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0))
	t.Cleanup(func() { noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0)); noxflags.SetGamePlay(old) })
	type row struct {
		Class            [2]uint32
		Flags            [2]uint32
		Callbacks, Teams int
		Gameplay         uint32
		Direct, Team     int32
	}
	var rows []row
	flags := []object.Flags{0, 1, 4, 0x10, 0x11, 0x20, 0x40, 0x80, 0xa0, 0xc0, 0x400, 0x4000, 0x4080}
	for _, ac := range []object.Class{0, 1, 2, 4, 0x80} {
		for _, bc := range []object.Class{0, 1, 2, 4, 0x80} {
			for _, af := range flags {
				for _, bf := range flags {
					for teams := 0; teams < 3; teams++ {
						*a = sa
						*b = sb
						a.ObjClass = ac
						b.ObjClass = bc
						a.ObjFlags = af
						b.ObjFlags = bf
						a.Collide = o.callback
						b.Collide = o.callback
						a.ObjOwner = nil
						b.ObjOwner = nil
						a.TeamVal.ID = 0
						b.TeamVal.ID = 0
						if teams != 0 {
							a.TeamVal.ID = 1
							b.TeamVal.ID = server.TeamID(teams)
						}
						direct := legacy.PortTestSpatialEligible(a, b, false)
						team := legacy.PortTestSpatialEligible(a, b, true)
						if bc&1 != 0 && direct != 0 {
							t.Fatal("projectile target admitted")
						}
						if (af&0x20 != 0 || bf&0x60 != 0) && direct != 0 {
							t.Fatal("disabled collision admitted")
						}
						if teams == 1 && team != 0 {
							t.Fatal("team filter admitted same team")
						}
						rows = append(rows, row{[2]uint32{uint32(ac), uint32(bc)}, [2]uint32{uint32(af), uint32(bf)}, 3, teams, 0, direct, team})
					}
				}
			}
		}
	}
	// Explicit callback absence and gameplay override, including reverse argument order.
	for callbacks := 0; callbacks < 4; callbacks++ {
		for gameplay := uint32(0); gameplay < 2; gameplay++ {
			*a = sa
			*b = sb
			a.ObjClass = 0
			b.ObjClass = 0
			a.ObjFlags = 0
			b.ObjFlags = 0
			a.ObjOwner = nil
			b.ObjOwner = nil
			a.Collide = nil
			b.Collide = nil
			if callbacks&1 != 0 {
				a.Collide = o.callback
			}
			if callbacks&2 != 0 {
				b.Collide = o.callback
			}
			a.TeamVal.ID = 1
			b.TeamVal.ID = 1
			noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0))
			noxflags.SetGamePlay(noxflags.GameplayFlag(gameplay))
			direct := legacy.PortTestSpatialEligible(a, b, false)
			team := legacy.PortTestSpatialEligible(a, b, true)
			want := int32(0)
			if callbacks == 3 {
				want = 1
			}
			if direct != want || team != want*int32(gameplay) {
				t.Fatal("callback/gameplay contract", callbacks, gameplay, direct, team)
			}
			rows = append(rows, row{[2]uint32{}, [2]uint32{}, callbacks, 1, gameplay, direct, team})
		}
	}
	spellbookCapture(t, "spatial-targeting-eligibility", rows, "f45ca90d37a9b929841b61082f74ac98863a54fc3d74ce70ee0cf84764067cc2")
}
