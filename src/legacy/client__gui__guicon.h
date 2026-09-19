#ifndef NOX_PORT_CLIENT_GUI_GUICON
#define NOX_PORT_CLIENT_GUI_GUICON

#include "client__gui__window.h"
#include "nox_wchar.h"
#include <stdlib.h>

enum {
	NOX_CONSOLE_BLACK = 1,
	NOX_CONSOLE_DARK_GREY = 2,
	NOX_CONSOLE_LIGHT_GREY = 3,
	NOX_CONSOLE_WHITE = 4,
	NOX_CONSOLE_DARK_RED = 5,
	NOX_CONSOLE_RED = 6,
	NOX_CONSOLE_LIGHT_RED = 7,
	NOX_CONSOLE_DARK_GREEN = 8,
	NOX_CONSOLE_GREEN = 9,
	NOX_CONSOLE_LIGHT_GREEN = 10,
	NOX_CONSOLE_DARK_BLUE = 11,
	NOX_CONSOLE_BLUE = 12,
	NOX_CONSOLE_LIGHT_BLUE = 13,
	NOX_CONSOLE_DARK_YELLOW = 14,
	NOX_CONSOLE_YELLOW = 15,
	NOX_CONSOLE_LIGHT_YELLOW = 16
};

int nox_gui_console_Hide_4512B0();
int nox_gui_console_flagXxx_451410();

int nox_gui_console_Print_450B90(unsigned char cl, wchar2_t* str);
void nox_gui_console_PrintOrError_450C30(unsigned char cl, wchar2_t* str);

static int nox_gui_console_Printf_450C00(unsigned char cl, wchar2_t* fmt, ...) {
	wchar2_t local[512];
	wchar2_t* text = local;
	va_list va, measure;
	va_start(va, fmt);
	va_copy(measure, va);
	int length = nox_vsnwprintf(local, 512, fmt, measure);
	va_end(measure);
	if (length < 0 || (size_t)length > SIZE_MAX / sizeof(*text) - 1) {
		va_end(va);
		return 0;
	}
	if (length >= 512) {
		text = malloc(((size_t)length + 1) * sizeof(*text));
		if (!text) {
			va_end(va);
			return 0;
		}
		nox_vsnwprintf(text, (size_t)length + 1, fmt, va);
	}
	va_end(va);
	int result = nox_gui_console_Print_450B90(cl, text);
	if (text != local) free(text);
	return result;
}

#endif // NOX_PORT_CLIENT_GUI_GUICON
