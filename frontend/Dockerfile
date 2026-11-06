# Build stage
FROM node:18-alpine AS builder
WORKDIR /app

# Copy sources and install dependencies. `npm ci` runs after the copy so it always
# replaces any node_modules folder coming from the host (e.g. Windows binaries).
COPY . .
RUN npm ci

# Same-origin API: nginx proxies /api to the backend (see nginx.conf)
ENV VITE_API_URL=/api

# Build the app
RUN npm run build

# Production stage
FROM nginx:stable-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
