from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    payment = world["payment_instructions"].get(args.payment_id)

    if payment is None:
        print(f"Payment does not exist: {args.payment_id}")
        return

    payment["status"] = "RECIPIENT_REJECTED"
    save_world(world)
    print(f"Recipient rejected payment: {args.payment_id}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "reject-recipient", help="Reject an incoming payment at recipient bank."
    )

    parser.add_argument("payment_id")

    parser.set_defaults(func=run)
