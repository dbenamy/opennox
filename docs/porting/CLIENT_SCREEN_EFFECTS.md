# Screen effects and remaining drawing port

Scope: 16 routines in the five remaining client drawing C files (harpoon/rope,
maiden, undead effect, waypoint and screen particle), screen-particle allocation /
list / deletion / callback walk, rope and circle raster helpers, and three shared
distance/table-square-root primitives. The selected C source occupies 561 lines.
Shared object/monster rendering and renderer-bounds setter remain dependencies.

## C prerequisite and decisions for later review

The rope constructor stores endpoint IDs with bit 15 marking static objects;
drawable creation stores the cleared low 15-bit network ID. The original rope
callback tested the marker but passed the encoded ID into exact static lookup,
so it never found actual static endpoints. Clear the marker before both lookups,
as the constructor already does. The independent original-C contract failed on
the first static case; corrected C renders all four static/dynamic combinations.
This reversible correction changes no physical C line count (92,869 / 106 files).
Review this prerequisite separately from the language conversion.

Preserve unsigned multiplication and logical right shifts in the circle helper,
including the resulting large coordinates and clipping. Its visual behavior may
warrant a later correction, but this translation does not redefine it. Preserve
the table-based distance approximation, 32-bit squared-coordinate wrap and RNG
ordering. These are observable in existing callers.

## Owned tests

Nine root tests contain six matrices (149,268 results) and three independent
contracts/probes. They use the qualified real drawable pool/list/spatial index,
material-indexed RLE images, rasterizer, font and lighting owners. Screen-particle
tests use the actual allocation class and linked-list head/tail, with capacities
0, 1, 2 and 7, exhaustion/tail reuse, child generation and callback mutation.
They capture all 52 particle bytes with normalized links, rendered pixels, full
renderer state, bounds, GUI flags and both RNG indices. Independent contracts
check nil viewport, stationary decay, preserved upper bytes on reuse, final shrink
and strict viewport deletion. GUI and canonical sqrt-table ownership were added
after early fixture failures; no rendering or allocation algorithm is stubbed.

Maiden tests own actual NPC storage, readonly server object lists/update colors,
global client/server callbacks, gameplay flags, chat/ally state and full vector
animations. They call the actual monster renderer and verify visible pixels and
material changes; missing NPC, absent/nonmatching/matching server lists, palettes,
animations, directions and frames are captured. This does not qualify every
player/team/name/enchantment branch in the shared monster renderer.

Rope tests use real dynamic/static network lookup, raw coordinates, missing
endpoints, direction offsets and viewport transforms. Other draw cases cover
harpoon slave animation, undead effect ages 69/70/71, frame wrap, RNG/spawn state,
distance-dependent radii including zero, and rotating waypoint pixels. Primitive
cases cover every input below 65,536, all scaled table bucket edges, randomized
32-bit values, coordinate wrap, exact power-of-four anchors, and real raster
returns/pixels/state across alpha values, clipping, endpoints and circle radii.

Final C captures c-g/c-h match byte-for-byte. One attempted repeat omitted the
386 environment and failed at discovery; it supplied no baseline data and was
rerun with the correct environment. All qualified captures use 386/SSE2/CGO.
Only hashes/expectations are tracked; raw captures and qualification manifests are
local in build/port-client-screen-effects. No C algorithm is kept solely for tests.

| Capture group | Results | SHA-256 |
| --- | ---: | --- |
| screen-effects-distance | 138366 | `1e2c756f945ea9f98383bf55d6b20e6e4c9350847694e786445d6ac4d116fe83` |
| screen-effects-drawables | 1632 | `b3c8897faa5ff57c99b2d71ea2f446e4df2ac4830d1eafb501b4471bbad6c628` |
| screen-effects-maiden | 1728 | `ad6a3815de864ed82b56a1b1c1f868705babe7f483a7bbade2fb14e77158701f` |
| screen-effects-particle-lifecycle | 6048 | `d21c6e420b68c92a842a92948104828a13d69eb8e612c4c06039c105f7c83123` |
| screen-effects-raster | 234 | `bb8f569d31fb9889f81a593e53441795df1747cc292e13c889c24cb53d83faf2` |
| screen-effects-rope | 1260 | `237972bd0664fcf9a99e92a7943e20416b8404a32676ca1a206f9bfc49595305` |

## Qualification

C qualification passed: focus 9, accumulated default 695
(one expected optional prerequisite skip), server/highres 62/62.
All selected tests started and completed. Production C build passed in
51.069s; fresh headless gameplay passed in 35.928s
with null audio and reference comparison enabled.

Native conversion passed all 149,268 frozen results and three contracts on its
first attempt; no implementation corrections or expectation changes were needed. The
16 routines now live in Go; 562 C lines and five whole C files were removed.
Production C is 92,307 lines / 101 files, with zero test-reference C. Existing
Go callers invoke private Go helpers directly. Eight required ABI bridges remain
(six draw callbacks, particle creation, two-coordinate distance); eight internal
C entries and their declarations are retired. Shared pool/head/tail globals stay
compatible with the existing initialization/reset owner. Callback walk saves its
next pointer before dispatch, and tail reuse preserves the untouched upper bytes.

Native qualification: focus 9, accumulated default
695 (one expected optional skip), server/highres
62/62; all selected tests started and completed.
Three production builds passed and were verified as ELF32/i386, SSE2, CGO enabled,
with all required exports present and retired entries/test helpers absent. The
full asset suite matches the existing failure multiset and package outcomes;
fresh headless gameplay passed in 38.049s with null audio and
reference comparison enabled. Accumulated frozen coverage is 300,569 results in
1,085 capture groups plus independent contracts. Local qualification.json records
exact timings, binary hashes, suite comparison and gameplay result.
