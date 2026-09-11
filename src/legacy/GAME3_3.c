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
#include "GAME4_1.h"
#include "GAME4_2.h"
#include "GAME4_3.h"
#include "GAME5.h"
#include "GAME5_2.h"
#include "common__net_list.h"
#include "common__random.h"
#include "common__system__team.h"
#include "operators.h"
#include "server__ability__ability.h"
#include "server__magic__plyrspel.h"
#include "server__object__health.h"
#include "server__system__server.h"

#include "common__gamemech__pausefx.h"

#include "client__gui__window.h"
#include "client__video__draw_common.h"
#include "common/fs/nox_fs.h"
#include "common__crypt.h"
#include "common__magic__speltree.h"
#include "server__script__script.h"
#include "server__script__activator.h"

extern uint32_t nox_xxx_respawnAllow_587000_205200;
extern uint32_t dword_5d4594_1567960;
extern uint32_t dword_5d4594_1568280;
extern uint32_t dword_5d4594_1568288;
extern uint32_t dword_5d4594_1563320;
extern uint32_t dword_5d4594_1568308;
extern uint32_t dword_5d4594_1567988;
extern uint32_t dword_5d4594_1565628;
extern uint32_t dword_5d4594_1565632;
extern uint32_t dword_5d4594_1565520;
extern uint32_t nox_server_needInitNetCodeCache;
extern uint32_t dword_5d4594_1565516;
extern uint32_t dword_5d4594_1567928;
extern void* nox_alloc_respawn_1568020;
extern uint32_t dword_5d4594_1565616;

extern uint32_t dword_5d4594_2649712;
extern uint64_t qword_581450_10176;
extern uint64_t qword_581450_10256;
extern uint64_t qword_5d4594_1567940;
extern uint64_t qword_581450_9544;
extern uint32_t dword_5d4594_1568024;
extern uint32_t dword_5d4594_1565512;
extern uint32_t dword_5d4594_2650652;
extern unsigned int gameex_flags;

//----- (004E3CA0) --------------------------------------------------------
double sub_4E3CA0() { return *getMemFloatPtr(0x587000, 202024); }

//----- (004E3CB0) --------------------------------------------------------
int sub_4E3CB0(float a1) {
	int result; // eax

	result = LODWORD(a1);
	*getMemFloatPtr(0x587000, 202024) = a1;
	return result;
}

//----- (004E3CC0) --------------------------------------------------------
int nox_game_getQuestStage_4E3CC0() { return *getMemU32Ptr(0x587000, 202028); }

//----- (004E3CD0) --------------------------------------------------------
void nox_game_setQuestStage_4E3CD0(int a1) { *getMemU32Ptr(0x587000, 202028) = a1; }

//----- (004E3CE0) --------------------------------------------------------
int nox_xxx_player_4E3CE0() {
	int v0; // ebp
	int v1; // esi
	int v2; // edi

	v0 = 0;
	v1 = nox_xxx_getFirstPlayerUnit_4DA7C0();
	if (!v1) {
		return 0;
	}
	do {
		v2 = *(uint32_t*)(v1 + 748);
		if ((!nox_common_gameFlags_check_40A5C0(1) ||
			 !nox_common_getEngineFlag(NOX_ENGINE_FLAG_DISABLE_GRAPHICS_RENDERING) ||
			 *(uint8_t*)(*(uint32_t*)(v2 + 276) + 2064) != 31) &&
			*(uint32_t*)(*(uint32_t*)(v2 + 276) + 4792) == 1) {
			++v0;
		}
		v1 = nox_xxx_getNextPlayerUnit_4DA7F0(v1);
	} while (v1);
	return v0;
}

//----- (004E3D50) --------------------------------------------------------
int sub_4E3D50() {
	int v0;   // esi
	float v2; // [esp+0h] [ebp-8h]

	if (!*getMemU32Ptr(0x5D4594, 1563928)) {
		*getMemFloatPtr(0x5D4594, 1563912) = nox_xxx_gamedataGetFloat_419D40("PlayerDifficultyDelta");
		*getMemU32Ptr(0x5D4594, 1563928) = 1;
	}
	v0 = nox_xxx_player_4E3CE0();
	v2 = (double)(unsigned int)nox_game_getQuestStage_4E3CC0() *
		 ((double)(v0 - 1) * *getMemFloatPtr(0x5D4594, 1563912) + 1.0);
	return sub_4E3CB0(v2);
}

//----- (004E3DD0) --------------------------------------------------------
void* nox_xxx_objectTypeByIndHealthData(int a1);
short sub_4E3DD0() {
	int v0;             // eax
	int v1;             // esi
	int v2;             // ebx
	int v3;             // ecx
	short v4;           // cx
	int v5;             // edi
	unsigned short v6;  // bp
	unsigned short v7;  // ax
	short v8;           // cx
	int v9;             // ebp
	int v10;            // ecx
	short v11;          // di
	short v12;          // ax
	int v13;            // ecx
	float v15;          // [esp+0h] [ebp-1Ch]
	float v16;          // [esp+0h] [ebp-1Ch]
	float v17;          // [esp+0h] [ebp-1Ch]
	float v18;          // [esp+0h] [ebp-1Ch]
	float v19;          // [esp+Ch] [ebp-10h]
	float v20;          // [esp+Ch] [ebp-10h]
	float v21;          // [esp+10h] [ebp-Ch]
	float v22;          // [esp+14h] [ebp-8h]
	unsigned short v23; // [esp+18h] [ebp-4h]

	v21 = nox_xxx_gamedataGetFloat_419D40("GeneratorMaxHealth");
	v23 = nox_float2int(v21);
	if (!*getMemU32Ptr(0x5D4594, 1563932)) {
		*getMemFloatPtr(0x5D4594, 1563908) = nox_xxx_gamedataGetFloat_419D40("PlayerDamageDiffInit");
		*getMemFloatPtr(0x5D4594, 1563916) = nox_xxx_gamedataGetFloat_419D40("SystemHealthDiffInit");
		*getMemFloatPtr(0x5D4594, 1563920) = nox_xxx_gamedataGetFloat_419D40("PlayerDamageDiffCoeff");
		*getMemFloatPtr(0x5D4594, 1563924) = nox_xxx_gamedataGetFloat_419D40("SystemHealthDiffCoeff");
		*getMemU32Ptr(0x5D4594, 1563932) = 1;
	}
	v19 = (sub_4E3CA0() - 1.0) * *getMemFloatPtr(0x5D4594, 1563920) + *getMemFloatPtr(0x5D4594, 1563908);
	sub_4E4080(v19);
	v20 = (sub_4E3CA0() - 1.0) * *getMemFloatPtr(0x5D4594, 1563924) + *getMemFloatPtr(0x5D4594, 1563916);
	sub_4E40C0(v20);
	v0 = nox_server_getFirstObject_4DA790();
	v1 = v0;
	if (v0) {
		do {
			v2 = nox_server_getNextObject_4DA7A0(v1);
			v0 = *(uint32_t*)(v1 + 8);
			if (v0 & 0x20000 && (v3 = *(uint32_t*)(v1 + 16), (v3 & 0x8000) == 0)) {
				v0 = *(uint32_t*)(v1 + 556);
				v4 = *(uint16_t*)(v0 + 4);
				if (v4) {
					LOWORD(v0) = *(uint16_t*)v0;
					if ((uint16_t)v0) {
						if ((uint16_t)v0 == v4) {
							v5 = nox_xxx_objectTypeByIndHealthData(*(unsigned short*)(v1 + 4));
							v15 = sub_4E40F0() * (double)*(unsigned short*)(v5 + 4);
							v6 = nox_float2int16_abs(v15);
							v16 = sub_4E40F0() * (double)*(unsigned short*)(v5);
							v7 = nox_float2int16_abs(v16);
							if (!v7) {
								v7 = 1;
							}
							if (!v6) {
								v6 = 1;
							}
							if (v7 > v23) {
								v7 = v23;
							}
							if (v6 > v23) {
								v6 = v23;
							}
							LOWORD(v0) = nox_xxx_unitSetHP_4E4560(v1, v7);
							*(uint16_t*)(*(uint32_t*)(v1 + 556) + 4) = v6;
						}
					}
				}
			} else if (v0 & 2) {
				v0 = *(uint32_t*)(v1 + 16);
				if ((v0 & 0x8000) == 0) {
					v0 = *(uint32_t*)(v1 + 556);
					v8 = *(uint16_t*)(v0 + 4);
					if (v8) {
						LOWORD(v0) = *(uint16_t*)v0;
						if ((uint16_t)v0) {
							if ((uint16_t)v0 == v8) {
								int v0a = nox_xxx_objectTypeByIndHealthData(*(unsigned short*)(v1 + 4));
								v9 = *(uint32_t*)(v1 + 748);
								v10 = *(uint32_t*)(v9 + 484);
								LOWORD(v0) = v10 ? *(uint16_t*)(v10 + 72) : *(uint16_t*)(v0a + 4);
								if ((signed char)*(uint8_t*)(v9 + 1440) >= 0) {
									v22 = (double)(unsigned short)v0;
									v17 = sub_4E40F0() * v22;
									v11 = nox_float2int16_abs(v17);
									v18 = sub_4E40F0() * v22;
									v12 = nox_float2int16_abs(v18);
									if (!v12) {
										v12 = 1;
									}
									if (!v11) {
										v11 = 1;
									}
									nox_xxx_unitSetHP_4E4560(v1, v12);
									v0 = v9 + 412;
									v13 = 32;
									*(uint16_t*)(*(uint32_t*)(v1 + 556) + 4) = v11;
									do {
										v0 += 2;
										--v13;
										*(uint16_t*)(v0 - 2) = **(uint16_t**)(v1 + 556);
									} while (v13);
								}
							}
						}
					}
				}
			}
			v1 = v2;
		} while (v2);
	}
	return v0;
}

//----- (004E4080) --------------------------------------------------------
void sub_4E4080(float a1) {
	double v1; // st7

	v1 = nox_xxx_gamedataGetFloat_419D40("PlayerDamageCap");
	if (a1 <= v1) {
		*getMemFloatPtr(0x587000, 202032) = a1;
	} else {
		*getMemFloatPtr(0x587000, 202032) = v1;
	}
}

//----- (004E40B0) --------------------------------------------------------
double sub_4E40B0() { return *getMemFloatPtr(0x587000, 202032); }

//----- (004E40C0) --------------------------------------------------------
void sub_4E40C0(float a1) {
	double v1; // st7

	v1 = nox_xxx_gamedataGetFloat_419D40("SystemHealthCap");
	if (a1 <= v1) {
		*getMemFloatPtr(0x587000, 202036) = a1;
	} else {
		*getMemFloatPtr(0x587000, 202036) = v1;
	}
}

//----- (004E40F0) --------------------------------------------------------
double sub_4E40F0() { return *getMemFloatPtr(0x587000, 202036); }

//----- (004E4100) --------------------------------------------------------
int sub_4E4100() {
	unsigned int v0; // ebp
	int v1;          // esi
	int v2;          // edi
	int result;      // eax

	v0 = 0;
	v1 = nox_xxx_getFirstPlayerUnit_4DA7C0();
	if (!v1) {
		return 1;
	}
	do {
		v2 = *(uint32_t*)(v1 + 748);
		if (!nox_common_gameFlags_check_40A5C0(1) ||
			!nox_common_getEngineFlag(NOX_ENGINE_FLAG_DISABLE_GRAPHICS_RENDERING) ||
			*(uint8_t*)(*(uint32_t*)(v2 + 276) + 2064) != 31) {
			if (*(uint32_t*)(*(uint32_t*)(v2 + 276) + 4792)) {
				++v0;
			}
		}
		result = nox_xxx_getNextPlayerUnit_4DA7F0(v1);
		v1 = result;
	} while (result);
	if (v0 < 6) {
		return 1;
	}
	return result;
}

//----- (004E41B0) --------------------------------------------------------
int sub_4E41B0(char* a1) {
	FILE* v1; // edi

	v1 = nox_fs_open_text(a1);
	if (v1) {
		do {
			while (1) {
				do {
					if (!nox_fs_fgets(v1, (char*)getMemAt(0x5D4594, 1563936), 1024)) {
						nox_fs_close(v1);
						return 1;
					}
				} while (!strtok((char*)getMemAt(0x5D4594, 1563936), "\r\t\n"));
				if (strcmp((const char*)getMemAt(0x5D4594, 1563936), "[Banned Users]")) {
					break;
				}
				if (!sub_4E42C0(v1)) {
					goto LABEL_11;
				}
			}
		} while (strcmp((const char*)getMemAt(0x5D4594, 1563936), "[Allowed Users]") || sub_4E4390(v1));
	LABEL_11:
		nox_fs_close(v1);
	}
	return 0;
}

//----- (004E42C0) --------------------------------------------------------
int sub_4E42C0(FILE* a1) {
	char* v1;       // eax
	char* v2;       // eax
	wchar2_t v4[26]; // [esp+Ch] [ebp-34h]

	while (1) {
		if (!nox_fs_fgets(a1, (char*)getMemAt(0x5D4594, 1563936), 1024)) {
			return 1;
		}
		v1 = strtok((char*)getMemAt(0x5D4594, 1563936), "\r\t\n");
		if (!v1) {
			return 1;
		}
		nox_swprintf(v4, L"%S", v1);
		if (!nox_fs_fgets(a1, (char*)getMemAt(0x5D4594, 1563936), 1024)) {
			return 1;
		}
		v2 = strtok((char*)getMemAt(0x5D4594, 1563936), "\r\t\n");
		if (!v2) {
			break;
		}
		if (!strcmp(v2, "0")) {
			sub_416770(0, v4, 0);
		} else {
			sub_416770(0, v4, v2);
		}
	}
	return 0;
}

//----- (004E4390) --------------------------------------------------------
int sub_4E4390(FILE* a1) {
	char* v1;       // eax
	wchar2_t v3[26]; // [esp+4h] [ebp-34h]

	while (nox_fs_fgets(a1, (char*)getMemAt(0x5D4594, 1563936), 1024)) {
		v1 = strtok((char*)getMemAt(0x5D4594, 1563936), "\r\t\n");
		if (!v1) {
			break;
		}
		nox_swprintf(v3, L"%S", v1);
		sub_4168A0(v3);
	}
	return 1;
}

//----- (004E43F0) --------------------------------------------------------
FILE* sub_4E43F0(char* a1) {
	FILE* result; // eax
	FILE* v2;     // edi
	int* i;       // esi
	int* j;       // esi

	result = nox_fs_create_text(a1);
	v2 = result;
	if (result) {
		nox_fs_fprintf(result, "%s\n", getMemAt(0x587000, 202212));
		for (i = sub_416900(); i; i = sub_416910(i)) {
			if (!*((uint64_t*)i + 8)) {
				nox_fs_fprintf(v2, "%S\n", i + 3);
				if (*((uint8_t*)i + 72)) {
					nox_fs_fprintf(v2, "%s\n", i + 18);
				} else {
					nox_fs_fprintf(v2, "0\n");
				}
			}
		}
		nox_fs_fprintf(v2, "\n%s\n", getMemAt(0x587000, 202228));
		for (j = sub_4168E0(); j; j = sub_4168F0(j)) {
			nox_fs_fprintf(v2, "%S\n", j + 3);
		}
		nox_fs_close(v2);
		result = (FILE*)1;
	}
	return result;
}

//----- (004E4DE0) --------------------------------------------------------
int sub_4E4DE0() {
	int v0;            // edi
	unsigned char* v1; // esi
	int result;        // eax

	if (*getMemU32Ptr(0x5D4594, 1565508)) {
		nox_free_alloc_class(*(void**)getMemAt(0x5D4594, 1565508));
	}
	*getMemU32Ptr(0x5D4594, 1565508) = 0;
	memset(getMemAt(0x5D4594, 1565524), 0, 0x40u);
	memset(getMemAt(0x5D4594, 1565124), 0, 384);
	dword_5d4594_1565512 = 0;
	dword_5d4594_1565516 = 0;
	v0 = 0;
	v1 = getMemAt(0x5D4594, 1565124);
	do {
		v1[1] = 2;
		*v1 = 1;
		v1[2] = nox_xxx_rateGet_40A6C0();
		result = sub_4E4E50(v0);
		v1 += 12;
		++v0;
	} while ((int)v1 < (int)getMemAt(0x5D4594, 1565508));
	return result;
}

//----- (004E4E50) --------------------------------------------------------
int sub_4E4E50(int a1) {
	unsigned int v1;   // esi
	unsigned char v2;  // bl
	unsigned char* v3; // ecx
	int result;        // eax

	v1 = 1;
	if (dword_5d4594_2650652 == 1) {
		v1 = nox_xxx_rateGet_40A6C0();
	}
	v2 = getMemByte(0x5D4594, 1565125 + 12 * a1);
	v3 = getMemAt(0x5D4594, 1565124 + 12 * a1);
	if (v2 > 2u) {
		*((uint32_t*)v3 + 2) = *v3 * (v2 - 1) * (gameFPS() / v1);
	} else {
		*((uint32_t*)v3 + 2) = 0;
	}
	result = v2 * *v3 * (gameFPS() / v1);
	*((uint32_t*)v3 + 1) = result;
	return result;
}

//----- (004E4ED0) --------------------------------------------------------
int sub_4E4ED0() {
	int result; // eax

	result = 0;
	memset(getMemAt(0x5D4594, 1565524), 0, 0x40u);
	return result;
}

//----- (004E4EF0) --------------------------------------------------------
int sub_4E4EF0() {
	unsigned char* v0; // esi
	int v1;            // edi
	int result;        // eax

	v0 = getMemAt(0x5D4594, 1565124);
	memset(getMemAt(0x5D4594, 1565124), 0, 384);
	v1 = 0;
	do {
		v0[1] = 2;
		*v0 = 1;
		v0[2] = nox_xxx_rateGet_40A6C0();
		result = sub_4E4E50(v1);
		v0 += 12;
		++v1;
	} while ((int)v0 < (int)getMemAt(0x5D4594, 1565508));
	return result;
}

//----- (004E4F30) --------------------------------------------------------
int sub_4E4F30(int a1) {
	int result; // eax

	result = a1;
	*getMemU16Ptr(0x5D4594, 1565524 + 2 * a1) = 0;
	return result;
}

//----- (004E4F40) --------------------------------------------------------
int nox_xxx_playerResetImportantCtr_4E4F40(int a1) {
	int v1; // esi

	v1 = 12 * a1;
	*getMemU8Ptr(0x5D4594, 1565125 + v1) = 2;
	*getMemU8Ptr(0x5D4594, 1565124 + v1) = 1;
	*getMemU8Ptr(0x5D4594, 1565126 + v1) = nox_xxx_rateGet_40A6C0();
	return sub_4E4E50(a1);
}

//----- (004E4F80) --------------------------------------------------------
int sub_4E4F80() {
	int result;       // eax
	unsigned char v1; // cl
	int v2;           // esi

	result = dword_5d4594_1565512;
	if (dword_5d4594_1565512) {
		do {
			v1 = *(uint8_t*)(result + 251);
			v2 = *(uint32_t*)(result + 408);
			if (v1 >= 0x31u && v1 <= 0x33u) {
				sub_4E4FC0(result);
			}
			result = v2;
		} while (v2);
	}
	return result;
}

//----- (004E4FC0) --------------------------------------------------------
void sub_4E4FC0(int a1) {
	int v1; // ecx
	int v2; // ecx

	v1 = *(uint32_t*)(a1 + 412);
	if (v1) {
		*(uint32_t*)(v1 + 408) = *(uint32_t*)(a1 + 408);
	} else {
		dword_5d4594_1565512 = *(uint32_t*)(a1 + 408);
	}
	v2 = *(uint32_t*)(a1 + 408);
	if (v2) {
		*(uint32_t*)(v2 + 412) = *(uint32_t*)(a1 + 412);
	} else {
		dword_5d4594_1565516 = *(uint32_t*)(a1 + 412);
	}
	nox_alloc_class_free_obj_first(*(unsigned int**)getMemAt(0x5D4594, 1565508), (uint64_t*)a1);
}

//----- (004E5030) --------------------------------------------------------
int nox_xxx_netSendPacket_4E5030(int a1, const void* a2, signed int a3, int a4, int a5, char a6) {
	char* v7;           // eax
	char* v8;           // edx
	uint16_t* v9;       // esi
	char v10;           // cl
	unsigned char* v11; // eax
	unsigned char* v12; // eax
	int v13;            // edi
	int v14;            // ecx

	if (a1 == 255 || (a1 & 0x80u) == 0 || dword_5d4594_2649712 & ~(1 << (a1 & 0x7F))) {
		if (a3 > 150) {
			return 0;
		}
		v7 = *(char**)getMemAt(0x5D4594, 1565508);
		if (!*getMemU32Ptr(0x5D4594, 1565508)) {
			if (nox_common_gameFlags_check_40A5C0(2048)) {
				dword_5d4594_1565520 = 512;
			} else {
				dword_5d4594_1565520 = nox_common_gameFlags_check_40A5C0(1) ? 3072 : 256;
			}
			if (nox_common_gameFlags_check_40A5C0(2048)) {
				v7 = nox_new_alloc_class_dynamic("importantClass", 416, *(int*)&dword_5d4594_1565520);
			} else {
				v7 = nox_new_alloc_class("importantClass", 416, *(int*)&dword_5d4594_1565520);
			}
			*getMemU32Ptr(0x5D4594, 1565508) = v7;
		}
		v8 = (char*)nox_alloc_class_new_obj_zero(v7);
		if (!v8) {
			if (nox_xxx_importantCheckRate_4E52B0() != 1) {
				return 0;
			}
			v8 = (char*)nox_alloc_class_new_obj_zero(*(uint32_t**)getMemAt(0x5D4594, 1565508));
			if (!v8) {
				return 0;
			}
		}
		memcpy(v8 + 251, a2, a3);
		v8[401] = a3;
		*((uint32_t*)v8 + 101) = a4;
		v8[250] = a1;
		*((uint32_t*)v8 + 45) = a5;
		v8[184] = a6;
		*((uint32_t*)v8 + 43) = 0;
		*((uint32_t*)v8 + 42) = 0;
		v8[164] = 0;
		*((uint32_t*)v8 + 44) = dword_5d4594_2649712;
		memset(v8 + 4, 0, 0x80u);
		v9 = v8 + 186;
		*(uint32_t*)v8 = gameFrame();
		memset(v8 + 186, 0, 0x40u);
		if (v8[184]) {
			if (a1 == 255) {
				v10 = 0;
				v11 = getMemAt(0x5D4594, 1565524);
				do {
					if ((1 << v10) & *((uint32_t*)v8 + 44)) {
						*v9 = (*(uint16_t*)v11)++;
					}
					v11 += 2;
					++v10;
					++v9;
				} while ((int)v11 < (int)getMemAt(0x5D4594, 1565588));
			} else if ((a1 & 0x80u) == 0) {
				*(uint16_t*)&v8[2 * a1 + 186] = (*getMemU16Ptr(0x5D4594, 1565524 + 2 * a1))++;
			} else {
				v12 = getMemAt(0x5D4594, 1565524);
				v13 = v8[250] & 0x7F;
				v14 = 0;
				do {
					if (v14 != v13 && (1 << v14) & *((uint32_t*)v8 + 44)) {
						*v9 = (*(uint16_t*)v12)++;
					}
					v12 += 2;
					++v14;
					++v9;
				} while ((int)v12 < (int)getMemAt(0x5D4594, 1565588));
			}
		}
		*((uint32_t*)v8 + 103) = 0;
		*((uint32_t*)v8 + 102) = dword_5d4594_1565512;
		if (dword_5d4594_1565512) {
			*(uint32_t*)(dword_5d4594_1565512 + 412) = v8;
			dword_5d4594_1565512 = v8;
			return 1;
		}
		dword_5d4594_1565516 = v8;
		dword_5d4594_1565512 = v8;
	}
	return 1;
}

//----- (004E52B0) --------------------------------------------------------
int nox_xxx_importantCheckRate_4E52B0() {
	int v0;            // ebp
	int v1;            // edx
	int v2;            // ebx
	unsigned short v3; // si
	unsigned int v4;   // edi
	char v5;           // al
	int v6;            // eax
	short v8[32];      // [esp+10h] [ebp-40h]

	v0 = 0;
	v1 = dword_5d4594_1565512;
	v8[0] = 0;
	memset(&v8[1], 0, 0x3Cu);
	v2 = -1;
	v3 = 0;
	v8[31] = 0;
	v4 = 999999999;
	if (!dword_5d4594_1565512) {
		return 0;
	}
	do {
		v5 = *(uint8_t*)(v1 + 250);
		if (v5 != -1 && v5 >= 0 && v5 != 31) {
			v6 = *(unsigned char*)(v1 + 250);
			if ((unsigned short)++v8[v6] > v3) {
				v2 = v6;
				v3 = v8[v6];
			}
		}
		if (*(uint32_t*)v1 < v4) {
			v4 = *(uint32_t*)v1;
			v0 = v1;
		}
		v1 = *(uint32_t*)(v1 + 408);
	} while (v1);
	if (v2 != -1) {
		nullsub_24(getMemAt(0x587000, 202360));
		nox_xxx_playerKickDueToRate_4E5360(v2);
	}
	if (!v0) {
		return 0;
	}
	sub_4E4FC0(v0);
	return 1;
}
// 4E5AB0: using guessed type void  nullsub_24(uint32_t);

