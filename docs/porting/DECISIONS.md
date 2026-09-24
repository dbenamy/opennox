# Porting decisions for later review

## Working policy (user instruction, 2026-09-12)

When reasonably confident in the right answer and reversal would not require
substantial effort, make the decision and record it here or in the batch doc
for later review. Continue without asking for confirmation. Record behavior
changes explicitly, with evidence and validation; do not label them exact
compatibility. Ask when confidence is insufficient or reversal would be costly.

## Inventory display fixture identity normalization

Pass the original return separately to each snapshot so a canonical drawable
identity is normalized only once. The old nested snapshot could reinterpret it as
an image handle, causing an address-dependent capture mismatch. A deterministic
collision test fails on pre-conversion `e64ff24e` and passes with this fixture-only
fix, as do all existing display captures. Goldens remain unchanged. This small
reversible fixture correction is covered by the standing policy; final batch
qualification passes in all three profiles, safe build and production gates. See [GO_PRIMITIVE_INTERFACES.md](GO_PRIMITIVE_INTERFACES.md).

## Green-bolt effect record — correction qualified

Original client message 152 adds 432 to a `nox_drawable*`. Since that type is
512 bytes (asserted in defs.h), this advances 221,184 bytes, not the intended
432-byte effect-data offset. Independent message contracts found the allocated
record untouched. The production handler now checks allocation success and adds
432 bytes before writing the existing 13-byte packed record. Allocation failure
returns the normal consumed length, 11, without effect writes.

This is an explicit bug fix, not exact preservation of undefined out-of-object
writes. Chosen under the standing reversible-decision policy. Preserve unsigned
coordinate midpoint truncation, cache initialization before the connection gate,
input bytes and packed field ordering. Fresh three-target and production qualification pass, including ABI, the exact
known full-suite results, and headless gameplay/save-load. See the tracked
[qualification](green-bolt-correction-qualification.json).
Status and validation: [GAME_MESSAGES.md](GAME_MESSAGES.md), PORTING_STATE.md.

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


### Journal shared storage and presentation

Journal list/report/rendering logic is Go, with libc allocation/free for entries
shared by the remaining C decoder and save/load layout. This preserves the old
nil-on-allocation-failure behavior and C-compatible lifetime without retaining a
C algorithm. Direct Go inputs retain first-NUL semantics and 63-byte storage.
Go strings replace the presentation scratch buffers; frozen valid-input layout
and pixels remain identical. Local flag updates still omit cached-height rebuild,
as did both original callers; no correction is claimed. Review these choices if
journal ownership/layout or update behavior changes. See JOURNAL.md (C baseline
2c3ea111) for actual save/load, report, render and gameplay evidence.


### Briefing presentation compatibility

Use typed Go chapter/score layouts backed by the existing shared storage and
stable descending sorting of the original participating-record count. The server
packs its six score slots; preserve bounded sparse-input behavior without adding
client compaction. Keep observed equal-score order. Go strings replace formatting
scratch buffers, preserving UTF-16 code-unit name clipping and known presentation
quirks (the `XX1` label, SoulGate paragraph edge and mixed prompt fonts). Review
these if changing briefing presentation, rather than silently correcting them
during translation. All thirteen routines are Go; actual window lifecycle,
coordinate and sprite callback owners remain integration dependencies.

The new chapter scenario samples early frames because the existing E2E harness
shortens briefing duration. Visual inspection rejected the initial late-frame
capture as briefing evidence; the corrected scenario and independent repeat are
tracked/recoverable. Statistics and instructions use real render/GUI/sprite
owners with the shipped Briefing.wnd, while ordinary chapter gameplay supplies
the full-game integration check. See BRIEFING.md for qualification and limits.


### Briefing lifecycle completion

Preserve existing resource-failure ownership, sprite-deletion order, nil event
responses and post-draw state reads when replacing the remaining callbacks with
Go. Keep the actual save-menu boundary outside this batch: quest-mode transition
fixtures exercise its early return, and ordinary chapter gameplay provides the
full-game check. Additional save-selector pixels and audible credits remain
separate coverage work. Review these limits when changing transitions.

Run the complete accumulated corpus in all three targets at this closed briefing
subsystem milestone. All selected roots completed with no failures and the one
expected optional skip in each target. Keep the 600s per-package timeout; no
infrastructure adjustment was necessary for this completed run. Highres root
execution was 598.823s, so give the next complete-corpus run a longer explicit
timeout before adding further coverage; retain strict root completion/failure
checks. Following batches resume the documented
affected-corpus policy until another milestone or uncertain dependency scope.


### Scoreboard compatibility and scratch ownership

Preserve stable score selection, unsigned rank comparisons, observer score
adjustments and their restoration, mode-specific team/player counts and raw
UTF-16 record copying. Long names can overlap the adjacent team field before that
field is rewritten; the narrow-font fixtures preserve the existing bytes and
rendered result. A bounded-name behavioral correction belongs in separate review,
as do unsupported sort sentinels and inconsistent headless-host counts. These
limits are documented rather than silently changed during translation.

Use existing mapped scratch for dynamic score/time text instead of permanently
interning each changing string. Retain only the five interfaces with real C callers;
move private callers to Go and retire the newly private briefing-stage bridge.
Use affected three-target qualification after the prior complete briefing milestone,
plus all production builds, exact known-failure comparison and both hosted-scoreboard
and chapter gameplay. Post-review callback-adapter cleanup repeats the full scoreboard
family in all three targets, with exact production fingerprints to reuse the broader
checks. See SCOREBOARD.md for evidence and limits.


### Minimap compatibility and private interfaces

Preserve the polygon predicate's existing first-query corner-ray miss and subsequent
floor caching; freeze that behavior separately from ordinary asymmetric polygon
cases. A predicate correction is a separate review item, not part of rendering
translation. Preserve zoom wrap/signed comparisons, integer projection, three-sided
shadow borders, current team/objective/observer visibility and real message drawing.

Include the private AI-debug monster traversal with its only renderer caller. Move
zoom/iterator state to Go, invoke Go owners directly and retire both private AI-path
C bridges. Test fixtures use actual path/object/wall/drawable/team owners, with normal
palette definitions supplied as inputs to the lightweight test server. Retain no C
algorithms solely as an oracle. Use affected three-target qualification after the
briefing milestone, plus builds/ABI, exact full-assets comparison and twelve fresh
minimap gameplay frames. Solo gameplay does not establish remote multiplayer
coverage; populated team/objective and debug behavior have actual-owner fixtures.


### World-wall ownership and visibility compatibility

Use existing Viewport methods for both private projection helpers and route
briefing through them. Keep actual static callback identities, player/team lookups,
vision buff eligibility, FOV scanline intersection and shared image-interval state.
Copy the first light sample before obtaining the second: the real sampler reuses
its buffer. Preserve signed-byte interlacing parity, front/back/translucent options,
sprite overrides, consumed wall flags and explored state. Retire all nine private
C interfaces; retain the actual edge renderer and shared C configuration owners.

The server's actual FOV clipper is intentionally unreachable. Run those 528 cases
only in client targets; all other 5,734 records execute on all three targets,
including ordinary/edge rendering. Do not install a replacement clipper for tests.
Pin the real Go pixel buffer while the production C row table retains its addresses;
cleanup detaches/frees the table before unpinning. Full-assets failure comparison
remains exactly the known 1,553 entries, and fresh solo chapter/minimap gameplay
matches twelve frames. Review pixel-row ownership consolidation with the remaining
C tile callers later; it is separate from this behavioral translation.


### Wall-edge rasterizer compatibility and private API

The only production caller discards the original edge routine's scratch/pointer
return, and its flags argument is unused. Use a private void Go API without the
flags argument; freeze flag-independence before conversion. Retire the edge C
interface plus light-multiplication and pitch exports whose only C caller was the
edge renderer. The underlying Go owners remain, and no test-only C algorithm stays.

Preserve low-resolution copy length exactly: a leading partial transparent run
contributes nothing, while later transparent runs count their full length, including
a run extending past the horizontal interval. Distinct background rows, partial-run
cases and independent alternate-row contracts make this behavior visible. A visual
change to that legacy behavior is a separate review item. Keep the real shared C
pixel-row/clip/configuration owners until their remaining tile callers move. Invalid
RLE buffers are outside the valid asset contract; actual image-owner load failures
remain unchanged. Three-target frozen comparisons, actual caller tests, builds/ABI,
exact known failures and fresh 12-frame gameplay qualify this conversion.


### Tile raster callback ABI, fill compatibility and integration

Normalize the two retained callbacks to their actual void/three-argument dispatch
slot. No caller consumes the decompiled scratch returns; fill narrows the tile
argument to uint16, confirmed with nonzero upper-word inputs before conversion.
Retire private raster/setup/getter C interfaces, retaining their Go owners and
shared C configuration storage. No test-only C algorithm remains.

Preserve fast-path full-width fill phase and split-path pattern restart; uniform
16-bit filling would silently change existing behavior. Compact loops use the same
46-row diamond layout, backed by independent whole-buffer contracts and frozen C.

E2E intentionally ignores nox.cfg. Flat-floor integration must use the real options
checkbox with a hover tick before clicking; a config-only replay or instantaneous
click does not establish that mode. Preserve the actual GUI scenario and pixel
manifest, and compare both normal12 and GUI-flat14 frames. Flagged tile definitions
may temporarily restore textures, as before. Invalid assets remain outside these
valid-owner contracts; this batch makes no image-decoder behavior change.


### Tile composition baseline and milestone timeout

Preserve the original redraw predicate's unsigned X calculation near the map
origin, independently of full redraw's signed calculation; fixtures establish the
difference before translation. A behavior correction would be a separate review
item. Overlay inputs must respect the existing row-start and one-initial-wrap
preconditions. Retain evidence from invalid fixture positions and do not count
those crashes as qualified output. Shared image ownership matters: tile and overlay
handles must be registered in the same actual bag.

Move the private callback selection/state to Go with its final composition callers;
retire the two raster exports, edge callback, no-op and C slots. Preserve actual
shared counter/grid/definition owners and all frozen behavioral expectations.

Run the complete accumulated port corpus at the native floor-rendering milestone.
The previous highres root took 598.823s under a 600-second timeout; expose a positive
--timeout-seconds option (default 600 unchanged), record it in results and use 900
for this larger milestone. Do not silently omit newer test families: the tracked
complete pattern now includes wall-edge, tile raster and composition tests.

## Floor/edge asset-reader compatibility and empty-table correction

Port the connected definition, skip, image-binding and free helpers together.
Before freezing C, make empty floor-definition binding return the existing
missing-definition failure after reading the name. The original path used a
pointer-derived index; the edge binder already rejects count zero. This three-line
correction is narrow and reversible, with an independent cursor/no-mutation test;
review it separately if changing loader error semantics later. Normal loading
order initializes definitions first. Fresh C production qualification is required.

Preserve the differing END consumption, partial failed-edge writes, C-string name
termination/padding, format-byte behavior and five-byte edge allocation extent.
Use the real startup facade strings/pointers and actual image resolver in tests;
inline names currently resolve nil. Keep raw calloc/free ownership compatible,
without retaining C reader algorithms solely for the oracle. Retire the now-unused
width/height C adapters in the same Go batch, retaining their live Go owners.

## Things-section traversal and public scratch-buffer views

Port the connected section readers together, including the aligned MemFile helper
whose only callers are in the wall reader. Preserve absolute eight-byte alignment,
full eight-byte consumption for count fields, partial scratch writes and early
returns. Keep one AVNT inner export while the client event loader still calls it.
Original-C contracts cover all tag bytes, count/name boundaries, concatenated
sections and public wall buffer views. Do not retain C algorithms just for tests.

