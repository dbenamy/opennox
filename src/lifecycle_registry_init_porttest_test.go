//go:build porttest

package opennox

import (
	"encoding/json"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

var lifecycleInitHashes = map[string]string{
	"lifecycle-registry-player-init":    "54e1c89cec57ae8ea1d664abd094585ada670cd3d853148026cd34a495b6cacc",
	"lifecycle-registry-reward-init-00": "71d4631b436690a7a594c397b7fe480641deb6db1e19ebe9ac82aeacd91d0b27",
	"lifecycle-registry-reward-init-01": "639d32f6f96572dfd227eb52ad0da52243d5c9679ddc70da29d8520fc825a922",
	"lifecycle-registry-reward-init-02": "b0fc5e3defd240404557b32025349145e4766ded7197e95064efb57d069b3574",
	"lifecycle-registry-reward-init-03": "ac3a0312e9819df56559ec11723cc036b8e0341db4a2f4d148c1a42dfff13a9a",
	"lifecycle-registry-reward-init-04": "c5217dba764ac7f4fb268d15ed519f6a75bb7f481b62d869dd093306a97cb3b1",
	"lifecycle-registry-reward-init-05": "ba0869eb1b20ee96ce9d9b1b94520c6d4ed4d55983ef4201afd4f8c3da1b76d3",
	"lifecycle-registry-reward-init-06": "8b54c7244308a65f3db9297808059ed427944cf6fc5919e7454f4fdc508d23f0",
	"lifecycle-registry-reward-init-07": "42921151496f20a2b1d28da1080f7d59350070b713902022d3efeb1f47daa591",
	"lifecycle-registry-reward-init-08": "96e7afa0bdbbbf6cb85b35e18dbc2685fd9df0f1c26e48ed6ecf5ff91a7d529c",
}

func lifecycleInitHash(t *testing.T, name string, rows []legacy.PortTestRoamResult) {
	t.Helper()
	callbackHash(t, name, rows, lifecycleInitHashes[name])
	counts := legacy.PortTestInitRegistryTakeCounts()
	if len(counts) == 0 {
		t.Fatal("no registered initialization executed")
	}
	b, e := json.Marshal(counts)
	if e != nil {
		t.Fatal(e)
	}
	t.Logf("LIFECYCLE_INIT_COUNTS %s", b)
}
func TestLifecycleRegistryRewardInit(t *testing.T) {
	for op := 0; op < 9; op++ {
		var cases []legacy.PortTestRoamSpec
		for i := 0; i < 32; i++ {
			s := rewardBase(op)
			s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Reward.RegisteredInit = true
			s.Seed = i + 1
			a := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack
			r := a.Reward
			r.InitWords[0] = uint32(i % 9)
			if op == 6 {
				r.InitWords[0] = 0
				if i%4 == 0 {
					r.InitWords[0] = 77
				}
			}
			a.ActorWords[20] = uint32(i % 16)
			a.ActorWords[12] = uint32(i % 16)
			r.UpdateWords[80] = uint32(i%5) << 24
			r.UpdateWords[84] = 0x03020100
			r.GeneratorStage = uint32(i % 4)
			r.UpdateWords[16] = 0 // empty type name for direction/name initializer.
			cases = append(cases, s)
		}
		lifecycleInitHash(t, fmt.Sprintf("lifecycle-registry-reward-init-%02d", op), effectsTimedRun(t, cases))
	}
}
func TestLifecycleRegistryPlayerInit(t *testing.T) {
	for _, op := range []int{20} {
		var cases []legacy.PortTestRoamSpec
		for class := byte(0); class < 3; class++ {
			for mode := 0; mode < 4; mode++ {
				s := controlsBase(op)
				s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Controls.RegisteredInit = true
				controlsStats(&s)
				p := s.Callbacks.Shop
				p.Resources.PlayerClass = class
				a := p.TemporaryUpdates.World.Objectives.Attack
				a.Controls.Equipment = true
				a.Controls.X = int32(mode & 1)
				a.Controls.Y = int32(mode >> 1)
				if op == 32 || op == 49 {
					a.Controls.Corpse = true
				}
				if op == 18 || op == 20 {
					if mode&1 != 0 {
						a.Controls.ByteReturn = 1
					} else if class < 2 && (op == 20 || mode>>1 == 0) {
						a.Controls.ByteReturn = 2
					}
				}
				if op == 49 {
					a.Controls.Bot = true
					a.Controls.BotWords = map[int]uint32{548: 0}
					if mode&1 != 0 {
						p.Resources.HP = 0
					}
				}
				p.TemporaryUpdates.World.Objectives.PlayerDataWords[0] = map[int]uint32{3684: 5, 4700: uint32(mode & 1)}
				// CallInit is void; original direct-helper tests retain byte-return assertions.
				a.Controls.ByteReturn = 0
				// The isolated setup supplies local respawn positioning and item callbacks.
				cases = append(cases, s)
			}
		}
		lifecycleInitHash(t, "lifecycle-registry-player-init", controlsRun(t, cases))
	}
}
