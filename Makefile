lint:
	golangci-lint -p format -p error -p comment -p performance -p import -p metalinter run ./... --fix