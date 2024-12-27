# docker build . -t rg.nl-ams.scw.cloud/teritori/teritorid:latest
# docker run --rm -it rg.nl-ams.scw.cloud/teritori/teritorid:latest /bin/sh
FROM golang:1.22-alpine3.18 AS go-builder
ARG arch=x86_64

# this comes from standard alpine nightly file
#  https://github.com/rust-lang/docker-rust-nightly/blob/master/alpine3.12/Dockerfile
# with some changes to support our toolchain, etc
RUN set -eux; apk add --no-cache ca-certificates build-base;

RUN apk add git
# NOTE: add these to run with LEDGER_ENABLED=true
# RUN apk add libusb-dev linux-headers

WORKDIR /code
COPY . /code/
# See https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.3.1/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.3.1/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 9f5f387bddbd51a549f2d8eb2cc0f08030f311610c43004b96c141a733b681bc
RUN sha256sum /lib/libwasmvm_muslc.x86_64.a | grep 129da0e50eaa3da093eb84c3a8f48da20b9573f3afdf83159b3984208c8ec5c3

# Copy the library you want to the final location that will be found by the linker flag `-lwasmvm_muslc`
RUN cp /lib/libwasmvm_muslc.${arch}.a /lib/libwasmvm_muslc.a

# force it to use static lib (from above) not standard libgo_cosmwasm.so file
RUN LEDGER_ENABLED=false BUILD_TAGS=muslc LINK_STATICALLY=true make build
RUN echo "Ensuring binary is statically linked ..." \
    && (file /code/build/teritorid | grep "statically linked")

# --------------------------------------------------------
FROM alpine:3.17

COPY --from=go-builder /code/build/teritorid /usr/bin/teritorid

#COPY docker/* /opt/
#RUN chmod +x /opt/*.sh

WORKDIR /opt

# rest server
EXPOSE 1317
# tendermint p2p
EXPOSE 26656
# tendermint rpc
EXPOSE 26657

CMD ["/usr/bin/teritorid", "version"]
