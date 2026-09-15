//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestClientMetersZeroMaximumContract(t *testing.T) {
	o := newMeterOwner(t)
	o.plain(t)
	r := &o.meters.Records[0]
	r.Current, r.Maximum = 0, 0
	if got := legacy.PortTestMeterCall(22, r.Window, 0, 0, 0, 0); got != 1 {
		t.Fatal("empty mini-bar return")
	}
	if got := legacy.PortTestMeterCall(36, r.Window, 0, 0, 0, 0); got != 1 {
		t.Fatal("empty tube return")
	}
	if r.Current != 0 || r.Maximum != 0 {
		t.Fatal("drawing changed zero-valued meter")
	}
}
func TestClientMetersLabelMatrix(t *testing.T) {
	o := newMeterOwner(t)
	var out []meterResult
	for id, value := range []uint32{0, 1, 99, 999, 1000, 65535, 0x7fffffff, 0x80000000, 0xffffffff} {
		o.plain(t)
		r := &o.meters.Records[6]
		r.Current = value
		out = append(out, o.invoke(t, id, 0, 19, r.Window, 0, 0, 0, 0))
	}
	meterCapture(t, "meters-labels", out, len(out), "8eb258f374da1bc564213590686f59d68893bbeff761adfa82517e0740d20014")
}
func TestClientMetersBarsMatrix(t *testing.T) {
	o := newMeterOwner(t)
	var out []meterResult
	id := 0
	for _, maximum := range []uint32{0, 1, 2, 31, 100, 255, 65535} {
		for _, current := range []uint32{0, 1, maximum / 2, maximum} {
			for record := 0; record < 2; record++ {
				o.plain(t)
				r := &o.meters.Records[record]
				r.Current, r.Maximum = current, maximum
				*o.meters.NamedWord("nox_client_renderBubbles_80844") = 0
				out = append(out, o.invoke(t, id, 0, 22, r.Window, 0, 0, 0, 0))
				out = append(out, o.invoke(t, id, 1, 36, r.Window, 0, 0, 0, 0))
				id++
			}
		}
	}
	meterCapture(t, "meters-bars", out, len(out), "5c02f7ab2dfcdb0f2070e2b70b3bc2d16ac9cfd73ae777cbd35a1988578060e5")
}
