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
#include "GAME5.h"
#include "GAME5_2.h"
#include "client__gui__guigen.h"
#include "common__binfile.h"
#include "common__crypt.h"
#include "common__net_list.h"
#include "common__random.h"
#include "operators.h"
#include "input.h"

#include "client__video__draw_common.h"
#include "common__magic__speltree.h"
#include "server__script__script.h"

#include <time.h>

extern uint32_t dword_5d4594_2487656;
extern uint32_t dword_5d4594_3835368;
extern uint32_t dword_5d4594_2487804;
extern uint32_t dword_5d4594_2487632;
extern uint32_t dword_5d4594_3835364;
extern uint32_t dword_5d4594_2487884;
extern uint32_t dword_5d4594_3835372;
extern uint32_t dword_5d4594_3835392;
extern uint32_t dword_5d4594_2487580;
extern uint32_t dword_5d4594_2487568;
extern uint32_t dword_5d4594_2487536;
extern uint32_t dword_5d4594_2487628;
extern uint32_t dword_5d4594_2487584;
extern uint32_t dword_5d4594_2487564;
extern uint32_t dword_5d4594_2487672;
extern uint32_t dword_5d4594_2487676;
extern uint32_t dword_5d4594_3835388;
extern uint32_t dword_5d4594_2487652;
extern uint32_t dword_5d4594_3835348;
extern uint32_t dword_5d4594_2487576;
extern uint32_t dword_5d4594_2487624;
extern uint32_t dword_5d4594_3835352;
extern uint32_t dword_5d4594_2487620;
extern uint32_t dword_5d4594_2487708;
extern uint32_t nox_xxx_energyBoltTarget_5d4594_2487880;
extern uint32_t dword_5d4594_2487532;
extern uint32_t dword_5d4594_2487248;

extern uint32_t dword_5d4594_2487560;
extern uint32_t dword_5d4594_2487540;
extern uint32_t dword_5d4594_2487712;
extern uint32_t dword_5d4594_2487524;
extern uint32_t dword_5d4594_2487556;
extern obj_5D4594_2650668_t** ptr_5D4594_2650668;

float get_nox_xxx_warriorMaxHealth_587000_312784();
float get_nox_xxx_wizardMaxHealth_587000_312816();
float get_nox_xxx_conjurerMaxHealth_587000_312800();

float get_nox_xxx_warriorMaxMana_587000_312788();
float get_nox_xxx_wizardMaximumMana_587000_312820();
float get_nox_xxx_conjurerMaxMana_587000_312804();

extern nox_tileDef_t nox_tile_defs_arr[176];

//----- (0051DEA0) --------------------------------------------------------
int nox_xxx_mapCountWallsMB_51DEA0(int a1) {
	int result; // eax

	if ((int)*(unsigned char*)(a1 + 5) < *getMemIntPtr(0x5D4594, 2487252)) {
		*getMemU32Ptr(0x5D4594, 2487252) = *(unsigned char*)(a1 + 5);
	}
	result = *(unsigned char*)(a1 + 6);
	if (result < *getMemIntPtr(0x5D4594, 2487256)) {
		*getMemU32Ptr(0x5D4594, 2487256) = *(unsigned char*)(a1 + 6);
	}
	return result;
}

//----- (0051DED0) --------------------------------------------------------
int sub_51DED0() {
	int* v0;            // edi
	char* v1;           // eax
	float* v3;          // esi
	int v4;             // ebx
	int v5;             // eax

	v0 = (int*)sub_45A060();
	if (!v0) {
		return 1;
	}
	do {
		if (!sub_4E3AD0(v0[27]) && sub_4E3B80(v0[27])) {
			v1 = (char*)nox_get_thing_name(v0[27]);
			v3 = (float*)nox_xxx_newObjectByTypeID_4E3810(v1);
			v4 = *((uint32_t*)v3 + 9);
			v3[14] = (double)v0[3] + 0.5;
			v3[15] = (double)v0[4] + 0.5;
			v5 = v0[32];
			*((uint32_t*)v3 + 10) = v5;
			*((uint32_t*)v3 + 11) = v5;
			*((uint32_t*)v3 + 9) = v5;
			*((uint32_t*)v3 + 4) = v0[30];
			*((uint32_t*)v3 + 5) = v0[70];
			nox_xxx_xfer_saveObj_51DF90((int)v3);
			*((uint32_t*)v3 + 9) = v4;
			nox_xxx_objectFreeMem_4E38A0((int)v3);
		}
		v0 = (int*)nox_drawable_next_45A070((int)v0);
	} while (v0);
	return 1;
}

//----- (0051E010) --------------------------------------------------------
int nox_xxx_nxzCompressFile_57BDD0(char* a1, char* a2);
int nox_xxx_mapSaveMap_51E010(char* a1, int a2) {
	char* v2;         // edi
	unsigned char v3; // cl
	int result;       // eax
	int v5;           // esi
	int v7;           // [esp+10h] [ebp-804h]
	char v8[1024];    // [esp+14h] [ebp-800h]
	char Mem[1024];   // [esp+414h] [ebp-400h]

	v7 = -86050098;
	strcpy(Mem, a1);
	v8[0] = 0;
	strncat(v8, a1, 1024-1);
	v8[strlen(v8)-4] = 0;
	v2 = &v8[strlen(v8) + 1];
	v3 = getMemByte(0x587000, 253116);
	*(uint32_t*)--v2 = *getMemU32Ptr(0x587000, 253112);
	v2[4] = v3;
	result = nox_xxx_cryptOpen_426910(Mem, 0, 19);
	if (result) {
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v7, 4u);
		v5 = nox_xxx_cryptFlush_4268E0();
		if (nox_xxx_map_51E140()) {
			sub_4268F0(v5);
			nox_xxx_cryptClose_4269F0();
			if (!a2 || (result = nox_xxx_nxzCompressFile_57BDD0(Mem, (int)v8)) != 0) {
				result = 1;
			}
		} else {
			nox_xxx_cryptClose_4269F0();
			result = 0;
		}
	}
	return result;
}
// 51E0D5: variable 'v6' is possibly undefined

//----- (0051E140) --------------------------------------------------------
void nox_xxx_map_5004F0();
void nox_xxx_mapSetWallInGlobalDir0pr1_5004D0();
int nox_xxx_map_51E140() {
	int result; // eax
	char v2;    // [esp+1h] [ebp-1h]

	*getMemU32Ptr(0x5D4594, 2487252) = 256;
	*getMemU32Ptr(0x5D4594, 2487256) = 256;
	nox_xxx_wallForeachFn_410640(nox_xxx_mapCountWallsMB_51DEA0, 0);
	nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x5D4594, 2487252), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread(getMemAt(0x5D4594, 2487256), 4u);
	nox_xxx_mapWall_426A80(getMemIntPtr(0x5D4594, 2487252));
	nox_xxx_mapSetWallInGlobalDir0pr1_5004D0();
	if (nox_xxx_mapWriteSectionsMB_426E20(0)) {
		nox_xxx_map_5004F0();
		v2 = 0;
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v2, 1u);
		result = 1;
	} else {
		nox_xxx_cryptClose_4269F0();
		result = 0;
	}
	return result;
}

//----- (00522FF0) --------------------------------------------------------
int nox_xxx_netSendPointFx_522FF0(char a1, float2* a2) {
	short v2;    // ax
	float v3;    // edx
	char a2a[5]; // [esp+4h] [ebp-8h]

	a2a[0] = a1;
	v2 = nox_float2int(a2->field_0);
	v3 = a2->field_4;
	*(uint16_t*)&a2a[1] = v2;
	*(uint16_t*)&a2a[3] = nox_float2int(v3);
	return nox_xxx_netSendFxAllCli_523030(a2, a2a, 5);
}

//----- (00523030) --------------------------------------------------------
int nox_xxx_netSendFxAllCli_523030(float2* a1, const void* a2, int a3) {
	int result; // eax
	int i;      // esi
	int v5;     // ecx
	int v6;     // eax
	float v7;   // edx
	float v8;   // eax
	double v9;  // st7
	double v10; // st7
	double v11; // st6
	float v12;  // [esp+4h] [ebp-18h]
	float v13;  // [esp+Ch] [ebp-10h]
	float v14;  // [esp+10h] [ebp-Ch]
	float v15;  // [esp+18h] [ebp-4h]

	result = nox_xxx_getFirstPlayerUnit_4DA7C0();
	for (i = result; result; i = result) {
		v5 = *(uint32_t*)(*(uint32_t*)(i + 748) + 276);
		if (*(uint8_t*)(v5 + 3680) & 3 && (v6 = *(uint32_t*)(v5 + 3628)) != 0) {
			v7 = *(float*)(v6 + 56);
			v8 = *(float*)(v6 + 60);
			v12 = v7;
		} else {
			v8 = *(float*)(i + 60);
			v12 = *(float*)(i + 56);
		}
		v9 = (double)*(unsigned short*)(v5 + 10);
		v13 = v12 - v9 - 50.0;
		v10 = v9 + v12 + 50.0;
		v11 = (double)*(unsigned short*)(v5 + 12);
		if (a1->field_0 > (double)v13 && v10 > a1->field_0) {
			v14 = v8 - v11 - 50.0;
			if (a1->field_4 > (double)v14) {
				v15 = v11 + v8 + 50.0;
				if (a1->field_4 < (double)v15) {
					nox_netlist_addToMsgListCli_40EBC0(*(unsigned char*)(v5 + 2064), 1, a2, a3);
				}
			}
		}
		result = nox_xxx_getNextPlayerUnit_4DA7F0(i);
	}
	return result;
}

//----- (00523150) --------------------------------------------------------
int sub_523150(char a1, char a2, float* a3) {
	char v4[6]; // [esp+4h] [ebp-8h]

	v4[0] = a1;
	v4[1] = a2;
	*(uint16_t*)&v4[2] = nox_float2int(*a3);
	*(uint16_t*)&v4[4] = nox_float2int(a3[1]);
	return nox_xxx_netSendFxAllCli_523030((float2*)a3, v4, 6);
}

//----- (005231B0) --------------------------------------------------------
int nox_xxx_netSparkExplosionFx_5231B0(float* a1, char a2) {
	short v2;   // ax
	float v3;   // ecx
	char v5[6]; // [esp+4h] [ebp-8h]

	v5[0] = -109;
	v2 = nox_float2int(*a1);
	v3 = a1[1];
	*(uint16_t*)&v5[1] = v2;
	*(uint16_t*)&v5[3] = nox_float2int(v3);
	v5[5] = a2;
	return nox_xxx_netSendFxAllCli_523030((float2*)a1, v5, 6);
}

//----- (00523200) --------------------------------------------------------
void nox_xxx_sendGeneratorBreakFX_523200(float* a1, char a2) {
	short v2;   // ax
	float v3;   // ecx
	char v5[7]; // [esp+4h] [ebp-8h]

	v5[0] = -16;
	v5[1] = 25;
	v2 = nox_float2int(*a1);
	v3 = a1[1];
	*(uint16_t*)&v5[2] = v2;
	*(uint16_t*)&v5[4] = nox_float2int(v3);
	v5[6] = a2;
	nox_xxx_netSendFxAllCli_523030((float2*)a1, v5, 7);
}

//----- (00523270) --------------------------------------------------------
int nox_xxx_netSendVampFx_523270(char a1, short* a2, short a3) {
	short v3;          // dx
	unsigned short v4; // cx
	unsigned short v5; // ax
	float2 a1a;        // [esp+0h] [ebp-14h]
	char a2a[11];      // [esp+8h] [ebp-Ch]

	a2a[0] = a1;
	v3 = a2[2];
	*(uint16_t*)&a2a[1] = *a2;
	v4 = a2[4];
	v5 = a2[6];
	*(uint16_t*)&a2a[3] = v3;
	*(uint16_t*)&a2a[5] = v4;
	*(uint16_t*)&a2a[7] = v5;
	a1a.field_0 = (double)v4;
	*(uint16_t*)&a2a[9] = a3;
	a1a.field_4 = (double)v5;
	return nox_xxx_netSendFxAllCli_523030(&a1a, a2a, 11);
}

