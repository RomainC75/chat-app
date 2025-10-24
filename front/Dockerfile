# FROM node:22-alpine AS build-stage
FROM node:22.14-alpine3.20 AS build-stage

WORKDIR /app

RUN corepack enable

COPY package.json pnpm-lock.yaml ./

RUN pnpm install --frozen-lockfile

COPY . .

ENV NODE_ENV=development
ENV PUBLIC_PATH=/

CMD ["pnpm", "run", "dev", "--host"]

