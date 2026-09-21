//go:build porttest

package opennox

import (
	"encoding/binary"
	"testing"

	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestClientRandomNames(t *testing.T) {
	s := new(server.Server)
	old := legacy.GetServer
	t.Cleanup(func() { legacy.GetServer = old })
	legacy.GetServer = func() legacy.Server { return &audioEventRandomOwner{s: s} }
	count := serverConfigOwnBytes(t, 0x5D4594, 814516, 4)
	names := []string{"Dweezle", "Glork", "Floogle", "Goombah", "Kraun", "Kloog", "Zurg", "Darg", "Arfingle", "Buurl", "Gurgin", "Grok", "Hurlong", "Luric", "Lupis", "Mallik", "Thrall", "Norwood", "Nulik", "Orin", "Olaf", "Orguk", "Pervis", "Paavik", "Qix", "Xevin", "Xurcon", "Markoan", "Yuric", "Yoovis", "Yalek", "Zug", "Zivik"}
	type row struct {
		Seed, Initial, Step, Index int
		Count                      uint32
		Name                       string
	}
	var rows []row
	seen := map[string]bool{}
	for _, seed := range []int{0, 1, 17, 31, 0x7fffffff} {
		for _, initial := range []int{0, 1, 2, 17, 33} {
			s.Rand.Logic, s.Rand.Other = prand.New(7), prand.New(seed)
			reference := prand.New(seed)
			binary.LittleEndian.PutUint32(count, uint32(initial))
			limit := initial
			if limit == 0 {
				limit = 33
			}
			for step := 0; step < 128; step++ {
				want := names[reference.Int(0, limit-1)]
				got := legacy.Nox_xxx_getRandomName_4358A0()
				if got != want || s.Rand.Other.Index() != reference.Index() || s.Rand.Logic.Index() != prand.New(7).Index() || binary.LittleEndian.Uint32(count) != uint32(limit) {
					t.Fatal("random name/cache/stream", seed, initial, step, got, want)
				}
				if initial == 0 {
					seen[got] = true
				}
				rows = append(rows, row{seed, initial, step, int(reference.Index()), uint32(limit), got})
			}
		}
	}
	if len(seen) != 33 {
		t.Fatal("name table coverage", len(seen))
	}
	interactionCapture(t, "client-random-names", rows)
}
