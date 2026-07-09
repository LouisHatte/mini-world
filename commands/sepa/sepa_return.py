from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    payment_id = args.payment_id

    if payment_id not in world["payment_instructions"]:
        print(f"Payment does not exist: {payment_id}")
        return

    payment = world["payment_instructions"][payment_id]
    payment["status"] = "RETURNED"
    payment["return_reason"] = args.reason

    sender_account = world["accounts"].get(payment["sender_account_id"])

    if sender_account is not None:
        sender_account["booked_balance"] += payment["amount"]

    save_world(world)

    print(f"Returned SEPA payment: {payment_id}")
    print(f"Reason: {args.reason}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser("sepa-return", help="Return or reject a SEPA payment.")

    parser.add_argument("payment_id")
    parser.add_argument("reason")

    parser.set_defaults(func=run)
