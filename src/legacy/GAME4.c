#include <errno.h>
#include <math.h>

#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME1_2.h"
#include "GAME2.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_1.h"
#include "GAME4_2.h"
#include "GAME4_3.h"
#include "GAME5_2.h"
#include "common__binfile.h"
#include "common__crypt.h"
#include "common__net_list.h"
#include "common__random.h"
#include "common__system__team.h"
#include "operators.h"
#include "server__magic__plyrspel.h"
#include "server__magic__spell__execdur.h"
#include "server__script__file.h"
#include "server__script__script.h"

#include "client__gui__window.h"
#include "client__video__draw_common.h"
#include "common/fs/nox_fs.h"
#include "common__magic__speltree.h"

extern uint32_t dword_5d4594_3835368;
extern uint32_t nox_server_kickQuestPlayerMinVotes_229992;
extern uint32_t nox_server_resetQuestMinVotes_229988;
extern uint32_t dword_5d4594_1599644;
extern uint32_t dword_5d4594_3835392;
extern uint32_t dword_5d4594_3835312;
extern uint32_t dword_5d4594_1568868;
extern void* nox_alloc_magicEnt_1569668;
extern void* nox_alloc_vote_1599652;
extern uint32_t dword_5d4594_1599616;
extern uint32_t dword_5d4594_1569672;
extern uint32_t dword_5d4594_3835396;
extern uint32_t dword_5d4594_1599596;
extern uint32_t dword_5d4594_1599576;
extern uint32_t dword_5d4594_1599656;
extern uint32_t dword_5d4594_2650652;


FILE* nox_file_8 = 0;

int nox_cheat_charmall = 0;
uint32_t dword_5d4594_1599480 = 0;
uint32_t dword_5d4594_1599476 = 0;
void* dword_5d4594_1599540 = 0;
void* dword_5d4594_1599532 = 0;
void* dword_5d4594_1599556 = 0;
void* dword_5d4594_1599548 = 0;
void* dword_5d4594_1599588 = 0;
void* dword_5d4594_1599592 = 0;

int nox_setImaginaryCaster();
int sub_57AEE0(int a1, nox_object_t* a2);
//----- (00500540) --------------------------------------------------------


//----- (005005E0) --------------------------------------------------------


//----- (005006B0) --------------------------------------------------------


//----- (00500750) --------------------------------------------------------


//----- (00500770) --------------------------------------------------------


//----- (00500790) --------------------------------------------------------


//----- (005007E0) --------------------------------------------------------


//----- (005009B0) --------------------------------------------------------


//----- (00500A60) --------------------------------------------------------


//----- (00500B70) --------------------------------------------------------

// 500B70: using guessed type char var_100[256];

//----- (00500C70) --------------------------------------------------------
// Sends information to the player that an unit order happened
int nox_xxx_orderUnitLocal_500C70(int owner, int orderType) {
	*((uint32_t*)nox_common_playerInfoFromNum_417090(owner) + 912) = orderType;
	return nox_xxx_netCreatureCmd_4D7EE0(owner, orderType);
}

nox_object_t* nox_xxx_unitDoSummonAt_5016C0(int a1, float* a2, nox_object_t* a3, unsigned char a4);
//----- (00502670) --------------------------------------------------------
void nox_server_scriptExecuteFnForEachGroupObj_502670(unsigned char* groupPtr, int expectedType, void (*a3)(int, int),
													  int a4) {
	int* i;             // esi
	int v5;             // eax
	int* j;             // esi
	uint32_t* v7;       // eax
	int* k;             // esi
	int v9;             // eax
	int* l;             // esi
	unsigned char* v11; // eax

	if (!groupPtr) {
		return;
	}
	switch (*groupPtr) {
	case 0u:
		if (expectedType != 0) {
			break;
		}
		for (i = (int*)*((uint32_t*)groupPtr + 21); i; i = (int*)i[2]) {
			v5 = nox_xxx_netGetUnitByExtent_4ED020(*i);
			if (v5) {
				a3(v5, a4);
			}
		}
		break;
	case 1u:
		if (expectedType != 1) {
			break;
		}
		for (j = (int*)*((uint32_t*)groupPtr + 21); j; j = (int*)j[2]) {
			v7 = nox_server_getWaypointById_579C40(*j);
			if (v7) {
				a3((int)v7, a4);
			}
		}
		break;
	case 2u:
		if (expectedType != 2) {
			break;
		}
		for (k = (int*)*((uint32_t*)groupPtr + 21); k; k = (int*)k[2]) {
			v9 = nox_server_getWallAtGrid_410580(*k, k[1]);
			if (v9) {
				a3(v9, a4);
			}
		}
		// fallthrough
	case 3u:
		for (l = (int*)*((uint32_t*)groupPtr + 21); l; l = (int*)l[2]) {
			v11 = (unsigned char*)nox_server_scriptGetGroup_57C0A0(*l);
			if (v11) {
				nox_server_scriptExecuteFnForEachGroupObj_502670(v11, expectedType, a3, a4);
			}
		}
		break;
	default:
		break;
	}
}

//----- (00502790) --------------------------------------------------------
int nox_xxx_mapgenMakeScript_502790(FILE* a1, char* a2) {
	int result;     // eax
	int i;          // ebx
	int v4;         // edi
	int v5;         // eax
	int v6;         // [esp+8h] [ebp-410h]
	int v7;         // [esp+Ch] [ebp-40Ch]
	int v8;         // [esp+10h] [ebp-408h]
	int v9;         // [esp+14h] [ebp-404h]
	char v10[1024]; // [esp+18h] [ebp-400h]

	nox_binfile_fread_408E40((char*)&v8, 4, 1, a1);
	nox_binfile_fread_408E40(v10, 1, v8, a1);
	nox_binfile_fread_408E40(a2, 4, 1, a1);
	nox_binfile_fread_408E40((char*)&v7, 4, 1, a1);
	result = v7;
	for (i = 0; i < v7; ++i) {
		nox_binfile_fread_408E40((char*)&v6, 1, 1, a1);
		nox_binfile_fseek_409050(a1, 1, SEEK_CUR);
		v4 = 0;
		v5 = 268 * (unsigned char)v6;
		if (getMemByte(0x587000, 218640 + v5)) {
			do {
				switch (*getMemU32Ptr(0x587000, 218648 + 8 * v4 + v5)) {
				case 0:
				case 3:
				case 4:
				case 5:
				case 6:
					nox_binfile_fseek_409050(a1, 4, SEEK_CUR);
					break;
				case 1:
					nox_binfile_fseek_409050(a1, 8, SEEK_CUR);
					break;
				case 2:
				case 7:
					nox_binfile_fread_408E40((char*)&v9, 1, 1, a1);
					nox_binfile_fseek_409050(a1, (unsigned char)v9, SEEK_CUR);
					break;
				default:
					break;
				}
				++v4;
				v5 = 268 * (unsigned char)v6;
			} while (v4 < getMemByte(0x587000, 218640 + v5));
		}
		result = v7;
	}
	return result;
}

//----- (005029A0) --------------------------------------------------------
int sub_5029A0(char* a1) {
	int v1; // edi
	int i;  // esi

	v1 = 0;
	if (*(int*)&dword_5d4594_1599596 <= 0) {
		return -1;
	}
	for (i = 0; nox_strcmpi(a1, (const char*)(i + dword_5d4594_1599576)); i += 76) {
		if (++v1 >= *(int*)&dword_5d4594_1599596) {
			return -1;
		}
	}
	return v1;
}

//----- (005029F0) --------------------------------------------------------
int sub_5029F0(int a1) {
	int result; // eax

	if (a1 < 0 || a1 > *(int*)&dword_5d4594_1599596) {
		result = 0;
	} else {
		result = dword_5d4594_1599576 + 76 * a1;
	}
	return result;
}

//----- (00502A20) --------------------------------------------------------
int sub_502A20() { return dword_5d4594_1599596; }

//----- (00502A50) --------------------------------------------------------
int sub_502A50(char* a1) {
	int result; // eax

	sub_502DF0();
	if (a1) {
		strncpy(dword_5d4594_1599588, a1, 0x7FFu);
		result = 1;
	} else {
		**(uint8_t**)&dword_5d4594_1599588 = getMemByte(0x5D4594, 1599608);
		result = 0;
	}
	return result;
}

//----- (00502AB0) --------------------------------------------------------
int sub_502AB0(char* a1) {
	int result; // eax

	if (a1) {
		strncpy(dword_5d4594_1599592, a1, 0x7FFu);
		result = 1;
	} else {
		**(uint8_t**)&dword_5d4594_1599592 = getMemByte(0x5D4594, 1599612);
		result = 0;
	}
	return result;
}

