#include <math.h>

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
#include "GAME3_3.h"
#include "GAME4_1.h"
#include "GAME5_2.h"
#include "client__draw__debugdraw.h"
#include "client__draw__staticdraw.h"
#include "client__drawable__drawable.h"

#include "client__gui__gadgets__listbox.h"
#include "client__gui__guibook.h"
#include "client__gui__guicon.h"
#include "client__gui__guishop.h"
#include "client__gui__servopts__guiserv.h"
#include "client__gui__tooltip.h"
#include "client__gui__window.h"

#include "client__video__draw_common.h"

#include "common/fs/nox_fs.h"
#include "common__binfile.h"
#include "common__magic__speltree.h"
#include "common__net_list.h"
#include "common__strman.h"
#include "input.h"
#include "input_common.h"
#include "operators.h"

extern uint32_t dword_8531A0_2576;
extern uint32_t dword_5d4594_1098456;
extern uint32_t dword_5d4594_1096636;
extern uint32_t dword_5d4594_1098620;
extern uint32_t dword_5d4594_1123520;
extern uint32_t dword_5d4594_1193188;
extern uint32_t dword_5d4594_1098596;
extern uint32_t dword_5d4594_1098600;
extern uint32_t dword_5d4594_1098616;
extern uint32_t dword_5d4594_1098604;
extern uint32_t dword_5d4594_3804684;
extern uint32_t nox_xxx_waypointCounterMB_587000_154948;
extern uint32_t dword_5d4594_1098592;
extern uint32_t dword_5d4594_1098580;
extern uint32_t dword_5d4594_3798812;
extern uint32_t dword_5d4594_3798800;
extern uint32_t dword_5d4594_3798828;
extern uint64_t qword_581450_9552;
extern uint64_t qword_581450_9544;
extern uint32_t dword_5d4594_1098624;
extern uint32_t dword_5d4594_3798816;
extern uint32_t dword_5d4594_3798808;
extern uint32_t dword_5d4594_3798832;
extern uint32_t dword_5d4594_1193384;
extern uint32_t dword_5d4594_1193360;
extern uint32_t dword_5d4594_3798836;
extern uint32_t dword_5d4594_1107036;
extern uint32_t dword_5d4594_251572;
extern uint32_t dword_5d4594_1098628;
extern uint32_t dword_5d4594_3798804;
extern uint32_t dword_5d4594_1098576;
extern uint32_t dword_5d4594_3798820;
extern uint32_t dword_5d4594_3798824;
extern void* dword_587000_155144;
extern uint32_t dword_5d4594_3798840;
extern void* dword_5d4594_1123524;
extern uint32_t dword_5d4594_1193380;
extern int nox_win_width;
extern int nox_win_height;

extern uint32_t nox_color_white_2523948;
extern uint32_t nox_color_yellow_2589772;

extern nox_render_data_t* nox_draw_curDrawData_3799572;

extern obj_5D4594_2650668_t** ptr_5D4594_2650668;

extern uint32_t nox_tile_def_cnt;
extern nox_tileDef_t nox_tile_defs_arr[176];

uint32_t dword_5d4594_1193156 = 0;
uint8_t** nox_pixbuffer_rows_3798784 = 0;

void* dword_5d4594_1096640 = 0;

void (*func_587000_154940)(int2*, uint32_t, uint32_t) = nox_xxx_tileDraw_4815E0;
int (*func_587000_154944)(int, int) = nox_xxx_drawTexEdgesProbably_481900;

void* nox_client_spriteUnderCursorXxx_1096644 = 0;
uint32_t nox_client_highResFloors_154952 = 1;
void* nox_video_tileBuf_ptr_3798796 = 0;
void* nox_video_tileBuf_end_3798844 = 0;

//----- (00476080) --------------------------------------------------------
int sub_476080(unsigned char* a1) {
	int v1;     // esi
	int v2;     // ecx
	int result; // eax
	int v4;     // edx
	int v5;     // ecx

	if (!*getMemU32Ptr(0x852978, 8)) {
		return 23 * a1[6] + 11;
	}
	switch (*a1) {
	case 0u:
	case 3u:
	case 0xBu:
		v1 = -23;
		v2 = 23 * a1[5] + 22;
		result = 23 * a1[6];
		break;
	case 1u:
	case 4u:
	case 0xCu:
		v1 = 23;
		v2 = 23 * a1[5];
		result = 23 * a1[6];
		break;
	default:
		return 23 * a1[6] + 11;
	}
	v4 = *(uint32_t*)(*getMemU32Ptr(0x852978, 8) + 12) - v2;
	v5 = v1 * (*(uint32_t*)(*getMemU32Ptr(0x852978, 8) + 16) - result) - 23 * v4;
	if (v1 < 0) {
		v5 = 23 * v4 - v1 * (*(uint32_t*)(*getMemU32Ptr(0x852978, 8) + 16) - result);
	}
	if (v5 < 0) {
		result += 22;
	}
	return result;
}

//----- (004761B0) --------------------------------------------------------
int sub_4761B0(nox_drawable* a1p) {
	int a1 = a1p;
	int result; // eax
	int v2;     // edx
	int v3;     // ecx
	int v4;     // edx

	if (!*getMemU32Ptr(0x852978, 8)) {
		return *(uint32_t*)(a1 + 16) + *getMemIntPtr(0x587000, 196188 + 8 * *(unsigned char*)(a1 + 299)) / 2;
	}
	result = *(uint32_t*)(a1 + 16);
	v2 = 8 * *(unsigned char*)(a1 + 299);
	v3 = (*(uint32_t*)(*getMemU32Ptr(0x852978, 8) + 16) - result) * *getMemIntPtr(0x587000, 196184 + v2) -
		 (*(uint32_t*)(*getMemU32Ptr(0x852978, 8) + 12) - *(uint32_t*)(a1 + 12)) * *getMemIntPtr(0x587000, 196188 + v2);
	if (*getMemIntPtr(0x587000, 196184 + v2) < 0) {
		v3 = (*(uint32_t*)(*getMemU32Ptr(0x852978, 8) + 12) - *(uint32_t*)(a1 + 12)) *
				 *getMemIntPtr(0x587000, 196188 + 8 * *(unsigned char*)(a1 + 299)) -
			 (*(uint32_t*)(*getMemU32Ptr(0x852978, 8) + 16) - result) *
				 *getMemIntPtr(0x587000, 196184 + 8 * *(unsigned char*)(a1 + 299));
	}
	v4 = result + *getMemIntPtr(0x587000, 196188 + 8 * *(unsigned char*)(a1 + 299));
	if (v3 >= 0) {
		if (v4 <= result) {
			result += *getMemIntPtr(0x587000, 196188 + 8 * *(unsigned char*)(a1 + 299));
		}
	} else if (v4 > result) {
		result += *getMemIntPtr(0x587000, 196188 + 8 * *(unsigned char*)(a1 + 299));
	}
	return result;
}

//----- (00476AE0) --------------------------------------------------------
void* sub_476AE0(nox_draw_viewport_t* vp, nox_drawable* dr) {
	unsigned char* a2 = dr;
	unsigned char* v2;                        // ebx
	int result;                               // eax
	int (*result2)(int*, int);                // eax
	int v4;                                   // eax
	int v5;                                   // edi
	int v6;                                   // ecx
	int v7;                                   // ebp
	char* v8;                                 // esi
	int v9;                                   // edx
	int v10;                                  // ecx
	unsigned int v11;                         // edi
	uint8_t* v12;                             // esi
	int v13;                                  // ebx
	int v14;                                  // edx
	unsigned int v15;                         // ecx
	unsigned int v16;                         // ebx
	int v17;                                  // ecx
	int v18;                                  // eax
	int v19;                                  // eax
	int v20;                                  // ebp
	int v21;                                  // edi
	int v22;                                  // ebp
	int v23;                                  // ebx
	int v24;                                  // esi
	int v25;                                  // [esp+10h] [ebp-1Ch]
	int v26;                                  // [esp+14h] [ebp-18h]
	int v27;                                  // [esp+18h] [ebp-14h]
	int v28;                                  // [esp+18h] [ebp-14h]
	unsigned int v29;                         // [esp+1Ch] [ebp-10h]
	int v30;                                  // [esp+20h] [ebp-Ch]
	int v31;                                  // [esp+24h] [ebp-8h]
	int v32;                                  // [esp+28h] [ebp-4h]
	void (*v33)(unsigned int, uint8_t*, int); // [esp+34h] [ebp+8h]

	v2 = a2;
	result2 = (int (*)(int*, int)) * ((uint32_t*)a2 + 75);
	if (result2 == nox_thing_static_draw) {
		if (*((uint32_t*)a2 + 28) & 0x40000 && !(*((uint32_t*)a2 + 30) & 0x1000000)) {
			return result2;
		}
		v4 = *(uint32_t*)(*((uint32_t*)a2 + 76) + 4);
	} else {
		v4 = *(uint32_t*)(*(uint32_t*)(*((uint32_t*)a2 + 76) + 4) + 4 * *((uint32_t*)a2 + 77));
	}
	result = nox_video_getImagePixdata_42FB30(v4);
	if (result) {
		v33 = sub_476D70;
		v5 = *(uint32_t*)result;
		v6 = *((uint32_t*)result + 1);
		v27 = *((uint32_t*)result + 1);
		v7 = *((uint32_t*)result + 2) + *((uint32_t*)v2 + 3) - *v2;
		v8 = (char*)result + 16;
		result =
			(int)(*((uint32_t*)result + 3) + *((uint32_t*)v2 + 4) - *((short*)v2 + 53) - *((short*)v2 + 52) - v2[1]);
		v31 = v5;
		if (v7 < *(int*)&dword_5d4594_3798820 || v7 + v5 >= *(int*)&dword_5d4594_3798820 + dword_5d4594_3798800 ||
			(v9 = dword_5d4594_3798824, (int)result < *(int*)&dword_5d4594_3798824) ||
			(int)result + v6 >= *(int*)&dword_5d4594_3798824 + dword_5d4594_3798808) {
			*((uint32_t*)v2 + 86) = 0;
		} else {
			v10 = nox_xxx_waypointCounterMB_587000_154948;
			if (*(int*)&nox_xxx_waypointCounterMB_587000_154948 <= 0) {
				*((uint32_t*)v2 + 86) = 0;
				v9 = dword_5d4594_3798824;
				v10 = nox_xxx_waypointCounterMB_587000_154948;
			}
			if (v10 - *((int*)v2 + 86) > 1 || v10 <= 0) {
				v11 = nox_video_tileBuf_end_3798844;
				v12 = v8 + 1;
				v29 = nox_video_tileBuf_end_3798844;
				v26 = (uint32_t)nox_video_tileBuf_end_3798844 - (uint32_t)nox_video_tileBuf_ptr_3798796;
				v13 = dword_5d4594_3798804 * ((uint32_t)result + dword_5d4594_3798840 - v9);
				v14 = v7 + dword_5d4594_3798836 - dword_5d4594_3798820;
				v15 = v13 + (uint32_t)nox_video_tileBuf_ptr_3798796 + 2 * v14;
				v25 = v13 + (uint32_t)nox_video_tileBuf_ptr_3798796 + 2 * v14;
				if (v15 >= nox_video_tileBuf_end_3798844) {
					v15 -= v26;
					v25 = v15;
				}
				result = (int)(v27 - 1);
				if (v27) {
					v30 = v27;
					do {
						v16 = v15;
						v28 = v31;
						if (v31 > 0) {
							do {
								v17 = (unsigned char)v12[1];
								v18 = *v12 & 0xF;
								v12 += 2;
								v19 = v18 - 1;
								v32 = v17;
								v20 = 2 * v17;
								if (v19) {
									if (v19 == 1) {
										if (v16 >= v11 || v16 + v20 < v11) {
											v33(v16, v12, 2 * v17);
											v12 += v20;
											v16 += v20;
										} else {
											v21 = v16 + v20 - v29;
											v22 = v29 - v16;
											v33(v16, v12, v29 - v16);
											v23 = nox_video_tileBuf_ptr_3798796;
											v24 = (int)&v12[v22];
											v33(nox_video_tileBuf_ptr_3798796, (uint8_t*)v24, v21);
											v12 = (uint8_t*)(v21 + v24);
											v16 = v21 + v23;
											v11 = v29;
										}
									}
								} else {
									v16 += v20;
									if (v16 >= v11) {
										v16 -= v26;
									}
								}
								v28 -= v32;
							} while (v28 > 0);
							v15 = v25;
						}
						v15 += dword_5d4594_3798804;
						v25 = v15;
						if (v15 >= v11) {
							v15 -= v26;
							v25 = v15;
						}
						result = (int)--v30;
					} while (v30);
				}
			} else {
				*((uint32_t*)v2 + 86) = v10;
			}
		}
	}
	return result;
}

