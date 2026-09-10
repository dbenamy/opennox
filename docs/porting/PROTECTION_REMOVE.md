# Protection record deletion and cleanup — 2026-09-10

This chunk replaces the delete-by-ID implementation, delete-and-clear wrapper,
and manager cleanup from GAME5_2.c. Only `sub_56F4F0` (delete-and-clear) still has
C callers, in the player cleanup path. It retains a C-to-Go export. The public Go
cleanup wrapper calls Go directly; obsolete internal C declarations/symbols for
`sub_56F510` and `sub_56F3B0` are removed.

Go finds the first matching record, unlinks it, updates both endpoints and the
XOR checksum, decrements the uint16 count, and frees the original C allocation.
A missing record changes nothing. The caller's handle becomes zero only after a
successful deletion. Cleanup saves each next pointer before freeing the current
record, then resets checksum/count/key/endpoints. It preserves the handle sequence.
C calloc/free ownership remains unchanged; the 386 record layout assertions apply.

Pre-conversion tests are committed at `56ad678c`. The same 1,000 deterministic
operation sequences run against C before replacement and Go afterward. They
verify full surviving payloads, forward/back links, head/tail, checksum, count,
handle values and sequence state after each operation, then cleanup. Cases cover
empty/singleton/multiple lists, head/middle/tail removal, duplicate IDs, misses,
repeated deletion, zero/UINT_MAX handles, zero keys and uint16 counter wrap. The
fixtures restore globals and own individual C allocations. Valid acyclic lists
and non-null handle pointers remain required; undefined invalid-pointer behavior
is not part of the compatibility contract. No C reference copy is added.

Commands from src with RECOVERY.md's target environment:

```bash
go test -count=1 ./internal/protection ./internal/e2etest
go test -tags porttest -count=1 -run '^TestProtection(RemoveABI|CreateABI|BitsetABI|RecordsABI|ABI)$' .
go test -tags 'server porttest' -count=1 -run '^TestProtection(RemoveABI|CreateABI|BitsetABI|RecordsABI|ABI)$' .
go test -tags 'highres porttest' -count=1 -run '^TestProtection(RemoveABI|CreateABI|BitsetABI|RecordsABI|ABI)$' .
go run ./internal/noxbuild -o ../build/port-remove/bin
```

Post-conversion checks pass under default/server/highres on 386; pure Go tests
also pass on amd64. All three production binaries build. Local evidence is under
build/port-remove. This validates record removal and cleanup, not the rest of the
manager's random insertion/rekey logic. Existing unrelated suite failures remain.
Production C: **142,393 physical lines**, down **65**, in 153 files, with **0**
test-reference C lines. See [C_LOC.md](C_LOC.md).
Fresh warrior scenario `remove-port` exits 0 against both preserved screenshots
with overrides disabled. Symbol checks confirm only the needed delete-and-clear
C entry remains among these operations.
