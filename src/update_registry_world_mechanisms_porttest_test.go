//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func TestUpdateRegistryWorldTriggers(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestWorld53B060, legacy.PortTestWorld53B1B0} {
		for _, state := range []uint32{0, 1, 3, 5, 255} {
			for _, bits := range []uint32{0, 1, 2, 3, 8, 9, 10, 11} {
				for _, on := range []bool{false, true} {
					for _, frame := range []uint32{99, 100, 101, 0xffffffff} {
						s := updateRegistryWorldBase()
						p := s.Callbacks.Shop
						w := p.TemporaryUpdates
						s.Owner.Frame = frame
						if !on {
							p.Items[0].Flags = 4
						}
						w.ItemWords[0][136] = 100
						w.UpdateWords[0] = map[int]uint32{0: bits, 8: state, 36: 237, 40: 239}
						w.UpdateRefs[0] = map[int]int{4: 1}
						p.Sequence = []legacy.PortTestShopAction{{Op: op}, {Op: op}}
						cases = append(cases, s)
					}
				}
			}
		}
	}
	updateRegistrySpecsHash(t, "update-registry-world-triggers", cases)
}

func TestUpdateRegistryWorldDoors(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, sub := range []uint32{0, 1, 4, 0x1000, 0x1005} {
		for _, material := range []uint32{0, 8} {
			for _, angles := range [][3]uint32{{0, 0, 0}, {0, 1, 0}, {0, 0, 2}, {0, 4, 2}, {31, 31, 0}, {1, 1, 31}} {
				for _, frame := range []uint32{105, 106, 0, 0xffffffff} {
					s := updateRegistryWorldBase()
					p := s.Callbacks.Shop
					w := p.TemporaryUpdates
					s.Owner.Frame = frame
					p.Items[0].Subclass = sub
					w.ItemWords[0][24] = material
					w.UpdateWords[0] = map[int]uint32{4: angles[0], 8: angles[1], 12: angles[2], 40: angles[2] * 8, 44: 90}
					w.Updatable = 3
					p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestWorld53AC50}, {Op: legacy.PortTestWorld53AC50}}
					cases = append(cases, s)
				}
			}
		}
	}
	updateRegistrySpecsHash(t, "update-registry-world-doors", cases)
}

func TestUpdateRegistryWorldSwitches(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestWorld53B300, legacy.PortTestWorld53B320} {
		for _, flags := range []uint32{0, 4, 0x40, 0x44, 0x1000000, 0x1000004, 0x1000044} {
			for _, sync := range []uint32{0, 1, 0xffffffff} {
				s := updateRegistryWorldBase()
				p := s.Callbacks.Shop
				p.Items[0].Flags = flags
				p.TemporaryUpdates.ItemWords[0][132] = sync
				p.Sequence = []legacy.PortTestShopAction{{Op: op}, {Op: op}}
				cases = append(cases, s)
			}
		}
	}
	updateRegistrySpecsHash(t, "update-registry-world-switches", cases)
}

func TestUpdateRegistryWorldElevators(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, state := range []uint32{0, 1, 2, 3, 4} {
		for _, height := range []uint32{0, 1, 2, 19, 20, 21, 30, 31, 32, 33, 62, 63, 64, 65, 0xffffffff} {
			for _, on := range []bool{false, true} {
				s := updateRegistryWorldBase()
				p := s.Callbacks.Shop
				w := p.TemporaryUpdates
				if !on {
					p.Items[0].Flags = 4
				}
				w.ItemWords[0][136] = 69
				w.UpdateWords[0] = map[int]uint32{12: state, 16: height}
				w.UpdateRefs[0] = map[int]int{4: 4}
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestWorld53B5D0}, {Op: legacy.PortTestWorld53B5D0}}
				cases = append(cases, s)
			}
		}
	}
	for _, state := range []uint32{0, 1, 2, 3, 4} {
		for _, prior := range []uint32{0, 1, 3} {
			for _, height := range []uint32{31, 32, 33, 64} {
				s := updateRegistryWorldBase()
				p := s.Callbacks.Shop
				w := p.TemporaryUpdates
				w.UpdateRefs[0] = map[int]int{4: 4}
				w.UpdateWords[0][12] = prior
				w.UpdateWords[1] = map[int]uint32{12: state, 16: height}
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestWorld53B380}}
				cases = append(cases, s)
			}
		}
	}
	updateRegistrySpecsHash(t, "update-registry-world-elevators", cases)
}

