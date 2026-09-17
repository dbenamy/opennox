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
extern uint32_t dword_5d4594_1045576;
extern uint32_t dword_5d4594_1047532;
extern uint32_t dword_5d4594_1045596;
extern uint32_t dword_5d4594_1049692;
extern uint32_t dword_5d4594_1046652;
extern uint32_t dword_5d4594_1047536;
extern uint32_t dword_5d4594_1045420;
extern uint32_t dword_5d4594_1045556;
extern uint32_t dword_5d4594_831240;
extern uint32_t dword_5d4594_1045580;
extern uint32_t dword_5d4594_1045548;
extern uint32_t dword_5d4594_1045640;
extern uint32_t dword_5d4594_831260;
extern uint32_t dword_587000_129656;
extern uint32_t dword_5d4594_832480;
extern uint32_t dword_5d4594_1046360;
extern uint32_t dword_5d4594_1045540;
extern uint32_t dword_587000_126996;
extern uint32_t dword_5d4594_1045428;
extern uint32_t dword_5d4594_1045544;
extern uint32_t dword_5d4594_1045520;
extern uint32_t dword_5d4594_831244;
extern uint32_t dword_5d4594_831076;
extern uint32_t dword_5d4594_1045552;
extern uint32_t dword_5d4594_1049524;
extern uint32_t dword_5d4594_1047524;
extern uint32_t dword_5d4594_1046356;
extern uint32_t dword_5d4594_1046864;
extern uint32_t dword_5d4594_831256;
extern uint32_t dword_5d4594_1049536;
extern uint32_t dword_5d4594_1046944;
extern uint32_t dword_5d4594_1045424;
extern uint32_t dword_5d4594_1045584;
extern uint32_t dword_5d4594_1049696;
extern uint32_t dword_5d4594_1046540;
extern uint32_t dword_5d4594_1045588;
extern uint32_t dword_5d4594_831276;
extern uint32_t dword_5d4594_1046648;
extern uint32_t dword_5d4594_1046640;
extern uint32_t dword_5d4594_1045692;
extern uint32_t dword_5d4594_831084;
extern uint32_t dword_5d4594_1046956;
extern uint32_t dword_5d4594_1047936;
extern uint32_t dword_5d4594_1046536;
extern uint32_t dword_5d4594_1045436;
extern uint32_t dword_5d4594_831220;
extern uint32_t dword_5d4594_1045460;
extern uint32_t dword_587000_122848;
extern uint32_t dword_5d4594_1046932;
extern uint32_t dword_5d4594_1046528;
extern uint32_t dword_5d4594_1046948;
extern uint32_t dword_5d4594_830864;
extern uint32_t dword_5d4594_1047528;
extern uint32_t dword_5d4594_1046636;
extern uint32_t dword_5d4594_832520;
extern uint32_t dword_5d4594_832500;
extern uint32_t dword_5d4594_832528;
extern uint32_t dword_5d4594_1045688;
extern uint32_t dword_5d4594_832524;
extern uint32_t dword_5d4594_832512;
extern uint32_t dword_5d4594_832496;
extern uint32_t nox_xxx_aNox_cfg_0_587000_132136;
extern uint32_t dword_5d4594_832516;
extern uint32_t dword_5d4594_1046520;
extern uint32_t dword_5d4594_1049532;
extern uint32_t dword_5d4594_832508;
extern uint32_t dword_5d4594_832504;
extern uint32_t dword_5d4594_831224;
extern uint32_t dword_5d4594_1046928;
extern uint32_t dword_5d4594_1045508;
extern uint32_t dword_5d4594_1046504;
extern uint32_t dword_5d4594_832492;
extern uint32_t dword_5d4594_1047512;
extern uint32_t dword_5d4594_832532;
extern uint32_t dword_5d4594_1046656;
extern uint32_t dword_5d4594_1049484;
extern uint32_t dword_5d4594_1046516;
extern uint32_t dword_5d4594_832536;
extern uint32_t dword_5d4594_1046952;
extern void* dword_587000_81128;
extern void* dword_587000_122852;
extern uint32_t dword_5d4594_1046532;
extern uint32_t dword_5d4594_1049496;
extern uint32_t dword_5d4594_1047932;
extern uint32_t dword_5d4594_1049512;
extern uint32_t nox_wnd_briefing_831232;
extern uint32_t dword_5d4594_1045432;
extern uint32_t dword_5d4594_1046924;
extern void* dword_587000_127004;
extern uint32_t dword_5d4594_1045468;
extern uint32_t dword_5d4594_1045532;
extern uint32_t dword_5d4594_1045536;
extern uint32_t dword_5d4594_1046508;
extern uint32_t dword_5d4594_1045480;
extern uint32_t dword_5d4594_1047520;
extern uint32_t nox_xxx_aNox_cfg_0_587000_132132;
extern uint32_t dword_5d4594_1045484;
extern uint32_t dword_5d4594_1045464;
extern uint32_t dword_5d4594_1049520;
extern uint32_t dword_5d4594_1045528;
extern uint32_t dword_5d4594_1046936;
extern uint32_t dword_5d4594_1047540;
extern uint32_t dword_5d4594_831236;
extern uint32_t dword_5d4594_1046496;
extern uint32_t dword_5d4594_1046872;
extern uint32_t nox_gameDisableMapDraw_5d4594_2650672;
extern uint32_t dword_5d4594_1045636;
extern uint32_t dword_5d4594_1049508;
extern uint32_t dword_5d4594_1047516;
extern uint32_t dword_5d4594_1049500;
extern uint32_t dword_5d4594_1045604;
extern uint32_t dword_5d4594_1046868;
extern uint32_t dword_5d4594_1049504;
extern uint32_t dword_5d4594_832484;
extern uint32_t dword_5d4594_1045684;
extern uint32_t dword_5d4594_1045516;
extern void* nox_xxx_aClosewoodengat_587000_133480;
extern uint32_t dword_5d4594_1046492;
extern int nox_win_width;
extern int nox_win_height;

extern uint32_t nox_color_white_2523948;
extern uint32_t nox_color_black_2650656;

nox_window* nox_win_unk1 = 0;
nox_window* dword_5d4594_1046512 = 0;

uint32_t dword_587000_122856 = 0x1;
uint32_t dword_5d4594_831092 = 0;
uint32_t nox_player_netCode_85319C = 0;

//----- (0044D960) --------------------------------------------------------
void sub_44D960() { dword_587000_122848 = 0; }

//----- (0044D970) --------------------------------------------------------
int sub_44D970() {
	int result; // eax

	result = dword_5d4594_831092;
	if (dword_5d4594_831092) {
		dword_587000_122848 = 1;
	}
	return result;
}

//----- (0044D990) --------------------------------------------------------
int sub_44D990() { return dword_587000_122848; }

//----- (00450750) --------------------------------------------------------
unsigned char sub_450750() { return getMemByte(0x5D4594, 831252); }

//----- (00450760) --------------------------------------------------------
char sub_450760(char a1) {
	char result; // al

	result = a1;
	*getMemU8Ptr(0x5D4594, 831252) = a1;
	return result;
}

//----- (00451850) --------------------------------------------------------
int sub_451850(int a2, void* a3p) {
	int a3 = a3p;
	int v2;            // edi
	unsigned char* v3; // esi
	int result;        // eax

	v2 = 0;
	v3 = getMemAt(0x5D4594, 840712);
	do {
		sub_451920((uint32_t*)v3 - 21);
		*(uint32_t*)v3 = nox_xxx_getSndName_40AF80(v2);
		v3 += 200;
		++v2;
	} while ((int)v3 < (int)getMemAt(0x5D4594, 1045312));
	dword_5d4594_1045420 = a3;
	dword_5d4594_1045428 = a2;
	if (a3) {
		dword_5d4594_1045424 = sub_4BD340(a3, 0x100000, 200, 0x2000);
		dword_5d4594_1045436 = sub_4BD280(200, 576);
	}
	if (!dword_5d4594_1045424 || !dword_5d4594_1045420 || !dword_5d4594_1045428 || !dword_5d4594_1045436) {
		return 0;
	}
	nox_common_list_clear_425760(getMemAt(0x5D4594, 840612));
	sub_4864A0(getMemAt(0x5D4594, 1045228));
	result = 1;
	*(uint32_t*)(dword_5d4594_1045428 + 184) = getMemAt(0x5D4594, 1045228);
	dword_5d4594_1045432 = 1;
	return result;
}

//----- (00451920) --------------------------------------------------------
int sub_451920(uint32_t* a2) {
	*a2 = 0;
	a2[1] = 0;
	a2[2] = 0;
	a2[14] = 0;
	a2[15] = 0;
	a2[19] = 0;
	a2[20] = 0;
	a2[12] = 1;
	a2[48] = 0;
	a2[18] = 0;
	a2[17] = 0;
	a2[25] = 0;
	a2[26] = 0;
	a2[16] = 600;
	return sub_4862E0((int)(a2 + 4), 0x4000);
}

//----- (00451970) --------------------------------------------------------
void sub_451970() {
	sub_4521F0();
	sub_452230();
	if (dword_5d4594_1045424) {
		sub_4BD3C0(*(void**)&dword_5d4594_1045424);
		dword_5d4594_1045424 = 0;
	}
	if (dword_5d4594_1045436) {
		sub_4BD2D0(*(void**)&dword_5d4594_1045436);
		dword_5d4594_1045436 = 0;
	}
	dword_5d4594_1045432 = 0;
}

//----- (004519C0) --------------------------------------------------------
void sub_4519C0() {
	int result;        // eax
	int v1;            // esi
	int v2;            // eax
	int v3;            // ebp
	int v4;            // eax
	unsigned char* v5; // edi
	unsigned char* v6; // esi
	unsigned char* v7; // edi
	int v8;            // eax
	int v9;            // eax
	int v10;           // eax

	result = dword_5d4594_1045432;
	if (!dword_5d4594_1045432) {
		return;
	}
	result = *getMemU32Ptr(0x5D4594, 1045448);
	if (*getMemU32Ptr(0x5D4594, 1045448)) {
		return;
	}
	*getMemU32Ptr(0x5D4594, 1045448) = 1;
	sub_486520(*(unsigned int**)&dword_587000_127004);
	v1 = *getMemU32Ptr(0x5D4594, 840612);
	++*getMemU32Ptr(0x5D4594, 1045440);
	if (*(unsigned char**)getMemAt(0x5D4594, 840612) != getMemAt(0x5D4594, 840612)) {
		do {
			v2 = *(uint32_t*)(v1 + 36);
			if (*(uint32_t*)(v2 + 100) != *getMemU32Ptr(0x5D4594, 1045440)) {
				nox_common_list_clear_425760((uint32_t*)(v2 + 88));
				*(uint32_t*)(*(uint32_t*)(v1 + 36) + 52) = 0;
				*(uint32_t*)(*(uint32_t*)(v1 + 36) + 100) = *getMemU32Ptr(0x5D4594, 1045440);
			}
			sub_486520((unsigned int*)(v1 + 184));
			if (*(uint32_t*)(v1 + 28) != 4) {
				sub_451BE0(v1);
			}
			v1 = *(uint32_t*)v1;
		} while ((unsigned char*)v1 != getMemAt(0x5D4594, 840612));
		v1 = *getMemU32Ptr(0x5D4594, 840612);
		if (*(unsigned char**)getMemAt(0x5D4594, 840612) != getMemAt(0x5D4594, 840612)) {
			do {
				sub_452510(v1);
				v1 = *(uint32_t*)v1;
			} while ((unsigned char*)v1 != getMemAt(0x5D4594, 840612));
			v1 = *getMemU32Ptr(0x5D4594, 840612);
		}
	}
	v3 = 0;
	sub_452010();
	if ((unsigned char*)v1 != getMemAt(0x5D4594, 840612)) {
		do {
			v4 = *(uint32_t*)(v1 + 176);
			v5 = *(unsigned char**)v1;
			if (!v4 || v1 != *(uint32_t*)(v4 + 152)) {
				sub_4523D0((uint32_t*)v1);
			}
			if (*(uint8_t*)(v1 + 24) & 1) {
				sub_451FE0(v1);
			} else {
				v3 += (unsigned int)(33 * (*(uint32_t*)(*(uint32_t*)(v1 + 36) + 20) >> 16)) >> 14;
				sub_452050((uint32_t*)v1);
			}
			v1 = (int)v5;
		} while (v5 != getMemAt(0x5D4594, 840612));
	}
	if (v3 <= 100) {
		sub_486350((int)getMemAt(0x5D4594, 1045228), 0x4000);
	} else {
		sub_486350((int)getMemAt(0x5D4594, 1045228), 0x190000u / v3);
	}
	result = sub_486520(getMemUintPtr(0x5D4594, 1045228));
	v6 = *(unsigned char**)getMemAt(0x5D4594, 840612);
	if (*(unsigned char**)getMemAt(0x5D4594, 840612) != getMemAt(0x5D4594, 840612)) {
		do {
			v7 = *(unsigned char**)v6;
			result = *((uint32_t*)v6 + 7);
			if (result == 1) {
				sub_451DC0((int)v6);
				v8 = sub_451CA0(v6);
				*((uint32_t*)v6 + 74) = v8;
				if (!v8) {
					do {
						if (!sub_452120((int)v6)) {
							break;
						}
						v7 = *(unsigned char**)v6;
						sub_451DC0((int)v6);
						v9 = sub_451CA0(v6);
						*((uint32_t*)v6 + 74) = v9;
					} while (!v9);
				}
				v10 = sub_451CA0(v6);
				*((uint32_t*)v6 + 74) = v10;
				if (!v10 || (result = sub_452490(v6)) == 0) {
					sub_4523D0(v6);
					result = sub_451FE0((int)v6);
				}
			}
			v6 = v7;
		} while (v7 != getMemAt(0x5D4594, 840612));
	}
	*getMemU32Ptr(0x5D4594, 1045448) = 0;
}

//----- (00451BE0) --------------------------------------------------------
int sub_451BE0(int a1) {
	int v1;          // eax
	int v2;          // edi
	unsigned int v3; // ebx
	uint32_t* v4;    // esi
	int v5;          // eax
	int v6;          // eax
	uint32_t* v7;    // ebx
	int result;      // eax
	int v9;          // esi
	uint32_t* v10;   // esi

	v1 = a1;
	v2 = *(uint32_t*)(a1 + 36);
	v3 = *(uint32_t*)(a1 + 188) >> 16;
	v4 = *(uint32_t**)(v2 + 88);
	if (v4 != (uint32_t*)(v2 + 88)) {
		do {
			v5 = (v4[44] >> 16) - v3;
			if (v5 < 0) {
				v5 = v3 - (v4[44] >> 16);
			}
			if (v5 >= (*(uint32_t*)(v2 + 20) >> 16) / 10) {
				if (v4[44] >> 16 < v3) {
					break;
				}
			} else {
				v6 = v4[4];
				if (*(uint8_t*)(v2 + 4) & 0x10) {
					if (v6) {
						break;
					}
				} else if (!v6) {
					break;
				}
			}
			v4 = (uint32_t*)*v4;
		} while (v4 != (uint32_t*)(v2 + 88));
		v1 = a1;
	}
	v7 = (uint32_t*)(v1 + 12);
	sub_425770((uint32_t*)(v1 + 12));
	nox_common_list_append_4258E0((int)v4, v7);
	result = *(uint32_t*)(v2 + 56);
	v9 = *(uint32_t*)(v2 + 52) + 1;
	*(uint32_t*)(v2 + 52) = v9;
	if (result) {
		if (v9 > result) {
			v10 = (uint32_t*)(*(uint32_t*)(v2 + 92) - 12);
			nox_common_list_remove_425920(*(uint32_t***)(v2 + 92));
			sub_4523D0(v10);
			result = *(uint32_t*)(v2 + 52) - 1;
			*(uint32_t*)(v2 + 52) = result;
		}
	}
	return result;
}

