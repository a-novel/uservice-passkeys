#!/bin/bash

COMPOSE_FILE="docker-compose.yaml"

PORT=4003
POSTGRES_PORT=5003

# Ensure containers are properly shut down when the program exits abnormally.
int_handler()
{
    docker compose -f ${COMPOSE_FILE} down
}
trap int_handler INT

export POSTGRES_PORT=${POSTGRES_PORT}

# Setup test containers.
docker compose -f ${COMPOSE_FILE} up -d

export ANOVEL_LOGGER_DYNAMIC=true
export PORT=${PORT}
export DSN="postgres://postgres:postgres@localhost:${POSTGRES_PORT}/postgres?sslmode=disable"

# Execute tests.
go run cmd/server/main.go

# Normal execution: containers are shut down.
docker compose -f ${COMPOSE_FILE} down
