//go:build porttest

package opennox

import (
	"context"
	"github.com/opennox/opennox/v1/legacy"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestGridBoundsRejectInvalid(t *testing.T) {
	if raw := os.Getenv("OPENNOX_PORT_GRID_BOUNDS_CHILD"); raw != "" {
		values := strings.Split(raw, ",")
		x, e1 := strconv.ParseUint(values[0], 16, 32)
		y, e2 := strconv.ParseUint(values[1], 16, 32)
		if e1 != nil || e2 != nil {
			t.Fatal("child input")
		}
		if got := legacy.PortTestGridBounds(uint32(x), uint32(y)); got != -1 {
			t.Fatalf("expected invalid tile -1, got %d", got)
		}
		return
	}
	// Ordinary invalid coordinates reject without touching a grid.
	for _, pair := range [][2]uint32{{0, 0}, {0xc47a0000, 0x43000000}, {0x43000000, 0xc47a0000}, {0x461c4000, 0x43000000}} {
		if got := legacy.PortTestGridBounds(pair[0], pair[1]); got != -1 {
			t.Fatalf("ordinary bounds %x: %d", pair, got)
		}
	}
	// These inputs convert to INT32_MIN. Reject them before any grid access.
	for _, raw := range []string{"7fc12345,43000000", "43000000,7fc12345", "7f800000,43000000", "43000000,ff800000", "7f7fffff,43000000", "ff7fffff,43000000"} {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestGridBoundsRejectInvalid$", "-test.count=1")
		cmd.Env = append(os.Environ(), "OPENNOX_PORT_GRID_BOUNDS_CHILD="+raw)
		_, err := cmd.CombinedOutput()
		timedOut := ctx.Err() != nil
		cancel()
		if timedOut || err != nil {
			t.Fatalf("expected safe bounds rejection for %s (timeout=%v err=%v)", raw, timedOut, err)
		}
	}
}
