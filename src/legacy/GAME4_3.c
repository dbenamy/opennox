#include <math.h>
#include <stdio.h>

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
#include "common__random.h"
#include "server__ability__ability.h"
#include "server__magic__plyrgide.h"
#include "server__magic__plyrspel.h"
#include "server__object__health.h"

#include "common__gamemech__pausefx.h"

#include "MixPatch.h"
#include "client__video__draw_common.h"
#include "common/fs/nox_fs.h"
#include "common__magic__speltree.h"
#include "defs.h"
#include "operators.h"
#include "server__script__builtin.h"
#include "server__script__script.h"

extern uint32_t dword_5d4594_2488620;
extern uint32_t dword_5d4594_2488656;
extern uint32_t dword_5d4594_3835360;
extern uint32_t dword_5d4594_2488728;
extern uint32_t dword_5d4594_2489436;
extern uint32_t dword_5d4594_2488724;
extern uint32_t dword_5d4594_2489160;
extern uint32_t dword_5d4594_2488720;
extern uint32_t dword_5d4594_2487932;
extern uint32_t nox_xxx_lightningOwner_5d4594_2487900;
extern uint32_t dword_587000_261388;
extern uint32_t dword_5d4594_2487948;
extern uint32_t dword_5d4594_2488652;
extern uint32_t dword_5d4594_3835348;
extern uint32_t dword_5d4594_3835352;
extern uint32_t nox_xxx_lightningClosestTargetDistance_5d4594_2487912;
extern uint32_t dword_5d4594_3835356;
extern uint32_t dword_5d4594_2487248;
extern uint32_t dword_5d4594_2488660;
extern uint64_t qword_581450_10176;
extern uint64_t qword_581450_9512;
extern uint64_t qword_581450_9544;
extern uint32_t nox_xxx_lightningTargetArrayIndex_5d4594_2487904;
extern uint32_t nox_xxx_lightningTarget_5d4594_2487908;
extern uint32_t dword_5d4594_251572;
extern uint32_t dword_5d4594_2650652;
extern obj_5D4594_2650668_t** ptr_5D4594_2650668;
extern unsigned int gameex_flags;

void nox_xxx_lightningSpellDuration_52FFD0(int a1, int a2, int a3);
int nox_xxx_waypoint_579F00(float2* a1, nox_object_t* a2);
void sub_4FF310(nox_object_t* a1);
char nox_xxx_monsterPlayHurtSound_532800(nox_object_t* a1p) {
	int a1 = a1p;
	int v1; // eax
	int v2; // edi

	LOBYTE(v1) = *(uint8_t*)(a1 + 8);
	v2 = *(uint32_t*)(a1 + 748);
	if (v1 & 2) {
		LOBYTE(v1) = (unsigned char)gameFrame();
		if (gameFrame() >= *(int*)(v2 + 532)) {
			*(uint32_t*)(v2 + 532) =
				gameFrame() + nox_common_randomInt_415FA0(2 * gameFPS(), 4 * gameFPS());
			v1 = nox_xxx_monsterGetSoundSet_424300(a1);
			if (v1) {
				nox_xxx_aud_501960(*(uint32_t*)(v1 + 8), a1, 0, 0);
			}
		}
	}
	return v1;
}

int sub_532880(int a1) {
	int result; // eax

	result = *(uint32_t*)(a1 + 748);
	if (!*(uint32_t*)(result + 520)) {
		*(uint32_t*)(result + 520) = gameFrame();
	}
	return result;
}

int sub_534020(int a1) { return (*(uint32_t*)(a1 + 12) >> 10) & 1; }

void nox_ai_debug_print(char* str);
void nox_ai_debug_printf_5341A0(char* a1, ...) {
	va_list va; // [esp+8h] [ebp+8h]

	va_start(va, a1);
	if (nox_common_getEngineFlag(NOX_ENGINE_FLAG_ENABLE_SHOW_AI)) {
		nox_vsprintf((char*)getMemAt(0x5D4594, 2487996), a1, va);
		nox_ai_debug_print((char*)getMemAt(0x5D4594, 2487996));
	}
}

