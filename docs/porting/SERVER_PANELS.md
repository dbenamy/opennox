# Server panels

## Scope and status

Corrected C baseline qualification passes. Parent **f5737e6a** is the
qualified/pushed server-options conversion. No panel translation is installed yet.
Current production C: **48,755 physical lines / 79 files / zero reference C**.

The batch contains **57 live functions / 2,098 corrected C body lines**, plus
one **37-line orphan** to remove. It covers the remaining five servopts translation
units and GAME2 453080..4559B0, GAME3 4AD570..4AD820, and GAME3_1 4BDC10..4BE610.
The neighboring 4BD280..4BDC00 allocation/cache family is a separate owner and is
excluded. The rule-file picker 4CEBA0 remains a dependency exercised through its
actual implementation. See [the scope](server-panels-scope.json).

`sub_4AD4B0` has only its definition, header declaration and one call inside literal
`if (0)` in general.c. Remove it with the conversion. Disabled construction branches
elsewhere do not automatically make their event handlers unreachable; retain live
callbacks and interfaces according to their actual callers.

## C prerequisites and decisions to review

Five constructors now return zero immediately if their root window resource fails
to load: 4530C0, 453850, 4541D0, 4BDC10 and 4BDFD0. Separate-process contracts first
reproduced six crashes (weapon and armor share a constructor). General options
already handled a missing root. All seven missing-resource/retry cases pass after
these reversible fixes, which add **15 physical C lines**. Evidence is in
`build/port-server-panels/focused-3.log` and `focused-4.log`. Because production C
changed, this baseline requires fresh production qualification.

Preserve the following observed behavior during translation; these are review
items for later game-behavior cleanup, not extra corrections in this batch:

- Weapon-mask queries inspect the highest nonzero byte of positive masks. For
  nonpositive masks they inspect the preceding byte at 1045451. Fixtures own that
  byte explicitly, including zero and the signed boundary.
- Generic spell-bit helpers truncate the word quotient to one byte and use 386
  shift counts. Guarded fixtures allocate the complete 256-word addressable range;
  boundary tests never rely on an address outside the owner.
- Object checkbox display skips class indices 0, 1, 4 and 5 in both object modes;
  click events resolve the displayed row's name instead. Preserve both paths.
- The upper level-limit checkbox path combines an unshifted edit value with the
  low nibble, while the numeric-edit path shifts it into the high nibble. Both
  paths have separate expectations.
- Player lookup is case-insensitive; finding a displayed row for removal is exact.
- Advanced audio settings retain the requested signed value in settings/text,
  while the actual audio threshold clamps it to 0..100.

## Fixtures and contracts

The fixtures reuse serverOptionsOwner, actual GUI parsing, list/entry widgets,
software rendering/fonts, teams/players, reliable queue and existing rule owners.
Controlled resources use shipped IDs/types. New panel fixtures explicitly match
multi-selection ownership of class/user lists in shipped `access.wnd`; previous
server-options fixtures retain their original resource configuration. Optional
old audio controls are taken from `lwadvsrv.wnd` to exercise their live event path.
Original assets and the archive are unchanged.

Each private state word and backing region has explicit snapshot/restore ownership.
The first shared-mask preflight caught guard checks written as extracted-global
lookups; those checks now inspect owned raw backing slices. Player selection first
exposed a fixture error: active players must occupy their real slot indices for
iteration. This was corrected without changing production code. C-allocated
admission-list entries are released through their real remove event before sentinel
storage is restored, including assertion-failure cleanup.

**22 roots / 22 frozen captures / 13,066 captured records** currently pass in the
focused C run (3.057s). These include independent assertions beyond capture hashes:

- 196,608 exhaustive byte-mask operations, word masks, pointer identities, exact
  copy width and no aliasing of caller storage; all 256 class indices; signed and
  quotient/shift boundaries for guarded spell-bit storage.
- Actual spell definitions and shipped weapon/armor tables: eligibility, enabled
  state, every final-return branch of spell snapshot, application bounds and
  preservation of unrelated words/definitions.
- Missing root and successful retry for seven panel variants; nine language
  indices on each side of the large-font resource switch.
- Spell/object population and labels, host/client controls, bulk events, checkbox
  state, scrolled spell toggles, empty rows and restricted-spell notices observed
  at the existing dialog API boundary.
- Admission numeric completion/focus events, clamping and truncation, neighboring
  fields, settings notification, checkboxes, multi-class selection, real name-list
  add/remove, player selection/lookup/departure and closed-panel notifications.
- General toggles across game modes, advanced tab cycles and repeated lifecycle,
  live mask refresh through the shipped callback table, settings-copy boundaries,
  advanced edits and actual audio-threshold mutation.
- Drawing at nonzero/clipped positions with opaque, transparent and image paths;
  general drawing's no-op behavior; admission action availability during drawing.

The baseline does not claim exhaustive end-to-end coverage of every networking,
rule-file or dialog implementation. Existing lower-level tests and production
scenarios supplement the bounded panel contracts. Review arithmetic widths,
return conventions, ownership and callback dispatch against C during translation.

## Qualification

Frozen expectations are in `src/server_panels_*porttest_test.go`. The focused
selection is [server-panels-focused-tests.txt](server-panels-focused-tests.txt);
the manifest is [server-panels-c-batch.json](server-panels-c-batch.json).

Default/server/highres each pass **168 roots / 12,914 leaves**, no skips, with
**104 identical captures / 26,048 records** including the parent corpus.
Durations: **72.43 / 163.07 / 83.80 seconds**. All gates use the same unchanged
**2,144-file source manifest**.

Fresh production passes in **391.28 seconds**: three builds and ABI/interface
checks; the exact known asset failure set (**1,553 entries; 15 passing / 3 failing /
32 no-test packages**); headless options gameplay, save/load and forced flat-map
regeneration. Client SHA-256:
`d9486bf2dc88515630d9b0b025b9c2f8d6ea4f2460c4a9a432d23968af4f8a71`.
All sessions joined. Static memory preflight and whitespace checks pass.

The external-reference audit identifies nine function interfaces still required by
production C: `sub_4AD840`, `sub_4535E0`, `sub_4535F0`, `sub_4536B0`,
`sub_453710`, `sub_453F70`, `sub_454040`, `sub_4540E0`, `sub_455800`.
Twenty-six private globals can move to Go; migrate the advanced-server getter and
other Go callers directly. The online-mode flag remains shared. The old advanced
refresh callback table can become native once its only consumer is translated.

Local evidence is under `build/port-server-panels/`. Freeze and installed fixture
scripts are consumed: do not recopy stale drafts or regenerate hashes to conceal a
difference. Commit/push the qualified C baseline before installing native Go, then
qualify, document C LOC, commit/push the conversion and continue.
