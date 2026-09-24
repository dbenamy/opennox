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
[Baseline evidence](server-fixture-bridges-baseline.json). No conversion is installed.

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
