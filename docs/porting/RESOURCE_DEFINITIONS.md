# Resource definitions — qualified C baseline

Qualified parent **45c36984**. Production remains original C: **18,662 physical
lines /61 files /zero reference C**. This scope has27 bodies /476 body lines:
24 live bodies /438 lines and three orphaned bodies /38 lines.

The live routines connect server object use/update/death/collision registrations,
client thing fields, sound-set loading, and reward/modifier catalog linking.
The two bitmask-name helpers and NPC voice setter have no external references
in the whole-repository literal audit. Remove these orphaned routines during
conversion. One export remains for the hurt-sound C caller: monster sound lookup.
The other23 live interfaces become private Go helpers. The sound-set list head
also moves to Go, together with its existing creature-xfer fixture owner.
Shared light scale/bias constants still have other users and remain unchanged.

## Qualification

- All149 affected roots pass on default/server/highres without skips. All63
  artifacts /30,040 records match across targets. All2,651 source fingerprints agree.
- Fifteen new frozen captures cover2,250 records. The original14 captures repeated
  unchanged in focused and affected runs before freezing. An additional numeric
  review added decimal overflow and direct float32 rounding contracts; these also
  pass against frozen expectations on all targets. No existing goldens changed.
- Static memory checks pass. Production source is identical to45c36984 apart from
  twelve new build-tag-only fixture files. Parent binary hashes were reverified,
  so its production/ABI/full-suite/gameplay/save-load evidence is reused for C.
  The native conversion must receive fresh production qualification.

See [selection](resource-definitions-selection.json),
[literal callers](resource-definitions-callers.json),
[interface plan](resource-definitions-interface-plan.json),
[captures](resource-definitions-captures.json),
[qualification](resource-definitions-c-qualification.json),
[production identity](resource-definitions-production-identity.json), and
[batch manifest](resource-definitions-c-batch.json).
Local final runs: `build/port-resource-definitions/c-qualified-{default,server,highres}`.
All sessions are joined. Installed fixture drafts, freeze.py and qualify-c.py are
consumed. Native drafts are prepared only under build; production is still C.

## Contracts and review decisions

Coverage includes full output buffers and adjacent-field preservation, signed and
byte boundaries, decimal overflow, defined incomplete scans, float32 midpoint and
range behavior, all21 server registrations, names and token delimiters, wand timing
at15/30/60 ticks per second, partial wand writes, all valid light angles, client
update lookup order, actual indexed/external image loading and MemFile positions,
catalog sentinels/strides/missing names, monster sound flags, encrypted sound-set
input, duplicate/list order, field offsets, comments, EOF and partial failures.
Owners use real catalogs, parser registrations, image handles, files and list data.
The affected selection includes gameplay/data consumers of these definitions.

Preserve two separately reviewable legacy behaviors:

- The MonsterArrow parser restarts tokenization on its original buffer and stores
  its first numeric value twice.
- Binfile.SkipLine stops on the first **non-newline** byte. Sound-set comments thus
  leave most comment text available as tokens. Do not reuse monsterDefinitionToken
  unchanged: its comment behavior differs. Two initial fixture assumptions failed
  on this distinction and were corrected before freezing; production was unchanged.

Malformed C inputs can access missing tokens or uninitialized stack bytes. Frozen
records cover defined partial writes, not random memory. Missing indexed images
are not covered because ImageByIndex has no bounds guard; current external named
image lookup returns nil and is covered. Future native fallbacks for undefined
inputs must be documented separately. On the current386 runtime, signed decimal
overflow clamps to int32 endpoints. Float parsing must round directly to float32;
parsing via float64 first loses a tested value just above a rounding midpoint.

Completed prior player-file/book-award scenario copies were deduplicated after
hash/inactivity checks, reclaiming2,225,495,402 bytes. Restore manifests remain
with the four runs; the helper's audit/apply modes are consumed. Original assets,
archive, saves, captures and production binaries remain intact.
