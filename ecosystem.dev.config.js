module.exports = {
  apps: [
    {
      name: 'jeopardy-backend-dev',
      cwd: './be-jeopardy',
      script: process.env.HOME + '/go/bin/air',
      args: '',
      env: {
        PORT: '8080',
        GIN_MODE: 'debug',
        DATABASE_URL: 'postgresql://postgres:postgres@localhost:5434/postgres?sslmode=disable',
        ALLOW_ORIGIN: 'http://localhost:4200'
      },
      autorestart: true,
      watch: false,  // air handles its own watching
      max_memory_restart: '500M',
      error_file: '../logs/backend-dev-error.log',
      out_file: '../logs/backend-dev-out.log',
      log_date_format: 'YYYY-MM-DD HH:mm:ss Z'
    },
    {
      name: 'jeopardy-frontend-dev',
      cwd: './fe-jeopardy',
      script: 'pnpm',
      args: 'run ng serve',
      autorestart: true,
      watch: false,  // Angular CLI handles its own watching
      max_memory_restart: '1G',
      error_file: '../logs/frontend-dev-error.log',
      out_file: '../logs/frontend-dev-out.log',
      log_date_format: 'YYYY-MM-DD HH:mm:ss Z'
    }
  ]
};
