#ifndef NOX_CLIENT_GUI_WINDOW_H
#define NOX_CLIENT_GUI_WINDOW_H

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#include "static_assert.h"
#include "nox_wchar.h"

typedef enum { NOX_WIN_HIDDEN = 0x10, NOX_WIN_LAYER_FRONT = 0x20, NOX_WIN_LAYER_BACK = 0x40 } nox_window_flags;

typedef struct nox_window nox_window;
typedef struct nox_window_data {
	unsigned int field_0; // 0, 0 (36)
	int group;            // 1, 4 (40)
	int style;            // 2, 8 (44)
	int status;           // 3, 12 (48)
	nox_window* win;      // 4, 16 (52)
	uint32_t bg_color;    // 5, 20 (56)
	void* bg_image;       // 6, 24 (60), nox_video_bag_image_t
	uint32_t en_color;    // 7, 28 (64)
	void* en_image;       // 8, 32 (68), nox_video_bag_image_t
	uint32_t hl_color;    // 9, 36 (72)
	void* hl_image;       // 10, 40 (76), nox_video_bag_image_t
	uint32_t dis_color;   // 11, 44 (80)
	void* dis_image;      // 12, 48 (84), nox_video_bag_image_t
	uint32_t sel_color;   // 13, 52 (88)
	void* sel_image;      // 14, 56 (92), nox_video_bag_image_t
	int img_px;           // 15, 60 (96)
	int img_py;           // 16, 64 (100)
	uint32_t text_color;  // 17, 68 (104)
	wchar2_t text[64];     // 18, 72 (108)
	void* font;           // 50, 200 (236)
	wchar2_t tooltip[64];  // 51, 204 (240)
} nox_window_data;
_Static_assert(sizeof(nox_window_data) == 332, "wrong size of nox_window_data structure!");

typedef struct nox_window {
	int id;                                                  // 0, 0
	nox_window_flags flags;                                  // 1, 4
	int width;                                               // 2, 8
	int height;                                              // 3, 12
	int off_x;                                               // 4, 16
	int off_y;                                               // 5, 20
	int end_x;                                               // 6, 24
	int end_y;                                               // 7, 28
	void* widget_data;                                       // 8, 32; different types
	nox_window_data draw_data;                               // 9, 36
	unsigned int field_92;                                   // 92, 368
	int (*field_93)(nox_window*, int, int, int);             // 93
	int (*field_94)(nox_window*, int, int, int);             // 94, 376
	int (*draw_func)(nox_window*, nox_window_data*);         // 95, 380, second arg is &field_9
	int (*tooltip_func)(nox_window*, nox_window_data*, int); // 96, 384
	nox_window* prev;                                        // 97, 388
	nox_window* next;                                        // 98, 392
	nox_window* parent;                                      // 99, 396
	nox_window* field_100;                                   // 100, 400
} nox_window;
_Static_assert(sizeof(nox_window) == 404, "wrong size of nox_window structure!");
typedef struct nox_window_ref nox_window_ref;
typedef struct nox_window_ref {
	nox_window* win;
	nox_window_ref* next;
} nox_window_ref;

#endif // NOX_CLIENT_GUI_WINDOW_H
