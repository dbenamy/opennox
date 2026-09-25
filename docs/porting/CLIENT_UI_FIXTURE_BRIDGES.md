# Client UI fixture C bridges

## Scope

Retire 71 fixture-only C export bridges across 19 UI owner files, covering window
helpers, inventory, shop/trade, spellbook, quickbar, binding controls, tooltips,
meters and drawing helpers. Preserve owner behavior, ABI widths, raw pointer
values, snapshot identities and all original assertions/captures. External
SDL2/OpenGL/OpenAL backends remain unchanged.

Primary independently matched all 390 reference occurrences across 3,063 current
Go/C/header/assembly files. The larger helper scan includes other tracked source
files as well. References are definitions, prototypes and fixture calls/getters;
two extra production preambles contain declarations only. Repeated symbol names
on one line count as separate occurrences. No production call or registration
uses were found for the 71 exported symbols.

Two small adjacent cleanups are included: replace the remaining C integer cast
in the inventory-place owner with an equivalent `int32` cast, and replace the
fixture's sole call to non-exported `sub_465DE0` with its existing Go owner. The
latter wrapper is not an additional C export. Together with three include-only
owner cleanups, the final profile inventory confirms five production cgo imports removed.

[Reachability](client-ui-fixture-bridges-reachability.json),
[qualification manifest](client-ui-fixture-bridges-batch.json).
Local work: `build/port-client-ui-fixture-bridges/`.

## Baseline and test selection

The original source is qualified server-fixture commit `a798ad1c`. Existing profile
binaries have matching source and supplemental-input fingerprints; original-path
baseline tests run directly through those verified binaries, avoiding a rebuild.

Initial 261 default/highres and 260 server tests passed with exact name sets and
no skips. `TestClientInventoryWindowWorldSelection` is explicitly excluded from
server builds. An additional 122 shared-fixture caller tests passed in all profiles, giving
383 default/highres and 382 server roots with exact name sets and no skips.
Source, binary and environment identity were verified for both runs. The conversion has completed qualification. The
indirect-caller review found no additional root beyond this union.

[Accepted original runs](client-ui-fixture-bridges-baseline.json),
[exact selection](client-ui-fixture-bridges-tests.txt),
[shared-caller audit](client-ui-fixture-bridges-selection-review.json).

After conversion, run the full combined affected selection in all three profiles,
a full default corpus, safe/static checks, three production/ABI builds, exact
known-suite comparison, a fresh headless save/load/resume scenario and original
asset integrity. The full corpus is useful after two consecutive fixture-bridge
batches because these UI constructors serve many other test families.

## Boundary review

Primary read all 71 original wrapper bodies. Preserve window nil results (`-2`),
output writes and aliasing, old flag return values, self-owner fallback, raw image
handles, and inclusive loops at signed-integer limits. Book notification/kind
pointers represent pointer values rather than strings. Inventory point helpers
have different signed/unsigned input words; preserve their 386 conversions and
relative-coordinate outputs. Trade arguments retain low 16 bits. Drawing helpers
retain color words, size and geometry behavior.

Fixture callback maps supply snapshot normalization. Preserve their names,
indices, existing nil slots and distinct used identities. Trace every getter
consumer before replacing function addresses; identity keys must never enter a
raw callback execution path. Keep independent C observers and other live fixture
C interfaces.

## Delegation and review

One GPT-6 Luna helper produced the reachability and initial selection audit.
Primary verified references, definitions and constraints, then found additional
shared-fixture callers beyond the proposed family selection. A conservative
free-function scan identified 122 additional roots; the helper checked all 321 caller-edge references, which primary independently
matched to current source lines and hashes. The 122 comprise 82 identity-getter
roots, two direct candidate-dispatch API roots and 38 conservative shared-owner
roots. Primary found no function-valued uses of the affected public fixture APIs;
the 14 affected root-owner methods on 11 types use constructors covered by the
selected families and added caller roots. This is a bounded source audit, not a
claim of a complete dynamic call graph. Full-corpus qualification remains required. Keep this audit separate from
runtime acceptance; a passing family or matching count does not prove coverage.

## Housekeeping

Reclaimed 559,919,104 allocated bytes from 1,654 verified original-asset duplicates
in the completed server-fixture scenario and 425,369,600 bytes from nine obsolete
root/legacy cache archives. Host-use, hash and stat checks passed. Originals,
current caches/binaries, saves and restoration records remain. Recovery commands
are in PORTING_STATE.md; cleanup scripts are consumed.

## Implementation and current qualification

Primary installed 46 changed/new source files after syntax formatting, retained
production-body review and exact getter mapping checks. The only retained
production function edit is the equivalent geometry-result cast; the private
window-level wrapper and two stale preamble declarations are removed too.

Thirty-five used fixture adapters preserve retired export boundaries; 30 match
original bodies modulo fixed-width scalar types, while five preserve equivalent
point/window pointer representations and conversions. Fifteen window adapters
retain original bodies with fixed-width scalar types, including sequential size
output writes. Four drawing/tooltip routes call existing Go owners directly.
Twenty-eight static byte identities preserve all six getter mappings, names,
indices and nil slots. Existing root assertions and captures are unchanged.

Luna supplied production removals, most fixture routes and all 28 identity keys.
Primary finished the window/render/tooltip routes locally as draft completion
started delaying integration, then stopped the helper and froze its partial
overlay. Primary removed 36 unused copied adapters (including duplicates of the
locally completed routes), completed imports and adjacent cleanups, and removed
two stale production declarations. An unused draft window-size adapter coalesced
two output writes; the installed adapter preserves their original order. These
were pre-build corrections, not changes to expectations after a test failure.
The bounded reachability audit was useful; next implementation handoffs should
list actual adapter calls and reserve subtle window/output behavior for primary
work. No measured subscription saving is claimed.

All 383 default/highres and 382 server focused roots passed, with exact baseline
name sets and no skips. The full default corpus passed 2,457 roots with the one
established diagnostic skip (`TestMapPopulationPrerequisiteProbe`), matching the
independently recorded 2,458-name set. Safe/static, all three production/ABI
builds, exact known asset-suite comparison and fresh headless save/load/resume
passed. All 1,654 original asset hashes are unchanged. Accepted phases share
identical source fingerprints and all 46 source paths match primary review.

All three focused builds passed on the first attempt. No failure-driven source
or expectation changes were required. Selected cgo files: **175→170**; selected
legacy exports: **556→485**. One additional fixture cgo import was removed.
Standalone production/test-reference C remains **0 lines**. External native
bindings are unchanged.

[Qualification](client-ui-fixture-bridges-qualification.json),
[updated inventory](client-ui-fixture-bridges-inventory-after.json).
Original baseline commit: `49cf5daf`.

The 157 tracked headers contain 3,242 physical lines.