//----- (004E5360) --------------------------------------------------------
char* nox_xxx_playerKickDueToRate_4E5360(int a1) {
	char* result; // eax
	char* v2;     // esi

	result = nox_common_playerInfoFromNum_417090(a1);
	v2 = result;
	if (result) {
		sub_4E55F0(a1);
		result = (char*)nox_xxx_netNeedTimestampStatus_4174F0((int)v2, 128);
	}
	return result;
}

//----- (004E5390) --------------------------------------------------------
int nox_xxx_netSendPacket1_4E5390(int a1, int a2, int a3, int a4, int a5) {
	return nox_xxx_netSendPacket_4E5030(a1, (const void*)a2, a3, a4, a5, 1);
}

//----- (004E53C0) --------------------------------------------------------
int nox_xxx_netClientSend2_4E53C0(int a1, const void* a2, int a3, int a4, int a5) {
	int result; // eax

	if (nox_common_gameFlags_check_40A5C0(1)) {
		nox_netlist_addToMsgListCli_40EBC0(a1, 0, a2, a3);
		result = 1;
	} else if (a1 == 255 || (a1 & 0x80u) != 0) {
		result = 0;
	} else {
		result = nox_xxx_netSendPacket0_4E5420(a1, a2, a3, a4, a5);
	}
	return result;
}

//----- (004E5420) --------------------------------------------------------
int nox_xxx_netSendPacket0_4E5420(int a1, const void* a2, signed int a3, int a4, int a5) {
	return nox_xxx_netSendPacket_4E5030(a1, a2, a3, a4, a5, 0);
}

//----- (004E5450) --------------------------------------------------------
int sub_4E5450(int a1, char* a2, signed int a3, int a4, int a5) {
	char* v5; // edi
	int v6;   // eax
	int v7;   // esi
	char v9;  // [esp+10h] [ebp+8h]

	v5 = a2;
	v9 = *a2;
	v6 = dword_5d4594_1565512;
	if (dword_5d4594_1565512) {
		do {
			v7 = *(uint32_t*)(v6 + 408);
			if (v9 == *(uint8_t*)(v6 + 251)) {
				if (a1 == 255 || (a1 & 0x80u) != 0) {
					sub_4E4FC0(v6);
				} else {
					sub_4E54D0(1 << a1, v6, a1);
				}
			}
			v6 = v7;
		} while (v7);
	}
	return nox_xxx_netSendPacket0_4E5420(a1, v5, a3, a4, a5);
}

//----- (004E54D0) --------------------------------------------------------
void sub_4E54D0(int a1, int a2, int a3) {
	int v3;  // ecx
	char v4; // dl
	int v5;  // esi
	char v6; // cl
	int v7;  // ecx
	int v8;  // edx
	int v9;  // ecx

	v3 = *(uint32_t*)(a2 + 404);
	if (v3) {
		v4 = *(uint8_t*)(a2 + 251);
		if (v4 == 49 || v4 == 50 || v4 == 51) {
			*(uint32_t*)(v3 + 148) &= ~a1;
		}
	}
	v5 = dword_5d4594_2649712 & *(uint32_t*)(a2 + 176);
	v6 = *(uint8_t*)(a2 + 250);
	if (v6 == -1) {
		v7 = a1 | *(uint32_t*)(a2 + 168);
		*(uint32_t*)(a2 + 168) = v7;
		if ((v5 & v7) == v5) {
			sub_4E4FC0(a2);
		}
	} else if (v6 >= 0) {
		if (*(unsigned char*)(a2 + 250) == a3) {
			sub_4E4FC0(a2);
		}
	} else {
		v8 = 1 << (v6 & 0x7F);
		v9 = a1 | *(uint32_t*)(a2 + 168);
		*(uint32_t*)(a2 + 168) = v9;
		if ((v5 & ~v8 & v9) == (v5 & ~v8)) {
			sub_4E4FC0(a2);
		}
	}
}

//----- (004E55A0) --------------------------------------------------------
int nox_net_importantACK_4E55A0(int a1, int a2) {
	int result; // eax
	int v3;     // edi

	result = dword_5d4594_1565512;
	if (dword_5d4594_1565512) {
		do {
			v3 = *(uint32_t*)(result + 408);
			if (*(uint32_t*)(result + 4 * a1 + 4) == a2) {
				sub_4E54D0(1 << a1, result, a1);
			}
			result = v3;
		} while (v3);
	}
	return result;
}

//----- (004E55F0) --------------------------------------------------------
int sub_4E55F0(unsigned char a1) {
	int result; // eax
	int v2;     // esi

	result = dword_5d4594_1565512;
	if (dword_5d4594_1565512) {
		do {
			v2 = *(uint32_t*)(result + 408);
			sub_4E54D0(1 << a1, result, a1);
			result = v2;
		} while (v2);
	}
	return result;
}

//----- (004E5670) --------------------------------------------------------
unsigned int nox_xxx_importantCheckRate2_4E5670(unsigned char a1) {
	int v1;              // ebp
	int v2;              // eax
	int v3;              // ecx
	char v4;             // al
	unsigned char* v5;   // esi
	unsigned char v6;    // al
	unsigned int result; // eax
	int v8;              // eax
	unsigned char v9;    // al

	v1 = 0;
	nox_common_playerInfoFromNum_417090(a1);
	v2 = dword_5d4594_1565516;
	if (dword_5d4594_1565516) {
		do {
			v3 = *(uint32_t*)(v2 + 412);
			if (!(*(uint32_t*)(v2 + 168) & (1 << a1))) {
				v4 = *(uint8_t*)(v2 + 250);
				if (v4 == -1) {
					++v1;
					goto LABEL_9;
				}
				if (v4 >= 0) {
					if (v4 == a1) {
						++v1;
						goto LABEL_9;
					}
				} else if (a1 != (v4 & 0x7F)) {
					++v1;
					goto LABEL_9;
				}
			}
		LABEL_9:
			v2 = v3;
		} while (v3);
	}
	v5 = getMemAt(0x5D4594, 1565124 + 12 * a1);
	if (nox_xxx_rateGet_40A6C0() != v5[2]) {
		v5[2] = nox_xxx_rateGet_40A6C0();
	}
	if (v1 <= *((uint32_t*)v5 + 1)) {
		v8 = *((uint32_t*)v5 + 2);
		if (v8 > 0 && v1 < v8) {
			if (*v5 == 2) {
				*v5 = 1;
				return sub_4E4E50(a1);
			}
			v9 = v5[1] - 1;
			v5[1] = v9;
			if (v9 < 2u) {
				v5[1] = 2;
			}
		}
		result = sub_4E4E50(a1);
	} else {
		v6 = v5[1] + 1;
		v5[1] = v6;
		if (v6 > 5u) {
			if (*v5 == 2) {
				nox_xxx_playerKickDueToRate_4E5360(a1);
			}
			v5[1] = 5;
			*v5 = 2;
		}
		v5[2] = nox_xxx_rateGet_40A6C0();
		result = sub_4E4E50(a1);
	}
	return result;
}

//----- (004E5770) --------------------------------------------------------
void nox_xxx_netImportant_4E5770(unsigned char a1, int a2) {
	int v2;                                    // edi
	char* v3;                                  // esi
	int v4;                                    // ebp
	int v5;                                    // eax
	char v6;                                   // al
	int v7;                                    // eax
	char v8;                                   // al
	char v9;                                   // cl
	int v10;                                   // eax
	int v11;                                   // eax
	char v12[1];                               // [esp+13h] [ebp-1Dh]
	int (*v13)(int, int, unsigned char*, int); // [esp+14h] [ebp-1Ch]
	int v14;                                   // [esp+18h] [ebp-18h]
	int v15;                                   // [esp+1Ch] [ebp-14h]
	int v16;                                   // [esp+20h] [ebp-10h]
	int v17;                                   // [esp+24h] [ebp-Ch]
	char v18[5];                               // [esp+28h] [ebp-8h]

	v2 = 1 << a1;
	v15 = 1;
	v14 = 0;
	v16 = 1 << a1;
	v3 = nox_common_playerInfoFromNum_417090(a1);
	v13 = nox_netlist_addToMsgListCli_40EBC0;
	if (a1 != 31) {
		v13 = nox_netlist_clientSendWrap_40ECA0;
	}
	if (!v3 || !nox_common_gameFlags_check_40A5C0(1) || v3[3680] & 0x10) {
		v4 = dword_5d4594_1565516;
		if (dword_5d4594_1565516) {
			while (1) {
				v17 = *(uint32_t*)(v4 + 412);
				v5 = *(uint32_t*)(v4 + 404);
				if (v5 && *(uint8_t*)(v5 + 16) & 0x20) {
					*(uint32_t*)(v4 + 404) = 0;
				}
				if (v2 & *(uint32_t*)(v4 + 168)) {
					goto LABEL_39;
				}
				v6 = *(uint8_t*)(v4 + 250);
				if (v6 != -1) {
					if (v6 >= 0) {
						if (v6 != a1) {
							goto LABEL_39;
						}
					} else if (a1 == (v6 & 0x7F)) {
						goto LABEL_39;
					}
				}
				v7 = *(uint32_t*)(v4 + 404);
				if (v7 && !(v2 & *(uint32_t*)(v7 + 148)) || *(uint32_t*)(v4 + 180) && !(v2 & dword_5d4594_2649712)) {
					sub_4E54D0(v2, v4, a1);
					return;
				}
				if (!(v2 & *(uint32_t*)(v4 + 172))) {
					goto LABEL_24;
				}
				v8 = *(uint8_t*)(a1 + v4 + 132);
				if (v8) {
					*(uint8_t*)(a1 + v4 + 132) = v8 - 1;
					goto LABEL_39;
				}
				if (v14 >= getMemByte(0x5D4594, 1565124 + 12 * a1)) {
					goto LABEL_39;
				}
				v9 = *(uint8_t*)(v4 + 164) + 1;
				v10 = v14 + 1;
				*(uint32_t*)(v4 + 172) &= ~v2;
				*(uint8_t*)(v4 + 164) = v9;
				v14 = v10;
			LABEL_24:
				if (v15) {
					if (a2 && a1 != 31) {
						v12[0] = -86;
						if (!v13(a1, a2, v12, 1)) {
							return;
						}
					} else {
						v18[0] = -86;
						*(uint32_t*)&v18[1] = gameFrame();
						if (!v13(a1, a2, v18, 5)) {
							return;
						}
					}
					v15 = 0;
				}
				if (!*(uint8_t*)(v4 + 184)) {
					v11 = v13(a1, a2, (const void*)(v4 + 251), *(unsigned char*)(v4 + 401));
					goto LABEL_36;
				}
				*getMemU8Ptr(0x5D4594, 1564964) = -52;
				*getMemU16Ptr(0x5D4594, 1564965) = *(uint16_t*)(v4 + 2 * a1 + 186);
				*getMemU8Ptr(0x5D4594, 1564967) = *(uint8_t*)(v4 + 401);
				memcpy(getMemAt(0x5D4594, 1564968), (const void*)(v4 + 251), *(unsigned char*)(v4 + 401));
				v11 = v13(a1, a2, getMemAt(0x5D4594, 1564964), *(unsigned char*)(v4 + 401) + 4);
				v2 = v16;
			LABEL_36:
				if (v11) {
					*(uint32_t*)(v4 + 172) |= v2;
					*(uint8_t*)(a1 + v4 + 132) =
						gameFPS() * (unsigned int)getMemByte(0x5D4594, 1565125 + 12 * a1) / nox_xxx_rateGet_40A6C0();
					*(uint32_t*)(v4 + 4 * a1 + 4) = gameFrame();
					if (nox_common_getEngineFlag(NOX_ENGINE_FLAG_REPLAY_READ)) {
						sub_4E54D0(v2, v4, a1);
					}
				}
			LABEL_39:
				v4 = v17;
				if (!v17) {
					goto LABEL_40;
				}
			}
		}
	LABEL_40:
		if (nox_common_gameFlags_check_40A5C0(1) &&
			!(gameFrame() % (gameFPS() * (unsigned int)getMemByte(0x5D4594, 1565125 + 12 * a1)))) {
			nox_xxx_importantCheckRate2_4E5670(a1);
		}
	}
}

//----- (004E5F40) --------------------------------------------------------
int sub_4E5F40(int a1) {
	int result; // eax
	int i;      // esi

	if (!*getMemU32Ptr(0x5D4594, 1565600)) {
		*getMemU32Ptr(0x5D4594, 1565600) = nox_xxx_getNameId_4E3AA0("Glyph");
	}
	result = nox_server_getFirstObject_4DA790();
	for (i = result; result; i = result) {
		if (nox_xxx_unitHasThatParent_4EC4F0(i, a1) && *(unsigned short*)(i + 4) == *getMemU32Ptr(0x5D4594, 1565600) &&
			!(*(uint8_t*)(i + 16) & 0x20)) {
			nox_xxx_netSendPointFx_522FF0(129, (float2*)(i + 56));
			nox_xxx_delayedDeleteObject_4E5CC0(i);
		}
		result = nox_server_getNextObject_4DA7A0(i);
	}
	return result;
}

//----- (004E5FC0) --------------------------------------------------------
void sub_4E5FC0(int a1) {
	int v1; // ebx
	int v2; // ebp
	int v3; // esi
	int v4; // edi

	v1 = *(uint32_t*)(a1 + 516);
	if (v1) {
		do {
			v2 = *(uint32_t*)(v1 + 512);
			if (nox_xxx_creatureIsMonitored_500CC0(a1, v1)) {
				v3 = nox_xxx_inventoryGetFirst_4E7980(v1);
				if (v3) {
					do {
						v4 = nox_xxx_inventoryGetNext_4E7990(v3);
						nox_xxx_delayedDeleteObject_4E5CC0(v3);
						v3 = v4;
					} while (v4);
				}
				nox_xxx_netSendPointFx_522FF0(129, (float2*)(v1 + 56));
				nox_xxx_delayedDeleteObject_4E5CC0(v1);
			}
			v1 = v2;
		} while (v2);
	}
}

//----- (004E6150) --------------------------------------------------------
nox_object_t* sub_4E6150(nox_playerInfo* a1p) {
	int a1 = a1p;
	int v1;     // eax
	int v2;     // eax
	int i;      // esi
	char* v4;   // eax
	int result; // eax
	char* v6;   // eax

	if (!dword_5d4594_1565616) {
		dword_5d4594_1565616 = nox_xxx_getNameId_4E3AA0("GameBall");
	}
	v1 = *(uint32_t*)(a1 + 3628);
	if (v1) {
		if (*(uint8_t*)(v1 + 8) & 4) {
			v2 = nox_xxx_getNextPlayerUnit_4DA7F0(*(uint32_t*)(a1 + 3628));
			goto LABEL_7;
		}
	} else if (nox_common_gameFlags_check_40A5C0(64)) {
		v2 = sub_4E6230();
		if (v2) {
			goto LABEL_7;
		}
	}
	v2 = nox_xxx_getFirstPlayerUnit_4DA7C0();
LABEL_7:
	i = v2;
	if (!v2) {
		goto LABEL_16;
	}
	if (*(uint8_t*)(v2 + 8) & 4) {
		while (1) {
			v4 = nox_common_playerInfoGetByID_417040(*(uint32_t*)(i + 36));
			if (!(*(uint8_t*)(i + 16) & 0x20) && !(v4[3680] & 1)) {
				break;
			}
			i = nox_xxx_getNextPlayerUnit_4DA7F0(i);
			if (!i) {
				goto LABEL_16;
			}
		}
	}
	if (i) {
		return i;
	}
LABEL_16:
	result = sub_4E6230();
	if (result) {
		return result;
	}
	for (i = nox_xxx_getFirstPlayerUnit_4DA7C0(); i; i = nox_xxx_getNextPlayerUnit_4DA7F0(i)) {
		v6 = nox_common_playerInfoGetByID_417040(*(uint32_t*)(i + 36));
		if (!(*(uint8_t*)(i + 16) & 0x20) && !(v6[3680] & 1)) {
			break;
		}
	}
	return i;
}

//----- (004E6230) --------------------------------------------------------
int sub_4E6230() {
	int result; // eax

	if (!dword_5d4594_1565616) {
		dword_5d4594_1565616 = nox_xxx_getNameId_4E3AA0("GameBall");
	}
	result = nox_server_getFirstObject_4DA790();
	if (!result) {
		return 0;
	}
	while (*(unsigned short*)(result + 4) != dword_5d4594_1565616) {
		result = nox_server_getNextObject_4DA7A0(result);
		if (!result) {
			return 0;
		}
	}
	return result;
}

//----- (004E6280) --------------------------------------------------------
nox_object_t* nox_xxx_playerObserverFindGoodSlave0_4E6280(nox_playerInfo* a1p) {
	int a1 = a1p;
	int result; // eax

	if (*(uint32_t*)(a1 + 3628)) {
		result = nox_xxx_playerObserverFindGoodSlave_4EC420(*(uint32_t*)(a1 + 3628));
	} else {
		result = nox_xxx_playerObserverFindGoodSlave2_4EC3E0(*(uint32_t*)(a1 + 2056));
	}
	if (result) {
		while (*(uint32_t*)(result + 16) & 0x8020) {
			result = nox_xxx_playerObserverFindGoodSlave_4EC420(result);
			if (!result) {
				goto LABEL_9;
			}
		}
		return result;
	}
LABEL_9:
	for (result = nox_xxx_playerObserverFindGoodSlave2_4EC3E0(*(uint32_t*)(a1 + 2056)); result;
		 result = nox_xxx_playerObserverFindGoodSlave_4EC420(result)) {
		if (!(*(uint32_t*)(result + 16) & 0x8020)) {
			break;
		}
	}
	return result;
}

//----- (004E6AA0) --------------------------------------------------------
void nox_xxx_playerLeaveObserver_0_4E6AA0(nox_playerInfo* pl) {
	int a1 = pl;
	int v1;      // esi
	int v2;      // edx
	int v3;      // eax
	uint32_t* i; // esi

	if (a1) {
		v1 = *(uint32_t*)(a1 + 2056);
		if (v1) {
			if (*(int (**)(uint32_t*))(v1 + 744) != nox_xxx_updatePlayerMonsterBot_4FAB20) {
				nox_xxx_playerUnsetStatus_417530(a1, 289);
				nox_xxx_spellBuffOff_4FF5B0(v1, 0);
				v2 = *(uint32_t*)(v1 + 16);
				*(uint32_t*)(v1 + 744) = nox_xxx_updatePlayer_4F8100;
				*(uint32_t*)(v1 + 16) = v2 & 0xFFFFFFBF;
				nox_xxx_monsterMarkUpdate_4E8020(*(uint32_t*)(a1 + 2056));
				if (nox_common_gameFlags_check_40A5C0(16)) {
					if (nox_xxx_CheckGameplayFlags_417DA0(4)) {
						v3 = *((uint32_t*)nox_xxx_getTeamByID_418AB0(
								   *(unsigned char*)(*(uint32_t*)(a1 + 2056) + 52)) +
							   19);
						if (v3) {
							if (!*(uint32_t*)(v3 + 492)) {
								sub_4F3400(*(uint32_t*)(a1 + 2056), v3, 1);
							}
						}
					}
				}
				if (nox_common_gameFlags_check_40A5C0(49152) && !sub_509D80(a1)) {
					sub_509C30(a1);
				}
				if (nox_common_gameFlags_check_40A5C0(4096)) {
					for (i = (uint32_t*)nox_xxx_getFirstPlayerUnit_4DA7C0(); i;
						 i = (uint32_t*)nox_xxx_getNextPlayerUnit_4DA7F0((int)i)) {
						if (*(uint32_t*)(*(uint32_t*)(i[187] + 276) + 4792) == 1) {
							nox_xxx_netReportEnchant_4D8F90(*(unsigned char*)(a1 + 2064), i);
						}
					}
				}
			}
		}
	}
}

int sub_50B510();
//----- (004E80C0) --------------------------------------------------------
int sub_4E80C0(char a1) {
	int result; // eax
	int v2;     // esi
	int v3;     // edx

	result = nox_server_getFirstObject_4DA790();
	if (result) {
		v2 = ~(1 << a1);
		do {
			v3 = v2 & *(uint32_t*)(result + 140);
			*(uint32_t*)(result + 144) &= v2;
			*(uint32_t*)(result + 140) = v3;
			result = nox_server_getNextObject_4DA7A0(result);
		} while (result);
	}
	return result;
}

//----- (004E82C0) --------------------------------------------------------
int sub_4E82C0(unsigned char a1, char a2, char a3, short a4) {
	int v4; // eax

	v4 = 6 * a1;
	*getMemU8Ptr(0x5D4594, 1567740 + v4) = a1;
	*getMemU8Ptr(0x5D4594, 1567741 + v4) = a3;
	*getMemU8Ptr(0x5D4594, 1567742 + v4) = a2;
	*getMemU16Ptr(0x5D4594, 1567744 + v4) = a4;
	return nox_xxx_netSendFlagStatus_4D95A0(255, a1, a2, a3, a4);
}

//----- (004E8310) --------------------------------------------------------
char* sub_4E8310() { return (char*)getMemAt(0x5D4594, 1567736); }

//----- (004E8320) --------------------------------------------------------
unsigned char* sub_4E8320(unsigned char a1) { return getMemAt(0x5D4594, 1567740 + 6 * a1); }

uint32_t nox_xxx_wallFlags(int i);
//----- (004E86E0) --------------------------------------------------------
void sub_4E86E0(int a1, float* a2) {
	int v2;    // ecx
	float* v3; // eax
	double v4; // st7
	double v5; // st6
	double v6; // st5
	double v7; // st4
	float v8;  // [esp+0h] [ebp-8h]
	float v9;  // [esp+4h] [ebp-4h]
	float v10; // [esp+Ch] [ebp+4h]
	float v11; // [esp+Ch] [ebp+4h]
	float v12; // [esp+10h] [ebp+8h]
	float v13; // [esp+10h] [ebp+8h]

	v2 = a1;
	if (a1) {
		v3 = a2;
		if (a2) {
			v12 = *(float*)(a1 + 120);
			v10 = v3[30];
			v8 = v10 + v12;
			v4 = (v12 - v10) / v8;
			v9 = (v10 + v10) / v8;
			v5 = *(float*)(v2 + 84) * v4 + v9 * v3[21];
			v6 = (v10 - v12) / v8;
			v7 = (v12 + v12) / v8;
			v11 = v6 * v3[20] + v7 * *(float*)(v2 + 80);
			v13 = v6 * v3[21] + v7 * *(float*)(v2 + 84);
			*(float*)(v2 + 80) = *(float*)(v2 + 80) * v4 + v9 * v3[20];
			*(float*)(v2 + 84) = v5;
			v3[21] = v13;
			v3[20] = v11;
		}
	}
}

