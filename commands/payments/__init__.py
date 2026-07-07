from argparse import _SubParsersAction

from commands.payments.internal_transfer import register as register_internal_transfer


def register_payment_commands(subparsers: _SubParsersAction) -> None:
    register_internal_transfer(subparsers)
