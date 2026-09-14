//go:build porttest

package legacy

import (
	"testing"
	"unsafe"
)

func TestMapOrchestrationDiagnosticIdentity(t *testing.T) {
	var engine, diagnostic [4]uint32
	f := &mapRoomTestFixture{}
	r := f.register(unsafe.Pointer(&engine[0]), 16, "engine", false)
	d := f.register(unsafe.Pointer(&diagnostic[0]), 16, "diagnostic", false)
	d.captureOnly = true
	for off := uint32(0); off < 16; off += 4 {
		raw := uint32(uintptr(d.ptr)) + off
		if got := f.normalize(raw); got != raw {
			t.Fatalf("diagnostic address reinterpreted as engine reference: %d", off)
		}
		if got := f.normalize(uint32(uintptr(r.ptr)) + off); got != r.id+off {
			t.Fatalf("engine reference lost: %d", off)
		}
	}
}
