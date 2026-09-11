//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

var objectStateHashes = map[string]string{
	"object-state-attributes":      "7036d632c9a7d4285cc8308aaa45f37c3da5820eeb6db19952a85e02616aefae",
	"object-state-bits-01":         "6072068d2d8755919b6a02c5b76273840fa57cf36d1f348e36bdceedc389c86c",
	"object-state-bits-02":         "6ff784cdd27873fc2223169dfc2663c32da15f131c0bdf4cb544f1b889f22b77",
	"object-state-bits-04":         "ab0e16527da1e48882947921907cf5b9125f49d372a37d6747ca58d6d7b70eec",
	"object-state-bits-05":         "1d8b66e4b0c96a9a71728dbf024a3724ca43f75e74276703a89a926919f974e2",
	"object-state-checksum-inputs": "6e45d2a33cc0e305cd2ca6cc3e6bb8e08a9f4c6b122dd82b17e0c51492831df3",
	"object-state-cleanup":         "02f1c914f53dab0537f9006aa32d7acd85a485f11c13dcf55b208ec4a6de1cea",
	"object-state-collide-40":      "4f409adaf8e7b1620ec98b6dd01bdb56a0e0c4d2f64441190251b7674d78dda9",
	"object-state-collide-41":      "48593cb7d5f15c89ae42dc6e3a92690459846489cd53e017490651ebce79568f",
	"object-state-collide-42":      "6e8439f3e88273d4313bb48462240ee6586a5c97b8f7b6d165e1b7839804f1b1",
	"object-state-direction-edges": "b2a17f6a0dd5124c324be202bd292d4c911909f29580587b6193a2aaf4bda271",
	"object-state-distance-edges":  "06cc9c83bb8b2d41c69b9ec584be12bccbd0e9a4f5c1aac2eb032bdc5103525a",
	"object-state-doors-38":        "eb129387fcdbd55faf8220fe94f629f667328d6d2bf479f9c969e320b2350598",
	"object-state-doors-39":        "4bfc423de483438598bd46fe5e5a35ccc694d1a736f72e84c7f762d36c7bb89a",
	"object-state-freeze-25":       "8062d5dce27e92a6e05c26d923d96ad42a09c074649e4ab8adc07babbb318782",
	"object-state-freeze-26":       "3cf1726afd29b73634364e0c57b58d5659e7cf4f778370841b6af83a8d856f01",
	"object-state-freeze-27":       "693a8ded69fa9abc8e54f118d4409b52fc535834d122ebd2142fd3d2f56640e5",
	"object-state-geometry-13":     "948af1fc29d2c5bb6d91534901865726b27dab3db7598e1c36c60e8f04c36d1a",
	"object-state-geometry-14":     "6c421b4101d5f76dc8dcad3cf69684aaf6cfd123a56f31fcd3abe3928da3a355",
	"object-state-geometry-15":     "f0d673d52514bcbbb6bd548315f36c38a238ea3a1fed646b3111a876d80b35e3",
	"object-state-loot":            "d1e9448ef68eb541b56185d9933179d3f24ef97ebdcffbf6a8e4f55eaca6e0ca",
	"object-state-nil":             "3324b684f699abf5f261ae909319cc2bed45ad34ec75df014fe92dee72d7936e",
	"object-state-owned-08":        "223b8637ab8a98bf65f5615dc70e591885132b7d2a8afe573db074f2d8273fd7",
	"object-state-owned-09":        "223b8637ab8a98bf65f5615dc70e591885132b7d2a8afe573db074f2d8273fd7",
	"object-state-owned-10":        "1f2023758e88caf9bb3cadbfb7a87f11d8f7eadb3ca495cdf666af1f6da88633",
	"object-state-owned-23":        "4a7191d21c67505d865caed600a6dae8ad84f6ae30b771db229f28093cdbc8ad",
	"object-state-owned-24":        "f2399804a1ff8e30c9f9532d41aa6a8bee0b9afb949ad41045fba572c5e0502f",
	"object-state-owned-30":        "3faac36cfa92af5737f51719ec69951ef8b3df2eea42fc1c76f1942b499d36ab",
	"object-state-owned-31":        "f127538916b8018f0a5596d44edb6cf477991e1fd83599871d6ffc51c4c7a46a",
	"object-state-owned-32":        "dcf66149e722a122d978bc47b787eaa5299c5f0bdc101d5f1839f50d620f8db9",
	"object-state-owned-33":        "3eb44336eb834881d702ebc25308cf3b6bb853b2e71d9b35bd3462efefbdfbed",
	"object-state-owned-34":        "46def1439e53b6c9ee7ee04d33b1f1ed024dce43ec5a548859df64b41cc9d9d0",
	"object-state-owned-35":        "fc3309dec98e2d0cbb875345c562de2a596951719cc2abb37bc0a12760a36620",
	"object-state-owned-37":        "6e1390c2abc8e181a8f2c69f08f7ca1696ad738b224173370614ebc689f54611",
	"object-state-owned-43":        "223b8637ab8a98bf65f5615dc70e591885132b7d2a8afe573db074f2d8273fd7",
	"object-state-players-19":      "93887205226c44812fd2a82bb3e131db2749b6c7a2f73308af5da9da3aa6caba",
	"object-state-players-28":      "5c0ceb2a173b6cfc5836417fd8934253df4606a274f2dc579ef13268f3c45cbb",
	"object-state-players-29":      "4364ec53b03f2c485220c0da7e0694efd8b609e0a3101d59f1ac38bd2d72d839",
	"object-state-players-36":      "996f27b53a52fe7be817c329907d2ef97635be4f80a8cd3ff6af476339acb57c",
	"object-state-positive":        "ce98ee2da86dded7b06616a5bbb2af2f4bd8c99198c402d3fd02942e7e84576b",
	"object-state-recent":          "c36eca0c0f26667388b6ea6ba64351c597626a700c71dd37e1d203b72985e5ea",
	"object-state-sync-00":         "ebe7fbfe04d48f4c62a28c44cc51045601013acd858e9f53caea08f32f9156a4",
	"object-state-sync-01":         "eec321ce28d7160dd20f9a4919107fb40f3a2336552dd6d8e8f4b2c548b834d2",
	"object-state-sync-02":         "2f99aaa6819c2b7f23b081417ee2f57781d53a33395a653dc0bed020d81549ff",
	"object-state-sync-03":         "ef572fef0d064d11cc59c327c8195058ec725486d45b078c8f704cd8c42206a7",
	"object-state-sync-04":         "732ca80ed5605d0dd26fa05550f2fceb27eb9d54d7436631ee43f7dc03d87de7",
	"object-state-sync-05":         "d41bef4f3fc46a13833807719581c3a398f7441f3f1fad341103af42791ca5fd",
	"object-state-sync-07":         "6d11a48fed3fb93af6de06ed2675c373ea55beb985ae2b708e1e29d40f21131b",
	"object-state-sync-17":         "f77743e159acda53f719792db427d4dde4224ed49fd33a74330ba9105d037eaf",
	"object-state-sync-20":         "bc2c423432f01d922d2556aa0b930c3fd7939943dcaa8638696da58957c9bce2",
	"object-state-sync-21":         "76340aac8c201e21a33904da9b1f0083b5c15f351658888279822b87ee2fe4ad",
	"object-state-sync-22":         "32a4fc0a75883ee655fd7410d495e22cf50322f711ede835cba4e74337b70073",
	"object-state-teleport":        "32b8518d44aa235619e9e541ab0a3ea81f7c9089e9b524aa779313317f226e14",
}

