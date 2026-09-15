#include "GAME2_1.h"
#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME1_2.h"
#include "GAME1_3.h"
#include "GAME2.h"
#include "GAME2_2.h"
#include "GAME2_3.h"
#include "GAME3.h"
#include "GAME3_1.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4_1.h"
#include "GAME5_2.h"
#include "client__draw__debugdraw.h"
#include "client__draw__staticdraw.h"
#include "client__drawable__drawable.h"
#include "client__gui__window.h"
#include "common__net_list.h"
#include "common__object__modifier.h"
#include "common__system__team.h"
#include "operators.h"

#include "client__gui__gamewin__gamewin.h"
#include "client__gui__gui_ctf.h"
#include "client__gui__guicon.h"
#include "client__gui__guiinput.h"
#include "client__gui__guiinv.h"
#include "client__gui__guijourn.h"
#include "client__gui__guimeter.h"
#include "client__gui__guiquit.h"
#include "client__gui__guirank.h"
#include "client__gui__guisave.h"
#include "client__gui__guispell.h"
#include "client__gui__guisumn.h"
#include "client__gui__guitrade.h"
#include "client__gui__tooltip.h"
#include "client__video__draw_common.h"

#include "client__system__ctrlevnt.h"
#include "common/fs/nox_fs.h"
#include "common__magic__speltree.h"
#include "input_common.h"

extern uint32_t dword_8531A0_2576;
extern uint32_t dword_8531A0_2572;
extern uint32_t dword_5d4594_1062552;
extern uint32_t dword_5d4594_1049844;
extern uint32_t dword_5d4594_1050008;
extern uint32_t dword_5d4594_1096272;
extern uint32_t dword_5d4594_1049516;
extern uint32_t dword_5d4594_1062520;
extern uint32_t dword_5d4594_1096276;
extern uint32_t dword_5d4594_1049976;
extern uint32_t dword_5d4594_1090284;
extern uint32_t dword_5d4594_1062484;
extern uint32_t dword_5d4594_1049992;
extern uint32_t dword_5d4594_1062556;
extern uint32_t dword_5d4594_1096264;
extern uint32_t dword_5d4594_1062564;
extern uint32_t dword_5d4594_1090280;
extern uint32_t dword_5d4594_1062560;
extern uint32_t dword_5d4594_1096280;
extern uint32_t dword_587000_145672;
extern uint32_t dword_5d4594_1064868;
extern uint32_t dword_5d4594_1049996;
extern uint32_t dword_5d4594_1062492;
extern uint32_t dword_5d4594_1062496;
extern uint32_t dword_5d4594_1096260;
extern uint32_t nox_client_gui_flag_1556112;
extern uint32_t dword_5d4594_1096284;
extern uint32_t dword_587000_145664;
extern uint32_t dword_5d4594_1096288;
extern uint32_t dword_5d4594_1047936;
extern uint32_t dword_5d4594_1062488;
extern uint32_t dword_5d4594_1062468;
extern uint32_t dword_5d4594_1064860;
extern uint32_t dword_5d4594_1096252;
extern uint32_t dword_5d4594_1064864;
extern uint32_t dword_587000_136184;
extern uint32_t dword_5d4594_1049532;
extern uint32_t dword_5d4594_1049484;
extern uint32_t dword_5d4594_1062516;
extern uint32_t nox_client_translucentFrontWalls_805844;
extern uint64_t qword_581450_9512;
extern uint64_t qword_581450_9544;
extern uint32_t dword_5d4594_1090276;
extern uint32_t dword_5d4594_1062476;
extern uint32_t dword_5d4594_1047932;
extern uint32_t dword_5d4594_1049512;
extern uint32_t nox_client_highResFloors_154952;
extern uint32_t dword_5d4594_1062528;
extern uint32_t dword_5d4594_1062524;
extern uint32_t dword_5d4594_1064856;
extern uint32_t dword_5d4594_1049856;
extern uint32_t dword_5d4594_1049520;
extern uint32_t nox_client_highResFrontWalls_80820;
extern uint32_t dword_5d4594_1049800_inventory_click_row_index;
extern uint32_t nox_xxx_minimap_587000_149232;
extern uint32_t dword_5d4594_1062456;
extern uint32_t dword_5d4594_1063636;
extern uint32_t dword_5d4594_1049796_inventory_click_column_index;
extern uint32_t dword_5d4594_1049508;
extern nox_window* dword_5d4594_1090048;
extern uint32_t dword_5d4594_1049500;
extern uint32_t dword_5d4594_1062512;
extern uint32_t dword_5d4594_1049864;
extern uint32_t dword_5d4594_1062508;
extern uint32_t dword_5d4594_1049504;
extern uint32_t dword_5d4594_1090120;
extern uint32_t dword_5d4594_1063116;
extern uint32_t dword_5d4594_1062480;
extern uint32_t nox_player_netCode_85319C;
extern void* nox_xxx_aClosewoodengat_587000_133480;
extern int nox_win_width;
extern int nox_win_height;
extern uint32_t array_5D4594_1049872[9];

extern uint32_t nox_color_white_2523948;
extern uint32_t nox_color_blue_2650684;
extern uint32_t nox_color_yellow_2589772;
extern uint32_t nox_color_violet_2598268;
extern uint32_t nox_color_black_2650656;

nox_window* nox_win_unk5 = 0;
nox_window* dword_5d4594_1062452 = 0;


obj_5D4594_2650668_t** ptr_5D4594_2650668 = 0;
const int ptr_5D4594_2650668_cap = 128;

nox_inventory_cell_t nox_client_inventory_grid_1050020[NOX_INVENTORY_CELLS_MAX] = {0};

//----- (00460D40) --------------------------------------------------------
int sub_460D40() { return dword_5d4594_1049508 != 0; }

//----- (00460D50) --------------------------------------------------------
int sub_460D50() {
	unsigned char* v0; // edi
	uint32_t** v1;     // esi
	int v2;            // ebx
	unsigned char* v3; // esi
	int result;        // eax

	nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1049500);
	dword_5d4594_1049500 = 0;
	nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1049504);
	dword_5d4594_1049504 = 0;
	nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1049520);
	dword_5d4594_1049520 = 0;
	nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1049508);
	dword_5d4594_1049508 = 0;
	nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1049512);
	dword_5d4594_1049512 = 0;
	nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1049516);
	dword_5d4594_1049516 = 0;
	v0 = getMemAt(0x5D4594, 1048404);
	do {
		nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)v0);
		*(uint32_t*)v0 = 0;
		v1 = (uint32_t**)(v0 + 24);
		v2 = 5;
		do {
			nox_xxx_windowDestroyMB_46C4E0(*(v1 - 5));
			*(v1 - 5) = 0;
			nox_xxx_windowDestroyMB_46C4E0(*v1);
			*v1 = 0;
			++v1;
			--v2;
		} while (v2);
		v0 += 256;
	} while ((int)v0 < (int)getMemAt(0x5D4594, 1049684));
	nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)getMemAt(0x5D4594, 1048148));
	*getMemU32Ptr(0x5D4594, 1048148) = 0;
	v3 = getMemAt(0x5D4594, 1048152);
	do {
		result = nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)v3);
		*(uint32_t*)v3 = 0;
		v3 += 4;
	} while ((int)v3 < (int)getMemAt(0x5D4594, 1048164));
	dword_5d4594_1049532 = 0;
	*getMemU32Ptr(0x5D4594, 1047928) = 0;
	dword_5d4594_1047932 = 0;
	return result;
}

//----- (00460E60) --------------------------------------------------------
int nox_xxx_cliPrepareGameplay1_460E60() {
	int result; // eax

	if (sub_460D40()) {
		sub_460D50();
	}
	result = nox_xxx_quickBarCreate_45E190();
	if (result) {
		sub_460EA0(nox_client_getRenderGUI());
		result = 1;
	}
	return result;
}

//----- (00460EA0) --------------------------------------------------------
int sub_460EA0(int a1) { return sub_460B90(a1); }

//----- (00460EB0) --------------------------------------------------------
void sub_460EB0(int a1, char a2) {
	if (a1 < 0 || a1 >= 140) {
		return;
	}
	*getMemU8Ptr(0x5D4594, 1049544 + a1) = a2;
}

//----- (00461010) --------------------------------------------------------
void sub_461010() {
	if (!dword_5d4594_1049484) {
		return;
	}
	nox_window_set_hidden(*getMemIntPtr(0x5D4594, 1048148), 1);
	nox_window_set_hidden(*(int*)&dword_5d4594_1049512, 0);
	sub_46AE10(*(int*)&dword_5d4594_1049500, 0);
	nox_xxx_clientPlaySoundSpecial_452D80(797, 100);
	dword_5d4594_1049484 = 0;
}

//----- (00461060) --------------------------------------------------------
void sub_461060() {
	if (dword_5d4594_1049484 == 1) {
		sub_461010();
		return;
	}
	if (*getMemU32Ptr(0x5D4594, 1049476) == 1) {
		nox_xxx_quickBarClose_4606B0();
	}
	nox_window_set_hidden(*getMemIntPtr(0x5D4594, 1048148), 0);
	nox_window_set_hidden(*(int*)&dword_5d4594_1049512, 1);
	sub_46AE10(*(int*)&dword_5d4594_1049500, 1);
	nox_xxx_clientPlaySoundSpecial_452D80(796, 100);
	dword_5d4594_1049484 = 1;
}

//----- (00461090) --------------------------------------------------------
char* sub_461090(int a1, int a2) {
	int v2;       // edx
	char* result; // eax

	v2 = gameFrame();
	result = (char*)getMemAt(0x5D4594, 1047764 + 24*1 + 20);
	do {
		if (*((uint32_t*)result - 5) == a1) {
			*(uint32_t*)result = a2 == 0 ? v2 : 0;
			*((uint32_t*)result - 3) = a2;
		}
		result += 24;
	} while ((int)result < (int)getMemAt(0x5D4594, 1047928));
	return result;
}

