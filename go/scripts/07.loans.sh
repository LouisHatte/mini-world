#!/bin/sh
set -eu

cd "$(dirname "$0")/.."

MINI_WORLD_FILE="${MINI_WORLD_FILE:-mini_world.json}"
export MINI_WORLD_FILE

go run ./cmd/mini init

go run ./cmd/mini create-central-bank ecb EUR
go run ./cmd/mini create-bank bank1
go run ./cmd/mini create-human alice
go run ./cmd/mini open-account alice bank1 EUR

go run ./cmd/mini issue-cash ecb 200
go run ./cmd/mini seed-cash ecb alice 100
go run ./cmd/mini deposit-cash alice bank1 EUR 100

go run ./cmd/mini register-asset alice_house alice EUR 1500
go run ./cmd/mini grant-loan bank1 alice EUR 1000 --collateral alice_house
go run ./cmd/mini accrue-interest loan_000001 50
go run ./cmd/mini repay-loan alice bank1 loan_000001 50
go run ./cmd/mini repay-loan alice bank1 loan_000001 1000

go run ./cmd/mini grant-loan bank1 alice EUR 500
go run ./cmd/mini accrue-interest loan_000002 25
go run ./cmd/mini default-loan loan_000002

go run ./cmd/mini check-world
