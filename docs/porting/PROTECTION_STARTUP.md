# Protection initialization — 2026-09-10

Startup now runs in native Go through the existing public Sub_56F1C0 wrapper.
It seeds the unchanged C floating RNG from wall-clock Unix seconds (low signed
32 bits), resets head/tail/count/key/sequence, XORs one floating draw into the
frame-derived key, then creates the initial zero record, seven reserved zero
records and final one record. The two handle slots and final private helper
return retain the original IDs. It preserves swap/rekey counters, consumes eight
Logic RNG draws for insertion and leaves Other RNG unchanged. As before, this
initializer expects callers to have handled the previous manager; it does not
free an existing list before discarding its head.

Its sole production caller was the public Go wrapper in session startup, so the
C initializer entry/declaration are removed rather than retained for tests. The
five now-unused C manager declarations and time.h include are removed too.
The floating seed/draw helpers remain C for the next bounded conversion.

Original-C baseline: `878cd8bc`, 200 scenarios through both the C entry and public
wrapper (400 runs). Tests independently construct the nine expected records and
randomized order, check both handle slots, checksum/key/sequence/count, complete
links, final return, counters and both RNG indices. Frames include zero, one,
signed boundaries and maximum uint32; RNG seeds/counters include negative and
wrap boundaries. Stale manager fields and slots must be reset. Allocation
failure injection is not added; the existing constructor contracts are retained.

The fixture brackets wall-clock seconds and compares the final floating state
against unchanged C seed/draw helpers for candidates in that interval. An
anomalous clock jump triggers a bounded retry rather than a guessed seed.
It saves/restores all manager fields, slots, counters, server and floating state,
including the five shipped constants from blob_581450.dat. Explicit constants
are necessary because ordinary unit-test setup leaves those blob bytes zero;
the test asserts the final floating state is nonzero to retain that coverage.

Production C after conversion: **141,914 physical lines (−27)** in 153 files;
test-reference C: **0**. Local artifacts: `build/port-init/`.

All accumulated protection tests pass on 386 default/server/highres; all three
production binaries build. Symbol checks confirm the native initializer and
absence of the retired C entry. The `init-port` warrior scenario accepts both
preserved screenshots with overrides disabled. The last full-suite milestone
was the object chunk, matching the known failures exactly.
