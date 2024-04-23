FROM golang:1.22.2-alpine AS builder

# WORKDIR /app

# COPY go.mod go.sum ./
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
 


# RUN go mod tidy
# RUN go mod vendor

# COPY . .

# RUN go build -o main main.go

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate

COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration




EXPOSE 8000

CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]