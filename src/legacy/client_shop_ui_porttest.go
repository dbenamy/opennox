//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME2_2.h"
#include "client__gui__guishop.h"
void sub_479680(void);
static uint32_t portTestShopUI(int op, uintptr_t* a) {
 switch(op) {
 case 0: return (uint32_t)(uintptr_t)sub_478030();
 case 1: return (uint32_t)(uintptr_t)sub_478040();
 case 2: return (uint32_t)(uintptr_t)sub_478080((int)a[0]);
 case 3: return (uint32_t)(uintptr_t)sub_4780A0((int)a[0]);
 case 4: return (uint32_t)(uintptr_t)sub_478110();
 case 5: return (uint32_t)(uintptr_t)sub_478480((int)a[0], (int)a[1], (int*)a[2], (int)a[3]);
 case 6: return (uint32_t)(uintptr_t)sub_478650((int)a[0], (int)a[1], (unsigned int)a[2]);
 case 7: sub_478850((int)a[0], (short)a[1], (int)a[2], (int)a[3]); return 0;
 case 8: return (uint32_t)(uintptr_t)sub_478970();
 case 9: return (uint32_t)(uintptr_t)sub_478A70((int2*)a[0]);
 case 10: return (uint32_t)(uintptr_t)sub_478C80();
 case 11: return (uint32_t)(uintptr_t)sub_478E50((int)a[0], (int)a[1], (unsigned int)a[2]);
 case 12: return (uint32_t)(uintptr_t)sub_478F10();
 case 13: return (uint32_t)(uintptr_t)sub_478F80();
 case 14: return (uint32_t)(uintptr_t)nox_xxx_getShopPic_4790F0((int)a[0]);
 case 15: sub_479280(); return 0;
 case 16: return (uint32_t)(uintptr_t)sub_479300((int)a[0], (int)a[1], (int)a[2], (short)a[3], (int)a[4]);
 case 17: return (uint32_t)(uintptr_t)sub_4793C0((int)a[0]);
 case 18: return (uint32_t)(uintptr_t)sub_479430();
 case 19: return (uint32_t)(uintptr_t)sub_479480((int)a[0]);
 case 20: return (uint32_t)(uintptr_t)sub_4794D0((int)a[0], (int)a[1]);
 case 21: return (uint32_t)(uintptr_t)sub_479590();
 case 22: sub_4795A0((int)a[0]); return 0;
 case 23: return (uint32_t)(uintptr_t)sub_479690((int)a[0], (short)a[1], (short)a[2], (int)a[3]);
 case 24: return (uint32_t)(uintptr_t)nox_client_tradeXxxSellAccept_4796D0((short)a[0]);
 case 25: return (uint32_t)(uintptr_t)sub_479700((short)a[0], (char)a[1]);
 case 26: sub_479810(); return 0;
 case 27: return (uint32_t)(uintptr_t)sub_479820((int)a[0], (short)a[1]);
 case 28: return (uint32_t)(uintptr_t)sub_479840((short)a[0]);
 case 29: return (uint32_t)(uintptr_t)sub_479870();
 case 30: return (uint32_t)(uintptr_t)sub_479880((uint32_t*)a[0]);
 case 31: return (uint32_t)(uintptr_t)sub_4798A0((uint32_t*)a[0]);
 case 32: sub_478730((int*)a[0]); return 0;
 case 33: nox_client_tradeXxxBuyAccept_478880((int)a[0], (short)a[1]); return 0;
 case 34: sub_4788F0((int)a[0], (int)a[1]); return 0;
 case 35: return (uint32_t)(uintptr_t)sub_478B10((int2*)a[0]);
 case 36: return (uint32_t)(uintptr_t)sub_478BC0((int*)a[0]);
 case 37: return (uint32_t)(uintptr_t)nox_xxx_cliStartShopDlg_478FD0((const wchar2_t*)a[0], (char*)a[1], (int)a[2]);
 case 38: sub_479520((int)a[0]); return 0;
 case 39: sub_479680(); return 0;
 case 40: return (uint32_t)(uintptr_t)sub_4795E0((int)a[0], (int)a[1]);
 case 41: sub_479740((int)a[0], (unsigned int)a[1]); return 0;
 default: return 0;
 }
}
*/
import "C"
import "unsafe"

// PortTestShopUI only dispatches to actual production C routines.
func PortTestShopUI(op int, args ...uintptr) uint32 {
	var a [10]uintptr
	copy(a[:], args)
	return uint32(C.portTestShopUI(C.int(op), (*C.uintptr_t)(unsafe.Pointer(&a[0]))))
}
