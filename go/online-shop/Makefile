include .env

# build: build all binaries
build: clean build_web build_api
	@echo "All binaries have been built"

# clean: clean all binaries and run go clean
clean:
	@echo "Cleaning..."
	@rm -f .build/*
	@go clean
	@echo "Cleaned"

# build_web: build the Web server
build_web:
	@echo "Building the Web server..."
	@go build -o .build/web ./cmd/web
	@echo "The Web server has been built"

# build_api: build the API server
build_api:
	@echo "Building the API server..."
	@go build -o .build/api ./cmd/api
	@echo "The API server has been built"

# start: start both the Web server and the API server
start: start_db start_web start_api
	@echo "All services have been started"

# stop: stop both the
stop: stop_db stop_web stop_api
	@echo "All services have been stopped"

# start_web: start the Web server
start_web: stop_web build_web
	@echo "Starting the Web server..."
	@env STRIPE_KEY="${STRIPE_KEY}" STRIPE_SECRET="${STRIPE_SECRET}" .build/web -port=${WEB_PORT} &
	@echo "The Web server is running"

# stop_web: stop the Web server
stop_web:
	@echo "Stopping the Web server..."
	@-pkill -SIGTERM -f "web -port=${WEB_PORT}"
	@echo "The Web server has been stopped"

# start_api: start the API server
start_api: stop_api build_api
	@echo "Starting the api server..."
	@env API_DSN="${API_DSN}" STRIPE_KEY="${STRIPE_KEY}" STRIPE_SECRET="${STRIPE_SECRET}" .build/api -port=${API_PORT} &
	@echo "The API server is running"

# stop_api: stop the API server
stop_api:
	@echo "Stopping the API server..."
	@-pkill -SIGTERM -f "api -port=${API_PORT}"
	@echo "The API server has been stopped"

# start_db: start the MySQL database
start_db: stop_db
	@echo "Starting the MySQL database..."
	@docker run -d --name ${MYSQL_NAME} -p 3306:3306 -v "$(PWD)/.data":"/var/lib/mysql" -e MYSQL_ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD}" mysql:${MYSQL_VERSION}
	@echo "The MySQL database is running"

# stop_db: stop the MySQL database
stop_db:
	@echo "Stopping the MySQL database..."
	@-docker stop ${MYSQL_NAME}
	@-docker rm ${MYSQL_NAME}
	@echo "The MySQL database has been stopped"

start_api_db: start_db start_api
	@echo "The MySQL database and API server are running"

stop_api_db: stop_db stop_api
	@echo "The MySQL database and API server have been stopped"
