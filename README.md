## Generating Go Code from .proto Files
Run the following command to generate Go code:
```sh
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       pkg/apis/payment/payment.proto
```
