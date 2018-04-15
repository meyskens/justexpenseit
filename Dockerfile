ARG ARCH
FROM golang AS gobuild

ARG GOARCH
ARG GOARM

COPY ./ /go/src/github.com/meyskens/justexpenseit
WORKDIR /go/src/github.com/meyskens/justexpenseit

RUN GOARCH=${GOARCH} GOARM=${GOARM} go build ./

ARG ARCH
FROM multiarch/alpine:${ARCH}-edge

RUN apk add --no-cache ca-certificates

COPY --from=gobuild /go/src/github.com/meyskens/justexpenseit/justexpenseit /usr/local/bin/justexpenseit

ENTRYPOINT [ "justexpenseit" ]
