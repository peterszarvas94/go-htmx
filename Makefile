# Define environment variables file
ENV_FILE := .env

# Include environment variables from .env file (if it exists)
ifneq (,$(wildcard $(ENV_FILE)))
	include $(ENV_FILE)
	export $(shell sed 's/=.*//' $(ENV_FILE))
endif

# Run the development server with gin
dev:
	gin -i --appPort 8080 --port 3000 run main.go

# Generate tailwindcss classes for development
tw-watch:
	tailwindcss -i tailwind.css -o static/style.css --watch
