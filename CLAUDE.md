# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Real-time multi-player Jeopardy game with over 100,000 questions. This is a full-stack application with a Go backend (WebSocket server) and Angular frontend.

Live deployment: https://playjeopardy.netlify.app

## Architecture

### Backend (be-jeopardy)
- **Language**: Go
- **Framework**: Gin HTTP router with WebSocket support
- **Database**: PostgreSQL (via Docker Compose)
- **Authentication**: JWT-based auth with Supabase integration

**Key packages** (in `internal/`):
- `jeopardy/`: Core game logic, player management, bot behavior, question handling, analytics
- `handlers/`: HTTP route handlers for API endpoints
- `socket/`: WebSocket connection management for real-time gameplay
- `db/`: Database queries and connection management
- `auth/`: JWT token validation and authentication
- `logic/`: Game flow and state management

The server runs a background goroutine (1-hour ticker) that cleans up inactive games.

### Frontend (fe-jeopardy)
- **Framework**: Angular 17.3
- **Styling**: LESS
- **State Management**: RxJS services
- **Authentication**: Supabase (Google, GitHub, email)

**Key directories** (in `src/app/`):
- `game/`: Main game UI components and gameplay screens
- `services/`: Core services including `game-state.service.ts` (game state management), `websocket.service.ts` (server communication), `api.service.ts` (HTTP requests), `auth.service.ts` (authentication)
- `auth/`, `profile/`, `analytics/`, `leaderboards/`: Feature modules
- `model/`: TypeScript interfaces for game entities

## Initial Setup (Fresh Environment)

### Prerequisites
- Go 1.25+ (`brew install go` on macOS)
- Node.js 18+ (currently using v23.8.0)
- PNPM (`npm install -g pnpm`)
- Docker Desktop (for PostgreSQL)
- PM2 for process management (`npm install -g pm2`)
- air for Go hot reload (`go install github.com/air-verse/air@latest`)

### Backend Setup

1. **Install dependencies:**
   ```bash
   cd be-jeopardy
   go mod tidy
   ```

2. **Generate JWT keys:**
   ```bash
   mkdir -p .keys
   ./gen_keys.sh
   ```
   This creates RSA key pair in `.keys/jwtRS512.key` and `.keys/jwtRS512.key.pub`

3. **Create .env file:**
   ```bash
   cat > .env << 'EOF'
   PORT=8080
   GIN_MODE=debug
   DATABASE_URL=postgresql://postgres:postgres@localhost:5434/postgres?sslmode=disable
   ALLOW_ORIGIN=http://localhost:4200
   EOF
   ```

