//go:build porttest

package opennox

import (
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func (o *meterOwner) constructorStart(t *testing.T, class int, flags uint32) {
	o.plain(t)
	o.c.GUI.DestroyAll()
	o.c.GUI.FreeDestroyed()
	o.c.GUI.FreeDestroyed()
	clear(o.meters.Records)
	o.parent = nil
	o.named = nil
	*o.meters.NamedWord("dword_5d4594_1090276") = 0
	noxflags.ResetGame()
	noxflags.SetGame(noxflags.GameFlag(flags))
	if class < 0 {
		*o.meters.NamedWord("dword_8531A0_2576") = 0
	} else {
		*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 2251)) = byte(class)
	}
}

func TestClientMetersConstructorContract(t *testing.T) {
	o := newMeterOwner(t)
	for class := -1; class <= 2; class++ {
		o.constructorStart(t, class, 0)
		got := legacy.PortTestMeterCall(35, nil, 0, 0, 0, 0)
		o.collect()
		if class < 0 {
			if got != 0 || len(o.windows) != 0 {
				t.Fatal("missing player created meter windows")
			}
			continue
		}
		if got != 1 || len(o.named) != 12 {
			t.Fatalf("class %d: constructor return/images %d/%v", class, got, o.named)
		}
		if o.meters.Records[0].Window == nil || o.meters.Records[2].Window == nil {
			t.Fatal("missing health windows")
		}
		if (o.meters.Records[1].Window != nil) != (class != 0) || (o.meters.Records[3].Window != nil) != (class != 0) {
			t.Fatal("class-specific mana windows")
		}
		want := 6
		if class != 0 {
			want = 9
		}
		if len(o.windows) != want {
			t.Fatalf("class %d: %d windows, want %d", class, len(o.windows), want)
		}
	}
}

func TestClientMetersConstructorMatrix(t *testing.T) {
	o := newMeterOwner(t)
	var out []meterResult
	id := 0
	for class := -1; class <= 2; class++ {
		for _, flags := range []uint32{0, 4096} {
			o.constructorStart(t, class, flags)
			out = append(out, o.invoke(t, id, 0, 35, nil, 0, 0, 0, 0))
			if class >= 0 {
				out = append(out, o.invoke(t, id, 1, 31, nil, 0, 0, 0, 0))
				out = append(out, o.invoke(t, id, 2, 32, nil, 0, 0, 0, 0))
				legacy.PortTestMeterCall(4, nil, 50, 100, 0, 0)
				legacy.PortTestMeterCall(8, nil, 25, 100, 0, 0)
				*o.meters.NamedWord("nox_client_renderBubbles_80844") = 0
				o.c.GUI.Draw()
				out = append(out, o.invoke(t, id, 3, 0, nil, 0, 0, 0, 0))
				out = append(out, o.invoke(t, id, 4, 27, nil, 0, 0, 0, 0))
				o.c.GUI.FreeDestroyed()
				o.c.GUI.FreeDestroyed()
			}
			id++
		}
	}
	for _, op := range []int{16, 34} {
		o.plain(t)
		out = append(out, o.invoke(t, id, 0, op, o.parent, 20, 30, 40, 61))
		id++
	}
	meterCapture(t, "meters-constructors", out, len(out), "b452721f6fe2d1a6dff89767a9231459e52d471f8d238b2abe378ff9b179061c")
}
