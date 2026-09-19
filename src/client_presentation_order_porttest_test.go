//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
)

func TestClientPresentationWallOrder(t *testing.T) {
	o := newObjectRenderOwner(t)
	player := o.drawable(7, image.Point{})
	clear(serverConfigOwnBytes(t, 0x852978, 8, 4))
	type record struct {
		Kind         byte
		Tile, Player image.Point
		Local        bool
		Y            int32
	}
	var rows []record
	offsets := []image.Point{{-24, -24}, {-1, 0}, {0, -1}, {0, 0}, {0, 1}, {1, 0}, {23, 23}, {24, 0}, {0, 24}}
	for _, local := range []bool{false, true} {
		if local {
			*memmap.PtrPtr(0x852978, 8) = player.C()
		} else {
			*memmap.PtrPtr(0x852978, 8) = nil
		}
		for _, tile := range []image.Point{{0, 0}, {1, 2}, {255, 255}} {
			for kind := 0; kind < 17; kind++ {
				for _, off := range offsets {
					player.PosVec = image.Pt(23*tile.X, 23*tile.Y).Add(off)
					p := [8]byte{byte(kind), 0, 0, 0, 0, byte(tile.X), byte(tile.Y), 0}
					want := int32(23*tile.Y + 11)
					if local {
						dx := 0
						ax := 23 * tile.X
						switch kind {
						case 0, 3, 11:
							dx = -23
							ax += 22
						case 1, 4, 12:
							dx = 23
						}
						if dx != 0 {
							// Choose the visible edge using the oriented wall segment, including its boundary.
							cross := int32(int64(dx)*int64(player.PosVec.Y-23*tile.Y) - 23*int64(player.PosVec.X-ax))
							if dx < 0 {
								cross = -cross
							}
							want = int32(23 * tile.Y)
							if cross < 0 {
								want += 22
							}
						}
					}
					got := legacy.PortTestPresentationWallY(&p)
					if got != want {
						t.Fatalf("kind=%d tile=%v player=%v local=%v: %d want %d", kind, tile, player.PosVec, local, got, want)
					}
					rows = append(rows, record{byte(kind), tile, player.PosVec, local, got})
				}
			}
		}
	}
	spellbookCapture(t, "client-presentation-wall-order", rows, "ca2bc375ae21e4975a8c7e70341c048092896877c0d3251f128b2ea35babed9a")
}
func TestClientPresentationDrawableOrder(t *testing.T) {
	o := newObjectRenderOwner(t)
	player := o.drawable(7, image.Point{})
	dr := o.drawable(8, image.Point{})
	clear(serverConfigOwnBytes(t, 0x852978, 8, 4))
	type record struct {
		Direction                byte
		Position, Player, Vector image.Point
		Local                    bool
		Y                        int32
	}
	var rows []record
	offsets := []image.Point{{-24, -24}, {-1, 0}, {0, -1}, {0, 0}, {0, 1}, {1, 0}, {23, 23}, {24, 0}, {0, 24}}
	for _, local := range []bool{false, true} {
		if local {
			*memmap.PtrPtr(0x852978, 8) = player.C()
		} else {
			*memmap.PtrPtr(0x852978, 8) = nil
		}
		for _, pos := range []image.Point{{0, 0}, {123, -456}, {100000, -100000}} {
			dr.PosVec = pos
			for dir := 0; dir < 32; dir++ {
				dr.Field_74_4 = byte(dir)
				v := image.Pt(int(memmap.Int32(0x587000, 196184+8*uintptr(dir))), int(memmap.Int32(0x587000, 196188+8*uintptr(dir))))
				if v == (image.Point{}) {
					t.Fatalf("missing original direction %d", dir)
				}
				for _, off := range offsets {
					player.PosVec = pos.Add(off)
					want := int32(pos.Y + v.Y/2)
					if local {
						side := int32(int64(off.Y)*int64(v.X) - int64(off.X)*int64(v.Y))
						if v.X < 0 {
							side = -side
						}
						want = int32(pos.Y)
						end := int32(pos.Y + v.Y)
						if side >= 0 {
							want = min(want, end)
						} else {
							want = max(want, end)
						}
					}
					got := legacy.PortTestPresentationDrawableY(dr)
					if got != want {
						t.Fatalf("%s: %d want %d", fmt.Sprint(local, pos, dir, off), got, want)
					}
					rows = append(rows, record{byte(dir), pos, player.PosVec, v, local, got})
				}
			}
		}
	}
	spellbookCapture(t, "client-presentation-drawable-order", rows, "d840310aac294772c1afdb25616e6209357cbd349c36ccdf8262e93855579a42")
}
