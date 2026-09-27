# Fixture callback retirement

Status: original baseline accepted; reviewed draft not installed.

Remove two generic C dispatch signatures whose remaining uses are fixture-only:
`CallVoidUPtr2` and `CallVoidInt2`. Keep every live production callback route.

- Radial collision/motion fixtures already call native production functions with
  Go closures. Record their ordered pointer/code words directly in those closures,
  preserving KeepAlive, empty results, callback count, pointer width and order.
  The motion observer retains its 32-row bound; exceeding this fixture-only bound
  now reports a Go panic instead of C.abort. Production behavior is unchanged.
- Retire `Nox_call_objectType_new_go`, which is called only by its own assertion
  block in `TestLegacyCallbackAdapters`. Keep that test's drawable assertions and
  C observer unchanged; the drawable callback route is still live.
- Retire `prefabGroupEach`, called only recursively and by operation 0 of the
  prefab fixture dispatcher. Remove that operation, its C observer and its sole
  traversal test. Other operation numbers and their assertions remain unchanged.
- Remove the two signatures from the generator allowlist and regenerate ccall.go.
  No replacement dispatch table or C implementation is added.

Whole tracked-code reference checks include function values and fixture callers.
The prefab dispatcher also has an indirect painting callback: its placement and
instantiation owners use retained operations 18/21 and 15 and are included in the
selection. Variable operation arguments in paths/metadata are {5,6}/{12,13}.
Only the retired traversal test invokes operation 0. See the committed reference
inventory and baseline function review. Original AST review: four changed,
ten removed and 67 unchanged functions, zero additions.

Original baseline uses exact-source `b2c975dc` root results with source,
supplemental-source, binary SHA256, runtime environment and discovered-name checks.
All 29 selected roots passed in each of default/server/highres. Fresh original
private adapter tests passed in default/server/highres/safe; two radial contracts
also passed in safe. These are established frozen assertions, so no new repeated
capture set is needed. The sole removed root exclusively tests the retired helper;
require all other 28 exact roots after conversion, plus private adapter and safe
radial checks. Frozen live-behavior expectations do not change.

Qualification requires focused radial preflight/static, the matching 28-root
three-profile selection, private adapter checks in four profiles, safe radial
contracts, safe build, three production builds/ABI, exact known-suite outcomes,
fresh headless save/load and original-asset integrity. Check both retired C
signature substrings are absent in production symbols. Regeneration and source
review must prove all retained dispatcher bodies are unchanged.

Expected qualified counts: five client / six server production cgo files,
17 embedded C bodies, zero legacy exports and zero standalone C lines. Two
fixture C imports disappear: 103 to 101 source files across build tags, distinct
from selected production counts. Headers remain 157 files / 2,731 lines.

Primary drafted/reviewed locally; Luna remains quota-unavailable. Artifacts and
consumed baseline scripts are in `build/port-callback-retirement-audit/`.
Seven obsolete string-batch executables were removed only after committed-source,
qualified-replacement/hash and host-use checks (387,117,056 allocated bytes).
Rebuild `32df9553` using retained commands/source records; current baseline uses
qualified spell binaries. Original logs and source maps remain.
