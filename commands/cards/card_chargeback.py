from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    payment = world["payment_instructions"].get(args.payment_id)

    if payment is None:
        print(f"Payment does not exist: {args.payment_id}")
        return

    auth = world["card_authorizations"].get(payment.get("auth_id"))

    if auth is not None:
        account = world["accounts"].get(auth["account_id"])

        if account is not None:
            account["booked_balance"] += payment["amount"]

        auth["status"] = "CHARGED_BACK"

    payment["status"] = "CHARGED_BACK"
    save_world(world)
    print(f"Charged back card payment: {args.payment_id}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "card-chargeback", help="Reverse a disputed card payment."
    )

    parser.add_argument("payment_id")

    parser.set_defaults(func=run)
