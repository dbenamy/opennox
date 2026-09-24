//go:build porttest

package opennox

import (
	"encoding/json"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"reflect"
	"testing"
)

var collisionRegistryOwnerHashes = map[string]string{
	"collision-registry-objectives-ball-pickup":    "a7c48f6cac89cdc32d264e9d7f80f14c4cc25aa5aef3cee2bdaf560a09d9fa08",
	"collision-registry-objectives-crown-dispatch": "c16fc1336e61c1a3169eccf0dad3e9294532b0a462e9e2f9add6213f651ee325",
	"collision-registry-objectives-flag-pickup":    "040aa0d0e738401e3ff9cb4c970a4cfd9c83af0fd46dc9e3c3128c3f391d0162",
	"collision-registry-objectives-home-score":     "b5d9d1e7524578a21f1a2c0664205ed67262a94bdd6d539faa51665c2785e40f",
	"collision-registry-objectives-positive":       "719eeeaea002ca25bb327e9412b6abb8ce5eeb9b4baaf9aaab0718a039ec7d81",
	"collision-registry-pickup":                    "d960decbf8d058cb1dd8a5e1d79f179563347453b241e8bfe0106358b03cfa4b",
	"collision-registry-spell-award":               "2a2ba357c95b24b9b0d8548aba29bcae92b26a79e8e90e6fcc85c23acb726741",
	"collision-registry-state-40":                  "4f409adaf8e7b1620ec98b6dd01bdb56a0e0c4d2f64441190251b7674d78dda9",
	"collision-registry-state-41":                  "48593cb7d5f15c89ae42dc6e3a92690459846489cd53e017490651ebce79568f",
	"collision-registry-state-42":                  "6e8439f3e88273d4313bb48462240ee6586a5c97b8f7b6d165e1b7839804f1b1",
	"collision-registry-trigger":                   "afe829f8057a68e7b1e132ba5e902432e9530e238f71a3a32c859bca66f7c388",
	"collision-registry-world-ankh-history":        "5a8075de055bcbab01d162400c32f4522ce528e861520e6a30198893ae4baca9",
	"collision-registry-world-chest-contents":      "78da915780f503688296b92367d629f4290ffb0059da4c2e4a4fee6ab7e8fc06",
	"collision-registry-world-chest-key":           "ba76f6791f72e9594cb80a051718e25988cf2c8a14a502aad658a4f206c49e5c",
	"collision-registry-world-clock-throttle":      "d58bc83c541e2bf0a765aa9d100848a498590eef2ccf4b424db812c526efa55b",
	"collision-registry-world-coop-exit-save":      "1f0bc07575e272375bd4c0b2638cd38d3b7d6b030f744f38224870c1cba1abb7",
	"collision-registry-world-door-keys":           "64ab2d4e3b6894ab846f90fd72e98ea3f9648105ff96dc038d58d541558d4c30",
	"collision-registry-world-door-magic":          "d41231ffbd2e616fd150710374198423df217da8601731f7b778f215ab346cde",
	"collision-registry-world-door-missing-key":    "ede7ef606d68833b7202b8a7af85698c7da6fc3e26cd26c6ff1c9b16e97727d5",
	"collision-registry-world-exit-admission":      "1ea73cdff61c3f25fbd3abb773a7b8fff5372bfb3862df2f915002985907cb80",
	"collision-registry-world-exit-glyphs":         "2333be7981f9345ab7a40a408a745706355e6eacfe75fde09a3937e28c305ee2",
	"collision-registry-world-quest-exit":          "82478c13ff1a911963f5cae944622e47690f2792bc02210186a321079d2448c5",
	"collision-registry-world-soul-gate":           "b4fe0e9b7260639b8c3fd4e37f697a3a3b525b8456fd194609d5a3e8a25b5b15",
	"collision-registry-world-spell-projectile":    "f2e8d2910d41d5b1658ce2200cd8161d6607744c251574aab8f4375b5cd6272b",
	"collision-registry-world-spell-wall":          "45c7b106a44784966fa0fc49057ff5d0c8ed6ae674ad86a195136a9c542536ac",
	"collision-registry-world-trap-geometry":       "373a3d97827b104b69663daf050ab41fe3fb9424fecb1550064f264af12df6ed",
	"collision-registry-world-undead-damage":       "d392885370f8910928f5673916e451a59125292fc830de9ef8d35f22978285d5",
}

func collisionRegistryCounts(t *testing.T) {
	t.Helper()
	counts := legacy.PortTestCollisionRegistryTakeCounts()
	if len(counts) == 0 {
		t.Fatal("registered collision fixture executed no callbacks")
	}
	b, err := json.Marshal(counts)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("COLLISION_REGISTRY_COUNTS %s", b)
}
func collisionRegistryCapture(t *testing.T, name string, rows any) {
	t.Helper()
	t.Logf("COLLISION_REGISTRY_ROWS %s %d", name, reflect.ValueOf(rows).Len())
	spellbookCapture(t, name, rows, collisionRegistryOwnerHashes[name])
	collisionRegistryCounts(t)
}
func collisionRegistryHash(t *testing.T, name string, rows []legacy.PortTestRoamResult, want string) {
	t.Helper()
	callbackHash(t, name, rows, want)
	collisionRegistryCounts(t)
}
func collisionRegistryWorld(op int, a, b *server.Object, normal *types.Pointf) uint32 {
	names := map[int]string{1: "DoorCollide", 7: "ExitCollide", 8: "SpellProjectileCollide", 9: "ChestCollide", 14: "TrapDoorCollide", 17: "UndeadKillerCollide", 19: "SoulGateCollide", 20: "AnkhCollide"}
	if name := names[op]; name != "" {
		legacy.PortTestRegisteredCollision(a, b, normal, name)
		return 0
	}
	if op == 3 {
		return legacy.PortTestWorldCollision(op, a, b, normal)
	}
	panic("unmapped registered world collision operation")
}
