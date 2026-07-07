from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()

    loan_id = args.loan_id

    if loan_id not in world["customer_loans"]:
        print(f"Loan does not exist: {loan_id}")
        return

    loan = world["customer_loans"][loan_id]

    if loan["status"] == "REPAID":
        print(f"Loan is already repaid: {loan_id}")
        return

    if loan["status"] == "WRITTEN_OFF":
        print(f"Loan is already written off: {loan_id}")
        return

    if loan["status"] != "DEFAULTED":
        print(f"Loan must be DEFAULTED before write-off: {loan_id}")
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

    principal_written_off = loan["outstanding_principal"]
    interest_written_off = loan["outstanding_interest"]
    total_written_off = principal_written_off + interest_written_off

    if total_written_off <= 0:
        print(f"Loan has nothing left to write off: {loan_id}")
        return

    loan["written_off_principal"] += principal_written_off
    loan["written_off_interest"] += interest_written_off

    loan["outstanding_principal"] = 0
    loan["outstanding_interest"] = 0
    loan["status"] = "WRITTEN_OFF"

    human["loans"][loan_id] = 0

    bank["loan_loss_expense"][currency] = (
        bank["loan_loss_expense"].get(currency, 0) + total_written_off
    )

    bank["equity"][currency] = bank["equity"].get(currency, 0) - total_written_off

    save_world(world)

    print(f"Wrote off loan: {loan_id}")
    print(f"Principal written off: {principal_written_off} {currency}")
    print(f"Interest written off: {interest_written_off} {currency}")
    print(f"Total written off: {total_written_off} {currency}")
    print(
        f"{bank_id} loan_loss_expense[{currency}]: {bank['loan_loss_expense'][currency]} {currency}"
    )
    print(f"{bank_id} equity[{currency}]: {bank['equity'][currency]} {currency}")
    print(f"{human_id} loans[{loan_id}]: {human['loans'][loan_id]} {currency}")
    print(f"{loan_id} status: {loan['status']}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "write-off-loan",
        help="Write off a defaulted customer loan and reduce bank equity.",
    )

    parser.add_argument("loan_id", help="Loan ID, for example: loan_000001")

    parser.set_defaults(func=run)
