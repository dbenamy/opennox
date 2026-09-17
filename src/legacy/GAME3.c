#include <math.h>
#include <sys/stat.h>

#include "compat.h"

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
#include "GAME4.h"
#include "GAME4_1.h"
#include "GAME5_2.h"
#include "client__draw__animdraw.h"
#include "client__drawable__drawable.h"
#include "client__video__draw_common.h"
#include "common__system__team.h"

#include "client__gui__guicon.h"
#include "client__gui__guiquit.h"
#include "client__gui__window.h"
#include "client__shell__mainmenu.h"
#include "client__shell__noxworld.h"
#include "client__shell__optsback.h"
#include "client__shell__selchar.h"
#include "client__shell__selclass.h"
#include "client__shell__selcolor.h"

#include "client__draw__drawrays.h"
#include "client__draw__fx.h"
#include "client__draw__lvupdraw.h"

#include "common/fs/nox_fs.h"
#include "common__binfile.h"
#include "common__crypt.h"
#include "input.h"
#include "input_common.h"
#include "operators.h"

extern uint32_t dword_8531A0_2572;
extern uint32_t dword_5d4594_1307792;
extern uint32_t dword_5d4594_1313788;
extern uint32_t dword_5d4594_1308136;
extern uint32_t nox_xxx_normalWndBits_587000_172880;
extern uint32_t dword_5d4594_1308104;
extern uint32_t dword_5d4594_1313532;
extern uint32_t dword_5d4594_1313564;
extern uint32_t dword_5d4594_1308096;
extern uint32_t dword_5d4594_1307724;
extern uint32_t dword_5d4594_1308112;
extern uint32_t nox_server_sendMotd_108752;
extern uint32_t dword_5d4594_1308148;
extern uint32_t dword_587000_171388;
extern uint32_t dword_5d4594_1308140;
extern uint32_t dword_5d4594_1308144;
extern uint32_t dword_5d4594_1308152;
extern uint32_t dword_5d4594_1313536;
extern uint32_t dword_5d4594_1308116;
extern uint32_t dword_5d4594_1313740;
extern uint32_t dword_5d4594_1308100;
extern uint32_t dword_5d4594_1308132;
extern uint32_t dword_5d4594_1308108;
extern uint32_t dword_5d4594_1308120;
extern uint32_t dword_5d4594_1308128;
extern uint32_t dword_5d4594_1308124;
extern uint32_t dword_5d4594_1309736;
extern uint32_t dword_5d4594_1309756;
extern uint32_t dword_5d4594_1309832;
extern uint32_t dword_5d4594_1313540;
extern uint32_t dword_5d4594_1309824;
extern uint32_t dword_5d4594_1307720;
extern uint32_t dword_5d4594_1307736;
extern uint32_t nox_server_connectionType_3596;
extern void* dword_587000_122852;
extern uint32_t dword_5d4594_1309828;
extern uint32_t dword_5d4594_1309836;
extern uint32_t dword_5d4594_1309728;
extern uint32_t dword_5d4594_1309732;
extern uint64_t qword_581450_9552;
extern void* dword_587000_93164;
extern uint32_t dword_5d4594_1307716;
extern void* dword_587000_127004;
extern uint32_t dword_5d4594_1308088;
extern uint32_t dword_5d4594_1309748;
extern uint32_t dword_5d4594_1309720;
extern uint32_t dword_5d4594_1308084;
extern uint32_t dword_5d4594_1309820;
extern uint32_t dword_5d4594_2650652;
extern uint32_t dword_5d4594_1307784;
extern int nox_win_width;
extern int nox_win_height;

extern uint32_t nox_color_white_2523948;
extern uint32_t nox_color_blue_2650684;
extern uint32_t nox_color_violet_2598268;
extern uint32_t nox_color_black_2650656;

nox_gui_animation* nox_wnd_xxx_1307732 = 0;
nox_gui_animation* nox_wnd_xxx_1308092 = 0;
nox_gui_animation* nox_wnd_xxx_1309740 = 0;

void* dword_5d4594_1308156 = 0;
void* dword_5d4594_1308160 = 0;
void* dword_5d4594_1308164 = 0;

//----- (004A2560) --------------------------------------------------------
int sub_4A2560(uint32_t* a1, int a2) {
	double v2; // st7
	double v3; // st6

	v2 = (double)(*(short*)(a2 + 44) - *a1);
	v3 = (double)(*(short*)(a2 + 46) - a1[1]);
	return sqrt(v3 * v3 + v2 * v2) <= *getMemDoublePtr(0x581450, 9720);
}

//----- (004A25C0) --------------------------------------------------------
int sub_4A25C0(uint32_t* a1, int* a2) {
	int v2;  // edi
	int* v3; // esi

	v2 = 0;
	v3 = nox_common_list_getFirstSafe_425890(a2);
	if (!v3) {
		return 0;
	}
	do {
		if (sub_4A2560(a1, (int)v3)) {
			++v2;
		}
		v3 = nox_common_list_getNextSafe_4258A0(v3);
	} while (v3);
	return v2;
}

//----- (004A2610) --------------------------------------------------------
int sub_4A2610(int a1, uint32_t* a2, int* a3) {
	int* i;             // esi
	int v4;             // eax
	uint32_t* v5;       // esi
	uint32_t* v6;       // ebx
	uint32_t* v7;       // edi
	uint32_t* v8;       // ebp
	char* v9;           // eax
	int v10;            // ebx
	unsigned char* v11; // esi
	uint32_t* v13;      // [esp+Ch] [ebp-150h]
	char* v14;          // [esp+10h] [ebp-14Ch]
	int v15[2];         // [esp+14h] [ebp-148h]
	char v16[64];       // [esp+1Ch] [ebp-140h]
	wchar2_t v17[128];   // [esp+5Ch] [ebp-100h]

	dword_5d4594_1307720 = 0;
	for (i = nox_common_list_getFirstSafe_425890(a3); i; i = nox_common_list_getNextSafe_4258A0(i)) {
		if (sub_4A2560(a2, (int)i)) {
			v4 = dword_5d4594_1307720;
			*getMemU32Ptr(0x5D4594, 1307316 + 4 * dword_5d4594_1307720) = i;
			dword_5d4594_1307720 = v4 + 1;
		}
	}
	if (dword_5d4594_1307720 > 0) {
		dword_5d4594_1307716 = nox_new_window_from_file("proxlist.wnd", *(uint32_t*)(a1 + 376));
		sub_4A2830(*a2 + 216, a2[1] + 27, v15);
		nox_window_setPos_46A9B0(*(uint32_t**)&dword_5d4594_1307716, v15[0], v15[1]);
		nox_xxx_wnd_46B280(*(int*)&dword_5d4594_1307716, a1);
		v5 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1307716, 10064);
		v6 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1307716, 10062);
		v13 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1307716, 10063);
		v7 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1307716, 10061);
		v8 = (uint32_t*)v7[8];
		v14 = nox_xxx_gLoadImg_42F970("UISlider");
		v9 = nox_xxx_gLoadImg_42F970("UISliderLit");
		sub_4B5700((int)v5, 0, 0, (int)v14, (int)v9, (int)v9);
		nox_xxx_wnd_46B280((int)v5, (int)v7);
		nox_xxx_wnd_46B280((int)v6, (int)v7);
		nox_xxx_wnd_46B280((int)v13, (int)v7);
		v8[9] = v5;
		v8[7] = v6;
		v8[8] = v13;
		*(uint32_t*)(v5[100] + 8) = 16;
		*(uint32_t*)(v5[100] + 12) = 10;
		v10 = 0;
		if (dword_5d4594_1307720 > 0) {
			v11 = getMemAt(0x5D4594, 1307316);
			do {
				if (*(uint8_t*)(*(uint32_t*)v11 + 120)) {
					strncpy(v16, (const char*)(*(uint32_t*)v11 + 120), 0xFu);
					v16[15] = 0;
				} else {
					nox_sprintAddrPort_43BC80(*(uint32_t*)v11 + 12, *(uint16_t*)(*(uint32_t*)v11 + 109), v16);
				}
				nox_swprintf(v17, L"%S   %dms", v16, *(uint32_t*)(*(uint32_t*)v11 + 96));
				nox_window_call_field_94((int)v7, 16397, (int)v17, -1);
				++v10;
				v11 += 4;
			} while (v10 < *(int*)&dword_5d4594_1307720);
		}
	}
	return dword_5d4594_1307716;
}

