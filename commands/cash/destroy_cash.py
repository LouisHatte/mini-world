from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()

    central_bank_id = args.central_bank_id
    amount = args.amount

    if amount <= 0:
        print("Amount must be greater than 0.")
        return

    if central_bank_id not in world["central_banks"]:
        print(f"Central bank does not exist: {central_bank_id}")
        return

    central_bank = world["central_banks"][central_bank_id]
    currency = central_bank["currency"]

    if central_bank["cash_vault"] < amount:
        print(
            f"Not enough physical cash in {central_bank_id}'s vault. "
            f"Available: {central_bank['cash_vault']} {currency}"
        )
        return

    central_bank["cash_vault"] -= amount
    central_bank["cash_issued"] -= amount

    save_world(world)

    print(f"Destroyed {amount} {currency} cash at {central_bank_id}.")
    print(f"{central_bank_id} cash_issued: {central_bank['cash_issued']} {currency}")
    print(f"{central_bank_id} cash_vault: {central_bank['cash_vault']} {currency}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "destroy-cash", help="Destroy physical cash held in a central bank vault."
    )

    parser.add_argument("central_bank_id", help="Central bank ID, for example: ecb")

    parser.add_argument("amount", type=int, help="Amount of physical cash to destroy.")

    parser.set_defaults(func=run)
