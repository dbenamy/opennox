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
#include "client__drawable__drawable.h"
#include "common__system__team.h"

#include "client__gui__chathelp.h"
#include "client__gui__gadgets__listbox.h"
#include "client__gui__guicon.h"
#include "client__gui__guiinv.h"
#include "client__gui__guiquit.h"
#include "client__gui__guivote.h"
#include "client__gui__window.h"
#include "client__network__cdecode.h"
#include "client__shell__noxworld.h"

#include "client__draw__fx.h"
#include "client__gui__guibook.h"
#include "client__video__draw_common.h"

#include "common/fs/nox_fs.h"
#include "common__magic__speltree.h"
#include "common__net_list.h"
#include "defs.h"
#include "input_common.h"
#include "operators.h"

extern uint32_t dword_8531A0_2572;
extern uint32_t dword_5d4594_1303508;
extern uint32_t dword_5d4594_1200776;
extern uint32_t dword_5d4594_1200796;
extern uint32_t dword_5d4594_1305788;
extern uint32_t nox_server_sanctuaryHelp_54276;
extern uint32_t dword_5d4594_1305748;
extern uint32_t dword_5d4594_1197352;
extern uint32_t dword_5d4594_1197356;
extern uint32_t dword_5d4594_1193712;
extern uint32_t nox_server_connectionType_3596;
extern void* nox_alloc_pixelSpan_1301844;
extern uint32_t nox_wol_servers_sorting_166704;
extern uint32_t dword_5d4594_1305680;
extern uint32_t dword_5d4594_1301848;
extern uint32_t dword_5d4594_1303452;
extern uint32_t dword_5d4594_1305684;
extern uint32_t nox_player_netCode_85319C;
extern int nox_win_width;
extern int nox_win_height;

extern uint32_t nox_color_white_2523948;
extern uint32_t nox_color_yellow_2589772;
extern uint32_t nox_color_black_2650656;

nox_render_data_t* nox_draw_curDrawData_3799572 = 0;

nox_list_item_t nox_gui_wol_servers_list = {0};

void nox_gui_winParentsReset_4A0CF0();
nox_window* nox_gui_winParentsTop_4A14F0();
nox_window* nox_gui_winParentsPop_4A18A0();
void nox_gui_winParentsPush_4A18C0(nox_window* win);

//----- (0048C580) --------------------------------------------------------
void sub_48C580(pixel8888* a1, int num) {
	unsigned int* pix = (unsigned int*)a1;
	for (int i = num - 1; i >= 0; i--) {
		unsigned int result = *pix;
		for (unsigned int* it = &pix[i]; it > pix; --it) {
			if (result > *it) {
				result = __sync_lock_test_and_set((volatile signed int*)it, result);
			}
		}
		*pix = result;
		++pix;
	}
}

//----- (0048CA70) --------------------------------------------------------
int nox_xxx_showObserverWindow_48CA70(int a1) { return nox_window_set_hidden(*(int*)&dword_5d4594_1193712, a1); }





//----- (0048D4B0) --------------------------------------------------------
int sub_48D4B0(int a1) {
	int result; // eax

	*getMemU32Ptr(0x5D4594, 1197304) = a1;
	if (a1 == 1) {
		result = sub_4C3460(0);
	} else {
		result = sub_4C3460(1);
	}
	return result;
}

//----- (0048D4F0) --------------------------------------------------------
int sub_48D4F0(unsigned short a1, unsigned short a2) {
	unsigned short v2; // cx

	v2 = 10000;
	if (a1 - 10000 < 0) {
		if (a2 >= 0xFFFF - (unsigned short)(10000 - a1)) {
			return 1;
		}
		v2 = a1;
	}
	return a2 < a1 && a2 >= a1 - v2;
}

//----- (0048D560) --------------------------------------------------------
int sub_48D560(unsigned short a1) {
	int* v1; // eax

	v1 = nox_common_list_getFirstSafe_425890(getMemIntPtr(0x5D4594, 1197340));
	if (!v1) {
		return 0;
	}
	while (v1[2] != a1) {
		v1 = nox_common_list_getNextSafe_4258A0(v1);
		if (!v1) {
			return 0;
		}
	}
	return 1;
}

//----- (0048D5A0) --------------------------------------------------------
uint32_t* sub_48D5A0(int a1) {
	uint32_t* result; // eax
	uint32_t* v2;     // ebx

	result = (uint32_t*)sub_48D4F0(*getMemU16Ptr(0x5D4594, 1197360), *(uint16_t*)(a1 + 1));
	if (!result) {
		result = (uint32_t*)sub_48D560(*(uint16_t*)(a1 + 1));
		if (!result) {
			result = calloc(*(unsigned char*)(a1 + 3) + 32, 1u);
			v2 = result;
			if (result) {
				sub_425770(result);
				v2[2] = *(unsigned short*)(a1 + 1);
				*((uint16_t*)v2 + 12) = *(unsigned char*)(a1 + 3);
				*((uint64_t*)v2 + 2) = nox_platform_get_ticks();
				memcpy(v2 + 8, (const void*)(a1 + 4), *(unsigned char*)(a1 + 3));
				if (*getMemU16Ptr(0x5D4594, 1197360) == *(uint16_t*)(a1 + 1)) {
					dword_5d4594_1197352 = v2;
				}
				result = (uint32_t*)sub_425790(getMemIntPtr(0x5D4594, 1197340), v2);
			}
		}
	}
	return result;
}

//----- (0048D660) --------------------------------------------------------
int sub_48D660() {
	unsigned long long v0; // rax
	int* v1;               // esi
	int* v2;               // edi

	LODWORD(v0) = dword_5d4594_1197352;
	if (!dword_5d4594_1197352) {
		if (dword_5d4594_1197356) {
			v0 = nox_platform_get_ticks() - *(uint64_t*)(dword_5d4594_1197356 + 16);
			if (v0 > 0x7530) {
				*getMemU16Ptr(0x5D4594, 1197360) = *(uint16_t*)(dword_5d4594_1197356 + 8);
				dword_5d4594_1197352 = dword_5d4594_1197356;
				LODWORD(v0) = nox_common_list_getNextSafe_4258A0(*(int**)&dword_5d4594_1197356);
				dword_5d4594_1197356 = v0;
			}
		}
	}
	v1 = *(int**)&dword_5d4594_1197352;
	if (dword_5d4594_1197352) {
		do {
			LODWORD(v0) = v1[2];
			if ((uint32_t)v0 != *getMemU16Ptr(0x5D4594, 1197360)) {
				break;
			}
			v2 = nox_common_list_getNextSafe_4258A0(v1);
			if (!v2) {
				v2 = nox_common_list_getFirstSafe_425890(getMemIntPtr(0x5D4594, 1197340));
				if (v2 == v1) {
					v2 = 0;
				}
			}
			nox_xxx_netOnPacketRecvCli_48EA70(31, (unsigned int)(v1 + 8), *((unsigned short*)v1 + 12));
			++*getMemU16Ptr(0x5D4594, 1197360);
			nox_common_list_remove_425920((uint32_t**)v1);
			free(v1);
			v1 = v2;
		} while (v2);
	}
	dword_5d4594_1197356 = v1;
	dword_5d4594_1197352 = 0;
	return v0;
}

//----- (0048D740) --------------------------------------------------------
int sub_48D740() {
	int result; // eax

	nox_common_list_clear_425760(getMemAt(0x5D4594, 1197340));
	result = 0;
	dword_5d4594_1197352 = 0;
	dword_5d4594_1197356 = 0;
	*getMemU16Ptr(0x5D4594, 1197360) = 0;
	return result;
}

//----- (0048D760) --------------------------------------------------------
void sub_48D760() {
	int* v0; // esi
	int* v1; // edi

	v0 = nox_common_list_getFirstSafe_425890(getMemIntPtr(0x5D4594, 1197340));
	if (v0) {
		do {
			v1 = nox_common_list_getNextSafe_4258A0(v0);
			nox_common_list_remove_425920((uint32_t**)v0);
			free(v0);
			v0 = v1;
		} while (v1);
	}
	nox_common_list_clear_425760(getMemAt(0x5D4594, 1197340));
	*getMemU16Ptr(0x5D4594, 1197360) = 0;
}

