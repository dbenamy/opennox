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
extern uint32_t dword_5d4594_2488608;
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
extern uint32_t dword_5d4594_2488604;
extern uint32_t nox_xxx_lightningTargetArrayIndex_5d4594_2487904;
extern uint32_t nox_xxx_lightningTarget_5d4594_2487908;
extern uint32_t dword_5d4594_251572;
extern uint32_t dword_5d4594_2650652;
extern obj_5D4594_2650668_t** ptr_5D4594_2650668;
extern unsigned int gameex_flags;

void nox_xxx_lightningSpellDuration_52FFD0(int a1, int a2, int a3);
int nox_xxx_waypoint_579F00(float2* a1, nox_object_t* a2);
void sub_4FF310(nox_object_t* a1);
//----- (00532800) --------------------------------------------------------
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

//----- (00532880) --------------------------------------------------------
int sub_532880(int a1) {
	int result; // eax

	result = *(uint32_t*)(a1 + 748);
	if (!*(uint32_t*)(result + 520)) {
		*(uint32_t*)(result + 520) = gameFrame();
	}
	return result;
}

//----- (00533080) --------------------------------------------------------
int nox_xxx_projAddVelocitySmth_533080(int a1, int a2, float a3, int a4) {
	int result;     // eax
	double v5;      // st7
	double v6;      // st6
	long double v7; // st7

	result = a2;
	v5 = *(float*)(a2 + 56) - *(float*)(a1 + 56);
	v6 = *(float*)(a2 + 60) - *(float*)(a1 + 60);
	v7 = sqrt(v6 * v6 + v5 * v5) / a3;
	*(float*)a4 = v7 * *(float*)(a2 + 80) + *(float*)(a2 + 56);
	*(float*)(a4 + 4) = v7 * *(float*)(a2 + 84) + *(float*)(a2 + 60);
	return result;
}

//----- (00534020) --------------------------------------------------------
int sub_534020(int a1) { return (*(uint32_t*)(a1 + 12) >> 10) & 1; }

//----- (005341A0) --------------------------------------------------------
void nox_ai_debug_print(char* str);
void nox_ai_debug_printf_5341A0(char* a1, ...) {
	va_list va; // [esp+8h] [ebp+8h]

	va_start(va, a1);
	if (nox_common_getEngineFlag(NOX_ENGINE_FLAG_ENABLE_SHOW_AI)) {
		nox_vsprintf((char*)getMemAt(0x5D4594, 2487996), a1, va);
		nox_ai_debug_print((char*)getMemAt(0x5D4594, 2487996));
	}
}

//----- (005345B0) --------------------------------------------------------
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

//----- (005345F0) --------------------------------------------------------
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

//----- (00534650) --------------------------------------------------------
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

//----- (00534670) --------------------------------------------------------
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

//----- (005361B0) --------------------------------------------------------
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

//----- (00536260) --------------------------------------------------------
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

//----- (005364E0) --------------------------------------------------------
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
// 5364E0: using guessed type char var_100[256];

//----- (00536550) --------------------------------------------------------
int sub_536550(char* a1, uint32_t* a2) {
	sscanf(a1, "%f %f", a2, a2 + 2);
	a2[1] = *a2;
	return 1;
}

//----- (00536580) --------------------------------------------------------
int sub_536580(char* a1, int a2) {
	sscanf(a1, "%d %d %d", a2, a2 + 4, a2 + 8);
	return 1;
}

//----- (005365B0) --------------------------------------------------------
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

//----- (00536600) --------------------------------------------------------
int sub_536600(char* a1, int a2) {
	sscanf(a1, "%d", a2);
	return 1;
}

//----- (00536B40) --------------------------------------------------------
int sub_536B40(char* a1, int a2) {
	char v3[64]; // [esp+4h] [ebp-40h]

	sscanf(a1, "%s %s", a2, v3);
	*(uint32_t*)(a2 + 128) = nox_xxx_utilFindSound_40AF50(v3);
	return 1;
}

//----- (00536D80) --------------------------------------------------------
int sub_536D80(char* a1, int a2) {
	sscanf(a1, "%d", a2);
	return 1;
}

//----- (00536DA0) --------------------------------------------------------
int sub_536DA0(char* a1, int* a2) {
	int v2;       // eax
	char v4[256]; // [esp+0h] [ebp-100h]

	sscanf(a1, "%s", v4);
	v2 = nox_xxx_utilFindSound_40AF50(v4);
	*a2 = v2;
	return v2 != 0;
}

//----- (00536DE0) --------------------------------------------------------
int sub_536DE0(char* a1, uint8_t* a2) {
	sscanf(a1, "%d", &a1);
	*a2 = (uint8_t)a1;
	return 1;
}

//----- (00536E10) --------------------------------------------------------
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

