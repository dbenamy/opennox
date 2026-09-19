# Porting checkpoint

Read [PORT.md](PORT.md) for the working plan. This is the resume checkpoint.

**Qualified C remaining: about 14.5k lines** — **14,457 physical lines in 58
production `.c` files**, zero reference C; +6 prerequisite correction lines.
See [C_LOC.md](docs/porting/C_LOC.md).

<!-- current-checkpoint -->

## Current — session dialogs corrected C baseline qualified

Qualified parent `e6245c82`; the connected32-body /907-line C selection covers
server filters, MOTD, disconnect and quit menus. The original-C baseline and
independent contracts are complete. See [SESSION_DIALOGS.md](docs/porting/SESSION_DIALOGS.md).

Corrected two fixed MOTD buffers before freezing: input-sized temporary storage
preserves line handling and byte widening; the listbox retains its255-unit row
limit. Working/qualified C is14,457 /58 /zero reference (+6 prerequisite lines).
Default/highres:307 roots /227 captures; server:306 /226, because the existing
world-selection test is !server. All shared captures match, all142 parent captures
are unchanged, and nine new dialog captures repeat independently. All2,715 source
fingerprints agree. Static, fresh three-target builds/ABI, exact known suite
(1,553 failure entries;15 pass /3 fail /32 skip packages),13-screen filter,
character-creation/gameplay and save/load all pass.

All sessions joined. Evidence:build/port-session-dialogs/c-{default,server,highres,production}
and the tracked C qualification/capture manifests. Commit/push this baseline,
then convert31 live bodies and remove the proven no-op4896E0 cleanup. Preserve
its Go wrapper's independent context reset. Four selected interfaces still have
unselected C callers:446360,446780,446950,445C40. Audit callback slots before
retiring the others. Move actual briefing/inventory/options owner views with them.

An initial native filter draft is under build/port-session-dialogs/native-filter-draft.go;
NOT installed, formatted, compiled or accepted. All installed fixture drafts and
qualify-c.py are CONSUMED. Do not replay stale source drafts or selection offsets.
About2.7GiB is free; hash-audit completed C scenario copies before the native builds
need more room. Preserve original assets/archive, captures, changed saves and logs.

Fixture review notes: the server owner seeds two teams; quit contracts now control
counts0/2/256. MOTD's quest gate reads runtime state1556160 independently of game
flags. Empty/nil row slices both mean no rows. These were corrected test assumptions,
not production behavior changes. No unresolved user question.

## Qualified parent — orphaned configuration callbacks

Browser parent **c857b71f** is committed/pushed. Twelve unreachable callbacks,
216 function-body lines plus 24 separator lines, and their declaration/registration
plumbing are removed. Commits `391c59fd` and `d15e0aec` removed their caller/parser
in 2022; current literal/dynamic/raw-blob audit finds no table reader. Live browser
coordinate storage remains. See
[CONFIG_CALLBACK_RETIREMENT.md](docs/porting/CONFIG_CALLBACK_RETIREMENT.md).

All 197 affected roots pass on default/server/highres, no skips, with 142 exact
original-C artifacts. All 2,703 source fingerprints agree. Static checks, three
fresh binaries/ABI inventories (all twelve retired symbols absent), exact known
suite (1,553 failure entries; 15 pass /3 fail /32 skip packages), browser,
character-creation/gameplay and save/load pass. No goldens changed and no C
algorithms kept for tests. All sessions joined; final qualification evidence is
under build/port-config-callbacks/final-* and the tracked qualification report.
remove.py and qualification/finalization scripts are CONSUMED.

Qualified browser asset audit18221/apply39728 reclaimed 1,669,136,451 bytes;
deduplicate-final-browser-assets.py deletion is CONSUMED. Preserve restoration
manifests, changed saves and all screenshots/logs plus original assets/archive.

The cleanup is pushed; the active session-dialog baseline is described above.

## Qualified parent — server-browser Go conversion

Original-C baseline **04061130** and qualified production parent **d8133587** are
committed/pushed. The connected browser batch is native: 64 live algorithms,
one nine-line orphan removed, 58 private selected interfaces and 34 C global owners
retired. Seven selected exports remain. The empty noxworld C file and previously
deferred character initializer adapter/header are removed. Two shared C connection
flags remain. See [SERVER_BROWSER.md](docs/porting/SERVER_BROWSER.md).

All 197 affected roots pass on default/server/highres without skips; all142 artifacts
/100,275 records match original C, including18 browser captures /72,402 records.
All 2,703 source fingerprints agree across tests and production. Static checks,
three fresh binaries/ABI inventories, exact known suite (1,553 failure entries;
15 pass /3 fail /32 skip packages), six-screen browser, character-creation/gameplay
and explicit save/load scenarios pass. No goldens changed.

Evidence: build/port-server-browser/native-final-{default,server,highres,production},
static-native-final.log and tracked native-qualification.json. All sessions joined.
All source installers, baseline freezers, boundary-retirement and qualification
scripts are CONSUMED. Original selection byte offsets are stale; never replay them.

Review notes: preserve byte-widened UI strings, signed-short high-port lookup,
unsigned geometry behavior, unconditional C highres description and distinct
32/64-bit deadline additions. Keep one 169-byte selected-record snapshot and one
12-byte sentinel across browser closure; gameplay still reads the endpoint.
Zero formerly undefined map-polygon metadata and password-header padding.
Original C leaked abandoned lists; native sorting reuses nodes.

Initial headless runs exposed freed static-label text and an offset applied to
the scrollbar parent rather than its thumb. Persistent text and correct child
target fix both. The label-lifetime regression demonstrated the failure before
fixing; all final gates use corrected source and unchanged expectations. Earlier
failed native-release/native-reviewed runs remain diagnostic evidence only.

Disk cleanup: completed C browser copies reclaimed 1,112,777,500 bytes; old inactive
compiler-cache entries reclaimed 7,527,475,663 bytes; three completed native browser
copies reclaimed 1,669,166,250 bytes. Deletion scripts are consumed; keep restoration
manifests, captures, screenshots/logs, originals/archive and qualified binaries.

Next: commit/push this native batch, then audit the 12 old config callbacks. Read-only
preview/table-reference/dynamic-access inventories are under build/port-config-callbacks.
Literal registrations exist but no table reader or incoming raw pointer has been
found. Finish tracing dynamic accesses before removing callbacks. No next-batch
source changes are installed; no user decision pending.

## Qualified parent — character-creation Go conversion qualified

C baseline **dd4debb6** is committed/pushed. All 23 selected bodies are Go;
16 private interfaces and 27 named C owners retire, and two C files are removed.
Seven thin selected C exports remain for live entries and animation slots.
See [CHARACTER_CREATION.md](docs/porting/CHARACTER_CREATION.md).

All 174 affected roots pass on default/server/highres without skips. All 142
captures /13,807 records match C, including 16 new captures /889 records. All
2,672 source fingerprints agree across tests and production. Static memory checks,
three fresh ELF32/386/SSE2/CGO binaries/ABI inventories, exact known full suite
(1,553 failure entries;15 pass /3 fail /32 skip packages), extended character
creation/gameplay and explicit save/load qualify. No goldens changed.

Final evidence: build/port-character-creation/native-reviewed-{default,server,
highres,production}, static-native-reviewed.log and native-qualification.json.
All sessions joined. All installers, freezers and qualification generators are
CONSUMED. Source and frozen expectations govern; never replay old installers.

Review notes: the first headless native run found event22's numeric child ID
interpreted as a pointer. Fixed by checking event type before decoding arguments;
an eight-ID independent regression now passes, and all gates were rerun. Preserve
legacy first-character name cleanup and UTF-16-as-byte-string filenames. Native
save records zero two unused, nonserialized padding bytes. Shared animation
callback storage stays unchanged; callback logic itself is Go.

Disk: old cache cleanup reclaimed 4,924,575,701 bytes from 18,704 checked old
entries. Three completed native scenario copies reclaimed another 1,669,136,451 bytes;
audit11424/apply65743 joined and are CONSUMED. About 5.7 GiB free;
keep restoration manifests/helper, captures, saves, assets/archive and binaries.
Conversion **d8133587** is committed/pushed; the next baseline has started.
The initial21-body browser candidate has since expanded to65 bodies /2,067 lines;
see the current checkpoint and tracked server-browser plan above.
The now-unused sub_4A5E90_A C adapter (implementation already Go) and empty
selclass header can retire with that next source batch; retain the active Go hook.

## Qualified parent — resource-definition Go conversion qualified

C baseline **b078430c** is pushed. All24 live bodies are Go; three orphaned bodies
and a private C list head are removed. One C export remains for monster sounds.
Go registrations/callers and the creature-xfer fixture use the native owners.

All152 affected roots pass on default/server/highres without skips. All63 artifacts
/30,040 records match original C;15 new frozen captures contain2,250 records.
All2,659 source fingerprints agree. Static checks and fresh production qualify:
three ELF32/386/SSE2/CGO binaries/ABI inventories, exact known full-suite results
(1,553 failure entries;15 pass /3 fail /32 skip packages), gameplay and save/load.
See [RESOURCE_DEFINITIONS.md](docs/porting/RESOURCE_DEFINITIONS.md) and
resource-definitions-native-qualification.json. Final evidence:
`build/port-resource-definitions/native-release-{default,server,highres,production}`.
All sessions joined; no builds/tests running. All installers, copied drafts,
freeze/qualification generators are CONSUMED. No frozen expectations changed.

Review decisions: preserve legacy comment skipping and MonsterArrow first-value
reuse. Close failed sound-set readers on every exit, retaining partial results and
original diagnostics. Use narrow deterministic fallbacks for undefined malformed
inputs. Retained-C ABI, repeated reader cleanup and228 libc lexical contracts pass.

Next: character creation,23 connected UI bodies /1,042 body lines. Read-only
selection/global audit and name fixture drafts exist under build/port-character-creation;
none installed yet. Adjacent13 configuration routines /217 body lines are deferred.
Review actual UI resources, animation, palettes, names, preview pixels and saved
player files. Literal global audit must also follow the three vardefs.go modifier
setters used by root modifiers.go. No next-batch source edits before this commit.

Disk: prior four completed scenario copies reclaimed2,225,495,402 bytes; this
batch's two completed copies reclaimed1,112,747,701 bytes. Audit/apply helpers
are CONSUMED; keep manifests and --restore helpers. Original assets/archive,
saves, captures and production binaries remain. About3.9 GiB free after cleanup.

## Qualified parent — book-award Go conversion

C baseline **27554f1f** is committed/pushed. All26 selected bodies and the private
counter are Go; seven C interfaces remain, nineteen private interfaces and one C
counter retire. Four C translation units are removed. Go callers use native helpers.

All108 affected roots pass on default/server/highres without skips. All81 captures
/12,248 records match original C exactly; the15 new captures cover6,763 records.
All2,639 source fingerprints agree, static-native-final passes, and no frozen
expectations changed. Fresh production validates three ELF32/386/SSE2/CGO binaries,
interface inventories, exact known full-suite results (1,553 failure entries;
15 pass /3 fail /32 skip packages), headless gameplay and explicit save/load.
All book-award sessions are joined.

See [BOOK_AWARDS.md](docs/porting/BOOK_AWARDS.md) and
book-awards-native-qualification.json. Evidence:
`build/port-book-awards/native-final-{default,server,highres,production}`.
All drafts, move-callers.py and finalize-native.py are CONSUMED; never replay.

Review decisions: preserve spell-family cap/bookkeeping on the original ID, and
unknown field-guide item consumption when admission passes though award rejects
ID0. No gameplay correction was made. Existing stats and book-tooltip contracts
were added to the native affected selection for callers moved to Go.

Next candidate: `build/port-resource-definitions/{selection-draft,callers-draft}.json`
and plan-draft.md.27 bodies /476 body lines; three appear orphaned. Verify roots,
then original-C contracts for the24 live functions. Next-batch fixture sources are installed; see active checkpoint above.

