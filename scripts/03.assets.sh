#!/bin/sh
set -eu

cd "$(dirname "$0")/.."

MINI_WORLD_FILE="${MINI_WORLD_FILE:-mini_world.json}"
export MINI_WORLD_FILE

go run ./cmd/mini init

go run ./cmd/mini create-central-bank ecb EUR
go run ./cmd/mini create-bank bank1
go run ./cmd/mini create-human alice
go run ./cmd/mini create-human bob

go run ./cmd/mini issue-cash ecb 1000
go run ./cmd/mini seed-cash ecb alice 700

go run ./cmd/mini register-asset gold1 ecb EUR 1000
go run ./cmd/mini register-asset bond1 bank1 EUR 750
go run ./cmd/mini register-asset painting1 bob EUR 500
go run ./cmd/mini buy-asset-cash alice painting1 450
go run ./cmd/mini revalue-asset painting1 550
