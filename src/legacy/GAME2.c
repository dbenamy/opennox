#include <math.h>
#include <stdio.h>
#include <time.h>

#include "client/audio/ail/compat_mss.h"

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
#include "GAME4_1.h"
#include "GAME5_2.h"
#include "client__drawable__drawable.h"
#include "common__system__settings.h"
#include "common__system__team.h"

#include "client__gui__gadgets__listbox.h"
#include "client__gui__gamewin__gamewin.h"
#include "client__gui__gui_ctf.h"
#include "client__gui__guibook.h"
#include "client__gui__guibrief.h"
#include "client__gui__guicon.h"
#include "client__gui__guisave.h"
#include "client__gui__guispell.h"
#include "client__gui__servopts__access.h"
#include "client__gui__servopts__general.h"
#include "client__gui__servopts__guiserv.h"
#include "client__gui__servopts__playrlst.h"
#include "client__gui__window.h"
#include "client__shell__mainmenu.h"
#include "common__gamemech__pausefx.h"
#include "common__net_list.h"
#include "common__strman.h"
#include "common__crypt.h"
#include "operators.h"

#include "client__audio__audevent.h"
#include "client__video__draw_common.h"

#include "client__system__ctrlevnt.h"
#include "common__magic__speltree.h"
#include "input_common.h"

extern uint32_t dword_8531A0_2576;
extern uint32_t dword_5d4594_1046852;
extern uint32_t dword_5d4594_1047532;
extern uint32_t dword_5d4594_1049692;
extern uint32_t dword_5d4594_1046652;
extern uint32_t dword_5d4594_1047536;
extern uint32_t dword_5d4594_1045420;
extern uint32_t dword_5d4594_831240;
extern uint32_t dword_5d4594_831260;
extern uint32_t dword_5d4594_832480;
extern uint32_t dword_587000_126996;
extern uint32_t dword_5d4594_1045428;
extern uint32_t dword_5d4594_831244;
extern uint32_t dword_5d4594_831076;
extern uint32_t dword_5d4594_1049524;
extern uint32_t dword_5d4594_1047524;
extern uint32_t dword_5d4594_1046864;
extern uint32_t dword_5d4594_831256;
extern uint32_t dword_5d4594_1049536;
extern uint32_t dword_5d4594_1046944;
extern uint32_t dword_5d4594_1049696;
extern uint32_t dword_5d4594_831276;
extern uint32_t dword_5d4594_1046648;
extern uint32_t dword_5d4594_1046640;
extern uint32_t dword_5d4594_831084;
extern uint32_t dword_5d4594_1046956;
extern uint32_t dword_5d4594_1047936;
extern uint32_t dword_5d4594_831220;
extern uint32_t dword_587000_122848;
extern uint32_t dword_5d4594_1046932;
extern uint32_t dword_5d4594_1046948;
extern uint32_t dword_5d4594_830864;
extern uint32_t dword_5d4594_1047528;
extern uint32_t dword_5d4594_1046636;
extern uint32_t dword_5d4594_832520;
extern uint32_t dword_5d4594_832500;
extern uint32_t dword_5d4594_832528;
extern uint32_t dword_5d4594_832524;
extern uint32_t dword_5d4594_832512;
extern uint32_t dword_5d4594_832496;
extern uint32_t nox_xxx_aNox_cfg_0_587000_132136;
extern uint32_t dword_5d4594_832516;
extern uint32_t dword_5d4594_1049532;
extern uint32_t dword_5d4594_832508;
extern uint32_t dword_5d4594_832504;
extern uint32_t dword_5d4594_831224;
extern uint32_t dword_5d4594_1046928;
extern uint32_t dword_5d4594_832492;
extern uint32_t dword_5d4594_1047512;
extern uint32_t dword_5d4594_832532;
extern uint32_t dword_5d4594_1046656;
extern uint32_t dword_5d4594_1049484;
extern uint32_t dword_5d4594_832536;
extern uint32_t dword_5d4594_1046952;
extern void* dword_587000_81128;
extern void* dword_587000_122852;
extern uint32_t dword_5d4594_1049496;
extern uint32_t dword_5d4594_1047932;
extern uint32_t dword_5d4594_1049512;
extern uint32_t nox_wnd_briefing_831232;
extern uint32_t dword_5d4594_1045432;
extern uint32_t dword_5d4594_1046924;
extern void* dword_587000_127004;
extern uint32_t dword_5d4594_1047520;
extern uint32_t nox_xxx_aNox_cfg_0_587000_132132;
extern uint32_t dword_5d4594_1049520;
extern uint32_t dword_5d4594_1046936;
extern uint32_t dword_5d4594_1047540;
extern uint32_t dword_5d4594_831236;
extern uint32_t dword_5d4594_1046872;
extern uint32_t nox_gameDisableMapDraw_5d4594_2650672;
extern uint32_t dword_5d4594_1049508;
extern uint32_t dword_5d4594_1047516;
extern uint32_t dword_5d4594_1049500;
extern uint32_t dword_5d4594_1046868;
extern uint32_t dword_5d4594_1049504;
extern uint32_t dword_5d4594_832484;
extern void* nox_xxx_aClosewoodengat_587000_133480;
extern int nox_win_width;
extern int nox_win_height;

