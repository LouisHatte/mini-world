from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    payment_id = args.payment_id

    if payment_id not in world["payment_instructions"]:
        print(f"Payment does not exist: {payment_id}")
        return

    payment = world["payment_instructions"][payment_id]
    payment["status"] = "CLEARED"

    save_world(world)

    print(f"STEP2 cleared payment: {payment_id}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "step2-clear", help="Clear a SEPA payment in STEP2."
    )

    parser.add_argument("payment_id")

    parser.set_defaults(func=run)
