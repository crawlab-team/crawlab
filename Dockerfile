# images
#FROM python:latest
FROM ubuntu:latest

# source files
ADD . /opt/crawlab

# add dns
#RUN echo -e "nameserver 180.76.76.76" >> /etc/resolv.conf
#ADD ./resolv.conf /etc
RUN cat /etc/resolv.conf

# install python
RUN apt-get update
RUN apt-get install -y python3 python3-pip net-tools iputils-ping

# soft link
RUN ln -s /usr/bin/pip3 /usr/local/bin/pip
RUN ln -s /usr/bin/python3 /usr/local/bin/python

# install required libraries
RUN pip install -U setuptools
RUN pip install -r /opt/crawlab/requirements.txt

# execute apps
WORKDIR /opt/crawlab
CMD python ./bin/run_worker.py
CMD python app.py
