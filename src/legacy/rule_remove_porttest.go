//go:build porttest

package legacy

/*
#include <stdlib.h>
#include "GAME5_2.h"
*/
import "C"

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"unsafe"

	"github.com/opennox/libs/ifs"
)

// PortTestRuleRemoveSpec describes a standalone tree below Dir. Paths use '/'
// solely for fixture construction; Map and File are passed unmodified to C.
type PortTestRuleRemoveSpec struct {
	Dir       string
	Map, File string
	Dirs      []string
	Files     map[string]string
}

type PortTestRuleRemoveEntry struct {
	Path string
	Dir  bool
	Data string
}

type PortTestRuleRemoveResult struct {
	Result        int
	Before, After []PortTestRuleRemoveEntry
}

// PortTestRuleRemove calls the still-native C helper. It captures the complete
// test tree, including empty directories, before and after the call.
func PortTestRuleRemove(spec PortTestRuleRemoveSpec) (out PortTestRuleRemoveResult, err error) {
	if spec.Dir == "" {
		return out, fmt.Errorf("missing fixture directory")
	}
	for _, name := range spec.Dirs {
		if err := os.MkdirAll(filepath.Join(spec.Dir, filepath.FromSlash(name)), 0700); err != nil {
			return out, err
		}
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
	snapshot := func() ([]PortTestRuleRemoveEntry, error) {
		var entries []PortTestRuleRemoveEntry
		err := filepath.WalkDir(spec.Dir, func(path string, d os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if path == spec.Dir {
				return nil
			}
			rel, err := filepath.Rel(spec.Dir, path)
			if err != nil {
				return err
			}
			e := PortTestRuleRemoveEntry{Path: filepath.ToSlash(rel), Dir: d.IsDir()}
			if !e.Dir {
				data, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				e.Data = string(data)
			}
			entries = append(entries, e)
			return nil
		})
		sort.Slice(entries, func(i, j int) bool { return entries[i].Path < entries[j].Path })
		return entries, err
	}
	if out.Before, err = snapshot(); err != nil {
		return out, err
	}
	old, err := ifs.Workdir()
	if err != nil {
		return out, err
	}
	if err = ifs.Chdir(spec.Dir); err != nil {
		return out, err
	}
	defer func() {
		if restoreErr := ifs.Chdir(old); err == nil {
			err = restoreErr
		}
	}()
	mapName := C.CString(spec.Map)
	fileName := C.CString(spec.File)
	defer C.free(unsafe.Pointer(mapName))
	defer C.free(unsafe.Pointer(fileName))
	out.Result = int(C.sub_57A9F0(mapName, fileName))
	out.After, err = snapshot()
	return out, err
}
