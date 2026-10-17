# Hyjal + Black Temple upgrades — Koerry (Fury Warrior, 5/5 T5)

Every Hyjal and Black Temple drop a Fury warrior could want, simmed as a single-slot swap on
Koerry's gear at 20,000 iterations (`sim/warriorhyjalbtsr_test.go`). This is the groundwork for
the soft-reserve plan. The reserve picks come once the softres list is posted, using
`softres_value.py`.

## Status: locked — <https://softres.it/raid/9b68rzRe>, 29 raiders (checked 2026-10-14)

Re-scored on Koerry's **live gear** (armory 2026-10-13) with Destroyer Gauntlets back on for a 5/5
baseline (2325 DPS; 2341 as actually worn, 4/5 with Pillager's Gauntlets). Pillager's is counted as
owned (`--owned 32608`), so tier-slot wins are valued as "swap the T5 gloves back on". These are
**3,000-iteration** numbers, ±2-3 DPS.

Koerry's locked reserves:

| Reserve | Slot | DPS | Reserved by others | Chance to win if it drops | Value |
|---------|------|-----|--------------------|---------------------------|-------|
| Leggings of Divine Retribution | Legs ·T | +48 | Páka | 1 in 2 | 4.9 |
| Dreadboots of the Legion | Feet | +39 | Raîns | 1 in 2 | 5.7 |
| Grips of Silent Justice | Hands ·T | +33 (+14 over Pillager's) | nobody | certain | 3.1 |
| Cursed Vision of Sargeras | Head ·T | +40 | Aloraz, Hardtøbeat, Keil, Keilhunt, Shizzy | 1 in 6 | 1.0 |
| Stormrage Signet Ring | Ring | +27 | Hardtøbeat, Keil, Keilhunt, Páka, Raîns ×3, Shizzy | 1 in 9 | 0.7 |
| | | | | **combined** | **15.3** |

The best set available at lock would have scored 26.3: Softstep Boots, Band of Devastation,
Dreadboots, Leggings of Divine Retribution and Deadly Cuffs. The two Illidan picks are where the
value went. Since the list is locked, **roll on these if they drop unreserved**: Softstep Boots of
Tracking (+26), Band of Devastation (+15), Deadly Cuffs (+17), Choker of Serrated Blades (+11).

## Brewfest: Empty Mug of Direbrew beats Madness of the Betrayer

Coren Direbrew's Anniversary trinket (281739, ilvl 141: +98 AP, use +376 AP for 20s, 2 min
cooldown) only reached this repo's database in the 2026-10-14 merge of `origin/master`. Simmed on
the same 5/5 gear (`sim/warriorbrewfest_test.go`, 20,000 iterations, post-merge baseline 2225):

| Trinkets | DPS | vs Badge + DST |
|----------|-----|----------------|
| **Empty Mug of Direbrew + Dragonspine Trophy** | 2276 | **+51** |
| Madness of the Betrayer + Dragonspine Trophy | 2263 | +39 |
| Bloodlust Brooch + Dragonspine Trophy | 2252 | +28 |
| Madness of the Betrayer + Empty Mug of Direbrew | 2238 | +14 |
| Empty Mug of Direbrew + Bloodlust Brooch | 2209 | -16 |
| Badge of the Swarmguard + Empty Mug of Direbrew | 2202 | -23 |

- **The Mug replaces Badge of the Swarmguard, not Dragonspine Trophy.** Keep DST.
- **Once you have the Mug, Madness of the Betrayer is no upgrade at all.** Mug + DST beats every
  pairing that includes Madness. That's a reason to leave trinkets off the soft-reserve list.

## Baseline: 5/5 Destroyer — 2220 DPS

The live armory now shows a **PvP snapshot** (Merciless Gladiator plate, Stormherald), so the
baseline is the PvE set from [result_warrior_t5.md](result_warrior_t5.md): Aug 31 armory gear plus
Dragonspine Trophy, with **Destroyer Battle-Helm** for a full 5/5 Destroyer. Fury talents and APL
throughout.

5/5 rather than the goggles because it costs nothing (the two helms are a wash: 2220 vs 2218), and
it makes your first tier-slot upgrade free. From 4/5, every shoulder/chest/hands/legs upgrade reads
~21 DPS worse, because that swap kills the 4pc. Those numbers are still in the raw file, and
`softres_value.py --baseline 4/5` scores them if ever needed.

## TL;DR

1. **Your off hand is the weakest slot.** Merciless Gladiator's Slicer (2.6 speed) loses ~40 DPS to
   a fast off hand: **Blade of Savagery** (Mother Shahraz, +41) or **Swiftsteel Bludgeon** (BT
   trash, +41). You only need one. A second one is still +12 in the main hand.
2. **The big non-weapon upgrades are all in BT:** Leggings of Divine Retribution (+37), Madness of the
   Betrayer (+37), Cursed Vision of Sargeras (+33), Dreadboots of the Legion (+31), Shadowmaster's
   Boots (+27).
3. **Hyjal is thin for you.** Tracker's Blade (+23, off hand), Midnight Chestguard (+22) and
   Legguards of Endless Rage (+21) are the best it has.
4. **The T5 set bonus is worth only ~27 DPS**, less than any top-10 upgrade. Never pass on an item
   to keep it.
5. **T6 tokens are mediocre for you** (+12 to +20). Each one is beaten by a non-tier item in
   the same slot, and every warrior, hunter and shaman in the raid can roll on them.

## How the T5 set bonus is handled

- **Measured directly.** The harness also sims your gear with set-less copies of the Destroyer pieces
  (identical stats, no set membership). Losing 2pc + 4pc costs **26.7 DPS**. The 2pc (Overpower AP)
  is near-worthless to Fury in Berserker stance, so almost all of it is the 4pc (-5 rage on
  Bloodthirst).
- **The first tier-slot upgrade is free.** With five pieces, replacing any one of head, shoulder,
  chest, hands or legs (marked **·T** below) keeps the 4pc. Every number in the tables is that
  first, free swap.
- **The second one costs ~27.** If you win two tier-slot items in the same clear, equipping both
  drops you to 3/5. `softres_value.py` scores the tier slots together, outcome by outcome. It
  equips both only when that beats wearing just the better one.

In practice: two tier-slot upgrades together are worth their sum minus ~27. Leggings of Divine
Retribution + Cursed Vision is still +43 (37 + 33 − 27). Two small ones, like a T6 token plus
Blood-stained Pauldrons, mostly cancel out.

Bold = 15+ DPS. For weapons, the other hand is shown in brackets.
## Every Hyjal drop that is an upgrade

| Item | Boss | Slot | DPS |
|---|---|---|---|
| Tracker's Blade | Rage Winterchill | Off hand (main hand -28) | **+23** |
| Midnight Chestguard | Archimonde | Chest ·T | **+22** |
| Legguards of Endless Rage | Archimonde | Legs ·T | **+21** |
| Boundless Agony | Azgalor | Off hand (main hand -31) | **+20** |
| Bow-stitched Leggings | Azgalor | Legs ·T | **+19** |
| Fist of Molten Fury | Trash | Off hand | **+18** |
| Blade of Infamy | Anetheron | Main hand (off hand +11) | +15 |
| Gloves of the Forgotten Protector → Onslaught Gauntlets (T6) | Azgalor | Hands ·T | +14 |
| Deadly Cuffs | Rage Winterchill | Wrist | +13 |
| Helm of the Forgotten Protector → Onslaught Battle-Helm (T6) | Archimonde | Head ·T | +12 |
| Blood-stained Pauldrons | Rage Winterchill | Shoulder ·T | +11 |
| Quickstrider Moccasins | Anetheron | Feet | +7 |
| Furious Shackles | Rage Winterchill | Wrist | +6 |
| Don Alejandro's Money Belt | Anetheron | Waist | +6 |
| Shady Dealer's Pantaloons | Azgalor | Legs ·T | +5 |
| Choker of Serrated Blades | Trash | Neck | +5 |
| Claw of Molten Fury | Trash | Main hand | +1 |

## Every Black Temple drop that is an upgrade

| Item | Boss | Slot | DPS |
|---|---|---|---|
| Warglaive of Azzinoth | Illidan Stormrage | Off hand (main hand +44) | **+62** |
| Swiftsteel Bludgeon | Trash | Off hand (main hand +12) | **+41** |
| Blade of Savagery | Mother Shahraz | Off hand (main hand +12) | **+41** |
| Leggings of Divine Retribution | Gurtogg Bloodboil | Legs ·T | **+37** |
| Madness of the Betrayer | Illidari Council | Trinket (replaces Badge of the Swarmguard) | **+37** |
| Cursed Vision of Sargeras | Illidan Stormrage | Head ·T | **+33** |
| Dreadboots of the Legion | Reliquary of the Lost | Feet | **+31** |
| Shadowmaster's Boots | Mother Shahraz | Feet | **+27** |
| Stormrage Signet Ring | Illidan Stormrage | Ring | **+23** |
| Grips of Silent Justice | Shade of Akama | Hands ·T | **+22** |
| Rising Tide | Naj'entus | Main hand (off hand +2) | **+22** |
| Unstoppable Aggressor's Ring | Gurtogg Bloodboil | Ring | **+21** |
| Pauldrons of the Forgotten Protector → Onslaught Shoulderblades (T6) | Mother Shahraz | Shoulder ·T | **+20** |
| Softstep Boots of Tracking | Teron Gorefiend | Feet | **+19** |
| Leggings of the Forgotten Protector → Onslaught Greaves (T6) | Illidari Council | Legs ·T | **+18** |
| Helm of the Illidari Shatterer | Illidari Council | Head ·T | **+18** |
| Shadowmoon Destroyer's Drape | Teron Gorefiend | Back | **+16** |
| Grips of Damnation | Reliquary of the Lost | Hands ·T | **+15** |
| Choker of Endless Nightmares | Supremus | Neck | **+15** |
| Chestguard of the Forgotten Protector → Onslaught Breastplate (T6) | Illidan Stormrage | Chest ·T | **+15** |
| Fists of Mukoa | Naj'entus | Hands ·T | +14 |
| Band of Devastation | Trash | Ring | +13 |
| The Brutalizer | Supremus | Off hand (main hand +2) | +13 |
| Messenger of Fate | Gurtogg Bloodboil | Off hand (main hand -44) | +10 |
| Pillager's Gauntlets | Trash | Hands ·T | +9 |
| Heartshatter Breastplate | Mother Shahraz | Chest ·T | +8 |
| Insidious Bands | Teron Gorefiend | Wrist | +7 |
| Ring of Deceitful Intent | Shade of Akama | Ring | +7 |
| Boneweave Girdle | Reliquary of the Lost | Waist | +7 |
| Mantle of Darkness | Naj'entus | Shoulder ·T | +7 |
| Shoulders of the Hidden Predator | Shade of Akama | Shoulder ·T | +6 |
| Syphon of the Nathrezim | Supremus | Main hand (off hand -6) | +6 |
| Choker of Serrated Blades | Trash | Neck | +5 |
| Twisted Blades of Zarak | Teron Gorefiend | Ranged | +2 |
| Forest Prowler's Helm | Illidari Council | Head ·T | +1 |
| Shadow-walker's Cord | Shade of Akama | Waist | +1 |
| Wraps of Precise Flight | Supremus | Wrist | +0 |

**Warglaives of Azzinoth as a pair: +198.** One alone is +44 MH / +62 OH; the two-piece haste proc is
the rest. In the last combined raid they were hard-reserved, so they are left out of the reserve
maths. Claw + Fist of Molten Fury together are +19, so that pair adds nothing over the Fist alone.

Ring rows show the better of the two ring slots: replacing Shapeshifter's Signet wins by 1-4 DPS
despite it being your only expertise. Rings and weapons are scored as two slots, so two rings, or
an off hand plus a main hand, do stack in the reserve maths.

## Things worth knowing before reserving

- **Two fast off hands is not a waste.** The second one goes in the main hand at +12. That is
  modest, and Rising Tide is better there (+22).
- **Main-hand drops compete with the S3 vendor.** [result_warrior_t5.md](result_warrior_t5.md)
  has Vengeful Gladiator's Cleaver at +48 in the main hand (Orc axe expertise). If you are buying
  it, do not spend a reserve on Rising Tide or Blade of Infamy.
- **Swiftsteel Bludgeon is a trash drop.** It ties Blade of Savagery, but the trash drop chance
  (0.5 per clear) is the softest number in the model. Blade of Savagery is the safer reserve if
  both are open.
- **Don't double up on boots.** Dreadboots of the Legion (+31) and Shadowmaster's Boots (+27) share a
  slot; the reserve maths will only pick both if the second is uncontested and cheap.

## What happens when you post the softres link

```
curl -s https://softres.it/api/raid/<id> -o raid.json
python3 -I softres_value.py --sims result_warrior_hyjal_bt.raw.txt --raid raid.json --me Koerry
#   add --exclude Weapons,Trinket,Ranged to leave those slots out
```

For every positive item, it reports who else reserved it and the expected DPS a reserve buys.
It then lists the best N-reserve combinations and scores Koerry's current reserves, if any.
The model is the one the shaman reports used:

- `value = P(drop) × [P(win | reserved) − P(win | not reserved)] × DPS gain`.
- **Reserved by others** (n entries): not reserving means you can't get it at all; reserving gives
  1/(n+1).
- **Reserved by nobody**: it goes to open roll. P(win) there is 1/(1 + raiders whose spec would
  roll main-spec on it, capped at 4). Plate counts warriors and rets; leather counts rogues,
  ferals, enh shamans and warriors. Tokens count every warrior, hunter and shaman.
- **Drop rate** (2026-10-14): a regular boss drops **2** items (unchanged in Anniversary) from its table in
  `assets/database/db.json` (Wowhead, Supremus: ~15% each over 13 items). A token boss (Azgalor,
  Archimonde, Shahraz, Council, Illidan) drops **2 regular items + 3 tier tokens** of random type in
  Anniversary, so a Protector token shows up 70% of the time. Wowhead's figures for Shahraz (16%
  per item, 55% per token) fit 1 item + 2 tokens instead; they are still mostly 2021 TBC Classic
  kills. `--drop-model classic` uses that version. Trash epics drop at ~1.3% per mob (Wowhead); over an
  estimated 100 trash kills per clear that is ~73% (`--trash-mobs` to change the count). Earlier
  sections used 4 drops per boss.