extern uint32_t nox_color_white_2523948;
extern uint32_t nox_color_black_2650656;

nox_window* nox_win_unk1 = 0;

uint32_t dword_587000_122856 = 0x1;
uint32_t dword_5d4594_831092 = 0;
uint32_t nox_player_netCode_85319C = 0;

// 457F15: variable 'v17' is possibly undefined

// 458291: variable 'v5' is possibly undefined

// 458F5A: variable 'v34' is possibly undefined

//----- (00459150) --------------------------------------------------------
//----- (00459DB0) --------------------------------------------------------
int sub_459DB0(nox_drawable* dr) {
	int a1 = dr;
	return *(uint32_t*)(a1 + 112) & 0x400000 && *(uint8_t*)(a1 + 116) & 8;
}

//----- (00459EC0) --------------------------------------------------------
int nox_xxx_cliNextMinimapObj_459EC0(int a1) {
	int next = *(uint32_t*)(a1 + 408);
	if (a1 && a1 == next) {
		printf("nox_xxx_cliNextMinimapObj_459EC0: infinite loop!\n");
		abort();
		return 0;
	}
	return next;
}

//----- (0045A010) --------------------------------------------------------
nox_drawable* sub_45A010(nox_drawable* dr) { return dr->field_104; }

//----- (0045A070) --------------------------------------------------------
nox_drawable* nox_drawable_next_45A070(nox_drawable* a1) {
	int result; // eax

	if (a1) {
		result = *(uint32_t*)((int)a1 + 368);
	} else {
		result = 0;
	}
	return result;
}

//----- (0045A990) --------------------------------------------------------
int nox_xxx_spriteSetActiveMB_45A990_drawable(int a1) {
	int result; // eax

	result = a1;
	*(uint32_t*)(a1 + 120) |= 4u;
	return result;
}

//----- (0045A9B0) --------------------------------------------------------
void sub_45A9B0(nox_drawable* a1p, nox_drawable* a2p) {
	int a1 = a1p;
	int a2 = a2p;
	int v2;          // esi
	int v3;          // ebp
	char* v4;        // eax
	char* v5;        // edi
	int* result = 0; // eax
	int v7;          // edi
	int v8;          // ebx
	int v9;          // eax
	int v10;         // esi
	long long v11;   // rax
	int v12;         // eax
	int* v13;        // esi
	int* v14;        // edi
	int* v15;        // edi
	int v16;         // [esp+Ch] [ebp-Ch]
	char* v17;       // [esp+10h] [ebp-8h]
	int* v18;        // [esp+14h] [ebp-4h]

	v2 = a1;
	v3 = 0;
	v16 = 0;
	v4 = nox_xxx_draw_452270(*(uint32_t*)(a1 + 492));
	v5 = v4;
	v17 = v4;
	result = (int*)nox_draw_getViewport_437250();
	v18 = result;
	if (v5 && result) {
		if (*(uint32_t*)(a1 + 120) & 0x1000000 && !(*(uint8_t*)(a1 + 280) & 0xC)) {
			v7 = *(uint32_t*)(a2 + 12) - *(uint32_t*)(a1 + 12);
			v8 = *(uint32_t*)(a2 + 16) - *(uint32_t*)(a1 + 16);
			v9 = sub_4522A0((int)v17);
			v10 = v9;
			if (v7 < v9 && v8 < v9 && v9 > 0) {
				v11 = (long long)sqrt((double)(v8 * v8 + v7 * v7 + 1));
				if ((int)v11 < v10) {
					v12 = 100 * (v10 - (int)v11) / v10;
					v3 = v12;
					if (v12 <= 100) {
						if (v12 < 0) {
							v3 = 0;
						}
					} else {
						v3 = 100;
					}
					v16 = 50 * (*(int*)(a1 + 12) - v18[6] - *v18) / (nox_win_width / 2);
				}
			}
			v2 = a1;
		}
		v13 = (int*)(v2 + 496);
		result = (int*)sub_452EB0(v13);
		v14 = result;
		if (v3) {
			if (result) {
				sub_452FE0((int)result, v16);
				result = (int*)sub_452F50((int)v14, v3);
			} else {
				result = nox_xxx_draw_452300(v17);
				v15 = result;
				if (result) {
					sub_452EE0((int)result, v3);
					sub_452F80((int)v15, v16);
					result = (int*)sub_452E90(v13, (int)v15);
				}
			}
		} else if (result) {
			result = (int*)sub_4523D0(result);
		}
	}
}

//----- (0045AB80) --------------------------------------------------------
int nox_xxx_spriteSetFrameMB_45AB80(int a1, int a2) {
	int result; // eax

	result = a1;
	if (!(*(uint8_t*)(a1 + 112) & 2) || !(*(uint32_t*)(a1 + 116) & 0x40000) || *(uint32_t*)(a1 + 276) != 8) {
		*(uint32_t*)(a1 + 312) = *(uint32_t*)(a1 + 308);
		*(uint32_t*)(a1 + 308) = a2;
	}
	return result;
}
