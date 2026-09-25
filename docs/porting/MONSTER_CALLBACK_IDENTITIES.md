# Monster callback identities

## Scope and behavior

Retire the 26 engine C wrappers in the shipped strike, dying and dead callback
tables: 25 wrappers in `legacy/ai_callbacks.go` and the Bomber forwarding hook in
`legacy/object.go`. Keep the Go owner algorithms and Linux 386 record layout.
All consumers are in legacy, so a private registry is sufficient. Use separate
result and void-discard dispatch helpers. Melee branches on the integer result;
dying uses the integer calling convention but ignores the result; dead uses the
void calling convention. Preserve raw C fallback for foreign callbacks, including
the independent combat/lifecycle observers.

Use distinct identity keys even for owners with identical behavior. Strike wrappers
reinterpret the original float argument's bits as an object address; the new typed
route passes that object directly, without numeric floating-point conversion.
Bomber's strike returns one without an owner call. Death adapters perform their
original side effects before returning one. Bomber's dead callback reads its mutable
root hook at dispatch time and preserves its signed result.

## Baseline and contracts

The conservative affected test closure includes direct callback owners, combat and
lifecycle consumers, resource definition loaders and all connected PortTestRoam
fixtures. The operation map remains sparse: 0–10 strike; 11–19 dead slots
[0,1,2,3,4,6,7,8,9]; 20–24 dying. Do not shift stable capture IDs or regenerate
expectations.

A new original-path Bomber forwarding contract covers all 26 table identities,
exact restoration of the table region, mutable hook replacement, nil/live object
arguments, signed 32-bit limits and both result/discard routes. It stubs the root
hook; the unchanged root Bomber gameplay owner's item/no-item branches are not
newly qualified by that contract. Existing captured owner tests cover the other
callback operations. The original Bomber contract passed twice in separate processes per profile,
and all 564 affected roots pass in each profile with exact names and no skips.
Original baseline production source is identical to qualified `024b2632`; its
production evidence was reused for that test-only baseline. The conversion is now fully qualified.

## Primary and helper review

Luna supplied a bounded caller audit and original-hook test draft. Primary added
unconditional deferred table restoration and nil coverage for void dispatch. The
first dispatch audit incorrectly inferred several table positions from wrapper
source order. Primary caught this before implementation and independently extracted
all 26 actual bindings through GAME_data_init.go and blobdata/blob_init.go, comparing
them with the isolated fixture writes. Preserve the original audit and corrected
version; only the verified table mapping is implementation input.

Local evidence: `build/port-monster-callback-identities/primary/` and
`build/port-monster-callback-scout-20260925/`. Original baseline committed/pushed as `28642775`; conversion fully qualified.


## Conversion review and qualification scope

The reviewed overlay spans ten source paths. Primary verified all original/draft
hashes, formatted the Go drafts, compared functions through the Go AST, and checked
every owner call, return, prototype removal and production/fixture binding.
Gameplay owner algorithms and existing assertions are unchanged. Primary also moved handler-map population into init, following existing registry
setup and keeping gameplay-owner dependencies out of package-variable initialization.
No owner or callback mapping required correction in the final overlay.
Its revision label was stale; actual source hashes match `024b2632` plus the two
original-only contract files. Acceptance uses those verified hashes.

Qualification covers the 564 affected roots on all three converted profiles, safe/static checks,
three production/ABI builds, the exact known asset-suite comparison, fresh
preflight/final headless save/load scenarios and original asset hashes. The
private dispatcher has three audited production consumers, covered by the affected
selection; the preceding shared drawable-update batch supplies the last full root
corpus run. Do not present that earlier result as a run of this conversion.


## Qualified result

All 564 affected roots pass on default/server/highres against exact original names,
without failures/skips. Safe/static, three production ABI checks, exact known
asset-suite comparison, both fresh headless save/load scenarios and all 1,654
original asset hashes pass. See [qualification](monster-callback-identities-qualification.json).

The private dispatcher has three audited production consumers; this batch uses
focused coverage plus production qualification. The immediately preceding drawable
update batch `024b2632` supplies the last full default corpus run, explicitly on
that earlier source rather than this conversion.

Legacy exports fall 359→333; selected production cgo files remain 152 client and
153 server because those files retain unrelated C interfaces. One isolated table
fixture no longer imports C. The 157 tracked headers contain
3,095 physical lines. Embedded production C bodies remain 77;
standalone production and test-reference C remain zero. External bindings are unchanged.