//----- (004E8AC0) --------------------------------------------------------
void nox_xxx_collideDoor_4E8AC0(int a2, int a3) {
	int v2;                 // edi
	int v3;                 // esi
	int v4;                 // eax
	unsigned long long v5;  // rax
	char v6;                // al
	unsigned long long v7;  // rax
	int v8;                 // ebx
	int v9;                 // eax
	int v10;                // ecx
	int v11;                // esi
	double v12;             // st7
	double v13;             // st6
	int v14;                // ecx
	int v15;                // eax
	int v16;                // eax
	unsigned long long v17; // rax
	int a3a[2];             // [esp+10h] [ebp-18h]
	float4 a1;              // [esp+18h] [ebp-10h]
	float v20;              // [esp+2Ch] [ebp+4h]

	v2 = a2;
	v3 = *(uint32_t*)(a2 + 748);
	if (a3 && *(uint32_t*)(v3 + 12) == *(uint32_t*)(v3 + 4)) {
		v4 = *(uint32_t*)(a2 + 508);
		if (v4) {
			if (*(uint32_t*)(a2 + 136) <= gameFrame()) {
				*(uint32_t*)(a2 + 508) = 0;
			} else if (v4 != a3) {
				v5 = nox_platform_get_ticks() - *(uint64_t*)&qword_5d4594_1567940;
				a3a[1] = HIDWORD(v5);
				if (v5 <= 0x5DC) {
					return;
				}
				if (*(uint8_t*)(a2 + 12) & 4) {
					nox_xxx_aud_501960(244, a2, 0, 0);
					nox_xxx_netPriMsgToPlayer_4DA2C0(a3, "objcoll.c:GateLockedMagic", 0);
				} else {
					nox_xxx_aud_501960(240, a2, 0, 0);
					nox_xxx_netPriMsgToPlayer_4DA2C0(a3, "objcoll.c:DoorLockedMagic", 0);
				}
				*(uint64_t*)&qword_5d4594_1567940 = nox_platform_get_ticks();
				return;
			}
		}
		v6 = *(uint8_t*)(v3 + 1);
		if (!v6) {
			return;
		}
		if (v6 == 5) {
			v7 = nox_platform_get_ticks() - *(uint64_t*)&qword_5d4594_1567940;
			a3a[1] = HIDWORD(v7);
			if (v7 <= 0x5DC) {
				return;
			}
			if (*(uint8_t*)(a2 + 12) & 4) {
				nox_xxx_aud_501960(244, a2, 0, 0);
				nox_xxx_netPriMsgToPlayer_4DA2C0(a3, "objcoll.c:GateLockedMechanism", 0);
			} else {
				nox_xxx_aud_501960(240, a2, 0, 0);
				nox_xxx_netPriMsgToPlayer_4DA2C0(a3, "objcoll.c:DoorLockedMechanism", 0);
			}
			*(uint64_t*)&qword_5d4594_1567940 = nox_platform_get_ticks();
			return;
		}
		v8 = nox_xxx_doorGetSomeKey_4E8910(a3, a2);
		if (!v8) {
			v17 = nox_platform_get_ticks() - *(uint64_t*)&qword_5d4594_1567940;
			a3a[1] = HIDWORD(v17);
			if (v17 <= 0x5DC) {
				return;
			}
			if (*(uint8_t*)(a2 + 12) & 4) {
				nox_xxx_aud_501960(244, a2, 0, 0);
				sub_4FADD0(a3, "objcoll.c:GateLockedKey", *(uint8_t*)(v3 + 1));
			} else {
				nox_xxx_aud_501960(240, a2, 0, 0);
				sub_4FADD0(a3, "objcoll.c:DoorLockedKey", *(uint8_t*)(v3 + 1));
			}
			*(uint64_t*)&qword_5d4594_1567940 = nox_platform_get_ticks();
			return;
		}
		v9 = *(uint32_t*)(v3 + 16);
		v10 = *(uint32_t*)(v3 + 20);
		*(uint8_t*)(v3 + 1) = 0;
		v11 = *(uint32_t*)(v3 + 4);
		v12 = (double)(23 * v9);
		a1.field_0 = v12 - 34.0;
		v13 = (double)(23 * v10);
		v20 = v13;
		a1.field_4 = v13 - 34.0;
		a1.field_8 = v12 + 34.0;
		a1.field_C = v20 + 34.0;
		switch (v11) {
		case 0:
			v14 = v10 - 1;
			a3a[0] = v9 - 1;
			a3a[1] = v14;
			break;
		case 8:
			v14 = v10 - 1;
			a3a[0] = v9 + 1;
			a3a[1] = v14;
			break;
		case 16:
			v15 = v9 + 1;
			a3a[0] = v15;
			v14 = v10 + 1;
			a3a[1] = v14;
			break;
		case 24:
			v15 = v9 - 1;
			a3a[0] = v15;
			v14 = v10 + 1;
			a3a[1] = v14;
			break;
		default:
			break;
		}
		if (nox_common_gameFlags_check_40A5C0(4096)) {
			sub_4E8390(v2);
			sub_4D71E0(gameFrame());
		}
		nox_xxx_getUnitsInRect_517C10(&a1, nox_xxx_fnFindCloseDoors_4E8340, (int)a3a);
		nox_xxx_aud_501960(234, v2, 0, 0);
		v16 = *(uint32_t*)(v8 + 492);
		if (v16 && a3 != v16 && *(uint8_t*)(v16 + 8) & 4 && nox_common_gameFlags_check_40A5C0(4096) &&
			sub_4D72C0() == 1) {
			nox_xxx_netPriMsgToPlayer_4DA2C0(*(uint32_t*)(v8 + 492), "GeneralPrint:KeyShared1", 0);
		}
		nox_xxx_delayedDeleteObject_4E5CC0(v8);
	}
}

//----- (004E8DF0) --------------------------------------------------------
int nox_xxx_collidePickup_4E8DF0(int a1, int a2) {
	int result; // eax
	int v3;     // ecx

	result = a2;
	if (a2) {
		v3 = *(uint32_t*)(a2 + 8);
		if (!(v3 & 2) && (unsigned int)(gameFrame() - *(uint32_t*)(a1 + 128)) >= (int)gameFPS() >> 1 &&
			(!(v3 & 4) || *(uint8_t*)(*(uint32_t*)(a2 + 748) + 240) & 1)) {
			result = nox_xxx_inventoryServPlace_4F36F0(a2, a1, 1, 1);
		}
	}
	return result;
}

//----- (004E8E50) --------------------------------------------------------
unsigned char* sub_4E8E50() { return getMemAt(0x5D4594, 1567844); }

//----- (004E8E60) --------------------------------------------------------
int sub_4E8E60() {
	int v0;          // esi
	int v1;          // edi
	int v2;          // ebx
	int v3;          // eax
	int v4;          // edx
	unsigned int v5; // eax
	int v6;          // eax
	int v7;          // ecx
	int result;      // eax
	int v9;          // eax
	int v10;         // ecx
	float v11;       // [esp+0h] [ebp-1Ch]
	float v12;       // [esp+0h] [ebp-1Ch]
	int v13;         // [esp+10h] [ebp-Ch]
	int v14;         // [esp+18h] [ebp-4h]

	v11 = nox_xxx_gamedataGetFloat_419D40("QuestExitTimerStart");
	v0 = nox_float2int(v11);
	v1 = 0;
	v14 = v0;
	v2 = v0;
	v13 = 0;
	if (sub_40A220()) {
		v3 = sub_40A230();
		v4 = (int)v3 / 1000;
		v5 = 0;
		v0 = v5 + v4;
		v14 = v5 + v4;
		v2 = v5 + v4;
	}
	v6 = nox_xxx_getFirstPlayerUnit_4DA7C0();
	if (!v6) {
		return sub_40A1F0(0);
	}
	do {
		v7 = *(uint32_t*)(v6 + 748);
		if (*(uint32_t*)(*(uint32_t*)(v7 + 276) + 4792) == 1) {
			++v1;
			if (*(uint32_t*)(v7 + 312)) {
				++v13;
			}
		}
		v6 = nox_xxx_getNextPlayerUnit_4DA7F0(v6);
	} while (v6);
	if (!v1) {
		return sub_40A1F0(0);
	}
	v12 = (double)v13 / (double)v1 * (double)v14;
	v9 = nox_float2int(v12);
	v10 = v0 - v9;
	if (v0 - v9 < v0 && (v0 -= v9, v2 != v10) || (result = sub_40A300()) == 0) {
		nox_xxx_servStartCountdown_40A2A0(v0, "objcoll.c:ExitCountdown");
		result = nox_xxx_netGauntlet_4D9E70(255);
	}
	return result;
}

//----- (004E8F60) --------------------------------------------------------
bool nox_server_questAllowDefault();
bool nox_server_questMaybeWarp_4E8F60() {
	unsigned int curLvl = nox_game_getQuestStage_4E3CC0();
	unsigned int toLvl = nox_server_questNextStageThreshold_4D74F0(curLvl);
	int cnt = 0;
	bool allow = nox_server_questAllowDefault();
	for (void* unit = nox_xxx_getFirstPlayerUnit_4DA7C0(); unit; unit = nox_xxx_getNextPlayerUnit_4DA7F0(unit)) {
		int v4 = *(uint32_t*)((int)unit + 748);
		if (!nox_common_gameFlags_check_40A5C0(1) ||
			!nox_common_getEngineFlag(NOX_ENGINE_FLAG_DISABLE_GRAPHICS_RENDERING) ||
			*(uint8_t*)(*(uint32_t*)(v4 + 276) + 2064) != 31) {
			int v5 = *(uint32_t*)(v4 + 276);
			if (*(uint32_t*)(v5 + 4792)) {
				++cnt;
				if (!*(uint32_t*)(v4 + 316)) {
					return 0;
				}
				if (*(uint32_t*)(v5 + 4696) >= toLvl) {
					allow = 1;
				}
			}
		}
	}
	if (!cnt) {
		return 0;
	}
	return allow;
}

//----- (004E9010) --------------------------------------------------------
int sub_4E9010() {
	int v0; // ebp
	int v1; // edi
	int v2; // esi

	v0 = 0;
	v1 = nox_xxx_getFirstPlayerUnit_4DA7C0();
	if (v1) {
		while (1) {
			v2 = *(uint32_t*)(v1 + 748);
			if (!nox_common_gameFlags_check_40A5C0(1) ||
				!nox_common_getEngineFlag(NOX_ENGINE_FLAG_DISABLE_GRAPHICS_RENDERING) ||
				*(uint8_t*)(*(uint32_t*)(v2 + 276) + 2064) != 31) {
				if (*(uint32_t*)(*(uint32_t*)(v2 + 276) + 4792)) {
					++v0;
					if (!*(uint32_t*)(v2 + 312)) {
						break;
					}
				}
			}
			v1 = nox_xxx_getNextPlayerUnit_4DA7F0(v1);
			if (!v1) {
				if (!v0) {
					return 0;
				}
				return 1;
			}
		}
	}
	return 0;
}

//----- (004E9090) --------------------------------------------------------
void sub_4DCBF0(int a1);
void nox_xxx_collideExit_4E9090(int a1, int a2, int a3) {
	int v4;           // ebp
	int v5;           // esi
	uint8_t* v6;      // edi
	int i;            // edi
	char v8;          // al
	int v9;           // ebx
	char* v10;        // edi
	int v11;          // edi
	int v12;          // edi
	int v13;          // eax
	unsigned int v14; // edi
	int v15;          // eax
	int j;            // esi
	int v17;          // eax
	int v18;          // eax
	const char* v19;  // [esp+8h] [ebp+8h]

	if (!dword_5d4594_1567960) {
		dword_5d4594_1567960 = nox_xxx_getNameId_4E3AA0("Glyph");
	}
	v4 = a2;
	if (!a2) {
		return;
	}
	if (!(*(uint8_t*)(a2 + 8) & 4)) {
		return;
	}
	v5 = *(uint32_t*)(a2 + 748);
	v6 = *(uint8_t**)(a1 + 700);
	v19 = *(const char**)(a1 + 700);
	if (nox_common_gameFlags_check_40A5C0(4096) && (*(char**)(v5 + 312) != 0 || *(char**)(v5 + 316) != 0)) {
		return;
	}
	if ((*(uint8_t*)(a1 + 12) & 2) && sub_4D75E0() == 0) {
		return;
	}
	if (sub_4DCC90() == 1) {
		return;
	}
	if (!sub_4DCC10(v4)) {
		return;
	}
	if (nox_xxx_checkGameFlagPause_413A50() == 1) {
		return;
	}
	if (!*v6 && nox_common_gameFlags_check_40A5C0(4096) == 0) {
		return;
	}
	if (*(uint8_t*)(*(uint32_t*)(v5 + 276) + 2251) == 1) {
		for (i = *(uint32_t*)(v4 + 516); i; i = *(uint32_t*)(i + 512)) {
			if (*(unsigned short*)(i + 4) == dword_5d4594_1567960 &&
				!*(uint32_t*)(i + 492)) {
				nox_xxx_delayedDeleteObject_4E5CC0(i);
				v8 = *(uint8_t*)(v5 + 244);
				if (v8) {
					*(uint8_t*)(v5 + 244) = v8 - 1;
				}
			}
		}
	}
	sub_4DCBF0(1);
	if (nox_common_gameFlags_check_40A5C0(2048)) {
		nox_setSaveFileName_4DB130("WORKING");
		sub_4DB170(1, a1, 0);
		return;
	}
	v9 = 1;
	if (nox_common_gameFlags_check_40A5C0(4096)) {
		v10 = nox_xxx_getQuestMapFile_4D0F60();
		v19 = v10;
	} else {
		v10 = (char*)v19;
	}
	if (!nox_common_gameFlags_check_40A5C0(4096)) {
		goto LABEL_47;
	}
	v11 = 1;
	do {
		if (nox_common_playerIsAbilityActive_4FC250(v4, v11)) {
			sub_4FC300((uint32_t*)v4, v11);
		}
		++v11;
	} while (v11 < 6);
	v12 = a1;
	v13 = *(uint32_t*)(a1 + 12);
	if (v13 & 1) {
		v14 = nox_game_getQuestStage_4E3CC0() + 1;
		sub_4D60E0(v4);
		v15 = *(uint32_t*)(v5 + 276);
		if (*(uint32_t*)(v15 + 4696) < v14) {
			*(uint32_t*)(v15 + 4696) = v14;
			sub_4D7450(*(unsigned char*)(*(uint32_t*)(v5 + 276) + 2064),
					   *(uint32_t*)(*(uint32_t*)(v5 + 276) + 4696));
		}
		*(uint32_t*)(v5 + 312) = a1;
		*(uint32_t*)(v5 + 316) = 0;
		nox_xxx_playerSetState_4FA020((uint32_t*)v4, 13);
		nox_xxx_playerGoObserver_4E6860(*(uint32_t*)(v5 + 276), 0, 0);
		nox_xxx_netInformTextMsg2_4DA180(18, (uint8_t*)(v4 + 36));
		v12 = a1;
		v9 = sub_4E9010();
	} else if (v13 & 2) {
		*(uint32_t*)(v5 + 312) = 0;
		*(uint32_t*)(v5 + 316) = a1;
		sub_4D75F0(gameFrame());
		nox_xxx_playerSetState_4FA020((uint32_t*)v4, 13);
		nox_xxx_playerGoObserver_4E6860(*(uint32_t*)(v5 + 276), 0, 0);
		for (j = nox_xxx_getFirstPlayerUnit_4DA7C0(); j;
			 j = nox_xxx_getNextPlayerUnit_4DA7F0(j)) {
			if (j != v4) {
				nox_xxx_netInformTextMsg_4DA0F0(
					*(unsigned char*)(*(uint32_t*)(*(uint32_t*)(j + 748) + 276) + 2064), 19,
					(int*)(v4 + 36));
			}
		}
		nox_xxx_netPriMsgToPlayer_4DA2C0(v4, "objcoll.c:PlayerEntersWarp", 0);
		v9 = nox_server_questMaybeWarp_4E8F60();
	}
	nox_xxx_aud_501960(1003, v12, 0, 0);
	if (!(*(uint8_t*)(v12 + 12) & 2) && !v9) {
		sub_4E8E60();
	}
	strcpy((char*)getMemAt(0x5D4594, 1567844), v19);
	if (v9 != 1) {
		return;
	}
	v10 = (char*)v19;
LABEL_47:
	if (v10 && *v10) {
		if (*(uint8_t*)(a1 + 12) & 2) {
			v17 = nox_game_getQuestStage_4E3CC0();
			v18 = nox_server_questNextStageThreshold_4D74F0(v17);
			nox_game_setQuestStage_4E3CD0(v18 - 1);
			sub_4D76E0(1);
			sub_4D60B0();
		}
		nox_xxx_mapLoad_4D2450(v10);
	}
}

//----- (004E9500) --------------------------------------------------------

// Mix Patch new version with fixups

/*void  nox_xxx_spellFlyCollide_4E9500(int a1, int a2, float *a3)
{
  int v3; // esi
  float2 *v4; // edi
  int v5; // ebx
  int v6; // ebp
  int v7; // ecx
  int v8; // ebx
  int v9; // eax
  int v10; // [esp+10h] [ebp-10h]
  int v11[3]; // [esp+14h] [ebp-Ch]

  v3 = a2;
  v4 = (float2 *)a1;
  v5 = *(uint32_t *)(a1 + 748);
  v10 = *(uint32_t *)(a1 + 748);
  if ( a2 )
  {
	if ( *(uint32_t *)(a2 + 16) & 0x8020 )
	  return;
	if ( !(*(uint8_t *)(a2 + 8) & 4) )
	  goto LABEL_23;
	v6 = *(uint32_t *)(a2 + 748);
	if ( *(uint8_t *)(v6 + 88) == 16 && nox_xxx_checkReflectShield_4E6E50((float2 *)(a2 + 56), *(short *)(a2 + 124),
(float2 *)(a1 + 72)) & 1 )
	{
	  nox_xxx_aud_501960(878, v3, 0, 0);
	  nox_xxx_projectileReflect_4E0A70((int)v4, v3);
	  nox_xxx_changeOwner_52BE40((int)v4, v3);
	  return;
	}
	if ( *(uint8_t *)(v6 + 88) == 13 || (!*(uint8_t *)(v6 + 88) && dword_980858[0] & 0x40000 ) )
	{
	  v7 = *(uint32_t *)(*(uint32_t *)(v6 + 276) + 4);
	  if ( v7 & 0x400 )
	  {
	if ( nox_xxx_checkReflectShield_4E6E50((float2 *)(v3 + 56), *(short *)(v3 + 124), v4 + 9) & 1 )
	{
	  v8 = nox_xxx_randInt_415FA0(18, 20);
	  nox_xxx_aud_501960(890, v3, 0, 0);
	  nox_xxx_projectileReflect_4E0A70((int)v4, v3);
	  nox_xxx_changeOwner_52BE40((int)v4, v3);
	  nox_xxx_playerSetState_4FA020((uint32_t *)v3, v8);
	  v9 = nox_xxx_playerMapAnimState_4FA2B0(v3);
	  nox_xxx_animPlayerGetFrameRange_4F9F90(v9, &a1, &a2);
	  v5 = v10;
	  *(uint8_t *)(v6 + 236) = a1 - 1;
	}
	  }
	}
M_LABEL_13:
	if ( nox_xxx_checkInversionEffect_4FA4F0(v3, (int)v4) )
	{
	  nox_xxx_changeOwner_52BE40((int)v4, v3);
	  return;
	}
	if ( nox_xxx_testUnitBuffs_4FF350(v3, 27) && nox_xxx_checkReflectShield_4E6E50((float2 *)(v3 + 56), *(short *)(v3 +
124), v4 + 7) & 1 )
	{
	  nox_xxx_changeOwner_52BE40((int)v4, v3);
	  nox_xxx_aud_501960(122, v3, 0, 0);
	}
	else
	{
LABEL_23:
	  if ( *(uint32_t *)(v5 + 4) == v3 )
	  {
	v11[0] = v3;
	nox_xxx_spellAccept_4FD400(*(uint32_t *)(v5 + 12), *(uint32_t *)(v5 + 8), *(uint32_t **)v5, (int)v4, v11, *(uint32_t
*)(v5 + 16)); nox_xxx_delayedDeleteObject_4E5CC0((int)v4);
	  }
	  else
	  {
	nox_xxx_projectileReflect_4E0A70((int)v4, v3);
	  }
	}
  }
  else if ( a3 )
  {
	nox_xxx_collideReflect_57B810(a3, a1 + 80);
  }
}*/
void nox_xxx_spellFlyCollide_4E9500(int a1, int a2, float* a3) {
	int v3;     // esi
	float2* v4; // edi
	int v5;     // ebx
	int v6;     // ebp
	int v7;     // ecx
	int v8;     // ebx
	int v9;     // eax
	char v10;   // al
	int v11;    // [esp+10h] [ebp-10h]
	int v12[3]; // [esp+14h] [ebp-Ch]

	v3 = a2;
	v4 = (float2*)a1;
	v5 = *(uint32_t*)(a1 + 748);
	v11 = *(uint32_t*)(a1 + 748);
	if (!a2) {
		if (a3) {
			nox_xxx_collideReflect_57B810(a3, a1 + 80);
		}
		return;
	}
	if (*(uint32_t*)(a2 + 16) & 0x8020) {
		return;
	}
	if (*(uint8_t*)(a2 + 8) & 4) {
		v6 = *(uint32_t*)(a2 + 748);
		if (*(uint8_t*)(v6 + 88) != 16) {
			if (*(uint8_t*)(v6 + 88) != 1 || nox_common_mapPlrActionToStateId_4FA2B0(a2) != 45 ||
				(v10 = *(uint8_t*)(*(uint32_t*)(v6 + 276) + 3), v10 != 1) && v10 != 2) {
				goto LABEL_22;
			}
			if (!(gameex_flags & 0x10)) { // BERSERKER_SHIED_BLOCK
				goto LABEL_13;
			}
		}
		if (nox_server_testTwoPointsAndDirection_4E6E50((float2*)(a2 + 56), *(short*)(a2 + 124), (float2*)(a1 + 72)) &
			1) {
			nox_xxx_aud_501960(878, v3, 0, 0);
			nox_xxx_projectileReflect_4E0A70((int)v4, v3);
			nox_xxx_changeOwner_52BE40((int)v4, v3);
			return;
		}
	LABEL_22:
		if (*(uint8_t*)(v6 + 88) == 13 ||
			!*(uint8_t*)(v6 + 88) && gameex_flags & 0x4) // GREAT_SWORD_BLOKING_WALK
		{
			v7 = *(uint32_t*)(*(uint32_t*)(v6 + 276) + 4);
			if (v7 & 0x400 &&
				nox_server_testTwoPointsAndDirection_4E6E50((float2*)(v3 + 56), *(short*)(v3 + 124), v4 + 9) &
					1) {
				v8 = nox_common_randomInt_415FA0(18, 20);
				nox_xxx_aud_501960(890, v3, 0, 0);
				nox_xxx_projectileReflect_4E0A70((int)v4, v3);
				nox_xxx_changeOwner_52BE40((int)v4, v3);
				nox_xxx_playerSetState_4FA020((uint32_t*)v3, v8);
				v9 = nox_common_mapPlrActionToStateId_4FA2B0(v3);
				nox_xxx_animPlayerGetFrameRange_4F9F90(v9, &a1, &a2);
				v5 = v11;
				*(uint8_t*)(v6 + 236) = a1 - 1;
			}
		}
	LABEL_13:
		if (nox_xxx_checkInversionEffect_4FA4F0(v3, (int)v4)) {
			nox_xxx_changeOwner_52BE40((int)v4, v3);
			return;
		}
		if (nox_xxx_testUnitBuffs_4FF350(v3, 27) &&
			nox_server_testTwoPointsAndDirection_4E6E50((float2*)(v3 + 56), *(short*)(v3 + 124), v4 + 7) & 1) {
			nox_xxx_changeOwner_52BE40((int)v4, v3);
			nox_xxx_aud_501960(122, v3, 0, 0);
			return;
		}
	}
	if (*(uint32_t*)(v5 + 4) == v3) {
		v12[0] = v3;
		nox_xxx_spellAccept_4FD400(*(uint32_t*)(v5 + 12), *(uint32_t*)(v5 + 8), *(uint32_t**)v5, (int)v4, v12,
								   *(uint32_t*)(v5 + 16));
		nox_xxx_delayedDeleteObject_4E5CC0((int)v4);
	} else {
		nox_xxx_projectileReflect_4E0A70((int)v4, v3);
	}
}

//----- (004E9C40) --------------------------------------------------------
void nox_xxx_collideChest_4E9C40(uint32_t* a1, int a2) {
	int v2;                // esi
	int v3;                // eax
	int v4;                // eax
	int i;                 // ebp
	void (*v6)(uint32_t*); // eax

	v2 = a2;
	if (!a2) {
		return;
	}
	if ((*(uint8_t*)(a2 + 8) & 4) == 0) {
		return;
	}
	v3 = a1[4];
	if ((v3 & 0x8000) != 0) {
		return;
	}
	if (nox_common_gameFlags_check_40A5C0(4096) && (v4 = a1[3], BYTE1(v4) & 0xF)) {
		for (i = nox_xxx_inventoryGetFirst_4E7980(a2); i; i = nox_xxx_inventoryGetNext_4E7990(i)) {
			if (*(uint8_t*)(i + 8) & 0x40) {
				if (!strcmp((const char*)nox_xxx_getUnitName_4E39D0(i), "SilverKey")) {
					nox_xxx_delayedDeleteObject_4E5CC0(i);
					nox_xxx_aud_501960(234, (int)a1, 0, 0);
					v2 = a2;
					goto LABEL_14;
				}
				v2 = a2;
			}
		}
		if ((unsigned long long)(nox_platform_get_ticks() - *(uint64_t*)&qword_5d4594_1567940) > 0x5DC) {
			nox_xxx_aud_501960(1012, (int)a1, 0, 0);
			nox_xxx_netPriMsgToPlayer_4DA2C0(v2, "objcoll.c:ChestLockedSilver", 0);
			*(uint64_t*)&qword_5d4594_1567940 = nox_platform_get_ticks();
		}
		return;
	}
LABEL_14:
	v6 = (void (*)(uint32_t*))a1[181];
	if (v6) {
		v6(a1);
	}
	nox_xxx_chest_4EDF00((int)a1, v2);
	nox_xxx_dropAllItems_4EDA40(a1);
	return;
}

