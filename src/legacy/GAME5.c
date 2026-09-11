#include <float.h>
#include <math.h>
#include <stdio.h>

#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME1_3.h"
#include "GAME2_2.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_1.h"
#include "GAME4_2.h"
#include "GAME4_3.h"
#include "GAME5.h"
#include "GAME5_2.h"
#include "common__random.h"
#include "common__crypt.h"
#include "operators.h"
#include "server__dbase__objdb.h"
#include "server__script__script.h"
#include "common__system__team.h"

#include "client__gui__window.h"
#include "client__video__draw_common.h"
extern uint32_t dword_5d4594_2491716;
extern void* nox_alloc_hit_2491548;
extern uint32_t dword_5d4594_2491580;
extern uint32_t dword_587000_292488;
extern uint32_t dword_5d4594_2491676;
extern uint32_t dword_5d4594_2491588;
extern uint32_t dword_5d4594_2491592;
extern uint32_t dword_5d4594_2491704;
extern uint32_t dword_5d4594_2490508;
extern uint32_t dword_5d4594_2489460;
extern uint32_t dword_5d4594_2491544;
extern uint32_t dword_5d4594_2491552;
extern uint32_t dword_5d4594_2491616;
extern uint32_t dword_587000_292492;
extern uint32_t dword_5d4594_2650652;
extern uint32_t nox_player_netCode_85319C;

typedef struct {
	uint32_t f0;
	float f4;
	float f8;
	uint32_t f12;
	float f16;
	float f20;
} struct313272;

struct313272 table_313272[11] = {
	{f0:0x0, f4:0, f8:0, f12:0x1, f16:0, f20:23},
	{f0:0x1, f4:0, f8:23, f12:0x0, f16:0, f20:0},
	{f0:0x1, f4:0, f8:23, f12:0x1, f16:0, f20:23},
	{f0:0x1, f4:0, f8:11.5, f12:0x1, f16:0, f20:23},
	{f0:0x1, f4:0, f8:23, f12:0x1, f16:11.5, f20:23},
	{f0:0x1, f4:11.5, f8:23, f12:0x1, f16:0, f20:23},
	{f0:0x1, f4:0, f8:23, f12:0x1, f16:0, f20:11.5},
	{f0:0x1, f4:0, f8:11.5, f12:0x1, f16:11.5, f20:23},
	{f0:0x1, f4:11.5, f8:23, f12:0x1, f16:11.5, f20:23},
	{f0:0x1, f4:11.5, f8:23, f12:0x1, f16:0, f20:11.5},
	{f0:0x1, f4:0, f8:11.5, f12:0x1, f16:0, f20:11.5},
};

//----- (00547DB0) --------------------------------------------------------
int sub_547DB0(int a1, float2* a2) {
	double v3; // st7
	double v4; // st6
	float v5;  // [esp+0h] [ebp-10h]
	float v6;  // [esp+4h] [ebp-Ch]
	float v7;  // [esp+8h] [ebp-8h]
	float v8;  // [esp+Ch] [ebp-4h]

	if (*(uint32_t*)(a1 + 172) == 1) {
		return 0;
	}
	if (*(uint32_t*)(a1 + 172) != 2) {
		if (*(uint32_t*)(a1 + 172) == 3) {
			v5 = *(float*)(a1 + 192) + *(float*)(a1 + 64);
			v6 = *(float*)(a1 + 196) + *(float*)(a1 + 68);
			if ((v6 - v5 + a2->field_0 - a2->field_4) * 0.70709997 < 0.0 &&
				(*(float*)(a1 + 204) + *(float*)(a1 + 68) - (*(float*)(a1 + 200) + *(float*)(a1 + 64)) + a2->field_0 -
				 a2->field_4) *
						0.70709997 >
					0.0) {
				v7 = *(float*)(a1 + 208) + *(float*)(a1 + 64);
				v8 = *(float*)(a1 + 212) + *(float*)(a1 + 68);
				if ((v8 + v7 - a2->field_0 - a2->field_4) * 0.70709997 > 0.0 &&
					(v6 + v5 - a2->field_0 - a2->field_4) * 0.70709997 < 0.0) {
					return 1;
				}
			}
		}
		return 0;
	}
	v3 = *(float*)(a1 + 64) - a2->field_0;
	v4 = *(float*)(a1 + 68) - a2->field_4;
	if (v4 * v4 + v3 * v3 > *(float*)(a1 + 180)) {
		return 0;
	}
	return 1;
}

//----- (00548100) --------------------------------------------------------
void sub_548100(int2* a1, int a2) {
	int v2;    // eax
	int v3;    // ecx
	int v4;    // edx
	int v5;    // esi
	int v6;    // edx
	char* v7;  // eax
	int v8;    // eax
	float2 v9; // [esp+0h] [ebp-8h]
	int v10;   // [esp+Ch] [ebp+4h]

	v2 = nox_server_getWallAtGrid_410580(a1->field_0, a1->field_4);
	if (v2) {
		if (*(uint8_t*)(v2 + 4) & 4) {
			if (*(uint8_t*)(a2 + 8) & 6) {
				v3 = *(uint32_t*)(v2 + 28);
				if (*(uint8_t*)(v3 + 21) == 1) {
					if (*(uint8_t*)(v3 + 20) & 2) {
						*(uint8_t*)(v3 + 21) = 4;
						v4 = *(uint32_t*)(v3 + 4);
						*(uint8_t*)(v3 + 22) = 0;
						v5 = 23 * v4;
						v10 = 23 * *(uint32_t*)(v3 + 8);
						v6 = *(unsigned char*)(v2 + 1);
						v9.field_0 = (double)v5 + 11.5;
						v9.field_4 = (double)v10 + 11.5;
						v7 = nox_xxx_wallFindOpenSound_410EE0(v6);
						v8 = nox_xxx_utilFindSound_40AF50(v7);
						nox_xxx_audCreate_501A30(v8, &v9, 0, 0);
					}
				}
			}
		}
	}
}

//----- (005481C0) --------------------------------------------------------
void sub_5481C0(int a1) {
	*(uint32_t*)(a1 + 96) = 0;
	*(uint32_t*)(a1 + 100) = 0;
	if (!(*(uint8_t*)(a1 + 16) & 0x60)) {
		nox_xxx_getUnitsInRect_517C10((float4*)(a1 + 232), sub_548220, a1);
		if (!(*(uint8_t*)(a1 + 16) & 8)) {
			if (*(uint32_t*)(a1 + 172) == 2) {
				sub_54FEF0(a1);
			} else if (*(uint32_t*)(a1 + 172) == 3) {
				sub_5504B0(a1);
			}
		}
	}
}

//----- (00548220) --------------------------------------------------------
void sub_548220(int* a1, float* a2) {
	int v2; // eax
	int v3; // ecx
	int v4; // eax

	if (sub_548360((int)a2, (int)a1)) {
		v2 = *((uint32_t*)a2 + 2);
		if (v2 & 0x4000) {
			sub_551AE0((int)a2, (int)a1, 0);
		} else {
			v3 = a1[2];
			if (v3 & 0x4000) {
				sub_551AE0((int)a1, (int)a2, 1);
			} else if ((v2 & 0x8000) == 0) {
				if ((v3 & 0x8000) == 0) {
					if ((v2 & 0x80u) == 0) {
						if (*((uint32_t*)a2 + 43) == 2) {
							if ((v3 & 0x80u) == 0) {
								if (a1[43] == 2) {
									nox_xxx_collisionCheckCircleCircle_550D00((int)a2, (int)a1);
								} else if (a1[43] == 3) {
									sub_54AD50((int)a2, (int)a1, 0);
								}
							} else {
								sub_5488B0(a1, a2, 1);
							}
						} else if (*((uint32_t*)a2 + 43) == 3) {
							if ((v3 & 0x80u) == 0) {
								if (a1[43] == 2) {
									sub_54AD50((int)a1, (int)a2, 1);
								} else if (a1[43] == 3) {
									sub_550F80(a2, (int)a1);
								}
							} else {
								sub_551250((unsigned int)a1, a2, 1);
							}
						}
					} else {
						v4 = a1[43];
						if (v4 == 2) {
							sub_5488B0((int*)a2, (float*)a1, 0);
						} else if (v4 == 3) {
							sub_551250((unsigned int)a2, (float*)a1, 0);
						}
					}
				} else {
					sub_551C40((int)a1, (int)a2);
				}
			} else {
				sub_551C40((int)a2, (int)a1);
			}
		}
	}
}

//----- (00548360) --------------------------------------------------------
int sub_548360(int a1, int a2) {
	int v2;         // eax
	int v3;         // ebp
	int (*v4)(int); // ecx
	int (*v5)(int); // esi
	int v6;         // ebx
	int v7;         // eax
	int v8;         // edx
	short v9;       // cx
	int result;     // eax
	int v11;        // ecx
	int v12;        // [esp+0h] [ebp-Ch]
	int v13;        // [esp+4h] [ebp-8h]
	int v14;        // [esp+14h] [ebp+8h]

	if (dword_5d4594_2490508) {
		v2 = *getMemU32Ptr(0x5D4594, 2490516);
	} else {
		dword_5d4594_2490508 = nox_xxx_getNameId_4E3AA0("Trigger");
		*getMemU32Ptr(0x5D4594, 2490512) = nox_xxx_getNameId_4E3AA0("BlackPowder");
		v2 = nox_xxx_getNameId_4E3AA0("TelekinesisHand");
		*getMemU32Ptr(0x5D4594, 2490516) = v2;
	}
	v3 = a2;
	v4 = *(int (**)(int))(a1 + 696);
	if (v4 == nox_xxx_collidePentagram_4EAB20 && *(unsigned short*)(a2 + 4) == v2) {
		return 0;
	}
	v5 = *(int (**)(int))(a2 + 696);
	if (v5 == nox_xxx_collidePentagram_4EAB20 && *(unsigned short*)(a1 + 4) == v2) {
		return 0;
	}
	v6 = *(uint32_t*)(a1 + 8);
	v12 = *(uint32_t*)(a1 + 8) & 0x80;
	if (v12) {
		if (*(uint8_t*)(a2 + 16) & 9) {
			return 0;
		}
	}
	v14 = *(uint32_t*)(a2 + 8);
	if (v14 & 0x80) {
		if (*(uint8_t*)(a1 + 16) & 9) {
			return 0;
		}
	}
	v7 = *(uint32_t*)(a1 + 16);
	v13 = *(uint32_t*)(a1 + 16);
	if (v7 & 0x60) {
		return 0;
	}
	v8 = *(uint32_t*)(v3 + 16);
	if (v8 & 0x60 || a1 == v3 || !v4 || !v5) {
		return 0;
	}
	if (v6 & 2 && v7 & 0x4000) {
		v9 = v14;
		if (v14 & 0xCC00) {
			return sub_5485B0(a1, v3);
		}
	} else {
		v9 = v14;
	}
	if (v9 & 2 && v8 & 0x4000 && v6 & 0xCC00) {
		return sub_5485B0(v3, a1);
	}
	if (v6 & 0x2000 && *(unsigned short*)(v3 + 4) == *getMemU32Ptr(0x5D4594, 2490512)) {
		return 1;
	}
	if (v9 & 0x2000 && *(unsigned short*)(a1 + 4) == *getMemU32Ptr(0x5D4594, 2490512)) {
		return 1;
	}
	if ((v11 = *(uint32_t*)(a1 + 16), v13 & 8) && v8 & 8 || v12 && v8 & 8 || v14 & 0x80 && v13 & 8 ||
		*(unsigned short*)(a1 + 4) != dword_5d4594_2490508 && *(unsigned short*)(v3 + 4) != dword_5d4594_2490508 &&
			(v13 & 0x11 && v8 & 0x24000 || v8 & 0x11 && v11 & 0x24000) ||
		v12 && v14 & 0x80 || v6 & 0x4000 && v14 & 0x4000 || (v6 & 0x8000) != 0 && (v14 & 0x8000) != 0 ||
		v11 & 0x400 && nox_xxx_unitsHaveSameTeam_4EC520(a1, v3)) {
		return 0;
	} else {
		return 1;
	}
}

//----- (005485B0) --------------------------------------------------------
int sub_5485B0(int a1, int a2) {
	int v2;      // ecx
	int v3;      // eax
	int v4;      // edx
	uint32_t* i; // ecx

	v2 = *(uint32_t*)(a1 + 748);
	v3 = 0;
	v4 = *(unsigned char*)(v2 + 2172);
	if (v4 <= 0) {
		return 0;
	}
	for (i = (uint32_t*)(v2 + 2140); *i != *(uint32_t*)(a2 + 36); ++i) {
		if (++v3 >= v4) {
			return 0;
		}
	}
	return 1;
}

//----- (00548630) --------------------------------------------------------
void nox_xxx_collSysAddCollision_548630(int a1, unsigned int a2, float2* a3) {
	int v3;       // esi
	int v4;       // esi
	uint32_t* v5; // eax
	int v6;       // ecx
	int* v7;      // eax
	int v8;       // ecx

	v3 = *(uint32_t*)(a1 + 36);
	if (a2 > 6) {
		v3 += *(uint32_t*)(a2 + 36);
	}
	v4 = v3 % 256;
	v5 = *(uint32_t**)getMemAt(0x5D4594, 2490520 + 4 * v4);
	if (v5) {
		while (1) {
			v6 = v5[2];
			if (v6 == a1 && v5[3] == a2) {
				break;
			}
			if (v6 == a2 && v5[3] == a1) {
				break;
			}
			v5 = (uint32_t*)*v5;
			if (!v5) {
				goto LABEL_9;
			}
		}
		return;
	}
LABEL_9:
	v7 = (int*)nox_alloc_class_new_obj_zero(*(uint32_t**)&nox_alloc_hit_2491548);
	if (v7) {
		v7[2] = a1;
		v7[3] = a2;
		*((float2*)v7 + 2) = *a3;
		v7[6] = v4;
		*v7 = *getMemU32Ptr(0x5D4594, 2490520 + 4 * v4);
		v8 = dword_5d4594_2491544;
		*getMemU32Ptr(0x5D4594, 2490520 + 4 * v4) = v7;
		v7[1] = v8;
		dword_5d4594_2491544 = v7;
	}
}

//----- (005486D0) --------------------------------------------------------
void nox_xxx_allocHitArray_5486D0() {
	char* v0; // edx
	int i;    // eax

	v0 = *(char**)&nox_alloc_hit_2491548;
	if (!nox_alloc_hit_2491548) {
		v0 = nox_new_alloc_class("Hit", 28, 1024);
		nox_alloc_hit_2491548 = v0;
		memset(getMemAt(0x5D4594, 2490520), 0, 0x400u);
	}
	for (i = dword_5d4594_2491544; i; i = *(uint32_t*)(i + 4)) {
		*getMemU32Ptr(0x5D4594, 2490520 + 4 * *(uint32_t*)(i + 24)) = 0;
	}
	nox_alloc_class_free_all(v0);
	dword_5d4594_2491544 = 0;
}

//----- (00548740) --------------------------------------------------------
void nox_xxx_collide_548740() {
	int i;           // esi
	unsigned int v1; // ecx
	int v2;          // eax
	float2 v3;       // [esp+4h] [ebp-8h]

	for (i = dword_5d4594_2491544; i; i = *(uint32_t*)(i + 4)) {
		v1 = *(uint32_t*)(i + 12);
		if (v1 > 6 || !v1) {
			(*(void (**)(uint32_t, unsigned int, int))(*(uint32_t*)(i + 8) + 696))(*(uint32_t*)(i + 8), v1, i + 16);
			if (*(uint32_t*)(i + 12)) {
				nox_xxx_collide_4FDF90(*(uint32_t*)(i + 8), *(uint32_t*)(i + 12));
			}
		}
		v2 = *(uint32_t*)(i + 12);
		if (v2 == 6) {
			(*(void (**)(uint32_t, uint32_t, uint32_t, int, int))(*(uint32_t*)(i + 8) + 716))(*(uint32_t*)(i + 8), 0, 0,
																							  2, 12);
			nox_xxx_unitHasCollideOrUpdateFn_537610(*(uint32_t*)(i + 8));
		} else if (v2) {
			v3.field_0 = -*(float*)(i + 16);
			v3.field_4 = -*(float*)(i + 20);
			(*(void (**)(uint32_t, uint32_t, float2*))(*(uint32_t*)(i + 12) + 696))(*(uint32_t*)(i + 12),
																					*(uint32_t*)(i + 8), &v3);
			nox_xxx_collide_4FDF90(*(uint32_t*)(i + 12), *(uint32_t*)(i + 8));
			if (*(uint8_t*)(*(uint32_t*)(i + 8) + 16) & 8) {
				nox_xxx_unitHasCollideOrUpdateFn_537610(*(uint32_t*)(i + 12));
			} else if (*(uint8_t*)(*(uint32_t*)(i + 12) + 16) & 8) {
				nox_xxx_unitHasCollideOrUpdateFn_537610(*(uint32_t*)(i + 8));
			}
		}
		nullsub_30(*(uint32_t*)(i + 8));
	}
}
// 5485F0: using guessed type void  nullsub_30(uint32_t);

//----- (00548830) --------------------------------------------------------
void sub_548830(int a1) {
	if (!*(uint32_t*)(a1 + 28)) {
		*(uint32_t*)(a1 + 36) = dword_5d4594_2491552;
		dword_5d4594_2491552 = a1;
		*(uint32_t*)(a1 + 28) = 1;
	}
}

//----- (00548860) --------------------------------------------------------
void sub_548860(int a1, short a2) {
	int v2; // eax

	v2 = *(uint32_t*)(a1 + 748);
	for (*(uint16_t*)(v2 + 40) += a2; *(uint16_t*)(v2 + 40) < 0; *(uint16_t*)(v2 + 40) += 256) {
		;
	}
	for (; *(uint16_t*)(v2 + 40) >= 256; *(uint16_t*)(v2 + 40) -= 256) {
		;
	}
	sub_548830(v2);
}

//----- (005488B0) --------------------------------------------------------
void sub_57C790(float4* a1, float2* a2, float2* a3, float a4);
void sub_5488B0(int* a1, float* a2, int a3) {
	int* v3;                     // edi
	int v4;                      // ebx
	float v5;                    // ecx
	int v6;                      // eax
	float* v7;                   // esi
	float* v8;                   // ebp
	double v9;                   // st7
	double v10;                  // st6
	long double v11;             // st6
	double v12;                  // st7
	float v13;                   // et1
	long double v14;             // st7
	float v15;                   // et1
	long double v16;             // st7
	int v17;                     // eax
	float* v18;                  // eax
	double v19;                  // st7
	double v20;                  // st7
	float v21;                   // [esp+0h] [ebp-44h]
	float v22;                   // [esp+4h] [ebp-40h]
	int v23;                     // [esp+4h] [ebp-40h]
	float v24;                   // [esp+18h] [ebp-2Ch]
	int* v25;                    // [esp+1Ch] [ebp-28h]
	volatile unsigned char* v26; // [esp+20h] [ebp-24h]
	float2 v27;                  // [esp+24h] [ebp-20h]
	float2 a3a;                  // [esp+2Ch] [ebp-18h]
	float4 a1a;                  // [esp+34h] [ebp-10h]
	signed int v30;              // [esp+48h] [ebp+4h]
	float v31;                   // [esp+48h] [ebp+4h]
	float v32;                   // [esp+4Ch] [ebp+8h]
	float v33;                   // [esp+4Ch] [ebp+8h]

	v3 = a1;
	v4 = a1[187];
	v5 = *((float*)a1 + 17);
	LODWORD(a1a.field_0) = a1[16];
	a1a.field_4 = v5;
	v6 = *(short*)(v4 + 40) + 160;
	if (v6 >= 256) {
		v6 = *(short*)(v4 + 40) - 96;
	}
	v26 = getMemAt(0x587000, 192092 + 8 * v6);
	v7 = a2;
	v25 = getMemIntPtr(0x587000, 192088 + 8 * v6);
	v30 = 2 * *(int*)v26;
	v8 = a2 + 16;
	a1a.field_8 = (double)(2 * *getMemIntPtr(0x587000, 192088 + 8 * v6)) + a1a.field_0;
	a1a.field_C = (double)v30 + a1a.field_4;
	sub_57C790(&a1a, (float2*)a2 + 8, &a3a, 32.0);
	v9 = a2[16] - a3a.field_0;
	v10 = a2[17] - a3a.field_4;
	v27.field_4 = v10;
	v11 = sqrt(v10 * v27.field_4 + v9 * v9);
	v31 = v11;
	if (v11 == 0.0) {
		v31 = 0.1;
	}
	if (v31 < (double)a2[44]) {
		v27.field_0 = v9 / v31;
		v27.field_4 = v27.field_4 / v31;
		nox_xxx_collSysAddCollision_548630((int)a2, (unsigned int)v3, &v27);
		*(uint32_t*)(v4 + 44) = gameFrame();
		if (a3 == 1) {
			v32 = a2[44] - v31;
			v24 = -(v27.field_4 * v7[21]) - v27.field_0 * v7[20];
			v12 = nox_xxx_objectGetMass_4E4A70((int)v7);
			v13 = *(float*)&dword_587000_292492;
			v14 = sqrt(v12 * v13 * 4.0);
			v15 = *(float*)&dword_587000_292492;
			v16 = v14 * v24 * 0.25 + v15 * v32;
			v22 = v16 * v27.field_4;
			v21 = v16 * v27.field_0;
			sub_548600((int)v7, v21, v22);
		}
		v17 = *((uint32_t*)v7 + 4);
		if (v17 & 0x8000000) {
			if (!(v17 & 8)) {
				nox_xxx_unitHasCollideOrUpdateFn_537610((int)v7);
			}
			*((uint32_t*)v7 + 4) &= 0xF7FFFFFF;
		}
		nox_xxx_unitHasCollideOrUpdateFn_537610((int)v3);
		if (!nox_xxx_servObjectHasTeam_419130((int)(v3 + 12)) || *(uint32_t*)(v4 + 12) != *(uint32_t*)(v4 + 4) ||
			nox_xxx_servCompareTeams_419150((int)(v3 + 12), (int)(v7 + 12))) {
			if (!a3 && !*(uint8_t*)(v4 + 1)) {
				v18 = (float*)v3[127];
				if (!v18 || v18 == v7) {
					v33 = v7[44] - v31;
					v19 = v33 * v7[30];
					if ((double)*v25 * (v7[17] - *((float*)v3 + 17)) -
							(double)*(int*)v26 * (*v8 - *((float*)v3 + 16)) <=
						0.0) {
						v20 = v19 + *(float*)(v4 + 32);
					} else {
						v20 = *(float*)(v4 + 32) - v19;
					}
					*(float*)(v4 + 32) = v20;
					sub_548830(v4);
					nox_xxx_unitAddToUpdatable_4DA8D0((int)v3);
				}
			}
		} else if (gameFrame() > (unsigned int)v3[34]) {
			v23 = *((unsigned char*)v3 + 52);
			v3[34] = gameFrame() + gameFPS();
			nox_xxx_getTeamByID_418AB0(v23);
			nox_xxx_netPriMsgToPlayer_4DA2C0((int)v7, "objcoll.c:GateLockedMechanism", 0);
		}
	}
}