//----- (004610D0) --------------------------------------------------------
char* sub_4610D0(unsigned char a1) {
	int* v1;      // esi
	char* result; // eax

	if (a1 != 6) {
		return sub_461090(*getMemU32Ptr(0x5D4594, 1047764 + 24*a1), 1);
	}
	v1 = getMemIntPtr(0x5D4594, 1047764 + 24*1);
	do {
		result = sub_461090(*v1, 1);
		v1 += 6;
	} while ((int)v1 < (int)getMemAt(0x5D4594, 1047908));
	return result;
}

//----- (00461120) --------------------------------------------------------
char* sub_461120(int a1, int a2) {
	int v2;       // edx
	char* result; // eax

	v2 = 1 << a1;
	result = (char*)getMemAt(0x5D4594, 1047764 + 24*1 + 12);
	do {
		if (*((uint32_t*)result - 3) == a1) {
			if (a2) {
				*(uint32_t*)result |= v2;
			} else {
				*(uint32_t*)result &= ~v2;
			}
		}
		result += 24;
	} while ((int)result < (int)getMemAt(0x5D4594, 1047920));
	return result;
}

//----- (00461160) --------------------------------------------------------
int sub_461160(int a1) {
	int v1;            // edx
	unsigned char* v2; // eax

	v1 = 1;
	v2 = getMemAt(0x5D4594, 1047764 + 24*1);
	while (*(uint32_t*)v2 != a1) {
		v2 += 24;
		++v1;
		if ((int)v2 >= (int)getMemAt(0x5D4594, 1047908)) {
			return 0;
		}
	}
	return ((1 << a1) & *getMemU32Ptr(0x5D4594, 1047764 + 24*v1 + 12)) != 0;
}

//----- (004611A0) --------------------------------------------------------
int sub_4611A0() { return dword_5d4594_1047932; }

//----- (004611B0) --------------------------------------------------------
int sub_4611B0() {
	int result; // eax

	result = dword_5d4594_1047936;
	if (dword_5d4594_1047936) {
		result = nox_xxx_clientSendAbil_45DAF0(*(int*)&dword_5d4594_1047936);
		dword_5d4594_1047936 = 0;
		dword_5d4594_1047932 = 0;
	}
	return result;
}

//----- (004611E0) --------------------------------------------------------
void nox_xxx_netAbilityRewardCli_4611E0(int a1, int a2, char* a3) {
	unsigned char* v3; // esi

	if (a1 >= 1 && a1 < 6) {
		v3 = getMemAt(0x5D4594, 1047764 + 24*1 + 16);
		do {
			if (*((uint32_t*)v3 - 4) == a1 && *(uint32_t*)v3 != a2) {
				if (nox_common_gameFlags_check_40A5C0(2) && dword_8531A0_2576) {
					*(uint32_t*)(dword_8531A0_2576 + 4 * a1 + 3696) = a2;
				}
				*(uint32_t*)v3 = a2;
				if (a2) {
					nox_xxx_abilityReward_45D290(a1, a3, (int)a3);
				}
			}
			v3 += 24;
		} while ((int)v3 < (int)getMemAt(0x5D4594, 1047924));
	}
}

//----- (00461250) --------------------------------------------------------
int nox_xxx_buttonFindFirstEmptySlot_461250() {
	int v0;       // ecx
	int v1;       // esi
	uint32_t* v2; // eax

	v0 = *(unsigned char*)((uint32_t)nox_xxx_aClosewoodengat_587000_133480 + 200);
	do {
		v1 = 0;
		v2 = (uint32_t*)((uint32_t)nox_xxx_aClosewoodengat_587000_133480 + 40 * v0);
		do {
			if (!*v2) {
				nox_xxx_clientUpdateButtonRow_45E110(v0);
				return v1;
			}
			++v1;
			v2 += 2;
		} while (v1 < 5);
		if (++v0 >= 5) {
			v0 = 0;
		}
	} while (v0 != *(unsigned char*)((uint32_t)nox_xxx_aClosewoodengat_587000_133480 + 200));
	return -1;
}

//----- (004612A0) --------------------------------------------------------
int sub_4612A0() {
	int result;  // eax
	uint32_t* i; // ecx

	result = 0;
	for (i = (uint32_t*)((uint32_t)nox_xxx_aClosewoodengat_587000_133480 +
						 40 * *(unsigned char*)((uint32_t)nox_xxx_aClosewoodengat_587000_133480 + 200));
		 *i; i += 2) {
		if (++result >= 5) {
			return -1;
		}
	}
	return result;
}

//----- (004612D0) --------------------------------------------------------
int nox_xxx_buttonHaveSpellInBarMB_4612D0(int a1) {
	int v1;       // edx
	int v2;       // eax
	uint32_t* v3; // ecx

	v1 = *(unsigned char*)((uint32_t)nox_xxx_aClosewoodengat_587000_133480 + 200);
	do {
		v2 = 0;
		v3 = (uint32_t*)((uint32_t)nox_xxx_aClosewoodengat_587000_133480 + 40 * v1);
		do {
			if (*v3 == a1) {
				return 1;
			}
			++v2;
			v3 += 2;
		} while (v2 < 5);
		if (++v1 >= 5) {
			v1 = 0;
		}
	} while (v1 != *(unsigned char*)((uint32_t)nox_xxx_aClosewoodengat_587000_133480 + 200));
	return 0;
}

//----- (00461320) --------------------------------------------------------
void nox_xxx_buttonSetImgMB_461320(int a1, uint32_t* a2) {
	if (a2) {
		if (a1 >= 0 && a1 < 5) {
			nox_client_wndGetPosition_46AA60(
				*(uint32_t**)((uint32_t)nox_xxx_aClosewoodengat_587000_133480 + 4 * a1 + 212), a2, a2 + 1);
		}
	}
}

//----- (00461360) --------------------------------------------------------
int sub_461360(int a1) {
	int v1;     // edx
	int v2;     // ecx
	int v3;     // ebx
	int v4;     // esi
	int result; // eax

	v1 = nox_xxx_aClosewoodengat_587000_133480;
	v2 = *(unsigned char*)((uint32_t)nox_xxx_aClosewoodengat_587000_133480 + 200);
	v3 = *(unsigned char*)((uint32_t)nox_xxx_aClosewoodengat_587000_133480 + 200);
	do {
		v4 = 5;
		result = 40 * v2;
		do {
			if (*(uint32_t*)(result + v1) == a1) {
				*(uint32_t*)(result + v1) = 0;
				v1 = nox_xxx_aClosewoodengat_587000_133480;
			}
			result += 8;
			--v4;
		} while (v4);
		if (++v2 >= 5) {
			v2 = 0;
		}
	} while (v2 != v3);
	return result;
}

//----- (00461400) --------------------------------------------------------
int sub_461400() {
	int i;      // esi
	int result; // eax
	int v2;     // ecx

	for (i = 0; i < 40; i += 8) {
		result = i;
		v2 = 5;
		do {
			*(uint32_t*)(result + (uint32_t)nox_xxx_aClosewoodengat_587000_133480) =
				*getMemU32Ptr(0x5D4594, 1047564 + result);
			*(uint8_t*)(result + (uint32_t)nox_xxx_aClosewoodengat_587000_133480 + 4) =
				getMemByte(0x5D4594, 1047568 + result);
			result += 40;
			--v2;
		} while (v2);
	}
	return result;
}

//----- (00461440) --------------------------------------------------------
int sub_461440(int a1) {
	int result; // eax

	result = a1;
	*getMemU32Ptr(0x5D4594, 1049688) = a1;
	return result;
}

//----- (00461450) --------------------------------------------------------
int sub_461450() { return *getMemU32Ptr(0x5D4594, 1049688); }

//----- (00461460) --------------------------------------------------------
void nox_xxx_playerInitColors_461460(nox_playerInfo* pl) {
	int a1 = pl;
	int v1;     // eax
	char v2;    // cl
	char v3;    // dl
	int v4;     // eax
	char v5;    // cl
	char v6;    // dl
	int v7;     // eax
	char v8;    // cl
	char v9;    // dl
	int v10;    // eax
	char v11;   // cl
	char v12;   // dl

	v1 = nox_color_rgb_4344A0(*(uint8_t*)(a1 + 2253), *(uint8_t*)(a1 + 2254), *(uint8_t*)(a1 + 2255));
	v2 = *(uint8_t*)(a1 + 2257);
	v3 = *(uint8_t*)(a1 + 2256);
	*(uint32_t*)(a1 + 2296) = v1;
	v4 = nox_color_rgb_4344A0(v3, v2, *(uint8_t*)(a1 + 2258));
	v5 = *(uint8_t*)(a1 + 2260);
	v6 = *(uint8_t*)(a1 + 2259);
	*(uint32_t*)(a1 + 2292) = v4;
	v7 = nox_color_rgb_4344A0(v6, v5, *(uint8_t*)(a1 + 2261));
	v8 = *(uint8_t*)(a1 + 2263);
	v9 = *(uint8_t*)(a1 + 2262);
	*(uint32_t*)(a1 + 2300) = v7;
	v10 = nox_color_rgb_4344A0(v9, v8, *(uint8_t*)(a1 + 2264));
	v11 = *(uint8_t*)(a1 + 2266);
	v12 = *(uint8_t*)(a1 + 2265);
	*(uint32_t*)(a1 + 2304) = v10;
	*(uint32_t*)(a1 + 2308) = nox_color_rgb_4344A0(v12, v11, *(uint8_t*)(a1 + 2267));
	*(uint32_t*)(a1 + 2312) = nox_color_white_2523948;
}

//----- (00461520) --------------------------------------------------------
char* sub_461520() {
	char* result; // eax
	int i;        // esi

	result = nox_common_playerInfoGetFirst_416EA0();
	for (i = (int)result; result; i = (int)result) {
		nox_xxx_playerInitColors_461460(i);
		result = nox_common_playerInfoGetNext_416EE0(i);
	}
	return result;
}

//----- (00469B90) --------------------------------------------------------
int sub_469B90(int* a1) {
	int result; // eax

	*getMemU32Ptr(0x587000, 142296) = *a1;
	*getMemU32Ptr(0x587000, 142300) = a1[1];
	result = a1[2];
	*getMemU32Ptr(0x587000, 142304) = a1[2];
	return result;
}

