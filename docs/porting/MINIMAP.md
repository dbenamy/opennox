# Minimap renderer

Status: **qualified C baseline** after scoreboard **ddf44816**.
Production remains **77,120 C lines / 91 files / zero test-reference C**. The Go
translation is an ignored draft and has not been applied.

## Scope and ownership

Seventeen routines / **734 C function-block lines** in GAME2_1.c and GAME4_1.c:
zoom controls, remembered polygon level, minimap drawing, wall/door strokes,
objective glyphs, player dots, the minimap/message caller and the two private
AI-debug monster iterators. The completed audit finds no remaining C callers
outside this connected batch. Four existing Go wrappers will call Go directly;
the private AI path-buffer C bridges can retire with their only caller.

Fixtures reuse real renderer, fonts, drawable allocation/list/spatial-index,
player/team and wall owners. Polygons use the production allocator and predicate
with owned vertex inputs. AI debug uses the actual server object pool/list and
path buffer. Test adapters supply inputs and call production code; no substitute
C algorithms. Every changed mapped region, named global and input table is restored.
The lightweight rendering server needs explicit normal team palette definitions
for its actual team-color lookup; the fixture owns and restores that input table.

## Frozen C expectations

**1,599 results / ten groups and root tests**, captured from original C. Independent
contracts check zoom signedness/overflow, setter preservation, supported wall/glyph
visibility, exact pixel translation, clipping, floor caching, wall immutability,
wrapper suppression/clip restoration and unchanged AI path data/iterator completion.

| Group | Results | SHA-256 |
| --- | ---: | --- |
| corner-ray | 5 | `78f6e192d052e1837d2eb3da7d04104173e11279cee02c9aeabfb2d081c53604` |
| debug-overlay | 64 | `d19b1faffd6908755042403f75924a39f57871485d0e0d819f5eb2e72019f50b` |
| door-walls | 320 | `2cce62e05865184c7694799a5dea64349632846a464d4ba51cd77566c157157f` |
| floor-cache | 8 | `c725cbf3b5c1e911e9b7bdb1f5358d020ff92bd58012a9fba5b6e784afe1a719` |
| full-rendering | 144 | `ea3c79d1e250c8cc01d63544be4422f2aa1809c9723c0538c8b0e6dcc8efac8e` |
| objective-visibility | 72 | `0c41cf00a5a44d1291ae59022238f47ef539dfa964bc961942e3ed56b6be0644` |
| primitive-contracts | 120 | `f1e50548eda7abb2b9267dee83e3cb94f88ff08c28c0b0473ba33b48b892ddb8` |
| primitives | 792 | `9494a590aa021cdba721c48dbaab95e8869d63eac6212a8103f7a58b007883aa` |
| visibility-wrappers | 26 | `7ddc6839a902f95722c7f11c77dccb8905664423d1bd891d18ebaa2b4595fc3f` |
| zoom | 48 | `9015214205f48b5e23b92e84a804190a307ad83c87768da20aea2b25bd7f9af3` |

Door cases cover all four directions, diagonal and orthogonal revealed/unrevealed
neighbors, an unrelated-neighbor negative control, missing neighbors and four zooms.
Full rendering varies six zooms, three viewport boundaries, walls and levels,
local/remote players, observers, teammates and crown buffs. Additional tracked
object cases cover ordinary/special dots, crown and ball team colors, local/remote
player marks and flag-bearing-player mode. Debug cases vary disabled/enabled,
empty/short/fractional paths and empty/nonmonster/mixed monster lists.
Top-level cases exercise GUI suppression, map visibility, poison and real messages.

## Original-C development and compatibility findings

- Development-a: four roots / 186.885s; floor contract exposed a preexisting
  polygon corner-ray miss. Development-b reproduced it in 22.880s.
- The predicate's first ray passes through a square corner and counts both incident
  edges. Ordinary floor-cache contracts use asymmetric bounds; a separate explicit
  compatibility case freezes levels **1,2,2,2,2**. Polygon correction is outside
  this batch and remains a separate review item.
- Development-c: five roots / 997 results / 22.515s, all passed.
- Development-d completed eight roots, failing a fixture's fifth neighbor lookup:
  only four adjacency pairs exist, followed by string bytes. The negative control
  now has explicit coordinates; the fixture owns exactly the original 32-byte table.
- Development-e: nine roots / 187.827s, all passed including pixel contracts.
- Development-f stopped on missing team-color fixture input. It completed only
  four of ten roots and is **not** qualification. Supply the palette to the real
  lookup; do not change game behavior or treat missing tests as passes.
- Development-g: ten roots / **190.053s**, all passed, supplying the frozen results.

## Qualification and gameplay

Independent affected-corpus qualification passed for default, server and
highres: **206 / 205 / 206 roots**, respectively, in **112.869 / 211.633 /
117.658s**. All 1,599 frozen results match in all three targets; all selected roots
started and completed. Each includes the expected prerequisite-probe skip. It includes minimap, renderer/drawable, UI clipping, particles, objectives,
team inventory, gameplay text/report and wall/population owners. Full discovery
must match started/completed roots; every minimap group must match the frozen hash.
Exact fingerprints for 1079 production files permit reuse of ddf44816's qualified binaries/ABI
and full-assets result (1,553 known failure entries, 15 pass / 3 fail / 32 skipped
packages). This is reuse of unchanged production evidence, not a new green suite.

Fresh **client-minimap-develop** captured twelve chapter/gameplay/minimap frames.
Visual inspection confirmed visible, enlarged, reduced and closed map states.
Independent **client-minimap-c-repeat** passed all twelve in **51.050s**, updates
false, fresh copied assets/save, seeded RNG, null audio and Xvfb. The tracked
[minimap scenario](minimap-chapter.yaml) and [decoded hashes](minimap-chapter-pixels.json)
preserve recovery. Original assets and archive are unchanged.

Limits: gameplay is solo warrior, not a remote multiplayer session. Team/objective,
level and AI-debug behavior is exercised by real-owner fixtures. Invalid zero zoom
is tested through the setter, never passed into rendering's division paths.
The existing polygon predicate behavior is preserved, not repaired here.

Raw captures (losslessly compressed), scope/caller audits, fingerprints and logs
are in ignored build/port-minimap. Frozen hashes and fixtures are tracked; no
permanent test-only C implementation is retained after conversion.
