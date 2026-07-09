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
