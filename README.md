# go-microservices-demo

## How to run

### Prerequisites
- Install [Docker](https://www.docker.com/get-started)
- Install [Docker Compose](https://docs.docker.com/compose/install/)
- Install [Go](https://golang.org/doc/install)
- Install [Make](https://www.gnu.org/software/make/)

### Run the application
- Clone the repository
- Run `make up` from the `project` directory in order to create docker images and run the backend services
- Run `make start` to start the frontend service
- Open your browser and navigate to `http://localhost:80`

### Stop the application
- Run `make down` from the `project` directory in order to stop and remove the docker containers
- Run `make stop` to stop the frontend service