//----- (00502B10) --------------------------------------------------------
int sub_502B10() {
	int result;   // eax
	int v1;       // ebp
	int v2;       // ecx
	int v3;       // ebp
	char v4;      // [esp+12h] [ebp-56h]
	char v5;      // [esp+13h] [ebp-55h]
	int v6;       // [esp+14h] [ebp-54h]
	int v7;       // [esp+18h] [ebp-50h]
	int v8;       // [esp+1Ch] [ebp-4Ch]
	float v9;     // [esp+20h] [ebp-48h]
	float v10;    // [esp+24h] [ebp-44h]
	char v11[64]; // [esp+28h] [ebp-40h]

	dword_5d4594_1599596 = 0;
	if (!dword_5d4594_1599588) {
		dword_5d4594_1599588 = calloc(1, 0x800u);
	}
	if (!dword_5d4594_1599592) {
		dword_5d4594_1599592 = calloc(1, 0x800u);
	}
	if (!dword_5d4594_1599576) {
		dword_5d4594_1599576 = calloc(1, 0x26000u);
	}
	result = 0;
	if (strlen(dword_5d4594_1599588)) {
		result = sub_502DA0(dword_5d4594_1599588);
		if (result) {
			nox_fs_fread(nox_file_8, &v8, 4);
			if (v8 == -889266515) {
				while (1) {
					v6 = 0;
					nox_fs_fread(nox_file_8, &v6, 4);
					v1 = v6;
					if (!v6) {
						break;
					}
					if (dword_5d4594_1599596 >= 2048) {
						sub_502DF0();
						return 0;
					}
					*(uint32_t*)(dword_5d4594_1599576 + 76 * dword_5d4594_1599596 + 72) = nox_fs_ftell(nox_file_8) - 4;
					nox_fs_fread(nox_file_8, &v7, 1);
					nox_fs_fread(nox_file_8, v11, (unsigned char)v7);
					v2 = -1 - (unsigned char)v7;
					v11[(unsigned char)v7] = 0;
					v3 = v2 + v1;
					strcpy((char*)(dword_5d4594_1599576 + 76 * dword_5d4594_1599596), v11);
					nox_fs_fread(nox_file_8, &v4, 1);
					nox_fs_fread(nox_file_8, &v5, 1);
					nox_fs_fread(nox_file_8, &v9, 4);
					nox_fs_fread(nox_file_8, &v10, 4);
					*(float*)(dword_5d4594_1599576 + 76 * dword_5d4594_1599596 + 64) = v9;
					*(float*)(dword_5d4594_1599576 + 76 * (dword_5d4594_1599596)++ + 68) = v10;
					nox_fs_fseek_cur(nox_file_8, v3 - 10);
				}
				sub_502DF0();
				return 1;
			} else {
				sub_502DF0();
				return 0;
			}
		}
	}
	return result;
}
// 502B10: using guessed type char var_40[64];

//----- (00502D70) --------------------------------------------------------
int sub_502D70(int a1) {
	if (a1 < 0 || a1 >= *(int*)&dword_5d4594_1599596) {
		return 0;
	}
	dword_5d4594_3835396 = a1;
	return nox_xxx_mapgenSaveMap_503830(a1) == 0;
}

//----- (00502DA0) --------------------------------------------------------
FILE* sub_502DA0(char* a1) {
	FILE* result; // eax

	result = nox_file_8;
	if (nox_file_8 || (result = (FILE*)nox_xxx_cryptOpen_426910(a1, 1, -1)) != 0 &&
						  (result = nox_xxx_mapgenGetSomeFile_426A60(), (nox_file_8 = result) != 0)) {
		nox_fs_fseek_start(result, 0);
		result = (FILE*)1;
	}
	return result;
}

//----- (00502DF0) --------------------------------------------------------
FILE* sub_502DF0() {
	FILE* result; // eax

	result = nox_file_8;
	if (nox_file_8) {
		nox_xxx_cryptClose_4269F0();
		nox_file_8 = 0;
	}
	return result;
}

//----- (00502E10) --------------------------------------------------------
FILE* sub_502E10(int a1) {
	if (!nox_file_8 || a1 < 0 || a1 >= *(int*)&dword_5d4594_1599596) {
		return 0;
	}
	nox_fs_fseek_start(nox_file_8, *(uint32_t*)(dword_5d4594_1599576 + 76 * a1 + 72));
	return nox_file_8;
}

//----- (00502E70) --------------------------------------------------------
double sub_502E70(int a1) {
	double result; // st7

	if (a1 < 0 || a1 >= *(int*)&dword_5d4594_1599596) {
		result = -1.0;
	} else {
		result = *(float*)(dword_5d4594_1599576 + 76 * a1 + 64);
	}
	return result;
}

//----- (00502EA0) --------------------------------------------------------
double sub_502EA0(int a1) {
	double result; // st7

	if (a1 < 0 || a1 >= *(int*)&dword_5d4594_1599596) {
		result = -1.0;
	} else {
		result = *(float*)(dword_5d4594_1599576 + 76 * a1 + 68);
	}
	return result;
}

//----- (00503830) --------------------------------------------------------
int nox_xxx_mapgenSaveMap_503830(int a1) {
	FILE* v1;         // esi
	uint32_t* v2;     // eax
	int v3;           // esi
	unsigned char v5; // [esp+Fh] [ebp-19Dh]
	int v6;           // [esp+10h] [ebp-19Ch]
	char v7;          // [esp+17h] [ebp-195h]
	int v8;           // [esp+18h] [ebp-194h]
	int v9;           // [esp+1Ch] [ebp-190h]
	int v10;          // [esp+20h] [ebp-18Ch]
	int v11[8];       // [esp+24h] [ebp-188h]
	int v19[2];       // [esp+44h] [ebp-168h]
	int v21;          // [esp+4Ch] [ebp-160h]
	int v22;          // [esp+50h] [ebp-15Ch]
	int v23;          // [esp+54h] [ebp-158h]
	char v24[4];      // [esp+58h] [ebp-154h]
	int4 v25;         // [esp+5Ch] [ebp-150h]
	char v26[64];     // [esp+6Ch] [ebp-140h]
	char v27[256];    // [esp+ACh] [ebp-100h]

	if (a1 < 0) {
		return 0;
	}
	if (a1 >= *(int*)&dword_5d4594_1599596) {
		return 0;
	}
	nox_xxx_free_503F40();
	*getMemU32Ptr(0x5D4594, 1599572) = -1;
	dword_5d4594_1599644 = 0;
	sub_502DA0(dword_5d4594_1599588);
	if (!sub_502E10(a1)) {
		return 0;
	}
	v1 = nox_xxx_mapgenGetSomeFile_426A60();
	nox_fs_fread(v1, &v22, 4);
	nox_fs_fread(v1, &v9, 1);
	nox_fs_fread(v1, v26, (unsigned char)v9);
	nox_fs_fread(v1, &v7, 1);
	nox_fs_fread(v1, &v5, 1);
	nox_fs_fread(v1, &v21, 4);
	nox_fs_fread(v1, &v23, 4);
	if (v5 > 1u) {
		nox_fs_fread(v1, &v6, 4);
		nox_fs_fseek_cur(v1, v6);
	}
	nox_fs_fread(v1, &v10, 4);
	if (v10 != -889266515) {
		return 0;
	}
	nox_fs_fread(v1, v19, 4);
	nox_fs_fread(v1, &v19[1], 4);
	nox_xxx_mapWall_426A80(v19);
	nox_fs_fread(v1, v11, 4);
	nox_fs_fread(v1, &v11[1], 4);
	nox_fs_fread(v1, &v11[6], 4);
	nox_fs_fread(v1, &v11[7], 4);
	nox_fs_fread(v1, &v11[2], 4);
	nox_fs_fread(v1, &v11[3], 4);
	nox_fs_fread(v1, &v11[4], 4u);
	nox_fs_fread(v1, &v11[5], 4u);
	sub_4D3C80(v11);
	memcpy(getMemAt(0x5D4594, 1599500), v11, 0x20u);
	sub_428170(&v11, &v25);
	nox_xxx_cryptSetTypeMB_426A50(1);
	while (1) {
		v6 = 0;
		LOBYTE(v8) = 0;
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v8, 1u);
		if (!(uint8_t)v8) {
			break;
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread(v27, (unsigned char)v8);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v24, 4u);
		if (!nox_xxx_mapReadSection_426EA0((int)&v11, v27, &v6)) {
			if (v6 == 1) {
				sub_502DF0();
				return 0;
			}
			v2 = nox_xxx_newObjectByTypeID_4E3810(v27);
			v3 = (int)v2;
			if (!v2) {
				return 0;
			}
			if (!((int (*)(uint32_t*, int4*))v2[176])(v2, &v25)) {
				nox_xxx_objectFreeMem_4E38A0(v3);
				sub_502DF0();
				return 0;
			}
			nox_xxx_servMapLoadPlaceObj_4F3F50(v3, 0, &v25.field_0);
		}
	}
	nox_xxx_cryptSetTypeMB_426A50(0);
	dword_5d4594_1599480 = a1;
	dword_5d4594_1599476 = 0;
	dword_5d4594_3835396 = a1;
	sub_502DF0();
	return 1;
}

