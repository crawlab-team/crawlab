FROM node:8.16.0-alpine AS frontend-build

ADD . /app
WORKDIR /app

# install frontend
RUN npm install -g yarn \
	&& yarn install --registry=https://registry.npm.taobao.org

RUN npm run build:prod

FROM alpine

#RUN apk update
RUN apk add nginx
COPY --from=frontend-build /app/dist /app/dist
COPY --from=frontend-build /app/conf/crawlab.conf /etc/nginx/conf.d
#RUN nginx -s start
#COPY ./dist /usr/share/nginx/html

#EXPOSE 80
#EXPOSE 8080

