FROM golang:latest
ARG branch=master
ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH $GOROOT/bin:$PATH
#作者
MAINTAINER Razil "jiawang06@163.com"

RUN apt-get -qq update && apt-get -qq install curl unzip

ADD https://releases.hashicorp.com/consul/1.2.1/consul_1.2.1_linux_amd64.zip /tmp/consul.zip
RUN cd /usr/sbin && unzip /tmp/consul.zip && chmod +x /usr/sbin/consul && rm /tmp/consul.zip
#CMD consul agent -dev > consul.log && tail -F consul.log && sleep 3
#CMD [ "/usr/sbin/consul", "agent", "-dev","-D" ]
  
ADD https://github.com/bottos-project/magiccube/raw/master/vendor/micro /usr/sbin/micro
RUN cd /usr/sbin && chmod +x /usr/sbin/micro

 
#设置工作目录
WORKDIR $GOPATH/src/github.com/bottos-project

#将服务器的go工程代码加入到docker容器中
RUN git clone -b $branch https://github.com/bottos-project/bottos.git --recursive
RUN git clone https://github.com/bottos-project/crypto-go.git

#go构建可执行文件
WORKDIR $GOPATH/src/github.com/bottos-project/bottos
RUN go build .
 
RUN chmod +x ./docker/botNode.sh

#暴露端口
EXPOSE 53/udp 8300 8301 8301/udp 8302 8302/udp 8400 8500 8080
VOLUME /go/src/github.com/bottos-project/bottos/datadir
#最终运行docker的命令
WORKDIR $GOPATH/src/github.com/bottos-project/bottos/docker/
ENTRYPOINT  ["./botNode.sh"]
CMD ["start"]
