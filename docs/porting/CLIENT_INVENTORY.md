# Client inventory queries and state

Status: corrected C baseline frozen and fully qualified; ready for conversion. Meter conversion
`30553db2` is pushed. This batch selects 31 routines / 396 C lines from GAME2_1.c:
equipped-item and inventory-stack searches, slot/type counts, scalar controls,
secondary weapon queries, durability/charge updates and potion use. The larger
inventory-mode transition sub_467650 is excluded. No conversion is applied yet.

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
