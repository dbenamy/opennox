#include <math.h>

#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME1_2.h"
#include "GAME1_3.h"
#include "GAME2.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME5.h"
#include "GAME5_2.h"
#include "client__system__parsecmd.h"
#include "common__log.h"
#include "common__object__armrlook.h"
#include "common__object__weaplook.h"
#include "common__system__team.h"

#include "client__gui__servopts__playrlst.h"
#include "client__system__ctrlevnt.h"
#include "client__video__draw_common.h"
#include "common__binfile.h"
#include "common__magic__speltree.h"
#include "common__net_list.h"
#include "common__random.h"
#include "common__strman.h"

#include "client__gui__window.h"
#include "client__io__win95__focus.h"
#include "common/fs/nox_fs.h"
#include "input_common.h"
#include "operators.h"

extern uint32_t dword_5d4594_527656;
extern uint32_t dword_5d4594_251712;
extern uint32_t dword_5d4594_251708;
extern uint32_t dword_5d4594_251716;
extern uint32_t dword_5d4594_251720;
extern uint32_t dword_5d4594_251744;
extern uint32_t dword_5d4594_3484;
extern void* dword_5d4594_251560;
extern uint32_t nox_tile_def_cnt;
extern uint32_t dword_5d4594_251572;
extern uint32_t nox_player_netCode_85319C;

int nox_server_gameSettingsUpdated; // If you define it as 1-byte bool, the game will crash

extern obj_5D4594_2650668_t** ptr_5D4594_2650668;
extern int ptr_5D4594_2650668_cap;

uint32_t nox_tile_def_cnt = 0;
nox_tileDef_t nox_tile_defs_arr[176] = {0};



uint8_t* nox_xxx_doorAttachWall_410360(int a1, int a2, int a3) {
	uint8_t* result; // eax
	char v4;         // cl

	result = nox_xxx_wallCreateAt_410250(a2, a3);
	if (result) {
		v4 = result[4] | 0x10;
		*((uint32_t*)result + 7) = a1;
		result[4] = v4;
	}
	return result;
}

uint32_t* sub_410390(int a1, int a2, int a3) {
	uint32_t* v3;     // esi
	uint32_t* result; // eax
	int v5;           // edx
	char v6;          // cl
	int v7[2];        // [esp+Ch] [ebp-8h]

	v3 = (uint32_t*)nox_xxx_wall_4105E0(a2, a3);
	if (v3 || (result = nox_xxx_wallCreateAt_410250(a2, a3), (v3 = result) != 0)) {
		v5 = *(uint32_t*)(a1 + 16);
		v7[0] = *(uint32_t*)(a1 + 12);
		v6 = *((uint8_t*)v3 + 4);
		v3[8] = a1;
		v7[1] = v5;
		*((uint8_t*)v3 + 4) = v6 | 0x10;
		result = nox_xxx_polygonIsPlayerInPolygon_4217B0((int2*)v7, 0);
		if (result || (result = sub_421990((int2*)v7, 10.0, 0)) != 0) {
			*((uint8_t*)v3 + 8) = *((uint8_t*)result + 130);
		} else {
			*((uint8_t*)v3 + 8) = 1;
		}
	}
	return result;
}

int sub_410550(short a1) {
	int* v1; // eax

	v1 = (int*)nox_xxx_wallSecretGetFirstWall_410780();
	if (!v1) {
		return 0;
	}
	while (*(uint16_t*)(v1[3] + 10) != a1) {
		v1 = (int*)nox_xxx_wallSecretNext_410790(v1);
		if (!v1) {
			return 0;
		}
	}
	return v1[3];
}

uint32_t* sub_410730() {
	uint32_t* result; // eax
	uint32_t* v1;     // esi

	result = dword_5d4594_251560;
	if (dword_5d4594_251560) {
		do {
			v1 = (uint32_t*)*result;
			free(result);
			result = v1;
		} while (v1);
		dword_5d4594_251560 = 0;
	} else {
		dword_5d4594_251560 = 0;
	}
	return result;
}

