# Dockerfile para el Frontend de React/Next.js (CertiTrack)

# Usa la imagen oficial de Node.js como base para la etapa de construcción
# Actualizado a Node.js 24
FROM node:24-alpine AS builder

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos de configuración de paquetes
COPY package.json pnpm-lock.yaml ./

# Instala las dependencias usando pnpm
RUN npm install -g pnpm && pnpm install --frozen-lockfile

# Copia el resto del código fuente de la aplicación
COPY . .

# Construye la aplicación frontend (ej. para Next.js, esto generará los archivos estáticos)
# Asegúrate de que tu script de build en package.json sea 'pnpm build'
RUN pnpm build

# Usa una imagen base ligera para la etapa final (para servir los archivos estáticos)
# Actualizado a Node.js 24
FROM node:24-alpine AS runner

# Establece el directorio de trabajo
WORKDIR /app

# Instala pnpm en la etapa runner también
RUN npm install -g pnpm

# Copia los archivos de construcción desde la etapa 'builder'
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json

# Expone el puerto en el que el frontend escuchará (ej. el puerto por defecto de Next.js)
EXPOSE 3000

# Comando para ejecutar la aplicación cuando el contenedor se inicie
# Para Next.js, esto iniciará el servidor de producción
CMD ["pnpm", "start"]