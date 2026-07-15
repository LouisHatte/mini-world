#!/bin/sh
set -eu

cd "$(dirname "$0")/.."

MINI_WORLD_FILE="${MINI_WORLD_FILE:-mini_world.json}"
export MINI_WORLD_FILE

go run ./cmd/mini init

go run ./cmd/mini create-central-bank ecb EUR
go run ./cmd/mini create-bank bank1
go run ./cmd/mini create-bank bank2
go run ./cmd/mini open-reserve-account ecb bank1
go run ./cmd/mini open-reserve-account ecb bank2

go run ./cmd/mini create-human alice
go run ./cmd/mini create-human bob
go run ./cmd/mini open-account alice bank1 EUR
go run ./cmd/mini open-account bob bank2 EUR

go run ./cmd/mini issue-cash ecb 1000
go run ./cmd/mini seed-cash ecb alice 500
go run ./cmd/mini deposit-cash alice bank1 EUR 400

go run ./cmd/mini lend-reserves ecb bank1 EUR 300

go run ./cmd/mini open-correspondent-account bank1 bank2 EUR
go run ./cmd/mini fund-correspondent-account corr_bank1_bank2_eur 200

go run ./cmd/mini swift-mt103 alice bank1 bob bank2 EUR 100
go run ./cmd/mini settle-swift payment_000001

go run ./cmd/mini swift-mt103 alice bank1 bob bank2 EUR 50
go run ./cmd/mini reject-swift payment_000002 COMPLIANCE_REJECTED

go run ./cmd/mini check-world