//----- (00476D70) --------------------------------------------------------
short sub_476D70(uint32_t* a1, int* a2, unsigned int a3) {
	uint32_t* v3;  // edi
	signed int v4; // ecx
	int* v5;       // esi
	int v6;        // eax

	v3 = a1;
	v4 = a3 >> 2;
	v5 = a2;
	if (a3 >> 2) {
		do {
			v6 = *v5;
			++v5;
			*v3 = v6;
			++v3;
		} while (v4-- > 1);
	}
	if (a3 & 3) {
		LOWORD(v6) = *(uint16_t*)v5;
		*(uint16_t*)v3 = *(uint16_t*)v5;
	}
	return v6;
}

//----- (00476E00) --------------------------------------------------------
int nox_client_setPhonemeFrame_476E00(int a1) {
	int result; // eax

	result = a1;
	*getMemU32Ptr(0x5D4594, 1096596 + 4 * a1) = gameFrame();
	return result;
}

//----- (00476E20) --------------------------------------------------------
uint32_t* sub_476E20() {
	int v0;       // esi
	char* v1;     // eax
	uint32_t* v2; // esi

	v0 = 0;
	while (1) {
		v1 = nox_xxx_gLoadImg_42F970(*(const char**)getMemAt(0x587000, 151272 + v0));
		*getMemU32Ptr(0x5D4594, 1096564 + v0) = v1;
		if (!v1) {
			break;
		}
		v0 += 4;
		if (v0 >= 32) {
			v2 = nox_window_new(0, 64, (nox_win_width - 100) / 2, (nox_win_height - 100) / 2, 1, 1, 0);
			nox_window_set_all_funcs(v2, 0, sub_476E90, 0);
			return v2;
		}
	}
	return 0;
}

//----- (00476E90) --------------------------------------------------------
int sub_476E90() {
	unsigned char* v0; // edi
	int v1;            // esi

	v0 = getMemAt(0x587000, 151208);
	v1 = 0;
	do {
		if (*getMemU32Ptr(0x5D4594, 1096596 + v1)) {
			nox_client_drawImageAt_47D2C0(*getMemU32Ptr(0x5D4594, 1096564 + v1),
										  nox_win_width / 2 + *(uint32_t*)v0 - 16,
										  *((uint32_t*)v0 + 1) + nox_win_height / 2 - 41);
			if ((unsigned int)(gameFrame() - *getMemU32Ptr(0x5D4594, 1096596 + v1)) > 3) {
				*getMemU32Ptr(0x5D4594, 1096596 + v1) = 0;
			}
		}
		v0 += 8;
		v1 += 4;
	} while ((int)v0 < (int)getMemAt(0x587000, 151272));
	return 1;
}

//----- (00476F40) --------------------------------------------------------
unsigned int nox_xxx_packetGetMarshall_476F40() {
	unsigned int result; // eax

	if (dword_5d4594_1096640) {
		result = nox_xxx_netGetUnitCodeCli_578B00(*(int*)&dword_5d4594_1096640);
	} else {
		result = 0;
	}
	return result;
}

//----- (00476FA0) --------------------------------------------------------
void nox_xxx_clientEnumHover_476FA0() {
	int4 v2; // [esp+10h] [ebp-10h]

	if (!*getMemU32Ptr(0x5D4594, 1096632)) {
		*getMemU32Ptr(0x5D4594, 1096632) = nox_xxx_getNameId_4E3AA0("Glyph");
	}
	nox_point mpos = nox_client_getMousePos_4309F0();
	sub_473970(&mpos, &mpos);
	dword_5d4594_1096640 = 0;
	nox_client_spriteUnderCursorXxx_1096644 = 0;
	*getMemU32Ptr(0x5D4594, 1096628) = 0;
	v2.field_0 = mpos.x - 96;
	v2.field_8 = mpos.x + 96;
	v2.field_C = mpos.y + 96;
	v2.field_4 = mpos.y - 96;
	dword_5d4594_1096636 = 0;
	nox_xxx_forEachSprite_49AB00(&v2, nox_xxx_clientOnCursorHover_477050, &mpos);
}

//----- (00477050) --------------------------------------------------------
int nox_xxx_client_4984B0_drawable(nox_drawable* dr);
void nox_xxx_clientOnCursorHover_477050(int arg0, int a2) {
	int v2;    // esi
	int v3;    // eax
	int v4;    // eax
	int v5;    // ecx
	char* v6;  // eax
	int v7;    // eax
	float v8;  // ebp
	int v9;    // edi
	int v10;   // ebx
	int v11;   // eax
	int v12;   // edx
	int v13;   // edi
	int v14;   // ebx
	float v15; // eax
	int v16;   // edx
	int v17;   // eax
	int v18;   // edi
	int v19;   // ebx
	int v20;   // eax
	int v21;   // ecx
	int v22;   // edi
	float v23; // [esp+0h] [ebp-24h]
	float v24; // [esp+0h] [ebp-24h]
	float v25; // [esp+0h] [ebp-24h]
	float v26; // [esp+0h] [ebp-24h]
	float2 a3; // [esp+14h] [ebp-10h]
	float2 a1; // [esp+1Ch] [ebp-8h]
	int v29;   // [esp+28h] [ebp+4h]

	if (!*getMemU32Ptr(0x5D4594, 1096648)) {
		*getMemU32Ptr(0x5D4594, 1096648) = nox_xxx_getTTByNameSpriteMB_44CFC0("Polyp");
	}
	v2 = arg0;
	if (arg0 == *getMemU32Ptr(0x852978, 8)) {
		return;
	}
	v3 = *(uint32_t*)(arg0 + 120);
	if (!((v3 & 0x8000) == 0 && (!nox_client_drawable_testBuff_4356C0(arg0, 0) ||
								 nox_client_drawable_testBuff_4356C0(*getMemIntPtr(0x852978, 8), 21)))) {
		return;
	}
	v4 = *(uint32_t*)(arg0 + 112);
	if (!(!(v4 & 2) || (v5 = *(uint32_t*)(arg0 + 116), !(v5 & 0x4000)))) {
		return;
	}
	if (!(v4 & 0x80400206 || *(uint32_t*)(arg0 + 108) == *getMemU32Ptr(0x5D4594, 1096648))) {
		return;
	}
	if (!nox_xxx_client_4984B0_drawable(arg0)) {
		return;
	}
	if (!(!(*(uint8_t*)(arg0 + 112) & 4) ||
		  (v6 = nox_common_playerInfoGetByID_417040(*(uint32_t*)(arg0 + 128))) != 0 && !(v6[3680] & 1))) {
		return;
	}
	v7 = *(uint32_t*)(arg0 + 112);
	if (!((!(v7 & 0x400000) || (*(uint8_t*)(arg0 + 116) & 0x80)) && (!(v7 & 2) || *(uint32_t*)(arg0 + 276) != 10))) {
		return;
	}
	v23 = (double)*(int*)(arg0 + 16) - *(float*)(arg0 + 100) - (double)*(short*)(arg0 + 104);
	v29 = nox_float2int(v23);
	v24 = (double)*(int*)(v2 + 16) - *(float*)(v2 + 96) - (double)*(short*)(v2 + 104);
	v8 = COERCE_FLOAT(nox_float2int(v24));
	a3.field_0 = v8;
	if (*(uint32_t*)(v2 + 44) == 2) {
		v25 = *(float*)(v2 + 48) * *(float*)(v2 + 48);
		LODWORD(a3.field_0) = nox_float2int(v25);
		v17 = nox_float2int(*(float*)(v2 + 48));
		v18 = *(uint32_t*)(v2 + 12);
		v19 = *(uint32_t*)(v2 + 12) - v17;
		v20 = v18 + nox_float2int(*(float*)(v2 + 48));
		v21 = *(uint32_t*)(a2 + 4);
		if (v21 <= SLODWORD(v8)) {
			v8 = *(float*)&v29;
			if (v21 >= v29) {
				if (*(int*)a2 <= v19 || *(int*)a2 >= v20) {
					return;
				}
				goto LABEL_38;
			}
		}
		v15 = a3.field_0;
		v16 = (*(uint32_t*)a2 - v18) * (*(uint32_t*)a2 - v18) + (v21 - LODWORD(v8)) * (v21 - LODWORD(v8));
	} else {
		if (*(uint32_t*)(v2 + 44) != 3) {
			return;
		}
		a1.field_0 = (double)*(int*)(v2 + 12);
		a1.field_4 = (double)SLODWORD(a3.field_0);
		a3.field_0 = (double)*(int*)a2;
		a3.field_4 = (double)*(int*)(a2 + 4);
		if (nox_xxx_map_57B850(&a1, (float*)(v2 + 44), &a3) ||
			(a1.field_4 = (double)v29, nox_xxx_map_57B850(&a1, (float*)(v2 + 44), &a3)) ||
			(v9 = *(uint32_t*)(v2 + 12) + nox_float2int(*(float*)(v2 + 72)),
			 v10 = v29 + nox_float2int(*(float*)(v2 + 76)), v11 = LODWORD(v8) + nox_float2int(*(float*)(v2 + 76)),
			 *(uint32_t*)a2 > v9) &&
				*(uint32_t*)a2 < *(int*)(v2 + 12) && (v12 = *(uint32_t*)(a2 + 4), v12 > v10) && v12 < v11) {
			goto LABEL_38;
		}
		v13 = *(uint32_t*)(v2 + 12) + nox_float2int(*(float*)(v2 + 80));
		v14 = v29 + nox_float2int(*(float*)(v2 + 84));
		LODWORD(v15) = LODWORD(v8) + nox_float2int(*(float*)(v2 + 84));
		if (*(int*)a2 < *(int*)(v2 + 12)) {
			return;
		}
		if (*(int*)a2 >= v13) {
			return;
		}
		v16 = *(uint32_t*)(a2 + 4);
		if (v16 <= v14) {
			return;
		}
	}
	if (v16 >= SLODWORD(v15)) {
		return;
	}
LABEL_38:
	v26 = (double)*(short*)(v2 + 104) + (double)*(int*)(v2 + 16) + *(float*)(v2 + 96);
	v22 = nox_float2int(v26);
	if (v22 > *getMemIntPtr(0x5D4594, 1096628)) {
		*getMemU32Ptr(0x5D4594, 1096628) = v22;
		dword_5d4594_1096640 = v2;
	}
	if (v2 != *getMemU32Ptr(0x852978, 8) && v22 > *(int*)&dword_5d4594_1096636 && nox_xxx_client_57B400(v2)) {
		if (dword_8531A0_2576 && *(uint8_t*)(dword_8531A0_2576 + 2251) == 1 &&
			*(uint32_t*)(v2 + 108) == *getMemU32Ptr(0x5D4594, 1096632)) {
			if (!nox_client_spriteUnderCursorXxx_1096644) {
				nox_client_spriteUnderCursorXxx_1096644 = v2;
				dword_5d4594_1096636 = 0;
			}
		} else {
			dword_5d4594_1096636 = v22;
			nox_client_spriteUnderCursorXxx_1096644 = v2;
		}
	}
}

//----- (00477600) --------------------------------------------------------
int nox_xxx_guiCursor_477600() { return *getMemU32Ptr(0x5D4594, 1096672); }

//----- (00479950) --------------------------------------------------------
int sub_479950() {
	void* v2; // [esp+0h] [ebp-4h]

	if (wndIsShown_nox_xxx_wndIsShown_46ACC0(*(int*)&dword_5d4594_1123524) == 1) {
		return 0;
	}
	*getMemU8Ptr(0x5D4594, 1123516) = 0;
	LOWORD(v2) = 720;
	BYTE2(v2) = 0;
	nox_netlist_addToMsgListCli_40EBC0(31, 0, &v2, 3);
	return 1;
}