//----- (00523530) --------------------------------------------------------
int nox_xxx_netClientPredictLinear_523530(int a1) {
	short v1;      // ax
	double v2;     // st7
	long long v3;  // rax
	double v4;     // st7
	long long v5;  // rax
	double v6;     // st7
	short v7;      // cx
	long long v8;  // rax
	double v9;     // st7
	long long v10; // rax
	double v11;    // st7
	int result;    // eax
	int i;         // esi
	char v14[14];  // [esp+4h] [ebp-10h]

	v14[0] = -75;
	v1 = nox_xxx_netGetUnitCodeServ_578AC0((uint32_t*)a1);
	v2 = *(float*)(a1 + 56);
	*(uint16_t*)&v14[1] = v1;
	*(uint16_t*)&v14[3] = *(uint16_t*)(a1 + 4);
	v3 = (long long)v2;
	v4 = *(float*)(a1 + 60);
	*(uint16_t*)&v14[5] = v3;
	v5 = (long long)v4;
	v6 = *(float*)(a1 + 112) * 16.0;
	v7 = *(uint16_t*)(a1 + 124);
	*(uint16_t*)&v14[7] = v5;
	*(uint16_t*)&v14[9] = v7;
	v8 = (long long)v6;
	v9 = *(float*)(a1 + 80) * 16.0;
	v14[11] = v8;
	v10 = (long long)v9;
	v11 = *(float*)(a1 + 84) * 16.0;
	v14[12] = v10;
	v14[13] = (long long)v11;
	result = nox_xxx_getFirstPlayerUnit_4DA7C0();
	for (i = result; result; i = result) {
		nox_netlist_addToMsgListCli_40EBC0(*(unsigned char*)(*(uint32_t*)(*(uint32_t*)(i + 748) + 276) + 2064), 1, v14,
										   14);
		result = nox_xxx_getNextPlayerUnit_4DA7F0(i);
	}
	return result;
}

//----- (00523670) --------------------------------------------------------
int nox_xxx_netSendShieldFx_523670(int a1, float* a2) {
	char v2;    // al
	int v3;     // eax
	char v5[4]; // [esp+4h] [ebp-Ch]
	float2 v6;  // [esp+8h] [ebp-8h]

	v5[0] = -128;
	*(uint16_t*)&v5[1] = nox_xxx_netGetUnitCodeServ_578AC0((uint32_t*)a1);
	if (a2) {
		v6.field_0 = *(float*)(a1 + 56) - *a2;
		v6.field_4 = *(float*)(a1 + 60) - a2[1];
		v3 = nox_xxx_math_509ED0(&v6);
		v2 = nox_xxx_math_509EA0(v3);
	} else {
		v2 = nox_xxx_math_509EA0(*(short*)(a1 + 124));
	}
	v5[3] = v2;
	return nox_xxx_netSendFxAllCli_523030((float2*)(a1 + 56), v5, 4);
}

//----- (005236F0) --------------------------------------------------------
int nox_xxx_sendSummonStartFX_5236F0(short a1, float* a2, char a3, short a4, short a5) {
	double v5;    // st7
	long long v6; // rax
	double v7;    // st7
	char v9[12];  // [esp+4h] [ebp-Ch]

	v9[0] = 126;
	*(uint16_t*)&v9[5] = a1;
	*(uint16_t*)&v9[7] = a4;
	v5 = *a2;
	v9[9] = a3;
	v6 = (long long)v5;
	v7 = a2[1];
	*(uint16_t*)&v9[1] = v6;
	*(uint16_t*)&v9[3] = (long long)v7;
	*(uint16_t*)&v9[10] = a5;
	return nox_xxx_netSendPacket0_4E5420(255, v9, 12, 0, 1);
}

//----- (00523760) --------------------------------------------------------
int nox_xxx_sendSummonCancelFX_523760(short a1) {
	char v3[3]; // [esp+0h] [ebp-4h]
	v3[0] = 127;
	*(uint16_t*)&v3[1] = a1;
	return nox_xxx_netSendPacket0_4E5420(255, v3, 3, 0, 1);
}

//----- (00523830) --------------------------------------------------------
void nox_xxx_sendGeneratorSpawnFX_523830(int4* a1, short a2) {
	double v2;    // st7
	short v3;     // cx
	short v4;     // dx
	double v5;    // st7
	short v6;     // cx
	float2 a1a;   // [esp+0h] [ebp-14h]
	char a2a[12]; // [esp+8h] [ebp-Ch]

	a2a[0] = -16;
	a2a[1] = 16;
	v2 = (double)a1->field_8;
	v3 = a1->field_0;
	*(uint16_t*)&a2a[4] = a1->field_4;
	v4 = a1->field_C;
	a1a.field_0 = v2;
	v5 = (double)a1->field_C;
	*(uint16_t*)&a2a[2] = v3;
	v6 = a1->field_8;
	*(uint16_t*)&a2a[8] = v4;
	*(uint16_t*)&a2a[6] = v6;
	a1a.field_4 = v5;
	*(uint16_t*)&a2a[10] = a2;
	nox_xxx_netSendFxAllCli_523030(&a1a, a2a, 12);
}

//----- (005238A0) --------------------------------------------------------
void nox_xxx_sendArrowTrapFX_5238A0(float* a1, char a2) {
	short v2;   // ax
	float v3;   // ecx
	char v5[6]; // [esp+4h] [ebp-8h]

	v5[0] = -95;
	v2 = nox_float2int16(*a1);
	v3 = a1[1];
	*(uint16_t*)&v5[1] = v2;
	*(uint16_t*)&v5[3] = nox_float2int16(v3);
	v5[5] = a2;
	nox_xxx_netSendFxAllCli_523030((float2*)a1, v5, 6);
}

int sub_526CA0(char* a1);

//----- (00527E50) --------------------------------------------------------
int nox_xxx_netUpdateObjectSpecial_527E50(nox_object_t* a1p, nox_object_t* a2p) {
	int a1 = a1p;
	uint32_t* a2 = a2p;
	unsigned int v2; // edi
	int v3;          // eax

	v2 = *(unsigned char*)(*(uint32_t*)(*(uint32_t*)(a1 + 748) + 276) + 2064);
	if (a2 && v2 < 0x20) {
		v3 = a2[v2 + 140];
		if (!(v3 & 0xFFF0000)) {
			return 0;
		}
		if ((v3 & 0x10000) == 0x10000) {
			nox_xxx_netReportAnimFrame_4D81F0(v2, a2);
			a2[v2 + 140] &= 0xFFFEFFFF;
		}
		if ((a2[v2 + 140] & 0x20000) == 0x20000) {
			if (!nox_xxx_unitIsEnemyTo_5330C0(a1, (int)a2)) {
				nox_xxx_netReportUnitCurrentHP_4D8620(v2, a2);
			}
			a2[v2 + 140] &= 0xFFFDFFFF;
		}
		if ((a2[v2 + 140] & 0x40000) == 0x40000) {
			nox_xxx_netReportObjHidden_4D8FD0(v2, a2);
			a2[v2 + 140] &= 0xFFFBFFFF;
		}
		if ((a2[v2 + 140] & 0x80000) == 0x80000) {
			nox_xxx_netReportXStatus_4D8230(v2, a2);
			a2[v2 + 140] &= 0xFFF7FFFF;
		}
		if ((a2[v2 + 140] & 0x400000) == 0x400000) {
			nox_xxx_netReportUnitHeight_4D9020(v2, (int)a2);
			a2[v2 + 140] &= 0xFFBFFFFF;
		}
		if ((0x800000 & a2[v2 + 140]) == 0x800000) {
			nox_xxx_netReportEnchant_4D8F90(v2, a2);
			a2[v2 + 140] &= 0xFF7FFFFF;
		}
		if ((a2[v2 + 140] & 0x2000000) == 0x2000000) {
			nox_xxx_netReportTeamBase_4D92D0(v2, (int)a2);
			a2[v2 + 140] &= 0xFDFFFFFF;
		}
		if ((a2[v2 + 140] & 0x4000000) == 0x4000000) {
			nox_xxx_netSendReportNPC_4D93A0(v2, (int)a2);
			a2[v2 + 140] &= 0xFBFFFFFF;
		}
	}
	return 1;
}

//----- (00528030) --------------------------------------------------------
short sub_528030(int a1) {
	int v1;             // ebp
	int v2;             // esi
	unsigned short* v3; // edi
	int v4;             // ebx
	int v5;             // edi
	int v6;             // edx
	unsigned int v7;    // ebx
	int v8;             // eax
	short result;       // ax
	int v10;            // edx
	int v11;            // [esp+18h] [ebp+4h]

	v1 = a1;
	v2 = *(uint32_t*)(a1 + 748);
	v3 = *(unsigned short**)(a1 + 556);
	v4 = *(uint32_t*)(v2 + 276);
	v11 = *(unsigned char*)(v4 + 2064);
	if (*(uint16_t*)(v2 + 10) == *v3) {
		v5 = gameFrame();
		v7 = gameFPS();
	} else if (abs(*v3 - *(unsigned short*)(v2 + 10)) >= v3[2] / 10 ||
		(v5 = gameFrame(), v6 = *(uint32_t*)(v4 + 2176), v7 = gameFPS(),
		 (unsigned int)(gameFrame() - v6) > (int)gameFPS() >> 2)) {
		nox_xxx_netSendPlrHealthToTeam_4D86E0(v11);
		v8 = *(uint32_t*)(v2 + 276);
		*(uint16_t*)(v2 + 10) = **(uint16_t**)(v1 + 556);
		*(uint32_t*)(v8 + 2176) = gameFrame();
		v5 = gameFrame();
		v7 = gameFPS();
	}
	result = *(uint16_t*)(v2 + 4);
	if (*(uint16_t*)(v2 + 6) != result) {
		result = 0;
		if (abs(*(unsigned short*)(v2 + 4) - *(unsigned short*)(v2 + 6)) >= *(unsigned short*)(v2 + 8) / 10 ||
			v5 - *(uint32_t*)(*(uint32_t*)(v2 + 276) + 2180) > v7 >> 2) {
			nox_xxx_netReportMana_4D8930(v11, v1);
			v10 = *(uint32_t*)(v2 + 276);
			*(uint16_t*)(v2 + 6) = *(uint16_t*)(v2 + 4);
			result = (unsigned short)gameFrame();
			*(uint32_t*)(v10 + 2180) = gameFrame();
		}
	}
	return result;
}

//----- (00528190) --------------------------------------------------------
int nox_xxx_checkIsKillable_528190(nox_object_t* a1p) {
	int a1 = a1p;
	uint16_t* v1; // eax
	bool v2;      // zf

	v1 = *(uint16_t**)(a1 + 556);
	if (!v1) {
		return 0;
	}
	if (*v1) {
		v2 = v1[2] == 0;
		if (v1[2]) {
			return 1;
		}
	} else {
		v2 = v1[2] == 0;
	}
	if (v2) {
		return 1;
	} else {
		return 0;
	}
}

//----- (005281D0) --------------------------------------------------------
int nox_xxx_frameCounterSetCopyToNextFrame_5281D0() {
	int result; // eax

	result = gameFrame() + 1;
	*getMemU32Ptr(0x5D4594, 2487684) = gameFrame() + 1;
	return result;
}

//----- (005281E0) --------------------------------------------------------
int nox_xxx_frameCounterSetCopy_5281E0() {
	int result; // eax

	result = gameFrame();
	*getMemU32Ptr(0x5D4594, 2487684) = gameFrame();
	return result;
}