Review found a compatibility edge in the preceding Go floor/edge readers. Their
original wrappers check backing capacity and pass `&buf[0]` to C, so any nonempty
view with enough capacity behaves like the full buffer; zero length is rejected
before execution. Actual game callers pass full buffers, which is why replay did
not reveal the shorter-view difference. Preserve that public contract in the
native section batch: keep capacity guards, check element zero before reslicing,
and expose capacity to the Go reader. This applies to five floor/edge wrappers
and the new wall reader. Add the independent floor-view equivalence/rejection
contract; retain all existing frozen floor expectations. This is a reversible
compatibility correction within the standing authorization, not a new parsing
policy or relaxed input requirement.

### Map decoder valid files and malformed/empty input (process trial round 2)

Move the complete NXZ decoder into local Go with owned arrays; retain the C
compressor because it remains production code. Preserve the dictionary, signed
16-bit frequencies and symbol order across byte-aligned blocks. Do not directly
substitute the pinned libs/nxz reader: its end-block handling does not perform the
alignment observed in C. Freeze valid C files and use independent map bytes and
hand-encoded block boundaries to qualify the replacement.

Review later: return an error for truncated/invalid bodies before creating the
destination, and accept a zero-output-size header as an empty map. The existing
C wrapper can panic on zero allocation; its unchecked body reads are not a sound
malformed-input oracle. These bounded, reversible behavior corrections fall under
the user's standing authorization. Keep malformed-input and empty-output tests
explicit; do not change any valid-file expectation. Empty-input compression is
outside this decoder scope and still has its existing allocation limitation.

### Map compressor ownership and bounded lookahead

Complete NXZ compression in Go and remove its last two production C files.
Retain exact encoded bytes, including following-file lookahead at 500,000-byte
block boundaries. Final-file lookahead is explicitly padded; the C version could
read beyond its allocation while maintaining hashes unused after the final block.
Review later: empty input now produces the four-byte zero-size header, and inputs
exceeding the wrapper's addressable allocation or format size return errors before
creating the destination. These replace allocation panic/unchecked narrowing;
independent tests cover both. No valid-file baseline was changed.

### Quest rotation with no eligible candidate

Before freezing the catalog baseline, add a fallback to uniformly select from the
full nonempty catalog when family and recent-history exclusions remove every
candidate. This occurs in valid small catalogs; the original reversed RNG bounds
returned -1 without a draw and produced a name pointer before the table. Review
later: fallback may repeat a recent family, and now consumes one RNG draw for a
catalog larger than one. Empty/singleton and successful filtered choices preserve
their existing contracts. The correction is three C lines, reversible, and tested
independently before freezing 1,789 catalog/cycle records.

Preserve the cycle cursor's strict `index > count` comparison and unchecked first
non-section mapcycle line. These are observable existing behavior; changing them
is outside this compatibility batch.

### Shared list owner and player-group names

Consolidate the duplicate root list implementation on the Go owner used by legacy
callers and retained C interfaces. Preserve the canonical C behavior for nil
successors; the old root Go helper could dereference nil on a zeroed list. Add an
independent root-wrapper contract. Preserve partial append behavior on an
uninitialized head, signed sort keys, duplicate order and exact link layout.

Review later: bound group names to nine raw UTF16 units plus a terminator in the
existing ten-word field. Valid C names are unchanged, including raw surrogate
units; overlong names have no valid C-buffer oracle. Native contracts cover
lengths 0/1/9/10/11/128. Keep the shared C-compatible allocation/free ownership.
Use the full accumulated corpus as this shared batch's all-target behavior gate,
in place of a redundant affected-only sweep; retain production and replay gates.

## Spellbook baseline — preserved behavior and replay preparation

Preserve the original quest-mode difference between guide-icon drawing and
pressing: drawing requires a known summon spell; pressing permits dragging
without it. This is captured, not corrected, during the mechanical port.
A fresh warrior has an empty book, so the gameplay scenario uses the existing
`cheat spells` developer game command to grant its five abilities. `NOX_DEV=true`
is scoped to this scenario in production qualification. No production fixture
backdoor or changes to the shipped game data are needed. Key taps use one frame;
ASCII console text uses individual key events. Pointer inputs are in the 1280x960
window coordinate space, scaled to the 1024x768 gameplay renderer.
See [SPELLBOOK.md](SPELLBOOK.md) for frozen tests and rejected replay evidence.

### Quickbar collapse restores the row saved by expansion

During the quickbar baseline audit, the complete real-window transition confirmed
that expansion saved the selected row at 1047912 while collapse read unwritten
1047908, returning every starting row to zero. Read 1047912 on collapse. This is a
reversible prerequisite correction before freezing the final port baseline.
Seven focused C roots pass; only restored-row/pointer/direction-indicator fields
change, and six unrelated capture hashes remain exact. Original behavior is
preserved in fdc5048e and quickbar-original-captures.json. See QUICKBAR.md for the
remaining batch qualification gates.

## Summon command menu bottom clamp — reversible C prerequisite

The original C command-menu constructor uses screen width when clamping the
vertical position at the bottom edge. A 640×480 corner contract produces a menu
rectangle (0,562)–(148,639), entirely below the screen. Use screen height in that
one calculation before freezing the summon baseline. The twelve corner contracts
cover 640×480, 1024×768 and portrait 480×640 screens and require the menu to fit.
This is an intentional behavior correction, not an exact match to the old bug.
Reversal is a one-word source change plus its explicitly changed expectations.
Evidence: build/port-summon/menu-bounds-before.log and TestSummonMenuBounds;
see [SUMMON.md](SUMMON.md). C LOC is unchanged.

### Refresh binding-editor prompt text when opening it — review later

Headless gameplay showed a blank prompt in the C binding editor. Construction
initializes the static-text widget with an empty buffer, which clears its text
pointer. Updating the buffer later does not update the widget. A regression test
reproduces this in both the in-game and main-menu editors.

Refresh child 981 after formatting each prompt, in the two existing C routers
before translation. This keeps the correction local and preserves the existing
static-text API for other callers. The visible prompt should name the selected
action every time, including after reopening it. See [BINDINGS.md](BINDINGS.md)
for qualification status. This reversible correction is authorized by the standing
instruction to act on confident decisions and record them for later review.

### Compiler-cache probe after binding editors

The installed ccache 4.12.3 can reuse an unchanged GAME3.c object across different
Go-style temporary paths. With GCC's switch recording disabled as in the Go build,
all probe object hashes agree: uncached 0.821s, direct hit 0.006s, changed-work-path
hit 0.097s. However, changing the compiler seed causes a miss (0.975s). Go 1.26
passes its package action ID as `-frandom-seed`, so a source revision changes that
seed. Evidence: build/port-compiler-cache/result-seeds.json and driver-seeds.log.

Decision for later review: keep the current plain-GCC environment. The probe does
not establish savings across port revisions. Ignoring the seed would relax cache
identity and needs broader object/build validation before adoption; one unchanged
object is insufficient evidence. This is a deferred optimization, not a blocker.
No toolchain, compiler flags, production source or qualification expectations were
changed by this experiment. The isolated cache is limited to 512MiB.

### Options volume checkbox dispatcher prerequisite

Use the existing window event dispatcher for the options panels' twenty calls to
mute-checkbox handlers. The direct C function-pointer reads bypass handlers stored
in Go window extensions. The independent zero-volume regression fails for all
three channels in both panels against 0649e67a's C and passes after this scoped
correction. The full volume matrix includes actual checkboxes, absent handlers,
audio readiness, preview ordering and real timer state.

The corrected C baseline qualifies 3,130 frozen records on all three targets plus
repeat, all 56 affected checks, and 41 repeated real gameplay frames. See
OPTIONS.md. This is an intentional reversible behavior correction to review;
it precedes Go translation. It removes 96 C lines of raw-pointer plumbing and
redundant empty-handler branches; record that separately from ported C LOC.
The transient clipping frame is recorded for a separate renderer investigation.

## Client map reader defaults and preserved edge cases

Initialize the old common record's team byte and extra-flags word to zero when
older layouts omit them. All 42 independent regression cases fail against the
unmodified C reader and pass after the two initializers. This reversible C
prerequisite is qualified before translation; no physical C lines are removed.
See [MAP_DRAWABLES.md](MAP_DRAWABLES.md).

Preserve old count narrowing, signed-X/unsigned-Y door positioning, and distinct
legacy versus modern team-registration exceptions. The 2,484-record baseline
uses actual tables, file/drawable/shape/light/modifier/wall/team owners and seven
frozen capture groups.

Keep the existing allocation-failure section-framing behavior during this port.
A typed reader can consume a partial record and return zero, after which the
section owner skips the full declared length. Controlled cases reproduce the
resulting stop positions and load-error flag. A later fix should define recovery
for the complete map-loading owner; it is separate from defaulting absent fields.

## Colored-light degenerate direction prerequisite

While preparing the colored-light animation batch after b593c7ae, eight independent
regressions showed that coincident targets and zero-width rotation arcs change the
angle via NaN-to-integer conversion. Preserve the existing light state in both
cases with two early returns, before freezing C expectations. This is a small,
reversible behavior correction authorized by the working plan; review later if
zero-width arcs should explicitly reset to their configured starting angle.
See COLOR_LIGHT.md for before/after evidence and qualification status.

## Object-transfer rejection ownership — prerequisite correction

Return a newly allocated inventory object to the real object pool when its
transfer callback rejects the record. Also clear a rejected parent's inventory
head after disposing its children, before disposing the parent. The first defect
leaves an unreachable live object; the second lets the normal destructor traverse
children that were already freed. Both changes are in GAME3_3.c before freezing
the server serialization baseline. Four C lines are added.

The failed-callback regression observes two live objects instead of one. The
placement regression faults when the destructor follows the freed inventory.
Evidence: build/port-object-xfer/c-failed-callback-before.log and
c-rejected-inventory-before.log. After-change default/server/highres checks pass: ten roots, 681 subcases plus
two ownership regressions, no skips. Static checking also passes; see
[OBJECT_XFER.md](OBJECT_XFER.md). These are deliberate ownership fixes,
not exact preservation of the broken failure paths. Successful-load formats and
return/stream behavior are preserved. Existing general buffer cleanup policy is
outside this correction; the failed-callback regression isolates the pool slot
using a type with no side buffers. Reversal is small, but retaining a leak and
a freed-object traversal is not recommended.

## C object-factory adapter — stale type prerequisite

Route nox_xxx_newObjectWithTypeInd_4E3450 through Server.NewObjectByTypeInd, which
already returns nil for a missing type. Its previous Go adapter called the pool
factory directly with Types.ByInd's nil result. The new stale-TOC regression
reaches that path from the real inventory loader and panics in ObjectType.Ind2.
The intended loader failure path already checks for a nil allocation.

This extends the guarded-factory choice previously used by the native map-object
admission port to remaining C callers. Valid allocations retain the same factory
and object initialization. It is a reversible adapter correction, with no C LOC
change. Before-change evidence is c-final-default/tests.jsonl and the matching
target checks under build/port-object-xfer; corrected c-final-default2, c-final-repeat,
c-final-server2 and c-final-highres2 each pass all 21 roots and twelve hashes.
Current-source gameplay, actual save/load and flat/map regeneration also pass.
Frozen expectations for previously covered valid records are unchanged.

## Object-transfer trigger callback signature

