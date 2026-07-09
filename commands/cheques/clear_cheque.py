from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    cheque = world["cheques"].get(args.cheque_id)

    if cheque is None:
        print(f"Cheque does not exist: {args.cheque_id}")
        return

    cheque["status"] = "CLEARED"
    save_world(world)
    print(f"Cleared cheque: {args.cheque_id}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser("clear-cheque", help="Clear an interbank cheque.")

    parser.add_argument("cheque_id")

    parser.set_defaults(func=run)