//----- (00536E50) --------------------------------------------------------
int sub_536E50(char* a1, uint8_t* a2) {
	char* v2; // eax

	v2 = strtok(a1, " ");
	*a2 = atoi(v2);
	return 1;
}

//----- (00536E80) --------------------------------------------------------
int sub_536E80(char* a1, int* a2) {
	char* v2; // eax
	char* v3; // eax

	v2 = strtok(a1, " ");
	*a2 = atoi(v2);
	v3 = strtok(a1, " ");
	a2[1] = atoi(v3);
	return 1;
}

//----- (005374B0) --------------------------------------------------------
int nox_xxx_traceRay_5374B0(float4* a1) { return nox_xxx_mapTraceRay_535250(a1, 0, 0, 9); }

//----- (00537580) --------------------------------------------------------
int sub_537580(int a1) { return *(uint8_t*)(a1 + 464) & 1; }

//----- (005375A0) --------------------------------------------------------
void sub_5375A0(int a1) {
	int v1;  // eax
	int v2;  // ecx
	char v3; // al

	if (*(uint8_t*)(a1 + 464) & 1) {
		v1 = dword_5d4594_2488604;
		v2 = 0;
		if (dword_5d4594_2488604) {
			while (v1 != a1) {
				v2 = v1;
				v1 = *(uint32_t*)(v1 + 460);
				if (!v1) {
					return;
				}
			}
			if (v1) {
				if (v2) {
					*(uint32_t*)(v2 + 460) = *(uint32_t*)(a1 + 460);
				} else {
					dword_5d4594_2488604 = *(uint32_t*)(a1 + 460);
				}
				if (a1 == dword_5d4594_2488608) {
					dword_5d4594_2488608 = v2;
				}
				v3 = *(uint8_t*)(a1 + 464);
				*(uint32_t*)(a1 + 460) = -1;
				*(uint8_t*)(a1 + 464) = v3 & 0xFE;
			}
		}
	}
}

//----- (00537610) --------------------------------------------------------
void sub_50B500();
char nox_xxx_unitHasCollideOrUpdateFn_537610(nox_object_t* a1p) {
	int a1 = a1p;
	int v1;               // eax
	int v2;               // edx
	int v3;               // edi
	void (*v4)(int, int); // ecx
	int v5;               // ecx

	v1 = *(uint32_t*)(a1 + 744);
	if (v1 || (v1 = *(uint32_t*)(a1 + 696)) != 0 && !(*(uint8_t*)(a1 + 16) & 0x40)) {
		if ((v2 = *(uint32_t*)(a1 + 8), !(v2 & 0x400000)) && !(*(uint8_t*)(a1 + 16) & 8) ||
			(v3 = nox_xxx_getNameId_4E3AA0("Spike"), v1 = nox_xxx_getNameId_4E3AA0("PeriodicSpike"),
			 v2 = *(uint32_t*)(a1 + 8), v2 & 0xE080) ||
			(v4 = *(void (**)(int, int))(a1 + 696), v4 == nox_xxx_collideFist_4EADF0) ||
			v4 == nox_xxx_collideUndeadKiller_4EBD40 || (v5 = *(unsigned short*)(a1 + 4), (unsigned short)v5 == v3) ||
			v5 == v1) {
			if (*(uint8_t*)(a1 + 16) & 4) {
				if (v2 & 0x2008) {
					sub_50B500();
				}
				nullsub_30(a1);
				LOBYTE(v1) = *(uint8_t*)(a1 + 464);
				if (!(v1 & 1)) {
					if (dword_5d4594_2488608) {
						*(uint32_t*)(dword_5d4594_2488608 + 460) = a1;
					} else {
						dword_5d4594_2488604 = a1;
					}
					dword_5d4594_2488608 = a1;
					LOBYTE(v1) = *(uint8_t*)(a1 + 464) | 1;
					*(uint32_t*)(a1 + 460) = 0;
					*(uint8_t*)(a1 + 464) = v1;
				}
			}
		}
	}
	return v1;
}
// 5485F0: using guessed type void  nullsub_30(uint32_t);

//----- (00537700) --------------------------------------------------------
nox_object_t* sub_537700() {
	int result;   // eax
	uint32_t* v1; // ecx

	result = dword_5d4594_2488604;
	v1 = (uint32_t*)(dword_5d4594_2488604 + 460);
	dword_5d4594_2488604 = *(uint32_t*)(dword_5d4594_2488604 + 460);
	if (!dword_5d4594_2488604) {
		dword_5d4594_2488608 = 0;
	}
	*v1 = -1;
	*(uint8_t*)(result + 464) &= 0xFEu;
	return result;
}

//----- (00537740) --------------------------------------------------------
int sub_537740() { return dword_5d4594_2488604; }

//----- (00537750) --------------------------------------------------------
int sub_537750(int a1) {
	int result; // eax

	result = a1;
	if (a1) {
		result = *(uint32_t*)(a1 + 460);
	}
	return result;
}