uint32_t* nox_xxx_wallSecretBlock_410760(uint32_t* a1) {
	uint32_t* result; // eax

	result = a1;
	*a1 = dword_5d4594_251560;
	dword_5d4594_251560 = a1;
	return result;
}

void* nox_xxx_wallSecretGetFirstWall_410780() { return dword_5d4594_251560; }

int nox_xxx_wallSecretNext_410790(int* a1) {
	int result; // eax

	if (a1) {
		result = *a1;
	} else {
		result = 0;
	}
	return result;
}

int* sub_4107A0(void* lpMem) {
	int* result; // eax
	int* v2;     // esi

	result = dword_5d4594_251560;
	v2 = 0;
	if (dword_5d4594_251560) {
		while (result != lpMem) {
			v2 = result;
			result = (int*)nox_xxx_wallSecretNext_410790(result);
			if (!result) {
				return result;
			}
		}
		if (result == dword_5d4594_251560) {
			dword_5d4594_251560 = nox_xxx_wallSecretNext_410790(result);
		} else {
			*v2 = nox_xxx_wallSecretNext_410790(result);
		}
		free(lpMem);
	}
	return result;
}

int nox_xxx_tileAlloc_410F60_init() {
	ptr_5D4594_2650668 = calloc(ptr_5D4594_2650668_cap, sizeof(void*));
	if (!ptr_5D4594_2650668) {
		return 0;
	}

	for (int i = 0; i < ptr_5D4594_2650668_cap; i++) {
		ptr_5D4594_2650668[i] = (obj_5D4594_2650668_t*)calloc(ptr_5D4594_2650668_cap, sizeof(obj_5D4594_2650668_t));
		if (!ptr_5D4594_2650668[i]) {
			return 0;
		}
	}
	return 1;
}

void nox_xxx_tileFree_410FC0_free() {
	for (int i = 0; i < ptr_5D4594_2650668_cap; i++) {
		obj_5D4594_2650668_t* ptr = ptr_5D4594_2650668[i];
		if (ptr) {
			free(ptr);
		}
	}
}

int nox_xxx_tileNFromPoint_411160(float2* a1) {
	float v12 = (a1->field_0 + 11.5) * 0.021739131;
	float v13 = (a1->field_4 + 11.5) * 0.021739131;

	int i = nox_float2int(v12);
	int j = nox_float2int(v13);

	float v14 = a1->field_0 + 11.5;
	float v15 = a1->field_4 + 11.5;

	int v4 = nox_float2int(v14) % 46;
	int v5 = nox_float2int(v15) % 46;

	if (i <= 1 || i >= 127 || j <= 1 || j >= 127) {
		return -1;
	}

	int result = 0;
	int v16[2] = {0};
	if (v4 <= v5) {
		if (46 - v4 <= v5) {
			obj_5D4594_2650668_t* obj = &ptr_5D4594_2650668[i][j];
			result = obj->field_6;
			if (obj->field_10) {
				v16[0] = v4;
				v16[1] = v5 - 23;
				result = sub_411350(obj->field_10, v16, result);
			}
		} else {
			obj_5D4594_2650668_t* obj = &ptr_5D4594_2650668[i - 1][j];
			result = obj->field_1;
			if (obj->field_5) {
				v16[1] = v5;
				v16[0] = v4 + 23;
				result = sub_411350(obj->field_5, v16, result);
			}
		}
	} else if (46 - v4 <= v5) {
		obj_5D4594_2650668_t* obj = &ptr_5D4594_2650668[i][j];
		result = obj->field_1;
		if (obj->field_5) {
			v16[1] = v5;
			v16[0] = v4 - 23;
			result = sub_411350(obj->field_5, v16, result);
		}
	} else {
		obj_5D4594_2650668_t* obj = &ptr_5D4594_2650668[i][j - 1];
		result = obj->field_6;
		if (obj->field_10) {
			v16[0] = v4;
			v16[1] = v5 + 23;
			result = sub_411350(obj->field_10, v16, result);
		}
	}
	return result;
}

