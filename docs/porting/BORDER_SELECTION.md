# Border selection investigation — 2026-09-11

The next port is 543FB0/544020/544070/5440A0. Production C remains **140,730
physical lines**, 153 files, zero reference C; this baseline removes no C.
The original-C fixture passes direct lookup, name selection, primary selection,
secondary selection and signed-count checks. It saves and verifies actual table,
guards and four shared state words, restores them, and checks inputs unchanged.
This is an investigation baseline, not yet exhaustive port qualification.

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

## Proposed repair

Preserve the disabled flag's success/no-write short circuit. For active selection,
validate the selected row against positive signed active count and physical
64-row capacity, then compare the nonnegative variation to that selected row's
limit. This also avoids indexing the table with large variation arguments.

Original-C baseline: `dd4a9f69`. The proposed regression fails against original
C on the valid-variation case, demonstrating that it detects the wrong-row bug.
The reviewable [repair patch](proposals/border-selection-selected-row.patch)
contains the C repair, changed expectations for both reproductions, and 25,600
boundary calls across active counts, selected IDs, flag values, row limits and
variation values. The fixture's input/table/state/guard checks remain enabled.
The proposed repair passes all focused border tests in default/server/highres
386 variants. It was then removed from the working source; only the proposal
patch is retained. Production builds, gameplay and full-suite validation are
still required when adopting the repair and completing the port.

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
Next port the quartet, retaining only the three ABI entries with C callers.
Full chunk validation, source count, documentation, commit and push follow.
