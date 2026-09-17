//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
)

func TestWorldMotionSentryBeam(t *testing.T) {
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
						rv := legacy.PortTestWorldMotionList("sentry-update", sentry, 0)
						if destroyed {
							if rv != uint32(uintptr(sentry.CObj())) {
								t.Fatal("destroyed sentry return")
							}
							rv = 1001
						}
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
	spellbookCapture(t, "world-motion-sentry-beam", rows, "0e94eebfaf4411f215eb3c61a0028277163eb7bd29cc62964cc58b3dd8fe1d59")
}

func TestWorldMotionSentryContacts(t *testing.T) {
	o := newCollisionCoreOwner(t)
	sentry, target := newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s)
	ss, st := *sentry, *target
	t.Cleanup(func() { *sentry = ss; *target = st })
	callback, observed, restore := legacy.PortTestWorldDamageObserver()
	t.Cleanup(restore)
	hp := collisionCoreGuarded(t, o, 12)
	ids := collisionCoreIDs(sentry, target, &o.units[0])
	type row struct {
		Quest, Parent, Health bool
		Class, Flags          uint32
		Position              [2]uint32
		Damage                [6]uint32
		Sounds                []int
	}
	var rows []row
	for _, quest := range []bool{false, true} {
		restoreFlags := noxflags.PortTestGameFlags(0)
		t.Cleanup(restoreFlags)
		if quest {
			noxflags.SetGame(noxflags.GameModeQuest)
		}
		for _, parent := range []bool{false, true} {
			for _, health := range []bool{false, true} {
				for _, cls := range []object.Class{object.ClassSimple, object.ClassMonster} {
					for _, flags := range []object.Flags{0, 1, 0x40, 0x10, 0x8010} {
						for _, pos := range []types.Pointf{{100, 100}, {150, 100}, {150, 109.999}, {150, 110}, {150, 110.001}, {200, 100}, {201, 100}, {99, 100}} {
							*observed = [6]uint32{}
							o.s.PortTestCombatAudioReset()
							*sentry = ss
							*target = st
							worldGeometryResetObject(sentry, 1001, 100, 100, false)
							worldGeometryResetObject(target, 1002, pos.X, pos.Y, false)
							target.ObjClass = cls
							target.ObjFlags = flags
							target.Damage = callback
							target.HealthData = nil
							if health {
								target.HealthData = (*server.HealthData)(hp)
							}
							sentry.ObjOwner = nil
							if parent {
								sentry.ObjOwner = &o.units[0]
							}
							ray := [4]float32{100, 100, 200, 100}
							legacy.PortTestWorldMotionSentryCandidate(target, sentry, &ray)
							got := *observed
							for i := 1; i <= 3; i++ {
								got[i] = collisionCoreID(t, ids, got[i])
							}
							hit := health && flags&0x41 == 0 && (flags&0x10 == 0 || quest && cls&2 != 0 && flags&0x8000 == 0) && pos.X >= 100 && pos.X <= 200 && pos.Y < 110
							if (got[0] == 1) != hit || got[0] > 1 {
								t.Fatal("sentry strict-radius/segment/eligibility contract", quest, parent, health, cls, flags, pos, got)
							}
							if hit {
								owner := uint32(1001)
								if parent {
									owner = 1003
								}
								if got != ([6]uint32{1, 1002, owner, 1001, 500, 16}) {
									t.Fatal("sentry damage attribution", got)
								}
							}
							var sounds []int
							for _, e := range o.s.PortTestCombatAudioSnapshot() {
								if e.Obj != target {
									t.Fatal("sentry audio owner")
								}
								sounds = append(sounds, int(e.ID))
							}
							if len(sounds) != int(got[0]) || len(sounds) == 1 && sounds[0] != 298 {
								t.Fatal("sentry contact audio", sounds)
							}
							rows = append(rows, row{quest, parent, health, uint32(cls), uint32(flags), motionBits(pos), got, sounds})
						}
					}
				}
			}
		}
		restoreFlags()
	}
	spellbookCapture(t, "world-motion-sentry-contacts", rows, "2db2ee309bdcf5c60275dd723a6d32104d0c86fa68b7387c4d3dc32da82894e0")
}

func TestWorldMotionSentryReporting(t *testing.T) {
	o := newCollisionCoreOwner(t)
	_, restoreScratch := legacy.PortTestWorldMotionRoundScratch()
	t.Cleanup(restoreScratch)
	u := newObjectXferSimple(t, o.s)
	saved := *u
	data := collisionCoreGuarded(t, o, 12)
	words, restore := legacy.PortTestWorldMotionListGlobals()
	t.Cleanup(restore)
	pl := o.s.Players.ByInd(1)
	oldPlayer := *pl
	t.Cleanup(func() { *u = saved; *pl = oldPlayer })
	type row struct {
		Start, End [2]uint32
		Powered    bool
		Angle      uint32
		Packet     []byte
	}
	var rows []row
	for _, powered := range []bool{false, true} {
		for _, ray := range [][4]float32{{50, 50, 150, 150}, {50, 50, 50, 150}, {-50, -50, 50, 50}, {-50, 100, 50, 100}, {50, 100, 150, 100}, {150, 150, 50, 50}, {100, 100, 100, 100}, {0.49, 0.5, 199.49, 199.5}} {
			*u = saved
			u.UpdateData = data
			u.ObjFlags = 0
			u.PosVec = types.Pointf{ray[0], ray[1]}
			u.Pos39 = types.Pointf{ray[2], ray[3]}
			u.InvNextItem = nil
			if powered {
				u.ObjFlags |= 0x1000000
			}
			*(*[3]float32)(data) = [3]float32{1.25, 0.375, 0.1}
			*words["sentry"] = uint32(uintptr(u.CObj()))
			pl.Field10 = 50
			pl.Field12 = 50
			pl.Pos3632Vec = types.Pointf{100, 100}
			o.s.NetList.ResetAll()
			legacy.PortTestWorldMotionSentryReport(1)
			packet := o.s.NetList.CopyPacketsA(1, netlist.Kind1)
			visible := powered && min(ray[0], ray[2]) < 150 && max(ray[0], ray[2]) > 50 && min(ray[1], ray[3]) < 150 && max(ray[1], ray[3]) > 50
			if (len(packet) == 9) != visible || len(packet) != 0 && (len(packet) != 9 || packet[0] != 0x95) {
				t.Fatal("sentry viewport strict overlap", ray, packet)
			}
			if visible && ray == ([4]float32{50, 50, 150, 150}) {
				for i, w := range []uint16{50, 50, 150, 150} {
					if binary.LittleEndian.Uint16(packet[1+2*i:]) != w {
						t.Fatal("sentry packet coordinates")
					}
				}
			}
			angle := (*[3]float32)(data)[0]
			if !powered && angle != 0.375 {
				t.Fatal("reporting resets inactive sentry angle")
			}
			rows = append(rows, row{motionBits(u.PosVec), motionBits(u.Pos39), powered, math.Float32bits(angle), append([]byte(nil), packet...)})
		}
	}
	spellbookCapture(t, "world-motion-sentry-reporting", rows, "f30dba050bd6a3a64b64f211530a782b14175ba4f41cf59644011b67cf6639f4")
}