//----- (00503B30) --------------------------------------------------------
void nox_script_readWriteZzz_541670(char* path, char* path2, char* dst);
void nox_xxx_waypoint_5799C0();
void sub_579D20();
nox_waypoint_t* sub_579890();
int sub_503B30(float2* a1) {
	int result; // eax
	int v2;     // edi
	double v3;  // st7
	float v4;   // ecx
	char* v5;   // eax
	char* v6;   // ecx
	int v7;     // eax
	int v8;     // esi
	int v9;     // edi
	int i;      // eax
	int j;      // eax
	int k;      // eax
	float2 v13; // [esp+Ch] [ebp-50h]
	float2 v14; // [esp+14h] [ebp-48h]
	float2 a2;  // [esp+1Ch] [ebp-40h]
	int2 v16;   // [esp+24h] [ebp-38h]
	int4 v17;   // [esp+2Ch] [ebp-30h]
	int v18[8]; // [esp+3Ch] [ebp-20h]

	result = nox_xxx_mapGenFixCoords_4D3D90(a1, &a2);
	if (result) {
		v2 = dword_5d4594_3835396;
		if (dword_5d4594_1599480 != dword_5d4594_3835396 || *(int*)&dword_5d4594_1599480 == -1 ||
			dword_5d4594_1599476 == 1) {
			result = nox_xxx_mapgenSaveMap_503830(*(int*)&dword_5d4594_3835396);
			if (!result) {
				return result;
			}
			v2 = dword_5d4594_3835396;
		}
		v18[2] = (long long)a2.field_0;
		v18[3] = (long long)a2.field_4;
		v13.field_0 = *(float*)(dword_5d4594_1599576 + 76 * v2 + 64) + a1->field_0;
		v13.field_4 = *(float*)(dword_5d4594_1599576 + 76 * v2 + 68) + a1->field_4;
		nox_xxx_mapGenFixCoords_4D3D90(&v13, &v14);
		v18[4] = (long long)v14.field_0;
		v18[5] = (long long)v14.field_4;
		v3 = *(float*)(dword_5d4594_1599576 + 76 * dword_5d4594_3835396 + 64) + a1->field_0;
		v13.field_4 = a1->field_4;
		v13.field_0 = v3;
		nox_xxx_mapGenFixCoords_4D3D90(&v13, &v14);
		v18[0] = (long long)v14.field_0;
		v4 = a1->field_0;
		v18[1] = (long long)v14.field_4;
		v13.field_0 = v4;
		v13.field_4 = *(float*)(dword_5d4594_1599576 + 76 * dword_5d4594_3835396 + 68) + a1->field_4;
		nox_xxx_mapGenFixCoords_4D3D90(&v13, &v14);
		v18[6] = (long long)v14.field_0;
		v18[7] = (long long)v14.field_4;
		sub_4D3C80(v18);
		sub_428170(v18, &v17);
		v5 = nox_xxx_mapGetWallSize_426A70();
		v6 = v5;
		v7 = *(uint32_t*)v5;
		*getMemU32Ptr(0x5D4594, 1599484) = v7;
		*getMemU32Ptr(0x5D4594, 1599488) = *((uint32_t*)v6 + 1);
		*getMemFloatPtr(0x5D4594, 1599492) = (double)(23 * v7);
		*getMemFloatPtr(0x5D4594, 1599496) = (double)(int)(23 * *getMemU32Ptr(0x5D4594, 1599488));
		v8 = (long long)(a2.field_0 - (double)*getMemIntPtr(0x5D4594, 1599508));
		v9 = (long long)(a2.field_4 - (double)*getMemIntPtr(0x5D4594, 1599512));
		result = nox_xxx_tileInit_504150(v8, v9);
		if (result) {
			result = sub_504330(v8, v9);
			if (result) {
				result = sub_504560(v8, v9);
				if (result) {
					result = sub_504910(v8, v9);
					if (result) {
						sub_579D20();
						for (i = sub_579890(); i; i = sub_5798A0(i)) {
							*(uint32_t*)(i + 480) |= 0x80000000;
						}
						dword_5d4594_3835392 = nox_xxx_interesting_xfer_4D0010(&v17, *(int*)&dword_5d4594_3835392);
						result = sub_504720(v8, v9);
						if (result) {
							for (j = sub_579890(); j; j = sub_5798A0(j)) {
								*(uint32_t*)(j + 4) = 0;
							}
							for (k = nox_server_getFirstObjectUninited_4DA870(); k;
								 k = nox_server_getNextObjectUninited_4DA880(k)) {
								*(uint32_t*)(k + 44) = 0;
							}
							nox_xxx_waypoint_5799C0();
							nox_xxx_unitClearPendingMB_4DB030();
							dword_5d4594_1599476 = 1;
							if (dword_5d4594_1599644) {
								++*getMemU32Ptr(0x973F18, 35880);
								sub_542BF0(*(int*)&dword_5d4594_3835312, v8, v9);
								v16.field_0 = v8;
								v16.field_4 = v9;
								sub_543110((const char*)getMemAt(0x973F18, 30760), &v16);
								if (*getMemU32Ptr(0x5D4594, 1599580)) {
									nox_fs_remove((const char*)getMemAt(0x973F18, 36008));
									nox_fs_move((const char*)getMemAt(0x973F18, 38056),
										   (const char*)getMemAt(0x973F18, 36008));
									nox_script_readWriteZzz_541670((const char*)getMemAt(0x973F18, 36008),
																   (const char*)getMemAt(0x973F18, 30760),
																   (const char*)getMemAt(0x973F18, 38056));
								} else {
									*getMemU32Ptr(0x5D4594, 1599580) = 1;
									nox_fs_move((const char*)getMemAt(0x973F18, 30760),
										   (const char*)getMemAt(0x973F18, 38056));
								}
							}
							++dword_5d4594_3835312;
							result = 1;
						}
					}
				}
			}
		}
	}
	return result;
}

//----- (00503EC0) --------------------------------------------------------
int sub_503EC0(int a1, float* a2) {
	float2 a1a; // [esp+0h] [ebp-18h]
	float2 v4;  // [esp+8h] [ebp-10h]
	float2 a2a; // [esp+10h] [ebp-8h]

	if (dword_5d4594_1599480 != dword_5d4594_3835396 || *(int*)&dword_5d4594_1599480 == -1 ||
		dword_5d4594_1599476 == 1) {
		return 0;
	}
	a1a.field_0 = (double)*getMemIntPtr(0x5D4594, 1599508);
	a1a.field_4 = (double)*getMemIntPtr(0x5D4594, 1599512);
	sub_4D3E30(&a1a, &a2a);
	sub_4D3E30((float2*)(a1 + 56), &v4);
	*a2 = v4.field_0 - a2a.field_0;
	a2[1] = v4.field_4 - a2a.field_4;
	return 1;
}

//----- (005040A0) --------------------------------------------------------
uint32_t* nox_xxx_tileAllocTileInCoordList_5040A0(int a1, int a2, float a3) {
	uint32_t* result; // eax
	uint32_t* v4;     // esi
	void* v5;         // eax
	double v6;        // st7
	bool v7;          // zf
	float v8;         // [esp+10h] [ebp+Ch]

	result = calloc(1, 0x18u);
	v4 = result;
	if (result) {
		v5 = calloc(1, 0x14u);
		*v4 = v5;
		if (v5) {
			v4[5] = 0;
			v4[4] = dword_5d4594_1599556;
			if (dword_5d4594_1599556) {
				*(uint32_t*)((uint32_t)dword_5d4594_1599556 + 20) = v4;
			}
			dword_5d4594_1599556 = v4;
			v6 = (double)a1 * 46.0;
			++*getMemU32Ptr(0x5D4594, 1599560);
			v7 = LOBYTE(a3) == 1;
			*((uint8_t*)v4 + 12) = LOBYTE(a3);
			*((float*)v4 + 1) = v6;
			v8 = (double)a2 * 46.0;
			*((float*)v4 + 2) = v8;
			result = v4;
			if (v7) {
				*((float*)v4 + 1) = v6 + 23.0;
			} else {
				*((float*)v4 + 2) = v8 + 23.0;
			}
		} else {
			free(v4);
			result = 0;
		}
	}
	return result;
}

//----- (00504150) --------------------------------------------------------
extern uint32_t nox_tile_def_cnt;
extern nox_tileDef_t nox_tile_defs_arr[176];
int nox_xxx_tileInit_504150(int a1, int a2) {
	int v5;         // edi
	int* i;         // esi
	float2 a1a;     // [esp+Ch] [ebp-50h]
	char v8[72];    // [esp+14h] [ebp-48h]
	float v9;       // [esp+60h] [ebp+4h]
	float v10;      // [esp+64h] [ebp+8h]

	if (*getMemIntPtr(0x587000, 229704) == -1) {
		if (nox_tile_def_cnt > 0) {
			int v2 = 0;
			for (int i = 0; i < nox_tile_def_cnt; i++) {
				nox_tileDef_t* p = &nox_tile_defs_arr[i];
				if (strcmp(&p->name[0], "TransparentFloor") == 0) {
					*getMemU32Ptr(0x587000, 229704) = i;
					v2 = i;
					break;
				}
			}
			if (v2 == -1) {
				return 0;
			}
		}
	}
	memcpy(v8, sub_4D3C70(), sizeof(v8));
	v5 = dword_5d4594_1599556;
	if (dword_5d4594_1599556) {
		v9 = (double)a1;
		v10 = (double)a2;
		do {
			a1a.field_0 = v9 + *(float*)(v5 + 4);
			a1a.field_4 = v10 + *(float*)(v5 + 8);
			nox_xxx_tileCheckImage_51D540(**(uint32_t**)v5);
			nox_xxx_tileCheckImageVari_51D570(*(uint32_t*)(*(uint32_t*)v5 + 4));
			nox_xxx_tile_51D5C0(1);
			if (**(uint32_t**)v5 != *getMemU32Ptr(0x587000, 229704)) {
				sub_51D8F0(&a1a);
			}
			for (i = *(int**)(*(uint32_t*)v5 + 16); i; i = (int*)i[4]) {
				nox_xxx_tileCheckByte3_544070(i[2]);
				nox_xxx_tileCheckByte4_5440A0(i[3]);
				nox_xxx_tileCheckImage_51D540(*i);
				nox_xxx_tileCheckImageVari_51D570(i[1]);
				nox_xxx_tileSubtile_544310(&a1a);
			}
			v5 = *(uint32_t*)(v5 + 16);
		} while (v5);
	}
	nox_xxx_tileInitdataClear_4D3C50(v8);
	return 1;
}

//----- (00504290) --------------------------------------------------------
uint32_t* sub_504290(char a1, char a2) {
	uint32_t* result; // eax
	uint32_t* v3;     // esi
	uint8_t* v4;      // eax

	result = calloc(1, 0xCu);
	v3 = result;
	if (result) {
		result[2] = 0;
		result[1] = dword_5d4594_1599532;
		if (dword_5d4594_1599532) {
			*(uint32_t*)((uint32_t)dword_5d4594_1599532 + 8) = result;
		}
		dword_5d4594_1599532 = result;
		v4 = calloc(1u, 0x24u);
		*v3 = v4;
		v4[5] = a1;
		*(uint8_t*)(*v3 + 6) = a2;
		result = v3;
	}
	return result;
}

//----- (005042F0) --------------------------------------------------------
uint32_t* nox_xxx_cliWallGet_5042F0(int a1, int a2) {
	uint32_t* result; // eax

	result = dword_5d4594_1599532;
	if (!dword_5d4594_1599532) {
		return 0;
	}
	while (*(unsigned char*)(*result + 5) != a1 || *(unsigned char*)(*result + 6) != a2) {
		result = (uint32_t*)result[1];
		if (!result) {
			return 0;
		}
	}
	return result;
}