//----- (00469BB0) --------------------------------------------------------
char* nox_xxx_getAmbientColor_469BB0() { return (char*)getMemAt(0x587000, 142296); }

//----- (00469FA0) --------------------------------------------------------
int sub_469FA0() { return *getMemU32Ptr(0x5D4594, 1064848); }

//----- (0046A430) --------------------------------------------------------
void nox_client_chatStart_46A430(int a1) {
	if (!nox_common_gameFlags_check_40A5C0(2048)) {
		if (!dword_5d4594_1064868) {
			**(uint16_t**)&dword_5d4594_1064864 = 0;
			*(uint16_t*)(dword_5d4594_1064864 + 1052) = 0;
			nox_xxx_wndShowModalMB_46A8C0(*(int*)&dword_5d4594_1064856);
			sub_46C690(*(int*)&dword_5d4594_1064856);
			nox_xxx_windowFocus_46B500(*(int*)&dword_5d4594_1064860);
			dword_5d4594_1064868 = 1;
			*getMemU32Ptr(0x5D4594, 1064872) = a1;
		}
	}
}

//----- (0046A4A0) --------------------------------------------------------
int sub_46A4A0() { return dword_5d4594_1064868; }

//----- (0046A4B0) --------------------------------------------------------
size_t nox_xxx_cmdSayDo_46A4B0(wchar2_t* a1, int a2) {
	uint32_t* v2;      // ebp
	size_t v3;         // edi
	size_t result;     // eax
	const wchar2_t* v5; // edi
	char v6;           // al
	int v7;            // eax
	char v8[520];      // [esp+Ch] [ebp-208h]

	v2 = nox_xxx_netSpriteByCodeDynamic_45A6F0(nox_player_netCode_85319C);
	v3 = nox_wcsspn(a1, L" ");
	result = nox_wcslen(a1);
	if (v3 != result) {
		v5 = &a1[v3];
		v8[0] = -88; // MSG_TEXT_MESSAGE
		*(uint16_t*)&v8[9] = 0;
		*(uint16_t*)&v8[1] = nox_player_netCode_85319C;
		v8[3] = 0;
		if (nox_xxx_cliCanTalkMB_4100F0((short*)a1)) {
			v6 = v8[3] | 2;
		} else {
			v6 = v8[3] | 4;
		}
		v8[3] = v6;
		if (a2) {
			v8[3] |= 1u;
		}
		v8[8] = nox_wcslen(v5) + 1;
		if (v8[3] & 4) {
			nox_wcscpy((wchar2_t*)&v8[11], v5);
			v7 = 2;
		} else {
			nox_sprintf(&v8[11], "%S", v5);
			v7 = 1;
		}
		if (v2) {
			*(uint16_t*)&v8[4] = *((uint16_t*)v2 + 6);
			*(uint16_t*)&v8[6] = *((uint16_t*)v2 + 8);
		} else {
			*(uint16_t*)&v8[6] = -1;
			*(uint16_t*)&v8[4] = -1;
		}
		result = nox_netlist_addToMsgListCli_40EBC0(31, 0, v8, v7 * (unsigned char)v8[8] + 11);
	}
	return result;
}

//----- (0046A5D0) --------------------------------------------------------
int sub_46A5D0(uint32_t* a1, int a2) {
	int v2;  // ecx
	bool v3; // sf
	int v5;  // [esp+4h] [ebp-8h]
	int v6;  // [esp+8h] [ebp-4h]

	v5 = 0;
	v6 = 0;
	nox_xxx_wndShowModalMB_46A8C0(*(int*)&dword_5d4594_1064856);
	nox_xxx_windowFocus_46B500(*(int*)&dword_5d4594_1064860);
	nox_xxx_drawGetStringSize_43F840(0, *(unsigned short**)&dword_5d4594_1064864, &v5, 0, 0);
	nox_xxx_drawGetStringSize_43F840(0, (unsigned short*)(dword_5d4594_1064864 + 512), &v6, 0, 0);
	v3 = v5 + v6 - 90 < 0;
	v5 += v6 + 10;
	v2 = v5;
	if (v5 < 100) {
		v2 = 100;
		v5 = v2;
	} else if (v5 > 320) {
		v2 = 320;
		v5 = v2;
	}
	nox_window_setPos_46A9B0(*(uint32_t**)&dword_5d4594_1064856, (nox_win_width - v2) / 2, *(uint32_t*)(dword_5d4594_1064856 + 20));
	sub_46AB20(a1, v5, 20);
	return nox_xxx_wndEditDrawNoImage_488160((int)a1, a2);
}

//----- (0046A6A0) --------------------------------------------------------
int sub_46A6A0() {
	if (wndIsShown_nox_xxx_wndIsShown_46ACC0(*(int*)&dword_5d4594_1064856)) {
		return 0;
	}
	if (nox_xxx_wndGetFocus_46B4F0() == dword_5d4594_1064860) {
		nox_xxx_windowFocus_46B500(0);
	}
	nox_xxx_wnd_46C6E0(*(int*)&dword_5d4594_1064856);
	nox_window_set_hidden(*(int*)&dword_5d4594_1064856, 1);
	*(uint32_t*)(dword_5d4594_1064856 + 4) &= 0xFFFFFFF7;
	*(uint32_t*)(dword_5d4594_1064860 + 4) &= 0xFFFFFFF7;
	set_dword_5d4594_3799468(1);
	dword_5d4594_1064868 = 0;
	return 1;
}

//----- (0046A730) --------------------------------------------------------
uint32_t* sub_46A730() {
	uint32_t* result; // eax

	*getMemU32Ptr(0x5D4594, 1064876) = nox_win_width / 2;
	*getMemU32Ptr(0x5D4594, 1064880) = 2 * nox_win_height / 3;
	result = nox_new_window_from_file("GuiChat.wnd", sub_46A820);
	dword_5d4594_1064856 = result;
	if (result) {
		nox_window_setPos_46A9B0(result, *getMemIntPtr(0x5D4594, 1064876), *getMemIntPtr(0x5D4594, 1064880));
		result = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1064856, 9201);
		dword_5d4594_1064860 = result;
		if (result) {
			nox_xxx_wndSetDrawFn_46B340((int)result, sub_46A5D0);
			nox_xxx_wndSetWindowProc_46B300(*(int*)&dword_5d4594_1064860, sub_46A7E0);
			result = *(uint32_t**)&dword_5d4594_1064856;
			dword_5d4594_1064864 = *(uint32_t*)(dword_5d4594_1064860 + 32);
		}
	}
	return result;
}

//----- (0046A7E0) --------------------------------------------------------
int sub_46A7E0(uint32_t* a1, int a2, int a3, int a4) {
	if (a2 != 21 || a3 != 1) {
		return nox_xxx_wndEditProc_487D70(a1, a2, a3, a4);
	}
	if (a4 == 2) {
		nox_xxx_consoleEsc_49B7A0();
	}
	return 1;
}

//----- (0046A820) --------------------------------------------------------
int sub_46A820(int a1, int a2, int a3, int a4) {
	if (a2 == 16415) {
		if (*(uint16_t*)(dword_5d4594_1064864 + 1052)) {
			nox_xxx_cmdSayDo_46A4B0(*(wchar2_t**)&dword_5d4594_1064864, *getMemIntPtr(0x5D4594, 1064872));
		}
		sub_46A6A0();
	}
	return 0;
}

//----- (0046A860) --------------------------------------------------------
int sub_46A860() {
	int result; // eax

	result = dword_5d4594_1064856;
	if (dword_5d4594_1064856) {
		result = nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1064856);
		dword_5d4594_1064856 = 0;
	}
	dword_5d4594_1064860 = 0;
	dword_5d4594_1064864 = 0;
	dword_5d4594_1064868 = 0;
	*getMemU32Ptr(0x5D4594, 1064872) = 0;
	return result;
}

//----- (004724E0) --------------------------------------------------------
void nox_client_mapZoomIn_4724E0() {
	nox_xxx_minimap_587000_149232 -= 10;
	if (*(int*)&nox_xxx_minimap_587000_149232 < 500) {
		nox_xxx_minimap_587000_149232 = 500;
	}
}

//----- (00472500) --------------------------------------------------------
void nox_client_mapZoomOut_472500() {
	nox_xxx_minimap_587000_149232 += 10;
	if (nox_xxx_minimap_587000_149232 > 4000) {
		nox_xxx_minimap_587000_149232 = 4000;
	}
}

//----- (00472520) --------------------------------------------------------
int nox_xxx_cliSetMinimapZoom_472520(int a1) {
	int result; // eax

	result = a1;
	nox_xxx_minimap_587000_149232 = a1;
	return result;
}

//----- (00472540) --------------------------------------------------------
int sub_472540(int a1) {
	int v1;     // edx
	int v2;     // eax
	int result; // eax
	int2 a1a;   // [esp+0h] [ebp-8h]

	if (a1 == *getMemU32Ptr(0x852978, 8)) {
		nox_xxx_getSomeCoods_435670(&a1a);
	} else {
		v1 = *(uint32_t*)(a1 + 16);
		a1a.field_0 = *(uint32_t*)(a1 + 12);
		a1a.field_4 = v1;
	}
	v2 = nox_xxx_polygonGetIdxA_421790(&a1a, *getMemIntPtr(0x5D4594, 1096312));
	if (v2) {
		*getMemU32Ptr(0x5D4594, 1096312) = v2;
	} else {
		v2 = *getMemU32Ptr(0x5D4594, 1096312);
	}
	if (v2) {
		result = (unsigned char)nox_xxx_polygonGetByIdx_4214A0(v2)[130];
	} else {
		result = 1;
	}
	return result;
}

//----- (004725C0) --------------------------------------------------------
void nox_xxx_drawMinimap4Sprite_4725C0(int a1) {
	int4* result; // eax

	result = (int4*)nox_client_drawable_testBuff_4356C0(*getMemIntPtr(0x852978, 8), 2);
	if (!result) {
		sub_437260();
		*getMemU32Ptr(0x5D4594, 1096316) = sub_472540(a1);
		nox_xxx_cliDrawMinimap_472600(a1, *getMemIntPtr(0x5D4594, 1096316));
		sub_437290();
	}
}