//----- (0048D7B0) --------------------------------------------------------
int* sub_48D7B0() {
	int* result; // eax

	for (result = nox_common_list_getFirstSafe_425890(getMemIntPtr(0x5D4594, 1197340)); result;
		 result = nox_common_list_getNextSafe_4258A0(result)) {
		;
	}
	return result;
}

//----- (004947E0) --------------------------------------------------------
char* sub_4947E0(int a1) {
	short v1;     // ax
	int v2;       // edi
	char* result; // eax
	int i;        // esi

	if (nox_common_gameFlags_check_40A5C0(1)) {
		v1 = nox_common_gameFlags_getVal_40A5B0();
		v2 = (unsigned short)nox_xxx_servGamedataGet_40A020(v1);
	} else {
		v2 = *((unsigned short*)nox_xxx_cliGamedataGet_416590(0) + 27);
	}
	result = nox_common_playerInfoGetFirst_416EA0();
	for (i = (int)result; result; i = (int)result) {
		if (!(*(uint8_t*)(i + 3680) & 1)) {
			if (i == a1) {
				if (nox_common_gameFlags_check_40A5C0(1024)) {
					if (*(uint32_t*)(i + 2140) >= v2) {
						*(uint32_t*)(i + 2140) = v2 - 1;
					}
				} else {
					*(uint32_t*)(i + 2136) = v2;
				}
			} else if (nox_common_gameFlags_check_40A5C0(1024)) {
				if (*(uint32_t*)(i + 2140) < v2) {
					*(uint32_t*)(i + 2140) = v2;
				}
			} else if (*(uint32_t*)(i + 2136) >= v2) {
				*(uint32_t*)(i + 2136) = v2 - 1;
			}
		}
		result = nox_common_playerInfoGetNext_416EE0(i);
	}
	return result;
}

//----- (004948B0) --------------------------------------------------------
int sub_4948B0(int a1) {
	short v1;   // ax
	int v2;     // edi
	char* i;    // esi
	int result; // eax
	int j;      // ebp
	char* v6;   // eax
	char* v7;   // esi
	int k;      // ebp
	char* v9;   // eax
	char* v10;  // esi

	if (nox_common_gameFlags_check_40A5C0(1)) {
		v1 = nox_common_gameFlags_getVal_40A5B0();
		v2 = (unsigned short)nox_xxx_servGamedataGet_40A020(v1);
	} else {
		v2 = *((unsigned short*)nox_xxx_cliGamedataGet_416590(0) + 27);
	}
	for (i = nox_server_teamFirst_418B10(); i; i = nox_server_teamNext_418B60((int)i)) {
		if (i == (char*)a1) {
			if (!nox_common_gameFlags_check_40A5C0(1024)) {
				*((uint32_t*)i + 13) = v2;
			}
		} else if (!nox_common_gameFlags_check_40A5C0(1024) && *((uint32_t*)i + 13) >= v2) {
			*((uint32_t*)i + 13) = v2 - 1;
		}
	}
	if (nox_common_gameFlags_check_40A5C0(1)) {
		result = nox_xxx_getFirstPlayerUnit_4DA7C0();
		for (j = result; result; j = result) {
			if (!nox_xxx_teamCompare2_419180(j + 48, *(uint8_t*)(a1 + 57))) {
				v6 = nox_common_playerInfoGetByID_417040(*(uint32_t*)(j + 36));
				v7 = v6;
				if (v6) {
					if (!(v6[3680] & 1)) {
						if (nox_common_gameFlags_check_40A5C0(1024)) {
							if (*((uint32_t*)v7 + 535) < v2) {
								*((uint32_t*)v7 + 535) = v2;
							}
						} else if (*((uint32_t*)v7 + 534) >= v2) {
							*((uint32_t*)v7 + 534) = v2 - 1;
						}
					}
				}
			}
			result = nox_xxx_getNextPlayerUnit_4DA7F0(j);
		}
	} else {
		result = nox_xxx_cliGetSpritePlayer_45A000();
		for (k = result; result; k = result) {
			if (!nox_xxx_teamCompare2_419180(k + 24, *(uint8_t*)(a1 + 57))) {
				v9 = nox_common_playerInfoGetByID_417040(*(uint32_t*)(k + 128));
				v10 = v9;
				if (v9) {
					if (!(v9[3680] & 1)) {
						if (nox_common_gameFlags_check_40A5C0(1024)) {
							if (*((uint32_t*)v10 + 535) < v2) {
								*((uint32_t*)v10 + 535) = v2;
							}
						} else if (*((uint32_t*)v10 + 534) >= v2) {
							*((uint32_t*)v10 + 534) = v2 - 1;
						}
					}
				}
			}
			result = sub_45A010(k);
		}
	}
	return result;
}

//----- (00494A60) --------------------------------------------------------
int nox_xxx_netCliProcUpdateStream_494A60(unsigned char* a1, int a2, uint32_t* a3) {
	unsigned short v3;  // di
	unsigned short v4;  // bp
	unsigned char* v5;  // esi
	int v6;             // ebx
	int v7;             // eax
	unsigned short v8;  // ax
	unsigned short v9;  // cx
	unsigned char* v10; // esi
	unsigned char v11;  // bl
	unsigned char* v12; // esi
	unsigned char v13;  // dl
	unsigned char v14;  // dl
	int v15;            // esi
	uint32_t* v16;      // eax
	unsigned char v17;  // cl
	int v18;            // edx
	unsigned short v20; // [esp+10h] [ebp-18h]
	unsigned short v21; // [esp+14h] [ebp-14h]
	unsigned char* v22; // [esp+18h] [ebp-10h]
	char v23[10];       // [esp+1Ch] [ebp-Ch]
	unsigned char v24;  // [esp+2Ch] [ebp+4h]
	unsigned char v25;  // [esp+2Ch] [ebp+4h]
	unsigned char v26;  // [esp+30h] [ebp+8h]

	v22 = a1;
	if (*a1 == 0xFFu) {
		v3 = *(uint16_t*)(a1 + 1);
		v4 = *(uint16_t*)(a1 + 3);
		v5 = a1 + 5;
		v6 = *(unsigned short*)(a1 + 3);
		v24 = nox_xxx_cliGenerateAlias_57B9A0((int)getMemAt(0x5D4594, 1198020), v3, v6, gameFrame());
		if (v24 != 0xFFu) {
			sub_57BA10((int)getMemAt(0x5D4594, 1198020 + 8 * v24), v3, v6, -1);
			v23[0] = 0xA5; // MSG_NEW_ALIAS
			v23[1] = v24;
			*(uint16_t*)&v23[2] = v3;
			*(uint16_t*)&v23[4] = v4;
			*(uint32_t*)&v23[6] = -1;
			nox_netlist_addToMsgListCli_40EBC0(a2, 0, v23, 10);
		}
	} else {
		v7 = 8 * *a1;
		v5 = a1 + 1;
		v3 = *getMemU16Ptr(0x5D4594, 1198020 + v7);
		v4 = *getMemU16Ptr(0x5D4594, 1198022 + v7);
	}
	v8 = *(uint16_t*)v5;
	v9 = *((uint16_t*)v5 + 1);
	v20 = *(uint16_t*)v5;
	v10 = v5 + 4;
	v21 = v9;
	v11 = *v10;
	v12 = v10 + 1;
	if ((v11 & 0x80u) == 0) {
		v25 = 0;
	} else {
		v13 = *v12++;
		v25 = v13;
	}
	v14 = *v12;
	v15 = (int)(v12 + 1);
	v26 = v14;
	if (v3 || v4) {
		v16 = nox_xxx_spriteCreate_48E970(v4, v3, v8, v9);
		if (v16) {
			v16[72] = gameFrame();
			v17 = (v11 >> 4) & 7;
			*((uint8_t*)v16 + 297) = v17;
			if (v17 > 3u) {
				*((uint8_t*)v16 + 297) = v17 + 1;
			}
			if (v16[69] != v26) {
				v18 = gameFrame();
				v16[69] = v26;
				v16[79] = v18;
			}
			nox_xxx_spriteSetFrameMB_45AB80((int)v16, v25);
		}
	}
	a3[0] = v20;
	a3[1] = v21;
	nox_xxx_cliUpdateCameraPos_435600(v20, v21);
	return v15 - (uint32_t)v22;
}