//----- (005281F0) --------------------------------------------------------
int nox_xxx_unitCanSee_536FB0(nox_object_t* a1, nox_object_t* a2, int a3);
void nox_xxx_unitUpdateSightMB_5281F0(nox_object_t* a1p) {
	uint32_t a1 = a1p;
	uint32_t v1;   // edi
	int v2;     // eax
	int v3;     // ebp
	double v4;  // st7
	int v5;     // esi
	int* v6;    // ebx
	double v7;  // st7
	double v8;  // st6
	double v9;  // st7
	double v10; // st6
	int v11;    // eax
	int v12;    // eax
	int v13;    // esi
	int v14;    // eax
	int v15;    // eax
	int v16;    // esi
	int v17;    // [esp+10h] [ebp-10h]
	float v18;  // [esp+10h] [ebp-10h]
	int v19;    // [esp+14h] [ebp-Ch]
	float v20;  // [esp+18h] [ebp-8h]
	float v21;  // [esp+24h] [ebp+4h]

	v1 = a1;
	v17 = 0;
	v2 = *(uint32_t*)(a1 + 16);
	v3 = *(uint32_t*)(a1 + 748);
	if ((v2 & 0x8000) != 0 && !nox_xxx_unitIsZombie_534A40(a1)) {
		return;
	}
	if (nox_common_gameFlags_check_40A5C0(4096)) {
		v4 = 640.0;
	} else {
		v4 = 250.0;
	}
	if (v4 >= *(float*)(v3 + 1312)) {
		v21 = v4;
	} else {
		v21 = *(float*)(v3 + 1312);
	}
	if (gameFrame() - *(uint32_t*)(v3 + 1212) <= (unsigned int)(2 * gameFPS())) {
		v19 = 0;
	} else {
		v19 = 1;
		*(uint32_t*)(v3 + 1212) = gameFrame();
	}
	v5 = 0;
	if (*(uint8_t*)(v3 + 1129)) {
		v6 = (int*)(v3 + 1132);
		do {
			if (*(uint32_t*)(*v6 + 16) & 0x8020 || !nox_xxx_unitCanSee_536FB0(v1, *v6, 0) ||
				(v7 = *(float*)(v1 + 56) - *(float*)(*v6 + 56),
				 v8 = *(float*)(v1 + 60) - *(float*)(*v6 + 60), v20 = (v21 + 30.0) * (v21 + 30.0),
				 v8 * v8 + v7 * v7 > v20) ||
				(v9 = *(float*)(v1 + 56) - *(float*)(v1 + 72),
				 v10 = *(float*)(v1 + 60) - *(float*)(v1 + 76), v10 * v10 + v9 * v9 > 1000.0) ||
				v19 && !nox_xxx_unitCanInteractWith_5370E0(v1, *v6, 0)) {
				nox_xxx_aiLostSight_528560(v1, v5--);
				v17 = 1;
				--v6;
			}
			++v5;
			++v6;
		} while (v5 < *(unsigned char*)(v3 + 1129));
	}
	v11 = *(uint32_t*)(v3 + 1196);
	if (v11 && nox_xxx_testUnitBuffs_4FF350(v11, 28)) {
		v17 = 1;
	}
	if ((!*(uint32_t*)(v3 + 1196) ||
		 gameFrame() - *(uint32_t*)(v3 + 1204) > (unsigned int)(2 * gameFPS())) &&
		(*(uint32_t*)(v3 + 1208) <= gameFrame() ||
		 gameFrame() == *getMemU32Ptr(0x5D4594, 2487684))) {
		nox_xxx_unitsGetInCircle_517F90((float2*)(v1 + 56), v21, nox_xxx_monsterUpdateSeenEnemies_5286D0, v1);
		*(uint32_t*)(v3 + 1204) = gameFrame();
		*(uint32_t*)(v3 + 1212) = gameFrame();
		v17 = 1;
	}
	if (v17) {
		v12 = *(uint32_t*)(v3 + 1196);
		if (v12) {
			v13 = *(uint32_t*)(v12 + 36);
		} else {
			v13 = 0;
		}
		sub_528610(v1);
		v14 = *(uint32_t*)(v3 + 1196);
		if (v14 && v13 && v13 != *(uint32_t*)(v14 + 36)) {
			*(uint32_t*)(v3 + 1200) = v13;
		}
	}
	if (*(uint32_t*)(v3 + 1204) == gameFrame()) {
		v15 = *(uint32_t*)(v3 + 1440);
		if (v15 & 0x400 || nox_common_gameFlags_check_40A5C0(0x2000) || *(uint32_t*)(v3 + 1196)) {
			*(uint32_t*)(v3 + 1208) = gameFrame() + nox_common_randomInt_415FA0(5, 10);
		} else {
			v16 = 5 * gameFPS();
			v18 = sub_5336D0(v1);
			*(float*)(v3 + 524) = v18;
			if (v18 < 0.0) {
				*(uint32_t*)(v3 + 1208) = v16 + gameFrame();
			} else if (v18 > (double)v21) {
				*(uint32_t*)(v3 + 1208) = (unsigned long long)(long long)((v18 - v21) * (double)v16 / (1000.0 - v21)) +
										  10 + gameFrame();
			} else {
				*(uint32_t*)(v3 + 1208) = nox_common_randomInt_415FA0(5, 10) + gameFrame();
			}
		}
	}
}

//----- (00528560) --------------------------------------------------------
int nox_xxx_aiLostSight_528560(int a1, int a2) {
	int v2;     // esi
	int v3;     // eax
	int* v4;    // edi
	int v5;     // eax
	int v6;     // eax
	int v7;     // edx
	int v8;     // ecx
	int result; // eax
	int v10;    // [esp-4h] [ebp-14h]

	v2 = *(uint32_t*)(a1 + 748);
	v3 = *(uint32_t*)(v2 + 4 * a2 + 1132);
	v4 = (int*)(v2 + 4 * a2 + 1132);
	v10 = *(uint32_t*)(v3 + 36);
	v5 = nox_xxx_getUnitName_4E39D0(v3);
	nox_ai_debug_printf_5341A0("%d: Lost sight of %s(#%d)\n", gameFrame(), v5, v10);
	nox_xxx_scriptCallByEventBlock_502490((int*)(v2 + 1296), *v4, a1, 15);
	v6 = *(uint32_t*)(v2 + 1196);
	if (*v4 == v6) {
		v7 = *(uint32_t*)(v6 + 36);
		*(uint32_t*)(v2 + 1196) = 0;
		*(uint32_t*)(v2 + 1200) = v7;
	}
	v8 = a2;
	LOBYTE(result) = *(uint8_t*)(v2 + 1129) - 1;
	*(uint8_t*)(v2 + 1129) = result;
	result = (unsigned char)result;
	if (a2 < (unsigned char)result) {
		result = v2 + 4 * a2 + 1132;
		do {
			++v8;
			*(uint32_t*)result = *(uint32_t*)(result + 4);
			result += 4;
		} while (v8 < *(unsigned char*)(v2 + 1129));
	}
	return result;
}

//----- (00528610) --------------------------------------------------------
void sub_528610(int a1) {
	int v1;    // ebx
	int v2;    // esi
	int v3;    // ebp
	char v4;   // al
	int* i;    // edi
	double v6; // st7
	double v7; // st6
	double v8; // st5
	float v9;  // [esp+14h] [ebp+4h]

	v1 = a1;
	v2 = *(uint32_t*)(a1 + 748);
	v3 = 0;
	v9 = 100000000.0;
	v4 = *(uint8_t*)(v2 + 1129);
	*(uint32_t*)(v2 + 1196) = 0;
	if (v4) {
		for (i = (int*)(v2 + 1132); *i != *(uint32_t*)(v2 + 1216); ++i) {
			if (nox_xxx_unitIsEnemyTo_5330C0(v1, *i) && nox_xxx_checkIsKillable_528190(*i)) {
				v6 = *(float*)(*i + 56) - *(float*)(v1 + 56);
				v7 = *(float*)(*i + 60) - *(float*)(v1 + 60);
				v8 = v7 * v7 + v6 * v6;
				if (v8 < v9) {
					v9 = v8;
					*(uint32_t*)(v2 + 1196) = *i;
				}
			}
			if (++v3 >= *(unsigned char*)(v2 + 1129)) {
				return;
			}
		}
		*(uint32_t*)(v2 + 1196) = *(uint32_t*)(v2 + 1216);
	}
}

//----- (005286D0) --------------------------------------------------------
void nox_xxx_monsterUpdateSeenEnemies_5286D0(int a1, int a2) {
	int v2;    // esi
	int v3;    // ebx
	int v4;    // eax
	int v5;    // eax
	double v6; // st7
	double v7; // st6
	float* v8; // eax
	float v9;  // [esp+10h] [ebp-4h]
	float v10; // [esp+1Ch] [ebp+8h]

	v2 = a2;
	v3 = *(uint32_t*)(a2 + 748);
	if (a2 != a1) {
		if (*(uint8_t*)(a1 + 8) & 6) {
			if (!(*(uint32_t*)(a1 + 16) & 0x8020)) {
				v4 = *(uint32_t*)(v3 + 1440);
				if ((v4 & 0x400 || nox_xxx_unitIsEnemyTo_5330C0(a2, a1)) && !sub_528950(a2, a1)) {
					v5 = *(uint32_t*)(v3 + 1440);
					if (v5 & 0x100 ||
						(v6 = *(float*)(a1 + 56) - *(float*)(a2 + 56), v7 = *(float*)(a1 + 60) - *(float*)(a2 + 60),
						 v9 = v7, v8 = getMemFloatPtr(0x587000, 194136 + 8 * *(short*)(a2 + 124)),
						 v10 = sqrt(v7 * v9 + v6 * v6) + 0.001, v9 / v10 * v8[1] + v6 / v10 * *v8 >= 0.5)) {
						if (nox_xxx_unitCanInteractWith_5370E0(v2, a1, 0)) {
							nox_xxx_monsterVisionSeeEnemy_5287B0(v2, a1);
						}
					}
				}
			}
		}
	}
}

//----- (005287B0) --------------------------------------------------------
void nox_xxx_monsterVisionSeeEnemy_5287B0(int a1, int a2) {
	int v2;     // esi
	double v3;  // st7
	int v4;     // edi
	int v5;     // ebx
	int v6;     // ebp
	int v7;     // edx
	double v8;  // st6
	double v9;  // st5
	double v10; // st4
	int v11;    // ebx
	double v12;
	double v13;
	unsigned int v15; // eax
	int v16;          // eax
	int v17;          // eax
	float v18;        // [esp+14h] [ebp+4h]

	v2 = a1;
	v3 = 0.0;
	v4 = *(uint32_t*)(a1 + 748);
	v5 = 16;
	v6 = 0;
	if (*(uint8_t*)(v4 + 1129) == 16) {
		v7 = v4 + 1132;
		do {
			v8 = *(float*)(*(uint32_t*)v7 + 56) - *(float*)(v2 + 56);
			v9 = *(float*)(*(uint32_t*)v7 + 60) - *(float*)(v2 + 60);
			v10 = v9 * v9 + v8 * v8;
			if (v10 > v3) {
				v18 = v10;
				v3 = v18;
				v6 = *(uint32_t*)v7;
			}
			v7 += 4;
			--v5;
		} while (v5);
		v11 = a2;
		v12 = *(float*)(a2 + 56) - *(float*)(v2 + 56);
		v13 = *(float*)(a2 + 60) - *(float*)(v2 + 60);
		if (v3 <= v13 * v13 + v12 * v12) {
			return;
		}
		sub_528910(v2, v6);
	} else {
		v11 = a2;
	}
	*(uint32_t*)(v4 + 4 * *(unsigned char*)(v4 + 1129) + 1132) = v11;
	v15 = *(uint32_t*)(v4 + 536);
	++*(uint8_t*)(v4 + 1129);
	if (gameFrame() > v15) {
		if (nox_xxx_unitIsEnemyTo_5330C0(v2, v11)) {
			if (!nox_xxx_unitIsZombie_534A40(v2) || (v16 = *(uint32_t*)(v2 + 16), (v16 & 0x8000) == 0)) {
				v17 = nox_xxx_monsterGetSoundSet_424300(v2);
				if (v17) {
					nox_xxx_aud_501960(*(uint32_t*)(v17 + 68), v2, 0, 0);
				}
				*(uint32_t*)(v4 + 536) =
					gameFrame() + nox_common_randomInt_415FA0(2 * gameFPS(), 4 * gameFPS());
			}
		}
	}
	nox_xxx_scriptCallByEventBlock_502490((int*)(v4 + 1232), v11, v2, 14);
}

//----- (00528910) --------------------------------------------------------
int sub_528910(int a1, int a2) {
	int result;  // eax
	int v3;      // edx
	int v4;      // ecx
	uint32_t* i; // edx

	result = 0;
	v3 = *(uint32_t*)(a1 + 748);
	v4 = *(unsigned char*)(v3 + 1129);
	if (v4 > 0) {
		for (i = (uint32_t*)(v3 + 1132); *i != a2; ++i) {
			if (++result >= v4) {
				return result;
			}
		}
		result = nox_xxx_aiLostSight_528560(a1, result);
	}
	return result;
}

