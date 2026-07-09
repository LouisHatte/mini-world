from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    auth = world["card_authorizations"].get(args.auth_id)

    if auth is None:
        print(f"Authorization does not exist: {args.auth_id}")
        return

    if auth["status"] != "AUTHORIZED":
        print(f"Authorization is not reversible: {auth['status']}")
        return

    world["holds"][auth["hold_id"]]["status"] = "RELEASED"
    auth["status"] = "REVERSED"
    save_world(world)
    print(f"Reversed card authorization: {args.auth_id}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "card-reverse", help="Cancel a card authorization."
    )

    parser.add_argument("auth_id")

    parser.set_defaults(func=run)
