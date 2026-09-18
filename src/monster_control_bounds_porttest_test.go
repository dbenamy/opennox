//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"strings"
	"testing"
)

// The original fixture preserves defined field layouts. This independent Go
// contract checks the new bounded-input policy and rejected-record ownership;
// it deliberately does not run overlong fields through the original C parser.
func TestMonsterControlDefinitionBoundsAndCleanup(t *testing.T) {
	monsterControlDefsOwner(t)
	t.Chdir(t.TempDir())
	base := alloc.PortTestAllocationCount()
	cases := []struct {
		text string
		want int
	}{
		{strings.Repeat("A", 63) + " END ", 1},
		{strings.Repeat("A", 64) + " END ", 0},
		{"PortMonsterA MISSILE_NAME " + strings.Repeat("x", 63) + " HEALTH 17 END ", 1},
		{"PortMonsterA MISSILE_NAME " + strings.Repeat("x", 64) + " HEALTH 17 END ", 0},
		{"PortMonsterA FLEE_RANGE " + strings.Repeat("1", 256) + " END ", 0},
		{"PortMonsterA END PortMonsterB MISSILE_NAME " + strings.Repeat("x", 64) + " END ", 1},
		{"PortMonsterA MELEE_STRIKE_FUNCTION unknown END ", 0},
		{"PortMonsterA DIE_FUNCTION unknown END ", 0},
		{"PortMonsterA DEAD_FUNCTION unknown END ", 0},
		{"PortMonsterA MELEE_ATTACK_DAMAGE_TYPE unknown END ", 0},
		{"PortMonsterA UNKNOWN 17 END ", 0},
	}
	for i, c := range cases {
		monsterControlWriteBin(t, c.text)
		if legacy.PortTestMonsterDefs("load", 0) != 1 {
			t.Fatal("opened file rejected at outer layer", i)
		}
		got := legacy.PortTestMonsterDefsSnapshot()
		if len(got) != c.want {
			t.Fatal("bounded definition policy", i, len(got), c.want)
		}
		if i == 2 && got[0][17] != 17 {
			t.Fatal("maximum-length missile name overlapped numeric fields")
		}
		if count := alloc.PortTestAllocationCount(); count != base+c.want {
			t.Fatal("rejected definition leaked owned allocation", i, count-base, c.want)
		}
		legacy.PortTestMonsterDefs("free", 0)
		if count := alloc.PortTestAllocationCount(); count != base {
			t.Fatal("definition cleanup leak", i, count-base)
		}
	}
}
