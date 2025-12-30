.PHONY: dev build-frontend build-backend prod clean

dev:
	@echo "Starting development servers..."
	@cd backend && go run -tags dev . & \
	cd frontend && pnpm run dev

build-frontend:
	@echo "Building frontend..."
	@cd frontend && pnpm run build
	@echo "This file is only kept so go builds work" > backend/frontend/dist/embedable-file.txt

build-backend: build-frontend
	@echo "Building backend for Linux AMD64..."
	@cd backend && GOOS=linux GOARCH=amd64 go build -o app

prod: build-backend
	@echo "Running production server..."
	@cd backend && ./app

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf backend/app backend/frontend/dist
	@mkdir -p backend/frontend/dist
	@echo "This file is only kept so go builds work" > backend/frontend/dist/embedable-file.txt
