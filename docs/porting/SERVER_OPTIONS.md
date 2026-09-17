# Server-options UI in Go

Following team UI, port 33 live functions / 1,256 corrected C body lines from guiserv.c
and GAME2.c 457460..459DA0. Remove the separate 28-line sub_457FE0 orphan:
its only two call sites are inside literal `if (0)` blocks; no registrations,
Go wrappers or other callers exist. Raw regex reachability overcounts this edge.
Do not include the unrelated drawable helpers beginning at 459DB0.

Owner boundaries: parsed GUI resources and actual widgets/rendering via entryOwner;
real settings records, map list and rule files; existing team/runtime owner for
team actions; real subpanel implementations where reached. Subpanel access,
spell/object/general/advanced windows remain outside this conversion but their
construction/destruction edges require coverage. Existing event adapters may
observe map-transition and dialog boundaries without copying implementation.

Contracts before translation:
- Game mode lookup exhausts 16-bit inputs, mask 0x17F0, seven table modes, Arena
  fallback, hidden Quest entry, lazy localized title loading and reverse lookup.
- Actual dropdown population/order/dimensions under variable fonts, selection,
  capture, mouse position, raw geometry and pixel drawing.
- Settings read/write: 58-byte record, 15-byte server names, signed conversion,
  empty/negative/255/256/65535/65536 numeric edits, map-tab splitting, bit retention,
  client/host and CTF/Highlander/Quest control enablement. Check full record bytes.
- Map lists: real enabled/disabled entries, compatible/incompatible modes,
  empty/no-match/case-matched current map, ordering, recommended counts, selection,
  scrolling and modified-rule palette; actual distinct default/user rule values.
- Constructor: all nine locales and font-height override; repeated open/close,
  failed resource/retry, host/client, child parentage, ownership and clean teardown.
  Constructor currently dereferences failed resource: reproduce and correct if
  confirmed, then qualify/freeze corrected C baseline separately.
- Tabs and nested panels: repeated switches across player/access/general, focus,
  checked state and cleanup; preserve valid pointers until owned cleanup finishes.
- Event matrix including irrelevant/unknown events and child IDs, commit/cancel,
  team thresholds, confirmation callbacks, dirty flag and reliable queue effects.
- Apply settings: empty/same/different map, map cycle, rule writer/reader effects,
  team auto-creation, game flags, server name, limits and hide behavior. Normalize
  only incidental pointer return identities; retain semantic returns and effects.
- Team count labels, tooltip flag bytes, dirty/visible/open accessors and actual
  external callers. Do not substitute zero-filled constant tables for shipped data.

State audit: named state appears private to this scope, except root accessor in
legacy/player.go (move caller with owner). Check numeric mapped-memory references
before retiring extracted globals. shared-state-audit.json is only a candidate
list, not sufficient proof. Preserve settings-owned blobs and list roots shared
with settings code. Native mode table replaces C table/cache only after baseline.

C baseline can reuse the qualified team-UI production result only if production
source is identical. Any prerequisite fix requires a fresh C production gate.
Freeze independent contract and original-C captures, commit/push; then native
translation with unchanged captures, three-target broader sweep and fresh
production/gameplay/save-load/flat gate. Record exact C LOC and review decisions.

Parent production baseline: **3d6436bb**, fully qualified/pushed team UI.
First C mode contracts are installed; mode fixture drafts are now consumed.
The two prerequisites below are the only production changes so far. C is
50,181 / 80 files / zero reference C. The completed baseline qualification is recorded below.

## Baseline progress and prerequisite

The exhaustive 16-bit mode contract passes against original C. The first fixture
run also passes limit formatting, server-name truncation and dirty/visibility
contracts. Dropdown/read-settings failures were fixture mistakes: Count is the
list allocation capacity (Field_11_0 is the live row count), and shipped control
10119 is a PUSHBUTTON whose draw-data text holds the selected mode, not a static
text widget. Corrected those fixture types/observations and arranged nested
controls under their actual owning panels.

An isolated missing-resource contract reproduced SIGSEGV at address 0x8 in
nox_xxx_guiServerOptsLoad_457500: the constructor reads width immediately after
the parser returns NULL. Added an immediate return 0 before layout/child access.
This reversible prerequisite adds three C lines: **50,180 / 80 files / zero
reference C** while baseline work continues. Because production source changes,
the corrected C baseline requires a fresh production qualification; parent
production evidence cannot be reused. Reproduction: focused-3.log. Do not freeze
the initial captures; the fixture matrix is incomplete.

The tab lifecycle contract also reproduced a stale general-panel pointer in the
parent options owner: closing destroys the actual general panel but leaves
1046540 nonzero. Added the missing clear beside the existing panel cleanup.
This is an ownership-state correction; no new membership or settings behavior is
introduced. Reproduction: focused-12.log. The two prerequisites now add **four C
lines**, leaving **50,181 / 80 files / zero reference C** during baseline work.
The selected live bodies become 1,256 lines after these corrections.

Apply contracts use the actual settings, rule writer/reader and close algorithms;
only the existing SwitchMap API is observed, so per-case map requests can be checked
without booting another game. Fresh production scenarios still exercise real game
loading. Cancel tests install the actual reliable sender and queue owner: the chat
cancel branch sends the existing team-clear message. Checkbox tests also use real
key events; parent notification happens before the checkbox toggles, so direct
handler inputs represent the previous checked state.

