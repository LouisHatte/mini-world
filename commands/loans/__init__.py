from argparse import _SubParsersAction

from commands.loans.grant_loan import register as register_grant_loan
from commands.loans.repay_loan import register as register_repay_loan
from commands.loans.accrue_interest import register as register_accrue_interest
from commands.loans.default_loan import register as register_default_loan
from commands.loans.cure_loan import register as register_cure_loan
from commands.loans.write_off_loan import register as register_write_off_loan


def register_loan_commands(subparsers: _SubParsersAction) -> None:
    register_grant_loan(subparsers)
    register_repay_loan(subparsers)
    register_accrue_interest(subparsers)
    register_default_loan(subparsers)
    register_cure_loan(subparsers)
    register_write_off_loan(subparsers)