//----- (004A2830) --------------------------------------------------------
uint32_t* sub_4A2830(int a1, int a2, uint32_t* a3) {
	uint32_t* result; // eax

	result = a3;
	*a3 = a1 - 100;
	a3[1] = a2 - 20;
	if (a1 - 100 + 200 > 600) {
		*a3 = 400;
	}
	if (a2 - 20 + 200 > 451) {
		a3[1] = 251;
	}
	if (a3[1] < 27) {
		a3[1] = 27;
	}
	if ((int)*a3 < 216) {
		*a3 = 216;
	}
	return result;
}

//----- (004A2890) --------------------------------------------------------
int sub_4A2890() {
	int result; // eax

	result = dword_5d4594_1307716;
	if (dword_5d4594_1307716) {
		result = nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1307716);
		dword_5d4594_1307716 = 0;
	}
	return result;
}

//----- (004A28B0) --------------------------------------------------------
int sub_4A28B0() { return dword_5d4594_1307716 != 0; }

//----- (004A28C0) --------------------------------------------------------
int sub_4A28C0(int a1) {
	int result; // eax

	if (a1 < *(int*)&dword_5d4594_1307720) {
		result = *getMemU32Ptr(0x5D4594, 1307316 + 4 * a1);
	} else {
		result = 0;
	}
	return result;
}

//----- (004A4840) --------------------------------------------------------
int nox_game_showSelClass_4A4840() {
	int result;   // eax
	uint32_t* v1; // eax
	uint32_t* v2; // eax
	uint32_t* v3; // eax

	sub_5007E0("*:*");
	sub_4A1BE0(1);
	dword_5d4594_1307724 = nox_xxx_getHostInfoPtr_431770();
	nox_game_addStateCode_43BDD0(600);
	result = nox_new_window_from_file("SelClass.wnd", sub_4A4A20);
	dword_5d4594_1307736 = result;
	if (result) {
		nox_xxx_wndSetWindowProc_46B300(result, sub_4A18E0);
		result = nox_gui_makeAnimation_43C5B0(*(uint32_t**)&dword_5d4594_1307736, 0, 0, 0, -460, 0, 20, 0, -40);
		nox_wnd_xxx_1307732 = result;
		if (result) {
			nox_wnd_xxx_1307732->field_0 = 600;
			nox_wnd_xxx_1307732->field_12 = sub_4A4970;
			nox_wnd_xxx_1307732->fnc_done_out = sub_4A49A0;
			v1 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1307736, 601);
			nox_xxx_wndSetDrawFn_46B340((int)v1, sub_4A49D0);
			v2 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1307736, 603);
			nox_xxx_wndSetDrawFn_46B340((int)v2, sub_4A49D0);
			v3 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1307736, 602);
			nox_xxx_wndSetDrawFn_46B340((int)v3, sub_4A49D0);
			*getMemU32Ptr(0x5D4594, 1307728) = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1307736, 610);
			nox_xxx_wndRetNULL_46A8A0();
			*getMemU32Ptr(0x5D4594, 1307740) = 0;
			sub_4A19F0("OptsBack.wnd:Back");
			sub_4602F0();
			result = 1;
		}
	}
	return result;
}
// 4A18E0: using guessed type int  sub_4A18E0(int, int, int, int);

//----- (004A4970) --------------------------------------------------------
int sub_4A4970() {
	nox_wnd_xxx_1307732->state = NOX_GUI_ANIM_OUT;
	sub_43BE40(2);
	nox_xxx_clientPlaySoundSpecial_452D80(923, 100);
	return 1;
}

//----- (004A49A0) --------------------------------------------------------
int sub_4A49A0() {
	int (*v0)(void); // esi

	v0 = nox_wnd_xxx_1307732->field_13;
	nox_gui_freeAnimation_43C570(nox_wnd_xxx_1307732);
	nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1307736);
	v0();
	return 1;
}

//----- (004A49D0) --------------------------------------------------------
int sub_4A49D0(int yTop, int a2) {
	uint32_t* v1; // esi
	int xLeft;    // [esp+4h] [ebp-4h]

	v1 = (uint32_t*)yTop;
	if (*getMemU32Ptr(0x5D4594, 1307740) != *(uint32_t*)yTop) {
		nox_client_wndGetPosition_46AA60((uint32_t*)yTop, &xLeft, &yTop);
		nox_client_drawRectFilledAlpha_49CF10(xLeft, yTop, v1[6] - v1[4], v1[7] - v1[5]);
	}
	return 1;
}

//----- (004A5E90) --------------------------------------------------------
void sub_4A5E90_A();
int sub_4A5E90() {
	const char** i; // eax
	uint32_t* v1;   // eax
	uint32_t* v2;   // eax
	uint32_t* v3;   // eax
	uint32_t* v4;   // eax
	uint32_t* v5;   // eax
	uint32_t* v6;   // eax
	uint32_t* v7;   // eax
	uint32_t* v8;   // eax
	uint32_t* v9;   // eax
	uint32_t* v10;  // eax
	int result;     // eax

	sub_4A5E90_A();
	v1 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1308084, 720);
	dword_5d4594_1308096 = v1;
	v1[8] = 131074;
	v2 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1308084, 721);
	dword_5d4594_1308100 = v2;
	v2[8] = 589825;
	v3 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1308084, 722);
	dword_5d4594_1308104 = v3;
	v3[8] = 589825;
	v4 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1308084, 723);
	dword_5d4594_1308108 = v4;
	v4[8] = 589825;
	v5 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1308084, 724);
	dword_5d4594_1308112 = v5;
	v5[8] = 589825;
	v6 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1308084, 725);
	dword_5d4594_1308116 = v6;
	v6[8] = *getMemU16Ptr(0x587000, 171372) << 16;
	v7 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1308084, 726);
	dword_5d4594_1308120 = v7;
	v7[8] = *getMemU16Ptr(0x587000, 171374) << 16;
	v8 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1308084, 727);
	dword_5d4594_1308124 = v8;
	v8[8] = *getMemU16Ptr(0x587000, 171376) << 16;
	v9 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1308084, 728);
	dword_5d4594_1308128 = v9;
	v9[8] = *getMemU16Ptr(0x587000, 171378) << 16;
	v10 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1308084, 729);
	dword_5d4594_1308132 = v10;
	v10[8] = *getMemU16Ptr(0x587000, 171380) << 16;
	dword_5d4594_1308136 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1308084, 711);
	dword_5d4594_1308140 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1308084, 712);
	dword_5d4594_1308144 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1308084, 713);
	dword_5d4594_1308148 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1308084, 714);
	dword_5d4594_1308152 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1308084, 751);
	result = dword_587000_171388;
	if (!dword_587000_171388) {
		nox_window_call_field_94(*(int*)&dword_5d4594_1308152, 16414, *(int*)&dword_5d4594_1307784, 0);
		sub_4A61E0(*(uint32_t**)&dword_5d4594_1308096, 2, (uint8_t*)(dword_5d4594_1307784 + 71));
		sub_4A61E0(*(uint32_t**)&dword_5d4594_1308100, 1, (uint8_t*)(dword_5d4594_1307784 + 68));
		sub_4A61E0(*(uint32_t**)&dword_5d4594_1308104, 1, (uint8_t*)(dword_5d4594_1307784 + 74));
		sub_4A61E0(*(uint32_t**)&dword_5d4594_1308108, 1, (uint8_t*)(dword_5d4594_1307784 + 77));
		sub_4A61E0(*(uint32_t**)&dword_5d4594_1308112, 1, (uint8_t*)(dword_5d4594_1307784 + 80));
		*(uint32_t*)(dword_5d4594_1308116 + 32) = *(unsigned char*)(dword_5d4594_1307784 + 83) << 16;
		*(uint32_t*)(dword_5d4594_1308120 + 32) = *(unsigned char*)(dword_5d4594_1307784 + 84) << 16;
		*(uint32_t*)(dword_5d4594_1308124 + 32) = *(unsigned char*)(dword_5d4594_1307784 + 85) << 16;
		*(uint32_t*)(dword_5d4594_1308128 + 32) = *(unsigned char*)(dword_5d4594_1307784 + 86) << 16;
		result = *(unsigned char*)(dword_5d4594_1307784 + 87) << 16;
		*(uint32_t*)(dword_5d4594_1308132 + 32) = result;
	}
	return result;
}

