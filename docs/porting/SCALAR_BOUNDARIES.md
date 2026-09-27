# Native scalar owner and caller boundaries

Status: original qualification accepted; conversion is an uninstalled draft.
Production baseline: `a7dc3a36`.

Seventeen production files still import C solely for character, floating-point,
boolean and integer types. Replace these with the exact native representations
emitted by the qualified Linux 386 cgo/GCC toolchain: `char` is `int8`, `uchar`
is `uint8`, `float` is `float32`, `bool` is `bool`, `longlong` is `int64`,
`ushort` is `uint16`, `short` is `int16`, and the 32-bit signed/unsigned types
retain their widths. The isolated generated-type proof is under
`build/port-scalar-boundaries/type-probe/`.

The draft covers 24 private C-typed signatures and their callers. Package AST
inspection found no non-call references to these functions. Six string-reading
owners and one return-value caller use the existing `GoStringP` bridge, retaining
its C.GoString behavior. Real libc calls and string allocation/ownership remain.
Message and damage calls use the existing Go owners with equivalent arguments,
narrowing and ignored returns. Arithmetic and existing float store/rounding
boundaries remain unchanged. One newly empty fixture C import and two unused prefab fixture prototypes are
also removed. Whole-source header/preamble searches found no live C callers.

The original selection includes 968 default/high-resolution roots and 966 server
roots. Server omits the two `!server` hover tests. The established opt-in
`TestMapPopulationPrerequisiteProbe` diagnostic is explicitly excluded. Selection
combines a reverse package-function caller audit, transitive fixture/root helper
calls, entire affected assertion files and conservative owner-family inclusion.
An initial overbroad graph that confused common method/field names with package
functions was rejected before execution. This selection deliberately includes
extra related owner tests rather than claiming every root is a direct caller.

Existing contracts cover signed boundary values, flag predicates, clocks, text
buffer identity and guards, byte-string lookup, prefab record precision and
ownership, damage, poison, saved-player framing and serialized messages. No new
width-mirroring tests or changed root assertions/goldens are needed for this
representation-only batch. The 25-root preflight covers representative scalar
boundaries before the complete converted selection. Qualification still requires
all three production/ABI builds, safe/static checks, exact known-suite comparison,
a fresh headless save/load/resume and original-asset integrity.

Baseline reuse is restricted to identical source/supplemental fingerprints,
verified binaries and matching runtime settings. The prior integer selection
supplies 166/164/166 first observations; 802 missing roots per profile and the
complete selections passed in separate processes before installation. Each
selected root has two passing original observations with exact name sets and no
failures/skips. See [baseline evidence](scalar-boundaries-baseline.json).

Luna remains quota-limited; primary handles this batch without a substitute model.
Before testing, 198 inactive Linux 386 Go cache archives were removed after exact
path/stat/hash and host-use checks, reclaiming 4,976,803,840 allocated bytes.
Source, assets, module downloads and qualified executables remain. Cache misses
rebuild normally. Journals: `build/port-scalar-boundaries/cache-headroom-*`;
cleanup script consumed. Proposed old-executable archival plans were not executed.

Qualified counts remain 52/53 selected project cgo files, zero legacy C exports,
20 embedded C bodies, 157 headers / 2,731 physical lines, and zero standalone
production/test-reference C lines until converted qualification is accepted.
