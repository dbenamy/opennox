//go:build porttest

// Verbatim functions (apart from names) from GAME5_2.c at 0e9d2e1f.
// Test-only historical reference; not linked into ordinary builds.
_Static_assert(sizeof(int) == 4, "reference needs 32-bit int");

int reference_checksum(int* a1, unsigned int a2) {
	int* v2;        // ecx
	int result;     // eax
	unsigned int i; // edx
	int v5;         // esi

	v2 = a1;
	result = 0;
	for (i = a2 >> 2; i; --i) {
		v5 = *v2;
		++v2;
		result ^= v5;
	}
	return result;
}

//----- (0056FAE0) --------------------------------------------------------
int reference_checksum_nullable(int* a1, unsigned int a2) {
	int result; // eax

	result = 0;
	if (a1) {
		result = reference_checksum(a1, a2);
	}
	return result;
}