//----- (00528950) --------------------------------------------------------
int sub_528950(int a1, int a2) {
	int v2;      // ecx
	int v3;      // eax
	int v4;      // edx
	uint32_t* i; // ecx

	v2 = *(uint32_t*)(a1 + 748);
	v3 = 0;
	v4 = *(unsigned char*)(v2 + 1129);
	if (v4 <= 0) {
		return 0;
	}
	for (i = (uint32_t*)(v2 + 1132); *i != a2; ++i) {
		if (++v3 >= v4) {
			return 0;
		}
	}
	return 1;
}

//----- (00528990) --------------------------------------------------------
nox_object_t* nox_xxx_getFirstUpdatableObject_4DA8A0();
nox_object_t* nox_xxx_getNextUpdatableObject_4DA8B0(nox_object_t* obj);
int sub_528990(nox_object_t* a1) {
	int result; // eax
	int i;      // esi

	result = nox_xxx_getFirstUpdatableObject_4DA8A0();
	for (i = result; result; i = result) {
		if (*(uint8_t*)(i + 8) & 2) {
			if (!(*(uint8_t*)(i + 16) & 0x20)) {
				sub_528910(i, a1);
				sub_528610(i);
			}
		}
		result = nox_xxx_getNextUpdatableObject_4DA8B0(i);
	}
	return result;
}

//----- (005289D0) --------------------------------------------------------
void nox_xxx_netReportDestroyObject_5289D0(nox_object_t* a1p) {
	int a1 = a1p;
	char* result; // eax
	int i;        // edi
	int v4;       // [esp+0h] [ebp-4h]

	result = nox_common_playerInfoGetFirst_416EA0();
	for (i = (int)result; result; i = (int)result) {
		if ((1 << *(uint8_t*)(i + 2064)) & *(uint32_t*)(a1 + 148)) {
			LOBYTE(v4) = ((unsigned char)*(uint32_t*)(a1 + 20) >> 6) | 0x31;
			*(uint16_t*)((char*)&v4 + 1) = nox_xxx_netGetUnitCodeServ_578AC0((uint32_t*)a1);
			nox_xxx_netSendPacket0_4E5420(*(unsigned char*)(i + 2064), &v4, 3, 0, 1);
		}
		if (*(uint8_t*)(a1 + 8) & 6) {
			nox_xxx_netFriendAddRemove_4D97A0(*(unsigned char*)(i + 2064), (uint32_t*)a1, 0);
		}
		result = nox_common_playerInfoGetNext_416EE0(i);
	}
}

//----- (00528A60) --------------------------------------------------------
int nox_xxx_netObjectOutOfSight_528A60(int a1, uint32_t* a2) {
	char v4[3]; // [esp+0h] [ebp-4h]
	v4[0] = 50;
	*(uint16_t*)&v4[1] = nox_xxx_netGetUnitCodeServ_578AC0(a2);
	return nox_xxx_netSendPacket0_4E5420(a1, v4, 3, 0, 1);
}

//----- (00528A90) --------------------------------------------------------
int nox_xxx_netObjectInShadows_528A90(int a1, uint32_t* a2) {
	char v4[3]; // [esp+0h] [ebp-4h]
	v4[0] = 51;
	*(uint16_t*)&v4[1] = nox_xxx_netGetUnitCodeServ_578AC0(a2);
	return nox_xxx_netSendPacket0_4E5420(a1, v4, 3, 0, 1);
}

//----- (00528BD0) --------------------------------------------------------
int nox_xxx_monsterCmdSend_528BD0(int unit, int source, const char* command, short a4) {
	short v4;      // ax
	double v5;     // st7
	long long v6;  // rax
	double v7;     // st7
	int result;    // eax
	int i;         // esi
	char v10[520]; // [esp+Ch] [ebp-208h]

	v10[0] = -88; // MSG_TEXT_MESSAGE
	v10[3] = 8;
	v4 = nox_xxx_netGetUnitCodeServ_578AC0((uint32_t*)unit);
	v5 = *(float*)(unit + 56);
	*(uint16_t*)&v10[1] = v4;
	v6 = (long long)v5;
	v7 = *(float*)(unit + 60);
	*(uint16_t*)&v10[4] = v6;
	*(uint16_t*)&v10[6] = (long long)v7;
	*(uint16_t*)&v10[9] = a4;
	v10[8] = strlen(command) + 1;
	result = source;
	strcpy(&v10[11], command);
	if (source) {
		if (*(uint8_t*)(source + 8) & 4) { // if source is player / local player ?
			result = nox_netlist_addToMsgListCli_40EBC0(
				*(unsigned char*)(*(uint32_t*)(*(uint32_t*)(source + 748) + 276) + 2064), 1, v10,
				(unsigned char)v10[8] + 11);
		}
	} else {
		result = nox_xxx_getFirstPlayerUnit_4DA7C0();
		for (i = result; result; i = result) {
			nox_netlist_addToMsgListCli_40EBC0(*(unsigned char*)(*(uint32_t*)(*(uint32_t*)(i + 748) + 276) + 2064), 1,
											   v10, (unsigned char)v10[8] + 11);
			result = nox_xxx_getNextPlayerUnit_4DA7F0(i);
		}
	}
	return result;
}

//----- (00528D60) --------------------------------------------------------
int nox_xxx_destroyEveryChatMB_528D60() {
	int result; // eax
	int i;      // esi
	char v2[3]; // [esp+4h] [ebp-4h]

	v2[0] = -54;
	*(uint16_t*)&v2[1] = -8531;
	result = nox_xxx_getFirstPlayerUnit_4DA7C0();
	for (i = result; result; i = result) {
		nox_netlist_addToMsgListCli_40EBC0(*(unsigned char*)(*(uint32_t*)(*(uint32_t*)(i + 748) + 276) + 2064), 1, v2,
										   3);
		result = nox_xxx_getNextPlayerUnit_4DA7F0(i);
	}
	return result;
}

//----- (00528DB0) --------------------------------------------------------
int nox_xxx_XFerMonster_528DB0(nox_object_t* a1p) {
	int a1 = a1p;
	int v1;            // esi
	int result;        // eax
	char* v3;          // edi
	char* v4;          // eax
	char* v5;          // eax
	char* v6;          // eax
	char* v7;          // eax
	char* v8;          // eax
	char* v9;          // eax
	char* v10;         // eax
	char* v11;         // eax
	char* v12;         // eax
	char* v13;         // edi
	int v14;           // eax
	short v15;         // ax
	int v16;           // edx
	uint32_t* v17;     // edi
	uint32_t* v18;     // eax
	int v19;           // ecx
	int v20;           // ebx
	char* v21;         // ebp
	int i;             // edi
	int v23;           // eax
	int* v24;          // ebp
	char* v25;         // ebx
	char* v26;         // ebx
	int v27;           // ebp
	unsigned char v28; // bl
	uint8_t* v29;      // edi
	int v30;           // eax
	int v31;           // eax
	int v32;           // edi
	int v33;           // ebx
	int j;             // eax
	int v35;           // esi
	int v36;           // ecx
	uint32_t* v37;     // eax
	uint32_t* v38;     // eax
	int v39;           // edi
	int* v40;          // ecx
	uint32_t* v41;     // edx
	int v42;           // esi
	int v43;           // eax
	int v44;           // [esp+10h] [ebp-128h]
	int v45;           // [esp+14h] [ebp-124h]
	int v46;           // [esp+18h] [ebp-120h]
	uint8_t* v47;      // [esp+1Ch] [ebp-11Ch]
	int v48;           // [esp+20h] [ebp-118h]
	int v49;           // [esp+24h] [ebp-114h]
	int v50;           // [esp+28h] [ebp-110h]
	int v51;           // [esp+2Ch] [ebp-10Ch]
	uint32_t v52[2];   // [esp+30h] [ebp-108h]
	char v53[256];     // [esp+38h] [ebp-100h]

	v1 = *(uint32_t*)(a1 + 748);
	v51 = *(uint32_t*)(a1 + 136);
	if (!*getMemU32Ptr(0x5D4594, 2487692)) {
		*getMemU32Ptr(0x5D4594, 2487692) = nox_xxx_getNameId_4E3AA0("Glyph");
	}
	v45 = 64;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v45, 2u);
	if ((short)v45 > 64) {
		return 0;
	}
	result = nox_xxx_mapReadWriteObjData_4F4530((int*)a1, (short)v45);
	if (!result) {
		return result;
	}
	if (!nox_crypt_IsReadOnly()) {
		nox_xxx_xferIndexedDirection_509E20(*(short*)(a1 + 124), (int2*)v52);
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread(v52, 8u);
	if ((short)v45 >= 3) {
		v3 = *(char**)(a1 + 756);
		if (v3) {
			v4 = v3 + 640;
		} else {
			v4 = 0;
		}
		nox_xxx_xferReadScriptHandler_4F5580(v1 + 1232, v4);
		if (v3) {
			v5 = v3 + 896;
		} else {
			v5 = 0;
		}
		nox_xxx_xferReadScriptHandler_4F5580(v1 + 1264, v5);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1220), 2u);
		if (v3) {
			v6 = v3 + 768;
		} else {
			v6 = 0;
		}
		nox_xxx_xferReadScriptHandler_4F5580(v1 + 1224, v6);
		if ((short)v45 >= 31) {
			v7 = v3 ? v3 + 1024 : 0;
			nox_xxx_xferReadScriptHandler_4F5580(v1 + 1240, v7);
			v8 = v3 ? v3 + 1152 : 0;
			nox_xxx_xferReadScriptHandler_4F5580(v1 + 1248, v8);
			v9 = v3 ? v3 + 1280 : 0;
			nox_xxx_xferReadScriptHandler_4F5580(v1 + 1256, v9);
			v10 = v3 ? v3 + 1408 : 0;
			nox_xxx_xferReadScriptHandler_4F5580(v1 + 1272, v10);
			v11 = v3 ? v3 + 1536 : 0;
			nox_xxx_xferReadScriptHandler_4F5580(v1 + 1280, v11);
			v12 = v3 ? v3 + 1664 : 0;
			nox_xxx_xferReadScriptHandler_4F5580(v1 + 1288, v12);
			if ((short)v45 >= 52) {
				if (v3) {
					v13 = v3 + 1792;
				} else {
					v13 = 0;
				}
				nox_xxx_xferReadScriptHandler_4F5580(v1 + 1296, v13);
			}
		}
	} else {
		sub_4F5540(v1 + 1232);
		sub_4F5540(v1 + 1264);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1220), 2u);
		sub_4F5540(v1 + 1224);
	}
	v14 = nox_crypt_IsReadOnly();
	if (nox_crypt_IsReadOnly() != 1 ||
		(v15 = nox_xxx_xferDirectionToAngle_509E00(v52), *(uint16_t*)(a1 + 126) = v15, *(uint16_t*)(a1 + 124) = v15,
		 v14 = nox_crypt_IsReadOnly(), nox_crypt_IsReadOnly() != 1)) {
		if (!v14) {
			v46 = 0;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v46, 4u);
		}
	} else if ((short)v45 >= 11) {
		v46 = 0;
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v46, 4u);
	}
	if ((short)v45 >= 31) {
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1332), 1u);
		if ((short)v45 < 51) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v46, 2u);
			*(uint32_t*)(v1 + 1440) = (unsigned short)v46;
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1440), 4u);
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1352), 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1336), 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1344), 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1312), 4u);
		if ((short)v45 < 33) {
			nox_xxx_cryptSeekCur_40E0A0(2);
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1304), 4u);
		*(uint32_t*)(v1 + 1308) = *(uint32_t*)(v1 + 1304);
		if ((short)v45 < 34) {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1360), 4u);
		}
		LOBYTE(v48) = strlen((const char*)(v1 + 1364));
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v48, 1u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1364), (unsigned char)v48);
		*(uint8_t*)((unsigned char)v48 + v1 + 1364) = 0;
		if ((short)v45 >= 34) {
			if (nox_crypt_IsReadOnly()) {
				memset((void*)(v1 + 1488), 0, 0x224u);
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 4u);
				for (i = 0; i < v44; ++i) {
					nox_xxx_fileReadWrite_426AC0_file3_fread(&v48, 1u);
					nox_xxx_fileReadWrite_426AC0_file3_fread(v53, (unsigned char)v48);
					v53[(unsigned char)v48] = 0;
					v23 = nox_xxx_spellNameToN_4243F0(v53);
					nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 4 * v23 + 1488), 4u);
				}
			} else {
				v16 = 0;
				v17 = (uint32_t*)(v1 + 1488);
				v44 = 0;
				v18 = (uint32_t*)(v1 + 1488);
				v19 = 137;
				do {
					if (*v18) {
						++v16;
					}
					++v18;
					--v19;
				} while (v19);
				v44 = v16;
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 4u);
				v20 = 0;
				v47 = (uint8_t*)(v1 + 1488);
				do {
					if (*v17) {
						v21 = nox_xxx_spellNameByN_424870(v20);
						LOBYTE(v46) = strlen(v21);
						nox_xxx_fileReadWrite_426AC0_file3_fread(&v46, 1u);
						nox_xxx_fileReadWrite_426AC0_file3_fread(v21, (unsigned char)v46);
						nox_xxx_fileReadWrite_426AC0_file3_fread(v47, 4u);
					}
					++v20;
					v17 = v47 + 4;
					v47 += 4;
				} while (v20 < 137);
			}
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1488), 0x224u);
		}
		if ((short)v45 < 46) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
			*(uint16_t*)(v1 + 1448) = (unsigned char)v44;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
			*(uint16_t*)(v1 + 1450) = (unsigned char)v44;
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1448), 2u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1450), 2u);
		}
		if ((short)v45 <= 32) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v49, 4u);
		}
		if ((short)v45 < 46) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
			*(uint16_t*)(v1 + 1456) = (unsigned char)v44;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
			*(uint16_t*)(v1 + 1458) = (unsigned char)v44;
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1456), 2u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1458), 2u);
		}
		if ((short)v45 <= 32) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v49, 4u);
		}
		if ((short)v45 < 46) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
			*(uint16_t*)(v1 + 1464) = (unsigned char)v44;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
			*(uint16_t*)(v1 + 1466) = (unsigned char)v44;
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1464), 2u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1466), 2u);
		}
		if ((short)v45 <= 32) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v49, 4u);
		}
		if ((short)v45 < 46) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
			*(uint16_t*)(v1 + 1472) = (unsigned char)v44;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
			*(uint16_t*)(v1 + 1474) = (unsigned char)v44;
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1472), 2u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1474), 2u);
		}
		if ((short)v45 <= 32) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v49, 4u);
		}
		if ((short)v45 < 46) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
			*(uint16_t*)(v1 + 1480) = (unsigned char)v44;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
			*(uint16_t*)(v1 + 1482) = (unsigned char)v44;
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1480), 2u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1482), 2u);
		}
		if ((short)v45 > 32 || (nox_xxx_fileReadWrite_426AC0_file3_fread(&v49, 4u), (short)v45 >= 32)) {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1316), 4u);
		}
		if ((short)v45 >= 33) {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 2040), 4u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1320), 4u);
			if ((short)v45 < 42) {
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 2u);
				if (!(uint16_t)v44) {
					*(uint8_t*)(v1 + 1445) = 1;
				}
			}
			if ((short)v45 < 53) {
				nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 2044), 4u);
				nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 2048), 4u);
				nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 2052), 4u);
			} else {
				v24 = (int*)(v1 + 2044);
				v46 = 3;
				do {
					if (nox_crypt_IsReadOnly() == 1) {
						nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
						nox_xxx_fileReadWrite_426AC0_file3_fread(v53, (unsigned char)v44);
						v53[(unsigned char)v44] = 0;
						*v24 = nox_xxx_spellNameToN_4243F0(v53);
					} else {
						v25 = nox_xxx_spellNameByN_424870(*v24);
						LOBYTE(v44) = strlen(v25);
						nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
						nox_xxx_fileReadWrite_426AC0_file3_fread(v25, (unsigned char)v44);
					}
					++v24;
					--v46;
				} while (v46);
			}
		}
		if ((short)v45 >= 34) {
			if (nox_crypt_IsReadOnly()) {
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
				nox_xxx_fileReadWrite_426AC0_file3_fread(v53, (unsigned char)v44);
				v53[(unsigned char)v44] = 0;
				*(uint32_t*)(v1 + 1360) = nox_xxx_actionNByNameMB_5345F0(v53);
			} else {
				v26 = sub_5345B0(*(uint32_t*)(v1 + 1360));
				LOBYTE(v44) = strlen(v26);
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
				nox_xxx_fileReadWrite_426AC0_file3_fread(v26, (unsigned char)v44);
			}
		}
	}
	if ((short)v45 >= 41) {
		result = nox_xxx_XFer_ActionData_529CE0(a1);
		if (!result) {
			return result;
		}
	}
	if ((short)v45 >= 42) {
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1445), 1u);
	}
	if ((short)v45 >= 43 && *(uint8_t*)(a1 + 12) & 8) {
		LOBYTE(v44) = 0;
		v27 = *(uint32_t*)(a1 + 692);
		if ((short)v45 >= 50) {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v27 + 1716), 4u);
		}
		if ((short)v45 >= 61) {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v27 + 1720), 4u);
		}
		if ((short)v45 >= 48) {
			LOBYTE(v47) = strlen((const char*)(v27 + 1684));
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v47, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v27 + 1684), (unsigned char)v47);
			*(uint8_t*)((unsigned char)v47 + v27 + 1684) = 0;
		}
		if (nox_crypt_IsReadOnly() == 1) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
			if (!v27) {
				goto LABEL_137;
			}
			*(uint8_t*)v27 = v44;
		} else {
			if (v27) {
				LOBYTE(v44) = *(uint8_t*)v27;
			}
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 1u);
		}
		if (v27) {
			LOBYTE(v44) = 0;
			if (*(uint8_t*)v27) {
				do {
					if (nox_crypt_IsReadOnly() == 1) {
						nox_xxx_XFer_ReadShopItem_52A840(v27 + 28 * (unsigned char)v44 + 4, (short)v45);
					} else {
						nox_xxx_XFer_WriteShopItem_52A5F0(v27 + 28 * (unsigned char)v44 + 4);
					}
					LOBYTE(v44) = v44 + 1;
				} while ((unsigned char)v44 < *(uint8_t*)v27);
			}
		}
	}
