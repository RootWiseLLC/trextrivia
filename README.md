# Jeopardy

Real-time multi-player Jeopardy game with over 100,000 questions.

Play here: https://playjeopardy.netlify.app

<details>
  <summary>Screenshots</summary>

![game](imgs/game.png)

<details>
  <summary>Home page</summary>

![home](imgs/home.png)

</details>

<details>
  <summary>Profile page</summary>

![profile](imgs/profile.png)

</details>

<details>
  <summary>Config page</summary>
  
![config](imgs/config.png)

</details>

<details>
  <summary>Analytics page</summary>
  
![analytics](imgs/analytics.png)

</details>

</details>

### Features

- Sign in with Google and GitHub accounts or email

- Player accounts with privacy controls and analytics (all-time record, totals and highs)

- Game configuration

  - Choose the categories you want to play with
  - Play against other people or against bots
  - Play solo or with up to 6 players
  - Play 1 or 2 round games
  - Play with or without penalties for incorrect answers
  - Play in public or private games

- In-game chat and emoji reactions

- Allows for players to pause the game

- Handles players disconnecting and rejoining the game

- Leaderboards

  - By % correct answers
  - By win %
  - By # of wins
  - By # of games
  - By total points (all-time)
  - By correct answers (all-time)
  - By total points (1 game)
  - By correct answers (1 game)

- Game analytics

  - Total number of games played
  - Average score after each round
  - Buzz-in rate for each round
  - Correctness rate for each round

### Development

#### First-time setup

**Backend setup:**

```bash
cd be-jeopardy

# Install Go dependencies
go mod tidy

# Generate JWT keys
mkdir -p .keys
./gen_keys.sh

# Create .env file
cat > .env << 'EOF'
PORT=8080
GIN_MODE=debug
DATABASE_URL=postgresql://postgres:postgres@localhost:5434/postgres?sslmode=disable
ALLOW_ORIGIN=http://localhost:4200
EOF

# Start PostgreSQL database
docker compose up -d postgres
```

**Frontend setup:**

```bash
cd fe-jeopardy

# Install dependencies (using PNPM)
pnpm install
```

#### Running the app

**Option 1: PM2 with hot reloading (recommended for development)**

First, install PM2 and air (Go hot reloader):
```bash
npm install -g pm2
go install github.com/air-verse/air@latest
```

Then start everything in dev mode:
```bash
# Start PostgreSQL
cd be-jeopardy && docker compose up -d postgres && cd ..

# Start backend and frontend with PM2 (with hot reloading)
pm2 start ecosystem.dev.config.js
```

Both backend (via air) and frontend (via Angular CLI) will automatically reload when you make code changes!

Manage your processes:
```bash
pm2 status              # View status of all processes
pm2 logs                # View all logs in real-time
pm2 logs jeopardy-backend   # View only backend logs
pm2 logs jeopardy-frontend  # View only frontend logs
pm2 restart all         # Restart all processes
pm2 stop all            # Stop all processes
pm2 delete all          # Stop and remove all processes
pm2 monit               # Real-time monitoring dashboard
```

Log files are stored in:
- Dev mode: `be-jeopardy/logs/backend-dev-*.log` and `fe-jeopardy/logs/frontend-dev-*.log`
- Production mode: `be-jeopardy/logs/backend-*.log` and `fe-jeopardy/logs/frontend-*.log`

**Production mode (no hot reloading):**

For running in production without hot reloading:
```bash
pm2 start ecosystem.config.js
```

Set to auto-start on boot (optional):
```bash
pm2 startup             # Follow the instructions
pm2 save                # Save current process list
```

**Option 2: Simple startup script**

```bash
./start.sh              # Start everything in background
./stop.sh               # Stop everything
tail -f backend.log     # View logs
```

**Option 3: Manual (separate terminals)**

```bash
# Terminal 1 - Backend
cd be-jeopardy
docker compose up -d postgres
source .env && make run

# Terminal 2 - Frontend
cd fe-jeopardy
pnpm run ng serve
```

Then open http://localhost:4200

#### Loading questions (optional)

To populate the database with Jeopardy questions:

```bash
cd be-jeopardy

# Install uv (Python package manager)
brew install uv

# Create Python virtual environment
uv venv

# Install Python dependencies
uv pip install requests beautifulsoup4 psycopg2-binary

# Scrape questions (this will take a while for full dataset)
cd scrapers
../.venv/bin/python jeopardy.py

# Load questions into database
cd ..
mv scrapers/season40.tsv clues/
.venv/bin/python insert_clues.py
.venv/bin/python add_alternatives.py
```

