//go:build porttest

package legacy

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/server"
	"math"
	"unsafe"
)

type PortTestUnitRewardSpec struct {
	Mode   int
	Target uint32
	WantXP uint32
}

func (p *portTestShopPools) unitRewardContract() []uint32 {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Controls.UnitReward
	u := p.resources.unit
	coef := memmap.PtrFloat32(0x587000, 206148)
	oldCoef := *coef
	defer func() { *coef = oldCoef }()
	*coef = blobdata.PortTestUnitExperienceCoefficient()
	victim, first, second := p.items[0].u, p.items[1].u, p.items[2].u
	victim.Obj130 = u
	*(*uint32)(unsafe.Add(victim.CObj(), 28)) = sp.Target
	first.ObjOwner = u
	second.ObjOwner = u
	first.ObjClass, second.ObjClass = object.ClassMonster, object.ClassMonster
	first.ObjFlags, second.ObjFlags = 0, 0
	oldFirst, oldSecond := first.UpdateData, second.UpdateData
	defer func() { first.UpdateData, second.UpdateData = oldFirst, oldSecond }()
	first.UpdateData = p.objectiveRegion(int(unsafe.Sizeof(server.MonsterUpdateData{})))
	second.UpdateData = p.objectiveRegion(int(unsafe.Sizeof(server.MonsterUpdateData{})))
	first.UpdateDataMonster().StatusFlags = object.MonStatusSummoned
	second.UpdateDataMonster().StatusFlags = object.MonStatusSummoned
	arg := victim
	switch sp.Mode {
	case 0: // direct player owner
	case 1:
		arg = nil
	case 2:
		victim.Obj130 = nil
	case 3:
		victim.Obj130 = first
		first.ObjOwner = nil
	case 4:
		victim.Obj130 = first
	case 5:
		victim.Obj130 = first
		first.UpdateDataMonster().StatusFlags = 0
	case 6:
		victim.Obj130 = first
		first.ObjFlags = object.FlagDead
	case 7:
		victim.Obj130 = first
		first.ObjClass = object.ClassSimple
	case 8:
		victim.Obj130 = first
		first.ObjOwner = second
		second.UpdateDataMonster().StatusFlags = 0
	case 9:
		victim.Obj130 = first
		first.ObjOwner = second
		first.UpdateDataMonster().StatusFlags = 0
	case 10:
		victim.Obj130 = first
		first.ObjClass = object.ClassSimple
		first.ObjOwner = second
	default:
		panic("reward mode")
	}
	unitMonsterReward(arg)
	got := math.Float32bits(*(*float32)(unsafe.Add(u.CObj(), 28)))
	if got != sp.WantXP {
		panic(fmt.Sprintf("reward XP %08x want %08x", got, sp.WantXP))
	}
	return []uint32{got, uint32(*controlByte(controlPlayer(u), 3684))}
}
