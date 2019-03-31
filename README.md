# ciphers

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
git clone git@github.com:jamesjoshuahill/ciphers.git
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

## Create a cipher

**POST /v1/ciphers**

Example request:
```bash
curl \
  --cacert acceptance/fixtures/cert.pem \
  https://127.0.0.1:8080/v1/ciphers \
  -X POST \
  -d '{"id":"some-id","data":"some plain text"}'
```

Example response:
```json
{
  "key":"1bc50ee2992feba6c1d9e384b3c8e9203dcfc0eed50c032dfc2821ca2aa0cfa5",
}
```

_The client is responsible for storing the id and key._

## Get a cipher

**GET /v1/ciphers/{id}**

Example request:
```bash
curl \
  --cacert acceptance/fixtures/cert.pem \
  https://127.0.0.1:8080/v1/ciphers/some-id \
  -X GET \
  -d '{"key":"1bc50ee2992feba6c1d9e384b3c8e9203dcfc0eed50c032dfc2821ca2aa0cfa5"}'
```

Example response:
```json
{
  "data": "some plain text"
}
```
