FROM golang:latest

WORKDIR /app/
COPY . /app/

RUN go mod download
RUN go mod tidy

RUN GOOS=linux GOARCH=amd64 go build \
    --ldflags "-s -w" \
    -o spirit \
    -tags sqlite \
    ./cmd/spirit/main.go

# Stage 1 done
RUN apt-get update -y \
 && apt-get install -y curl ca-certificates openssl git tar fontconfig tzdata iproute2 locales \
 && useradd -d /home/container -m container

USER container
ENV  USER=container HOME=/home/container

WORKDIR     /home/container

COPY        entrypoint.sh /entrypoint.sh

CMD         ["/bin/bash", "/entrypoint.sh"]

