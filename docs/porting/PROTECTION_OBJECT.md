# Protection object checksum and toggles — 2026-09-10

The object checksum now reads the existing Go Object layout: uint32 NetCode,
uint16 TypeInd and optional uint16 current health, plus word-only checksums of
initialization data and NUL-excluded object ID bytes. The initialization size is
interpreted as signed int32, matching the former C getter; absent types and
nonpositive sizes skip the data. Nil objects return zero. Object and buffer
contents remain read-only.

Both toggle C entries remain exported. They use signed handle eligibility and
first-match lookup. Both an ineligible and a missing ID return the original raw
ID; a successful call XORs the digest into the stored value and manager checksum
and returns the updated checksum. Repeated identical toggles restore the value
and checksum. No rekey/RNG/list/key/sequence changes occur.

The sole-use C field getters (4E4C00/10/30/80), Go size getter C export (4E4C50),
standalone object-checksum C entry, and FAC0 checksum bridge are retired. FAE0
remains exported for its remaining C caller; its ABI tests and independent byte
oracle remain. The FC50 internal header loses its const qualifier to match the
generated Go export; its caller objects are compatible and behavior is read-only.

Original-C tests are recoverable at `e4127e22`. 520 scenarios run through both
toggle entries, with direct digest comparisons and 1–3 toggles per run. The
C-allocated object fixture has a porttest-only temporary type-table helper.
Checks cover nil object/health/data/name, zero and high-bit fields, type index
zero and the missing-type sentinel, exact/trailing data bytes, embedded NULs,
lookup misses/duplicates, repeated restoration and unchanged object/buffer bytes.
Guard pages prove that rejected/missing IDs do not evaluate objects, and that
missing types, nonpositive sizes or partial words do not read initialization data.
Full manager/list/RNG snapshots remain checked throughout.

Final C baseline: build/port-object/c-before-final.log. Production C after the
conversion is **141,984 physical lines**, down **131**, in 153 files; C references
**0**. See [C_LOC.md](C_LOC.md). All accumulated protection tests pass on 386 default/server/highres; all three
production binaries build. Symbol checks confirm both toggles are Go-backed and
the retired checksum/getter entries are absent in all three binaries. The
`object-port` warrior scenario accepts both preserved screenshots with overrides
disabled. Full-suite results match the rekey baseline exactly: 15 passing, 3
known failing and 32 skipped/no-test packages, with no changed failure entries.
Local verification artifacts are under `build/port-object/`.
