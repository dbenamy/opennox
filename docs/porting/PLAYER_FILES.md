# Player-file sections — qualified Go conversion

C baseline **583caebb**; qualified production parent **d7f52707**. All18 selected
bodies are now Go. Production C is **19,482 physical lines /65 files /zero reference
C**, down **1,601 lines**. Six C interfaces remain for live callers and mapped
callbacks; twelve private interfaces and the server section table retire.

All146 affected roots pass on default/server/highres. All66 captures /6,567 records
match the frozen C baseline exactly; all2,630 source fingerprints agree. Fresh
production binaries pass ABI checks, the full suite matches its known failure and
package-result sets exactly, and headless gameplay plus explicit save/load pass.
Static memory-map checks pass. See [native qualification](player-files-native-qualification.json).
The36 new player-file captures contain1,951 records; no frozen expectations changed.

Earlier progress notes below are historical. The production-identity report applies
to the test-only C baseline; the conversion has fresh production evidence.

The batch selects 18 functions /1,540 body lines in GAME1_1.c: player attributes,
status, inventory, field guide/spellbook/enchantments, journal, game and GUI data,
file metadata, music state, section framing and client player-save requests.
The adjacent briefing gate helpers are unrelated and excluded. Whole-repository
review finds 16 external roots and all 18 bodies reachable. The local C section
table is also owned by this batch; distinguish its internal callback references
from interfaces needed by outside C callers during translation.

Reuse real player/object/inventory, journal, quickbar, music, timer and cryptfile
owners. Verify exact section bytes, historical version gates, read/write mutations,
flags, optional fields, Unicode lengths, file positions, callback order and cleanup.
File metadata uses the wall clock: identify and validate that field explicitly
before normalizing any writer capture. Use deterministic input timestamps on reads.

Production source is currently unchanged. Reuse the qualified parent production
baseline only after checking identity; corrections would require fresh qualification.
No expectations are frozen and no conversion has started. The actual-C dispatcher
and focused manifest are installed. Remaining work is listed below.

The first music C sweep passes both roots: 24 write and 112 read records cover
exact state-word order, current and historical signed version gates, optional
section flags, queue levels, counts through the adjacent-slot boundary, untouched
backing bytes, file position and real music module changes. `music-initial` is
joined. Player attribute contracts pass (25 write /54 read records). They cover
field order, UTF-16 name lengths, supported/historical versions, name rejection,
mode flags, legacy skipped words and guards with an actual PlayerInfo allocation.
The attributes draft is consumed; use tracked source. No expectations are frozen.

Book contracts add exact spell/ability write bytes for all three classes, sparse
and full entry sets, flag gates, and read version/count/presence boundaries for
both books. Two discovery attempts found incorrect fixture imports; corrected
before books-third. Production remains unchanged.

books-third (six roots), fieldbook-initial (seven) and status-initial (eight)
passed and are joined. Field-guide names use explicitly owned lookup strings,
including empty, UTF-8 and 255-byte names. Status writers cover unsigned word
boundaries, current/max health and mana, saved temporary words, poison fields
and object guards. Their drafts are consumed. Spellbook restoration now tests
real award calls, bookkeeping records and reliable report queues; all 90 cases
pass. No frozen baseline or conversion yet.

Journal/inventory contracts pass (journal-inventory-second, 13 roots), including
existing journal save/load coverage. The first attempt used a multiline pattern;
the driver requires one regex line, corrected before running tests.
Guide-family restoration initially failed because the fixture passed an initial
value to alloc.New, whose argument supplies only the type. Explicitly initializing
the allocated array fixed it; no production change. Metadata reads and enclosing
GUI section contracts pass together in gui-guide-third (17 roots).

Metadata writes and recovery-repeat both pass all 18 roots. Seventeen new captures
/1,006 records repeat byte-for-byte; the eighteenth root reuses the existing
journal save/load golden. Timestamp normalization follows actual clock bounds and
weekday checks. The static scanner required literal 10980 instead of a local
constant in two metadata accessors; static-recovery now passes. All runs are joined.
This is a partial-baseline recovery checkpoint, not a frozen or qualified conversion.
The tracked recovery JSON records capture hashes and counts; expectations in new
tests remain unfrozen. Production source is identical to d7f52707: only porttest
files and documentation have been added/changed. Full baseline review remains.

