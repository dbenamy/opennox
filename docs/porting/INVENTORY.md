# Inventory insertion, removal, pickup and drop

## Connected scope in progress

33 original C functions / 1,421 physical C lines. Production is unchanged from
resource commit `0f74e2b3`: 130,474 C lines / 150 files / zero reference C.

The scope includes removal owner 4ED0C0; drop policy/dispatch/placement and chest
handling from 4ED290 through 4EE370 (excluding already native intervening roots);
insertion 4F3070; specialized pickups 4F3350..4F3DD0; weapon pickup 53A720,
Oblivion pickup 53A9C0, weapon drop 53AB10, armor pickup/drop 53E7F0/53EB70 and
droppable predicate 53EBF0. Neighboring cheat storage and equipment sweeps are
outside the scope. Exact caller audit and function list: ignored
build/port-inventory/scope.json and callers.txt, recoverable from these addresses.

## Original-C fixture work

inventory_porttest.go and inventory_porttest_test.go extend the existing shop
and resource fixture. The C dispatcher only invokes original functions; it
contains no replacement algorithms. Item use/drop callbacks record arguments
and controlled outcomes at retained service boundaries. Default pickup records
its inputs and optionally executes the actual insertion owner for integration
sequences. Full resource/player/object/init/use records and packet/protection
traces remain observable, alongside placement results, caches and callbacks.

Create-at capture treats existing pool-owned items as borrowed: moving an item
does not allocate it again, change its identity or transfer cleanup ownership.
This branch is confined to inventory cases; earlier fixtures retain their
previous behavior and locked hashes.

The original-C baseline covers 2,198 cases in 16 groups. Complete captures match
byte-for-byte across separate processes. The combined inventory/shop/resource
run passes in 34.072s with all 8,577 existing contracts unchanged; the focused
inventory run takes 5.030s. Hashes are locked in inventory_porttest_test.go.
No inventory production code has been converted yet.

Recovery: source build/baseline/env.sh, cd src, then run
`go test -tags porttest -run '^Test(Inventory|Shop|Resources)' -count=1 .`.
Set OPENNOX_CALLBACK_CAPTURE to an ignored absolute prefix to regenerate full
JSON captures; original algorithm sources are present in this baseline commit.
After conversion keep these hashes unchanged, qualify the connected batch,
update C_LOC/recovery checkpoints and commit/push. No C algorithm is retained
solely for tests after conversion.

### Expanded baseline review

The corpus now covers all 33 entry points, including real drop-all/insertion
sequences, ammo byte overflow and modifier matching, original food audio tables,
special item classes, armor replacement, and crown/treasure team objectives.
Owned-item links are explicit initial inputs. Borrowed player slabs are attached
to the isolated server so retained ownership/minimap services execute normally;
the runtime-only server handle is checked unchanged and normalized to an identity.
Minimap state covers all 32 player-info slots, including active slots without a
unit. Guarded init/use/health/update buffers detect writes beyond their allocation.
Crown/flag packet headers are filled by the actual serializer, so every transmitted
byte is compared; no new undefined-packet exclusions are needed.

Retained services: create-at records movement without reallocating existing
items, default pickup records admission and optionally calls real inventory
insertion, and the root player-state hook records requests. Oblivion's retained
pause effect is configured already busy, avoiding unrelated real-time/root-server
setup. Original food lookup tables occupy exactly 40 bytes each and are restored.

Geometry limitation found during original-C capture: the retained Go map-ray
implementation can loop while traversing from (100,100) to approximately
(317.9601,-34.2735), reached by radius 256 placement. The stack remained in
server/wall.go mapTraceRayImpl line 1055. This predates the inventory conversion;
random-placement boundary cases use origins far enough inside the map to keep
these rays valid. Track the negative-coordinate traversal issue separately;
do not silently change the ray service as part of this behavior-preserving port.
Pure shape-radius cases still cover signed zero, subnormal values, infinities,
and quiet/signaling NaNs. The 2,198 cases in 16 groups match byte-for-byte across two separate original-C runs. Hashes are locked in inventory_porttest_test.go.