//----- (004EAAA0) --------------------------------------------------------
void sub_4EAAA0(int a1) {
	if (gameFrame() > (unsigned int)(*(uint32_t*)(a1 + 136) + 3)) {
		*(uint32_t*)(a1 + 136) = gameFrame();
		nox_xxx_aud_501960(281, a1, 0, 0);
	}
}

//----- (004EAAD0) --------------------------------------------------------
void sub_4EAAD0(int a1, int a2) {
	if (a2 && *(uint8_t*)(a2 + 8) & 4) {
		if (gameFrame() > (unsigned int)(*(uint32_t*)(a1 + 136) + 30)) {
			*(uint32_t*)(a1 + 136) = gameFrame();
			nox_xxx_aud_501960(**(uint32_t**)(a1 + 700), a1, 0, 0);
		}
	}
}

//----- (004EAB20) --------------------------------------------------------
int nox_xxx_collidePentagram_4EAB20(int a1) {
	int result; // eax

	result = a1;
	*(uint32_t*)(*(uint32_t*)(a1 + 748) + 4) = 1;
	return result;
}

//----- (004EAB40) --------------------------------------------------------
void nox_xxx_collideSign_4EAB40(int a1, int a2) {

	if (a2) {
		if (*(uint8_t*)(a2 + 8) & 4) {
			(*(int (**)(int, int))(a1 + 732))(a2, a1);
		}
	}
}

//----- (004EAB60) --------------------------------------------------------
void nox_xxx_collideTrapDoor_4EAB60(int a1, int a2) {
	int v2;    // ebx
	int v3;    // eax
	int v4;    // eax
	double v5; // st7

	v2 = *(uint32_t*)(a1 + 700);
	if (a2) {
		v3 = *(uint32_t*)(a2 + 8);
		if ((v3 & 0x80u) == 0) {
			if (*(uint32_t*)(a1 + 16) & 0x1000000) {
				v4 = *(uint32_t*)(a2 + 172);
				if (v4 == 3) {
					if (*(float*)(a1 + 184) < (double)*(float*)(a2 + 184) ||
						*(float*)(a1 + 188) < (double)*(float*)(a2 + 188)) {
						return;
					}
					if (nox_xxx_map_57B850((float2*)(a1 + 56), (float*)(a1 + 172), (float2*)(a2 + 56))) {
						*(uint32_t*)(a2 + 16) |= 0x60000u;
						*(float*)(a2 + 164) = (double)*(int*)(v2 + 8);
						*(float*)(a2 + 168) = (double)*(int*)(v2 + 12);
						*(uint32_t*)(a2 + 156) = *(uint32_t*)(a1 + 56);
						*(uint32_t*)(a2 + 160) = *(uint32_t*)(a1 + 60);
					}
					return;
				}
				if (v4 != 2) {
					if (nox_xxx_map_57B850((float2*)(a1 + 56), (float*)(a1 + 172), (float2*)(a2 + 56))) {
						*(uint32_t*)(a2 + 16) |= 0x60000u;
						*(float*)(a2 + 164) = (double)*(int*)(v2 + 8);
						*(float*)(a2 + 168) = (double)*(int*)(v2 + 12);
						*(uint32_t*)(a2 + 156) = *(uint32_t*)(a1 + 56);
						*(uint32_t*)(a2 + 160) = *(uint32_t*)(a1 + 60);
					}
					return;
				}
				v5 = *(float*)(a2 + 176) + *(float*)(a2 + 176);
				if (v5 <= *(float*)(a1 + 184) && v5 <= *(float*)(a1 + 188)) {
					if (nox_xxx_map_57B850((float2*)(a1 + 56), (float*)(a1 + 172), (float2*)(a2 + 56))) {
						*(uint32_t*)(a2 + 16) |= 0x60000u;
						*(float*)(a2 + 164) = (double)*(int*)(v2 + 8);
						*(float*)(a2 + 168) = (double)*(int*)(v2 + 12);
						*(uint32_t*)(a2 + 156) = *(uint32_t*)(a1 + 56);
						*(uint32_t*)(a2 + 160) = *(uint32_t*)(a1 + 60);
					}
					return;
				}
			} else if (!*(uint32_t*)(v2 + 24) && (!(v3 & 4) || !nox_common_playerIsAbilityActive_4FC250(a2, 4))) {
				if (*(uint16_t*)(v2 + 20)) {
					*(uint32_t*)(v2 + 16) = gameFrame() + *(unsigned short*)(v2 + 20);
				}
				nox_xxx_scriptCallByEventBlock_502490((int*)v2, a2, a1, 20);
				*(uint32_t*)(v2 + 24) = 1;
			}
		}
	}
}

//----- (004EACA0) --------------------------------------------------------
void sub_4EACA0(int a1, int a2) {
	int* v3;          // ebx

	v3 = *(int**)(a1 + 700);
	if (a2) {
		if ((signed char)*(uint8_t*)(a2 + 8) >= 0) {
			nox_xxx_netSendPointFx_522FF0(138, (float2*)(a2 + 56));
			nox_xxx_aud_501960(147, a2, 0, 0);
			*(float*)(a2 + 164) = (double)*v3;
			*(float*)(a2 + 168) = (double)v3[1];
			nox_xxx_teleportToMB_4E7190((uint8_t*)a2, (float*)(a2 + 164));
			nox_xxx_netSendPointFx_522FF0(137, (float2*)(a2 + 56));
			nox_xxx_aud_501960(147, a2, 0, 0);
		}
	}
}

//----- (004EAD20) --------------------------------------------------------
int nox_xxx_collideSpellPedestal_4EAD20(int a1, int a2) {
	int result; // eax

	result = a2;
	if (a2) {
		result = nox_xxx_spellGrantToPlayer_4FB550(a2, **(uint32_t**)(a1 + 700), 1, 0, 0);
	}
	return result;
}

//----- (004EBD40) --------------------------------------------------------
void nox_xxx_collideUndeadKiller_4EBD40(int a1, int a2, int a3) {
	int v3;            // ebx
	int* v4;           // esi
	unsigned short v5; // ax
	int v6;            // esi
	int v7;            // edi
	int v8;            // esi
	int v9;            // eax
	int v10;           // eax
	int v11;           // [esp+14h] [ebp+4h]

	if (a2) {
		if (*(uint8_t*)(a2 + 8) & 2 && *(uint8_t*)(a2 + 12) & 0x40) {
			v3 = a1;
			v4 = *(int**)(a1 + 700);
			v5 = nox_xxx_unitGetHP_4EE780(a2);
			v6 = *v4;
			v11 = v6;
			v7 = v5;
			v8 = *(uint32_t*)(v6 + 72);
			if (v8 <= v5) {
				if (v8) {
					v10 = nox_xxx_findParentChainPlayer_4EC580(v3);
					(*(void (**)(int, int, int, int, int))(a2 + 716))(a2, v10, v3, v8, 6);
					nox_xxx_delayedDeleteObject_4E5CC0(v3);
					*(uint32_t*)(v11 + 72) -= v8;
				} else {
					nox_xxx_delayedDeleteObject_4E5CC0(v3);
				}
			} else {
				v9 = nox_xxx_findParentChainPlayer_4EC580(v3);
				(*(void (**)(int, int, int, int, int))(a2 + 716))(a2, v9, v3, v7, 6);
				*(uint32_t*)(v11 + 72) = v8 - v7;
			}
		}
	} else if (!a3) {
		nox_xxx_delayedDeleteObject_4E5CC0(a1);
	}
}

//----- (004EBE10) --------------------------------------------------------
void nox_xxx_collideMonsterGen_4EBE10(int a1, int a2) {
	if (a2) {
		if (*(uint8_t*)(a2 + 8) & 4) {
			nox_xxx_scriptCallByEventBlock_502490((int*)(*(uint32_t*)(a1 + 748) + 72), a2, a1, 19);
		}
	}
}

//----- (004EBE40) --------------------------------------------------------
void sub_4EBE40(int a1, int a2) {
	float2* v2;   // edi
	int v3;       // esi
	int v4;       // eax
	int v5;       // ecx
	int v6;       // esi
	uint32_t* v7; // [esp+Ch] [ebp+4h]

	v2 = (float2*)a1;
	v7 = *(uint32_t**)(a1 + 700);
	if (nox_common_gameFlags_check_40A5C0(4096) && a2 && *(uint8_t*)(a2 + 8) & 4) {
		sub_4D7520(0);
		v3 = 1;
		v4 = nox_xxx_getFirstPlayerUnit_4DA7C0();
		if (!v4) {
			sub_4D71E0(gameFrame());
		} else {
			do {
				v5 = *(uint32_t*)(v4 + 748);
				if (*(uint32_t*)(*(uint32_t*)(v5 + 276) + 4792) == 1 && *(uint32_t*)(v5 + 308)) {
					v3 = 0;
				}
				v4 = nox_xxx_getNextPlayerUnit_4DA7F0(v4);
			} while (v4);
			if (v3 == 1) {
				sub_4D71E0(gameFrame());
			}
		}
		v6 = *(uint32_t*)(a2 + 748);
		if (*(float2**)(v6 + 308) != v2 || (unsigned int)(gameFrame() - *v7) > (int)gameFPS()) {
			nox_xxx_aud_501960(1005, (int)v2, 0, 0);
			nox_xxx_netSendPointFx_522FF0(130, v2 + 7);
			nox_xxx_netPriMsgToPlayer_4DA2C0(a2, "objcoll.c:SoulGateCollide", 0);
		}
		*(uint32_t*)(v6 + 308) = v2;
		*v7 = gameFrame();
	}
}

//----- (004EBF40) --------------------------------------------------------
void nox_xxx_collideAnkhQuest_4EBF40(int a1, int a2) {
	int v2;            // esi
	int v3;            // ebp
	int v4;            // ebx
	int v5;            // ecx
	uint32_t* v6;      // eax
	int v7;            // edi
	int v8;            // eax
	int v9;            // edx
	int v10;           // eax
	uint32_t* v11;     // ecx
	uint32_t* v12;     // eax
	int v13;           // edx
	int v14;           // eax
	uint32_t* v15;     // ecx
	unsigned char v16; // al
	float v17;         // [esp+0h] [ebp-20h]
	int v18;           // [esp+14h] [ebp-Ch]
	int v19;           // [esp+18h] [ebp-8h]

	v2 = a2;
	v3 = *(uint32_t*)(a1 + 692);
	v19 = 0;
	if (a2 && *(uint8_t*)(a2 + 8) & 4) {
		v4 = *(uint32_t*)(a2 + 748);
		v5 = 0;
		v6 = (uint32_t*)(*(uint32_t*)(v4 + 276) + 4796);
		while (*v6 != a1) {
			++v5;
			++v6;
			if (v5 >= 5) {
				goto LABEL_8;
			}
		}
		v19 = 1;
	LABEL_8:
		v18 = 0;
		v7 = v3 + 50;
		do {
			if (gameFrame() - *(uint32_t*)(v7 + 26) > (unsigned int)(240 * gameFPS())) {
				nox_wcscpy((wchar2_t*)(v7 - 50), (const wchar2_t*)getMemAt(0x5D4594, 1568012));
				*(uint8_t*)(v7 + 1) = getMemByte(0x5D4594, 1568016);
				*(uint8_t*)v7 = 0;
				*(uint32_t*)(v7 + 26) = 0;
			}
			v8 = *(uint32_t*)(v4 + 276);
			if (*(uint8_t*)v7 == *(uint8_t*)(v8 + 2251) &&
				!nox_wcscmp((const wchar2_t*)(v7 - 50), (const wchar2_t*)(v8 + 2185))) {
				if (!strcmp((const char*)(v7 + 1), (const char*)(*(uint32_t*)(v4 + 276) + 2112))) {
					v9 = *(uint32_t*)(v4 + 276);
					v10 = 0;
					v11 = (uint32_t*)(v9 + 4796);
					while (*v11) {
						++v10;
						++v11;
						if (v10 >= 5) {
							v2 = a2;
							goto LABEL_17;
						}
					}
					v2 = a2;
					*(uint32_t*)(v9 + 4 * v10 + 4796) = a1;
				LABEL_17:
					if ((unsigned long long)(nox_platform_get_ticks() - *(uint64_t*)&qword_5d4594_1567940) <= 0x5DC) {
						return;
					}
					nox_xxx_netPriMsgToPlayer_4DA2C0(v2, "objcoll.c:ExtraLifeAlreadyAwarded", 0);
					nox_xxx_aud_501960(925, v2, 0, 0);
					*(uint64_t*)&qword_5d4594_1567940 = nox_platform_get_ticks();
					return;
				}
				v2 = a2;
			}
			v7 += 80;
			++v18;
		} while (v18 < 64);
		if (v19 == 1) {
			if ((unsigned long long)(nox_platform_get_ticks() - *(uint64_t*)&qword_5d4594_1567940) <= 0x5DC) {
				return;
			}
			nox_xxx_netPriMsgToPlayer_4DA2C0(v2, "objcoll.c:ExtraLifeAlreadyAwarded", 0);
			nox_xxx_aud_501960(925, v2, 0, 0);
			*(uint64_t*)&qword_5d4594_1567940 = nox_platform_get_ticks();
			return;
		}
		v17 = nox_xxx_gamedataGetFloat_419D40("MaxExtraLives");
		if (*(uint32_t*)(v4 + 320) < nox_float2int(v17)) {
			v12 = nox_xxx_newObjectByTypeID_4E3810("AnkhTradable");
			if (v12) {
				((void (*)(int, uint32_t*, int, uint32_t))v12[177])(v2, v12, 1, 0);
			}
			*(uint32_t*)(a1 + 136) = gameFrame();
			nox_xxx_aud_501960(1004, a1, 0, 0);
			nox_xxx_netSendPointFx_522FF0(130, (float2*)(a1 + 56));
			nox_xxx_netPriMsgToPlayer_4DA2C0(v2, "objcoll.c:AwardExtraLife", 0);
			v13 = *(uint32_t*)(v4 + 276);
			v14 = 0;
			v15 = (uint32_t*)(v13 + 4796);
			while (*v15) {
				++v14;
				++v15;
				if (v14 >= 5) {
					goto LABEL_35;
				}
			}
			*(uint32_t*)(v13 + 4 * v14 + 4796) = a1;
		LABEL_35:
			nox_wcscpy((wchar2_t*)(v3 + 80 * *(unsigned char*)(v3 + 5120)),
					   (const wchar2_t*)(*(uint32_t*)(v4 + 276) + 2185));
			*(uint8_t*)(80 * *(unsigned char*)(v3 + 5120) + v3 + 50) = *(uint8_t*)(*(uint32_t*)(v4 + 276) + 2251);
			strcpy((char*)(80 * *(unsigned char*)(v3 + 5120) + v3 + 51), (const char*)(*(uint32_t*)(v4 + 276) + 2112));
			*(uint32_t*)(80 * *(unsigned char*)(v3 + 5120) + v3 + 76) = gameFrame();
			v16 = *(uint8_t*)(v3 + 5120) + 1;
			*(uint8_t*)(v3 + 5120) = v16;
			if (v16 >= 0x40u) {
				*(uint8_t*)(v3 + 5120) = 0;
			}
		} else if ((unsigned long long)(nox_platform_get_ticks() - *(uint64_t*)&qword_5d4594_1567940) > 0x5DC) {
			nox_xxx_netPriMsgToPlayer_4DA2C0(v2, "pickup.c:MaxTradableAnkhsReached", 0);
			nox_xxx_aud_501960(925, v2, 0, 0);
			*(uint64_t*)&qword_5d4594_1567940 = nox_platform_get_ticks();
			return;
		}
	}
}

//----- (004EC3E0) --------------------------------------------------------
int nox_xxx_playerObserverFindGoodSlave2_4EC3E0(int a1) {
	int result; // eax

	if (!a1) {
		return 0;
	}
	result = *(uint32_t*)(a1 + 516);
	if (!result) {
		return 0;
	}
	while (!(*(uint8_t*)(result + 8) & 2) || !(*(uint8_t*)(*(uint32_t*)(result + 748) + 1440) & 0x80)) {
		result = *(uint32_t*)(result + 512);
		if (!result) {
			return 0;
		}
	}
	return result;
}

//----- (004EC420) --------------------------------------------------------
int nox_xxx_playerObserverFindGoodSlave_4EC420(int a1) {
	int result; // eax

	if (!a1) {
		return 0;
	}
	if (!*(uint32_t*)(a1 + 508)) {
		return 0;
	}
	result = *(uint32_t*)(a1 + 512);
	if (!result) {
		return 0;
	}
	while (!(*(uint8_t*)(result + 8) & 2) || !(*(uint8_t*)(*(uint32_t*)(result + 748) + 1440) & 0x80)) {
		result = *(uint32_t*)(result + 512);
		if (!result) {
			return 0;
		}
	}
	return result;
}

//----- (004EC470) --------------------------------------------------------
void nox_xxx_unitRemoveChild_4EC470(nox_object_t* a1) {
	int v1; // eax
	int v2; // ecx

	if (a1) {
		v1 = *(uint32_t*)((int)a1 + 516);
		if (v1) {
			do {
				v2 = *(uint32_t*)(v1 + 512);
				*(uint32_t*)(v1 + 508) = 0;
				*(uint32_t*)(v1 + 512) = 0;
				v1 = v2;
			} while (v2);
		}
		*(uint32_t*)((int)a1 + 516) = 0;
	}
}

//----- (004EC4B0) --------------------------------------------------------
void nox_xxx_unitTransferSlaves_4EC4B0(nox_object_t* a1p) {
	int a1 = a1p;
	int v1; // eax
	int v2; // esi

	if (a1) {
		v1 = *(uint32_t*)(a1 + 516);
		if (v1) {
			do {
				v2 = *(uint32_t*)(v1 + 512);
				nox_xxx_unitSetOwner_4EC290(*(uint32_t*)(a1 + 508), v1);
				v1 = v2;
			} while (v2);
		}
	}
}

//----- (004EC520) --------------------------------------------------------
int nox_xxx_unitsHaveSameTeam_4EC520(nox_object_t* a1p, nox_object_t* a2p) {
	int a1 = a1p;
	int a2 = a2p;
	int v2; // edi
	int v3; // esi

	v2 = a1;
	if (a1 && a2) {
		while (2) {
			v3 = a2;
			do {
				if (nox_xxx_servCompareTeams_419150(v2 + 48, v3 + 48) || v2 == v3) {
					return 1;
				}
				v3 = *(uint32_t*)(v3 + 508);
			} while (v3);
			v2 = *(uint32_t*)(v2 + 508);
			if (v2) {
				continue;
			}
			break;
		}
	}
	return 0;
}

//----- (004EC5B0) --------------------------------------------------------
void sub_4EC5B0() {
	dword_5d4594_1568024 = 0;
	nox_alloc_class_free_all(*(uint32_t**)&nox_alloc_respawn_1568020);
	nox_xxx_respawnAllow_587000_205200 = 1;
}

//----- (004EC5E0) --------------------------------------------------------
uint32_t* nox_xxx_respawnAdd_4EC5E0(nox_object_t* a1p) {
	int a1 = a1p;
	uint32_t* result;  // eax
	uint32_t* v2;      // ebx
	unsigned short v3; // cx
	uint8_t* v4;       // ebp

	result = *(uint32_t**)&nox_xxx_respawnAllow_587000_205200;
	if (nox_xxx_respawnAllow_587000_205200) {
		result = nox_alloc_class_new_obj_zero(*(uint32_t**)&nox_alloc_respawn_1568020);
		v2 = result;
		if (result) {
			v3 = *(uint16_t*)(a1 + 4);
			result[1] = a1;
			*result = v3;
			result[2] = *(uint32_t*)(a1 + 56);
			result[3] = *(uint32_t*)(a1 + 60);
			*((uint16_t*)result + 8) = *(uint16_t*)(a1 + 124);
			if (*(uint32_t*)(a1 + 8) & 0x13001000) {
				memcpy(result + 7, *(const void**)(a1 + 692), 0x14u);
			}
			if (*(uint32_t*)(a1 + 8) & 0x1000000 && nox_xxx_weaponInventoryEquipFlags_415820(a1) & 0x82) {
				v4 = *(uint8_t**)(a1 + 736);
				*((uint8_t*)v2 + 48) = v4[1];
				*((uint8_t*)v2 + 49) = *v4;
			}
			v2[14] = 0;
			v2[13] = dword_5d4594_1568024;
			result = *(uint32_t**)&dword_5d4594_1568024;
			if (dword_5d4594_1568024) {
				*(uint32_t*)(dword_5d4594_1568024 + 56) = v2;
			}
			dword_5d4594_1568024 = v2;
		}
	}
	return result;
}

//----- (004EC6A0) --------------------------------------------------------
void sub_4EC6A0(int a1) {
	int v1;       // eax
	uint64_t* v2; // ecx
	int v3;       // eax
	int v4;       // ecx
	int v5;       // ecx

	v1 = dword_5d4594_1568024;
	if (*(uint32_t*)(dword_5d4594_1568024 + 4) == a1) {
		v2 = *(uint64_t**)&dword_5d4594_1568024;
		dword_5d4594_1568024 = *(uint32_t*)(dword_5d4594_1568024 + 52);
		v3 = *(uint32_t*)(v1 + 52);
		if (v3) {
			*(uint32_t*)(v3 + 56) = 0;
		}
		nox_alloc_class_free_obj_first(*(unsigned int**)&nox_alloc_respawn_1568020, v2);
	} else if (dword_5d4594_1568024) {
		while (*(uint32_t*)(v1 + 4) != a1) {
			v1 = *(uint32_t*)(v1 + 52);
			if (!v1) {
				return;
			}
		}
		v4 = *(uint32_t*)(v1 + 56);
		if (v4) {
			*(uint32_t*)(v4 + 52) = *(uint32_t*)(v1 + 52);
		}
		v5 = *(uint32_t*)(v1 + 52);
		if (v5) {
			*(uint32_t*)(v5 + 56) = *(uint32_t*)(v1 + 56);
		}
		nox_alloc_class_free_obj_first(*(unsigned int**)&nox_alloc_respawn_1568020, (uint64_t*)v1);
	}
}

//----- (004ECA60) --------------------------------------------------------
int nox_xxx_allocItemRespawnArray_4ECA60() {
	nox_alloc_respawn_1568020 = nox_new_alloc_class("Respawn", 60, 384);
	return nox_alloc_respawn_1568020 != 0;
}

//----- (004ECA90) --------------------------------------------------------
void sub_4ECA90() { nox_free_alloc_class(*(void**)&nox_alloc_respawn_1568020); }

//----- (004ECCB0) --------------------------------------------------------
nox_object_t* nox_server_getObjectFromNetCode_4ECCB0(int a1) {
	int result; // eax
	int v2;     // edi
	int v3;     // esi
	char* v4;   // eax
	int v5;     // ecx

	result = nox_server_netCodeCache_lookupObj_4ECD90(a1);
	if (result) {
		return result;
	}
	v2 = nox_server_getFirstObject_4DA790();
	if (v2) {
		while (1) {
			if (!(*(uint8_t*)(v2 + 16) & 0x20) && *(uint32_t*)(v2 + 36) == a1) {
				nox_server_netCodeCache_addObj_4ECEA0(v2);
				return v2;
			}
			v3 = *(uint32_t*)(v2 + 504);
			if (v3) {
				while (*(uint8_t*)(v3 + 16) & 0x20 || *(uint32_t*)(v3 + 36) != a1) {
					v3 = *(uint32_t*)(v3 + 496);
					if (!v3) {
						goto LABEL_9;
					}
				}
				nox_server_netCodeCache_addObj_4ECEA0(v3);
				return v3;
			}
		LABEL_9:
			v2 = nox_server_getNextObject_4DA7A0(v2);
			if (!v2) {
				goto LABEL_10;
			}
		}
	}
LABEL_10:
	v3 = nox_server_getFirstObjectUninited_4DA870();
	if (v3) {
		while (*(uint8_t*)(v3 + 16) & 0x20 || *(uint32_t*)(v3 + 36) != a1) {
			v3 = nox_server_getNextObjectUninited_4DA880(v3);
			if (!v3) {
				goto LABEL_17;
			}
		}
		nox_server_netCodeCache_addObj_4ECEA0(v3);
		return v3;
	}
LABEL_17:
	v4 = nox_common_playerInfoGetFirst_416EA0();
	if (!v4) {
		return 0;
	}
	while (1) {
		v5 = *((uint32_t*)v4 + 514);
		if (v5) {
			if (!(*(uint8_t*)(v5 + 16) & 0x20) && *(uint32_t*)(v5 + 36) == a1) {
				break;
			}
		}
		v4 = nox_common_playerInfoGetNext_416EE0((int)v4);
		if (!v4) {
			return 0;
		}
	}
	return *((uint32_t*)v4 + 514);
}