Use an object-pointer parameter for the Go-backed TriggerXfer export. The old C
float declaration immediately reinterprets its first 386 stack slot as an object
address and later reuses that variable as numeric scratch. Registered transfer
callbacks pass object pointers; the remaining C reference compares function
identity. No numeric-float caller was found. Keep the symbol and callback identity,
correct its header to nox_object_t*, and never numerically convert an address to
float. This is an ABI description correction for the supported 386 target;
focused registered-callback and full production ABI checks qualify it.

## Item serialization prerequisites — generator child and old wand attributes

The next typed-serialization baseline exposes two existing defects. Free a newly
allocated generator child when its registered transfer callback rejects the record;
the before-change contract observes two live pool objects instead of one. Define
the final attribute word in pre-version-11 weapon records as 0xffffffff, matching
the two 0xffff words supplied by the version-11 path. A real charged wand copies
uninitialized stack bytes before this correction; four old-version cases fail,
while the version-11 control passes.

The initial attribute fixture used ClassWeapon and missed the charged-wand branch;
corrected before-change evidence uses ClassWand and is recorded separately. Normal
armor does not take that branch, so its early-return path is left unchanged.
These two C lines are small reversible corrections authorized by the working plan.
They precede frozen expectations; they do not change the current writer format.
After-change default/server/highres qualification passes 29 roots / 1,569 leaf
cases per target without skips, and static checking passes. See ITEM_XFER.md
for evidence, the fixture correction and the remaining full-baseline work.

## Item serialization — preserve established records and callback ABI

The item/reward port keeps historical field widths and stream operation boundaries,
including zero-length writes. Reward-marker counts include only mask value 1 while
the writer emits names for any nonzero mask; explicit writer-only contracts preserve
this quirk instead of inventing a format migration. Reads merge masks and keep
already accepted entries on later failure. Sparse generator children compact within
each row on reload and retain ownership of earlier successfully loaded children.

Nested generator tests use the same cipher framing as real saves. An initial XOR
fixture exposed a writable seek/backpatch limitation in that alternate mode; changing
shared stream behavior was unnecessary for qualifying this port. Defined serialized
fields are checked independently, with full wire hashes also pinning compatibility
padding. The empty obelisk hook remains for objective-update callers; the translated
serializer omits its no-op call. Existing callback header declarations and symbol
identities remain unchanged, with address conversion restricted to the C bridges.
See ITEM_XFER.md for frozen evidence and final qualification status.

## Visibility/effects — orphan removal and entry-point contracts

Remove the 47-line `sub_528030` throttle helper after a full reference audit finds
only its definition, declaration and the new baseline fixture. No production
caller, callback registration or dynamic symbol-lookup path reaches it. Its C
baseline and 3,360 historical records remain recoverable at `e2616865`; the native
implementation and test corpus omit this unreachable helper. The signed-mana
finding in that baseline is therefore not an observed player-facing problem.
Perform this reachability check before building future baseline fixtures.

Preserve each effect entry point's coordinate domain: vampire effects cull using
unsigned destination low words; generator-spawn effects use full signed destination
coordinates. Reuse the existing Go viewport method only after paired C/Go contracts
verify strict boundaries, observer-camera selection and asymmetric float stores.
Six live C entry points remain; 24 C symbols are retired, with Go callers routed
directly. No live gameplay behavior correction is included in this batch.

Supplemental original-C scan-delay tests use real circular collider bounds. The
old minimal fixture left bounds at the origin, which accidentally lay inside its
small test region. Correct bounds preserve the original spatial capture and allow
864 independent distant/empty-search scheduling contracts. See
[VISIBILITY_EFFECTS.md](VISIBILITY_EFFECTS.md) and
[VISIBILITY_SCAN_DELAY.md](VISIBILITY_SCAN_DELAY.md) for evidence and validation.


### Object and recipient reports

Qualified after C baseline `8b370caa`: move the private minimap-count helper with
its only caller and retire the now-unnecessary audio/visibility C bridges. Preserve
health-word narrowing, history changes even when a later queue write fails, and
signed polygon-level behavior. Polygon vertex/ray edge behavior is captured and
left to the existing geometry owner; no unrelated geometry fix is folded in.

Reuse a preceding qualified production baseline when new changes are test-only
and production source is identical. The new C contracts still run across all three
targets and repeat independently. Completed native batches still receive the broad
affected corpus, fresh builds/ABI, exact known suite and headless integration gates.
See [OBJECT_REPORTS.md](OBJECT_REPORTS.md) for the evidence and fixture corrections.


### Reliable message queue pressure ownership

Fix the independently reproduced original-C cleanup crash before freezing the next
baseline. Slow-player removal can already free the oldest message selected by the
pressure scan. Check the surviving list before unlinking that selection again;
retain acknowledgement/related-bit updates and return success when removal already
released it. Tests verify pool reuse and surviving-list variants. This is a
reversible prerequisite correction within the authorized port workflow, not an
expectation change made to accommodate Go. Fresh C production checks are required
because production source changed. See [RELIABLE_REPORTS.md](RELIABLE_REPORTS.md).


### Reliable queue conversion and evidence counts

Preserve the original shared C-owned allocation/list/rate layout while replacing
queue algorithms, and keep only seven entry points required by C callers. The
pressure ownership fix is already qualified in C baseline `f95e7aee`; the Go
conversion changes no frozen expectation. Go callers avoid the C round trip and
extra payload allocation. Treat unique terminal test names as leaf cases: recounting
the original visibility/object-report logs corrects an 18-case overstatement in
prior documentation. No tests were removed. See [RELIABLE_REPORTS.md](RELIABLE_REPORTS.md).


### World-collision arithmetic and callback contracts

Preserve the qualified C binary's actual floating-point boundaries for mass
exchange, including its float32 coefficient spill and PC53 intermediates. The
initial literal translation differed in 22 of 175 frozen rows; inspecting the
original instructions resolved it without changing any expectation. All focused
contracts and affected target sweeps now match. Callback ordering, raw UTF-16
identity/history bytes, countdown signedness and clock truncation remain explicit.

Retain the twenty C entry points used by callers or callback registration and
retire the private mass-exchange interface. Invalid door angles formerly left
adjacent coordinates uninitialized; Go initializes them to zero. Spell argument
fields that C left uninitialized are also zeroed. These are reversible choices
outside the established valid-input contract, recorded for later review. No C
algorithm is kept solely for tests. See [WORLD_COLLISIONS.md](WORLD_COLLISIONS.md).


### Quest score constant width and fixture data

The shipped quest score exponent is an eight-byte double near 1.9. Original C
read it as long double; the independent shipped-data regression returned
2,147,483,648 instead of 37 for stage 2/ten points. Change the read to
getMemDoublePtr before translating the quest runtime. This reversible prerequisite
fix is within the authorized workflow and requires fresh C production gates.

The first fixture left the exponent zero and masked this bug. It was caught in
source review before the C baseline commit. All quest fixtures now install and
restore the shipped constant. Only three scoring groups are refrozen from corrected
C; the other 20 groups remain identical. Numeric/table input audits are now explicit
in PORT.md so nontrivial production data is checked before freezing expectations.
See [QUEST_RUNTIME.md](QUEST_RUNTIME.md) for qualification and limitations.


### Quest runtime native interfaces and absolute conversion

The native conversion preserves all frozen corrected-C behavior, including signed
maximum-health comparisons, float32 rounding points, all-roster score aggregation,
and participation equality/nonzero distinctions. Retire the private C absolute
converter with its last four production callers; keep independent Go assertions
and the two C converters that still have production callers.

The health-scaling C export now returns void. All remaining callers ignored the
old decompiler return, whose temporary value had no coherent result contract.
This reversible interface cleanup leaves the covered mutations unchanged. The
stage-message builder bounds each name to its 32-byte field; the preserved valid
contract is at most 31 name bytes plus a terminator. See
[QUEST_RUNTIME.md](QUEST_RUNTIME.md) for complete three-target and production gates.

## Match roster prerequisites and C protocol version — review after conversion

Initialize the roster's reused scratch buffer inside each player iteration and
zero the settings server-name padding. Independent long/short/empty-name sequences
reproduced both uninitialized fields before this correction. For a Flagball timeout
without a winning team, send the existing generic flag draw message (87 / 65535 /
timeout marker 1); the Flagball winner reporter requires a real team in both its
old C and current Go implementations. Empty-match testing reproduced the old crash.

Add three explicit bytes after `Player.Active` so the Go identifier starts at 2096,
matching C, with a compile-time offset assertion. Previously it started at 2093,
which the new full C/Go encoder comparison exposed. Total size/later offsets are
unchanged because the explicit padding replaces implicit padding after that array.
These are intentional, reversible corrections under the standing authorization.

Preserve the selected C settings message's protocol value 0x000F039A on every
target. The legacy highres compiler flag is unconditional, whereas the root Go
version is target-specific. Correcting that compiler flag would also change C
rendering constants; defer that broader compatibility change for review.
Evidence, qualification and conversion scope: [MATCH_ROSTER.md](MATCH_ROSTER.md).

## Team runtime message prerequisites — review after conversion

Before freezing the team baseline, zero four fixed-width local message buffers:
rename, team change and the two client requests. Real-C regression cases reproduce
unspecified trailing bytes; preserve lengths and populated fields. The join-member
message separately writes its object type into an unrelated local `short`, not
its outgoing array. Store it at byte 8 of the existing 10-byte record and remove
the unused local. Header bytes already matched the independent contract; the
incorrect type word was reproduced separately. These corrections are authorized,
reversible and intended to make the existing protocol fields deterministic and
correct. The corrected C baseline `0b4b853c` passes all three targets and fresh
production qualification. See [TEAM_RUNTIME.md](TEAM_RUNTIME.md).

The group/rebalance contract additionally reproduces stale team member counts:
the clear-player caller removes membership nodes without decrementing its counter.
Decrement only after a successful unlink; a detached matching-ID entry must not
change the count. Keep the lower-level unlink API unchanged because other callers
already own their decrements. This is another reversible prerequisite correction;
combined corrected-C and production qualification passes in TEAM_RUNTIME.md.

## Team score dispatch ownership — conversion review

Preserve the legacy score sender's direct reliable-queue dispatch. An initial
reuse of Server.TeamChangeLessons matched the default output but routed through
its replaceable send hook, changing the existing objective fixtures' complete
observations. Keep the public Go setter unchanged; the converted legacy owner
stores the score and sends the existing record directly. A new dispatch-owner
contract distinguishes the two paths, and both original objective-score hashes
are restored without changing goldens. PORT.md now calls for this ownership check
when reusing existing Go APIs. This is translation fidelity, not a gameplay rule
change; no additional user decision is needed.

## Team UI prerequisites — review after conversion

Before freezing the team HUD/player-list baseline, independent real-widget
contracts reproduce five existing C issues. Correct them under the standing
permission for confident reversible changes:

- Refresh frees both old backing row lists before rebuilding, keeping metadata
  and displayed rows aligned instead of accumulating duplicates.
- Team lookup compares complete bare names, preserving spaces.
- Rename finds the row by team ID, updates its backing name and preserves the
  current selection; it must not rename a different selected team or insert an
  extra row when nothing is selected.
- Assign-button eligibility compares the player selection as signed, so the -1
  sentinel disables the control.
- Failed player-list resource loading returns zero before accessing child widgets.

The corrected original-C UI corpus passes 25 roots with 19 frozen captures.
Three-target and production qualification are recorded in [TEAM_UI.md](TEAM_UI.md).
These changes add 22 C lines before conversion (51,225 remaining). They change no
message layouts or membership rules. Keep existing empty-name acceptance and
ASCII case folding; this runtime's C comparison does not fold Cyrillic case.
Broader Unicode case handling is a separate review item, not part of this port.

### Server-options baseline prerequisites (review after port)

