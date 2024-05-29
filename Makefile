vet:
	go vet ./...

test:
	go test -vet=off -failfast -race -coverprofile=coverage.out ./...

atreugo:
	go run ./examples/atreugo

fiber:
	go run ./examples/fiber

http:
	go run ./examples/http