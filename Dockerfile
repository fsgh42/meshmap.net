FROM golang:latest AS generate-protobuf
COPY . /app
WORKDIR /app
RUN \
    apt-get update && \
    apt-get install -y protobuf-compiler patch && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN \
    cd /tmp && \
    git clone --progress --depth 1 https://github.com/meshtastic/protobufs.git && \
    patch -p1 < /app/patch/remove-nanopb.patch && \
    rm -rf /app/internal/meshtastic/generated && \
    protoc -I=protobufs --go_out=/app/internal/meshtastic --go_opt=module=github.com/meshtastic/go protobufs/meshtastic/*.proto

FROM golang:latest AS meshobserv
COPY --from=generate-protobuf /app /build
WORKDIR /build
COPY cmd/meshobserv ./
COPY internal ./internal/
RUN go get ./... && go mod tidy
RUN go build -v -o /usr/bin/meshobserv && chmod 755 /usr/bin/meshobserv

FROM debian:bookworm-slim
COPY --from=build --chmod=666 /build/configs/blocklist.txt /blocklist.txt
COPY --from=build /usr/bin/meshobserv /usr/bin/meshobserv
ENTRYPOINT ["/usr/bin/meshobserv"]
