import copy
from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    name = args.name

    if name not in world["snapshots"]:
        print(f"Snapshot does not exist: {name}")
        return

    snapshots = world["snapshots"]
    restored = copy.deepcopy(snapshots[name])
    restored["snapshots"] = snapshots
    save_world(restored)
    print(f"Loaded snapshot: {name}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "load-snapshot", help="Restore a previous snapshot."
    )

    parser.add_argument("name")

    parser.set_defaults(func=run)