4. **Start PostgreSQL:**
   ```bash
   docker compose up -d postgres
   ```
   PostgreSQL runs on port 5434 (mapped from container's 5432)

### Frontend Setup

1. **Install dependencies:**
   ```bash
   cd fe-jeopardy
   pnpm install
   ```

### Loading Questions into Database

After PostgreSQL is running, load questions using Python scripts with uv:

**1. Set up Python environment with uv:**
```bash
brew install uv                    # Install uv package manager
cd be-jeopardy
uv venv                            # Create virtual environment in .venv
uv pip install requests beautifulsoup4 psycopg2-binary
```

**2. Modify scraper for quick testing** (optional):
For testing, edit `scrapers/jeopardy.py` to limit games:
```python
games = 1
MAX_GAMES = 10  # Add this line
for row in rows:
    if games > MAX_GAMES:  # Add this line
        break              # Add this line
    print(f'{games}/{MAX_GAMES}')  # Update this line
```

**3. Run scraper:**
```bash
mkdir -p clues
cd scrapers
../.venv/bin/python jeopardy.py    # Scrapes j-archive.com
# Alternative sources:
# ../.venv/bin/python jetpunk.py
# ../.venv/bin/python opentdb.py
```

**4. Fix database port in scripts:**
The Python scripts default to port 5432, but docker-compose uses 5434. Update both:
- `insert_clues.py` line 15: change `port='5432'` to `port='5434'`
- `add_alternatives.py` line 8: change `port='5432'` to `port='5434'`

**5. Load questions into database:**
```bash
cd be-jeopardy
mv scrapers/season40.tsv clues/     # Move scraped data
.venv/bin/python insert_clues.py    # Creates tables and inserts questions
```

**6. Process answer alternatives:**
```bash
.venv/bin/python add_alternatives.py  # Adds variations (removes "the", handles parentheses, etc.)
```

**Result:** Database populated with ~600 questions (10 games) or more depending on scraper settings.

## Development Commands

### Running the App

**Development mode with hot reloading (recommended):**
```bash
# From project root
cd be-jeopardy && docker compose up -d postgres && cd ..
pm2 start ecosystem.dev.config.js

# View logs
pm2 logs

# Stop
pm2 stop all
pm2 delete all
```

**Production mode (no hot reloading):**
```bash
pm2 start ecosystem.config.js
```

**Manual (separate terminals):**
```bash
# Terminal 1 - Backend
cd be-jeopardy
source .env && make run

# Terminal 2 - Frontend
cd fe-jeopardy
pnpm run ng serve
```

### Backend
```bash
cd be-jeopardy
go mod tidy                      # Install dependencies
docker compose up -d postgres    # Start PostgreSQL database
make build                       # Build server binary
make run                         # Build and run server
make clean                       # Clean build artifacts
go test ./...                    # Run tests
```

### Frontend
```bash
cd fe-jeopardy
pnpm install                     # Install dependencies
pnpm run ng serve                # Run dev server with hot reload
pnpm run ng build                # Production build
pnpm run ng test                 # Run Karma/Jasmine tests
```

## Critical Architecture Details

### In-Memory Game State
All active games are stored in memory (not database) for low latency:
- `privateGames map[string]*Game` - Private games by name
- `publicGames map[string]*Game` - Matchmaking pool
- `playerGames map[string]*Game` - Player ID → Game lookup
- Games are cleaned up hourly if inactive
- **Important**: Server restart loses all active games (acceptable tradeoff for performance)

### Channel-Based Event System
Each game uses Go channels for concurrent event handling:
- `msgChan` - Player actions (pick, buzz, answer, wager)
- `disconnectChan` - Handle player disconnections
- `restartChan` - Trigger game restart
- `chatChan` - Broadcast chat messages
- `reactChan` - Broadcast emoji reactions

All game state modifications must happen in the main game goroutine to avoid race conditions.

### WebSocket Communication
- **Three separate connections**: Game state, chat, reactions (separation of concerns)
- `SafeConn` wrapper prevents concurrent write panics (critical - wraps `websocket.Conn` with mutex)
- Ping every 50 seconds to prevent connection drops
- Backend: `internal/socket/socket.go`
- Frontend: `services/websocket.service.ts`

### JWT Authentication Flow
- Backend reads keys from `.keys/jwtRS512.key` (local) OR environment variables (production)
- Auth code in `internal/auth/auth.go` tries env vars first, falls back to file reading
- Generated on game join, validated on every WebSocket connection and HTTP request
- Contains player ID in `sub` claim
- RS512 algorithm (RSA-based, more secure than HS256)

### Game State Machine
Frontend state transitions (`model/model.ts:GameState`):
```
PreGame → BoardIntro → RecvPick ↔ RecvBuzz ↔ RecvAns ↔ RecvWager
                                     ↓
                               RecvDispute
                                     ↓
                                PostGame
```

Each state corresponds to a component in `game/` directory. State updates are broadcast via RxJS Subject in `game-state.service.ts`. Components subscribe to state changes and render accordingly.

### Bot Implementation
Bots (`internal/jeopardy/bot.go`):
- Implement `GamePlayer` interface (same as human players)
- Always know correct answers (read from `g.CurQuestion.Answer`)
- Add realistic delays: 3s for picking, 5s for answering, 10s for Daily Doubles
- Use game-theory optimal wagering strategies (e.g., wager enough to beat leader)
- Names: Eager Eagle, Golden Gorilla, Sharp Shark, Smart Snake, Tough Tiger

### Player Disconnection Handling
- Game auto-pauses when player disconnects
- Player can rejoin via same JWT token
- Disconnected players can be replaced with bots
- If all players disconnect, game is removed from memory after timeout

## Answer Validation Logic

The game uses **fuzzy matching** via Levenshtein distance (`internal/jeopardy/question.go`):

- Case-insensitive comparison
- Tolerance scales with answer length:
  - ≤5 chars: exact match
  - 5-7 chars: 1 character difference allowed
  - 7-9 chars: 2 characters
  - 9-12 chars: 3 characters
  - 12-15 chars: 4 characters
  - 15+ chars: 5 characters

Answers from j-archive.com include variations:
- `(Harry) Houdini` → parentheses mean "Harry" is optional
- `Gordie Howe (or Bobby Hull)` → multiple acceptable answers

The `add_alternatives.py` script processes these patterns into the database:
1. Accent normalization (`Café → Cafe`)
2. Article removal (`The Beatles → Beatles`)
3. Parentheses handling (generates both forms)
4. Dynamic learning: Disputed answers added to alternatives

## Database Schema

**jeopardy_clues** (100,000+ rows):
- `round int` - 1=Jeopardy, 2=Double Jeopardy, 3=Final
- `clue_value int` - Point value (200/400/600/800/1000 or 400/800/...)
- `category text` - Category name
- `question text` - The answer (what contestants see)
- `answer text` - The question (correct response)
- `alternatives text[]` - Valid answer variations
- `incorrect text[]` - Dynamically added wrong answers
- `air_date text` - Original episode date

**Question Selection** (`internal/db/sql/get_questions.sql`):
- Randomly picks 6 categories per round (must have exactly 5 questions each)
- Daily Doubles use weighted random based on real Jeopardy statistics (see `daily_double_occurrence_bounds.sql`)
- Custom categories: Players can search and select specific categories

**Analytics Tables**:
- `jeopardy_analytics` - Game-level analytics (JSONB for flexibility)
- `player_games` - Per-player cumulative stats (email → wins, points, accuracy)

## Real-time Game Flow

1. Players create/join games via HTTP API (`handlers/handlers.go`)
2. WebSocket connection established for real-time updates
3. Game state synchronized through `game-state.service.ts` on frontend
4. Backend manages game logic, question selection, scoring, and timers in `internal/jeopardy/game.go`
5. Chat messages and emoji reactions broadcast via separate WebSocket connections
6. Player analytics and leaderboards updated in PostgreSQL at game end

## Important Notes

- **PNPM required**: Frontend uses PNPM (not NPM) per global config
- **uv required**: Backend Python scripts use uv for package management (not pip/virtualenv)
- **Multi-line env vars**: JWT keys are read from files (not env vars) for local development
- **air location**: May be installed in `~/go/bin/air` if not in PATH
- **Node types**: Frontend requires `@types/node` and `"types": ["node"]` in tsconfig for `NodeJS.Timeout`
- **Docker port mapping**: PostgreSQL container exposes port 5434 (not default 5432)
- **Game cleanup**: Runs hourly to remove stale games
- **WebSocket sync**: Message types defined in both frontend (`model/model.ts`) and backend (`internal/jeopardy/message.go`) - keep in sync
- **Tests in Docker**: Frontend tests require Chrome launcher, won't work in Docker
- **Concurrency safety**: Always use `SafeConn` for WebSocket writes; never write to `websocket.Conn` directly

## Common Development Tasks

### Adding a New Message Type
1. Add enum to `fe-jeopardy/src/app/model/model.ts:MessageType`
2. Add to `be-jeopardy/internal/jeopardy/message.go`
3. Handle in `fe-jeopardy/src/app/services/websocket.service.ts`
4. Handle in `be-jeopardy/internal/jeopardy/game.go:processMessage`

### Adding a New Game State
1. Add enum to `fe-jeopardy/src/app/model/model.ts:GameState`
2. Create component in `fe-jeopardy/src/app/game/`
3. Add to component switcher in `game.component.html`
4. Handle transitions in `be-jeopardy/internal/jeopardy/game.go`

### Modifying Question Selection
- SQL queries in `be-jeopardy/internal/db/sql/`
- Embedded via `//go:embed` directive in `db.go`
- **Important**: Requires rebuild after SQL changes

### Updating Answer Validation
- Algorithm in `be-jeopardy/internal/jeopardy/question.go:CheckAnswer`
- Alternative generation in `add_alternatives.py`
- Test with `go test ./internal/jeopardy/...`

### Debugging WebSocket Issues
- Check `SafeConn` usage - concurrent writes will panic without mutex
- Use `pm2 logs` to view real-time WebSocket messages
- Frontend DevTools → Network → WS to inspect WebSocket frames
- Backend logs show all incoming/outgoing messages when `GIN_MODE=debug`