//----- (004799A0) --------------------------------------------------------
int sub_4799A0() {
	int result;    // eax
	uint32_t* v1;  // edi
	uint32_t* v2;  // ebp
	uint32_t* v3;  // eax
	int v4;        // esi
	uint32_t* v5;  // ebx
	char* v6;      // eax
	uint32_t* v7;  // eax
	char* v8;      // [esp-18h] [ebp-1Ch]
	char* v9;      // [esp-14h] [ebp-18h]
	uint32_t* v10; // [esp+0h] [ebp-4h]

	*getMemU32Ptr(0x5D4594, 1107052) = nox_color_rgb_4344A0(240, 128, 64);
	result = nox_new_window_from_file("Dialog.wnd", nox_xxx_guiDialog_479B00);
	dword_5d4594_1123524 = result;
	if (result) {
		nox_xxx_wndSetWindowProc_46B300(result, sub_479BE0);
		nox_xxx_wndSetDrawFn_46B340(*(int*)&dword_5d4594_1123524, sub_479CB0);
		nox_gui_winSetFunc96_46B070(*(int*)&dword_5d4594_1123524, sub_479D00);
		v1 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1123524, 3904);
		v2 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1123524, 3903);
		v10 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1123524, 3902);
		v3 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1123524, 3901);
		v4 = (int)v3;
		v5 = (uint32_t*)v3[8];
		v9 = nox_xxx_gLoadImg_42F970("UISliderLit");
		v8 = nox_xxx_gLoadImg_42F970("UISliderLit");
		v6 = nox_xxx_gLoadImg_42F970("UISlider");
		sub_4B5700((int)v1, 0, 0, (int)v6, (int)v8, (int)v9);
		nox_xxx_wnd_46B280((int)v1, v4);
		nox_xxx_wnd_46B280((int)v2, v4);
		nox_xxx_wnd_46B280((int)v10, v4);
		v5[9] = v1;
		v5[7] = v2;
		v5[8] = v10;
		v7 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1123524, 3906);
		nox_xxx_wndSetDrawFn_46B340((int)v7, sub_479C40);
		nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1123524, 0);
		nox_window_set_hidden(*(int*)&dword_5d4594_1123524, 1);
		dword_5d4594_1123520 = 0;
		result = 1;
	}
	return result;
}

//----- (00479B00) --------------------------------------------------------
int nox_xxx_guiDialog_479B00(int a1, int a2, int* a3, int a4) {
	int v3;     // esi
	int result; // eax

	if (a2 != 16391) {
		return 0;
	}
	v3 = nox_xxx_wndGetID_46B0A0(a3);
	if (sub_45D9B0()) {
		return 0;
	}
	nox_xxx_clientPlaySoundSpecial_452D80(766, 100);
	switch (v3) {
	case 3906:
		sub_479950();
		result = 0;
		break;
	case 3907:
		nox_xxx_playDialogFile_44D900(*getMemIntPtr(0x5D4594, 1115312), 100);
		result = 0;
		break;
	case 3908:
		*getMemU8Ptr(0x5D4594, 1123516) = 1;
		LOWORD(a2) = 720;
		BYTE2(a2) = 1;
		nox_netlist_addToMsgListCli_40EBC0(31, 0, &a2, 3);
		result = 0;
		break;
	case 3909:
		*getMemU8Ptr(0x5D4594, 1123516) = 2;
		BYTE2(a2) = 2;
		LOWORD(a2) = 720;
		nox_netlist_addToMsgListCli_40EBC0(31, 0, &a2, 3);
		return 0;
	default:
		return 0;
	}
	return result;
}
// 479B4D: variable 'v4' is possibly undefined

//----- (00479BE0) --------------------------------------------------------
int sub_479BE0(uint32_t* a1, int a2, unsigned int a3, int a4) {
	switch (a2) {
	case 5:
	case 6:
	case 7:
	case 9:
	case 10:
	case 11:
	case 13:
	case 14:
	case 15:
		nox_xxx_wndPointInWnd_46AAB0(a1, (unsigned short)a3, a3 >> 16);
		break;
	default:
		return 1;
	}
	return 1;
}

//----- (00479C40) --------------------------------------------------------
int nox_xxx_wndButtonDrawNoImg_4A81D0(int a1, int a2);
int sub_479C40(uint32_t* a1, int a2) {
	char v2;   // bl
	int yTop;  // [esp+8h] [ebp-8h]
	int xLeft; // [esp+Ch] [ebp-4h]

	v2 = nox_xxx_bookGet_430B40_get_mouse_prev_seq();
	if (!sub_44D930() && (v2 & 0x7Fu) < 0x1E && v2 & 8) {
		nox_client_wndGetPosition_46AA60(a1, &xLeft, &yTop);
		sub_49CD30(xLeft, yTop, a1[2], a1[3] - 2, *getMemIntPtr(0x5D4594, 1107052), 4);
	}
	return nox_xxx_wndButtonDrawNoImg_4A81D0((int)a1, a2);
}

//----- (00479CB0) --------------------------------------------------------
int sub_479CB0(int a1, int a2) {
	int v2; // esi
	int v4; // [esp+4h] [ebp-8h]
	int v5; // [esp+8h] [ebp-4h]

	v2 = *(uint32_t*)(a2 + 24);
	nox_client_wndGetPosition_46AA60(*(uint32_t**)&dword_5d4594_1123524, &v4, &v5);
	nox_client_drawImageAt_47D2C0(v2, nox_win_width - NOX_DEFAULT_WIDTH, nox_win_height - NOX_DEFAULT_HEIGHT);
	return 1;
}

//----- (00479D00) --------------------------------------------------------
int sub_479D00() { return 1; }

//----- (00479D10) --------------------------------------------------------
int sub_479D10() {
	int result; // eax

	nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1123524);
	result = 0;
	dword_5d4594_1123524 = 0;
	dword_5d4594_1123520 = 0;
	return result;
}

//----- (0047A260) --------------------------------------------------------
int sub_47A260() { return dword_5d4594_1123520; }

//----- (0047DBC0) --------------------------------------------------------
unsigned char sub_47DBC0() { return getMemByte(0x5D4594, 1193128); }

//----- (0047FCE0) --------------------------------------------------------
int sub_47FCE0(uint32_t* a1, int a2) {
	int v2;            // edx
	unsigned char* v3; // eax
	int v4;            // eax
	int v5;            // esi

	v2 = 0;
	if (*(int*)&dword_5d4594_3804684 > 0) {
		v3 = getMemAt(0x973F18, 6092);
		while (*((uint32_t*)v3 - 1) != a1[3] || *(uint32_t*)v3 != a1[2] || *((uint32_t*)v3 + 1) != a1[21] ||
			   *((uint32_t*)v3 + 2) != a1[26]) {
			++v2;
			v3 += 16;
			if (v2 >= *(int*)&dword_5d4594_3804684) {
				goto LABEL_8;
			}
		}
		return 1;
	}
LABEL_8:
	v4 = 16 * dword_5d4594_3804684;
	v5 = dword_5d4594_3804684 + 1;
	*getMemU32Ptr(0x973F18, 6088 + v4) = a1[3];
	*getMemU32Ptr(0x973F18, 6092 + v4) = a1[2];
	*getMemU32Ptr(0x973F18, 6096 + v4) = a1[21];
	*getMemU32Ptr(0x973F18, 6100 + v4) = a1[26];
	dword_5d4594_3804684 = v5;
	return 1;
}

// 4514E0: using guessed type void  nullsub_4(uint32_t, uint32_t, uint32_t, uint32_t);

//----- (00480220) --------------------------------------------------------
uint8_t* sub_480220(uint8_t* a1, uint8_t* a2) {
	unsigned int v2; // edx
	uint8_t* result; // eax

	result = a2;
	*a1 = 8 * *a2;
	LOWORD(v2) = *(uint16_t*)a2;
	a1[1] = (v2 >> 3) & 0xFC;
	a1[2] = a2[1] & 0xF8;
	return result;
}
// 480232: variable 'v2' is possibly undefined

//----- (00480250) --------------------------------------------------------
uint16_t* sub_480250(uint8_t* a1, uint16_t* a2) {
	uint16_t* result; // eax

	result = a2;
	*a2 = (*a1 >> 3) | (8 * (a1[1] & 0xFC | (32 * (a1[2] & 0xF8))));
	return result;
}

//----- (00481410) --------------------------------------------------------
void sub_481410() { nox_xxx_waypointCounterMB_587000_154948 = -1; }

//----- (00481900) --------------------------------------------------------
char nox_xxx_drawTexEdgesProbably_481900(uint32_t* a1, uint32_t* a2) {
	int v2;            // ebx
	int v3;            // esi
	int v4;            // edx
	int v5;            // ebp
	int v6;            // eax
	int v7;            // edi
	int v8;            // edi
	int v9;            // eax
	uint8_t* v10;      // ebx
	int v11;           // ebp
	unsigned int v12;  // esi
	unsigned int v13;  // edi
	char* v14;         // ebp
	char* v15;         // edx
	int v16;           // ebx
	char* v17;         // edi
	unsigned int v18;  // ecx
	char* v19;         // esi
	unsigned int v20;  // eax
	int v22;           // eax
	int v24;           // [esp+10h] [ebp-1Ch]
	unsigned char v25; // [esp+14h] [ebp-18h]
	char* v26;         // [esp+14h] [ebp-18h]
	char* v27;         // [esp+18h] [ebp-14h]
	int v28;           // [esp+1Ch] [ebp-10h]
	unsigned char v29; // [esp+20h] [ebp-Ch]
	int v30;           // [esp+28h] [ebp-4h]
	unsigned int v31;  // [esp+30h] [ebp+4h]
	unsigned char v32; // [esp+34h] [ebp+8h]
	char* v33;         // [esp+34h] [ebp+8h]
	int* addr;

	nox_tileDef_t* tile = &nox_tile_defs_arr[a2[0]];

	v2 = dword_5d4594_3798824;
	v3 = tile->data_32[a2[1] + tile->field_46];
	v4 = a2[2];
	v5 = dword_5d4594_3798836;
	addr = (*getMemU32Ptr(0x85B3FC, 28676 + 60 * v4) + 4 * (a2[3] + *getMemU16Ptr(0x85B3FC, 28690 + 60 * v4)));
	v6 = *(uint32_t*)addr;
	v7 = dword_5d4594_3798840;
	*getMemU32Ptr(0x5D4594, 2523980 + 4 * v4) = 1;
	v8 = dword_5d4594_3798804 * (v7 + a1[1] - v2) + (uint32_t)nox_video_tileBuf_ptr_3798796 +
		 ((v5 + *a1 - dword_5d4594_3798820) << getMemByte(0x973F18, 7696));
	v9 = nox_video_getImagePixdata_42FB30(v6);
	if (v9) {
		v32 = *(uint8_t*)v9;
		v25 = *(uint8_t*)(v9 + 1);
		v10 = (uint8_t*)(v9 + 2);
		v9 = nox_video_getImagePixdata_42FB30(v3);
		v11 = v9;
		if (v9) {
			v12 = nox_video_tileBuf_end_3798844;
			v9 = v32;
			v13 = dword_5d4594_3798804 * v32 + v8;
			v31 = v13;
			v14 = (char*)((*getMemU32Ptr(0x973CE0, +4 * v32) << getMemByte(0x973F18, 7696)) + v11);
			v27 = v14;
			if (v13 >= nox_video_tileBuf_end_3798844) {
				v13 += (uint32_t)nox_video_tileBuf_ptr_3798796 - (uint32_t)nox_video_tileBuf_end_3798844;
				v31 = v13;
			}
			v28 = v32;
			v30 = v25;
			if (v32 <= (int)v25) {
				do {
					v24 = *getMemU32Ptr(0x973CE0, 384 + 4 * v9);
					v26 = v14;
					v15 = (char*)(v13 + (*getMemU32Ptr(0x973CE0, 192 + 4 * v9) << getMemByte(0x973F18, 7696)));
					if (v24 > 0) {
						do {
							LOBYTE(v9) = *v10;
							v29 = v10[1];
							v33 = v10 + 2;
							v16 = v29 << getMemByte(0x973F18, 7696);
							switch ((uint8_t)v9) {
							case 2:
								if ((unsigned int)&v15[v16] < v12) {
									v19 = v33;
									v18 = v29 << getMemByte(0x973F18, 7696);
									v17 = v15;
								} else {
									memcpy(v15, v33, v12 - (uint32_t)v15);
									v17 = nox_video_tileBuf_ptr_3798796;
									v18 = v16 - (v12 - (uint32_t)v15);
									v19 = &v33[v12 - (uint32_t)v15];
								}
								memcpy(v17, v19, v18);
								v12 = nox_video_tileBuf_end_3798844;
								v13 = v31;
								v33 += v16;
								break;
							case 3:
								if ((unsigned int)&v15[v16] < v12) {
									memcpy(v15, v14, v16);
								} else {
									v20 = v12 - (uint32_t)v15;
									memcpy(v15, v14, v20);
									v14 = v26;
									memcpy(nox_video_tileBuf_ptr_3798796, &v26[v20], v16 - v20);
								}
								v12 = nox_video_tileBuf_end_3798844;
								v13 = v31;
								break;
							case 1:
								break;
							default:
								return v9;
							}
							v15 += v16;
							v22 = v24 - v29;
							v14 += v16;
							v24 -= v29;
							v26 = v14;
							if ((unsigned int)v15 >= v12) {
								v15 += (uint32_t)nox_video_tileBuf_ptr_3798796 - v12;
							}
							v10 = v33;
						} while (v22 > 0);
						v9 = v28;
					}
					v14 = &v27[*getMemU32Ptr(0x973CE0, 384 + 4 * v9) << getMemByte(0x973F18, 7696)];
					v13 += dword_5d4594_3798804;
					v27 += *getMemU32Ptr(0x973CE0, 384 + 4 * v9) << getMemByte(0x973F18, 7696);
					v31 = v13;
					if (v13 >= v12) {
						v13 += (uint32_t)nox_video_tileBuf_ptr_3798796 - v12;
						v31 = v13;
					}
					v28 = ++v9;
				} while (v9 <= v30);
			}
		}
	}
	return v9;
}

