# sudo docker build -t finddocker .
# sudo docker run -it -p 18003:8003 -p 11883:1883 -v /path/to/host/data/folder:/data finddocker bash
FROM ubuntu:16.04

# Get basics
RUN apt-get update
RUN apt-get -y upgrade
RUN apt-get install -y golang git wget curl vim
RUN mkdir /usr/local/work
ENV GOPATH /usr/local/work

# Install SVM
WORKDIR "/tmp"
RUN wget http://www.csie.ntu.edu.tw/~cjlin/cgi-bin/libsvm.cgi?+http://www.csie.ntu.edu.tw/~cjlin/libsvm+tar.gz -O libsvm.tar.gz
RUN tar -xvzf libsvm.tar.gz
RUN mv libsvm-*/* ./
RUN make
RUN cp svm-scale /usr/local/bin/
RUN cp svm-predict /usr/local/bin/
RUN cp svm-train /usr/local/bin/
RUN rm -rf *

# Install mosquitto
RUN apt-get install -y mosquitto-clients mosquitto

# Install FIND
WORKDIR "/root"
RUN go get github.com/schollz/find
RUN git clone https://github.com/schollz/find.git
WORKDIR "/root/find"
RUN mkdir mosquitto
RUN touch mosquitto/conf
RUN git pull
RUN go build
#ENTRYPOINT git pull && go build && mosquitto -c /root/find/mosquitto/conf -d && ./find -mqtt localhost:1883 -mqttadmin admin -mqttadminpass 123 -mosquitto `pgrep mosquitto` -data /data > log & bash

# Setup supervisor
RUN apt-get update && apt-get install -y supervisor
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# Add Tini
ENV TINI_VERSION v0.9.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini
ENTRYPOINT ["/tini", "--"]

# Startup
CMD ["/usr/bin/supervisord"]