//----- (00451CA0) --------------------------------------------------------
int sub_451CA0(uint32_t* a1) {
	int v1;       // ecx
	int v3;       // eax
	uint32_t* v4; // ecx

	v1 = a1[42];
	a1[108] = v1;
	if (!v1) {
		return 0;
	}
	v3 = 0;
	if (v1 > 0) {
		v4 = a1 + 76;
		do {
			*v4 = v3++;
			++v4;
		} while (v3 < a1[108]);
	}
	a1[43] = -1;
	return sub_451CF0(a1);
}

//----- (00451F30) --------------------------------------------------------
int sub_451F30(int a1, int a2) {
	int v2;     // edx
	int result; // eax

	*(uint32_t*)(a1 + 4 * *(uint32_t*)(a1 + 168) + 40) =
		sub_4BD470(*(uint32_t***)&dword_5d4594_1045424, *(short*)(*(uint32_t*)(a1 + 36) + 2 * a2 + 128));
	v2 = *(uint32_t*)(a1 + 168);
	result = *(uint32_t*)(a1 + 4 * v2 + 40);
	if (result) {
		sub_4BD650(*(uint32_t*)(a1 + 4 * v2 + 40));
		result = *(uint32_t*)(a1 + 168) + 1;
		*(uint32_t*)(a1 + 168) = result;
	}
	return result;
}

//----- (00451F90) --------------------------------------------------------
int sub_451F90(int a1) {
	int v1;     // edi
	int result; // eax
	int* v3;    // esi

	v1 = 0;
	result = *(uint32_t*)(a1 + 168);
	if (result <= 0) {
		*(uint32_t*)(a1 + 168) = 0;
	} else {
		v3 = (int*)(a1 + 40);
		do {
			sub_4BD660(*v3);
			*v3 = 0;
			result = *(uint32_t*)(a1 + 168);
			++v1;
			++v3;
		} while (v1 < result);
		*(uint32_t*)(a1 + 168) = 0;
	}
	return result;
}

//----- (00451FE0) --------------------------------------------------------
int sub_451FE0(int a1) {
	nox_common_list_remove_425920((uint32_t**)a1);
	*(uint32_t*)(a1 + 280) = 0;
	return sub_4BD300(*(uint32_t**)&dword_5d4594_1045436, a1);
}

//----- (00452010) --------------------------------------------------------
int sub_452010() {
	unsigned char* v0; // esi
	int v1;            // ebx
	int v2;            // edi

	v0 = getMemAt(0x5D4594, 839892);
	v1 = 6;
	do {
		v2 = 10;
		do {
			nox_common_list_clear_425760(v0);
			v0 += 12;
			--v2;
		} while (v2);
		--v1;
	} while (v1);
	return ++*getMemU32Ptr(0x5D4594, 1045444);
}

//----- (00452050) --------------------------------------------------------
void sub_452050(uint32_t* a1) {
	uint32_t* v1;      // esi
	int v2;            // edi
	unsigned int v3;   // ebx
	unsigned char* v4; // ebp
	uint32_t* result;  // eax
	uint32_t** v6;     // esi
	uint32_t** v7;     // esi
	uint32_t* v8;      // esi

	v1 = (uint32_t*)a1[9];
	v2 = v1[12] + a1[75];
	v3 = (a1[47] >> 16) / 0x666u;
	v4 = getMemAt(0x5D4594, 839892 + 120 * v2);
	if (v1[26] == *getMemU32Ptr(0x5D4594, 1045444)) {
		result = (uint32_t*)v1[27];
		if (v2 <= (int)result) {
			if ((uint32_t*)v2 == result && v3 > v1[31]) {
				v1[31] = v3;
				v7 = (uint32_t**)(v1 + 28);
				nox_common_list_remove_425920(v7);
				nox_common_list_append_4258E0((int)&v4[12 * v3], v7);
			}
		} else {
			v1[27] = v2;
			v1[31] = v3;
			v6 = (uint32_t**)(v1 + 28);
			nox_common_list_remove_425920(v6);
			nox_common_list_append_4258E0((int)&v4[12 * v3], v6);
		}
	} else {
		v1[26] = *getMemU32Ptr(0x5D4594, 1045444);
		v1[27] = v2;
		v1[31] = v3;
		v8 = v1 + 28;
		sub_425770(v8);
		nox_common_list_append_4258E0((int)&v4[12 * v3], v8);
	}
}

//----- (00452120) --------------------------------------------------------
int* sub_452120(int a1) {
	int v1;            // ebp
	int* result;       // eax
	int* v3;           // ebx
	unsigned char* v4; // esi
	unsigned char* v5; // edi

	v1 = 0;
	result = sub_4521A0(*(uint32_t*)(a1 + 300) + *(uint32_t*)(*(uint32_t*)(a1 + 36) + 48));
	v3 = result;
	if (result) {
		sub_452190((int)result);
		v4 = *(unsigned char**)getMemAt(0x5D4594, 840612);
		if (*(unsigned char**)getMemAt(0x5D4594, 840612) != getMemAt(0x5D4594, 840612)) {
			do {
				v5 = *(unsigned char**)v4;
				if (*((int**)v4 + 9) == v3) {
					sub_4523D0(v4);
					sub_451FE0((int)v4);
					v1 = 1;
				}
				v4 = v5;
			} while (v5 != getMemAt(0x5D4594, 840612));
		}
		result = (int*)v1;
	}
	return result;
}

//----- (00452190) --------------------------------------------------------
void sub_452190(int a1) { nox_common_list_remove_425920((uint32_t**)(a1 + 112)); }

//----- (004521A0) --------------------------------------------------------
int* sub_4521A0(int a1) {
	int v1;            // ebp
	unsigned char* v2; // ebx
	int v3;            // edi
	int* v4;           // esi
	int* v5;           // eax

	v1 = 0;
	v2 = getMemAt(0x5D4594, 839892);
	if (a1 > 0) {
		while (1) {
			v3 = 0;
			v4 = (int*)v2;
			do {
				v5 = nox_common_list_getFirstSafe_425890(v4);
				if (v5) {
					return v5 - 28;
				}
				++v3;
				v4 += 3;
			} while (v3 < 10);
			++v1;
			v2 += 120;
			if (v1 < a1) {
				continue;
			}
			break;
		}
	}
	return 0;
}

//----- (004521F0) --------------------------------------------------------
int sub_4521F0() {
	int result;        // eax
	unsigned char* v1; // esi
	unsigned char* v2; // edi

	result = dword_5d4594_1045432;
	if (dword_5d4594_1045432) {
		v1 = *(unsigned char**)getMemAt(0x5D4594, 840612);
		if (*(unsigned char**)getMemAt(0x5D4594, 840612) != getMemAt(0x5D4594, 840612)) {
			do {
				v2 = *(unsigned char**)v1;
				sub_4523D0(v1);
				result = sub_451FE0((int)v1);
				v1 = v2;
			} while (v2 != getMemAt(0x5D4594, 840612));
		}
	}
	return result;
}

//----- (00452230) --------------------------------------------------------
int***** sub_452230() {
	int***** result; // eax
	int**** v1;      // esi

	result = *(int******)&dword_5d4594_1045432;
	if (dword_5d4594_1045432) {
		result = *(int******)getMemAt(0x5D4594, 840612);
		if (*(unsigned char**)getMemAt(0x5D4594, 840612) != getMemAt(0x5D4594, 840612)) {
			do {
				v1 = *result;
				if ((uint8_t)result[6] & 1) {
					sub_451FE0((int)result);
				}
				result = (int*****)v1;
			} while (v1 != (int****)getMemAt(0x5D4594, 840612));
		}
	}
	return result;
}

//----- (00452300) --------------------------------------------------------
uint32_t* nox_xxx_draw_452300(uint32_t* a1) {
	uint32_t* v1; // esi

	if (!dword_5d4594_1045432) {
		return 0;
	}
	if (!dword_587000_126996) {
		return 0;
	}
	if (!*a1) {
		return 0;
	}
	v1 = sub_4BD2E0(*(uint32_t***)&dword_5d4594_1045436);
	if (!v1) {
		sub_452230();
		v1 = sub_4BD2E0(*(uint32_t***)&dword_5d4594_1045436);
		if (!v1) {
			return 0;
		}
	}
	memset(v1, 0, 0x240u);
	v1[9] = a1;
	sub_425770(v1);
	v1[7] = 0;
	v1[75] = 0;
	v1[142] = 0;
	v1[108] = 0;
	v1[42] = 0;
	sub_4864A0(v1 + 46);
	nox_common_list_append_4258E0((int)getMemAt(0x5D4594, 840612), v1);
	v1[70] = (*getMemU32Ptr(0x587000, 127000))++;
	return v1;
}

//----- (004523D0) --------------------------------------------------------
int sub_4523D0(void* a1p) {
	uint32_t* a1 = a1p;
	int result = 0; // eax

	if (!(a1[6] & 1)) {
		sub_452410((int)a1);
		sub_451F90((int)a1);
		a1[7] = 4;
		a1[70] = 0;
		result = a1[6];
		LOBYTE(result) = result | 1;
		a1[6] = result;
	}
	return result;
}

//----- (00452410) --------------------------------------------------------
int sub_452410(int a1) {
	int result; // eax

	result = *(uint32_t*)(a1 + 176);
	if (result && a1 == *(uint32_t*)(result + 152)) {
		if (*(uint8_t*)(a1 + 24) & 2) {
			sub_4BDA80(*(uint32_t*)(a1 + 176));
		}
		sub_4BDB30(*(uint32_t*)(a1 + 176));
		*(uint32_t*)(*(uint32_t*)(a1 + 176) + 152) = 0;
		*(uint32_t*)(*(uint32_t*)(a1 + 176) + 148) = 0;
		result = *(uint32_t*)(a1 + 176);
		*(uint32_t*)(result + 140) = 0;
		*(uint32_t*)(*(uint32_t*)(a1 + 176) + 144) = 0;
		*(uint32_t*)(*(uint32_t*)(a1 + 176) + 112) = 0;
		*(uint32_t*)(a1 + 176) = 0;
	}
	return result;
}

//----- (00452490) --------------------------------------------------------
int sub_452490(uint32_t* a1) {
	int v1; // eax
	int v3; // edi
	int v4; // eax

	v1 = a1[44];
	if (a1 != *(uint32_t**)(v1 + 152)) {
		return 0;
	}
	v3 = a1[74];
	sub_4BDB90((uint32_t*)v1, (uint32_t*)a1[74]);
	a1[7] = 3;
	v4 = a1[6];
	LOBYTE(v4) = v4 | 2;
	a1[6] = v4;
	a1[74] = 0;
	if (!sub_4BDB40(a1[44])) {
		return 1;
	}
	a1[7] = 1;
	a1[74] = v3;
	a1[6] &= 0xFFFFFFFD;
	return 0;
}

//----- (00452510) --------------------------------------------------------
void sub_452510(int a3) {
	int v1; // eax
	int v2; // eax

	if (!dword_587000_126996) {
		*(uint32_t*)(a3 + 28) = 4;
	}
	while (1) {
		v1 = *(uint32_t*)(a3 + 28);
		if (!v1) {
			break;
		}
		v2 = v1 - 2;
		if (v2) {
			uint32_t v3 = v2 - 2;
			if (!v3) {
				sub_4523D0((uint32_t*)a3);
			}
			return;
		}
		if (nox_platform_get_ticks() <= *(uint64_t*)(a3 + 288)) {
			return;
		}
		*(uint32_t*)(a3 + 28) = *(uint32_t*)(a3 + 32);
	}
	if (!sub_452580((uint32_t*)a3)) {
		sub_4523D0((uint32_t*)a3);
	}
}

//----- (00452690) --------------------------------------------------------
long long sub_452690(int a3, long long a4, int a5) {
	long long result; // rax

	*(uint32_t*)(a3 + 32) = a5;
	result = a4 + nox_platform_get_ticks();
	*(uint64_t*)(a3 + 288) = result;
	*(uint32_t*)(a3 + 28) = 2;
	return result;
}

//----- (004526D0) --------------------------------------------------------
int sub_4526D0(int a1) {
	*(uint32_t*)(*(uint32_t*)(a1 + 152) + 28) = 4;
	return 0;
}

//----- (004526F0) --------------------------------------------------------
int sub_4526F0(int a1) {
	uint32_t* v1; // esi
	int v2;       // eax

	v1 = *(uint32_t**)(a1 + 152);
	v1[6] &= 0xFFFFFFFD;
	v2 = 4;
	if (v1[7] != 4) {
		if (v1[74] || v1[142]) {
			v2 = 1;
		} else {
			v1[71] = 0;
		}
		if (v1[71]) {
			sub_452690((int)v1, (unsigned int)v1[71], v2);
			v1[71] = 0;
			return 0;
		}
		v1[7] = v2;
	}
	return 0;
}

//----- (00452810) --------------------------------------------------------
int* sub_452810(int a1, char a2) {
	int* v2; // esi
	int* v3; // eax

	v2 = 0;
	if (dword_5d4594_1045428) {
		v3 = sub_487810(*(int*)&dword_5d4594_1045428, 1);
		v2 = v3;
		if (v3) {
			if (v3[31] & 0x15 && v3[30] > a1) {
				return 0;
			}
			sub_4BDA80((int)v3);
			v2[29] = dword_587000_127004;
			v2[30] = a1;
			if (a2 & 1) {
				v2[32] = -1;
			} else {
				v2[32] = 0;
			}
			sub_486320(v2 + 4, 0x4000);
		}
	}
	return v2;
}

//----- (00452D80) --------------------------------------------------------
void nox_xxx_clientPlaySoundSpecial_452D80(int a1, int a2) {
#ifdef NOX_PORT_TEST_CLIENT_SOUND
	extern void nox_porttest_client_sound(int id, int volume);
	nox_porttest_client_sound(a1, a2);
#endif
	uint32_t* result; // eax
	uint32_t* v3;     // esi

	result = nox_xxx_draw_452270(a1);
	if (!result) {
		return;
	}
	result = nox_xxx_draw_452300(result);
	v3 = result;
	if (!result) {
		return;
	}
	sub_452EE0((int)result, a2);
	sub_452510((int)v3);
}

