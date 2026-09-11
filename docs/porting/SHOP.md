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

## Expanded original-C baseline

The complete batch baseline now locks **4,238 cases in 21 groups**. All new
captures repeat exactly; the original thirteen hashes remain unchanged.
Added 144 repair sequences, 96 sales sequences, 16 quest-cached session
sequences, 129 stock-loading cases (including quantity 255), 168 quest-loading
cases, 24 multi-node offer removals, 256 quest rounding cases and 16 stock
count/index boundaries up to the full 60-entry vendor record.

Independent assertions cover repair current/max health, sales gold credits,
completed peer-trade gold and inventory delivery, exact packet metadata,
stock counts, first/last indices, pool limits and repeated withdrawal misses.
Repair tests include insufficient funds, full/damaged/overfull/zero health,
empty/full wand charges, missing codes and repeated quotes. The retained
repair behavior is preserved; the fixture does not add an affordability gate.

Stock loading uses real type factories and registered spell/ability/guide Xfer
identities. Quest loading retains the real reward generator with controlled
category and spell-eligibility tables, stage boundaries/wraparound, missing
named items, the ankh cutoff, marker presence and probability variants.
Spell/ability rewards are enabled; weapon/armor/guide reward categories return
no generated object in this fixture. Their non-quest stock creation is covered,
and the reward-generation dependency itself remains C.

The original stock loader copies a five-word modifier array after initializing
only its first four words. The fixture excludes exactly that uninitialized
fifth word when it was copied. All defined modifiers, object data, list order,
prices, packets and RNG state remain captured. The native loader will initialize
this formerly undefined word to zero, as the earlier callback port does.

Factory item data, health, full player state, packet recipient/ordering fields,
protected-gold records, cached sessions and object lifetimes are captured at
each action. Fixture-created health/init allocations and delayed reward markers
are explicitly reclaimed. Query input records remain byte-exact read-only.
Large capacity cases retain full-object digests and explicit list snapshots.

Complete C run: 11.741s. Expanded locked shop plus adjacent generator/spawn/
penalty/callback/creation/death regressions: 27.405s. Artifacts: c-locked-final-*,
c-expanded-*, c-final-* and locked-expanded-adjacent.log under build/port-shop.

## Conversion and qualification next

Replace all 36 C bodies together, keep required C ABI roots, and route private
helpers and Go callers directly to Go. Then run the locked contracts, accumulated
variants, production builds, exact full-suite comparison and fresh headless
scenario once for this connected batch. Production C remains **133,272** physical
lines / 152 files, with zero reference C until the conversion is applied.

Primary external-caller inventory is build/port-shop/callers.json. Retain both
shopExit and tradeAccept exports: the first has packet-decoder and spell-owner
callers; the second has a packet-decoder caller. addItemToShopSession is private
once this whole interval is converted. The helper audit got these rows wrong.
The ignored price draft needs primary corrections: quest modifier/sell products
stay double until their original float32 spills; quest-protected and NaN prices
follow the C clamp; do not invent a string length cap. Do not integrate the draft
without the locked contract comparison.
