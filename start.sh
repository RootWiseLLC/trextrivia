#!/bin/bash

# Start Jeopardy backend and frontend
set -e

echo "Starting Jeopardy..."
echo ""

# Check if PostgreSQL is running
echo "Checking PostgreSQL..."
cd be-jeopardy
docker compose up -d postgres
echo "✓ PostgreSQL is running"
echo ""

# Start backend in background
echo "Starting backend server..."
cd /Users/david/build/jeopardy/be-jeopardy
source .env && make run > ../backend.log 2>&1 &
BACKEND_PID=$!
echo "✓ Backend started (PID: $BACKEND_PID)"
echo "  Logs: tail -f backend.log"
echo ""

# Wait a moment for backend to start
sleep 2

# Start frontend in background
echo "Starting frontend server..."
cd /Users/david/build/jeopardy/fe-jeopardy
pnpm run ng serve > ../frontend.log 2>&1 &
FRONTEND_PID=$!
echo "✓ Frontend started (PID: $FRONTEND_PID)"
echo "  Logs: tail -f frontend.log"
echo ""

# Save PIDs to file for stopping later
cd /Users/david/build/jeopardy
echo $BACKEND_PID > .backend.pid
echo $FRONTEND_PID > .frontend.pid

echo "================================"
echo "Jeopardy is running!"
echo "================================"
echo ""
echo "Frontend: http://localhost:4200"
echo "Backend:  http://localhost:8080"
echo ""
echo "To view logs:"
echo "  Backend:  tail -f backend.log"
echo "  Frontend: tail -f frontend.log"
echo ""
echo "To stop:"
echo "  ./stop.sh"
echo ""
