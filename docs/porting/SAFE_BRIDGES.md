# Optional safe-profile forwarding shims

Status: actual-C runtime baseline passes against pushed entry conversion
`e007a425`. Production wrappers are unchanged; standalone C remains 45 lines/four
files. Three fresh processes produce the same 142-case capture, frozen as
`2b34e573e3b8f5b62ced16258ec2de832284fbe9d0cecc3a0c5120edf22eb59a`.
The frozen contract and 129-case shop consumer pass under `safe,porttest`, with
both discovered roots completed and no skips. See
[safe-bridges-c-qualification.json](safe-bridges-c-qualification.json).

The optional `safe` build routes legacy C allocation and memory/string calls
through Go allocator exports and enables AddressSanitizer plus mapped-memory
checks. Six C wrappers currently adapt const-qualified pointers to `_go` exports.
The proposed change exports those Go implementations directly using named C
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

The intended baseline reuse is limited to unchanged production builds: new files
are tagged `safe && porttest` and must be shown excluded from production file
selection. Verify every previous source hash and reused binary hash before
acceptance. The new safe runtime contracts/captures must run fresh. After the
conversion, rerun safe runtime/build checks and fresh production qualification.
