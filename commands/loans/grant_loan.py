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


def build_loan_id(world: dict) -> str:
    next_number = len(world["customer_loans"]) + 1
    return f"loan_{next_number:06d}"


def run(args: Namespace) -> None:
    world = load_world()

    bank_id = args.bank_id
    human_id = args.human_id
    currency = args.currency.upper()
    amount = args.amount

    if amount <= 0:
        print("Amount must be greater than 0.")
        return

    if bank_id not in world["banks"]:
        print(f"Bank does not exist: {bank_id}")
        return

    if human_id not in world["humans"]:
        print(f"Human does not exist: {human_id}")
        return

    account_id = find_account_id(world, human_id, bank_id, currency)

    if account_id is None:
        print(f"No active {currency} account for {human_id} at {bank_id}.")
        print(
            "Open an account first, for example: "
            f"python3.11 mini.py open-account {human_id} {bank_id} {currency}"
        )
        return

    bank = world["banks"][bank_id]
    human = world["humans"][human_id]
    account = world["accounts"][account_id]

    loan_id = build_loan_id(world)

    loan = {
        "id": loan_id,
        "bank_id": bank_id,
        "borrower_human_id": human_id,
        "currency": currency,
        "original_principal": amount,
        "outstanding_principal": amount,
        "outstanding_interest": 0,
        "total_interest_accrued": 0,
        "written_off_principal": 0,
        "written_off_interest": 0,
        "status": "ACTIVE",
    }

    # Bank asset:
    # Alice now owes money to Bank 1.
    world["customer_loans"][loan_id] = loan
    bank["customer_loans"].append(loan_id)

    # Human liability:
    # Alice records that she owes this loan.
    human["loans"][loan_id] = amount

    # Bank liability / human asset:
    # Bank 1 credits Alice's deposit account.
    # This creates spendable bank money.
    account["booked_balance"] += amount

    save_world(world)

    print(f"Granted {amount} {currency} loan from {bank_id} to {human_id}.")
    print(f"Loan: {loan_id}")
    print(
        f"{loan_id} outstanding_principal: {loan['outstanding_principal']} {currency}"
    )
    print(f"{loan_id} outstanding_interest: {loan['outstanding_interest']} {currency}")
    print(f"{account_id} booked_balance: {account['booked_balance']} {currency}")
    print(f"{human_id} loans[{loan_id}]: {human['loans'][loan_id]} {currency}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "grant-loan",
        help="Grant a commercial bank loan and credit the borrower's deposit account.",
    )

    parser.add_argument("bank_id", help="Commercial bank ID, for example: bank1")

    parser.add_argument("human_id", help="Borrower human ID, for example: alice")

    parser.add_argument("currency", help="Currency, for example: EUR or USD")

    parser.add_argument("amount", type=int, help="Loan amount.")

    parser.set_defaults(func=run)
