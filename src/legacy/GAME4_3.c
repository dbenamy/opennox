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

//----- (0052F8A0) --------------------------------------------------------
void nox_xxx_lightningSpellDuration_52FFD0(int a1, int a2, int a3);
int nox_xxx_onFrameLightning_52F8A0(float a1) {
	int source;                // esi
	int v2;                    // eax
	int v4;                    // eax
	int v5;                    // eax
	int v6;                    // edi
	int v7;                    // ecx
	int owner;                 // edx
	int spellLevel;            // ebx
	int v10;                   // eax
	int v11;                   // ecx
	int v12;                   // edi
	int target;                // eax
	int index;                 // ecx
	int v15;                   // ecx
	int v16;                   // ecx
	int v17;                   // ecx
	int v18;                   // ecx
	int secondBounceTarget;    // eax
	int v20;                   // ebx
	uint32_t* v21;             // edi
	uint32_t* j;               // ebp
	int v23;                   // eax
	int v24;                   // eax
	int v25;                   // ecx
	int v26;                   // eax
	int v27;                   // edi
	char v28;                  // al
	unsigned char v29;         // al
	int v30;                   // eax
	int v31;                   // ebp
	int i;                     // edi
	int v33;                   // eax
	int v34;                   // edi
	float range1;              // [esp+0h] [ebp-20h]
	float range2;              // [esp+0h] [ebp-20h]
	float range3;              // [esp+0h] [ebp-20h]
	float range4;              // [esp+0h] [ebp-20h]
	float lightningSearchTime; // [esp+8h] [ebp-18h]
	float range5;              // [esp+1Ch] [ebp-4h]
	float lightningRange;      // [esp+24h] [ebp+4h]
	float damage;              // [esp+24h] [ebp+4h]

	source = LODWORD(a1);
	v2 = *(uint32_t*)(LODWORD(a1) + 16);
	if (v2) {
		if (nox_xxx_testUnitBuffs_4FF350(v2, 8)) {
			return 1;
		}
	} else if (!*(uint32_t*)(LODWORD(a1) + 20)) {
		return 1;
	}
	lightningRange = nox_xxx_gamedataGetFloat_419D40("LightningRange");
	if (*(uint32_t*)(source + 20)) { // Is the source a trap?
		*getMemU32Ptr(0x5D4594, 2487820) = *(uint32_t*)(source + 28);
		*getMemU32Ptr(0x5D4594, 2487824) = *(uint32_t*)(source + 32);
		nox_xxx_unitsGetInCircle_517F90((float2*)(source + 28), lightningRange, nox_xxx_lightningSpellTrapEffect_530020,
										*(uint32_t*)(source + 16));
		return 1;
	}
	if (*(uint8_t*)(*(uint32_t*)(source + 16) + 8) & 4 && !nox_xxx_unitGetOldMana_4EEC80(*(uint32_t*)(source + 16))) {
		return 1;
	}
	if ((unsigned int)(gameFrame() - *(uint32_t*)(source + 60)) > 2 &&
		sub_4E6BD0(*(uint32_t*)(source + 16))) {
		return 1;
	}
	v4 = *(uint32_t*)(source + 16);
	if (*(uint8_t*)(v4 + 8) & 2 && sub_4FEA70(v4, (float2*)(source + 28))) {
		return 1;
	}
	v5 = *(uint32_t*)(source + 104);
	if (v5) {
		do {
			v6 = *(uint32_t*)(v5 + 116);
			sub_4FE980(v5);
			v5 = v6;
		} while (v6);
	}
	v7 = *(uint32_t*)(source + 8);
	*(uint32_t*)(source + 104) = *(uint32_t*)(source + 108);
	*(uint32_t*)(source + 108) = 0;
	nox_xxx_lightningTarget_5d4594_2487908 = 0;
	nox_xxx_lightningTargetArrayIndex_5d4594_2487904 = 0;
	owner = *(uint32_t*)(source + 16);
	*getMemU32Ptr(0x5D4594, 2487844) = 0;
	spellLevel = *getMemU32Ptr(0x587000, 260380 + 4 * v7);
	*getMemU32Ptr(0x5D4594, 2487848) = 0;
	nox_xxx_lightningOwner_5d4594_2487900 = owner;
	*getMemU32Ptr(0x5D4594, 2487852) = 0;
	*getMemU32Ptr(0x5D4594, 2487856) = 0;
	*getMemU32Ptr(0x5D4594, 2487860) = 0;
	v10 = *(uint32_t*)(source + 16);
	if (!(*(uint8_t*)(v10 + 8) & 4) || (v11 = *(uint32_t*)(v10 + 748), (v12 = *(uint32_t*)(v11 + 288)) == 0) ||
		(!nox_xxx_unitIsEnemyTo_5330C0(v10, *(uint32_t*)(v11 + 288)) ||
				 nox_xxx_calcDistance_4E6C00(*(uint32_t*)(source + 16), v12) > lightningRange
			 ? (target = nox_xxx_lightningTarget_5d4594_2487908)
			 : (target = v12, nox_xxx_lightningTarget_5d4594_2487908 = v12),
		 !target)) {
		*(float*)&nox_xxx_lightningClosestTargetDistance_5d4594_2487912 = lightningRange * lightningRange;
		nox_xxx_unitsGetInCircle_517F90((float2*)(source + 28), lightningRange, nox_xxx_lightningCanAttackCheck_52FF10,
										*(uint32_t*)(source + 16));
		target = nox_xxx_lightningTarget_5d4594_2487908;
		if (!nox_xxx_lightningTarget_5d4594_2487908) {
			for (i = *(uint32_t*)(source + 104); i; i = *(uint32_t*)(i + 116)) {
				if (*(uint32_t*)(i + 48)) {
					nox_xxx_netStopRaySpell_4FEF90(i, *(uint32_t**)(i + 48));
				}
			}
			v33 = *(uint32_t*)(source + 104);
			if (v33) {
				do {
					v34 = *(uint32_t*)(v33 + 116);
					sub_4FE980(v33);
					v33 = v34;
				} while (v34);
			}
			*(uint32_t*)(source + 104) = 0;
			return 0;
		}
	}
	// Should be 0
	index = nox_xxx_lightningTargetArrayIndex_5d4594_2487904;
	// ARRAY! First member of 2487844 is the target itself...
	*getMemU32Ptr(0x5D4594, 2487844 + 4 * nox_xxx_lightningTargetArrayIndex_5d4594_2487904) = target;

	nox_xxx_lightningTargetArrayIndex_5d4594_2487904 = index + 1;
	// Second level of magic allows ONE jump from main target
	if (spellLevel > 1) {
		nox_xxx_lightningTarget_5d4594_2487908 = 0;
		*(float*)&nox_xxx_lightningClosestTargetDistance_5d4594_2487912 = lightningRange * lightningRange;
		range1 = lightningRange * 0.94999999;
		nox_xxx_unitsGetInCircle_517F90((float2*)(*getMemU32Ptr(0x5D4594, 2487844) + 56), range1,
										nox_xxx_lightningCanAttackCheck_52FF10, *getMemIntPtr(0x5D4594, 2487844));
		if (nox_xxx_lightningTarget_5d4594_2487908) {
			v15 = nox_xxx_lightningTargetArrayIndex_5d4594_2487904;
			// THIS SETS *getMemU32Ptr(0x5D4594, 2487848)!!!
			// Second member of  2487844 is the first bounce target from main target
			*getMemU32Ptr(0x5D4594, 2487844 + 4 * nox_xxx_lightningTargetArrayIndex_5d4594_2487904) =
				nox_xxx_lightningTarget_5d4594_2487908;
			nox_xxx_lightningTargetArrayIndex_5d4594_2487904 = v15 + 1;
		}
	}

	// Third level of magic allows TWO jump from main target
	if (spellLevel > 2) {
		nox_xxx_lightningTarget_5d4594_2487908 = 0;
		*(float*)&nox_xxx_lightningClosestTargetDistance_5d4594_2487912 = lightningRange * lightningRange;
		range2 = lightningRange * 0.89999998;
		nox_xxx_unitsGetInCircle_517F90((float2*)(*getMemU32Ptr(0x5D4594, 2487844) + 56), range2,
										nox_xxx_lightningCanAttackCheck_52FF10, *getMemIntPtr(0x5D4594, 2487844));
		if (nox_xxx_lightningTarget_5d4594_2487908) {
			v16 = nox_xxx_lightningTargetArrayIndex_5d4594_2487904;
			// THIS SETS *getMemU32Ptr(0x5D4594, 2487852)!!!
			// Since the range is less than the first bounce target check, it always sets third member
			// Third member of  2487844 is the second bounce target from main target
			*getMemU32Ptr(0x5D4594, 2487844 + 4 * nox_xxx_lightningTargetArrayIndex_5d4594_2487904) =
				nox_xxx_lightningTarget_5d4594_2487908;
			nox_xxx_lightningTargetArrayIndex_5d4594_2487904 = v16 + 1;
		}
	}

	if (*getMemU32Ptr(0x5D4594, 2487848)) {
		// Fourth level of magic allows TWO jump from main target and ONE from secondary target
		// OR one jump from main target and one jump from secondary target
		if (spellLevel > 3) {
			nox_xxx_lightningTarget_5d4594_2487908 = 0;
			*(float*)&nox_xxx_lightningClosestTargetDistance_5d4594_2487912 = lightningRange * lightningRange;
			range3 = lightningRange * 0.85000002;
			nox_xxx_unitsGetInCircle_517F90((float2*)(*getMemU32Ptr(0x5D4594, 2487848) + 56), range3,
											nox_xxx_lightningCanAttackCheck_52FF10,
											*(int*)&*getMemU32Ptr(0x5D4594, 2487848));
			if (nox_xxx_lightningTarget_5d4594_2487908) {
				v17 = nox_xxx_lightningTargetArrayIndex_5d4594_2487904;
				// THIS SETS *getMemU32Ptr(0x5D4594, 2487852)
				// OR *getMemU32Ptr(0x5D4594, 2487856), if we have a second bounce already!!!
				// third member of 2487844 is the second bounce target.
				*getMemU32Ptr(0x5D4594, 2487844 + 4 * nox_xxx_lightningTargetArrayIndex_5d4594_2487904) =
					nox_xxx_lightningTarget_5d4594_2487908;
				nox_xxx_lightningTargetArrayIndex_5d4594_2487904 = v17 + 1;
			}
		}
	}
	if (*getMemU32Ptr(0x5D4594, 2487852)) {
		if (spellLevel > 4) {
			// FIFTH level of magic allows TWO jump from main target and ONE from EACH secondary target
			// OR one jump from main target and TWO consecutive jumps!!!
			nox_xxx_lightningTarget_5d4594_2487908 = 0;
			range5 = lightningRange * lightningRange; // Is this intentional???
			*(float*)&nox_xxx_lightningClosestTargetDistance_5d4594_2487912 = lightningRange * lightningRange;
			range4 = lightningRange * 0.80000001;
			nox_xxx_unitsGetInCircle_517F90((float2*)(*getMemU32Ptr(0x5D4594, 2487852) + 56), range4,
											nox_xxx_lightningCanAttackCheck_52FF10,
											*(int*)&*getMemU32Ptr(0x5D4594, 2487852));
			if (nox_xxx_lightningTarget_5d4594_2487908) {
				v18 = nox_xxx_lightningTargetArrayIndex_5d4594_2487904;
				// THIS SETS 2487860 OR 2487856
				*getMemU32Ptr(0x5D4594, 2487844 + 4 * nox_xxx_lightningTargetArrayIndex_5d4594_2487904) =
					nox_xxx_lightningTarget_5d4594_2487908;
				nox_xxx_lightningTargetArrayIndex_5d4594_2487904 = v18 + 1;
			}
		}
	}

	nox_xxx_lightningSpellDuration_52FFD0(source, *(uint32_t*)(source + 16), *getMemIntPtr(0x5D4594, 2487844));

	if (spellLevel > 1 && *getMemU32Ptr(0x5D4594, 2487848)) {
		nox_xxx_lightningSpellDuration_52FFD0(source, *getMemIntPtr(0x5D4594, 2487844),
											  *getMemU32Ptr(0x5D4594, 2487848));
	}

	secondBounceTarget = *getMemU32Ptr(0x5D4594, 2487852);
	if (spellLevel > 2 && *getMemU32Ptr(0x5D4594, 2487852)) {
		nox_xxx_lightningSpellDuration_52FFD0(source, *getMemIntPtr(0x5D4594, 2487844),
											  *getMemU32Ptr(0x5D4594, 2487852));
		secondBounceTarget = *getMemU32Ptr(0x5D4594, 2487852);
	}

	if (spellLevel > 3 && *getMemU32Ptr(0x5D4594, 2487856)) {
		if (*getMemU32Ptr(0x5D4594, 2487848)) {
			nox_xxx_lightningSpellDuration_52FFD0(source, *getMemU32Ptr(0x5D4594, 2487848), *getMemU32Ptr(0x5D4594, 2487856));
			secondBounceTarget = *getMemU32Ptr(0x5D4594, 2487852);
			goto LABEL_55;
		}
		if (secondBounceTarget) {
			nox_xxx_lightningSpellDuration_52FFD0(source, secondBounceTarget, *getMemU32Ptr(0x5D4594, 2487856));
			secondBounceTarget = *getMemU32Ptr(0x5D4594, 2487852);
			goto LABEL_55;
		}
	}
LABEL_55:
	if (spellLevel > 4) {
		// This is the fourth jump!!!
		if (*getMemU32Ptr(0x5D4594, 2487860)) {
			if (secondBounceTarget || (secondBounceTarget = *getMemU32Ptr(0x5D4594, 2487848)) != 0) {
				nox_xxx_lightningSpellDuration_52FFD0(source, secondBounceTarget, *getMemIntPtr(0x5D4594, 2487860));
			}
		}
	}
	// No main target
	if (!*getMemU32Ptr(0x5D4594, 2487844)) {
		return 0;
	}

	damage = nox_xxx_gamedataGetFloat_419D40("LightningDamage") + *(float*)(source + 76);
	v20 = nox_float2int(damage);
	*(float*)(source + 76) = damage - (double)v20;
	v21 = *(uint32_t**)(source + 108);
	for (j = *(uint32_t**)(source + 104); v21; v21 = (uint32_t*)v21[29]) {
		if (j) {
			v23 = j[12];
			if (v21[12] != v23 || v21[4] != j[4]) {
				if (v23) {
					nox_xxx_netStopRaySpell_4FEF90((int)j, (uint32_t*)j[12]);
				}
				nox_xxx_netStartDurationRaySpell_4FF130((int)v21);
			}
			j = (uint32_t*)j[29];
		} else {
			nox_xxx_netStartDurationRaySpell_4FF130((int)v21);
		}
		if (v20 > 0) {
			(*(void (**)(uint32_t, uint32_t, uint32_t, int, int))(v21[12] + 716))(v21[12], *(uint32_t*)(source + 16), 0,
																				  v20, 17);
		}
		v24 = v21[12];
		if (*(uint32_t*)(v24 + 16) & 0x8020) {
			nox_xxx_netSendPointFx_522FF0(129, (float2*)(v24 + 56));
		}
	}
	for (; j; j = (uint32_t*)j[29]) {
		if (j[12]) {
			nox_xxx_netStopRaySpell_4FEF90((int)j, (uint32_t*)j[12]);
		}
	}
	v25 = *(uint32_t*)(source + 16);
	if (*(uint8_t*)(v25 + 8) & 4) {
		v26 = *(uint32_t*)(source + 72);
		if (v26) {
			v27 = *(uint32_t*)(v26 + 736);
			v28 = *(uint8_t*)(v27 + 108);
			if (!v28) {
				return 1;
			}
			v29 = v28 - 1;
			*(uint8_t*)(v27 + 108) = v29;
			*(uint32_t*)(v27 + 112) = 100 * v29 / *(unsigned char*)(v27 + 109);
			v30 = *(uint32_t*)(source + 16);
			if (v30 && *(uint8_t*)(v30 + 8) & 4) {
				v31 = *(uint32_t*)(v30 + 748);
				nox_xxx_playerSetState_4FA020((uint32_t*)v30, 22);
				nox_xxx_netReportCharges_4D82B0(*(unsigned char*)(*(uint32_t*)(v31 + 276) + 2064),
												*(uint32_t**)(source + 72), *(uint8_t*)(v27 + 108),
												*(uint8_t*)(v27 + 109));
			}
			if (!*(uint8_t*)(v27 + 108)) {
				return 1;
			}
		} else {
			nox_xxx_playerSetState_4FA020((uint32_t*)v25, 10);
			nox_xxx_playerManaSub_4EEBF0(*(uint32_t*)(source + 16), 1);
			if (!nox_xxx_unitGetOldMana_4EEC80(*(uint32_t*)(source + 16))) {
				return 1;
			}
		}
	}
	if (!(gameFrame() % (gameFPS() / 3u))) {
		nox_xxx_aud_501960(78, *(uint32_t*)(source + 16), 0, 0);
		nox_xxx_aud_501960(78, *getMemIntPtr(0x5D4594, 2487844), 0, 0);
	}

	lightningSearchTime = nox_xxx_gamedataGetFloat_419D40("LightningSearchTime");
	*(uint32_t*)(source + 68) = gameFrame() + nox_float2int(lightningSearchTime);

	return 0;
}

