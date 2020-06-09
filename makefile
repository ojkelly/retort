
# Watch for source changes and reload
watch:
	watchman-make -p '**/*.go' -t dev

dev:
	rm debug.log || true
	go run example/hackernews/main.go

# Build the example app
demo:
	rm debug.log || true
	go run example/cmd/main.go


# Build the hn app
hn:
	go run example/hackernews/main.go


# Run the example app with the race detector
race:
	rm debug.log || true
	go run -race example/cmd/main.go 2>&1 | tee race.log

DOCS_PORT=6060
docs:
	godoc -http=":$(DOCS_PORT)" & open http://localhost:$(DOCS_PORT)/pkg/retort.dev/

# Serve the retort.dev website for development
# requires npm and npx on the system
# its just a static page, load it in a browser preview if you want
serve:
	npx serve .

# Install watchman on macOS
# also upgrade pywatchman, to one that works with python3
# https://github.com/facebook/watchman/issues/631#issuecomment-541255161
install-watchman-macos:
	pip install pywatchman
	brew install watchman
