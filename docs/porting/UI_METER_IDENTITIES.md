# UI meter bridge baseline and next conversion

## Scope under review

The immediate candidate is the 17 remaining meter exports in the state, items
and draw files, together with two meter constructors and the adjacent inventory
weapon-draw adapter. These have Go owners and internal GUI/fixture consumers.
The helper's initial 17-export recommendation retained those adjacent C bridges;
primary review favors including them if the complete caller/fixture audit confirms
that they can be removed together. No UI production conversion is installed yet.

Production source is qualified audio revision `9eac418e`. Two test-only files add
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

## Proposed implementation boundaries

Use existing typed `gui.WindowFunc` and `gui.WindowDrawFunc` support for meter
events/drawing; a broad change to the generic C wrappers is unnecessary. Preserve
nil-versus-nonzero event response conventions, 386 argument/result widths, callback
capture identities, owner algorithms, rendering order and mapped record layouts.

Tooltip dispatch needs a native route while retaining the foreign callback path.
Primary favors a small callback-identity registry consulted by `Window.TooltipFunc`
after its existing nil/dead checks; this preserves the stored callback word and
normalization while avoiding changes to the C-compatible Window layout. The
helper instead suggested an extension-field setter. Finalize this reversible
implementation choice during the complete caller audit. Add native registry
contracts when that code exists; the new original-path contracts protect the
existing installed route and foreign fallback.

The inventory weapon draw edge has a C-typed wrapper with fixture consumers.
Either migrate all consumers in the same batch or extract the typed owner while
retaining that bridge; do not leave a native key entering raw C dispatch. The
constructor's two adjacent C exports have only inventoried Go/fixture callers;
check the full reference graph before removing their prototypes.

## Resume checkpoint

Luna reached its usage limit during the read-only scout. Its partial proposal and
reference inventory are under `build/port-ui-meter-bridge-scout-20260925/`; its
reported Git head predates the audio commit, so verify actual recorded file hashes.
It did not draft or install a UI conversion. Primary completed and accepted the original baseline, proving production source
equality to `9eac418e`. This test/documentation checkpoint is ready for resumption;
no conversion or converted qualification is claimed.

Requalify converted source with the selected profiles, safe/static, three
production/ABI builds, exact known-suite comparison, both fresh headless save/load
runs and original asset hashes. Reuse audio production evidence only for the
original test-only baseline, after proving production source equality.
