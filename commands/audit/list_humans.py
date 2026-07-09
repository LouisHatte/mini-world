from argparse import Namespace, _SubParsersAction

from world import load_world


def run(args: Namespace) -> None:
    for human_id in load_world()["humans"]:
        print(human_id)


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser("list-humans", help="List humans.")

    parser.set_defaults(func=run)
