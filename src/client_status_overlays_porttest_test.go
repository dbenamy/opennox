//go:build porttest && !server

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"strings"
	"testing"
	"unicode/utf16"
	"unsafe"

	noxcolor "github.com/opennox/libs/color"
	"github.com/opennox/libs/noximage"
	"github.com/opennox/libs/player"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"golang.org/x/image/font/basicfont"
)

func clientOverlayOwner(t *testing.T) (*objectDrawingOwner, map[string]*uint32) {
	t.Helper()
	o := newObjectDrawingOwner(t)
	tables, pointers := blobdata.PortTestClientOverlayTables()
	for _, table := range tables {
		copy(serverConfigOwnBytes(t, table.Base, table.Offset, len(table.Data)), table.Data)
	}
	for _, pair := range pointers {
		binary.LittleEndian.PutUint32(serverConfigOwnBytes(t, 0x587000, pair[0], 4), uint32(uintptr(memmap.PtrOff(0x587000, pair[1]))))
	}

	globals, restore := legacy.PortTestChatBubbleRenderGlobals()
	t.Cleanup(restore)
	*globals["width"], *globals["height"], *globals["white"] = 320, 240, uint32(noxcolor.RGB5551Color(255, 255, 255))
	rect := image.Rect(0, 0, 320, 240)
	o.pix = noximage.NewImage16(rect)
	o.c.r.SetPixBuffer(o.pix)
	o.c.r.Data().SetClipRect(rect)
	o.c.r.Data().SetClipRect2(image.Rect(0, 0, 319, 239))
	o.c.r.Data().SetRect3(rect)
	_, restoreFont := o.c.r.GetFonts().PortTestWindowFont(basicfont.Face7x13, "default", "large")
	t.Cleanup(restoreFont)
	return o, globals
}

func TestClientLoadingOverlay(t *testing.T) {
	o, _ := clientOverlayOwner(t)
	set, restore := o.c.srv.Server.PortTestMeterStrings(strman.Entry{ID: "client.c:InProgress", Vals: []strman.Variant{{Str: "Loading game"}}})
	t.Cleanup(restore)
	set(0)
	cache := serverConfigOwnBytes(t, 0x5D4594, 814540, 4)
	stamp := serverConfigOwnBytes(t, 0x5D4594, 811920, 4)
	words, restore := legacy.PortTestClientInteractionWords()
	t.Cleanup(restore)
	*words["dword_5d4594_1305680"] = 0
	sessions, restore := legacy.PortTestSessionDialogWords()
	t.Cleanup(restore)
	*sessions["quit"], *sessions["motd"] = 0, 0
	options, restore := legacy.PortTestServerOptionsWords()
	t.Cleanup(restore)
	*options["root"] = 0
	scoreboard, restore := legacy.PortTestScoreboardWords()
	t.Cleanup(restore)
	*scoreboard["dword_5d4594_1090048"], *scoreboard["dword_5d4594_1090120"] = 0, 0
	oldConsole := legacy.Nox_gui_console_flagXxx_451410
	t.Cleanup(func() { legacy.Nox_gui_console_flagXxx_451410 = oldConsole })
	legacy.Nox_gui_console_flagXxx_451410 = func() int { return 0 }
	oldLoad := legacy.Nox_xxx_gLoadImg
	t.Cleanup(func() { legacy.Nox_xxx_gLoadImg = oldLoad })
	loads := 0
	legacy.Nox_xxx_gLoadImg = func(name string) *noxrender.Image {
		if name != "MenuSystemBG" {
			t.Fatal("loading image", name)
		}
		loads++
		return o.images[0]
	}
	type row struct {
		Frame, Old, Stored uint32
		Repeat, Loads      int
		Enabled            bool
		Pixels             string
	}
	var rows []row
	for _, pair := range [][2]uint32{{0, 0}, {0, 1}, {1, 0}, {2, 999}, {3, 2}, {3, 1}, {0xffffffff, 0xfffffffe}, {0, 0xffffffff}} {
		clear(cache)
		loads = 0
		o.c.srv.SetFrame(pair[0])
		binary.LittleEndian.PutUint32(stamp, pair[1])
		for repeat := 0; repeat < 2; repeat++ {
			clear(o.pix.Pix)
			before := binary.LittleEndian.Uint32(stamp)
			legacy.Nox_xxx_clientDrawAll_436100_draw_A()
			enabled := noxflags.HasEngine(noxflags.EngineFlag9)
			want := pair[0] == 2 || pair[0]-before == 1
			if enabled != want || loads != 1 || binary.LittleEndian.Uint32(cache) != uint32(uintptr(o.images[0].C())) {
				t.Fatal("loading gate/cache", pair, repeat, enabled, want, loads)
			}
			textPixels := false
			for _, pixel := range o.pix.Pix[80*o.pix.Stride:] {
				if pixel != 0 {
					textPixels = true
					break
				}
			}
			if textPixels != enabled {
				t.Fatal("loading text raster", pair, repeat, enabled, textPixels)
			}
			stored := before
			if pair[0] == 2 {
				stored = 2
			}
			if binary.LittleEndian.Uint32(stamp) != stored {
				t.Fatal("loading timestamp")
			}
			rows = append(rows, row{pair[0], pair[1], stored, repeat, loads, enabled, effectsPixelHash(o.pix)})
		}
	}
	interactionCapture(t, "client-loading-overlay", rows)
}

