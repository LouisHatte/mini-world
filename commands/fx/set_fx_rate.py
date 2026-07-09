from argparse import _SubParsersAction

from commands.fx.create_fx_market import run


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser("set-fx-rate", help="Update an FX rate.")

    parser.add_argument("from_currency")
    parser.add_argument("to_currency")
    parser.add_argument("rate", type=float)

    parser.set_defaults(func=run)
