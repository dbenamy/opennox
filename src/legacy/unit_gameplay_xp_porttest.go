//go:build porttest

package legacy

import (
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"math"
	"unsafe"
)

type PortTestUnitExperienceSpec struct {
	Presentation        bool
	MissingPresentation bool
	WantAward           *uint32
	Op                  int
	Amount, WantXP      uint32
	WantLevel           byte
	Saving, HasSave     bool
}

func (p *portTestShopPools) unitExperienceContract() []uint32 {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Controls.UnitExperience
	u := p.resources.unit
	checkPresentation, restorePresentation := p.unitExperiencePresentation(sp)
	defer restorePresentation()
	oldSaving, oldSave := Nox_xxx_gameGet_4DB1B0, Sub_4DB1C0
	defer func() { Nox_xxx_gameGet_4DB1B0, Sub_4DB1C0 = oldSaving, oldSave }()
	Nox_xxx_gameGet_4DB1B0 = func() bool { return sp.Saving }
	Sub_4DB1C0 = func() unsafe.Pointer {
		if sp.HasSave {
			return u.CObj()
		}
		return nil
	}
	coef := memmap.PtrFloat32(0x587000, 206148)
	oldCoef := *coef
	defer func() { *coef = oldCoef }()
	*coef = blobdata.PortTestUnitExperienceCoefficient()
	if math.Float32bits(*coef) != 0x3c23d70a {
		panic("shipped XP coefficient changed")
	}
	before := *(*float32)(unsafe.Add(u.CObj(), 28))
	amount := math.Float32frombits(sp.Amount)
	result := float64(0)
	switch sp.Op {
	case 0:
		unitExperienceLevel(u)
	case 1:
		unitGiveExperience(u, amount)
	case 2:
		result = unitRewardExperience(u, amount)
		want := float64(0)
		if before < amount {
			want = float64(float32(float64(float32(amount-before))*float64(float32(0.01)) + 1))
		}
		if sp.WantAward != nil {
			want = float64(math.Float32frombits(*sp.WantAward))
		}
		if result != want {
			panic(fmt.Sprintf("XP return %g want %g", result, want))
		}
	default:
		panic("unknown XP operation")
	}
	got := math.Float32bits(*(*float32)(unsafe.Add(u.CObj(), 28)))
	level := *controlByte(controlPlayer(u), 3684)
	if got != sp.WantXP || level != sp.WantLevel {
		panic(fmt.Sprintf("XP state %08x/%d want %08x/%d", got, level, sp.WantXP, sp.WantLevel))
	}
	checkPresentation()
	bits := math.Float64bits(result)
	return []uint32{got, uint32(level), uint32(bits), uint32(bits >> 32)}
}
