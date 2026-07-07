from argparse import _SubParsersAction

from commands.world import register_world_commands
from commands.setup import register_setup_commands
from commands.cash import register_cash_commands
from commands.reserves import register_reserve_commands
from commands.payments import register_payment_commands
from commands.loans import register_loan_commands


def register_all_commands(subparsers: _SubParsersAction) -> None:
    register_world_commands(subparsers)
    register_setup_commands(subparsers)
    register_cash_commands(subparsers)
    register_reserve_commands(subparsers)
    register_payment_commands(subparsers)
    register_loan_commands(subparsers)
