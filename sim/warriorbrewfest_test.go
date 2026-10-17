package sim

import (
	"fmt"
	"testing"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
)

// Brewfest trinket check for Koerry (Fury Warrior), on the 5/5 Destroyer gear
// from sim/warriorupgrade_test.go. Empty Mug of Direbrew (281739) is the
// Anniversary ilvl 141 version from Coren Direbrew: +98 AP, use +376 AP for 20s
// on a 2 min cooldown (sim/common/tbc/stat_bonus_cds_auto_gen.go). It only exists
// in the item database after the October 2026 origin/master merge.
//
//	WARRIOR_ITERATIONS=300 go test -tags with_db ./sim/ -run TestWarriorBrewfestTrinket -v

const (
	itemBadgeOfTheSwarmguard = 21670
	itemDragonspineTrophy    = 28830
	itemEmptyMugOfDirebrew   = 281739
	itemMadnessOfTheBetrayer = 32505
)

func TestWarriorBrewfestTrinket(t *testing.T) {
	RegisterAll()
	iterations := warriorIterationCount()
	encounter := core.MakeSingleTargetEncounter(5)
	warriorBaseGear = warriorT5Gear

	run := func(gear *proto.EquipmentSpec) simResult {
		return runWarrior(t, gear, talentsFuryStandard, "fury", encounter, iterations)
	}
	trinkets := func(t1, t2 int32) *proto.EquipmentSpec {
		g := cloneGear(warriorT5Gear())
		g.Items[slotTrinket1] = itemSpec(t1, 0)
		g.Items[slotTrinket2] = itemSpec(t2, 0)
		return g
	}

	base := run(warriorT5Gear())
	fmt.Printf("ITERATIONS %d\n", iterations)
	fmt.Printf("BREWBASE | Badge of the Swarmguard + Dragonspine Trophy (5/5 T5) | DPS %.1f (stdev %.1f)\n", base.dps, base.stdev)

	for _, c := range []struct {
		label  string
		t1, t2 int32
	}{
		{"Empty Mug of Direbrew + Dragonspine Trophy", itemEmptyMugOfDirebrew, itemDragonspineTrophy},
		{"Badge of the Swarmguard + Empty Mug of Direbrew", itemBadgeOfTheSwarmguard, itemEmptyMugOfDirebrew},
		{"Madness of the Betrayer + Dragonspine Trophy", itemMadnessOfTheBetrayer, itemDragonspineTrophy},
		{"Madness of the Betrayer + Empty Mug of Direbrew", itemMadnessOfTheBetrayer, itemEmptyMugOfDirebrew},
		{"Bloodlust Brooch + Dragonspine Trophy", itemBloodlustBrooch, itemDragonspineTrophy},
		{"Empty Mug of Direbrew + Bloodlust Brooch", itemEmptyMugOfDirebrew, itemBloodlustBrooch},
	} {
		r := run(trinkets(c.t1, c.t2))
		fmt.Printf("BREW | %-48s | %8.1f | %+7.1f\n", c.label, r.dps, r.dps-base.dps)
	}
}
