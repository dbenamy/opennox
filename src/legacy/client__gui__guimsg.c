#include "client__gui__guimsg.h"
#include "client__gui__guicon.h"
#include "common__strman.h"
extern uint32_t dword_5d4594_825736;

//----- (00445490) --------------------------------------------------------
void nox_xxx_printCentered_445490(wchar2_t* a1) {
	int v1;             // eax
	int v2;             // eax
	unsigned short* v3; // eax

	if (a1) {
		v1 = ++*(unsigned int*)&dword_5d4594_825736;
		if (*(unsigned int*)&dword_5d4594_825736 == 3) {
			v1 = 0;
			*(unsigned int*)&dword_5d4594_825736 = 0;
		}
		// Each row has 318 text units before its expiry word, including the terminator.
		wchar2_t* text = getMemU16Ptr(0x5D4594, 823804 + 644 * v1);
		size_t length = 0;
		while (length < 317 && a1[length]) {
			text[length] = a1[length];
			++length;
		}
		text[length] = 0;
		v2 = 644 * *(unsigned int*)&dword_5d4594_825736;
		*getMemUintPtr(0x5D4594, 824440 + v2) = gameFrame() + 4 * gameFPS() + gameFPS();
		*getMemU8Ptr(0x5D4594, 824444 + v2) = 0;
		v3 = nox_strman_loadString_40F1D0("systemmsg", 0, "C:\\NoxPost\\src\\Client\\Gui\\guimsg.c", 69);
		nox_gui_console_Printf_450C00(NOX_CONSOLE_RED, v3, a1);
	}
}
