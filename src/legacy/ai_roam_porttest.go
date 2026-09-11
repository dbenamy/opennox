//go:build porttest

package legacy

import (
	"bytes"
	"encoding/binary"
	"math"
	"time"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/prand"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/unit/ai"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

type PortTestRoamSpec struct {
	Lifecycle                  *PortTestLifecycleSpec
	Combat                     *PortTestCombatSpec
	Path                       *PortTestPathSpec
	Navigation                 *PortTestNavigationSpec
	GuardEscort                *PortTestGuardEscortSpec
	Owner                      *PortTestRoamOwnerSpec
	Op, Seed                   int
	Index, Insert, Count, Mask byte
	Stack                      int8
	History                    [16]byte
	Neighbors                  [32]byte
	Flags                      [34]byte
	Enabled                    [34]bool
}
type PortTestRoamResult struct {
	Lifecycle          *PortTestLifecycleResult `json:",omitempty"`
	Combat             *PortTestCombatResult    `json:",omitempty"`
	Nanos              int64                    `json:"-"`
	Trace              []uint32                 `json:",omitempty"`
	History            [16]byte
	Index, Arg, Field2 uint32
	Return             int
	Stack              int8
	Logic, Other       int
	Changed            bool
	Changes            []uint32
	Intact             bool
}

// PortTestRoam runs a shared guarded fixture through registered actions and native helpers.
// Pointer-bearing output is normalized to stable waypoint IDs before hashing.
func PortTestRoam(specs []PortTestRoamSpec) []PortTestRoamResult {
	oldGame := noxflags.GetGame()
	offsets := []uintptr{2490500, 2489452, 2489444}
	for _, sp := range specs {
		if sp.Path != nil {
			offsets = append(offsets, 2386204)
			restorePath := portTestPathEnvironment()
			defer restorePath()
			break
		}
	}
	oldWords := make([]uint32, len(offsets))
	for i, off := range offsets {
		oldWords[i] = *memmap.PtrUint32(0x5D4594, off)
	}
	defer func() {
		noxflags.ResetGame()
		noxflags.SetGame(oldGame)
		for i, off := range offsets {
			*memmap.PtrUint32(0x5D4594, off) = oldWords[i]
		}
	}()
	hb, fh := alloc.Make([]byte{}, int(unsafe.Sizeof(server.HealthData{}))+16)
	defer fh()
	health := (*server.HealthData)(unsafe.Pointer(&hb[8]))
	for _, sp := range specs {
		if sp.GuardEscort != nil {
			restore := portTestGuardEscortEnvironment()
			defer restore()
			break
		}
	}
	core := new(server.Server)
	core.SetFrame(123)
	configureWalls, wallsUnchanged, freeWalls := core.PortTestPathWalls()
	defer freeWalls()
	restoreTypes := core.PortTestObjectInitSize(1, 0)
	defer restoreTypes()
	oldGet, oldFlags := GetServer, noxflags.GetEngine()
	proxy := &portTestRoamOwnerServer{portTestRandomServer: portTestRandomServer{core: core}}
	for _, sp := range specs {
		if sp.Combat != nil {
			restore := portTestCombatEnvironment(proxy)
			defer restore()
			break
		}
	}
	for _, sp := range specs {
		if sp.Lifecycle != nil {
			restore := portTestLifecycleEnvironment(proxy)
			defer restore()
			break
		}
	}
	GetServer = func() Server { return proxy }
	noxflags.UnsetEngine(noxflags.EngineShowAI)
	defer func() { GetServer = oldGet; noxflags.ResetEngine(); noxflags.SetEngine(oldFlags) }()
	ob, fo := alloc.Make([]byte{}, int(unsafe.Sizeof(server.Object{}))+16)
	defer fo()
	ub, fu := alloc.Make([]byte{}, int(unsafe.Sizeof(server.MonsterUpdateData{}))+16)
	defer fu()
	stride := int(unsafe.Sizeof(server.Waypoint{})) + 16
	wb, fw := alloc.Make([]byte{}, 34*stride)
	defer fw()
	db, fd := alloc.Make([]byte{}, int(unsafe.Sizeof(server.MonsterDef{}))+16)
	defer fd()
	tb, ft := alloc.Make([]byte{}, int(unsafe.Sizeof(server.Object{}))+16)
	defer ft()
	def := (*server.MonsterDef)(unsafe.Pointer(&db[8]))
	target := (*server.Object)(unsafe.Pointer(&tb[8]))
	obj := (*server.Object)(unsafe.Pointer(&ob[8]))
	ud := (*server.MonsterUpdateData)(unsafe.Pointer(&ub[8]))
	detach := server.PortTestAttachAI(core, obj, target)
	defer detach()
	raw := func(id byte) uint32 {
		if id == 0 {
			return 0
		}
		return uint32(uintptr(unsafe.Pointer(&wb[int(id)*stride+8])))
	}
	ids := map[uint32]uint32{0: 0, uint32(uintptr(unsafe.Pointer(target))): 100}
	if proxy.combat != nil {
		ids[uint32(uintptr(unsafe.Pointer(obj)))] = 101
		ids[uint32(uintptr(unsafe.Pointer(proxy.combat.weapon)))] = 102
		for i, b := range proxy.combat.extra {
			ids[uint32(uintptr(unsafe.Pointer(&b[8])))] = uint32(400 + i)
		}
	}
	if proxy.life != nil {
		proxy.life.ids = ids
		for i := range proxy.life.players {
			ids[uint32(uintptr(unsafe.Pointer(&proxy.life.players[i])))] = uint32(200 + i)
		}
	}
	var configurePlayers func(int)
	playersUnchanged := func() bool { return true }
	for _, sp := range specs {
		if sp.GuardEscort != nil {
			units, configure, unchanged, free := core.PortTestEscortPlayers()
			defer free()
			configurePlayers, playersUnchanged = configure, unchanged
			for i := range units {
				ids[uint32(uintptr(unsafe.Pointer(&units[i])))] = uint32(200 + i)
			}
			break
		}
	}
	scriptName, freeScriptName := alloc.Make([]byte{}, 128)
	defer freeScriptName()

	for id := byte(1); id < 34; id++ {
		ids[raw(id)] = uint32(id)
	}
	normalize := func(v uint32) uint32 {
		if id, ok := ids[v]; ok {
			return id
		}
		return v
	}
	put := func(b []byte, off int, v uint32) { binary.LittleEndian.PutUint32(b[off:], v) }
	get := func(b []byte, off int) uint32 { return binary.LittleEndian.Uint32(b[off:]) }
	guard := func(b []byte) {
		for i := 0; i < 8; i++ {
			b[i] = 0xa5
			b[len(b)-8+i] = 0x5a
		}
	}
	intact := func(b []byte) bool {
		for i := 0; i < 8; i++ {
			if b[i] != 0xa5 || b[len(b)-8+i] != 0x5a {
				return false
			}
		}
		return true
	}
	out := make([]PortTestRoamResult, 0, len(specs))
	for _, sp := range specs {
		clear(ob[8:780])
		clear(ub)
		clear(wb)
		clear(db)
		clear(hb)
		guard(hb)
		clear(tb[8:780])
		guard(db)
		guard(tb)
		proxy.trace = nil
		proxy.pathEndpoints = [2]*server.Waypoint{}
		proxy.endpointCall = 0
		guard(ob)
		guard(ub)
		for i := 0; i < 34; i++ {
			guard(wb[i*stride : (i+1)*stride])
		}
		obj.TypeInd = 1
		obj.ObjClass = object.ClassMonster
		obj.UpdateData = unsafe.Pointer(ud)
		ud.AIStackInd = sp.Stack
		ud.Field2 = 99
		ud.Field91 = uint32(sp.Index)
		for i := 0; i <= int(sp.Stack); i++ {
			ud.AIStack[i].Action = uint32(ai.ACTION_WAIT)
		}
		head := &ud.AIStack[sp.Stack]
		head.Action = uint32(ai.ACTION_ROAM)
		head.Args[0] = uintptr(raw(33))
		head.Args[2] = uintptr(sp.Mask)
		for i, id := range sp.History {
			put(ub, 8+300+4*i, raw(id))
		}
		for id := byte(1); id < 34; id++ {
			b := wb[int(id)*stride+8 : (int(id)+1)*stride-8]
			b[477] = sp.Flags[id]
			if sp.Enabled[id] {
				put(b, 480, 1)
			}
		}
		root := wb[8 : stride-8]
		root[476] = sp.Count
		for i, id := range sp.Neighbors {
			put(root, 92+i*8, raw(id))
		}
		if owner := sp.Owner; owner != nil {
			core.SetFrame(owner.Frame)
			core.SetTickRate(owner.FPS)
			core.PortTestAIEmptyMap()
			obj.PosVec = types.Pointf{X: math.Float32frombits(owner.X), Y: math.Float32frombits(owner.Y)}
			obj.Buffs = owner.Buffs
			obj.Frame134 = 1
			ud.Aggression = math.Float32frombits(owner.Aggression)
			ud.StatusFlags = object.MonsterStatus(owner.Status)
			ud.MonsterDef = def
			def.MoveSndFrameA100 = 250
			def.MoveSndFrameB104 = 251
			target.PosVec = types.Pointf{X: 321, Y: 654}
			if owner.Enemy {
				ud.CurrentEnemy = target
			}
			ud.Field2 = 0
			if owner.ExistingPath {
				ud.Field2 = 1
				ud.Field67 = 1
			}
			head.Args[0] = uintptr(raw(owner.Current))
			for id := byte(1); id < 34; id++ {
				wp := roamWaypoint(raw(id))
				wp.PosVec = types.Pointf{X: math.Float32frombits(owner.WX) + float32(id-1)*40, Y: math.Float32frombits(owner.WY)}
				wp.PointsCnt = sp.Count
				for j, n := range sp.Neighbors {
					wp.Points[j].Waypoint = roamWaypoint(raw(n))
				}
			}
			if owner.Register {
				core.Map.Nox_xxx_waypointMapRegister_5179B0(roamWaypoint(raw(1)))
			}
			proxy.mode = owner.PathMode
			proxy.fallback = roamWaypoint(raw(owner.Fallback))
		}
		if sp.GuardEscort != nil {
			portTestGuardEscortPrepare(proxy, obj, target, sp.GuardEscort)
			configurePlayers(sp.GuardEscort.Players)
			clear(scriptName)
			copy(scriptName, sp.GuardEscort.ScriptName)
			core.Objs.List = nil
			core.Objs.Pending = nil
			if sp.GuardEscort.ScriptName != "" {
				target.IDPtr = unsafe.Pointer(&scriptName[0])
				if sp.GuardEscort.Pending {
					core.Objs.Pending = target
				} else {
					core.Objs.List = target
				}
			}
		}
		if sp.Navigation != nil {
			portTestNavigationPrepare(proxy, obj, target, health, sp.Navigation)
		}
		if sp.Combat != nil {
			portTestCombatPrepare(proxy, obj, target, health, sp.Combat)
			configureWalls(sp.Combat.Wall)
		}
		if sp.Lifecycle != nil {
			portTestLifecyclePrepare(proxy, obj, target, health, sp.Lifecycle)
		}
		if sp.Path != nil {
			configureWalls(sp.Path.Wall)
			portTestPathPrepare(proxy, obj, sp.Path, raw)
		}
		core.Rand.Logic, core.Rand.Other = prand.New(sp.Seed), prand.New(sp.Seed+1)
		core.AI.StackChanged = false
		beforeO, beforeU, beforeW := bytes.Clone(ob), bytes.Clone(ub), bytes.Clone(wb)
		beforeD, beforeT := bytes.Clone(db), bytes.Clone(tb)
		beforeName := bytes.Clone(scriptName)
		beforeH := bytes.Clone(hb)
		ret := 0
		var nanos int64
		var combatResult *PortTestCombatResult
		var lifeResult *PortTestLifecycleResult
		switch sp.Op {
		case 11:
			ret = int(normalize(uint32(portTestLifecycleCall(obj, sp.Lifecycle))))
			lifeResult = portTestLifecycleTrace(proxy, health, normalize)
			combatResult = portTestCombatTrace(proxy, normalize)
			copy(beforeH[8:len(beforeH)-8], hb[8:len(hb)-8])
		case 10:
			portTestCombatCall(obj, sp.Combat)
			combatResult = portTestCombatTrace(proxy, normalize)
		case 0:
			server.GetAIAction(ai.ACTION_ROAM).Start(obj)
		case 1:
			server.GetAIAction(ai.ACTION_ROAM).Cancel(obj)
		case 2:
			roamInsert(ud, roamWaypoint(raw(sp.Insert)))
		case 3:
			ret = int(normalize(roamWaypointWord(roamPrevious(ud, sp.Mask))))
		case 4:
			ret = int(normalize(roamWaypointWord(roamSuccessor(ud, (*server.Waypoint)(unsafe.Pointer(&root[0])), sp.Mask))))
		case 5:
			ret = bool2int(roamDeadEnd(obj, (*server.Waypoint)(unsafe.Pointer(&root[0]))))
		case 6:
			start := time.Now()
			for repeat := 0; repeat < max(1, sp.Owner.Repeat); repeat++ {
				server.GetAIAction(ai.ACTION_ROAM).Update(obj)
			}
			nanos = time.Since(start).Nanoseconds()
		case 9:
			start := time.Now()
			for repeat := 0; repeat < max(1, sp.Owner.Repeat); repeat++ {
				ret = portTestPathCall(obj, sp.Path)
			}
			nanos = time.Since(start).Nanoseconds()
		case 8:
			ret = int(portTestNavigationCall(obj, sp.Navigation))
		case 7:
			start := time.Now()
			for repeat := 0; repeat < max(1, sp.Owner.Repeat); repeat++ {
				ret = int(normalize(portTestGuardEscortCall(obj, sp.GuardEscort)))
			}
			nanos = time.Since(start).Nanoseconds()
		default:
			panic("invalid roam operation")
		}
		if (sp.Navigation != nil || sp.Combat != nil) && !bytes.Equal(tb, beforeT) {
			for j := 8; j < len(tb)-8; j += 4 {
				if get(tb, j) != get(beforeT, j) {
					proxy.trace = append(proxy.trace, 0xff000000+uint32(j-8), get(tb, j))
				}
			}
			// Spatial iteration may write only its two visitation tokens.
			copy(beforeT[8+248:8+256], tb[8+248:8+256])
		}
		var pathChanges []uint32
		if sp.Path != nil {
			for id := 0; id < 34; id++ {
				for off := 504; off <= 512; off += 4 {
					j := id*stride + 8 + off
					if get(wb, j) != get(beforeW, j) {
						pathChanges = append(pathChanges, uint32(8192+id*stride+off), normalize(get(wb, j)))
						put(beforeW, j, get(wb, j))
					}
				}
			}
		}
		r := PortTestRoamResult{Lifecycle: lifeResult, Combat: combatResult, Nanos: nanos, Trace: proxy.trace, Index: ud.Field91, Arg: normalize(uint32(head.Args[0])), Field2: ud.Field2, Return: ret, Stack: ud.AIStackInd, Logic: core.Rand.Logic.Index(), Other: core.Rand.Other.Index(), Changed: core.AI.StackChanged, Intact: intact(ob) && intact(ub) && bytes.Equal(wb, beforeW) && bytes.Equal(db, beforeD) && bytes.Equal(tb, beforeT) && playersUnchanged() && bytes.Equal(scriptName, beforeName) && bytes.Equal(hb, beforeH)}
		if sp.Navigation != nil || sp.Path != nil {
			for _, off := range offsets {
				r.Trace = append(r.Trace, uint32(off), normalize(*memmap.PtrUint32(0x5D4594, off)))
			}
		}
		if lifeResult != nil {
			r.Intact = r.Intact && proxy.life.playersUnchanged()
		}
		if combatResult != nil {
			r.Intact = r.Intact && combatResult.Intact && wallsUnchanged()
		}
		if sp.Path != nil {
			r.Changes = append(r.Changes, pathChanges...)
			r.Trace = append(r.Trace, portTestPathGraphState(normalize)...)
			r.Intact = r.Intact && wallsUnchanged()
		}
		for i := range r.History {
			r.History[i] = byte(normalize(get(ub, 8+300+4*i)))
		}
		for region, pair := range [][2][]byte{{beforeO, ob}, {beforeU, ub}} {
			for j := 8; j < len(pair[0])-8; j += 4 {
				a, b := get(pair[0], j), get(pair[1], j)
				if a != b {
					r.Changes = append(r.Changes, uint32(region*4096+j-8), normalize(b))
				}
			}
		}
		out = append(out, r)
	}
	return out
}
