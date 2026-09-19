# Resource definitions — qualified Go conversion

C baseline **b078430c** is committed/pushed. All24 live bodies /438 body lines
are now Go; three orphaned bodies /38 lines and the private C sound-set list head
are removed. Production C is **18,044 physical lines /61 files /zero reference C**,
down **618**:476 body lines, one global, and141 obsolete comments/markers/blanks.

The routines connect server object use/update/death/collision registrations,
client thing fields, sound-set loading, and reward/modifier catalog linking.
One export remains for the monster hurt-sound C caller. The other23 live interfaces
are private Go helpers; Go callers/registrations invoke them directly. The existing
creature-xfer fixture now owns the native sound-set head. Shared light constants
still have other users and remain unchanged. No C test algorithm is retained.
The two bitmask-name helpers and NPC voice setter were orphaned by literal audit.

## Qualification

- All152 affected roots pass on default/server/highres without skips. All63
  artifacts /30,040 records match original C exactly. All2,659 source fingerprints
  agree between targets and production; no frozen expectations changed.
- Fifteen frozen captures cover2,250 new records. Native-only contracts cover the
  retained C ABI,32 repeated failed-reader loads without descriptor growth, and
  228 independently observed Linux/386 libc floating-point lexical cases.
- Static memory checks pass. Fresh production validates three ELF32/386/SSE2/CGO
  binaries and retained/retired interface inventories. Full-suite results match
  the known baseline exactly:1,553 failure entries;15 pass /3 fail /32 skip packages.
  Fresh headless gameplay and explicit save/load comparisons pass.
- The C baseline had149 affected roots and2,651 source fingerprints. Its production
  identity matched parent45c36984, allowing reuse of that parent's qualified
  production evidence. The native conversion received fresh production checks.

See [selection](resource-definitions-selection.json),
[literal callers](resource-definitions-callers.json),
[interface plan](resource-definitions-interface-plan.json),
[captures](resource-definitions-captures.json),
[C qualification](resource-definitions-c-qualification.json),
[native qualification](resource-definitions-native-qualification.json),
[float contracts](resource-definitions-float-lexical.json), and
[native manifest](resource-definitions-native-batch.json).
Final native runs: `build/port-resource-definitions/native-release-{default,server,highres,production}`.
All sessions are joined. All source installers, freeze/qualification generators
and copied drafts are consumed; do not replay them.

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
image lookup returns nil and is covered. Native fallbacks for undefined
inputs are documented below. On the current386 runtime, signed decimal
overflow clamps to int32 endpoints. Float parsing must round directly to float32;
parsing via float64 first loses a tested value just above a rounding midpoint.

Reversible cleanup for review: Go sound-set loading closes the reader on every
exit, preserving partial list/return behavior; C leaked it on unknown fields.
The independent descriptor contract verifies cleanup. Original open/key/line-skip
error diagnostics are preserved. Narrow deterministic fallbacks replace undefined
C cases: missing required tokens return failure, failed byte scans leave the byte
unchanged, and oversized token reads stop rather than overwrite the destination.
These are outside frozen defined-C equivalence.

The228-case lexer review corrected incomplete exponent/infinity matching, numeric
NaN payloads, and ASCII-only keyword handling (Unicode İNF must be rejected).
The standalone comparison reports zero differences; tests retain those independently
observed values without keeping C algorithms. This review happened too late and
caused repeated target sweeps. PORT.md now calls for libc/file-adapter review
before long milestone gates.

Disk cleanup deduplicated completed player-file/book-award asset copies
(2,225,495,402 bytes) and this batch's two completed scenario copies
(1,112,747,701 bytes). All removed copies matched the original asset hashes;
restoration manifests/helpers remain. Original assets/archive, changed files,
saves, captures and production binaries are intact. Both audit/apply helpers
are consumed; their --restore modes remain available.
