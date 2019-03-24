FROM golang:1.10

RUN apt-get update
RUN apt-get -y install curl g++ make bzip2 nano

WORKDIR /go/src/go-binlog-replication
COPY . .

#installing dep and vendors
RUN go get -u github.com/golang/dep/...
CMD ["sh", "-c", "cd /go/src/go-binlog-replication/src && dep ensure -update && /bin/bash"]