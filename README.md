# go-grpc-example
Ping gRPC server/client in GO with TLS encryption, proxy api and auth.

## How to generate certs

```
openssl genrsa -out cert/server.key 2048
openssl req -new -x509 -sha256 -key cert/server.key -out cert/server.crt -days 3650
openssl req -new -sha256 -key cert/server.key -out cert/server.csr
openssl x509 -req -sha256 -in cert/server.csr -signkey cert/server.key -out cert/server.crt -days 3650
```

## See makefile for build instructions.
```
> make help

api                            Auto-generate grpc go sources
build_client                   Build the binary file for client
build_server                   Build the binary file for server
clean                          Remove previous builds
dep                            Get the dependencies
help                           Display this help screen
```

## Reference:
https://medium.com/pantomath/how-we-use-grpc-to-build-a-client-server-system-in-go-dd20045fa1c2