int nox_xxx_mapTileAllowTeleport_411A90(float2* a1) {
	if (*getMemIntPtr(0x587000, 26520) == -1) {
		for (int i = 0; i < 176; i++) {
			nox_tileDef_t* p = &nox_tile_defs_arr[i];
			if (!strcmp(&p->name[0], "WaterNoTeleport")) {
				*getMemU32Ptr(0x587000, 26516) = i;
			} else if (!strcmp(&p->name[0], "WaterDeepNoTeleport")) {
				*getMemU32Ptr(0x587000, 26520) = i;
			} else if (!strcmp(&p->name[0], "WaterShallowNoTeleport")) {
				*getMemU32Ptr(0x587000, 26524) = i;
			} else if (!strcmp(&p->name[0], "WaterSwampDeepNoTeleport")) {
				*getMemU32Ptr(0x587000, 26528) = i;
			} else if (!strcmp(&p->name[0], "WaterSwampShallowNoTeleport")) {
				*getMemU32Ptr(0x587000, 26532) = i;
			}
		}
	}
	int v3 = nox_xxx_tileNFromPoint_411160(a1);
	return v3 == *getMemU32Ptr(0x587000, 26516) || v3 == *getMemU32Ptr(0x587000, 26520) ||
		   v3 == *getMemU32Ptr(0x587000, 26524) || v3 == *getMemU32Ptr(0x587000, 26528) ||
		   v3 == *getMemU32Ptr(0x587000, 26532);
}

obj_412ae0_t* nox_xxx_modifNext_4133C0(obj_412ae0_t* a1) { return a1->field_34; }

int sub_4133D0(nox_object_t* a1p) {
	int a1 = a1p;
	int v1; // eax
	int v2; // esi
	int v3; // eax

	v1 = *getMemU32Ptr(0x5D4594, 251620);
	v2 = *(uint32_t*)(a1 + 692);
	if (!*getMemU32Ptr(0x5D4594, 251620)) {
		v3 = nox_xxx_modifGetIdByName_413290("Material7");
		v1 = nox_xxx_modifGetDescById_413330(v3);
		*getMemU32Ptr(0x5D4594, 251620) = v1;
	}
	return *(uint32_t*)(a1 + 8) & 0x13001000 && *(uint32_t*)(v2 + 4) == v1;
}

int sub_413420(char a1) {
	unsigned char* v1; // esi
	int v2;            // ecx
	unsigned char* v3; // eax

	if (!*getMemU32Ptr(0x5D4594, 251624)) {
		v1 = getMemAt(0x587000, 27340);
		do {
			*(uint32_t*)v1 = nox_xxx_gLoadImg_42F970(*((const char**)v1 - 1));
			v1 += 20;
		} while ((int)v1 < (int)getMemAt(0x587000, 27460));
		*getMemU32Ptr(0x5D4594, 251624) = 1;
	}
	v2 = 0;
	v3 = getMemAt(0x587000, 27332);
	while (*v3 != a1) {
		v3 += 20;
		++v2;
		if ((int)v3 >= (int)getMemAt(0x587000, 27452)) {
			return 0;
		}
	}
	return *getMemU32Ptr(0x587000, 27340 + 20 * v2);
}

int sub_4134D0() {
	dword_5d4594_251708 = 0;
	dword_5d4594_251712 = 0;
	dword_5d4594_251716 = 0;
	dword_5d4594_251720 = 0;
	return 0;
}

void sub_4137E0() {
	if (dword_5d4594_251720) {
		dword_5d4594_251720 = 0;
		sub_43DBE0();
	}
}

