from argparse import _SubParsersAction

from commands.cash.issue_cash import register as register_issue_cash
from commands.cash.seed_cash import register as register_seed_cash
from commands.cash.supply_cash import register as register_supply_cash
from commands.cash.deposit_cash import register as register_deposit_cash
from commands.cash.withdraw_cash import register as register_withdraw_cash
from commands.cash.transfer_cash import register as register_transfer_cash
from commands.cash.move_cash import register as register_move_cash
from commands.cash.return_cash import register as register_return_cash
from commands.cash.destroy_cash import register as register_destroy_cash


def register_cash_commands(subparsers: _SubParsersAction) -> None:
    register_issue_cash(subparsers)
    register_seed_cash(subparsers)
    register_supply_cash(subparsers)
    register_deposit_cash(subparsers)
    register_withdraw_cash(subparsers)
    register_transfer_cash(subparsers)
    register_move_cash(subparsers)
    register_return_cash(subparsers)
    register_destroy_cash(subparsers)