func objectStateBase(op int) legacy.PortTestRoamSpec {
	s := attackBase()
	p := s.Callbacks.Shop
	a := p.TemporaryUpdates.World.Objectives.Attack
	a.Actor = 1
	a.State = &legacy.PortTestObjectStateSpec{Target: 4, X: 1, Y: 4, Z: 1, FloatBits: math.Float32bits(2.5)}
	a.RecordWords = map[int]uint32{0: math.Float32bits(540), 4: math.Float32bits(512)}
	a.RecordRefs = nil
	p.Inventory.Linked = nil
	p.Inventory.Owned = nil
	p.Equipment.ActiveWeapon = 0
	p.Equipment.WeaponFlags = 0
	p.Sequence = []legacy.PortTestShopAction{{Op: 1200 + op}}
	return s
}
func objectStateHash(t *testing.T, name string, s []legacy.PortTestRoamSpec) {
	t.Helper()
	callbackHash(t, name, effectsTimedRun(t, s), objectStateHashes[name])
}
func TestObjectStateSync(t *testing.T) {
	for _, op := range []int{0, 1, 2, 3, 4, 5, 7, 17, 20, 21, 22} {
		var cases []legacy.PortTestRoamSpec
		for _, subject := range []int{1, 2, 3} {
			for _, mask := range []uint32{0, 1, 0x80000000, 0xffffffff} {
				for _, value := range []int32{0, 1, -1} {
					s := objectStateBase(op)
					p := s.Callbacks.Shop
					a := p.TemporaryUpdates.World.Objectives.Attack
					p.Resources.Subject = subject
					p.Resources.SyncSeed = 0x12345000
					a.ActorWords = map[int]uint32{148: mask}
					a.State.X = value
					if op == 3 {
						a.State.FloatBits = math.Float32bits(float32(value))
					}
					p.Sequence = append(p.Sequence, p.Sequence[0])
					cases = append(cases, s)
				}
			}
		}
		objectStateHash(t, fmt.Sprintf("object-state-sync-%02d", op), cases)
	}
}
func TestObjectStateGeometry(t *testing.T) {
	for _, op := range []int{13, 14, 15} {
		var cases []legacy.PortTestRoamSpec
		for _, shape := range []uint32{1, 2, 3} {
			for _, pos := range [][2]float32{{512, 512}, {513, 512}, {512, 513}, {511, 512}, {512, 511}, {540, 540}, {490, 540}, {540, 490}, {490, 490}} {
				for _, dir := range []int32{0, 32, 64, 128, 192, 255} {
					s := objectStateBase(op)
					p := s.Callbacks.Shop
					a := p.TemporaryUpdates.World.Objectives.Attack
					a.State.X = dir
					a.ActorWords = map[int]uint32{172: shape, 176: math.Float32bits(10), 184: math.Float32bits(30), 188: math.Float32bits(20)}
					a.RecordWords = map[int]uint32{0: math.Float32bits(pos[0]), 4: math.Float32bits(pos[1])}
					p.TemporaryUpdates.ItemWords[1][56] = math.Float32bits(pos[0])
					p.TemporaryUpdates.ItemWords[1][60] = math.Float32bits(pos[1])
					p.TemporaryUpdates.ItemWords[1][172] = shape
					cases = append(cases, s)
				}
			}
		}
		objectStateHash(t, fmt.Sprintf("object-state-geometry-%02d", op), cases)
	}
}
func TestObjectStateAttributes(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, class := range []uint32{1, 4, 0x1000, 0x1000000, 0x2000000} {
		for _, mask := range []uint32{0, 1, 15} {
			for _, subclass := range []uint32{0, 0x10000} {
				s := objectStateBase(6)
				p := s.Callbacks.Shop
				a := p.TemporaryUpdates.World.Objectives.Attack
				a.Actor = 3
				p.Items[0].Class = class
				p.Items[0].Subclass = subclass
				a.RecordWords = map[int]uint32{16: 0x12345678}
				a.RecordRefs = map[int]int{}
				for i := 0; i < 4; i++ {
					if mask&(1<<i) != 0 {
						a.RecordRefs[4*i] = 4
					}
				}
				cases = append(cases, s)
			}
		}
	}
	objectStateHash(t, "object-state-attributes", cases)
}