//----- (00537760) --------------------------------------------------------
unsigned int sub_537760() { return dword_5d4594_2488620 != 0 ? (unsigned int)getMemAt(0x5D4594, 2488612) : 0; }

//----- (00537770) --------------------------------------------------------
void sub_537770(nox_object_t* a1p) {
	int a1 = a1p;
	int v1;    // eax
	int v3;    // [esp+0h] [ebp-Ch]
	float2 v4; // [esp+4h] [ebp-8h]

	LOBYTE(v1) = getMemByte(0x5D4594, 2488624);
	if (!*getMemU32Ptr(0x5D4594, 2488624)) {
		*getMemU32Ptr(0x5D4594, 2488624) = nox_xxx_getNameId_4E3AA0("SmallFist");
		*getMemU32Ptr(0x5D4594, 2488628) = nox_xxx_getNameId_4E3AA0("MediumFist");
		v1 = nox_xxx_getNameId_4E3AA0("LargeFist");
		*getMemU32Ptr(0x5D4594, 2488632) = v1;
	}
	if (!(*(uint8_t*)(a1 + 16) & 0x60)) {
		dword_5d4594_2488620 = 0;
		LOBYTE(v1) = nox_xxx_projectileTraceHit_537850(a1, &v3, &v4);
		if ((uint8_t)v1) {
			if (!v3 || (v1 = *(unsigned short*)(v3 + 4), (unsigned short)v1 != *getMemU32Ptr(0x5D4594, 2488624)) &&
						   v1 != *getMemU32Ptr(0x5D4594, 2488628) && v1 != *getMemU32Ptr(0x5D4594, 2488632)) {
				(*(void (**)(int, int, float2*))(a1 + 696))(a1, v3, &v4);
				LOBYTE(v1) = v3;
				dword_5d4594_2488620 = 0;
				if (v3) {
					v4.field_0 = -v4.field_0;
					v4.field_4 = -v4.field_4;
					LOBYTE(v1) = (*(int (**)(int, int, float2*))(v3 + 696))(v3, a1, &v4);
				}
			}
		}
	}
}

//----- (00537850) --------------------------------------------------------
int sub_57CDB0(int2* a1, float* a2, float2* a3);
char nox_xxx_projectileTraceHit_537850(int a1, int* a2, float2* a3) {
	int v3;     // esi
	int v4;     // ebx
	double v5;  // st7
	double v6;  // st7
	int v7;     // edx
	float v8;   // eax
	int v9;     // eax
	float v10;  // eax
	float v11;  // ecx
	float v12;  // edx
	bool v13;   // al
	int v14;    // edx
	int v15;    // edi
	float v16;  // ecx
	double v17; // st7
	int v18;    // ebp
	double v19; // st7
	double v20; // st6
	double v21; // st7
	double v22; // st6
	double v23; // st5
	float v25;  // edx
	float v26;  // [esp+1Ch] [ebp-44h]
	float2 v27; // [esp+20h] [ebp-40h]
	float2 v28; // [esp+28h] [ebp-38h]
	int v29[2]; // [esp+30h] [ebp-30h]
	float2 a2a; // [esp+38h] [ebp-28h]
	float2 v31; // [esp+40h] [ebp-20h]
	int2 a3a;   // [esp+48h] [ebp-18h]
	float4 a1a; // [esp+50h] [ebp-10h]
	int v34;    // [esp+64h] [ebp+4h]
	int v35;    // [esp+64h] [ebp+4h]

	v3 = a1;
	v4 = 0;
	*(float*)&v34 = *(float*)(a1 + 64) - *(float*)(a1 + 56);
	v5 = *(float*)(v3 + 68) - *(float*)(v3 + 60);
	v26 = v5;
	v6 = v5 * v26 + *(float*)&v34 * *(float*)&v34;
	if (v6 <= 36.0) {
		v7 = *(uint32_t*)(v3 + 56);
		v8 = *(float*)(v3 + 64);
		v28.field_4 = *(float*)(v3 + 68);
		v29[0] = v7;
		v28.field_0 = v8;
		v29[1] = *(uint32_t*)(v3 + 60);
		v9 = sub_54E810(v3, &v28, (int)v29);
		if (v9) {
			v4 = v9;
			v27.field_0 = *(float*)v29 - *(float*)(v9 + 56);
			v27.field_4 = *(float*)&v29[1] - *(float*)(v9 + 60);
		}
	} else {
		v15 = nox_double2int(sqrt(v6 * 0.027777778)) + 1;
		v16 = *(float*)(v3 + 60);
		v17 = (double)v15;
		v29[0] = *(uint32_t*)(v3 + 56);
		v18 = 0;
		*(float*)&v29[1] = v16;
		LODWORD(v28.field_0) = v29[0];
		v28.field_4 = v16;
		*(float*)&v35 = *(float*)&v34 / v17;
		v27.field_0 = v26 / v17;
		if (v15 > 0) {
			while (1) {
				v28.field_0 = v28.field_0 + *(float*)&v35;
				v28.field_4 = v28.field_4 + v27.field_0;
				v9 = sub_54E810(v3, &v28, (int)v29);
				if (v9) {
					v4 = v9;
					v27.field_0 = *(float*)v29 - *(float*)(v9 + 56);
					v27.field_4 = *(float*)&v29[1] - *(float*)(v9 + 60);
					break;
				}
				++v18;
				*(float2*)v29 = v28;
				if (v18 >= v15) {
					break;
				}
			}
		}
	}
	v10 = *(float*)(v3 + 60);
	v11 = *(float*)(v3 + 64);
	a1a.field_0 = *(float*)(v3 + 56);
	v12 = *(float*)(v3 + 68);
	a1a.field_4 = v10;
	a1a.field_8 = v11;
	a1a.field_C = v12;
	if (nox_xxx_mapTraceRay_535250(&a1a, &a2a, &a3a, 5)) {
		v13 = 0;
	} else {
		*(int2*)getMemAt(0x5D4594, 2488612) = a3a;
		*(float2*)&a1a.field_8 = a2a;
		dword_5d4594_2488620 = 1;
		v13 = sub_57CDB0(&a3a, &a1a.field_0, &v31) != 0;
		v14 = *(uint32_t*)(v3 + 60);
		*(uint32_t*)(v3 + 64) = *(uint32_t*)(v3 + 56);
		*(uint32_t*)(v3 + 68) = v14;
	}
	if (v4) {
		if (!v13) {
			*a3 = v27;
			*a2 = v4;
			return 1;
		}
		v19 = *(float*)(v3 + 56) - *(float*)(v4 + 56);
		v20 = *(float*)(v3 + 60) - *(float*)(v4 + 60);
		v21 = v20 * v20 + v19 * v19;
		v22 = *(float*)(v3 + 56) - a2a.field_0;
		v23 = *(float*)(v3 + 60) - a2a.field_4;
		if (v21 < v23 * v23 + v22 * v22) {
			*a3 = v27;
			*a2 = v4;
			return 1;
		}
		v25 = v31.field_4;
		a3->field_0 = v31.field_0;
		a3->field_4 = v25;
		*a2 = 0;
		return 1;
	}
	if (v13) {
		v25 = v31.field_4;
		a3->field_0 = v31.field_0;
		a3->field_4 = v25;
		*a2 = 0;
		return 1;
	}
	return 0;
}
// 537A87: variable 'v24' is possibly undefined