//----- (00494C30) --------------------------------------------------------
unsigned char* nox_xxx_netCliUpdateStream2_494C30(unsigned char* a1, int a2, int* a3) {
	unsigned char* v3; // esi
	unsigned char v4;  // al
	unsigned short v6; // di
	uint16_t* v7;      // esi
	unsigned short v8; // bp
	unsigned short v9; // bx
	short* v10;        // esi
	int v11;           // eax
	int v12;           // ecx
	int* v13;          // ebx
	int v14;           // esi
	int v15;           // eax
	int v16;           // ecx
	int v17;           // eax
	uint32_t* v18;     // eax
	int v19;           // edi
	char v20;          // cl
	unsigned char v21; // al
	unsigned char v22; // dl
	int v23;           // ecx
	unsigned char v24; // [esp+10h] [ebp-18h]
	int v25;           // [esp+14h] [ebp-14h]
	char v26[10];      // [esp+1Ch] [ebp-Ch]
	char v27;          // [esp+34h] [ebp+Ch]
	unsigned char v28; // [esp+34h] [ebp+Ch]

	v3 = a1;
	v4 = a1[0];
	v25 = 0;
	if (!v4) {
		v4 = a1[1];
		v3 = a1 + 1;
		if (!v4 && !a1[2]) {
			return (unsigned char*)-3;
		}
		v25 = 1;
	}
	if (v4 == 0xFFu) {
		v6 = *(uint16_t*)(v3 + 1);
		v7 = v3 + 3;
		v8 = *v7;
		v9 = *v7;
		v10 = v7 + 1;
		v24 = nox_xxx_cliGenerateAlias_57B9A0((int)getMemAt(0x5D4594, 1198020), v6, v9, gameFrame());
		if (v24 != 0xFFu) {
			sub_57BA10((int)getMemAt(0x5D4594, 1198020 + 8 * v24), v6, v9, gameFrame() + 60);
			v26[0] = 0xA5; // MSG_NEW_ALIAS
			v26[1] = v24;
			*(uint16_t*)&v26[2] = v6;
			*(uint16_t*)&v26[4] = v8;
			*(uint32_t*)&v26[6] = gameFrame() + 60;
			nox_netlist_addToMsgListCli_40EBC0(a2, 0, v26, 10);
		}
	} else {
		v11 = 8 * v4;
		v10 = (short*)(v3 + 1);
		v6 = *getMemU16Ptr(0x5D4594, 1198020 + v11);
		v8 = *getMemU16Ptr(0x5D4594, 1198022 + v11);
	}
	if (v25) {
		v12 = *v10;
		v13 = a3;
		v14 = (int)(v10 + 2);
		*a3 = v12;
		a3[1] = *(short*)(v14 - 2);
	} else {
		v13 = a3;
		v15 = *(char*)v10;
		v14 = (int)(v10 + 1);
		*a3 += v15;
		a3[1] += *(char*)(v14 - 1);
	}
	v16 = v13[0];
	if (v16 < 0) {
		return &a1[-v14];
	}
	if (v16 > 6000) {
		return &a1[-v14];
	}
	v17 = v13[1];
	if (v17 < 0) {
		return &a1[-v14];
	}
	if (v17 > 6000) {
		return &a1[-v14];
	}
	v18 = nox_xxx_spriteCreate_48E970(v8, v6, v16, v17);
	v19 = (int)v18;
	if (!v18) {
		return &a1[-v14];
	}
	if (v18[28] & 0x200000) {
		v18[72] = gameFrame();
		v27 = *(uint8_t*)v14;
		v20 = *(uint8_t*)v14;
		v21 = (*(uint8_t*)v14 >> 4) & 7;
		*(uint8_t*)(v19 + 297) = v21;
		if (v21 > 3u) {
			*(uint8_t*)(v19 + 297) = v21 + 1;
		}
		if (v20 < 0) {
			nox_xxx_spriteSetFrameMB_45AB80(v19, *(unsigned char*)++v14);
			v20 = v27;
		}
		if (*(uint8_t*)(v19 + 112) & 4) {
			v22 = *(uint8_t*)++v14;
			v28 = v22;
		} else {
			v28 = v20 & 0xF;
		}
		if (*(uint32_t*)(v19 + 276) != v28) {
			v23 = gameFrame();
			*(uint32_t*)(v19 + 276) = v28;
			*(uint32_t*)(v19 + 316) = v23;
		}
		++v14;
	} else {
		v18[72] = gameFrame();
		nox_xxx_sprite_49AA00_drawable(v18);
	}
	*v13 = *(uint32_t*)(v19 + 12);
	v13[1] = *(uint32_t*)(v19 + 16);
	return (unsigned char*)(v14 - (uint32_t)a1);
}

//----- (00494F00) --------------------------------------------------------
int sub_494F00() {
	int result; // eax
	int v1;     // esi
	int v2;     // eax

	*getMemU32Ptr(0x5D4594, 1200772) = nox_xxx_getTTByNameSpriteMB_44CFC0("Spark");
	if (!*getMemU32Ptr(0x5D4594, 1200772)) {
		return 0;
	}
	result = nox_xxx_getTTByNameSpriteMB_44CFC0("BlueSpark");
	dword_5d4594_1200776 = result;
	if (result) {
		result = nox_xxx_getTTByNameSpriteMB_44CFC0("YellowSpark");
		*getMemU32Ptr(0x5D4594, 1200780) = result;
		if (result) {
			result = nox_xxx_getTTByNameSpriteMB_44CFC0("CyanSpark");
			*getMemU32Ptr(0x5D4594, 1200784) = result;
			if (result) {
				result = nox_xxx_getTTByNameSpriteMB_44CFC0("GreenSpark");
				*getMemU32Ptr(0x5D4594, 1200788) = result;
				if (result) {
					result = nox_xxx_getTTByNameSpriteMB_44CFC0("Puff");
					*getMemU32Ptr(0x5D4594, 1200792) = result;
					if (result) {
						v1 = 0;
						while (1) {
							v2 = nox_xxx_getTTByNameSpriteMB_44CFC0(*(char**)getMemAt(0x587000, 161216 + v1));
							*getMemU32Ptr(0x5D4594, 1200812 + v1) = v2;
							if (!v2) {
								break;
							}
							v1 += 4;
							if (v1 >= 20) {
								dword_5d4594_1200796 = nox_xxx_getTTByNameSpriteMB_44CFC0("VioletSpark");
								return dword_5d4594_1200796 != 0;
							}
						}
						return 0;
					}
				}
			}
		}
	}
	return result;
}

//----- (00499360) --------------------------------------------------------
int nox_xxx_loadReflSheild_499360() {
	int v0;            // eax
	int v1;            // eax
	int v2;            // eax
	int v3;            // eax
	int v4;            // eax
	int v5;            // eax
	int v6;            // eax
	int v7;            // eax
	unsigned char* v8; // eax

	v0 = nox_xxx_getTTByNameSpriteMB_44CFC0("ReflectiveShieldNW");
	*getMemU32Ptr(0x5D4594, 1217468) = nox_new_drawable_for_thing(v0);
	v1 = nox_xxx_getTTByNameSpriteMB_44CFC0("ReflectiveShieldN");
	*getMemU32Ptr(0x5D4594, 1217472) = nox_new_drawable_for_thing(v1);
	v2 = nox_xxx_getTTByNameSpriteMB_44CFC0("ReflectiveShieldNE");
	*getMemU32Ptr(0x5D4594, 1217476) = nox_new_drawable_for_thing(v2);
	v3 = nox_xxx_getTTByNameSpriteMB_44CFC0("ReflectiveShieldW");
	*getMemU32Ptr(0x5D4594, 1217480) = nox_new_drawable_for_thing(v3);
	*getMemU32Ptr(0x5D4594, 1217484) = 0;
	v4 = nox_xxx_getTTByNameSpriteMB_44CFC0("ReflectiveShieldE");
	*getMemU32Ptr(0x5D4594, 1217488) = nox_new_drawable_for_thing(v4);
	v5 = nox_xxx_getTTByNameSpriteMB_44CFC0("ReflectiveShieldSW");
	*getMemU32Ptr(0x5D4594, 1217492) = nox_new_drawable_for_thing(v5);
	v6 = nox_xxx_getTTByNameSpriteMB_44CFC0("ReflectiveShieldS");
	*getMemU32Ptr(0x5D4594, 1217496) = nox_new_drawable_for_thing(v6);
	v7 = nox_xxx_getTTByNameSpriteMB_44CFC0("ReflectiveShieldSE");
	*getMemU32Ptr(0x5D4594, 1217500) = nox_new_drawable_for_thing(v7);
	v8 = getMemAt(0x5D4594, 1217468);
	while (1) {
		if (v8 != getMemAt(0x5D4594, 1217484)) {
			if (*(uint32_t*)v8) {
				*(uint32_t*)(*(uint32_t*)v8 + 120) |= 0x1000000u;
				goto LABEL_5;
			}
			return 0;
		}
	LABEL_5:
		v8 += 4;
		if ((int)v8 >= (int)getMemAt(0x5D4594, 1217504)) {
			*getMemU32Ptr(0x5D4594, 1217504) = 0;
			return 1;
		}
	}
}

