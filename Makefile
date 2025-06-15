.PHONY: tidy gazelle run-server proto

tidy:
	bazel run @rules_go//go -- mod tidy -v
	bazel mod tidy

gazelle:
	bazel run //:gazelle

run-server:
	bazel run //src/cmd/server:server

proto:
	@if [ -z "$$ARG" ]; then \
		echo "Usage: make proto ARG=<service>"; \
		exit 1; \
	fi
	protoc --proto_path=idl/$(ARG)/ \
	       --go_out=src/$(ARG)/proto/ --go_opt=paths=source_relative \
	       --go-grpc_out=src/$(ARG)/proto/ --go-grpc_opt=paths=source_relative \
	       idl/$(ARG)/$(ARG).proto

