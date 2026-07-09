from argparse import Namespace, _SubParsersAction

from commands.common import next_id, require_active_account
from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    currency = args.currency.upper()
    account_id, account = require_active_account(
        world, args.drawer_human_id, args.drawer_bank_id, currency
    )

    if account_id is None or account is None:
        print("Drawer account does not exist.")
        return

    cheque_id = next_id(world["cheques"], "cheque")
    world["cheques"][cheque_id] = {
        "id": cheque_id,
        "drawer_human_id": args.drawer_human_id,
        "drawer_bank_id": args.drawer_bank_id,
        "drawer_account_id": account_id,
        "payee_human_id": args.payee_human_id,
        "currency": currency,
        "amount": args.amount,
        "status": "WRITTEN",
    }
    save_world(world)
    print(f"Wrote cheque: {cheque_id}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser("write-cheque", help="Write a cheque.")

    parser.add_argument("drawer_human_id")
    parser.add_argument("drawer_bank_id")
    parser.add_argument("payee_human_id")
    parser.add_argument("currency")
    parser.add_argument("amount", type=int)

    parser.set_defaults(func=run)
