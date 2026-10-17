# Black Temple soft-reserve plan — Koerrii (Enhancement Shaman)

Optimising the **4** soft reserves for <https://softres.it/raid/6dBB6Vnd> (Black Temple, 25-man,
Horde, TBC). Companion to [result_shaman_hyjal_sr.md](result_shaman_hyjal_sr.md); same method,
same baseline.

**Baseline: 2237.6 DPS** on Koerrii's current armory gear. Every BT drop an Enhancement Shaman can
equip was simmed as a single-slot swap at 20,000 iterations (`sim/shamanbtsr_test.go`).

## Status: ACTION NEEDED — refreshed 2026-09-01, 22/25, still open

The list jumped 14 → 22 reserves and **three of Koerrii's four picks picked up rivals**. Current
reserves are worth **17.67**, down from 25.78. The best available set is now **19.62**, so there is
about **2 DPS per clear** to recover — and the list is still unlocked, so it is recoverable.

**Cuica (Ret Paladin) took two of the four**: Shadow-walker's Cord and Band of Devastation.
**Wfr (BM Hunter)** took Insidious Bands, and **Hanebuth (BM Hunter)** joined Madness of the
Betrayer.

### Recommended change

| Action | Item | Slot | DPS | Rivals | Value |
|--------|------|------|-----|--------|-------|
| **Drop** | Shadow-walker's Cord | Waist | +29.3 | 1 (Cuica) | 4.19 |
| **Add** | **Boneweave Girdle** | Waist | +29.1 | **0** | **5.97** |
| **Drop** | Band of Devastation | Finger | +17.6 | 1 (Cuica) | 4.40 |
| **Add** | **Mantle of Darkness** | Shoulder | +21.0 | **0** | **4.50** |
| Keep | Madness of the Betrayer | Trinket | +39.3 | 3 | 4.37 |
| Keep | Insidious Bands | Wrist | +28.3 | 1 (Wfr) | 4.72 |

Final list: **Madness of the Betrayer · Boneweave Girdle · Insidious Bands · Mantle of Darkness**
(19.55).

The waist swap is the unambiguous one: Boneweave Girdle is the same slot at essentially the same
DPS (+29.1 vs +29.3) and nobody is on it. Dropping Band of Devastation also removes the pick that
rested on this model's softest input — the assumed 0.5 trash drop rate.

The top six combinations all sit between 19.45 and 19.62, which is well inside the model's
uncertainty, so the fourth slot is a genuine toss-up. Keeping **Madness of the Betrayer** (19.55)
over the pure EV-max set (19.62, which swaps it for Grips of Damnation) costs 0.07 and buys a 25%
shot at the biggest single upgrade in the instance — a trinket that stays BiS well past this raid.

### What is still uncontested

Even at 22/25 these remain reserved by nobody: **Boneweave Girdle** (+29.1), **Mantle of Darkness**
(+21.0), **Grips of Damnation** (+19.2), **Shoulders of the Hidden Predator** (+18.7), **Wraps of
Precise Flight** (+17.6), **Nether Shadow Tunic** (+15.7), **Unstoppable Aggressor's Ring** (+14.0),
**Shadowmoon Destroyer's Drape** (+11.7). The raid's melee are piling onto head, neck, rings and
trinkets and ignoring the armour slots.

### Items that have become traps

| Item | DPS | Rivals | Value | Note |
|------|-----|--------|-------|------|
| Cursed Vision of Sargeras | +34.2 | **6** | 1.15 | Six reservers on Illidan's long table — hopeless |
| Stormrage Signet Ring | +24.0 | 3 | 1.41 | Also Illidan; both hunters and both warriors want it |
| Choker of Endless Nightmares | +19.4 | 3 | 1.49 | Wfr joined Sliumpius and Srhóx |

> Rival counts exclude Koerrii, so they read as "who else wants this".

## Original recommendation (2026-09-01, at 10/25 reserved)

| # | Reserve | Boss | Slot | DPS | Rivals | Why |
|---|---------|------|------|-----|--------|-----|
| 1 | **Madness of the Betrayer** | Illidari Council | Trinket | **+39.3** | 2 | Biggest single BT upgrade, and the Council has a short loot table |
| 2 | **Insidious Bands** | Teron Gorefiend | Wrist | **+28.3** | **0** | Highest-value reserve on the list — big gain, nobody on it, mid-raid boss |
| 3 | **Shadow-walker's Cord** | Shade of Akama | Waist | **+29.3** | **0** | Near-identical to Boneweave Girdle but from an earlier, easier boss |
| 4 | **Band of Devastation** | Trash | Finger | **+17.6** | **0** | Trash drops need no boss kill, so it is the cheapest slot to fill |