//----- (00472600) --------------------------------------------------------
void* sub_4106A0(int a1);
int sub_50CB00();
void* sub_50CB10();
int nox_xxx_cliDrawMinimap_472600(int a1, int a2) {
	char* v2;                           // ebp
	int v3;                             // esi
	int v4;                             // kr08_4
	int v5;                             // ebx
	int v6;                             // et1
	int v7;                             // esi
	int v8;                             // ebx
	int v9;                             // ebp
	int v10;                            // eax
	int v11;                            // edi
	char v12;                           // al
	int v13;                            // eax
	int v14;                            // esi
	int v15;                            // ebp
	unsigned char* v16;                 // esi
	int v17;                            // eax
	int v18;                            // eax
	int v19;                            // eax
	int v20;                            // et1
	int v21;                            // ecx
	int v22;                            // ebx
	int v23;                            // ebp
	int v24;                            // esi
	int v25;                            // ecx
	int v26;                            // et1
	char v27;                           // al
	char v28;                           // dl
	int v29;                            // edi
	char* v30;                          // esi
	float* v31;                         // esi
	int v32;                            // et1
	double v33;                         // st7
	int v34;                            // et1
	double v35;                         // st7
	float* j;                           // esi
	int v37;                            // et1
	double v38;                         // st7
	uint32_t* v39;                      // eax
	int k;                              // esi
	nox_player_polygon_check_data* v41; // eax
	int v42;                            // et1
	int v43;                            // eax
	int v44;                            // eax
	int v45;                            // edi
	uint32_t* v46;                      // eax
	char* v47;                          // eax
	int v48;                           // eax
	uint32_t* v49;                      // eax
	uint32_t* v50;                      // edi
	char* v51;                          // eax
	int v52;                           // eax
	int v53;                            // eax
	int v54;                            // eax
	char* v55;                          // eax
	char* v56;                          // eax
	uint32_t* v57;                      // eax
	int v58;                           // eax
	int l;                              // esi
	int v60;                            // eax
	int v61;                            // edx
	uint32_t* v62;                      // edi
	nox_player_polygon_check_data* v63; // eax
	int v64;                            // et1
	char* v65;                          // eax
	int v66;                           // eax
	int v68;                            // [esp-10h] [ebp-70h]
	int v69;                            // [esp+10h] [ebp-50h]
	int v70;                            // [esp+10h] [ebp-50h]
	int i;                              // [esp+14h] [ebp-4Ch]
	int v72;                            // [esp+14h] [ebp-4Ch]
	int v73;                            // [esp+14h] [ebp-4Ch]
	int v74;                            // [esp+18h] [ebp-48h]
	char* v75;                          // [esp+18h] [ebp-48h]
	int2 v76;                           // [esp+20h] [ebp-40h]
	int v77;                            // [esp+28h] [ebp-38h]
	int v78;                            // [esp+2Ch] [ebp-34h]
	int v79;                            // [esp+30h] [ebp-30h]
	int v80;                            // [esp+34h] [ebp-2Ch]
	int v81;                            // [esp+38h] [ebp-28h]
	int2 xLeft;                         // [esp+40h] [ebp-20h]
	int yTop;                           // [esp+4Ch] [ebp-14h]
	int2 v84;                           // [esp+50h] [ebp-10h]
	int v85;                            // [esp+5Ch] [ebp-4h]

	v2 = nox_draw_getViewport_437250();
	if (!getMemByte(0x5D4594, 1096300)) {
		*getMemU8Ptr(0x5D4594, 1096300) = nox_xxx_wallTileByName_410D60("InvisibleWallSet");
		*getMemU8Ptr(0x5D4594, 1096301) = nox_xxx_wallTileByName_410D60("InvisibleBlockingWallSet");
	}
	nox_client_drawEnableAlpha_434560(0);
	nox_xxx_wndDraw_49F7F0();
	v3 = nox_win_width / 6;
	v4 = nox_win_height - nox_win_width / 6;
	yTop = v4 / 2;
	nox_client_copyRect_49F6F0(0, 0, nox_win_width, nox_win_height);
	v5 = *(uint32_t*)v2;
	if (*(uint32_t*)v2 <= 0) {
		nox_client_drawRectFilledAlpha_49CF10(0, v4 / 2, v3, v3);
	} else {
		nox_client_drawSetColor_434460(nox_color_black_2650656);
		if (v5 >= v3) {
			nox_client_drawRectFilledOpaque_49CE30(0, v4 / 2, v3, v3);
		} else {
			nox_client_drawRectFilledOpaque_49CE30(0, v4 / 2, v5, v3);
			nox_client_drawRectFilledAlpha_49CF10(v5, v4 / 2, v3 - v5, v3);
		}
	}
	nox_client_drawEnableAlpha_434560(1);
	nox_client_drawSetColor_434460(nox_color_black_2650656);
	nox_client_drawSetAlpha_434580(0x5Au);
	nox_client_drawRectLines_473510(-1, yTop - 1, v3 + 2, v3 + 2);
	nox_client_drawSetAlpha_434580(0x3Cu);
	nox_client_drawRectLines_473510(-2, yTop - 2, v3 + 4, v3 + 4);
	nox_client_drawSetAlpha_434580(0x28u);
	nox_client_drawRectLines_473510(-3, yTop - 3, v3 + 6, v3 + 6);
	nox_client_drawEnableAlpha_434560(0);
	nox_client_copyRect_49F6F0(0, yTop, v3, v3);
	v6 = nox_xxx_minimap_587000_149232;
	v7 = v3 * v6 / 100;
	nox_xxx_getSomeCoods_435670(&v84);
	v8 = v84.field_0 - v7 / 2;
	v9 = v84.field_4 - v7 / 2;
	xLeft.field_0 = v84.field_0 - v7 / 2;
	xLeft.field_4 = v9;
	v81 = v8 / 23;
	v77 = (v8 + v7) / 23;
	v78 = (v9 + v7) / 23;
	v74 = 23 * (v8 / 23);
	v80 = 23 * (v9 / 23);
	v10 = v9 / 23;
	for (i = v9 / 23; i <= v78; ++i) {
		v11 = sub_4106A0(v10);
		while (v11) {
			v12 = *(uint8_t*)(v11 + 1);
			if (v12 == getMemByte(0x5D4594, 1096300)) {
				goto LABEL_37;
			}
			if (v12 == getMemByte(0x5D4594, 1096301)) {
				goto LABEL_37;
			}
			if (*(uint8_t*)(v11 + 8) && *(unsigned char*)(v11 + 8) != a2) {
				goto LABEL_37;
			}
			v13 = *(unsigned char*)(v11 + 5);
			if (v13 < v81) {
				goto LABEL_37;
			}
			if (v13 > v77) {
				break;
			}
			v14 = v74 + 23 * (v13 - v81);
			if (*(uint8_t*)(v11 + 4) & 0x10) {
				v15 = *(uint32_t*)(v11 + 32);
				if (!v15) {
					goto LABEL_37;
				}
				v69 = 0;
				v16 = getMemAt(0x587000, 149244);
				while (1) {
					v17 = nox_server_getWallAtGrid_410580(*((uint32_t*)v16 - 1) + *(unsigned char*)(v11 + 5),
														  *(uint32_t*)v16 + *(unsigned char*)(v11 + 6));
					if (v17) {
						if (*(uint32_t*)(v17 + 12)) {
							if (v69 < 4) {
								v20 = nox_xxx_minimap_587000_149232;
								v21 = 8 * *(unsigned char*)(v15 + 299);
								v22 = 100 * (*(int*)(v15 + 12) - xLeft.field_0) / v20;
								v85 = 100 * (*(int*)(v15 + 16) - xLeft.field_4) / v20;
								v23 = 100 * *getMemIntPtr(0x587000, 196184 + v21) / v20;
								v24 = 100 * *getMemIntPtr(0x587000, 196188 + v21) / v20;
								nox_client_drawSetColor_434460(*getMemIntPtr(0x85B3FC, 940));
								nox_client_drawAddPoint_49F500(v22, yTop + v85);
								nox_xxx_rasterPointRel_49F570(v23, v24);
								nox_client_drawLineFromPoints_49E4B0();
								v8 = xLeft.field_0;
							}
							goto LABEL_37;
						}
					} else {
						v18 = nox_server_getWallAtGrid_410580(*((uint32_t*)v16 - 1) + *(unsigned char*)(v11 + 5),
															  *(unsigned char*)(v11 + 6));
						if (v18 && *(uint32_t*)(v18 + 12) ||
							(v19 = nox_server_getWallAtGrid_410580(
								 *(unsigned char*)(v11 + 5), *(uint32_t*)v16 + *(unsigned char*)(v11 + 6))) != 0 &&
								*(uint32_t*)(v19 + 12)) {
							if (v69 < 4) {
								v20 = nox_xxx_minimap_587000_149232;
								v21 = 8 * *(unsigned char*)(v15 + 299);
								v22 = 100 * (*(int*)(v15 + 12) - xLeft.field_0) / v20;
								v85 = 100 * (*(int*)(v15 + 16) - xLeft.field_4) / v20;
								v23 = 100 * *getMemIntPtr(0x587000, 196184 + v21) / v20;
								v24 = 100 * *getMemIntPtr(0x587000, 196188 + v21) / v20;
								nox_client_drawSetColor_434460(*getMemIntPtr(0x85B3FC, 940));
								nox_client_drawAddPoint_49F500(v22, yTop + v85);
								nox_xxx_rasterPointRel_49F570(v23, v24);
								nox_client_drawLineFromPoints_49E4B0();
								v8 = xLeft.field_0;
							}
							goto LABEL_37;
						}
					}
					v16 += 8;
					++v69;
					if ((int)v16 >= (int)getMemAt(0x587000, 149276)) {
						goto LABEL_37;
					}
				}
			}
			if (nox_common_gameFlags_check_40A5C0(0x10000) || *(uint32_t*)(v11 + 12)) {
				v26 = nox_xxx_minimap_587000_149232;
				v25 = v26;
				v76.field_0 = 100 * (v14 - v8) / v26;
				v76.field_4 = yTop + 100 * (v80 - v9) / v26;
				v27 = *(uint8_t*)(v11 + 4);
				if (!(v27 & 4) || (v28 = *(uint8_t*)(*(uint32_t*)(v11 + 28) + 21), v28 != 3) && v28 != 2) {
					if (!(v27 & 0x20)) {
						sub_4730D0(&v76, *(uint8_t*)v11, 2300 / v25);
					}
				}
			}
		LABEL_37:
			v11 = *(uint32_t*)(v11 + 24);
			v9 = xLeft.field_4;
		}
		v10 = i + 1;
		v80 += 23;
	}
	if (nox_common_getEngineFlag(NOX_ENGINE_FLAG_ENABLE_SHOW_AI)) {
		v29 = sub_50CB00();
		v30 = (char*)sub_50CB10();
		if (v29 >= 2) {
			nox_client_drawSetColor_434460(dword_8531A0_2572);
			if (v29 - 1 > 0) {
				v31 = (float*)(v30 + 8);
				v72 = v29 - 1;
				do {
					v32 = nox_xxx_minimap_587000_149232;
					v33 = *(v31 - 1);
					xLeft.field_0 = (int)(100 * ((unsigned long long)(long long)*(v31 - 2) - v8)) / v32;
					nox_client_drawAddPoint_49F500(xLeft.field_0,
												   yTop + (int)(100 * ((unsigned long long)(long long)v33 - v9)) / v32);
					v34 = nox_xxx_minimap_587000_149232;
					v35 = v31[1];
					xLeft.field_0 = (int)(100 * ((unsigned long long)(long long)*v31 - v8)) / v34;
					nox_client_drawAddPoint_49F500(xLeft.field_0,
												   yTop + (int)(100 * ((unsigned long long)(long long)v35 - v9)) / v34);
					nox_client_drawLineFromPoints_49E4B0();
					v31 += 2;
					--v72;
				} while (v72);
			}
		}
		for (j = (float*)nox_xxx_minimapFirstMonster_50AAE0(); j; j = (float*)nox_xxx_minimapNextMonster_50AB10()) {
			nox_client_drawSetColor_434460(*getMemIntPtr(0x85B3FC, 940));
			v37 = nox_xxx_minimap_587000_149232;
			v38 = j[1];
			xLeft.field_0 = (int)(100 * ((unsigned long long)(long long)*j - v8)) / v37;
			nox_xxx_minimapDrawPoint_473570(xLeft.field_0,
											yTop + (int)(100 * ((unsigned long long)(long long)v38 - v9)) / v37);
		}
	}
	v73 = 0;
	if (!*getMemU32Ptr(0x5D4594, 1096304)) {
		*getMemU32Ptr(0x5D4594, 1096304) = nox_xxx_getTTByNameSpriteMB_44CFC0("Crown");
		*getMemU32Ptr(0x5D4594, 1096308) = nox_xxx_getTTByNameSpriteMB_44CFC0("GameBall");
	}
	v39 = nox_xxx_objGetTeamByNetCode_418C80(nox_player_netCode_85319C);
	v70 = (int)v39;
	if (v39 && nox_xxx_servObjectHasTeam_419130((int)v39)) {
		v73 = 1;
	}
	for (k = nox_xxx_cliFirstMinimapObj_459EB0(); k; k = nox_xxx_cliNextMinimapObj_459EC0(k)) {
		v41 = nox_xxx_polygonIsPlayerInPolygon_4217B0((int2*)(k + 12), 0);
		if (v41) {
			if (BYTE2(v41->field_0[32]) != a2) {
				continue;
			}
		} else if (a2 != 1) {
			continue;
		}
		v42 = nox_xxx_minimap_587000_149232;
		xLeft.field_0 = 100 * (*(int*)(k + 12) - v8) / v42;
		xLeft.field_4 = yTop + 100 * (*(int*)(k + 16) - v9) / v42;
		if (!(*(uint32_t*)(k + 112) & 0x400000) || (v43 = nox_color_blue_2650684, !(*(uint8_t*)(k + 116) & 8))) {
			v43 = *getMemU32Ptr(0x85B3FC, 940);
		}
		nox_client_drawSetColor_434460(v43);
		v44 = *(uint32_t*)(k + 108);
		if (v44 == *getMemIntPtr(0x5D4594, 1096304)) {
			if (nox_server_teamFirst_418B10() || (v45 = nox_xxx_cliGetSpritePlayer_45A000()) == 0) {
				// nop
			} else {
				while (!nox_client_drawable_testBuff_4356C0(v45, 30)) {
					v45 = sub_45A010(v45);
					if (!v45) {
						goto LABEL_64;
					}
				}
				continue;
			}
		LABEL_64:
			nox_client_drawSetColor_434460(dword_8531A0_2572);
			v46 = nox_xxx_objGetTeamByNetCode_418C80(*(uint32_t*)(k + 128));
			if (v46) {
				v47 = nox_xxx_getTeamByID_418AB0(*((unsigned char*)v46 + 4));
				if (v47) {
					v48 = nox_xxx_materialGetTeamColor_418D50((int)v47);
					nox_client_drawSetColor_434460(v48);
				}
			}
			sub_473420(&xLeft);
			continue;
		}
		if (v44 == *getMemIntPtr(0x5D4594, 1096308)) {
			v49 = nox_xxx_objGetTeamByNetCode_418C80(*(uint32_t*)(k + 128));
			v50 = v49;
			if (v49 && nox_xxx_servObjectHasTeam_419130((int)v49)) {
				v51 = nox_xxx_getTeamByID_418AB0(*((unsigned char*)v50 + 4));
				if (v51) {
					v52 = nox_xxx_materialGetTeamColor_418D50((int)v51);
					nox_client_drawSetColor_434460(v52);
				}
			} else {
				nox_client_drawSetColor_434460(nox_color_white_2523948);
			}
			nox_video_drawCircleRad3_4734F0(&xLeft.field_0);
			continue;
		}
		v53 = *(uint32_t*)(k + 112);
		if (v53 & 0x10000000) {
			if (*(uint32_t*)(k + 120) & 0x1000000) {
				nox_client_drawSetColor_434460(nox_color_white_2523948);
				v54 = sub_4B9470(*(const char***)(k + 436));
				v55 = nox_xxx_getTeamByID_418AB0(v54);
				if (v55) {
					v58 = nox_xxx_materialGetTeamColor_418D50((int)v55);
					nox_client_drawSetColor_434460(v58);
					sub_4733B0(&xLeft);
					continue;
				}
				sub_4733B0(&xLeft);
				continue;
			}
		} else {
			if (!(v53 & 4)) {
				nox_xxx_minimapDrawPoint_473570(xLeft.field_0, xLeft.field_4);
				continue;
			}
			if (!nox_common_gameFlags_check_40A5C0(32)) {
				if (k == *getMemIntPtr(0x852978, 8)) {
					sub_4735C0(xLeft.field_0, xLeft.field_4);
				} else {
					nox_xxx_minimapDrawPoint_473570(xLeft.field_0, xLeft.field_4);
				}
				continue;
			}
			v56 = nox_common_playerInfoGetByID_417040(*(uint32_t*)(k + 128));
			if (v56) {
				v81 = *((uint32_t*)v56 + 1) & 1;
				if (v81) {
					nox_client_drawSetColor_434460(nox_color_white_2523948);
					v57 = nox_xxx_objGetTeamByNetCode_418C80(*(uint32_t*)(k + 128));
					if (v57) {
						v55 = *((uint8_t*)v57 + 4) == 1 ? nox_xxx_getTeamByID_418AB0(2) : nox_xxx_getTeamByID_418AB0(1);
						if (v55) {
							v58 = nox_xxx_materialGetTeamColor_418D50((int)v55);
							nox_client_drawSetColor_434460(v58);
							sub_4733B0(&xLeft);
							continue;
						}
					}
					sub_4733B0(&xLeft);
					continue;
				}
			}
		}
	}
	v79 = dword_8531A0_2572;
	for (l = nox_xxx_cliGetSpritePlayer_45A000(); l; l = sub_45A010(l)) {
		v60 = nox_client_drawable_testBuff_4356C0(l, 30);
		v61 = *(uint32_t*)(l + 128);
		v77 = v60;
		v62 = nox_xxx_objGetTeamByNetCode_418C80(v61);
		v68 = *(uint32_t*)(l + 128);
		v81 = (int)v62;
		v75 = nox_common_playerInfoGetByID_417040(v68);
		if (v70 && v62 && v73) {
			v76.field_0 = nox_xxx_servCompareTeams_419150(v70, (int)v62);
			if (v76.field_0) {
				goto LABEL_103;
			}
		} else {
			v76.field_0 = 0;
		}
		if (!(v77 || l == *getMemU32Ptr(0x852978, 8) || *(uint8_t*)(dword_8531A0_2576 + 3680) & 1)) {
			continue;
		}
	LABEL_103:
		v63 = nox_xxx_polygonIsPlayerInPolygon_4217B0((int2*)(l + 12), 0);
		if ((!v63 || BYTE2(v63->field_0[32]) == a2) && v75 && (v75[3680] & 1) != 1) {
			v64 = nox_xxx_minimap_587000_149232;
			xLeft.field_0 = 100 * (*(int*)(l + 12) - v8) / v64;
			xLeft.field_4 = yTop + 100 * (*(int*)(l + 16) - v9) / v64;
			if (l == *getMemU32Ptr(0x852978, 8) || v76.field_0) {
				nox_client_drawSetColor_434460(v79);
			} else {
				nox_client_drawSetColor_434460(*getMemIntPtr(0x85B3FC, 940));
			}
			if (v77) {
				if (v81) {
					v65 = nox_xxx_getTeamByID_418AB0(*(unsigned char*)(v81 + 4));
					if (v65) {
						v66 = nox_xxx_materialGetTeamColor_418D50((int)v65);
						nox_client_drawSetColor_434460(v66);
					}
				}
				sub_473420(&xLeft);
			} else {
				nox_xxx_minimapDrawPoint_473570(xLeft.field_0, xLeft.field_4);
			}
		}
	}
	return sub_49F860();
}

