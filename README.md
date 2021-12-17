# Geographical Org Chart

## Setup
`$ docker-compose up --build`

in another terminal window run this from the project root

ensure you have go-migrate installed

`$ export POSTGRES_URI="postgres://geographer:go_figure_it_out@0.0.0.0:5432/orgchart?sslmode=disable"`
`migrate -database ${POSTGRES_URI} -path ./migrations up`
