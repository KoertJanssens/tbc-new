# Phase 4: Glaives vs Non-Glaives — Fury Warrior vs Rogue

20,000-iteration comparison of the legendary Warglaives of Azzinoth against a non-Glaive
weapon substitute, for Fury Warrior and Combat Rogue on their **Phase 4** BiS gear, across
1/2/4-target encounters. Unlike the Phase 3 trinket benchmark (see [results.md](results.md)),
trinkets are not cycled here — each spec keeps the fixed trinkets from its own P4 gear preset.
A third option, **Sword Spec Fury**, is also included for Warrior Fury: same Warglaives, but
with 5 points respec'd from Improved Berserker Stance into Sword Specialization (legal here
since Warglaives of Azzinoth are itemized as Swords in this repo's item database).

## Assumptions

- **Iterations / seed**: 20,000 iterations, `RandomSeed: 101`, `core.RunRaidSimConcurrent`,
  for every spec/weapon/target-count combination.
- **Encounters**: `core.MakeSingleTargetEncounter(5)` (180s ±5s), extended to 2 and 4
  `core.NewDefaultTarget()` dummies for the cleave scenarios — same encounter construction
  as the trinket benchmark's cleave/4-target scenarios in results.md.
- **Buffs/debuffs**: `FullRaidBuffs` / `FullPartyBuffs` / `FullIndividualBuffs` /
  `FullDebuffs`. No bonus hit rating is added this time (unlike the trinket benchmark) —
  each P4 BiS preset's own gear hit rating is used as-is, since this is a same-spec
  weapon-swap comparison rather than a cross-spec one.
- **Gear**: each spec's existing/newly-built Phase 4 BiS preset, with only the
  MainHand/OffHand item IDs swapped between the two weapon variants (existing enchants on
  those slots are preserved for both variants):
  - Warrior Fury: `ui/warrior/dps/gear_sets/p4_fury.gear.json` (already existed; equips
    Warglaives by default, Trinket1/2 = Madness of the Betrayer + Dragonspine Trophy).
  - Rogue: `ui/rogue/dps/gear_sets/p4.gear.json` — **did not exist in this repo and was
    built for this comparison.** No P4 gear preset previously existed for Rogue (only
    `p1`/`p2`/`p3`/`preraid`). It was derived from the existing, already-verified
    `p3.gear.json` (same enchant-ID scheme as this repo) by cross-referencing the historical
    public wowsims/tbc source tree's P3 → P4 Rogue preset diff, which changes only 3 of 17
    slots — back, ring1, ring2 — with every enchant ID identical between P3 and P4:
    - Back: Consortium's cloak → **Cloak of Fiends** (33590)
    - Ring1: **Signet of Primal Wrath** (33496, new Zul'Aman ring)
    - Ring2: Stormrage Signet Ring (32497, was P3's ring1, shifted to ring2)
    - None of these 3 items have gem sockets, so no gem changes were needed. All other
      slots (head/neck/shoulder/chest/wrist/hands/waist/legs/feet/trinkets/weapons/ranged)
      are unchanged from `p3.gear.json`, including its existing Trinket1/2 (Dragonspine
      Trophy + Warp-Spring Coil) and default Warglaives.
- **Weapons**:
  - **Glaives**: Warglaives of Azzinoth, MH 32837 / OH 32838 (both specs).
  - **Non-Glaive**: Season 3 (Vengeful Gladiator) fast one-handers — the same substitutes
    used in the Phase 3 trinket benchmark project (results.md). Per this repo's item
    database `phase` tags, Vengeful Gladiator gear is tagged `phase: 3` while Brutal
    Gladiator (S4) gear is tagged `phase: 5`, so Vengeful remains the best non-legendary
    fast-1H substitute available at Phase 4 by this repo's own content-phase model:
    - Warrior Fury: MH Vengeful Gladiator's Cleaver (33669), OH Vengeful Gladiator's
      Chopper (34015)
    - Rogue: MH Vengeful Gladiator's Slicer (33762), OH Blade of Savagery (32369)
- **Talents/consumables/rotations**: same as the Phase 3 trinket benchmark project — Fury
  Warrior talents `3500501130201-05050005505012050115` with the `fury` APL; Rogue talents
  `00532012502-023305200005015002321151` with the `swords` APL. Consumables/race/professions
  also match that project's configs.
- **Sword Spec Fury variant**: same gear/APL/consumables as standard Warrior Fury, always
  equipped with Warglaives, but talents `350050113020105-0505000550501205011` — the standard
  Fury build with 5 points moved from Improved Berserker Stance (Fury tree) to Sword
  Specialization (Arms tree), verified field-by-field against `proto.WarriorTalents`'
  generated field numbers and this repo's `warrior.TalentTreeSizes = [3]int{23, 21, 22}` (so
  every other talent point is unchanged from the standard build). Warglaive of Azzinoth is
  `weaponType: 9` (Sword) in `assets/database/db.json`, so Sword Specialization's proc is
  live with this weapon.
- **Scratch benchmark code**: `sim/glaivebench_test.go` (temporary, deleted after extracting
  these results).

## Results table

| Scenario | Spec | Weapon | DPS | Hit Rating |
|---|---|---|---:|---:|
| 1 target | Warrior Fury | Glaives | 2795.3 | 220.0 |
| 1 target | Warrior Fury | Non-Glaive | 2667.0 | 215.0 |
| 1 target | Warrior Fury (Sword Spec) | Glaives | 2845.8 | 220.0 |
| 1 target | Rogue | Glaives | 2786.0 | 299.0 |
| 1 target | Rogue | Non-Glaive | 2590.2 | 301.0 |
| 2 target | Warrior Fury | Glaives | 4257.7 | 220.0 |
| 2 target | Warrior Fury | Non-Glaive | 4191.2 | 215.0 |
| 2 target | Warrior Fury (Sword Spec) | Glaives | 4419.2 | 220.0 |
| 2 target | Rogue | Glaives | 3348.2 | 299.0 |
| 2 target | Rogue | Non-Glaive | 3100.0 | 301.0 |
| 4 target | Warrior Fury | Glaives | 5026.6 | 220.0 |
| 4 target | Warrior Fury | Non-Glaive | 4977.9 | 215.0 |
| 4 target | Warrior Fury (Sword Spec) | Glaives | 5172.5 | 220.0 |
| 4 target | Rogue | Glaives | 3348.3 | 299.0 |
| 4 target | Rogue | Non-Glaive | 3101.8 | 301.0 |

## Glaives advantage (Glaives − Non-Glaive)

| Scenario | Spec | DPS gain from Glaives | % gain |
|---|---|---:|---:|
| 1 target | Warrior Fury | +128.3 | +4.81% |
| 1 target | Warrior Fury (Sword Spec) | +178.8 | +6.70% |
| 1 target | Rogue | +195.8 | +7.56% |
| 2 target | Warrior Fury | +66.5 | +1.59% |
| 2 target | Warrior Fury (Sword Spec) | +228.0 | +5.44% |
| 2 target | Rogue | +248.2 | +8.01% |
| 4 target | Warrior Fury | +48.7 | +0.98% |
| 4 target | Warrior Fury (Sword Spec) | +194.6 | +3.91% |
| 4 target | Rogue | +246.5 | +7.95% |

Warrior Fury (Sword Spec)'s "Non-Glaive" baseline here is the standard-talents Non-Glaive
DPS (2667.0 / 4191.2 / 4977.9), not a Sword Spec + Non-Glaive number — Sword Spec can't
proc on the Axe-type Non-Glaive substitute, so no one would realistically respec into it
without also having Glaives. This row shows the full gain of switching from "standard talents,
no Glaives" to "Glaives + Sword Spec" in one step.

## Findings

- **Fury Warrior's multi-target DPS scales strongly with target count; Rogue's does not.**
  Warrior Fury gains +52.3% going from 1→2 targets and a further +18.1% going from 2→4
  (2795.3 → 4257.7 → 5026.6), because Whirlwind/Bloodthirst cleave hits *every* nearby
  enemy. Rogue gains a one-time +20.2% jump from 1→2 targets (2786.0 → 3348.2) but then is
  essentially flat from 2→4 targets (3348.2 → 3348.3) — consistent with results.md's earlier
  finding that Blade Flurry cleaves to only **one** extra nearby target in TBC regardless of
  total enemy count, so a Rogue gets its full "cleave" benefit already at 2 targets and
  gains nothing further from a 3rd or 4th.
- **The Glaives' DPS advantage (in %) shrinks for Fury Warrior as target count rises**
  (+4.81% → +1.59% → +0.98%), since the extra weapon-damage-based throughput Warglaives
  provide doesn't scale with target count the way Whirlwind cleave does — it becomes a
  proportionally smaller slice of a bigger AOE pie. **For Rogue, the Glaives' advantage (in
  %) stays essentially constant across all three scenarios** (+7.56% → +8.01% → +7.95%),
  since Rogue's rotation doesn't cleave beyond the 2-target Blade Flurry cap, so relative
  weapon value doesn't get diluted by additional targets.
