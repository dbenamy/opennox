//go:build porttest

package legacy

/*
#include "defs.h"
#include <stdio.h>
#include <stdint.h>
extern uint32_t dword_5d4594_2487524;
*/
import "C"
import (
	"bytes"
	"io"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

// Reserved action values resolve to real file and token-buffer pointers.
const PortTestThemeFile = 0x60000001
const PortTestThemeToken = 0x60000002

func PortTestMapTheme(cases []PortTestPaintSpec, owner func(*server.Server) (Server, func())) []PortTestPaintResult {
	ext := &paintTestExtension{globals: map[string]*uint32{
		"themeLine":     memmap.PtrUint32(0x5D4594, 2487520),
		"themeTemplate": (*uint32)(unsafe.Pointer(&C.dword_5d4594_2487524)),
	}}
	for _, name := range []string{"themeInputPath", "themeFile", "themeOffset", "themeError", "themeClock", "themePlayers", "themeOpenFiles", "themeClosedFiles"} {
		p, free := alloc.New(uint32(0))
		defer free()
		ext.globals[name] = p
	}
	token := memmap.PtrOff(0x5D4594, 2487264)
	table := memmap.PtrOff(0x587000, 253144)
	oldToken := bytes.Clone(unsafe.Slice((*byte)(token), 256))
	oldTable := bytes.Clone(unsafe.Slice((*byte)(table), 368))
	defer func() {
		copy(unsafe.Slice((*byte)(token), 256), oldToken)
		copy(unsafe.Slice((*byte)(table), 368), oldTable)
	}()
	copy(unsafe.Slice((*byte)(table), 368), blobdata.PortTestMapThemeTableData())
	*memmap.PtrPtr(0x587000, 253144) = memmap.PtrOff(0x587000, 253296)
	*memmap.PtrPtr(0x587000, 253148) = memmap.PtrOff(0x587000, 253304)
	*memmap.PtrPtr(0x587000, 253152) = memmap.PtrOff(0x587000, 253312)
	*memmap.PtrPtr(0x587000, 253156) = memmap.PtrOff(0x587000, 253320)
	*memmap.PtrPtr(0x587000, 253160) = memmap.PtrOff(0x587000, 253328)
	*memmap.PtrPtr(0x587000, 253164) = memmap.PtrOff(0x587000, 253336)
	*memmap.PtrPtr(0x587000, 253172) = memmap.PtrOff(0x587000, 253340)
	*memmap.PtrPtr(0x587000, 253176) = memmap.PtrOff(0x587000, 253348)
	*memmap.PtrPtr(0x587000, 253180) = memmap.PtrOff(0x587000, 253356)
	*memmap.PtrPtr(0x587000, 253184) = memmap.PtrOff(0x587000, 253364)
	*memmap.PtrPtr(0x587000, 253188) = memmap.PtrOff(0x587000, 253372)
	*memmap.PtrPtr(0x587000, 253192) = memmap.PtrOff(0x587000, 253380)
	*memmap.PtrPtr(0x587000, 253200) = memmap.PtrOff(0x587000, 253388)
	*memmap.PtrPtr(0x587000, 253204) = memmap.PtrOff(0x587000, 253396)
	*memmap.PtrPtr(0x587000, 253208) = memmap.PtrOff(0x587000, 253412)
	*memmap.PtrPtr(0x587000, 253216) = memmap.PtrOff(0x587000, 253432)
	*memmap.PtrPtr(0x587000, 253220) = memmap.PtrOff(0x587000, 253440)
	*memmap.PtrPtr(0x587000, 253224) = memmap.PtrOff(0x587000, 253448)
	*memmap.PtrPtr(0x587000, 253228) = memmap.PtrOff(0x587000, 253464)
	*memmap.PtrPtr(0x587000, 253232) = memmap.PtrOff(0x587000, 253472)
	*memmap.PtrPtr(0x587000, 253236) = memmap.PtrOff(0x587000, 253480)
	*memmap.PtrPtr(0x587000, 253244) = memmap.PtrOff(0x587000, 253488)
	*memmap.PtrPtr(0x587000, 253248) = memmap.PtrOff(0x587000, 253496)

	var active *paintTestFixture
	var input *binfile.Binfile
	var oldHandles map[unsafe.Pointer]bool
	var seenHandles map[unsafe.Pointer]bool
	var playersUnchanged func() bool
	var restorePlayers func()
	ext.setup = func(p *server.PortTestPaintOwners) func() {
		oldAlloc, oldFree := themeTestOnAllocate, themeTestOnRelease
		themeTestOnAllocate = func(ptr unsafe.Pointer, size int) {
			if active != nil {
				r := active.register(ptr, size, "themeAllocation", false)
				active.owned[r] = true
			}
		}
		themeTestOnRelease = func(ptr unsafe.Pointer) {
			if active != nil {
				if r := active.known(ptr); r != nil {
					r.alive = false
					delete(active.owned, r)
				}
			}
		}
		return func() { themeObserve(false, 0); themeTestOnAllocate = oldAlloc; themeTestOnRelease = oldFree }
	}
	ext.before = func(f *paintTestFixture, sp PortTestPaintSpec) {
		active = f
		seenHandles = map[unsafe.Pointer]bool{}
		if ptr := *ext.globals["themePlayers"]; ptr != 0 {
			count := *populationWord(ptr, 0)
			if count > 32 {
				panic("theme fixture player count")
			}
			classes := make([]byte, count)
			for i := range classes {
				classes[i] = byte(*populationWord(ptr, 4+4*i))
			}
			playersUnchanged, restorePlayers = f.owners.S.PortTestThemePlayers(classes)
		}
		clear(unsafe.Slice((*byte)(token), 256))
		f.register(token, 256, "themeToken", false)
		f.register(table, 368, "themeTables", false)
		oldHandles = map[unsafe.Pointer]bool{}
		files.RLock()
		for h := range files.byHandle {
			oldHandles[h] = true
		}
		files.RUnlock()
		input = nil
		if path := *ext.globals["themeInputPath"]; path != 0 {
			var err error
			input, err = binfile.BinfileOpen(populationString(path), binfile.ReadOnly)
			if err != nil {
				panic(err)
			}
			handle := NewFileHandle(input.File)
			*ext.globals["themeFile"] = mapRoomRaw(unsafe.Pointer(handle))
			f.register(unsafe.Pointer(handle), 0, "themeFile", false)
			seenHandles[unsafe.Pointer(handle)] = true
		}
	}
	ext.invoke = func(f *paintTestFixture, op int, args [6]uint32) uint32 {
		for i, v := range args {
			switch v {
			case PortTestThemeFile:
				args[i] = *ext.globals["themeFile"]
			case PortTestThemeToken:
				args[i] = mapRoomRaw(token)
			}
		}
		themeObserve(true, *ext.globals["themeClock"])
		defer themeObserve(false, 0)
		return themeInvokeGo(op, args)
	}
	ext.after = func(f *paintTestFixture, op int, ret uint32) {
		if input != nil {
			offset, err := input.File.File.Seek(0, io.SeekCurrent)
			if err != nil {
				panic(err)
			}
			*ext.globals["themeOffset"] = uint32(offset)
			*ext.globals["themeError"] = uint32(bool2int(input.File.Err != nil))
		}
		if playersUnchanged != nil {
			f.intact = playersUnchanged() && f.intact
		}
		*ext.globals["themeOpenFiles"], *ext.globals["themeClosedFiles"] = 0, 0
		files.RLock()
		defer files.RUnlock()
		for h, file := range files.byHandle {
			if !oldHandles[h] {
				if file.File.Fd() == ^uintptr(0) {
					*ext.globals["themeClosedFiles"]++
				} else {
					*ext.globals["themeOpenFiles"]++
				}
			}
			if !oldHandles[h] && !seenHandles[h] {
				seenHandles[h] = true
				r := f.register(h, 0, "themeFile", false)
				if file.File.Fd() == ^uintptr(0) {
					r.alive = false
				}
			}
		}
	}
	ext.finish = func(f *paintTestFixture) {
		themeObserve(false, 0)
		active = nil
		if restorePlayers != nil {
			restorePlayers()
			restorePlayers = nil
			playersUnchanged = nil
		}
		files.Lock()
		defer files.Unlock()
		for h, file := range files.byHandle {
			if !oldHandles[h] {
				_ = file.Close()
				delete(files.byHandle, h)
			}
		}
	}
	return portTestMapPainting(cases, owner, ext)
}
