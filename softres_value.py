"""Soft-reserve planner for Koerry (Fury Warrior) in Hyjal / Black Temple.

Scores every simmed drop for a softres.it raid and finds the best set of
reserves, using the same model as result_shaman_*_sr.md:

    value = P(drop) * [P(win | reserved) - P(win | not reserved)] * DPS gain

  - Reserved by others (n entries): not reserving means you cannot win it at
    all; reserving gives 1/(n+1).
  - Reserved by nobody: it goes to open roll. P(win) there is 1/(1+open
    rivals), where open rivals = raiders whose spec would roll MS on it,
    capped at OPEN_RIVAL_CAP. Reserving makes it 1.
  - P(drop) = regular drops / table size, table size = items the boss has in
    assets/database/db.json. Regular bosses drop 2 items, unchanged in Anniversary (Wowhead, Supremus:
    ~15% each over 13). Token bosses depend on --drop-model:
      anniversary (default): 2 regular items + 3 tier tokens of random type,
        so P(a Protector token) = 1 - (2/3)^3 = 70%.
      classic: 1 regular item + 2 tokens, P = 1 - (2/3)^2 = 56% - what
        Wowhead's mostly-2021 numbers show (Shahraz: 16% per item, 55% token).
    Trash: 1.3% per mob (Wowhead) over an estimated 100 trash kills per clear
    (--trash-mobs), so ~73% per clear.

Combinations are scored by expected DPS per slot (two wins in one slot do not
stack; rings and weapons fill two slots), not by summing single values.

Destroyer (T5) set bonus: the harness measures 2pc+4pc on its own (WSRSET).
The slots holding a Destroyer piece (WSRTIER line: which slots, how many) are
scored jointly, outcome by outcome, equipping whatever subset of won items is
actually best. With n pieces worn, n-4 tier-slot upgrades are free; one more
costs the set once:
  - 5/5: one tier-slot upgrade is free; two or more cost the set once.
  - 4/5: any tier-slot upgrade costs the set once. Each 4/5 delta already paid
    for it, so the set is added back to get the item's set-free gain first.
Two items that each look negative from 4/5 (e.g. Grips of Damnation -6) can
therefore be worth reserving together with one that breaks the set anyway.

Inputs: the WSR lines printed by sim/warriorhyjalbtsr_test.go, and a softres
raid JSON (https://softres.it/api/raid/<id>).

    python3 -I softres_value.py --sims result_warrior_hyjal_bt.raw.txt \
        --raid <raid.json> [--me Koerry] [--baseline Oct13|T5|4/5] [--reserves N]
"""

import argparse
import itertools
import json
import os
import sys
import unicodedata

OPEN_RIVAL_CAP = 4
# While the list is not full, open-roll rivals are at least what a typical 25-man
# brings for that category (2-3 rogues, 2-3 hunters, 1-2 enh, 1 ret, 1 feral, 1-2
# other warriors); otherwise an early list makes every item look uncontested.
TYPICAL_OPEN_RIVALS = {"plate": 2, "mail": 3, "leather": 3, "jewellery": 4, "sword": 3, "axe": 3,
                       "mace": 3, "dagger": 3, "fist": 3, "bow/gun/xbow": 3, "thrown": 2, "token": 4}
# Wowhead, Supremus: ~15% per item = 2 drops over his 13. Overridable with --drops-per-boss.
DROPS_PER_BOSS = 2
# Wowhead: a trash epic drops at ~1.3% per mob. Per clear that is 1 - (1 - rate)^mobs;
# the mob count is an estimate (BT and Hyjal each run to 100+ trash kills).
TRASH_DROP_PER_MOB = 0.013
TRASH_MOBS_PER_CLEAR = 100
TRASH_P_DROP = 1 - (1 - TRASH_DROP_PER_MOB) ** TRASH_MOBS_PER_CLEAR
# Token bosses: (regular items, tier tokens) per kill. "classic" matches Wowhead's
# mostly-2021 data for Mother Shahraz (16% per item over 6, 55% per token type).
TIER_BOSS_DROPS = {"anniversary": (2, 3), "classic": (1, 2)}
TOKEN_BOSSES = {"Azgalor", "Archimonde", "Mother Shahraz", "The Illidari Council", "Illidan Stormrage"}