char* sub_5345B0(int a1) {
	unsigned char* v1; // ecx
	int v2;            // eax

	v1 = getMemAt(0x587000, 262056);
	while (1) {
		if (*(uint32_t*)v1 == a1) {
			v2 = 0;
			while (1) {
				if (v2 == a1) {
					return *(char**)getMemAt(0x587000, 261768 + 4 * v2);
				}
				if (++v2 >= 39) {
					goto LABEL_6;
				}
			}
		}
	LABEL_6:
		v1 += 4;
		if ((int)v1 >= (int)getMemAt(0x587000, 262072)) {
			return *(char**)getMemAt(0x587000, 261920);
		}
	}
}

int nox_xxx_actionNByNameMB_5345F0(const char* a1) {
	int v1;          // ebp
	const char** v2; // edi

	v1 = 0;
	v2 = (const char**)getMemAt(0x587000, 261768);
	while (strcmp(*v2, a1)) {
		++v2;
		++v1;
		if ((int)v2 >= (int)getMemAt(0x587000, 261924)) {
			return 38;
		}
	}
	return v1;
}

char* sub_534650(int a1) {
	int v1; // eax

	v1 = 0;
	while (v1 != a1) {
		if (++v1 >= 72) {
			return 0;
		}
	}
	return *(char**)getMemAt(0x587000, 261768 + 4 * v1);
}

int nox_xxx_actionByName_534670(const char* a1) {
	int v1;          // ebp
	const char** v2; // edi

	v1 = 0;
	v2 = (const char**)getMemAt(0x587000, 261768);
	while (strcmp(*v2, a1)) {
		++v2;
		++v1;
		if ((int)v2 >= (int)getMemAt(0x587000, 262056)) {
			return 0;
		}
	}
	return v1;
}

char* sub_5361B0(char* a1, int a2) {
	char* result; // eax
	char v3;      // al
	double v4;    // [esp+4h] [ebp-8h]

	*(uint32_t*)a2 = 1;
	result = strtok(a1, " ");
	if (result) {
		v3 = atoi(result);
		*(uint8_t*)(a2 + 109) = v3;
		*(uint8_t*)(a2 + 108) = v3;
		*(uint32_t*)(a2 + 112) = 100;
		result = strtok(0, " ");
		if (result) {
			v4 = (double)gameFPS();
			*(uint32_t*)(a2 + 100) = (long long)(v4 / atof(result));
			result = strtok(0, " ");
			if (result) {
				*(uint32_t*)(a2 + 92) = nox_xxx_spellNameToN_4243F0(result);
				result = (char*)1;
			}
		}
	}
	return result;
}

char* sub_536260(char* a1, int a2) {
	char* result; // eax
	char v3;      // al
	char* v4;     // eax
	int v5;       // eax
	char* v6;     // eax
	double v7;    // [esp+Ch] [ebp-8h]
	double v8;

	*(uint32_t*)a2 = 0;
	result = strtok(a1, " ");
	if (result) {
		v3 = atoi(result);
		*(uint8_t*)(a2 + 109) = v3;
		*(uint8_t*)(a2 + 108) = v3;
		*(uint32_t*)(a2 + 112) = 100;
		result = strtok(0, " ");
		if (result) {
			strcpy((char*)(a2 + 4), result);
			*(uint32_t*)(a2 + 84) = 0;
			result = strtok(0, " ");
			if (result) {
				v7 = (double)gameFPS();
				v8 = atof(result);
				if (v8 == 0.0) {
					*(uint32_t*)(a2 + 100) = 0;
				} else {
					*(uint32_t*)(a2 + 100) = (long long)(v7 / v8);
				}
				v4 = strtok(0, " ");
				if (v4 && !strcmp(v4, "MULTI_SHOT")) {
					v5 = *(uint32_t*)(a2 + 96);
					LOBYTE(v5) = v5 | 1;
					*(uint32_t*)(a2 + 96) = v5;
				}
				v6 = strtok(0, " ");
				if (v6) {
					*(uint32_t*)(a2 + 88) = nox_xxx_utilFindSound_40AF50(v6);
				}
				result = (char*)1;
			}
		}
	}
	return result;
}

int sub_5364E0(char* a1, int a2) {
	unsigned int v2; // ecx
	char v3;         // al
	char* v4;        // edi
	char* v5;        // esi
	int result;      // eax
	char v7[256];    // [esp+Ch] [ebp-100h]

	sscanf(a1, "%s", v7);
	v2 = strlen(v7) + 1;
	v3 = v2;
	v2 >>= 2;
	memcpy((void*)(a2 + 16), v7, 4 * v2);
	v5 = &v7[4 * v2];
	v4 = (char*)(a2 + 16 + 4 * v2);
	LOBYTE(v2) = v3;
	result = 1;
	memcpy(v4, v5, v2 & 3);
	*(uint32_t*)(a2 + 12) = 0;
	return result;
}

