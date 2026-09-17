//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME4.h"
#include "GAME4_1.h"
extern unsigned int dword_587000_230092;
*/
import "C"

import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

// Call the production Go entry points after retiring their C bridges.
func PortTestObjectReports(op int, a, b *server.Object, player, x, y int, pos *types.Pointf) int {
	switch op {
	case 0:
		return objectReportInit()
	case 1:
		return objectReportPhantom(player, b)
	case 2:
		return objectReportSimple(player, b)
	case 3:
		return objectReportMonster(player, b, x)
	case 4:
		return objectReportSprite(a, player, b)
	case 5:
		return objectReportPlayer(a, b, x, y)
	case 6:
		return objectReportRecipient(a, b)
	case 7:
		return objectReportSchedule(a.UpdateDataPlayer())
	case 8:
		return objectReportSoundLevel(*pos, b)
	case 9:
		objectReportRemoteAudio(a)
		return 0
	default:
		panic("object report operation")
	}
}
func PortTestObjectReportsDirectionThreshold() *int32 {
	return (*int32)(unsafe.Pointer(&C.dword_587000_230092))
}
