# Shop and trade port

## Scope and baseline checkpoint

The connected batch is GAME4_1.c 50E2A0 through 510E20, ending before 510E50:
36 functions / 1,337 physical C lines. Session/item pools, prices, stock,
offers, gold balancing, packets, completion/cancellation, repair and sales share
this boundary. Retained trade-engine and packet-decoder callers stay outside.
No production C has been removed in this baseline checkpoint.

The first 3,389 original-C cases are locked across thirteen byte-exact repeated
groups in src/shop_porttest_test.go. Their locked contracts and adjacent
generator/spawn/penalty/callback/creation/death tests pass together (25.396s).
Original captures and repeat logs are in build/port-shop/c-fourth-* and
c-repeat-fourth-*. Raw pointer-bearing diagnostic artifacts stay ignored.

Coverage so far:

- 384 base prices, 960 modifier/health prices, 576 ammo/wand charge prices,
  140 reward/gem prices and 384 floating-point rounding cases. Inputs include
  unsigned high-bit worth, signed modifier prices, float32 integer boundaries,
  nonfinite multipliers, zero capacities, quest flags, gem caches, both vendor
  orientations, missing and valid guide/spell definitions, all price modes.
- 288 stock-match/index cases and 448 category/price priority keys. Includes
  duplicate matches, exact type width, each modifier slot, reward forms,
  missing records and every priority-table category.
- 13 session sequences and 27 stock sequences, with head/middle/tail removal,
  insertion order and ties, misses, bulk reset, destruction and slot reuse.
  One additional 501-item sequence reaches the real 500-node capacity and
  checks failed insertion, whole-list release and successful reuse. Session
  sequences reach the real 64-session capacity (128 gold objects).
- 64 packet sequences capture all eight packet helpers, exact bytes, length,
  recipient, ordered flag, sequence numbers and queue order.
- 80 gold-balance and 24 completion/withdrawal/cancellation sequences capture
  intermediate session/list/object/player state, gold and protected-gold
  records. The original inventory and gold services are exercised directly.

The fixture layers on the existing callback/player/network infrastructure.
It owns guarded C query records, real trade pools and a 768-object factory pool,
and restores the prior pools, tables, players and server definitions. Full C
object prefixes and relevant input slabs are normalized and compared. Large
capacity sequences hash complete normalized object snapshots to keep artifacts
bounded, while preserving explicit session/node topology. Unused vendor health
and update pointers were removed from the fixture after repeat checks caught
address-dependent snapshots. Query records are byte-for-byte read-only.

An existing lifecycle quirk is preserved: bulk reset frees trade pool records
and shop stock, but does not release the two gold objects per session. Tests
observe the resulting live-object count and explicitly clean up those orphaned
allocations afterward. Ordinary session destruction does free its gold objects.
Allocation-class creation failure is not injected; the retained Go allocator
panics there. Fixed-pool exhaustion is tested through the real allocator.

## Remaining work before conversion

Extend the C baseline for stock loading (including quest generation), repair
quotes/repair/sales, quest-cached sessions and more overlapping offer sequences.
Add independent assertions for gold/inventory outcomes alongside exact hashes.
Commit those baselines before replacing the production C bodies, then qualify
one combined native batch. C remains 133,272 physical lines / 152 files, with
zero reference C.

Primary external-caller inventory is build/port-shop/callers.json. Retain both
shopExit and tradeAccept exports: the first has packet-decoder and spell-owner
callers; the second has a packet-decoder caller. The helper audit initially
omitted those two rows. An ignored price draft also needs primary corrections:
quest modifier/sell products stay double until their original float32 spills;
quest-protected and NaN prices follow the C clamp; do not invent a string length
cap. Do not integrate that draft without the locked contract comparison.
