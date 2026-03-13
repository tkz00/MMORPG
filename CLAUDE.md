# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

An MMORPG framework with a **Unity client** (C#) and a **Go backend**. Real-time multiplayer via WebSockets with backend-authoritative game logic. Character state is persisted in PostgreSQL.

## Commands

### Running the game
```bash
cp .env.example .env        # first time setup
make db-up                  # start PostgreSQL via Docker
make run                    # start Go backend
# Then open Unity project and press Play
```

### Backend (run from repo root)
```bash
go run -C backend cmd/server/main.go   # start server
go test -C backend ./test/integration_test/...   # integration tests (requires Docker)
```

Integration tests spin up Docker automatically, start the full server stack, and tear down when done.

### Environment variables
Copy `.env.example` to `.env` and fill in PostgreSQL credentials. Tests auto-load `.env` if present, otherwise expect env vars set externally (CI).

## Architecture

### Backend servers (three run concurrently)
- **WebSocket game server** — `connection/native_server.go`, port **3009**, handles real-time gameplay
- **Configurator REST API** — `pkg/configurator/`, port **8080**, manages ability definitions
- **HTTP REST API** — `pkg/handlers/`, port **8081**, exposes character data (Gin)

### Game loop
`config/config.go` sets the tick rate at **50ms**. Each tick:
1. `connection/native_server.go` fires a `tick` channel with `deltaTime`
2. `pkg/game/gameplay/game_loop.go::UpdateState()` processes all players, projectiles, AoEs, spawners, NPCs, and delayed mechanics
3. A diff of the updated `GameState` is broadcast to all connected WebSocket clients

### Core packages
| Path | Purpose |
|---|---|
| `pkg/game/entities/` | Domain model: `GameState`, `Character`, `Projectile`, `AoE`, `Mechanic`, etc. |
| `pkg/game/gameplay/game_loop.go` | Per-tick update logic and player connection lifecycle |
| `pkg/game/repository/` | PostgreSQL persistence via GORM |
| `pkg/configurator/` | Ability CRUD backed by `abilities.json` and `playersInitialAbilities.json` |
| `api/dtos/` | WebSocket message DTOs and serialization (shared by game server and client protocol) |
| `connection/` | WebSocket server with channel-based concurrency for client add/remove/tick |

### Ability / Mechanic system
Abilities are stored in `backend/abilities.json` and loaded at startup. Each ability contains a list of **Mechanics** — composable effects like `damage`, `heal`, `create_projectile`, `create_AoE`, `delay`, `buff_stat`. Mechanics are registered in `pkg/game/entities/game_state.go::StartGameState()` and implemented in `pkg/game/entities/mechanics.go`.

The configurator server exposes a REST API used by the Unity editor window (`Window > Abilities Editor Window`) to create, edit, and delete abilities without restarting the server.

**Hard constraint**: The first ability (`id: "1"`, Sword Slash) is used by skeleton NPCs and must not be deleted. Players can have at most 4 abilities assigned at once.

### Client (Unity 2022.3.53f1)
- `Assets/Scripts/GameManager.cs` — singleton that handles all incoming WS messages, instantiates/destroys game objects, and delegates to subsystems
- `Assets/Scripts/Services/WebSocketConnection.cs` — WebSocket connection management with typed message handlers
- `Assets/Scripts/Dtos/` — C# DTOs mirroring the Go `api/dtos/` package (must stay in sync)
- `Assets/Scripts/Player/` — player-specific scripts added dynamically to the main player's GameObject

### WebSocket message protocol
Messages are JSON with shape `{ "actionType": string, "body": object }`.

**Client → Server:** `position`, `ability_cast`, `respawn`, `use_item`, `equip_item`, `unequip_item`

**Server → Client:** `Player` (on connect), `GameState` (every tick diff), `respawn`

### Persistence
Characters are saved to PostgreSQL every 5 seconds via a background goroutine in `gameplay/game_loop.go::saveGameState()`. New characters are created on first connection using their `character` query parameter as the name/ID.

### Diff-based game state updates
Inventory diffs are already implemented; full game state currently sends all data each tick. The pattern in `api/dtos/mapper.go` is the place to implement further diff logic.
