FROM node:14 AS build

ADD . /app
WORKDIR /app
RUN rm /app/.npmrc

# install frontend
RUN npm i -g pnpm
RUN pnpm install
RUN pnpm run build

FROM alpine:3.14

# copy files
COPY --from=build /app/dist /app/dist