Expected value: **~25.8 DPS per full clear** — the best of any 4-item combination.

If your guild treats trash loot as unreserveable or it drops rarely, swap #4 for **Mantle of
Darkness** (Naj'entus, shoulder, +21.0, 0 rivals) — that set is worth ~24.0.

## Every BT drop, ranked

Single-slot swap against your current gear. Rings and trinkets are listed at their better slot.

| Item | Boss | Slot | DPS |
|------|------|------|-----|
| Madness of the Betrayer | Illidari Council | Trinket | **+39.3** |
| Cursed Vision of Sargeras | Illidan Stormrage | Head | **+34.2** |
| Shadow-walker's Cord | Shade of Akama | Waist | **+29.3** |
| Boneweave Girdle | Reliquary of the Lost | Waist | **+29.1** |
| Insidious Bands | Teron Gorefiend | Wrist | **+28.3** |
| Stormrage Signet Ring | Illidan Stormrage | Finger | **+24.0** |
| Mantle of Darkness | Naj'entus | Shoulder | +21.0 |
| Choker of Endless Nightmares | Supremus | Neck | +19.4 |
| Grips of Damnation | Reliquary of the Lost | Hands | +19.2 |
| Shoulders of the Hidden Predator | Shade of Akama | Shoulder | +18.7 |
| Fists of Mukoa | Naj'entus | Hands | +18.1 |
| Band of Devastation | Trash | Finger | +17.6 |
| Wraps of Precise Flight | Supremus | Wrist | +17.6 |
| Shadowmaster's Boots | Mother Shahraz | Feet | +16.0 |
| Nether Shadow Tunic | Supremus | Chest | +15.7 |
| Unstoppable Aggressor's Ring | Gurtogg Bloodboil | Finger | +14.0 |
| Softstep Boots of Tracking | Teron Gorefiend | Feet | +13.0 |
| Shadowmoon Destroyer's Drape | Teron Gorefiend | Back | +11.7 |
| Forest Prowler's Helm | Illidari Council | Head | +8.8 |
| Rising Tide | Naj'entus | Main hand | +8.3 |
| Choker of Serrated Blades | Trash | Neck | +6.6 |
| Vest of Mounting Assault | Gurtogg Bloodboil | Chest | -6.0 |
| Treads of the Den Mother | Trash | Feet | -9.1 |
| Band of the Abyssal Lord | Supremus | Finger | -23.5 |
| Pendant of Titans | Reliquary of the Lost | Neck | -28.0 |
| Syphon of the Nathrezim | Supremus | Main hand | -33.1 |
| Shard of Azzinoth | Illidan Stormrage | Off hand | -91.1 |
| Swiftsteel Bludgeon | Trash | Off hand | -112.4 |
| The Brutalizer | Supremus | Off hand | -126.0 |
| Messenger of Fate | Gurtogg Bloodboil | Main hand | -174.6 |

## The reserve maths

Same model as the Hyjal plan: a reserve is worth `P(drop) × [P(win with SR) − P(win without SR)] ×
DPS gain`. An item somebody else already reserved is protected away from you entirely, so your
reserve buys the full roll; an unreserved item you would often win on an open roll anyway, so your
reserve only buys the difference.

| Item | Gain | Rivals | P(drop) | P(win) no SR | P(win) with SR | **Value** |
|------|------|--------|---------|--------------|----------------|-----------|
| Insidious Bands | +28.3 | 0 | 0.33 | 0.25 | 1.00 | **7.08** |
| Band of Devastation | +17.6 | 0 | 0.50 | 0.25 | 1.00 | **6.60** |
| Shadow-walker's Cord | +29.3 | 0 | 0.29 | 0.25 | 1.00 | **6.28** |
| Boneweave Girdle | +29.1 | 0 | 0.31 | 0.33 | 1.00 | **5.97** |
| Madness of the Betrayer | +39.3 | 2 | 0.44 | 0.00 | 0.33 | **5.82** |
| Mantle of Darkness | +21.0 | 0 | 0.29 | 0.25 | 1.00 | 4.50 |
| Grips of Damnation | +19.2 | 0 | 0.31 | 0.25 | 1.00 | 4.43 |
| Wraps of Precise Flight | +17.6 | 0 | 0.31 | 0.33 | 1.00 | 3.61 |
| Nether Shadow Tunic | +15.7 | 0 | 0.31 | 0.25 | 1.00 | 3.62 |
| Shoulders of the Hidden Predator | +18.7 | 0 | 0.29 | 0.33 | 1.00 | 3.56 |
| Shadowmaster's Boots | +16.0 | 1 | 0.44 | 0.00 | 0.50 | 3.56 |
| Unstoppable Aggressor's Ring | +14.0 | 0 | 0.31 | 0.25 | 1.00 | 3.23 |
| Shadowmoon Destroyer's Drape | +11.7 | 0 | 0.33 | 0.25 | 1.00 | 2.92 |
| Fists of Mukoa | +18.1 | 1 | 0.29 | 0.00 | 0.50 | 2.59 |
| Softstep Boots of Tracking | +13.0 | 1 | 0.33 | 0.00 | 0.50 | 2.17 |
| Cursed Vision of Sargeras | +34.2 | 3 | 0.24 | 0.00 | 0.25 | 2.01 |
| Choker of Endless Nightmares | +19.4 | 2 | 0.31 | 0.00 | 0.33 | 1.99 |
| Stormrage Signet Ring | +24.0 | 2 | 0.24 | 0.00 | 0.33 | 1.88 |

Best four-item sets, accounting for same-slot items not stacking:

| Combination | Value |
|-------------|-------|
| **Madness + Shadow-walker's Cord + Insidious Bands + Band of Devastation** | **25.78** |
| Madness + Boneweave Girdle + Insidious Bands + Band of Devastation | 25.47 |
| Shadow-walker's Cord + Insidious Bands + Mantle of Darkness + Band of Devastation | 24.45 |
| Shadow-walker's Cord + Insidious Bands + Grips of Damnation + Band of Devastation | 24.38 |
| Madness + Insidious Bands + Mantle of Darkness + Band of Devastation | 24.00 |

## Traps to avoid

- **Don't reserve Cursed Vision of Sargeras**, even though it is your second-biggest upgrade at
  +34.2. It is the most contested item in the raid (Snipemaster has it twice, Srhóx once), and
  Illidan has the longest loot table in the instance. Your reserve is worth **2.01** — less than a
  third of what Insidious Bands buys. Roll on it if it drops; don't spend a slot on it.
- **Don't reserve both waist items.** Shadow-walker's Cord (+29.3) and Boneweave Girdle (+29.1)
  are nearly identical and both uncontested, so reserving both mostly buys the same upgrade twice.
  Take the Shade of Akama one — it is the third boss rather than the sixth.
- **Don't reserve a weapon.** Every BT one-hander is a downgrade on your Merciless Gladiator's
  Cleavers, including Syphon of the Nathrezim (-33.1) — which this repo's own P3 preset equips.
  You are an Orc; the axes and their +5 expertise win. Rising Tide is the only positive at +8.3,
  and it is not worth a slot.
- **Don't reserve a tier 6 token.** The BT Protector tokens are 31091 / 31100 / 31103 (shoulders,
  chest, legs). Skyshatter Pauldrons sim at +19.6, Tunic at +2.7, Pants at +7.2 — the shoulders are
  respectable but still lose to Mantle of Darkness (+21.0), which nobody has reserved.

