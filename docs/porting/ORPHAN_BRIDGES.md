# Unused memory and GUI bridges

Status: baseline qualified; removal not yet applied.

The preceding formatting conversion is qualified and pushed as `59bf181d`.
This batch removes unreachable C memory accessors and five GUI wrappers, together
with their private Go exports. It replaces the durability test's sole remaining
C memory lookup with `memmap.PtrUint64(0x581450, 9608)`. The durability algorithm,
its live C ABI, threshold mutation/restoration, and batched fixture stay unchanged.

## Reachability and scope

A whole-source literal search, including headers and Go preambles, found:

- Twelve `mem_get*` accessors: definitions/declarations, the tagged durability
  quarter accessor, and offline source-analysis metadata/tests only.
- Five GUI adapters and their five `_go` exports: definitions/declarations only.
  Existing Go GUI methods retain old C names in historical comments.
- Two `safe` memory lookup trampolines and their exports: used only by the removed
  accessors. Preserve allocator/string/memory substitutions, ASan flags and
  `memmap.SetRuntimeChecks(true)`.

The ignored machine-readable search is
`build/port-final-orphan-bridges/before-references.json`. References in
`src/internal/{blobs,callgraph,noxfactor,offalign}` are offline tooling; references
in `src/client/gui/window.go` and `src/save.go` are comments. The remaining
`common/memmap` Go registry and mapped data remain live and unchanged.

Retire 24 C interfaces. Keep all shared global storage, callback identities,
window structure definitions and live window callbacks. Unreachable algorithms
need no replacement implementation or new behavioral fixtures.

## Validation

Four existing test roots cover the changed durability accessor and nearby native
window creation, callbacks, deferred destruction and drawing. The GUI tests are
integration smoke coverage, not claims that the removed adapters were reachable.
See [orphan-bridges-tests.txt](orphan-bridges-tests.txt).

Default baseline: all four roots pass, no skips; static checks pass. The optional
`safe` build failed before removal because `nox_free` passed `unsafe.Pointer` to
the now-generic `alloc.Free`. Use the existing `alloc.FreePtr` API instead. This
one-line repair is limited to the `safe` build tag; the qualified default/server/
highres production behavior is unchanged. The repaired safe build passes. Source manifests differ only in that tagged
file; reuse the preceding qualified default/server/highres production evidence.
See [orphan-bridges-baseline.json](orphan-bridges-baseline.json).

After removal, run the four roots and static checks in default/server/highres,
rebuild `safe`, and run fresh production/ABI, exact known-suite comparison,
headless gameplay and save/load. Also audit retired symbols in source and linked
binaries. Do not count an optional profile's build as runtime qualification.

## Delegation

GPT-6 Luna prepared the guarded installer and existing-test shortlist. Primary
review found and corrected an invalid deletion marker before execution. The
installer preflights every input before writing and has a no-write `--check`.
Primary independently searched references and selected validation. The shortlist
was precise and distinguished adjacent GUI smoke tests from direct coverage.

## Storage cleanup

After the formatting gates joined, verified duplicate assets in the completed
`text-format-native-save` run were removed: 556,358,986 bytes. Restore using
`python3 build/port-final-formatting/deduplicate-format-save.py --restore text-format-native-save`.
Forty-eight identical completed native captures now share verified canonical
`build/port-final-formatting/c-probe3` files, saving 362,014,506 logical bytes.
The sharing manifest is `build/port-final-formatting/native-capture-sharing.json`.
Preserve canonical targets; original assets and result logs are untouched.
