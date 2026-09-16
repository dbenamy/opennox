# Things-section cursor readers

Status: **qualified Go conversion**. Frozen C baseline: **950c536b**.
Removed **398 C lines**; **72,542 lines /91 files /zero reference C** remain.
Previous production baseline: floor/edge conversion **bea20162**, with
**72,940 C lines /91 files /zero reference C**. Scope: eight connected functions
and 397 C block lines. Root `things.go` uses these while traversing sections for different loading
passes: audio definitions/events, spell/ability definitions, image records and wall
records. This batch preserves traversal; actual definition/image owners remain.

The wall reader aligns its count fields relative to the beginning of the MemFile:
it skips to the next eight-byte boundary and consumes an entire eight-byte block,
using its first byte. It accepts up to eight material names; larger counts stop
after the eighth name. Fifteen orientation groups each contain four image faces
per entry. Scratch-buffer writes and END consumption are observable.

Spell/ability/image/audio counts are signed; nonpositive values consume only the
count word. Length fields and inline/indexed image references have distinct skip
rules. Image kind2 supplies an animation count and extra fields; all other kinds
skip one image reference. AVNT tag0 succeeds, known tags skip their specified
payload and continue, and unknown tags stop immediately. Preserve the remaining
C event-reader caller of the AVNT inner helper with one narrow Go export.

Actual-owner contracts use C-compatible MemFile allocations, complete input and
guard bytes, exact cursor progression, whole-input immutability and full scratch
contents. Matrices cover all AVNT tag bytes, repeated/reversed event sequences,
zero/max byte lengths, positive signed-16-bit length boundaries, zero/negative
record counts, animation variants, every absolute alignment offset, material
limits, orientation counts and good/bad END. Concatenated sections in both orders
check that a reader leaves the next one at the correct cursor.

Do not test impossible truncated records or negative signed string lengths against
the original unchecked readers. Those inputs are outside this qualification.
No original C algorithm will remain solely for tests. Case counts describe the
matrix size, not the quality of coverage by themselves.

C production is unchanged from the preceding qualified floor/edge reader port.
Reuse its builds/full-suite/gameplay evidence only after exact production source
fingerprints match, and audit the seven original C interfaces in saved binaries.
Run all affected root tests and new C captures in separate default/server/highres
processes. Go qualification requires fresh production builds/ABI, the exact known
failure set and both normal and GUI flat-floor gameplay replays.

Ignored source drafts and evidence are under build/port-thing-skips. No preparation
or application script may be rerun after it has mutated source. Join every test
and build reader before editing Go/C/header files. Original assets remain intact.

Initial C capture passed seven roots in 187.476s. Added public wall scratch-buffer
contracts pass with all eight roots in 36.144s; the seven earlier group bytes are
unchanged. **7,737 records/eight groups are frozen**; three-target C qualification
is complete. Native implementation is now applied and qualified. Preparation/freezing scripts
are stale and must not be rerun.

The final caller audit includes the aligned MemFile helper: its only two callers
are in the wall reader. Retire its 18 C lines and declaration with those callers;
wall contracts already cover the actual aligned-byte behavior. Preserve the AVNT
inner export for its separate live C caller.


## Buffer-view compatibility follow-up

The original floor/edge/wall wrappers check backing capacity, then pass `&buf[0]`
to C. A nonempty short slice with sufficient capacity therefore behaves exactly
like a full-length slice; an empty slice is rejected before the reader executes.
The prior Go floor/edge port used the shorter length for terminator indexing.
Actual game callers pass full buffers, so the existing gameplay qualification did
not exercise this API edge.

In this native batch, retain the capacity checks, preserve empty-slice rejection
and expose the backing buffer to the Go readers. The wall behavior is now frozen
from actual C, including cursor and whole-buffer rejection contracts. An added
independent floor/edge test will compare accepted views with full-buffer behavior
and assert that rejected views preserve owner state. The original a4d69a64 bridge
establishes their equivalence by passing the same address regardless of nonempty
slice length. No frozen floor/edge expectations will change. The fix/test are now applied and qualified alongside all existing floor records.

## Qualified C checkpoint

| Group | Records | SHA-256 |
| --- | ---: | --- |
| audio | 480 | `c9ea00dc6ae76204249978d9fbe9e416a87786a5ac3d6487dcaf222d82360e24` |
| concatenated | 112 | `198057bc2e0c28bab1385a6a52f2a437c73c7db526b09bf112e0a5d6d44bd127` |
| event-sequences | 144 | `69c6b2f8281f179c14a944cf5916efbb6b301e5acb57141c8c9a8cf9c51b84d0` |
| event-tags | 2048 | `eaf550105afc29ecbf1579f004c32119341eace14399dd9044cd87770bf70b07` |
| images | 720 | `9221f54959b72482dab6df3101b01a4c889fab7bd5736bf577440c8fb8df30d6` |
| scratch-capacity | 9 | `6bfdd59d5861b06669856b724decc01a413bf782ba2419f143bb8d6217eafe38` |
| spells-abilities | 1152 | `62b79a08d4088d3101357d93d270e575bb1e0b877306b60a30e6675507daec1d` |
| walls | 3072 | `657885a2ecd7241a10a7564ea541b30ed3b8034b1286ffe080cfc0427553faac` |

Affected roots **306 / 304 / 306** completed in **169.272 / 254.122 / 175.717s**
(default/server/highres). All 7,737 records repeat exactly in separate processes
across all targets. Production fingerprints exactly match **bea20162**, so its
three production builds, exact known 1,553 failure entries/package outcomes and
26 matching gameplay frames remain valid evidence. Saved binary hashes match
their original qualification, and all eight original C interfaces are present.
At the C checkpoint, production was unchanged and the buffer-view correction
was still an unapplied native-stage plan. Source fingerprints remain
unchanged throughout. Evidence: build/port-thing-skips/c-qualification.json.

## Go conversion and buffer-view regression coverage

Go now owns all seven section functions and the aligned-byte behavior. Seven C
interfaces retire; the AVNT inner helper retains one narrow export for its actual
remaining C event-loader caller. No test-only C algorithm remains.
Removing 397 function-block lines and the now-terminal separator blank removes
398 physical C lines.

The public wall reader and five preceding floor/edge readers preserve the original
capacity-based buffer contract. Existing capacity checks remain, empty views are
rejected before calling a reader, and accepted nonempty views expose their backing
capacity. Original-C wall captures cover this directly. The independent floor/edge
root checks all five readers, accepted short/full views and unchanged owner state
on rejection; its original-C basis is the address-only bridge in a4d69a64. Existing
4,396 floor/edge frozen records remain unchanged.

Focused 16 roots passed in **183.071s**.
Affected roots **307 / 305 / 307** completed in **248.978 / 252.756 / 179.232s**
(default/server/highres); all 7,737 section records match frozen C. All three
production builds pass ELF32/i386/SSE2/CGO, retained/retired-interface checks and
absence of test helpers. Exact known failures remain 1,553 entries with unchanged
package outcomes. Fresh normal 12 and GUI flat-floor 14 gameplay frames match.
Source fingerprints remain unchanged throughout. Evidence:
build/port-thing-skips/native-qualification.json and native-evidence/.
