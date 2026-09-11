//go:build porttest

package opennox

import (
	"bytes"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

const (
	tileMinInt32 = -1 << 31
	tileMaxInt32 = 1<<31 - 1
)

func TestTileSelectionABI(t *testing.T) {
	if unsafe.Sizeof(server.TileDef{}) != 60 || unsafe.Offsetof(server.TileDef{}.NameBuf) != 0 || unsafe.Offsetof(server.TileDef{}.Field52) != 52 || unsafe.Offsetof(server.TileDef{}.Field53) != 53 {
		t.Fatal("tile layout mismatch")
	}
	before := legacy.PortTestTileSelectionState()
	checkRestored := func() {
		t.Helper()
		if !bytes.Equal(before, legacy.PortTestTileSelectionState()) {
			t.Fatal("tile fixture did not restore table/count/state/guards")
		}
	}
	defer checkRestored()

	names := []legacy.PortTestTileNameSpec{{Name: []byte("ALPHA")}, {Name: []byte("none")}, {Name: nil}, {Name: []byte("missing")}, {Name: []byte{0x80, 'X'}}, {Name: []byte{0x81, 'X'}}, {Name: []byte("aLpHa")}, {Name: []byte("1234567890123456789012345678901")}, {Name: []byte{'A', 0, 'l', 'p', 'h', 'a'}}, {Name: []byte("1234567890123456789012345678901x")}, {Name: []byte("Alpha\x00ignored")}, {Name: []byte{0x80, 'x'}}}
	wantName := []struct {
		ret      int
		selected uint32
	}{{1, 17}, {1, 255}, {1, 175}, {0, 0}, {1, 41}, {1, 42}, {1, 17}, {1, 93}, {0, 0}, {0, 0}, {1, 17}, {1, 41}}
	gotName := legacy.PortTestTileNames(names)
	checkRestored()
	if len(gotName) != len(names) {
		t.Fatal("name result length")
	}
	for i, got := range gotName {
		w := wantName[i]
		if got.Return != w.ret || got.Selected != w.selected || got.Variation != 0x11223344 || got.Flag != 1 || got.Count != 1 || !got.TableUnchanged || !got.InputUnchanged || !got.GuardsOK {
			t.Fatalf("name %d got=%+v want ret/selected=%d/%d", i, got, w.ret, w.selected)
		}
	}

	ints := []int32{tileMinInt32, -1}
	for i := int32(0); i < 176; i++ {
		ints = append(ints, i)
	}
	ints = append(ints, 176, tileMaxInt32)
	gotImage := legacy.PortTestTileScalars(ints, false)
	checkRestored()
	if len(gotImage) != len(ints) {
		t.Fatal("image result length")
	}
	for i, g := range gotImage {
		v := ints[i]
		want := 0
		sel := uint32(0)
		if v >= 0 && v < 176 {
			want = 1
			sel = uint32(v)
		}
		if g.Return != want || g.Selected != sel || g.Variation != 0x11223344 || g.Flag != 0xa5a5a5a5 || g.Count != 7 || !g.TableUnchanged || !g.GuardsOK {
			t.Fatalf("image %d got=%+v", v, g)
		}
	}
	gotFlag := legacy.PortTestTileScalars(ints, true)
	checkRestored()
	if len(gotFlag) != len(ints) {
		t.Fatal("flag result length")
	}
	for i, g := range gotFlag {
		v := ints[i]
		want := 0
		flag := uint32(0xa5a5a5a5)
		if v == 0 || v == 1 {
			want = 1
			flag = uint32(v)
		}
		if g.Return != want || g.Selected != 99 || g.Variation != 0x11223344 || g.Flag != flag || g.Count != 7 || !g.TableUnchanged || !g.GuardsOK {
			t.Fatalf("flag %d got=%+v", v, g)
		}
	}

	var vars []legacy.PortTestTileVariationSpec
	for w := 0; w < 256; w++ {
		for h := 0; h < 256; h++ {
			product := int64(w * h)
			for _, v := range []int64{-1, 0, product - 1, product, product + 1, tileMinInt32, tileMaxInt32} {
				vars = append(vars, legacy.PortTestTileVariationSpec{Selected: uint32((w + h) % 176), Width: byte(w), Height: byte(h), Value: int32(v)})
			}
		}
	}
	gotVars := legacy.PortTestTileVariations(vars)
	checkRestored()
	if len(gotVars) != len(vars) {
		t.Fatal("variation result length")
	}
	for i, s := range vars {
		g := gotVars[i]
		product := int64(s.Width) * int64(s.Height)
		want := 0
		variation := uint32(0)
		if int64(s.Value) <= product-1 {
			want = 1
			variation = uint32(s.Value)
		}
		if g.Return != want || g.Selected != s.Selected || g.Variation != variation || g.Flag != 1 || g.Count != 7 || !g.TableUnchanged || !g.GuardsOK {
			t.Fatalf("variation %d spec=%+v got=%+v want=%d/%08x", i, s, g, want, variation)
		}
	}
}