//----- (00504330) --------------------------------------------------------
int sub_504330(int a1, int a2) {
	unsigned char** v2; // edi
	int v3;             // ebp
	int v4;             // ecx
	int v6;             // ebx
	unsigned char* v7;  // esi
	unsigned char v8;   // dl
	unsigned char v9;   // al
	uint32_t* v11;      // [esp-4h] [ebp-14h]
	int v12;            // [esp-4h] [ebp-14h]

	v2 = dword_5d4594_1599532;
	if (!dword_5d4594_1599532) {
		return 1;
	}
	while (1) {
		v3 = (a1 + 23 * (*v2)[5]) / 23;
		v4 = a2 + 23 * (*v2)[6];
		v6 = v4 / 23;
		v7 = (unsigned char*)nox_server_getWallAtGrid_410580(v3, v6);
		if (v7) {
			v8 = **v2;
			if (dword_5d4594_3835368) {
				*v7 = nox_xxx_wall_42A6C0(*v7, v8);
			} else {
				*v7 = v8;
			}
			goto LABEL_8;
		}
		v7 = (unsigned char*)nox_xxx_wallCreateAt_410250(v3, v6);
		if (!v7) {
			return 0;
		}
		*v7 = **v2;
	LABEL_8:
		v7[1] = (*v2)[1];
		v7[2] = (*v2)[2];
		v7[7] = (*v2)[7];
		if (((*v2)[4] & 0x80u) != 0) {
			v7[4] |= 0x80u;
		}
		if (v7[2] >= nox_xxx_mapWallMaxVariation_410DD0(v7[1], *v7, 0)) {
			v7[2] = 0;
		}
		if (v7[4] & 4) {
			sub_4107A0(*((void**)v7 + 7));
		}
		if ((*v2)[4] & 4) {
			v7[4] |= 4u;
			v11 = (uint32_t*)*((uint32_t*)*v2 + 7);
			*((uint32_t*)v7 + 7) = v11;
			nox_xxx_wallSecretBlock_410760(v11);
		}
		if ((*v2)[4] & 8) {
			v9 = v7[4];
			if (!(v9 & 8)) {
				v12 = *((uint32_t*)v7 + 7);
				v7[4] = v9 | 8;
				nox_xxx_wallBreackableListAdd_410840(v12);
			}
		}
		if ((*v2)[4] & 0x40) {
			v7[4] |= 0x40u;
		}
		v2 = (unsigned char**)v2[1];
		if (!v2) {
			return 1;
		}
	}
}

//----- (005044B0) --------------------------------------------------------
uint32_t* sub_5044B0(int a1, float a2, float a3) {
	uint32_t* result; // eax
	uint32_t* v4;     // esi
	uint32_t* v5;     // eax
	uint32_t* v6;     // eax

	result = calloc(1, 0xCu);
	v4 = result;
	if (result) {
		v5 = sub_579E70();
		*v4 = v5;
		if (v5) {
			v4[2] = 0;
			v4[1] = dword_5d4594_1599548;
			if (dword_5d4594_1599548) {
				*(uint32_t*)((uint32_t)dword_5d4594_1599548 + 8) = v4;
			}
			dword_5d4594_1599548 = v4;
			*(uint32_t*)*v4 = a1;
			*(float*)(*v4 + 8) = a2;
			*(float*)(*v4 + 12) = a3;
			*(uint32_t*)(*v4 + 488) = 0;
			v6 = (uint32_t*)v4[1];
			if (v6) {
				*(uint32_t*)(*v4 + 484) = *v6;
				*(uint32_t*)(*(uint32_t*)v4[1] + 488) = *v4;
				result = v4;
			} else {
				result = v4;
				*(uint32_t*)(*v4 + 484) = 0;
			}
		} else {
			free(v4);
			result = 0;
		}
	}
	return result;
}

//----- (00504560) --------------------------------------------------------
void sub_579E90(nox_waypoint_t* a1);
int sub_504560(int a1, int a2) {
	int* v2;  // esi
	float v4; // [esp+8h] [ebp+4h]
	float v5; // [esp+Ch] [ebp+8h]

	v2 = dword_5d4594_1599548;
	if (dword_5d4594_1599548) {
		v4 = (double)a1;
		v5 = (double)a2;
		do {
			*(float*)(*v2 + 8) = v4 + *(float*)(*v2 + 8);
			*(float*)(*v2 + 12) = v5 + *(float*)(*v2 + 12);
			sub_579E90(*v2);
			v2 = (int*)v2[1];
		} while (v2);
	}
	return 1;
}

//----- (005048A0) --------------------------------------------------------
uint32_t* nox_xxx_unitAddToList_5048A0(int a1) {
	uint32_t* result; // eax
	uint32_t* v2;     // ecx

	result = calloc(1, 0xCu);
	if (!result) {
		return 0;
	}
	result[2] = 0;
	*result = a1;
	result[1] = dword_5d4594_1599540;
	if (dword_5d4594_1599540) {
		*(uint32_t*)((uint32_t)dword_5d4594_1599540 + 8) = result;
	}
	dword_5d4594_1599540 = result;
	*(uint32_t*)(*result + 448) = 0;
	v2 = (uint32_t*)result[1];
	if (v2) {
		*(uint32_t*)(*result + 444) = *v2;
		*(uint32_t*)(*(uint32_t*)result[1] + 448) = *result;
	} else {
		*(uint32_t*)(*result + 444) = 0;
	}
	return result;
}

//----- (00504910) --------------------------------------------------------
int sub_504910(int a1, int a2) {
	int* v2;  // esi
	float v4; // [esp+8h] [ebp+4h]
	float v5; // [esp+Ch] [ebp+8h]

	v2 = dword_5d4594_1599540;
	if (dword_5d4594_1599540) {
		v4 = (double)a1;
		v5 = (double)a2;
		do {
			*(float*)(*v2 + 56) = v4 + *(float*)(*v2 + 56);
			*(float*)(*v2 + 60) = v5 + *(float*)(*v2 + 60);
			nox_xxx_createAt_4DAA50(*v2, 0, *(float*)(*v2 + 56), *(float*)(*v2 + 60));
			*(uint32_t*)(*v2 + 16) |= 0x80000000;
			v2 = (int*)v2[1];
		} while (v2);
	}
	return 1;
}

//----- (00504980) --------------------------------------------------------
int sub_504980() {
	int result; // eax

	if ((dword_5d4594_1599480 == dword_5d4594_3835396 && *(int*)&dword_5d4594_1599480 != -1 &&
			 dword_5d4594_1599476 != 1 ||
		 nox_xxx_mapgenSaveMap_503830(*(int*)&dword_5d4594_3835396)) &&
		dword_5d4594_1599540) {
		result = **(uint32_t**)&dword_5d4594_1599540;
	} else {
		result = 0;
	}
	return result;
}

//----- (005049C0) --------------------------------------------------------
int sub_5049C0(int a1) {
	int result; // eax

	result = a1;
	if (a1) {
		result = *(uint32_t*)(a1 + 444);
	}
	return result;
}

//----- (005049D0) --------------------------------------------------------
void* sub_5049D0() { return dword_5d4594_1599540; }

//----- (005049E0) --------------------------------------------------------
int sub_5049E0(int a1) {
	int result; // eax

	result = a1;
	if (a1) {
		result = *(uint32_t*)(a1 + 4);
	}
	return result;
}

//----- (00504A10) --------------------------------------------------------
int sub_504A10(int a1) {
	int* v1; // esi
	int v3;  // eax
	int v4;  // eax
	int v5;  // ecx
	int v6;  // ecx

	if (!a1) {
		return 0;
	}
	v1 = dword_5d4594_1599540;
	if (!dword_5d4594_1599540) {
		return 0;
	}
	while (*v1 != a1) {
		v1 = (int*)v1[1];
		if (!v1) {
			return 0;
		}
	}
	v3 = v1[1];
	if (v3) {
		*(uint32_t*)(v3 + 8) = v1[2];
	}
	v4 = v1[2];
	if (v4) {
		*(uint32_t*)(v4 + 4) = v1[1];
	}
	if (v1 == dword_5d4594_1599540) {
		dword_5d4594_1599540 = v1[1];
	}
	v5 = *(uint32_t*)(*v1 + 444);
	if (v5) {
		*(uint32_t*)(v5 + 448) = *(uint32_t*)(*v1 + 448);
	}
	v6 = *(uint32_t*)(*v1 + 448);
	if (v6) {
		*(uint32_t*)(v6 + 444) = *(uint32_t*)(*v1 + 444);
	}
	nox_xxx_objectFreeMem_4E38A0(*v1);
	free(v1);
	return 1;
}

//----- (00505060) --------------------------------------------------------
void* sub_505060() {
	void* result; // eax

	result = *(void**)&dword_5d4594_1599616;
	if (dword_5d4594_1599616) {
		free(*(void**)&dword_5d4594_1599616);
		dword_5d4594_1599616 = 0;
	}
	return result;
}

