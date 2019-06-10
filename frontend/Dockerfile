# images
FROM node:8.12

# source files
ADD . /opt/crawlab/frontend

# environment variables
#ENV NVM_DIR /usr/local/nvm
#ENV NODE_VERSION 8.12.0
#ENV WORK_DIR /opt/crawlab/frontend

# install git curl
RUN apt-get update
RUN apt-get install -y nginx
#RUN apt-get install -y git curl

# install nvm
#RUN curl https://raw.githubusercontent.com/creationix/nvm/v0.24.0/install.sh | bash \
#    && . $NVM_DIR/nvm.sh \
#    && nvm install v$NODE_VERSION \
#    && nvm use v$NODE_VERSION \
#    && nvm alias default v$NODE_VERSION
#ENV NODE_PATH $NVM_DIR/versions/node/v$NODE_VERSION/lib/node_modules
#ENV PATH $NVM_DIR/versions/node/v$NODE_VERSION/bin:$PATH

# install frontend
RUN npm install -g yarn pm2 --registry=https://registry.npm.taobao.org
RUN cd /opt/crawlab/frontend && yarn install --registry=https://registry.npm.taobao.org

# nginx config & start frontend
RUN cp $WORK_DIR/conf/crawlab.conf /etc/nginx/conf.d
RUN service nginx reload

CMD ["npm", "run", "build:prod"]
