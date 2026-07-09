from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    cheque = world["cheques"].get(args.cheque_id)

    if cheque is None:
        print(f"Cheque does not exist: {args.cheque_id}")
        return

    account = world["accounts"].get(cheque.get("payee_account_id"))

    if account is not None:
        account["booked_balance"] -= cheque["amount"]

    cheque["status"] = "REVERSED"
    save_world(world)
    print(f"Reversed cheque credit: {args.cheque_id}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "reverse-cheque-credit", help="Reverse provisional cheque credit."
    )

    parser.add_argument("cheque_id")

    parser.set_defaults(func=run)
