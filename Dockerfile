FROM golang:1.26-alpine@sha256:f85330846cde1e57ca9ec309382da3b8e6ae3ab943d2739500e08c86393a21b1 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
	-ldflags "-s -w" \
	-a -installsuffix cgo -o /go/bin/envtpl ./cmd/envtpl/.

	RUN ./test/test.sh

FROM scratch

COPY --from=builder /go/bin/envtpl /bin/envtpl

ENTRYPOINT [ "/bin/envtpl" ]
