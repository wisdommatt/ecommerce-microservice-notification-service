protoc:
	protoc proto/email.proto --go-grpc_out=. --go_out=.
	
run:
	go run server/main.go

tests:
	go test ./... -race -cover

watch:
	go install github.com/cespare/reflex@latest
	reflex -s -- sh -c 'clear && APP_ENV=dev PORT=2020 go run server/main.go'