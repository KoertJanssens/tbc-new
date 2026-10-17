package sim

import (
	"fmt"
	"sort"
	"testing"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
)

// Soft-reserve helper for Black Temple, same shape as sim/shamanhyjalsr_test.go:
// every BT drop an Enhancement Shaman can equip, simmed as a single-slot swap
// against Koerrii's current gear (the baseline from result_shaman.md).
//
// Loot table from assets/database/db.json zone 3959, filtered to cloth/leather/
// mail plus one-hand axe/mace/dagger/fist. Excluded on purpose:
//   - Warglaives of Azzinoth (32837/32838) - swords, shamans cannot use them.
//   - Twisted Blades of Zarak (32326) - rangedWeaponType 4 is Thrown; shaman
//     totems are rangedWeaponType 8, so BT drops no ranged upgrade at all.
//
//	SHAMAN_ITERATIONS=300 go test -tags with_db ./sim/ -run TestShamanBTSR -v

func btDrops() []hyjalDrop {
	return []hyjalDrop{
		// --- High Warlord Naj'entus ---
		{"Naj'entus", candidate{slotHands, "Fists of Mukoa", 32234, nil, 0, ""}},
		{"Naj'entus", candidate{slotShoulder, "Mantle of Darkness", 32377, nil, 0, ""}},
		{"Naj'entus", candidate{slotMainHand, "Rising Tide", 32236, nil, 0, "axe 2.6"}},
		{"Naj'entus", candidate{slotOffHand, "Rising Tide", 32236, nil, 0, "axe 2.6"}},
		// --- Supremus ---
		{"Supremus", candidate{slotWrist, "Wraps of Precise Flight", 32251, nil, 0, ""}},
		{"Supremus", candidate{slotChest, "Nether Shadow Tunic", 32252, nil, 0, ""}},
		{"Supremus", candidate{slotNeck, "Choker of Endless Nightmares", 32260, nil, 0, ""}},
		{"Supremus", candidate{slotFinger1, "Band of the Abyssal Lord", 32261, nil, 0, ""}},
		{"Supremus", candidate{slotFinger2, "Band of the Abyssal Lord", 32261, nil, 0, ""}},
		{"Supremus", candidate{slotMainHand, "The Brutalizer", 32254, nil, 0, "axe 1.6"}},
		{"Supremus", candidate{slotOffHand, "The Brutalizer", 32254, nil, 0, "axe 1.6"}},
		{"Supremus", candidate{slotMainHand, "Syphon of the Nathrezim", 32262, nil, 0, "mace 2.8"}},
		{"Supremus", candidate{slotOffHand, "Syphon of the Nathrezim", 32262, nil, 0, "mace 2.8"}},
		// --- Shade of Akama ---
		{"Shade of Akama", candidate{slotShoulder, "Shoulders of the Hidden Predator", 32264, nil, 0, ""}},
		{"Shade of Akama", candidate{slotWaist, "Shadow-walker's Cord", 32265, nil, 0, ""}},
		{"Shade of Akama", candidate{slotFinger1, "Ring of Deceitful Intent", 32266, nil, 0, ""}},
		{"Shade of Akama", candidate{slotFinger2, "Ring of Deceitful Intent", 32266, nil, 0, ""}},
		// --- Teron Gorefiend ---
		{"Teron Gorefiend", candidate{slotBack, "Shadowmoon Destroyer's Drape", 32323, nil, 0, ""}},
		{"Teron Gorefiend", candidate{slotWrist, "Insidious Bands", 32324, nil, 0, ""}},
		{"Teron Gorefiend", candidate{slotFeet, "Softstep Boots of Tracking", 32510, nil, 0, ""}},
		// --- Gurtogg Bloodboil ---
		{"Gurtogg Bloodboil", candidate{slotChest, "Vest of Mounting Assault", 32334, nil, 0, ""}},
		{"Gurtogg Bloodboil", candidate{slotFinger1, "Unstoppable Aggressor's Ring", 32335, nil, 0, ""}},
		{"Gurtogg Bloodboil", candidate{slotFinger2, "Unstoppable Aggressor's Ring", 32335, nil, 0, ""}},
		{"Gurtogg Bloodboil", candidate{slotMainHand, "Messenger of Fate", 32269, nil, 0, "dagger 1.4"}},
		{"Gurtogg Bloodboil", candidate{slotOffHand, "Messenger of Fate", 32269, nil, 0, "dagger 1.4"}},
		// --- Reliquary of the Lost ---
		{"Reliquary of the Lost", candidate{slotWaist, "Boneweave Girdle", 32346, nil, 0, ""}},
		{"Reliquary of the Lost", candidate{slotHands, "Grips of Damnation", 32347, nil, 0, ""}},
		{"Reliquary of the Lost", candidate{slotNeck, "Pendant of Titans", 32362, nil, 0, ""}},
		// --- Mother Shahraz ---
		{"Mother Shahraz", candidate{slotFeet, "Shadowmaster's Boots", 32366, nil, 0, ""}},
		// --- The Illidari Council ---
		{"Illidari Council", candidate{slotHead, "Forest Prowler's Helm", 32376, nil, 0, ""}},
		{"Illidari Council", candidate{slotTrinket1, "Madness of the Betrayer", 32505, nil, 0, ""}},
		{"Illidari Council", candidate{slotTrinket2, "Madness of the Betrayer", 32505, nil, 0, ""}},
		// --- Illidan Stormrage ---
		{"Illidan Stormrage", candidate{slotHead, "Cursed Vision of Sargeras", 32235, nil, 0, ""}},
		{"Illidan Stormrage", candidate{slotFinger1, "Stormrage Signet Ring", 32497, nil, 0, ""}},
		{"Illidan Stormrage", candidate{slotFinger2, "Stormrage Signet Ring", 32497, nil, 0, ""}},
		{"Illidan Stormrage", candidate{slotMainHand, "Shard of Azzinoth", 32471, nil, 0, "dagger 1.9"}},
		{"Illidan Stormrage", candidate{slotOffHand, "Shard of Azzinoth", 32471, nil, 0, "dagger 1.9"}},
		// --- Trash ---
		{"Trash", candidate{slotFinger1, "Band of Devastation", 32526, nil, 0, ""}},
		{"Trash", candidate{slotFinger2, "Band of Devastation", 32526, nil, 0, ""}},
		{"Trash", candidate{slotNeck, "Choker of Serrated Blades", 32591, nil, 0, ""}},
		{"Trash", candidate{slotFeet, "Treads of the Den Mother", 32593, nil, 0, ""}},
		{"Trash", candidate{slotMainHand, "Swiftsteel Bludgeon", 32943, nil, 0, "mace 1.5"}},
		{"Trash", candidate{slotOffHand, "Swiftsteel Bludgeon", 32943, nil, 0, "mace 1.5"}},
	}
}

