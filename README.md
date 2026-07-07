## Personnal notes

```json
{
  "ecb": {
    // I owe Bank 1 €5,000 in reserves. (1) (liabilities)
    "reserve_accounts": {
      "bank1": 5000
    },
    // Bank 1 owes me €5,000 as a loan. (2) (assets)
    "loans_to_banks": {
      "bank1": 5000
    }
  }
}
```

```json
{
  "bank1": {
    // I own €5,000 in reserves at ECB. (1) (assets)
    "reserve_balances": {
      "ecb": 5000
    },
    // I owe ECB €5,000 as a loan. (2) (liabilities)
    "loans_from_central_banks": {
      "ecb": 5000
    }
  }
}
```

## Commands

```sh
# ---------------------------------------- WORLD ----------------------------------------
init                    Create empty world
check-world             Check mirror consistency

# ---------------------------------------- SETUP ----------------------------------------
create-central-bank     Create ECB/Fed/etc.
create-bank             Create commercial bank
create-human            Create Alice/Bob
open-account            Open a bank deposit account

# ---------------------------------------- CASH ----------------------------------------
issue-cash              Create physical cash
destroy-cash            Central bank removes damaged cash from circulation

seed-cash               Give initial existing cash to a human
deposit-cash            Human deposits cash into a bank
withdraw-cash           Converts bank deposit into physical cash
transfer-cash           A human physically gives cash to another human

supply-cash             Convert bank reserves into physical cash
return-cash             Bank returns physical cash to central bank

move-cash               Physical cash transport between banks

# ---------------------------------------- RESERVES ----------------------------------------
lend-reserves           Create reserves by lending to a bank
repay-reserves          Bank repays central bank loan and destroys reserves
settle-reserves         Move reserves from one bank to another

check-reserves          Show Bank 1’s reserve position
reserve-shortfall       Test whether a bank can settle a payment

# ---------------------------------------- PAYMENTS ----------------------------------------
internal-transfer       Transfer inside the same bank

# ---------------------------------------- LOANS ----------------------------------------
grant-loan              Bank creates loan asset and deposit liability
repay-loan              Human repays part of loan
accrue-interest         Add interest to loan balance
default-loan            Mark loan as defaulted
write-off-loan          Remove bad loan from bank assets

# TODO

Commercial bank deposits and internal transfers
Command Purpose
place-hold alice bank1 EUR 100 Reserve part of Alice’s available balance
release-hold hold_id Cancel a hold
capture-hold hold_id Turn a hold into a booked debit
close-account account_id Close a bank account
freeze-account account_id Block account activity
unfreeze-account account_id Restore account activity

SEPA / STEP2 / T2
Command Purpose
create-step2 step2 EUR Create SEPA CSM infrastructure
sepa-transfer alice bank1 bob bank2 EUR 100 Full SEPA payment flow
step2-receive message_id STEP2 receives Bank 1 message
step2-clear payment_id STEP2 validates/routes/clears
t2-settle payment_id T2 moves central bank money
step2-deliver payment_id STEP2 delivers settled payment to Bank 2
bank-credit-incoming bank2 payment_id Bank 2 credits Bob
sepa-return payment_id reason Return/reject a SEPA payment
sepa-reconcile payment_id Compare Bank 1, STEP2, T2, Bank 2 states

SWIFT / correspondent banking
Command Purpose
create-correspondent-account bank1 bankX USD Bank 1 opens nostro/vostro relationship
swift-transfer alice bank1 bob bank2 USD 100 International SWIFT-style transfer
swift-message bank1 bank2 MT103 Create SWIFT payment instruction
correspondent-debit bankX bank1 USD 100 Correspondent bank debits nostro account
correspondent-credit bankY bank2 USD 100 Correspondent chain credits recipient bank
deduct-fee bankX payment_id 10 Intermediary bank takes fee
fx-convert bank1 EUR USD 100 FX conversion
swift-reconcile payment_id Check message chain vs ledger movements

Bonds / securities
Command Purpose
issue-bond issuer bond1 EUR 1000 Government/company issues a bond
buy-security ecb bank1 bond1 1000 Central bank buys bond, creates reserves
sell-security ecb bank1 bond1 1000 Central bank sells bond, destroys reserves
pay-coupon bond1 holder 50 Bond pays interest
redeem-bond bond1 holder 1000 Bond principal is repaid
mark-to-market bond1 950 Change market value of security

Cheques
Command Purpose
write-cheque alice bank1 bob EUR 100 Alice writes cheque to Bob
deposit-cheque bob bank2 cheque_id Bob deposits cheque at Bank 2
clear-cheque cheque_id Interbank cheque clearing
settle-cheque cheque_id Central bank settlement
bounce-cheque cheque_id Insufficient funds / invalid cheque
reverse-cheque-credit cheque_id Remove provisional credit from Bob

Cards / holds
Command Purpose
card-authorize alice merchant EUR 50 Create authorization hold
card-capture auth_id Finalize card payment
card-reverse auth_id Cancel authorization
card-chargeback payment_id Reverse disputed card payment

FX and cross-border
Command Purpose
create-fx-market EUR USD 1.10 Define FX rate
set-fx-rate EUR USD 1.08 Update FX rate
fx-convert alice bank1 EUR USD 100 Convert account balance
cross-border-transfer alice bank1 bob bank2 EUR USD 100 Transfer with FX
nostro-transfer bank1 bankX USD 100 Move funds via correspondent account

Failure / retry simulation
Command Purpose
simulate-network-failure message_id Message not delivered
retry-message message_id Retry outbox message
duplicate-message message_id Test idempotency
fail-settlement payment_id T2 settlement fails
reject-recipient payment_id Bank 2 rejects incoming payment
reconcile-all Search for inconsistent states
repair-mismatch Manual correction command

Audit / debugging
Command Purpose
show-json Pretty-print the raw JSON
show-ledger entity_id Show ledger entries for one entity
show-balance-sheet bank1 Show assets/liabilities/equity
show-account account_id Show one account
show-payment payment_id Show payment state machine
list-banks List banks
list-humans List humans
list-accounts List accounts
snapshot name Save current world state
load-snapshot name Restore previous state
```
