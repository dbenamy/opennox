//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestWorldCollisionsDoorKeys(t *testing.T) {
	o := newWorldCollisionOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	names := []string{"SilverKey", "GoldKey", "RubyKey", "SapphireKey"}
	t.Cleanup(o.s.PortTestAttackTypes(1, nil, names...))
	t.Cleanup(o.s.PortTestSpellEffectTypes(names, nil, nil))
	a := newObjectXferSimple(t, o.s)
	near := newObjectXferSimple(t, o.s)
	b := &o.units[1]
	a.UpdateData = o.record(t, 64)
	near.UpdateData = o.record(t, 64)
	t.Cleanup(func() { a.UpdateData = nil; near.UpdateData = nil; b.InvFirstItem = nil; o.s.Objs.DeletedList = nil })
	near.ObjClass = object.ClassDoor
	near.ObjFlags = object.FlagActive
	near.PosVec = types.Ptf(115, 115)
	near.NewPos = near.PosVec
	near.Shape.Kind = server.ShapeKindCircle
	near.Shape.Circle.R = 1
	near.Shape.Circle.R2 = 1
	o.s.Map.AddObjectToIndex(near)
	var keys []*server.Object
	for _, name := range names {
		key := o.s.NewObjectByTypeID(name)
		if key == nil {
			t.Fatal("key factory")
		}
		key.ObjClass = object.ClassKey
		keys = append(keys, key)
		t.Cleanup(func() { key.InvHolder = nil; key.InvNextItem = nil; o.s.Objs.FreeObject(key) })
	}
	var rows []struct {
		Name                               string
		Lock, NearLock, Notify, NearNotify byte
		Deleted                            bool
		Sounds                             []int
	}
	defer func() {
		spellbookCapture(t, "world-collisions-door-keys", rows, "64ab2d4e3b6894ab846f90fd72e98ea3f9648105ff96dc038d58d541558d4c30")
	}()
	for lock := 1; lock <= 4; lock++ {
		for _, angle := range []uint32{0, 8, 16, 24} {
			for _, quest := range []bool{false, true} {
				for _, adjacent := range []bool{false, true} {
					name := fmt.Sprintf("key%d/angle%d/quest%t/adjacent%t", lock, angle, quest, adjacent)
					t.Run(name, func(t *testing.T) {
						o.reset()
						o.s.PortTestCombatAudioReset()
						flags := noxflags.GameFlag(0)
						if quest {
							flags = noxflags.GameModeQuest
						}
						defer noxflags.PortTestGameFlags(flags)()
						clear(unsafe.Slice((*byte)(a.UpdateData), 64))
						clear(unsafe.Slice((*byte)(near.UpdateData), 64))
						*(*byte)(unsafe.Add(a.UpdateData, 1)) = byte(lock)
						*(*byte)(unsafe.Add(near.UpdateData, 1)) = byte(lock)
						objectXferSetWord(a.UpdateData, 4, angle)
						objectXferSetWord(a.UpdateData, 12, angle)
						objectXferSetWord(a.UpdateData, 16, 5)
						objectXferSetWord(a.UpdateData, 20, 5)
						x, y := uint32(4), uint32(4)
						if angle == 8 || angle == 16 {
							x = 6
						}
						if angle == 16 || angle == 24 {
							y = 6
						}
						if !adjacent {
							x = 9
						}
						objectXferSetWord(near.UpdateData, 16, x)
						objectXferSetWord(near.UpdateData, 20, y)
						key := keys[lock-1]
						key.ObjFlags = 0
						key.DeletedNext = nil
						key.InvHolder = b
						b.InvFirstItem = key
						o.s.Objs.DeletedList = nil
						legacy.PortTestWorldCollision(1, a, b, nil)
						if *(*byte)(unsafe.Add(a.UpdateData, 1)) != 0 || !key.ObjFlags.Has(object.FlagDestroyed) {
							t.Fatal("key unlock/consume")
						}
						wantNear := byte(lock)
						if adjacent {
							wantNear = 0
						}
						if *(*byte)(unsafe.Add(near.UpdateData, 1)) != wantNear {
							t.Fatal("adjacent door unlock")
						}
						notify := *(*byte)(unsafe.Add(a.UpdateData, 48))
						nearNotify := *(*byte)(unsafe.Add(near.UpdateData, 48))
						if (notify == 1) != quest || (nearNotify == 1) != (quest && adjacent) {
							t.Fatal("quest door notification")
						}
						var sounds []int
						for _, e := range o.s.PortTestCombatAudioSnapshot() {
							sounds = append(sounds, int(e.ID))
						}
						if len(sounds) != 1 || sounds[0] != 234 {
							t.Fatal("unlock sound")
						}
						rows = append(rows, struct {
							Name                               string
							Lock, NearLock, Notify, NearNotify byte
							Deleted                            bool
							Sounds                             []int
						}{name, 0, wantNear, notify, nearNotify, true, sounds})
					})
				}
			}
		}
	}
}

