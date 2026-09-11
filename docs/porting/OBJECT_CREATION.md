# Object initialization and small death callbacks

GAME5 54C0C0–54CBB0 contains 11 functions and 340 physical C lines: monster
spell defaults; weapon/armor, obelisk, animation, trigger, generator and reward
initialization; ImpEgg, Polyp and Potion death callbacks. Stop before unrelated
player death-inventory cleanup 54CBD0. The auto-spell initializer has a live C
caller and a Go wrapper. The ten other callbacks are registered by Go and still
require their C ABI addresses.

## Original-C baseline

The shared callback fixture uses guarded C-owned actor, update, initializer,
use and health buffers, real type/weapon/armor/effect lookups, real balance
loads and object allocation, normalized pointers and restored global caches.
It captures all update-data mutations, 256 initializer bytes, 128 use bytes,
health words, return values, audio/deletion/creation order and both RNG indices.
Existing callback hashes pass unchanged after extending the fixture.

There are 2,048 generated cases, 11 smoke cases, 200 auto-spell cases,
640 equipment cases, 388 independent contracts, 84 cloud cases and one compiled
precision discriminator (3,372 total). Coverage includes cold/warm/partial/aliased
caches, type matches, absent modifier definitions, nil health, game modes,
16-bit FPS/durability wrap, byte ammo/charge writes, raw float32 bit patterns,
negative/NaN/infinite balance values and signed cloud FPS conversion.

Cloud allocation tests use valid types: removing a type is an invalid input to
the retained C-to-Go allocator bridge and panics rather than returning nil.
No allocation-failure behavior is claimed by these tests.

Compiled x87 rounds durability/staff balance loads to float32, multiplies in
double precision and spills each product to float32 before nox_float2int.
Staff doubling stays in double precision. The precision test uses MaxFloat32
and zero charges to distinguish this from overflowing the doubled multiplier
in float32. Assembly/captures are ignored under build/port-object-creation.

Hashes are locked in src/object_creation_porttest_test.go:

- creation-smoke: `369dea251575244f090a0f35caaf64ec5a61c05d549d39e5f87d928d9cafdda1`.
- creation-auto: `f181718a31e4bd9f548da39120ddd20e09eebc0479d7fea263fdcdf9ccb69e17`.
- creation-equipment: `0173bc94d782ea6c63cff2d7615360fbb721896b7dc930e671a72c84291c209d`.
- creation-contracts: `2a3b17723d0baad95c04351d2130ced551bdfb1b97eab6a586383a95e75be56c`.
- creation-cloud: `814d05ad41ae7560735566cb115907b4a72e2f3f1fd261af6c9c9d8143e1da25`.
- creation-corpus: `6bf58d5246c63bef4e22f9be339513b53df6ac01280234d1545235ec310b9ee9`.
- creation-precision: `cdbbad34bf6bf71b987296520e88aa6bc54a5b5b3158532305081cd002eff176`.

Production C remains 135,294 lines, 153 files, zero reference C.
Next: convert all eleven owners, compare the locked hashes, then qualify once.
