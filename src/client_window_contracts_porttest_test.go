//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestClientWindowNilOutputsContract(t *testing.T) {
	for _, op := range []int{0, 1, 2} {
		ret, xy := legacy.PortTestWindowHelper(op, nil, 0, 0, nil, 0)
		if ret != -2 {
			t.Fatalf("op%d nil return=%d", op, ret)
		}
		want := [2]uint32{}
		if op == 1 {
			want = [2]uint32{0x13572468, 0x89abcdef}
		}
		if xy != want {
			t.Fatalf("op%d output=%x want%x", op, xy, want)
		}
	}
}
func TestClientWindowResizeContract(t *testing.T) {
	o := newWindowHelperOwner(t)
	o.create()
	result := o.invoke(0, 0, 4, o.win, 31, -7, nil, 0)
	if result.Return != 0 || len(result.Notices) != 1 {
		t.Fatalf("resize result=%+v", result)
	}
	notice := result.Notices[0]
	if notice.Code != 16388 || notice.A != 31 || int32(notice.B) != -7 || notice.Geometry != [6]int32{10, 12, 31, -7, 41, 5} {
		t.Fatalf("resize callback did not observe updated geometry: %+v", notice)
	}
}
func TestClientWindowBoundsAndAncestryContract(t *testing.T) {
	o := newWindowHelperOwner(t)
	o.create()
	for _, v := range [][3]int{{13, 16, 1}, {73, 76, 1}, {74, 76, 0}, {13, 77, 0}, {12, 16, 0}} {
		got, _ := legacy.PortTestWindowHelper(3, o.win, v[0], v[1], nil, 0)
		if got != v[2] {
			t.Fatalf("point %v: got%d", v, got)
		}
	}
	for _, v := range [][3]int{{1, 1, 0}, {1, 4, 1}, {4, 1, 0}, {2, 5, 0}, {0, 5, 1}} {
		got, _ := legacy.PortTestWindowHelper(9, o.windows[v[0]], 0, 0, o.windows[v[1]].C(), 0)
		if got != v[2] {
			t.Fatalf("ancestry %v: got%d", v, got)
		}
	}
}

func TestClientWindowMaximumIDRangeContract(t *testing.T) {
	o := newWindowHelperOwner(t)
	o.create()
	child := o.windows[4]
	child.SetID(0x7fffffff)
	legacy.PortTestWindowHelper(22, o.win, 0x7fffffff, 0x7fffffff, nil, 1)
	if child.Flags&16 == 0 {
		t.Fatal("inclusive maximum-ID range did not hide child")
	}
	legacy.PortTestWindowHelper(23, o.win, 0x7fffffff, 0x7fffffff, nil, 0)
	if child.Flags&8 != 0 {
		t.Fatal("inclusive maximum-ID range did not disable child")
	}
	legacy.PortTestWindowHelper(22, o.win, 0x7fffffff, 0x7fffffff, nil, 0)
	if child.Flags&16 != 0 {
		t.Fatal("inclusive maximum-ID range did not show child")
	}
}