//----- (00452DC0) --------------------------------------------------------
void sub_452DC0(int a1, int a2, int a3) {
	uint32_t* result; // eax
	uint32_t* v4;     // esi

	result = nox_xxx_draw_452270(a1);
	if (!result) {
		return;
	}
	result = nox_xxx_draw_452300(result);
	v4 = result;
	if (!result) {
		return;
	}
	sub_452EE0((int)result, a2);
	sub_452F80((int)v4, a3);
	sub_452510((int)v4);
}

//----- (00452E10) --------------------------------------------------------
void sub_452E10(int a1, int a2, int a3) {
	uint32_t* result; // eax
	uint32_t* v4;     // esi

	result = nox_xxx_draw_452270(a1);
	if (!result) {
		return;
	}
	result = nox_xxx_draw_452300(result);
	v4 = result;
	if (!result) {
		return;
	}
	sub_452EE0((int)result, a2);
	sub_452F80((int)v4, a3);
	v4[75] = 2;
	sub_452510((int)v4);
}

//----- (00452E90) --------------------------------------------------------
int sub_452E90(uint32_t* a1, int a2) {
	int result; // eax

	result = a2;
	*a1 = a2;
	if (a2) {
		a1[1] = *(uint32_t*)(a2 + 280);
		result = *(uint32_t*)(a2 + 36);
		a1[2] = result;
	}
	return result;
}

//----- (00452EB0) --------------------------------------------------------
int sub_452EB0(int* a1) {
	int result; // eax

	result = *a1;
	if (*a1 && (a1[2] != *(uint32_t*)(result + 36) || a1[1] != *(uint32_t*)(result + 280))) {
		result = 0;
		*a1 = 0;
	}
	return result;
}

//----- (00452EE0) --------------------------------------------------------
int sub_452EE0(int a1, int a2) {
	int v2; // eax

	v2 = sub_452F10(a1, a2);
	sub_486320((uint32_t*)(a1 + 184), v2);
	return sub_4863B0((unsigned int*)(a1 + 184));
}

//----- (00452F10) --------------------------------------------------------
unsigned int sub_452F10(int a1, int a2) {
	int v2; // ecx

	v2 = a2;
	if (a2 <= 100) {
		if (a2 < 0) {
			v2 = 0;
		}
	} else {
		v2 = 100;
	}
	return (unsigned int)(163 * v2 * (*(uint32_t*)(*(uint32_t*)(a1 + 36) + 20) >> 16)) >> 14;
}

//----- (00452F50) --------------------------------------------------------
int sub_452F50(int a1, int a2) {
	int v2; // eax

	v2 = sub_452F10(a1, a2);
	return sub_486350(a1 + 184, v2);
}

//----- (00452F80) --------------------------------------------------------
uint32_t* sub_452F80(int a1, int a2) {
	int v2; // eax

	v2 = sub_452FA0(a2);
	return sub_486320((uint32_t*)(a1 + 248), v2);
}

//----- (00452FA0) --------------------------------------------------------
int sub_452FA0(int a1) {
	int v1; // eax

	v1 = a1;
	if (a1 <= 50) {
		if (a1 < -50) {
			v1 = -50;
		}
	} else {
		v1 = 50;
	}
	return (v1 * 8192) / 50 + 8192;
}

//----- (00452FE0) --------------------------------------------------------
int sub_452FE0(int a1, int a2) {
	int v2; // eax

	v2 = sub_452FA0(a2);
	return sub_486350(a1 + 248, v2);
}

//----- (00453050) --------------------------------------------------------
void sub_453050() { dword_587000_126996 = 0; }

//----- (00453060) --------------------------------------------------------
void nox_xxx____setargv_9_453060() { dword_587000_126996 = 1; }

//----- (00453070) --------------------------------------------------------
int sub_453070() { return dword_587000_126996; }

//----- (00453080) --------------------------------------------------------
int sub_453080(char a1) {
	int result; // eax

	if (dword_5d4594_1045460) {
		result = sub_453690(1 << a1);
	} else {
		result = sub_453660(1 << a1);
	}
	return result;
}

//----- (004532E0) --------------------------------------------------------
uint32_t* sub_4532E0() {
	int v0;           // esi
	int v1;           // edi
	int v2;           // ebx
	uint32_t* result; // eax

	v0 = nox_xxx_guiFontHeightMB_43F320(*(uint32_t*)(dword_5d4594_1045464 + 236)) + 1;
	sub_46AB20(*(uint32_t**)&dword_5d4594_1045464, *(uint32_t*)(dword_5d4594_1045464 + 8), 15 * v0 + 2);
	v1 = 1520;
	v2 = *(uint32_t*)(dword_5d4594_1045464 + 20) + v0 + 2;
	do {
		result = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045468, v1);
		result[5] = v2;
		v2 += v0;
		++v1;
	} while (v1 < 1534);
	return result;
}

//----- (00453350) --------------------------------------------------------
int sub_453350(int a1, int a2) {
	int result; // eax
	int xLeft;  // [esp+4h] [ebp-8h]
	int yTop;   // [esp+8h] [ebp-4h]

	nox_client_wndGetPosition_46AA60((uint32_t*)a1, &xLeft, &yTop);
	if ((signed char)*(uint8_t*)(a1 + 4) >= 0) {
		if (*(uint32_t*)(a2 + 20) != 0x80000000) {
			nox_client_drawRectFilledAlpha_49CF10(xLeft, yTop, *(uint32_t*)(a1 + 8), *(uint32_t*)(a1 + 12));
		}
		result = 1;
	} else {
		nox_client_drawImageAt_47D2C0(*(uint32_t*)(a2 + 24), xLeft, yTop);
		result = 1;
	}
	return result;
}

//----- (004533D0) --------------------------------------------------------
int sub_4533D0(int a1, int a2, int a3, int a4) {
	int v3;       // esi
	int v5;       // eax
	int v6;       // eax
	int v7;       // edx
	int v8;       // ecx
	int v9;       // edx
	int v10;      // ecx
	wchar2_t* v11; // [esp-4h] [ebp-Ch]

	if (a2 == 0x4000) {
		if ((uint32_t*)a3 == nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045468, 1513) ||
			(uint32_t*)a3 == nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045468, 1514)) {
			nox_window_call_field_94(*(int*)&dword_5d4594_1045464, 0x4000, a3, 0);
			sub_453750();
		}
		return 0;
	}
	if (a2 != 16391) {
		return 0;
	}
	v3 = nox_xxx_wndGetID_46B0A0((int*)a3);
	switch (v3) {
	case 1513:
	case 1514:
		nox_window_call_field_94(*(int*)&dword_5d4594_1045464, 0x4000, a3, 0);
		sub_453750();
		return 0;
	case 1515:
		if (dword_5d4594_1045460) {
			*getMemU32Ptr(0x5D4594, 1045456) = -1;
		} else {
			*getMemU32Ptr(0x5D4594, 1045452) = -1;
		}
		sub_453750();
		sub_459D50(1);
		break;
	case 1516:
		if (dword_5d4594_1045460) {
			*getMemU32Ptr(0x5D4594, 1045456) = 0;
		} else {
			*getMemU32Ptr(0x5D4594, 1045452) = 0;
		}
		sub_453750();
		sub_459D50(1);
		break;
	case 1520:
	case 1521:
	case 1522:
	case 1523:
	case 1524:
	case 1525:
	case 1526:
	case 1527:
	case 1528:
	case 1529:
	case 1530:
	case 1531:
	case 1532:
	case 1533:
		v5 = sub_4A4800(*(uint32_t*)(dword_5d4594_1045464 + 32));
		v11 = (wchar2_t*)nox_window_call_field_94(*(int*)&dword_5d4594_1045464, 16406, v5 + v3 - 1520, 0);
		if (dword_5d4594_1045460) {
			v6 = sub_415DA0(v11);
		} else {
			v6 = sub_415960(v11);
		}
		if (*(uint8_t*)(a3 + 36) & 4) {
			if (dword_5d4594_1045460) {
				sub_453640(getMemAt(0x5D4594, 1045456), v6, 0);
			} else {
				v7 = 0;
				if (v6 > 0) {
					do {
						v8 = v6 >> 8;
						if (v6 >> 8 > 0) {
							v6 >>= 8;
						}
						++v7;
					} while (v8 > 0);
				}
				sub_453620(getMemAt(0x5D4594, 1045451 + v7), v6, 0);
			}
		} else if (dword_5d4594_1045460) {
			sub_453640(getMemAt(0x5D4594, 1045456), v6, 1);
		} else {
			v9 = 0;
			if (v6 > 0) {
				do {
					v10 = v6 >> 8;
					if (v6 >> 8 > 0) {
						v6 >>= 8;
					}
					++v9;
				} while (v10 > 0);
			}
			sub_453620(getMemAt(0x5D4594, 1045451 + v9), v6, 1);
		}
		sub_459D50(1);
		break;
	default:
		break;
	}
	nox_xxx_clientPlaySoundSpecial_452D80(766, 100);
	return 0;
}

//----- (004535E0) --------------------------------------------------------
int* sub_4535E0(int* a1) {
	int* result; // eax

	result = a1;
	*getMemU32Ptr(0x5D4594, 1045452) = *a1;
	return result;
}

//----- (004535F0) --------------------------------------------------------
int sub_4535F0(int a1) {
	int result; // eax

	result = a1;
	*getMemU32Ptr(0x5D4594, 1045456) = a1;
	return result;
}

//----- (00453600) --------------------------------------------------------
char* sub_453600() { return (char*)getMemAt(0x5D4594, 1045452); }

//----- (00453610) --------------------------------------------------------
int sub_453610() { return *getMemU32Ptr(0x5D4594, 1045456); }

//----- (00453620) --------------------------------------------------------
uint8_t* sub_453620(uint8_t* a1, char a2, int a3) {
	uint8_t* result; // eax

	result = a1;
	if (a3) {
		*a1 |= a2;
	} else {
		*a1 &= ~a2;
	}
	return result;
}

//----- (00453640) --------------------------------------------------------
uint32_t* sub_453640(uint32_t* a1, int a2, int a3) {
	uint32_t* result; // eax

	result = a1;
	if (a3) {
		*a1 |= a2;
	} else {
		*a1 &= ~a2;
	}
	return result;
}

//----- (00453660) --------------------------------------------------------
int sub_453660(int a1) {
	int v1; // ecx
	int v2; // edx
	int v3; // eax

	v1 = a1;
	v2 = 0;
	if (a1 > 0) {
		do {
			v3 = v1 >> 8;
			if (v1 >> 8 > 0) {
				v1 >>= 8;
			}
			++v2;
		} while (v3 > 0);
	}
	return ((unsigned char)v1 & getMemByte(0x5D4594, 1045451 + v2)) != 0;
}

//----- (00453690) --------------------------------------------------------
int sub_453690(int a1) { return (a1 & *getMemU32Ptr(0x5D4594, 1045456)) != 0; }

//----- (004536B0) --------------------------------------------------------
int sub_4536B0(uint32_t* a1) {
	int v1;       // ebx
	int v2;       // edi
	uint32_t* v3; // esi
	int v4;       // ebp
	int result;   // eax

	v1 = 1;
	*a1 = -1;
	v2 = 1;
	v3 = a1;
	v4 = 27;
	do {
		result = sub_415840((char*)v2);
		if (result) {
			result = nox_xxx_getUnitDefDd10_4E3BA0(result);
			if (!result) {
				result = (unsigned char)~(uint8_t)v1;
				*(uint8_t*)v3 &= ~(uint8_t)v1;
			}
		}
		if (v1 == 128) {
			v1 = 1;
			v3 = (uint32_t*)((char*)v3 + 1);
		} else {
			v1 *= 2;
		}
		v2 *= 2;
		--v4;
	} while (v4);
	return result;
}

//----- (00453710) --------------------------------------------------------
int sub_453710() {
	int v0; // edi
	int v1; // esi
	int v2; // ebx
	int v3; // eax

	v0 = -1;
	v1 = 1;
	v2 = 26;
	do {
		v3 = sub_415CD0((char*)v1);
		if (v3 && !nox_xxx_getUnitDefDd10_4E3BA0(v3)) {
			v0 &= ~v1;
		}
		v1 *= 2;
		--v2;
	} while (v2);
	return v0;
}

//----- (00453750) --------------------------------------------------------
char sub_453750() {
	int v0;      // esi
	int i;       // ebx
	uint32_t* j; // edi
	bool v3;     // zf
	int v4;      // eax

	v0 = sub_4A4800(*(uint32_t*)(dword_5d4594_1045464 + 32));
	for (i = 1520; i <= 1533; ++i) {
		for (j = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045468, i); (1 << v0) & 0x33; ++v0) {
			;
		}
		if (v0 >= *getMemIntPtr(0x5D4594, 1045472 + 4 * dword_5d4594_1045460)) {
			LOBYTE(v4) = nox_window_set_hidden((int)j, 1);
		} else {
			nox_window_set_hidden((int)j, 0);
			v3 = !sub_453080(v0);
			v4 = j[9];
			if (v3) {
				LOBYTE(v4) = v4 & 0xFB;
			} else {
				LOBYTE(v4) = v4 | 4;
			}
			j[9] = v4;
		}
		++v0;
	}
	return v4;
}

//----- (00453B00) --------------------------------------------------------
uint32_t* sub_453B00() {
	int v0;           // esi
	int v1;           // edi
	int v2;           // ebx
	uint32_t* result; // eax

	v0 = nox_xxx_guiFontHeightMB_43F320(*(uint32_t*)(dword_5d4594_1045480 + 236)) + 1;
	sub_46AB20(*(uint32_t**)&dword_5d4594_1045480, *(uint32_t*)(dword_5d4594_1045480 + 8), 15 * v0 + 2);
	sub_46AB20(*(uint32_t**)&dword_5d4594_1045508, *(uint32_t*)(dword_5d4594_1045508 + 8), 15 * v0 + 2);
	v1 = 1120;
	v2 = *(uint32_t*)(dword_5d4594_1045480 + 20) + v0 + 2;
	do {
		result = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045484, v1);
		result[5] = v2;
		v2 += v0;
		++v1;
	} while (v1 < 1134);
	return result;
}

//----- (00453B80) --------------------------------------------------------
int sub_453B80(int a1, int a2) {
	int result; // eax
	int xLeft;  // [esp+4h] [ebp-8h]
	int yTop;   // [esp+8h] [ebp-4h]

	nox_client_wndGetPosition_46AA60((uint32_t*)a1, &xLeft, &yTop);
	if ((signed char)*(uint8_t*)(a1 + 4) >= 0) {
		if (*(uint32_t*)(a2 + 20) != 0x80000000) {
			nox_client_drawRectFilledAlpha_49CF10(xLeft, yTop, *(uint32_t*)(a1 + 8), *(uint32_t*)(a1 + 12));
		}
		result = 1;
	} else {
		nox_client_drawImageAt_47D2C0(*(uint32_t*)(a2 + 24), xLeft, yTop);
		result = 1;
	}
	return result;
}

//----- (00453F70) --------------------------------------------------------
void sub_453F70(const void* a1) { memcpy(getMemAt(0x5D4594, 1045488), a1, 0x14u); }

//----- (00453F90) --------------------------------------------------------
char* sub_453F90() { return (char*)getMemAt(0x5D4594, 1045488); }

