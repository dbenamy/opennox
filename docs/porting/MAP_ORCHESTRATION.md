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
| step-boundaries | 24 | 635b72c0db8747b78fa33790d65a35c0c45e5886ce3f3ca7632a29fbbfd5de32 |

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

## Expanded C contracts

Empty, one-byte and 63-byte names pass. The 24 additional real step cases cover
normal/backdrop painting and impossible required-decoration failures over eight
seeds. All 16 valid maps return1; eight impossible cases return2. Backdrop maps
contain97 painted cells versus15 without a backdrop for each paired seed.

Fifteen outer contracts cover main/alternate success, missing-theme abort and both
100-attempt retry limits using real generation. They check metadata clearing,
flags and saved objects during attempts, object restoration, serializer path/flags,
service-call order and actual nc.obj/blend.obj contents. File cases cover absence,
current/blend/both files, nonempty-directory removal failure, alternate missing
source and map-save failure. The serializer alone has a scoped tagged wrapper with
real fallback; generation, filesystem operations and object ownership remain real.
Metadata and map-switch callbacks are recorded external service boundaries.

All175 cases/four mandatory captures and the new contracts pass three repetitions
(root16.382s, legacy0.046s). Qualification of all22,658 map cases/281 groups plus
contracts passes default/server/highres in74.179s/148.546s/82.656s, with every
historical hash unchanged. No orchestration is converted.

The expanded C baseline is ready to commit/push before conversion. The caller
audit finds no C consumers outside this group. Existing Go debug/gameloop wrappers
can call native helpers directly; retire all five C declarations and the unused
alternate-start reference in `keep.go`. Audit floating-point rounding in radius
and backdrop coordinates. Qualify accumulated variants, production builds,
known full-suite failures and fresh unchanged gameplay; record the new C LOC.

## Unflagged backdrop prerequisite

The first expanded C integration failed while painting a nil backdrop decoration.
New backdrops have ThemeFlags zero; the inherited selector used the room address
as a random weight for that case. The correction sums matching unrestricted
weights for flags0, leaving all six normal phase branches unchanged. Direct tests
cover eligibility filtering, both weighted choices and empty/ineligible lists.
The earlier room port faithfully preserved the C defect; this correction is
separate from orchestration translation and has no C LOC delta.

Evidence: c-expanded.log (original failure), c-outer.log (15 original outer
contracts pass in1.075s), step-boundaries-audit.json, c-expanded-locked-repeat.log,
and c-expanded-qualification.log under build/port-map-orchestration. Qualified
captures are compressed with verified SHA-256 manifests.

## Separate behavior reviews

Current C failure paths can leave the generation flag set and, for main start,
leave the original object list detached in saved storage. Preserve and test that
behavior for equivalence; any cleanup correction should be separate.

Existing theme cleanup leaves nested wall/floor records allocated. The fixture
records those survivors and releases them at teardown. Deep cleanup requires an
ownership/sharing review; current tests do not claim full production cleanup.

Outer fixtures and step boundaries are applied and qualified. Native step/outer
helpers and prepare-native.py under build/port-map-orchestration remain unapplied.
Final application qualification scripts are in build/port-map-orchestration-native.
