FROM golang:1.16 as builder
WORKDIR /app
COPY go.mod go.mod
RUN go mod download
COPY app /app
RUN CGO_ENABLED=1 GOOS=linux GOFLAGS=-mod=mod go build -a -installsuffix cgo -o main .

FROM debian:latest
WORKDIR /docker
COPY --from=builder /app/main /docker/main
# Reserved Exporter Port for this exporter
EXPOSE 8080
run mkdir -p /docker/server/ /docker/ccd /docker/data && \
	chown nobody:nogroup -R /docker/server /docker/ccd /docker/data
USER nobody:nogroup
VOLUME /docker/server
VOLUME /docker/ccd
COPY public /docker/public
CMD ["/docker/main"]