//----- (00548B60) --------------------------------------------------------
void sub_548B60() {
	int v0;   // esi
	int v1;   // eax
	int v2;   // ecx
	short v3; // ax
	short v4; // ax
	short v5; // ax
	int v6;   // eax
	int v7;   // eax
	int v8;   // eax
	float v9; // [esp+0h] [ebp-14h]

	v0 = dword_5d4594_2491552;
	if (dword_5d4594_2491552) {
		while (1) {
			v9 = *(float*)(v0 + 32) + *(float*)(v0 + 32) + 0.5;
			v1 = nox_float2int(v9);
			v2 = v1;
			if (v1 < 0) {
				v2 = -v1;
			}
			if (v2 > 4) {
				LOWORD(v2) = 4;
			}
			if (*(float*)(v0 + 32) >= -0.0099999998) {
				if (*(float*)(v0 + 32) <= 0.0099999998) {
					goto LABEL_15;
				}
				if (*(uint32_t*)(v0 + 12) == (*(uint32_t*)(v0 + 4) + 12) % 32) {
					goto LABEL_15;
				}
				*(uint16_t*)(v0 + 40) += v2;
				v5 = *(uint16_t*)(v0 + 40);
				if (v5 < 256) {
					goto LABEL_15;
				}
				v4 = v5 - 256;
			} else {
				if (*(int*)(v0 + 12) == *(int*)(v0 + 4) - 12 + (*(int*)(v0 + 4) - 12 < 0 ? 0x20 : 0)) {
					goto LABEL_15;
				}
				*(uint16_t*)(v0 + 40) -= v2;
				v3 = *(uint16_t*)(v0 + 40);
				if (v3 >= 0) {
					goto LABEL_15;
				}
				v4 = v3 + 256;
			}
			*(uint16_t*)(v0 + 40) = v4;
		LABEL_15:
			v6 = 32 * *(short*)(v0 + 40) / 256;
			*(uint32_t*)(v0 + 12) = v6;
			if (v6 < 0) {
				do {
					v7 = *(uint32_t*)(v0 + 12);
					*(uint32_t*)(v0 + 12) = v7 + 32;
				} while (v7 + 32 < 0);
			}
			if (*(uint32_t*)(v0 + 12) >= 32) {
				do {
					v8 = *(uint32_t*)(v0 + 12) - 32;
					*(uint32_t*)(v0 + 12) = v8;
				} while (v8 >= 32);
			}
			*(uint32_t*)(v0 + 28) = 0;
			*(uint32_t*)(v0 + 32) = 0;
			v0 = *(uint32_t*)(v0 + 36);
			if (!v0) {
				dword_5d4594_2491552 = 0;
				return;
			}
		}
	}
	dword_5d4594_2491552 = 0;
}

//----- (00548CD0) --------------------------------------------------------
void nox_xxx_script_forcedialog_548CD0(nox_object_t* a1p, nox_object_t* a2p) {
	int a1 = a1p;
	int a2 = a2p;
	int v2; // ecx
	int v3; // eax

	if (a1) {
		if (a2) {
			if (*(uint8_t*)(a1 + 8) & 4) {
				if (*(uint8_t*)(a2 + 8) & 2) {
					if (!(*(uint32_t*)(a1 + 16) & 0x8020)) {
						v2 = *(uint32_t*)(a2 + 748);
						if (*(uint8_t*)(a2 + 20) & 0x10) {
							v3 = *(uint32_t*)(v2 + 2096);
							if (v3 != -1 && *(int*)(v2 + 2100) != -1) {
								nox_script_callByIndex_507310(v3, a1, a2);
							}
						}
					}
				}
			}
		}
	}
}

//----- (00548D30) --------------------------------------------------------
void nox_xxx_scriptDialog_548D30(nox_object_t* a1p, char a2) {
	int a1 = a1p;
	int v2; // ebx
	int v3; // edi
	int v5; // ebp
	int v6; // esi

	v2 = a1;
	v3 = *(uint32_t*)(a1 + 748);
	nox_xxx_unitUnFreeze_4E7A60(a1, 0);
	v5 = *(uint32_t*)(v3 + 284);
	if (v5) {
		v6 = *(uint32_t*)(v5 + 748);
		if (*(int*)(v6 + 2096) != -1 && *(int*)(v6 + 2100) != -1) {
			LOWORD(a1) = 1232;
			nox_xxx_netSendPacket0_4E5420(*(unsigned char*)(*(uint32_t*)(v3 + 276) + 2064), &a1, 2, 0, 1);
			*(uint32_t*)(v3 + 284) = 0;
			if (*(uint8_t*)(v6 + 2104) == 1) {
				*(uint8_t*)(v6 + 2105) = a2;
			} else {
				*(uint8_t*)(v6 + 2105) = 0;
			}
			nox_script_callByIndex_507310(*(uint32_t*)(v6 + 2100), v2, v5);
		}
	}
}

//----- (0054A990) --------------------------------------------------------
double sub_54A990(float2* a1, float a2, int a3, float2* a4) {
	int v4;          // ecx
	float2* v5;      // edx
	double result;   // st7
	char v7;         // al
	int v8;          // eax
	double v9;       // st7
	long double v10; // st7
	double v11;      // st7
	long double v12; // st7
	double v13;      // st7
	long double v14; // st7
	char v15;        // [esp+0h] [ebp-30h]
	float v16;       // [esp+0h] [ebp-30h]
	float v17;       // [esp+4h] [ebp-2Ch]
	float v18;       // [esp+8h] [ebp-28h]
	float v19;       // [esp+8h] [ebp-28h]
	float v20;       // [esp+Ch] [ebp-24h]
	float v21;       // [esp+10h] [ebp-20h]
	float v22;       // [esp+14h] [ebp-1Ch]
	float v23;       // [esp+18h] [ebp-18h]
	float v24;       // [esp+1Ch] [ebp-14h]
	float v25;       // [esp+20h] [ebp-10h]
	float v26;       // [esp+24h] [ebp-Ch]
	float v27;       // [esp+28h] [ebp-8h]
	float v28;       // [esp+2Ch] [ebp-4h]
	float v29;       // [esp+34h] [ebp+4h]
	float v30;       // [esp+34h] [ebp+4h]
	float v31;       // [esp+3Ch] [ebp+Ch]
	float v32;       // [esp+3Ch] [ebp+Ch]

	v4 = a3;
	v5 = a1;
	result = -1.0;
	v15 = 0;
	v21 = *(float*)(a3 + 192) + *(float*)(a3 + 64);
	v22 = *(float*)(a3 + 196) + *(float*)(a3 + 68);
	v25 = *(float*)(a3 + 200) + *(float*)(a3 + 64);
	v26 = *(float*)(a3 + 204) + *(float*)(a3 + 68);
	v23 = *(float*)(a3 + 208) + *(float*)(a3 + 64);
	v24 = *(float*)(a3 + 212) + *(float*)(a3 + 68);
	v27 = *(float*)(a3 + 216) + *(float*)(a3 + 64);
	v28 = *(float*)(a3 + 220) + *(float*)(a3 + 68);
	v31 = (v22 - v21 + a1->field_0 - a1->field_4) * 0.70709997;
	v29 = (v26 - v25 + a1->field_0 - a1->field_4) * 0.70709997;
	v17 = (v24 + v23 - v5->field_0 - v5->field_4) * 0.70709997;
	v18 = (v22 + v21 - v5->field_0 - v5->field_4) * 0.70709997;
	if (v31 <= 0.0) {
		if (v29 < 0.0) {
			v15 = 2;
		}
	} else {
		v15 = 1;
	}
	if (v17 >= 0.0) {
		if (v18 <= 0.0) {
			goto LABEL_10;
		}
		v7 = v15 | 4;
	} else {
		v7 = v15 | 8;
	}
	v15 = v7;
LABEL_10:
	switch (v15) {
	case 0:
		v19 = v5->field_0 - *(float*)(v4 + 64);
		v20 = v5->field_4 - *(float*)(v4 + 68);
		if (v19 == 0.0 && v20 == 0.0) {
			v8 = nox_common_randomInt_415FA0(0, 3);
			v19 = *getMemFloatPtr(0x587000, 289928 + 8 * v8);
			v20 = *getMemFloatPtr(0x587000, 289932 + 8 * v8);
		}
		result = sqrt(v20 * v20 + v19 * v19);
		if (result == 0.0) {
			result = 0.1;
		}
		a4->field_0 = v19 / result;
		a4->field_4 = v20 / result;
		return result;
	case 1:
		result = a2 - v31;
		a4->field_0 = 0.70709997;
		a4->field_4 = -0.70709997;
		return result;
	case 2:
		result = a2 - -v29;
		a4->field_0 = -0.70709997;
		a4->field_4 = 0.70709997;
		return result;
	case 4:
		result = a2 - v18;
		a4->field_0 = -0.70709997;
		a4->field_4 = -0.70709997;
		return result;
	case 5:
		v30 = v5->field_0 - v21;
		v9 = v5->field_4 - v22;
		v16 = v9;
		v10 = sqrt(v9 * v16 + v30 * v30);
		v32 = v10;
		if (v10 == 0.0) {
			v32 = 0.1;
		}
		result = a2 - v32;
		if (result >= 0.0) {
			a4->field_0 = v30 / v32;
			a4->field_4 = v16 / v32;
		}
		return result;
	case 6:
		v30 = v5->field_0 - v25;
		v13 = v5->field_4 - v26;
		v16 = v13;
		v14 = sqrt(v13 * v16 + v30 * v30);
		v32 = v14;
		if (v14 == 0.0) {
			v32 = 0.1;
		}
		result = a2 - v32;
		if (result >= 0.0) {
			a4->field_0 = v30 / v32;
			a4->field_4 = v16 / v32;
		}
		return result;
	case 8:
		result = a2 - -v17;
		a4->field_0 = 0.70709997;
		a4->field_4 = 0.70709997;
		return result;
	case 9:
		v30 = v5->field_0 - v23;
		v9 = v5->field_4 - v24;
		v16 = v9;
		v10 = sqrt(v9 * v16 + v30 * v30);
		v32 = v10;
		if (v10 == 0.0) {
			v32 = 0.1;
		}
		result = a2 - v32;
		if (result >= 0.0) {
			a4->field_0 = v30 / v32;
			a4->field_4 = v16 / v32;
		}
		return result;
	case 0xA:
		v30 = v5->field_0 - v27;
		v11 = v5->field_4 - v28;
		v16 = v11;
		v12 = sqrt(v11 * v16 + v30 * v30);
		v32 = v12;
		if (v12 == 0.0) {
			v32 = 0.1;
		}
		result = a2 - v32;
		if (result >= 0.0) {
			a4->field_0 = v30 / v32;
			a4->field_4 = v16 / v32;
		}
		return result;
	default:
		return result;
	}
}

//----- (0054AD50) --------------------------------------------------------
void sub_54AD50(int a1, int a2, int a3) {
	int v3;     // esi
	int v4;     // edi
	int v5;     // ebp
	double v6;  // st7
	float v7;   // eax
	float v8;   // edx
	float v9;   // eax
	int v10;    // ebx
	int v11;    // eax
	double v12; // st7
	double v13; // st7
	double v14; // st6
	double v15; // st7
	double v16; // st7
	float v17;  // [esp+0h] [ebp-34h]
	float v18;  // [esp+4h] [ebp-30h]
	float v19;  // [esp+18h] [ebp-1Ch]
	float2 a4;  // [esp+1Ch] [ebp-18h]
	float4 a1a; // [esp+24h] [ebp-10h]
	int v22;    // [esp+38h] [ebp+4h]
	int v23;    // [esp+38h] [ebp+4h]
	float a3c;  // [esp+3Ch] [ebp+8h]
	int a3a;    // [esp+3Ch] [ebp+8h]
	int a3b;    // [esp+3Ch] [ebp+8h]
	float v27;  // [esp+40h] [ebp+Ch]

	v3 = a1;
	v4 = a2;
	v5 = 1;
	v6 = sub_54A990((float2*)(a1 + 64), *(float*)(a1 + 176), a2, &a4);
	if (v6 >= 0.0) {
		if (!(*(uint32_t*)(a1 + 8) & 0x2204) || !(*(uint32_t*)(a2 + 8) & 0x2204) ||
			(v7 = *(float*)(a1 + 56), v8 = *(float*)(a2 + 56), a1a.field_4 = *(float*)(a1 + 60), a1a.field_0 = v7,
			 v9 = *(float*)(a2 + 60), a1a.field_8 = v8, a1a.field_C = v9, nox_xxx_mapTraceRay_535250(&a1a, 0, 0, 0))) {
			nox_xxx_collSysAddCollision_548630(a1, a2, &a4);
			if ((*(uint8_t*)(a1 + 16) & 8) == 8 || (*(uint8_t*)(a2 + 16) & 8) == 8) {
				v5 = 0;
			}
			v10 = a3;
			if (a3 || !(*(uint8_t*)(a1 + 8) & 6) || (v11 = *(uint32_t*)(a2 + 16), !(v11 & 0x2000))) {
				if (v5) {
					a3c = v6;
					v12 = *(float*)&dword_587000_292488 * a3c;
					a1a.field_4 = a4.field_0;
					*(float*)&v22 = a4.field_0 * v12;
					*(float*)&a3a = a4.field_4 * v12;
					v13 = *(float*)(v3 + 80) - *(float*)(v4 + 80);
					v14 = *(float*)(v3 + 84) - *(float*)(v4 + 84);
					a1a.field_0 = -a4.field_4;
					v19 = a1a.field_0 * v13 + v14 * a4.field_0;
					v27 = nox_xxx_objectGetMass_4E4A70(v3);
					if (nox_xxx_objectGetMass_4E4A70(v4) <= v27) {
						v15 = nox_xxx_objectGetMass_4E4A70(v4);
					} else {
						v15 = nox_xxx_objectGetMass_4E4A70(v3);
					}
					v16 = v15 * v19;
					*(float*)&v23 = *(float*)&v22 - v16 * a1a.field_0 * 0.69999999;
					*(float*)&a3b = *(float*)&a3a - v16 * a1a.field_4 * 0.69999999;
					if (v10) {
						v18 = -*(float*)&a3b;
						v17 = -*(float*)&v23;
						sub_548600(v4, v17, v18);
					} else {
						sub_548600(v3, *(float*)&v23, *(float*)&a3b);
					}
				}
			}
			if (*(uint32_t*)(v3 + 16) & 0x8000000) {
				nox_xxx_unitHasCollideOrUpdateFn_537610(v3);
				*(uint32_t*)(v3 + 16) &= 0xF7FFFFFF;
			}
			if (*(uint32_t*)(v4 + 16) & 0x8000000) {
				nox_xxx_unitHasCollideOrUpdateFn_537610(v4);
				*(uint32_t*)(v4 + 16) &= 0xF7FFFFFF;
			}
		}
	}
}

//----- (0054AF40) --------------------------------------------------------
nox_object_t* nox_xxx_findObjectAtCursor_54AF40(nox_object_t* a1p) {
	int a1 = a1p;
	int v1;     // eax
	double v2;  // st7
	float2 a1a; // [esp+0h] [ebp-8h]

	v1 = *(uint32_t*)(*(uint32_t*)(a1 + 748) + 276);
	a1a.field_0 = (double)*(int*)(v1 + 2284);
	v2 = (double)*(int*)(v1 + 2288);
	dword_5d4594_2491592 = a1;
	a1a.field_4 = v2;
	*getMemU32Ptr(0x5D4594, 2491596) = 0;
	*getMemU32Ptr(0x5D4594, 2491600) = 0;
	nox_xxx_unitsGetInCircle_517F90(&a1a, 100.0, (int)nox_xxx_playerCursorScanFn_54AFB0, (int)&a1a);
	return *getMemU32Ptr(0x5D4594, 2491596);
}

//----- (0054AFB0) --------------------------------------------------------
void nox_xxx_playerCursorScanFn_54AFB0(int a1, float* a2) {
	float* v2;  // esi
	char* v3;   // eax
	int v4;     // eax
	double v5;  // st7
	float v6;   // ecx
	int v7;     // eax
	double v8;  // st7
	double v9;  // st6
	double v10; // st7
	double v11; // st6
	double v12; // st7
	float v13;  // [esp+0h] [ebp-28h]
	float v14;  // [esp+10h] [ebp-18h]
	float2 a3;  // [esp+18h] [ebp-10h]
	float2 a1a; // [esp+20h] [ebp-8h]
	float v17;  // [esp+2Ch] [ebp+4h]

	if (!*getMemU32Ptr(0x5D4594, 2491604)) {
		*getMemU32Ptr(0x5D4594, 2491604) = nox_xxx_getNameId_4E3AA0("Polyp");
	}
	v2 = (float*)a1;
	if (a1 != dword_5d4594_2491592 && !(*(uint32_t*)(a1 + 16) & 0x8020) &&
		(!nox_xxx_testUnitBuffs_4FF350(a1, 0) || nox_xxx_testUnitBuffs_4FF350(*(int*)&dword_5d4594_2491592, 21)) &&
		(*(uint32_t*)(a1 + 8) & 0x80000206 || *(unsigned short*)(a1 + 4) == *getMemU32Ptr(0x5D4594, 2491604))) {
		if (nox_xxx_mapCheck_537110(a1, *(int*)&dword_5d4594_2491592)) {
			if (!(*(uint8_t*)(a1 + 8) & 4) ||
				(*(uint32_t*)(a1 + 36) != nox_player_netCode_85319C ||
				 !nox_common_getEngineFlag(NOX_ENGINE_FLAG_DISABLE_GRAPHICS_RENDERING)) &&
					(v3 = nox_common_playerInfoGetByID_417040(*(uint32_t*)(a1 + 36))) != 0 && !(v3[3680] & 1)) {
				if ((*(uint32_t*)(a1 + 8) & 0x200) != 512 || (v4 = *(uint32_t*)(a1 + 16), BYTE1(v4) & 0x40)) {
					v5 = *(float*)(a1 + 60) - *(float*)(a1 + 104);
					v14 = v5;
					v17 = v5;
					if (*((uint32_t*)v2 + 43) == 2) {
						v13 = v2[44] * v2[44];
						v7 = nox_float2int(v13);
						a3.field_0 = v2[14] - v2[44];
						a1a.field_0 = v2[44] + v2[14];
						if (a2[1] > (double)v17) {
							v8 = *a2 - v2[14];
							v9 = a2[1] - v17;
							if ((double)v7 <= v9 * v9 + v8 * v8) {
								return;
							}
							v12 = v2[26] + v2[15];
							if (v12 > *getMemFloatPtr(0x5D4594, 2491600)) {
								*getMemFloatPtr(0x5D4594, 2491600) = v12;
								*getMemU32Ptr(0x5D4594, 2491596) = v2;
							}
							return;
						}
						v10 = *a2;
						if (a2[1] < (double)v14) {
							v11 = a2[1] - v14;
							if ((double)v7 <= v11 * v11 + (v10 - v2[14]) * (v10 - v2[14])) {
								return;
							}
							v12 = v2[26] + v2[15];
							if (v12 > *getMemFloatPtr(0x5D4594, 2491600)) {
								*getMemFloatPtr(0x5D4594, 2491600) = v12;
								*getMemU32Ptr(0x5D4594, 2491596) = v2;
							}
							return;
						}
						if (v10 > a3.field_0 && *a2 < (double)a1a.field_0) {
							v12 = v2[26] + v2[15];
							if (v12 > *getMemFloatPtr(0x5D4594, 2491600)) {
								*getMemFloatPtr(0x5D4594, 2491600) = v12;
								*getMemU32Ptr(0x5D4594, 2491596) = v2;
							}
							return;
						}
					} else if (*((uint32_t*)v2 + 43) == 3) {
						v6 = v2[14];
						a1a.field_4 = v5;
						a3.field_0 = *a2;
						a1a.field_0 = v6;
						a3.field_4 = a2[1];
						if (nox_xxx_map_57B850(&a1a, v2 + 43, &a3)) {
							v12 = v2[26] + v2[15];
							if (v12 > *getMemFloatPtr(0x5D4594, 2491600)) {
								*getMemFloatPtr(0x5D4594, 2491600) = v12;
								*getMemU32Ptr(0x5D4594, 2491596) = v2;
							}
							return;
						}
						a1a.field_4 = v5;
						if (nox_xxx_map_57B850(&a1a, v2 + 43, &a3) || v2[50] + v2[14] < *a2 && *a2 < (double)v2[14] &&
																		  v14 + v2[51] < a2[1] &&
																		  v17 + v2[51] > a2[1]) {
							v12 = v2[26] + v2[15];
							if (v12 > *getMemFloatPtr(0x5D4594, 2491600)) {
								*getMemFloatPtr(0x5D4594, 2491600) = v12;
								*getMemU32Ptr(0x5D4594, 2491596) = v2;
							}
							return;
						}
						if (*a2 >= (double)v2[14] && v2[52] + v2[14] > *a2 && v14 + v2[53] < a2[1] &&
							v17 + v2[53] > a2[1]) {
							v12 = v2[26] + v2[15];
							if (v12 > *getMemFloatPtr(0x5D4594, 2491600)) {
								*getMemFloatPtr(0x5D4594, 2491600) = v12;
								*getMemU32Ptr(0x5D4594, 2491596) = v2;
							}
							return;
						}
					}
				}
			}
		}
	}
}

