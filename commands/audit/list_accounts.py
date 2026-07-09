from argparse import Namespace, _SubParsersAction

from world import load_world


def run(args: Namespace) -> None:
    world = load_world()

    for account_id, account in world["accounts"].items():
        print(
            f"{account_id} {account['owner_human_id']} {account['bank_id']} "
            f"{account['currency']} {account['booked_balance']} {account['status']}"
        )


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser("list-accounts", help="List accounts.")

    parser.set_defaults(func=run)
