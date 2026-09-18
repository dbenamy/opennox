#include <math.h>
#include <time.h>

#include "common__system__team.h"
#include "server__script__script.h"
#include "server__system__server.h"
#include "server__script__activator.h"

#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME1_2.h"
#include "GAME1_3.h"
#include "GAME2.h"
#include "GAME2_1.h"
#include "GAME2_2.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_1.h"
#include "GAME4_2.h"
#include "GAME4_3.h"
#include "GAME5.h"
#include "GAME5_2.h"
#include "client__gui__guiquit.h"
#include "client__video__draw_common.h"
#include "common/fs/nox_fs.h"
#include "common__log.h"
#include "common__magic__speltree.h"
#include "common__net_list.h"
#include "common__random.h"
#include "common__strman.h"
#include "common__crypt.h"
#include "common__system__settings.h"
#include "input_common.h"
#include "operators.h"

extern uint32_t dword_5d4594_600116;
extern uint32_t dword_5d4594_1548480;
extern uint32_t dword_5d4594_608316;
extern uint32_t dword_5d4594_1569756;
extern uint32_t dword_5d4594_1548476;

//----- (00426060) --------------------------------------------------------
void sub_426060() {
	char* v0;     // eax
	char* i;      // eax
	char* v2;     // eax
	char* j;      // esi
	void* result; // eax
	char* v5;     // edi
	char* v6;     // esi
	int v7;       // [esp-4h] [ebp-4h]

	dword_5d4594_608316 = 0;
	dword_5d4594_600116 = time(0);
	v7 = sub_5545A0();
	v0 = sub_554230();
	sub_4282D0(v0, v7);
	if (!nox_common_gameFlags_check_40A5C0(0x2000) || nox_common_gameFlags_check_40A5C0(4096)) {
		result = (void*)nox_common_gameFlags_check_40A5C0(4096);
		if (result) {
			result = (void*)nox_common_gameFlags_check_40A5C0(4096);
			if (result) {
				v5 = sub_416640();
				v6 = nox_xxx_cliGamedataGet_416590(0);
				*getMemU16Ptr(0x5D4594, 739396) = sub_40A770();
				*getMemU32Ptr(0x5D4594, 739400) = *((uint32_t*)v5 + 10);
				*getMemU32Ptr(0x5D4594, 739404) = sub_4200E0();
				*getMemU8Ptr(0x5D4594, 739412) = (v6[53] & 0xC0) != 0;
				*getMemU32Ptr(0x5D4594, 739408) = 5;
				strncpy((char*)getMemAt(0x5D4594, 739676), v6 + 9, 0xFu);
				*getMemU8Ptr(0x5D4594, 739691) = 0;
				strncpy((char*)getMemAt(0x5D4594, 739420), v6, 8u);
				*getMemU8Ptr(0x5D4594, 739428) = 0;
				sub_4289D0((void**)getMemAt(0x5D4594, 739396));
			}
		}
	} else {
		for (i = nox_common_playerInfoGetFirst_416EA0(); i; i = nox_common_playerInfoGetNext_416EE0((int)i)) {
			*((uint32_t*)i + 1162) = -1;
		}
		v2 = nox_common_playerInfoFromNum_417090(31);
		if (v2) {
			sub_425F10((int)v2);
		}
		for (j = nox_common_playerInfoGetFirst_416EA0(); j; j = nox_common_playerInfoGetNext_416EE0((int)j)) {
			if (j[2064] != 31) {
				sub_425F10((int)j);
			}
		}
		sub_426150();
		sub_428810((int)getMemAt(0x5D4594, 599476), 0);
		*getMemU16Ptr(0x5D4594, 599482) = 0;
	}
}
