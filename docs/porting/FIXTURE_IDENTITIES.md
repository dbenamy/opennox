# Independent fixture identity scopes

Status: qualified across default, server and high-resolution profiles.
Production code and frozen capture expectations are unchanged.

## Failure and investigation

Remaining-storage qualification stopped in server consumer group seven at
`TestPlayerControlsStats`, capture `controls-stats-12`. The expected digest was
`1eb30e0b28bfd067925f8232ea8b451b87157747abe0cec2d972f14bb46901f9`;
the observed digest was
`664b79cbc45503d5444a96d740bf9a8e69dc3a8bb3a4369302930862a45b7101`.
The group completed all 200 selected roots without skips. Earlier stats operations
matched. Actual callback JSON was not enabled for that run, limiting diagnosis.

The unchanged source subsequently passed three isolated stats repetitions,
twenty independent stats processes, and two complete repetitions of the same
200-test group with the exact qualification asset environment. A separate probe
without the asset environment passed stats but skipped a shipped-definition test;
it is diagnostic evidence only. No expected hashes were regenerated.

The shared callback helper now accepts `OPENNOX_CALLBACK_FAILURE_CAPTURE`: on a
mismatch it saves the actual JSON before reporting the original failure. Passing
tests need not dump large captures. Qualification enables this option.

## Concrete fixture defect

`PortTestRoam` shares an address-to-ID map across a corpus of independent cases.
The local normalization closure and `proxy.life.ids` alias the same map.
`portTestShopPools.identify` adds transient allocation addresses to both the local
shop map and that shared map. Objective regions and player resource records are
released after their case, but only the shop-local map was cleared. The shared
map retained addresses from released allocations.

Whole-word snapshot normalization could therefore interpret a later scalar word
as an address from an earlier case. This lifetime defect is confirmed independently
of the intermittent digest: a case-entry assertion against the initial persistent
identity map fails deterministically on the second stats case before restoration.
The red source and log are retained under `build/port-remaining-storage/`.

The repair clones the persistent identity map after corpus setup, then restores
it in place after every complete case and all its snapshots. Updated aliases for
original persistent addresses are retained; new transient keys are removed. In-place restoration
preserves aliases held by the closure and lifecycle fixture. Clearing in shop
cleanup would be too early because outer snapshots still use current-case IDs.
An entry assertion retains the independent-case invariant. Per-case setup registers
its current callback and allocation identities as before.

## Limits and review

The old digest mismatch is not proven to have arisen from this defect. Restoring
map lifetime also does not eliminate the separate possibility that a *live*
pointer word equals an ordinary scalar. Typed pointer-field normalization would
be a separate fixture change if concrete evidence warrants it. Preserve failure
captures and investigate any recurrence; do not normalize away mismatched values
or weaken the frozen checks.

GPT-6 Luna performed the bounded read-only lifetime review and identified the
shared map as a plausible source of nondeterminism. The primary verified the
allocation/free and alias paths, built the deterministic red, chose the repair,
and owns qualification. No subscription savings are inferred.

The first broad probe caught an overbroad restoration: seven captures changed
only because the live player alias 54000 reverted to its initial ID 200. Complete
JSON comparisons found no other differing values. The repair now retains alias
updates for the original persistent address set. It still drops all new transient
keys. Frozen expected captures remain unchanged; full qualification is pending.

## Qualified result

The final fixture passes all 2,291/2,280/2,291 selected consumer roots across
all three profiles, without skips or changed expected hashes. Both storage
captures and the audio GC regression also pass. Fresh safe build, production/ABI,
known-suite comparison, headless gameplay and explicit save/load pass with the
same final source fingerprints. See [raw-storage-native-qualification.json](raw-storage-native-qualification.json)
for the accepted evidence and the recorded investigation disposition.
