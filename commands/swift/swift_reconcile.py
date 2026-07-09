from argparse import Namespace, _SubParsersAction

from world import load_world


def run(args: Namespace) -> None:
    world = load_world()
    payment = world["payment_instructions"].get(args.payment_id)

    if payment is None:
        print(f"Payment does not exist: {args.payment_id}")
        return

    print(f"SWIFT payment: {args.payment_id}")
    print(f"Status: {payment['status']}")
    print(f"Messages: {', '.join(payment.get('message_ids', [])) or 'NONE'}")
    print(f"Fees: {payment.get('fees', 0)} {payment['currency']}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "swift-reconcile", help="Check SWIFT message chain vs ledger movements."
    )

    parser.add_argument("payment_id")

    parser.set_defaults(func=run)
