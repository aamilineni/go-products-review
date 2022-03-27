test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

