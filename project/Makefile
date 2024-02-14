FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp

.PHONY: init_go_work up_back up_build down_back build_broker build_front start_front stop_front start_all stop_all

init_go_work:
	cd .. && go work init && go work use ./front-end ./broker-service ./authentication

## up: starts all containers in the background without forcing build
up_back:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down_back:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

## build_front: builds the frone end binary
build_front:
	@echo "Building front end binary..."
	cd ../front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"

## start: starts the front end
start_front: build_front
	@echo "Starting front end"
	cd ../front-end && ./${FRONT_END_BINARY} &

## stop: stop the front end
stop_front:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"

## all: runs both backend and frontend
start_all: up_back start_front

## stop_all: stops both backend and frontend
stop_all: stop_front down_back