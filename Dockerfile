FROM golang:alpine as build

RUN     apk update \
    &&  apk upgrade \
    &&  apk add --no-cache \
            git \
            make \
            bash

WORKDIR $GOPATH/src/github.com/xkortex/justget

COPY . ./

RUN go get

RUN make

FROM build as inline_test

RUN ./tests/get.sh

FROM scratch

COPY --from=build /go/bin/justget /go/bin/justget

ENTRYPOINT ["/go/bin/justget"]
