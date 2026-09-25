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

Two small adjacent cleanups are proposed: replace the remaining C integer cast
in the inventory-place owner with an equivalent `int32` cast, and replace the
fixture's sole call to non-exported `sub_465DE0` with its existing Go owner. The
latter wrapper is not an additional C export. Together with three include-only
owner cleanups, these may remove five production cgo imports; measure the result.

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
Source, binary and environment identity were verified for both runs. No conversion
is installed. The indirect-caller review found no additional root beyond this union.

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
