# secret

[![Build Status](https://travis-ci.org/jamesjoshuahill/secret.svg?branch=master)](https://travis-ci.org/jamesjoshuahill/secret)
[![Go Report Card](https://goreportcard.com/badge/github.com/jamesjoshuahill/secret)](https://goreportcard.com/report/github.com/jamesjoshuahill/secret)

A microservice written in Go that stores secrets encrypted with AES.

Ciphertext is stored in memory and can be retrieved with the correct AES key.

## Get

_This module requires Go 1.15 or greater._

```bash
git clone https://github.com/jamesjoshuahill/secret.git
cd secret
```

## Lint

Install [golangci-lint](https://golangci-lint.run/usage/install/#local-installation).

Run linters:

```bash
golangci-lint run
```

## Test

Run the tests using the [Ginkgo](https://onsi.github.io/ginkgo/) test runner:

```bash
go install github.com/onsi/ginkgo/ginkgo
ginkgo -race -cover -r
```

or:

```bash
go test -race -cover ./...
```

## Run

_The server requires TLS configuration._

For example, use the self-signed certificate and private key used by the test suite:

```bash
go run cmd/secret-server/main.go \
  --host localhost \
  --port 8080 \
  --cert acceptance_test/testdata/cert.pem \
  --key acceptance_test/testdata/key.pem
```

Then, create a secret:

```bash
curl \
  --cacert acceptance_test/testdata/cert.pem \
  https://localhost:8080/v1/secrets \
  -X POST \
  -H 'Content-Type: application/json' \
  -d '{"id":"some-id","data":"some plain text"}'
```

and retrieve it using the AES key:

```bash
curl \
  --cacert acceptance_test/testdata/cert.pem \
  https://localhost:8080/v1/secrets/some-id \
  -X GET \
  -H 'Content-Type: application/json' \
  -d '{"key":"AES KEY for secret"}'
```

## Client

The `client` package provides a `Client` to interact with the server.

Please refer to the [package documentation](https://pkg.go.dev/github.com/jamesjoshuahill/secret/pkg/client).

## API

Please refer to the [API specification](API.md).
