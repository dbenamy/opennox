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

Current status: C references locked and related baseline passed. Production C
remains **95,146 lines /130 files /zero test-reference C**. The native translation
is staged locally and has not replaced C yet. The C baseline is recoverable from
the fixture commit accompanying this document. Local scope/signatures/original
sections and raw evidence are in `build/port-client-draw-particles`.


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