//----- (004730D0) --------------------------------------------------------
int sub_4730D0(int2* a1, unsigned char a2, int a3) {
	int result = 0; // eax
	int v4;     // ebx
	int v5;     // edi
	int2* v6;   // ebp

	if (nox_xxx_minimap_587000_149232 <= 2000) {
		v4 = *getMemU32Ptr(0x85B3FC, 956);
		result = a2;
		v5 = a3 / 2;
		switch (a2) {
		case 0u:
			result =
				sub_473380(a1->field_0, a3 + a1->field_4, a1->field_0 + a3, a1->field_4, *getMemIntPtr(0x85B3FC, 956));
			break;
		case 1u:
			result =
				sub_473380(a1->field_0, a1->field_4, a1->field_0 + a3, a1->field_4 + a3, *getMemIntPtr(0x85B3FC, 956));
			break;
		case 2u:
			sub_473380(a1->field_0, a3 + a1->field_4, a1->field_0 + a3, a1->field_4, *getMemIntPtr(0x85B3FC, 956));
			result = sub_473380(a1->field_0, a1->field_4, a1->field_0 + a3, a1->field_4 + a3, v4);
			break;
		case 3u:
			sub_473380(a1->field_0, a1->field_4, a1->field_0 + v5, a1->field_4 + v5, *getMemIntPtr(0x85B3FC, 956));
			result = sub_473380(a1->field_0, a3 + a1->field_4, a3 + a1->field_0, a1->field_4, v4);
			break;
		case 4u:
			sub_473380(a1->field_0, a1->field_4, a1->field_0 + a3, a1->field_4 + a3, *getMemIntPtr(0x85B3FC, 956));
			result = sub_473380(v5 + a1->field_0, v5 + a1->field_4, a3 + a1->field_0, a1->field_4, v4);
			break;
		case 5u:
			sub_473380(a1->field_0, a3 + a1->field_4, a1->field_0 + a3, a1->field_4, *getMemIntPtr(0x85B3FC, 956));
			result = sub_473380(v5 + a1->field_0, v5 + a1->field_4, a3 + a1->field_0, a1->field_4 + a3, v4);
			break;
		case 6u:
			v6 = a1;
			sub_473380(a1->field_0, a1->field_4, a1->field_0 + a3, a1->field_4 + a3, *getMemIntPtr(0x85B3FC, 956));
			result = sub_473380(v6->field_0, a3 + v6->field_4, v5 + v6->field_0, v6->field_4 + v5, v4);
			break;
		case 7u:
			sub_473380(a1->field_0, a1->field_4, a1->field_0 + v5, a1->field_4 + v5, *getMemIntPtr(0x85B3FC, 956));
			result = sub_473380(v5 + a1->field_0, v5 + a1->field_4, a3 + a1->field_0, a1->field_4, v4);
			break;
		case 8u:
			sub_473380(v5 + a1->field_0, v5 + a1->field_4, a1->field_0 + a3, a1->field_4 + a3,
					   *getMemIntPtr(0x85B3FC, 956));
			result = sub_473380(v5 + a1->field_0, v5 + a1->field_4, a3 + a1->field_0, a1->field_4, v4);
			break;
		case 9u:
			sub_473380(v5 + a1->field_0, v5 + a1->field_4, a1->field_0 + a3, a1->field_4 + a3,
					   *getMemIntPtr(0x85B3FC, 956));
			result = sub_473380(v5 + a1->field_0, v5 + a1->field_4, a1->field_0, a1->field_4 + a3, v4);
			break;
		case 0xAu:
			v6 = a1;
			sub_473380(a1->field_0, a1->field_4, a1->field_0 + v5, a1->field_4 + v5, *getMemIntPtr(0x85B3FC, 956));
			result = sub_473380(v6->field_0, a3 + v6->field_4, v5 + v6->field_0, v6->field_4 + v5, v4);
			break;
		default:
			return result;
		}
	} else {
		nox_client_drawSetColor_434460(*getMemIntPtr(0x85B3FC, 956));
		nox_client_drawPixel_49EFA0(a1->field_0, a1->field_4);
	}
	return result;
}

