gen:
	protoc --proto_path=./helloworld \
		   --proto_path=./third_party \
		   --go_out=paths=source_relative:./helloworld \
		   --go-grpc_out=paths=source_relative:./helloworld \
		   --grpc-gateway_out ./helloworld \
		   --grpc-gateway_opt logtostderr=true \
		   --grpc-gateway_opt paths=source_relative \
		   --grpc-gateway_opt generate_unbound_methods=true \
		   helloworld/helloworld.proto