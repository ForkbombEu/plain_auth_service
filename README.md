# User Authentication Service

This is a simple Go application that exposes an HTTP endpoint for user authentication. It uses SQLite to store user credentials and verifies the username and password against the database.

## Features

- Exposes a single HTTP endpoint `/auth` for authentication.
- Validates user credentials stored in an SQLite database.
- Reads user credentials from a CSV file to populate the database.
- Performs double-check user/password verification.

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/forkbombeu/plain_auth_service.git
    cd plain_auth_service
    ```

2. Install the required Go packages:
    ```sh
    go mod tidy
    ```

## Running the Application

1. Run the application:
    ```sh
    go run main.go
    ```

The server will start on `localhost:8080`.

## Usage

Send a POST request to the `/auth` endpoint with a JSON body containing the username and password.

Example:

```bash
curl http://test.auth.forkbomb.eu:8080/auth \
    -d '{"username":"admin", "password":"admin123"}'
```

## Slangroom usage

#### `auth.slang`

```json 
{
  "address": "http://test.auth.forkbomb.eu:8080/auth",
  "auth": {
    "username":"admin",
    "password":"admin123"
  }
}
```
#### `auth.data.json`

```gherkin
Rule unknown ignore
Given I connect to 'address' and send object 'auth' and do post and output into 'result'
Given I have a 'string dictionary' named 'result'
Then print the 'result'
```
#### use with `slangroom-exec`

```bash
slexfe -F auth | slangroom-exec | jq
```

expected output:

```json
{
  "result": {
    "result": "Success",
    "status": "200"
  }
}
```

## Response

- `200 OK` with message `Success` if the credentials are correct.
- `401 Unauthorized` if the credentials are incorrect.

# License

Unlicensed, and written in 5min with the help of AI generated creative content

