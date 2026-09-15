# Porting decisions for later review

## Working policy (user instruction, 2026-09-12)

When reasonably confident in the right answer and reversal would not require
substantial effort, make the decision and record it here or in the batch doc
for later review. Continue without asking for confirmation. Record behavior
changes explicitly, with evidence and validation; do not label them exact
compatibility. Ask when confidence is insufficient or reversal would be costly.

## Locked-door message padding — approved, review with Go conversion

Initialize unused bytes in the fixed 52-byte notification. The client interprets
only the header, NUL-terminated localization key and selector. Original C left
padding uninitialized; a zero-initialized C array now establishes deterministic
output for the upcoming Go conversion. No meaningful fields or string acceptance
rules change. Reversal is one initializer, though preserving indeterminate
padding is not recommended. The user explicitly approved this correction.
Evidence and reproducible tests: [PLAYER_CONTROLS.md](PLAYER_CONTROLS.md).

## Default-equipment modifier word — review with controls conversion

Initialize the entire 20-byte C modifier array in 4EF7D0. Only four descriptor
words were assigned, but the attribute helper copies a fifth word into items.
Chosen under the standing policy: use deterministic zero for that uninitialized
word, matching the future Go array default. Reversal is one initializer. Retain
all descriptor selection, ordering, callbacks and defined values. Evidence and
validation are recorded in [PLAYER_CONTROLS.md](PLAYER_CONTROLS.md).

## Spell phoneme class offset — review with spell lifecycle conversion

Correct the server phoneme helper's class-byte access before locking its C
baseline. `getObjectFromNetCode` returns `nox_object_t*`; adding 8 before the cast
advanced eight whole objects (6,176 bytes), instead of reading byte offset 8.
Nonplayer phoneme cases exposed an out-of-bounds read and crash. Cast to a byte
pointer before adding the offset, matching the client branch's class-field check.
The new corpus checks male/female and nonplayer sounds in both server and client
paths. This is a deliberate bug correction, not exact preservation of the invalid
read. Chosen under the standing policy; reversal is one expression.

## Buff power accessor — review with spell lifecycle conversion

Change server.Object.EnchantPower from BuffsDur to BuffsPower. The C accessor
uses the byte power array; the existing Go method incorrectly returned the
16-bit duration. The locked C corpus independently distinguishes duration 1,234
from power 7 and covers signed C byte returns. Native callers now receive power.
Chosen under the standing policy; reversal is one field access. This intentional
correction is separate from exact C/native lifecycle compatibility.

## Creature-tag caster read — review with sustained spells conversion

Move the caster-data read below the existing nil-caster guard in 530160 before
locking the C baseline. The source previously dereferenced the caster before
testing it. A dedicated nil-caster test already passes with the current compiler;
this is not evidence of an observed runtime crash. The change makes the intended
rejection defined in the source and independent of optimization. Positive input
behavior is unchanged. Chosen under the standing policy; reversal moves one line.
Evidence: `build/port-sustained-spells/c-tag-nil-before.log` and the committed
`TestSustainedSpellsTagNilCaster` regression case when this baseline is locked.

## Plasma direction predicate — preserve now, review as a gameplay change

Preserve the original expression in 531920: OR-ing the direction mask with 0xC
makes that part of the predicate always true. Enemy and interaction checks still
apply. The native helper evaluates the direction owner and preserves selection
behavior; all locked spatial captures remain unchanged. Any intended restriction
to targets in front should be a separate gameplay fix with explicit tests, not an
incidental change during conversion. This choice is reversible in one predicate
and follows the standing policy of documenting such decisions without pausing.

## Map-generation random range — review the compatibility correction

Constrain `nox_platform_rand` to `platform.RandInt() & 0x7fff` before locking the
room-generation C baseline. Both remaining production callers scale the result
as a 15-bit CRT random value; the real Go platform returns a wider integer.
A seeded original-C request for a float in [-5, 8] returned 575182.218727404.
The independent random-range corpus failed, and its rejection-sampling cases
took 52.35 seconds. The full Go platform API is unchanged; only its C compatibility
export is narrowed. Low 15 bits also avoid architecture-dependent int width.

This deliberately changes generated layouts for a given seed, fixing out-of-range
placements and heavily biased/sluggish selection. It is not exact preservation
of the broken adapter. Chosen under the standing policy: the appropriate range
is explicit in both consumers and reversal is one expression. Retain this as a
seed-compatibility decision to review later. Evidence: build/port-map-rooms/
rng-before-fix.json, c-smoke-map-rooms-smoke-59.json and c-boundaries.log.

## Map-painting stack records — review the compatibility correction

Replace separate decompiler locals consumed as contiguous coordinates, dimensions
or runtime tile patterns with explicit arrays (and one mixed int/float union).
Ten painting functions have this layout defect. Also initialize the border
entry's eight-word pattern before assigning its tile and border fields; its
anchor-mode field was otherwise uninitialized. This is a prerequisite C repair,
not an attempt to preserve compiler-dependent behavior during the Go conversion.

Repeated original-C processes with identical inputs disagreed in 26 capture
groups. A positive 1x1 rectangle alternated between painting and no change.
Optimized i386 assembly for sub_5245A0 retains only the X local at stack offset
40, places the stack canary at 44, and removes the intended Y assignments.
The callee reads two floats through the X address, so Y comes from unrelated
stack data. Independent corner-mask checks also fail before the repair.
Evidence: build/port-map-painting/c-expanded{,-b} captures/logs and game4_2.s.

