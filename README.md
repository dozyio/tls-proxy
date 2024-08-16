# TLS Proxy

A simple HTTP2 TLS proxy for connecting to local or remote services. Ideal for testing and development. Your mileage may vary when connecting to remote domains that require certificates, but it should be fine for proxying to local services.

## Installing

```sh
go install github.com/dozyio/tls-proxy@latest
```

## Running

Listen on 0.0.0.0:9000 and proxy to http://127.0.0.1:3000, using the supplied certificate and key.

```sh
tls-proxy -l "0.0.0.0:9000" -t "http://127.0.0.1:3000" -c "192.168.1.2.crt" -k "192.168.1.2.key"
```

## Generating Certificates

```sh
openssl genrsa -out 192.168.1.2.key 2048
openssl req -new -key 192.168.1.2.key -out 192.168.1.2.csr
openssl x509 -req -days 3650 -in 192.168.1.2.csr -signkey 192.168.1.2.key -out 192.168.1.2.crt
```
