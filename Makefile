default:
	@go build -o build/kafka-scout cmd/serve/serve.go && cd ui && yarn build

serve:
	@./build/kafka-scout

run:
	@go run cmd/serve/serve.go