Chosen under the standing decision policy: the intended record layouts are
explicit in the consumers and decompiler offsets; reverting is inexpensive.
Generated layouts can deliberately differ from the broken C implementation.
Repeated corrected-C captures and independent painting/corner contracts must
pass before the new C baseline is locked. Review this with the prior map RNG
compatibility fix if old seed/layout reproduction becomes a requirement.

## Native map-object admission — review the defined failure paths

Reject a nil object-selection name before string comparison, and use the existing
server NewObjectByTypeInd guard when a selected type index is stale. The C name
routine had a nil check after its unguarded strcmpi call; its placement routine
checked the allocation result, but the C factory adapter could dereference a
missing type first. Native callers now receive the intended zero/nil result.
Defined-input behavior remains covered by unchanged C capture hashes. Two
additional native admission contracts cover these formerly unsafe inputs.
This reversible choice follows the standing decision policy; it is separate from
byte-for-byte preservation of the valid-input painting corpus.


## Population stack records — review generated-layout compatibility

Original C execution confirms a spell-name crash, a point-output stack abort,
and PlayerStart placement at a clamped map corner instead of the room center.
Use explicit buffers of the sizes required by their consumers in spell-name
lookup, population item placement and the population finale. Reject spell names
that do not fit the existing 60-byte input record. The formatted output needs
66 bytes including the `SPELL_` prefix and terminator. These are prerequisites
to the corrected-C baseline, not differences hidden inside conversion hashes.
The source fix is reversible; retain the corrected behavior unless reproduction
of old broken generation becomes an explicit requirement. MAP_POPULATION.md
records execution evidence and qualification status. After repairing spell-name lookup, the isolated invalid-book probe aborts with
`free(): invalid size`. Route disposal through the existing engine object-free
service; verify both the zero result and restored live-object count. This keeps
object-pool ownership intact. It does not redesign the engine disposal service.


## Prefab candidate arrays — same stack-record prerequisite

526550 has two six-entry candidate records represented as individual scalars
followed by unrelated five-entry arrays. An original-C two-room probe crashes;
use explicit six-entry arrays before locking the population baseline. This is
the same reversible storage repair as the point/name fixes. Candidate ordering,
nearest-six selection and link topology need independent checks as well as
complete corrected-C captures. No other candidate-selection policy is changed.

## Population item attributes — initialize the complete copied record

5221A0 initializes four words of a five-word local item-attribute buffer; the
existing attribute setter copies all five words into the object. Independent
complete C captures agree in 26 groups, but 356 enchanted-item cases differ only
in the fifth word. Initialize the complete 20-byte buffer before constructing
modifiers. The explicit trailing-word regression fails on the original C. This
is a prerequisite correction, not behavior to reproduce in Go. Review later if
that fifth field should acquire an explicit nonzero gameplay meaning; arbitrary
stack contents are not a defined value. Evidence: modifier-repeat-diagnosis.json,
modifier-initialization-original.log and modifier-fixed.log under
build/port-map-population. No physical C line-count change from this initializer.

## Hallway second-corridor storage — 2026-09-14, review later

Use the blob corridor array consistently in the ten second-corridor connection
reads across 54B810/54BB20/54BD90/54BF20. Original bent-route regression passes a
null room because the separate named global is never assigned by construction.
This cheap reversible correction follows the standing policy; retaining that
failure would prevent ordinary multi-segment routes. See [MAP_HALLWAYS.md](MAP_HALLWAYS.md)
for original evidence, regression and validation status. No unrelated named/blob
globals are merged or treated as aliases.

## Short-gap hallway shapes — 2026-09-14, review separately

Preserve existing signed corridor-length arithmetic in the hallway conversion.
The expanded gap 1–3 matrix records 491 admitted zero/negative-length corridors;
these are defined record values, with no invalid memory access observed. Rejecting
them would alter fallback and route selection beyond the pointer-storage repair.
Keep the exact C captures and flag the gameplay behavior for a later correction;
see [MAP_HALLWAYS.md](MAP_HALLWAYS.md). This is a compatibility choice, not an
endorsement of the resulting shapes.

## Theme value buffer and modifier counters — 2026-09-14, review later

Restore contiguous local storage before qualifying the theme parser: a 60-byte
algorithm value buffer and four modifier counters in one array. Both original
normal-input probes abort, documented in [MAP_THEME.md](MAP_THEME.md). These
small reversible repairs follow the standing policy; preserving the aborts
would prevent ordinary settings and equipment parsing. Broader parsing and
route behavior remain outside this prerequisite change.

## Inherited modifier deletion — 2026-09-14, review later

Copy the next inherited modifier into the vacated slot before advancing the
source pointer. The original loop skips that value and reads beyond populated
scratch entries; the alpha/beta/gamma removal regression fails. The corrected
72-case matrix checks weapon/armor slots and cleanup. This small copy-order
repair follows the standing policy; see [MAP_THEME.md](MAP_THEME.md).

### Theme parser shared ownership and token state (2026-09-14; review later)

Preserve the C allocator ownership of theme records, shallow decoration-copy
children, and current cleanup scope while porting the 33 parser routines. The
remaining generator still consumes these records; changing ownership at the same
time would require a larger caller/lifetime audit. Keep standard libc atoi/atof
and time semantics behind small production Go helpers. No old parser algorithms
remain for testing. Complete repeated C captures cover record contents/disposal,
tokens, file state and random draws. The edging table walk continues comparing
against the shared token after each nested read, matching the original source.

### Growth direction array and merge configuration (2026-09-14; review later)

Restore a contiguous four-element direction array in nox_xxx_mapGenFillRoom_4D53B0.
The original stronger probe observed an invalid hallway kind; checking only the
return value had missed it. Preserve normalized record bytes before disposal in
the growth fixture so rejected candidates are verified too. All 16 blocked masks
across eight seeds now exercise only the intended remaining directions.

