//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestWorldMotionTriggerDisabledScript(t *testing.T) {
	o := newCollisionCoreOwner(t)
	u, v := newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s)
	savedU, savedV := *u, *v
	data := collisionCoreGuarded(t, o, 56)
	t.Cleanup(func() { *u = savedU; *v = savedV })
	_, accept, restore := o.s.NoxScriptVM.PortTestWorldMotionPredicates()
	t.Cleanup(restore)
	u.UpdateData = data
	u.ObjFlags = 0x1000000
	v.Mass = 2
	words := (*[14]uint32)(data)
	*words = [14]uint32{}
	cb := (*server.ScriptCallback)(unsafe.Add(data, 12))
	cb.Func = accept
	cb.Flags = 2
	legacy.PortTestWorldMotionTimed("trigger", u, v, 0)
	if cb.Flags != 3 || words[0]&1 == 0 || words[1] != uint32(uintptr(v.CObj())) {
		t.Fatal("one-shot trigger first contact")
	}
	before := *words
	t.Log("one-shot script accepted first contact and disabled itself; second contact must tolerate no script result")
	legacy.PortTestWorldMotionTimed("trigger", u, v, 0)
	if *words != before {
		t.Fatal("disabled trigger callback changed contact state")
	}
	// A script disabled before its first contact should leave the trigger unarmed.
	words[0] = 0
	words[1] = 0
	cb.Flags = 1
	legacy.PortTestWorldMotionTimed("trigger", u, v, 0)
	if words[0] != 0 || words[1] != 0 {
		t.Fatal("disabled predicate admitted contact")
	}
}
