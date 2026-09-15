# Inventory drag cancellation ownership

Status: **ownership cleanup qualified**.
The preceding window conversion is `24f3f67e`. C LOC remains **82,632 / 96 files /
zero reference C**.

The original cancellation code restores/copies an inventory drag into a separate
cell drawable, then clears the temporary drag pointer without deleting it. Normal
mouse release already deletes this copy. Equipment drags borrow an existing
drawable and must never delete it. The conversion deliberately preserved this
behavior; its frozen captures and committed C baseline remain recoverable.

Independent contracts use the real drawable factory, pool counts and deletion
ledger. They cover 600 repeated cycles across stack sizes 1/2/32 and direct-cancel
or window-close paths; 200 borrowed-equipment cycles; fallback into another cell;
and full-inventory restoration failure. They also check item identities, cleared
cursor/mouse capture, exactly-once deletion and idempotent cancellation.

The first discovery attempt failed on an unused draft import, before executing
tests. After correction, the unchanged implementation completed all three root
tests in 21.131s: borrowed equipment passed both 100-cycle paths; all eight
non-equipped subcases failed the retained-temporary assertion. All jobs joined
before source edits.

The fix deletes the temporary after the restoration attempt, guarded by the
existing equipment-borrow flag. It leaves the restored inventory drawable and
borrowed equipment alive. Full-inventory failure still reports the existing error
and clears the drag; its temporary allocation is now also released.

The first fixed run completed all 55 window/display root tests in 180.622s. All
new ownership contracts passed, including all 800 repeated cycles. Two expected
capture mismatches remained: drag-drop and cancel-fallback. No other tests failed.

## Exact expectation review

All 29 capture groups were compared with the qualified conversion's raw captures.
Exactly six results changed: drag-drop rows 9/21/33/45 and cancel-fallback rows 1/3.
For each result, an independently constructed expected correction:

- Marks only the temporary object's ledger entry dead and appends its one deletion.
- Removes only that drawable from the renderer's live-object snapshot and appends
  the same identity to the renderer's deletion list.
- Decrements drawable count and live pool count by one and increments free pool
  count by one.

Applying just those changes to the old captures yields the new captures exactly.
Every other field—including pixels, GUI state, item codes, equipment flags, sounds
and requests—is unchanged. The remaining 27 groups are byte-identical. Two frozen
hashes were updated only after that verification; the C baseline and old hashes
remain recoverable. The shared failure message now says “frozen expectation” to
cover both original C results and this explicitly reviewed correction.

Evidence is under build/port-client-inventory-cancel: before-b-result.json,
after-a-result.json, capture-deltas.json, verify-reviewed-deltas.py and
reviewed-hashes.json. The focused repeat passed all 55 window/display root tests in 57.969s, with
all 29 repeated capture hashes matching the reviewed expectations.

| Group | Original C hash | Reviewed ownership hash |
|---|---|---|
| cancel-fallback | `83117e86628d2a3dfcf6091ab15579236d6ebf1403da861d4950f7ffaf87e58d` | `02249a85ac4fd7bc542dc340eb45ad29b48024bb2dc44e0a4bb0cad6da05a56f` |
| drag-drop | `41ca0fa4c4a5a72d1b8b4f4a795d10c396922d940873da843adbd6a4768012f3` | `ada8d7b8d6e2499ff5a3dac275cd601b7e0de83b9ff952e6683ad4d450f9a9d2` |

## Final qualification

Final checks passed: accumulated default **877 / 469.897s**,
affected server **314 / 229.749s**, highres
**316 / 155.701s**. All selected tests started and
finished. All three production binaries verify ELF32/i386/SSE2/CGO, the existing
11 retained / 23 retired window interfaces, and no test helpers. Full assets retain
exactly 1,553 known failures and 15 pass / 3 fail / 32 skip packages. The nine-screen
headless scenario matches the original C gameplay reference in
**50.395s**, with reference replacement disabled. All **1,581 source
fingerprints** remain unchanged. C LOC is **82,632 / 96 files / zero reference C**.

Final evidence: native-qualification.json, native-{default,server,highres}-result.json,
native-source-fingerprints.json, native-binary-verification.json and
full-suite-comparison.json under build/port-client-inventory-cancel. Gameplay is
build/baseline/runs/client-inventory-cancel-port. Completed captures have verified
lossless compression/restoration manifests; active data and original assets were
left unchanged. All source-reading jobs were joined before committing.
