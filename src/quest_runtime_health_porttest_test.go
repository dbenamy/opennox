//go:build porttest

package opennox

import (
	"fmt"
	"math"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/server"
)

func TestQuestRuntimeHealthScaling(t *testing.T) {
	o := newQuestRuntimeOwner(t)
	u := newCreatureXferObject(t, o.s, "Monster")
	t.Cleanup(o.s.PortTestSpellEffectTypes(nil, []string{"PortCreatureMonster"}, u))
	definition := o.s.Types.ByInd(int(u.TypeInd)).Health()
	*definition = server.HealthData{Cur: 80, Max: 100}
	oldList := o.s.Objs.First()
	o.s.Objs.SetObjects(u)
	t.Cleanup(func() { o.s.Objs.SetObjects(oldList) })
	override := o.record(t, 80)
	objectXferSetWord(override, 72, 60)
	originalData := u.UpdateData
	t.Cleanup(func() { u.ObjClass = object.ClassMonster })
	type row struct {
		Name           string
		HP             [3]uint16
		Damage, Health uint32
		History        [32]uint16
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "quest-runtime-health-scaling", rows, "0f79db5ba9238e69309e4aaeb235a6cb29016ded196b49f1c963ef557e346834")
	}()
	for _, class := range []object.Class{object.ClassSimple, object.ClassMonster, object.Class(0x20000)} {
		for _, disabled := range []bool{false, true} {
			for _, hp := range [][2]uint16{{0, 0}, {0, 100}, {50, 100}, {100, 100}, {32767, 32767}, {32768, 32768}, {65535, 65535}} {
				for _, difficulty := range []float32{0, 1, 2, 10} {
					for _, special := range []int{0, 1, 2} {
						name := fmt.Sprintf("class%x/disabled%t/hp%v/difficulty%g/special%d", uint32(class), disabled, hp, difficulty, special)
						t.Run(name, func(t *testing.T) {
							u.ObjClass = class
							u.ObjFlags = 0
							if disabled {
								u.ObjFlags = object.Flags(0x8000)
							}
							*u.HealthData = server.HealthData{Cur: hp[0], Field2: 0xabcd, Max: hp[1]}
							for i := 0; i < 32; i++ {
								*(*uint16)(unsafe.Add(originalData, 412+2*i)) = 0x1234
							}
							*(*unsafe.Pointer)(unsafe.Add(originalData, 484)) = nil
							*(*byte)(unsafe.Add(originalData, 1440)) = 0
							if special == 1 {
								*(*unsafe.Pointer)(unsafe.Add(originalData, 484)) = override
							}
							if special == 2 {
								*(*byte)(unsafe.Add(originalData, 1440)) = 0x80
							}
							*o.quest["202024"] = math.Float32bits(difficulty)
							*o.quest["1563932"] = 0
							o.balance(map[string]float64{"GeneratorMaxHealth": 150, "PlayerDamageDiffInit": 1, "SystemHealthDiffInit": 1, "PlayerDamageDiffCoeff": 0.25, "SystemHealthDiffCoeff": 0.5, "PlayerDamageCap": 2, "SystemHealthCap": 3})
							questRuntimeCall("sub_4E3DD0", nil)
							var got row
							got.Name = name
							got.HP = [3]uint16{u.HealthData.Cur, u.HealthData.Field2, u.HealthData.Max}
							got.Damage = *o.quest["202032"]
							got.Health = *o.quest["202036"]
							for i := range got.History {
								got.History[i] = *(*uint16)(unsafe.Add(originalData, 412+2*i))
							}
							unchanged := disabled || hp[0] == 0 || hp[1] == 0 || hp[0] != hp[1] || hp[1] >= 32768 || class == object.ClassSimple || class == object.ClassMonster && special == 2
							if unchanged && (got.HP != [3]uint16{hp[0], 0xabcd, hp[1]}) {
								t.Fatal("ineligible health changed", got.HP)
							}
							if got.HP[1] != 0xabcd {
								t.Fatal("adjacent health word changed")
							}
							if !unchanged {
								if got.HP[0] == 0 || got.HP[2] == 0 {
									t.Fatal("scaled health must remain positive")
								}
								if class == object.Class(0x20000) && (got.HP[0] > 150 || got.HP[2] > 150) {
									t.Fatal("generator cap exceeded")
								}
								if class == object.ClassMonster {
									for _, v := range got.History {
										if v != got.HP[0] {
											t.Fatal("monster health history not refreshed")
										}
									}
								}
							}
							rows = append(rows, got)
						})
					}
				}
			}
		}
	}
}
