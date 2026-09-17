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
extern uint32_t dword_5d4594_1567988;
extern uint32_t dword_5d4594_1565628;
extern uint32_t dword_5d4594_1565632;
extern uint32_t dword_5d4594_1565520;
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

