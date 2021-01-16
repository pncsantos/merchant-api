# merchant-api

## Compiling
This is designed as a go module aware program and thus requires go 1.11 or better You can clone it anywhere, just run make inside the cloned directory to build.

## Requirements
This does require a mongodb database to be setup and reachable. It will attempt to create and migrate the database upon starting.

## Data Storage
Data is stored in a mongodb database by default.

## Run project by:
If you already have Docker Application in your machine, you can just simply run the application by:

docker-compose up --build

## Run test
go test -v ./... --run