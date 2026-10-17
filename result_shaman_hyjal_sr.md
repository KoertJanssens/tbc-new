# Hyjal soft-reserve plan — Koerrii (Enhancement Shaman)

Optimising the 3 soft reserves for <https://softres.it/raid/NbW6qYWa> (Hyjal Summit, 25-man,
Horde, TBC), using measured DPS from this repo's engine rather than item level.

**Baseline: 2237.6 DPS** on Koerrii's current armory gear — the same baseline as
[result_shaman.md](result_shaman.md). Every Hyjal drop an Enhancement Shaman can equip was simmed
as a single-slot swap at 20,000 iterations (`sim/shamanhyjalsr_test.go`).

## Status: locked — data refreshed 2026-09-01

**This list is now locked at 25/25 reserves and Koerrii is on it.** Nothing below is actionable any
more; it is kept as the record of how the picks were made.

Reserved: **Valestalker Girdle**, **Bow-stitched Leggings**, **Midnight Chestguard**.

That takes the one pick that mattered — Valestalker Girdle, still the only reserve on it — and
trades the other two for higher raw DPS at worse odds than recommended:

| Reserved | Slot | DPS | Rivals | P(drop) | P(win) | Value |
|----------|------|-----|--------|---------|--------|-------|
| Valestalker Girdle | Waist | +29.0 | **0** | 0.33 | 1.00 | **6.44** |
| Bow-stitched Leggings | Legs | +24.1 | 3 | 0.44 | 0.25 | 2.68 |
| Midnight Chestguard | Chest | +27.3 | 2 | 0.27 | 0.33 | 2.43 |
| | | | | | | **11.55** |