Point four growth merge-rate reads at the same configuration blob that the theme
parser writes. With config mergeRate=100, the original probe produced zero merges
when an unrelated named word was zero and one merge when it was 100. This is a
deliberate storage correction, consistent with the earlier hallway prerequisite.
Four directions, 0/50/100 rates, eight seeds and independently varied named values
verify that actual configuration controls the result and random draws.

### Door scans read integer grid coordinates (2026-09-14; review later)

Restore integer GridX/GridY loads in the two directional door scans. The source
used a float pointer to a shared room record and converted the integer words as
floats; an otherwise equivalent translated-room probe produced no doors/waypoints.
The expanded contract checks both object counts and projected coordinate shifts
for positive/negative translations in all four directions. Qualify this small C
correction with the full growth baseline before native conversion.

## Cache invariant lookups in map capture fixtures (2026-09-14)

Larger generator integration cases spent excessive time normalizing released grid
records. Cache the two fixed C transfer-function addresses once per fixture case,
and index stable wall-address ranges by page. Keep the original range checks,
interior-byte offsets, pool order and canonical IDs; rebuild the wall index if the
fixture pool grows. These are porttest-only changes and do not alter production.

The first cached run preserves all growth/hallway hashes and completes 128 normal
map integrations. All prior map hashes pass in default/server/highres (66.162 / 144.538 / 74.101s
wall). Matching growth test durations sum to 24.96s before and 17.47s after in
observed runs, not a dedicated benchmark.
Evidence: build/port-map-orchestration/cache-{variants,growth-timing}.json and
cache-*.log. Commit separately from the subsequently exposed wall-list defect.


## 2026-09-14 — correct wall-list traversal before generator orchestration

Decision for later review: `serverWalls.find` must traverse the global `Next20`
link from the global head. Its prior use of the row link dropped live walls on
ordinary deletion and caused a ring-generator integration to hang. This is a
one-line reversible correctness fix, established by a failing direct regression.
Twelve deletion/reuse contracts cover different rows, one row and bucket collisions.
The intentional update to 11 historical painting hashes is field-audited: only
520 global-list links and 31 heads change; every other field remains identical.
See [WALL_LIST.md](WALL_LIST.md) for qualification and evidence.

Keep diagnostic release storage outside the engine's allocator. C-allocated trace
buffers reused disposed engine records and made canonical pointer captures depend
on allocator layout. Go-owned retained buffers plus saved grid-row identities
remove that observer effect. This is test-only; no engine allocation ownership
changes. Repeated complete captures and eight allocation layouts check stability.
The separate shallow theme-cleanup ownership issue remains deferred for review.


## 2026-09-14 — unflagged backdrop decoration weights

Real backdrop integration exposed an inherited selector bug: a newly created
backdrop has ThemeFlags zero, while the selection default used the room address as
its random weight. Depending on allocation address it returned nil (followed by a
painting failure) or performed an enormous address-dependent draw/retry sequence.
The earlier room conversion preserved that C behavior; this is a separate repair.

Decision for later review: for ThemeFlags zero, sum weights of matching unrestricted
decorations. Keep the existing six phase-specific weight paths unchanged. This is
a bounded reversible change with direct selection tests and real backdrop maps.
Empty/ineligible lists return nil. All22,658 map cases/281 groups and contracts
pass default/server/highres (74.179s/148.546s/82.656s), all historical hashes
unchanged. No orchestration C routines are converted. The initial real-C failure is c-expanded.log under
build/port-map-orchestration, and the corrected rerun is c-backdrop-fixed.log.

## 2026-09-14 — initialize unused journal message bytes

Decision for later review: initialize the three 68-byte journal add/remove/update
buffers before filling their defined fields. Short names previously left unused
bytes dependent on prior stack contents. The original-C regression reproduces
this for all three operations; corrected tests pass. A recursive before/after
audit permits differences only after the name terminator through byte65 for
add/update or byte67 for remove. Nine captured message snapshots change in254
unused bytes; every defined field, routing value and other captured state is
identical. No message masking or C reference implementation was added.

Local evidence: `build/port-gameplay-reports/journal-padding-audit.json` and its
reproducible audit script, comparing `journal-isolated-original.json` against
`journal-corrected.json`. Earlier `journal-original.json` predates a fixture
record-list isolation fix and is not the comparison oracle. Reporting baseline
qualification remains in progress. Physical production C is unchanged at101,128
lines in148files, with zero reference C.

## 2026-09-14 — gameplay-reporting baseline qualification

Use a full accumulated default port run plus focused reporting/player-control
server and highres runs before conversion, then full accumulated runs on allthree
variants after conversion. All7,809 new C cases/28 captures repeat exactly before
locking; the original73 routines remain the implementation during this baseline.
The default accumulated run passes81,468 cases/1,010 groups and contracts. This
keeps full variant validation of the final change while avoiding a second full
variant matrix before the conversion. Revisit if a baseline result or variant
difference identifies a broader dependency.

The reporting fixture reuses real queues, guarded records, player iteration and
team/member lists. Elimination countdown transitions are recorded at the existing
service boundary; the two empty headless text slots are initialized/restored with
Settings.c:SuddenDeathImminent from the existing Go handler. Network recipient
masks are optional per-case inputs and restored afterward. Reporting actions
check and preserve x87 precision/rounding on a locked OS thread. These fixture
additions leave all prior capture hashes unchanged.