LABEL_137:
	if ((short)v45 >= 44) {
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)v1, 4u);
	}
	if ((short)v45 >= 45) {
		v50 = *(uint32_t*)(a1 + 12) & 0x180;
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v50, 4u);
		*(uint32_t*)(a1 + 12) |= v50;
	}
	if ((short)v45 >= 49) {
		nox_xxx_fileReadWrite_426AC0_file3_fread(*(uint8_t**)(a1 + 556), 2u);
	}
	if ((short)v45 >= 51) {
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1348), 1u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1340), 1u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1444), 1u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 2036), 1u);
	}
	if (*(uint8_t*)(a1 + 12) & 0x20 && (short)v45 >= 54) {
		v28 = 0;
		v29 = (uint8_t*)(v1 + 2076);
		LOBYTE(v46) = 0;
		do {
			nox_xxx_fileReadWrite_426AC0_file3_fread(v29, 3u);
			if (nox_crypt_IsReadOnly() == 1) {
				nox_xxx_setNPCColor_4E4A90(a1, v46, (int)v29);
			}
			++v28;
			v29 += 3;
			LOBYTE(v46) = v28;
		} while (v28 < 6u);
	}
	if ((short)v45 >= 55 && *(uint8_t*)(a1 + 12) & 0x20) {
		nox_xxx_readNPCVoiceSet_52AD10(a1);
	}
	if (!((short)v45 < 62 || (result = nox_xxx_XFer_ReadMonsterBuffs_52AAB0((uint32_t*)a1)) != 0)) {
		return result;
	}
	if ((short)v45 >= 63 && *(uint32_t*)(a1 + 12) & 0x80000) {
		nox_xxx_readNPCVoiceSet_52AD10(a1);
	}
	if ((short)v45 >= 64) {
		LOBYTE(v47) = *(uint8_t*)(a1 + 540);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v47, 1u);
		if (nox_crypt_IsReadOnly() != 1) {
			goto LABEL_171;
		}
		if (!(uint8_t)v47) {
			goto LABEL_164;
		}
		nox_xxx_setSomePoisonData_4EEA90(a1, (unsigned char)v47);
	}
	if (nox_crypt_IsReadOnly() != 1) {
		goto LABEL_171;
	}
LABEL_164:
	if (!*(uint8_t*)(v1 + 1445)) {
		v30 = *(uint32_t*)(a1 + 16);
		if ((v30 & 0x8000) != 0 && !nox_xxx_unitIsZombie_534A40(a1)) {
			*(uint32_t*)(a1 + 16) |= 0x40u;
		}
		goto LABEL_171;
	}
	if (nox_common_gameFlags_check_40A5C0(1)) {
		*(uint16_t*)(*(uint32_t*)(a1 + 556) + 4) = 0;
		**(uint16_t**)(a1 + 556) = 0;
	}
	if (nox_crypt_IsReadOnly() == 1) {
		v30 = *(uint32_t*)(a1 + 16);
		if ((v30 & 0x8000) != 0 && !nox_xxx_unitIsZombie_534A40(a1)) {
			*(uint32_t*)(a1 + 16) |= 0x40u;
		}
	}
LABEL_171:
	if (*(uint32_t*)(a1 + 136)) {
		if (nox_crypt_IsReadOnly() != 1) {
			*(uint32_t*)(a1 + 136) = v51;
			return 1;
		}
		result = nox_xxx_xfer_4F3E30(v45, a1, *(uint32_t*)(a1 + 136));
		if (!result) {
			return result;
		}
	}
	if (nox_crypt_IsReadOnly() == 1) {
		if (nox_common_gameFlags_check_40A5C0(0x200000) || !nox_xxx_gameIsSwitchToSolo_4DB240()) {
			nox_xxx_monsterOnSpawnSpellcaster_529BC0(a1);
		}
		if (nox_crypt_IsReadOnly() == 1) {
			if (*(uint8_t*)(a1 + 8) & 2) {
				v31 = *(uint32_t*)(a1 + 12);
				if (v31 & 0x2000) {
					v32 = 0;
					v33 = 0;
					for (j = nox_xxx_inventoryGetFirst_4E7980(a1); j;
						 j = nox_xxx_inventoryGetNext_4E7990(j)) {
						if (*(unsigned short*)(j + 4) == *getMemU32Ptr(0x5D4594, 2487692)) {
							v32 = 1;
						}
					}
					v35 = v1 + 2044;
					v36 = 3;
					v37 = (uint32_t*)v35;
					do {
						if (*v37) {
							++v33;
						}
						++v37;
						--v36;
					} while (v36);
					if (!v32 && v33) {
						v38 = nox_xxx_newObjectByTypeID_4E3810("Glyph");
						v46 = (int)v38;
						if (v38) {
							v39 = v38[173];
							if (v33 > 0) {
								v40 = (int*)v35;
								v41 = (uint32_t*)v38[173];
								v42 = v33;
								do {
									v43 = *v40;
									++v40;
									*v41 = v43;
									++v41;
									--v42;
								} while (v42);
								v38 = (uint32_t*)v46;
							}
							*(uint8_t*)(v39 + 20) = v33;
							*(uint32_t*)(v39 + 24) = 0;
							*(uint32_t*)(v39 + 28) = *(uint32_t*)(a1 + 56);
							*(uint32_t*)(v39 + 32) = *(uint32_t*)(a1 + 60);
						}
						nox_xxx_inventoryPutImpl_4F3070(a1, (int)v38, 1);
					}
				}
			}
		}
	}
	*(uint32_t*)(a1 + 136) = v51;
	return 1;
}
// 528DB0: using guessed type char var_100[256];

