//go:build porttest

package legacy

/*
#include "GAME5.h"
extern uint32_t dword_5d4594_2491716;
*/
import "C"
import (
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"math"
	"unsafe"
)

type PortTestGeneratorSpec struct {
	SpawnPolicy                                      *PortTestGeneratorSpawnPolicySpec
	Inventory                                        *PortTestGeneratorInventorySpec
	Frame                                            uint32
	Sources                                          byte
	UseDef                                           bool
	DefHealth, SourceSpellWord                       uint32
	PlayerPos                                        [2]uint32
	View                                             [2]uint16
	PlayerFlags, PlayerStatus                        uint32
	Joined                                           bool
	Killer, Beholder                                 bool
	QuestState, HealthScale, SourceFill              uint32
	Cache                                            [8]uint32
	Balance                                          map[string]float64
	Stage, Level, GenObjFlags, GenXStatus, LastSpawn uint32
	Current, Limit, RateClass                        byte
	BlockerAtSource                                  bool
	Op                                               int
	Tile                                             int
	Cold                                             bool
	Flags                                            uint32
	Radius                                           uint32
	Point                                            [2]uint32
	TowardPlayer                                     bool
	BlockerRadius                                    uint32
}
type PortTestGeneratorResult struct {
	Objects     *PortTestGeneratorObjectsResult `json:",omitempty"`
	Data, Cache []uint32
	TileCache   [5]uint32
	Point       [2]uint32
	Occupied    uint32
	Intact      bool
}
type portTestGeneratorState struct {
	objects          *portTestGeneratorObjects
	configureBalance func(map[string]float64)
	configureTiles   func(int, bool)
	tilesIntact      func() bool
	point            []uint32
	spec             *PortTestGeneratorSpec
}

