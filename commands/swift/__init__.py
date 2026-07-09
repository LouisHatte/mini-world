from argparse import _SubParsersAction

from commands.swift.correspondent_credit import register as register_correspondent_credit
from commands.swift.correspondent_debit import register as register_correspondent_debit
from commands.swift.create_correspondent_account import (
    register as register_create_correspondent_account,
)
from commands.swift.deduct_fee import register as register_deduct_fee
from commands.swift.swift_message import register as register_swift_message
from commands.swift.swift_reconcile import register as register_swift_reconcile
from commands.swift.swift_transfer import register as register_swift_transfer


def register_swift_commands(subparsers: _SubParsersAction) -> None:
    register_create_correspondent_account(subparsers)
    register_swift_transfer(subparsers)
    register_swift_message(subparsers)
    register_correspondent_debit(subparsers)
    register_correspondent_credit(subparsers)
    register_deduct_fee(subparsers)
    register_swift_reconcile(subparsers)
