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
#include "GAME5_2.h"
#include "client__draw__animdraw.h"
#include "client__drawable__drawable.h"
#include "client__gui__guisumn.h"
#include "client__gui__servopts__advserv.h"
#include "common__net_list.h"
#include "common__system__team.h"

#include "client__gui__guiinput.h"
#include "client__gui__servopts__objlst.h"
#include "client__gui__servopts__spelllst.h"
#include "client__gui__tooltip.h"
#include "client__gui__window.h"
#include "client__shell__inputcfg__inputcfg.h"
#include "client__shell__mainmenu.h"

#include "client__draw__plasma.h"
#include "client__drawable__update__charmup.h"
#include "client__drawable__update__fireball.h"
#include "client__system__ctrlevnt.h"

#include "client__shell__optsback.h"
#include "client__video__draw_common.h"
#include "common/fs/nox_fs.h"
#include "common__strman.h"
#include "input_common.h"
#include "operators.h"

extern uint32_t dword_8531A0_2572;
extern uint32_t dword_8531A0_2576;
extern uint32_t dword_587000_180476;
extern uint32_t dword_5d4594_1319260;
extern uint32_t dword_5d4594_1321196;
extern uint32_t dword_5d4594_1313880;
extern uint32_t dword_5d4594_1320944;
extern uint32_t dword_5d4594_1320948;
extern uint32_t dword_587000_183460;
extern uint32_t dword_5d4594_1316412;
extern uint32_t dword_5d4594_1522968;
extern uint32_t dword_5d4594_1321520;
extern uint32_t dword_5d4594_1319268;
extern uint32_t dword_5d4594_1321024;
extern uint32_t dword_5d4594_1320936;
extern uint32_t dword_587000_183456;
extern uint32_t dword_5d4594_1320988;
extern uint32_t dword_5d4594_1319248;
extern uint32_t dword_587000_180480;
extern uint32_t dword_5d4594_1319236;
extern uint32_t dword_5d4594_1320972;
extern uint32_t dword_5d4594_1321208;
extern uint32_t dword_5d4594_1321800;
extern uint32_t dword_5d4594_1321224;
extern uint32_t dword_5d4594_1320932;
extern uint32_t dword_5d4594_1321032;
extern uint32_t dword_5d4594_1319264;
extern uint32_t dword_5d4594_1321044;
extern uint32_t dword_5d4594_1319232;
extern nox_window* dword_5d4594_1321236;
extern nox_window* dword_5d4594_1321240;
extern nox_window* dword_5d4594_1321248;
extern nox_window* dword_5d4594_1321244;
nox_window* dword_5d4594_1522616 = 0;
nox_window* dword_5d4594_1522620 = 0;
nox_window* dword_5d4594_1522624 = 0;
nox_window* dword_5d4594_1522628 = 0;
extern uint32_t dword_5d4594_1320992;
extern uint32_t dword_5d4594_1321204;
extern uint32_t dword_5d4594_1316408;
extern uint64_t qword_581450_9512;
extern uint64_t qword_581450_9544;
extern uint32_t dword_5d4594_1320968;
extern uint32_t dword_5d4594_1522632;
extern uint32_t dword_5d4594_1321252;
extern uint32_t dword_5d4594_1522612;
extern uint32_t dword_5d4594_1522604;
extern uint32_t dword_5d4594_1321232;
extern uint32_t dword_5d4594_1320964;
extern uint32_t dword_5d4594_1321228;
extern uint32_t dword_5d4594_1321040;
extern uint32_t dword_5d4594_1320940;
extern uint32_t nox_player_netCode_85319C;
extern int nox_win_width;
extern int nox_win_height;

extern uint32_t nox_color_white_2523948;
extern uint32_t nox_color_blue_2650684;
extern uint32_t nox_color_yellow_2589772;
extern uint32_t nox_color_violet_2598268;
extern uint32_t nox_color_black_2650656;

nox_gui_animation* nox_wnd_xxx_1522608 = 0;

void* nox_gui_itemAmount_item_1319256 = 0;
void* nox_gui_itemAmount_dialog_1319228 = 0;

// 487CF0: using guessed type void  nullsub_10(uint32_t);

// 487CA0: using guessed type void  nullsub_9(uint32_t);

//----- (004BF9F0) --------------------------------------------------------

// 4BFA4C: variable 'v10' is possibly undefined
// 4BFA4C: variable 'v11' is possibly undefined

//----- (004CA650) --------------------------------------------------------

// 4CA67E: variable 'v6' is possibly undefined

nox_window* dword_5d4594_1321236 = 0;
nox_window* dword_5d4594_1321240 = 0;
nox_window* dword_5d4594_1321244 = 0;
nox_window* dword_5d4594_1321248 = 0;
