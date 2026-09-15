# Client inventory transactions

Status: converted and fully qualified. The corrected C baseline is `9676bcf0`.
All 26 routines / 1,025 C lines are replaced: stack allocation, append/removal and
compaction, equipment classification/lists, pickup/drop/drag/place, secondary
weapon changes, capacity and item requests. Window construction, rendering and
event dispatch remain separate. Production C is **85,399 lines / 97 files / zero
reference C**.

## Native qualification

| Check | Result |
| --- | --- |
| Frozen comparison | All 3,606 results / eight groups unchanged on first native comparison |
| Focused native tests | All 16 passed, 173.242s |
| Accumulated default | All 821 root tests passed, 431.028s |
| Server variant | All 187 selected root tests passed, 185.802s |
| Highres variant | All 188 selected root tests passed, 85.676s |
| Three production binaries | ELF32/i386/SSE2/CGO; 19 retained interfaces, seven retired; no test helpers |
| Full asset suite | Exact 1,553 known failure entries; 15 pass / 3 fail / 32 skip packages |
| Fresh warrior gameplay | Passed, 35.377s; comparisons enabled, null audio |
| Source integrity | All 1,537 Go/C/header fingerprints unchanged throughout qualification |
| C reduction | 86,424 → 85,399 lines, a reduction of 1,025 |

Four native files share the existing typed inventory storage and actual drawable,
net-list and meter implementations. Three production Go callers invoke Go directly.
The public pickup header drops its pointer's `const` qualifier to agree with the
CGO-generated declaration; pointer layout/calling convention and read-only input
behavior are unchanged. Native compaction uses typed cells without the old
transient out-of-array rolling pointers, retaining stack order and observed
return offsets. No C algorithm is retained solely for testing.

Evidence: build/port-client-inventory-transactions/native-qualification.json,
native-capture-comparison.json, native-source-fingerprints.json, native-*.jsonl
and native-binary-verification.json. These check durations are not gameplay
performance measurements. Accumulated frozen coverage is 450,881 results /
1,144 groups, plus independent contracts.

The earlier inventory-query hallway mismatch did not recur: all 821 accumulated
tests, including hallway routes, pass. Its cause remains unexplained; see
[CLIENT_INVENTORY.md](CLIENT_INVENTORY.md). Do not change that oracle on recurrence;
failed hallway comparisons automatically preserve full captures.

## Ownership and coverage

The fixture delegates actual drawable creation/deletion through the existing
client interface, records lifetime identities and verifies real allocator/list
counts. Case objects are separate from protected base templates. Genuine fixed-
pool exhaustion exercises allocation failure without mocking the constructor.
Inventory, meter, player, GUI, string manager, console and net-list owners are
real implementations with controlled inputs and restored state.

Eight groups cover stacks (1,488), equipment (408), compaction (64), capacity
(1,296), notifications (63), secondary selection (56), drag/place (129), and
pickup/follow-up transitions (102). Sixteen tests include independent contracts
for coordinates, stack ownership/code conservation, metadata/durability copies,
request bytes, real allocation failure, visible-grid priority, error-ring wrap
and expiry, equipment classification/list consistency, selection follow-ups,
sounds and cursor clearing. Captures include meter/potion state, game requests,
error text, lifetime events and allocator count deltas.

## C prerequisites to review

Existing-stack pickup in `nox_xxx_spritePickup_461660` previously read uninitialized
coordinates when recording the last weapon cell. Derive column and row from the
returned cell pointer before freezing. This adds two C lines and is checked
across every visible cell. Undefined stack contents are not compatibility
expectations.

The first compaction conservation test aborted in fortified `memcpy`. GDB found
a 148-byte copy from grid byte offset 1,776 to 148, both in bounds, while generated
`__memcpy_chk` received destination capacity zero. The old rolling pointer moves
before the array when retrying a merged stack; the generated bound remained zero
when the pointer advanced back. Recompute the destination from existing row/column
indices at the copy. Fortification stays enabled; source, length and traversal
are unchanged, with no C LOC change. Evidence:
build/port-client-inventory-transactions/compaction-failure-gdb.txt.

Both corrections preceded the frozen baseline and were made under the standing
policy for confident reversible changes. Layout conservation, full/empty grids,
and moving the selected cell all pass against corrected C and native Go.

## C baseline qualification

Affected tests pass 187 default / 186 server / 187 highres root tests
(72.020s / 185.245s / 84.839s). The production client verifies all 26 original C
routines and excludes helpers. The full asset suite matches the exact known
failure set, and fresh warrior gameplay passes in 37.103s. All 1,532 initial
source fingerprints stayed unchanged throughout those checks.

The final pickup-transition test was added afterward, with production and all
existing fixtures/oracles unchanged. Final 16-test default repeats pass in
21.122s / 9.137s; server and highres pass in 19.785s / 20.847s. All 3,606 results
match and all 1,533 final source fingerprints are unchanged. The test-only extension
is recorded separately in final-c-qualification.json; the earlier broad checks
are not presented as having included that final test. Production build, full-suite
and gameplay evidence is reused because production source did not change.

## Artifacts and recovery

The qualified C executable is verified in bin/opennox-c-baseline.gz. Completed
C/native gameplay copies were deduplicated only after success, with verified
asset hashes and restoration manifests. The earlier qualified hallway route
capture is a verified gzip archive with its oracle and complete data preserved.
Old compiled cache data was selectively pruned after all source-reading jobs
joined, freeing about 2 GiB; no original assets or active run data were removed.
The committed baseline and expectations remain sufficient to repeat qualification
if ignored artifacts are lost; use the repository recovery instructions.
