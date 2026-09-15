# Sprite animation and draw-data parsing

## Scope and status

The corrected C baseline is fully qualified. The production conversion has not
yet been applied. Local evidence/scripts: `build/port-client-sprite-animation/`.

The batch covers six complete C files, 640 physical lines: `animdraw` (234),
`canidraw` (131), `staticdraw` (55), draw/parse/parse (137), `boulderdraw` (59),
and `slavedraw` (24). It contains seventeen C routines. Sixteen move to Go or
reuse an existing Go helper; the unused C vector-array loader is removed because
production already uses `src/drawable_vector.go`. Thirteen C callbacks/parser
entry points remain for current callers. Three private helper entry points and
the unused vector-array entry point can retire.

Production C before conversion: **94,415 lines /124 files /zero reference C**.
Expected after conversion: **93,775 lines /118 files /zero reference C**.

The shared `nox_xxx_drawObject_4C4770_draw` renderer and draw-data cleanup stay in
production C. The new fixture establishes real image rendering for ordinary
non-player sprites with fixed lighting; it does not qualify the shared renderer's
player/team, invisibility, or wall-clipping branches for conversion.

## Prerequisite correction

Conditional Random selected from the inclusive interval `[0, frameCount]`, but
its parser allocates exactly `frameCount` entries. Correct the bound to
`frameCount - 1`, as ordinary animation already does. The independent contract
uses an owned sentinel after a one-frame array, exercises both active states over
128 seeds, verifies the real callback never selects the sentinel, and checks
exactly one RNG value is consumed. It also verifies these seeds distinguish the
old upper bound. This is a deliberate behavior correction, not a new C oracle
for out-of-bounds reads. The change does not alter C LOC.

## Fixture and baseline

The tagged image owner installs valid synthetic type-6 images into the actual
`RenderSprites` index and handle maps. Normal `Image.Pixdata` interns their data.
It preserves previous maps and frees only its images; the production raw-image
C-handle guard remains intact. An isolated handle arena restores prior state.
The image contract checks real pixels, distinct image results, retained drawable
image metadata, and distinct static-random/slave callback identities.

The existing effects fixture supplies real drawable allocation, lists, spatial
index, RNG, viewport, renderer and framebuffer. Its snapshot now normalizes only
explicitly registered image/data pointers, in addition to its existing owned
list/callback pointers. Other unexpected pointers still fail. Old fixtures do
not register the new pointer classes. The accumulated C suite checks that this
shared fixture extension preserves their established oracles.

Four result groups are captured, inspected and repeated before native replacement:

| Group | Results | SHA-256 |
| --- | ---: | --- |
| Animation frames | 16,096 | `648484bf7f0f989ee84d59ff9e38f35a8251a8fefc5b3d5cbf9ddb4f574e8d55` |
| Static/conditional/basic parsers | 675 | `fe24ac26e809b9b3544791850d802ac5011f588d3e78bbb8b674e82972eaf91c` |
| Vector/state readers | 315 | `cad133b687a342d021a0316636a25ffc6082f6151bc11ce3a6008b90c831d661` |
| Boulder frames | 520 | `a5aa6343cc7b85d21c13caf0d2bf49d6f838ab47e94480bb5fd8404bd1647a0f` |

Total: **17,606 results /four groups**, plus image and random-bound contracts.
Captures include pixels, render words, normalized drawable state, deletions,
RNG positions, relevant last-draw globals, parser metadata, frame arrays, scratch
bytes and remaining input. Parser tests assert exact consumption and callbacks;
state parsing also checks its size header and installed kind.

Coverage includes all basic animation modes and unsupported modes; active/inactive
conditional states with distinct images and delays; frame-counter wrap; counts
1/2/7/31/255; delays 0/1/7/255; one-shot/fade deletion; class/game-flag suppression;
manual indices; edge clipping; boulder movement below, at and above squared distance 100, and directional
frame wrap. Parsers additionally exercise zero frames without drawing them,
zero through five conditional states, all named animation kinds and unknown
strings, indexed images and the existing missing-external-image result, vector
headers and all three state slots. Impossible zero-frame loop/random drawing and
invalid manual indices are outside the fixture's lifecycle invariants.

## Compatibility and review notes

- Fade animation's early deletion leaves alpha enabled in the original callback.
  Preserve this in the port; review any visual correction separately.
- Conditional draw-data reports size 16 despite its 56-byte allocation. Preserve
  the header, byte counts/delays and five frame-array layout for C consumers.
- Keep parsed arrays on the C heap while the production C cleanup owns them,
  including the existing zero-count allocation convention. Allocation-failure
  cleanup and malformed-asset policy are not redesigned by this batch.
- The unused C vector-array loader is not kept solely for tests. Existing Go
  vector loading remains its production implementation.
- Door locks, arrow tails, glyph alpha and summon effects need their own connected
  owners/contracts. They are subsequent work, not implicitly qualified by this
  non-player animation fixture.

## Validation plan

Repeat focused C captures; run accumulated C regression tests because the shared
fixture changed, plus affected server/highres variants; commit/push the baseline.
Then replace C, require unchanged captures/contracts, run accumulated and affected
variant regressions, build all three production binaries, audit ELF32/SSE2 and
symbols, compare full asset-suite failures to the established oracle, and run a
fresh unchanged headless gameplay scenario. No C algorithm is retained for tests.

## Qualified C checkpoint

All 17,606 results repeat byte-for-byte (c-e/c-f). Six focused tests pass, including
both independent contracts. Accumulated standard: 672 selected/completed tests
(671 pass, one optional prerequisite skip), 362.349s. Affected server and
highres each pass 39 tests (164.168s/32.001s). The corrected C production
build passes and fresh unchanged headless gameplay passes in 35.569s.
The prerequisite bound correction is therefore qualified before native replacement.
No native Go drafts were applied during these checks. Local result manifests are
c-qualification.json and c-gameplay-qualification.json.
