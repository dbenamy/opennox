//go:build porttest

package opennox

import (
	"reflect"
	"testing"
	"unsafe"

	noxcolor "github.com/opennox/libs/color"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func resourceExpectedColors(c server.PlayerColors, white uint32) server.PlayerUnitColors {
	pack := func(c types.RGB) uint32 { return noxcolor.RGB5551Color(c.R, c.G, c.B).Color32() }
	return server.PlayerUnitColors{Skin: pack(c.Skin), Hair: pack(c.Hair), Mustache: pack(c.Mustache), Goatee: pack(c.Goatee), Beard: pack(c.Beard), UnkColor: white}
}
func TestClientResourcePlayerColors(t *testing.T) {
	globals, restore := legacy.PortTestChatBubbleRenderGlobals()
	t.Cleanup(restore)
	p, free := alloc.New(server.Player{})
	t.Cleanup(free)
	type row struct {
		Sample int
		White  uint32
		Input  server.PlayerColors
		Output server.PlayerUnitColors
	}
	var rows []row
	for _, white := range []uint32{0, 0x7fff, 0x7fff7fff, 0x12345678} {
		for sample := 0; sample < 256; sample++ {
			*globals["white"] = white
			colors := server.PlayerColors{Hair: types.RGB{byte(sample), 0, 255}, Skin: types.RGB{0, byte(sample), 255}, Mustache: types.RGB{255, 0, byte(sample)}, Goatee: types.RGB{byte(sample), byte(sample), byte(sample)}, Beard: types.RGB{byte(sample * 17), byte(sample * 31), byte(sample * 61)}, Pants: 17, Shirt1: 33, Shirt2: 65, Shoes1: 129, Shoes2: 255}
			p.Info().Colors = colors
			p.Colors = server.PlayerUnitColors{}
			legacy.Nox_xxx_playerInitColors_461460(p)
			if want := resourceExpectedColors(colors, white); p.Colors != want || p.Info().Colors != colors {
				t.Fatal("player colors", sample, white, p.Colors, want)
			}
			rows = append(rows, row{sample, white, colors, p.Colors})
		}
	}
	interactionCapture(t, "client-resource-player-colors", rows)
}
func TestClientResourceAllPlayerColors(t *testing.T) {
	o := newObjectDrawingOwner(t)
	players, free := o.c.srv.PortTestObjectRenderPlayers()
	t.Cleanup(free)
	globals, restore := legacy.PortTestChatBubbleRenderGlobals()
	t.Cleanup(restore)
	*globals["white"] = 0x7fff7fff
	type row struct {
		Pattern int
		Active  []byte
		Colors  []server.PlayerUnitColors
	}
	var rows []row
	for pattern := 0; pattern < 5; pattern++ {
		var expected []server.PlayerUnitColors
		var active []byte
		for i := range players {
			p := &players[i]
			p.PlayerInd = byte(i)
			p.Active = 0
			if pattern == 1 || pattern == 2 && i%2 == 0 || pattern == 3 && i == 31 || pattern == 4 && i%7 == 0 {
				p.Active = 1
			}
			p.Info().Colors = server.PlayerColors{Hair: types.RGB{byte(i * 3), byte(i * 7), byte(i * 11)}, Skin: types.RGB{128, 64, 32}}
			p.Colors = server.PlayerUnitColors{Skin: 0xdeadbeef, Hair: 0x12345678}
			want := p.Colors
			if p.Active != 0 {
				want = resourceExpectedColors(p.Info().Colors, *globals["white"])
			}
			expected = append(expected, want)
			active = append(active, p.Active)
		}
		legacy.Sub_461520()
		var got []server.PlayerUnitColors
		for i := range players {
			got = append(got, players[i].Colors)
		}
		if !reflect.DeepEqual(got, expected) {
			t.Fatal("active player color traversal", pattern)
		}
		rows = append(rows, row{pattern, active, got})
	}
	interactionCapture(t, "client-resource-all-player-colors", rows)
}
func TestClientResourceAnimationCache(t *testing.T) {
	cache := serverConfigOwnBytes(t, 0x5D4594, 1096456, 8)
	refs, free := alloc.Make([]legacy.ImageRef{}, 2)
	t.Cleanup(free)
	old := legacy.Nox_xxx_gLoadAnim
	t.Cleanup(func() { legacy.Nox_xxx_gLoadAnim = old })
	type row struct {
		Mask, Step int
		Calls      []string
		Stored     [2]uint32
	}
	var rows []row
	for mask := 0; mask < 4; mask++ {
		var calls []string
		legacy.Nox_xxx_gLoadAnim = func(name string) *legacy.ImageRef {
			calls = append(calls, name)
			index := 0
			if name == "SphericalShieldAnim" {
				index = 1
			} else if name != "ConfusedBirdies" {
				t.Fatal(name)
			}
			if mask&(1<<index) != 0 {
				return nil
			}
			return &refs[index]
		}
		for step := 0; step < 4; step++ {
			if step < 2 {
				legacy.Sub_473930()
			} else {
				legacy.Sub_473960()
			}
			r := row{Mask: mask, Step: step, Calls: append([]string(nil), calls...)}
			words := unsafe.Slice((*uint32)(unsafe.Pointer(&cache[0])), 2)
			for i, v := range words {
				want := uint32(0)
				if step < 2 && mask&(1<<i) == 0 {
					want = uint32(uintptr(unsafe.Pointer(&refs[i])))
				}
				if v != want {
					t.Fatal("animation cache", mask, step, i, v, want)
				}
				if v != 0 {
					r.Stored[i] = uint32(i + 1)
				}
			}
			if len(calls) != 2*min(step+1, 2) {
				t.Fatal("animation reload/clear calls", step, calls)
			}
			rows = append(rows, r)
		}
	}
	interactionCapture(t, "client-resource-animation-cache", rows)
}
