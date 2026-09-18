# Floor and wall map sections

## Scope and baseline

Native statistics checkpoint **6035e191** is committed and pushed. The next batch
selects **18 connected C bodies /1,629 body lines**: FloorMap, WallMap,
WindowWalls, DestructableWalls and SecretWalls, their historical format readers,
private wall-iteration callbacks and counter resets. The five roots are registered
in src/maps.go and invoked through legacy/maps.go. All are live; historical
readers remain reachable through format-version dispatch. No conversion installed.
See [map-sections-selection.json](map-sections-selection.json).

Production C is **28,794 physical lines /67 files /zero reference C** after the
small prerequisites (+4 lines); no selected body has yet been converted.
Baseline fixtures are being prepared. Reuse the existing map-painting extension
owner for tile grids, subtile pools, normal and generation wall ownership,
normalized pointers and full state snapshots. Exercise actual cryptfile IO with
plaintext test files. Keep independent small wire-format contracts alongside C
captures; report file consumption, return values, mutations and serialization.

## Planned contracts

- Tile records: scalar narrowing, signed variations, subtile order, zero/nonzero
  chains, pool-block boundaries and byte-sized count wrap. Writing also normalizes
  the main tile's scalar fields; observe that mutation explicitly.
- Floor sections: versions1–4 and rejection, whole-grid/region paths, odd/even
  halves, subtile chains, relocated bounds, generation scratch ownership and
  exact roundtrip bytes/state. Audit older branches before choosing inputs.
- Wall sections: versions1–7 and rejection, coordinate/direction/material/
  variation fields, region relocation, skipped/missing walls and ordinary versus
  generation owners. Exercise the historical sprite-dependent branch with actual
  definitions rather than zero-filled stand-ins.
- Window, breakable and secret metadata: flags, counters, ordering, versioned
  fields, region filtering, missing entries and secret-state defaults. Source
  review found paths that retain a previous wall pointer on a missing scratch
  lookup; establish their behavior with independent contracts before freezing.
- After baseline and conversion, run affected map/geometry/painting/layout/save
  checks in all three targets. Fresh native production must pass builds/ABI,
  exact known full-suite comparison and gameplay/save/flat-map regeneration.
  Reuse prior production for a test-only C baseline only with explicit production
  source identity. Any prerequisite C correction requires fresh qualification.

## Disk and recovery

Completed statistics native scenario asset copies were compared byte-for-byte
with originals and deduplicated, reclaiming **1,660,044,319 bytes**. Restoration
manifests, changed maps/saves, screenshots and logs remain. The original assets
and archive are unchanged. Audit/apply modes of
`build/port-game-statistics/deduplicate-game-statistics-native-assets.py` are
consumed; restore mode remains available.

## Initial fixture work

The shared adapter invokes original C section functions through real cryptfile
IO and the existing painting ownership/snapshot fixture. The first324 tile cases
cover read/write, versions1/3/4, nine chain lengths through256 and six scalar
patterns. Independent byte, file-position, return, chain-length and mutation
contracts accompany the pending capture. No hash is frozen yet. The first run
is under build/port-map-sections/tile-first. Further section coverage remains.

## Reachable format paths

The root floor writer always starts with version4, and the root wall writer with
version7. Their historical helper dispatches (floor<=3, wall<6) can therefore only
occur when reading. Both helpers have no other production callers. Historical
write branches inside those helpers are unreachable and need not be recreated.
Test historical formats through the live roots with real read-mode IO; do not add
artificial direct write calls just to retain dead internal branches. Floor
versions below3 are rejected by the historical helper.

The current whole-grid floor encoding uses0xFFFF as its end marker. Cell(127,127)
with both halves present produces the same encoded word. Preserve and explicitly
capture this existing format limitation; silently changing the marker would
change compatibility. Region encoding uses different coordinate semantics and
needs separate coverage. Reading clears only the low flags byte in the grid,
which the complete-state contracts should distinguish from clearing a whole word.

The initial tile and current-floor contracts pass (2 roots /438 entries), including
the end-marker collision. The adapter now observes actual allocations for prefab
scratch records using the existing observer, preserving allocation/free behavior.
Unfrozen metadata cases check a valid scratch wall followed by a missing one, a
missing wall with a region pointer, and an unattached secret allocation.
Source review also found the current wall writer compares an unsigned flags byte
against zero when deciding to encode its0x80 flag; a pending independent wire
contract will check this against the flag bit already supported by its reader.

## Fixture stall and bounded observation

