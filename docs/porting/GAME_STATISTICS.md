# Game statistics and report serialization

## Scope and current state

The caller audit selects **44 connected C bodies / 2,015 body lines**, following
qualified native server orchestration **bdf8cdcf**. The selection includes player
participation and event recording, session initialization/finalization, report
array ownership, field construction/serialization and the report's local encoding
helpers. Unrelated map loading routines in the same address range are excluded.
See [game-statistics-selection.json](game-statistics-selection.json).

The original-C baseline is in progress; no production conversion is installed. A three-line C prerequisite correction
is now installed and awaiting qualification.
C remains **30,819 lines / 68 files / zero reference C**.

## Caller and ownership audit

Most outside callers are existing Go wrappers. `sub_425CA0` also remains called
by C player-death processing in GAME5.c and will need a C entrypoint until those
callers move. It is also called by the Go objective-scoring owner. Participation
updates are called by the Go team/player owners. Private record functions have
no outside live callers in the source audit. The record cleanup forwarding body
`sub_42CC50` lives in GAME5_2.c and belongs to this batch.

Reporting endpoints use `sub_420360` to look up game-result service metadata.
The current send stubs abort if an old service is configured. Baseline fixtures
must exercise report construction locally, and the no-service path with owned
service state; they must not contact an external reporting service. This port
will not add an online service or silently change that existing branch.

## Contract plan

- Real C-owned field records: constructors, append wrappers, replacing/freeing
  payloads, linked order, repeated serialization, exact header/padding/endian
  bytes and unchanged host records after serialization.
- Scalar signed-byte argument narrowing into byte/word/dword fields, including
  negative values and larger original values. Strings, embedded zeroes, tag
  truncation/padding, empty payloads and length boundaries.
- Complete match/quest reports: distinct field values and real pointer arrays,
  player ordering, sequence counters and overflow, report modes, mutable tag data,
  and independent decoding of output records.
- Local encoding helpers: repeated/escaped runs, byte distributions, boundary
  lengths and deterministic seeded state. Own the nonzero numeric scale constant;
  do not let uninitialized mapped state produce a trivial oracle. Time-dependent
  output needs a justified clock seam or semantic decoding, not blanket omission.
- Player/session reporting: real player and network owners, flags/observer rules,
  counts/indexes, event ordering, report array rebuild/cleanup and elapsed time.
  Reconsider boundaries and scope as the C baseline exposes actual behavior.

Observed legacy details to preserve until separately reviewed: numeric record
constructors accept `char` even when storing two/four-byte values; serializer
prepends fields and restores host byte order after writing; report sequence
counters are wider than the serialized numeric argument. These need explicit
contracts rather than an assumed modern report format.

## Disk and recovery

Completed server-orchestration native scenario copies were hash-checked and
deduplicated, reclaiming **1,660,044,319 bytes**. Each run retains its restoration
manifest, changed maps, saves, screenshots and logs. Original assets/archive and
production binaries are unchanged. The audit/apply modes of
`build/port-game-statistics/deduplicate-server-orchestration-native-assets.py`
are consumed; only restore mode is reusable.

## Baseline development notes

Record contracts include independent exact bytes, length truncation through
16-bit fields, total-length header wrapping, host-state restoration and repeated
serialization. Complete report fixtures now initialize the shipped report-name
tables; a data audit caught zero-filled tags before freezing. Both intermediate
and final match reports increment the sequence counter.

The run encoder compares an unsigned input byte against signed `char`; high-bit
bytes do not form runs like low bytes. Low-byte runs also wrap their byte counter
at 256. Capture and preserve this behavior; changing the report format is outside
this compatibility conversion. The seeded random generator can produce values
outside [0,1) for some large seeds; do not assume a stronger range than its C
implementation. Direct arithmetic/state and byte-consumption contracts supplement
frozen sequences.

The statistics fixture reuses the existing map-theme time/allocation observer
for deterministic C wall-clock values, with allocation callbacks disabled and
restored afterward. A first attempt to add a separate wrapper hit a duplicate
link symbol; it was replaced by this existing owner. No shared wrapper or build
flag changes are needed, and production is unchanged. Clock-controlled session
and encoding qualification is still in progress.

## Event registration prerequisite (review item)

