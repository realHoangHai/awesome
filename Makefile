.PHONY: init
# init env
init:
	go install \
        github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
        github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
        google.golang.org/protobuf/cmd/protoc-gen-go \
        google.golang.org/grpc/cmd/protoc-gen-go-grpc

.PHONY: gen.error gen.api gen.wire

gen.error:
	protoc --proto_path=./pkg/errors \
		   --proto_path=./third_party \
		   --go_out=paths=source_relative:./pkg/errors \
		   --go-grpc_out=paths=source_relative:./pkg/errors \
		   pkg/errors/errors.proto

gen.api:
	protoc --proto_path=./api/user/v1 \
		   --proto_path=./third_party \
		   --go_out=paths=source_relative:./api/user/v1 \
		   --go-errors_out=paths=source_relative:./api/user/v1 \
		   --go-grpc_out=paths=source_relative:./api/user/v1 \
		   --grpc-gateway_out ./api/user/v1 \
		   --grpc-gateway_opt logtostderr=true \
		   --grpc-gateway_opt paths=source_relative \
		   --grpc-gateway_opt generate_unbound_methods=true \
		   --openapiv2_out ./api/user/v1 \
           --openapiv2_opt logtostderr=true \
           --openapiv2_opt json_names_for_fields=false \
		   api/user/v1/*.proto


# generate
gen.wire:
	wire ./cmd