//----- (004A61E0) --------------------------------------------------------
unsigned char* sub_4A61E0(uint32_t* a1, int a2, unsigned char* a3) {
	unsigned int v3;       // esi
	unsigned char* result; // eax
	uint32_t* v5;          // eax
	uint32_t* v6;          // eax

	v3 = 0;
	result = getMemAt(0x5D4594, 1307797 + 96 * a2);
	while (*(result - 1) != *a3 || *result != a3[1] || result[1] != a3[2]) {
		++v3;
		result += 3;
		if (v3 >= 0x20) {
			goto LABEL_9;
		}
	}
	if (a2 == 1) {
		v5 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1308084, *a1 - 10);
		nox_xxx_wnd_46ABB0((int)v5, 1);
		result = (unsigned char*)nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1308084, *a1 + 10);
		*((uint32_t*)result + 9) |= 6u;
	}
LABEL_9:
	if (v3 == 32 && a2 == 1) {
		v6 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1308084, *a1 - 10);
		result = (unsigned char*)nox_xxx_wnd_46ABB0((int)v6, 0);
		LOWORD(v3) = 9;
	}
	a1[8] = (unsigned short)a2 | ((unsigned short)v3 << 16);
	return result;
}

//----- (004A6890) --------------------------------------------------------
int sub_4A6890() {
	nox_wnd_xxx_1308092->state = NOX_GUI_ANIM_OUT;
	sub_43BE40(2);
	nox_xxx_clientPlaySoundSpecial_452D80(923, 100);
	sub_4A68C0();
	return 1;
}

//----- (004A6B50) --------------------------------------------------------
int sub_4A6B50(wchar2_t* a1) {
	wchar2_t* v1;    // esi
	int v2;         // ebp
	int v3;         // ebx
	signed int v4;  // eax
	short* v5;      // edi
	wchar2_t v6;     // ax
	signed int v7;  // esi
	wchar2_t* i;     // edi
	int v10;        // [esp+10h] [ebp-3Ch]
	signed int v11; // [esp+14h] [ebp-38h]
	short v12[26];  // [esp+18h] [ebp-34h]

	v1 = a1;
	v2 = 0;
	v3 = 1;
	v10 = 0;
	v4 = nox_wcslen(a1);
	if (v4 >= 1) {
		v5 = v12;
		v11 = v4;
		do {
			if (iswspace(*v1)) {
				if (!v3) {
					*v5 = *v1;
					++v5;
					++v10;
				}
			} else {
				if (v3 == 1) {
					v6 = *v1;
					if (*v1 == 42 || v6 == 63 || v6 == 60 || v6 == 62 || v6 == 92 || v6 == 47 || v6 == 58 || v6 == 34 ||
						v6 == 124) {
						*v1 = 45;
					}
				}
				*v5 = *v1;
				++v5;
				++v10;
				v3 = 0;
				v2 = 1;
			}
			++v1;
			--v11;
		} while (v11);
		v12[v10] = 0;
		if (v2) {
			nox_wcscpy(a1, (const wchar2_t*)v12);
			v7 = nox_wcslen((const wchar2_t*)v12) - 1;
			if ((int)v7 >= 0) {
				for (i = &a1[v7]; iswspace(*i); --i) {
					if (--v7 < 0) {
						return v2;
					}
				}
				a1[v7 + 1] = 0;
			}
		}
	}
	return v2;
}
// 4A6B50: using guessed type wchar2_t var_34[26];

//----- (004A6C90) --------------------------------------------------------
int sub_4A6C90() {
	int (*v0)(void); // esi

	v0 = nox_wnd_xxx_1308092->field_13;
	nox_gui_freeAnimation_43C570(nox_wnd_xxx_1308092);
	nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1308084);
	nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1308088);
	if (v0) {
		v0();
	} else {
		nox_client_resetScreenParticles_431510();
		nox_gui_draw();
		if (!nox_common_gameFlags_check_40A5C0(0x2000)) {
			nox_client_guiXxxDestroy_4A24A0();
			return 1;
		}
		nox_xxx_serverHost_43B4D0();
		return 1;
	}
	return 1;
}

//----- (004A6D20) --------------------------------------------------------
int sub_4A6D20(int a1, int a2) {
	int v1;            // esi
	unsigned short v2; // di
	int v3;            // ebx
	int v4;            // ecx
	int v5;            // edx
	int v6;            // edi
	int v7;            // ebp
	int v8;            // eax
	int v10;           // [esp-10h] [ebp-28h]
	int yTop;          // [esp+10h] [ebp-8h]
	int xLeft;         // [esp+14h] [ebp-4h]
	int v13;           // [esp+1Ch] [ebp+4h]

	v1 = a1;
	v2 = *(uint16_t*)(a1 + 32);
	v3 = *(uint32_t*)(a1 + 32) >> 16;
	nox_client_wndGetPosition_46AA60((uint32_t*)a1, &xLeft, &yTop);
	v4 = *(uint32_t*)(a1 + 20);
	v5 = *(uint32_t*)(a1 + 16);
	v6 = (unsigned short)v3 + 32 * v2;
	v13 = *(uint32_t*)(a1 + 28) - v4;
	v7 = *(uint32_t*)(v1 + 24) - v5;
	LOBYTE(v5) = getMemByte(0x5D4594, 1307798 + 3 * v6);
	LOBYTE(v4) = getMemByte(0x5D4594, 1307797 + 3 * v6);
	v10 = v5;
	LOBYTE(v5) = getMemByte(0x5D4594, 1307796 + 3 * v6);
	v8 = nox_color_rgb_4344A0(v5, v4, v10);
	nox_client_drawSetColor_434460(v8);
	nox_client_drawRectFilledOpaque_49CE30(xLeft, yTop, v7, v13);
	return 1;
}

