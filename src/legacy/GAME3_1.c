#include <math.h>

#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME1_2.h"
#include "GAME1_3.h"
#include "GAME2.h"
#include "GAME2_1.h"
#include "GAME2_2.h"
#include "GAME2_3.h"
#include "GAME3.h"
#include "GAME3_1.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME5_2.h"
#include "client__draw__animdraw.h"
#include "client__drawable__drawable.h"
#include "client__gui__guisumn.h"
#include "client__gui__servopts__advserv.h"
#include "common__net_list.h"
#include "common__system__team.h"

#include "client__gui__guiinput.h"
#include "client__gui__servopts__objlst.h"
#include "client__gui__servopts__spelllst.h"
#include "client__gui__tooltip.h"
#include "client__gui__window.h"
#include "client__shell__inputcfg__inputcfg.h"
#include "client__shell__mainmenu.h"

#include "client__draw__plasma.h"
#include "client__drawable__update__charmup.h"
#include "client__drawable__update__fireball.h"
#include "client__system__ctrlevnt.h"

#include "client__shell__optsback.h"
#include "client__video__draw_common.h"
#include "common/fs/nox_fs.h"
#include "common__strman.h"
#include "input_common.h"
#include "operators.h"

extern uint32_t dword_8531A0_2572;
extern uint32_t dword_8531A0_2576;
extern uint32_t dword_587000_180476;
extern uint32_t dword_5d4594_1319260;
extern uint32_t dword_5d4594_1321196;
extern uint32_t dword_5d4594_1313880;
extern uint32_t dword_5d4594_1320944;
extern uint32_t dword_5d4594_1320948;
extern uint32_t dword_587000_183460;
extern uint32_t dword_5d4594_1316412;
extern uint32_t dword_5d4594_1522968;
extern uint32_t dword_5d4594_1321520;
extern uint32_t dword_5d4594_1319268;
extern uint32_t dword_5d4594_1321024;
extern uint32_t dword_5d4594_1320936;
extern uint32_t dword_587000_183456;
extern uint32_t dword_5d4594_1320988;
extern uint32_t dword_5d4594_1319248;
extern uint32_t dword_587000_180480;
extern uint32_t dword_5d4594_1319236;
extern uint32_t dword_5d4594_1320972;
extern uint32_t dword_5d4594_1316704;
extern uint32_t dword_5d4594_1321208;
extern uint32_t dword_5d4594_1321800;
extern uint32_t dword_5d4594_1321224;
extern uint32_t dword_5d4594_1319056;
extern uint32_t dword_5d4594_1320932;
extern uint32_t dword_5d4594_1321032;
extern uint32_t dword_5d4594_1319264;
extern uint32_t dword_5d4594_1321044;
extern uint32_t dword_5d4594_1319232;
extern uint32_t dword_5d4594_1321216;
extern nox_window* dword_5d4594_1321236;
extern nox_window* dword_5d4594_1321240;
extern nox_window* dword_5d4594_1321248;
extern nox_window* dword_5d4594_1321244;
nox_window* dword_5d4594_1522616 = 0;
nox_window* dword_5d4594_1522620 = 0;
nox_window* dword_5d4594_1522624 = 0;
nox_window* dword_5d4594_1522628 = 0;
extern uint32_t dword_5d4594_1320992;
extern uint32_t dword_5d4594_1321204;
extern uint32_t dword_5d4594_1316408;
extern uint64_t qword_581450_9512;
extern uint64_t qword_581450_9544;
extern uint32_t dword_5d4594_1320968;
extern uint32_t dword_5d4594_1319060;
extern uint32_t dword_5d4594_1522632;
extern uint32_t dword_5d4594_1321252;
extern uint32_t dword_5d4594_1522612;
extern uint32_t dword_5d4594_1522604;
extern uint32_t dword_5d4594_1321232;
extern uint32_t dword_5d4594_1316712;
extern uint32_t dword_5d4594_1320964;
extern uint32_t dword_5d4594_1316708;
extern uint32_t dword_5d4594_1321228;
extern uint32_t dword_5d4594_1321040;
extern uint32_t dword_5d4594_1316972;
extern uint32_t dword_5d4594_1320940;
extern uint32_t nox_player_netCode_85319C;
extern int nox_win_width;
extern int nox_win_height;

extern uint32_t nox_color_white_2523948;
extern uint32_t nox_color_blue_2650684;
extern uint32_t nox_color_yellow_2589772;
extern uint32_t nox_color_violet_2598268;
extern uint32_t nox_color_black_2650656;

nox_gui_animation* nox_wnd_xxx_1522608 = 0;

void* nox_gui_itemAmount_item_1319256 = 0;
void* nox_gui_itemAmount_dialog_1319228 = 0;

//----- (004BD280) --------------------------------------------------------
uint32_t* sub_4BD280(int a1, int a2) {
	int v2;           // esi
	uint32_t* result; // eax
	uint32_t* v4;     // ecx
	int v5;           // edi

	v2 = a2 + 4;
	result = calloc(1, a1 * (a2 + 4) + 4);
	if (result) {
		v4 = result + 1;
		*result = result + 1;
		if (a1 != 1) {
			v5 = a1 - 1;
			do {
				--v5;
				*v4 = (char*)v4 + v2;
				v4 = (uint32_t*)((char*)v4 + v2);
			} while (v5);
		}
		*v4 = 0;
	}
	return result;
}

//----- (004BD2D0) --------------------------------------------------------
void sub_4BD2D0(void* lpMem) { free(lpMem); }

//----- (004BD2E0) --------------------------------------------------------
uint32_t* sub_4BD2E0(uint32_t** a1) {
	uint32_t* result; // eax
	uint32_t* v2;     // edx

	result = *a1;
	if (*a1) {
		v2 = (uint32_t*)*result;
		++result;
		*a1 = v2;
	}
	return result;
}

//----- (004BD300) --------------------------------------------------------
int sub_4BD300(uint32_t* a1, int a2) {
	int result; // eax

	result = a2 - 4;
	*(uint32_t*)(a2 - 4) = *a1;
	*a1 = a2 - 4;
	return result;
}

//----- (004BD340) --------------------------------------------------------
uint32_t* sub_4BD340(int a1, int a2, int a3, int a4) {
	uint32_t* v4; // esi

	v4 = calloc(1, 0x1Cu);
	memset(v4, 0, 0x1Cu);
	*v4 = a1;
	v4[6] = a4;
	v4[1] = sub_4BD280(a2 / (a4 + 24), a4 + 24);
	v4[2] = sub_4BD280(a3, 84);
	nox_common_list_clear_425760(v4 + 3);
	if (v4[1] && v4[2]) {
		return v4;
	}
	sub_4BD3C0(v4);
	return 0;
}