//----- (00473380) --------------------------------------------------------
int sub_473380(int a1, int a2, int a3, int a4, int a5) {
	nox_client_drawSetColor_434460(a5);
	nox_client_drawAddPoint_49F500(a1, a2);
	nox_client_drawAddPoint_49F500(a3, a4);
	return nox_client_drawLineFromPoints_49E4B0();
}

//----- (004733B0) --------------------------------------------------------
int sub_4733B0(uint32_t* a1) {
	int v1; // esi
	int v2; // edi

	v1 = a1[1] + 4;
	v2 = *a1 - 2;
	nox_client_drawAddPoint_49F500(v2, v1);
	v1 -= 8;
	nox_client_drawAddPoint_49F500(v2, v1);
	nox_client_drawLineFromPoints_49E4B0();
	nox_client_drawAddPoint_49F500(v2, v1);
	v2 += 4;
	nox_client_drawAddPoint_49F500(v2, v1);
	nox_client_drawLineFromPoints_49E4B0();
	nox_client_drawAddPoint_49F500(v2, v1);
	v1 += 4;
	nox_client_drawAddPoint_49F500(v2, v1);
	nox_client_drawLineFromPoints_49E4B0();
	nox_client_drawAddPoint_49F500(v2, v1);
	nox_client_drawAddPoint_49F500(v2 - 4, v1);
	return nox_client_drawLineFromPoints_49E4B0();
}

//----- (00473420) --------------------------------------------------------
int sub_473420(uint32_t* a1) {
	int v1; // edi
	int v2; // esi

	v1 = a1[1] + 6;
	v2 = *a1 - 4;
	nox_client_drawAddPoint_49F500(v2, v1);
	v1 -= 12;
	v2 -= 2;
	nox_client_drawAddPoint_49F500(v2, v1);
	nox_client_drawLineFromPoints_49E4B0();
	nox_client_drawAddPoint_49F500(v2, v1);
	v1 += 6;
	v2 += 4;
	nox_client_drawAddPoint_49F500(v2, v1);
	nox_client_drawLineFromPoints_49E4B0();
	nox_client_drawAddPoint_49F500(v2, v1);
	v1 -= 6;
	v2 += 2;
	nox_client_drawAddPoint_49F500(v2, v1);
	nox_client_drawLineFromPoints_49E4B0();
	nox_client_drawAddPoint_49F500(v2, v1);
	v1 += 6;
	v2 += 2;
	nox_client_drawAddPoint_49F500(v2, v1);
	nox_client_drawLineFromPoints_49E4B0();
	nox_client_drawAddPoint_49F500(v2, v1);
	v1 -= 6;
	v2 += 4;
	nox_client_drawAddPoint_49F500(v2, v1);
	nox_client_drawLineFromPoints_49E4B0();
	nox_client_drawAddPoint_49F500(v2, v1);
	v1 += 12;
	v2 -= 2;
	nox_client_drawAddPoint_49F500(v2, v1);
	nox_client_drawLineFromPoints_49E4B0();
	nox_client_drawAddPoint_49F500(v2, v1);
	nox_client_drawAddPoint_49F500(v2 - 8, v1);
	return nox_client_drawLineFromPoints_49E4B0();
}

//----- (004734F0) --------------------------------------------------------
void nox_video_drawCircle_4B0B90(int a1, int a2, int a3);
void nox_video_drawCircleRad3_4734F0(int* a1) { nox_video_drawCircle_4B0B90(a1[0], a1[1], 3); }

//----- (00473510) --------------------------------------------------------
int nox_client_drawRectLines_473510(int a1, int a2, int a3, int a4) {
	int v4; // esi
	int v5; // edi

	nox_client_drawAddPoint_49F500(a1, a2);
	v4 = a1 + a3 - 1;
	nox_client_drawAddPoint_49F500(v4, a2);
	nox_client_drawLineFromPoints_49E4B0();
	nox_client_drawAddPoint_49F500(v4, a2);
	v5 = a4 - 1 + a2;
	nox_client_drawAddPoint_49F500(v4, v5);
	nox_client_drawLineFromPoints_49E4B0();
	nox_client_drawAddPoint_49F500(v4, v5);
	nox_client_drawAddPoint_49F500(a1, v5);
	return nox_client_drawLineFromPoints_49E4B0();
}

//----- (00473570) --------------------------------------------------------
void nox_xxx_minimapDrawPoint_473570(int xLeft, int yTop) {
	if (nox_xxx_minimap_587000_149232 > 1200) {
		nox_xxx_drawPointMB_499B70(xLeft, yTop, (nox_xxx_minimap_587000_149232 < 1750) + 2);
	} else {
		nox_xxx_drawPointMB_499B70(xLeft, yTop, 4);
	}
}

//----- (004735C0) --------------------------------------------------------
void sub_4735C0(int xLeft, int yTop) {
	if (nox_xxx_minimap_587000_149232 > 1200) {
		nox_xxx_drawPointMB_499B70(xLeft, yTop, (nox_xxx_minimap_587000_149232 < 1750) + 4);
	} else {
		nox_xxx_drawPointMB_499B70(xLeft, yTop, 6);
	}
}

//----- (004738E0) --------------------------------------------------------
int nox_xxx_drawMinimapAndLines_4738E0() {
	int result;   // eax
	uint32_t* v1; // eax

	result = 1;
	if (nox_client_gui_flag_1556112 != 1) {
		if (getMemByte(0x5D4594, 1096424) & 1) {
			v1 = nox_xxx_netSpriteByCodeDynamic_45A6F0(nox_player_netCode_85319C);
			nox_xxx_drawMinimap4Sprite_4725C0((int)v1);
		}
		result = nox_xxx_drawMessageLines_445530();
	}
	return result;
}

//----- (00473920) --------------------------------------------------------
void nox_xxx____setargv_11_473920() { *getMemU32Ptr(0x5D4594, 1096520) = 1; }

//----- (00473930) --------------------------------------------------------
char* sub_473930() {
	char* result; // eax

	*getMemU32Ptr(0x5D4594, 1096456) = nox_xxx_gLoadAnim_42FA20("ConfusedBirdies");
	result = nox_xxx_gLoadAnim_42FA20("SphericalShieldAnim");
	*getMemU32Ptr(0x5D4594, 1096460) = result;
	return result;
}

