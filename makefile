
dev:
	rm debug.log || true
	go run example/cmd/main.go

race:
	rm debug.log || true
	go run -race example/cmd/main.go 2>&1 | tee race.log