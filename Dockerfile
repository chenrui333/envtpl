FROM golang:1.25-alpine@sha256:f6751d823c26342f9506c03797d2527668d095b0a15f1862cddb4d927a7a4ced AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
	-ldflags "-s -w" \
	-a -installsuffix cgo -o /go/bin/envtpl ./cmd/envtpl/.

	RUN ./test/test.sh

FROM scratch

COPY --from=builder /go/bin/envtpl /bin/envtpl

ENTRYPOINT [ "/bin/envtpl" ]
