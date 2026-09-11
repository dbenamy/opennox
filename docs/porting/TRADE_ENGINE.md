# Remaining server trade engine

## Batch boundary

The next connected batch is all of server__system__trade.c: 18 functions and
815 physical C lines. Six packet-decoder roots start a session, add an offer,
buy one/by type, quote a sale and sell one. Twelve private helpers handle stock
removal, cached shop creation, introduction/offer packets, offer admission,
stock-definition counts, shortfall messages and quest gem policy.

The shop core is native and pushed as `d8892083`. This batch reuses its pools,
players, inventory/gold services, packet capture and definition fixtures. The
remaining engine algorithms are still C. The caller audit is ignored locally at
build/port-trade-engine/callers.json. Inspect previously retained shop-core ABI
roots again after conversion to remove those whose final C callers disappear.

## Locked original-C baseline

New fixture: legacy/shop_engine_porttest.go; new contracts:
src/shop_engine_porttest_test.go. The expanded corpus has 748 cases in twelve groups: stock removal and repeated
misses; four-type offer admission; packet boundaries and all modifier slots;
cold quest caches and vendor counts; offer mutation; single/bulk purchases;
quotes/sales; opening; food/staff limits; cached reopening; real 500-node pool
exhaustion/reuse; and first/second-player already-trading rejection. All twelve groups repeat byte-exact. The 748-case C run passes in 8.647s;
all 4,238 existing shop-core contracts plus the locked engine cases pass
together in 20.147s. Hashes are in tradeEngineHashes; captures and repeat
comparison metadata stay under build/port-trade-engine.

Each engine step captures the existing full session/node/object/player state,
plus the vendor's complete 1,724-byte record and twelve engine cache words.
All three sparse player records are preserved, including the third rejection
participant. Actual text-message queues for indices 1, 7 and 31 are captured.
A real immutable StringManager table exercises localized message formatting.
Engine cases permit intended input mutations while preserving slab guards;
query-only shop cases retain their original byte-exact read-only checks.

Two introduction packets transmit unused stack tails in C. For the 86-byte
shop introduction, exclude only bytes after the terminating byte string at
payload offset 54. For the 52-byte peer introduction, exclude only bytes after
the terminating UTF-16 name at offset 2. All defined name bytes, opcodes, sizes,
recipients, ordering, queue state and other packet payloads remain observed.
The exclusion is applied to both the shop-step and outer lifecycle packet
snapshots. The native implementation will initialize those formerly
indeterminate tails.

Offer admission's original four-slot arrays assume at most four distinct item
types in an existing offer list. Cover all valid list sizes, matches, modified
items and rejection of a fifth type; do not feed overflowing invalid lists to
the C baseline. This does not relax the public add-offer capacity checks.

No production C has been removed in this batch yet. Current count remains
131,935 physical C lines / 152 files / zero test-reference C lines.

The pickup dependency is a recording C callback: every invocation preserves
unit/item identity, both arguments and order. Its algorithm is outside this
batch. Single purchases bypass pickup for food/reward classes, whereas bulk
purchases call an available pickup function; both paths are covered. The real
inventory and gold implementations run. Bulk affordability deliberately tests
the original one-time gold snapshot across multiple purchases.

Player animation state is another retained dependency. Engine-only cases record
its requested unit/state and return success; they do not model animation. The
real unit-freeze, raise, status-packet and path-reset services still execute.
The hook is restored after each case. The missing full game server initially
caused this dependency to panic; no production behavior was changed to fix it.

Cached reopening exposed a fixture ownership issue: a released session/object
address may be reused. Adoption now checks live ownership, while maintaining
stable identities and original lifetime observations. All fixture corrections preceded the engine hash lock. The existing 4,238
shop-core hashes remain unchanged and pass alongside this new baseline.

Independent assertions cover repeated stock-removal misses, empty-offer
admission, first opening and rejection, exact sale/repeated-sale gold, cached
session identity, and pool exhaustion followed by successful slot reuse.

The first combined capture run filled the VM disk with duplicate diagnostics.
Old reproducible Go cache files (over twelve hours old) were removed; the asset
archive and gameplay evidence were preserved. The full combined run then
passed without duplicate capture output. This was an artifact-write failure,
not a test-behavior mismatch.

## Conversion next

Replace all eighteen engine bodies and remove the now-empty C file. Preserve
six packet-decoder roots and retire eight shop-core exports whose last C
callers disappear. Route internal and fixture calls directly to Go. Re-run all
locked cases, accumulated variants, production builds, full-suite comparison
and fresh headless gameplay once for the completed connected batch.