//----- (0054B2D0) --------------------------------------------------------
int sub_54B2D0(int* a1, int a2, uint32_t* a3) {
	int v3;   // ebp
	int* v4;  // edi
	int v5;   // edx
	int v6;   // ebx
	int v7;   // eax
	int v8;   // kr04_4
	int* v9;  // ebx
	int* v11; // edi
	int v12;  // ebx
	int v13;  // eax
	int v14;  // kr0C_4
	int* v15; // ebx
	int* v16; // edi
	int v17;  // ebx
	int v18;  // ebx
	int v19;  // eax
	int v20;  // kr14_4
	int* v21; // ebx
	int v22;  // ecx
	int* v23; // edi
	int v24;  // ebx
	int v25;  // eax
	int v26;  // ecx
	int v27;  // ebx
	int v28;  // eax
	int v29;  // kr1C_4
	int* v30; // ebx
	int v31;  // ecx
	int v32;  // [esp-4h] [ebp-1Ch]
	int v33;  // [esp-4h] [ebp-1Ch]
	int2 a4;  // [esp+10h] [ebp-8h]
	int v35;  // [esp+20h] [ebp+8h]
	int v36;  // [esp+20h] [ebp+8h]
	int v37;  // [esp+20h] [ebp+8h]
	int v38;  // [esp+20h] [ebp+8h]
	int a2a;  // [esp+24h] [ebp+Ch]
	int a2b;  // [esp+24h] [ebp+Ch]
	int a2c;  // [esp+24h] [ebp+Ch]
	int a2d;  // [esp+24h] [ebp+Ch]
	int a2e;  // [esp+24h] [ebp+Ch]
	int a2f;  // [esp+24h] [ebp+Ch]
	int a2g;  // [esp+24h] [ebp+Ch]
	int a2h;  // [esp+24h] [ebp+Ch]

	v3 = a1[4 * a2 + 22];
	switch (a2) {
	case 0:
		v4 = a3;
		a2a = 0;
		v5 = v4[4];
		v6 = v4[3] - v3;
		a4.field_0 = v4[1] + v4[3] / 2;
		a4.field_4 = v5 + v4[2] - 1;
		if (v6 >= 0) {
			do {
				if (sub_54B810(a1[37], (int)v4, a1 + 20, &a4, a1[22])) {
					return 1;
				}
				v7 = v4[1];
				if (++a4.field_0 > v7 + v6) {
					a4.field_0 = v7;
				}
				++a2a;
			} while (a2a <= v6);
		}
		a2b = 0;
		v8 = v4[4];
		v35 = v8 - v3;
		a4.field_4 = v4[2] + v8 / 2;
		if (v8 - v3 < 0) {
			return 0;
		}
		v9 = a1 + 22;
		do {
			a4.field_0 = v4[1];
			if (sub_54BD90(a1[37], (int)v4, a1 + 20, &a4.field_0, *v9)) {
				return 1;
			}
			v32 = *v9;
			a4.field_0 = a4.field_0 + v4[3] - 1;
			if (sub_54BD90(a1[37], (int)v4, a1 + 20, &a4.field_0, v32)) {
				return 1;
			}
			if (++a4.field_4 > v4[2] + v35) {
				a4.field_4 = v4[2];
			}
			++a2b;
		} while (a2b <= v35);
		return 0;
	case 1:
		v11 = a3;
		a2c = 0;
		v12 = v11[3] - v3;
		a4.field_0 = v11[1] + v11[3] / 2;
		a4.field_4 = v11[2];
		if (v12 < 0) {
			goto LABEL_20;
		}
		while (1) {
			if (sub_54B810((int)v11, a1[37], &a4.field_0, (int2*)a1 + 12, a1[26])) {
				return 1;
			}
			v13 = v11[1];
			if (++a4.field_0 > v13 + v12) {
				a4.field_0 = v13;
			}
			if (++a2c > v12) {
				break;
			}
		}
	LABEL_20:
		a2d = 0;
		v14 = v11[4];
		v36 = v14 - v3;
		a4.field_4 = v11[2] + v14 / 2;
		if (v14 - v3 < 0) {
			return 0;
		}
		v15 = a1 + 26;
		while (1) {
			a4.field_0 = v11[1];
			if (sub_54BF20(a1[37], (int)v11, a1 + 24, &a4.field_0, *v15)) {
				break;
			}
			v33 = *v15;
			a4.field_0 = a4.field_0 + v11[3] - 1;
			if (sub_54BF20(a1[37], (int)v11, a1 + 24, &a4.field_0, v33)) {
				break;
			}
			if (++a4.field_4 > v11[2] + v36) {
				a4.field_4 = v11[2];
			}
			if (++a2d > v36) {
				return 0;
			}
		}
		return 1;
	case 2:
		v16 = a3;
		a2e = 0;
		v17 = v16[3];
		a4.field_0 = v16[1];
		v18 = v17 - v3;
		a4.field_4 = v16[2] + v16[4] / 2;
		if (v18 < 0) {
			goto LABEL_33;
		}
		while (1) {
			if (sub_54BB20(a1[37], (int)v16, a1 + 28, &a4, a1[30])) {
				return 1;
			}
			v19 = v16[2];
			if (++a4.field_4 > v19 + v18) {
				a4.field_4 = v19;
			}
			if (++a2e > v18) {
				break;
			}
		}
	LABEL_33:
		a2f = 0;
		v20 = v16[3];
		v37 = v20 - v3;
		a4.field_0 = v16[1] + v20 / 2;
		if (v20 - v3 < 0) {
			return 0;
		}
		v21 = a1 + 30;
		while (1) {
			a4.field_4 = v16[2];
			if (sub_54BD90((int)v16, a1[37], &a4.field_0, a1 + 28, *v21)) {
				break;
			}
			v22 = *v21;
			a4.field_4 = a4.field_4 + v16[4] - 1;
			if (sub_54BF20((int)v16, a1[37], &a4.field_0, a1 + 28, v22)) {
				break;
			}
			if (++a4.field_0 > v16[1] + v37) {
				a4.field_0 = v16[1];
			}
			if (++a2f > v37) {
				return 0;
			}
		}
		return 1;
	case 3:
		v23 = a3;
		a2g = 0;
		v24 = v23[3];
		v25 = v23[4];
		v26 = v23[2];
		a4.field_0 = v23[1] + v24 - 1;
		v27 = v24 - v3;
		a4.field_4 = v26 + v25 / 2;
		if (v27 < 0) {
			goto LABEL_46;
		}
		break;
	default:
		return 0;
	}
	while (1) {
		if (sub_54BB20((int)v23, a1[37], &a4.field_0, a1 + 32, a1[34])) {
			return 1;
		}
		v28 = v23[2];
		if (++a4.field_4 > v28 + v27) {
			a4.field_4 = v28;
		}
		if (++a2g > v27) {
			break;
		}
	}
LABEL_46:
	a2h = 0;
	v29 = v23[3];
	v38 = v29 - v3;
	a4.field_0 = v23[1] + v29 / 2;
	if (v29 - v3 < 0) {
		return 0;
	}
	v30 = a1 + 34;
	while (1) {
		a4.field_4 = v23[2];
		if (sub_54BD90((int)v23, a1[37], &a4.field_0, a1 + 32, *v30)) {
			break;
		}
		v31 = *v30;
		a4.field_4 = a4.field_4 + v23[4] - 1;
		if (sub_54BF20((int)v23, a1[37], &a4.field_0, a1 + 32, v31)) {
			break;
		}
		if (++a4.field_0 > v23[1] + v38) {
			a4.field_0 = v23[1];
		}
		if (++a2h > v38) {
			return 0;
		}
	}
	return 1;
}

//----- (0054B810) --------------------------------------------------------
int sub_54B810(int a1, int a2, int* a3, int2* a4, int a5) {
	int* v5;    // ebx
	int2* v6;   // esi
	int v7;     // eax
	float* v8;  // eax
	int v10;    // ecx
	int v11;    // edx
	int v12;    // eax
	int v13;    // ecx
	int v14;    // edi
	float* v15; // eax
	double v16; // st7
	float* v17; // eax
	float* v18; // eax
	float2 a2a; // [esp+10h] [ebp-8h]
	int v20;    // [esp+24h] [ebp+Ch]
	int2* v21;  // [esp+28h] [ebp+10h]
	int2* v22;  // [esp+28h] [ebp+10h]
	int v23;    // [esp+2Ch] [ebp+14h]
	int v24;    // [esp+2Ch] [ebp+14h]

	v5 = a3;
	v6 = a4;
	v7 = a4->field_0 - *a3;
	v20 = a4->field_0 - *a3;
	if (v20 < 0) {
		v7 = -v7;
	}
	*getMemU32Ptr(0x5D4594, 2491608) = 0;
	if (v7) {
		if (v7 < 3) {
			return 0;
		}
		v10 = v5[1];
		v11 = a4->field_4;
		v12 = v5[1];
		*getMemU32Ptr(0x5D4594, 2491608) = 0;
		v13 = v10 - (v12 - v11) / 2;
		v14 = v13 - 1;
		v21 = (int2*)(v13 - 1);
		v15 = sub_523E30(2, a5, v5[1] - (v13 - 1));
		*getMemU32Ptr(0x5D4594, 2491612 + 4 * *getMemU32Ptr(0x5D4594, 2491608)) = v15;
		a2a.field_0 = (double)*v5 * 32.526913;
		v16 = (double)(int)v21 * 32.526913;
		a2a.field_4 = v16;
		nox_xxx_mapGenSetRoomPos_521880(v15, &a2a);
		++*getMemU32Ptr(0x5D4594, 2491608);
		if (v20 <= 0) {
			v17 = sub_523E30(5, a5, *v5 - v6->field_0);
			a2a.field_4 = v16;
			*getMemU32Ptr(0x5D4594, 2491612 + 4 * *getMemU32Ptr(0x5D4594, 2491608)) = v17;
			a2a.field_0 = (double)v6->field_0 * 32.526913;
			sub_521A70(*getMemIntPtr(0x5D4594, 2491612), *(int*)&dword_5d4594_2491616, 3);
		} else {
			*getMemU32Ptr(0x5D4594, 2491612 + 4 * *getMemU32Ptr(0x5D4594, 2491608)) =
				sub_523E30(4, a5, v6->field_0 - *v5);
			v22 = (int2*)(a5 + *v5);
			a2a.field_4 = v16;
			a2a.field_0 = (double)(int)v22 * 32.526913;
			sub_521A70(*getMemIntPtr(0x5D4594, 2491612), *(int*)&dword_5d4594_2491616, 2);
		}
		nox_xxx_mapGenSetRoomPos_521880(*(uint32_t**)getMemAt(0x5D4594, 2491612 + 4 * *getMemU32Ptr(0x5D4594, 2491608)),
										&a2a);
		++*getMemU32Ptr(0x5D4594, 2491608);
		v18 = sub_523E30(2, a5, v14 - v6->field_4 - 1);
		*getMemU32Ptr(0x5D4594, 2491612 + 4 * *getMemU32Ptr(0x5D4594, 2491608)) = v18;
		v24 = v6->field_4 + 1;
		a2a.field_0 = (double)v6->field_0 * 32.526913;
		a2a.field_4 = (double)v24 * 32.526913;
		nox_xxx_mapGenSetRoomPos_521880(v18, &a2a);
		sub_521A70(*(int*)&dword_5d4594_2491616, *getMemIntPtr(0x5D4594, 2491620), 0);
	} else {
		v8 = sub_523E30(2, a5, v5[1] - a4->field_4 - 1);
		*getMemU32Ptr(0x5D4594, 2491612 + 4 * *getMemU32Ptr(0x5D4594, 2491608)) = v8;
		v23 = a4->field_4 + 1;
		a2a.field_0 = (double)a4->field_0 * 32.526913;
		a2a.field_4 = (double)v23 * 32.526913;
		nox_xxx_mapGenSetRoomPos_521880(v8, &a2a);
	}
	++*getMemU32Ptr(0x5D4594, 2491608);
	return sub_54BA60(a1, a2, 0, 1);
}

//----- (0054BA60) --------------------------------------------------------
int sub_54BA60(int a1, int a2, int a3, int a4) {
	int v4;            // eax
	int v5;            // esi
	int* v6;           // edi
	int v7;            // esi
	unsigned char* v8; // edi
	int result;        // eax
	int v10;           // edi
	void** v11;        // esi

	v4 = *getMemU32Ptr(0x5D4594, 2491608);
	v5 = 0;
	if (*getMemIntPtr(0x5D4594, 2491608) > 0) {
		v6 = getMemIntPtr(0x5D4594, 2491612);
		while (!sub_521200(*v6)) {
			v4 = *getMemU32Ptr(0x5D4594, 2491608);
			++v5;
			++v6;
			if (v5 >= *getMemIntPtr(0x5D4594, 2491608)) {
				goto LABEL_5;
			}
		}
		v10 = 0;
		if (*getMemU32Ptr(0x5D4594, 2491608) > 0) {
			v11 = (void**)getMemAt(0x5D4594, 2491612);
			do {
				sub_521A10(*v11);
				++v10;
				++v11;
			} while (v10 < *getMemIntPtr(0x5D4594, 2491608));
		}
		return 0;
	}
LABEL_5:
	v7 = 0;
	if (v4 > 0) {
		v8 = getMemAt(0x5D4594, 2491612);
		do {
			nox_xxx_mapGenAddNewRoom_521730(*(uint32_t**)v8);
			++v7;
			v8 += 4;
		} while (v7 < *getMemIntPtr(0x5D4594, 2491608));
	}
	sub_521A70(a1, *getMemIntPtr(0x5D4594, 2491612), a3);
	sub_521A70(a2, *getMemU32Ptr(0x5D4594, 2491608 + 4 * *getMemU32Ptr(0x5D4594, 2491608)), a4);
	return 1;
}

//----- (0054BB20) --------------------------------------------------------
int sub_54BB20(int a1, int a2, int* a3, uint32_t* a4, int a5) {
	uint32_t* v5; // ebp
	int v6;       // eax
	float* v7;    // eax
	int v9;       // ecx
	int v10;      // eax
	int v11;      // edi
	int v12;      // ebx
	float* v13;   // eax
	float* v14;   // eax
	float2 a2a;   // [esp+10h] [ebp-8h]
	int v16;      // [esp+28h] [ebp+10h]
	int v17;      // [esp+2Ch] [ebp+14h]

	v5 = a4;
	v6 = a4[1] - a3[1];
	v16 = v6;
	if (v6 < 0) {
		v6 = -v6;
	}
	*getMemU32Ptr(0x5D4594, 2491608) = 0;
	if (v6) {
		if (v6 < 3) {
			return 0;
		}
		v9 = *a3;
		v10 = *v5 - *a3;
		*getMemU32Ptr(0x5D4594, 2491608) = 0;
		v11 = v9 + v10 / 2;
		v12 = a5;
		v13 = sub_523E30(4, a5, v11 - *a3 + a5 - 1);
		*getMemU32Ptr(0x5D4594, 2491612 + 4 * *getMemU32Ptr(0x5D4594, 2491608)) = v13;
		a2a.field_0 = (double)(*a3 + 1) * 32.526913;
		a2a.field_4 = (double)a3[1] * 32.526913;
		nox_xxx_mapGenSetRoomPos_521880(v13, &a2a);
		++*getMemU32Ptr(0x5D4594, 2491608);
		if (v16 <= 0) {
			*getMemU32Ptr(0x5D4594, 2491612 + 4 * *getMemU32Ptr(0x5D4594, 2491608)) = sub_523E30(2, a5, a3[1] - v5[1]);
			a2a.field_0 = (double)v11 * 32.526913;
			a2a.field_4 = (double)(int)v5[1] * 32.526913;
			sub_521A70(*getMemIntPtr(0x5D4594, 2491612), *(int*)&dword_5d4594_2491616, 0);
		} else {
			*getMemU32Ptr(0x5D4594, 2491612 + 4 * *getMemU32Ptr(0x5D4594, 2491608)) = sub_523E30(3, a5, v5[1] - a3[1]);
			v17 = a3[1] + a5;
			a2a.field_0 = (double)v11 * 32.526913;
			a2a.field_4 = (double)v17 * 32.526913;
			sub_521A70(*getMemIntPtr(0x5D4594, 2491612), *(int*)&dword_5d4594_2491616, 1);
		}
		nox_xxx_mapGenSetRoomPos_521880(*(uint32_t**)getMemAt(0x5D4594, 2491612 + 4 * *getMemU32Ptr(0x5D4594, 2491608)),
										&a2a);
		++*getMemU32Ptr(0x5D4594, 2491608);
		v14 = sub_523E30(4, v12, *v5 - v11 - v12);
		*getMemU32Ptr(0x5D4594, 2491612 + 4 * *getMemU32Ptr(0x5D4594, 2491608)) = v14;
		a2a.field_0 = (double)(v12 + v11) * 32.526913;
		a2a.field_4 = (double)(int)v5[1] * 32.526913;
		nox_xxx_mapGenSetRoomPos_521880(v14, &a2a);
		sub_521A70(*(int*)&dword_5d4594_2491616, *getMemIntPtr(0x5D4594, 2491620), 2);
	} else {
		v7 = sub_523E30(4, a5, *v5 - *a3 - 1);
		*getMemU32Ptr(0x5D4594, 2491612 + 4 * *getMemU32Ptr(0x5D4594, 2491608)) = v7;
		a2a.field_0 = (double)(*a3 + 1) * 32.526913;
		a2a.field_4 = (double)a3[1] * 32.526913;
		nox_xxx_mapGenSetRoomPos_521880(v7, &a2a);
	}
	++*getMemU32Ptr(0x5D4594, 2491608);
	return sub_54BA60(a1, a2, 2, 3);
}

//----- (0054BD90) --------------------------------------------------------
int sub_54BD90(int a1, int a2, int* a3, int* a4, int a5) {
	int v6;     // ecx
	int v7;     // ebx
	int v8;     // ebx
	float* v9;  // eax
	int v10;    // esi
	float2 a2a; // [esp+Ch] [ebp-8h]

	if (a4[1] > a3[1] - a5) {
		return 0;
	}
	v6 = *a3;
	v7 = *a4;
	*getMemU32Ptr(0x5D4594, 2491608) = 0;
	v8 = v7 - v6;
	v9 = sub_523E30(2, a5, a3[1] - a4[1]);
	*getMemU32Ptr(0x5D4594, 2491612 + 4 * *getMemU32Ptr(0x5D4594, 2491608)) = v9;
	a2a.field_0 = (double)*a3 * 32.526913;
	a2a.field_4 = (double)a4[1] * 32.526913;
	nox_xxx_mapGenSetRoomPos_521880(v9, &a2a);
	++*getMemU32Ptr(0x5D4594, 2491608);
	if (v8 <= 0) {
		*getMemU32Ptr(0x5D4594, 2491612 + 4 * *getMemU32Ptr(0x5D4594, 2491608)) = sub_523E30(5, a5, *a3 - *a4 - 1);
		a2a.field_0 = (double)(*a4 + 1) * 32.526913;
		a2a.field_4 = (double)a4[1] * 32.526913;
		sub_521A70(*getMemIntPtr(0x5D4594, 2491612), *(int*)&dword_5d4594_2491616, 3);
		v10 = 2;
	} else {
		*getMemU32Ptr(0x5D4594, 2491612 + 4 * *getMemU32Ptr(0x5D4594, 2491608)) = sub_523E30(4, a5, *a4 - *a3 - a5);
		a2a.field_0 = (double)(a5 + *a3) * 32.526913;
		a2a.field_4 = (double)a4[1] * 32.526913;
		sub_521A70(*getMemIntPtr(0x5D4594, 2491612), *(int*)&dword_5d4594_2491616, 2);
		v10 = 3;
	}
	nox_xxx_mapGenSetRoomPos_521880(*(uint32_t**)getMemAt(0x5D4594, 2491612 + 4 * *getMemU32Ptr(0x5D4594, 2491608)),
									&a2a);
	++*getMemU32Ptr(0x5D4594, 2491608);
	return sub_54BA60(a1, a2, 0, v10);
}

//----- (0054BF20) --------------------------------------------------------
int sub_54BF20(int a1, int a2, int* a3, int* a4, int a5) {
	int* v5;    // esi
	int v7;     // ebx
	int v8;     // ebp
	float* v9;  // eax
	int v10;    // esi
	float2 a2a; // [esp+8h] [ebp-8h]
	int v12;    // [esp+1Ch] [ebp+Ch]

	v5 = a3;
	if (a4[1] <= a3[1]) {
		return 0;
	}
	v7 = *a3;
	v8 = *a4;
	*getMemU32Ptr(0x5D4594, 2491608) = 0;
	v9 = sub_523E30(3, a5, a4[1] - a3[1] + a5 - 1);
	*getMemU32Ptr(0x5D4594, 2491612 + 4 * *getMemU32Ptr(0x5D4594, 2491608)) = v9;
	v12 = a3[1] + 1;
	a2a.field_0 = (double)*v5 * 32.526913;
	a2a.field_4 = (double)v12 * 32.526913;
	nox_xxx_mapGenSetRoomPos_521880(v9, &a2a);
	++*getMemU32Ptr(0x5D4594, 2491608);
	if (v8 - v7 <= 0) {
		*getMemU32Ptr(0x5D4594, 2491612 + 4 * *getMemU32Ptr(0x5D4594, 2491608)) = sub_523E30(5, a5, *v5 - *a4 - 1);
		a2a.field_0 = (double)(*a4 + 1) * 32.526913;
		a2a.field_4 = (double)a4[1] * 32.526913;
		sub_521A70(*getMemIntPtr(0x5D4594, 2491612), *(int*)&dword_5d4594_2491616, 3);
		v10 = 2;
	} else {
		*getMemU32Ptr(0x5D4594, 2491612 + 4 * *getMemU32Ptr(0x5D4594, 2491608)) = sub_523E30(4, a5, *a4 - *v5 - a5);
		a2a.field_0 = (double)(a5 + *v5) * 32.526913;
		a2a.field_4 = (double)a4[1] * 32.526913;
		sub_521A70(*getMemIntPtr(0x5D4594, 2491612), *(int*)&dword_5d4594_2491616, 2);
		v10 = 3;
	}
	nox_xxx_mapGenSetRoomPos_521880(*(uint32_t**)getMemAt(0x5D4594, 2491612 + 4 * *getMemU32Ptr(0x5D4594, 2491608)),
									&a2a);
	++*getMemU32Ptr(0x5D4594, 2491608);
	return sub_54BA60(a1, a2, 1, v10);
}

