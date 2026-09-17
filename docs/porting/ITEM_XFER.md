# Item and reward serialization

## Scope

Continue after qualified common/world-object serialization (b4ff2519). Move the
twelve adjacent GAME4.c callbacks from 004F5F30 through RewardMarker 004F74D0:
spell/ability rewards, field guides, weapons, armor, ammunition, team items, gold,
obelisks, toxic clouds, monster generators and reward markers. The initial scope
is **1,063 C lines**. Preserve the two unrelated declarations following the block.
Current qualified C: **59,602 lines / 82 production files / zero reference C**.

All callback symbols remain registered or compared by existing callers, including
shop and root transfer code. Keep their identities as Go-backed exports. Reuse the
qualified common/inventory stream helpers and real object/type ownership fixtures.
No C algorithm will remain solely as a test reference.

## Baseline and qualification plan

Use registered callbacks with real object pools, type templates, stream owners,
modifier descriptors and client/minimap owners. Test current exact bytes and
round trips, signed/historical versions, name boundaries, payload mutations,
stream positions, lifetime restoration and failure ownership. Add detailed
contracts for weapon charges/HP policy, modifiers, nested generator children,
reward masks/names and obelisk minimap state. Normalize only identified pointer
fields; check pointer identity and order independently.

Repeat C captures before freezing expectations. The preceding common serialization
batch already completed the full accumulated milestone. Use an audited affected
selection across default/server/highres for this batch; broaden if shared changes
or failures warrant it. Qualify all three production builds/ABIs and normal,
actual save/load and flat rendering/map regeneration. Preserve the exact known
full-suite failure set when running production qualification.

## Audit findings to resolve

- The generator appears to leave a newly allocated child unowned when its transfer
  callback rejects the record. Demonstrate with the real callback/pool before
  selecting the same small release-on-rejection correction used by inventory.
- Weapon/armor version<11 initializes only sixteen of twenty attribute scratch
  bytes. Ordinary nil-modifier items skip copying them, but charged weapon
  subclasses may copy the uninitialized tail. Add a discriminating C contract
  before deciding how to define it; do not freeze random stack bytes.
- Reward-marker counts include exactly-one bytes, while its writer emits every
  nonzero entry. Preserve the observed format unless an independent correction
  is justified; do not assume noncanonical masks round-trip.
- Stream checksums depend on operation boundaries. Keep field grouping and
  zero-length operations aligned with C, as in the preceding batch.

## Current progress

Prerequisites are committed and pushed as **bf1afd10**. All twelve item callbacks
remain in C. Expanded contracts are being qualified before repeated captures and
baseline freezing. Ignored build/port-item-xfer contains consumed fixture drafts
and an uninstalled translation draft; these are not accepted implementation.

## Prerequisite evidence

The initial generator regression fails against C: a rejected child leaves two
live pool objects instead of the parent alone. Evidence:
build/port-item-xfer/c-ownership-before.log.

The old-attribute fixture first selected ClassWeapon, which does not enter the
charged-item copy branch; that was a fixture-selection failure, not evidence of
uninitialized copying. With the actual ClassWand bit and charged subtype, all
four old-version cases copy a stack-derived word into the attribute tail
(c-attributes-before2.log). A version-11 control is being checked to confirm its
existing two-0xffff default. The intended correction defines that same default
for pre-version-11 wand records. Ordinary items retain their current skip-copy
behavior. Armor uses a similar scratch buffer but does not take the charged-wand
copy path for its normal class; no armor change is currently proposed.

These are small, reversible prerequisite corrections, recorded for review before
baseline freezing. No item callback has yet been translated to Go.

The version-11 control passes for charged and ordinary items. Two C lines now
release the rejected generator child and initialize the old wand attribute tail.
Both regressions pass in the first after-change run. Other initial failures were
fixture setup: alloc.New uses its argument as a type hint and zero-fills, and the
C ability-name resolver needs its actual pointer table initialized. The fixtures
now assign health values explicitly and own/restore the ability table using the
same production names as the existing quickbar fixture.

