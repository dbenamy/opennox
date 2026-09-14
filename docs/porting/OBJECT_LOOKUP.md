# Object lookup and net-code cache

Thirteen routines are now native Go. The conversion removes 438 physical C
lines from GAME3_2.c, GAME3_3.c and vardefs.c, bringing production C to **99,071
lines in 147 files**, with zero reference-only C. The qualified original-C
baseline is `55ab415b`.

The cache's sixteen nodes, list heads and initialization flag are owned by Go.
Only net-code lookup and script-ID lookup retain C exports for existing callers.
Eleven private C entry points, the C cache storage and its private C types are
removed. Go callers and the objective test fixture use the native owner directly.
The historical address interval audit found only memmap metadata, with no direct
accesses to that storage. The unrelated shadow-list routines remain outside scope.

The original C baseline has **1,074 cases /eight repeated locked captures**, plus
a focused probe. Cases use real main, pending, missile and player lists and
capture every cache node and list head/tail after each operation. Names cover nil
versus empty IDs, embedded NUL, case, non-ASCII bytes, qualified names, multiple
colons and one inventory level. Search cases distinguish pending inventories,
missiles and player fallback, destroyed objects and full-width IDs. Independent
ordered-slice oracles check promotion, capacity, eviction, invalidation, reset
and flush; eight mixed cases contain 1,280 operations. Mutation cases cover changed
codes, destroyed cached objects and duplicate cache entries.

The native focused lookup/objective comparison passes in 11.506s with every hash
unchanged. The accumulated **85,333 captured cases /1,029 groups**, plus applicable
contracts, pass standard/server/highres. The guarded runner verifies that all
Go-discovered root tests execute and complete:

| Variant | Root tests | Wall time |
| --- | ---: | ---: |
| Standard | 623 | 357.408s |
| Server | 622 | 411.876s |
| High-resolution | 623 | 346.540s |

All three production builds pass and verify as i386/SSE2, with the two required
exports present, eleven retired symbols and C cache globals absent, and no test
helpers. The full suite matches all 1,553 known failure entries exactly (15 passing,
three failing, 32 skipped packages). Fresh unchanged Xvfb/null-audio gameplay
passes in 35.566s with golden override disabled. Builds ran alongside the variant
matrix using separate outputs and fixed source.

Review later: six side-effect-only private helpers have unused C return values
(prepend, remove, initialize, add, invalidate and flush). All callers discard
these values or propagate them into another discarded return. Their native
helpers are void; tests compare their complete state effects. Meaningful lookup
and free-node results remain checked.

The existing server.Object name matcher has different colon and empty-ID
semantics, so the native path uses a dedicated compatibility helper. It also
preserves cached destroyed-object hits, uncached player fallback and
initialization's retained node values/stale free links.

Local evidence is under build/port-object-lookup: qualification.json, captures.json,
ABI/address audits, guarded matrix results, binary-verification.json and
full-suite-comparison.json. Captures have verified gzip archive manifests.
Gameplay evidence is build/baseline/runs/object-lookup-port/result.json.
