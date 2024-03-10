FROM golang:1.19 as builder
WORKDIR /app
COPY go.mod go.mod
RUN go mod download
COPY app /app
RUN CGO_ENABLED=1 GOOS=linux GOFLAGS=-mod=mod go build -a -installsuffix cgo -o main .

FROM debian:latest
WORKDIR /docker
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates
COPY --from=builder /app/main /docker/main
# Reserved Exporter Port for this exporter
EXPOSE 8080
run mkdir -p /docker/server/ /docker/ccd /docker/data && \
	chown nobody:nogroup -R /docker/server /docker/ccd /docker/data
USER nobody:nogroup
VOLUME /docker/server
VOLUME /docker/ccd
CMD ["/docker/main"]
