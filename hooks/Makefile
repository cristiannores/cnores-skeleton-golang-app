setup-hooks:
	@echo "\033[0;33mCreting commit message validation hook\033[0m"
	@cp ./hooks/commit-msg .git/hooks/commit-msg
	@chmod +x .git/hooks/commit-msg
	@echo "\033[0;33mCreting pre-commit validation hook\033[0m"
	@cp ./hooks/pre-commit .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "\033[0;33mCreting pre-push validation hook\033[0m"
	@cp ./hooks/pre-push .git/hooks/pre-push
	@chmod +x .git/hooks/pre-push

unit-test:
	go test -v ./...

unit-test-coverage:
	go test -v -cover ./...

unit-test-report:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

run-dev:
	go run ./main.go

mock:
	@chmod +x ./hooks/mockgenerator
	@bash ./hooks/mockgenerator $(file) $(interface)

mocks:
	@chmod +x ./hooks/find-and-create-mocks ./hooks/mockgenerator
	@bash ./hooks/find-and-create-mocks
