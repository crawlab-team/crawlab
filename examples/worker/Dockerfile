FROM ubuntu:latest

ENV DEBIAN_FRONTEND=noninteractive

# 添加依赖描述文件
ADD requirements.txt /opt/crawlab/
# 添加二进制包
ADD crawlab /usr/local/bin/

RUN chmod +x /usr/local/bin/crawlab

# 安装基本环境
RUN apt-get update \
	&& apt-get install -y curl git net-tools iputils-ping ntp python3 python3-pip \
	&& apt-get clean \
	&& ln -s /usr/bin/pip3 /usr/local/bin/pip \
	&& ln -s /usr/bin/python3 /usr/local/bin/python

# 安装依赖
RUN pip install -r /opt/crawlab/requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple

WORKDIR /opt/crawlab

ENTRYPOINT ["crawlab"]