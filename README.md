# go-grpc-example
Ping gRPC server/client in GO with TLS encryption, proxy api and auth.

## How generate certs

```
openssl genrsa -out cert/server.key 2048
openssl req -new -x509 -sha256 -key cert/server.key -out cert/server.crt -days 3650
openssl req -new -sha256 -key cert/server.key -out cert/server.csr
openssl x509 -req -sha256 -in cert/server.csr -signkey cert/server.key -out cert/server.crt -days 3650
```

## See makefile for build instructions.

## Reference:
https://medium.com/pantomath/how-we-use-grpc-to-build-a-client-server-system-in-go-dd20045fa1c2
