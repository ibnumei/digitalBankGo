#build stage
FROM golang:1.21-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

#Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .

EXPOSE 8080
CMD [ "/app/main" ]

#docker build -t digitalbank:latest .
#docker run --name digitalbank -p 8080:8080 -e GIN_MODE=release digitalbank:latest
#docker container inspect postgres12 (IPAddress for postgres 172.17.0.2)
#docker run --name digitalbank -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:mysecretpassword@172.17.0.2:5432/master_bank?sslmode=disable" digitalbank:latest