//go:build porttest

package opennox

import (
	"fmt"
	"math"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

var collisionRegistryProjectileHashes = map[string]string{
	"collision-registry-projectile-arrow-kinds":     "930eb0b839a789b8be1bef431890b29c6766ac065be3fe2dc98a70c56a3091b1",
	"collision-registry-projectile-arrow-modifiers": "d0ca525209c43f16ea2092c0cec754b337e01125581cad2121cec22a5811346d",
	"collision-registry-projectile-basic-00":        "fd58bfa5381475f5c24bdbc1a0f780c2498338d69adec6bd80cac6311e57fee5",
	"collision-registry-projectile-basic-01":        "486a9e2c2226aeaa5643e8853ae0073c99b8172960221ab19b6dd9189389c985",
	"collision-registry-projectile-basic-02":        "c45cf4736175bc6fe5601e926af18ea5a272b3ecc6c45f8102a9f37af7a75275",
	"collision-registry-projectile-basic-03":        "f36dceed6dfdf8f84ded7cc1947a0201403771989c91e750fe69dc60710a6e83",
	"collision-registry-projectile-basic-06":        "82f6ee8efa3f820600299f3e6c7efe26ebfe20550d95b966be6e8f10c8df18e8",
	"collision-registry-projectile-basic-09":        "23cc0165b0b0e8ae9fa8817260dc3e484cb8a23c1cfb395adbddd1c943345ded",
	"collision-registry-projectile-basic-10":        "84d2c2c96edddfaa95f1369d02486f5a82ba1857db5554420e3cf83327d261a1",
	"collision-registry-projectile-basic-11":        "46994269c469cfb324e4307de838057ad6d18c333612f3e8600b36d0190e51f2",
	"collision-registry-projectile-basic-12":        "2a547eb505b4750dcaf50219fade4977eb957c093ab043ca8bede74050e7c2c7",
	"collision-registry-projectile-basic-13":        "e03e592fb7122158ed63176a99425595c272273d165a099ff950b6f1d96ddd70",
	"collision-registry-projectile-basic-14":        "b2f910b9bde78ffede9dae5ff52f4bc3a660ae7338b579a79c0dec03571064b9",
	"collision-registry-projectile-basic-15":        "23cc0165b0b0e8ae9fa8817260dc3e484cb8a23c1cfb395adbddd1c943345ded",
	"collision-registry-projectile-basic-16":        "f48c8516c7ce9926b5c372ad150b239026382d8a998c3c21b462a0f9389b67aa",
	"collision-registry-projectile-basic-17":        "d64c6561b28653ce58aeaed6ca14321db31fc399cecf69ebcfaa8782c6f9100a",
	"collision-registry-projectile-basic-18":        "fde322a90af59da1b9e9edea6a2095f48e555142695c0b7fcc2d4c7f85d653ed",
	"collision-registry-projectile-basic-23":        "60fc9f8741637c0b5363a3c79bff2913ce82178ca8b10b5abd17c65fed59d8f4",
	"collision-registry-projectile-basic-24":        "d55820ca0f862d0c2522454c5c75f02936498c5eb18f3d2eb87c20fcefc53fcc",
	"collision-registry-projectile-basic-25":        "e77e437d295f754ddec930060b39e76fdd3235177e88118e685c4c5108de4fd9",
	"collision-registry-projectile-basic-26":        "66200f608be755d71eecb3e3c2542a84861d286f93145c9c035b7fdb808b742a",
	"collision-registry-projectile-bomb":            "0983c79b62342f546554bd4492df95e4f0759221f0723e939daba9d92b7d25e8",
	"collision-registry-projectile-buffs-05":        "4d2bd05dcf4f57bda20e5aaabb6d96b7fa6f5546db5151f014fd43e598449dd1",
	"collision-registry-projectile-buffs-08":        "ec12a673a654ff81ae5cad3b8a09a912a4856f25da23d743265a7756a2f41ea1",
	"collision-registry-projectile-buffs-12":        "eda20dfe39bedf76677d843113c40988beb53c483b67a0f20d82376d72494386",
	"collision-registry-projectile-buffs-16":        "e08ad412edb54e30f1049877727210139b2913ba82793857155b9aa484f125aa",
	"collision-registry-projectile-buffs-18":        "829259713cfff516461d24fd7bccf1b82f4a8881a62161ba75f246732511623b",
	"collision-registry-projectile-buffs-25":        "ea50faa02c3c2921816768bc564c905d68648714fe9abbc3721214cfdcc4bf24",
	"collision-registry-projectile-buffs-26":        "a39440a7e942adbac7ce34ca54c557d8978d02932cc25c56e605ccb87ee0669e",
	"collision-registry-projectile-chakram":         "e8a0e4177e0c6a127b46e5edb7162f1bec6f47cdd9e0afba4b875dfaace16465",
	"collision-registry-projectile-chakram-owner":   "c45f0ff1c07832d116c1747eef59ca4e60b90e14370efa6b374b972a57629b12",
	"collision-registry-projectile-damage-clock":    "c406efa01ca507029d1344d8db31022fb42a1d74a6d17d8b86391a54bac7a289",
	"collision-registry-projectile-death":           "ecaa58153a297c10373c5b2af0b69a8720b8b18a055011130bee79561ffaa7c8",
	"collision-registry-projectile-mana-clock":      "513b5d98f58604edb6be71b4904152769bc5b1df23fdef10ba313ea147176812",
	"collision-registry-projectile-named-damage":    "2b8ff32d5e1f00f0e2434938130c426c0631a5afc10b2f1ce85627e906d5b6f0",
	"collision-registry-projectile-normal-01":       "3d05a9f2edd2bb62722b9269b7bca2639c4cd94f2f05ccdcf56fd12d3f17ea7b",
	"collision-registry-projectile-normal-05":       "7175ea0dba63fc9265bca0af38621439e33a0b49d17fcaa94161a1e8cde1eeeb",
	"collision-registry-projectile-normal-09":       "2f5b7ad4649204f92f044fe76c3cec7ced9075b4e5a7d6baeacc73092d887977",
	"collision-registry-projectile-normal-10":       "c028f3711eb3d679ab8f377b13f1b52d15e70f612e243c57296c680e6e4e4eff",
	"collision-registry-projectile-normal-11":       "3bafcb9ccf10d9081ba2e190e27baa5de0c3b315ad11e871f98aa955cc80ed16",
	"collision-registry-projectile-normal-12":       "b164b5ab7cee831d2b023daade6ffcb92211f69d9d0362cfa25b4777a37dee81",
	"collision-registry-projectile-normal-13":       "2f5b7ad4649204f92f044fe76c3cec7ced9075b4e5a7d6baeacc73092d887977",
	"collision-registry-projectile-normal-15":       "2f5b7ad4649204f92f044fe76c3cec7ced9075b4e5a7d6baeacc73092d887977",
	"collision-registry-projectile-positive":        "f976096fb1a974940a6f7eb840b1ad9392e396360847cfea1afad61006b10c61",
	"collision-registry-projectile-splash":          "80722f4708ecf4ac24e675237c5fea1b77455d68c2374820604061b0088a0c4a",
}

func projectileCollisionRegistryBase() legacy.PortTestRoamSpec {
	s := projectileCollisionBase()
	s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Collision.RegisteredCollision = true
	return s
}

func projectileCollisionRegistryHash(t *testing.T, name string, s []legacy.PortTestRoamSpec) {
	t.Helper()
	collisionRegistryHash(t, name, effectsTimedRun(t, s), collisionRegistryProjectileHashes[name])
}

// These sibling roots use the real registered callback route. The original roots remain the frozen direct-C baseline.
func TestCollisionRegistryProjectileBasic(t *testing.T) {
	for _, op := range []int{0, 1, 2, 3, 6, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 23, 24, 25, 26} {
		t.Run(fmt.Sprint(op), func(t *testing.T) {
			var cases []legacy.PortTestRoamSpec
			for _, target := range []int{0, 4, 100, 101} {
				for _, accepted := range []int32{0, 1, 256, -1} {
					s := projectileCollisionRegistryBase()
					p := s.Callbacks.Shop
					w := p.TemporaryUpdates
					s.Callbacks.DamageResult = &accepted
					w.Target = target
					if op == 18 {
						w.World.CollideWords[0] = map[int]uint32{0: math.Float32bits(540), 4: math.Float32bits(520)}
						w.World.Objectives.Attack.Collision.CollideRefs = nil
					}
					p.Sequence = []legacy.PortTestShopAction{{Op: 1000 + op}}
					cases = append(cases, s)
				}
			}
			projectileCollisionRegistryHash(t, fmt.Sprintf("collision-registry-projectile-basic-%02d", op), cases)
		})
	}
}

func TestCollisionRegistryProjectileNormals(t *testing.T) {
	for _, op := range []int{1, 5, 9, 10, 11, 12, 13, 15} {
		t.Run(fmt.Sprint(op), func(t *testing.T) {
			var cases []legacy.PortTestRoamSpec
			for _, normal := range [][2]float32{{1, 0}, {0, 1}, {1, 1}, {1, -1}, {0, 0}, {0.6, 0.8}} {
				for _, wall := range []bool{false, true} {
					s := projectileCollisionRegistryBase()
					p := s.Callbacks.Shop
					w := p.TemporaryUpdates
					a := w.World.Objectives.Attack
					w.Target = 0
					a.Collision.Normal = &[2]uint32{math.Float32bits(normal[0]), math.Float32bits(normal[1])}
					a.Collision.WallContact = wall
					a.Collision.WallXY = [2]int32{6, 4}
					p.Sequence = []legacy.PortTestShopAction{{Op: 1000 + op}}
					cases = append(cases, s)
				}
			}
			projectileCollisionRegistryHash(t, fmt.Sprintf("collision-registry-projectile-normal-%02d", op), cases)
		})
	}
}

func TestCollisionRegistryProjectileBuffs(t *testing.T) {
	for _, op := range []int{5, 8, 12, 16, 18, 25, 26} {
		t.Run(fmt.Sprint(op), func(t *testing.T) {
			var cases []legacy.PortTestRoamSpec
			for _, target := range []int{100, 101} {
				for _, buff := range []uint32{0, 1 << 23, 1 << 27, 1 << 14, 1} {
					for _, dir := range []uint32{0, 128} {
						for _, flags := range []uint32{1, 2049, 4097} {
							s := projectileCollisionRegistryBase()
							p := s.Callbacks.Shop
							w := p.TemporaryUpdates
							o := w.World.Objectives
							s.Lifecycle.GameFlags = flags
							w.Target = target
							o.PlayerWords[1][56] = math.Float32bits(540)
							o.PlayerWords[1][60] = math.Float32bits(512)
							o.PlayerWords[target-100][340] = buff
							o.PlayerWords[target-100][124] = dir
							if op == 18 {
								w.World.CollideWords[0] = map[int]uint32{0: math.Float32bits(540), 4: math.Float32bits(520)}
								o.Attack.Collision.CollideRefs = nil
							}
							p.Sequence = []legacy.PortTestShopAction{{Op: 1000 + op}}
							cases = append(cases, s)
						}
					}
				}
			}
			projectileCollisionRegistryHash(t, fmt.Sprintf("collision-registry-projectile-buffs-%02d", op), cases)
		})
	}
}

func TestCollisionRegistryProjectileChakram(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, target := range []int{0, 4, 100, 101} {
		for _, count := range []uint32{0, 1, 2, 4} {
			for _, mode := range []uint32{0, 1, 2} {
				for _, held := range []bool{false, true} {
					s := projectileCollisionRegistryBase()
					p := s.Callbacks.Shop
					w := p.TemporaryUpdates
					a := w.World.Objectives.Attack
					w.Target = target
					w.UpdateWords[0][4] = count
					w.UpdateWords[0][24] = mode
					a.Collision.Normal = &[2]uint32{math.Float32bits(1), 0}
					if held {
						p.Items[2].Type = 23
						a.ActorRefs[504] = 5
						w.ItemRefs[2] = map[int]int{492: 3}
					}
					w.ItemWords[2][16] = 4
					p.Sequence = []legacy.PortTestShopAction{{Op: 1019}}
					cases = append(cases, s)
				}
			}
		}
	}
	projectileCollisionRegistryHash(t, "collision-registry-projectile-chakram", cases)
}

func TestCollisionRegistryProjectileBomb(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, target := range []int{0, 4, 100, 101} {
		for _, flags := range []uint32{1, 2049} {
			for _, dead := range []uint32{0, 0x20, 0x8000} {
				s := projectileCollisionRegistryBase()
				p := s.Callbacks.Shop
				w := p.TemporaryUpdates
				a := w.World.Objectives.Attack
				a.Actor = 1
				p.Resources.Subject = 3
				w.Target = target
				s.Lifecycle.GameFlags = flags
				a.ActorRefs = nil
				a.UpdateWords = map[int]uint32{1272: 0, 1276: 0}
				p.Items[1].Flags = 4 | dead
				w.World.Objectives.PlayerWords[1][16] = 4 | dead
				p.Sequence = []legacy.PortTestShopAction{{Op: 1004}}
				cases = append(cases, s)
			}
		}
	}
	projectileCollisionRegistryHash(t, "collision-registry-projectile-bomb", cases)
}

func TestCollisionRegistryProjectileArrowModifiers(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, mask := range []uint32{0, 1, 4, 8, 15} {
		for _, pre := range []uint32{0, 4, 8, 12} {
			for _, damage := range []float32{0, 0.49, 0.5, 1.5, 25.75} {
				s := projectileCollisionRegistryBase()
				p := s.Callbacks.Shop
				a := p.TemporaryUpdates.World.Objectives.Attack
				a.AttackMask = mask
				a.PreMask = pre
				a.SetOutput = true
				a.Output = math.Float32bits(damage)
				p.Items[0].Mods = [4]bool{true, true, true, true}
				p.Sequence = []legacy.PortTestShopAction{{Op: 1023}}
				cases = append(cases, s)
			}
		}
	}
	projectileCollisionRegistryHash(t, "collision-registry-projectile-arrow-modifiers", cases)
}

func TestCollisionRegistryProjectilePositive(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for i, op := range []int{0, 0, 13, 16, 25, 26, 3, 19} {
		s := projectileCollisionRegistryBase()
		p := s.Callbacks.Shop
		w := p.TemporaryUpdates
		a := w.World.Objectives.Attack
		switch i {
		case 1:
			v := int32(0)
			s.Callbacks.DamageResult = &v
		case 2:
			w.Target = 0
			a.Collision.Normal = &[2]uint32{math.Float32bits(1), math.Float32bits(1)}
		case 3, 4, 5:
			w.Target = 101
		case 6:
			w.Target = 101
			w.World.Objectives.PlayerUpdateWords[1] = map[int]uint32{4: 20 | 20<<16, 8: 100}
		case 7:
			p.Items[2].Type = 23
			w.Target = 100
			a.ActorRefs[504] = 5
			w.ItemRefs[2] = map[int]int{492: 3}
			w.ItemWords[2][16] = 4
		}
		p.Sequence = []legacy.PortTestShopAction{{Op: 1000 + op}}
		cases = append(cases, s)
	}
	r := effectsTimedRun(t, cases)
	deleted := func(i int) bool {
		trace := r[i].Trace
		for j := 0; j+1 < len(trace); j++ {
			if trace[j] == 32 && trace[j+1] == 70000 {
				return true
			}
		}
		return false
	}
	if len(r[0].Callbacks.Damage) != 5 || r[0].Callbacks.Damage[3] != 12 || r[0].Callbacks.Damage[4] != 11 || !deleted(0) {
		t.Fatal("accepted projectile hit must deal 12 damage and request deletion")
	}
	if len(r[1].Callbacks.Damage) != 5 || deleted(1) {
		t.Fatal("rejected projectile hit must preserve projectile")
	}
	u := worldObject(t, r[2], 0, 70000)
	if u[20] != math.Float32bits(-4) || u[21] != math.Float32bits(-3) {
		t.Fatalf("bounce velocity: %x %x", u[20], u[21])
	}
	unit, _, _ := objectivePlayerData(r[3].Callbacks.Shop.Sequence[0], 1)
	if unit[85]&(1<<4) == 0 {
		t.Fatal("webbing slow missing")
	}
	unit, _, _ = objectivePlayerData(r[4].Callbacks.Shop.Sequence[0], 1)
	if unit[85]&(1<<5|1<<14) != (1<<5|1<<14) || len(r[4].Lifecycle.Created) != 1 {
		t.Fatal("bear trap must spawn and stun/anchor target")
	}
	if len(r[5].Lifecycle.Created) != 1 {
		t.Fatal("gas trap must create cloud")
	}
	_, ud, _ := objectivePlayerData(r[6].Callbacks.Shop.Sequence[0], 1)
	if ud[1]&0xffff != 8 {
		t.Fatalf("mana after drain: %d", ud[1]&0xffff)
	}
	unit, ud, _ = objectivePlayerData(r[7].Callbacks.Shop.Sequence[0], 0)
	if unit[126] != 70002 || ud[26] != 70002 || !deleted(7) {
		t.Fatal("returning chakram must restore and equip its weapon")
	}
	collisionRegistryHash(t, "collision-registry-projectile-positive", r, collisionRegistryProjectileHashes["collision-registry-projectile-positive"])
}

func TestCollisionRegistryProjectileDamageClock(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, frame := range []uint32{0, 1, 2, 0xffffffff} {
		for _, damage := range []uint32{0, 1, 2, 3, 255} {
			for _, health := range []bool{false, true} {
				for _, kind := range []uint32{0, 2, 11} {
					s := projectileCollisionRegistryBase()
					p := s.Callbacks.Shop
					s.Owner.Frame = frame
					p.Items[1].Health = health
					p.TemporaryUpdates.World.CollideWords[0] = map[int]uint32{0: damage, 4: kind}
					p.TemporaryUpdates.World.Objectives.Attack.Collision.CollideRefs = nil
					p.Sequence = []legacy.PortTestShopAction{{Op: 1002}}
					cases = append(cases, s)
				}
			}
		}
	}
	projectileCollisionRegistryHash(t, "collision-registry-projectile-damage-clock", cases)
}

func TestCollisionRegistryProjectileManaClock(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, frame := range []uint32{0, 1, 15, 16, 32767, 32768, 65535, 65536, 0xffffffff} {
		for _, stamp := range []uint32{0, 32767, 32768, 65535} {
			for _, mana := range []uint32{0, 1, 20} {
				s := projectileCollisionRegistryBase()
				p := s.Callbacks.Shop
				w := p.TemporaryUpdates
				a := w.World.Objectives.Attack
				s.Owner.Frame = frame
				s.Owner.FPS = 30
				w.Target = 101
				a.ActorWords = map[int]uint32{540: stamp << 16}
				w.World.Objectives.PlayerUpdateWords[1] = map[int]uint32{4: mana | mana<<16, 8: 100}
				p.Sequence = []legacy.PortTestShopAction{{Op: 1003}, {Op: 1003}}
				cases = append(cases, s)
			}
		}
	}
	projectileCollisionRegistryHash(t, "collision-registry-projectile-mana-clock", cases)
}

func TestCollisionRegistryProjectileNamedDamage(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, name := range []string{"", "ThrowingStone", "ImpShot"} {
		for _, damage := range []float64{-1, 0, 0.49, 0.5, 1.9, 100, 65536} {
			for _, accepted := range []int32{0, 1, 256, -1} {
				s := projectileCollisionRegistryBase()
				p := s.Callbacks.Shop
				s.Callbacks.DamageResult = &accepted
				if name != "" {
					p.TemporaryUpdates.World.ItemNames = []string{name}
				}
				p.EffectsUse.Balance["UrchinStoneDamage"] = []float64{damage}
				p.EffectsUse.Balance["ImpShotDamage"] = []float64{damage}
				p.Sequence = []legacy.PortTestShopAction{{Op: 1000}}
				cases = append(cases, s)
			}
		}
	}
	projectileCollisionRegistryHash(t, "collision-registry-projectile-named-damage", cases)
}

func TestCollisionRegistryProjectileDeath(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, target := range []int{0, 4, 100, 101} {
		for _, callback := range []bool{false, true} {
			for _, flags := range []uint32{0, 0x20, 0x8000} {
				s := projectileCollisionRegistryBase()
				p := s.Callbacks.Shop
				w := p.TemporaryUpdates
				w.Target = target
				w.DieCallback = callback
				p.Items[0].Flags = flags
				p.Sequence = []legacy.PortTestShopAction{{Op: 1006}}
				cases = append(cases, s)
			}
		}
	}
	projectileCollisionRegistryHash(t, "collision-registry-projectile-death", cases)
}

func TestCollisionRegistryProjectileArrowKinds(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, name := range []string{"", "ArcherArrow", "ArcherBolt", "WeakArcherArrow"} {
		for _, flags := range []uint32{1, 2049, 4097} {
			for _, accepted := range []int32{0, 1} {
				for _, hp := range []uint16{0, 1} {
					for _, max := range []uint16{0, 100} {
						s := projectileCollisionRegistryBase()
						p := s.Callbacks.Shop
						w := p.TemporaryUpdates
						a := w.World.Objectives.Attack
						s.Lifecycle.GameFlags = flags
						s.Callbacks.DamageResult = &accepted
						p.Items[1].HP = hp
						p.Items[1].MaxHP = max
						if name != "" {
							w.World.ItemNames = []string{name}
							a.Collision.DefinitionRefs = map[int]int{23: 3}
						}
						p.EffectsUse.Balance["BoltSoloDamageMin"] = []float64{12}
						p.Sequence = []legacy.PortTestShopAction{{Op: 1023}}
						cases = append(cases, s)
					}
				}
			}
		}
	}
	projectileCollisionRegistryHash(t, "collision-registry-projectile-arrow-kinds", cases)
}

func TestCollisionRegistryProjectileSplash(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{5, 8} {
		for _, target := range []int{0, 4} {
			for _, radius := range []float64{0, 5, 30, 60} {
				for _, indexed := range [][]int{nil, {4}, {4, 5}} {
					s := projectileCollisionRegistryBase()
					p := s.Callbacks.Shop
					w := p.TemporaryUpdates
					w.Target = target
					w.Indexed = indexed
					p.Items[2].Health = true
					p.Items[2].HP = 100
					p.Items[2].MaxHP = 100
					w.ItemWords[2][56] = math.Float32bits(550)
					w.ItemWords[2][60] = math.Float32bits(512)
					p.EffectsUse.Balance["MagicMissileRange"] = []float64{radius}
					p.EffectsUse.Balance["MagicMissilePushRange"] = []float64{radius}
					w.World.CollideWords[0] = map[int]uint32{0: uint32(radius)}
					p.Sequence = []legacy.PortTestShopAction{{Op: 1000 + op}}
					cases = append(cases, s)
				}
			}
		}
	}
	projectileCollisionRegistryHash(t, "collision-registry-projectile-splash", cases)
}

func TestCollisionRegistryProjectileChakramOwner(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, owner := range []int{0, 100} {
		for _, flags := range []uint32{4, 0x24, 0x8004} {
			for _, weaponFlags := range []uint32{4, 0x24} {
				s := projectileCollisionRegistryBase()
				p := s.Callbacks.Shop
				w := p.TemporaryUpdates
				a := w.World.Objectives.Attack
				a.ActorRefs[508] = owner
				a.ActorRefs[504] = 5
				w.ItemRefs[2] = map[int]int{492: 3}
				p.Items[2].Type = 23
				w.ItemWords[2][16] = weaponFlags
				w.World.Objectives.PlayerWords[0][16] = flags
				p.Sequence = []legacy.PortTestShopAction{{Op: 1019}}
				cases = append(cases, s)
			}
		}
	}
	projectileCollisionRegistryHash(t, "collision-registry-projectile-chakram-owner", cases)
}
