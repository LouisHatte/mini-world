from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()

    human_id = args.human_id

    if human_id in world["humans"]:
        print(f"Human already exists: {human_id}")
        return

    world["humans"][human_id] = {
        "id": human_id,
        "name": human_id,
        "cash_wallet": {},
        "bank_accounts": [],
        "loans": {},
    }

    save_world(world)

    print(f"Created human: {human_id}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser("create-human", help="Create a human actor.")

    parser.add_argument("human_id", help="Unique human ID, for example: alice")

    parser.set_defaults(func=run)
