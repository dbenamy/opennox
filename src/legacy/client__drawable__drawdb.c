#include <stdio.h>

#include "client__draw__debugdraw.h"
#include "client__drawable__drawable.h"

#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME1_2.h"
#include "GAME2.h"
#include "GAME2_1.h"
#include "GAME2_2.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_3.h"
#include "client__gui__window.h"
#include "client__video__draw_common.h"
#include "common__strman.h"
#include "operators.h"

extern uint64_t qword_581450_9544;
extern uint64_t qword_581450_9552;
extern uint32_t dword_5d4594_251572;

extern int nox_parse_thing_draw_funcs_cnt;

void* nox_xxx_draw_44C780(int a1) {
	int i;        // esi
	int v2;       // eax
	void* result; // eax

	for (i = 0; i < 32; i += 4) {
		v2 = i;
		if (i >= 16) {
			v2 = i + 4;
		}
		result = *(void**)(v2 + a1);
		if (result) {
			free(result);
		}
	}
	return result;
}

void* sub_44C7B0(int a1) {
	void** v1;    // ebx
	int v2;       // ebp
	void** v3;    // esi
	int v4;       // edi
	void** v5;    // esi
	int v6;       // edi
	void* result; // eax

	v1 = (void**)(a1 + 52);
	v2 = 55;
	do {
		if (*v1) {
			nox_xxx_draw_44C780((int)*v1 + 4);
			free(*v1);
		}
		v3 = v1 + 1;
		v4 = 26;
		do {
			if (*v3) {
				nox_xxx_draw_44C780((int)*v3 + 4);
				free(*v3);
			}
			++v3;
			--v4;
		} while (v4);
		v5 = v1 + 27;
		v6 = 27;
		do {
			result = *v5;
			if (*v5) {
				nox_xxx_draw_44C780((int)result + 4);
				free(*v5);
			}
			++v5;
			--v6;
		} while (v6);
		v1 += 66;
		--v2;
	} while (v2);
	return result;
}

void nox_xxx_draw_44C650_free_kind(void* lpMem, int kind) {
	void** v7 = 0;
	int v8 = 0;
	char* v9 = 0;
	int v10 = 0;
	char* v11 = 0;
	int v12 = 0;

	switch (kind) {
	case 2:
	case 3:
		if (*((uint32_t*)lpMem + 1)) {
			free(*((void**)lpMem + 1));
		}
		free(lpMem);
		break;
	case 4:
		v7 = (void**)((char*)lpMem + 4);
		v8 = 5;
		do {
			if (*v7) {
				free(*v7);
			}
			++v7;
			--v8;
		} while (v8);
		free(lpMem);
		break;
	case 5:
		nox_xxx_draw_44C780((int)lpMem + 4);
		free(lpMem);
		break;
	case 6:
		sub_44C7B0((int)lpMem);
		free(lpMem);
		break;
	case 7:
		v9 = (char*)lpMem + 8;
		v10 = 16;
		do {
			nox_xxx_draw_44C780((int)v9);
			v9 += 48;
			--v10;
		} while (v10);
		free(lpMem);
		break;
	case 8:
		v11 = (char*)lpMem + 8;
		v12 = 3;
		do {
			nox_xxx_draw_44C780((int)v11);
			v11 += 48;
			--v12;
		} while (v12);
		free(lpMem);
		break;
	default:
		free(lpMem);
		break;
	}
}