//----- (00505080) --------------------------------------------------------
int nox_server_mapRWMapIntro_505080() {
	FILE* v0;         // ebx
	int v1;           // ebp
	char* v2;         // edi
	short v3;         // dx
	unsigned char v4; // al
	char* v5;         // edi
	char* v6;         // edi
	unsigned char v7; // cl
	short v8;         // si
	int v9;           // edi
	FILE* v10;        // eax
	FILE* v11;        // esi
	int result;       // eax
	uint8_t* v13;     // eax
	uint8_t* v14;     // edi
	int i;            // esi
	char v16;         // [esp+13h] [ebp-409h]
	size_t v17;       // [esp+14h] [ebp-408h]
	int v18;          // [esp+18h] [ebp-404h]
	char v19[1024];   // [esp+1Ch] [ebp-400h]

	v0 = 0;
	v18 = 1;
	v17 = 0;
	v1 = nox_common_gameFlags_check_40A5C0(0x200000);
	sub_505060();
	v2 = nox_fs_root();
	v3 = *getMemU16Ptr(0x587000, 229832);
	strcpy(v19, v2);
	v4 = getMemByte(0x587000, 229834);
	v5 = &v19[strlen(v19)];
	*(uint32_t*)v5 = *getMemU32Ptr(0x587000, 229828);
	*((uint16_t*)v5 + 2) = v3;
	v5[6] = v4;
	strcat(v19, nox_xxx_mapGetMapName_409B40());
	*(uint16_t*)&v19[strlen(v19)] = *getMemU16Ptr(0x587000, 229836);
	strcat(v19, nox_xxx_mapGetMapName_409B40());
	v6 = &v19[strlen(v19) + 1];
	v7 = getMemByte(0x587000, 229844);
	v8 = v18;
	*(uint32_t*)--v6 = *getMemU32Ptr(0x587000, 229840);
	v6[4] = v7;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v18, 2u);
	if (v8 > (short)v18) {
		return 0;
	}
	v9 = 0;
	if (nox_crypt_IsReadOnly()) {
		if (nox_crypt_IsReadOnly() != 1) {
			return 0;
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v17, 4u);
		if ((int)v17 <= 0) {
			return 1;
		}
		if (nox_common_gameFlags_check_40A5C0(0x400000)) {
			nox_xxx_cryptSeekCur_40E0A0(v17);
			return 1;
		}
		if (v1) {
			v0 = nox_fs_create(v19);
			if (!v0) {
				return 0;
			}
			v14 = (uint8_t*)v18;
		} else {
			v13 = calloc(1, v17);
			dword_5d4594_1599616 = v13;
			if (!v13) {
				return 0;
			}
			v14 = v13;
		}
		for (i = 0; i < (int)v17; ++i) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v16, 1u);
			if (v1) {
				nox_fs_fwrite(v0, &v16, 1);
			} else {
				*v14++ = v16;
			}
		}
		if (v0) {
			nox_fs_close(v0);
		}
		return 1;
	}
	if (v1 && (v10 = nox_fs_open(v19), (v11 = v10) != 0)) {
		v17 = nox_fs_fsize(v11);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v17, 4u);
		if ((int)v17 > 0) {
			do {
				nox_fs_fread(v11, &v16, 1);
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v16, 1u);
				++v9;
			} while (v9 < (int)v17);
		}
		nox_fs_close(v11);
		result = 1;
	} else {
		v17 = 0;
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v17, 4u);
		result = 1;
	}
	return result;
}
// 505080: using guessed type char var_400[1024];

//----- (00505C30) --------------------------------------------------------
int sub_57C130(uint32_t* a1, int a2);
int nox_server_mapLoadAddGroup_57C0C0(char* a1, unsigned int a2, unsigned char a3);
int nox_server_mapRWGroupData_505C30() {
	char v0;      // bp
	int i;        // eax
	int j;        // eax
	char v7;      // al
	bool v8;      // zf
	char* v9;     // eax
	char* v10;    // esi
	int v12;      // [esp+10h] [ebp-15Ch]
	char v14[2];  // [esp+1Ah] [ebp-152h]
	int v15;      // [esp+1Ch] [ebp-150h]
	int v16;      // [esp+20h] [ebp-14Ch]
	int v17;      // [esp+24h] [ebp-148h]
	int v18;      // [esp+28h] [ebp-144h]
	int v19[2];   // [esp+2Ch] [ebp-140h]
	int v21;      // [esp+34h] [ebp-138h]
	int v22;      // [esp+38h] [ebp-134h]
	char v23[76]; // [esp+3Ch] [ebp-130h]
	char v24[76]; // [esp+88h] [ebp-E4h]
	char v25[76]; // [esp+D4h] [ebp-98h]
	char v26[76]; // [esp+120h] [ebp-4Ch]

	v15 = 3;
	v0 = nox_xxx_wallGet_426A30();
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v15, 2u);
	if ((short)v15 > 3) {
		return 0;
	}

	if (nox_crypt_IsReadOnly()) {
		v21 = 0;
		v22 = 0;
		int v13;
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v13, 4u);
		if (v13 <= 0) {
			return 1;
		}
		for (int v2 = 0; v2 < v13; ++v2) {
			v17 = 0;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v17, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v23, (unsigned char)v17);
			v23[v17] = 0; // null-term missing
			v8 = (uint16_t)v15 == 2;
			if ((short)v15 < 2) {
				if (strlen(nox_server_currentMapGetFilename_409B30()) + 1 + strlen(v23) >= 0x35) {
					return 0;
				}
				v9 = nox_server_currentMapGetFilename_409B30();
				nox_sprintf(v25, "%s.map:%s", v9, v23);
				strcpy(v23, v25);
				v8 = (uint16_t)v15 == 2;
			}
			if (v8) {
				strcpy(v14, ":");
				strcpy(v24, nox_server_currentMapGetFilename_409B30());
				*(uint16_t*)&v24[strlen(v24)] = *getMemU16Ptr(0x587000, 229976);
				strcpy(v26, v23);
				strtok(v26, v14);
				v10 = strtok(0, v14);
				if (strlen(nox_server_currentMapGetFilename_409B30()) + 1 + strlen(v10) >= 0x35) {
					return 0;
				}
				strcat(v24, v10);
				strcpy(v23, v24);
			}
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v18, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v16, 4u);
			if (!(v0 & 4)) {
				if (nox_common_gameFlags_check_40A5C0(0x400000)) {
					sub_504600(v23, v16, v18);
				} else if (nox_common_gameFlags_check_40A5C0(2097153)) {
					nox_server_mapLoadAddGroup_57C0C0(v23, v16, v18);
				}
			}
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v12, 4u);

			for (int v11 = 0; v11 < v12; ++v11) {
				if ((uint8_t)v18) {
					if ((uint8_t)v18 == 1) {
						nox_xxx_fileReadWrite_426AC0_file3_fread(&v19[0], 4u);
					} else if ((uint8_t)v18 != 2) {
						if ((uint8_t)v18 != 3) {
							return 0;
						}
						nox_xxx_fileReadWrite_426AC0_file3_fread(&v19[0], 4u);
					} else {
						nox_xxx_fileReadWrite_426AC0_file3_fread(&v19[0], 4u);
						nox_xxx_fileReadWrite_426AC0_file3_fread(&v19[1], 4u);
					}
				} else {
					nox_xxx_fileReadWrite_426AC0_file3_fread(&v19[0], 4u);
				}
				if (!(v0 & 4)) {
					if (nox_common_gameFlags_check_40A5C0(0x400000)) {
						sub_5046A0(v19, v16);
					} else {
						sub_57C130(v19, v16);
					}
				}
			}
		}
		return 1;
	}
	int v13 = 0;
	for (i = nox_server_getFirstMapGroup_57C080(); i; i = nox_server_getNextMapGroup_57C090(i)) {
		++v13;
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v13, 4u);
	for (int v4 = nox_server_getFirstMapGroup_57C080(); v4; v4 = nox_server_getNextMapGroup_57C090(v4)) {
		LOBYTE(v17) = strlen((const char*)(v4 + 8)) + 1;
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v17, 1u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v4 + 8), (unsigned char)v17);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)v4, 1u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v4 + 4), 4u);
		v12 = 0;
		for (j = *(uint32_t*)(v4 + 84); j; j = *(uint32_t*)(j + 8)) {
			++v12;
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v12, 4u);
		for (int v6 = *(uint32_t*)(v4 + 84); v6; v6 = *(uint32_t*)(v6 + 8)) {
			v7 = *(uint8_t*)v4;
			if (!*(uint8_t*)v4 || v7 == 1) {
				nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)v6, 4u);
			} else if (v7 != 2) {
				if (v7 != 3) {
					return 0;
				}
				nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)v6, 4u);
			} else {
				nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)v6, 4u);
				nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v6 + 4), 4u);
			}
		}
	}
	return 1;
}
// 505C30: using guessed type char var_E4[76];