int sub_536550(char* a1, uint32_t* a2) {
	sscanf(a1, "%f %f", a2, a2 + 2);
	a2[1] = *a2;
	return 1;
}

int sub_536580(char* a1, int a2) {
	sscanf(a1, "%d %d %d", a2, a2 + 4, a2 + 8);
	return 1;
}

int sub_5365B0(char* a1, int a2) {
	char* v2; // eax
	char* v3; // eax

	v2 = strtok(a1, " ");
	if (v2) {
		*(uint32_t*)(a2 + 36) = nox_xxx_utilFindSound_40AF50(v2);
	}
	v3 = strtok(0, " ");
	if (v3) {
		*(uint32_t*)(a2 + 40) = nox_xxx_utilFindSound_40AF50(v3);
	}
	return 1;
}

int sub_536600(char* a1, int a2) {
	sscanf(a1, "%d", a2);
	return 1;
}

int sub_536B40(char* a1, int a2) {
	char v3[64]; // [esp+4h] [ebp-40h]

	sscanf(a1, "%s %s", a2, v3);
	*(uint32_t*)(a2 + 128) = nox_xxx_utilFindSound_40AF50(v3);
	return 1;
}

int sub_536D80(char* a1, int a2) {
	sscanf(a1, "%d", a2);
	return 1;
}

int sub_536DA0(char* a1, int* a2) {
	int v2;       // eax
	char v4[256]; // [esp+0h] [ebp-100h]

	sscanf(a1, "%s", v4);
	v2 = nox_xxx_utilFindSound_40AF50(v4);
	*a2 = v2;
	return v2 != 0;
}

int sub_536DE0(char* a1, uint8_t* a2) {
	sscanf(a1, "%d", &a1);
	*a2 = (uint8_t)a1;
	return 1;
}

int nox_xxx_collideDamageLoad_536E10(char* a1, int a2) {
	char* v2; // eax
	char* v3; // eax
	int v4;   // eax

	v2 = strtok(a1, " ");
	*(uint8_t*)a2 = atoi(v2);
	v3 = strtok(0, " ");
	v4 = nox_xxx_parseDamageTypeByName_4E0A00(v3);
	*(uint32_t*)(a2 + 4) = v4;
	return v4 != 18;
}

int sub_536E50(char* a1, uint8_t* a2) {
	char* v2; // eax

	v2 = strtok(a1, " ");
	*a2 = atoi(v2);
	return 1;
}

int sub_536E80(char* a1, int* a2) {
	char* v2; // eax
	char* v3; // eax

	v2 = strtok(a1, " ");
	*a2 = atoi(v2);
	v3 = strtok(a1, " ");
	a2[1] = atoi(v3);
	return 1;
}

void nox_xxx_updateUndeadKiller_53E190(int a1) {
	int v1; // eax

	v1 = **(uint32_t**)(a1 + 700);
	if (v1 && *(uint8_t*)(v1 + 88) & 1) {
		nox_xxx_delayedDeleteObject_4E5CC0(a1);
	} else if ((unsigned int)(gameFrame() - *(uint32_t*)(a1 + 136)) > 0x46) {
		nox_xxx_delayedDeleteObject_4E5CC0(a1);
	}
}

int nox_xxx_useRead_53F7C0(int a1, int a2) {
	int v2; // esi
	int v3; // ecx

	if (*(uint8_t*)(a1 + 8) & 4) {
		v2 = *(uint32_t*)(a2 + 736);
		v3 = *(uint32_t*)(v2 + 256);
		if ((gameFrame() - v3 > (unsigned int)(3 * gameFPS()) || !v3) &&
			nox_xxx_mapCheck_537110(a1, a2) == 1) {
			nox_xxx_netPriMsgToPlayer_4DA2C0(a1, (const char*)v2, 1);
			*(uint32_t*)(v2 + 256) = gameFrame();
		}
	}
	return 1;
}

