# Porting decisions for later review

## Working policy (user instruction, 2026-09-12)

When reasonably confident in the right answer and reversal would not require
substantial effort, make the decision and record it here or in the batch doc
for later review. Continue without asking for confirmation. Record behavior
changes explicitly, with evidence and validation; do not label them exact
compatibility. Ask when confidence is insufficient or reversal would be costly.

## Locked-door message padding — approved, review with Go conversion

Initialize unused bytes in the fixed 52-byte notification. The client interprets
only the header, NUL-terminated localization key and selector. Original C left
padding uninitialized; a zero-initialized C array now establishes deterministic
output for the upcoming Go conversion. No meaningful fields or string acceptance
rules change. Reversal is one initializer, though preserving indeterminate
padding is not recommended. The user explicitly approved this correction.
Evidence and reproducible tests: [PLAYER_CONTROLS.md](PLAYER_CONTROLS.md).

## Default-equipment modifier word — review with controls conversion

Initialize the entire 20-byte C modifier array in 4EF7D0. Only four descriptor
words were assigned, but the attribute helper copies a fifth word into items.
Chosen under the standing policy: use deterministic zero for that uninitialized
word, matching the future Go array default. Reversal is one initializer. Retain
all descriptor selection, ordering, callbacks and defined values. Evidence and
validation are recorded in [PLAYER_CONTROLS.md](PLAYER_CONTROLS.md).

## Spell phoneme class offset — review with spell lifecycle conversion

Correct the server phoneme helper's class-byte access before locking its C
baseline. `getObjectFromNetCode` returns `nox_object_t*`; adding 8 before the cast
advanced eight whole objects (6,176 bytes), instead of reading byte offset 8.
Nonplayer phoneme cases exposed an out-of-bounds read and crash. Cast to a byte
pointer before adding the offset, matching the client branch's class-field check.
The new corpus checks male/female and nonplayer sounds in both server and client
paths. This is a deliberate bug correction, not exact preservation of the invalid
read. Chosen under the standing policy; reversal is one expression.
