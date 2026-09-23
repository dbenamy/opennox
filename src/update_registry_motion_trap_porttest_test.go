//go:build porttest

package opennox

import (
	"bytes"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestUpdateRegistryMotionTrap(t *testing.T) {
	o := newWorldMotionTrapOwner(t)
	target := &o.units[0]
	saved := *target
	type initial struct {
		Name       string
		Scan, Shot uint32
		State, On  byte
	}
	cases := []initial{{"first-power", 99, 99, 0, 0}, {"scan-now", 0, 0, 0, 1}, {"fire-ready", 5, 0, 1, 1}, {"cooldown", 5, 2, 1, 1}, {"rescan", 1, 30, 1, 1}, {"other-state", 5, 0, 2, 1}, {"wrapped", 0xffffffff, 0xffffffff, 1, 1}}
	type row struct {
		Name, Trap      string
		Target, Powered bool
		Step            int
		Return          int32
		Data            [16]uint32
		Arrows          []worldMotionArrow
		Packets         [][]byte
	}
	var rows []row
	for _, trap := range []string{"ArrowTrap1", "ArrowTrap2", "Trigger"} {
		for _, hasTarget := range []bool{false, true} {
			for _, powered := range []bool{false, true} {
				for _, init := range cases {
					o.s.PortTestAIEmptyMap()
					o.configure(0)
					*target = saved
					worldGeometryResetObject(o.trap, 1001, 100, 100, false)
					o.trap.TypeInd = uint16(o.s.Types.IndByID(trap))
					o.trap.Direction1 = 0
					o.trap.ObjOwner = nil
					o.trap.ObjFlags = 4
					if powered {
						o.trap.ObjFlags |= 0x1000000
					}
					worldGeometryResetObject(target, 1002, 200, 100, false)
					target.ObjClass = object.ClassPlayer
					target.ObjFlags = 4
					if hasTarget {
						o.s.Map.AddObjectToIndex(target)
					}
					*o.data = [16]uint32{}
					o.data[0] = init.Scan
					o.data[1] = init.Shot
					o.data[2] = uint32(init.State)
					o.data[3] = uint32(o.s.Types.IndByID("MotionTestArrow"))
					o.data[12] = uint32(init.On)
					o.s.SetTickRate(30)
					for step := 0; step < 3; step++ {
						before := *o.data
						o.s.NetList.ResetAll()
						o.s.Objs.Pending = nil
						o.s.PortTestCombatAudioReset()
						legacy.PortTestRegisteredUpdate(o.trap, "SkullUpdate")
						rv := int32(0)
						arrows := o.created(t)
						packets := visibilityEffectsPackets(o.s)
						if !powered {
							want := before
							want[12] &^= 0xff
							if *o.data != want || len(arrows) != 0 {
								t.Fatal("unpowered trap contract")
							}
						}
						if powered && step == 0 && init.Name == "first-power" {
							if o.data[0] != 29 || (o.data[2]&0xff == 1) != hasTarget || (len(arrows) == 1) != hasTarget {
								t.Fatal("trap power-on reset/acquisition")
							}
						}
						if powered && step == 0 && init.Name == "fire-ready" && (len(arrows) != 1 || o.data[0] != 4 || o.data[1] != 29) {
							t.Fatal("trap firing cooldown contract")
						}
						wantFX := byte(0)
						if len(arrows) > 0 && trap == "ArrowTrap1" {
							wantFX = 1
						}
						if len(arrows) > 0 && trap == "ArrowTrap2" {
							wantFX = 2
						}
						packetCount := 0
						for _, p := range packets {
							if len(p) > 0 {
								packetCount++
								if wantFX == 0 || !bytes.Equal(p, []byte{161, 100, 0, 100, 0, wantFX}) {
									t.Fatal("trap effect packet", p, wantFX)
								}
							}
						}
						if wantFX != 0 && packetCount == 0 {
							t.Fatal("trap effect not delivered to actual players")
						}
						rows = append(rows, row{init.Name, trap, hasTarget, powered, step, rv, *o.data, arrows, packets})
					}
				}
			}
		}
	}
	spellbookCapture(t, "update-registry-motion-trap", rows, updateRegistryHashes["update-registry-motion-trap"])
	updateRegistryReportCounts(t)
}
