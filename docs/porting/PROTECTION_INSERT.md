# Protection randomized insertion — 2026-09-10

This chunk replaces `sub_56F2F0` in GAME5_2.c. Its only production caller was
the Go record constructor, so insertion becomes a private Go call and the C
function and declarations are removed. Records remain C allocations with the
existing asserted 386 layout.

An empty list receives its first record without consuming randomness. A nonempty
list consumes exactly one Logic RNG draw, including a singleton, and inserts
before the chosen record. Other RNG state and the existing tail are unchanged.
The uint16 record count retains its wrapping behavior. The removed zero-key
callback was an empty diagnostic function.

Pre-port tests are recoverable at `06f1cec9`; no separate C reference is kept.
The same tests run against original C and the conversion: 500 deterministic
sequences cover empty/single/multiple records, zero/random keys, integer/float
construction, exact decoded payload order, forward/back links, endpoints,
checksum and both RNG indices after every insertion. Two prepopulated lists
cover 32,768 and 65,535 records, including count wrap to zero. Pure helper tests
check insertion before the head, middle and tail while preserving payloads.
Valid allocated records and consistent acyclic lists remain preconditions.

The implementation was delegated to Terra as a bounded trial; the primary agent
reviewed the original C semantics, caller audit and test coverage and owns final
validation. No measured billing savings are claimed.

Validation artifacts are local under `build/port-insert`; final outcomes and C
line counts are recorded after checks complete.

All accumulated `^TestProtection` tests pass with `porttest`, `server porttest`
and `highres porttest` on 386. Pure protection tests pass on 386/amd64; E2E oracle
unit tests pass. `go run ./internal/noxbuild -o ../build/port-insert/bin` builds all
three targets. Symbol checks confirm the obsolete insertion entry is absent and
the required integer-constructor entry remains. Fresh `insertion-port` warrior
scenario exits 0 against both preserved screenshots, with overrides disabled.
Existing unrelated full-suite failures remain as previously documented.

Production C: **142,351 physical lines**, down **42**, in 153 files; test-reference
C: **0**. See [C_LOC.md](C_LOC.md).