//----- (0054D2B0) --------------------------------------------------------
int nox_xxx_diePlayer_54D2B0(int a1) {
	int v1;                // edi
	int v2;                // ebp
	int v3;                // esi
	uint32_t* v4;          // eax
	char* v5;              // eax
	int v6;                // ebx
	char v7;               // cl
	int v8;                // eax
	int v9;                // edx
	short v10;             // ax
	short v11;             // dx
	int v12;               // eax
	int v13;               // edx
	short v14;             // ax
	int v15;               // edx
	int v16;               // eax
	int v17;               // edx
	int result;            // eax
	int v19;               // eax
	int v20;               // eax
	float v21;             // [esp+0h] [ebp-28h]
	char* v22;             // [esp+14h] [ebp-14h]
	unsigned char v23[14]; // [esp+18h] [ebp-10h]
	int v24;               // [esp+2Ch] [ebp+4h]

	v1 = a1;
	v2 = 0;
	v3 = *(uint32_t*)(a1 + 748);
	if (!*getMemU32Ptr(0x5D4594, 2491688)) {
		*getMemU32Ptr(0x5D4594, 2491688) = nox_xxx_getNameId_4E3AA0("AnkhTradable");
	}
	if (nox_common_gameFlags_check_40A5C0(2048)) {
		sub_4DB170(0, 0, 0);
	}
	v24 = *(uint32_t*)(a1 + 520);
	if (v24) {
		v24 = nox_xxx_findParentChainPlayer_4EC580(v24);
	}
	v4 = *(uint32_t**)(v3 + 276);
	if (v4[900] && gameFrame() - v4[902] < (unsigned int)(10 * gameFPS())) {
		v5 = nox_common_playerInfoFromNum_417090(v4[901]);
		v6 = (int)v5;
		v22 = v5;
		if (v5) {
			if (*((uint32_t*)v5 + 523) && *((uint32_t*)v5 + 514)) {
				v2 = nox_server_getObjectFromNetCode_4ECCB0(*((uint32_t*)v5 + 515));
			} else {
				v22 = 0;
				v2 = 0;
				v6 = 0;
			}
			if (v24 == v2) {
				v2 = 0;
			}
			if (v1 == v2) {
				v2 = 0;
			}
		}
	} else {
		v22 = 0;
		v6 = 0;
	}
	if (!nox_common_gameFlags_check_40A5C0(0x2000)) {
		goto LABEL_38;
	}
	v7 = 0;
	v23[10] = 0;
	*(uint32_t*)v23 = 0;
	*(uint32_t*)&v23[4] = 0;
	*(uint16_t*)&v23[8] = 0;
	if (v24 && *(uint8_t*)(v24 + 8) & 4) {
		*(uint16_t*)&v23[2] = *(uint16_t*)(v24 + 36);
	}
	v8 = *(uint32_t*)(v1 + 520);
	if (v8) {
		v9 = *(uint32_t*)(v8 + 8);
		if (v9 & 2) {
			v23[10] = 1;
			v10 = *(uint16_t*)(v8 + 4);
			*(uint16_t*)&v23[8] = v10;
			goto LABEL_31;
		}
		if (v9 & 4) {
			v7 = *(uint8_t*)(v3 + 304);
			v11 = *(uint16_t*)(v3 + 300);
			v23[10] = *(uint8_t*)(v3 + 304);
			*(uint16_t*)&v23[8] = v11;
		} else {
			v12 = *(uint32_t*)(v8 + 508);
			if (v12) {
				v13 = *(uint32_t*)(v12 + 8);
				if (v13 & 4) {
					v7 = *(uint8_t*)(v3 + 304);
					v14 = *(uint16_t*)(v3 + 300);
					v23[10] = *(uint8_t*)(v3 + 304);
					*(uint16_t*)&v23[8] = v14;
				} else if (v13 & 2) {
					v23[10] = 1;
					*(uint16_t*)&v23[8] = *(uint16_t*)(v12 + 4);
					goto LABEL_31;
				}
			}
		}
	}
	if (!v7) {
		v10 = *(uint16_t*)(v3 + 300);
		v23[10] = *(uint8_t*)(v3 + 304);
		*(uint16_t*)&v23[8] = v10;
	}
LABEL_31:
	if (v2 && *(uint8_t*)(v2 + 8) & 4) {
		*(uint16_t*)&v23[4] = *(uint16_t*)(v2 + 36);
	}
	*(uint16_t*)&v23[6] = *(uint16_t*)(v1 + 36);
	nox_xxx_netInformTextMsg2_4DA180(14, v23);
	if (v23[10] == 2 && *(uint16_t*)&v23[8] == 2) {
		sub_4FC0B0(v24, 1);
	}
	v6 = (int)v22;
	*(uint32_t*)(v3 + 304) = 0;
LABEL_38:
	if (*(uint32_t*)(v1 + 524) == 16) {
		nox_xxx_aud_501960(299, v1, 0, 0);
	} else if (*(uint8_t*)(*(uint32_t*)(v3 + 276) + 2252)) {
		nox_xxx_aud_501960(331, v1, 0, 0);
	} else {
		nox_xxx_aud_501960(321, v1, 0, 0);
	}
	v15 = *(uint32_t*)(v1 + 16);
	BYTE1(v15) |= 0x80u;
	*(uint32_t*)(v1 + 16) = v15;
	nox_xxx_playerSetState_4FA020((uint32_t*)v1, 3);
	*(uint8_t*)(v3 + 188) = 0;
	*(uint32_t*)(v3 + 216) = 0;
	*(uint32_t*)(v3 + 192) = 0;
	*(uint32_t*)(v3 + 196) = 0;
	*(uint32_t*)(v3 + 200) = 0;
	*(uint32_t*)(v3 + 204) = 0;
	*(uint32_t*)(v3 + 208) = 0;
	*(uint8_t*)(v3 + 212) = 0;
	v16 = nox_xxx_gamePlayIsAnyPlayers_40A8A0();
	if (v16) {
		if (nox_common_gameFlags_check_40A5C0(256)) {
			nox_xxx_playerUpdateScore_54D980(v1, v24, v2, v6);
		} else if (nox_common_gameFlags_check_40A5C0(16)) {
			nox_xxx_playerHandleKotrDeath_54DC40(v1, v24);
		} else if (nox_common_gameFlags_check_40A5C0(1024)) {
			nox_xxx_playerHandleElimDeath_54D7A0(v1, v24);
		}
	}
	if (nox_common_gameFlags_check_40A5C0(1024) && nox_xxx_servGamedataGet_40A020(1024) &&
		*(uint32_t*)(*(uint32_t*)(v3 + 276) + 2140) >= (int)(unsigned short)nox_xxx_servGamedataGet_40A020(1024)) {
		nox_xxx_playerRemoveSpawnedStuff_4E5AD0(v1);
	}
	*(uint32_t*)(v1 + 16) |= 0x10u;
	nox_xxx_action_4DA9F0((uint32_t*)v1);
	if (!nox_common_gameFlags_check_40A5C0(4096)) {
		nox_xxx_dropAllItems_4EDA40((uint32_t*)v1);
	}
	nox_xxx_netNotifyPlayerDied_54DF00(v1);
	v17 = *(uint32_t*)(v3 + 276);
	*(uint16_t*)(v3 + 4) = 0;
	nox_xxx_protectMana_56F9E0(*(uint32_t*)(v17 + 4596), 0);
	nox_xxx_setUnitBuffFlags_4E48F0(v1, 0);
	nox_xxx_playerCancelAbils_4FC180(v1);
	*(uint32_t*)(*(uint32_t*)(v3 + 276) + 3600) = 0;
	nox_xxx_playerCancelSpells_4FEAE0(v1);
	nox_xxx_unitClearBuffs_4FF580(v1);
	if (*(uint32_t*)(v3 + 280)) {
		nox_xxx_shopCancelSession_510DC0(*(uint32_t**)(v3 + 280));
	}
	*(uint32_t*)(v3 + 280) = 0;
	result = nox_common_gameFlags_check_40A5C0(4096);
	if (result) {
		v19 = *(uint32_t*)(v3 + 320);
		if (v19) {
			*(uint32_t*)(v3 + 320) = v19 - 1;
			result = sub_4D6130(v1);
		} else {
			v20 = *(uint32_t*)(v3 + 276);
			*(uint32_t*)(v3 + 548) = gameFrame();
			v23[0] = -16;
			v23[1] = 2;
			*(uint16_t*)&v23[8] = *(uint16_t*)(v20 + 4688);
			*(uint16_t*)&v23[2] = *(uint16_t*)(v20 + 4668);
			*(uint16_t*)&v23[6] = *(uint16_t*)(v20 + 4664);
			*(uint16_t*)&v23[4] = *(uint16_t*)(v20 + 4672);
			*(uint32_t*)&v23[10] = 0;
			nox_xxx_netSendPacket0_4E5420(*(unsigned char*)(v20 + 2064), v23, 14, 0, 1);
			sub_4D6000(v1);
			sub_54CBD0(v1);
			v21 = nox_xxx_gamedataGetFloat_419D40("QuestGameStartingExtraLives");
			*(uint32_t*)(v3 + 320) = nox_float2int(v21);
			result = *(uint32_t*)(v3 + 276);
			*(uint8_t*)(*(unsigned char*)(result + 2064) + v3 + 452) = *(uint8_t*)(v3 + 320);
		}
	}
	return result;
}
// 54D56D: variable 'v16' is possibly undefined

//----- (0054D7A0) --------------------------------------------------------
void nox_xxx_playerHandleElimDeath_54D7A0(int a1, int a2) {
	int v2;   // edi
	int v3;   // ebx
	char* v4; // ebp
	int v5;   // eax
	int v6;   // [esp-4h] [ebp-18h]
	int v7;   // [esp-4h] [ebp-18h]
	char* v8; // [esp+10h] [ebp-4h]
	int v9;   // [esp+18h] [ebp+4h]

	v2 = a1;
	v3 = 0;
	v4 = 0;
	v8 = 0;
	v6 = a1 + 48;
	v9 = *(uint32_t*)(a1 + 748);
	if (nox_xxx_servObjectHasTeam_419130(v6)) {
		v8 = nox_xxx_getTeamByID_418AB0(*(unsigned char*)(v2 + 52));
	}
	if (a2) {
		v3 = *(uint32_t*)(a2 + 748);
		if (nox_xxx_servObjectHasTeam_419130(a2 + 48)) {
			v4 = nox_xxx_getTeamByID_418AB0(*(unsigned char*)(a2 + 52));
		}
	}
	if (a2 == v2) {
		nox_xxx_playerSubLessons_4D8EC0(v2, 1);
		nox_xxx_playerIncrementElimDeath_4D8D40(v2);
		nox_xxx_netReportLesson_4D8EF0(v2);
		if (v8) {
			nox_xxx_netChangeTeamID_419090((int)v8, *((uint32_t*)v8 + 13) + 1);
		}
		if (dword_5d4594_2650652) {
			if (v3) {
				sub_425CA0(*(uint32_t*)(v3 + 276), *(uint32_t*)(v3 + 276));
			}
		}
		return;
	}
	if (a2) {
		if (*(uint8_t*)(a2 + 8) & 4) {
			if (v4) {
				if (v4 == v8) {
					nox_xxx_playerSubLessons_4D8EC0(a2, 1);
					nox_xxx_netReportLesson_4D8EF0(a2);
					if (dword_5d4594_2650652 && v3) {
						sub_425CA0(*(uint32_t*)(v3 + 276), *(uint32_t*)(v3 + 276));
					}
					goto LABEL_32;
				}
			} else if (!v8) {
				nox_xxx_changeScore_4D8E90(a2, 1);
				nox_xxx_netReportLesson_4D8EF0(a2);
				if (dword_5d4594_2650652 && v3 && v9) {
					sub_425CA0(*(uint32_t*)(v3 + 276), *(uint32_t*)(v9 + 276));
				}
				goto LABEL_32;
			}
			nox_xxx_changeScore_4D8E90(a2, 1);
			nox_xxx_netReportLesson_4D8EF0(a2);
			if (dword_5d4594_2650652) {
				if (v3 && v9) {
					v5 = *(uint32_t*)(v3 + 276);
					v7 = *(uint32_t*)(v9 + 276);
					sub_425CA0(v5, v7);
					goto LABEL_32;
				}
			}
		}
	} else if (dword_5d4594_2650652 && v9) {
		v5 = *(uint32_t*)(v9 + 276);
		v7 = *(uint32_t*)(v9 + 276);
		sub_425CA0(v5, v7);
	}
LABEL_32:
	nox_xxx_playerIncrementElimDeath_4D8D40(v2);
	nox_xxx_netReportLesson_4D8EF0(v2);
	if (v8) {
		nox_xxx_netChangeTeamID_419090((int)v8, *((uint32_t*)v8 + 13) + 1);
	}
}

//----- (0054D980) --------------------------------------------------------
void nox_xxx_playerUpdateScore_54D980(int a1, int a2, int a3, int a4) {
	int v4;       // ebx
	char* v5;     // edi
	int v6;       // ebp
	int v7;       // eax
	char* result; // eax
	char* v9;     // esi
	int v10;      // edi
	int v11;      // ecx
	int v12;      // [esp-4h] [ebp-20h]
	int v13;      // [esp-4h] [ebp-20h]
	char* v14;    // [esp+10h] [ebp-Ch]
	char* v15;    // [esp+14h] [ebp-8h]
	char* v16;    // [esp+18h] [ebp-4h]
	int v17;      // [esp+20h] [ebp+4h]

	v4 = a1;
	v5 = 0;
	v12 = a1 + 48;
	v14 = 0;
	v15 = 0;
	v17 = *(uint32_t*)(a1 + 748);
	v6 = 0;
	v16 = 0;
	if (nox_xxx_servObjectHasTeam_419130(v12)) {
		v14 = nox_xxx_getTeamByID_418AB0(*(unsigned char*)(v4 + 52));
	}
	if (a2) {
		v6 = *(uint32_t*)(a2 + 748);
		if (nox_xxx_servObjectHasTeam_419130(a2 + 48)) {
			v5 = nox_xxx_getTeamByID_418AB0(*(unsigned char*)(a2 + 52));
		}
	}
	if (a4) {
		if (a3) {
			v16 = *(char**)(a3 + 748);
			if (nox_xxx_servObjectHasTeam_419130(a3 + 48)) {
				v15 = nox_xxx_getTeamByID_418AB0(*(unsigned char*)(a3 + 52));
			}
		}
	}
	if (a2 == v4) {
		goto LABEL_31;
	}
	if (a2) {
		if (!(*(uint8_t*)(a2 + 8) & 4)) {
			nox_xxx_playerIncrementElimDeath_4D8D40(v4);
			result = (char*)nox_xxx_netReportLesson_4D8EF0(v4);
			v9 = v14;
			goto LABEL_36;
		}
		if (v5) {
			if (v5 == v14) {
				nox_xxx_playerSubLessons_4D8EC0(a2, 1);
				nox_xxx_netReportLesson_4D8EF0(a2);
				nox_xxx_netChangeTeamID_419090((int)v5, *((uint32_t*)v5 + 13) - 1);
				if (!dword_5d4594_2650652 || !v6) {
					nox_xxx_playerIncrementElimDeath_4D8D40(v4);
					result = (char*)nox_xxx_netReportLesson_4D8EF0(v4);
					v9 = v14;
					goto LABEL_36;
				}
				v7 = *(uint32_t*)(v6 + 276);
				v13 = *(uint32_t*)(v6 + 276);
				sub_425CA0(v7, v13);
				nox_xxx_playerIncrementElimDeath_4D8D40(v4);
				result = (char*)nox_xxx_netReportLesson_4D8EF0(v4);
				v9 = v14;
				goto LABEL_36;
			}
		} else if (!v14) {
			nox_xxx_changeScore_4D8E90(a2, 1);
			nox_xxx_netReportLesson_4D8EF0(a2);
			if (!dword_5d4594_2650652 || !v6 || !v17) {
				nox_xxx_playerIncrementElimDeath_4D8D40(v4);
				result = (char*)nox_xxx_netReportLesson_4D8EF0(v4);
				v9 = v14;
				goto LABEL_36;
			}
			v7 = *(uint32_t*)(v6 + 276);
			v13 = *(uint32_t*)(v17 + 276);
			sub_425CA0(v7, v13);
			nox_xxx_playerIncrementElimDeath_4D8D40(v4);
			result = (char*)nox_xxx_netReportLesson_4D8EF0(v4);
			v9 = v14;
			goto LABEL_36;
		}
		nox_xxx_changeScore_4D8E90(a2, 1);
		nox_xxx_netReportLesson_4D8EF0(a2);
		nox_xxx_netChangeTeamID_419090((int)v5, *((uint32_t*)v5 + 13) + 1);
		if (dword_5d4594_2650652 && v6 && v17) {
			sub_425CA0(*(uint32_t*)(v6 + 276), *(uint32_t*)(v17 + 276));
		}
		nox_xxx_playerIncrementElimDeath_4D8D40(v4);
		result = (char*)nox_xxx_netReportLesson_4D8EF0(v4);
		v9 = v14;
		goto LABEL_36;
	}
	if (a3) {
		nox_xxx_playerIncrementElimDeath_4D8D40(v4);
		result = (char*)nox_xxx_netReportLesson_4D8EF0(v4);
		v9 = v14;
		goto LABEL_36;
	}
LABEL_31:
	nox_xxx_playerSubLessons_4D8EC0(v4, 1);
	nox_xxx_netReportLesson_4D8EF0(v4);
	v9 = v14;
	if (v14) {
		nox_xxx_netChangeTeamID_419090((int)v14, *((uint32_t*)v14 + 13) - 1);
	}
	result = *(char**)&dword_5d4594_2650652;
	if (dword_5d4594_2650652 && v6) {
		result = sub_425CA0(*(uint32_t*)(v6 + 276), *(uint32_t*)(v6 + 276));
	}
LABEL_36:
	if (!a3) {
		return;
	}
	if (v5) {
		result = v15;
		if (v5 == v15) {
			return;
		}
		v10 = (int)v15;
	} else {
		v10 = (int)v15;
		if (!v15) {
			goto LABEL_44;
		}
	}
	if (v9) {
		if (v9 == (char*)v10) {
			return;
		}
	}
	if (v10) {
		nox_xxx_changeScore_4D8E90(a3, 1);
		nox_xxx_netReportLesson_4D8EF0(a3);
		if (v10) {
			nox_xxx_netChangeTeamID_419090(v10, *(uint32_t*)(v10 + 52) + 1);
		}
		result = *(char**)&dword_5d4594_2650652;
		if (dword_5d4594_2650652) {
			result = v16;
			if (v16) {
				v11 = v17;
				if (v17) {
					sub_425CA0(*((uint32_t*)result + 69), *(uint32_t*)(v11 + 276));
					return;
				}
			}
		}
		return;
	}
LABEL_44:
	nox_xxx_changeScore_4D8E90(a3, 1);
	nox_xxx_netReportLesson_4D8EF0(a3);
	if (dword_5d4594_2650652) {
		result = v16;
		if (v16) {
			v11 = v17;
			if (v17) {
				sub_425CA0(*((uint32_t*)result + 69), *(uint32_t*)(v11 + 276));
				return;
			}
		}
	}
}

//----- (0054DC40) --------------------------------------------------------
void nox_xxx_playerHandleKotrDeath_54DC40(int a1, int a2) {
	char* v2;     // edi
	char* v3;     // ebx
	char* result; // eax
	int v5;       // ebp
	double v6;    // st7
	int v7;       // ebx
	int v8;       // ebx
	int v9;       // eax
	double v10;   // st7
	int v11;      // eax
	float v12;    // [esp+0h] [ebp-18h]
	float v13;    // [esp+0h] [ebp-18h]
	int v14;      // [esp+0h] [ebp-18h]
	float v15;    // [esp+0h] [ebp-18h]
	int v16;      // [esp+14h] [ebp-4h]

	v2 = 0;
	v3 = 0;
	v16 = *(uint32_t*)(a1 + 748);
	result = (char*)nox_xxx_servObjectHasTeam_419130(a1 + 48);
	if (result) {
		result = nox_xxx_getTeamByID_418AB0(*(unsigned char*)(a1 + 52));
		v3 = result;
	}
	if (a2) {
		v5 = *(uint32_t*)(a2 + 748);
		result = (char*)nox_xxx_servObjectHasTeam_419130(a2 + 48);
		if (result) {
			result = nox_xxx_getTeamByID_418AB0(*(unsigned char*)(a2 + 52));
			v2 = result;
		}
		if (*(uint8_t*)(a2 + 8) & 4) {
			if (a2 == a1 || v3 == v2 && v3) {
				if (!nox_xxx_unitIsCrown_4E7BE0(a2)) {
					nox_xxx_playerIncrementElimDeath_4D8D40(a1);
					nox_xxx_netReportLesson_4D8EF0(a1);
					return;
				}
				nox_xxx_playerSubLessons_4D8EC0(a2, 1);
				nox_xxx_netReportLesson_4D8EF0(a2);
				if (v2) {
					nox_xxx_netChangeTeamID_419090((int)v2, *((uint32_t*)v2 + 13) - 1);
				}
				if (!dword_5d4594_2650652 || !v5) {
					nox_xxx_playerIncrementElimDeath_4D8D40(a1);
					nox_xxx_netReportLesson_4D8EF0(a1);
					return;
				}
				v9 = *(uint32_t*)(v5 + 276);
				v14 = *(uint32_t*)(v5 + 276);
			} else {
				if (!v2 || v2 == v3) {
					if (nox_xxx_unitIsCrown_4E7BE0(a2) || nox_xxx_unitIsCrown_4E7BE0(a1)) {
						if (nox_xxx_unitIsCrown_4E7BE0(a2)) {
							v10 = nox_xxx_gamedataGetFloat_419D40("KotRKingKillsPawnPoints");
						} else {
							v10 = nox_xxx_gamedataGetFloat_419D40("KotRPawnKillsKingPoints");
						}
						v15 = v10;
						v11 = nox_float2int(v15);
						nox_xxx_changeScore_4D8E90(a2, v11);
						nox_xxx_netReportLesson_4D8EF0(a2);
						if (dword_5d4594_2650652 && v5 && v16) {
							sub_425CA0(*(uint32_t*)(v5 + 276), *(uint32_t*)(v16 + 276));
						}
						if (!nox_xxx_CheckGameplayFlags_417DA0(4) && nox_xxx_unitIsCrown_4E7BE0(a1)) {
							sub_4ED050(a1, a2);
						}
					}
					nox_xxx_playerIncrementElimDeath_4D8D40(a1);
					nox_xxx_netReportLesson_4D8EF0(a1);
					return;
				}
				if (nox_xxx_unitIsCrown_4E7BE0(a2)) {
					if (nox_xxx_unitIsCrown_4E7BE0(a1)) {
						v6 = nox_xxx_gamedataGetFloat_419D40("KotRKingKillsKingPoints");
					} else {
						v6 = nox_xxx_gamedataGetFloat_419D40("KotRKingKillsPawnPoints");
					}
					v12 = v6;
					v7 = nox_float2int(v12);
					nox_xxx_changeScore_4D8E90(a2, v7);
					nox_xxx_netChangeTeamID_419090((int)v2, v7 + *((uint32_t*)v2 + 13));
					nox_xxx_netReportLesson_4D8EF0(a2);
					if (dword_5d4594_2650652 && v5) {
						if (v16) {
							sub_425CA0(*(uint32_t*)(v5 + 276), *(uint32_t*)(v16 + 276));
						}
					}
					nox_xxx_playerIncrementElimDeath_4D8D40(a1);
					nox_xxx_netReportLesson_4D8EF0(a1);
					return;
				}
				if (!nox_xxx_unitIsCrown_4E7BE0(a1) ||
					(v13 = nox_xxx_gamedataGetFloat_419D40("KotRPawnKillsKingPoints"), v8 = nox_float2int(v13),
					 nox_xxx_changeScore_4D8E90(a2, v8),
					 nox_xxx_netChangeTeamID_419090((int)v2, v8 + *((uint32_t*)v2 + 13)),
					 nox_xxx_netReportLesson_4D8EF0(a2), !dword_5d4594_2650652) ||
					!v5 || !v16) {
					nox_xxx_playerIncrementElimDeath_4D8D40(a1);
					nox_xxx_netReportLesson_4D8EF0(a1);
					return;
				}
				v9 = *(uint32_t*)(v5 + 276);
				v14 = *(uint32_t*)(v16 + 276);
			}
			sub_425CA0(v9, v14);
			nox_xxx_playerIncrementElimDeath_4D8D40(a1);
			nox_xxx_netReportLesson_4D8EF0(a1);
			return;
		}
	}
}

