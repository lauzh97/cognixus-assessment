FROM debian:latest

WORKDIR /usr/app

COPY config.yaml ./

RUN apt-get update
RUN apt-get -y install wget

RUN wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
ENV PATH="${PATH}:/usr/local/go/bin"

RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# RUN apt-get -y install postgresql postgresql-contrib

COPY . /usr/app

EXPOSE 8081
EXPOSE 8090

WORKDIR /usr/app/cmd/server
RUN go build

WORKDIR /usr/app
CMD [ "./cmd/server/server" ]