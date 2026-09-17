//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestCreatureXferTimestamps(t *testing.T) {
	var rows []struct {
		Case                           string
		Input, Delta, Stored, Returned uint32
	}
	defer func() {
		spellbookCapture(t, "creature-xfer-timestamps", rows, "5f3c1e03a96daf814ea8e85d2646a2fb09072184a8058895608dd29c0c83ee71")
	}()
	values := []uint32{0, 1, 2, 10, 0x7ffffffe, 0x7fffffff, 0x80000000, 0x80000001, 0xfffffffe, 0xffffffff}
	for _, value := range values {
		for _, delta := range values {
			t.Run(fmt.Sprintf("value%08x-delta%08x", value, delta), func(t *testing.T) {
				stored, ret := legacy.PortTestCreatureXferAdjust(value, delta)
				want := value + delta
				wantStored := want
				if int32(want) < 1 {
					wantStored = 1
				}
				if stored != wantStored || ret != want {
					t.Fatalf("stored=%08x return=%08x want=%08x/%08x", stored, ret, wantStored, want)
				}
				rows = append(rows, struct {
					Case                           string
					Input, Delta, Stored, Returned uint32
				}{t.Name(), value, delta, stored, ret})
			})
		}
	}
}