Independent C contracts found two small lifecycle defects. Missing window resources
cause the server-options constructor to dereference NULL while reading width;
return 0 before layout/child access. Closing options destroys its general panel
but leaves the parent's general-panel pointer set; clear that pointer with the
other panel state. Both corrections are reversible and covered by missing-resource
retry and tab/close contracts. They add four physical C lines before conversion.
Because production source changes, qualify a fresh corrected C baseline rather
than reusing the previous team's production result. See SERVER_OPTIONS.md.


### Server-panel baseline prerequisites and compatibility notes (review after port)

Five panel constructors dereference missing window resources. Independent child-process
contracts reproduced failures for object (both modes), spell, access, advanced-tab
and advanced-server panels. Return zero immediately after a failed root load, then
allow a successful retry. General options already handles the failure. These five
reversible guards add 15 C lines; qualify fresh production before freezing the
baseline. See [SERVER_PANELS.md](SERVER_PANELS.md).

Keep the observed object checkbox index/name mismatch and the upper level-limit
checkbox's unshifted value during this port; separate contracts cover each path.
They are candidates for later behavior cleanup. Preserve byte-oriented weapon-mask
queries and 386 spell shift/index conventions. The panel fixtures now use shipped
multi-selection ownership where required; this is a fixture correction. Player
lookup is case-insensitive, while displayed-row removal is exact.


Server-panel native review also preserves exact UTF-16 code-unit row matching,
change-only button enabling (avoiding recursive changes to already-enabled
children), direct checkbox Y-offset writes, and the advanced window's position
from the original nil-root lookup. Nine thin C interfaces remain for actual C
callers; Go callers and the advanced callback table use Go directly. The numeric
parser, name-list owner and rule picker remain outside this batch. All frozen
captures and fresh production qualification pass; no further behavior correction
was needed during native conversion.


### Server-configuration baseline prerequisites (review after port)

Independent contracts reproduced five existing C defects: a missing rule-picker
resource causes child access through a null root; consecutive expired admission
entries are skipped after removal; libc wide formatting omits the game's UTF-16
names from admission files; Linux rule enumeration reads an empty alternate-name
field; and settings refresh copies a 12-byte map name across its 9-byte field into
the server-name prefix. Correct these locally before freezing: guard the missing
root, advance the list index only for retained entries, narrow names through the
existing game formatter, enumerate actual filenames with adequate local buffers,
and copy exactly eight map-name bytes plus a terminator. They add two net C lines.
See [SERVER_CONFIG.md](SERVER_CONFIG.md) for reproduction and qualification.

Preserve legacy byte encoding and the rule-entry completion callback's narrow view
of wide text. Also preserve repeated dirty notifications caused by signed weapon
bytes, headless empty counts and video flags: identical record bytes do not imply
an unchanged return value. These are explicit later-review items, not silent
conversion fixes. The old disabled report function and its private setter are
proven dead and are scheduled for removal with their no-op calls during translation.


Server-configuration native review preserves the opaque rule-picker background,
including the renderer's existing color; its pixel capture rejected an initially
reused blended-background helper. All frozen expectations remain unchanged.
The port removes two proven disabled helpers and their no-op callers, eight private
C globals and 41 function interfaces. Thirty-two thin exports retain actual C
callers; list allocations/layout and shared configuration storage remain compatible.
The Cgo admission header drops a const qualifier without changing the ABI. See
SERVER_CONFIG.md for qualification and deferred compatibility cleanup items.

### Map-polygon compatibility and ownership (review after port)

Preserve the existing finite-ray containment parity at shared vertices, remote
player behavior when leaving every region, and the host's zero-cache level
handling. Independent contracts make those branches explicit. Correcting them
would change gameplay and should be a separate change. See
[MAP_POLYGONS.md](MAP_POLYGONS.md).

Keep the polygon allocations on the C heap while remaining owners and borrowed
arrays require it. Reset detaches borrowed arrays without freeing them. If editor
metadata allocation fails, the Go constructor returns nil instead of attempting
to free the fixed polygon record in backing storage. This is a narrow, reversible
correction to an invalid cleanup path; forced allocation failure is not tested.

The nearest-vertex draft initially misread a double temporary as an integer. An
independent large-coordinate contract caught it; corrected float/double widths
match the original C captures. The earlier C-baseline report's description of this
temporary was wrong and has been corrected. Expected test outputs were unchanged.
The constructor also retains its initial decimal ID write before copying defaults,
so short default strings preserve the same trailing record bytes.

## Geometry collision prerequisites — qualified, review with conversion

Read the first object's class/flag words as bits in box collision (550F80), matching
object layout and the neighboring circle response. Three numeric float conversions
lost force-suppression bits and failed to clear the wake flag. Independent C
contracts reproduced all three failures; the corresponding second-object paths
already use raw reads. The corrected contracts pass before the Go conversion.

Restore the computed X less-than/equal condition in wall quadrant selection
(550CB0). The decompiled condition byte was initialized to zero and never assigned,
so two masks were unreachable despite the existing X comparisons. Four-quadrant
contracts reproduce the defect; threshold-neighbor contracts keep the original
16.263456 comparison and float width. Both decisions change existing behavior and
are small, reversible corrections under the standing policy. Fresh production
qualification and three-target regression checks pass. See
[WORLD_GEOMETRY.md](WORLD_GEOMETRY.md) for evidence.


### Geometry arithmetic and interface ownership

Preserve the qualified C calculation's observable rounding points. C float locals
can remain wide in x87 registers; their declared types alone are insufficient to
choose Go intermediate widths. The geometry port uses explicit wide calculations
and float32 stores/reloads, checked against the compiled baseline and frozen
captures. In particular, the existing Go ShapeBox.Calc rounds earlier and cannot
replace the selected C corner calculation. Its other callers remain unchanged.

Keep the two collision force coefficients shared while remaining C code reads
them. Move the private direction threshold and circle wall spans into Go, retire
unused C interfaces, and retain twelve exports required by actual C callers. Gate
locked-team notifications use the existing private-message adapter path so hooks
and dispatch ownership remain consistent. See WORLD_GEOMETRY.md for qualification.

## Collision-core ownership and compatibility

Move the Hit indices, activation/angular heads, collision type caches and both
force coefficients to Go once their final C readers move. Keep Hit records in the
existing fixed C-backed 1,024-record class: callback-visible addresses and pool
exhaustion/reset behavior are part of the contract. Ten exports remain for actual
C callers; eighteen obsolete interfaces are retired, including five prior geometry
exports. No C collision implementation remains solely to supply test expectations.

Preserve the existing circle/box interior-distance rule and angular/contact order.
Independent force/list contracts supplement the frozen comparisons; this batch
makes no intentional game-rule correction. Review x87 stores/reloads for explicit
Go arithmetic widths, including the shaft's wide height subtraction before its
float32 store. Keep the retained absolute-value helper's observable scratch write.
Use pointer-typed callback arguments for temporary contact normals so their lifetime
remains visible to cgo. The radial contract exercises the remaining production C
caller across the new double-return export. See COLLISION_CORE.md for qualification.

## World-motion corrected-C prerequisites — review after conversion

Fix sentry removal's membership test to use bit 31: its unsigned `<0` comparison
compiled out the only unlink branch. Real three-object registration/removal
reproduced retained head/middle/tail membership. The corrected list contract
covers permutations, repeated removal and destroyed-object updates.

Treat a nil trigger-script callback result as no admission. An actual one-shot
script accepts the first contact, disables itself, and returns nil on the second;
the previous C dereference crashed. Preserve prior contact state when no result
is returned. Both corrections are small and reversible under the standing policy.
The corrected-C baseline is frozen and pushed as e7174c35 after all-target and
fresh production qualification. Both corrections carry into the Go conversion. See
[WORLD_MOTION.md](WORLD_MOTION.md). Preserve the separate signed allow-team quirk;
values 128–255 do not match an unsigned object team byte.

## World-motion ownership and rounding — review after conversion

Keep the corrected-C sentry unlink and disabled-trigger fixes above. Move decay
and sentry list heads plus velocity type IDs into Go; retain C object allocation
and the shared trace flag while their other C users remain. Six world-motion
exports serve real callbacks/C callers. Remove all 24 obsolete batch interfaces
and nine collision-core interfaces rather than retaining adapters for tests.

Use the compiled 386/SSE2 baseline to place float32 rounding: several declared C
float locals stay wide in x87 registers. Preserve the actual spills and raw float
results instead of mechanically narrowing every intermediate. The trace adapter
still needs a C-backed temporary record for the remaining spatial C callback;
its allocation can disappear with the next connected spatial-targeting batch.
No C algorithm is kept solely as a reference implementation.


## Spatial targeting and opaque callback tokens — review after conversion

Keep the eleven-function spatial batch connected rather than padding its LOC count.
Use Go calls for projectile/cursor candidates and wall normals, retiring obsolete
C interfaces and the temporary C allocations in tracing and aim prediction. Retain
one ray export for its real spell caller. Preserve tile-byte narrowing, compiled
float spills, cursor ties and the existing door-line helper behavior.

Correct the older curve callback bridge exposed by qualification: an opaque integer
token must cross as an integer, while temporary point arguments remain typed
pointers. A deterministic test uses the integer bits of a live Go address containing
pointers, reproduced the original panic, and passes after the correction. This
preserves the callback ABI and curve captures; it does not weaken cgo checking.
The small, reversible correction follows the standing authorization. See
[SPATIAL_TARGETING.md](SPATIAL_TARGETING.md) for full evidence and qualification.


## Monster control ownership and bounded definitions — review after conversion

Move pending-owner pool/list, monster definition list and script-cache cold flag
into Go. Retain mapped cache record layout and two exports required by actual C
callers. Remove 62 obsolete interfaces, including 21 upstream exports whose last
C caller disappeared in this batch; existing fixtures call the same Go owners.
Preserve the animation delay-255 behavior, empty-head word, action layouts, list
quirks and compiled arithmetic. No game-rule change is intended.

Free incomplete definitions on rejected callback/damage/field input and close the
opened file on every exit. Bound names/missiles to 63 bytes plus NUL and general
tokens to 255 plus NUL. Oversized fields stop the record without overwriting its
neighbors; preceding accepted records and the existing load return convention are
preserved. These reversible ownership/input corrections follow standing user
authorization. Independent boundary/allocation-balance contracts supplement the
unchanged original-C captures. See [MONSTER_CONTROL.md](MONSTER_CONTROL.md).


## Quest progress ownership and input limits — review after conversion

Keep the private quest-variable list in Go, preserving its qualified record layout,
serialization order, old-kind retention on updates, and first-occurrence suffix
matching. Three exports remain for actual C reset/save/load callers; 24 obsolete
interfaces and the list-head C global are removed. The namespace and mapped table
state remain compatible with their existing owners. Registered VM builtins invoke
Go directly; separate modern-script API TODOs are outside this behavior-preserving
batch.

Reject overlong names/namespaces instead of overwriting adjacent fields, and reject
truncated save records or excessive counts at EOF. A later malformed record keeps
previous accepted entries; load still clears the old list before rejecting a
version. Independent tests cover every truncation position, size boundaries,
embedded NUL and large counts. These small reversible corrections follow standing
authorization; no defined-input C capture changed.

Boss current HP truncates the wide product to int64 and then uint16; maximum HP
uses a separately rounded float32 product and the original converter. Compiled C
confirms both paths. Preserve byte-wrapped generator caps, signed minion-stage
admission and exact RNG order. See [QUEST_PROGRESS.md](QUEST_PROGRESS.md).


## Prefab prerequisites — review with the map-runtime conversion

