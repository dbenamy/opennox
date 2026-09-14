# Quest eligibility

The next connected batch contains 14 routines /563 C section lines in GAME3_3.c,
from 4F24E0 through before 4F3E30. Production remains C while the new baseline is
captured. Four entry points serve the common server join path; ten private C
entry points can retire. Existing Go quest-penalty callers will use native
helpers. The current production count remains 99,071 C lines /147 files, with
zero reference-only C.

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
qualification scripts are under build/port-quest-eligibility. Native source, its installer and qualification scripts are staged there. No C
eligibility algorithm has been replaced yet.

The native fixture will call the four retained production C ABIs directly and
private Go helpers for the ten retired entry points. Thus the retained argument
and return conversions remain exercised without keeping a C algorithm solely
for testing. Native source is still staged at this baseline checkpoint.
