//go:build porttest

package opennox

import (
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
)

func TestClientResourceMainMenuDraw(t *testing.T) {
	o := newEntryOwner(t)
	table := serverConfigOwnBytes(t, 0x587000, 168832, 4*48)
	win := o.c.GUI.NewWindowRaw(o.parent, 8, 10, 12, 50, 40, nil)
	draw := win.DrawData()
	draw.SetImagePoint(image.Pt(3, 5))
	draw.BgImageHnd = o.images[0].C()
	draw.HlImageHnd = o.images[1].C()
	call := gui.WrapDrawFuncC(legacy.Get_sub_4A22A0())
	type row struct {
		Records, Names, Advance, Style, Seed, Pattern, Step int
		Random                                              uint32
		State                                               []uint32
		Pixels                                              string
	}
	var rows []row
	clear(o.pix.Pix)
	blank := effectsPixelHash(o.pix)
	initial := [][4]uint32{{0, 1, 0, 0}, {1, 1, 0, 1}, {0, 2, 1, 2}, {0, 0, 0, 0}, {1, 0xffffffff, 0, 0}, {0, 3, 2, 1}}
	for _, count := range []int{0, 1, 3} {
		for names := 0; names < 2; names++ {
			for advance := 0; advance < 2; advance++ {
				for style := 0; style < 4; style++ {
					for _, seed := range []int{1, 17, 91} {
						for pattern, start := range initial {
							clear(table)
							words := unsafe.Slice((*uint32)(unsafe.Pointer(&table[0])), len(table)/4)
							for i := 0; i < count; i++ {
								w := words[i*12:][:12]
								w[0] = 1
								w[1] = uint32(uintptr(o.images[i+2].C()))
								w[2] = uint32(10 + 12*i)
								w[3] = uint32(35 + 9*i)
								w[4] = start[0]
								w[5] = 2
								w[6] = 4
								w[7] = 1
								w[8] = 3
								w[9] = start[1]
								w[10] = start[2]
								w[11] = start[3]
							}
							if names == 0 {
								words[0] = 0
							}
							if advance == 0 {
								words[1] = 0
							}
							win.Flags = 8
							draw.Field0 = 0
							draw.BgColorVal = 0x12345678
							switch style {
							case 1:
								draw.BgColorVal = 0x80000000
							case 2:
								win.Flags |= gui.StatusImage
							case 3:
								win.Flags |= gui.StatusImage
								draw.Field0 = 2
							}
							o.c.srv.Rand.Other = prand.New(seed)
							for step := 0; step < 4; step++ {
								clear(o.pix.Pix)
								before := append([]byte(nil), table...)
								index := o.c.srv.Rand.Other.Index()
								if got := call(win, draw); got != 1 {
									t.Fatal("main menu draw return", got)
								}
								if style != 1 && effectsPixelHash(o.pix) == blank {
									t.Fatal("menu backdrop did not draw", style)
								}
								if advance == 0 || count == 0 {
									if o.c.srv.Rand.Other.Index() != index {
										t.Fatal("inactive menu animation consumed randomness")
									}
									for i, v := range before {
										if table[i] != v {
											t.Fatal("inactive menu animation changed table")
										}
									}
								}
								normalized := append([]uint32(nil), words...)
								for i := 0; i < count; i++ {
									if normalized[i*12+1] != 0 {
										normalized[i*12+1] = uint32(i + 1)
									}
								}
								rows = append(rows, row{count, names, advance, style, seed, pattern, step, uint32(o.c.srv.Rand.Other.Index()), normalized, effectsPixelHash(o.pix)})
							}
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "client-resource-main-menu-draw", rows)
}