Connect the previously unimplemented C section bridge to the existing Go registry.
Set both recognized/error results consistently and retain the caller's unknown-name
object fallback. Match the Go cache-node destructor to the existing C allocation
and removal paths; group references retain their separate tracked allocator. This
is narrower than changing all constructors/callers to a new ownership convention.
Both prior failures were reproduced through real entrypoints and are covered by
independent tests, all three affected sweeps and fresh production qualification.
No C algorithm or frozen expectation changed.

Preserve the existing DebugData reader's permissive short reads in this bridge
repair; consider stricter section validation separately. Audit tile/wall/waypoint
payload lifetime in the larger batch: the current cleanup regression establishes
node-wrapper ownership and explicitly owns its payload allocations. See
[PREFAB_RUNTIME.md](PREFAB_RUNTIME.md).


## Prefab baseline — waypoint ownership and failed-file cleanup

Real prefab-waypoint transfer followed by normal server teardown reproduced an
allocator mismatch. Use the receiving server owner's tracked allocator and release
unplaced waypoint payloads with that same allocator; raw cache wrappers retain
their matching allocator. Independent allocation-balance tests cover both placement
states. Complete prefab-file tests also reproduced unclosed files on bad magic and
unknown object type. Close those two failure paths using the existing helper.
These reversible corrections are qualified before freezing the larger C baseline.

The remaining tile/wall payload leak requires a coordinated ownership repair:
placed secret-wall data retains a back-reference to the cached wall. Defer freeing
that record until transfer and partial-failure behavior are independently tested
in the native work. See [PREFAB_RUNTIME.md](PREFAB_RUNTIME.md).

## Prefab cache payload and special-wall ownership — qualified native correction

Transfer secret-wall data to the actual world wall, updating its coordinates and
back-reference; clear the cache pointer even if later placement fails. Register
breakable walls using the world-wall pointer, preserve their identity, and clear
old secret state when replacing it. Release cached tile/wall payloads and secret
data that was never transferred. Independent real-owner tests reproduced the old
references and leaks before correction; all eighteen frozen captures already match
the initial native translation. See [PREFAB_RUNTIME.md](PREFAB_RUNTIME.md) for
qualification status and explicit undefined-input bounds choices. This is an
authorized reversible correctness change, not exact preservation of the C defects.

Prefab ownership corrections now pass all three targets and fresh production;
77 existing captures remain byte-identical. See the batch native qualification.


### Prefab script prerequisites and reachability

Independent contracts reproduced operand narrowing, unconditional name replacement,
a reserved-name scalar buffer, merge destination/handle lifecycle failures, an
unadvanced pending-object loop, and swapped VM coordinate suffix parsing. Correct
these reversibly before freezing the baseline; retain the C clock bridge's 32-bit
truncation. All three targets and fresh production/integration qualify the repaired
baseline. Remove two selection helpers with no production callers during conversion.
Evidence and review details: [PREFAB_SCRIPTS.md](PREFAB_SCRIPTS.md).


### Prefab script native boundaries and dispatch

Preserve the clock adapter's 32-bit truncation, mixed-sign bounds and actual
builtin-remapping predicate dispatch. Reject incomplete instruction streams and
invalid generation directories; stop a failed source copy before using a stale
backup. Remove two proven orphan helpers and four unused C adapters. Independent
contracts, frozen C comparisons and fresh three-target/production integration
qualify the final change. See [PREFAB_SCRIPTS.md](PREFAB_SCRIPTS.md).

### Voting baseline corrections

Independent contracts reproduced uninitialized padding in both fixed52-byte
name messages and kind 1 withdrawal removing a kind 0 vote for the same player.
Initialize both message temporaries and pass the dispatched kind to the private
withdrawal helper before freezing. Wire lengths/actions/names and external
dispatch interfaces remain unchanged. Keep quest admission settings distinct
from the record's literal 6 minimum; their names alone do not justify changing
that behavior. Three-target/fresh production qualification passes. Details and
fixture limits: [VOTES.md](VOTES.md).

### Voting native boundaries and interface retirement

Bound client name messages to 24 UTF-16 units plus terminator and selection
snapshots to 32 names of 27 units plus terminator. Preserve raw UTF-16 comparisons
and row text. Return without showing the vote window if its local team object is
missing. Clear the allocation-class handle after shutdown so repeated close and
restart are safe. Independent native contracts cover these reversible decisions.

The initial scope audit missed GUI disposal and choice reset; add original-C
contracts from the committed baseline before translating them. Retain choice reset
for the client decoder and move disposal's only caller directly to Go. Cast dispatch
and window-show C returns become void because all actual callers discard their
incidental results. See [VOTES.md](VOTES.md) for evidence and qualification status.


### Console command prerequisites and compatibility

Two independent baseline contracts reproduced incorrect respawn-off flag handling
and a mode-restricted remote-command return leaving its borrowed sender set.
Correct both before freezing; C baseline15df0163 passes all targets and fresh
production/integration. Keep the original numeric narrowing order, C-locale name
comparison, low-byte UTF16 server-name conversion and final 15-byte name limit.

The native registry calls Go handlers directly. Retain one C dispatcher entrypoint
for the quit dialog's observer action, discovered by the first native compile;
move the other private interfaces and scratch strings entirely into Go. Validate
that entrypoint with the actual nil-text calling convention. Preserve defined
`%%` / `%!` formatting and C's ignored string width/precision, verified with an
additional committed-C fixture after review found the first draft's mismatch.
Unsupported argument-taking formats without arguments have no defined C result;
Go's no-argument path treats them as text rather than reading nonexistent arguments.
See [CONSOLE_COMMANDS.md](CONSOLE_COMMANDS.md) for evidence and qualification state.


### Player-state team-count prerequisite

The original active-competitor team loop counted an empty team whenever any
eligible player existed elsewhere. Independent tests with real memberships failed
35 team-mode cases and no non-team cases. Add the same membership predicate used
by the neighboring per-team counter before freezing the C baseline. Preserve
non-team player-record counting, including active players without units, and the
distinct status0x20 rule in the multiple-participants query. This is a reversible
correctness fix for review; evidence and qualification state are in
[PLAYER_STATE.md](PLAYER_STATE.md).


### Player-state native compatibility

Keep the two different eligibility rules (active competitors versus multiple
participants), signed admission limits, wrapping frame subtraction and the exact
status-report mask. Preserve raw UTF-16 name units with C-locale ASCII folding,
fixed equipment-slot capacity and each slot's untouched sixth word. Reuse actual
modifier, minimap and reliable-message owners. Retain only the 11 C entrypoints
needed by the decoder and server lifecycle; two unreferenced scalar helpers retire.
All targets, frozen captures and fresh production qualify. See
[PLAYER_STATE.md](PLAYER_STATE.md).


### Session entry baseline: save discovery and defined settings snapshots

Independent contracts found missing manual saves after gaps, repeated character
counts after empty/rejected metadata, and repeated settings broadcasts from
uninitialized snapshot bytes. Before freezing C, count slots1..13 independently,
clear each character output/check load success, and zero-initialize the GUI
snapshot. Source review also found a negative copy length for path basenames
shorter than four bytes; apply the existing filename branch's zero clamp there.
These reversible corrections preserve normal shipped inputs and are review items.
All three targets and fresh production/headless qualification pass. Details and
failure evidence: [SESSION_ENTRY.md](SESSION_ENTRY.md).

### Item-respawn empty-list removal

Before translating item respawn, add an early return when the C list head is null.
The original sub_4EC6A0 dereferenced head+4 before its later head test. This is a
source-proven invalid access; the new empty/repeated-removal contracts check the
same no-op semantics already used for missing objects in nonempty lists. The
three-line correction is reversible and follows the standing authorization for
confident fixes. See [ITEM_RESPAWN.md](ITEM_RESPAWN.md). Fresh C production
qualification precedes conversion.

### Game-statistics event registration prerequisite

Before freezing statistics C, fix three incorrect references in `sub_425CA0`:
select the actor address using its own host-slot test, and store the target's new
name/host address in the target row. Original C both selected the wrong address
and panicked for a new host paired with a remote player. The independent fixture
uses real connection records with nonzero addresses and checks indices/names/IPs.
The correction also admits the routine's existing absent-target branch for fresh
actors. This is reversible and requires fresh production qualification; no C LOC
change. See [GAME_STATISTICS.md](GAME_STATISTICS.md) for evidence and status.

### Game-statistics reachability and native ownership

Inspecting full caller conditions found nine statistics routines reachable only
through constant-false branches. Retire those routines and disabled callers;
translate the 35 live routines reached from event recording and participation.
Keep the original full C contracts in commit05033129; project the live subset
from those C captures and verify it on C before conversion. Do not preserve
orphan algorithms solely to satisfy tests. The batch workflow now requires this
caller-context audit before fixture design.

Native report records use Go ownership; mapped report arrays use the existing
tracked allocator. Extend each player-name allocation beyond the old ten-byte
minimum when its string needs more room, and avoid the run encoder's unused
final lookahead read. These allocation/access corrections preserve report bytes;
all17 focused captures match C. Keep the legacy field widths, format quirks and
configured-service behavior. Full qualification status is recorded in
[GAME_STATISTICS.md](GAME_STATISTICS.md).

## Map-section prerequisites — review with conversion

Independent original-C contracts exposed missing scratch-wall lookups reusing a
previous wall/region pointer, unattached secret records leaking, the current wall
writer losing its0x80 flag, and v3 floor regions omitting23*Y in both tile-half
coordinate formulas. Correct these before freezing the C baseline. Preserve the
existing wire layouts, signed/count narrowing and sentinel limitations. These
are deliberate behavior corrections, not claims of exact original-C behavior.
They are small and reversible;19 focused roots /4,144 entries pass, and19 captures
/4,126 records repeat identically across processes. All three targets and fresh
production/headless qualification pass. See [MAP_SECTIONS.md](MAP_SECTIONS.md) for failures,
fixture corrections and scope.

## Map-section IO boundaries — preserve existing checksum behavior

The existing cryptfile checksum complements once per ReadWrite operation. Changing
an8-byte coordinate operation to two4-byte operations preserves file bytes but
changes map checksums. The first native fixture caught257 such metadata differences;
restore the original grouping and keep every frozen expectation. All bytes, state,
returns and checksums now match corrected C in all targets, and production/save-load
qualification passes. Do not “simplify” these boundaries without an explicit format
compatibility decision. See [MAP_SECTIONS.md](MAP_SECTIONS.md).

## Arena scoring prerequisite — unteamed killer

An original-C contract confirms a null-team read when an unteamed player kills a
teamed player. Guard the absent killer team's score update while preserving the
player's score increase. This is a small, reversible correction before the next
C baseline, requiring fresh production qualification. Preserve the separate
historical suicide/environment death-counter behavior. See
[PLAYER_DEATH.md](PLAYER_DEATH.md) for evidence and completed qualification.

## Speech-bubble prerequisite — explicit tail removal

The independent original-C remove/append test exposes a stale tail after deleting
the last bubble. Update the tail to the predecessor, matching expiry removal.
This prevents appending through freed storage and is a reversible production
correction to review later. Corrected C and native lifecycle contracts, repeated
three-target captures and fresh production qualification now pass; see
[CHAT_BUBBLES.md](CHAT_BUBBLES.md).

## Combat-overlay prerequisite — literal player names

Independent glyph observations show the original C kill feed turns a player name
containing two percent signs into one. Three paths incorrectly use the name as a
format string. Copy the already bounded player-name field literally instead; the
assist prefix still uses its fixed formatting template. This is a reversible UI
correction made before freezing the C baseline. Corrected qualification is pending;
see [COMBAT_OVERLAYS.md](COMBAT_OVERLAYS.md).

