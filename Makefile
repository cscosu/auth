all: build

tailwind:
	@mkdir -p tmp
	@if [ ! -f tmp/tailwindcss ]; then curl -sL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 -o tmp/tailwindcss; fi
	@chmod +x tmp/tailwindcss

prepare: tailwind

build: tailwind
	@tmp/tailwindcss -i styles.css -o static/tailwind.css --minify
	@go build -o tmp/main .
