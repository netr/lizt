.PHONY: test bench

test:
	@echo "Running tests..."
	go test -v ./...
	@echo "Done."
lint:
	@echo "Running linter..."
	golangci-lint -p format -p error -p comment -p performance -p import -p metalinter run ./... --fix
	@echo "Done."
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...
	@echo "Done."