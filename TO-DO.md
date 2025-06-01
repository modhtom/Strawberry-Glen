To‑Do List

1. Pre‑Production & Documentation Game Design Document (GDD) --DONE--

Outline vision, core loops, target audience, art style, UI wireframes, and technical scope
Define win/lose conditions, progression milestones, and failure states Paper/Code Prototype
Rapidly mock up core interaction (walking, planting, harvesting).

2. Core Mechanics & Systems Tile & Collision System --DONE--

Finalize tile‐collision metadata (walkable, impassable, interactive)
and responses (e.g., “You can’t walk in water!” or slipping in mud).

Build a semantic map layer (object entities, crops, shop counters) separate from rendering layer.

Time & Progression System
Implement day/night cycle,
growth timers,
bake queues,
cow milking,
collecting eggs.

Game Loops
Farming loop: Plant → water → grow → harvest → store/sell.
Baking loop: Ingredient input → recipe process → finished goods → sell.
Economy & Progression: Cashflow, unlocks (new seeds, recipes, upgrades).

Interaction System “use” key (‘E’) for talking, picking, planting, opening shop.

3. Asset Pipeline Art & Animation --DONE--

Audio Integration ⚠️--NOT DONE YET--⚠️
Curate SFX for footsteps, planting, vending; BGM for day, night, shop ambience.
Set up audio manager (volume settings, mute toggle, layering).

Art
Fix main menu, Items textures, Desk, crops and Cows.

4. User Interface & UX HUD & Menus --DONE--

Inventory screen: item slots, counts, tooltips.  
Shop interface: buy/sell dialog, pricing, stock.
Pause/settings: volume sliders, keybindings, save/load.

5. Narrative & Worldbuilding Story & Goals ⚠️--NOT DONE YET--⚠️

Player backstory, long‑term objectives (e.g., expand farm, win bakery contest).

NPC profiles, routines, and failure/fallback states (missed harvest, debts).

6. Testing, Polish & Release ⚠️--NOT DONE YET--⚠️

Save/Load System
Persistent world state (tile changes, inventory, money, day count).
