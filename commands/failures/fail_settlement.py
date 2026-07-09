from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    payment = world["payment_instructions"].get(args.payment_id)

    if payment is None:
        print(f"Payment does not exist: {args.payment_id}")
        return

    payment["status"] = "SETTLEMENT_FAILED"
    settlement_id = payment.get("settlement_id")

    if settlement_id in world["settlements"]:
        world["settlements"][settlement_id]["status"] = "FAILED"

    save_world(world)
    print(f"Failed settlement for payment: {args.payment_id}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser("fail-settlement", help="Mark settlement as failed.")

    parser.add_argument("payment_id")

    parser.set_defaults(func=run)
