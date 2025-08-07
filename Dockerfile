FROM golang:1.24-alpine@sha256:c8c5f95d64aa79b6547f3b626eb84b16a7ce18a139e3e9ca19a8c078b85ba80d AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
	-ldflags "-s -w" \
	-a -installsuffix cgo -o /go/bin/envtpl ./cmd/envtpl/.

	RUN ./test/test.sh

FROM scratch

COPY --from=builder /go/bin/envtpl /bin/envtpl

ENTRYPOINT [ "/bin/envtpl" ]