## Who you are competing with

**Sliumpius (Enhancement Shaman)** reserved Fists of Mukoa, Choker of Endless Nightmares, Softstep
Boots of Tracking and Forest Prowler's Helm — **none of your top five**. Their list is noticeably
worse than yours; nothing in it exceeds +19.4 for you.

The real competition is the two Fury Warriors: **Axxen** and **Srhóx** both hold Madness of the
Betrayer and Stormrage Signet Ring, and Srhóx also has Cursed Vision and Choker of Endless
Nightmares. **Snipemaster (Combat Rogue)** has Cursed Vision twice plus Shadowmaster's Boots and
Ring of Deceitful Intent.

That leaves the entire waist, wrist, shoulder and hands range completely unreserved — which is
where three of your four picks come from.

## Assumptions

- **Drop rate**: 4 items per boss. Table sizes from `assets/database/db.json` (zone 3959) plus 3
  tier tokens each on Mother Shahraz, the Illidari Council and Illidan, whose database tables are
  visibly short. That gives 0.24-0.44 per item depending on boss. Trash is modelled at a flat 0.5
  per clear, which is the softest number here — Band of Devastation's ranking depends on it.
- **Open-roll rivals** for unreserved items: 2 for mail (Sliumpius plus a hunter), 3 for leather
  and jewellery (Snipemaster, plus the two Fury Warriors, who demonstrably roll on leather — both
  reserved Cursed Vision).
- **Reserves refreshed 2026-09-01**: 22 of 25 raiders, still unlocked. Rival counts exclude Koerrii.
  The tables further down this file are the original 10/25 snapshot; the Status section at the top
  carries the current numbers.
- **Single target, patchwerk**, same as [result_shaman.md](result_shaman.md).
- Black Temple drops **no ranged upgrade** for you: Twisted Blades of Zarak is Thrown
  (`rangedWeaponType: 4`) and shaman totems are `rangedWeaponType: 8`. Warglaives of Azzinoth are
  swords and unusable by a shaman.