func portTestGeneratorEnvironment(proxy *portTestRoamOwnerServer) func() {
	configure, intact, freeTiles := portTestGeneratorTileEnvironment()
	_, configureBalance, freeServer := proxy.core.PortTestGeneratorEnvironment()
	oldHardcore := C.dword_5d4594_2491716
	offsets := []uintptr{2491712, 2491720, 2491724, 2491728, 2491732, 2491736, 2491740, 2491744, 2388660}
	oldWords := make([]uint32, len(offsets))
	for i, o := range offsets {
		oldWords[i] = *memmap.PtrUint32(0x5D4594, o)
	}
	oldStage := *memmap.PtrUint32(0x587000, 202028)
	p, freePoint := alloc.Make([]uint32{}, 6)
	oldOccupied := *memmap.PtrUint32(0x5D4594, 2491708)
	proxy.callbacks.generation = &portTestGeneratorState{configureTiles: configure, tilesIntact: intact, point: p, configureBalance: configureBalance}
	freeObjects := portTestGeneratorObjectsEnvironment(proxy)
	return func() {
		*memmap.PtrUint32(0x5D4594, 2491708) = oldOccupied
		C.dword_5d4594_2491716 = oldHardcore
		for i, o := range offsets {
			*memmap.PtrUint32(0x5D4594, o) = oldWords[i]
		}
		*memmap.PtrUint32(0x587000, 202028) = oldStage
		freeObjects()
		freeServer()
		freePoint()
		freeTiles()
	}
}
func portTestGeneratorPrepare(proxy *portTestRoamOwnerServer, u *server.Object, sp *PortTestGeneratorSpec) {
	st := proxy.callbacks.generation
	st.spec = sp
	proxy.core.SetFrame(sp.Frame)
	st.configureTiles(sp.Tile, sp.Cold)
	// Prime only tile-independent vacancy cases so their guard check expects
	// the same known cache values without inventing a generator tile lookup.
	if sp.Op == 0 || sp.Op >= 3 {
		st.configureTiles(sp.Tile, false)
	}
	st.point[0], st.point[1], st.point[4], st.point[5] = 0xa5a5a5a5, 0xa5a5a5a5, 0x5a5a5a5a, 0x5a5a5a5a
	st.point[2], st.point[3] = sp.Point[0], sp.Point[1]
	*memmap.PtrUint32(0x5D4594, 2491708) = 0x12345678
	*(*uint32)(unsafe.Add(u.UpdateData, 92)) = sp.Flags
	t := proxy.combat.target
	if t.ObjFlags.Has(object.FlagPartitioned) {
		for i := uint32(0); i < t.ObjIndexCur; i++ {
			proxy.core.Map.Sub5178E0(true, &t.ObjIndex[i])
		}
		proxy.core.Map.Sub5178E0(false, &t.ObjIndexBase)
		t.ObjFlags &^= object.FlagPartitioned
		t.ObjIndexCur = 0
		t.ObjIndex = [4]server.ObjectIndex{}
		t.ObjIndexBase = server.ObjectIndex{}
	}
	if sp.BlockerAtSource {
		t.PosVec = u.PosVec
	}
	t.Shape.Kind = server.ShapeKindCircle
	t.Shape.Circle.R = math.Float32frombits(sp.BlockerRadius)
	t.Shape.Circle.R2 = t.Shape.Circle.R * t.Shape.Circle.R
	t.NewPos = t.PosVec
	t.CollideP1 = types.Pointf{X: t.PosVec.X - t.Shape.Circle.R, Y: t.PosVec.Y - t.Shape.Circle.R}
	t.CollideP2 = types.Pointf{X: t.PosVec.X + t.Shape.Circle.R, Y: t.PosVec.Y + t.Shape.Circle.R}
	if sp.BlockerRadius != 0 {
		proxy.core.Map.AddObjectToIndex(t)
	}
	proxy.life.ids[uint32(uintptr(unsafe.Pointer(&st.point[2])))] = 996
	st.configureBalance(sp.Balance)
	C.dword_5d4594_2491716 = C.uint32_t(sp.Cache[0])
	for i := 1; i < 8; i++ {
		*memmap.PtrUint32(0x5D4594, 2491716+uintptr(i*4)) = sp.Cache[i]
	}
	*memmap.PtrUint32(0x587000, 202028) = sp.Stage
	*memmap.PtrUint32(0x5D4594, 2388660) = sp.Level
	if sp.Op == 3 {
		clear(unsafe.Slice((*byte)(u.UpdateData), 164))
		u.ObjFlags = object.Flags(sp.GenObjFlags)
		u.Field5 = sp.GenXStatus
		*(*byte)(unsafe.Add(u.UpdateData, 80+sp.Level)) = sp.RateClass
		*(*byte)(unsafe.Add(u.UpdateData, 86)) = sp.Current
		*(*byte)(unsafe.Add(u.UpdateData, 87)) = sp.Limit
		*(*uint32)(unsafe.Add(u.UpdateData, 88)) = sp.LastSpawn
	}
	portTestGeneratorObjectsPrepare(proxy, u, sp)
}
func portTestGeneratorCall(proxy *portTestRoamOwnerServer, u *server.Object) uint32 {
	st := proxy.callbacks.generation
	sp := st.spec
	p := unsafe.Pointer(&st.point[2])
	t := proxy.combat.target
	if sp.Op >= 4 {
		return portTestGeneratorObjectsCall(proxy, u)
	}
	switch sp.Op {
	case 0:
		return generatorOccupied(*(*types.Pointf)(p))
	case 1:
		return uint32(generatorRadial(math.Float32frombits(sp.Radius), u.PosVec, (*types.Pointf)(p), t))
	case 2:
		var player *server.Object
		if sp.TowardPlayer {
			player = &proxy.life.players[0]
		}
		return uint32(generatorPlace(u, (*types.Pointf)(p), player, t))
	case 3:
		return uint32(int32(C.nox_xxx_updateMonsterGenerator_54E930((*C.uint32_t)(u.CObj()))))
	default:
		panic("generator operation")
	}
}
func portTestGeneratorTrace(proxy *portTestRoamOwnerServer, normalize func(uint32) uint32) *PortTestGeneratorResult {
	st := proxy.callbacks.generation
	p := st.point
	r := &PortTestGeneratorResult{Point: [2]uint32{p[2], p[3]}, Occupied: *memmap.PtrUint32(0x5D4594, 2491708), Intact: p[0] == 0xa5a5a5a5 && p[1] == 0xa5a5a5a5 && p[4] == 0x5a5a5a5a && p[5] == 0x5a5a5a5a && st.tilesIntact()}
	for i := range r.TileCache {
		r.TileCache[i] = *memmap.PtrUint32(0x587000, 26516+uintptr(i*4))
	}
	for off := uintptr(0); off < 164; off += 4 {
		r.Data = append(r.Data, normalize(*(*uint32)(unsafe.Add(proxy.combat.actor.UpdateData, off))))
	}
	r.Cache = append(r.Cache, uint32(C.dword_5d4594_2491716))
	for i := 1; i < 8; i++ {
		r.Cache = append(r.Cache, *memmap.PtrUint32(0x5D4594, 2491716+uintptr(i*4)))
	}
	if st.spec.Op >= 4 {
		r.Objects = portTestGeneratorObjectsTrace(proxy, normalize)
		r.Intact = r.Intact && r.Objects.Intact
	}
	return r
}