func TestObjectStateOwnership(t *testing.T) {
	for _, op := range []int{8, 9, 10, 23, 24, 30, 31, 32, 33, 34, 35, 37, 43} {
		var cases []legacy.PortTestRoamSpec
		for _, subject := range []int{1, 2, 3} {
			for _, mode := range []int{0, 1, 2, 3} {
				s := objectStateBase(op)
				p := s.Callbacks.Shop
				a := p.TemporaryUpdates.World.Objectives.Attack
				p.Resources.Subject = subject
				if mode != 0 {
					p.Inventory.Linked = []int{0, 1}
					p.Inventory.Owned = []int{0, 1}
					p.Items[0].Class = 2
					p.Items[0].Subclass = 0x80
					p.Items[1].Subclass = 0x80
				}
				if mode == 2 {
					p.TemporaryUpdates.World.ItemNames = []string{"Crown"}
				}
				if mode == 3 {
					p.TemporaryUpdates.World.ItemNames = []string{"GameBall"}
				}
				a.State.X = 2
				a.State.Y = 0x80
				if op == 24 || op == 34 || op == 37 {
					a.Actor = 3
				}
				if op == 34 {
					a.State.Target = 4
					p.Items[0].Class = 1
					p.Items[1].Class = 1
					if mode&1 != 0 {
						p.Items[1].Type = p.Items[0].Type
					}
				}
				if op == 37 {
					p.TemporaryUpdates.World.ItemNames = []string{"Pixie"}
					a.UpdateWords = map[int]uint32{4: 123}
				}
				cases = append(cases, s)
			}
		}
		objectStateHash(t, fmt.Sprintf("object-state-owned-%02d", op), cases)
	}
}
func TestObjectStateRecent(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, frame := range []uint32{0, 1, 2, 3, 0xffffffff} {
		for _, last := range []uint32{0, 1, 2, 0xfffffffe, 0xffffffff} {
			for _, hp := range []bool{false, true} {
				s := objectStateBase(12)
				s.Owner.Frame = frame
				p := s.Callbacks.Shop
				p.Resources.NoHealth = !hp
				p.TemporaryUpdates.World.Objectives.Attack.ActorWords = map[int]uint32{536: last}
				cases = append(cases, s)
			}
		}
	}
	objectStateHash(t, "object-state-recent", cases)
}
func TestObjectStateFreeze(t *testing.T) {
	for _, op := range []int{25, 26, 27} {
		var cases []legacy.PortTestRoamSpec
		for _, subject := range []int{1, 2, 3} {
			for _, flags := range []uint32{0, 2, 0x8000, 0x8002} {
				for _, force := range []int32{0, 1, 256, -1} {
					s := objectStateBase(op)
					p := s.Callbacks.Shop
					a := p.TemporaryUpdates.World.Objectives.Attack
					p.Resources.Subject = subject
					p.Resources.Flags = flags
					a.State.X = force
					a.State.Globals = map[int]uint32{1567712: 17}
					p.Sequence = append(p.Sequence, p.Sequence[0])
					cases = append(cases, s)
				}
			}
		}
		objectStateHash(t, fmt.Sprintf("object-state-freeze-%02d", op), cases)
	}
}
func TestObjectStatePlayers(t *testing.T) {
	for _, op := range []int{19, 28, 29, 36} {
		var cases []legacy.PortTestRoamSpec
		for _, target := range []int{0, 1, 3, 101} {
			for _, frame := range []uint32{0, 100, 0xffffffff} {
				s := objectStateBase(op)
				p := s.Callbacks.Shop
				a := p.TemporaryUpdates.World.Objectives.Attack
				a.State.Target = target
				a.State.X = 1
				s.Owner.Frame = frame
				p.TemporaryUpdates.World.Objectives.ObjectList = []int{1, 3, 4}
				cases = append(cases, s)
			}
		}
		objectStateHash(t, fmt.Sprintf("object-state-players-%02d", op), cases)
	}
}

