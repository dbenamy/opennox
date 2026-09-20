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
#include "client__drawable__drawable.h"
#include "common__system__team.h"

#include "client__gui__chathelp.h"
#include "client__gui__gadgets__listbox.h"
#include "client__gui__guicon.h"
#include "client__gui__guiinv.h"
#include "client__gui__guiquit.h"
#include "client__gui__guivote.h"
#include "client__gui__window.h"
#include "client__network__cdecode.h"
#include "client__shell__noxworld.h"

#include "client__draw__fx.h"
#include "client__gui__guibook.h"
#include "client__video__draw_common.h"

#include "common/fs/nox_fs.h"
#include "common__magic__speltree.h"
#include "common__net_list.h"
#include "defs.h"
#include "input_common.h"
#include "operators.h"

extern uint32_t dword_8531A0_2572;
extern uint32_t nox_server_sanctuaryHelp_54276;
extern uint32_t dword_5d4594_1305748;
extern uint32_t nox_server_connectionType_3596;
extern void* nox_alloc_pixelSpan_1301844;
extern uint32_t dword_5d4594_1301848;
extern uint32_t nox_player_netCode_85319C;
extern int nox_win_width;
extern int nox_win_height;

extern uint32_t nox_color_white_2523948;
extern uint32_t nox_color_yellow_2589772;
extern uint32_t nox_color_black_2650656;

nox_render_data_t* nox_draw_curDrawData_3799572 = 0;

void nox_gui_winParentsReset_4A0CF0();
nox_window* nox_gui_winParentsTop_4A14F0();
nox_window* nox_gui_winParentsPop_4A18A0();
void nox_gui_winParentsPush_4A18C0(nox_window* win);

//----- (0048C580) --------------------------------------------------------
void sub_48C580(pixel8888* a1, int num) {
	unsigned int* pix = (unsigned int*)a1;
	for (int i = num - 1; i >= 0; i--) {
		unsigned int result = *pix;
		for (unsigned int* it = &pix[i]; it > pix; --it) {
			if (result > *it) {
				result = __sync_lock_test_and_set((volatile signed int*)it, result);
			}
		}
		*pix = result;
		++pix;
	}
}

//----- (004947E0) --------------------------------------------------------


//----- (004948B0) --------------------------------------------------------


//----- (00494F00) --------------------------------------------------------

// 49A025: variable 'v11' is possibly undefined
// 49A025: variable 'v12' is possibly undefined

//----- (0049AEA0) --------------------------------------------------------
int sub_49AEA0() {
	if (nox_alloc_pixelSpan_1301844) {
		nox_free_alloc_class(*(void**)&nox_alloc_pixelSpan_1301844);
		nox_alloc_pixelSpan_1301844 = 0;
	}
	if (dword_5d4594_1301848) {
		free(*(void**)&dword_5d4594_1301848);
		dword_5d4594_1301848 = 0;
	}
	return 1;
}

//----- (0049B7A0) --------------------------------------------------------

// 49B874: variable 'v1' is possibly undefined
// 49B881: variable 'v2' is possibly undefined

// 49C248: variable 'v2' is possibly undefined
// 49C256: variable 'v4' is possibly undefined