# Onslaught Battlegear piece -> Protector token that buys it, and its boss.
TOKENS = {
    30969: (31094, "Gloves of the Forgotten Protector", "Azgalor"),
    30972: (31095, "Helm of the Forgotten Protector", "Archimonde"),
    30979: (31103, "Pauldrons of the Forgotten Protector", "Mother Shahraz"),
    30977: (31100, "Leggings of the Forgotten Protector", "The Illidari Council"),
    30975: (31091, "Chestguard of the Forgotten Protector", "Illidan Stormrage"),
}

ZONES = {"hyjal": "Hyjal", "blacktemple": "BT"}
ZONE_IDS = {"Hyjal": 3606, "BT": 3959}

# WoW spec ids, as softres.it reports them.
SPECS = {
    71: ("Warrior", "Arms"), 72: ("Warrior", "Fury"), 73: ("Warrior", "Protection"),
    65: ("Paladin", "Holy"), 66: ("Paladin", "Protection"), 70: ("Paladin", "Retribution"),
    253: ("Hunter", "Beast Mastery"), 254: ("Hunter", "Marksmanship"), 255: ("Hunter", "Survival"),
    259: ("Rogue", "Assassination"), 260: ("Rogue", "Combat"), 261: ("Rogue", "Subtlety"),
    256: ("Priest", "Discipline"), 257: ("Priest", "Holy"), 258: ("Priest", "Shadow"),
    262: ("Shaman", "Elemental"), 263: ("Shaman", "Enhancement"), 264: ("Shaman", "Restoration"),
    62: ("Mage", "Arcane"), 63: ("Mage", "Fire"), 64: ("Mage", "Frost"),
    265: ("Warlock", "Affliction"), 266: ("Warlock", "Demonology"), 267: ("Warlock", "Destruction"),
    102: ("Druid", "Balance"), 103: ("Druid", "Feral"), 104: ("Druid", "Guardian"), 105: ("Druid", "Restoration"),
}

DPS_WARRIOR = {71, 72}
RET = {70}
HUNTER = {253, 254, 255}
ROGUE = {259, 260, 261}
ENH = {263}
FERAL = {103}
PHYSICAL = DPS_WARRIOR | RET | HUNTER | ROGUE | ENH | FERAL

# Which specs roll main-spec on an item, by category.
ROLLERS = {
    "plate": DPS_WARRIOR | RET,
    "mail": HUNTER | ENH | DPS_WARRIOR,
    "leather": ROGUE | FERAL | ENH | DPS_WARRIOR,
    "jewellery": PHYSICAL,
    "sword": DPS_WARRIOR | ROGUE | HUNTER,
    "axe": DPS_WARRIOR | HUNTER | ENH,
    "mace": DPS_WARRIOR | ROGUE | ENH,
    "dagger": DPS_WARRIOR | ROGUE | HUNTER,
    "fist": DPS_WARRIOR | ROGUE | ENH | HUNTER,
    "bow/gun/xbow": HUNTER | DPS_WARRIOR | ROGUE,
    "thrown": DPS_WARRIOR | ROGUE,
    "token": {71, 72, 73} | HUNTER | {262, 263, 264},  # every Protector-class spec
}

# Slots holding a Destroyer piece in each baseline (the 4/5 head is goggles).
# Fallback for raw files from before the harness printed WSRTIER lines.
TIER_SLOTS = {"T5": {"Head", "Shoulder", "Chest", "Hands", "Legs"}, "4/5": {"Shoulder", "Chest", "Hands", "Legs"}}

SLOT_GROUP = {"Finger1": "Rings", "Finger2": "Rings", "MainHand": "Weapons", "OffHand": "Weapons",
              "Trinket1": "Trinket", "Trinket2": "Trinket"}


def fold(s):
    return unicodedata.normalize("NFKD", s).encode("ascii", "ignore").decode().lower()


def load_sims(path, baseline):
    """Parse WSR lines -> {item id: {...}} with a gain per slot."""
    items = {}
    base = None
    set_value = 0.0
    tier = None
    with open(path) as f:
        for line in f:
            parts = [p.strip() for p in line.split("|")]
            if parts[0] == "WSRBASE" and parts[1] == baseline:
                base = parts[3]
            if parts[0] == "WSRSET" and parts[1] == baseline:
                set_value = -float(parts[4])
            if parts[0] == "WSRTIER" and parts[1] == baseline:
                tier = set(parts[3].split(",")) if parts[3] else set()
            if parts[0] != "WSR" or parts[1] != baseline:
                continue
            _, _, zone, boss, slot, label, iid, _dps, delta, token, note = parts[:11]
            iid = int(iid)
            it = items.setdefault(iid, {"id": iid, "label": label, "zone": zone, "boss": boss,
                                        "token": token, "note": note, "gains": {}})
            it["gains"][slot] = float(delta)
            it["label"] = label.replace(" (MH)", "").replace(" (OH)", "")
    if not items:
        sys.exit(f"no WSR lines for baseline {baseline!r} in {path}")
    if tier is None:
        tier = TIER_SLOTS[baseline]
    return items, base, set_value, tier


