//go:build porttest

package opennox

import (
	"bytes"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
)

func TestQuestRuntimeScoreShippedConstant(t *testing.T) {
	_ = newQuestRuntimeOwner(t)
	raw := unsafe.Slice(memmap.PtrUint8(0x581450, 10088), 12)
	old := bytes.Clone(raw)
	t.Cleanup(func() { copy(raw, old) })
	copy(raw, blobdata.PortTestQuestScoreConstant())
	// With the shipped exponent near1.9, ten points at stage2 produce37 after
	// truncation. This distinguishes a real exponent from a zero-filled fixture
	// and a double constant accidentally read as extended precision.
	got := questRuntimeCall("sub_4D66E0", nil, 1, 0, 0, 2)
	if got != 37 {
		t.Fatalf("shipped exponent: stage2 score=%d want37", got)
	}
}
