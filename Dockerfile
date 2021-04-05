FROM golang:1.16.3 as builder
WORKDIR /app
COPY go.mod go.mod
RUN go mod download
COPY app /app
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

FROM debian:latest
WORKDIR /docker
COPY --from=builder /app/main /docker/main
# Reserved Exporter Port for this exporter
EXPOSE 8080
USER nobody:nogroup
CMD ["/docker/main"]
