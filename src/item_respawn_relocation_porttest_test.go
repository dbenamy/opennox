//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
	"unsafe"
)

func TestItemRespawnRelocation(t *testing.T) {
	o := newMatchRosterOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	t.Cleanup(o.s.PortTestRewardTypes([]string{"Crown", "PortTestRewardWeapon", "RespawnFixture"}, nil, true, 0, 0x82))
	_, restore := legacy.PortTestItemRespawnGlobals()
	t.Cleanup(restore)
	t.Cleanup(o.s.PortTestCombatAudioReset)
	type row struct {
		Name     string
		Position types.Pointf
		Health   uint16
		Ammo     [2]byte
		Charge   uint32
		Audio    []types.Pointf
		Packets  [][]byte
	}
	var rows []row
	for _, class := range []uint32{0x1000, 0x1000000, 0x2000000} {
		for _, frame := range []uint32{149, 150, 151} {
			for _, distance := range []float32{math.Nextafter32(50, 0), 50, math.Nextafter32(50, 100), 60} {
				name := fmt.Sprintf("class%x/frame%d/distance%08x", class, frame, math.Float32bits(distance))
				t.Run(name, func(t *testing.T) {
					o.reset()
					o.s.SetFrame(frame)
					o.s.PortTestCombatAudioReset()
					u := o.s.NewObjectByTypeID("PortTestRewardWeapon")
					if u == nil {
						t.Fatal("item allocation")
					}
					oldUse, oldHealth := u.UseData, u.HealthData
					defer func() { noxServer.sub_517870(u); u.UseData, u.HealthData = oldUse, oldHealth; o.s.Objs.FreeObject(u) }()
					u.UseData.Ptr = o.record(t, 128)
					u.HealthData = (*server.HealthData)(o.record(t, int(unsafe.Sizeof(server.HealthData{}))))
					hp := (*[6]uint16)(unsafe.Pointer(u.HealthData))
					hp[0] = 7
					hp[2] = 88
					data := unsafe.Slice((*byte)(u.UseData.Ptr), 128)
					data[0] = 17
					data[1] = 29
					data[108] = 4
					data[109] = 40
					*(*uint32)(unsafe.Add(u.UseData.Ptr, 112)) = 10
					u.ObjClass = object.Class(class)
					u.ObjFlags = 0
					u.InvHolder = nil
					u.PosVec = types.Pointf{X: 100, Y: 100}
					u.Field32 = 0
					cleanup := legacy.PortTestTeamRuntimeRespawns([]*server.Object{u})
					defer cleanup()
					u.PosVec = types.Pointf{X: 100 + distance, Y: 100}
					data[0] = 1
					data[1] = 2
					before := u.PosVec
					legacy.PortTestItemRespawn("tick", nil)
					// Quantize the actual input coordinates before the independent distance check.
					dx := float64(before.X) - 100
					move := frame > 150 && dx*dx > 2500
					want := before
					if move {
						want = types.Pointf{X: 100, Y: 100}
					}
					if u.PosVec != want {
						t.Fatal("relocation threshold", u.PosVec, want)
					}
					wantHP := uint16(7)
					if move {
						wantHP = 88
					}
					if hp[0] != wantHP {
						t.Fatal("health restoration", hp[0], wantHP)
					}
					ammo := [2]byte{data[0], data[1]}
					wantAmmo := [2]byte{1, 2}
					if move && class&0x1000000 != 0 {
						wantAmmo = [2]byte{17, 29}
					}
					if ammo != wantAmmo {
						t.Fatal("ammo restoration", ammo, wantAmmo)
					}
					charge := *(*uint32)(unsafe.Add(u.UseData.Ptr, 112))
					wantCharge := uint32(10)
					if move && class&0x1000 != 0 {
						wantCharge = 100
					}
					if charge != wantCharge {
						t.Fatal("charge restoration", charge, wantCharge)
					}
					r := row{Name: name, Position: u.PosVec, Health: hp[0], Ammo: ammo, Charge: charge, Packets: visibilityEffectsPackets(o.s)}
					sounds := o.s.PortTestCombatAudioSnapshot()
					wantSounds := 0
					if move {
						wantSounds = 2
					}
					if len(sounds) != wantSounds {
						t.Fatal("sound count", len(sounds), wantSounds)
					}
					for _, a := range sounds {
						if a.ID != 283 || !a.ByPos || a.Kind != 0 || a.Code != 0 {
							t.Fatal("relocation audio", a)
						}
						r.Audio = append(r.Audio, a.Pos)
					}
					if move && (r.Audio[0] != before || r.Audio[1] != want) {
						t.Fatal("audio positions", r.Audio)
					}
					rows = append(rows, r)
				})
			}
		}
	}
	spellbookCapture(t, "item-respawn-relocation", rows, "cabc108e6d22035170b003a11f3d6b46a1e033932283dd450ef37127d3a094c4")
}