New item snapshots include complete health records, named modifier identities,
normalized nested-child pointers and stream checksums in addition to the full
object and type-owned buffers. This makes operation grouping observable. These
captures are still development data, not frozen expectations; qualification is
in progress.

## Qualified prerequisite checkpoint

Default/server/highres each pass **29 roots / 1,569 leaf cases**, without skips:
272 new item cases plus the 1,297 qualified object-transfer cases. Test-driver
seconds are 29.208 / 115.061 / 37.825. All phases record unchanged source; static
checking passes (static-prerequisite.log). The three new development capture
groups match across targets: current defaults 12, historical 186, name boundaries
24. They are not yet the final frozen item oracle.

The two C corrections add two lines, leaving **59,604 production C lines / 82
files / zero reference C**. The scoped item callback block is now **1,065 lines**.
This qualifies the prerequisite and initial contracts; modifier/charge/HP policy,
nonempty generator ownership, reward-mask/name and obelisk-minimap coverage, full
repeated C captures and current-source integration remain due before replacement.

## Expanded contracts and review notes

Development C checks pass 336 charged-wand policy cases, 792 weapon/armor health
cases, 64 modifier-mask cases, 96 obelisk/minimap cases, and 17 generator cases.
These exercise real modifier definitions, player bitsets, drawable/minimap owners,
object allocation and nested callbacks. Reward-mask/name coverage is next.

The first nonempty generator writer fixture used XOR mode to obtain simple
four-byte sections. That mode could not backpatch its section length through the
existing writable-binfile seek adapter. The fixture now uses the actual map cipher
mode. Independent contracts check aligned section lengths and decoded fields;
a separate wire hash retains compatibility padding. Sparse children compact within
each row, retain row order and allocate distinct owned objects on reload. Failure
cases check that earlier accepted children remain owned and rejected children are
released. No shared stream behavior was changed to accommodate this fixture.

Review later: modifier descriptors with empty names serialize like absent modifiers
but take a zero-length writer operation, producing a different checksum. Captures
retain both read and write checksums so the port preserves that distinction.

## Frozen C baseline

All three targets match **1,757 new item records / twelve groups** and the existing
**1,174 object records / twelve groups**. Expectations are now pinned both in tests
and the batch manifest. Fresh frozen default/server/highres runs each pass **38
roots / 3,104 leaf cases**, with no skips and unchanged source. Static memory-access
checking passes. Artifacts: c-expanded-{default,server,highres},
c-frozen-{repeat,server,highres}, static-c-frozen.log.

Reward coverage adds 230 cases: every valid ID, mixed masks, invalid names through
255 bytes, rejected ID zero, partial mutations, duplicate names, retained existing
bits and versioned tails. Noncanonical writers are checked separately from readers;
counts include only value 1 while names include every nonzero value.

The C production client matches the prior 41-frame gameplay and seven-frame actual
save/load references. Flat rendering/map regeneration qualification is in progress.
Its first launch used an incorrect local compressor path and stopped before starting
the game; the corrected launch uses the existing qualified compressor. No engine
change was needed.

The conversion selection contains 161 root tests: item/common serialization,
object creation/state, resources, rewards, shop engine, map population/painting and
minimap. The existing opt-in map-population diagnostic is excluded explicitly;
all selected tests must run without skips. This follows the established affected
selection policy after the previous full common-serialization milestone.

The empty obelisk callback still has two objective-update callers. Keep its shared
C definition/header while eliminating the now-unnecessary call from the translated
obelisk serializer. This batch still removes exactly 1,065 C lines.

Flat replay now passes **14 frames** and exact War01A regeneration. Current C
client SHA-256: 87b6326a0043ee952c075d5fbfeb4bb76fde0cdec20883d4117f3c1aa0607ea9.
Frozen-phase source proof matches all **1,908 staged source files**. The baseline
is ready for its recovery commit; conversion and native qualification follow.
