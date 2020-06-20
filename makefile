
# Watch for source changes and reload
watch:
	watchman-make -p '**/*.go' -t dev

dev:
	rm debug.log || true
	go run example/hackernews/main.go

test:
	go test ./...

# go get github.com/dnephin/filewatcher
watch-demo:
	filewatcher -x .git make demo

# Build the example app
demo:
	rm debug.log || true
	go run example/cmd/main.go


# Build the hn app
hn:
	rm debug.log || true
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


prepare-site:
	find ./** -type d -exec cp redirect.html {}/index.html \;

remove-redirect-html:
	rm -rf ./**/index.html