//----- (00481BF0) --------------------------------------------------------
void nox_xxx_tileCallDrawEdges_481BF0(int a1, int a2) {
	int i; // esi

	for (i = a2; i; i = *(uint32_t*)(i + 16)) {
		func_587000_154944(a1, i);
	}
}
// 5ACD40: invalid function type has been ignored

//----- (00481C20) --------------------------------------------------------
void nox_xxx_tileDrawMB_481C20_A(nox_draw_viewport_t* vp, int v3) {
	int v17; // esi
	int v63; // [esp-8h] [ebp-54h]
	int v16; // edi
	int v62; // [esp-8h] [ebp-54h]
	int v14; // esi
	int v15; // edi
	int v13; // ecx
	int v12; // eax
	int v11; // ebx
	int v4;  // eax
	int v6;  // ecx
	int v10; // esi
	int v9;  // edx
	int v7;  // edx
	int v8;  // eax
	int j;
	int2 v68; // [esp+1Ch] [ebp-30h]
	if (v3 >= *(int*)&dword_5d4594_3798820 + 23) {
		int v71 = nox_getBackbufWidth() + v3;
		if (v71 <= *(int*)&dword_5d4594_3798800 + dword_5d4594_3798820 - 46 ||
			*(int*)&dword_5d4594_3798812 + *(int*)&dword_5d4594_3798828 - 1 >= 128) {
			return;
		}
		if (v71 > *(int*)&dword_5d4594_3798800 + *(int*)&dword_5d4594_3798820) {
			nox_xxx_tileDrawImpl_4826A0(vp);
			return;
		}
		v7 = dword_5d4594_3798828 + 1;
		dword_5d4594_3798828 = v7;
		v8 = dword_5d4594_3798820 + 46;
		dword_5d4594_3798820 += 46;
		j = dword_5d4594_3798812 + v7 - 2;
		v9 = dword_5d4594_3798836 + 46;
		dword_5d4594_3798836 += 46;
		if (*(int*)&dword_5d4594_3798836 >= *(int*)&dword_5d4594_3798800) {
			dword_5d4594_3798836 = v9 - dword_5d4594_3798800;
			v10 = ++dword_5d4594_3798840;
			if (*(int*)&dword_5d4594_3798840 >= *(int*)&dword_5d4594_3798808) {
				dword_5d4594_3798840 = v10 - dword_5d4594_3798808;
			}
		}
		v4 = dword_5d4594_3798800 + v8 - 92;
	} else {
		if (*(int*)&dword_5d4594_3798828 <= 0) {
			return;
		}
		if (v3 < *(int*)&dword_5d4594_3798820 - 23) {
			nox_xxx_tileDrawImpl_4826A0(vp);
			return;
		}
		v4 = dword_5d4594_3798820 - 46;
		j = --dword_5d4594_3798828;
		dword_5d4594_3798820 -= 46;
		dword_5d4594_3798836 -= 46;
		if (*(int*)&dword_5d4594_3798836 < 0) {
			v6 = *(int*)&dword_5d4594_3798840 - 1;
			bool v5 = *(int*)&dword_5d4594_3798840 - 1 < 0;
			dword_5d4594_3798836 += dword_5d4594_3798800;
			--*(int*)&dword_5d4594_3798840;
			if (v5) {
				*(int*)&dword_5d4594_3798840 = *(int*)&dword_5d4594_3798808 + v6;
			}
		}
	}
	int v76 = v4;
	sub_481410();
	v11 = dword_5d4594_3798824;
	unsigned int v74 = dword_5d4594_3798832;
	if (*(int*)&dword_5d4594_3798832 < *(int*)&dword_5d4594_3798832 + *(int*)&dword_5d4594_3798816) {
		v12 = 44 * dword_5d4594_3798832;
		for (int i = 44 * dword_5d4594_3798832;; v12 = i) {
			HIWORD(v13) = HIWORD(ptr_5D4594_2650668); // TODO: why it's doing it?
			v14 = v12 + (uint32_t)(ptr_5D4594_2650668[j]);
			if (*(uint8_t*)v14 & 2) {
				LOWORD(v13) = *(uint16_t*)(v14 + 24);
				v15 = *(unsigned short*)(v14 + 24);
				v62 = nox_tile_defs_arr[v15].data_32[*(uint32_t*)(v14 + 28) + nox_tile_defs_arr[v15].field_46];
				v68.field_0 = v76;
				v68.field_4 = v11 + 23;
				func_587000_154940(&v68, v62, v13);
				*getMemU32Ptr(0x85B3FC, 228 + 4 * v15) = 1;
				if (*(uint32_t*)(v14 + 40)) {
					nox_xxx_tileCallDrawEdges_481BF0((int)&v68, *(uint32_t*)(v14 + 40));
				}
			}
			if (*(uint8_t*)v14 & 1) {
				LOWORD(v13) = *(uint16_t*)(v14 + 4);
				v16 = *(unsigned short*)(v14 + 4);
				v63 = nox_tile_defs_arr[v16].data_32[*(uint32_t*)(v14 + 8) + nox_tile_defs_arr[v16].field_46];
				v68.field_0 = v76 + 23;
				v68.field_4 = v11;
				func_587000_154940(&v68, v63, v13);
				*getMemU32Ptr(0x85B3FC, 228 + 4 * v16) = 1;
				v17 = *(uint32_t*)(v14 + 20);
				if (v17) {
					nox_xxx_tileCallDrawEdges_481BF0((int)&v68, v17);
				}
			}
			v11 += 46;
			bool v18 = (int)++v74 < *(int*)&dword_5d4594_3798832 + dword_5d4594_3798816;
			i += 44;
			if (!v18) {
				break;
			}
		}
	}
}

void nox_xxx_tileDrawMB_481C20_B(nox_draw_viewport_t* vp, int v78) {
	int v33;  // esi
	int v65;  // [esp-8h] [ebp-54h]
	int v32;  // edi
	int2 v68; // [esp+1Ch] [ebp-30h]
	int v64;  // [esp-8h] [ebp-54h]
	int v31;  // edi
	int v28;  // esi
	char v29; // al
	int v30;  // esi
	int v27;  // ecx
	int v25;  // eax
	int v26;  // ebx
	int v21;  // eax
	int v20;  // edi
	int v19;  // esi
	int v23;  // edi
	int v24;  // ecx
	int v22;  // esi
	int v76;
	if ((int)v78 >= *(int*)&dword_5d4594_3798824 + 23) {
		if ((int)v78 + nox_getBackbufHeight() <= *(int*)&dword_5d4594_3798824 + *(int*)&dword_5d4594_3798808) {
			return;
		}
		v22 = dword_5d4594_3798832;
		if (*(int*)&dword_5d4594_3798832 + *(int*)&dword_5d4594_3798816 >= 128) {
			return;
		}
		if ((int)v78 + nox_getBackbufHeight() > *(int*)&dword_5d4594_3798824 + *(int*)&dword_5d4594_3798808 + 46) {
			nox_xxx_tileDrawImpl_4826A0(vp);
			return;
		}
		++dword_5d4594_3798832;
		v23 = dword_5d4594_3798824 + 46;
		v19 = dword_5d4594_3798816 + v22;
		v24 = dword_5d4594_3798840 + 46;
		dword_5d4594_3798824 += 46;
		dword_5d4594_3798840 += 46;
		if (*(int*)&dword_5d4594_3798840 >= *(int*)&dword_5d4594_3798808) {
			dword_5d4594_3798840 = v24 - dword_5d4594_3798808;
		}
		v76 = v23 + dword_5d4594_3798808 - 46;
	} else {
		if (*(int*)&dword_5d4594_3798832 <= 0) {
			return;
		}
		if ((int)v78 < *(int*)&dword_5d4594_3798824 - 23) {
			nox_xxx_tileDrawImpl_4826A0(vp);
			return;
		}
		v19 = dword_5d4594_3798832 - 1;
		v20 = dword_5d4594_3798824 - 46;
		v21 = dword_5d4594_3798840 - 46;
		bool v5 = *(int*)&dword_5d4594_3798840 - 46 < 0;
		--dword_5d4594_3798832;
		dword_5d4594_3798824 -= 46;
		dword_5d4594_3798840 -= 46;
		if (v5) {
			dword_5d4594_3798840 = dword_5d4594_3798808 + v21;
		}
		v76 = v20;
	}
	sub_481410();
	v25 = dword_5d4594_3798828;
	v26 = dword_5d4594_3798820;
	int i = dword_5d4594_3798828;
	int j = dword_5d4594_3798812 + dword_5d4594_3798828 - 1;
	if (dword_5d4594_3798828 < j) {
		v27 = 44 * v19;
		int v71 = 44 * v19;
		while (1) {
			v28 = ptr_5D4594_2650668[v25];
			v29 = *(uint8_t*)(v28 + v27);
			v30 = v27 + v28;
			if (v29 & 2) {
				LOWORD(v27) = *(uint16_t*)(v30 + 24);
				v31 = *(unsigned short*)(v30 + 24);
				v64 = nox_tile_defs_arr[v31].data_32[*(uint32_t*)(v30 + 28) + nox_tile_defs_arr[v31].field_46];
				v68.field_0 = v26;
				v68.field_4 = v76 + 23;
				func_587000_154940(&v68, v64, v27);
				*getMemU32Ptr(0x85B3FC, 228 + 4 * v31) = 1;
				if (*(uint32_t*)(v30 + 40)) {
					nox_xxx_tileCallDrawEdges_481BF0((int)&v68, *(uint32_t*)(v30 + 40));
				}
			}
			if (*(uint8_t*)v30 & 1) {
				LOWORD(v27) = *(uint16_t*)(v30 + 4);
				v32 = *(unsigned short*)(v30 + 4);
				v65 = nox_tile_defs_arr[v32].data_32[*(uint32_t*)(v30 + 8) + nox_tile_defs_arr[v32].field_46];
				v68.field_0 = v26 + 23;
				v68.field_4 = v76;
				func_587000_154940(&v68, v65, v27);
				*getMemU32Ptr(0x85B3FC, 228 + 4 * v32) = 1;
				v33 = *(uint32_t*)(v30 + 20);
				if (v33) {
					nox_xxx_tileCallDrawEdges_481BF0((int)&v68, v33);
				}
			}
			v26 += 46;
			if (++i >= j) {
				break;
			}
			v27 = v71;
			v25 = i;
		}
	}
}