//----- (00499450) --------------------------------------------------------
int sub_499450() {
	unsigned char* v0; // esi
	int result;        // eax

	v0 = getMemAt(0x5D4594, 1217468);
	do {
		result = *(uint32_t*)v0;
		if (*(uint32_t*)v0) {
			result = nox_xxx_spriteDelete_45A4B0(*(uint64_t**)v0);
		}
		*(uint32_t*)v0 = 0;
		v0 += 4;
	} while ((int)v0 < (int)getMemAt(0x5D4594, 1217504));
	*getMemU32Ptr(0x5D4594, 1217504) = 0;
	return result;
}

//----- (00499810) --------------------------------------------------------
int nox_xxx_drawShield_499810(nox_draw_viewport_t* vp, nox_drawable* dr) {
	int a1 = vp;
	int a2 = dr;
	int v3; // [esp-4h] [ebp-8h]

	*(uint32_t*)(*getMemU32Ptr(0x5D4594, 1217468 + 4 * *(unsigned char*)(a2 + 297)) + 12) =
		*(uint32_t*)(a2 + 12) + *getMemU32Ptr(0x587000, 161776 + 8 * *(unsigned char*)(a2 + 297));
	*(uint32_t*)(*getMemU32Ptr(0x5D4594, 1217468 + 4 * *(unsigned char*)(a2 + 297)) + 16) =
		*(uint32_t*)(a2 + 16) + *(short*)(a2 + 104) + *getMemU32Ptr(0x587000, 161780 + 8 * *(unsigned char*)(a2 + 297));
	v3 = *getMemU32Ptr(0x5D4594, 1217468 + 4 * *(unsigned char*)(a2 + 297));
	(*(void (**)(int, int))(v3 + 300))(a1, v3);
	return 0;
}

//----- (00499880) --------------------------------------------------------
uint32_t* nox_xxx_fxDrawTurnUndead_499880(short* a1) {
	int i;            // ebx
	uint32_t* result; // eax
	uint32_t* v3;     // esi
	int v4;           // eax
	double v5;        // st7

	if (!*getMemU32Ptr(0x5D4594, 1217508)) {
		*getMemU32Ptr(0x5D4594, 1217508) = nox_xxx_getTTByNameSpriteMB_44CFC0("UndeadKiller");
	}
	for (i = 0; i < 256; i += 6) {
		result = (uint32_t*)nox_xxx_spriteLoadAdd_45A360_drawable(*getMemIntPtr(0x5D4594, 1217508), *a1, a1[1]);
		v3 = result;
		if (result) {
			v4 = 8 * (short)i;
			*((uint16_t*)v3 + 254) = i;
			*((float*)v3 + 117) = *getMemFloatPtr(0x587000, 194136 + v4) * 4.0;
			v5 = *getMemFloatPtr(0x587000, 194140 + v4) * 4.0;
			v3[119] = 0;
			*((float*)v3 + 118) = v5;
			v3[79] = gameFrame();
			v3[81] = *a1;
			v3[82] = a1[1];
			v3[115] = nox_xxx_sprite_4CA540;
			nox_xxx_spriteToList_49BC80_drawable(v3);
			nox_xxx_spriteToSightDestroyList_49BAB0_drawable(v3);
		}
	}
	return result;
}

//----- (00499CF0) --------------------------------------------------------
void nox_xxx_bookRewardCli_499CF0(int* a1, int a2, int a3) {
	unsigned int result; // eax
	int v4;              // esi
	int2 a3a;            // [esp+8h] [ebp-8h]

	if (!nox_common_gameFlags_check_40A5C0(2048) ||
		(result = nox_xxx_bookGet_430B40_get_mouse_prev_seq() - *getMemU32Ptr(0x5D4594, 1217504), result >= 2)) {
		*getMemU32Ptr(0x5D4594, 1217504) = nox_xxx_bookGet_430B40_get_mouse_prev_seq();
		if (a1 == (int*)2) {
			v4 = 0;
		} else {
			v4 = (a1 == (int*)3) + 2;
		}
		a3a.field_0 = 5;
		a3a.field_4 = nox_win_height / 3;
		nox_xxx_bookSetForward_45D200(a1, a2, &a3a);
		nox_xxx_draw_499E70(v4, a3a.field_0, a3a.field_4, 271, 166, 1, 1);
		nox_xxx_draw_499E70(v4, a3a.field_0, a3a.field_4, 135, 166, 2, 1);
		nox_xxx_draw_499E70(v4, a3a.field_0, a3a.field_4 + 166, 135, 166, 2, 1);
		nox_xxx_draw_499E70(v4, a3a.field_0 + 271, a3a.field_4, 271, 166, 1, 2);
		nox_xxx_draw_499E70(v4, a3a.field_0 + 135, a3a.field_4, 135, 166, 2, 2);
		nox_xxx_draw_499E70(v4, a3a.field_0 + 135, a3a.field_4 + 166, 135, 166, 2, 2);
		if (a1 != (int*)4 && a3 == 1) {
			nox_xxx_bookFillAll_45D570((int)a1, a2);
		}
	}
}

//----- (00499F60) --------------------------------------------------------
void sub_499F60(int a1, int a2, int a3, short a4, char a5, char a6, char a7, char a8, char a9, int a10) {
	uint32_t* result; // eax
	int v11;          // edx
	int v12;          // ecx
	uint32_t* v13;    // esi
	int v14;          // eax

	if (!*getMemU32Ptr(0x5D4594, 1217512)) {
		*getMemU32Ptr(0x5D4594, 1217512) = nox_xxx_getTTByNameSpriteMB_44CFC0("RedBubbleParticle");
		*getMemU32Ptr(0x5D4594, 1217516) = nox_xxx_getTTByNameSpriteMB_44CFC0("WhiteBubbleParticle");
		*getMemU32Ptr(0x5D4594, 1217520) = nox_xxx_getTTByNameSpriteMB_44CFC0("LightBlueBubbleParticle");
		*getMemU32Ptr(0x5D4594, 1217524) = nox_xxx_getTTByNameSpriteMB_44CFC0("OrangeBubbleParticle");
		*getMemU32Ptr(0x5D4594, 1217528) = nox_xxx_getTTByNameSpriteMB_44CFC0("GreenBubbleParticle");
		*getMemU32Ptr(0x5D4594, 1217532) = nox_xxx_getTTByNameSpriteMB_44CFC0("VioletBubbleParticle");
		*getMemU32Ptr(0x5D4594, 1217536) = nox_xxx_getTTByNameSpriteMB_44CFC0("LightVioletBubbleParticle");
		*getMemU32Ptr(0x5D4594, 1217540) = nox_xxx_getTTByNameSpriteMB_44CFC0("YellowBubbleParticle");
	}
	result = (uint32_t*)nox_xxx_spriteLoadAdd_45A360_drawable(a1, a2, a3);
	v13 = result;
	if (result) {
		BYTE1(v11) = HIBYTE(a4);
		LOBYTE(result) = *((uint8_t*)result + 160);
		LOBYTE(v12) = *((uint8_t*)v13 + 156);
		*((uint16_t*)v13 + 52) = a4;
		LOBYTE(v11) = *((uint8_t*)v13 + 152);
		v13[108] = nox_color_rgb_4344A0(v11, v12, (int)result);
		if (a1 == *getMemU32Ptr(0x5D4594, 1217512)) {
			v14 = nox_color_rgb_4344A0(255, 128, 128);
		} else if (a1 == *getMemU32Ptr(0x5D4594, 1217516)) {
			v14 = nox_color_rgb_4344A0(255, 255, 255);
		} else if (a1 == *getMemU32Ptr(0x5D4594, 1217524)) {
			v14 = nox_color_rgb_4344A0(255, 100, 50);
		} else if (a1 == *getMemU32Ptr(0x5D4594, 1217528)) {
			v14 = nox_color_rgb_4344A0(64, 255, 64);
		} else if (a1 == *getMemU32Ptr(0x5D4594, 1217532)) {
			v14 = nox_color_rgb_4344A0(255, 100, 255);
		} else if (a1 == *getMemU32Ptr(0x5D4594, 1217536)) {
			v14 = nox_color_rgb_4344A0(255, 200, 255);
		} else if (a1 == *getMemU32Ptr(0x5D4594, 1217540)) {
			v14 = nox_color_rgb_4344A0(255, 255, 200);
		} else {
			v14 = nox_color_rgb_4344A0(200, 200, 255);
		}
		v13[109] = v14;
		*((uint8_t*)v13 + 440) = a5;
		*((uint8_t*)v13 + 443) = a6;
		*((uint8_t*)v13 + 442) = a6;
		*((uint8_t*)v13 + 441) = 1;
		*((uint8_t*)v13 + 444) = a8;
		*((uint8_t*)v13 + 445) = a9;
		*((uint8_t*)v13 + 446) = a7;
		nox_xxx_spriteToSightDestroyList_49BAB0_drawable(v13);
		nox_xxx_spriteTransparentDecay_49B950(v13, a10);
		nox_xxx_sprite_45A110_drawable(v13);
	}
}
// 49A025: variable 'v11' is possibly undefined
// 49A025: variable 'v12' is possibly undefined

