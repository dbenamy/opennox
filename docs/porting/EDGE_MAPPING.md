# Border edge mapping — 2026-09-11

Scope: 543E60 and 543EB0, both called by the C tile-list owner 543C50. Keep both
ABI entries until that owner moves. Normalization 411490 remains C in this chunk.
The edge table is 64×60 bytes at 85B3FC+28644, width/height at row+52/+53.
The merge table is 12×12 uint32 at 587000+282736. Production C before this chunk:
**140,672 physical lines**, 153 files, zero reference C.

## Original-C baseline

1,145,856 operations pass: 786,432 direct transformations (all 256² dimension
pairs × 12 edge inputs) and 359,424 dependent merges (all 64 physical rows,
12 normalized classes, 12 mapping columns, 13 mapping values, three dimension
profiles). Profiles are 3×3, 5×6 and 255×255. Direct inputs include signed extremes
and default wrapping. C uses -fno-strict-overflow; the Go oracle uses explicit
uint32 modular addition for the default branch.

Exact 3×3 is identity with no RNG use; edge0 is always zero with no draw.
Edges1/6 use width-based IntClamp; edges3/4 use height-based IntClamp. Existing
Go IntClamp returns its maximum without a draw when max-min==-1 (dimension2);
other reversed spans retain its signed-modulo behavior. Tests cover these exact
results and Logic RNG index after every call, including repeated table wrap.
Other RNG stays unchanged. RNG is supplied by an isolated temporary server;
the original server getter is restored.

The fixture configures mapping cells from independently supplied expected
normalization classes. All other cells are sentinel255, preventing a wrong row
or column from accidentally matching a previous test. It checks every table
byte/word and adjacent guards after each call, and verifies actual restoration.
C-owned records have guards and poison words; only current-edge word+12 may
change. Mapping255 returns0 with no change or RNG use. Otherwise it returns1,
including when the resulting edge equals the original word.

Dependent inputs stay within the caller's normalized 12×12 mapping contract.
No copies of C algorithms or assets are added. Local artifacts:
build/port-edge-mapping. Original-C baseline: `606f5aa7`.

## Native implementation

Both live ABI entries now execute Go. The generator preserves the existing Logic
RNG helper and fixed int32 arithmetic; merge retains the C normalization call,
sentinel early return, and writes only when the value changes. Production
C is **140,613 physical lines (−59)**, 153 files, zero reference C. All 1,145,856 focused calls match the C baseline. Accumulated default/server/
highres tests pass, all three production binaries build as ELF32 Intel80386,
and fresh edge-mapping-port gameplay passes both preserved screenshot checks
with overrides off. The preceding border chunk verified the full suite's exact
known1,553 failure-entry multiset; this immediately following chunk did not repeat it.