//----- (004A6DC0) --------------------------------------------------------
int sub_4A6DC0(uint32_t* a1, int a2) {
	int v1;   // eax
	int v2;   // ecx
	int v3;   // edx
	int v4;   // esi
	int v5;   // eax
	int v6;   // ecx
	int v7;   // edx
	int v8;   // eax
	int v9;   // edx
	int v10;  // eax
	int v11;  // ecx
	int v12;  // eax
	int v13;  // edx
	int v14;  // ecx
	int v15;  // eax
	int v16;  // eax
	int v17;  // edx
	int v18;  // ecx
	int v19;  // eax
	int v20;  // eax
	int v21;  // ebx
	int v22;  // eax
	int v23;  // edx
	int* v24; // ebp
	int v25;  // ecx
	int v26;  // eax
	int v27;  // edi
	int i;    // esi
	int v29;  // eax
	int v30;  // eax
	int v31;  // ecx
	int v32;  // eax
	int v33;  // edi
	int j;    // esi
	int v35;  // eax
	int v36;  // eax
	int v37;  // ecx
	int v38;  // eax
	int v39;  // eax
	int v40;  // ecx
	int v41;  // edx
	int v42;  // edi
	int k;    // esi
	int v44;  // eax
	int v45;  // eax
	int v46;  // ecx
	int v47;  // eax
	int v48;  // eax
	int v49;  // ecx
	int v50;  // edx
	int v52;  // [esp-24h] [ebp-3Ch]
	int v53;  // [esp-24h] [ebp-3Ch]
	int v54;  // [esp+10h] [ebp-8h]
	int v55;  // [esp+14h] [ebp-4h]

	nox_client_wndGetPosition_46AA60(a1, &v54, &v55);
	v1 = ((unsigned short)(*(uint32_t*)(dword_5d4594_1308096 + 32) >> 16)) +
		 32 * *(unsigned short*)(dword_5d4594_1308096 + 32);
	v2 = getMemByte(0x5D4594, 1307798 + 3 * v1);
	LOBYTE(v3) = getMemByte(0x5D4594, 1307797 + 3 * v1);
	LOBYTE(v1) = getMemByte(0x5D4594, 1307796 + 3 * v1);
	v4 = nox_color_rgb_4344A0(v1, v3, v2);
	nox_draw_setMaterial_4341D0(6, v4);
	v5 = ((unsigned short)(*(uint32_t*)(dword_5d4594_1308100 + 32) >> 16)) +
		 32 * *(unsigned short*)(dword_5d4594_1308100 + 32);
	v6 = getMemByte(0x5D4594, 1307798 + 3 * v5);
	LOBYTE(v7) = getMemByte(0x5D4594, 1307797 + 3 * v5);
	LOBYTE(v5) = getMemByte(0x5D4594, 1307796 + 3 * v5);
	v8 = nox_color_rgb_4344A0(v5, v7, v6);
	nox_draw_setMaterial_4341D0(1, v8);
	v9 = dword_5d4594_1308144;
	if (*(uint8_t*)(dword_5d4594_1308144 + 4) & 8) {
		v10 = (unsigned short)(*(uint32_t*)(dword_5d4594_1308108 + 32) >> 16);
		v11 = v10 + 32 * *(unsigned short*)(dword_5d4594_1308108 + 32);
		LOBYTE(v9) = getMemByte(0x5D4594, 1307797 + 3 * v11);
		v12 = nox_color_rgb_4344A0(
			getMemByte(0x5D4594, 1307796 + 3 * v11), v9,
			getMemByte(0x5D4594, 1307798 + 3 * (v10 + 32 * *(unsigned short*)(dword_5d4594_1308108 + 32))));
		nox_draw_setMaterial_4341D0(2, v12);
	} else {
		nox_draw_setMaterial_4341D0(2, v4);
	}
	nox_draw_setMaterial_4341D0(3, v4);
	v13 = dword_5d4594_1308148;
	if (*(uint8_t*)(dword_5d4594_1308148 + 4) & 8) {
		v14 = (unsigned short)(*(uint32_t*)(dword_5d4594_1308112 + 32) >> 16);
		v15 = v14 + 32 * *(unsigned short*)(dword_5d4594_1308112 + 32);
		LOBYTE(v13) = getMemByte(0x5D4594, 1307797 + 3 * v15);
		LOBYTE(v15) = getMemByte(0x5D4594, 1307796 + 3 * v15);
		v16 = nox_color_rgb_4344A0(
			v15, v13, getMemByte(0x5D4594, 1307798 + 3 * (v14 + 32 * *(unsigned short*)(dword_5d4594_1308112 + 32))));
		nox_draw_setMaterial_4341D0(4, v16);
	} else {
		nox_draw_setMaterial_4341D0(4, v4);
	}
	v17 = dword_5d4594_1308140;
	if (*(uint8_t*)(dword_5d4594_1308140 + 4) & 8) {
		v18 = (unsigned short)(*(uint32_t*)(dword_5d4594_1308104 + 32) >> 16);
		v19 = v18 + 32 * *(unsigned short*)(dword_5d4594_1308104 + 32);
		LOBYTE(v17) = getMemByte(0x5D4594, 1307797 + 3 * v19);
		LOBYTE(v19) = getMemByte(0x5D4594, 1307796 + 3 * v19);
		v20 = nox_color_rgb_4344A0(
			v19, v17, getMemByte(0x5D4594, 1307798 + 3 * (v18 + 32 * *(unsigned short*)(dword_5d4594_1308104 + 32))));
		nox_draw_setMaterial_4341D0(5, v20);
	} else {
		nox_draw_setMaterial_4341D0(5, v4);
	}
	if (*(uint8_t*)(dword_5d4594_1308136 + 4) & 8) {
		nox_client_drawImageAt_47D2C0(*getMemU32Ptr(0x973A20, 16 + 4 * *(unsigned char*)(dword_5d4594_1307784 + 67)),
									  v54, v55);
	} else {
		nox_client_drawImageAt_47D2C0(*getMemU32Ptr(0x973A20, 24 + 4 * *(unsigned char*)(dword_5d4594_1307784 + 67)),
									  v54, v55);
	}
	v21 = 0;
	v22 = *(unsigned char*)(dword_5d4594_1307784 + 67);
	v23 = 3 * v22;
	v24 = getMemIntPtr(0x973A20, 32 + 104 * v22);
	do {
		v25 = v21;
		v26 = 1 << v21;
		if (1 << v21 == 4) {
			v27 = 1;
			for (i = 3; i < 21; i += 3) {
				v29 = dword_5d4594_1308156;
				LOBYTE(v25) = *(uint8_t*)(i + (uint32_t)dword_5d4594_1308156 + 14);
				LOBYTE(v23) = *(uint8_t*)(i + (uint32_t)dword_5d4594_1308156 + 13);
				LOBYTE(v29) = *(uint8_t*)(i + (uint32_t)dword_5d4594_1308156 + 12);
				nox_draw_setMaterial_4340A0(v27++, v29, v23, v25);
			}
			v30 = ((unsigned short)(*(uint32_t*)(dword_5d4594_1308116 + 32) >> 16)) +
				  32 * *(unsigned short*)(dword_5d4594_1308116 + 32);
			v31 = getMemByte(0x5D4594, 1307798 + 3 * v30);
			LOBYTE(v23) = getMemByte(0x5D4594, 1307797 + 3 * v30);
			LOBYTE(v30) = getMemByte(0x5D4594, 1307796 + 3 * v30);
			v32 = nox_color_rgb_4344A0(v30, v23, v31);
			nox_draw_setMaterial_4341D0(*(uint32_t*)((uint32_t)dword_5d4594_1308156 + 40), v32);
		} else if (v26 == 1024) {
			v33 = 1;
			for (j = 3; j < 21; j += 3) {
				v35 = dword_5d4594_1308160;
				LOBYTE(v25) = *(uint8_t*)((uint32_t)dword_5d4594_1308160 + j + 14);
				LOBYTE(v23) = *(uint8_t*)((uint32_t)dword_5d4594_1308160 + j + 13);
				LOBYTE(v35) = *(uint8_t*)((uint32_t)dword_5d4594_1308160 + j + 12);
				nox_draw_setMaterial_4340A0(v33++, v35, v23, v25);
			}
			v36 = ((unsigned short)(*(uint32_t*)(dword_5d4594_1308120 + 32) >> 16)) +
				  32 * *(unsigned short*)(dword_5d4594_1308120 + 32);
			v37 = getMemByte(0x5D4594, 1307798 + 3 * v36);
			LOBYTE(v23) = getMemByte(0x5D4594, 1307797 + 3 * v36);
			LOBYTE(v36) = getMemByte(0x5D4594, 1307796 + 3 * v36);
			v38 = nox_color_rgb_4344A0(v36, v23, v37);
			nox_draw_setMaterial_4341D0(*(uint32_t*)((uint32_t)dword_5d4594_1308160 + 40), v38);
			v39 = ((unsigned short)(*(uint32_t*)(dword_5d4594_1308124 + 32) >> 16)) +
				  32 * *(unsigned short*)(dword_5d4594_1308124 + 32);
			v40 = getMemByte(0x5D4594, 1307798 + 3 * v39);
			LOBYTE(v41) = getMemByte(0x5D4594, 1307797 + 3 * v39);
			LOBYTE(v39) = getMemByte(0x5D4594, 1307796 + 3 * v39);
			v52 = nox_color_rgb_4344A0(v39, v41, v40);
			nox_draw_setMaterial_4341D0(*(uint32_t*)((uint32_t)dword_5d4594_1308160 + 44), v52);
		} else {
			if (v26 != 1) {
				goto LABEL_27;
			}
			v42 = 1;
			for (k = 3; k < 21; k += 3) {
				v44 = dword_5d4594_1308164;
				LOBYTE(v25) = *(uint8_t*)(k + (uint32_t)dword_5d4594_1308164 + 14);
				LOBYTE(v23) = *(uint8_t*)(k + (uint32_t)dword_5d4594_1308164 + 13);
				LOBYTE(v44) = *(uint8_t*)(k + (uint32_t)dword_5d4594_1308164 + 12);
				nox_draw_setMaterial_4340A0(v42++, v44, v23, v25);
			}
			v45 = ((unsigned short)(*(uint32_t*)(dword_5d4594_1308128 + 32) >> 16)) +
				  32 * *(unsigned short*)(dword_5d4594_1308128 + 32);
			v46 = getMemByte(0x5D4594, 1307798 + 3 * v45);
			LOBYTE(v23) = getMemByte(0x5D4594, 1307797 + 3 * v45);
			LOBYTE(v45) = getMemByte(0x5D4594, 1307796 + 3 * v45);
			v47 = nox_color_rgb_4344A0(v45, v23, v46);
			nox_draw_setMaterial_4341D0(*(uint32_t*)((uint32_t)dword_5d4594_1308164 + 40), v47);
			v48 = ((unsigned short)(*(uint32_t*)(dword_5d4594_1308132 + 32) >> 16)) +
				  32 * *(unsigned short*)(dword_5d4594_1308132 + 32);
			v49 = getMemByte(0x5D4594, 1307798 + 3 * v48);
			LOBYTE(v50) = getMemByte(0x5D4594, 1307797 + 3 * v48);
			LOBYTE(v48) = getMemByte(0x5D4594, 1307796 + 3 * v48);
			v53 = nox_color_rgb_4344A0(v48, v50, v49);
			nox_draw_setMaterial_4341D0(*(uint32_t*)((uint32_t)dword_5d4594_1308164 + 36), v53);
		}
		nox_client_drawImageAt_47D2C0(*v24, v54, v55);
	LABEL_27:
		++v21;
		++v24;
	} while (v21 < 26);
	return 1;
}
// 4A6E13: variable 'v3' is possibly undefined
// 4A6E5A: variable 'v7' is possibly undefined
// 4A7032: variable 'v23' is possibly undefined
// 4A7032: variable 'v25' is possibly undefined
// 4A7158: variable 'v41' is possibly undefined
// 4A7227: variable 'v50' is possibly undefined