nox_server_netCodeCacheStruct nox_server_netCodeCache;

//----- (004ECD90) --------------------------------------------------------
int nox_server_netCodeCache_lookupObj_4ECD90(int a1) {
	uint32_t* v1; // esi

	if (nox_server_needInitNetCodeCache) {
		nox_server_netCodeCache_initArray_4ECE50();
	}
	v1 = *(uint32_t**)&nox_server_netCodeCache.firstUsedObject;
	if (!*(uint32_t*)&nox_server_netCodeCache.firstUsedObject) {
		return 0;
	}
	while (*(uint32_t*)(*v1 + 36) != a1) { // Search for netCode/objectCode/extent in cache
		v1 = (uint32_t*)v1[2];
		if (!v1) {
			return 0;
		}
	}
	sub_4ECE10(&nox_server_netCodeCache.firstUsedObject, (int)v1);
	sub_4ECDE0(&nox_server_netCodeCache.firstUsedObject, (int)v1);
	return *v1;
}

//----- (004ECDE0) --------------------------------------------------------
int sub_4ECDE0(uint32_t* a1, int a2) {
	int result; // eax

	result = a2;
	*(uint32_t*)(a2 + 4) = 0;
	*(uint32_t*)(a2 + 8) = *a1;
	if (*a1) {
		*(uint32_t*)(*a1 + 4) = a2;
	} else {
		a1[1] = a2;
	}
	*a1 = a2;
	return result;
}

//----- (004ECE10) --------------------------------------------------------
int sub_4ECE10(uint32_t* a1, int a2) {
	int result; // eax
	int v3;     // ecx
	int v4;     // ecx

	result = a2;
	v3 = *(uint32_t*)(a2 + 8);
	if (v3) {
		*(uint32_t*)(v3 + 4) = *(uint32_t*)(a2 + 4);
	} else {
		a1[1] = *(uint32_t*)(a2 + 4);
	}
	v4 = *(uint32_t*)(a2 + 4);
	if (v4) {
		*(uint32_t*)(v4 + 8) = *(uint32_t*)(a2 + 8);
	} else {
		result = *(uint32_t*)(a2 + 8);
		*a1 = result;
	}
	return result;
}

//----- (004ECE50) --------------------------------------------------------
int nox_server_netCodeCache_initArray_4ECE50() {
	unsigned char* v0; // esi
	int result;        // eax

	v0 = &nox_server_netCodeCache.objArray[0]; // actually start of the array! size=234-44=(192/12)=16
	nox_server_netCodeCache.firstUsedObject = 0;
	nox_server_netCodeCache.lastUsedObject = 0;
	nox_server_netCodeCache.firstFreeObject = 0;
	nox_server_netCodeCache.lastFreeObject = 0;
	do {
		result = sub_4ECDE0(&nox_server_netCodeCache.firstFreeObject, (int)v0);
		v0 += 12;
	} while ((int)v0 < &nox_server_netCodeCache.firstUsedObject);
	nox_server_needInitNetCodeCache = 0;
	return result;
}

//----- (004ECEA0) --------------------------------------------------------
int nox_server_netCodeCache_addObj_4ECEA0(int a1) {
	int* v1;    // eax
	int v2;     // esi
	int result; // eax
	int v4;     // [esp-8h] [ebp-8h]

	v1 = (int*)nox_server_netCodeCache_nextUnused_4ECEF0();
	if (v1) {
		*v1 = a1;
		result = sub_4ECDE0(&nox_server_netCodeCache.firstUsedObject, (int)v1);
	} else {
		v2 = nox_server_netCodeCache.lastUsedObject;
		v4 = nox_server_netCodeCache.lastUsedObject;
		nox_server_netCodeCache.lastUsedObject->value = a1;
		sub_4ECE10(&nox_server_netCodeCache.firstUsedObject, v4);
		result = sub_4ECDE0(&nox_server_netCodeCache.firstUsedObject, v2);
	}
	return result;
}

//----- (004ECEF0) --------------------------------------------------------
int nox_server_netCodeCache_nextUnused_4ECEF0() {
	int result; // eax

	result = nox_server_netCodeCache.firstFreeObject;
	if (!nox_server_netCodeCache.firstFreeObject) {
		return 0;
	}
	nox_server_netCodeCache.firstFreeObject =
		nox_server_netCodeCache.firstFreeObject->prev; //*(uint32_t*)(*(uint32_t*)&netCodeCache.firstFreeObject + 8);
	return result;
}

//----- (004ECF10) --------------------------------------------------------
int sub_4ECF10(int a1) {
	int result = nox_server_getFirstObject_4DA790();
	while (result) {
		if (!(*(uint8_t*)(result + 16) & 0x20 || *(uint32_t*)(result + 44) != a1)) {
			return result;
		}

		int v2 = *(uint32_t*)(result + 504);

		while (v2) {
			if (!(*(uint8_t*)(v2 + 16) & 0x20 || *(uint32_t*)(v2 + 44) != a1)) {
				return v2;
			}
			v2 = *(uint32_t*)(v2 + 496);
		}
		result = nox_server_getNextObject_4DA7A0(result);
	}

	result = nox_server_getFirstObjectUninited_4DA870();
	while (result != 0) {
		if (!(*(uint8_t*)(result + 16) & 0x20) && *(uint32_t*)(result + 44) == a1) {
			return result;
		}
		result = nox_server_getNextObjectUninited_4DA880(result);
	}

	nox_object_t* obj = nox_xxx_getFirstUpdatable2Object_4DA840();
	while (obj) {
		if (!(obj->obj_flags & 0x20) && obj->script_id == a1) {
			return obj;
		}

		obj = nox_xxx_getNextUpdatable2Object_4DA850(obj);
	}

	return obj;
}

//----- (004ECFA0) --------------------------------------------------------
int sub_4ECFA0(nox_object_t* a1) {
	int result;   // eax
	uint32_t* v2; // esi

	result = nox_server_needInitNetCodeCache;
	if (!nox_server_needInitNetCodeCache) {
		v2 = *(uint32_t**)&nox_server_netCodeCache.firstUsedObject;
		if (*(uint32_t*)&nox_server_netCodeCache.firstUsedObject) {
			result = a1;
			while (*v2 != a1) {
				v2 = (uint32_t*)v2[2];
				if (!v2) {
					return result;
				}
			}
			sub_4ECE10(&nox_server_netCodeCache.firstUsedObject, (int)v2);
			result = sub_4ECDE0(&nox_server_netCodeCache.firstFreeObject, (int)v2);
		}
	}
	return result;
}

//----- (004ECFE0) --------------------------------------------------------
int sub_4ECFE0() {
	int result; // eax
	int v1;     // esi
	int v2;     // edi

	result = nox_server_needInitNetCodeCache;
	if (!nox_server_needInitNetCodeCache) {
		v1 = *(uint32_t*)&nox_server_netCodeCache.firstUsedObject;
		if (*(uint32_t*)&nox_server_netCodeCache.firstUsedObject) {
			do {
				v2 = *(uint32_t*)(v1 + 8);
				sub_4ECE10(&nox_server_netCodeCache.firstUsedObject, v1);
				result = sub_4ECDE0(&nox_server_netCodeCache.firstFreeObject, v1);
				v1 = v2;
			} while (v2);
		}
	}
	return result;
}

//----- (004ED050) --------------------------------------------------------
void sub_4ED050(int a1, int a2) {
	int v2; // eax
	int i;  // esi
	int v4; // edi

	LOWORD(v2) = *getMemU16Ptr(0x5D4594, 1568248);
	if (!*getMemU32Ptr(0x5D4594, 1568248)) {
		v2 = nox_xxx_getNameId_4E3AA0("Crown");
		*getMemU32Ptr(0x5D4594, 1568248) = v2;
	}
	for (i = *(uint32_t*)(a1 + 516); i; i = *(uint32_t*)(i + 512)) {
		LOWORD(v2) = *(uint16_t*)(i + 4);
		if ((unsigned short)v2 == *getMemU32Ptr(0x5D4594, 1568248)) {
			v4 = *(uint32_t*)(i + 748);
			LOWORD(v2) = nox_xxx_dropCrown_4ED5E0(a1, i, (int*)(a1 + 56));
			*(uint32_t*)(v4 + 4) = a2;
		}
	}
}

//----- (004EED40) --------------------------------------------------------
void nox_xxx_abilGivePlayerAll_4EED40(int a1, char a2, int a3) {
	int* v3;      // esi
	uint32_t* v4; // edi
	int v5;       // ebx

	if (a1 && a2 > 0) {
		v3 = getMemIntPtr(0x587000, 206108);
		v4 = (uint32_t*)(*(uint32_t*)(*(uint32_t*)(a1 + 748) + 276) + 3696);
		v5 = a2;
		do {
			if (*v3) {
				if (nox_common_gameFlags_check_40A5C0(4096) || nox_xxx_isQuest_4D6F50() || sub_4D6F70()) {
					*v4 = 0;
				} else {
					nox_xxx_abilityRewardServ_4FB9C0_ability(a1, *v3, a3);
				}
			}
			++v3;
			++v4;
			--v5;
		} while (v5);
	}
}

//----- (004EEDC0) --------------------------------------------------------
int nox_xxx_plrReadVals_4EEDC0(nox_object_t* a1p, int a2) {
    int a1 = a1p;
	int v2;     // ebx
	int v3;     // edi
	int v4;     // esi
	float4 v5a;  // ebp
	float* v5;  // ebp
	short v6;   // ax
	char v7;    // al
	double v8;  // st7
	short v9;   // ax
	double v10; // st7
	int i;      // eax
	int v12;    // edx
	int v13;    // eax
	size_t v14; // eax
	int result; // eax
	float v16;  // [esp-1Ch] [ebp-40h]
	float v17;  // [esp+0h] [ebp-24h]
	float v18;  // [esp+0h] [ebp-24h]
	float v19;  // [esp+0h] [ebp-24h]
	float v20;  // [esp+0h] [ebp-24h]
	int v21;    // [esp+0h] [ebp-24h]
	float v22;  // [esp+14h] [ebp-10h]
	int v23;    // [esp+18h] [ebp-Ch]
	char v24;   // [esp+1Ch] [ebp-8h]
	float4 v25a; // [esp+20h] [ebp-4h]
	float* v25; // [esp+20h] [ebp-4h]
	float4 v26a; // [esp+28h] [ebp+4h]
	float* v26; // [esp+28h] [ebp+4h]

	v2 = a1;
	v3 = *(uint32_t*)(a1 + 748);
	v23 = 0;
	v4 = *(uint32_t*)(v3 + 276);
	v26a = sub_57B350();
	v5a = nox_xxx_plrGetMaxVarsPtr_57B360(*(unsigned char*)(v4 + 2251));
	v25a = nox_xxx_plrGetMaxVarsPtr_57B360(0);
	v26 = &v26a;
	v5 = &v5a;
	v25 = &v25a;
	if (nox_common_gameFlags_check_40A5C0(0x2000)) {
		*(uint16_t*)(*(uint32_t*)(v2 + 556) + 4) = nox_float2int(*v5);
		v6 = nox_float2int16_abs(*v5);
		nox_xxx_unitSetHP_4E4560(v2, v6);
		*(uint16_t*)(v3 + 8) = nox_float2int(v5[1]);
		*(uint16_t*)(v3 + 4) = nox_float2int(v5[1]);
		*(uint32_t*)(v4 + 2239) = nox_float2int(v5[3]);
		v16 = v5[2];
		*(float*)(v2 + 548) = v5[2] * 0.000099999997;
		*(uint32_t*)(v4 + 2235) = nox_float2int(v16);
		if (!*(uint8_t*)(v4 + 2251) && !nox_common_gameFlags_check_40A5C0(4096) && !sub_4D6F30()) {
			nox_xxx_abilGivePlayerAll_4EED40(v2, 10, 0);
		}
	} else {
		v7 = *(uint8_t*)(v4 + 3684);
		if (v7 > NOX_PLAYER_MAX_LEVEL) {
			v7 = NOX_PLAYER_MAX_LEVEL;
		}
		v24 = v7;
		v22 = (double)(v7 - 1);
		v17 = (*v5 - *v26) * v22 / ((double)NOX_PLAYER_MAX_LEVEL - 1) + *v26 + 0.5;
		*(uint16_t*)(*(uint32_t*)(v2 + 556) + 4) = nox_float2int(v17);
		nox_xxx_unitSetHP_4E4560(v2, *(uint16_t*)(*(uint32_t*)(v2 + 556) + 4));
		v8 = (v5[1] - v26[1]) * v22 / ((double)NOX_PLAYER_MAX_LEVEL - 1) + v26[1];
		if (v8 > v5[1]) {
			v8 = v5[1];
		}
		v18 = v8 + 0.5;
		v9 = nox_float2int(v18);
		*(uint16_t*)(v3 + 8) = v9;
		*(uint16_t*)(v3 + 4) = v9;
		v19 = (v5[3] - v26[3]) * v22 / ((double)NOX_PLAYER_MAX_LEVEL - 1) + v26[3] + 0.5;
		*(uint32_t*)(v4 + 2239) = nox_float2int(v19);
		v10 = (v5[2] - v26[2]) * v22 / ((double)NOX_PLAYER_MAX_LEVEL - 1) + v26[2];
		*(float*)(v2 + 548) = v10 * 0.000099999997;
		v20 = v10 + 0.5;
		*(uint32_t*)(v4 + 2235) = nox_float2int(v20);
		if (!*(uint8_t*)(v4 + 2251)) {
			nox_xxx_abilGivePlayerAll_4EED40(v2, v24, a2);
		}
	}
	*(float*)(v2 + 120) = (double)*(int*)(v4 + 2239) / v25[3] * 20.0 + 10.0;
	*(uint16_t*)(*(uint32_t*)(v3 + 276) + 3652) =
		(long long)(((double)*(int*)(v4 + 2239) / v25[3] * 1250.0 + 750.0) * *getMemDoublePtr(0x581450, 10216));
	*(uint16_t*)(v2 + 490) = *(uint16_t*)(*(uint32_t*)(v3 + 276) + 3652);
	sub_56F780(*(uint32_t*)(*(uint32_t*)(v3 + 276) + 4624), *(uint32_t*)(v4 + 2239));
	sub_56F780(*(uint32_t*)(*(uint32_t*)(v3 + 276) + 4620), *(uint32_t*)(v4 + 2235));
	nox_xxx_protectPlayerHPMana_56F870(*(uint32_t*)(*(uint32_t*)(v3 + 276) + 4600), *(uint16_t*)(v3 + 8));
	nox_xxx_protectPlayerHPMana_56F870(*(uint32_t*)(*(uint32_t*)(v3 + 276) + 4592),
									   *(uint16_t*)(*(uint32_t*)(v2 + 556) + 4));
	for (i = *(uint32_t*)(v2 + 504); i; v23 += v12) {
		v12 = *(unsigned char*)(i + 488);
		i = *(uint32_t*)(i + 496);
	}
	*(uint32_t*)(v4 + 3656) = v23 > *(unsigned short*)(v2 + 490);
	v13 = *(uint32_t*)(v3 + 276);
	v21 = *(uint32_t*)(v13 + 4628);
	v14 = nox_wcslen((const wchar2_t*)(v13 + 2185));
	result = sub_56FB00((int*)(*(uint32_t*)(v3 + 276) + 2185), 2 * v14, v21);
	*(uint8_t*)(v4 + 2184) = 1;
	return result;
}

//----- (004EF140) --------------------------------------------------------
int sub_4EF140(int a1) {
	int v1;     // edi
	int v2;     // ecx
	int result; // eax
	int i;      // ebx
	int v5;     // [esp-Ch] [ebp-14h]

	v1 = *(uint32_t*)(*(uint32_t*)(a1 + 748) + 276);
	if (nox_common_gameFlags_check_40A5C0(0x2000)) {
		v2 = *(uint32_t*)(v1 + 4644);
		*(uint8_t*)(v1 + 3684) = NOX_PLAYER_MAX_LEVEL;
		sub_56F820(v2, 0xAu);
		result = nox_xxx_plrReadVals_4EEDC0(a1, 0);
	} else {
		for (i = 0; i <= NOX_PLAYER_MAX_LEVEL; ++i) {
			if (nox_xxx_gamedataGetFloatTable_419D70("XPTable", i) > *(float*)(a1 + 28)) {
				break;
			}
		}
		v5 = *(uint32_t*)(v1 + 4644);
		*(uint8_t*)(v1 + 3684) = i - 1;
		sub_56F820(v5, i - 1);
		result = nox_xxx_plrReadVals_4EEDC0(a1, 0);
	}
	return result;
}

//----- (004EF1E0) --------------------------------------------------------
double nox_xxx_calcBoltDamage_4EF1E0(int a1, int a2) {
	double result; // st7

	if (!*getMemU32Ptr(0x5D4594, 1568264)) {
		*getMemU32Ptr(0x5D4594, 1568264) = nox_xxx_getNameId_4E3AA0("ArcherBolt");
	}
	if (!nox_common_gameFlags_check_40A5C0(2048) || *(uint32_t*)(a2 + 4) != *getMemU32Ptr(0x5D4594, 1568264)) {
		result = (double)(a1 - *(unsigned short*)(a2 + 60)) * *(float*)(a2 + 64) + (double)*(unsigned short*)(a2 + 72);
	} else {
		result = nox_xxx_gamedataGetFloat_419D40("BoltSoloDamageMin") +
				 (double)(a1 - *(unsigned short*)(a2 + 60)) * *(float*)(a2 + 64);
	}
	return result;
}

//----- (004EF410) --------------------------------------------------------
void sub_4EF410(int a1, unsigned char a2) {
	signed char v2; // bl
	int v3;         // esi
	int v4;         // edi
	uint32_t* v5;   // esi
	int v6;         // [esp-10h] [ebp-24h]
	float v7;       // [esp+0h] [ebp-14h]

	v2 = a2;
	v3 = *(uint32_t*)(*(uint32_t*)(a1 + 748) + 276);
	if ((char)a2 > NOX_PLAYER_MAX_LEVEL) {
		v2 = NOX_PLAYER_MAX_LEVEL;
		a2 = NOX_PLAYER_MAX_LEVEL;
	}
	*(float*)(a1 + 28) = nox_xxx_gamedataGetFloatTable_419D70("XPTable", v2);
	v7 = nox_xxx_gamedataGetFloatTable_419D70("XPTable", v2);
	sub_56F8C0(*(uint32_t*)(v3 + 4604), v7);
	sub_4D81A0(a1);
	v6 = *(uint32_t*)(v3 + 4644);
	*(uint8_t*)(v3 + 3684) = v2;
	sub_56F820(v6, a2);
	nox_xxx_plrReadVals_4EEDC0(a1, 0);
	if (nox_common_gameFlags_check_40A5C0(2048) && !*(uint8_t*)(v3 + 2251)) {
		v4 = 1;
		v5 = (uint32_t*)(v3 + 3700);
		do {
			if (*v5) {
				nox_xxx_book_45DBE0((void*)3, v4, v4 - 1);
			}
			++v4;
			++v5;
		} while (v4 < 6);
	}
	if (nox_common_gameFlags_check_40A5C0(2048)) {
		sub_57AF30(a1, 0);
	}
}

//----- (004EF580) --------------------------------------------------------
char nox_xxx_getRespawnWeaponFlags_4EF580() {
	char v0; // bl
	int v1;  // eax
	int v2;  // eax
	int v3;  // eax
	int v4;  // eax
	int v5;  // eax
	int v6;  // eax
	int v7;  // eax
	int v8;  // eax

	v0 = 0;
	v1 = sub_415CD0((char*)0x400);
	if (nox_xxx_getUnitDefDd10_4E3BA0(v1)) {
		v0 = 1;
	}
	v2 = sub_415CD0((char*)4);
	if (nox_xxx_getUnitDefDd10_4E3BA0(v2)) {
		v0 |= 2u;
	}
	v3 = sub_415CD0((char*)1);
	if (nox_xxx_getUnitDefDd10_4E3BA0(v3)) {
		v0 |= 4u;
	}
	v4 = sub_415840((char*)0x8000);
	if (nox_xxx_getUnitDefDd10_4E3BA0(v4)) {
		v0 |= 8u;
	}
	v5 = sub_415CD0((char*)0x4000);
	if (nox_xxx_getUnitDefDd10_4E3BA0(v5)) {
		v0 |= 0x10u;
	}
	v6 = sub_415840((char*)0x100);
	if (nox_xxx_getUnitDefDd10_4E3BA0(v6)) {
		v0 |= 0x20u;
	}
	v7 = sub_415840((char*)0x200);
	if (nox_xxx_getUnitDefDd10_4E3BA0(v7)) {
		v0 |= 0x40u;
	}
	v8 = sub_415CD0((char*)0x1000000);
	if (nox_xxx_getUnitDefDd10_4E3BA0(v8)) {
		v0 |= 0x80u;
	}
	return v0;
}

//----- (004EF6F0) --------------------------------------------------------
int sub_4EF6F0(int a1) {
	int v1; // esi
	int i;  // eax

	if (!*getMemU32Ptr(0x5D4594, 1568268)) {
		*getMemU32Ptr(0x5D4594, 1568268) = nox_xxx_getNameId_4E3AA0("Glyph");
	}
	v1 = 0;
	for (i = nox_xxx_inventoryGetFirst_4E7980(a1); i; i = nox_xxx_inventoryGetNext_4E7990(i)) {
		if (*(unsigned short*)(i + 4) == *getMemU32Ptr(0x5D4594, 1568268)) {
			++v1;
		}
	}
	return v1;
}

//----- (004EF750) --------------------------------------------------------
nox_object_t* nox_xxx_playerRespawnItem_4EF750(nox_object_t* a1p, char* a2, int* a3, int a4, int a5) {
	int a1 = a1p;
	uint32_t* v5;                    // eax
	uint32_t* v6;                    // esi
	void (*v7)(uint32_t*, uint32_t); // eax
	int v8;                          // eax

	v5 = nox_xxx_newObjectByTypeID_4E3810(a2);
	v6 = v5;
	if (v5) {
		v7 = (void (*)(uint32_t*, uint32_t))v5[172];
		if (v7) {
			v7(v6, 0);
		}
		if (a3) {
			nox_xxx_modifSetItemAttrs_4E4990((int)v6, a3);
		}
		nox_xxx_inventoryServPlace_4F36F0(a1, (int)v6, a4, a5);
		v8 = v6[2];
		v6[4] &= 0xFFF7FFFF;
		if (v8 & 0x3001000) {
			*(uint32_t*)(v6[187] + 4) |= 1u;
		}
	}
	return v6;
}

