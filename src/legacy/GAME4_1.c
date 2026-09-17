#include <math.h>

#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME1_2.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_1.h"
#include "GAME4_2.h"
#include "GAME4_3.h"
#include "GAME5.h"
#include "GAME5_2.h"
#include "client__system__ctrlevnt.h"
#include "common__binfile.h"
#include "common__crypt.h"
#include "common__net_list.h"
#include "common__random.h"
#include "operators.h"
#include "common__system__team.h"
#include "server__gamemech__explevel.h"
#include "server__magic__plyrspel.h"
#include "server__system__trade.h"

#include "client__gui__window.h"
#include "client__video__draw_common.h"
#include "common__magic__speltree.h"
#include "server__script__script.h"

extern uint32_t dword_5d4594_3835364;
extern uint32_t dword_5d4594_1599708;
extern uint32_t dword_587000_234176;
extern uint32_t dword_5d4594_2487244;
extern uint32_t dword_587000_237036;
extern void* nox_alloc_pendingOwn_2386916;
extern uint32_t dword_5d4594_2386228;
extern void* nox_alloc_spawn_2386216;
extern uint32_t dword_5d4594_3835348;
extern void* nox_alloc_tradeSession_2386492;
extern uint32_t dword_5d4594_2386920;
extern void* nox_alloc_monsterList_2386220;
extern uint32_t dword_5d4594_2386500;
extern uint32_t dword_5d4594_2386212;
extern void* nox_alloc_tradeItems_2386496;
extern uint32_t dword_5d4594_2386224;

extern uint32_t nox_tile_def_cnt;
extern nox_tileDef_t nox_tile_defs_arr[176];

void* nox_monsterBin_head_2386924 = 0;













//----- (00509FF0) --------------------------------------------------------
int sub_509FF0(int a1) {
	int result; // eax

	result = a1;
	if (*(uint8_t*)(*(uint32_t*)a1 + 16) & 0x20) {
		*(uint32_t*)a1 = 0;
	}
	return result;
}

//----- (0050A010) --------------------------------------------------------
int nox_xxx_monsterActionIsCondition_50A010(int a1) {
	int result; // eax

	result = a1 < 39;
	LOBYTE(result) = a1 > 39;
	return result;
}

//----- (0050A020) --------------------------------------------------------
int nox_xxx_mobActionGet_50A020(int a1) {
	return *(uint32_t*)(*(uint32_t*)(a1 + 748) + 24 * (*(char*)(*(uint32_t*)(a1 + 748) + 544) + 23));
}

//----- (0050A090) --------------------------------------------------------
int nox_xxx_monsterIsActionScheduled_50A090(int a1, int a2) {
	int v2;      // ecx
	int v3;      // eax
	uint32_t* i; // ecx

	v2 = *(uint32_t*)(a1 + 748);
	v3 = *(char*)(v2 + 544) - 1;
	if (v3 < 0) {
		return 0;
	}
	for (i = (uint32_t*)(v2 + 8 * (3 * v3 + 69)); *i != a2; i -= 6) {
		if (--v3 < 0) {
			return 0;
		}
	}
	return 1;
}

//----- (0050A360) --------------------------------------------------------
int* nox_xxx_monsterAction_50A360(int a1, int a2) {
	int* result; // eax

	if (*(uint8_t*)(a1 + 8) & 2 &&
		*(uint32_t*)(*(uint32_t*)(a1 + 748) + 24 * (*(char*)(*(uint32_t*)(a1 + 748) + 544) + 23)) != a2) {
		result = nox_xxx_monsterPushAction_50A260(a1, a2);
	} else {
		result = 0;
	}
	return result;
}

//----- (0050A3D0) --------------------------------------------------------
int nox_xxx_monsterCallDieFn_50A3D0(uint32_t* a1) {
	int v1;     // ebx
	int i;      // edi
	int result; // eax
	int v4;     // ecx
	int v5;     // eax
	int v6;     // eax
	int v7;     // ecx
	int v8;     // edi
	int v9;     // ecx
	int v10;    // eax
	int v11;    // eax
	int v12;    // esi

	v1 = a1[187];
	if (nox_common_gameFlags_check_40A5C0(4096)) {
		sub_50E1E0((int)a1);
	}
	for (i = nox_xxx_getFirstPlayerUnit_4DA7C0(); i; i = nox_xxx_getNextPlayerUnit_4DA7F0(i)) {
		if ((uint32_t*)nox_xxx_playerGetPossess_4DDF30(i) == a1) {
			nox_xxx_playerObserveClear_4DDEF0(i);
		}
	}
	nox_xxx_monsterClearActionStack_50A3A0((int)a1);
	nox_xxx_monsterPushAction_50A260((int)a1, 31);
	nox_xxx_monsterPushAction_50A260((int)a1, 30);
	result = nox_xxx_unitIsZombie_534A40((int)a1);
	if (!result) {
		v4 = a1[4];
		LOBYTE(v4) = v4 & 0x7F;
		a1[4] = v4;
		nox_xxx_action_4DA9F0(a1);
		nox_xxx_unitClearBuffs_4FF580((int)a1);
		if ((signed char)*(uint8_t*)(v1 + 1440) >= 0) {
			if (!nox_common_gameFlags_check_40A5C0(4096)) {
				goto LABEL_13;
			}
			v5 = nox_common_randomInt_415FA0(5, 8);
		} else {
			v5 = nox_common_randomInt_415FA0(10, 20);
		}
		nox_xxx_unitSetDecayTime_511660(a1, gameFPS() * v5);
	LABEL_13:
		v6 = a1[127];
		if (v6 && *(uint8_t*)(v6 + 8) & 4) {
			v7 = a1[3];
			v8 = *(uint32_t*)(v6 + 748);
			LOBYTE(v7) = v7 & 0x7F;
			a1[3] = v7;
			nox_xxx_netFxShield_0_4D9200(*(unsigned char*)(*(uint32_t*)(v8 + 276) + 2064), (int)a1);
			nox_xxx_netUnmarkMinimapObj_417300(*(unsigned char*)(*(uint32_t*)(v8 + 276) + 2064), (int)a1, 1);
		}
		v9 = a1[3];
		BYTE1(v9) &= 0xFEu;
		a1[3] = v9;
		nox_xxx_unitTransferSlaves_4EC4B0(a1);
		nox_xxx_unitClearOwner_4EC300(a1);
		v10 = a1[3];
		if (!(v10 & 0x2000)) {
			nox_xxx_dropAllItems_4EDA40(a1);
		}
		if (!nox_common_gameFlags_check_40A5C0(2048) && *(uint32_t*)(v1 + 2188) == 2 && *(uint32_t*)(v1 + 2184) == 2) {
			if (a1[130]) {
				v11 = nox_xxx_findParentChainPlayer_4EC580(a1[130]);
				if (*(uint8_t*)(v11 + 8) & 4) {
					sub_4FC0B0(v11, 1);
				}
			}
		}
		result = nox_common_gameFlags_check_40A5C0(4096);
		if (result) {
			v12 = a1[130];
			if (v12) {
				result = nox_xxx_findParentChainPlayer_4EC580(v12);
				if (*(uint8_t*)(result + 8) & 4) {
					result = sub_4D6170(result);
				}
			}
		}
	}
	return result;
}