//----- (0052FF10) --------------------------------------------------------
void nox_xxx_lightningCanAttackCheck_52FF10(int target, int source) {
	int owner;                         // eax
	int v3;                            // ecx
	int index;                         // eax
	unsigned char* ptrTargetFromArray; // ecx
	double xDistance;                  // st7
	double yDistance;                  // st6
	double distance;                   // st5

	if (*(uint32_t*)(target + 8) & 0x20006) { // Type check
		// Checks if the target is an enemy of the owner, if not, stop
		owner = nox_xxx_lightningOwner_5d4594_2487900;
		if (nox_xxx_lightningOwner_5d4594_2487900) {
			if (!nox_xxx_unitIsEnemyTo_5330C0(*(int*)&nox_xxx_lightningOwner_5d4594_2487900, target)) {
				return;
			}
			owner = nox_xxx_lightningOwner_5d4594_2487900;
		}
		if (!(*(uint8_t*)(target + 8) & 2) || (v3 = *(uint32_t*)(target + 12), (v3 & 0x8000) == 0)) {
			if (!(*(uint32_t*)(target + 16) & 0x8020) && target != source && target != owner) {
				index = 0;
				if (*(int*)&nox_xxx_lightningTargetArrayIndex_5d4594_2487904 > 0) {
					ptrTargetFromArray = getMemAt(0x5D4594, 2487844);
					while (*(uint32_t*)ptrTargetFromArray != target) {
						++index;
						ptrTargetFromArray += 4;
						if (index >= *(int*)&nox_xxx_lightningTargetArrayIndex_5d4594_2487904) {
							goto LABEL_14;
						}
					}
					return;
				}
			LABEL_14:
				if (nox_xxx_unitCanInteractWith_5370E0(source, target, 0)) {
					xDistance = *(float*)(target + 56) - *(float*)(source + 56);
					yDistance = *(float*)(target + 60) - *(float*)(source + 60);
					distance = yDistance * yDistance + xDistance * xDistance;
					if (distance < *(float*)&nox_xxx_lightningClosestTargetDistance_5d4594_2487912) {
						*(float*)&nox_xxx_lightningClosestTargetDistance_5d4594_2487912 = distance;
						nox_xxx_lightningTarget_5d4594_2487908 = target;
					}
				}
			}
		}
	}
}

//----- (00530020) --------------------------------------------------------
void nox_xxx_lightningSpellTrapEffect_530020(int a1, int a2) {
	int v2;    // eax
	int v3;    // eax
	float v4;  // edx
	float v5;  // eax
	int v6;    // eax
	int v7;    // eax
	float v8;  // [esp+0h] [ebp-20h]
	float4 v9; // [esp+10h] [ebp-10h]

	if (a1 != a2) {
		v2 = *(uint32_t*)(a1 + 8);
		if (v2 & 6) {
			if (!(*(uint32_t*)(a1 + 16) & 0x8020)) {
				if (!(v2 & 2) || (v3 = *(uint32_t*)(a1 + 12), (v3 & 0x8000) == 0)) {
					if (!a2 || nox_xxx_unitIsEnemyTo_5330C0(a2, a1)) {
						v4 = *(float*)(a1 + 56);
						v9.field_4 = *getMemFloatPtr(0x5D4594, 2487824);
						v9.field_0 = *getMemFloatPtr(0x5D4594, 2487820);
						v5 = *(float*)(a1 + 60);
						v9.field_8 = v4;
						v9.field_C = v5;
						if (nox_xxx_mapTraceRay_535250(&v9, 0, 0, 9)) {
							v8 = nox_xxx_gamedataGetFloat_419D40("LightningGlyphDamage");
							v6 = nox_float2int(v8);
							(*(void (**)(int, uint32_t, uint32_t, int, int))(a1 + 716))(a1, 0, 0, v6, 17);
							nox_xxx_netSendPointFx_522FF0(129, (float2*)(a1 + 56));
							v7 = nox_xxx_spellGetAud44_424800(43, 0);
							nox_xxx_aud_501960(v7, a1, 0, 0);
						}
					}
				}
			}
		}
	}
}

//----- (00530100) --------------------------------------------------------
char sub_530100(uint32_t* a1) {
	int v1; // eax
	int v2; // esi
	int v3; // eax
	int v4; // esi
	int v5; // edi
	int v6; // edi

	v1 = a1[27];
	if (v1) {
		do {
			v2 = *(uint32_t*)(v1 + 116);
			sub_4FE980(v1);
			v1 = v2;
		} while (v2);
	}
	v3 = a1[26];
	a1[27] = 0;
	if (v3) {
		do {
			v4 = *(uint32_t*)(v3 + 116);
			sub_4FE980(v3);
			v3 = v4;
		} while (v4);
	}
	a1[26] = 0;
	v5 = a1[18];
	if (v5) {
		v6 = *(uint32_t*)(v5 + 736);
		v3 = *(uint32_t*)(v6 + 96);
		LOBYTE(v3) = v3 & 0xFB;
		*(uint32_t*)(v6 + 96) = v3;
	}
	return v3;
}

//----- (00530160) --------------------------------------------------------
int nox_xxx_spellTagCreature_530160(uint32_t* a1) {
	int v1;       // eax
	int v2;       // edi
	int v3;       // eax
	uint32_t* v4; // eax
	short v5;     // ax
	int v6;       // ecx
	short v7;     // dx
	float v9;     // [esp+0h] [ebp-14h]
	char v10[7];  // [esp+Ch] [ebp-8h]

	v1 = a1[4];
	v2 = *(uint32_t*)(v1 + 748);
	if (!v1) {
		return 1;
	}
	if (*(uint32_t*)(v1 + 16) & 0x8020) {
		return 1;
	}
	if (!(*(uint8_t*)(v1 + 8) & 4)) {
		return 1;
	}
	v3 = a1[12];
	if (!v3 || *(uint32_t*)(v3 + 16) & 0x8020) {
		return 1;
	}
	v9 = nox_xxx_gamedataGetFloat_419D40("TagDurationPerLevel");
	a1[17] = gameFrame() + a1[2] * nox_float2int(v9);
	nox_xxx_netMarkMinimapObject_417190(*(unsigned char*)(*(uint32_t*)(v2 + 276) + 2064), a1[12], 1);
	v4 = (uint32_t*)a1[12];
	v10[0] = -46;
	v5 = nox_xxx_netGetUnitCodeServ_578AC0(v4);
	v6 = a1[12];
	*(uint16_t*)&v10[1] = v5;
	v7 = *(uint16_t*)(v6 + 4);
	v10[5] = 1;
	*(uint16_t*)&v10[3] = v7;
	v10[6] = 1;
	nox_xxx_netSendPacket0_4E5420(*(unsigned char*)(*(uint32_t*)(v2 + 276) + 2064), v10, 7, 0, 1);
	return 0;
}

//----- (00530250) --------------------------------------------------------
unsigned int sub_530250(int a1) {
	int v1;              // eax
	unsigned int result; // eax

	v1 = *(uint32_t*)(a1 + 48);
	if (v1) {
		result = ((*(uint32_t*)(v1 + 16) & 0xFFu) >> 5) & 1;
	} else {
		result = 1;
	}
	return result;
}

//----- (00530270) --------------------------------------------------------
int sub_530270(int a1) {
	int result;   // eax
	int v2;       // edi
	uint32_t* v3; // edx
	short v4;     // cx
	char v5[7];   // [esp+8h] [ebp-8h]

	result = *(uint32_t*)(a1 + 16);
	if (result) {
		if (*(uint8_t*)(result + 8) & 4) {
			v2 = *(uint32_t*)(result + 748);
			result = *(uint32_t*)(a1 + 48);
			if (result) {
				if (!(*(uint8_t*)(result + 8) & 4)) {
					nox_xxx_netUnmarkMinimapObj_417300(*(unsigned char*)(*(uint32_t*)(v2 + 276) + 2064), result, 1);
				}
				v3 = *(uint32_t**)(a1 + 48);
				v5[0] = -46;
				*(uint16_t*)&v5[1] = nox_xxx_netGetUnitCodeServ_578AC0(v3);
				v4 = *(uint16_t*)(*(uint32_t*)(a1 + 48) + 4);
				v5[5] = 2;
				*(uint16_t*)&v5[3] = v4;
				v5[6] = 1;
				result = nox_xxx_netSendPacket0_4E5420(*(unsigned char*)(*(uint32_t*)(v2 + 276) + 2064), v5, 7, 0, 1);
			}
		}
	}
	return result;
}

//----- (00530310) --------------------------------------------------------
int nox_xxx_spellBlink2_530310(uint32_t* a1) {
	int result; // eax
	float v2;   // [esp+0h] [ebp-8h]

	if (!nox_common_gameFlags_check_40A5C0(4096) || (result = 1, a1[5] != 1)) {
		if (nox_common_gameFlags_check_40A5C0(2048)) {
			v2 = nox_xxx_gamedataGetFloatTable_419D70("TeleportDelay", a1[2] - 1);
			a1[17] = gameFrame() + nox_float2int(v2);
			result = 0;
		} else {
			result = 0;
			a1[17] = gameFrame() + 1;
		}
	}
	return result;
}

//----- (00530380) --------------------------------------------------------
int nox_xxx_waypoint_579F00(float2* a1, nox_object_t* a2);
int nox_xxx_spellBlink1_530380(int* a1) {
	int v1;     // eax
	int v3;     // eax
	float2* v4; // ecx
	int v5;     // eax
	bool v6;    // zf
	int v7;     // eax
	int v8;     // eax
	int v9;     // [esp-Ch] [ebp-18h]
	int v10;    // [esp-Ch] [ebp-18h]
	int v11;    // [esp-8h] [ebp-14h]
	int v12;    // [esp-4h] [ebp-10h]
	float2 v13; // [esp+4h] [ebp-8h]

	v1 = a1[12];
	if (!v1) {
		return 1;
	}
	if (*(uint32_t*)(v1 + 16) & 0x8020) {
		return 1;
	}
	if (a1[17] - 1 != gameFrame()) {
		return 0;
	}
	if (nox_xxx_testUnitBuffs_4FF350(v1, 14)) {
		nox_xxx_aud_501960(231, a1[12], 0, 0);
		return 1;
	}
	if (nox_common_gameFlags_check_40A5C0(4096) && (v3 = a1[12], *(uint8_t*)(v3 + 8) & 4)) {
		v4 = *(float2**)(*(uint32_t*)(v3 + 748) + 308);
		if (v4) {
			sub_4ED970(60.0, v4 + 7, &v13);
		} else {
			nox_xxx_mapFindPlayerStart_4F7AB0(&v13, a1[12]);
		}
	} else if (!nox_xxx_waypoint_579F00(&v13, a1[12])) {
		nox_xxx_mapFindPlayerStart_4F7AB0(&v13, a1[12]);
	}
	nox_xxx_spellTeleportCreateWake_530560(a1[12], (int*)(a1[12] + 56), &v13);
	nox_xxx_netSendPointFx_522FF0(137, (float2*)(a1[12] + 56));
	v9 = a1[12];
	v5 = nox_xxx_spellGetAud44_424800(a1[1], 0);
	nox_xxx_aud_501960(v5, v9, 0, 0);
	if (!nox_xxx_testUnitBuffs_4FF350(a1[12], 0)) {
		nox_xxx_netSendPointFx_522FF0(137, (float2*)(a1[12] + 56));
		nox_xxx_netSendPointFx_522FF0(137, &v13);
	}
	nox_xxx_teleportToMB_4E7190((uint8_t*)a1[12], &v13.field_0);
	v6 = !nox_xxx_testUnitBuffs_4FF350(a1[12], 0);
	v7 = a1[12];
	if (v6) {
		nox_xxx_netSendPointFx_522FF0(137, (float2*)(v7 + 56));
		v12 = 0;
		v11 = 0;
		v10 = a1[12];
	} else {
		if (!(*(uint8_t*)(v7 + 8) & 4)) {
			sub_4E7540(a1[4], a1[12]);
			return 1;
		}
		v12 = *(uint32_t*)(v7 + 36);
		v11 = 2;
		v10 = a1[12];
	}
	v8 = nox_xxx_spellGetAud44_424800(a1[1], 0);
	nox_xxx_aud_501960(v8, v10, v11, v12);
	sub_4E7540(a1[4], a1[12]);
	return 1;
}

//----- (00530560) --------------------------------------------------------
uint32_t* nox_xxx_spellTeleportCreateWake_530560(int a1, int* a2, uint32_t* a3) {
	int v3;           // eax
	uint32_t* result; // eax
	uint32_t* v5;     // esi
	uint32_t* v6;     // eax

	v3 = *getMemU32Ptr(0x5D4594, 2487916);
	if (!*getMemU32Ptr(0x5D4594, 2487916)) {
		v3 = nox_xxx_getNameId_4E3AA0("TeleportWake");
		*getMemU32Ptr(0x5D4594, 2487916) = v3;
	}
	result = nox_xxx_newObjectWithTypeInd_4E3450(v3);
	v5 = result;
	if (result) {
		v6 = (uint32_t*)result[175];
		*v6 = *a3;
		v6[1] = a3[1];
		nox_xxx_createAt_4DAA50((int)v5, a1, *(float*)a2, *((float*)a2 + 1));
		result = (uint32_t*)(gameFrame() + gameFPS());
		v5[34] = gameFrame() + gameFPS();
	}
	return result;
}

//----- (005305D0) --------------------------------------------------------
int sub_5305D0(uint32_t* a1) {
	int v1;   // eax
	int v2;   // eax
	float v4; // [esp+0h] [ebp-8h]

	v1 = a1[3];
	if (*(uint8_t*)(v1 + 8) & 4) {
		v2 = *(uint32_t*)(v1 + 748);
		if (!v2 || !*(uint32_t*)(v2 + 4 * a1[1] - 372)) {
			return 1;
		}
	}
	if (nox_common_gameFlags_check_40A5C0(2048)) {
		v4 = nox_xxx_gamedataGetFloatTable_419D70("TeleportDelay", a1[2] - 1);
		a1[17] = gameFrame() + nox_float2int(v4);
	} else {
		a1[17] = gameFrame() + 1;
	}
	return 0;
}

//----- (00530650) --------------------------------------------------------
int sub_530650(int* a1) {
	int v1;     // ecx
	int v2;     // eax
	int v4;     // eax
	int v5;     // ebp
	int v6;     // edi
	int v7;     // edx
	int v8;     // ecx
	int v9;     // eax
	int v10;    // eax
	int v11;    // eax
	char v12;   // al
	int v13;    // eax
	int v14;    // [esp-10h] [ebp-24h]
	int v15;    // [esp-Ch] [ebp-20h]
	int v16;    // [esp-8h] [ebp-1Ch]
	int v17;    // [esp-4h] [ebp-18h]
	float2 v18; // [esp+Ch] [ebp-8h]

	v1 = a1[3];
	if (!v1) {
		return 1;
	}
	if (*(uint8_t*)(v1 + 16) & 0x20) {
		return 1;
	}
	v2 = a1[12];
	if (!v2 || *(uint32_t*)(v2 + 16) & 0x8020 || v2 != v1 && !a1[5]) {
		return 1;
	}
	if (a1[17] - 1 != gameFrame()) {
		return 0;
	}
	if (nox_xxx_testUnitBuffs_4FF350(v2, 14)) {
		nox_xxx_aud_501960(231, a1[12], 0, 0);
		return 1;
	}
	v4 = a1[3];
	if (*(uint8_t*)(v4 + 8) & 4) {
		v5 = *(uint32_t*)(v4 + 748);
		v6 = a1[1] - 122;
		v7 = *(uint32_t*)(v5 + 4 * v6 + 116);
		if (!v7) {
			return 1;
		}
		v8 = a1[12];
		v18 = *(float2*)(v8 + 56);
		nox_xxx_spellTeleportCreateWake_530560(v8, (int*)(v8 + 56), (uint32_t*)(v7 + 56));
		v14 = a1[12];
		v9 = nox_xxx_spellGetAud44_424800(a1[1], 1);
		nox_xxx_aud_501960(v9, v14, 0, 0);
		nox_xxx_teleportToMB_4E7190((uint8_t*)a1[12], (float*)(*(uint32_t*)(v5 + 4 * v6 + 116) + 56));
		if (nox_xxx_testUnitBuffs_4FF350(a1[12], 0)) {
			v10 = a1[12];
			if (!(*(uint8_t*)(v10 + 8) & 4)) {
				goto LABEL_18;
			}
			v17 = *(uint32_t*)(v10 + 36);
			v16 = 2;
			v15 = a1[12];
		} else {
			nox_xxx_netSendPointFx_522FF0(137, &v18);
			nox_xxx_netSendPointFx_522FF0(137, (float2*)(a1[12] + 56));
			v17 = 0;
			v16 = 0;
			v15 = a1[12];
		}
		v11 = nox_xxx_spellGetAud44_424800(a1[1], 1);
		nox_xxx_aud_501960(v11, v15, v16, v17);
	LABEL_18:
		*(uint32_t*)(*(uint32_t*)(v5 + 4 * v6 + 116) + 136) = gameFrame();
		v12 = *(uint8_t*)(v6 + v5 + 156) - 1;
		*(uint8_t*)(v6 + v5 + 156) = v12;
		if (!v12) {
			nox_xxx_netSendPointFx_522FF0(129, (float2*)(*(uint32_t*)(v5 + 4 * v6 + 116) + 56));
			nox_xxx_delayedDeleteObject_4E5CC0(*(uint32_t*)(v5 + 4 * v6 + 116));
			*(uint32_t*)(v5 + 4 * v6 + 116) = 0;
		}
	}
	v13 = a1[4];
	if (!v13 || *(uint8_t*)(v13 + 16) & 0x20) {
		return 1;
	}
	sub_4E7540(v13, a1[12]);
	return 1;
}

//----- (00530820) --------------------------------------------------------
int nox_xxx_castTele_530820(int a1) {
	int result; // eax
	float v2;   // [esp+0h] [ebp-8h]

	if (nox_common_gameFlags_check_40A5C0(2048)) {
		v2 = nox_xxx_gamedataGetFloatTable_419D70("TeleportDelay", *(uint32_t*)(a1 + 8) - 1);
		*(uint32_t*)(a1 + 68) = gameFrame() + nox_float2int(v2);
		result = 0;
	} else {
		result = 0;
		*(uint32_t*)(a1 + 68) = gameFrame() + 1;
	}
	return result;
}