//----- (004EF7D0) --------------------------------------------------------
int sub_4DE4D0(char a1);
char nox_xxx_playerMakeDefItems_4EF7D0(int a1, int a2, int a3) {
	int v3;                // esi
	int v4;                // edi
	int v5;                // eax
	int v6;                // ecx
	int v7;                // eax
	uint32_t* v8;          // ebp
	int v9;                // eax
	int v10;               // eax
	int v11;               // ebp
	int v12;               // ecx
	int v13;               // eax
	int v14;               // eax
	int v15;               // eax
	int v16;               // edi
	unsigned char v18[20]; // [esp+10h] [ebp-14h]
	unsigned char* v19;    // [esp+28h] [ebp+4h]
	uint32_t* v20;         // [esp+2Ch] [ebp+8h]

	v3 = a1;
	v4 = *(uint32_t*)(a1 + 748);
	v19 = (unsigned char*)(*(uint32_t*)(v4 + 276) + 2185);
	if (a2) {
		nox_xxx_removePoison_4EE9D0(v3);
		nox_xxx_unitHPsetOnMax_4EE6F0(v3);
		nox_xxx_playerManaRefresh_4EECF0(v3);
	}
	nox_xxx_playerCancelAbils_4FC180(v3);
	sub_4D7E50(v3);
	*(uint32_t*)(v4 + 312) = 0;
	*(uint32_t*)(v4 + 316) = 0;
	*(uint8_t*)(v3 + 541) = 0;
	*(uint32_t*)(v4 + 84) = 0;
	*(uint16_t*)(v4 + 78) = 0;
	*(uint16_t*)(v4 + 76) = 0;
	v5 = v4 + 12;
	v6 = 32;
	do {
		v5 += 2;
		--v6;
		*(uint16_t*)(v5 - 2) = **(uint16_t**)(v3 + 556);
	} while (v6);
	*(uint32_t*)(v3 + 16) &= 0xFFEB3FE7;
	nox_xxx_playerSetState_4FA020((uint32_t*)v3, 13);
	nox_xxx_unitClearBuffs_4FF580(v3);
	*(uint8_t*)(v4 + 188) = 0;
	*(uint32_t*)(v4 + 216) = 0;
	*(uint32_t*)(v4 + 192) = 0;
	*(uint32_t*)(v4 + 196) = 0;
	*(uint32_t*)(v4 + 200) = 0;
	*(uint32_t*)(v4 + 204) = 0;
	*(uint32_t*)(v4 + 208) = 0;
	*(uint8_t*)(v4 + 212) = 0;
	*(uint32_t*)(v4 + 136) = 0;
	*(uint32_t*)(v4 + 132) = 0;
	*(uint32_t*)(v4 + 268) = 0;
	sub_4F7950(v3);
	*(uint32_t*)(v3 + 520) = 0;
	if (nox_common_gameFlags_check_40A5C0(0x2000)) {
		sub_4DE4D0(*(uint8_t*)(*(uint32_t*)(v4 + 276) + 2064));
	}
	v7 = *(uint32_t*)(v4 + 276);
	if (v7 && !*(uint32_t*)(v7 + 4700)) {
		nox_xxx_netReportTotalHealth_4D85C0(*(unsigned char*)(v7 + 2064), (uint32_t*)v3);
		nox_xxx_netReportTotalMana_4D88C0(*(unsigned char*)(*(uint32_t*)(v4 + 276) + 2064), v3);
		if (a3) {
			LOBYTE(v7) = nox_xxx_netSendPlayerRespawn_4EFC30(v3, 0);
		} else {
			v8 = *(uint32_t**)(v3 + 504);
			if (v8) {
				do {
					v20 = (uint32_t*)v8[124];
					if (sub_53E2D0((int)v8) || (v9 = v8[4], !(v9 & 0x100)) ||
						v8[2] & 0x2000000 && nox_xxx_unitArmorInventoryEquipFlags_415C70((int)v8) & 0x808) {
						nox_xxx_delayedDeleteObject_4E5CC0((int)v8);
					}
					v8 = v20;
				} while (v20);
			}
			nox_xxx_netSendPlayerRespawn_4EFC30(v3, 1);
			v10 = nox_xxx_modifGetIdByName_413290("UserColor1");
			v11 = nox_xxx_modifGetDescById_413330(v10);
			if (nox_common_gameFlags_check_40A5C0(2560) || *(uint8_t*)(*(uint32_t*)(v4 + 276) + 2251)) {
				v12 = **(uint32_t**)(v4 + 276);
				if (!(v12 & 0x400)) {
					*(uint32_t*)v18 = 0;
					*(uint32_t*)&v18[4] = nox_xxx_modifGetDescById_413330(*(uint32_t*)(v11 + 4) + v19[84]);
					*(uint32_t*)&v18[8] = nox_xxx_modifGetDescById_413330(*(uint32_t*)(v11 + 4) + v19[85]);
					*(uint32_t*)&v18[12] = 0;
					nox_xxx_playerRespawnItem_4EF750(v3, "StreetShirt", (int*)v18, 1, 0);
				}
			}
			if (!(**(uint8_t**)(v4 + 276) & 4)) {
				*(uint32_t*)v18 = 0;
				*(uint32_t*)&v18[4] = nox_xxx_modifGetDescById_413330(*(uint32_t*)(v11 + 4) + v19[83]);
				*(uint32_t*)&v18[8] = 0;
				*(uint32_t*)&v18[12] = 0;
				nox_xxx_playerRespawnItem_4EF750(v3, "StreetPants", (int*)v18, 1, 0);
			}
			if (!(**(uint8_t**)(v4 + 276) & 1)) {
				*(uint32_t*)v18 = nox_xxx_modifGetDescById_413330(*(uint32_t*)(v11 + 4) + v19[87]);
				*(uint32_t*)&v18[4] = nox_xxx_modifGetDescById_413330(*(uint32_t*)(v11 + 4) + v19[86]);
				*(uint32_t*)&v18[8] = 0;
				*(uint32_t*)&v18[12] = 0;
				nox_xxx_playerRespawnItem_4EF750(v3, "StreetSneakers", (int*)v18, 1, 0);
			}
			if (nox_common_gameFlags_check_40A5C0(2048)) {
				v13 = nox_xxx_modifGetIdByName_413290("ArmorQuality1");
				*(uint32_t*)v18 = nox_xxx_modifGetDescById_413330(v13);
				*(uint32_t*)&v18[8] = 0;
				*(uint32_t*)&v18[12] = 0;
				if (*(uint8_t*)(*(uint32_t*)(v4 + 276) + 2251)) {
					*(uint32_t*)&v18[4] = 0;
				} else {
					v14 = nox_xxx_modifGetIdByName_413290("Material1");
					*(uint32_t*)&v18[4] = nox_xxx_modifGetDescById_413330(v14);
				}
				LOBYTE(v7) = (unsigned int)nox_xxx_playerRespawnItem_4EF750(
					v3, *(char**)getMemAt(0x587000, 206376 + 4 * *(unsigned char*)(*(uint32_t*)(v4 + 276) + 2251)),
					(int*)v18, 1, 0);
			} else if (nox_common_gameFlags_check_40A5C0(4096) && sub_4CFE00() >= 0) {
				*(uint32_t*)v18 = 0;
				*(uint32_t*)&v18[4] = 0;
				*(uint32_t*)&v18[8] = 0;
				*(uint32_t*)&v18[12] = 0;
				if (*(uint8_t*)(*(uint32_t*)(v4 + 276) + 2251) == 1) {
					v15 = nox_xxx_modifGetIdByName_413290("Replenishment1");
					*(uint32_t*)&v18[8] = nox_xxx_modifGetDescById_413330(v15);
				}
				LOBYTE(v7) = (unsigned int)nox_xxx_playerRespawnItem_4EF750(
					v3, *(char**)getMemAt(0x587000, 206388 + 4 * *(unsigned char*)(*(uint32_t*)(v4 + 276) + 2251)),
					(int*)v18, 1, 0);
			} else {
				LOBYTE(v7) = *(uint8_t*)(*(uint32_t*)(v4 + 276) + 2251);
				if ((uint8_t)v7) {
					if ((uint8_t)v7 == 1) {
						LOBYTE(v7) = (unsigned int)nox_xxx_playerRespawnItem_4EF750(v3, "WizardRobe", 0, 1, 0);
					}
				} else {
					nox_xxx_playerRespawnItem_4EF750(v3, "Longsword", 0, 1, 0);
					LOBYTE(v7) = (unsigned int)nox_xxx_playerRespawnItem_4EF750(v3, "WoodenShield", 0, 1, 0);
				}
			}
		}
		v16 = *(uint32_t*)(v4 + 276);
		if (v16) {
			*(uint32_t*)(v16 + 4700) = 1;
		}
	}
	return v7;
}

//----- (004EFC30) --------------------------------------------------------
int nox_xxx_netSendPlayerRespawn_4EFC30(int a1, char a2) {
	char v3[9]; // [esp+0h] [ebp-Ch]

	v3[0] = -23;
	*(uint32_t*)&v3[3] = gameFrame();
	*(uint16_t*)&v3[1] = *(uint16_t*)(a1 + 36);
	v3[7] = nox_xxx_getRespawnWeaponFlags_4EF580();
	v3[8] = a2;
	return nox_xxx_netSendPacket1_4E5390(255, (int)v3, 9, 0, 0);
}

//----- (004EFE80) --------------------------------------------------------
char nox_xxx_unitInitPlayer_4EFE80(nox_object_t* a1p) {
	int a1 = a1p;
	int v1;   // edi
	int v2;   // eax
	float v4; // [esp+0h] [ebp-Ch]

	v1 = *(uint32_t*)(a1 + 748);
	v2 = nox_object_getGold_4FA6D0(a1);
	nox_xxx_playerSubGold_4FA5D0(a1, v2);
	sub_4EF140(a1);
	nox_xxx_spellAwardAll1_4EFD80(*(uint32_t*)(v1 + 276));
	nox_xxx_spellAwardAll2_4EFC80(*(uint32_t*)(v1 + 276));
	nox_xxx_plrReadVals_4EEDC0(a1, 0);
	nox_xxx_spellAwardAll3_4EFE10(*(uint32_t*)(v1 + 276));
	if (nox_common_gameFlags_check_40A5C0(4096)) {
		v4 = nox_xxx_gamedataGetFloat_419D40("QuestGameStartingExtraLives");
		*(uint32_t*)(v1 + 320) = nox_float2int(v4);
	}
	return nox_xxx_playerMakeDefItems_4EF7D0(a1, 1, 0);
}

//----- (004EFF10) --------------------------------------------------------
int sub_4EFF10(int a1) {
	int v1;            // edi
	unsigned short v2; // ax
	int v3;            // ecx
	unsigned int v4;   // ecx
	int result;        // eax

	v1 = *(uint32_t*)(a1 + 748);
	nox_xxx_spellAwardAll1_4EFD80(*(uint32_t*)(v1 + 276));
	nox_xxx_spellAwardAll2_4EFC80(*(uint32_t*)(v1 + 276));
	*(uint8_t*)(*(uint32_t*)(v1 + 276) + 3684) = 1;
	nox_xxx_playerCancelAbils_4FC180(a1);
	nox_xxx_plrReadVals_4EEDC0(a1, 0);
	nox_xxx_spellAwardAll3_4EFE10(*(uint32_t*)(v1 + 276));
	v2 = *(uint16_t*)(v1 + 8);
	v3 = *(uint32_t*)(v1 + 276);
	*(uint16_t*)(v1 + 4) = v2;
	*(uint16_t*)(v1 + 6) = v2;
	nox_xxx_protectPlayerHPMana_56F870(*(uint32_t*)(v3 + 4596), v2);
	*(uint32_t*)(v1 + 192) = 0;
	*(uint32_t*)(v1 + 196) = 0;
	*(uint32_t*)(v1 + 200) = 0;
	*(uint32_t*)(v1 + 204) = 0;
	*(uint32_t*)(v1 + 208) = 0;
	*(uint8_t*)(v1 + 212) = 0;
	nox_xxx_unitHPsetOnMax_4EE6F0(a1);
	v4 = *(uint32_t*)(a1 + 16) & 0xFFEB3FE7;
	*(uint8_t*)(a1 + 541) = 0;
	*(uint32_t*)(a1 + 16) = v4;
	nox_xxx_playerSetState_4FA020((uint32_t*)a1, 13);
	nox_xxx_unitClearBuffs_4FF580(a1);
	nox_xxx_playerCancelSpells_4FEAE0(a1);
	nox_xxx_removePoison_4EE9D0(a1);
	sub_4F7950(a1);
	nox_xxx_netReportTotalHealth_4D85C0(*(unsigned char*)(*(uint32_t*)(v1 + 276) + 2064), (uint32_t*)a1);
	nox_xxx_netReportTotalMana_4D88C0(*(unsigned char*)(*(uint32_t*)(v1 + 276) + 2064), a1);
	*(uint32_t*)(a1 + 520) = 0;
	result = -559023410;
	*(uint32_t*)(*(uint32_t*)(v1 + 276) + 3664) = -559023410;
	*(uint32_t*)(*(uint32_t*)(v1 + 276) + 3660) = -559023410;
	return result;
}

//----- (004F24E0) --------------------------------------------------------
int sub_4F24E0(int a1) {
	int v1;            // eax
	unsigned char* v2; // ecx

	v1 = *getMemU32Ptr(0x587000, 207108);
	if (!*getMemU32Ptr(0x587000, 207108)) {
		return 0;
	}
	v2 = getMemAt(0x587000, 207108);
	while (v1 != a1 || !*((uint32_t*)v2 + 1)) {
		v1 = *((uint32_t*)v2 + 3);
		v2 += 12;
		if (!v1) {
			return 0;
		}
	}
	if (a1 && a1 != 34 && a1 != 27 && a1 != 9 && a1 != 41) {
		return 1;
	} else {
		return 0;
	}
}

//----- (004F2530) --------------------------------------------------------
int sub_4F2530(int a1) {
	int v1;           // eax
	unsigned char* i; // ecx

	v1 = *getMemU32Ptr(0x587000, 207796);
	if (!*getMemU32Ptr(0x587000, 207796)) {
		return 0;
	}
	for (i = getMemAt(0x587000, 207796); v1 != a1 || !*((uint32_t*)i + 1); i += 12) {
		v1 = *((uint32_t*)i + 3);
		if (!v1) {
			return 0;
		}
	}
	return a1 != 0;
}

//----- (004F2570) --------------------------------------------------------
int sub_4F2570(int a1) { return a1 > 0 && a1 < 6; }

//----- (004F2590) --------------------------------------------------------
int sub_4F2590(int a1) {
	int v1;            // ebx
	int v3;            // ecx
	int v4;            // edx
	unsigned char* v5; // eax
	int v6;            // edi

	if (!*getMemU32Ptr(0x5D4594, 1568328)) {
		*getMemU32Ptr(0x5D4594, 1568328) = nox_xxx_getNameId_4E3AA0("Diamond");
		*getMemU32Ptr(0x5D4594, 1568332) = nox_xxx_getNameId_4E3AA0("Emerald");
		*getMemU32Ptr(0x5D4594, 1568336) = nox_xxx_getNameId_4E3AA0("Ruby");
		*getMemU32Ptr(0x5D4594, 1568340) = nox_xxx_getNameId_4E3AA0("SulphorousFlareWand");
		*getMemU32Ptr(0x5D4594, 1568344) = nox_xxx_getNameId_4E3AA0("StreetSneakers");
		*getMemU32Ptr(0x5D4594, 1568352) = nox_xxx_getNameId_4E3AA0("StreetShirt");
		*getMemU32Ptr(0x5D4594, 1568348) = nox_xxx_getNameId_4E3AA0("StreetPants");
	}
	v1 = *(uint32_t*)(a1 + 8);
	if (v1 & 0x40) {
		return 0;
	}
	if (v1 & 0x10) {
		return (*(uint32_t*)(a1 + 12) & 0x1FF78) != 0;
	}
	if (v1 & 0x100) {
		return sub_4F2700(a1);
	}
	v3 = *(unsigned short*)(a1 + 4);
	if ((unsigned short)v3 != *getMemU32Ptr(0x5D4594, 1568328) && v3 != *getMemU32Ptr(0x5D4594, 1568332) &&
		v3 != *getMemU32Ptr(0x5D4594, 1568336)) {
		v4 = 0;
		if (*getMemU32Ptr(0x587000, 208180)) {
			v5 = getMemAt(0x587000, 208180);
			while (*((uint32_t*)v5 + 1) != v3) {
				v6 = *((uint32_t*)v5 + 5);
				v5 += 20;
				if (!v6) {
					goto LABEL_18;
				}
			}
			v4 = 1;
		}
	LABEL_18:
		if (v3 == *getMemU32Ptr(0x5D4594, 1568340) || v3 == *getMemU32Ptr(0x5D4594, 1568344) ||
			v3 == *getMemU32Ptr(0x5D4594, 1568352) || v3 == *getMemU32Ptr(0x5D4594, 1568348)) {
			return sub_4F2B60(a1);
		}
		if (v4 != 1) {
			return 0;
		}
		if (v1 & 0x1000000) {
			return sub_4F27A0(a1);
		}
		if (v1 & 0x2000000) {
			return sub_4F2B20(a1);
		}
	}
	return 1;
}

//----- (004F2700) --------------------------------------------------------
int sub_4F2700(int a1) {
	int v1;           // ecx
	int v2;           // eax
	unsigned char* i; // ecx
	int v5;           // eax
	int v6;           // ecx
	unsigned char* j; // edx
	unsigned char v8; // al

	v1 = *(uint32_t*)(a1 + 12);
	if (v1 & 1) {
		if (*getMemU32Ptr(0x587000, 207108)) {
			v2 = *getMemU32Ptr(0x587000, 207108);
			for (i = getMemAt(0x587000, 207108); v2 != **(unsigned char**)(a1 + 736) || !*((uint32_t*)i + 1); i += 12) {
				v2 = *((uint32_t*)i + 3);
				if (!v2) {
					return 0;
				}
			}
			return 1;
		}
		return 0;
	}
	if (v1 & 2) {
		v5 = nox_xxx_guide_427010(*(const char**)(a1 + 736));
		v6 = *getMemU32Ptr(0x587000, 207796);
		if (*getMemU32Ptr(0x587000, 207796)) {
			for (j = getMemAt(0x587000, 207796); v6 != v5 || !*((uint32_t*)j + 1); j += 12) {
				v6 = *((uint32_t*)j + 3);
				if (!v6) {
					return 0;
				}
			}
			return 1;
		}
		return 0;
	}
	if (!(v1 & 4)) {
		return 0;
	}
	v8 = **(uint8_t**)(a1 + 736);
	if (!v8 || v8 >= 6u) {
		return 0;
	}
	return 1;
}

//----- (004F27A0) --------------------------------------------------------
int sub_4F27A0(int a1) {
	int result; // eax

	result = sub_4F27E0(a1);
	if (result) {
		result = sub_4F28C0(a1);
		if (result) {
			result = sub_4F2960(a1) != 0;
		}
	}
	return result;
}

//----- (004F27E0) --------------------------------------------------------
int sub_4F27E0(int a1) {
	int v1;            // edi
	int v2;            // eax
	int v3;            // esi
	int v4;            // edx
	unsigned char* i;  // ecx
	int v6;            // ebx
	unsigned char* v8; // ecx
	int v9;            // edx
	unsigned char* j;  // ecx
	int v11;           // ebx

	v1 = **(uint32_t**)(a1 + 692);
	if (!v1) {
		return 1;
	}
	if (*(uint32_t*)(a1 + 8) & 0x1000000) {
		v2 = nox_xxx_weaponInventoryEquipFlags_415820(a1);
	} else {
		v2 = nox_xxx_unitArmorInventoryEquipFlags_415C70(a1);
	}
	v3 = *(uint32_t*)(a1 + 8) & 0x1000000;
	if (v3) {
		v4 = 0;
		if (*getMemU32Ptr(0x587000, 210712)) {
			for (i = getMemAt(0x587000, 210712); *((uint32_t*)i - 1) != v1; i += 24) {
				v6 = *((uint32_t*)i + 6);
				++v4;
				if (!v6) {
					return 0;
				}
			}
			v8 = getMemAt(0x587000, 210704 + 24 * v4);
			goto LABEL_18;
		}
		return 0;
	}
	v9 = 0;
	if (!*getMemU32Ptr(0x587000, 210856)) {
		return 0;
	}
	for (j = getMemAt(0x587000, 210856); *((uint32_t*)j - 1) != v1; j += 24) {
		v11 = *((uint32_t*)j + 6);
		++v9;
		if (!v11) {
			return 0;
		}
	}
	v8 = getMemAt(0x587000, 210848 + 24 * v9);
LABEL_18:
	if (!v8) {
		return 0;
	}
	if (v3) {
		if (!(v2 & *(uint32_t*)(v1 + 28)) || v2 & *((uint32_t*)v8 + 5)) {
			return 0;
		}
	} else if (!(v2 & *(uint32_t*)(v1 + 32)) || v2 & *((uint32_t*)v8 + 4)) {
		return 0;
	}
	return 1;
}

//----- (004F28C0) --------------------------------------------------------
int sub_4F28C0(int a1) {
	int v1;            // esi
	int v2;            // eax
	int v3;            // edx
	unsigned char* i;  // ecx
	int v5;            // ebp
	unsigned char* v7; // ecx

	v1 = *(uint32_t*)(*(uint32_t*)(a1 + 692) + 4);
	if (v1) {
		if (*(uint32_t*)(a1 + 8) & 0x1000000) {
			v2 = nox_xxx_weaponInventoryEquipFlags_415820(a1);
		} else {
			v2 = nox_xxx_unitArmorInventoryEquipFlags_415C70(a1);
		}
		v3 = 0;
		if (!*getMemU32Ptr(0x587000, 211000)) {
			return 0;
		}
		for (i = getMemAt(0x587000, 211000); *((uint32_t*)i - 1) != v1; i += 24) {
			v5 = *((uint32_t*)i + 6);
			++v3;
			if (!v5) {
				return 0;
			}
		}
		v7 = getMemAt(0x587000, 210992 + 24 * v3);
		if (!v7) {
			return 0;
		}
		if (*(uint32_t*)(a1 + 8) & 0x1000000) {
			if (!(v2 & *(uint32_t*)(v1 + 28)) || v2 & *((uint32_t*)v7 + 5)) {
				return 0;
			}
		} else if (!(v2 & *(uint32_t*)(v1 + 32)) || v2 & *((uint32_t*)v7 + 4)) {
			return 0;
		}
	}
	return 1;
}

//----- (004F2960) --------------------------------------------------------
int sub_4F2960(int a1) {
	int v1;             // eax
	int v2;             // eax
	int v3;             // eax
	int v4;             // eax
	int v5;             // eax
	int v6;             // ebx
	int v7;             // ecx
	char v8;            // dl
	unsigned char* v9;  // edi
	int v10;            // esi
	unsigned char* v11; // edx
	int v12;            // esi
	uint8_t* v13;       // ecx
	unsigned char* v14; // edx
	int v16;            // [esp+10h] [ebp-4h]

	v16 = *(uint32_t*)(a1 + 692);
	if (!dword_5d4594_1568308) {
		v1 = nox_xxx_modifGetIdByName_413290("Replenishment1");
		dword_5d4594_1568308 = nox_xxx_modifGetDescById_413330(v1);
		v2 = nox_xxx_modifGetIdByName_413290("Replenishment2");
		*getMemU32Ptr(0x5D4594, 1568312) = nox_xxx_modifGetDescById_413330(v2);
		v3 = nox_xxx_modifGetIdByName_413290("Replenishment3");
		*getMemU32Ptr(0x5D4594, 1568316) = nox_xxx_modifGetDescById_413330(v3);
		v4 = nox_xxx_modifGetIdByName_413290("Replenishment4");
		*getMemU32Ptr(0x5D4594, 1568320) = nox_xxx_modifGetDescById_413330(v4);
	}
	if (*(uint32_t*)(a1 + 8) & 0x1000000) {
		v5 = nox_xxx_weaponInventoryEquipFlags_415820(a1);
	} else {
		v5 = nox_xxx_unitArmorInventoryEquipFlags_415C70(a1);
	}
	v6 = 2;
	while (1) {
		v7 = *(uint32_t*)(v16 + 4 * v6);
		if (!v7) {
			goto LABEL_40;
		}
		v8 = *(uint8_t*)(v7 + 36);
		v9 = 0;
		if (v6 == 2) {
			if (!(v8 & 1)) {
				return 0;
			}
		} else if (!(v8 & 2)) {
			return 0;
		}
		v10 = 0;
		if (!*getMemU32Ptr(0x587000, 209344)) {
			// nop
		} else {
			v11 = getMemAt(0x587000, 209344);
			while (1) {
				if (*((uint32_t*)v11 - 1) == v7) {
					v9 = getMemAt(0x587000, 209336 + 24 * v10);
					break;
				}
				v11 += 24;
				++v10;
				if (!*(uint32_t*)v11) {
					break;
				}
			}
		}
		if (!v9) {
			if (v7 == dword_5d4594_1568308) {
				goto LABEL_A;
			}
			if (v7 != *getMemU32Ptr(0x5D4594, 1568312) && v7 != *getMemU32Ptr(0x5D4594, 1568316) &&
				v7 != *getMemU32Ptr(0x5D4594, 1568320)) {
				return 0;
			}
		}
		if (v7 == dword_5d4594_1568308 || v7 == *getMemU32Ptr(0x5D4594, 1568312) ||
			v7 == *getMemU32Ptr(0x5D4594, 1568316) || v7 == *getMemU32Ptr(0x5D4594, 1568320)) {
			goto LABEL_A;
		}
		if (*(uint32_t*)(a1 + 8) & 0x1000000) {
			if (!(v5 & *(uint32_t*)(v7 + 28)) || v5 & *((uint32_t*)v9 + 5)) {
				return 0;
			}
		} else if (!(v5 & *(uint32_t*)(v7 + 32)) || v5 & *((uint32_t*)v9 + 4)) {
			return 0;
		}
		goto LABEL_40;
	LABEL_A:
		v12 = 0;
		if (!*getMemU32Ptr(0x587000, 208180)) {
			return 0;
		}
		v13 = *(uint8_t**)getMemAt(0x587000, 208180);
		v14 = getMemAt(0x587000, 208180);
		do {
			if (*v13 == 35 && *((uint32_t*)v14 + 1) == *(unsigned short*)(a1 + 4)) {
				v12 = 1;
			}
			v13 = (uint8_t*)*((uint32_t*)v14 + 5);
			v14 += 20;
		} while (v13);
		if (!v12) {
			return 0;
		}
	LABEL_40:
		if (++v6 >= 4) {
			return 1;
		}
	}
}