//----- (004A7270) --------------------------------------------------------
int sub_4A7270(int a1, int a2, unsigned int a3, int a4) {
	int v3;  // eax
	int2 v5; // [esp+0h] [ebp-8h]

	if (a2 != 5) {
		return 0;
	}
	v5.field_4 = a3 >> 16;
	v5.field_0 = (unsigned short)a3;
	v3 = nox_xxx_pointInRect_4281F0(&v5, (int4*)(a1 + 16));
	if (v3) {
		return 0;
	}
	sub_4A72D0(0xDEADu);
	return 1;
}
// 4A72AA: variable 'v3' is possibly undefined

//----- (004A72D0) --------------------------------------------------------
uint32_t* sub_4A72D0(unsigned short a1) {
	uint32_t* result; // eax

	nox_xxx_wnd_46C6E0(*(int*)&dword_5d4594_1308088);
	result = (uint32_t*)nox_window_set_hidden(*(int*)&dword_5d4594_1308088, 1);
	if (a1 < 0x20u) {
		result = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1308088, *(int*)&dword_5d4594_1307792);
		if (result) {
			result[8] = *getMemU16Ptr(0x5D4594, 1307788) | (a1 << 16);
		}
	}
	return result;
}

//----- (004A7330) --------------------------------------------------------
int sub_4A7330(int a1, int a2, int* a3, unsigned int a4) {
	int result;   // eax
	int v5;       // edi
	int v6;       // eax
	uint32_t* v7; // eax

	if (a2 == 16389) {
		nox_xxx_clientPlaySoundSpecial_452D80(920, 100);
		result = 1;
	} else if (a2 == 16391) {
		v5 = a4 >> 16;
		v6 = nox_xxx_wndGetID_46B0A0(a3);
		switch (v6) {
		case 720:
			dword_5d4594_1307792 = v6;
			sub_4A7530(2u);
			sub_4A7580((unsigned short)a4, v5);
			break;
		case 721:
		case 722:
		case 723:
		case 724:
			dword_5d4594_1307792 = v6;
			sub_4A7530(1u);
			sub_4A7580((unsigned short)a4, v5);
			break;
		case 725:
		case 726:
		case 727:
		case 728:
		case 729:
			dword_5d4594_1307792 = v6;
			sub_4A7530(0);
			sub_4A7580((unsigned short)a4, v5);
			break;
		case 731:
		case 732:
		case 733:
		case 734:
			v7 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1308084, v6 - 20);
			if (v7) {
				nox_xxx_wnd_46ABB0((int)v7, ((unsigned int)~v7[1] >> 3) & 1);
			}
			break;
		case 761:
		case 762:
		case 763:
		case 764:
		case 765:
		case 766:
		case 767:
		case 768:
		case 769:
		case 770:
		case 771:
		case 772:
		case 773:
		case 774:
		case 775:
		case 776:
		case 777:
		case 778:
		case 779:
		case 780:
		case 781:
		case 782:
		case 783:
		case 784:
		case 785:
		case 786:
		case 787:
		case 788:
		case 789:
		case 790:
		case 791:
		case 792:
			sub_4A72D0(v6 - 761);
			break;
		case 799:
			if (*getMemU32Ptr(0x5D4594, 1308168) == 1) {
				nox_game_decStateInd_43BDC0();
			}
			nox_game_decStateInd_43BDC0();
			nox_game_decStateInd_43BDC0();
			dword_587000_171388 = 1;
			if (sub_4A75C0()) {
				if (*(uint8_t*)(dword_5d4594_1307784 + 66) == 0) {
					nox_xxx_gameSetMapPath_409D70("war01a.map");
				} else if (*(uint8_t*)(dword_5d4594_1307784 + 66) == 1) {
					nox_xxx_gameSetMapPath_409D70("wiz01a.map");
				} else if (*(uint8_t*)(dword_5d4594_1307784 + 66) == 2) {
					nox_xxx_gameSetMapPath_409D70("con01a.map");
				}
				sub_4A24C0(0);
				sub_4A6890();
				nox_wnd_xxx_1308092->field_13 = 0;
			}
			break;
		default:
			break;
		}
		nox_xxx_clientPlaySoundSpecial_452D80(921, 100);
		result = 1;
	} else {
		result = 0;
	}
	return result;
}

//----- (004A7530) --------------------------------------------------------
uint32_t* sub_4A7530(unsigned short a1) {
	int v1;           // esi
	uint32_t* result; // eax

	v1 = 761;
	*getMemU32Ptr(0x5D4594, 1307788) = a1;
	do {
		result = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1308088, v1);
		if (result) {
			result[8] = a1 | ((unsigned short)(v1 - 761) << 16);
		}
		++v1;
	} while (v1 <= 792);
	return result;
}