func TestWorldCollisionsDoorMagic(t *testing.T) {
	o := newWorldCollisionOwner(t)
	a := newObjectXferSimple(t, o.s)
	b := &o.units[1]
	data := o.record(t, 64)
	a.UpdateData = data
	t.Cleanup(func() { a.UpdateData = nil; a.ObjOwner = nil })
	var rows []struct {
		Name     string
		Owner    uint32
		Sounds   []int
		Throttle uint64
	}
	defer func() {
		spellbookCapture(t, "world-collisions-door-magic", rows, "d41231ffbd2e616fd150710374198423df217da8601731f7b778f215ab346cde")
	}()
	for _, frame := range []uint32{0, 123, 0xffffffff} {
		for _, expires := range []uint32{0, 122, 123, 124, 0xffffffff} {
			for _, owner := range []int{0, 1, 2} {
				for _, closed := range []bool{false, true} {
					for _, gate := range []bool{false, true} {
						name := fmt.Sprintf("frame%d/expires%d/owner%d/closed%t/gate%t", frame, expires, owner, closed, gate)
						t.Run(name, func(t *testing.T) {
							o.reset()
							o.s.SetFrame(frame)
							o.s.PortTestCombatAudioReset()
							*o.throttle = 0
							a.ObjOwner = nil
							if owner == 1 {
								a.ObjOwner = b
							}
							if owner == 2 {
								a.ObjOwner = &o.units[2]
							}
							a.Field34 = expires
							a.ObjSubClass = 0
							if gate {
								a.ObjSubClass = 4
							}
							objectXferSetWord(data, 4, 0)
							objectXferSetWord(data, 12, 1)
							if closed {
								objectXferSetWord(data, 12, 0)
							}
							legacy.PortTestWorldCollision(1, a, b, nil)
							wantOwner := owner
							if closed && expires <= frame {
								wantOwner = 0
							}
							gotOwner := 0
							if a.ObjOwner == b {
								gotOwner = 1
							} else if a.ObjOwner == &o.units[2] {
								gotOwner = 2
							} else if a.ObjOwner != nil {
								t.Fatal("magic lock owner")
							}
							if gotOwner != wantOwner {
								t.Fatal("magic lock expiry")
							}
							var sounds []int
							for _, e := range o.s.PortTestCombatAudioSnapshot() {
								sounds = append(sounds, int(e.ID))
							}
							fired := closed && owner == 2 && expires > frame
							if (len(sounds) == 1) != fired || len(sounds) > 1 {
								t.Fatal("magic lock sound count")
							}
							if fired {
								want := 240
								if gate {
									want = 244
								}
								if sounds[0] != want || *o.throttle != 10000 {
									t.Fatal("magic lock notification")
								}
							} else if *o.throttle != 0 {
								t.Fatal("magic lock throttle")
							}
							rows = append(rows, struct {
								Name     string
								Owner    uint32
								Sounds   []int
								Throttle uint64
							}{name, uint32(gotOwner), sounds, *o.throttle})
						})
					}
				}
			}
		}
	}
}

func TestWorldCollisionsDoorMissingKey(t *testing.T) {
	o := newWorldCollisionOwner(t)
	a := newObjectXferSimple(t, o.s)
	b := &o.units[1]
	data := o.record(t, 64)
	a.UpdateData = data
	t.Cleanup(func() { a.UpdateData = nil })
	var rows []struct {
		Key    int
		Gate   bool
		Packet []byte
	}
	for key := 1; key <= 4; key++ {
		for _, gate := range []bool{false, true} {
			o.reset()
			*o.throttle = 0
			a.ObjSubClass = 0
			if gate {
				a.ObjSubClass = 4
			}
			*(*byte)(unsafe.Add(data, 1)) = byte(key)
			b.InvFirstItem = nil
			legacy.PortTestWorldCollision(1, a, b, nil)
			nodes := o.state().Nodes
			if len(nodes) != 1 || len(nodes[0].Data) != 52 || nodes[0].Data[0] != 240 || nodes[0].Data[1] != 33 || nodes[0].Data[51] != byte(key) || nodes[0].To != 7 {
				t.Fatal("missing key message encoding")
			}
			rows = append(rows, struct {
				Key    int
				Gate   bool
				Packet []byte
			}{key, gate, nodes[0].Data})
		}
	}
	spellbookCapture(t, "world-collisions-door-missing-key", rows, "ede7ef606d68833b7202b8a7af85698c7da6fc3e26cd26c6ff1c9b16e97727d5")
}
