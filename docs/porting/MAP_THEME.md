# Theme parser

Active scope: 33 routines / 1,888 C section lines, 51E260..520D50 in GAME4_2.c.
Hallway native conversion **f1b92d3f** is committed/pushed and fully qualified.
Only mapGenReadTheme and sub_520D50 have external callers; 31 helpers can retire.
Current physical C remains **104,307 / 148 files / zero reference C**.

## Fixture bootstrap

The initial tagged fixture uses real Go-backed legacy file handles and shared
record snapshots. A test-only calloc/free observer records allocations and
release before free, active only during the parser invocation; a test-only clock
makes the default seed deterministic. Other tests pass through the real services.
The fixture initializes/restores actual startup string tables from blobdata and
tracks the token buffer, line count, template pointer and input position/error.
The adapters must be absent from production builds and must preserve older
mandatory capture hashes.

Current smoke: TestMapThemeTokensProbe in build/port-map-theme/tokens-original.log.
The isolated TestMapThemeAlgorithmProbe exercises a normal numeric setting before
repairing the suspected scalar token buffer. The equipment parser also traverses
four scalar counters as an array; probe this before any contiguous-record repair.
No production theme C changes or locked theme hashes yet. Staged allocator,
callbacks, dispatcher and blob-table getter have been applied; do not reapply.

## Plan

Cover whitespace/comments/EOF and counters, nested conditions and operators,
player classes/counts, experience level and percentage draws, all algorithm
settings/defaults, equipment/spell sets/modifiers, decoration types and copies,
occurrence/frequency/room constraints, doors, prefabs and cleanup. Use synthetic
full theme files through the actual file layer. The supplied image has no .thm
files or AreaMap.lib; do not claim actual theme/prefab content integration.
Ordinary asset-backed gameplay remains available.

Repeat complete corrected-C captures in default/server/highres, lock mandatory
hashes and commit/push before conversion. Then qualify the native implementation
with complete captures, accumulated variants, production binaries/symbol audit,
known full-suite comparison and fresh unchanged headless gameplay. Count C,
update recovery docs, commit/push, summarize and continue.

## Local evidence storage

Completed population/hallway gameplay directories now keep changed/new files,
logs and recordings plus deduplicated-assets.json. Identical asset copies were
verified by SHA-256 against build/assets/extracted/drive_c/Nox before removal.
Restore an entire data tree with the command in its manifest:
`python3 build/baseline/deduplicate-run-assets.py --restore map-hallways-port`
(or map-population-port). The script verifies source and restored hashes; the
original supplied archive is untouched. This reclaimed about 1.1 GB while keeping
all gameplay evidence reconstructable. New gameplay runs still make fresh copies.

## Contiguous-record prerequisite — validation running

The token fixture smoke passes (0.086s). Both original focused probes abort:
algorithm-original.log for `midHallLength 12 END`, and modifier-original.log for
one named value in each of the four modifier slots. The real C parser was used;
no production parser behavior was replaced by the fixture.

Restore a 60-byte local value buffer in genReadAlgData_51EBB0 and an explicit
four-element count array in sub_51F230. The former matches the original stack
extent; the latter makes the existing indexed slot traversal defined. This is a
deliberate prerequisite repair, not preservation of the original aborts.
Corrected probes and existing hallway/population/painting/room checks are running
in prerequisite-default.log. Qualify all variants and commit/push the prerequisite
before broad theme captures and conversion. No theme baseline is locked yet.

Corrected default map checks pass in **41.556s**. The expanded **72-case** theme
prerequisite smoke passes in **0.823s**: five numeric token lengths up to 59
characters, all 16 slot-presence masks, one/three strings per slot, both input
orders, and the original three probes. Assertions include exact slot counts,
allocation extents and every stored string byte, plus guard/control state.
Server/highres map checks are running under prerequisite-variants.py; inspect
prerequisite-variants.log and prerequisite-variants.json. Working C is **104,300 /
148 files / zero reference C**, −7 lines from consolidating counter declarations
and initialization. No conversion has occurred in this batch.

All prerequisite variants pass: server **124.688s**, highres **53.839s** wall time
including builds, alongside the existing mandatory map captures. Default broad
checks pass in 41.556s; the expanded 72-case prerequisite smoke passes in 0.823s.
Record the −7-line C prerequisite checkpoint and commit/push before extending
the full theme baseline. No previous expected hashes changed.

## Inherited modifier removal prerequisite — qualification running

Buffer/counter prerequisite **ff4de53d** is committed/pushed. Full synthetic
file parsing and cleanup now pass, with a fixed default seed and actual file
closure. New stream/player/settings fixtures cover all byte values, 255-byte
tokens, EOF/comment state, nested conditionals, actual sparse player iteration,
all operators, algorithm keys/numeric boundaries and decoration frequency totals.

An original template-removal probe fails: removing alpha from alpha/beta/gamma
does not preserve beta/gamma. The C shift loop advanced its source before copying,
skipping the next value and reading beyond the populated scratch entries.
Move the increment after strcpy. Evidence: template-removal-original.log.
The 72-case regression spans weapon/armor, all four slots, first/middle/last and
case-insensitive removal, missing/repeated values, additions/deduplication and
cleanup. This is a deliberate correction, not compatibility with that failure.

All current checks pass in **50.195s**. The independent repeat passes in
**52.482s** with **2,146 cases / eight complete captures** identical byte-for-byte;
existing hallway/population/painting/room hashes also remain unchanged. Server
and highres comparisons are running via qualify-removal.py (qualify-removal.log).
Working C remains **104,300 / 148 files / zero reference**; the copy-order fix
changes no line count. Lock the eight verified hashes and commit/push this
prerequisite checkpoint before expanding the remaining parser sections.

Current capture labels and hashes: removal-captures.json; main raw copies:
c-removal-*.json. Repeated variant copies are losslessly compressed after exact
comparison, with manifest removal-capture-archive.json. No native theme routines
yet. The full-file and player draft files were applied and the fixture now tracks
new/closed file handles and verifies player records remain byte-for-byte unchanged.

All eight captures match in default repeat/server/highres: **52.482s / 129.368s /
60.201s** wall time, including builds and existing map checks. The eight hashes
are now mandatory in map_theme_baseline_porttest_test.go. The mandatory-hash smoke passes in **8.485s** (locked-streams.log). New C-only corpus labels may
be captured with OPENNOX_MAP_THEME_EXTEND_C=1; this never bypasses an existing
hash. Remove that extension allowance before native conversion, once all labels
are locked. No native theme algorithms are written yet.
