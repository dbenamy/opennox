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
extern uint32_t dword_5d4594_1045424;
extern uint32_t dword_5d4594_1049696;
extern uint32_t dword_5d4594_831276;
extern uint32_t dword_5d4594_1046648;
extern uint32_t dword_5d4594_1046640;
extern uint32_t dword_5d4594_831084;
extern uint32_t dword_5d4594_1046956;
extern uint32_t dword_5d4594_1047936;
extern uint32_t dword_5d4594_1045436;
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
