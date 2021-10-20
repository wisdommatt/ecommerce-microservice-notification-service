run:
	go run server/main.go

tests:
	go test ./... -race -cover