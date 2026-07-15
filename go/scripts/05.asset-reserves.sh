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
go run ./cmd/mini open-account alice bank1 EUR

go run ./cmd/mini issue-cash ecb 1000
go run ./cmd/mini seed-cash ecb alice 600
go run ./cmd/mini deposit-cash alice bank1 EUR 500

go run ./cmd/mini register-asset collateral1 bank1 EUR 1000
go run ./cmd/mini lend-reserves ecb bank1 EUR 500 --collateral collateral1
go run ./cmd/mini reserve-transfer ecb bank1 bank2 EUR 300

go run ./cmd/mini register-asset bank1_bond bank1 EUR 100
go run ./cmd/mini register-asset bank1_bill bank1 EUR 100
go run ./cmd/mini register-asset ecb_bond ecb EUR 100

go run ./cmd/mini buy-asset-reserves ecb bank2 bank1_bond 100
go run ./cmd/mini buy-asset-reserves ecb ecb bank1_bill 100
go run ./cmd/mini buy-asset-reserves ecb bank1 ecb_bond 100

go run ./cmd/mini check-world
