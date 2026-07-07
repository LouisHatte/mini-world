from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()

    central_bank_id = args.central_bank_id
    currency = args.currency.upper()

    if central_bank_id in world["central_banks"]:
        print(f"Central bank already exists: {central_bank_id}")
        return

    world["central_banks"][central_bank_id] = {
        "id": central_bank_id,
        "name": central_bank_id,
        "currency": currency,
        "cash_issued": 0,
        "cash_vault": 0,
        "reserve_accounts": {},
        "loans_to_banks": {},
        "securities": {},
    }

    save_world(world)

    print(f"Created central bank: {central_bank_id} ({currency})")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser("create-central-bank", help="Create a central bank.")

    parser.add_argument(
        "central_bank_id", help="Unique central bank ID, for example: ecb"
    )

    parser.add_argument(
        "currency", help="Currency issued by this central bank, for example: EUR or USD"
    )

    parser.set_defaults(func=run)
