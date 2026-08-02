# Etapa de compilación
FROM node:24-alpine AS compilacion

WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm install

COPY . .
RUN npm run build

# Etapa de ejecución (servidor estático)
FROM nginx:1.29-alpine

COPY --from=compilacion /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
