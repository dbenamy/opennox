# Internal C-glue removal

Scope is defined in [PORT.md](../../PORT.md#goal-and-target). Keep native client
libraries/bindings; remove the engine's internal C dependencies while preserving
Linux 386/SSE2 behavior. This is a dependency plan, not a calendar estimate.

## Production dependency baseline

The inventory uses the same `./cmd/opennox` entrypoint as `internal/noxbuild`,
with default, `highres` and `server` tags. Before this phase, all three select
463 project cgo files in six packages: 458 in `legacy`, and one in each package
below. This differs from the 469-file whole-repository count because build
profiles exclude some files. There are still zero standalone production/test C
lines; 79 embedded production callback bodies remain.

| Package | Direct reason for cgo | Removal boundary |
| --- | --- | --- |
| `internal/binfile` | Empty import in `memfile.go` | Remove import; allocation remains transitively dependent on `alloc`. |
| `internal/netstr` | Linux `FIONREAD` constant from a C header | Use the equivalent Go syscall constant; preserve ioctl result and errno. Windows adapter is separate platform work. |
| `server` | Empty import in `audio_event.go` | Remove import; server still depends transitively on allocation and callbacks. |
| `legacy/common/alloc` | libc allocation, memory and string operations | Preserve lifetime, address stability, zeroing, bounds and ownership. |
| `legacy/common/ccall` | 76 generic C function-pointer dispatchers | Migrate all owners/identities before removing shared raw fallback. |
| `legacy` | Exports, C scalar/struct types, callback addresses, libc calls, three specialized adapters and build directives | Remove by connected owner groups, retiring interfaces and declarations together. |

Default/highres additionally select `go-gl/gl/v3.3-core/gl`, `go-sdl2/sdl` and
`go-openal/openal`. OpenGL is imported transitively by `libs/client/seat/opengl`;
it is not an unused go.mod entry. Server excludes these native client packages.
Standard-library `net` and `runtime/cgo` also appear in cgo-enabled metadata;
they are not engine-owned C glue.

Reproduce the inventory after sourcing the build environment:

```bash
python3 tools/porting/cgo_inventory.py --out build/port-cgo-audit-new
```

The tool records selected files, source hashes, importers, lexical C-selector
mentions and `//export` counts. It runs metadata discovery with cgo both on and
off; `go list -e` can succeed while packages contain errors or missing definitions.
In particular, server discovery with cgo off is not a successful server build.
Use actual builds/tests for acceptance. No dependency downloads or source changes
are performed by the inventory.

## Removal order

1. **Leaf imports and Linux socket constant.** Remove the three small direct cgo
   dependencies above. Qualify unchanged file/audio consumers and real loopback
   socket contracts, then fresh production builds/ABI and gameplay/save-load.
   This reduces directly cgo-using project packages from six to three, without
   claiming the server is cgo-free. Review further empty-import candidates separately.
2. **Allocation and libc ownership.** Inventory every direct allocation/free pair,
   including calls outside `alloc`. The six bounded memory/string helpers are
   now Go (see [GO_MEMORY.md](GO_MEMORY.md)), without changing allocation ownership.
   The [allocation centralization](RAW_ALLOCATION.md) moves 49 legacy calls and
   four tracked backend calls behind one libc boundary while preserving normal/raw
   versus safe/tracked domains. Review remaining string allocations, then replace
   the allocator behind that boundary with
   explicit lifetime and failure contracts. Do not substitute Go heap storage indiscriminately: raw 32-bit address
   words can outlive Go references, and pointer-containing layouts interact with GC.
   A stable non-Go-heap allocator is a candidate to evaluate, not an accepted design.
   Preserve failure behavior or record a justified correction before conversion.
3. **Callback identity and legacy type/export removal.** Use typed Go calls and
   registries; qualify a representation for persistent callback identities before
   retiring C addresses. Cover equality, null slots, aliases, mutable handlers,
   return conventions and actual owner calls. Ordinary Go function values are not
   interchangeable with single-word C callback addresses. Batch migrations by
   ownership; remove a fallback only after every reachable user is accounted for.
   Translate remaining C scalar/layout dependencies and retire unused exports and
   headers alongside those owners. Work that merely adds another parallel dispatch
   table must have a clear path to deleting the old route.
4. **Close the internal dependency graph.** Remove remaining build directives and
   engine libc calls, adapt test-only raw-C observers and retire generated bridge
   assumptions in qualification tooling. Complete a real server build with cgo off
   and qualify it. Then qualify both clients with external-library cgo retained.
   Use fresh reference comparisons, known-suite outcomes and meaningful performance
   checks; do not accept symbol deletion or compilation alone as completion.

These boundaries can overlap when one owner spans allocation, types and callbacks.
The prepared 26-monster-callback draft is deferred until it fits this removal
order; completing its parallel dispatch alone would not retire the shared bridge.

## Completion and progress measures

The milestone requires no engine-owned production C implementation, preamble
callback bodies, C-type dependency, exported C trampoline or libc helper path.
External library bindings are an explicit exception; do not count standard-library
or third-party cgo against internal completion. Keep x86/32-bit layouts where
useful. Investigate any actual external consumer before removing an interface.

Record direct project cgo package/file counts, remaining callback bodies and live
exports/types as well as standalone C LOC. Keep whole-repository and selected-build
counts distinct. Test-only C glue must be retired or explicitly documented as
remaining qualification work, rather than hidden by a production-only count.
Existing `safe`/ABI checks apply while bridges exist; replace checks tied to retired
C machinery with appropriate Go ownership/layout contracts as the boundary moves.

## Review notes

- The Luna external source audit missed transitive OpenGL use and inferred too much
  from the absence of direct repository imports. Primary dependency metadata caught
  this. Require selected dependency graphs for future external-dependency claims.
- External backend choice is deliberately deferred. Internal removal does not
  authorize dropping rendering/audio behavior or stubbing client services.