## Qualified parent — player-file conversion

C baseline **583caebb** is pushed. All18 selected bodies are Go, with six retained
C interfaces and twelve private interfaces retired. The owned server section table
is native dispatch; the client mapped table keeps its exported callbacks.

All146 affected roots pass on each target with no skips; all66 artifacts /6,567
records match C exactly and all2,630 source fingerprints match. The36 new captures
contain1,951 records. Static-native-qualified passes. Fresh production validates
three ELF32/386/SSE2/CGO binaries and their interface inventories, exact known full
suite (1,553 failure entries;15 pass /3 fail /32 skip packages), headless gameplay
and explicit save/load. **All sessions are joined; no builds/tests running.**

The independent fieldbook writer contract caught a native byte/uint32 layout error,
now fixed. The caller audit had misclassified a return-call as a prototype; the
client writer C export is retained for character creation. No frozen expectations
or gameplay/file-format behavior changed. See [PLAYER_FILES.md](docs/porting/PLAYER_FILES.md)
and player-files-native-qualification.json for evidence and review limits.

The qualified conversion is pushed as a8015048; continue the book-award baseline.
Local evidence: build/port-player-files/native-final-{default,server,highres,production},
native-fifth and static-native-qualified.log. All drafts, installers and freeze.py
are **CONSUMED**; never replay them. Actual source and frozen expectations govern.

Disk cleanup reclaimed **8,604,874,683 bytes (8.014 GiB)** from98,529 audited regular
Go cache files older than12 hours; four empty cache directories were retained.
Audit/partial/applied manifests under build/port-player-files are consumed. Assets,
module cache, captures, saves and production binaries remain intact; about10 GiB
free before the final scenarios. Preserve assets and the local original archive.

## Qualified parent — native client audio events

C baseline **e30b952e**, qualification **994d7b3c**, conversion **d7f52707** are
pushed. All62 bodies are Go;21 C interfaces remain and 41 private interfaces retire.
All three targets, ABI/builds, exact known full-suite results and fresh gameplay/
save-load pass. C is 21,083 lines /65 files /zero reference, down 1,327 lines.
Details: [CLIENT_AUDIO_EVENTS.md](docs/porting/CLIENT_AUDIO_EVENTS.md). Completed
C/native scenario asset copies were verified/deduplicated, reclaiming 2,225,495,402
bytes. Their cleanup scripts are consumed; preserve restoration manifests.

## Qualified parent — native client audio streams

**50f0711d**, C baseline **2f4dbfea**, qualification **419a85dc**. Seventy bodies,
22 retained interfaces, 48 private interfaces and one C global retired. Sixteen
frozen captures /927 records passed the first native build; affected targets,
fresh production/gameplay/save-load qualified. C was 22,400 /66 files /zero
reference (−1,261 from corrected C). See CLIENT_AUDIO_STREAMS.md.
All installers/freezers are consumed. Four completed stream C/native scenario
copies were verified/deduplicated, reclaiming 2,225,495,402 bytes; preserve their
restoration manifests. Earlier cache cleanup reclaimed 8,618,371,198 bytes and is
consumed. Original assets, archive, captures, saves and binaries remain.

## Qualified parent — native combat overlays

**15836cdb**, C baseline **6998b8ba**: 32 bodies, 21 private interfaces and eight
private globals retired. Eleven live C exports remain. Nineteen focused captures
/707 records and 84 affected captures /18,691 records match; three fresh production
binaries, exact known full-suite results, gameplay and save/load qualified.
C was 24,411 /66 files /zero reference (−957). Pre-baseline fixes preserved literal
feed names, cleared missing-victim text and corrected the Go DrawableFX link offset.
See COMBAT_OVERLAYS.md. All installers/freezers are consumed. Older compiler-cache
and scenario deduplication manifests remain; superseded binaries can be rebuilt
from their recorded revisions. Original assets/archive remain untouched.

## Qualified parent — native speech bubbles

**f588409c**, C baseline **bf03e4d6**: 13 bodies and two private globals replaced;
12 captures /552 records, all target comparisons and fresh production qualified.
C 25,368 /67 files /zero reference (−779). See CHAT_BUBBLES.md. All its mutation
scripts are consumed. Its C/native asset restoration helpers and compressed-log
manifests remain under build/port-chat-bubbles; restore individual gzip logs with
gzip -dk. Preserve original assets and the untracked archive.

## Qualified parent — native player death

**b16f1a90**, C baseline 4c3191c1: seven C bodies replaced, six private interfaces
removed. Fourteen focused captures /1,101 records; each target 318 roots /24,386
entries; 149 captures /33,295 records match. Fresh production/gameplay/save-load
qualified. C 26,145 /67 files (−693). See [PLAYER_DEATH.md](docs/porting/PLAYER_DEATH.md).
All its copied drafts/installers and cleanup apply modes are consumed; restoration
manifests remain under build/port-player-death.

## Qualified parent — native map metadata

**00131036** is committed and pushed; C baseline f3181558. Five C bodies and six
adapters retired. Five focused roots /478 entries; each target 515 roots /47,472
entries; 287 captures /103,626 records match. Fresh production, gameplay, save/load
and compressed flat-map regeneration qualify. C 26,836 /67 files (−165).
See [MAP_METADATA.md](docs/porting/MAP_METADATA.md). All its scripts are consumed.

## Qualified parent — native floor/wall map sections

**90fbd480** is committed and pushed; corrected C baseline71633546. Eighteen C
bodies and two globals retired.19 focused roots /4,144 entries; each target510
roots /46,994 entries;282 captures /103,153 records match. Fresh production and
headless integration qualify. C count27,001 /67 files (−1,793 including174 obsolete
headings/blanks). See [MAP_SECTIONS.md](docs/porting/MAP_SECTIONS.md). All scripts
are consumed. Preserve original IO grouping: split coordinate IO changed only
checksums on the initial native run; grouping was restored without golden changes.

## Qualified parent — native game statistics

**6035e191** is committed and pushed; C baseline858bc315. Converted35 live routines,
removed nine orphan routines,43 C interfaces, seven globals and the final
server__system__server.c file.17 focused roots /1,307 entries; each target491 roots
/42,850 entries;263 captures /99,027 records match. Fresh production/integration
qualified. C count28,790 /67 files (−2,029); current C prerequisites add4 lines.
See [GAME_STATISTICS.md](docs/porting/GAME_STATISTICS.md). All statistics installers,
qualification scripts and C/native asset audit/apply modes are consumed; restore
modes remain available. Never replay older mutation scripts.

## Qualified parent — server map and round orchestration

Native checkpoint **bdf8cdcf** is committed and pushed. Ten live routines converted;
thirteen C function symbols and two globals retired. Its C baseline is ec5ffeb4.
Qualification: 34 focused roots /1,837 entries; each target 427 roots /41,317 entries;
245 captures /78,253 records match C. Fresh builds/ABI, exact known suite and three
headless scenarios pass. C count after this parent: 30,819 lines /68 files.
See [SERVER_ORCHESTRATION.md](docs/porting/SERVER_ORCHESTRATION.md).

## Qualified parent — native console commands

Commit **a72dd6c6** records the completed native console conversion. Repaired C baseline
**15df0163 is committed and pushed**. All 47 selected C bodies are replaced; two
C files, 46 function interfaces, two globals and three private scratch buffers
retire. One C-to-Go dispatcher entrypoint remains for the quit dialog. See
[CONSOLE_COMMANDS.md](docs/porting/CONSOLE_COMMANDS.md).

Final qualification: 25 focused roots /539 tests; 18 captures /581 records match C.
Each default/server/highres gate passes 300 roots /20,221 tests without skips;
169 captures /32,019 records match across targets and C evidence. All four gates
have identical source. Static mapped-memory, three fresh builds/ABI, exact known
full-suite results and three headless scenarios pass. All tool sessions are joined;
source is editable. Evidence: `docs/porting/console-commands-native-qualification.json`
and `build/port-console-commands/native-final-{default,server,highres}` / `native-production`.
The formatter extension has an independently repeated C fixture patch and report.

The player-state batch now follows this qualified snapshot.
Ignored audit material in `build/port-player-state` selects 22 candidate routines /
619 body lines in GAME1.c. Two scalar getter/setter bodies appear orphaned; audit
remaining wrappers/callbacks before final scope. Consider the connected quit-menu
module where ownership fits. Player-state fixtures and the prerequisite C correction are now in progress.

Disk: verified asset deduplication and lossless compression of completed old evidence
reclaimed space. Restoration manifests remain local; original assets and the archive
are unchanged. See the console report for details. Do not replay consumed cleanup
scripts or installers. The untracked asset archive remains outside commits.


## Qualified parent — native client/server voting

The repaired C baseline **5572b500 is committed and pushed**. Thirty-five live
routines are native Go; 36 C bodies removed (one orphan), 31 interfaces retired,
five C entrypoints retained and 12 globals moved. Two neighboring GUI lifecycle
helpers were separately C-qualified from that committed baseline before conversion.
Their reproducible fixture patch and evidence are tracked; see [VOTES.md](docs/porting/VOTES.md).

All 22 focused roots pass and nine captures /387 records match C. Each target
passes 224 roots /43,580 cases without skips; all 169 captures /48,975 records match
C and each other. Four gates share unchanged 2,352-file source. Static checks and
fresh production pass: three builds/ABI, exact known full-suite results, gameplay,
save/load and flat regeneration. All sessions are joined; source is editable.
No user decision is pending. Evidence is under build/port-votes/native-* and in
[votes-native-qualification.json](docs/porting/votes-native-qualification.json).

Next: establish and qualify the console-command C baseline.
Read-only candidate audit: build/port-console-commands/audit-plan.md,
selection-draft.json and references.json. No candidate source changes yet.
Both voting install-native.py and extend-lifecycle.py are consumed; never replay.
Original assets/archive are unchanged. Completed voting C scenarios have verified
asset restoration manifests after reclaiming about1.55GiB; changed maps/saves,
reports and binaries remain. The isolated C lifecycle checkout remains for evidence.
Previous prefab installers/deletion passes are also consumed.

## Qualified parent — prefab script native conversion (8e8db9a1, pushed)

The repaired C baseline **5c83d11a is committed and pushed**. Eighteen live
algorithms are native, two orphan helpers removed, 24 C interfaces retired and
five counters moved to Go. No test-reference C is retained. Details and review
items: [PREFAB_SCRIPTS.md](docs/porting/PREFAB_SCRIPTS.md).

Final default/server/highres each pass156roots /2,057 including subtests, with
90 byte-identical captures /34,943 records matching C. Seventeen focused roots
and static checks pass. Fresh production passes three builds/ABI/symbol checks,
the exact known full-suite failure set, gameplay, save/load and flat regeneration.
All sessions are joined; source is editable. Evidence: native-final-* directories
and native-capture-audit.json under build/port-prefab-scripts.

## Qualified parent — prefab/map-runtime native conversion (d5f460cc, pushed)

Forty native algorithms replace the qualified C baseline **86654f0f**. Thirty C
interfaces, fourteen globals and two orphan bodies are retired; ten C exports and
two shared counters remain for live callers. No C reference algorithms remain.
The tile/wall payload and special-wall ownership corrections are documented for
review in [PREFAB_RUNTIME.md](docs/porting/PREFAB_RUNTIME.md).

All three targets pass **136 roots / 1,791 tests including subtests**, no skips;
**77 captures / 30,578 records** match the original C baseline byte for byte.
Fresh production builds/ABI/interfaces, exact known asset failures, headless
gameplay, save/load and flat-map regeneration pass. All gates share unchanged
**2,323-file source**; static mapped-memory checks pass. All sessions are joined.
See [native qualification](docs/porting/prefab-runtime-native-qualification.json).
C is **36,917 / 73 files / zero reference**, **−1,583**.