The first metadata attempt stalled during the preceding floor fixture's owner
reset. A SIGQUIT stack capture showed Go GC mark termination and a runnable
locked-thread test goroutine; no metadata assertion ran. The owned test process
was stopped and joined. Evidence is metadata-original/stacks.txt. This does not
establish a C algorithm failure or a specific runtime root cause.

The allocation observer had been enabled for the whole fixture. Match the
existing theme/growth fixture pattern by enabling it only around selected section
calls and scratch allocation, then disabling it before normal owner/file work.
No shared observer, production source, GC setting or toolchain was changed.
Retry metadata-bounded selects the suspected metadata defects and independent
wall-wire contracts. Region-read and wall-write drafts are installed; their
expectations remain unfrozen.

## Independent C defects confirmed

The bounded run completed both roots. Four missing-wall contracts failed: a
breakable/secret missing scratch lookup reused the preceding wall (counter19
instead of18); a missing window lookup changed the region word from0xb8 to0xf8;
and a missing secret lookup retained an unattached32-byte allocation. The wall
writer failed64 flag-preservation cases. Its unsigned-byte comparison always
took the flag-clear branch, despite the reader supporting the high bit.
Evidence: build/port-map-sections/metadata-bounded.

Apply small reversible prerequisites before freezing: clear lookup results when
a generation wall is absent; release an unattached secret record; interpret the
wall flag byte as signed in the existing encoding predicate. This changes no
wire layout. These are behavior corrections for later review, not translation
differences to hide in updated goldens. Fresh C production qualification is now
required. Corrected-first reruns all installed contracts, including24 region
cases and576 current wall-read cases. Hashes remain unfrozen.

## Expanded contracts

Corrected-first passed6 roots /1,238 entries. Historical-first passed10 roots
/2,435 entries, including original versions1–5 wall readers, version2 sprite
presence/absence in real definitions, full-width v3 floor payloads and unsupported
version rejection. Sprite handles use owned backing and are restored before
wall-definition snapshots. Breakable lists now have explicit observation and
cleanup in the new adapter.

Metadata-full passed all read contracts, including counter wrap and secret-state
defaults. Nine writer assertions initially omitted the existing owner iterator's
rule that excludes door/broken walls before metadata callbacks. Corrected the
independent expectations to include that rule; no production change. Current
region-expanded adds wall relocation and historical floor-region placement.
Source review identified a suspected missing factor23 in the historical odd-row
floor coordinate calculation; these contracts test the same placement as current
format loading before deciding whether to correct it.

The historical-region run confirmed24 placement failures, covering both live and
generation owners and even/odd rows. The v3 reader used grid Y directly where its
formulas require23*Y. Correct both expressions to match current-format world/grid
conversion. This is a second reversible prerequisite; no bytes/layout change.
Evidence: region-expanded/tests.jsonl. The same run passed current wall relocation
and all prior metadata contracts after fixing the writer expectations.

Corrected-expanded additionally tests region floor writing and1,014 existing-wall
overlay cases (versions1/3/7, merge off/on, all13x13 direction combinations).
Overlay expectations use the existing qualified composition table and independently
check owner reuse and retention of the wall's high flag.

Corrected-expanded passed the historical coordinate fixes and region-writer
contracts. Its504 failing overlay expectations read mapped bytes before the
fixture installed the immutable startup direction table. The actual fixture
correctly supplies that table, so fix the independent expected values to use
blobdata.PortTestMapPaintingTables; no production change. Focused-final repeats
all contracts and adds signed metadata-count boundaries and multiple-floor
serialization order. The latter checks that upper flag bytes do not affect bounds. Its first
expectation misread the C X-scan pointer type: it is uint8_t*, so both axes
inspect the low byte. Corrected the independent header bounds after reviewing
the declaration and observed output; record ordering itself already matched.

A reviewed cleanup removed33 superseded binaries from successful old production
runs, reclaiming1,603,731,468 bytes. No executable in use was selected. The latest
qualified statistics binaries, all logs/captures and original assets remain.
Paths, SHA256, sizes, original revisions and successful-run reports are recorded
in build/port-map-sections/superseded-binaries-cleanup.json; regeneration uses
the retained manifests at those revisions. Cleanup is consumed.

## Repeated C capture checkpoint

Focused-pass and focused-repeat both pass **18 roots /4,125 entries**, no skips.
All **18 captures /4,108 records** are byte-identical across processes and are now
frozen in the tests and map-sections-captures.json. Static mapped-memory checks
pass after freezing. This is a recovery checkpoint; accumulated three-target
and fresh production qualification remain pending. No Go conversion is installed.
freeze.py is consumed. Native IO/floor sketches under build are uninstalled and
unqualified; actual source and committed captures take precedence.
