build:
	go build -o /dev/null -v cmd/service/main.go

run: 
	docker compose up --build --force-recreate --remove-orphans

test:
	docker compose -f test.compose.yml up --build --force-recreate --abort-on-container-exit

lint: 
	docker run --rm \
	--volume $$(pwd):/src \
	--volume ~/.cache:/root./.cache \
	$$(docker build --quiet --file lint.Dockerfile .) 

.PHONY: lambda-build-rest-api
lambda-build-rest-api:
	GOOS=linux GOARCH=amd64 go build -o bin/rest-api/bootstrap cmd/lambda/rest-api/handler.go

.PHONY: lambda-build-post-signup-confirmation
lambda-build-post-signup-confirmation:
	GOOS=linux GOARCH=amd64 go build -o bin/post-signup/bootstrap cmd/lambda/post-signup-confirmation/handler.go

generate:
	go generate ./...

migrate-up:
	migrate -source file://migrations -database ${DB_CONNECTION_URI} up

migrate-down:
	migrate -source file://migrations -database ${DB_CONNECTION_URI} down 

get-posts:
	curl curl http://localhost:8000/v1/posts 