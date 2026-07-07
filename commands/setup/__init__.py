from argparse import _SubParsersAction

from commands.setup.create_central_bank import register as register_create_central_bank
from commands.setup.create_bank import register as register_create_bank
from commands.setup.create_human import register as register_create_human
from commands.setup.open_account import register as register_open_account


def register_setup_commands(subparsers: _SubParsersAction) -> None:
    register_create_central_bank(subparsers)
    register_create_bank(subparsers)
    register_create_human(subparsers)
    register_open_account(subparsers)