int sub_53F830(int a1, int a2) {
	int v2; // esi
	int v3; // ebx
	int v4; // edi
	int v5; // ecx
	int v6; // eax
	int v7; // eax
	int v8; // edx

	v2 = a1;
	if (*(uint8_t*)(a1 + 8) & 4) {
		v3 = *(uint32_t*)(a1 + 748);
		v4 = *(uint32_t*)(a2 + 736);
		v5 = *(uint32_t*)(v4 + 256);
		if ((gameFrame() - v5 > (unsigned int)(3 * gameFPS()) || !v5) &&
			nox_xxx_mapCheck_537110(a1, a2) == 1) {
			if (sub_4D75E0()) {
				v6 = nox_game_getQuestStage_4E3CC0();
				v7 = nox_server_questNextStageThreshold_4D74F0(v6);
				v8 = *(uint32_t*)(v3 + 276);
				a1 = v7;
				nox_xxx_netInformTextMsg_4DA0F0(*(unsigned char*)(v8 + 2064), 21, &a1);
			} else {
				nox_xxx_netPriMsgToPlayer_4DA2C0(v2, "GeneralPrint:WarpClosed", 1);
			}
			*(uint32_t*)(v4 + 256) = gameFrame();
		}
	}
	return 1;
}

int nox_xxx_castPixies_540440(int a1, int a2, int a3, int a4, int a5, int a6) {
	int v6;        // ebx
	int v7;        // eax
	int v8;        // esi
	int v9;        // eax
	int v10;       // ebp
	float v11;     // eax
	double v12;    // st7
	uint32_t* v13; // esi
	int* v14;      // edi
	int v15;       // eax
	float4 v17;    // [esp+Ch] [ebp-10h]
	float v18;     // [esp+2Ch] [ebp+10h]
	int v19;       // [esp+34h] [ebp+18h]

	v6 = a4;
	v7 = *getMemU32Ptr(0x5D4594, 2489140);
	v18 = *(float*)(a4 + 176) + 4.0;
	if (!*getMemU32Ptr(0x5D4594, 2489140)) {
		v7 = nox_xxx_getNameId_4E3AA0("Pixie");
		*getMemU32Ptr(0x5D4594, 2489140) = v7;
	}
	v8 = nox_xxx_unitIsUnitTT_4E7C80(a3, v7);
	if (v8 < (int)(long long)nox_xxx_gamedataGetFloatTable_419D70("PixieCount", a6 - 1)) {
		v9 = (unsigned long long)(long long)nox_xxx_gamedataGetFloatTable_419D70("PixieCount", a6 - 1) - v8;
		if (v9 > 0) {
			v19 = v9;
			do {
				v10 = nox_common_randomInt_415FA0(0, 255);
				v11 = *(float*)(v6 + 60);
				v12 = v18 * *getMemFloatPtr(0x587000, 194136 + 8 * v10) + *(float*)(v6 + 56);
				v17.field_0 = *(float*)(v6 + 56);
				v17.field_4 = v11;
				v17.field_8 = v12;
				v17.field_C = v18 * *getMemFloatPtr(0x587000, 194140 + 8 * v10) + *(float*)(v6 + 60);
				if (nox_xxx_mapTraceRay_535250(&v17, 0, 0, 5)) {
					v13 = nox_xxx_newObjectByTypeID_4E3810("Pixie");
					if (v13) {
						v14 = (int*)v13[187];
						nox_xxx_createAt_4DAA50((int)v13, a3, v17.field_8, v17.field_C);
						*((uint16_t*)v13 + 63) = v10;
						*((uint16_t*)v13 + 62) = v10;
						v13[20] = 0;
						v13[21] = 0;
						v14[1] = nox_xxx_spellFlySearchTarget_540610(0, (int)v13, 32, 600.0, 0, a3);
						*v14 = a3;
						v14[3] = a1;
						v13[39] = *(uint32_t*)(v6 + 56);
						v13[40] = *(uint32_t*)(v6 + 60);
						v14[5] = gameFrame() + gameFPS() * nox_common_randomInt_415FA0(30, 90);
						v14[6] = gameFrame();
					}
				}
				--v19;
			} while (v19);
		}
		v15 = nox_xxx_spellGetAud44_424800(a1, 0);
		nox_xxx_aud_501960(v15, v6, 0, 0);
	}
	return 1;
}

void* nox_objectTypeGetXfer(char* id);
