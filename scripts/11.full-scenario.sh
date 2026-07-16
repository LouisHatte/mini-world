#!/bin/sh
set -eu

cd "$(dirname "$0")/.."

MINI_WORLD_FILE="${MINI_WORLD_FILE:-mini_world.json}"
export MINI_WORLD_FILE

go run ./cmd/mini init

# Setup
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
go run ./cmd/mini create-human charlie
go run ./cmd/mini open-account alice bank1 EUR
go run ./cmd/mini open-account alice bank1 USD
go run ./cmd/mini open-account bob bank1 EUR
go run ./cmd/mini open-account bob bank1 USD
go run ./cmd/mini open-account bob bank2 EUR
go run ./cmd/mini open-account charlie bank2 EUR

# Cash
go run ./cmd/mini issue-cash ecb 5000
go run ./cmd/mini issue-cash fed 3000
go run ./cmd/mini seed-cash ecb alice 2000
go run ./cmd/mini seed-cash ecb bob 500
go run ./cmd/mini seed-cash fed bob 1000
go run ./cmd/mini transfer-cash alice bob EUR 100
go run ./cmd/mini deposit-cash alice bank1 EUR 1200
go run ./cmd/mini deposit-cash bob bank1 EUR 200
go run ./cmd/mini deposit-cash bob bank1 USD 500
go run ./cmd/mini withdraw-cash alice bank1 EUR 100
go run ./cmd/mini move-cash bank1 bank2 EUR 50

# Reserves and central-bank cash logistics
go run ./cmd/mini lend-reserves ecb bank1 EUR 1500
go run ./cmd/mini lend-reserves ecb bank2 EUR 300
go run ./cmd/mini lend-reserves fed bank1 USD 500
go run ./cmd/mini lend-reserves fed bank2 USD 500
go run ./cmd/mini return-cash ecb bank2 25
go run ./cmd/mini supply-cash ecb bank2 10
go run ./cmd/mini sell-cash ecb bank1 bank2 10
go run ./cmd/mini reserve-transfer ecb bank1 bank2 EUR 200
go run ./cmd/mini repay-reserve-loan bank1 reserve_loan_ecb_bank1_1 100

# Assets
go run ./cmd/mini register-asset gold1 ecb EUR 1000
go run ./cmd/mini register-asset bank1_bond bank1 EUR 100
go run ./cmd/mini register-asset bank1_bill bank1 EUR 100
go run ./cmd/mini register-asset ecb_bond ecb EUR 100
go run ./cmd/mini register-asset painting1 bob EUR 500
go run ./cmd/mini register-asset alice_house alice EUR 1500
go run ./cmd/mini buy-asset-cash alice painting1 300
go run ./cmd/mini revalue-asset painting1 550
go run ./cmd/mini buy-asset-reserves ecb bank2 bank1_bond 100
go run ./cmd/mini buy-asset-reserves ecb ecb bank1_bill 100
go run ./cmd/mini buy-asset-reserves ecb bank1 ecb_bond 100

# Payments
go run ./cmd/mini internal-transfer alice bob bank1 EUR 100
go run ./cmd/mini interbank-payment alice bank1 charlie bank2 EUR 150
go run ./cmd/mini pay bob bank1 charlie bank2 EUR 50
go run ./cmd/mini pay alice bank1 bob bank1 EUR 25

# Loans
go run ./cmd/mini grant-loan bank1 alice EUR 1000 --collateral alice_house
go run ./cmd/mini accrue-interest loan_000001 50
go run ./cmd/mini repay-loan alice bank1 loan_000001 50
go run ./cmd/mini repay-loan alice bank1 loan_000001 300
go run ./cmd/mini grant-loan bank1 alice EUR 500
go run ./cmd/mini accrue-interest loan_000002 25
go run ./cmd/mini default-loan loan_000002

# SEPA
go run ./cmd/mini sepa-credit-transfer alice bank1 bob bank2 EUR 100
go run ./cmd/mini settle-sepa payment_000005
go run ./cmd/mini sepa-instant bob bank2 alice bank1 EUR 40
go run ./cmd/mini sepa-credit-transfer alice bank1 bob bank2 EUR 25
go run ./cmd/mini reject-sepa payment_000007 ACCOUNT_CLOSED

# FX
go run ./cmd/mini set-fx-rate EUR USD 1.1
go run ./cmd/mini fx-convert-deposit alice bank1 EUR USD 100
go run ./cmd/mini fx-convert-cash alice bank1 EUR USD 50
go run ./cmd/mini fx-bank-trade bank1 bank2 EUR USD 100

# SWIFT
go run ./cmd/mini open-correspondent-account bank1 bank2 EUR
go run ./cmd/mini fund-correspondent-account corr_bank1_bank2_eur 200
go run ./cmd/mini swift-mt103 alice bank1 bob bank2 EUR 100
go run ./cmd/mini settle-swift payment_000008
go run ./cmd/mini swift-mt103 alice bank1 bob bank2 EUR 50
go run ./cmd/mini reject-swift payment_000009 COMPLIANCE_REJECTED

go run ./cmd/mini destroy-cash ecb 100
go run ./cmd/mini check-world
