//go:build porttest

package opennox

import (
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func (o *meterOwner) weapons(t *testing.T) (*client.Drawable, *client.Drawable) {
	names := []string{"meter-unused-zero"}
	for i := 1; o.c.Things.TypeByInd(i) != nil; i++ {
		names = append(names, o.c.Things.TypeByInd(i).ID())
	}
	t.Cleanup(o.c.srv.Server.PortTestRewardTypes(names, nil, true, 0, 0))
	for i := 1; i < len(names); i++ {
		if o.c.srv.Types.IndByID(names[i]) != i {
			t.Fatal("client/server item type indices differ")
		}
	}
	bowType, quiverType := o.c.Things.TypeByID("Bow"), o.c.Things.TypeByID("Quiver")
	bowType.ObjClass, quiverType.ObjClass = 0x1000000, 0x1000000
	bowType.ObjSubClass = 4
	o.weapon.TypeInd = uint32(bowType.Index())
	name, freeName := alloc.CString16("Bow")
	t.Cleanup(freeName)
	o.weapon.Desc8 = name
	q, freeQ := alloc.New(server.Modifier{})
	t.Cleanup(freeQ)
	name, freeName = alloc.CString16("Quiver")
	t.Cleanup(freeName)
	q.TypeInd, q.Desc8 = uint32(quiverType.Index()), name
	o.weapon.Next80 = q
	t.Cleanup(func() { o.weapon.Next80 = nil })
	o.c.srv.Types.ClientTypeByID = func(id string) int {
		if typ := o.c.Things.TypeByID(id); typ != nil {
			return typ.Index()
		}
		return 0
	}
	o.c.srv.Server.Nox_xxx_equipWeapon_4157C0()
	bow, quiver := o.newUnlinked(t), o.newUnlinked(t)
	o.c.DrawableLinkThing(bow, bowType.Index())
	o.c.DrawableLinkThing(quiver, quiverType.Index())
	return bow, quiver
}
func (o *meterOwner) equip(bow, quiver *client.Drawable, health uint16) {
	o.equipment[7] = uint32(uintptr(bow.C()))
	if quiver != nil {
		o.equipment[0] = uint32(uintptr(quiver.C()))
	}
	o.slot(0, bow, 1, 123)
	bow.NetCode32 = 123
	*(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 4)) = 4
	h := unsafe.Slice((*uint16)(unsafe.Add(bow.C(), 292)), 2)
	h[0], h[1] = health, 100
	*(*noxrender.Viewport)(memmap.PtrOff(0x5D4594, 1049732)) = *o.c.Viewport()
}
func TestClientMetersWeaponTooltipContract(t *testing.T) {
	o := newMeterOwner(t)
	bow, quiver := o.weapons(t)
	o.plain(t)
	o.equip(bow, quiver, 100)
	legacy.PortTestMeterCall(33, nil, 0, 0, 0, 0)
	if got := alloc.GoString16((*uint16)(memmap.PtrOff(0x5D4594, 1096676))); got != "Bow\nQuiver" {
		t.Fatalf("equipped tooltip %q", got)
	}
	clear(o.equipment)
	legacy.PortTestMeterCall(33, nil, 0, 0, 0, 0)
	if got := alloc.GoString16((*uint16)(memmap.PtrOff(0x5D4594, 1096676))); got != "Current weapon" {
		t.Fatalf("empty weapon tooltip %q", got)
	}
}
func TestClientMetersWeaponMatrix(t *testing.T) {
	o := newMeterOwner(t)
	bow, quiver := o.weapons(t)
	var out []meterResult
	id := 0
	for mode := 0; mode < 3; mode++ {
		for _, health := range []uint16{0, 24, 25, 49, 50, 100} {
			for _, charge := range [][2]uint32{{0, 0}, {0, 100}, {1, 100}, {99, 100}, {100, 100}, {101, 100}} {
				o.plain(t)
				if mode > 0 {
					q := quiver
					if mode == 1 {
						q = nil
					}
					o.equip(bow, q, health)
				}
				r := &o.meters.Records[4]
				r.Current, r.Maximum = charge[0], charge[1]
				out = append(out, o.invoke(t, id, 0, 17, r.Window, 0, 0, 0, 0))
				out = append(out, o.invoke(t, id, 1, 33, nil, 0, 0, 0, 0))
				id++
			}
		}
	}
	meterCapture(t, "meters-weapons", out, len(out), "1c74529d57bfaa0eb6020679a052bac6b3107e9d3d6a5d0583a659992ad1a36d")
}
