FROM node:12 AS build

ADD . /app
WORKDIR /app
RUN rm /app/.npmrc

# install frontend
RUN yarn install && yarn run build:docker

FROM alpine:3.14

# copy files
COPY --from=build /app/dist /app/dist
