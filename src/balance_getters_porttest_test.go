//go:build porttest

package opennox

import (
	"math"
	"testing"

	"github.com/opennox/libs/balance"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

type balanceGetterServer struct {
	legacy.Server
	s *server.Server
}

func (o *balanceGetterServer) S() *server.Server { return o.s }

func TestBalanceGetterCAndGoContracts(t *testing.T) {
	s := new(server.Server)
	negativeZero := math.Copysign(0, -1)
	preciseDouble := math.Nextafter(1, 2)
	restoreBalance := s.PortTestBalanceOverlay(
		balance.Config{
			"globalonly":    balance.Float(6.75),
			"scalarcase":    balance.Float(999.5), // tagged values override this
			"precisedouble": balance.Float(preciseDouble),
			"widedouble":    balance.Float(1e100),
			"negativezero":  balance.Float(negativeZero),
			"globaltable":   balance.Array{3.125, -7.5, 0.0},
			"emptytable":    balance.Array{},
		},
		balance.Config{
			"scalarcase": balance.Float(19.375),
			"tagtable":   balance.Array{1.25, -2.5, 0, 1e100},
		},
		balance.Config{
			"scalarcase": balance.Float(-0.125),
			"tagtable":   balance.Array{negativeZero, 8.75},
		},
	)
	t.Cleanup(restoreBalance)
	oldGetServer := legacy.GetServer
	legacy.GetServer = func() legacy.Server { return &balanceGetterServer{s: s} }
	t.Cleanup(func() { legacy.GetServer = oldGetServer })
	oldFlags := flags.GetGame()
	t.Cleanup(func() {
		flags.ResetGame()
		flags.SetGame(oldFlags)
	})

	type check struct {
		name       string
		index      int
		wantScalar float64
		wantIndex  float64
	}
	for _, mode := range []struct {
		name   string
		coop   bool
		checks []check
	}{
		{
			name: "arena",
			checks: []check{
				{"scalarcase", 0, 19.375, 19.375},
				{"scalarcase", 2, 19.375, 0},
				{"globalonly", 0, 6.75, 6.75},
				{"globaltable", 0, 3.125, 3.125},
				{"globaltable", 1, 3.125, -7.5},
				{"globaltable", 2, 3.125, 0},
				{"globaltable", 3, 3.125, 0},
				{"tagtable", 0, 1.25, 1.25},
				{"tagtable", 1, 1.25, -2.5},
				{"tagtable", 2, 1.25, 0},
				{"tagtable", 3, 1.25, 1e100},
				{"tagtable", 4, 1.25, 0},
				{"emptytable", 0, 0, 0},
				{"missingkey", 0, 0, 0},
				{"widedouble", 0, 1e100, 1e100},
				{"precisedouble", 0, preciseDouble, preciseDouble},
				{"ScAlArCaSe", 0, 19.375, 19.375},
				{"", 0, 0, 0},
				{"negativezero", 0, negativeZero, negativeZero},
			},
		},
		{
			name: "solo",
			coop: true,
			checks: []check{
				{"scalarcase", 0, -0.125, -0.125},
				{"scalarcase", 1, -0.125, 0},
				{"tagtable", 0, negativeZero, negativeZero},
				{"tagtable", 1, negativeZero, 8.75},
				{"tagtable", 2, negativeZero, 0},
				{"tagtable", 3, negativeZero, 0},
				{"globalonly", 0, 6.75, 6.75},
			},
		},
	} {
		flags.ResetGame()
		if mode.coop {
			flags.SetGame(flags.GameModeCoop)
		}
		for _, tc := range mode.checks {
			for _, idx := range []int{tc.index, -1, int(-1 << 31), int(1<<31 - 1)} {
				want := tc.wantIndex
				if idx != tc.index {
					want = 0 // all these indices are outside each test table
				}
				got := legacy.PortTestBalanceGetters(tc.name, idx)
				if math.Float64bits(got.ScalarC) != math.Float64bits(tc.wantScalar) ||
					math.Float64bits(got.ScalarGo) != math.Float64bits(tc.wantScalar) ||
					math.Float64bits(got.IndexC) != math.Float64bits(want) ||
					math.Float64bits(got.IndexGo) != math.Float64bits(want) {
					t.Fatalf("mode=%s key=%q index=%d got C(%.17g,%.17g) Go(%.17g,%.17g), want scalar %.17g index %.17g",
						mode.name, tc.name, idx, got.ScalarC, got.IndexC, got.ScalarGo, got.IndexGo, tc.wantScalar, want)
				}
			}
		}
	}
}
