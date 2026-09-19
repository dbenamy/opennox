# Client world, spell and item presentation

The 22 selected C bodies are now Go. Four exports remain for actual client decoder
calls; 18 private interfaces and the private chant-tree C global retire. Production
C is **23,648 physical lines in 66 files**, with zero reference C: **−764** from the
corrected baseline. See [the native qualification](client-presentation-native-qualification.json).

The batch covers wall/drawable ordering, circular floor sprite composition and its
copy helper, phoneme icons and chant playback, reflective shields, Turn Undead
particles, persistent/transient ray ownership, book rewards, bubble particles and
NPC equipment appearance. The caller audit follows the main drawable frame pipeline
as well as decoder callbacks. Every selected body had a live caller or private
helper path; this batch does not prune unreachable algorithms.

## Baseline correction

The original persistent-ray allocator checked a counter with no writer. After all
96 pointer slots were filled, another creation allocated an untracked sprite. An
independent repeated-creation test reproduced the live count growing from 98 to 99,
including two endpoints. The existing free-slot search now precedes allocation.
Full tables reject creation; duplicate removal, first-free-slot reuse and bulk
cleanup are checked separately. The sole production caller ignores the return.
This reversible correction is recorded for later review in [DECISIONS.md](DECISIONS.md).

Frozen checkpoint **fc842bfa**, qualified by **9e9c8419**, includes this correction.
The pre-baseline comment increased C by one line, from 24,411 to 24,412. Fresh
production qualification was required; prior production artifacts were not reused.
The original failure remains in
`build/port-client-spell-presentation/presentation-rays-fields/tests.jsonl`.

## Contracts

Sixteen frozen captures contain **7,744 records**. Independent assertions use real
allocation/list, drawable/render, floor composition, GUI, NPC, modifier, spell,
audio-request and book owners. Lookup tables and gesture names come from the
original blob, with the gesture pointer table relocated as in production.

| Capture | Records | Main contracts |
| --- | ---: | --- |
| Wall ordering | 918 | Wall kinds, tile bounds, local-player absence and oriented-side boundaries |
| Drawable ordering | 1,728 | All 32 directions, signed positions and local-player ordering |
| Copy | 160 | Whole-pixel spans, unaligned source/destination addresses and guard bytes |
| Floor composition | 1,346 | Literal/transparent spans, clipping, circular wrapping, image/sprite offsets, signed heights, frame stamps and early exits |
| Equipment | 300 | Weapon/armor capacity, full arrays and holes, aggregate masks and modifier identity |
| Phoneme initialization | 9 | Every image-loading failure boundary, loaded prefix and actual window placement |
| Phoneme drawing | 160 | All eight icons, real pixels, frame zero/wrap, draw-before-expiry and subsequent absence |
| Turn Undead | 108 | All 43 angles per successful emission, allocation failures, velocity bits, origins, frame and list ownership; no RNG consumption |
| Shield ownership | 9 | Partial allocation, flag-setting prefix, counter reset and repeated cleanup |
| Shield drawing | 120 | Eight directions, signed heights/positions, real callback and pixels |
| Ray events | 528 | Kind selection, static/dynamic endpoints, missing endpoints, allocation failure, midpoint truncation and stored payload |
| Ray capacity | 193 | Allocation cap, duplicate removal order, first-free-slot reuse and complete cleanup |
| Transient rays | 5 | Empty/full lists, deletion order, persistent-ray isolation, beam reset and repeated cleanup |
| Bubble particles | 1,296 | Named/default colors, source RGB byte narrowing, parameter layout, signed height, deadline and actual list ownership |
| Chant | 720 | Real phoneme tree, timer wrap, cancellation, icons and voice-specific sound requests |
| Book reward | 144 | Cooperative throttling, sequence wrap, placement, reward kinds, particle count and autofill |

Valid-input boundaries remain explicit: sprite streams supply complete whole-pixel
runs, private ray counts describe at most 96 entries, and chant data supplies a
valid tree path. This port does not define new behavior for corrupted private state.
The four live decoder exports are exercised through C in the native test bridge.
Their three unused pointer-valued results became void in both export and header.

Fixture-only setup corrections supplied a real GUI, explicitly typed a constant
for 386 and converted a named image-handle type. No frozen expectation changed
during translation. The first native build passed all focused captures unchanged.

## Qualification

All shared captures match original C and each other. Source fingerprints remain
unchanged throughout each target/production gate: 2,545 baseline source files and
2,549 native source files.

| Target | Affected test roots | Captures | Records |
| --- | ---: | ---: | ---: |
| Default | 265 | 184 | 278,640 |
| High resolution | 265 | 184 | 278,640 |
| Server | 264 | 183 | 271,728 |

The server exclusion is the established client-only object-render occlusion test.
The accumulated selection includes affected combat overlays, render owners,
particles/effects, floor raster/composition, spellbooks and spell lifecycle.
The production gate initially rejected a manifest classification mistake: the
original C bodies had been listed as Go-backed exports. Correcting `retained_c`
resolved it without changing source.

Memory-map static checks pass. Three fresh 386/SSE2/CGO binaries pass interface
retention/retirement and test-helper exclusion checks. Direct `nm` checks also
confirm removal of the private chant-tree C global. The full suite matches the
known 1,553 failure entries and package outcomes exactly (32 skip, 15 pass, 3 fail).
Fresh headless gameplay and explicit save, reload and resumed-gameplay checks pass.

## Recovery and disk

Manifests: [C qualification](client-presentation-c-qualification.json),
[native qualification](client-presentation-native-qualification.json),
[frozen captures](client-presentation-captures.json),
[C batch](client-presentation-c-batch.json) and
[native batch](client-presentation-native-batch.json).
Local evidence is under `build/port-client-spell-presentation`, including
`c-final-*`, `c-qualified-production`, `native-focused`, `native-final-*`,
`static-native.log` and `native-global-retirement.json`.

All copied fixture/native drafts, the installer and freezer are **consumed**.
The installed source is authoritative. Qualifier scripts verify this completed
state; do not regenerate goldens or replay installation during later work.

Completed C and native scenario asset copies were deduplicated only after hash and
inactive-process checks, reclaiming 1,112,747,701 bytes per pair. Their
`deduplicate-client-presentation-{c,native}-assets.py --restore NAME` helpers and
manifests remain; audit/apply modes are consumed. Original assets and the archive
remain untouched. A further six superseded combat-overlay binaries were removed
when disk space became tight; hashes, revisions and rebuild instructions remain in
`build/port-combat-overlays/removed-superseded-binaries.json`. All captures and the
current C/native production binaries remain available.
