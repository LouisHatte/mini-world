from argparse import Namespace, _SubParsersAction

from world import WORLD_FILE, create_empty_world, save_world, world_exists


def run(args: Namespace) -> None:
    if world_exists() and not args.reset:
        print(f"World already exists: {WORLD_FILE}")
        print("Use --reset if you want to overwrite it.")
        return

    world = create_empty_world()
    save_world(world)

    if args.reset:
        print(f"Reset world: {WORLD_FILE}")
    else:
        print(f"Initialized empty world: {WORLD_FILE}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser("init", help="Create an empty banking world.")

    parser.add_argument(
        "--reset", action="store_true", help="Overwrite the existing world file."
    )

    parser.set_defaults(func=run)
