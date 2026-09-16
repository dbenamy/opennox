# World-wall rendering and drawable visibility

Status: **qualified, frozen C baseline**, after qualified/pushed minimap conversion
**33e4247c**. Production remains **76,382 C lines / 91 files / zero reference C**.
No world-wall production changes. Expectations are frozen and qualified.

Scope: nine routines / 620 function-block lines in GAME2_1.c and GAME2_2.c: wall rendering, two
viewport projection helpers, local/team visibility, four drawing-pass predicates and the private image-interval
helper. That helper has no remaining caller after the renderer is converted.
The projection helpers overlap existing Viewport methods; one remains a private
C call from briefing presentation. The caller audit found no other C callers for
the selected batch; reassess callbacks and static draw-function identity on port.

Applied fixtures use actual viewport/drawable/player/team/wall/render/image owners.
Predicate cases compare actual static/animated callback addresses and require
partitioning of eligible active drawables, no-callback exclusion, inactive-pass
exclusion and no drawable mutation. Projection cases check signed overflow and
compare existing Go viewport methods with original C. Wall cases cover ordinary
and edge rendering, all wall directions, front/back/translucent and high-resolution
configuration inputs. Type-3 opaque RLE images go through the real image/handle
owner; edge drawing uses the production pixel-row initializer and renderer buffer.

## Baseline development and qualification

Development-a passed three roots / 5,104 records in 188.398s. The pixel buffer is
now pinned while C's actual row table retains its addresses; cleanup detaches and
frees that table before unpinning. Development-b passed five roots / 5,668 records
in 22.695s, adding real FOV scanlines and local/team visibility; all prior hashes
were unchanged. Development-c passed seven roots / 5,778 records in 22.506s, adding
wall variants, nonuniform light samples, clipped pixels and nil-wall guards.

Development-d passed eight roots in 187.251s after adding independent interval
clamping/cache tests for the private helper. Frozen **6,262 records /eight groups**
passed independent three-target affected qualification. Server has **5,734
records /seven groups**: FOV clipping is client-only because the actual server
clipper is deliberately unreachable. Ordinary/edge wall rendering and all other
contracts execute on server too; no replacement clipper is installed.

C qualification includes wall, minimap, briefing, sprite/renderer, UI clipping,
AI path, wall deletion, map painting/population, objectives/team and gameplay text
owners. Exact production fingerprints must match qualified 33e4247c before reusing
its three binaries/ABI, exact full-assets known-failure comparison and fresh
twelve-frame chapter/minimap replay. That scene displays the world walls and was
validated with updates disabled against independent C references. The Go wall
conversion will require new builds, failure comparison and fresh gameplay.


Scope, callers, drafts, logs and captures are under ignored build/port-world-walls.
Applied drafts become stale; use actual source and PORTING_STATE.md. Original game
assets/archive remain unchanged.

## Frozen expectation hashes

| Group | Records | SHA-256 |
| --- | ---: | --- |
| drawable-passes | 4096 | `e479a374b5fb1e9675eeeedbebdf87cfda0ba0acf649074b1ce3aab2e14ae1c9` |
| field-of-view | 528 | `08bc8250f4ea083893469e27c94b575da89834afec5cbb13925db95ef1f30330` |
| image-intervals | 484 | `0d151547bb49a85dba116b14df3ecbedb133c2c746611f3d127232a9ad95a96e` |
| nil-input | 2 | `6898cad19c98a98f8f3dc515f15393d46c4ba518eb03b2fece4023cab3d04a45` |
| player-visibility | 36 | `fdfd24d199a7ad96fced54d2b512796c56fab6166aac3f378a2425599c2996cd` |
| projection | 128 | `ddce5d8f5e2b8f4bc3e624561834096e511e8a3eff1925e7a57872c0a550f214` |
| rendering | 880 | `8f4eb579e0681cf61f69161f1ca8ffcf55f00466c11eae41ee7694b69c654702` |
| variants-lighting-clipping | 108 | `bdbefd26349eb6451f1e2fe22d67e990e76a09d738f422250f7312004bb731ec` |

Go drafts are prepared but unapplied. The renderer preserves the sampler's reused
light buffer by copying its first color before the second query, delegates FOV
intersection to the actual client, and retains the shared C edge-rendering owner.
The private image-interval helper moves with its caller; shared state stays with
its existing owner. No test expectation is regenerated from a Go draft.

Completed C qualification: **246 / 244 / 246 roots** (default/server/highres),
all started and completed, with the expected prerequisite-probe skip in each.
Times: **131.375 / 213.129 / 135.790s**. All applicable hashes match exactly.
Exact production fingerprints match 33e4247c, validating reuse of its three
binaries/ABI, full-assets known-failure comparison and fresh twelve-frame gameplay.
All readers joined before committing the baseline and applying Go.
