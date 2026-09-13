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
