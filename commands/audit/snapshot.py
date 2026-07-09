import copy
from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    name = args.name
    snapshot_world = copy.deepcopy(world)
    snapshot_world.pop("snapshots", None)
    world["snapshots"][name] = snapshot_world
    save_world(world)
    print(f"Saved snapshot: {name}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser("snapshot", help="Save current world state.")

    parser.add_argument("name")

    parser.set_defaults(func=run)
