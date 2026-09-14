# Quest eligibility

This completed batch contains 14 routines /563 C section lines in GAME3_3.c,
from 4F24E0 through before 4F3E30. The native conversion is fully qualified. Four entry points serve the common server join path; ten private C
entry points are retired. Existing Go quest-penalty callers use native
helpers. The production count is now 98,506 C lines /147 files, down 565 lines
(563 section lines and two obsolete declarations), with zero reference-only C.

The guarded fixture uses actual type, modifier, equipment and guide registries,
plus bounded records in the production table regions. It saves and restores
those tables, registry owners and cached IDs; object and modifier allocations
use the existing guarded objective owner. The first original-C run passes every behavioral assertion across 10,793 cases
in nine groups (27.989s). Only the expected unlocked-capture checks fail. Both original-C runs produced byte-for-byte identical captures; all nine
hashes are now locked. The accumulated standard baseline passes all 632 selected root tests
(369.139s wall), with every selected test executed and completed. Sixteen
additional double-input rounding cases pass and repeat exactly, detecting omission
of C's intermediate float32 assignment. Their locked focused run verifies its
one selected root test executes and completes (17.503s wall). The final baseline
is **10,809 cases /ten locked capture groups**. The original 10,793 cases ran in
the complete matrix; the additional 16 root-test-only cases ran separately.

The groups cover scalar spell/beast/ability admission, books, modifier
masks and exclusions, duplicate/sentinel rows, effect slots, composite checks,
item class precedence, special equipment, inventory limits and destroyed items.
The inventory tests include all twelve potion limits, nested inventories,
nil/nonplayer callers and finite/nonfinite float-to-integer staff limits.
Expected behavior is checked separately before capture hashes can be locked.

The source audit found two easily confused cases: guide names use case-sensitive
strcmp, while the special UserColo equipment prefix uses an eight-byte
case-insensitive comparison. Grouped guide IDs skip their first marker and can
match their terminating zero; ordinary eligibility tables stop at zero.

Use repeated original-C captures and the complete accumulated standard baseline.
For the completed native batch, use the full standard suite and affected
server/highres suites (eligibility, penalty, inventory, equipment and rewards),
all three production builds and fresh gameplay. A differing variant must first
be replayed against the recoverable C baseline. Broaden for variant-sensitive
changes or unexplained discrepancies; see DECISIONS.md for this testing scope.

Local plan, original sections, ABI audit, staged fixture and test sources, and
qualification scripts are under build/port-quest-eligibility. The qualified C baseline is 8c1c2a94. Native focused eligibility and penalty
tests pass in 29.040s, matching all ten hashes unchanged. The native full standard matrix passes all 633 selected root tests (379.259s),
covering 96,142 accumulated captured cases /1,039 groups plus contracts.
Affected server/highres suites each pass all 67 selected tests (171.678s and
77.589s). All selections execute and complete. All three production builds pass
ELF32/SSE2 and symbol checks: ten retired symbols absent, four required exports
present, test helpers absent. The full suite has the exact same 1,553 failure
entries (15 passing, three failing, 32 skipped packages). Fresh headless gameplay
with unchanged goldens passes in 36.344s. Full local metadata is recorded in
build/port-quest-eligibility/qualification.json.

The native fixture calls the four retained production C ABIs directly and
private Go helpers for the ten retired entry points. Thus the retained argument
and return conversions remain exercised without keeping a C algorithm solely
for testing. No original C eligibility algorithm remains. The existing C-backed cache word
and shared table data remain production state, as in neighboring native ports.
