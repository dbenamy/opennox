# Procedural-map orchestration

The next connected scope is five remaining routines in `legacy/GAME3_2.c`,
207 section lines: map-name setter/getter, main start, generation step and
alternate start. **All five remain C.** The wall-list prerequisite is documented
in [WALL_LIST.md](WALL_LIST.md); growth is already native and qualified.

## Original C step baseline

The fixture uses real keyed synthetic themes, real generation and real painting,
object and waypoint owners. Available game assets do not include generator themes.
The unchanged gameplay scenario therefore complements these synthetic integration
cases; it does not establish procedural-map coverage by itself.

Current mandatory capture groups:

| Group | Cases | SHA-256 |
| --- | ---: | --- |
| normal-maps | 128 | c6fc4af3bd0c2935bf7c5c3607302ca0e066950237e2d455f5170a42967ef6ab |
| ring-maps | 16 | d8f1e06ff469e2c6b6dce2f9583c375a0c62f936831875547d06c3803acfb93d |
| theme-failures | 7 | e0bdd5828f0ebcc7a584699b059b6826c8b19281160eac5351eb9600c4da89c6 |

Normal cases vary map size, room size, recursion, branching and seed. Ring cases
vary recursion and seed. Failure cases include missing, empty and malformed themes.
Captures retain room-release contents, generated floor/wall/object/waypoint state,
globals, return values, guard state, floating-point control and RNG tails. The
positive probe requires a PlayerStart and generated room releases. Eight varied
allocation layouts each compare complete captures for a normal map and two rings.

Release snapshots live in retained Go buffers that are never passed to C. They
are marked `captureOnly` so engine-pointer normalization cannot mistake ordinary
values for diagnostic pointers. Grid-row IDs are retained before release. See the
wall prerequisite document for the audited tracing corrections and repetitions.

## Remaining baseline work before conversion

Add empty, one-byte and 63-byte map-name contracts. Extend real step captures for
backdrop painting and impossible required-decoration failures. Cover main and
alternate start using real success, missing-theme abort and the 100-attempt retry
limit. Count attempts through the existing progress callback; do not substitute
the generator step. Check metadata clearing, generation flags and object-list
ownership during attempts as well as final restoration.

Exercise absent/current/blend/both files, a nonempty-directory removal failure,
alternate-start missing-source failure and map-save failure. Verify file contents
and service-call order. Map serialization may use a tagged wrapper with fallback
to the real serializer outside these contracts; record its path/flags/result.
Metadata and map-switch callbacks are external service boundaries. Generation,
file operations and object-list detach/restore remain real.

Repeat and lock the expanded C baseline, commit/push it, then convert. The caller
audit finds no C consumers outside this group. Existing Go debug/gameloop wrappers
can call native helpers directly; retire all five C declarations and the unused
alternate-start reference in `keep.go`. Audit floating-point rounding in radius
and backdrop coordinates. Qualify accumulated variants, production builds,
known full-suite failures and fresh unchanged gameplay; record the new C LOC.

## Separate behavior reviews

Current C failure paths can leave the generation flag set and, for main start,
leave the original object list detached in saved storage. Preserve and test that
behavior for equivalence; any cleanup correction should be separate.

Existing theme cleanup leaves nested wall/floor records allocated. The fixture
records those survivors and releases them at teardown. Deep cleanup requires an
ownership/sharing review; current tests do not claim full production cleanup.

Local drafts under `build/port-map-orchestration` are unapplied: outer service and
contract fixtures, step boundaries, native step/outer helpers and preparation
scripts. They are convenience drafts, not a qualified baseline or conversion.
