//go:build porttest

package opennox

import (
	"bytes"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func gameplayReportsBase(op int, args ...legacy.PortTestGameplayReportArg) legacy.PortTestRoamSpec {
	s := controlsBase(0)
	p := s.Callbacks.Shop
	sp := &legacy.PortTestGameplayReportsSpec{}
	copy(sp.Args[:], args)
	p.TemporaryUpdates.World.Objectives.Attack.Controls.Reports = sp
	p.Sequence = []legacy.PortTestShopAction{{Op: 1800 + op}}
	return s
}
func reportValue(v uint32) legacy.PortTestGameplayReportArg {
	return legacy.PortTestGameplayReportArg{Value: v}
}
func reportObject(ref int) legacy.PortTestGameplayReportArg {
	return legacy.PortTestGameplayReportArg{Kind: "object", Ref: ref}
}

func TestGameplayReportsCreatureOrderProbe(t *testing.T) {
	s := gameplayReportsBase(3, reportValue(7), reportValue(5))
	out := controlsRun(t, []legacy.PortTestRoamSpec{s})[0]
	if !out.Intact || !out.Callbacks.Shop.Intact {
		t.Fatal("gameplay report guards")
	}
	step := out.Callbacks.Shop.Sequence[0]
	if step.GameplayReports == nil {
		t.Fatal("missing report snapshot")
	}
	found := 0
	for _, p := range step.Packets {
		if bytes.Equal(p.Data, []byte{237, 5}) {
			if p.Recipient != 7 {
				t.Fatalf("recipient %d", p.Recipient)
			}
			found++
		}
	}
	if found != 1 {
		t.Fatalf("creature-order messages %d, total packets %d", found, len(step.Packets))
	}
	t.Logf("original C creature order: one message, recipient7, return%d", step.Return)
}
