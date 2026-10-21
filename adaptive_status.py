"""Current standings of a (possibly still running) TestWarriorHyjalBTSRAdaptive log.

Reads WSRSTEP lines (the latest estimate per candidate) and WSR lines (final
verdicts), so it works on a partial log after any round. Drop chances come from
softres_value.py's model (Anniversary drops, Wowhead-calibrated) and are shown
for Hyjal and Black Temple only.

    python3 -I adaptive_status.py <log> [--split] [--top 15] [--sort dps|ev]
        [--exclude Weapons,Trinket,Ranged] [--owned 32608]

--owned lists items you already have that the baseline does not wear. They are
left out of the tables and add a "vs today" column: an item in an owned item's
slot is worth its delta minus the owned item's; if that slot is a tier slot,
every tier-slot item also loses min(owned delta, set bonus) - you either put
the tier piece back or give up the set.
"""

import argparse
import json
import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import softres_value as sv  # noqa: E402

SLOT_GROUP = {"Finger1": "Rings", "Finger2": "Rings", "MainHand": "Weapons", "OffHand": "Weapons",
              "Trinket1": "Trinket", "Trinket2": "Trinket", "Ranged": "Ranged"}
SLOT_NAME = {"Finger1": "Ring", "Finger2": "Ring", "Trinket1": "Trinket", "Trinket2": "Trinket",
             "MainHand": "Main hand", "OffHand": "Off hand"}
GROUPS = [("Hyjal", {"Hyjal", "Hyjal+BT"}), ("Black Temple", {"BT", "Hyjal+BT"}), ("Zul'Aman", {"ZA"}),
          ("External (reputation, Season 3, Phase 4 without a source)", {"Rep", "S3", "P4"})]


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("log")
    ap.add_argument("--top", type=int, default=15)
    ap.add_argument("--exclude", default="")
    ap.add_argument("--owned", default="")
    ap.add_argument("--split", action="store_true", help="one table per instance / external")
    ap.add_argument("--sort", default="dps", choices=["dps", "ev"])
    args = ap.parse_args()
    excluded = {g.strip() for g in args.exclude.split(",") if g.strip()}
    owned = {int(x) for x in args.owned.split(",") if x.strip()}

    here = os.path.dirname(os.path.abspath(__file__))
    db = json.load(open(os.path.join(here, "assets/database/db.json")))
    pdrop = sv.drop_model(db, sv.DROPS_PER_BOSS)

    meta, rounds, est, tier, set_value = {}, [], {}, set(), 0.0
    for line in open(args.log):
        p = [x.strip() for x in line.split("|")]
        if p[0] in ("WSRBASE", "WSRSUMMARY"):
            meta[p[0]] = " | ".join(p[2:])
        elif p[0] == "WSRSET":
            set_value = -float(p[4])
        elif p[0] == "WSRTIER":
            tier = set(p[3].split(",")) if p[3] else set()
        elif p[0] == "WSRROUND":
            rounds.append(" | ".join(p[2:]))
        elif p[0] == "WSRSTEP":
            slot, label, iid, n, delta, ci = p[2:8]
            zone, boss = (p[8], p[9]) if len(p) >= 10 else ("", "")
            prev = est.get((slot, label), {})
            est[(slot, label)] = dict(slot=slot, label=label, id=int(iid), n=int(n), delta=float(delta),
                                      ci=float(ci), zone=zone or prev.get("zone", ""),
                                      boss=boss or prev.get("boss", ""), state="open")
        elif p[0] == "WSR" and len(p) >= 14 and (p[4], p[5]) in est:
            e = est[(p[4], p[5])]
            e["state"], e["zone"], e["boss"] = p[13], p[2], p[3]

    print(meta.get("WSRBASE", "(baseline not finished)"))
    print(meta.get("WSRSUMMARY") or (rounds[-1] if rounds else "(round 1 in progress)"))
    print(f"Tier set bonus {set_value:.1f}; tier slots {', '.join(sorted(tier))}")

    # "vs today" adjustments from owned items.
    adjust_slot, adjust_tier = {}, 0.0
    for e in est.values():
        if e["id"] in owned:
            adjust_slot[e["slot"]] = e["delta"]
            if e["slot"] in tier:
                adjust_tier = max(adjust_tier, min(e["delta"], set_value))
            print(f"Owned: {e['label']} ({e['slot']}) {e['delta']:+.1f} ± {e['ci']:.1f}, n={e['n']}")
    print()

    def today(e):
        d = e["delta"] - adjust_slot.get(e["slot"], 0.0)
        if e["slot"] in tier and e["slot"] not in adjust_slot:
            d -= adjust_tier
        return d

    def drop(e):
        if not ({"Hyjal", "BT", "Hyjal+BT"} & {e["zone"]}):
            return None
        if e["boss"] == "Trash":
            return sv.TRASH_P_DROP
        return pdrop.get(e["id"], (None, ""))[0]

    rows = [e for e in est.values() if e["id"] not in owned
            and SLOT_GROUP.get(e["slot"], e["slot"]) not in excluded]

    def table(title, entries):
        key = (lambda e: -(max(today(e), 0) * (drop(e) or 0))) if args.sort == "ev" else (lambda e: -today(e))
        entries = sorted(entries, key=key)[:args.top]
        print(f"### {title}\n")
        head = "| # | Item | Source | Slot | DPS | ±95% | Iterations |"
        head += " vs today |" if owned else ""
        print(head + " Drop | DPS × drop |")
        print("|" + "---|" * (head.count("|") + 1))
        for i, e in enumerate(entries, 1):
            p = drop(e)
            t = today(e)
            slot = SLOT_NAME.get(e["slot"], e["slot"]) + (" ·T" if e["slot"] in tier else "")
            line = f"| {i} | {e['label']} | {e['boss']} | {slot} | {e['delta']:+.1f} | {e['ci']:.1f} | {e['n']} |"
            line += f" {t:+.1f} |" if owned else ""
            line += f" {p:.2f} | {max(t, 0) * p:.1f} |" if p is not None else " — | — |"
            if e["state"] == "open":
                line += " *(open)*"
            print(line)
        print()

    if args.split:
        for title, zones in GROUPS:
            table(title, [e for e in rows if e["zone"] in zones])
    else:
        table("All candidates", rows)


if __name__ == "__main__":
    main()
