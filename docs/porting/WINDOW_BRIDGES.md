# Window helper bridge retirement

## Scope

Retire 19 exports: six geometry/ancestry/text interfaces in `client_ui_window.go`,
nine wrappers in `gui_window.go`, and four in `gui.go`. Whole-source review finds
no runtime callers for the latter 13: only definitions, headers, translator rules
and method comments. Retire those wrappers without replacement; remove the unused
field-94 forwarding macro with its declaration. Preserve translator mappings and
the unrelated root-package `sub_46AEE0` string helper.

Two of the six live interfaces serve inventory: point containment and strict
ancestry. Point containment already has a native owner. Move ancestry to a native
window owner, preserving nil handling. The four fixture-only offset/position/text
boundaries become native fixture adapters, consistent with the earlier window
retirements. Preserve fixture operations, independent assertions and frozen captures.

In particular, preserve sequential output writes when position output pointers
alias; local Y and each ancestor's X/Y accumulate in the same word. Nil local-offset
queries zero both outputs; nil global-position queries leave outputs untouched.
Preserve signed 32-bit coordinates, inclusive edges, text-event ordering and raw
text-pointer results. Keep shared C window types/layout and unrelated private
helpers for their remaining users.

The connected inventory input file's last C dependency after point migration is
an integer cast in a nonzero word comparison. The equivalent native comparison
allows that import to retire along with `client_ui_window.go`'s import.
Dependency discovery and qualification confirmed both removals.

## Original baseline and qualification

[Baseline](window-bridges-baseline.json) records exact source/environment identity
with qualified quantity revision `5585ab65`. Reuse its 394 client / 391 server
roots and production gates. Run seven additional, disjoint window roots twice per
profile in separate processes, using the same qualified binaries; all pass with
no skips. No test source or expectations changed, and no recompilation was needed.
The converted selection is the complete union: 401 client / 398 server roots.

The seven roots cover geometry/state, tree/labels, nil outputs, resize event order,
bounds/ancestry, maximum ID ranges and aliased position outputs. Inventory, widget,
shop/trade, quickbar and other GUI contracts remain in the reused selection.
After conversion require all selected roots in three profiles, safe/static, three
production/ABI builds, exact known-suite comparison, two headless save/load runs
and unchanged original asset hashes. Converted results below are accepted.

## Delegation and review

One Luna helper audited the 19 names and drafted a bounded ignored overlay;
primary owns reachability, source review and qualification. The two independent
scans found 93 references across 13 files. Primary corrected the scout's claim
that inventory input would retain other C calls, and that `gui_window.go` had
another live C string helper; the actual remaining dependencies are described
above. Shared C record types remain outside this batch. No measured usage savings
are claimed.

## Qualified conversion

All 401 client / 398 server roots pass with exact original names and no
failures/skips. Safe/static, three production/ABI builds, exact known-suite
comparison, two headless save/load scenarios and all 1,654 original asset hashes
pass on unchanged reviewed source. See [qualification](window-bridges-qualification.json).

Exports fall 210→191; selected production cgo files fall 130/131→128/129.
Headers remain 157 files / 2,947 physical lines: 19 prototypes, one unused macro
and three blank lines retired. Embedded production C bodies remain 77; standalone
production and test-reference C remain zero.

Primary reviewed the eight-file overlay, every changed function, remaining symbol
references and the unchanged C window record prefix. Raw offsets, sequential
aliased writes, nil guards, strict ancestry, signed coordinates and text-pointer
results are preserved. Pre-compile review corrected a stale fixture comment and
removed an unused import. Root assertions and captures remain unchanged.
