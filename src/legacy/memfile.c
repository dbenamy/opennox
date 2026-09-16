#include "memfile.h"
#include "common__binfile.h"
#include <stdlib.h>
#include <string.h>

int8_t nox_memfile_read_i8(nox_memfile* f) {
	if (!f->data)
		return 0;
	int8_t v = *(int8_t*)f->cur;
	f->cur++;
	return v;
}

uint8_t nox_memfile_read_u8(nox_memfile* f) {
	if (!f->data)
		return 0;
	uint8_t v = *(uint8_t*)f->cur;
	f->cur++;
	return v;
}

int16_t nox_memfile_read_i16(nox_memfile* f) {
	if (!f->data)
		return 0;
	int16_t v = *(int16_t*)f->cur;
	f->cur += 2;
	return v;
}

uint16_t nox_memfile_read_u16(nox_memfile* f) {
	if (!f->data)
		return 0;
	uint16_t v = *(uint16_t*)f->cur;
	f->cur += 2;
	return v;
}

int32_t nox_memfile_read_i32(nox_memfile* f) {
	if (!f->data)
		return 0;
	int32_t v = *(int32_t*)f->cur;
	f->cur += 4;
	return v;
}

uint32_t nox_memfile_read_u32(nox_memfile* f) {
	if (!f->data)
		return 0;
	uint32_t v = *(uint32_t*)f->cur;
	f->cur += 4;
	return v;
}

void nox_memfile_skip(nox_memfile* f, int n) {
	if (!f->data)
		return;
	f->cur += n;
}

//----- (0040ACC0) --------------------------------------------------------
unsigned int nox_memfile_read(void* dst, const unsigned int sz, const int cnt, nox_memfile* f) {
	unsigned int n = cnt * sz;
	if (f->cur + n > f->end)
		n = f->end - f->cur;
	memcpy(dst, f->cur, n);
	f->cur += n;
	return n / sz;
}
