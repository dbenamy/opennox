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
extern uint32_t dword_5d4594_1200776;
extern uint32_t dword_5d4594_1200796;
extern uint32_t nox_server_sanctuaryHelp_54276;
extern uint32_t dword_5d4594_1305748;
extern uint32_t dword_5d4594_1197352;
extern uint32_t dword_5d4594_1197356;
extern uint32_t dword_5d4594_1193712;
extern uint32_t nox_server_connectionType_3596;
extern void* nox_alloc_pixelSpan_1301844;
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

// 49A025: variable 'v11' is possibly undefined
// 49A025: variable 'v12' is possibly undefined

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

// 49C248: variable 'v2' is possibly undefined
// 49C256: variable 'v4' is possibly undefined

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


//----- (0049FF20) --------------------------------------------------------


//----- (0049FFA0) --------------------------------------------------------


//----- (004A0020) --------------------------------------------------------


//----- (004A0030) --------------------------------------------------------


//----- (004A0290) --------------------------------------------------------


//----- (004A0360) --------------------------------------------------------


//----- (004A0390) --------------------------------------------------------


//----- (004A0410) --------------------------------------------------------


//----- (004A0490) --------------------------------------------------------


//----- (004A04C0) --------------------------------------------------------
