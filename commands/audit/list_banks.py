from argparse import Namespace, _SubParsersAction

from world import load_world


def run(args: Namespace) -> None:
    for bank_id in load_world()["banks"]:
        print(bank_id)


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser("list-banks", help="List banks.")

    parser.set_defaults(func=run)
