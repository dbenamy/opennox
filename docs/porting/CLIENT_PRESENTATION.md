# Client world, spell and item presentation — C baseline in progress

Qualified parent: **15836cdb**, native combat overlays. Production C remains
**24,411 physical lines /66 files /zero reference**.

Scope:22 bodies /715 C body lines across GAME2_2 and GAME2_3. The caller audit
includes the main drawable frame pipeline: wall/drawable ordering, cached floor
sprite composition and its private copy helper, phoneme icons and playback,
reflective shields, Turn Undead particles, ray ownership, book rewards, bubble
particles and equipment appearance. Four initial candidates turned out to be
world-rendering helpers rather than spell helpers. Keep them in this connected
presentation batch: they are live in the same frame pipeline and established
floor/drawable fixtures already own their inputs.

Every candidate has a live external caller or callback/private-helper path.
See client-presentation-selection.json. No reachability pruning is proposed.

Reuse real tile raster/composition, drawable/render, GUI, NPC and spell owners.
Contract priorities: relative ordering and side boundaries; exact copy spans and
scroll/wrap/cache behavior; frame wrap and phoneme timing; partial creation and
cleanup; ray capacity/duplicates; equipment capacity and modifier identity;
particle colors, RNG and callback state; book reward throttling and placement.
Audit each original lookup table before capturing results. Compare real pixels
and normalized state, with independent assertions before freezing repeated C.

The parent production qualification may be reused only while production source
is identical. Any pre-baseline production correction requires fresh qualification.
Seven original-C contract groups pass (4,619 records): ordering, copy, floor
composition, equipment, phoneme loading and phoneme rendering/expiry. Fixture-only
corrections supplied a real GUI, explicitly typed an unsigned constant for 386,
and converted a named image-handle pointer before comparison. Production remains
unchanged. Turn Undead contracts are being checked next. Original table/name bytes
are supplied explicitly; phoneme pointers match blob_init.go relocation. No golden
is frozen. Copied bake and selection drafts are consumed; tracked source wins.

Pre-baseline correction under verification: the persistent-ray allocator could
create a 97th sprite after its 96 pointer slots were full. No writer updates the
old capacity counter. An independent repeated-creation test reproduced98 ->99
live sprites (including two endpoints); the final sprite was untracked. Move the
existing slot search ahead of allocation. The sole production caller ignores its
return. Duplicate removal and free-slot reuse remain covered. Fresh production
qualification is now required; the preceding production baseline cannot be reused.
See presentation-rays-fields/tests.jsonl for the original failure. Current C is
24,412 lines (+1 explanatory comment); conversion has not started.

## Frozen recovery checkpoint

16 captures /7,744 records match default, server and high-resolution repeats with
identical source fingerprints. Hashes are now frozen in the focused tests and
client-presentation-captures.json. All16 independent contracts pass, including
the corrected capacity case. The affected265-root corpus and fresh production
qualification are still pending; this checkpoint is not a completed conversion.
All copied fixture drafts and freeze.py are consumed. The native ray implementation
draft is uninstalled and requires review. Current selected C is716 body lines.

Completed native combat scenario copies were deduplicated after content-hash and
inactive-process checks, reclaiming1,112,747,701 bytes. The restoration helper is
build/port-combat-overlays/deduplicate-combat-overlays-native-assets.py; only its
--restore NAME mode should be reused. Original assets/archive remain untouched.
