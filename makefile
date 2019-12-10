
# Watch for source changes and reload
watch:
	watchman-make -p '**/*.go' -t dev

# Build the example app
dev:
	rm debug.log || true
	go run example/cmd/main.go

# Build the hn app
hn:
	go run example/hackernews/main.go


# Run the example app with the race detector
race:
	rm debug.log || true
	go run -race example/cmd/main.go 2>&1 | tee race.log


# Install watchman on macOS
# also upgrade pywatchman, to one that works with python3
# https://github.com/facebook/watchman/issues/631#issuecomment-541255161
install-watchman-macos:
	pip install pywatchman
	brew install watchman
