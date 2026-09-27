# Native quantity and image completion callbacks

Status: original baseline accepted; seven-file native draft remains uninstalled.
Parent qualified revision: `5233d084`.

Replace quantity and image completion's foreign fallbacks with native identity
registries. Preserve six quantity handlers, two book handlers, stable keys,
caller-side nil guards, mutation order, integer widths and allocation lifetime.
The quantity callback still receives the same five raw words; its temporary point
remains allocated by the shared allocator and freed after callback completion.
The image callback receives the same image reference. Register controlled fixture
observers in init, preserving hook replacement/restoration and existing assertions.
Retire CallVoidUPtr5 from its generator and generated Go/C code. CallVoidPtr remains
for other owners. SDL2/OpenGL/OpenAL bindings are unchanged.

## Producer and representation review

Four production uiAmountShow sites (buy, sell, repair, inventory drop) provide
native accept/cancel identities or nil. The two stored fields are offsets 1319160
and 1319100; their writer is uiAmountShow and readers are uiAmountCallback and
uiShopCancelQuantity. The latter compares unchanged stable accept identities.
The interaction escape fixture clears the cancel word. Other observer identities
come from PortTestTradeUIObserve through the existing controlled fixture facade.
CallVoidUPtr5 has no other production consumer.

ImageRefAnim producers initialize OnEnd to nil in things.go or set one of the two
bookImageEndKey identities in gui_book_construct.go/gui_book_pages.go. Existing
fixtures use those identities, nil, or PortTestObserveImageEnd. Audio voice/event
fields also named OnEnd belong to different records and are outside this change.
The sole image completion dispatcher is guarded in gui_cursor.go. ImageRefAnim's
16-byte layout and all raw callback fields remain unchanged.

The draft review compares all six quantity case bodies with the registered
closures, both book mappings, and unchanged allocation code. Generator output
matches the draft byte-for-byte. Fixture observations retain the two signed point
words and four uint32 metadata words; opaque integer values stay integers.

Unregistered keys deliberately panic instead of entering a raw foreign callback.
This is a reversible invalid-input behavior change after tracing producers, not a
claim that panic diagnostics/recoverability equal C crashes. Do not apply this
conclusion to other callback families without their own audit.

## Qualification scope

Select 133 default/highres and 132 server roots from 66 owner files: trade/shop quantity lifecycle, inventory
window callers, spellbook owners, image completion and escape controllers.
Run fresh original selections using the exact-source qualified parent binaries,
checking binary/source/environment identity. Repeat six focused contracts in each
profile: quantity accept/cancel/nil and GC plus image timing/reference, mutation
order, noncompletion and book owners. Existing captures and assertions remain
unchanged. There is no need for another C observer: these contracts already cover
the boundaries being replaced.

Require matching converted selections in default/server/highres, the focused
contracts repeated separately, safe build/static, three production builds/ABI,
exact known-suite outcomes, fresh save/load and all 1,654 original asset hashes.
Record exact names, not counts alone. Existing rendering fixture owners are not
claimed safe-mode contract coverage; the safe production build remains a gate.

Expected counts: production cgo files unchanged at 5 client / 6 server;
embedded C bodies 15→14; fixture C-import files 100→98; standalone C 0;
legacy exports 0; headers 157 files / 2,731 lines.

Primary only: Luna's previously reported quota remains exhausted. Artifacts:
`build/port-ui-completion/`. Completed scripts are single-use.

The first baseline controller expected 133 tests on server too and stopped after
the successful 132-test server run. TestClientInventoryWindowWorldSelection is
explicitly !server. Correct the controller count, verify compiled-name sets per
profile, retain those successful unchanged-source runs and resume the remaining
profiles/repeats. No test or source correction is required.

All 133/132/133 original roots passed without skips in default/server/highres;
the six focused contracts passed again in each profile. Exact compiled-name sets,
binary/source hashes and relevant environment settings were verified. Evidence:
[baseline](ui-completion-baseline.json). Native qualification remains pending.
