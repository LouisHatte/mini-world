from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()

    loan_id = args.loan_id
    amount = args.amount

    if amount <= 0:
        print("Amount must be greater than 0.")
        return

    if loan_id not in world["customer_loans"]:
        print(f"Loan does not exist: {loan_id}")
        return

    loan = world["customer_loans"][loan_id]

    if loan["status"] != "ACTIVE":
        print(f"Loan is not active: {loan_id}")
        return

    bank_id = loan["bank_id"]
    human_id = loan["borrower_human_id"]
    currency = loan["currency"]

    if bank_id not in world["banks"]:
        print(f"Bank does not exist: {bank_id}")
        return

    if human_id not in world["humans"]:
        print(f"Human does not exist: {human_id}")
        return

    bank = world["banks"][bank_id]
    human = world["humans"][human_id]

    # Bank loan asset increases.
    loan["outstanding_interest"] += amount
    loan["total_interest_accrued"] += amount

    # Human liability mirror increases.
    human["loans"][loan_id] = (
        loan["outstanding_principal"] + loan["outstanding_interest"]
    )

    save_world(world)

    print(f"Accrued {amount} {currency} interest on {loan_id}.")
    print(
        f"{loan_id} outstanding_principal: {loan['outstanding_principal']} {currency}"
    )
    print(f"{loan_id} outstanding_interest: {loan['outstanding_interest']} {currency}")
    print(f"{human_id} loans[{loan_id}]: {human['loans'][loan_id]} {currency}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "accrue-interest", help="Add interest to an active customer loan."
    )

    parser.add_argument("loan_id", help="Loan ID, for example: loan_000001")

    parser.add_argument("amount", type=int, help="Interest amount to accrue.")

    parser.set_defaults(func=run)
