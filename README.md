# lena

A session management service

## How to modify port, endpoint, and storage

```
$ export LENA_PORT=5454 // any number within the range of port
$ export LENA_ENDPOINT=http // http | grpc
$ export LENA_STORAGE=inmemory // inmemory | sqlite
```

## How to generate go files from the proto file

```
$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative endpoints/grpcendpoint/lena.proto
```
