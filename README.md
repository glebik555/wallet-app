# Wallet App - Wallet CRUD and Transaction RESTful API Service

## Intro
Implementation of a RESTful Service API that provides endpoints to manage wallets and perform basic operations such as deposits and withdrawals. The service ensures data consistency and supports transactional operations.

## Features

- Create/Get wallet balances
- Deposit and Withdraw operations with transactional safety
- Graceful error handling for insufficient funds
- REST API endpoints:
  - `POST /api/v1/wallet` - perform deposit/withdraw operations
  - `GET /api/v1/wallets/{uuid}` - get wallet balance

## Environment Variables

Configurable via `.env` file or Docker environment variables:

- PORT
- PG_HOST
- PG_PORT
- PG_USER
- PG_PASSWORD
- PG_DB
- PG_SSLMODE
- MAX_DB_CONNS


## Database

- PostgreSQL 15
- SQL migrations in `internal/db/migrations/001_init.sql`
- Table `wallets` with columns:
  - `uuid` (primary key)
  - `balance`
  - `created_at`
  - `updated_at`

## Running with Docker Compose

```bash
# Start the application and database
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop the services
docker-compose down

API Example
Deposit
POST /api/v1/wallet
Content-Type: application/json

{
  "walletId": "123e4567-e89b-12d3-a456-426614174000",
  "operationType": "DEPOSIT",
  "amount": 1000
}


Get Balance
GET /api/v1/wallets/123e4567-e89b-12d3-a456-426614174000

Response:
{
  "walletId": "123e4567-e89b-12d3-a456-426614174000",
  "balance": 1000
}