//----- (00506260) --------------------------------------------------------
nox_waypoint_t* nox_xxx_waypointNewNotMap_579970(int a1, float a2, float a3);
int nox_server_mapRWWaypoints_506260(uint32_t* a1) {
	float* v2;        // esi
	uint32_t* v3;     // edi
	char* v4;         // esi
	int v5;           // ebx
	uint8_t* v6;      // edi
	char* v7;         // esi
	char** v8;        // eax
	char* v9;         // ebp
	uint8_t* v10;     // esi
	int v11;          // ebx
	uint8_t* v12;     // ebp
	uint8_t* v13;     // edi
	int v14;          // ebx
	uint8_t* v15;     // ebp
	uint8_t* v16;     // edi
	float* v17;       // [esp+10h] [ebp-9Ch]
	int v18;          // [esp+14h] [ebp-98h]
	int v19;          // [esp+18h] [ebp-94h]
	int v20;          // [esp+1Ch] [ebp-90h]
	unsigned int v21; // [esp+20h] [ebp-8Ch]
	float v22;        // [esp+24h] [ebp-88h]
	float v23;        // [esp+28h] [ebp-84h]
	int v24;          // [esp+2Ch] [ebp-80h]
	int v25;          // [esp+30h] [ebp-7Ch]
	int2 v26;         // [esp+34h] [ebp-78h]
	int v27;          // [esp+3Ch] [ebp-70h]
	long long v28;    // [esp+40h] [ebp-6Ch]
	long long v29;    // [esp+48h] [ebp-64h]
	int4 v30;         // [esp+50h] [ebp-5Ch]
	char v31[76];     // [esp+60h] [ebp-4Ch]

	v18 = 4;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v18, 2u);
	if ((short)v18 > 4) {
		return 0;
	}
	if (nox_crypt_IsReadOnly()) {
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v19, 4u);
		v24 = 0;
		if (v19 > 0) {
			while (1) {
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v25, 4u);
				if ((short)v18 < 4) {
					nox_xxx_fileReadWrite_426AC0_file3_fread(&v21, 4u);
					v28 = v21;
					v22 = (double)v21;
					nox_xxx_fileReadWrite_426AC0_file3_fread(&v21, 4u);
					v29 = v21;
					v23 = (double)v21;
				} else {
					nox_xxx_fileReadWrite_426AC0_file3_fread(&v22, 4u);
					nox_xxx_fileReadWrite_426AC0_file3_fread(&v23, 4u);
				}
				if ((short)v18 >= 3) {
					nox_xxx_fileReadWrite_426AC0_file3_fread(&v20, 1u);
					if ((uint8_t)v20) {
						nox_xxx_fileReadWrite_426AC0_file3_fread(v31, (unsigned char)v20);
						v31[(unsigned char)v20] = 0;
					} else {
						v31[0] = getMemByte(0x5D4594, 1599648);
					}
				}
				if (a1) {
					v7 = nox_xxx_mapGetWallSize_426A70();
					sub_428170(a1, &v30);
					v22 = v22 - (double)(int)(23 * *(uint32_t*)v7) + (double)v30.field_0 - 11.0;
					v23 = v23 - (double)(int)(23 * *((uint32_t*)v7 + 1)) + (double)v30.field_4 - 11.0;
				}
				if (nox_common_gameFlags_check_40A5C0(0x400000)) {
					v8 = (char**)sub_5044B0(v25, v22, v23);
					v9 = *v8;
					v17 = (float*)*v8;
				} else {
					v17 = nox_xxx_waypointNewNotMap_579970(v25, v22, v23);
					v9 = (char*)v17;
				}
				if (!v9) {
					break;
				}
				if ((short)v18 >= 3) {
					strcpy(v9 + 16, v31);
				}
				nox_xxx_fileReadWrite_426AC0_file3_fread(v9 + 480, 4u);
				if ((short)v18 < 4) {
					nox_xxx_fileReadWrite_426AC0_file3_fread(&v21, 4u);
					v10 = v9 + 476;
					v9[476] = v21;
				} else {
					v10 = v9 + 476;
					nox_xxx_fileReadWrite_426AC0_file3_fread(v9 + 476, 1u);
				}
				if ((short)v18 >= 2) {
					v14 = 0;
					if (*v10) {
						v15 = v9 + 96;
						v16 = v17 + 87;
						do {
							nox_xxx_fileReadWrite_426AC0_file3_fread(v16, 4u);
							nox_xxx_fileReadWrite_426AC0_file3_fread(v15, 1u);
							++v14;
							v16 += 4;
							v15 += 8;
						} while (v14 < (unsigned char)*v10);
					}
				} else {
					v11 = 0;
					if (*v10) {
						v12 = v9 + 96;
						v13 = v17 + 87;
						do {
							nox_xxx_fileReadWrite_426AC0_file3_fread(v13, 4u);
							*v12 = 2;
							++v11;
							v13 += 4;
							v12 += 8;
						} while (v11 < (unsigned char)*v10);
					}
				}
				if (++v24 >= v19) {
					return 1;
				}
			}
			return 0;
		}
		return 1;
	}
	v19 = 0;
	v2 = (float*)nox_xxx_waypointGetList_579860();
	if (v2) {
		do {
			v3 = a1;
			if (!a1 ||
				(v26.field_0 = (long long)v2[2], v26.field_4 = (long long)v2[3], nox_xxx_wallMath_427F30(&v26, a1))) {
				++v19;
			}
			v2 = (float*)nox_xxx_waypointNext_579870((int)v2);
		} while (v2);
	} else {
		v3 = a1;
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v19, 4u);
	v4 = (char*)nox_xxx_waypointGetList_579860();
	if (!v4) {
		return 1;
	}
	do {
		if (!v3 || (v26.field_0 = (long long)*((float*)v4 + 2), v26.field_4 = (long long)*((float*)v4 + 3),
					nox_xxx_wallMath_427F30(&v26, v3))) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(v4, 4u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v4 + 8, 4u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v4 + 12, 4u);
			LOBYTE(v20) = strlen(v4 + 16);
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v20, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v4 + 16, (unsigned char)v20);
			v27 = *((uint32_t*)v4 + 120) & 1;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v27, 4u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v4 + 476, 1u);
			v5 = 0;
			if (v4[476]) {
				v6 = v4 + 96;
				do {
					nox_xxx_fileReadWrite_426AC0_file3_fread(*((uint8_t**)v6 - 1), 4u);
					nox_xxx_fileReadWrite_426AC0_file3_fread(v6, 1u);
					++v5;
					v6 += 8;
				} while (v5 < (unsigned char)v4[476]);
			}
			v3 = a1;
		}
		v4 = (char*)nox_xxx_waypointNext_579870((int)v4);
	} while (v4);
	return 1;
}
// 506260: using guessed type char var_4C[76];

//----- (005066D0) --------------------------------------------------------
int nox_xxx_allocVoteArray_5066D0() {
	int result; // eax

	result = nox_new_alloc_class("VoteClass", 52, 64);
	nox_alloc_vote_1599652 = result;
	if (result) {
		dword_5d4594_1599656 = 0;
		result = 1;
	}
	return result;
}

//----- (00506720) --------------------------------------------------------
int sub_506720() {
	int result; // eax

	nox_free_alloc_class(*(void**)&nox_alloc_vote_1599652);
	result = 0;
	nox_alloc_vote_1599652 = 0;
	dword_5d4594_1599656 = 0;
	return result;
}

//----- (00506740) --------------------------------------------------------
int sub_506740(nox_object_t* a1p) {
	int a1 = a1p;
	int result; // eax
	int v2;     // esi
	int v3;     // ecx
	int v4;     // edi

	result = a1;
	if (a1) {
		if (*(uint8_t*)(a1 + 8) & 4) {
			result = dword_5d4594_1599656;
			v2 = 1 << *(uint8_t*)(*(uint32_t*)(*(uint32_t*)(a1 + 748) + 276) + 2064);
			if (dword_5d4594_1599656) {
				do {
					v3 = *(uint32_t*)(result + 8);
					v4 = *(uint32_t*)(result + 44);
					if (v3 & v2) {
						*(uint32_t*)(result + 8) = ~v2 & v3;
						--*(uint8_t*)(result + 4);
					}
					if (!*(uint8_t*)(result + 4)) {
						sub_5067B0(result);
					}
					result = v4;
				} while (v4);
			}
		}
	}
	return result;
}

//----- (005067B0) --------------------------------------------------------
void sub_5067B0(int a1) {
	int v1; // esi

	if (a1) {
		if (*(uint32_t*)a1 == 2) {
			v1 = 0;
			do {
				if ((1 << v1) & *(uint32_t*)(a1 + 8)) {
					nox_xxx_netSendVote_506840(v1);
				}
				++v1;
			} while (v1 < 32);
		}
		sub_506810(a1);
		nox_alloc_class_free_obj_first(*(unsigned int**)&nox_alloc_vote_1599652, (uint64_t*)a1);
		if (!dword_5d4594_1599656) {
			sub_507190(255, 0);
		}
	}
}

//----- (00506810) --------------------------------------------------------
int sub_506810(int a1) {
	int result; // eax
	int v2;     // ecx
	int v3;     // ecx

	result = a1;
	v2 = *(uint32_t*)(a1 + 44);
	if (v2) {
		*(uint32_t*)(v2 + 48) = *(uint32_t*)(a1 + 48);
	}
	v3 = *(uint32_t*)(a1 + 48);
	if (v3) {
		result = *(uint32_t*)(a1 + 44);
		*(uint32_t*)(v3 + 44) = result;
	} else {
		dword_5d4594_1599656 = *(uint32_t*)(a1 + 44);
	}
	return result;
}

//----- (00506840) --------------------------------------------------------
int nox_xxx_netSendVote_506840(int a1) {
	char v2[2]; // [esp+0h] [ebp-2h]

	v2[0] = -18;
	v2[1] = 7;
	return nox_xxx_netSendPacket1_4E5390(a1, (int)v2, 2, 0, 1);
}

//----- (00506870) --------------------------------------------------------
char sub_506870(int a1, int a2, wchar2_t* a3) {
	char result; // al

	result = a2;
	if (a2 && *(uint8_t*)(a2 + 8) & 4) {
		switch (a1) {
		case 0:
			result = sub_5068E0(0, a2, a3);
			break;
		case 1:
			result = sub_5068E0(1, a2, a3);
			break;
		case 2:
			result = (unsigned int)sub_506B00(2, a2);
			break;
		case 3:
			result = (unsigned int)sub_506B80(3, a2, a3);
			break;
		default:
			return result;
		}
	}
	return result;
}

//----- (005068E0) --------------------------------------------------------
char sub_5068E0(int a1, int a2, wchar2_t* a3) {
	int v3; // eax
	int v4; // ebp
	int v5; // esi
	int v6; // edi
	int v7; // esi

	LOBYTE(v3) = getMemByte(0x587000, 229980);
	if (*getMemU32Ptr(0x587000, 229980) > 0x20u) {
		return v3;
	}
	if (*getMemU32Ptr(0x587000, 229980) == 0) {
		return v3;
	}
	if (!a3) {
		return v3;
	}
	if (!a2) {
		return v3;
	}
	v4 = 1 << *(uint8_t*)(*(uint32_t*)(*(uint32_t*)(a2 + 748) + 276) + 2064);
	v3 = nox_common_playerInfoGetFirst_416EA0();
	v5 = v3;
	if (!v3) {
		return v3;
	}
	while (1) {
		if (*(uint32_t*)(v5 + 2092) == 1) {
			v3 = nox_wcscmp((const wchar2_t*)(v5 + 4704), a3);
			if (!v3) {
				break;
			}
		}
		v3 = nox_common_playerInfoGetNext_416EE0(v5);
		v5 = v3;
		if (!v3) {
			return v3;
		}
	}
	if (*(uint8_t*)(v5 + 2064) == 31) {
		return v3;
	}
	v6 = *(uint32_t*)(v5 + 2056);
	if (!v6) {
		return v3;
	}
	if (a2 == v6) {
		return v3;
	}
	if (!nox_xxx_CheckGameplayFlags_417DA0(4) || (v3 = nox_xxx_servCompareTeams_419150(a2 + 48, v6 + 48)) != 0) {
		v7 = dword_5d4594_1599656;
		if (dword_5d4594_1599656) {
			while (*(uint32_t*)v7 != a1 || *(uint32_t*)(v7 + 28) != v6) {
				v7 = *(uint32_t*)(v7 + 44);
				if (!v7) {
					break;
				}
			}
		}
		if (!v7) {
			v3 = sub_506A20(a1, a2);
			v7 = v3;
			if (!v3) {
				return v3;
			}
			*(uint32_t*)(v3 + 28) = v6;
			if (nox_xxx_CheckGameplayFlags_417DA0(4)) {
				*(uint32_t*)(v7 + 20) = 1;
			}
		}
		v3 = *(uint32_t*)(v7 + 8);
		if (!(v4 & v3)) {
			LOBYTE(v3) = *(uint8_t*)(v7 + 4) + 1;
			*(uint32_t*)(v7 + 8) |= v4;
			*(uint8_t*)(v7 + 4) = v3;
		}
	}
	return v3;
}

