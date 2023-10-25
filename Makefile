#!make

include .env
export $(grep -v '^#' .env | sed 's/=.*//' | xargs)

# Run the development server with gin
dev:
	gin -i --appPort 8080 --port 3000 run main.go

# Generate tailwindcss classes for development
tw:
	tailwindcss -i tailwind.css -o static/style.css --watch

# Build
build:
	tailwindcss -i tailwind.css -o static/style.css --minify
	docker build -t go-hmtx .

# Run
run:
	docker run -p 8080:8080 --env-file .env go-hmtx
