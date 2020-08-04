# Currency Service
This is a simple test of GRPC, the idea is to enrich this server later.

# Installing dependencies
You will have to install **[protoc](https://developers.google.com/protocol-buffers/docs/downloads)** to generate the code and **[grpcurl](https://github.com/fullstorydev/grpcurl)** to test the code.

### Linux
```shell
sudo apt install protobuf-compiler
```

### Mac
```shell
brew install protoc
```

Then run the build command:

```shell
protoc -I protos/ protos/currency.proto --go_out=plugins=grpc:protos/currency
```
or
```
make protos
```

## Running the application
As this is a POC, all ports and configurations are hardcoded the app will run on port **:9098**

```shell script
go run main.go
```

With grpcurl installed we can test our server

### List Services
```
grpcurl --plaintext localhost:9092 list
Currency
grpc.reflection.v1alpha.ServerReflection
```

### List Methods
```
grpcurl --plaintext localhost:9098 list Currency        
Currency.GetRate
```

### Method detail for GetRate
```
grpcurl --plaintext localhost:9098 describe Currency.GetRate
Currency.GetRate is a method:
rpc GetRate ( .RateRequest ) returns ( .RateResponse );
```

### RateRequest detail
```
grpcurl --plaintext localhost:9098 describe .RateRequest    
RateRequest is a message:
message RateRequest {
  string Base = 1;
  string Destination = 2;
}
```

### Execute a request
```
grpcurl --plaintext -d '{"Base": "BR", "Destination": "USD"}' localhost:9098 Currency/GetRate
{
  "rate": 5.37
}
```
{
	"Base" : "BRL",
	"Destination" : "USD"
}

{
	"Base" : "BRL",
	"Destination" : "EUR"
}


 grpcurl --plaintext -d @ localhost:9098 Currency/SubscribeRates