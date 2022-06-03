FROM crawlabteam/crawlab-backend:latest AS backend-build

FROM crawlabteam/crawlab-frontend:latest AS frontend-build

FROM crawlabteam/crawlab-public-plugins:latest AS public-plugins-build

# images
FROM crawlabteam/crawlab-base:latest

# add files
COPY ./backend/conf /app/backend/conf
COPY ./nginx /app/nginx
COPY ./bin /app/bin

# copy backend files
RUN mkdir -p /opt/bin
COPY --from=backend-build /go/bin/crawlab /opt/bin
RUN cp /opt/bin/crawlab /usr/local/bin/crawlab-server

# copy frontend files
COPY --from=frontend-build /app/dist /app/dist

# copy public-plugins files
COPY --from=public-plugins-build /app/plugins /app/plugins

# copy nginx config files
COPY ./nginx/crawlab.conf /etc/nginx/conf.d

# start backend
CMD ["/bin/bash", "/app/bin/docker-init.sh"]