//----- (0054DF00) --------------------------------------------------------
void nox_xxx_netNotifyPlayerDied_54DF00(int a1) {
	int v1 = 0; // ecx
	short v2;   // cx
	int v4;     // [esp+0h] [ebp-4h]

	v4 = v1;
	v2 = *(uint16_t*)(a1 + 36);
	LOBYTE(v4) = -24;
	*(uint16_t*)((char*)&v4 + 1) = v2;
	nox_xxx_netSendPacket1_4E5390(255, (int)&v4, 3, 0, 0);
}

//----- (0054E630) --------------------------------------------------------
void nox_xxx_dieMonsterGen_54E630(int a1) {
	int v1;       // edi
	int v2;       // eax
	uint32_t* v3; // eax

	v1 = *(uint32_t*)(a1 + 748);
	sub_4D71E0(gameFrame());
	sub_4D7520(0);
	nox_xxx_scriptCallByEventBlock_502490((int*)(v1 + 56), *(uint32_t*)(a1 + 520), a1, 3);
	nox_xxx_aud_501960(1000, a1, 0, 0);
	nox_xxx_sendGeneratorBreakFX_523200((float*)(a1 + 56), 200);
	if (nox_common_gameFlags_check_40A5C0(4096)) {
		if (*(uint32_t*)(a1 + 520)) {
			v2 = nox_xxx_findParentChainPlayer_4EC580(*(uint32_t*)(a1 + 520));
			if (*(uint8_t*)(v2 + 8) & 4) {
				sub_4D61B0(v2);
			}
		}
	}
	v3 = nox_xxx_newObjectByTypeID_4E3810("DestroyedGenerator");
	if (v3) {
		nox_xxx_createAt_4DAA50((int)v3, 0, *(float*)(a1 + 56), *(float*)(a1 + 60));
	}
	nox_xxx_delayedDeleteObject_4E5CC0(a1);
}

//----- (0054E6F0) --------------------------------------------------------
int sub_54E6F0(int a1, int a2) {
	int result; // eax

	result = sub_54E730(a2, a1);
	if (result) {
		result = !nox_xxx_unitsHaveSameTeam_4EC520(a1, a2) || nox_xxx_GetGameplayFlags_417D90() & 1;
	}
	return result;
}

//----- (0054E730) --------------------------------------------------------
int sub_54E730(int a1, int a2) {
	int v2;     // ecx
	int v3;     // eax
	int result; // eax
	int v5;     // eax

	if (*(uint8_t*)(a2 + 8) & 1) {
		return 0;
	}
	v2 = *(uint32_t*)(a1 + 16);
	if (v2 & 0x20) {
		return 0;
	}
	v3 = *(uint32_t*)(a2 + 16);
	if (v3 & 0x20 || !*(uint32_t*)(a1 + 696) || !*(uint32_t*)(a2 + 696) || v3 & 0x40) {
		return 0;
	}
	if ((v3 & 0x80u) != 0) {
		return 1;
	}
	if (v2 & 0x11 && v3 & 0x4000 || v3 & 0x11 && v2 & 0x4000 ||
		(v2 & 0x400 || v3 & 0x400) && nox_xxx_unitsHaveSameTeam_4EC520(a2, a1) ||
		(v5 = *(uint32_t*)(a1 + 508)) != 0 && *(uint8_t*)(a1 + 8) & 1 && !(*(uint8_t*)(a1 + 12) & 2) &&
			*(uint8_t*)(v5 + 8) & 2 && *(uint8_t*)(a2 + 8) & 2 &&
			(!nox_xxx_unitIsEnemyTo_5330C0(v5, a2) || nox_xxx_unitsHaveSameTeam_4EC520(a2, *(uint32_t*)(a1 + 508)))) {
		return 0;
	} else {
		return 1;
	}
}

//----- (0054E810) --------------------------------------------------------
int sub_54E810(int a1, float2* a2, int a3) {
	int a3a[4]; // [esp+0h] [ebp-10h]

	a3a[0] = a1;
	a3a[1] = 0;
	a3a[2] = a3;
	a3a[3] = (int)a2;
	sub_517B70(a2, sub_54E850, (int)a3a);
	return a3a[1];
}

//----- (0054E850) --------------------------------------------------------
void sub_54E850(int a1, int a2) {
	int v2;       // eax
	float2* v3;   // ecx
	float2* v4;   // edx
	float v5;     // edx
	uint32_t* v6; // edx
	uint32_t* v7; // eax
	float4 a2a;   // [esp+8h] [ebp-20h]
	float4 a1a;   // [esp+18h] [ebp-10h]

	if ((signed char)*(uint8_t*)(a1 + 8) >= 0) {
		if (sub_54E730(*(uint32_t*)a2, a1) && sub_547DB0(a1, *(float2**)(a2 + 12))) {
			*(uint32_t*)(a2 + 4) = a1;
		}
	} else {
		v2 = *(uint32_t*)(a1 + 748);
		if (!(*(uint8_t*)(a1 + 12) & 4)) {
			v3 = *(float2**)(a2 + 8);
			a1a.field_0 = v3->field_0;
			v4 = *(float2**)(a2 + 12);
			a1a.field_4 = v3->field_4;
			a1a.field_8 = v4->field_0;
			v5 = v4->field_4;
			a2a.field_0 = *(float*)(a1 + 56);
			a1a.field_C = v5;
			a2a.field_4 = *(float*)(a1 + 60);
			a2a.field_8 = (double)*getMemIntPtr(0x587000, 196184 + 8 * *(uint32_t*)(v2 + 12)) + a2a.field_0;
			a2a.field_C = (double)*getMemIntPtr(0x587000, 196188 + 8 * *(uint32_t*)(v2 + 12)) + a2a.field_4;
			if (sub_427980(&a1a, &a2a)) {
				v6 = *(uint32_t**)(a2 + 8);
				v7 = *(uint32_t**)(a2 + 12);
				*v7 = *v6;
				v7[1] = v6[1];
				*(uint32_t*)(a2 + 4) = a1;
			}
		}
	}
}

//----- (0054E930) --------------------------------------------------------
char nox_xxx_updateMonsterGenerator_54E930(uint32_t* a1) {
	unsigned int v1; // esi
	int v2;          // edi
	int v3;          // ebp
	unsigned int v4; // eax
	unsigned int v5; // ebx
	unsigned int v6; // esi
	int v7;          // ecx
	double v8;       // st7
	int v9;          // ecx
	uint32_t* v10;   // eax
	int v11;         // edx
	int v12;         // esi
	float v14;       // [esp+0h] [ebp-24h]
	float v15;       // [esp+0h] [ebp-24h]
	float v16;       // [esp+0h] [ebp-24h]
	float v17;       // [esp+0h] [ebp-24h]
	float v18;       // [esp+0h] [ebp-24h]
	float v19;       // [esp+0h] [ebp-24h]
	float v20;       // [esp+0h] [ebp-24h]
	float v21;       // [esp+0h] [ebp-24h]
	long long v22;   // [esp+14h] [ebp-10h]
	float2 a2;       // [esp+1Ch] [ebp-8h]

	v1 = nox_game_getQuestStage_4E3CC0();
	v2 = a1[187];
	v3 = nox_xxx_getQuestStage_51A930();
	if (!dword_5d4594_2491716) {
		v14 = nox_xxx_gamedataGetFloat_419D40("QuestHardcoreStage");
		dword_5d4594_2491716 = nox_float2int(v14);
		v15 = nox_xxx_gamedataGetFloat_419D40("QuestHardcoreSpawnRateIncrease");
		*getMemU32Ptr(0x5D4594, 2491720) = nox_float2int(v15);
		*getMemFloatPtr(0x5D4594, 2491744) = nox_xxx_gamedataGetFloat_419D40("QuestHardcoreSpawnCap");
		v16 = nox_xxx_gamedataGetFloat_419D40("SpawnRateHighValue");
		*getMemU32Ptr(0x5D4594, 2491724) = nox_float2int(v16);
		v17 = nox_xxx_gamedataGetFloat_419D40("SpawnRateNormalValue");
		*getMemU32Ptr(0x5D4594, 2491728) = nox_float2int(v17);
		v18 = nox_xxx_gamedataGetFloat_419D40("SpawnRateLowValue");
		*getMemU32Ptr(0x5D4594, 2491732) = nox_float2int(v18);
		v19 = nox_xxx_gamedataGetFloat_419D40("SpawnRateVeryLowValue");
		*getMemU32Ptr(0x5D4594, 2491736) = nox_float2int(v19);
		v20 = nox_xxx_gamedataGetFloat_419D40("SpawnRateVeryVeryLowValue");
		*getMemU32Ptr(0x5D4594, 2491740) = nox_float2int(v20);
	}
	v4 = a1[5];
	if (!(v4 & 0x800)) {
		v4 = a1[4];
		if (v4 & 0x1000000) {
			if (!(v4 & 0x8020)) {
				nox_xxx_unitNeedSync_4E44F0((int)a1);
				v5 = gameFrame() - *(uint32_t*)(v2 + 88);
				switch (*(unsigned char*)(v2 + v3 + 80)) {
				case 0u:
					v4 = *getMemU32Ptr(0x5D4594, 2491724);
					break;
				case 1u:
					v4 = *getMemU32Ptr(0x5D4594, 2491728);
					break;
				case 2u:
					v4 = *getMemU32Ptr(0x5D4594, 2491732);
					break;
				case 3u:
					v4 = *getMemU32Ptr(0x5D4594, 2491736);
					break;
				case 4u:
					v4 = *getMemU32Ptr(0x5D4594, 2491740);
					break;
				default:
					v4 = (unsigned int)a1;
					break;
				}
				if (v1 >= *(int*)&dword_5d4594_2491716) {
					v6 = *getMemU32Ptr(0x5D4594, 2491720) * (v1 - dword_5d4594_2491716 + 1);
					if (v6 > v4) {
						v7 = 0;
					} else {
						v7 = v4 - v6;
					}
					v22 = v4;
					v8 = (double)v4;
					LODWORD(v22) = v7;
					if ((double)v22 / v8 < *getMemFloatPtr(0x5D4594, 2491744)) {
						v21 = v8 * *getMemFloatPtr(0x5D4594, 2491744);
						v7 = nox_float2int(v21);
					}
					v4 = v7;
				}
				if (v5 > v4) {
					LOBYTE(v4) = *(uint8_t*)(v2 + 86);
					if ((unsigned char)v4 < *(uint8_t*)(v2 + 87)) {
						if ((unsigned char)gameFrame() & 8) {
							v9 = 0;
							v10 = (uint32_t*)(v2 + 16 * v3);
							v11 = 4;
							do {
								if (*v10) {
									++v9;
								}
								++v10;
								--v11;
							} while (v11);
							v12 = *(uint32_t*)(v2 + 4 * (nox_common_randomInt_415FA0(0, v9 - 1) + 4 * v3));
							v4 = nox_xxx_mobGeneratorPick_54EBA0(a1, &a2, v12);
							if (v4 == 1) {
								LOBYTE(v4) = (unsigned int)nox_xxx_mobGeneratorSpawn_54F070((int)a1, (int)&a2, v12);
								*(uint32_t*)(v2 + 88) = gameFrame();
							}
						}
					}
				}
			}
		}
	}
	return v4;
}

//----- (0054EBA0) --------------------------------------------------------
int nox_xxx_mobGeneratorPick_54EBA0(uint32_t* a1, float2* a2, int a4) {
	float* v3;   // esi
	int v4;      // ebx
	int v5;      // edx
	int v6;      // eax
	float v7;    // edx
	float v8;    // ecx
	float v9;    // edx
	int v10;     // eax
	char* v11;   // eax
	int v13;     // [esp+18h] [ebp-A0h]
	float v14;   // [esp+1Ch] [ebp-9Ch]
	int* v15;    // [esp+20h] [ebp-98h]
	int v16;     // [esp+24h] [ebp-94h]
	float4 a1a;  // [esp+28h] [ebp-90h]
	int v18[32]; // [esp+38h] [ebp-80h]

	v13 = 0;
	v16 = a1[187];
	v3 = (float*)nox_xxx_getFirstPlayerUnit_4DA7C0();
	if (!v3) {
		return 0;
	}
	v15 = v18;
	do {
		v4 = *((uint32_t*)v3 + 187);
		if (!((uint32_t)v3[4] & 0x8020)) {
			v5 = *(uint32_t*)(v4 + 276);
			if (!(*(uint8_t*)(v5 + 3680) & 1)) {
				v6 = *(uint32_t*)(v16 + 92);
				if (v6 & 1) {
					if (nox_xxx_calcDistance_4E6C00((int)a1, (int)v3) <= 300.0) {
						return nox_xxx_mgenSetCreaturePos_54ED50((int)a1, a2, 0, a4);
					}
				} else {
					if (!(v6 & 2)) {
						return 0;
					}
					v14 = nox_double2float(sqrt((double)(*(unsigned short*)(v5 + 12) * *(unsigned short*)(v5 + 12) +
														 *(unsigned short*)(v5 + 10) * *(unsigned short*)(v5 + 10))));
					if (nox_xxx_calcDistance_4E6C00((int)a1, (int)v3) <= v14) {
						v7 = *((float*)a1 + 14);
						v8 = v3[14];
						a1a.field_4 = *((float*)a1 + 15);
						a1a.field_0 = v7;
						v9 = v3[15];
						a1a.field_8 = v8;
						a1a.field_C = v9;
						if (nox_xxx_mapTraceRay_535250(&a1a, 0, 0, 69)) {
							*v15 = *(unsigned char*)(*(uint32_t*)(v4 + 276) + 2064);
							++v13;
							++v15;
						}
					}
				}
			}
		}
		v3 = (float*)nox_xxx_getNextPlayerUnit_4DA7F0((int)v3);
	} while (v3);
	if (!v13) {
		return 0;
	}
	v10 = nox_common_randomInt_415FA0(0, v13 - 1);
	v11 = nox_common_playerInfoFromNum_417090(v18[v10]);
	return nox_xxx_mgenSetCreaturePos_54ED50((int)a1, a2, *((uint32_t*)v11 + 514), a4);
}
// 54EBA0: using guessed type int var_80[32];

//----- (0054ED50) --------------------------------------------------------
int nox_xxx_mgenSetCreaturePos_54ED50(int a1, float2* a2, int a3, int a4) {
	int v4;     // ecx
	char v5;    // al
	float v6;   // ecx
	short v7;   // ax
	float2 a1a; // [esp+10h] [ebp-18h]
	float4 v10; // [esp+18h] [ebp-10h]

	if (*(uint8_t*)(*(uint32_t*)(a1 + 748) + 92) & 2 && a3) {
		a1a.field_0 = *(float*)(a3 + 56) - *(float*)(a1 + 56);
		a1a.field_4 = *(float*)(a3 + 60) - *(float*)(a1 + 60);
		if (a1a.field_0 == 0.0) {
			a1a.field_0 = a1a.field_0 + 1.0;
		}
		if (a1a.field_4 == 0.0) {
			a1a.field_4 = a1a.field_4 + 1.0;
		}
		nox_xxx_utilNormalizeVector_509F20(&a1a);
		a1a.field_0 = a1a.field_0 * 45.0 + *(float*)(a1 + 56);
		a1a.field_4 = a1a.field_4 * 45.0 + *(float*)(a1 + 60);
		if (!sub_54EF00(&a1a.field_0)) {
			v10.field_8 = a1a.field_0;
			v10.field_0 = *(float*)(a1 + 56);
			v4 = *(uint32_t*)(a4 + 16);
			v10.field_C = a1a.field_4;
			v5 = 1;
			v10.field_4 = *(float*)(a1 + 60);
			if (v4 & 0x4000) {
				v5 = 5;
			}
			if (nox_xxx_mapTraceRay_535250(&v10, 0, 0, v5) && !nox_xxx_mapTileAllowTeleport_411A90(&a1a)) {
				v6 = a1a.field_4;
				a2->field_0 = a1a.field_0;
				a2->field_4 = v6;
				a1a.field_0 = *(float*)(a3 + 56) - a1a.field_0;
				a1a.field_4 = *(float*)(a3 + 60) - a1a.field_4;
				v7 = nox_xxx_math_509ED0(&a1a);
				*(uint16_t*)(a4 + 124) = v7;
				*(uint16_t*)(a4 + 126) = v7;
				return 1;
			}
		}
	}
	if (sub_54EF90(45.0, a1 + 56, (int)a2, a4) == 1) {
		if (!a3) {
			return 1;
		}
		v10.field_0 = *(float*)(a3 + 56) - *(float*)(a1 + 56);
		v10.field_4 = *(float*)(a3 + 60) - *(float*)(a1 + 60);
		v7 = nox_xxx_math_509ED0((float2*)&v10);
		*(uint16_t*)(a4 + 124) = v7;
		*(uint16_t*)(a4 + 126) = v7;
		return 1;
	}
	return 0;
}

//----- (0054EF00) --------------------------------------------------------
int sub_54EF00(float* a3) {
	float4 a1; // [esp+0h] [ebp-10h]

	*getMemU32Ptr(0x5D4594, 2491708) = 0;
	a1.field_0 = *a3 - 15.0;
	a1.field_4 = a3[1] - 15.0;
	a1.field_8 = *a3 + 15.0;
	a1.field_C = a3[1] + 15.0;
	nox_xxx_getUnitsInRect_517C10(&a1, sub_54EF60, (int)a3);
	return *getMemU32Ptr(0x5D4594, 2491708);
}

//----- (0054EF60) --------------------------------------------------------
void sub_54EF60(float* a1, int a2) {
	int v2; // ecx

	if (!((uint32_t)a1[2] & 0x20000)) {
		v2 = *((uint32_t*)a1 + 5);
		if (!(v2 & 0x800) && sub_547DB0((int)a1, (float2*)a2) == 1) {
			*getMemU32Ptr(0x5D4594, 2491708) = 1;
		}
	}
}

//----- (0054EF90) --------------------------------------------------------
int sub_54EF90(float a1, int a2, int a3, int a4) {
	float* v4;      // esi
	float v5;       // ecx
	char v6;        // bl
	int v7;         // eax
	int v8;         // edi
	long double v9; // st7
	float v11;      // edx
	float4 v12;     // [esp+Ch] [ebp-10h]
	int v13;        // [esp+24h] [ebp+8h]

	v4 = (float*)a2;
	v5 = *(float*)(a2 + 4);
	v6 = 1;
	v12.field_0 = *(float*)a2;
	v12.field_4 = v5;
	*(float*)&v13 = nox_common_randomFloat_416030(-3.1415927, 3.1415927);
	v7 = *(uint32_t*)(a4 + 16);
	if (v7 & 0x4000) {
		v6 = 5;
	}
	v8 = 0;
	while (1) {
		v9 = *(float*)&v13 + 1.8849558;
		*(float*)&v13 = v9;
		v12.field_8 = cos(v9) * a1 + *v4;
		v12.field_C = sin(*(float*)&v13) * a1 + v4[1];
		if (nox_xxx_mapTraceRay_535250(&v12, 0, 0, v6)) {
			if (!sub_54EF00(&v12.field_8) && !nox_xxx_mapTileAllowTeleport_411A90((float2*)&v12.field_8)) {
				break;
			}
		}
		if (++v8 >= 32) {
			return 0;
		}
	}
	v11 = v12.field_C;
	*(float*)a3 = v12.field_8;
	*(float*)(a3 + 4) = v11;
	return 1;
}

//----- (0054F070) --------------------------------------------------------
void* nox_xxx_objectTypeByIndHealthData(int a1);
uint32_t* nox_xxx_mobGeneratorSpawn_54F070(int a1, int a2, int a3) {
	int v3;           // edi
	uint32_t* result; // eax
	uint32_t* v5;     // esi
	int v6;           // ebp
	int v7;           // ebp
	int v8;           // ebp
	uint16_t* v9;     // eax
	int v10;          // eax
	int v11;          // eax
	float v12;        // edx
	int v13;          // eax
	double v14;       // st7
	float v15;        // [esp+0h] [ebp-30h]
	float v16;        // [esp+0h] [ebp-30h]
	float v17;        // [esp+0h] [ebp-30h]
	float v18;        // [esp+4h] [ebp-2Ch]
	float v19;        // [esp+4h] [ebp-2Ch]
	int v20;          // [esp+14h] [ebp-1Ch]
	float2 a1a;       // [esp+18h] [ebp-18h]
	int4 v22;         // [esp+20h] [ebp-10h]
	int v23;          // [esp+34h] [ebp+4h]

	v3 = a1;
	v20 = *(uint32_t*)(a1 + 748);
	nox_xxx_getQuestStage_51A930();
	if (!*getMemU32Ptr(0x5D4594, 2491712)) {
		*getMemU32Ptr(0x5D4594, 2491712) = nox_xxx_getNameId_4E3AA0("Beholder");
	}
	result = (uint32_t*)sub_50DE80(a1, (float*)a2);
	if (result) {
		result = nox_xxx_newObjectWithTypeInd_4E3450(*(unsigned short*)(a3 + 4));
		v5 = result;
		if (result) {
			v6 = result[187];
			v23 = result[187];
			nox_xxx_unitCreatureCopyUC_54F2B0(a3, (int)result);
			v7 = *(uint32_t*)(v6 + 484);
			if (v7) {
				v15 = sub_4E40F0() * (double)*(int*)(v7 + 72);
				*(uint16_t*)v5[139] = nox_float2int(v15);
				v16 = sub_4E40F0() * (double)*(int*)(v7 + 72);
			} else {
				v8 = nox_xxx_objectTypeByIndHealthData(*((unsigned short*)v5 + 2));
				v17 = sub_4E40F0() * (double)*(unsigned short*)v8;
				*(uint16_t*)v5[139] = nox_float2int(v17);
				v16 = sub_4E40F0() * (double)*(unsigned short*)(v8 + 4);
			}
			*(uint16_t*)(v5[139] + 4) = nox_float2int(v16);
			v9 = (uint16_t*)v5[139];
			if (!*v9) {
				*v9 = 1;
			}
			v10 = v5[139];
			if (!*(uint16_t*)(v10 + 4)) {
				*(uint16_t*)(v10 + 4) = 1;
			}
			if (*((unsigned short*)v5 + 2) == *getMemU32Ptr(0x5D4594, 2491712)) {
				*(uint32_t*)(v23 + 1504) = 0;
			}
			if (sub_50E030(v3, v5)) {
				nox_xxx_createAt_4DAA50((int)v5, 0, *(float*)a2, *(float*)(a2 + 4));
				nox_xxx_scriptCallByEventBlock_502490((int*)(v20 + 64), (int)v5, v3, 2);
				v11 = nox_float2int(*(float*)(v3 + 56));
				v12 = *(float*)(v3 + 60);
				v22.field_0 = v11;
				v13 = nox_float2int(v12);
				v14 = *(float*)a2 - *(float*)(v3 + 56);
				v22.field_4 = v13 - 50;
				a1a.field_0 = v14;
				a1a.field_4 = *(float*)(a2 + 4) - *(float*)(v3 + 60);
				nox_xxx_utilNormalizeVector_509F20(&a1a);
				v18 = a1a.field_0 * 30.0 + *(float*)a2;
				v22.field_8 = nox_float2int(v18);
				v19 = a1a.field_4 * 30.0 + *(float*)(a2 + 4);
				v22.field_C = nox_float2int(v19);
				nox_xxx_sendGeneratorSpawnFX_523830(&v22, 10);
				nox_xxx_aud_501960(1002, (int)v5, 0, 0);
			} else {
				result = (uint32_t*)nox_xxx_objectFreeMem_4E38A0((int)v5);
			}
		}
	}
	return result;
}

