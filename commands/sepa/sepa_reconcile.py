from argparse import Namespace, _SubParsersAction

from world import load_world


def run(args: Namespace) -> None:
    world = load_world()
    payment_id = args.payment_id

    if payment_id not in world["payment_instructions"]:
        print(f"Payment does not exist: {payment_id}")
        return

    payment = world["payment_instructions"][payment_id]
    message = world["messages"].get(payment.get("message_id"))
    settlement = world["settlements"].get(payment.get("settlement_id"))

    print(f"SEPA payment: {payment_id}")
    print(f"Payment status: {payment['status']}")
    print(f"Message status: {message['status'] if message else 'MISSING'}")
    print(f"Settlement status: {settlement['status'] if settlement else 'MISSING'}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "sepa-reconcile",
        help="Compare SEPA message, clearing, settlement, and account states.",
    )

    parser.add_argument("payment_id")

    parser.set_defaults(func=run)