//----- (004BD3C0) --------------------------------------------------------
void sub_4BD3C0(void* lpMem) {
	int i; // eax

	for (i = nox_common_list_getNext_425940((int*)lpMem + 3); i; i = nox_common_list_getNext_425940((int*)lpMem + 3)) {
		sub_4BD690(i);
	}
	if (*((uint32_t*)lpMem + 1)) {
		sub_4BD2D0(*((void**)lpMem + 1));
	}
	if (*((uint32_t*)lpMem + 2)) {
		sub_4BD2D0(*((void**)lpMem + 2));
	}
	free(lpMem);
}

//----- (004BD420) --------------------------------------------------------
uint32_t* sub_4BD420(int a1, int a2) {
	uint32_t* result; // eax

	result = *(uint32_t**)(a1 + 12);
	if (result == (uint32_t*)(a1 + 12)) {
		return 0;
	}
	while (result[4] != a2 || !result[5]) {
		result = (uint32_t*)*result;
		if (result == (uint32_t*)(a1 + 12)) {
			return 0;
		}
	}
	return result;
}

//----- (004BD470) --------------------------------------------------------
uint32_t* sub_4BD470(uint32_t** a1, int a2) {
	uint32_t* v2;  // eax
	uint32_t* v3;  // edi
	uint32_t* v5;  // ebx
	uint32_t* v6;  // eax
	uint32_t* v7;  // ebp
	char* v8;      // edi
	signed int v9; // eax
	uint32_t* v10; // [esp+18h] [ebp+8h]

	v2 = sub_4BD420((int)a1, a2);
	v3 = v2;
	if (v2) {
		nox_common_list_remove_425920((uint32_t**)v2);
		sub_425900(a1 + 3, v3);
		return v3;
	}
	if (!sub_486B60((int)*a1, a2)) {
		return 0;
	}
	v5 = sub_4BD2E0((uint32_t**)a1[2]);
	if (!v5) {
		sub_4BD600((int)a1);
		v5 = sub_4BD2E0((uint32_t**)a1[2]);
		if (!v5) {
			sub_486E00((int)*a1);
			return 0;
		}
	}
	v5[4] = a2;
	v5[13] = a1;
	sub_425770(v5);
	v5[3] = 0;
	sub_487C30(v5 + 6);
	nullsub_10(v5 + 14);
	v5[11] = v5 + 14;
	v6 = (uint32_t*)(*a1)[71];
	v10 = (uint32_t*)(*a1)[71];
	if (!v6) {
		sub_486AA0(*a1, v5[4], v5 + 14);
		sub_425900(a1 + 3, v5);
		v5[5] = 1;
		sub_486E00((int)*a1);
		return v5;
	}
	while (1) {
		v7 = a1[6];
		if ((int)v7 > (int)v6) {
			v7 = v6;
		}
		v8 = (char*)sub_4BD2E0((uint32_t**)a1[1]);
		if (!v8) {
			while (sub_4BD600((int)a1)) {
				v8 = (char*)sub_4BD2E0((uint32_t**)a1[1]);
				if (v8) {
					goto LABEL_17;
				}
			}
			sub_4BD690((int)v5);
			return 0;
		}
	LABEL_17:
		sub_487D30(v8, (int)(v8 + 24), (int)v7);
		sub_487C50((int)(v5 + 6), v8);
		v9 = sub_486DB0((int)*a1, v8 + 24, (signed int)v7);
		if ((uint32_t*)v9 != v7) {
			sub_4BD690((int)v5);
			return 0;
		}
		v10 = (uint32_t*)((char*)v10 - v9);
		if (!v10) {
			sub_486AA0(*a1, v5[4], v5 + 14);
			sub_425900(a1 + 3, v5);
			v5[5] = 1;
			sub_486E00((int)*a1);
			return v5;
		}
		v6 = v10;
	}
}
// 487CF0: using guessed type void  nullsub_10(uint32_t);

//----- (004BD600) --------------------------------------------------------
int sub_4BD600(int a1) {
	int v1; // esi

	v1 = sub_425960(a1 + 12);
	if (!v1) {
		return 0;
	}
	while (sub_4BD680(v1)) {
		v1 = sub_425960(v1);
		if (!v1) {
			return 0;
		}
	}
	sub_4BD690(v1);
	return 1;
}

//----- (004BD650) --------------------------------------------------------
int sub_4BD650(int a1) {
	int result; // eax

	result = a1;
	++*(uint32_t*)(a1 + 12);
	return result;
}

//----- (004BD660) --------------------------------------------------------
int sub_4BD660(int a1) {
	int result; // eax

	result = *(uint32_t*)(a1 + 12) - 1;
	*(uint32_t*)(a1 + 12) = result;
	if (result < 0) {
		*(uint32_t*)(a1 + 12) = 0;
	}
	return result;
}

//----- (004BD680) --------------------------------------------------------
int sub_4BD680(int a1) { return *(uint32_t*)(a1 + 12); }

//----- (004BD690) --------------------------------------------------------
int sub_4BD690(int a1) {
	uint32_t** i; // esi

	if (*(uint32_t*)(a1 + 4) != a1) {
		nox_common_list_remove_425920((uint32_t**)a1);
	}
	for (i = (uint32_t**)nox_common_list_getNext_425940((int*)(a1 + 32)); i;
		 i = (uint32_t**)nox_common_list_getNext_425940((int*)(a1 + 32))) {
		nox_common_list_remove_425920(i);
		sub_487D60((int)i);
		sub_4BD300(*(uint32_t**)(*(uint32_t*)(a1 + 52) + 4), (int)i);
	}
	nullsub_9(a1 + 24);
	return sub_4BD300(*(uint32_t**)(*(uint32_t*)(a1 + 52) + 8), a1);
}
// 487CA0: using guessed type void  nullsub_9(uint32_t);

//----- (004BD710) --------------------------------------------------------
int sub_4BD710(int a1) { return a1 + 24; }

