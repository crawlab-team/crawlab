FROM crawlabteam/crawlab-backend:latest AS backend-build

FROM crawlabteam/crawlab-frontend:latest AS frontend-build

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

# copy nginx config files
COPY ./nginx/crawlab.conf /etc/nginx/conf.d

# install plugins
RUN /bin/bash /app/bin/docker-install-plugins.sh

# start backend
CMD ["/bin/bash", "/app/bin/docker-init.sh"]