//----- (00506A20) --------------------------------------------------------
uint32_t* sub_506A20(int a1, int a2) {
	int v2;       // ebx
	uint32_t* v3; // esi

	v2 = 0;
	if (!a2 || !(*(uint8_t*)(a2 + 8) & 4)) {
		return 0;
	}
	if (!dword_5d4594_1599656) {
		v2 = 1;
	}
	v3 = nox_alloc_class_new_obj_zero(*(uint32_t**)&nox_alloc_vote_1599652);
	if (!v3) {
		return 0;
	}
	*v3 = a1;
	v3[6] = gameFrame();
	v3[4] = a2 + 48;
	switch (a1) {
	case 0:
	case 1:
		*((uint8_t*)v3 + 12) = getMemByte(0x587000, 229980);
		break;
	case 2:
	case 3:
		*((uint8_t*)v3 + 12) = 6;
		break;
	default:
		*((uint8_t*)v3 + 12) = getMemByte(0x587000, 229984);
		break;
	}
	nox_xxx_voteAddMB_506AD0((int)v3);
	if (v2) {
		sub_507190(255, 1);
	}
	return v3;
}

//----- (00506AD0) --------------------------------------------------------
int nox_xxx_voteAddMB_506AD0(int a1) {
	int result; // eax

	result = a1;
	*(uint32_t*)(a1 + 48) = 0;
	*(uint32_t*)(a1 + 44) = dword_5d4594_1599656;
	if (dword_5d4594_1599656) {
		*(uint32_t*)(dword_5d4594_1599656 + 48) = a1;
	}
	dword_5d4594_1599656 = a1;
	return result;
}

//----- (00506B00) --------------------------------------------------------
uint32_t* sub_506B00(int a1, int a2) {
	uint32_t* result; // eax
	int v3;           // esi
	char v4;          // cl

	result = *(uint32_t**)&nox_server_resetQuestMinVotes_229988;
	if (nox_server_resetQuestMinVotes_229988) {
		if (a2) {
			result = *(uint32_t**)(*(uint32_t*)(a2 + 748) + 276);
			v3 = 1 << *((uint8_t*)result + 2064);
			if (result[1198]) {
				result = *(uint32_t**)&dword_5d4594_1599656;
				if (dword_5d4594_1599656) {
					while (*result != a1) {
						result = (uint32_t*)result[11];
						if (!result) {
							break;
						}
					}
				}
				if (!result) {
					result = sub_506A20(a1, a2);
					if (!result) {
						return result;
					}
					result[5] = 0;
				}
				if (!(result[2] & v3)) {
					v4 = *((uint8_t*)result + 4) + 1;
					result[2] |= v3;
					*((uint8_t*)result + 4) = v4;
				}
			}
		}
	}
	return result;
}