//----- (00482570) --------------------------------------------------------
int nox_xxx_tileCheckRedrawMB_482570(nox_draw_viewport_t* vp) {
	uint32_t* a1 = vp;
	int v1;       // esi
	int v2;       // ebx
	int v3;       // eax
	int v4;       // edx
	int i;        // ebp
	int v6;       // edi
	uint32_t* v7; // esi
	int v8;       // eax
	int v10;      // [esp+10h] [ebp-8h]
	int v11;      // [esp+14h] [ebp-4h]
	int v12;      // [esp+1Ch] [ebp+4h]

	v1 = a1[5] - a1[1];
	v2 = (a1[4] - *a1 - 11) / 46;
	if (v2 < 0) {
		v2 = 0;
	}
	v12 = dword_5d4594_3798812 + v2 - 1;
	if (v12 >= 128) {
		v12 = 127;
		v2 = 127 - dword_5d4594_3798812;
	}
	v3 = (v1 - 11) / 46 - 1;
	if (v3 < 0) {
		v3 = 0;
	}
	v4 = dword_5d4594_3798816 + v3;
	v10 = dword_5d4594_3798816 + v3;
	if (*(int*)&dword_5d4594_3798816 + v3 >= 128) {
		v10 = 127;
		v3 = 127 - dword_5d4594_3798816;
		v4 = 127;
	}
	v11 = v3;
	if (v3 >= v4) {
		return 0;
	}
	for (i = 44 * v3; 1; i += 44) {
		v6 = v2;
		if (v2 < v12) {
			v7 = &ptr_5D4594_2650668[v2];
			while (1) {
				v8 = i + *v7;
				if (*(uint8_t*)v8 & 1) {
					if ((nox_tile_defs_arr[*(uint32_t*)(v8 + 4)].field_58 & 1) == 1) {
						return 1;
					}
				}
				if (*(uint8_t*)v8 & 2 && (nox_tile_defs_arr[*(uint32_t*)(v8 + 24)].field_58 & 1) == 1) {
					return 1;
				}
				++v6;
				++v7;
				if (v6 >= v12) {
					v4 = v10;
					goto LABEL_19;
				}
			}
		}
	LABEL_19:
		if (++v11 >= v4) {
			return 0;
		}
	}
}

//----- (004826A0) --------------------------------------------------------
int nox_xxx_tileDrawImpl_4826A0(nox_draw_viewport_t* vp) {
	uint32_t* a1 = vp;
	int v1;     // esi
	int v2;     // ebx
	int result; // eax
	int v4;     // ecx
	int v5;     // esi
	int v6;     // ebx
	int v7;     // ecx
	int v8;     // ebp
	int v9;     // eax
	int v10;    // ecx
	int v11;    // esi
	char v12;   // al
	int v13;    // esi
	int v14;    // edi
	int v15;    // edi
	int v16;    // esi
	bool v17;   // zf
	int v18;    // [esp-8h] [ebp-34h]
	int v19;    // [esp-8h] [ebp-34h]
	int v20;    // [esp+10h] [ebp-1Ch]
	int v21;    // [esp+14h] [ebp-18h]
	int v22;    // [esp+18h] [ebp-14h]
	int v23;    // [esp+1Ch] [ebp-10h]
	int2 v24;   // [esp+24h] [ebp-8h]
	int v25;    // [esp+30h] [ebp+4h]

	v1 = a1[4] - *a1;
	v2 = a1[5] - a1[1];
	dword_5d4594_3798836 = 0;
	dword_5d4594_3798840 = 0;
	sub_481410();
	memset(getMemAt(0x85B3FC, 228), 0, 0x2C0u);
	memset(getMemAt(0x5D4594, 2523980), 0, 0x100u);
	v25 = (v1 - 11) / 46;
	if (v25 < 0) {
		v25 = 0;
	}
	v20 = dword_5d4594_3798812 + v25 - 1;
	if (v20 >= 128) {
		v20 = 127;
		v25 = 127 - dword_5d4594_3798812;
	}
	result = (v2 - 11) / 46 - 1;
	if (result < 0) {
		result = 0;
	}
	v4 = dword_5d4594_3798816 + result;
	if (*(int*)&dword_5d4594_3798816 + result >= 128) {
		v4 = 127;
		result = 127 - dword_5d4594_3798816;
	}
	v5 = v25;
	dword_5d4594_3798824 = 46 * result - 11;
	v6 = 46 * result - 57;
	dword_5d4594_3798828 = v25;
	dword_5d4594_3798820 = 46 * v25 - 11;
	dword_5d4594_3798832 = result;
	if (result < v4) {
		v21 = 44 * result;
		v23 = v4 - result;
		v7 = v20;
		do {
			v8 = 46 * v25 - 57;
			v6 += 46;
			v9 = v5;
			v22 = v5;
			if (v5 < v7) {
				do {
					HIWORD(v10) = HIWORD(ptr_5D4594_2650668); // TODO: why it's doing it?
					v8 += 46;
					v11 = ptr_5D4594_2650668[v9];
					v12 = *(uint8_t*)(v11 + v21);
					v13 = v21 + v11;
					if (v12) {
						if (v12 & 2) {
							LOWORD(v10) = *(uint16_t*)(v13 + 24);
							v14 = *(unsigned short*)(v13 + 24);
							v18 = nox_tile_defs_arr[v14].data_32[*(uint32_t*)(v13 + 28) + nox_tile_defs_arr[v14].field_46];
							v24.field_0 = v8;
							v24.field_4 = v6 + 23;
							func_587000_154940(&v24, v18, v10);
							*getMemU32Ptr(0x85B3FC, 228 + 4 * v14) = 1;
							if (*(uint32_t*)(v13 + 40)) {
								nox_xxx_tileCallDrawEdges_481BF0((int)&v24, *(uint32_t*)(v13 + 40));
							}
						}
						if (*(uint8_t*)v13 & 1) {
							LOWORD(v10) = *(uint16_t*)(v13 + 4);
							v15 = *(unsigned short*)(v13 + 4);
							v19 = nox_tile_defs_arr[v15].data_32[*(uint32_t*)(v13 + 8) + nox_tile_defs_arr[v15].field_46];
							v24.field_0 = v8 + 23;
							v24.field_4 = v6;
							func_587000_154940(&v24, v19, v10);
							*getMemU32Ptr(0x85B3FC, 228 + 4 * v15) = 1;
							v16 = *(uint32_t*)(v13 + 20);
							if (v16) {
								nox_xxx_tileCallDrawEdges_481BF0((int)&v24, v16);
							}
						}
					}
					v7 = v20;
					v9 = ++v22;
				} while (v22 < v20);
				v5 = v25;
			}
			result = v23 - 1;
			v17 = v23 == 1;
			v21 += 44;
			--v23;
		} while (!v17);
	}
	return result;
}
// 4828A6: variable 'v10' is possibly undefined

//----- (00485B30) --------------------------------------------------------
int nox_thing_read_floor_485B30(nox_memfile* f, char* a2) {
	int a1 = f;
	int v2;            // esi
	uint8_t* v3;       // edi
	uint8_t* v9;       // edx
	int v10;           // edi
	int* v12;          // eax
	int v13;           // ecx
	char* v14;         // eax
	char v15;          // cl
	int v16;           // [esp-4h] [ebp-40h]
	unsigned char v17; // [esp+10h] [ebp-2Ch]
	int i;             // [esp+10h] [ebp-2Ch]
	unsigned char v19; // [esp+14h] [ebp-28h]
	const char* v21;   // [esp+18h] [ebp-24h]
	char v22[32];      // [esp+1Ch] [ebp-20h]
	unsigned char v23; // [esp+40h] [ebp+4h]

	v2 = a1;
	v16 = a1;
	v3 = (uint8_t*)(*(uint32_t*)(a1 + 8) + 4);
	*(uint32_t*)(a1 + 8) = v3;
	LOBYTE(a1) = *v3;
	*(uint32_t*)(v2 + 8) = v3 + 1;
	nox_memfile_read(v22, 1u, (unsigned char)a1, v16);
	v22[(unsigned char)a1] = 0;
	int v7 = a1;
	if (nox_tile_def_cnt > 0) {
		int v5 = 0;
		for (v5 = 0; v5 < nox_tile_def_cnt; v5++) {
			nox_tileDef_t* p = &nox_tile_defs_arr[v5];
			if (strcmp(&p->name[0], v22) == 0) {
				v7 = v5;
				break;
			}
		}
		if (v5 == nox_tile_def_cnt) {
			return 0;
		}
	}
	v9 = (uint8_t*)(*(uint32_t*)(v2 + 8) + 12);
	*(uint32_t*)(v2 + 8) = v9;
	LOBYTE(v21) = *v9;
	*(uint32_t*)(v2 + 8) = v9 + 1;
	v19 = v9[1];
	*(uint32_t*)(v2 + 8) = v9 + 2;
	v17 = v9[2];
	*(uint32_t*)(v2 + 8) = v9 + 4;
	v10 = (unsigned char)v21 * v19 * v17;
	nox_tile_defs_arr[v7].data_32 = calloc(v10, 4);
	int v11 = 0;
	for (i = 0; v11 < v10; ) {
		v12 = *(int**)(v2 + 8);
		v13 = *v12;
		*(uint32_t*)(v2 + 8) = v12 + 1;
		*a2 = getMemByte(0x5D4594, 1193192);
		if (v13 == -1) {
			v14 = *(char**)(v2 + 8);
			v15 = *v14++;
			*(uint32_t*)(v2 + 8) = v14;
			LOBYTE(v21) = v15;
			v23 = *v14;
			*(uint32_t*)(v2 + 8) = v14 + 1;
			nox_memfile_read(a2, 1u, v23, v2);
			v13 = -1;
			a2[v23] = 0;
			v11 = i;
		}
		nox_tile_defs_arr[v7].data_32[v11] = nox_xxx_readImgMB_42FAA0(v13, v21, a2);
		v11++;
		i = v11;
	}
	return 1;
}
// 485CB9: variable 'v21' is possibly undefined
// 485B30: using guessed type char var_20[32];

//----- (00485D40) --------------------------------------------------------
int nox_thing_read_edge_485D40(nox_memfile* f, char* a2) {
	int a1 = f;
	int v2;            // esi
	uint8_t* v3;       // edi
	int v4;            // eax
	int v5;            // ebp
	const char* v6;    // ebx
	int v7;            // ebx
	int result;        // eax
	uint8_t* v9;       // edi
	char v10;          // cl
	unsigned char v11; // dl
	int v12;           // eax
	unsigned char v13; // di
	int v14;           // edi
	int v15;           // ebx
	int* v16;          // eax
	int v17;           // ecx
	char* v18;         // eax
	char v19;          // cl
	int* v20;          // eax
	int v21;           // ecx
	int v22;           // [esp-4h] [ebp-40h]
	int i;             // [esp+10h] [ebp-2Ch]
	int v24;           // [esp+14h] [ebp-28h]
	unsigned int v25;  // [esp+18h] [ebp-24h]
	char v26[64];      // [esp+1Ch] [ebp-20h]
	unsigned char v27; // [esp+40h] [ebp+4h]

	v2 = a1;
	v22 = a1;
	v3 = (uint8_t*)(*(uint32_t*)(a1 + 8) + 4);
	*(uint32_t*)(a1 + 8) = v3;
	LOBYTE(a1) = *v3;
	*(uint32_t*)(v2 + 8) = v3 + 1;
	nox_memfile_read(v26, 1u, (unsigned char)a1, v22);
	v4 = dword_5d4594_251572;
	v5 = 0;
	v26[(unsigned char)a1] = 0;
	if (v4 <= 0) {
		v7 = a1;
	} else {
		v6 = (const char*)getMemAt(0x85B3FC, 28644);
		while (1) {
			if (!strcmp(v6, v26)) {
				v7 = v5;
				break;
			}
			++v5;
			v6 += 60;
			if (v5 >= *(int*)&dword_5d4594_251572) {
				v7 = a1;
				break;
			}
		}
	}
	if (v5 == dword_5d4594_251572) {
		return 0;
	}
	v9 = (uint8_t*)(*(uint32_t*)(v2 + 8) + 9);
	*(uint32_t*)(v2 + 8) = v9;
	v25 = *v9;
	*(uint32_t*)(v2 + 8) = v9 + 2;
	v10 = v9[2];
	*(uint32_t*)(v2 + 8) = v9 + 3;
	if (v10 == 1) {
		return 0;
	}
	v11 = v9[3];
	v12 = (int)(v9 + 4);
	*(uint32_t*)(v2 + 8) = v9 + 4;
	v13 = v9[4];
	*(uint32_t*)(v2 + 8) = v12 + 1;
	v14 = 2 * v25 * (v11 + v13);
	result = (int)calloc(v14, 5);
	memset((int*)result, 0, 5 * v14);
	v24 = 15 * v7;
	*getMemU32Ptr(0x85B3FC, 28676 + 60 * v7) = result;
	if (result) {
		v15 = 0;
		for (i = 0; v15 < v14; *(uint32_t*)(*getMemU32Ptr(0x85B3FC, 28676 + 4 * v24) + 4 * v15 - 4) =
								   nox_xxx_readImgMB_42FAA0(v17, v25, a2)) {
			v16 = *(int**)(v2 + 8);
			v17 = *v16;
			*(uint32_t*)(v2 + 8) = v16 + 1;
			*a2 = getMemByte(0x5D4594, 1193196);
			if (v17 == -1) {
				v18 = *(char**)(v2 + 8);
				v19 = *v18++;
				*(uint32_t*)(v2 + 8) = v18;
				LOBYTE(v25) = v19;
				v27 = *v18;
				*(uint32_t*)(v2 + 8) = v18 + 1;
				nox_memfile_read(a2, 1u, v27, v2);
				v17 = -1;
				a2[v27] = 0;
				v15 = i;
			}
			i = ++v15;
		}
		v20 = *(int**)(v2 + 8);
		v21 = *v20;
		*(uint32_t*)(v2 + 8) = v20 + 1;
		result = v21 == 1162757152;
	}
	return result;
}
// 485EEB: variable 'v25' is possibly undefined
// 485D40: using guessed type char var_20[32];

