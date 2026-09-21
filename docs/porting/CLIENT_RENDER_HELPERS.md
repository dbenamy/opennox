# Client presentation and lifecycle

These helpers are now Go and fully qualified: **3,543 C lines remain in 28 files**,
915 fewer than the baseline. All 6,129 captured cases match C; three additional
safety contracts and fresh production/ABI/gameplay/save-load checks pass.

The original-C baseline covers nine test roots and eight captures /6,129 cases:

| Behavior | Cases | Independent checks |
| --- | ---: | --- |
| Palette expansion, sorting and compaction | 20 | All 256 entries; packed unsigned ordering, index ties, source preservation and tag bytes |
| Drawable looping audio | 2,340 | Range/gates, signed arithmetic, new/existing/stopped events, volume and pan |
| Loading overlay | 16 | Frame gate, timestamp, image cache, visible text and pixels |
| Winner overlay | 72 | Both modes, image cache, window sizes, multiline/wrapped text and unchanged input |
| Debug overlay | 405 | Map text, real drawable/player ownership, all classes, signed byte levels and viewport origins |
| Random names | 3,200 | All 33 names, cached count, exact random stream consumption |
| Server-address copy | 40 | Byte truncation, zero padding, embedded NUL, UTF-8 boundaries and untouched tail |
| Audio context lifecycle | 36 | Nil/repeated cleanup, all voice kinds, callback order, stop versus destroy, device references |

An additional scalar/host-info contract checks mapped state and stable host-info
identity without normalizing or capturing addresses. Server builds select six
roots/five captures; the three rendering roots are client-only.

The original baseline had **4,458 physical C lines in 33 files**, zero reference C.
The baseline adds only tagged test fixtures. Production source is compared with
qualified queue revision `2ce59fea`; its production/ABI, known-suite and headless
scenario evidence can be reused. Native conversion will need fresh qualification.
See [the baseline manifest](client-render-helpers-c-batch.json) and
[qualification](client-render-helpers-c-qualification.json).

## Scope and compatibility decisions

Translate the connected presentation, palette, random-name, state/host helpers
and audio context wrappers together. Move the live shared global definitions when
an otherwise-empty C translation unit is removed.

Whole-source symbol and mapped-offset searches find the old focus callback arrays
(`754088`, `754092`) and their counts only in definitions, registration writes,
shutdown frees and metadata. Nothing reads or invokes registered callbacks. Remove
that registry and its private pause/empty audio callbacks. Preserve the live
`dword_5d4594_251744` reset beside registration; this state is used by input/session
code. Audit the input-reset callback interface before retiring it.

The pixel-span allocator globals (`1301844`, `1301848`) are initialized to zero
and never assigned anywhere outside their own cleanup. Remove their inactive
cleanup after preserving the neighboring live render-data global.

Defer animation-data teardown to a separate resource-ownership batch. Its nested
allocation graph deserves direct ownership coverage; broad allocator instrumentation
solely to meet a line-count target would add unnecessary work.

Preserve the C looping-audio pan coordinate: viewport word 6 is `World.Max.X`.
Do not substitute the usual world-to-screen helper. Preserve C signed 32-bit
arithmetic and the qualified target's floating-point conversion for overflow cases.
Audio context disposal must unlink the context before freeing it; calling the
lower-level free helper alone leaves a dangling global list entry.

Winner messages in the C baseline stay within its 127-unit local UTF-16 capacity.
If the Go implementation safely accepts longer text, cover that separately rather
than exercising the original C stack overflow. Only winner modes 0 and 1 have live
callers. Random-name cached counts cover valid values 0 through 33.

## Baseline fixture repairs

Before freezing, the audio expectation was corrected from `World.Min.X` to
`World.Max.X`. Winner/class-name fixtures now bind the exact shipped tables and
startup pointers instead of reading zero-filled memory. Overlay white is constructed
with the real RGB5551 encoder: literal `0xffff` sets the transparency bit and made
text invisible. The visible-text contract caught this before goldens were frozen.
The earlier checkpoint's claim that loading passed was incorrect; the original
`c-overlays1` log also records that failure. No production fix was needed.

The first manifest attempts compiled successfully but rejected a multiline test
list; the runner requires a single regex. Corrected runs use identical source and
compare every applicable capture hash. Failed probe logs remain under
`build/port-client-render-helpers`.

## Qualified native implementation

The installed conversion removes five C translation units, moves their live global
definitions to the existing shared owner, and reduces C by 915 physical lines to
3,543 in 28 files. Sixteen live helpers (including the input-state initializer) move
to Go; six inert helpers disappear. Fifteen obsolete Go-backed exports retire;
[the symbol list](client-render-helpers-retired.json) records them. The original-C
baseline is pushed as 78bef3a6. No C algorithms remain solely for tests.

Initial native probe2 passes all 11 selected roots and matches all eight C capture
hashes. Probe1 failed to compile because the flags package's declared name differs
from its import path; an explicit import alias fixes it. Default/server/highres pass 96/91/96 affected roots with no skips and all applicable
C captures unchanged. Fresh headless gameplay, all three production builds and ABI
checks, the exact known-suite comparison and explicit save/load pass. The preflight
and qualified default binaries have identical hashes. See
[the native qualification](client-render-helpers-native-qualification.json).

Additional reversible boundary behavior: winner messages can exceed the former
127-unit stack buffer; invalid winner modes return without drawing; zero-width
(or ±1-width) audio viewports use centered pan rather than dividing by zero. The
debug overlay caps text to its 80-unit scratch allocation, always terminates it,
and leaves neighboring memory untouched. Three independent native test roots
cover these cases. Valid captured C behavior is unchanged; goldens are not updated
to accept the implementation.

The first preflight command built the root library package and therefore could not
execute (exit 126). The corrected command targets `./cmd/opennox`; this was a runner
setup error, with production source unchanged. Keep both logs for recovery.
