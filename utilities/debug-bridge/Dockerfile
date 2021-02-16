FROM busybox:uclibc AS busybox
FROM falcosecurity/kilt-utilities AS kilt-utils
FROM golang:alpine AS builder

ENV CGO_ENABLED=0
COPY . /debug-bridge/
WORKDIR /debug-bridge/debug-bridge
RUN go build -o /debug-ssh .

WORKDIR /debug-bridge/process-trace
RUN go build -o /process-trace .

FROM scratch
COPY --from=kilt-utils /kilt/waitforever /debug-bridge/waitforever
COPY --from=builder /debug-ssh /debug-bridge/debug-bridge
COPY --from=builder /process-trace /debug-bridge/process-trace
COPY --from=busybox /bin/ /debug-bridge/

VOLUME ["/debug-bridge"]

CMD ["/debug-bridge/waitforever"]