//----- (004BD720) --------------------------------------------------------
uint32_t* sub_4BD720(int a1) {
	uint32_t* v1; // esi

	v1 = calloc(1, 0x138u);
	memset(v1, 0, 0x138u);
	sub_425770(v1);
	sub_4BDC00((int)(v1 + 30));
	sub_4864A0(v1 + 44);
	sub_4BD7C0(v1);
	v1[33] = a1;
	v1[43] = *(uint32_t*)(a1 + 256);
	if (!(*(int (**)(uint32_t*))(*(uint32_t*)(a1 + 256) + 4))(v1)) {
		return v1;
	}
	if (v1) {
		sub_4BD7A0(v1);
	}
	return 0;
}

//----- (004BD7A0) --------------------------------------------------------
void sub_4BD7A0(void* lpMem) {
	(*(void (**)(void*))(*((uint32_t*)lpMem + 43) + 8))(lpMem);
	free(lpMem);
}

//----- (004BD7C0) --------------------------------------------------------
uint32_t* sub_4BD7C0(uint32_t* a1) {
	uint32_t* result; // eax

	a1[69] = sub_4BD8C0;
	a1[70] = sub_4BD940;
	a1[71] = sub_4BD9B0;
	a1[34] = 0;
	a1[35] = 0;
	a1[36] = 0;
	a1[38] = 0;
	a1[3] = 1;
	sub_4BDC00((int)(a1 + 30));
	a1[30] = 0;
	a1[29] = *getMemU32Ptr(0x5D4594, 1193340);
	a1[28] = 0;
	result = sub_4864A0(a1 + 4);
	a1[72] = 0;
	return result;
}

//----- (004BD840) --------------------------------------------------------
void sub_4BD840(int a3) {
	int v1;           // ebp
	unsigned int* v2; // edi
	uint32_t* v3;     // esi
	int result;       // eax

	v1 = *(uint32_t*)(a3 + 132);
	v2 = (unsigned int*)(a3 + 176);
	sub_4864A0((uint32_t*)(a3 + 176));
	sub_486570((unsigned int*)(a3 + 176), (uint32_t*)(a3 + 16));
	sub_486620((uint32_t*)(a3 + 16));
	if (*(uint32_t*)(a3 + 112)) {
		sub_486570(v2, *(uint32_t**)(a3 + 112));
		sub_486620(*(uint32_t**)(a3 + 112));
	}
	v3 = *(uint32_t**)(a3 + 116);
	if (v3) {
		sub_486570(v2, v3);
	}
	sub_486570(v2, (uint32_t*)(v1 + 88));
	result = *(uint32_t*)(v1 + 184);
	if (result) {
		 sub_486570(v2, *(uint32_t**)(v1 + 184));
	}
}

//----- (004BD8C0) --------------------------------------------------------
int sub_4BD8C0(int a1) {
	int (*v1)(int); // eax
	int result;     // eax
	int v3;         // eax
	int v4;         // eax

	v1 = *(int (**)(int))(a1 + 136);
	if (v1) {
		result = v1(a1);
		if (result) {
			*(uint32_t*)(a1 + 300) = 0;
			*(uint32_t*)(a1 + 304) = 0;
			*(uint32_t*)(a1 + 296) = 0;
			return result;
		}
	} else {
		if (*(uint32_t*)(a1 + 292)) {
			v3 = nox_common_list_getNext_425940(*(int**)(a1 + 292));
			*(uint32_t*)(a1 + 292) = v3;
			if (v3) {
				*(uint32_t*)(a1 + 296) = *(uint32_t*)(v3 + 12);
				v4 = *(uint32_t*)(v3 + 16);
				*(uint32_t*)(a1 + 300) = v4;
				*(uint32_t*)(a1 + 304) = v4;
				return 0;
			}
		}
		*(uint32_t*)(a1 + 300) = 0;
	}
	return 0;
}

//----- (004BD940) --------------------------------------------------------
int sub_4BD940(int a1) {
	void (*v1)(int); // eax

	if (*(uint32_t*)(a1 + 128)) {
		if (*(int*)(a1 + 128) != -1) {
			--*(uint32_t*)(a1 + 128);
		}
		sub_4BDB90((uint32_t*)a1, *(uint32_t**)(a1 + 288));
	} else {
		sub_4BDB90((uint32_t*)a1, 0);
	}
	v1 = *(void (**)(int))(a1 + 140);
	if (v1) {
		v1(a1);
	}
	if (*(uint32_t*)(a1 + 288)) {
		(*(void (**)(int))(*(uint32_t*)(a1 + 172) + 36))(a1);
	}
	return 0;
}

//----- (004BD9B0) --------------------------------------------------------
int sub_4BD9B0(uint32_t* a2) {
	int v1;               // eax
	int (*v2)(uint32_t*); // eax
	int result;           // eax

	a2[72] = 0;
	v1 = a2[31];
	LOBYTE(v1) = v1 & 0xFA;
	a2[31] = v1;
	a2[32] = 0;
	sub_4864A0(a2 + 4);
	v2 = (int (*)(uint32_t*))a2[36];
	if (v2) {
		result = v2(a2);
	} else {
		result = 0;
	}
	return result;
}

//----- (004BDA60) --------------------------------------------------------
void sub_4BDA60(void* lpMem) {
	sub_4BDA80((int)lpMem);
	sub_486E90((int)lpMem);
	sub_4BD7A0(lpMem);
}

//----- (004BDA80) --------------------------------------------------------
int sub_4BDA80(int a1) {
	int result = 0;      // eax
	int (*result2)(int); // eax

	if (*(uint8_t*)(a1 + 124) & 5) {
		(*(void (**)(int))(*(uint32_t*)(a1 + 172) + 16))(a1);
	}
	result2 = *(int (**)(int))(a1 + 148);
	if (result2) {
		result = (int)result2(a1);
	}
	*(uint32_t*)(a1 + 288) = 0;
	return result;
}

//----- (004BDB20) --------------------------------------------------------
int sub_4BDB20(int a1) {
	int result; // eax

	result = a1;
	*(uint32_t*)(a1 + 124) |= 0x10u;
	return result;
}

//----- (004BDB30) --------------------------------------------------------
int sub_4BDB30(int a1) {
	int result; // eax

	result = a1;
	*(uint32_t*)(a1 + 124) &= 0xFFFFFFEF;
	return result;
}

//----- (004BDB40) --------------------------------------------------------
int sub_4BDB40(int a2) {
	int result; // eax

	if (*(uint8_t*)(a2 + 124) & 5) {
		return -2146500608;
	}
	if (!*(uint32_t*)(a2 + 288)) {
		return -2147024896;
	}
	sub_486520((unsigned int*)(a2 + 16));
	sub_4BD840(a2);
	result = (*(int (**)(int))(*(uint32_t*)(a2 + 172) + 12))(a2);
	if (!result) {
		*(uint32_t*)(a2 + 124) |= 1u;
	}
	return result;
}

