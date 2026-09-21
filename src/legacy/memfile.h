#ifndef NOX_MEMFILE_H
#define NOX_MEMFILE_H

#include <stdint.h>

typedef struct {
	char* data; // 0, 0x0, 0
	int size;   // 1, 0x4, 4
	char* cur;  // 2, 0x8, 8
	char* end;  // 3, 0xC, 12
} nox_memfile;


#endif // NOX_MEMFILE_H