//----- (00530880) --------------------------------------------------------
int sub_530880(int* a1) {
	int v1;       // eax
	int v2;       // eax
	int v4;       // eax
	int v5;       // edi
	int v6;       // edx
	int v7;       // ecx
	uint32_t* v8; // eax
	int v9;       // eax
	int v10;      // esi
	int v11;      // eax
	char v12;     // al
	int v13;      // [esp-Ch] [ebp-20h]
	float2 v14;   // [esp+Ch] [ebp-8h]

	v1 = a1[4];
	if (!v1) {
		return 1;
	}
	if (*(uint32_t*)(v1 + 16) & 0x8020) {
		return 1;
	}
	v2 = a1[12];
	if (!v2 || *(uint32_t*)(v2 + 16) & 0x8020) {
		return 1;
	}
	if (a1[17] - 1 != gameFrame()) {
		return 0;
	}
	if (nox_xxx_testUnitBuffs_4FF350(v2, 14)) {
		nox_xxx_aud_501960(231, a1[12], 0, 0);
		return 1;
	}
	v4 = a1[4];
	if (*(uint8_t*)(v4 + 8) & 4) {
		v5 = *(uint32_t*)(v4 + 748);
		v6 = 0;
		v7 = 4;
		v8 = (uint32_t*)(v5 + 116);
		do {
			if (*v8) {
				++v6;
			}
			++v8;
			--v7;
		} while (v7);
		if (!v6) {
			return 1;
		}
		v9 = a1[12];
		v14 = *(float2*)(a1[12] + 56);
		do {
			v10 = nox_common_randomInt_415FA0(0, 3);
		} while (!*(uint32_t*)(v5 + 4 * v10 + 116));
		nox_xxx_spellTeleportCreateWake_530560(a1[12], (int*)(a1[12] + 56),
											   (uint32_t*)(*(uint32_t*)(v5 + 4 * v10 + 116) + 56));
		nox_xxx_teleportToMB_4E7190((uint8_t*)a1[12], (float*)(*(uint32_t*)(v5 + 4 * v10 + 116) + 56));
		if (!nox_xxx_testUnitBuffs_4FF350(a1[12], 0)) {
			nox_xxx_netSendPointFx_522FF0(137, &v14);
			nox_xxx_netSendPointFx_522FF0(137, (float2*)(a1[12] + 56));
		}
		v13 = a1[12];
		v11 = nox_xxx_spellGetAud44_424800(a1[1], 0);
		nox_xxx_aud_501960(v11, v13, 0, 0);
		v12 = *(uint8_t*)(v10 + v5 + 156) - 1;
		*(uint8_t*)(v10 + v5 + 156) = v12;
		if (!v12) {
			nox_xxx_netSendPointFx_522FF0(129, (float2*)(*(uint32_t*)(v5 + 4 * v10 + 116) + 56));
			nox_xxx_delayedDeleteObject_4E5CC0(*(uint32_t*)(v5 + 4 * v10 + 116));
			*(uint32_t*)(v5 + 4 * v10 + 116) = 0;
		}
	}
	sub_4E7540(a1[4], a1[12]);
	return 1;
}

//----- (00530B70) --------------------------------------------------------
int nox_xxx_castTTT_530B70(int* a1) {
	int v1;  // ecx
	int v2;  // eax
	int v4;  // eax
	int v5;  // eax
	int v6;  // eax
	int v7;  // [esp-Ch] [ebp-14h]
	int v8;  // [esp-Ch] [ebp-14h]
	int v9;  // [esp-8h] [ebp-10h]
	int v10; // [esp-4h] [ebp-Ch]

	v1 = a1[12];
	if (!v1) {
		return 1;
	}
	v2 = a1[4];
	if (!v2 || *(uint32_t*)(v1 + 16) & 0x8020 || *(uint32_t*)(v2 + 16) & 0x8020) {
		return 1;
	}
	if (a1[17] - 1 != gameFrame()) {
		return 0;
	}
	if (nox_xxx_testUnitBuffs_4FF350(v1, 14)) {
		nox_xxx_aud_501960(231, a1[12], 0, 0);
		return 1;
	}
	nox_xxx_netSendPointFx_522FF0(137, (float2*)(a1[12] + 56));
	v7 = a1[12];
	v4 = nox_xxx_spellGetAud44_424800(a1[1], 0);
	nox_xxx_aud_501960(v4, v7, 0, 0);
	nox_xxx_spellTeleportCreateWake_530560(a1[12], (int*)(a1[12] + 56), a1 + 13);
	nox_xxx_teleportToMB_4E7190((uint8_t*)a1[12], (float*)a1 + 13);
	if (nox_xxx_testUnitBuffs_4FF350(a1[12], 0)) {
		v6 = a1[12];
		if (!(*(uint8_t*)(v6 + 8) & 4)) {
			sub_4E7540(a1[4], a1[12]);
			return 1;
		}
		v10 = *(uint32_t*)(v6 + 36);
		v9 = 2;
		v8 = a1[12];
		v5 = nox_xxx_spellGetAud44_424800(a1[1], 0);
	} else {
		nox_xxx_netSendPointFx_522FF0(137, (float2*)(a1[12] + 56));
		v10 = 0;
		v9 = 0;
		v8 = a1[12];
		v5 = nox_xxx_spellGetAud44_424800(a1[1], 0);
	}
	nox_xxx_aud_501960(v5, v8, v9, v10);
	sub_4E7540(a1[4], a1[12]);
	return 1;
}

//----- (00530CA0) --------------------------------------------------------
int sub_530CA0(int a1) {
	int v1;     // ecx
	int v2;     // eax
	int v3;     // ecx
	int result; // eax
	float v5;   // [esp+0h] [ebp-8h]

	v1 = *(uint32_t*)(a1 + 16);
	if (!v1) {
		return 1;
	}
	if (*(uint8_t*)(a1 + 88) & 0x20) {
		return 1;
	}
	v2 = *(uint32_t*)(a1 + 48);
	if (!v2) {
		return 1;
	}
	if (*(uint32_t*)(v2 + 16) & 0x8020) {
		return 1;
	}
	if (v1 == v2) {
		return 1;
	}
	if (*(uint8_t*)(v2 + 8) & 2) {
		v3 = *(uint32_t*)(v2 + 12);
		if (v3 & 0x4000) {
			return 1;
		}
	}
	if (nox_common_gameFlags_check_40A5C0(2048)) {
		v5 = nox_xxx_gamedataGetFloatTable_419D70("TeleportDelay", *(uint32_t*)(a1 + 8) - 1);
		*(uint32_t*)(a1 + 68) = gameFrame() + nox_float2int(v5);
		result = 0;
	} else {
		result = 0;
		*(uint32_t*)(a1 + 68) = gameFrame() + 1;
	}
	return result;
}

//----- (00530D30) --------------------------------------------------------
int sub_530D30(int* a1) {
	int* v1;    // esi
	int v2;     // ecx
	int v3;     // eax
	int v5;     // eax
	int v6;     // ecx
	float2* v7; // ecx
	int v8;     // eax
	int v9;     // eax
	int v10;    // eax
	int v11;    // [esp-Ch] [ebp-24h]
	int v12;    // [esp-Ch] [ebp-24h]
	float* v13; // [esp-4h] [ebp-1Ch]
	float4 v14; // [esp+8h] [ebp-10h]
	int v18;    // [esp+1Ch] [ebp+4h]

	v1 = a1;
	v2 = a1[12];
	if (!v2) {
		return 1;
	}
	v3 = a1[4];
	if (!v3 || *(uint32_t*)(v2 + 16) & 0x8020 || *(uint32_t*)(v3 + 16) & 0x8020) {
		return 1;
	}
	if (a1[17] - 1 != gameFrame()) {
		return 0;
	}
	if (nox_xxx_testUnitBuffs_4FF350(v2, 14) || nox_xxx_testUnitBuffs_4FF350(a1[4], 14)) {
		nox_xxx_aud_501960(231, a1[12], 0, 0);
		nox_xxx_aud_501960(231, a1[4], 0, 0);
		return 1;
	}
	if (a1[5]) {
		goto LABEL_23;
	}
	if (!nox_xxx_unitCanInteractWith_5370E0(a1[4], a1[12], 0)) {
		nox_xxx_netPriMsgToPlayer_4DA2C0(a1[4], "ExecDur.c:NeedClearLOSForSwap", 0);
		return 1;
	}
	v5 = a1[4];
	if (!(!(*(uint8_t*)(v5 + 8) & 4) ||
		(v6 = *(uint32_t*)(v5 + 748),
		 v14.field_0 = *(float*)(v5 + 56) - (double)*(unsigned short*)(*(uint32_t*)(v6 + 276) + 10),
		 v14.field_4 = *(float*)(v5 + 60) - (double)*(unsigned short*)(*(uint32_t*)(v6 + 276) + 12),
		 v14.field_8 = (double)*(unsigned short*)(*(uint32_t*)(v6 + 276) + 10) + *(float*)(v5 + 56),
		 v18 = *(unsigned short*)(*(uint32_t*)(v6 + 276) + 12), v7 = (float2*)(v1[12] + 56),
		 v14.field_C = (double)v18 + *(float*)(v5 + 60), sub_428220(v7, &v14)))) {
		nox_xxx_netPriMsgToPlayer_4DA2C0(v1[4], "ExecDur.c:NeedClearLOSForSwap", 0);
		return 1;
	}
LABEL_23:
	v8 = v1[12];
	v14.field_0 = *(float*)(v8 + 56);
	v13 = (float*)(v1[4] + 56);
	v14.field_4 = *(float*)(v8 + 60);
	nox_xxx_teleportToMB_4E7190((uint8_t*)v8, v13);
	nox_xxx_teleportToMB_4E7190((uint8_t*)v1[4], &v14);
	if (!nox_xxx_testUnitBuffs_4FF350(v1[12], 0) && !nox_xxx_testUnitBuffs_4FF350(v1[4], 0)) {
		nox_xxx_netSendPointFx_522FF0(137, (float2*)(v1[4] + 56));
		nox_xxx_netSendPointFx_522FF0(137, (float2*)(v1[12] + 56));
	}
	v11 = v1[12];
	v9 = nox_xxx_spellGetAud44_424800(v1[1], 0);
	nox_xxx_aud_501960(v9, v11, 0, 0);
	v12 = v1[4];
	v10 = nox_xxx_spellGetAud44_424800(v1[1], 0);
	nox_xxx_aud_501960(v10, v12, 0, 0);
	sub_4E7540(v1[4], v1[12]);
	return 1;
}

//----- (00530F90) --------------------------------------------------------
int nox_xxx_manaBomb_530F90(uint32_t* a1) {
	int v1;       // ecx
	uint32_t* v2; // ebx
	float v4;     // [esp+0h] [ebp-10h]
	float v5;     // [esp+0h] [ebp-10h]

	if (!*getMemU32Ptr(0x5D4594, 2487920)) {
		*getMemU32Ptr(0x5D4594, 2487920) = nox_xxx_getNameId_4E3AA0("ManaBombCharge");
	}
	v4 = nox_xxx_gamedataGetFloatTable_419D70("ManaBombInitPower", a1[2] - 1);
	a1[18] = nox_float2int(v4);
	a1[19] = 0;
	a1[21] = 0;
	v1 = a1[4];
	if (v1 && !(*(uint8_t*)(v1 + 8) & 4) || a1[5]) {
		v5 = nox_xxx_gamedataGetFloat_419D40("ManaBombGlyphDuration");
		a1[17] = gameFrame() + nox_float2int(v5);
	} else {
		nox_xxx_buffApplyTo_4FF380(v1, 5, 10 * (uint16_t)gameFPS(), 5);
		nox_xxx_buffApplyTo_4FF380(a1[4], 14, 10 * (uint16_t)gameFPS(), 5);
		nox_xxx_buffApplyTo_4FF380(a1[4], 29, 10 * (uint16_t)gameFPS(), 5);
		a1[20] = *(uint32_t*)(a1[4] + 120);
		*(uint32_t*)(a1[4] + 120) = 1203982323;
		*(uint32_t*)(a1[4] + 84) = 0;
		*(uint32_t*)(a1[4] + 80) = 0;
		*(uint32_t*)(a1[4] + 92) = 0;
		*(uint32_t*)(a1[4] + 88) = 0;
		*(uint32_t*)(a1[4] + 100) = 0;
		*(uint32_t*)(a1[4] + 96) = 0;
	}
	v2 = nox_xxx_newObjectWithTypeInd_4E3450(*getMemIntPtr(0x5D4594, 2487920));
	if (v2) {
		nox_xxx_createAt_4DAA50((int)v2, 0, *((float*)a1 + 7), *((float*)a1 + 8));
		a1[19] = v2;
	}
	return 0;
}

//----- (005310C0) --------------------------------------------------------
int nox_xxx_manaBombBoom_5310C0(int* a1) {
	int v1;     // edi
	int v2;     // eax
	int v4;     // eax
	bool v5;    // zf
	int v6;     // edi
	float v7;   // edx
	int v8;     // eax
	int v9;     // eax
	float v10;  // [esp+0h] [ebp-28h]
	float v11;  // [esp+4h] [ebp-24h]
	int v12;    // [esp+10h] [ebp-18h]
	float v13;  // [esp+14h] [ebp-14h]
	float v14;  // [esp+14h] [ebp-14h]
	float2 v15; // [esp+20h] [ebp-8h]

	v1 = 0;
	v2 = a1[4];
	if (!v2 && !a1[5]) {
		return 1;
	}
	if (a1[17] - 1 == gameFrame()) {
		v1 = 1;
	}
	if (v2 && *(uint8_t*)(v2 + 8) & 4 && !nox_xxx_unitGetOldMana_4EEC80(a1[4])) {
		v1 = 1;
	}
	if (a1[19]) {
		if (!a1[5] && (v4 = a1[4]) != 0 && *(uint8_t*)(v4 + 8) & 4) {
			if ((unsigned short)nox_xxx_unitGetOldMana_4EEC80(a1[4]) >= 0xFu) {
				goto LABEL_18;
			}
		} else if ((unsigned int)(a1[17] - gameFrame()) >= 0xA) {
			goto LABEL_18;
		}
		nox_xxx_delayedDeleteObject_4E5CC0(a1[19]);
		a1[19] = 0;
	}
LABEL_18:
	v5 = v1 == 0;
	v6 = a1[18];
	if (!v5) {
		if (a1[5]) {
			v7 = *((float*)a1 + 8);
			v15.field_0 = *((float*)a1 + 7);
		} else {
			v8 = a1[4];
			if (!v8) {
				return 1;
			}
			v15.field_0 = *(float*)(v8 + 56);
			v7 = *(float*)(v8 + 60);
		}
		v15.field_4 = v7;
		nox_xxx_gameSetWallsDamage_4E25A0(1);
		v12 = a1[4];
		v11 = nox_xxx_gamedataGetFloat_419D40("ManaBombInRadius");
		v10 = nox_xxx_gamedataGetFloat_419D40("ManaBombOutRadius");
		nox_xxx_mapDamageUnitsAround_4E25B0(&v15, v10, v11, v6, 15, v12, 0);
		v13 = nox_xxx_gamedataGetFloat_419D40("ManaBombShakeMag");
		v9 = nox_float2int(v13);
		nox_xxx_earthquakeSend_4D9110(&v15, v9);
		nox_xxx_netSendPointFx_522FF0(129, &v15);
		nox_xxx_netSendPointFx_522FF0(154, &v15);
		nox_xxx_aud_501960(81, a1[4], 0, 0);
		a1[21] = 1;
		return 1;
	}
	v14 = nox_xxx_gamedataGetFloatTable_419D70("ManaBombDeltaPower", a1[2] - 1);
	a1[18] = v6 + nox_float2int(v14);
	if (!a1[5]) {
		if (sub_4E7BC0(a1[4])) {
			nox_xxx_playerManaSub_4EEBF0(a1[4], a1[2]);
		}
	}
	return 0;
}

//----- (00531290) --------------------------------------------------------
int sub_531290(int a1) {
	int v1;     // eax
	int result; // eax

	if (*(uint32_t*)(a1 + 76)) {
		nox_xxx_delayedDeleteObject_4E5CC0(*(uint32_t*)(a1 + 76));
		*(uint32_t*)(a1 + 76) = 0;
	}
	if (!*(uint32_t*)(a1 + 20)) {
		v1 = *(uint32_t*)(a1 + 16);
		if (v1) {
			if (*(uint8_t*)(v1 + 8) & 4) {
				nox_xxx_spellBuffOff_4FF5B0(v1, 5);
				nox_xxx_spellBuffOff_4FF5B0(*(uint32_t*)(a1 + 16), 14);
				nox_xxx_spellBuffOff_4FF5B0(*(uint32_t*)(a1 + 16), 29);
				*(uint32_t*)(*(uint32_t*)(a1 + 16) + 120) = *(uint32_t*)(a1 + 80);
			}
		}
	}
	result = *(uint32_t*)(a1 + 84);
	if (!result) {
		result = nox_xxx_netSendPointFx_522FF0(163, (float2*)(a1 + 28));
	}
	return result;
}

//----- (00531310) --------------------------------------------------------
int nox_xxx_spellTurnUndeadCreate_531310(uint32_t* a1) {
	int v1;       // ecx
	int v2;       // edi
	uint32_t* v3; // eax
	uint32_t* v4; // esi
	int v5;       // eax
	double v6;    // st7
	float v8;     // [esp+0h] [ebp-18h]
	float2 v9;    // [esp+10h] [ebp-8h]

	v8 = nox_xxx_gamedataGetFloatTable_419D70("TurnUndeadKillPoints", a1[2] - 1);
	a1[18] = nox_float2int(v8);
	if (!*getMemU32Ptr(0x5D4594, 2487924)) {
		*getMemU32Ptr(0x5D4594, 2487924) = nox_xxx_getNameId_4E3AA0("UndeadKiller");
	}
	if (a1[5]) {
		v1 = a1[6];
	} else {
		v1 = a1[4];
	}
	v2 = 0;
	v9 = *(float2*)(v1 + 56);
	do {
		v3 = nox_xxx_newObjectWithTypeInd_4E3450(*getMemIntPtr(0x5D4594, 2487924));
		v4 = v3;
		if (v3) {
			*(uint32_t*)v3[175] = a1;
			nox_xxx_createAt_4DAA50((int)v3, a1[4], v9.field_0, v9.field_4);
			*((uint16_t*)v4 + 63) = v2;
			v5 = 8 * (short)v2;
			*((uint16_t*)v4 + 62) = v2;
			*((float*)v4 + 20) = *getMemFloatPtr(0x587000, 194136 + v5) * 4.0;
			v6 = *getMemFloatPtr(0x587000, 194140 + v5) * 4.0;
			v4[28] = 0;
			*((float*)v4 + 21) = v6;
		}
		v2 += 6;
	} while (v2 < 256);
	nox_xxx_netSendPointFx_522FF0(160, &v9);
	return 0;
}

//----- (00531410) --------------------------------------------------------
int nox_xxx_spellTurnUndeadUpdate_531410() { return 0; }