//----- (0050A850) --------------------------------------------------------
char nox_xxx_updateNPCAnimData_50A850(nox_object_t* a1p) {
	int a1 = a1p;
	uint8_t* v1;       // esi
	unsigned char* v2; // eax
	unsigned char v3;  // cl
	unsigned char v4;  // cl
	unsigned char v5;  // dl
	unsigned char v6;  // cl

	v1 = *(uint8_t**)(a1 + 748);
	if (*(uint8_t*)(a1 + 12) & 0x10) {
		v2 = *(unsigned char**)&v1[24 * (char)v1[544] + 552];
		if (v2 == (unsigned char*)16 || v2 == (unsigned char*)17) {
			v1[483] = 0;
			return (char)v2;
		}
	}
	LOBYTE(v2) = v1[483];
	if (!(uint8_t)v2) {
		v2 = nox_xxx_unitNPCActionToAnim_533D00(a1);
		if (v2) {
			v1[480] = v2[9];
			if (!v2[9]) {
				v1[483] = 1;
				return (char)v2;
			}
			v3 = v1[482] + 1;
			v1[482] = v3;
			if (v3 >= v2[10] + 1) {
				v4 = v1[481];
				v1[482] = 0;
				v1[481] = ++v4;
				v5 = v4;
				v6 = v2[9];
				if (v5 >= v6) {
					if (*((uint32_t*)v2 + 3)) {
						v1[481] = 0;
						return (char)v2;
					}
					v1[481] = v6 - 1;
					v1[483] = 1;
					return (char)v2;
				}
			}
		}
	}
	return (char)v2;
}

//----- (0050A910) --------------------------------------------------------
int nox_xxx_mobAction_50A910(nox_object_t* a1p) {
	int a1 = a1p;
	int v1;       // ecx
	int v2;       // eax
	int result;   // eax
	uint32_t* v4; // esi
	int v5;       // edi
	uint32_t* v6; // ebx
	int v7;       // eax
	int v8;       // eax
	int v9;       // ecx
	int v10;      // eax
	int v11;      // ecx
	int v12;      // [esp+0h] [ebp-8h]
	int v13;      // [esp+4h] [ebp-4h]

	v1 = *(uint32_t*)(a1 + 748);
	v12 = v1;
	v2 = *(uint32_t*)(v1 + 1216);
	if (v2 && *(uint32_t*)(v2 + 16) & 0x8020) {
		*(uint32_t*)(v1 + 1216) = 0;
	}
	result = *(char*)(v1 + 544);
	if (result < 0) {
		return result;
	}
	v4 = (uint32_t*)(v1 + 8 * (3 * result + 69));
	v13 = result + 1;
	while (1) {
		v5 = 0;
		if (*getMemU32Ptr(0x587000, 230388 + 16 * *v4) > 0) {
			v6 = v4 + 1;
			do {
				if (*getMemU32Ptr(0x587000, 230392 + 4 * (v5 + 4 * *v4)) == 1 && *v6) {
					sub_509FF0((int)v6);
				}
				++v5;
				v6 += 2;
			} while (v5 < *getMemIntPtr(0x587000, 230388 + 16 * *v4));
			v1 = v12;
		}
		switch (*v4) {
		case 3:
			v7 = v4[3];
			if (v7) {
				v11 = *(uint32_t*)(v7 + 56);
				v4[1] = v11;
				v4[2] = *(uint32_t*)(v7 + 60);
			}
			break;
		case 7:
		case 8:
			v10 = v4[3];
			if (v10) {
				if (nox_xxx_unitCanInteractWith_5370E0(a1, v10, 0) || nox_xxx_checkMobAction_50A0D0(a1, 3)) {
					v7 = v4[3];
					v11 = *(uint32_t*)(v7 + 56);
					v4[1] = v11;
					v4[2] = *(uint32_t*)(v7 + 60);
				} else {
					v4[3] = 0;
				}
			}
			break;
		case 0xF:
			v7 = *(uint32_t*)(v1 + 1196);
			if (v7) {
				v11 = *(uint32_t*)(v7 + 56);
				v4[1] = v11;
				v4[2] = *(uint32_t*)(v7 + 60);
			}
			break;
		case 0x11:
			v8 = v4[3];
			if (v8 && nox_xxx_unitCanInteractWith_5370E0(a1, v8, 0)) {
				v9 = v4[3];
				v4[1] = *(uint32_t*)(v9 + 56);
				v4[2] = *(uint32_t*)(v9 + 60);
			}
			break;
		default:
			break;
		}
		v4 -= 6;
		result = --v13;
		if (!v13) {
			break;
		}
		v1 = v12;
	}
	return result;
}

//----- (0050CAC0) --------------------------------------------------------
void sub_50CAC0(int a1, int a2) {
	if (dword_5d4594_1599708 != 1) {
		if (nox_xxx_unitIsEnemyTo_5330C0(a2, a1)) {
			dword_5d4594_1599708 = 1;
		}
	}
}














