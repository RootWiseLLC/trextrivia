#!/bin/bash

# Stop Jeopardy backend and frontend

echo "Stopping Jeopardy..."

cd /Users/david/build/jeopardy

if [ -f .backend.pid ]; then
    BACKEND_PID=$(cat .backend.pid)
    if ps -p $BACKEND_PID > /dev/null 2>&1; then
        kill $BACKEND_PID
        echo "✓ Backend stopped (PID: $BACKEND_PID)"
    else
        echo "✓ Backend was not running"
    fi
    rm .backend.pid
fi

if [ -f .frontend.pid ]; then
    FRONTEND_PID=$(cat .frontend.pid)
    if ps -p $FRONTEND_PID > /dev/null 2>&1; then
        kill $FRONTEND_PID
        echo "✓ Frontend stopped (PID: $FRONTEND_PID)"
    else
        echo "✓ Frontend was not running"
    fi
    rm .frontend.pid
fi

echo ""
echo "Jeopardy stopped."
