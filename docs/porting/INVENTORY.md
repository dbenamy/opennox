# Inventory pickup and drop port

Completed 2026-09-11: 33 functions / 1,421 physical C lines removed.
Production C: **129,053 lines / 150 files / zero test-reference C**.
Original-C baselines are pushed as `0763837a` and `d71388bd`.

## Scope and implementation

The owners now live in legacy/inventory.go, inventory_drop.go,
inventory_placement.go, inventory_pickup.go, inventory_equipment.go and
inventory_objectives.go. inventory_exports.go retains 31 ABI entry points for
C callers and callback addresses. Private eligibility and shape-radius bridges
are retired; their fixture actions call Go directly. Native pickup/drop
registrations and callers in shop, trade, spawn and monster debris code call Go
directly. Unknown C drop callbacks still use the existing ABI.

Converted addresses:

- src/legacy/GAME3_3.c: 004ED0C0, 004ED290, 004ED500, 004ED580, 004ED5E0, 004ED710, 004ED790, 004ED810, 004ED930, 004ED970, 004EDA40, 004EDCD0, 004EDDE0, 004EDE50, 004EDF00, 004EE2A0, 004EE370, 004F3070, 004F3350, 004F3400, 004F34D0, 004F3510, 004F3580, 004F3B00, 004F3C60, 004F3CE0, 004F3DD0.
- src/legacy/GAME4_3.c: 0053A720, 0053A9C0, 0053AB10, 0053E7F0, 0053EB70, 0053EBF0.

Neighboring cheat storage, equipment sweeps, table initialization, and retained
map/equipment/root services are outside this batch. The next connected candidate
is the weapon/armor equip/dequip family and its helpers.

## Original-C contracts

The fixture reuses the shop/resource player, object, item, network and protection
snapshots. All 2,288 cases in 17 groups match complete original-C capture files
byte-for-byte. Coverage includes linked-list insertion/removal and ownership;
weight/protection updates; drop admission and rejection; use/pickup outcomes;
default and specialized drops; ray/placement/chest sequences; drop-all fallback;
ammo merging/byte overflow/modifier matching; food audio tables; equipment and
armor replacement; crown, treasure, and team membership; shape radius including
signed zero, subnormals, infinities, and quiet/signaling NaNs.

Item init/use/health/update buffers are guarded. Minimap state covers all 32
player-info slots, including active slots without units. Complete packets,
player/item memory, callback/RNG/deletion/creation effects and caches are compared.
No new packet-byte exclusions were introduced: the actual serializer initializes
both leading bytes of the crown/flag messages.

The first 2,198 original-C cases were repeated before conversion. A native
mismatch exposed an incorrect simplification of team comparison: 419180 validates
list membership as well as the numeric ID. The port calls that retained service.
The additional 90 actual-membership cases were captured and repeated in an
isolated original-C checkout before their hash was locked in `d71388bd`. The
original numeric-ID-without-membership cases remain unchanged.

The fixture position now uses C-owned memory, preserving all locked outputs
while permitting exported C functions to return its address. Borrowed player
slabs are attached to the isolated server; runtime server handles are checked
unchanged and normalized as identities. No cgo checks are disabled. Owned-item
links are explicit fixture inputs, not setup-time ownership notifications.

Retained service boundaries: create-at records movement without reallocating
borrowed items; default pickup records admission and optionally invokes real
inventory insertion; the root player-state hook records requests. Oblivion's
retained pause effect is configured already busy, avoiding unrelated time/root
setup. Original food tables occupy exactly 40 bytes each and are restored.

Preserved details include callback/list ordering, ammo byte wrapping, armor
replacement priority, full minimap/player state changes, drop-all spiral state,
food-table priority and floating-point spill points. Placement retains libc
sin/cos as mathematical primitives, not C placement algorithms. Drop-all's C
fallback `v1 + 7` uses float2 pointer arithmetic: position at byte 56.

## Qualification

- All 10,865 inventory/resource/shop/trade contracts pass in 37.646s, with all
  prior 8,577 hashes unchanged. Final inventory captures match original C exactly.
- Accumulated default/server/highres port tests pass in 89.586s / 90.305s / 88.066s.
- All three production builds pass: ELF32/Intel 80386, GO386=sse2, CGO enabled.
- Full-suite failure multiset is exactly unchanged: 1,553 entries; 15 pass,
  3 fail and 32 skipped packages.
- Fresh inventory-port gameplay passes in 36.485s against unchanged repeat-a
  goldens, overrides disabled, Xvfb and null audio.

Reproduce with build/baseline/env.sh, then from src:
`go test -tags porttest -run '^Test(Inventory|Shop|Resources)' -count=1 .`.
Original sources and fixtures are recoverable at the baseline commits. Stable
local captures are build/port-inventory/c-members-inventory-*.json and
native-final-inventory-*.json. Qualification scripts/results live in that ignored
artifact directory; full-suite logs may contain secrets, so report metadata only.
No C algorithm remains solely for tests.

## Separate retained map-ray issue

Original-C testing exposed a loop in the retained Go map-ray traversal when a
ray from (100,100) crossed into negative coordinates, approximately
(317.9601,-34.2735), during radius-256 placement. The stack remained in
server/wall.go mapTraceRayImpl, called by MapTraceRayAt. This predates inventory
conversion. Random-placement cases use origins far enough inside the map to keep
rays valid; drop-all fallback cases with out-of-map origins also pass. Track the
negative-coordinate traversal issue separately rather than change map semantics
inside this behavior-preserving port. It remains unresolved.
