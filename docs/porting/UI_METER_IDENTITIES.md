# UI meter bridge conversion

## Scope and baseline

Remove 20 exports: 17 meter bridges in the state/items/draw files, the two meter
constructors, and the adjacent inventory weapon-draw bridge. Primary audited all
116 symbol-reference sites in 17 source files, including preambles, prototypes,
callback installation and fixture normalization. The constructors have Go/fixture
callers; the adjacent draw address is otherwise only a fixture normalization key.
Move those consumers together. Keep owner algorithms and record layouts unchanged.

The original baseline production source matches qualified audio revision `9eac418e`. Two test-only files add
contracts for the installed weapon tooltip and the foreign tooltip ABI. The new
installed-route test constructs the actual weapon window, checks its registered
callback, and checks equipped/empty tooltip text with varied argument words.
The independent C observer checks all three foreign argument words, exactly one
call, and suppression for nil windows, nil/dead callback words and destroyed
windows. Its C body remains test-only evidence for the retained foreign fallback.

The selected 363 client / 360 server roots include meter, inventory, interaction
and trade UI owners, plus reverse references from unique constructors and tooltip
dispatch. Three
hover/world-selection tests explicitly require `!server`; this is build selection,
not a runtime skip. Existing pixel/state captures remain unchanged.

All 363 client and 360 server original roots pass without failures/skips;
both new roots also pass in separate repeat processes per profile. See the
[accepted baseline](ui-meter-identities-baseline.json) and
[test selection](ui-meter-identities-tests.txt).

## Implementation boundaries

Use existing typed `gui.WindowFunc` and `gui.WindowDrawFunc` support for meter
events/drawing; a broad change to the generic C wrappers is unnecessary. Preserve
nil-versus-nonzero event response conventions, 386 argument/result widths, callback
capture identities, owner algorithms, rendering order and mapped record layouts.

Tooltip dispatch needs a native route while retaining the foreign callback path.
Use a small callback-identity registry consulted by `Window.TooltipFunc`
after its existing nil/dead checks; this preserves the stored callback word and
normalization while avoiding changes to the C-compatible Window layout. The
helper initially suggested an extension-field setter; the registry preserves the
existing callback-word identity without adding another per-window field. This is
a reversible implementation choice. Add native registry contracts for exact
arguments, distinct keys, replacement, foreign fallback, invalid registrations
and destroyed-window guards; the original contracts protect the installed route.

The inventory weapon-draw bridge and its fixture callers migrate together. Keep
its normalization entry at the same index with a stable identity. Preserve all
10 meter callback operation IDs and every existing native fixture case. Event
adapters retain the original event-code/argument evaluation and nil-versus-raw
response conventions; draw functions retain their integer results.

## Qualification status

Original baseline `4792e80b` is committed and pushed. Luna prepared the bounded
14-path overlay; primary verified original/draft hashes, every source diff, exact
prototype removals, callback mappings and retained function bodies. The installed
conversion adds two primary contracts: native tooltip dispatch boundaries and the
actual installed event callback's evaluation order. They supplement the original
baseline and are not claimed as original-C differential captures.

Primary caught one draft ordering difference: argument evaluation preceded event
code evaluation. Luna corrected it before compilation. The new contract mutates
the event code while evaluating arguments, so reversing that order fails directly.
The bounded implementation handoff was useful; final acceptance remains with the
primary. No subscription savings have been measured.

Converted qualification passed: 365 client / 362 server roots, safe/static,
three production/ABI builds, exact known-suite comparison, both fresh headless
save/load runs and original asset hashes. Audio production evidence is reused only
for the original test-only baseline, whose production source is identical.


## Qualified result

All 365 client / 362 server roots pass: the original exact names plus the two new
native tooltip/event contracts. Safe/static, three production/ABI builds, exact
known-suite comparison, both fresh headless save/load scenarios and all 1,654 asset
hashes pass on unchanged reviewed source. See [qualification](ui-meter-identities-qualification.json).
Exports fall 262→242; selected production cgo files fall 146→142 client and 147→143
server. Headers remain 157 files / 3,002 physical lines; embedded production C
bodies remain 77; standalone production/test-reference C remains zero.
