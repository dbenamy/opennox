# Retire C interfaces without native callers

Status: baseline frozen at `ac3a842c` from qualified source `f6f5ee4c`;
conversion fully qualified. See [baseline](go-only-exports-baseline.json) and
[batch manifest](go-only-exports-batch.json).

## Scope and reachability

Retire 265 C export directives and 244 header prototypes across 69 Go files and
36 headers. Remove 126 unused wrapper functions completely; preserve the other
139 functions and every Go caller. Remove 15 now-unnecessary C imports/preambles
and three unused `unsafe` imports. No full source files, tests, globals, callback
IDs or external native bindings are removed. Expected selected counts are
397→382 project cgo files and 1,444→1,179 legacy exports in each profile.
Standalone production/reference C remains zero; 79 embedded callback bodies
remain outside this batch.

Primary independently scanned lexical Go identifiers, C selectors, comments,
literals, C preambles, headers, assembly and all other tracked non-documentation
inputs. All 265 names lack native callers. The 39 comment/literal references
are historical annotations or source-rewriter/parser test inputs, not runtime
native lookups; preserve them. No ignored/untracked Go/C/header/assembly inputs
were found under `src`.

Package resolution matters: the initial textual inventory treated same-named
functions in different Go packages as callers. Primary separated 126 unexported
functions with no owning-package Go identifier reference beyond their definition.
Delete those wrappers; preserve all 139 others conservatively. Header declarations
are accepted only as exact simple prototypes. The qualified executable build
surface is the supported interface; no external shared-library API is promised.

## Review and baseline

Luna produced a bounded deletion draft. Primary reconstructed all 105 outputs
from accepted original lines, verified all 126 function ends using Go's parser,
and compared every remaining non-import declaration through the AST. A second
AST comparison after import cleanup/formatting also passes. Header edits equal
exactly the 244 accepted prototype lines. No fixtures or frozen outputs change.

C-import candidates have no remaining C selectors, exports or build directives.
Their preambles contain only includes, externs, prototypes or typedefs; local
header closure review checks for executable definitions and side effects. Primary
caught an omitted multi-blank-line preamble in the helper's `config.go` audit and
reviewed its extra headers before accepting removal. `unsafe` imports have no
package-initialization effects. Keep all live Go declarations in those files.

Reuse the preceding completed qualification by exact production/test source
fingerprint identity in all phases, rather than rebuilding an identical baseline.
The baseline report records the source commit and evidence. All new qualification
must run after installation: complete root porttest corpus in all three profiles,
exact discovered/started/completed name sets and only the known diagnostic skip;
safe build/static; fresh three-profile production/ABI with all 265 retired names
absent; exact known suite failures/package outcomes; and headless character
creation/save/load/resume. Six formerly required ABI names move to the retired
set; 290 other required names remain.

## Delegation outcome

The bounded Luna edit draft was useful and mechanically checkable. Primary kept
reachability algorithm design, package-reference resolution, acceptance and final
qualification. Helper audit corrections included a wrong caller file path and an
omitted C preamble. These reinforce checking actual source and retaining textual
scan limitations; no measured cost-saving claim is made.

Local drafts/review: `build/port-go-only-exports/`. `draft-v2/` and `reviewed-v2/`
replace the earlier directive-only drafts. Completed scripts are consumed; do
not rerun them against later source.

### Follow-on audit finding

During this batch's qualification, a separate Luna frontier audit proposed 25
object-state export removals but omitted calls inside `object_state_porttest.go`'s
C preamble. Primary rejected that proposal before any installation; those wrappers
remain untouched. The accepted 265-export cohort had independently scanned all
comment/literal references and its tests compile. Narrow subsequent helper work to
explicit scalar-edit spans; keep C-preamble reachability classification with the
primary. The ignored frontier is incomplete evidence, not the next accepted batch.

## Completed qualification

All gates pass on identical source fingerprints. Default/highres each pass 2,425
root tests plus the expected map-population diagnostic skip; server passes 2,414
plus that skip. Discovered, started and completed name sets exactly match the
frozen complete-corpus inventory. There are no failed events or unexpected skips.
The prebuilt controller took 2,141 seconds including sequential compilation and
at most two concurrent root sweeps. Safe build and static checks pass.

Three fresh production binaries pass ABI checks; all 265 newly retired C symbols
are absent, while required interfaces remain. The broader suite exactly matches
304 known failure events and 17 passing/two failing/32 skipped packages. Headless
character creation and explicit save/load/resume pass. All test/fixture source
inputs and 1,654 original asset hashes remain unchanged. No safe runtime contract
run is claimed for this export-only batch.

Measured in every production profile: project cgo files 397→382; legacy exports
1,444→1,179. Phase totals: 81/463 cgo files eliminated on net and 711/1,890 exports
retired. Headers: 157 files / 3,902 physical lines. Three direct project cgo
packages, 79 embedded callback bodies and external native bindings remain.
Standalone production/reference C stays zero.

Evidence: [qualification](go-only-exports-qualification.json) and
[selected inventory](go-only-exports-inventory-after.json).
Local completed gates: `build/port-go-only-exports/{contracts,safe,production}/`.
