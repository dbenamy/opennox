//go:build porttest

package opennox

import (
	"fmt"
	"math"
	"testing"
)

func TestQuestRuntimeFloatScalars(t *testing.T) {
	o := newQuestRuntimeOwner(t)
	type row struct {
		Name         string
		Return, Read uint64
		Stored       uint32
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "quest-runtime-float-scalars", rows, "8341aa3435c35b42aa110ff9efcaa28baba40ea1e8a2f0a6c0cb813a6b66322f")
	}()
	bits := []uint32{0, 0x80000000, 1, 0x80000001, 0x3f000000, 0x3f800000, 0xbf800000, 0x4b800001, 0x7f7fffff, 0xff7fffff, 0x7f800000, 0xff800000, 0x7fc01234}
	for _, cap := range []float64{-1, 0, 1, 1.00000001, 10000} {
		o.balance(map[string]float64{"PlayerDamageCap": cap, "SystemHealthCap": cap})
		for _, b := range bits {
			for _, p := range [][3]string{{"sub_4E3CB0", "sub_4E3CA0", "202024"}, {"sub_4E4080", "sub_4E40B0", "202032"}, {"sub_4E40C0", "sub_4E40F0", "202036"}} {
				name := fmt.Sprintf("%s/cap%g/bits%x", p[0], cap, b)
				t.Run(name, func(t *testing.T) {
					rv := questRuntimeCall(p[0], nil, b)
					rd := questRuntimeCall(p[1], nil)
					stored := *o.quest[p[2]]
					f := math.Float32frombits(b)
					want := f
					if p[0] != "sub_4E3CB0" && !(float64(f) <= cap) {
						want = float32(cap)
					}
					if !math.IsNaN(float64(want)) && stored != math.Float32bits(want) {
						t.Fatalf("stored%x want%x", stored, math.Float32bits(want))
					}
					if p[0] == "sub_4E3CB0" && rv != uint64(b) {
						t.Fatalf("float-bit setter return%x want%x", rv, b)
					}
					if !math.IsNaN(float64(want)) && rd != math.Float64bits(float64(want)) {
						t.Fatal("getter precision")
					}
					rows = append(rows, row{name, rv, rd, stored})
				})
			}
		}
	}
}
func TestQuestRuntimeStatisticResetAll(t *testing.T) {
	o := newQuestRuntimeOwner(t)
	type row struct {
		Name  string
		Stats [3][11]uint32
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "quest-runtime-reset-all", rows, "9b08cbbea8ea881f25a4573fb260909c09d61a031ad028000eb1b6720db032c2")
	}()
	for mask := 0; mask < 8; mask++ {
		for _, stage := range []uint32{0, 1, 65536, 0xffffffff} {
			name := fmt.Sprintf("mask%d/stage%x", mask, stage)
			t.Run(name, func(t *testing.T) {
				*o.quest["202028"] = stage
				var want [3][11]uint32
				for i := range o.units {
					u := &o.units[i]
					pl := u.UpdateDataPlayer().Player
					pl.Active = 0
					pl.PlayerUnit = nil
					for j := 0; j < 11; j++ {
						objectXferSetWord(pl.C(), 4652+4*j, uint32(100+i+j))
					}
					want[i] = questRuntimeStats(pl.C())
					if mask&(1<<uint(i)) != 0 {
						pl.Active = 1
						pl.PlayerUnit = u
						want[i] = [11]uint32{}
						want[i][9] = stage
						want[i][10] = 63
					}
				}
				if rv := questRuntimeCall("sub_4D60B0", nil); rv != 0 {
					t.Fatal("reset-all result", rv)
				}
				var got [3][11]uint32
				for i := range o.units {
					got[i] = questRuntimeStats(o.units[i].UpdateDataPlayer().Player.C())
				}
				if got != want {
					t.Fatal("reset-all scope", got, want)
				}
				rows = append(rows, row{name, got})
			})
		}
	}
}