func TestUpdateRegistryWorldTeleports(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, state := range []uint32{0, 1, 2, 3} {
		for _, count := range []uint32{0, 1, 3, 4, 0xff} {
			for _, anim := range []uint32{0, 7, 8, 9} {
				for _, activate := range []uint32{0, 1} {
					s := updateRegistryWorldBase()
					p := s.Callbacks.Shop
					w := p.TemporaryUpdates
					w.UpdateWords[0] = map[int]uint32{0: state, 4: activate, 8: count | (count << 8), 20: anim}
					w.UpdateRefs[0] = map[int]int{12: 4}
					w.ItemWords[1][56] = math.Float32bits(640)
					p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestWorld53BEF0}, {Op: legacy.PortTestWorld53BEF0}}
					cases = append(cases, s)
				}
			}
		}
	}
	updateRegistrySpecsHash(t, "update-registry-world-teleports", cases)
}

func TestUpdateRegistryWorldTrapDoors(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, state := range []uint32{0, 2, 4, 6, 8, 10, 12} {
		for _, on := range []bool{false, true} {
			for _, stamp := range []uint32{0, 99, 100, 101, 0xffffffff} {
				s := updateRegistryWorldBase()
				p := s.Callbacks.Shop
				w := p.TemporaryUpdates
				if !on {
					p.Items[0].Flags = 4
				}
				w.ItemWords[0][20] = state
				w.World.CollideWords[0] = map[int]uint32{16: stamp, 24: 17}
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestWorld53DE80}, {Op: legacy.PortTestWorld53DE80}}
				cases = append(cases, s)
			}
		}
	}
	updateRegistrySpecsHash(t, "update-registry-world-trap-doors", cases)
}

func TestUpdateRegistryWorldPhantomAndAngleQueue(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, distance := range []float32{0, 399.99997, 400, 400.00003, 401} {
		for _, dir := range []uint32{0, 1, 127, 128, 255} {
			s := updateRegistryWorldBase()
			p := s.Callbacks.Shop
			w := p.TemporaryUpdates
			w.UpdateRefs[0] = map[int]int{0: 4}
			w.UpdateWords[0][4] = math.Float32bits(512 - distance)
			w.UpdateWords[0][8] = math.Float32bits(512)
			w.ItemWords[1][124] = dir
			p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestWorld53B860}}
			cases = append(cases, s)
		}
	}
	for _, angle := range []uint32{0, 1, 127, 128, 254, 255, 256, 65535} {
		for _, delta := range []int{-257, -2, -1, 0, 1, 2, 257} {
			s := updateRegistryWorldBase()
			p := s.Callbacks.Shop
			p.TemporaryUpdates.UpdateWords[0][40] = angle
			p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestWorld548860, Value: uint32(delta)}, {Op: legacy.PortTestWorld548830}, {Op: legacy.PortTestWorld548860, Item: 1, Value: uint32(delta)}}
			cases = append(cases, s)
		}
	}
	updateRegistrySpecsHash(t, "update-registry-world-phantom-queue", cases)
}

func TestUpdateRegistryWorldForces(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestWorld53B030, legacy.PortTestWorld53C160, legacy.PortTestWorld53C240} {
		for _, dir := range []uint32{0, 32, 64, 96, 128, 160, 192, 224, 255} {
			for _, delta := range [][2]float32{{0, 0}, {20, 0}, {-20, 0}, {0, 20}, {0, -20}, {20, 20}, {399.8, 0}, {400, 0}} {
				s := updateRegistryWorldBase()
				p := s.Callbacks.Shop
				w := p.TemporaryUpdates
				w.ItemWords[0][124] = dir
				w.ItemWords[1][56] = math.Float32bits(512 + delta[0])
				w.ItemWords[1][60] = math.Float32bits(512 + delta[1])
				w.UpdateWords[0] = map[int]uint32{0: math.Float32bits(400), 8: math.Float32bits(30)}
				w.Indexed = []int{4}
				p.Sequence = []legacy.PortTestShopAction{{Op: op}}
				cases = append(cases, s)
			}
		}
	}
	updateRegistrySpecsHash(t, "update-registry-world-forces", cases)
}