Next: qualify the connected prefab-script and
map-generation candidate: **20 C functions / 1,363 body lines**, covering all
server__script__file.c helpers, object callback/name remapping, pending-object
references, generation initialization and bounds ordering. Read-only selection,
original bodies, 194 references and audit notes are in build/port-prefab-scripts.
No next-batch source is installed. Reproduce suspected name-rewrite and script
buffer defects before freezing; use actual file/script/object owners and record
reversible corrections. No user decision is pending.

`install-native.py`, `repair-native-ownership.py`, `finalize-native.py` and prior
baseline finalizers are **consumed**. Do not replay ignored drafts. Preserve original
assets/archive and frozen captures. First two native compile attempts needed stale
fixture globals and the retired C-loader linker wrapper corrected; no captured
expectation changed. The latter is now a build-tagged Go fixture hook.

## Qualified baseline — prefab/map-runtime (86654f0f, pushed)


The connected **40-function / 1,371-body-line** baseline is qualified, with
**18 new frozen captures / 8,880 records**, repeated in separate C processes.
All three target sweeps pass **133 roots / 1,780 tests including subtests**, with
**77 identical captures / 30,578 records** per target. All gates share unchanged
**2,314-file source**. Static mapped-memory checks pass. Fresh production binaries,
ABI/interfaces, exact known asset failures, headless gameplay, save/load and flat
regeneration pass. See [PREFAB_RUNTIME.md](docs/porting/PREFAB_RUNTIME.md) and its
[C qualification report](docs/porting/prefab-runtime-c-qualification.json).

The waypoint allocation/disposal mismatch and two loader failure paths leaving
files open are repaired before freezing. C is **38,500 lines / 74 files / zero
reference C** (+2 temporary cleanup statements). No selected C algorithm is ported
yet. The initial production run stopped on a reused scenario directory name;
corrected names resumed using hash-verified fresh build/suite evidence.

Next: port and retire the selected C bodies/interfaces, compare every frozen
capture unchanged, qualify three targets and fresh production, record C LOC,
commit/push and continue. Ten exports appear necessary for external C callers.
Native design and an uninstalled initial state/helper draft are in
build/port-prefab-runtime/native-design.md and native-state.go. Review drafts before
use. Audit cached tile/wall payload release together with secret-wall data's
back-reference and partial-placement ownership; do not simply free a record still
referenced by the world. The known leak is documented for native ownership work.

freeze-captures.py, finish-c-baseline.py and finish-prerequisite.py are consumed;
never replay them. Do not rerun old quest/monster installers. No user decision is
pending. Preserve original assets/archive and frozen captures.


Disk maintenance reclaimed **3.092 GiB** from six completed prerequisite/quest
scenario copies after SHA-256 comparison with original assets. Changed maps/saves,
reports, binaries and original assets/archive remain. Per-run restoration manifests
are saved; build/port-prefab-runtime/deduplicate-completed-assets.py --apply is
consumed and must not be repeated. About **8.1 GiB** was free after cleanup.

## Qualified parent — quest-progress native conversion (ae6fbe52, pushed)

C baseline **1c915175 is committed and pushed**. Twenty-four live functions are
in Go and the orphan getter is removed. Final focused run passes twelve roots,
including all **11 original captures / 3,156 records** unchanged and the new
bounded-name/partial-save contract. All three broader targets pass without skips,
with **227 identical captures / 150,824 records** each. Fresh production passes
all builds/ABI/interfaces, exact known asset failures, headless gameplay, save/load
and flat-map regeneration. All four gates share unchanged source; sessions are
joined. See [QUEST_PROGRESS.md](docs/porting/QUEST_PROGRESS.md) and its native
qualification report for counts and timings.

This retires **24 interfaces** and one C global, retaining only three actual C
entrypoints. The list is Go owned. Reversible name/file limits and the separately
rounded boss-health conversions are documented for review. No goldens changed.
Production C: **38,498 / 74 files / zero reference**.

Next: qualify the connected prefab/map-runtime candidate under
build/port-prefab-runtime: **40 live functions / 1,369 body lines** plus the orphan
waypoint setter. Reuse actual population/file/group/waypoint owners; finish the
manual callback, numeric table and source-reference audit before baseline fixtures.
No next-batch source is installed. The orphan nox_strnicmp lost its last caller
in quest reset and can be audited alongside the next cleanup without rerunning
this completed qualification solely for eleven source lines.

All quest freeze/install/retirement/finalization scripts are consumed; never replay.
Original assets/archive remain intact. No pending user decision.

## Qualified parent — monster-control native conversion (4a8d73ed, pushed)

C baseline **77cc1fd2 is committed and pushed**. All forty live functions are in
Go and three proven orphans are removed. Final focused run passes sixteen roots,
including all **15 original captures / 20,997 records** unchanged and the new
bounded-parser/allocation contract. All three broader targets pass without skips:
**176 identical captures / 135,370 records** each. Fresh production passes all
three builds/ABI/interfaces, exact known asset-suite failures, headless gameplay,
save/load and flat-map regeneration. All four gates share unchanged source;
all sessions are joined. See [MONSTER_CONTROL.md](docs/porting/MONSTER_CONTROL.md)
and monster-control-native-qualification.json for exact counts and timings.

This retires **62 interfaces** and four C globals, retaining only two actual C
entrypoints. Reversible parser limits and rejected-record cleanup are documented
for review. No goldens changed. Production C: **39,193 / 74 files / zero reference**.

Next: prepare a connected quest-progress C baseline covering journal variables,
serialization and remaining quest-stage/boss-spawn owners. Caller and fixture
review is underway; no next-batch source is installed. The conservative whole-src
C reachability proposal under build/port-reachability is evidence for manual
review only, not authorization for blindly deleting its candidates.

All monster freeze/install/resume/retirement/finalization scripts are consumed;
do not replay. Original assets/archive are intact. No user decision is required.

## Qualified parent — spatial-targeting Go conversion (0bbaba1d, pushed)

C baseline **76ba09f8 is committed and pushed**. The native conversion replaces
**11 functions / 554 body lines**, retains one ray export for a real C spell caller,
and retires eleven obsolete interfaces plus one C global. Go owns cursor state;
tracing and AI prediction no longer allocate temporary C point records.

Focused native-2 passes **13 roots** (0.299s): eleven spatial captures / **26,142
records** unchanged from C, existing curve captures and the opaque-token regression.
All three broader targets pass with no skips: **706/705/706 roots** and
**257 captures / 157,554 records** each. Fresh production passes all gates,
including the exact known asset failures, gameplay, save/load and flat regeneration.
All four gates share unchanged source; all sessions are joined. See
[SPATIAL_TARGETING.md](docs/porting/SPATIAL_TARGETING.md) and
spatial-targeting-native-qualification.json for counts, timings and binaries.

Qualification caught an older curve callback bridge treating integer userdata as
a pointer. A deterministic regression reproduced it before correction; typed point
pointers plus an integer token now preserve the callback ABI. No goldens changed.
Production C: **40,218 / 74 files / zero reference (−559)**.

Next: complete monster-control C contracts. The proposal
covers **40 live functions / 1,001 body lines** and three proven orphans / 31 lines.
Comment-only name matches are not callers. Complete dynamic/preamble reachability
and shipped table review, then reuse AI owners for C contracts. Proposal/reference/
table audits are under build/port-monster-control. No next-batch source is installed.

**Spatial install-native.py is consumed.** Work from src; never reinstall drafts
or rerun old cleanup scripts. Current artifacts: build/port-spatial-targeting/
native2-{default,server,highres}, native-production, native-focused-2. Original
assets/archive are preserved. No user blocker or pending product decision.

## Qualified parent — world-motion Go conversion (6541718f, pushed)

Corrected-C baseline **e7174c35 is committed and pushed**. Native conversion now
passes all gates: **30 C functions replaced**, six required exports retained,
33 obsolete interfaces and two C list globals retired. Go owns the list heads and
velocity type cache. No C algorithm remains solely for tests.

Native-3 passes **25 roots / 24 captures / 9,633 records** unchanged from C
(0.373s). Default/server/highres pass **690/689/690 roots** and
**43,124/43,123/43,124 tests including subtests**, no skips. All **245 captures /
114,004 records** match C. Durations: **278.76/386.65/330.78s**. Fresh production
passes in **372.74s**: three binaries/ABI, exact known 1,553 asset failure entries
(15 pass / 3 fail / 32 no-test packages), gameplay, save/load and flat regeneration.
All four gates share an unchanged **2,244-file source manifest**. All sessions are
joined; no active build or user blocker. See WORLD_MOTION.md and
world-motion-native-qualification.json.

The sentry unlink and disabled one-shot trigger corrections are carried forward.
Native translation errors (delayed-delete owner and alloc.New initialization)
were caught and fixed without changing goldens. Compiled-C rounding boundaries
are preserved; see the report and DECISIONS.md.

Next: complete the connected spatial-targeting C baseline.
Read-only proposal: **11 functions / 554 C body lines**, including projectile and
cursor candidates, wall normals and quadrant helpers. It removes the remaining
spatial C calls from motionTrace and its temporary C allocation. Proposal,
reference audit and table audit are in build/port-spatial-targeting. Supply the
shipped door table at 0x587000:196184; reuse collision/index/wall/player owners.
Geometry test adapters and contracts are installed. This smaller connected scope is preferable
to adding unrelated player-death or parser code just to reach a LOC target.

**install-native.py is consumed.** Native drafts are installed; work from src.
Do not rerun old fixture integration or cleanup scripts. C/native captures and
source manifests remain in build/port-world-motion; preserve the untracked archive.
Disk cleanup reclaimed 14.96 GiB of rebuildable cache; about **18 GiB** remains free
following qualification. Original assets, module downloads and evidence remain.

## Qualified parent — collision-core Go conversion (950f8f22, pushed)

C baseline **bc4a791e is committed/pushed**. The native conversion replaces
**23 live functions / 993 body lines**, retires eighteen function interfaces and
eight private C globals, and retains ten exports for actual C callers. Queue indices,
force coefficients and type caches are Go-owned; Hit records retain the real fixed
C-backed pool. Six production Go files implement the batch.

Focused native-3 passes **16 roots/captures / 22,848 records** unchanged from C
(0.378s); static native-2 passes. Broader default/server/highres pass
**660/659/660 roots**, **42,918/42,917/42,918 tests including subtests**, zero skips.
**221 captures / 104,371 records** match C across all targets. Durations:
**275.41/387.91/328.61s**. Fresh production passes **375.82s**: three binaries,
ABI/interfaces, exact known 1,553 asset failures (15 pass / 3 fail / 32 no-test),
gameplay, save/load and flat-map regeneration. All four gates share an unchanged
**2,219-file source manifest**. All sessions are joined; no active job or blocker.

C: **41,886 / 74 files / zero reference (−1,096)**. Client SHA:
71758f918f3d4ae516cbb3de252bd056f322013095e35da334f087c03753cc33.
See docs/porting/COLLISION_CORE.md and the scope/batch/test manifests. Artifacts:
build/port-collision-core/native-{default,server,highres,production}, native-audit.json,
native-interface-audit.json. All drafts/install/freeze scripts are consumed.
Never reinstall or regenerate frozen root expectations.

No intentional gameplay rule change. Review preserves compiled-C arithmetic widths,
including shaft subtraction; original circle/box interior behavior; unsigned angular
stop modulo; and pointer-typed temporary-normal callback lifetime. The first link
check caught springs.go's still-required cgo import; restored. The C baseline's
default driver had misnamed inherited capture paths after tests passed; independent
corrected-manifest and cross-target audit verified every capture. Reports retain
both issues transparently.

Next: world-motion baseline described above.

Disk cleanup reclaimed **4.638 GiB** from nine newly completed polygon/geometry
run copies. Each removed file matched the original asset hash; per-run manifests
record restoration. **deduplicate-completed-assets.py --apply is consumed; never
repeat deletion mode.** Original assets/archive, changed files and results remain.
Free space afterward: about 12 GiB. See completed-assets-{plan,audit}.json and
completed-assets-applied.log under the new batch directory.

## Qualified parent — world geometry/collision Go conversion (011fcb73, pushed)

