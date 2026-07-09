from argparse import _SubParsersAction

from commands.cheques.bounce_cheque import register as register_bounce_cheque
from commands.cheques.clear_cheque import register as register_clear_cheque
from commands.cheques.deposit_cheque import register as register_deposit_cheque
from commands.cheques.reverse_cheque_credit import (
    register as register_reverse_cheque_credit,
)
from commands.cheques.settle_cheque import register as register_settle_cheque
from commands.cheques.write_cheque import register as register_write_cheque


def register_cheque_commands(subparsers: _SubParsersAction) -> None:
    register_write_cheque(subparsers)
    register_deposit_cheque(subparsers)
    register_clear_cheque(subparsers)
    register_settle_cheque(subparsers)
    register_bounce_cheque(subparsers)
    register_reverse_cheque_credit(subparsers)
