from argparse import _SubParsersAction

from commands.reserves.lend_reserves import register as register_lend_reserves
from commands.reserves.settle_reserves import register as register_settle_reserves
from commands.reserves.repay_reserves import register as register_repay_reserves
from commands.reserves.check_reserves import register as register_check_reserves
from commands.reserves.reserve_shortfall import register as register_reserve_shortfall


def register_reserve_commands(subparsers: _SubParsersAction) -> None:
    register_lend_reserves(subparsers)
    register_settle_reserves(subparsers)
    register_repay_reserves(subparsers)
    register_check_reserves(subparsers)
    register_reserve_shortfall(subparsers)
