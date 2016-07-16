# sudo docker build -t findserver . && sudo docker run -it findserver
# Run headless using: sudo docker run -itd findserver
FROM ubuntu:16.04

EXPOSE 8003 1883

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
RUN go build
RUN mkdir mosquitto
RUN touch mosquitto/conf
ENTRYPOINT mosquitto -c /root/find/mosquitto/conf -d && ./find -mqtt localhost:1883 -mqttadmin admin -mqttadminpass 123 -mosquitto `pgrep mosquitto` > log & bash