## Combat-overlay prerequisite — missing victim lookup

Consecutive original-C console notifications reuse the prior victim's name when
the next victim is absent or unknown. The independent notification contract
reproduces18 failures. Initialize that local name buffer to empty before lookup,
preserving known-player formatting. This reversible correction is included before
baseline capture; see [COMBAT_OVERLAYS.md](COMBAT_OVERLAYS.md).

### Combat overlays: correct the shared effect link layout before freezing

Independent real-allocation contracts reproduced failed membership and incomplete
cleanup with Go DrawableFX.Next at byte16. Original C attach/detach stores the
per-drawable next link at64; byte16 belongs to movement history. Move Next to64,
retain byte16 as uint32 and assert the link offset at compile time. Remove the
small-pointer cleanup workaround that masked the wrong field. Both caller
contracts pass in effect-layout-corrected; wider qualification remains. This is
a reversible prerequisite repair for review, not an intentional trail algorithm
change. Production C LOC remains25,368/67 files/zero reference.

### Combat overlays: qualified native outcome

C baseline 6998b8ba and its three prerequisite repairs are qualified. The native
32-body conversion passes all 19 focused groups /707 records and all 84 affected
captures /18,691 records unchanged. Three fresh binaries, exact known full-suite
results, gameplay and explicit save/load pass. The first native discovery needed
only a test-adapter signature correction; the first behavioral comparison passed.
Twenty-one private interfaces and eight private C globals retire, leaving eleven
exports for actual C callers. Current C:24,411 lines /66 files /zero reference,
−957. See COMBAT_OVERLAYS.md and its native qualification report.


Client presentation pre-baseline correction: test-created persistent rays filled
all96 pointer slots, then a97th allocation escaped ownership (live count98 ->99
including two endpoints). No code updates the old limit counter. Move the existing
free-slot scan before allocation; full tables now reject creation. Duplicate
removal, first-free-slot reuse and full cleanup pass independent contracts. This
reversible correction is recorded for review and qualified before translation.
Frozen baseline fc842bfa and CLIENT_PRESENTATION.md contain the evidence.


Client presentation translation qualifies with four remaining decoder exports.
The three unused pointer-valued results became void in both headers and exports;
C callers and the test bridge still exercise the actual interfaces. Eighteen
private interfaces and the private chant-tree C global retire. All frozen C
captures passed on the first native build, and all affected targets plus fresh
production/gameplay/save-load qualify. C:23,648 /66 files /zero reference,
−764 from the corrected baseline. See CLIENT_PRESENTATION.md.


Client audio stream baseline corrections (qualified): the override
reader's decompiled local path array provided only36 bytes, aborting on an ordinary
long directory name. Use a separate bounded280-byte path; unsupported lengths fall
back to the packed audio file. Require a complete WAV format chunk and nonzero
channels before using its metadata; otherwise close the override and preserve bag
fallback. Keep the 32-bit C clock narrowing and existing cache/zero-request behavior.
Independent contracts and reproduction logs are in CLIENT_AUDIO_STREAMS.md. These
reversible choices are recorded for later review; native translation is pending.


Client audio streams now qualify natively: all70 bodies translated,22 required C
interfaces retained,48 private interfaces and the driver C global retired. All
frozen captures passed on the first native build. Keep explicit shared layouts,
C-heap ownership,32-bit clock narrowing and existing cache/zero-request behavior.
Affected all-target checks, fresh binaries/ABI, exact known suite and gameplay/
save-load qualify. Current C22400 /66 files /zero reference, −1261 from corrected C.
See CLIENT_AUDIO_STREAMS.md.


Client audio events qualify natively: 62 bodies translated, 21 required exports
retained, 41 private interfaces and two cache/pool C globals retired. All frozen
captures passed on the first native build. Preserve 32-bit clock narrowing, strict
deadline comparison, repeated manager RNG selections, empty-chunk termination,
the initial zero-serial handle convention, and deferred pan updates. The recording
fixture delivers the real device completion notification before voice reuse. No
production algorithms were corrected. Three targets, fresh binaries/ABI, exact
known suite and gameplay/save-load qualify. C:21,083 /65 files /zero reference,
−1,327 from the C baseline (which included ten test-adapter lines). See
[CLIENT_AUDIO_EVENTS.md](CLIENT_AUDIO_EVENTS.md) for evidence and review details.


## Server-browser ownership and deterministic initialization

Re-sort reuses nodes and snapshots the selected 169-byte record before replacing
or freeing it. Keep that snapshot and one 12-byte list sentinel for the browser/
session lifetime, because gameplay reads the endpoint after browser close. This
replaces the original abandoned-list leaks. Zero formerly undefined password
header padding and map-polygon temporary metadata. Preserve defined text, payload,
clock-width, coordinate and high-port behavior; these reversible choices are
recorded for review. Exact C captures and fresh browser/gameplay/save-load qualify.

Native validation caught two translation bugs: temporary static-label text was
freed although the widget retains its pointer, and the scrollbar image offset was
applied to its parent instead of its thumb child. Persistent GUI text and the
correct child target now match the unchanged original-C screens. A real-widget
lifetime regression reproduced the first failure before the fix. See
[SERVER_BROWSER.md](SERVER_BROWSER.md) for contracts and complete evidence.

## Session-dialog MOTD buffer correction

The file/transfer path accepts message lines longer than the display routine's
256-byte split buffer and 256-unit wide-format buffer. Replace these two locals
with input-sized temporary allocations before capturing the port baseline.
Preserve byte widening and line handling; the existing listbox still owns and
truncates each displayed row to255 units. Do not freeze undefined writes beyond
the old buffers. A regression through the actual dialog covers lengths through
4,096 bytes and verifies subsequent lines and unchanged input. This is a
reversible prerequisite under the user's standing authorization; see
[SESSION_DIALOGS.md](SESSION_DIALOGS.md). All-target captures, fresh builds/ABI, exact known suite and filter/gameplay/
save-load qualification pass.


## Session-dialog native ownership and boundary review

Move nine window/file owners to Go while preserving the C heap lifetime of MOTD
file contents. Reject file lengths that cannot fit the signed 32-bit allocation
and handle allocation failure by returning failure; the old wrapped-size/nil
access cases were undefined and are not compatibility goldens. This reversible
choice is recorded for review. Normal file counts, short-read zero tails, byte
widening and ownership match the qualified C baseline.

Disconnect centering reads EndPos at C offsets24/28, including initial resource
position. The shipped origin-zero resource masked an initial SizeVal translation;
corrected Go uses EndPos and an independent three-position regression. Preserve
unrestricted-filter return-before-read ordering and the raw MOTD visibility flags.
All frozen captures and final fresh production scenarios remain unchanged. See
[SESSION_DIALOGS.md](SESSION_DIALOGS.md).


### Client interaction text prerequisites — qualified C baseline

Direct Go console callers can exceed GUI input limits. Before freezing the C
baseline, bound centered text to317 UTF-16 units plus NUL (full text still reaches
the console), allocate the chat-format temporary from input length, and replace
the console formatter's static512-unit unbounded write with bounded local formatting
and a correctly copied variadic retry into allocated storage. Preserve ordinary
text, encoding/byte-count wrap, leading spaces, control words and queue results.
Return failure for impossible allocation sizes or allocation failure. These are
reversible prerequisite corrections to previously undefined out-of-array writes;
review the truncation choice after the port. Existing UTF-16 code-unit semantics
are preserved, including a split surrogate at the centered-text boundary.

Focused independent contracts and the affected corpus pass on all three targets,
including mixed-argument console formatting. Fresh production builds/ABI, exact
known full-suite outcomes, and dialog/gameplay/save-load/inventory scenarios pass.
The shared console helper's localized client-decoding caller remains in scope.
Qualified C13,456 /57 files /zero reference (+14 lines); headers are excluded.
See CLIENT_INTERACTION.md and client-interaction-c-qualification.json. The
centered-text truncation choice remains recorded for later review.

### Client interaction: unreachable dialog and scalar interface

Retire the connection-type dialog's unreachable creation/update/callback graph,
its sole owner, and its constant-false visibility getter. Whole-repository caller
and callback review found no live entry. The former MOTD fixture deliberately
injected that dead owner; remove its gate6 explicitly rather than retaining dead
production state for testing. Preserve all72 surviving rows, unchanged and in
original order, from the80-row C baseline. The native manifest records the exact
projection while the original C capture hashes stay frozen.

Change the retained key-update C argument from uint32_t* to uint32_t. Its sole C
caller passes a numeric byte, and the implementation compares numeric state
without dereferencing it. This preserves the32-bit value and avoids constructing
invalid Go pointers from scalar state. Existing boundary contracts cover the
values. Both decisions are reversible and recorded for later review; neither
changes a reachable game feature. See CLIENT_INTERACTION.md.

### Legacy online session reachability and shrinking mapping inventory

Keep live retry/status/briefing/map-marker behavior, while retiring list/service
branches whose population code is gone. The account selector remains at−1 and all
twelve shipped legacy queue pointers remain zero; preserve their live failure
callbacks rather than inventing unreachable state for tests. Root-visible state
moves to Go, and the final C log formatter retires with its three callers.
C timer captures include28,800 boundary cases and exact signed-frame log order.

The full-suite mapping-reader test required at least1,396 C variables. Removing17
owners crosses that historical threshold without breaking parsing. Replace this
count with actual mapping round-trip preservation and a small independent parser
fixture. This is a reversible test-infrastructure correction, not a change to
production mapping behavior. Other known suite failures remain recorded. See
[ONLINE_SESSION.md](ONLINE_SESSION.md).

## Missing-object script predicates — review with script-binding conversion

Builtins184/185 now return false when object resolution fails. Their C bodies
previously read a parameter that the actual no-argument callback dispatcher never
supplied; both legacy-table and normal VM routes reproduce nonzero missing-object
results. Signatures now match the dispatcher. SetRoamFlag also loses its unpassed
parameter and unused initializer; its live popped-byte behavior is unchanged.
This is an intentional, reversible correction under the standing policy, not
preservation of undefined C behavior. Valid-object bits and higher-priority VM
compatibility overrides are unchanged. See [SCRIPT_BINDINGS.md](SCRIPT_BINDINGS.md)
for independent contracts, original evidence and qualification status.

## Remaining client presentation boundaries — review later

The render-helper conversion removes the focus registry after whole-source symbol
and mapped-offset searches found no callback readers. Preserve the separate live
input-state reset. The pixel-span cleanup had no allocator/population path and is
also removed. Shared live globals keep their existing owners and initial values.

Native winner text can exceed C's 127-unit local buffer; invalid modes do not draw.
A zero-width audio viewport uses centered pan. Debug text is capped to its 80-unit
scratch buffer with a terminator and unchanged neighboring memory. Three independent
native contracts cover these reversible corrections. All 6,129 valid captured C
cases remain unchanged, including World.Max.X pan and signed player levels.
See [CLIENT_RENDER_HELPERS.md](CLIENT_RENDER_HELPERS.md).

## Client resource lifecycle boundaries — review later

Preserve raw-libc ownership for sprite parser data; the Go teardown uses a thin
libc free boundary. The frame sampler retains unsigned 32-bit clock narrowing
before 64-bit history arithmetic. Player color initialization reuses the existing
Go color encoder while copying the mutable white word exactly.

Remove the packed-color and video/material callback registrations after exact
symbol and mapped-offset searches found no readers. Move the main-menu callback
directly to Go, preserving timer underflow, sentinel gates and RNG call order.
Repeated modal creation keeps the old GUI-owned window, matching current behavior.
These are reversible ownership/implementation choices; captured C behavior remains
the qualification contract. See [CLIENT_RESOURCES.md](CLIENT_RESOURCES.md).

