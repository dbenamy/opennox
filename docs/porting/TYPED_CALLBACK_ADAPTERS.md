# Typed callback adapter consolidation

Status: original-adapter baseline qualified in default/server/highres/safe.
All 24 calls per profile pass; production source and retained binaries are
unchanged from `83dc68ac`. See
[baseline evidence](typed-callback-adapters-c-qualification.json). Keep the exported Go APIs
`legacy.CallDrawFunc` and `legacy.Nox_call_objectType_new_go`, and route their
calls through the existing `ccall.CallIntPtr2` and `ccall.CallVoidPtr` dispatchers.
This removes duplicate C invocation bodies without changing callback signatures
or requiring an architecture-specific function-pointer calling mechanism.

The draw shim reads `nox_drawable.draw_func`; the replacement will read the
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

Standalone C remains zero lines/files. Before conversion, production C preamble
bodies remain 81; after qualified consolidation the expected count is 79.
Test-only C fixture bodies are excluded from the production count and from the
standalone `.c` metric.

Primary extended the draft to valid null data-pointer cases, moved the test to
external package `legacy_test` to avoid root-package build/init overhead, and
enabled server coverage because both production wrappers compile there.