//----- (00453FA0) --------------------------------------------------------
int sub_453FA0(int a1, int a2, int a3) {
	bool v3;          // sf
	char v4;          // cl
	int result;       // eax
	unsigned char v6; // [esp+8h] [ebp+8h]

	v4 = a2 & 0x1F;
	v3 = (a2 & 0x8000001F & 0x80000000) != 0;
	v6 = a2 / 32;
	if (v3) {
		v4 = ((v4 - 1) | 0xE0) + 1;
	}
	result = 1 << v4;
	if (a3) {
		*(uint32_t*)(a1 + 4 * v6) |= result;
	} else {
		result = ~result;
		*(uint32_t*)(a1 + 4 * v6) &= result;
	}
	return result;
}

//----- (00454000) --------------------------------------------------------
int sub_454000(int a1, int a2) { return (*(uint32_t*)(a1 + 4 * ((a2 / 32) & 0xFF)) & (1 << (a2 % 32))) != 0; }

//----- (00454040) --------------------------------------------------------
int sub_454040(uint32_t* a1) {
	int v1;           // esi
	int v2;           // edi
	int result;       // eax
	unsigned char v4; // [esp+Ch] [ebp-4h]

	v1 = 1;
	*a1 = -1;
	v4 = 0;
	v2 = 1;
	a1[1] = -1;
	a1[2] = -1;
	a1[3] = -1;
	a1[4] = -1;
	do {
		if (v1 == 0x80000000) {
			v1 = 1;
			++v4;
		} else {
			v1 *= 2;
		}
		result = nox_xxx_spellIsValid_424B50(v2);
		if (result) {
			result = nox_xxx_spellFlags_424A70(v2);
			if (result & 0x7000000) {
				result = nox_xxx_spellIsEnabled_424B70(v2);
				if (!result) {
					result = (int)&a1[v4];
					*(uint32_t*)result = ~v1 & a1[v4];
				}
			}
		}
		++v2;
	} while (v2 < 137);
	return result;
}

//----- (004540E0) --------------------------------------------------------
int sub_4540E0(int a1) {
	int v1;     // esi
	int result; // eax

	v1 = 1;
	do {
		if (sub_454000(a1, v1)) {
			result = nox_xxx_spellEnable_424B90(v1);
		} else {
			result = nox_xxx_spellDisable_424BB0(v1);
		}
		++v1;
	} while (v1 < 137);
	return result;
}

//----- (00454120) --------------------------------------------------------
char sub_454120() {
	int v0;       // ebp
	int v1;       // ebx
	uint32_t* v2; // esi
	int v3;       // edi
	int v4;       // ecx
	bool v5;      // zf
	int v6;       // eax
	int v8;       // [esp+10h] [ebp-4h]

	v8 = *(uint32_t*)(dword_5d4594_1045480 + 32);
	v0 = 1120;
	v1 = 524 * sub_4A4800(v8);
	do {
		v2 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045484, v0);
		v3 = 0;
		v4 = *(uint32_t*)(v8 + 24);
		if (v4 + v1 != -4 && *(uint16_t*)(v4 + v1 + 4)) {
			v3 = nox_xxx_spellByTitle_424960((wchar2_t*)(v4 + v1 + 4));
			nox_window_set_hidden((int)v2, 0);
		} else {
			nox_window_set_hidden((int)v2, 1);
		}
		v1 += 524;
		if (v3) {
			v5 = !sub_454000((int)getMemAt(0x5D4594, 1045488), v3);
			v6 = v2[9];
			if (!v5) {
				LOBYTE(v6) = v6 | 4;
				goto LABEL_11;
			}
		} else {
			v6 = v2[9];
		}
		LOBYTE(v6) = v6 & 0xFB;
	LABEL_11:
		++v0;
		v2[9] = v6;
	} while (v0 <= 1133);
	return v6;
}

//----- (004541D0) --------------------------------------------------------
int nox_xxx_guiServerAccessLoad_4541D0(int a1) {
	int v2;        // esi
	uint32_t* v3;  // ebx
	char* v4;      // ebp
	char* v5;      // esi
	uint32_t* v6;  // edi
	uint32_t* v7;  // ebx
	uint32_t* v8;  // edi
	uint32_t* v9;  // ebx
	uint32_t* v10; // edi
	uint32_t* v11; // eax
	uint32_t* v12; // [esp+0h] [ebp-4h]
	uint32_t* v13; // [esp+0h] [ebp-4h]
	int v14;       // [esp+0h] [ebp-4h]
	uint32_t* v15; // [esp+8h] [ebp+4h]
	uint32_t* v16; // [esp+8h] [ebp+4h]
	uint32_t* v17; // [esp+8h] [ebp+4h]

	if (dword_5d4594_1045516) {
		return 0;
	}
	v2 = nox_strman_get_lang_code();
	if (nox_xxx_guiFontHeightMB_43F320(0) > 10) {
		v2 = 2;
	}
	dword_5d4594_1045516 =
		nox_new_window_from_file(*(const char**)getMemAt(0x587000, 127824 + 4 * v2), nox_xxx_windowAccessProc_454BA0);
	nox_draw_setTabWidth_43FE20(100);
	sub_46B120(*(uint32_t**)&dword_5d4594_1045516, a1);
	dword_5d4594_1045520 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10102);
	*getMemU32Ptr(0x5D4594, 1045524) = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10103);
	dword_5d4594_1045532 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10109);
	dword_5d4594_1045528 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10105);
	dword_5d4594_1045536 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10200);
	dword_5d4594_1045540 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10111);
	dword_5d4594_1045544 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10112);
	dword_5d4594_1045548 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10113);
	dword_5d4594_1045556 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10104);
	*getMemU32Ptr(0x5D4594, 1045560) = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10125);
	*getMemU32Ptr(0x5D4594, 1045564) = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10127);
	*getMemU32Ptr(0x5D4594, 1045568) = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10129);
	*getMemU32Ptr(0x5D4594, 1045572) = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10131);
	dword_5d4594_1045576 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10126);
	dword_5d4594_1045580 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10128);
	dword_5d4594_1045584 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10130);
	dword_5d4594_1045588 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10132);
	dword_5d4594_1045552 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10123);
	*getMemU32Ptr(0x5D4594, 1045592) = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10133);
	v3 = *(uint32_t**)(dword_5d4594_1045532 + 32);
	v4 = nox_xxx_gLoadImg_42F970("UISlider");
	v5 = nox_xxx_gLoadImg_42F970("UISliderLit");
	v6 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10190);
	v15 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10188);
	v12 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10189);
	*(uint32_t*)(v6[100] + 8) = 16;
	*(uint32_t*)(v6[100] + 12) = 10;
	sub_4B5700((int)v6, 0, 0, (int)v4, (int)v5, (int)v5);
	nox_xxx_wnd_46B280((int)v6, *(int*)&dword_5d4594_1045532);
	nox_xxx_wnd_46B280((int)v15, *(int*)&dword_5d4594_1045532);
	nox_xxx_wnd_46B280((int)v12, *(int*)&dword_5d4594_1045532);
	v3[9] = v6;
	v3[7] = v15;
	v3[8] = v12;
	v7 = *(uint32_t**)(dword_5d4594_1045528 + 32);
	v8 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10187);
	v16 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10185);
	v13 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10186);
	*(uint32_t*)(v8[100] + 8) = 16;
	*(uint32_t*)(v8[100] + 12) = 10;
	sub_4B5700((int)v8, 0, 0, (int)v4, (int)v5, (int)v5);
	nox_xxx_wnd_46B280((int)v8, *(int*)&dword_5d4594_1045528);
	nox_xxx_wnd_46B280((int)v16, *(int*)&dword_5d4594_1045528);
	nox_xxx_wnd_46B280((int)v13, *(int*)&dword_5d4594_1045528);
	v7[9] = v8;
	v7[7] = v16;
	v7[8] = v13;
	v9 = *(uint32_t**)(dword_5d4594_1045536 + 32);
	v10 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10203);
	v17 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10201);
	v11 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10202);
	*(uint32_t*)(v10[100] + 8) = 16;
	v14 = (int)v11;
	*(uint32_t*)(v10[100] + 12) = 10;
	sub_4B5700((int)v10, 0, 0, (int)v4, (int)v5, (int)v5);
	nox_xxx_wnd_46B280((int)v10, *(int*)&dword_5d4594_1045536);
	nox_xxx_wnd_46B280((int)v17, *(int*)&dword_5d4594_1045536);
	nox_xxx_wnd_46B280(v14, *(int*)&dword_5d4594_1045536);
	v9[9] = v10;
	v9[7] = v17;
	v9[8] = v14;
	sub_454740();
	sub_454640();
	if (nox_common_gameFlags_check_40A5C0(1)) {
		nox_xxx_wndSetDrawFn_46B340(*(int*)&dword_5d4594_1045516, sub_454A90);
	}
	nox_xxx_wndRetNULL_46A8A0();
	return dword_5d4594_1045516;
}

//----- (00454A90) --------------------------------------------------------
int sub_454A90(int a1, int a2) {
	uint16_t* v2; // eax
	int v3;       // eax
	char v4;      // cl
	int xLeft;    // [esp+8h] [ebp-8h]
	int yTop;     // [esp+Ch] [ebp-4h]

	nox_client_wndGetPosition_46AA60((uint32_t*)a1, &xLeft, &yTop);
	if ((signed char)*(uint8_t*)(a1 + 4) >= 0) {
		if (*(uint32_t*)(a2 + 20) != 0x80000000) {
			nox_client_drawRectFilledAlpha_49CF10(xLeft, yTop, *(uint32_t*)(a1 + 8), *(uint32_t*)(a1 + 12));
		}
	} else {
		nox_client_drawImageAt_47D2C0(*(uint32_t*)(a2 + 24), xLeft, yTop);
	}
	if (*(uint8_t*)(dword_5d4594_1045540 + 4) & 8) {
		v2 = (uint16_t*)nox_window_call_field_94(*(int*)&dword_5d4594_1045540, 16413, 0, 0);
		if (v2 && *v2) {
			if (!(*(uint8_t*)(dword_5d4594_1045544 + 4) & 8)) {
				nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1045544, 1);
			}
		} else if (*(uint8_t*)(dword_5d4594_1045544 + 4) & 8) {
			nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1045544, 0);
		}
	}
	v3 = nox_window_call_field_94(*(int*)&dword_5d4594_1045596, 16404, 0, 0);
	v4 = *(uint8_t*)(dword_5d4594_1045548 + 4);
	if (v3 < 0) {
		if (v4 & 8) {
			nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1045548, 0);
		}
	} else if (!(v4 & 8)) {
		nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1045548, 1);
		return 1;
	}
	return 1;
}