char* sub_413890() {
	unsigned char* v0; // edi
	int v1;            // ecx
	short v2;          // dx
	unsigned char v3;  // al

	*getMemU8Ptr(0x5D4594, 251636) = 0;
	*getMemU8Ptr(0x5D4594, 251637) = 0;
	if (!getMemByte(0x5D4594, 251636)) {
		return 0;
	}
	v0 = getMemAt(0x5D4594, 251637 + strlen((const char*)getMemAt(0x5D4594, 251636)));
	v1 = *getMemU32Ptr(0x587000, 32316);
	v2 = *getMemU16Ptr(0x587000, 32320);
	*(uint32_t*)--v0 = *getMemU32Ptr(0x587000, 32312);
	v3 = getMemByte(0x587000, 32322);
	*((uint32_t*)v0 + 1) = v1;
	*((uint16_t*)v0 + 4) = v2;
	v0[10] = v3;
	return (char*)getMemAt(0x5D4594, 251636);
}

void sub_4138E0(int a1) {
	*getMemU32Ptr(0x5D4594, 251740) = nox_xxx_checkGameFlagPause_413A50();
	sub_413A00(1);
}

void sub_413900(int a1) {
	if (!nox_video_inFadeTransition_44E0D0()) {
		if (!*getMemU32Ptr(0x5D4594, 251740)) {
			sub_413A00(0);
		}
	}
}

int sub_413920() {
	sub_42EBB0(1u, sub_413900, 0, "Pause");
	sub_42EBB0(2u, sub_4138E0, 0, "Pause");
	dword_5d4594_251744 = 0;
	return 1;
}

int sub_4139B0() { return dword_5d4594_251744 != 0; }

int nox_xxx_checkGameFlagPause_413A50() { return nox_common_gameFlags_check_40A5C0(0x40000); }

void nox_xxx_gameLoopMemDump_413E30() {
	signed int v0;     // ebx
	int v1;            // esi
	unsigned char* v2; // eax
	int v3;            // edx
	signed int v4;     // edi
	const char* v5;    // ebp
	unsigned char* v6; // eax
	signed int v7;     // ecx
	int v8;            // ebx
	signed int v9;     // [esp+8h] [ebp-8h]

	v0 = 0;
	v1 = 0;
	v9 = 0;
	v2 = getMemAt(0x5D4594, 252364);
	do {
		*(uint32_t*)v2 = 0;
		v2 += 84;
	} while ((int)v2 < (int)getMemAt(0x5D4594, 338380));
	v3 = *getMemU32Ptr(0x5D4594, 338304);
	if (*getMemU32Ptr(0x5D4594, 338304)) {
		do {
			v4 = 0;
			if (v0 > 0) {
				v5 = (const char*)getMemAt(0x5D4594, 252284);
				while (strcmp(v5, (const char*)(v3 + 20))) {
					v0 = v9;
					++v4;
					v5 += 84;
					if (v4 >= v9) {
						goto LABEL_10;
					}
				}
				v0 = v9;
				*getMemU32Ptr(0x5D4594, 252364 + 84 * v4) += *(uint32_t*)(v3 + 16);
			}
		LABEL_10:
			if (v4 == v0) {
				*getMemU32Ptr(0x5D4594, 252364 + 84 * v4) = *(uint32_t*)(v3 + 16);
				strcpy((char*)getMemAt(0x5D4594, 252284 + 84 * v4), (const char*)(v3 + 20));
				v9 = ++v0;
			}
			v3 = *(uint32_t*)(v3 + 4);
		} while (v3);
		v1 = 0;
	}
	qsort(getMemAt(0x5D4594, 252284), v0, 0x54u, sub_413F60);
	if (v0 > 0) {
		v6 = getMemAt(0x5D4594, 252364);
		v7 = v0;
		do {
			v8 = *(uint32_t*)v6;
			v6 += 84;
			v1 += v8;
			--v7;
		} while (v7);
	}
}

int sub_413F60(const void* a1, const void* a2) { return *((uint32_t*)a1 + 20) - *((uint32_t*)a2 + 20); }