func TestClientWinnerOverlay(t *testing.T) {
	o, globals := clientOverlayOwner(t)
	text := serverConfigOwnBytes(t, 0x5D4594, 811376, 512)
	cache := serverConfigOwnBytes(t, 0x5D4594, 811888, 8)
	mode := serverConfigOwnBytes(t, 0x5D4594, 811060, 4)
	oldLoad := legacy.Nox_xxx_gLoadImg
	t.Cleanup(func() { legacy.Nox_xxx_gLoadImg = oldLoad })
	var loads []string
	legacy.Nox_xxx_gLoadImg = func(name string) *noxrender.Image {
		if name == "" {
			t.Fatal("missing winner image name")
		}
		loads = append(loads, name)
		return o.images[len(loads)%len(o.images)]
	}
	type row struct {
		Mode, Width, Height, Repeat int
		Text                        string
		Loads                       []string
		Pixels                      string
	}
	var rows []row
	for _, size := range []image.Point{{320, 240}, {321, 241}, {160, 120}} {
		*globals["width"], *globals["height"] = uint32(size.X), uint32(size.Y)
		for kind := 0; kind < 2; kind++ {
			for _, message := range []string{"", "Victory", "one\ntwo", "\r\none\n\n\rtwo\r", "long words wrap across the notice area before the next line", strings.Repeat("word ", 25) + "ab"} {
				clear(cache)
				clear(text)
				loads = nil
				binary.LittleEndian.PutUint32(mode, uint32(kind))
				units := utf16.Encode([]rune(message))
				if len(units) > 127 {
					t.Fatal("C local text capacity")
				}
				for i, u := range units {
					binary.LittleEndian.PutUint16(text[2*i:], u)
				}
				input := bytes.Clone(text)
				for repeat := 0; repeat < 2; repeat++ {
					clear(o.pix.Pix)
					legacy.Nox_xxx_clientDrawAll_436100_draw_B()
					if !bytes.Equal(text, input) || len(loads) != 1 {
						t.Fatal("winner source/cache", kind, message, repeat, loads)
					}
					if binary.LittleEndian.Uint32(cache[4*kind:]) == 0 {
						t.Fatal("winner cache slot")
					}
					rows = append(rows, row{kind, size.X, size.Y, repeat, message, append([]string(nil), loads...), effectsPixelHash(o.pix)})
				}
			}
		}
	}
	interactionCapture(t, "client-winner-overlay", rows)
}

func TestClientDebugOverlay(t *testing.T) {
	o, _ := clientOverlayOwner(t)
	entries := []strman.Entry{{ID: "client.c:PlayerInfo", Vals: []strman.Variant{{Str: "Level %d %s"}}}}
	classNames := []string{"Warrior", "Wizard", "Conjurer"}
	for _, name := range classNames {
		entries = append(entries, strman.Entry{ID: strman.ID("client.c:" + name), Vals: []strman.Variant{{Str: name}}})
	}
	set, restore := o.c.srv.Server.PortTestMeterStrings(entries...)
	t.Cleanup(restore)
	set(0)
	filename := serverConfigOwnBytes(t, 0x5D4594, 2598188, 80)
	scratch := serverConfigOwnBytes(t, 0x5D4594, 811120, 160)
	players, freePlayers := o.c.srv.PortTestObjectRenderPlayers()
	t.Cleanup(freePlayers)
	oldPlayer := legacy.Get_dword_8531A0_2576()
	t.Cleanup(func() { legacy.Set_dword_8531A0_2576(oldPlayer) })
	dr, free := alloc.New(client.Drawable{})
	t.Cleanup(free)
	active := serverConfigOwnBytes(t, 0x852978, 8, 4)
	type row struct {
		Map                   string
		Present, Class, Level int
		Origin, Position      image.Point
		Text, Pixels          string
	}
	var rows []row
	for _, mapName := range []string{"", "arena.map", strings.Repeat("x", 79)} {
		clear(filename)
		copy(filename, mapName)
		for present := 0; present < 3; present++ {
			clear(active)
			legacy.Set_dword_8531A0_2576(nil)
			if present != 0 {
				binary.LittleEndian.PutUint32(active, uint32(uintptr(unsafe.Pointer(dr))))
			}
			if present == 2 {
				legacy.Set_dword_8531A0_2576(&players[0])
			}
			for class := 0; class < 3; class++ {
				players[0].Info().SetPlayerClass(player.Class(class))
				for _, level := range []byte{0, 1, 127, 128, 255} {
					players[0].Level = level
					for _, origin := range []image.Point{{0, 0}, {10, 20}, {-7, 5}} {
						o.c.Viewport().Screen.Min = origin
						dr.PosVec = image.Pt(-123, 456)
						clear(o.pix.Pix)
						clear(scratch)
						legacy.Sub_436F50()
						got := alloc.GoString16((*uint16)(unsafe.Pointer(&scratch[0])))
						want := mapName
						if present == 2 {
							want = fmt.Sprintf("Level %d %s", int8(level), classNames[class])
						}
						if got != want {
							t.Fatal("debug overlay text", present, class, level, got, want)
						}
						rows = append(rows, row{mapName, present, class, int(level), origin, dr.PosVec, got, effectsPixelHash(o.pix)})
					}
				}
			}
		}
	}
	interactionCapture(t, "client-debug-overlay", rows)
}