C baseline **9a945523** and supplemental player-wall contract **7d363edc** are
pushed. The conversion replaces **32 functions / 1,385 corrected C body lines**,
keeps twelve exports for actual C callers, retires twenty function interfaces and
two private data definitions, and routes Go callers directly. Two force coefficients
remain shared with C readers. Seven production Go files implement the batch.

Focused native-7 passes **17 roots / 16 captures / 15,902 records** (0.352s);
static-2 passes. Broader default/server/highres pass **544/543/544 roots**,
**42,760/42,759/42,760 tests including subtests**, no skips; **205 identical
captures / 81,523 records** match C. Durations **251.18/375.39/290.70s**.
Fresh production passes **398.08s**: three binaries/ABI/interfaces, exact known
1,553 asset failures (15 pass / 3 fail / 32 no-test packages), gameplay, save/load
and flat-map regeneration. All four gates share an unchanged **2,198-file source**.
All sessions are joined; no active tests or user blocker.

C is **42,982 / 74 files / zero reference (−1,414)**. Client SHA:
3b0754cf3af539d8447c294af8a64a48429441760495861838d69429233a7737.
See [WORLD_GEOMETRY.md](docs/porting/WORLD_GEOMETRY.md), scope and batch manifests.
Artifacts: build/port-world-geometry/native-{default,server,highres,production},
native-audit.json and native-interface-source-audit.json.

Arithmetic review preserves compiled C's x87 store/reload boundaries; declared
float locals sometimes remain wide. ShapeBox.Calc is not interchangeable, so the
port uses a compatible corner calculation. floatToInt32 is reused. Frozen root
expectations unchanged. All native drafts, install-native.py and finish-width-review.py
are consumed; never reinstall. A local comparison helper now rejects missing
captures instead of reporting an empty set as matching.

Next: prepare the collision-core C baseline described above. The read-only next candidate has **23 live functions / 993 C body lines**
in collision event queues, dispatch, activation lists and remaining contact
geometry. See build/port-collision-core/{proposal.json,combined-reachability.json,plan-draft.md,c-disassembly}.
Reuse actual owners and
capture the separate shipped circle/box coincidence table at 0x587000:289928.
At most three heavy jobs; source build/baseline/env.sh for every Go command;
GOMAXPROCS=2 (highres=1), GOMEMLIMIT=768MiB. Original assets/archive untouched.

## Qualified parent — map-polygon native conversion (10d3294d, pushed)

C baseline **afdb26e7 is committed/pushed**. The native conversion replaces
**33 functions / 1,000 C body lines**, retires 26 function interfaces and three
private globals, and retains seven thin C exports. All 16 focused groups pass
(0.223s); static checks pass. Broader default/server/highres pass **269/268/269
roots**, **34,988/34,987/34,988 tests including subtests**, no skips. **87 identical
captures / 51,433 records** match C. Durations **108.70/215.00/152.76s**.

Fresh production passes **371.15s**: three builds/ABI/interface checks, exact known
1,553 asset failures (15 pass / 3 fail / 32 no-test), options gameplay, save/load and
flat-map regeneration. All four gates share an unchanged **2,180-file source
manifest**; all sessions joined. Client SHA:
0a9413f9dfb01e79c28603eecc3bed74967b537e222fc069a629de873fd399f2.
C: **44,401 / 74 files / zero reference (−1,072)**.

See [MAP_POLYGONS.md](docs/porting/MAP_POLYGONS.md), map-polygons-batch.json,
map-polygons-tests.txt and map-polygons-scope.json. Artifacts under
build/port-polygons/native-{default,server,highres,production}, native-audit.json
and native-interface-audit.json. All drafts/install/finish-docs scripts are consumed; never
reinstall them. Frozen root tests/expectations were not changed. Nearest-vertex
width review, ray parity, player cache quirks and allocation ownership decisions
are documented for review. No active sessions or user blocker.

Disk cleanup before the polygon batch reclaimed **6.184 GiB** from twelve completed
panel/configuration runs. Per-run manifests permit restoration; never repeat
deletion mode. Original assets/archive and modified files remain intact.

## Qualified parent — server-configuration native conversion

Corrected C baseline **ce1399b2** is committed/pushed. Native conversion replaces
**64 live functions / 816 corrected C body lines**, removes two disabled helpers
and their no-op callers, and retires eight private globals / 41 function interfaces.
Thirty-two thin exports remain for actual C callers. Go callers and callbacks use
native functions. See [SERVER_CONFIG.md](docs/porting/SERVER_CONFIG.md).

All 26 focused roots pass (2.127s); static preflight passes (0.398s). Broader
native default/server/highres pass **535/534/535 roots**, **58,651/58,650/58,651
leaves**, no skips. **240 identical captures / 76,002 records** match corrected C.
Durations: **178.11/283.13/222.88s**. Fresh production passes in **377.75s**:
three builds and ABI/interface checks, the exact known 1,553 asset failures
(15 passing / 3 failing / 32 no-test packages), options gameplay, save/load and
forced flat-map regeneration. All four gates share an unchanged **2,167-file source
manifest**; all sessions joined. Client SHA:
7c3a16ffb2531aac2c9b6a4368380b41982b89de2af36483c5c5b9d30017ef95.
C is **45,473 / 74 files / zero reference C (−920)**.

Artifacts: build/port-server-config/native-{default,server,highres,production},
native-audit.json, native-interface-audit.json. Manifests: server-config-batch.json,
server-config-tests.txt and server-config-scope.json. All fixture/native drafts,
installation and export-preparation scripts are **consumed**; never recopy them.
Frozen expectations were not changed during conversion. The pixel capture caught
and corrected use of a blended background where the picker requires opaque fill.

The conversion is committed/pushed and duplicate cleanup is complete. Preliminary combined
scope: **33 functions / 1,000 C body lines** (core lifecycle/lookup/render/events plus
map serialization). Read-only scope/audit/plan are in build/port-polygons;
use candidate-with-serialization.json, not the smaller core-only candidate.json.
The reachability closure includes all 33 functions. Polygon C fixtures are in progress; no native conversion is installed. Preserve original assets/archive; no user decision/blocker.

## Qualified parent — server-panel conversion

This checkpoint converts **57 live functions / 2,098 corrected C body lines**.
Corrected C baseline **c2bd3a34** was committed/pushed before translation; parent
server-options conversion is **f5737e6a**. All five remaining servopts C units,
the 37-line orphan and 26 private C globals are removed. Nine thin C exports remain;
49 function interfaces and two obsolete blob callback registrations are retired.
Go callers/fixtures invoke Go directly. Advanced refresh callbacks are native;
the shared online-mode flag remains. See [SERVER_PANELS.md](docs/porting/SERVER_PANELS.md).

All **22 focused roots** pass with unchanged expectations (3.142s); static memory
preflight passes (0.378s). Broader default/server/highres pass **509/508/509 roots**,
**58,212/58,211/58,212 leaves**, no skips. **215 identical captures / 74,168 records**
match both corrected C and the preceding native corpus. Durations:
**186.42/297.11/228.63s**. Fresh production passes in **374.28s**: three builds and
ABI/interface checks, exact known 1,553 asset failures (15 pass / 3 fail / 32 no-test),
options gameplay, save/load and forced flat-map regeneration. All four gates share
one unchanged **2,148-file source manifest**; all sessions joined. Client SHA:
81fdbc91cebc252b4ef2e87a5edac4462707163affa62f750f8ecc40685a3527.
C is **46,391 / 74 files / zero reference C (−2,364 from corrected baseline)**.

Artifacts: build/port-server-panels/native-{default,server,highres,production},
native-audit.json, interface-audit.json. Manifest: server-panels-batch.json;
selection: server-panels-tests.txt. All drafts/install/freeze scripts are consumed:
never reinstall them or regenerate frozen expectations.

Next candidate: server configuration, admission-list persistence and rule picker
(66 functions / 836 C body lines before reachability cleanup). Read-only scope,
external-reference audit, test plan and uninstalled initial core fixture drafts
are in build/port-server-config. Review before installing; no baseline tests for
this next batch are installed yet. Audit the no-op sub_416690 and its dead private
setter sub_4164F0; preserve actual behavior rather than restoring disabled code.
The rule-picker missing-resource path needs a child-process contract. Continue
corrected-C baseline/captures/qualification/commit, then native conversion.
No blocker or user decision. About 17 GB free; preserve original assets/archive.

### Qualified parent — server options

Native conversion **f5737e6a is committed/pushed**. Corrected C baseline
**55ccc2b7 was committed/pushed before translation**. See
[SERVER_OPTIONS.md](docs/porting/SERVER_OPTIONS.md) for scope and review decisions.
All **33 functions / 1,256 live C body lines** are native; the 28-line orphan,
guiserv.c translation unit, 16 private globals and C mode table/cache are gone.
Twelve C interfaces remain (ten outside callers and two GUI tooltip callbacks);
22 function interfaces are retired. Go callers and fixtures call Go directly.

All **35 options roots / 847 leaves** pass with **23 unchanged C captures /
1,322 records**, in 4.852s. Broader default/server/highres pass **487 / 486 / 487
roots**, **57,405 / 57,404 / 57,405 leaves**, no skips, in **204.41 / 329.54 /
230.07s**. All **193 captures / 61,102 records** match, including both the corrected-C
baseline and preceding team-UI corpus. The server excludes client-only occlusion.

Fresh production passes in **416.35s**: three builds/ABI/interfaces, exact known
asset failures (1,553 entries; 15 pass / 3 fail / 32 no-test), gameplay, save/load
and forced flat-map regeneration. All four gates share one unchanged **2,128-file
source manifest**. Static memory and interface audits pass; all sessions joined.
Client SHA: f1641ca43485c15532894867eedcaa2bf7aab0f7d061d0d2a976b0e5c35fb058.
C is **48,740 / 79 files / zero reference C (−1,441)**.

Artifacts: build/port-server-options/native-{default,server,highres,production},
native-audit.json and interface-audit.json. Manifest: server-options-batch.json;
selection: server-options-tests.txt. Frozen expectations are unchanged. First-run
hit-test and negative-edit differences were fixed against C. Apply preserves the
selected-name alias across settings copies. All native drafts/install scripts are
consumed or stale: never recopy them.

Completed-run asset deduplication is finished: 12 runs, 19,648 identical files
(6.18 GiB logical duplicates), with per-run restoration manifests. Original assets,
changed data and screenshots remain. About 21 GB is free.

Next scope is underway: server panels. Read-only draft scope: build/port-server-panels/candidate.json,
58 functions / 2,120 body lines across five remaining servopts files and related
GAME2/GAME3/GAME3_1 blocks. plan-draft.md records owners/tests; external-refs.json
is advisory and requires review. Candidate audit predates the newly installed scalar fixture. No blocker or
user question. Preserve archive/assets.

### Qualified parent — team UI

Corrected C baseline **e86e6b21** was committed/pushed before conversion. All
**35 live functions / 943 body lines** are native. Two C translation units and five
private C globals are gone; 14 Go-backed interfaces remain / 21 retired. Go callers
and fixtures call Go directly. Flagball root stays shared with unported C callback
sub_456140. Row storage preserves its 72-byte layout and address-stable allocation;
list algorithms and ownership are Go. No C algorithm remains solely for testing.

All three broader sweeps pass without skips: **415 / 414 / 415 roots**,
**56,434 / 56,433 / 56,434 leaves**, **170 identical captures / 59,780 records**,
in **292.34 / 288.66 / 354.17s**. The server excludes the client-only occlusion root.
Fresh production passes in **301.60s**: all three builds/ABI, exact known
asset-suite failures (1,553 entries; 15 pass / 3 fail / 32 no-test), gameplay,
save/load and flat rendering. Forced expansion removes 51 maps and regenerates
one exactly. All four gates share the same unchanged **2,097-file source manifest**;
all sessions are joined. Client SHA-256:
`226200f02134e47708ec3239a20c72d185ca39e7931e9fd05cf665c4128a2bc1`.
Artifacts: build/port-team-ui/native-{default,server,highres,production-final},
native-audit.json and interface-audit.json. Manifest: team-ui-batch.json;
selection: team-ui-tests.txt. Frozen original-C expectations are unchanged.

