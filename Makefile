unit-tests:
	go test -v ./... -coverprofile=coverage_report_from_container.out

integration-tests:
	go test -v ./tests/integration/api_tests.go

all: unit-tests integration-tests