//----- (00454BA0) --------------------------------------------------------
int nox_xxx_windowAccessProc_454BA0(int a1, int a2, int* a3, int a4) {
	int v4;                 // eax
	int result;             // eax
	char* v6;               // eax
	int v7;                 // esi
	char* v8;               // esi
	const wchar2_t* v9;      // eax
	char* v10;              // esi
	const wchar2_t* v11;     // eax
	char* v12;              // esi
	const wchar2_t* v13;     // eax
	char* v14;              // esi
	const wchar2_t* v15;     // eax
	char* v16;              // eax
	char v17;               // dl
	uint32_t* v18;          // esi
	uint32_t* v19;          // eax
	uint32_t* v20;          // eax
	char* v21;              // esi
	wchar2_t* v22;           // esi
	int v23;                // esi
	int v24;                // esi
	char* v25;              // eax
	int* v26;               // eax
	int* v27;               // ebp
	int v28;                // eax
	wchar2_t* v29;           // eax
	wchar2_t* v30;           // edi
	char* v31;              // eax
	char* v32;              // esi
	int* v33;               // eax
	int* v34;               // edi
	int v35;                // eax
	wchar2_t* v36;           // eax
	char* v37;              // eax
	char* v38;              // esi
	char* v39;              // ebp
	const wchar2_t* v40;     // eax
	wchar2_t* v41;           // esi
	int v42;                // ebx
	int v43;                // ebp
	uint32_t* v44;          // esi
	const wchar2_t* v45;     // eax
	wchar2_t* v46;           // edi
	int v47;                // ebx
	wchar2_t WideCharStr[8]; // [esp+10h] [ebp-10h]
	char* v49;              // [esp+28h] [ebp+8h]
	char* v50;              // [esp+30h] [ebp+10h]
	char* v51;              // [esp+30h] [ebp+10h]

	switch (a2) {
	case 16387:
		v43 = a4;
		v44 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, a4);
		v51 = sub_416630();
		v49 = sub_416640();
		if (!v44 || (unsigned short)a3 == 1) {
			return 0;
		}
		v45 = (const wchar2_t*)nox_window_call_field_94((int)v44, 16413, 0, 0);
		v46 = (wchar2_t*)v45;
		if (!v45 || !*v45 || (v47 = nox_wcstol(v45, 0, 10), v47 < 0)) {
			v47 = 0;
		}
		switch (v43) {
		case 10104:
			nox_wcscpy((wchar2_t*)v49 + 39, v46);
			return 0;
		case 10126:
			if (v47 > 14) {
				LOBYTE(v47) = 14;
				nox_itow(14, WideCharStr, 10);
				nox_window_call_field_94((int)v44, 16414, (int)WideCharStr, -1);
			}
			v51[1] = v47 | v51[1] & 0xF0;
			result = 0;
			break;
		case 10128:
			if (v47 > 14) {
				LOBYTE(v47) = 14;
				nox_itow(14, WideCharStr, 10);
				nox_window_call_field_94((int)v44, 16414, (int)WideCharStr, -1);
			}
			v51[1] = (16 * v47) | v51[1] & 0xF;
			result = 0;
			break;
		case 10130:
			*(uint16_t*)(v51 + 5) = v47;
			result = 0;
			break;
		case 10132:
			*(uint16_t*)(v51 + 7) = v47;
			result = 0;
			break;
		case 10133:
			if (v47 <= 32) {
				if (v47 < 1) {
					v47 = 1;
					nox_itow(1, WideCharStr, 10);
					nox_window_call_field_94((int)v44, 16414, (int)WideCharStr, -1);
				}
			} else {
				v47 = 32;
				nox_itow(32, WideCharStr, 10);
				nox_window_call_field_94((int)v44, 16414, (int)WideCharStr, -1);
			}
			nox_xxx_servSetPlrLimit_409F80(v47);
			result = 0;
			v51[4] = v47;
			break;
		case 10136:
			nox_xxx_sysopSetPass_40A610(v46);
			result = 0;
			break;
		default:
			return 0;
		}
		break;
	case 16391:
		v7 = nox_xxx_wndGetID_46B0A0(a3);
		nox_xxx_clientPlaySoundSpecial_452D80(766, 100);
		switch (v7) {
		case 10102:
			v16 = sub_416630();
			v17 = *v16 ^ 0x10;
			*v16 = v17;
			if (v17 & 0x10) {
				v20 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10206);
				nox_xxx_wnd_46ABB0((int)v20, 1);
			} else {
				v18 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10206);
				nox_xxx_wnd_46ABB0((int)v18, 0);
				v18[9] &= 0xFFFFFFFB;
				v19 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10207);
				v19[9] |= 4u;
				nox_window_set_hidden(*(int*)&dword_5d4594_1045532, 1);
				sub_46ACE0(*(uint32_t**)&dword_5d4594_1045516, 10188, 10190, 1);
				nox_window_set_hidden(*(int*)&dword_5d4594_1045528, 0);
				sub_46ACE0(*(uint32_t**)&dword_5d4594_1045516, 10185, 10187, 0);
			}
			return 0;
		case 10103:
			v21 = sub_416630();
			nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1045556,
							   ((unsigned int)~*(uint32_t*)(*getMemU32Ptr(0x5D4594, 1045524) + 36) >> 2) & 1);
			*v21 ^= 0x20u;
			return 0;
		case 10112:
			v22 = (wchar2_t*)nox_window_call_field_94(*(int*)&dword_5d4594_1045540, 16413, 0, 0);
			if (wndIsShown_nox_xxx_wndIsShown_46ACC0(*(int*)&dword_5d4594_1045528)) {
				sub_4168A0(v22);
			} else {
				sub_416770(0, v22, 0);
			}
			nox_window_call_field_94(*(int*)&dword_5d4594_1045540, 16414, (int)getMemAt(0x5D4594, 1045600), 0);
			return 0;
		case 10113:
			if (*(uint8_t*)(dword_5d4594_1045520 + 36) & 4) {
				v23 = nox_window_call_field_94(*(int*)&dword_5d4594_1045532, 16404, 0, 0);
				nox_window_call_field_94(*(int*)&dword_5d4594_1045532, 16398, v23, 0);
				sub_416860(v23);
			} else {
				v24 = nox_window_call_field_94(*(int*)&dword_5d4594_1045528, 16404, 0, 0);
				nox_window_call_field_94(*(int*)&dword_5d4594_1045528, 16398, v24, 0);
				sub_416820(v24);
			}
			return 0;
		case 10124:
			v25 = sub_416630();
			v25[2] ^= 0x80u;
			return 0;
		case 10125:
			v8 = sub_416630();
			if (*(uint8_t*)(*getMemU32Ptr(0x5D4594, 1045560) + 36) & 4) {
				nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1045576, 0);
				v8[1] |= 0xFu;
				return 0;
			}
			nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1045576, 1);
			v9 = (const wchar2_t*)nox_window_call_field_94(*(int*)&dword_5d4594_1045576, 16413, 0, 0);
			if (!*v9) {
				return 0;
			}
			v8[1] = v8[1] & 0xF0 | nox_wcstol(v9, 0, 10);
			return 0;
		case 10127:
			v10 = sub_416630();
			if (*(uint8_t*)(*getMemU32Ptr(0x5D4594, 1045564) + 36) & 4) {
				nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1045580, 0);
				v10[1] |= 0xF0u;
				return 0;
			}
			nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1045580, 1);
			v11 = (const wchar2_t*)nox_window_call_field_94(*(int*)&dword_5d4594_1045580, 16413, 0, 0);
			if (!*v11) {
				return 0;
			}
			v10[1] = v10[1] & 0xF | nox_wcstol(v11, 0, 10);
			return 0;
		case 10129:
			v12 = sub_416630();
			if (*(uint8_t*)(*getMemU32Ptr(0x5D4594, 1045568) + 36) & 4) {
				nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1045584, 0);
				*(uint16_t*)(v12 + 5) = -1;
				return 0;
			}
			nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1045584, 1);
			v13 = (const wchar2_t*)nox_window_call_field_94(*(int*)&dword_5d4594_1045584, 16413, 0, 0);
			if (!*v13) {
				return 0;
			}
			*(uint16_t*)(v12 + 5) = nox_wcstol(v13, 0, 10);
			return 0;
		case 10131:
			v14 = sub_416630();
			if (*(uint8_t*)(*getMemU32Ptr(0x5D4594, 1045572) + 36) & 4) {
				nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1045588, 0);
				*(uint16_t*)(v14 + 7) = -1;
				return 0;
			}
			nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1045588, 1);
			v15 = (const wchar2_t*)nox_window_call_field_94(*(int*)&dword_5d4594_1045588, 16413, 0, 0);
			if (!*v15) {
				return 0;
			}
			*(uint16_t*)(v14 + 7) = nox_wcstol(v15, 0, 10);
			result = 0;
			break;
		case 10191:
			v33 = (int*)nox_window_call_field_94(*(int*)&dword_5d4594_1045536, 16404, 0, 0);
			v34 = v33;
			v35 = *v33;
			if (v35 < 0) {
				return 0;
			}
			do {
				v36 = (wchar2_t*)nox_window_call_field_94(*(int*)&dword_5d4594_1045536, 16406, v35, 0);
				v37 = nox_xxx_playerByName_4170D0(v36);
				v38 = v37;
				if (v37 && v37[2064] != 31) {
					if (nox_common_gameFlags_check_40A5C0(4096)) {
						sub_4DCFB0(*((uint32_t*)v38 + 514));
					} else {
						nox_xxx_playerCallDisconnect_4DEAB0((unsigned char)v38[2064], 4);
					}
				}
				v35 = v34[1];
				++v34;
			} while (v35 >= 0);
			return 0;
		case 10192:
			v26 = (int*)nox_window_call_field_94(*(int*)&dword_5d4594_1045536, 16404, 0, 0);
			v27 = v26;
			v28 = *v26;
			if (v28 < 0) {
				return 0;
			}
			do {
				v29 = (wchar2_t*)nox_window_call_field_94(*(int*)&dword_5d4594_1045536, 16406, v28, 0);
				v30 = v29;
				v31 = nox_xxx_playerByName_4170D0(v29);
				v32 = v31;
				if (v31 && v31[2064] != 31) {
					if (nox_common_gameFlags_check_40A5C0(4096)) {
						sub_4DCFB0(*((uint32_t*)v32 + 514));
					} else {
						nox_xxx_playerDisconnByPlrID_4DEB00((unsigned char)v32[2064]);
					}
					sub_416770(0, v30, v32 + 2112);
				}
				v28 = v27[1];
				++v27;
			} while (v28 >= 0);
			return 0;
		case 10206:
			nox_window_set_hidden(*(int*)&dword_5d4594_1045532, 0);
			sub_46ACE0(*(uint32_t**)&dword_5d4594_1045516, 10188, 10190, 0);
			nox_window_set_hidden(*(int*)&dword_5d4594_1045528, 1);
			sub_46ACE0(*(uint32_t**)&dword_5d4594_1045516, 10185, 10187, 1);
			dword_5d4594_1045596 = dword_5d4594_1045532;
			return 0;
		case 10207:
			nox_window_set_hidden(*(int*)&dword_5d4594_1045532, 1);
			sub_46ACE0(*(uint32_t**)&dword_5d4594_1045516, 10188, 10190, 1);
			nox_window_set_hidden(*(int*)&dword_5d4594_1045528, 0);
			sub_46ACE0(*(uint32_t**)&dword_5d4594_1045516, 10185, 10187, 0);
			dword_5d4594_1045596 = dword_5d4594_1045528;
			return 0;
		default:
			return 0;
		}
		break;
	case 16400:
		v4 = nox_xxx_wndGetID_46B0A0(a3) - 10123;
		if (v4) {
			if (v4 != 77) {
				return 0;
			}
			if (sub_455770()) {
				sub_46AD20(*(uint32_t**)&dword_5d4594_1045516, 10191, 10192, 1);
			} else {
				sub_46AD20(*(uint32_t**)&dword_5d4594_1045516, 10191, 10192, 0);
			}
			result = 0;
		} else {
			v6 = sub_416640();
			v6[100] ^= 1 << a4;
			result = 0;
		}
		break;
	case 16415:
		v39 = sub_416630();
		v50 = sub_416640();
		v40 = (const wchar2_t*)nox_window_call_field_94((int)a3, 16413, 0, 0);
		v41 = (wchar2_t*)v40;
		if (!v40 || !*v40 || (v42 = nox_wcstol(v40, 0, 10), v42 < 0)) {
			v42 = 0;
		}
		switch (nox_xxx_wndGetID_46B0A0(a3)) {
		case 10104:
			nox_wcscpy((wchar2_t*)v50 + 39, v41);
			return 0;
		case 10126:
			if (v42 > 14) {
				LOBYTE(v42) = 14;
				nox_itow(14, WideCharStr, 10);
				nox_window_call_field_94((int)a3, 16414, (int)WideCharStr, -1);
			}
			result = 0;
			v39[1] = v42 | v39[1] & 0xF0;
			break;
		case 10128:
			if (v42 > 14) {
				LOBYTE(v42) = 14;
				nox_itow(14, WideCharStr, 10);
				nox_window_call_field_94((int)a3, 16414, (int)WideCharStr, -1);
			}
			result = 0;
			v39[1] = (16 * v42) | v39[1] & 0xF;
			break;
		case 10130:
			*(uint16_t*)(v39 + 5) = v42;
			result = 0;
			break;
		case 10132:
			*(uint16_t*)(v39 + 7) = v42;
			result = 0;
			break;
		case 10133:
			if (v42 <= 32) {
				if (v42 < 1) {
					v42 = 1;
					nox_itow(1, WideCharStr, 10);
					nox_window_call_field_94((int)a3, 16414, (int)WideCharStr, -1);
				}
			} else {
				v42 = 32;
				nox_itow(32, WideCharStr, 10);
				nox_window_call_field_94((int)a3, 16414, (int)WideCharStr, -1);
			}
			nox_xxx_servSetPlrLimit_409F80(v42);
			v39[4] = v42;
			result = 0;
			break;
		case 10136:
			nox_xxx_sysopSetPass_40A610(v41);
			result = 0;
			break;
		default:
			return 0;
		}
		break;
	default:
		return 0;
	}
	return result;
}

//----- (00455770) --------------------------------------------------------
int sub_455770() {
	int* v0;     // eax
	int* v1;     // esi
	int v2;      // eax
	wchar2_t* v3; // eax
	char* v4;    // eax

	v0 = (int*)nox_window_call_field_94(*(int*)&dword_5d4594_1045536, 16404, 0, 0);
	v1 = v0;
	v2 = *v0;
	if (v2 < 0) {
		return 0;
	}
	while (1) {
		v3 = (wchar2_t*)nox_window_call_field_94(*(int*)&dword_5d4594_1045536, 16406, v2, 0);
		v4 = nox_xxx_playerByName_4170D0(v3);
		if (v4) {
			if (v4[2064] != 31) {
				break;
			}
		}
		v2 = v1[1];
		++v1;
		if (v2 < 0) {
			return 0;
		}
	}
	return 1;
}

//----- (004557D0) --------------------------------------------------------
int sub_4557D0(int a1) {
	int result; // eax

	result = a1;
	if (a1) {
		result = dword_5d4594_1045516;
		if (dword_5d4594_1045516) {
			result = nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1045516);
		}
	}
	dword_5d4594_1045516 = 0;
	return result;
}

//----- (00455800) --------------------------------------------------------
int* sub_455800() {
	int* result;             // eax
	int v1;                  // eax
	int* i;                  // esi
	int* j;                  // esi
	wchar2_t WideCharStr[10]; // [esp+4h] [ebp-48h]
	wchar2_t v5[26];          // [esp+18h] [ebp-34h]

	result = *(int**)&dword_5d4594_1045516;
	if (dword_5d4594_1045516) {
		v1 = nox_xxx_servGetPlrLimit_409FA0();
		nox_itow(v1, WideCharStr, 10);
		nox_window_call_field_94(*getMemIntPtr(0x5D4594, 1045592), 16414, (int)WideCharStr, 0);
		result = (int*)nox_common_gameFlags_check_40A5C0(1);
		if (result) {
			nox_window_call_field_94(*(int*)&dword_5d4594_1045528, 16399, 0, 0);
			nox_window_call_field_94(*(int*)&dword_5d4594_1045532, 16399, 0, 0);
			for (i = sub_4168E0(); i; i = sub_4168F0(i)) {
				nox_window_call_field_94(*(int*)&dword_5d4594_1045532, 16397, (int)(i + 3), -1);
			}
			result = sub_416900();
			for (j = result; result; j = result) {
				if (*((uint8_t*)j + 72)) {
					nox_window_call_field_94(*(int*)&dword_5d4594_1045528, 16397, (int)(j + 3), -1);
				} else {
					nox_swprintf(v5, L"*%s", j + 3);
					nox_window_call_field_94(*(int*)&dword_5d4594_1045528, 16397, (int)v5, -1);
				}
				result = sub_416910(j);
			}
		}
	}
	return result;
}

//----- (00455920) --------------------------------------------------------
int sub_455920(wchar2_t* a1) {
	int result;   // eax
	uint32_t* v2; // eax

	result = dword_5d4594_1045516;
	if (dword_5d4594_1045516) {
		v2 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045516, 10200);
		result = nox_window_call_field_94((int)v2, 16397, a1, 3);
	}
	return result;
}

//----- (00455950) --------------------------------------------------------
void sub_455950(wchar2_t* a1) {
	int v1; // eax

	if (dword_5d4594_1045516) {
		v1 = sub_4559B0(a1);
		if (v1 != -1) {
			nox_window_call_field_94(*(int*)&dword_5d4594_1045536, 16398, v1, 0);
			if (!sub_455770()) {
				sub_46AD20(*(uint32_t**)&dword_5d4594_1045516, 10191, 10192, 0);
			}
		}
	}
}

//----- (004559B0) --------------------------------------------------------
int sub_4559B0(wchar2_t* a1) {
	int v1;            // esi
	int v2;            // ebx
	const wchar2_t* v3; // eax

	v1 = 0;
	v2 = *(uint32_t*)(dword_5d4594_1045536 + 32);
	if (*(short*)(v2 + 44) <= 0) {
		return -1;
	}
	while (1) {
		v3 = (const wchar2_t*)nox_window_call_field_94(*(int*)&dword_5d4594_1045536, 16406, v1, 0);
		if (!nox_wcscmp(v3, a1)) {
			break;
		}
		if (++v1 >= *(short*)(v2 + 44)) {
			return -1;
		}
	}
	return v1;
}

//----- (00455A00) --------------------------------------------------------
int sub_455A00(int a1) {
	int result; // eax

	if (a1 && *getMemU32Ptr(0x5D4594, 1045608) && nox_xxx_wndGetFlags_46ADA0(*(int*)&dword_5d4594_1045604) & 0x10) {
		result = nox_xxx_wndShowModalMB_46A8C0(*(int*)&dword_5d4594_1045604);
	} else {
		result = nox_window_set_hidden(*(int*)&dword_5d4594_1045604, 1);
	}
	return result;
}