//----- (00537AF0) --------------------------------------------------------
void nox_xxx_sMakeScorch_537AF0(float* a1, int a2) {
	uint32_t* result; // eax
	uint32_t* v3;     // esi
	int v4;           // eax
	int v5;           // [esp-18h] [ebp-18h]
	int v6;           // [esp-18h] [ebp-18h]
	int v7;           // [esp-18h] [ebp-18h]

	if (!*getMemU32Ptr(0x5D4594, 2488636)) {
		nox_xxx_scorchInit_537BD0();
	}
	if (a2) {
		if (a2 == 1) {
			v6 = *getMemU32Ptr(0x587000, 276836 + 8 * nox_common_randomInt_415FA0(0, 0));
			result = nox_xxx_newObjectWithTypeInd_4E3450(v6);
		} else {
			result = (uint32_t*)(a2 - 2);
			if (a2 != 2) {
				return;
			}
			v5 = *getMemU32Ptr(0x587000, 276844 + 8 * nox_common_randomInt_415FA0(0, 0));
			result = nox_xxx_newObjectWithTypeInd_4E3450(v5);
		}
	} else {
		v7 = *getMemU32Ptr(0x587000, 276828 + 8 * nox_common_randomInt_415FA0(0, 0));
		result = nox_xxx_newObjectWithTypeInd_4E3450(v7);
	}
	v3 = result;
	if (result) {
		nox_xxx_createAt_4DAA50((int)result, 0, *(float*)a1, *((float*)a1 + 1));
		if (nox_common_gameFlags_check_40A5C0(4096)) {
			v4 = nox_common_randomInt_415FA0(5, 8);
		} else {
			v4 = nox_common_randomInt_415FA0(10, 20);
		}
		nox_xxx_unitSetDecayTime_511660(v3, gameFPS() * v4);
	}
}

//----- (00537BD0) --------------------------------------------------------
int nox_xxx_scorchInit_537BD0() {
	int result; // eax

	*getMemU32Ptr(0x587000, 276828) = nox_xxx_getNameId_4E3AA0(*(char**)getMemAt(0x587000, 276824));
	*getMemU32Ptr(0x587000, 276836) = nox_xxx_getNameId_4E3AA0(*(char**)getMemAt(0x587000, 276832));
	result = nox_xxx_getNameId_4E3AA0(*(char**)getMemAt(0x587000, 276840));
	*getMemU32Ptr(0x587000, 276844) = result;
	*getMemU32Ptr(0x5D4594, 2488636) = 1;
	return result;
}

