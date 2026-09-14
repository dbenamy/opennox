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
