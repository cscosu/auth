all: build

tailwind:
	@mkdir -p tmp

prepare: tailwind

build: tailwind
	@go build -o tmp/main .
