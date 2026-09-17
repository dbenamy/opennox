//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/prand"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
	"unsafe"
)

func TestVisibilityEffectsScanDelay(t *testing.T) {
	s, u, others, calls, setup := visibilitySeenOwner(t)
	s.PortTestAIEmptyMap()
	t.Cleanup(s.Map.Free)
	for off, data := range blobdata.PortTestCombatTables() {
		dst := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, off)), len(data))
		old := append([]byte(nil), dst...)
		copy(dst, data)
		t.Cleanup(func() { copy(dst, old) })
	}
	forced := memmap.PtrUint32(0x5D4594, 2487684)
	old := *forced
	*forced = 0
	t.Cleanup(func() { *forced = old })
	u.PosVec = types.Pointf{1000, 1000}
	u.PrevPos = u.PosVec
	u.Direction1 = 0
	u.TeamPtr().ID = 1
	v := others[0]
	v.TeamPtr().ID = 2
	type row struct {
		Name     string
		Times    [4]uint32
		Distance uint32
		Random   int
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "visibility-effects-scan-delay", rows, "0db5f83c5cf17aaf3c9552fc3b924247ea22c628453e5a9617362b066355b50d")
	}()
	for _, quest := range []bool{false, true} {
		for _, frame := range []uint32{100, 0xfffffffe} {
			for _, fps := range []uint32{0, 30, 60} {
				for _, custom := range []int{0, 250, 640, 999, 1000, 1200} {
					for _, distance := range []int{-1, 0, 1, 249, 250, 251, 639, 640, 641, 999, 1000, 1001} {
						name := fmt.Sprintf("quest%t/frame%x/fps%d/radius%d/distance%d", quest, frame, fps, custom, distance)
						t.Run(name, func(t *testing.T) {
							flags := noxflags.GameFlag(0x200000)
							if quest {
								flags |= 4096
							}
							defer noxflags.PortTestGameFlags(flags)()
							s.PortTestAIEmptyMap()
							setup(0)
							s.SetFrame(frame)
							s.SetTickRate(fps)
							s.Rand.Logic = prand.New(12345)
							*(*float32)(unsafe.Add(u.UpdateData, 1312)) = float32(custom)
							if distance >= 0 {
								v.PosVec = types.Pointf{float32(1000 - distance), 1000}
								v.NewPos = v.PosVec
								v.ObjFlags = object.FlagActive
								v.Shape.Kind = server.ShapeKindCircle
								v.Shape.Circle.R = 1
								v.Shape.Circle.R2 = 1
								s.Map.AddObjectToIndex(v)
							}
							rv := legacy.PortTestVisibilityEffects(17, u, nil, nil, nil, [5]int32{}, nil, "")
							radius := 250
							if quest {
								radius = 640
							}
							if custom > radius {
								radius = custom
							}
							found := distance >= 0 && distance <= 1000
							expectedRNG := prand.New(12345)
							delay := uint32(5 * fps)
							wantDistance := float32(-1)
							if found {
								wantDistance = float32(distance)
								if distance > radius {
									delay = uint32((distance-radius)*int(5*fps)/(1000-radius) + 10)
								} else {
									delay = uint32(expectedRNG.IntClamp(5, 10))
								}
							}
							r := row{Name: name, Distance: objectXferGetWord(u.UpdateData, 524), Random: s.Rand.Logic.Index()}
							for i, off := range []int{1204, 1208, 1212, 1200} {
								r.Times[i] = objectXferGetWord(u.UpdateData, off)
							}
							if rv != 0 || *(*byte)(unsafe.Add(u.UpdateData, 1129)) != 0 || len(*calls) != 0 {
								t.Fatalf("unexpected enemy/callback return%d calls%v", rv, *calls)
							}
							if r.Times != [4]uint32{frame, frame + delay, frame, 0} || r.Distance != math.Float32bits(wantDistance) || r.Random != expectedRNG.Index() {
								t.Fatalf("state%+v expected deadline%x distance%g index%d", r, frame+delay, wantDistance, expectedRNG.Index())
							}
							rows = append(rows, r)
						})
					}
				}
			}
		}
	}
}