func TestObjectStateCleanup(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, mode := range []int32{0, 1, 2} {
		for _, flags := range []uint32{1, 2049, 8193} {
			for _, owned := range []bool{false, true} {
				s := objectStateBase(11)
				p := s.Callbacks.Shop
				a := p.TemporaryUpdates.World.Objectives.Attack
				a.State.X = mode
				s.Lifecycle.GameFlags = flags
				p.TemporaryUpdates.World.Objectives.ObjectList = []int{1, 4}
				a.State.MissileList = []int{3, 5}
				p.Items[0].Class = 1
				p.TemporaryUpdates.World.ItemNames = []string{"Pixie"}
				p.Items[1].Class = 2
				p.Items[1].Subclass = 0x100
				if owned {
					a.ActorRefs = map[int]int{516: 3}
					p.TemporaryUpdates.ItemRefs[0] = map[int]int{508: 1}
				}
				cases = append(cases, s)
			}
		}
	}
	objectStateHash(t, "object-state-cleanup", cases)
}
func TestObjectStateTeleport(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, subject := range []int{1, 2, 3} {
		for _, flags := range []uint32{1, 2049, 4097} {
			for _, buff := range []uint32{0, 1 << 14} {
				for _, disabled := range []uint32{0, 2} {
					s := objectStateBase(16)
					p := s.Callbacks.Shop
					p.Resources.Subject = subject
					p.Resources.Subclass = 8
					p.Resources.Flags = disabled
					p.Resources.Buffs = buff
					s.Lifecycle.GameFlags = flags
					cases = append(cases, s)
				}
			}
		}
	}
	objectStateHash(t, "object-state-teleport", cases)
}
func TestObjectStateLoot(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, name := range []string{"BarrelPortTest", "CratePortTest"} {
		for seed := 1; seed <= 100; seed++ {
			s := objectStateBase(18)
			s.Seed = seed
			a := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack
			a.Actor = 3
			a.State.ActorName = name
			cases = append(cases, s)
		}
	}
	objectStateHash(t, "object-state-loot", cases)
}
func TestObjectStateDoors(t *testing.T) {
	for _, op := range []int{38, 39} {
		var cases []legacy.PortTestRoamSpec
		for _, class := range []uint32{1, 0x80} {
			for _, match := range []bool{false, true} {
				for _, flags := range []uint32{1, 4097} {
					s := objectStateBase(op)
					p := s.Callbacks.Shop
					a := p.TemporaryUpdates.World.Objectives.Attack
					a.Actor = 3
					p.Items[0].Class = class
					s.Lifecycle.GameFlags = flags
					a.RecordWords = map[int]uint32{0: 6, 4: 4}
					a.UpdateWords = map[int]uint32{0: 0x100, 16: 6, 20: 4}
					if !match {
						a.UpdateWords[20] = 5
					}
					cases = append(cases, s)
				}
			}
		}
		objectStateHash(t, fmt.Sprintf("object-state-doors-%02d", op), cases)
	}
}
func TestObjectStateCollide(t *testing.T) {
	for _, op := range []int{40, 41, 42} {
		var cases []legacy.PortTestRoamSpec
		for _, target := range []int{0, 4, 101} {
			for _, enabled := range []bool{false, true} {
				for _, frame := range []uint32{0, 100, 0xffffffff} {
					s := objectStateBase(op)
					p := s.Callbacks.Shop
					a := p.TemporaryUpdates.World.Objectives.Attack
					a.State.Target = target
					s.Owner.Frame = frame
					if op != 42 {
						p.Resources.Subject = 3
						a.UpdateWords = map[int]uint32{1272: 0xffffffff}
						if enabled {
							a.UpdateWords[1272] = 17
						} else if op == 41 {
							a.ActorRefs = map[int]int{508: 101}
						}
					} else if enabled {
						p.Equipment.ActiveAbilities = 2
					}
					p.Items[1].Health = true
					p.Items[1].HP = 100
					p.Items[1].MaxHP = 100
					p.Items[1].Flags = 0
					p.Items[1].Class = 0x400000
					a.State.HealthRefs = []int{101}
					p.EffectsUse.Balance["BerserkerDamage"] = []float64{12}
					p.EffectsUse.Balance["BerserkerStunDuration"] = []float64{30}
					p.EffectsUse.Balance["BerserkerPainRatio"] = []float64{0.25}
					cases = append(cases, s)
				}
			}
		}
		objectStateHash(t, fmt.Sprintf("object-state-collide-%02d", op), cases)
	}
}

