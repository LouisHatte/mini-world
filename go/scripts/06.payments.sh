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
go run ./cmd/mini create-human charlie
go run ./cmd/mini open-account alice bank1 EUR
go run ./cmd/mini open-account bob bank1 EUR
go run ./cmd/mini open-account charlie bank2 EUR

go run ./cmd/mini issue-cash ecb 1000
go run ./cmd/mini seed-cash ecb alice 700
go run ./cmd/mini deposit-cash alice bank1 EUR 600

go run ./cmd/mini register-asset bond1 bank1 EUR 1000
go run ./cmd/mini lend-reserves ecb bank1 EUR 500 --collateral bond1

go run ./cmd/mini internal-transfer alice bob bank1 EUR 100
go run ./cmd/mini interbank-payment alice bank1 charlie bank2 EUR 150
go run ./cmd/mini pay bob bank1 charlie bank2 EUR 50
go run ./cmd/mini pay alice bank1 bob bank1 EUR 25

go run ./cmd/mini check-world
