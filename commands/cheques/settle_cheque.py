from argparse import Namespace, _SubParsersAction

from commands.common import assert_account_can_debit
from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    cheque = world["cheques"].get(args.cheque_id)

    if cheque is None:
        print(f"Cheque does not exist: {args.cheque_id}")
        return

    drawer_account = world["accounts"][cheque["drawer_account_id"]]
    error = assert_account_can_debit(world, drawer_account, cheque["amount"])

    if error is not None:
        print(error)
        cheque["status"] = "BOUNCED"
        save_world(world)
        return

    drawer_account["booked_balance"] -= cheque["amount"]
    cheque["status"] = "SETTLED"
    save_world(world)
    print(f"Settled cheque: {args.cheque_id}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser("settle-cheque", help="Settle a cheque.")

    parser.add_argument("cheque_id")

    parser.set_defaults(func=run)
