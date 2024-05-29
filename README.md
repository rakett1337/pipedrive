# Pipedrive Deals API Proxy

This is a Golang application that serves as a proxy for the Pipedrive Deals API. It provides endpoints to retrieve, create, and update deals in Pipedrive using the Pipedrive API.

# Prerequisites

To run this app you need to have at least one of the following installed on your machine:

- Docker
- Go

# Getting started

```bash
git clone https://github.com/rakett1337/pipedrive.git
cd pipedrive
```

# Running the app

## Using Docker

```bash
docker build -t pipedrive .
docker run -p 80:80 pipedrive
```

## Using Go

```bash
go run cmd/server/main.go
```

or build and run the binary

```bash
go build -o ./main cmd/server/main.go && ./main
```

# API Endpoints

The following API endpoints are available:

- GET /deals: Retrieves a list of deals from Pipedrive.
- POST /deals: Creates a new deal in Pipedrive.
- PUT /deals/:id: Updates an existing deal in Pipedrive.

To fully utilize the deals API, you need to obtain an API token from Pipedrive. If you have an account with Pipedrive, you can get your personal API token [here](https://app.pipedrive.com/settings/api)

As the instructions were to build a "simple application", any requests made to the /deals endpoint will be forwarded as-is to the Pipedrive API and the response will be returned to the client. No attempt to verify the request body or headers is made before forwarding the request.
Refer to the [Pipedrive API documentation](https://developers.pipedrive.com/docs/api/v1/Deals) for more information on how interact with the Pipedrive API.

### Example usage

```bash
curl -X GET http://localhost/deals -H "x-api-token: YOUR_API_TOKEN"
```

# Metrics

The app exposes metrics in the /metrics endpoint about the number of requests made to the Pipedrive API and the average response time of those requests.
Additionally all requests will be logged to the console. Metrics are not stored persistently and will be reset when the app is restarted.
This is handled by the metrics middleware which is applied to all requests.

```bash
curl http://localhost/metrics
```

# Testing

The deals and metrics handlers include unit tests. Deals handler test requires a valid Pipedrive API token to run, otherwise the tests will be failed. To run the tests, use the following command:

```bash
API_TOKEN=YOUR_API_TOKEN go test ./... -v
```

# GitHub Actions

The repository includes GitHub Actions workflows:

- Test and Lint: runs on every commit pushed to a pull request and checks if the code is properly formatted and all tests pass. API_TOKEN is required to run the tests and is stored as a secret in the repository.

- Deploy: runs only when a pull-request is merged to main.
