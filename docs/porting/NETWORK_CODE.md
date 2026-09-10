# Client unit-code encoding and bit helpers — 2026-09-10

Three C functions (578B00/30/70) now execute Go. The client encoder reads the
existing typed client.Drawable layout, returns zero for nil or full uint32 codes
at least 0x8000, and adds bit 15 for class bits 22 or 29. It leaves the drawable
unchanged. The bit helpers preserve signed-short input truncation followed by
0x7fff masking, and unsigned-int bit-15 extraction.

Existing C signatures are retained. In particular the encoder accepts the
legacy 32-bit integer address at its boundary, converting those raw bits to a
Drawable pointer before field access. Its 13 decompiled C call sites remain ABI
compatible. The many C bit-helper callers retain their symbols; the pre-existing
root Go bit helpers remain native and are cross-checked in the tests. Dynamic
extent resolution (578B40) is the next separate chunk.

Original-C baseline: `23c84897`. All 65,536 low-word inputs are tested with five
upper-word patterns, yielding 327,680 paired C bit-helper checks plus checks of
the existing native Go helpers. The client encoder has 2,988 cases covering nil,
full-word code boundaries, each class bit and its complement, mixed classes,
and 2,048 deterministic random cases. A C-allocated typed Drawable fixture
checks its entire byte representation remains unchanged. Test oracles use
arithmetic division/modulo and individual class-bit tests independently of the
implementation masks.

Production C: **141,821 physical lines (−23)** in 153 files; test-reference C: **0**.
Local artifacts: `build/port-netcode/`.

All accumulated protection and new network-code tests pass on 386 default,
server and highres. All three production binaries build and expose the three
C entries through Go-backed bridges. The `netcode-port` warrior scenario accepts
both preserved screenshots with overrides disabled. The immediately preceding
RNG full-suite milestone matched the known failure set exactly.
