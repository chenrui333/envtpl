FROM golang:1.20-alpine
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
	-ldflags "-s -w" \
	-a -installsuffix cgo -o /go/bin/envtpl ./cmd/envtpl/.
RUN ./test/test.sh

FROM scratch
COPY --from=0 /go/bin/envtpl /bin/envtpl
ENTRYPOINT [ "/bin/envtpl" ]
