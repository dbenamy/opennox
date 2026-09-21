# Client presentation and lifecycle

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

Production remains **4,458 physical C lines in 33 files**, zero reference C.
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
