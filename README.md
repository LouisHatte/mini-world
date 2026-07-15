# mini-world

## Run

From this directory:

```sh
go run ./cmd/mini --help
go run ./cmd/mini init
```

## Test

Run the milestone scenario:

```sh
sh scripts/01.setup.sh
sh scripts/02.cash.sh
sh scripts/03.assets.sh
sh scripts/04.reserves.sh
sh scripts/05.asset-reserves.sh
sh scripts/06.payments.sh
sh scripts/07.loans.sh
sh scripts/08.sepa.sh
sh scripts/09.fx.sh
sh scripts/10.swift.sh
```

## Commands

```sh
# ---------------------------------------- WORLD ----------------------------------------

init                    Create empty world
check-world             Check mirror consistency

# ---------------------------------------- SETUP ----------------------------------------

create-central-bank     Create central bank
create-bank             Create commercial bank
create-human            Create human
open-account            Open a bank deposit account
open-reserve-account    Open a commercial bank reserve account

# ---------------------------------------- CASH ----------------------------------------

issue-cash              Create physical cash in a central bank vault
destroy-cash            Central bank removes physical cash from circulation

seed-cash               Give initial existing cash from central bank to a human

deposit-cash            Human deposits physical cash into a commercial bank
withdraw-cash           Human converts bank deposit into physical cash
transfer-cash           Human physically gives cash to another human

supply-cash             Commercial bank converts reserves into physical cash from central bank
return-cash             Commercial bank returns physical cash to central bank and receives reserves

move-cash               Physical cash transport between commercial banks, without reserve settlement
sell-cash               One commercial bank sells physical cash to another commercial bank, settled with reserves

# ---------------------------------------- ASSETS ----------------------------------------

register-asset          Setup/admin: register an already existing asset
buy-asset-cash          Human buys an existing asset with physical cash
buy-asset-reserves      Bank or central bank buys an asset, settled with reserves
revalue-asset           Change estimated asset value

# ---------------------------------------- RESERVES ----------------------------------------

lend-reserves           Central bank lends reserves, optionally with --collateral asset_id
repay-reserve-loan      Commercial bank repays central bank reserve loan
reserve-transfer        Transfer reserves between commercial banks

# ---------------------------------------- PAYMENTS ----------------------------------------

internal-transfer       Transfer deposits between two humans inside the same commercial bank
interbank-payment       Transfer deposits between humans at different banks, settled with reserves
pay                     High-level payment command, auto-detects internal or interbank payment

# ---------------------------------------- LOANS ----------------------------------------

grant-loan              Bank grants loan and creates deposit, optionally with --collateral asset_id
accrue-interest         Accrue interest on loans
repay-loan              Human repays loan
default-loan            Mark loan as defaulted

# ---------------------------------------- SEPA ----------------------------------------

sepa-credit-transfer    Create SEPA credit transfer
sepa-instant            Create and settle SEPA instant payment
settle-sepa             Settle SEPA transfers with reserves
reject-sepa             Reject a non-settled SEPA payment and refund sender

# ---------------------------------------- FX ----------------------------------------

set-fx-rate             Set exchange rate between two currencies
fx-convert-deposit      Human converts bank deposit from one currency to another through a bank
fx-convert-cash         Human exchanges physical cash from one currency to another through a bank
fx-bank-trade           Two banks exchange currencies, settled with reserves

# ---------------------------------------- SWIFT ----------------------------------------

open-correspondent-account   Open a correspondent account between two banks
fund-correspondent-account   Fund a correspondent account using reserves
swift-mt103                  Create SWIFT customer payment message
settle-swift                 Settle SWIFT payment through correspondent accounts
reject-swift                 Reject a non-settled SWIFT payment
```
