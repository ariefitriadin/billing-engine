# Billing-Engine Test

## Overview
This project is a billing engine written in Go, using PostgreSQL as the database.

## Prerequisites
- **Go**: Ensure you have Go installed on your machine.
- **PostgreSQL**: Ensure you have PostgreSQL installed and running.

## Installation

### Install Docker
using docker run : 

```
docker-compose up -d   
```
for up and run postgresql


### Install dbmate
Follow the instructions to install dbmate from the [dbmate GitHub repository](https://github.com/amacneil/dbmate?tab=readme-ov-file#installation).

### Install sqlc
Follow the instructions to install sqlc from the [sqlc documentation](https://docs.sqlc.dev/en/stable/overview/install.html).

### Optional: Install Make
Make is optional but recommended for running build commands. Install it using your package manager.

## Running Locally
1. Clone the repository:
    git clone https://github.com/yourusername/billing-engine.git
    
    cd billing-engine


2. Run the following command to set up the project:

    make all


3. Import the Insomnia file `Insomnia_2024-05-20.json` into Insomnia for testing API endpoints.

## Using cURL

### Create New Loan
```
curl --request POST \
--url http://localhost:8080/loans \
--header 'Content-Type: application/json' \
--header 'User-Agent: insomnia/9.2.0' \
--data '{
"borrower_id": 2,
"amount": 3000000,
"interest_rate": 10,
"duration_weeks": 5
}
```

### Check if Loan is Delinquent

```
curl --request GET \
--url http://localhost:8080/loans/39/delinquent \
--header 'User-Agent: insomnia/9.2.0'
```

### Get Outstanding
```
curl --request GET \
  --url http://localhost:8080/loans/39/outstanding \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/9.2.0'
```

### Make Payment
```
curl --request POST \
  --url http://localhost:8080/loans/39/payment \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/9.2.0' \
  --data '{
	"amount": 183334
}'
```