//----- (00455A50) --------------------------------------------------------
char sub_455A50(char a1) {
	uint32_t* v1; // edi
	int v2;       // eax
	int v3;       // esi
	int v4;       // edx

	v1 = 0;
	nox_xxx_objGetTeamByNetCode_418C80(nox_player_netCode_85319C);
	if (!(dword_5d4594_1045604 || (v2 = sub_455C30()) != 0)) {
		return v2;
	}
	int cur_w, cur_h, cur_d;
	nox_xxx_gameGetScreenBoundaries_43BEB0_get_video_mode(&cur_w, &cur_h, &cur_d);

	int max_w, max_h;
	nox_xxx_screenGetSize_430C50_get_video_max(&max_w, &max_h);

	if (cur_w > max_w) {
		cur_w = max_w;
	}
	if (cur_h > max_h) {
		cur_h = max_h;
	}

	*getMemU32Ptr(0x5D4594, 1045612) = 0;
	*getMemU32Ptr(0x5D4594, 1045616) = 0;
	*getMemU32Ptr(0x5D4594, 1045620) = 0;
	*getMemU32Ptr(0x5D4594, 1045608) = 1;
	*getMemU32Ptr(0x5D4594, 1045624) = 0;
	sub_46ACE0(*(uint32_t**)&dword_5d4594_1045604, 8811, 8826, 1);
	v3 = 0;
	*getMemU8Ptr(0x5D4594, 1045628) = a1;
	if (a1) {
		do {
			v1 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045604, v3 + 8811);
			nox_window_set_hidden((int)v1, 0);
			++v3;
		} while (v3 < getMemByte(0x5D4594, 1045628));
	}
	if (getMemByte(0x5D4594, 1045628) <= 4u) {
		v4 = cur_w - *(uint32_t*)(*(int*)&dword_5d4594_1045604 + 8) / 2;
	} else {
		v4 = cur_w - *(uint32_t*)(dword_5d4594_1045604 + 8);
	}
	*(uint32_t*)(dword_5d4594_1045604 + 16) = v4 - 91;
	*(uint32_t*)(dword_5d4594_1045604 + 24) =
		*(uint32_t*)(dword_5d4594_1045604 + 8) + *(uint32_t*)(dword_5d4594_1045604 + 16);
	LOBYTE(v2) = getMemByte(0x5D4594, 1045628);
	if (getMemByte(0x5D4594, 1045628) <= 8u) {
		if (getMemByte(0x5D4594, 1045628) <= 4u) {
			if (!v1) {
				*getMemU32Ptr(0x5D4594, 1045608) = 0;
				return v2;
			}
			*(uint32_t*)(dword_5d4594_1045604 + 20) = cur_h - 40 * getMemByte(0x5D4594, 1045628);
		} else {
			*(uint32_t*)(dword_5d4594_1045604 + 20) = cur_h - *(int*)(dword_5d4594_1045604 + 12) / 2;
		}
	} else {
		*(uint32_t*)(dword_5d4594_1045604 + 20) = cur_h - *(uint32_t*)(dword_5d4594_1045604 + 12);
	}
	*(uint32_t*)(dword_5d4594_1045604 + 28) =
		*(uint32_t*)(dword_5d4594_1045604 + 12) + *(uint32_t*)(dword_5d4594_1045604 + 20);
	LOBYTE(v2) = nox_xxx_wndShowModalMB_46A8C0(*(int*)&dword_5d4594_1045604);
	return v2;
}

//----- (00455C10) --------------------------------------------------------
int sub_455C10() {
	*getMemU32Ptr(0x5D4594, 1045608) = 0;
	return sub_455A00(0);
}

//----- (00455CD0) --------------------------------------------------------
int sub_455CD0(uint8_t* a1, uint32_t* a2) {
	uint8_t* v2;      // esi
	int v3;           // eax
	uint8_t* v5;      // [esp-Ch] [ebp-18h]
	int v6;           // [esp+4h] [ebp-8h]
	int v7;           // [esp+8h] [ebp-4h]
	unsigned char v8; // [esp+10h] [ebp+4h]

	v2 = a1;
	v5 = a1;
	v8 = *a1 - 107;
	nox_client_wndGetPosition_46AA60(v5, &v6, &v7);
	if (getMemByte(0x5D4594, 1045612 + v8) == 1) {
		v3 = a2[8];
		if (v3) {
			nox_client_drawImageAt_47D2C0(v3, v6, v7);
			goto LABEL_9;
		}
	} else if (getMemByte(0x5D4594, 1045612 + v8) == 2) {
		v3 = a2[12];
		if (v3) {
			nox_client_drawImageAt_47D2C0(v3, v6, v7);
			goto LABEL_9;
		}
	} else {
		v3 = a2[6];
		if (v3) {
			nox_client_drawImageAt_47D2C0(v3, v6, v7);
			goto LABEL_9;
		}
	}
LABEL_9:
	if (v2[4] & 0x20 && *getMemU32Ptr(0x5D4594, 1045632)) {
		nox_client_drawImageAt_47D2C0(*getMemIntPtr(0x5D4594, 1045632), v6 - 4, v7 - 4);
	}
	return 1;
}

//----- (00455E70) --------------------------------------------------------
int sub_455E70(unsigned char a1) {
	uint32_t* v2; // eax
	uint32_t* v3; // eax

	for (int i = 0; i < getMemByte(0x5D4594, 1045628); i++) {
		v2 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045604, i + 8811);
		nox_xxx_wndClearFlag_46AD80((int)v2, 32);
	}
	v3 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045604, a1 + 8810);
	return nox_xxx_wnd_46AD60((int)v3, 32);
}

//----- (00455EE0) --------------------------------------------------------
int sub_455EE0() {
	int result; // eax

	result = dword_5d4594_1045604;
	if (dword_5d4594_1045604) {
		result = nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1045604);
	}
	dword_5d4594_1045604 = 0;
	*getMemU32Ptr(0x5D4594, 1045608) = 0;
	return result;
}

//----- (00455F10) --------------------------------------------------------
int sub_455F10(int a1) {
	int result; // eax

	if (a1 && dword_5d4594_1045640 && nox_xxx_wndGetFlags_46ADA0(*(int*)&dword_5d4594_1045636) & 0x10) {
		result = nox_xxx_wndShowModalMB_46A8C0(*(int*)&dword_5d4594_1045636);
	} else {
		result = nox_window_set_hidden(*(int*)&dword_5d4594_1045636, 1);
	}
	return result;
}

//----- (00455F60) --------------------------------------------------------
int sub_455F60() {
	int result; // eax
	int v1;     // ecx
	int v2;     // [esp+0h] [ebp-14h]
	int v3;     // [esp+4h] [ebp-10h]
	int v4;     // [esp+8h] [ebp-Ch]
	int v5;     // [esp+Ch] [ebp-8h]
	int v6;     // [esp+10h] [ebp-4h]

	nox_xxx_objGetTeamByNetCode_418C80(nox_player_netCode_85319C);
	if (dword_5d4594_1045636 || (result = sub_456070()) != 0) {
		nox_xxx_gameGetScreenBoundaries_43BEB0_get_video_mode(&v3, &v2, &v6);
		nox_xxx_screenGetSize_430C50_get_video_max(&v4, &v5);
		v1 = v3;
		if (v3 > v4) {
			v1 = v4;
			v3 = v4;
		}
		if (v2 > v5) {
			v2 = v5;
		}
		dword_5d4594_1045640 = 1;
		*(uint32_t*)(dword_5d4594_1045636 + 16) = v1 - *(uint32_t*)(dword_5d4594_1045636 + 8) / 3 - 91;
		*(uint32_t*)(dword_5d4594_1045636 + 24) =
			*(uint32_t*)(dword_5d4594_1045636 + 8) + *(uint32_t*)(dword_5d4594_1045636 + 16);
		*(uint32_t*)(dword_5d4594_1045636 + 20) = v2 - 120;
		*(uint32_t*)(dword_5d4594_1045636 + 28) =
			*(uint32_t*)(dword_5d4594_1045636 + 12) + *(uint32_t*)(dword_5d4594_1045636 + 20);
		*getMemU8Ptr(0x5D4594, 1045644) = 0;
		result = nox_xxx_wndShowModalMB_46A8C0(*(int*)&dword_5d4594_1045636);
	}
	return result;
}

//----- (00456050) --------------------------------------------------------
int sub_456050() {
	dword_5d4594_1045640 = 0;
	return sub_455F10(0);
}

//----- (00456070) --------------------------------------------------------
int sub_456070() {
	if (!dword_5d4594_1045636) {
		dword_5d4594_1045636 = nox_new_window_from_file("gui_fb.wnd", 0);
		nox_window_set_all_funcs(*(uint32_t**)&dword_5d4594_1045636, 0, sub_4560D0, 0);
		if (!dword_5d4594_1045636) {
			return 0;
		}
		sub_455F10(0);
		*getMemU32Ptr(0x5D4594, 1045648) = nox_xxx_gLoadImg_42F970("FlagTeamBorder");
	}
	return 1;
}

//----- (004560D0) --------------------------------------------------------
int sub_4560D0(int a1, int a2) {
	int v2; // eax
	int v4; // [esp+4h] [ebp-8h]
	int v5; // [esp+8h] [ebp-4h]

	nox_client_wndGetPosition_46AA60((uint32_t*)a1, &v4, &v5);
	v2 = *(uint32_t*)(a2 + 24);
	if (v2) {
		nox_client_drawImageAt_47D2C0(v2, v4, v5);
	}
	if (*(uint8_t*)(a1 + 4) & 0x20 && *getMemU32Ptr(0x5D4594, 1045648)) {
		nox_client_drawImageAt_47D2C0(*getMemIntPtr(0x5D4594, 1045648), v4 - 4, v5 - 4);
	}
	return 1;
}

//----- (00456240) --------------------------------------------------------
int sub_456240() {
	int result; // eax

	result = dword_5d4594_1045636;
	if (dword_5d4594_1045636) {
		result = nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1045636);
	}
	dword_5d4594_1045636 = 0;
	dword_5d4594_1045640 = 0;
	return result;
}

//----- (00456500) --------------------------------------------------------
int sub_456500() {
	uint32_t* v0; // ebp
	uint32_t* v1; // ebx
	char* i;      // esi
	char* j;      // esi
	uint32_t* v4; // edi
	int result;   // eax

	// Refresh owns both the widget rows and their parallel metadata lists.
	for (int list = 0; list < 2; ++list) {
		int* head = getMemIntPtr(0x5D4594, 1045652 + 16 * list);
		int* row = nox_common_list_getFirstSafe_425890(head);
		while (row) {
			int* next = nox_common_list_getNextSafe_4258A0(row);
			nox_common_list_remove_425920((uint32_t**)row);
			free(row);
			row = next;
		}
	}

	v0 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045684, 10501);
	v1 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045684, 10502);
	nox_window_call_field_94((int)v0, 16399, 0, 0);
	nox_window_call_field_94((int)v1, 16399, 0, 0);
	nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1045688, 0);
	nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1045692, 0);
	for (i = nox_server_teamFirst_418B10(); i; i = nox_server_teamNext_418B60((int)i)) {
		sub_457230((wchar2_t*)i);
	}
	for (j = nox_common_playerInfoGetFirst_416EA0(); j; j = nox_common_playerInfoGetNext_416EE0((int)j)) {
		if (j[2064] != 31 || !nox_common_getEngineFlag(NOX_ENGINE_FLAG_DISABLE_GRAPHICS_RENDERING)) {
			sub_457140(*((uint32_t*)j + 515), (wchar2_t*)j + 2352);
			v4 = nox_xxx_objGetTeamByNetCode_418C80(*((uint32_t*)j + 515));
			if (nox_xxx_servObjectHasTeam_419130((int)v4)) {
				sub_4571A0(*((uint32_t*)j + 515), *((unsigned char*)v4 + 4));
			}
		}
	}
	if (nox_common_gameFlags_check_40A5C0(4096)) {
		nox_xxx_wnd_46ABB0((int)v0, 0);
		result = nox_xxx_wnd_46ABB0((int)v1, 0);
	} else {
		nox_xxx_wnd_46ABB0((int)v0, 1);
		result = nox_xxx_wnd_46ABB0((int)v1, 1);
	}
	return result;
}

//----- (00456640) --------------------------------------------------------
int sub_456640(int a1, int a2) {
	uint32_t* v2; // eax
	uint32_t* v3; // eax
	uint32_t* v4; // esi
	uint32_t* v5; // eax
	uint32_t* v6; // eax
	uint32_t* v7; // eax
	int xLeft;    // [esp+4h] [ebp-8h]
	int yTop;     // [esp+8h] [ebp-4h]

	nox_xxx_objGetTeamByNetCode_418C80(nox_player_netCode_85319C);
	nox_client_wndGetPosition_46AA60((uint32_t*)a1, &xLeft, &yTop);
	if ((signed char)*(uint8_t*)(a1 + 4) >= 0) {
		if (*(uint32_t*)(a2 + 20) != 0x80000000) {
			nox_client_drawRectFilledAlpha_49CF10(xLeft, yTop, *(uint32_t*)(a1 + 8), *(uint32_t*)(a1 + 12));
		}
	} else {
		nox_client_drawImageAt_47D2C0(*(uint32_t*)(a2 + 24), xLeft, yTop);
	}
	if (nox_common_gameFlags_check_40A5C0(1) && !nox_common_gameFlags_check_40A5C0(0x8000) && !sub_40A740()) {
		v2 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045684, 10502);
		if (nox_window_call_field_94((int)v2, 16404, 0, 0) < 0) {
			v6 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045684, 10503);
			nox_xxx_wnd_46ABB0((int)v6, 0);
		} else {
			v3 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045684, 10501);
			v4 = (uint32_t*)nox_window_call_field_94((int)v3, 16404, 0, 0);
			v5 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045684, 10503);
			nox_xxx_wnd_46ABB0((int)v5, (int32_t)*v4 >= 0);
		}
	}
	v7 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045684, 10502);
	if (nox_window_call_field_94((int)v7, 16404, 0, 0) < 0) {
		nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1045688, 0);
		nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1045692, 0);
	}
	return 1;
}

//----- (00456BB0) --------------------------------------------------------
void sub_456BB0(int a1) {
	uint32_t* v1; // eax
	int v2;       // esi
	int v3;       // edi
	uint32_t* v4; // esi
	float2 v6;    // [esp+8h] [ebp-8h]

	if (!dword_8531A0_2576 || !(*(uint8_t*)(dword_8531A0_2576 + 4) & 1) &&
											  !(*(uint8_t*)(dword_8531A0_2576 + 3680) & 1)) {
		v2 = nox_common_gameFlags_check_40A5C0(1);
		v1 = nox_xxx_objGetTeamByNetCode_418C80(nox_player_netCode_85319C);
		v3 = (int)v1;
		if (v1) {
			if (nox_xxx_servObjectHasTeam_419130((int)v1)) {
				if (!v2) {
					sub_419960(a1, v3, nox_player_netCode_85319C);
					++*getMemU32Ptr(0x5D4594, 1045696);
					return;
				}
				sub_4196D0(v3, a1, nox_player_netCode_85319C, 1);
				v1 = (uint32_t*)nox_common_gameFlags_check_40A5C0(128);
				if (!v1) {
					v1 = nox_common_playerInfoGetByID_417040(nox_player_netCode_85319C);
					v4 = v1;
					if (v1) {
						nox_xxx_mapFindPlayerStart_4F7AB0(&v6, v1[514]);
						nox_xxx_unitMove_4E7010(v4[514], &v6);
						++*getMemU32Ptr(0x5D4594, 1045696);
						return;
					}
				}
			} else {
				if (v2) {
					nox_xxx_createAtImpl_4191D0(*(uint8_t*)(a1 + 57), v3, v2, nox_player_netCode_85319C, 1);
					++*getMemU32Ptr(0x5D4594, 1045696);
					return;
				}
				sub_419900(a1, v3, nox_player_netCode_85319C);
			}
			++*getMemU32Ptr(0x5D4594, 1045696);
		}
	}
}

