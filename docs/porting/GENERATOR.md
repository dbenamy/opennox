# Monster generator policy

Scope: nine functions in GAME5: death 54E630 and update, player selection,
placement, vacancy enumeration/candidate, radial search, spawn, and creature
copy at 54E930–54F2B0. Together they occupy 468 physical C lines. Four shared
visibility helpers between death and update remain outside this batch.

## Original-C contracts

The baseline contains 2,264 generator cases plus 128 tile-fixture known answers.
The case groups are 3 placement smoke, 768 placement, 45 vacancy, 600 update
gates, 4 object smoke, 48 death, 280 player selection, 72 spawn, 64 creature
copy, 192 integrated update/spawn, 60 definition-health, and 128 inventory copy.
Hashes are locked in src/generator_porttest_test.go before C removal. Two
independent C runs match; locked generator and existing callback, creation,
penalty and death regressions pass together (11.291s).

Fixtures use real spatial indexing, guarded tile rows and definitions, ray
tracing, player records, balance lookup, object allocation, SpawnClass lists,
inventory insertion, modifier copying, NPC equipment, script recording, audio
queues, and packet encoding. Captures include complete 164-byte generator and
2200-byte creature records, 772-byte C object layouts, health, player state,
created objects, spawn references, caches, guards and both random indices.
Fixture-owned resources and prior server/blob state are restored.

Placement tests distinguish successful placement from all 32 attempts blocked;
192 full-blocker cases fail and the other 576 succeed. A radius-50 circle at the
source covers the search while fitting the spatial index's four partitions.
Resized objects are unlinked and reindexed. Tile checks include cold caches,
water/no-teleport definitions, and points outside the real grid interior.

Update tests vary stage, rate class, cache initialization, flags, frame, capacity,
and packed source arrays of length 1–4 at all three levels. Integrated cases
produce exactly 12 real spawns; admission limits can suppress creation after
selection. Health tests cover default type health, signed definition health,
zero clamping, fractional/negative scales, overflow, infinities and NaN.
Inventory tests verify every copied modifier word, equipment dispatch, reversed
insertion order and replacement of a previously equipped weapon.

Pointer values receive stable identities. Update's char return truncates the
spawned pointer: the fixture first verifies the actual low byte, then records a
stable marker. Definition pointers and generator source pointers overlaid on
shared fixture fields are normalized as well. No generator algorithm is retained
in C solely for tests after conversion.

## Arithmetic and scope limits

The original radial assembly stores the incremented angle as float for sin but
passes the unrounded double sum to cos. Preserve this asymmetry, the double
1.8849558 increment, float result stores, and Logic RNG consumption. Existing
libm primitives may remain C calls. The original vacancy and player filtering
contain numeric conversions of float views of class/flags; preserve those
operations rather than silently repairing the decompiled behavior.

The retained SpawnClass allocator, glyph creation, modifier effects, inventory
and equipment engines are dependencies, not converted algorithms here. Tests
exercise successful spawn registration and admission rejection; they do not
exhaust the SpawnClass pool or simulate impossible factory type IDs. Source
arrays are packed as expected by the original routine. Synthetic modifier
callbacks are nil; their copied descriptors and real equip policy are observed.
