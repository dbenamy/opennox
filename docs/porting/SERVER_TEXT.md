# Extension and server-listing helpers

Status: native conversion qualified. C falls from the corrected baseline's 1,760
lines /16 files to **1,219 /12**, removing 541 lines (521 net since `eb01acc9`).
No reference C remains. Corrected C baseline: `d87a3e02`.

## Native qualification

The first native probe matched all five captures. Default/server/highres pass
173/173/173 affected roots, no skips and all 21,135 frozen cases. The selection
covers new helpers and accumulated equipment, inventory, controls, server config,
browser, team, reliable reports, server messages and quest runtime contracts.
Three production builds and ABI/retired-symbol checks pass; the full suite matches
known failures, and fresh headless gameplay plus explicit save/load pass. Source
fingerprints agree, and preflight/final default binaries have the same SHA-256.
See [the native record](server-text-native-qualification.json).

Names preserve exact low-byte UTF-16 matching and output write order. Weapon
cycling uses the native equipment owner and stops at the first eligible candidate,
including failure. Trap drop preserves first-match success independent of the
drop service result. The listing assembles a temporary header before copying,
keeping overlapping input/output tokens valid. Flag indices outside 0..31 now
explicitly return zero; the UI only passes 1..5. Values below -31 had no defined
C result; this deterministic extension is recorded for review.

Fifteen C interfaces are removed. Twenty-two storage definitions retain their
exact types/initializers in vardefs.c; seven remaining callback identities retain
their signatures and bodies. The memory-file ABI struct remains; its eight dead
C methods and prototypes are removed. No old algorithm remains solely for tests.
All local installers/finalizers/native drafts are consumed or stale.

## Scope and reachability

- GameEx player-name lookup, weapon cycling, trap drop and settings flag index.
  The name paths are live through `EXTENSION_MESSAGES` config and settings UI;
  the default extension mask alone is not a reachability argument.
- Server-listing formatter and its Go slice boundary.
- Remove eight orphan C memory-file readers/skippers. Literal references are
  definitions, declarations and old source-conversion/callgraph tool mappings.
  Preserve the live memory-file ABI struct and Go owners.
- Consolidate remaining declaration-only C translation units without changing
  variable types, initial values or addresses. Keep registered no-op callbacks
  whose function addresses are identities.
- File enumeration and varargs/string formatting remain a following batch: both
  have live callers and require their own compatibility contracts.

## Intentional corrections before freezing C

The player-name parsers used a one-byte local destination for an entire converted
name. Use the actual 28-unit `nox_playerInfo.name_final` width plus a terminator.
Preserve historical low-byte conversion of each UTF-16 code unit, including
embedded low-byte NUL behavior; this is not a change to UTF-8. Reject nil records,
the existing -2 sentinel, nil output pointers and missing team lookup results
without partial writes.

Initialize all 72 listing-header bytes. Reserved bytes 22/23 and 70/71, and the
non-quest stage bytes 68/69, now deterministically contain zero. The wire map
field is nine bytes including its terminator, confirmed against the installed
`noxnet/discover.MsgServerInfo` decoder. Limit its name to eight bytes. Keep the
existing overlapping writes: slot byte 47 is overwritten by the first map byte.
The public Go wrapper rejects source lengths below 12 and destinations smaller
than 73 plus the server-name length. Remove its unused C input copy.

These are reversible corrections authorized by the port process, recorded for
review. They change production C, so the preceding production qualification
cannot be reused as identical source for this baseline.

## Contracts

- Names: real player records and object lookup; duplicate names, first active
  match, exact case/byte matching, empty/non-ASCII/embedded-NUL/full-width names,
  missing objects, all active subsets, class/team byte boundaries, output-address
  identity and untouched record/output guards.
- Listing: 12,544 cases cover hosting/block gates, dedicated counts and narrowing,
  quest flags/stage, settings/slot fields, video bits, map/name boundaries,
  input token, short/exact/large output buffers and unchanged input/guards. The
  independent packet decoder checks the resulting name, map and token.
- Flag index: -31 through 31; production UI passes only 1 through 5. Inputs below
  -31 cause C division by zero and are outside the frozen defined-input contract.
- Weapon cycling: 3,920 cases over status/state gates, both directions and signed
  nonzero direction bytes, no/current/end weapons, zero-mask/ammunition/class/strength
  skips, actual equip/dequip failure, first eligible failure and selected weapon.
- Trap drop: nil owners, status/state gates, all subsets and forward/reverse/empty
  inventory order, exact class-byte matches, callback failure, first match and
  nearby/reach-limited cursor placement. Real drop dispatch and callbacks remain.

Inventory results include the entire existing normalized owner/effect snapshot.
The capture stores a SHA-256 of the JSON for each result to avoid duplicating
hundreds of MB of unchanged fixture memory. No result fields are omitted before
hashing; independent assertions run first. Other new captures retain full rows.

## Probe issues and local recovery

Probe1: listing fixture forgot the real setter's 15-byte name limit/terminator;
corrected the fixture and isolated raw game flags from lifecycle hooks.
Probe2: listing and flag contracts passed; the missing-object name fixture only
removed the world list, but actual lookup also searches active players. Give the
record an absent net code for that case.
Probe3: weapon cases passed; trap position assertion missed the real 75-unit
reach limit. Use exact-axis near/far cases with independent expected positions.
Probe4 passed all five groups /20,575 cases. Probe5 adds the zero-weapon-mask
boundary, bringing weapon cases to 3,920. All five groups /21,135 cases pass.
The captures are frozen in `server-text-c-batch.json`; both default processes, server, highres and fresh C gameplay preflight pass.
All five roots run without skips and every frozen hash matches. See
[the baseline record](server-text-c-qualification.json). Three-profile production/
ABI, exact full-suite comparison and explicit save/load remain required after
conversion; this is a recoverable baseline, not completed native qualification. Native installation followed this committed baseline. The draft files are now
stale; tracked source is authoritative.

Local artifacts: `build/port-final-text-helpers`. Probe3's completed large capture
is gzip archived with round-trip SHA-256 verification (`c-probe3/archive.json`).
The older runtime save-run duplicate assets were verified against original
assets and deduplicated: 556,358,986 bytes. Restore with
`python3 build/port-final-text-helpers/deduplicate-runtime-save.py --restore server-runtime-native-save`.
The cleanup script's audit/apply are consumed; preserve original assets.

Eight obsolete root-test archives without the new extension roots were removed
after checking age, archive symbols, hashes and open files: 652,232,214 bytes.
Production/library caches and all current extension-test archives were retained.

Native-preflight asset copies were also verified/deduplicated, 556,388,715 bytes;
restore with `python3 build/port-final-text-helpers/deduplicate-text-native.py --restore server-text-native`.
C-preflight copies use `deduplicate-text-c.py --restore server-text-c` in that
same directory. All cleanup jobs are joined and applies consumed.