//----- (004BDB90) --------------------------------------------------------
void sub_4BDB90(uint32_t* a1, uint32_t* a2) {
	int v2;       // eax
	int v3;       // eax
	uint32_t* v4; // eax
	int v5;       // eax

	a1[72] = a2;
	if (a2) {
		v2 = sub_487C80((int)a2);
		a1[73] = v2;
		if (v2) {
			a1[74] = *(uint32_t*)(v2 + 12);
			v3 = *(uint32_t*)(v2 + 16);
			a1[75] = v3;
			a1[76] = v3;
			*a2 = 0;
		} else {
			v4 = (uint32_t*)a1[72];
			a1[74] = *v4;
			v5 = v4[1];
			a1[75] = v5;
			a1[76] = v5;
		}
	}
}

//----- (004BDC00) --------------------------------------------------------
int sub_4BDC00(int a1) {
	int result; // eax

	result = a1;
	*(uint32_t*)(a1 + 8) = 0;
	*(uint32_t*)(a1 + 4) = 0;
	return result;
}

//----- (004BDC10) --------------------------------------------------------
int nox_xxx_loadAdvancedWnd_4BDC10(int* a1) {
	dword_5d4594_1316708 = nox_new_window_from_file("advanced.wnd", nox_xxx_windowAdvancedServProc_4BDDB0);
	if (!dword_5d4594_1316708) {
		return 0;
	}
	sub_46B120(*(uint32_t**)&dword_5d4594_1316708, 0);
	sub_46C690(*(int*)&dword_5d4594_1316708);
	nox_xxx_wndShowModalMB_46A8C0(*(int*)&dword_5d4594_1316708);
	nox_xxx_wndSetWindowProc_46B300(*(int*)&dword_5d4594_1316708, sub_4BDDA0);
	nox_xxx_windowFocus_46B500(*(int*)&dword_5d4594_1316708);
	return sub_4BDC70(a1);
}

//----- (004BDC70) --------------------------------------------------------
int sub_4BDC70(int* a1) {
	uint32_t* v1; // eax
	uint32_t* v2; // eax
	uint32_t* v3; // eax

	if (nox_common_gameFlags_check_40A5C0(1)) {
		v1 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1316708, 10167);
		v1[9] |= 4u;
		dword_5d4594_1316704 = 0;
	} else {
		v2 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1316708, 10164);
		v2[9] |= 4u;
		v3 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1316708, 10167);
		nox_xxx_wnd_46ABB0((int)v3, 0);
		dword_5d4594_1316704 = 1;
	}
	sub_453F70(a1 + 6);
	sub_4535E0(a1 + 11);
	sub_4535F0(a1[12]);
	return sub_4BDD10();
}

//----- (004BDD10) --------------------------------------------------------
int sub_4BDD10() {
	int v1;   // eax
	char* v2; // eax

	switch (dword_5d4594_1316704) {
	case 0:
		v2 = sub_4165B0();
		v1 = sub_4CEBA0(*(int*)&dword_5d4594_1316708, v2);
		dword_5d4594_1316712 = v1;
		break;
	case 1:
		dword_5d4594_1316712 = nox_xxx_guiSpelllistLoad_453850(*(int*)&dword_5d4594_1316708);
		return nox_xxx_windowFocus_46B500(*(int*)&dword_5d4594_1316712);
	case 2:
		v1 = nox_xxx_guiObjlistLoad_4530C0(*(int*)&dword_5d4594_1316708, 0x1000000);
		dword_5d4594_1316712 = v1;
		break;
	case 3:
		v1 = nox_xxx_guiObjlistLoad_4530C0(*(int*)&dword_5d4594_1316708, 0x2000000);
		dword_5d4594_1316712 = v1;
		break;
	default:
		return nox_xxx_windowFocus_46B500(*(int*)&dword_5d4594_1316712);
	}
	return nox_xxx_windowFocus_46B500(*(int*)&dword_5d4594_1316712);
}

//----- (004BDDA0) --------------------------------------------------------
int sub_4BDDA0() { return 1; }

//----- (004BDDB0) --------------------------------------------------------
int nox_xxx_windowAdvancedServProc_4BDDB0(int a1, int a2, int* a3, int a4) {
	int v3;     // esi
	int result; // eax
	char* v5;   // ebx

	if (a2 == 23) {
		return 1;
	}
	if (a2 != 16391) {
		return 1;
	}
	v3 = nox_xxx_wndGetID_46B0A0(a3);
	nox_xxx_clientPlaySoundSpecial_452D80(766, 100);
	switch (v3) {
	case 10148:
		v5 = sub_4165B0();
		memcpy(v5 + 24, sub_453F90(), 0x14u);
		*((uint32_t*)v5 + 11) = *(uint32_t*)sub_453600();
		*((uint32_t*)v5 + 12) = sub_453610();
		sub_4BDF30();
		return 1;
	case 10164:
		if (dword_5d4594_1316712) {
			nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1316712);
			dword_5d4594_1316712 = 0;
		}
		dword_5d4594_1316704 = 1;
		sub_4BDD10();
		result = 1;
		break;
	case 10165:
		if (dword_5d4594_1316712) {
			nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1316712);
			dword_5d4594_1316712 = 0;
		}
		dword_5d4594_1316704 = 2;
		sub_4BDD10();
		result = 1;
		break;
	case 10166:
		if (dword_5d4594_1316712) {
			nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1316712);
			dword_5d4594_1316712 = 0;
		}
		dword_5d4594_1316704 = 3;
		sub_4BDD10();
		result = 1;
		break;
	case 10167:
		if (dword_5d4594_1316712) {
			nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1316712);
			dword_5d4594_1316712 = 0;
		}
		dword_5d4594_1316704 = 0;
		sub_4BDD10();
		result = 1;
		break;
	default:
		return 1;
	}
	return result;
}

//----- (004BDF30) --------------------------------------------------------
int sub_4BDF30() {
	int result; // eax

	result = dword_5d4594_1316708;
	if (dword_5d4594_1316708) {
		nox_xxx_wnd_46C6E0(*(int*)&dword_5d4594_1316708);
		nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1316708);
		dword_5d4594_1316708 = 0;
		dword_5d4594_1316712 = 0;
		result = nox_xxx_windowFocus_46B500(0);
	}
	return result;
}