The first production attempt had a manifest-only error: the shared Flagball C
variable was in the Go-export list. Corrected to retained_c; final production
passes without source or golden changes. Five original-C UI fixes and ASCII-only
name folding are documented in TEAM_UI.md / DECISIONS.md for review.

Next: complete the server-options C fixture/contract matrix and qualify baseline.
Read-only
scope audit in build/port-server-options selects **33 live functions / 1,252 body
lines**, plus a **28-line orphan** (sub_457FE0, called only in two literal if(0)
blocks). plan-draft.md describes the test/owner matrix. Candidate reachability and
state audits are advisory; numeric references and callbacks require manual review.
Initial mode fixture drafts there are consumed/installed. No user decision blocks work.

Source build/baseline/env.sh for Go. Maximum three heavy jobs; do not edit source
while qualification reads it. Team UI ignored drafts/install/freeze scripts are
consumed or stale: never recopy them. Preserve original assets, module cache and
archive. Prior Go-cache pruning and scenario-dedup manifests are consumed. Around
22 GB remains free; the latest team-UI scenarios have not been deduplicated.

### Completed — team runtime

Corrected C baseline `0b4b853c` and native conversion `6963545e` are pushed: 36 live functions ported, one 62-line orphan removed, 25 C
interfaces retained / 12 retired, two private globals moved to Go. Frozen hashes
are unchanged. See [TEAM_RUNTIME.md](docs/porting/TEAM_RUNTIME.md).

Default/server/highres pass **390 / 389 / 390 roots**, **56,105 / 56,104 / 56,105
leaves**, **151 groups / 58,911 records** identical, no skips. Times **158.90 /
261.03 / 197.54s**. Fresh production passes in **369.47s**: three builds/ABI, exact
known asset failures, gameplay, save/load and flat rendering; 51 maps removed, one
regenerated. All four final gates share an unchanged 2,083-file manifest; all joined.

One broader-test correction: the legacy score sender must go directly to the
reliable queue, while the existing public Go setter uses a replaceable send hook.
Preserve both owners; new dispatch contract and original objective-score hashes pass.
PORT.md records the lesson for future Go reuse. No gameplay rule changed here.

Manifests: team-runtime-c-batch.json / team-runtime-batch.json. Artifacts under
build/port-team-runtime; final sweeps native-{default,server,highres}-2 and
native-production. Ignored fixture/native drafts and install/rewrite scripts are
consumed/stale; do not rerun them. Preserve original assets and archive.

### Last completed qualification

Quest conversion `c1f80856` and corrected C baseline `bacee82b` are committed/pushed.
The native match/roster conversion is qualified; this checkpoint accompanies its
conversion commit.
It removes 28 C functions (one proven orphan), two private globals and trailing
separators: **773 physical C lines**. Ten actual C interfaces remain; eighteen
retired symbols have no remaining source references. No reference C is retained.

Frozen baseline: **24 roots / 10,027 leaves / 21 groups / 10,034 records**.
Native focused checks pass (136.42s). Broader default/server/highres checks pass
(115.52s / 225.84s / 153.84s): **283 roots / 52,837 leaves / 127 groups /
55,862 records**, without skips; every capture group matches across targets.
Fresh native production passes (388.13s): three builds and ABI checks, exact known
asset failures (1,553 entries; 15 pass / 3 fail / 32 no-test packages), gameplay,
actual save/load and flat rendering. The flat scene removes 51 loose maps and
regenerates the selected map exactly. All four gates share 2,060 source fingerprints.
Artifacts: build/port-match-roster/native-{focused-3,default,server,highres,production},
native-coverage-audit.json; full detail in docs/porting/MATCH_ROSTER.md.
All match/roster qualification processes have completed and joined.
The new team probe has its own source freeze above.

The C baseline corrected roster/settings padding, the Flagball no-winner crash,
and Go's player identifier offset (2096). Native bring-up corrected a helper name
and unordered wall-message metadata; no frozen expectations changed. Preserve the
unconditional C highres protocol version pending separate review.

Historical next step after match/roster: audit connected team management.
Read-only candidate inventory is build/port-match-roster/next-team-candidates.json:
42 bodies / 1,096 C lines covering membership, team messages and objective setup.
This is a candidate list, not an accepted scope or frozen baseline. Verify live
callers/callbacks and use actual membership/object owners before capturing cases.
Full player-arrival orchestration and GUI settings remain separate owners.

install-native.py and freeze-baseline.py are CONSUMED; old ignored drafts are stale.
Never recopy them or regenerate expectations to hide a difference. Preserve original
assets and untracked nox-iso-from-archive-org.7z. Completed scene deduplication keeps
per-scene restoration manifests; never repeat consumed cleanup scripts.
Source build/baseline/env.sh for Go. No source edits during builds; no user blocker
and no active agent delegation.

<!-- /current-checkpoint -->



## GitHub backup and recovery

