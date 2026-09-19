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
extern uint32_t dword_5d4594_1313788;
extern uint32_t nox_xxx_normalWndBits_587000_172880;
extern uint32_t dword_5d4594_1313532;
extern uint32_t dword_5d4594_1313564;
extern uint32_t nox_server_sendMotd_108752;
extern uint32_t dword_5d4594_1313536;
extern uint32_t dword_5d4594_1313740;
extern uint32_t dword_5d4594_1309736;
extern uint32_t dword_5d4594_1309756;
extern uint32_t dword_5d4594_1309832;
extern uint32_t dword_5d4594_1313540;
extern uint32_t dword_5d4594_1309824;
extern uint32_t nox_server_connectionType_3596;
extern void* dword_587000_122852;
extern uint32_t dword_5d4594_1309828;
extern uint32_t dword_5d4594_1309836;
extern uint32_t dword_5d4594_1309728;
extern uint32_t dword_5d4594_1309732;
extern uint64_t qword_581450_9552;
extern void* dword_587000_93164;
extern void* dword_587000_127004;
extern uint32_t dword_5d4594_1309748;
extern uint32_t dword_5d4594_1309720;
extern uint32_t dword_5d4594_1309820;
extern uint32_t dword_5d4594_2650652;
extern int nox_win_width;
extern int nox_win_height;

extern uint32_t nox_color_white_2523948;
extern uint32_t nox_color_blue_2650684;
extern uint32_t nox_color_violet_2598268;
extern uint32_t nox_color_black_2650656;

nox_gui_animation* nox_wnd_xxx_1309740 = 0;


//----- (004A2560) --------------------------------------------------------


//----- (004A25C0) --------------------------------------------------------


//----- (004A2610) --------------------------------------------------------


//----- (004A2830) --------------------------------------------------------


//----- (004A2890) --------------------------------------------------------


//----- (004A28B0) --------------------------------------------------------


//----- (004A28C0) --------------------------------------------------------


//----- (004A7EF0) --------------------------------------------------------


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