//----- (00486640) --------------------------------------------------------
int sub_486640(void* a1p, int a2) {
	int a1 = a1p;
	return a2 * (*(uint32_t*)(a1 + 36) >> 16) / 100;
}

//----- (004866D0) --------------------------------------------------------
int sub_4866D0(uint32_t* a1, int a2) { return *a1 + 36 * a2; }

//----- (00486A10) --------------------------------------------------------
unsigned int sub_486A10(int a1, void* a2) {
	void* v2;            // eax
	unsigned int result; // eax

	v2 = bsearch(a2, *(const void**)a1, *(uint32_t*)(a1 + 4), 0x24u, (int (*)(const void*, const void*))nox_strcmpi);
	if (v2) {
		result = ((unsigned int)v2 - *(uint32_t*)a1) / 0x24;
	} else {
		result = -1;
	}
	return result;
}

//----- (00486AA0) --------------------------------------------------------
int sub_486AA0(uint32_t* a1, int a2, uint32_t* a3) {
	uint32_t* v3; // eax
	int result;   // eax

	v3 = (uint32_t*)sub_4866D0(a1, a2);
	*a3 = 4;
	a3[2] = v3[6];
	a3[3] = ((v3[7] & 1) != 0) + 1;
	a3[6] = v3[8];
	if (v3[7] & 8) {
		result = 2;
		a3[1] = 2;
		a3[4] = 2;
	} else {
		a3[1] = 0;
		result = ((v3[7] & 4) != 0) + 1;
		a3[4] = result;
	}
	return result;
}

//----- (00486B60) --------------------------------------------------------
int sub_486B60(int a1, int a2) {
	int v2;       // ebp
	FILE* v3;     // eax
	FILE* v6;     // eax
	FILE* v7;     // esi
	int v8;       // edi
	int v9;       // eax
	int v10;      // ecx
	int v12;      // [esp+10h] [ebp-12Ch]
	char v13[8];  // [esp+14h] [ebp-128h]
	char v14[16]; // [esp+1Ch] [ebp-120h]
	int v15[12];  // [esp+2Ch] [ebp-110h]

	v12 = 1;
	v2 = sub_4866D0((uint32_t*)a1, a2);
	sub_486E00(a1);
	v3 = *(FILE**)(a1 + 268);
	*(uint32_t*)(a1 + 280) = v3;
	*(uint32_t*)(a1 + 284) = *(uint32_t*)(v2 + 20);
	if (nox_fs_fseek_start(v3, *(uint32_t*)(v2 + 16))) {
		v12 = 0;
	}
	if (!*(uint32_t*)(v2 + 20)) {
		v12 = 0;
	}
	if (!*(uint32_t*)(a1 + 276)) {
		return v12;
	}
	strcpy((char*)&v15[3], (const char*)(a1 + 8));
	strcat((char*)&v15[3], (const char*)v2);
	strcat((char*)&v15[3], ".wav");
	v6 = nox_fs_open((const char*)&v15[3]);
	v7 = v6;
	v8 = 0;
	*(uint32_t*)(a1 + 272) = v6;
	if (!v6) {
		return v12;
	}
	if (nox_binfile_fread_raw_40ADD0((char*)v15, 0xCu, 1u, v6) != 1 || v15[0] != 1179011410 || v15[2] != 1163280727) {
		printf("error: '%s' is bad - cannot read\n", &v15[3]);
		if (*(uint32_t*)(a1 + 272)) {
			nox_fs_close(*(FILE**)(a1 + 272));
			*(uint32_t*)(a1 + 272) = 0;
		}
		return v12;
	}
	if (nox_binfile_fread_raw_40ADD0(v13, 8u, 1u, v7) != 1) {
		goto LABEL_18;
	}
	while (1) {
		if (*(uint32_t*)v13 == 544501094) {
			nox_binfile_fread_raw_40ADD0(v14, 0x10u, 1u, v7);
			nox_fs_fseek_cur(v7, *(uint32_t*)&v13[4] - 16);
			goto LABEL_15;
		}
		if (*(uint32_t*)v13 == 1635017060) {
			break;
		}
		nox_fs_fseek_cur(v7, *(int*)&v13[4]);
	LABEL_15:
		if (nox_binfile_fread_raw_40ADD0(v13, 8u, 1u, v7) != 1) {
			goto LABEL_18;
		}
	}
	v8 = *(uint32_t*)&v13[4];
LABEL_18:
	*(uint32_t*)(v2 + 28) = 2;
	if (*(unsigned short*)&v14[12] / (int)*(unsigned short*)&v14[2] == 2) {
		*(uint32_t*)(v2 + 28) = 6;
	}
	if (*(uint16_t*)&v14[2] == 2) {
		v9 = *(uint32_t*)(v2 + 28);
		LOBYTE(v9) = v9 | 1;
		*(uint32_t*)(v2 + 28) = v9;
	}
	*(uint32_t*)(v2 + 24) = *(uint32_t*)&v14[4];
	v10 = *(uint32_t*)(a1 + 272);
	*(uint32_t*)(a1 + 284) = v8;
	*(uint32_t*)(a1 + 280) = v10;
	return 1;
}

//----- (00486DB0) --------------------------------------------------------
signed int sub_486DB0(int a1, char* a2, signed int a3) {
	signed int result; // eax
	signed int v4;     // eax

	if (!*(uint32_t*)(a1 + 280)) {
		return 0;
	}
	v4 = a3;
	if (a3 > *(int*)(a1 + 284)) {
		v4 = *(uint32_t*)(a1 + 284);
	}
	if (v4 <= 0 || (result = nox_binfile_fread_raw_40ADD0(a2, 1u, v4, *(FILE**)(a1 + 280)), result < 0)) {
		result = 0;
	}
	*(uint32_t*)(a1 + 284) -= result;
	return result;
}

//----- (00486E00) --------------------------------------------------------
FILE* sub_486E00(int a1) {
	FILE* result; // eax

	result = *(FILE**)(a1 + 272);
	*(uint32_t*)(a1 + 280) = 0;
	if (result) {
		nox_fs_close(result);
		result = 0;
		*(uint32_t*)(a1 + 272) = 0;
	}
	return result;
}

//----- (00486E30) --------------------------------------------------------
int sub_486E30(int a1, uint32_t* a2) {
	int result; // eax

	a2[33] = a1;
	++*(uint32_t*)(a1 + 192);
	++*(uint32_t*)(a1 + 212);
	nox_common_list_append_4258E0(a1 + 200, a2);
	result = *(uint32_t*)(a1 + 212) - 1;
	*(uint32_t*)(a1 + 212) = result;
	if (result < 0) {
		*(uint32_t*)(a1 + 212) = 0;
	}
	return result;
}

//----- (00486E90) --------------------------------------------------------
int sub_486E90(int a1) {
	int v1;     // esi
	int result; // eax

	v1 = *(uint32_t*)(a1 + 132);
	nox_common_list_remove_425920((uint32_t**)a1);
	--*(uint32_t*)(v1 + 192);
	++*(uint32_t*)(v1 + 212);
	nox_common_list_remove_425920((uint32_t**)a1);
	result = *(uint32_t*)(v1 + 212) - 1;
	*(uint32_t*)(v1 + 212) = result;
	if (result < 0) {
		*(uint32_t*)(v1 + 212) = 0;
	}
	return result;
}

//----- (00486FA0) --------------------------------------------------------
uint32_t* sub_486FA0(int a1) {
	uint32_t* result; // eax
	uint32_t* v2;     // edi
	int v3;           // eax

	result = sub_486FE0(a1);
	v2 = result;
	if (result) {
		v3 = *(uint32_t*)(a1 + 12);
		LOBYTE(v3) = v3 | 1;
		*(uint32_t*)(a1 + 12) = v3;
		sub_487050(v2);
		if (*(uint8_t*)(a1 + 8) & 2) {
			*getMemU32Ptr(0x5D4594, 1193332) = 1;
		}
		result = v2;
	}
	return result;
}

//----- (00486FE0) --------------------------------------------------------
uint32_t* sub_486FE0(int a1) {
	uint32_t* v1; // esi

	v1 = calloc(1, 0x58u);
	memset(v1, 0, 0x58u);
	sub_425770(v1);
	v1[4] = 0;
	v1[3] = a1;
	if (!(*(int (**)(uint32_t*))(a1 + 20))(v1)) {
		return v1;
	}
	if (v1) {
		sub_487030(v1);
	}
	return 0;
}

//----- (00487030) --------------------------------------------------------
void sub_487030(void* lpMem) {
	(*(void (**)(void*))(*((uint32_t*)lpMem + 3) + 24))(lpMem);
	*(uint32_t*)(*((uint32_t*)lpMem + 3) + 12) &= 0xFFFFFFFE;
	free(lpMem);
}

//----- (00487050) --------------------------------------------------------
void sub_487050(uint32_t* a1) { nox_common_list_append_4258E0(*(int*)&dword_587000_155144, a1); }

//----- (00487070) --------------------------------------------------------
void sub_487070(void* lpMem) {
	sub_487090((uint32_t**)lpMem);
	sub_487030(lpMem);
	*getMemU32Ptr(0x5D4594, 1193332) = 0;
}

//----- (00487090) --------------------------------------------------------
void sub_487090(uint32_t** a1) { nox_common_list_remove_425920(a1); }

//----- (004870A0) --------------------------------------------------------
void sub_4870A0() {
	int* v1; // edi
	int* v2; // esi
	int* v3; // [esp+4h] [ebp-4h]

	v1 = sub_4870E0((int*)&v3);
	if (v1) {
		do {
			v2 = sub_487100(&v3);
			sub_487070(v1);
			v1 = v2;
		} while (v2);
	}
}

//----- (004870E0) --------------------------------------------------------
int* sub_4870E0(int* a1) {
	int* result; // eax

	result = nox_common_list_getFirstSafe_425890(*(int**)&dword_587000_155144);
	*a1 = (int)result;
	return result;
}

//----- (00487100) --------------------------------------------------------
int* sub_487100(int** a1) {
	if (*a1) {
		*a1 = nox_common_list_getNextSafe_4258A0(*a1);
	}
	return *a1;
}

//----- (00487150) --------------------------------------------------------
uint32_t* sub_487150(int a1, const void* a2) {
	int v2;       // edi
	uint32_t* v3; // esi
	uint32_t* v4; // eax
	int v6;       // [esp+8h] [ebp-4h]

	v2 = a1;
	if (a1 == -1) {
		v2 = 0;
	}
	sub_487360(v2, (int**)&a1, &v6);
	if (!a1) {
		return 0;
	}
	v3 = *(uint32_t**)(a1 + 4 * v6 + 24);
	if (!v3) {
		v4 = sub_4871C0(a1, v6, a2);
		v3 = v4;
		if (!v4) {
			return 0;
		}
		v4[47] = v2;
		sub_487310(v4);
	}
	++v3[4];
	return v3;
}

//----- (004871C0) --------------------------------------------------------
uint32_t* sub_4871C0(int a1, int a2, const void* a3) {
	int v3;       // ebp
	uint32_t* v4; // esi

	v3 = *(uint32_t*)(a1 + 12);
	v4 = calloc(1, 0x108u);
	memset(v4, 0, 0x108u);
	sub_425770(v4);
	v4[6] = a2;
	v4[5] = a1;
	v4[4] = 0;
	++*(uint32_t*)(a1 + 16);
	*(uint32_t*)(a1 + 4 * a2 + 24) = v4;
	v4[64] = *(uint32_t*)(v3 + 36);
	nox_common_list_clear_425760(v4 + 50);
	sub_4864A0(v4 + 22);
	v4[53] = 0;
	v4[56] = 33;
	v4[60] = 0;
	v4[58] = 0;
	v4[62] = 0;
	v4[54] = sub_4873C0;
	v4[57] = 0;
	v4[61] = 0;
	v4[59] = 0;
	v4[63] = 0;
	nullsub_10(v4 + 15);
	nullsub_10(v4 + 8);
	if (a3) {
		sub_487590((int)v4, a3);
	}
	if (!(*(int (**)(uint32_t*))(v3 + 28))(v4)) {
		return v4;
	}
	if (v4) {
		sub_4872C0(v4);
	}
	return 0;
}
// 487CF0: using guessed type void  nullsub_10(uint32_t);

