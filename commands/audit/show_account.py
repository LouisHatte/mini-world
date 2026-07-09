import json
from argparse import Namespace, _SubParsersAction

from world import load_world


def run(args: Namespace) -> None:
    world = load_world()
    account = world["accounts"].get(args.account_id)

    if account is None:
        print(f"Account does not exist: {args.account_id}")
        return

    print(json.dumps(account, indent=2, ensure_ascii=False))


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser("show-account", help="Show one account.")

    parser.add_argument("account_id")

    parser.set_defaults(func=run)