- **In absolute DPS, Glaives are worth more to Rogue than to Fury Warrior at 2 and 4
  targets** (+248.2/+246.5 DPS for Rogue vs. +66.5/+48.7 DPS for Warrior), despite Warglaives
  mattering more for Warrior at 1 target (+128.3 vs. +195.8 — closer, but still Rogue-
  favored). This reflects Combat Rogue's higher reliance on raw weapon damage/speed
  (finishers scale off combo points built by weapon-damage-based builders) versus Fury
  Warrior's Whirlwind cleave increasingly dominating the multi-target damage profile
  regardless of which weapon is equipped.
- **Sword Spec Fury is only worth taking if you actually have Warglaives (or another
  Sword-type weapon).** Its Glaives-vs-Non-Glaive-baseline delta (+6.70% to +3.91%, see
  the Glaives advantage table above) is consistently bigger than standard Fury's (+4.81% to
  +0.98%), because Sword Specialization simply can't proc without a Sword equipped — the
  respec makes a Warglaive-wielding Fury Warrior *more*, not less, dependent on keeping
  those weapons equipped.
- **Stacked together, Glaives + Sword Spec gain +3.9% to +6.7% total DPS over the
  Non-Glaive/standard-talents baseline** — biggest at 1 target (+178.8 DPS, +6.70%),
  shrinking through 2 targets (+228.0 DPS, +5.44%) and 4 targets (+194.6 DPS, +3.91%). The %
  shrinks with target count for the same reason the Glaives-alone advantage does (Whirlwind
  cleave dilutes any fixed weapon-damage-based gain across more targets), even though the
  absolute DPS gain stays roughly flat (~+180-230 DPS) across all three scenarios.
