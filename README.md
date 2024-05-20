# Billing-Engine Test

## Overview
This project is a billing engine written in Go, using PostgreSQL as the database.

## Prerequisites
- Docker

## Get Up and Running

### Install Docker if not installed yet


### Run using docker-compose : 

```
docker-compose up
```

- wait until eveything up and running 

- Import the Insomnia file `Insomnia_2024-05-20.json` into Insomnia for testing API endpoints.

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



