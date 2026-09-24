# Shop and trade fixture owners

## Scope and baseline

Use existing native shop/trade owners directly from the pooled fixtures. Remove
three shop C exports and their prototypes, six C-typed shop adapters, five native
trade forwarding adapters, and C imports made unnecessary by these migrations.
Production algorithms, frozen expectations and pointer-word normalization stay
unchanged. Preserve the price fixture's null session, signed 32-bit inputs and raw
32-bit code/count/result words in repair and trade operations.

The original baseline reuses 33 shop roots from each freshly qualified
inventory/resource profile at `9dac58bf`. Every root passed without skips, source
fingerprints match exactly, and recorded execution environments are retained.
No new capture or expectation generation is needed. See
[baseline](shop-trade-owners-baseline.json), [selection](shop-trade-owners-tests.txt)
and [manifest](shop-trade-owners-batch.json).

GPT-6 Luna owns the bounded uninstalled draft; the primary owns reconstruction,
semantic review, integration and qualification. Local artifacts:
`build/port-shop-trade-owners/`. Conversion is not yet installed or qualified.

## Qualification plan

Rerun the exact 33-root selection in default/server/highres, then safe/static,
three fresh production builds and ABI checks, exact known-suite comparison and
headless character creation/save/load/resume. Verify source identity throughout,
unchanged original assets and external native-library bindings. Record the actual
post-conversion cgo/export/header counts after qualification.
