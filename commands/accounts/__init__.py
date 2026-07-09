from argparse import _SubParsersAction

from commands.accounts.capture_hold import register as register_capture_hold
from commands.accounts.close_account import register as register_close_account
from commands.accounts.freeze_account import register as register_freeze_account
from commands.accounts.place_hold import register as register_place_hold
from commands.accounts.release_hold import register as register_release_hold
from commands.accounts.unfreeze_account import register as register_unfreeze_account


def register_account_commands(subparsers: _SubParsersAction) -> None:
    register_place_hold(subparsers)
    register_release_hold(subparsers)
    register_capture_hold(subparsers)
    register_close_account(subparsers)
    register_freeze_account(subparsers)
    register_unfreeze_account(subparsers)
