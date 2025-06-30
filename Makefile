# Start both client and server
dev:
	@echo "Go API Running on :8080"
	@echo "Svelte App running on :3003"
	@APP_ENV=development \
		go run ./server & \
		cd client && npm run dev -- --host & \
		wait

# backend only (prod mode will server SPA too)
server:
	go run ./server

client:
	cd client && npm run dev -- --host

test:
	go test ./... -v

build: tidy
	GOOS=linux GOARCH=amd64 go build -o bin/server ./server
	cd client && npm run build

tidy: 
	go vet ./...
	go mod tidy

.PHONY: dev server client test build tidy