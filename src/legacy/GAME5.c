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
extern uint32_t dword_5d4594_2491580;
extern uint32_t dword_5d4594_2491676;
extern uint32_t dword_5d4594_2491588;
extern uint32_t dword_5d4594_2491704;
extern uint32_t dword_5d4594_2489460;
extern uint32_t dword_5d4594_2650652;
extern uint32_t nox_player_netCode_85319C;

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

//----- (0054AF40) --------------------------------------------------------

//----- (0054AFB0) --------------------------------------------------------


//----- (0054E6F0) --------------------------------------------------------

//----- (0054E730) --------------------------------------------------------

//----- (0054E810) --------------------------------------------------------

//----- (0054E850) --------------------------------------------------------