func TestShamanBTSR(t *testing.T) {
	RegisterAll()
	iterations := iterationCount()
	encounter := core.MakeSingleTargetEncounter(5)

	run := func(gear *proto.EquipmentSpec) simResult {
		return runShaman(t, gear, talentsKoerriiPvE, encounter, iterations)
	}

	base := run(currentGear())
	fmt.Printf("ITERATIONS %d\n", iterations)
	fmt.Printf("BTBASE | current armory gear | DPS %.1f (stdev %.1f)\n", base.dps, base.stdev)

	type scoredDrop struct {
		d hyjalDrop
		r simResult
	}
	all := []scoredDrop{}
	for _, d := range btDrops() {
		if d.c.slot == slotOffHand && mainHandOnly[d.c.id] {
			continue
		}
		gear := cloneGear(currentGear())
		d.c.apply(gear)
		r := run(gear)
		all = append(all, scoredDrop{d, r})
		fmt.Printf("BT | %-21s | %-9s | %-32s | %6d | %8.1f | %+7.1f | %s\n",
			d.boss, slotNames[d.c.slot], d.c.label, d.c.id, r.dps, r.dps-base.dps, d.c.note)
	}

	sort.Slice(all, func(i, j int) bool { return all[i].r.dps > all[j].r.dps })
	fmt.Println("\n=== RANKED ===")
	for i, s := range all {
		fmt.Printf("BTRANK | %2d | %-21s | %-9s | %-32s | %6d | %+7.1f\n",
			i+1, s.d.boss, slotNames[s.d.c.slot], s.d.c.label, s.d.c.id, s.r.dps-base.dps)
	}
}
