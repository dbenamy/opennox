# Optional safe-profile forwarding shims

Status: direct Go exports qualified against the actual-C runtime baseline.
Standalone C is now **25 physical lines in three files**, down 20 lines and one
translation unit; zero reference C. The six public functions are Go exports with
const-qualified pointer typedefs and explicit C comparison return types. Macro
redirection, sanitizer flags, foreign allocation and function semantics remain.

All 142 frozen cases and the 129-case shop consumer pass under safe,porttest.
The safe production build and static check pass. Three fresh production binaries
pass ABI checks; the known-suite failure/package sets match exactly; headless
character creation and explicit save/load pass. The six old `_go` names are absent
and all ten live callback addresses stay distinct. See
[safe-bridges-native-qualification.json](safe-bridges-native-qualification.json).

The actual-C runtime baseline was captured against pushed entry conversion
`e007a425`. That baseline kept production wrappers unchanged at 45 C lines/four
files. Three fresh processes produce the same 142-case capture, frozen as
`2b34e573e3b8f5b62ced16258ec2de832284fbe9d0cecc3a0c5120edf22eb59a`.
The frozen contract and 129-case shop consumer pass under `safe,porttest`, with
both discovered roots completed and no skips. See
[safe-bridges-c-qualification.json](safe-bridges-c-qualification.json).

The optional `safe` build routes legacy C allocation and memory/string calls
through Go allocator exports and enables AddressSanitizer plus mapped-memory
checks. Six C wrappers previously adapted const-qualified pointers to `_go` exports.
The conversion exports those Go implementations directly using named C
const-pointer typedefs, retaining the public names, libc/allocator behavior,
macro redirection and sanitizer settings. A compiled isolated 386 cgo probe
confirms that generated C prototypes preserve the typedefs' const qualification.
This probe establishes type compatibility, not game/runtime correctness.

The new fixture calls the actual six C shims on two separately allocated foreign
buffers. Inputs and operations are bounded to 64-byte usable regions, surrounded
by 16-byte sentinels. Results snapshot both entire buffers and normalize pointer
returns by exact destination equality. Independent root tests check all memcpy
sizes 0..64, mixed alignment, bounded comparisons and unsigned byte order, NUL
termination, high bytes, empty strings, exact capacities, return values, source
preservation and guards. Comparison sign assertions express the portable C
contract; raw integers remain in the current-target capture for parity. No
reference C bodies are copied into tests.

`TestShopStockLoading` also reaches the real production strcpy branch: its
porttest type table installs FieldGuideXfer at index 29, object creation preserves
that Xfer pointer, and shopLoad copies the reward name into the field-guide use
buffer. The fixture captures that data. This chain was verified from source;
its 129 captured cases now also pass under the safe profile.

Luna drafted the bounded test matrix and root test against a primary-defined API.
Primary review corrected a strlen input-side mismatch and added mixed-alignment,
bounded-comparison and exact-capacity cases before execution. Primary wrote the
actual-C bridge and owns acceptance. Source reports also required regex escaping
corrections. No net subscription savings have been measured.

Baseline reuse was limited to unchanged production builds: official Go file
selection excluded both new `safe && porttest` fixtures from all four production
variants. Every preceding source hash and reused binary hash was verified, along
with the ten callback identities. Safe runtime checks ran fresh. After conversion,
all production/gameplay gates ran fresh; only cgo_safe.go and deletion of
cgo_safe.c differ from the baseline source fingerprint. No test/golden changed.