//----- (00529BC0) --------------------------------------------------------
void nox_xxx_monsterOnSpawnSpellcaster_529BC0(int a1) {
	int v1;       // esi
	uint32_t* v2; // eax
	int i;        // ecx
	int v4;       // edx
	int v5;       // edi
	int v6;       // edi
	int v7;       // ecx

	if (a1) {
		v1 = *(uint32_t*)(a1 + 748);
		v2 = nox_xxx_monsterDefByTT_517560(*(unsigned short*)(a1 + 4));
		if (v2) {
			if (!*(uint8_t*)(v1 + 1445)) {
				**(uint16_t**)(a1 + 556) = *((uint16_t*)v2 + 34);
				*(uint16_t*)(*(uint32_t*)(a1 + 556) + 4) = *((uint16_t*)v2 + 34);
				*(uint16_t*)(*(uint32_t*)(a1 + 556) + 2) = *((uint16_t*)v2 + 34);
			}
			if (*(uint8_t*)(v1 + 1444) == 1) {
				*(uint32_t*)(v1 + 1440) = v2[23];
			} else {
				for (i = 0; i < 22; ++i) {
					v4 = 1 << i;
					if (!((1 << i) & 0x19C40)) {
						v5 = *(uint32_t*)(v1 + 1440);
						if (v5 & v4 && !(v4 & v2[23])) {
							*(uint32_t*)(v1 + 1440) = v5 & ~v4;
						}
						v6 = *(uint32_t*)(v1 + 1440);
						if (!(v6 & v4) && v4 & v2[23]) {
							*(uint32_t*)(v1 + 1440) = v4 | v6;
						}
					}
				}
			}
			v7 = *(uint32_t*)(v1 + 1440);
			if (!(v7 & 0x20)) {
				BYTE1(v7) &= 0xE7u;
				*(uint32_t*)(v1 + 1440) = v7;
			}
			if (*(uint8_t*)(v1 + 1340) == 1) {
				*(uint32_t*)(v1 + 1336) = v2[20];
			}
			if (*(uint8_t*)(v1 + 1348) == 1) {
				*(uint32_t*)(v1 + 1344) = v2[21];
			}
			if (*(uint8_t*)(v1 + 2036) == 1) {
				nox_xxx_monsterAutoSpells_54C0C0(a1);
			}
		}
	}
}

//----- (00529CE0) --------------------------------------------------------
nox_waypoint_t* sub_579C80(unsigned int a1);
int nox_xxx_XFer_ActionData_529CE0(int a1) {
	int v1;             // ebp
	int v3;             // ebx
	bool v4;            // cc
	uint8_t** v5;       // esi
	int v6;             // esi
	bool v7;            // zf
	uint8_t* v8;        // esi
	int v9;             // ecx
	int v10;            // eax
	int v11;            // eax
	int v12;            // eax
	int v13;            // eax
	int v14;            // ebx
	unsigned char* v15; // edi
	int* v16;           // ebp
	int v17;            // eax
	int v18;            // [esp+4h] [ebp-118h]
	char v19;           // [esp+Bh] [ebp-111h]
	int v20;            // [esp+Ch] [ebp-110h]
	int v21;            // [esp+10h] [ebp-10Ch]
	int v22;            // [esp+14h] [ebp-108h]
	int v23;            // [esp+18h] [ebp-104h]
	char v24[256];      // [esp+1Ch] [ebp-100h]

	v1 = *(uint32_t*)(a1 + 748);
	v21 = 4;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v21, 2u);
	if ((short)v21 > 4) {
		return 0;
	}
	v19 = 1;
	if ((short)v21 < 2) {
		goto LABEL_67;
	}
	v19 = 0;
	if (nox_common_gameFlags_check_40A5C0(1) && !nox_common_gameFlags_check_40A5C0(0x400000)) {
		v19 = 1;
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v19, 1u);
	if (!(v19 || (uint16_t)v21 == 1)) {
		return 1;
	}
LABEL_67:
	v23 = gameFrame();
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v23, 4u);
	v3 = gameFrame() - v23;
	v18 = 0;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v18, 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 8), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 12), 8 * *(uint32_t*)(v1 + 8));
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 268), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 272), 8u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 280), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 284), 1u);
	if (nox_crypt_IsReadOnly() == 1) {
		nox_xxx_AssignIfGreater_52A420((int*)(v1 + 280), v3);
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 296), 4u);
	v4 = *(int*)(v1 + 296) <= 0;
	v18 = 0;
	if (!v4) {
		v5 = (uint8_t**)(v1 + 300);
		do {
			if (nox_crypt_IsReadOnly()) {
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v22, 4u);
				*v5 = sub_579C80(v22);
			} else {
				nox_xxx_fileReadWrite_426AC0_file3_fread(*v5, 4u);
			}
			++v5;
			v4 = ++v18 < *(int*)(v1 + 296);
		} while (v4);
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 364), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 368), 8u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 376), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 380), 8u);
	strcpy(v24, nox_xxx_getSndName_40AF80(*(uint32_t*)(v1 + 388)));
	LOBYTE(v20) = strlen(v24);
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v20, 1u);
	nox_xxx_fileReadWrite_426AC0_file3_fread(v24, (unsigned char)v20);
	v24[(unsigned char)v20] = 0;
	*(uint32_t*)(v1 + 388) = nox_xxx_utilFindSound_40AF50(v24);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 396), 8u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 404), 4u);
	if (nox_crypt_IsReadOnly() == 1) {
		nox_xxx_AssignIfGreater_52A420((int*)(v1 + 404), v3);
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 481), 1u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 482), 1u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 483), 1u);
	if ((short)v21 < 3) {
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v18, 4u);
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 496), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 500), 8u);
	if (nox_crypt_IsReadOnly() == 1) {
		nox_xxx_AssignIfGreater_52A420((int*)(v1 + 496), v3);
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v18, 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 536), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 540), 4u);
	nox_xxx_AssignIfGreater_52A420((int*)(v1 + 536), v3);
	nox_xxx_AssignIfGreater_52A420((int*)(v1 + 540), v3);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 544), 1u);
	v6 = 0;
	if ((signed char)*(uint8_t*)(v1 + 544) >= 0) {
		v18 = v1 + 552;
		do {
			sub_52A440(a1, v18, v3);
			++v6;
			v18 += 24;
		} while (v6 <= *(char*)(v1 + 544));
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1129), 1u);
	v7 = *(uint8_t*)(v1 + 1129) == 0;
	v18 = 0;
	if (!v7) {
		v8 = (uint8_t*)(v1 + 1132);
		do {
			if (nox_crypt_IsReadOnly()) {
				nox_xxx_fileReadWrite_426AC0_file3_fread(v8, 4u);
			} else {
				if (!*getMemU32Ptr(0x5D4594, 2487688)) {
					*getMemU32Ptr(0x5D4594, 2487688) = nox_xxx_getNameId_4E3AA0("NewPlayer");
				}
				nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(*(uint32_t*)v8 + 44), 4u);
			}
			v9 = *(unsigned char*)(v1 + 1129);
			v8 += 4;
			++v18;
		} while (v18 < v9);
	}
	if (nox_crypt_IsReadOnly()) {
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1196), 4u);
	} else {
		v10 = *(uint32_t*)(v1 + 1196);
		if (v10) {
			v22 = *(uint32_t*)(v10 + 44);
		} else {
			v22 = 0;
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v22, 4u);
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1204), 4u);
	nox_xxx_AssignIfGreater_52A420((int*)(v1 + 1204), v3);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 2096), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 2100), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 2104), 1u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 2105), 1u);
	LOBYTE(v20) = strlen((const char*)(v1 + 2106));
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v20, 1u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 2106), (unsigned char)v20);
	*(uint8_t*)((unsigned char)v20 + v1 + 2106) = 0;
	if ((short)v21 < 4) {
		return 1;
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 4), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 288), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 292), 4u);
	v11 = nox_server_getObjectFromNetCode_4ECCB0(*(uint32_t*)(v1 + 392));
	if (v11) {
		v18 = *(uint32_t*)(v11 + 44);
	} else {
		v18 = 0;
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v18, 4u);
	if (nox_crypt_IsReadOnly() == 1) {
		*(uint32_t*)(v1 + 392) = v18;
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 492), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 508), 4u);
	nox_xxx_AssignIfGreater_52A420((int*)(v1 + 508), v3);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 512), 4u);
	nox_xxx_AssignIfGreater_52A420((int*)(v1 + 512), v3);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 516), 4u);
	nox_xxx_AssignIfGreater_52A420((int*)(v1 + 516), v3);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 520), 4u);
	nox_xxx_AssignIfGreater_52A420((int*)(v1 + 520), v3);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 528), 4u);
	nox_xxx_AssignIfGreater_52A420((int*)(v1 + 528), v3);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 532), 4u);
	nox_xxx_AssignIfGreater_52A420((int*)(v1 + 532), v3);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 524), 4u);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 548), 4u);
	nox_xxx_AssignIfGreater_52A420((int*)(v1 + 548), v3);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1128), 1u);
	v12 = nox_server_getObjectFromNetCode_4ECCB0(*(uint32_t*)(v1 + 1200));
	if (v12) {
		v18 = *(uint32_t*)(v12 + 44);
	} else {
		v18 = 0;
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v18, 4u);
	if (nox_crypt_IsReadOnly() == 1) {
		*(uint32_t*)(v1 + 1200) = v18;
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1208), 4u);
	nox_xxx_AssignIfGreater_52A420((int*)(v1 + 1208), v3);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1212), 4u);
	nox_xxx_AssignIfGreater_52A420((int*)(v1 + 1212), v3);
	v13 = *(uint32_t*)(v1 + 1216);
	v14 = 0;
	if (v13) {
		v18 = *(uint32_t*)(v13 + 44);
	} else {
		v18 = 0;
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v18, 4u);
	if (nox_crypt_IsReadOnly() == 1) {
		*(uint32_t*)(v1 + 1216) = v18;
	}
	v15 = (unsigned char*)(v1 + 2172);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 2172), 1u);
	if (*(uint8_t*)(v1 + 2172)) {
		v16 = (int*)(v1 + 2140);
		do {
			v17 = nox_server_getObjectFromNetCode_4ECCB0(*v16);
			if (v17) {
				v18 = *(uint32_t*)(v17 + 44);
			} else {
				v18 = 0;
			}
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v18, 4u);
			if (nox_crypt_IsReadOnly() == 1) {
				*v16 = v18;
			}
			++v14;
			++v16;
		} while (v14 < *v15);
	}
	return 1;
}
// 529CE0: using guessed type char var_100[256];

//----- (0052A420) --------------------------------------------------------
int nox_xxx_AssignIfGreater_52A420(int* a1, int a2) {
	int result; // eax

	result = a2 + *a1;
	if (result >= 1) {
		*a1 = result;
	} else {
		*a1 = 1;
	}
	return result;
}

//----- (0052A440) --------------------------------------------------------
int sub_52A440(int a1, int a2, int a3) {
	int v3;        // eax
	int v4;        // edi
	int* v5;       // esi
	int result;    // eax
	size_t v7;     // [esp-4h] [ebp-120h]
	int v8;        // [esp+10h] [ebp-10Ch]
	int v9;        // [esp+14h] [ebp-108h]
	int v10;       // [esp+18h] [ebp-104h]
	char v11[256]; // [esp+1Ch] [ebp-100h]

	strcpy(v11, sub_534650(*(uint32_t*)a2));
	LOBYTE(v10) = strlen(v11);
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v10, 1u);
	nox_xxx_fileReadWrite_426AC0_file3_fread(v11, (unsigned char)v10);
	v11[(unsigned char)v10] = 0;
	v3 = nox_xxx_actionByName_534670(v11);
	*(uint32_t*)a2 = v3;
	LOBYTE(v8) = getMemByte(0x587000, 255604 + 16 * v3);
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v8, 1u);
	v4 = 0;
	if ((unsigned char)v8 > 0u) {
		v5 = (int*)(a2 + 4);
		while (1) {
			result = *getMemU32Ptr(0x587000, 255608 + 4 * (v4 + 4 * *(uint32_t*)a2));
			switch (result) {
			case 0:
				v7 = 8;
				nox_xxx_fileReadWrite_426AC0_file3_fread(v5, v7);
				break;
			case 1:
				if (nox_crypt_IsReadOnly()) {
					v7 = 4;
					nox_xxx_fileReadWrite_426AC0_file3_fread(v5, v7);
					break;
				}
				if (*v5) {
					nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(*v5 + 44), 4u);
				} else {
					v9 = 0;
					nox_xxx_fileReadWrite_426AC0_file3_fread(&v9, 4u);
				}
				break;
			case 2:
				if (nox_crypt_IsReadOnly()) {
					v7 = 4;
					nox_xxx_fileReadWrite_426AC0_file3_fread(v5, v7);
					break;
				}
				if (*v5) {
					nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)*v5, 4u);
				} else {
					v9 = 0;
					nox_xxx_fileReadWrite_426AC0_file3_fread(&v9, 4u);
				}
				break;
			case 3:
			case 4:
			case 6:
				v7 = 4;
				nox_xxx_fileReadWrite_426AC0_file3_fread(v5, v7);
				break;
			case 5:
				nox_xxx_fileReadWrite_426AC0_file3_fread(v5, 4u);
				if (nox_crypt_IsReadOnly() == 1) {
					nox_xxx_AssignIfGreater_52A420(v5, a3);
				}
				break;
			case 7:
				v7 = 1;
				nox_xxx_fileReadWrite_426AC0_file3_fread(v5, v7);
				break;
			default:
				return result;
			}
			++v4;
			v5 += 2;
			if (v4 >= (unsigned char)v8) {
				return nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(a2 + 20), 4u);
			}
		}
	}
	return nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(a2 + 20), 4u);
}
// 52A440: using guessed type char var_100[256];

