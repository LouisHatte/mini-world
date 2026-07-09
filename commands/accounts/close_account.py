from argparse import _SubParsersAction

from commands.accounts.status import set_account_status


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "close-account", help="Close a zero-balance bank account."
    )

    parser.add_argument("account_id")

    parser.set_defaults(func=lambda args: set_account_status(args, "CLOSED"))
