# 🍓 Strawberry Glen

A whimsical 2D farm-bakery simulator where you inherit a magical farm, grow enchanted strawberries, bake quirky recipes, and unravel doughy mysteries. Earn coins, befriend moody cows, and win the Grand Glen Bake-Off!

![Gameplay Preview](assets/preview.png)

---

## 🌟 Features

- **Farming & Crafting**: Plant seeds, harvest crops, and bake magical goods like _Strawberry Milk Cake_ and _Experimental Jam_.
- **Cow Companions**: Manage cows with unique personalities—grumpy Moozart loves jazz, while Cowculus philosophizes in the rain.
- **Dynamic Shop**: Buy/sell items, upgrade your store, and attract ghost customers with _Ghost Acceptance™_ upgrades.
- **Inventory System**: Manage tools, seeds, and baked goods across 5 slots. Reusable tools and stackable items supported!
- **Quests & Secrets**: Collect _Eldermint Leaves_ and craft the _Cow Flute_ to calm angry cows and uncover hidden lore.
- **Whimsical Tone**: Burnt pies, sentient strawberries, and a mushroom syndicate debt—charm meets chaos!

---

## 🛠️ Installation

1. **Prerequisites**:

   - [Go](https://golang.org/dl/) (1.16+)
   - [Raylib](https://www.raylib.com/) (C library) and [raylib-go](https://github.com/gen2brain/raylib-go) bindings.

2. **Run the Game**:
   ```bash
   git clone https://github.com/yourusername/strawberry-glen.git
   cd strawberry-glen
   go run main.go inventory.go shop.go
   ```

---

## 🥧 How to Play

1. **Farming**:
   - Collect seeds (e.g., `Wheat Seeds`).
   - Plant, water, and harvest crops over 3 days.
2. **Baking**:

   - Combine ingredients (e.g., `Strawberry + Milk + Wheat`) at the oven.
   - Sell baked goods in your shop for coins!

3. **Cow Care**:

   - Feed cows daily to collect milk.
   - Use the _Cow Flute_ to calm grumpy cows.

4. **Shop Management**:

   - Buy low, sell high! Restock occurs every 5 minutes.
   - Unlock upgrades like _Loyalty Cards_ and _Ghost Acceptance™_.

5. **Quests**:
   - Complete tasks (e.g., gather _Eldermint Leaves_) to progress the story.

---

## 🎮 Master the Keys - Control Cheat Sheet

**Movement & Exploration**  
🕹️ `WASD`/`Arrow Keys` – Move your character  
🔍 `Z`/`X` – Zoom in/out to see details or the big picture

**Quick Actions**  
🎒 `I` – Open/close your **inventory** (manage seeds, tools, and goodies)  
🏪 `B` – Toggle the **shop** (buy low, sell high!)  
🏪 `E` – Interact with things
⏸️ `ESC` – Pause game to access settings or save progress

**Inventory Management**  
🖱️ `Mouse Hover` – Preview item details in your inventory  
✅ `Enter`/`Left Click` – Use selected item (plant seeds, water crops, play the cow flute!)  
⬅️➡️ `Arrow Keys` – Navigate inventory slots when menu is open

---

## 📦 Your Pocket Universe - Item System Explained

### 🛠️ Tools (Reusable forever!)

⚒️ **Watering Can (ID:3)** - Hydrate crops daily  
⚒️ **Hoe (ID:4)** - Till soil for planting  
⚒️ **Axe (ID:5)** - Clear obstacles

### 🌱 Seeds (Plant in tilled soil)

🌾 **Wheat Seeds (ID:1)** - Grows in 3 days  
🌾 **Strawberry Seeds (ID:2)** - Sweet profits!

### 🍓 Crops (Harvest to bake/sell)

🧺 **Wheat (ID:10)** - Base for bread  
🧺 **Strawberry (ID:11)** - For tarts and cakes

### 🥧 Baked Goods (Sell for $$$)

🧁 **Bread (ID:30)** - Basic but reliable  
🧁 **Strawberry Tart (ID:31)** - Customer favorite  
🧁 **Burnt Pie (ID:33)** - Oops! Still sells to Glump the Ogre

### 🐄 Special Items (Unlock secrets!)

🎵 **Cow Flute (ID:41)** - Calms angry cows instantly  
🍃 **Eldermint Leaves (ID:40)** - Quest item for Berry's tea

---

**Pro Tips:**  
🔸 **Stack smart**: Seeds/crops stack in inventory (e.g. 5 Wheat Seeds = 1 slot)  
🔹 **Experiment**: Try combos like 🍓+🥛+🌾 in oven for Strawberry Milk Cake!  
🔻 **Quick sell**: "Experimental Jam (ID:34)" vanishes after use - sell it fast!

---

## 🖥️ Technical Details

- **Engine**: Built with Go and Raylib for 2D rendering.
- **Code Structure**:
  - `main.go`: Core gameplay loop, rendering, and input handling.
  - `inventory.go`: Inventory management and UI.
  - `shop.go`: Buy/sell logic and shop interface.
- **Assets**: Textures for items, tiles, and characters in `assets/`.

---

## 🙌 Credits

- **Assets**:
- **Music**:
- **Inspiration**: Stardew Valley + Animal Crossing

---

## 📜 License

MIT License. See [LICENSE](LICENSE) for details.
