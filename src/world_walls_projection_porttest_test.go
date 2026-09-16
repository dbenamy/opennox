//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"os"
	"testing"
)

func worldWallsCapture(t *testing.T, label string, rows any, want string) {
	t.Helper()
	raw, err := json.Marshal(rows)
	if err != nil {
		t.Fatal(err)
	}
	if prefix := os.Getenv("OPENNOX_WORLD_WALLS_CAPTURE"); prefix != "" {
		if err := os.WriteFile(prefix+"-"+label+".json", raw, 0600); err != nil {
			t.Fatal(err)
		}
	}
	got := fmt.Sprintf("%x", sha256.Sum256(raw))
	t.Logf("%s: %s", label, got)
	// Expectations were captured from original production C.
	if want != got {
		t.Fatalf("%s hash %s want %s", label, got, want)
	}
}
func TestWorldWallsProjection(t *testing.T) {
	type result struct {
		Op                           int
		Screen, World, Input, Output image.Point
		Return                       int
	}
	var rows []result
	for _, screen := range []image.Point{{0, 0}, {8, 12}, {-17, 23}, {2147483647, -2147483648}} {
		for _, world := range []image.Point{{0, 0}, {182, 182}, {-100, -200}, {-2147483648, 2147483647}} {
			vp := noxrender.Viewport{Screen: image.Rectangle{Min: screen}, World: image.Rectangle{Min: world}}
			for _, p := range []image.Point{{0, 0}, {-1, 1}, {500, -500}, {2147483647, -2147483648}} {
				for op := 0; op < 2; op++ {
					ret, out := legacy.PortTestWorldWalls(op, &vp, nil, nil, p)
					want := vp.ToScreenPos(p)
					if op == 1 {
						want = vp.ToWorldPos(p)
					}
					if ret != p.Y || out != want {
						t.Fatalf("op%d input%v -> %v,%d want%v,%d", op, p, out, ret, want, p.Y)
					}
					rows = append(rows, result{op, screen, world, p, out, ret})
				}
			}
		}
	}
	worldWallsCapture(t, "projection", rows, "ddce5d8f5e2b8f4bc3e624561834096e511e8a3eff1925e7a57872c0a550f214")
}
