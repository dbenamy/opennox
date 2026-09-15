# Client inventory transactions

Status: C baseline qualified and ready to commit. Production remains C.
Previous completed query conversion: `019a5d26`; failure diagnostics: `e3b47068`.

Scope: 26 routines covering stack allocation, append/removal/compaction, equipment
classification and lists, pickup/drop/drag/place, secondary weapon changes and
item requests. Window construction, rendering and event dispatch remain separate.
The audited scope is 1,025 C lines including the prerequisite below.

## Prerequisite to review

Existing-stack pickup in `nox_xxx_spritePickup_461660` previously read uninitialized
local coordinates when recording the last weapon cell. Derive column and row
from the returned cell pointer. This adds two C lines before baseline capture;
production C is 86,424 lines in 97 files, with no test-reference C. Undefined
stack contents are not compatibility expectations. A grid-wide contract
exercises the corrected path.

## Ownership and coverage

Reuse real inventory, meter, player, strings, GUI and net-list owners. Observe
actual drawable creation/deletion through the existing client interface while
delegating to production implementations. Track case object identities and pool
counts, preserve the base fixture templates, and use genuine fixed-pool exhaustion
for failure cases. No test implementation of the C algorithm is retained.

The corpus has 3,606 results in eight frozen groups: stacks (1,488), equipment
(408), compaction (64), capacity (1,296), notifications (63), secondary selection
(56), drag/place (129), and pickup/follow-up transitions (102). Sixteen tests
include independent checks for coordinates, stack ownership and code conservation,
metadata/durability copies, message bytes, real allocation failure, visible-grid
priority, error-ring wrap/expiry, equipment classification and list consistency,
selection follow-ups, sounds and cursor clearing. Captures include meter/potion
state, game requests, error text, lifetime events and allocator count deltas.

The last transition test was added after broader C qualification. Production and
all existing fixtures/oracles are unchanged. Repeat all 16 focused tests in default,
server and highres; reuse the already-qualified production build, full-suite and
gameplay evidence. Record this test-only extension separately rather than claiming
the earlier broad run included it. Qualify and commit before translation.

The earlier isolated hallway qualification mismatch remains unexplained; see
CLIENT_INVENTORY.md. Its oracle is unchanged and failure captures are automatic.

### Compaction copy prerequisite

The first real compaction conservation test aborted in fortified `memcpy`.
GDB found a 148-byte copy from grid byte offset 1,776 to 148, both in bounds,
but the generated `__memcpy_chk` received destination capacity zero. The original
rolling destination pointer temporarily moves before the array when retrying a
merged stack; compiler object-size tracking then retains zero capacity after the
pointer is advanced back into the array. Recompute the copy destination from
its existing row/column indices at the copy. Fortification remains enabled;
copy length, source, traversal and algorithm are unchanged. No C LOC change.
Evidence: build/port-client-inventory-transactions/compaction-failure-gdb.txt.
This correction precedes all frozen captures. All 64 layout conservation cases,
full/empty grids, and moving the selected cell pass against corrected C.

## C qualification

The affected suite passes 187 default / 186 server / 187 highres root tests
(72.020s / 185.245s / 84.839s). The production client is ELF32/i386/SSE2/CGO;
all 26 C routines are present and test helpers are absent. The full asset suite
has the exact known multiset: 1,553 failure entries, 15 pass / 3 fail / 32 skip
packages. Fresh warrior gameplay passes in 37.103s with comparisons enabled.
All 1,532 initial source fingerprints remained unchanged throughout those checks.

Final 16-test default repeats pass in 21.122s / 9.137s; server and highres pass
in 19.785s / 20.847s. All 3,606 results match exactly and all 1,533 final source
fingerprints are unchanged. The test-only transition extension is recorded in
build/port-client-inventory-transactions/final-c-qualification.json.
Its schema explicitly distinguishes those checks from the earlier broader runs.

## Conversion plan

Replace all 26 routines together (1,025 C lines); retain 19 real C interfaces and
retire seven private helpers. Move three production Go callers directly to Go.
Use the shared typed inventory cells, actual drawable factory/deletion, existing
native meter helpers and real net-list owner. Preserve traversal, stack-code
ordering, item metadata, list ordering, and the legacy derived-pointer returns.
Do not regenerate frozen expectations during translation.
