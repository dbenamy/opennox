//go:build porttest

package legacy

/*
#include "GAME1.h"
#include "GAME3_2.h"
int nox_xxx_netSendChat_528AC0(nox_object_t*, wchar2_t*, wchar2_t);
static uint32_t gameplayTextInvoke(int op, uint32_t* a) {
 switch(op) {
 case 0: return nox_xxx_netSendLineMessage_4D9EB0(a[0], (wchar2_t*)a[1], a[2], a[3]);
 case 1: return nox_xxx_printToAll_4D9FD0(a[0], (wchar2_t*)a[1], a[2], a[3]);
 case 2: return nox_xxx_netInformTextMsg_4DA0F0(a[0],a[1],(int*)a[2]);
 case 3: return nox_xxx_netInformTextMsg2_4DA180(a[0],(uint8_t*)a[1]);
 case 4: nox_xxx_netPriMsgToPlayer_4DA2C0((nox_object_t*)a[0],(const char*)a[1],a[2]); return 0;
 case 5: return nox_xxx_netPrintLineToAll_4DA390((const char*)a[0]);
 case 6: return (uint32_t)nox_xxx_getFirstPlayerUnit_4DA7C0();
 case 7: return (uint32_t)nox_xxx_getNextPlayerUnit_4DA7F0((const nox_object_t*)a[0]);
 case 8: return nox_xxx_cliCanTalkMB_4100F0((short*)a[0]);
 case 9: return nox_xxx_netSendChat_528AC0((nox_object_t*)a[0],(wchar2_t*)a[1],a[2]);
 default: return 0xDEADBEEF;
 }
}
*/
import "C"

import "unsafe"

func gameplayTextInvoke(op int, args [5]uint32) uint32 {
	return uint32(C.gameplayTextInvoke(C.int(op), (*C.uint32_t)(unsafe.Pointer(&args[0]))))
}
