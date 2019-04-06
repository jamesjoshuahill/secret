# ciphers

[![Go Report Card](https://goreportcard.com/badge/github.com/jamesjoshuahill/ciphers)](https://goreportcard.com/report/github.com/jamesjoshuahill/ciphers)

A microservice written in Go that creates ciphers using AES encryption.

Ciphers are stored in memory and can be retrieved with the correct key.

## Get

_This module requires Go 1.11 or greater._

Download into GOPATH:

```bash
go get github.com/jamesjoshuahill/ciphers
export GO111MODULE=on
```

or elsewhere:

```bash
git clone https://github.com/jamesjoshuahill/ciphers.git
```

## Test

Run the tests using the [Ginkgo](https://onsi.github.io/ginkgo/) test runner:

```bash
go get github.com/onsi/ginkgo/ginkgo
ginkgo -r
```

or:

```bash
go test ./...
```

## Run

_The server requires TLS configuration_

For example, use the self-signed certificate and private key used by the test suite:

```bash
go run cmd/server/main.go \
  --port 8080 \
  --cert acceptance/fixtures/cert.pem \
  --key acceptance/fixtures/key.pem
```

Then, create a cipher:

```bash
curl \
  --cacert acceptance/fixtures/cert.pem \
  https://127.0.0.1:8080/v1/ciphers \
  -X POST \
  -H 'Content-Type: application/json' \
  -d '{"id":"some-id","data":"some plain text"}'
```

and retrieve it:

```bash
curl \
  --cacert acceptance/fixtures/cert.pem \
  https://127.0.0.1:8080/v1/ciphers/some-id \
  -X GET \
  -H 'Content-Type: application/json' \
  -d '{"key":"AES KEY for cipher"}'
```

## Client

The `client` package provides a `Client` to interact with the server.

Please refer to the package [documentation](https://godoc.org/github.com/jamesjoshuahill/ciphers/pkg/client).

## API

Please refer to the [API specification](API.md).