//----- (00473960) --------------------------------------------------------
int sub_473960() {
	int result; // eax

	result = 0;
	*getMemU32Ptr(0x5D4594, 1096456) = 0;
	*getMemU32Ptr(0x5D4594, 1096460) = 0;
	return result;
}

//----- (004739E0) --------------------------------------------------------
int sub_4739E0(uint32_t* a1, int2* a2, int2* a3) {
	int result; // eax

	a3->field_0 = a2->field_0 + *a1 - a1[4];
	result = a2->field_4;
	a3->field_4 = result + a1[1] - a1[5];
	return result;
}

//----- (00473A10) --------------------------------------------------------
int sub_473A10(uint32_t* a1, int2* a2, uint32_t* a3) {
	int result; // eax

	*a3 = a2->field_0 + a1[4] - *a1;
	result = a2->field_4;
	a3[1] = result + a1[5] - a1[1];
	return result;
}

//----- (00473C10) --------------------------------------------------------
uint32_t nox_xxx_wallFlags(int i);
void nox_xxx_drawWalls_473C10(nox_draw_viewport_t* vp, void* data) {
	uint32_t* a1 = vp;
	unsigned char* a2 = data;
	unsigned char* v3; // esi
	unsigned char v4;  // dl
	int v5;            // ecx
	int v6;            // ebx
	int v7;            // ebp
	int v8;            // eax
	int v9;            // edi
	int v10;           // eax
	int v11;           // ecx
	int v12;           // eax
	int v13;           // ecx
	int v14;           // eax
	int v15;           // edx
	int v16;           // eax
	int v17;           // ecx
	int v18;           // eax
	int v19;           // ecx
	bool v20;          // zf
	int v21;           // edx
	int v22;           // ebx
	int v23;           // ecx
	int v24;           // edx
	int v25;           // eax
	unsigned char v26; // al
	unsigned char v27; // al
	char v28;          // cl
	int v29;           // eax
	int v30;           // eax
	int v31;           // eax
	int* v32;          // edi
	int v33;           // eax
	int v34;           // eax
	int v35;           // eax
	int v36;           // edx
	int v37;           // eax
	int v38;           // eax
	int v39;           // eax
	int v40;           // edx
	int v41;           // eax
	int v42;           // eax
	int v43;           // eax
	int v44;           // edx
	int v45x;         // eax
	int v45y;         // eax
	int v46;           // ebx
	int v47;           // ebp
	int v48;           // eax
	int v49;           // ecx
	int v50;           // edx
	int v51;           // eax
	int v52;           // eax
	int v53;           // eax
	uint8_t* v54;      // edi
	int v55x;         // eax
	int v55y;         // eax
	int v56;           // ebx
	int v57;           // ebp
	int v58;           // eax
	int v59;           // edx
	int v60;           // ecx
	int v61;           // eax
	int v63;           // [esp-18h] [ebp-80h]
	int v64;           // [esp-14h] [ebp-7Ch]
	int v65;           // [esp-10h] [ebp-78h]
	int v66;           // [esp-Ch] [ebp-74h]
	int v67;           // [esp-8h] [ebp-70h]
	int v68;           // [esp-4h] [ebp-6Ch]
	int v69;           // [esp-4h] [ebp-6Ch]
	int a3;            // [esp+10h] [ebp-58h]
	int a4;            // [esp+14h] [ebp-54h]
	int v72;           // [esp+18h] [ebp-50h]
	int v73;           // [esp+1Ch] [ebp-4Ch]
	int v74;           // [esp+20h] [ebp-48h]
	int v75;           // [esp+24h] [ebp-44h]
	int v76;           // [esp+28h] [ebp-40h]
	int2 v77;          // [esp+2Ch] [ebp-3Ch]
	int2 a2a;          // [esp+34h] [ebp-34h]
	int2 a1a;          // [esp+3Ch] [ebp-2Ch]
	int2 v80;          // [esp+44h] [ebp-24h]
	int2 v81;          // [esp+4Ch] [ebp-1Ch]
	int v82;           // [esp+54h] [ebp-14h]
	int v83[3];        // [esp+5Ch] [ebp-Ch]
	int v84;           // [esp+70h] [ebp+8h]

	v3 = a2;
	a4 = nox_win_width;
	v72 = 0;
	a3 = 0;
	if (!a2) {
		return;
	}
	v4 = a2[4];
	if (!(v4 & 1)) {
		return;
	}
	v5 = a2[6];
	v6 = *a1 + 23 * a2[5] - a1[4];
	v82 = *a1 + 23 * a2[5] - a1[4];
	v7 = a1[1] + 23 * v5 - a1[5];
	v74 = *getMemU32Ptr(0x587000, 149364 + 4 * a2[3]);
	v8 = v74;
	if (v74 == -1) {
		v8 = *a2;
		v74 = *a2;
	}
	v84 = v8;
	if (v8) {
		if (v8 == 1 && v4 & 0x40) {
			v84 = 12;
		}
	} else if (v4 & 0x40) {
		v84 = 11;
	}
	if (*getMemU32Ptr(0x587000, 80808)) {
		v9 = 16 * v74;
		v10 = *getMemU32Ptr(0x587000, 85440 + 16 * v74);
		v11 = *getMemU32Ptr(0x587000, 85448 + 16 * v74);
		a1a.field_4 = v7 + *getMemU32Ptr(0x587000, 85444 + 16 * v74);
		v12 = v6 + v10;
		a2a.field_4 = v7 + *getMemU32Ptr(0x587000, 85452 + 16 * v74);
		v13 = v6 + v11;
		a1a.field_0 = v12;
		a2a.field_0 = v13;
		if (v74 == 7 || v74 == 9) {
			if (sub_4C42A0(&a1a, &a2a, &a3, &a4)) {
				v22 = 1;
			} else {
				v22 = 0;
				a3 = a2a.field_0;
			}
			v23 = *getMemU32Ptr(0x587000, 85508 + v9);
			a1a.field_0 = v82 + *getMemU32Ptr(0x587000, 85504 + v9);
			v24 = v82 + *getMemU32Ptr(0x587000, 85512 + v9);
			v25 = *getMemU32Ptr(0x587000, 85516 + v9);
			a1a.field_4 = v7 + v23;
			a2a.field_0 = v24;
			a2a.field_4 = v7 + v25;
			if (sub_4C42A0(&a1a, &a2a, &a3, &a4)) {
				v19 = a3;
			} else {
				if (!v22) {
					return;
				}
				if (a4 > a1a.field_0) {
					a4 = a1a.field_0;
				}
				v19 = a3;
			}
		} else {
			if (v74 != 8 && v74 != 10) {
				if (!sub_4C42A0(&a1a, &a2a, &a3, &a4)) {
					return;
				}
				v19 = a3;
			} else {
				v76 = v13;
				v75 = v12;
				if (sub_4C42A0(&a1a, &a2a, &v75, &v76)) {
					v73 = v76 - v75 >= 3;
				} else {
					v73 = 0;
				}
				v14 = *getMemU32Ptr(0x587000, 85504 + 16 * v74);
				v15 = *getMemU32Ptr(0x587000, 85516 + 16 * v74);
				v80.field_4 = v7 + *getMemU32Ptr(0x587000, 85508 + 16 * v74);
				v16 = v6 + v14;
				v17 = v6 + *getMemU32Ptr(0x587000, 85512 + 16 * v74);
				v80.field_0 = v16;
				a3 = v16;
				v81.field_0 = v17;
				a4 = v17;
				v81.field_4 = v7 + v15;
				v18 = sub_4C42A0(&v80, &v81, &a3, &a4);
				v19 = a3;
				v20 = v18 == 0;
				if (v20) {
					v21 = 0;
				} else {
					v21 = a4 - a3 >= 3;
				}
				if (v73) {
					if (v21) {
						if (a3 > v75) {
							v19 = v75;
							a3 = v75;
						}
						if (v19 <= v80.field_0) {
							v19 = 0;
							a3 = 0;
						}
						if (a4 < v76) {
							a4 = v76;
						}
						if (a4 >= v81.field_0) {
							a4 = nox_win_width;
						}
					} else {
						v19 = v75;
						a3 = v75;
						a4 = v76;
						if (v74 != 8) {
							v84 = 1;
							if (v19 == v80.field_0) {
								v19 = 0;
								a3 = 0;
							}
						} else {
							v84 = 0;
							if (v76 == v81.field_0) {
								a4 = nox_win_width;
							}
						}
					}
				} else {
					if (!v21) {
						return;
					}
					v84 = (v74 != 8) + 13;
					if (a4 == v81.field_0) {
						a4 = nox_win_width;
					}
					if (v19 == v80.field_0) {
						v19 = 0;
						a3 = 0;
					}
				}
			}
		}
		if (v19 >= a4) {
			v26 = v3[4];
			v3[3] = 0;
			v3[4] = v26 & 0xFC;
			return;
		}
	}
	v27 = v3[4];
	v28 = v3[4] & 2;
	if (!v28) {
		v29 = (v3[4] >> 2) & 2;
		goto LABEL_64;
	}
	if (*getMemU32Ptr(0x5D4594, 805848) && nox_client_translucentFrontWalls_805844) {
		if (!nox_client_highResFrontWalls_80820 && nox_client_highResFloors_154952) {
			v72 |= 4u;
			goto LABEL_61;
		}
		v72 = 8;
	}
	if (!nox_client_highResFrontWalls_80820) {
		v72 |= 4u;
	}
LABEL_61:
	v29 = (v27 & 8 | 4u) >> 2;
LABEL_64:
	v73 = v29;
	if (v28 && nox_client_translucentFrontWalls_805844 && !(nox_xxx_wallFlags(v3[1]) & 4)) {
		v30 = v72;
		LOBYTE(v30) = v72 | 2;
		v72 = v30;
	} else {
		v72 |= 1u;
	}
	if (*getMemU32Ptr(0x587000, 80816)) {
		switch (v74) {
		case 0:
		case 3:
			v31 = v3[6];
			v77.field_0 = 23 * v3[5];
			v77.field_4 = 23 * (v31 + 1);
			v32 = sub_469920(&v77);
			if (v32 != (int*)31) {
				v83[0] = *v32;
				v83[1] = v32[1];
				v33 = v32[2];
				v32 = v83;
				v83[2] = v33;
			}
			v77.field_0 += 23;
			v77.field_4 -= 23;
			v34 = sub_469920(&v77);
			break;
		case 1:
		case 4:
			v35 = v3[6];
			v77.field_0 = 23 * v3[5];
			v77.field_4 = 23 * v35;
			v32 = sub_469920(&v77);
			if (v32 != (int*)31) {
				v83[0] = *v32;
				v83[1] = v32[1];
				v36 = v32[2];
				v32 = v83;
				v83[2] = v36;
			}
			v77.field_0 += 23;
			v77.field_4 += 23;
			v34 = sub_469920(&v77);
			break;
		case 7:
			v37 = v3[6];
			v77.field_0 = 23 * v3[5];
			v77.field_4 = 23 * v37;
			v32 = sub_469920(&v77);
			if (v32 != (int*)31) {
				v83[0] = *v32;
				v83[1] = v32[1];
				v38 = v32[2];
				v32 = v83;
				v83[2] = v38;
			}
			v77.field_0 += 23;
			v34 = sub_469920(&v77);
			break;
		case 8:
			v39 = v3[6];
			v77.field_0 = 23 * v3[5] + 11;
			v77.field_4 = 23 * v39 + 11;
			v32 = sub_469920(&v77);
			if (v32 != (int*)31) {
				v83[0] = *v32;
				v83[1] = v32[1];
				v40 = v32[2];
				v32 = v83;
				v83[2] = v40;
			}
			v77.field_0 -= 34;
			v77.field_4 -= 34;
			v34 = sub_469920(&v77);
			break;
		case 10:
			v41 = v3[6];
			v77.field_0 = 23 * v3[5];
			v77.field_4 = 23 * (v41 + 1);
			v32 = sub_469920(&v77);
			if (v32 != (int*)31) {
				v83[0] = *v32;
				v83[1] = v32[1];
				v42 = v32[2];
				v32 = v83;
				v83[2] = v42;
			}
			v77.field_0 += 11;
			v77.field_4 -= 11;
			v34 = sub_469920(&v77);
			break;
		default:
			v43 = v3[6];
			v77.field_0 = 23 * v3[5];
			v77.field_4 = 23 * (v43 + 1);
			v32 = sub_469920(&v77);
			if (v32 != (int*)31) {
				v83[0] = *v32;
				v83[1] = v32[1];
				v44 = v32[2];
				v32 = v83;
				v83[2] = v44;
			}
			v77.field_0 += 23;
			v34 = sub_469920(&v77);
			break;
		}
		v74 = v34;
		nox_xxx_getWallDrawOffset_46A3F0(v3[1], v84, v3[2], v73, &v45x, &v45y);
		v46 = v82 + v45x - 51;
		v47 = -73 - v45y + v7;
		sub_4345F0(1);
		v48 = *((uint8_t*)v32 + 8);
		LOBYTE(v49) = *((uint8_t*)v32 + 4);
		LOBYTE(v50) = *(uint8_t*)v32;
		nox_draw_setColorMultAndIntensityRGB_433CD0(v50, v49, v48);
		if (!(v72 & 2)) {
			v69 = v72;
			v66 = a4;
			v65 = a3;
			v64 = nox_win_height;
			v63 = v74;
			v52 = nox_xxx_getWallSprite_46A3B0(v3[1], v84, v3[2], v73);
			nox_xxx_edgeDraw_480EF0(v52, v46, v47, v32, v63, v64, v65, v66, 0, v69);
			goto LABEL_106;
		}
		if (!sub_47D380(a3, a4)) {
			goto LABEL_106;
		}
		nox_client_drawEnableAlpha_434560(1);
		nox_client_drawSetAlpha_434580(0x80u);
		sub_47D400(nox_client_highResFrontWalls_80820 == 0, a1[5]);
		v68 = v47;
		v67 = v46;
		v51 = nox_xxx_getWallSprite_46A3B0(v3[1], v84, v3[2], v73);
	} else {
		v53 = v3[6];
		v77.field_0 = 23 * v3[5] + 11;
		v77.field_4 = 23 * v53 + 11;
		v54 = sub_469920(&v77);
		nox_xxx_getWallDrawOffset_46A3F0(v3[1], v84, v3[2], v73, &v55x, &v55y);
		v56 = v82 + v55x - 50;
		v57 = -72 - v55y + v7;
		sub_4345F0(1);
		LOBYTE(v59) = v54[8];
		v58 = v54[4];
		LOBYTE(v60) = *v54;
		nox_draw_setColorMultAndIntensityRGB_433CD0(v60, v58, v59);
		if (!(v72 & 2)) {
			if (sub_47D380(a3, a4)) {
				sub_47D400(nox_client_highResFrontWalls_80820 == 0, a1[5]);
				v61 = nox_xxx_getWallSprite_46A3B0(v3[1], v84, v3[2], v73);
				nox_client_drawImageAt_47D2C0(v61, v56, v57);
				sub_47D400(0, 0);
			}
			goto LABEL_106;
		}
		if (!sub_47D380(a3, a4)) {
			goto LABEL_106;
		}
		nox_client_drawEnableAlpha_434560(1);
		nox_client_drawSetAlpha_434580(0x80u);
		sub_47D400(nox_client_highResFrontWalls_80820 == 0, a1[5]);
		v68 = v57;
		v67 = v56;
		v51 = nox_xxx_getWallSprite_46A3B0(v3[1], v84, v3[2], v73);
	}
	nox_client_drawImageAt_47D2C0(v51, v67, v68);
	sub_47D400(0, 0);
	nox_client_drawEnableAlpha_434560(0);
LABEL_106:
	sub_4345F0(0);
	v3[3] = 0;
	v3[4] &= 0xFC;
	*((uint32_t*)v3 + 3) = 1;
	return;
}
// 474366: variable 'v50' is possibly undefined
// 474366: variable 'v49' is possibly undefined
// 4744A3: variable 'v60' is possibly undefined
// 4744A3: variable 'v59' is possibly undefined

