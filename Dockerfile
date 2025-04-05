FROM golang:1.24-alpine@sha256:7772cb5322baa875edd74705556d08f0eeca7b9c4b5367754ce3f2f00041ccee AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
	-ldflags "-s -w" \
	-a -installsuffix cgo -o /go/bin/envtpl ./cmd/envtpl/.

	RUN ./test/test.sh

FROM scratch

COPY --from=builder /go/bin/envtpl /bin/envtpl

ENTRYPOINT [ "/bin/envtpl" ]
