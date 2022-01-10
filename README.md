# Geographical Org Chart

## Intro
A Geo Based Org Chart for slalom written in Go PostgreSQL and Docker. 

## Setup
`$ docker-compose up --build`

in another terminal window run this from the project root

ensure you have go-migrate installed

`$ export POSTGRES_URI="postgres://geographer:go_figure_it_out@0.0.0.0:5432/orgchart?sslmode=disable"`
`migrate -database ${POSTGRES_URI} -path ./migrations up`

## Roadmap TODOs
* Add PostGIS support to the database and database container.
* Add a modern Web Frontend (React/Vue/Angular) to replace the Hacky one in place now.
* Add the ability to upload and download
* Moar tests
* Sessions stored in the DB
* Uploadable profile images
* JWT authentication