//----- (0052AAB0) --------------------------------------------------------
int nox_xxx_XFer_ReadMonsterBuffs_52AAB0(uint32_t* a1) {
	int v1;   // ebp
	char* v2; // ebx
	int v3;   // eax
	int v5;   // ebx
	int v6;   // edi
	int v7;   // ecx
	int v8;   // eax
	int v9;   // eax
	int v10;  // [esp-4h] [ebp-138h]
	int v11;  // [esp+10h] [ebp-124h]
	int v12;  // [esp+14h] [ebp-120h]
	int v13;  // [esp+18h] [ebp-11Ch]
	int v14;  // [esp+1Ch] [ebp-118h]
	int v15;  // [esp+20h] [ebp-114h]
	int v16;  // [esp+24h] [ebp-110h]
	int v17[3];
	char v20[256]; // [esp+34h] [ebp-100h]

	v13 = 2;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v13, 2u);
	if ((short)v13 > 2 || (short)v13 <= 0) {
		return 0;
	}
	LOBYTE(v16) = sub_424CB0((int)a1);
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v16, 1u);
	if (nox_crypt_IsReadOnly()) {
		v5 = 0;
		if (!(uint8_t)v16) {
			return 1;
		}
		while (1) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v11, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v20, (unsigned char)v11);
			v20[(unsigned char)v11] = 0;
			v6 = nox_xxx_enchantByName_424880(v20);
			if (v6 == -1) {
				break;
			}
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v15, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v14, 4u);
			v7 = a1[15];
			v17[1] = a1[14];
			v10 = (unsigned char)v15;
			v17[0] = a1;
			v17[2] = v7;
			v8 = nox_xxx_getEnchantSpell_424920(v6);
			nox_xxx_spellAccept_4FD400(v8, (int)a1, a1, (int)a1, v17, v10);
			*((uint16_t*)a1 + v6 + 172) = v14;
			if (v6 == 26 && (short)v13 >= 2) {
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v12, 4u);
				v9 = sub_4FF2D0(51, (int)a1);
				if (v9) {
					*(uint32_t*)(v9 + 72) = v12;
				}
			}
			if (++v5 >= (unsigned char)v16) {
				return 1;
			}
		}
		return 0;
	}
	v1 = sub_424D00();
	if (v1 == -1) {
		return 1;
	}
	do {
		if (nox_xxx_testUnitBuffs_4FF350((int)a1, v1)) {
			v2 = nox_xxx_getEnchantName_4248F0(v1);
			LOBYTE(v11) = strlen(v2);
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v11, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v2, (unsigned char)v11);
			LOBYTE(v15) = nox_xxx_buffGetPower_4FF570((int)a1, v1);
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v15, 1u);
			v14 = nox_xxx_unitGetBuffTimer_4FF550((int)a1, v1);
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v14, 4u);
			if (v1 == 26) {
				v3 = sub_4FF2D0(51, (int)a1);
				if (v3) {
					v12 = *(uint32_t*)(v3 + 72);
				} else {
					v12 = 100;
				}
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v12, 4u);
			}
		}
		v1 = sub_424D20(v1);
	} while (v1 != -1);
	return 1;
}
// 52AAB0: using guessed type char var_100[256];

//----- (0052AD10) --------------------------------------------------------
size_t nox_xxx_readNPCVoiceSet_52AD10(int a1) {
	const char** v1; // eax
	const char** v2; // esi
	size_t result;   // eax
	const char** v4; // eax
	int v5;          // [esp+4h] [ebp-104h]
	char v6[256];    // [esp+8h] [ebp-100h]

	if (nox_crypt_IsReadOnly()) {
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v5, 1u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(v6, (unsigned char)v5);
		v6[(unsigned char)v5] = 0;
		v4 = nox_xxx_getDefaultSoundSet_424350(v6);
		result = nox_xxx_setNPCVoiceSet_424320(a1, (int)v4);
	} else {
		v1 = (const char**)nox_xxx_monsterGetSoundSet_424300(a1);
		v2 = v1;
		if (v1) {
			LOBYTE(v5) = strlen(*v1);
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v5, 1u);
			result = nox_xxx_fileReadWrite_426AC0_file3_fread(*v2, (unsigned char)v5);
		} else {
			LOBYTE(v5) = 0;
			result = nox_xxx_fileReadWrite_426AC0_file3_fread(&v5, 1u);
		}
	}
	return result;
}
// 52AD10: using guessed type char var_100[256];

