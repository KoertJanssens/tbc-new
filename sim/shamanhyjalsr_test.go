package sim

import (
	"fmt"
	"sort"
	"testing"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
)

// Soft-reserve helper for Mount Hyjal: sims every Hyjal drop an Enhancement
// Shaman can actually use as a single-slot swap against Koerrii's current gear
// (same baseline as sim/shamanupgrade_test.go and result_shaman.md), so the
// reserve list can be ranked on measured DPS rather than item level.
//
// Loot table derived from assets/database/db.json, zone 3606, filtered to
// cloth/leather/mail plus one-hand axe/mace/dagger/fist. Items with no melee
// stats at all (healer/caster pieces) are skipped - they cannot be upgrades.
//
//	SHAMAN_ITERATIONS=300 go test -tags with_db ./sim/ -run TestShamanHyjalSR -v

// handType 1 in assets/database/db.json = main hand only.
var mainHandOnly = map[int32]bool{30918: true, 32946: true}

type hyjalDrop struct {
	boss string
	c    candidate
}

func hyjalDrops() []hyjalDrop {
	return []hyjalDrop{
		// --- Rage Winterchill ---
		{"Rage Winterchill", candidate{slotWrist, "Deadly Cuffs", 30863, nil, 0, ""}},
		{"Rage Winterchill", candidate{slotWrist, "Bracers of the Pathfinder", 30864, nil, 0, ""}},
		{"Rage Winterchill", candidate{slotMainHand, "Tracker's Blade", 30865, nil, 0, "dagger 1.5"}},
		{"Rage Winterchill", candidate{slotOffHand, "Tracker's Blade", 30865, nil, 0, "dagger 1.5"}},
		// --- Anetheron ---
		{"Anetheron", candidate{slotWaist, "Don Alejandro's Money Belt", 30879, nil, 0, ""}},
		{"Anetheron", candidate{slotFeet, "Quickstrider Moccasins", 30880, nil, 0, ""}},
		// --- Kaz'rogal ---
		{"Kaz'rogal", candidate{slotFeet, "Black Featherlight Boots", 30891, nil, 0, ""}},
		// Beast-tamer's Shoulders (30892) is deliberately omitted: it is a hunter pet
		// item, and sim/hunter/item_sets.go:243 registers its effect with an
		// unconditional agent.(hunter.HunterAgent) cast, so equipping it on any other
		// class panics the engine. Its stats (agi/AP, no crit or hit) lose to
		// Razorfury Mantle for a shaman regardless.
		{"Kaz'rogal", candidate{slotShoulder, "Razorfury Mantle", 30917, nil, 0, ""}},
		{"Kaz'rogal", candidate{slotWaist, "Valestalker Girdle", 30919, nil, 0, ""}},
		{"Kaz'rogal", candidate{slotMainHand, "Hammer of Atonement", 30918, nil, 0, "mace 1.8"}},
		// --- Azgalor ---
		{"Azgalor", candidate{slotLegs, "Shady Dealer's Pantaloons", 30898, nil, 0, ""}},
		{"Azgalor", candidate{slotLegs, "Bow-stitched Leggings", 30900, nil, 0, ""}},
		{"Azgalor", candidate{slotMainHand, "Boundless Agony", 30901, nil, 0, "dagger 1.8"}},
		{"Azgalor", candidate{slotOffHand, "Boundless Agony", 30901, nil, 0, "dagger 1.8"}},
		// --- Archimonde ---
		{"Archimonde", candidate{slotChest, "Midnight Chestguard", 30905, nil, 0, ""}},
		{"Archimonde", candidate{slotChest, "Mail of Fevered Pursuit", 30907, nil, 0, ""}},
		// --- Trash ---
		{"Trash", candidate{slotNeck, "Choker of Serrated Blades", 32591, nil, 0, ""}},
		{"Trash", candidate{slotChest, "Chestguard of Relentless Storms", 32592, nil, 0, ""}},
		{"Trash", candidate{slotMainHand, "Claw of Molten Fury", 32946, nil, 0, "fist 2.7"}},
	}
}

func TestShamanHyjalSR(t *testing.T) {
	RegisterAll()
	iterations := iterationCount()
	encounter := core.MakeSingleTargetEncounter(5)

	run := func(gear *proto.EquipmentSpec) simResult {
		return runShaman(t, gear, talentsKoerriiPvE, encounter, iterations)
	}

	base := run(currentGear())
	fmt.Printf("ITERATIONS %d\n", iterations)
	fmt.Printf("HYJALBASE | current armory gear | DPS %.1f (stdev %.1f)\n", base.dps, base.stdev)

	type scoredDrop struct {
		d hyjalDrop
		r simResult
	}
	all := []scoredDrop{}
	for _, d := range hyjalDrops() {
		// Main-hand-only weapons cannot legally go in the off hand; the engine does
		// not enforce it and would report a nonsense number for an empty off hand.
		if d.c.slot == slotOffHand && mainHandOnly[d.c.id] {
			continue
		}
		gear := cloneGear(currentGear())
		d.c.apply(gear)
		r := run(gear)
		all = append(all, scoredDrop{d, r})
		fmt.Printf("HYJAL | %-16s | %-9s | %-32s | %6d | %8.1f | %+7.1f | %s\n",
			d.boss, slotNames[d.c.slot], d.c.label, d.c.id, r.dps, r.dps-base.dps, d.c.note)
	}

	sort.Slice(all, func(i, j int) bool { return all[i].r.dps > all[j].r.dps })
	fmt.Println("\n=== RANKED ===")
	for i, s := range all {
		fmt.Printf("HYJALRANK | %2d | %-16s | %-9s | %-32s | %6d | %+7.1f\n",
			i+1, s.d.boss, slotNames[s.d.c.slot], s.d.c.label, s.d.c.id, s.r.dps-base.dps)
	}
}
