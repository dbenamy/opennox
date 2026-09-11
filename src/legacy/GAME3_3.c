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

//----- (004E44F0) --------------------------------------------------------
void nox_xxx_unitNeedSync_4E44F0(nox_object_t* a1) {
	*(uint32_t*)((int)a1 + 152) = -1;
}

//----- (004E4500) --------------------------------------------------------
int* sub_4E4500(nox_object_t* a1p, int a2, int a3, int a4) {
	int a1 = a1p;
	int v4;      // ecx
	int* result; // eax
	int v6;      // edx

	v4 = 0;
	result = (int*)(a1 + 560);
	do {
		if (a4) {
			*result |= a3;
		} else {
			*result &= ~a3;
		}
		v6 = *result;
		if ((1 << v4) & *(uint32_t*)(a1 + 148)) {
			*result = a2 | v6;
		} else if (!(v6 & a3)) {
			*result = v6 & ~a2;
		}
		++v4;
		++result;
	} while (v4 < 32);
	return result;
}

//----- (004E4670) --------------------------------------------------------
int* nox_xxx_unitSetOnOff_4E4670(int a1, int a2) {
	int v2;          // eax
	unsigned int v3; // eax
	int* result;     // eax
	int v5;          // ecx
	int v6;          // edx
	int v7;          // eax

	nox_xxx_unitNeedSync_4E44F0(a1);
	v2 = *(uint32_t*)(a1 + 16);
	if (a2) {
		v3 = v2 | 0x1000000;
	} else {
		v3 = v2 & 0xFEFFFFFF;
	}
	*(uint32_t*)(a1 + 16) = v3;
	if (*(uint32_t*)(a1 + 8) & 0x20400004) {
		result = (int*)(a1 + 560);
		v5 = 32;
		do {
			v6 = *result;
			++result;
			--v5;
			*(result - 1) = v6 & 0xFFFFF000 | 0x40000;
		} while (v5);
	} else {
		v7 = sub_4E4C90(a1, 4u);
		result = sub_4E4500(a1, 0x40000, 4, v7);
	}
	return result;
}

//----- (004E46F0) --------------------------------------------------------
void nox_xxx_unitRaise_4E46F0(nox_object_t* obj, float a2) {
	int a1 = obj;
	int* v2; // eax
	int v3;  // ecx
	int v4;  // edx
	int v5;  // eax

	if (*(float*)(a1 + 104) != a2) {
		nox_xxx_unitNeedSync_4E44F0(a1);
		*(float*)(a1 + 104) = a2;
		if (*(uint32_t*)(a1 + 8) & 0x20400004) {
			v2 = (int*)(a1 + 560);
			v3 = 32;
			do {
				v4 = *v2;
				++v2;
				--v3;
				*(v2 - 1) = v4 & 0xFFFFF000 | 0x400000;
			} while (v3);
		} else {
			v5 = sub_4E4C90(a1, 0x40u);
			sub_4E4500(a1, 0x400000, 64, v5);
		}
	}
}

//----- (004E4880) --------------------------------------------------------
int* nox_xxx_servMarkObjAnimFrame_4E4880(int a1, int a2) {
	int* result; // eax
	int v3;      // ecx
	int v4;      // edx
	int v5;      // eax

	nox_xxx_unitNeedSync_4E44F0(a1);
	*(uint32_t*)(a1 + 132) = a2;
	if (*(uint32_t*)(a1 + 8) & 0x20400004) {
		result = (int*)(a1 + 560);
		v3 = 32;
		do {
			v4 = *result;
			++result;
			--v3;
			*(result - 1) = v4 & 0xFFFFF000 | 0x10000;
		} while (v3);
	} else {
		v5 = sub_4E4C90(a1, 1u);
		result = sub_4E4500(a1, 0x10000, 1, v5);
	}
	return result;
}

//----- (004E48F0) --------------------------------------------------------
int* nox_xxx_setUnitBuffFlags_4E48F0(int a1, int a2) {
	char v2;     // cl
	int* result; // eax
	int v4;      // ecx
	int v5;      // edx
	int v6;      // eax

	nox_xxx_unitNeedSync_4E44F0(a1);
	v2 = *(uint8_t*)(a1 + 8);
	*(uint32_t*)(a1 + 340) = a2;
	if (v2 & 4) {
		nox_xxx_playerResetProtectionCRC_56F7D0(*(uint32_t*)(*(uint32_t*)(*(uint32_t*)(a1 + 748) + 276) + 4612), a2);
	}
	if (*(uint32_t*)(a1 + 8) & 0x20400004) {
		result = (int*)(a1 + 560);
		v4 = 32;
		do {
			v5 = *result;
			++result;
			--v4;
			*(result - 1) = 0x800000 | v5 & 0xFFFFF000;
		} while (v4);
	} else {
		v6 = sub_4E4C90(a1, 0x80u);
		result = sub_4E4500(a1, 0x800000, 128, v6);
	}
	return result;
}

//----- (004E4990) --------------------------------------------------------
int* nox_xxx_modifSetItemAttrs_4E4990(nox_object_t* a1p, int* a2) {
	int a1 = a1p;
	int v2;      // eax
	int v3;      // edx
	int* result; // eax
	int v5;      // ecx
	int v6;      // ecx
	int v7;      // edx
	int v8;      // eax

	v2 = *(uint32_t*)(a1 + 8);
	if (v2 & 0x1000 && *(uint32_t*)(a1 + 12) & 0x47F0000) {
		goto LABEL_19;
	}
	v3 = 0;
	result = a2;
	v5 = 4;
	do {
		if (*result) {
			v3 = 1;
		}
		++result;
		--v5;
	} while (v5);
	if (!v3) {
		return result;
	}
LABEL_19:
	result = *(int**)getMemAt(0x5D4594, 1564960);
	if (!*getMemU32Ptr(0x5D4594, 1564960)) {
		result = (int*)nox_xxx_getNameId_4E3AA0("TeamBase");
		*getMemU32Ptr(0x5D4594, 1564960) = result;
	}
	if (*(uint32_t*)(a1 + 8) & 0x13001000 || (int*)*(unsigned short*)(a1 + 4) == result) {
		nox_xxx_unitNeedSync_4E44F0(a1);
		memcpy(*(void**)(a1 + 692), a2, 0x14u);
		if (*(uint32_t*)(a1 + 8) & 0x20400004) {
			result = (int*)(a1 + 560);
			v6 = 32;
			do {
				v7 = *result;
				++result;
				--v6;
				*(result - 1) = v7 & 0xFFFFF000 | 0x2000000;
			} while (v6);
		} else {
			v8 = sub_4E4C90(a1, 0x200u);
			result = sub_4E4500(a1, 0x2000000, 512, v8);
		}
	}
	return result;
}

//----- (004E4A70) --------------------------------------------------------
double nox_xxx_objectGetMass_4E4A70(int a1) { return *(float*)(a1 + 120); }

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

//----- (004E5AD0) --------------------------------------------------------
void nox_xxx_playerRemoveSpawnedStuff_4E5AD0(nox_object_t* a1p) {
	int a1 = a1p;
	char v1; // al
	int v2;  // esi
	int v3;  // edi

	if (*(uint8_t*)(a1 + 8) & 4) {
		v1 = *(uint8_t*)(*(uint32_t*)(*(uint32_t*)(a1 + 748) + 276) + 2251);
		if (v1 == 1) {
			sub_4E5F40(a1);
		} else if (v1 == 2) {
			sub_4E5FC0(a1);
		}
	}
	v2 = *(uint32_t*)(a1 + 516);
	if (v2) {
		do {
			v3 = *(uint32_t*)(v2 + 512);
			if (*(uint8_t*)(v2 + 8) & 1 || !sub_4E3B80(*(unsigned short*)(v2 + 4))) {
				nox_xxx_delayedDeleteObject_4E5CC0(v2);
			}
			v2 = v3;
		} while (v3);
	}
}

//----- (004E5B50) --------------------------------------------------------
int nox_xxx_isUnit_4E5B50(nox_object_t* a1p) {
	int a1 = a1p;
	int v1;     // ecx
	int result; // eax

	result = 0;
	if (!nox_common_gameFlags_check_40A5C0(0x2000)) {
		if (*(uint8_t*)(a1 + 8) & 2) {
			v1 = *(uint32_t*)(a1 + 12);
			if (v1 & 0x100) {
				result = 1;
			}
		}
	}
	return result;
}

//----- (004E5B80) --------------------------------------------------------
int sub_4E5B80(nox_object_t* a1p) {
	int a1 = a1p;
	int v1;     // eax
	int result; // eax

	if (!*getMemU32Ptr(0x5D4594, 1565592)) {
		*getMemU32Ptr(0x5D4594, 1565592) = nox_xxx_getNameId_4E3AA0("Pixie");
	}
	result = 0;
	if (a1) {
		if (*(uint8_t*)(a1 + 8) & 1) {
			if (nox_common_gameFlags_check_40A5C0(2048)) {
				if (*(unsigned short*)(a1 + 4) == *getMemU32Ptr(0x5D4594, 1565592)) {
					v1 = nox_xxx_findParentChainPlayer_4EC580(a1);
					if (v1) {
						if (*(uint8_t*)(v1 + 8) & 4) {
							result = 1;
						}
					}
				}
			}
		}
	}
	return result;
}