//----- (00456D00) --------------------------------------------------------
wchar2_t* sub_456D00(int a1, wchar2_t* a2) {
	uint32_t* v2;      // eax
	const wchar2_t* v3; // eax
	wchar2_t* v4;       // eax
	wchar2_t v6[56];    // [esp+0h] [ebp-70h]

	v2 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045684, 10502);
	v3 = (const wchar2_t*)nox_window_call_field_94((int)v2, 16406, a1, 0);
	nox_wcscpy(v6, v3);
	v4 = nox_wcstok(v6, L"\t\n\r");
	return nox_wcscpy(a2, v4);
}

//----- (00456D60) --------------------------------------------------------
int* sub_456D60(int a1) {
	int* v1;     // esi
	int* v2;     // edi
	int* result; // eax
	int* v4;     // esi
	int* v5;     // edi

	if (a1 && dword_5d4594_1045684) {
		nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1045684);
	}
	v1 = nox_common_list_getFirstSafe_425890(getMemIntPtr(0x5D4594, 1045652));
	if (v1) {
		do {
			v2 = nox_common_list_getNextSafe_4258A0(v1);
			nox_common_list_remove_425920((uint32_t**)v1);
			free(v1);
			v1 = v2;
		} while (v2);
	}
	result = nox_common_list_getFirstSafe_425890(getMemIntPtr(0x5D4594, 1045668));
	v4 = result;
	if (result) {
		do {
			v5 = nox_common_list_getNextSafe_4258A0(v4);
			nox_common_list_remove_425920((uint32_t**)v4);
			free(v4);
			v4 = v5;
		} while (v5);
	}
	dword_5d4594_1045684 = 0;
	return result;
}

//----- (00456DF0) --------------------------------------------------------
int sub_456DF0(int a1) {
	int result;   // eax
	int v2;       // esi
	uint32_t* v3; // eax

	result = dword_5d4594_1045684;
	if (dword_5d4594_1045684) {
		result = sub_456E40(a1, 1);
		v2 = result;
		if (result != -1) {
			v3 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045684, 10501);
			result = nox_window_call_field_94((int)v3, 16398, v2, 0);
		}
	}
	return result;
}

//----- (00456E40) --------------------------------------------------------
int sub_456E40(int a1, int a2) {
	int* v2; // esi
	int v3;  // edi

	v2 = nox_common_list_getFirstSafe_425890(getMemIntPtr(0x5D4594, 1045652));
	v3 = 0;
	if (!v2) {
		return -1;
	}
	while (v2[15] != a1) {
		v2 = nox_common_list_getNextSafe_4258A0(v2);
		++v3;
		if (!v2) {
			return -1;
		}
	}
	if (a2) {
		nox_common_list_remove_425920((uint32_t**)v2);
		free(v2);
	}
	return v3;
}

//----- (00456EA0) --------------------------------------------------------
int sub_456EA0(wchar2_t* a1) {
	int result;   // eax
	uint32_t* v2; // edi
	int v3;       // eax
	int v4;       // esi

	result = dword_5d4594_1045684;
	if (dword_5d4594_1045684) {
		v2 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045684, 10502);
		v3 = sub_456F10(a1, 1);
		v4 = v3;
		if (v3 != -1) {
			sub_4258C0((uint32_t**)getMemAt(0x5D4594, 1045668), v3);
			nox_window_call_field_94((int)v2, 16404, 0, 0);
			nox_window_call_field_94((int)v2, 16398, v4, 0);
		}
		result = sub_456500();
	}
	return result;
}

//----- (00456F10) --------------------------------------------------------
int sub_456F10(wchar2_t* a1, int a2) {
	int* v2;        // esi
	int v3;         // ebx

	v2 = nox_common_list_getFirstSafe_425890(getMemIntPtr(0x5D4594, 1045668));
	v3 = 0;
	if (!v2) {
		return -1;
	}
	while (1) {
		// Metadata stores the bare team name, including any spaces.
		if (!_nox_wcsicmp((const wchar2_t*)v2 + 6, a1)) {
			break;
		}
		v2 = nox_common_list_getNextSafe_4258A0(v2);
		++v3;
		if (!v2) {
			return -1;
		}
	}
	if (a2) {
		nox_common_list_remove_425920((uint32_t**)v2);
		free(v2);
	}
	return v3;
}

//----- (00456FA0) --------------------------------------------------------
int sub_456FA0() {
	int result; // eax
	int* v1;    // esi
	int* v2;    // edi

	result = dword_5d4594_1045684;
	if (dword_5d4594_1045684) {
		v1 = nox_common_list_getFirstSafe_425890(getMemIntPtr(0x5D4594, 1045668));
		if (v1) {
			do {
				v2 = nox_common_list_getNextSafe_4258A0(v1);
				nox_common_list_remove_425920((uint32_t**)v1);
				free(v1);
				v1 = v2;
			} while (v2);
		}
		nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1045688, 0);
		nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1045692, 0);
		result = sub_456500();
	}
	return result;
}

//----- (00457120) --------------------------------------------------------
unsigned char sub_457120(int a1) { return getMemByte(0x587000, 128968 + 8 * (*(unsigned char*)(a1 + 56) % 10)); }

//----- (00457140) --------------------------------------------------------
int sub_457140(int a1, wchar2_t* a2) {
	int result;   // eax
	wchar2_t* v3;  // esi
	uint32_t* v4; // eax

	result = dword_5d4594_1045684;
	if (dword_5d4594_1045684) {
		v3 = (wchar2_t*)calloc(1, 0x48u);
		sub_425770(v3);
		nox_common_list_append_4258E0((int)getMemAt(0x5D4594, 1045652), v3);
		nox_wcscpy(v3 + 6, a2);
		*((uint32_t*)v3 + 15) = a1;
		v4 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045684, 10501);
		result = nox_window_call_field_94((int)v4, 16397, (int)a2, 3);
	}
	return result;
}

//----- (004571A0) --------------------------------------------------------
int sub_4571A0(int a1, int a2) {
	int result; // eax
	int v3;     // ebp
	int v4;     // esi
	int v5;     // edi
	int* i;     // eax

	result = dword_5d4594_1045684;
	if (dword_5d4594_1045684) {
		v3 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1045684, 10501)->widget_data;
		result = sub_456E40(a1, 0);
		v4 = result;
		if (result != -1) {
			v5 = 3;
			if (a2) {
				for (i = nox_common_list_getFirstSafe_425890(getMemIntPtr(0x5D4594, 1045668)); i;
					 i = nox_common_list_getNextSafe_4258A0(i)) {
					if (i[15] == a2) {
						v5 = *((unsigned char*)i + 64);
					}
				}
			}
			result = *(uint32_t*)(v3 + 24);
			*(uint32_t*)(result + 524 * v4 + 516) = **(uint32_t**)getMemAt(0x85B3FC, 132 + 4 * v5);
		}
	}
	return result;
}

//----- (004573A0) --------------------------------------------------------
int sub_4573A0() {
	int result; // eax

	result = dword_5d4594_1045684;
	if (dword_5d4594_1045684) {
		result = sub_456500();
	}
	return result;
}

//----- (004573B0) --------------------------------------------------------
void sub_4573B0() { *getMemU32Ptr(0x5D4594, 1045696) = 0; }

//----- (00457460) --------------------------------------------------------
int sub_457460(int a1) {
	wchar2_t WideCharStr[16]; // [esp+4h] [ebp-20h]

	nox_itow(*(unsigned short*)(a1 + 54), WideCharStr, 10);
	nox_window_call_field_94(*(int*)&dword_5d4594_1046516, 16414, (int)WideCharStr, 0);
	nox_itow(*(unsigned char*)(a1 + 56), WideCharStr, 10);
	return nox_window_call_field_94(*(int*)&dword_5d4594_1046520, 16414, (int)WideCharStr, 0);
}

//----- (004574D0) --------------------------------------------------------
int nox_xxx_guiServerOptionsTryHide_4574D0() {
	if (!dword_5d4594_1046492) {
		return 0;
	}
	nox_xxx_guiServerOptionsHide_4597E0(0);
	dword_5d4594_1046492 = 0;
	return 1;
}

//----- (00457B60) --------------------------------------------------------
int sub_457B60(int a1) {
	char* v1;                // eax
	uint32_t* v2;            // eax
	char* v3;                // eax
	unsigned short v4;       // di
	char* v5;                // esi
	uint32_t* v6;            // eax
	unsigned short v7;       // ax
	unsigned char v8;        // al
	uint32_t* v9;            // eax
	uint32_t* v10;           // eax
	uint32_t* v11;           // edi
	int v12;                 // eax
	uint32_t* v13;           // edi
	int v14;                 // eax
	uint32_t* v15;           // edi
	int v16;                 // eax
	unsigned int v17;        // ecx
	wchar2_t WideCharStr[16]; // [esp+0h] [ebp-E8h]
	wchar2_t v20[100];        // [esp+20h] [ebp-C8h]

	v1 = nox_xxx_serverOptionsGetServername_40A4C0();
	nox_swprintf(v20, L"%S", v1);
	nox_window_call_field_94(dword_5d4594_1046512, 16414, (int)v20, 0);
	sub_459CD0();
	if (nox_common_gameFlags_check_40A5C0(128)) {
		v2 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1046492, 10122);
	} else {
		v2 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1046492, 10119);
	}
	nox_window_set_hidden((int)v2, 1);
	if (nox_common_gameFlags_check_40A5C0(1)) {
		nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1046508, 1);
		if (nox_common_gameFlags_check_40A5C0(128)) {
			v3 = nox_xxx_cliGamedataGet_416590(1);
			v4 = *((uint16_t*)v3 + 26) & 0x17F0;
			nox_client_guiserv_updateMapList_458230(v4, v3, 1);
			v5 = sub_4165B0();
			*((uint16_t*)v5 + 26) = v4 | *((uint16_t*)v5 + 26) & 0xE80F;
			if (sub_40A740()) {
				nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1046508, 0);
			} else if (nox_xxx_CheckGameplayFlags_417DA0(4)) {
				v6 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1046508, 10330);
				v6[9] |= 4u;
			} else {
				sub_46AD20(*(uint32_t**)&dword_5d4594_1046508, 10331, 10333, 0);
			}
		} else {
			if (nox_xxx_CheckGameplayFlags_417DA0(4)) {
				sub_46AD20(*(uint32_t**)&dword_5d4594_1046508, 10330, 10331, 0);
			} else {
				sub_46AD20(*(uint32_t**)&dword_5d4594_1046508, 10330, 10333, 0);
			}
			v5 = (char*)a1;
			nox_client_guiserv_updateMapList_458230(*(unsigned short*)(a1 + 52), (char*)a1, 0);
		}
		v7 = nox_xxx_servGamedataGet_40A020(*((uint16_t*)v5 + 26));
		nox_itow(v7, WideCharStr, 10);
		nox_window_call_field_94(*(int*)&dword_5d4594_1046516, 16414, (int)WideCharStr, 0);
		v8 = sub_40A180(*((uint16_t*)v5 + 26));
		nox_itow(v8, WideCharStr, 10);
		nox_window_call_field_94(*(int*)&dword_5d4594_1046520, 16414, (int)WideCharStr, 0);
		sub_4580E0((int)v5);
	} else {
		v9 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1046492, 10161);
		nox_xxx_wnd_46ABB0((int)v9, 0);
		nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1046536, 0);
		*(uint32_t*)(dword_5d4594_1046536 + 4) |= 8u;
		nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1046504, 1);
		sub_46AD20(*(uint32_t**)&dword_5d4594_1046504, 10134, 10135, 0);
		sub_46ACE0(*(uint32_t**)&dword_5d4594_1046492, 10141, 10141, 1);
		v5 = (char*)a1;
		nox_client_guiserv_updateMapList_458230(*(unsigned short*)(a1 + 52), (char*)a1, 0);
		nox_itow(*(unsigned short*)(a1 + 54), WideCharStr, 10);
		nox_window_call_field_94(*(int*)&dword_5d4594_1046516, 16414, (int)WideCharStr, 0);
		nox_itow(*(unsigned char*)(a1 + 56), WideCharStr, 10);
		nox_window_call_field_94(*(int*)&dword_5d4594_1046520, 16414, (int)WideCharStr, 0);
		if (0) {
			v10 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1046492, 10210);
			nox_window_set_hidden((int)v10, 0);
			sub_457FE0();
		}
	}
	v11 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1046508, 10331);
	if (nox_xxx_CheckGameplayFlags_417DA0(2)) {
		v12 = v11[9];
		LOBYTE(v12) = v12 | 4;
		v11[9] = v12;
	}
	v13 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1046508, 10333);
	if (nox_xxx_CheckGameplayFlags_417DA0(1)) {
		v14 = v13[9];
		LOBYTE(v14) = v14 | 4;
		v13[9] = v14;
	}
	v15 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1046492, 10119);
	v16 = nox_xxx_guiServerOptionsGetGametypeName_4573C0(*((uint16_t*)v5 + 26));
	nox_window_call_field_94((int)v15, 16385, v16, 0);
	sub_459880((int)v5);
	LOWORD(v17) = *((uint16_t*)v5 + 26);
	return sub_457F30((v17 >> 12) & 1);
}
// 457F15: variable 'v17' is possibly undefined

//----- (00457F30) --------------------------------------------------------
int sub_457F30(int a1) {
	int v1;       // esi
	uint32_t* v2; // eax
	uint32_t* v3; // eax
	uint32_t* v4; // esi
	uint32_t* v5; // eax

	v1 = a1 != 1;
	v2 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1046504, 10152);
	nox_xxx_wnd_46ABB0((int)v2, v1);
	v3 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1046504, 10141);
	nox_xxx_wnd_46ABB0((int)v3, v1);
	nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1046516, v1);
	nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1046520, v1);
	if (a1 == 1) {
		v4 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1046492, 10122);
		if (!wndIsShown_nox_xxx_wndIsShown_46ACC0((int)v4)) {
			nox_window_set_hidden((int)v4, 1);
		}
	}
	v5 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1046492, 10183);
	return nox_window_set_hidden((int)v5, a1);
}

// 458291: variable 'v5' is possibly undefined

//----- (00458500) --------------------------------------------------------
int nox_xxx_windowServerOptionsDrawProc_458500(uint32_t* a1, int a2) {
	uint32_t* v2; // esi
	int xLeft;    // [esp+8h] [ebp-8h]
	int v5;       // [esp+Ch] [ebp-4h]

	nox_point mpos = nox_client_getMousePos_4309F0();
	nox_client_wndGetPosition_46AA60(a1, &xLeft, &v5);
	nox_client_drawRectFilledAlpha_49CF10(xLeft, v5 + 25, a1[2], a1[3] - 25);
	v2 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1046492, 10120);
	if (!wndIsShown_nox_xxx_wndIsShown_46ACC0((int)v2) && !nox_xxx_wndPointInWnd_46AAB0(v2, mpos.x, mpos.y)) {
		nox_xxx_wndClearCaptureMain_46ADE0((int)v2);
		nox_window_set_hidden((int)v2, 1);
	}
	return 1;
}