Remaining before freezing: status reads / real-unit attributes, complete inventory
read/write, enchantment state and historical ability tails, game/script sections,
outer file framing/extraction and save request, remaining quest-mode book gates.
Review coverage and real constant/table owners, repeat captures in fresh processes,
qualify affected targets, verify production identity, then commit frozen baseline.
All installed drafts (attributes, fieldbook, status, journal gates, inventory helpers,
guide restore, metadata read/write, GUI) are consumed; use source, never replay.

Recovery checkpoint 33a1248a is pushed. Status reads subsequently pass all 216
cases, including current-versus-saved health/mana, poison state/timestamp,
experience report, signed versions, optional section gates and full owner guards.
The draft is consumed. Enchantment writers are installed/running, with explicit
iteration order/count ownership, active ability records, cooldowns and shield
duration records. These are actual production owners, not algorithm substitutes.

Enchantment write (96 records), read gates/empty ability tails (128), game-state
write/read (9/90), save requests (10) and extraction (13) pass; extraction-initial
finishes 25 roots. Actual guide enumeration/count, ability map/server clock, shield
record and quest-variable owners are explicit. Version 4's warrior tail exists
outside the optional enchantment body; version 5's tail is inside it. A stored
ability-list entry counts as active even when its Active field is zero.

The game section preserves the existing doubled byte-length map field, including
following record bytes, and ignores the stored object ID on reads. Tests cover
map lengths through 31 within the real Player record and verify the untouched
surrounding record. This format convention is not corrected in this port.

Extraction uses real key-27 files, mixed known/unknown sections and sizes around
8-byte boundaries. Destination padding and guards remain unchanged. Missing files
leave output untouched. Client load framing with real callbacks is now running.
Remaining: populated enchantment restoration/active ability replay, complete
inventory read/write, real-unit attributes and quest-mode book gates, client writer
and server loader framing, repeated baseline/affected target qualification.

Client-load-initial passes (26 roots), exercising actual callbacks, unknown sections,
GUI/music composition, section rejection and file closure/path-update rules.
Server loader framing is under test with real player stats, XP values, ability
owners and health/mana state. Its first fixture expected previous mana because of
the legacy name unitGetOldMana; the actual dependency returns current mana. The
fixture now expects 11, preserving observable resource restoration. The second
sweep also tests real infravision application rather than a callback substitute.

Server-load-enchant-second, ability-restore-initial and attributes-unit-initial
pass and are joined: 30 roots. Infravision restoration covers historical default
power, clamping and zero-duration behavior (140 records); warrior replay covers
strict active flags, cooldown omissions, remaining duration and actual reports
(36). Real-unit attributes cover names/colors, old and rejected life counts,
quest reset/stage fields and report bytes (56). The name/colors mutations happen
before excessive lives reject the section. Registered setter bookkeeping remains
outside this fixture; its dependency contracts are separate.

Client writer framing is installed and being checked. Main inventory read/write
and quest-mode book validity/level gates remain before baseline review, repeat
captures, affected targets and freezing. No production source has changed.

Client-write-second passes 31 roots; all 29 earlier captures repeat unchanged.
Eight writer cases cover metadata-only versus all client sections, nonzero mode
values, real key-27 framing, GUI/music payloads, terminator, file closure, copied
input/guards and failed opens. The first fixture assumed blank names for empty
quickbar slots; the existing contract requires SPELL_INVALID, now corrected.
Only validated temporary paths and the independently checked timestamp field are
normalized in the writer capture. Static-writer and static-inventory pass.

Inventory-write-second passes 32 roots. Fifteen cases cover optional-section
flags, empty/list/grid order, actual Gold/Ammo xfers, missing grid codes, and
secondary/quiver ScriptID conversion. The fixture's initial compile used uint32
for the Go int ScriptID field; corrected before the passing run. Read gates and
independent loaded-report checks are now under test. Populated reads and quest
filter/equipment integration remain; the inventory contract is not yet complete.

Inventory-read-initial and recovery-second-repeat pass all 33 roots; all 32
captures /1842 records match exactly. Read gates add 112 cases spanning
signed versions, optional/mode gates, real gold changes, unknown types, quest
count rejection, historical tail fields and loaded-report emission. The limit
table points to the actual controlled Gold type and the staff limit is explicitly
3. Existing independent item-transfer contracts supply the writer byte oracle.
No populated successful inventory read is claimed yet. All sessions are joined.
Production source identity against d7f52707 is unchanged; this is a recovery
checkpoint with unfinished C baseline work, not a completed port.

Recovery af59461e is pushed. Quest-books-initial passes 35 roots, adding 80 spell
and nine guide cases. Spell level comparison is signed: 0xffffffff is -1 and
passes the <=3 quest gate. Actual enabled reward rows, disabled rows and guide
family eligibility are covered with real award/checksum/report owners.

