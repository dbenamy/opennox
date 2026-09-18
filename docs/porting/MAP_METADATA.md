# Map metadata sections

## Scope

Close the remaining C map-section dispatch boundary: MapInfo, AmbientData,
ObjectTOC, and the two connected ambient getter/setter services. Five bodies,
142 physical body lines. This deliberately smaller batch reuses the preceding
map-section IO and existing real renderer/object owners rather than padding the
scope with unrelated code. Parent: `90fbd480` (27,001 C lines /67 files).

## Contracts

- MapInfo: versions 0–3, rejected future/high-bit versions, fixed field bytes,
  historical defaults, both strings at lengths 0/1/15/31, preserved tails and
  surrounding bytes. The historical default path writes only the first byte of
  each string. Current writer always selects version 3.
- AmbientData: signed version boundary (positive future versions accepted), full
  32-bit channels, renderer propagation under the actual game flag mask, exact
  bytes, surrounding storage, and renderer state when propagation is disabled.
- ObjectTOC: signed version acceptance, server/client lookup routing, empty,
  unknown and repeated names, embedded terminators, maximum byte length, code
  zero/high-bit values, clearing and preservation on rejected versions. Writer
  contracts cover actual object, inventory, pending and missile lists, repeated
  types, encounter-code assignment and type-index serialization order.

All IO goes through the actual plaintext cryptfile fixture; captures include its
checksum and file position. Preserve original individual IO calls: the existing
checksum depends on call boundaries. No algorithms are copied into the fixture.
Only valid bounded MapInfo strings are supplied; malformed/truncated input
semantics are outside this translation's behavior contract.

## Qualification status

C baseline qualified: five focused roots /478 test entries, five repeated captures
/473 records, identical on default/server/highres. Count contracts additionally
cover 0/1/32767/32768/65535 entries. Static mapped-memory checks pass. All 2,478
source files are identical across final frozen gates. Every parent source file
is unchanged; exactly five additive porttest-only files permit reuse of the
qualified parent production binaries/scenarios. See
[qualification](map-metadata-c-qualification.json).

No C prerequisite correction was needed. Initial writer expectations used the
wrong numeric value for GameFlag22; checking its declaration corrected the test
and added both neighboring flag values before freezing. Six C-to-Go adapters
become unreferenced with this conversion and will retire with the five C bodies.

The freeze and C qualification scripts are consumed. Native implementation is
still an ignored draft, not installed.
