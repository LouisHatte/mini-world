from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    hold_id = args.hold_id

    if hold_id not in world["holds"]:
        print(f"Hold does not exist: {hold_id}")
        return

    hold = world["holds"][hold_id]

    if hold["status"] != "HELD":
        print(f"Hold is not releasable: {hold_id} ({hold['status']})")
        return

    hold["status"] = "RELEASED"

    save_world(world)

    print(f"Released hold: {hold_id}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "release-hold", help="Cancel a held account reservation."
    )

    parser.add_argument("hold_id")

    parser.set_defaults(func=run)
