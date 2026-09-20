//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func gameMessageStateGuards(t *testing.T, rows []legacy.PortTestRoamResult) {
	t.Helper()
	for i, r := range rows {
		if !r.Intact || !r.Callbacks.Intact || !r.Callbacks.Shop.Intact || !r.Spells.Intact || !r.Main.Intact || !r.MonsterState.Intact || !r.Combat.Intact {
			t.Fatalf("game-message state case %d changed a fixture guard", i)
		}
	}
}