Qualified changes and these documents are pushed to `dev` in
[dbenamy/opennox](https://github.com/dbenamy/opennox). The explicit SSH push URL is
`git@github.com:dbenamy/opennox.git`; credentials are provisioned separately.

[Recovery instructions](docs/porting/RECOVERY.md) and the tracked
[warrior scenario](docs/porting/warrior-smoke.yaml) preserve setup and validation
steps if the ignored `build/` directory is lost. Raw captures, logs, binaries,
screenshots, saves and game assets are local artifacts, not committed source.

## Historical notes

The sections below preserve earlier measurements and decisions. Their C counts,
run status and then-next actions describe the recorded stage, **not current work**.
Use the current checkpoint above, the [C size history](docs/porting/C_LOC.md) and
individual batch reports for subsequent progress. This log is not an exhaustive
chronology of every completed batch.

### Historical status after infrastructure repairs

Plan and infrastructure changes are committed. Screenshot checks now fail
reliably, automatic test writes are isolated, and stale API/vet failures are
repaired. Full default suite with assets is down from seven failing packages to
three: blobs tooling, renderer goldens and audio goldens. None are suppressed.
No C-to-Go conversion had begun at that stage. The following measurements record
that infrastructure baseline; later conversion reports supersede its progress status.

### Local artifacts and reproducibility

Everything under build/ is ignored by Git. These artifacts persist in this VM,
not in commits. See build/baseline/README.md for details and commands.

- build/baseline/env.sh: target flags, i386 pkg-config and build/cache Go caches.
  Pinned modules downloaded successfully; no go.mod/go.sum changes.
- build/baseline/bin/{opennox,opennox-hd,opennox-server}: pristine baseline builds.
  binaries.json records SHA-256, ELF32/Intel 80386 headers and help exit codes (0).
  environment.json records toolchain, revision, flags and source archive hash.
- logs/build-pristine.log: all targets built successfully in 7.66 seconds warm,
  peak RSS 455332 KiB. No cold-build or gameplay performance claim.
- The user-confirmed asset archive actually contains Nox.wsquashfs, with an
  installed drive_c/Nox tree. Extracted only that subtree to
  build/assets/extracted/drive_c/Nox. Original archive untouched. Test runs use
  separate copies and empty save directories; no bundled executables were run.
- Full JSON test outputs, exit files and build graphs are in build/baseline/logs.
  test-summary.json gives package/test event counts.

### Test outcomes: baseline is red

All three full suite variants were attempted with assets, GOARCH=386,
CGO_ENABLED=1 and allowed C flags. Each exits 1:

| Variant | Passed packages | Failed packages | Packages without tests |
| --- | ---: | ---: | ---: |
| default | 11 | 7 | 30 |
| server | 10 | 7 | 30 |
| highres | 11 | 7 | 30 |

One explicit skipped test with assets: internal/blobs/TestSplitBlob. Failed
package counts include compilation/vet failures. Test/subtest counts are not
independent scenarios.

Failures:
- cmd/noxmovie uses old sdl.New signature (missing logger).
- internal/netstr tests have obsolete callback types and missing netmsg import.
- internal/offalign tests pass uint where uintptr is required.
- Root package vet rejects log.Printf's %w in maps.go.
- internal/blobs has obsolete memmap.go path and a formatter failure.
- noxrender image/particle goldens fail. Particle tests fail even without assets.
  They hash encoded PNG, not just pixels; root cause is not yet established.
- Audio PCM hashes mismatch supplied data. Asset/decoder/toolchain differences
  have not been isolated; do not simply regenerate the expected hashes.
- server-tag legacy/dialog expects "empty" and receives an empty string.

Important: some existing tests MODIFY SOURCE. The first full run changed ten
tracked files and emitted token files. Saved its diff to
logs/test-generated-changes.patch, restored only those generated edits, moved
emitted files to logs, and rebuilt all binaries from clean source afterward.
Subsequent suites ran in build/baseline/test-checkout. That disposable copy was
reused across variants; future runs must recreate it per variant to avoid any
cross-run contamination. The working checkout is clean apart from planning docs
and the original archive.

### Gameplay scenario and oracle

- Existing src/e2e*.go harness supports YAML input, simulated time, platform RNG,
  screen comparisons and save hashing. Use it before inventing another harness.
- build/baseline/e2e/warrior-smoke.yaml: title → Solo → new warrior → war01a →
  captain dialogue → short walk → screenshots → clean quit.
- run-scenario.py RUN_NAME [BINARY] creates a fresh data/save copy, uses Xvfb
  1280x960 and OpenAL null backend (audio enabled), enforces 120s timeout, and
  records command/exit/elapsed/logs in runs/RUN_NAME. Run names must be new.
- runs/repeat-a and runs/repeat-b both exited 0 in about 40 seconds. Visually
  inspected gameplay and moved-player frames. Decoded NRGBA pixels match exactly
  at both checkpoints (frame-comparison.json). This covers this short scenario,
  not broad gameplay correctness or physical audio/display quality.
- compare-frames.go and its binary perform independent decoded-pixel comparison.
  It correctly rejects a deliberately replaced frame (comparator-negative/).
- Why independent at baseline? The original Screen passed nil to e2eError on a
  mismatch and auto-created missing goldens. This was repaired in 3575f443; current
  checks require explicit golden updates and reject mismatches.
- HD client also completed the same fresh scenario with exit 0 in 39.3 seconds.
  Both captured frames match the standard-client pixels at this 1024x768 game
  resolution (hd-comparison.json); higher resolutions are not covered.
- Save comparison: war01a.map bytes match across runs; Player.plr differs in both
  WORKING and AUTOSAVE. See save-comparison.json. Cause and save-load compatibility
  are unverified. Missing-script-object warnings also remain in baseline logs.

### C inventory started

Saved go list -deps -json graphs for default/server/highres. Local compiled C:
152 for clients (148 legacy + 3 cnxz + 1 ail), 151 for server (no ail C file).
This is compiled translation-unit inventory, not linked/reachable/active logic.
Indirect callbacks, retained symbols and subsystem classification remain pending.

### Follow-up identified at the initial baseline

1. Diagnose PNG/PCM goldens without masking real behavior changes. Investigate
   Player.plr nondeterminism and add a save-load scenario.
2. Use build graphs for a bounded dependency audit; choose a cohesive C leaf and
   establish differential coverage before porting it.

Do not call the suite green. Dedicated-server map/tick scenarios, multiplayer,
replay validation, sanitizer compatibility and performance work remain pending.

### Infrastructure changes after baseline

- Plan/checkpoint committed as af08739d; repository-local Git author configured
  from the user's supplied identity.
- Screenshot oracle: moved comparison into internal/e2etest. Normal checks fail
  for missing goldens, pixel differences and dimension differences; only explicit
  NOX_E2E_OVERRIDE updates goldens. Decoder formats/origins are normalized, actual
  and diff frames are retained, input buffers are not mutated, and file errors
  propagate. Focused 386 tests pass; all three targets build into build/infra-bin.
  A real headless client with a deliberately wrong-size golden exits 2 with a
  screen mismatch (build/screen-negative), confirming integration fails visibly.
- Source-rewriting blobs tests now copy src into t.TempDir and restore the global
  tool path afterward. Token diagnostics also go into t.TempDir. Package execution
  left tracked engine files untouched. Token tests pass; blobs retains the known
  obsolete memmap.go-path and formatter parse failures (tests-isolated.log).
  Those failures are not skipped or relabeled as passing.
- Full-suite follow-up identified internal/noxfactor/TestNoxFactor as another
  source rewriter. Its generated enum substitutions were saved to
  logs/noxfactor-generated.patch and restored; the test now uses a temporary
  source copy. Default rendering/heatmap PNGs also use temporary output paths.
  Focused noxfactor, memmap and primitive-render tests pass on 386 and leave
  engine source untouched (tests-isolation-followup.log).
- Updated obsolete movie-player SDL call, offalign fixture integer types, and
  maps.go's invalid logging format. Reworked the netstr test for current typed
  callbacks, synchronous server binding, loopback payload verification, a bounded
  wait and goroutine cleanup. Target tests pass, including five repeated netstr
  runs; supplementary amd64 netstr race test passes. Root package passes default
  vet/compilation. No production network behavior changed.
- Full default suite with assets now has 14 passing packages, 3 failing packages,
  32 packages without tests, and no compilation/vet failures. Remaining failures
  are internal/blobs, client/noxrender, and legacy/client/audio/ail, already known
  from baseline. See logs/tests-infra-default.*. Final test-isolation follow-up
  above was validated separately after this full run.
- Tagged follow-up: the movie command lacked !server even though its movie library
  is client-only. Added that matching constraint; it is omitted by server go list
  ./... and still compiles for default/highres. Focused server root/netstr/offalign/
  e2etest checks pass; highres focused checks pass. Initial variant logs preserve
  this discovered setup failure (their shell's final exit reflected highres only).
- Final production validation: all three targets build after infrastructure/API
  changes (logs/build-infra-final.log). The patched standard client completed
  the fresh warrior scenario against both pre-existing baseline goldens with
  NOX_E2E_OVERRIDE=false and exit 0 in 38.9s (runs/infra-positive). This exercises
  the repaired in-engine screenshot oracle's success path; screen-negative
  exercises its failure path. Only the original media archive remains untracked.

Next session should start with the bounded active-C dependency audit and select
an independent leaf whose relevant tests pass. The remaining blob-tool and
render/audio failures need diagnosis before touching those areas, but are not a
blanket blocker for unrelated conversions. Baseline Player.plr differences also
remain unexplained. No C implementation has been replaced yet.

### Bounded failure diagnosis — 2026-09-10

- Blob formatter: combining two dynamic Go offset terms dropped the joining +,
  e.g. uintptr(x)*13+71276+uintptr(y) became invalid Go. Added regression cases
  covering operand order, nested sums, subtraction and zero offsets, checking
  parsing, idempotence and independently evaluated arithmetic. Tests fail before
  the fix and pass afterward; TestFormatAccesses now passes on the source copy.
- Remaining ReadBlobs failure is a separate obsolete storage-format assumption:
  it expects root memmap.go/memmap.c/GAME_data.c, while current code uses legacy
  shims plus embedded .dat files and generated pointer initialization. Simply
  changing paths to cgo_blobs.c would incorrectly treat initialized data as zero.
  Updating the split/write tool to that format is separate work, not needed for
  the initial dependency inventory. Do not use it to rewrite current blobs yet.
- Particle rendering diagnosis: 386 and amd64 produce identical pixels for all
  six particle cases. Current PNGs re-encode identically with Go 1.19.13, 1.23.12
  and 1.26.0. The old color library used RGB max 248; the pinned current library
  expands that value to 255. Applying the old expansion reproduces all six
  original PNG goldens exactly. Migrated ONLY those independently verified cases
  to dimensioned, little-endian framebuffer-word SHA-256 hashes; pixel/size
  mutation and subimage-stride checks pass, as do particle tests on 386/amd64.
  No production rendering behavior or sprite goldens were changed.
- Sprite diagnosis: the exact golden-era revision c62202f5 passes the sampled
  APA00001/default case with the same assets. In a disposable current dependency
  copy, restoring only historical color/rgba5551.go makes the entire sprite test
  matrix pass. Export-only normalization was insufficient because the conversion
  also affects intermediate inputs. Keep sprite goldens unchanged pending an
  explicit color-behavior decision before porting that path.
- Audio diagnosis: three dialogue files have equal sample counts on 386/amd64;
  only 206/205/52 samples differ respectively, each by at most one int16 unit.
  Diagnostic 386 SSE floating-point flags produce byte-identical amd64 PCM and
  original hashes. Production flags and goldens remain unchanged. This finding
  covers these samples, not all audio. See the tracked detailed report below.
- Final full default 386/CGO suite with assets: 14 passing packages, 3 failing
  packages, 32 without tests, with no compilation/vet failures. Remaining failing
  packages are blobs (only TestReadBlobs), noxrender (sprite references), and ail
  (audio references). Log: build/diagnosis/final-suite.jsonl. Tracked engine files
  remained unchanged by execution. No new game behavior was introduced, so the
  previously verified three builds and gameplay baseline were not rerun for
  these formatter/test-only changes.

The bounded diagnosis is complete. Durable findings, measurements, reproduction
commands and follow-up decisions are in
[docs/porting/FAILURE_DIAGNOSIS.md](docs/porting/FAILURE_DIAGNOSIS.md).
Next: use the saved build graphs for a bounded compiled/linked C inventory,
choose an independent leaf, and establish its C-reference differential tests
before conversion. Do not require the entire baseline suite to be green, and do
not use obsolete blob writers or unresolved render/audio goldens as port oracles.

### First conversion: protection checksum — 2026-09-10

- Bounded build inventory: standard/highres select 152 repository C translation
  units, server 151. Both checksum symbols are retained in all baseline binaries.
  Read-only reproduction tool and dependency scope: docs/porting/C_INVENTORY.md.
- Committed pre-conversion reference tests and inventory as 00228a81. Existing
  Go checksum agrees with the untouched historical C functions on 386.
- Replaced the two production C checksum definitions with Go exports calling a
  shared internal/protection implementation. Retained the original C only behind
  porttest. C ABI width, return bits, null handling, word/tail boundaries,
  unaligned buffers, chunk boundaries and non-mutation checks pass. Differential
  tests pass for default/server/highres; both bounded fuzz runs pass.
- All three production builds succeed. Test reference symbols are absent from
  their binaries. Fresh warrior scenario checksum-port exits 0 against both
  preserved screenshots with overrides disabled. Full suite has 15 passing,
  3 known failing, 32 no-test packages, with no compile/vet failures.
- Remaining C-to-Go calls have measurable overhead; the local benchmark and
  interpretation are recorded in docs/porting/PROTECTION_CHECKSUM.md. Do not claim
  this conversion improves performance or covers every surrounding caller.
- Production .c physical LOC is now 142,637 (−28), 153 files. Test-only reference
  is 33 lines separately. The user requested counts after EVERY conversion chunk;
  tools/porting/c_loc.py and docs/porting/C_LOC.md define and track this measure.

Next session: read docs/porting/PROTECTION_CHECKSUM.md, retain its differential
reference tests, and select the next cohesive leaf with caller/state evidence.
One checksum implementation plus its nullable wrapper has now moved out of C;
no broad protection-manager or render/audio conversion has been attempted.

Conversion commit: 66fa7bd4; pushed to dbenamy/opennox dev with the preceding
reference-test commit 00228a81. Only the original media archive is untracked.

### Retire checksum C test reference — 2026-09-10

At the user's request, removed internal/protectionref after the completed
conversion comparisons. The historical C implementation and differential harness
remain recoverable from 66fa7bd4. Keep the permanent Go fixed-value/property tests
and tagged ABI tests; the latter now calculate expected results independently by
byte lane, preserving alignment, chunk, mutation and nullable-length coverage.
The benchmark retains direct Go and C-to-Go paths only. Production code is unchanged.

Validation: 386 Go unit tests and TestProtectionABI pass; a five-second ABI fuzz
run passes 551,441 cases. Logs: build/port-checksum/retire-{unit,abi,fuzz}.log.
Production C count stays 142,637 lines across 153 files; test-reference C drops
from 33 to 0. C_LOC.md and the handoff reflect this retirement. No engine rebuild
or gameplay rerun was needed for this test-only removal.

### Protection record helpers — baseline checkpoint, 2026-09-10

Selected sub_56F590 (decoded-ID lookup), sub_56F6F0 (index lookup), and
sub_56F720 (payload swap). All are retained in the standard baseline binary.
Their only diagnostic callback, nullsub_31, is an empty C function. Tests execute
2,000 deterministic scenarios against current C, with temporary C-heap records
and restored globals: empty/single/multiple lists, duplicates, high-bit keys/IDs,
missing and extreme indices, null/self/adjacent/non-adjacent swaps, preserved
links and modulo-32-bit counter increments. Current C and new pure Go helpers
pass separately before rewiring the ABI. Logs: build/port-records/c-before.log
and unit.log. Production C count is still 142,637; no extra C reference is needed.

Protection record helper conversion completed: the same ABI scenarios pass on
386 for default/server/highres; pure Go tests pass on 386/amd64, all three targets
build with Go export bridges, and records-port exits 0 against both preserved
screenshots. The global layout assertions pass. No C reference was added.
Production C: 142,570 lines (−67 this chunk), 153 files; test-reference C: 0.
See docs/porting/PROTECTION_RECORDS.md for scope, commands and limitations.
Next chunk: inspect and test the protection spell/ability bitset operations;
keep allocation, rekeying and floating-point state outside that scope.

### Protection bitset operations — baseline checkpoint, 2026-09-10

Current C passes 5,000 deterministic state scenarios plus 12,291 direct bit
checks before conversion. Tests cover signed handle thresholds, empty/missing
records, key zero/high bits, enabled/disabled awards, checksum deltas, successful
and failed validation, ignored entry zero, modulo-32 collisions, signed truthy
values and count <= 1 with null data. Inputs and non-payload state stay unchanged.
Logs: build/port-bitset/c-before.log and unit.log. No C reference is copied.

Protection bitset conversion completed: the same C-before/Go-after ABI checks
pass for default/server/highres; pure Go tests pass on 386/amd64. All three
binaries build with the Go exports, and bitset-port exits 0 against both preserved
screenshots. Existing Go wrappers avoid a C round trip. Production C is 142,503
lines (−67 this chunk), 153 files, with zero C reference lines. Details:
docs/porting/PROTECTION_BITSET.md. Next inspect integer/float record construction,
including exact float bit patterns and allocation-failure state handling.

### Protection constructors — baseline checkpoint, 2026-09-10

Current C constructors pass 16,224 cases (1,014 bit patterns × four keys × four
C/Go integer/float call paths) before replacement. Patterns include signed zero,
subnormal boundaries, infinities and NaN payloads plus seeded random values.
Tests use an empty C-owned manager, verify exact words/checksum/list endpoints,
and restore globals/free records. Empty insertion draws no randomness. The new
pure Go initializer separately checks nil-allocation state and reset links.
Logs: build/port-create/c-before.log and unit.log. No C reference is copied.

Protection construction completed. A generated C-to-Go float export failed the
signaling-NaN test (7f800001 became 7fc00001); no expected values were relaxed.
A caller audit showed no remaining C caller for the float constructor once its
Go wrapper calls the shared initializer directly. Removed that unused C entry
point/declarations rather than retaining a float shim. Integer C entry remains.
All 12,168 live-path constructor cases and earlier ABI tests pass under all three
tags; pure Go tests pass on 386/amd64. All accepted binaries build and have the
expected symbols. The create-port scenario exits 0 against both screenshots.
Full suite: 15 passing, 3 known failing, 32 no-test packages, no compile/vet errors.
Use build/port-create/accepted-* artifacts; earlier outputs are diagnostics.
Production C is 142,458 lines (−45), 153 files; test-reference C is zero.
See docs/porting/PROTECTION_CREATE.md. Next: record deletion and manager cleanup;
only the delete-and-clear operation has remaining C callers, so preserve that
ABI while routing existing Go cleanup directly to Go.

### Protection deletion/cleanup — baseline checkpoint, 2026-09-10

Current C passes 1,000 deterministic removal/cleanup sequences, including exact
surviving payloads/links/endpoints, head/middle/tail and duplicate-ID removal,
misses/repeated deletion, zero and UINT_MAX IDs, checksum updates, uint16 count
wrap and cleanup resets. The handle sequence remains unchanged. Fixtures own
individual C allocations, call the public cleanup wrapper and restore globals.
Pure Go unlink tests pass separately. Final pre-port log:
build/port-remove/c-before-final.log. Only delete-and-clear has remaining C
callers; cleanup and the internal delete-by-ID entry can become direct Go.

Protection deletion/cleanup completed: original-C and Go-after state sequences
pass, all accumulated ABI checks pass for default/server/highres, and pure Go
checks pass on 386/amd64. All three builds have the expected retained/removed
symbols. remove-port exits 0 against both preserved screenshots. Production C
is 142,393 lines (−65 this chunk), 153 files; C references remain zero.
See docs/porting/PROTECTION_REMOVE.md. Next: randomized record insertion, with
explicit comparison of list order and RNG index/consumption under fixed seeds.

### Protection randomized insertion — baseline checkpoint, 2026-09-10

Original C passes 500 deterministic insertion sequences plus prepopulated
32,768/65,535-record boundaries. Tests compare exact list order, back links,
endpoints, checksum, count wrapping, both constructor paths and Logic/Other RNG
indices. Pure Go InsertBefore tests pass separately. Logs:
build/port-insert/c-before-final.log and unit.log. Production C remains 142,393
lines. The sole production caller is the Go constructor; remove the obsolete C
entry point after equivalence validation. Trial delegation: Terra implements
this bounded conversion; the primary agent reviews and runs integration checks.

Randomized insertion completed: primary review accepted Terra's implementation
without corrections. All accumulated protection ABI tests pass under three tags,
pure helpers pass on 386/amd64, all production binaries build with expected
symbols, and insertion-port exits 0 against both screenshots. Evidence:
build/port-insert. Production C: 142,351 lines (−42), 153 files; C references: 0.
See docs/porting/PROTECTION_INSERT.md. Next: reserved-record initialization and
handle allocation, including uint32 sequence wrap and return-value semantics.

### Protection reserved records/handles — baseline checkpoint, 2026-09-10

Original C passes 400 deterministic mixed operation sequences against full list,
checksum, handle sequence, return-value and RNG snapshots. Cases include empty
and prepopulated lists, zero/high-bit IDs, the threshold and uint32 wraparound.
Final pre-port evidence: build/port-handles/c-before.log. Both C entries remain
needed by sub_56F1C0. Reserved initialization increments the sequence after each
attempt; ordinary allocation increments only on success. Allocator exhaustion is
not injected. Production C remains 142,351 lines; reference C remains zero.

Reserved records/handles completed: original-C and Go-after sequences pass,
including all accumulated protection tests for default/server/highres. All three
production builds have the required Go-backed C exports. handles-port exits 0
against both preserved screenshots. Production C: 142,327 lines (−24), 153
files; reference C: 0. Evidence: build/port-handles; see
docs/porting/PROTECTION_HANDLES.md. Next: record rekey/shuffle, testing exact
payload order, unchanged links, checksum resets and RNG/counter consumption.

### Protection rekey/shuffle — baseline checkpoint, 2026-09-10

Original C passes 400 deterministic scenarios in C-export and Go-wrapper modes.
Checks compare exact shuffled decoded values, unchanged node identities/links,
checksum reset, raw key/return, count/sequence, wrapping counters and both server
RNG indices. An independently called unchanged C floating-RNG helper supplies
the expected draw and raw post-state; the fixture restores pre-state before
calling rekey. Range vardefs are saved via their C addresses, not blob offsets.
Pure Rekey tests pass on 386/amd64. Evidence: build/port-rekey/c-before.log and
unit-*.log. C is still 142,327 lines; no reference copy is added.

Rekey/shuffle completed: all accumulated protection tests pass in three target
configurations, pure helpers pass on 386/amd64, and all three binaries build with
expected exports. rekey-port exits 0 against both preserved screenshots. Full
default suite: 15 passing / 3 known failing / 32 no-test packages, with the exact
same failing package/test entries as the accepted constructor checkpoint.
Production C: 142,265 lines (−62), 153 files; reference C: 0. See
docs/porting/PROTECTION_REKEY.md and build/port-rekey. Next retire the now-unused
integer struct-constructor C bridge (remaining production callers are native Go),
then continue protected-value validation/mutation.

### Unused protection bridge cleanup — completed, 2026-09-10

Retired integer struct-constructor and single-bit C exports/prototypes after
caller audits found only native Go production paths. Go APIs and behavior tests
remain; constructor cases now number 8,112 across its two live paths. All
accumulated protection tests pass under three tags, all binaries build with
expected symbols, and bridges-port exits 0 against both screenshots. Production
C remains 142,265 lines (delta 0), 153 files; references 0. See
docs/porting/PROTECTION_BRIDGES.md and build/port-bridges. Next: integer/byte/word
protected-value setters; their decompiled pointer returns are raw scalar bits.
Caller audit found no dereferences, so uint32 C return declarations are suitable.

### Protection integer/byte/word setters — baseline checkpoint, 2026-09-10

Original C passes 2,500 setter calls (500 scenarios × four C entries and one
Go-wrapper path), plus the rekey regression checks after shared fixture reuse.
Tests cover signed eligibility, misses/duplicates, truncation/return bits and full
manager/RNG post-state, with unchanged-state assertions on unsuccessful paths.
Baseline: build/port-setters/c-before.log. Caller audit confirms pointer-typed
returns are scalar bits, with no dereferences/function-pointer uses, so correct
these four declarations to uint32_t while preserving 386 return behavior.
Production C remains 142,265 lines. Terra's bounded Go draft and primary-owned
tests both passed review; no production replacement has been made yet.

Integer/byte/word setters completed: all accumulated protection checks pass for
three configurations, all production targets build with the four Go-backed C
exports, and setters-port exits 0 against both preserved screenshots. C return
declarations now reflect scalar bits; call-site uses remain compatible.
Production C: 142,189 lines (−76), 153 files; reference C: 0. See
docs/porting/PROTECTION_SET.md and build/port-setters. Next: additive protection
updates (int32, signed int16 mana, unsigned uint8 level), with wraparound tests.

### Protection additive updates — baseline checkpoint, 2026-09-10

Original C passes 2,000 calls across integer, signed-short mana, unsigned-byte
level and Go mana-wrapper paths. A widened signed oracle checks modulo addition;
full manager/RNG snapshots cover boundary/random values and failed lookups.
Baseline evidence: build/port-add/c-before.log. C is still 142,189 lines.
The missing F980 declaration is added to its existing header for test access.
After conversion, scalar return declarations require five explicit C casts at
existing pointer-typed surrounding results; these do not change caller APIs.

Additive updates completed: original-C and Go-after arithmetic/state expectations
pass along with accumulated protection tests in all three configurations. All
production builds have the expected three Go-backed C entries; add-port exits 0
against both screenshots. Production C: 142,130 lines (−59), 153 files; reference
C: 0. See docs/porting/PROTECTION_ADD.md and build/port-add. Next: buffer checksum
validation (sub_56FB00), including signed eligibility, first-match lookup,
partial-word handling and proof that rejected/missing IDs never read the buffer.

### Protection buffer validation — baseline checkpoint, 2026-09-10

Original C passes 515 validation cases: aligned/unaligned buffers and trailing
bytes, first-match duplicates, signed eligibility, missing IDs, nil huge lengths,
and PROT_NONE guard pages proving short-circuit/no-partial-word reads. Full
manager/RNG state and readable input bytes stay unchanged. Final baseline:
build/port-validate/c-before-final.log (zero-key cases include live matches).
Production C remains 142,130 lines; draft and tests were reviewed before port.

Buffer validation completed: all protection tests and production builds pass
in three configurations; the C entry is Go-backed and validate-port exits 0
against both screenshots. Production C: 142,115 lines (−15), 153 files; reference
C: 0. See docs/porting/PROTECTION_VALIDATE.md and build/port-validate. Next object
checksum/toggles require object/type fixtures and preserve raw-ID returns on
missing records, unlike the setter functions.

### Protection object checksum/toggles — baseline checkpoint, 2026-09-10

Final original C passes 520 scenarios through both toggle entries (1,040 runs),
with direct checksum comparisons and 1–3 repeated toggles per run. Fixtures use
C-allocated objects/data and a porttest-only temporary type table. They compare
object/health/init/name bytes, manager/list/RNG state, raw-ID returns on misses,
and XOR restoration. Guards prove no object reads on rejected/missing IDs and
no init-data reads for absent types, nonpositive signed sizes or partial words.
Type index zero with positive init data has an explicit case. Evidence:
build/port-object/c-before-final.log. Production C remains 142,115 lines.
Retire CRC/getter bridges with no remaining C callers; the length-aware checksum
entry still has a C caller. FC50's const parameter becomes non-const in the
internal declaration to match the generated Go export; behavior remains read-only.

### Object checksum/toggles completed — 2026-09-10

Ported the digest and both toggles, removed sole-use object getter and checksum
bridges, and preserved missing-ID return behavior and read-only object access.
Original-C baseline is `e4127e22`; 1,040 toggle scenarios plus direct digest and
guard-page checks pass after the port. See [details](docs/porting/PROTECTION_OBJECT.md).
Production C: **141,984 physical lines (−131)**; test-reference C: **0**.

### Float updates completed — 2026-09-10

Ported both float updates after 3,600 original-C ABI calls and independent
precision-53 arbitrary-precision tests. Actual hosted x87 precision corrected
the standalone C probe assumption before the port. See [details](docs/porting/PROTECTION_FLOAT.md).
Production C: **141,941 physical lines (−43)**; test-reference C: **0**.

### Initialization completed — 2026-09-10

Ported startup after 400 original-C/wrapper baseline runs with the shipped
floating constants; retired the sole-use C bridge. See [details](docs/porting/PROTECTION_STARTUP.md).
Production C: **141,914 physical lines (−27)**; test-reference C: **0**.

### Floating RNG/state completed — 2026-09-10

Ported the remaining four helpers and moved their private state to Go after
46,046 original-C snapshots and independent arbitrary-precision tests. Retired
all four C bridges and both C range globals. See [details](docs/porting/PROTECTION_RANDOM.md).
Production C: **141,844 physical lines (−70)**; test-reference C: **0**.

### Client unit-code/bit helpers completed — 2026-09-10

Ported three C entries after exhaustive low-word/upper-pattern bit tests and
2,988 read-only drawable cases. See [details](docs/porting/NETWORK_CODE.md).
Production C: **141,821 physical lines (−23)**; test-reference C: **0**.

### Dynamic unit-code/extent lookup completed — 2026-09-10

Ported two C entries through the existing typed server lookup after exhaustive
unmarked-code bypass and randomized C-backed object-list baselines. See
[details](docs/porting/NETWORK_EXTENT.md).
Production C: **141,786 physical lines (−35)**; test-reference C: **0**.

### Waypoint helpers completed — 2026-09-10

Ported allocation, duplicate next-link entries and composite flag predicate,
retiring its sole-use mask bridge after 459,008 predicate cases and byte-level
link/allocation baselines. See [details](docs/porting/WAYPOINT_HELPERS.md).
Production C: **141,745 physical lines (−41)**; test-reference C: **0**.

### Completed — map-rule loading/parsing (2026-09-10)

Ported 57A1B0/57A1E0/57A3F0/57A4D0/57A620 after original-C baselines at
2bf05750 and 9c86046c. Go owns context, file reading, tokenization and directive
application; only header lookup and top-level loader retain live C bridges.
Independent file/selection/settings/encoding/list tests pass in all variants.
All binaries build, rules-port gameplay passes, and the full-suite failure
multiset exactly matches the prior milestone. Production C: **141,455 (−290)**,
153 files, zero test-reference C. See docs/porting/RULE_LOADING.md.

### Writer baseline — historical decision checkpoint (2026-09-10)

Confirmed original online writer memory-layout/output instability; saved the
asset-free diagnostic and compiled-offset analysis in docs/porting/RULE_WRITER.md.
Added 81 stable offline baseline cases. All accumulated ABI tests pass on 386
default/server/highres with those cases. No writer production changes; C remains
**141,455** physical lines. Pending user choice: fix the online bug as part of the
writer port, or postpone that chunk and continue elsewhere.

### Completed — rule writer and approved online fix (2026-09-10)

Ported 57AAA0, replacing overlapping decompiler-split temporaries with two full
Settings2 values. The user explicitly approved fixing the unstable online output.
All 665 writer cases and accumulated ABI tests pass in all variants; three builds,
writer-port gameplay and exact known-failure full-suite comparison pass. C is
**141,351 (−104)** physical lines, 153 files, zero test-reference C. See
RULE_WRITER.md for independent online expectations and original offline baseline.

### Completed — rule-file deletion (2026-09-10)

Ported 57A9F0 after original-C baseline b861ab46. Ten full-tree/exact-return cases
and all accumulated ABI tests pass on 386 default/server/highres; all binaries
build and rule-remove-port passes both gameplay screenshots. Production C:
**141,340 (−11)** physical lines, 153 files, zero test-reference C. Full-suite
milestone remains the immediately preceding writer chunk. See RULE_REMOVAL.md.

### Completed — command-rule loading/dispatch (2026-09-10)

Ported 57A950/4D0550/4D0670/57AE30 after original-C baseline b803c933. Preserved
any-bit mode checks, exact headers, byte widening/254-byte physical-line copies,
callback effects, literal path quirks and file precedence. Added guards only for
undefined short-path/read-error cases. All targeted variants, all binaries and
rule-command-port gameplay pass; full-suite failure multiset remains exactly
unchanged. Production C: **141,215 (−125)** physical lines, 153 files, zero
reference C. See COMMAND_RULES.md.

### Spell-class eligibility completed — 2026-09-10

Ported 57AEA0 with original-C baseline 2970e5e9, preserving full-width class input,
real spell flag lookup and exact 0/9 returns. Removed unused C chat predicate.
All three accumulated test variants, builds and fresh gameplay checks pass.
Production C: **141,180 lines (−35)**; details in docs/porting/SPELL_CLASS.md.

### Player-ping aggregates completed — 2026-09-10

Ported 554290/554300 with original-C baseline fea6ca7b. Preserved active-player
order, host exclusion, two timing reads per qualifying player, unsigned minimum,
32-bit wrapped sum and signed average division. Retired both unused C bridges.
All three accumulated test variants, builds and fresh gameplay checks pass.
Production C: **141,126 lines (−54)**; see docs/porting/PING_AGGREGATES.md.
Next alias-table work has a pending user decision documented in the top checkpoint.

### Network aliases and exhaustion fix completed — 2026-09-11

Ported reset/select/write with original-C helper baseline 19d02832. Actual-caller
regressions reproduced the approved bug and now pass with both sentinel checks
fixed. Exact packet/sprite/camera behavior continues on exhausted tables. All
three accumulated test variants/builds and fresh gameplay pass; full-suite
failure multiset exactly matches baseline (1,553 entries). Production C:
**141,082 lines (−44)**. Details: docs/porting/NETWORK_ALIASES.md.

### Glyph/item eligibility completed — 2026-09-11

Ported both predicates and their caches with original-C baseline d6d7c136.
Preserved lookup-before-gates, glyph restriction before cheat, callback ordering,
and observed 386 class-shift behavior. Retired unused item C bridge; shared C
cheat flag remains live. All three accumulated test variants/builds and gameplay
pass. Production C: **141,042 lines (−40)**. Details: docs/porting/GLYPH_ELIGIBILITY.md.

### Collision primitives completed — 2026-09-11

Ported reflection/containment with original-C baseline f85e37ee, exact raw-bit
reflection and compact containment result fixture. Preserved PC53 arithmetic,
strict boundaries, NaN quieting, overlapping inputs and live C ABI. All three
accumulated test variants/builds, fresh gameplay and full-suite comparison pass;
same 1,553 known failure entries. Production C: **141,000 lines (−42)**.
Details: docs/porting/COLLISION_PRIMITIVES.md.

### Completed — player controls, respawning, observers and bot transitions

Baseline `1a87b410` was committed/pushed before conversion. All 56 functions /
1,971 C lines are native. All 3,205 cases / 58 full captures match C exactly
(10.976s). C is **115,985 lines / 149 files / zero reference C**. See
[PLAYER_CONTROLS.md](docs/porting/PLAYER_CONTROLS.md) for both initialization fixes,
the XP-protection encoding finding, caller connections and evidence.

Accumulated default/server/highres tests, including 44,559 focused cases / 422
groups, pass: 192.132s / 179.762s / 196.813s. Three production binaries verified
ELF32/i386/SSE2/CGO. Full suite with assets has the exact known failures: 1,553
entries; 15 packages pass, 3 fail, 32 skip. Fresh unchanged repeat-a gameplay
passes in 42.972s. Evidence: build/port-player-controls and baseline/runs/player-controls-port.

At that checkpoint, the next action was to retire 30 exports unused by production callers and route their fixture
operations directly to Go. Draft build/port-controls-bridges/retire.py is not yet
applied; review before use. Keep all controls hashes unchanged, verify all targets
and symbols, document unchanged C LOC, commit/push and continue the next batch.
The next algorithm scope was staged in build/port-spell-lifecycle: 29 spell-casting/buff
blocks, 1,328 address-block lines before declaration audit. Do not rerun stale
controls source generators. No pending question; no new agents; preserve archive.

### Completed — player-controls obsolete exports

Controls conversion `0a3f446d` is pushed. Follow-up removes 30 unused exports and
header declarations, with direct Go fixture dispatch and unchanged expected
hashes. All 3,205 cases / 58 groups pass in default/server/highres (9.281s /
9.427s / 14.374s); three ELF32/i386/SSE2/CGO builds succeed and all 30 symbols are
absent. Fresh unchanged gameplay passes in 34.974s. C remains 115,985 / 149 files /
zero reference C. Evidence: build/port-controls-bridges. Committed/pushed 9fd4ebbf.

### Spell casting and buff lifecycle — historical C baseline

At this baseline: 29 functions / 1,326 removable C lines; production was still C. Optional guarded
fixtures and 2,246 cases / 60 expected hashes are applied. First full C run passes
in 8.743s; final locked spell + controls corpus passes twice in 38.765s. Phoneme typed-
pointer offset bug corrected before locking (DECISIONS.md). See
[batch notes](docs/porting/SPELL_LIFECYCLE.md). The then-next tasks were baseline commit/push, conversion, the EnchantPower
duration-array correction and export qualification; see the linked batch report
for their completion. Stable controls/reward
captures are now losslessly compressed .json.gz; tracked hashes unchanged.

<!-- next-scope-draft -->
Read-only next-scope candidate: build/port-server-panels/candidate.json and
plan-draft.md, 58 functions / 2,120 body lines from remaining servopts panels and
related GAME2/GAME3/GAME3_1 blocks. No next-batch source changes yet. External
reference graph reaches every candidate; review literal branches/registrations.

<!-- next-collision-core-draft -->
Read-only next candidate: build/port-collision-core/proposal.json,
combined-reachability.json and plan-draft.md. Twenty-three live functions /
993 C body lines in collision dispatch/queues/activation and remaining contact
geometry. No next-batch source or fixture changes installed.

Disk cleanup reclaimed **6.184 GiB** from twelve completed scenario copies after
verifying every removed file against the original asset hash. Per-run restoration
manifests preserve how to recreate them. Original assets/archive, changed files
and all reports remain. About **20 GiB** is free. The spatial-targeting
**deduplicate-completed-assets.py --apply is consumed; never repeat deletion mode.**
See build/port-spatial-targeting/completed-assets-{plan,audit}.json and applied log.

<!-- item-respawn-disk -->
Session-entry native asset deduplication is complete:1,660,044,319 bytes reclaimed.
Audit29716/apply32358 are joined; the script under build/port-item-respawn is
consumed. Original assets/archive are unchanged.

Latest disk cleanup:72 verified superseded compressed binaries removed,
1,648,354,499bytes reclaimed. Audit/apply in build/port-player-death/cleanup-compressed.py
are consumed; preserve superseded-compressed-binaries.json and all reports.

Chat-bubble disk cleanup: 49 completed historical test logs were losslessly
compressed after verifying decompression against originals, reclaiming
2,239,764,149 bytes. Manifests: build/port-chat-bubbles/compressed-completed-logs.json
and compressed-historical-logs.json. Restore individual original paths with gzip -dk.
Both compression scripts are consumed; captures/reports/assets/archive remain.
About 2.9 GiB was free after cleanup, before C production qualification.

Read-only next candidate: build/port-client-spell-presentation/selection-draft.json
and plan-draft.md. Twenty-two client spell/item presentation bodies; caller audit
and boundary review remain before accepting the next batch. No source installed.


Read-only next candidate: build/port-client-audio-events/selection-draft.json,
56 bodies /1,031 body lines across audio event scheduling/playback and dialog
queue helpers in GAME2.c, GAME1_3.c and client__audio__audevent.c. Whole-repository
reachability and ownership audit remain; no next-batch source is installed.

Prepared but NOT RUN: deduplicate-client-audio-streams-completed-assets.py can
audit completed C/native scenario copies only after production finishes. Its
apply mode must not run until the independent hash/inactivity audit is reviewed.


## Active continuation — client audio events C baseline

Qualified parent **50f0711d** is pushed. Tracked scope/audit: 62 bodies / 1,182 body
lines, 35 external roots, all reachable. The C dispatcher and actual global owners
are installed. Format/volume/pan/switch and music/default contracts passed
(primitives-initial and music-defaults-initial); both sessions are joined.
Refill-initial is joined: ordinary refill and scratch guards passed; the fixture
was corrected to assert the existing empty-chunk termination convention.

selection-refill-second is active (session54366); no source edits until joined.
It adds actual-RNG sample selection and reruns corrected refill boundaries. The
C refill bodies use a test-only external-device observer, with normal AIL fallback
when no fixture owner exists. The conditional adapter block changes the C file;
production algorithms remain unchanged. Reaudit baseline identity or rerun C
production qualification. No goldens frozen. Remaining: connected event ownership,
scheduling/spatial contracts, fixture review/repeats, baseline qualification/commit.
Installed build-dispatcher.py and music/selection drafts are CONSUMED.


Read-only next candidate: build/port-character-creation/selection-draft.json and
plan-draft.md,36 bodies /1,259 body lines across GAME3 class/color/config helpers
and client shell selclass/selcolor. Caller/callback/global review remains. No source
installed for this candidate. Mapped appearance parser callbacks are live roots.

Character-creation candidate refined by read-only review:23 UI bodies /1,042 body
lines (ui-selection-draft.json);13 adjacent configuration bodies /217 lines are
separate and deferred. Name helper/test drafts exist only under build, NOT installed
or run. Original36-body selection remains as provenance. No next-batch source edits.

<!-- next-candidate -->
Historical initial server-browser candidate:21 bodies /1,174 lines. Superseded
by the current65-body baseline and tracked plan; do not resume the initial scope.
Review the apparently unconsumed old configuration callback table and disabled
online branches before deciding reachability; coordinate getter has live callers.