//----- (0053E190) --------------------------------------------------------
void nox_xxx_updateUndeadKiller_53E190(int a1) {
	int v1; // eax

	v1 = **(uint32_t**)(a1 + 700);
	if (v1 && *(uint8_t*)(v1 + 88) & 1) {
		nox_xxx_delayedDeleteObject_4E5CC0(a1);
	} else if ((unsigned int)(gameFrame() - *(uint32_t*)(a1 + 136)) > 0x46) {
		nox_xxx_delayedDeleteObject_4E5CC0(a1);
	}
}

//----- (0053F7C0) --------------------------------------------------------
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

//----- (0053F830) --------------------------------------------------------
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

//----- (0053F930) --------------------------------------------------------
int sub_53F930(int a1, int a2) {
	int v2;     // ebx
	int v3;     // esi
	int result; // eax

	if (!(*(uint8_t*)(a1 + 8) & 4)) {
		return 0;
	}
	v2 = *(uint32_t*)(a1 + 748);
	v3 = nox_xxx_guide_427010(*(const char**)(a2 + 736));
	if (nox_common_gameFlags_check_40A5C0(4096) && *(uint8_t*)(*(uint32_t*)(v2 + 276) + 2251) != 2) {
		nox_xxx_netPriMsgToPlayer_4DA2C0(a1, "pickup.c:ObjectEquipClassFail", 0);
		return 0;
	}
	if (*(uint32_t*)(*(uint32_t*)(v2 + 276) + 4 * v3 + 4244)) {
		nox_xxx_netPriMsgToPlayer_4DA2C0(a1, "objcoll.c:AlreadyHaveGuide", 0);
		result = 0;
	} else {
		nox_xxx_awardBeastGuide_4FAE80_magic_plyrgide(a1, v3, 1);
		nox_xxx_delayedDeleteObject_4E5CC0(a2);
		result = 1;
	}
	return result;
}

//----- (0053F9E0) --------------------------------------------------------
int nox_xxx_useSpellReward_53F9E0(int a1, int a2) {
	unsigned char* v2; // ebx
	int v3;            // ebp
	int v4;            // edi
	int v5;            // ecx
	char v6;           // al

	v2 = *(unsigned char**)(a2 + 736);
	v3 = 0;
	v4 = *(uint32_t*)(a1 + 748);
	if (!(*(uint8_t*)(a1 + 8) & 4)) {
		return 0;
	}
	v5 = *(uint32_t*)(v4 + 276);
	v6 = *(uint8_t*)(v5 + 2251);
	if (v6 != 1 && v6 != 2) {
		nox_xxx_netPriMsgToPlayer_4DA2C0(a1, "use.c:SpellRewardClassFail", 0);
		nox_xxx_aud_501960(925, a1, 2, *(uint32_t*)(a1 + 36));
		return 0;
	}
	if (nox_xxx_playerCheckSpellClass_57AEA0(*(unsigned char*)(v5 + 2251), *v2)) {
		nox_xxx_netPriMsgToPlayer_4DA2C0(a1, "use.c:SpellRewardClassFail", 0);
		nox_xxx_aud_501960(925, a1, 2, *(uint32_t*)(a1 + 36));
		return 0;
	}
	if (nox_common_gameFlags_check_40A5C0(6144) && !*(uint32_t*)(*(uint32_t*)(v4 + 276) + 4 * *v2 + 3696)) {
		v3 = 1;
	}
	if (nox_xxx_spellGrantToPlayer_4FB550(a1, *v2, 1, v3, 0)) {
		nox_xxx_delayedDeleteObject_4E5CC0(a2);
	} else {
		nox_xxx_aud_501960(925, a1, 2, *(uint32_t*)(a1 + 36));
	}
	return 1;
}

//----- (0053FAE0) --------------------------------------------------------
int nox_xxx_useAbilityReward_53FAE0(int a1, int a2) {
	unsigned char* v2; // ebx
	int v3;            // ebp
	int v4;            // edi
	int result;        // eax

	v2 = *(unsigned char**)(a2 + 736);
	v3 = 0;
	v4 = *(uint32_t*)(a1 + 748);
	if (!(*(uint8_t*)(a1 + 8) & 4)) {
		return 0;
	}
	if (*(uint8_t*)(*(uint32_t*)(v4 + 276) + 2251)) {
		nox_xxx_netPriMsgToPlayer_4DA2C0(a1, "pickup.c:ObjectEquipClassFail", 0);
		nox_xxx_aud_501960(925, a1, 2, *(uint32_t*)(a1 + 36));
		result = 0;
	} else {
		if (nox_common_gameFlags_check_40A5C0(6144) && !*(uint32_t*)(*(uint32_t*)(v4 + 276) + 4 * *v2 + 3696)) {
			v3 = 1;
		}
		if (nox_xxx_abilityRewardServ_4FB9C0_ability(a1, *v2, v3)) {
			nox_xxx_delayedDeleteObject_4E5CC0(a2);
		} else {
			nox_xxx_aud_501960(925, a1, 2, *(uint32_t*)(a1 + 36));
		}
		result = 1;
	}
	return result;
}

