#!/bin/sh
set -eu

cd "$(dirname "$0")/.."

MINI_WORLD_FILE="${MINI_WORLD_FILE:-mini_world.json}"
export MINI_WORLD_FILE

go run ./cmd/mini init

go run ./cmd/mini create-central-bank ecb EUR
go run ./cmd/mini create-central-bank fed USD
go run ./cmd/mini create-bank bank1
go run ./cmd/mini create-bank bank2
go run ./cmd/mini open-reserve-account ecb bank1
go run ./cmd/mini open-reserve-account ecb bank2
go run ./cmd/mini open-reserve-account fed bank1
go run ./cmd/mini open-reserve-account fed bank2

go run ./cmd/mini create-human alice
go run ./cmd/mini create-human bob
go run ./cmd/mini open-account alice bank1 EUR
go run ./cmd/mini open-account alice bank1 USD
go run ./cmd/mini open-account bob bank1 USD

go run ./cmd/mini issue-cash ecb 1000
go run ./cmd/mini seed-cash ecb alice 500
go run ./cmd/mini deposit-cash alice bank1 EUR 300

go run ./cmd/mini issue-cash fed 1000
go run ./cmd/mini seed-cash fed bob 300
go run ./cmd/mini deposit-cash bob bank1 USD 200

go run ./cmd/mini lend-reserves ecb bank1 EUR 300
go run ./cmd/mini lend-reserves fed bank2 USD 300

go run ./cmd/mini set-fx-rate EUR USD 1.1
go run ./cmd/mini fx-convert-deposit alice bank1 EUR USD 100
go run ./cmd/mini fx-convert-cash alice bank1 EUR USD 50
go run ./cmd/mini fx-bank-trade bank1 bank2 EUR USD 100

go run ./cmd/mini check-world
