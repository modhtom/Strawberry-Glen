# Strawberry Glen

**Strawberry Glen** is a cozy, quirky 2D farm‑bakery sim where you inherit an enchanted bakery‑farm, raise strawberries and cows, bake whimsical recipes, and uncover a floury conspiracy—all on your quest to earn 2 million coins and win the Grand Glen Bake‑Off.

---

## 📖 Table of Contents

1. [Game Vision & Narrative Hooks](#game-vision--narrative-hooks)
2. [Core Loop & Paper Prototype](#core-loop--paper-prototype)
3. [Systems & Mechanics](#systems--mechanics)
   - Recipe Discovery
   - Cow Personalities
   - Shop Upgrades
   - Baking Contest

---

## 🎨 Game Vision & Narrative Hooks

**Working Title:** _Strawberry Glen_

**Core Fantasy:** You’re the reluctant heir to Aunt Marmalade’s magical bakery‑farm. She vanished into the “dough dimension,” leaving behind a moody cow, sentient strawberries, and debt to a mushroom syndicate. Plant, bake, and barter your way to 2M coins, win the Bake‑Off, and unearth the bakery’s secrets.

**Tone:** Wholesome but weird

**Win Condition:** Earn 2 million coins, win the Grand Glen Bake‑Off, and unlock a secret ending by maxing out cow friendship.

---

## 🧪 Core Loop & Paper Prototype

Before art or polish, prove the loop is fun. Prototype in a single scene:

1. **Farm Loop:** Plant → Water → Grow (3 days) → Harvest → Inventory
2. **Cow Loop:** Feed daily → Milk (2‑day cycle) → Miss feed = grumpy cow
3. **Baking Loop:** Combine 2–3 ingredients at oven → Wait 5 s → Get product
4. **Sales Loop:** Place on shop counter → Auto‑sell → Update coins & UI
5. **Time Loop:** 1 s = 1 game hour, 12 ticks = new day, trigger growth & cow reset

**Prototype Checklist:**

- [ ] 10×10 field movement & tile metadata
- [ ] Planting, watering, harvesting logic
- [ ] Cow entity with feed/milk interaction
- [ ] Baking UI & combination logic
- [ ] Shop counter auto‑sell
- [ ] Day/night transitions
- [ ] Coin & inventory HUD + basic SFX

---

## 🛠️ Systems & Mechanics

### 🍰 Recipe Discovery System

- **Mechanic:** Use 2–3 ingredients in the oven
- **Known recipes:** Starter combos loaded by default
- **Discoverable:** Unrecognized combos generate funny names and log as new recipes
- **Examples:**
  - 🍓+🌾 → Tart
  - 🥛+🌾 → Cream Bun
  - 🍓+🥛+🌾 → Strawberry Milk Cake

### 🐮 Cow Personalities

Each cow has a **name**, **mood**, and **quirk**. Mood affects milk yield & dialogue.

- **Moozart:** Grumpy until serenaded (jazz unlocks bonus milk)
- **Cowculus:** Philosophical in rain, produces flavored milk if happy
- **Dairyssa:** Dramatic, requires baked gift to boost mood

### 🏪 Shop Upgrades

Spend coins to improve shop:

| Upgrade            | Cost | Effect                             |
| ------------------ | ---- | ---------------------------------- |
| Bigger Counter     | 500  | Hold 3 items at once               |
| Cute Signage       | 750  | +10% sell price                    |
| Ghost Acceptance™ | 1000 | Ghost customers arrive at night    |
| Loyalty Card       | 1500 | Regular customers daily            |
| Orders Board       | 2000 | Daily quest orders → bonus payouts |

### 🏆 Baking Contest

- Occurs every few days in town square
- Submit one item; judged on Taste, Presentation, Creativity
- **Judges:**
  - Baron Flan (fancy names)
  - Loafilda (vegan)
  - Glump the Ogre (burnt goods fan)
  - Sir Meowster (only pink treats)
- **Rewards:** Coins, Fame (↑shop traffic), Exclusive ingredients