//----- (00506B80) --------------------------------------------------------
uint32_t* sub_506B80(int a1, int a2, wchar2_t* a3) {
	uint32_t* result;  // eax
	int v4;            // edi
	const wchar2_t* v5; // esi
	int v6;            // esi
	char v7;           // cl

	result = *(uint32_t**)&nox_server_kickQuestPlayerMinVotes_229992;
	if (nox_server_kickQuestPlayerMinVotes_229992) {
		if (a3) {
			result = (uint32_t*)a2;
			if (a2) {
				result = *(uint32_t**)(*(uint32_t*)(a2 + 748) + 276);
				v4 = 1 << *((uint8_t*)result + 2064);
				if (result[1198]) {
					result = nox_common_playerInfoGetFirst_416EA0();
					v5 = (const wchar2_t*)result;
					if (result) {
						while (1) {
							if (*((uint32_t*)v5 + 523) == 1) {
								result = (uint32_t*)nox_wcscmp(v5 + 2352, a3);
								if (!result) {
									break;
								}
							}
							result = nox_common_playerInfoGetNext_416EE0((int)v5);
							v5 = (const wchar2_t*)result;
							if (!result) {
								return result;
							}
						}
						if (*((uint8_t*)v5 + 2064) != 31) {
							result = (uint32_t*)*((uint32_t*)v5 + 1198);
							if (result) {
								v6 = *((uint32_t*)v5 + 514);
								if (v6) {
									if (a2 != v6) {
										result = *(uint32_t**)&dword_5d4594_1599656;
										if (dword_5d4594_1599656) {
											while (*result != a1 || result[7] != v6) {
												result = (uint32_t*)result[11];
												if (!result) {
													break;
												}
											}
										}
										if (!result) {
											result = sub_506A20(a1, a2);
											if (!result) {
												return result;
											}
											result[7] = v6;
										}
										if (!(result[2] & v4)) {
											v7 = *((uint8_t*)result + 4) + 1;
											result[2] |= v4;
											*((uint8_t*)result + 4) = v7;
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return result;
}

//----- (00506C90) --------------------------------------------------------
void sub_506C90(int a1, int a2, wchar2_t* a3) {
	if (a2 && *(uint8_t*)(a2 + 8) & 4) {
		switch (a1) {
		case 0:
			sub_506D00(a2, a3);
			break;
		case 1:
			sub_506D00(a2, a3);
			break;
		case 2:
			sub_506DE0(a2);
			break;
		case 3:
			sub_506E50(a2, a3);
			break;
		default:
			return;
		}
	}
}

//----- (00506D00) --------------------------------------------------------
void sub_506D00(int a1, wchar2_t* a2) {
	char* v2; // esi
	int v3;   // esi
	int v4;   // eax
	int v5;   // edx
	int v6;   // esi
	bool v7;  // zf

	if (a1) {
		if (a2) {
			if (*(uint8_t*)(a1 + 8) & 4) {
				v2 = nox_common_playerInfoGetFirst_416EA0();
				if (v2) {
					while (*((uint32_t*)v2 + 523) != 1 || nox_wcscmp((const wchar2_t*)v2 + 2352, a2)) {
						v2 = nox_common_playerInfoGetNext_416EE0((int)v2);
						if (!v2) {
							return;
						}
					}
					if (v2[2064] != 31) {
						v3 = *((uint32_t*)v2 + 514);
						if (v3) {
							v4 = dword_5d4594_1599656;
							v5 = 1 << *(uint8_t*)(*(uint32_t*)(*(uint32_t*)(a1 + 748) + 276) + 2064);
							if (dword_5d4594_1599656) {
								while (*(uint32_t*)v4 || *(uint32_t*)(v4 + 28) != v3 || !(v5 & *(uint32_t*)(v4 + 8))) {
									v4 = *(uint32_t*)(v4 + 44);
									if (!v4) {
										return;
									}
								}
								if (v4) {
									v6 = ~v5 & *(uint32_t*)(v4 + 8);
									v7 = (*(uint8_t*)(v4 + 4))-- == 1;
									*(uint32_t*)(v4 + 8) = v6;
									if (v7) {
										sub_5067B0(v4);
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

//----- (00506DE0) --------------------------------------------------------
void sub_506DE0(int a1) {
	int result;
	int v2;  // edx
	char v3; // cl
	int v4;  // esi

	if (a1) {
		if (*(uint8_t*)(a1 + 8) & 4) {
			result = dword_5d4594_1599656;
			v2 = 1 << *(uint8_t*)(*(uint32_t*)(*(uint32_t*)(a1 + 748) + 276) + 2064);
			if (dword_5d4594_1599656) {
				while (*(uint32_t*)result != 2) {
					result = *(uint32_t*)(result + 44);
					if (!result) {
						return;
					}
				}
				if (result && v2 & *(uint32_t*)(result + 8)) {
					v3 = *(uint8_t*)(result + 4) - 1;
					v4 = ~v2 & *(uint32_t*)(result + 8);
					*(uint8_t*)(result + 4) = v3;
					*(uint32_t*)(result + 8) = v4;
					if (!v3) {
						sub_5067B0(result);
					}
				}
			}
		}
	}
}

//----- (00506E50) --------------------------------------------------------
void sub_506E50(int a1, wchar2_t* a2) {
	char* v2; // esi
	int v3;   // esi
	int v4;   // eax
	int v5;   // edx
	int v6;   // esi
	bool v7;  // zf

	if (a1) {
		if (a2) {
			if (*(uint8_t*)(a1 + 8) & 4) {
				v2 = nox_common_playerInfoGetFirst_416EA0();
				if (v2) {
					while (*((uint32_t*)v2 + 523) != 1 || nox_wcscmp((const wchar2_t*)v2 + 2352, a2)) {
						v2 = nox_common_playerInfoGetNext_416EE0((int)v2);
						if (!v2) {
							return;
						}
					}
					if (v2[2064] != 31) {
						v3 = *((uint32_t*)v2 + 514);
						if (v3) {
							v4 = dword_5d4594_1599656;
							v5 = 1 << *(uint8_t*)(*(uint32_t*)(*(uint32_t*)(a1 + 748) + 276) + 2064);
							if (dword_5d4594_1599656) {
								while (*(uint32_t*)v4 != 3 || *(uint32_t*)(v4 + 28) != v3 ||
									   !(v5 & *(uint32_t*)(v4 + 8))) {
									v4 = *(uint32_t*)(v4 + 44);
									if (!v4) {
										return;
									}
								}
								if (v4) {
									v6 = ~v5 & *(uint32_t*)(v4 + 8);
									v7 = (*(uint8_t*)(v4 + 4))-- == 1;
									*(uint32_t*)(v4 + 8) = v6;
									if (v7) {
										sub_5067B0(v4);
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

//----- (00506F80) --------------------------------------------------------
void sub_506F80(int a1) {
	int v1; // esi
	int v3; // esi

	v1 = *(uint32_t*)(a1 + 28);
	if (*(uint8_t*)(v1 + 16) & 0x20) {
		sub_5067B0(a1);
		return;
	}
	*(uint32_t*)(a1 + 16) = v1 + 48;
	if (sub_507000(a1) == 1) {
		v3 = *(uint32_t*)(v1 + 748);
		nox_xxx_playerCallDisconnect_4DEAB0(*(unsigned char*)(*(uint32_t*)(v3 + 276) + 2064), 4);
		sub_416770(15, (wchar2_t*)(*(uint32_t*)(v3 + 276) + 4704), (const char*)(*(uint32_t*)(v3 + 276) + 2112));
		sub_5067B0(a1);
	}
}

//----- (00507000) --------------------------------------------------------
int sub_507000(int a1) {
	int v1; // edi
	int i;  // esi
	int j;  // eax

	v1 = 0;
	if (*(uint8_t*)(a1 + 4) >= *(uint8_t*)(a1 + 12)) {
		return 1;
	}
	if (*(uint32_t*)(a1 + 20) == 1) {
		for (i = nox_xxx_getFirstPlayerUnit_4DA7C0(); i; i = nox_xxx_getNextPlayerUnit_4DA7F0(i)) {
			if (nox_xxx_servCompareTeams_419150(*(uint32_t*)(a1 + 16), i + 48)) {
				++v1;
			}
		}
	} else {
		for (j = nox_xxx_getFirstPlayerUnit_4DA7C0(); j; j = nox_xxx_getNextPlayerUnit_4DA7F0(j)) {
			++v1;
		}
	}
	return *(unsigned char*)(a1 + 4) >= (unsigned int)(v1 - 1) && *(uint8_t*)(a1 + 4) >= 2u;
}

//----- (00507090) --------------------------------------------------------
void sub_507090(int a1) {
	int i;    // esi
	int v3;   // eax
	char* v4; // eax

	nox_xxx_player_4E3CE0();
	if (*(unsigned char*)(a1 + 4) >= nox_xxx_player_4E3CE0()) {
		for (i = nox_xxx_getFirstPlayerUnit_4DA7C0(); i; i = nox_xxx_getNextPlayerUnit_4DA7F0(i)) {
			v3 = *(uint32_t*)(*(uint32_t*)(i + 748) + 276);
			if (*(uint32_t*)(v3 + 4792) == 1) {
				nox_xxx_playerRespawn_4F7EF0(*(uint32_t*)(v3 + 2056));
			}
		}
		nox_game_setQuestStage_4E3CD0(0);
		v4 = nox_xxx_getQuestMapFile_4D0F60();
		nox_xxx_mapLoad_4D2450(v4);
		sub_5067B0(a1);
	}
}

//----- (00507100) --------------------------------------------------------
void sub_507100(int a1) {
	int v1;              // edi
	int v2;              // ebx
	unsigned int v3;     // eax
	unsigned int result; // eax

	v1 = *(uint32_t*)(a1 + 28);
	if (!v1) {
		sub_5067B0(a1);
		return;
	}
	if (*(uint8_t*)(v1 + 16) & 0x20) {
		sub_5067B0(a1);
		return;
	}
	v2 = *(uint32_t*)(v1 + 748);
	if (!*(uint32_t*)(*(uint32_t*)(v2 + 276) + 4792)) {
		sub_5067B0(a1);
		return;
	}
	if (*(uint8_t*)(a1 + 4) >= *(uint8_t*)(a1 + 12)) {
		goto LABEL_8;
	}
	v3 = nox_xxx_player_4E3CE0();
	if (v3 <= 1) {
		sub_5067B0(a1);
		return;
	}
	result = v3 - 1;
	if (!(*(unsigned char*)(a1 + 4) >= result && *(uint8_t*)(a1 + 4) >= 2u)) {
		return;
	}
LABEL_8:
	sub_4DCFB0(v1);
	sub_416770(15, (wchar2_t*)(*(uint32_t*)(v2 + 276) + 4704), (const char*)(*(uint32_t*)(v2 + 276) + 2112));
	sub_5067B0(a1);
	return;
}

//----- (00507190) --------------------------------------------------------
int sub_507190(int a1, char a2) {
	char v4[3]; // [esp+0h] [ebp-4h]
	v4[0] = -18;
	v4[1] = 6;
	v4[2] = a2;
	return nox_xxx_netSendPacket1_4E5390(a1, (int)v4, 3, 0, 1);
}

//----- (005071C0) --------------------------------------------------------
int sub_5071C0() { return dword_5d4594_1599656 != 0; }

//----- (00509120) --------------------------------------------------------
void sub_509120(uint32_t* a1, int a2, const char* a3) {
	char* v3;     // ebx
	int v4;       // ecx
	uint32_t* v5; // esi
	char* v6;     // edx
	uint32_t* v7; // esi
	int v8;       // esi
	uint32_t* v9; // esi

	v3 = (char*)a1[189];
	if (v3) {
		if (a2 == 14) {
			if (nox_common_gameFlags_check_40A5C0(6291456)) {
				strcpy(v3, a3);
			} else {
				a1[192] = nox_script_indexByEvent(a3);
			}
			return;
		}
		v4 = a1[2];
		if (v4 & 0x200) {
			v5 = (uint32_t*)a1[187];
			if (a2) {
				if (a2 == 1) {
					if (!nox_common_gameFlags_check_40A5C0(6291456)) {
						v5[6] = nox_script_indexByEvent(a3);
						return;
					}
					v6 = v3 + 256;
				} else {
					if (a2 != 2) {
						return;
					}
					if (!nox_common_gameFlags_check_40A5C0(6291456)) {
						v5[8] = nox_script_indexByEvent(a3);
						return;
					}
					v6 = v3 + 384;
				}
			} else {
				if (!nox_common_gameFlags_check_40A5C0(6291456)) {
					v5[4] = nox_script_indexByEvent(a3);
					return;
				}
				v6 = v3 + 512;
			}
			strcpy(v6, a3);
			return;
		}
		if (v4 & 2) {
			v7 = (uint32_t*)a1[187];
			switch (a2) {
			case 3:
				if (nox_common_gameFlags_check_40A5C0(6291456)) {
					v6 = v3 + 640;
					strcpy(v6, a3);
					return;
				}
				v7[309] = nox_script_indexByEvent(a3);
				break;
			case 4:
				if (nox_common_gameFlags_check_40A5C0(6291456)) {
					v6 = v3 + 768;
					strcpy(v6, a3);
					return;
				}
				v7[307] = nox_script_indexByEvent(a3);
				break;
			case 5:
				if (nox_common_gameFlags_check_40A5C0(6291456)) {
					v6 = v3 + 896;
					strcpy(v6, a3);
					return;
				}
				v7[317] = nox_script_indexByEvent(a3);
				break;
			case 6:
				if (nox_common_gameFlags_check_40A5C0(6291456)) {
					v6 = v3 + 1024;
					strcpy(v6, a3);
					return;
				}
				v7[311] = nox_script_indexByEvent(a3);
				break;
			case 7:
				if (nox_common_gameFlags_check_40A5C0(6291456)) {
					v6 = v3 + 1152;
					strcpy(v6, a3);
					return;
				}
				v7[313] = nox_script_indexByEvent(a3);
				break;
			case 8:
				if (nox_common_gameFlags_check_40A5C0(6291456)) {
					v6 = v3 + 1280;
					strcpy(v6, a3);
					return;
				}
				v7[315] = nox_script_indexByEvent(a3);
				break;
			case 9:
				if (nox_common_gameFlags_check_40A5C0(6291456)) {
					v6 = v3 + 1408;
					strcpy(v6, a3);
					return;
				}
				v7[319] = nox_script_indexByEvent(a3);
				break;
			case 10:
				if (nox_common_gameFlags_check_40A5C0(6291456)) {
					v6 = v3 + 1536;
					strcpy(v6, a3);
					return;
				}
				v7[321] = nox_script_indexByEvent(a3);
				break;
			case 11:
				if (nox_common_gameFlags_check_40A5C0(6291456)) {
					v6 = v3 + 1664;
					strcpy(v6, a3);
					return;
				}
				v7[323] = nox_script_indexByEvent(a3);
				break;
			case 13:
				if (nox_common_gameFlags_check_40A5C0(6291456)) {
					v6 = v3 + 1792;
					strcpy(v6, a3);
					return;
				}
				v7[325] = nox_script_indexByEvent(a3);
				break;
			default:
				return;
			}
		} else {
			if (v4 & 0x800) {
				v8 = a1[175];
				if (a2 != 12) {
					return;
				}
				if (!nox_common_gameFlags_check_40A5C0(6291456)) {
					*(uint32_t*)(v8 + 4) = nox_script_indexByEvent(a3);
					return;
				}
				v6 = v3 + 128;
				strcpy(v6, a3);
				return;
			}
			if (v4 & 0x20000) {
				v9 = (uint32_t*)a1[187];
				switch (a2) {
				case 15:
					if (nox_common_gameFlags_check_40A5C0(6291456)) {
						v6 = v3 + 1920;
						strcpy(v6, a3);
						return;
					}
					v9[13] = nox_script_indexByEvent(a3);
					break;
				case 16:
					if (nox_common_gameFlags_check_40A5C0(6291456)) {
						v6 = v3 + 2048;
						strcpy(v6, a3);
						return;
					}
					v9[15] = nox_script_indexByEvent(a3);
					break;
				case 17:
					if (nox_common_gameFlags_check_40A5C0(6291456)) {
						v6 = v3 + 2304;
						strcpy(v6, a3);
						return;
					}
					v9[17] = nox_script_indexByEvent(a3);
					break;
				case 18:
					if (nox_common_gameFlags_check_40A5C0(6291456)) {
						v6 = v3 + 2176;
						strcpy(v6, a3);
						return;
					}
					v9[19] = nox_script_indexByEvent(a3);
					break;
				default:
					return;
				}
			}
		}
	}
}
