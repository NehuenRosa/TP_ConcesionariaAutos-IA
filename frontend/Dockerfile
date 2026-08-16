# Etapa de compilación
FROM node:24-alpine AS compilacion

ARG VITE_API_URL=http://localhost:8080/api
ENV VITE_API_URL=$VITE_API_URL

WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm ci

COPY . .
RUN npm run build

# Etapa de ejecución (servidor estático no-root, filesystem de solo lectura)
FROM nginx:1.29-alpine

# Eliminar el server block por defecto (los scripts de entrypoint intentan
# modificarlo y fallan con el filesystem read-only).
RUN rm -f /etc/nginx/conf.d/default.conf

COPY --from=compilacion /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
