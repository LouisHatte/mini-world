from argparse import _SubParsersAction

from commands.accounts.status import set_account_status


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "freeze-account", help="Block activity on a bank account."
    )

    parser.add_argument("account_id")

    parser.set_defaults(func=lambda args: set_account_status(args, "FROZEN"))