def drop_model(db, drops_per_boss, model="anniversary"):
    """item id -> (P(drop), boss name) for Hyjal/BT drops."""
    npcs = {n["id"]: n["name"] for n in db["npcs"]}
    table = {}
    owner = {}
    for it in db["items"]:
        for s in it.get("sources", []):
            d = s.get("drop")
            if not d or d.get("zoneId") not in ZONE_IDS.values():
                continue
            npc = d.get("npcId")
            if npc:
                table[npc] = table.get(npc, 0) + 1
                owner[it["id"]] = npc
    p = {}
    for iid, npc in owner.items():
        name = npcs.get(npc, str(npc))
        drops = TIER_BOSS_DROPS[model][0] if name in TOKEN_BOSSES else drops_per_boss
        p[iid] = (min(1.0, drops / table[npc]), name)
    # Tokens are extra drops on their boss, not part of the regular loot.
    for piece, (_tok, _tname, boss) in TOKENS.items():
        p[piece] = (1 - (2 / 3) ** TIER_BOSS_DROPS[model][1], boss)
    return p


def category(db_item, piece_id):
    if piece_id in TOKENS:
        return "token"
    t = db_item.get("type")
    if t == 13:
        return {1: "axe", 2: "dagger", 3: "fist", 4: "mace", 9: "sword"}.get(db_item.get("weaponType"), "sword")
    if t == 14:
        return "thrown" if db_item.get("rangedWeaponType") == 4 else "bow/gun/xbow"
    return {2: "leather", 3: "mail", 4: "plate"}.get(db_item.get("armorType", 0), "jewellery")


