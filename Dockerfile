# images
FROM ubuntu:latest

# set as non-interactive
ENV DEBIAN_FRONTEND noninteractive

# environment variables
ENV NVM_DIR /usr/local/nvm  
ENV NODE_VERSION 8.12
ENV WORK_DIR /opt/crawlab

# source files
ADD . /opt/crawlab

# install python
RUN apt-get update
RUN apt-get install -y python3 python3-pip net-tools iputils-ping redis-server git nginx ntp curl

# python soft link
RUN ln -s /usr/bin/pip3 /usr/local/bin/pip
RUN ln -s /usr/bin/python3 /usr/local/bin/python

# install mongodb
RUN echo "Asia/Shanghai" > /etc/timezone && dpkg-reconfigure -f noninteractive tzdata
RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 9DA31620334BD75D9DCB49F368818C72E52529D4
RUN echo "deb [ arch=amd64 ] https://repo.mongodb.org/apt/ubuntu bionic/mongodb-org/4.0 multiverse" | tee /etc/apt/sources.list.d/mongodb-org-4.0.list
RUN apt-get update
RUN apt-get install -y mongodb-org

# install backend
RUN pip install -U setuptools -i https://pypi.tuna.tsinghua.edu.cn/simple
RUN pip install -r /opt/crawlab/crawlab/requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple

# install nvm
RUN curl https://raw.githubusercontent.com/creationix/nvm/v0.24.0/install.sh | bash \  
    && bash $NVM_DIR/nvm.sh \
    && nvm install v$NODE_VERSION \
    && nvm use v$NODE_VERSION \
    && nvm alias default v$NODE_VERSION

# install frontend
RUN npm install -g yarn pm2
RUN cd /opt/crawlab/frontend && yarn install

# nginx config & start frontend
RUN cp $WORK_DIR/crawlab.conf /etc/nginx/conf.d
RUN service nginx reload

# start mongodb
CMD mongod

# start redis
CMD redis-server

# start backend
WORKDIR /opt/crawlab/crawlab
CMD python $WORK_DIR/crawlab/app.py 
CMD python $WORK_DIR/crawlab/flower.py 
CMD python $WORK_DIR/crawlab/worker.py 
#CMD pm2 start $WORK_DIR/crawlab/app.py 
#CMD pm2 start $WORK_DIR/crawlab/flower.py 
#CMD pm2 start $WORK_DIR/crawlab/worker.py 

EXPOSE 8080
EXPOSE 8000