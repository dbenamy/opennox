#include <float.h>
#include <math.h>

#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME1_2.h"
#include "GAME1_3.h"
#include "GAME2.h"
#include "GAME2_2.h"
#include "GAME2_3.h"
#include "GAME3_1.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_1.h"
#include "GAME4_2.h"
#include "GAME4_3.h"
#include "GAME5.h"
#include "GAME5_2.h"
#include "client__system__parsecmd.h"
#include "common__binfile.h"
#include "common__random.h"
#include "common__strman.h"

#include "client__gui__guiinv.h"
#include "client__shell__mainmenu.h"
#include "common__system__team.h"

#include "MixPatch.h"
#include "client__video__draw_common.h"
#include "common/fs/nox_fs.h"
#include "common__magic__speltree.h"
#include "operators.h"

extern uint32_t dword_5d4594_2523804;
extern uint32_t dword_5d4594_2523776;
extern uint64_t qword_581450_9544;
extern uint32_t dword_5d4594_2523780;
extern uint32_t dword_5d4594_2650652;
extern uint32_t dword_8531A0_2576;

//----- (00554040) --------------------------------------------------------
unsigned int nox_server_makeServerInfoPacket_554040(const char* inBuf, int inSz, char* out) {
	char buf[72];

	char* v3 = sub_416640();
	char* game = nox_xxx_cliGamedataGet_416590(0);
	if (!sub_43AF30() || sub_4D6F30()) {
		return 0;
	}
	char playerLimit = nox_xxx_servGetPlrLimit_409FA0();
	char playerCount = nox_common_playerInfoCount_416F40();
	if (nox_common_getEngineFlag(NOX_ENGINE_FLAG_DISABLE_GRAPHICS_RENDERING)) {
		--playerCount;
		--playerLimit;
	}
	char* srvName = nox_xxx_serverOptionsGetServername_40A4C0();
	buf[0] = 0;
	buf[1] = 0;
	buf[2] = 13;
	buf[3] = playerCount;
	buf[4] = playerLimit;
	buf[5] = v3[101] & 0xF;
	buf[6] = ((unsigned char)v3[101]) >> 4;
	*(uint32_t*)&buf[7] = *((uint32_t*)game + 11);
	strcpy(&buf[10], nox_xxx_mapGetMapName_409B40());
	buf[19] = v3[102] | sub_43BE50_get_video_mode_id();
	buf[20] = v3[100];
	buf[21] = v3[100] & 0x10;
	*(uint32_t*)&buf[24] = *((uint32_t*)v3 + 12);
	unsigned int gameFlags = nox_common_gameFlags_getVal_40A5B0();
	if (nox_xxx_isQuest_4D6F50()) {
		gameFlags = (gameFlags & 0xFFFFFF7Fu) | 0x1000u;
		*(uint16_t*)&buf[68] = nox_game_getQuestStage_4E3CC0();
	}
	*(uint32_t*)&buf[28] = gameFlags;
	*(uint32_t*)&buf[32] = *((uint32_t*)game + 12);
	*(uint16_t*)&buf[36] = *(uint16_t*)(v3 + 105);
	*(uint16_t*)&buf[38] = *(uint16_t*)(v3 + 107);
	*(uint32_t*)&buf[40] = *((uint32_t*)v3 + 11);
	*(uint32_t*)&buf[44] = *(uint32_t*)(&inBuf[8]); // timestamp of the packet
	memcpy(&buf[48], game + 24, 20);
	memcpy(&out[0], buf, 72);
	strcpy(&out[72], srvName);
	return 72 + strlen(srvName) + 1;
}

//----- (0057ADF0) --------------------------------------------------------


//----- (0057AF20) --------------------------------------------------------


//----- (0057B0A0) --------------------------------------------------------


//----- (0057B180) --------------------------------------------------------


//----- (0057B3D0) --------------------------------------------------------
int nox_cheat_allowall = 0;

//----- (0057C090) --------------------------------------------------------


//----- (0057CDB0) --------------------------------------------------------

//----- (0057F1D0) --------------------------------------------------------

//----- (0057F2A0) --------------------------------------------------------

void nullsub_10(uint32_t a1) {}





void nullsub_30(uint32_t a1) {}
void nullsub_29(void) {}
void nullsub_35(uint32_t a1, uint32_t a2) {}
void nullsub_24(uint32_t a1) {}
void nullsub_31(uint32_t a1) {}
void nullsub_9(uint32_t a1) {}
