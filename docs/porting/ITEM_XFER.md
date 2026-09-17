# Item and reward serialization

## Qualified result

Twelve GAME4.c callbacks are now Go implementations: spell/ability rewards, field
guides, weapons, armor, ammunition, team items, gold, obelisks, toxic clouds,
monster generators and reward markers. The conversion removes **1,065 C lines**,
leaving **58,539 physical lines / 82 production files / zero reference C**.
The unrelated prefix and two declarations following the block are unchanged.

All **24 frozen groups / 2,931 state records** match: **1,757 item records** plus
**1,174 common/world-object records**. Default/server/highres each pass **178 roots /
3,294 leaf cases**, with no skips. All three production binaries and ABI audits
pass. The full asset suite preserves exactly **1,553 known failure entries** and
package results **15 pass / 3 fail / 32 skip**. Gameplay matches **41 frames**,
actual save/load **seven frames**, and flat rendering **14 frames** with exact
War01A regeneration. No new failure is accepted or suppressed.

The twelve existing callback symbols remain Go-backed C exports, preserving the
original mixed int/int-pointer declarations and function identities. Go code uses
the existing common/inventory stream helpers and real modifier/resource/object
owners directly. No C algorithm remains solely for testing.

## Recovery and evidence

- Prerequisite fixes: **bf1afd10**, committed and pushed before freezing.
- Frozen C baseline: **e4852239**, committed and pushed before translation.
- Expectations: test-local hashes plus [item-xfer-batch.json](item-xfer-batch.json).
- Selection: [item-xfer-affected-tests.txt](item-xfer-affected-tests.txt); item tests
  are also included in the accumulated milestone pattern.
- Final artifacts: build/port-item-xfer/native-final-{default,server,highres,production}.
- Final source proof: native-index-proof.json covers **1,912 staged source files**,
  including all affected phases and production. Static memory checking passes.
- Client SHA-256: 657f95d64f1c78e25cb8ee5b1d3b30f3c85ff38e58484b2b87a3d53c132ff4d7.

The C default/server/highres captures matched before pinning expectations. Fresh
frozen runs each passed **38 roots / 3,104 leaf cases** without skips and matched
all 24 groups. C-source proof covered 1,908 files. The C client also matched the
preceding gameplay/save-load/flat references before replacement. Its SHA-256 is
87b6326a0043ee952c075d5fbfeb4bb76fde0cdec20883d4117f3c1aa0607ea9.

Final affected test-driver seconds: **101.587 / 100.195 / 111.189**. These include
bounded concurrent qualification and are not runtime performance claims. An earlier
native run also passed all groups before the fixture-only tooling adjustment below.
Final production binaries are byte-identical to those earlier builds.

## Contracts and notable behavior

Fixtures use registered callbacks, real object/type pools, modifier definitions,
player bitsets, and drawable/minimap owners. Snapshots include the full object,
type-owned buffers, health data, modifier identities, nested children and stream
checksums. Only identified pointer fields are normalized; independent contracts
check actual pointer identity, ordering, ownership, bytes and stream positions.

Expanded item coverage includes:

- **336 charge-policy cases:** historical versions, ordinary/special wand subtypes,
  quest validation, count/capacity boundaries and signed charge amounts.
- **792 health-policy cases:** weapon/armor versions, solo/switch/quest/player modes,
  present/missing definitions, clamping and low-word definition durability.
- **64 modifier cases:** every four-slot mask, real descriptor identity, empty and
  255-byte names, attribute-tail mutations and team-position copying.
- **96 obelisk cases:** client flags, static/dynamic/missing sprites, minimap lists
  including non-head matches, and mana boundaries.
- **17 generator cases:** sparse/dense child lists, row compaction, cipher section
  framing, identity/order, lifetime, allocation counts and partial failures.
- **230 reward cases:** every valid ID, mixed/noncanonical masks, invalid names up
  to 255 bytes, rejected ID zero, duplicates, retained bits and versioned tails.

Historical early returns retain their observed lifetime behavior. Stream checksum
updates depend on operation boundaries, including zero-length calls. An empty-named
modifier descriptor writes the same bytes as an absent modifier but takes an extra
zero-length write; snapshots preserve the resulting checksum distinction.

Reward counts include only mask value 1, but the writer emits names for every
nonzero value. Preserve that established behavior. Noncanonical writers have
separate contracts rather than an assumed round trip. Reads merge existing bits,
retain earlier accepted names on later failure, and stop at the rejected name.

Sparse generator children compact within each row on reload. Earlier successfully
loaded children remain owned when a later child fails. The empty obelisk callback
still has two objective-update callers, so its shared C definition/header remain;
the translated serializer omits its no-op call.

The affected selection covers item/common serialization, object creation/state,
resources, rewards, equipment, shop engine, map population/painting and minimap.
The existing opt-in map-population diagnostic is explicitly excluded; every selected
test must execute without skips. The preceding common-serialization batch already
ran the full accumulated milestone, so this batch uses the established affected
selection policy plus complete production qualification.

## Prerequisite corrections

The original generator left a newly allocated child unowned when its callback
rejected the record. The before-change contract observed two live objects instead
of the parent alone (c-ownership-before.log). Release that rejected child before
returning failure, matching the preceding inventory correction.

Pre-version-11 weapon records initialized only sixteen of twenty attribute bytes.
A real charged wand copied the uninitialized final word. Four old-version cases
failed (c-attributes-before2.log), while version 11 supplied two 0xffff words and
passed its control. Define the old tail as the same 0xffffffff. Normal nil-modifier
items retain their skip-copy behavior; normal armor does not enter the charged-wand
path and received no C correction.

These two reversible fixes added two C lines before baseline freezing. Initial
prerequisite default/server/highres runs each passed 29 roots / 1,569 leaf cases.
See [DECISIONS.md](DECISIONS.md) for the recorded rationale.

## Fixture and tooling corrections

The first old-attribute fixture used ClassWeapon and missed the charged-wand branch;
its failure was not engine evidence. The corrected fixture uses ClassWand. Other
initial setup corrections explicitly assigned health after alloc.New (which uses its
argument as a type hint and zero-fills) and owned/restored the actual ability-name
table. None changed an expectation to accommodate translated behavior.

The first nonempty generator writer fixture used XOR mode for four-byte framing.
That alternate mode could not backpatch its section length through the writable
binfile seek adapter. The fixture now uses the actual map cipher mode. Independent
contracts check aligned lengths and decoded fields; a full wire hash also pins
compatibility padding. Shared stream behavior was unchanged.

The first flat-C launch used the wrong local compressor path and stopped before
starting the game. The corrected launch used the existing qualified compressor.
Initial native compilation required the correct ability-resolver header and exact
export/header signatures. These were corrected before behavior qualification.

The first production full-suite run found one new tooling failure: noxfactor's
legacy C-cast tokenizer rewrote the fixture type `func(int) string` into the invalid
`funcint(string)`. An isolated copy reproduced it. Naming the parameter as
`func(index int) string` preserves the Go type and fixture behavior while avoiding
that ambiguous token sequence. Both refactoring-tool packages pass; final affected
and production gates were repeated, with all frozen expectations unchanged. A
general tokenizer redesign is deferred for later review.

Completed scenario copies retain screenshots, logs, saves, changed maps and verified
asset-restoration manifests. Duplicate unchanged assets and identical production
binaries share or release local storage; original assets remain untouched.
