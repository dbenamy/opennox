# Client inventory queries and state

Status: converted and qualified. The corrected C baseline is `0d563b10`.
The batch replaces 31 routines / 396 C lines: equipped-item and stack searches,
slot/type counts, scalar controls, secondary weapon queries, durability/charge
updates and potion use. The larger inventory-mode transition sub_467650 is
excluded. Production C remaining: **86,422 lines / 97 files / zero reference C**.

## Native qualification

| Check | Result |
| --- | --- |
| Frozen comparison | All 7,137 results / 6 groups unchanged |
| Focused native tests | 11 passed, 171.929s; includes seven independent contracts |
| Accumulated default repeat | All 805 root tests passed, 420.943s |
| Server variant | All 171 selected root tests passed, 222.455s |
| Highres variant | All 172 selected root tests passed, 87.043s |
| Three production binaries | ELF32/i386/SSE2/CGO; 26 native C interfaces, five retired, no test helpers |
| Full asset suite | Exact 1,553 known failure entries; 15 pass / 3 fail / 32 skip packages |
| Fresh warrior gameplay | Passed, 36.119s; reference comparison enabled, null audio |
| Source integrity | All 1,524 Go/C/header fingerprints unchanged throughout checks |

The first accumulated run had one hallway-route capture mismatch. It did not
recur in hallway-alone, inventory-plus-hallway, the exact 380-test preceding
prefix, or the complete 805-test repeat. Its cause remains **unexplained**; it is
not claimed fixed or attributed to scheduling. No inventory code, test input or
expected hash changed during this investigation. Keep the failure and repeat
evidence and investigate any recurrence. A separate diagnostic follow-up will
save full hallway mismatch captures automatically.

Evidence: build/port-client-inventory/native-qualification.json includes the
initial failure and all isolation results; native-capture-comparison.json,
native-binary-verification.json and native-default-repeat-*.json preserve final
comparisons. C baseline binary/captures and isolated hallway captures have
verified gzip archives with restore instructions. Test/build elapsed times are
not runtime performance benchmarks. Accumulated frozen coverage is now 447,275
results / 1,136 groups, plus independent contracts.

## Boundary correction before baseline

sub_461EF0 scanned row indices 0 through21 with an inclusive <=21 limit, while
the actual array declares21 rows x4 columns (84 cells). The last iteration
aliases later columns and reaches cell84 outside the array. Use <21. The final
valid row20 remains searchable; the coordinate-based UI queries intentionally
expose only rows0..19. This is a local reversible correction under standing user
authorization. Do not execute the old out-of-array read as a golden.

## Fixture and planned acceptance

Reuse the actual meter window/renderer, player, item definitions, drawable,
inventory grid, equipment lists and outgoing-message owners. Own and restore
additional lookup-result, dragged-item, scalar/float and window-handle state.
Thin tagged calls enter the production C ABI. No reference algorithm is added.

Matrices cover all84 cells; empty,1,2,31,32-entry stacks; first/last/missing codes;
visible-grid boundaries; row-major search order over column-major storage;
nine equipped lists with varying lengths/flags; Bow selection; signed bytes,
float bits and byte indices; durability/charge updates and dragged-item fallback;
potion pause/cursor/message behavior and weapon selection. Independent contracts
check these behaviors without relying solely on frozen output. Pointer returns
are normalized only for owned objects/storage; the original signed low16 pointer
return from durability updates is asserted before normalization.

Initial c-a passed all 10 tests in 171.823s; fresh-process c-b passed in 8.901s.
All 7,137 results / 6 groups repeated byte-for-byte, covering every selected
operation. Six independent contracts include four separate contract tests and
assertions embedded in potion-use and weapon-selection matrices. Hashes are now
frozen and qualified: 171 client / 170 server / 171 highres tests pass in
83.106s / 182.665s / 85.329s. The production C client is ELF32/i386/SSE2/CGO
with all 31 selected interfaces and no test helpers. The full asset suite has
the exact known 1,553 failure entries (15 pass / 3 fail / 32 skip packages).
Fresh warrior gameplay passes in 35.614s with reference comparison enabled.
All 1,521 Go/C/header fingerprints are unchanged. Evidence: c-qualification.json. Current artifacts: build/port-client-inventory. Actual source
is authoritative; ignored drafts and caller-audit line numbers are historical.

After freezing: translate the connected scope, keep only required C interfaces,
run accumulated/default and affected server/highres checks, all production builds
and symbol checks, known full-suite failure comparison and fresh gameplay; update
C LOC and documents, commit/push and continue.

## Native draft review

The native draft uses a typed148-byte inventory cell: drawable pointer,32 item
codes, equipped/secondary fields, count/flags and tail word. Actual shared C
storage remains because unconverted inventory code uses it. Five interfaces
have no remaining C caller after this batch; retire them and route meter calls
directly to native helpers. Retain26 genuine C interfaces.

For a count above32, bound the native code search to the cell's32 code slots.
The old search reads beyond the code array for such a count. Valid counts0..32
retain their frozen behavior; an additional native contract must check malformed
counts without executing the invalid old C read. Other count getters preserve
the stored byte. Record this reversible choice for later review.

## Local artifact maintenance

The completed meter captures are verified lossless gzip archives with a restore
manifest. Low disk space also required pruning 2.00GiB of the oldest compiled Go
cache artifacts (102 files larger than10MiB, unused for at least24 hours by cache
mtime). This affects rebuild speed only; Go recreates cache entries. Current
source, module downloads, qualification evidence and original assets are intact.
The manifest is build/baseline/old-build-cache-prune.json. The full cache was not
cleared. Completed gameplay asset copies use the verified deduplication manifest.

## Native comparison checkpoint

Native-a passed all 11 selected tests in 171.929s; all 7,137 captured results
match the unchanged C hashes. Final qualification is complete; see the results and limitation above. Typed cells and lookup
records view the existing shared C storage. Meter calls now use native helpers;
the production C exports for potion use, first-item/type count, mode getter and
selected-weapon query are retired. A native-only contract checks the actual
148-byte cell/8-byte result layout, malformed count255, duplicate equipment
priority, and zero-code lookup versus ignored zero-code charge updates.

## Qualification scheduling trial

For native qualification, overlap the accumulated default suite with one
server-then-highres sequence. At most two test drivers run concurrently; the
production builder also runs with -p2. Tests use GOMAXPROCS=2. Each has separate
logs/results and actual test-start/completion checks. Always join every test job
and the builder, including failures, before source edits. Full asset tests and
fresh gameplay run afterward. This changes scheduling only, not coverage or
acceptance criteria, and uses no additional agents. The initial joined test phase took450.507s, but its one unexplained hallway
failure means this trial is not a clean acceptance of the scheduling change.
Keep it optional and preserve captures if repeated; do not weaken any test.

## First native qualification finding

All805 accumulated root tests executed and finished in450.417s. One failed:
TestMapHallwaysRoutes produced515216c24d12014d81f2633afd4257858be633eab2591023ec56daf2f972d0ec
instead of frozen712a46430f47c5d55a5f1bcbe9ef4dfc638fae8e205ddab827b6266c7a1bfdbc.
Other accumulated tests and server171/highres172 tests passed; all three builds
succeeded. The qualification driver joined all jobs and stopped before full
asset/gameplay checks. This is not accepted or attributed to scheduling yet.

With identical source, hallway-alone passes in19.840s and inventory-plus-hallway
passes all12 selected tests in20.754s. The exact380-test preceding prefix passed in193.911s, then the complete805-test
repeat passed in420.943s with full hallway captures. All three hallway groups
matched the original hashes. Frozen expectations remain unchanged.
