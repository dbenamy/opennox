//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"unsafe"
)

type PortTestWorldTablesResult struct {
	Sine, Acos    []uint32
	ColdReturn    int8
	WarmReturns   []int8
	WarmUnchanged bool
	GuardsOK      bool
	Constants     [2]uint64
}

func PortTestWorldTables() (r PortTestWorldTablesResult) {
	restoreEnv := portTestWorldNumericEnv()
	defer restoreEnv()
	sine := unsafe.Slice(memmap.PtrUint32(0x85B3FC, 12260), 4096)
	acos := unsafe.Slice(memmap.PtrUint32(0x5D4594, 338472), 8192)
	flag := memmap.PtrUint8(0x5D4594, 371240)
	constants := memmap.Slice(0x581450, 7184)[:24]
	oldSine, oldAcos := append([]uint32(nil), sine...), append([]uint32(nil), acos...)
	oldFlag, oldConstants := *flag, append([]byte(nil), constants...)
	guards := []*uint32{memmap.PtrUint32(0x85B3FC, 12256), memmap.PtrUint32(0x85B3FC, 28644), memmap.PtrUint32(0x5D4594, 338468)}
	var oldGuards [3]uint32
	for i, p := range guards {
		oldGuards[i] = *p
		*p = 0x13572468
	}
	defer func() {
		copy(sine, oldSine)
		copy(acos, oldAcos)
		*flag = oldFlag
		copy(constants, oldConstants)
		for i, p := range guards {
			*p = oldGuards[i]
		}
	}()
	copy(constants, blobdata.PortTestWorldTrigConstants())
	r.Constants = [2]uint64{*memmap.PtrUint64(0x581450, 7184), *memmap.PtrUint64(0x581450, 7200)}
	clear(sine)
	clear(acos)
	*flag = 0
	r.ColdReturn = int8(worldInitTrigTables())
	r.Sine = append([]uint32(nil), sine...)
	r.Acos = append([]uint32(nil), acos...)
	r.WarmUnchanged = *flag == 1
	for _, f := range []byte{1, 2, 127, 128, 255} {
		*flag = f
		r.WarmReturns = append(r.WarmReturns, int8(worldInitTrigTables()))
		r.WarmUnchanged = r.WarmUnchanged && *flag == f
		for i, v := range sine {
			r.WarmUnchanged = r.WarmUnchanged && v == r.Sine[i]
		}
		for i, v := range acos {
			r.WarmUnchanged = r.WarmUnchanged && v == r.Acos[i]
		}
	}
	r.GuardsOK = true
	for _, p := range guards {
		r.GuardsOK = r.GuardsOK && *p == 0x13572468
	}
	return r
}
