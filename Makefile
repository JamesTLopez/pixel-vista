run: build
	@./bin/pixelvista

install:

	# @go install github.com/a-h/templ/cmd/templ@latest
	@go get github.com/a-h/templ/cmd/templ@latest
	@go get ./...
	@go mod vendor
	@go mod tidy
	@go mod download
	@npm install -D tailwindcss
	@npm install -D daisyui@latest

build:
	npx tailwindcss -i view/css/input.css -o public/styles.css
	@templ generate view
	@go build -o bin/pixelvista main.go

up:
	@go run cmd/migrate/main.go up

reset:
	@go run cmd/reset/main.go up
	@go run cmd/migrate/main.go up
down:
	@go run cmd/migrate/main.go down

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

seed:
	@go run cmd/seed/main.go
