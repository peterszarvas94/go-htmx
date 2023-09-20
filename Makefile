dev:
	gin -i --appPort 8080 --port 3000 run main.go

tw-watch:
	tailwindcss -i tailwind.css -o static/style.css --watch
