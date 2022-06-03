FROM node:12

ADD ./frontend /app
WORKDIR /app
RUN rm /app/.npmrc

# install frontend
RUN yarn install && yarn run build:docker