//----- (004F2B20) --------------------------------------------------------
int sub_4F2B20(int a1) {
	int result; // eax

	result = sub_4F27E0(a1);
	if (result) {
		result = sub_4F28C0(a1);
		if (result) {
			result = sub_4F2960(a1) != 0;
		}
	}
	return result;
}

//----- (004F2B60) --------------------------------------------------------
int sub_4F2B60(int a1) {
	int v1;           // eax
	uint32_t* v2;     // eax
	short v4;         // ax
	const char*** v5; // esi
	int v6;           // edi

	if (!*getMemU32Ptr(0x5D4594, 1568324)) {
		v1 = nox_xxx_modifGetIdByName_413290("Replenishment1");
		*getMemU32Ptr(0x5D4594, 1568324) = nox_xxx_modifGetDescById_413330(v1);
	}
	if (*(uint32_t*)(a1 + 8) & 0x1000000 && nox_xxx_weaponInventoryEquipFlags_415820(a1) & 0x10000) {
		v2 = *(uint32_t**)(a1 + 692);
		if (*v2) {
			return 0;
		}
		if (v2[1]) {
			return 0;
		}
		if (v2[2] != *getMemU32Ptr(0x5D4594, 1568324)) {
			return 0;
		}
		if (v2[3]) {
			return 0;
		}
	}
	if (*(uint32_t*)(a1 + 8) & 0x2000000) {
		v4 = nox_xxx_unitArmorInventoryEquipFlags_415C70(a1);
		v5 = *(const char****)(a1 + 692);
		if (v4 & 0x405) {
			v6 = 0;
			while (!*v5 || !nox_strnicmp(**v5, "UserColo", 8u)) {
				++v6;
				++v5;
				if (v6 >= 4) {
					return 1;
				}
			}
			return 0;
		}
	}
	return 1;
}

//----- (004F2C30) --------------------------------------------------------
int sub_4F2C30(int a1) {
	int v1;   // edi
	float v3; // [esp+0h] [ebp-Ch]

	if (!*getMemU32Ptr(0x5D4594, 1568356)) {
		*getMemU32Ptr(0x5D4594, 1568356) = nox_xxx_getNameId_4E3AA0("RedPotion");
		*getMemU32Ptr(0x5D4594, 1568360) = nox_xxx_getNameId_4E3AA0("BluePotion");
		*getMemU32Ptr(0x5D4594, 1568364) = nox_xxx_getNameId_4E3AA0("CurePoisonPotion");
		*getMemU32Ptr(0x5D4594, 1568368) = nox_xxx_getNameId_4E3AA0("HastePotion");
		*getMemU32Ptr(0x5D4594, 1568372) = nox_xxx_getNameId_4E3AA0("InvisibilityPotion");
		*getMemU32Ptr(0x5D4594, 1568376) = nox_xxx_getNameId_4E3AA0("ShieldPotion");
		*getMemU32Ptr(0x5D4594, 1568380) = nox_xxx_getNameId_4E3AA0("VampirismPotion");
		*getMemU32Ptr(0x5D4594, 1568384) = nox_xxx_getNameId_4E3AA0("FireProtectPotion");
		*getMemU32Ptr(0x5D4594, 1568388) = nox_xxx_getNameId_4E3AA0("ShockProtectPotion");
		*getMemU32Ptr(0x5D4594, 1568392) = nox_xxx_getNameId_4E3AA0("PoisonProtectPotion");
		*getMemU32Ptr(0x5D4594, 1568396) = nox_xxx_getNameId_4E3AA0("InvulnerabilityPotion");
		*getMemU32Ptr(0x5D4594, 1568400) = nox_xxx_getNameId_4E3AA0("InfravisionPotion");
		*getMemU32Ptr(0x5D4594, 1568404) = nox_xxx_getNameId_4E3AA0("InfinitePainWand");
	}
	if (!a1 || !(*(uint8_t*)(a1 + 8) & 4)) {
		return 1;
	}
	if (nox_xxx_inventoryCountObjects_4E7D30(a1, *getMemIntPtr(0x5D4594, 1568356)) > 9 ||
		nox_xxx_inventoryCountObjects_4E7D30(a1, *getMemIntPtr(0x5D4594, 1568360)) > 9 ||
		nox_xxx_inventoryCountObjects_4E7D30(a1, *getMemIntPtr(0x5D4594, 1568364)) > 9 ||
		nox_xxx_inventoryCountObjects_4E7D30(a1, *getMemIntPtr(0x5D4594, 1568368)) > 9 ||
		nox_xxx_inventoryCountObjects_4E7D30(a1, *getMemIntPtr(0x5D4594, 1568372)) > 9 ||
		nox_xxx_inventoryCountObjects_4E7D30(a1, *getMemIntPtr(0x5D4594, 1568376)) > 9 ||
		nox_xxx_inventoryCountObjects_4E7D30(a1, *getMemIntPtr(0x5D4594, 1568380)) > 9 ||
		nox_xxx_inventoryCountObjects_4E7D30(a1, *getMemIntPtr(0x5D4594, 1568384)) > 9 ||
		nox_xxx_inventoryCountObjects_4E7D30(a1, *getMemIntPtr(0x5D4594, 1568388)) > 9 ||
		nox_xxx_inventoryCountObjects_4E7D30(a1, *getMemIntPtr(0x5D4594, 1568392)) > 9 ||
		nox_xxx_inventoryCountObjects_4E7D30(a1, *getMemIntPtr(0x5D4594, 1568396)) > 9 ||
		nox_xxx_inventoryCountObjects_4E7D30(a1, *getMemIntPtr(0x5D4594, 1568400)) > 9) {
		return 0;
	}
	v3 = nox_xxx_gamedataGetFloat_419D40("ForceOfNatureStaffLimit");
	v1 = nox_float2int(v3);
	return nox_xxx_inventoryCountObjects_4E7D30(a1, *getMemIntPtr(0x5D4594, 1568404)) <= v1;
}

//----- (004F2E70) --------------------------------------------------------
int nox_xxx_spell_4F2E70(int a1) {
	int v1;            // eax
	int v2;            // esi
	unsigned char* v3; // ecx

	v1 = *getMemU32Ptr(0x587000, 207108);
	v2 = 0;
	if (*getMemU32Ptr(0x587000, 207108)) {
		v3 = getMemAt(0x587000, 207108);
		while (v1 != a1 || !*((uint32_t*)v3 + 1)) {
			v1 = *((uint32_t*)v3 + 3);
			v3 += 12;
			if (!v1) {
				goto LABEL_8;
			}
		}
		v2 = 1;
	}
LABEL_8:
	if (a1 == 46 || a1 == 47 || a1 == 48 || a1 == 49 || a1 == 122 || a1 == 123 || a1 == 124 || a1 == 125) {
		v2 = 1;
	}
	return a1 >= 75 && a1 <= 114 || v2;
}

//----- (004F2EF0) --------------------------------------------------------
int sub_4F2EF0(int a1) {
	int v1;            // eax
	int v2;            // edi
	unsigned char* v3; // ecx
	int* v4;           // eax
	unsigned char* v5; // edx
	int v6;            // ecx
	int* v7;           // eax
	int v8;            // ecx

	v1 = *getMemU32Ptr(0x587000, 207796);
	v2 = 0;
	if (*getMemU32Ptr(0x587000, 207796)) {
		v3 = getMemAt(0x587000, 207796);
		while (v1 != a1 || !*((uint32_t*)v3 + 1)) {
			v1 = *((uint32_t*)v3 + 3);
			v3 += 12;
			if (!v1) {
				goto LABEL_8;
			}
		}
		v2 = 1;
	}
LABEL_8:
	v4 = *(int**)getMemAt(0x587000, 207032);
	v5 = getMemAt(0x587000, 207032);
	if (*getMemU32Ptr(0x587000, 207032)) {
		do {
			v6 = *v4;
			v7 = v4 + 1;
			if (v6) {
				while (1) {
					v8 = *v7;
					if (*v7 == a1) {
						break;
					}
					++v7;
					if (!v8) {
						goto LABEL_14;
					}
				}
				v2 = 1;
			}
		LABEL_14:
			v4 = (int*)*((uint32_t*)v5 + 1);
			v5 += 4;
		} while (v4);
	}
	return v2 != 0;
}

//----- (004F3E30) --------------------------------------------------------
int nox_xxx_xfer_4F3E30(unsigned short a1, nox_object_t* a2p, int a3) {
	int a2 = a2p;
	int v3;            // ebp
	unsigned short v4; // si
	uint32_t* v5;      // eax
	uint32_t* v6;      // esi
	int v7;            // edx
	int v8;            // eax
	int v10;           // [esp+10h] [ebp-10Ch]
	int v11;           // [esp+14h] [ebp-108h]
	int v12;           // [esp+18h] [ebp-104h]
	char v13[256];     // [esp+1Ch] [ebp-100h]

	v3 = 0;
	if (a3 <= 0) {
		return 1;
	}
	while (1) {
		if (a1 < 0x3Cu) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v10, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v13, (unsigned char)v10);
			v13[(unsigned char)v10] = 0;
			v4 = nox_xxx_getNameId_4E3AA0(v13);
			if (!v4) {
				return 0;
			}
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v11, 2u);
			v4 = nox_xxx_objectTOCgetTT_42C2B0(v11);
			if (!v4) {
				return 0;
			}
		}
		nox_xxx_fileCryptReadCrcMB_426C20(&v12, 4u);
		v5 = nox_xxx_newObjectWithTypeInd_4E3450(v4);
		v6 = v5;
		if (!v5 || !((int (*)(uint32_t*, uint32_t))v5[176])(v5, 0)) {
			break;
		}
		v7 = *(uint32_t*)(a2 + 504);
		v6[125] = 0;
		v6[124] = v7;
		v8 = *(uint32_t*)(a2 + 504);
		if (v8) {
			*(uint32_t*)(v8 + 500) = v6;
		}
		++v3;
		*(uint32_t*)(a2 + 504) = v6;
		v6[123] = a2;
		if (v3 >= a3) {
			return 1;
		}
	}
	return 0;
}
// 4F3E30: using guessed type char var_100[256];

//----- (004F3F50) --------------------------------------------------------
int nox_xxx_servMapLoadPlaceObj_4F3F50(nox_object_t* a1p, int a2, void* a3p) {
	int a1 = a1p;
	int* a3 = a3p;
	int v3;     // eax
	int v4;     // edi
	int result; // eax
	char* v6;   // eax
	int v7;     // eax
	int v8;     // edi

	if (nox_common_gameFlags_check_40A5C0(0x400000) || nox_xxx_getUnitDefDd10_4E3BA0(*(unsigned short*)(a1 + 4))) {
		if (a3) {
			v6 = nox_xxx_mapGetWallSize_426A70();
			*(float*)(a1 + 56) = *(float*)(a1 + 56) - (double)(int)(23 * *(uint32_t*)v6) + (double)*a3 - 11.0;
			*(float*)(a1 + 60) = *(float*)(a1 + 60) - (double)(int)(23 * *((uint32_t*)v6 + 1)) + (double)a3[1] - 11.0;
		}
		if (nox_common_gameFlags_check_40A5C0(0x400000)) {
			nox_xxx_unitAddToList_5048A0(a1);
			result = 1;
		} else if (nox_common_gameFlags_check_40A5C0(0x200000) || sub_4E3AD0(*(unsigned short*)(a1 + 4))) {
			nox_xxx_createAt_4DAA50(a1, a2, *(float*)(a1 + 56), *(float*)(a1 + 60));
			result = 1;
		} else {
			v7 = *(uint32_t*)(a1 + 504);
			if (v7) {
				do {
					v8 = *(uint32_t*)(v7 + 496);
					nox_xxx_objectFreeMem_4E38A0(v7);
					v7 = v8;
				} while (v8);
			}
			*(uint32_t*)(a1 + 504) = 0;
			nox_xxx_objectFreeMem_4E38A0(a1);
			result = 0;
		}
	} else {
		v3 = *(uint32_t*)(a1 + 504);
		if (v3) {
			do {
				v4 = *(uint32_t*)(v3 + 496);
				nox_xxx_objectFreeMem_4E38A0(v3);
				v3 = v4;
			} while (v4);
		}
		nox_xxx_objectFreeMem_4E38A0(a1);
		result = 0;
	}
	return result;
}

//----- (004F4170) --------------------------------------------------------
int nox_xxx_readObjectOldVer_4F4170(int a1, int a2, int a3) {
	uint8_t** v3;       // esi
	uint8_t* v4;        // ebx
	unsigned int v5;    // edx
	uint8_t* v6;        // eax
	int v7;             // eax
	int v8;             // ebx
	int v9;             // ebp
	uint32_t* v10;      // edi
	uint8_t* v11;       // edx
	int result;         // eax
	uint8_t* v13;       // eax
	int* v14;           // ebx
	uint8_t* v15;       // edi
	unsigned short v16; // ax
	uint8_t* j;         // edi
	unsigned int i;     // edi
	int v19;            // [esp+10h] [ebp-18h]
	int v20;            // [esp+14h] [ebp-14h]
	int v21;            // [esp+18h] [ebp-10h]
	int v22;            // [esp+1Ch] [ebp-Ch]
	int v23[2];         // [esp+20h] [ebp-8h]

	v3 = (uint8_t**)a1;
	LOBYTE(v20) = 0;
	if (nox_crypt_IsReadOnly() == 1) {
		*(uint32_t*)(a1 + 136) = 0;
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)v3 + 40, 4u);
	v4 = v3[4];
	v21 = (unsigned int)v3[4] & 0x11408162;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v21, 4u);
	v5 = (unsigned int)v3[4] & 0xEEBF7E9D;
	v3[4] = (uint8_t*)v5;
	v6 = (uint8_t*)v5;
	if ((unsigned char)v4 & 0x40) {
		LOBYTE(v6) = v5 | 0x40;
		v3[4] = v6;
	}
	v7 = v21;
	v3[4] = (uint8_t*)(v21 | (unsigned int)v3[4]);
	if (nox_crypt_IsReadOnly() == 1) {
		if (v7 & 0x1000000) {
			nox_xxx_objectSetOn_4E75B0((int)v3);
		} else {
			nox_xxx_objectSetOff_4E7600((int)v3);
		}
	}
	v8 = a3;
	v9 = a2;
	if (nox_crypt_IsReadOnly()) {
		if (a3 < 40 || a2 < 4) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(v23, 8u);
			v10 = v3 + 14;
			*((float*)v3 + 14) = (double)v23[0];
			*((float*)v3 + 15) = (double)v23[1];
		} else {
			v10 = v3 + 14;
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)v3 + 56, 4u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)v3 + 60, 4u);
		}
		v11 = (uint8_t*)v10[1];
		v3[16] = (uint8_t*)*v10;
		v3[17] = v11;
	} else {
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)v3 + 56, 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)v3 + 60, 4u);
	}
	if (v8 >= 10) {
		if (*v3) {
			LOBYTE(v19) = strlen(*v3);
		} else {
			LOBYTE(v19) = 0;
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v19, 1u);
		if (nox_crypt_IsReadOnly() == 1) {
			if ((uint8_t)v19) {
				result = (int)calloc(1, (unsigned char)v19 + 1);
				*v3 = (uint8_t*)result;
				if (!result) {
					return result;
				}
			}
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread(*v3, (unsigned char)v19);
		if (*v3) {
			(*v3)[(unsigned char)v19] = 0;
		}
	}
	if (v8 >= 20) {
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)v3 + 52, 1u);
	}
	if (v8 >= 30) {
		v13 = v3[126];
		for (LOBYTE(v20) = 0; v13; LOBYTE(v20) = v20 + 1) {
			v13 = (uint8_t*)*((uint32_t*)v13 + 124);
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v20, 1u);
		if (nox_crypt_IsReadOnly() == 1) {
			v3[34] = (uint8_t*)(unsigned char)v20;
		}
	}
	if (v8 >= 40) {
		v14 = (int*)(v3 + 11);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)v3 + 44, 4u);
		if (!v3[11] && nox_crypt_IsReadOnly() == 1 && !nox_common_gameFlags_check_40A5C0(0x200000) &&
			!nox_common_gameFlags_check_40A5C0(0x400000)) {
			*v14 = nox_server_NextObjectScriptID();
		}
		if (v9 >= 2) {
			v15 = v3[129];
			v16 = 0;
			a1 = 0;
			if (v15) {
				do {
					if (!(v15[16] & 0x20) && sub_4E3B80(*((unsigned short*)v15 + 2))) {
						++a1;
					}
					v15 = (uint8_t*)*((uint32_t*)v15 + 128);
				} while (v15);
				v16 = a1;
			}
			if (v9 < 5) {
				v22 = v16;
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v22, 4u);
				a1 = v22;
			} else {
				nox_xxx_fileReadWrite_426AC0_file3_fread(&a1, 2u);
			}
			if (nox_crypt_IsReadOnly()) {
				for (i = 0; i < (unsigned short)a1; ++i) {
					nox_xxx_fileReadWrite_426AC0_file3_fread(&a3, 4u);
					if (!nox_common_gameFlags_check_40A5C0(0x200000) && !nox_common_gameFlags_check_40A5C0(0x400000)) {
						sub_516F90(*v14, a3);
					}
				}
			} else {
				for (j = v3[129]; j; j = (uint8_t*)*((uint32_t*)j + 128)) {
					if (!(j[16] & 0x20) && sub_4E3B80(*((unsigned short*)j + 2))) {
						nox_xxx_fileReadWrite_426AC0_file3_fread(j + 44, 4u);
					}
				}
			}
		}
		if (v9 >= 3) {
			v21 = (unsigned int)v3[5] & 0x5E;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v21, 4u);
			nox_xxx_unitUnsetXStatus_4E4780((int)v3, 94);
			nox_xxx_unitSetXStatus_4E4800((int)v3, (int*)v21);
		}
	}
	return 1;
}

//----- (004F4530) --------------------------------------------------------
int nox_xxx_mapReadWriteObjData_4F4530(nox_object_t* a1p, int a2) {
	int* a1 = a1p;
	int* v2;          // esi
	int v3;           // edi
	int v4;           // ecx
	short v5;         // ax
	int result;       // eax
	int* v7;          // ebp
	int* v8;          // edi
	int v9;           // ecx
	int v10;          // ebx
	unsigned int v11; // edx
	unsigned int v12; // eax
	int v13;          // eax
	int v14;          // eax
	int v15;          // edi
	int k;            // edi
	unsigned int j;   // edi
	int v18;          // [esp+10h] [ebp-1Ch]
	int v19;          // [esp+14h] [ebp-18h]
	int i;            // [esp+18h] [ebp-14h]
	int v21;          // [esp+1Ch] [ebp-10h]
	int v22[2];       // [esp+20h] [ebp-Ch]
	int v23;          // [esp+28h] [ebp-4h]

	v2 = a1;
	v3 = a2;
	v4 = a1[34];
	v5 = 0;
	LOBYTE(v19) = 0;
	v18 = 0;
	v23 = v4;
	if (a2 >= 40 || !nox_crypt_IsReadOnly()) {
		v18 = 64;
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v18, 2u);
		v5 = v18;
		if ((short)v18 > 64) {
			return 0;
		}
	}
	if (v3 < 40 || v5 < 61) {
		return nox_xxx_readObjectOldVer_4F4170((int)v2, v5, v3);
	}
	if (nox_crypt_IsReadOnly() == 1) {
		v2[34] = 0;
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)v2 + 40, 4u);
	v7 = v2 + 11;
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)v2 + 44, 4u);
	if (!v2[11] && nox_crypt_IsReadOnly() == 1 && !nox_common_gameFlags_check_40A5C0(0x200000) &&
		!nox_common_gameFlags_check_40A5C0(0x400000)) {
		*v7 = nox_server_NextObjectScriptID();
	}
	if (nox_crypt_IsReadOnly()) {
		if ((short)v18 < 4) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(v22, 8u);
			v8 = v2 + 14;
			*((float*)v2 + 14) = (double)v22[0];
			*((float*)v2 + 15) = (double)v22[1];
		} else {
			v8 = v2 + 14;
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)v2 + 56, 4u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)v2 + 60, 4u);
		}
		v9 = v8[1];
		v2[16] = *v8;
		v2[17] = v9;
	} else {
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)v2 + 56, 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)v2 + 60, 4u);
	}
	LOBYTE(a1) = sub_4F40A0((int)v2);
	nox_xxx_fileReadWrite_426AC0_file3_fread(&a1, 1u);
	if (!(uint8_t)a1) {
		return 1;
	}
	v10 = v2[4];
	v21 = v2[4] & 0x11408162;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v21, 4u);
	v11 = v2[4] & 0xEEBF7E9D;
	v2[4] = v11;
	v12 = v11;
	if (v10 & 0x40) {
		LOBYTE(v12) = v11 | 0x40;
		v2[4] = v12;
	}
	v13 = v21;
	v2[4] |= v21;
	if (nox_crypt_IsReadOnly() == 1) {
		if (v13 & 0x1000000) {
			nox_xxx_objectSetOn_4E75B0((int)v2);
		} else {
			nox_xxx_objectSetOff_4E7600((int)v2);
		}
	}
	if (*v2) {
		LOBYTE(a2) = strlen((const char*)*v2);
	} else {
		LOBYTE(a2) = 0;
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread(&a2, 1u);
	if (nox_crypt_IsReadOnly() != 1 || !(uint8_t)a2 ||
		(result = (int)calloc(1, (unsigned char)a2 + 1), (*v2 = result) != 0)) {
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)*v2, (unsigned char)a2);
		if (*v2) {
			*(uint8_t*)((unsigned char)a2 + *v2) = 0;
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)v2 + 52, 1u);
		v14 = v2[126];
		for (LOBYTE(v19) = 0; v14; LOBYTE(v19) = v19 + 1) {
			v14 = *(uint32_t*)(v14 + 496);
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v19, 1u);
		if (nox_crypt_IsReadOnly() == 1) {
			v2[34] = (unsigned char)v19;
		}
		v15 = v2[129];
		for (i = 0; v15; v15 = *(uint32_t*)(v15 + 512)) {
			if (!(*(uint8_t*)(v15 + 16) & 0x20) && sub_4E3B80(*(unsigned short*)(v15 + 4))) {
				++i;
			}
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread(&i, 2u);
		if (nox_crypt_IsReadOnly()) {
			for (j = 0; j < (unsigned short)i; ++j) {
				nox_xxx_fileReadWrite_426AC0_file3_fread(v22, 4u);
				if (!nox_common_gameFlags_check_40A5C0(0x200000) && !nox_common_gameFlags_check_40A5C0(0x400000)) {
					sub_516F90(*v7, v22[0]);
				}
			}
		} else {
			for (k = v2[129]; k; k = *(uint32_t*)(k + 512)) {
				if (!(*(uint8_t*)(k + 16) & 0x20) && sub_4E3B80(*(unsigned short*)(k + 4))) {
					nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(k + 44), 4u);
				}
			}
		}
		v21 = v2[5] & 0x5E;
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v21, 4u);
		nox_xxx_unitUnsetXStatus_4E4780((int)v2, 94);
		nox_xxx_unitSetXStatus_4E4800((int)v2, (int*)v21);
		if ((short)v18 < 63 || (result = nox_xxx_xferReadScriptHandler_4F5580((int)(v2 + 191), (char*)v2[189])) != 0) {
			if ((short)v18 >= 64) {
				v22[0] = v23 - gameFrame();
				nox_xxx_fileReadWrite_426AC0_file3_fread(v22, 4u);
				if (v22[0] > 0 && nox_crypt_IsReadOnly() == 1) {
					if (v2[4] & 0x400000) {
						v2[32] = v22[0];
					}
				}
			}
			return 1;
		}
	}
	return result;
}

//----- (004F4A20) --------------------------------------------------------
int nox_xxx_XFerSpellPagePedistal_4F4A20(int a1) {
	int v1;     // esi
	int v2;     // edi
	int result; // eax

	v1 = a1;
	v2 = *(uint32_t*)(a1 + 136);
	a1 = 60;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&a1, 2u);
	if ((short)a1 > 60) {
		return 0;
	}
	result = nox_xxx_mapReadWriteObjData_4F4530((int*)v1, (short)a1);
	if (result) {
		nox_xxx_fileReadWrite_426AC0_file3_fread(*(uint8_t**)(v1 + 700), 4u);
		if (!*(uint32_t*)(v1 + 136) || nox_crypt_IsReadOnly() != 1 ||
			(result = nox_xxx_xfer_4F3E30(a1, v1, *(uint32_t*)(v1 + 136))) != 0) {
			*(uint32_t*)(v1 + 136) = v2;
			result = 1;
		}
	}
	return result;
}

