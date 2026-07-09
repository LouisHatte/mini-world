from argparse import _SubParsersAction

from commands.audit.list_accounts import register as register_list_accounts
from commands.audit.list_banks import register as register_list_banks
from commands.audit.list_humans import register as register_list_humans
from commands.audit.load_snapshot import register as register_load_snapshot
from commands.audit.show_account import register as register_show_account
from commands.audit.show_balance_sheet import register as register_show_balance_sheet
from commands.audit.show_json import register as register_show_json
from commands.audit.show_ledger import register as register_show_ledger
from commands.audit.show_payment import register as register_show_payment
from commands.audit.snapshot import register as register_snapshot


def register_audit_commands(subparsers: _SubParsersAction) -> None:
    register_show_json(subparsers)
    register_show_ledger(subparsers)
    register_show_balance_sheet(subparsers)
    register_show_account(subparsers)
    register_show_payment(subparsers)
    register_list_banks(subparsers)
    register_list_humans(subparsers)
    register_list_accounts(subparsers)
    register_snapshot(subparsers)
    register_load_snapshot(subparsers)
