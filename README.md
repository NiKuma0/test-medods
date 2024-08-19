# Token Service

This project is a Go-based token service that generates JWT access tokens and refresh tokens. It also validates tokens and saves refresh tokens to a PostgreSQL database. The refresh tokens are hashed and can only be used once.

## Features

- Generate JWT access tokens
- Generate refresh tokens
- Validate access tokens
- Validate refresh tokens
- Save hashed refresh tokens to the database
- Ensure refresh tokens can only be used once

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/NiKuma0/test-medods.git
    cd test-medods
    ```

2. Install dependencies:

    ```sh
    go mod tidy
    ```

3. Set up your PostgreSQL database and update the DSN in [.env](.env)

## Running in dev mode:

1. Initialize the database:

    ```
    make migrate-dev
    ```

2. Run the application:

    ```sh
    make run
    ```

## API Endpoints

### Generate Tokens

**Endpoint:** `POST /token`

**Request Body:**

```json
{
    "user_id": "your-user-UUID"
}
```

**Returning body:**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.asdfasdfasfdasdfasdfasdfasd.EFYZG7Z9y7BPH-bblFjZTMOYbwGTX7GbRFjAwlygHhg",
  "refresh_token": "ODllOWZmNDgtMTk5Ny00YmRmLTkyOGItMWMyNzY3NWNkMDRm"
}
```

### 

**Endpoint:** `POST /refresh`

**Request Body:**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.asdfasdfasfdasdfasdfasdfasd.EFYZG7Z9y7BPH-bblFjZTMOYbwGTX7GbRFjAwlygHhg",
  "refresh_token": "ODllOWZmNDgtMTk5Ny00YmRmLTkyOGItMWMyNzY3NWNkMDRm"
}
```

**Returning body:**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.asdfasdfasfdasdfasdfasdfasd.EFYZG7Z9y7BPH-bblFjZTMOYbwGTX7GbRFjAwlygHhg",
  "refresh_token": "ODllOWZmNDgtMTk5Ny00YmRmLTkyOGItMWMyNzY3NWNkMDRm"
}
```