//----- (004BDF70) --------------------------------------------------------
int sub_4BDF70(int* a1) {
	int result; // eax

	result = *(int*)&dword_5d4594_1316708;
	if (dword_5d4594_1316708) {
		result = sub_4BDF90(a1);
	}
	return result;
}

//----- (004BDF90) --------------------------------------------------------
int sub_4BDF90(int* a1) {
	int result = 0;       // eax
	int (*result2)(void); // eax

	sub_453F70(a1 + 6);
	sub_4535E0(a1 + 11);
	sub_4535F0(a1[12]);
	result2 = *(int (**)(void))getMemAt(0x587000, 180016 + 4 * dword_5d4594_1316704);
	if (result2) {
		result = result2();
	}
	return result;
}

//----- (004BDFD0) --------------------------------------------------------
int sub_4BDFD0() {
	char* v0;     // eax
	int v1;       // esi
	int v2;       // edi
	uint32_t* v3; // eax
	uint32_t* v4; // eax
	int v5;       // eax
	uint32_t* v6; // esi
	uint32_t* v7; // eax
	int v9;       // [esp+8h] [ebp-8h]
	int v10;      // [esp+Ch] [ebp-4h]

	v0 = sub_416640();
	v1 = nox_strman_get_lang_code();
	v2 = (int)v0;
	if (nox_xxx_guiFontHeightMB_43F320(0) > 10) {
		v1 = 2;
	}
	if (0) {
		v3 = nox_new_window_from_file(*(const char**)getMemAt(0x587000, 180088 + 4 * v1), sub_4BE330);
	} else {
		v3 = nox_new_window_from_file(*(const char**)getMemAt(0x587000, 180048 + 4 * v1), sub_4BE330);
	}
	dword_5d4594_1316972 = v3;
	if (!dword_5d4594_1316972) {
		return 0;
	}
	sub_46B120(v3, 0);
	sub_46C690(*(int*)&dword_5d4594_1316972);
	nox_xxx_wndShowModalMB_46A8C0(*(int*)&dword_5d4594_1316972);
	nox_xxx_windowFocus_46B500(*(int*)&dword_5d4594_1316972);
	nox_xxx_wndSetWindowProc_46B300(*(int*)&dword_5d4594_1316972, sub_4BE320);
	v4 = nox_xxx_wndGetChildByID_46B0C0(0, 10100);
	nox_gui_getWindowOffs_46AA20((int)v4, &v10, &v9);
	if (0) {
		v5 = v9 + 55;
	} else {
		v5 = v9 + 80;
	}
	nox_window_setPos_46A9B0(*(uint32_t**)&dword_5d4594_1316972, v10 + 15, v5);
	v6 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1316972, 2104);
	nox_xxx_wnd_46B280((int)v6, *(int*)&dword_5d4594_1316972);
	nox_xxx_wndSetProc_46B2C0((int)v6, sub_4BE330);
	if (0) {
		v7 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1316972, 2119);
		nox_xxx_wnd_46B280((int)v7, *(int*)&dword_5d4594_1316972);
	}
	return sub_4BE120(v2);
}

//----- (004BE120) --------------------------------------------------------
int sub_4BE120(int a1) {
	uint32_t* v1;    // eax
	int v2;          // ecx
	unsigned int v3; // ecx
	uint32_t* v4;    // eax
	int v5;          // ecx
	unsigned int v6; // ecx
	uint32_t* v7;    // esi
	uint32_t* v8;    // eax
	uint32_t* v9;    // eax
	uint32_t* v10;   // eax
	uint32_t* v11;   // eax
	uint32_t* v12;   // eax
	int result;      // eax
	uint32_t* v14;   // eax
	int v15;         // [esp-14h] [ebp-40h]
	wchar2_t v16[16]; // [esp+Ch] [ebp-20h]

	v1 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1316972, 2102);
	v2 = v1[9];
	if (*(uint32_t*)(a1 + 58)) {
		v3 = v2 | 4;
	} else {
		v3 = v2 & 0xFFFFFFFB;
	}
	v1[9] = v3;
	v4 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1316972, 2103);
	v5 = v4[9];
	if (*(uint32_t*)(a1 + 62)) {
		v6 = v5 | 4;
	} else {
		v6 = v5 & 0xFFFFFFFB;
	}
	v4[9] = v6;
	v7 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1316972, 2110);
	nox_swprintf(v16, L"%d", *(uint32_t*)(a1 + 70));
	nox_window_call_field_94((int)v7, 16414, (int)v16, -1);
	switch (*(uint32_t*)(a1 + 66)) {
	case 0:
		v8 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1316972, 2110);
		nox_xxx_wnd_46ABB0((int)v8, 0);
		v15 = 2106;
		v10 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1316972, v15);
		v10[9] |= 4u;
		break;
	case 1:
		v9 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1316972, 2110);
		nox_xxx_wnd_46ABB0((int)v9, 0);
		v10 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1316972, 2107);
		v10[9] |= 4u;
		break;
	case 2:
		v11 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1316972, 2110);
		nox_xxx_wnd_46ABB0((int)v11, 0);
		v10 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1316972, 2108);
		v10[9] |= 4u;
		break;
	case 3:
		v12 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1316972, 2110);
		nox_xxx_wnd_46ABB0((int)v12, 1);
		v15 = 2109;
		v10 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1316972, v15);
		v10[9] |= 4u;
		break;
	default:
		break;
	}
	if (0) {
		v14 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1316972, 2119);
		nox_window_call_field_94((int)v14, 16394, *(uint32_t*)(a1 + 74), 0);
		result = sub_4BE2C0(*(uint32_t*)(a1 + 74));
	}
	return 0;
}

//----- (004BE320) --------------------------------------------------------
int sub_4BE320() { return 1; }