func TestUpdateRegistryWorldPositiveEffects(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	s := updateRegistryWorldBase()
	p := s.Callbacks.Shop
	p.TemporaryUpdates.UpdateWords[0] = map[int]uint32{0: 9, 8: 0, 36: 237}
	p.TemporaryUpdates.UpdateRefs[0] = map[int]int{4: 1}
	s.Owner.Frame = 101
	p.TemporaryUpdates.ItemWords[0][136] = 100
	p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestWorld53B060}}
	cases = append(cases, s)
	for _, op := range []int{legacy.PortTestWorld53B410, legacy.PortTestWorld53B750} {
		s := updateRegistryWorldBase()
		worldPlatform(&s, 0, 32)
		p := s.Callbacks.Shop
		w := p.TemporaryUpdates
		w.Target = 5
		w.UpdateRefs[0] = map[int]int{4: 4}
		w.UpdateWords[0][16] = 64
		w.UpdateWords[1][16] = 64
		w.ItemWords[1][56] = math.Float32bits(600)
		w.ItemWords[1][60] = math.Float32bits(600)
		w.ItemWords[1][184] = math.Float32bits(64)
		w.ItemWords[1][188] = math.Float32bits(64)
		if op == legacy.PortTestWorld53B750 {
			w.ItemWords[2][104] = math.Float32bits(64)
		}
		p.Sequence = []legacy.PortTestShopAction{{Op: op}}
		cases = append(cases, s)
	}
	s = updateRegistryWorldBase()
	p = s.Callbacks.Shop
	p.TemporaryUpdates.ItemWords[1][56] = math.Float32bits(532)
	p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestWorld53C240}}
	cases = append(cases, s)
	s = updateRegistryWorldBase()
	p = s.Callbacks.Shop
	p.TemporaryUpdates.Target = 1
	p.EffectsUse.Target = 1
	p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestWorld53C140}}
	cases = append(cases, s)
	r := effectsTimedRun(t, cases)
	want := []uint32{10, 10, 20, 54000, 70000}
	for i, v := range want {
		if len(r[0].Trace) <= i || r[0].Trace[i] != v {
			t.Fatalf("script event/caller/trigger: %v", r[0].Trace)
		}
	}
	for i, z := range []float32{64, 0} {
		o := worldObject(t, r[i+1], 0, 70002)
		if o[14] != math.Float32bits(600) || o[15] != math.Float32bits(600) || o[26] != math.Float32bits(z) {
			t.Fatalf("elevator %d position/height: %x %x %x", i, o[14], o[15], o[26])
		}
	}
	o := worldObject(t, r[3], 0, 70001)
	fx, fy := math.Float32frombits(o[22]), math.Float32frombits(o[23])
	if !(fx > 0) || math.IsInf(float64(fx), 0) || fy != 0 {
		t.Fatalf("directional force: %v,%v", fx, fy)
	}
	// EffectsUse captures the real resource player as its target.
	d := r[4].Callbacks.Shop.Sequence[0].EffectsUseData
	// No equipment accepts: seven scalar words, thirty modifier words, target-present marker.
	start := 38
	if d[start+14] != math.Float32bits(512) || d[start+15] != math.Float32bits(512) {
		t.Fatalf("player teleport: %x,%x", d[start+14], d[start+15])
	}
	updateRegistryHash(t, "update-registry-world-positive", r, updateRegistryHashes["update-registry-world-positive"])
}

func TestUpdateRegistryWorldTypeAndClockBoundaries(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, name := range []string{"Trigger", "PressurePlate", ""} {
		for _, fps := range []int{1, 4, 30, 60} {
			for _, frame := range []uint32{0, 1, 99, 100, 101, 0xffffffff} {
				s := updateRegistryWorldBase()
				s.Owner.FPS = uint32(fps)
				s.Owner.Frame = frame
				p := s.Callbacks.Shop
				w := p.TemporaryUpdates
				w.World.ItemNames = []string{name}
				w.UpdateWords[0] = map[int]uint32{0: 9, 8: 0, 36: 237}
				w.UpdateRefs[0] = map[int]int{4: 1}
				w.ItemWords[0][136] = 100
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestWorld53B1B0}, {Op: legacy.PortTestWorld53B1B0}}
				cases = append(cases, s)
			}
		}
	}
	for _, fps := range []int{1, 4, 30, 60} {
		for _, frame := range []uint32{0, 1, 0xfffffffe, 0xffffffff} {
			for _, op := range []int{legacy.PortTestWorld53B5D0, legacy.PortTestWorld53DE80, legacy.PortTestWorld53AC50} {
				s := updateRegistryWorldBase()
				s.Owner.FPS = uint32(fps)
				s.Owner.Frame = frame
				p := s.Callbacks.Shop
				w := p.TemporaryUpdates
				w.ItemWords[0][136] = frame - uint32(fps)
				w.UpdateWords[0] = map[int]uint32{4: 0, 8: 2, 12: 2, 40: 16, 44: frame - uint32(fps)/2}
				w.World.CollideWords[0] = map[int]uint32{16: frame, 24: 17}
				p.Sequence = []legacy.PortTestShopAction{{Op: op}}
				cases = append(cases, s)
			}
		}
	}
	updateRegistrySpecsHash(t, "update-registry-world-type-clock", cases)
}

