# Internal audio bridge identities

## Scope and original contracts

Remove 26 engine-owned forwarding exports: 13 stream/device table callbacks,
eight event/music bridges and five unreferenced AIL compatibility wrappers.
Keep the AIL Go backend and external OpenAL/SDL2/OpenGL bindings. Five unused
private event wrappers can also disappear with their export-only file.

Original production source is qualified `6f981385`. Three new test-only files
check original C callback identities, mutable hooks and music-slot addresses.
The selected 130 roots cover audio owners and their callers in each supported
profile. The two new roots repeat in separate processes. Existing independent
assertions, captures and foreign-callback observers remain unchanged.
See [baseline](audio-bridge-identities-baseline.json) and
[test selection](audio-bridge-identities-tests.txt).

## Conversion contract and review decisions

Keep the 13-word table mapping, including three unused voice words, as distinct
Go identities. Their original stop/abort/zero behavior remains available through
typed dispatch; the new table contract checks their identities but does not invoke
those dormant bodies. Source tracing finds no reader of the three raw words.
This preserves layout and latent behavior without retaining C entrypoints solely
for an unobserved caller. Revisit if an actual consumer is discovered.

Look up mutable forwarding hooks at call time. Preserve each caller's integer
versus discarded-return convention and the foreign C callback fallback. Migrate
every selected descriptor/voice/event field consumer before registering Go keys.
Preserve signed 32-bit and signed-byte results extended to 64-bit fixture captures;
pointer results instead retain zero-extended 32-bit address bits. Music-slot index
arithmetic retains its existing 386 narrowing/wrap behavior.

The five removed AIL compatibility exports have no in-tree consumers or references
in the scanned cached dependency Go/C/header/assembly sources. Actual backend
methods stay. The event device fixture's three private C-typed helpers can use Go
uintptr/uint32/int32 types without changing device hooks or buffer ownership.

Luna prepares an ignored overlay; primary acceptance requires independent mapping,
whole-reference and source-hash review. The production table mapping was independently
traced through both initializer variables and blob offsets. The legacy raw-address
GC test exercises unchanged chunk/buffer binding, outside this callback conversion.

## Qualification status

Original baseline `403efa06` is committed and pushed: all 130 roots pass per
profile with zero skips, and both new contracts pass twice per profile. An initial
run omitted asset variables and skipped RealCatalog; only its two asset-independent
new-contract passes are reused as repetition evidence. The accepted affected run
uses every runtime variable from the batch manifest.

The 14-path Luna overlay passed primary hash, AST, table and caller reconstruction.
Primary caught and restored the known-name zero fallback before installation, and
requested an explicit int32 device-ready result (int was equivalent on 386).
Existing owner algorithms and root assertions remain unchanged.

Conversion qualification is complete.
Passed converted gates: all affected roots on default/server/highres, safe/static,
three production/ABI builds, exact known-suite comparison, fresh preflight/final
headless save/load, asset hashes and measured cgo/export/header inventory.
No current-source full-root corpus result is claimed.


All 130 converted roots pass per profile; production and scenario gates pass with
unchanged reviewed source. See [qualification](audio-bridge-identities-qualification.json).
Exports fall 288→262; selected production cgo files fall 149→146 client and 150→147
server. Headers remain 157 files with 3,024 physical lines. Embedded production C
bodies remain 77; standalone production/test-reference C remains zero.
