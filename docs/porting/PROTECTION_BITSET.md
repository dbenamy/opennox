# Protection spell/ability bitsets — 2026-09-10

This chunk replaces `sub_56FCB0`, `nox_xxx_playerAwardSpellProtectionCRC_56FCE0`
and `nox_xxx_playerApplyProtectionCRC_56FD50` from GAME5_2.c. Go helpers implement
the bit operations; legacy exports retain the C ABI and update existing C-owned
records and checksum state. Existing Go spell wrappers now call Go directly,
removing their otherwise unnecessary Go→C→Go round trip.

The compatibility contract intentionally preserves:

- Nonzero values set a bit; zero does not clear one already set.
- Nonnegative indices fold modulo 32. This is not a collision-free bitset.
- Array validation ignores entry zero and folds only entries 1 through count−1.
- Counts <= 1 inspect no array memory and validate against a zero decoded mask.
- Handle eligibility uses a signed comparison with 657757279. Missing/ineligible
  awards return the original handle; validation returns zero.
- Awards update the encrypted payload and XOR checksum using old/new payloads;
  they do not allocate, relink, rekey, consume randomness or mutate input flags.

Negative spell/ability indices are outside the defined C shift contract; the
current callers use nonnegative indices. No claim is made to preserve undefined
negative C shifts. The existing 386 record layout assertions remain applicable.

## Tests and recovery

Pre-conversion tests are committed at `96dec42b` and passed against C. They run
5,000 deterministic state scenarios and 12,291 direct bit checks: disabled and
signed truthy values, duplicate folded bits, index/word boundaries, signed handle
thresholds, empty/missing records, zero/high-bit keys, successes and mismatches,
null array pointers for count <= 1, exact checksum deltas and unchanged input/
non-payload state. The same assertions run after conversion through the actual
C ABI. Fixed Go unit cases cover the bit contract independently. No C reference
copy is added; recover the old path from the pre-conversion test commit.

Commands from src with RECOVERY.md's environment:

```bash
go test -count=1 ./internal/protection ./internal/e2etest
go test -tags porttest -count=1 -run '^TestProtection(BitsetABI|RecordsABI|ABI)$' .
go test -tags 'server porttest' -count=1 -run '^TestProtection(BitsetABI|RecordsABI|ABI)$' .
go test -tags 'highres porttest' -count=1 -run '^TestProtection(BitsetABI|RecordsABI|ABI)$' .
go run ./internal/noxbuild -o ../build/port-bitset/bin
```

Local evidence is under build/port-bitset. This extends coverage for these
operations, not the entire protection manager or every spell/gameplay path.
Production C: **142,503 physical lines**, down **67** this chunk, 153 files.
Test-reference C: **0**. See [C_LOC.md](C_LOC.md).

Post-conversion checks pass for default/server/highres on 386 and pure Go tests
on amd64. All three binaries build and contain all three Go export bridges.
The fresh standard-client warrior scenario bitset-port exits 0 against both
preserved screenshots with overrides disabled. Known unrelated suite failures
are unchanged; this chunk used the focused tests and integration/build checks.