//----- (00511C50) --------------------------------------------------------
nox_object_t* nox_xxx_script_511C50(int a1) {
	uint32_t* v1; // esi

	if (dword_587000_237036) {
		sub_511D20();
	}
	v1 = *(uint32_t**)getMemAt(0x5D4594, 2386820);
	if (!*getMemU32Ptr(0x5D4594, 2386820)) {
		return 0;
	}
	while (*(uint8_t*)(*v1 + 16) & 0x20 || *(uint32_t*)(*v1 + 44) != a1) {
		v1 = (uint32_t*)v1[2];
		if (!v1) {
			return 0;
		}
	}
	sub_511CE0(getMemAt(0x5D4594, 2386820), (int)v1);
	sub_511CB0(getMemAt(0x5D4594, 2386820), (int)v1);
	return *v1;
}

//----- (00511CB0) --------------------------------------------------------
int sub_511CB0(uint32_t* a1, int a2) {
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

//----- (00511CE0) --------------------------------------------------------
int sub_511CE0(uint32_t* a1, int a2) {
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

//----- (00511D20) --------------------------------------------------------
int sub_511D20() {
	unsigned char* v0; // esi
	int result;        // eax

	v0 = getMemAt(0x5D4594, 2386628);
	*getMemU32Ptr(0x5D4594, 2386820) = 0;
	*getMemU32Ptr(0x5D4594, 2386824) = 0;
	*getMemU32Ptr(0x5D4594, 2386620) = 0;
	*getMemU32Ptr(0x5D4594, 2386624) = 0;
	do {
		result = sub_511CB0(getMemAt(0x5D4594, 2386620), (int)v0);
		v0 += 12;
	} while ((int)v0 < (int)getMemAt(0x5D4594, 2386820));
	dword_587000_237036 = 0;
	return result;
}

//----- (00511D70) --------------------------------------------------------
int nox_xxx_scriptPrepareFoundUnit_511D70(nox_object_t* obj) {
	int a1 = obj;
	int* v1;    // eax
	int v2;     // esi
	int result; // eax
	int v4;     // [esp-8h] [ebp-8h]

	v1 = (int*)sub_511DC0();
	if (v1) {
		*v1 = a1;
		result = sub_511CB0(getMemAt(0x5D4594, 2386820), (int)v1);
	} else {
		v2 = *getMemU32Ptr(0x5D4594, 2386824);
		v4 = *getMemU32Ptr(0x5D4594, 2386824);
		**(uint32_t**)getMemAt(0x5D4594, 2386824) = a1;
		sub_511CE0(getMemAt(0x5D4594, 2386820), v4);
		result = sub_511CB0(getMemAt(0x5D4594, 2386820), v2);
	}
	return result;
}

//----- (00511DC0) --------------------------------------------------------
int sub_511DC0() {
	int result; // eax

	result = *getMemU32Ptr(0x5D4594, 2386620);
	if (!*getMemU32Ptr(0x5D4594, 2386620)) {
		return 0;
	}
	*getMemU32Ptr(0x5D4594, 2386620) = *(uint32_t*)(*getMemU32Ptr(0x5D4594, 2386620) + 8);
	return result;
}

//----- (00511DE0) --------------------------------------------------------
int sub_511DE0(nox_object_t* a1) {
	int result;   // eax
	uint32_t* v2; // esi

	result = dword_587000_237036;
	if (!dword_587000_237036) {
		v2 = *(uint32_t**)getMemAt(0x5D4594, 2386820);
		if (*getMemU32Ptr(0x5D4594, 2386820)) {
			result = a1;
			while (*v2 != a1) {
				v2 = (uint32_t*)v2[2];
				if (!v2) {
					return result;
				}
			}
			sub_511CE0(getMemAt(0x5D4594, 2386820), (int)v2);
			result = sub_511CB0(getMemAt(0x5D4594, 2386620), (int)v2);
		}
	}
	return result;
}

//----- (00511E20) --------------------------------------------------------
int sub_511E20() {
	int result; // eax
	int v1;     // esi
	int v2;     // edi

	result = dword_587000_237036;
	if (!dword_587000_237036) {
		v1 = *getMemU32Ptr(0x5D4594, 2386820);
		if (*getMemU32Ptr(0x5D4594, 2386820)) {
			do {
				v2 = *(uint32_t*)(v1 + 8);
				sub_511CE0(getMemAt(0x5D4594, 2386820), v1);
				result = sub_511CB0(getMemAt(0x5D4594, 2386620), v1);
				v1 = v2;
			} while (v2);
		}
	}
	return result;
}

//----- (005125A0) --------------------------------------------------------
float* nox_xxx_monsterLookAt_5125A0(nox_object_t* obj, int a2) {
	int a1 = obj;
	float* result; // eax
	int v3;        // edx
	float v4;      // [esp+0h] [ebp-8h]
	float v5;      // [esp+4h] [ebp-4h]

	result = (float*)nox_xxx_mathDirection4ToAngle_509E90(a2);
	if (*(uint8_t*)(a1 + 8) & 2) {
		v3 = *(uint32_t*)(a1 + 16);
		if ((v3 & 0x8000) == 0) {
			v4 = *getMemFloatPtr(0x587000, 194136 + 8 * (uint32_t)result) * 10.0 + *(float*)(a1 + 56);
			v5 = *getMemFloatPtr(0x587000, 194140 + 8 * (uint32_t)result) * 10.0 + *(float*)(a1 + 60);
			result = (float*)nox_xxx_monsterPushAction_50A260(a1, 25);
			if (result) {
				result[1] = v4;
				result[2] = v5;
			}
		}
	}
	return result;
}

//----- (00514110) --------------------------------------------------------
void nox_xxx_monsterWalkTo_514110(nox_object_t* obj, float x, float y) {
	int a1 = obj;
	int* result; // eax
	int* v4;     // eax

	result = *(int**)(a1 + 16);
	if (SBYTE1(result) >= 0 && *(uint8_t*)(a1 + 8) & 2) {
		nox_xxx_monsterClearActionStack_50A3A0(a1);
		v4 = nox_xxx_monsterPushAction_50A260(a1, 32);
		if (v4) {
			v4[1] = 8;
		}
		result = nox_xxx_monsterPushAction_50A260(a1, 8);
		if (result) {
			((float*)result)[1] = x;
			((float*)result)[2] = y;
			result[3] = 0;
		}
	}
}

//----- (00515680) --------------------------------------------------------
void nox_xxx_monsterGoPatrol_515680(nox_object_t* a1p, void* a2p) {
	int a1 = a1p;
	int a2 = a2p;
	int v2;    // ebx
	int v3;    // eax
	int* v4;   // edi
	float2 v5; // [esp+8h] [ebp-8h]

	v2 = *(uint32_t*)(a1 + 748);
	if (a1) {
		if (*(uint8_t*)(a1 + 8) & 2) {
			v3 = *(uint32_t*)(a1 + 16);
			if ((v3 & 0x8000) == 0) {
				v5.field_0 = *(float*)(a2 + 8) - *(float*)a2;
				v5.field_4 = *(float*)(a2 + 12) - *(float*)(a2 + 4);
				nox_xxx_monsterClearActionStack_50A3A0(a1);
				v4 = nox_xxx_monsterPushAction_50A260(a1, 4);
				if (v4) {
					v4[1] = *(uint32_t*)a2;
					v4[2] = *(uint32_t*)(a2 + 4);
					v4[3] = nox_xxx_math_509ED0(&v5);
				}
				*(uint32_t*)(v2 + 1312) = *(uint32_t*)(a2 + 16);
			}
		}
	}
}

//----- (005157A0) --------------------------------------------------------
void nox_xxx_unitHunt_5157A0(nox_object_t* obj) {
	int a1 = obj;
	int result; // eax

	if (a1 && *(uint8_t*)(a1 + 8) & 2) {
		result = *(int*)(a1 + 16);
		if (SBYTE1(result) >= 0) {
			nox_xxx_monsterClearActionStack_50A3A0(a1);
			nox_xxx_monsterPushAction_50A260(a1, 5);
		}
	}
}

//----- (00515820) --------------------------------------------------------
void nox_xxx_unitIdle_515820(nox_object_t* obj) {
	int a1 = obj;
	int result; // eax

	if (a1 && *(uint8_t*)(a1 + 8) & 2) {
		result = *(int*)(a1 + 16);
		if (SBYTE1(result) >= 0) {
			nox_xxx_monsterClearActionStack_50A3A0(a1);
			nox_xxx_monsterPushAction_50A260(a1, 0);
		}
	}
}

//----- (005158C0) --------------------------------------------------------
void nox_xxx_unitSetFollow_5158C0(nox_object_t* obj1, nox_object_t* obj2) {
	int a1 = obj1;
	int a2 = obj2;
	int v2;  // eax
	int* v3; // eax

	if (a1) {
		if (a2) {
			if (*(uint8_t*)(a1 + 8) & 2) {
				if (a1 != a2) {
					v2 = *(uint32_t*)(a1 + 16);
					if ((v2 & 0x8000) == 0) {
						nox_xxx_monsterClearActionStack_50A3A0(a1);
						v3 = nox_xxx_monsterPushAction_50A260(a1, 3);
						if (v3) {
							v3[1] = *(uint32_t*)(a2 + 56);
							v3[2] = *(uint32_t*)(a2 + 60);
							v3[3] = a2;
						}
					}
				}
			}
		}
	}
}

//----- (00515A30) --------------------------------------------------------
void nox_xxx_monsterActionMelee_515A30(nox_object_t* a1p, float2* a2) {
	int a1 = a1p;
	int v2;    // eax
	int* v3;   // eax
	float* v4; // edi
	float* v5; // eax

	if (a1) {
		if (*(uint8_t*)(a1 + 8) & 2) {
			v2 = *(uint32_t*)(a1 + 16);
			if ((v2 & 0x8000) == 0) {
				if (nox_xxx_monsterCanMelee_534220(a1)) {
					nox_xxx_monsterClearActionStack_50A3A0(a1);
					v3 = nox_xxx_monsterPushAction_50A260(a1, 32);
					if (v3) {
						v3[1] = 16;
					}
					nox_xxx_monsterPushAction_50A260(a1, 16);
					v4 = (float*)nox_xxx_monsterPushAction_50A260(a1, 51);
					if (v4) {
						v4[1] = sub_534470(a1) + *(float*)(a1 + 176);
						v4[3] = a2->field_0;
						v4[4] = a2->field_4;
					}
					v5 = (float*)nox_xxx_monsterPushAction_50A260(a1, 7);
					if (v5) {
						v5[1] = a2->field_0;
						v5[2] = a2->field_4;
						v5[3] = 0;
					}
				}
			}
		}
	}
}

//----- (00515B80) --------------------------------------------------------
void nox_xxx_monsterMissileAttack_515B80(nox_object_t* a1p, float2* a2) {
	int a1 = a1p;
	int v2;  // eax
	int* v3; // eax
	int* v4; // eax

	if (a1) {
		if (a2) {
			if (*(uint8_t*)(a1 + 8) & 2) {
				v2 = *(uint32_t*)(a1 + 16);
				if ((v2 & 0x8000) == 0) {
					if (nox_xxx_monsterCanShoot_534280(a1)) {
						nox_xxx_monsterClearActionStack_50A3A0(a1);
						v3 = nox_xxx_monsterPushAction_50A260(a1, 32);
						if (v3) {
							v3[1] = 17;
						}
						v4 = nox_xxx_monsterPushAction_50A260(a1, 17);
						if (v4) {
							v4[1] = *(uint32_t*)(&a2->field_0);
							v4[2] = *(uint32_t*)(&a2->field_4);
							v4[3] = 0;
						}
					}
				}
			}
		}
	}
}

//----- (00515C80) --------------------------------------------------------
int sub_515C80(int a1, uint8_t* a2) {
	int result; // eax

	result = a1;
	if (a1) {
		if (*(uint8_t*)(a1 + 8) & 2) {
			result = *(uint32_t*)(a1 + 748);
			*(uint8_t*)(result + 1332) = *a2;
		}
	}
	return result;
}

//----- (00515D30) --------------------------------------------------------
void nox_xxx_mobSetFightTarg_515D30(nox_object_t* a1p, nox_object_t* a2p) {
	int a1 = a1p;
	int a2 = a2p;
	int v2;  // ebx
	int v3;  // eax
	int* v4; // eax
	int* v5; // eax

	v2 = *(uint32_t*)(a1 + 748);
	if (a1) {
		if (a2) {
			if (*(uint8_t*)(a1 + 8) & 2) {
				if (a1 != a2) {
					v3 = *(uint32_t*)(a1 + 16);
					if ((v3 & 0x8000) == 0) {
						nox_xxx_monsterClearActionStack_50A3A0(a1);
						*(uint32_t*)(v2 + 1216) = a2;
						nox_xxx_frameCounterSetCopyToNextFrame_5281D0();
						v4 = nox_xxx_monsterPushAction_50A260(a1, 32);
						if (v4) {
							v4[1] = 15;
						}
						v5 = nox_xxx_monsterPushAction_50A260(a1, 15);
						if (v5) {
							v5[1] = *(uint32_t*)(a2 + 56);
							v5[2] = *(uint32_t*)(a2 + 60);
							v5[3] = gameFrame();
						}
					}
				}
			}
		}
	}
}

//----- (00515F70) --------------------------------------------------------
void nox_server_scriptFleeFrom_515F70(nox_object_t* a1p, void* a2p) {
	int a1 = a1p;
	uint32_t* a2 = a2p;
	int v2;  // eax
	int* v3; // eax
	int* v4; // eax
	int* v5; // eax
	int v6;  // edx
	int v7;  // edx

	if (a1) {
		if (*(uint8_t*)(a1 + 8) & 2) {
			v2 = *(uint32_t*)(a1 + 16);
			if ((v2 & 0x8000) == 0 && !nox_xxx_checkMobAction_50A0D0(a1, 24)) {
				v3 = nox_xxx_monsterPushAction_50A260(a1, 32);
				if (v3) {
					v3[1] = 24;
				}
				v4 = nox_xxx_monsterPushAction_50A260(a1, 41);
				if (v4) {
					v4[1] = gameFrame() + a2[1];
				}
				v5 = nox_xxx_monsterPushAction_50A260(a1, 24);
				if (v5) {
					v6 = *a2;
					v5[1] = *(uint32_t*)(*a2 + 56);
					v7 = *(uint32_t*)(v6 + 60);
					v5[3] = 0;
					v5[2] = v7;
				}
			}
		}
	}
}

//----- (00516090) --------------------------------------------------------
void sub_516090(nox_object_t* a1p, uint32_t a2) {
	int a1 = a1p;
	int v2;  // eax
	int* v3; // eax
	int* v4; // eax

	if (a1) {
		if (*(uint8_t*)(a1 + 8) & 2) {
			v2 = *(uint32_t*)(a1 + 16);
			if ((v2 & 0x8000) == 0) {
				v3 = nox_xxx_monsterPushAction_50A260(a1, 32);
				if (v3) {
					v3[1] = 1;
				}
				v4 = nox_xxx_monsterPushAction_50A260(a1, 1);
				if (v4) {
					v4[1] = gameFrame() + a2;
				}
			}
		}
	}
}

//----- (00516570) --------------------------------------------------------
extern uint32_t nox_gameDisableMapDraw_5d4594_2650672;
int sub_516570() {
	uint8_t* v2 = 0;
	int v1 = nox_xxx_getFirstPlayerUnit_4DA7C0();
	if (v1) {
		do {
			v2 = *(uint8_t**)(*(uint32_t*)(v1 + 748) + 276);
			if (v2[2064] == 31) {
				break;
			}
			v1 = nox_xxx_getNextPlayerUnit_4DA7F0(v1);
		} while (v1);
	}
	nox_gameDisableMapDraw_5d4594_2650672 = 1;
	if (!v2) {
		return 0;
	}
	return nox_xxx_netSendChapterEnd_4D9560((unsigned char)v2[2064], getMemByte(0x5D4594, 2386828), *getMemIntPtr(0x5D4594, 2386832));
}

//----- (00516D00) --------------------------------------------------------
unsigned int sub_516D00(nox_object_t* a1p) {
	int a1 = a1p;
	unsigned int result; // eax
	int v2;              // edx

	result = a1;
	if (a1 && *(uint8_t*)(a1 + 8) & 2) {
		v2 = *(uint32_t*)(a1 + 16);
		if ((v2 & 0x8000) != 0) {
			*(uint32_t*)(*(uint32_t*)(a1 + 748) + 1440) &= 0xFFEFFFFF;
			result = nox_xxx_mobRaiseZombie_534AB0(a1);
		}
	}
	return result;
}

//----- (00516EE0) --------------------------------------------------------
int nox_xxx_allocPendingOwnsArray_516EE0() {
	dword_5d4594_2386920 = 0;
	nox_alloc_pendingOwn_2386916 = nox_new_alloc_class("PendingOwn", 12, 512);
	return nox_alloc_pendingOwn_2386916 != 0;
}

//----- (00516F10) --------------------------------------------------------
int sub_516F10() {
	int result; // eax

	nox_free_alloc_class(*(void**)&nox_alloc_pendingOwn_2386916);
	result = 0;
	nox_alloc_pendingOwn_2386916 = 0;
	dword_5d4594_2386920 = 0;
	return result;
}

//----- (00516F30) --------------------------------------------------------
void sub_516F30() {
	nox_alloc_class_free_all(*(uint32_t**)&nox_alloc_pendingOwn_2386916);
	dword_5d4594_2386920 = 0;
}

//----- (00516F90) --------------------------------------------------------
uint32_t* sub_516F90(int a1, int a2) {
	uint32_t* result; // eax

	result = nox_alloc_class_new_obj_zero(*(uint32_t**)&nox_alloc_pendingOwn_2386916);
	if (result) {
		*result = a1;
		result[1] = a2;
		result[2] = dword_5d4594_2386920;
		dword_5d4594_2386920 = result;
	}
	return result;
}

//----- (00516FC0) --------------------------------------------------------
void sub_516FC0() {
	int* v0; // esi
	int v1;  // edi
	int v2;  // eax

	v0 = *(int**)&dword_5d4594_2386920;
	if (dword_5d4594_2386920) {
		do {
			v1 = sub_4ECF10(*v0);
			v2 = sub_4ECF10(v0[1]);
			if (v1 && v2) {
				nox_xxx_unitSetOwner_4EC290(v1, v2);
			}
			v0 = (int*)v0[2];
		} while (v0);
	}
	sub_516F30();
}

//----- (00517010) --------------------------------------------------------
int nox_xxx_loadMonsterBin_517010() {
	int result;   // eax
	FILE* v1;     // esi
	char v2[256]; // [esp+4h] [ebp-100h]

	nox_monsterBin_head_2386924 = 0;
	result = nox_binfile_open_408CC0("monster.bin", 0);
	v1 = (FILE*)result;
	if (result) {
		result = nox_binfile_cryptSet_408D40(result, 23);
		if (result) {
			while (nox_xxx_readStr_517090(v1, v2) && nox_xxx_servParseMonsterDef_517170(v1, v2)) {
				;
			}
			nox_binfile_close_408D90(v1);
			result = 1;
		}
	}
	return result;
}

//----- (00517090) --------------------------------------------------------
int nox_xxx_readStr_517090(FILE* a1, uint8_t* a2) {
	uint8_t* v2;          // ebx
	int v3;               // ecx
	int v4;               // edi
	int v5;               // ebp
	int v6;               // eax
	uint16_t CharType[2]; // [esp+10h] [ebp-4h]

	v2 = a2;
	v3 = 0;
	*(uint32_t*)CharType = 0;
	v4 = 1;
	do {
		while (1) {
			v5 = v3;
			nox_binfile_fread_408E40((char*)CharType, 1, 1, a1);
			if (nox_binfile_lastErr_409370(a1) == -1) {
				return 0;
			}
			v3 = *(uint32_t*)CharType;
			v6 = isspace(CharType[0]);
			if (v6) {
				break;
			}
			v4 = 0;
			if (v3 != 47 || v5 != 47) {
				*v2++ = v3;
			} else {
				sub_517140(a1);
				v2 = a2;
				v3 = *(uint32_t*)CharType;
				v4 = 1;
			}
		}
	} while (v4);
	*v2 = 0;
	return 1;
}

//----- (00517140) --------------------------------------------------------
int sub_517140(FILE* a1) {
	FILE* v1;   // esi
	int result; // eax

	v1 = a1;
	do {
		LOBYTE(a1) = 0;
		nox_binfile_fread_408E40((char*)&a1, 1, 1, v1);
		result = nox_binfile_lastErr_409370(v1);
	} while (result != -1 && (uint8_t)a1 != 10);
	return result;
}

//----- (00517170) --------------------------------------------------------
int nox_xxx_servParseMonsterDef_517170(FILE* a1, const char* a2) {
	int result;        // eax
	uint32_t* v3;      // ebx
	unsigned char* v4; // esi
	int v5;            // eax
	char* v6;          // edi
	int v7;            // edi
	unsigned char* v8; // esi
	int v9;            // [esp+10h] [ebp-104h]
	char v10[256];     // [esp+14h] [ebp-100h]

	result = (int)calloc(1u, 0xF8u);
	v3 = (uint32_t*)result;
	if (!result) {
		return 0;
	}
	strcpy((char*)result, a2);
	while (1) {
		if (!nox_xxx_readStr_517090(a1, v10) || !nox_strcmpi("END", v10)) {
			v3[61] = nox_monsterBin_head_2386924;
			nox_monsterBin_head_2386924 = v3;
			return 1;
		}
		if (nox_common_gameFlags_check_40A5C0(2048) || nox_common_gameFlags_check_40A5C0(0x200000)) {
			if (nox_strcmpi("ARENA", v10)) {
				if (nox_strcmpi("SOLO", v10)) {
					goto LABEL_10;
				}
			} else {
				sub_517140(a1);
			}
			continue;
		}
	LABEL_10:
		if (!nox_common_gameFlags_check_40A5C0(0x2000)) {
			goto LABEL_14;
		}
		if (!nox_strcmpi("SOLO", v10)) {
			sub_517140(a1);
			continue;
		}
		if (!nox_strcmpi("ARENA", v10)) {
			continue;
		}
	LABEL_14:
		v4 = getMemAt(0x587000, 248192);
		if (*getMemU32Ptr(0x587000, 248192)) {
			do {
				if (!nox_strcmpi(*(const char**)v4, v10)) {
					break;
				}
				v5 = *((uint32_t*)v4 + 3);
				v4 += 12;
			} while (v5);
		}
		if (!*(uint32_t*)v4) {
			free(v3);
			return 0;
		}
		v6 = (char*)v3 + *((uint32_t*)v4 + 2);
		switch (*((uint32_t*)v4 + 1)) {
		case 0:
			nox_xxx_readStr_517090(a1, v10);
			*(uint32_t*)v6 = atoi(v10);
			continue;
		case 1:
			nox_xxx_readStr_517090(a1, v10);
			*(float*)v6 = atof(v10);
			continue;
		case 2:
			nox_xxx_readStr_517090(a1, v10);
			*(uint32_t*)v6 = nox_xxx_utilFindSound_40AF50(v10);
			continue;
		case 3:
			nox_xxx_readStr_517090(a1, v10);
			if (nox_xxx_monsterLoadStrikeFn_549040((int)v3, v10)) {
				continue;
			}
			return 0;
		case 4:
			nox_xxx_readStr_517090(a1, v10);
			if (nox_xxx_monsterLoadDieFn_5490E0((int)v3, v10)) {
				continue;
			}
			return 0;
		case 5:
			nox_xxx_readStr_517090(a1, v10);
			if (nox_xxx_monsterLoadDeadFn_549180((int)v3, v10)) {
				continue;
			}
			return 0;
		case 6:
			v9 = 0;
			nox_xxx_readStr_517090(a1, v10);
			set_bitmask_flags_from_plus_separated_names_423930(
				v10, &v9, (const char**)getMemAt(0x587000, 247536));
			*(uint16_t*)v6 = v9;
			continue;
		case 7:
			nox_xxx_readStr_517090(a1, v6);
			if (!strcmp("NULL", v6)) {
				*v6 = 0;
			}
			continue;
		case 8:
			nox_xxx_readStr_517090(a1, v10);
			v3[31] = 18;
			v7 = 0;
			v8 = getMemAt(0x587000, 247464);
			break;
		default:
			continue;
		}
		while (nox_strcmpi(v10, (const char*)(*(uint32_t*)v8 + 7))) {
			v8 += 4;
			++v7;
			if ((int)v8 >= (int)getMemAt(0x587000, 247536)) {
				goto LABEL_27;
			}
		}
		v3[31] = v7;
	LABEL_27:
		if (v3[31] == 18) {
			return 0;
		}
	}
	return result;
}

//----- (005174F0) --------------------------------------------------------
uint32_t* nox_xxx_monsterListFree_5174F0() {
	uint32_t* result; // eax
	uint32_t* v1;     // esi

	result = nox_monsterBin_head_2386924;
	if (nox_monsterBin_head_2386924) {
		do {
			v1 = (uint32_t*)result[61];
			free(result);
			result = v1;
		} while (v1);
		nox_monsterBin_head_2386924 = 0;
	}
	return result;
}

//----- (00517520) --------------------------------------------------------
int nox_xxx_monsterList_517520() {
	int v0; // esi
	int v1; // eax

	v0 = nox_monsterBin_head_2386924;
	if (!nox_monsterBin_head_2386924) {
		return 1;
	}
	while (1) {
		v1 = nox_xxx_getNameId_4E3AA0((char*)v0);
		*(uint32_t*)(v0 + 240) = v1;
		if (!v1) {
			break;
		}
		v0 = *(uint32_t*)(v0 + 244);
		if (!v0) {
			return 1;
		}
	}
	nox_xxx_monsterListFree_5174F0();
	return 0;
}

//----- (00517560) --------------------------------------------------------
void* nox_xxx_monsterDefByTT_517560(int a1) {
	uint32_t* result; // eax

	result = nox_monsterBin_head_2386924;
	if (!nox_monsterBin_head_2386924) {
		return 0;
	}
	while (result[60] != a1) {
		result = (uint32_t*)result[61];
		if (!result) {
			return 0;
		}
	}
	return result;
}



//----- (0051A500) --------------------------------------------------------
int sub_51A500(int a1) {
	int v1;           // ecx
	int v2;           // edx
	unsigned char* i; // eax
	int v4;           // esi

	if (!*getMemU32Ptr(0x5D4594, 2388664)) {
		sub_51A550();
	}
	if (!a1) {
		return 0;
	}
	v1 = 0;
	if (!*getMemU32Ptr(0x587000, 249904)) {
		return 0;
	}
	HIWORD(v2) = 0;
	for (i = getMemAt(0x587000, 249904);; i += 16) {
		LOWORD(v2) = *(uint16_t*)(a1 + 4);
		if (*((uint32_t*)i - 1) == v2) {
			break;
		}
		v4 = *((uint32_t*)i + 4);
		++v1;
		if (!v4) {
			return 0;
		}
	}
	return *getMemU32Ptr(0x587000, 249908 + 16 * v1);
}

//----- (0051A550) --------------------------------------------------------
char* sub_51A550() {
	char* result;      // eax
	unsigned char* v1; // esi

	result = *(char**)getMemAt(0x587000, 249904);
	if (*getMemU32Ptr(0x587000, 249904)) {
		v1 = getMemAt(0x587000, 249896);
		do {
			*((uint32_t*)v1 + 3) = nox_xxx_getNameId_4E3AA0(result);
			*((uint32_t*)v1 + 1) = nox_xxx_getNameId_4E3AA0(*(char**)v1);
			result = (char*)*((uint32_t*)v1 + 6);
			v1 += 16;
		} while (result);
		*getMemU32Ptr(0x5D4594, 2388664) = 1;
	} else {
		*getMemU32Ptr(0x5D4594, 2388664) = 1;
	}
	return result;
}

//----- (0051A5A0) --------------------------------------------------------
void* nox_xxx_objectTypeByIndHealthData(int a1);
void nox_xxx_spawnHecubahQuest_51A5A0(int* a1) {
	uint32_t* v1;  // edi
	uint32_t* v2;  // esi
	int v3;        // eax
	int v4;        // eax
	double v5;     // st7
	int v6;        // eax
	uint32_t* v7;  // esi
	int v8;        // eax
	uint32_t* v9;  // eax
	int v10;       // eax
	uint32_t* v11; // eax
	int v12;       // eax
	uint32_t* v13; // eax
	int v14;       // eax
	uint32_t* v15; // eax
	float v16;     // [esp+8h] [ebp-8h]
	float v17;     // [esp+Ch] [ebp-4h]

	v1 = nox_xxx_newObjectByTypeID_4E3810("Hecubah");
	v16 = sub_4E40F0();
	if (v1) {
		v2 = (uint32_t*)v1[187];
		v3 = v2[121];
		if (v3) {
			v4 = *(uint32_t*)(v3 + 72);
		} else {
			v4 = *(unsigned short*)((int)nox_xxx_objectTypeByIndHealthData(*((unsigned short*)v1 + 2)) + 4);
		}
		if (v16 < 1.0) {
			v16 = 1.0;
		}
		v5 = (double)v4 * v16;
		v17 = v5;
		nox_xxx_unitSetHP_4E4560((int)v1, (long long)v5);
		*(uint16_t*)(v1[139] + 4) = nox_float2int(v17);
		if (!*(uint16_t*)v1[139]) {
			nox_xxx_unitSetHP_4E4560((int)v1, 1u);
		}
		v6 = v1[139];
		if (!*(uint16_t*)(v6 + 4)) {
			*(uint16_t*)(v6 + 4) = 1;
		}
		v2[411] = 0x10000000;
		v2[423] = 0x10000000;
		v2[340] = 4;
		v2[326] = 1062501089;
		v2[510] = 3;
		v2[410] = 0x8000000;
		v2[444] = 0x20000000;
		v2[388] = 0x40000000;
		v2[415] = 0x40000000;
		nox_xxx_gamedataGetFloat_419D40("HecubahQuestSkill");
		v2[330] = 1062836634;
		nox_xxx_createAt_4DAA50((int)v1, 0, *(float*)a1, *((float*)a1 + 1));
		v7 = nox_xxx_newObjectByTypeID_4E3810("RewardMarker");
		if (v7) {
			v8 = nox_game_getQuestStage_4E3CC0();
			v9 = nox_server_rewardgen_activateMarker_4F0720((int)v7, v8 + 2);
			if (v9) {
				nox_xxx_inventoryPutImpl_4F3070((int)v1, (int)v9, 0);
			}
			v10 = nox_game_getQuestStage_4E3CC0();
			v11 = nox_server_rewardgen_activateMarker_4F0720((int)v7, v10 + 2);
			if (v11) {
				nox_xxx_inventoryPutImpl_4F3070((int)v1, (int)v11, 0);
			}
			v12 = nox_game_getQuestStage_4E3CC0();
			v13 = nox_server_rewardgen_activateMarker_4F0720((int)v7, v12 + 2);
			if (v13) {
				nox_xxx_inventoryPutImpl_4F3070((int)v1, (int)v13, 0);
			}
			v14 = nox_game_getQuestStage_4E3CC0();
			v15 = nox_server_rewardgen_activateMarker_4F0720((int)v7, v14 + 2);
			if (v15) {
				nox_xxx_inventoryPutImpl_4F3070((int)v1, (int)v15, 0);
			}
			nox_xxx_objectFreeMem_4E38A0((int)v7);
		}
	}
}

//----- (0051A7A0) --------------------------------------------------------
void nox_xxx_spawnNecroQuest_51A7A0(int* a1) {
	uint32_t* v1; // edi
	uint32_t* v2; // esi
	int v3;       // eax
	int v4;       // eax
	double v5;    // st7
	int v6;       // eax
	uint32_t* v7; // esi
	int v8;       // eax
	uint32_t* v9; // eax
	float v10;    // [esp+8h] [ebp-8h]
	float v11;    // [esp+Ch] [ebp-4h]

	v1 = nox_xxx_newObjectByTypeID_4E3810("Necromancer");
	v10 = sub_4E40F0();
	if (v1) {
		v2 = (uint32_t*)v1[187];
		v3 = v2[121];
		if (v3) {
			v4 = *(uint32_t*)(v3 + 72);
		} else {
			v4 = *(unsigned short*)((int)nox_xxx_objectTypeByIndHealthData(*((unsigned short*)v1 + 2)) + 4);
		}
		if (v10 < 1.0) {
			v10 = 1.0;
		}
		v5 = (double)v4 * v10;
		v11 = v5;
		nox_xxx_unitSetHP_4E4560((int)v1, (long long)v5);
		*(uint16_t*)(v1[139] + 4) = nox_float2int(v11);
		if (!*(uint16_t*)v1[139]) {
			nox_xxx_unitSetHP_4E4560((int)v1, 1u);
		}
		v6 = v1[139];
		if (!*(uint16_t*)(v6 + 4)) {
			*(uint16_t*)(v6 + 4) = 1;
		}
		v2[340] = 4;
		v2[411] = 0x10000000;
		v2[423] = 0x10000000;
		v2[326] = 1062501089;
		v2[510] = 1;
		v2[410] = 0x8000000;
		v2[444] = 0x20000000;
		v2[415] = 0x40000000;
		nox_xxx_createAt_4DAA50((int)v1, 0, *(float*)a1, *((float*)a1 + 1));
		v7 = nox_xxx_newObjectByTypeID_4E3810("RewardMarker");
		if (v7) {
			v8 = nox_game_getQuestStage_4E3CC0();
			v9 = nox_server_rewardgen_activateMarker_4F0720((int)v7, v8 + 2);
			if (v9) {
				nox_xxx_inventoryPutImpl_4F3070((int)v1, (int)v9, 0);
			}
			nox_xxx_objectFreeMem_4E38A0((int)v7);
		}
	}
}

//----- (0051A930) --------------------------------------------------------
int nox_xxx_getQuestStage_51A930() { return *getMemU32Ptr(0x5D4594, 2388660); }

//----- (0051A940) --------------------------------------------------------
int sub_51A940(int a1) {
	int result; // eax

	result = a1;
	*getMemU32Ptr(0x5D4594, 2388656) = a1;
	return result;
}

//----- (0051A950) --------------------------------------------------------
int sub_51A950() { return *getMemU32Ptr(0x5D4594, 2388656); }




//----- (0051D0E0) --------------------------------------------------------
void sub_51D0E0() { dword_5d4594_2487244 = 0; }

//----- (0051D0F0) --------------------------------------------------------
int sub_51D0F0(char a1) {
	*getMemU8Ptr(0x973F18, 35972) = a1;
	return 1;
}

//----- (0051D100) --------------------------------------------------------
int sub_51D100(int a1) {
	if (a1 != 1 && a1) {
		return 0;
	}
	*getMemU32Ptr(0x973F18, 35976) = a1;
	return 1;
}

//----- (0051D120) --------------------------------------------------------
nox_waypoint_t* nox_xxx_waypointNew_5798F0(float a1, float a2);
uint32_t* sub_51D120(float* a1) {
	uint32_t* result; // eax
	uint32_t* v2;     // esi
	float2 v3;        // [esp+4h] [ebp-8h]

	result = (float*)nox_xxx_mapGenFixCoords_4D3D90((float2*)a1, &v3);
	if (result) {
		result = nox_xxx_waypointNew_5798F0(v3.field_0, v3.field_4);
		v2 = result;
		if (result) {
			if (dword_5d4594_2487244) {
				if (*getMemU32Ptr(0x973F18, 35976) == 1) {
					sub_51D300(*(int*)&dword_5d4594_2487244, (int)result, getMemByte(0x973F18, 35972));
					sub_51D300((int)v2, *(int*)&dword_5d4594_2487244, getMemByte(0x973F18, 35972));
				}
			}
			dword_5d4594_2487244 = v2;
			result = v2;
		}
	}
	return result;
}

//----- (0051D1A0) --------------------------------------------------------
nox_waypoint_t* sub_579AD0(float a1, float a2);
float* sub_51D1A0(float2* a1) {
	float* result; // eax
	float2 a2;     // [esp+0h] [ebp-8h]

	result = (float*)nox_xxx_mapGenFixCoords_4D3D90(a1, &a2);
	if (result) {
		result = sub_579AD0(a2.field_0, a2.field_4);
	}
	return result;
}

//----- (0051D3F0) --------------------------------------------------------
float2* sub_51D3F0(float2* a1, float2* a2) {
	float2* result; // eax
	float2* v3;     // esi

	result = a1;
	if (a1) {
		if (a2) {
			result = (float2*)sub_51D1A0(a1);
			v3 = result;
			if (result) {
				result = (float2*)sub_51D1A0(a2);
				if (result) {
					result = (float2*)sub_51D2C0((int)v3, (int)result);
				}
			}
		} else {
			result = 0;
		}
	}
	return result;
}
