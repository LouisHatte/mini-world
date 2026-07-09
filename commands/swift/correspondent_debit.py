from argparse import _SubParsersAction

from commands.swift.correspondent import adjust_correspondent


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "correspondent-debit", help="Debit a nostro balance through a correspondent."
    )

    parser.add_argument("correspondent_bank_id")
    parser.add_argument("owner_bank_id")
    parser.add_argument("currency")
    parser.add_argument("amount", type=int)

    parser.set_defaults(
        func=lambda args: adjust_correspondent(args, "nostro_balance", "debit")
    )