## Event and dependency coverage

Advanced host/client modals run their actual constructors and cleanup. Repeated
top-level and general/access/player tab switches check parent/child links and
panel lifetime. Team controls use the real team owner and reliable sender;
confirmation thresholds run zero through four teams in CTF, Flagball, Arena and
KoTR, including invoking the actual confirmation callback and checking that team
creation occurs only after acceptance. Mode selection checks current-selection
validity independently of event-row payload, preserving that original distinction.
Map selection loads actual distinct base/user rules; reset removes the complete
user rule file and reloads the base rules while preserving row text and selection.

The shipped mode-to-limit index table is explicitly installed from blobdata, with
distinct nonzero score/time values per mode. Byte-wrapping team count labels and
group-count labels are independently asserted. Startup contracts check first-open
consumption and the existing quest/server console-command boundary, including
reopen without executing the command again. Static mapped-memory preflight passes.
These original-C contracts supplied the frozen expectations below.

The final descriptive-label contract caught a fixture omission: GameTypeIs reads
its localization filename from blob 131072, which was empty in the fixture.
Supplying guiserv.c restores the real lookup. All 23 frozen captures remain
unchanged after this correction; final target sweeps include the label assertions.
This is test-only and requires no production rebuild: the fresh corrected-C
production source is identical, with only two porttest files changing afterward.

## Qualified C baseline

All 35 server-options roots pass (847 leaf cases, plus exhaustive 65,536-mode
lookup within one case), with **23 frozen captures / 1,322 records**. Affected
rules, entry/list widgets, teams, HUD and roster dependencies bring each target to
**146 roots / 12,107 leaves** and **82 identical captures / 12,982 records**, no
skips. Final default/server/highres sweeps pass in **37.85 / 65.10 / 81.80s**.

Fresh corrected-C production passes in **401.66s**: all three production builds
and ABI checks, exact known asset failures (1,553 entries; 15 pass / 3 fail /
32 no-test), gameplay, save/load and forced flat-map expansion. The latter removes
51 maps and regenerates the selected map exactly. Client SHA-256:
`22b3ffb0711a237534b145cda0fb474546ee57695146d5a5002ce4ba26240977`.

Final sweeps share an unchanged 2,119-file source manifest. The production gate's
manifest differs only in the two porttest files for the late label contract and
filename fixture correction; production source identity is verified explicitly.
All sessions are joined. Evidence: build/port-server-options/c-{default,server,
highres}-final, c-production, c-audit.json and static-final.log. Reproduce through
server-options-c-batch.json and server-options-focused-tests.txt.

**C: 50,181 physical lines / 80 files / zero reference C**. Commit this corrected
baseline before translating the selected 33 live functions / 1,256 body lines and
removing the separate 28-line orphan. Frozen expectations stay unchanged.

## Native conversion

Baseline **55ccc2b7** was committed/pushed before translation. All 33 functions
are installed as Go; 16 private C globals and the mode table/cache are native.
The 28-line orphan is removed. Twelve C interfaces remain: ten actual external
callers and two tooltip callbacks stored through the GUI's C-pointer ABI. Twenty-two
function interfaces are retired; Go callers and fixtures invoke Go directly.
Current C is **48,740 / 79 files / zero reference C (−1,441)**, qualified below.

First compile succeeded. Frozen captures caught two translation differences:
window.PointIn uses different geometry from the legacy absolute-position/live-size
hit test, and clamping a negative edit must not rewrite the displayed text. Both
were corrected without changing expectations. Manual review also preserved the
selected-name alias across record copying and the edited-name text color.

## Native qualification

All **35 options roots / 847 leaves** pass against the unchanged C contracts and
23 frozen captures in 4.852s. The broader dependency selection now includes
rules and actual entry/list widgets as well as the prior UI, gameplay/report,
team and roster corpus. Default/server/highres pass **487 / 486 / 487 roots** and
**57,405 / 57,404 / 57,405 leaves**, with no skips, in **204.41 / 329.54 / 230.07s**.
The single server exclusion is the existing client-only occlusion root.

All **193 captures / 61,102 records** are identical across targets, including
all 82 corrected-C baseline captures and all 170 preceding team-UI captures.
Fresh production passes in **416.35s**: all three builds/ABI and interface checks,
exact 1,553 known failure entries (15 pass / 3 fail / 32 no-test), gameplay,
save/load and forced flat-map regeneration. Client SHA-256:
`f1641ca43485c15532894867eedcaa2bf7aab0f7d061d0d2a976b0e5c35fb058`.

All four final gates share the same unchanged **2,128-file source manifest**;
all sessions are joined. Static memory preflight and retired-interface audits pass.
Artifacts: build/port-server-options/native-{default,server,highres,production},
native-audit.json and interface-audit.json. Reproduction manifest:
server-options-batch.json; selection: server-options-tests.txt.

Final C is **48,740 physical lines / 79 files / zero reference C**, a reduction of
**1,441 lines** including the C mode table, private globals, orphan, declarations
and separators. The original C implementation remains recoverable at 55ccc2b7;
no C algorithm is retained solely for testing. Next candidate: related server
access/general/advanced/spell/object panels, with the current fixtures reusable.
