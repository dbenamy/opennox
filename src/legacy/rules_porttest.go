//go:build porttest

package legacy

/*
#include <stdlib.h>
#include "GAME1_1.h"
#include "GAME5_2.h"
extern unsigned int dword_5d4594_2650652;
*/
import "C"

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf16"
	"unsafe"

	"github.com/opennox/libs/ifs"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"github.com/opennox/opennox/v1/server"
)

var _ = [1]struct{}{}[60-unsafe.Sizeof(server.Settings2{})]

var PortTestRuleModes = []struct {
	Name  string
	Flags uint32
}{
	{"[ELIMINATION]", 0x400}, {"[DEATHMATCH]", 0x100}, {"[CAPTURE_THE_FLAG]", 0x20},
	{"[KING_OF_THE_REALM]", 0x10}, {"[FLAGBALL]", 0x40}, {"[COMMON]", 0x17f0}, {"[QUEST]", 0x1000},
}

func portTestRuleTable() func() {
	raw := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 312208)), 56)
	old := append([]byte(nil), raw...)
	var frees []func()
	for i, m := range PortTestRuleModes {
		p, free := alloc.CString(m.Name)
		frees = append(frees, free)
		*memmap.PtrPtr(0x587000, uintptr(312208+8*i)) = unsafe.Pointer(p)
		*memmap.PtrUint32(0x587000, uintptr(312212+8*i)) = m.Flags
	}
	return func() {
		copy(raw, old)
		for _, free := range frees {
			free()
		}
	}
}

type PortTestRulesSpec struct {
	Kind                                 string // load, file, lines, tokens, or write
	Initial                              [60]byte
	Selection                            uint8
	Flags                                uint16
	Context, Online                      uint32
	Host, WithRejected, Wrapper, UserNil bool
	User, Dir, Path                      string
	Files                                map[string]string
	Lines                                []string
	Tokens                               [][]string
	SeedRejected                         []string
}

type PortTestRulesState struct {
	Settings                                                  [60]byte
	Context                                                   uint32
	Result                                                    uint8
	Rejected                                                  []string
	LinksValid, GuardsValid, TableUnchanged, HandlesUnchanged bool
}

type PortTestRulesResult struct {
	Catalog        server.PortTestRuleServer
	Steps          []PortTestRulesState
	FilesUnchanged bool
	Written        string
	WriteExists    bool
}

func PortTestRuleHeaders(values []uint16) []int {
	restore := portTestRuleTable()
	defer restore()
	out := make([]int, 0, len(values))
	for _, v := range values {
		p := unsafe.Pointer(C.sub_57A1B0(C.short(v)))
		index := -1
		if p != nil {
			index = -2
			for i := range PortTestRuleModes {
				if p == *memmap.PtrPtr(0x587000, uintptr(312208+8*i)) {
					index = i
					break
				}
			}
		}
		out = append(out, index)
	}
	return out
}

