//go:build porttest

package opennox

import (
	"encoding/json"
	"os"
	"testing"
)

func meterCapture(t *testing.T, label string, out []meterResult, count int, want string) {
	t.Helper()
	if prefix := os.Getenv("OPENNOX_METERS_CAPTURE"); prefix != "" {
		data, err := json.Marshal(out)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(prefix+"-"+label+".json", data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	effectsCapture(t, label, out, count, want)
}
