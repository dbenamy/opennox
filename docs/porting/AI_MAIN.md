# Main monster AI decisions and defensive reactions

This connected batch owns the main monster decision loop (547210), idle
vocalization (5469B0), dangerous-unit classification (547120), attack unwinding
(5471B0), dodge selection (547C50), and shield-threat owner/callback
(533E70/533EB0). Cast, damage, visibility, map and inventory engines remain shared.
Native callers should use Go directly; these C bodies have no external C owner
once the batch is converted.

## Original-C fixture

The shared AI fixture supplies guarded object, monster-update, definition,
health and action-stack memory. It reuses the existing audio/script/placement
recorders and adds real host-player input, two synthetic cloud types, a real
128-by-128 tile grid, and spatially indexed missile/cursor objects. Spell-user
calls record the spell, source and three argument words at the existing Go
boundary. Retained visibility, obstacle, tile, selection and morph engines run.

The corpus exercises cadence, class/status/capability gates, wrapping frame
arithmetic, aggression, health, buffs, stack unwinding, direction tables and
missile selection. Independent contracts check action and side-effect ordering,
logic RNG consumption, cursor reaction distance, health/flee boundaries,
teleport cooldown decisions, stack capacity, cancellation, movement frustration,
nearest missile ties, and late tile rejection after five dodge attempts.

Fixture corrections are part of establishing the C baseline: spatial objects
need the active flag before indexing; copied missile indices must be reset;
callback and embedded spatial-node pointers need stable snapshot IDs. A generated
DEAD action without its dead flag caused the original C to reject a confusion
push and dereference its null result. Main-loop cases now use coherent lifecycle
flags. This does not change production behavior. RNG indices start at the seed;
these helpers use the logic stream, not the other stream.

## Arithmetic evidence

Compiled C keeps shield position deltas and distance in double precision through
the angular test. Only normalized velocity X spills to float32 at that stage;
distance spills before the interaction call and nearest-distance comparison.
Two independently calculated velocity vectors distinguish an incorrect early
distance spill. The decompiled float declarations alone are insufficient.

Movement frustration has an independent case where a full double delta is
slightly greater than 15 but a float32 delta is exactly 15. Dodge magnitude rounds
to float32 before the random sign call; generated endpoints and both RNG indices
are captured. C x87 flags stay unchanged while Go uses 386/SSE2.

Ignored baseline, assembly and draft artifacts: build/port-ai-main. Production C
at the start of the batch: **137,297 physical lines**, 153 files, zero reference C.

## Locked checkpoint

The 1,792 generated cases and 119 independent contracts pass repeatedly:

- Corpus: `9a775ff5fa2e30b043c7d7f076e9d0d0eaf0d7fde28d47914de2a97a471adc6f`.
- Contracts: `8a46dbb2e63ae4253721a2358dd5a5f879f146fd1c9d86922d8614c9698c831d`.

Existing combat, lifecycle and monster-state hashes also pass unchanged.
Additional contracts exercise quest dodge, usable-item morph/restore, use after
placement, guard aggro, and started-action cancellation. Production C is unchanged
at this checkpoint. From src with the baseline environment loaded, run
`go test -tags porttest -run '^TestAIMain' .`. OPENNOX_MAIN_CAPTURE optionally
writes ignored JSON snapshots; OPENNOX_MAIN_CASE narrows generated-case debugging.

## Native conversion

The seven C bodies are replaced by legacy/ai_main.go. Production C falls by
**556 physical lines** to **136,741**, across 153 files with zero reference C.
The unit-AI and public shield wrappers, plus the combat block handler, call Go
directly. Obsolete C declarations are removed; no test-only C algorithm remains.

Both locked hashes match natively. The main loop preserves the captured action
pointer, including rebinding it to the pushed confusion dependency, and the
original sound-set pointer. Cast arguments use C-owned temporary storage across
the retained engine's Go callback. Native enemy aggro and existing state/food
helpers replace unnecessary C round trips. Shield arithmetic follows the compiled
spills recorded above. The completed batch qualification is recorded below.

## Qualification

Original-C checkpoint: `ee96e478`. Accumulated port tests pass for default
(69.604s), server (37.968s) and highres (41.370s). All three production binaries
build and report ELF32, Intel 80386 and GO386=sse2. The full asset-backed suite
matches exactly the established 1,553 failure entries: 15 passing / 3 failing /
32 skipped packages, no added failures. Fresh `ai-main-port` headless gameplay
exits 0 against both preserved screenshots, overrides disabled.

Next: the connected spell-decision/cast-action family 5408A0–541490. Permission
slots are 136 words for spell IDs 1–136; the temporary candidate arrays merely
have capacity 137. Frame comparisons use uint32 conversion under C's usual
arithmetic rules. Reuse the main fixture's cast/morph/map/health boundaries.