//----- (004872C0) --------------------------------------------------------
void sub_4872C0(void* lpMem) {
	int v1; // eax
	int v2; // ecx

	sub_487910((int)lpMem, -1);
	(*(void (**)(void*))(*(uint32_t*)(*((uint32_t*)lpMem + 5) + 12) + 32))(lpMem);
	*(uint32_t*)(*((uint32_t*)lpMem + 5) + 4 * *((uint32_t*)lpMem + 6) + 24) = 0;
	v1 = *((uint32_t*)lpMem + 5);
	v2 = *(uint32_t*)(v1 + 16) - 1;
	*(uint32_t*)(v1 + 16) = v2;
	if (v2 < 0) {
		*(uint32_t*)(*((uint32_t*)lpMem + 5) + 16) = 0;
	}
	free(lpMem);
}

//----- (00487310) --------------------------------------------------------
int sub_487310(uint32_t* a1) {
	int result; // eax

	++*(uint32_t*)((uint32_t)dword_587000_155144 + 24);
	nox_common_list_append_4258E0((uint32_t)dword_587000_155144 + 12, a1);
	result = *(uint32_t*)((uint32_t)dword_587000_155144 + 24) - 1;
	*(uint32_t*)((uint32_t)dword_587000_155144 + 24) = result;
	if (result < 0) {
		*(uint32_t*)((uint32_t)dword_587000_155144 + 24) = 0;
	}
	return result;
}

//----- (00487360) --------------------------------------------------------
int* sub_487360(int a1, int** a2, int* a3) {
	int* result; // eax
	int i;       // esi
	int v5;      // ecx
	int* v6;     // [esp+4h] [ebp-4h]

	result = sub_4870E0((int*)&v6);
	for (i = a1; result; result = sub_487100(&v6)) {
		v5 = result[5];
		if (i < v5) {
			break;
		}
		i -= v5;
	}
	*a2 = result;
	if (result) {
		result = a3;
		*a3 = i;
	} else {
		*a3 = -1;
	}
	return result;
}

//----- (004873C0) --------------------------------------------------------
int sub_4873C0(int a3) {
	int v1;           // esi
	long long v3;     // rax
	unsigned int v4;  // ecx
	int v5;           // ebp
	bool v6;          // cf
	unsigned int v7;  // ebx
	unsigned int v8;  // ecx
	int v9;           // edi
	unsigned int v10; // eax
	uint32_t* v11;    // ebx
	int v12;          // edi
	int v13;          // ebp
	uint32_t* v14;    // esi
	int v15;          // [esp+10h] [ebp-Ch]
	int v16;          // [esp+18h] [ebp-4h]
	int v17;          // [esp+20h] [ebp+4h]

	v1 = a3;
	if (*(uint32_t*)(a3 + 212)) {
		return -2146304000;
	}
	v3 = nox_platform_get_ticks();
	v4 = *(uint32_t*)(a3 + 248);
	v5 = v3;
	v6 = (unsigned int)v3 < v4;
	v7 = v3 - v4;
	v8 = *(uint32_t*)(a3 + 224);
	v16 = HIDWORD(v3);
	v9 = HIDWORD(v3) - (v6 + *(uint32_t*)(a3 + 252));
	v10 = *(uint32_t*)(a3 + 228);
	if (__PAIR64__(v9, v7) >= __PAIR64__(v10, v8)) {
		*(uint32_t*)(a3 + 232) = v7;
		*(uint32_t*)(a3 + 236) = v9;
		if (*(uint64_t*)(a3 + 240) > 10 * __PAIR64__(v10, v8)) {
			*(uint32_t*)(a3 + 240) = 0;
			*(uint32_t*)(a3 + 244) = 0;
		}
		if (__PAIR64__(v9, v7) > *(uint64_t*)(a3 + 240)) {
			*(uint32_t*)(a3 + 240) = v7;
			*(uint32_t*)(a3 + 244) = v9;
		}
		v11 = (uint32_t*)(a3 + 88);
		v15 = a3 + 88;
		sub_486520((unsigned int*)(a3 + 88));
		if (*(uint32_t*)(a3 + 184)) {
			sub_486520(*(unsigned int**)(a3 + 184));
		}
		if (!*(uint32_t*)(a3 + 184) || !(v17 = sub_486550(*(uint8_t**)(a3 + 184)))) {
			v17 = sub_486550((uint8_t*)(v1 + 88));
		}
		*(uint32_t*)(v1 + 248) = v5;
		*(uint32_t*)(v1 + 252) = v16;
		v12 = *(uint32_t*)(v1 + 200);
		if (v12 != v1 + 200) {
			do {
				v13 = *(uint32_t*)v12;
				if (*(uint8_t*)(v12 + 124) & 1 && *(uint32_t*)(v12 + 288)) {
					if ((sub_486520((unsigned int*)(v12 + 16)), v17) || sub_486550((uint8_t*)(v12 + 16)) ||
						*(uint32_t*)(v12 + 116) && sub_486550(*(uint8_t**)(v12 + 116)) ||
						*(uint32_t*)(v12 + 112) && sub_486550(*(uint8_t**)(v12 + 112))) {
						sub_4BD840(v12);
						(*(void (**)(int))(*(uint32_t*)(v12 + 172) + 32))(v12);
					}
				}
				v12 = v13;
			} while (v13 != v1 + 200);
			v11 = (uint32_t*)v15;
		}
		v14 = *(uint32_t**)(v1 + 184);
		if (v14) {
			sub_486620(v14);
		}
		sub_486620(v11);
	}
	return 0;
}

//----- (00487590) --------------------------------------------------------
int sub_487590(int a1, const void* a2) {
	int result; // eax

	result = a1;
	memcpy((void*)(a1 + 60), a2, 0x1Cu);
	return result;
}

//----- (004875B0) --------------------------------------------------------
int* sub_4875B0(int* a1) {
	int* result; // eax

	result = nox_common_list_getFirstSafe_425890((int*)((uint32_t)dword_587000_155144 + 12));
	*a1 = (int)result;
	return result;
}

//----- (004875D0) --------------------------------------------------------
int* sub_4875D0(int** a1) {
	if (*a1) {
		*a1 = nox_common_list_getNextSafe_4258A0(*a1);
	}
	return *a1;
}

//----- (004875F0) --------------------------------------------------------
int sub_4875F0() {
	int* v0;    // edi
	int* v1;    // esi
	int result; // eax
	int* v3;    // [esp+4h] [ebp-4h]

	++*(uint32_t*)((uint32_t)dword_587000_155144 + 24);
	v0 = sub_4875B0((int*)&v3);
	if (v0) {
		do {
			v1 = sub_4875D0(&v3);
			sub_487680(v0);
			v0 = v1;
		} while (v1);
	}
	result = *(uint32_t*)((uint32_t)dword_587000_155144 + 24) - 1;
	*(uint32_t*)((uint32_t)dword_587000_155144 + 24) = result;
	if (result < 0) {
		*(uint32_t*)((uint32_t)dword_587000_155144 + 24) = 0;
	}
	return result;
}

//----- (00487680) --------------------------------------------------------
void sub_487680(void* lpMem) {
	sub_4876A0((uint32_t**)lpMem);
	sub_4872C0(lpMem);
}

//----- (004876A0) --------------------------------------------------------
void* sub_4876A0(uint32_t** a1) {
	void* result; // eax

	++*(uint32_t*)((uint32_t)dword_587000_155144 + 24);
	nox_common_list_remove_425920(a1);
	result = (void*)(*(uint32_t*)((uint32_t)dword_587000_155144 + 24) - 1);
	*(uint32_t*)((uint32_t)dword_587000_155144 + 24) = result;
	if ((int)result < 0) {
		result = *(void**)&dword_587000_155144;
		*(uint32_t*)((uint32_t)dword_587000_155144 + 24) = 0;
	}
	return result;
}

//----- (00487750) --------------------------------------------------------
uint32_t* sub_487750(int a1) {
	uint32_t* v1; // eax
	uint32_t* v2; // esi

	if (*(int*)(a1 + 192) >= *(int*)(a1 + 196)) {
		return 0;
	}
	v1 = sub_4BD720(a1);
	v2 = v1;
	if (!v1) {
		return 0;
	}
	sub_486E30(a1, v1);
	return v2;
}

//----- (00487790) --------------------------------------------------------
int sub_487790(int a1, int a2) {
	int v2; // esi
	int v3; // edi

	v2 = 0;
	if (sub_487750(a1)) {
		v3 = a2;
		do {
			++v2;
			--v3;
		} while (v3 && sub_487750(a1));
	}
	return v2;
}

//----- (004877D0) --------------------------------------------------------
int* sub_4877D0(int a1, int* a2) {
	int* result; // eax

	result = nox_common_list_getFirstSafe_425890((int*)(a1 + 200));
	*a2 = (int)result;
	return result;
}

//----- (004877F0) --------------------------------------------------------
int* sub_4877F0(int** a1) {
	if (*a1) {
		*a1 = nox_common_list_getNextSafe_4258A0(*a1);
	}
	return *a1;
}

//----- (00487810) --------------------------------------------------------
int* sub_487810(int a1, int a2) {
	unsigned int v2; // esi
	int v3;          // edi
	int* v4;         // ebp
	int* result;     // eax
	int v6;          // ecx
	unsigned int v7; // edx
	int v8;          // [esp+10h] [ebp-8h]
	int* v9;         // [esp+14h] [ebp-4h]

	v2 = -1;
	if (a2 == -1) {
		a2 = 1;
	}
	v3 = 127;
	v4 = 0;
	v8 = 127;
	v9 = 0;
	for (result = sub_4877D0(a1, &a1); result; result = sub_4877F0((int**)&a1)) {
		if (result[3] == a2) {
			if (!(result[31] & 0x15)) {
				return result;
			}
			v6 = result[30];
			if (result[31] & 1) {
				if (v6 >= v3) {
					if (v6 == v3) {
						v7 = result[45];
						if (v7 < v2 && v2 - v7 >= 0x666) {
							v3 = result[30];
							v4 = result;
							v2 = result[45];
						}
					}
				} else {
					v2 = result[45];
					v3 = result[30];
					v4 = result;
				}
			} else if (v6 < v8) {
				v8 = result[30];
				v9 = result;
			}
		}
	}
	result = v9;
	if (!v9 || v8 > v3) {
		result = v4;
	}
	return result;
}

//----- (00487910) --------------------------------------------------------
int sub_487910(int a1, int a2) {
	int* v2; // edi
	int v3;  // ebx
	int* v4; // esi

	v2 = sub_4877D0(a1, &a1);
	if (!v2) {
		return 0;
	}
	v3 = a2;
	do {
		v4 = sub_4877F0((int**)&a1);
		if (v3 == -1 || v2[3] == v3) {
			sub_4BDA60(v2);
		}
		v2 = v4;
	} while (v4);
	return 0;
}

//----- (00487970) --------------------------------------------------------
int* sub_487970(int a1, int a2) {
	int* result; // eax
	int* v3;     // edi
	int v4;      // ebx
	int* v5;     // esi

	result = sub_4877D0(a1, &a1);
	v3 = result;
	if (result) {
		v4 = a2;
		do {
			result = sub_4877F0((int**)&a1);
			v5 = result;
			if (v4 == -1 || v3[3] == v4) {
				result = (int*)sub_4BDA80((int)v3);
			}
			v3 = v5;
		} while (v5);
	}
	return result;
}

//----- (00487C30) --------------------------------------------------------
void sub_487C30(uint32_t* a1) {
	*a1 = 0;
	a1[1] = 0;
	a1[5] = 0;
	a1[6] = 0;
	nox_common_list_clear_425760(a1 + 2);
}

//----- (00487C50) --------------------------------------------------------
int sub_487C50(int a1, uint32_t* a2) {
	int result; // eax

	nox_common_list_append_4258E0(a1 + 8, a2);
	result = a2[4] + *(uint32_t*)(a1 + 4);
	*(uint32_t*)(a1 + 4) = result;
	a2[5] = a1;
	return result;
}

//----- (00487C80) --------------------------------------------------------
int sub_487C80(int a1) { return nox_common_list_getNext_425940((int*)(a1 + 8)); }

