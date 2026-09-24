//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/internal/netlist"
	"testing"
	"unsafe"
)

func TestCollisionRegistryWorldAnkhHistory(t *testing.T) {
	o := newWorldCollisionOwner(t)
	t.Cleanup(o.s.PortTestWorldAnkhType())
	a := newObjectXferSimple(t, o.s)
	b := &o.units[1]
	other := &o.units[2]
	data := o.record(t, 5124)
	oldData := a.InitData
	a.InitData = data
	t.Cleanup(func() { a.InitData = oldData })
	player := b.UpdateDataPlayer().Player.C()
	copy(unsafe.Slice((*byte)(unsafe.Add(player, 2185)), 50), []byte{'A', 0, 0, 0})
	*(*byte)(unsafe.Add(player, 2251)) = 1
	copy(unsafe.Slice((*byte)(unsafe.Add(player, 2112)), 25), []byte("fixture\x00"))
	var rows []struct {
		Name         string
		Lives, Frame uint32
		Index        byte
		History      [5]uint32
		Records      []byte
		Sounds       []int
		Bytes        []byte
	}
	defer func() {
		collisionRegistryCapture(t, "collision-registry-world-ankh-history", rows)
	}()
	for _, mode := range []string{"new", "identity", "visited", "full-history"} {
		for _, frame := range []uint32{7199, 7200, 7201} {
			for _, index := range []byte{0, 63} {
				for _, lives := range []uint32{4, 5} {
					for _, ticks := range []uint64{1500, 1501} {
						name := fmt.Sprintf("%s/frame%d/index%d/lives%d/ticks%d", mode, frame, index, lives, ticks)
						t.Run(name, func(t *testing.T) {
							o.reset()
							o.s.SetFrame(frame)
							o.s.PortTestCombatAudioReset()
							o.ticks = ticks
							*o.throttle = 0
							clear(unsafe.Slice((*byte)(data), 5124))
							*(*byte)(unsafe.Add(data, 5120)) = index
							objectXferSetWord(b.UpdateData, 320, lives)
							a.Field34 = 17
							for i := 0; i < 5; i++ {
								objectXferSetWord(player, 4796+i*4, 0)
							}
							if mode == "visited" {
								objectXferSetWord(player, 4796, uint32(uintptr(a.CObj())))
							}
							if mode == "full-history" {
								for i := 0; i < 5; i++ {
									objectXferSetWord(player, 4796+i*4, uint32(uintptr(other.CObj())))
								}
							}
							if mode == "identity" {
								record := unsafe.Add(data, 7*80)
								copy(unsafe.Slice((*byte)(record), 50), unsafe.Slice((*byte)(unsafe.Add(player, 2185)), 50))
								*(*byte)(unsafe.Add(record, 50)) = 1
								copy(unsafe.Slice((*byte)(unsafe.Add(record, 51)), 25), unsafe.Slice((*byte)(unsafe.Add(player, 2112)), 25))
							}
							defer func() {
								for u := o.s.Objs.DeletedList; u != nil; {
									next := u.DeletedNext
									o.s.Objs.FreeObject(u)
									u = next
								}
								o.s.Objs.DeletedList = nil
							}()
							collisionRegistryWorld(20, a, b, nil)
							seen := mode == "visited" || mode == "identity" && frame <= 7200
							awarded := !seen && lives < 5
							wantLives, wantFrame, wantIndex := lives, uint32(17), index
							if awarded {
								wantLives++
								wantFrame = frame
								wantIndex = (index + 1) % 64
							}
							if objectXferGetWord(b.UpdateData, 320) != wantLives || a.Field34 != wantFrame || *(*byte)(unsafe.Add(data, 5120)) != wantIndex {
								t.Fatal("ankh award state")
							}
							var history [5]uint32
							for i := range history {
								v := objectXferGetWord(player, 4796+i*4)
								switch v {
								case 0:
								case uint32(uintptr(a.CObj())):
									history[i] = 1
								case uint32(uintptr(other.CObj())):
									history[i] = 2
								default:
									t.Fatal("unknown history pointer")
								}
							}
							if mode == "full-history" {
								for _, v := range history {
									if v != 2 {
										t.Fatal("full history overwrite")
									}
								}
							} else if (seen || awarded) && history[0] != 1 {
								t.Fatal("missing history")
							}
							var sounds []int
							for _, e := range o.s.PortTestCombatAudioSnapshot() {
								sounds = append(sounds, int(e.ID))
							}
							wantSounds := 0
							if awarded {
								wantSounds = 2
							} else if ticks > 1500 {
								wantSounds = 1
							}
							if len(sounds) != wantSounds {
								t.Fatalf("ankh sounds %v want%d", sounds, wantSounds)
							}
							if awarded {
								if sounds[0] != 1004 || sounds[1] != 1004 || o.s.Objs.DeletedList == nil {
									t.Fatal("tradable pickup")
								}
							} else if len(sounds) > 0 && sounds[0] != 925 {
								t.Fatal("ankh rejection sound")
							}
							rows = append(rows, struct {
								Name         string
								Lives, Frame uint32
								Index        byte
								History      [5]uint32
								Records      []byte
								Sounds       []int
								Bytes        []byte
							}{name, objectXferGetWord(b.UpdateData, 320), a.Field34, *(*byte)(unsafe.Add(data, 5120)), history, append([]byte(nil), unsafe.Slice((*byte)(data), 5120)...), sounds, o.s.NetList.CopyPacketsA(7, netlist.Kind1)})
						})
					}
				}
			}
		}
	}
	collisionRegistryWorld(20, a, nil, nil)
}