char nox_xxx_initSinCosTables_414C90() {
	long long v0;      // rax
	unsigned char* v1; // esi
	unsigned char* v2; // esi
	int v4;            // [esp+0h] [ebp-4h]
	int v5;            // [esp+0h] [ebp-4h]

	LOBYTE(v0) = getMemByte(0x5D4594, 371240);
	if (!getMemByte(0x5D4594, 371240)) {
		*getMemU8Ptr(0x5D4594, 371240) = 1;
		v4 = 0;
		v1 = getMemAt(0x85B3FC, 12260);
		do {
			*(uint32_t*)v1 = (long long)(sin((double)v4 * 0.0015339808) * *getMemDoublePtr(0x581450, 7200));
			v1 += 4;
			++v4;
		} while ((int)v1 < (int)getMemAt(0x85B3FC, 28644));
		v5 = 0;
		v2 = getMemAt(0x5D4594, 338472);
		do {
			v0 = (long long)(acos((double)v5 * 0.00024414062 - 1.0) * *getMemDoublePtr(0x581450, 7184));
			*(uint32_t*)v2 = v0;
			v2 += 4;
			++v5;
		} while ((int)v2 < (int)getMemAt(0x5D4594, 371240));
	}
	return v0;
}

int sub_415960(wchar2_t* a1) {
	int v1;             // edi
	const wchar2_t** v2; // eax
	unsigned char* v3;  // esi
	int v4;             // ecx

	v1 = 0;
	if (!*getMemU32Ptr(0x587000, 33392)) {
		return 0;
	}
	v2 = (const wchar2_t**)getMemAt(0x587000, 33392);
	v3 = getMemAt(0x587000, 33392);
	while (_nox_wcsicmp(a1, *v2)) {
		v4 = *((uint32_t*)v3 + 3);
		v3 += 12;
		++v1;
		v2 = (const wchar2_t**)v3;
		if (!v4) {
			return 0;
		}
	}
	return *getMemU32Ptr(0x587000, 33400 + 12 * v1);
}

int sub_4159F0(int a1) {
	int v1;           // ecx
	unsigned char* i; // eax
	int v3;           // esi

	v1 = 0;
	if (!*getMemU32Ptr(0x587000, 33392)) {
		return 0;
	}
	for (i = getMemAt(0x587000, 33392); *((uint32_t*)i + 2) != a1; i += 12) {
		v3 = *((uint32_t*)i + 3);
		++v1;
		if (!v3) {
			return 0;
		}
	}
	return *getMemU32Ptr(0x587000, 33392 + 12 * v1);
}

double sub_415BD0(int a1) {
	float* v1;     // eax
	double result; // st7

	if (*(uint32_t*)(a1 + 8) & 0x2000000 &&
		(v1 = (float*)nox_xxx_equipClothFindDefByTT_413270(*(unsigned short*)(a1 + 4))) != 0) {
		result = v1[16];
	} else {
		result = 0.0;
	}
	return result;
}

int sub_415DA0(wchar2_t* a1) {
	int v1;             // edi
	const wchar2_t** v2; // eax
	unsigned char* v3;  // esi
	int v4;             // ecx

	v1 = 0;
	if (!*getMemU32Ptr(0x587000, 35496)) {
		return 0;
	}
	v2 = (const wchar2_t**)getMemAt(0x587000, 35496);
	v3 = getMemAt(0x587000, 35496);
	while (_nox_wcsicmp(a1, *v2)) {
		v4 = *((uint32_t*)v3 + 3);
		v3 += 12;
		++v1;
		v2 = (const wchar2_t**)v3;
		if (!v4) {
			return 0;
		}
	}
	return *getMemU32Ptr(0x587000, 35504 + 12 * v1);
}

int sub_415E80(int a1) {
	int v1;           // ecx
	unsigned char* i; // eax
	int v3;           // esi

	v1 = 0;
	if (!*getMemU32Ptr(0x587000, 35496)) {
		return 0;
	}
	for (i = getMemAt(0x587000, 35496); *((uint32_t*)i + 2) != a1; i += 12) {
		v3 = *((uint32_t*)i + 3);
		++v1;
		if (!v3) {
			return 0;
		}
	}
	return *getMemU32Ptr(0x587000, 35496 + 12 * v1);
}
