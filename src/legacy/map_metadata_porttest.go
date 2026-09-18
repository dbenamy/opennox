//go:build porttest

package legacy

/*
#include "GAME1_2.h"
*/
import "C"

import (
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"os"
	"path/filepath"
)

// Own only the file; callers supply real mapped state and object/render owners.
func PortTestMapMetadata(spec PortTestMapSectionIO, dir string) PortTestMapSectionWire {
	old := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	defer func() { cryptfile.Close(); cryptfile.SetGlobal(old) }()
	path := filepath.Join(dir, "metadata.bin")
	mode := cryptfile.WriteOnly
	if spec.Read {
		mode = cryptfile.ReadOnly
		if err := os.WriteFile(path, spec.Data, 0600); err != nil {
			panic(err)
		}
	}
	if err := cryptfile.OpenGlobal(path, mode, -1); err != nil {
		panic(err)
	}
	var ret uint32
	switch spec.Function {
	case "info":
		ret = uint32(C.nox_server_mapRWMapInfo_42A6E0())
	case "ambient":
		ret = uint32(C.nox_server_mapRWAmbientData_429200())
	case "toc":
		ret = uint32(C.nox_server_mapRWObjectTOC_428B30())
	default:
		panic(spec.Function)
	}
	pos, err := cryptfile.Global().File.Seek(0, 1)
	if err != nil {
		panic(err)
	}
	sum := cryptfile.Global().PortTestChecksum()
	if err := cryptfile.Close(); err != nil {
		panic(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return PortTestMapSectionWire{Data: data, Position: pos, Checksum: sum, Return: ret}
}