//----- (004A7580) --------------------------------------------------------
int sub_4A7580(int a1, int a2) {
	nox_xxx_wndShowModalMB_46A8C0(*(int*)&dword_5d4594_1308088);
	sub_46C690(*(int*)&dword_5d4594_1308088);
	return nox_window_setPos_46A9B0(*(uint32_t**)&dword_5d4594_1308088, a1 - *(uint32_t*)(dword_5d4594_1308088 + 8),
									a2 - *(uint32_t*)(dword_5d4594_1308088 + 12) / 2);
}

//----- (004A7A60) --------------------------------------------------------
int sub_4A7A60(int a1) {
	int result; // eax

	result = a1;
	dword_587000_171388 = a1;
	return result;
}

//----- (004A7A70) --------------------------------------------------------
int sub_4A7A70(int a1) {
	int result; // eax

	result = a1;
	*getMemU32Ptr(0x5D4594, 1308168) = a1;
	return result;
}

//----- (004A7A80) --------------------------------------------------------
int sub_4A7A80(const char* a1) {
	int result;        // eax
	unsigned int v2;   // ecx
	char v3;           // al
	unsigned char* v4; // edi
	const char* v5;    // esi

	if (!a1) {
		return 0;
	}
	v2 = strlen(a1) + 1;
	v3 = v2;
	v2 >>= 2;
	memcpy(getMemAt(0x5D4594, 1308644), a1, 4 * v2);
	v5 = &a1[4 * v2];
	v4 = getMemAt(0x5D4594, 1308644 + 4 * v2);
	LOBYTE(v2) = v3;
	result = 1;
	memcpy(v4, v5, v2 & 3);
	return result;
}

//----- (004A7AC0) --------------------------------------------------------
int sub_4A7AC0(const char* a1) {
	int result;        // eax
	unsigned int v2;   // ecx
	char v3;           // al
	unsigned char* v4; // edi
	const char* v5;    // esi

	if (!a1) {
		return 0;
	}
	v2 = strlen(a1) + 1;
	v3 = v2;
	v2 >>= 2;
	memcpy(getMemAt(0x5D4594, 1308172), a1, 4 * v2);
	v5 = &a1[4 * v2];
	v4 = getMemAt(0x5D4594, 1308172 + 4 * v2);
	LOBYTE(v2) = v3;
	result = 1;
	memcpy(v4, v5, v2 & 3);
	return result;
}

//----- (004A7B00) --------------------------------------------------------
int sub_4A7B00(const char* a1) {
	int result;        // eax
	unsigned int v2;   // ecx
	char v3;           // al
	unsigned char* v4; // edi
	const char* v5;    // esi

	if (!a1) {
		return 0;
	}
	v2 = strlen(a1) + 1;
	v3 = v2;
	v2 >>= 2;
	memcpy(getMemAt(0x5D4594, 1308352), a1, 4 * v2);
	v5 = &a1[4 * v2];
	v4 = getMemAt(0x5D4594, 1308352 + 4 * v2);
	LOBYTE(v2) = v3;
	result = 1;
	memcpy(v4, v5, v2 & 3);
	return result;
}

//----- (004A7B40) --------------------------------------------------------
int sub_4A7B40(char* a1) {
	const char* v2;    // eax
	int v3;            // edi
	unsigned char* v4; // esi

	if (!a1) {
		return 0;
	}
	v2 = *(const char**)getMemAt(0x587000, 171856);
	v3 = 0;
	if (*getMemU32Ptr(0x587000, 171856)) {
		v4 = getMemAt(0x587000, 171856);
		while (nox_strcmpi(v2, a1)) {
			v2 = (const char*)*((uint32_t*)v4 + 2);
			v4 += 8;
			++v3;
			if (!v2) {
				return 1;
			}
		}
		*getMemU32Ptr(0x5D4594, 1308184) = *getMemU32Ptr(0x587000, 171860 + 8 * v3);
	}
	return 1;
}

//----- (004A7BA0) --------------------------------------------------------
int sub_4A7BA0(char* a1) {
	int result; // eax

	result = (int)a1;
	if (a1) {
		*getMemU32Ptr(0x5D4594, 1308740) = atoi(a1);
		result = 1;
	}
	return result;
}

//----- (004A7BC0) --------------------------------------------------------
int sub_4A7BC0(const char* a1) {
	int result;        // eax
	unsigned int v2;   // ecx
	char v3;           // al
	unsigned char* v4; // edi
	const char* v5;    // esi

	if (!a1) {
		return 0;
	}
	v2 = strlen(a1) + 1;
	v3 = v2;
	v2 >>= 2;
	memcpy(getMemAt(0x5D4594, 1308324), a1, 4 * v2);
	v5 = &a1[4 * v2];
	v4 = getMemAt(0x5D4594, 1308324 + 4 * v2);
	LOBYTE(v2) = v3;
	result = 1;
	memcpy(v4, v5, v2 & 3);
	return result;
}

//----- (004A7C00) --------------------------------------------------------
int sub_4A7C00(const char* a1) {
	int result;        // eax
	unsigned int v2;   // ecx
	char v3;           // al
	unsigned char* v4; // edi
	const char* v5;    // esi

	if (!a1) {
		return 0;
	}
	v2 = strlen(a1) + 1;
	v3 = v2;
	v2 >>= 2;
	memcpy(getMemAt(0x5D4594, 1308364), a1, 4 * v2);
	v5 = &a1[4 * v2];
	v4 = getMemAt(0x5D4594, 1308364 + 4 * v2);
	LOBYTE(v2) = v3;
	result = 1;
	memcpy(v4, v5, v2 & 3);
	return result;
}

//----- (004A7C40) --------------------------------------------------------
int sub_4A7C40(char* a1) {
	int result; // eax

	result = (int)a1;
	if (a1) {
		*getMemU32Ptr(0x5D4594, 1308188) = atoi(a1);
		result = 1;
	}
	return result;
}

//----- (004A7C60) --------------------------------------------------------
int sub_4A7C60(char* a1) {
	char* v1;   // eax
	char* v2;   // eax

	if (!a1) {
		return 0;
	}
	*getMemU32Ptr(0x5D4594, 1308736) = 0;
	*getMemU32Ptr(0x5D4594, 1308732) = 0;
	v1 = strtok(a1, ",\t\r\n");
	if (v1) {
		*getMemU32Ptr(0x5D4594, 1308732) = atoi(v1);
	}
	v2 = strtok(0, " \t\r\n");
	if (v2) {
		*getMemU32Ptr(0x5D4594, 1308736) = atoi(v2);
	}
	if (*getMemU32Ptr(0x5D4594, 1308732) && *getMemU32Ptr(0x5D4594, 1308736)) {
		return 1;
	} else {
		return 0;
	}
}

//----- (004A7CE0) --------------------------------------------------------
int sub_4A7CE0(char* a1) {
	int result; // eax

	result = (int)a1;
	if (a1) {
		*getMemU32Ptr(0x5D4594, 1308728) = atoi(a1);
		result = 1;
	}
	return result;
}

//----- (004A7D00) --------------------------------------------------------
int sub_4A7D00(const char* a1) {
	int result;        // eax
	unsigned int v2;   // ecx
	char v3;           // al
	unsigned char* v4; // edi
	const char* v5;    // esi

	if (!a1) {
		return 0;
	}
	result = 0;
	if (strlen(a1) <= 0x80) {
		v2 = strlen(a1) + 1;
		v3 = v2;
		v2 >>= 2;
		memcpy(getMemAt(0x5D4594, 1308192), a1, 4 * v2);
		v5 = &a1[4 * v2];
		v4 = getMemAt(0x5D4594, 1308192 + 4 * v2);
		LOBYTE(v2) = v3;
		result = 1;
		memcpy(v4, v5, v2 & 3);
	}
	return result;
}

//----- (004A7D50) --------------------------------------------------------
int sub_4A7D50(char* a1) {
	int result; // eax

	result = (int)a1;
	if (a1) {
		*getMemU32Ptr(0x5D4594, 1308348) = atoi(a1);
		result = 1;
	}
	return result;
}

