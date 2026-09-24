# Server fixture C bridges

## Scope and original behavior

Retire 81 C export bridges across 23 legacy owner files whose remaining in-tree
callers are test fixtures. Production algorithms already have Go owners. The
change moves fixture calls and identity maps to those owners, preserving sparse
operation IDs, captured expectations, arguments, results and side effects.

The whole-source audit found 446 reference rows. The primary independently matched
them across 3,063 tracked source files, including C preambles and headers. Five
quest names also identify independent root-package Go functions/methods; those do
not call the C bridges. Preserve those owners. This audit establishes in-tree
reachability, not an external ABI promise.

Eighteen production files may lose their C imports after conversion. Five retain
other C dependencies. Measure the final selected-profile inventory after
qualification. Standalone production and test-reference C remain zero.

[Reachability](server-fixture-bridges-reachability.json),
[test selection](server-fixture-bridges-tests.txt),
[qualification manifest](server-fixture-bridges-batch.json).

## Baseline and coverage

The original bridge source is qualified modifier commit `99b65896`.
The reviewed selection contains 288 root tests. Reuse the source-identical full
default corpus for all 288 and focused server/highres results for 57 each. Run the
remaining 231 highres and 230 server roots using the fingerprinted original
binaries. Verify exact test-name sets, source, supplemental inputs, binary hashes
and environment before accepting this baseline. The field-of-view test is explicitly excluded by `porttest && !server`;
server coverage therefore contains 287 roots. No new source contracts or
modified assertions are currently required by the boundary review.

Existing tests cover direct wrapper results and mutations, signed/truncated words,
float bit patterns, pointer identity, state restoration and guarded storage.
Nested world/objective/projectile and attack fixtures also populate the affected
identity maps even when their operations use registered production dispatch.

Baseline accepted: exact 288-root client / 287-root server coverage, with no selected
skips or failures. Additional 231-root highres and 230-root server runs passed on unchanged source.
[Baseline evidence](server-fixture-bridges-baseline.json). The conversion has completed qualification; results are recorded below.

## Decisions and review points

- Keep unknown/live C callback fallbacks and independent observers. Test-only
  identities may become distinct stable keys only after every consumer is traced.
- Preserve signed C-char results, short truncation, pointer-shaped 32-bit results,
  single/double precision boundaries and the order of guard checks and reads.
- Preserve nil rule filenames separately from empty filenames; retain C-string
  termination semantics for strings passed through the old adapters.
- Keep map-section error flag writes/logging, protection-ID clearing on successful
  deletion, and rules-file removal path construction and success results.
- Production algorithms and existing assertions/captures stay unchanged. Moving a
  wrapper's small boundary behavior to a fixture adapter does not authorize an
  algorithm rewrite or regenerating expected results.
- Run the full affected selection in all three profiles after conversion, then
  safe/static, three production/ABI builds, exact known-suite comparison, fresh
  headless save/load/resume and original-asset integrity. Reassess whether a full
  corpus is needed if the actual diff reaches shared production dispatch.

## Delegation

One GPT-6 Luna helper provided reachability and selection review. Primary checked
the whole-source references and nested fixture constructor chains. Luna identified
54 additional roots beyond the preliminary 234; all exact names and source hashes
were independently checked before the baseline selection was frozen. Luna is now
drafting the bounded conversion in an ignored overlay; primary owns baseline
acceptance, integration, boundary review and qualification. No usage savings are
inferred from this outcome.

Local drafts and evidence: `build/port-server-fixture-bridges/`;
prior selection audit: `build/port-after-modifier-identities/`.

## Housekeeping

Before conversion, reclaimed 559,865,856 allocated bytes from 1,654 verified
original-asset duplicates in the completed modifier scenario and 340,205,568 bytes
from five obsolete root/legacy cache archives. Host-use, hash and stat checks
passed. Originals, current binaries/caches, saves and restoration manifests remain;
see the checkpoint for recovery commands. Cleanup scripts are consumed.

## Reviewed conversion

Primary verified the frozen 55-path draft against its source hashes. All retained
production owner functions match the original AST representation. The 80 copied
fixture adapters preserve original bodies and conversions; durability uses its
existing Go owner directly. All 81 header deletions remove selected declarations.

Before compilation, primary added three missing fixture imports and removed an
unused production import. Twelve numeric-only fixtures now use equivalent Go
widths and drop cgo imports; the identity table has only its 19 used keys. World
and floating-point observers retain their original C bodies with independent
standard-header declarations. Other live fixture C interfaces remain.

No production algorithm, assertion or captured expectation changed. The expected
production cgo decrease is 18; the measured post-qualification inventory governs
the final reported count.

## Qualified results

All 288 default/highres and 287 server roots passed, with exact expected name sets
and no skips. Safe/static checks, three production/ABI builds, the exact known
asset-suite results and a fresh headless save/load/resume scenario passed. All
1,654 original asset hashes remain unchanged. Accepted phases have identical
source fingerprints, and all 55 changed source paths match primary review.

No shared production algorithm or dispatch API changed, so the affected selection
in all profiles plus production gates qualifies this batch. The prior full default
corpus (2,457 passes plus one known diagnostic skip) belongs to modifier commit
99b65896; it supplied source-identical original-path baseline evidence here, not
a claim of a post-conversion full-corpus run. Existing assertions and captures
remain unchanged.

Selected production cgo files: **193→175**; legacy C exports: **637→556**.
Production and test-reference standalone C remain **0 lines**. The 157 tracked
headers contain 3,313 physical lines. External native-library bindings are
unchanged; this does not establish a whole-build cgo-free executable.

[Qualification](server-fixture-bridges-qualification.json),
[updated inventory](server-fixture-bridges-inventory-after.json).
Original baseline: 9a833ff4.

The Luna draft was useful for the bounded wrapper/header/caller migration. Primary
review caught three missing adapter imports and one unused owner import before
compilation, and removed avoidable test cgo types/unused identity slots. These
corrections reinforce explicit compiler-visible import and consumer review; no
subscription savings were measured.
