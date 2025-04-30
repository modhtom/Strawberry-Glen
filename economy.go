package main

type Progression struct {
	TotalGoldEarned  int
	DaysPlayed       int
	CropsHarvested   int
	ItemsSold        int
	UnlockedSeeds    map[int]bool
	UnlockedRecipes  map[int]bool
	UnlockedUpgrades map[string]bool
	ActiveUpgrades   map[string]bool
}

var progression = Progression{
	UnlockedSeeds:    make(map[int]bool),
	UnlockedRecipes:  make(map[int]bool),
	UnlockedUpgrades: make(map[string]bool),
	ActiveUpgrades:   make(map[string]bool),
}

var unlockables = []struct {
	Condition func() bool
	Unlock    func()
	Message   string
}{
	{
		Condition: func() bool { return progression.TotalGoldEarned >= 500 },
		Unlock:    func() { progression.UnlockedSeeds[2] = true },
		Message:   "Unlocked Strawberry Seeds!",
	},
	{
		Condition: func() bool { return progression.CropsHarvested >= 20 },
		Unlock:    func() { progression.UnlockedUpgrades["Fertilizer"] = true },
		Message:   "Unlocked Fertilizer (Growth Speed +25%%)",
	},
	{
		Condition: func() bool { return progression.TotalGoldEarned >= 1500 },
		Unlock:    func() { progression.UnlockedRecipes[32] = true },
		Message:   "Unlocked Strawberry Milk Cake Recipe!",
	},
	{
		Condition: func() bool { return progression.DaysPlayed >= 7 },
		Unlock:    func() { progression.UnlockedUpgrades["LargeBackpack"] = true },
		Message:   "Unlocked Large Backpack (Inventory +2 Slots!)",
	},
	{
		Condition: func() bool { return progression.TotalGoldEarned >= 5000 },
		Unlock: func() {
			progression.UnlockedSeeds[40] = true
			progression.UnlockedRecipes[41] = true
		},
		Message: "Unlocked Eldermint Seeds and Mystic Pie Recipe!",
	},
}

func checkProgression() {
	for _, u := range unlockables {
		if u.Condition() && !progression.UnlockedUpgrades[u.Message] {
			u.Unlock()
			progression.UnlockedUpgrades[u.Message] = true
			showMessages(u.Message, 3.0)
		}
	}
}
