//go:build porttest

package opennox

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestClientTooltipRequestContract(t *testing.T) {
	o := newTooltipOwner(t)
	for _, tc := range []struct {
		sub, sentinel uint32
		kind          byte
	}{{1, 137, 1}, {2, 41, 2}, {4, 6, 4}, {7, 137, 1}, {6, 41, 2}} {
		o.prepare(0x400100, tc.sub, 0, 0)
		o.dr.NetCode32 = 0x1234
		first := o.invoke(t, 0, 0, o.dr)
		want := [][]byte{{0xe2, 0x34, 0x92, tc.kind}}
		if !reflect.DeepEqual(first.Messages, want) {
			t.Fatalf("sub%d request %x want%x", tc.sub, first.Messages, want)
		}
		if got := unsafe.Slice((*uint32)(unsafe.Pointer(o.dr)), 128)[108]; got != tc.sentinel {
			t.Fatalf("cached metadata %d want%d", got, tc.sentinel)
		}
		second := o.invoke(t, 0, 1, o.dr)
		if !reflect.DeepEqual(first.Messages, second.Messages) || len(second.Text) != 0 {
			t.Fatal("pending tooltip repeated a request or produced a title")
		}
	}
}
func TestClientTooltipLocalizedOrderContract(t *testing.T) {
	o := newTooltipOwner(t)
	for _, tc := range []struct {
		lang int
		want string
	}{{0, "M0 M1 Weapon M2 S3"}, {2, "Weapon M2 S3M1  M0"}, {3, "M0 Weapon M1 M2 S3"}, {5, "Weapon M1 M0 M2 S3"}, {6, "M0 M2 S3  M1Weapon"}} {
		o.language(tc.lang)
		o.prepare(0x1000, 0, 0, 15)
		r := o.invoke(t, 0, 0, o.dr)
		want := []uint16{}
		for _, c := range tc.want {
			want = append(want, uint16(c))
		}
		if !reflect.DeepEqual(r.Text, want) {
			t.Fatalf("language%d text%x want%x", tc.lang, r.Text, want)
		}
	}
}
func TestClientTooltipBorrowedNameContract(t *testing.T) {
	o := newTooltipOwner(t)
	o.prepare(0, 0, 0, 0)
	p := o.c.Things.TypeByInd(4).PrettyName
	if got := legacy.PortTestTooltip(o.dr); got != p {
		t.Fatal("ordinary item name must remain borrowed")
	}
	if o.scratch[0] != 0 || o.scratch[1] != 0x8101 {
		t.Fatal("fallback scratch initialization overwrote tail")
	}
	if got := legacy.PortTestTooltip(nil); got != &o.scratch[0] {
		t.Fatal("nil item must return scratch")
	}
}
func TestClientTooltipCursorBoundsContract(t *testing.T) {
	o := newTooltipOwner(t)
	text := make([]uint16, 257)
	for i := range text {
		text[i] = uint16(0xd800 + i)
	}
	p := tooltipWide(t, text)
	legacy.PortTestTooltipCursor(&p[0])
	if !reflect.DeepEqual(o.cursor[:255], p[:255]) || o.cursor[255] != 0 {
		t.Fatal("cursor truncation changed raw units or failed to terminate")
	}
	old := append([]uint16(nil), o.cursor...)
	legacy.PortTestTooltipCursor(nil)
	if o.cursor[0] != 0 || !reflect.DeepEqual(o.cursor[1:], old[1:]) {
		t.Fatal("nil cursor setter changed the tail")
	}
}

func TestClientTooltipBookOrderContract(t *testing.T) {
	o := newTooltipOwner(t)
	for _, tc := range []struct {
		lang int
		sub  uint32
		want string
	}{{0, 1, "BookOf Spark"}, {6, 1, "Spark BookOf"}, {0, 2, "Guide1 LoreScroll"}, {3, 2, "LoreScroll Guide1"}, {5, 2, "LoreScroll Guide1"}, {0, 4, "BookOf Berserk"}, {6, 4, "Berserk BookOf"}} {
		o.language(tc.lang)
		o.prepare(0x100, tc.sub, 1, 0)
		got := alloc.GoString16(legacy.PortTestTooltip(o.dr))
		if got != tc.want {
			t.Fatalf("language%d subtype%d got%q want%q", tc.lang, tc.sub, got, tc.want)
		}
	}
}
func TestClientTooltipMissingDefinitionContract(t *testing.T) {
	o := newTooltipOwner(t)
	o.prepare(0x1000, 0, 0, 15)
	o.weapon.TypeInd = 5
	got := alloc.GoString16(legacy.PortTestTooltip(o.dr))
	want := "Missing: " + alloc.GoString(o.c.Things.TypeByInd(4).Name)
	if got != want {
		t.Fatalf("missing definition got%q want%q", got, want)
	}
}
func TestClientTooltipLongNameContract(t *testing.T) {
	o := newTooltipOwner(t)
	old := o.weapon.Desc8
	defer func() { o.weapon.Desc8 = old }()
	for _, n := range []int{0, 1, 255, 511, 1000, 1023} {
		text := make([]uint16, n)
		for i := range text {
			text[i] = uint16(0xd800 + i%0x700)
		}
		p := tooltipWide(t, text)
		o.weapon.Desc8 = &p[0]
		for _, lang := range []int{0, 1, 2, 3, 4, 5, 6} {
			o.language(lang)
			o.prepare(0x1000, 0, 0, 0)
			got := legacy.PortTestTooltip(o.dr)
			if got != &o.scratch[0] || !reflect.DeepEqual(o.scratch[:n], text) || o.scratch[n] != 0 {
				t.Fatalf("length%d language%d raw name changed", n, lang)
			}
			for i := n + 1; i < len(o.scratch); i++ {
				if o.scratch[i] != uint16(0x8100+i%256) {
					t.Fatal("name copy changed unused tail")
				}
			}
		}
	}
}
