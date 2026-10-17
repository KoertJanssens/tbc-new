"""Weapon gains as a percentage of baseline DPS, side by side for 1-4 targets.

Each log is a TestWarriorHyjalBTSRAdaptive run (single target, or a
WARRIOR_TARGETS=n weapons-only cleave run). The latest WSRSTEP estimate per
candidate is used, so partial logs work; a percentage is delta / baseline DPS
of the same run, and the ± is the 95% half-width on the same scale.

    python3 -I weapon_cleave_table.py 1=<log> 2=<log> 3=<log> 4=<log> [--hand OffHand] [--top 20]
"""

import argparse
import os
import re

SLOT_NAME = {"MainHand": "Main hand", "OffHand": "Off hand"}


def read(path):
    base, est, final = None, {}, {}
    for line in open(path):
        p = [x.strip() for x in line.split("|")]
        if p[0] == "WSRBASE":
            base = float(re.search(r"DPS ([\d.]+)", p[3]).group(1))
        elif p[0] == "WSRSTEP" and p[2] in SLOT_NAME:
            est[(p[2], p[3])] = (float(p[6]), float(p[7]), int(p[5]), p[8] if len(p) > 8 else "")
        elif p[0] == "WSR" and len(p) >= 14 and p[4] in SLOT_NAME:
            final[(p[4], p[5])] = p[13]
    return base, est, final


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("logs", nargs="+", help="targets=path, e.g. 1=single.txt 2=cleave2.txt")
    ap.add_argument("--hand", default="", help="MainHand or OffHand only")
    ap.add_argument("--top", type=int, default=25)
    args = ap.parse_args()

    runs = []
    for spec in args.logs:
        n, path = spec.split("=", 1)
        if not os.path.exists(path):
            continue  # queued run not started yet
        base, est, final = read(path)
        if base:
            runs.append((int(n), base, est, final))
    runs.sort()

    keys = set()
    for _, _, est, _ in runs:
        keys |= set(est)
    if args.hand:
        keys = {k for k in keys if k[0] == args.hand}

    def pct(run, key):
        _, base, est, _ = run
        if key not in est:
            return None
        d, ci, n, _ = est[key]
        return 100 * d / base, 100 * ci / base, n

    # Sort by the first run's percentage (single target when given).
    def sort_key(key):
        v = pct(runs[0], key)
        return -(v[0] if v else -999)

    print("Baselines: " + " | ".join(f"{n} target{'s' if n > 1 else ''} {base:.0f} DPS" for n, base, _, _ in runs) + "\n")
    head = "| Weapon | Hand | Source | " + " | ".join(f"{n}T" for n, *_ in runs) + " |"
    print(head)
    print("|" + "---|" * (head.count("|") - 1))
    for key in sorted(keys, key=sort_key)[:args.top]:
        cells = []
        src = ""
        for run in runs:
            v = pct(run, key)
            src = src or run[2].get(key, (0, 0, 0, ""))[3]
            if v is None:
                cells.append("—")
            else:
                still_open = key not in run[3]
                cells.append(f"{v[0]:+.2f}% ±{v[1]:.2f}" + (" …" if still_open else ""))
        print(f"| {key[1]} | {SLOT_NAME[key[0]]} | {src} | " + " | ".join(cells) + " |")
    print("\n… = still running in that log")


if __name__ == "__main__":
    main()
