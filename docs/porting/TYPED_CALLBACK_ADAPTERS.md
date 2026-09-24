# Typed callback adapter consolidation

Status: baseline and shared-dispatch conversion qualified in default/server/highres/safe.
All 24 calls per profile pass. The original-adapter baseline reused production
evidence after verifying unchanged source and binary hashes from `83dc68ac`. See
[baseline evidence](typed-callback-adapters-c-qualification.json). The exported Go APIs
`legacy.CallDrawFunc` and `legacy.Nox_call_objectType_new_go` now route their
calls through the existing `ccall.CallIntPtr2` and `ccall.CallVoidPtr` dispatchers.
This removes duplicate C invocation bodies without changing callback signatures
or requiring an architecture-specific function-pointer calling mechanism.

The draw shim reads `nox_drawable.draw_func`; the replacement reads the
existing Go `Drawable.DrawFuncPtr` field. Its placement is therefore part of the
ABI contract. The object wrapper has no current in-tree caller, but preserving
its exported API avoids inferring downstream compatibility from a source search.

Luna drafted a focused test and C callback fixture; primary reviewed them before
installation. Tests invoke the existing production wrappers with C-heap objects.
Two distinct draw callbacks record exact viewport/drawable pointers and return
all five signed boundary values: INT32_MIN, −1, zero, one, INT32_MAX. Each call, with a live or null viewport pointer,
must select the correct callback and invoke it exactly once. The object callback
records pointer identity and invocation count across four calls alternating live
and null object pointers.

The fixture callbacks do not dereference the supplied objects. Assigning a
callback through `DrawFuncPtr` and invoking the original C shim explicitly tests
the Go/C field layout mapping. Null function pointers are not valid inputs and
are not invoked. This is callback ABI testing, not a claim of renderer-pixel or
object-initialization behavior coverage. Existing consumer tests provide that
higher-level evidence where applicable.

Standalone C remains zero lines/files. Production C preamble
bodies fall from 81 to 79 after qualified consolidation.
Test-only C fixture bodies are excluded from the production count and from the
standalone `.c` metric.

Primary extended the draft to valid null data-pointer cases, moved the test to
external package `legacy_test` to avoid root-package build/init overhead, and
enabled server coverage because both production wrappers compile there.

## Conversion acceptance

Both APIs now use existing shared dispatchers. Only `legacy/drawable.go` and
`legacy/object_type.go` change production behavior. All 24 calls per profile
pass unchanged. Twelve rendering/object consumer roots pass in both client
profiles without skips; safe build/static checks pass. Fresh client/highres/server
builds and four-binary ABI checks pass, with both retired shim names absent.
Full suite matches exactly: 304 known failures, 17 pass/2 fail/32 skip packages.
Headless character creation/settings and save/load match existing references.
See [native qualification](typed-callback-adapters-native-qualification.json).

One initial consumer invocation rejected a two-line pattern before tests started.
Its logs are preserved separately; final phases use the corrected single-line
pattern. No golden or existing expected failure changed.

The original consolidation left 79 production preamble bodies. The later
[native-record storage batch](GO_NATIVE_RECORD_STORAGE.md) retired the unused
sprite-iteration adapter after a complete reference review. **78 now remain:**
76 generic ccall dispatchers plus the spell and curve-segment adapters.
No project-owned header function bodies remained in the reviewed inventory.
These callback boundaries remain compatible with the current 386/SSE2 plus cgo
target; retiring shared dispatch requires migrating its live callback identities.