Post-removal ABI audit:45 routines retain C callers;28 entry points can be retired.
No production Go address references require retaining those28 symbols. Production
reporting will use native helpers, with thin Go-backed C bridges for actual C
callers and no retained C algorithms for testing.

## 2026-09-14 — native reporting buffer ownership

Use temporary Go byte buffers for gameplay reporting and pass them directly to
the retained C queue insertion routines. Source review confirms both normal and
coalescing insertion copy input synchronously into owned queue storage; neither
retains the input pointer. This avoids an additional temporary C allocation while
keeping all engine objects and queued records under their existing ownership.
Direct message-list insertion uses the existing native queue. Revisit this choice
if queue ownership changes; current byte/routing/state captures remain exact.

The focused native reporting run matches all28 C hashes. Correction: the original
accumulated native variant runs selected no tests because of a trailing newline
in the saved pattern. Their claimed coverage is superseded by the corrected
text-batch matrix (see GAMEPLAY_REPORTS.md). Keep the45 C bridges required by
actual C callers and remove28 obsolete entry points. The full source audit finds
no retired-name references. Allthree production builds and symbol audits pass; the full suite matches1,553
known failure entries exactly and fresh unchanged headless gameplay passes35.208s.

## 2026-09-14 — reject empty or incomplete port-test selections

A saved pattern ending in `$` plus a newline made allthree original broad native
reporting runs select zero root tests. Go returned success, and the original
runner accepted it. Focused reporting, C baselines, builds, full-suite comparison
and gameplay were unaffected. Correct the historical claims instead of treating
those runs as evidence. The text-batch pattern had already stripped the newline. Its corrected full
default/server/highres matrix passes; all builds/symbol checks pass, the full
suite matches1,553 known failure entries, and unchanged gameplay passes34.664s.

The reusable tools/porting/run_tests.py trims surrounding whitespace, rejects
empty/multiline patterns, discovers root tests with Go itself, rejects zero
selection, then checks `go test -json` events for every discovered root test to
start and finish. Logs remain local. Validate it against an actual test selected
by a newline-terminated pattern and against a pattern selecting nothing. Use this
runner for subsequent accumulated matrices and retain selected/executed counts.

Guard validation passes: an actual newline-terminated selection ran and completed
its one root test; a zero-match pattern was rejected; a controlled Go-output
fixture that listed two tests but ran only one was rejected with the missing test
named. Discovery audits select603 root tests for reporting and614 for reporting
plus text; the latter is exactly the former plus11 text tests. Local evidence is
runner-positive/empty/partial-result.json and selection-audit.json under
build/port-gameplay-text. These checks do not replace the running full matrix.

The final per-variant discovery audit selects614/613/614 root tests. Server omits
only TestFloorEligibility, whose source explicitly requires !server. Actual root
logs for all corrected broad runs contain successful nonempty execution. This
corrected matrix supersedes the invalid earlier reporting matrix; the original
local qualification record is annotated accordingly. The new JSON runner is the
required path for subsequent accumulated matrices, rather than repeating an
already completed valid matrix solely to change its log format.

## 2026-09-14 — scope variant regression work to the changed subsystem

Review later: after object lookup's complete three-variant accumulated matrix,
use one complete accumulated standard matrix plus affected server/highres suites
for the next quest-eligibility batch. Its fourteen routines have no variant
branches; the four outside C entry points are in the common server join path,
and Go callers are the already-captured quest-penalty routines. Include the new
eligibility corpus and existing penalty, inventory, equipment and reward suites
in both affected variants. Keep repeated original-C captures, verified test
selection/execution, all three production builds and gameplay qualification.

This applies the existing subsystem-level testing plan to reduce repetition of
unaffected historical corpora. Run complete variant matrices again for changes
to shared variant-sensitive behavior, fixture ownership affecting those paths,
or an unexplained discrepancy. The current object lookup matrix remains full
in all three variants; this decision applies to following work.

For that tag-independent batch, establish the original-C baseline with repeated
focused captures and the full standard accumulated suite. Run affected variant
suites on the completed native version against those same locked hashes. If a
variant differs, replay the recoverable C baseline in that variant before
attributing the difference to the port. No extra C-variant matrix is required
when there is no variant-dependent implementation or observed discrepancy.

## Client effects renderer reference — 2026-09-14

Preserve current original-C effects output through the current renderer using
exact RGB5551 framebuffer references. Keep the historical sprite color backend
policy unresolved as documented in FAILURE_DIAGNOSIS.md; do not regenerate its
goldens as part of this higher-level conversion. This keeps a recoverable C
reference for effect geometry, state and RNG without coupling the batch to a
separate rendering policy change. Diagnostic PNGs may accompany raw references.
Chosen under the standing reversible-decision policy; review before changing
the color backend. See CLIENT_EFFECTS.md for fixture and qualification scope.

## Missing plasma endpoints — 2026-09-14

Skip plasma drawing when either referenced endpoint is absent, returning the
normal live-drawable result. Its C fallback had continued with local variables
that did not contain endpoint coordinates. The other four ray renderers already
return early in this situation. This correction is needed before establishing
a deterministic plasma reference. A focused regression checks absent source,
absent target and both absent: unchanged drawable/framebuffer/scratch state, no
particle creation and no RNG consumption. This is a deliberate C prerequisite
correction, not exact preservation of the broken branch. Chosen under the user's
standing reversible-decision policy; review alongside CLIENT_EFFECTS.md.

## Coordinate chain-lightning particle endpoint — 2026-09-14

