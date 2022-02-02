FROM alpine:3.15.0 AS builder
RUN apk add --no-cache tzdata

FROM scratch
COPY ./builders/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

USER nobody
ARG SERVICE
COPY ./bin/${SERVICE} /bin/service

ENV TZ Asia/Bangkok
ENV PORT 8701
ENTRYPOINT ["/bin/service"]