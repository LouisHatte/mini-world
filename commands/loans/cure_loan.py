from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()

    loan_id = args.loan_id

    if loan_id not in world["customer_loans"]:
        print(f"Loan does not exist: {loan_id}")
        return

    loan = world["customer_loans"][loan_id]

    if loan["status"] == "ACTIVE":
        print(f"Loan is already active: {loan_id}")
        return

    if loan["status"] == "REPAID":
        print(f"Loan is already repaid: {loan_id}")
        return

    if loan["status"] == "WRITTEN_OFF":
        print(f"Loan is already written off: {loan_id}")
        return

    if loan["status"] != "DEFAULTED":
        print(f"Loan is not defaulted: {loan_id}")
        return

    loan["status"] = "ACTIVE"

    save_world(world)

    currency = loan["currency"]
    total_due = loan["outstanding_principal"] + loan["outstanding_interest"]

    print(f"Cured loan: {loan_id}")
    print(
        f"{loan_id} outstanding_principal: {loan['outstanding_principal']} {currency}"
    )
    print(f"{loan_id} outstanding_interest: {loan['outstanding_interest']} {currency}")
    print(f"{loan_id} total_due: {total_due} {currency}")
    print(f"{loan_id} status: {loan['status']}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "cure-loan", help="Move a defaulted customer loan back to active status."
    )

    parser.add_argument("loan_id", help="Loan ID, for example: loan_000001")

    parser.set_defaults(func=run)
