![Build Status](https://github.com/eddie023/bYd/actions/workflows/main.yml/badge.svg?branch=main)

# Bootstrap Your Dream (bYd)
Serverless starter kit built with Golang, AWS, Terraform. The goal of this project is to provide cost-effective enterprise grade starting point to launch a new project into production. Largely inspired by my own experience working on production systems.

## Running The Project Locally

### Connecting to local DB 
1. Docker compose up command will start a postgres db
2. export local DB config using `export DB_CONNECTION_URI="postgres://root:postgres@localhost:5432/postgres?sslmode=disable"` and run `make migrate-up` 
3. Connect to your postgres using psql command `psql --host localhost --port 5432 --user root --db postgres` and run your seed script. 