Original C at **bdf8cdcf** selected the first player's address using the second
player's host-slot test. A newly recorded host with a remote target reached the
existing address bridge with index32 and panicked (`eighth.jsonl`). Independent
name/address assertions then showed a remote actor paired with an existing host
receiving the host address instead (`event-original-contract.jsonl`, mask1,
actor0, target2). Inspection also found the target's new name and host address
being written to the actor's row index rather than the target's newly allocated
index.

Correct three references in `sub_425CA0`: test the actor's own player index for its
address, and use the target's new index for its name and host address. This also
makes the explicitly supported absent-target branch valid for a fresh actor.
The expanded contracts cover all actor/target identities, absent participants,
new/existing rows and flags, with nonzero real connection addresses. This is a
confident reversible correction under the standing authorization; no golden was
frozen or changed to hide the bug. C LOC is unchanged. Fresh production/headless
qualification is required; the preceding production evidence cannot be reused.

## Repeated C capture checkpoint

**21 focused roots / 1,519 test entries** pass in separate processes. All **21
captures / 3968 records** match exactly and their hashes are frozen in the
fixtures. See [game-statistics-captures.json](game-statistics-captures.json).
Static mapped-memory checks pass. The automatic event-flush boundary is covered
at 253/254/255 existing rows, with actual report arrays and roster rebuilding.
Embedded-zero string versus binary payloads have separate contracts.

At capture checkpoint05033129 this was **not yet a fully qualified C baseline**.
The three broader targets and fresh production/headless qualification remain.
Use [game-statistics-c-batch.json](game-statistics-c-batch.json). Its selection adds
statistics, objective scoring, damage, object death and AI combat to the preceding
qualified target pattern. Do not begin conversion before those gates pass.

Five old inactive failed-run asset copies were also audited and deduplicated,
reclaiming **2,755,182,757 bytes**. Their recorded exit statuses were preserved;
only files matching the original assets byte-for-byte were removed. Restoration
manifests, failure logs, changed maps/saves and screenshots remain. The audit/apply
passes of `build/port-game-statistics/deduplicate-old-failed-assets.py` are consumed.

## Conversion reachability refinement

The full Go callsite review found literal-false guards around statistics session
initialization, finalization, periodic reporting, join registration and departure
completion. A textual reference is not sufficient evidence of a live root.
Follow the two remaining live entrypoints: event recording (`sub_425CA0`, C player
death and Go objectives) and participation (`sub_425ED0`, Go team membership).
Their private-helper closure contains **35 C bodies / 1,472 body lines**.

**Nine C bodies / 543 body lines are orphaned** and will retire with their disabled
callers, not be translated solely to preserve test coverage. This also removes
the final body from server__system__server.c. The original C checkpoint preserves
the complete source and contracts. See
[game-statistics-reachability.json](game-statistics-reachability.json).

Before conversion, retire four orphan-only test roots and project the two mixed
captures from the committed C data: keep participation cases and non-quest
reports. Verify those projected hashes on C. Seventeen live-contract roots /
3,366 captured records remain; event-triggered flushing already covers the live
reporting lifecycle. Generic modes within the still-live match serializer retain
their frozen contracts. Do not regenerate expectations from Go results.

The first production attempt built successfully but its manifest incorrectly
classified original-C routines as Go-backed exports. They are now listed under
`retained_c`; all44 symbols were verified present in the first binary. The fresh
production rerun is under `c-production-final`. This is a qualification-metadata
correction, not an application failure or production-source change.

## Qualified C baseline and live projection

All three targets pass **495 roots /43,062 test entries** without skips. The
267 captures /99,629 records (including combat) agree across targets; all21 new
statistics captures match the frozen hashes. All four gates have the same2,450
source-file hashes. Three fresh production builds/ABI, exact known full-suite
failures, gameplay, save/load and flat-map regeneration pass. See
[game-statistics-c-qualification.json](game-statistics-c-qualification.json).

The live-only test projection is now independently verified on that same C
implementation: **17 roots /1307 entries**,17 captures /3,366 records. Two hashes
were computed by filtering qualified C capture records; all remaining hashes
are unchanged. No production source changed during projection. The full21-root
three-target and production evidence above remains applicable; no duplicate
production run is needed for this test-only change. See
[game-statistics-live-captures.json](game-statistics-live-captures.json).
