   # Dockerfile
   FROM golang:1.21-alpine

   WORKDIR /app

   # Install dependencies
   RUN apk add --no-cache git curl postgresql-client

   # Install sqlc
   RUN curl -L https://github.com/kyleconroy/sqlc/releases/download/v1.10.0/sqlc_1.10.0_linux_amd64.tar.gz | tar -xz -C /usr/local/bin

   # Install dbmate
   RUN curl -L https://github.com/amacneil/dbmate/releases/download/v1.11.0/dbmate-linux-amd64 -o /usr/local/bin/dbmate && \
       chmod +x /usr/local/bin/dbmate

   COPY go.mod ./
   COPY go.sum ./
   RUN go mod download

   COPY . ./

   RUN chmod +x wait-for-it.sh
   # Copy .env.docker to .env
   RUN cp .env.docker .env

   RUN go build -o /billing-engine main.go

   EXPOSE 8080

   CMD ["/billing-engine"]
