.PHONY: dev-backend dev-frontend test build

dev-backend:
	cd backend && go run ./cmd/api

dev-frontend:
	cd frontend && npm run dev

test:
	cd backend && go test -v ./...
	cd frontend && npm run typecheck

build:
	cd backend && go build -o bin/copely-api ./cmd/api
	cd frontend && npm run build
