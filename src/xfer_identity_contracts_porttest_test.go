//go:build porttest

package opennox

import (
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestXferIdentityLifetime(t *testing.T) {
	typ, free := alloc.New(server.ObjectType{})
	defer free()
	obj, free2 := alloc.New(server.Object{})
	defer free2()
	names := legacy.PortTestXferRegistryNames()
	if len(names) != 28 {
		t.Fatal("transfer registration count", len(names))
	}
	keys := make(map[string]unsafe.Pointer, len(names))
	seen := make(map[unsafe.Pointer]bool, len(names))
	for _, name := range names {
		key := legacy.PortTestXferRegistryPointer(name)
		if key == nil || seen[key] {
			t.Fatal("transfer identity", name)
		}
		seen[key] = true
		keys[name] = key
		typ.Xfer = key
		obj.Xfer = typ.Xfer
		collisionRegistryGrow(128)
		if typ.Xfer != key || obj.Xfer != key || legacy.PortTestXferRegistryPointer(name) != key {
			t.Fatal("stored transfer identity", name)
		}
	}
	if server.DefaultXfer != keys["DefaultXfer"] {
		t.Fatal("default transfer key")
	}
	for _, row := range []struct {
		name string
		get  func() unsafe.Pointer
	}{
		{"DefaultXfer", legacy.Get_nox_xxx_XFerDefault_4F49A0},
		{"FieldGuideXfer", legacy.Get_nox_xxx_XFerFieldGuide_4F6390},
		{"AbilityRewardXfer", legacy.Get_nox_xxx_XFerAbilityReward_4F6240},
		{"InvisibleLightXfer", legacy.Get_nox_xxx_XFerInvLight_4F5AA0},
	} {
		if row.get() != keys[row.name] {
			t.Fatal("transfer getter", row.name)
		}
	}
}

// The startup getter also controls per-object script storage allocation.
func TestXferInvisibleLightClassification(t *testing.T) {
	s := newObjectXferOwner(t)
	s.PortTestObjectXferAdmission(0, 0, true)
	typ := s.Types.ByInd(1)
	light := legacy.Get_nox_xxx_XFerInvLight_4F5AA0()
	s.Objs.XFerInvLight = light
	for _, flags := range []noxflags.GameFlag{0, noxflags.GameFlag22, noxflags.GameFlag23, noxflags.GameFlag22 | noxflags.GameFlag23} {
		for _, key := range []unsafe.Pointer{nil, server.DefaultXfer, light} {
			restore := noxflags.PortTestGameFlags(flags)
			typ.Xfer = key
			u := newObjectXferSimple(t, s)
			restore()
			if u.Xfer != key || (u.Field189 != nil) != (flags != 0 && key == light) {
				t.Fatal("invisible-light allocation", flags, key, u.Xfer, u.Field189)
			}
		}
	}
}
