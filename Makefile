vet:
	go vet ./...

test:
	go test -vet=off -failfast -race -coverprofile=coverage.out ./...