//----- (0049A3D0) --------------------------------------------------------
char* nox_xxx_clientEquip_49A3D0(char a1, int a2, int a3, int a4) {
	char* npc;   // eax
	uint32_t* k; // edx
	char* v7;    // ecx
	char** v8;   // edi
	int l;       // esi
	uint32_t* i; // edx
	char* v12;   // ecx
	char** v13;  // edi
	int j;       // esi

	npc = nox_npc_by_id(a2);
	if (!npc) {
		return 0;
	}
	if (a1 == 81 || a1 == 80) {
		int v10 = 0;
		for (i = npc + 32; *i; i += 6) {
			if (++v10 >= 27) {
				return npc;
			}
		}
		v12 = &npc[24 * v10];
		*((uint32_t*)v12 + 8) = a3;
		v13 = (char**)(v12 + 36);
		*((uint32_t*)npc + 326) |= a3;
		for (j = 0; j < 4; ++j) {
			npc = (char*)nox_xxx_modifGetDescById_413330(*(unsigned char*)(j + a4));
			*v13 = npc;
			++v13;
		}
	} else {
		int v5 = 0;
		for (k = npc + 680; *k; k += 6) {
			if (++v5 >= 26) {
				return npc;
			}
		}
		v7 = &npc[24 * v5];
		*((uint32_t*)v7 + 170) = a3;
		v8 = (char**)(v7 + 684);
		*((uint32_t*)npc + 327) |= a3;
		for (l = 0; l < 4; ++l) {
			npc = (char*)nox_xxx_modifGetDescById_413330(*(unsigned char*)(l + a4));
			*v8 = npc;
			++v8;
		}
	}
	return npc;
}

//----- (0049AEA0) --------------------------------------------------------
int sub_49AEA0() {
	if (nox_alloc_pixelSpan_1301844) {
		nox_free_alloc_class(*(void**)&nox_alloc_pixelSpan_1301844);
		nox_alloc_pixelSpan_1301844 = 0;
	}
	if (dword_5d4594_1301848) {
		free(*(void**)&dword_5d4594_1301848);
		dword_5d4594_1301848 = 0;
	}
	return 1;
}

//----- (0049B3E0) --------------------------------------------------------
int sub_49B3E0() {
	int result; // eax

	result = nox_new_window_from_file("GGOver.wnd", sub_49B420);
	dword_5d4594_1303452 = result;
	if (result) {
		nox_window_set_hidden(result, 1);
		nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1303452, 0);
		result = 1;
	}
	return result;
}

//----- (0049B420) --------------------------------------------------------
int sub_49B420(int a1, int a2, int* a3, int a4) {
	int v3; // esi

	if (a2 == 16391) {
		v3 = nox_xxx_wndGetID_46B0A0(a3);
		nox_xxx_clientPlaySoundSpecial_452D80(766, 100);
		if (v3 == 10701) {
			nox_client_quit_4460C0();
			sub_49B6B0();
		} else if (v3 == 10702) {
			LOWORD(a2) = 1008;
			nox_xxx_netClientSend2_4E53C0(31, &a2, 2, 0, 1);
			sub_49B6B0();
			return 0;
		}
	}
	return 0;
}

//----- (0049B490) --------------------------------------------------------
int sub_49B490() {
	int result; // eax

	result = nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1303452);
	dword_5d4594_1303452 = 0;
	return result;
}

//----- (0049B6B0) --------------------------------------------------------
int sub_49B6B0() {
	nox_window_set_hidden(*(int*)&dword_5d4594_1303452, 1);
	nox_xxx_wnd_46ABB0(*(int*)&dword_5d4594_1303452, 0);
	return nox_xxx_windowFocus_46B500(0);
}

//----- (0049B7A0) --------------------------------------------------------
void nox_xxx_consoleEsc_49B7A0() {
	int v0; // esi

	v0 = 0;
	if (!nox_xxx_guiCursor_477600() && !nox_video_inFadeTransition_44E0D0()) {
		if (sub_460660()) {
			v0 = 1;
		}
		if (!sub_46A6A0() && v0 != 1) {
			if (sub_45D9B0() == 1) {
				sub_45D870();
			} else {
				if (sub_4BFE40()) {
					v0 = 1;
				}
				if (nox_xxx_quickBarClose_4606B0()) {
					v0 = 1;
				}
				if (sub_462740()) {
					v0 = 1;
				}
				if (!sub_44A4E0() && v0 != 1) {
					if (sub_479590() == 2) {
						sub_4795A0(1);
					} else if (sub_479590() == 3) {
						sub_4795A0(1);
					} else if (sub_479590() == 4) {
						sub_4795A0(1);
					} else if (!sub_478040() && !sub_479950()) {
						if (sub_467C10()) {
							v0 = 1;
						}
						if (nox_xxx_bookHideMB_45ACA0(0)) {
							v0 = 1;
						}
						if (nox_gui_console_Hide_4512B0()) {
							v0 = 1;
						}
						if (sub_446780()) {
							v0 = 1;
						}
						if (nox_xxx_guiServerOptionsTryHide_4574D0()) {
							v0 = 1;
						}
						if (sub_48CAD0()) {
							v0 = 1;
						}
						if (sub_4AD9B0(1)) {
							v0 = 1;
						}
						if (sub_4C35B0(1)) {
							v0 = 1;
						}
						if (!sub_46D6F0() && v0 != 1) {
							if (nox_xxx_game_4DCCB0()) {
								sub_445C40();
							} else {
								nox_xxx_clientPlaySoundSpecial_452D80(231, 100);
							}
						}
					}
				}
			}
		}
	}
}
// 49B874: variable 'v1' is possibly undefined
// 49B881: variable 'v2' is possibly undefined

//----- (0049BB80) --------------------------------------------------------
void* sub_49BB80(char a1) {
	void* result; // eax

	*getMemU8Ptr(0x5D4594, 1303504) = a1;
	*getMemU8Ptr(0x5D4594, 1303512) = 0;
	*getMemU32Ptr(0x5D4594, 1303516) = gameFrame();
	result = nox_xxx_spellGetDefArrayPtr_424820();
	dword_5d4594_1303508 = result;
	return result;
}

//----- (0049BBB0) --------------------------------------------------------
void sub_49BBB0() { *getMemU8Ptr(0x5D4594, 1303504) = 0; }