//----- (004BE330) --------------------------------------------------------
int sub_4BE330(int a1, unsigned int a2, int* a3, int a4) {
	uint32_t* v4;       // esi
	char* v5;           // edi
	int result;         // eax
	const wchar2_t* v7;  // eax
	int v8;             // eax
	int v9;             // esi
	char* v10;          // eax
	char* v11;          // eax
	char* v12;          // esi
	uint32_t* v13;      // eax
	char* v14;          // esi
	uint32_t* v15;      // eax
	char* v16;          // esi
	uint32_t* v17;      // eax
	char* v18;          // esi
	uint32_t* v19;      // eax
	char* v20;          // edi
	const wchar2_t* v21; // eax
	int v22;            // esi

	if (a2 > 0x4007) {
		if (a2 == 16393) {
			sub_4BE2C0(a4);
			nox_xxx_gameSetAudioFadeoutMb_501AC0(a4);
		} else if (a2 == 16415) {
			v20 = sub_416640();
			v21 = (const wchar2_t*)nox_window_call_field_94((int)a3, 16413, 0, 0);
			if (v21) {
				if (*v21) {
					v22 = nox_wcstol(v21, 0, 10);
					if (v22 < 0) {
						v22 = 0;
					}
					if (nox_xxx_wndGetID_46B0A0(a3) == 2110) {
						*(uint32_t*)(v20 + 70) = v22;
						return 1;
					}
				}
			}
		}
		return 1;
	}
	if (a2 != 16391) {
		if (a2 != 23 && a2 == 16387) {
			v4 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1316972, a4);
			v5 = sub_416640();
			if (!v4) {
				return 0;
			}
			if ((unsigned short)a3 == 1) {
				return 0;
			}
			v7 = (const wchar2_t*)nox_window_call_field_94((int)v4, 16413, 0, 0);
			if (v7 && *v7) {
				v8 = nox_wcstol(v7, 0, 10);
				if (v8 < 0) {
					v8 = 0;
				}
				if (a4 == 2110) {
					*(uint32_t*)(v5 + 70) = v8;
					return 1;
				}
			}
		}
		return 1;
	}
	v9 = nox_xxx_wndGetID_46B0A0(a3);
	nox_xxx_clientPlaySoundSpecial_452D80(766, 100);
	switch (v9) {
	case 2102:
		v10 = sub_416640();
		*(uint32_t*)(v10 + 58) ^= 1u;
		result = 1;
		break;
	case 2103:
		v11 = sub_416640();
		*(uint32_t*)(v11 + 62) ^= 1u;
		result = 1;
		break;
	case 2106:
		v12 = sub_416640();
		v13 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1316972, 2110);
		nox_xxx_wnd_46ABB0((int)v13, 0);
		*(uint32_t*)(v12 + 66) = 0;
		sub_40A6A0(1);
		result = 1;
		break;
	case 2107:
		v14 = sub_416640();
		v15 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1316972, 2110);
		nox_xxx_wnd_46ABB0((int)v15, 0);
		*(uint32_t*)(v14 + 66) = 1;
		sub_40A6A0(1);
		result = 1;
		break;
	case 2108:
		v16 = sub_416640();
		v17 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1316972, 2110);
		nox_xxx_wnd_46ABB0((int)v17, 0);
		*(uint32_t*)(v16 + 66) = 2;
		sub_40A6A0(1);
		result = 1;
		break;
	case 2109:
		v18 = sub_416640();
		v19 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1316972, 2110);
		nox_xxx_wnd_46ABB0((int)v19, 1);
		*(uint32_t*)(v18 + 66) = 3;
		sub_40A6A0(1);
		result = 1;
		break;
	case 2130:
		sub_4BE610();
		result = 1;
		break;
	default:
		return 1;
	}
	return result;
}

//----- (004BE610) --------------------------------------------------------
int sub_4BE610() {
	int result; // eax

	result = dword_5d4594_1316972;
	if (dword_5d4594_1316972) {
		nox_xxx_wnd_46C6E0(*(int*)&dword_5d4594_1316972);
		nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1316972);
		dword_5d4594_1316972 = 0;
		result = nox_xxx_windowFocus_46B500(0);
	}
	return result;
}

//----- (004BF010) --------------------------------------------------------
int nox_xxx_clientReportSecondaryWeapon_4BF010(int a1) {
	char v3[3]; // [esp+0h] [ebp-4h]
	v3[0] = -32;
	*(uint16_t*)&v3[1] = nox_xxx_netGetUnitCodeCli_578B00(a1);
	return nox_xxx_netClientSend2_4E53C0(31, v3, 3, 0, 1);
}

//----- (004BF7E0) --------------------------------------------------------
short sub_4BF7E0(uint32_t* a1) {
	int v1;             // esi
	int v2;             // ebp
	int v3;             // eax
	int v4;             // edi
	int v5;             // ebx
	int v6;             // eax
	int v7;             // eax
	int v8;             // ebx
	int v9;             // eax
	int v10;            // ebx
	int v11;            // eax
	int v13;            // [esp+10h] [ebp-4h]
	unsigned char* v14; // [esp+18h] [ebp+4h]
	unsigned char* v15; // [esp+18h] [ebp+4h]

	v1 = a1[1] + 15;
	v2 = *a1 + 11;
	v13 = a1[1] + 15;
	nox_client_drawSetColor_434460(nox_color_black_2650656);
	nox_client_drawRectFilledOpaque_49CE30(v2, v1, 200, 200);
	LOWORD(v3) = *getMemU16Ptr(0x852978, 8);
	if (*getMemU32Ptr(0x852978, 8)) {
		v4 = dword_8531A0_2576;
		if (dword_8531A0_2576) {
			nox_draw_setMaterial_4341D0(1, *(uint32_t*)(dword_8531A0_2576 + 2296));
			nox_draw_setMaterial_4341D0(2, *(uint32_t*)(v4 + 2304));
			nox_draw_setMaterial_4341D0(3, *(uint32_t*)(v4 + 2312));
			nox_draw_setMaterial_4341D0(4, *(uint32_t*)(v4 + 2308));
			nox_draw_setMaterial_4341D0(5, *(uint32_t*)(v4 + 2300));
			nox_draw_setMaterial_4341D0(6, *(uint32_t*)(v4 + 2292));
			if (*(uint32_t*)(v4 + 2292) == *(uint32_t*)(v4 + 2296)) {
				nox_client_drawImageAt_47D2C0(*getMemU32Ptr(0x973A20, 24 + 4 * *(unsigned char*)(v4 + 2252)), v2, v1);
			} else {
				nox_client_drawImageAt_47D2C0(*getMemU32Ptr(0x973A20, 16 + 4 * *(unsigned char*)(v4 + 2252)), v2, v1);
			}
			v5 = 0;
			v14 = getMemAt(0x973A20, 32 + 104 * *(unsigned char*)(v4 + 2252));
			do {
				if (*(uint32_t*)v4 & (1 << v5) && !((1 << v5) & 0x3000000)) {
					v6 = sub_415CD0((char*)(1 << v5));
					sub_4BF9F0(1 << v5, v6, v2, v13, (int)v14, v5, 0);
				}
				++v5;
			} while (v5 < 26);
			if (*(uint8_t*)v4 & 2) {
				v7 = sub_415CD0((char*)2);
				sub_4BF9F0(2, v7, v2, v13, (int)v14, 0, 1);
			}
			v8 = 0;
			do {
				if (*(uint32_t*)v4 & (1 << v8) && (1 << v8) & 0x3000000) {
					v9 = sub_415CD0((char*)(1 << v8));
					sub_4BF9F0(1 << v8, v9, v2, v13, (int)v14, v8, 0);
				}
				++v8;
			} while (v8 < 26);
			v10 = 0;
			v15 = getMemAt(0x973A20, 256 + 108 * *(unsigned char*)(v4 + 2252));
			do {
				v3 = *(uint32_t*)(v4 + 4);
				if (v3 & (1 << v10)) {
					v11 = sub_415840((char*)(1 << v10));
					LOWORD(v3) = sub_4BF9F0(1 << v10, v11, v2, v13, (int)v15, v10, 0);
				}
				++v10;
			} while (v10 < 27);
		}
	}
	return v3;
}

