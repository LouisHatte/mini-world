from argparse import _SubParsersAction

from commands.world import register_world_commands
from commands.setup import register_setup_commands
from commands.cash import register_cash_commands
from commands.reserves import register_reserve_commands
from commands.payments import register_payment_commands
from commands.loans import register_loan_commands
from commands.accounts import register_account_commands
from commands.sepa import register_sepa_commands
from commands.swift import register_swift_commands
from commands.securities import register_security_commands
from commands.cheques import register_cheque_commands
from commands.cards import register_card_commands
from commands.fx import register_fx_commands
from commands.failures import register_failure_commands
from commands.audit import register_audit_commands


def register_all_commands(subparsers: _SubParsersAction) -> None:
    register_world_commands(subparsers)
    register_setup_commands(subparsers)
    register_cash_commands(subparsers)
    register_reserve_commands(subparsers)
    register_payment_commands(subparsers)
    register_loan_commands(subparsers)
    register_account_commands(subparsers)
    register_sepa_commands(subparsers)
    register_swift_commands(subparsers)
    register_security_commands(subparsers)
    register_cheque_commands(subparsers)
    register_card_commands(subparsers)
    register_fx_commands(subparsers)
    register_failure_commands(subparsers)
    register_audit_commands(subparsers)