//----- (0054F2B0) --------------------------------------------------------
void nox_xxx_unitCreatureCopyUC_54F2B0(int a1, int a2) {
	int v2;       // ebp
	int v3;       // edi
	uint32_t* v4; // eax
	uint32_t* v5; // esi
	int v6;       // eax
	short result; // ax

	v2 = a1;
	memcpy(*(void**)(a2 + 748), *(const void**)(a1 + 748), 0x898u);
	if (*(uint8_t*)(a1 + 12) & 0x10) {
		v3 = nox_xxx_inventoryGetFirst_4E7980(a1);
		if (v3) {
			do {
				v4 = nox_xxx_newObjectWithTypeInd_4E3450(*(unsigned short*)(v3 + 4));
				v5 = v4;
				if (v4) {
					if (v4[2] & 0x13001000) {
						nox_xxx_modifSetItemAttrs_4E4990((int)v4, *(int**)(v3 + 692));
					}
					nox_xxx_inventoryPutImpl_4F3070(a2, (int)v5, 0);
					if (*(uint32_t*)(v3 + 16) & 0x100) {
						v6 = v5[2];
						if (v6 & 0x1001000) {
							nox_xxx_NPCEquipWeapon_53A2C0(a2, (int)v5);
						} else if (v6 & 0x2000000) {
							nox_xxx_NPCEquipArmor_53E520(a2, v5);
						}
					}
				}
				v3 = nox_xxx_inventoryGetNext_4E7990(v3);
			} while (v3);
			v2 = a1;
		}
	}
	result = *(uint16_t*)(v2 + 124);
	*(uint16_t*)(a2 + 124) = result;
	*(uint16_t*)(a2 + 126) = result;
}

//----- (0054F740) --------------------------------------------------------
void nox_xxx_unitUpdateMover_54F740(int a1) {
	float* v1;         // edi
	unsigned char* v2; // esi
	int v3;            // eax
	int v4;            // ebp
	int v5;            // eax
	uint32_t* v6;      // ebp
	double v7;         // st7
	double v8;         // st7
	int v9;            // ecx
	char v10;          // al
	int v11;           // eax
	int v12;           // esi
	double v13;        // st7
	double v14;        // st6
	float v15;         // [esp+14h] [ebp-4h]
	float v16;         // [esp+1Ch] [ebp+4h]

	v1 = (float*)a1;
	v2 = *(unsigned char**)(a1 + 748);
	if (!*((uint32_t*)v2 + 8)) {
		nox_xxx_unitRemoveFromUpdatable_4DA920((uint32_t*)a1);
		return;
	}
	if (!*((uint32_t*)v2 + 7)) {
		v3 = nox_xxx_netGetUnitByExtent_4ED020(*((uint32_t*)v2 + 8));
		*((uint32_t*)v2 + 7) = v3;
		if (!v3) {
			nox_xxx_unitRemoveFromUpdatable_4DA920((uint32_t*)a1);
			return;
		}
	}
	if (*((uint32_t*)v2 + 4) && !*((uint32_t*)v2 + 3)) {
		*((uint32_t*)v2 + 3) = nox_server_getWaypointById_579C40(*((uint32_t*)v2 + 4));
	}
	if (*((uint32_t*)v2 + 6) && !*((uint32_t*)v2 + 5)) {
		*((uint32_t*)v2 + 5) = nox_server_getWaypointById_579C40(*((uint32_t*)v2 + 6));
	}
	v4 = *((uint32_t*)v2 + 7);
	v5 = *(uint32_t*)(v4 + 16);
	if (!(v5 & 4) || v5 & 0x20) {
		*((uint32_t*)v2 + 7) = 0;
		nox_xxx_unitRemoveFromUpdatable_4DA920((uint32_t*)a1);
		return;
	}
	switch (*v2) {
	case 0u:
		if (*(uint32_t*)(a1 + 16) & 0x1000000) {
			v6 = nox_server_getWaypointById_579C40(*((uint32_t*)v2 + 2));
			if (v6) {
				nox_xxx_unitMove_4E7010(a1, (float2*)(*((uint32_t*)v2 + 7) + 56));
				v7 = (double)*((int*)v2 + 1);
				*v2 = 1;
				*((uint32_t*)v2 + 3) = v6;
				*((uint32_t*)v2 + 5) = 0;
				v8 = v7 * 0.25;
				*(float*)(a1 + 548) = v8;
				*(float*)(a1 + 544) = v8;
			} else {
				*v2 = 3;
			}
		}
		break;
	case 1u:
		if (*(uint32_t*)(a1 + 16) & 0x1000000) {
			if (*(float*)(a1 + 80) != 0.0 && *(float*)(a1 + 84) != 0.0) {
				v9 = *((uint32_t*)v2 + 3);
				if ((*(float*)(v9 + 12) - *(float*)(a1 + 60)) * *(float*)(a1 + 84) +
						(*(float*)(v9 + 8) - *(float*)(a1 + 56)) * *(float*)(a1 + 80) <=
					0.0) {
					v10 = *(uint8_t*)(v9 + 476);
					if (v10) {
						if (v10 == 1) {
							*((uint32_t*)v2 + 5) = v9;
							*((uint32_t*)v2 + 3) = *(uint32_t*)(v9 + 92);
						} else {
							do {
								v11 = nox_common_randomInt_415FA0(0, *(unsigned char*)(v9 + 476) - 1);
								v9 = *((uint32_t*)v2 + 3);
							} while (*(uint32_t*)(v9 + 8 * v11 + 92) == *((uint32_t*)v2 + 5));
							*((uint32_t*)v2 + 5) = v9;
							*((uint32_t*)v2 + 3) = *(uint32_t*)(v9 + 8 * v11 + 92);
						}
					} else {
						*v2 = 3;
						nox_xxx_unitMove_4E7010(v4, (float2*)(v9 + 8));
					}
				}
			}
			if (*v2 == 1) {
				nox_xxx_unitMove_4E7010(*((uint32_t*)v2 + 7), (float2*)(a1 + 56));
				v12 = *((uint32_t*)v2 + 3);
				v13 = *(float*)(v12 + 8) - *(float*)(a1 + 56);
				v14 = *(float*)(v12 + 12) - *(float*)(a1 + 60);
				v15 = v14;
				v16 = sqrt(v14 * v15 + v13 * v13) + 0.1;
				v1[20] = v13 * (double)v1[136] / v16;
				if (v1[20] == 0.0) {
					v1[20] = FLT_MIN; // We need a minimal value so we don't end up checked as
									  // "uninitialized"
				}
				v1[21] = v15 * (double)v1[136] / v16;
				if (v1[21] == 0.0) {
					v1[21] = FLT_MIN; // We need a minimal value so we don't end up checked as
									  // "uninitialized"
				}
			}
		} else {
			*v2 = 2;
		}
		break;
	case 2u:
		if (*(uint32_t*)(a1 + 16) & 0x1000000) {
			nox_xxx_unitMove_4E7010(a1, (float2*)(v4 + 56));
			*v2 = 1;
		}
		break;
	case 3u:
		nox_xxx_unitRemoveFromUpdatable_4DA920((uint32_t*)a1);
		return;
	default:
		return;
	}
}

//----- (0054F9A0) --------------------------------------------------------
int nox_xxx_updateShootingTrap_54F9A0(int a1) {
	int result; // eax
	int* v2;    // esi
	int v3;     // eax
	int v4;     // eax

	result = *(uint32_t*)(a1 + 16);
	v2 = *(int**)(a1 + 748);
	if (result & 0x1000000) {
		if (!*((uint8_t*)v2 + 48)) {
			*v2 = 0;
			v2[1] = 0;
		}
		v3 = *v2;
		*((uint8_t*)v2 + 48) = 1;
		if (!v3) {
			sub_54FBB0(a1);
			*v2 = gameFPS();
		}
		if (*((uint8_t*)v2 + 8) == 1 && !v2[1]) {
			if (!*getMemU32Ptr(0x5D4594, 2491780)) {
				*getMemU32Ptr(0x5D4594, 2491780) = nox_xxx_getNameId_4E3AA0("ArrowTrap1");
				*getMemU32Ptr(0x5D4594, 2491784) = nox_xxx_getNameId_4E3AA0("ArrowTrap2");
			}
			nox_xxx_createArrowTrapProjectile_54FA80(a1, v2[3]);
			v4 = *(unsigned short*)(a1 + 4);
			if ((unsigned short)v4 == *getMemU32Ptr(0x5D4594, 2491780)) {
				nox_xxx_sendArrowTrapFX_5238A0((float*)(a1 + 56), 1);
			} else if (v4 == *getMemU32Ptr(0x5D4594, 2491784)) {
				nox_xxx_sendArrowTrapFX_5238A0((float*)(a1 + 56), 2);
			}
			v2[1] = 30;
		}
		if (*v2) {
			--*v2;
		}
		result = v2[1];
		if (result) {
			v2[1] = --result;
		}
	} else {
		*((uint8_t*)v2 + 48) = 0;
	}
	return result;
}

//----- (0054FA80) --------------------------------------------------------
void nox_xxx_createArrowTrapProjectile_54FA80(int a1, int a2) {
	double v2; // st7
	int v3;    // eax
	int v4;    // eax
	int v5;    // esi
	short v6;  // ax
	int* v7;   // esi
	float v9;  // [esp+0h] [ebp-18h]
	float v10; // [esp+10h] [ebp-8h]
	float v11; // [esp+14h] [ebp-4h]

	v2 = *(float*)(a1 + 176) + 4.0;
	v3 = 8 * *(short*)(a1 + 124);
	v10 = v2 * *getMemFloatPtr(0x587000, 194136 + v3) + *(float*)(a1 + 56);
	v11 = v2 * *getMemFloatPtr(0x587000, 194140 + v3) + *(float*)(a1 + 60);
	v4 = (int)nox_xxx_newObjectWithTypeInd_4E3450(a2);
	v5 = v4;
	if (v4) {
		nox_xxx_createAt_4DAA50(v4, a1, v10, v11);
		v6 = *(uint16_t*)(a1 + 124);
		*(uint16_t*)(v5 + 124) = v6;
		*(uint16_t*)(v5 + 126) = v6;
		*(float*)(v5 + 80) = *getMemFloatPtr(0x587000, 194136 + 8 * *(short*)(a1 + 124)) * *(float*)(v5 + 544);
		*(float*)(v5 + 84) = *getMemFloatPtr(0x587000, 194140 + 8 * *(short*)(a1 + 124)) * *(float*)(v5 + 544);
		if (!*getMemU32Ptr(0x5D4594, 2491768)) {
			*getMemU32Ptr(0x5D4594, 2491768) = nox_xxx_getNameId_4E3AA0("MercArcherArrow");
			*getMemU32Ptr(0x5D4594, 2491772) = nox_xxx_getNameId_4E3AA0("ArrowTrap1");
			*getMemU32Ptr(0x5D4594, 2491776) = nox_xxx_getNameId_4E3AA0("ArrowTrap2");
		}
		v4 = *(unsigned short*)(a1 + 4);
		if ((unsigned short)v4 == *getMemU32Ptr(0x5D4594, 2491772) || v4 == *getMemU32Ptr(0x5D4594, 2491776)) {
			v7 = *(int**)(v5 + 700);
			v9 = nox_xxx_gamedataGetFloat_419D40("ArrowTrapDamage");
			v4 = nox_float2int(v9);
			*v7 = v4;
			v7[1] = v4;
		}
		if (a2 == *getMemU32Ptr(0x5D4594, 2491768)) {
			nox_xxx_aud_501960(889, a1, 0, 0);
		}
	}
}

//----- (0054FBB0) --------------------------------------------------------
void sub_54FBB0(int a1) {
	int v1 = *(uint32_t*)(a1 + 748);
	if (sub_54FBF0(a1)) {
		if (*(uint8_t*)(v1 + 8) != 1) {
			*(uint8_t*)(v1 + 8) = 1;
			*(uint32_t*)(v1 + 4) = 0;
		}
	} else {
		*(uint8_t*)(v1 + 8) = 0;
	}
}

//----- (0054FBF0) --------------------------------------------------------
int sub_54FBF0(int a3) {
	double v1; // st7
	float4 a1; // [esp+0h] [ebp-10h]

	a1.field_0 = *(float*)(a3 + 56) - 350.0;
	a1.field_4 = *(float*)(a3 + 60) - 350.0;
	a1.field_8 = *(float*)(a3 + 56) + 350.0;
	v1 = *(float*)(a3 + 60) + 350.0;
	*getMemU32Ptr(0x5D4594, 2491764) = 0;
	a1.field_C = v1;
	nox_xxx_getUnitsInRect_517C10(&a1, nox_xxx_unitIsAttackReachable_54FC50, a3);
	return *getMemU32Ptr(0x5D4594, 2491764);
}

//----- (0054FC50) --------------------------------------------------------
void nox_xxx_unitIsAttackReachable_54FC50(int a1, int a2) {
	if (*(uint8_t*)(a1 + 8) & 6 && !(*(uint32_t*)(a1 + 16) & 0x8020) && nox_xxx_unitIsEnemyTo_5330C0(a2, a1) &&
		nox_xxx_mapCheck_537110(a1, a2)) {
		if (nox_server_testTwoPointsAndDirection_4E6E50((float2*)(a2 + 56), *(short*)(a2 + 124), (float2*)a1 + 7) & 1) {
			*getMemU32Ptr(0x5D4594, 2491764) = 1;
		}
	}
}

//----- (0054FCD0) --------------------------------------------------------
void nox_xxx_collideTrigger_54FCD0(int a1, int a2) {
	int* v2; // esi
	int v3;  // eax
	int v4;  // eax
	int v5;  // eax
	char v6; // al
	char v7; // al
	int v8;  // eax

	v2 = *(int**)(a1 + 748);
	if (*(uint32_t*)(a1 + 16) & 0x1000000) {
		if (*((uint8_t*)v2 + 8) != 5) {
			if (a2) {
				if (nox_xxx_objectGetMass_4E4A70(a2) > 0.0) {
					v4 = v2[11];
					if (!v4 || v4 & *(uint32_t*)(a2 + 8)) {
						v5 = v2[12];
						if (!v5 || !(v5 & *(uint32_t*)(a2 + 8))) {
							v6 = *((uint8_t*)v2 + 52);
							if (!v6 || *(uint8_t*)(a2 + 52) == v6) {
								v7 = *((uint8_t*)v2 + 53);
								if ((!v7 || *(char*)(a2 + 52) != v7) &&
									(v2[4] == -1 ||
									 *(uint32_t*)nox_xxx_scriptCallByEventBlock_502490(v2 + 3, a2, a1, 1))) {
									v8 = *v2;
									v2[1] = a2;
									LOBYTE(v8) = v8 | 1;
									*v2 = v8;
								}
							}
						}
					}
				}
			}
		}
	} else {
		v3 = *v2;
		v2[1] = 0;
		LOBYTE(v3) = v3 & 0xFE;
		*v2 = v3;
	}
}

//----- (0054FD80) --------------------------------------------------------
float* nox_xxx_createSpark_54FD80(float a1, float a2, int a3, int a4, float a5, float a6, float a7, int a8) {
	float* v8;        // esi
	uint32_t* v9;     // ebx
	int* v10;         // edi
	int v11;          // eax
	int v12;          // eax
	unsigned int v13; // ecx
	float* result;    // eax

	v8 = (float*)nox_xxx_newObjectByTypeID_4E3810("Spark");
	if (!v8) {
		return 0;
	}
	v9 = (uint32_t*)*((uint32_t*)v8 + 175);
	v10 = (int*)*((uint32_t*)v8 + 187);
	nox_xxx_createAt_4DAA50((int)v8, a8, a1, a2);
	*(uint32_t*)&v8[34] = gameFrame();
	*v10 = a4;
	v10[1] = a4;
	v10[3] = a3;
	v11 = *((uint32_t*)v8 + 2);
	BYTE1(v11) &= 0xDFu;
	v12 = v11 | 0x80000;
	v13 = (uint32_t)v8[4] & 0xFF7FFFBF;
	*((uint32_t*)v8 + 2) = v12;
	*((uint32_t*)v8 + 4) = v13;
	switch (a3) {
	case 0:
		*((uint32_t*)v8 + 4) = v13 | 0x40;
		*((uint32_t*)v8 + 2) = v12 & 0xFFF7FFFF;
		break;
	case 1:
		BYTE1(v12) |= 0x20u;
		*((uint32_t*)v8 + 2) = v12;
		*((uint32_t*)v8 + 4) = 0x800000 | v13;
		*v9 = 3;
		nox_xxx_unitRaise_4E46F0((int)v8, 28.0);
		v8[27] = a7;
		v8[29] = 7.0;
		v8[20] = a5;
		v8[21] = a6;
		return v8;
	case 2:
		*((uint32_t*)v8 + 4) = 0x800040 | v13;
		*v9 = 0;
		nox_xxx_unitRaise_4E46F0((int)v8, 28.0);
		v8[27] = a7;
		v8[29] = 7.0;
		v8[20] = a5;
		v8[21] = a6;
		return v8;
	case 4:
		*((uint32_t*)v8 + 2) = v12 & 0xFFF7FFFF;
		break;
	default:
		break;
	}
	*v9 = 0;
	nox_xxx_unitRaise_4E46F0((int)v8, 28.0);
	v8[27] = a7;
	v8[29] = 0.0;
	v8[20] = a5;
	v8[21] = a6;
	result = v8;
	return result;
}

//----- (0054FEF0) --------------------------------------------------------
void sub_54FEF0(int a2) {
	float* v1; // edi
	int v2;    // esi
	int v3;    // ebx
	int v4;    // eax
	int i;     // ebp
	int v6;    // esi
	float v7;  // [esp+0h] [ebp-20h]
	float v8;  // [esp+0h] [ebp-20h]
	float v9;  // [esp+0h] [ebp-20h]
	float v10; // [esp+0h] [ebp-20h]
	int v11;   // [esp+14h] [ebp-Ch]
	int2 a1;   // [esp+18h] [ebp-8h]
	int a2a;   // [esp+24h] [ebp+4h]

	v1 = (float*)a2;
	if (!(*(uint32_t*)(a2 + 8) & 0x400000)) {
		v7 = *(float*)(a2 + 232) * 0.043478262;
		a2a = nox_float2int(v7);
		v8 = v1[59] * 0.043478262;
		v2 = nox_float2int(v8);
		v9 = v1[60] * 0.043478262;
		v3 = nox_float2int(v9);
		v10 = v1[61] * 0.043478262;
		v4 = nox_float2int(v10);
		v11 = v4;
		for (i = v2; i <= v4; ++i) {
			v6 = a2a;
			a1.field_4 = i;
			if (a2a <= v3) {
				do {
					a1.field_0 = v6;
					if (sub_54FFC0(&a1, (int)v1)) {
						sub_548100(&a1, (int)v1);
					}
					++v6;
				} while (v6 <= v3);
				v4 = v11;
			}
		}
	}
}

//----- (0054FFC0) --------------------------------------------------------
int sub_54FFC0(int2* a1, int a2) {
	int v2;           // ebx
	int v3;           // eax
	int v4;           // edi
	int v6;           // eax
	int v7;           // ecx
	unsigned char v8; // dl
	int v9;           // eax
	int v10;          // ebx
	int v11;          // eax
	int v12;          // esi
	char v13;         // [esp+Ch] [ebp-20h]
	int a5;           // [esp+10h] [ebp-1Ch]
	int a4;           // [esp+14h] [ebp-18h]
	int v16;          // [esp+18h] [ebp-14h]
	float2 a1a;       // [esp+1Ch] [ebp-10h]
	float2 a7;        // [esp+24h] [ebp-8h]

	v2 = a2;
	a4 = 0;
	a5 = 0;
	v16 = 0;
	v3 = *(uint32_t*)(a2 + 16);
	if (!(v3 & 0x4000) || *(uint32_t*)(a2 + 172) != 2 || (v13 = 0, *(float*)(a2 + 176) > 9.0)) {
		v13 = 64;
	}
	v4 = (unsigned char)sub_57B500(a1->field_0, a1->field_4, v13);
	if (v4 == 255) {
		return 0;
	}
	v6 = a1->field_0;
	v7 = a1->field_4;
	a1a.field_0 = (double)(23 * a1->field_0);
	v8 = table_313272[v4].f0;
	a1a.field_4 = (double)(23 * v7);
	if (v8) {
		if (getMemByte(0x587000, 292496 + v4) & 2 && sub_57B500(v6 - 1, v7 - 1, v13) == -1) {
			a4 = 1;
		}
		if (getMemByte(0x587000, 292496 + v4) & 4 && sub_57B500(a1->field_0 + 1, a1->field_4 + 1, v13) == -1) {
			a5 = 1;
		}
		if (v4 == 7 || v4 == 10) {
			a5 = 1;
		} else if (v4 == 8 || v4 == 9) {
			a4 = 1;
			v9 = sub_550280((int)&a1a, table_313272[v4].f4, table_313272[v4].f8, 1, a5, a2 + 64, (int)&a7);
			goto LABEL_21;
		}
		v9 = sub_550280((int)&a1a, table_313272[v4].f4, table_313272[v4].f8, a4, a5, a2 + 64, (int)&a7);
	LABEL_21:
		if (v9) {
			if (*(uint8_t*)(a2 + 8) & 4) {
				v10 = *(uint32_t*)(a2 + 748);
				if (v10) {
					*(uint32_t*)(v10 + 296) = nox_server_getWallAtGrid_410580(a1->field_0, a1->field_4);
				}
				v2 = a2;
			}
			if (sub_550380(v4, v2, &a7)) {
				v16 = 1;
			}
		}
	}
	if (table_313272[v4].f12) {
		if (getMemByte(0x587000, 292496 + v4) & 8 && sub_57B500(a1->field_0 - 1, a1->field_4 + 1, v13) == -1) {
			a4 = 1;
		}
		if (getMemByte(0x587000, 292496 + v4) & 1 && sub_57B500(a1->field_0 + 1, a1->field_4 - 1, v13) == -1) {
			a5 = 1;
		}
		if (v4 == 7 || v4 == 8) {
			a4 = 1;
		} else if (v4 == 9 || v4 == 10) {
			v11 = sub_5502F0(&a1a, table_313272[v4].f16, table_313272[v4].f20, a4, 1, (float2*)(v2 + 64), &a7);
			goto LABEL_42;
		}
		v11 = sub_5502F0(&a1a, table_313272[v4].f16, table_313272[v4].f20, a4, a5, (float2*)(v2 + 64), &a7);
	LABEL_42:
		if (v11) {
			if (*(uint8_t*)(v2 + 8) & 4) {
				v12 = *(uint32_t*)(v2 + 748);
				if (v12) {
					*(uint32_t*)(v12 + 296) = nox_server_getWallAtGrid_410580(a1->field_0, a1->field_4);
				}
			}
			if (sub_550380(v4, v2, &a7)) {
				v16 = 1;
			}
		}
	}
	return v16;
}

