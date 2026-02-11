FROM golang:1.26-alpine@sha256:d4c4845f5d60c6a974c6000ce58ae079328d03ab7f721a0734277e69905473e5 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
	-ldflags "-s -w" \
	-a -installsuffix cgo -o /go/bin/envtpl ./cmd/envtpl/.

	RUN ./test/test.sh

FROM scratch

COPY --from=builder /go/bin/envtpl /bin/envtpl

ENTRYPOINT [ "/bin/envtpl" ]