func TestUpdateRegistryWorldOwnerIntegration(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, shaft := range []bool{false, true} {
		s := updateRegistryWorldBase()
		worldPlatform(&s, 0, 32)
		p := s.Callbacks.Shop
		w := p.TemporaryUpdates
		p.Items[2].Class = uint32(object.ClassSimple)
		w.Indexed = []int{5}
		w.UpdateRefs[0] = map[int]int{4: 4}
		w.ItemWords[1][56] = math.Float32bits(600)
		w.ItemWords[1][60] = math.Float32bits(600)
		w.ItemWords[1][184] = math.Float32bits(64)
		w.ItemWords[1][188] = math.Float32bits(64)
		op := legacy.PortTestWorld53B5D0
		w.UpdateWords[0] = map[int]uint32{12: 3, 16: 62}
		w.ItemWords[2][104] = math.Float32bits(64)
		if shaft {
			op = legacy.PortTestWorld53B380
			w.UpdateWords[1] = map[int]uint32{12: 1, 16: 32}
			w.ItemWords[2][104] = math.Float32bits(-32)
		}
		p.Sequence = []legacy.PortTestShopAction{{Op: op}}
		cases = append(cases, s)
	}
	s := updateRegistryWorldBase()
	p := s.Callbacks.Shop
	w := p.TemporaryUpdates
	s.Lifecycle.GameFlags = 2048
	p.Items[2].Class = uint32(object.ClassSimple)
	w.Indexed = []int{5}
	w.UpdateWords[0] = map[int]uint32{0: 1, 8: 0, 20: 7}
	w.UpdateRefs[0] = map[int]int{12: 4}
	w.ItemWords[1][56] = math.Float32bits(640)
	p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestWorld53BEF0}}
	cases = append(cases, s)
	s = updateRegistryWorldBase()
	p = s.Callbacks.Shop
	w = p.TemporaryUpdates
	w.Indexed = []int{4}
	w.ItemWords[1][56] = math.Float32bits(532)
	p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestWorld53C160}}
	cases = append(cases, s)
	r := effectsTimedRun(t, cases)
	for i, z := range []float32{0, 32} {
		o := worldObject(t, r[i], 0, 70002)
		if o[14] != math.Float32bits(600) || o[15] != math.Float32bits(600) || o[26] != math.Float32bits(z) {
			t.Fatalf("owner elevator %d: %x,%x,%x", i, o[14], o[15], o[26])
		}
	}
	o := worldObject(t, r[2], 0, 70002)
	if o[14] != math.Float32bits(640) || o[15] != math.Float32bits(512) {
		t.Fatalf("visible teleporter: %x,%x", o[14], o[15])
	}
	o = worldObject(t, r[3], 0, 70001)
	if !(math.Float32frombits(o[22]) > 0) || o[23] != 0 {
		t.Fatalf("blow owner force: %x,%x", o[22], o[23])
	}
	updateRegistryHash(t, "update-registry-world-owner-integration", r, updateRegistryHashes["update-registry-world-owner-integration"])
}

func TestUpdateRegistryWorldInvisibleTeleport(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestWorld53C0C0} {
		for _, class := range []uint32{1, 2, 4, 0x20000, 0x400000, 0x420000} {
			for _, on := range []bool{false, true} {
				s := updateRegistryWorldBase()
				p := s.Callbacks.Shop
				w := p.TemporaryUpdates
				p.Items[1].Class = class
				if !on {
					p.Items[0].Flags = 4
				}
				w.UpdateWords[0][4] = 1
				w.UpdateRefs[0] = map[int]int{12: 5}
				w.ItemWords[2][56] = math.Float32bits(640)
				w.Indexed = []int{4}
				if class == 4 {
					// Use the actual resource player and its guarded player data.
					p.Items[1].Class = 1
					w.Target = 1
					w.Indexed = []int{1}
					if op == legacy.PortTestWorld53C0C0 {
						w.ItemWords[0][56] = math.Float32bits(180)
						w.ItemWords[0][60] = math.Float32bits(100)
					}
				}
				p.Sequence = []legacy.PortTestShopAction{{Op: op}}
				cases = append(cases, s)
			}
		}
	}
	results := effectsTimedRun(t, cases)
	for i, r := range results {
		d := r.Callbacks.Shop.Sequence[0].TemporaryUpdatesData
		// Three prefix words, callback trace, seven cache words, then update data.
		if d[3+int(d[2])+7+1] != 0 {
			t.Fatalf("case%d activation not cleared", i)
		}
	}
	for i, wantX := range []float32{512, 640} {
		target := worldObject(t, results[2+i], 0, 70001)
		if target[14] != math.Float32bits(wantX) || target[15] != math.Float32bits(512) {
			t.Fatalf("invisible teleport powered%d target %x,%x", i, target[14], target[15])
		}
	}
	updateRegistryHash(t, "update-registry-world-invisible-teleport", results, updateRegistryHashes["update-registry-world-invisible-teleport"])
}