//----- (0049BBC0) --------------------------------------------------------
void sub_49BBC0() {
	int v0;           // eax
	unsigned char v1; // [esp+0h] [ebp-4h]

	if (getMemByte(0x5D4594, 1303504)) {
		v1 = nox_xxx_spellPhonemes_424A20(getMemByte(0x5D4594, 1303504), getMemByte(0x5D4594, 1303512));
		if (gameFrame() >= *getMemIntPtr(0x5D4594, 1303516)) {
			v0 = nox_xxx_spellGetPhoneme_4FE1C0(nox_player_netCode_85319C, v1);
			nox_xxx_clientPlaySoundSpecial_452D80(v0, 100);
			nox_client_setPhonemeFrame_476E00(*getMemU32Ptr(0x587000, 163576 + 4 * v1));
			*getMemU32Ptr(0x5D4594, 1303516) = gameFrame() + 3;
			dword_5d4594_1303508 = nox_xxx_updateSpellRelated_424830(*(int*)&dword_5d4594_1303508, v1);
			++*getMemU8Ptr(0x5D4594, 1303512);
		}
		if (**(uint32_t**)&dword_5d4594_1303508 == getMemByte(0x5D4594, 1303504)) {
			sub_49BBB0();
		}
	}
}

//----- (0049C160) --------------------------------------------------------
uint32_t* nox_xxx_clientAddRayEffect_49C160(int a1) {
	uint32_t* result;   // eax
	int v2;             // eax
	int v3;             // esi
	int v4;             // eax
	int v5;             // ebx
	uint32_t* v6;       // eax
	uint32_t* v7;       // edi
	int v8;             // esi
	int v9;             // edi
	int v10;            // kr00_4
	int v11;            // ecx
	int v12;            // edx
	int v13;            // edx
	unsigned char* v14; // ecx

	result = *(uint32_t**)getMemAt(0x5D4594, 1304312);
	if (*getMemIntPtr(0x5D4594, 1304312) < 96) {
		if (!*getMemU32Ptr(0x5D4594, 1304352)) {
			*getMemU32Ptr(0x5D4594, 1304352) = nox_xxx_getTTByNameSpriteMB_44CFC0("DynamicLightning");
			*getMemU32Ptr(0x5D4594, 1304356) = nox_xxx_getTTByNameSpriteMB_44CFC0("DynamicChainLightning");
			*getMemU32Ptr(0x5D4594, 1304360) = nox_xxx_getTTByNameSpriteMB_44CFC0("DynamicEnergyBolt");
			*getMemU32Ptr(0x5D4594, 1304364) = nox_xxx_getTTByNameSpriteMB_44CFC0("OrbRay");
			*getMemU32Ptr(0x5D4594, 1304368) = nox_xxx_getTTByNameSpriteMB_44CFC0("PlasmaRay");
			*getMemU32Ptr(0x5D4594, 1304372) = nox_xxx_getTTByNameSpriteMB_44CFC0("DrainManaRay");
			*getMemU32Ptr(0x5D4594, 1304376) = nox_xxx_getTTByNameSpriteMB_44CFC0("HealRay");
			*getMemU32Ptr(0x5D4594, 1304380) = nox_xxx_getTTByNameSpriteMB_44CFC0("CharmRay");
			*getMemU32Ptr(0x5D4594, 1304384) = nox_xxx_getTTByNameSpriteMB_44CFC0("DrainManaOrb");
			*getMemU32Ptr(0x5D4594, 1304388) = nox_xxx_getTTByNameSpriteMB_44CFC0("HealOrb");
			*getMemU32Ptr(0x5D4594, 1304392) = nox_xxx_getTTByNameSpriteMB_44CFC0("CharmOrb");
			*getMemU32Ptr(0x5D4594, 1304396) = nox_xxx_getTTByNameSpriteMB_44CFC0("HarpoonRope");
		}
		v2 = nox_xxx_netClearHighBit_578B30(*(uint16_t*)(a1 + 3));
		v3 = v2;
		v4 = nox_xxx_netClearHighBit_578B30(*(uint16_t*)(a1 + 5));
		v5 = v4;
		v6 = nox_xxx_netTestHighBit_578B70(*(unsigned short*)(a1 + 3)) ? nox_xxx_netSpriteByCodeStatic_45A720(v3)
																	   : nox_xxx_netSpriteByCodeDynamic_45A6F0(v3);
		v7 = v6;
		result = nox_xxx_netTestHighBit_578B70(*(unsigned short*)(a1 + 5)) ? nox_xxx_netSpriteByCodeStatic_45A720(v5)
																		   : nox_xxx_netSpriteByCodeDynamic_45A6F0(v5);
		if (v7 && result) {
			v8 = v7[3];
			v9 = v7[4];
			v10 = result[4] - v9;
			v11 = v8 + ((int)result[3] - v8) / 2;
			result = (uint32_t*)(v9 + v10 / 2);
			switch (*(unsigned char*)(a1 + 1)) {
			case 1u:
				v12 = *getMemU32Ptr(0x5D4594, 1304368);
				break;
			case 2u:
				v12 = *getMemU32Ptr(0x5D4594, 1304380);
				break;
			case 3u:
				v12 = *getMemU32Ptr(0x5D4594, 1304356);
				break;
			case 4u:
				v12 = *getMemU32Ptr(0x5D4594, 1304360);
				break;
			case 5u:
				v12 = *getMemU32Ptr(0x5D4594, 1304372);
				break;
			case 6u:
				v12 = *getMemU32Ptr(0x5D4594, 1304376);
				break;
			case 7u:
				v12 = *getMemU32Ptr(0x5D4594, 1304396);
				break;
			case 0x8Cu:
				v12 = *getMemU32Ptr(0x5D4594, 1304352);
				break;
			default:
				return result;
			}
			result = (uint32_t*)nox_xxx_spriteLoadAdd_45A360_drawable(v12, v11, v9 + v10 / 2);
			if (!result) {
				return result;
			}
			*((uint8_t*)result + 432) = 1;
			*(uint32_t*)((char*)result + 437) = *(unsigned short*)(a1 + 3);
			*(uint32_t*)((char*)result + 441) = *(unsigned short*)(a1 + 5);
			v13 = 0;
			*(uint32_t*)((char*)result + 433) = *(unsigned char*)(a1 + 2);
			v14 = getMemAt(0x5D4594, 1303924);
			while (*(uint32_t*)v14) {
				v14 += 4;
				++v13;
				if ((int)v14 >= (int)getMemAt(0x5D4594, 1304308)) {
					return result;
				}
			}
			*getMemU32Ptr(0x5D4594, 1303924 + 4 * v13) = result;
		}
	}
	return result;
}
// 49C248: variable 'v2' is possibly undefined
// 49C256: variable 'v4' is possibly undefined

//----- (0049C450) --------------------------------------------------------
void nox_xxx_clientRemoveRayEffect_49C450(int a1) {
	int v1;  // esi
	int* v2; // ecx

	v1 = 0;
	v2 = getMemIntPtr(0x5D4594, 1303924);
	while (1) {
		int result = *v2;
		if (*v2) {
			if (*(unsigned short*)(a1 + 3) == *(uint32_t*)(result + 437) &&
				*(unsigned short*)(a1 + 5) == *(uint32_t*)(result + 441)) {
				break;
			}
		}
		++v2;
		++v1;
		if ((int)v2 >= (int)getMemAt(0x5D4594, 1304308)) {
			return;
		}
	}
	nox_xxx_spriteDeleteStatic_45A4E0_drawable(*v2);
	*getMemU32Ptr(0x5D4594, 1303924 + 4 * v1) = 0;
}

//----- (0049C4B0) --------------------------------------------------------
void nox_xxx_spriteDeleteSomeList_49C4B0() {
	int v0;  // esi
	int* v1; // edi

	v0 = 0;
	if (*getMemU32Ptr(0x5D4594, 1304308) > 0) {
		v1 = getMemIntPtr(0x5D4594, 1303540);
		do {
			nox_xxx_spriteDeleteStatic_45A4E0_drawable(*v1);
			++v0;
			++v1;
		} while (v0 < *getMemIntPtr(0x5D4594, 1304308));
	}
	*getMemU32Ptr(0x5D4594, 1304308) = 0;
	sub_4C5050();
}

