# List and player-group C bridge removal

## Scope and baseline

Retire the 13 wrappers in `legacy/lists_exports.go` and their header prototypes.
Their remaining consumers are three porttest fixtures. Production root-package
wrappers already call native Go APIs; keep those and all owner algorithms unchanged.
Whole tracked-source references, including complete preambles and address uses,
were audited. No live in-tree production C caller or callback registration remains.

Original source is qualified `1ca13c2a`. All 19 roots pass on default/server/highres,
without skips/failures, using its verified same-source binaries and runtime settings.
Its production evidence is reused for this documentation-only baseline. Tests cover
list ordering/mixed/nullable behavior, player-group lifetime/boundaries/bounded names,
root wrappers, client sequence delivery/wrap and rule-file list ownership.
See [baseline](list-glue-identities-baseline.json) and
[test selection](list-glue-identities-tests.txt).

## Review and qualification plan

Use the existing native owners directly in list/group, rules and session fixtures.
Preserve the 12-byte node layout, 32-bit address results, signed indices, nil traversal
and existing C allocation/free ownership. Group removal may return an address inside
a just-freed group; compare only address bits as before. Keep every assertion and
frozen expectation. No callback registry or algorithm translation is needed.

Luna supplied a bounded caller audit and five-path uninstalled overlay. Primary
verified original/draft hashes, exact 13 prototype removals, wrapper-to-owner mapping,
whole-reference coverage and AST changes; only formatting was needed. An initial
name-only primary test closure over-selected unrelated methods named Close/State/Find.
The accepted selection follows unique fixture constructors and actual rules/session
entrypoints instead. All 19 roots are available on every profile.

After baseline commit, install and run all 19 affected roots per profile, safe/static,
three production/ABI checks, exact known-suite comparison, fresh preflight/final
headless save/load, original asset hashes and measured inventory. No current-source
full-root corpus is claimed; last full run was the earlier shared drawable batch.
Conversion remains pending.
