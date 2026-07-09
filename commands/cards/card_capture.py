from argparse import Namespace, _SubParsersAction

from commands.common import next_id
from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    auth = world["card_authorizations"].get(args.auth_id)

    if auth is None:
        print(f"Authorization does not exist: {args.auth_id}")
        return

    if auth["status"] != "AUTHORIZED":
        print(f"Authorization is not capturable: {auth['status']}")
        return

    account = world["accounts"][auth["account_id"]]
    hold = world["holds"][auth["hold_id"]]
    account["booked_balance"] -= auth["amount"]
    hold["status"] = "CAPTURED"
    auth["status"] = "CAPTURED"
    payment_id = next_id(world["payment_instructions"], "payment")
    world["payment_instructions"][payment_id] = {
        "id": payment_id,
        "rail": "CARD",
        "status": "CAPTURED",
        "auth_id": args.auth_id,
        "merchant": auth["merchant"],
        "currency": auth["currency"],
        "amount": auth["amount"],
    }
    auth["payment_id"] = payment_id
    save_world(world)
    print(f"Captured card authorization: {args.auth_id}")
    print(f"Payment: {payment_id}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser("card-capture", help="Finalize a card payment.")

    parser.add_argument("auth_id")

    parser.set_defaults(func=run)
