# go-microservices-demo

## How to run

### Prerequisites
- Install [Docker](https://www.docker.com/get-started)
- Install [Docker Compose](https://docs.docker.com/compose/install/)
- Install [Go](https://golang.org/doc/install)
- Install [Make](https://www.gnu.org/software/make/)

### Run the application
- Clone the repository
- Run `make start_all` from the `project` directory in order to create docker images and run the backend services
- Open your browser and navigate to `http://localhost:80`

### Stop the application
- Run `make stop_all` from the `project` directory in order to stop and remove the docker containers