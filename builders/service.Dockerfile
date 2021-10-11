FROM golang:1.17 AS builder

RUN apt-get -qq update && apt-get -yqq install upx

ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG GO_TAGS=local
ARG SERVICE

RUN go build \
  -tags ${GO_TAGS} \
  -trimpath \
  -ldflags "-s -w -extldflags '-static'" \
  -installsuffix cgo \
  -o /bin/service \
  ./cmd/${SERVICE}

RUN strip /bin/service
RUN upx -q -9 /bin/service

RUN echo "nobody:x:65534:65534:Nobody:/:" > /etc_passwd

FROM scratch
COPY --from=builder /etc_passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /bin/service /bin/service

ENV TZ Asia/Bangkok
ENV PORT 8701

USER nobody

ENTRYPOINT ["/bin/service"]