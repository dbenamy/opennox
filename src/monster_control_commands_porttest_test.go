//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
	"unsafe"
)

func TestMonsterControlCommands(t *testing.T) {
	o := newCollisionCoreOwner(t)
	serverConfigOwnBytes(t, 0x5D4594, 2487684, 4)
	a, b := newCreatureXferObject(t, o.s, "Monster"), newCreatureXferObject(t, o.s, "Monster")
	sa, sb := *a, *b
	ud := a.UpdateDataMonster()
	savedUD := *ud
	t.Cleanup(func() { *a = sa; *b = sb; *ud = savedUD })
	def := (*server.MonsterDef)(collisionCoreGuarded(t, o, int(unsafe.Sizeof(server.MonsterDef{}))))
	oldChanged := o.s.AI.StackChanged
	t.Cleanup(func() { o.s.AI.StackChanged = oldChanged })
	ids := collisionCoreIDs(a, b)
	normalize := func(raw uint32) uint32 {
		if raw == 0 {
			return 0
		}
		if raw == uint32(uintptr(a.UpdateData)) {
			return 2001
		}
		if id, ok := ids[unsafe.Pointer(uintptr(raw))]; ok {
			return id
		}
		start := uint32(uintptr(unsafe.Pointer(&ud.AIStack[0])))
		end := start + uint32(unsafe.Sizeof(ud.AIStack))
		if raw >= start && raw < end {
			if (raw-start)%24 != 0 {
				t.Fatal("unaligned returned action")
			}
			return 3000 + (raw-start)/24
		}
		return raw
	}
	type row struct {
		Op           string
		Class, Flags uint32
		InitialStack int8
		Variant      int
		Return       uint32
		Stack        int8
		Actions      [24][6]uint32
		Fields       [12]uint32
		Changed      bool
	}
	var rows []row
	ops := []string{"ensure", "look", "walk", "patrol", "hunt", "idle", "follow", "melee", "missile", "control-byte", "fight", "flee", "wait"}
	for _, op := range ops {
		for _, class := range []object.Class{object.ClassMonster, object.ClassSimple} {
			for _, flags := range []object.Flags{4, 4 | 0x20, 4 | 0x8000} {
				for _, stack := range []int8{0, 1, 22, 23} {
					for variant := 0; variant < 3; variant++ {
						*a = sa
						*b = sb
						*ud = server.MonsterUpdateData{}
						*def = server.MonsterDef{MeleeAttackRange112: 32, MissileAttackRange212: 128}
						def.MissileName148[0] = 'B'
						worldGeometryResetObject(a, 1001, 100.125, 200.5, false)
						worldGeometryResetObject(b, 1002, 133.25, 250.75, false)
						a.ObjClass = class
						a.ObjFlags = flags
						a.ObjSubClass = 0
						ud.MonsterDef = def
						ud.AIStackInd = stack
						for i := 0; i <= int(stack); i++ {
							ud.AIStack[i] = server.AIStackItem{Action: 1, Args: [4]uintptr{11, 22, 33, 44}}
						}
						if stack == 0 {
							ud.AIStack[0].Action = 0
						}
						ud.Field333 = 0xaabbccdd
						ud.SightRange = 77
						ud.Field304 = 0
						ud.Field2 = 17
						ud.Field67 = 19
						ud.Field120_1 = 3
						ud.Field120_2 = 4
						ud.Field120_3 = 5
						if variant == 1 {
							def.MeleeAttackRange112 = 0
							def.MissileName148[0] = 0
						}
						if variant == 2 {
							ud.AIStack[stack].Action = 31
						}
						o.s.SetFrame(0xfffffff0)
						o.s.AI.StackChanged = false
						*memmap.PtrUint32(0x5D4594, 2487684) = 0x12345678
						args := [8]uint32{math.Float32bits(3.25), math.Float32bits(-5.5), math.Float32bits(23.75), math.Float32bits(48.125), math.Float32bits(31.5)}
						arg := []int32{0, 8, 255}[variant]
						if op == "ensure" {
							arg = []int32{1, 17, 30}[variant]
						}
						if op == "flee" {
							args[0] = uint32(uintptr(b.CObj()))
							args[1] = []uint32{0, 32, 0xffffffff}[variant]
						}
						if op == "wait" {
							arg = []int32{0, 32, -1}[variant]
						}
						if op == "control-byte" {
							args[0] = []uint32{0, 0x12345680, 0xffffffff}[variant]
						}
						beforeArgs, beforeUD := args, *ud
						rv := legacy.PortTestMonsterControl(op, a, b, &args, arg)
						if args != beforeArgs {
							t.Fatal("command changed argument record", op)
						}
						if class != object.ClassMonster && *ud != beforeUD {
							t.Fatal("command modified nonmonster", op)
						}
						if flags&0x8000 != 0 && op != "control-byte" && op != "ensure" && *ud != beforeUD {
							t.Fatal("command modified dead monster", op)
						}
						if class == object.ClassMonster && op == "control-byte" && ud.Field333 != (0xaabbcc00|args[0]&255) {
							t.Fatal("control byte overwrote adjacent bytes")
						}
						if class == object.ClassMonster && flags&0x8000 == 0 && op == "walk" && (ud.AIStackInd != 1 || ud.AIStack[0].Action != 32 || ud.AIStack[0].Args[0] != 8 || ud.AIStack[1].Action != 8 || uint32(ud.AIStack[1].Args[0]) != args[0] || uint32(ud.AIStack[1].Args[1]) != args[1]) {
							t.Fatal("walk stack contract", stack, variant, ud.AIStackInd)
						}
						r := row{Op: op, Class: uint32(class), Flags: uint32(flags), InitialStack: stack, Variant: variant, Return: normalize(rv), Stack: ud.AIStackInd, Changed: o.s.AI.StackChanged}
						for i, s := range ud.AIStack {
							r.Actions[i] = [6]uint32{s.Action, uint32(s.Args[0]), uint32(s.Args[1]), normalize(uint32(s.Args[2])), uint32(s.Args[3]), s.Field5}
						}
						r.Fields = [12]uint32{ud.Field2, ud.Field67, uint32(ud.Field120_1), uint32(ud.Field120_2), uint32(ud.Field120_3), ud.Field124, ud.Field137, math.Float32bits(ud.SightRange), ud.Field333, normalize(ud.Field304), o.s.Frame(), memmap.Uint32(0x5D4594, 2487684)}
						rows = append(rows, r)
					}
				}
			}
		}
	}
	spellbookCapture(t, "monster-control-commands", rows, "835364ffeba09032be9d1bdd62b507b838cc66314ae79cc6b61fa2848a8e4923")
}

