FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./
RUN go build -o /app/tasks ./services/tasks/cmd/tasks

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/tasks /app/tasks
EXPOSE 8082
ENTRYPOINT ["/app/tasks"]
