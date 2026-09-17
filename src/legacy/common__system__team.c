#include "common__system__team.h"
// FIXME
#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME2.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "client__gui__guicon.h"
#include "client__gui__guimsg.h"
#include "client__gui__servopts__guiserv.h"
#include "common__strman.h"
#include "operators.h"

extern uint32_t nox_player_netCode_85319C;

typedef struct {
	const char* name;
	wchar2_t* title;
	int code;
	uint32_t* color;
} nox_team_info_t;

uint32_t nox_color_white_2523948 = 0;
uint32_t nox_color_red_2589776 = 0;
uint32_t nox_color_blue_2650684 = 0;
uint32_t nox_color_green_2614268 = 0;
uint32_t nox_color_cyan_2649820 = 0;
uint32_t nox_color_yellow_2589772 = 0;
uint32_t nox_color_violet_2598268 = 0;
uint32_t nox_color_black_2650656 = 0;
uint32_t nox_color_orange_2614256 = 0;
