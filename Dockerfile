# docker build . -t rg.nl-ams.scw.cloud/teritori/teritorid:latest
# docker run --rm -it rg.nl-ams.scw.cloud/teritori/teritorid:latest /bin/sh
FROM golang:1.21-alpine3.17 AS go-builder
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
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.5.1/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.5.1/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep b89c242ffe2c867267621a6469f07ab70fc204091809d9c6f482c3fdf9293830
RUN sha256sum /lib/libwasmvm_muslc.x86_64.a | grep c0f4614d0835be78ac8f3d647a70ccd7ed9f48632bc1374db04e4df2245cb467

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
