FROM golang:1.18 as builder

ENV GO111MODULE=on GOPROXY=https://goproxy.cn,direct CGO_ENABLED=0

WORKDIR /go/src/remote

COPY . /go/src/remote/

RUN make fmt && make build

FROM scratch 

LABEL maintainer="<yangj280@chinaunicom.cn>"

COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /go/src/remote/remote /opt/remote-agent/bin/

ENTRYPOINT ["/opt/remote-agent/bin/remote"]
