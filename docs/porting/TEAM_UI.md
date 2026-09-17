# Team HUD and server player-list UI

## Scope

Start from qualified team-runtime conversion `6963545e` (51,203 physical C lines).
The corrected baseline has **51,225 physical C lines / 82 files / zero reference
C**. Five prerequisite fixes add 22 lines. The conversion scope is **35 live
functions / 943 C body lines** (originally 921). All are reachable from external
callers plus the selected callback/helper graph; there is no orphan in this batch.

Include CTF HUD construction, layout, visibility, selection, tooltips and drawing;
Flagball HUD construction/layout/drawing; server player/team list construction,
formatting, selection, mutations and event handlers; and the two deferred map-entry
wrappers. Exclude neighboring ban-list/server-options helpers and the adjacent
Flagball status setter `sub_456140`. That C setter still reads the Flagball root
pointer: retain shared storage for that pointer until its last reader is ported.

Exact corrected bodies: build/port-team-ui/corrected-scope.json. Initial references
and reachability: scope-audit.json / reachability.json in the same directory.
No native UI implementation is installed in this baseline.

## Contracts and owners

Use actual GUI/window parsing, listboxes, renderer, allocated players/teams,
object allocator, intrusive row lists, reliable messages and host netlists.
Controlled images, fonts, resources, strings and nonzero palettes supply inputs.
The dialog observer records the existing API boundary; selected-name extraction,
validation, rename, team storage and messages remain real.

The 25 UI roots cover:

- Missing resources/retry, repeated construction, all nine language indices and
  cap-height boundaries, close with/without window destruction and reopen.
- 224 CTF tooltip combinations, raw state bytes, selection clearing, layout
  thresholds, screen bounds, visibility and actual CTF/Flagball pixels.
- Player add/find/remove, repeated roster/team refresh, real team removal in both
  orders, backing/display row agreement, full names, ASCII case comparison,
  Unicode round trips and fixed-name boundaries.
- 108 selection-control combinations, 25 flag-label formatting cases, all 256
  color-byte values, actual palette colors, draw-based button eligibility,
  headless filtering and quest-mode disabled lists.
- Client join/change requests with observer/missing/same-team cases, host single
  and multiple-player assignment, observer exclusion, selection clearing, rename
  dialog arguments and acceptance, and duplicate-name rejection.
- Map-wrapper zero/positive flag counts (0/1/2/8/16), HUD side effects and return
  values. Without a ball-start object, Flagball opens its HUD but returns zero.

Fixture review matters: locale layout reads the cached default font's cap height,
not merely a named face's height. Client requests use the reliable queue; hosts
use the netlist. Multiselect event 16403 replaces selection; 16405 adds/toggles it.
Own extracted globals separately from backing-blob regions. Static mapped-memory
validation passes. Freeze only complete successful runs.

## Prerequisite fixes for review

Independent original-C contracts reproduced these before freezing:

1. Refresh cleared widgets but accumulated backing metadata. Free both old lists
   before rebuilding, so repeated refresh stays consistent.
2. Name lookup split bare metadata at spaces. Compare the complete stored name.
3. Rename targeted the selected row and left metadata stale. Find by team ID,
   update the backing name, and preserve selection. With no selection, the old
   helper inserted an extra row; with another selected team, it renamed that row.
4. An unsigned comparison accepted the -1 selection sentinel. Compare as signed
   so Assign is disabled when no player is selected.
5. Failed window loading continued into child access. Return zero immediately.

These are authorized reversible correctness prerequisites, not changes to message
formats or membership rules. Preserve existing empty-name acceptance and the C
runtime's ASCII case folding. Broader Unicode case handling is a separate review
item; no source setlocale call changes this runtime behavior. See DECISIONS.md.
Reproductions: players-c, players-c-2, players-draw-construction and lifecycle-c
under build/port-team-ui. The fully reviewed UI run c-reviewed-2 passes in 34.63s.

## Corrected C qualification

Manifest: [team-ui-c-batch.json](team-ui-c-batch.json); selection:
[team-ui-focused-tests.txt](team-ui-focused-tests.txt). **19 UI captures are frozen**
alongside 40 unchanged preceding captures. Include the preceding team/runtime,
roster and objective-dispatch contracts in the affected sweep.

| Target | Roots | Leaf cases | Capture groups / records | Seconds |
| --- | ---: | ---: | ---: | ---: |
| Default | 74 | 11,136 | 59 / 11,660 | 50.47 |
| Server | 74 | 11,136 | 59 / 11,660 | 136.09 |
| High resolution | 74 | 11,136 | 59 / 11,660 | 61.61 |

All pass without skips. Every capture is identical across the three processes;
all share the same unchanged 2,093-file source manifest. Audit: c-audit.json.
Fresh production qualification passes in **375.08s**: three builds/ABI, exact
1,553 known asset-suite failure entries (15 passing / 3 failing / 32 no-test
packages), gameplay, save/load and flat rendering. The flat scenario removes
51 loose maps and regenerates one compressed map. All four gates share the same
unchanged source manifest. Client SHA-256:
`52e539f531d51cb7dccb0ff54ed9a599a8987a6847063c420f3439f498e9ff48`.
Production gate: c-production. All sessions are joined.

Commit/push this qualified corrected baseline, then translate without changing
frozen expectations and qualify the conversion. Never edit sources while a build
or test reads them. Ignored fixture/freeze scripts are consumed; do not rerun them.
The ignored native-hud-draft.go and translation-review.md are uninstalled drafts.

## Disk maintenance

Verified prior team scenario deduplication reclaimed 3.32 GB. With builds joined,
pruning only generated Go cache data unused for over eight hours reclaimed
23,333,730,778 bytes (11,011 files), leaving about 27 GB free at that point. The
local Go 1.26 cache refreshes last-use mtimes within one hour. Current dependencies,
module cache, source, qualification records, original assets and archive remain.
Prune manifest: build/port-team-ui/pruned-go-cache.json; it and prior dedup manifests
are consumed. Older removed cache entries can be rebuilt.