## Runtime helpers and orphaned state — review later

Remove the unreferenced bit encoder, modifier-next/material-draw/map-group helpers,
and memory-accounting graph with no production population or output reader. Keep
the console command route and remove its artificial accounting-list fixture.
Remove two inert startup/game-loop calls whose four words have only zero writers.
The old audio-path helper always returns an empty string; retain that public result.

Preserve the material lookup before its class gate, including null-cache equality;
ASCII-only equipment-name folding in the hosted C locale; signed modifier-key
bytes; clock narrowing to 32 bits; and child traversal while ownership changes
reverse the destination list. Keep registered no-op modifier callback identities:
two are compared explicitly, so empty bodies do not imply dead registrations.
Two shared definitions move unchanged, including the four-byte settings flag.

The baseline and first native probe agree on all 55,986 captured cases. Default/server/highres pass 304/303/304
affected tests; final production qualification
is recorded in [SERVER_RUNTIME.md](SERVER_RUNTIME.md). No C converter or other
algorithm remains solely for testing; x87 control-word instrumentation remains.


## Extension/name and listing baseline corrections

Before freezing this C baseline, replace the player-name parsers' one-byte local
buffer with the actual 28-unit field capacity plus terminator; retain low-byte
UTF-16 conversion and exact first-active matching. Missing team/output owners
return failure without partial writes. Initialize the server-listing header,
keep its eight-byte map payload plus terminator, and reject short slices at the
Go boundary. These reversible behavior corrections follow the authorized process;
review them as corrections, not claims of preserving undefined stack contents.

Five groups /21,135 cases pass repeated default, server and highres captures and
fresh headless C gameplay. Native production/ABI, known-suite and save/load are
still pending. Inventory captures hash every complete normalized result per case
because verbose fixtures otherwise repeat hundreds of MB of unchanged memory;
independent behavior assertions remain. See [SERVER_TEXT.md](SERVER_TEXT.md).


The extension/listing native conversion is qualified: 173 affected roots in each
profile and all 21,135 cases, production/ABI, known-suite comparison and headless
gameplay/save-load. Remove the eight unreachable memory-file C methods while
keeping their live ABI type; consolidate 22 shared definitions without changing
storage types. Keep the seven registered no-op callback identities. Define flag
indices outside 0..31 as zero (the UI uses 1..5; the old negative extreme domain
could divide by zero). No C algorithm is retained for these tests. See
[the qualification](server-text-native-qualification.json).


Formatting/scalar/audio-directory conversion is qualified: 375/373/375 affected roots across
default/server/highres, 668,876 frozen cases, production/ABI, exact known-suite comparison and
headless gameplay/save-load. Preserve the custom formatter and literal processing;
use Go-owned temporary text and terminate fixed buffers at actual capacity rather
than reproduce out-of-bounds C writes. Port only live base10 decimal semantics and
exact normalized directory presence; retire unreachable glob/string APIs. Keep
shared catalog storage unchanged. Tagged raw words become pointers only for string
arguments. See [TEXT_FORMAT.md](TEXT_FORMAT.md) for evidence and delegation lessons.


Retire unused memory lookup and GUI C adapters after whole-source reachability
review; reuse the existing Go registry for the sole tagged threshold getter.
Keep offline symbol-name metadata and live native window APIs. Preserve safe-mode
allocator/ASan behavior; repair its stale generic free call before the removal.
This qualifies 117 fewer C lines and 24 fewer interfaces. The optional safe build
passes; runtime qualification remains default/server/highres. Continue one bounded
Luna helper after the two-batch review, with primary evidence checks and integration.
See [ORPHAN_BRIDGES.md](ORPHAN_BRIDGES.md).


## Exact integer damage dispatch

Add a separate int32 registry alongside the established boolean damage API.
Register canonical exact helpers in both maps, deriving bool only from the exact
result. Boolean-only registrations leave the value map alone and retain original
C fallback for raw callers. This preserves results such as256 for both full-word
and low-byte consumers. Raw calls require a configured callback as before. Scope:
five raw sites, eleven canonical registrations; no layout/C-address changes.
Original132-root, three-profile baseline and2071 frozen owner cases and conversion
are fully qualified, including fresh production and scenarios. Reviewable/reversible API decision authorized by standing user
instructions. See DAMAGE_VALUES.md.


## Old reproducible compiler-cache eviction

During exact damage-value qualification, disk space remained tight after archiving
older qualified binaries. Marker-specific cache audits found no more candidates.
Pruned40 older repository Go compiler archives (1,882,456,060 bytes), using a
reviewed pre20:00 cutoff and independent hash/stat/archive/module/host-use checks.
This policy permits rebuilding older versions rather than retaining every compiler
cache entry; it does not delete source, module downloads, original assets, goldens,
logs, binaries or qualification reports. No build was active during pruning.
The bounded audit came from Luna; the primary reviewed and executed the cleanup.


## Batch all registered object updates together

The next registry batch covers all53 object updates through one CallUpdate boundary,
rather than completing another tiny adjacent damage-sound batch. Bind existing Go
export wrappers directly, preserving their complete bodies, names/C-addresses,
data sizes and parser order. Six mutable handlers remain looked up at invocation.
Original-owner sibling tests discard return values because CallUpdate is void;
original direct-return tests remain intact. Coverage uses actual registrations
and explicit execution counts, with the existing registered UndeadKiller owner
covering the53rd name. Baseline and conversion281-root qualification are complete in three profiles; fresh production and headless scenarios pass.
See UPDATE_REGISTRY.md. This reversible API/process choice follows standing user
authorization for larger coherent batches with independent contracts.


## Preserve both initializer calling conventions

Creation uses typed Go registry dispatch after normal object-template copying.
Initialization uses a direct Go handler map with explicit raw fallbacks: existing
CallInit keeps its one-argument/nil behavior, while CallInitWithArg preserves the
two-argument call used by pending activation, player arrival and item respawning.
The latter keeps the configured-slot precondition and original caller guards.
The respawn call loads Init by offset 688; named-field searches alone missed it.
Monster/Shopkeeper share a C callback address but retain distinct data sizes, and
both replaceable monster handlers remain resolved at invocation. All 208 roots per
profile and fresh production/scenarios pass. This reversible design follows the
standing authorization; see LIFECYCLE_REGISTRY.md for evidence and fixture corrections.


## Preserve collision argument and return conventions

Collision registration binds complete existing export wrappers with callback-specific
argument conversions. The public integer-word CallCollide API remains; queue/trace
owners use CallCollideWith with real object/normal pointers, explicit raw fallback
and KeepAlive. Compiler escape diagnostics and stack-growth/GC contracts qualify
the conversion. The UndeadKiller third word remains a nonzero condition; normal
consumers keep their existing pointer interpretations. Three default aliases keep
one callback identity and distinct per-name data sizes.

The temporaryMagicMissile owner consumes an integer return and stays raw rather
than receiving an invented result from a void API. Two damage-sound callbacks are
also deferred for their own nil-default contract. These reversible scope decisions
follow standing authorization. See COLLISION_REGISTRY.md for all qualification.


## Preserve distinct transfer and sound default contracts

Transfer registration dispatch calls the complete existing wrappers. The server
API keeps exact nonzero success and generic failure text; the outer DefaultXfer
shortcut continues calling its local implementation with detailed errors. It does
not gain the legacy wrapper's replaceable-handler semantics. The damage owner
keeps its fixed nil-slot default wrapper rather than resolving the mutable
DefaultDamageSound registry pointer. Registered sound returns remain ignored.

Explicit Go maps retain raw pointer-call fallbacks and configured-slot
preconditions. Callback identities and layouts are unchanged. These reversible
choices preserve the original dispatch contracts; qualification is recorded in
XFER_SOUND_REGISTRY.md.

## Internal glue before external bindings

The immediate phase removes engine-owned C glue while retaining SDL2/OpenGL/OpenAL
and similar external bindings. Whole-build cgo-free clients are deferred for a
separate discussion. Selected dependency metadata, rather than standalone C LOC,
now guides removal order. The leaf cleanup batches 33 unused imports with the
equivalent Linux socket ioctl constant; Go bodies and frozen expectations stay
unchanged. This avoids an extra qualification cycle for trivial import-only work.
See INTERNAL_C_GLUE.md and CGO_LEAVES.md.

## Libc memory comparison ordering

The Go memory-helper batch preserves comparison sign, not incidental libc return
magnitude. Original Linux/386 strcmp magnitudes vary by optimized path and
mismatch position; memcmp also differs by span length. Whole-source caller review
found no engine magnitude consumer. Raw original captures remain in baseline
commit `5cc27785`; the original implementations were restored before capturing
normalized ordering three times and freezing the revised contracts. This is a
deliberate reversible compatibility choice, not an unchanged exact-return claim.
See GO_MEMORY.md and go-memory-c-qualification.json.

## Go memory-helper performance boundary

Keep Go copy/clear and a doubling-copy nonzero fill. The production direct callers
use 60-byte copies and 32-byte zero fills, both faster in the measured VM; larger
Go copies regress in the bounded 386 microbenchmarks. No current large production
consumer was found. Record the limitation and revisit if a real workload warrants
a specialized implementation; do not claim a whole-game speedup. See GO_MEMORY.md.

## Preserve allocation domains during centralization

Use build-tag-specific legacy allocation helpers: normal calls remain raw libc,
while `safe` calls use the existing tracked allocator. The runtime `NOX_SAFE`
setting does not select the old C macros, so it must not select this route either.
Centralize tracked backend calls in the same raw libc file without changing
failure bookkeeping, pointer layouts or ownership. Defer known string-lifetime
issues and failed-realloc bookkeeping to explicit owner changes; this batch
preserves their behavior. The allocator remains native storage until a separate
qualified backend replacement. See RAW_ALLOCATION.md.

## Retire unused internal string interfaces

Remove the unreferenced input string-buffer get/free pair, filesystem normalization
export, CStringArray, CBytes wrapper and CWStrSlice. Repository-wide source/header,
registration and build/documentation review found no runtime consumers or supported
external library API. Preserve live string allocation and callback entrypoints.
Retaining dead wrappers solely to fix/test their ownership quirks would obstruct
internal glue removal. This reversible internal API cleanup is recorded in
STRING_BOUNDARY.md, including the limits of external-consumer discovery.

## Retire proven-unused internal exports

Remove 65 forwarding exports and 63 prototypes after whole-source reachability
review including C preambles. Preserve live logger/sentinel declarations beside
retired wrappers. Executable builds and remaining ABI callbacks pass; no supported
external shared-library consumer was found. This reversible retirement preserves
all live implementations and tests. See UNUSED_EXPORTS.md.

## Complete regression selection and valid pointer fixtures

Use `^Test` for the full root porttest corpus, replacing a stale manual selector
that omitted 899 default/highres roots and 892 server roots. Keep focused owner
patterns separately. The expanded server sweep exposed an invalid integer seed
in the player-reset fixture's pointer-only Obj130 field. Seed live objects there
and in the other pointer fields, preserving numeric seeds, byte assertions and
frozen output; retain the engine's typed pointer clear. A focused original run did
not reproduce the heap-state-dependent crash, so full qualification remains the
gate. See COMPLETE_PORT_CORPUS.md for completed versus pending evidence.

## Bounded concurrent root-profile sweeps

Try at most two prebuilt root porttest processes after sequential compilation.
Require matching binary/source/profile records, isolated outputs, unset optional
capture/diagnostic paths, and joined jobs before source changes. The bounded
96-root/profile trial passes; it does not establish full-corpus isolation or a
measured general speedup. Keep production/headless qualification sequential and
fall back to one process if resource/isolation failures arise. See PREBUILT_PROFILES.md.

