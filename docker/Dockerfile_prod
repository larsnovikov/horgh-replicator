FROM golang:1.10

RUN apt-get update
RUN apt-get -y install curl g++ make bzip2 nano unixodbc unixodbc-dev mysql-client

WORKDIR /go/src/horgh-replicator
COPY . .

COPY files/vertica-client-7.2.0-0.x86_64.tar.gz /vertica-client.tar.gz
RUN tar -xvf /vertica-client.tar.gz -C /

#installing dep and vendors
RUN go get -u github.com/golang/dep/...

RUN ["sh", "-c", "cd /go/src/horgh-replicator/src \
    && dep ensure -update \
    && go build main.go \
    && mv main horgh-replicator"]

WORKDIR /go/src/horgh-replicator/src

ENTRYPOINT ["./horgh-replicator", "listen"]