from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def find_account_id(
    world: dict, human_id: str, bank_id: str, currency: str
) -> str | None:
    for account_id in world["humans"][human_id]["bank_accounts"]:
        account = world["accounts"][account_id]

        if (
            account["owner_human_id"] == human_id
            and account["bank_id"] == bank_id
            and account["currency"] == currency
            and account["status"] == "ACTIVE"
        ):
            return account_id

    return None


def run(args: Namespace) -> None:
    world = load_world()

    human_id = args.human_id
    bank_id = args.bank_id
    loan_id = args.loan_id
    amount = args.amount

    if amount <= 0:
        print("Amount must be greater than 0.")
        return

    if human_id not in world["humans"]:
        print(f"Human does not exist: {human_id}")
        return

    if bank_id not in world["banks"]:
        print(f"Bank does not exist: {bank_id}")
        return

    if loan_id not in world["customer_loans"]:
        print(f"Loan does not exist: {loan_id}")
        return

    loan = world["customer_loans"][loan_id]

    if loan["bank_id"] != bank_id:
        print(f"Loan {loan_id} does not belong to {bank_id}.")
        return

    if loan["borrower_human_id"] != human_id:
        print(f"Loan {loan_id} does not belong to {human_id}.")
        return

    if loan["status"] != "ACTIVE":
        print(f"Loan is not active: {loan_id}")
        return

    currency = loan["currency"]
    outstanding_principal = loan["outstanding_principal"]
    outstanding_interest = loan["outstanding_interest"]
    total_due = outstanding_principal + outstanding_interest

    if total_due < amount:
        print(
            f"Repayment amount is greater than total due. "
            f"Total due: {total_due} {currency}"
        )
        return

    account_id = find_account_id(world, human_id, bank_id, currency)

    if account_id is None:
        print(f"No active {currency} account for {human_id} at {bank_id}.")
        return

    bank = world["banks"][bank_id]
    human = world["humans"][human_id]
    account = world["accounts"][account_id]

    account_balance = account["booked_balance"]

    if account_balance < amount:
        print(
            f"Not enough money in {human_id}'s account. "
            f"Available: {account_balance} {currency}"
        )
        return

    account["booked_balance"] -= amount

    remaining_payment = amount

    interest_paid = min(remaining_payment, loan["outstanding_interest"])
    loan["outstanding_interest"] -= interest_paid
    remaining_payment -= interest_paid

    principal_paid = min(remaining_payment, loan["outstanding_principal"])
    loan["outstanding_principal"] -= principal_paid
    remaining_payment -= principal_paid

    bank["interest_income"][currency] = (
        bank["interest_income"].get(currency, 0) + interest_paid
    )

    bank["equity"][currency] = bank["equity"].get(currency, 0) + interest_paid

    human["loans"][loan_id] = (
        loan["outstanding_principal"] + loan["outstanding_interest"]
    )

    if human["loans"][loan_id] == 0:
        loan["status"] = "REPAID"

    save_world(world)

    print(f"{human_id} repaid {amount} {currency} on {loan_id}.")
    print(f"Interest paid: {interest_paid} {currency}")
    print(f"Principal paid: {principal_paid} {currency}")
    print(f"{account_id} booked_balance: {account['booked_balance']} {currency}")
    print(
        f"{loan_id} outstanding_principal: {loan['outstanding_principal']} {currency}"
    )
    print(f"{loan_id} outstanding_interest: {loan['outstanding_interest']} {currency}")
    print(f"{human_id} loans[{loan_id}]: {human['loans'][loan_id]} {currency}")
    print(f"{loan_id} status: {loan['status']}")
    print(
        f"{bank_id} interest_income[{currency}]: {bank['interest_income'][currency]} {currency}"
    )
    print(f"{bank_id} equity[{currency}]: {bank['equity'][currency]} {currency}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "repay-loan",
        help="Repay a commercial bank loan using the borrower's deposit balance.",
    )

    parser.add_argument("human_id", help="Borrower human ID, for example: alice")

    parser.add_argument("bank_id", help="Commercial bank ID, for example: bank1")

    parser.add_argument("loan_id", help="Loan ID, for example: loan_000001")

    parser.add_argument("amount", type=int, help="Repayment amount.")

    parser.set_defaults(func=run)
