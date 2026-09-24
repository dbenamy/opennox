# Remaining unused C export bridges

Status: original baseline qualified and frozen; conversion not yet installed.
Complete-corpus evidence is reused by exact source identity from `4a0ab0dc`;
production evidence is reused from `9adcf3d3` with the sole qualified test-fixture
difference explicitly recorded. See remaining-unused-exports-baseline.json.

Retire 378 unused C export wrappers from 75 mixed Go files and exactly 333 simple
prototypes from 44 headers. Preserve the underlying Go implementations, globals,
types, private functions, tests, registries and stable callback IDs. Two touched
files no longer need a C import or preamble. No complete source file is deleted.

## Reachability and review

Primary scanned all tracked source text, including C preambles, with an exact
identifier dictionary. The accepted cohort has no references beyond its own
export/definition and simple header declarations. Supplemental scans of tracked
source inputs and repository build/configuration files found no additional matches.
The executable build surface remains the supported interface; no external shared
library API is being maintained by these unused wrappers.

Prototype classification accepts only header lines with a real return-type prefix.
An initial filter incorrectly accepted a `return function(...)` call; another
accepted a whitespace-prefixed sample call. Both were corrected before selecting
the final cohort. Preserve `sub_42CD90` and the three exports with additional C
preamble prototypes (`nox_xxx_playDialogFile_44D900`,
`nox_xxx_unitsNewAddToList_4DAC00`, `nox_xxx_spellCastByPlayer_4FEEF0`) for separate
review. Their exclusion is conservative, not proof they are live.

Luna generated bounded function/prototype deletion spans and unused imports.
Primary reconstructed the changes from original bytes, independently verified
every function end with Go's parser, and compared every remaining non-import
declaration through the AST. Header output equals the original minus the exact
333 recorded lines. Removed imports remain represented elsewhere in their package,
so package initialization is not intentionally dropped. Primary additionally
removed the now-unused C imports/preambles from `game_statistics.go` and
`map_growth_exports.go`, preserving all their Go declarations.

## Qualification plan

Reuse the freshly qualified complete root corpus only with exact source identity.
The production baseline is the qualified 65-export cleanup; its sole subsequent
source difference is the qualified player-reset test fixture, not engine code.
All existing fixtures and frozen expectations remain unchanged in this removal.

After conversion, require the complete porttest corpus in default/server/highres,
with exact discovered/executed/completed root-name sets and only the documented
map-population diagnostic skip. Then safe build/static checks, three fresh
production builds, ABI checks requiring all 378 symbols absent, the exact known
suite failure/package outcomes, and headless character creation/save/load/resume.
Remaining required ABI names are preserved; 66 formerly required names in this
proven-unused cohort move explicitly into the retired set.

The expected dependency reduction is two selected cgo files and 378 selected
legacy C export bridges. Measure the actual counts after qualification. The 79
embedded callback bodies and external native backends are outside this batch.

## Delegation

The bounded deletion draft was useful and mechanically reviewable. Primary kept
inventory algorithm design, cohort acceptance, import cleanup, integration and
qualification. No cost-saving claim is made from a successful draft alone.

Local review artifacts: `build/port-next-exports/`.