def best_assignment(won, items, group):
    """Best total set-free gain from a set of won items in one slot group."""
    if not won:
        return 0.0
    if group in ("Rings", "Weapons"):
        slots = ("Finger1", "Finger2") if group == "Rings" else ("MainHand", "OffHand")
        best = 0.0
        for a in [None, *won]:
            for b in [None, *won]:
                if a is not None and a == b:
                    continue
                g = max(0.0, items[a]["neutral"].get(slots[0], 0.0)) if a else 0.0
                g += max(0.0, items[b]["neutral"].get(slots[1], 0.0)) if b else 0.0
                best = max(best, g)
        return best
    return max(0.0, max(max(items[i]["neutral"].values()) for i in won))


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--sims", required=True)
    ap.add_argument("--raid", required=True, help="softres raid JSON file")
    ap.add_argument("--db", default=os.path.join(os.path.dirname(os.path.abspath(__file__)), "assets/database/db.json"))
    ap.add_argument("--me", default="Koerry")
    ap.add_argument("--baseline", default="Oct13", help="Oct13 (live armory), T5 or 4/5")
    ap.add_argument("--reserves", type=int, help="override the raid's reserve limit")
    ap.add_argument("--top", type=int, default=10)
    ap.add_argument("--drops-per-boss", type=int, default=DROPS_PER_BOSS)
    ap.add_argument("--drop-model", default="anniversary", choices=sorted(TIER_BOSS_DROPS))
    ap.add_argument("--trash-mobs", type=int, default=TRASH_MOBS_PER_CLEAR,
                    help="trash mobs killed per clear; each rolls 1.3%% for a given trash epic")
    ap.add_argument("--owned", default="", help="comma-separated item ids you already have but the "
                    "baseline does not wear (e.g. 32608 Pillager's Gauntlets on a 5/5 baseline): "
                    "treated as already won, never reserved")
    ap.add_argument("--exclude", default="", help="comma-separated slot groups to leave out, "
                    "e.g. Weapons,Trinket,Ranged")
    args = ap.parse_args()

    db = json.load(open(args.db))
    dbitems = {it["id"]: it for it in db["items"]}
    raid = json.load(open(args.raid))
    items, base, set_value, tier_slots = load_sims(args.sims, args.baseline)
    free_tier_swaps = max(0, len(tier_slots) - 4)
    pdrop = drop_model(db, args.drops_per_boss, args.drop_model)
    trash_p = 1 - (1 - TRASH_DROP_PER_MOB) ** args.trash_mobs

    zones = {ZONES[i["slug"]] for i in raid["instances"] if i["slug"] in ZONES}
    hard = set(raid.get("hard_reserves") or [])
    k = args.reserves or raid["reserve_limit"]

    me = None
    others = []
    for r in raid["reserves"]:
        if fold(r["name"]) == fold(args.me):
            me = r
        else:
            others.append(r)

    # Reserve entries by others, keyed by the id that is actually reserved.
    reservers = {}
    for r in others:
        for iid in r["items"]:
            reservers.setdefault(iid, []).append(r["name"])
    specs = [r["spec"] for r in others]
    raid_size = max(i.get("slots", 25) for i in raid["instances"])
    partial = len(raid["reserves"]) < raid_size

    # Set-free gain per slot: from 4/5 every tier-slot delta paid for breaking the
    # 4pc, so add it back; the joint tier scoring below charges it once.
    for it in items.values():
        it["neutral"] = {sl: g + (set_value if free_tier_swaps == 0 and sl in tier_slots else 0.0)
                         for sl, g in it["gains"].items()}

    excluded = {g.strip() for g in args.exclude.split(",") if g.strip()}
    owned = {int(x) for x in args.owned.split(",") if x.strip()}
    rows = []
    for iid, it in items.items():
        if iid not in owned and not any(z in zones for z in it["zone"].split("+")):
            continue
        res_id = TOKENS[iid][0] if iid in TOKENS else iid
        gain = max(it["gains"].values())
        neutral = max(it["neutral"].values())
        if iid in hard or neutral <= 0:
            continue
        p, boss = pdrop.get(iid, (trash_p, "Trash"))
        if it["boss"] == "Trash":
            p, boss = trash_p, "Trash"
        cat = category(dbitems.get(iid, {}), iid)
        rivals = reservers.get(res_id, [])
        if iid in owned:
            p, p_no, p_sr = 1.0, 1.0, 1.0
        elif rivals:
            p_no, p_sr = 0.0, 1.0 / (len(rivals) + 1)
        else:
            open_rivals = sum(1 for s_ in specs if s_ in ROLLERS[cat])
            if partial:
                open_rivals = max(open_rivals, TYPICAL_OPEN_RIVALS[cat])
            open_rivals = min(OPEN_RIVAL_CAP, open_rivals)
            p_no, p_sr = 1.0 / (open_rivals + 1), 1.0
        group = SLOT_GROUP.get(next(iter(it["gains"])), next(iter(it["gains"])))
        if group in excluded:
            continue
        name = TOKENS[iid][1] + f" -> {it['label']}" if iid in TOKENS else it["label"]
        rows.append(dict(id=iid, res_id=res_id, name=name, zone=it["zone"], boss=boss, group=group,
                         slots="/".join(sl for sl, g in it["neutral"].items() if g > 0) or group,
                         gain=gain, neutral=neutral, cat=cat, rivals=rivals, p=p, p_no=p_no, p_sr=p_sr))

    byid = {r["id"]: r for r in rows}
    groups = {}
    for r in rows:
        groups.setdefault(r["group"], []).append(r["id"])
    tier_groups = [g for g in groups if g in tier_slots]
    tier_ids = {i for g in tier_groups for i in groups[g]}

    def outcomes(g, reserved):
        """(probability, best set-free gain) for every win/lose outcome in a slot group."""
        ids = groups[g]
        probs = [byid[i]["p"] * (byid[i]["p_sr"] if i in reserved else byid[i]["p_no"]) for i in ids]
        for outcome in itertools.product((0, 1), repeat=len(ids)):
            pr = 1.0
            for o, q in zip(outcome, probs):
                pr *= q if o else 1 - q
            if pr:
                yield pr, best_assignment([i for i, o in zip(ids, outcome) if o], items, g)

    def ev_group(g, reserved):
        return sum(pr * v for pr, v in outcomes(g, reserved))

    def ev_tier(reserved):
        dists = []
        for g in tier_groups:
            d = {}
            for pr, v in outcomes(g, reserved):
                d[v] = d.get(v, 0.0) + pr
            dists.append(list(d.items()))
        total = 0.0
        for combo in itertools.product(*dists):
            pr = 1.0
            for _, q in combo:
                pr *= q
            vals = sorted((v for v, _ in combo if v > 0), reverse=True)
            # Either stop at the free swaps and keep the 4pc, or take them all and lose it.
            best = sum(vals[:free_tier_swaps])
            if len(vals) > free_tier_swaps:
                best = max(best, sum(vals) - set_value)
            total += pr * best
        return total

    base_ev = {g: ev_group(g, set()) for g in groups if g not in tier_slots}
    cache = {}

    def combo_value(combo):
        tkey = ("tier", frozenset(i for i in combo if i in tier_ids))
        for key in (tkey, ("tier", frozenset())):
            if key not in cache:
                cache[key] = ev_tier(key[1])
        total = cache[tkey] - cache[("tier", frozenset())]
        for g in {byid[i]["group"] for i in combo} - set(tier_groups):
            key = (g, frozenset(i for i in combo if byid[i]["group"] == g))
            if key not in cache:
                cache[key] = ev_group(g, key[1]) - base_ev[g]
            total += cache[key]
        return total

    for r in rows:
        r["value"] = combo_value((r["id"],))
    rows = [r for r in rows if r["id"] not in owned]

    rows.sort(key=lambda r: -r["value"])

    print(f"# Soft-reserve scoring — {args.me} (Fury Warrior), baseline {args.baseline}: {base}\n")
    print(f"Destroyer 2pc+4pc worth {set_value:.1f} DPS in this baseline; {len(tier_slots)} pieces worn "
          f"({', '.join(sorted(tier_slots))}), {free_tier_swaps} free tier swap(s).\n")
    print(f"Raid {raid['id']} · instances: {', '.join(i['name'] for i in raid['instances'])} · "
          f"{len(raid['reserves'])} raiders · {k} reserves each · locked: {raid.get('locked')}")
    if partial:
        print(f"\n**List not full ({len(raid['reserves'])}/{raid_size}):** open-roll rivals are floored at a "
              f"typical 25-man raid's; re-run as more people reserve.")
    if raid.get("instance_reserve_limits"):
        print(f"\n**Per-instance reserve limits present (not modelled):** {raid['instance_reserve_limits']}")
    if hard:
        print(f"Hard-reserved (excluded): {sorted(hard)}")
    if excluded:
        print(f"Slots left out by request: {', '.join(sorted(excluded))}")
    if owned:
        print(f"Already owned (counted as won): {', '.join(items[i]['label'] for i in sorted(owned) if i in items)}")
    if me:
        print(f"\n{args.me} is on the list, reserved: {me['items']}")
    else:
        print(f"\n{args.me} is not on the list yet.")
    comp = {}
    for s in specs:
        c, sp = SPECS.get(s, ("?", str(s)))
        comp[f"{sp} {c}"] = comp.get(f"{sp} {c}", 0) + 1
    print("Others: " + ", ".join(f"{n}× {k_}" for k_, n in sorted(comp.items(), key=lambda x: -x[1])))

    print("\n## Every positive item\n")
    print("| Item | Boss | Slot | DPS | Reserved by | P(drop) | P(win) no SR | P(win) SR | Value |")
    print("|---|---|---|---|---|---|---|---|---|")
    for r in rows:
        who = ", ".join(r["rivals"]) if r["rivals"] else f"— (open: {r['cat']})"
        dps = f"{r['gain']:+.1f}" + (f" ({r['neutral']:+.1f} set-free)" if abs(r["neutral"] - r["gain"]) > 0.05 else "")
        print(f"| {r['name']} | {r['zone']} — {r['boss']} | {r['slots']} | {dps} | {who} | "
              f"{r['p']:.2f} | {r['p_no']:.2f} | {r['p_sr']:.2f} | **{r['value']:.2f}** |")

    pool = [r["id"] for r in rows[:24]]
    n = min(k, len(pool))
    combos = sorted(((combo_value(c), c) for c in itertools.combinations(pool, n)), key=lambda x: -x[0])
    print(f"\n## Best {n}-reserve combinations (slot-aware)\n")
    print("| # | Combination | Value |")
    print("|---|---|---|")
    for i, (v, c) in enumerate(combos[:args.top], 1):
        print(f"| {i} | {' + '.join(byid[x]['name'] for x in c)} | **{v:.2f}** |")

    if me:
        mine = []
        for iid in me["items"]:
            piece = next((p for p, t in TOKENS.items() if t[0] == iid), iid)
            if piece in byid:
                mine.append(piece)
        print(f"\nCurrent reserves score **{combo_value(tuple(mine)):.2f}** "
              f"(not simmed, hard-reserved or never an upgrade, scored 0: "
              f"{[i for i in me['items'] if next((p for p, t in TOKENS.items() if t[0] == i), i) not in byid]}).")


if __name__ == "__main__":
    main()
