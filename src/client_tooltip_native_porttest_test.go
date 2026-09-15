//go:build porttest

package opennox

import (
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestClientTooltipNativeWrappersContract(t *testing.T) {
	o := newTooltipOwner(t)
	o.prepare(0, 0, 0, 0)
	if got := legacy.Nox_xxx_clientAskInfoMb_4BF050(o.dr); got != "Pretty" {
		t.Fatalf("Go name wrapper got%q", got)
	}
	legacy.Nox_xxx_cursorSetTooltip_4776B0("λ")
	if o.cursor[0] != 'λ' || o.cursor[1] != 0 || o.cursor[2] != 0x9102 {
		t.Fatal("Go cursor wrapper changed text/tail")
	}
	legacy.Nox_xxx_cursorSetTooltip_4776B0("")
	if o.cursor[0] != 0 || o.cursor[2] != 0x9102 {
		t.Fatal("Go cursor clear changed tail")
	}
}

// The former unbounded C concatenation has no valid execution for a name that
// exceeds its mapped scratch allocation. The native assembler truncates it.
func TestClientTooltipNativeCapacityContract(t *testing.T) {
	o := newTooltipOwner(t)
	text := make([]uint16, 2048)
	for i := range text {
		text[i] = uint16(0xd800 + i%0x700)
	}
	p := tooltipWide(t, text)
	o.weapon.Desc8 = &p[0]
	o.prepare(0x1000, 0, 0, 15)
	if got := legacy.PortTestTooltip(o.dr); got != &o.scratch[0] {
		t.Fatal("oversized name returned different storage")
	}
	for i := 0; i < 1023; i++ {
		// Default order is M0 M1 followed by the oversized base name.
		want := uint16(0)
		if i < 6 {
			want = []uint16{'M', '0', ' ', 'M', '1', ' '}[i]
		} else {
			want = text[i-6]
		}
		if o.scratch[i] != want {
			t.Fatalf("truncated unit%d changed", i)
		}
	}
	if o.scratch[1023] != 0 {
		t.Fatal("oversized assembled name is not terminated")
	}
}