Initialize both target coordinates from the ray's endpoint packet. The C branch
assigned target Y to X and left Y uninitialized. Repeated references differed
in eight unpaused coordinate cases, including particle positions, allocation
counts and RNG consumption; object-bound and paused cases were stable. This is
undefined input to particle generation, so do not freeze a process-dependent
capture or emulate it in Go. A dedicated contract compares coordinate and
object bindings with identical endpoints and requires identical particle
positions/types and RNG consumption. This is a deliberate prerequisite fix,
chosen under the standing reversible-decision policy; review with the plasma
endpoint repair. Earlier ray-drawing hashes are withdrawn.

Preserve two separate existing behaviors for later review: non-plasma ray
callbacks pass flagged static codes unchanged to lookup, while plasma strips
the flag; the forward-difference curve rasterizer stores two coefficients in
named globals distinct from their historical mapped slots. Neither is silently
changed in the effects port. Tests exercise the current storage and lookup
behavior explicitly.

## Effects floating-point and clipping compatibility — 2026-09-14

The first native comparison localized discrepancies to curve/plasma intermediate
rounding and moving-spark left-edge clipping. Read the original compiler output
to preserve the effective boundaries: Hermite X uses wider t/t²/t³ values, Y
uses their float32 spills with wider accumulation; plasma stores normalized Y
as float32 but uses the wider division result in its next product. Moving-spark
X clipping and the bottom edge compare unsigned viewport words, unlike the
random sparkle path. Preserve these current behaviors rather than changing the
locked C references or general renderer policy. Native float-to-int calls use
the already qualified Go helper. Review these quirks before any later geometry
cleanup; the port's compatibility tests remain the constraint.


### Drawable tail allocation retry — 2026-09-14

Guard the two unchecked tail-link allocations before capturing the next C
baseline. Failure preserves the previous anchor for a later retry and adds no
deadline entry. Magic-trail sparks remain independent of tail allocation. This
reversible repair is authorized without another confirmation; revisit the retry
policy if allocation policy changes. See [client updates](CLIENT_UPDATES.md).

Reuse the qualified effects owner for this state-only callback batch. Run the
accumulated standard suite and affected effects/update server/highres tests, all
production builds and unchanged headless gameplay. The full asset suite was
just qualified at the effects subsystem milestone; repeat it on evidence of
shared regression rather than automatically for every callback batch.


### Particle drawing and color startup boundary — 2026-09-14

Include the color initializer with the palette helper it calls. Its only remaining
caller is Go, so both C entry points can retire and maps.go can call Go directly.
Together with the falling-spark helper, three private boundaries retire in the
19-routine /731-line batch. Preserve signed-byte palette behavior; review an
intended gradient separately rather than changing the captured visual baseline.

New fixtures reuse the existing owner without modifying its API or production
code. Repeated C capture plus the related effects/update/drawing standard suite
qualify the baseline. Native qualification broadens to accumulated standard,
affected variants, all builds, full asset-suite comparison and fresh gameplay.
The full suite is warranted here by shared lighting/color startup changes.

## 2026-09-15 — sprite animation baseline and conditional random bound

Correct the conditional animation RNG bound from count to count-1 before locking
its C baseline. The parser allocates count entries and RNG bounds are inclusive;
a sentinel/128-seed contract checks both states and RNG consumption. Use an
isolated, tagged real RenderSprites image owner and retain drawObject as a C
dependency until its broader player/team/clipping paths are owned. Preserve the
existing fade-on-deletion alpha state and conditional data size header for later
review. Remove the unused C vector-array loader instead of porting its already
existing Go replacement again. See CLIENT_SPRITE_ANIMATION.md.

## Object drawing prerequisites (2026-09-15)

Before the next C baseline, failed optional arrow tails now preserve the main
sprite and retry anchor; default glyph opacity is explicitly255 instead of the
drawable address's low byte; generator Random uses the declared last frame index
(count-1) with the inclusive RNG. These are reversible corrections under the
user's standing authorization. Review the intended optional-effect/opacity policy
later; do not reinstate address-dependent pixels to reproduce memory addresses.
Independent rendered-pixel, allocation/lifetime and RNG contracts accompany them.
See [the object drawing checkpoint](CLIENT_OBJECT_DRAWING.md).

## Screen effects baseline (2026-09-15)

Clear encoded rope endpoint IDs before actual static/dynamic drawable lookup;
the producer preserves the static marker while drawable NetCode32 stores the
cleared code. The original-C rendered contract fails for static endpoints and
passes after this reversible prerequisite. Preserve the circle helper's unsigned
trig arithmetic and table-based distance approximation for this conversion;
review unusual circle clipping separately. Own actual particle allocation/list
mutation and real maiden/monster rendering, including gameplay flags. Keep the
adjacent renderer-bounds setter as a shared C dependency. See
[screen effects](CLIENT_SCREEN_EFFECTS.md) for scope and evidence.

## Shared object renderer baseline

Own real player/team, animation, renderer, clipping and sight state before moving
the ten related routines together. Preserve the existing acceptance of short
nonhorizontal occlusion paths and ignored returned bounds; review this behavior
separately. The scanline matrix targets client/highres because server implementations
are intentionally unreachable. Keep shared C team lookup as a dependency; retire
six private renderer entries and retain only four ABI bridges needed by C callers.
See [the shared object renderer report](CLIENT_OBJECT_RENDER.md) for qualification
and the explicit hosted x87 precision check.

## UI rendering and progress-bar batch

Port thirteen connected UI rendering routines with real window/renderer owners.
The pointer-declared rectangle-copy return is scalar 0/1 from noxCopyRect; preserve
those bits and correct its declaration to int when converting. All remaining C
callers ignore the result. Keep particle-border tests within known rasterizer
preconditions (positive radius, contained unclipped images); radius-zero division
and out-of-buffer unclipped rendering are pre-existing issues, not fixes hidden
inside this port. See [UI rendering](CLIENT_UI_RENDER.md) for evidence and limits.

