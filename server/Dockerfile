FROM golang:1.21-bullseye AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/godocker .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/godocker /app/godocker

EXPOSE 3000

CMD ["/app/godocker"]
