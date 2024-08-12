build-lambdas:
	@find internal/lambda -type d -mindepth 1 -maxdepth 1 | while read dir; do \
		echo "Building $$dir..."; \
		(cd $$dir && GOOS=linux GOARCH=amd64 go build -o bootstrap main.go) || exit 1; \
	done
deploy:
	make build-lambdas
	cd infrastructure && cdk deploy --require-approval never
	make migrate-up

destroy:
	cd infrastructure && cdk destroy --force
diff:
	cd infrastructure && cdk diff

run-dev:
	aws s3 cp files/txns.csv s3://infrastructurestack-transactionsbucket77a27bfc-mpzvy34bxssi/

update-email:
	@if [ -z "$(email)" ]; then \
		echo "Usage: make update-email email=<new-email@example.com>"; \
		exit 1; \
	fi
	curl -X PUT https://41slfl71z8.execute-api.us-east-2.amazonaws.com/develop/param \
		-H "Content-Type: application/json" \
		-d '{"notificationEmail": "$(email)"}'

migrate-up:
	go run db/scripts/migrate.go up

migrate-down:
	go run db/scripts/migrate.go down
