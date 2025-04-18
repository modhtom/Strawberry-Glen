To‑Do List

1. Pre‑Production & Documentation Game Design Document (GDD)

Outline vision, core loops, target audience, art style, UI wireframes, and technical scope
Define win/lose conditions, progression milestones, and failure states Paper/Code Prototype
Rapidly mock up core interaction (walking, planting, harvesting).

------------------DONE⬆️----------------------------------

2. Core Mechanics & Systems Tile & Collision System

Finalize tile‐collision metadata (walkable, impassable, interactive)
and responses (e.g., “You can’t walk in water!” or slipping in mud).

Build a semantic map layer (object entities, crops, shop counters) separate from rendering layer.

Time & Progression System Implement day/night cycle, growth timers, bake queues,
and scheduling for NPC routines.

Game Loops
Farming loop: Plant → water → grow → harvest → store/sell.
Baking loop: Ingredient input → recipe process → finished goods → sell.
Economy & Progression: Cashflow, unlocks (new seeds, recipes, upgrades).

Interaction System
Contextual “use” key (e.g., ‘E’) for talking, picking, planting, opening shop.

3. Asset Pipeline Art & Animation

Audio Integration
Curate SFX for footsteps, planting, vending; BGM for day, night, shop ambience.
Set up audio manager (volume settings, mute toggle, layering).

4. User Interface & UX HUD & Menus

Inventory screen: item slots, counts, tooltips.
Shop interface: buy/sell dialog, pricing, stock.
Pause/settings: volume sliders, keybindings, save/load.

Dialogue & NPC Interaction Dialogue boxes, choice prompts, simple quest log.

5. Narrative & Worldbuilding Story & Goals

Player backstory, long‑term objectives (e.g., expand farm, win bakery contest).

NPC profiles, routines, and failure/fallback states (missed harvest, debts).

6. Testing, Polish & Release

Save/Load System
Persistent world state (tile changes, inventory, money, day count).

Playtesting & QA
Internal alpha tests for bug‑finding, balance tuning, and loop tuning.
External playtests for usability, fun factor, and narrative clarity.

Performance & Optimization
Profile frame rate, memory usage; optimize sprite batching and map culling.

Polish & Feedback
Add VFX (particle effects for harvest), SFX layering, UI animations.
Iterate based on analytics or direct tester feedback.

7. Deployment & Post‑Launch

Build Pipeline & Distribution
Scripts for packaging assets/executables, platform builds.

Marketing & Community
Trailer, screenshots, dev blog updates, play‑throughlets.

Updates & Roadmap
Plan DLC/content updates, seasonal events, player feedback channels.