func PortTestRules(spec PortTestRulesSpec) (out PortTestRulesResult, err error) {
	if spec.Kind == "load" || spec.Kind == "file" || spec.Kind == "write" {
		restoreHandles := handles.PortTestInit()
		defer restoreHandles()
	}
	if spec.Dir != "" {
		for name, data := range spec.Files {
			path := filepath.Join(spec.Dir, filepath.FromSlash(name))
			if err = os.MkdirAll(filepath.Dir(path), 0700); err != nil {
				return out, err
			}
			if err = os.WriteFile(path, []byte(data), 0600); err != nil {
				return out, err
			}
		}
		old, e := ifs.Workdir()
		if e != nil {
			return out, e
		}
		if err = ifs.Chdir(spec.Dir); err != nil {
			return out, err
		}
		defer func() {
			if e := ifs.Chdir(old); err == nil {
				err = e
			}
		}()
	}
	core, catalog := server.PortTestRuleServerSetup()
	out.Catalog = catalog
	oldGet := GetServer
	GetServer = func() Server { return &portTestRandomServer{core: core} }
	defer func() { GetServer = oldGet }()
	gameFlags := noxflags.GameFlag(0)
	if spec.Host {
		gameFlags = noxflags.GameHost
	}
	restoreFlags := noxflags.PortTestGameFlags(gameFlags)
	defer restoreFlags()
	restoreTable := portTestRuleTable()
	defer restoreTable()
	table := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 312208)), 56)
	tableBefore := append([]byte(nil), table...)
	oldOnline, oldContext := C.dword_5d4594_2650652, ruleLoaderContext
	C.dword_5d4594_2650652, ruleLoaderContext = C.uint32_t(spec.Online), spec.Context
	defer func() { C.dword_5d4594_2650652, ruleLoaderContext = oldOnline, oldContext }()

	buf, freeBuf := alloc.Make([]byte{}, 76)
	defer freeBuf()
	for i := range buf {
		buf[i] = 0xa5
	}
	copy(buf[8:68], spec.Initial[:])
	settings := (*server.Settings2)(unsafe.Pointer(&buf[8]))
	var head *C.nox_list_item_t
	if spec.WithRejected {
		head = (*C.nox_list_item_t)(C.calloc(1, 12))
		if head == nil {
			panic("fixture allocation failed")
		}
		C.nox_common_list_clear_425760(head)
		defer func() { C.sub_57ADF0((*C.int)(unsafe.Pointer(head))); C.free(unsafe.Pointer(head)) }()
		for _, line := range spec.SeedRejected {
			text := utf16.Encode([]rune(line))
			if len(text) > 255 {
				return out, fmt.Errorf("oversized fixture rejected line")
			}
			p := C.calloc(1, 524)
			if p == nil {
				panic("fixture allocation failed")
			}
			copy(unsafe.Slice((*uint16)(unsafe.Add(p, 12)), 256), text)
			C.nox_common_list_append_4258E0(head, (*C.nox_list_item_t)(p))
		}
	}
	files.RLock()
	beforeHandles := len(files.byHandle)
	files.RUnlock()
	snapshot := func(result uint8) PortTestRulesState {
		s := PortTestRulesState{Context: uint32(ruleLoaderContext), Result: result, LinksValid: true, GuardsValid: true, TableUnchanged: bytes.Equal(table, tableBefore)}
		copy(s.Settings[:], buf[8:68])
		for _, v := range append(append([]byte(nil), buf[:8]...), buf[68:]...) {
			s.GuardsValid = s.GuardsValid && v == 0xa5
		}
		if head != nil {
			prev := head
			for p, n := head.field_0, 0; p != head; p, n = p.field_0, n+1 {
				if p == nil || n > 4096 {
					s.LinksValid = false
					break
				}
				s.LinksValid = s.LinksValid && p.field_1 == prev && p.field_2 == nil
				text := unsafe.Slice((*uint16)(unsafe.Add(unsafe.Pointer(p), 12)), 256)
				s.Rejected = append(s.Rejected, alloc.GoString16S(text))
				prev = p
			}
			s.LinksValid = s.LinksValid && head.field_1 == prev && head.field_2 == head
		}
		files.RLock()
		s.HandlesUnchanged = len(files.byHandle) == beforeHandles
		files.RUnlock()
		return s
	}
	switch spec.Kind {
	case "write":
		name, free := alloc.CString(spec.User)
		var result uint8
		if spec.Wrapper {
			Sub_57AAA0(spec.User, settings, unsafe.Pointer(head))
		} else {
			result = uint8(C.sub_57AAA0((*C.char)(unsafe.Pointer(name)), (*C.char)(unsafe.Pointer(settings)), (*C.int)(unsafe.Pointer(head))))
		}
		free()
		out.Steps = append(out.Steps, snapshot(result))
		data, e := os.ReadFile(filepath.Join(spec.Dir, filepath.FromSlash(spec.Path)))
		out.Written, out.WriteExists = string(data), e == nil
	case "load":
		var user *C.char
		if !spec.UserNil {
			p, free := alloc.CString(spec.User)
			defer free()
			user = (*C.char)(unsafe.Pointer(p))
		}
		var result uint8
		if spec.Wrapper {
			Sub_57A1E0(settings, spec.User, unsafe.Pointer(head), int(spec.Selection), noxflags.GameFlag(spec.Flags))
		} else {
			result = uint8(C.sub_57A1E0((*C.int)(unsafe.Pointer(settings)), user, (*C.int)(unsafe.Pointer(head)), C.char(spec.Selection), C.short(spec.Flags)))
		}
		out.Steps = append(out.Steps, snapshot(result))
	case "file":
		result := uint8(ruleReadFile(ruleTestPath(spec.Path), settings, head, uint32(spec.Flags)))
		out.Steps = append(out.Steps, snapshot(result))
	case "lines":
		for _, line := range spec.Lines {
			wide := utf16.Encode([]rune(ruleTestPath(line)))
			ruleParseLine(wide, settings, head, uint32(spec.Flags))
			out.Steps = append(out.Steps, snapshot(0))
		}
	case "tokens":
		for _, tokens := range spec.Tokens {
			if len(tokens) == 0 || len(tokens) > 32 {
				return out, fmt.Errorf("invalid fixture token count")
			}
			wide := make([][]uint16, len(tokens))
			for i, token := range tokens {
				wide[i] = utf16.Encode([]rune(ruleTestPath(token)))
			}
			result := uint8(ruleApply(wide, settings, uint32(spec.Flags)))
			out.Steps = append(out.Steps, snapshot(result))
		}
	default:
		return out, fmt.Errorf("unknown rules fixture kind %q", spec.Kind)
	}
	out.FilesUnchanged = true
	for name, data := range spec.Files {
		got, e := os.ReadFile(filepath.Join(spec.Dir, filepath.FromSlash(name)))
		out.FilesUnchanged = out.FilesUnchanged && e == nil && bytes.Equal(got, []byte(data))
	}
	return out, nil
}

func ruleTestPath(s string) string {
	if i := strings.IndexByte(s, 0); i >= 0 {
		return s[:i]
	}
	return s
}
