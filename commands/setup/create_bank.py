from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()

    bank_id = args.bank_id

    if bank_id in world["banks"]:
        print(f"Bank already exists: {bank_id}")
        return

    world["banks"][bank_id] = {
        "id": bank_id,
        "name": bank_id,
        "cash_vault": {},
        "reserve_balances": {},
        "loans_from_central_banks": {},
        "customer_accounts": [],
        "customer_loans": [],
        "interest_income": {},
        "loan_loss_expense": {},
        "equity": {},
    }

    save_world(world)

    print(f"Created bank: {bank_id}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser("create-bank", help="Create a commercial bank.")

    parser.add_argument("bank_id", help="Unique bank ID, for example: bank1")

    parser.set_defaults(func=run)
