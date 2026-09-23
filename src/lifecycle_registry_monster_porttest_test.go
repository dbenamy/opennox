//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/unit/ai"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
	"unsafe"
)

func TestLifecycleRegistryMonsterInit(t *testing.T) {
	type row struct {
		Name   string
		NPC    bool
		Before uint16
		Health server.HealthData
		Speed  uint32
		Update []byte
	}
	var rows []row
	for _, name := range []string{"MonsterInit", "ShopkeeperInit"} {
		for _, npc := range []bool{false, true} {
			for _, hp := range []uint16{80, 100} {
				t.Run(fmt.Sprintf("%s/npc%v/hp%d", name, npc, hp), func(t *testing.T) {
					s := newObjectXferOwner(t)
					u := newObjectXferSimple(t, s)
					ud, freeUD := alloc.New(server.MonsterUpdateData{})
					defer freeUD()
					health, freeHealth := alloc.New(server.HealthData{})
					defer freeHealth()
					*health = server.HealthData{Cur: hp, Max: 100}
					ud.AIAction340 = uint32(ai.ACTION_INVALID)
					ud.AIStackInd = -1
					ud.Field338 = 0.5
					ud.Field332 = 0.5
					ud.StatusFlags = object.MonStatusHoldYourGround | object.MonStatusAlwaysRun
					u.ObjClass = object.ClassMonster
					if npc {
						u.ObjSubClass = object.SubClass(object.MonsterNPC)
					}
					u.HealthData = health
					u.UpdateData = unsafe.Pointer(ud)
					defer func() { u.HealthData = nil; u.UpdateData = nil }()
					u.PosVec = types.Pointf{X: 41, Y: 53}
					u.Direction1 = 17
					u.SpeedBase = 20
					u.Init = legacy.PortTestLifecycleRegistryPointer(name, false)
					lifecyclePending(t, u)
					want := hp
					if hp == 100 {
						want = 50
					}
					if health.Cur != want || health.Field2 != want || health.Max != 100 {
						t.Fatalf("monster health %+v", *health)
					}
					for _, v := range ud.HealthGraph103 {
						if v != want {
							t.Fatal("monster health graph")
						}
					}
					if ud.Direction94 != 17 || ud.Pos95 != u.PosVec || ud.FleeRange != 0 || !ud.StatusFlags.Has(object.MonStatusRunning) {
						t.Fatal("monster direction/status")
					}
					if npc {
						if u.SpeedBase != float32(1.95) {
							t.Fatal("NPC speed")
						}
					} else if u.SpeedBase < 18.99 || u.SpeedBase > 21.01 {
						t.Fatal("monster random speed")
					}
					rows = append(rows, row{name, npc, hp, *health, math.Float32bits(u.SpeedBase), bytes.Clone(unsafe.Slice((*byte)(unsafe.Pointer(ud)), int(unsafe.Sizeof(*ud))))})
				})
			}
		}
	}
	spellbookCapture(t, "lifecycle-registry-monster-init", rows, "9825ff7f8d6a9f55b0a1efef7e1a0746eff6ccdf88f13a6b4bc29ae98fd7a075")
}