//----- (004E5BF0) --------------------------------------------------------
void sub_4E5BF0(int a1) {
	int v1; // esi
	int v2; // edi
	int v3; // eax
	int v4; // eax

	if (!*getMemU32Ptr(0x5D4594, 1565596)) {
		*getMemU32Ptr(0x5D4594, 1565596) = nox_xxx_getNameId_4E3AA0("Moonglow");
	}
	v1 = nox_server_getFirstObject_4DA790();
	if (v1) {
		do {
			v2 = nox_server_getNextObject_4DA7A0(v1);
			if (!a1 || !(*(uint8_t*)(v1 + 8) & 4) &&
						   ((v3 = *(uint32_t*)(v1 + 492)) == 0 || !(*(uint8_t*)(v3 + 8) & 4)) &&
						   (*(unsigned short*)(v1 + 4) != *getMemU32Ptr(0x5D4594, 1565596) ||
							(v4 = *(uint32_t*)(v1 + 508)) == 0 || !(*(uint8_t*)(v4 + 8) & 4)) &&
						   !nox_xxx_isUnit_4E5B50(v1)) {
				nox_xxx_delayedDeleteObject_4E5CC0(v1);
			}
			v1 = v2;
		} while (v2);
	}
	nox_object_t* obj = nox_xxx_getFirstUpdatable2Object_4DA840();
	if (obj) {
		do {
			nox_object_t* v6 = nox_xxx_getNextUpdatable2Object_4DA850(obj);
			if (a1 != 1 || !sub_4E5B80(obj)) {
				nox_xxx_delayedDeleteObject_4E5CC0(obj);
			}
			obj = v6;
		} while (obj);
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

//----- (004E6BD0) --------------------------------------------------------
int sub_4E6BD0(int a1) {
	return *(uint32_t*)(a1 + 556) && (unsigned int)(gameFrame() - *(uint32_t*)(a1 + 536)) <= 1;
}

//----- (004E6C00) --------------------------------------------------------
double nox_xxx_calcDistance_4E6C00(nox_object_t* a1p, nox_object_t* a2p) {
	int a1 = a1p;
	int a2 = a2p;
	double v2;          // st7
	double v3;          // st6
	int v4;             // edx
	double v5;     // st5
	double result; // st7
	double v7;          // st6
	int v8;             // eax
	double v9;          // st6
	float v10;          // [esp+4h] [ebp+4h]
	float v11;          // [esp+4h] [ebp+4h]

	v2 = *(float*)(a1 + 56) - *(float*)(a2 + 56);
	v3 = *(float*)(a1 + 60) - *(float*)(a2 + 60);
	v4 = *(uint32_t*)(a1 + 172);
	v5 = sqrt(v3 * v3 + v2 * v2);
	result = v5;
	if (v4 == 2) {
		result = v5 - *(float*)(a1 + 176);
	} else if (v4 == 3) {
		v7 = *(float*)(a1 + 184) * 0.5;
		v10 = *(float*)(a1 + 188) * 0.5;
		if (v7 <= v10) {
			v7 = v10;
		}
		result = v5 - v7;
	}
	v8 = *(uint32_t*)(a2 + 172);
	if (v8 == 2) {
		result = result - *(float*)(a2 + 176);
	} else if (v8 == 3) {
		v9 = *(float*)(a2 + 184) * 0.5;
		v11 = *(float*)(a2 + 188) * 0.5;
		if (v9 <= v11) {
			v9 = v11;
		}
		result = result - v9;
	}
	if (result < 0.0099999998) {
		result = 0.0099999998;
	}
	return result;
}

//----- (004E6CE0) --------------------------------------------------------
int sub_4E6CE0(float2* a1, float2* a2) {
	double v2;  // st7
	int v3;     // ecx
	int v4;     // eax
	int v5;     // ecx
	int v6;     // eax
	int result; // eax

	*(float*)&dword_5d4594_1565628 = a2->field_0 - a1->field_0;
	*(float*)&dword_5d4594_1565632 = a2->field_4 - a1->field_4;
	*getMemFloatPtr(0x5D4594, 1565652) = *(float*)&dword_5d4594_1565628 * 0.41304299 - *(float*)&dword_5d4594_1565632;
	v2 = *(float*)&dword_5d4594_1565628 * 2.4210529 - *(float*)&dword_5d4594_1565632;
	*getMemFloatPtr(0x5D4594, 1565656) = v2;
	*getMemFloatPtr(0x5D4594, 1565636) = *(float*)&dword_5d4594_1565628 * -2.4210529 - *(float*)&dword_5d4594_1565632;
	*getMemFloatPtr(0x5D4594, 1567708) = *(float*)&dword_5d4594_1565628 * -0.41304299 - *(float*)&dword_5d4594_1565632;
	if (*getMemFloatPtr(0x5D4594, 1565652) < 0.0) {
		v3 = 0;
	} else {
		v3 = 8;
	}
	if (v2 < 0.0) {
		v4 = 0;
	} else {
		v4 = 4;
	}
	v5 = v4 | v3;
	if (*getMemFloatPtr(0x5D4594, 1565636) < 0.0) {
		v6 = 0;
	} else {
		v6 = 2;
	}
	*getMemU32Ptr(0x5D4594, 1565640) = (*getMemFloatPtr(0x5D4594, 1567708) >= 0.0) | v6 | v5;
	switch (*getMemU32Ptr(0x5D4594, 1565640)) {
	case 0:
		result = 2;
		break;
	case 2:
		result = 6;
		break;
	case 3:
		result = 4;
		break;
	case 4:
		result = 10;
		break;
	case 0xB:
		result = 5;
		break;
	case 0xC:
		result = 8;
		break;
	case 0xD:
		result = 9;
		break;
	case 0xF:
		result = 1;
		break;
	default:
		result = 0;
		break;
	}
	return result;
}

//----- (004E6E50) --------------------------------------------------------
int nox_server_testTwoPointsAndDirection_4E6E50(float2* a1, int a2, float2* a3) {
	int v3;    // esi
	int v5[2]; // [esp+4h] [ebp-8h]

	nox_xxx_xferIndexedDirection_509E20(a2, (int2*)v5);
	v3 = v5[1] + v5[0] + 2 * v5[1] + 4;
	return *getMemU32Ptr(0x587000, 202504 + 4 * (sub_4E6CE0(a1, a3) + 16 * v3));
}

//----- (004E7190) --------------------------------------------------------
void nox_xxx_teleportToMB_4E7190(uint8_t* a1, float* a2) {
	if (!nox_xxx_testUnitBuffs_4FF350((int)a1, 14) && !(a1[16] & 2) &&
		(!nox_common_gameFlags_check_40A5C0(4096) || !(a1[8] & 2) || !(a1[12] & 8)) &&
		(nox_common_gameFlags_check_40A5C0(2048) || a1[8] & 6)) {
		nox_xxx_unitMove_4E7010((int)a1, (float2*)a2);
	}
}

//----- (004E7290) --------------------------------------------------------
int nox_xxx_objectUnkUpdateCoords_4E7290(nox_object_t* a1p) {
	int a1 = a1p;
	int result; // eax
	int v2;     // ecx
	int v3;     // edx
	int v4;     // ecx

	result = a1;
	switch (*(uint32_t*)(a1 + 172)) {
	case 1:
		v2 = *(uint32_t*)(a1 + 56);
		*(uint32_t*)(a1 + 232) = v2;
		v3 = v2;
		v4 = *(uint32_t*)(a1 + 60);
		*(uint32_t*)(a1 + 240) = v3;
		*(uint32_t*)(a1 + 236) = v4;
		*(uint32_t*)(a1 + 244) = v4;
		break;
	case 2:
		*(float*)(a1 + 232) = *(float*)(a1 + 56) - *(float*)(a1 + 176);
		*(float*)(a1 + 236) = *(float*)(a1 + 60) - *(float*)(a1 + 176);
		*(float*)(a1 + 240) = *(float*)(a1 + 176) + *(float*)(a1 + 56);
		*(float*)(a1 + 244) = *(float*)(a1 + 176) + *(float*)(a1 + 60);
		break;
	case 3:
		*(float*)(a1 + 232) = *(float*)(a1 + 200) + *(float*)(a1 + 56);
		*(float*)(a1 + 236) = *(float*)(a1 + 196) + *(float*)(a1 + 60);
		*(float*)(a1 + 240) = *(float*)(a1 + 208) + *(float*)(a1 + 56);
		*(float*)(a1 + 244) = *(float*)(a1 + 220) + *(float*)(a1 + 60);
		break;
	}
	return result;
}

//----- (004E7470) --------------------------------------------------------
void nox_xxx_spawnSomeBarrel_4E7470(int a1, int a2) {
	const char* v2;    // eax
	unsigned char* v3; // edi
	int v4;            // eax
	int v5;            // edx
	unsigned char* v6; // esi
	unsigned char* v7; // ecx
	int v8;            // ebx
	char* result;      // eax
	int v10;           // edi
	int i;             // ebx
	uint32_t* v12;     // esi
	float2 a3;         // [esp+Ch] [ebp-8h]

	v2 = (const char*)nox_xxx_getUnitName_4E39D0(a1);
	v3 = getMemAt(0x587000, 203080);
	if (strncmp(v2, "Barrel", 6u)) {
		v3 = getMemAt(0x587000, 203240);
	}
	v4 = nox_common_randomInt_415FA0(0, 99);
	v5 = 0;
	if (*(uint32_t*)v3) {
		v6 = v3;
		v7 = v3;
		do {
			if (*((uint32_t*)v6 + 2) > v4) {
				break;
			}
			v8 = *((uint32_t*)v7 + 3);
			v7 += 12;
			++v5;
			v6 = v7;
		} while (v8);
	}
	result = *(char**)&v3[12 * v5];
	v10 = (int)&v3[12 * v5];
	if (result) {
		result = *(char**)(v10 + 4);
		for (i = 0; i < (int)result; ++i) {
			v12 = nox_xxx_newObjectByTypeID_4E3810(*(char**)v10);
			if (v12) {
				sub_4ED970(35.0, (float2*)a2, &a3);
				nox_xxx_createAt_4DAA50((int)v12, 0, a3.field_0, a3.field_4);
			}
			result = *(char**)(v10 + 4);
		}
	}
}

//----- (004E7540) --------------------------------------------------------
void sub_4E7540(nox_object_t* a1p, nox_object_t* a2p) {
	int a1 = a1p;
	int a2 = a2p;
	int v3; // eax

	if (a1) {
		if (a2) {
			if (*(uint8_t*)(a1 + 8) & 4 && *(uint8_t*)(a2 + 8) & 4 && a1 != a2) {
				v3 = *(uint32_t*)(a2 + 748);
				*(uint32_t*)(*(uint32_t*)(v3 + 276) + 3604) =
					*(unsigned char*)(*(uint32_t*)(*(uint32_t*)(a1 + 748) + 276) + 2064);
				*(uint32_t*)(*(uint32_t*)(v3 + 276) + 3608) = gameFrame();
				*(uint32_t*)(*(uint32_t*)(v3 + 276) + 3600) = 1;
			}
		}
	}
}

//----- (004E75B0) --------------------------------------------------------
char nox_xxx_objectSetOn_4E75B0(nox_object_t* obj) {
	int a1 = obj;
	int v1; // eax
	int v2; // eax

	if (!(*(uint32_t*)(a1 + 16) & 0x1000000)) {
		v1 = *(uint32_t*)(a1 + 8);
		if (v1 & 0x4000) {
			nox_xxx_aud_501960(235, a1, 0, 0);
		}
	}
	nox_xxx_unitSetOnOff_4E4670(a1, 1);
	v2 = *(uint32_t*)(a1 + 8);
	if (v2 & 0x10042000) {
		*(uint32_t*)(a1 + 16) &= 0xFFFFFFBF;
	}
	if (!(v2 & 1)) {
		LOBYTE(v2) = nox_xxx_unitHasCollideOrUpdateFn_537610(a1);
	}
	return v2;
}

//----- (004E7600) --------------------------------------------------------
int nox_xxx_objectSetOff_4E7600(nox_object_t* obj) {
	int a1 = obj;
	int v1;     // eax
	int result; // eax

	if (*(uint32_t*)(a1 + 16) & 0x1000000) {
		v1 = *(uint32_t*)(a1 + 8);
		if (v1 & 0x4000) {
			nox_xxx_aud_501960(236, a1, 0, 0);
		}
	}
	nox_xxx_unitSetOnOff_4E4670(a1, 0);
	result = *(uint32_t*)(a1 + 8);
	if (result & 0x10042000) {
		result = *(uint32_t*)(a1 + 16);
		LOBYTE(result) = result | 0x40;
		*(uint32_t*)(a1 + 16) = result;
	}
	return result;
}

//----- (004E7700) --------------------------------------------------------
int sub_4E7700(int a1) {
	unsigned short* v1; // ecx
	int result;         // eax

	v1 = *(unsigned short**)(a1 + 556);
	result = *(uint32_t*)(a1 + 340) ^ *(uint32_t*)(a1 + 248) ^ *(uint32_t*)(a1 + 120) ^ *(uint32_t*)(a1 + 128) ^
			 *(uint32_t*)(a1 + 132) ^ *(uint32_t*)(a1 + 136) ^ *(uint32_t*)(a1 + 148) ^ *(uint32_t*)(a1 + 152) ^
			 *(short*)(a1 + 124) ^ *(short*)(a1 + 126) ^ *(uint32_t*)(a1 + 108) ^ *(uint32_t*)(a1 + 104) ^
			 *(uint32_t*)(a1 + 100) ^ *(uint32_t*)(a1 + 96) ^ *(uint32_t*)(a1 + 92) ^ *(uint32_t*)(a1 + 88) ^
			 *(uint32_t*)(a1 + 84) ^ *(uint32_t*)(a1 + 80) ^ *(uint32_t*)(a1 + 76) ^ *(uint32_t*)(a1 + 72) ^
			 *(uint32_t*)(a1 + 68) ^ *(uint32_t*)(a1 + 64) ^ *(uint32_t*)(a1 + 60) ^ *(uint32_t*)(a1 + 16) ^
			 *(uint32_t*)(a1 + 20) ^ *(uint32_t*)(a1 + 36) ^ *(uint32_t*)(a1 + 40) ^ *(uint32_t*)(a1 + 44) ^
			 *(uint32_t*)(a1 + 56) ^ *(unsigned short*)(a1 + 4) ^ *(unsigned char*)(a1 + 52);
	if (v1) {
		result ^= *v1 ^ v1[1] ^ v1[2];
	}
	return result;
}

//----- (004E7980) --------------------------------------------------------
int nox_xxx_inventoryGetFirst_4E7980(int a1) { return *(uint32_t*)(a1 + 504); }

//----- (004E7990) --------------------------------------------------------
int nox_xxx_inventoryGetNext_4E7990(int a1) {
	int result; // eax

	if (a1) {
		result = *(uint32_t*)(a1 + 496);
	} else {
		result = 0;
	}
	return result;
}

//----- (004E79B0) --------------------------------------------------------
int sub_4E79B0(int a1) {
	int result; // eax

	result = a1;
	*getMemU32Ptr(0x5D4594, 1567712) = a1;
	return result;
}

//----- (004E79C0) --------------------------------------------------------
int sub_50B510();
char nox_xxx_unitFreeze_4E79C0(nox_object_t* obj, int a2) {
	int a1 = obj;
	int v2; // eax
	int i;  // esi

	v2 = *(uint32_t*)(a1 + 16);
	if (!(v2 & 2)) {
		LOBYTE(v2) = v2 | 2;
		*(uint32_t*)(a1 + 16) = v2;
		if (*(uint8_t*)(a1 + 8) & 4) {
			if (!*getMemU32Ptr(0x5D4594, 1567712)) {
				*getMemU32Ptr(0x5D4594, 1567712) = a2;
			}
			nox_xxx_netReportPlrStatus_4D8270(a1);
			nox_xxx_playerSetState_4FA020((uint32_t*)a1, 13);
			nox_xxx_unitRaise_4E46F0(a1, 0.0);
			sub_50B510();
			for (i = *(uint32_t*)(a1 + 516); i; i = *(uint32_t*)(i + 512)) {
				if (*(uint8_t*)(i + 8) & 2 && (*(uint8_t*)(*(uint32_t*)(i + 748) + 1440) & 0x80)) {
					nox_xxx_unitFreeze_4E79C0(i, a2);
				}
			}
		}
		LOBYTE(v2) = *(uint8_t*)(a1 + 8);
		if (v2 & 2) {
			v2 = *(uint32_t*)(a1 + 16);
			if ((v2 & 0x8000) == 0) {
				LOBYTE(v2) = (unsigned int)nox_xxx_monsterPushAction_50A260(a1, 0);
			}
		}
	}
	return v2;
}

//----- (004E7A60) --------------------------------------------------------
char nox_xxx_unitUnFreeze_4E7A60(nox_object_t* obj, int a2) {
	int a1 = obj;
	int v2; // eax
	int i;  // esi

	v2 = *(uint32_t*)(a1 + 16);
	if (v2 & 2) {
		if (*(uint8_t*)(a1 + 8) & 4) {
			LOBYTE(v2) = getMemByte(0x5D4594, 1567712);
			if (*getMemU32Ptr(0x5D4594, 1567712) && !a2) {
				return v2;
			}
			*getMemU32Ptr(0x5D4594, 1567712) = 0;
			*(uint32_t*)(a1 + 16) &= 0xFFFFFFFD;
			LOBYTE(v2) = nox_xxx_netReportPlrStatus_4D8270(a1);
			for (i = *(uint32_t*)(a1 + 516); i; i = *(uint32_t*)(i + 512)) {
				if (*(uint8_t*)(i + 8) & 2) {
					v2 = *(uint32_t*)(i + 748);
					if (*(uint8_t*)(v2 + 1440) & 0x80) {
						LOBYTE(v2) = nox_xxx_unitUnFreeze_4E7A60(i, a2);
					}
				}
			}
		} else {
			LOBYTE(v2) = v2 & 0xFD;
			*(uint32_t*)(a1 + 16) = v2;
		}
		if (*(uint8_t*)(a1 + 8) & 2) {
			v2 = *(uint32_t*)(a1 + 16);
			if ((v2 & 0x8000) == 0) {
				LOBYTE(v2) = nox_xxx_monsterPopAction_50A160(a1);
			}
		}
	}
	return v2;
}

//----- (004E7B00) --------------------------------------------------------
void nox_xxx_unitBecomePet_4E7B00(int a1, int a2) {
	int v2; // ecx
	int v3; // edi

	if (a1) {
		if (a2) {
			v2 = *(uint32_t*)(a2 + 12);
			v3 = *(uint32_t*)(a1 + 748);
			LOBYTE(v2) = v2 | 0x80;
			*(uint32_t*)(a2 + 12) = v2;
			nox_xxx_netMonitorCreature_4D9250(*(unsigned char*)(*(uint32_t*)(v3 + 276) + 2064), a2);
			nox_xxx_netMarkMinimapObject_417190(*(unsigned char*)(*(uint32_t*)(v3 + 276) + 2064), a2, 1);
			nox_xxx_unitSetOwner_4EC290(a1, a2);
		}
	}
}

//----- (004E7B60) --------------------------------------------------------
void nox_xxx_monsterRemoveMonitors_4E7B60(nox_object_t* a1p, nox_object_t* a2p) {
	int a1 = a1p;
	uint32_t* a2 = a2p;
	int v2; // edi
	int v3; // edx

	v2 = *(uint32_t*)(a1 + 748);
	if (a1) {
		if (a2) {
			v3 = a2[3];
			LOBYTE(v3) = v3 & 0x7F;
			a2[3] = v3;
			nox_xxx_netSendUnMonitorCrea_4D92A0(*(unsigned char*)(*(uint32_t*)(v2 + 276) + 2064), a2);
			nox_xxx_netUnmarkMinimapObj_417300(*(unsigned char*)(*(uint32_t*)(v2 + 276) + 2064), (int)a2, 1);
			nox_xxx_unitClearOwner_4EC300((int)a2);
		}
	}
}

//----- (004E7BC0) --------------------------------------------------------
int sub_4E7BC0(int a1) {
	int result; // eax

	result = a1;
	if (a1) {
		result = (*(uint32_t*)(a1 + 8) >> 2) & 1;
	}
	return result;
}

//----- (004E7BE0) --------------------------------------------------------
int nox_xxx_unitIsCrown_4E7BE0(int a1) {
	int v1; // eax
	int v2; // ecx

	v1 = *getMemU32Ptr(0x5D4594, 1567716);
	if (!*getMemU32Ptr(0x5D4594, 1567716)) {
		v1 = nox_xxx_getNameId_4E3AA0("Crown");
		*getMemU32Ptr(0x5D4594, 1567716) = v1;
	}
	v2 = *(uint32_t*)(a1 + 516);
	if (!v2) {
		return 0;
	}
	while (*(unsigned short*)(v2 + 4) != v1) {
		v2 = *(uint32_t*)(v2 + 512);
		if (!v2) {
			return 0;
		}
	}
	return 1;
}

//----- (004E7C30) --------------------------------------------------------
int nox_xxx_unitIsGameball_4E7C30(int a1) {
	int v1; // eax
	int v2; // ecx

	v1 = *getMemU32Ptr(0x5D4594, 1567720);
	if (!*getMemU32Ptr(0x5D4594, 1567720)) {
		v1 = nox_xxx_getNameId_4E3AA0("GameBall");
		*getMemU32Ptr(0x5D4594, 1567720) = v1;
	}
	v2 = *(uint32_t*)(a1 + 516);
	if (!v2) {
		return 0;
	}
	while (*(unsigned short*)(v2 + 4) != v1) {
		v2 = *(uint32_t*)(v2 + 512);
		if (!v2) {
			return 0;
		}
	}
	return 1;
}

//----- (004E7CF0) --------------------------------------------------------
int nox_xxx_unitCountSlaves_4E7CF0(int a1, int a2, int a3) {
	int result;  // eax
	uint32_t* i; // ecx

	result = 0;
	if (!a1 || !a2 || !a3) {
		return 0;
	}
	for (i = *(uint32_t**)(a1 + 516); i; i = (uint32_t*)i[128]) {
		if (a2 & i[2]) {
			if (a3 & i[3]) {
				++result;
			}
		}
	}
	return result;
}

//----- (004E7DE0) --------------------------------------------------------
int sub_4E7DE0(int a1, nox_object_t* item) {
	int v2;       // ebx
	uint32_t* v3; // eax
	int v4;       // ecx
	int v5;       // edx
	int v6;       // eax
	bool v7;      // zf

	if (!a1 || !item || *(uint16_t*)(a1 + 4) != *(uint16_t*)&item->typ_ind) {
		return 0;
	}
	v2 = *(uint32_t*)(a1 + 8);
	if (v2 & 0x13001000) {
		v3 = *(uint32_t**)(a1 + 692);
		v4 = 0;
		v5 = (int)item->init_data - (uint32_t)v3;
		while (*v3 == *(uint32_t*)((char*)v3 + v5)) {
			++v4;
			++v3;
			if (v4 >= 4) {
				goto LABEL_8;
			}
		}
		return 0;
	}
LABEL_8:
	if (!(v2 & 0x100)) {
		return 1;
	}
	v6 = *(uint32_t*)(a1 + 12);
	if (v6 & 1) {
		v7 = **(uint8_t**)(a1 + 736) == *(uint8_t*)item->use_data;
	} else {
		if (!(v6 & 2)) {
			if (**(uint8_t**)(a1 + 736) != *(uint8_t*)item->use_data) {
				return 0;
			}
			return 1;
		}
		v7 = strcmp(*(const char**)(a1 + 736), item->use_data) == 0;
	}
	if (!v7) {
		return 0;
	}
	return 1;
}

//----- (004E7F10) --------------------------------------------------------
char* nox_xxx_unitPostCreateNotify_4E7F10(nox_object_t* a1p) {
	int a1 = a1p;
	char* result; // eax
	int i;        // edi
	int v3;       // eax
	int v4;       // esi
	int v5;       // eax

	*(uint32_t*)(a1 + 140) = 0;
	*(uint32_t*)(a1 + 144) = 0;
	result = nox_common_playerInfoGetFirst_416EA0();
	for (i = (int)result; result; i = (int)result) {
		v3 = *(uint32_t*)(i + 2056);
		v4 = 1 << *(uint8_t*)(i + 2064);
		if (v3) {
			if (nox_xxx_unitIsHostileMimic_4E7F90(v3, a1) == 1) {
				v5 = v4 | *(uint32_t*)(a1 + 144);
				*(uint32_t*)(a1 + 140) |= v4;
				*(uint32_t*)(a1 + 144) = v5;
			}
		}
		result = nox_common_playerInfoGetNext_416EE0(i);
	}
	return result;
}

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

//----- (004E8110) --------------------------------------------------------
char* sub_4E8110(int a1) {
	int v1;       // edi
	char* result; // eax
	char* v3;     // ebx
	char* v4;     // esi
	int v5;       // ebp
	char v6;      // al
	int v7;       // ecx
	int v8;       // eax
	bool v9;      // zf
	int v10;      // eax
	int v11;      // ecx

	v1 = 1 << a1;
	result = nox_common_playerInfoFromNum_417090(a1);
	v3 = result;
	if (result) {
		result = (char*)nox_server_getFirstObject_4DA790();
		v4 = result;
		if (result) {
			v5 = ~v1;
			do {
				v6 = v4[8];
				v7 = v5 & *((uint32_t*)v4 + 35);
				*((uint32_t*)v4 + 36) &= v5;
				*((uint32_t*)v4 + 35) = v7;
				if ((v6 & 6)) {
					v8 = *((uint32_t*)v3 + 514);
					if (v8) {
						v9 = nox_xxx_unitIsHostileMimic_4E7F90(v8, (int)v4) == 1;
						v10 = *((uint32_t*)v4 + 36);
						if (v9) {
							if (!(v10 & v1)) {
								v11 = v1 | v10;
								*((uint32_t*)v4 + 36) = v11;
								*((uint32_t*)v4 + 35) |= v1;
							}
						} else if (v10 & v1) {
							v11 = v5 & *((uint32_t*)v4 + 36);
							*((uint32_t*)v4 + 36) = v11;
							*((uint32_t*)v4 + 35) |= v1;
						}
					}
				}
				result = (char*)nox_server_getNextObject_4DA7A0((int)v4);
				v4 = result;
			} while (result);
		}
	}
	return result;
}

//----- (004E81D0) --------------------------------------------------------
int sub_4E81D0(nox_object_t* a1p) {
	int a1 = a1p;
	int result; // eax

	result = *getMemU32Ptr(0x5D4594, 1567728);
	if (!*getMemU32Ptr(0x5D4594, 1567728)) {
		result = nox_xxx_getNameId_4E3AA0("Pixie");
		*getMemU32Ptr(0x5D4594, 1567728) = result;
	}
	if (a1) {
		if (*(unsigned short*)(a1 + 4) == result) {
			result = *(uint32_t*)(a1 + 748);
			*(uint32_t*)(result + 4) = 0;
		}
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

//----- (004E8340) --------------------------------------------------------
void nox_xxx_fnFindCloseDoors_4E8340(float* a1, int a2) {
	int v2; // eax

	if (*((uint8_t*)a1 + 8) & 0x80) {
		v2 = *((uint32_t*)a1 + 187);
		if (*(uint32_t*)(v2 + 16) == *(uint32_t*)a2 && *(uint32_t*)(v2 + 20) == *(uint32_t*)(a2 + 4)) {
			*(uint8_t*)(v2 + 1) = 0;
			if (nox_common_gameFlags_check_40A5C0(4096)) {
				sub_4E8390((int)a1);
			}
		}
	}
}

//----- (004E8390) --------------------------------------------------------
int sub_4E8390(int a1) {
	*(uint8_t*)(*(uint32_t*)(a1 + 748) + 48) = 1;
	return sub_4D6A20(255, a1);
}

int nox_objectCollideDefault(int a1, int a2, float* a3) {
	return 0;
}

//----- (004E83B0) --------------------------------------------------------
unsigned char* nox_xxx_collideMonsterEventProc_4E83B0(int a1, int a2) {
	return nox_xxx_scriptCallByEventBlock_502490((int*)(*(uint32_t*)(a1 + 748) + 1272), a2, a1, 22);
}

//----- (004E83D0) --------------------------------------------------------
unsigned char* nox_xxx_collideMimic_4E83D0(int a1, int a2) {
	int v2;  // eax
	int* v3; // eax
	int* v4; // eax

	if (a2) {
		v2 = *(uint32_t*)(a2 + 16);
		if ((v2 & 0x8000) == 0 && *(uint8_t*)(a2 + 8) & 6 && nox_xxx_unitIsEnemyTo_5330C0(a1, a2) &&
			!nox_xxx_monsterIsActionScheduled_50A090(a1, 15)) {
			v3 = nox_xxx_monsterPushAction_50A260(a1, 43);
			if (v3) {
				v3[1] = gameFrame();
			}
			v4 = nox_xxx_monsterPushAction_50A260(a1, 15);
			if (v4) {
				v4[1] = *(uint32_t*)(a2 + 56);
				v4[2] = *(uint32_t*)(a2 + 60);
				v4[3] = gameFrame();
			}
		}
	}
	return nox_xxx_collideMonsterEventProc_4E83B0(a1, a2);
}

//----- (004E8460) --------------------------------------------------------
uint32_t nox_xxx_wallFlags(int i);
void nox_xxx_collidePlayer_4E8460(int a1, int a2) {
	int v2;       // ecx
	uint16_t* v3; // eax
	int v4;       // ebx
	int v5;       // eax
	short v6;     // ax
	int v7;       // eax
	short v8;     // ax
	int v9;       // eax
	int v10;      // eax
	int v11;      // eax
	char v12;     // al
	float v13;    // [esp+0h] [ebp-1Ch]
	float v14;    // [esp+4h] [ebp-18h]
	int v15;      // [esp+4h] [ebp-18h]
	float v16;    // [esp+Ch] [ebp-10h]
	float v17;    // [esp+Ch] [ebp-10h]
	float v18;    // [esp+Ch] [ebp-10h]
	float v19;    // [esp+10h] [ebp-Ch]

	if (!nox_common_playerIsAbilityActive_4FC250(a1, 1)) {
		goto LABEL_26;
	}
	if (!a2) {
		goto LABEL_14;
	}
	v2 = *(uint32_t*)(a2 + 8);
	if (!(v2 & 6) || (v3 = *(uint16_t**)(a2 + 556), !*v3) && v3[2]) {
		if (!(v2 & 0x400000) && (signed char)*(uint8_t*)(a2 + 16) >= 0 &&
			*(float*)(a2 + 120) <= (double)*(float*)(a1 + 120)) {
			goto LABEL_26;
		}
	}
	if ((v2 & 0x80u) == 0) {
		if (!(*(uint8_t*)(a2 + 16) & 9) && !(v2 & 1)) {
			goto LABEL_14;
		} else {
			goto LABEL_26;
		}
	} else if (*(uint8_t*)(*(uint32_t*)(a2 + 748) + 1)) {
		goto LABEL_14;
	} else {
		goto LABEL_26;
	}
LABEL_14:
	nox_xxx_playerSetState_4FA020((uint32_t*)a1, 13);
	nox_xxx_earthquakeSend_4D9110((float*)(a1 + 56), 10);
	sub_4FC300((uint32_t*)a1, 1);
	if (a2) {
		v16 = nox_xxx_gamedataGetFloat_419D40("BerserkerDamage");
		v4 = nox_float2int(v16);
		if (!(*(uint32_t*)(a2 + 8) & 0x400000)) {
			sub_4E86E0(a1, (float*)a2);
		}
		v5 = nox_xxx_findParentChainPlayer_4EC580(a1);
		(*(void (**)(int, int, int, int, int))(a2 + 716))(a2, v5, a1, v4, 2);
		if (*(uint32_t*)(a2 + 8) & 0x20006) {
			nox_xxx_unitMove_4E7010(a1, (float2*)(a1 + 72));
			goto LABEL_26;
		}
		v17 = nox_xxx_gamedataGetFloat_419D40("BerserkerStunDuration");
		v6 = nox_float2int(v17);
		nox_xxx_buffApplyTo_4FF380(a1, 5, v6, 5);
	} else {
		v7 = *(uint32_t*)(*(uint32_t*)(a1 + 748) + 296);
		if (v7 && !(nox_xxx_wallFlags(*(unsigned char*)(v7 + 1)) & 5)) {
			nox_xxx_unitMove_4E7010(a1, (float2*)(a1 + 72));
			goto LABEL_26;
		}
		nox_xxx_aud_501960(171, a1, 0, 0);
		v18 = nox_xxx_gamedataGetFloat_419D40("BerserkerStunDuration");
		v8 = nox_float2int(v18);
		nox_xxx_buffApplyTo_4FF380(a1, 5, v8, 5);
		v14 = *(float*)(a1 + 68) * 0.043478262;
		v15 = nox_float2int(v14);
		v13 = *(float*)(a1 + 64) * 0.043478262;
		v9 = nox_float2int(v13);
		nox_xxx_damageToMap_534BC0(v9, v15, 100, 2, a1);
	}
	v19 = nox_xxx_gamedataGetFloat_419D40("BerserkerPainRatio") * (double)**(unsigned short**)(a1 + 556);
	v10 = nox_float2int(v19);
	if (v10 < 1) {
		v10 = 1;
	}
	nox_xxx_unitDamageClear_4EE5E0(a1, v10);
	nox_xxx_unitMove_4E7010(a1, (float2*)(a1 + 72));
LABEL_26:
	if (a2) {
		if (*(uint8_t*)(a2 + 8) & 4) {
			v11 = *(uint32_t*)(a2 + 16);
			if ((v11 & 0x8000) == 0) {
				if (nox_xxx_testUnitBuffs_4FF350(a1, 16)) {
					if (nox_xxx_unitGetBuffTimer_4FF550(a1, 16) < (unsigned int)(14 * gameFPS())) {
						v12 = nox_xxx_buffGetPower_4FF570(a1, 16);
						nox_xxx_buffApplyTo_4FF380(a2, 16, 15 * (uint16_t)gameFPS(), v12);
						nox_xxx_spellBuffOff_4FF5B0(a1, 16);
					}
				}
			}
		}
	}
}

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

//----- (004F0390) --------------------------------------------------------
uint32_t* nox_xxx_unitSparkInit_4F0390(int a1) {
	uint32_t* result; // eax

	result = *(uint32_t**)(a1 + 748);
	result[1] = 32;
	*result = 32;
	return result;
}

//----- (004F03B0) --------------------------------------------------------
int nox_xxx_initFrog_4F03B0(int a1) {
	uint8_t* v1; // esi
	int result;  // eax

	v1 = *(uint8_t**)(a1 + 748);
	*v1 = nox_common_randomInt_415FA0(55, 60);
	v1[1] = 1;
	v1[2] = 0;
	result = nox_common_randomInt_415FA0(0, 255);
	*(uint16_t*)(a1 + 126) = result;
	return result;
}

//----- (004F0400) --------------------------------------------------------
int* nox_xxx_initChest_4F0400(int a1) {
	int* result; // eax

	result = (int*)a1;
	if (!(*(uint8_t*)(a1 + 20) & 0xE)) {
		nox_xxx_unitSetXStatus_4E4800(a1, (int*)2);
	}
	return result;
}

//----- (004F0420) --------------------------------------------------------
uint32_t* nox_xxx_unitBoulderInit_4F0420(uint32_t* a1) {
	uint32_t* result; // eax
	int v2;           // edx

	result = a1;
	v2 = a1[15];
	a1[39] = a1[14];
	a1[40] = v2;
	return result;
}

//----- (004F0450) --------------------------------------------------------
int sub_4F0450(int a1) {
	int v1;     // edi
	short v2;   // ax
	int result; // eax

	v1 = *(uint32_t*)(a1 + 748);
	v2 = nox_xxx_xferDirectionToAngle_509E00(*(uint32_t**)(a1 + 692));
	*(uint16_t*)(a1 + 126) = v2;
	*(uint16_t*)(a1 + 124) = v2;
	result = nox_xxx_getNameId_4E3AA0((char*)(v1 + 16));
	*(uint32_t*)(v1 + 12) = result;
	return result;
}

//----- (004F0490) --------------------------------------------------------
int sub_4F0490(int a1) {
	int result; // eax

	result = nox_xxx_xferDirectionToAngle_509E00(*(uint32_t**)(a1 + 692));
	*(uint16_t*)(a1 + 126) = result;
	*(uint16_t*)(a1 + 124) = result;
	return result;
}

//----- (004F04B0) --------------------------------------------------------
int nox_xxx_unitInitGold_4F04B0(int a1) {
	int result;   // eax
	uint32_t* v2; // edi
	char* v3;     // eax
	int v4;       // esi
	int v5;       // ecx
	int v6;       // esi
	int v7;       // [esp+4h] [ebp-4h]
	float v8;     // [esp+Ch] [ebp+4h]
	float v9;     // [esp+Ch] [ebp+4h]

	result = a1;
	v2 = *(uint32_t**)(a1 + 692);
	if (!*v2) {
		v3 = nox_common_playerInfoGetFirst_416EA0();
		v4 = 0;
		v7 = 0;
		v8 = 0.0;
		if (v3) {
			do {
				v5 = *((uint32_t*)v3 + 514);
				if (v5) {
					v8 = v8 + *(float*)(v5 + 28);
				}
				v3 = nox_common_playerInfoGetNext_416EE0((int)v3);
				++v4;
			} while (v3);
			v7 = v4;
		}
		v9 = v8 / (double)v7;
		v6 = nox_common_randomInt_415FA0((long long)(v9 * *(double*)&qword_581450_10256),
										 (long long)(v9 * *getMemDoublePtr(0x581450, 10264))) -
			 (unsigned long long)(long long)(v9 * *getMemDoublePtr(0x581450, 10248));
		result = nox_common_randomInt_415FA0(15, 30);
		*v2 = result + v6;
	}
	return result;
}

//----- (004F0570) --------------------------------------------------------
int* nox_xxx_breakInit_4F0570(int a1) {
	int* result; // eax

	result = (int*)a1;
	if (!(*(uint8_t*)(a1 + 20) & 0xE)) {
		nox_xxx_unitSetXStatus_4E4800(a1, (int*)2);
	}
	return result;
}

//----- (004F0590) --------------------------------------------------------
int nox_xxx_unitInitGenerator_4F0590(int a1) {
	int v1;     // esi
	double v2;  // st7
	int result; // eax

	v1 = *(uint32_t*)(a1 + 748);
	switch (*(unsigned char*)(nox_xxx_getQuestStage_51A930() + v1 + 83)) {
	case 0u:
		v2 = nox_xxx_gamedataGetFloat_419D40("GeneratorMaxActiveCreaturesHigh");
		*(uint8_t*)(v1 + 87) = (long long)v2;
		break;
	case 1u:
		v2 = nox_xxx_gamedataGetFloat_419D40("GeneratorMaxActiveCreaturesNormal");
		*(uint8_t*)(v1 + 87) = (long long)v2;
		break;
	case 2u:
		v2 = nox_xxx_gamedataGetFloat_419D40("GeneratorMaxActiveCreaturesLow");
		*(uint8_t*)(v1 + 87) = (long long)v2;
		break;
	case 3u:
		v2 = nox_xxx_gamedataGetFloat_419D40("GeneratorMaxActiveCreaturesSingular");
		*(uint8_t*)(v1 + 87) = (long long)v2;
		break;
	default:
		break;
	}
	result = *(uint32_t*)(a1 + 12);
	if (result & 1) {
		result = nox_xxx_mathDirection4ToAngle_509E90(0);
		*(uint16_t*)(a1 + 124) = result;
		*(uint16_t*)(a1 + 126) = *(uint16_t*)(a1 + 124);
		return result;
	}
	if (result & 2) {
		result = nox_xxx_mathDirection4ToAngle_509E90(2);
		*(uint16_t*)(a1 + 124) = result;
		*(uint16_t*)(a1 + 126) = *(uint16_t*)(a1 + 124);
		return result;
	}
	if (result & 4) {
		result = nox_xxx_mathDirection4ToAngle_509E90(8);
		*(uint16_t*)(a1 + 124) = result;
		*(uint16_t*)(a1 + 126) = *(uint16_t*)(a1 + 124);
		return result;
	}
	if (result & 8) {
		result = nox_xxx_mathDirection4ToAngle_509E90(6);
		*(uint16_t*)(a1 + 124) = result;
		*(uint16_t*)(a1 + 126) = *(uint16_t*)(a1 + 124);
		return result;
	}
	*(uint16_t*)(a1 + 126) = *(uint16_t*)(a1 + 124);
	return result;
}

//----- (004F0720) --------------------------------------------------------
uint32_t* nox_server_rewardgen_activateMarker_4F0720(int a1, unsigned int a2) {
	int v2;             // eax
	int v3;             // ebx
	int* v4;            // ebp
	uint32_t* result;   // eax
	int v6;             // esi
	char v7;            // cl
	unsigned char* v8;  // eax
	unsigned int v9;    // eax
	int v10;            // ebp
	unsigned int v11;   // edi
	char v12;           // cl
	unsigned char* v13; // esi
	unsigned int v14;   // edx

	v2 = *getMemU32Ptr(0x5D4594, 1568276);
	v3 = a1;
	v4 = *(int**)(a1 + 692);
	if (!*getMemU32Ptr(0x5D4594, 1568276)) {
		v2 = nox_xxx_getNameId_4E3AA0("RewardMarkerPlus");
		*getMemU32Ptr(0x5D4594, 1568276) = v2;
	}
	if (*(unsigned short*)(a1 + 4) == v2) {
		a2 += 2;
	}
	switch (v4[53]) {
	case 1:
		if (nox_common_randomInt_415FA0(0, 100) > 75) {
			return 0;
		}
		break;
	case 2:
		if (nox_common_randomInt_415FA0(0, 100) > 50) {
			return 0;
		}
		break;
	case 3:
		if (nox_common_randomInt_415FA0(0, 100) > 25) {
			return 0;
		}
		break;
	case 4:
		if (nox_common_randomInt_415FA0(0, 100) > 5) {
			return 0;
		}
		break;
	default:
		break;
	}
	v6 = 0;
	v7 = 0;
	v8 = getMemAt(0x587000, 207044);
	do {
		if ((1 << v7) & *v4) {
			v6 += *v8;
		}
		v8 += 8;
		++v7;
	} while ((int)v8 < (int)getMemAt(0x587000, 207108));
	if (!v6) {
		return 0;
	}
	v9 = nox_common_randomInt_415FA0(1, v6);
	v10 = *v4;
	v11 = 0;
	v12 = 0;
	v13 = getMemAt(0x587000, 207044);
	while (1) {
		v14 = 1 << v12;
		if (v10 & (1 << v12)) {
			v11 += *v13;
			v3 = a1;
			if (v11 >= v9) {
				break;
			}
		}
		v13 += 8;
		++v12;
		if ((int)v13 >= (int)getMemAt(0x587000, 207108)) {
			v14 = a2;
			break;
		}
	}
	switch (v14) {
	case 1u:
		result = nox_xxx_rewardSpellBook_4F09F0(v3, a2);
		break;
	case 2u:
		result = nox_xxx_rewardAbilityBook_4F0C70(v3);
		break;
	case 4u:
		result = nox_xxx_rewardFieldGuide_4F0D20(v3, a2);
		break;
	case 8u:
		result = (uint32_t*)nox_xxx_rewardMakeWeapon_4F14E0(v3, a2);
		break;
	case 0x10u:
		result = nox_xxx_rewardMakeArmor_4F0E80(v3, a2);
		break;
	case 0x20u:
		result = nox_xxx_createGem_4F1D30(v3, a2);
		break;
	case 0x40u:
		result = nox_xxx_rewardMakePotion_4F1C40(v3, a2);
		break;
	case 0x80u:
		result = nox_xxx_createGem2_4F1F00(v3, a2);
		break;
	default:
		result = nox_xxx_createGem_4F1D30(v3, a2);
		break;
	}
	return result;
}

//----- (004F09F0) --------------------------------------------------------
uint32_t* nox_xxx_rewardSpellBook_4F09F0(int a1, unsigned int a2) {
	int v2;             // esi
	int v3;             // ecx
	int i;              // eax
	int v5;             // eax
	int v6;             // edx
	int v7;             // ecx
	uint32_t* result;   // eax
	int v9;             // ebx
	int v10;            // ebp
	int v11;            // ecx
	unsigned char* v12; // eax
	int v13;            // edx
	int v14;            // edi
	int v15;            // esi
	int v16;            // edx
	unsigned char* j;   // ecx
	int v18;            // eax

	v2 = *(uint32_t*)(a1 + 692);
	if (*(uint8_t*)(v2 + 4) & 1) {
		v3 = 0;
		for (i = 0; i < 137; ++i) {
			if (*(uint8_t*)(v2 + i + 8) == 1) {
				++v3;
			}
		}
		if (v3) {
			v5 = nox_common_randomInt_415FA0(0, v3 - 1);
			v6 = 0;
			v7 = 0;
			while (1) {
				if (*(uint8_t*)(v2 + v7 + 8) == 1) {
					if (v6 == v5) {
						v9 = v7;
						goto LABEL_27;
					}
					++v6;
				}
				if (++v7 >= 137) {
					return 0;
				}
			}
		}
		return 0;
	}
	v10 = nox_server_rewardGen_pickRandomSlots_4F0B60(a2);
	v11 = 0;
	if (!*getMemU32Ptr(0x587000, 207108)) {
		return 0;
	}
	v12 = getMemAt(0x587000, 207104);
	do {
		if (v10 & *((uint32_t*)v12 + 2)) {
			v11 += *v12;
		}
		v13 = *((uint32_t*)v12 + 4);
		v12 += 12;
	} while (v13);
	if (!v11) {
		return 0;
	}
	v14 = nox_common_randomInt_415FA0(0, v11 - 1);
	v15 = 0;
	v16 = 0;
	if (!*getMemU32Ptr(0x587000, 207108)) {
		return 0;
	}
	for (j = getMemAt(0x587000, 207104);; j += 12) {
		if (v10 & *((uint32_t*)j + 2)) {
			v15 += *j;
			if (v14 < v15) {
				break;
			}
		}
		v18 = *((uint32_t*)j + 4);
		++v16;
		if (!v18) {
			return 0;
		}
	}
	v9 = *getMemU32Ptr(0x587000, 207108 + 12 * v16);
LABEL_27:
	if (!v9) {
		return 0;
	}
	if (nox_xxx_playerCheckSpellClass_57AEA0(1, v9) || nox_xxx_playerCheckSpellClass_57AEA0(2, v9)) {
		if (nox_xxx_playerCheckSpellClass_57AEA0(1, v9)) {
			if (nox_xxx_playerCheckSpellClass_57AEA0(2, v9)) {
				return 0;
			}
			result = nox_xxx_newObjectByTypeID_4E3810("ConjurerSpellBook");
		} else {
			result = nox_xxx_newObjectByTypeID_4E3810("WizardSpellBook");
		}
	} else {
		result = nox_xxx_newObjectByTypeID_4E3810("CommonSpellBook");
	}
	if (!result) {
		return 0;
	}
	*(uint8_t*)result[184] = v9;
	return result;
}

//----- (004F0B60) --------------------------------------------------------
int nox_server_rewardGen_pickRandomSlots_4F0B60(unsigned int a1) {
	unsigned int v2; // eax
	float* v3;       // eax
	int v4;          // eax
	int v5;          // ecx
	float* v6;       // edx
	double v7;       // st7
	int v8;          // [esp+0h] [ebp-18h]
	float v9[5];     // [esp+4h] [ebp-14h]
	float v10;       // [esp+1Ch] [ebp+4h]

	v9[0] = 0.0;
	v9[1] = 0.0;
	v9[2] = 0.0;
	v9[3] = 0.0;
	v9[4] = 0.0;
	if (a1 > 0xA) {
		return 16;
	}
	switch (a1) {
	case 0u:
		return 1;
	case 1u:
		v9[0] = 87.5;
		v9[1] = 12.5;
		goto LABEL_12;
	case 9u:
		v9[3] = 12.5;
		v9[4] = 87.5;
		goto LABEL_12;
	case 0xAu:
		return 16;
	}
	if (a1 & 1) {
		v2 = a1 >> 1;
		v9[v2] = 75.0;
		*(int*)((char*)&v9 + v2 * 4) = 1095237632; // it's a float actually
		v9[v2 + 1] = 12.5;
	} else {
		v3 = &v9[a1 >> 1];
		*(v3 - 1) = 50.0;
		*v3 = 50.0;
	}
LABEL_12:
	v8 = 545;
	v4 = nox_common_randomInt_415FA0(0, 200);
	v5 = 0;
	v6 = v9;
	v7 = 0.0;
	while (1) {
		v7 = v7 + *v6;
		v10 = (double)v4 * 0.5;
		if (v10 <= v7) {
			break;
		}
		++v5;
		++v6;
		if (v5 >= 5) {
			return 1;
		}
	}
	return 1 << v5;
}

//----- (004F0C70) --------------------------------------------------------
uint32_t* nox_xxx_rewardAbilityBook_4F0C70(int a1) {
	int v1;           // esi
	int v2;           // ecx
	int i;            // eax
	uint32_t* result; // eax
	int v5;           // eax
	int v6;           // edx
	int v7;           // ecx
	int v8;           // ebx

	v1 = *(uint32_t*)(a1 + 692);
	if (!(*(uint8_t*)(v1 + 4) & 2)) {
		v8 = nox_common_randomInt_415FA0(1, 5);
		goto LABEL_16;
	}
	v2 = 0;
	for (i = 0; i < 6; ++i) {
		if (*(uint8_t*)(v1 + i + 145) == 1) {
			++v2;
		}
	}
	if (!v2) {
		return 0;
	}
	v5 = nox_common_randomInt_415FA0(0, v2 - 1);
	v6 = 0;
	v7 = 0;
	while (1) {
		if (*(uint8_t*)(v1 + v7 + 145) == 1) {
			if (v6 != v5) {
				++v6;
				goto LABEL_12;
			}
			v8 = v7;
			break;
		}
	LABEL_12:
		if (++v7 >= 6) {
			return 0;
		}
	}
LABEL_16:
	if (!v8) {
		return 0;
	}
	result = nox_xxx_newObjectByTypeID_4E3810("AbilityBook");
	if (result) {
		*(uint8_t*)result[184] = v8;
	}
	return result;
}

//----- (004F0D20) --------------------------------------------------------
uint32_t* nox_xxx_rewardFieldGuide_4F0D20(int a1, unsigned int a2) {
	int v2;             // esi
	int v3;             // ecx
	int i;              // eax
	uint32_t* result;   // eax
	int v6;             // eax
	int v7;             // edx
	int v8;             // ecx
	int v9;             // esi
	int v10;            // ebp
	int v11;            // ecx
	unsigned char* v12; // eax
	int v13;            // edx
	int v14;            // edi
	int v15;            // esi
	int v16;            // edx
	unsigned char* j;   // ecx
	int v18;            // eax
	uint32_t* v19;      // ebx
	char* v20;          // ebp
	char* v21;          // edi

	v2 = *(uint32_t*)(a1 + 692);
	if (*(uint8_t*)(v2 + 4) & 4) {
		v3 = 0;
		for (i = 0; i < 41; ++i) {
			if (*(uint8_t*)(v2 + i + 151) == 1) {
				++v3;
			}
		}
		if (v3) {
			v6 = nox_common_randomInt_415FA0(0, v3 - 1);
			v7 = 0;
			v8 = 0;
			while (1) {
				if (*(uint8_t*)(v2 + v8 + 151) == 1) {
					if (v7 == v6) {
						v9 = v8;
						goto LABEL_29;
					}
					++v7;
				}
				if (++v8 >= 41) {
					return 0;
				}
			}
		}
		return 0;
	}
	v10 = nox_server_rewardGen_pickRandomSlots_4F0B60(a2);
	v11 = 0;
	if (!*getMemU32Ptr(0x587000, 207796)) {
		return 0;
	}
	v12 = getMemAt(0x587000, 207792);
	do {
		if (v10 & *((uint32_t*)v12 + 2)) {
			v11 += *v12;
		}
		v13 = *((uint32_t*)v12 + 4);
		v12 += 12;
	} while (v13);
	if (!v11) {
		return 0;
	}
	v14 = nox_common_randomInt_415FA0(0, v11 - 1);
	v15 = 0;
	v16 = 0;
	if (!*getMemU32Ptr(0x587000, 207796)) {
		return 0;
	}
	for (j = getMemAt(0x587000, 207792);; j += 12) {
		if (v10 & *((uint32_t*)j + 2)) {
			v15 += *j;
			if (v14 < v15) {
				break;
			}
		}
		v18 = *((uint32_t*)j + 4);
		++v16;
		if (!v18) {
			return 0;
		}
	}
	v9 = *getMemU32Ptr(0x587000, 207796 + 12 * v16);
LABEL_29:
	if (!v9) {
		return 0;
	}
	result = nox_xxx_newObjectByTypeID_4E3810("FieldGuide");
	v19 = result;
	if (result) {
		v20 = (char*)result[184];
		v21 = nox_xxx_guideNameByN_427230(v9);
		result = v19;
		strcpy(v20, v21);
	}
	return result;
}

//----- (004F0E80) --------------------------------------------------------
uint32_t* nox_xxx_rewardMakeArmor_4F0E80(int a1, unsigned int a2) {
	int v2;                // ebx
	int v3;                // edi
	unsigned char* v4;     // esi
	int v5;                // eax
	int v7;                // ebx
	int v8;                // ebp
	int v9;                // edi
	unsigned char* i;      // esi
	int v11;               // eax
	int v12;               // esi
	int v13;               // ebp
	uint32_t* v14;         // edi
	int v15;               // eax
	int v16;               // eax
	short v17;             // bx
	int v18;               // eax
	int v19;               // ecx
	unsigned char* v20;    // eax
	int v21;               // edx
	int v22;               // ecx
	unsigned char* v23;    // eax
	int v24;               // edx
	int v25;               // edx
	unsigned char* v26;    // eax
	int v27;               // ecx
	int v28;               // ecx
	unsigned char* v29;    // eax
	int v30;               // edx
	int v31;               // eax
	int v32;               // esi
	int v33;               // edx
	unsigned char* v34;    // ecx
	int v35;               // edi
	int v36;               // ecx
	unsigned char* v37;    // eax
	int v38;               // esi
	int v39;               // eax
	int v40;               // esi
	int v41;               // edx
	unsigned char* v42;    // ecx
	int v43;               // edi
	int v44;               // edx
	unsigned char* v45;    // eax
	int v46;               // ecx
	int v47;               // eax
	int v48;               // edi
	int v49;               // esi
	unsigned char* v50;    // ecx
	int v51;               // edx
	signed int v52;        // eax
	int v53;               // ecx
	int v54;               // eax
	int v55;               // ebx
	int v56;               // edx
	unsigned char* v57;    // eax
	int v58;               // ecx
	int v59;               // eax
	int v60;               // edi
	int v61;               // esi
	unsigned char* v62;    // ecx
	int v63;               // edx
	int v64;               // [esp+10h] [ebp-20h]
	int v65;               // [esp+14h] [ebp-1Ch]
	int v66;               // [esp+18h] [ebp-18h]
	unsigned char v67[20]; // [esp+1Ch] [ebp-14h]

	v65 = 0;
	v2 = nox_server_rewardGen_pickRandomSlots_4F0B60(a2);
	v3 = 0;
	v64 = v2;
	if (!*getMemU32Ptr(0x587000, 208180)) {
		return 0;
	}
	v4 = getMemAt(0x587000, 208192);
	do {
		if (*(v4 - 4) & 2 && v2 & *(uint32_t*)v4 && nox_xxx_getUnitDefDd10_4E3BA0(*((uint32_t*)v4 - 2))) {
			v3 += *(v4 - 16);
		}
		v5 = *((uint32_t*)v4 + 2);
		v4 += 20;
	} while (v5);
	if (!v3) {
		return 0;
	}
	v7 = nox_common_randomInt_415FA0(0, v3 - 1);
	v8 = 0;
	v9 = 0;
	if (!*getMemU32Ptr(0x587000, 208180)) {
		return 0;
	}
	for (i = getMemAt(0x587000, 208192);; i += 20) {
		if (*(i - 4) & 2) {
			if (v64 & *(uint32_t*)i) {
				if (nox_xxx_getUnitDefDd10_4E3BA0(*((uint32_t*)i - 2))) {
					v8 += *(i - 16);
					if (v7 < v8) {
						break;
					}
				}
			}
		}
		v11 = *((uint32_t*)i + 2);
		++v9;
		if (!v11) {
			return 0;
		}
	}
	v12 = *getMemU32Ptr(0x587000, 208184 + 20 * v9);
	if (!v12) {
		return 0;
	}
	v13 = sub_415D10(*(char**)getMemAt(0x587000, 208184 + 20 * v9));
	v14 = nox_xxx_newObjectWithTypeInd_4E3450(v12);
	v66 = (int)v14;
	if (!v14) {
		return 0;
	}
	switch (v64) {
	case 2:
		v15 = nox_common_randomInt_415FA0(0, 1);
		break;
	case 4:
		v15 = nox_common_randomInt_415FA0(0, 2);
		break;
	case 8:
		v15 = nox_common_randomInt_415FA0(1, 3);
		break;
	case 16:
		v15 = nox_common_randomInt_415FA0(2, 4);
		break;
	default:
		return v14;
	}
	if (!v15) {
		return v14;
	}
	*(uint32_t*)v67 = 0;
	*(uint32_t*)&v67[4] = 0;
	*(uint32_t*)&v67[8] = 0;
	*(uint32_t*)&v67[12] = 0;
	*(uint16_t*)&v67[16] = 0;
	*(uint16_t*)&v67[18] = 0;
	switch (v15) {
	case 1:
		v16 = nox_common_randomInt_415FA0(1, 100);
		if (v16 > 20) {
			v17 = (v16 > 50) + 1;
		} else {
			v17 = 4;
		}
		break;
	case 2:
		v18 = nox_common_randomInt_415FA0(1, 100);
		if (v18 > 12) {
			v17 = v18 > 25 ? 3 : 6;
		} else {
			v17 = 5;
		}
		break;
	case 3:
		v17 = 7;
		break;
	default:
		v17 = 15;
		if (v15 != 4) {
			v17 = a2;
		}
		break;
	}
	if (v17 & 1) {
		v19 = 0;
		if (!*getMemU32Ptr(0x587000, 210856)) {
			// nop
		} else {
			v20 = getMemAt(0x587000, 210852);
			do {
				if (v64 & *((uint32_t*)v20 + 2) && v13 & *(uint32_t*)(*(uint32_t*)v20 + 32) &&
					!(v13 & *((uint32_t*)v20 + 3))) {
					++v19;
				}
				v21 = *((uint32_t*)v20 + 7);
				v20 += 24;
			} while (v21);
		}
		if (!v19) {
			if (v17 & 2) {
				if (v17 & 4) {
					if (!(v17 & 8)) {
						v17 |= 8u;
					}
				} else {
					v17 |= 4u;
				}
			} else {
				v17 |= 2u;
			}
			v17 &= 0xFFFEu;
		}
	}
	if (v17 & 2) {
		v22 = 0;
		if (!*getMemU32Ptr(0x587000, 211000)) {
			// nop
		} else {
			v23 = getMemAt(0x587000, 210996);
			do {
				if (v64 & *((uint32_t*)v23 + 2) && v13 & *(uint32_t*)(*(uint32_t*)v23 + 32) &&
					!(v13 & *((uint32_t*)v23 + 3))) {
					++v22;
				}
				v24 = *((uint32_t*)v23 + 7);
				v23 += 24;
			} while (v24);
		}
		if (!v22) {
			if (v17 & 4) {
				if (!(v17 & 8)) {
					v17 |= 8u;
				}
			} else {
				v17 |= 4u;
			}
			v17 &= 0xFFFDu;
		}
	}
	if (v17 & 4) {
		v25 = 0;
		if (!*getMemU32Ptr(0x587000, 209344)) {
			// nop
		} else {
			v26 = getMemAt(0x587000, 209340);
			do {
				if (v64 & *((uint32_t*)v26 + 2) && v13 & *(uint32_t*)(*(uint32_t*)v26 + 32) &&
					!(v13 & *((uint32_t*)v26 + 3)) && *(uint8_t*)(*(uint32_t*)v26 + 36) & 1) {
					++v25;
				}
				v27 = *((uint32_t*)v26 + 7);
				v26 += 24;
			} while (v27);
		}
		if (!v25) {
			v17 &= 0xFFF3u;
		}
	}
	if (!v17) {
		return v14;
	}
	if (!(v17 & 1)) {
		goto LABEL_103;
	}
	v28 = 0;
	if (!*getMemU32Ptr(0x587000, 210856)) {
		goto LABEL_103;
	}
	v29 = getMemAt(0x587000, 210852);
	do {
		if (*((uint32_t*)v29 + 2) & v64 && v13 & *(uint32_t*)(*(uint32_t*)v29 + 32) &&
			!(v13 & *((uint32_t*)v29 + 3))) {
			++v28;
		}
		v30 = *((uint32_t*)v29 + 7);
		v29 += 24;
	} while (v30);
	if (!v28) {
		goto LABEL_103;
	}
	v31 = nox_common_randomInt_415FA0(0, v28 - 1);
	v32 = 0;
	v33 = 0;
	if (!*getMemU32Ptr(0x587000, 210856)) {
		goto LABEL_103;
	}
	v34 = getMemAt(0x587000, 210852);
	while (1) {
		if (!(v64 & *((uint32_t*)v34 + 2)) || !(v13 & *(uint32_t*)(*(uint32_t*)v34 + 32)) ||
			v13 & *((uint32_t*)v34 + 3)) {
			goto LABEL_100;
		}
		if (v32 == v31) {
			break;
		}
		++v32;
	LABEL_100:
		v35 = *((uint32_t*)v34 + 7);
		v34 += 24;
		++v33;
		if (!v35) {
			goto LABEL_103;
		}
	}
	*(uint32_t*)v67 = *getMemU32Ptr(0x587000, 210852 + 24 * v33);
LABEL_103:
	if (v17 & 2) {
		v36 = 0;
		if (*getMemU32Ptr(0x587000, 211000)) {
			v37 = getMemAt(0x587000, 210996);
			do {
				if (v64 & *((uint32_t*)v37 + 2) && v13 & *(uint32_t*)(*(uint32_t*)v37 + 32) &&
					!(v13 & *((uint32_t*)v37 + 3))) {
					++v36;
				}
				v38 = *((uint32_t*)v37 + 7);
				v37 += 24;
			} while (v38);
			if (v36) {
				v39 = nox_common_randomInt_415FA0(0, v36 - 1);
				v40 = 0;
				v41 = 0;
				if (*getMemU32Ptr(0x587000, 211000)) {
					v42 = getMemAt(0x587000, 210996);
					do {
						if (v64 & *((uint32_t*)v42 + 2) && v13 & *(uint32_t*)(*(uint32_t*)v42 + 32) &&
							!(v13 & *((uint32_t*)v42 + 3))) {
							if (v40 == v39) {
								*(uint32_t*)&v67[4] = *getMemU32Ptr(0x587000, 210996 + 24 * v41);
								break;
							}
							++v40;
						}
						v43 = *((uint32_t*)v42 + 7);
						v42 += 24;
						++v41;
					} while (v43);
				}
			}
		}
	}
	if (v17 & 4) {
		v44 = 0;
		if (*getMemU32Ptr(0x587000, 209344)) {
			v45 = getMemAt(0x587000, 209340);
			do {
				if (v64 & *((uint32_t*)v45 + 2) && v13 & *(uint32_t*)(*(uint32_t*)v45 + 32) &&
					!(v13 & *((uint32_t*)v45 + 3)) && *(uint8_t*)(*(uint32_t*)v45 + 36) & 1) {
					++v44;
				}
				v46 = *((uint32_t*)v45 + 7);
				v45 += 24;
			} while (v46);
			if (v44) {
				v47 = nox_common_randomInt_415FA0(0, v44 - 1);
				v48 = 0;
				v49 = 0;
				if (*getMemU32Ptr(0x587000, 209344)) {
					v50 = getMemAt(0x587000, 209340);
					do {
						if (v64 & *((uint32_t*)v50 + 2) && v13 & *(uint32_t*)(*(uint32_t*)v50 + 32) &&
							!(v13 & *((uint32_t*)v50 + 3)) && *(uint8_t*)(*(uint32_t*)v50 + 36) & 1) {
							if (v48 == v47) {
								v65 = v49;
								*(uint32_t*)&v67[8] = *getMemU32Ptr(0x587000, 209340 + 24 * v49);
								break;
							}
							++v48;
						}
						v51 = *((uint32_t*)v50 + 7);
						v50 += 24;
						++v49;
					} while (v51);
				}
			}
		}
	}
	if (v17 & 8) {
		v52 = a2 >> 1;
		if ((int)(a2 >> 1) >= 1) {
			if (v52 >= 5) {
				v52 = 4;
			}
		} else {
			v52 = 1;
		}
		v53 = v52 - 1;
		v54 = v52 + 1;
		if (v53 < 1) {
			v53 = 1;
		}
		if (v54 >= 5) {
			v54 = 4;
		}
		v55 = nox_common_randomInt_415FA0(v53, v54);
		v56 = 0;
		if (*getMemU32Ptr(0x587000, 209344)) {
			v57 = getMemAt(0x587000, 209340);
			do {
				if (v55 & *((uint32_t*)v57 + 2) && v13 & *(uint32_t*)(*(uint32_t*)v57 + 32) &&
					!(v13 & *((uint32_t*)v57 + 3)) && *(uint8_t*)(*(uint32_t*)v57 + 36) & 2) {
					++v56;
				}
				v58 = *((uint32_t*)v57 + 7);
				v57 += 24;
			} while (v58);
			if (v56) {
				v59 = nox_common_randomInt_415FA0(0, v56 - 1);
				v60 = 0;
				v61 = 0;
				if (*getMemU32Ptr(0x587000, 209344)) {
					v62 = getMemAt(0x587000, 209340);
					do {
						if (v55 & *((uint32_t*)v62 + 2) && v13 & *(uint32_t*)(*(uint32_t*)v62 + 32) &&
							!(v13 & *((uint32_t*)v62 + 3)) && *(uint8_t*)(*(uint32_t*)v62 + 36) & 2) {
							if (v60 == v59) {
								if (getMemByte(0x587000, 209336 + 24 * v65) !=
									getMemByte(0x587000, 209336 + 24 * v61)) {
									*(uint32_t*)&v67[12] = *getMemU32Ptr(0x587000, 209340 + 24 * v61);
								}
								break;
							}
							++v60;
						}
						v63 = *((uint32_t*)v62 + 7);
						v62 += 24;
						++v61;
					} while (v63);
				}
			}
		}
	}
	nox_xxx_modifSetItemAttrs_4E4990(v66, (int*)v67);
	return (uint32_t*)v66;
}

//----- (004F14E0) --------------------------------------------------------
int nox_xxx_rewardMakeWeapon_4F14E0(int a1, unsigned int a2) {
	int v2;                // ebx
	int v3;                // edi
	unsigned char* v4;     // esi
	int v5;                // eax
	int v7;                // ebx
	int v8;                // edi
	int v9;                // ebp
	unsigned char* v10;    // esi
	int v11;               // eax
	char* v12;             // esi
	uint32_t* v13;         // eax
	int v14;               // edi
	int v15;               // eax
	int v16;               // eax
	int v17;               // eax
	int v18;               // eax
	short v19;             // bx
	int v20;               // eax
	int v21;               // ecx
	int v22;               // ebp
	unsigned char* v23;    // eax
	int v24;               // edx
	int v25;               // ecx
	unsigned char* v26;    // eax
	int v27;               // edx
	int v28;               // edx
	unsigned char* v29;    // eax
	int v30;               // ecx
	int v31;               // ecx
	unsigned char* v32;    // eax
	int v33;               // edx
	int v34;               // eax
	int v35;               // esi
	int v36;               // edx
	unsigned char* v37;    // ecx
	int v38;               // edi
	int v39;               // ecx
	unsigned char* v40;    // eax
	int v41;               // esi
	int v42;               // eax
	int v43;               // esi
	int v44;               // edx
	unsigned char* v45;    // ecx
	int v46;               // edi
	signed int v47;        // eax
	int v48;               // ecx
	int v49;               // eax
	int v50;               // ebp
	int v51;               // edx
	unsigned char* v52;    // eax
	int v53;               // ecx
	int v54;               // eax
	int v55;               // edi
	int v56;               // esi
	unsigned char* v57;    // ecx
	int v58;               // edx
	signed int v59;        // eax
	int v60;               // ecx
	int v61;               // eax
	int v62;               // ebp
	int v63;               // edx
	unsigned char* v64;    // eax
	int v65;               // ecx
	int v66;               // eax
	int v67;               // edi
	int v68;               // esi
	unsigned char* v69;    // ecx
	int v70;               // edx
	short v71 = 0;         // [esp+10h] [ebp-28h]
	int v72;               // [esp+14h] [ebp-24h]
	int v73;               // [esp+18h] [ebp-20h]
	int v74;               // [esp+1Ch] [ebp-1Ch]
	int v75;               // [esp+20h] [ebp-18h]
	unsigned char v76[20]; // [esp+24h] [ebp-14h]

	v74 = 0;
	v2 = nox_server_rewardGen_pickRandomSlots_4F0B60(a2);
	v3 = 0;
	v72 = v2;
	if (!*getMemU32Ptr(0x587000, 208180)) {
		return 0;
	}
	v4 = getMemAt(0x587000, 208192);
	do {
		if (*(v4 - 4) & 1 && v2 & *(uint32_t*)v4 && nox_xxx_getUnitDefDd10_4E3BA0(*((uint32_t*)v4 - 2))) {
			v3 += *(v4 - 16);
		}
		v5 = *((uint32_t*)v4 + 2);
		v4 += 20;
	} while (v5);
	if (!v3) {
		return 0;
	}
	v7 = nox_common_randomInt_415FA0(0, v3 - 1);
	v8 = 0;
	v9 = 0;
	if (*getMemU32Ptr(0x587000, 208180)) {
		v10 = getMemAt(0x587000, 208192);
		while (1) {
			if (*(v10 - 4) & 1) {
				if (v72 & *(uint32_t*)v10) {
					if (nox_xxx_getUnitDefDd10_4E3BA0(*((uint32_t*)v10 - 2))) {
						v8 += *(v10 - 16);
						if (v7 < v8) {
							v12 = *(char**)getMemAt(0x587000, 208184 + 20 * v9);
							break;
						}
					}
				}
			}
			v11 = *((uint32_t*)v10 + 2);
			v10 += 20;
			++v9;
			if (!v11) {
				v12 = (char*)a2;
				break;
			}
		}
	} else {
		v12 = (char*)a2;
	}
	if (!v12) {
		return 0;
	}
	v73 = nox_xxx_ammoCheck_415880(v12);
	v13 = nox_xxx_newObjectWithTypeInd_4E3450((int)v12);
	v14 = (int)v13;
	v75 = (int)v13;
	if (!v13) {
		return 0;
	}
	v15 = v13[2];
	if (v15 & 0x1000 && *(uint32_t*)(v14 + 12) & 0x47F0000) {
		if (**(uint8_t**)getMemAt(0x587000, 208180 + 20 * v9) == 35) {
			*(uint32_t*)v76 = 0;
			*(uint32_t*)&v76[4] = 0;
			v16 = nox_xxx_modifGetIdByName_413290("Replenishment1");
			*(uint32_t*)&v76[8] = nox_xxx_modifGetDescById_413330(v16);
			*(uint32_t*)&v76[12] = 0;
			*(uint16_t*)&v76[16] = 0;
			*(uint16_t*)&v76[18] = 0;
			nox_xxx_modifSetItemAttrs_4E4990(v14, (int*)v76);
		}
		return v14;
	}
	switch (v72) {
	case 2:
		v17 = nox_common_randomInt_415FA0(0, 1);
		break;
	case 4:
		v17 = nox_common_randomInt_415FA0(0, 2);
		break;
	case 8:
		v17 = nox_common_randomInt_415FA0(1, 3);
		break;
	case 16:
		v17 = nox_common_randomInt_415FA0(2, 4);
		break;
	default:
		return v14;
	}
	if (!v17) {
		return v14;
	}
	*(uint32_t*)v76 = 0;
	*(uint32_t*)&v76[4] = 0;
	*(uint32_t*)&v76[8] = 0;
	*(uint32_t*)&v76[12] = 0;
	*(uint16_t*)&v76[16] = 0;
	*(uint16_t*)&v76[18] = 0;
	switch (v17) {
	case 1:
		v18 = nox_common_randomInt_415FA0(1, 100);
		if (v18 > 20) {
			v19 = (v18 > 50) + 1;
			LOBYTE(v71) = (v18 > 50) + 1;
		} else {
			v19 = 4;
			LOBYTE(v71) = 4;
		}
		break;
	case 2:
		v20 = nox_common_randomInt_415FA0(1, 100);
		if (v20 > 12) {
			v19 = v20 > 25 ? 3 : 6;
			LOBYTE(v71) = v20 > 25 ? 3 : 6;
		} else {
			v19 = 5;
			LOBYTE(v71) = 5;
		}
		break;
	case 3:
		v19 = 7;
		LOBYTE(v71) = 7;
		break;
	case 4:
		v19 = 15;
		LOBYTE(v71) = 15;
		break;
	default:
		v19 = v71;
		break;
	}
	if (!(v19 & 1)) {
		v22 = v73;
		goto LABEL_67;
	}
	v21 = 0;
	if (!*getMemU32Ptr(0x587000, 210712)) {
		v22 = v73;
	} else {
		v22 = v73;
		v23 = getMemAt(0x587000, 210708);
		do {
			if (v72 & *((uint32_t*)v23 + 2) && v73 & *(uint32_t*)(*(uint32_t*)v23 + 28) &&
				!(v73 & *((uint32_t*)v23 + 4))) {
				++v21;
			}
			v24 = *((uint32_t*)v23 + 7);
			v23 += 24;
		} while (v24);
	}
	if (!v21) {
		if (v19 & 2) {
			if (v19 & 4) {
				if (!(v19 & 8)) {
					v19 |= 8u;
				}
				v19 &= 0xFFFEu;
				LOBYTE(v71) = v19;
			} else {
				v19 = (v19 | 4) & 0xFFFE;
				LOBYTE(v71) = v19;
			}
		} else {
			v19 = (v19 | 2) & 0xFFFE;
			LOBYTE(v71) = v19;
		}
	}
LABEL_67:
	if (v19 & 2) {
		v25 = 0;
		if (!*getMemU32Ptr(0x587000, 211000)) {
			// nop
		} else {
			v26 = getMemAt(0x587000, 210996);
			do {
				if (v72 & *((uint32_t*)v26 + 2) && v22 & *(uint32_t*)(*(uint32_t*)v26 + 28) &&
					!(v22 & *((uint32_t*)v26 + 4))) {
					++v25;
				}
				v27 = *((uint32_t*)v26 + 7);
				v26 += 24;
			} while (v27);
		}
		if (!v25) {
			if (v19 & 4) {
				if (!(v19 & 8)) {
					v19 |= 8u;
				}
			} else {
				v19 |= 4u;
			}
			v19 &= 0xFFFDu;
			LOBYTE(v71) = v19;
		}
	}
	if (v19 & 4) {
		v28 = 0;
		if (!*getMemU32Ptr(0x587000, 209344)) {
			// nop
		} else {
			v29 = getMemAt(0x587000, 209340);
			do {
				if (v72 & *((uint32_t*)v29 + 2) && v22 & *(uint32_t*)(*(uint32_t*)v29 + 28) &&
					!(v22 & *((uint32_t*)v29 + 4)) && *(uint8_t*)(*(uint32_t*)v29 + 36) & 1) {
					++v28;
				}
				v30 = *((uint32_t*)v29 + 7);
				v29 += 24;
			} while (v30);
		}
		if (!v28) {
			v19 &= 0xFFF3u;
			LOBYTE(v71) = v19;
		}
	}
	if (!v19) {
		return v14;
	}
	if (v19 & 1) {
		v31 = 0;
		if (*getMemU32Ptr(0x587000, 210712)) {
			v32 = getMemAt(0x587000, 210708);
			do {
				if (*((uint32_t*)v32 + 2) & v72 && v22 & *(uint32_t*)(*(uint32_t*)v32 + 28) &&
					!(v22 & *((uint32_t*)v32 + 4))) {
					++v31;
				}
				v33 = *((uint32_t*)v32 + 7);
				v32 += 24;
			} while (v33);
			if (v31) {
				v34 = nox_common_randomInt_415FA0(0, v31 - 1);
				v35 = 0;
				v36 = 0;
				if (*getMemU32Ptr(0x587000, 210712)) {
					v37 = getMemAt(0x587000, 210708);
					do {
						if (v72 & *((uint32_t*)v37 + 2) && v22 & *(uint32_t*)(*(uint32_t*)v37 + 28) &&
							!(v22 & *((uint32_t*)v37 + 4))) {
							if (v35 == v34) {
								*(uint32_t*)v76 = *getMemU32Ptr(0x587000, 210708 + 24 * v36);
								break;
							}
							++v35;
						}
						v38 = *((uint32_t*)v37 + 7);
						v37 += 24;
						++v36;
					} while (v38);
				}
			}
		}
	}
	if (v19 & 2) {
		v39 = 0;
		if (*getMemU32Ptr(0x587000, 211000)) {
			v40 = getMemAt(0x587000, 210996);
			do {
				if (v72 & *((uint32_t*)v40 + 2) && v22 & *(uint32_t*)(*(uint32_t*)v40 + 28) &&
					!(v22 & *((uint32_t*)v40 + 4))) {
					++v39;
				}
				v41 = *((uint32_t*)v40 + 7);
				v40 += 24;
			} while (v41);
			if (v39) {
				v42 = nox_common_randomInt_415FA0(0, v39 - 1);
				v43 = 0;
				v44 = 0;
				if (*getMemU32Ptr(0x587000, 211000)) {
					v45 = getMemAt(0x587000, 210996);
					do {
						if (v72 & *((uint32_t*)v45 + 2) && v22 & *(uint32_t*)(*(uint32_t*)v45 + 28) &&
							!(v22 & *((uint32_t*)v45 + 4))) {
							if (v43 == v42) {
								*(uint32_t*)&v76[4] = *getMemU32Ptr(0x587000, 210996 + 24 * v44);
								break;
							}
							++v43;
						}
						v46 = *((uint32_t*)v45 + 7);
						v45 += 24;
						++v44;
					} while (v46);
				}
			}
		}
	}
	if (v19 & 4) {
		v47 = a2 >> 1;
		if ((int)(a2 >> 1) >= 1) {
			if (v47 >= 5) {
				v47 = 4;
			}
		} else {
			v47 = 1;
		}
		v48 = v47 - 1;
		v49 = v47 + 1;
		if (v48 < 1) {
			v48 = 1;
		}
		if (v49 >= 5) {
			v49 = 4;
		}
		v50 = nox_common_randomInt_415FA0(v48, v49);
		v51 = 0;
		if (*getMemU32Ptr(0x587000, 209344)) {
			v52 = getMemAt(0x587000, 209340);
			do {
				if (v50 & *((uint32_t*)v52 + 2) && v73 & *(uint32_t*)(*(uint32_t*)v52 + 28) &&
					!(v73 & *((uint32_t*)v52 + 4)) && *(uint8_t*)(*(uint32_t*)v52 + 36) & 1) {
					++v51;
				}
				v53 = *((uint32_t*)v52 + 7);
				v52 += 24;
			} while (v53);
			if (v51) {
				v54 = nox_common_randomInt_415FA0(0, v51 - 1);
				v55 = 0;
				v56 = 0;
				if (*getMemU32Ptr(0x587000, 209344)) {
					v57 = getMemAt(0x587000, 209340);
					do {
						if (v50 & *((uint32_t*)v57 + 2)) {
							if (v73 & *(uint32_t*)(*(uint32_t*)v57 + 28) && !(v73 & *((uint32_t*)v57 + 4)) &&
								*(uint8_t*)(*(uint32_t*)v57 + 36) & 1) {
								if (v55 == v54) {
									LOBYTE(v19) = v71;
									v74 = v56;
									*(uint32_t*)&v76[8] = *getMemU32Ptr(0x587000, 209340 + 24 * v56);
									break;
								}
								++v55;
							}
							LOBYTE(v19) = v71;
						}
						v58 = *((uint32_t*)v57 + 7);
						v57 += 24;
						++v56;
					} while (v58);
				}
			}
		}
	}
	if (v19 & 8) {
		v59 = a2 >> 1;
		if ((int)(a2 >> 1) >= 1) {
			if (v59 >= 5) {
				v59 = 4;
			}
		} else {
			v59 = 1;
		}
		v60 = v59 - 1;
		v61 = v59 + 1;
		if (v60 < 1) {
			v60 = 1;
		}
		if (v61 >= 5) {
			v61 = 4;
		}
		v62 = nox_common_randomInt_415FA0(v60, v61);
		v63 = 0;
		if (*getMemU32Ptr(0x587000, 209344)) {
			v64 = getMemAt(0x587000, 209340);
			do {
				if (v62 & *((uint32_t*)v64 + 2) && v73 & *(uint32_t*)(*(uint32_t*)v64 + 28) &&
					!(v73 & *((uint32_t*)v64 + 4)) && *(uint8_t*)(*(uint32_t*)v64 + 36) & 2) {
					++v63;
				}
				v65 = *((uint32_t*)v64 + 7);
				v64 += 24;
			} while (v65);
			if (v63) {
				v66 = nox_common_randomInt_415FA0(0, v63 - 1);
				v67 = 0;
				v68 = 0;
				if (*getMemU32Ptr(0x587000, 209344)) {
					v69 = getMemAt(0x587000, 209340);
					do {
						if (v62 & *((uint32_t*)v69 + 2) && v73 & *(uint32_t*)(*(uint32_t*)v69 + 28) &&
							!(v73 & *((uint32_t*)v69 + 4)) && *(uint8_t*)(*(uint32_t*)v69 + 36) & 2) {
							if (v67 == v66) {
								if (getMemByte(0x587000, 209336 + 24 * v74) !=
									getMemByte(0x587000, 209336 + 24 * v68)) {
									*(uint32_t*)&v76[12] = *getMemU32Ptr(0x587000, 209340 + 24 * v68);
								}
								break;
							}
							++v67;
						}
						v70 = *((uint32_t*)v69 + 7);
						v69 += 24;
						++v68;
					} while (v70);
				}
			}
		}
	}
	nox_xxx_modifSetItemAttrs_4E4990(v75, (int*)v76);
	return v75;
}
// 4F1776: variable 'v71' is possibly undefined

//----- (004F1C40) --------------------------------------------------------
uint32_t* nox_xxx_rewardMakePotion_4F1C40(int a1, unsigned int a2) {
	int v2;            // ebx
	int v3;            // edi
	unsigned char* v4; // esi
	int v5;            // eax
	uint32_t* result;  // eax
	int v7;            // ebx
	int v8;            // ebp
	int v9;            // edi
	unsigned char* i;  // esi
	int v11;           // eax
	int v12;           // [esp+18h] [ebp+8h]

	v2 = nox_server_rewardGen_pickRandomSlots_4F0B60(a2);
	v3 = 0;
	v12 = v2;
	if (!*getMemU32Ptr(0x587000, 208180)) {
		return 0;
	}
	v4 = getMemAt(0x587000, 208192);
	do {
		if (*(v4 - 4) & 4 && v2 & *(uint32_t*)v4 && nox_xxx_getUnitDefDd10_4E3BA0(*((uint32_t*)v4 - 2))) {
			v3 += *(v4 - 16);
		}
		v5 = *((uint32_t*)v4 + 2);
		v4 += 20;
	} while (v5);
	if (!v3) {
		return 0;
	}
	v7 = nox_common_randomInt_415FA0(0, v3 - 1);
	v8 = 0;
	v9 = 0;
	if (!*getMemU32Ptr(0x587000, 208180)) {
		return 0;
	}
	for (i = getMemAt(0x587000, 208192);; i += 20) {
		if (*(i - 4) & 4) {
			if (*(uint32_t*)i & v12) {
				if (nox_xxx_getUnitDefDd10_4E3BA0(*((uint32_t*)i - 2))) {
					v8 += *(i - 16);
					if (v7 < v8) {
						break;
					}
				}
			}
		}
		v11 = *((uint32_t*)i + 2);
		++v9;
		if (!v11) {
			return 0;
		}
	}
	if (*getMemU32Ptr(0x587000, 208184 + 20 * v9)) {
		return nox_xxx_newObjectWithTypeInd_4E3450(*getMemU32Ptr(0x587000, 208184 + 20 * v9));
	} else {
		return 0;
	}
}

//----- (004F1D30) --------------------------------------------------------
uint32_t* nox_xxx_createGem_4F1D30(int a1, unsigned int a2) {
	unsigned int v2;  // ebx
	int v3;           // eax
	uint32_t* result; // eax
	uint32_t* v5;     // edi
	int* v6;          // esi

	v2 = nox_server_rewardGen_pickRandomSlots_4F0B60(a2);
	if (v2 < 4 || nox_common_randomInt_415FA0(1, 100) <= 90) {
		if (nox_common_randomInt_415FA0(1, 2) == 1) {
			result = nox_xxx_newObjectByTypeID_4E3810("QuestGoldChest");
		} else {
			result = nox_xxx_newObjectByTypeID_4E3810("QuestGoldPile");
		}
		v5 = result;
		if (result) {
			v6 = (int*)result[173];
			switch (v2) {
			case 2u:
				*v6 = nox_common_randomInt_415FA0(*getMemIntPtr(0x587000, 211144), *getMemIntPtr(0x587000, 211148));
				result = v5;
				break;
			case 4u:
				*v6 = nox_common_randomInt_415FA0(*getMemIntPtr(0x587000, 211152), *getMemIntPtr(0x587000, 211156));
				result = v5;
				break;
			case 8u:
				*v6 = nox_common_randomInt_415FA0(*getMemIntPtr(0x587000, 211160), *getMemIntPtr(0x587000, 211164));
				result = v5;
				break;
			case 0x10u:
				*v6 = nox_common_randomInt_415FA0(*getMemIntPtr(0x587000, 211168), *getMemIntPtr(0x587000, 211172));
				result = v5;
				break;
			default:
				*v6 = nox_common_randomInt_415FA0(*getMemIntPtr(0x587000, 211136), *getMemIntPtr(0x587000, 211140));
				result = v5;
				break;
			}
		}
	} else {
		v3 = nox_common_randomInt_415FA0(1, 100);
		if (v3 >= 50) {
			if (v3 >= 90) {
				result = nox_xxx_newObjectByTypeID_4E3810("DiamondGem");
			} else {
				result = nox_xxx_newObjectByTypeID_4E3810("EmeraldGem");
			}
		} else {
			result = nox_xxx_newObjectByTypeID_4E3810("RubyGem");
		}
	}
	return result;
}

//----- (004F1F00) --------------------------------------------------------
uint32_t* nox_xxx_createGem2_4F1F00(int a1, unsigned int a2) { return nox_xxx_createGem_4F1D30(a1, a2); }

//----- (004F2110) --------------------------------------------------------
void sub_4F2110() {
	int v0;       // esi
	int i;        // eax
	int v2;       // ecx
	int v3;       // ebx
	int v4;       // edi
	int j;        // esi
	int v6;       // eax
	uint32_t* v7; // eax

	v0 = 0;
	if (!dword_5d4594_1568280) {
		dword_5d4594_1568280 = nox_xxx_getNameId_4E3AA0("RewardMarker");
		*getMemU32Ptr(0x5D4594, 1568284) = nox_xxx_getNameId_4E3AA0("RewardMarkerPlus");
	}
	for (i = nox_server_getFirstObject_4DA790(); i; i = nox_server_getNextObject_4DA7A0(i)) {
		v2 = *(unsigned short*)(i + 4);
		if (((unsigned short)v2 == dword_5d4594_1568280 || v2 == *getMemU32Ptr(0x5D4594, 1568284)) &&
			(**(uint8_t**)(i + 692) & 0x80)) {
			++v0;
		}
	}
	v3 = nox_common_randomInt_415FA0(0, v0 - 1);
	v4 = 0;
	for (j = nox_server_getFirstObject_4DA790(); j; j = nox_server_getNextObject_4DA7A0(j)) {
		v6 = *(unsigned short*)(j + 4);
		if (((unsigned short)v6 == dword_5d4594_1568280 || v6 == *getMemU32Ptr(0x5D4594, 1568284)) &&
			(**(uint8_t**)(j + 692) & 0x80)) {
			if (v4 == v3) {
				v7 = nox_xxx_newObjectByTypeID_4E3810("Ankh");
				if (v7) {
					nox_xxx_createAt_4DAA50((int)v7, 0, *(float*)(j + 56), *(float*)(j + 60));
					nox_xxx_delayedDeleteObject_4E5CC0(j);
					return;
				}
			} else {
				++v4;
			}
		}
	}
}

//----- (004F2210) --------------------------------------------------------
int sub_4F2210() {
	uint32_t* v0;     // ebp
	int v1;           // esi
	int v2;           // edi
	int v3;           // edi
	int v4;           // esi
	int result;       // eax
	int v6;           // ecx
	unsigned int v7;  // edi
	uint32_t* v8;     // esi
	uint32_t* v9;     // ebx
	int v10;          // ecx
	int v11;          // edx
	char v12;         // cl
	int v13;          // esi
	int v14;          // ebx
	uint32_t* v15;    // edi
	int v16;          // eax
	int v17;          // ecx
	int* v18;         // ecx
	int v19;          // edx
	int v20;          // eax
	unsigned int v21; // ebx
	int v22;          // esi
	uint32_t* v23;    // edi
	int v24;          // eax
	int v25;          // ecx
	unsigned int v26; // esi
	int* v27;         // edi
	float v28 = 0;    // [esp+10h] [ebp-14h]
	uint32_t* lpMem;  // [esp+14h] [ebp-10h]
	unsigned int v30; // [esp+18h] [ebp-Ch]
	unsigned int v31; // [esp+18h] [ebp-Ch]

	v0 = 0;
	lpMem = 0;
	v1 = nox_game_getQuestStage_4E3CC0();
	v2 = nox_xxx_player_4E3CE0();
	if (!dword_5d4594_1568288) {
		dword_5d4594_1568288 = nox_xxx_getNameId_4E3AA0("RewardMarker");
		*getMemU32Ptr(0x5D4594, 1568292) = nox_xxx_getNameId_4E3AA0("RewardMarkerPlus");
		*getMemU32Ptr(0x5D4594, 1568296) = nox_xxx_getNameId_4E3AA0("RedPotion");
	}
	if (v1 == 1) {
		v28 = 0.5;
	} else {
		switch (v2) {
		case 1:
		case 2:
			v28 = 0.40000001;
			break;
		case 3:
		case 4:
			v28 = 0.69999999;
			break;
		case 5:
		case 6:
			v28 = 1.0;
			break;
		default:
			break;
		}
	}
	v3 = 0;
	v4 = 0;
	result = nox_server_getFirstObject_4DA790();
	if (result) {
		do {
			v6 = *(unsigned short*)(result + 4);
			if ((unsigned short)v6 == dword_5d4594_1568288) {
				if (!(*(uint8_t*)(*(uint32_t*)(result + 692) + 216) & 1)) {
					++v3;
				}
			} else if (v6 == *getMemU32Ptr(0x5D4594, 1568296)) {
				++v4;
			}
			result = nox_server_getNextObject_4DA7A0(result);
		} while (result);
		if (v3) {
			v0 = calloc(v3, 4);
			if (!v4) {
				goto LABEL_21;
			}
		} else if (!v4) {
			return result;
		}
		lpMem = calloc(v4, 4);
	LABEL_21:
		v7 = 0;
		v30 = 0;
		result = nox_server_getFirstObject_4DA790();
		if (result) {
			v8 = lpMem;
			v9 = v0;
			do {
				v10 = *(unsigned short*)(result + 4);
				if ((unsigned short)v10 == dword_5d4594_1568288) {
					v11 = *(uint32_t*)(result + 692);
					v12 = *(uint8_t*)(v11 + 216);
					if (v12 & 1) {
						*(uint8_t*)(v11 + 216) = v12 | 0x80;
					} else {
						*v9 = result;
						++v7;
						++v9;
					}
				} else if (v10 == *getMemU32Ptr(0x5D4594, 1568292)) {
					*(uint8_t*)(*(uint32_t*)(result + 692) + 216) |= 0x80u;
				} else if (v10 == *getMemU32Ptr(0x5D4594, 1568296)) {
					*v8 = result;
					++v8;
					++v30;
				}
				result = nox_server_getNextObject_4DA7A0(result);
			} while (result);
			if (v7) {
				v13 = v7 - 1;
				v14 = (long long)((double)v7 * v28 + 0.5);
				if (v7 != 1) {
					v15 = &v0[v13];
					do {
						v16 = nox_common_randomInt_415FA0(0, v13);
						v17 = v0[v16];
						v0[v16] = *v15;
						*v15 = v17;
						--v13;
						--v15;
					} while (v13);
				}
				if (v14) {
					v18 = v0;
					v19 = v14;
					do {
						v20 = *v18;
						++v18;
						--v19;
						*(uint8_t*)(*(uint32_t*)(v20 + 692) + 216) |= 0x80u;
					} while (v19);
				}
				free(v0);
			}
			v21 = v30;
			if (v30) {
				v22 = v30 - 1;
				v31 = (long long)((double)v30 * v28 + 0.5);
				if (v21 != 1) {
					v23 = &lpMem[v22];
					do {
						v24 = nox_common_randomInt_415FA0(0, v22);
						v25 = lpMem[v24];
						lpMem[v24] = *v23;
						*v23 = v25;
						--v22;
						--v23;
					} while (v22);
				}
				v26 = 0;
				if (v21) {
					v27 = lpMem;
					do {
						if (v26 >= v31) {
							nox_xxx_delayedDeleteObject_4E5CC0(*v27);
						}
						++v26;
						++v27;
					} while (v26 < v21);
				}
				free(lpMem);
			}
		}
	}
	return result;
}
// 4F23B8: variable 'v28' is possibly undefined

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
