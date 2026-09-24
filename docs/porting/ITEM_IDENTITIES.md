# Native item use, drop and pickup identities

## Original baseline and scope

All 181 selected roots pass without skips in default, server and high-resolution
profiles against `8c9b3fc6` plus three test-only files. Production fingerprints
match that qualified checkpoint exactly; its production, ABI, known-suite and
save/load evidence supplies the production baseline.

Replace 41 distinct C addresses: fourteen Use, eleven Drop and sixteen Pickup
registrations. Retire their 41 exports and two adjacent Go-only inventory fixture
adapters. Shared owners retain distinct registration identities. Preserve data
sizes, nil AmmoUse/BowUse registrations, ten mutable hooks and six script-facing
Use getters, plus the existing Drop getter.

New contracts cover all keys stored in foreign ObjectType/Object records across
GC/stack growth, late replacement of all mutable hooks, getter identities and
full signed Use return words through the actual item-use path. The configurable
independent C observer defaults to its previous result of one. Existing assertions
and captures are unchanged. Focused roots include inventory, item effects, world
collisions, book awards, scripts, player files, unit gameplay, orchestration and
resource definitions.

The typed Use.Get API returns bool, but effectsUse consumes the full C integer.
Add a native int32 result registry and explicit result/discard methods. Unknown
addresses retain the original integer or void C calling convention respectively.
Food pickup still discards the result; original caller gates and nil behavior stay
unchanged. Keep object layout and the existing bool cache. This reversible API
choice requires a full default port corpus sweep after focused qualification.

Primary owns API design, baseline acceptance and fixture migration; one GPT-6
Luna helper drafts the bounded production overlay. Review exact keys, owners,
input narrowing, result widths, sparse capture IDs and remaining references before
qualification. Add independent native/raw dispatch contracts after conversion.
Run all focused profiles, full default port corpus, safe/static checks, three
production/ABI builds, exact known-suite comparison, headless save/load, source
fingerprints and original-asset hashes before accepting the conversion.

[Original evidence](item-identities-baseline.json),
[commands](item-identities-batch.json), [selection](item-identities-tests.txt).
Local artifacts: `build/port-item-identities/`.
