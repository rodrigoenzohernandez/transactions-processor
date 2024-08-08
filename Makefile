build-lambdas:
	@find internal/lambda -type d -mindepth 1 -maxdepth 1 | while read dir; do \
		echo "Building $$dir..."; \
		(cd $$dir && GOOS=linux GOARCH=amd64 go build -o bootstrap main.go) || exit 1; \
	done
deploy:
	make build-lambdas && cd infrastructure && cdk deploy --verbose --require-approval never
destroy:
	cd infrastructure && cdk destroy --force
diff:
	cd infrastructure && cdk diff