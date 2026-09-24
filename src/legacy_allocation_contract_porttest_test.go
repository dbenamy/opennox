//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestLegacyAllocationDomains(t *testing.T) {
	type capture struct {
		Count, Size int
		Result      legacy.PortAllocationCapture
	}
	var captures []capture
	for _, count := range []int{1, 2, 3, 7} {
		for _, size := range []int{1, 2, 3, 8, 15, 16, 31, 64, 256} {
			n := count * size
			resizes := []int{n + 33, 1, n + 4096, n}
			got := legacy.PortTestLegacyAllocation(count, size, resizes)
			if !bytes.Equal(got.Initial, make([]byte, n)) {
				t.Fatalf("calloc count=%d size=%d not zero", count, size)
			}
			wantLive := legacy.PortTestAllocationUsesTracker()
			wantDelta := 0
			if wantLive {
				wantDelta = 1
			}
			if got.InitialLive != wantLive || got.InitialCountDelta != wantDelta {
				t.Fatalf("initial domain count=%d size=%d: %+v", count, size, got)
			}
			for stage, step := range got.Steps {
				preserved := n
				if resizes[stage] < preserved {
					preserved = resizes[stage]
				}
				want := make([]byte, preserved)
				for i := range want {
					want[i] = byte(i*37 + stage*53 + 19)
				}
				if step.Size != resizes[stage] || !bytes.Equal(step.Prefix, want) || step.Live != wantLive || step.CountDelta != wantDelta || !step.OldReleased {
					t.Fatalf("realloc count=%d size=%d stage=%d: %+v", count, size, stage, step)
				}
				n = resizes[stage]
			}
			if len(got.Steps) != len(resizes) || got.FinalLive || got.FinalCountDelta != 0 {
				t.Fatalf("release count=%d size=%d: %+v", count, size, got)
			}
			captures = append(captures, capture{count, size, got})
		}
	}
	data, err := json.Marshal(captures)
	if err != nil {
		t.Fatal(err)
	}
	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	t.Logf("legacy allocation domains capture=%s cases=%d", hash, len(captures))
	want := "e0b7a4b9ab7fff0890a9ab38b3274c1ffe4749cbd7125f29c273b25553dc8525"
	if legacy.PortTestAllocationUsesTracker() {
		want = "a868e1790dd29a7291feb4e97d834084734117b244741cdb6981efcc8e56e077"
	}
	if hash != want {
		t.Fatalf("allocation contract hash %s want %s", hash, want)
	}
}

func TestLegacyCStringOwnership(t *testing.T) {
	for _, input := range []string{"", "a", "\x00", "\x80\xffx", "alpha\x00omega", string(bytes.Repeat([]byte{0x9a}, 255))} {
		got := legacy.PortTestLegacyCStringOwnership(input)
		wantLive := legacy.PortTestAllocationUsesTracker()
		wantDelta := 0
		if wantLive {
			wantDelta = 1
		}
		if !bytes.Equal(got.Initial, append([]byte(input), 0)) || got.InitialLive != wantLive || got.InitialCountDelta != wantDelta || got.FinalLive || got.FinalCountDelta != 0 {
			t.Fatalf("CString/StrFree input=%x: %+v", input, got)
		}
	}
}