//----- (00531420) --------------------------------------------------------
int nox_xxx_spellTurnUndeadDelete_531420(int a1) {
	int result; // eax
	int i;      // esi

	if (!*getMemU32Ptr(0x5D4594, 2487928)) {
		*getMemU32Ptr(0x5D4594, 2487928) = nox_xxx_getNameId_4E3AA0("UndeadKiller");
	}
	result = nox_server_getFirstObject_4DA790();
	for (i = result; result; i = result) {
		if (*(unsigned short*)(i + 4) == *getMemU32Ptr(0x5D4594, 2487928) && **(uint32_t**)(i + 700) == a1) {
			nox_xxx_delayedDeleteObject_4E5CC0(i);
		}
		result = nox_server_getNextObject_4DA7A0(i);
	}
	return result;
}

//----- (00531490) --------------------------------------------------------
void sub_4FF310(nox_object_t* a1);
int sub_531490(uint32_t* a1) {
	int v1;     // eax
	int v2;     // edi
	int result; // eax

	v1 = a1[12];
	v2 = 20 * a1[2] * gameFPS();
	if (!v1) {
		return 1;
	}
	if (!(*(uint8_t*)(v1 + 8) & 4)) {
		return 1;
	}
	sub_4FF310(a1[12]);
	nox_xxx_buffApplyTo_4FF380(a1[12], 27, v2, a1[2]);
	result = 0;
	a1[17] = v2 + gameFrame();
	return result;
}

//----- (005314F0) --------------------------------------------------------
int sub_5314F0(int a1) {
	int v1;     // eax
	int result; // eax
	int v3;     // eax

	v1 = *(uint32_t*)(a1 + 48);
	if (!v1) {
		return 1;
	}
	if (nox_xxx_testUnitBuffs_4FF350(v1, 8)) {
		return 1;
	}
	v3 = *(uint32_t*)(a1 + 48);
	if (v3 && *(uint8_t*)(v3 + 8) & 2 && sub_4FEA70(v3, (float2*)(a1 + 28))) {
		result = 1;
	} else {
		result = (*(uint32_t*)(*(uint32_t*)(a1 + 48) + 16) & 0x8020) != 0;
	}
	return result;
}

//----- (00531560) --------------------------------------------------------
int sub_531560(int a1) {
	int result; // eax

	result = *(uint32_t*)(a1 + 48);
	if (result) {
		result = nox_xxx_spellBuffOff_4FF5B0(result, 27);
	}
	return result;
}

//----- (00531580) --------------------------------------------------------
int nox_xxx_plasmaSmth_531580(int a1) {
	int v1; // eax
	int v2; // eax

	*(uint32_t*)(a1 + 72) = 0;
	*(uint32_t*)(a1 + 76) = 0;
	*(uint32_t*)(a1 + 48) = 0;
	nox_xxx_netSendPointFx_522FF0(131, (float2*)(a1 + 28));
	v1 = *(uint32_t*)(a1 + 16);
	if (!v1 || !(*(uint8_t*)(v1 + 8) & 4)) {
		return 0;
	}
	v2 = *(uint32_t*)(*(uint32_t*)(v1 + 748) + 104);
	if (v2 && *(uint32_t*)(v2 + 12) & 0x4000000) {
		if (*(uint8_t*)(*(uint32_t*)(v2 + 736) + 96) & 4) {
			*(uint32_t*)(a1 + 72) = v2;
			*(uint8_t*)(a1 + 88) |= 2u;
		}
		return 0;
	}
	return 1;
}

//----- (00531600) --------------------------------------------------------
int nox_xxx_plasmaShot_531600(int a1) {
	int v1;            // eax
	int v2;            // eax
	int v3;            // edi
	int v4;            // eax
	int v5;            // eax
	int v6;            // ecx
	int v7;            // eax
	int v8;            // eax
	double v9;         // st7
	int v10;           // eax
	int v11;           // eax
	int v12;           // eax
	int v13;           // edi
	unsigned char v14; // al
	unsigned char v15; // al
	int v16;           // eax
	int v17;           // ebp
	uint32_t* v19;     // [esp-4h] [ebp-18h]
	float v20;         // [esp+0h] [ebp-14h]
	float v21;         // [esp+0h] [ebp-14h]

	if (!*getMemU32Ptr(0x5D4594, 2487936)) {
		*getMemU32Ptr(0x5D4594, 2487936) = nox_xxx_getNameId_4E3AA0("Hecubah");
		*getMemU32Ptr(0x5D4594, 2487940) = nox_xxx_getNameId_4E3AA0("HecubahWithOrb");
	}
	v1 = *(uint32_t*)(a1 + 16);
	if (!v1 || *(uint8_t*)(a1 + 88) & 0x20 || *(uint8_t*)(v1 + 8) & 2 && sub_4FEA70(v1, (float2*)(a1 + 28))) {
		return 1;
	}
	v2 = *(uint32_t*)(a1 + 48);
	if (!v2) {
		goto LABEL_16;
	}
	if (*(uint32_t*)(v2 + 16) & 0x8020) {
		*(uint32_t*)(a1 + 48) = 0;
		goto LABEL_16;
	}
	if (nox_server_testTwoPointsAndDirection_4E6E50((float2*)(*(uint32_t*)(a1 + 16) + 56), *(short*)(*(uint32_t*)(a1 + 16) + 124), (float2*)(v2 + 56)) & 2) {
		if (!*(uint32_t*)(a1 + 76)) {
			*(uint32_t*)(a1 + 76) = 3 * gameFPS();
		}
	} else {
		*(uint32_t*)(a1 + 76) = 0;
	}
	if (nox_xxx_calcDistance_4E6C00(*(uint32_t*)(a1 + 48), *(uint32_t*)(a1 + 16)) > 400.0 ||
		!nox_xxx_unitCanInteractWith_5370E0(*(uint32_t*)(a1 + 16), *(uint32_t*)(a1 + 48), 0)) {
		*(uint32_t*)(a1 + 48) = 0;
		goto LABEL_16;
	} else {
		goto LABEL_A;
	}
LABEL_16:
	v3 = *(uint32_t*)(*(uint32_t*)(*(uint32_t*)(a1 + 16) + 748) + 288);
	dword_5d4594_2487932 = 0;
	if (!v3) {
		*getMemU32Ptr(0x587000, 260404) = 1209810944;
		nox_xxx_unitsGetInCircle_517F90((float2*)(*(uint32_t*)(a1 + 16) + 56), 400.0, sub_531920, *(uint32_t*)(a1 + 16));
	} else {
		if (nox_xxx_unitIsEnemyTo_5330C0(*(uint32_t*)(a1 + 16), v3) &&
			nox_xxx_calcDistance_4E6C00(*(uint32_t*)(a1 + 16), v3) <= 400.0) {
			dword_5d4594_2487932 = v3;
		}
		if (!dword_5d4594_2487932) {
			*getMemU32Ptr(0x587000, 260404) = 1209810944;
			nox_xxx_unitsGetInCircle_517F90((float2*)(*(uint32_t*)(a1 + 16) + 56), 400.0, sub_531920, *(uint32_t*)(a1 + 16));
		}
	}
	*(uint32_t*)(a1 + 48) = dword_5d4594_2487932;
	*(uint32_t*)(a1 + 76) = 0;
LABEL_A:
	v4 = *(uint32_t*)(a1 + 76);
	if (v4) {
		v5 = v4 - 1;
		*(uint32_t*)(a1 + 76) = v5;
		if (!v5) {
			*(uint32_t*)(a1 + 48) = 0;
		}
	}
	v6 = *(uint32_t*)(a1 + 48);
	v7 = *(uint32_t*)(a1 + 36);
	if (v6) {
		if (v6 != v7) {
			if (v7) {
				nox_xxx_netStopRaySpell_4FEF90(a1, *(uint32_t**)(a1 + 36));
			}
			nox_xxx_netStartDurationRaySpell_4FF130(a1);
		}
		v8 = *(unsigned short*)(*(uint32_t*)(a1 + 48) + 4);
		if ((unsigned short)v8 == *getMemU32Ptr(0x5D4594, 2487936) || v8 == *getMemU32Ptr(0x5D4594, 2487940)) {
			v9 = nox_xxx_gamedataGetFloat_419D40("PlasmaDamageHecubah");
		} else {
			v9 = nox_xxx_gamedataGetFloat_419D40("PlasmaDamage");
		}
		v20 = v9;
		v10 = nox_float2int(v20);
		(*(void (**)(uint32_t, uint32_t, uint32_t, int, int))(*(uint32_t*)(a1 + 48) + 716))(
			*(uint32_t*)(a1 + 48), *(uint32_t*)(a1 + 16), 0, v10, 14);
		v11 = *(uint32_t*)(a1 + 48);
		if (*(uint32_t*)(v11 + 16) & 0x8020) {
			nox_xxx_netSendPointFx_522FF0(131, (float2*)(v11 + 56));
		}
		v19 = *(uint32_t**)(a1 + 16);
		*(uint32_t*)(a1 + 36) = *(uint32_t*)(a1 + 48);
		nox_xxx_playerSetState_4FA020(v19, 22);
		if (!(gameFrame() % (gameFPS() / 3u))) {
			nox_xxx_aud_501960(98, *(uint32_t*)(a1 + 16), 0, 0);
			nox_xxx_aud_501960(98, *(uint32_t*)(a1 + 48), 0, 0);
		}
		if (!*(uint32_t*)(a1 + 76)) {
			v21 = nox_xxx_gamedataGetFloat_419D40("PlasmaSearchTime");
			*(uint32_t*)(a1 + 68) = gameFrame() + nox_float2int(v21);
		}
		v12 = *(uint32_t*)(a1 + 72);
		if (v12) {
			v13 = *(uint32_t*)(v12 + 736);
			v14 = *(uint8_t*)(v13 + 108);
			if (v14 <= 0u) {
				return 1;
			}
			v15 = v14 - 1;
			*(uint8_t*)(v13 + 108) = v15;
			*(uint32_t*)(v13 + 112) = 100 * v15 / *(unsigned char*)(v13 + 109);
			v16 = *(uint32_t*)(a1 + 16);
			if (v16) {
				if (*(uint8_t*)(v16 + 8) & 4) {
					v17 = *(uint32_t*)(v16 + 748);
					nox_xxx_playerSetState_4FA020((uint32_t*)v16, 22);
					nox_xxx_netReportCharges_4D82B0(*(unsigned char*)(*(uint32_t*)(v17 + 276) + 2064),
													*(uint32_t**)(a1 + 72), *(uint8_t*)(v13 + 108),
													*(uint8_t*)(v13 + 109));
				}
			}
			if (*(uint8_t*)(v13 + 108) <= 0u) {
				return 1;
			}
		}
	} else if (v7) {
		nox_xxx_netStopRaySpell_4FEF90(a1, *(uint32_t**)(a1 + 36));
		*(uint32_t*)(a1 + 36) = 0;
	}
	return 0;
}

//----- (00531920) --------------------------------------------------------
void sub_531920(int a1, int a2) {
	int v2;    // eax
	int v3;    // eax
	double v4; // st7
	double v5; // st6
	double v6; // st5

	v2 = *(uint32_t*)(a1 + 8);
	if (v2 & 0x20006) {
		if (!(*(uint32_t*)(a1 + 16) & 0x8020) && a1 != a2) {
			if (!(v2 & 2) || (v3 = *(uint32_t*)(a1 + 12), (v3 & 0x8000) == 0)) {
				if (nox_xxx_unitIsEnemyTo_5330C0(a2, a1) &&
					nox_server_testTwoPointsAndDirection_4E6E50((float2*)(a2 + 56), *(short*)(a2 + 124),
																(float2*)(a1 + 56)) &
							1 |
						0xC &&
					nox_xxx_unitCanInteractWith_5370E0(a2, a1, 0)) {
					v4 = *(float*)(a1 + 56) - *(float*)(a2 + 56);
					v5 = *(float*)(a1 + 60) - *(float*)(a2 + 60);
					v6 = v5 * v5 + v4 * v4;
					if (v6 < *getMemFloatPtr(0x587000, 260404)) {
						*getMemFloatPtr(0x587000, 260404) = v6;
						dword_5d4594_2487932 = a1;
					}
				}
			}
		}
	}
}

//----- (005319E0) --------------------------------------------------------
int sub_5319E0(int a1) {
	int result; // eax

	result = *(uint32_t*)(a1 + 72);
	if (result) {
		result = *(uint32_t*)(result + 736);
		*(uint32_t*)(result + 96) &= 0xFFFFFFFB;
	}
	return result;
}

//----- (00531A00) --------------------------------------------------------
int nox_xxx_spellCreateMoonglow_531A00(uint32_t* a1) {
	short v1;     // di
	int v2;       // eax
	uint32_t* v3; // eax
	int v4;       // ecx
	int v5;       // edx
	float v7;     // [esp+0h] [ebp-14h]
	float v8;     // [esp+Ch] [ebp-8h]
	float v9;     // [esp+10h] [ebp-4h]

	v7 = nox_xxx_gamedataGetFloat_419D40("MoonglowEnchantmentDuration");
	v1 = nox_float2int(v7);
	v2 = a1[12];
	if (!v2 || *(uint32_t*)(v2 + 16) & 0x8020) {
		return 1;
	}
	if ((*(uint8_t*)(v2 + 8) & 4) != 4) {
		nox_xxx_buffApplyTo_4FF380(v2, 15, v1, a1[2]);
		return 1;
	}
	v3 = nox_xxx_newObjectByTypeID_4E3810("Moonglow");
	a1[18] = v3;
	if (!v3) {
		return 1;
	}
	v4 = a1[12];
	v5 = *(uint32_t*)(*(uint32_t*)(v4 + 748) + 276);
	if (*(uint8_t*)(v5 + 3680) & 0x10) {
		v8 = (double)*(int*)(v5 + 2284);
		v9 = (double)*(int*)(v5 + 2288);
	} else {
		v8 = 2944.0;
		v9 = 2944.0;
	}
	nox_xxx_createAt_4DAA50((int)v3, v4, v8, v9);
	nox_xxx_buffApplyTo_4FF380(a1[12], 1, v1, a1[2]);
	return 0;
}

//----- (00531AF0) --------------------------------------------------------
int sub_531AF0(int a1) {
	int result; // eax

	result = *(uint32_t*)(a1 + 48);
	if (result) {
		if ((*(uint8_t*)(result + 8) & 4) == 4) {
			if (*(uint32_t*)(a1 + 72)) {
				nox_xxx_delayedDeleteObject_4E5CC0(*(uint32_t*)(a1 + 72));
			}
			*(uint32_t*)(a1 + 72) = 0;
			result = nox_xxx_spellBuffOff_4FF5B0(*(uint32_t*)(a1 + 48), 1);
		} else {
			result = nox_xxx_spellBuffOff_4FF5B0(result, 15);
		}
	}
	return result;
}

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

//----- (00538290) --------------------------------------------------------
int nox_xxx_playerPreAttackEffects_538290(int a1, int a2, int a3, int a4) {
	int result;                          // eax
	int v5;                              // esi
	int v6;                              // ebp
	int* v7;                             // esi
	int v8;                              // eax
	void (*v9)(int, int, int, int, int); // ecx
	int v10;                             // [esp+10h] [ebp+4h]

	result = a3;
	if (a3) {
		v5 = *(uint32_t*)(a3 + 692);
		v6 = a1;
		if (nox_xxx_CheckGameplayFlags_417DA0(1) || !a2 || !(*(uint8_t*)(a2 + 8) & 6) ||
			(result = nox_xxx_unitIsEnemyTo_5330C0(a2, a1)) != 0) {
			result = nox_xxx_testUnitBuffs_4FF350(a1, 23);
			if (!result) {
				result = nox_xxx_testUnitBuffs_4FF350(a1, 27);
				if (!result) {
					v7 = (int*)(v5 + 8);
					v10 = 2;
					do {
						v8 = *v7;
						if (*v7) {
							v9 = *(void (**)(int, int, int, int, int))(v8 + 52);
							if (v9) {
								v9(v8, a3, a2, v6, a4);
							}
						}
						++v7;
						result = --v10;
					} while (v10);
				}
			}
		}
	}
	return result;
}

//----- (00538330) --------------------------------------------------------
int nox_xxx_playerTraceAttack_538330(int a1, int a2) {
	int v2;    // ebp
	int v3;    // eax
	int v4;    // ecx
	int v5;    // eax
	int v6;    // eax
	int v8;    // eax
	int v9;    // eax
	int v12;   // eax
	float a3;  // [esp+0h] [ebp-34h]
	float v15; // [esp+Ch] [ebp-28h]
	float v16; // [esp+Ch] [ebp-28h]
	float v17; // [esp+Ch] [ebp-28h]
	float v18; // [esp+Ch] [ebp-28h]
	float v19; // [esp+Ch] [ebp-28h]
	int v20;   // [esp+10h] [ebp-24h]
	int v21;   // [esp+20h] [ebp-14h]
	int4 a1a;  // [esp+24h] [ebp-10h]
	float v23; // [esp+38h] [ebp+4h]

	v2 = a1;
	v21 = 0;
	if (!a1 || !a2) {
		return 0;
	}
	dword_5d4594_2488656 = 0;
	dword_5d4594_2488660 = 0;
	v3 = *(uint32_t*)(a2 + 28);
	if (v3 && (v4 = *(uint32_t*)(v3 + 12), BYTE1(v4) & 0x40)) {
		sub_518040(a2 + 16, *(float*)(a2 + 8), sub_538510, a2);
		v21 = 25;
	} else {
		dword_5d4594_2488660 = 0;
		dword_5d4594_2488652 = *(uint32_t*)(a2 + 8);
		sub_518040(a1 + 56, *(float*)(a2 + 8), sub_5386A0, a1);
		if (dword_5d4594_2488660) {
			sub_538510(*(int*)&dword_5d4594_2488660, a2);
		}
	}
	v23 = (double)v21;
	v15 = *(float*)(a2 + 16) - *(float*)(a2 + 8) - v23;
	v5 = nox_float2int(v15);
	v16 = *(float*)(a2 + 20) - *(float*)(a2 + 8) - v23;
	a1a.field_0 = v5 / 23;
	v6 = nox_float2int(v16);
	v17 = *(float*)(a2 + 8) + *(float*)(a2 + 16) + v23;
	a1a.field_4 = v6 / 23;
	v8 = nox_float2int(v17);
	v18 = *(float*)(a2 + 20) + *(float*)(a2 + 8) + v23;
	a1a.field_8 = v8 / 23;
	v9 = nox_float2int(v18);
	v12 = *(uint32_t*)(a2 + 28);
	a1a.field_C = v9 / 23;
	if (!v12) {
		v12 = v2;
	}
	a3 = v23 + *(float*)(a2 + 8);
	nox_xxx_mapDamageToWalls_534FC0(&a1a, a2 + 16, a3, (long long)(*(float*)a2 + 0.5), *(unsigned char*)(a2 + 4), v12);
	if (*(uint32_t*)(a2 + 28)) {
		if (dword_5d4594_2488656) {
			v20 = *(unsigned char*)(a2 + 4);
			v19 = nox_xxx_gamedataGetFloat_419D40("ItemDamagePercentage") * *(float*)a2;
			nox_xxx_playerDamageWeapon_4E1560(*(uint32_t*)(a2 + 28), *(uint32_t*)(a2 + 12),
											  *(int*)&dword_5d4594_2488660, *(int*)&dword_5d4594_2488660, v19, v20);
		}
	}
	return dword_5d4594_2488656;
}

