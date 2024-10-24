# Stori Challenge

This project is a solution for the Stori Challenge. It processes CSV files containing transaction data, calculates summaries, and sends the results via email.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [Endpoints](#endpoints)
- [Environment Variables](#environment-variables)
- [License](#license)

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/stori-challenge.git
    cd stori-challenge
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    ```

3. Set up your environment variables (see [Environment Variables](#environment-variables)).

## Usage

To run the project locally:

First, create the build
```sh
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go
```
Next, start the sam server:
```sh
sam local start-api --log-file logfile.txt
```
Finally, invoke the lambda function:
```sh
sam invoke "storiCsvLoader" -l logfile.txt 
```

## Project Structure

```plaintext
stori-challenge/
├── application/
│   └── summary.go
├── domain/
│   ├── entities/
│   │   └── transaction.go
│   └── repositories/
│       ├── transaction_repository.go
│       └── pg_transaction_repository.go
├── infrastructure/
│   └── utils/
│       └── csv_reader.go
├── controller/
│   └── handler.go
├── main.go
└── go.mod
```

## Endpoints

- **POST /process-csv**: Processes a CSV file and returns a summary of the transactions.
```sh
curl --location 'http://localhost:3000/storiCsvLoader' \
--form 'transactions=@"/your-path/stori-challenge-v1/transactions.csv"' \
--form 'destination="davidsaldarriaga.pardo@gmail.com"'
```
## Environment Variables

The following environment variables need to be set:

- `DB_HOST`: Database host
- `DB_PORT`: Database port
- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `APP_PASSWORD_SMTP` : Password generated for use the SMTP server 
- `SENDER_EMAIL` : Sender email (davidsaldarriaga.pardo@gmail.com)

