# Protection reserved records and handle allocation — 2026-09-10

This chunk replaces `sub_56F250` and `nox_xxx_protectionCreateInt_56F400`.
Both retain C exports because manager initialization (`sub_56F1C0`) still calls
them. Reserved initialization attempts exactly seven zero-value records and
advances the uint32 sequence after every attempt, returning only the last
constructor result. Ordinary allocation returns the old sequence and increments
it only after successful creation; failure returns zero without an increment.
A successful allocation at sequence zero also returns zero, preserving the
original ambiguity. The failure diagnostic is an empty callback.

The existing root Go `protectInt`/`protectFloat32` wrappers already implement
handle allocation for player fields. This chunk keeps those paths scoped apart
from the remaining C initialization callers; it does not change their APIs.

Terra prepares the pre-port state tests; the primary agent reviews the test
oracle and implements/integrates the conversion. Production records remain
C-allocated. Forced allocator exhaustion is not exercised; the constructor's
nil result path is reviewed along with the distinct sequence increment rules.
Final validation results and source counts follow after the checks.

Pre-port tests are committed at `20f6ffb8`. The same 400 deterministic sequences
pass against original C and Go. They mix ordinary and reserved allocations from
empty/prepopulated lists with zero/random keys and sequence values at zero, the
handle threshold, signed boundaries and uint32 wraparound. Every operation
checks the full ordered record contents, links/endpoints, checksum, count,
sequence, raw return bits and both RNG positions.

All accumulated `^TestProtection` tests pass with `porttest`, `server porttest`
and `highres porttest` on 386. Logs: build/port-handles/go-after-*.log.
Production C: **142,327 physical lines**, down **24**, in 153 files; C references:
**0**. See [C_LOC.md](C_LOC.md).

All three production targets build through `go run ./internal/noxbuild -o
../build/port-handles/bin`. Symbol checks verify both retained C entries now have
Go export bridges. Fresh `handles-port` warrior gameplay exits 0 with both
preserved screenshots accepted and overrides disabled. Existing unrelated
full-suite failures remain documented in the earlier checkpoint.