//----- (004BF9F0) --------------------------------------------------------
short sub_4BF9F0(int a1, int a2, int a3, int a4, int a5, int a6, int a7) {
	int v7;        // eax
	int v8;        // esi
	uint32_t* v9;  // eax
	int v10;       // edx
	int v11;       // ecx
	uint32_t* v12; // ebx
	int v13;       // ebp
	int v14;       // edi
	uint8_t* v15;  // esi
	int* v16;      // edi
	uint8_t** v17; // esi
	int v18;       // ebx
	uint8_t* v19;  // eax

	v7 = sub_461600(a2);
	v8 = v7;
	if (v7) {
		if (*(uint32_t*)(v7 + 112) & 0x2000000) {
			v9 = nox_xxx_equipClothFindDefByTT_413270(*(uint32_t*)(v7 + 108));
		} else {
			v9 = nox_xxx_getProjectileClassById_413250(*(uint32_t*)(v7 + 108));
		}
		v12 = v9;
		if (v9) {
			v13 = v8 + 432;
			v14 = 1;
			v15 = v9 + 4;
			do {
				LOBYTE(v9) = v15[1];
				LOBYTE(v11) = *v15;
				LOBYTE(v10) = *(v15 - 1);
				nox_draw_setMaterial_4340A0(v14++, v10, v11, (int)v9);
				v15 += 3;
			} while (v14 < 7);
			v16 = v12 + 9;
			v17 = (uint8_t**)v13;
			v18 = 4;
			do {
				v19 = *v17;
				if (*v17) {
					LOBYTE(v11) = v19[26];
					LOBYTE(v10) = v19[25];
					LOBYTE(v19) = v19[24];
					nox_draw_setMaterial_4340A0(*v16, (int)v19, v10, v11);
				}
				++v17;
				++v16;
				--v18;
			} while (v18);
		}
		if (a7) {
			nox_client_drawImageAt_47D2C0(*getMemIntPtr(0x5D4594, 1319052), a3, a4);
		} else {
			nox_client_drawImageAt_47D2C0(*(uint32_t*)(a5 + 4 * a6), a3, a4);
		}
	}
	return v7;
}
// 4BFA4C: variable 'v10' is possibly undefined
// 4BFA4C: variable 'v11' is possibly undefined

//----- (004BFAD0) --------------------------------------------------------
int sub_4BFAD0() {
	int v0;         // esi
	int v1;         // ebx
	int i;          // edi
	char* v3;       // eax
	const char* v4; // ecx
	int v5;         // ebp
	int v6;         // ebp

	v0 = 0;
	v1 = 0;
	for (i = 0; i < 8; i += 4) {
		v3 = nox_xxx_gLoadImg_42F970(*(const char**)getMemAt(0x587000, 180960 + i));
		v4 = *(const char**)getMemAt(0x587000, 180968 + i);
		*getMemU32Ptr(0x973A20, 16 + i) = v3;
		*getMemU32Ptr(0x973A20, 24 + i) = nox_xxx_gLoadImg_42F970(v4);
		v5 = 26;
		do {
			*getMemU32Ptr(0x973A20, 32 + v1) = nox_xxx_gLoadImg_42F970(*(const char**)getMemAt(0x587000, 180976 + v1));
			v1 += 4;
			--v5;
		} while (v5);
		v6 = 27;
		do {
			*getMemU32Ptr(0x973A20, 256 + v0) = nox_xxx_gLoadImg_42F970(*(const char**)getMemAt(0x587000, 181184 + v0));
			v0 += 4;
			--v6;
		} while (v6);
	}
	*getMemU32Ptr(0x5D4594, 1319052) = nox_xxx_gLoadImg_42F970("MaleMedievalCloakTop");
	return 1;
}

//----- (004BFB70) --------------------------------------------------------
void sub_4BFB70(int a1) {
	if (dword_5d4594_1319056) {
		dword_5d4594_1319056 = a1;
	} else {
		if (a1 == 1) {
			nox_xxx_clientPlaySoundSpecial_452D80(1022, 100);
		}
		dword_5d4594_1319056 = a1;
	}
}

//----- (004BFBB0) --------------------------------------------------------
void sub_4BFBB0(uint32_t* a1) {
	if (dword_5d4594_1319056) {
		if (dword_5d4594_1319056 == 1) {
			if (!a1) {
				sub_4BFC70();
				sub_4BFB70(0);
			}
		}
	} else if (a1 == (uint32_t*)1) {
		sub_4BFBF0();
		sub_4BFB70(1);
	}
}

//----- (004BFBF0) --------------------------------------------------------
int sub_4BFBF0() {
	int result; // eax
	int v1;     // [esp+0h] [ebp-8h]
	int v2;     // [esp+4h] [ebp-4h]

	result = dword_5d4594_1319060;
	if (dword_5d4594_1319060) {
		nox_window_set_hidden(*(int*)&dword_5d4594_1319060, 0);
		nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1319060, 1);
		nox_window_get_size(*(int*)&dword_5d4594_1319060, &v2, &v1);
		nox_window_setPos_46A9B0(*(uint32_t**)&dword_5d4594_1319060, nox_win_width / 2 - v2 / 2,
								 nox_win_height / 2 - v1 / 2);
		result = nox_xxx_windowFocus_46B500(0);
	}
	return result;
}

//----- (004BFC70) --------------------------------------------------------
int sub_4BFC70() {
	int result; // eax

	result = dword_5d4594_1319060;
	if (dword_5d4594_1319060) {
		nox_window_set_hidden(*(int*)&dword_5d4594_1319060, 1);
		result = nox_xxx_windowFocus_46B500(0);
	}
	return result;
}