Review found that the new inventory reader owner had omitted warrior strength,
which is the carrying-capacity denominator even for a wizard. It now supplies
40 explicitly; production is unchanged. Populated inventory restoration exposed
a missing fixture spatial index in the actual pickup path. The existing empty-map
owner is now installed, with map cleanup and item unlinking before release.
The second sweep includes quest writer filtering, equipped records for excluded
items, selected quivers and actual nontrivial inventory-limit acceptance/rejection.
Static-quest-inventory passes; the runtime sweep is in progress.

Inventory-restore-third passes 37 roots /36 new captures /1,951 records. The
second attempt found two fixture expectations to correct from source: the
constant 1161297920 encodes float2944 (not2000), and the filter bit0x40 is distinct
from the named ClassImmobile enum. Populated read ownership, values, order,
positions, saved IDs/quest renumbering and completion reports now pass. Twelve
quest writes check class/Glyph exclusion, unfiltered equipped/quiver tails and
limits accepting2 items at3 versus rejecting at1. No production correction.

The affected selection (146 root declarations) covers the new sections, journal
and quickbar save owners, object/item transfers, server inventory/equipment,
quest eligibility, resources, player stats/name, inventory notifications, quest
variables and spell lifecycle/duration. UI grid/selection queries and secondary
notifications cover the writer's client inventory dependencies. It deliberately
omits unrelated drawing, tooltip and input suites. c-affected-default passes;
all 36 player-file captures match inventory-restore-third and c-focused-highres
exactly, with identical source fingerprints. Server focused qualification is
still running before expectations can be frozen.

Final frozen C qualification passes all 146 roots on each target with no skips:
default 136.643s, server 132.037s, highres 132.064s. All66 artifacts /6,567 records
match exactly across targets; all 2,623 source fingerprints match current source.
The production-identity report verifies every added/changed source file is behind
the porttest constraint. The qualified d7f52707 production binaries, known full
suite result, fresh gameplay and explicit save/load are reused as planned.
Static-quest-inventory passes. All sessions are joined. No C algorithms changed.

Review limits: full equipped-weapon restoration, registered stat/name bookkeeping,
nonempty magic-wall sections and gameplay save/load also rely on their existing
dependency contracts and production integration. The batch does not promise safe
handling of arbitrary corrupt files where the original uses uninitialized reads
or unchecked indices. Native code must preserve supported/historical behavior;
any intentional format or error-handling change needs separate recorded evidence.
The capture index freezes only validated timestamp and temporary-path
normalization. All other section bytes, effects and reports remain exact.

C baseline583caebb is committed/pushed. The18-body Go conversion is installed;
six outside-C/callback exports remain and twelve private interfaces retire.
The shared object-transfer stream preserves serialized field call boundaries;
server section dispatch is native, while the client mapped callback table retains
its actual exported functions. Go callers and the journal fixture call Go directly.
Working C is19,482 physical lines /65 files /zero reference, down1,601 from C.
Static-native-initial passes; initial native qualification began against locked expectations.
No native qualification is claimed yet. All native drafts/install script are consumed.

Initial native discovery failed on a Cgo extern spelling inconsistent with the
existing unsigned-int local-player-code declaration; corrected. The second build
found a real C caller of sub_41CEE0 in character creation (selcolor.c). The original
caller audit's broad prototype regex had treated a return-call as a declaration.
The tracked audit is corrected, and literal scanning across all remaining C files
confirms three C entry points plus three mapped callbacks. Six exports remain;
twelve private interfaces and two server-table symbols retire. This changes only
the interface inventory, not the batch algorithms or frozen expected behavior.

Native-fourth compiles and passes 36/37 focused roots. The independent fieldbook
write contract caught a native width error: learned-guide flags are uint32 rows
starting at player offset4248 (index1), not bytes. Corrected the native view to
uint32 at4244; frozen captures and contracts remain unchanged. The preceding
third attempt only required explicit Go player-index casts.

Final native qualification is complete. Default/server/highres affected checks
passed in74.897s/234.752s/138.128s; all sessions are joined. Fresh production took
303.975s and verifies three Go-backed ABI inventories, absence of retired symbols
and porttest helpers, and the exact known full-suite result (1,553 failure entries;
15 pass /3 fail /32 skipped packages). Gameplay and explicit save/load scenarios
both pass. No intentional gameplay or file-format change was introduced.
