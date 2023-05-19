include .env

## build: builds all binaries
build: clean build_web build_api
	@echo "All binaries have been built"

## clean: cleans all binaries and runs go clean
clean:
	@echo "Cleaning..."
	@rm -f .build/*
	@go clean
	@echo "Cleaned"

## build_front: builds the front end
build_web:
	@echo "Building the Web server..."
	@go build -o .build/web ./cmd/web
	@echo "The Web server has been built"

## build_back: builds the back end
build_api:
	@echo "Building the API server..."
	@go build -o .build/api ./cmd/api
	@echo "The API server has been built"

## start: starts front and back end
start: start_web start_api
	@echo "All servers have been started"

## start_front: starts the front end
start_web: stop_web build_web
	@echo "Starting the Web server..."
	@env STRIPE_KEY=${STRIPE_KEY} STRIPE_SECRET=${STRIPE_SECRET} .build/web -port=${WEB_PORT} &
	@echo "The Web server is running"

## start_back: starts the back end
start_api: stop_api build_api
	@echo "Starting the api server..."
	@env STRIPE_KEY=${STRIPE_KEY} STRIPE_SECRET=${STRIPE_SECRET} .build/api -port=${API_PORT} &
	@echo "The API server is running"

## stop: stops the front and back end
stop: stop_web stop_api
	@echo "All servers have been stopped"

## stop_front: stops the front end
stop_web:
	@echo "Stopping the Web server..."
	@pkill -SIGTERM -f "web -port=${WEB_PORT}" || true
	@echo "The Web server has been stopped"

## stop_back: stops the back end
stop_api:
	@echo "Stopping the API server..."
	@pkill -SIGTERM -f "api -port=${API_PORT}" || true
	@echo "The API server has been stopped"
