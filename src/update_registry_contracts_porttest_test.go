//go:build porttest

package opennox

import (
	"encoding/json"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

// Frozen from actual registered Object.CallUpdate owners before typed dispatch.
var updateRegistryHashes = map[string]string{
	"update-registry-generator-update-gates":         "dc753db5ded7d537d80f66fc0595f4db5cff90a276aa4407a5ce581fe96cc6d1",
	"update-registry-generator-update-spawn":         "299330844c6096df687b2b20a47651027a57215d371f8391607afe37658a98fe",
	"update-registry-motion-mover":                   "77cfed5e8d3e5e56b21b9d47c3d6138d18280dba1f5074bade9c9f059010fa6d",
	"update-registry-motion-sentry":                  "323db817d302737613de329371e22f227fd1706c904f5223d116add67d573140",
	"update-registry-motion-trap":                    "123545df8a5e94b445e1f189cb290e5f16cb785f31cf35704ee1dee9683d70d2",
	"update-registry-objectives-ball-clock":          "4df5e94496eae0e20a156cfab61f11e7c2ebce58cde5dbd6011ab95f7d83e0a9",
	"update-registry-objectives-crown-dispatch":      "0c24d4eed559bb7729333f54123a728142509485e1a271e95b5cf693d4c4da13",
	"update-registry-objectives-flag-deadline":       "92435e14ce02529ec368bb880d29ecab16ee47c20ec799a0771ec16d987a2933",
	"update-registry-objectives-obelisk-eligibility": "7ceff3a1bdd066ceb085e712aef1e98976df8341ebd7267a8b38f920810d5ecf",
	"update-registry-objectives-obelisk-idle":        "e402b16f9e6e8586e97b04f02d5f6247081a27c7da044eeb01f4104497381620",
	"update-registry-objectives-obelisk-transfer":    "b5c46594f05d2169d1908c348391093461b3735a45e48920b1faf2a96c2396b0",
	"update-registry-objectives-obelisk-wands":       "b069bc7f664d71d6a41153f8ceb3cbf37c86be63c6036208030cc71981e01d37",
	"update-registry-temporary-acquisition":          "0e8a000ce16de999cc34356f2c3bf17164ef7a8010053759a422459046497a3e",
	"update-registry-temporary-area-owners":          "5fcd97779256069eb4a68e4b8dad4b0ffc61464127c882c39a395318790be90d",
	"update-registry-temporary-break":                "de8371d8c52f0ab8e7f6900936990185a1f9de28df3cd0c531bc597d5cd6a0e2",
	"update-registry-temporary-frame-wrap":           "3c31aafcd7142fed084829a510e8c1cf4ba09dbcbc3a58bfd8974e0dada0ba04",
	"update-registry-temporary-homing":               "0dec953e506ba8dd9f5fb1b05a051cd09e9399683feda90a45d442eb8e2fff67",
	"update-registry-temporary-lifetime":             "f1050e32749cb74765493399f56ebfbfde939d07bafceefec60e0b8fe9e4f97a",
	"update-registry-temporary-missing-owners":       "38bb8ad4d8aef15a62dafa008a591a85db32faba52ad90c814888bc9fc8a5fef",
	"update-registry-temporary-owner-effects":        "b895bccccd6dc134cb30e984d3fa1822e4e09f7f00c4744db9ea377e83748761",
	"update-registry-temporary-random-streams":       "e54a9167c141e754c25246cc72122392235635034234612a3c7f236191947b90",
	"update-registry-temporary-spawners":             "1a25e476641e8b8cc9576520a5f01262696a790940c7b1c7199e9233f8f65e81",
	"update-registry-temporary-tick-rates":           "e208d2de7e7c5290603b25d7ca772a30d697bf51dc8eb1aa06c7f3063a3b7fcf",
	"update-registry-temporary-trails":               "e92e72b5d473c021ee05f216bd328e713a64b3bdd39856db7f43cc16a273badf",
	"update-registry-world-doors":                    "4ab11c2923a1a1fa31346a1092705fc12a9dae948613a40db6f1f9ed3115fb81",
	"update-registry-world-elevators":                "4eac9076eaf4b663ae9b06235c0cb6a86be6239774f9d31318ab246efe972d7a",
	"update-registry-world-forces":                   "2c439ed43e78719e38493dd0215615aae78c5dd67e7f07373c8afad8331ebfad",
	"update-registry-world-invisible-teleport":       "4693a945e6fa2e7885614c1b8431a2ace733c4a3fe64250deb01200b25ef8db8",
	"update-registry-world-owner-integration":        "005d9ea7898022dbe5684d28f152aa2ca82680e75b220b93681d38a72a9705b9",
	"update-registry-world-phantom-queue":            "9c4199a101a109cc23d1092d3efc938d69aeb71743764aa365a0b5e119818b0a",
	"update-registry-world-positive":                 "5206010f166eb46098b2d9e41dd7f07bf095ff969abf81c48e0e3d086ada47f8",
	"update-registry-world-switches":                 "885a594c122f4e3082108d5d2fffc5cf88fc01dd1025abcdb6528de8ced8e0c5",
	"update-registry-world-teleports":                "65cbd9ca486446f59e4fcbcbd3ae1bc7d80b194699420f36b40b17a4e60eab7b",
	"update-registry-world-trap-doors":               "5d77f69f9234d64eeb78da5934d34b099f910221b1090e70f1ae2c26287de022",
	"update-registry-world-triggers":                 "227329d562aee47313b52dc2486c83ac850a2d01240503d0f0962c20a4484607",
	"update-registry-world-type-clock":               "c9c7c5010ace7b9fb092b01d216e31ff775bbf8035fd47bea15e01a8df09c2fc",
}

func updateRegistryTemporaryBase() legacy.PortTestRoamSpec {
	s := temporaryBase()
	s.Callbacks.Shop.TemporaryUpdates.RegisteredUpdates = true
	return s
}
func updateRegistryWorldBase() legacy.PortTestRoamSpec {
	s := worldBase()
	s.Callbacks.Shop.TemporaryUpdates.RegisteredUpdates = true
	return s
}
func updateRegistryObjectiveBase() legacy.PortTestRoamSpec {
	s := objectiveBase()
	s.Callbacks.Shop.TemporaryUpdates.RegisteredUpdates = true
	return s
}
func updateRegistryReportCounts(t *testing.T) {
	t.Helper()
	counts := legacy.PortTestUpdateRegistryTakeCounts()
	if len(counts) == 0 {
		t.Fatal("registered update fixture did not execute any callbacks")
	}
	b, err := json.Marshal(counts)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("UPDATE_REGISTRY_COUNTS %s", b)
}
func updateRegistryHash(t *testing.T, name string, r []legacy.PortTestRoamResult, want string) {
	t.Helper()
	callbackHash(t, name, r, want)
	updateRegistryReportCounts(t)
}
func updateRegistrySpecsHash(t *testing.T, name string, s []legacy.PortTestRoamSpec) {
	t.Helper()
	updateRegistryHash(t, name, effectsTimedRun(t, s), updateRegistryHashes[name])
}

func TestUpdateRegistryNames(t *testing.T) {
	names := legacy.PortTestUpdateRegistryNames()
	if len(names) != 53 {
		t.Fatalf("registrations %d want53", len(names))
	}
	seen := map[string]bool{}
	for _, name := range names {
		if seen[name] {
			t.Fatal("duplicate registration", name)
		}
		seen[name] = true
	}
}
func TestUpdateRegistryRawForwarding(t *testing.T) {
	u, free := alloc.New(server.Object{})
	defer free()
	u.ObjFlags = 0x40
	callback := legacy.PortTestDeathForwardReset()
	u.CallUpdate()
	ptr, count := legacy.PortTestDeathForwardSnapshot()
	if ptr != nil || count != 0 {
		t.Fatal("nil update dispatched")
	}
	u.Update = callback
	u.CallUpdate()
	ptr, count = legacy.PortTestDeathForwardSnapshot()
	if ptr != u.CObj() || count != 1 || u.ObjFlags != 0x40 {
		t.Fatalf("raw update %p/%d/%x", ptr, count, u.ObjFlags)
	}
	u.Update = nil
	u.CallUpdate()
	ptr, count = legacy.PortTestDeathForwardSnapshot()
	if ptr != u.CObj() || count != 1 {
		t.Fatal("cleared slot dispatched")
	}
}
func TestUpdateRegistryDynamicHandlers(t *testing.T) {
	u, free := alloc.New(server.Object{})
	defer free()
	for _, tc := range []struct {
		name string
		slot *func(*server.Object)
	}{
		{"PlayerUpdate", &legacy.Nox_xxx_updatePlayer_4F8100},
		{"ProjectileUpdate", &legacy.Nox_xxx_updateProjectile_53AC10},
		{"MonsterUpdate", &legacy.Nox_xxx_unitUpdateMonster_50A5C0},
		{"PixieUpdate", &legacy.Nox_xxx_updatePixie_53CD20},
		{"DeathBallUpdate", &legacy.Nox_xxx_updateDeathBall_53D080},
		{"HarpoonUpdate", &legacy.Nox_xxx_updateHarpoon_54F380},
	} {
		t.Run(tc.name, func(t *testing.T) {
			old := *tc.slot
			defer func() { *tc.slot = old }()
			calls := 0
			for _, value := range []uint32{17, 91} {
				want := value
				*tc.slot = func(got *server.Object) {
					if got != u {
						t.Fatal("wrong update object")
					}
					calls++
					got.Field34 = want
				}
				legacy.PortTestRegisteredUpdate(u, tc.name)
				if u.Field34 != value {
					t.Fatalf("handler value%d want%d", u.Field34, value)
				}
			}
			if calls != 2 {
				t.Fatalf("handler calls%d want2", calls)
			}
		})
	}
	updateRegistryReportCounts(t)
}
