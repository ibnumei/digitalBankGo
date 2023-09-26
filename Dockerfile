#build stage
FROM golang:1.21-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
#install curl first, because alpine image doesnt have  curl
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz

#Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]

#docker build -t digitalbank:latest .
#docker run --name digitalbank -p 8080:8080 -e GIN_MODE=release digitalbank:latest
#docker container inspect postgres12 (IPAddress for postgres 172.17.0.2)
#docker run --name digitalbank -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:mysecretpassword@172.17.0.2:5432/master_bank?sslmode=disable" digitalbank:latest