//----- (004F4AB0) --------------------------------------------------------
int nox_xxx_XFerReadable_4F4AB0(int a1) {
	int* v1;    // esi
	int v2;     // ebx
	int v3;     // edi
	int result; // eax
	size_t v5;  // [esp+Ch] [ebp-4h]

	v1 = (int*)a1;
	v2 = *(uint32_t*)(a1 + 736);
	v3 = *(uint32_t*)(a1 + 136);
	v5 = strlen(*(const char**)(a1 + 736)) + 1;
	a1 = 60;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&a1, 2u);
	if ((short)a1 > 60) {
		return 0;
	}
	result = nox_xxx_mapReadWriteObjData_4F4530(v1, (short)a1);
	if (result) {
		if ((short)a1 >= 2) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v5, 4u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)v2, v5);
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)v2, v5);
		}
		if (nox_crypt_IsReadOnly() != 1 || (*(uint32_t*)(v2 + 256) = 0, !v1[34]) ||
			(result = nox_xxx_xfer_4F3E30(a1, (int)v1, v1[34])) != 0) {
			v1[34] = v3;
			result = 1;
		}
	}
	return result;
}

//----- (004F4B90) --------------------------------------------------------
int nox_xxx_XFerExit_4F4B90(int a1) {
	int* v1;     // ebp
	uint8_t* v2; // ebx
	int v3;      // edi
	int result;  // eax
	uint8_t* i;  // esi
	size_t v6;   // [esp+Ch] [ebp-4h]

	v1 = (int*)a1;
	v2 = *(uint8_t**)(a1 + 700);
	v3 = *(uint32_t*)(a1 + 136);
	v6 = strlen(*(const char**)(a1 + 700)) + 1;
	a1 = 60;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&a1, 2u);
	if ((short)a1 > 60) {
		return 0;
	}
	result = nox_xxx_mapReadWriteObjData_4F4530(v1, (short)a1);
	if (result) {
		if ((short)a1 >= 2) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v6, 4u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v2, v6);
		} else if (nox_crypt_IsReadOnly() == 1) {
			for (i = v2;; ++i) {
				nox_xxx_fileReadWrite_426AC0_file3_fread(i, 1u);
				if (!*i) {
					break;
				}
			}
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread(v2, v6);
		}
		if ((short)a1 >= 31) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(v2 + 80, 4u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v2 + 84, 4u);
		}
		if (!v1[34] || nox_crypt_IsReadOnly() != 1 || (result = nox_xxx_xfer_4F3E30(a1, (int)v1, v1[34])) != 0) {
			v1[34] = v3;
			result = 1;
		}
	}
	return result;
}

//----- (004F4CB0) --------------------------------------------------------
int nox_xxx_XFerDoor_4F4CB0(int a1) {
	int v1;     // edi
	int v2;     // esi
	int result; // eax
	int v4;     // ebx
	int v5;     // ebp
	int v6;     // [esp+8h] [ebp-14h]
	int v7;     // [esp+Ch] [ebp-10h]
	int v8;     // [esp+10h] [ebp-Ch]
	int v9;     // [esp+14h] [ebp-8h]
	int v10;    // [esp+18h] [ebp-4h]

	v1 = a1;
	v2 = *(uint32_t*)(a1 + 748);
	v10 = *(uint32_t*)(a1 + 136);
	v7 = 60;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v7, 2u);
	if ((short)v7 > 60) {
		return 0;
	}
	result = nox_xxx_mapReadWriteObjData_4F4530((int*)v1, (short)v7);
	if (result) {
		if (!nox_crypt_IsReadOnly()) {
			a1 = *(uint32_t*)(v2 + 12);
			v8 = *(unsigned char*)(v2 + 1);
			v6 = *(uint32_t*)(v2 + 4);
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread(&a1, 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v8, 4u);
		if ((short)v7 < 41) {
			v6 = a1;
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v6, 4u);
		}
		if (nox_crypt_IsReadOnly() == 1) {
			*(uint32_t*)(v2 + 12) = a1;
			*(uint16_t*)(v2 + 40) = (a1 << 8) / 32;
			*(uint32_t*)(v2 + 4) = v6;
			*(uint32_t*)(v2 + 8) = a1;
			v9 = *getMemIntPtr(0x587000, 196184 + 8 * v6) / 2;
			v4 = (long long)(((double)v9 + *(float*)(v1 + 56)) * 0.043478262);
			v9 = *getMemIntPtr(0x587000, 196188 + 8 * v6) / 2;
			v5 = (long long)(((double)v9 + *(float*)(v1 + 60)) * 0.043478262);
			nox_xxx_doorAttachWall_410360(v1, v4, v5);
			*(uint32_t*)(v2 + 16) = v4;
			*(uint32_t*)(v2 + 20) = v5;
			*(uint8_t*)(v2 + 1) = v8;
		}
		if (!*(uint32_t*)(v1 + 136) || nox_crypt_IsReadOnly() != 1 ||
			(result = nox_xxx_xfer_4F3E30(v7, v1, *(uint32_t*)(v1 + 136))) != 0) {
			*(uint32_t*)(v1 + 136) = v10;
			result = 1;
		}
	}
	return result;
}

//----- (004F4E50) --------------------------------------------------------
int nox_xxx_unitTriggerXfer_4F4E50(float a1) {
	int v1;      // edi
	uint8_t* v2; // esi
	int result;  // eax
	double v4;   // st7
	double v5;   // st7
	double v6;   // st7
	char* v7;    // ebp
	char* v8;    // eax
	char* v9;    // eax
	char* v10;   // eax
	int v11;     // [esp+Ch] [ebp-Ch]
	int v12;     // [esp+10h] [ebp-8h]
	int v13;     // [esp+14h] [ebp-4h]

	v1 = LODWORD(a1);
	v2 = *(uint8_t**)(LODWORD(a1) + 748);
	v13 = *(uint32_t*)(LODWORD(a1) + 136);
	v11 = 61;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v11, 2u);
	if ((short)v11 > 61) {
		return 0;
	}
	result = nox_xxx_mapReadWriteObjData_4F4530((int*)v1, (short)v11);
	if (result) {
		if (nox_crypt_IsReadOnly()) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&a1, 4u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v12, 4u);
			v5 = (double)SLODWORD(a1);
			a1 = v5;
			*(float*)(v1 + 184) = v5;
			v6 = (double)v12;
			*(float*)(v1 + 188) = v6;
			if (a1 > 60.0) {
				*(uint32_t*)(v1 + 184) = 1114636288;
			}
			if (v6 > 60.0) {
				*(uint32_t*)(v1 + 188) = 1114636288;
			}
		} else {
			v4 = *(float*)(v1 + 188);
			LODWORD(a1) = (long long)*(float*)(v1 + 184);
			v12 = (long long)v4;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&a1, 4u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v12, 4u);
		}
		nox_shape_box_calc((nox_shape*)(v1 + 172));
		if ((short)v11 < 41) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&a1, 3u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(&a1, 3u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(&a1, 3u);
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread(v2 + 54, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v2 + 55, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v2 + 56, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v2 + 57, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v2 + 58, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v2 + 59, 1u);
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread(v2, 4u);
		if ((short)v11 >= 3) {
			v7 = *(char**)(v1 + 756);
			if (v7) {
				v8 = v7 + 256;
			} else {
				v8 = 0;
			}
			nox_xxx_xferReadScriptHandler_4F5580((int)(v2 + 20), v8);
			if (v7) {
				v9 = v7 + 384;
			} else {
				v9 = 0;
			}
			nox_xxx_xferReadScriptHandler_4F5580((int)(v2 + 28), v9);
			if ((short)v11 >= 31) {
				if (v7) {
					v10 = v7 + 512;
				} else {
					v10 = 0;
				}
				nox_xxx_xferReadScriptHandler_4F5580((int)(v2 + 12), v10);
			}
		} else {
			sub_4F5540((int)(v2 + 20));
			sub_4F5540((int)(v2 + 28));
		}
		if (nox_crypt_IsReadOnly() == 1 && (short)v11 < 31) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&a1, 1u);
			nox_xxx_cryptSeekCur_40E0A0(4 * LOBYTE(a1));
			nox_xxx_fileReadWrite_426AC0_file3_fread(&a1, 1u);
			nox_xxx_cryptSeekCur_40E0A0(4 * LOBYTE(a1));
			nox_xxx_fileReadWrite_426AC0_file3_fread(&a1, 1u);
			nox_xxx_cryptSeekCur_40E0A0(4 * LOBYTE(a1));
			nox_xxx_fileReadWrite_426AC0_file3_fread(&a1, 1u);
			nox_xxx_cryptSeekCur_40E0A0(4 * LOBYTE(a1));
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread(v2 + 44, 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(v2 + 48, 4u);
		if (nox_crypt_IsReadOnly() == 1) {
			v2[52] = 0;
			v2[53] = 0;
		}
		if (!nox_crypt_IsReadOnly() || (short)v11 >= 21) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(v2 + 52, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v2 + 53, 1u);
		}
		if ((short)v11 >= 61) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(v2 + 8, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v2 + 9, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 132), 4u);
			if (nox_crypt_IsReadOnly() == 1) {
				nox_xxx_servMarkObjAnimFrame_4E4880(v1, *(uint32_t*)(v1 + 132));
			}
		}
		if (!*(uint32_t*)(v1 + 136) || nox_crypt_IsReadOnly() != 1 ||
			(result = nox_xxx_xfer_4F3E30(v11, v1, *(uint32_t*)(v1 + 136))) != 0) {
			result = 1;
			*(uint32_t*)(v1 + 136) = v13;
		}
	}
	return result;
}

//----- (004F51D0) --------------------------------------------------------
int nox_xxx_XFerHole_4F51D0(int a1) {
	int* v1;    // edi
	int v2;     // ebx
	int v3;     // esi
	int result; // eax
	char* v5;   // eax
	int v6;     // [esp+Ch] [ebp-4h]

	v1 = (int*)a1;
	v2 = *(uint32_t*)(a1 + 756);
	v3 = *(uint32_t*)(a1 + 700);
	v6 = *(uint32_t*)(a1 + 136);
	a1 = 60;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&a1, 2u);
	if ((short)a1 > 60) {
		return 0;
	}
	result = nox_xxx_mapReadWriteObjData_4F4530(v1, (short)a1);
	if (result) {
		if ((short)a1 < 42) {
			*(uint32_t*)(v3 + 24) = 0;
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v3 + 24), 4u);
		}
		if ((short)a1 < 41) {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v3 + 8), 8u);
			*(uint32_t*)(v3 + 4) = -1;
			*(uint32_t*)v3 = 0;
			*(uint32_t*)(v3 + 16) = 0;
			*(uint16_t*)(v3 + 20) = 0;
		} else {
			if (v2) {
				v5 = (char*)(v2 + 128);
			} else {
				v5 = 0;
			}
			nox_xxx_xferReadScriptHandler_4F5580(v3, v5);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v3 + 8), 8u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v3 + 16), 4u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v3 + 20), 2u);
		}
		if (!v1[34] || nox_crypt_IsReadOnly() != 1 || (result = nox_xxx_xfer_4F3E30(a1, (int)v1, v1[34])) != 0) {
			result = 1;
			v1[34] = v6;
		}
	}
	return result;
}

//----- (004F5300) --------------------------------------------------------
int nox_xxx_XFerTransporter_4F5300(int a1) {
	int* v1;    // esi
	int v2;     // edi
	int v3;     // ebx
	int result; // eax
	int v5;     // [esp+Ch] [ebp-4h]

	v1 = (int*)a1;
	v2 = *(uint32_t*)(a1 + 748);
	v3 = *(uint32_t*)(a1 + 136);
	a1 = 60;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&a1, 2u);
	if ((short)a1 > 60) {
		return 0;
	}
	result = nox_xxx_mapReadWriteObjData_4F4530(v1, (short)a1);
	if (result) {
		if (nox_crypt_IsReadOnly()) {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v2 + 16), 4u);
		} else {
			if (*(uint32_t*)(v2 + 12)) {
				v5 = *(uint32_t*)(v2 + 16);
			} else {
				v5 = 0;
			}
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v5, 4u);
		}
		if (!v1[34] || nox_crypt_IsReadOnly() != 1 || (result = nox_xxx_xfer_4F3E30(a1, (int)v1, v1[34])) != 0) {
			v1[34] = v3;
			result = 1;
		}
	}
	return result;
}

//----- (004F53D0) --------------------------------------------------------
int nox_xxx_XFerElevator_4F53D0(int a1) {
	int* v1;     // esi
	uint8_t* v2; // edi
	int v3;      // ebx
	int result;  // eax

	v1 = (int*)a1;
	v2 = *(uint8_t**)(a1 + 748);
	v3 = *(uint32_t*)(a1 + 136);
	a1 = 61;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&a1, 2u);
	if ((short)a1 > 61) {
		return 0;
	}
	result = nox_xxx_mapReadWriteObjData_4F4530(v1, (short)a1);
	if (result) {
		nox_xxx_fileReadWrite_426AC0_file3_fread(v2 + 8, 4u);
		if ((short)a1 >= 41) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(v2 + 16, 4u);
		}
		if ((short)a1 >= 61) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(v2 + 12, 1u);
		}
		if (!v1[34] || nox_crypt_IsReadOnly() != 1 || (result = nox_xxx_xfer_4F3E30(a1, (int)v1, v1[34])) != 0) {
			v1[34] = v3;
			result = 1;
		}
	}
	return result;
}

//----- (004F54A0) --------------------------------------------------------
int nox_xxx_XFerElevatorShaft_4F54A0(int a1) {
	int* v1;    // esi
	int v2;     // edi
	int v3;     // ebx
	int result; // eax

	v1 = (int*)a1;
	v2 = *(uint32_t*)(a1 + 748);
	v3 = *(uint32_t*)(a1 + 136);
	a1 = 60;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&a1, 2u);
	if ((short)a1 > 60) {
		return 0;
	}
	result = nox_xxx_mapReadWriteObjData_4F4530(v1, (short)a1);
	if (result) {
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v2 + 8), 4u);
		if (!v1[34] || nox_crypt_IsReadOnly() != 1 || (result = nox_xxx_xfer_4F3E30(a1, (int)v1, v1[34])) != 0) {
			v1[34] = v3;
			result = 1;
		}
	}
	return result;
}

//----- (004F5540) --------------------------------------------------------
int sub_4F5540(int a1) {
	int result; // eax
	FILE* v2;   // eax

	result = nox_crypt_IsReadOnly();
	if (nox_crypt_IsReadOnly() == 1) {
		v2 = nox_xxx_mapgenGetSomeFile_426A60();
		nox_xxx_mapgenMakeScript_502790(v2, (char*)a1);
		result = nox_common_gameFlags_check_40A5C0(0x400000);
		if (!result) {
			*(uint32_t*)(a1 + 4) = -1;
		}
	}
	return result;
}

//----- (004F5730) --------------------------------------------------------
int nox_xxx_XFerMover_4F5730(int a1) {
	int v1;     // edi
	int v2;     // esi
	int v3;     // ebp
	int result; // eax
	int* v5;    // eax
	int* v6;    // esi
	int v7;     // [esp+Ch] [ebp-4h]

	v1 = a1;
	v2 = *(uint32_t*)(a1 + 748);
	v3 = *(uint32_t*)(a1 + 136);
	v7 = 60;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v7, 2u);
	if ((short)v7 > 60) {
		return 0;
	}
	result = nox_xxx_mapReadWriteObjData_4F4530((int*)v1, (short)v7);
	if (result) {
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v2 + 4), 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v2 + 8), 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v2 + 32), 4u);
		if ((short)v7 >= 41) {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)v2, 1u);
			if (nox_crypt_IsReadOnly()) {
				nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v2 + 16), 4u);
				nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v2 + 24), 4u);
			} else {
				v5 = *(int**)(v2 + 12);
				a1 = 0;
				if (v5) {
					a1 = *v5;
				}
				nox_xxx_fileReadWrite_426AC0_file3_fread(&a1, 4u);
				v6 = *(int**)(v2 + 20);
				a1 = 0;
				if (v6) {
					a1 = *v6;
				}
				nox_xxx_fileReadWrite_426AC0_file3_fread(&a1, 4u);
			}
		}
		if ((short)v7 >= 42) {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 548), 4u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 544), 4u);
		}
		if (!*(uint32_t*)(v1 + 136) || nox_crypt_IsReadOnly() != 1 ||
			(result = nox_xxx_xfer_4F3E30(v7, v1, *(uint32_t*)(v1 + 136))) != 0) {
			*(uint32_t*)(v1 + 136) = v3;
			result = 1;
		}
	}
	return result;
}

//----- (004F5890) --------------------------------------------------------
int nox_xxx_XFerGlyph_4F5890(int a1) {
	int v1;        // edi
	int v2;        // ebp
	int v3;        // eax
	int result;    // eax
	uint8_t* v5;   // esi
	int v6;        // edi
	int* v7;       // ebp
	int v8;        // ebx
	char* v9;      // eax
	size_t v10;    // [esp-Ch] [ebp-128h]
	int v11;       // [esp+8h] [ebp-114h]
	int v12;       // [esp+Ch] [ebp-110h]
	int* v13;      // [esp+10h] [ebp-10Ch]
	int v14;       // [esp+14h] [ebp-108h]
	int v15;       // [esp+18h] [ebp-104h]
	char v16[256]; // [esp+1Ch] [ebp-100h]

	v1 = a1;
	v2 = *(uint32_t*)(a1 + 692);
	v3 = *(uint32_t*)(a1 + 136);
	v13 = *(int**)(a1 + 692);
	v14 = v3;
	v12 = 60;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v12, 2u);
	if ((short)v12 > 60) {
		return 0;
	}
	result = nox_xxx_mapReadWriteObjData_4F4530((int*)a1, (short)v12);
	if (!result) {
		return 0;
	}
	if ((short)v12 < 41) {
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v15, 4u);
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(a1 + 124), 1u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v2 + 28), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v2 + 32), 4u);
	v5 = (uint8_t*)(v2 + 20);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v2 + 20), 1u);
	if (nox_crypt_IsReadOnly() == 1) {
		if ((short)v12 < 31) {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)v2, 0x14u);
			goto LABEL_16;
		}
		v6 = 0;
		if (!*v5) {
			v1 = a1;
			v2 = (int)v13;
			goto LABEL_19;
		}
		v7 = v13;
		do {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v11, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v16, (unsigned char)v11);
			v16[(unsigned char)v11] = 0;
			*v7 = nox_xxx_spellNameToN_4243F0(v16);
			++v6;
			++v7;
		} while (v6 < (unsigned char)*v5);
	} else {
		v8 = 0;
		if (!*v5) {
			goto LABEL_16;
		}
		do {
			LOBYTE(v11) = strlen(nox_xxx_spellNameByN_424870(*(uint32_t*)v2));
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v11, 1u);
			v10 = (unsigned char)v11;
			v9 = nox_xxx_spellNameByN_424870(*(uint32_t*)v2);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v9, v10);
			++v8;
			v2 += 4;
		} while (v8 < (unsigned char)*v5);
	}
	v1 = a1;
	v2 = (int)v13;
LABEL_16:
	if (nox_crypt_IsReadOnly() != 1) {
		goto LABEL_20;
	}
LABEL_19:
	*(uint16_t*)(v1 + 126) = *(uint16_t*)(v1 + 124);
	*(uint32_t*)(v2 + 24) = 0;
LABEL_20:
	if (!*(uint32_t*)(v1 + 136) || nox_crypt_IsReadOnly() != 1 ||
		(result = nox_xxx_xfer_4F3E30(v12, v1, *(uint32_t*)(v1 + 136))) != 0) {
		result = 1;
		*(uint32_t*)(v1 + 136) = v14;
	}
	return result;
}
// 4F5890: using guessed type char var_100[256];

//----- (004F5AA0) --------------------------------------------------------
int nox_xxx_XFerInvLight_4F5AA0(int* a1) {
	int result;   // eax
	uint32_t* v2; // eax
	int v3;       // [esp+8h] [ebp-98h]
	int v4;       // [esp+Ch] [ebp-94h]
	int v5;       // [esp+10h] [ebp-90h]
	char v6[140]; // [esp+14h] [ebp-8Ch]

	v5 = a1[34];
	memset(v6, 0, sizeof(v6));
	v3 = 60;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v3, 2u);
	if ((short)v3 > 60) {
		return 0;
	}
	result = nox_xxx_mapReadWriteObjData_4F4530(a1, (short)v3);
	if (!result) {
		return 0;
	}
	if (nox_crypt_IsReadOnly()) {
		goto LABEL_14;
	}
	if (nox_common_gameFlags_check_40A5C0(6291456)) {
		v2 = (uint32_t*)sub_45A060();
		if (!v2) {
			goto LABEL_14;
		}
		while (v2[32] != a1[10]) {
			v2 = (uint32_t*)nox_drawable_next_45A070((int)v2);
			if (!v2) {
				goto LABEL_14;
			}
		}
	} else if (a1[2] & 0x20400000) {
		v2 = nox_xxx_netSpriteByCodeStatic_45A720(a1[10]);
	} else {
		v2 = nox_xxx_netSpriteByCodeDynamic_45A6F0(a1[9]);
	}
	if (!v2) {
		abort();
	}
	memcpy(v6, v2 + 34, sizeof(v6));
LABEL_14:
	if ((short)v3 >= 2) {
		nox_xxx_fileReadWrite_426AC0_file3_fread(v6, 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v6[4], 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v6[8], 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v6[12], 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v6[16], 0xCu);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v6[28], 2u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v6[30], 2u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v6[32], 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v6[40], 2u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v6[42], 0x30u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v6[90], 0x10u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v6[106], 0x10u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v6[122], 2u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v6[124], 2u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v6[126], 2u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v6[128], 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v6[134], 2u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v6[136], 2u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v6[138], 1u);
		if ((short)v3 > 40) {
			if ((short)v3 >= 42) {
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v6[36], 4u);
			} else {
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v4, 1u);
				*(uint32_t*)&v6[36] = (unsigned char)v4;
			}
			if (nox_crypt_IsReadOnly() == 1) {
				goto LABEL_20;
			}
			goto LABEL_22;
		}
		if (nox_crypt_IsReadOnly() != 1) {
			goto LABEL_22;
		}
		*(uint32_t*)&v6[36] = 0;
	} else {
		nox_xxx_fileReadWrite_426AC0_file3_fread(v6, 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v6[4], 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v6[8], 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v6[12], 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v6[16], 0xCu);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v6[28], 2u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v6[30], 2u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v6[32], 4u);
		*(uint16_t*)&v6[40] = 0;
		*(uint16_t*)&v6[122] = 0;
		*(uint16_t*)&v6[124] = 0;
		*(uint16_t*)&v6[126] = 0;
		*(uint32_t*)&v6[128] = 0;
		*(uint16_t*)&v6[134] = 0;
		v6[138] = -128;
		if (nox_crypt_IsReadOnly() != 1) {
			goto LABEL_22;
		}
		if (*(float*)&v6[4] > 63.0 ||
			(double)*(int*)&v6[12] * *getMemDoublePtr(0x581450, 9752) > *getMemDoublePtr(0x581450, 9744)) {
			sub_484CE0((int)v6, 63.0);
			if (nox_crypt_IsReadOnly() == 1) {
				goto LABEL_20;
			}
			goto LABEL_22;
		}
	}
LABEL_20:
	if (nox_common_gameFlags_check_40A5C0(6291456)) {
		memcpy((void*)(a1[189] + 2432), v6, 0x8Cu);
	}
LABEL_22:
	if (!a1[34] || nox_crypt_IsReadOnly() != 1 ||
		(result = nox_xxx_xfer_4F3E30(v3, (int)a1, a1[34])) != 0) {
		a1[34] = v5;
		result = 1;
	}
	return result;
}

//----- (004F5E50) --------------------------------------------------------
int nox_xxx_XFerSentry_4F5E50(int a1) {
	int* v1;     // edi
	uint8_t* v2; // esi
	int v3;      // ebp
	int result;  // eax

	v1 = (int*)a1;
	v2 = *(uint8_t**)(a1 + 748);
	v3 = *(uint32_t*)(a1 + 136);
	a1 = 61;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&a1, 2u);
	if ((short)a1 > 61) {
		return 0;
	}
	result = nox_xxx_mapReadWriteObjData_4F4530(v1, (short)a1);
	if (result) {
		nox_xxx_fileReadWrite_426AC0_file3_fread(v2 + 4, 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(v2 + 8, 4u);
		if (nox_crypt_IsReadOnly() == 1 || nox_common_gameFlags_check_40A5C0(0x200000)) {
			*(uint32_t*)v2 = *((uint32_t*)v2 + 1);
		}
		if ((short)a1 >= 61) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(v2, 4u);
		}
		if (!v1[34] || nox_crypt_IsReadOnly() != 1 || (result = nox_xxx_xfer_4F3E30(a1, (int)v1, v1[34])) != 0) {
			v1[34] = v3;
			result = 1;
		}
	}
	return result;
}