## Slider compatibility

Preserve the original nonzero-minimum position formulas, vertical image draw
no-op, tab-navigation no-ops and exact raw key-state comparisons. Float positions
convert through int64 then narrow to 32 bits, matching the hosted C behavior at
numeric boundaries. Window image/color callbacks are chosen at construction,
independently of later flag mutations. These are reversible compatibility choices
to review later; see [CLIENT_SLIDERS.md](CLIENT_SLIDERS.md).

## Deferred GUI cleanup and radio-owned data

Fix the shared GUI lifecycle before freezing the radio C baseline: queued windows
must retain their cleanup callback and receive it once from FreeDestroyed, even
though ordinary events are rejected after destruction. Preserve deferred timing
and queue order. Reject duplicate destruction using the actual dead/destroyed
state, not a getter that hides that state. The radio Go constructor also releases
its owned RadioButtonData. Independent failures reproduced the missing callback
and leaked data; corrected contracts cover radio and slider allocation release,
ordinary-event rejection, callback order and newly queued cleanup. This is a
reversible correctness decision under the user's standing authorization, subject
to broad qualification and later review. See [CLIENT_RADIO.md](CLIENT_RADIO.md).

## Text-entry input ownership and bounds

Before freezing the entry baseline, release active text-input ownership from the
entry destructor itself: shared GUI destruction suppresses ordinary focus loss.
The original focused-destruction contract failed; disabling text input only when
this entry owns it is a local correction that preserves other active entries.
Bound both composition copies to 255 UTF-16 units plus a terminator, and stop
text scrolling at the main-text terminator even if composition alone cannot fit.
The real input owner has no composition length limit, and the old draw loop had
no terminating case for that width condition. Do not use invalid memory reads or
writes as the compatibility oracle. Keep all other image/color clipping, signed
limits, raw code-unit editing, notification and key-state behavior. These are
reversible corrections to review later; qualification is recorded in
[CLIENT_ENTRY.md](CLIENT_ENTRY.md).


The entry conversion replaces the unused four-byte C context token with a Go
boolean and a typed active-window pointer. Context teardown retains its original
focus behavior, while entry destruction releases its own active input state.
Keep CRT wide-character classification for locale compatibility on the current
CGO target. All twelve routines qualified without changing the frozen C results;
review these compatibility choices separately from future platform work.


### Listbox baseline corrections — review later

Before translating listbox widgets, correct byte-offset scaling in middle row
insertion (independent original-C order regression), bounded row/label terminators,
selection sentinel allocation and overlapping shift, negative head-removal bounds,
capacity-limited scroll lookup, non-progressing auto-scroll and empty-string draw
clipping. Guard the full-selection empty-area click's extra sentinel write.
These are local reversible fixes; qualification and exact preserved quirks are in
[CLIENT_LISTBOX.md](CLIENT_LISTBOX.md). Keep the existing unusual head-removal
selection comparison for now. Represent the ABI selection union as a uint32 word,
which holds either a scalar index or an array address, preserving size/offsets.


### Window ID-range termination — review later

Break after processing the final inclusive ID in hide/enable range loops. The
original signed loop counter wraps at INT_MAX; first=last=INT_MAX should process
one child and finish. Source audit identified the nontermination; a corrected-C
contract checks the actual maximum-ID child. This adds two C prerequisite lines
and preserves ordinary ranges. See [CLIENT_WINDOW.md](CLIENT_WINDOW.md).


## Item tooltip native assembly bounds

The item-name assembler now uses its mapped1,024-unit scratch capacity and
truncates assembled names to1,023 raw UTF16 units plus a terminator. Original C
concatenation had no defined in-allocation behavior for oversized combined names;
those cases are not frozen as C goldens. Exact-capacity1,023-unit names are covered
by the original-C contract, and a native-only2,048-unit base-name contract checks
truncation with modifiers. This is a confident reversible choice to review later,
under the user's standing authorization. The shared variadic missing-equipment
formatter remains unchanged and outside this bounded assembly path; this decision
does not claim that the shared formatter has been ported or generally bounded.

## Meter label capacity and zero-maximum mini-bars

Before freezing the meter baseline, expand sub_471450's four-unit UTF16 decimal
label buffer to12 units, sufficient for signed32-bit values including their sign
and terminator. Establish the original overflow by its source capacity and
formatting contract, not by treating a stack overwrite as a golden. Exercise
999/1000,65535 and signed limits after the correction.

For nox_xxx_drawHealthManaBar_471C00, use zero fill height when maximum is zero.
The corresponding tube already has an empty-state branch, and newly cleared
meter records have a zero maximum. Preserve existing nonzero-maximum arithmetic.
Check the empty case explicitly through both draw paths. Both choices are local,
reversible and made under the user's standing authorization; review them later.
They do not change physical C LOC. Both corrections passed C-baseline and native
qualification; see CLIENT_METERS.md for the completed evidence.


## Meter combined tooltip storage — review later

The original weapon/quiver tooltip concatenates into mapped UTF16 storage at
5D4594+1091968; the next live field, a poison color, is at+1092992. The native
assembler keeps the text within512 units including its terminator. Ordinary
C-baseline names are unchanged. A native-only900-character name contract checks
511 stored characters, the255-character cursor bound, and an unchanged adjacent
color. The oversized original C write is outside the buffer and is not executed
as an oracle. This follows the preceding item-tooltip bounded-assembly decision.

