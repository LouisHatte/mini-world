import json
from argparse import Namespace, _SubParsersAction

from world import load_world


def run(args: Namespace) -> None:
    print(json.dumps(load_world(), indent=2, ensure_ascii=False))


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser("show-json", help="Pretty-print raw world JSON.")

    parser.set_defaults(func=run)
