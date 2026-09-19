# Player-file sections — C baseline in progress

Qualified parent: **d7f52707**, client audio events and playback. Production C
remains **21,083 physical lines /65 files /zero reference C**.

The batch selects 18 functions /1,540 body lines in GAME1_1.c: player attributes,
status, inventory, field guide/spellbook/enchantments, journal, game and GUI data,
file metadata, music state, section framing and client player-save requests.
The adjacent briefing gate helpers are unrelated and excluded. Whole-repository
review finds 16 external roots and all 18 bodies reachable. The local C section
table is also owned by this batch; distinguish its internal callback references
from interfaces needed by outside C callers during translation.

Reuse real player/object/inventory, journal, quickbar, music, timer and cryptfile
owners. Verify exact section bytes, historical version gates, read/write mutations,
flags, optional fields, Unicode lengths, file positions, callback order and cleanup.
File metadata uses the wall clock: identify and validate that field explicitly
before normalizing any writer capture. Use deterministic input timestamps on reads.

Production source is currently unchanged. Reuse the qualified parent production
baseline only after checking identity; corrections would require fresh qualification.
No expectations are frozen and no conversion has started. The actual-C dispatcher
and focused manifest are installed. Remaining work is listed below.

The first music C sweep passes both roots: 24 write and 112 read records cover
exact state-word order, current and historical signed version gates, optional
section flags, queue levels, counts through the adjacent-slot boundary, untouched
backing bytes, file position and real music module changes. `music-initial` is
joined. Player attribute contracts pass (25 write /54 read records). They cover
field order, UTF-16 name lengths, supported/historical versions, name rejection,
mode flags, legacy skipped words and guards with an actual PlayerInfo allocation.
The attributes draft is consumed; use tracked source. No expectations are frozen.

Book contracts add exact spell/ability write bytes for all three classes, sparse
and full entry sets, flag gates, and read version/count/presence boundaries for
both books. Two discovery attempts found incorrect fixture imports; corrected
before books-third. Production remains unchanged.

books-third (six roots), fieldbook-initial (seven) and status-initial (eight)
passed and are joined. Field-guide names use explicitly owned lookup strings,
including empty, UTF-8 and 255-byte names. Status writers cover unsigned word
boundaries, current/max health and mana, saved temporary words, poison fields
and object guards. Their drafts are consumed. Spellbook restoration now tests
real award calls, bookkeeping records and reliable report queues; all 90 cases
pass. No frozen baseline or conversion yet.

Journal/inventory contracts pass (journal-inventory-second, 13 roots), including
existing journal save/load coverage. The first attempt used a multiline pattern;
the driver requires one regex line, corrected before running tests.
Guide-family restoration initially failed because the fixture passed an initial
value to alloc.New, whose argument supplies only the type. Explicitly initializing
the allocated array fixed it; no production change. Metadata reads and enclosing
GUI section contracts pass together in gui-guide-third (17 roots).

Metadata writes and recovery-repeat both pass all 18 roots. Seventeen new captures
/1,006 records repeat byte-for-byte; the eighteenth root reuses the existing
journal save/load golden. Timestamp normalization follows actual clock bounds and
weekday checks. The static scanner required literal 10980 instead of a local
constant in two metadata accessors; static-recovery now passes. All runs are joined.
This is a partial-baseline recovery checkpoint, not a frozen or qualified conversion.
The tracked recovery JSON records capture hashes and counts; expectations in new
tests remain unfrozen. Production source is identical to d7f52707: only porttest
files and documentation have been added/changed. Full baseline review remains.

Remaining before freezing: status reads / real-unit attributes, complete inventory
read/write, enchantment state and historical ability tails, game/script sections,
outer file framing/extraction and save request, remaining quest-mode book gates.
Review coverage and real constant/table owners, repeat captures in fresh processes,
qualify affected targets, verify production identity, then commit frozen baseline.
All installed drafts (attributes, fieldbook, status, journal gates, inventory helpers,
guide restore, metadata read/write, GUI) are consumed; use source, never replay.