//----- (00474B40) --------------------------------------------------------
int sub_474B40(nox_drawable* dr) {
	int a1 = dr;
	uint32_t* v1; // edi
	uint32_t* v2; // eax
	int v3;       // eax

	v1 = nox_xxx_objGetTeamByNetCode_418C80(nox_player_netCode_85319C);
	if (v1) {
		v2 = nox_xxx_objGetTeamByNetCode_418C80(*(uint32_t*)(a1 + 128));
		if (v2) {
			if (nox_player_netCode_85319C == *(uint32_t*)(a1 + 128) ||
				nox_xxx_servCompareTeams_419150((int)v1, (int)v2)) {
				return 1;
			}
		}
	}
	v3 = *getMemU32Ptr(0x852978, 8);
	if (a1 == *getMemU32Ptr(0x852978, 8)) {
		return 1;
	}
	if (*getMemU32Ptr(0x852978, 8)) {
		if (!nox_client_drawable_testBuff_4356C0(*getMemIntPtr(0x852978, 8), 21)) {
			v3 = *getMemU32Ptr(0x852978, 8);
			goto LABEL_9;
		}
		return 1;
	}
LABEL_9:
	if (*(uint8_t*)(a1 + 112) & 4) {
		if (a1 != v3) {
			nox_common_playerInfoGetByID_417040(*(uint32_t*)(a1 + 128));
		}
	}
	return 0;
}

//----- (004756E0) --------------------------------------------------------
int nox_xxx_sprite_4756E0_drawable(nox_drawable* dr) {
	uint32_t* a1 = dr;
	int result;           // eax
	int (*v2)(int*, int); // esi
	int v3;               // edx
	int v4;               // ecx

	result = 0;
	v2 = (int (*)(int*, int))a1[75];
	if (v2) {
		v3 = a1[30];
		v4 = a1[28];
		if (!(v3 & 0x1000) && v3 & 1 && (v2 == nox_thing_static_draw || v2 == nox_thing_static_random_draw) &&
			!(v4 & 0x80800000) && (v3 & 0x48 || v4 & 0x400000) && !(v3 & 0x800)) {
			result = 1;
		}
	}
	return result;
}

//----- (00475740) --------------------------------------------------------
int nox_xxx_sprite_475740_drawable(nox_drawable* dr) {
	uint32_t* a1 = dr;
	int result;           // eax
	int (*v2)(int*, int); // edx
	int v3;               // ebx
	int v4;               // ecx

	result = 0;
	v2 = (int (*)(int*, int))a1[75];
	if (v2) {
		v3 = a1[30];
		v4 = a1[28];
		if (!(v3 & 0x1000)) {
			if (v3 & 1) {
				result = 1;
				if ((v2 == nox_thing_static_draw || v2 == nox_thing_static_random_draw) && !(v4 & 0x80800000) &&
					!(v3 & 0x800) && (v3 & 0x48 || v4 & 0x400000)) {
					result = 0;
				}
			}
		}
	}
	return result;
}

//----- (004757A0) --------------------------------------------------------
int nox_xxx_sprite_4757A0_drawable(nox_drawable* dr) {
	int a1 = dr;
	int result; // eax
	int v2;     // ecx

	result = 0;
	if (*(uint32_t*)(a1 + 300)) {
		v2 = *(uint32_t*)(a1 + 120);
		if (!(v2 & 0x1000)) {
			if (v2 & 0x4000) {
				if (v2 & 0x40) {
					result = 1;
				}
			}
		}
	}
	return result;
}

//----- (004757D0) --------------------------------------------------------
int sub_4757D0_drawable(nox_drawable* dr) {
	uint32_t* a1 = dr;
	int result; // eax
	int v2;     // ecx

	result = 0;
	if (a1[75]) {
		v2 = a1[30];
		if (!(v2 & 1) && (!(a1[28] & 0x2000) || v2 & 0x1000000) && !(v2 & 0x1000)) {
			result = 1;
		}
	}
	return result;
}