Move the private seven20-byte meter records into a Go array. The caller audit
found no remaining C consumer outside the converted batch. Keep only actual
remaining C callers and callback interfaces; the fixture borrows the native
owner after conversion instead of retaining an otherwise-unused C array.

## Client inventory search row bound — review later

Before freezing the inventory-query baseline, change sub_461EF0's row loop from
<= NOX_INVENTORY_ROW_COUNT to < NOX_INVENTORY_ROW_COUNT. Its declared84-cell
storage has rows0..20 in four columns; row21 aliases subsequent columns and then
reads past the final cell. Keep the extra valid21st row for searches and the
existing20-row limit for visible coordinate queries. Cover last-cell/last-stack
hits, missing codes and duplicate search priority. This fixes an invalid read
without changing valid cell behavior; the original invalid read is not an oracle.

## Native client inventory stack-code bound — review later

The shared148-byte cell contains32 stored item codes. Bound native code searches
to min(count,32) when the stored count byte is larger. The old C loop would read
past the code array. Preserve count getters and all valid0..32 behavior. A native
contract verifies first/last stored-code hits, missing-code termination and the
unchanged raw255 count getter; no invalid C read is used as a baseline.

## Hallway qualification mismatch — investigate on recurrence

The first inventory native accumulated run had one hallway-route hash mismatch.
No source or oracle changed before hallway-alone, inventory-plus-hallway, exact
preceding-prefix and full accumulated repeats all passed. Preserve the first
failure and full repeat captures in the inventory report. Its cause is unknown;
there is no evidence to call it fixed or attribute it to overlapping test jobs.
Proceed with the qualified reversible inventory conversion under standing user
authorization, and add automatic full mismatch capture as a separate diagnostic
follow-up so any recurrence has inspectable data without another replay.

## Client inventory transaction prerequisites — review later

Before freezing the transaction baseline, derive existing-stack pickup coordinates
from its returned cell; the old local pair was uninitialized on that path. This
adds two C lines, checked across all visible inventory cells. Also address the
compaction copy destination directly by existing cell row/column. GDB demonstrated
an in-bounds 148-byte copy whose generated fortified destination size was zero
with the old rolling pointer. Keep fortification enabled and preserve the copy's
source, length and traversal. Both are reversible local corrections, made under
the standing authorization. See [CLIENT_INVENTORY_TRANSACTIONS.md](CLIENT_INVENTORY_TRANSACTIONS.md).

## Inventory display localization input — review later

The fresh four-screen C comparison initially differed only in the shirt's random
localized description. Seed the test process's existing math/rand source with
GODEBUG=randautoseed=0, then capture and compare full screens. The seeded capture
and fresh comparison pass; production localization and the compared pixel area
are unchanged. Preserve the failed run and diff. This is deterministic test input,
not a claim that all gameplay randomness is controlled by this setting.

Display C baseline changes are all behind porttest build tags. Verify unchanged
production source before reusing the immediately preceding qualified three
binaries/full asset result. Run fresh affected tests on all three targets and the
new display scenario; after translation, build and qualify all targets again.

## Inventory identify header bound — review later

Native identification limits its heading to 255 UTF16 units plus a terminator in
the declared 256-unit buffer, using the existing bounded-copy helper. The old C
could leave its same-sized temporary unterminated and concatenate beyond the
heading on long item names. Preserve valid headings, and test the payload boundary,
long ASCII/Greek/emoji names and adjacent storage without using invalid C writes
as an oracle. Existing display hashes remain unchanged.

### Keep spell-force scalar callback words out of Go pointer slots

During inventory-display qualification, GC rejected raw distance word 0x3f8ccccd
in the earlier force port's CallVoidPtr3 frame. Use the existing uintptr callback
bridge and uintptr for its opaque argument throughout native callers. This is a
required correctness fix, preserves the 386 callback words, and does not disable
GC/checks or alter the oracle. Add collection inside the recorded C callback to
exercise the boundary. Reversible and authorized; review the shared callback
convention when the target ABI changes. Full failure evidence remains in the
inventory-display build directory; qualification is complete as recorded in
CLIENT_INVENTORY_DISPLAY.md.

### Define the inventory trade-hit boolean before freezing C

The inventory main mouse handler wrote only LOBYTE(v14) from sub_479880's bool
return and then tested all of uninitialized v14. Assign the complete boolean
word instead. The intended false path preserves the previous identified item;
the true path selects the actual trade-grid item and stack code. Add direct
contracts before freezing this batch. No C lines are added/removed. This is a
small reversible decompilation correction under standing authorization; review
with CLIENT_INVENTORY_WINDOW.md. The C baseline and native conversion are now
qualified in 0842b2d1 and 24f3f67e respectively.

### Inventory window baseline: diagnostics and cancellation ownership

The 34-routine window batch freezes actual owners and 2,120 results after a
byte-identical repeat; authored minimal resources exercise the real parser and
nine gameplay screens exercise the original assets. A broader identification
mismatch did not reproduce in isolated/order/broader repeats. Keep the original
hash and failure visible; automatic full failure capture now applies to display
tests. Do not infer a cause from passing repeats. See CLIENT_INVENTORY_WINDOW.md.

The source/ledger audit shows cancellation retains a temporary non-equipped drag
drawable. Preserve that deterministic allocation behavior for this conversion,
then fix it in a separate qualified Go cleanup chunk with lifetime contracts.
This sequencing is reversible and recorded for review; no behavior is hidden
from the frozen captures or pool checks.

### Release the owned inventory drag on cancellation — review later

