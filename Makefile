.PHONY: protos

protos:
	protoc -I protos/ protos/currency.proto --go_out=plugins=grpc:protos/currency
	protoc -I protos/ protos/healthcheck.proto --go_out=plugins=grpc:protos/healthcheck