//----- (00458590) --------------------------------------------------------
int nox_xxx_guiServerOptionsProc_458590(int a1, int a2, int a3, int a4) {
	if (a2 == 21) {
		if (a3 != 1) {
			return 0;
		}
		if (a4 == 2) {
			nox_xxx_clientPlaySoundSpecial_452D80(231, 100);
			nox_xxx_guiServerOptionsHide_4597E0(0);
		}
	}
	return 1;
}

// 458F5A: variable 'v34' is possibly undefined

//----- (00459150) --------------------------------------------------------
void sub_4537F0();
char sub_459150() {
	char* v0;          // ebp
	int i;             // edi
	char* v2;          // esi
	int v3;            // ebx
	wchar2_t* v4;       // eax
	char* v5;          // eax
	const char* v6;    // edi
	char* v7;          // eax
	int v8;            // ebx
	char* v9;          // esi
	unsigned char v10; // dl
	char* v11;         // edi
	char* v12;         // eax
	char* v13;         // eax
	int* v14;          // eax
	int* v16;          // [esp-Ch] [ebp-2Ch]
	int* v17;          // [esp-4h] [ebp-24h]
	short v18;         // [esp-4h] [ebp-24h]
	char v19[16];      // [esp+10h] [ebp-10h]

	v0 = sub_4165B0();
	if (v0[52] & 0x60) {
		for (i = nox_xxx_getTeamCounter_417DD0(); i < 2; ++i) {
			v2 = nox_xxx_teamCreate_4186D0(0);
			v3 = (i != 0) + 1;
			if (sub_40A740()) {
				sub_418800((wchar2_t*)v2, (wchar2_t*)getMemAt(0x5D4594, 1046564), 0);
			} else {
				v4 = nox_server_teamTitle_418C20(v3);
				sub_418800((wchar2_t*)v2, v4, 1);
			}
			v2[56] = v3;
			sub_4184D0((wchar2_t*)v2);
		}
		nox_xxx_SetGameplayFlag_417D50(4);
	}
	LOBYTE(v5) = v0[57];
	if ((uint8_t)v5) {
		v5 = sub_409B80();
		v6 = v5;
	} else {
		v6 = v0;
	}
	if (*v6) {
		v7 = nox_xxx_mapGetMapName_409B40();
		v8 = nox_strcmpi(v6, v7);
		sub_4165F0(1, 0);
		v9 = sub_4165D0(0);
		nox_xxx_gameSetServername_40A440(v9 + 9);
		if (!((*((uint16_t*)v9 + 26) >> 8) & 0x10)) {
			sub_409FB0_settings(*((uint16_t*)v9 + 26), *((uint16_t*)v9 + 27));
			sub_40A040_settings(*((uint16_t*)v9 + 26), v9[56]);
		}
		if (v8 && 0) {
			if (!nox_common_gameFlags_check_40A5C0(128)) {
				nox_xxx_net_4263C0();
				sub_4264D0();
				sub_416150(15, 0);
			}
			if (!v9[57]) {
				sub_41D650();
			}
		}
		if (nox_common_gameFlags_check_40A5C0(128)) {
			nox_common_gameFlags_unset_40A540(49152);
			nox_xxx_setGameFlags_40A4D0(*((uint16_t*)v9 + 26) & 0xC000);
		}
		v9[57] = 0;
		if (v8) {
			v19[0] = 0;
			*(uint32_t*)&v19[1] = 0;
			*(uint32_t*)&v19[5] = 0;
			*(uint32_t*)&v19[9] = 0;
			v10 = getMemByte(0x587000, 131672);
			strcpy(v19, v6);
			v11 = &v19[strlen(v19)];
			*(uint32_t*)v11 = *getMemU32Ptr(0x587000, 131668);
			v11[4] = v10;
			nox_xxx_mapLoad_4D2450(v19);
			sub_416690();
			sub_4165D0(1);
			LOBYTE(v5) = (unsigned int)nox_xxx_guiServerOptionsHide_4597E0(0);
		} else {
			nox_xxx_spellEnableAll_424BD0();
			sub_4537F0();
			v17 = (int*)sub_459870();
			v12 = nox_xxx_cliGamedataGet_416590(1);
			sub_57AAA0("user.rul", v12, v17);
			v13 = nox_server_currentMapGetFilename_409B30();
			sub_57A950(v13);
			v18 = nox_common_gameFlags_getVal_40A5B0();
			v16 = (int*)sub_459870();
			v14 = (int*)nox_xxx_cliGamedataGet_416590(0);
			sub_57A1E0(v14, "user.rul", v16, 3, v18);
			LOBYTE(v5) = (unsigned int)nox_xxx_guiServerOptionsHide_4597E0(0);
		}
	}
	return (char)v5;
}

//----- (004593B0) --------------------------------------------------------
int sub_4593B0(int a1) {
	uint32_t* v1; // esi
	uint32_t* v2; // ebx
	uint32_t* v3; // edi
	int result;   // eax

	v1 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1046492, 10163);
	v2 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1046492, 10161);
	v3 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1046492, 10162);
	if (a1) {
		if (a1 == 1) {
			v1[15] = dword_5d4594_1046360;
			v3[15] = dword_5d4594_1046356;
			if (dword_5d4594_1046532) {
				sub_456D60(1);
				dword_5d4594_1046532 = 0;
			}
			if (dword_5d4594_1046540) {
				sub_4AD820();
				dword_5d4594_1046540 = 0;
			}
			dword_5d4594_1046528 = nox_xxx_guiServerAccessLoad_4541D0(*(int*)&dword_5d4594_1046492);
			result = sub_459560(10161);
		} else {
			result = a1 - 2;
			if (a1 == 2) {
				v1[15] = dword_5d4594_1046360;
				if (nox_common_gameFlags_check_40A5C0(1)) {
					v2[15] = dword_5d4594_1046360;
				} else {
					v2[21] = dword_5d4594_1046360;
				}
				if (dword_5d4594_1046528) {
					sub_4557D0(1);
					dword_5d4594_1046528 = 0;
				}
				if (dword_5d4594_1046540) {
					sub_4AD820();
					dword_5d4594_1046540 = 0;
				}
				dword_5d4594_1046532 = nox_xxx_guiServerPlayersLoad_456270(*(int*)&dword_5d4594_1046492);
				result = sub_459560(10162);
			}
		}
	} else {
		if (nox_common_gameFlags_check_40A5C0(1)) {
			v2[15] = dword_5d4594_1046356;
		} else {
			v2[21] = dword_5d4594_1046356;
		}
		v3[15] = dword_5d4594_1046356;
		if (dword_5d4594_1046532) {
			sub_456D60(1);
			dword_5d4594_1046532 = 0;
		} else if (dword_5d4594_1046528) {
			sub_4557D0(1);
			dword_5d4594_1046528 = 0;
		}
		dword_5d4594_1046540 = nox_xxx_gui_4AD320(*(int*)&dword_5d4594_1046492);
		result = sub_459560(10163);
	}
	return result;
}

//----- (00459560) --------------------------------------------------------
int sub_459560(int a1) {
	uint32_t* v1; // ebx
	int v2;       // esi
	uint32_t* v3; // edi
	uint32_t* v4; // eax
	uint32_t* v5; // edi
	uint32_t* v6; // edi
	uint32_t* v7; // edi

	v1 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1046492, a1);
	v2 = v1[99];
	switch (a1) {
	case 10161:
		v6 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1046492, 10163);
		sub_46B120(v6, 0);
		sub_46B120(v6, v2);
		v4 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1046492, 10162);
		v7 = v4;
		sub_46B120(v4, 0);
		sub_46B120(v7, v2);
		break;
	case 10162:
		v5 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1046492, 10163);
		sub_46B120(v5, 0);
		sub_46B120(v5, v2);
		v4 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1046492, 10161);
		v7 = v4;
		sub_46B120(v4, 0);
		sub_46B120(v7, v2);
		break;
	case 10163:
		v3 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1046492, 10162);
		sub_46B120(v3, 0);
		sub_46B120(v3, v2);
		v4 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1046492, 10161);
		v7 = v4;
		sub_46B120(v4, 0);
		sub_46B120(v7, v2);
		break;
	}
	sub_46B120(v1, 0);
	return sub_46B120(v1, v2);
}

//----- (00459700) --------------------------------------------------------
int sub_459700() {
	char* v0;          // esi
	int v1;            // eax
	int v2;            // edi
	const wchar2_t* v3; // eax
	wchar2_t v5[256] = {0};   // [esp+8h] [ebp-200h]

	sub_416580();
	v0 = sub_4165B0();
	v1 = nox_window_call_field_94(*(int*)&dword_5d4594_1046496, 16404, 0, 0);
	v2 = v1;
	v3 = (const wchar2_t*)nox_window_call_field_94(*(int*)&dword_5d4594_1046496, 16406, v1, 0);
	if (v3) {
		nox_wcsncpy(v5, v3, 255);
		v5[255] = 0;
	} else {
		v5[0] = 0;
	}
	nox_window_call_field_94(*(int*)&dword_5d4594_1046496, 16398, v2, 0);
	nox_window_call_field_94(*(int*)&dword_5d4594_1046496, 16402, v2, 0);
	nox_window_call_field_94(*(int*)&dword_5d4594_1046496, 16397, (int)v5, -1);
	nox_window_call_field_94(*(int*)&dword_5d4594_1046496, 16403, v2, 0);
	sub_57A9F0(v0, "user.rul");
	sub_57A1E0((int*)v0, "user.rul", 0, 5, *((uint16_t*)v0 + 26));
	sub_459880((int)v0);
	return sub_459D50(1);
}

//----- (004597E0) --------------------------------------------------------
int* nox_xxx_guiServerOptionsHide_4597E0(int a1) {
	int v1;      // eax
	int* result; // eax

	if (dword_5d4594_1046492) {
		v1 = nox_xxx_wndGetFocus_46B4F0();
		if (nox_window_is_child(*(int*)&dword_5d4594_1046492, v1)) {
			nox_xxx_windowFocus_46B500(0);
		}
		nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1046492);
		dword_5d4594_1046492 = 0;
		sub_456D60(0);
		sub_4BE610();
		dword_5d4594_1046532 = 0;
		sub_4557D0(0);
		dword_5d4594_1046528 = 0;
		sub_4AD820();
		dword_5d4594_1046536 = 0;
	}
	result = (int*)a1;
	if (a1) {
		result = sub_57ADF0(getMemIntPtr(0x5D4594, 1045956));
		dword_587000_129656 = 1;
	}
	return result;
}

//----- (00459870) --------------------------------------------------------
char* sub_459870() { return (char*)getMemAt(0x5D4594, 1045956); }

//----- (00459A40) --------------------------------------------------------
int sub_459A40(char* a1) {
	char v2[100];    // [esp+0h] [ebp-12Ch]
	wchar2_t v3[100]; // [esp+64h] [ebp-C8h]

	strncpy(v2, a1, 0xFu);
	v2[15] = 0;
	nox_swprintf(v3, L"%S", v2);
	return nox_window_call_field_94(dword_5d4594_1046512, 16414, (int)v3, 0);
}

//----- (00459AA0) --------------------------------------------------------
char* sub_459AA0(void* a1p) {
	int a1 = a1p;
	int v1;            // eax
	const wchar2_t* v2; // eax
	const wchar2_t* v3; // eax
	char* result;      // eax
	int v5;            // eax
	const char* v6;    // edi
	char v7[100];      // [esp+Ch] [ebp-64h]

	v1 = nox_window_call_field_94(dword_5d4594_1046512, 16413, 0, 0);
	nox_sprintf(v7, "%S", v1);
	strncpy((char*)(a1 + 9), v7, 0xFu);
	*(uint16_t*)(a1 + 52) = sub_459C10();
	memcpy((void*)(a1 + 24), sub_453F90(), 0x14u);
	*(uint32_t*)(a1 + 44) = *(uint32_t*)sub_453600();
	*(uint32_t*)(a1 + 48) = sub_453610();
	v2 = (const wchar2_t*)nox_window_call_field_94(*(int*)&dword_5d4594_1046516, 16413, 0, 0);
	if (*v2) {
		*(uint16_t*)(a1 + 54) = nox_wcstol(v2, 0, 10);
	}
	v3 = (const wchar2_t*)nox_window_call_field_94(*(int*)&dword_5d4594_1046520, 16413, 0, 0);
	if (*v3) {
		*(uint8_t*)(a1 + 56) = nox_wcstol(v3, 0, 10);
	}
	*(uint8_t*)(a1 + 57) =
		(nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1046492, 10122)->draw_data.field_0 >> 2) & 1;
	result = (char*)nox_window_call_field_94(*(int*)&dword_5d4594_1046496, 16404, 0, 0);
	if ((int)result >= 0 && (v5 = nox_window_call_field_94(*(int*)&dword_5d4594_1046496, 16406, (int)result, 0),
							 nox_sprintf(v7, "%S", v5), (result = strtok(v7, "\t")) != 0)) {
		v6 = result;
		result = 0;
		strcpy((char*)a1, v6);
	} else {
		*(uint8_t*)a1 = 0;
	}
	return result;
}

//----- (00459C30) --------------------------------------------------------
int sub_459C30() {
	char* v0;     // esi
	int result;   // eax
	uint32_t* v2; // edi
	int v3;       // eax

	v0 = nox_xxx_cliGamedataGet_416590(1);
	sub_453F70(v0 + 24);
	sub_4535E0((int*)v0 + 11);
	sub_4535F0(*((uint32_t*)v0 + 12));
	result = dword_5d4594_1046492;
	if (dword_5d4594_1046492) {
		sub_459880((int)v0);
		sub_4BDF70((int*)v0);
		nox_client_guiserv_updateMapList_458230(*((unsigned short*)v0 + 26), v0, 0);
		sub_457460((int)v0);
		if (0) {
			sub_457FE0();
		}
		v2 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1046492, 10119);
		v3 = nox_xxx_guiServerOptionsGetGametypeName_4573C0(*((uint16_t*)v0 + 26));
		result = nox_window_call_field_94((int)v2, 16385, v3, 0);
	}
	return result;
}

//----- (00459D50) --------------------------------------------------------
int sub_459D50(int a1) {
	int result; // eax

	result = a1;
	*getMemU32Ptr(0x5D4594, 1046544) = a1;
	return result;
}

//----- (00459D60) --------------------------------------------------------
int sub_459D60() { return *getMemU32Ptr(0x5D4594, 1046544); }

//----- (00459D80) --------------------------------------------------------
int sub_459D80(int a1) {
	int result; // eax

	result = dword_5d4594_1046492;
	if (dword_5d4594_1046492) {
		result = nox_window_set_hidden(*(int*)&dword_5d4594_1046492, a1);
	}
	return result;
}

//----- (00459DA0) --------------------------------------------------------
int sub_459DA0() { return dword_5d4594_1046492 != 0; }

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
