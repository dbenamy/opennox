//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestPrefabScriptsCallbackAssignment(t *testing.T) {
	s := newObjectXferOwner(t)
	s.NoxScriptVM.Init(s)
	funcs := []prefabScriptFunction{{Name: "GLOBAL", Code: []uint32{72}}, {Name: "GLOBAL", Vars: []uint32{1, 1, 1, 1}, Code: []uint32{72}}, {Name: "Target", Code: []uint32{72}}}
	if err := s.NoxScriptVM.ReadScript(bytes.NewReader(prefabScriptsEncode(nil, funcs))); err != nil {
		t.Fatal(err)
	}
	u := newObjectXferSimple(t, s)
	sd := u.Field189
	if sd == nil {
		t.Fatal("real editor factory omitted script storage")
	}
	script := unsafe.Slice((*byte)(sd), 2572)
	update, freeUpdate := alloc.Make([]byte{}, 4096)
	defer freeUpdate()
	collide, freeCollide := alloc.Make([]byte{}, 88)
	defer freeCollide()
	oldClass := u.ObjClass
	u.UpdateData = unsafe.Pointer(&update[0])
	u.CollideData = unsafe.Pointer(&collide[0])
	defer func() { u.ObjClass = oldClass; u.Field189 = sd; u.UpdateData = nil; u.CollideData = nil }()
	type slot struct {
		script, update, collide int
		pickup                  bool
	}
	target := func(class object.Class, event int) (slot, bool) {
		if event == 14 {
			return slot{script: 0, pickup: true}, true
		}
		switch {
		case class.Has(object.ClassTrigger):
			v, ok := map[int]slot{0: {script: 512, update: 16}, 1: {script: 256, update: 24}, 2: {script: 384, update: 32}}[event]
			return v, ok
		case class.Has(object.ClassMonster):
			v, ok := map[int]slot{3: {script: 640, update: 1236}, 4: {script: 768, update: 1228}, 5: {script: 896, update: 1268}, 6: {script: 1024, update: 1244}, 7: {script: 1152, update: 1252}, 8: {script: 1280, update: 1260}, 9: {script: 1408, update: 1276}, 10: {script: 1536, update: 1284}, 11: {script: 1664, update: 1292}, 13: {script: 1792, update: 1300}}[event]
			return v, ok
		case class.Has(object.ClassHole):
			return slot{script: 128, collide: 4}, event == 12
		case class.Has(object.ClassMonsterGenerator):
			v, ok := map[int]slot{15: {script: 1920, update: 52}, 16: {script: 2048, update: 60}, 17: {script: 2304, update: 68}, 18: {script: 2176, update: 76}}[event]
			return v, ok
		}
		return slot{}, false
	}
	type row struct {
		Name                    string
		Script, Update, Collide [32]byte
		Pickup                  int32
	}
	var rows []row
	for _, class := range []object.Class{object.ClassSimple, object.ClassTrigger, object.ClassMonster, object.ClassHole, object.ClassMonsterGenerator, object.ClassTrigger | object.ClassMonster} {
		u.ObjClass = class
		for _, mode := range []noxflags.GameFlag{0, 0x200000, 0x400000} {
			restore := noxflags.PortTestGameFlags(mode)
			for _, present := range []bool{false, true} {
				for _, name := range []string{"Target", "Missing"} {
					cp, free := alloc.CString(name)
					for event := -1; event <= 19; event++ {
						for _, b := range [][]byte{script, update, collide} {
							for i := range b {
								b[i] = 0x5a
							}
						}
						u.ScriptPickup.Func = 12345
						u.Field189 = sd
						if !present {
							u.Field189 = nil
						}
						ws, wu, wc := append([]byte(nil), script...), append([]byte(nil), update...), append([]byte(nil), collide...)
						pickup := int32(12345)
						where, ok := target(class, event)
						if present && ok {
							if mode != 0 {
								copy(ws[where.script:], append([]byte(name), 0))
							} else {
								index := int32(-1)
								if name == "Target" {
									index = 2
								}
								switch {
								case where.pickup:
									pickup = index
								case where.collide != 0:
									binary.LittleEndian.PutUint32(wc[where.collide:], uint32(index))
								default:
									binary.LittleEndian.PutUint32(wu[where.update:], uint32(index))
								}
							}
						}
						legacy.PortTestPrefabScriptsCall(14, u.CObj(), unsafe.Pointer(cp), nil, uint32(int32(event)))
						label := fmt.Sprintf("class%x/mode%x/present%t/%s/event%d", uint32(class), uint32(mode), present, name, event)
						if !bytes.Equal(script, ws) || !bytes.Equal(update, wu) || !bytes.Equal(collide, wc) || u.ScriptPickup.Func != pickup {
							t.Errorf("%s changed unexpected callback fields", label)
						}
						if present && ok && (mode != 0 || name == "Target") {
							got, found := s.NoxScriptVM.Nox_script_objCallbackName_508CB0(u, event)
							if !found || got != name {
								t.Errorf("%s real callback getter returned %q/%t", label, got, found)
							}
						}
						rows = append(rows, row{label, sha256.Sum256(script), sha256.Sum256(update), sha256.Sum256(collide), u.ScriptPickup.Func})
					}
					free()
				}
			}
			restore()
		}
	}
	spellbookCapture(t, "prefab-scripts-callbacks", rows, "8b4d6b80aaac00b0fe268422e78fb195387e265d61cec1110e83e9355d5f8b85")
}