//----- (0052ADE0) --------------------------------------------------------
int nox_xxx_XFerNPC_52ADE0(nox_object_t* a1p) {
	int a1 = a1p;
	int v1;            // esi
	char* v2;          // edi
	int result;        // eax
	uint32_t* v4;      // eax
	char* v5;          // eax
	char* v6;          // eax
	char* v7;          // eax
	char* v8;          // eax
	char* v9;          // eax
	char* v10;         // eax
	char* v11;         // eax
	char* v12;         // eax
	char* v13;         // eax
	char* v14;         // edi
	int v15;           // eax
	short v16;         // ax
	unsigned char v17; // bl
	uint8_t* v18;      // edi
	int v19;           // ebx
	unsigned char i;   // bl
	int v21;           // edx
	uint16_t* v22;     // eax
	uint16_t* v23;     // eax
	int v24;           // edx
	uint8_t* v25;      // ebx
	uint32_t* v26;     // eax
	int v27;           // ecx
	char* v28;         // ebp
	bool v29;          // zf
	unsigned char j;   // bl
	int v31;           // eax
	int* v32;          // ebp
	char* v33;         // ebx
	char* v34;         // ebx
	int v35;           // edi
	int v36;           // eax
	int v37;           // eax
	int v38;           // eax
	int v39;           // eax
	uint32_t* k;       // esi
	int v41;           // eax
	int v42;           // eax
	int v43;           // [esp+10h] [ebp-238h]
	int v44;           // [esp+14h] [ebp-234h]
	int v45;           // [esp+18h] [ebp-230h]
	unsigned char v46; // [esp+1Fh] [ebp-229h]
	int v47;           // [esp+20h] [ebp-228h]
	int v48;           // [esp+24h] [ebp-224h]
	int v49;           // [esp+28h] [ebp-220h]
	int v50;           // [esp+2Ch] [ebp-21Ch]
	int v51;           // [esp+30h] [ebp-218h]
	int v52;           // [esp+34h] [ebp-214h]
	char v53[3];       // [esp+38h] [ebp-210h]
	int v54;           // [esp+3Ch] [ebp-20Ch]
	uint32_t v55[2];   // [esp+40h] [ebp-208h]
	char v56[256];     // [esp+48h] [ebp-200h]
	char v57[256];     // [esp+148h] [ebp-100h]

	v1 = a1p->data_update;
	v2 = *(char**)(a1 + 756);
	v54 = *(uint32_t*)(a1 + 136);
	v44 = 62;
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v44, 2u);
	if ((short)v44 > 62) {
		return 0;
	}
	if (nox_crypt_IsReadOnly() == 1) {
		v4 = nox_xxx_monsterDefByTT_517560(*(unsigned short*)(a1 + 4));
		*(uint32_t*)(v1 + 484) = v4;
		if (v4) {
			*(uint32_t*)(v1 + 1440) = v4[23];
		}
	}
	result = nox_xxx_mapReadWriteObjData_4F4530((int*)a1, (short)v44);
	if (!result) {
		return 0;
	}
	if (!nox_crypt_IsReadOnly()) {
		nox_xxx_xferIndexedDirection_509E20(*(short*)(a1 + 124), (int2*)v55);
	}
	nox_xxx_fileReadWrite_426AC0_file3_fread(v55, 8u);
	if (v2) {
		v5 = v2 + 640;
	} else {
		v5 = 0;
	}
	nox_xxx_xferReadScriptHandler_4F5580(v1 + 1232, v5);
	if (v2) {
		v6 = v2 + 896;
	} else {
		v6 = 0;
	}
	nox_xxx_xferReadScriptHandler_4F5580(v1 + 1264, v6);
	nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1220), 2u);
	if (v2) {
		v7 = v2 + 768;
	} else {
		v7 = 0;
	}
	nox_xxx_xferReadScriptHandler_4F5580(v1 + 1224, v7);
	if ((short)v44 >= 32) {
		v8 = v2 ? v2 + 1024 : 0;
		nox_xxx_xferReadScriptHandler_4F5580(v1 + 1240, v8);
		v9 = v2 ? v2 + 1152 : 0;
		nox_xxx_xferReadScriptHandler_4F5580(v1 + 1248, v9);
		v10 = v2 ? v2 + 1280 : 0;
		nox_xxx_xferReadScriptHandler_4F5580(v1 + 1256, v10);
		v11 = v2 ? v2 + 1408 : 0;
		nox_xxx_xferReadScriptHandler_4F5580(v1 + 1272, v11);
		v12 = v2 ? v2 + 1536 : 0;
		nox_xxx_xferReadScriptHandler_4F5580(v1 + 1280, v12);
		v13 = v2 ? v2 + 1664 : 0;
		nox_xxx_xferReadScriptHandler_4F5580(v1 + 1288, v13);
		if ((short)v44 >= 50) {
			if (v2) {
				v14 = v2 + 1792;
			} else {
				v14 = 0;
			}
			nox_xxx_xferReadScriptHandler_4F5580(v1 + 1296, v14);
		}
	}
	v47 = 0; // FIXME: set to direction? was uninitialized
	v15 = nox_crypt_IsReadOnly();
	if (nox_crypt_IsReadOnly() != 1 ||
		(v16 = nox_xxx_xferDirectionToAngle_509E00(v55), *(uint16_t*)(a1 + 126) = v16, *(uint16_t*)(a1 + 124) = v16,
		 v15 = nox_crypt_IsReadOnly(), nox_crypt_IsReadOnly() != 1)) {
		if (!v15) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v47, 4u);
		}
	} else {
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v47, 4u);
	}
	if (nox_crypt_IsReadOnly() == 1) {
		v17 = 0;
		LOBYTE(v48) = 0;
		do {
			nox_xxx_fileReadWrite_426AC0_file3_fread(v53, 3u);
			nox_xxx_setNPCColor_4E4A90(a1, v48, v53);
			LOBYTE(v48) = ++v17;
		} while (v17 < 6u);
	} else {
		v18 = (uint8_t*)(v1 + 2076);
		v19 = 6;
		do {
			nox_xxx_fileReadWrite_426AC0_file3_fread(v18, 3u);
			v18 += 3;
			--v19;
		} while (v19);
	}
	if (nox_crypt_IsReadOnly() == 1 && (uint16_t)v44 == 31) {
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v49, 2u);
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v46, 1u);
		for (i = 0; i < v46; ++i) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v50, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(v56, (unsigned char)v50);
			v56[(unsigned char)v50] = 0;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
		}
	}
	if ((short)v44 >= 32) {
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1332), 1u);
		if ((short)v44 < 49) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v47, 2u);
			*(uint32_t*)(v1 + 1440) = (unsigned short)v47;
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1440), 4u);
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1352), 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1336), 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1344), 4u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1312), 4u);
		v43 = 0;
		v22 = *(uint16_t**)(a1 + 556);
		if (v22) {
			LOWORD(v21) = *v22;
			v43 = v21;
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 2u);
		v23 = *(uint16_t**)(a1 + 556);
		if (v23) {
			*v23 = v43;
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1304), 4u);
		*(uint32_t*)(v1 + 1308) = *(uint32_t*)(v1 + 1304);
		if ((short)v44 < 35) {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1360), 4u);
		}
		LOBYTE(v49) = strlen((const char*)(v1 + 1364));
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v49, 1u);
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1364), (unsigned char)v49);
		*(uint8_t*)((unsigned char)v49 + v1 + 1364) = 0;
		if ((short)v44 >= 34) {
			if (nox_crypt_IsReadOnly()) {
				memset((void*)(v1 + 1488), 0, 0x224u);
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 4u);
				for (j = 0; j < v43; LOBYTE(v48) = j) {
					nox_xxx_fileReadWrite_426AC0_file3_fread(&v49, 1u);
					nox_xxx_fileReadWrite_426AC0_file3_fread(v56, (unsigned char)v49);
					v56[(unsigned char)v49] = 0;
					v31 = nox_xxx_spellNameToN_4243F0(v56);
					nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 4 * v31 + 1488), 4u);
					++j;
				}
			} else {
				v24 = 0;
				v25 = (uint8_t*)(v1 + 1488);
				v43 = 0;
				v26 = (uint32_t*)(v1 + 1488);
				v27 = 137;
				do {
					if (*v26) {
						++v24;
					}
					++v26;
					--v27;
				} while (v27);
				v43 = v24;
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 4u);
				v48 = 0;
				v47 = 137;
				do {
					if (*(uint32_t*)v25) {
						v28 = nox_xxx_spellNameByN_424870(v48);
						LOBYTE(v52) = strlen(v28);
						nox_xxx_fileReadWrite_426AC0_file3_fread(&v52, 1u);
						nox_xxx_fileReadWrite_426AC0_file3_fread(v28, (unsigned char)v52);
						nox_xxx_fileReadWrite_426AC0_file3_fread(v25, 4u);
					}
					v25 += 4;
					v29 = v47 == 1;
					++v48;
					--v47;
				} while (!v29);
			}
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1488), 0x224u);
		}
		if ((short)v44 < 47) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
			*(uint16_t*)(v1 + 1448) = (unsigned char)v43;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
			*(uint16_t*)(v1 + 1450) = (unsigned char)v43;
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1448), 2u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1450), 2u);
		}
		if ((short)v44 < 34) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v45, 4u);
		}
		if ((short)v44 < 47) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
			*(uint16_t*)(v1 + 1456) = (unsigned char)v43;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
			*(uint16_t*)(v1 + 1458) = (unsigned char)v43;
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1456), 2u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1458), 2u);
		}
		if ((short)v44 < 34) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v45, 4u);
		}
		if ((short)v44 < 47) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
			*(uint16_t*)(v1 + 1464) = (unsigned char)v43;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
			*(uint16_t*)(v1 + 1466) = (unsigned char)v43;
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1464), 2u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1466), 2u);
		}
		if ((short)v44 < 34) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v45, 4u);
		}
		if ((short)v44 < 47) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
			*(uint16_t*)(v1 + 1472) = (unsigned char)v43;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
			*(uint16_t*)(v1 + 1474) = (unsigned char)v43;
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1472), 2u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1474), 2u);
		}
		if ((short)v44 < 34) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v45, 4u);
		}
		if ((short)v44 < 47) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
			*(uint16_t*)(v1 + 1480) = (unsigned char)v43;
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
			*(uint16_t*)(v1 + 1482) = (unsigned char)v43;
		} else {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1480), 2u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1482), 2u);
		}
		if ((short)v44 < 34) {
			nox_xxx_fileReadWrite_426AC0_file3_fread(&v45, 4u);
		}
		if ((short)v44 >= 33) {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1316), 4u);
		}
		if ((short)v44 >= 34) {
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 2040), 4u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1324), 1u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1328), 4u);
			nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1320), 4u);
			if ((short)v44 < 42) {
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 2u);
				if (!(uint16_t)v43) {
					*(uint8_t*)(v1 + 1445) = 1;
				}
			}
			v32 = (int*)(v1 + 2044);
			v47 = 3;
			do {
				if (nox_crypt_IsReadOnly() == 1) {
					nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
					nox_xxx_fileReadWrite_426AC0_file3_fread(v56, (unsigned char)v43);
					v56[(unsigned char)v43] = 0;
					*v32 = nox_xxx_spellNameToN_4243F0(v56);
				} else {
					v33 = nox_xxx_spellNameByN_424870(*v32);
					LOBYTE(v43) = strlen(v33);
					nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
					nox_xxx_fileReadWrite_426AC0_file3_fread(v33, (unsigned char)v43);
				}
				++v32;
				--v47;
			} while (v47);
		}
		if ((short)v44 >= 35) {
			if (nox_crypt_IsReadOnly()) {
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
				nox_xxx_fileReadWrite_426AC0_file3_fread(v57, (unsigned char)v43);
				v57[(unsigned char)v43] = 0;
				*(uint32_t*)(v1 + 1360) = nox_xxx_actionNByNameMB_5345F0(v57);
			} else {
				v34 = sub_5345B0(*(uint32_t*)(v1 + 1360));
				LOBYTE(v43) = strlen(v34);
				nox_xxx_fileReadWrite_426AC0_file3_fread(&v43, 1u);
				nox_xxx_fileReadWrite_426AC0_file3_fread(v34, (unsigned char)v43);
			}
		}
	}
	if ((short)v44 < 41) {
		v35 = a1;
	} else {
		v35 = a1;
		result = nox_xxx_XFer_ActionData_529CE0(a1);
		if (!result) {
			return result;
		}
	}
	if ((short)v44 >= 42) {
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v1 + 1445), 1u);
	}
	if ((short)v44 >= 44) {
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)v1, 4u);
	}
	if ((short)v44 >= 45) {
		v36 = *(uint32_t*)(v35 + 556);
		v45 = 0;
		if (v36) {
			v45 = *(unsigned short*)(v36 + 4);
		}
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v45, 4u);
		v37 = *(uint32_t*)(v35 + 556);
		if (v37) {
			*(uint16_t*)(v37 + 4) = v45;
		}
	}
	if ((short)v44 >= 46) {
		v51 = *(uint32_t*)(v35 + 12) & 0x180;
		nox_xxx_fileReadWrite_426AC0_file3_fread(&v51, 4u);
		*(uint32_t*)(v35 + 12) |= v51;
	}
	if ((short)v44 >= 48) {
		nox_xxx_fileReadWrite_426AC0_file3_fread(*(uint8_t**)(v35 + 556), 2u);
	}
	if ((short)v44 >= 51) {
		nox_xxx_fileReadWrite_426AC0_file3_fread((uint8_t*)(v35 + 28), 4u);
	}
	if ((short)v44 >= 52) {
		nox_xxx_readNPCVoiceSet_52AD10(v35);
	}
	if (!((short)v44 < 61 || (result = nox_xxx_XFer_ReadMonsterBuffs_52AAB0((uint32_t*)v35)) != 0)) {
		return result;
	}
	if ((short)v44 < 62) {
		if (nox_crypt_IsReadOnly() == 1) {
			goto LABEL_149;
		}
		goto LABEL_156;
	}
	LOBYTE(v45) = *(uint8_t*)(v35 + 540);
	nox_xxx_fileReadWrite_426AC0_file3_fread(&v45, 1u);
	if (nox_crypt_IsReadOnly() != 1) {
		goto LABEL_156;
	}
	if ((uint8_t)v45) {
		nox_xxx_setSomePoisonData_4EEA90(v35, (unsigned char)v45);
		if (nox_crypt_IsReadOnly() == 1) {
			goto LABEL_149;
		}
		goto LABEL_156;
	}
LABEL_149:
	if (!*(uint8_t*)(v1 + 1445)) {
		goto LABEL_170;
	}
	if (nox_common_gameFlags_check_40A5C0(1)) {
		v38 = *(uint32_t*)(v35 + 556);
		if (v38) {
			*(uint16_t*)(v38 + 4) = 0;
			**(uint16_t**)(v35 + 556) = 0;
		}
	}
	if (nox_crypt_IsReadOnly() != 1) {
		goto LABEL_156;
	}
LABEL_170:
	v39 = *(uint32_t*)(v35 + 16);
	if ((v39 & 0x8000) != 0) {
		LOBYTE(v39) = v39 | 0x40;
		*(uint32_t*)(v35 + 16) = v39;
	}
LABEL_156:
	if (!*(uint32_t*)(v35 + 136) || nox_crypt_IsReadOnly() != 1 ||
		(result = nox_xxx_xfer_4F3E30(v44, v35, *(uint32_t*)(v35 + 136))) != 0) {
		sub_52BA70(v35);
		if (nox_crypt_IsReadOnly() == 1) {
			if (nox_common_gameFlags_check_40A5C0(1)) {
				for (k = *(uint32_t**)(v35 + 504); k; k = (uint32_t*)k[124]) {
					v41 = k[4];
					if (v41 & 0x100) {
						v42 = k[2];
						k[4] &= 0xFFFFFEFF;
						if (v42 & 0x1001000) {
							nox_xxx_NPCEquipWeapon_53A2C0(v35, (int)k);
						} else {
							nox_xxx_NPCEquipArmor_53E520(v35, k);
						}
					}
				}
			}
		}
		*(uint32_t*)(v35 + 136) = v54;
		return 1;
	}
	return result;
}
// 52B1CD: variable 'v21' is possibly undefined
// 52ADE0: using guessed type char var_200[256];
// 52ADE0: using guessed type char var_100[256];

//----- (0052BA70) --------------------------------------------------------
int sub_52BA70(int a1) {
	int result; // eax
	int v2;     // edi
	int v3;     // esi
	int v4;     // ecx
	int v5;     // edx

	result = a1;
	v2 = 0;
	v3 = 0;
	if (a1) {
		for (result = *(uint32_t*)(a1 + 504); result; result = *(uint32_t*)(result + 496)) {
			v4 = *(uint32_t*)(result + 16);
			if (v4 & 0x100) {
				v5 = *(uint32_t*)(result + 8);
				if (v5 & 0x1001000 && *(uint32_t*)(result + 12) & 0x7FFE40C) {
					if (v3 == 1) {
						BYTE1(v4) &= 0xFEu;
						*(uint32_t*)(result + 16) = v4;
					} else {
						v2 = 1;
					}
				} else if (v5 & 0x2000000 && *(uint8_t*)(result + 12) & 2) {
					if (v2 == 1) {
						BYTE1(v4) &= 0xFEu;
						*(uint32_t*)(result + 16) = v4;
					} else {
						v3 = 1;
					}
				}
			}
		}
	}
	return result;
}

//----- (0052BAF0) --------------------------------------------------------
int sub_52BAF0(int a1) {
	int v1;       // eax
	int result;   // eax
	uint32_t* v3; // ebx
	int v4;       // edi
	int v5;       // ebp
	int* v6;      // esi
	int v7;       // [esp+0h] [ebp-4h]
	int v8;       // [esp+8h] [ebp+4h]

	v1 = a1;
	v8 = 0;
	result = *(uint32_t*)(v1 + 748);
	v7 = result;
	if (*(uint8_t*)(result + 544) & 0x80) {
		return result;
	}
	v3 = (uint32_t*)(result + 552);
	do {
		v4 = 0;
		v5 = *getMemU32Ptr(0x587000, 255604 + 16 * *v3);
		if (v5 <= 0) {
			goto LABEL_14;
		}
		v6 = v3 + 1;
		do {
			if (*getMemU32Ptr(0x587000, 255608 + 4 * (v4 + 4 * *v3)) == 1) {
				if (*v6) {
					*v6 = sub_4ECF10(*v6);
					goto LABEL_12;
				}
			} else {
				if (*getMemU32Ptr(0x587000, 255608 + 4 * (v4 + 4 * *v3)) != 2) {
					goto LABEL_12;
				}
				if (*v6) {
					*v6 = nox_server_getWaypointById_579C40(*v6);
					goto LABEL_12;
				}
			}
			*v6 = 0;
		LABEL_12:
			++v4;
			v6 += 2;
		} while (v4 < v5);
		result = v7;
	LABEL_14:
		v3 += 6;
		++v8;
	} while (v8 <= *(char*)(result + 544));
	return result;
}