//----- (00550280) --------------------------------------------------------
int sub_550280(int a1, float a2, float a3, int a4, int a5, int a6, int a7) {
	double v7; // st7

	v7 = (*(float*)a6 - *(float*)a1 + *(float*)(a6 + 4) - *(float*)(a1 + 4)) * 0.70709997 * 0.70709997;
	if (v7 < a2) {
		if (!a4) {
			return 0;
		}
		v7 = a2;
	}
	if (v7 > a3) {
		if (!a5) {
			return 0;
		}
		v7 = a3;
	}
	*(float*)a7 = v7 + *(float*)a1;
	*(float*)(a7 + 4) = v7 + *(float*)(a1 + 4);
	return 1;
}

//----- (005502F0) --------------------------------------------------------
int sub_5502F0(float2* a1, float a2, float a3, int a4, int a5, float2* a6, float2* a7) {
	double v7; // st7

	v7 = 23.0 -
		 ((a6->field_4 - a1->field_4) * 0.70709997 - (a6->field_0 - (a1->field_0 + 23.0)) * 0.70709997) * 0.70709997;
	if (v7 < a2) {
		if (!a4) {
			return 0;
		}
		v7 = a2;
	}
	if (v7 > a3) {
		if (!a5) {
			return 0;
		}
		v7 = a3;
	}
	a7->field_0 = v7 + a1->field_0;
	a7->field_4 = 23.0 - v7 + a1->field_4;
	return 1;
}

//----- (00550380) --------------------------------------------------------
int sub_550380(int a1, int a2, float2* a3) {
	float* v3;      // esi
	double v4;      // st7
	double v5;      // st6
	long double v6; // st6
	double v7;      // st7
	double v8;      // st7
	double v9;      // st7
	float v11;      // [esp+0h] [ebp-14h]
	float v12;      // [esp+4h] [ebp-10h]
	float2 v13;     // [esp+Ch] [ebp-8h]
	float v14;      // [esp+1Ch] [ebp+8h]
	float v15;      // [esp+20h] [ebp+Ch]
	float v16;      // [esp+20h] [ebp+Ch]

	v3 = (float*)a2;
	v4 = *(float*)(a2 + 64) - a3->field_0;
	v5 = *(float*)(a2 + 68) - a3->field_4;
	v15 = v5;
	v6 = sqrt(v5 * v15 + v4 * v4);
	v14 = v6;
	if (v6 == 0.0) {
		v14 = 0.0099999998;
	}
	if (v14 >= (double)v3[44]) {
		return 0;
	}
	v13.field_0 = v4 / v14;
	v7 = v15 / v14;
	v13.field_4 = v7;
	v8 = -v7 * v3[21] + -v13.field_0 * v3[20];
	v16 = v8;
	if (v8 > 0.0 && !sub_550480((int)v3)) {
		v3[20] = v16 * v13.field_0 + v3[20];
		v3[21] = v16 * v13.field_4 + v3[21];
	}
	v9 = (v3[44] - v14) * *(float*)&dword_587000_292492;
	v12 = v9 * v13.field_4;
	v11 = v9 * v13.field_0;
	sub_548600((int)v3, v11, v12);
	nox_xxx_collSysAddCollision_548630((int)v3, 0, &v13);
	return 1;
}

//----- (00550480) --------------------------------------------------------
int sub_550480(int a1) {
	int v1; // eax

	v1 = *getMemU32Ptr(0x5D4594, 2491788);
	if (!*getMemU32Ptr(0x5D4594, 2491788)) {
		v1 = nox_xxx_getNameId_4E3AA0("GameBall");
		*getMemU32Ptr(0x5D4594, 2491788) = v1;
	}
	return *(unsigned short*)(a1 + 4) == v1;
}

//----- (005504B0) --------------------------------------------------------
void sub_5504B0(int a2) {
	float* v1; // edi
	int v2;    // esi
	int v3;    // ebx
	int v4;    // eax
	int i;     // ebp
	int v6;    // esi
	float v7;  // [esp+0h] [ebp-20h]
	float v8;  // [esp+0h] [ebp-20h]
	float v9;  // [esp+0h] [ebp-20h]
	float v10; // [esp+0h] [ebp-20h]
	int v11;   // [esp+14h] [ebp-Ch]
	int2 a1;   // [esp+18h] [ebp-8h]
	int a2a;   // [esp+24h] [ebp+4h]

	v1 = (float*)a2;
	if (!(*(uint32_t*)(a2 + 8) & 0x400000)) {
		v7 = *(float*)(a2 + 232) * 0.043478262;
		a2a = nox_float2int(v7);
		v8 = v1[59] * 0.043478262;
		v2 = nox_float2int(v8);
		v9 = v1[60] * 0.043478262;
		v3 = nox_float2int(v9);
		v10 = v1[61] * 0.043478262;
		v4 = nox_float2int(v10);
		v11 = v4;
		for (i = v2; i <= v4; ++i) {
			v6 = a2a;
			a1.field_4 = i;
			if (a2a <= v3) {
				do {
					a1.field_0 = v6;
					if (sub_550580(&a1, v1)) {
						sub_548100(&a1, (int)v1);
					}
					++v6;
				} while (v6 <= v3);
				v4 = v11;
			}
		}
	}
}

//----- (00550580) --------------------------------------------------------
int sub_550580(int2* a1, float* a2) {
	int v2;            // ebp
	int v3;            // ecx
	double v5;         // st7
	double v6;         // st6
	unsigned char* v7; // edi
	double v8;         // st5
	char v9;           // bl
	double v10;        // st7
	float2 a5;         // [esp+8h] [ebp-30h]
	float2 a1a;        // [esp+10h] [ebp-28h]
	float2 a2a;        // [esp+18h] [ebp-20h]
	float2 a3;         // [esp+20h] [ebp-18h]
	float4 a4;         // [esp+28h] [ebp-10h]
	float v16;         // [esp+3Ch] [ebp+4h]

	v2 = 0;
	v3 = (unsigned char)sub_57B500(a1->field_0, a1->field_4, 64);
	if (v3 == 255) {
		return 0;
	}
	v5 = (double)(23 * a1->field_0);
	v6 = (double)(23 * a1->field_4);
	v7 = getMemAt(0x587000, 292520 + 24 * v3);
	a1a.field_0 = (a2[17] + a2[16]) * 0.70710677;
	a1a.field_4 = (a2[17] - a2[16]) * 0.70710677;
	a3.field_0 = (a2[19] + a2[18]) * 0.70710677;
	a3.field_4 = (a2[19] - a2[18]) * 0.70710677;
	v8 = a2[46] * 0.5;
	a4.field_0 = a1a.field_0 - v8;
	v16 = a2[47] * 0.5;
	a4.field_4 = a1a.field_4 - v16;
	a4.field_8 = v8 + a1a.field_0;
	a4.field_C = v16 + a1a.field_4;
	a2a.field_0 = (v6 + v5) * 0.70710677;
	a2a.field_4 = (v6 - v5) * 0.70710677;
	v9 = sub_550CB0(&a1a, &a2a);
	if (*v7) {
		if (!((unsigned char)v9 & v7[1])) {
			a5.field_0 = a2a.field_0 + 16.263456;
			v10 = a2a.field_4 - 16.263456;
			a5.field_4 = v10;
			a5.field_4 = v10 + *((float*)v7 + 1);
			if (sub_550A10((int)a2, &a1a, &a3, &a4, &a5, *((float*)v7 + 2))) {
				v2 = 1;
			}
		}
	}
	if (v7[12] && !((unsigned char)v9 & v7[13])) {
		a5 = a2a;
		a5.field_0 = a2a.field_0 + *((float*)v7 + 4);
		if (sub_550760((int)a2, &a1a, &a3, &a4, &a5, *((float*)v7 + 5))) {
			v2 = 1;
		}
	}
	return v2;
}

//----- (00550760) --------------------------------------------------------
int sub_550760(int a1, float2* a2, float2* a3, float4* a4, float2* a5, float a6) {
	float2* v6;      // edx
	float4* v7;      // edi
	double v8;       // st7
	double v9;       // st7
	double v10;      // st6
	double v11;      // st7
	double v12;      // st7
	float v13;       // et1
	double v14;      // st7
	float v15;       // et1
	double v16;      // st7
	float v17;       // et1
	long double v18; // st7
	double v19;      // st7
	double v20;      // st6
	double v21;      // st7
	float2 v23;      // [esp+8h] [ebp-18h]
	float2 v24;      // [esp+10h] [ebp-10h]
	float2 v25;      // [esp+18h] [ebp-8h]
	float v26;       // [esp+28h] [ebp+8h]
	float v27;       // [esp+30h] [ebp+10h]
	float v28;       // [esp+30h] [ebp+10h]
	float v29;       // [esp+30h] [ebp+10h]
	float v30;       // [esp+34h] [ebp+14h]
	float v31;       // [esp+34h] [ebp+14h]
	float v32;       // [esp+38h] [ebp+18h]
	float v33;       // [esp+38h] [ebp+18h]

	v6 = a5;
	v7 = a4;
	if (a4->field_0 <= (double)a5->field_0) {
		v30 = a5->field_0;
	} else {
		v30 = a4->field_0;
	}
	v8 = a6 + v6->field_0;
	if (a4->field_8 >= v8) {
		v27 = v8;
	} else {
		v27 = a4->field_8;
	}
	if (v30 > v8 || v27 < (double)v6->field_0) {
		return 0;
	}
	v9 = a2->field_4 - v6->field_4;
	v10 = a3->field_4 - v6->field_4;
	if (v10 * v9 < 0.0) {
		v11 = v6->field_4;
		v32 = v10;
		if (v32 >= 0.0) {
			v12 = v11 + 2.0;
		} else {
			v12 = v11 - 2.0;
		}
		a2->field_4 = v12;
		*(float*)(a1 + 64) = (a2->field_0 - a2->field_4) * 0.70710677;
		*(float*)(a1 + 68) = (a2->field_0 + a2->field_4) * 0.70710677;
		v7->field_4 = a2->field_4 - *(float*)(a1 + 188) * 0.5;
		v7->field_C = *(float*)(a1 + 188) * 0.5 + a2->field_4;
		v9 = a2->field_4 - v6->field_4;
	}
	if (v7->field_4 > (double)v6->field_4 || v7->field_C < (double)v6->field_4) {
		return 0;
	}
	v26 = (v27 - v30) / (v7->field_8 - v7->field_0);
	v31 = (*(float*)(a1 + 84) - *(float*)(a1 + 80)) * 0.70710677;
	if (v9 >= 0.0) {
		v28 = v6->field_4 - v7->field_4;
		v15 = *(float*)&dword_587000_292492;
		v14 = v15 * v28;
	} else {
		v28 = v7->field_C - v6->field_4;
		v13 = *(float*)&dword_587000_292492;
		v14 = -(v13 * v28);
	}
	v33 = v14;
	v16 = nox_xxx_objectGetMass_4E4A70(a1);
	v17 = *(float*)&dword_587000_292492;
	v18 = (sqrt(v16 * v17 * 4.0) * -v31 * 0.5 + v33) * v26;
	v23.field_0 = v18 * -0.70710677;
	v23.field_4 = v18 * 0.70710677;
	if (v28 >= 0.0) {
		v19 = 0.70710677;
		v24.field_0 = -0.70710677;
	} else {
		v19 = -0.70710677;
		v24.field_0 = 0.70710677;
	}
	v20 = v24.field_0 * *(float*)(a1 + 80) + v19 * *(float*)(a1 + 84);
	if (v20 < 0.0) {
		*(float*)(a1 + 80) = *(float*)(a1 + 80) - v20 * v24.field_0;
		*(float*)(a1 + 84) = *(float*)(a1 + 84) - v20 * v19;
	}
	v21 = -v19;
	v25.field_0 = v21;
	v29 = v21 * *(float*)(a1 + 80) + v24.field_0 * *(float*)(a1 + 84);
	v23.field_0 = v23.field_0 - nox_xxx_objectGetMass_4E4A70(a1) * v29 * v25.field_0 * 0.69999999;
	v23.field_4 = v23.field_4 - nox_xxx_objectGetMass_4E4A70(a1) * v29 * v24.field_0 * 0.69999999;
	sub_548600(a1, v23.field_0, v23.field_4);
	nox_xxx_collSysAddCollision_548630(a1, 0, &v23);
	return 1;
}

//----- (00550A10) --------------------------------------------------------
int sub_550A10(int a1, float2* a2, float2* a3, float4* a4, float2* a5, float a6) {
	float2* v6;      // edx
	float4* v7;      // edi
	double v8;       // st7
	double v9;       // st7
	double v10;      // st7
	double v11;      // st6
	double v12;      // st7
	double v13;      // st7
	float v14;       // et1
	double v15;      // st7
	float v16;       // et1
	double v17;      // st7
	float v18;       // et1
	long double v19; // st7
	double v20;      // st7
	double v21;      // st6
	double v22;      // st7
	float2 v24;      // [esp+8h] [ebp-18h]
	float2 v25;      // [esp+10h] [ebp-10h]
	float2 v26;      // [esp+18h] [ebp-8h]
	float v27;       // [esp+28h] [ebp+8h]
	float v28;       // [esp+30h] [ebp+10h]
	float v29;       // [esp+30h] [ebp+10h]
	float v30;       // [esp+30h] [ebp+10h]
	float v31;       // [esp+34h] [ebp+14h]
	float v32;       // [esp+34h] [ebp+14h]
	float v33;       // [esp+38h] [ebp+18h]
	float v34;       // [esp+38h] [ebp+18h]

	v6 = a5;
	v7 = a4;
	if (a5->field_4 >= (double)a4->field_4) {
		v8 = a5->field_4;
	} else {
		v8 = a4->field_4;
	}
	v31 = v8;
	v9 = a6 + v6->field_4;
	if (a4->field_C >= v9) {
		v28 = v9;
	} else {
		v28 = a4->field_C;
	}
	if (v31 > v9 || v28 < (double)v6->field_4) {
		return 0;
	}
	v10 = a2->field_0 - v6->field_0;
	v11 = a3->field_0 - v6->field_0;
	if (v11 * v10 < 0.0) {
		v12 = v6->field_0;
		v33 = v11;
		if (v33 >= 0.0) {
			v13 = v12 + 2.0;
		} else {
			v13 = v12 - 2.0;
		}
		a2->field_0 = v13;
		*(float*)(a1 + 64) = (a2->field_0 - a2->field_4) * 0.70710677;
		*(float*)(a1 + 68) = (a2->field_4 + a2->field_0) * 0.70710677;
		v7->field_0 = a2->field_0 - *(float*)(a1 + 184) * 0.5;
		v7->field_8 = *(float*)(a1 + 184) * 0.5 + a2->field_0;
		v10 = a2->field_0 - v6->field_0;
	}
	if (v7->field_0 > (double)v6->field_0 || v7->field_8 < (double)v6->field_0) {
		return 0;
	}
	v27 = (v28 - v31) / (v7->field_C - v7->field_4);
	v32 = (*(float*)(a1 + 80) + *(float*)(a1 + 84)) * 0.70710677;
	if (v10 >= 0.0) {
		v29 = v6->field_0 - v7->field_0;
		v16 = *(float*)&dword_587000_292492;
		v15 = v16 * v29;
	} else {
		v29 = v7->field_8 - v6->field_0;
		v14 = *(float*)&dword_587000_292492;
		v15 = -(v14 * v29);
	}
	v34 = v15;
	v17 = nox_xxx_objectGetMass_4E4A70(a1);
	v18 = *(float*)&dword_587000_292492;
	v19 = (sqrt(v17 * v18 * 4.0) * -v32 * 0.5 + v34) * v27 * 0.70710677;
	v24.field_0 = v19;
	v24.field_4 = v19;
	if (v29 >= 0.0) {
		v20 = 0.70710677;
		v25.field_0 = 0.70710677;
	} else {
		v20 = -0.70710677;
		v25.field_0 = -0.70710677;
	}
	v21 = v25.field_0 * *(float*)(a1 + 80) + v20 * *(float*)(a1 + 84);
	if (v21 < 0.0) {
		*(float*)(a1 + 80) = *(float*)(a1 + 80) - v21 * v25.field_0;
		*(float*)(a1 + 84) = *(float*)(a1 + 84) - v21 * v20;
	}
	v22 = -v20;
	v26.field_0 = v22;
	v30 = v22 * *(float*)(a1 + 80) + v25.field_0 * *(float*)(a1 + 84);
	v24.field_0 = v24.field_0 - nox_xxx_objectGetMass_4E4A70(a1) * v30 * v26.field_0 * 0.69999999;
	v24.field_4 = v24.field_4 - nox_xxx_objectGetMass_4E4A70(a1) * v30 * v25.field_0 * 0.69999999;
	sub_548600(a1, v24.field_0, v24.field_4);
	nox_xxx_collSysAddCollision_548630(a1, 0, &v24);
	return 1;
}

//----- (00550CB0) --------------------------------------------------------
char sub_550CB0(float2* a1, float2* a2) {
	double v2;   // st7
	char v3 = 0; // fps^1
	bool v4;     // c0
	char v5;     // c2
	bool v6;     // c3
	char v7;     // ah
	bool v8;     // c0
	bool v9;     // c3
	char result; // al
	float v11;   // [esp+0h] [ebp-8h]

	v11 = a1->field_0 - a2->field_0;
	v2 = a1->field_4 - a2->field_4;
	v4 = v11 < 16.263456;
	v5 = 0;
	v6 = v11 == 16.263456;
	v7 = v3;
	v8 = v2 < 0.0;
	v9 = v2 == 0.0;
	if (v7 & 0x41) {
		result = 8;
		if (v8 || v9) {
			result = 1;
		}
	} else if (v8 || v9) {
		result = 2;
	} else {
		result = 4;
	}
	return result;
}

//----- (00550D00) --------------------------------------------------------
void nox_xxx_collisionCheckCircleCircle_550D00(int a1, int a2) {
	int v2;         // esi
	int v3;         // edi
	int v4;         // eax
	long double v5; // st7
	double v6;      // st7
	int v7;         // ebp
	float v8;       // eax
	float v9;       // edx
	float v10;      // eax
	double v11;     // st7
	double v12;     // st6
	double v13;     // st6
	double v14;     // st5
	float v15;      // [esp+Ch] [ebp-20h]
	float v16;      // [esp+Ch] [ebp-20h]
	float v17;      // [esp+10h] [ebp-1Ch]
	float2 v18;     // [esp+14h] [ebp-18h]
	float4 a1a;     // [esp+1Ch] [ebp-10h]
	float v20;      // [esp+30h] [ebp+4h]
	float v21;      // [esp+30h] [ebp+4h]
	float v22;      // [esp+34h] [ebp+8h]
	float v23;      // [esp+34h] [ebp+8h]
	float v24;      // [esp+34h] [ebp+8h]

	v2 = a1;
	v3 = a2;
	v18.field_0 = *(float*)(a2 + 64) - *(float*)(a1 + 64);
	v18.field_4 = *(float*)(a2 + 68) - *(float*)(a1 + 68);
	if (v18.field_0 == 0.0 && v18.field_4 == 0.0) {
		v4 = nox_common_randomInt_415FA0(0, 3);
		v18.field_0 = *getMemFloatPtr(0x587000, 292784 + 8 * v4);
		v18.field_4 = *getMemFloatPtr(0x587000, 292788 + 8 * v4);
	}
	v5 = sqrt(v18.field_4 * v18.field_4 + v18.field_0 * v18.field_0);
	v22 = v5;
	if (v5 == 0.0) {
		v22 = 0.0099999998;
	}
	v6 = *(float*)(a1 + 176) + *(float*)(v3 + 176) - v22;
	v20 = v6;
	if (v6 > 0.0) {
		v7 = 1;
		if (!(*(uint32_t*)(v2 + 8) & 0x2204) || !(*(uint32_t*)(v3 + 8) & 0x2204) ||
			(v8 = *(float*)(v2 + 56), v9 = *(float*)(v3 + 56), a1a.field_4 = *(float*)(v2 + 60), a1a.field_0 = v8,
			 v10 = *(float*)(v3 + 60), a1a.field_8 = v9, a1a.field_C = v10,
			 nox_xxx_mapTraceRay_535250(&a1a, 0, 0, 0))) {
			nox_xxx_collSysAddCollision_548630(v3, v2, &v18);
			if ((*(uint8_t*)(v2 + 16) & 8) == 8 || (*(uint8_t*)(v3 + 16) & 8) == 8) {
				v7 = 0;
			}
			if ((!(*(uint8_t*)(v2 + 8) & 6) || (*(uint32_t*)(v3 + 16) & 0x2000) != 0x2000) && v7) {
				a1a.field_0 = v18.field_0 / v22;
				a1a.field_4 = v18.field_4 / v22;
				v15 = *(float*)(v2 + 80) - *(float*)(v3 + 80);
				v17 = *(float*)(v2 + 84) - *(float*)(v3 + 84);
				v23 = nox_xxx_objectGetMass_4E4A70(v2);
				if (nox_xxx_objectGetMass_4E4A70(v3) <= v23) {
					v11 = nox_xxx_objectGetMass_4E4A70(v3);
				} else {
					v11 = nox_xxx_objectGetMass_4E4A70(v2);
				}
				v12 = *(float*)&dword_587000_292488 * v20;
				v21 = -(v12 * a1a.field_0);
				v24 = -(v12 * a1a.field_4);
				if (!(*(uint8_t*)(v2 + 8) & 6) || !(*(uint8_t*)(v3 + 8) & 6)) {
					v13 = -a1a.field_4;
					v14 = v13 * v15 + v17 * a1a.field_0;
					v16 = v14;
					v21 = v21 - v14 * v13 * v11 * 0.69999999;
					v24 = v24 - v16 * v11 * a1a.field_0 * 0.69999999;
				}
				sub_548600(v2, v21, v24);
			}
			if (*(uint32_t*)(v2 + 16) & 0x8000000) {
				nox_xxx_unitHasCollideOrUpdateFn_537610(v2);
				*(uint32_t*)(v2 + 16) &= 0xF7FFFFFF;
			}
			if (*(uint32_t*)(v3 + 16) & 0x8000000) {
				nox_xxx_unitHasCollideOrUpdateFn_537610(v3);
				*(uint32_t*)(v3 + 16) &= 0xF7FFFFFF;
			}
		}
	}
}

