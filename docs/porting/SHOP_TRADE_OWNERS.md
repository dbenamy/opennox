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
`build/port-shop-trade-owners/`. The six-file conversion is fully qualified.

## Qualification

All 33 roots pass without skips in default/server/highres. Safe/static checks,
three fresh production builds/ABI checks, the exact known-suite comparison and
headless character creation/save/load/resume pass. The suite retains its known
304 failure events, with 17 passing, two failing and 32 skipped packages.
All phases have identical source fingerprints; all six changed/deleted files
match independently reconstructed/formatted drafts. Original assets (1,654 hashes),
frozen expectations and external bindings remain unchanged.

Selected production cgo files fall 229→228 (235/463 eliminated on net), plus two
fixture cgo imports removed. C exports fall 1,081→1,078 (812/1,890 retired).
Production C callback bodies remain 78. Headers remain 157 files / 3,801 physical
lines. Standalone production/test-reference C lines remain zero. See
[qualification](shop-trade-owners-qualification.json) and
[inventory](shop-trade-owners-inventory-after.json).

Primary review simplified transport-only pointer/float roundtrips to the existing
object/session pointers and raw uint32 code/count values. It also checked signed
mode/type/index narrowing, nil session handling and returned bits against the
removed adapters. The independent trade-pickup C callback preamble is byte-identical.
The native public cancel wrapper remains; private conversion helpers have no callers
and are removed with the two adapter files. Luna's bounded draft was useful after
this review; its initial revision metadata was corrected before acceptance.
