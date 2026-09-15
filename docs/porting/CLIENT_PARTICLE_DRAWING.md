# Procedural particle drawing and lighting

The drawable-update conversion is pushed at `06133afd`. This next batch covers
19 routines /731 physical C lines: six whole files for magic, bubble, rain,
falling sparks, spider spit and vortex drawing; five light-property setters in
GAME2_2; and both shared palette/color initializers in GAME3. Moving the color initializer
with its palette helper permits retiring both C bridges and routing maps.go
directly to Go.

The new fixture reuses the existing effects renderer and drawable owner. It
installs actual draw callback addresses and invokes them through the drawable's
callback slot. It captures complete normalized drawable state, pixels, renderer
state, globals, ownership/deletions and RNG consumption. Creation cases include
failure injection, viewport offsets/clipping, frame progression and deadlines.
Additional contracts check light-field layout, scalar/pointer returns, clamping,
minimum/monotonic radius, palette samples, and explicit bubble/vortex lifecycles.

The fixture loads the real production light-radius coefficient table, saves and
restores the named light constants separately from mapped storage, and invokes
the production palette initializer. The first palette uses signed-char arithmetic;
its actual colors are the baseline, not a newly designed gradient. No sprite
image loading is needed for these procedural primitives. Animation/image parsing
and the separate screen-particle allocation list remain later connected work.

The first drawing capture exposed an impossible fixture state: shrinking at size
zero. The real bubble lifecycle switches phase when it reaches zero. Correcting
phase/size combinations resolved the renderer panic without a production change;
zero-size growth and deletion cases remain. Explicit lifecycle contracts pass.

Current status: native conversion fully qualified. C baseline `423a0a3a` is
pushed and recoverable. All 19 routines are translated; six whole C files and
three private C entry points are removed. Sixteen live callbacks/exports remain
and are exercised through their actual C ABI. maps.go calls the Go color
initializer directly. No C algorithms are retained solely for tests.


The expanded C capture repeats byte-for-byte (c-d/c-e): 4,647 results /four groups,
plus explicit bubble and vortex lifecycle contracts. References are now locked:

| Group | Results | SHA-256 |
| --- | ---: | --- |
| Drawing callbacks | 4,359 | `7390a7366081a3fed2a09630095746f735ad3f505f019875b2905aa2057f6b76` |
| Light properties | 272 | `1ceb08e8d9bd30d4612500c837222766cd96ad3be88ebf2638eb6efa3754192b` |
| Palette initialization | 4 | `617ccd2159c45848fb40dc4d6732f87f6731e5118a65987b284612663dac2857` |
| Color initialization, cold/warm | 12 | `095076e2cab218dab5797bb2f73a02949f6a9ccf6fe2d6b5e50c03acf18021c9` |

The related effects/update/drawing standard suite passed all 33 selected tests
in 22.350s before the C-baseline checkpoint. This fixture extends no existing owner API and makes no production
change; related regression coverage is sufficient for this baseline. Native
qualification will run accumulated standard, affected server/highres, all three
production builds and ABI checks, full asset-suite comparison to known failures,
and unchanged fresh headless gameplay. The full suite is included because this
batch moves shared lighting properties and production color initialization.

Review later: the source's first palette uses a signed-byte decrement, producing
mostly black entries and a final colored entry. Preserve that exact behavior in
this port; changing the intended gradient needs separate visual/reference evidence.


Final native validation:

| Check | Result |
| --- | --- |
| Focused comparison | All 4,647 locked results and lifecycle contracts match on first native run; 164.196s |
| Accumulated standard | 666 selected/completed; 665 pass and one optional prerequisite skip; 370.211s |
| Affected server | All 33 selected tests pass; 163.456s |
| Affected highres | All 33 selected tests pass; 28.651s |
| Production builds | Standard/server/highres pass; independent builds ran alongside tests |
| Binary checks | ELF32/i386, SSE2, CGO; three retired symbols absent, 16 retained exports present, test helpers absent |
| Full asset-backed suite | Exactly 1,553 known failure entries and the same package outcomes: 15 pass /3 fail /32 skip |
| Fresh gameplay | Original assets, unchanged warrior scenario, override=false, Xvfb/null audio; exit zero in 37.574s |

Production C is **94,415 physical lines /124 files /zero reference C**, down
**731**. Accumulated captures are **120,397 results /1,065 groups**, plus
independent contracts. Raw evidence and binary hashes are in the local
`qualification.json`; tracked tests/hashes and baseline Git revision support
recovery without those local artifacts.

Preserved details include separate named/mapped color storage, signed-byte
palette generation, exact floating-point light conversion and return widths,
clamping, unsigned vortex trailing-point arithmetic, phase/timer/deletion rules,
alpha state, viewport clipping and allocation/RNG order. No production behavior
repair was needed; the only prerequisite correction was to an invalid fixture.

Next investigate sprite data/animation and its real image-handle owner. A read-only
plan is in `build/port-client-sprite-animation/PLAN.md`; no next-batch fixture or
production edit is applied. Audit scope and image/data ownership before capture.
