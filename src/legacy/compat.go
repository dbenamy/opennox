package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"time"

	"github.com/opennox/libs/env"
	"github.com/opennox/libs/platform"
)

//export noxGetLocalTime
func noxGetLocalTime(p *C.noxSYSTEMTIME) {
	tm := time.Now()
	if env.IsE2E() {
		tm = time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC).Add(platform.Ticks())
	}
	p.wYear = C.ushort(tm.Year())
	p.wMonth = C.ushort(tm.Month())
	p.wDayOfWeek = C.ushort(tm.Weekday())
	p.wDay = C.ushort(tm.Day())
	p.wHour = C.ushort(tm.Hour())
	p.wMinute = C.ushort(tm.Minute())
	p.wSecond = C.ushort(tm.Second())
	p.wMilliseconds = C.ushort(tm.Nanosecond() / 1e6)
}