//----- (0049C4F0) --------------------------------------------------------
void nox_xxx_sprite_49C4F0() {
	int* v0 = getMemIntPtr(0x5D4594, 1303924);
	do {
		if (*v0) {
			nox_xxx_spriteDeleteStatic_45A4E0_drawable(*v0);
			*v0 = 0;
		}
		++v0;
	} while ((int)v0 < (int)getMemAt(0x5D4594, 1304308));
}

//----- (0049C520) --------------------------------------------------------
int sub_49C520(nox_drawable* a1p) {
	int a1 = a1p;
	unsigned char* v1; // eax
	int v2;            // eax
	unsigned char* i;  // ecx

	v1 = getMemAt(0x5D4594, 1303924);
	while (a1 != *(uint32_t*)v1) {
		v1 += 4;
		if ((int)v1 >= (int)getMemAt(0x5D4594, 1304308)) {
			v2 = 0;
			if (*getMemIntPtr(0x5D4594, 1304308) <= 0) {
				return 0;
			}
			for (i = getMemAt(0x5D4594, 1303540); a1 != *(uint32_t*)i; i += 4) {
				if (++v2 >= *getMemIntPtr(0x5D4594, 1304308)) {
					return 0;
				}
			}
			return 1;
		}
	}
	return 1;
}

//----- (0049C760) --------------------------------------------------------
int nox_xxx_wnd_49C760(int a1, int a2, int* a3, int a4) {
	int v3; // esi

	if (a2 == 16391) {
		v3 = nox_xxx_wndGetID_46B0A0(a3);
		nox_xxx_clientPlaySoundSpecial_452D80(766, 100);
		if (v3 == 4103) {
			sub_49C7A0();
		}
	}
	return 1;
}

//----- (0049C7A0) --------------------------------------------------------
int sub_49C7A0() {
	int result; // eax

	result = dword_5d4594_1305680;
	if (dword_5d4594_1305680) {
		nox_server_sanctuaryHelp_54276 =
			((unsigned int)~(
				 nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1305680, 4104)->draw_data.field_0) >>
			 2) &
			1;
		nox_xxx_wnd_46C6E0(*(int*)&dword_5d4594_1305680);
		nox_xxx_wndClearCaptureMain_46ADE0(*(int*)&dword_5d4594_1305680);
		nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1305680);
		dword_5d4594_1305680 = 0;
		nox_xxx_windowFocus_46B500(0);
		result = nox_common_gameFlags_check_40A5C0(1);
		if (result) {
			result = sub_459D80(0);
		}
	}
	return result;
}

//----- (0049C810) --------------------------------------------------------
int sub_49C810() { return dword_5d4594_1305680 != 0; }

//----- (0049CA60) --------------------------------------------------------
int sub_49CA60(int a1, int a2, int* a3, int a4) {
	int v3;       // esi
	uint32_t* v4; // eax
	int v5;       // eax
	int v6;       // eax

	if (a2 == 16391) {
		v3 = nox_xxx_wndGetID_46B0A0(a3);
		nox_xxx_clientPlaySoundSpecial_452D80(766, 100);
		if (v3 == 10353) {
			nox_xxx_wnd_46C6E0(*(int*)&dword_5d4594_1305684);
			nox_xxx_wndClearCaptureMain_46ADE0(*(int*)&dword_5d4594_1305684);
			if (nox_common_gameFlags_check_40A5C0(128) && nox_server_sanctuaryHelp_54276) {
				nox_xxx_cliShowHelpGui_49C560();
			} else {
				sub_459D80(0);
			}
			v4 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_1305684, 10352);
			v5 = nox_window_call_field_94((int)v4, 16404, 0, 0);
			nox_server_connectionType_3596 = v5 + 1;
			v6 = sub_40A710(v5 + 1);
			nox_xxx_rateUpdate_40A6D0(v6);
			nox_xxx_windowDestroyMB_46C4E0(*(uint32_t**)&dword_5d4594_1305684);
			dword_5d4594_1305684 = 0;
			nox_xxx_windowFocus_46B500(0);
		}
	}
	return 1;
}

//----- (0049CB40) --------------------------------------------------------
int sub_49CB40() { return dword_5d4594_1305684 != 0; }

//----- (0049FDB0) --------------------------------------------------------
void sub_49FDB0(int a1) {
	unsigned char* v1; // ebx
	int j;             // esi
	int v3;            // eax
	unsigned char* v4; // ebx
	int i;             // esi
	int v6;            // eax
	int v7;            // [esp+0h] [ebp-90h]
	char v8[140];      // [esp+4h] [ebp-8Ch]

	if (!dword_5d4594_1305788) {
		if (0) {
			v7 = 0;
			if (*getMemU32Ptr(0x587000, 166016 + 4 * a1) > 0) {
				v4 = getMemAt(0x587000, 166032 + 80 * a1);
				do {
					for (i = 0; i < (char)*v4; ++i) {
						v6 = 8 * (12 * a1 + (char)v4[i + 1]);
						sub_420DA0(*getMemFloatPtr(0x587000, 165360 + v6), *getMemFloatPtr(0x587000, 165364 + v6));
					}
					strcpy(&v8[4], *((const char**)v4 + 3));
					sub_4211D0((int)v8);
					sub_4214D0();
					v4 += 16;
					++v7;
				} while (v7 < *getMemIntPtr(0x587000, 166016 + 4 * a1));
			}
		} else {
			v1 = getMemAt(0x587000, 165744);
			do {
				for (j = 0; j < (char)*v1; ++j) {
					v3 = 8 * (char)v1[j + 1];
					sub_420DA0(*getMemFloatPtr(0x587000, 165104 + v3), *getMemFloatPtr(0x587000, 165108 + v3));
				}
				strcpy(&v8[4], *((const char**)v1 + 3));
				sub_4211D0((int)v8);
				sub_4214D0();
				v1 += 16;
			} while ((int)v1 < (int)getMemAt(0x587000, 166016));
		}
		dword_5d4594_1305788 = 1;
	}
}

//----- (0049FF20) --------------------------------------------------------
uint32_t* sub_49FF20() {
	uint32_t* result; // eax

	result = *(uint32_t**)&dword_5d4594_1305788;
	if (dword_5d4594_1305788) {
		result = sub_421B10();
		dword_5d4594_1305788 = 0;
	}
	return result;
}

//----- (0049FFA0) --------------------------------------------------------
int* sub_49FFA0(int a1) {
	int* result; // eax
	int* v2;     // esi
	int* v3;     // edi

	if (!*getMemU32Ptr(0x5D4594, 1305808)) {
		nox_common_list_clear_425760(&nox_gui_wol_servers_list);
	}
	result = nox_common_list_getFirstSafe_425890(&nox_gui_wol_servers_list);
	v2 = result;
	if (result) {
		do {
			v3 = nox_common_list_getNextSafe_4258A0(v2);
			nox_common_list_remove_425920((uint32_t**)v2);
			if (a1) {
				nox_xxx_windowDestroyMB_46C4E0((uint32_t*)v2[7]);
			}
			free(v2);
			v2 = v3;
		} while (v3);
		*getMemU32Ptr(0x5D4594, 1305808) = 1;
	} else {
		*getMemU32Ptr(0x5D4594, 1305808) = 1;
	}
	return result;
}

//----- (004A0020) --------------------------------------------------------
char* sub_4A0020() { return &nox_gui_wol_servers_list; }