After qualifying the unchanged C behavior in 24f3f67e, add deletion of the
non-equipped temporary drag after restoration attempts. Placement/pickup copies
into a separate drawable; equipment drags borrow a live drawable and are excluded
from deletion. Normal mouse release already makes this ownership distinction.
Full-inventory failure still reports the existing error and clears the drag.

Independent tests first demonstrate the old leak, then cover 800 repeated cycles,
fallback and full-inventory failure, identity preservation, exactly-once deletion
and idempotence. Six captured cancellation results change only the precise
lifetime/deletion/pool fields; applying that expected correction to the original
captures reproduces the new captures exactly. The other 27 groups and every
unrelated field stay unchanged. Only two hashes are updated. This is a small,
reversible correctness fix under standing authorization; see
[CLIENT_INVENTORY_CANCEL.md](CLIENT_INVENTORY_CANCEL.md) for qualification.


### Trade UI C baseline prerequisites — review later

Before replacing the connected quantity dialog/trade routines, independent
contracts justify complete grid-pointer returns/assignments, scalar packed mouse
coordinates, bounds and allocation checks, safe missing-item removal, quantity
replacement/destruction ownership, and clearing/finalizing trade drags on reset.
Quantity replacement allocates first so failure preserves the existing dialog.
Reset deletes a detached drag only when neither grid owns it, including a queued
item update arriving during the drag. Quantity destruction also clears its active
flag. Together these add 35 physical C lines before conversion.

Trade reset's three labels borrowed a stack text buffer. Give that local C buffer
static storage until conversion; the Go owner must provide equally stable text.
This narrowly repairs the demonstrated lifetime bug without changing shared
static-text widget ownership or normalizing away differing strings. Separate C
captures and direct label assertions establish stability. See
[CLIENT_TRADE_UI.md](CLIENT_TRADE_UI.md) for original failures and qualification.

A repeated lifecycle fixture also exhausted the GUI pool because it omitted the
normal deferred FreeDestroyed pass. Correcting that fixture is sufficient for
this batch. Review generic GUI allocation failure separately: NewWindowRaw
passes a nil allocation to setExt, whose panic can retain the extension mutex
and block cleanup. That generic path is not fixed by this port.

### Native quantity/trade interfaces — review later

Retain eleven actual C interfaces and retire twenty-five private ones after
moving Go callers to the new owners. The quantity callback receives a C-owned
point and four scalar argument words; forced collection inside the callback
checks this actual boundary. The add-report return is a numeric price, so its
header/export now declares uint32_t instead of char*. Existing decoder callers
ignore that result; all 32 bits remain unchanged, including high-bit prices.
The quantity modifier parameter is void* to match cgo's generated declaration;
the implementation still only copies the supplied bytes. Stable interned reset
labels preserve the corrected C buffer's borrowing lifetime.

These reversible type/lifetime choices preserve the frozen C behavior and avoid
representing scalar game data as Go pointers. Keep the shared production C
number parser until its own connected port; it is not a test-only reference.
See [CLIENT_TRADE_UI.md](CLIENT_TRADE_UI.md) for qualification and the explicit
remaining gameplay integration limits.

### Shop UI C prerequisites — review later

Before freezing the shop UI baseline, independent contracts show full 32-code
stacks are still selected, failed quantity allocations retain sell/repair pending
flags, and count*price overflow can offer unaffordable goods. Skip full stacks
and guard the append. Compare affordable count with gold/unit-price when price
is nonzero, retaining free-goods behavior. Arm a sell/repair request only after
the quantity owner publishes a replacement drawable; allocation failure leaves
the existing dialog intact and permits retry. These are small reversible
correctness changes under standing authorization. Preserve the original failing
runs and verify normal request/callback behavior before freezing the C oracle.
See [CLIENT_SHOP_UI.md](CLIENT_SHOP_UI.md) for observed cases and qualification.


### Shop closure owns its quantity dialog — review later

Independent C tests reproduced retained quantity drawables/windows when closing
or destroying a shop during buy, sell or repair input. Cancel through the actual
quantity lifecycle before shop teardown only when its accept callback is one of
the three known shop callbacks; preserve unrelated inventory quantity dialogs.
Clear both pending flags. Eight direct ownership cases cover close/destroy and
three owned plus one unrelated callback; 100 shop recreation cycles also pass.
This reversible prerequisite adds 11 C lines (22 total for this batch), precedes
baseline freezing, and passes all affected targets/full-assets/gameplay checks.
No algorithm is retained only for testing after translation.


### Qualification scope between subsystem milestones — review later

The complete default corpus at the shop boundary passed 940 roots in 637.877s;
affected server qualification passed 377 in 296.439s. These are wall times in
this VM, not an isolated benchmark of test selection. Repeating unrelated port
families for every smaller connected follow-up has substantial overhead.

For following batches, select accumulated default/server/highres coverage from
actual callers, shared owners/state and dependencies; preserve focused C oracles,
independent contracts, all production builds/interface checks and relevant fresh
gameplay comparisons. Record the exact pattern and selected/completed counts.
Run the complete accumulated corpus at subsystem milestones and for shared
infrastructure changes or uncertain regression scope. Keep the tests and their
expectations available; the change is run frequency, not removal of coverage.
The current shop qualification still runs the complete corpus as planned.

For an unchanged-production C baseline, reuse the immediately preceding qualified
production builds and full-asset results when source identity establishes that
only tagged fixtures have changed. New C behavior fixes invalidate that reuse.
Capture and repeat the new baseline independently before freezing and committing.
This reversible workflow adjustment follows the user's request to improve pace
and autonomously record reasonably confident implementation/testing choices.
