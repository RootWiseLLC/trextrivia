FROM nginx:alpine

# Copy built app
COPY dist/jeopardy /usr/share/nginx/html

# Copy nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