//----- (004BFC90) --------------------------------------------------------
int sub_4BFC90() {
	int result; // eax

	result = nox_new_window_from_file("SKey.wnd", sub_4BFCD0);
	dword_5d4594_1319060 = result;
	if (result) {
		sub_4BFB70(0);
		sub_4BFC70();
		result = 1;
	}
	return result;
}

//----- (004BFCD0) --------------------------------------------------------
int sub_4BFCD0(int a1, int a2, int* a3, int a4) {
	int v3; // esi

	if (a2 == 16391) {
		v3 = nox_xxx_wndGetID_46B0A0(a3);
		nox_xxx_clientPlaySoundSpecial_452D80(766, 100);
		if (v3 == 10803) {
			sub_4BFC70();
		}
	}
	return 0;
}

//----- (004BFD10) --------------------------------------------------------
void sub_4BFD10() {
	nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1319060);
	dword_5d4594_1319060 = 0;
	sub_4BFB70(0);
}

//----- (004BFD30) --------------------------------------------------------
int sub_4BFD30() { return dword_5d4594_1319056; }

//----- (004C3390) --------------------------------------------------------
int sub_4C3390() {
	*getMemU32Ptr(0x5D4594, 1321220) = nox_xxx_gLoadImg_42F970("VoteInProgress");
	dword_5d4594_1321216 = nox_window_new(0, 136, nox_win_width - 50, nox_win_height / 2 - 100, 50, 50, 0);
	nox_xxx_wndSetIcon_46AE60(*(int*)&dword_5d4594_1321216, *getMemIntPtr(0x5D4594, 1321220));
	nox_window_set_all_funcs(*(uint32_t**)&dword_5d4594_1321216, 0, sub_4C3410, 0);
	nox_window_set_hidden(*(int*)&dword_5d4594_1321216, 1);
	return 1;
}

//----- (004C3410) --------------------------------------------------------
int sub_4C3410(int* a1) {
	int* v1; // esi
	int v2;  // edx
	int v4;  // [esp+4h] [ebp-4h]

	v1 = a1;
	nox_client_wndGetPosition_46AA60(a1, &a1, &v4);
	v2 = v1[25];
	a1 = (int*)((char*)a1 + v1[24]);
	nox_client_drawImageAt_47D2C0(v1[15], (int)a1, v2 + v4);
	return 1;
}

//----- (004C3460) --------------------------------------------------------
int sub_4C3460(int a1) { return nox_window_set_hidden(*(int*)&dword_5d4594_1321216, a1); }

//----- (004C34A0) --------------------------------------------------------
int sub_4C34A0() {
	int result; // eax

	result = nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1321216);
	dword_5d4594_1321216 = 0;
	return result;
}

//----- (004CA540) --------------------------------------------------------
int nox_xxx_sprite_4CA540(uint32_t* a1, int a2) {
	uint32_t* v2; // esi
	double v3;    // st7
	double v4;    // st6
	double v5;    // st5
	double v6;    // st4
	int v7;       // eax
	int v8;       // edi
	int v9;       // eax
	float v11;    // [esp+0h] [ebp-14h]
	float v12;    // [esp+0h] [ebp-14h]
	float v13;    // [esp+10h] [ebp-4h]
	float v14;    // [esp+1Ch] [ebp+8h]

	v2 = (uint32_t*)a2;
	v3 = 0.0;
	v4 = *(float*)(a2 + 468);
	v5 = *(float*)(a2 + 472);
	v6 = 0.0;
	v7 = gameFrame() - *(uint32_t*)(a2 + 316) + 1;
	do {
		--v7;
		v13 = -(v5 * *(float*)(a2 + 476));
		v4 = v4 - v4 * *(float*)(a2 + 476);
		v5 = v5 + v13;
		v3 = v3 + v4;
		v6 = v6 + v5;
	} while (v7);
	v14 = v6;
	v11 = (double)(int)v2[81] + v3;
	v8 = nox_float2int(v11);
	v12 = (double)(int)v2[82] + v14;
	v9 = nox_float2int(v12);
	if (v8 > 0 && v9 > 0 && v8 < 5888 && v9 < 5888) {
		nox_xxx_updateSpritePosition_49AA90(v2, v8, v9);
		if (sub_4992B0(*a1 + v2[3] - a1[4], v2[4] + a1[1] - a1[5])) {
			return 1;
		}
	}
	nox_xxx_spriteDeleteStatic_45A4E0_drawable((int)v2);
	return 0;
}

//----- (004CA650) --------------------------------------------------------
int sub_4CA650(int a1, int a2) {
	int v2;             // esi
	int v3;             // eax
	int v4;             // ebx
	int v5;             // edi
	int v6;             // eax
	int v7;             // ecx
	int v8;             // eax
	int v9;             // ebp
	int v10;            // eax
	int v11;            // edi
	unsigned short v12; // dx
	int v13;            // ebp
	int v14;            // eax
	int result;         // eax
	int v16;            // [esp+10h] [ebp-10h]
	int v17;            // [esp+28h] [ebp+8h]

	v2 = a2;
	v3 = *(uint32_t*)(a2 + 16);
	v4 = *(unsigned short*)(a2 + 434) - v3;
	v5 = *(unsigned short*)(a2 + 432) - *(uint32_t*)(a2 + 12);
	v6 = sub_48C6B0(v5, v4);
	v7 = v6;
	v8 = *(unsigned char*)(a2 + 443);
	v17 = v8;
	++v7;
	v9 = *(uint32_t*)(v2 + 12);
	v10 = v5 * v8 / v7;
	v11 = *(uint32_t*)(v2 + 16);
	v16 = v9 + v10;
	v12 = *(uint16_t*)(v2 + 432);
	v13 = v9 - v12;
	v14 = v11 + v4 * v17 / v7;
	if (v7 <= 10 ||
		v13 * (v16 - v12) + (v11 - *(unsigned short*)(v2 + 434)) * (v14 - *(unsigned short*)(v2 + 434)) < 0) {
		nox_xxx_spriteDeleteStatic_45A4E0_drawable(v2);
		result = 0;
	} else {
		nox_xxx_updateSpritePosition_49AA90((uint32_t*)v2, v16, v14);
		result = 1;
	}
	return result;
}
// 4CA67E: variable 'v6' is possibly undefined


nox_window* dword_5d4594_1321236 = 0;
nox_window* dword_5d4594_1321240 = 0;
nox_window* dword_5d4594_1321244 = 0;
nox_window* dword_5d4594_1321248 = 0;