//----- (00550F80) --------------------------------------------------------
void sub_550F80(float* a1, int a2) {
	float* v2;  // esi
	int v3;     // ebx
	double v4;  // st7
	double v5;  // st7
	int v6;     // eax
	double v7;  // st7
	double v8;  // st6
	double v9;  // st5
	double v10; // st4
	double v11; // st7
	float v12;  // et1
	double v13; // st7
	float v14;  // et1
	double v15; // st6
	float2 v16; // [esp+Ch] [ebp-38h]
	float v17;  // [esp+14h] [ebp-30h]
	float v18;  // [esp+18h] [ebp-2Ch]
	float v19;  // [esp+1Ch] [ebp-28h]
	float v20;  // [esp+20h] [ebp-24h]
	float v21;  // [esp+24h] [ebp-20h]
	float v22;  // [esp+28h] [ebp-1Ch]
	float v23;  // [esp+2Ch] [ebp-18h]
	float v24;  // [esp+30h] [ebp-14h]
	float v25;  // [esp+34h] [ebp-10h]
	float v26;  // [esp+38h] [ebp-Ch]
	float v27;  // [esp+3Ch] [ebp-8h]
	float v28;  // [esp+40h] [ebp-4h]
	float v29;  // [esp+48h] [ebp+4h]
	float v30;  // [esp+48h] [ebp+4h]

	v2 = a1;
	v3 = 1;
	v17 = (a1[16] + a1[17]) * 0.70710677;
	v18 = (a1[17] - a1[16]) * 0.70710677;
	v19 = (*(float*)(a2 + 64) + *(float*)(a2 + 68)) * 0.70710677;
	v20 = (*(float*)(a2 + 68) - *(float*)(a2 + 64)) * 0.70710677;
	v4 = a1[46] * 0.5;
	v21 = v17 - v4;
	v29 = a1[47] * 0.5;
	v22 = v18 - v29;
	v23 = v4 + v17;
	v24 = v29 + v18;
	v5 = *(float*)(a2 + 184) * 0.5;
	v25 = v19 - v5;
	v30 = *(float*)(a2 + 188) * 0.5;
	v26 = v20 - v30;
	v27 = v5 + v19;
	v28 = v30 + v20;
	if (v21 <= (double)v27 && v22 <= (double)v28 && v23 >= (double)v25 && v24 >= (double)v26) {
		v16.field_0 = *(float*)(a2 + 64) - v2[16];
		v16.field_4 = *(float*)(a2 + 68) - v2[17];
		nox_xxx_collSysAddCollision_548630(a2, (unsigned int)v2, &v16);
		if (((uint8_t)v2[4] & 8) == 8 || (*(uint8_t*)(a2 + 16) & 8) == 8) {
			v3 = 0;
		}
		if (!((uint8_t)v2[2] & 6) || (v6 = *(uint32_t*)(a2 + 16), !(v6 & 0x2000))) {
			if (v3) {
				if (v21 <= (double)v25) {
					v7 = v25;
				} else {
					v7 = v21;
				}
				if (v23 >= (double)v27) {
					v8 = v27;
				} else {
					v8 = v23;
				}
				if (v22 <= (double)v26) {
					v9 = v26;
				} else {
					v9 = v22;
				}
				if (v24 >= (double)v28) {
					v10 = v28;
				} else {
					v10 = v24;
				}
				v16.field_0 = v8 - v7;
				v16.field_4 = v10 - v9;
				if (v16.field_0 >= (double)v16.field_4) {
					if (v18 >= (double)v20) {
						v20 = v16.field_4;
					} else {
						v20 = -v16.field_4;
					}
					v11 = 0.0;
				} else {
					v11 = v16.field_0;
					if (v17 < (double)v19) {
						v11 = -v11;
					}
					v20 = 0.0;
				}
				v12 = *(float*)&dword_587000_292488;
				v13 = v11 * v12;
				v14 = *(float*)&dword_587000_292488;
				v15 = v20 * v14;
				v19 = (v13 - v15) * 0.70710677;
				v20 = (v15 + v13) * 0.70710677;
				sub_548600((int)v2, v19, v20);
			}
		}
		if ((uint32_t)v2[4] & 0x8000000) {
			nox_xxx_unitHasCollideOrUpdateFn_537610((int)v2);
			*((uint32_t*)v2 + 4) &= 0xF7FFFFFF;
		}
		if (*(uint32_t*)(a2 + 16) & 0x8000000) {
			nox_xxx_unitHasCollideOrUpdateFn_537610(a2);
			*(uint32_t*)(a2 + 16) &= 0xF7FFFFFF;
		}
	}
}

//----- (00551250) --------------------------------------------------------
void sub_551250(unsigned int a1, float* a2, int a3) {
	unsigned int v3; // edi
	float* v4;       // esi
	int v5;          // ebp
	double v6;       // st7
	double v7;       // st6
	double v8;       // st6
	int v9;          // ebx
	double v10;      // st7
	double v11;      // st6
	long double v12; // st6
	long double v13; // st7
	double v14;      // st6
	long double v15; // st7
	double v16;      // st7
	float v17;       // et1
	long double v18; // st7
	float v19;       // et1
	long double v20; // st7
	int v21;         // eax
	int v22;         // eax
	float* v23;      // eax
	int v24;         // ebx
	double v25;      // st7
	double v26;      // st7
	float v27;       // [esp+0h] [ebp-78h]
	float v28;       // [esp+4h] [ebp-74h]
	float2 v29;      // [esp+18h] [ebp-60h]
	float v30;       // [esp+20h] [ebp-58h]
	float v31;       // [esp+24h] [ebp-54h]
	float2 a4;       // [esp+28h] [ebp-50h]
	float2 v33;      // [esp+30h] [ebp-48h]
	float4 a2a;      // [esp+38h] [ebp-40h]
	float4 a3a;      // [esp+48h] [ebp-30h]
	float4 a1a;      // [esp+58h] [ebp-20h]
	float4 v37;      // [esp+68h] [ebp-10h]
	float v38;       // [esp+7Ch] [ebp+4h]
	float v39;       // [esp+7Ch] [ebp+4h]
	float v40;       // [esp+7Ch] [ebp+4h]
	float v41;       // [esp+80h] [ebp+8h]

	v3 = a1;
	v4 = a2;
	v5 = *(uint32_t*)(a1 + 748);
	v6 = (*(float*)(a1 + 68) + *(float*)(a1 + 64)) * 0.70710677;
	v29.field_4 = (*(float*)(a1 + 68) - *(float*)(a1 + 64)) * 0.70710677;
	v7 = a2[17] + a2[16];
	a1a.field_4 = v29.field_4;
	v30 = v7 * 0.70710677;
	v31 = (a2[17] - a2[16]) * 0.70710677;
	v8 = a2[46] * 0.5;
	a3a.field_0 = v30 - v8;
	v38 = a2[47] * 0.5;
	a3a.field_4 = v31 - v38;
	a3a.field_8 = v8 + v30;
	a3a.field_C = v38 + v31;
	a1a.field_0 = v6;
	v9 = *(short*)(v5 + 40) + 128;
	if (v9 >= 256) {
		v9 = *(short*)(v5 + 40) - 128;
	}
	a1a.field_8 = *getMemFloatPtr(0x587000, 194136 + 8 * v9) * 32.0 + v6;
	a1a.field_C = *getMemFloatPtr(0x587000, 194140 + 8 * v9) * 32.0 + v29.field_4;
	if (v6 >= a1a.field_8) {
		a2a.field_8 = v6;
		a2a.field_0 = a1a.field_8;
	} else {
		a2a.field_0 = v6;
		a2a.field_8 = a1a.field_8;
	}
	if (v29.field_4 >= (double)a1a.field_C) {
		a2a.field_4 = a1a.field_C;
		a2a.field_C = v29.field_4;
	} else {
		a2a.field_C = a1a.field_C;
		a2a.field_4 = v29.field_4;
	}
	if (a2a.field_0 <= (double)a3a.field_8 && a2a.field_4 <= (double)a3a.field_C &&
		a2a.field_8 >= (double)a3a.field_0 && a2a.field_C >= (double)a3a.field_4) {
		if (sub_551960(&a1a, &a2a, &a3a, &a4)) {
			v10 = v30 - a4.field_0;
			v11 = v31 - a4.field_4;
			v29.field_4 = v11;
			v12 = sqrt(v11 * v29.field_4 + v10 * v10);
			v39 = v12;
			if (v12 != 0.0) {
				v37.field_0 = v30;
				v37.field_4 = v31;
				*(float2*)&v37.field_8 = a4;
				v29.field_0 = v10 / v39;
				v29.field_4 = v29.field_4 / v39;
				if (sub_5516A0(&v37, &a3a, &v33, 1, 1) == 1) {
					v13 = sqrt((v33.field_4 - v31) * (v33.field_4 - v31) + (v33.field_0 - v30) * (v33.field_0 - v30));
					if (v13 != 0.0) {
						v14 = (v29.field_0 - v29.field_4) * 0.70710677;
						v29.field_4 = (v29.field_4 + v29.field_0) * 0.70710677;
						v29.field_0 = v14;
						v15 = v13 - v39;
						v40 = v15;
						if (v15 > 0.0) {
							nox_xxx_collSysAddCollision_548630((int)a2, v3, &v29);
							*(uint32_t*)(v5 + 44) = gameFrame();
							if (a3 == 1) {
								v41 = -(v29.field_4 * a2[21]) - v29.field_0 * a2[20];
								v16 = nox_xxx_objectGetMass_4E4A70((int)v4);
								v17 = *(float*)&dword_587000_292492;
								v18 = sqrt(v16 * v17 * 4.0);
								v19 = *(float*)&dword_587000_292492;
								v20 = v18 * v41 * 0.25 + v40 * v19;
								v28 = v20 * v29.field_4;
								v27 = v20 * v29.field_0;
								sub_548600((int)v4, v27, v28);
							}
							v21 = *((uint32_t*)v4 + 4);
							if (v21 & 0x8000000) {
								if (!(v21 & 8)) {
									nox_xxx_unitHasCollideOrUpdateFn_537610((int)v4);
								}
								*((uint32_t*)v4 + 4) &= 0xF7FFFFFF;
							}
							nox_xxx_unitHasCollideOrUpdateFn_537610(v3);
							if (!nox_xxx_servObjectHasTeam_419130(v3 + 48) ||
								*(uint32_t*)(v5 + 12) != *(uint32_t*)(v5 + 4) ||
								nox_xxx_servCompareTeams_419150(v3 + 48, (int)(v4 + 12))) {
								if (!a3 && !*(uint8_t*)(v5 + 1)) {
									v23 = *(float**)(v3 + 508);
									if (!v23 || v23 == v4) {
										v24 = v9 + 32;
										if (v24 >= 256) {
											v24 -= 256;
										}
										v25 = v40 * v4[30];
										if ((double)*getMemIntPtr(0x587000, 192088 + 8 * v24) *
													(v4[17] - *(float*)(v3 + 68)) -
												(double)*getMemIntPtr(0x587000, 192092 + 8 * v24) *
													(v4[16] - *(float*)(v3 + 64)) <=
											0.0) {
											v26 = v25 + *(float*)(v5 + 32);
										} else {
											v26 = *(float*)(v5 + 32) - v25;
										}
										*(float*)(v5 + 32) = v26;
										sub_548830(v5);
										nox_xxx_unitAddToUpdatable_4DA8D0(v3);
									}
								}
							} else if (gameFrame() > *(uint32_t*)(v3 + 136)) {
								v22 = *(unsigned char*)(v3 + 52);
								*(uint32_t*)(v3 + 136) = gameFrame() + gameFPS();
								nox_xxx_getTeamByID_418AB0(v22);
								nox_xxx_netPriMsgToPlayer_4DA2C0((int)v4, "objcoll.c:GateLockedMechanism", 0);
							}
						}
					}
				}
			}
		}
	}
}

//----- (005516A0) --------------------------------------------------------
int sub_5516A0(float4* a1, float4* a2, float2* a3, int a4, int a6) {
	int v5;     // ecx
	int v6;     // eax
	float2* v7; // edi
	int v9;     // [esp+10h] [ebp-8h]
	int v10;    // [esp+14h] [ebp-4h]

	v5 = 0;
	v6 = 0;
	v9 = 0;
	v10 = 0;
	if (a4 > 0) {
		v7 = a3;
		do {
			if (v10 >= 4) {
				break;
			}
			switch (v10) {
			case 0:
				v6 = sub_551780(a1, a2->field_0, a2->field_8, a2->field_4, v7, a6);
				break;
			case 1:
				v6 = sub_551870(a1, a2->field_8, a2->field_4, a2->field_C, v7, a6);
				break;
			case 2:
				v6 = sub_551780(a1, a2->field_0, a2->field_8, a2->field_C, v7, a6);
				break;
			case 3:
				v6 = sub_551870(a1, a2->field_0, a2->field_4, a2->field_C, v7, a6);
				break;
			default:
				break;
			}
			v5 = v9;
			if (v6) {
				v5 = v9 + 1;
				++v7;
				++v9;
			}
			++v10;
		} while (v5 < a4);
	}
	return v5;
}

//----- (00551780) --------------------------------------------------------
int sub_551780(float4* a1, float a2, float a3, float a4, float2* a5, int a6) {
	float4* v6; // ecx
	double v7;  // st7
	double v8;  // st7
	float v10;  // [esp+4h] [ebp+4h]
	float v11;  // [esp+4h] [ebp+4h]

	v6 = a1;
	v7 = a1->field_4 - a1->field_C;
	if (v7 == 0.0) {
		return 0;
	}
	v10 = (a4 - a1->field_C) / v7;
	if (v10 < 0.0 && !a6) {
		return 0;
	}
	if (v10 > 1.0) {
		return 0;
	}
	v8 = a2 - a3;
	if (v8 == 0.0) {
		return 0;
	}
	v11 = ((1.0 - v10) * v6->field_8 + v10 * v6->field_0 - a3) / v8;
	if (v11 < 0.0 || v11 > 1.0) {
		return 0;
	}
	a5->field_4 = a4;
	a5->field_0 = (1.0 - v11) * a3 + v11 * a2;
	return 1;
}

//----- (00551870) --------------------------------------------------------
int sub_551870(float4* a1, float a2, float a3, float a4, float2* a5, int a6) {
	float4* v6; // ecx
	double v7;  // st7
	double v8;  // st7
	float v10;  // [esp+4h] [ebp+4h]
	float v11;  // [esp+4h] [ebp+4h]

	v6 = a1;
	v7 = a1->field_0 - a1->field_8;
	if (v7 == 0.0) {
		return 0;
	}
	v10 = (a2 - a1->field_8) / v7;
	if (v10 < 0.0 && !a6) {
		return 0;
	}
	if (v10 > 1.0) {
		return 0;
	}
	v8 = a3 - a4;
	if (v8 == 0.0) {
		return 0;
	}
	v11 = ((1.0 - v10) * v6->field_C + v10 * v6->field_4 - a4) / v8;
	if (v11 < 0.0 || v11 > 1.0) {
		return 0;
	}
	a5->field_0 = a2;
	a5->field_4 = (1.0 - v11) * a4 + v11 * a3;
	return 1;
}

//----- (00551960) --------------------------------------------------------
int sub_551960(float4* a1, float4* a2, float4* a3, float2* a4) {
	int result;    // eax
	int v5;        // eax
	int v6;        // eax
	double v7;     // st7
	float2* v8;    // eax
	double v9;     // st7
	float2 a3a[2]; // [esp+8h] [ebp-10h]

	if (a2->field_0 < (double)a3->field_0 || a2->field_8 > (double)a3->field_8 || a2->field_4 < (double)a3->field_4 ||
		a2->field_C > (double)a3->field_C) {
		v5 = sub_5516A0(a1, a3, a3a, 2, 0) - 1;
		if (v5) {
			if (v5 == 1) {
				a4->field_0 = (a3a[1].field_0 + a3a[0].field_0) * 0.5;
				a4->field_4 = (a3a[1].field_4 + a3a[0].field_4) * 0.5;
				result = 1;
			} else {
				result = 0;
			}
		} else {
			v6 = sub_551A90((float2*)a1, a3);
			v7 = a3a[0].field_0;
			if (v6) {
				v8 = a4;
				a4->field_0 = (v7 + a1->field_0) * 0.5;
				v9 = a3a[0].field_4 + a1->field_4;
			} else {
				v8 = a4;
				a4->field_0 = (v7 + a1->field_8) * 0.5;
				v9 = a3a[0].field_4 + a1->field_C;
			}
			v8->field_4 = v9 * 0.5;
			result = 1;
		}
	} else {
		a4->field_0 = (a1->field_8 + a1->field_0) * 0.5;
		result = 1;
		a4->field_4 = (a1->field_C + a1->field_4) * 0.5;
	}
	return result;
}

//----- (00551A90) --------------------------------------------------------
int sub_551A90(float2* a1, float4* a2) {
	return a1->field_0 >= (double)a2->field_0 && a1->field_0 <= (double)a2->field_8 &&
		   a1->field_4 >= (double)a2->field_4 && a1->field_4 <= (double)a2->field_C;
}

//----- (00551AE0) --------------------------------------------------------
void sub_551AE0(int a1, int a2, int a3) {
	int v3;   // edi
	int v4;   // eax
	float v5; // [esp+0h] [ebp-10h]
	float v6; // [esp+0h] [ebp-10h]

	v3 = *(uint32_t*)(a1 + 748);
	if (!*getMemU32Ptr(0x5D4594, 2491808)) {
		sub_551BF0();
	}
	if (a3) {
		v4 = *(unsigned short*)(a2 + 4);
		if ((unsigned short)v4 != *getMemU32Ptr(0x5D4594, 2491792) && v4 != *getMemU32Ptr(0x5D4594, 2491796) &&
			v4 != *getMemU32Ptr(0x5D4594, 2491800) && v4 != *getMemU32Ptr(0x5D4594, 2491804)) {
			v5 = *(float*)(a2 + 104) - (double)*(int*)(v3 + 16);
			if (sub_419A10(v5) > 10.0) {
				if ((double)*(int*)(v3 + 16) > *(float*)(a2 + 104)) {
					if (*(uint32_t*)(a2 + 172) == 2) {
						sub_54AD50(a2, a1, 0);
					} else if (*(uint32_t*)(a2 + 172) == 3) {
						sub_550F80((float*)a2, a1);
					}
				}
			} else if (nox_xxx_map_57B850((float2*)(a1 + 64), (float*)(a1 + 172), (float2*)(a2 + 64))) {
				*(uint32_t*)(a2 + 16) = *(uint32_t*)(a2 + 16) & 0xFFFBFFFF | 0x100000;
				v6 = (double)*(int*)(v3 + 16) + 4.0;
				nox_xxx_unitRaise_4E46F0(a2, v6);
				*(uint32_t*)(a2 + 108) = 0;
			}
		}
	}
}

//----- (00551BF0) --------------------------------------------------------
void sub_551BF0() {
	*getMemU32Ptr(0x5D4594, 2491792) = nox_xxx_getNameId_4E3AA0("SmallFist");
	*getMemU32Ptr(0x5D4594, 2491796) = nox_xxx_getNameId_4E3AA0("MediumFist");
	*getMemU32Ptr(0x5D4594, 2491800) = nox_xxx_getNameId_4E3AA0("LargeFist");
	*getMemU32Ptr(0x5D4594, 2491804) = nox_xxx_getNameId_4E3AA0("Meteor");
	*getMemU32Ptr(0x5D4594, 2491808) = 1;
}

//----- (00551C40) --------------------------------------------------------
void sub_551C40(int a1, int a2) {
	int v2;       // edi
	int v3;       // ebp
	int v4;       // eax
	int v5;       // eax
	double v6;    // st7
	uint32_t* v7; // ebx
	float v8;     // [esp+0h] [ebp-14h]
	float v9;     // [esp+0h] [ebp-14h]
	float v10;    // [esp+18h] [ebp+4h]

	v2 = a1;
	v3 = *(uint32_t*)(a1 + 748);
	if (!*getMemU32Ptr(0x5D4594, 2491808)) {
		sub_551BF0();
	}
	if (*(uint32_t*)(v3 + 4)) {
		v4 = *(unsigned short*)(a2 + 4);
		if ((unsigned short)v4 != *getMemU32Ptr(0x5D4594, 2491792) && v4 != *getMemU32Ptr(0x5D4594, 2491796) &&
			v4 != *getMemU32Ptr(0x5D4594, 2491800) && v4 != *getMemU32Ptr(0x5D4594, 2491804)) {
			v5 = *(uint32_t*)(a2 + 172);
			if (v5 == 3) {
				if (*(float*)(a1 + 184) < (double)*(float*)(a2 + 184) ||
					*(float*)(a1 + 188) < (double)*(float*)(a2 + 188)) {
					return;
				}
			} else if (v5 == 2) {
				v6 = *(float*)(a2 + 176) + *(float*)(a2 + 176);
				if (v6 > *(float*)(a1 + 184) || v6 > *(float*)(a1 + 188)) {
					return;
				}
			}
			v7 = (uint32_t*)(a1 + 64);
			if (nox_xxx_map_57B850((float2*)(a1 + 64), (float*)(a1 + 172), (float2*)(a2 + 64))) {
				v10 = (double)(int)(*(uint32_t*)(*(uint32_t*)(*(uint32_t*)(v3 + 4) + 748) + 16) - 64);
				v8 = *(float*)(a2 + 104) - v10;
				if (sub_419A10(v8) > 10.0) {
					if (v10 <= -10.0) {
						*(uint32_t*)(a2 + 16) |= 0x40000u;
						*(uint32_t*)(a2 + 156) = *v7;
						*(uint32_t*)(a2 + 160) = *(uint32_t*)(v2 + 68);
						*(uint32_t*)(a2 + 164) = *(uint32_t*)(*(uint32_t*)(v3 + 4) + 64);
						*(uint32_t*)(a2 + 168) = *(uint32_t*)(*(uint32_t*)(v3 + 4) + 68);
					}
				} else {
					v9 = v10 + 4.0;
					*(uint32_t*)(a2 + 16) = *(uint32_t*)(a2 + 16) & 0xFFFBFFFF | 0x100000;
					nox_xxx_unitRaise_4E46F0(a2, v9);
					*(uint32_t*)(a2 + 108) = 0;
				}
			}
		}
	}
}