//----- (00487D00) --------------------------------------------------------
int sub_487D00(uint32_t* a1) {
	int v1;     // edx
	int result; // eax

	v1 = a1[1];
	result = a1[2] * a1[3] * a1[4];
	a1[5] = result;
	if (v1 == 1) {
		result >>= 2;
		a1[5] = result;
	}
	return result;
}

//----- (00487D30) --------------------------------------------------------
uint32_t* sub_487D30(uint32_t* a1, int a2, int a3) {
	uint32_t* result; // eax

	a1[4] = a3;
	a1[3] = a2;
	result = sub_425770(a1);
	a1[5] = 0;
	return result;
}

//----- (00487D60) --------------------------------------------------------
int sub_487D60(int a1) {
	int result; // eax

	result = a1;
	*(uint32_t*)(a1 + 20) = 0;
	return result;
}

//----- (004896E0) --------------------------------------------------------
int sub_4896E0() {
	if (dword_5d4594_1193360) {
		free(*(void**)&dword_5d4594_1193360);
	}
	return 1;
}

//----- (00489870) --------------------------------------------------------
int sub_489870() {
	int v0;            // eax
	unsigned char* v1; // esi
	uint32_t* v2;      // eax
	const wchar2_t* v3; // eax
	unsigned int v4;   // eax
	int v5;            // edx
	char v6;           // cl
	int v7;            // eax

	v0 = 0;
	v1 = getMemAt(0x5D4594, 1193388 + 44 * v0);
	if (*getMemU32Ptr(0x5D4594, 1193372 + 4 * v0) == 2) {
		*(uint32_t*)v1 =
			(nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10028)->draw_data.field_0 >> 2) & 1;
		v2 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10031);
		v3 = (const wchar2_t*)nox_window_call_field_94((int)v2, 16413, 0, 0);
		*((uint32_t*)v1 + 4) = nox_wcstol(v3, 0, 10);
		*((uint32_t*)v1 + 1) =
			(nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10029)->draw_data.field_0 >> 2) & 1;
		v4 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10030)->draw_data.field_0;
		*((uint32_t*)v1 + 3) = 0;
		*((uint32_t*)v1 + 2) = (v4 >> 2) & 1;
		if (nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10015)->draw_data.field_0 & 4) {
			v5 = *((uint32_t*)v1 + 3);
			LOBYTE(v5) = v5 | 0x80;
			*((uint32_t*)v1 + 3) = v5;
			v6 = *((uint8_t*)(&nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10016)
								   ->draw_data.field_0));
			v7 = *((uint32_t*)v1 + 3);
			if (v6 & 4) {
				LOBYTE(v7) = v7 | 1;
			} else {
				LOBYTE(v7) = v7 | 2;
			}
			*((uint32_t*)v1 + 3) = v7;
		}
		*((uint32_t*)v1 + 5) =
			(nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10014)->draw_data.field_0 >> 2) & 1;
		*((uint32_t*)v1 + 10) =
			(nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10018)->draw_data.field_0 >> 2) & 1;
	}
	return nox_window_set_hidden(*(int*)&dword_5d4594_1193380, 1);
}

//----- (004899C0) --------------------------------------------------------
int nox_xxx_checkSomeFlagsOnJoin_4899C0(nox_gui_server_ent_t* srv) {
	int a1 = srv;
	int v1;            // eax
	int v2;            // edx
	int v3;            // eax
	unsigned char* v4; // ebp
	int v5;            // eax
	int v6;            // ebx
	unsigned int v7;   // eax
	char v9;           // al
	unsigned char v10; // cl
	unsigned char v11; // cl
	char v12;          // al
	int v13[15];       // [esp+Ch] [ebp-3Ch]
	unsigned char v14; // [esp+4Ch] [ebp+4h]
	unsigned char v15; // [esp+4Ch] [ebp+4h]

	v1 = 0;
	v2 = 11 * v1;
	v3 = *getMemU32Ptr(0x5D4594, 1193372 + 4 * v1);
	v4 = getMemAt(0x5D4594, 1193388 + 4 * v2);
	if (!v3) {
		return 1;
	}
	v5 = v3 - 1;
	if (v5) {
		if (v5 == 1) {
			v6 = a1;
			if (*(uint32_t*)v4) {
				v7 = *(uint32_t*)(a1 + 96);
				if (v7 > *((uint32_t*)v4 + 4) && v7 != 9999) {
					return 0;
				}
			}
			if (*((uint32_t*)v4 + 1) && *(uint8_t*)(a1 + 100) & 0x10) {
				return 0;
			}
			if (*((uint32_t*)v4 + 2) && *(uint8_t*)(a1 + 100) & 0x20) {
				return 0;
			}
			v9 = *(uint8_t*)(a1 + 102);
			if (v9 < 0 && *((uint32_t*)v4 + 3) > (v9 & 0x7F)) {
				return 0;
			}
			if (*((uint32_t*)v4 + 5)) {
				strcpy((char*)v13, (const char*)(a1 + 111));
				sub_57A1E0(v13, 0, 0, 5, *(uint16_t*)(a1 + 163));
				v10 = 0;
				v14 = 0;
				while (v13[v14 + 6] == *(uint32_t*)(4 * v14 + v6 + 135)) {
					v14 = ++v10;
					if (v10 >= 5u) {
						v11 = 0;
						v15 = 0;
						while (*((uint8_t*)&v13[11] + v15) == *(uint8_t*)(v15 + v6 + 155)) {
							v15 = ++v11;
							if (v11 >= 4u) {
								if (v13[12] == *(uint32_t*)(v6 + 159)) {
									goto LABEL_26;
								}
								return 0;
							}
						}
						return 0;
					}
				}
				return 0;
			}
		LABEL_26:
			if (*((uint32_t*)v4 + 10) && *(uint32_t*)(v6 + 48) != NOX_CLIENT_VERS_CODE) {
				return 0;
			}
		}
		return 1;
	}
	v12 = *(uint8_t*)(a1 + 100);
	if (v12 & 0x10) {
		return 0;
	}
	if (v12 & 0x20) {
		return 0;
	}
	return *(uint32_t*)(a1 + 48) == NOX_CLIENT_VERS_CODE;
}

//----- (00489B80) --------------------------------------------------------
uint32_t* sub_489B80(int a1) {
	uint32_t* result;  // eax
	int v2;            // ebx
	unsigned char* v3; // edi
	uint32_t* v4;      // eax
	uint32_t* v5;      // eax
	uint32_t* v6;      // eax
	uint32_t* v7;      // eax
	uint32_t* v8;      // eax
	uint32_t* v9;      // eax
	uint32_t* v10;     // esi
	int v11;           // ebx
	int v12;           // ebx
	wchar2_t v13[16];   // [esp+0h] [ebp-20h]

	result = nox_new_window_from_file("filter.wnd", nox_xxx_windowMplayFilterProc_489E70);
	dword_5d4594_1193380 = result;
	if (result) {
		dword_5d4594_1193384 = nox_xxx_wndGetChildByID_46B0C0(result, 10012);
		v2 = 0;
		v3 = getMemAt(0x5D4594, 1193388 + 44 * v2);
		sub_46B120(*(uint32_t**)&dword_5d4594_1193380, a1);
		sub_46B120(*(uint32_t**)&dword_5d4594_1193384, *(int*)&dword_5d4594_1193380);
		nox_xxx_wndSetProc_46B2C0(*(int*)&dword_5d4594_1193384, nox_xxx_windowMplayFilterProc_489E70);
		if (*(uint32_t*)v3) {
			v4 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10028);
			v4[9] |= 4u;
		}
		nox_swprintf(v13, L"%d", *((uint32_t*)v3 + 4));
		v5 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10031);
		nox_window_call_field_94((int)v5, 16414, (int)v13, -1);
		if (*((uint32_t*)v3 + 1)) {
			v6 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10029);
			v6[9] |= 4u;
		}
		if (*((uint32_t*)v3 + 2)) {
			v7 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10030);
			v7[9] |= 4u;
		}
		if (*((uint32_t*)v3 + 3)) {
			v8 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10015);
			v8[9] |= 4u;
		}
		if (v3[12] & 2) {
			v9 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10017);
		} else {
			v9 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10016);
		}
		v10 = v9;
		v9[9] |= 4u;
		if (*((uint32_t*)v3 + 5)) {
			v10 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10014);
			v10[9] |= 4u;
		}
		if (*((uint32_t*)v3 + 10)) {
			v10 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10018);
			v10[9] |= 4u;
		}
		v11 = *getMemU32Ptr(0x5D4594, 1193372 + 4 * v2);
		if (v11) {
			v12 = v11 - 1;
			if (v12) {
				if (v12 == 1) {
					v10 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10026);
					sub_489DC0();
				}
			} else {
				v10 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10025);
				nox_window_set_hidden(*(int*)&dword_5d4594_1193384, 1);
			}
		} else {
			v10 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10024);
			nox_window_set_hidden(*(int*)&dword_5d4594_1193384, 1);
		}
		v10[9] |= 4u;
		result = *(uint32_t**)&dword_5d4594_1193380;
	}
	return result;
}

//----- (00489DC0) --------------------------------------------------------
void sub_489DC0() {
	uint32_t* v0; // eax
	uint32_t* v1; // eax

	nox_window_set_hidden(*(int*)&dword_5d4594_1193384, 0);
	if (nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193384, 10028)->draw_data.field_0 & 4) {
		v0 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193384, 10031);
		nox_xxx_wnd_46ABB0((int)v0, 1);
	} else {
		v1 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193384, 10031);
		nox_xxx_wnd_46ABB0((int)v1, 0);
	}
	if (nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193384, 10015)->draw_data.field_0 & 4) {
		sub_46AD20(*(uint32_t**)&dword_5d4594_1193384, 10016, 10017, 1);
	} else {
		sub_46AD20(*(uint32_t**)&dword_5d4594_1193384, 10016, 10017, 0);
	}
}

//----- (00489E70) --------------------------------------------------------
int nox_xxx_windowMplayFilterProc_489E70(int a1, int a2, int* a3, int a4) {
	int v3;       // ebx
	int v4;       // esi
	int result;   // eax
	uint32_t* v6; // eax
	int v7;       // [esp-Ch] [ebp-10h]
	int v8;       // [esp-Ch] [ebp-10h]

	v3 = 0;
	if (a2 == 23) {
		return 1;
	}
	if (a2 != 16391) {
		return 0;
	}
	v4 = nox_xxx_wndGetID_46B0A0(a3);
	nox_xxx_clientPlaySoundSpecial_452D80(766, 100);
	switch (v4) {
	case 10015:
		sub_46AD20(*(uint32_t**)&dword_5d4594_1193380, 10016, 10017, ((unsigned int)~a3[9] >> 2) & 1);
		return 0;
	case 10024:
		v7 = dword_5d4594_1193384;
		*getMemU32Ptr(0x5D4594, 1193372 + 4 * v3) = 0;
		nox_window_set_hidden(v7, 1);
		result = 0;
		break;
	case 10025:
		v8 = dword_5d4594_1193384;
		*getMemU32Ptr(0x5D4594, 1193372 + 4 * v3) = 1;
		nox_window_set_hidden(v8, 1);
		result = 0;
		break;
	case 10026:
		*getMemU32Ptr(0x5D4594, 1193372 + 4 * v3) = 2;
		sub_489DC0();
		result = 0;
		break;
	case 10028:
		v6 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1193380, 10031);
		nox_xxx_wnd_46ABB0((int)v6, ((unsigned int)~a3[9] >> 2) & 1);
		result = 0;
		break;
	default:
		return 0;
	}
	return result;
}

//----- (00489FB0) --------------------------------------------------------
int sub_489FB0() {
	int result; // eax

	result = dword_5d4594_1193380;
	if (dword_5d4594_1193380) {
		sub_489870();
		nox_xxx_wndClearCaptureMain_46ADE0(*(int*)&dword_5d4594_1193380);
		result = nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1193380);
		dword_5d4594_1193380 = 0;
	}
	return result;
}

//----- (00489FF0) --------------------------------------------------------
int sub_489FF0(int a1, int a2, const void* a3) {
	int result; // eax

	*getMemU32Ptr(0x5D4594, 1193372 + 4 * a1) = a2;
	result = 11 * a1;
	memcpy(getMemAt(0x5D4594, 1193388 + 44 * a1), a3, 0x2Cu);
	return result;
}

//----- (0048A210) --------------------------------------------------------
int nox_xxx_setSomeFunc_48A210(int a1) {
	int result; // eax

	result = a1;
	*getMemU32Ptr(0x5D4594, 1193504) = a1;
	return result;
}
