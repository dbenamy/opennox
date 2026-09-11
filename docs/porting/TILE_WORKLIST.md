# Tile-fill worklist — 2026-09-11

Scope: push 51DD50 and pop 51DE30, owned by the C tile-fill routine 5437E0.
Keep both ABI entries while the owner remains C. The worklist has 500 three-word
records at 0x973F18+16200, overflow at +22200, and count in the shared C variable
dword_5d4594_2487248. The separate grid has 128 row pointers and 128 44-byte cells
per row; keys are the words at offsets 4 and 24 within a cell.

Push accepts coordinates 1..126 and requires a matching key in at least one
flag-selected half (`flags & 1` selects offset 4; `flags & 2` selects offset 24).
It also rejects y==1 for mask 1 and x==1 for mask 2. No grid is accessed for out-of-range
coordinates or flags with neither bit set. Duplicate identity is the exact
coordinate/flags triplet, excluding key. Other flag bits are retained.

Duplicate scanning uses signed count, while capacity uses unsigned count. At
capacity, an exact duplicate returns without setting overflow; a qualifying
unique item sets overflow to 1 without writing queue/count. High-bit counts skip
the signed scan and still meet unsigned capacity. Popping a signed nonpositive
count returns zero without touching outputs, including null output pointers.

Pop is LIFO for independent output buffers. It decrements count and then writes
X, Y, flags in order, reloading count and queue after each write. Outputs that
alias count or queue therefore have observable effects; disassembly confirms the
reloads. Pop does not clear inactive records. Overflow stays unchanged unless
an output pointer explicitly aliases it. Restrict successful pop tests to count
and alias writes that keep every subsequently accessed queue index valid.

## Original-C baseline

The installed fixture passes 9,000 original-C operations: 7,972 coordinate/flag/
key-match enqueue-pop steps, nine capacity/overflow transitions, six duplicate
identity steps, eight signed-negative-count steps, an empty null-output pop,
1,000 combinations of pop output aliases, and four explicit LIFO checks.
Matrix pushes are immediately drained so previous duplicates cannot hide a
broken eligibility branch. Overflow starts poisoned to detect unintended resets.

Two separately allocated C grids provide exact row-by-row comparison. Snapshots
check all 1,500 queue words, count, overflow, grid contents/pointer, output guards,
and adjacent queue guards. The fixture restores and verifies the actual original
global state. Safe aliases include separate outputs, shared outputs, count,
overflow, older queue entries and the current record's fields. The current X
value differs from the decremented count, proving the first count reload matters.

Positive counts above 500 are outside the tested worklist invariant because C
can scan beyond its physical records. No copied C algorithm or assets are used.
Local artifacts: build/port-tile-worklist. Production C before this chunk: **140,785 physical
lines**, 153 files, zero reference C.
