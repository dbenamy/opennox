# Object lookup and net-code cache

The next connected batch replaces 13 routines spanning 436 C section lines in
GAME3_2.c and GAME3_3.c. Production code is still C during baseline capture.
The current production count remains 99,509 lines in 147 files, with zero
reference-only C lines. Original production code is recoverable at 8acac7e5.

The guarded fixture exercises real main, pending, missile and player lists and
captures all sixteen cache nodes, both list heads/tails and the initialization
flag after each operation. Names cover nil versus empty IDs, embedded NUL,
case, non-ASCII bytes, qualified names, multiple colons and one inventory level.
Lookup tests distinguish main inventories, pending inventories, missiles and
player fallback; cover destroyed objects and full-width IDs. Independent ordered
slice oracles check cache promotion, capacity/eviction, invalidation, reset,
flush and mixed sequences. Mutation cases cover changed object codes, destroyed
cached objects and duplicate cache entries.

Only net-code lookup and script-ID lookup have remaining outside C callers.
Eleven entry points can retire; existing Go wrappers will call native helpers.
No production Go address references were found. The cache types are private.
A source scan for both historical address intervals and their interior offsets
found only the two memmap metadata annotations; there are no direct accesses to
that historical storage. Local audit: build/port-object-lookup/abi-audit.json and
global-address-audit.json.

Review later: six side-effect-only private helpers have unused C return values
(prepend, remove, initialize, add, invalidate and flush). All callers discard
these values or propagate them into another discarded return. The fixture
compares their full state effects and deliberately ignores their returns. Native
helpers will therefore be void. This avoids preserving decompiler return values
solely for tests; meaningful lookup and free-node results remain checked.

The existing server.Object name matcher has different colon and empty-ID
semantics. Use a dedicated compatibility helper rather than replacing this path
with it. Preserve one-level inventory traversal, cached destroyed-object hits,
uncached player fallback and initialization's retained node values/stale free
links. The unrelated shadow-list routines are outside this batch.

The original C produces 1,074 cases in eight locked, repeated SHA-256 capture
groups, plus a focused lookup probe. Eight mixed cases contain 1,280 operations.
Focused checks pass in 5.490s. The accumulated standard baseline passes with
623 Go-discovered root tests actually executed and completed (331.395s wall);
affected server/highres each pass all 91 selected root tests (172.723s/69.908s
wall). Every selected root test executed and completed. The baseline is ready
for conversion. All subsequent accumulated
matrices use tools/porting/run_tests.py, which verifies Go-discovered tests
actually execute and finish. Generated captures and raw logs remain local under
build/port-object-lookup.
