//go:build porttest

package legacy

/*
#include <stdlib.h>
#include "GAME3_2.h"
#include "GAME5_2.h"
*/
import "C"

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"unsafe"

	"github.com/opennox/libs/ifs"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

// PortTestCommandRulesSpec selects one native entry point. Path is passed as
// raw C bytes; only Mode "path" accepts NilPath. Files are below Dir.
type PortTestCommandRulesSpec struct {
	Mode              string // header, file, path, map, wrapper
	Dir               string
	Path, Map         string
	NilPath           bool
	Headers           []string
	Files             map[string]string
	Flags             uint32
	CallbackReturns   []bool
	AfterCommandFlags []uint32
}

type PortTestCommandRulesResult struct {
	SettingsContextSame  bool
	FinalFlags           uint32
	HeaderValues         []uint32
	Result               int
	Commands             []string
	CallbackResults      []bool
	Files                map[string]string
	FilesUnchanged       bool
	HeaderTableUnchanged bool
	FilenameBlobsSame    bool
	HandlesUnchanged     bool
}

// PortTestCommandRules isolates the C command-rule readers. The recorder is
// installed below the existing Go parse-command export, so it observes exactly
// the widened command text passed by the C reader.
func PortTestCommandRules(spec PortTestCommandRulesSpec) (out PortTestCommandRulesResult, err error) {
	if spec.Dir == "" {
		return out, fmt.Errorf("missing fixture directory")
	}
	for name, data := range spec.Files {
		path := filepath.Join(spec.Dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			return out, err
		}
		if err := os.WriteFile(path, []byte(data), 0600); err != nil {
			return out, err
		}
	}
	snapshotFiles := func() (map[string]string, error) {
		got := make(map[string]string)
		err := filepath.WalkDir(spec.Dir, func(path string, d os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() {
				return nil
			}
			rel, err := filepath.Rel(spec.Dir, path)
			if err != nil {
				return err
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			got[filepath.ToSlash(rel)] = string(data)
			return nil
		})
		return got, err
	}
	before, err := snapshotFiles()
	if err != nil {
		return out, err
	}
	oldDir, err := ifs.Workdir()
	if err != nil {
		return out, err
	}
	if err = ifs.Chdir(spec.Dir); err != nil {
		return out, err
	}
	defer func() {
		if restoreErr := ifs.Chdir(oldDir); err == nil {
			err = restoreErr
		}
	}()
	restoreHandles := handles.PortTestInit()
	defer restoreHandles()
	restoreTable := portTestRuleTable()
	defer restoreTable()
	table := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 312208)), 56)
	tableBefore := append([]byte(nil), table...)
	userBlob := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 191748)), 9)
	fallbackBlob := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 191760)), 5)
	oldUser, oldFallback := append([]byte(nil), userBlob...), append([]byte(nil), fallbackBlob...)
	copy(userBlob, []byte("user.rul\x00"))
	copy(fallbackBlob, []byte(".rul\x00"))
	defer func() { copy(userBlob, oldUser); copy(fallbackBlob, oldFallback) }()
	filenameBefore := append(append([]byte(nil), userBlob...), fallbackBlob...)

	oldSettingsContext := ruleLoaderContext
	ruleLoaderContext = 0xdeadcafe
	defer func() { ruleLoaderContext = oldSettingsContext }()
	restoreFlags := noxflags.PortTestGameFlags(noxflags.GameFlag(spec.Flags))
	defer func() { restoreFlags() }()
	oldExec := ExecConsoleCmd
	defer func() { ExecConsoleCmd = oldExec }()
	call := 0
	ExecConsoleCmd = func(_ context.Context, cmd string) bool {
		out.Commands = append(out.Commands, cmd)
		result := true
		if call < len(spec.CallbackReturns) {
			result = spec.CallbackReturns[call]
		}
		out.CallbackResults = append(out.CallbackResults, result)
		if call < len(spec.AfterCommandFlags) {
			restoreFlags()
			restoreFlags = noxflags.PortTestGameFlags(noxflags.GameFlag(spec.AfterCommandFlags[call]))
		}
		call++
		return result
	}
	files.RLock()
	handleCount := len(files.byHandle)
	files.RUnlock()

	switch spec.Mode {
	case "header":
		for _, header := range spec.Headers {
			p := C.CString(header)
			out.HeaderValues = append(out.HeaderValues, uint32(C.sub_57AE30(p)))
			C.free(unsafe.Pointer(p))
		}
	case "file", "path":
		if spec.Mode == "file" && spec.NilPath {
			return out, fmt.Errorf("nil path is only valid for path mode")
		}
		var p *C.char
		if !spec.NilPath {
			p = C.CString(spec.Path)
			defer C.free(unsafe.Pointer(p))
		}
		if spec.Mode == "file" {
			out.Result = int(C.sub_4D0670(p))
		} else {
			out.Result = int(C.sub_4D0550(p))
		}
	case "map":
		p := C.CString(spec.Map)
		defer C.free(unsafe.Pointer(p))
		out.Result = int(C.sub_57A950(p))
	case "wrapper":
		Sub_4D0550(spec.Path)
	default:
		return out, fmt.Errorf("unknown command-rule mode %q", spec.Mode)
	}
	out.SettingsContextSame = ruleLoaderContext == 0xdeadcafe
	out.FinalFlags = uint32(noxflags.GetGame())
	out.Files, err = snapshotFiles()
	out.FilesUnchanged = err == nil && reflect.DeepEqual(before, out.Files)
	out.HeaderTableUnchanged = bytes.Equal(table, tableBefore)
	out.FilenameBlobsSame = bytes.Equal(append(append([]byte(nil), userBlob...), fallbackBlob...), filenameBefore)
	files.RLock()
	out.HandlesUnchanged = len(files.byHandle) == handleCount
	files.RUnlock()
	return out, err
}
