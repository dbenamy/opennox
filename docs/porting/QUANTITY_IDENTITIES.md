# Quantity, shop and trade callback conversion

## Scope and preserved boundaries

Retire ten C export bridges: six deferred amount callbacks (buy, sell, repair,
sell cancellation, repair cancellation and inventory drop), the shop-active query,
shop and trade tooltips, and the outer amount-dialog adapter. Retire the private C-typed
trade-add fixture adapter. Move every producer, stored-identity comparison and
fixture normalization with those boundaries. The two export-only files can then
be deleted. Remove now-unused cgo preambles in their connected owners.

Preserve the dialog's callback slots and five argument words, including unsigned
clamping before dispatch. Keep its temporary point allocation and deferred release
in their original positions: release happens after the callback and dialog toggle.
Buy narrows code to 16 bits; sell narrows code/type; repair narrows code; inventory
drop interprets point coordinates and count as signed 32-bit values. Ignore the
same return words. Unknown foreign identities retain the exact five-word void
callback fallback. Shop closure recognizes only its three accepting identities.

The two tooltip adapters use the existing GUI registry and decode packed unsigned
32-bit points. Preserve window layout, algorithm bodies, packet bytes, fixture
operation numbers, assertions and frozen state/pixel captures.

## Original baseline and coverage

[Baseline](quantity-identities-baseline.json) reuses the just-qualified inventory
revision `785c2d54`: all production/test and supplemental fingerprints match,
with identical environment and exact selected names. No source was added or
changed for this baseline. All 394 client / 391 server roots passed in three
profiles. Production, ABI, safe/static, known-suite, two headless save/load runs
and asset hashes are also recorded in that preceding qualification. This is
reused evidence, not another execution.

The selected contracts exercise all six callbacks through actual dialog
acceptance/cancellation, close-time ownership, stack restoration and packet bytes.
Other cases cover nil and foreign callbacks, five argument words, unsigned limits,
re-entry through Go with forced GC, reopening, allocation failures and drawable
lifetime. Installed trade tooltips have explicit owner-effect assertions; shop
coverage combines constructor identity snapshots, independent grid-hover effects
and the existing shared-registry dispatch contracts. No new mirror tests are
needed for this boundary conversion. Primary also reviews adapter widths and
argument/defer ordering before qualification.

Converted acceptance requires the full affected selection in three profiles, safe/static,
three production/ABI builds, exact known-suite comparison, two headless save/load
scenarios and original asset hashes. The converted gates below are now accepted.

## Delegation and review

One Luna helper drafted a bounded overlay under `build/port-quantity-identities/`;
primary owns scope, caller review, installation, qualification and acceptance.
Primary corrected two scout claims: the export files contain no unrelated entries,
and the temporary point is freed after the dialog toggle. A whole-file dependency
review also identified four connected cgo imports that can retire with those two
files. Dependency discovery and qualification confirmed those removals.

## Qualified conversion

All 394 client / 391 server roots pass with exact original names and no
failures/skips. Safe/static, three production/ABI builds, exact known-suite
comparison, two headless save/load scenarios and all 1,654 original asset hashes
pass on the same reviewed source. See [qualification](quantity-identities-qualification.json).

Legacy exports fall 220→210; selected production cgo files fall 136/137→130/131.
Headers remain 157 files / 2,970 physical lines. Embedded production C bodies
remain 77; standalone production and test-reference C remain zero.

Luna supplied the bounded 18-path overlay. Primary verified all original/draft
hashes, whole-source consumers, every changed function and all fixture callback
table substitutions before installation. Pre-compile review removed the unused
ccall import and the now-empty fixture preambles, and corrected stale adapter
comments. Existing assertions and state/pixel captures were unchanged. No measured
subscription savings are claimed.
