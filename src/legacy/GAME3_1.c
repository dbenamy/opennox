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
extern uint32_t dword_5d4594_1320964;
extern uint32_t dword_5d4594_1321228;
extern uint32_t dword_5d4594_1321040;
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

// 487CF0: using guessed type void  nullsub_10(uint32_t);

// 487CA0: using guessed type void  nullsub_9(uint32_t);

//----- (004BF010) --------------------------------------------------------


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


//----- (004BFBB0) --------------------------------------------------------


//----- (004BFBF0) --------------------------------------------------------


//----- (004BFC70) --------------------------------------------------------


//----- (004BFC90) --------------------------------------------------------


//----- (004BFCD0) --------------------------------------------------------


//----- (004BFD10) --------------------------------------------------------


//----- (004BFD30) --------------------------------------------------------


//----- (004C3390) --------------------------------------------------------


//----- (004C3410) --------------------------------------------------------


//----- (004C3460) --------------------------------------------------------


//----- (004C34A0) --------------------------------------------------------


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