//----- (004A7EF0) --------------------------------------------------------
char* sub_4A7EF0() { return (char*)getMemAt(0x5D4594, 1308732); }

//----- (004A9C80) --------------------------------------------------------
int nox_xxx_compassGenStrings_4A9C80() {
	int v0;            // edi
	unsigned char* v1; // esi
	int v2;            // edi
	unsigned char* v3; // esi
	char v5[64];       // [esp+8h] [ebp-40h]

	*getMemU32Ptr(0x5D4594, 1309664) = 0;
	v0 = 0;
	v1 = getMemAt(0x5D4594, 1309644);
	do {
		nox_sprintf(v5, "Compass%d", ++v0);
		*(uint32_t*)v1 = nox_xxx_gLoadImg_42F970(v5);
		v1 += 4;
	} while ((int)v1 < (int)getMemAt(0x5D4594, 1309660));
	v2 = 0;
	v3 = getMemAt(0x5D4594, 1309516);
	do {
		nox_sprintf(v5, "CompassMainArrow%d", ++v2);
		*(uint32_t*)v3 = nox_xxx_gLoadImg_42F970(v5);
		v3 += 4;
	} while ((int)v3 < (int)getMemAt(0x5D4594, 1309644));
	return 1;
}

//----- (004AB260) --------------------------------------------------------
int sub_4AB260() {
	*getMemU32Ptr(0x5D4594, 1309752) = nox_xxx_gLoadImg_42F970("DisconnectIcon");
	dword_5d4594_1309756 = nox_window_new(0, 136, nox_win_width - 50, nox_win_height / 2 + 3, 50, 50, 0);
	nox_xxx_wndSetIcon_46AE60(*(int*)&dword_5d4594_1309756, *getMemIntPtr(0x5D4594, 1309752));
	nox_window_set_all_funcs(*(uint32_t**)&dword_5d4594_1309756, 0, sub_4AB420, 0);
	dword_5d4594_1309748 = nox_new_window_from_file("discon.wnd", sub_4AB390);
	nox_xxx_wndSetWindowProc_46B300(*(int*)&dword_5d4594_1309748, sub_4AB340);
	sub_46B120(*(uint32_t**)&dword_5d4594_1309748, 0);
	nox_window_setPos_46A9B0(*(uint32_t**)&dword_5d4594_1309748,
							 nox_win_width / 2 - *(uint32_t*)(dword_5d4594_1309748 + 24) / 2,
							 nox_win_height / 2 - *(uint32_t*)(dword_5d4594_1309748 + 28) / 2);
	return 1;
}

//----- (004AB340) --------------------------------------------------------
int sub_4AB340(int a1, int a2, int a3, int a4) {
	if (a2 != 21) {
		return 0;
	}
	if (a3 == 1) {
		return 1;
	}
	if (a3 == 57) {
		nox_point mpos = nox_client_getMousePos_4309F0();
		nox_window_call_field_93(a1, 5, mpos.x | (mpos.y << 16), 0);
	}
	return 0;
}

//----- (004AB390) --------------------------------------------------------
int sub_4AB390(int a1, int a2, int* a3, int a4) {
	int v3;     // eax
	int result; // eax

	if (a2 == 23) {
		return 1;
	}
	if (a2 != 16391) {
		return 0;
	}
	v3 = nox_xxx_wndGetID_46B0A0(a3) - 576;
	if (!v3) {
		sub_43CF40();
		return 0;
	}
	if (v3 != 1) {
		return 0;
	}
	sub_446380();
	if (dword_5d4594_2650652 && sub_41E2F0() == 9) {
		sub_41F4B0();
		sub_41EC30();
		sub_446490(0);
		nox_xxx____setargv_4_44B000();
		sub_4AB4D0(0);
		result = 0;
	} else {
		sub_43B750();
		sub_4AB4D0(0);
		result = 0;
	}
	return result;
}

//----- (004AB420) --------------------------------------------------------
int sub_4AB420(int* a1) {
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

//----- (004AB470) --------------------------------------------------------
int sub_4AB470() {
	int result; // eax

	nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1309748);
	nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1309756);
	result = 0;
	dword_5d4594_1309756 = 0;
	dword_5d4594_1309748 = 0;
	return result;
}

//----- (004AB4A0) --------------------------------------------------------
int sub_4AB4A0(int a1) {
	int result; // eax

	if (a1) {
		result = nox_window_set_hidden(*(int*)&dword_5d4594_1309756, 0);
	} else {
		result = nox_window_set_hidden(*(int*)&dword_5d4594_1309756, 1);
	}
	return result;
}

//----- (004AB4D0) --------------------------------------------------------
int sub_4AB4D0(int a1) {
	int result; // eax

	if (a1) {
		nox_video_stopAllFades_44E040();
		nox_window_set_hidden(*(int*)&dword_5d4594_1309748, 0);
		nox_xxx_wndShowModalMB_46A8C0(*(int*)&dword_5d4594_1309748);
		sub_46C690(*(int*)&dword_5d4594_1309748);
		nox_xxx_windowFocus_46B500(*(int*)&dword_5d4594_1309748);
		result = nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1309748, 1);
	} else {
		nox_window_set_hidden(*(int*)&dword_5d4594_1309748, 1);
		nox_xxx_wnd_46C6E0(*(int*)&dword_5d4594_1309748);
		nox_xxx_windowFocus_46B500(0);
		result = nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1309748, 0);
	}
	return result;
}

//----- (004AEE30) --------------------------------------------------------
long long sub_4AEE30() {
	int v0;            // edi
	unsigned char* v1; // esi
	long long result;  // rax

	v0 = 0;
	v1 = getMemAt(0x5D4594, 1309840);
	do {
		result =
			(long long)(sin((double)(v0 + 192) * *getMemDoublePtr(0x581450, 9768) * *getMemDoublePtr(0x581450, 9760)) *
						*(double*)&qword_581450_9552);
		*(uint32_t*)v1 = result;
		v1 += 4;
		++v0;
	} while ((int)v1 < (int)getMemAt(0x5D4594, 1311120));
	return result;
}

//----- (004B7C40) --------------------------------------------------------
uint32_t* nox_xxx_netHandleSummonPacket_4B7C40(short a1, unsigned short* a2, unsigned short a3, unsigned char a4,
											   short a5) {
	int v5;           // eax
	uint32_t* result; // eax
	uint32_t* v7;     // edi
	uint32_t* v8;     // esi
	int v9;           // [esp-8h] [ebp-1Ch]
	int v10;          // [esp-4h] [ebp-18h]
	int v11;          // [esp+10h] [ebp-4h]

	v10 = a2[1];
	v9 = *a2;
	v5 = nox_xxx_getTTByNameSpriteMB_44CFC0("SummonEffect");
	result = (uint32_t*)nox_xxx_spriteLoadAdd_45A360_drawable(v5, v9, v10);
	v7 = result;
	if (result) {
		result = nox_new_drawable_for_thing(a3);
		v8 = result;
		if (result) {
			result[3] = *a2;
			result[4] = a2[1];
			*((uint8_t*)result + 297) = nox_xxx_math_509EA0(a4);
			HIWORD(v11) = a1;
			LOWORD(v11) = a5;
			v8[69] = 8;
			v7[108] = v8;
			v7[109] = v11;
			result = gameFrame();
			v7[79] = gameFrame();
		}
	}
	return result;
}

