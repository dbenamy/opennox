# Floor and wall map sections

## Current result

FloorMap, WallMap, WindowWalls, DestructableWalls and SecretWalls now use Go.
All18 selected C bodies/interfaces and two private C globals retire. The five
registered Go callers invoke Go directly; no C algorithm is retained for tests.
Historical readers remain reachable; proven unreachable historical writer branches
are omitted. Selection/reachability evidence is in
[map-sections-selection.json](map-sections-selection.json).

Corrected C baseline71633546 is pushed. Native qualification passes **19 focused
roots /4,144 entries**, with19 captures /4,126 records matching C. Each target
passes **510 roots /46,994 entries**, no skips. All **282 captures /103,153 records**
match C and each other; all native gates have identical source. Static checks,
three fresh production binaries/ABI, exact known asset-suite outcomes, gameplay,
explicit save/load and compressed flat-map regeneration pass. See
[map-sections-native-qualification.json](map-sections-native-qualification.json).

C remaining: **27,001 physical lines /67 files /zero reference**. The−1,793
reduction includes−1,619 bodies/global definitions and−174 obsolete address-heading
and blank lines in the touched files. The implementation preserves the corrected
C behavior, including IO call boundaries that affect the existing map checksum.
No frozen expectations changed during conversion.

## Contract scope

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

## Three-target C gates and additive light contract

Capture checkpoint **4cc84212** is pushed. Accumulated default/server/highres
gates each pass **509 roots /46,975 entries**, no skips. All **281 captures
/103,135 records** match across targets, including the18 frozen section captures.
All three gates have identical2,468-file source. See
map-sections-core-c-qualification.json.

Translation review added a real-polygon contract: nondefault and zero level,
containment versus nearby-edge fallback, outside fallback100, and wall-region
inclusion/exclusion. Its18 cases pass twice and its capture is frozen. The new
total is **19 focused roots /4,144 entries**, **19 captures /4,126 records**.
This is an additive test-only change; qualify that group on each target and
verify that all preexisting source hashes are unchanged, rather than repeating
the same accumulated groups. Fresh production qualification uses the final source.
The native batch will run the whole affected selection including the new group.

A second reviewed cleanup removed33 more superseded successful-run binaries,
reclaiming1,606,402,944 bytes; manifest suffix cleanup-2.json. Both cleanup modes
are consumed. Original assets, latest qualified binaries and evidence remain.

## Qualified C baseline

All C gates now qualify: **510 roots /46,994 entries per target**, counted as the
509-root accumulated selection plus the independently qualified light addition.
All **282 captures /103,153 records** match across targets; the19 new section
captures contain4,126 records. Source identity is verified within the accumulated
gates and within the additive/production gates; their sole source difference is
the new light-test file. No existing production/test file changed between them.

Static checks, three fresh production binaries/ABI, the exact known1,553 failure
entries /3 packages, gameplay, explicit save/load and compressed flat-map
regeneration pass. See map-sections-c-qualification.json and c-production.
The C count is28,794 /67 files /zero reference. Go drafts are ready for integration;
no native source is installed at this checkpoint. qualify-c.py is consumed.

## Native integration in progress

Qualified C baseline **71633546** is pushed. Installed all18 native replacements,
retired their C interfaces and two private globals, and moved the five Go callers
and original fixture directly to Go. Historical writer-only branches disappear;
all live historical readers remain. Current C is27,001 physical lines /67 files
/zero reference:−1,619 from body/global removal and−174 obsolete address-heading
plus blank lines removed from the two touched C files (including prior empty
headings). Total−1,793 from corrected C.

The first native comparison matched all file bytes, return values and state.
Only metadata checksums differed (257 records across four captures): coordinates
had been split into two4-byte IO calls. Existing cryptfile checksum updates invert
once per call, making call boundaries observable. Restore one8-byte coordinate
operation for window/breakable reads and all metadata writes. No golden changed.
Evidence: native-focus; corrected run:native-focus-fixed.

Completed C asset copies were hash-verified and deduplicated, reclaiming
1,660,044,319 bytes. Audit/apply modes are consumed; restore is supported by
deduplicate-map-sections-c-assets.py and each run's restoration manifest.
Installer/draft copies are consumed; actual native source supersedes them.

The corrected native focused sweep passes all19 roots /4,144 entries; all19
captures /4,126 records match frozen C exactly. Static checks pass. Default
accumulated gate also passes; remaining gates and fresh production are underway.
A third audited superseded-binary cleanup reclaimed1,609,262,940 bytes; its
cleanup-3.json manifest is consumed. Original/latest artifacts remain protected.

## Final native qualification

All native gates and production are complete. Raw evidence is under
build/port-map-sections/native-{focus-fixed,default,server,highres,production}.
The initial native-focus checksum-only mismatch is retained as diagnostic evidence;
its corrected successor passes all frozen records. All installers and qualification
scripts are consumed. Original assets and committed C baselines provide recovery.

Native scenario copies were also hash-verified and deduplicated after qualification,
reclaiming1,660,044,319 bytes. Audit/apply modes of the native asset cleanup script
are consumed; restore remains available. Changed maps/saves, screenshots and all
comparison evidence remain. Original assets/archive are unchanged.