//----- (00538510) --------------------------------------------------------
void sub_538510(int a1, int a2) {
	int v2;     // eax
	int v3;     // eax
	int v4;     // eax
	float v5;   // ecx
	float v6;   // edx
	float v7;   // eax
	int v8;     // eax
	int v9;     // eax
	int v10;    // esi
	char v11;   // al
	float4 v12; // [esp+Ch] [ebp-10h]

	if (a1) {
		if (a2) {
			v2 = *(uint32_t*)(a2 + 12);
			if (v2) {
				if (a1 != v2) {
					if ((unsigned char)nox_server_testTwoPointsAndDirection_4E6E50(
							(float2*)(v2 + 56), *(short*)(v2 + 124), (float2*)(a1 + 56)) &
						*(uint8_t*)(a2 + 32)) {
						v3 = *(uint32_t*)(a1 + 16);
						if (!(v3 & 0x8040) && (*(uint32_t*)(a2 + 24) || !(v3 & 8))) {
							if (*(uint16_t*)(a1 + 24) != 0x4000) {
								dword_5d4594_2488656 = 1;
							}
							v4 = *(uint32_t*)(a2 + 12);
							v12.field_0 = *(float*)(v4 + 56);
							v5 = *(float*)(a1 + 60);
							v6 = *(float*)(v4 + 60);
							v7 = *(float*)(a1 + 56);
							v12.field_4 = v6;
							v12.field_8 = v7;
							v12.field_C = v5;
							if (nox_xxx_mapTraceRay_535250(&v12, 0, 0, 5)) {
								nox_xxx_playerPreAttackEffects_538290(a1, *(uint32_t*)(a2 + 12), *(uint32_t*)(a2 + 28),
																	  a2);
								(*(void (**)(int, uint32_t, uint32_t, uint32_t, uint32_t))(a1 + 716))(
									a1, *(uint32_t*)(a2 + 12), *(uint32_t*)(a2 + 28), (long long)(*(float*)a2 + 0.5),
									*(unsigned char*)(a2 + 4));
								if (nox_common_gameFlags_check_40A5C0(2048)) {
									if (*(uint8_t*)(*(uint32_t*)(a2 + 12) + 8) & 4) {
										if (!(*(uint8_t*)(a1 + 8) & 2)) {
											v8 = *(uint32_t*)(a1 + 556);
											if (v8) {
												if (*(uint16_t*)(v8 + 4)) {
													v9 = *(uint32_t*)(a1 + 16);
													if ((v9 & 0x8000) == 0 && !(v9 & 0x20)) {
														nox_xxx_netSendPointFx_522FF0(139, (float2*)(a1 + 56));
													}
												}
											}
										}
									}
								}
								if (!*(uint32_t*)(a2 + 28)) {
									v10 = *(uint32_t*)(a2 + 12);
									if (*(uint8_t*)(v10 + 8) & 4) {
										v11 = *(uint8_t*)(*(uint32_t*)(*(uint32_t*)(v10 + 748) + 276) + 8);
									} else {
										v11 = *(uint8_t*)(*(uint32_t*)(v10 + 748) + 2068);
									}
									if (v11 == 25) {
										nox_xxx_objectApplyForce_52DF80(v10 + 56, a1, 20.0);
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

//----- (005386A0) --------------------------------------------------------
void sub_5386A0(int a3, int a2) {
	int v2;    // eax
	int v3;    // eax
	double v4; // st7
	double v5; // st7
	float* v6; // eax
	double v7; // st7
	float2 a4; // [esp+8h] [ebp-8h]

	if (a2 != a3) {
		if (a2) {
			if (a3) {
				v2 = *(uint32_t*)(a3 + 16);
				if (!(v2 & 0x8049) && (*(uint8_t*)(a3 + 8) & 6 || !(v2 & 0x10) || (v2 & 0x80u) != 0) &&
					(nox_xxx_unitIsEnemyTo_5330C0(a2, a3) || !(*(uint8_t*)(a3 + 8) & 6))) {
					v3 = *(uint32_t*)(a3 + 16);
					if ((v3 & 0x8000) == 0) {
						if (nox_xxx_unitCanInteractWith_5370E0(a2, a3, 1)) {
							if (*(float*)&dword_5d4594_2488652 > 0.0) {
								a4.field_0 = *(float*)(a3 + 56) - *(float*)(a2 + 56);
								v4 = *(float*)(a3 + 60) - *(float*)(a2 + 60);
								a4.field_4 = v4;
								v5 = sqrt(v4 * a4.field_4 + a4.field_0 * a4.field_0);
								if (v5 != 0.0) {
									v6 = getMemFloatPtr(0x587000, 194136 + 8 * *(short*)(a2 + 124));
									if (a4.field_4 / v5 * v6[1] + a4.field_0 / v5 * *v6 > 0.5) {
										if (*(uint32_t*)(a3 + 172) == 2) {
											v5 = v5 - *(float*)(a3 + 176);
										} else if (*(uint32_t*)(a3 + 172) == 3) {
											v7 =
												sub_54A990((float2*)(a2 + 56), *(float*)&dword_5d4594_2488652, a3, &a4);
											if (v7 < 0.0) {
												return;
											}
											v5 = *(float*)&dword_5d4594_2488652 - v7;
										}
										if (v5 < 0.0) {
											v5 = 0.0;
										}
										if ((v5 < *(float*)&dword_5d4594_2488652 ||
											 dword_5d4594_2488660 && !(*(uint8_t*)(dword_5d4594_2488660 + 8) & 2) &&
												 (*(uint8_t*)(a3 + 8) & 2) == 2) &&
											(!dword_5d4594_2488660 || !(*(uint8_t*)(dword_5d4594_2488660 + 8) & 2))) {
											*(float*)&dword_5d4594_2488652 = v5;
											dword_5d4594_2488660 = a3;
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
}

//----- (00538840) --------------------------------------------------------
int nox_xxx_itemApplyAttackEffect_538840(int a1, int a2, int a3) {
	int v3;                                   // edi
	int* v4;                                  // esi
	int v5;                                   // eax
	void (*v6)(int, int, int, uint32_t, int); // ecx
	int result;                               // eax
	int v8;                                   // [esp+14h] [ebp+4h]

	v3 = a1;
	v8 = 4;
	v4 = *(int**)(v3 + 692);
	do {
		v5 = *v4;
		if (*v4) {
			v6 = *(void (**)(int, int, int, uint32_t, int))(v5 + 40);
			if (v6) {
				v6(v5, v3, a2, 0, a3);
			}
		}
		++v4;
		result = --v8;
	} while (v8);
	return result;
}

//----- (00538960) --------------------------------------------------------
void nox_xxx_castCounterSpell_52BBB0(int a1, int a2, int a3, int a4);
int nox_xxx_playerAttack_538960(nox_object_t* a1p) {
	int a1 = a1p;
	float* v1;         // edi
	int v2;            // eax
	int v3;            // eax
	int v4;            // ebp
	int v5;            // ebx
	char v6;           // dl
	float v7;          // edx
	float v8;          // eax
	char v9;           // al
	double v10;        // st7
	int v11;           // eax
	int v12;           // edx
	double v13;        // st7
	short v14;         // ax
	int v15;           // ecx
	int v16;           // edx
	double v17;        // st7
	int v18;           // eax
	int v19;           // eax
	char v20;          // cl
	int v21;           // ecx
	double v22;        // st7
	double v23;        // st7
	int v24;           // edi
	int v25;           // ebx
	int v26;           // edx
	double v27;        // st7
	double v28;        // st7
	double v29;        // st7
	int v30;           // eax
	float v31;         // edx
	double v32;        // st6
	uint32_t* v33;     // eax
	int v34;           // edi
	int v35;           // ebx
	short v36;         // ax
	double v37;        // st7
	int v38;           // eax
	float v39;         // edx
	char* v40;         // ebx
	double v41;        // st6
	uint32_t* v42;     // eax
	int v43;           // edi
	short v44;         // ax
	char v45;          // al
	int v46;           // edx
	double v47;        // st7
	double v48;        // st7
	int v49;           // edx
	double v50;        // st7
	double v51;        // st7
	int v52;           // edx
	double v53;        // st7
	double v54;        // st7
	int v55;           // eax
	int v56;           // eax
	int v57;           // edx
	double v58;        // st7
	double v59;        // st7
	int v60;           // edx
	double v61;        // st7
	double v62;        // st7
	uint32_t* v63;     // edi
	int v64;           // eax
	char v65;          // bl
	unsigned char v66; // bl
	uint32_t* v67;     // edi
	int v68;           // eax
	bool v69;          // cc
	int v70;           // eax
	int v71;           // ecx
	unsigned char v72; // al
	int v73;           // ecx
	float v75;         // [esp+0h] [ebp-C0h]
	short v76 = 0;     // [esp+17h] [ebp-A9h]
	int v77 = 0;       // [esp+1Ch] [ebp-A4h]
	int v78;           // [esp+20h] [ebp-A0h]
	int v79;           // [esp+24h] [ebp-9Ch]
	int v80;           // [esp+28h] [ebp-98h]
	int v81;           // [esp+2Ch] [ebp-94h]
	char v82[36];      // [esp+30h] [ebp-90h]
	float4 v83;
	int v86;      // [esp+64h] [ebp-5Ch]
	char v87[88]; // [esp+68h] [ebp-58h]

	v1 = 0;
	v2 = *(uint32_t*)(a1 + 8);
	if (v2 & 4) {
		v79 = *(uint32_t*)(a1 + 748);
		v3 = *(uint32_t*)(v79 + 276);
		v4 = *(uint32_t*)(v79 + 104);
		v5 = *(uint32_t*)(v3 + 4);
		v6 = *(uint8_t*)(v3 + 8);
		LOBYTE(v76) = *(uint8_t*)(v79 + 236);
	} else {
		if (!(v2 & 2) || !(*(uint8_t*)(a1 + 12) & 0x10)) {
			return 0;
		}
		v79 = *(uint32_t*)(a1 + 748);
		v5 = *(uint32_t*)(v79 + 2056);
		v4 = *(uint32_t*)(v79 + 2064);
		v6 = *(uint8_t*)(v79 + 2068);
		LOBYTE(v76) = *(uint8_t*)(v79 + 481);
	}
	LOBYTE(v80) = v6;
	if (v4) {
		v1 = (float*)nox_xxx_getProjectileClassById_413250(*(unsigned short*)(v4 + 4));
		if (!v1) {
			return 0;
		}
		*(uint32_t*)(v4 + 56) = *(uint32_t*)(a1 + 56);
		*(uint32_t*)(v4 + 60) = *(uint32_t*)(a1 + 60);
		*(uint32_t*)(v4 + 72) = *(uint32_t*)(a1 + 56);
		*(uint32_t*)(v4 + 76) = *(uint32_t*)(a1 + 60);
	} else if (!v6) {
		LOBYTE(v80) = nox_common_randomInt_415FA0(23, 24);
		if (!(*(uint8_t*)(a1 + 8) & 4)) {
			*(uint8_t*)(v79 + 2068) = v80;
		} else {
			if (!*(uint8_t*)(*(uint32_t*)(v79 + 276) + 2251) && nox_common_randomInt_415FA0(0, 100) >= 75) {
				LOBYTE(v80) = 25;
			}
			if (*(uint8_t*)(a1 + 8) & 4) {
				*(uint8_t*)(*(uint32_t*)(v79 + 276) + 8) = v80;
			} else {
				*(uint8_t*)(v79 + 2068) = v80;
			}
		}
	}
	*(uint32_t*)&v82[28] = v4;
	v81 = nox_xxx_unitGetStrength_4F9FD0(a1);
	if (nox_common_playerIsAbilityActive_4FC250(a1, 2) && nox_xxx_probablyWarcryCheck_4FC3E0(a1, 2)) {
		nox_xxx_animPlayerGetFrameRange_4F9F90(46, &v77, &v78);
		HIBYTE(v76) = (gameFrame() - *(uint32_t*)(a1 + 136)) / (unsigned int)(v78 + 1);
		if (v76 == 770) {
			v7 = *(float*)(a1 + 56);
			v8 = *(float*)(a1 + 60);
			LODWORD(v83.field_0) = a1;
			v83.field_4 = v7;
			v83.field_8 = v8;
			nox_xxx_earthquakeSend_4D9110((float*)(a1 + 56), 15);
			nox_xxx_unitsGetInCircle_517F90((float2*)(a1 + 56), 300.0, nox_xxx_warcryStunMonsters_539B90, a1);
			nox_xxx_castCounterSpell_52BBB0(13, a1, a1, a1);
		}
		if (HIBYTE(v76) >= v77) {
			sub_4FC440(a1, 2);
		}
		goto LABEL_159;
	}
	if (nox_common_playerIsAbilityActive_4FC250(a1, 1)) {
		if (!nox_xxx_testUnitBuffs_4FF350(a1, 25) && !nox_xxx_testUnitBuffs_4FF350(a1, 5)) {
			nox_xxx_animPlayerGetFrameRange_4F9F90(45, &v77, &v78);
			v9 = (gameFrame() - *(uint32_t*)(a1 + 136)) / (unsigned int)(v78 + 1);
			v10 = *(float*)(a1 + 548) * 6.0;
			*(float*)(a1 + 544) = v10;
			HIBYTE(v76) = v9;
			v11 = 8 * *(short*)(a1 + 124);
			v12 = v77 - 1;
			*(float*)(a1 + 88) = v10 * *getMemFloatPtr(0x587000, 194136 + v11) + *(float*)(a1 + 88);
			*(float*)(a1 + 92) = v10 * *getMemFloatPtr(0x587000, 194140 + v11) + *(float*)(a1 + 92);
			if (HIBYTE(v76) >= v12) {
				HIBYTE(v76) = 0;
			}
			goto LABEL_159;
		}
		return 0;
	}
	if (!v4) {
		nox_xxx_animPlayerGetFrameRange_4F9F90((unsigned char)v80, &v77, &v78);
		if (*(uint8_t*)(a1 + 8) & 4 && !*(uint32_t*)v79) {
			*(uint32_t*)v79 = gameFrame() + v77 * (v78 + 1);
		}
		HIBYTE(v76) = (gameFrame() - *(uint32_t*)(a1 + 136)) / (unsigned int)(v78 + 1);
		if (HIBYTE(v76) < v77) {
			goto LABEL_159;
		}
		if ((unsigned char)v80 < 0x17u) {
			v13 = 0.0;
			v14 = 0;
		} else if ((unsigned char)v80 <= 0x18u) {
			v13 = 0.039999999;
			v14 = 5;
		} else if ((uint8_t)v80 == 25) {
			v13 = 0.039999999;
			v14 = 10;
		} else {
			v13 = 0.0;
			v14 = 0;
		}
		v15 = *(uint32_t*)(a1 + 56);
		v16 = *(uint32_t*)(a1 + 60);
		*(float*)&v87[64] = v13;
		*(uint16_t*)&v87[72] = v14;
		*(uint32_t*)&v82[16] = v15;
		*(uint32_t*)&v82[20] = v16;
		*(uint16_t*)&v87[60] = 0;
		*(float*)v82 = nox_xxx_calcBoltDamage_4EF1E0(v81, (int)v87);
		v17 = *(float*)(a1 + 176) + 20.0;
		v82[4] = 10;
		*(uint32_t*)&v82[12] = a1;
		*(uint32_t*)&v82[24] = 0;
		*(float*)&v82[8] = v17;
		v82[32] = 1;
		*(uint32_t*)&v82[28] = 0;
		v18 = nox_xxx_playerTraceAttack_538330(a1, (int)v82);
		if (!v18) {
			nox_xxx_aud_501960(879, a1, 0, 0);
		}
		goto LABEL_159;
	}
	if (v5 & 0x47F8000) {
		v19 = *(uint32_t*)(v4 + 736);
		v86 = *(uint32_t*)(v4 + 736);
		if ((v5 & 0x8000) != 0 || (v20 = *(uint8_t*)(v19 + 96), v80 = 31, v20 & 2)) {
			v80 = 29;
		}
		nox_xxx_animPlayerGetFrameRange_4F9F90(v80, &v77, &v78);
		if (*(uint8_t*)(a1 + 8) & 4 && !*(uint32_t*)v79) {
			if (v80 == 29) {
				*(uint32_t*)v79 = gameFrame() + v77 * (v78 + 1);
			} else {
				*(uint32_t*)v79 = gameFrame();
			}
		}
		HIBYTE(v76) = (gameFrame() - *(uint32_t*)(a1 + 136)) / (unsigned int)(v78 + 1);
		if (HIBYTE(v76) >= v77) {
			if (v5 & 0x47F0000) {
				*(uint32_t*)(v86 + 96) &= 0xFFFFFFFD;
			}
			goto LABEL_159;
		}
		if (v80 != 29 || HIBYTE(v76) != v77 / 2 || HIBYTE(v76) <= (unsigned char)v76) {
			if (*(uint8_t*)(a1 + 8) & 2) {
				if (v80 != 29 && v76 == 256) {
					v24 = *(uint32_t*)(v4 + 736);
					nox_xxx_useByNetCode_53F8E0(a1, v4);
					if (!*(uint8_t*)(v24 + 108)) {
						if (*(uint8_t*)(v24 + 109)) {
							nox_xxx_equipWeaponNPC_53A030(a1, v4);
						}
					}
				}
			}
			goto LABEL_159;
		}
		v21 = *(uint32_t*)(a1 + 60);
		*(uint32_t*)&v82[16] = *(uint32_t*)(a1 + 56);
		*(uint32_t*)&v82[20] = v21;
		*(float*)v82 = nox_xxx_calcBoltDamage_4EF1E0(v81, (int)v1);
		v22 = *(float*)(a1 + 176);
		v82[4] = 0;
		*(uint32_t*)&v82[12] = a1;
		v23 = v22 + v1[17];
		*(uint32_t*)&v82[24] = 0;
		v82[32] = 1;
		*(float*)&v82[8] = v23;
		nox_xxx_itemApplyAttackEffect_538840(v4, a1, (int)v82);
		v18 = nox_xxx_playerTraceAttack_538330(a1, (int)v82);
		if (!v18) {
			nox_xxx_aud_501960(879, a1, 0, 0);
		}
		goto LABEL_159;
	}
	if (v5 & 0x7800000) {
		v25 = ((v5 & 0x3800000) != 0) + 31;
		nox_xxx_animPlayerGetFrameRange_4F9F90(v25, &v77, &v78);
		if (*(uint8_t*)(a1 + 8) & 4 && !*(uint32_t*)v79) {
			*(uint32_t*)v79 = gameFrame() + v77 * (v78 + 1);
		}
		HIBYTE(v76) = (gameFrame() - *(uint32_t*)(a1 + 136)) / (unsigned int)(v78 + 1);
		if (v25 != 32 || HIBYTE(v76) != v77 / 2 || HIBYTE(v76) <= (unsigned char)v76) {
			goto LABEL_159;
		}
		v26 = *(uint32_t*)(a1 + 60);
		*(uint32_t*)&v82[16] = *(uint32_t*)(a1 + 56);
		*(uint32_t*)&v82[20] = v26;
		*(float*)v82 = nox_xxx_calcBoltDamage_4EF1E0(v81, (int)v1);
		v27 = *(float*)(a1 + 176);
		v82[4] = 0;
		*(uint32_t*)&v82[12] = a1;
		v28 = v27 + v1[17];
		*(uint32_t*)&v82[24] = 0;
		v82[32] = 1;
		*(float*)&v82[8] = v28;
		nox_xxx_itemApplyAttackEffect_538840(v4, a1, (int)v82);
		v18 = nox_xxx_playerTraceAttack_538330(a1, (int)v82);
		if (!v18) {
			nox_xxx_aud_501960(879, a1, 0, 0);
		}
		goto LABEL_159;
	}
	if (!(v5 & 0x40)) {
		if ((v5 & 0x80u) == 0) {
			if (v5 & 0x200) {
				nox_xxx_animPlayerGetFrameRange_4F9F90(28, &v77, &v78);
				if (*(uint8_t*)(a1 + 8) & 4 && !*(uint32_t*)v79) {
					*(uint32_t*)v79 = gameFrame() + v77 * (v78 + 1);
				}
				HIBYTE(v76) = (gameFrame() - *(uint32_t*)(a1 + 136)) / (unsigned int)(v78 + 1);
				if (HIBYTE(v76) == v77 / 2 && HIBYTE(v76) > (unsigned char)v76) {
					v46 = *(uint32_t*)(a1 + 60);
					*(uint32_t*)&v82[16] = *(uint32_t*)(a1 + 56);
					*(uint32_t*)&v82[20] = v46;
					*(float*)v82 = nox_xxx_calcBoltDamage_4EF1E0(v81, (int)v1);
					v47 = *(float*)(a1 + 176);
					v82[4] = 0;
					*(uint32_t*)&v82[12] = a1;
					v48 = v47 + v1[17];
					*(uint32_t*)&v82[24] = 0;
					v82[32] = 1;
					*(float*)&v82[8] = v48;
					nox_xxx_itemApplyAttackEffect_538840(v4, a1, (int)v82);
					if (!nox_xxx_playerTraceAttack_538330(a1, (int)v82)) {
						nox_xxx_aud_501960(880, a1, 0, 0);
					}
				}
			} else if (v5 & 0x100) {
				nox_xxx_animPlayerGetFrameRange_4F9F90(27, &v77, &v78);
				if (*(uint8_t*)(a1 + 8) & 4 && !*(uint32_t*)v79) {
					*(uint32_t*)v79 = gameFrame() + v77 * (v78 + 1);
				}
				HIBYTE(v76) = (gameFrame() - *(uint32_t*)(a1 + 136)) / (unsigned int)(v78 + 1);
				if (HIBYTE(v76) == v77 / 2 && HIBYTE(v76) > (unsigned char)v76) {
					v49 = *(uint32_t*)(a1 + 60);
					*(uint32_t*)&v82[16] = *(uint32_t*)(a1 + 56);
					*(uint32_t*)&v82[20] = v49;
					*(float*)v82 = nox_xxx_calcBoltDamage_4EF1E0(v81, (int)v1);
					v50 = *(float*)(a1 + 176);
					v82[4] = 0;
					*(uint32_t*)&v82[12] = a1;
					v51 = v50 + v1[17];
					*(uint32_t*)&v82[24] = 0;
					v82[32] = 1;
					*(float*)&v82[8] = v51;
					nox_xxx_itemApplyAttackEffect_538840(v4, a1, (int)v82);
					if (!nox_xxx_playerTraceAttack_538330(a1, (int)v82)) {
						nox_xxx_aud_501960(881, a1, 0, 0);
					}
				}
			} else if (v5 & 0x400) {
				nox_xxx_animPlayerGetFrameRange_4F9F90(37, &v77, &v78);
				if (*(uint8_t*)(a1 + 8) & 4 && !*(uint32_t*)v79) {
					*(uint32_t*)v79 = gameFrame() + v77 * (v78 + 1);
				}
				HIBYTE(v76) = (gameFrame() - *(uint32_t*)(a1 + 136)) / (unsigned int)(v78 + 1);
				if (HIBYTE(v76) == v77 / 2 && HIBYTE(v76) > (unsigned char)v76) {
					v52 = *(uint32_t*)(a1 + 60);
					*(uint32_t*)&v82[16] = *(uint32_t*)(a1 + 56);
					*(uint32_t*)&v82[20] = v52;
					*(float*)v82 = nox_xxx_calcBoltDamage_4EF1E0(v81, (int)v1);
					v53 = *(float*)(a1 + 176);
					v82[4] = 0;
					*(uint32_t*)&v82[12] = a1;
					v54 = v53 + v1[17];
					*(uint32_t*)&v82[24] = 0;
					v82[32] = 1;
					*(float*)&v82[8] = v54;
					nox_xxx_itemApplyAttackEffect_538840(v4, a1, (int)v82);
					if (!nox_xxx_playerTraceAttack_538330(a1, (int)v82)) {
						nox_xxx_aud_501960(881, a1, 0, 0);
					}
				}
			} else if (v5 & 0x4000) {
				nox_xxx_animPlayerGetFrameRange_4F9F90(39, &v77, &v78);
				if (*(uint8_t*)(a1 + 8) & 4 && !*(uint32_t*)v79) {
					*(uint32_t*)v79 = gameFrame() + v77 * (v78 + 1);
				}
				HIBYTE(v76) = (gameFrame() - *(uint32_t*)(a1 + 136)) / (unsigned int)(v78 + 1);
				if (HIBYTE(v76) == v77 / 2 && HIBYTE(v76) > (unsigned char)v76) {
					v55 = 8 * *(short*)(a1 + 124);
					*(float*)&v82[16] = *getMemFloatPtr(0x587000, 194136 + v55) * 35.0 + *(float*)(a1 + 56);
					*(float*)&v82[20] = *getMemFloatPtr(0x587000, 194140 + v55) * 35.0 + *(float*)(a1 + 60);
					*(float*)v82 = nox_xxx_calcBoltDamage_4EF1E0(v81, (int)v1);
					v82[4] = 2;
					*(uint32_t*)&v82[12] = a1;
					*(float*)&v82[8] = v1[17];
					*(uint32_t*)&v82[24] = 1;
					v82[32] = 1;
					nox_xxx_itemApplyAttackEffect_538840(v4, a1, (int)v82);
					nox_xxx_playerTraceAttack_538330(a1, (int)v82);
					v75 = (double)v81 * 0.1;
					v56 = nox_float2int(v75);
					nox_xxx_earthquakeSend_4D9110((float*)(a1 + 56), v56);
					nox_xxx_aud_501960(882, a1, 0, 0);
				}
			} else if (v5 & 0x800) {
				nox_xxx_animPlayerGetFrameRange_4F9F90(26, &v77, &v78);
				if (*(uint8_t*)(a1 + 8) & 4 && !*(uint32_t*)v79) {
					*(uint32_t*)v79 = gameFrame() + v77 * (v78 + 1);
				}
				HIBYTE(v76) = (gameFrame() - *(uint32_t*)(a1 + 136)) / (unsigned int)(v78 + 1);
				if (HIBYTE(v76) == v77 / 2 && HIBYTE(v76) > (unsigned char)v76) {
					v57 = *(uint32_t*)(a1 + 60);
					*(uint32_t*)&v82[16] = *(uint32_t*)(a1 + 56);
					*(uint32_t*)&v82[20] = v57;
					*(float*)v82 = nox_xxx_calcBoltDamage_4EF1E0(v81, (int)v1);
					v58 = *(float*)(a1 + 176);
					v82[4] = 2;
					*(uint32_t*)&v82[12] = a1;
					v59 = v58 + v1[17];
					*(uint32_t*)&v82[24] = 1;
					v82[32] = 1;
					*(float*)&v82[8] = v59;
					nox_xxx_itemApplyAttackEffect_538840(v4, a1, (int)v82);
					if (!nox_xxx_playerTraceAttack_538330(a1, (int)v82)) {
						nox_xxx_aud_501960(884, a1, 0, 0);
					}
				}
			} else if (v5 & 0x3000) {
				nox_xxx_animPlayerGetFrameRange_4F9F90(35, &v77, &v78);
				if (*(uint8_t*)(a1 + 8) & 4 && !*(uint32_t*)v79) {
					*(uint32_t*)v79 = gameFrame() + v77 * (v78 + 1);
				}
				HIBYTE(v76) = (gameFrame() - *(uint32_t*)(a1 + 136)) / (unsigned int)(v78 + 1);
				if (HIBYTE(v76) == v77 / 2 && HIBYTE(v76) > (unsigned char)v76) {
					v60 = *(uint32_t*)(a1 + 60);
					*(uint32_t*)&v82[16] = *(uint32_t*)(a1 + 56);
					*(uint32_t*)&v82[20] = v60;
					*(float*)v82 = nox_xxx_calcBoltDamage_4EF1E0(v81, (int)v1);
					v61 = *(float*)(a1 + 176);
					v82[4] = 0;
					*(uint32_t*)&v82[12] = a1;
					v62 = v61 + v1[17];
					*(uint32_t*)&v82[24] = 1;
					v82[32] = 1;
					*(float*)&v82[8] = v62;
					nox_xxx_itemApplyAttackEffect_538840(v4, a1, (int)v82);
					if (!nox_xxx_playerTraceAttack_538330(a1, (int)v82)) {
						nox_xxx_aud_501960(883, a1, 0, 0);
					}
				}
			} else if (v5 & 4) {
				nox_xxx_animPlayerGetFrameRange_4F9F90(33, &v77, &v78);
				if (*(uint8_t*)(a1 + 8) & 4) {
					v63 = (uint32_t*)v79;
					if (!*(uint32_t*)v79) {
						v64 = nox_xxx_itemCheckReadinessEffect_4E0960(v4);
						*v63 = gameFrame() + v77 * (v78 + 1) - v64;
					}
				}
				v65 = (gameFrame() - *(uint32_t*)(a1 + 136)) / (unsigned int)(v78 + 1);
				v66 = nox_xxx_itemCheckReadinessEffect_4E0960(v4) + v65;
				HIBYTE(v76) = v66;
				if (v66 >= v77 - 1 && v66 > (unsigned char)v76) {
					nox_xxx_shootBowCrossbow1_539BD0(a1, v4);
					HIBYTE(v76) = v77;
				}
			} else if (v5 & 8) {
				nox_xxx_animPlayerGetFrameRange_4F9F90(34, &v77, &v78);
				if (*(uint8_t*)(a1 + 8) & 4) {
					v67 = (uint32_t*)v79;
					if (!*(uint32_t*)v79) {
						v68 = nox_xxx_itemCheckReadinessEffect_4E0960(v4);
						*v67 = gameFrame() + v77 * (v78 + 1) - v68;
					}
				}
				HIBYTE(v76) = (gameFrame() - *(uint32_t*)(a1 + 136)) / (unsigned int)(v78 + 1);
				v69 = (int)HIBYTE(v76) < 1;
				if (HIBYTE(v76) == 1) {
					if (!(uint8_t)v76) {
						nox_xxx_shootBowCrossbow1_539BD0(a1, v4);
					}
					v69 = 0;
				}
				if (!v69) {
					HIBYTE(v76) += nox_xxx_itemCheckReadinessEffect_4E0960(v4);
				}
			}
			goto LABEL_159;
		}
		nox_xxx_animPlayerGetFrameRange_4F9F90(44, &v77, &v78);
		if (*(uint8_t*)(a1 + 8) & 4 && !*(uint32_t*)v79) {
			*(uint32_t*)v79 = gameFrame() + v77 * (v78 + 1);
		}
		HIBYTE(v76) = (gameFrame() - *(uint32_t*)(a1 + 136)) / (unsigned int)(v78 + 1);
		if (HIBYTE(v76) != v77 / 2 || HIBYTE(v76) <= (unsigned char)v76) {
			goto LABEL_159;
		}
		v37 = *(float*)(a1 + 176) + 4.0;
		v38 = 8 * *(short*)(a1 + 124);
		v39 = *(float*)(a1 + 60);
		v40 = *(char**)(v4 + 736);
		v41 = v37 * *getMemFloatPtr(0x587000, 194136 + v38) + *(float*)(a1 + 56);
		v83.field_0 = *(float*)(a1 + 56);
		v83.field_4 = v39;
		v83.field_8 = v41;
		v83.field_C = v37 * *getMemFloatPtr(0x587000, 194140 + v38) + *(float*)(a1 + 60);
		if (nox_xxx_mapTraceRay_535250(&v83, 0, 0, 4)) {
			v42 = nox_xxx_newObjectByTypeID_4E3810("FanChakramInMotion");
			v43 = (int)v42;
			if (!v42) {
				return 0;
			}
			*(uint32_t*)(v42[175] + 4) = a1;
			nox_xxx_createAt_4DAA50((int)v42, a1, v83.field_8, v83.field_C);
			nox_xxx_modifSetItemAttrs_4E4990(v43, *(int**)(v4 + 692));
			*(float*)(v43 + 80) =
				*getMemFloatPtr(0x587000, 194136 + 8 * *(short*)(a1 + 124)) * *(float*)(v43 + 544);
			*(float*)(v43 + 84) =
				*getMemFloatPtr(0x587000, 194140 + 8 * *(short*)(a1 + 124)) * *(float*)(v43 + 544);
			v44 = *(uint16_t*)(a1 + 124);
			*(uint16_t*)(v43 + 124) = v44;
			*(uint16_t*)(v43 + 126) = v44;
			nox_xxx_aud_501960(891, a1, 0, 0);
			if (!v40[2]) {
				v45 = v40[1] - 1;
				v40[1] = v45;
				if (v45) {
					if (*(uint8_t*)(a1 + 8) & 4) {
						nox_xxx_netReportCharges_4D82B0(*(unsigned char*)(*(uint32_t*)(v79 + 276) + 2064),
														(uint32_t*)v4, v45, *v40);
					}
				} else {
					sub_4ED0C0(a1, (int*)v4);
					nox_xxx_delayedDeleteObject_4E5CC0(v4);
					sub_539FB0((uint32_t*)a1);
				}
			}
			goto LABEL_159;
		}
		nox_xxx_aud_501960(323, a1, 0, 0);
		goto LABEL_159;
	}
	nox_xxx_animPlayerGetFrameRange_4F9F90(44, &v77, &v78);
	if (*(uint8_t*)(a1 + 8) & 4 && !*(uint32_t*)v79) {
		*(uint32_t*)v79 = gameFrame() + v77 * (v78 + 1);
	}
	HIBYTE(v76) = (gameFrame() - *(uint32_t*)(a1 + 136)) / (unsigned int)(v78 + 1);
	if (HIBYTE(v76) != v77 / 2 || HIBYTE(v76) <= (unsigned char)v76) {
		goto LABEL_159;
	}
	v29 = *(float*)(a1 + 176) + 4.0;
	v30 = 8 * *(short*)(a1 + 124);
	v31 = *(float*)(a1 + 60);
	v32 = v29 * *getMemFloatPtr(0x587000, 194136 + v30) + *(float*)(a1 + 56);
	v83.field_0 = *(float*)(a1 + 56);
	v83.field_4 = v31;
	v83.field_8 = v32;
	v83.field_C = v29 * *getMemFloatPtr(0x587000, 194140 + v30) + *(float*)(a1 + 60);
	if (!nox_xxx_mapTraceRay_535250(&v83, 0, 0, 5)) {
		nox_xxx_aud_501960(323, a1, 0, 0);
		goto LABEL_159;
	}
	v33 = nox_xxx_newObjectByTypeID_4E3810("RoundChakramInMotion");
	v34 = (int)v33;
	if (!v33) {
		return 0;
	}
	v35 = v33[187];
	sub_4ED0C0(a1, (int*)v4);
	nox_xxx_createAt_4DAA50(v34, a1, v83.field_8, v83.field_C);
	nox_xxx_inventoryPutImpl_4F3070(v34, v4, 1);
	nox_xxx_modifSetItemAttrs_4E4990(v34, *(int**)(v4 + 692));
	*(float*)(v34 + 80) = *getMemFloatPtr(0x587000, 194136 + 8 * *(short*)(a1 + 124)) * *(float*)(v34 + 544);
	*(float*)(v34 + 84) = *getMemFloatPtr(0x587000, 194140 + 8 * *(short*)(a1 + 124)) * *(float*)(v34 + 544);
	v36 = *(uint16_t*)(a1 + 124);
	*(uint16_t*)(v34 + 124) = v36;
	*(uint16_t*)(v34 + 126) = v36;
	*(uint8_t*)(v35 + 4) = 4;
	*(uint32_t*)(v35 + 16) = *(uint32_t*)(a1 + 56);
	*(uint32_t*)(v35 + 20) = *(uint32_t*)(a1 + 60);
	*(uint8_t*)(v35 + 24) = 2;
	nox_xxx_aud_501960(891, a1, 0, 0);
LABEL_159:
	v70 = *(uint32_t*)(a1 + 8);
	if (v70 & 4) {
		v71 = v79;
		v72 = HIBYTE(v76);
		*(uint8_t*)(v79 + 236) = HIBYTE(v76);
		if (HIBYTE(v76) >= v77) {
			*(uint8_t*)(v71 + 236) = v77 - 1;
		}
	} else if (v70 & 2 && *(uint8_t*)(a1 + 12) & 0x10) {
		v73 = v79;
		*(uint8_t*)(v79 + 481) = HIBYTE(v76);
		v72 = HIBYTE(v76);
		if (HIBYTE(v76) >= v77) {
			*(uint8_t*)(v73 + 481) = v77 - 1;
		}
	} else {
		v72 = HIBYTE(v76);
	}
	return v72 < v77;
}
// 539AF2: variable 'v76' is possibly undefined

//----- (00539B90) --------------------------------------------------------
short nox_xxx_warcryStunMonsters_539B90(int a1, int a2) {
	short result; // ax

	result = a2;
	if (a2) {
		result = a1;
		if (a1) {
			if (*(uint8_t*)(a1 + 8) & 2 && *(uint32_t*)(a1 + 12) & 0x20000 && !(*(uint32_t*)(a1 + 16) & 0x8020)) {
				nox_xxx_buffApplyTo_4FF380(a1, 5, 90, 3);
			}
		}
	}
	return result;
}

//----- (00539BD0) --------------------------------------------------------
int nox_xxx_shootBowCrossbow1_539BD0(int a1, int a2) {
	int v2;       // eax
	uint8_t* v3;  // edi
	char* v4;     // ebx
	uint32_t* v6; // esi
	int v7;       // eax
	int v8;       // eax

	v2 = nox_xxx_weaponInventoryEquipFlags_415820(a2);
	v3 = *(uint8_t**)(a2 + 736);
	v4 = (char*)v2;
	if (!*v3) {
		v6 = *(uint32_t**)(a1 + 504);
		if (v6) {
			while (1) {
				v7 = v6[4];
				if (v7 & 0x100) {
					if (nox_xxx_weaponInventoryEquipFlags_415820((int)v6) == 2) {
						v8 = v6[184];
						if (*(uint8_t*)(v8 + 1) || *(uint8_t*)(v8 + 2) == 1) {
							break;
						}
					}
				}
				v6 = (uint32_t*)v6[124];
				if (!v6) {
					goto LABEL_14;
				}
			}
			nox_xxx_shootBowCrossbow2_539D80(a1, (int)v6, a2, v4);
			return 1;
		}
	LABEL_14:
		if (nox_common_gameFlags_check_40A5C0(4096) && v4 == (char*)4) {
			nox_xxx_shootBowCrossbow2_539D80(a1, 0, a2, (char*)4);
		}
		if (!(*(uint8_t*)(a1 + 8) & 4)) {
			return 0;
		}
		if (nox_xxx_playerTryReloadQuiver_539FF0((uint32_t*)a1) == 1) {
			nox_xxx_netPriMsgToPlayer_4DA2C0(a1, "pattack.c:ReloadQuiver", 0);
			*v3 = 0;
		} else {
			if (nox_common_gameFlags_check_40A5C0(4096) && v4 == (char*)4) {
				goto LABEL_25;
			}
			nox_xxx_netPriMsgToPlayer_4DA2C0(a1, "pattack.c:NoQuiver", 0);
		}
		if (v4 != (char*)4) {
			if (v4 == (char*)8) {
				nox_xxx_aud_501960(888, a1, 0, 0);
			}
			goto LABEL_29;
		}
	LABEL_25:
		if (!nox_common_gameFlags_check_40A5C0(4096)) {
			nox_xxx_aud_501960(887, a1, 0, 0);
		}
	LABEL_29:
		nox_xxx_playerSetState_4FA020((uint32_t*)a1, 13);
		return 0;
	}
	nox_xxx_netPriMsgToPlayer_4DA2C0(a1, "pattack.c:ReloadingQuiver", 0);
	if (v4 == (char*)4) {
		if (!nox_common_gameFlags_check_40A5C0(4096)) {
			nox_xxx_aud_501960(887, a1, 0, 0);
			--*v3;
			return 0;
		}
	} else if (v4 == (char*)8) {
		nox_xxx_aud_501960(888, a1, 0, 0);
	}
	--*v3;
	return 0;
}

//----- (00539D80) --------------------------------------------------------
uint32_t* nox_xxx_shootBowCrossbow2_539D80(int a1, int a2, int a3, char* a4) {
	double v4;        // st7
	char* v5;         // ebx
	int v6;           // eax
	double v7;        // st6
	double v8;        // st7
	float v9;         // eax
	uint32_t* result; // eax
	uint32_t* v11;    // eax
	int v12;          // edi
	short v13;        // ax
	int v14;          // edi
	char v15;         // al
	char v16;         // cl
	float v17;        // [esp+Ch] [ebp-18h]
	float v18;        // [esp+10h] [ebp-14h]
	float4 v19;       // [esp+14h] [ebp-10h]

	v4 = *(float*)(a1 + 176) + 4.0;
	if (a2) {
		v5 = *(char**)(a2 + 736);
	} else {
		v5 = a4;
	}
	v6 = 8 * *(short*)(a1 + 124);
	v7 = v4 * *getMemFloatPtr(0x587000, 194136 + v6);
	v19.field_4 = *(float*)(a1 + 60);
	v17 = v7 + *(float*)(a1 + 56);
	v8 = v4 * *getMemFloatPtr(0x587000, 194140 + v6);
	v9 = *(float*)(a1 + 56);
	v19.field_8 = v17;
	v19.field_0 = v9;
	v18 = v8 + *(float*)(a1 + 60);
	v19.field_C = v18;
	result = (uint32_t*)nox_xxx_mapTraceRay_535250(&v19, 0, 0, 5);
	if (result) {
		if (a2) {
			if (a4 == (char*)4) {
				v11 = nox_xxx_newObjectByTypeID_4E3810("ArcherArrow");
			} else {
				v11 = nox_xxx_newObjectByTypeID_4E3810("ArcherBolt");
			}
		} else {
			v11 = nox_xxx_newObjectByTypeID_4E3810("WeakArcherArrow");
		}
		v12 = (int)v11;
		if (v11) {
			*(uint32_t*)(v11[175] + 4) = a1;
			nox_xxx_createAt_4DAA50((int)v11, a1, v17, v18);
			if (a2) {
				nox_xxx_modifSetItemAttrs_4E4990(v12, *(int**)(a2 + 692));
			}
			nox_xxx_shootApplyEffects_539F40(a1, a3, v12);
			*(float*)(v12 + 80) = *getMemFloatPtr(0x587000, 194136 + 8 * *(short*)(a1 + 124)) * *(float*)(v12 + 544);
			*(float*)(v12 + 84) = *getMemFloatPtr(0x587000, 194140 + 8 * *(short*)(a1 + 124)) * *(float*)(v12 + 544);
			v13 = *(uint16_t*)(a1 + 124);
			*(uint16_t*)(v12 + 124) = v13;
			*(uint16_t*)(v12 + 126) = v13;
		}
		if (a2) {
			if (!v5[2]) {
				if (*(uint8_t*)(a1 + 8) & 4) {
					v14 = *(uint32_t*)(a1 + 748);
					v15 = v5[1] - 1;
					v16 = *v5;
					v5[1] = v15;
					nox_xxx_netReportCharges_4D82B0(*(unsigned char*)(*(uint32_t*)(v14 + 276) + 2064), (uint32_t*)a2,
													v15, v16);
					if (!v5[1]) {
						nox_xxx_delayedDeleteObject_4E5CC0(a2);
					}
				}
			}
		}
		if (a4 == (char*)4) {
			nox_xxx_aud_501960(885, a1, 0, 0);
		} else {
			nox_xxx_aud_501960(886, a1, 0, 0);
		}
	}
	return result;
}

//----- (00539F40) --------------------------------------------------------
int nox_xxx_shootApplyEffects_539F40(int a1, int a2, int a3) {
	int v3;                             // ebx
	int v4;                             // edi
	int* v5;                            // esi
	int v6;                             // eax
	int (*v7)(int, int, int, int, int); // ecx
	int result;                         // eax
	int v9;                             // [esp+18h] [ebp+8h]
	int v10;                            // [esp+1Ch] [ebp+Ch]

	v3 = a3;
	v4 = a2;
	v10 = *(uint32_t*)(a3 + 692);
	v9 = 2;
	v5 = (int*)(*(uint32_t*)(v4 + 692) + 8);
	do {
		v6 = *v5;
		if (*v5) {
			if (*(char (**)(int, int, int, int))(v6 + 52) == nox_xxx_recoilEffect_4E0640) {
				*(uint32_t*)(v10 + 12) = v6;
			} else {
				v7 = *(int (**)(int, int, int, int, int))(v6 + 40);
				if (v7 == nox_xxx_effectProjectileSpeed_4E09B0) {
					v7(v6, v4, a1, 0, v3);
				}
			}
		}
		++v5;
		result = --v9;
	} while (v9);
	return result;
}

//----- (00539FB0) --------------------------------------------------------
int sub_539FB0(uint32_t* a1) {
	nox_object_t* item;

	item = (nox_object_t*)a1[126];
	if (!item) {
		return 0;
	}
	while (nox_xxx_weaponInventoryEquipFlags_415820(item) != 128) {
		item = item->inv_next_item;
		if (!item) {
			return 0;
		}
	}
	return nox_xxx_playerEquipWeapon_53A420(a1, item, 1, 1);
}

//----- (00539FF0) --------------------------------------------------------
int nox_xxx_playerTryReloadQuiver_539FF0(uint32_t* a1) {
	nox_object_t* item;

	item = (nox_object_t*)a1[126];
	if (!item) {
		return 0;
	}
	while (nox_xxx_weaponInventoryEquipFlags_415820(item) != 2) {
		item = item->inv_next_item;
		if (!item) {
			return 0;
		}
	}
	return nox_xxx_playerEquipWeapon_53A420(a1, item, 1, 1);
}

//----- (0053C580) --------------------------------------------------------
float get_nox_xxx_warriorMaxMana_587000_312788();
float get_nox_xxx_wizardMaximumMana_587000_312820();
float get_nox_xxx_conjurerMaxMana_587000_312804();
signed int nox_xxx_updateObelisk_53C580(int a1) {
	int v1;             // edi
	signed int* v2;     // esi
	int v3;             // ebp
	int v4;             // eax
	int v5;             // esi
	double v6;          // st7
	double v7;          // st6
	int v8;             // ebp
	uint32_t* v9;       // eax
	int v10;            // ecx
	int v11;            // ebx
	short v12;          // bx
	char v13;           // al
	short v14;          // ax
	int v15;            // eax
	unsigned short v16; // ax
	int v17;            // edx
	signed int result;  // eax
	float v19;          // [esp+0h] [ebp-20h]
	float v20;          // [esp+4h] [ebp-1Ch]
	int v21;            // [esp+14h] [ebp-Ch]
	int v22;            // [esp+18h] [ebp-8h]
	signed int v23;     // [esp+1Ch] [ebp-4h]
	signed int* v24;    // [esp+24h] [ebp+4h]

	v1 = a1;
	v23 = 1;
	v2 = *(signed int**)(a1 + 748);
	v24 = *(signed int**)(a1 + 748);
	v3 = nox_xxx_getFirstPlayerUnit_4DA7C0();
	v21 = v3;
	if (!v3) {
		goto LABEL_50;
	}
	while (1) {
		v4 = *(uint32_t*)(v3 + 16);
		if ((v4 & 0x8000) != 0) {
			goto LABEL_47;
		}
		v5 = *(uint32_t*)(v3 + 748);
		if (nox_xxx_servObjectHasTeam_419130(v1 + 48)) {
			if (!nox_xxx_servCompareTeams_419150(v1 + 48, v3 + 48)) {
				goto LABEL_47;
			}
		}
		v6 = *(float*)(v1 + 56) - *(float*)(v3 + 56);
		v7 = *(float*)(v1 + 60) - *(float*)(v3 + 60);
		v22 = 0;
		if (v7 * v7 + v6 * v6 >= 2500.0 || !nox_xxx_mapCheck_537110(v1, v3)) {
			goto LABEL_47;
		}
		v23 = 0;
		if (*v24 >= 1) {
			v8 = nox_xxx_getRechargeRate_53C940(*(uint32_t**)(v5 + 104));
			if (!nox_common_gameFlags_check_40A5C0(0x2000) || nox_common_gameFlags_check_40A5C0(4096)) {
				if (!v8) {
					v3 = v21;
					goto LABEL_25;
				}
			} else {
				v8 = 1;
			}
			v9 = *(uint32_t**)(v5 + 104);
			if (v9) {
				v10 = v9[2];
				if (v10 & 0x1000) {
					if (v9[3] & 0x47F0000) {
						v11 = v9[184];
						if (*(int*)(v11 + 112) < 100) {
							if (nox_common_gameFlags_check_40A5C0(4096)) {
								if (!(gameFrame() % (gameFPS() >> 1))) {
									nox_xxx_aud_501960(230, v1, 0, 0);
								}
							} else {
								v22 = 1;
								--*v24;
							}
							if (nox_xxx_rechargeItem_53C520(*(uint32_t*)(v5 + 104), v8)) {
								nox_xxx_netReportCharges_4D82B0(*(unsigned char*)(*(uint32_t*)(v5 + 276) + 2064),
																*(uint32_t**)(v5 + 104), *(uint8_t*)(v11 + 108),
																*(uint8_t*)(v11 + 109));
							}
						}
					}
				}
			}
			v3 = v21;
			goto LABEL_25;
		}
	LABEL_25:
		if ((int)*v24 < 1 || *(uint16_t*)(v5 + 4) >= *(uint16_t*)(v5 + 8) || (int)*v24 <= 0) {
			if (!v22) {
				goto LABEL_47;
			}
			goto LABEL_43;
		}
		v12 = 1;
		if (nox_common_gameFlags_check_40A5C0(4096)) {
			v13 = *(uint8_t*)(*(uint32_t*)(v5 + 276) + 2251);
			if (v13) {
				if (v13 == 1) {
					v14 = nox_float2int(get_nox_xxx_wizardMaximumMana_587000_312820());
				} else {
					if (v13 != 2) {
						goto LABEL_36;
					}
					v14 = nox_float2int(get_nox_xxx_conjurerMaxMana_587000_312804());
				}
			} else {
				v14 = nox_float2int(get_nox_xxx_warriorMaxMana_587000_312788());
			}
			v12 = v14;
		}
	LABEL_36:
		v15 = *(uint32_t*)(v5 + 276);
		*(uint16_t*)(v5 + 4) += v12;
		nox_xxx_protectMana_56F9E0(*(uint32_t*)(v15 + 4596), v12);
		v16 = *(uint16_t*)(v5 + 8);
		if (*(uint16_t*)(v5 + 4) > v16) {
			v17 = *(uint32_t*)(v5 + 276);
			*(uint16_t*)(v5 + 4) = v16;
			nox_xxx_protectPlayerHPMana_56F870(*(uint32_t*)(v17 + 4596), v16);
		}
		if (nox_common_gameFlags_check_40A5C0(4096)) {
			if (!(gameFrame() % (gameFPS() >> 1))) {
				nox_xxx_aud_501960(230, v1, 0, 0);
			}
			if (!v22) {
				goto LABEL_47;
			}
			goto LABEL_43;
		}
		--*v24;
	LABEL_43:
		if (!(*v24 % 8)) {
			v19 = (double)(80 * *v24 / 50);
			nullsub_35(v1, LODWORD(v19));
			nox_xxx_unitNeedSync_4E44F0(v1);
		}
		if ((unsigned int)(gameFrame() - *(uint32_t*)(v1 + 136)) > (int)gameFPS() >> 1) {
			nox_xxx_aud_501960(230, v1, 0, 0);
			*(uint32_t*)(v1 + 136) = gameFrame();
		}
	LABEL_47:
		v21 = nox_xxx_getNextPlayerUnit_4DA7F0(v3);
		if (!v21) {
			break;
		}
		v3 = v21;
	}
	result = v23;
	if (!v23) {
		return result;
	}
		v2 = v24;
LABEL_50:
	result = gameFrame() / (gameFPS() >> 1);
	if (!(gameFrame() % (gameFPS() >> 1))) {
		result = *v2;
		if ((int)*v2 < 50) {
			if (!(result % 8)) {
				v20 = (double)(80 * result / 50);
				nullsub_35(v1, LODWORD(v20));
				nox_xxx_unitNeedSync_4E44F0(v1);
			}
			++*v2;
		}
	}
	return result;
}
// 4E4770: using guessed type void  nullsub_35(uint32_t, uint32_t);

//----- (0053DDF0) --------------------------------------------------------
int nox_xxx_updateFlag_53DDF0(int a1) {
	int v1;           // esi
	int v2;           // edi
	int result;       // eax
	int v4;           // eax
	int v5;           // edx
	char v6;          // [esp-18h] [ebp-24h]
	unsigned char v7; // [esp+8h] [ebp-4h]

	v1 = a1;
	v2 = *(uint32_t*)(a1 + 748);
	result = *(uint32_t*)(v2 + 8);
	if (result) {
		v4 = sub_4ECBD0(a1);
		v5 = *(uint32_t*)(v2 + 8);
		a1 = v4;
		v7 = *(uint8_t*)(v1 + 52);
		result = 3 * gameFPS();
		if (gameFrame() - v5 > (unsigned int)(30 * gameFPS())) {
			nox_xxx_aud_501960(305, v1, 0, 0);
			v6 = a1;
			*(uint32_t*)(v2 + 8) = 0;
			sub_4E82C0(v7, 0, v6, 0);
			nox_xxx_unitMove_4E7010(v1, (float2*)v2);
			result = nox_xxx_netInformTextMsg2_4DA180(8, &a1);
		}
	}
	return result;
}

//----- (0053DF40) --------------------------------------------------------
void nox_xxx_updateGameBall_53DF40(int a3) {
	int v1;                // edi
	unsigned long long v2; // rax
	int v3;                // eax
	int v4;                // eax
	int v5;                // eax
	double v6;             // st7
	double v7;             // st7
	double v8;             // st7
	double v9;             // st6
	int v10;               // eax
	float4 v11;            // [esp+8h] [ebp-10h]

	v1 = *(uint32_t*)(a3 + 748);
	*(uint32_t*)(a3 + 112) = 1008981770;
	if (*(uint32_t*)v1 && *(uint8_t*)(*(uint32_t*)v1 + 16) & 0x20) {
		sub_4EB9B0(a3, 0);
		nox_xxx_netChangeTeamMb_419570(a3 + 48, *(uint32_t*)(a3 + 36));
		sub_4E8290(1, 0);
	}
	v2 = nox_platform_get_ticks() - *(uint64_t*)(v1 + 8);
	v11.field_4 = *((float*)&v2 + 1);
	if (v2 <= 20000) {
		v3 = *(uint32_t*)(a3 + 508);
		if (v3) {
			if (v3 != *(uint32_t*)v1 ||
				(unsigned int)(gameFrame() - *(uint32_t*)(v1 + 16)) <= *(int*)(v1 + 20)) {
				*(uint32_t*)(a3 + 16) |= 0x40u;
				*(uint64_t*)(v1 + 8) = nox_platform_get_ticks();
				v5 = *(uint32_t*)(a3 + 508);
				if (*(uint32_t*)(v5 + 16) & 0x8020) {
					nox_xxx_unitClearOwner_4EC300(a3);
					sub_4EB9B0(a3, 0);
					nox_xxx_netChangeTeamMb_419570(a3 + 48, *(uint32_t*)(a3 + 36));
					sub_4E8290(1, 0);
				} else {
					v6 = *(float*)(v5 + 176) + *(float*)(a3 + 176);
					v11.field_0 = *(float*)(v5 + 56);
					v11.field_4 = *(float*)(v5 + 60);
					v7 = v6 + 10.0;
					v11.field_8 = v7 * *getMemFloatPtr(0x587000, 194136 + 8 * *(short*)(v5 + 124)) + *(float*)(v5 + 56);
					v11.field_C = v7 * *getMemFloatPtr(0x587000, 194140 + 8 * *(short*)(v5 + 124)) + *(float*)(v5 + 60);
					if (nox_xxx_mapTraceRay_535250(&v11, 0, 0, 5)) {
						nox_xxx_unitMove_4E7010(a3, (float2*)&v11.field_8);
					}
				}
			} else {
				*(uint32_t*)(a3 + 16) &= 0xFFFFFFBF;
				*(uint32_t*)(a3 + 520) = 0;
				v4 = *(short*)(*(uint32_t*)(a3 + 508) + 124) + nox_common_randomInt_415FA0(-32, 32);
				if (v4 < 0) {
					v4 += (unsigned int)(255 - v4) >> 8 << 8;
				}
				if (v4 >= 256) {
					v4 += -256 * ((unsigned int)v4 >> 8);
				}
				v11.field_0 = *(float*)(a3 + 56) - *getMemFloatPtr(0x587000, 194136 + 8 * v4) * 20.0;
				v11.field_4 = *(float*)(a3 + 60) - *getMemFloatPtr(0x587000, 194140 + 8 * v4) * 20.0;
				nox_xxx_objectApplyForce_52DF80((int)&v11, a3, 30.0);
				nox_xxx_unitClearOwner_4EC300(a3);
				sub_4E8290(1, 0);
				nox_xxx_aud_501960(926, a3, 0, 0);
			}
		} else {
			v8 = *(float*)(a3 + 84);
			v9 = *(float*)(a3 + 80);
			v10 = *(uint32_t*)(a3 + 16);
			LOBYTE(v10) = v10 & 0xBF;
			*(uint32_t*)(a3 + 16) = v10;
			if (*(float*)(v1 + 24) > sqrt(v9 * v9 + v8 * v8)) {
				*(uint32_t*)v1 = 0;
			}
		}
	} else {
		sub_417F50(a3);
	}
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

//----- (0053E1D0) --------------------------------------------------------
void nox_xxx_updateCrown_53E1D0(int a1) {
	uint32_t* v1; // ecx
	int v2;       // eax
	int v3;       // eax
	double v4;    // st7
	double v5;    // st7
	float4 v6;    // [esp+4h] [ebp-10h]

	v1 = *(uint32_t**)(a1 + 748);
	v2 = v1[1];
	if (!v2 || *(uint32_t*)(v2 + 16) & 0x8020) {
		if (*v1 && *(uint8_t*)(*v1 + 16) & 0x20) {
			*v1 = 0;
		}
		v3 = *(uint32_t*)(a1 + 508);
		if (v3) {
			if (*(uint32_t*)(v3 + 16) & 0x8020) {
				nox_xxx_unitClearOwner_4EC300(a1);
			} else {
				v4 = *(float*)(v3 + 176) + *(float*)(a1 + 176);
				v6.field_0 = *(float*)(v3 + 56);
				v6.field_4 = *(float*)(v3 + 60);
				v5 = v4 + 10.0;
				v6.field_8 = v5 * *getMemFloatPtr(0x587000, 194136 + 8 * *(short*)(v3 + 124)) + *(float*)(v3 + 56);
				v6.field_C = v5 * *getMemFloatPtr(0x587000, 194140 + 8 * *(short*)(v3 + 124)) + *(float*)(v3 + 60);
				if (nox_xxx_mapTraceRay_535250(&v6, 0, 0, 5)) {
					nox_xxx_unitMove_4E7010(a1, (float2*)&v6.field_8);
				}
			}
		}
	} else {
		sub_4F3400(v2, a1, 1);
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

//----- (00543680) --------------------------------------------------------
int sub_543680(float* a1) {
	int result; // eax
	float2 a2;  // [esp+0h] [ebp-30h]
	int2 v3;    // [esp+8h] [ebp-28h]
	int v4[8];  // [esp+10h] [ebp-20h]

	if (dword_5d4594_3835356 == 255) {
		return 1;
	}
	result = nox_xxx_mapGenFixCoords_4D3D90((float2*)a1, &a2);
	if (result) {
		v4[0] = *getMemU32Ptr(0x973F18, 35912);
		dword_5d4594_3835352 = 1;
		v4[1] = dword_5d4594_3835348;
		v4[6] = dword_5d4594_3835356;
		v4[7] = dword_5d4594_3835360;
		v3.field_0 = (long long)a2.field_0;
		v3.field_4 = (long long)a2.field_4;
		result = sub_5437E0(&v3.field_0, (int)v4, 46);
		dword_5d4594_3835352 = 0;
	}
	return result;
}

//----- (005437E0) --------------------------------------------------------
int sub_5437E0(int* a1, int a2, int a3) {
	double v3;          // st7
	double v4;          // st6
	int v5;             // ebx
	double v6;          // st6
	int v7;             // edi
	int v8;             // esi
	int v9;             // edx
	int v10;            // ecx
	int v11;            // esi
	int v12;            // edx
	int v13;            // ebx
	int v14;            // ecx
	int v15;            // edx
	int v16;            // ebx
	unsigned char* v17; // edi
	int v18;            // edx
	int v19;            // eax
	int* v20;           // ecx
	int v21;            // edi
	int v23;            // [esp+Ch] [ebp-8h]
	float v24;          // [esp+10h] [ebp-4h]

	v3 = (double)*a1 + 11.5;
	v4 = (double)a1[1];
	v5 = (long long)(v3 * 0.021739131);
	v23 = (long long)(v3 * 0.021739131);
	v6 = v4 + 11.5;
	v24 = v6;
	v7 = (long long)(v6 * 0.021739131);
	a1 = (int*)(long long)(v6 * 0.021739131);
	v8 = (int)(long long)v3 % 46;
	*getMemU32Ptr(0x973F18, 22200) = 0;
	v9 = (long long)v24 % 46;
	dword_5d4594_2487248 = 0;
	if (v5 <= 0 || v5 >= 127 || v7 <= 0 || v7 >= 127) {
		v11 = a3;
	} else if (v8 <= v9) {
		if (a3 - v8 <= v9) {
			v15 = ptr_5D4594_2650668[v5];
			v11 = *(uint32_t*)(v15 + 44 * v7 + 24);
			sub_51DD50(v5, v7, 2, *(uint32_t*)(v15 + 44 * v7 + 24));
		} else {
			v13 = v5 - 1;
			v14 = ptr_5D4594_2650668[v13];
			v11 = *(uint32_t*)(v14 + 44 * v7 + 4);
			sub_51DD50(v13, v7, 1, *(uint32_t*)(v14 + 44 * v7 + 4));
		}
	} else if (a3 - v8 <= v9) {
		v12 = ptr_5D4594_2650668[v5];
		v11 = *(uint32_t*)(v12 + 44 * v7 + 4);
		sub_51DD50(v5, v7, 1, *(uint32_t*)(v12 + 44 * v7 + 4));
	} else {
		v10 = ptr_5D4594_2650668[v5];
		v11 = *(uint32_t*)(v10 + 44 * v7 - 20);
		sub_51DD50(v5, v7 - 1, 2, *(uint32_t*)(v10 + 44 * v7 - 20));
	}
	v16 = 0;
	if (dword_5d4594_2487248 > 0) {
		v17 = getMemAt(0x973F18, 16204);
		do {
			v18 = *((uint32_t*)v17 + 1);
			v19 = *((uint32_t*)v17 - 1);
			v20 = *(int**)v17;
			v23 = *((uint32_t*)v17 - 1);
			a1 = v20;
			a3 = v18;
			if (v18 & 2) {
				sub_51DD50(v19, (int)v20, 1, v11);
				sub_51DD50(v23, (int)a1 + 1, 1, v11);
				sub_51DD50(v23 - 1, (int)a1, 1, v11);
				sub_51DD50(v23 - 1, (int)a1 + 1, 1, v11);
			} else if (v18 & 1) {
				sub_51DD50(v19 + 1, (int)v20, 2, v11);
				sub_51DD50(v23 + 1, (int)a1 - 1, 2, v11);
				sub_51DD50(v23, (int)a1, 2, v11);
				sub_51DD50(v23, (int)a1 - 1, 2, v11);
			}
			++v16;
			v17 += 12;
		} while (v16 < *(int*)&dword_5d4594_2487248);
	}
	if (sub_51DE30(&v23, &a1, &a3)) {
		v21 = a2;
		do {
			if (a3 & 2) {
				sub_543BC0(v23, (int)a1, 1, v11, v21, 1);
				sub_543BC0(v23, (int)a1 + 1, 1, v11, v21, 4);
				sub_543BC0(v23 - 1, (int)a1, 1, v11, v21, 3);
				sub_543BC0(v23 - 1, (int)a1 + 1, 1, v11, v21, 6);
				sub_543BC0(v23, (int)a1 - 1, 2, v11, v21, 0);
				sub_543BC0(v23 + 1, (int)a1, 2, v11, v21, 2);
				sub_543BC0(v23 - 1, (int)a1, 2, v11, v21, 5);
				sub_543BC0(v23, (int)a1 + 1, 2, v11, v21, 7);
			} else if (a3 & 1) {
				sub_543BC0(v23 + 1, (int)a1, 2, v11, v21, 4);
				sub_543BC0(v23 + 1, (int)a1 - 1, 2, v11, v21, 1);
				sub_543BC0(v23, (int)a1, 2, v11, v21, 6);
				sub_543BC0(v23, (int)a1 - 1, 2, v11, v21, 3);
				sub_543BC0(v23, (int)a1 - 1, 1, v11, v21, 0);
				sub_543BC0(v23 + 1, (int)a1, 1, v11, v21, 2);
				sub_543BC0(v23 - 1, (int)a1, 1, v11, v21, 5);
				sub_543BC0(v23, (int)a1 + 1, 1, v11, v21, 7);
			}
		} while (sub_51DE30(&v23, &a1, &a3));
	}
	return *getMemU32Ptr(0x973F18, 22200) == 0;
}

//----- (00543BC0) --------------------------------------------------------
void sub_543BC0(int a1, int a2, int a3, int a4, int a5, int a6) {
	int v7; // ebx

	if (a1 > 0 && a1 < 127) {
		if (a2 > 0 && a2 < 127) {
			v7 = a3;
			if ((!(a3 & 1) || a2 != 1) && (!(a3 & 2) || a1 != 1)) {
				if (a3 & 2) {
					if (a4 == *(uint32_t*)((uint32_t)(ptr_5D4594_2650668[a1]) + 44 * a2 + 24)) {
						return;
					}
					v7 = a3;
				}
				if (!(a3 & 1) || a4 != *(uint32_t*)((uint32_t)(ptr_5D4594_2650668[a1]) + 44 * a2 + 4)) {
					if (a5) {
						*(uint32_t*)(a5 + 28) = a6;
						sub_51DA70(a1, a2, a5, v7, 1);
					}
				}
			}
		}
	}
}

//----- (00543C50) --------------------------------------------------------
int nox_xxx_tile_543C50(uint32_t* a1, int a2, int a3, int a4, int a5, int a6) {
	uint32_t* v6;   // edi
	uint32_t* v7;   // esi
	uint32_t** v8;  // eax
	int result;     // eax
	uint32_t** v10; // eax
	int v11;        // ebx
	int v12;        // eax
	uint32_t* v13;  // ebx
	uint32_t** v14; // edi
	uint32_t* j;    // ebp
	uint32_t* v16;  // esi
	int v17;        // eax
	int v18;        // eax
	uint32_t* v19;  // esi
	int v20;        // eax
	int i;          // ecx
	int v22;        // [esp+24h] [ebp+14h]

	if (a2 == 255 || a4 == 255) {
		v19 = a1;
		v20 = a1[4];
		if (v20) {
			for (i = *(uint32_t*)(v20 + 16); i; i = *(uint32_t*)(i + 16)) {
				v19 = (uint32_t*)v20;
				v20 = i;
			}
			nox_xxx_tileFreeTileOne_4221E0(v20);
			v19[4] = 0;
		}
		return 1;
	}
	v6 = a1;
	if (*a1 == 255) {
		return 1;
	}
	v7 = a1;
	if (a6) {
		v10 = (uint32_t**)(a1 + 4);
		v11 = 0;
		if (!a1[4]) {
			v12 = nox_xxx_mapGenEdge_543EB0(a4, a5);
			v7[4] = nox_xxx_tileListAddNewSubtile_422160(a2, a3, a4, v12);
			return 1;
		}
		do {
			v7 = *v10;
			if (**v10 == a2 && v7[1] == a3 && v7[2] == a4 && sub_543E60((int)v7, a5)) {
				v11 = 1;
			}
			v10 = (uint32_t**)(v7 + 4);
		} while (v7[4]);
		if (!v11) {
			v12 = nox_xxx_mapGenEdge_543EB0(a4, a5);
			v7[4] = nox_xxx_tileListAddNewSubtile_422160(a2, a3, a4, v12);
			return 1;
		}
		while (1) {
			v13 = v6;
			v22 = 0;
			if (!v6[4]) {
				break;
			}
			do {
				v13 = (uint32_t*)v13[4];
				if (*v13 == a2 && v13[1] == a3 && v13[2] == a4) {
					v14 = (uint32_t**)(v13 + 4);
					for (j = v13; j[4]; v14 = (uint32_t**)(j + 4)) {
						v16 = *v14;
						if (**v14 == a2 && v16[1] == a3 && (v17 = v16[2], v17 == a4) &&
							(v18 = sub_411490(v17, v16[3]), sub_543E60((int)v13, v18))) {
							v22 = 1;
							*v14 = (uint32_t*)v16[4];
							nox_xxx_tileFreeTileOne_4221E0(v16);
						} else {
							j = v16;
						}
					}
					v6 = a1;
				}
			} while (v13[4]);
			if (!v22) {
				return 1;
			}
		}
		return 1;
	}
	v8 = (uint32_t**)(a1 + 4);
	if (a1[4]) {
		while (1) {
			v7 = *v8;
			if (**v8 == a2 && v7[1] == a3 && v7[2] == a4 && v7[3] == a5) {
				break;
			}
			v8 = (uint32_t**)(v7 + 4);
			if (!v7[4]) {
				v7[4] = nox_xxx_tileListAddNewSubtile_422160(a2, a3, a4, a5);
				return 1;
			}
		}
		return 0;
	} else {
		v7[4] = nox_xxx_tileListAddNewSubtile_422160(a2, a3, a4, a5);
		return 1;
	}
}

//----- (00544310) --------------------------------------------------------
int nox_xxx_tileSubtile_544310(float2* a1) {
	double v1;  // st7
	double v2;  // st6
	int v3;     // esi
	int v4;     // edi
	int v5;     // ebp
	int v6;     // edx
	int result; // eax
	int v8[8];  // [esp+10h] [ebp-20h]
	float v9;   // [esp+34h] [ebp+4h]

	dword_5d4594_3835352 = 1;
	v1 = a1->field_0 + 11.5;
	v8[1] = dword_5d4594_3835348;
	v8[0] = *getMemU32Ptr(0x973F18, 35912);
	v8[2] = 0;
	LOBYTE(v8[5]) = 0;
	v8[3] = -1;
	v8[4] = -1;
	v8[6] = dword_5d4594_3835356;
	v8[7] = dword_5d4594_3835360;
	v2 = a1->field_4 + 11.5;
	v3 = (long long)(v1 * 0.021739131);
	v9 = v2;
	v4 = (long long)(v2 * 0.021739131);
	v5 = (int)(long long)v1 % 46;
	v6 = (long long)v9 % 46;
	if (*getMemU32Ptr(0x973F18, 35912) == 255) {
		result = sub_51D9C0(v3, v4, v5, v6, 0);
	} else {
		result = sub_51D9C0(v3, v4, v5, v6, (int)v8);
	}
	dword_5d4594_3835352 = 0;
	return result;
}