func TestObjectStateIndividualBits(t *testing.T) {
	for _, op := range []int{1, 2, 4, 5} {
		var cases []legacy.PortTestRoamSpec
		for _, subject := range []int{1, 2, 3} {
			for bit := 0; bit < 32; bit++ {
				for _, set := range []int32{0, 1} {
					s := objectStateBase(op)
					p := s.Callbacks.Shop
					a := p.TemporaryUpdates.World.Objectives.Attack
					p.Resources.Subject = subject
					p.Resources.SyncSeed = 0xa5123490
					a.ActorWords = map[int]uint32{148: 1 << bit}
					a.State.X = int32(uint32(1) << bit)
					a.State.Y = int32(uint32(1) << uint(bit%12))
					a.State.Z = set
					cases = append(cases, s)
				}
			}
		}
		objectStateHash(t, fmt.Sprintf("object-state-bits-%02d", op), cases)
	}
}
func TestObjectStateDirectionEdges(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, x := range []float32{-1024, -1, -0.001, 0, 0.001, 1, 1024} {
		for _, slope := range []float64{-2.4210529, -0.41304299, 0, 0.41304299, 2.4210529} {
			edge := float32(float64(x) * slope)
			for _, y := range []float32{math.Nextafter32(edge, float32(math.Inf(-1))), edge, math.Nextafter32(edge, float32(math.Inf(1)))} {
				s := objectStateBase(14)
				a := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack
				a.ActorWords = map[int]uint32{56: 0, 60: 0}
				a.RecordWords = map[int]uint32{0: math.Float32bits(x), 4: math.Float32bits(y)}
				cases = append(cases, s)
			}
		}
	}
	objectStateHash(t, "object-state-direction-edges", cases)
}
func TestObjectStateDistanceEdges(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, shape1 := range []uint32{0, 1, 2, 3} {
		for _, shape2 := range []uint32{0, 1, 2, 3} {
			for _, dx := range []float32{0, 0.009, 0.01, 0.011, 9.99, 10, 10.01, 30, 100000} {
				s := objectStateBase(13)
				p := s.Callbacks.Shop
				a := p.TemporaryUpdates.World.Objectives.Attack
				a.ActorWords = map[int]uint32{56: 0, 60: 0, 172: shape1, 176: math.Float32bits(5), 184: math.Float32bits(10), 188: math.Float32bits(20)}
				p.TemporaryUpdates.ItemWords[1] = map[int]uint32{56: math.Float32bits(dx), 60: 0, 172: shape2, 176: math.Float32bits(5), 184: math.Float32bits(20), 188: math.Float32bits(10)}
				cases = append(cases, s)
			}
		}
	}
	objectStateHash(t, "object-state-distance-edges", cases)
}
func TestObjectStatePositive(t *testing.T) {
	s := objectStateBase(2)
	s.Callbacks.Shop.Sequence = append(s.Callbacks.Shop.Sequence, legacy.PortTestShopAction{Op: 1204})
	r := effectsTimedRun(t, []legacy.PortTestRoamSpec{s})
	for j, st := range r[0].Callbacks.Shop.Sequence {
		u := st.ResourceData[1]
		if u[152/4] != 0xffffffff || u[16/4]&0x1000000 == 0 {
			t.Fatalf("step %d sync/on flags", j)
		}
		if j == 1 && u[132/4] != 1 {
			t.Fatal("animation frame")
		}
	}
	callbackHash(t, "object-state-positive", r, objectStateHashes["object-state-positive"])
}

func TestObjectStateNil(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{10, 19, 24, 28, 30, 33, 34, 37, 43} {
		s := objectStateBase(op)
		s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Actor = 0
		cases = append(cases, s)
	}
	objectStateHash(t, "object-state-nil", cases)
}

func TestObjectStateChecksumInputs(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, off := range []int{340, 248, 120, 128, 132, 136, 148, 152, 108, 104, 100, 96, 92, 88, 84, 80, 76, 72, 68, 64, 60, 16, 20, 36, 40, 44, 56, 4, 52, 124} {
		for _, bits := range []uint32{1, 0x80008000, 0xffffffff} {
			s := objectStateBase(22)
			s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.ActorWords = map[int]uint32{off: bits}
			cases = append(cases, s)
		}
	}
	for _, hp := range []uint16{0, 1, 65535} {
		for _, none := range []bool{false, true} {
			s := objectStateBase(22)
			p := s.Callbacks.Shop
			p.Resources.HP = hp
			p.Resources.OldHP = hp
			p.Resources.MaxHP = hp
			p.Resources.NoHealth = none
			cases = append(cases, s)
		}
	}
	objectStateHash(t, "object-state-checksum-inputs", cases)
}
