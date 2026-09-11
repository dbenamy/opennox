# AI lifecycle and item searches

This batch covers the five remaining C-registered actions: HUNT, PICKUP_OBJECT,
DYING, DEAD and GET_UP, plus zombie revival, released souls, burn deletion,
monster-state reset, and edible/usable item searches. It retires 16 C bodies
(534A90/AB0 and 5449D0–544F70). Revival and both search owners still have C
callers, so their ABI entry points must remain.

The original C passes 5,632 generated complete-state cases and 174 independent
contracts using the shared guarded AI fixture. Generated hash:
`15d5986a08af2c44eb120f178e45366cb014817a42c4788a21db0d892ce4c485`.
Contract hash:
`b6509bbe7747b05c7b8f48f58bc8505f520ab2510743f731b33292cf8b42a6e8`.

The fixture captures object/monster-data changes, action stacks, RNG indices,
shared caches, health, deferred sounds, scripts, indirect C callbacks, created
object layouts, decay-list membership, deletion calls, and actual queued shield
and spark packets. Inventory placement and weapon-class eligibility are recorded
call boundaries with configurable results; their existing implementations remain
outside this batch. Map iteration, vision, object allocation and packet queues
use their production implementations. Minimap removal sees an initialized empty
list, and kill rewards see no last damager. C instrumentation records callbacks;
it contains no copied lifecycle algorithm.

Independent checks cover strict/wrapping revival deadlines, zombie types and
stay-dead flags, soul subclass/type gates, zero motion before the dead callback,
clearing Update (+744) while retaining UpdateData (+748), max-health restoration,
raise ABI return values, stack capacity, death sound/event arguments, item filters,
pickup equality rejection and argument rereading after placement (even if it
returns false), nearest candidates and first-visited ties, and burn packet bytes.
Duration conversion rounds the balance double to float32, then truncates through
the engine conversion routine. A search precision discriminator checks minimum
`0x3f800003`: both coordinate deltas remain double through squaring, while the
stored minimum is float32. Original compiled x87 instructions confirm this.

Artifacts are ignored under `build/port-ai-lifecycle`. At the baseline checkpoint,
production remains unchanged: **138,212 physical C lines**, 153 files, zero
reference C. Native conversion and batch qualification follow this checkpoint.

## Native conversion

Original-C checkpoint: `8bae3340`. Native action registration replaces the last
five C registrations and removes the obsolete C action adapter. The three live
C entry points use generated Go exports; native navigation calls food search
directly. Shared zombie classification remains C-backed to preserve the original
cache words until its caller family is ported. Packet/reward/scorch engines and
the float-to-int conversion remain shared callees.

Production C is now **137,882 physical lines (minus 330)**, 153 files, zero
reference C. The port uses typed velocity/force fields and the existing audio,
map, action, allocation and minimap APIs. Review corrected shared-cache and ABI
return handling; parity tests caught an inverted poison filter in the draft.
Repeated baseline testing also caught an unnormalized optional metadata pointer;
the fixture now records its identity and frees that allocation explicitly.

## Qualification

Both original-state hashes match. Accumulated port tests pass in default,
server and highres. All three production binaries build as ELF32/Intel80386
with GO386=sse2. The full suite matches exactly 1,553 known failure entries
(15 pass, 3 fail, 32 skip packages; no added or removed failures).
Fresh `ai-lifecycle-port` gameplay exits 0 against both preserved screenshots,
with overrides off under Xvfb/OpenAL null. All qualification is complete.