//----- (004A0030) --------------------------------------------------------
int nox_wol_servers_addResult_4A0030(nox_gui_server_ent_t* srv) {
	int* v3;     // edi
	wchar2_t* v6; // ebp
	wchar2_t* v7; // eax
	wchar2_t* v8; // ebp
	wchar2_t* v9; // eax

	nox_gui_server_ent_t* rec = calloc(1, sizeof(nox_gui_server_ent_t));
	memcpy(rec, srv, sizeof(nox_gui_server_ent_t));

	int v2 = 0;
	switch (nox_wol_servers_sorting_166704) {
	case 0: // by name (asc)
		if (nox_gui_wol_servers_list.field_1 == &nox_gui_wol_servers_list) {
			return sub_425790(&nox_gui_wol_servers_list, rec);
		}
		v3 = nox_common_list_getFirstSafe_425890(&nox_gui_wol_servers_list);
		if (!v3) {
			nox_common_list_append_4258E0(&nox_gui_wol_servers_list, rec);
			return 0;
		}
		do {
			if (nox_strcmpi(rec->server_name, (const char*)v3 + 120) <= 0) {
				nox_common_list_append_4258E0((int)v3, rec);
				return v2;
			}
			++v2;
			v3 = nox_common_list_getNextSafe_4258A0(v3);
		} while (v3);
		nox_common_list_append_4258E0(&nox_gui_wol_servers_list, rec);
		return v2;
	case 1: // by name (desc)
		if (nox_gui_wol_servers_list.field_1 == &nox_gui_wol_servers_list) {
			return sub_425790(&nox_gui_wol_servers_list, rec);
		}
		v3 = nox_common_list_getFirstSafe_425890(&nox_gui_wol_servers_list);
		if (!v3) {
			nox_common_list_append_4258E0(&nox_gui_wol_servers_list, rec);
			return 0;
		}
		while (nox_strcmpi(rec->server_name, (const char*)v3 + 120) < 0) {
			++v2;
			v3 = nox_common_list_getNextSafe_4258A0(v3);
			if (!v3) {
				nox_common_list_append_4258E0(&nox_gui_wol_servers_list, rec);
				return v2;
			}
		}
		nox_common_list_append_4258E0((int)v3, rec);
		return v2;
	case 2: // by players (asc)
		rec->sort_key = rec->players;
		return sub_425790(&nox_gui_wol_servers_list, rec);
	case 3: // by players (desc)
		rec->sort_key = 32 - rec->players;
		return sub_425790(&nox_gui_wol_servers_list, rec);
	case 4: // by mode (asc)
		if (nox_gui_wol_servers_list.field_1 == &nox_gui_wol_servers_list) {
			return sub_425790(&nox_gui_wol_servers_list, rec);
		}
		v6 = nox_gui_wol_gameModeString_43BCB0(rec->flags);
		v3 = nox_common_list_getFirstSafe_425890(&nox_gui_wol_servers_list);
		if (!v3) {
			nox_common_list_append_4258E0(&nox_gui_wol_servers_list, rec);
			return 0;
		}
		while (1) {
			v7 = nox_gui_wol_gameModeString_43BCB0(*(uint16_t*)((char*)v3 + 163));
			if (nox_wcscmp(v6, v7) <= 0) {
				nox_common_list_append_4258E0((int)v3, rec);
				return v2;
			}
			++v2;
			v3 = nox_common_list_getNextSafe_4258A0(v3);
			if (!v3) {
				nox_common_list_append_4258E0(&nox_gui_wol_servers_list, rec);
				return v2;
			}
		}
	case 5: // by mode (desc)
		if (nox_gui_wol_servers_list.field_1 == &nox_gui_wol_servers_list) {
			return sub_425790(&nox_gui_wol_servers_list, rec);
		}
		v8 = nox_gui_wol_gameModeString_43BCB0(rec->flags);
		v3 = nox_common_list_getFirstSafe_425890(&nox_gui_wol_servers_list);
		if (!v3) {
			nox_common_list_append_4258E0(&nox_gui_wol_servers_list, rec);
			return 0;
		}
		while (1) {
			v9 = nox_gui_wol_gameModeString_43BCB0(*(uint16_t*)((char*)v3 + 163));
			if (nox_wcscmp(v8, v9) >= 0) {
				break;
			}
			++v2;
			v3 = nox_common_list_getNextSafe_4258A0(v3);
			if (!v3) {
				nox_common_list_append_4258E0(&nox_gui_wol_servers_list, rec);
				return v2;
			}
		}
		nox_common_list_append_4258E0((int)v3, rec);
		return v2;
	case 6: // by ping (asc)
		rec->sort_key = rec->ping;
		return sub_425790(&nox_gui_wol_servers_list, rec);
	case 7: // by ping (desc)
		rec->sort_key = 1000 - rec->ping;
		return sub_425790(&nox_gui_wol_servers_list, rec);
	case 8: // by status (asc)
		rec->sort_key = rec->status & 0x30;
		return sub_425790(&nox_gui_wol_servers_list, rec);
	case 9: // by status (desc)
		rec->sort_key = 48 - (rec->status & 0x30);
		return sub_425790(&nox_gui_wol_servers_list, rec);
	default:
		return 0;
	}
}

//----- (004A0290) --------------------------------------------------------
void nox_wol_servers_sortBtnHandler_4A0290(int id) {
	switch (id - 10047) {
	case 0:
		nox_wol_servers_sorting_166704 = (nox_wol_servers_sorting_166704 == 0) + 0;
		break;
	case 1:
		nox_wol_servers_sorting_166704 = (nox_wol_servers_sorting_166704 == 2) + 2;
		break;
	case 2:
		nox_wol_servers_sorting_166704 = (nox_wol_servers_sorting_166704 == 4) + 4;
		break;
	case 3:
		nox_wol_servers_sorting_166704 = (nox_wol_servers_sorting_166704 == 6) + 6;
		break;
	case 4:
		nox_wol_servers_sorting_166704 = (nox_wol_servers_sorting_166704 == 8) + 8;
		break;
	}
}

//----- (004A0360) --------------------------------------------------------
int* sub_4A0360() {
	int* result; // eax
	int* i;      // esi

	result = nox_common_list_getFirstSafe_425890(&nox_gui_wol_servers_list);
	for (i = result; result; i = result) {
		nox_gui_wol_newServerLine_43B7C0(i);
		result = nox_common_list_getNextSafe_4258A0(i);
	}
	return result;
}

//----- (004A0390) --------------------------------------------------------
int* sub_4A0390() {
	uint32_t* v0; // ecx
	int* v1;      // esi
	int* v2;      // edi
	int v4[3];    // [esp+0h] [ebp-Ch]
	uint32_t* v5; // [esp+4h] [ebp-8h]

	nox_common_list_clear_425760(&v4);
	v0 = nox_gui_wol_servers_list.field_1;
	v4[0] = nox_gui_wol_servers_list.field_0;
	v5 = nox_gui_wol_servers_list.field_1;
	if (nox_gui_wol_servers_list.field_0) {
		*(uint32_t*)((uint32_t)nox_gui_wol_servers_list.field_0 + 4) = &v4;
		v0 = v5;
	}
	if (v0) {
		*v0 = &v4;
	}
	nox_common_list_clear_425760(&nox_gui_wol_servers_list);
	v1 = nox_common_list_getFirstSafe_425890(&v4);
	if (v1) {
		do {
			v2 = nox_common_list_getNextSafe_4258A0(v1);
			nox_wol_servers_addResult_4A0030(v1);
			v1 = v2;
		} while (v2);
	}
	return sub_4A0360();
}

//----- (004A0410) --------------------------------------------------------
int sub_4A0410(const char* a1, short a2) {
	int* v2; // edi

	v2 = nox_common_list_getFirstSafe_425890(&nox_gui_wol_servers_list);
	if (!v2) {
		return 1;
	}
	while (strcmp(a1, (const char*)v2 + 12) || a2 != *(uint16_t*)((char*)v2 + 109)) {
		v2 = nox_common_list_getNextSafe_4258A0(v2);
		if (!v2) {
			return 1;
		}
	}
	return 0;
}

//----- (004A0490) --------------------------------------------------------
int* sub_4A0490(int a1) {
	int* result; // eax

	result = nox_common_list_getFirstSafe_425890(&nox_gui_wol_servers_list);
	if (!result) {
		return 0;
	}
	while (result[9] != a1) {
		result = nox_common_list_getNextSafe_4258A0(result);
		if (!result) {
			return 0;
		}
	}
	return result;
}

//----- (004A04C0) --------------------------------------------------------
int* sub_4A04C0(int a1) {
	int v1;      // esi
	int* result; // eax

	v1 = 0;
	result = nox_common_list_getFirstSafe_425890(&nox_gui_wol_servers_list);
	if (!result) {
		return 0;
	}
	while (a1 != v1) {
		++v1;
		result = nox_common_list_getNextSafe_4258A0(result);
		if (!result) {
			return 0;
		}
	}
	return result;
}