## Remaining unused exports and complete concurrent corpus

Retire 378 additional wrappers and 333 simple header declarations after exact
whole-repository reachability review and independent deletion/AST checks. Preserve
four ambiguous candidates for separate review rather than expand this batch.
Remove only the two newly unnecessary C imports and ordinary unused imports;
all live Go declarations, fixtures and callback identities remain unchanged.
See REMAINING_UNUSED_EXPORTS.md for full qualification and counts.

The first complete root corpus passes with two prebuilt processes and sequential
compilation. Retain that bounded setting for this audited corpus, with exact test
name sets and source/binary/profile verification. Individual runtimes stayed near
the preceding sequential sweeps; the overlap reduces elapsed sweep time without
establishing a general engine performance improvement. Keep production and
headless qualification sequential.

## Retire C interfaces while preserving Go callers

The 265-export batch deletes 126 unreferenced wrappers but keeps 139 functions
used by Go, removing only their C exports. Resolve owning-package references
before deletion: same-named functions in another package are different owners.
Remove 244 exact header prototypes and 15 unused C imports after preamble/header
review; all complete regression, production, ABI and gameplay gates pass. See
GO_ONLY_EXPORTS.md.

Reject a separate follow-on object-state export proposal whose inventory missed
native calls inside a test preamble. C preambles are active code despite Go's
comment syntax. Narrow helper work to explicit edit spans and inspect whole C
selectors: the subsequent scalar draft initially changed the suffix but retained
C prefixes. Primary corrected it while reconstructing the draft; original storage contracts
and fresh qualification remain the conversion gates.

## Safe scalar contracts and the raw-blob fixture boundary

Qualify all 395 scalar owners under safe in addition to default/server/highres.
The optional original safe attempt at the separate raw-storage fixture correctly
tripped the runtime mapped-address guard: that fixture deliberately probes
registered offsets whose owners were extracted to Go. Preserve its original
failure evidence, guard and expectations; retain its raw-alias coverage in normal
profiles and run the scalar contract alone under safe. The latter passes the
existing frozen capture. This is an explicit fixture/profile boundary, not a
conversion regression or a claim of complete safe runtime qualification. See
GO_SCALAR_STORAGE.md and go-scalar-storage-baseline.json.

## Native internal calls and replaceable server hooks

Use existing native Go owners or exact extracted bodies for the 31 remaining
primitive-adapter caller files. Preserve exported C entrypoints while other
callers still need them. All 50 source edits qualified without fixture or golden
changes; selected production cgo files fall from 302 to 271. See
[GO_NATIVE_CALL_BOUNDARIES.md](GO_NATIVE_CALL_BOUNDARIES.md).

Preserve evaluation order when replacing a wrapper with a method call. The
replaceable `GetServer` hook can run before argument reads in a direct method
expression. Keep an exact native helper or an explicit temporary where needed;
string interning and subsequent reading must retain their original order too.
This is a compatibility constraint, not a new game behavior.


## Native geometry/state boundaries and private wrappers

Retire unused private Go wrappers when complete source/reference review proves
no runtime caller or registration, even if their bodies contain pointer/layout
operations. Keep historical translation-rule strings and same-named owners in
other packages distinct from runtime legacy calls. Migrate all callers before
retiring a live adapter. This connected batch retires 145 wrappers and 24 C
imports without changing any actual C export or callback identity.

Use target layout probes and complete field reader/writer audits for C scalar
and point conversions; preserve raw word aliases and original ownership. Primary
caught stale field selectors in a Luna scratch-array conversion before install.
All frozen storage, complete root, production and gameplay qualification passes.
See [GO_LAYOUT_BOUNDARIES.md](GO_LAYOUT_BOUNDARIES.md).


## Native record storage and fixture boundaries

Use existing native window/render/particle/list/inventory records where target
layout probes and alias contracts agree. Represent 44-byte world grid cells as
11 raw words, preserving pointer-word slots, signed fallback decoding and
unmanaged allocation/free ownership. Preserve the browser's existing 172-byte
Go record versus 169-byte packed copy; no packing correction is bundled here.

Keep independent C fixture allocation/observation where it still tests live
boundaries. Four draft type mismatches were corrected before installation. The
first compile separately caught an import-cleanup error: package names can differ
from directory basenames (`common/flags` is `noxflags`). Restore the import and
record failed logs; resolve package names before future pruning. Frozen captures
and root assertions remain unchanged. All completed-batch gates pass.
See [GO_NATIVE_RECORD_STORAGE.md](GO_NATIVE_RECORD_STORAGE.md).


## Native owner calls and compiled constants

Preserve legacy package-wide compiled constants even when root Go profile
constants differ: NOX_HIGH_RES applies to every legacy profile. Move private
callers directly to existing owners with explicit signed narrowing and the same
string lifetimes. Keep the public browser record raw-address accessor with an
unsafe.Pointer result, documenting the C-specific return type removal.

Use the audited affected-family selection for this localized caller batch; expand
helper proposals to include actual AI policy and message-drawing routes. All
three focused profiles, safe/static, fresh production/ABI, known-suite and headless
save/load gates pass without changing frozen expectations. See
[GO_NATIVE_OWNER_CONSTANTS.md](GO_NATIVE_OWNER_CONSTANTS.md).


## Object-state fixture identities and exports

Retire the 37 object-state exports used only by fixture dispatch/address maps;
keep four live collision callbacks. Before removal, two instrumented original
processes covered all 44 operations with zero consumed function-address IDs.
Audit shared map aliases and indirect raw-word normalization, not only the local
fixture normalizer. Reserve the 37 removed entries in reservedFunctionIDs so
map-size-dependent generated IDs keep their existing formula. Original raw map
sizes vary at three loot snapshots; frozen outputs and generated-ID counters
match and remain the acceptance boundary. Do not regenerate expectations to
hide a difference. All conversion gates pass. See
[OBJECT_STATE_OWNERS.md](OBJECT_STATE_OWNERS.md).


## Equipment, inventory and resource fixture routes

Retire 61 wrappers whose remaining routes can use existing Go owners, preserving
22 inventory callbacks and the gold-pickup callback still used by C. Keep C-int
narrowing, signed short results, unsigned gold and raw 32-bit return words exactly
as before. The retired functions have no fixture address registrations, so no
identity reservation change is needed. Independent C modifier/drop/death observers
remain. All three profiles and production gates pass without changing expectations.
See [INVENTORY_RESOURCE_OWNERS.md](INVENTORY_RESOURCE_OWNERS.md).


## Shop/trade transport-only adapters

Remove the obsolete shop price pointer-through-float transport and code/count
pointer roundtrips when calling existing Go owners. Those adapters do no float
arithmetic or dereferencing; native pointers and raw uint32 values preserve the
386 input words. Keep signed input narrowing and nil session handling. The native
public cancel API and independent trade-pickup observer remain. The same frozen
33-root profile selection and all production gates pass. See
[SHOP_TRADE_OWNERS.md](SHOP_TRADE_OWNERS.md).


## Spell/reward fixture addresses

Use raw addresses when probing function-map use: reward's normalized IDs overlap
the generated-ID range, so an ID-only filter can misattribute hits. Two original
processes show no consumed candidate function addresses; preserve their cardinality
with 13 spell and 12 reward reservations, keep nine live reward callbacks and their C test
routes. Original raw map sizes vary in three records while frozen outputs and
generated-ID counts agree. All converted gates pass without expectation changes.
See [SPELL_REWARD_OWNERS.md](SPELL_REWARD_OWNERS.md).


## Map room/painting fixture dispatch

Retire 64 numeric fixture-only C exports without changing owner algorithms or
normalization. Keep x87 control-word fixtures, allocation helpers and live painting
transfer addresses. Match sparse operation IDs and each numeric transport width;
mapRoomIsHall returns a uint32 despite its C pointer argument, while mapRoomOverlap
returns a pointer. Primary corrected both argument/result distinctions before
compilation. Pass the existing grid-config pointer directly to native Go. All
89 roots in three profiles and production gates pass with unchanged expectations.
See [MAP_ROOM_PAINT_OWNERS.md](MAP_ROOM_PAINT_OWNERS.md).


## Catalog/effect adapter retirement

Retire eleven fixture-only production bridges and seven private C-typed wrappers;
keep all live rendering/update callback routes and their C fixture coverage.
The curve's specialized C body only transported calls to a Go test observer, so
replace it with a local observer preserving ordered segments, opaque int32 tokens,
input guards and frozen expectations. No independent C reference algorithm is
lost. Add rain-orb creation boundary/failure/RNG contracts before conversion.
Explicit int32 conversion preserves map-cycle C-int results. No candidate has an
address registration, so identity reservations need no adjustment. All gates pass.
See [CATALOG_EFFECT_OWNERS.md](CATALOG_EFFECT_OWNERS.md).


## Native collision identities and result dispatch

Use distinct nonzero-sized static Go slots for all 51 registered collision keys,
retaining 53 names, aliases, data sizes and late-bound handlers. The existing object
layout remains 32-bit. One registry handles void and optional uint32 results;
magic-missile expiry preserves declared native results and exact raw C fallback.
Native void owners explicitly return zero when asked for an integer result: the
previous mismatched C invocation had no defined integer result, and production
update dispatch discards the enclosing return. Activation retains the low byte of
the current key, not a linker-specific address value. These are reversible choices.
Pin fixture pointer arguments across uintptr transport and nested Go calls.
Chest fixtures deliberately reuse a collision owner as Death; register that
fixture-only typed route explicitly rather than executing a native data address.
See [COLLISION_IDENTITIES.md](COLLISION_IDENTITIES.md) for qualification and scope.


## Native death identities and equipment text

Retire fourteen death C callback identities using distinct static Go slots and the
existing typed cache. Keep all names, closures, data sizes, parser hooks and raw
fallback behavior. Direct fixtures preserve ImpEgg's uint32 bits, Spawn's signed
short extension and GameBall's C-int result. Native uint16 pointers and the exact
underlying string-interning helper remove the equipment death owner's remaining
C type dependency without changing formatting or allocation behavior. Existing
Glyph override coverage is reused; the new contract covers object/type storage
and dispatch after GC. See [DEATH_IDENTITIES.md](DEATH_IDENTITIES.md).


## Native creation/init identities

Retire nineteen C callback addresses across twenty names, retaining the intentional
MonsterInit/Shopkeeper alias and distinct per-name sizes. Creation identity remains
on ObjectType; Init is copied to Object. Reuse existing typed caches and raw fallback.
Keep mutable hooks inside per-call closures. Extract five creation owners and
Boulder initialization without changing write order or return semantics. Direct
fixtures preserve sparse IDs, pointer/integer result bits and PlayerInit's signed
byte extension. Gold accessor classification now follows the native registered
key and has explicit signed/wraparound contracts. See
[CREATE_INIT_IDENTITIES.md](CREATE_INIT_IDENTITIES.md).


## Native damage/sound identities and direct owners

Retire thirteen registered C callback identities and fifteen adjacent wrappers
from the same fixture cohort, removing 28 exports in one qualification cycle.
Keep distinct Stone/Default keys, both damage result caches, late-bound sound
hooks and the fixed default sound owner for a nil callback slot. Preserve fixture
operation widths and first-NUL semantics. Reserve fifteen retired address entries
because map cardinality affects subsequent capture IDs; skip nil entries rather
than identifying them. Independent C observers remain. See
[DAMAGE_IDENTITIES.md](DAMAGE_IDENTITIES.md).