//----- (0053FBC0) --------------------------------------------------------
uint32_t* nox_xxx_respawnPlayerImpl_53FBC0(float* a1, int a2) {
	int v2;           // eax
	int v3;           // ebx
	float* v4;        // esi
	uint32_t* result; // eax
	uint32_t* v6;     // edi
	int v7;           // eax
	int v8;           // eax
	float v9;         // [esp+0h] [ebp-1Ch]
	float v10;        // [esp+4h] [ebp-18h]
	int* v11;         // [esp+24h] [ebp+8h]

	if (!*getMemU32Ptr(0x5D4594, 2488736)) {
		nox_xxx_createCorpse_53FCA0();
	}
	v2 = nox_xxx_math_509EA0(a2);
	v3 = 0;
	v4 = getMemFloatPtr(0x587000, 280376 + 88 * v2);
	v11 = getMemIntPtr(0x5D4594, 2488740 + 44 * v2);
	do {
		result = nox_xxx_newObjectWithTypeInd_4E3450(*v11);
		v6 = result;
		if (!result) {
			break;
		}
		if (dword_5d4594_2650652) {
			if (nox_common_gameFlags_check_40A5C0(0x2000)) {
				v7 = v6[4];
				LOBYTE(v7) = v7 | 0x40;
				v6[4] = v7;
			}
		}
		v10 = v4[1] + a1[1];
		v9 = *v4 + *a1;
		nox_xxx_createAt_4DAA50((int)v6, 0, v9, v10);
		v8 = nox_common_randomInt_415FA0(10, 20);
		result = (uint32_t*)nox_xxx_unitSetDecayTime_511660(v6, gameFPS() * v8);
		++v3;
		v4 += 2;
		++v11;
	} while (v3 < 11);
	return result;
}

//----- (0053FCA0) --------------------------------------------------------
void nox_xxx_createCorpse_53FCA0() {
	int i;             // ebx
	unsigned char* v1; // edi
	char v2[32];       // [esp+Ch] [ebp-20h]

	for (i = 0; i < 9; ++i) {
		switch (i) {
		case 0:
			v1 = getMemAt(0x587000, 281216);
			break;
		case 1:
			v1 = getMemAt(0x587000, 281220);
			break;
		case 2:
			v1 = getMemAt(0x587000, 281224);
			break;
		case 3:
			v1 = getMemAt(0x587000, 281228);
			break;
		case 5:
			v1 = getMemAt(0x587000, 281232);
			break;
		case 6:
			v1 = getMemAt(0x587000, 281236);
			break;
		case 7:
			v1 = getMemAt(0x587000, 281240);
			break;
		case 8:
			v1 = getMemAt(0x587000, 281244);
			break;
		default:
			continue;
		}
		nox_sprintf(v2, "CorpseSkull%s", v1);
		*getMemU32Ptr(0x5D4594, 2488740 + 44 * i) = nox_xxx_getNameId_4E3AA0(v2);
		nox_sprintf(v2, "CorpseRibCage%s", v1);
		*getMemU32Ptr(0x5D4594, 2488744 + 44 * i) = nox_xxx_getNameId_4E3AA0(v2);
		nox_sprintf(v2, "CorpsePelvis%s", v1);
		*getMemU32Ptr(0x5D4594, 2488748 + 44 * i) = nox_xxx_getNameId_4E3AA0(v2);
		nox_sprintf(v2, "CorpseLeftLowerLeg%s", v1);
		*getMemU32Ptr(0x5D4594, 2488752 + 44 * i) = nox_xxx_getNameId_4E3AA0(v2);
		nox_sprintf(v2, "CorpseLeftUpperLeg%s", v1);
		*getMemU32Ptr(0x5D4594, 2488756 + 44 * i) = nox_xxx_getNameId_4E3AA0(v2);
		nox_sprintf(v2, "CorpseLeftLowerArm%s", v1);
		*getMemU32Ptr(0x5D4594, 2488760 + 44 * i) = nox_xxx_getNameId_4E3AA0(v2);
		nox_sprintf(v2, "CorpseLeftUpperArm%s", v1);
		*getMemU32Ptr(0x5D4594, 2488764 + 44 * i) = nox_xxx_getNameId_4E3AA0(v2);
		nox_sprintf(v2, "CorpseRightLowerLeg%s", v1);
		*getMemU32Ptr(0x5D4594, 2488768 + 44 * i) = nox_xxx_getNameId_4E3AA0(v2);
		nox_sprintf(v2, "CorpseRightUpperLeg%s", v1);
		*getMemU32Ptr(0x5D4594, 2488772 + 44 * i) = nox_xxx_getNameId_4E3AA0(v2);
		nox_sprintf(v2, "CorpseRightLowerArm%s", v1);
		*getMemU32Ptr(0x5D4594, 2488776 + 44 * i) = nox_xxx_getNameId_4E3AA0(v2);
		nox_sprintf(v2, "CorpseRightUpperArm%s", v1);
		*getMemU32Ptr(0x5D4594, 2488780 + 44 * i) = nox_xxx_getNameId_4E3AA0(v2);
	}
	*getMemU32Ptr(0x5D4594, 2488736) = 1;
}

