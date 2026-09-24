# Retire C interfaces without native callers

Status: baseline frozen from qualified source `f6f5ee4c`; reviewed conversion
not yet installed. See [baseline](go-only-exports-baseline.json) and
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
