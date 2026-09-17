//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"runtime"
	"testing"
	"unsafe"
)

// Keep the fixture on the Go heap. Its address is used only as opaque integer
// bits; the production curve callback must never interpret those bits as a pointer.
var curveTokenFixture *[2]*int

func TestClientEffectsCurveOpaqueToken(t *testing.T) {
	old := curveTokenFixture
	payload := new([2]*int)
	payload[0], payload[1] = new(int), new(int)
	curveTokenFixture = payload
	t.Cleanup(func() { curveTokenFixture = old })
	for _, token := range []int32{0, -1, -2147483648, int32(uintptr(unsafe.Pointer(payload)))} {
		got := legacy.PortTestEffectsCurve(legacy.PortTestEffectsCurveSpec{Points: [4][2]int32{{0, 0}, {64, 32}, {0, 0}, {0, 0}}, Steps: 2, Token: token})
		if len(got) != 2 || got[0].Token != token || got[1].Token != token {
			t.Fatal("opaque callback token changed", token, got)
		}
		if got[0].From != ([2]int32{0, 0}) || got[1].To != ([2]int32{64, 32}) {
			t.Fatal("curve endpoints changed")
		}
	}
	runtime.KeepAlive(payload)
}
