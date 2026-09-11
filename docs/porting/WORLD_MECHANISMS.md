# World-object mechanisms

Next connected batch: 21 functions / 821 physical C lines. Door motion and its
queued angle helpers, toggle/trigger/switch state machines, elevator/shaft and
movement callbacks/audio, visible/invisible pentagrams with teleport callbacks,
push/blow and directional force callback, phantom positioning and trap doors.
Original bodies/addresses are recorded in scope.json; production remains C.

Reuse temporary/effects/inventory guarded fixtures and effectsTimedRun clock
assertions. Existing script recorder captures event, block offset and exact
caller/trigger identities. New optional mechanism fields/actions must leave all
20,314 previous focused contracts unchanged. Preserve actual spatial iteration,
shape containment, movement/teleport, synchronization and collision-list paths.

Fixture additions: guarded collision data for trap-door timers; optional object
shape setup after temporary fixture's center-index setup; save/restore door angle
queue and trigger/pressure-plate type caches; actual named type definitions; queue
and object identities in snapshots. Trace marker currently used by temporary
creation remains unchanged. Do not keep production algorithms in test C.

Cases: enabled/disabled transitions and competing state bits, all animation
states, exact/strict frame deadlines and wrap, FPS 1/4/30/60, door direction wrap
and material/subclass sound priority, elevator heights around 0/20/32/64,
shape inclusion and size rejection, partner missing/present, positive movement
and exact Z, teleport animation counters/activation/target filtering/FX order,
trap-door CD deadlines and side effects, force quadrants/range/edge angles/RNG,
phantom distance and opposite direction. Add independent positive assertions for
script delivery, actual movement, force changes, queue mutations, and sound order.
Keep original nonnull pointer preconditions explicit.

Repeat original C full captures, lock and push before conversion. Port the family
and private helpers together, route direct Go owners, compare full captures and
all previous hashes, then qualify once at the batch boundary. Record C LOC after
completion, commit/push, summarize and continue.

## Audited scope

- 0053AC50: `char nox_xxx_updateDoor_53AC50(int a1)` (75 C lines).
- 0053B030: `void nox_xxx_updatePush_53B030(int a1)` (6 C lines).
- 0053B060: `char nox_xxx_updateToggle_53B060(uint32_t* a1)` (52 C lines).
- 0053B1B0: `char nox_xxx_updateTrigger_53B1B0(int a1)` (55 C lines).
- 0053B300: `char sub_53B300(int a1)` (14 C lines).
- 0053B320: `char nox_xxx_updateSwitch_53B320(uint32_t* a1)` (26 C lines).
- 0053B380: `char nox_xxx_updateElevatorShaft_53B380(int a1)` (30 C lines).
- 0053B410: `void nox_xxx_fnElevatorShaft_53B410(int a1, int a2)` (19 C lines).
- 0053B490: `void nox_xxx_elevatorAud_53B490(int a1, int a2)` (53 C lines).
- 0053B5D0: `void nox_xxx_updateElevator_53B5D0(uint32_t* a1)` (66 C lines).
- 0053B750: `void nox_xxx_elevatorFn_53B750(int a1, int a2)` (40 C lines).
- 0053B860: `void nox_xxx_updatePhantomPlayer_53B860(int a1)` (29 C lines).
- 0053BEF0: `int nox_xxx_updateTeleportPentagram_53BEF0(int a1)` (86 C lines).
- 0053C060: `void nox_xxx_fnPentagramTeleport_53C060(float* a1, int a2)` (11 C lines).
- 0053C0C0: `int nox_xxx_updateInvisiblePentagram_53C0C0(int a1)` (25 C lines).
- 0053C140: `void sub_53C140(float* a1, int a2)` (7 C lines).
- 0053C160: `void nox_xxx_updateBlow_53C160(int a3)` (51 C lines).
- 0053C240: `void sub_53C240(float* a1, int arg4)` (118 C lines).
- 0053DE80: `int* nox_xxx_updateTrapDoor_53DE80(uint32_t* a1)` (35 C lines).
- 00548830: `void sub_548830(int a1)` (9 C lines).
- 00548860: `void sub_548860(int a1, short a2)` (14 C lines).

## Locked original-C baseline

The baseline contains **1,979 cases / 14 groups**. Two complete original-C
captures match byte-for-byte (5.533s / 5.167s), and all prior **20,314** focused
contracts pass unchanged (58.421s). Hashes are locked in
src/world_mechanisms_porttest_test.go. Stable captures:
build/port-world-mechanisms/c-locked-{source,repeat}-world-*.json.

Tests cover all 21 entries, including named Trigger/PressurePlate timing,
FPS 1/4/30/60, wrap, strict/exact deadlines, sound material/subclass priority,
script events and caller/trigger identities, direction and queue transitions,
elevator height/distance bounds, actual indexed objects, teleport counter states,
class filtering, trap-door timers, and force directions/range. Independent
positive assertions verify script delivery, direct and owner-driven elevator
movement/height, visible teleporter and real player teleportation, and direct
and owner-driven directional force. Finite mass and real platform geometry are
supplied; zero-distance force behavior is also captured.

Optional World fixture state adds guarded collision storage and snapshots,
trigger/pressure-plate types, door and collision-list globals, and the original
absolute-value scratch relocation at 0x587000+55744 -> 0x5d4594+527672. All are
saved/restored. Filter-only generator-class objects own no generator children;
after snapshots, teardown clears that class bit before generic allocation
cleanup. Player movement uses the actual guarded player, not an inventory
object with a player class bit. Ordinary circle queries use ordinary objects;
missile cases remain in the advanced rectangle/force corpus.

The collision event recorder is aligned to 256 bytes because a retained char
return can contain the low byte of its function address. This stabilizes that
pointer-derived byte across linking and ASLR without changing production code.
Full callback identity, return values, queue state, object bytes, side effects
and guard checks remain captured. No production C algorithm is copied to tests.

Temporary/projectile conversion is pushed as `8f96bbc6`. Production C is still
**126,124 lines / 149 files / zero reference C**. Commit/push this baseline before
conversion. Next: translate the connected family and private angle helpers,
route Go callers, compare all full captures, qualify once at the boundary, update
C_LOC/docs, commit/push, summarize and continue.
