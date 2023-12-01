FROM golang:1.12 as builder
RUN mkdir /build
COPY . /build/
WORKDIR /build
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -o dns_tls_proxy .

FROM golang:1.21.4-alpine3.18
COPY --from=builder /build/dns_tls_proxy /app/
WORKDIR /app

EXPOSE 53

CMD [ "./dns_tls_proxy" ]