Against the recommended set (Valestalker Girdle + Deadly Cuffs + Shady Dealer's Pantaloons) at
**13.48**, that is a difference of **~1.9 DPS per clear** — real but small, and the three reserves
sit in three different slots so none of the value is wasted on overlap. The decisive call, taking
the uncontested girdle, was made correctly.

Four raiders joined after the original analysis: Boltmacka (Destro Lock), Samaa (Shadow Priest),
Wfr (BM Hunter) and Koerrii. Only one touched this shaman's targets — **Wfr took Don Alejandro's
Money Belt**, pushing it to 5 reserves and confirming it was right to avoid.

## Recommendation (as made, before the list locked)

| # | Reserve | Boss | Slot | DPS | Rivals | Why |
|---|---------|------|------|-----|--------|-----|
| 1 | **Valestalker Girdle** | Kaz'rogal | Waist | **+29.0** | **0** | Biggest Hyjal upgrade you can get *and* nobody has reserved it |
| 2 | **Deadly Cuffs** | Rage Winterchill | Wrist | **+27.6** | 2 | Second-biggest, and it is protected *away* from you unless you reserve it |
| 3 | **Shady Dealer's Pantaloons** | Azgalor | Legs | **+13.4** | 0 | Smaller, but Azgalor has the shortest loot table in the zone — it drops often |

Expected value of that set: **~13.5 DPS per full clear**, the best of any 3-item combination.

The three runners-up are all within ~0.7 of it, so this is not a knife-edge decision — but do not
drop Valestalker Girdle under any circumstances. It is worth more than the other two combined.

## Every Hyjal drop, ranked

Single-slot swap against your current gear. Weapons are all large downgrades and are listed for
completeness — your Merciless Gladiator's Cleavers are better than anything Hyjal offers.

| Item | Boss | Slot | Armor | DPS |
|------|------|------|-------|-----|
| Valestalker Girdle | Kaz'rogal | Waist | Mail | **+29.0** |
| Deadly Cuffs | Rage Winterchill | Wrist | Leather | **+27.6** |
| Midnight Chestguard | Archimonde | Chest | Leather | **+27.3** |
| Don Alejandro's Money Belt | Anetheron | Waist | Leather | **+24.9** |
| Bow-stitched Leggings | Azgalor | Legs | Mail | **+24.1** |
| Shady Dealer's Pantaloons | Azgalor | Legs | Leather | **+13.4** |
| Bracers of the Pathfinder | Rage Winterchill | Wrist | Mail | +11.4 |
| Razorfury Mantle | Kaz'rogal | Shoulder | Leather | +10.9 |
| Quickstrider Moccasins | Anetheron | Feet | Mail | +7.1 |
| Choker of Serrated Blades | Trash | Neck | — | +6.6 |
| Mail of Fevered Pursuit | Archimonde | Chest | Mail | +3.9 |
| Black Featherlight Boots | Kaz'rogal | Feet | Leather | -1.5 |
| Claw of Molten Fury | Trash | Main hand | — | -42.3 |
| Chestguard of Relentless Storms | Trash | Chest | Mail | -65.2 |
| Boundless Agony | Azgalor | Off hand | — | -94.5 |
| Tracker's Blade | Rage Winterchill | Off hand | — | -115.7 |
| Boundless Agony | Azgalor | Main hand | — | -132.1 |
| Tracker's Blade | Rage Winterchill | Main hand | — | -144.8 |
| Hammer of Atonement | Kaz'rogal | Main hand | — | -287.1 |

## The reserve maths

An SR is only worth what it adds *over not reserving*. Those are two different situations:

- **Someone else has already reserved the item.** Then it is protected and you cannot get it at
  all without reserving it too. Your reserve is worth the full `1/(rivals+1)` roll.
- **Nobody has reserved it.** Then it goes to open roll if you leave it, and you would often win
  anyway. Your reserve only buys the difference.

So `value = P(item drops) × [P(win with SR) − P(win without SR)] × DPS gain`.

| Item | Gain | Rivals | P(drop) | P(win) no SR | P(win) with SR | **Value of reserving** |
|------|------|--------|---------|--------------|----------------|------------------------|
| Valestalker Girdle | +29.0 | 0 | 0.33 | 0.33 | 1.00 | **6.44** |
| Shady Dealer's Pantaloons | +13.4 | 0 | 0.44 | 0.33 | 1.00 | **3.97** |
| Deadly Cuffs | +27.6 | 2 | 0.33 | 0.00 | 0.33 | **3.07** |
| Bow-stitched Leggings | +24.1 | 3 | 0.44 | 0.00 | 0.25 | 2.68 |
| Bracers of the Pathfinder | +11.4 | 0 | 0.33 | 0.33 | 1.00 | 2.53 |
| Midnight Chestguard | +27.3 | 2 | 0.27 | 0.00 | 0.33 | 2.43 |
| Razorfury Mantle | +10.9 | 0 | 0.33 | 0.33 | 1.00 | 2.42 |
| Don Alejandro's Money Belt | +24.9 | 5 | 0.33 | 0.00 | 0.17 | 0.92 |
| Quickstrider Moccasins | +7.1 | 0 | 0.33 | 0.33 | 1.00 | 1.58 |
| Mail of Fevered Pursuit | +3.9 | 0 | 0.27 | 0.33 | 1.00 | 0.69 |

Best three-item combinations, accounting for the fact that two items in the same slot do not
stack (winning one makes the other nearly worthless):

| Combination | Value |
|-------------|-------|
| **Valestalker Girdle + Deadly Cuffs + Shady Dealer's Pantaloons** | **13.48** |
| Valestalker Girdle + Shady Dealer's Pantaloons + Bracers of the Pathfinder | 12.95 |
| Valestalker Girdle + Midnight Chestguard + Shady Dealer's Pantaloons | 12.84 |
| Valestalker Girdle + Shady Dealer's Pantaloons + Razorfury Mantle | 12.84 |
| Valestalker Girdle + Bow-stitched Leggings + Shady Dealer's Pantaloons | 12.65 |

## Traps to avoid

- **Don't reserve Don Alejandro's Money Belt.** +24.9 looks strong, but it has the most competition
  in the raid (5 reserves after the refresh) *and* it is a waist item — if you win Valestalker Girdle it is worth
  almost nothing to you. Two waist reserves is the single worst use of your three slots.
- **Don't reserve a tier 6 token.** The Hyjal tokens are Gloves and Helm of the Forgotten
  Protector. Skyshatter Cover sims at **+7.5** and Skyshatter Grips at **+4.4** for you (see
  [result_shaman.md](result_shaman.md)) — a quarter of what the girdle is worth — and taking one
  drops you from 5/5 to 4/5 Cataclysm Harness for no gain.
- **Midnight Chestguard is the worse of the two big contested items**, despite being nearly
  identical in DPS to Deadly Cuffs (+27.3 vs +27.6). Archimonde is the last and hardest boss, and
  has a 15-item table against Rage Winterchill's 12 — so it is both less likely to be killed and
  less likely to drop. Reserve the cuffs.
- **Don't reserve weapons.** Every Hyjal one-hander is a downgrade on your Merciless Gladiator's
  Cleavers, several of them catastrophically.

## Who you are competing with

**Sliumpius is another Enhancement Shaman** and has reserved Bow-stitched Leggings, Midnight
Chestguard and Deadly Cuffs — three of your top five, and now a direct rival on two of the three
you locked in. They did **not** take Valestalker Girdle, which is why it was such a clear first pick
and why it is still uncontested on the final list.

Srhóx (Fury Warrior) also holds Deadly Cuffs, Midnight Chestguard and Don Alejandro's Money Belt.
Cuica (Ret Paladin) has Bow-stitched Leggings twice.

## Assumptions

- **Drop rate**: TBC 25-man bosses drop 4 items. Table sizes come from `assets/database/db.json`
  (zone 3606) plus the tier tokens the database does not carry — 3 glove tokens on Azgalor, 3 helm
  tokens on Archimonde. That gives 4/12 for Rage Winterchill, Anetheron and Kaz'rogal, 4/9 for
  Azgalor, 4/15 for Archimonde. Azgalor's short table is why his legs are a better bet than their
  raw DPS suggests.
- **Open-roll rivals**: for an unreserved item I assume 2 other plausible rollers — Sliumpius
  (Enhancement Shaman) and Hanebuth (BM Hunter) for mail, Snipemaster (Combat Rogue) and Ertanias
  (Feral Druid) for leather. If your guild hands leather to rogues and druids first, Shady Dealer's
  Pantaloons is worth less than modelled and Bracers of the Pathfinder (mail) becomes the better
  third pick.
- **Reserves refreshed 2026-09-01**, list locked at 25/25. Rival counts in the tables below exclude
  Koerrii, so they read as "who else wants this".
- **Single target, patchwerk**: same encounter assumptions as [result_shaman.md](result_shaman.md).
- Choker of Serrated Blades is a trash drop, so it has no meaningful drop rate to model and is
  excluded from the reserve maths. At +6.6 it would not make the cut anyway.

## Two engine issues found while building this

1. **Beast-tamer's Shoulders (30892) panics the sim for any non-hunter.**
   `sim/hunter/item_sets.go:243` registers its pet effect with an unconditional
   `agent.(hunter.HunterAgent)` cast, and the item carries no `classAllowlist`, so equipping it on
   a shaman crashes the engine rather than ignoring the effect. Omitted from the run; its stats
   (agility/AP, no crit or hit) lose to Razorfury Mantle for a shaman anyway.
2. **The engine does not enforce main-hand-only weapons in the off hand.** Hammer of Atonement and
   Claw of Molten Fury are `handType: 1` in the database, but equipping them off-hand produces a
   result rather than an error — an apparent -880 and -682 DPS, which is really the character
   swinging with a broken off hand. The harness now skips those combinations; they are excluded
   from the ranked table above.