//----- (004B7EE0) --------------------------------------------------------
void sub_4B7EE0(short a1) {
	if (!*getMemU32Ptr(0x5D4594, 1313744)) {
		*getMemU32Ptr(0x5D4594, 1313744) = nox_xxx_getTTByNameSpriteMB_44CFC0("SummonEffect");
	}
	if (!dword_5d4594_1313740) {
		dword_5d4594_1313740 = nox_xxx_getTTByNameSpriteMB_44CFC0("BlueSpark");
	}

	int v2 = sub_45A060();
	if (!v2) {
		return;
	}

	while (*(uint32_t*)(v2 + 108) != *getMemU32Ptr(0x5D4594, 1313744) || *(uint16_t*)(v2 + 438) != a1) {
		v2 = nox_drawable_next_45A070(v2);
		if (!v2) {
			return;
		}
	}
	nox_xxx_makePointFxCli_499610(*(int*)&dword_5d4594_1313740, 50, 1000, 30, *(uint32_t*)(v2 + 12),
								  *(uint32_t*)(v2 + 16));
	nox_xxx_spriteDelete_45A4B0(*(uint64_t**)(v2 + 432));
	nox_xxx_spriteDeleteStatic_45A4E0_drawable(v2);
}

//----- (004B7F90) --------------------------------------------------------
int nox_xxx_spriteShieldLoad_4B7F90() {
	int result; // eax

	*getMemU32Ptr(0x5D4594, 1313748) = nox_xxx_getTTByNameSpriteMB_44CFC0("SphericalShieldNW");
	*getMemU32Ptr(0x5D4594, 1313752) = nox_xxx_getTTByNameSpriteMB_44CFC0("SphericalShieldN");
	*getMemU32Ptr(0x5D4594, 1313756) = nox_xxx_getTTByNameSpriteMB_44CFC0("SphericalShieldNE");
	*getMemU32Ptr(0x5D4594, 1313760) = nox_xxx_getTTByNameSpriteMB_44CFC0("SphericalShieldW");
	*getMemU32Ptr(0x5D4594, 1313768) = nox_xxx_getTTByNameSpriteMB_44CFC0("SphericalShieldE");
	*getMemU32Ptr(0x5D4594, 1313772) = nox_xxx_getTTByNameSpriteMB_44CFC0("SphericalShieldSW");
	*getMemU32Ptr(0x5D4594, 1313776) = nox_xxx_getTTByNameSpriteMB_44CFC0("SphericalShieldS");
	result = nox_xxx_getTTByNameSpriteMB_44CFC0("SphericalShieldSE");
	*getMemU32Ptr(0x5D4594, 1313780) = result;
	*getMemU32Ptr(0x5D4594, 1313764) = 0;
	*getMemU32Ptr(0x5D4594, 1313784) = 1;
	return result;
}
// 4B8040: variable 'v3' is possibly undefined
// 4B804D: variable 'v5' is possibly undefined

//----- (004B8090) --------------------------------------------------------
uint32_t* nox_xxx_fxShield_4B8090(unsigned int a1, int a2) {
	int v2;           // edi
	int v3;           // eax
	uint32_t* result; // eax
	int v5;           // eax
	uint32_t* v6;     // esi
	int4 v7;          // [esp+0h] [ebp-10h]

	if (!*getMemU32Ptr(0x5D4594, 1313784)) {
		nox_xxx_spriteShieldLoad_4B7F90();
	}
	v2 = a2;
	switch (a2) {
	case 0:
		nox_xxx_getTTByNameSpriteMB_44CFC0("SphericalShieldNW");
		break;
	case 1:
		nox_xxx_getTTByNameSpriteMB_44CFC0("ShpericalShieldN");
		break;
	case 2:
		nox_xxx_getTTByNameSpriteMB_44CFC0("SphericalShieldNW");
		break;
	case 3:
		nox_xxx_getTTByNameSpriteMB_44CFC0("SphericalShieldW");
		break;
	case 5:
		nox_xxx_getTTByNameSpriteMB_44CFC0("SphericalShieldE");
		break;
	case 6:
		nox_xxx_getTTByNameSpriteMB_44CFC0("SphericalShieldSW");
		break;
	case 7:
		nox_xxx_getTTByNameSpriteMB_44CFC0("SphericalShieldS");
		break;
	case 8:
		nox_xxx_getTTByNameSpriteMB_44CFC0("SphericalShieldSE");
		break;
	default:
		break;
	}
	if (nox_xxx_netTestHighBit_578B70(a1)) {
		v3 = nox_xxx_netClearHighBit_578B30(a1);
		result = nox_xxx_netSpriteByCodeStatic_45A720(v3);
	} else {
		v5 = nox_xxx_netClearHighBit_578B30(a1);
		result = nox_xxx_netSpriteByCodeDynamic_45A6F0(v5);
	}
	v6 = result;
	if (result) {
		v7.field_0 = result[3] - 10;
		v7.field_4 = result[4] - 10;
		v7.field_8 = result[3] + 10;
		v7.field_C = result[4] + 10;
		dword_5d4594_1313788 = 0;
		nox_xxx_forEachSprite_49AB00(&v7, nox_xxx_spriteScanForShield_4B81E0, (int)&a1);
		result = *(uint32_t**)&dword_5d4594_1313788;
		if (dword_5d4594_1313788 != 1) {
			result = (uint32_t*)nox_xxx_spriteLoadAdd_45A360_drawable(*getMemU32Ptr(0x5D4594, 1313748 + 4 * v2), v6[3],
																	  v6[4] + 3);
			if (result) {
				result[108] = a1;
			}
		}
	}
	return result;
}
// 4B810D: variable 'v3' is possibly undefined
// 4B811F: variable 'v5' is possibly undefined

//----- (004B81E0) --------------------------------------------------------
void nox_xxx_spriteScanForShield_4B81E0(int a1, int a2) {
	unsigned char* v2; // eax

	v2 = getMemAt(0x5D4594, 1313748);
	do {
		if (*(uint32_t*)(a1 + 108) == *(uint32_t*)v2 && *(uint32_t*)(a1 + 432) == *(uint32_t*)a2) {
			dword_5d4594_1313788 = 1;
		}
		v2 += 4;
	} while ((int)v2 < (int)getMemAt(0x5D4594, 1313784));
}

//----- (004B8E10) --------------------------------------------------------
uint32_t* sub_4B8E10(uint32_t* a1, char* a2) {
	uint32_t* result; // eax
	int v3;           // ebx
	int v4;           // eax
	int v5;           // edx
	int v6;           // ecx
	uint32_t* v7;     // ebp
	int v8;           // edi
	uint8_t* v9;      // esi
	uint32_t* v10;    // ecx
	int v11;          // eax
	int* v12;         // edi
	int v13;          // ebx
	uint32_t** v14;   // esi
	int v15;          // [esp-Ch] [ebp-14h]

	result = a1;
	v3 = 0;
	while ((char*)*result != a2) {
		++v3;
		result += 6;
		if (v3 >= 27) {
			return result;
		}
	}
	v4 = sub_415840(a2);
	result = nox_xxx_getProjectileClassById_413250(v4);
	v7 = result;
	if (result) {
		v8 = 1;
		v9 = result + 4;
		do {
			LOBYTE(result) = v9[1];
			LOBYTE(v6) = *v9;
			LOBYTE(v5) = *(v9 - 1);
			nox_draw_setMaterial_4340A0(v8++, v5, v6, (int)result);
			v9 += 3;
		} while (v8 < 7);
		v10 = a1;
		v11 = 3 * v3;
		v12 = v7 + 9;
		v13 = 4;
		v14 = (uint32_t**)&a1[2 * v11 + 1];
		do {
			result = *v14;
			if (*v14) {
				LOBYTE(v5) = *((uint8_t*)result + 26);
				LOBYTE(v10) = *((uint8_t*)result + 25);
				v15 = v5;
				LOBYTE(v5) = *((uint8_t*)result + 24);
				nox_draw_setMaterial_4340A0(*v12, v5, (int)v10, v15);
			}
			++v14;
			++v12;
			--v13;
		} while (v13);
	}
	return result;
}
// 4B8E57: variable 'v5' is possibly undefined
// 4B8E57: variable 'v6' is possibly undefined
// 4B8E90: variable 'v10' is possibly undefined