// Head reads include the legacy slot immediately before an empty stack. Keep its
// return explicit rather than silently replacing it with the native nil-head API.
func TestMonsterControlHead(t *testing.T) {
	o := newCollisionCoreOwner(t)
	a := newCreatureXferObject(t, o.s, "Monster")
	ud := a.UpdateDataMonster()
	saved := *ud
	t.Cleanup(func() { *ud = saved })
	type row struct {
		Stack  int8
		Value  uint32
		Return uint32
	}
	var rows []row
	for _, index := range []int8{-1, 0, 1, 22, 23} {
		for _, value := range []uint32{0, 1, 30, 39, 40, 71, 0x80000000, 0xffffffff} {
			*ud = server.MonsterUpdateData{}
			ud.AIStackInd = index
			word := (*uint32)(unsafe.Add(a.UpdateData, 552+24*int(index)))
			*word = value
			var args [8]uint32
			rv := legacy.PortTestMonsterControl("head", a, nil, &args, 0)
			if rv != value {
				t.Fatal("action head raw word", index, value, rv)
			}
			rows = append(rows, row{index, value, rv})
		}
	}
	spellbookCapture(t, "monster-control-head", rows, "d89afaad3e10ceaeb91b157561516975574cfe022b183328e37f911d9e73303e")
}

func TestMonsterControlCommandEdges(t *testing.T) {
	o := newCollisionCoreOwner(t)
	a := newCreatureXferObject(t, o.s, "Monster")
	sa := *a
	ud := a.UpdateDataMonster()
	saved := *ud
	t.Cleanup(func() { *a = sa; *ud = saved })
	type row struct {
		Op      string
		Mode    int
		Return  uint32
		Stack   int8
		Actions [24][6]uint32
	}
	var rows []row
	for _, op := range []string{"hunt", "idle", "follow", "melee", "missile", "control-byte", "flee", "wait"} {
		rv := legacy.PortTestMonsterControl(op, nil, nil, nil, 0)
		if rv != 0 {
			t.Fatal("nil unit admission", op, rv)
		}
		rows = append(rows, row{Op: op, Mode: 0, Return: rv})
	}
	for _, op := range []string{"follow", "fight", "flee", "missile"} {
		for mode := 1; mode <= 2; mode++ {
			*a = sa
			*ud = server.MonsterUpdateData{}
			a.ObjClass = object.ClassMonster
			a.ObjFlags = 4
			ud.AIStackInd = 1
			ud.AIStack[0] = server.AIStackItem{Action: 1}
			ud.AIStack[1] = server.AIStackItem{Action: 24, Args: [4]uintptr{11, 22, 33, 44}}
			before := *ud
			target := a
			if mode == 1 {
				target = nil
			}
			var args [8]uint32
			var point *[8]uint32
			if op == "flee" {
				args[0] = uint32(uintptr(a.CObj()))
				args[1] = 17
				point = &args
			}
			// Null point is an explicit missile no-op; an existing flee action is an
			// explicit flee no-op. Follow/fight reject nil and self targets.
			rv := legacy.PortTestMonsterControl(op, a, target, point, 0)
			if *ud != before {
				t.Fatal("command edge changed existing stack", op, mode)
			}
			r := row{Op: op, Mode: mode, Return: rv, Stack: ud.AIStackInd}
			for i, s := range ud.AIStack {
				r.Actions[i] = [6]uint32{s.Action, uint32(s.Args[0]), uint32(s.Args[1]), uint32(s.Args[2]), uint32(s.Args[3]), s.Field5}
			}
			rows = append(rows, r)
		}
	}
	spellbookCapture(t, "monster-control-command-edges", rows, "e2cad238996bc9106f1396f053e6e4365de4d32226360b97e26c9261c3620078")
}
