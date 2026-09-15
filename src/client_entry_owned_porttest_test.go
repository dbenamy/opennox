//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
)

func TestClientEntryOwnedDataContract(t *testing.T) {
	o := newEntryOwner(t)
	for _, lang := range []int{0, 6, 8} {
		o.create(t, lang, 8, nil)
		ptr := o.win.WidgetData
		if !alloc.PortTestAllocationLive(ptr) {
			t.Fatal("Go entry data is not registered as owned")
		}
		o.win.Destroy()
		if !alloc.PortTestAllocationLive(ptr) {
			t.Fatal("entry data released before deferred cleanup")
		}
		o.c.GUI.FreeDestroyed()
		o.c.GUI.FreeDestroyed()
		o.win = nil
		if alloc.PortTestAllocationLive(ptr) {
			t.Fatal("entry data survived cleanup")
		}
	}
}
