# mini-world

## Run

From this directory:

```sh
go run ./cmd/mini --help
MINI_WORLD_FILE=../mini_world.json go run ./cmd/mini init
```

## Test

Run the milestone scenario:

```sh
sh scripts/01.setup.sh
sh scripts/02.cash.sh
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
```

## TODO

```sh
# ---------------------------------------- WORLD ----------------------------------------

show-world Show world summary
list-entities List all entities

# ---------------------------------------- SETUP ----------------------------------------

create-currency Create currency

# ---------------------------------------- CASH ----------------------------------------

show-cash Show cash balances by holder
check-cash Check cash consistency

# ---------------------------------------- DEPOSITS ----------------------------------------

show-deposits Show all bank deposit liabilities
show-account Show one human bank account

# ---------------------------------------- RESERVES ----------------------------------------

issue-reserves Central bank creates reserves for a commercial bank
reserve-transfer Transfer reserves between commercial banks
show-reserves Show reserve balances
set-reserve-requirement Set minimum reserve ratio
check-reserve-ratio Check commercial bank reserve ratio

# ---------------------------------------- PAYMENTS ----------------------------------------

pay Human pays another human, auto internal/interbank
internal-transfer Transfer deposits inside same bank
interbank-payment Transfer between humans at different banks
show-payments Show payment history

# ---------------------------------------- LOANS ----------------------------------------

grant-loan Bank grants loan and creates deposit
show-loans Show loans
accrue-interest Accrue interest on loans
repay-loan Human repays loan
default-loan Mark loan as defaulted

# ---------------------------------------- COLLATERAL ----------------------------------------

create-collateral Create collateral owned by a human
pledge-collateral Attach collateral to a loan
release-collateral Release collateral after repayment
seize-collateral Bank seizes collateral after default
show-collateral Show collateral registry

# ---------------------------------------- SEPA ----------------------------------------

sepa-credit-transfer Create SEPA credit transfer
sepa-instant Create instant SEPA transfer
show-sepa-messages Show SEPA messages
settle-sepa Settle SEPA transfers with reserves

# ---------------------------------------- FX ----------------------------------------

set-fx-rate Set exchange rate
show-fx-rates Show exchange rates
fx-convert Convert human deposit between currencies
fx-bank-trade FX trade between banks
fx-revalue Revalue FX positions

# ---------------------------------------- SWIFT ----------------------------------------

create-correspondent-account Create correspondent account between banks
swift-mt103 Create SWIFT MT103 payment message
show-swift-messages Show SWIFT messages
settle-swift Settle SWIFT payment
show-correspondents Show correspondent banking links
```
