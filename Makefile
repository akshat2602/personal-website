.PHONY: build serve test clean

build:
	go run cmd/build/main.go

serve: build
	python3 -m http.server 8080 --directory public

test: build
	@pgrep -f "http.server 8080" > /dev/null || python3 -m http.server 8080 --directory public &
	@sleep 1
	node test-ssg.js

clean:
	rm -rf public

watch:
	@echo "Watching for changes... (requires entr)"
	find content ssg-templates ssg-static internal cmd -type f | entr -r go run cmd/build/main.go
