import json
from argparse import Namespace, _SubParsersAction

from world import load_world


def run(args: Namespace) -> None:
    world = load_world()
    payment = world["payment_instructions"].get(args.payment_id)

    if payment is None:
        print(f"Payment does not exist: {args.payment_id}")
        return

    print(json.dumps(payment, indent=2, ensure_ascii=False))


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "show-payment", help="Show one payment state machine."
    )

    parser.add_argument("payment_id")

    parser.set_defaults(func=run)
