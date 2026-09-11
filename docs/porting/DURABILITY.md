# Durability classification — 2026-09-11

Scope: `sub_57B190`, used by inventory description and server item-destruction
reporting. Keep the C entry for both live caller sites. Inputs are uint16 values;
returns are zero maximum → 4, full durability → 0, at least half → 1,
below quarter → 3, otherwise → 2. Overfull values remain in band 1.

The half threshold is the shared C `qword_581450_9544`; the quarter threshold
is blob 0x581450+9608. Preserve those reads and ordered comparisons. The similarly
named local constant in client/drawable.go does not modify this C global.

## Baseline design

Each threshold pair checks 786,424 current/maximum pairs: every uint16 maximum,
current 0/1/65535, maximum−1/maximum/maximum+1, and values immediately around its
quarter and half. Out-of-range candidates are omitted. Default .5/.25 thresholds
use an independent integer-rational oracle. Sixteen altered pairs cover reversed
thresholds, both signed zeros, infinities, quiet/signaling NaNs in either slot,
non-rational normal values, extreme positive doubles and negative thresholds.
These use ordered float64 comparison expectations; no NaN output bits escape.
Total: 13,369,208 classifications, with exact result codes checked.

A small C fixture loop invokes the live entry in batches. It temporarily sets
threshold bits, verifies that calls leave them unchanged, and restores both
original values. It contains no copied classifier/reference implementation and
requires no assets. Local artifacts: build/port-durability.

All cases pass against the original C before replacement. Production C before
this chunk: **140,903 physical lines**, 153 files, zero reference C.
