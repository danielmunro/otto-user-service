# Otto User Service

What is it?

An open source identity and access management microservice. Written in go, ready to run in AWS Lambda and backend by 
AWS Cognito.

## Running the server
To run the server, follow these steps:

Create an .env file:
```.env
COGNITO_CLIENT_ID=
COGNITO_CLIENT_SECRET= # if a secret is defined
USER_POOL_ID=
AWS_REGION=
```

Start the Postgres database (`-d` runs it as a background process):
```
docker-compose up -d
```

Run the server locally:
```
go run main.go
```

Run the tests:
```
go test ./internal/...
```

## Development

Generating models:
```
./bin/swagger-generate-models
```

## Todo

* better error handling
* groups
* related entities (email, password)
* versioned docs
* recruit contributors

