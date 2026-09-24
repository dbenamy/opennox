#ifndef NOX_PORT_CLIENT_SYSTEM_CTRLEVNT
#define NOX_PORT_CLIENT_SYSTEM_CTRLEVNT

#include "common__savegame.h"
#include "defs.h"
#include "static_assert.h"

#define NOX_CTRLEVENT_KEYS_NUM 8
typedef struct nox_ctrlevent_key_t nox_ctrlevent_key_t;
typedef struct nox_ctrlevent_key_t {
	unsigned int keys[NOX_CTRLEVENT_KEYS_NUM];    // 0, 0
	unsigned int keys_cnt;                        // 8, 32
	unsigned int field_9[NOX_CTRLEVENT_KEYS_NUM]; // 9, 36
	unsigned int field_9_cnt;                     // 17, 68
	nox_ctrlevent_key_t* field_18;                // 18, 72
	nox_ctrlevent_key_t* field_19;                // 19, 76
	nox_ctrlevent_key_t* field_20;                // 20, 80
	nox_ctrlevent_key_t* field_21;                // 21, 84
	unsigned char field_22_0;                     // 22, 88
	unsigned char field_22_1;                     // 22, 89
	unsigned short field_22_2;                    // 22, 90
	unsigned int frame;                           // 23, 92, // TODO: game frame
} nox_ctrlevent_key_t;
_Static_assert(sizeof(nox_ctrlevent_key_t) == 96, "wrong size of nox_ctrlevent_key_t structure!");


void sub_42CD90();

#endif // NOX_PORT_CLIENT_SYSTEM_CTRLEVNT