//----- (00540440) --------------------------------------------------------
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

//----- (00542BF0) --------------------------------------------------------
void* nox_objectTypeGetXfer(char* id);
char* sub_542BF0(int a1, int a2, int a3) {
	int** v3;         // eax
	int* v4;          // ebx
	char* v5;         // esi
	char* v6;         // edx
	const char* v7;   // eax
	int v8;           // esi
	char* v9;         // eax
	char* v10;        // eax
	int (*v11)(int*); // eax
	const char* v12;  // edx
	char* v13;        // eax
	const char* v14;  // edx
	char* v15;        // eax
	const char* v16;  // edx
	const char* v17;  // edx
	char* v18;        // eax
	const char* v19;  // edx
	char* v20;        // eax
	const char* v21;  // edx
	char* v22;        // eax
	const char* v23;  // edx
	char* v24;        // eax
	const char* v25;  // edx
	char* v26;        // eax
	const char* v27;  // edx
	char* v28;        // eax
	const char* v29;  // edx
	char* v30;        // eax
	const char* v31;  // edx
	char* v32;        // eax
	const char* v33;  // edx
	const char* v34;  // edx
	const char* v35;  // eax
	char* v36;        // eax
	const char* v37;  // eax
	char* v38;        // eax
	const char* v39;  // eax
	char* v40;        // eax
	const char* v41;  // eax
	char* result;     // eax
	char* i;          // ebx
	char* v44;        // [esp-14h] [ebp-28h]
	char* v45;        // [esp-14h] [ebp-28h]
	char* v46;        // [esp-14h] [ebp-28h]
	char* v47;        // [esp-14h] [ebp-28h]
	int v48;          // [esp+10h] [ebp-4h]

	v3 = (int**)sub_5049D0();
	v48 = (int)v3;
	if (v3) {
		while (1) {
			v4 = *v3;
			if (((*v3)[4] & 0x80000000) == 0x80000000) {
				if (*v4) {
					v5 = sub_543620(*v4, a1);
					v6 = (char*)realloc((void*)*v4, strlen(v5) + 1);
					*v4 = (int)v6;
					strcpy(v6, v5);
				}
				v7 = (const char*)nox_script_objCallbackName_508CB0(v4, 14);
				if (v7) {
					v8 = a3;
					if (strlen(v7)) {
						v9 = sub_5435C0((int)v7, a1, a2, a3);
						sub_509120(v4, 14, v9);
					}
				} else {
					v8 = a3;
				}
				v10 = (char*)nox_xxx_getUnitName_4E39D0((int)v4);
				v11 = nox_objectTypeGetXfer(v10);
				if (v11 == nox_xxx_unitTriggerXfer_4F4E50) {
					v12 = (const char*)nox_script_objCallbackName_508CB0(v4, 1);
					if (strlen(v12)) {
						v13 = sub_5435C0((int)v12, a1, a2, v8);
						sub_509120(v4, 1, v13);
					}
					v14 = (const char*)nox_script_objCallbackName_508CB0(v4, 2);
					if (strlen(v14)) {
						v15 = sub_5435C0((int)v14, a1, a2, v8);
						sub_509120(v4, 2, v15);
					}
					v16 = (const char*)nox_script_objCallbackName_508CB0(v4, 0);
					if (strlen(v16)) {
						v44 = sub_5435C0((int)v16, a1, a2, v8);
						sub_509120(v4, 0, v44);
					}
				} else if (v11 == nox_xxx_XFerMonster_528DB0) {
					v17 = (const char*)nox_script_objCallbackName_508CB0(v4, 3);
					if (strlen(v17)) {
						v18 = sub_5435C0((int)v17, a1, a2, v8);
						sub_509120(v4, 3, v18);
					}
					v19 = (const char*)nox_script_objCallbackName_508CB0(v4, 5);
					if (strlen(v19)) {
						v20 = sub_5435C0((int)v19, a1, a2, v8);
						sub_509120(v4, 5, v20);
					}
					v21 = (const char*)nox_script_objCallbackName_508CB0(v4, 4);
					if (strlen(v21)) {
						v22 = sub_5435C0((int)v21, a1, a2, v8);
						sub_509120(v4, 4, v22);
					}
					v23 = (const char*)nox_script_objCallbackName_508CB0(v4, 6);
					if (strlen(v23)) {
						v24 = sub_5435C0((int)v23, a1, a2, v8);
						sub_509120(v4, 6, v24);
					}
					v25 = (const char*)nox_script_objCallbackName_508CB0(v4, 7);
					if (strlen(v25)) {
						v26 = sub_5435C0((int)v25, a1, a2, v8);
						sub_509120(v4, 7, v26);
					}
					v27 = (const char*)nox_script_objCallbackName_508CB0(v4, 8);
					if (strlen(v27)) {
						v28 = sub_5435C0((int)v27, a1, a2, v8);
						sub_509120(v4, 8, v28);
					}
					v29 = (const char*)nox_script_objCallbackName_508CB0(v4, 9);
					if (strlen(v29)) {
						v30 = sub_5435C0((int)v29, a1, a2, v8);
						sub_509120(v4, 9, v30);
					}
					v31 = (const char*)nox_script_objCallbackName_508CB0(v4, 10);
					if (strlen(v31)) {
						v32 = sub_5435C0((int)v31, a1, a2, v8);
						sub_509120(v4, 10, v32);
					}
					v33 = (const char*)nox_script_objCallbackName_508CB0(v4, 11);
					if (strlen(v33)) {
						v45 = sub_5435C0((int)v33, a1, a2, v8);
						sub_509120(v4, 11, v45);
					}
				} else if (v11 == nox_xxx_XFerHole_4F51D0) {
					v34 = (const char*)nox_script_objCallbackName_508CB0(v4, 12);
					if (strlen(v34)) {
						v46 = sub_5435C0((int)v34, a1, a2, v8);
						sub_509120(v4, 12, v46);
					}
				} else if (v11 == nox_xxx_XFerMonsterGen_4F7130) {
					v35 = (const char*)nox_script_objCallbackName_508CB0(v4, 15);
					if (v35 && strlen(v35)) {
						v36 = sub_5435C0((int)v35, a1, a2, v8);
						sub_509120(v4, 15, v36);
					}
					v37 = (const char*)nox_script_objCallbackName_508CB0(v4, 16);
					if (v37 && strlen(v37)) {
						v38 = sub_5435C0((int)v37, a1, a2, v8);
						sub_509120(v4, 16, v38);
					}
					v39 = (const char*)nox_script_objCallbackName_508CB0(v4, 18);
					if (v39 && strlen(v39)) {
						v40 = sub_5435C0((int)v39, a1, a2, v8);
						sub_509120(v4, 18, v40);
					}
					v41 = (const char*)nox_script_objCallbackName_508CB0(v4, 17);
					if (v41 && strlen(v41)) {
						v47 = sub_5435C0((int)v41, a1, a2, v8);
						sub_509120(v4, 17, v47);
					}
				}
				v4[4] &= 0x7FFFFFFFu;
				v3 = (int**)v48;
			}
			v48 = sub_5049E0((int)v3);
			if (!v48) {
				break;
			}
			v3 = (int**)v48;
		}
	}
	result = (char*)nox_xxx_waypointGetList_579860();
	for (i = result; result; i = result) {
		if (*((int*)i + 120) < 0) {
			if ((int)strlen(i + 16) > 0) {
				strcpy(i + 16, sub_543620((int)(i + 16), a1));
			}
			*((uint32_t*)i + 120) &= 0x7FFFFFFFu;
		}
		result = (char*)nox_xxx_waypointNext_579870((int)i);
	}
	return result;
}
// 543110: using guessed type char NewFileName[2048];

//----- (005435C0) --------------------------------------------------------
char* sub_5435C0(int a1, int a2, int a3, int a4) {
	nox_sprintf((char*)getMemAt(0x5D4594, 2489164), "%s%%%d%%%d%%%d", a1, a2, a3, a4);
	strlen((const char*)getMemAt(0x5D4594, 2489164));
	strcpy((char*)getMemAt(0x5D4594, 2489164), "ERROR_NAME_TOO_LONG!");
	return (char*)getMemAt(0x5D4594, 2489164);
}

//----- (00543620) --------------------------------------------------------
char* sub_543620(int a1, int a2) {
	nox_sprintf((char*)getMemAt(0x5D4594, 2489164), "%s%%%d", a1, a2);
	strlen((const char*)getMemAt(0x5D4594, 2489164));
	strcpy((char*)getMemAt(0x5D4594, 2489164), "ERROR_NAME_TOO_LONG!");
	return (char*)getMemAt(0x5D4594, 2489164);
}
