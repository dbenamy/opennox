# Border selection — 2026-09-11

Scope: 543FB0/544020/544070/5440A0. The Go replacement preserves the three
ABI entries with remaining C callers; the private lookup no longer has a C
bridge. Production C after conversion: **140,672 physical lines (−58)** from
the preceding completed chunk, 153 files, zero reference C.

## Reproduced behavior defect

Byte4 (5440A0) checks its variation argument against the uint16 at
`0x85B3FC + 28688 + 60 * variation`. It should use the selected border row.
The loader in GAME1.c writes each row's limit as `2 * (width + height)`;
GAME4.c calls Byte3 with the border ID and Byte4 with its separate variation.
Placement in GAME4_3.c stores these as separate fields.

Both original-C probes call Byte3(7) then Byte4(16), with 64 active rows:

- Selected row 7 limit 20, row 16 limit 12: rejects valid variation 16 and
  retains the previous secondary value.
- Selected row 7 limit 12, row 16 limit 20: accepts invalid variation 16.

Limits 20 and 12 correspond to loader dimensions 5×5 and 3×3. These are synthetic
fixtures using realistic limits, not a claim that a particular shipped map has
these exact records. Existing placement callers ignore the validation return,
so false rejection can leave a stale variation. No visual defect has been
reproduced in gameplay.

## Approved repair

Preserve the disabled flag's success/no-write short circuit. For active selection,
validate the selected row against positive signed active count and physical
64-row capacity, then compare the nonnegative variation to that selected row's
limit. This also avoids indexing the table with large variation arguments.

Original-C baseline: `dd4a9f69`. The proposed regression fails against original
C on the valid-variation case, demonstrating that it detects the wrong-row bug.
The historical [repair patch](proposals/border-selection-selected-row.patch)
records the C change approved by the user on 2026-09-11. It was applied and
qualified before porting. Do not reapply it over the completed Go replacement.

## Approved repair baseline

The user approved the prepared repair on 2026-09-11. It is now applied to C,
and expanded qualification passes before Go replacement: 33,280 exact lookup/
name-selection checks, 25,600 repaired variation boundary calls, 63 primary
boundary checks and the original focused regressions. Lookup covers every active
count 0..64 and every physical row, first duplicates, empty names, non-ASCII
bytes and embedded NUL. Primary selection preserves signed count behavior even
when count exceeds the physical table because it does not access that table.
No new clamp is imposed there. Local artifacts: build/port-border-selection.

Temporary repaired C count: **140,733 (+3)**, 153 files, zero reference C.
Repaired/expanded C baseline commit: `e226f189`.

## Native implementation

`legacy/border_selection.go` preserves signed count and input behavior, the
first exact case-sensitive name match, C NUL semantics, and raw non-ASCII bytes.
Only NONE uses the shared C locale comparator. Name selection disables the flag
before lookup; misses preserve the previous primary and secondary. Numeric
primary selection rejects negative/out-of-count values without touching state.
Variation checks preserve any-nonzero flag activation and disabled success with
no writes, and include the approved selected-row bounds repair.

All expanded C checks pass against Go. Accumulated tests pass in default,
server and highres 386 variants; all three production binaries build as ELF32
Intel 80386. Fresh border-selection-port gameplay passes both preserved screenshot
checks with overrides disabled. Full suite matches the exact known 1,553 failure
entries (15 passing/3 known failing/32 skipped-no-test packages), with no added or
removed entries. Local artifacts: build/port-border-selection.
