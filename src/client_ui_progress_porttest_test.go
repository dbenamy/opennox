//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestClientUIRenderProgress(t *testing.T) {
	o := newObjectRenderOwner(t)
	c := o.c
	c.GUI = gui.New(c.Render())
	defer c.GUI.DestroyAll()
	draw, free := gui.NewWindowData()
	defer free()
	if legacy.Nox_gui_newProgressBar_4CAF10(nil, 0, 0, 0, 10, 10, draw) != nil {
		t.Fatal("progress constructor accepted missing style")
	}
	parent := c.GUI.NewWindowRaw(nil, 0, 3, 4, 90, 90, nil)
	type result struct {
		Case, Step  int
		Value       uint32
		EventReturn int
		Flags       uint32
		Pos, Size   image.Point
		Style       uint32
		OwnWindow   bool
		Smoothing   bool
		Render      objectRenderResult
	}
	var out []result
	id := 0
	for par := 0; par < 2; par++ {
		for mode := 0; mode < 2; mode++ {
			for smooth := 0; smooth < 2; smooth++ {
				for owner := 0; owner < 2; owner++ {
					for mask := 0; mask < 16; mask++ {
						for _, geom := range [][4]int{{10, 15, 61, 21}, {-4, 75, 17, 5}, {20, 20, 0, 0}} {
							id++
							o.resetRender(uint32(id), 120)
							c.r.SetTextSmooting(owner != 0)
							*draw = gui.WindowData{Style: 0x1000, BgColorVal: 0x21082108, EnColorVal: 0x7fff7fff, HlColorVal: 0x03e003e0, TextColorVal: 0x7c007c00, BgImageHnd: noxrender.ImageHandle(o.images[id%32].C())}
							colors := []*uint32{&draw.BgColorVal, &draw.EnColorVal, &draw.HlColorVal, &draw.TextColorVal}
							for i, p := range colors {
								if mask&(1<<i) != 0 {
									*p = 0x80000000
								}
							}
							if owner != 0 {
								draw.Window = parent
							}
							var p *gui.Window
							if par != 0 {
								p = parent
							}
							status := gui.StatusFlags(8 | mode*128 | smooth*8192)
							win := legacy.Nox_gui_newProgressBar_4CAF10(p, status, geom[0], geom[1], geom[2], geom[3], draw)
							if win == nil || win.Flags&8 != 0 || win.DrawData().Window != draw.Window || (owner == 0 && draw.Window != win) {
								t.Fatal("progress construction/copy/owner")
							}
							value := (*uint32)(unsafe.Add(win.C(), 32))
							*value = 37
							for step, v := range []int32{-1, 0, 1, 50, 99, 100, 101, 2147483647, -2147483648} {
								before := *value
								ret := gui.EventRespInt(win.Func94(gui.AsWindowEvent(16416, uintptr(uint32(v)), 0)))
								want := before
								if v >= 0 && v <= 100 {
									want = uint32(v)
								}
								if *value != want || ret != 0 {
									t.Fatal("progress update acceptance/return")
								}
								win.Func94(gui.AsWindowEvent(16417, 123, 0))
								if *value != want {
									t.Fatal("unrelated event changed progress")
								}
								win.Draw()
								out = append(out, result{id, step, *value, ret, uint32(win.Flags), win.GlobalPos(), win.SizeVal, uint32(win.DrawData().Style), win.DrawData().Window == win, c.r.PortTestUITextSmoothing(), o.renderResult(t, id, step, ret)})
							}
							win.Destroy()
							c.GUI.FreeDestroyed()
						}
					}
				}
			}
		}
	}
	effectsCapture(t, "ui-render-progress", out, len(out), "d503d5b21d7648b83286fd91b37e794fd5f2c941481dd19b2c0d6a14fa537700")
}

func TestClientUIRenderProgressPixels(t *testing.T) {
	o := newObjectRenderOwner(t)
	c := o.c
	c.GUI = gui.New(c.Render())
	defer c.GUI.DestroyAll()
	draw, free := gui.NewWindowData()
	defer free()
	draw.Style = 0x1000
	draw.BgColorVal = 0x21082108
	draw.HlColorVal = 0x03e003e0
	win := legacy.Nox_gui_newProgressBar_4CAF10(nil, 0, 10, 15, 60, 20, draw)
	win.Func94(gui.AsWindowEvent(16416, 0, 0))
	win.Draw()
	empty := effectsPixelHash(o.pix)
	win.Func94(gui.AsWindowEvent(16416, 100, 0))
	win.Draw()
	if effectsPixelHash(o.pix) == empty {
		t.Fatal("progress fill did not change actual framebuffer")
	}
}
