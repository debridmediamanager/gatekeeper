tailwind:
	npx tailwindcss -i ./web/input.css -o ./static/css/tailwind.css

run:
	go run ./cmd/gatekeeper/main.go