- **Combinations are scored slot by slot**: two wins in one slot don't stack, rings and weapons fill
  two slots, and the tier slots carry the set-bonus rule above.

## Engine issue found while building this

**Black Bow of the Betrayer (32336) panics the sim for any non-hunter.** This is the same bug class as
Beast-tamer's Shoulders in [result_shaman_hyjal_sr.md](result_shaman_hyjal_sr.md).
`sim/hunter/item_sets.go:253` casts the wearer to `HunterAgent` unconditionally. The bow's only
effect is hunter mana on ranged hits, so it is left out; as a warrior's ranged slot it would be a
stat stick at best.

## How this was measured

- **Harness**: `sim/warriorhyjalbtsr_test.go`, `TestWarriorHyjalBTSR`. Raw output (both baselines) is
  `result_warrior_hyjal_bt.raw.txt`. Re-run with
  `go test -tags with_db ./sim/ -run TestWarriorHyjalBTSR -v -timeout 0`; set `WARRIOR_ITERATIONS`
  (default 20000) to trade accuracy for speed. ~170 configurations, ~17 minutes.
- **Same sim setup as [result_warrior_t5.md](result_warrior_t5.md)**: Orc Fury
  (`3500501130201-05050005505012050115`), fury APL, single target 180s ±5s, full raid buffs and
  debuffs, `RandomSeed: 101`. Candidates are gemmed from their own sockets with P3 epic gems and
  inherit your enchant for the slot. The numbers match that report's single-slot swaps to within
  noise.
- **Candidate pool**: every Hyjal/BT drop with physical DPS stats in any armour type, all
  one-handers in both hands (main-hand-only and off-hand-only items in their own hand only),
  ranged weapons, and each Protector token as the Onslaught Battlegear piece it buys. Two-handers,
  caster/healer pieces and pure tank pieces are excluded.
- **Noise**: ±1 DPS per number at 20,000 iterations; differences under ~5 DPS are not meaningful.
