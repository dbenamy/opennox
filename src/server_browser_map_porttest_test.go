//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"testing"
	"unsafe"
)

func TestServerBrowserMapMarkers(t *testing.T) {
	o := newEntryOwner(t)
	words, restore := legacy.PortTestServerBrowserWords()
	defer restore()
	*words["nox_wol_wnd_world_814980"] = uint32(uintptr(o.parent.C()))
	*words["dword_5d4594_814984"] = uint32(uintptr(o.parent.C()))
	*words["dword_5d4594_814988"] = uint32(uintptr(o.parent.C()))
	regions := serverConfigOwnBytes(t, 0x587000, 87528, 40)
	for i := 0; i < 5; i++ {
		binary.LittleEndian.PutUint16(regions[i*8:], uint16(i*100))
		binary.LittleEndian.PutUint16(regions[i*8+2:], uint16(i*50))
		binary.LittleEndian.PutUint16(regions[i*8+4:], uint16(i*100+90))
		binary.LittleEndian.PutUint16(regions[i*8+6:], uint16(i*50+45))
		w := o.c.GUI.NewWindowRaw(o.parent, 8, 0, 0, 200, 200, nil)
		w.SetID(uint(10054 + i))
	}
	thresholds := serverConfigOwnBytes(t, 0x587000, 87484, 12)
	for i, v := range []uint32{50, 150, 300} {
		binary.LittleEndian.PutUint32(thresholds[4*i:], v)
	}
	images := serverConfigOwnBytes(t, 0x5D4594, 814556, 48)
	for i := 0; i < 12; i++ {
		binary.LittleEndian.PutUint32(images[i*4:], uint32(uintptr(o.images[i%len(o.images)].C())))
	}
	raw, free := alloc.Make([]byte{}, 172)
	defer free()
	type row struct {
		Region                  int32
		InputX, InputY          int16
		Ping                    int32
		Name                    string
		X, Y                    int16
		ID                      uint
		Off, Size               image.Point
		Tooltip                 string
		Normal, Hover, Disabled int
	}
	var rows []row
	for _, region := range []int32{-1, 0, 2, 4} {
		for _, pt := range [][2]int16{{10, 10}, {220, 110}, {410, 210}, {-1, -1}, {32767, -32768}} {
			for _, ping := range []int32{0, 50, 51, 150, 151, 9999, -1} {
				for _, name := range []string{"", "Example"} {
					clear(raw)
					*words["dword_587000_87412"] = uint32(region)
					copy(raw[12:28], "127.0.0.1")
					binary.LittleEndian.PutUint16(raw[109:], 18590)
					copy(raw[120:135], name)
					binary.LittleEndian.PutUint32(raw[36:], 17)
					binary.LittleEndian.PutUint16(raw[44:], uint16(pt[0]))
					binary.LittleEndian.PutUint16(raw[46:], uint16(pt[1]))
					binary.LittleEndian.PutUint32(raw[96:], uint32(ping))
					legacy.PortTestServerBrowserRow(unsafe.Pointer(&raw[0]))
					w := (*gui.Window)(unsafe.Pointer(uintptr(binary.LittleEndian.Uint32(raw[28:]))))
					if w == nil {
						t.Fatal("missing marker")
					}
					r := row{Region: region, InputX: pt[0], InputY: pt[1], Ping: ping, Name: name, X: int16(binary.LittleEndian.Uint16(raw[44:])), Y: int16(binary.LittleEndian.Uint16(raw[46:])), ID: w.ID(), Off: w.Off, Size: w.SizeVal, Tooltip: w.DrawData().Tooltip()}
					index := int(region)
					size, offset := 20, 10
					if region < 0 {
						index = 0
						size, offset = 10, 5
						for i := 0; i < 4; i++ {
							if int(pt[0]) > i*100 && int(pt[0]) < i*100+90 && int(pt[1]) > i*50 && int(pt[1]) < i*50+45 {
								index = i
								break
							}
						}
					}
					x, y := uint16(pt[0])-uint16(index*100), uint16(pt[1])-uint16(index*50)
					if region < 0 {
						x >>= 1
						y >>= 1
					}
					if r.X != int16(x) || r.Y != int16(y) || r.ID != 10087 || r.Off != image.Pt(int(int16(x))-offset, int(int16(y))-offset) || r.Size != image.Pt(size, size) {
						t.Fatalf("marker geometry: %+v", r)
					}
					title := name
					if title == "" {
						title = "127.0.0.1:18590"
					}
					tip := fmt.Sprintf("%s %dms", title, ping)
					if ping == 9999 {
						tip = title + " -- ms"
					}
					if r.Tooltip != tip {
						t.Fatalf("tooltip %q want %q", r.Tooltip, tip)
					}
					band := 0
					if uint32(ping) > 50 {
						band = 1
					}
					if uint32(ping) > 150 {
						band = 2
					}
					base := band * 4
					if region < 0 {
						base += 2
					}
					handles := []uintptr{uintptr(w.DrawData().BgImageHnd), uintptr(w.DrawData().HlImageHnd), uintptr(w.DrawData().DisImageHnd)}
					for i, handle := range handles {
						want := base
						if i == 1 {
							want++
						}
						if handle != uintptr(o.images[want%len(o.images)].C()) {
							t.Fatal("marker image", region, ping, i)
						}
					}
					r.Normal = base
					r.Hover = base + 1
					r.Disabled = base
					rows = append(rows, r)
					w.Destroy()
					o.c.GUI.FreeDestroyed()
				}
			}
		}
	}
	spellbookCapture(t, "server-browser-map-markers", rows, "cb06c5c8ad309b364b0f0d1cc28f3737bc5a2688d6609c440e7cee8f29b726ba")
}
